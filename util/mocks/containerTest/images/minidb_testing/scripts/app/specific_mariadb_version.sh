#!/bin/bash

############################################################
# apt_add_repo 为特定版本的数据库 repo
# apt_add_repo is to add repo for specific database version
#
# parameter 1: debian version
# parameter 2: mariadb version
# parameter 3: retry count
#
apt_add_repo () {
  print_xiaomi 0 "apt add repo of debian $1 and mariadb $2"
  for i in $(seq 1 1 "$3")
  do
    apt-get install -y software-properties-common dirmngr apt-transport-https

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
      if [ "$retry_count" -eq "$3" ]; then
        return 1
      fi
    fi
  done

  for i in $(seq 1 1 "$3")
  do
    apt-key adv --fetch-keys 'https://mariadb.org/mariadb_release_signing_key.asc'

    # 是否取得金钥成功 check if fetch key successfully
    result_code=$?
    retry_count=$i
    print_detail 3 "result code: $result_code"
    print_detail 3 "retry count: $retry_count"
    if [ $result_code -eq 0 ]; then
      print_success 2 "fetch key in $i time(s) successfully"
      break
    else
      print_fail 2 "fetch key in $i time(s) failed"
      if [ "$retry_count" -eq "$3" ]; then
        return 1
      fi
    fi
  done

  for i in $(seq 1 1 "$3")
  do
    add-apt-repository "deb [arch=amd64,i386,arm64,ppc64el] https://mirrors.aliyun.com/mariadb/repo/$2/debian $1 main"
    # 是否新增 repo 成功 check if add repo successfully
    result_code=$?
    retry_count=$i
    print_detail 3 "result code: $result_code"
    print_detail 3 "retry count: $retry_count"
    if [ $result_code -eq 0 ]; then
      print_success 2 "add repo in $i time(s) successfully"
      break
    else
      print_fail 2 "add repo in $i time(s) failed"
      if [ "$retry_count" -eq "$3" ]; then
        return 1
      fi
    fi
  done

  return 0
}
