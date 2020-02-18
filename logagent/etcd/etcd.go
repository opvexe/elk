package etcd

import (
	"fmt"
	etcd "go.etcd.io/etcd/clientv3"
	"strings"
	"time"
)

// etcd 客户端对象
type EtcdClient struct {
	client *etcd.Client
	keys []string //存储日志收集的key
}

var etcdClient *EtcdClient

// 初始化etcd
func NewEtcd(address []string,key string) (err error) {
	client,err := etcd.New(etcd.Config{
		Endpoints:            address,
		DialKeepAliveTimeout:5*time.Second,
	})
	if err!=nil{
		return
	}
	etcdClient = &EtcdClient{
		client: client,
	}
	if strings.HasSuffix(key,"/") == false{
		key = fmt.Sprintf("%s/",key)
	}
	//通过本地ip和配置文件中的前缀获取etcd数据
	for _, ip := range localIPSouce{
		etcdkey := fmt.Sprintf("%s%s",key,ip)
	}
}