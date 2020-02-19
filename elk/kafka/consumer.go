package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"shumin-project/elk-server/elk/es"
	"shumin-project/elk-server/elk/etcd"
	"sync"
)

var (
	Consumer  sarama.Consumer 		// 全局消费者
	Topics	  map[string]string		 // 存储已发现的topic
)

// 开始消费数据
func Consume(topic string) {
	defer Consumer.Close()
	for {
		partitionList, err := Consumer.Partitions(topic) //根据topic找到所有的分区
		if err != nil {
			fmt.Println(err)
			return
		}
		//遍历所有的分区
		wg := sync.WaitGroup{}
		for _, partition := range partitionList {
			pc, err := Consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
			if err != nil {
				fmt.Println(err)
			}
			defer pc.Close()
			//异步的消费数据
			wg.Add(1)
			go func(partitionConsumer sarama.PartitionConsumer) {
				defer wg.Done()
				for msg := range partitionConsumer.Messages() {
					// 将数据发往es的通道
					logData := es.EsDataStruct{Msg:string(msg.Value)}
					es.SendToEsChan(msg.Topic, etcd.EtcdConfig.Key, logData)
				}
			}(pc)
		}
		wg.Wait()
	}
}

// 初始化消费者
func InitConsumer() {
	var err error
	address := KafkaConfig.Ip + ":" + KafkaConfig.Port
	Consumer, err = sarama.NewConsumer([]string{address}, nil)
	if err != nil {
		fmt.Println("new consumer err: ", err)
		return
	}
	Topics = make(map[string]string, KafkaConfig.ChanSize)
}

// 接收最新的消费topic
func SendToConsumeChan(topic string) {
	_, ok := Topics[topic]
	if !ok {
		go Consume(topic)
		Topics[topic] = topic
	}
}
