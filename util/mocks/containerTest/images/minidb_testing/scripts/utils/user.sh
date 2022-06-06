#!/bin/bash

############################################################
# create_user 为创建一个用户和群组
# create_user is to create a user and a group
#
# parameter 1: user name (string)
#
create_user () {
  id "$1"
  result_code=$?
  if [ $result_code -eq 0 ]; then
    groupadd "$1"
    useradd -g "$1" -s /bin/false "$1"
  fi
}