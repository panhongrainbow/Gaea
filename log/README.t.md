# 日誌分流說明 multiFile

> 1. 原本小米的程式碼是支援把日誌寫入到 console 和 檔案裡，現在在基於其基礎，把日誌進行分流，程式碼有不同的區塊，每個區塊會寫入不同的檔案
> 2. 小米的設定檔中，在 etc 目錄下，有 namespace 資料夾，如果能做到每個 namespace 的日誌都分別集中到一個日誌檔，分別有自己的日誌檔，是很實用的功能

## 1 日誌分流程式

> 整個程式碼在編寫的過程中，會儘量把檔案 multiFile.go 裡面的程式碼儘量依附在 file.go 裡，這樣子可以避免重複的程式碼

### 1-1 設定檔的設定方式

這裡是要比較 把日誌寫入檔案 (在設定值為 log_output=file) 和 把日誌寫入檔案並分流  (在設定值為 log_output=multiFile)，在設定檔細節的不同

#### 1-1-1 temp

//

#### 1-1-2 原設定值

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

#### 1-1-3 新設定值 (統一設定)

一定要有兩個以上的設定檔，才能做日誌的分流，

如果只有一個檔案，那就不如使用 設定值為 log_output=file

反正如果要做日誌分流，至少要設定兩個以上的日誌檔案

在 /etc/gaea.ini 設定檔，設定值如下

```ini
;log config
log_path=./logs
log_level=Notice
log_filename=default,log1
service_name=svc1
log_output=multiFile
```

- 在這個設定檔，會自産生兩個日誌檔 logs/default.log 和 logs/log1.log，

- 因為等級只有一個值為 Notice，所有的檔案都採用此值作為預設值，都為為 Notice
- 服務名稱為統一設定，logs/default.log 和 logs/log1.log 的服務名稱都為 svc1

#### 1-1-4 新設定值 (各別設定)

一樣跟 版本一 一樣，要有兩個以上的設定檔，才能做日誌的分流

在 /etc/gaea.ini 設定檔，設定值如下

```ini
;log config
log_path=./logs
log_level=Notice,debug
log_filename=default,log1
service_name=svc1,svc2
log_output=multiFile
```

- 這時會産生兩個日誌檔 logs/default.log 和 logs/log1.log，

- 但是這次不同，使用兩個等級，一個為 Notice，另一個為 Debug
  相對應，logs/default.log 就會使用等級 Notice，logs/log1.log 就會使用等級 Debug
- 服務名稱為各別設定，logs/default.log 的服務名稱都為 svc1，logs/log1.log 的服務名稱都為 svc2

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

#### 1-2-1 日誌編號的問題

不管是 (1) 終端機 console 、(2) 檔案 file 還是 (3) 多檔案輸出 multiFile，都會有日誌編號

比如在

- (1) 終端機 console 顯示的日誌編號為 800000001
- (2) 檔案 file 顯示的日誌編號為 900000001

這次所修改的程式，一開始就指定多檔輸出的日誌編號為 1000000001

- (3) 檔案 multiFile 顯示的日誌編號為 1000000001

雖然 多檔輸出 multiFile 的程式碼是依附在 檔案 file，但是日誌編號還是依然可以正常指定，因為原本程式就有預留可以把傳入 日誌編號 參數的地方

```go
// Notice 显示 Notice 的资讯，格式为 档名::日志格式为第一个参数
func (ps *XMultiFileLog) Notice(format string, a ...interface{}) error {
    // >>>>> >>>>> >>>>> >>>>> >>>>> 預先處理的部份
    
	// 先拆分 format 字串
	logFile, newFormat := ps.prepareMultiFile(format)
    
    // >>>>> >>>>> >>>>> >>>>> >>>>> 保留原程式的部份

	// 以下程式码尽量保留
	if ps.multi[logFile].level > NoticeLevel {
		return nil
	}

    // 传入新的 newFormat 参数 (最後在這裡指定日誌編號)
	return ps.multi[logFile].noticex(XMultiFileLogDefaultLogID, newFormat, a...)
}
```

以上程式碼進行以下拆解

- 預先處理的部份
  在這個部份會把 format 字串，拆分成

  ​     1 日誌檔案變數 logFile
  ​     2 新的 format 格式字串 newFormat

  之後做後續處理

- 保留原程式的部份
  在這個部份整個邏輯會和 file.go 的程式碼很像
  就在最後 noticex 函式裡，會傳入日誌編號的參數 XMultiFileLogDefaultLogID 值為 1000000001

#### 1-2-2 多檔輸出的物件設計

物件的程式碼如下

```go
// XMultiFileLog is the multi file logger
type XMultiFileLog struct {
	// 这里设定预设写入日志的档名，不在这里作日志 Log 分流的相关设定，避免进行上锁
	defaultXLog string               // 预设要写入的日志档案
	multi       map[string]*XFileLog // 多档案的输出
}
```

- 有一個 defaultXLog 會記錄預設要寫入的預設日誌檔案，當程式無法得知要寫入那個日誌檔案時，就會直接把日誌寫入日誌檔案裡

- multi 元素把 檔案輸出 file 物件 和其檔案名稱做對應，

  也就是 XMultiFileLog (多檔日誌輸出) 是由多個 XFileLog (單檔日誌輸出) 所組成的

#### 1-2-3 多檔輸出的物件初始化

主程式 main 會配有一個 副程式 initXLog，副程式 initXLog 的作用如下：

- 整理和組成 init 函式的 config 參數
- 把現在的 日誌設定值 設置 成全域變數

#### 1-2-4 initXLog 程式碼如下

- 將會組成 cfg 變數後，依照設定值會傳入任一個 console、file 和 multiFile 的 init 的方法
- 所以知道單元測試的剖面會切在 副程式 initXLog 的後面
- 其實一直會有想把 cfg["filename"] = filename 這一行中的 filename 字串改成駝峰式命名 fileName，也就是
   cfg["filename"] = filename 改成 cfg["fileName"] = fileName，但後面看起來好像不協調，因為整個 initXLog 函式裡將會只有 filename 很明顯使用駝峰式命名 fileName，就先算了 

```go
func initXLog(output, path, filename, level, service string) error {
	cfg := make(map[string]string)
	cfg["path"] = path
	cfg["filename"] = filename
	cfg["level"] = level
	cfg["service"] = service
	cfg["skip"] = "5"
    // 设置xlog打印方法堆栈需要跳过的层数, 5目前为调用log.Debug()等方法的方法名, 比xlog默认值多一层.

	logger, err := xlog.CreateLogManager(output, cfg)
	if err != nil {
		return err
	}

	log.SetGlobalLogger(logger) // 設定成全域日誌變數
	return nil
}
```

#### 1-2-5 副程式 init 程式碼如下

> 多檔輸出 multiFuile 的副程式 init 程式碼如下

```go
func (ps *XMultiFileLog) Init(config map[string]string) (err error) {
	// 先初始化 ps 的 multi 的对应 map
	ps.multi = make(map[string]*XFileLog)

	// 有三个设定值使用逗号，分别是 fileName，service 和 level，要特别处理

	// 产生 fileName 阵列
	var filename []string
	fStr, ok := config["filename"] // 先确认 fileName 设定值是否存在
	if ok {                        // 如果 fileName 值 存在
		filename = strings.Split(fStr, ",") // 以逗点分隔开来
	}
	if !ok { // 如果 fileName 值 不 存在
		err = fmt.Errorf("init XFileLog failed, not found filename")
		return
	}
    
    // 以下程式碼略過
}
```

## 2 單元測試的運作方式
