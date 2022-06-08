#!/bin/bash

############################################################
# apt_add_repo 为特定版本的数据库 repo
# apt_add_repo is to add repo for specific database version
#
# parameter 1: repo key (string) EX: https://mariadb.org/mariadb_release_signing_key.asc
# parameter 2: repository url (string) EX: "deb [arch=amd64,i386,arm64,ppc64el] https://mirrors.aliyun.com/mariadb/repo/10.8/debian bullseye main"
# parameter 3: retry count (integer) EX: 10
#
apt_add_repo () {
  print_xiaomi 0 "apt add repo of ubuntu"
  for i in $(seq 1 1 "$3")
  do
    # 安裝 apt-key 工具; install apt-key tool
    apt-get install -y apt-transport-https curl # 在一开始安装这些软件 add these packages at the beginning

    # 检查是否安装成功; check if install successfully
    result_code=$?
    retry_count=$i
    print_process 3 "result code: $result_code"
    print_process 3 "retry count: $retry_count"
    if [ $result_code -eq 0 ]; then
      print_success 3 "apt install package in $retry_count time(s) successfully"
      break
    else
      print_fail 3 "apt install package in $retry_count time(s) failed"
      if [ "$retry_count" -eq "$3" ]; then
        return 1
      fi
    fi
  done

  for i in $(seq 1 1 "$3")
  do
    # 取得秘钥; get key
    # echo ""
    curl -o "/etc/apt/trusted.gpg.d/${1##*/}" "$1"

    # 是否取得金钥成功; check if fetch key successfully
    result_code=$?
    retry_count=$i
    print_process 3 "result code: $result_code"
    print_process 3 "retry count: $retry_count"
    if [ $result_code -eq 0 ]; then
      print_success 3 "fetch key in $retry_count time(s) successfully"
      break
    else
      print_fail 3 "fetch key in $retry_count time(s) failed"
      if [ "$retry_count" -eq "$3" ]; then
        return 1
      fi
    fi
  done

  for i in $(seq 1 1 "$3")
  do
    # 添加 repo; add repo
    sh -c "echo $2 >>/etc/apt/sources.list"

    # 是否新增 repo 成功; check if add repo successfully
    result_code=$?
    retry_count=$i
    print_process 3 "result code: $result_code"
    print_process 3 "retry count: $retry_count"
    if [ $result_code -eq 0 ]; then
      print_success 3 "add repo in $retry_count time(s) successfully"

      apt-get purge -y --auto-remove apt-transport-https curl # 在最后移除这些软件 remove these packages at the end
      if [ $result_code -eq 0 ]; then
        print_success 3 "remove software-properties-common successfully"
      fi

      break
    else
      print_fail 3 "add repo in $retry_count time(s) failed"
      if [ "$retry_count" -eq "$3" ]; then
        return 1
      fi
    fi
  done

  return 0
}
