package liteslice

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
