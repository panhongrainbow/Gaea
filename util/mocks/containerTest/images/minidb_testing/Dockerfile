FROM debian:bullseye-slim

# 安置脚本 add scripts
ADD ./scripts /tmp/scripts

# 安装 Install
RUN cd /tmp/scripts/; sh /tmp/scripts/install.sh; rm -rf /tmp/scripts

# 开埠 open port
EXPOSE 3306

# 启动数据库 Start
ENTRYPOINT ["/home/xiaomi/mysqld_init.sh"]
