package es

import (
	"context"
	"fmt"
	"github.com/olivere/elastic"
	"log"
	"os"
	"time"
)

// 数据结构体
type EsDataStruct struct {
	Msg string	`json:"msg"`
}

// 通道数据结构体
type SendDataStruct struct {
	index 		string
	typeName 	string
	msg 		EsDataStruct
}

// 存储kafka消费数据的通道
var DataChan chan SendDataStruct

// 初始化
func Init()  {
	address := "http://" + EsConf.Ip + ":" + EsConf.Port

	client, err := elastic.NewClient(
		elastic.SetURL(address),							//指定要连接的URL（默认值是http://127.0.0.1:9200）
		elastic.SetSniff(EsConf.Sniff),						//允许您指定弹性是否应该定期检查集群（默认为true）
		elastic.SetHealthcheckInterval(time.Duration(EsConf.HealthcheckInterval)*time.Second),			//指定间隔之间的两个健康检查（默认是60秒）
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),		//将日志记录器设置为用于错误消息（默认为NIL）
		elastic.SetGzip(true),							//启用或禁用请求端的压缩。默认情况下禁用
	)
	if err != nil {
		fmt.Println("init es err :", err)
		return
	}
	ctx := context.Background()
	DataChan = make(chan SendDataStruct, 10000)

	// 后台插入数据
	go InsertData(client, ctx)
}

// es插数据
func InsertData(client *elastic.Client, ctx context.Context)  {
	for {
		select {
		case sendData := <- DataChan:
			index := sendData.index
			typeName := sendData.typeName
			data:= sendData.msg

			fmt.Printf("%v\n",data)
			put, err := client.Index().Index(index).Type(typeName).BodyJson(data).Do(ctx)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Printf("Indexed  [id: %s] to [index: %s], [type: %s]\n", put.Id, put.Index, put.Type)
		default:
			time.Sleep(time.Microsecond)
		}
	}
}

// 将kafka消费的数据存储到通道
func SendToEsChan(index, topic string, esData EsDataStruct)  {
	msg := SendDataStruct{
		index:    index,
		typeName: topic,
		msg:	  esData,
	}
	DataChan <- msg
}
