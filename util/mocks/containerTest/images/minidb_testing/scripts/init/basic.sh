#!/bin/bash

############################################################
# load 为载入一些基础的函数
# load is to load some basic functions
#
# parameter 1: utils path (string)
#
load () {
  for i in $1
  do
    if [ -e "${i}" ]; then
      echo "load - processing $i"
      # shellcheck source=/dev/null
      . "${i}"
    fi
  done
}

############################################################
# compare_string 为用于比较字符串
# compare_string is to compare string
#
# parameter 1: string variable (string)
# parameter 2: wanted string (string)
#
compare_string () {
  if [ "$1" = "$2" ]; then
    return 0
  else
    return 1
  fi
}