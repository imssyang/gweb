#ifndef _PYFMT_H
#define _PYFMT_H

#include <Python.h>

#ifdef __cplusplus
extern "C" {
#endif

size_t PyfmtDumps(const char* mode, char* data, size_t size, size_t indent);

#ifdef __cplusplus
}
#endif

#endif // _PYFMT_H