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
  print_xiaomi 0 "replace string by using sed"
  print_list 3 "file: $1"
  print_list 5 "from: $2"
  print_list 5 "to: $3"

  if [ -f "$1" ]; then
    sed -i "s/$2/$3/g" "$1"

    # 检查是否安装成功 check if install successfully
    result_code=$?
    print_detail 3 "result code: $result_code"

    if [ $result_code -eq 0 ]; then
      print_success 2 "replace string successfully"
    else
      print_fail 2 "replace string failed"
    fi
  else
    # 档案不存在 file does not exist
    print_fail 2 "$1 does not exist"
  fi
}