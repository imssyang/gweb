.DEFAULT: env
.PYONY: clean

OS_TYPE := $(shell uname)
PROJECT_DIR=$(shell pwd)
PYTHON_HOME=$(shell pyenv prefix)
PYTHON_VER=$(shell ls ${PYTHON_HOME}/include)

ifeq ($(OS_TYPE), Linux)
	export LD_LIBRARY_PATH=${PYTHON_HOME}/lib
else ifeq ($(OS_TYPE), Darwin)
	export BASE_LDFLAGS=-Wl,-no_warn_duplicate_libraries
endif

export CGO_CFLAGS=\
	-I${PYTHON_HOME}/include/${PYTHON_VER}
export CGO_CXXFLAGS=\
	-I${PYTHON_HOME}/include/${PYTHON_VER} \
	-I${PROJECT_DIR}/third_party
export CGO_LDFLAGS=${BASE_LDFLAGS} \
	-L${PYTHON_HOME}/lib \
	-l${PYTHON_VER}
export PYTHONPATH=${PROJECT_DIR}/internal/api

formatui-deploy:
	mkdir -p public/img public/js public/css
	cp third_party/formatui/src/img/formatui.svg public/img/formatify.svg
	cp third_party/formatui/dist/index.min.js public/js/formatify.min.js
	cp third_party/formatui/dist/index.min.css public/css/formatify.min.css
	cp -r third_party/formatui/dist/plugins/* public/plugins

formatui-clean:
	rm -rf public/img/formatify.svg \
		public/js/formatify.min.js \
		public/css/formatify.min.css \
		public/plugins/bootstrap-icons@* \
		public/plugins/clipboard@* \
		public/plugins/json5@* \
		public/plugins/w2ui@*

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

run:
	go run cmd/gweb.go -p 5015 --debug

deploy: env formatui-deploy
	python -m compileall -b internal/api/formatify
	rsync -av --include="*/" --include="*.pyc" --exclude="*" \
		internal/api/formatify deploy
	go build -v -o deploy/gweb cmd/gweb.go
ifeq ($(OS_TYPE), Linux)
	patchelf --set-rpath '$$ORIGIN' deploy/gweb
	cp -v ${PYTHON_HOME}/lib/lib${PYTHON_VER}.so.1.0 deploy
endif

test:
	python -m unittest -v tests/formatify/test_pytext.py
	python -m unittest -v tests/formatify/test_pycmd.py
	python -m unittest -v tests/formatify/test_pyfmt.py

clean: formatui-clean
	find internal -name "*.pyc" -type f -delete
	find internal -type d -name "__pycache__" -exec rm -r {} +
	find tests -type d -name "__pycache__" -exec rm -r {} +
	rm -rf deploy/gweb \
		deploy/libpython* \
		deploy/formatify
