#!/bin/bash

############################################################
# check_or_install_database_version 为检查数据库版本是否支援
# check_or_install_database_version is to check whether the database version is supported
#
# parameter 1: database type (mariadb or mysql)
# parameter 2: database version
# parameter 3: correct or check
#
# return 0: supported
# return 1: not supported
#
check_or_install_database_version () {
  case $1 in
    mariadb)
      check_mariadb_version "$2" "$3"
      return $? # check by sub function
      ;;
    *)
      return 1 # not supported
      ;;
  esac
}

############################################################
# check_mariadb_version 为检查 mariadb 数据库版本是否支援
# check_mariadb_version is to check whether the database version is supported
#
# parameter 1: mariadb version
# parameter 2: correct or check
#
check_mariadb_version () {
  case $1 in
    10)
      compare_string "$2" "correct" && sh ./app/mariadb/correct_mariadb.sh || return 0
      return 0 # supported
      ;;
    *)
      return 1 # not supported
      ;;
  esac
}