package main

import (
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/config"
	log "github.com/sirupsen/logrus"
)

type appConfig struct {
	Chat     config.Server `yaml:"chat"`
	Accounts config.Server `yaml:"accounts"`
	KeyDir   string        `yaml:"keyDir"`
	DBFile   string        `yaml:"dbFile"`
}

var (
	conf = &appConfig{
		Chat: config.Server{
			Port:     8180,
			Host:     "",
			Mode:     "development",
			LogLevel: log.InfoLevel,
		},
		Accounts: config.Server{
			Port: 8080,
			Host: "",
		},
		KeyDir: "/etc/sro/auth",
		DBFile: "/etc/sro/db.yaml",
	}
)

func main() {

}
