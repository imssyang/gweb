#!/bin/bash

PROJ_DIR=$(dirname $(dirname $(readlink -f "$0")))
PROJ_NAME=$(basename $PROJ_DIR)

symlink_create() {
  if [[ ! -d $2 ]] && [[ ! -s $2 ]]; then
    ln -s $1 $2
    echo "($PROJ_NAME) create symlink: $2 -> $1"
  fi
}

symlink_delete() {
  if [[ -d $1 ]] || [[ -s $1 ]]; then
    rm -rf $1
    echo "($PROJ_NAME) delete symlink: $1"
  fi
}

