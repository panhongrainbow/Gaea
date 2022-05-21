#!/bin/sh

# set_dns 为设置 dns 服务
# set_dns is to set dns service
set_dns () {
    echo "nameserver $1" # > /etc/resolv.conf
}