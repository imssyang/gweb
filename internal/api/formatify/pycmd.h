#ifndef _PYCMD_H
#define _PYCMD_H

#include <Python.h>

#ifdef __cplusplus
extern "C" {
#endif

size_t PycmdDumps(char* cmd, size_t size, size_t indent);

#ifdef __cplusplus
}
#endif

#endif