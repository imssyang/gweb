package format

// #include <pyfmt.h>
import "C"
import (
	"fmt"
	"math"
	"unsafe"
)

type PyDumpsData struct {
	Mode      string
	Data      string
	Indent    int
	HasEscape bool
	cdata     *C.PyDumpsData
}

func NewPyDumpsData(mode, data string, indent int, hasEscape bool) (*PyDumpsData, error) {
	d := &PyDumpsData{
		Mode:      mode,
		Data:      data,
		Indent:    indent,
		HasEscape: hasEscape,
	}
	cdata, err := d.newCData()
	if err != nil {
		return d, err
	}
	d.cdata = cdata
	return d, nil
}

func (d *PyDumpsData) newCData() (*C.PyDumpsData, error) {
	HasEscapeI := 0
	if d.HasEscape {
		HasEscapeI = 1
	}
	d_ := &C.PyDumpsData{
		mode:       C.CString(d.Mode),
		data:       C.CString(d.Data),
		size:       C.size_t(len(d.Data)),
		indent:     C.size_t(d.Indent),
		has_escape: C.size_t(HasEscapeI),
	}
	return d_, nil
}

func (d *PyDumpsData) resetCData(desiredSize int) error {
	newSize := int(math.Max(float64(len(d.Data)), float64(desiredSize)) + 1)
	newSize_ := C.size_t(newSize)
	newData_ := (*C.char)(C.malloc(newSize_))
	if newData_ == nil {
		return fmt.Errorf("malloc(%v) failed to allocate memory", newSize_)
	}

	C.memset(unsafe.Pointer(newData_), 0, newSize_)
	C.strncpy(newData_, d.cdata.data, newSize_-1)
	C.free(unsafe.Pointer(d.cdata.data))
	d.cdata.data = newData_
	d.cdata.size = newSize_
	return nil
}

func (d *PyDumpsData) freeCData() {
	if d.cdata != nil {
		C.free(unsafe.Pointer(d.cdata.data))
		C.free(unsafe.Pointer(d.cdata.mode))
		d.cdata = nil
	}
}

func PyDumps(mode, data string, indent int, hasEscape bool) (string, error) {
	d, err := NewPyDumpsData(mode, data, indent, hasEscape)
	if err != nil {
		return data, err
	}
	defer d.freeCData()

	desiredSize := int(C.PyDumpsSize(d.cdata))
	if desiredSize == 0 {
		return data, fmt.Errorf("C.PyDumpsSize fail.")
	}

	err = d.resetCData(desiredSize)
	if err != nil {
		return data, fmt.Errorf("resetCData(%v) fail: %v", desiredSize, err)
	}

	errno := int(C.PyDumps(d.cdata))
	if errno != 0 {
		return data, fmt.Errorf("C.PyDumps error code: %v", errno)
	}

	return C.GoString(d.cdata.data), nil
}
