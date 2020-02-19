package main

import (
	"shumin-project/elk-server/elk/es"
	"shumin-project/elk-server/elk/etcd"
	"shumin-project/elk-server/elk/kafka"
	"shumin-project/elk-server/elk/taillog"
)

func main() {
	// 初始化etcd
	etcd.Init()

	//加载es配置
	es.LoadEsConfig()

	//初始化es
	es.Init()

	//加载Kafka配置
	kafka.LoadConfig()

	//初始化消费者
	kafka.InitConsumer()

	//初始化生产者
	kafka.InitProducer()

	//etcd 获取日志收集的path和kafka的topic
	LogConf:=etcd.GetConf()

	//开启任务
	taillog.Init(LogConf)

	//etcd watch
	newConfChan :=taillog.PushConfToChan()
	go etcd.WatchConf(newConfChan)

	ch :=make(chan struct{},1)
	<- ch
}
