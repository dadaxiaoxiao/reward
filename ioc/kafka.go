package ioc

import (
	"github.com/IBM/sarama"
	"github.com/dadaxiaoxiao/go-pkg/saramax"
	"github.com/dadaxiaoxiao/reward/internal/events"
	"github.com/spf13/viper"
)

// InitKafka 初始化 kafka 客户端
func InitKafka() sarama.Client {
	type Config struct {
		Addrs []string `yaml:"addrs"`
	}
	saramaCfg := sarama.NewConfig()
	saramaCfg.Producer.Return.Successes = true
	saramaCfg.Producer.Partitioner = sarama.NewConsistentCRCHashPartitioner
	var cfg Config
	err := viper.UnmarshalKey("kafka", &cfg)
	if err != nil {
		panic(err)
	}
	client, err := sarama.NewClient(cfg.Addrs, saramaCfg)
	if err != nil {
		panic(err)
	}
	return client
}

// NewConsumers 面临的问题依旧是所有的 Consumer 在这里注册一下
func NewConsumers(payment *events.PaymentEventConsumer) []saramax.Consumer {
	return []saramax.Consumer{payment}
}
