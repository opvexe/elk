package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	etcd "go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"shumin-project/elk-server/logagent/config"
	"shumin-project/elk-server/logagent/taillog"
	"strings"
	"time"
)

// etcd 客户端对象
type EtcdClient struct {
	client *etcd.Client
	keys   []string //存储日志收集的key
}

var etcdClient *EtcdClient

// 初始化etcd
func NewEtcd(address []string, key string) (err error) {
	client, err := etcd.New(etcd.Config{
		Endpoints:            address,
		DialKeepAliveTimeout: 5 * time.Second,
	})
	if err != nil {
		return
	}
	etcdClient = &EtcdClient{
		client: client,
	}
	if strings.HasSuffix(key, "/") == false {
		key = fmt.Sprintf("%s/", key)
	}
	//通过本地ip和配置文件中的前缀获取etcd数据
	for _, ip := range localIPSouce {
		etcdkey := fmt.Sprintf("%s%s", key, ip)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		resp, err := etcdClient.client.Get(ctx, etcdkey)
		cancel()

		if err != nil {
			fmt.Printf("get key:%s from etcd failed,err%s\n", etcdkey, err)
			continue
		}
		etcdClient.keys = append(etcdClient.keys, etcdkey)
		for _, v := range resp.Kvs {
			if string(v.Key) == etcdkey {
				err = json.Unmarshal(v.Value, &config.AgentConf.CollectionKey)
				if err != nil {
					fmt.Printf("json unmarsha key:%s failed,err:%s", v.Key, err)
					continue
				}

			}
		}
	}
	initEtcdWatch()
	return
}

// 初始化etcd监控
func initEtcdWatch() {
	for _, key := range etcdClient.keys {
		go etcdWatch(key)
	}
}

// etcd 监控处理
func etcdWatch(key string) {
	fmt.Println("start watch key:%s", key)
	for true {
		rech := etcdClient.client.Watch(context.Background(), key)
		var config []taillog.TailCollection
		var status bool = true
		for wresp := range rech {
			for _, ev := range wresp.Events {
				// key 删除
				if ev.Type == mvccpb.DELETE {
					fmt.Println("key [%s] is deleted", key)
					continue
				}
				// key 更新
				if ev.Type == mvccpb.PUT && string(ev.Kv.Key) == key {
					err := json.Unmarshal(ev.Kv.Value, &config)
					if err != nil {
						fmt.Println("key [%s], unmareshal [%s],err:", key, err)
						status = false
						continue
					}
				}
				fmt.Printf("get etcd config, %s %q : %q", ev.Type, ev.Kv.Key, ev.Kv.Value)
			}
			if status {
				break
			}
		}
		//更新tail任务
		err := taillog.UpdataTailTask(config)
		if err != nil {
			fmt.Printf("Update tailf task failed, connect: %s, err: %s", config, err)
		}
	}
}
