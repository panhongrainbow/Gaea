FROM debian:bullseye-slim

# 安装数据库
RUN echo "nameserver 8.8.8.8" > /etc/resolv.conf
RUN apt-get update && apt-get install -y --no-install-recommends mariadb-server && apt-get clean

# 修改数据库连线设置
# 修改数据库连线设置 (正规表达式为 bind-address(\s*?)=(\s*?)127\.0\.0\.1)
RUN sed -i "s/bind-address.*/bind-address=0.0.0.0/g" /etc/mysql/mariadb.conf.d/50-server.cnf

# 设定用户密码和修正
RUN mkdir -p /home/mariadb/
ADD mariadb /home/mariadb/

# 清除 man file
RUN rm -rf /usr/share/man/??;rm -rf /usr/share/man/??_*

# 进行修正
RUN chmod +x /home/mariadb/mysqld_init.sh

# 启动数据库
ENTRYPOINT ["/home/mariadb/mysqld_init.sh"]
