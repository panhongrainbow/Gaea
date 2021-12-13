# 日志分流说明 multiFile

> 1. 原本小米的程式码是支援把日志写入到 console 和 档案里，现在在基于其基础，把日志进行分流，程式码有不同的区块，每个区块会写入不同的档案
> 2. 小米的设定档中，在 etc 目录下，有 namespace 资料夹，如果能做到每个 namespace 的日志都分别集中到一个日志档，分别有自己的日志档，是很实用的功能

## 1 程式码修改过程

> 整个程式码在编写的过程中，会尽量把档案 multiFile.go 里面的程式码尽量依附在 file.go 里，这样子可以避免重复的程式码

### 1-1 设定档的不同

这里是要比较 把日志写入档案 (在设定值为 log_output=file) 和 把日志写入档案并分流  (在设定值为 log_output=multiFile)，在设定档细节的不同

#### 原设定值

在把 日志写入档案 (在设定值为 log_output=file)，其设定档如下 (在这次新增和修改程式码后，持续支援此设定值)

```ini
;log config
log_path=./logs
log_level=Notice
log_filename=gaea
log_output=file
```

在把 日志写入档案并分流 (在设定值为 log_output=multiFile)，其设定档如下

#### 新设定值 (版本一)

一定要有两个以上的设定档，才能做日志的分流，

如果只有一个档案，那就不如使用 设定值为 log_output=file

反正如果要做日志分流，至少要设定两个以上的日志档案

```ini
;log config
log_path=./logs
log_level=Notice
log_filename=default,log1
log_output=multiFile
```

在这个设定档，会自产生两个日志档 logs/default.log 和 logs/log1.log，

因为等级只有一个值为 Notice，所有的档案都采用此值作为预设值，都为为 Notice

#### 新设定值 (版本二)

一样跟 版本一 一样，要有两个以上的设定档，才能做日志的分流，

```ini
;log config
log_path=./logs
log_level=Notice,debug
log_filename=default,log1
log_output=multiFile
```

这时会产生两个日志档 logs/default.log 和 logs/log1.log，

但是这次不同，使用两个等级，一个为 Notice，另一个为 Debug

相对应，logs/default.log 就会使用等级 Notice，logs/log1.log 就会使用等级 Debug



