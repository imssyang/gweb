PROJECT_DIR=$(shell pwd)
PYTHON_HOME=$(shell pyenv prefix)
PYTHON_VER=$(shell ls ${PYTHON_HOME}/include)
FFMPEG_HOME=/opt/ffmpeg
OS_TYPE := $(shell uname)

export CGO_CFLAGS = -Wall -Wextra -O2 \
	-I${PYTHON_HOME}/include/${PYTHON_VER} \
	-I${FFMPEG_HOME}/include \
	-I${PROJECT_DIR}/third_party
export CGO_CXXFLAGS = -std=c++20 -O2 \
	-I${PYTHON_HOME}/include/${PYTHON_VER} \
	-I${FFMPEG_HOME}/include \
	-I${PROJECT_DIR}/third_party
export CGO_LDFLAGS = -Wl,-no_warn_duplicate_libraries \
	-L${PYTHON_HOME}/lib -l${PYTHON_VER} \
	-L${FFMPEG_HOME}/lib \
	-lavcodec -lavformat -lavutil -lswscale -lswresample
export PYTHONPATH=${PROJECT_DIR}/pkg
ifeq ($(OS_TYPE), Linux)
	export LD_LIBRARY_PATH=${PYTHON_HOME}/lib:${FFMPEG_HOME}/lib
endif

TARGET = gweb

all: $(TARGET)

$(TARGET): formatui mediaui
	python -m compileall -b pkg/format
	rsync -av --include="*/" --include="*.pyc" --exclude="*" \
		pkg/format deploy
	go build -v -o deploy/$@ cmd/gweb.go
ifeq ($(OS_TYPE), Linux)
	patchelf --set-rpath '$$ORIGIN' deploy/$@
	cp -v ${PYTHON_HOME}/lib/lib${PYTHON_VER}.so.1.0 deploy
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
	mkdir -p public/img public/js public/css

formatui: init
	cp third_party/formatui/src/img/formatui.svg public/img/format.svg
	cp third_party/formatui/dist/index.min.js public/js/format.min.js
	cp third_party/formatui/dist/index.min.css public/css/format.min.css
	cp -r third_party/formatui/dist/plugins/* public/plugins

mediaui: init
	cp third_party/mediaui/src/img/mediaui.svg public/img/media.svg
	cp third_party/mediaui/dist/index.min.js public/js/media.min.js
	cp third_party/mediaui/dist/index.min.css public/css/media.min.css

run:
	go run cmd/gweb.go -p 5015 --debug

test: env
	python -m unittest -v tests/format/test_pytext.py
	python -m unittest -v tests/format/test_pycmd.py
	python -m unittest -v tests/format/test_pyfmt.py

clean: clean-formatui clean-mediaui
	find pkg -name "*.pyc" -type f -delete
	find pkg -type d -name "__pycache__" -exec rm -r {} +
	find tests -type d -name "__pycache__" -exec rm -r {} +
	rm -rf deploy/gweb \
		deploy/libpython* \
		deploy/format

clean-formatui:
	rm -rf public/img/format.svg \
		public/js/format.min.js \
		public/css/format.min.css \
		public/plugins/bootstrap-icons@* \
		public/plugins/clipboard@* \
		public/plugins/json5@* \
		public/plugins/w2ui@*

clean-mediaui:
	rm -rf public/img/media.svg \
		public/js/media.min.js \
		public/css/media.min.css

.PYONY: all env init formatui mediaui run test clean