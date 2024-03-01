package literand

import (
	"math/rand"
	"time"
)

// 从数组串里面随机取一个值

func RandomChoice[T any](sliceArray []T) T {
	seed := rand.New(rand.NewSource(time.Now().Unix()))
	num := seed.Intn(len(sliceArray))
	return sliceArray[num]
}
