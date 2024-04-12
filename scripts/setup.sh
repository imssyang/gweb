#!/bin/bash

PROJ_DIR=$(dirname $(dirname $(readlink -f "$0")))
SYSD_DIR=/etc/systemd/system
source ${PROJ_DIR}/scripts/utils.sh

enable_service() {
  symlink_create $PROJ_DIR/init/$1 $SYSD_DIR/$1
  systemctl enable $1
  systemctl daemon-reload
}

disable_service() {
  systemctl disable $1
  systemctl daemon-reload
  symlink_delete $SYSD_DIR/$1
}

case "$1" in
  install)
    enable_service gweb.service
    ;;
  uninstall)
    disable_service gweb.service
    ;;
  *)
    echo "Usage: ${0##*/} {install|uninstall}"
    ;;
esac

exit 0
