package config

import "time"

type root struct {
	AppRestApi   AppRestApi   `env:"app_rest_api"`
	AppGrpcApi   AppGrpcApi   `env:"app_grpc_api"`
	AppScheduler AppScheduler `env:"app_scheduler"`
	Database     Database     `env:"database"`
}

type AppRestApi struct {
	Name      string `env:"name"`
	Env       string `env:"env"`
	DebugMode bool   `env:"debug_mode"`
	Port      int    `env:"port"`
	Pprof     Pprof  `env:"pprof"`
}

type AppGrpcApi struct {
	Name      string `env:"name"`
	Env       string `env:"env"`
	DebugMode bool   `env:"debug_mode"`
	Port      int    `env:"port"`
	Pprof     Pprof  `env:"pprof"`
}

type AppScheduler struct {
	Name                string `env:"name"`
	Env                 string `env:"env"`
	DebugMode           bool   `env:"debug_mode"`
	HealthCheckInterval string `env:"healthcheck_interval"`
	Pprof               Pprof  `env:"pprof"`
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
