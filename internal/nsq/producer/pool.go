package producer

import (
	"sync/atomic"

	"github.com/nsqio/go-nsq"
)

type Pool struct {
	Producers   []*nsq.Producer
	MaxAttempts int
	next        uint32
}

func NewPool(addrs []string, cfg *nsq.Config) (*Pool, error) {
	p := &Pool{
		Producers: make([]*nsq.Producer, len(addrs)),
	}
	for i, a := range addrs {
		np, err := nsq.NewProducer(a, cfg)
		if err != nil {
			return nil, err
		}
		p.Producers[i] = np
	}
	return p, nil
}

func (p *Pool) Publish(topic string, body []byte) error {
	n := atomic.AddUint32(&p.next, 1)
	l := len(p.Producers)

	for attempt := 0; (attempt <= p.MaxAttempts || p.MaxAttempts <= 0) && attempt < l; attempt++ {
		producer := p.Producers[(int(n)+attempt-1)%l]
		err := producer.Publish(topic, body)
		if err != nil {
			return err
		}
		break
	}

	return nil
}

func (p *Pool) PublishAsync(topic string, body []byte, doneChan chan *nsq.ProducerTransaction, args ...interface{}) error {
	n := int(atomic.AddUint32(&p.next, 1))

	maxAttempts := len(p.Producers)
	if maxAttempts > p.MaxAttempts {
		maxAttempts = p.MaxAttempts
	}

	go func(n int) {
		ch := make(chan *nsq.ProducerTransaction, 1)
		defer close(ch)

		for attempt := 0; attempt < maxAttempts; attempt++ {

			isLastAttempt := attempt+1 == maxAttempts

			index := (n + attempt - 1) % len(p.Producers)
			producer := p.Producers[index]

			if err := producer.PublishAsync(topic, body, ch, args...); err != nil {
				if isLastAttempt {
					doneChan <- &nsq.ProducerTransaction{Error: err, Args: args}
					break
				}
				continue
			}

			transaction := <-ch
			if transaction.Error != nil && !isLastAttempt {
				continue
			}

			doneChan <- transaction
			break
		}
	}(n)

	return nil
}
