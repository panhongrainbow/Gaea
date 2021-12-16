# 日誌分流說明 multiFile

> 1. 原本小米的程式碼是支援把日誌寫入到 console 和 檔案裡，現在在基於其基礎，把日誌進行分流，程式碼有不同的區塊，每個區塊會寫入不同的檔案
> 2. 小米的設定檔中，在 etc 目錄下，有 namespace 資料夾，如果能做到每個 namespace 的日誌都分別集中到一個日誌檔，分別有自己的日誌檔，是很實用的功能

## 1 日誌分流程式

> 整個程式碼在編寫的過程中，會儘量把檔案 multiFile.go 裡面的程式碼儘量依附在 file.go 裡，這樣子可以避免重複的程式碼

### 1-1 設定檔的設定方式

這裡是要比較 把日誌寫入檔案 (在設定值為 log_output=file) 和 把日誌寫入檔案並分流  (在設定值為 log_output=multiFile)，在設定檔細節的不同

#### 原設定值

在把 日誌寫入檔案 (在設定值為 log_output=file)，其設定檔如下 (在這次新增和修改程式碼後，持續支援此設定值)

在 /etc/gaea.ini 設定檔，設定值如下

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

在 /etc/gaea.ini 設定檔，設定值如下

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

一樣跟 版本一 一樣，要有兩個以上的設定檔，才能做日誌的分流

在 /etc/gaea.ini 設定檔，設定值如下

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

#### 預設檔案

預設日誌檔案為設定值 log_filename=default,log1 的第一個元素，也就是 default.log

當日誌找不到指定要寫入的日誌檔案時，會直接寫入預設的日誌檔案，也就是default.log

#### 指定日誌檔案的方式

在程式碼中，可以指定要寫入的日誌檔案

指定方式如下

1. 指定格式如下，在 format 參數裡，以雙分號 :: 為界線
2. 雙分號 :: 以前的為 指定的日誌檔案
3. 雙分號 :: 以後的為 指定的日誌格式，%s 也加在這裡

設定例子如下

```go
// 在 /etc/gaea.ini 設定檔，設定值 指定
// log_filename=default,log1
// log_output=multiFile

err = ps.Notice("record1") // 在沒有指定寫入的日誌檔時，會先寫到預設的日誌檔裡，也就是 default.log
err = ps.Notice("log1::record2") // 有指定寫入日誌檔 log1.log，所以會把日誌寫到檔案 log1.log 裡
err = ps.Notice("log2::record3") // 有指定寫入日誌檔 log2.log，但是 log2.log 並不存在於 gaea.ini 的設定裡，所以只能把日誌寫到預設的日誌檔裡，也就是 default.log
```

### 1-2 程式的設計邏輯

#### 日誌編號的問題

不管是 (1) 終端機 console 、(2) 檔案 file 還是 (3) 多檔案輸出 multiFile，都會有日誌編號

比如在

- (1) 終端機 console 顯示的日誌編號為 800000001
- (2) 檔案 file 顯示的日誌編號為 900000001

這次所修改的程式，一開始就指定多檔輸出的日誌編號為 1000000001

- (3) 檔案 multiFile 顯示的日誌編號為 1000000001

雖然 多檔輸出 multiFile 的程式碼是依附在 檔案 file，但是日誌編號還是依然可以正常指定，因為原本程式就有預留可以把傳入 日誌編號 參數的地方

```go
```









## 2 單元測試的運作方式
