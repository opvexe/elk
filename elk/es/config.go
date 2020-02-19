package es

import (
	"fmt"
	"gopkg.in/ini.v1"
)

type EsConfig struct {
	Es		`ini:"es"`
}

type Es struct {
	Ip						string	`ini:"ip"`
	Port					string	`ini:"port"`
	Sniff					bool	`ini:"sniff"`						//允许您指定弹性是否应该定期检查集群（默认为true）
	HealthcheckInterval		int		`ini:"HealthcheckInterval"`			//指定间隔之间的两个健康检查（默认是60秒）
}

var EsConf = new(EsConfig)

// 加载配置
func LoadEsConfig()  {
	err := ini.MapTo(EsConf,"./config/conf.ini")
	if err != nil {
		fmt.Println("load es config err: ", err)
		return
	}
}
