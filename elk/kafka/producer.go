package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"time"
)

type LogData struct {
	topic 	string
	data 	string
}

var ProductChan chan *LogData	// 存放每条任务的chan
var Producer sarama.SyncProducer	// 全局生产者

// 从通道获取消息并生产
func SendMsg()  {
	//构造消息
	message := &sarama.ProducerMessage{}
	logData := new(LogData)
	for {
		select {
		case logData = <- ProductChan:
			message.Topic = logData.topic
			message.Value = sarama.StringEncoder(logData.data)
			_, _, err := Producer.SendMessage(message)
			if err != nil {
				fmt.Println(err)

			}
		default:
			time.Sleep(time.Microsecond)
		}
	}
}

//初始化生产者
func InitProducer() {
	var err error
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	//连接kafka
	address := KafkaConfig.Kafka.Ip + ":" + KafkaConfig.Kafka.Port
	Producer, err = sarama.NewSyncProducer([]string{address}, config)
	if err != nil {
		fmt.Println("init producer err : ", err)
		return
	}
	// 初始化ConsumeChan
	ProductChan = make(chan *LogData, KafkaConfig.ChanSize)
	// 在后台发送数据
	go SendMsg()
}

func SendToProducerChan(topic, data string)  {
	logData := new(LogData)
	logData.topic = topic
	logData.data = data
	ProductChan <- logData
}
