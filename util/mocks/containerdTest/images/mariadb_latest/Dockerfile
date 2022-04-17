FROM debian:latest

# 安装数据库
RUN apt-get update
RUN apt-get install -y mariadb-server mariadb-client

# 进行修正
RUN mkdir /var/run/mysqld
# RUN useradd -m mysql
RUN chown mysql:mysql /var/run/mysqld
RUN chmod 777 /var/run/mysqld

# 启动数据库
ENTRYPOINT ["mysqld"]
