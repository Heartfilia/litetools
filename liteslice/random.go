package liteslice

import (
	"math/rand"
	"sync"
	"time"
)

// 从数组串里面随机取一个值

var (
	randomChoiceMu  sync.Mutex
	randomChoiceRng = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func RandomChoice[T any](sliceArray []T) T {
	if len(sliceArray) == 0 {
		panic("RandomChoice requires a non-empty slice")
	}

	randomChoiceMu.Lock()
	num := randomChoiceRng.Intn(len(sliceArray))
	randomChoiceMu.Unlock()
	return sliceArray[num]
}
