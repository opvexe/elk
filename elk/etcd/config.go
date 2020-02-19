package etcd

import "gopkg.in/ini.v1"

type Config struct {
	Etcd `ini:"etcd"`
}

type Etcd struct {
	Ip          string `ini:"ip"`
	Port        string `ini:"port"`
	DialTimeout int    `ini:"DialTimeout"`
	Key         string `ini:"key"` // etcd日志配置的key
	ChanSize    int    `ini:"chanSize"`
	TaskNum     int    `ini:"taskNum"` // 日志任务数量
}

var EtcdConfig = new(Config)

func LoadConfig() error {
	err := ini.MapTo(EtcdConfig, "./config/conf.ini")
	return err
}
