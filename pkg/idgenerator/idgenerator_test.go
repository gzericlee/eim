package idgenerator

import (
	"testing"
)

func init() {
	Init([]string{"127.0.0.1:7001", "127.0.0.1:7002", "127.0.0.1:7003", "127.0.0.1:7004", "127.0.0.1:7005"}, "pass@word1")
}

func BenchmarkNextId(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NextId()
	}
}
