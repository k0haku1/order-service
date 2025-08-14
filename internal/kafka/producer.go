package kafka

import (
	"github.com/IBM/sarama"
	"log"
)

type Producer struct {
	producer sarama.AsyncProducer
	topic    string
}

func NewProducer(brokers []string, topic string) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	prod := &Producer{
		producer: producer,
		topic:    topic,
	}

	go func() {
		for msg := range producer.Successes() {
			log.Println("SUCCESS", msg.Key, msg.Value)
		}
	}()
	go func() {
		for err := range producer.Errors() {
			log.Println("ERROR", err)
		}
	}()

	return prod, nil
}
func (p *Producer) Send(key string, value []byte) {
	p.producer.Input() <- &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(value),
	}
}

func (p *Producer) Close() error {
	err := p.producer.Close()
	if err != nil {
		return err
	}
	return nil
}
