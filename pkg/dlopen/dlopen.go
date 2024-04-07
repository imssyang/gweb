//go:build linux

package dlopen

/*
#include <stdlib.h>
#include <dlfcn.h>
#cgo LDFLAGS: -ldl
*/
import "C"
import (
	"errors"
	"fmt"
	"unsafe"
)

type Library struct {
	Name   string
	Handle unsafe.Pointer
}

func NewLibrary(names []string) (*Library, error) {
	for _, name := range names {
		l := &Library{}
		err := l.Open(name)
		if err == nil {
			return l, nil
		}
	}
	return nil, fmt.Errorf("unable to open the library: %v", names)
}

func (l *Library) Open(name string) error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	C.dlerror()
	handle := C.dlopen(cname, C.RTLD_LAZY)
	if handle == nil {
		err := C.dlerror()
		return fmt.Errorf("error opening %q: %v", name, errors.New(C.GoString(err)))
	}

	l.Name = name
	l.Handle = handle
	return nil
}

func (l *Library) Symbol(symbol string) (unsafe.Pointer, error) {
	sym := C.CString(symbol)
	defer C.free(unsafe.Pointer(sym))

	C.dlerror()
	p := C.dlsym(l.Handle, sym)
	err := C.dlerror()
	if err != nil {
		return nil, fmt.Errorf("error resolving %v: %v", symbol, errors.New(C.GoString(err)))
	}

	return p, nil
}

func (l *Library) Close() error {
	C.dlerror()
	C.dlclose(l.Handle)
	err := C.dlerror()
	if err != nil {
		return fmt.Errorf("error closing %v: %v", l.Name, errors.New(C.GoString(err)))
	}

	return nil
}
