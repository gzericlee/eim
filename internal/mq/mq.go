package mq

import "github.com/nsqio/go-nsq"

type base interface {
	printDetails()
}

type Producer interface {
	base
	Publish(topic string, body []byte) error
}

type Consumer interface {
	base
	Subscribe(topic string, channel string, handler nsq.Handler) error
}

func NewProducer(endpoints []string) (Producer, error) {
	return newNsqProducer(endpoints)
}

func NewConsumer(endpoints []string) (Consumer, error) {
	return newNsqConsumer(endpoints)
}
