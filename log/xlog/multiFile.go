package xlog

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// XMultiFileLog is the multi file logger
type XMultiFileLog struct {
	// 这里设定预设写入日志的档名，不在这里作日志 Log 分流的相关设定，避免进行上锁
	defaultXLog string               // 预设要写入的日志档案
	multi       map[string]*XFileLog // 多档案的输出
}

const (
	XMultiFileLogDefaultLogID = "1000000001" // XMultiFileLog 的固定值
)

// NewXMultiFileLog 为赋值函式
func NewXMultiFileLog() XLogger {
	return &XMultiFileLog{}
}

// Init 是用来设置 XLogger
// 设定档假设
// "path": /home/panhong/go/src/github.com/panhong/demo/logs
// "fileName": "result0,result1,result2"
// "level": "Notice,Notice,Notice"
// "service": "shard0,shard1,shard2"
// "skip": "5"
func (ps *XMultiFileLog) Init(config map[string]string) (err error) {
	// 先初始化 ps 的 multi 的对应 map
	// ps.multi = make(map[string]*XFileLog) // 移到后面才能确定初始化 multi map 的大小

	// 有三个设定值使用逗号，分别是 fileName，service 和 level，要特别处理

	// 目录交由储存物件去处理
	/*path, ok := config["path"]
	if !ok {
		err = fmt.Errorf("init XFileLog failed, not found path")
		return
	}*/

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

	// 先初始化 ps 的 multi 的对应 map
	ps.multi = make(map[string]*XFileLog, len(filename)) // 直接指定 multi map 的大小

	// 产生 level 阵列
	var level []string
	lStr, ok := config["level"] // 先确认 level 设定值是否存在
	if ok {                     // 如果 level 值 存在
		level = strings.Split(lStr, ",") // 以逗点分隔开来
	}
	if !ok { // 如果 level 值 不 存在
		err = fmt.Errorf("init XFileLog failed, not found level")
		return
	}

	// 产生 service 阵列
	var service []string
	sStr, ok := config["service"] // 先确认 service 设定值是否存在
	if ok {                       // 如果 service 值 存在
		service = strings.Split(sStr, ",") // 以逗点分隔开来
	}
	if !ok { // 如果 service 值 不 存在
		err = fmt.Errorf("init XFileLog failed, not found service")
		return
	}

	// 进行最后长度检查，因为是以 fileName 阵列为中心跑回圈，所以任何一个阵列的长度要大于 fileName 阵列

	// 检查 service 阵列
	// service 阵列可以只含有一个值，这时所有的档案该设定就会被统一设定
	// service 阵列可以只含多个值，但是数量要大于 fileName 阵列，这时所有的档案该设定就会被个别设定
	if len(service) < len(filename) && len(service) != 1 {
		err = fmt.Errorf("init XFileLog failed, lack service config")
		return
	}

	// fileName 没问题时，就可以直接指定预设的日志档案
	ps.defaultXLog = filename[0] // 第一个日志档案就是预设值

	// 检查 level 阵列
	// level 阵列可以只含有一个值，这时所有的档案该设定就会被统一设定
	// level 阵列可以只含多个值，但是数量要大于 fileName 阵列，这时所有的档案该设定就会被个别设定
	if len(level) < len(filename) && len(level) != 1 {
		err = fmt.Errorf("init XFileLog failed, lack level config")
		return
	}

	// 以 fileName 为中心做回圈
	for i := 0; i < len(filename); i++ {
		p := new(XFileLog)

		// fileName 设定
		config["filename"] = filename[i]

		// level 设定值可以支援单值和多值
		if len(level) >= len(filename) && (len(level[i]) > 0) {
			p.level, _ = strconv.Atoi(level[i]) // level 可以进行个别设定
			config["level"] = service[i]
		}
		if len(level) == 1 && (len(level[0]) > 0) {
			p.level, _ = strconv.Atoi(level[0]) // level 可以进行统一设定
			config["level"] = service[0]
		}

		// skip 设定
		skip, _ := config["skip"]
		if len(skip) > 0 {
			skipNum, err := strconv.Atoi(skip)
			if err == nil {
				p.skip = skipNum
			}
		}

		// runtime 设定
		runtime, ok := config["runtime"]
		if !ok || runtime == "true" || runtime == "TRUE" {
			p.runtime = true
		} else {
			p.runtime = false
		}

		// hostname 设定
		hostname, _ := os.Hostname()
		p.hostname = hostname

		// service 设定值可以支援单值和多值
		if len(service) >= len(filename) && (len(service[i]) > 0) {
			p.service = service[i] // service 可以进行个别设定
			config["service"] = service[i]
		}
		if len(service) == 1 && (len(service[0]) > 0) {
			p.service = service[0] // service 可以进行统一设定
			config["service"] = service[0]
		}

		// 文档和目录，交由储存物件去管理
		// p = filename[i] // fileName 可以进行个别设定
		// config["service"] = service[i]

		p.storage = NewLogStorageClient(config) // 关于储存的设定，比如 目录和文档，交由储存物件去处理

		// 错误回传
		err := p.ReOpen()
		if err != nil { // 一旦有错误，就回传错误
			return err
		}

		// 完整组合 multi map
		ps.multi[filename[i]] = p
	}

	// 正确回传
	return nil
}

