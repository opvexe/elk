package config

import (
	"errors"
	"github.com/astaxie/beego/config"
	"strings"
)

type Conf struct {
	LogPath       string   //日志路径
	CollectionKey string   // 前缀
	KafkaHost     []string //kafka
	EtcdHost      []string //etcd
}

var AgentConf *Conf

// 加载配置信息
func loadConf(confType, confPath string) (err error) {
	conf, err := config.NewConfig(confType, confPath)
	if err != nil {
		return
	}
	AgentConf = &Conf{}
	//获取基础配置
	err = getAgentConf(conf)
	if err != nil {
		return
	}
	return
}

// 获取配置
func getAgentConf(conf config.Configer) (err error) {
	// 获取日志路径
	logPath := conf.String("log::log_path")
	if len(logPath) == 0 {
		logPath = "/Users/Facebook/logs/bussiess.log"
	}
	AgentConf.LogPath = logPath

	//获取etcd 地址
	etcdHosts := conf.String("etcd::kafka_host")
	if len(etcdHosts) == 0 {
		err = errors.New("etcd hosts error")
		return
	}
	AgentConf.EtcdHost = strings.Split(etcdHosts, ",")

	// 获取kaafka 地址
	kafkaHosts := conf.String("kafka::kafka_host")
	if len(kafkaHosts) == 0 {
		err = errors.New("kafkaa hosts err")
		return
	}
	AgentConf.KafkaHost = strings.Split(kafkaHosts, ",")

	// 获取前缀
	collectionPre := conf.String("collection::collection_key")
	if len(collectionPre) == 0 {
		err = errors.New("collection key is err")
		return
	}
	AgentConf.CollectionKey = collectionPre
	return
}
