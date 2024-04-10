package formatify

/*
#cgo CFLAGS: -I/opt/python/pyenv/versions/3.9.18/include/python3.9
#cgo CXXFLAGS: -I/opt/python/pyenv/versions/3.9.18/include/python3.9 -I/opt/app/gweb/third_party -std=c++11
#cgo LDFLAGS: -L/opt/python/pyenv/versions/3.9.18/lib -lpython3.9
#include <Python.h>
#include "pyfmt.h"
*/
import "C"
import (
	"fmt"
	"unsafe"
)

func PyfmtDumps(mode, data string, indent int) (string, error) {
	mode_ := C.CString(mode)
	defer C.free(unsafe.Pointer(mode_))

	size_ := C.size_t(len(data)*2 + 8)
	formated_ := (*C.char)(C.malloc(size_))
	defer C.free(unsafe.Pointer(formated_))

	indent_ := C.size_t(indent)
	C.strcpy(formated_, C.CString(data))
	errno := int(C.PyfmtDumps(mode_, formated_, size_, indent_))
	if errno != 0 {
		return data, fmt.Errorf("C.PycmdDumps error code: %v", errno)
	}

	return C.GoString(formated_), nil
}
