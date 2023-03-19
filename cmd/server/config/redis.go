package config

type Redis struct {
	Addr        string `yaml:"addr"`
	Password    string `yaml:"password"`
	SessionSize int    `yaml:"session-size"`
}
