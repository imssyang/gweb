#ifndef _PYFMT_H
#define _PYFMT_H

#include <Python.h>

#ifdef __cplusplus
extern "C" {
#endif

struct PyDumpsData {
    char* mode;
    char* data;
    size_t size;
    size_t indent;
    size_t has_escape;
};

size_t PyfmtDesiredSize(const char* mode, const char* data, size_t indent, size_t has_escape);
size_t PyfmtDumps(const char* mode, char* data, size_t size, size_t indent, size_t has_escape);

#ifdef __cplusplus
}
#endif

#endif // _PYFMT_H