# 测试用容器镜像列表

> 以下为目前在测试时，列出所使用的容器，包含打包的方式

## 实验用的容器像

### 重新打包 Debian MariaDB 数据库镜像

> 因为现有的数据库容器都无法在 Containerd 上正常启动，所以要进行修正和重新打包，**使用以下命令重新打包镜像**

#### 创建帐户资料档案

档案位于 Gaea/util/mocks/containerTest/images/mariadb_testing/mariadb/user.sql

```sql
-- SQL 档案，当容器启动时，会立即执行以下命令，就会创建新用户 xiaomi ，和 root 相同权限
CREATE USER 'xiaomi'@'%' IDENTIFIED BY '12345';
GRANT ALL PRIVILEGES ON *.* TO 'xiaomi'@'%' WITH GRANT OPTION;
FLUSH PRIVILEGES;
```

#### 数据库容器的执行脚本

档案位于 Gaea/util/mocks/containerTest/images/mariadb_testing/mariadb/mysqld_init.sh

```bash
#!/bin/bash

# 一般原本 Docker 的数据库容器无法在 containerd 上执行，进行以下修正
mkdir /var/run/mysqld
# RUN useradd -m mysql
chown mysql:mysql /var/run/mysqld
chmod 777 /var/run/mysqld

# user.sql 为一开始执行 mysqld 服务时，所需要执行的 SQL 脚本，会创建一个用户 xiaomi，并且设置密码
mysqld --init-file=/home/mariadb/user.sql
```

#### 数据库容器的 Dockerfile

档案位于 Gaea/util/mocks/containerdTest/images/mariadb_testing/Dockerfile

```dockerfile
FROM debian:latest

# 安装数据库
RUN apt-get update && apt-get install -y mariadb-server && apt-get clean

# 修改数据库连线设置 (正规表达式为 bind-address(\s*?)=(\s*?)127\.0\.0\.1)
RUN sed -i "s/bind-address.*/bind-address=0.0.0.0/g" /etc/mysql/mariadb.conf.d/50-server.cnf

# 设定用户密码和修正
RUN mkdir -p /home/mariadb/
ADD mariadb /home/mariadb/

# 进行修正
RUN chmod +x /home/mariadb/mysqld_init.sh

# 启动数据库
ENTRYPOINT ["/home/mariadb/mysqld_init.sh"]
```

#### 进行容器镜像档打包

执行以下命令进行打包

```bash
# 安装打包工具 buildah
$ apt-get install -y buildah

# 进入镜像档目录
$ cd Gaea/util/mocks/containerTest/images/mariadb_testing

# 进行打包
$ buildah bud -t mariadb:testing .

# 查询打包结果
$ buildah images
# 会显示以下内容
# REPOSITORY        TAG    IMAGE ID     CREATED        SIZE
# localhost/mariadb latest e4fe0437050b 5 minutes ago  484 MB # 产生新的镜像

# 安装容器工具 podman
$ apt-get install -y podman

# 保存容器镜像为 tar 档 mariadb-latest.tar
$ podman image save localhost/mariadb:testing -o mariadb-testing.tar

# 安装容器工具 podman
$ apt-get install -y podman

# 保存容器镜像为 tar 档 mariadb-latest.tar
$ podman image save localhost/mariadb:testing -o mariadb-testing.tar


```

#### 镜像档打包后保存

执行以下命令进行保存

```bash


# 检查保存结果
$ ls
# 会显示以下内容
# Dockerfile  mariadb  minidb-testing.tar # 新的 tar 档产生
```

#### 上传镜像档到远端仓库

> 目前容器的远端仓库有 [docker hub](https://hub.docker.com/) 和 [qury io](https://quay.io/)，目前打算在 [docker hub](https://hub.docker.com/) 上进行测试，功能较完整的镜象档上传到 [qury io](https://quay.io/)

先把测试用的数据库镜像上传到 [docker hub](https://hub.docker.com/)

```bash
# 到容器目录之下
$ cd gaea/util/mocks/containerdTest/images/mariadb_testing

# 打包容器镜像档
$ buildah bud -t mariadb:testing .

# 上传容器镜象
# <token> 为在 docker 网站中产生的验证 token 
$ skopeo copy docker-archive:./mariadb-testing.tar docker://docker.io/panhongrainbow/mariadb:testing --dest-creds panhongrainbow:<token>
```

### 重新打包 Minideb MariaDB 数据库镜像

> - Minideb MariaDB 镜象目前和其他容器只差在使用的基底镜象不同，之后还要新增功能和精简镜像大小
> - 先只简单记录打包镜像的步骤和命令

#### 启动测试容器

请使用方式去启动 debian 测试容器

```bash
$ ctr namespace create test

# $ ctr -n test i pull "docker.io/library/debian:bullseye-slim"
$ ctr -n test i pull "docker.io/library/ubuntu:latest"

# $ ctr -n test c create docker.io/library/debian:bullseye-slim --with-ns=network:/var/run/netns/gaea-default debian
$ ctr -n test c create docker.io/library/ubuntu:latest --with-ns=network:/var/run/netns/gaea-default debian

$ mkdir /tmp/container

$ ctr -n test snapshot mounts /tmp/container debian | sh

$ ctr -n test task start -d debian

$ ctr -n test task exec -t --exec-id debian-sh debian sh


$ ctr -n test task kill -s SIGKILL debian

$ ctr -n test task rm debian

$ umount /tmp/container

$ ctr -n test c rm debian

```



#### 进行容器镜像档打包

执行以下命令进行打包

```bash
# >>>>> >>>>> >>>>> 根据 DockerfIle 把打包镜像

# 安装打包工具 buildah
$ apt-get install -y buildah

# 进入镜像档目录
$ cd Gaea/util/mocks/containerTest/images/minidb_testing

# 进行打包
$ buildah bud -t minidb:testing .

# 查询打包结果
$ buildah images
# 会显示以下内容
# REPOSITORY         TAG       IMAGE ID       CREATED         SIZE
# localhost/minidb   testing   39ace3b144b1   2 minutes ago   433 MB

# >>>>> >>>>> >>>>> 打包镜像成 tar 档

# 保存容器镜像为 tar 档 minidb-latest.tar
$ podman image save localhost/minidb:testing -o minidb-testing.tar

# 上传容器镜象
# <token> 为在 docker 网站中产生的验证 token 
$ skopeo copy docker-archive:./minidb-testing.tar docker://docker.io/panhongrainbow/minidb:testing --dest-creds panhongrainbow:<token>
```

## 较正式的容器镜像

目前没有















