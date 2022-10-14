package config

import (
	"fmt"
	"github.com/WilSimpson/ShatteredRealms/go-backend/pkg/repository"
	log "github.com/sirupsen/logrus"
)

const (
	ModeProduction  ServerMode = "production"
	ModeDebug       ServerMode = "debug"
	ModeVerbose     ServerMode = "verbose"
	ModeDevelopment ServerMode = "development"
)

var (
	AllModes = []ServerMode{ModeProduction, ModeDebug, ModeVerbose, ModeDevelopment}
)

type ServerMode string

type Server struct {
	Local    ServerAddress           `yaml:"local"`
	Remote   ServerAddress           `yaml:"remote"`
	Mode     ServerMode              `yaml:"mode"`
	LogLevel log.Level               `yaml:"logLevel"`
	DB       repository.DBPoolConfig `yaml:"db"`
}

type ServerAddress struct {
	Port uint   `yaml:"port"`
	Host string `yaml:"host"`
}

func (s *ServerAddress) Address() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}
