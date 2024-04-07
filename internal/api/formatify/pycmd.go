package formatify

// #cgo CFLAGS: -I/opt/python3/include/python3.8
// #cgo CXXFLAGS: -I/opt/python3/include/python3.8 -I/root/.blog/gweb/third_party -std=c++11
// #cgo LDFLAGS: -L/opt/python3/lib -lpython3.8
// #include <Python.h>
// #include "cmd.h"
import "C"
import "unsafe"

func Command() (int, error) {
	name := C.CString("Gopher")
	defer C.free(unsafe.Pointer(name))

	return int(C.calc()), nil
}

//include "cmd.h"
// CXXFLAGS: -I/root/.blog/gweb/third_party -std=c++11
// LDFLAGS: -L/opt/python3/lib -lpython3.8 -lstdc++
