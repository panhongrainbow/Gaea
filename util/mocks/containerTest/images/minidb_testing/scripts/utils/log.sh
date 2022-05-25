#!/bin/sh

############################################################
# set_log 为设置日志
# set_log is to set log
#
# parameter 1: log path
#
set_log () {
  print_xiaomi 0 "set log path"
  print_list 3 "log path: $1"
  # 全部日志重定向到文件里 redirect all log to file
  exec 1>"$1" 2>&1
  print_success 2 "set log path successfully"
}