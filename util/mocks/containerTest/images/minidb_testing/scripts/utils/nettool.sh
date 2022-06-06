#!/bin/bash

############################################################
# set_dns 为设置 dns 服务
# set_dns is to set dns service
#
# parameter 1: dns service address (string)
#
set_dns () {
  print_xiaomi 0 "set dns name server"
  print_list 3 "dns address: $1"
  echo "nameserver $1" > /etc/resolv.conf
  print_success 3 "set dns name server $1 successfully"
}