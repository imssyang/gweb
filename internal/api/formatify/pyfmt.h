#pragma once
#include <Python.h>

#ifdef __cplusplus
extern "C" {
#endif

typedef struct {
    char* mode;
    char* data;
    size_t size;
    size_t indent;
    size_t has_escape;
} PyDumpsData;

size_t PyDumpsSize(const PyDumpsData* pydata);
size_t PyDumps(PyDumpsData* pydata);

#ifdef __cplusplus
}
#endif