// ReOpen 可以进行重新开档
func (ps *XMultiFileLog) ReOpen() error {
	// 多个 xfile 组成 ps.multi，现在把 xfile 一个个打开
	for _, xfile := range ps.multi {
		if err := xfile.storage.client.ReOpen(); err != nil { // 开启日志文档也交由储存物件去管理
			return err // 只要有错误就中止
		}
	}

	// 正确回传
	return nil
}

// SetLevel 可以用来个别设定 multi file 的 日志 log 等级
func (ps *XMultiFileLog) SetLevel(level string) {
	// 产生 level 阵列
	var levelArr []string
	if len(level) > 0 { // 先确认 level 参数不为空
		levelArr = strings.Split(level, ",") // 以逗点分隔开来
	}
	if len(level) == 0 { // 如果 level 值 为空
		return // 立刻中断
	}

	// 检查 levelArr 阵列
	// levelArr 阵列可以只含有一个值，这时所有的档案该设定就会被统一设定
	// levelArr 阵列可以只含多个值，但是数量要大于 fileName 阵列，这时所有的档案该设定就会被个别设定
	if len(levelArr) < len(ps.multi) && len(levelArr) != 1 {
		return // 立刻中断
	}

	// 以 fileName 为中心做回圈
	for i := 0; i < len(ps.multi); i++ {
		for key := range ps.multi {
			ps.multi[key].level = LevelFromStr(levelArr[i])
		}
	}
}

// SetSkip 设定忽略重复的行数
func (ps *XMultiFileLog) SetSkip(skip int) {
	// 以 ps.multi 为中心做回圈
	for key := range ps.multi {
		ps.multi[key].skip = skip // 进行统一设定 skip
	}
}

// prepareMultiFile 为和 XMultiFileLog 的预先处理函式
// 有想到最壤的状况，如果 format 的参数会使 prepareMultiFile 发生错误，最后输出会是 file 为预设日志档案，newFormat 为空字串
func (ps *XMultiFileLog) prepareMultiFile(format string) (string, string) {
	// 先拆分 format 字串
	logFile, newFormat := FileFormatFromStr(format) // logFile 值为档名，newFormat 为字串格式

	// 检查拆分的结果
	_, ok := ps.multi[logFile] // 检查 logFile 值是否存在
	if logFile == "" || ok == false {
		logFile = ps.defaultXLog // 如果 logFile 值为不存在就采预设值
	}

	// 正确回传
	return logFile, newFormat
}

// >>>>> >>>>> >>>>> >>>>> >>>>> 警告日志

// Warn 显示 Warn 的资讯，格式为 档名::日志格式为第一个参数
func (ps *XMultiFileLog) Warn(format string, a ...interface{}) error {
	// 先拆分 format 字串
	logFile, newFormat := ps.prepareMultiFile(format)

	// 以下程式码尽量保留
	if ps.multi[logFile].level > WarnLevel {
		return nil
	}

	// 传入新的 newFormat 参数
	return ps.warnx(XMultiFileLogDefaultLogID, newFormat, a...)
}

