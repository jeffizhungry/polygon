package main

import "github.com/joeshaw/envdecode"

// Server configuration
type serverConfig struct {
	Hostname string `env:"HOSTNAME,default=localhost"`
	Port     int    `env:"PORT,default=8008"`
}

// Server config singleton
var Server serverConfig

// Parse environmental variables
func init() {
	envdecode.Decode(&Server)
}
