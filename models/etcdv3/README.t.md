# Gaea 小米數據庫中間件 Ectd V3 API 升級過程說明

> 源由 Etcd V3 API 使用 gRPC 作為溝通的協定，效能上會有所提升，有必要進行升級

## 1 Etcd 測試環境設置 

> 直接到 Etcd V3 https://github.com/etcd-io/etcd/releases 找到容器 etcd v3.5.1 的啟動方式，以下為啟動指令

```bash
# 連結到 Etcd GitHub
$ firefox https://github.com/etcd-io/etcd/releases

# 刪除之前的容器
$ docker stop etcd-gcr-v3.5.1
$ docker rm etcd-gcr-v3.5.1
$ sudo rm -rf /tmp/etcd-data.tmp

# 執行以下指令去執行 etcd v3.5.1 容器
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

# 執行以下指令進行初期測試
$ docker exec etcd-gcr-v3.5.1 /bin/sh -c "/usr/local/bin/etcd --version"
$ docker exec etcd-gcr-v3.5.1 /bin/sh -c "/usr/local/bin/etcdctl version"
$ docker exec etcd-gcr-v3.5.1 /bin/sh -c "/usr/local/bin/etcdctl endpoint health"
$ docker exec etcd-gcr-v3.5.1 /bin/sh -c "/usr/local/bin/etcdctl put foo bar"
$ docker exec etcd-gcr-v3.5.1 /bin/sh -c "/usr/local/bin/etcdctl get foo"
$ docker exec etcd-gcr-v3.5.1 /bin/sh -c "/usr/local/bin/etcdutl version"
```

這時發現，用原始程式碼會發生以下錯誤

```bash
# response is invalid json. The endpoint is probably not valid etcd cluster endpoint.
```

## 2 Etcd V3 GUI 介面

可以安裝 Etcd 圖形介面去觀察 Etcd 的寫入狀況

```bash
# 安裝 Etcd 圖形介面工具
$ snap install etcd-manager
```

## 3 Etcd 的連線測試

Etcd V3 只要接通後，就會一直保持連線，所以不用另外在寫 Ping 函式進行偵測

