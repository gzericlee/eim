package notify

import (
	"log"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"eim/pkg/cache"
)

const notifyTopic = "cache:change:notify"

type Notifier interface {
	OK() bool
	Pub(channel, payload string) error
	Sub(channel string, callback func(payload string)) error
}

var redisCli Notifier
var m sync.Map

func delAll(pool, key string) {
	if caches, _ := m.Load(pool); caches != nil {
		for _, c := range *(caches.(*[]*cache.Cache)) {
			c.Del(key)
		}
	}
}

func Init(r Notifier) {
	if redisCli != r {
		redisCli = r
		go func() {
			defer func() {
				if err := recover(); err != nil {
					log.Println(err)
					debug.PrintStack()
				}
			}()

			for {
				for r == nil || !r.OK() {
					time.Sleep(10 * time.Millisecond)
				}
				_ = r.Sub(notifyTopic, func(payload string) {
					vs := strings.Split(payload, ":")
					if len(vs) >= 2 {
						delAll(vs[0], vs[1])
					}
				})
			}
		}()
	}
}

func Bind(pool string, caches ...*cache.Cache) error {
	c, _ := m.LoadOrStore(pool, &[]*cache.Cache{})
	*(c.(*[]*cache.Cache)) = append(*(c.(*[]*cache.Cache)), caches...)
	return nil
}

func Del(pool, key string) error {
	r := redisCli
	if r != nil && r.Pub(notifyTopic, strings.Join([]string{pool, key}, ":")) == nil {
		return nil
	}
	delAll(pool, key)
	return nil
}
