package config

// Config ...
type Config struct {
	Translators Translators `yaml:"translators"`
	MySQL       MySQL       `yaml:"mysql"`
}

// Translators list
type Translators struct {
	Yandex Yandex `yaml:"yandex"`
}

// Yandex provider config
type Yandex struct {
	ApiVer string `yaml:"api-ver"`
	ApiKey string `yaml:"api-key"`
}

// MySQL storage config
type MySQL struct {
	Address  string `yaml:"address"`
	DBName   string `yaml:"db-name"`
	User     string `yaml:"user"`
	Password string `yaml:"pass"`
}
