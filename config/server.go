package main

import (
	"fmt"

	"github.com/joeshaw/envdecode"
)

// Server config info
type serverConfig struct {
	Hostname string `env:"HOSTNAME,default=localhost"`
	Port     int    `env:"PORT,default=8008"`
}

func (cfg serverConfig) Address() string {
	return fmt.Sprintf("%v:%v", cfg.Hostname, cfg.Port)
}

var Server serverConfig

func init() {
	envdecode.Decode(&Server)
}
