.DEFAULT: env
.PYONY: clean

PROJECT_DIR=$(shell pwd)
PYTHON_HOME=$(shell pyenv prefix)
PYTHON_VER=$(shell ls ${PYTHON_HOME}/include)
export CGO_CFLAGS=\
	-I${PYTHON_HOME}/include/${PYTHON_VER}
export CGO_CXXFLAGS=\
	-I${PYTHON_HOME}/include/${PYTHON_VER} \
	-I${PROJECT_DIR}/third_party
export CGO_LDFLAGS=\
	-L${PYTHON_HOME}/lib \
	-l${PYTHON_VER}
export LD_LIBRARY_PATH=${PYTHON_HOME}/lib
export PYTHONPATH=${PROJECT_DIR}/internal/api

env:
	@echo PROJECT_DIR=$(PROJECT_DIR)
	@echo PYTHON_HOME=$(PYTHON_HOME)
	@echo PYTHON_VER=$(PYTHON_VER)
	@echo PYTHONPATH=$(PYTHONPATH)
	@echo CGO_CFLAGS=$(CGO_CFLAGS)
	@echo CGO_CXXFLAGS=$(CGO_CXXFLAGS)
	@echo CGO_LDFLAGS=$(CGO_LDFLAGS)
	@echo LD_LIBRARY_PATH=$(LD_LIBRARY_PATH)

run:
	go run cmd/gweb.go -p 5015 --debug

deploy: env
	go build -v -o deploy/gweb cmd/gweb.go
	patchelf --set-rpath '$$ORIGIN' deploy/gweb
	python -m compileall internal/api/formatify
	cp -v ${PYTHON_HOME}/lib/lib${PYTHON_VER}.so.1.0 deploy
	rsync -av --include="*/" --include="*.pyc" --exclude="*" internal/api/formatify deploy

test:
	python -m unittest -v tests/formatify/test_pytext.py
	python -m unittest -v tests/formatify/test_pycmd.py

clean:
	rm -rf deploy/gweb deploy/libpython* deploy/formatify
	find internal -name "*.pyc" -type f -delete
	find internal -type d -name "__pycache__" -exec rm -r {} +
	find tests -type d -name "__pycache__" -exec rm -r {} +
