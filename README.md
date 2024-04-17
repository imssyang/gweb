<a href="https://github.com/imssyang/gweb">
  <h1 align="center">
    <picture>
	  <img alt="GWeb" src="https://github.com/imssyang/gweb/blob/main/public/img/favicon.svg" width="100" />
    </picture>
    <p>GWeb</p>
  </h1>
</a>

[![license](https://img.shields.io/go/l/gweb.svg)](https://github.com/imssyang/gweb/blob/main/LICENSE)

GWeb is a backend-framework based on [gin](https://gin-gonic.com) in golang, and support for C or CPP language based on cgo. With the [pybind11](https://github.com/pybind/pybind11) of the open source community, also support Python after simple encapsulation. Therefore it can be used as a general framework to experiment with various development scenarios. Nothing limits developers to doing interesting things!

## Feature

- Support golang, c, cpp and python code.
- As backend of [formatui](https://github.com/imssyang/formatui).

## Dependencies

* [Gin Web Framework](https://gin-gonic.com): A fastest full-featured web framework for Go. 
* [pybind11](https://github.com/pybind/pybind11): Seamless operability between C++11 and Python.
* [Python](https://www.python.org): Python is a programming language that lets you work quickly
and integrate systems more effectively.

### Commands

```bash
# Run by docker, and open http://localhost:5005 in browser
docker run -it -p 5005:5005 --rm ghcr.io/imssyang/gweb:latest
```

## Todo

- Support more frontend
