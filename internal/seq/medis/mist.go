package medis

import (
	"crypto/rand"
	"math/big"
)

// 随机因子二进制位数
const saltBit = uint(8)

// 随机因子移位数
const saltShift = uint(8)

// 高位数移位数
const increaseShift = saltBit + saltShift

func generate(increase int64) int64 {
	randA, _ := rand.Int(rand.Reader, big.NewInt(255))
	saltA := randA.Int64()
	randB, _ := rand.Int(rand.Reader, big.NewInt(255))
	saltB := randB.Int64()
	mist := (increase << increaseShift) | (saltA << saltShift) | saltB // 通过位运算实现自动占位
	return mist
}
