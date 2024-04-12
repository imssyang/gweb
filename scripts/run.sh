#!/bin/bash

PROJECT_DIR=/opt/app/gweb
PYTHON_HOME=/opt/python/pyenv/versions/3.9.18

pytest() {
  export PYTHONPATH=${PYTHONPATH}:${PROJECT_DIR}/internal/api
  python -m unittest -v ${PROJECT_DIR}/tests/formatify/test_pytext.py
  python -m unittest -v ${PROJECT_DIR}/tests/formatify/test_pycmd.py
}

gomain() {
  PYTHONPATH=${PYTHONPATH}:${PROJECT_DIR}/internal/api \
  LD_LIBRARY_PATH=${LD_LIBRARY_PATH}:${PYTHON_HOME}/lib \
  CGO_CFLAGS="-I${PYTHON_HOME}/include/python3.9" \
  CGO_CXXFLAGS="-I${PYTHON_HOME}/include/python3.9 -I${PROJECT_DIR}/third_party" \
  CGO_LDFLAGS="-L${PYTHON_HOME}/lib -lpython3.9" \
  go run ${PROJECT_DIR}/cmd/gweb.go
}

case "$1" in
  gomain)
    gomain "$@"
    ;;
  pytest)
    pytest "$@"
    ;;
  *)
    echo "Usage: ${0##*/} {gomain|pytest}"
    ;;
esac

exit 0
