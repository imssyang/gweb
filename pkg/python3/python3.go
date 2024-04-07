package python3

/*
#include <stdio.h>
#include <stdlib.h>

typedef int (*add_mod_t)(int a, int b, int mod);
int add_mod_wrap(void *fn, int a, int b, int mod) {
	add_mod_t add_mod = (add_mod_t)fn;
	return add_mod(a, b, mod);
}
*/
import "C"

import (
	"fmt"

	"github.com/imssyang/gweb/pkg/dlopen"
)

func AddMod(names []string, a, b, mod int) (int, error) {
	lib, err := dlopen.NewLibrary(names)
	if err != nil {
		return -1, fmt.Errorf(`couldn't get a handle to the library: %v`, err)
	}
	defer lib.Close()

	add_mod, err := lib.Symbol("add_mod")
	if err != nil {
		return -1, fmt.Errorf(`couldn't get add_mod: %v`, err)
	}

	res := C.add_mod_wrap(add_mod, C.int(a), C.int(b), C.int(mod))

	return int(res), nil
}
