package kafka

import (
	"log-collector/global/errcode"
	"sync"

	"github.com/shopify/sarama"
)

var (
	Producer *KafkaSyncProducer
	once     sync.Once
)

type KafkaSyncProducer struct {
	Producer sarama.SyncProducer
	Config   *sarama.Config
}

func (p *KafkaSyncProducer) SendMessag(topic string, value string) (int32, int64, error) {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(value),
	}

	return p.Producer.SendMessage(msg)
}

func (p *KafkaSyncProducer) Close() error {
	return p.Producer.Close()
}

func InitProducer(addrs []string) error {
	var err error
	once.Do(func() {
		producer, e := NewSyncProducer(addrs, nil)
		if e != nil {
			err = errcode.InitKafkaError.WithDetail(e.Error()).ToError()
			return
		}
		Producer = producer
	})

	return err
}

func NewSyncProducer(addrs []string, config *sarama.Config) (*KafkaSyncProducer, error) {
	if config == nil {
		config = sarama.NewConfig()
		// learder and followers must return ack
		config.Producer.RequiredAcks = sarama.WaitForAll
		config.Producer.Partitioner = sarama.NewRandomPartitioner
		config.Producer.Return.Successes = true
	}

	cli, err := sarama.NewClient(addrs, config)
	if err != nil {
		return nil, err
	}

	producer, err := sarama.NewSyncProducerFromClient(cli)
	if err != nil {
		return nil, err
	}

	return &KafkaSyncProducer{Producer: producer, Config: config}, nil
}
