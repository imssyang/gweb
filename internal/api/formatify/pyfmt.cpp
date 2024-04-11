#include <iostream>
#include <string>
#include <pybind11/embed.h>
#include "pyfmt.h"

namespace py = pybind11;

size_t PyfmtDumps(const char* mode, char* data, size_t size, size_t indent)
{
	py::scoped_interpreter guard{};
	py::module_ pyfmt = py::module_::import("formatify");
	py::object result = pyfmt.attr("dumps")(mode, data, indent);
	if (result.is_none()) {
        return -1;
    }
	std::string formatted = result.cast<std::string>();
	if (formatted.length() >= size) {
		py::print("[CGO] error result tool long. (", formatted.length(), " >= ", size, ")");
		return -2;
	}

	std::strcpy(data, formatted.data());
	return 0;
}
