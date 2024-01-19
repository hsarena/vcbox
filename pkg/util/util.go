package util

var isMock bool

func SetMock(b bool) {
	isMock = b
}

func IsMock() bool {
	return isMock
}
