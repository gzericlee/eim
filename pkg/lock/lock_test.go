package lock

import (
	"fmt"
	"testing"
)

func BenchmarkNewKeyLock(b *testing.B) {
	keyLock := NewKeyLock()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key-%d", i)
		keyLock.TryLock(key)
		keyLock.Unlock(key)
	}
}
