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

type Db struct {
	Instance *sql.DB
	dns      string
}

func NewMySQLDatabase() (*Db, error) {
	dns := fmt.Sprintf(
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

	return &Db{Instance: db, dns: dns}, nil
}

func (d *Db) Close() error {
	return d.Instance.Close()
}

const migrationsPath = "file://migrations"

func (d *Db) Migrate() error {
	m, err := migrate.New(migrationsPath, d.dns)
	if err != nil {
		return fmt.Errorf("failed to initialize migrations: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	fmt.Println("Migrations applied successfully!")
	return nil
}

func (d *Db) ConstructFindQuery(baseQuery, countQuery string, args types.FindArgs) (finalQuery string, totalCount int, params []any) {
	filterClauses := make([]string, 0)
	filterParams := make([]any, 0)

	filterKeys := make([]string, 0, len(args.Filter))
	for k := range args.Filter {
		filterKeys = append(filterKeys, k)
	}
	sort.Strings(filterKeys)

	for _, key := range filterKeys {
		filterClauses = append(filterClauses, fmt.Sprintf("%s = ?", key))
		filterParams = append(filterParams, args.Filter[key])
	}

	if len(filterClauses) > 0 {
		filterSQL := " WHERE " + strings.Join(filterClauses, " AND ")
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
		filterParams = append(filterParams, args.Limit)
	}
	if args.Offset > 0 {
		baseQuery += " OFFSET ?"
		filterParams = append(filterParams, args.Offset)
	}

	row := d.Instance.QueryRow(countQuery, filterParams...)
	if err := row.Scan(&totalCount); err != nil {
		totalCount = 0
	}

	return baseQuery, totalCount, filterParams
}
