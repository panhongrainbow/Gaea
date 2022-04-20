FROM debian:latest

# 安装数据库
RUN apt-get update
RUN apt-get install -y mariadb-server mariadb-client

# 修改數据库設定
RUN sed -i "s/bind-address.*/bind-address=0.0.0.0/g" /etc/mysql/mariadb.conf.d/50-server.cnf

# 设定用户密码
RUN mkdir -p /home/mysql/
ADD user.sql /home/mysql/

# 进行修正
ADD mysqld_init.sh /home/mysql/
RUN chmod +x /home/mysql/mysqld_init.sh

# 启动数据库
ENTRYPOINT ["/home/mysql/mysqld_init.sh"]
