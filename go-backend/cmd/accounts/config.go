package main

import "fmt"

type config struct {
	Port   uint   `yaml:"accounts.port"`
	Host   string `yaml:"accounts.host"`
	Mode   string `yaml:"accounts.mode"`
	KeyDir string `yaml:"accounts.keyDir"`
	DBFile string `yaml:"accounts.dbFile"`
}

func (c *config) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
