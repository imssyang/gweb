extern "C" {
	#include "cmd.h"
}

#include <iostream>
#include <pybind11/embed.h>

namespace py = pybind11;

void hello1() {
    py::exec(R"(
        kwargs = dict(name="World", number=42)
        message = "Hello1, {name}! The answer is {number}".format(**kwargs)
        print(message)
    )");
}

void hello2() {
	using namespace py::literals;
    auto kwargs = py::dict("name"_a="World", "number"_a=42);
    auto message = "Hello2, {name}! The answer is {number}"_s.format(**kwargs);
    py::print(message);
}

void hello3() {
	using namespace py::literals;
    auto locals = py::dict("name"_a="World", "number"_a=42);
    py::exec(R"(
        message = "Hello3, {name}! The answer is {number}".format(**locals())
    )", py::globals(), locals);

    auto message = locals["message"].cast<std::string>();
    std::cout << message << std::endl;
}

void sys_calc() {
	py::module_ sys = py::module_::import("sys");
	py::print(sys.attr("path"));

	py::module_ calc = py::module_::import("calc");
	py::object result = calc.attr("add")(1, 2);
	int n = result.cast<int>();
	assert(n == 3);
	std::cout << "sys_calc: " << n << std::endl;
}

void cmd_getInteger() {
	py::module_ sys = py::module_::import("sys");
	py::print(sys.attr("path"));

	py::module_ cmd = py::module_::import("cmd");
	py::object result = cmd.attr("getInteger")();
	int n = result.cast<int>();
	std::cout << n << std::endl;
}

void cmd_formatCommand() {
	py::module_ sys = py::module_::import("sys");
	py::print(sys.attr("path"));

	py::module_ cmd = py::module_::import("cmd");
	py::object result = cmd.attr("formatCommand")();
	std::string r = result.cast<std::string>();
	std::cout << r << std::endl;
}

int calc() {
	std::cout << "enter: " << 123 << std::endl;
    py::scoped_interpreter guard{};

	hello1();
	hello2();
	hello3();
	sys_calc();
	cmd_getInteger();
	cmd_formatCommand();
	std::cout << "exit: " << 123 << std::endl;

	return 12345;
}
