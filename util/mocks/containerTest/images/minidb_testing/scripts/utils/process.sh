#!/bin/bash

############################################################
# panic 为错误中断所有操作
# panic is to stop all operation immediately because of error
#
# parameter 1:
#
panic () {
  echo ">>>>>>>>>>>>"
  echo "$1"
  if [ "$1" = "panic" ]; then
    echo "panic - processing"
    print_fail 2 "panic - processing"
    exit
  fi
}