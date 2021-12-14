# 日志分流说明 multiFile

> 1. 原本小米的程式码是支援把日志写入到 console 和 档案里，现在在基于其基础，把日志进行分流，程式码有不同的区块，每个区块会写入不同的档案
> 2. 小米的设定档中，在 etc 目录下，有 namespace 资料夹，如果能做到每个 namespace 的日志都分别集中到一个日志档，分别有自己的日志档，是很实用的功能

## 1 日志分流程式

> 整个程式码在编写的过程中，会尽量把档案 multiFile.go 里面的程式码尽量依附在 file.go 里，这样子可以避免重复的程式码

### 1-1 设定档的设定方式

这里是要比较 把日志写入档案 (在设定值为 log_output=file) 和 把日志写入档案并分流  (在设定值为 log_output=multiFile)，在设定档细节的不同

#### 原设定值

在把 日志写入档案 (在设定值为 log_output=file)，其设定档如下 (在这次新增和修改程式码后，持续支援此设定值)

在 /etc/gaea.ini 设定档，设定值如下

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

在 /etc/gaea.ini 设定档，设定值如下

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

一样跟 版本一 一样，要有两个以上的设定档，才能做日志的分流

在 /etc/gaea.ini 设定档，设定值如下

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

#### 预设档案

预设日志档案为设定值 log_filename=default,log1 的第一个元素，也就是 default.log

当日志找不到指定要写入的日志档案时，会直接写入预设的日志档案，也就是default.log

#### 指定日志档案的方式

在程式码中，可以指定要写入的日志档案

指定方式如下

1. 指定格式如下，在 format 参数里，以双分号 :: 为界线
2. 双分号 :: 以前的为 指定的日志档案
3. 双分号 :: 以后的为 指定的日志格式，%s 也加在这里

设定例子如下

```go
// 在 /etc/gaea.ini 设定档，设定值 指定
// log_filename=default,log1
// log_output=multiFile

err = ps.Notice("record1") // 在没有指定写入的日志档时，会先写到预设的日志档里，也就是 default.log
err = ps.Notice("log1::record2") // 有指定写入日志档 log1.log，所以会把日志写到档案 log1.log 里
err = ps.Notice("log2::record3") // 有指定写入日志档 log2.log，但是 log2.log 并不存在于 gaea.ini 的设定里，所以只能把日志写到预设的日志档里，也就是 default.log
```

### 1-2 程式的设计逻辑

#### 日志编号的问题

不管是 终端机 console 、档案 file 还是多档输出 multiFile 



## 2 单元测试的运作方式
