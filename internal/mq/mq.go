package mq

type base interface {
	printDetails()
}

type Handler interface {
	HandleMessage(data []byte) error
}

type Producer interface {
	base
	Publish(subj string, body []byte) error
}

type Consumer interface {
	base
	Subscribe(subj string, queue string, handler Handler) error
}

func NewProducer(endpoints []string) (Producer, error) {
	return newNatsProducer(endpoints)
}

func NewConsumer(endpoints []string) (Consumer, error) {
	return newNatsConsumer(endpoints)
}