#!/bin/bash

############################################################
# load 为载入一些基础的函数
# load is to load some basic functions
#
# parameter 1: utils path
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
# post 为执行最后操作
# post is to post operation
#
# parameter 1: log path
#
exit_safely () {
  # shellcheck disable=SC2046
  exec > $(tty)
  tail -n 20 "$1"
}
post () {
  trap 'exit_safely '"$1" EXIT INT QUIT ILL
}

############################################################
# compare_string 为用于比较字符串
# compare_string is to compare string
#
# parameter 1: string variable
# parameter 2: wanted string
#
compare_string () {
  if [ "$1" = "$2" ]; then
    return 0
  else
    return 1
  fi
}