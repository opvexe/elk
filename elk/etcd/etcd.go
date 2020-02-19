package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"time"
)

var Client *clientv3.Client // Etcd全局客户端
var LogConfig []*LogEntry   // 所有的日志收集项

type LogEntry struct {
	Path  string // tail日志的路径
	Topic string // kafka的topic
}

//初始化
func Init() {
	err := LoadConfig()
	if err != nil {
		fmt.Println("load etcd config err :", err)
		return
	}
	endpoints := EtcdConfig.Etcd.Ip + ":" + EtcdConfig.Etcd.Port
	Client, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{endpoints},
		DialTimeout: time.Second * time.Duration(EtcdConfig.Etcd.DialTimeout),
	})
	if err != nil {
		fmt.Println("init etcd err :", err)
		return
	}
	LogConfig = make([]*LogEntry, EtcdConfig.TaskNum)
}

// 获取所有的日志配置
func GetConf() []*LogEntry {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	result, err := Client.Get(ctx, EtcdConfig.Key)
	if err != nil {
		fmt.Println("get err:", err)
	}
	for _, kv := range result.Kvs {
		err := json.Unmarshal(kv.Value, &LogConfig)
		if err != nil {
			fmt.Println("unmarshal err : ", err)
			return nil
		}
	}
	return LogConfig
}

// 监控etcd
func WatchConf(newConfChan chan<- []*LogEntry) {
	watchChan := Client.Watch(context.Background(), EtcdConfig.Key)
	for {
		for event := range watchChan {
			for _, ev := range event.Events {
				var newConf []*LogEntry
				if ev.Type != mvccpb.DELETE {
					// 判断是否为删除事件
					err := json.Unmarshal(ev.Kv.Value, &newConf)
					if err != nil {
						fmt.Println("unmarshal err : ", err)
						return
					}
					newConfChan <- newConf
				}
			}
		}
	}
}
