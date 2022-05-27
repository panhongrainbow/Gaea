#!/bin/bash

############################################################
# set_dns 为设置 dns 服务
# set_dns is to set dns service
#
# parameter 1: dns service address
#
set_dns () {
  print_xiaomi 0 "set dns name server"
  print_list 3 "dns address: $1"
  # if [ -f "/etc/resolv.conf" ]; then
  echo "nameserver $1" # > /etc/resolv.conf
  # fi
  print_success 2 "set dns name server $1 successfully"
}