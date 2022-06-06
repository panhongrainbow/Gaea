#!/bin/bash

############################################################
# correct_mariadb_container 为修正数据库容器的脚本，使其能够在 containerd 正常运行
# correct_mariadb_container is to correct mariadb container, make it workable in containerd normally
#
# parameter 1: no parameter
#
correct_mariadb_container () {
  # 以下为对 mariadb 数据库容器修正 correct the mariadb container in the following
  print_xiaomi 0 "correct mariadb database"
  print_list 3 "make xiaomi and mysql directory"

  # 建立 xiaomi 家目录 create xiaomi home directory
  mkdir -p /home/xiaomi/

  # 创建 mysqld 初始化脚本 create mysqld init script
  print_list 3 "create mysqld init script"
cat << EOF > /home/xiaomi/mysqld_init.sh
#!/bin/bash

# 有些数据库容器无法在 containerd 上执行，进行以下修正
# correct the database container because it is not workable on containerd

mkdir -p /var/run/mysqld
mkdir -p /var/lib/mysql/

chown mysql:mysql /var/run/mysqld
chown mysql:mysql /var/lib/mysqld

chmod 777 /var/run/mysqld

# user.sql 为一开始执行 mysqld 服务时，所需要执行的 SQL 脚本，会创建一个用户 xiaomi，并且设置密码
# user.sql is the SQL script to create a user xiaomi and set the password
mariadb --user=mysql --init-file=/home/xiaomi/user.sql
EOF

  # 建立数据库帐户和密码 create account and password for database
  print_list 3 "create user script"
cat << EOF > /home/xiaomi/user.sql
-- 建立帐户和密码 create account and password
CREATE USER 'xiaomi'@'%' IDENTIFIED BY '12345';
GRANT ALL PRIVILEGES ON *.* TO 'xiaomi'@'%' WITH GRANT OPTION;
FLUSH PRIVILEGES;
EOF

  # 让脚本可以执行 make the script executable
  print_list 3 "make mysqld init script executable"
  chmod +x /home/xiaomi/mysqld_init.sh
}