package formatify

/*
#cgo CXXFLAGS: -std=c++11
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

	size_ := C.size_t(len(data)*2 + 1024)
	formatted_ := (*C.char)(C.malloc(size_))
	defer C.free(unsafe.Pointer(formatted_))

	indent_ := C.size_t(indent)
	C.strcpy(formatted_, C.CString(data))
	errno := int(C.PyfmtDumps(mode_, formatted_, size_, indent_))
	if errno != 0 {
		return data, fmt.Errorf("C.PycmdDumps error code: %v", errno)
	}

	return C.GoString(formatted_), nil
}
