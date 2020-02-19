package taillog

import (
	"fmt"
	"gopkg.in/ini.v1"
)

type Config struct {
	Tail `ini:"tail"`
}

type Tail struct {
	Path string `ini:"path"`
}

func LoadConfig() *Config {
	config := new(Config)
	err := ini.MapTo(config, "./config/conf.ini")
	if err != nil {
		fmt.Println("tail config err:", err)
		return nil
	}
	return config
}
