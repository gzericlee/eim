package mq

import (
	"sync/atomic"

	"github.com/nsqio/go-nsq"
)

type producerPool struct {
	producers    []*nsq.Producer
	maxAttempts  int
	next         int32
	semaphore    chan struct{}
	successCount int64
	failedCount  int64
	retryCount   int64
}

func newProducerPool(address []string, cfg *nsq.Config, maxConcurrency int) (*producerPool, error) {
	p := &producerPool{
		producers:   make([]*nsq.Producer, len(address)),
		semaphore:   make(chan struct{}, maxConcurrency),
		maxAttempts: len(address),
	}
	for i, a := range address {
		np, err := nsq.NewProducer(a, cfg)
		if err != nil {
			return nil, err
		}
		p.producers[i] = np
	}
	return p, nil
}

func (its *producerPool) Close() {
	for _, producer := range its.producers {
		producer.Stop()
	}
}

func (its *producerPool) publishAsync(topic string, body []byte, doneChan chan *nsq.ProducerTransaction, args ...interface{}) error {
	n := int(atomic.AddInt32(&its.next, 1))

	maxAttempts := len(its.producers)
	if maxAttempts > its.maxAttempts {
		maxAttempts = its.maxAttempts
	}

	its.semaphore <- struct{}{}
	go func(n int) {
		defer func() { <-its.semaphore }()
		ch := make(chan *nsq.ProducerTransaction, 1)
		defer close(ch)
		for attempt := 0; attempt < maxAttempts; attempt++ {

			isLastAttempt := attempt+1 == maxAttempts

			index := (n + attempt) % len(its.producers)
			producer := its.producers[index]

			if err := producer.PublishAsync(topic, body, ch, args...); err != nil {
				if isLastAttempt {
					doneChan <- &nsq.ProducerTransaction{Error: err, Args: args}
					break
				}
				atomic.AddInt64(&its.retryCount, 1)
				continue
			}

			transaction := <-ch
			if transaction.Error == nil {
				atomic.AddInt64(&its.successCount, 1)
				doneChan <- transaction
				break
			} else if isLastAttempt {
				atomic.AddInt64(&its.failedCount, 1)
				doneChan <- transaction
				break
			}
		}
	}(n)

	return nil
}
