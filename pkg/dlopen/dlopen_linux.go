//go:build linux

package dlopen

/*
#include <stdlib.h>
#include <dlfcn.h>
#cgo LDFLAGS: -ldl
*/
import "C"
