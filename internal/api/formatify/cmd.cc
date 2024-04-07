extern "C" {
	#include "cmd.h"
}

#include <iostream>
#include <pybind11/embed.h>

namespace py = pybind11;

int calc() {
	std::cout << "enter: " << 123 << std::endl;
    py::scoped_interpreter guard{};
	py::module_ sys = py::module_::import("sys");
	py::print(sys.attr("path"));

	py::module_ cmd = py::module_::import("cmd");
	py::object result = cmd.attr("getInteger")();
	int n = result.cast<int>();
	std::cout << n;
	return 12345;
}