package jsonPath

type Result struct {
	resultString string
	Int64        int64
	Int          int
	Float64      float64
	Float32      float32
	Bool         bool
	Object       interface{}
}

func (r *Result) String() string {
	return r.resultString
}

func (r *Result) SetString(str string) {
	r.resultString = str
}
