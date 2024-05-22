package snowflake

import (
	"testing"
	"time"
)

var generator *Generator

// NodeCount = 1 单调递增，> 1 趋势递增
func init() {
	var err error
	generator, err = NewGenerator(GeneratorConfig{
		RedisEndpoints: []string{"127.0.0.1:7001", "127.0.0.1:7002", "127.0.0.1:7003", "127.0.0.1:7004", "127.0.0.1:7005"},
		RedisPassword:  "pass@word1",
		MaxWorkerId:    1023,
		MinWorkerId:    1,
		NodeCount:      5,
	})
	if err != nil {
		panic(err)
	}
}

func BenchmarkNextId(b *testing.B) {
	for i := 0; i < b.N; i++ {
		generator.NextId()
	}
}

func TestCheckId(t *testing.T) {
	idMap := make(map[int64]bool)
	lastId := int64(0)
	for i := 0; i < 10000; i++ {
		id := generator.NextId()
		if _, exists := idMap[id]; exists {
			t.Fatalf("Duplicate id found: %d", id)
		}
		if lastId > 0 && id < lastId {
			t.Fatalf("Invalid id found: %d", id)
		}
		idMap[id] = true
		t.Log(id)
		lastId = id
		time.Sleep(time.Millisecond * 20)
	}
}
