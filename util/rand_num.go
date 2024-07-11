package util

import (
	"math/rand"
	"time"
)

// GenerateRandomNumbers 生成 count 个范围在 [min, max) 之间的随机数
func GenerateRandomNumbers(count, min, max int) []int {
	rand.Seed(time.Now().UnixNano())
	numbers := make([]int, count)
	for i := 0; i < count; i++ {
		numbers[i] = rand.Intn(max-min) + min
	}
	return numbers
}
