#!/bin/bash

############################################################
# panic 为错误中断所有操作
# panic is to stop all operation immediately because of error
#
# parameter 1:
#
panic () {
  if [ "$1" = "panic" ]; then
    echo "panic - processing"
    print_fail 3 "panic - processing"
    exit
  fi
}