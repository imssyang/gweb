#pragma once
#include <iostream>
#include <functional>
#include <queue>
#include <string>
#include <pybind11/embed.h>

namespace py = pybind11;

template <typename T>
inline void hash_combine(std::size_t &seed, const T &val) {
    seed ^= std::hash<T>()(val) + 0x9e3779b9 + (seed << 6) + (seed >> 2);
}

template <typename... Types>
inline std::size_t hash_val(const Types &... args) {
    std::size_t seed = 0;
    (hash_combine(seed, args), ...);
    return seed;
}

class PyCache {
public:
    std::string check(std::size_t hash_req) {
        if (results_.find(hash_req) == results_.end()) {
            if (requests_.size() > 1) {
                std::size_t hash_key = requests_.front();
                requests_.pop();
                results_.erase(hash_key);
            }
            return "";
        }
        return results_[hash_req];
    }

    void cache(std::size_t hash_req, std::string result) {
        requests_.push(hash_req);
        results_[hash_req] = result;
    }

private:
    std::unordered_map<std::size_t, std::string> results_;
    std::queue<std::size_t> requests_;
};

class PyFmt : public PyCache {
public:
    std::string dumps(const std::string& mode, const std::string& data, size_t indent) {
        std::size_t hash_req = hash_val(mode, data, indent);
        auto checked_data = check(hash_req);
        if (!checked_data.empty()) {
            return checked_data;
        }

		py::scoped_interpreter guard{};
		py::module_ pyfmt = py::module_::import("formatify");
		py::object result = pyfmt.attr("dumps")(mode, data, indent);
		if (result.is_none()) {
            py::print("[CGO] python formatify.dumps(", mode, ") fail.");
			return "";
		}

		auto formatted_data = result.cast<std::string>();
        cache(hash_req, formatted_data);
        return formatted_data;
	}
};
