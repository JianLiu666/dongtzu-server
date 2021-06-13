package config

type arangoConfig struct {
	Addr             string `yaml:"addr,omitempty"`
	DBName           string `yaml:"dbName,omitempty"`
	Username         string `yaml:"username,omitempty"`
	Password         string `yaml:"password,omitempty"`
	RetryCount       int    `yaml:"retryCount,omitempty"`
	RetryDurationMin int64  `yaml:"retryDurationMin,omitempty"` // TODO: not using, check this
	RetryDurationMax int64  `yaml:"retryDurationMax,omitempty"` // TODO: not using, check this
}
