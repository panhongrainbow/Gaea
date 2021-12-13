# 日誌分流說明 multiFile

> 1. 原本小米的程式碼是支援把日誌寫入到 console 和 檔案裡，現在在基於其基礎，把日誌進行分流，程式碼有不同的區塊，每個區塊會寫入不同的檔案
> 2. 小米的設定檔中，在 etc 目錄下，有 namespace 資料夾，如果能做到每個 namespace 的日誌都分別集中到一個日誌檔，分別有自己的日誌檔，是很實用的功能

## 1 程式碼修改過程

> 整個程式碼在編寫的過程中，會儘量把檔案 multiFile.go 裡面的程式碼儘量依附在 file.go 裡，這樣子可以避免重複的程式碼

### 1-1 設定檔的不同

這裡是要比較 把日誌寫入檔案 (在設定值為 log_output=file) 和 把日誌寫入檔案並分流  (在設定值為 log_output=multiFile)，在設定檔細節的不同

#### 原設定值

在把 日誌寫入檔案 (在設定值為 log_output=file)，其設定檔如下 (在這次新增和修改程式碼後，持續支援此設定值)

```ini
;log config
log_path=./logs
log_level=Notice
log_filename=gaea
log_output=file
```

在把 日誌寫入檔案並分流 (在設定值為 log_output=multiFile)，其設定檔如下

#### 新設定值 (版本一)

一定要有兩個以上的設定檔，才能做日誌的分流，

如果只有一個檔案，那就不如使用 設定值為 log_output=file

反正如果要做日誌分流，至少要設定兩個以上的日誌檔案

```ini
;log config
log_path=./logs
log_level=Notice
log_filename=default,log1
log_output=multiFile
```

在這個設定檔，會自産生兩個日誌檔 logs/default.log 和 logs/log1.log，

因為等級只有一個值為 Notice，所有的檔案都採用此值作為預設值，都為為 Notice

#### 新設定值 (版本二)

一樣跟 版本一 一樣，要有兩個以上的設定檔，才能做日誌的分流，

```ini
;log config
log_path=./logs
log_level=Notice,debug
log_filename=default,log1
log_output=multiFile
```

這時會産生兩個日誌檔 logs/default.log 和 logs/log1.log，

但是這次不同，使用兩個等級，一個為 Notice，另一個為 Debug

相對應，logs/default.log 就會使用等級 Notice，logs/log1.log 就會使用等級 Debug



