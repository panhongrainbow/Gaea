#!/bin/bash

# 以下为对 mariadb 数据库修正

print_xiaomi 0 "correct mariadb database"

print_list 3 "make xiaomi directory"

mkdir -p /home/xiaomi/

print_list 3 "create mysqld init script"

cat << EOF > /home/xiaomi/mysqld_init.sh
#!/bin/bash

# 一般原本 Docker 的数据库容器无法在 Containerd 上执行，进行以下修正
mkdir /var/run/mysqld
# RUN useradd -m mysql
chown mysql:mysql /var/run/mysqld
chmod 777 /var/run/mysqld

# user.sql 为一开始执行 mysqld 服务时，所需要执行的 SQL 脚本，会创建一个用户 xiaomi，并且设置密码
mysqld --init-file=/home/xiaomi/user.sql
EOF

print_list 3 "create user script"

cat << EOF > /home/xiaomi/user.sql
CREATE USER 'xiaomi'@'%' IDENTIFIED BY '12345';
GRANT ALL PRIVILEGES ON *.* TO 'xiaomi'@'%' WITH GRANT OPTION;
FLUSH PRIVILEGES;
EOF

print_list 3 "make mysqld init script executable"

chmod +x /home/xiaomi/mysqld_init.sh
