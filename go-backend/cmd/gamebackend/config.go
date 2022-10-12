package main

import "fmt"

const (
	ModeLocalHost = "localhost"
)

type config struct {
	Port                uint   `yaml:"port"`
	Host                string `yaml:"host"`
	Mode                string `yaml:"mode"`
	KeyDir              string `yaml:"keyDir"`
	DBFile              string `yaml:"dbFile"`
	AccountsPort        uint   `yaml:"accountsPort"`
	AccountsHost        string `yaml:"accountsHost"`
	CharactersPort      uint   `yaml:"charactersPort"`
	CharactersHost      string `yaml:"charactersHost"`
	AgonesKeyFile       string `yaml:"agonesKeyFile"`
	AgonesCertFile      string `yaml:"agonesCertFile"`
	AgonesCaCertFile    string `yaml:"agonesCaCertFile"`
	AgonesNamespace     string `yaml:"agonesNamespace"`
	AgonesAllocatorHost string `yaml:"agonesAllocatorHost"`
	AgonesAllocatorPort uint   `yaml:"agonesAllocatorPort"`
}

func (c *config) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func (c *config) AccountsAddress() string {
	return fmt.Sprintf("%s:%d", c.AccountsHost, c.AccountsPort)
}

func (c *config) CharactersAddress() string {
	return fmt.Sprintf("%s:%d", c.CharactersHost, c.CharactersPort)
}

func (c *config) AgonesAllocatorAddress() string {
	return fmt.Sprintf("%s:%d", c.AgonesAllocatorHost, c.AgonesAllocatorPort)
}
