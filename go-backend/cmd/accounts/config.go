package main

import "fmt"

type config struct {
	Port   uint   `yaml:"port"`
	Host   string `yaml:"host"`
	Mode   string `yaml:"mode"`
	KeyDir string `yaml:"keyDir"`
	DBFile string `yaml:"dbFile"`
}

func (c *config) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
