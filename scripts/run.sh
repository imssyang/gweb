#!/bin/bash

PROJECT_DIR=/opt/app/gweb
PYTHON_HOME=/opt/python/pyenv/versions/3.9.18

pymain() {
  PYTHONPATH=${PYTHONPATH}:${PROJECT_DIR}/third_party \
  python ${PROJECT_DIR}/tests/pymain.py
}

gomain() {
  LD_LIBRARY_PATH=${LD_LIBRARY_PATH}:${PYTHON_HOME}/lib \
  PYTHONPATH=${PYTHONPATH}:${PROJECT_DIR}/third_party:${PROJECT_DIR}/internal/api/formatify \
  go run ${PROJECT_DIR}/cmd/gweb.go
}

gomain "$@"
