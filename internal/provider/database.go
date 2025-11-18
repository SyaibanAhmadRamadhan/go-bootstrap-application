package provider

import (
	"database/sql"
	"log/slog"

	"erp-directory-service/internal/config"

	"github.com/SyaibanAhmadRamadhan/go-foundation-kit/databases/sqlx"
	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	rdbms sqlx.RDBMS
	sqlDB *sql.DB
}

func NewDB(debugMode bool) (DB, error) {
	cfg := config.GetDatabase()

	db, err := sql.Open("mysql", cfg.DSN)
	if err != nil {
		return DB{}, err
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	if err := db.Ping(); err != nil {
		return DB{}, err
	}

	slog.Info("database connection established",
		"max_open_conns", cfg.MaxOpenConns,
		"max_idle_conns", cfg.MaxIdleConns,
		"conn_max_lifetime", cfg.ConnMaxLifetime,
		"conn_max_idle_time", cfg.ConnMaxIdleTime,
	)

	rdbms := sqlx.NewRDBMS(db,
		sqlx.UseDebug(debugMode),
	)

	return DB{
		rdbms: rdbms,
		sqlDB: db,
	}, nil
}

func (d *DB) RDBMS() sqlx.RDBMS {
	return d.rdbms
}

func (d *DB) SQLDB() *sql.DB {
	return d.sqlDB
}

func (d *DB) Close() error {
	slog.Info("Close db Connection")
	return d.rdbms.Close()
}
