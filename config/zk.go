package config

type ZooKeeperConfig struct {
	Servers        string `yaml:"servers" mapstructure:"servers"`
	Path           string `yaml:"path" mapstructure:"path"`
	Timeout        int    `yaml:"timeout" mapstructure:"timeout"`
	MaxRetries     int    `yaml:"max_retries" mapstructure:"max_retries"`
	InitialBackoff int    `yaml:"initial_backoff" mapstructure:"initial_backoff"`
	MaxBackoff     int    `yaml:"max_backoff" mapstructure:"max_backoff"`
}
