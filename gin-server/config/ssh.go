package config

type Ssh struct {
	User    string `mapstructure:"user" json:"user" yaml:"user"`
	Port    int    `mapstructure:"port" json:"port" yaml:"port"`
	KeyPath string `mapstructure:"keypath" json:"keypath" yaml:"keypath"`
}
