package medis

import (
	"sync"
)

var magazine *Magazine
var once sync.Once

type Magazine struct {
	ListKey      string     // 存储已生成数据的键名
	MaxKey       string     // 存储当前最大值的键名
	Channel      chan int64 // 信道
	Capacity     int64      // 信道容量
	Threshold    int64      // 信道容量阈值 低于此值时触发补充操作
	KvThreshold  int64      // 存储容量阈值 低于此值时触发补充操作
	KvSupplement int64      // 触发存储补充操作时具体的补充量
}

func initMagazine(empty bool) *Magazine {
	once.Do(func() {
		capacity := Capacity
		channel := make(chan int64, capacity)
		threshold := int64(float64(capacity) * Percent)
		// 初始化薄雾结构体时生成第一批值，同时将最大值存入 kv
		if empty == false {
			var i int64
			for i = 1; i < capacity+1; i++ {
				// 预生成 预存
				sequence := generate(i)
				channel <- sequence
			}
		}
		// 存储容量的阈值为信道阈值的指定倍数 存储补充量为信道容量指定倍数 倍数由配置决定
		magazine = &Magazine{
			ListKey: ListKey, MaxKey: MaxKey, Capacity: capacity,
			Threshold: threshold, Channel: channel, KvThreshold: threshold * Multiple,
			KvSupplement: capacity * Multiple,
		}
	})
	return magazine
}

func GetMagazine() *Magazine {
	return magazine
}
