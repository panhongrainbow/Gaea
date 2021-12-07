package xlog

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// XMultiFileLog is the multi file logger
type XMultiFileLog struct {
	multi map[string]*XFileLog // 多档案的输出
}

const (
	XMultiFileLogDefaultLogID = "1000000001" // XMultiFileLog 的固定值
)

// NewXMultiFileLog 为赋值函式
func NewXMultiFileLog() XLogger {
	return &XFileLog{
		skip: XLogDefSkipNum,
	}
}

// Init 是用来设置 XLogger
// 设定档假设
// "path": /home/panhong/"
// "filename": "result0,result1,result2"
// "level": "Notice,Notice,Notice"
// "service": "shard0,shard1,shard2"
// "skip": "5"
func (ps *XMultiFileLog) Init(config map[string]string) (err error) {
	// 先初始化 ps 的 multi 的对应 map
	ps.multi = make(map[string]*XFileLog)

	// 有三个设定值使用逗号，分别是 filename，service 和 level，要特别处理

	// 产生 filename 阵列
	var filename []string
	fStr, ok := config["filename"] // 先确认 filename 设定值是否存在
	if ok { // 如果 filename 值 存在
		filename = strings.Split(fStr, ",") // 以逗点分隔开来
	}
	if !ok { // 如果 filename 值 不 存在
		err = fmt.Errorf("init XFileLog failed, not found filename")
		return
	}

	// 产生 service 阵列
	var service []string
	sStr, ok := config["service"] // 先确认 service 设定值是否存在
	if ok { // 如果 service 值 存在
		service = strings.Split(sStr, ",") // 以逗点分隔开来
	}
	if !ok { // 如果 service 值 不 存在
		err = fmt.Errorf("init XFileLog failed, not found service")
		return
	}

	// 产生 level 阵列
	var level []string
	lStr, ok := config["level"] // 先确认 level 设定值是否存在
	if ok { // 如果 level 值 存在
		level = strings.Split(lStr, ",") // 以逗点分隔开来
	}
	if !ok { // 如果 level 值 不 存在
		err = fmt.Errorf("init XFileLog failed, not found level")
		return
	}

	// 进行最后长度检查，因为是以 filename 阵列为中心跑回圈，所以任何一个阵列的长度要大于 filename 阵列

	// 检查 service 阵列
	// service 阵列可以只含有一个值，这时所有的档案该设定就会被统一设定
	// service 阵列可以只含多个值，但是数量要大于 filename 阵列，这时所有的档案该设定就会被个别设定
	if len(service) < len(filename) && len(service) != 1 {
		err = fmt.Errorf("init XFileLog failed, lack service config")
		return
	}

	// 检查 level 阵列
	// level 阵列可以只含有一个值，这时所有的档案该设定就会被统一设定
	// level 阵列可以只含多个值，但是数量要大于 filename 阵列，这时所有的档案该设定就会被个别设定
	if len(level) < len(filename) && len(level) != 1 {
		err = fmt.Errorf("init XFileLog failed, lack level config")
		return
	}

	// 以 filename 为中心做回圈
	for i := 0; i < len(filename); i++ {
		p := new(XFileLog)

		path, ok := config["path"]
		if !ok {
			err = fmt.Errorf("init XFileLog failed, not found path")
			return
		}

		if len(service) >= len(filename) && (len(service[i]) > 0) {
			p.service = service[i] // service 可以进行个别设定
		}
		if len(service) == 1  && (len(service[0]) > 0) {
			p.service = service[0] // service 可以进行统一设定
		}

		runtime, ok := config["runtime"]
		if !ok || runtime == "true" || runtime == "TRUE" {
			p.runtime = true
		} else {
			p.runtime = false
		}

		skip, _ := config["skip"]
		if len(skip) > 0 {
			skipNum, err := strconv.Atoi(skip)
			if err == nil {
				p.skip = skipNum
			}
		}

		isDir, err := isDir(path)
		if err != nil || !isDir {
			err = os.MkdirAll(path, 0755)
			if err != nil {
				return newError("Mkdir failed, err:%v", err)
			}
		}

		p.path = path

		p.filename = filename[i] // filename 可以进行个别设定

		if len(level) >= len(filename) {
			p.level = LevelFromStr(level[i]) // level 可以进行个别设定
		}
		if len(level) == 1 {
			p.level = LevelFromStr(level[0]) // level 可以进行统一设定
		}

		hostname, _ := os.Hostname()
		p.hostname = hostname
		body := func() {
			go p.spliter()
		}
		doSplit, ok := config["dosplit"]
		if !ok {
			doSplit = "true"
		}
		if doSplit == "true" {
			p.split.Do(body)
		}

		// 错误回传
		if p.ReOpen() != nil { // 一旦有错误，就回传错误
			return p.ReOpen()
		}

		// 完整组合 multi map
		ps.multi[filename[i]] = p
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
	// levelArr 阵列可以只含多个值，但是数量要大于 filename 阵列，这时所有的档案该设定就会被个别设定
	if len(levelArr) < len(ps.multi) && len(levelArr) != 1 {
		return // 立刻中断
	}

	// 以 filename 为中心做回圈
	for i := 0; i < len(ps.multi); i++ {
		for key := range ps.multi {
			ps.multi[key].level = LevelFromStr(levelArr[i])
		}
	}
}

// SetSkip 设定忽略重复的行数
func (ps *XMultiFileLog) SetSkip(skip int) {
	// 以 filename 为中心做回圈
	for i := 0; i < len(ps.multi); i++ {
		for key := range ps.multi {
			ps.multi[key].skip = skip // 进行统一设定 skip
		}
	}
}

// Debug 函式这样做不对，再看看不如做最好
func (ps *XMultiFileLog) Debug(format string, a ...interface{}) error {
	for key := range ps.multi {
		if err := ps.multi[key].debugx(XFileLogDefaultLogID, format, a...); err != nil {
			// 错误回传
			return err
		}
	}

	// 正确回传
	return nil
}
