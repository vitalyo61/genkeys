package config

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	Server   *Server
	DataBase *DataBase `toml:"data_base"`
}

type Server struct {
	Address string
	Timeout int // in seconds
}

type DataBase struct {
	URL string `toml:"url"`
}

func Make() *Config {
	return new(Config)
}

func (c *Config) Open(name string) error {
	_, err := toml.DecodeFile(name, c)
	return err
}
