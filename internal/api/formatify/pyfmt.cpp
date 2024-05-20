#include <iostream>
#include <string>
#include "internal.h"
#include "pyfmt.h"

static PyFmt pyfmt;

size_t PyfmtDesiredSize(const char* mode, const char* data, size_t indent) {
	std::string result = pyfmt.dumps(mode, data, indent);
	if (result.length() == 0) {
		return 0;
	}
	return result.length() + 1;
}

size_t PyfmtDumps(const char* mode, char* data, size_t size, size_t indent) {
	std::string result = pyfmt.dumps(mode, data, indent);
	if (result.length() >= size) {
		std::cout << "[CGO-error] size of result too long. ("
					<< result.length() << " >= " <<  size << ")"
					<< std::endl;
		return -1;
	}

	std::strcpy(data, result.data());
	return 0;
}
