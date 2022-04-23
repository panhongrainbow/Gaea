#!/bin/bash

# 进行修正
mkdir /var/run/mysqld
# RUN useradd -m mysql
chown mysql:mysql /var/run/mysqld
chmod 777 /var/run/mysqld

# 启动数据库
mysqld --init-file=/home/mariadb/user.sql
