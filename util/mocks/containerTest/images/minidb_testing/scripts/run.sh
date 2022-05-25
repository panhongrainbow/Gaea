#!/bin/sh

# shellcheck disable=SC1090

# >>>>>>>> 定义变量 default variables >>>>>>>>>>
utils_path="./utils/*sh"     # 工具路径 utils path
mariadb_path="./mariadb/*sh" # mariadb 设定工具包 mariadb config path
dns_addr="8.8.8.8"           # DNS 地址 define DNS address
deb_package="mariadb-server" # 安装包 define deb package
log_path="./logs/log.txt"    # 日志路径 define log path
mysql_bind_config_path="/etc/mysql/mariadb.conf.d/50-server.cnf" # mysql 配置路径 define mysql config path
# <<<<<<<< 定义变量 default variables <<<<<<<<<<

# >>>>>>>> 载入工具 load utils >>>>>>>>>>
. ./init/basic.sh      # 载入基础函数 import basic functions
load "${utils_path}"   # 载入工具包 load utils
load "${mariadb_path}" # 载入 mariadb 设定工具包 load mariadb config tool
# <<<<<<<< 载入工具 load utils <<<<<<<<<<

# >>>>>>>> 初始化设置 initialize >>>>>>>>>>
set_log "${log_path}" # 设定日志 set log
post "${log_path}"    # 执行最后操作 post operation
# >>>>>>>> 初始化设置 initialize >>>>>>>>>>

# 以下为暂时测试
# apt_update 3
apt_install "vim nano" 3
replace ./tmp BBC CNN
echo -ne "\033]10;#676767\007"
return 0

# >>>>>>>> 设定容器环境 set debian env >>>>>>>>>>
set_dns "${dns_addr}" # 设定 dns 服务 set dns 服务
# <<<<<<<< 设定容器环境 set debian env <<<<<<<<<<

# >>>>>>>> 设定数据库 set mariadb >>>>>>>>>>
apt_install "${deb_package}" # 安装 mariadb 数据库 install mariadb database
# <<<<<<<< 设定数据库 set mariadb <<<<<<<<<<

# >>>>>>>> 设定数据库对外网路位置 set mariadb bind ip address >>>>>>>>>>
replace "$mysql_bind_config_path" "bind-address.*" "bind-address=0.0.0.0"
# <<<<<<<<<< 设定数据库对外网路位置 set mariadb bind ip address <<<<<<<<<<

return 0