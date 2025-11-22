package config

import "time"

type root struct {
	AppRestApi   AppRestApi   `env:"app_rest_api"`
	AppGrpcApi   AppGrpcApi   `env:"app_grpc_api"`
	AppScheduler AppScheduler `env:"app_scheduler"`
}

type AppRestApi struct {
	Name      string   `env:"name"`
	Env       string   `env:"env"`
	DebugMode bool     `env:"debug_mode"`
	Port      int      `env:"port"`
	Gin       Gin      `env:"gin"`
	Pprof     Pprof    `env:"pprof"`
	Database  Database `env:"database"`
}

type AppGrpcApi struct {
	Name      string   `env:"name"`
	Env       string   `env:"env"`
	DebugMode bool     `env:"debug_mode"`
	Port      int      `env:"port"`
	Pprof     Pprof    `env:"pprof"`
	Database  Database `env:"database"`
}

type AppScheduler struct {
	Name                string   `env:"name"`
	Env                 string   `env:"env"`
	DebugMode           bool     `env:"debug_mode"`
	HealthCheckInterval string   `env:"healthcheck_interval"`
	Pprof               Pprof    `env:"pprof"`
	Database            Database `env:"database"`
}

type Pprof struct {
	Enable      bool   `env:"enable"`
	Port        int    `env:"port"`
	StaticToken string `env:"static_token"`
}

type Database struct {
	DSN             string        `env:"dsn"`
	MaxOpenConns    int           `env:"max_open_conns"`
	MaxIdleConns    int           `env:"max_idle_conns"`
	ConnMaxLifetime time.Duration `env:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `env:"conn_max_idle_time"`
}

type Gin struct {
	Mode                 string  `env:"mode"`
	DisableConsoleColor  bool    `env:"disable_console_color"`
	DisableDefaultWriter bool    `env:"disable_default_writer"`
	DisableErrorWriter   bool    `env:"disable_error_writer"`
	Cors                 GinCors `env:"cors"`
	UseOtel              bool    `env:"use_otel"`
}

type GinCors struct {
	AllowOrigins     []string `env:"allow_origins"`
	AllowMethods     []string `env:"allow_methods"`
	AllowHeaders     []string `env:"allow_headers"`
	AllowCredentials bool     `env:"allow_credentials"`
	ExposeHeaders    []string `env:"expose_headers"`
	MaxAge           int      `env:"max_age"`
}
