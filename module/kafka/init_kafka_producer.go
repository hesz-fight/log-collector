package kafka

import (
	"log-collector/global/errcode"

	"github.com/shopify/sarama"
)

var Producer *KafkaSyncProducer

type KafkaSyncProducer struct {
	Producer sarama.SyncProducer
	Topic    string
	Config   *sarama.Config
}

func (p *KafkaSyncProducer) SendMessag(value string) (int32, int64, error) {
	msg := &sarama.ProducerMessage{
		Topic: p.Topic,
		Value: sarama.StringEncoder(value),
	}

	return p.Producer.SendMessage(msg)
}

func (p *KafkaSyncProducer) Close() error {
	return p.Producer.Close()
}

func InitProducer(addrs []string, topic string) error {
	var err error
	Producer, err = NewSyncProducer(addrs, topic, nil)
	if err != nil {
		return errcode.InitKafkaError.WithDetail(err.Error()).ToError()
	}

	return nil
}

func NewSyncProducer(addrs []string, topic string, config *sarama.Config) (*KafkaSyncProducer, error) {
	if config == nil {
		config = sarama.NewConfig()
		// learder and followers must return ack
		config.Producer.RequiredAcks = sarama.WaitForAll
		config.Producer.Partitioner = sarama.NewRandomPartitioner
		config.Producer.Return.Successes = true
	}
	producer, err := sarama.NewSyncProducer(addrs, config)
	if err != nil {
		return nil, err
	}

	return &KafkaSyncProducer{Producer: producer, Topic: topic, Config: config}, nil
}