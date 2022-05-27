#!/bin/bash

############################################################
# check_database_version 为检查数据库版本是否支援
# check_database_version is to check whether the database version is supported
#
# parameter 1: database type (mariadb or mysql)
# parameter 2: database version
#
# return 0: supported
# return 1: not supported
#
check_database_version () {
  case $1 in
    mariadb)
      check_mariadb_version "$2"
      return $? # check by sub function
      ;;
    *)
      return 1 # not supported
      ;;
  esac
}

############################################################
# mariadb_version 为检查 mariadb 数据库版本是否支援
# check_database_version is to check whether the database version is supported
#
# parameter 1: mariadb version
#
check_mariadb_version () {
  case $1 in
    10)
      return 0 # supported
      ;;
    *)
      return 1 # not supported
      ;;
  esac
}