package utils

type intNumber interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}
type uintNumber interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}
type floatNumber interface {
	~float32 | ~float64
}

type number interface {
	intNumber | uintNumber | floatNumber
}

func Max[T number](s ...T) T {
	m := s[0]
	for _, v := range s {
		if m < v {
			m = v
		}
	}
	return m
}

func Min[T number](s ...T) T {
	m := s[0]
	for _, v := range s {
		if m > v {
			m = v
		}
	}
	return m
}
