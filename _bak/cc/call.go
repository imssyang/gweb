package cc

// #cgo CFLAGS: -I/opt/python3/include/python3.8
// #cgo LDFLAGS: -L/opt/python3/lib -lpython3.8
// #include "pow.c"
import "C"

func Run() {
	C.main2()
}
