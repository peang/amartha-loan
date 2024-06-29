package configs

import (
	"database/sql"
	"fmt"
	"runtime"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

func LoadDatabase(c *Config) *bun.DB {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=%s",
		c.SQLUsername,
		c.SQLPassword,
		c.SQLHost,
		c.SQLDatabase,
		c.SQLSSL,
	)
	pgconn := pgdriver.NewConnector(
		pgdriver.WithDSN(dsn),
	)

	sqldb := sql.OpenDB(pgconn)
	if c.Env == "production" {
		maxOpenConns := 4 * runtime.GOMAXPROCS(0)
		sqldb.SetMaxOpenConns(maxOpenConns)
		sqldb.SetMaxIdleConns(maxOpenConns)
	}

	db := bun.NewDB(sqldb, pgdialect.New())

	if c.Env != "production" {
		db.AddQueryHook(bundebug.NewQueryHook(
			bundebug.WithVerbose(true),
			bundebug.FromEnv("BUNDEBUG"),
		))
	}

	return db
}
