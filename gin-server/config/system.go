package config

type System struct {
	Path       string `mapstructure:"path" json:"path" yaml:"path"`
	RemotePath string `mapstructure:"remotepath" json:"remotepath" yaml:"remotepath"`
}
