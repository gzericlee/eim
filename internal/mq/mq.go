package mq

import "github.com/nats-io/nats.go"

type IBase interface {
	printDetails()
}

type IHandler interface {
	Process(msg *nats.Msg) error
}

type IProducer interface {
	IBase
	Publish(subj string, body []byte) error
}

type IConsumer interface {
	IBase
	Subscribe(subj string, queue string, handler IHandler) error
}

func NewProducer(endpoints []string) (IProducer, error) {
	return newNatsProducer(endpoints)
}

func NewConsumer(endpoints []string) (IConsumer, error) {
	return newNatsConsumer(endpoints)
}