// Warnx 显示 Warnx 的资讯，格式为 档名::日志格式为第一个参数
func (ps *XMultiFileLog) Warnx(logID, format string, a ...interface{}) error {
	// 先拆分 format 字串
	logFile, newFormat := ps.prepareMultiFile(format)

	// 以下程式码尽量保留
	if ps.multi[logFile].level > WarnLevel {
		return nil
	}

	// 传入新的 newFormat 参数
	return ps.warnx(logID, newFormat, a...)
}

// warnx 为警告日志的写入函式
func (ps *XMultiFileLog) warnx(logID, format string, a ...interface{}) error {
	// 先拆分 format 字串
	logFile, newFormat := ps.prepareMultiFile(format)

	// 以下程式码尽量保留
	logText := formatValue(newFormat, a...) // 传入新的 newFormat 参数
	fun, filename, lineno := getRuntimeInfo(ps.multi[logFile].skip)
	logText = formatLineInfo(ps.multi[logFile].runtime, fun, filepath.Base(filename), logText, lineno)
	//logText = fmt.Sprintf("[%s:%s:%d] %s", fun, filepath.Base(fileName), lineno, logText)

	return ps.multi[logFile].write(WarnLevel, &logText, logID)
}

// >>>>> >>>>> >>>>> >>>>> >>>>> 严重错误日志

// Fatal 显示 Fatal 的资讯，格式为 档名::日志格式为第一个参数
func (ps *XMultiFileLog) Fatal(format string, a ...interface{}) error {
	// 先拆分 format 字串
	logFile, newFormat := ps.prepareMultiFile(format)

	// 以下程式码尽量保留
	if ps.multi[logFile].level > FatalLevel {
		return nil
	}

	// 传入新的 newFormat 参数
	return ps.fatalx(XMultiFileLogDefaultLogID, newFormat, a...)
}

// Fatalx 显示 Fatalx 的资讯，格式为 档名::日志格式为第一个参数
func (ps *XMultiFileLog) Fatalx(logID, format string, a ...interface{}) error {
	// 先拆分 format 字串
	logFile, newFormat := ps.prepareMultiFile(format)

	// 以下程式码尽量保留
	if ps.multi[logFile].level > FatalLevel {
		return nil
	}

	// 传入新的 newFormat 参数
	return ps.fatalx(logID, newFormat, a...)
}

// fatalx 为严重错误日志的写入函式
func (ps *XMultiFileLog) fatalx(logID, format string, a ...interface{}) error {
	// 先拆分 format 字串
	logFile, newFormat := ps.prepareMultiFile(format)

	// 以下程式码尽量保留
	logText := formatValue(newFormat, a...) // 传入新的 newFormat 参数
	fun, filename, lineno := getRuntimeInfo(ps.multi[logFile].skip)
	logText = formatLineInfo(ps.multi[logFile].runtime, fun, filepath.Base(filename), logText, lineno)
	//logText = fmt.Sprintf("[%s:%s:%d] %s", fun, filepath.Base(fileName), lineno, logText)

	return ps.multi[logFile].write(FatalLevel, &logText, logID)
}

// >>>>> >>>>> >>>>> >>>>> >>>>> 警告日志

// Notice 显示 Notice 的资讯，格式为 档名::日志格式为第一个参数
func (ps *XMultiFileLog) Notice(format string, a ...interface{}) error {
	// 先拆分 format 字串
	logFile, newFormat := ps.prepareMultiFile(format)

	// 以下程式码尽量保留
	if ps.multi[logFile].level > NoticeLevel {
		return nil
	}

	// 传入新的 newFormat 参数
	return ps.multi[logFile].noticex(XMultiFileLogDefaultLogID, newFormat, a...)
}

// Noticex 显示 Noticex 的资讯，格式为 档名::日志格式为第一个参数
func (ps *XMultiFileLog) Noticex(logID, format string, a ...interface{}) error {
	// 先拆分 format 字串
	logFile, newFormat := ps.prepareMultiFile(format)

	// 以下程式码尽量保留
	if ps.multi[logFile].level > NoticeLevel {
		return nil
	}

	// 传入新的 newFormat 参数
	return ps.noticex(logID, newFormat, a...)
}

