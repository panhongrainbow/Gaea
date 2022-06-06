#!/bin/bash

############################################################
# clean_mariadb_container 为用来最后清除數據庫容器较不必要的档案
# clean_mariadb_container is to remove unnecessary files inside mariadb container
#
# no parameters
#
clean_mariadb_container () {
  # 清除數據庫不必要的檔案 remove unnecessary files
  print_xiaomi 0 "clean mariadb library"
  rm -rf /var/lib/mysql/*
}