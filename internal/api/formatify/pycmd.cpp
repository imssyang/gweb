#include <iostream>
#include <string>
#include <pybind11/embed.h>
#include "pycmd.h"

namespace py = pybind11;

size_t PycmdDumps(char* cmd, size_t size, size_t indent)
{
	py::scoped_interpreter guard{};
	py::module_ pycmd = py::module_::import("pycmd_wrap");
	py::object result = pycmd.attr("dumps")(cmd, indent);
	std::string formated = result.cast<std::string>();
	if (formated.length() >= size) {
		py::print("error result tool long! (", formated.length(), " >= ", size, ")");
		return -1;
	}

	std::strcpy(cmd, formated.data());
	return 0;
}