// noticex 为注意日志的写入函式
func (ps *XMultiFileLog) noticex(logID, format string, a ...interface{}) error {
	// 先拆分 format 字串
	logFile, newFormat := ps.prepareMultiFile(format)

	// 以下程式码尽量保留
	logText := formatValue(newFormat, a...) // 传入新的 newFormat 参数
	fun, filename, lineno := getRuntimeInfo(ps.multi[logFile].skip)
	logText = formatLineInfo(ps.multi[logFile].runtime, fun, filepath.Base(filename), logText, lineno)

	// 传入新的 newFormat 参数
	return ps.multi[logFile].write(NoticeLevel, &logText, logID)
}

// >>>>> >>>>> >>>>> >>>>> >>>>> 追踪日志

// Trace 显示 Trace 的资讯，格式为 档名::日志格式为第一个参数
func (ps *XMultiFileLog) Trace(format string, a ...interface{}) error {
	// 先拆分 format 字串
	logFile, newFormat := ps.prepareMultiFile(format)

	// 以下程式码尽量保留
	return ps.multi[logFile].tracex(XMultiFileLogDefaultLogID, newFormat, a...) // 传入新的 newFormat 参数s
}

// Tracex 显示 Tracex 的资讯，格式为 档名::日志格式为第一个参数
func (ps *XMultiFileLog) Tracex(logID, format string, a ...interface{}) error {
	// 先拆分 format 字串
	logFile, newFormat := ps.prepareMultiFile(format)

	// 以下程式码尽量保留
	return ps.multi[logFile].tracex(logID, newFormat, a...) // 传入新的 newFormat 参数
}

// tracex 为追踪日志的写入函式
func (ps *XMultiFileLog) tracex(logID, format string, a ...interface{}) error {
	// 先拆分 format 字串
	logFile, newFormat := ps.prepareMultiFile(format)

	// 以下程式码尽量保留
	if ps.multi[logFile].level > TraceLevel {
		return nil
	}

	logText := formatValue(newFormat, a...) // 传入新的 newFormat 参数
	fun, filename, lineno := getRuntimeInfo(ps.multi[logFile].skip)
	logText = formatLineInfo(ps.multi[logFile].runtime, fun, filepath.Base(filename), logText, lineno)
	//logText = fmt.Sprintf("[%s:%s:%d] %s", fun, filepath.Base(fileName), lineno, logText)

	return ps.multi[logFile].write(TraceLevel, &logText, logID)
}

// >>>>> >>>>> >>>>> >>>>> >>>>> 除错日志

// Debug 显示 Debug 的资讯，格式为 档名::日志格式为第一个参数
func (ps *XMultiFileLog) Debug(format string, a ...interface{}) error {
	// 先拆分 format 字串
	logFile, newFormat := ps.prepareMultiFile(format)

	// 以下程式码尽量保留
	return ps.multi[logFile].debugx(XMultiFileLogDefaultLogID, newFormat, a...) // 传入新的 newFormat 参数
}

// Debugx 显示 Debugx 的资讯，格式为 档名::日志格式为第一个参数
func (ps *XMultiFileLog) Debugx(logID, format string, a ...interface{}) error {
	// 先拆分 format 字串
	logFile, newFormat := ps.prepareMultiFile(format)

	// 以下程式码尽量保留
	return ps.multi[logFile].debugx(logID, newFormat, a...) // 传入新的 newFormat 参数
}

// debugx 为除错日志的写入函式
func (ps *XMultiFileLog) debugx(logID, format string, a ...interface{}) error {
	// 先拆分 format 字串
	logFile, newFormat := ps.prepareMultiFile(format)

	// 以下程式码尽量保留
	if ps.multi[logFile].level > DebugLevel {
		return nil
	}

	logText := formatValue(newFormat, a...) // 传入新的 newFormat 参数
	fun, filename, lineno := getRuntimeInfo(ps.multi[logFile].skip)
	logText = formatLineInfo(ps.multi[logFile].runtime, fun, filepath.Base(filename), logText, lineno)

	return ps.multi[logFile].write(DebugLevel, &logText, logID)
}

// Close  可以关闭所有档案
func (ps *XMultiFileLog) Close() {
	// 多个 xfile 组成 ps.multi，现在把 xfile 一个个关闭
	for _, xfile := range ps.multi {
		_ = xfile.storage.client.Close() // 关闭文档交由储存物件去处理
	}
}
