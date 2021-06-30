package config

type Server struct {
	JWT     JWT     `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	Zap     Zap     `mapstructure:"zap" json:"zap" yaml:"zap"`
	Redis   Redis   `mapstructure:"redis" json:"redis" yaml:"redis"`
	System  System  `mapstructure:"system" json:"system" yaml:"system"`
	Mysql   Mysql   `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Mongo   Mongo   `mapstructure:"mongo" json:"mongo" yaml:"mongo"`
	Ssh     Ssh     `mapstructure:"ssh" json:"ssh" yaml:"ssh"`
}
