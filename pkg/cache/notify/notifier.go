package notify

import (
	"fmt"
	"sync"
	"time"

	"eim/pkg/cache"
)

type Notifier struct {
	provider IProvider
	pools    sync.Map
}

func NewNotifier(provider IProvider) *Notifier {
	return &Notifier{
		provider: provider,
	}
}

func (its *Notifier) Bind(pool string, callback func(payload []string), caches ...*cache.Cache) {
	if callback == nil {
		callback = func(payload []string) {}
	}
	c, _ := its.pools.LoadOrStore(pool, &[]*cache.Cache{})
	*(c.(*[]*cache.Cache)) = append(*(c.(*[]*cache.Cache)), caches...)

	go func() {
		for {
			for !its.provider.OK() {
				time.Sleep(10 * time.Millisecond)
			}
			_ = its.provider.Sub(notifyTopic, callback)
		}
	}()
}

func (its *Notifier) Change(pool string, payload []string) error {
	err := its.provider.Pub(notifyTopic, payload)
	if err != nil {
		return fmt.Errorf("notify change(%s) -> %w", payload[0], err)
	}
	return nil
}
