#include <iostream>
#include <string>
#include "internal.h"
#include "pyfmt.h"

static PyFormat pyfmt;

size_t PyDumpsSize(const PyDumpsData* pydata)
{
	auto& pyd = *pydata;
	std::string result = pyfmt.dumps(pyd.mode, pyd.data, pyd.indent, pyd.has_escape);
	if (result.length() == 0) {
		return 0;
	}
	return result.length() + 1;
}

size_t PyDumps(PyDumpsData* pydata)
{
	auto& pyd = *pydata;
	std::string result = pyfmt.dumps(pyd.mode, pyd.data, pyd.indent, pyd.has_escape);
	if (result.length() >= pyd.size) {
		std::cout << "[CGO-error] size of result too long. ("
					<< result.length() << " >= " << pyd.size << ")"
					<< std::endl;
		return -1;
	}

	std::strcpy(pyd.data, result.data());
	return 0;
}
