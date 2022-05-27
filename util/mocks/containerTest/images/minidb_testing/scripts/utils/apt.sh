#!/bin/bash

############################################################
# apt_update 为更新 apt repo
# apt_update is to update apt repo
#
# parameter 1: retry count
#
apt_update () {
  print_xiaomi 0 "apt update start"
  for i in $(seq 1 1 "$1")
  do
    apt-get update
    if [ $? -eq 0 ]; then
      print_success 2 "apt update in $i time successfully"
      return 0
    else
      print_fail 2 "apt update in $i time failed"
    fi
  done
}

############################################################
# apt_install 为安装 apt 包
# apt_install is to install apt packages
#
# parameter 1: debian package(s) (to split packages separated By blank)
# parameter 2: retry count
#
apt_install () {
  print_xiaomi 0 "apt install package(s)"

  # 重试进行安装 install again and again if failed
  for i in $(seq 1 1 "$2")
  do
    # $1 不加双引号不会报错
    apt-get install -y --no-install-recommends $1

    # 检查是否安装成功 check if install successfully
    result_code=$?
    retry_count=$i
    print_detail 3 "result code: $result_code"
    print_detail 3 "retry count: $retry_count"

    if [ $result_code -eq 0 ]; then
      print_success 2 "apt install package in $i time(s) successfully"
      break
    else
      print_fail 2 "apt install package in $i time(s) failed"
      if [ "$retry_count" -eq "$2" ]; then
        break
      fi
    fi
  done

  # 印出套件列表 print list
  print_list 5 "installed: $1"
  return 0
}

############################################################
# apt_clean 为清理 apt 下载包
# apt_clean is to clean apt packages
#
# no parameters
#
apt_clean () {
  print_xiaomi 0 "apt clean"
  apt-get clean
  print_success 2 "apt clean successfully"
  return 0
}