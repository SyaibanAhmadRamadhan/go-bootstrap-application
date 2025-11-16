package config

import "time"

type root struct {
	App      App      `env:"app"`
	Database Database `env:"database"`
}

type App struct {
	Name      string `env:"name"`
	Port      int    `env:"port"`
	Env       string `env:"env"`
	DebugMode bool   `env:"debug_mode"`
}

type Database struct {
	DSN             string        `env:"dsn"`
	MaxOpenConns    int           `env:"max_open_conns"`
	MaxIdleConns    int           `env:"max_idle_conns"`
	ConnMaxLifetime time.Duration `env:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `env:"conn_max_idle_time"`
}
