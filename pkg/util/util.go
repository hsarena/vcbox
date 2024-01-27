package util

var isMock bool

func SetMock(b bool) {
	isMock = b
}

func IsMock() bool {
	return isMock
}

func ToF64(i64 []int64) []float64 {
	f64 := make([]float64, len(i64))
	var ii int64
	var i int
	for i, ii = range i64 {
		f64[i] = float64(ii / 1000)
	}
	return f64
}
