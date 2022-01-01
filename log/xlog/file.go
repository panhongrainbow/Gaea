// Copyright 2019 The Gaea Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package xlog

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// XFileLog is the file logger
type XFileLog struct { // 单档案的输出
	level    int         // 日志等级
	skip     int         // 略过等级
	runtime  bool        // 本机资讯
	hostname string      // 本机名称
	service  string      // 服务名称
	storage  *LogStorage // 服务名称
}

// constants of XFileLog
const (
	XFileLogDefaultLogID = "900000001"
	SpliterDelay         = 5
	CleanDays            = -3
)

// NewXFileLog is the constructor of XFileLog
func NewXFileLog() XLogger {
	return &XFileLog{
		skip: XLogDefSkipNum,
	}
}

// Init implements XLogger
func (p *XFileLog) Init(config map[string]string) (err error) {

	// 在这里指定储存方式
	p.storage = NewLogStorageClient(config)

	// 以下会有部份的设定值移到储存代码里

	/*path, ok := config["path"] // (移到储存)
	if !ok {
		err = fmt.Errorf("init XFileLog failed, not found path")
		return
	}*/

	/*filename, ok := config["filename"] // (移到储存)
	if !ok {
		err = fmt.Errorf("init XFileLog failed, not found filename")
		return
	}*/

	level, ok := config["level"]
	if !ok {
		err = fmt.Errorf("init XFileLog failed, not found level")
		return
	}

	service, _ := config["service"]
	if len(service) > 0 {
		p.service = service
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

	/*isDir, err := isDir(path) // (移到储存)
	if err != nil || !isDir {
		err = os.MkdirAll(path, 0755)
		if err != nil {
			return newError("Mkdir failed, err:%v", err)
		}
	}*/

	/*p.path = path // (移到储存)
	p.filename = filename*/

	p.level = LevelFromStr(level)
	hostname, _ := os.Hostname()
	p.hostname = hostname

	/*body := func() { // (移到储存)
		go p.spliter()
	}
	doSplit, ok := config["dosplit"]
	if !ok {
		doSplit = "true"
	}
	if doSplit == "true" {
		p.split.Do(body)
	}*/

	return p.ReOpen()
}

// SetLevel implements XLogger
func (p *XFileLog) SetLevel(level string) {
	p.level = LevelFromStr(level)
}

// SetSkip implements XLogger
func (p *XFileLog) SetSkip(skip int) {
	p.skip = skip
}

// ReOpen implements XLogger
// 用于重新开档的函式，在这里会真的把所有的日志档案开启
// 如果是在单元测试，会检查所有的模拟的双向通道是否存在
func (p *XFileLog) ReOpen() error {
	return p.storage.client.ReOpen() // 交由储存物件去开档
}

// Close implements XLogger
func (p *XFileLog) Close() {
	_ = p.storage.client.Close()
}

// Warn implements XLogger
func (p *XFileLog) Warn(format string, a ...interface{}) error {
	if p.level > WarnLevel {
		return nil
	}

	return p.warnx(XFileLogDefaultLogID, format, a...)
}

// Warnx implements XLogger
func (p *XFileLog) Warnx(logID, format string, a ...interface{}) error {
	if p.level > WarnLevel {
		return nil
	}

	return p.warnx(logID, format, a...)
}

func (p *XFileLog) warnx(logID, format string, a ...interface{}) error {
	logText := formatValue(format, a...)
	fun, filename, lineno := getRuntimeInfo(p.skip)
	logText = formatLineInfo(p.runtime, fun, filepath.Base(filename), logText, lineno)
	//logText = fmt.Sprintf("[%s:%s:%d] %s", fun, filepath.Base(fileName), lineno, logText)

	return p.write(WarnLevel, &logText, logID)
}

// Fatal implements XLogger
func (p *XFileLog) Fatal(format string, a ...interface{}) error {
	if p.level > FatalLevel {
		return nil
	}

	return p.fatalx(XFileLogDefaultLogID, format, a...)
}

// Fatalx implements XLogger
func (p *XFileLog) Fatalx(logID, format string, a ...interface{}) error {
	if p.level > FatalLevel {
		return nil
	}

	return p.fatalx(logID, format, a...)
}

func (p *XFileLog) fatalx(logID, format string, a ...interface{}) error {
	logText := formatValue(format, a...)
	fun, filename, lineno := getRuntimeInfo(p.skip)
	logText = formatLineInfo(p.runtime, fun, filepath.Base(filename), logText, lineno)
	//logText = fmt.Sprintf("[%s:%s:%d] %s", fun, filepath.Base(fileName), lineno, logText)

	return p.write(FatalLevel, &logText, logID)
}

// Notice implements XLogger
func (p *XFileLog) Notice(format string, a ...interface{}) error {
	if p.level > NoticeLevel {
		return nil
	}
	return p.noticex(XFileLogDefaultLogID, format, a...)
}

// Noticex implements XLogger
func (p *XFileLog) Noticex(logID, format string, a ...interface{}) error {
	if p.level > NoticeLevel {
		return nil
	}
	return p.noticex(logID, format, a...)
}

func (p *XFileLog) noticex(logID, format string, a ...interface{}) error {
	logText := formatValue(format, a...)
	fun, filename, lineno := getRuntimeInfo(p.skip)
	logText = formatLineInfo(p.runtime, fun, filepath.Base(filename), logText, lineno)

	return p.write(NoticeLevel, &logText, logID)
}

// Trace implements XLogger
func (p *XFileLog) Trace(format string, a ...interface{}) error {
	return p.tracex(XFileLogDefaultLogID, format, a...)
}

// Tracex implements XLogger
func (p *XFileLog) Tracex(logID, format string, a ...interface{}) error {
	return p.tracex(logID, format, a...)
}

func (p *XFileLog) tracex(logID, format string, a ...interface{}) error {
	if p.level > TraceLevel {
		return nil
	}

	logText := formatValue(format, a...)
	fun, filename, lineno := getRuntimeInfo(p.skip)
	logText = formatLineInfo(p.runtime, fun, filepath.Base(filename), logText, lineno)
	//logText = fmt.Sprintf("[%s:%s:%d] %s", fun, filepath.Base(fileName), lineno, logText)

	return p.write(TraceLevel, &logText, logID)
}

// Debug implements XLogger
func (p *XFileLog) Debug(format string, a ...interface{}) error {
	return p.debugx(XFileLogDefaultLogID, format, a...)
}

// Debugx implements XLogger
func (p *XFileLog) Debugx(logID, format string, a ...interface{}) error {
	return p.debugx(logID, format, a...)
}

func (p *XFileLog) debugx(logID, format string, a ...interface{}) error {
	if p.level > DebugLevel {
		return nil
	}

	logText := formatValue(format, a...)
	fun, filename, lineno := getRuntimeInfo(p.skip)
	logText = formatLineInfo(p.runtime, fun, filepath.Base(filename), logText, lineno)

	return p.write(DebugLevel, &logText, logID)
}

// write 为最后的写入函式，会把日志写入档案
// 如果是在单元测试，会把日志写入模拟的双向通道内
func (p *XFileLog) write(level int, msg *string, logID string) error {
	levelText := levelTextArray[level]
	time := time.Now().Format("2006-01-02 15:04:05")

	logText := formatLog(msg, time, p.service, p.hostname, levelText, logID)

	if level >= WarnLevel {
		_ = p.storage.client.WriteErr([]byte(logText)) // 交由储存物件写入错误日志
	}

	_ = p.storage.client.Write([]byte(logText)) // 交由储存物件写入日志

	// 正确回传
	return nil
}
