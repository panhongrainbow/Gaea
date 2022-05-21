#!/bin/sh

# shellcheck disable=SC2181

############################################################
# apt_update 为更新 apt repo
# apt_update is to update apt repo
#
# parameter 1: retry count
#
apt_update () {
  echo "${ORANGE}[XiaoMi] apt update start !${NC}"
  for i in $(seq 1 1 "$1")
  do
    apt-get update > /dev/null 2>&1
    if [ $? -eq 0 ]; then
      echo "${BLUE}  => apt update in $i time successfully !${NC}"
      return 0
    else
      echo "${RED}  => apt update in $i time failed !${NC}"
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
    echo "${ORANGE}[XiaoMi] apt update install package !${NC}"

    for i in $(seq 1 1 "$2")
    do
      apt-get install -y --no-install-recommends "$1" > /dev/null 2>&1
      if [ $? -eq 0 ]; then
        echo "${BLUE}  => apt install package in $i time successfully !${NC}"
        break
      else
        echo "${RED}  => apt install package in $i time failed !${NC}"
        if [ $? -eq "$2" ]; then
          break
        fi
      fi
    done

    echo "${PURPLE}    => installed: $1 !${NC}"
    return 0
}