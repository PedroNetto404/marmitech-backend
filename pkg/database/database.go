package database

import (
	"database/sql"
	"fmt"
	"sort"
	"strings"

	"github.com/PedroNetto404/marmitech-backend/internal/config"
	"github.com/PedroNetto404/marmitech-backend/pkg/types"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate"
)

type Database struct {
	Instance *sql.DB
}

var dns string

func New() (*Database, error) {
	dns = fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Env.MySqlUser,
		config.Env.MySqlPass,
		config.Env.MySqlHost,
		config.Env.MySqlPort,
		config.Env.MySqlDatabase,
	)

	db, err := sql.Open("mysql", dns)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Database{Instance: db}, nil
}

func (d *Database) Close() error {
	return d.Instance.Close()
}

const migrationsPath = "file://migrations"

func (d *Database) Migrate() error {
	m, err := migrate.New(migrationsPath, dns)
	if err != nil {
		return fmt.Errorf("failed to initialize migrations: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	fmt.Println("Migrations applied successfully!")
	return nil
}

func (d *Database) BuildFindQuery(baseQuery, countQuery string, args types.FindArgs) (finalQuery string, totalCount int, params []any) {
	filterClauses := make([]string, 0)
	params = make([]any, 0)

	filterKeys := make([]string, 0, len(args.Filter))
	for k := range args.Filter {
		filterKeys = append(filterKeys, k)
	}
	sort.Strings(filterKeys)

	for _, key := range filterKeys {
		filterClauses = append(filterClauses, fmt.Sprintf("%s = ?", key))
		params = append(params, args.Filter[key])
	}

	if len(filterClauses) > 0 {
		filterSQL := " AND " + strings.Join(filterClauses, " AND ")
		baseQuery += filterSQL
		countQuery += filterSQL
	}

	if args.SortBy != "" {
		baseQuery += fmt.Sprintf(" ORDER BY %s", args.SortBy)
		if args.SortAsc {
			baseQuery += " ASC"
		} else {
			baseQuery += " DESC"
		}
	}

	if args.Limit > 0 {
		baseQuery += " LIMIT ?"
		params = append(params, args.Limit)
	}
	if args.Offset > 0 {
		baseQuery += " OFFSET ?"
		params = append(params, args.Offset)
	}

	row := d.Instance.QueryRow(countQuery, params[:len(params)-2]...) 
	if err := row.Scan(&totalCount); err != nil {
		totalCount = 0 
	}

	return baseQuery, totalCount, params
}
