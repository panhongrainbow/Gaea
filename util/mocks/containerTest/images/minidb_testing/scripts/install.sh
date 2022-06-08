#!/bin/bash

# >>>>> 定义变量 default variables
database_type="mariadb"         # 使用 mariadb 或 mysql 数据库; use mariadb or mysql database
database_version="10.8"         # 数据库版本; database version
database_repo_key="https://mariadb.org/mariadb_release_signing_key.asc"
                                # 数据库 repo 密钥; database repo key
database_repo='deb [arch=amd64,arm64,ppc64el] https://mirrors.aliyun.com/mariadb/repo/10.8/ubuntu impish main'
                                # 数据库 repo 配置; database repo configuration
utils_path="./utils/*sh"        # 工具路径; utils path
app_path="./app/*sh"            # 软件路径; app path
dns_addr="8.8.8.8"              # 名称服务器地址; define dns address
deb_package="mariadb-server-${database_version}"
                                # 定义要安装的套件; define deb package
mysql_bind_config_path="/etc/mysql/mariadb.conf.d/50-server.cnf"
                                # mysql 配置路径; define mysql config path

# >>>>> 载入工具 load utils
. ./init/print_color.sh         # 载入颜色印出函数; import color printing functions
. ./init/basic.sh               # 载入基础函数; import basic functions
load "${utils_path}"            # 载入工具包; load utils
load "${app_path}"              # 载入 app 设定工具包; load app config tool

# >>>>> 环境初始化设置 environment initialize
set_dns "${dns_addr}"           # 设定名称服务器位罝; set dns services
apt_update 10                   # 更新 apt 源; update apt source
apt_add_repo "$database_repo_key" "$database_repo" 10
                                # 添加数据库 apt 源; add apt database source
apt_update 10                   # 更新 apt 源; update apt source
check_or_correct_database_version "${database_type}" "${database_version}" "${mysql_bind_config_path}" "check"
                                # 检查数据库版本; check database version

# >>>>> 设定数据库 set mariadb
apt_install "${deb_package}" 10 # 安装 mariadb 数据库; install mariadb database
apt_clean                       # 清理 apt 快取; clean apt cache
create_user "mysql"             # 创建 mysql 用户; create mysql user
replace "$mysql_bind_config_path" "bind-address.*" "bind-address=0.0.0.0"
                                # 设定数据库对外网路位置; set mariadb bind ip address
check_or_correct_database_version "${database_type}" "${database_version}" "${mysql_bind_config_path}" "correct"
                                # 检查数据库版本; check database version
create_user "mysql"             # 创建 mysql 用户; create mysql user

# >>>>> 安装完成和清除 install finished and clean container
clean_mariadb_container         # 清除数据库容器不必要的档案; clean mariadb container unnecessary files
clean_container                 # 清除容器不必要的档案; clean container unnecessary files
