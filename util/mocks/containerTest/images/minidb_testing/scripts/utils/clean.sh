#!/bin/bash

############################################################
# clean_container 为用来最后清除容器较不必要的档案
# clean_container is to remove unnecessary files inside container
#
# no parameters
#
clean_container () {
  print_xiaomi 0 "clean container"

  # 清除 apt 設置 remove apt settings
  print_success 3 "apt purge apt successfully"
  apt-get autoclean
  apt-get autoremove
  apt_purge "apt*"
  rm -rf /var/lib/apt/lists/*

  # 清除文檔 remove doc files
  print_success 3 "remove doc successfully"
  rm -rf /usr/share/doc/*

  # 清除日誌 remove log files
  print_success 3 "remove log successfully"
  rm -rf /var/log/*
}