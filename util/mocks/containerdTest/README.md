# 容器单元测试

## 说明

### 会面临的问题如下

- 有一个大问题，**数据库测试数据**要如何载入到容器内，目前偏向的作法如下
  制作一个本地的印象档，里面会有一个脚本会在启动时，自动载入测试数据到容器内，
  好处是可以减少依赖，不需要用 orm 去载入测试数据，能少用一个依赖就少用一个依赖，为了是要增加单元测试的稳定性
- docker 未来的重用性将会下降，未来将会直接操作 **containerd** ，不经由 docker 去控制，这次测试准备把 containerd 整合到单元测试内，测试包命名为 **containerdTest**

## Containerd 容器的安装

### 简介

此为 containerd 的早期版本，已经去除 docker 那一层去进行效能的提升<img src="./assets/image-20220331141145598.png" alt="image-20220331141145598" style="zoom:100%;" /> 

### containerd 主体安装

> 参考以下网址进行安装 [containerd安装文件](https://containerd.io/downloads/)

```bash
# 先安装依赖，libseccomp2 套件将会指定 process 去呼叫特定的 system call 
$ sudo apt-get update
$ sudo apt-get install libseccomp2

# 下载 containerd 套件，目前最新版本为 1.6.2 版，
$ wget https://github.com/containerd/containerd/releases/download/v1.6.2/cri-containerd-cni-1.6.2-linux-amd64.tar.gz

# 这份压缩包含 containerd 运行时所需要的 runc 
$ tar -tf cri-containerd-cni-1.6.2-linux-amd64.tar.gz | grep runc
usr/local/bin/containerd-shim-runc-v2
usr/local/bin/containerd-shim-runc-v1
usr/local/sbin/runc

# 进行安装整个 containerd
$ sudo tar -C / -xzf cri-containerd-cni-1.6.2-linux-amd64.tar.gz

# 檢查 ctr 指令是否存在
$ tar -tf cri-containerd-cni-1.6.2-linux-amd64.tar.gz | grep ctr
usr/local/bin/ctr
```

### containerd 插件設定

```bash
# 產生 config.toml 設定
$ tar -tf cri-containerd-cni-1.6.2-linux-amd64.tar.gz | grep config.toml # 压缩包未含设定档
$ containerd config default > /etc/containerd/config.toml # 用命令产生预设文档





```







### cni 网路设定



## ContainerdTest 单元测试



### 运作说明



### 操作演示