GOENV_ROOT := /opt/go/goenv
PYENV_ROOT := /opt/python/pyenv
PATH := \
	$(GOENV_ROOT)/bin:$(GOENV_ROOT)/shims: \
	$(PYENV_ROOT)/bin:$(PYENV_ROOT)/shims: \
	$(PATH)

PROJECT_DIR=$(shell pwd)
LIBMEDIA_HOME=${PROJECT_DIR}/third_party/libmedia/release

PYTHON_HOME=$(shell pyenv prefix 2>/dev/null || echo /opt/python/pyenv/versions/3.12.2)
PYTHON_VER=$(shell ls ${PYTHON_HOME}/include 2>/dev/null || echo python3.12)
FFMPEG_HOME=/opt/ffmpeg
OS_TYPE := $(shell uname)

ifeq ($(OS_TYPE), Linux)
	CXXFLAGS = -std=c++2a
	LDFLAGS=-Wl,-rpath,'$ORIGIN'
	export LD_LIBRARY_PATH=${PYTHON_HOME}/lib:${FFMPEG_HOME}/lib:${LIBMEDIA_HOME}/lib
else ifeq ($(OS_TYPE), Darwin)
	CXXFLAGS = -std=c++20
	LDFLAGS=-Wl,-no_warn_duplicate_libraries
endif

export CGO_CFLAGS = -Wall -Wextra -O2 \
	-I${PYTHON_HOME}/include/${PYTHON_VER} \
	-I${FFMPEG_HOME}/include \
	-I${PROJECT_DIR}/third_party
export CGO_CXXFLAGS = ${CXXFLAGS} -O2 \
	-I${PYTHON_HOME}/include/${PYTHON_VER} \
	-I${FFMPEG_HOME}/include \
	-I${PROJECT_DIR}/third_party
export CGO_LDFLAGS = ${LDFLAGS} \
	-L${PYTHON_HOME}/lib -l${PYTHON_VER} \
	-L${FFMPEG_HOME}/lib \
	-L${LIBMEDIA_HOME}/lib \
	-lmedia -lavcodec -lavformat -lavutil -lswscale -lswresample
export GOPROXY := https://goproxy.cn
export PYTHONPATH=${PROJECT_DIR}/pkg

TARGET = gweb

all: $(TARGET)

$(TARGET): build
	mkdir -p deploy
	go build -v -o deploy/$(TARGET) cmd/gweb.go
	rsync -av --include="*/" --include="*.pyc" --exclude="*" \
		pkg/format deploy
ifeq ($(OS_TYPE), Linux)
	patchelf --set-rpath '$$ORIGIN' deploy/$@
	find ${PYTHON_HOME}/lib -maxdepth 1 \( -type f -o -type l \) \
		-regex '.*/lib${PYTHON_VER}\.so\.[0-9]+\.[0-9]+' \
		-exec cp -v {} deploy/ \;
	find ${FFMPEG_HOME}/lib -maxdepth 1 \( -type f -o -type l \) \
		-regextype posix-extended \
		-regex '.*/lib(avutil|avformat|avfilter|avcodec|swscale|swresample)\.so\.[0-9]+' \
		-exec cp -v {} deploy/ \;
	cp -v ${LIBMEDIA_HOME}/lib/libmedia.so deploy
else ifeq ($(OS_TYPE), Darwin)
	install_name_tool -add_rpath ${FFMPEG_HOME}/lib/ deploy/$(TARGET)
endif

build: formatui mediaui
	python -m compileall -b pkg/format
	go build -x -v -o $(TARGET) cmd/gweb.go
ifeq ($(OS_TYPE), Darwin)
	install_name_tool -add_rpath ${FFMPEG_HOME}/lib/ $(TARGET)
endif

env:
	@echo OS_TYPE=$(OS_TYPE)
	@echo PROJECT_DIR=$(PROJECT_DIR)
	@echo PYTHON_HOME=$(PYTHON_HOME)
	@echo PYTHON_VER=$(PYTHON_VER)
	@echo PYTHONPATH=$(PYTHONPATH)
	@echo CGO_CFLAGS=$(CGO_CFLAGS)
	@echo CGO_CXXFLAGS=$(CGO_CXXFLAGS)
	@echo CGO_LDFLAGS=$(CGO_LDFLAGS)
ifeq ($(OS_TYPE), Linux)
	@echo LD_LIBRARY_PATH=$(LD_LIBRARY_PATH)
else ifeq ($(OS_TYPE), Darwin)
	@echo DYLD_LIBRARY_PATH=$(DYLD_LIBRARY_PATH)
endif

init: env
	mkdir -p public/img public/js public/css public/plugins

formatui: init
	cp third_party/formatui/dist/img/formatui.svg public/img/format.svg
	cp third_party/formatui/dist/index.min.js public/js/format.min.js
	cp third_party/formatui/dist/index.min.css public/css/format.min.css
	cp -r third_party/formatui/dist/plugins/* public/plugins

mediaui: init
	cp third_party/mediaui/dist/mediaui.svg public/img/mediaui.svg
	cp third_party/mediaui/dist/mediaui.js public/js/mediaui.js
	cp third_party/mediaui/dist/mediaui.css public/css/mediaui.css

run:
	go run cmd/gweb.go -p 5015 --debug

tool: env
	go tool cgo -debug-gcc pkg/format/format.go
	#go tool cgo pkg/media/media.go

test: env
	pushd tests/format && \
	python -m unittest -v test_pytext.py && \
	python -m unittest -v test_pycmd.py && \
	python -m unittest -v test_pyfmt.py && \
	popd

clean: delformatui delmediaui
	find pkg -name "*.pyc" -type f -delete
	find pkg -type d -name "__pycache__" -exec rm -r {} +
	find tests -type d -name "__pycache__" -exec rm -r {} +
	find deploy/* -name "gweb.yaml" -prune -o -exec rm -rf {} +
	rm -rf gweb

delformatui:
	rm -rf public/img/format.svg \
		public/js/format.min.js \
		public/css/format.min.css \
		public/plugins/bootstrap-icons@* \
		public/plugins/clipboard@* \
		public/plugins/json5@* \
		public/plugins/w2ui@*

delmediaui:
	rm -rf public/img/mediaui.svg \
		public/js/mediaui.js \
		public/css/mediaui.css

.PYONY: all env init formatui mediaui run test clean
