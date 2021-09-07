# 中文文件繁简的互换方法(简体版) 

> 因为有些人是读简体中文，有些人是读繁体中文，这篇文件说明如何把所有程式码进行繁简互换

## 安装相关的套件

安装相关的中文繁简互换套件

```bash
# 安装 OpenCC 套件
$ apt-get install -y opencc libopencc2-data
```

## 繁体中文转成简体中文

执行以下指令，把繁体中文转成简体中文

```bash
# 进入专案文件资料夹
$ cd Gaea/doc

# 备份文件
$ mv goland-develop.md goland-develop.old.md

# 把文件由繁体中文转成简体中文
$ opencc -i goland-develop.t.md -o goland-develop.md -c t2s.json
```

## 简体中文转成繁体中文

执行以下指令，把简体中文转成繁体中文

```bash
# 进入专案文件资料夹
$ cd Gaea/doc

# 备份文件
$ mv goland-develop.md goland-develop.old.md

# 把文件由简体中文转成繁体中文
$ opencc -i goland-develop.md -o goland-develop.t.md -c s2t.json
```
