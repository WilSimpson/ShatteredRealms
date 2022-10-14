package config

import "fmt"

type UptraceConfig struct {
	Host  string `yaml:"host"`
	Port  uint   `yaml:"port"`
	Id    string `yaml:"id"`
	Token string `yaml:"token"`
}

func (c *UptraceConfig) DSN() string {
	return fmt.Sprintf("http://%s@%s:%d/%s", c.Token, c.Host, c.Port, c.Id)
}
