package config

import (
	"fmt"
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
	Port     uint       `yaml:"port"`
	Host     string     `yaml:"host"`
	Mode     ServerMode `yaml:"mode"`
	LogLevel log.Level  `yaml:"logLevel"`
}

func (s *Server) Address() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}
