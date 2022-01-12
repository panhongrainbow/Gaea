# Gaea 小米数据库中间件 Ectd V3 API 升级过程说明

> 源由 Etcd V3 API 使用 gRPC 作为沟通的协定，效能上会有所提升，有必要进行升级

## 1 Etcd 测试环境设置 

> 直接到 Etcd V3 https://github.com/etcd-io/etcd/releases 找到容器 etcd v3.5.1 的启动方式，以下为启动指令

```bash
# 连结到 Etcd GitHub
$ firefox https://github.com/etcd-io/etcd/releases

# 删除之前的容器
$ docker stop etcd-gcr-v3.5.1
$ docker rm etcd-gcr-v3.5.1
$ sudo rm -rf /tmp/etcd-data.tmp

# 执行以下指令去执行 etcd v3.5.1 容器
$ rm -rf /tmp/etcd-data.tmp && mkdir -p /tmp/etcd-data.tmp && \
  docker rmi gcr.io/etcd-development/etcd:v3.5.1 || true && \
  docker run \
  -p 2379:2379 \
  -p 2380:2380 \
  --mount type=bind,source=/tmp/etcd-data.tmp,destination=/etcd-data \
  --name etcd-gcr-v3.5.1 \
  gcr.io/etcd-development/etcd:v3.5.1 \
  /usr/local/bin/etcd \
  --name s1 \
  --data-dir /etcd-data \
  --listen-client-urls http://0.0.0.0:2379 \
  --advertise-client-urls http://0.0.0.0:2379 \
  --listen-peer-urls http://0.0.0.0:2380 \
  --initial-advertise-peer-urls http://0.0.0.0:2380 \
  --initial-cluster s1=http://0.0.0.0:2380 \
  --initial-cluster-token tkn \
  --initial-cluster-state new \
  --log-level info \
  --logger zap \
  --log-outputs stderr

# 执行以下指令进行初期测试
$ docker exec etcd-gcr-v3.5.1 /bin/sh -c "/usr/local/bin/etcd --version"
$ docker exec etcd-gcr-v3.5.1 /bin/sh -c "/usr/local/bin/etcdctl version"
$ docker exec etcd-gcr-v3.5.1 /bin/sh -c "/usr/local/bin/etcdctl endpoint health"
$ docker exec etcd-gcr-v3.5.1 /bin/sh -c "/usr/local/bin/etcdctl put foo bar"
$ docker exec etcd-gcr-v3.5.1 /bin/sh -c "/usr/local/bin/etcdctl get foo"
$ docker exec etcd-gcr-v3.5.1 /bin/sh -c "/usr/local/bin/etcdutl version"
```

这时发现，用原始程式码会发生以下错误

```bash
# response is invalid json. The endpoint is probably not valid etcd cluster endpoint.
```

## 2 Etcd V3 GUI 介面

可以安装 Etcd 图形介面去观察 Etcd 的写入状况

```bash
# 安装 Etcd 图形介面工具
$ snap install etcd-manager
```

## 3 Etcd 的连线测试

Etcd V3 只要接通后，就会一直保持连线，所以不用另外在写 Ping 函式进行侦测

