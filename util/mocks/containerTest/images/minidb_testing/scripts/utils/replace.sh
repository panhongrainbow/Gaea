#!/bin/bash

############################################################
# replace 为取代字符串
# replace is to replace string
#
# parameter 1: file path (string)
# parameter 2: old string pattern (string)
# parameter 3: new string pattern (string)
#
replace () {
  print_xiaomi 0 "replace string by using sed" # 取代字符串
  print_list 3 "file: $1"
  print_list 3 "from: $2"
  print_list 3 "to: $3"

  if [ -f "$1" ]; then # 如果文件不存在 if file not exist
    # 取代字符串 replace string
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
  else # 如果文件存在 if file exist
    print_fail 3 "file not exist" # 文件不存在
    return 1
  fi
}