package config

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	Server *Server
}

type Server struct {
	Address string
	Timeout int // in seconds
}

func Make() *Config {
	return new(Config)
}

func (c *Config) Open(name string) error {
	_, err := toml.DecodeFile(name, c)
	return err
}
