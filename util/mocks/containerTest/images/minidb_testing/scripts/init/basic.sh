#!/bin/sh

# shellcheck disable=SC1090

# load 为载入一些基础的函数
# load is to load some basic functions
load () {
    for i in $1
    do
        if [ -e "${i}" ]; then
            echo "load - processing $i"
            . "${i}"
        fi
    done
}

############################################################
# post 为执行最后操作
# post is to post operation
#
# no parameters
#
exit_safely (){
  # shellcheck disable=SC2046
  exec > $(tty)
  tail -n 20 ./logs/log.txt
}
post() {
  trap exit_safely EXIT
}