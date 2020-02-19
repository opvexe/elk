package kafka

import (
	"fmt"
	"gopkg.in/ini.v1"
)

type Config struct {
	Kafka	`ini:"kafka"`
}

type Kafka struct {
	Ip 			string 		`ini:"ip"`
	Port 		string		`ini:"port"`
	ChanSize	int			`ini:"chanSize"`
}

var KafkaConfig = new(Config)

func LoadConfig() {
	err := ini.MapTo(KafkaConfig, "./config/conf.ini")
	if err != nil {
		fmt.Println("ini config err:", err)
		return
	}
}
