package formatify

/*
#cgo CXXFLAGS: -std=c++17
#include <Python.h>
#include "pyfmt.h"
*/
import "C"
import (
	"fmt"
	"math"
	"unsafe"
)

func PyfmtDumps(mode, data string, indent int, has_escape bool) (string, error) {
	has_escape_i := 0
	if has_escape {
		has_escape_i = 1
	}
	mode_ := C.CString(mode)
	data_ := C.CString(data)
	indent_ := C.size_t(indent)
	has_escape_ := C.size_t(has_escape_i)
	defer C.free(unsafe.Pointer(mode_))
	defer C.free(unsafe.Pointer(data_))

	desired_size := int(C.PyfmtDesiredSize(mode_, data_, indent_, has_escape_))
	if desired_size == 0 {
		return data, fmt.Errorf("C.PyfmtDesiredSize fail.")
	}

	size_ := C.size_t(math.Max(float64(len(data)), float64(desired_size)))
	result_ := (*C.char)(C.malloc(size_))
	defer C.free(unsafe.Pointer(result_))

	C.strcpy(result_, data_)
	errno := int(C.PyfmtDumps(mode_, result_, size_, indent_, has_escape_))
	if errno != 0 {
		return data, fmt.Errorf("C.PycmdDumps error code: %v", errno)
	}

	return C.GoString(result_), nil
}
