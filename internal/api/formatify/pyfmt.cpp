#include <iostream>
#include <string>
#include <pybind11/embed.h>
#include "pyfmt.h"

namespace py = pybind11;

size_t PyfmtDumps(const char* mode, char* data, size_t size, size_t indent)
{
	py::scoped_interpreter guard{};
	py::module_ pyfmt = py::module_::import("pyfmt");
	py::object result = pyfmt.attr("dumps")(mode, data, indent);
	if (result.is_none()) {
		py::print("[CGO] error result is none for", data);
        return -1;
    }
	std::string formated = result.cast<std::string>();
	if (formated.length() >= size) {
		py::print("[CGO] error result tool long. (", formated.length(), " >= ", size, ")");
		return -2;
	}

	std::strcpy(data, formated.data());
	return 0;
}
