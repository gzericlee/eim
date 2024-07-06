package cache

import "testing"

var c *Cache[string, interface{}]

func init() {
	var err error
	c, err = NewCache[string, interface{}]("test", 1000000)
	if err != nil {
		panic(err)
	}
}

func TestCache_Set(t *testing.T) {
	c.Set("key", "value")
	t.Log(c.Get("key"))
	t.Log(c.Get("key1"))
}

func BenchmarkCache_Set(b *testing.B) {
	for i := 0; i < b.N; i++ {
		c.Set("key", "value")
	}
}
