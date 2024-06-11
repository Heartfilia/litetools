package literand

import (
	"math/rand"
	"time"
)

// 从数组串里面随机取一个值

func RandomChoice[T any](sliceArray []T) T {
	seed := rand.New(rand.NewSource(time.Now().UnixMicro()))
	num := seed.Intn(len(sliceArray))
	return sliceArray[num]
}

// 移除数组里面的指定内容

func SliceRemove[T comparable](slice []T, key T) []T {
	for i := 0; i < len(slice); i++ {
		if slice[i] == key {
			slice = append(slice[:i], slice[i+1:]...)
			i--
		}
	}
	return slice
}
