#!/bin/bash

############################################################
# 色彩定莪表
# color table
#
# Black='\033[0;30m'       # Black (黑色) 0;30
RED='\033[0;31m'         # Red (红色) 0;31
# GREEN='\033[0;32m'       # Green (绿色) 0;32
# BROWN='\033[0;33m'       # Brown/Orange (橙色) 0;33
ORANGE='\033[0;33m'      # Brown/Orange (橙色) 0;33
BLUE='\033[0;34m'        # Blue (蓝色) 0;34
PURPLE='\033[0;35m'      # Purple (紫色) 0;35
# CYAN='\033[0;36m'        # Cyan (青色) 0;36
# LightGRAY='\033[0;37m'   # Light Gray (浅灰) 0;37
# DarkGRAY='\033[1;30m'    # Dark Gray (深灰) 1;30
# LightRed='\033[1;31m'    # Light Red (浅红) 1;31
# LightGreen='\033[1;32m'  # Light Green (浅绿) 1;32
# YELLOW='\033[1;32m'      # Yellow (黄色) 1;33
# LightBlue='\033[1;34m'   # Light Blue (浅蓝色) 1;34
# LightPurple='\033[1;35m' # Light Purple (浅紫) 1;35
LightCyan='\033[1;36m'   # Light Cyan (浅青)  1;36
# WHITE='\033[1;37m'       # White (白色) 1;37
NC='\033[0m'             # No Color (无色)

############################################################
# 印出空白
# print blank
#
# parameter 1: blank count in the beginning
#
print_blank () {
  # 印出空白 print black
  for i in $(seq 1 1 "$1")
  do
    printf " "
  done
}

############################################################
# 打印彩色标题
# print color xiaomi title
#
# parameter 1: blank count in the beginning
# parameter 2: title
#
print_xiaomi () {
  # 印出空白 print black
  print_blank "$1"

  # 印出彩色标题 print color title
  echo "${ORANGE}[XiaoMi] $2 !${NC}"
}

############################################################
# 打印执行成功的彩色信息
# print color
#
# parameter 1: blank count in the beginning
# parameter 2: content
#
print_success () {
  # 印出空白 print black
  print_blank "$1"

  # 印出彩色内容 print color content
  echo "${BLUE}=> $2 !${NC}"
}

############################################################
# 打印执行错误的彩色信息
# print fail
#
# parameter 1: blank count in the beginning
# parameter 2: content
#
print_fail () {
  # 印出空白 print black
  print_blank "$1"

  # 印出彩色内容 print color content
  echo "${RED}=> $2 !${NC}"
}

############################################################
# 打印色彩明细
# print list
#
print_list () {
  # 印出空白 print black
  print_blank "$1"

  # 印出彩色内容 print color content
  echo "${PURPLE}=> $2 !${NC}"
}

############################################################
# 打印执行细节
# print process
#
print_process () {
  # 印出空白 print black
  print_blank "$1"

  # 印出彩色内容 print color content
  echo "${LightCyan}=> $2 !${NC}"
}