#!/bin/sh

############################################################
# set_log 为设置日志
# set_log is to set log
#
# parameter 1: log path
#
set_log () {
  # 全部日志重定向到文件里 redirect all log to file
  exec 1>"$1" 2>&1
}