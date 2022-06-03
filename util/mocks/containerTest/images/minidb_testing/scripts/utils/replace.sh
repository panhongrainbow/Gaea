#!/bin/bash

############################################################
# replace 为取代字符串
# replace is to replace string
#
# parameter 1: file path
# parameter 2: old string pattern
# parameter 3: new string pattern
# parameter 4: backup bool
#
replace () {
  # 如果文件不存在，则创建文件
  if [ ! -f "$1" ]; then
    print_xiaomi 0 "create file by using touch"
    print_list 3 "file: $1"
    print_list 3 "from: $2"
    print_list 3 "to: $3"

    mkdir -p "${1%/*}"
    touch "$1"
    result_code=$?
    if [ $result_code -eq 0 ]; then
      print_success 3 "create file successfully"
      echo "$3" >> "$1"
      print_success 3 "edit file successfully"
      return 0
    else
      # 档案不存在 file does not exist
      print_fail 3 "$1 does not exist"
      return 1
    fi
  fi

  # 如果文件存在，则取代字符串
  if [ -f "$1" ]; then
    print_xiaomi 0 "replace string by using sed"
    print_list 3 "file: $1"
    print_list 3 "from: $2"
    print_list 3 "to: $3"

    sed -i "s/$2/$3/g" "$1"

    # 检查是否安装成功 check if install successfully
    result_code=$?
    print_process 3 "result code: $result_code"

    if [ $result_code -eq 0 ]; then
      print_success 3 "replace string successfully"
      return 0
    else
      print_fail 3 "replace string failed"
      return 1
    fi
  else
    # 档案不存在 file does not exist
    print_fail 3 "$1 does not exist"
    return 1
  fi
}