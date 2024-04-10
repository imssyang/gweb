package formatify

/*
#cgo CFLAGS: -I/opt/python/pyenv/versions/3.9.18/include/python3.9
#cgo CXXFLAGS: -I/opt/python/pyenv/versions/3.9.18/include/python3.9 -I/opt/app/gweb/third_party -std=c++11
#cgo LDFLAGS: -L/opt/python/pyenv/versions/3.9.18/lib -lpython3.9
#include <Python.h>
#include "pycmd.h"
*/
import "C"
import (
	"fmt"
	"unsafe"
)

func PycmdDumps(cmd string, indent int) (string, error) {
	size := C.size_t(len(cmd)*2 + 8)
	formated := (*C.char)(C.malloc(size))
	defer C.free(unsafe.Pointer(formated))

	C.strcpy(formated, C.CString(cmd))
	errno := int(C.PycmdDumps(formated, size, C.size_t(indent)))
	if errno != 0 {
		return cmd, fmt.Errorf("C.PycmdDumps error code: %v", errno)
	}

	return C.GoString(formated), nil
}
