#!/bin/bash

# shellcheck disable=SC1090

# >>>>>>>> 定义变量 default variables >>>>>>>>>>
database_type="mariadb"      # 使用 mariadb 或 mysql 数据库 use mariadb or mysql database
database_version="10"        # 数据库版本 database version
utils_path="./utils/*sh"     # 工具路径 utils path
app_path="./app/*sh"         # 软件路径 app path
dns_addr="8.8.8.8"           # DNS 地址 define DNS address
deb_package="mariadb-server" # 安装包 define deb package
log_path="./logs/log.txt"    # 日志路径 define log path
mysql_bind_config_path="/etc/mysql/mariadb.conf.d/50-server.cnf" # mysql 配置路径 define mysql config path
# <<<<<<<< 定义变量 default variables <<<<<<<<<<

# >>>>>>>> 设定色彩 set color >>>>>>>>>>
. ./init/print_color.sh # 载入颜色印出函数 import color printing functions
. ./init/set_color.sh   # 载入颜色设置函数 import color setting functions
# <<<<<<<< 设定色彩 set color <<<<<<<<<<

# >>>>>>>> 载入工具 load utils >>>>>>>>>>
. ./init/basic.sh      # 载入基础函数 import basic functions
set_default_font_color 103 103 103 # 设置默认字体颜色 set default font color
load "${utils_path}"   # 载入工具包 load utils
load "${app_path}"     # 载入 app 设定工具包 load app config tool
# <<<<<<<< 载入工具 load utils <<<<<<<<<<

# >>>>>>>> 初始化设置 initialize >>>>>>>>>>
set_log "${log_path}" # 设定日志 set log
post "${log_path}"    # 执行离开后操作 post operation after exit
apt_install bc 2      # 安装 mariadb 数据库 install mariadb database
panic "$(check_database_version "${database_type}" "${database_version}" "check" && echo "continue" || echo "panic")"
# >>>>>>>> 初始化设置 initialize >>>>>>>>>>

# >>>>>>>> 设定容器环境 set debian env >>>>>>>>>>
set_dns "${dns_addr}" # 设定 dns 服务 set dns 服务
apt_update 3          # 更新 apt 源 update apt source
# <<<<<<<< 设定容器环境 set debian env <<<<<<<<<<

# >>>>>>>> 设定数据库 set mariadb >>>>>>>>>>
apt_install "${deb_package}" 2 # 安装 mariadb 数据库 install mariadb database
apt_clean                      # 清理 apt 源 clean apt source
replace "$mysql_bind_config_path" "bind-address.*" "bind-address=0.0.0.0" # 设定数据库对外网路位置 set mariadb bind ip address
check_database_version "${database_type}" "${database_version}" "correct"
# <<<<<<<< 设定数据库 set mariadb <<<<<<<<<<