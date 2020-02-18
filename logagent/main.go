package main

import (
	"flag"
	"fmt"
	"shumin-project/elk-server/logagent/config"
	"shumin-project/elk-server/logagent/etcd"
	"shumin-project/elk-server/logagent/kafka"
	"shumin-project/elk-server/logagent/taillog"
)

var (
	logPath string
	confType string = "ini"
)

func init() {
	flag.StringVar(&logPath,"p","","请输入日志日志")
	flag.Parse()
	if len(logPath) == 0 {
		flag.PrintDefaults()
	}
}

func main() {
	//加载配置文件
	err:=config.LoadConf(confType,logPath)
	if err!=nil{
		fmt.Println("load conf err:",err)
		return
	}
	conf := config.AgentConf

	// 初始化etcd
	 err =etcd.NewEtcd(conf.EtcdHost,conf.CollectionKey)
	 if err!=nil{
	 	fmt.Println("init etcd err:",err)
		 return
	 }

	 // 初始化tail
	 err =taillog.NewTail(conf.Collects,conf.Chansize,conf.IP)
	 if err!=nil{
	 	fmt.Println("init tail err:",err)
		 return
	 }

	 // 初始化kafka
	 err =kafka.NewKafka(conf.KafkaHost)
	 if err!=nil{
	 	fmt.Println("init kafka err:",err)
		 return
	 }

}
