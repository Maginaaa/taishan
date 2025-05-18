package dal

import (
	"data/config"
	"github.com/Shopify/sarama"
	"log"
	"strings"
)

var (
	kafkaProducer sarama.SyncProducer
)

// NewKafkaProducer 构建生产者
func NewKafkaProducer() {
	var err error

	kafkaConf := sarama.NewConfig()
	kafkaConf.Producer.RequiredAcks = sarama.WaitForAll                                                   // 发送完数据需要leader和follow都确认
	kafkaConf.Producer.Partitioner = sarama.NewManualPartitioner                                          // 设置选择分区的策略为Hash,当设置key时，所有的key的消息都在一个分区Partitioner里
	kafkaConf.Producer.Return.Successes = true                                                            // 成功交付的消息将在success channel返回
	kafkaProducer, err = sarama.NewSyncProducer(strings.Split(config.Conf.Kafka.Address, ","), kafkaConf) // 生产者客户端
	if err != nil {
		log.Fatal(err)
	}

	return
}
