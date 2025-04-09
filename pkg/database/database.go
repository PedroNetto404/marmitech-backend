package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/PedroNetto404/marmitech-backend/internal/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/source/file"
	"github.com/Masterminds/squirrel"
)

type Db struct {
	Instance *sql.DB
	dns      string
	queryTimeout int
	retryAttempts int
	retryDelay   int
}

func NewDb(
	queryTimeout,
	retryAttempts,
	retryDelay int,
) (*Db, error) {
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

	return &Db{
		dns: dns,
		retryAttempts: retryAttempts,
		retryDelay: retryDelay,
		queryTimeout: queryTimeout,
	}, nil
}

func (d *Db) Close() error {
	return d.Instance.Close()
}

const migrationsPath = "file://migrations"

func (d *Db) Migrate() error {
	driver, err := mysql.WithInstance(d.Instance, &mysql.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(migrationsPath, "mysql", driver)
	if err != nil {
		return fmt.Errorf("failed to initialize migrations: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	fmt.Println("Migrations applied successfully!")
	return nil
}

type Mapper func(rows *sql.Rows) (any, error)

// Query é um wrapper que possui timeout e retry
func (d *Db) QueryPage(
	builder sq.SelectBuilder,
	mapper Mapper,
	limit int,
	offset int,
) (any, error) {
	builder = builder.Limit(uint64(limit)).Offset(uint64(offset))

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}
	
	rows, err := d.Instance.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]any, 0)
	for rows.Next() {
		item, err := mapper(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}
