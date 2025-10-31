package config

type root struct {
	App App `env:"app"`
}

type App struct {
	Name      string `env:"name"`
	Port      int    `env:"port"`
	Env       string `env:"env"`
	DebugMode bool   `env:"debug_mode"`
}
