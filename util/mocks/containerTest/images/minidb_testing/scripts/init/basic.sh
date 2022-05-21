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