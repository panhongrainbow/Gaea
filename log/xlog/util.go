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
	"bytes"
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

// log level
const (
	DebugLevel = iota
	TraceLevel
	NoticeLevel
	WarnLevel
	FatalLevel
	NoneLevel
)

// log skip num
const (
	XLogDefSkipNum = 4
)

var (
	levelTextArray = []string{
		DebugLevel:  "DEBUG",
		TraceLevel:  "TRACE",
		NoticeLevel: "NOTICE",
		WarnLevel:   "WARN",
		FatalLevel:  "FATAL",
	}
)

// LevelFromStr get log level from level string
func LevelFromStr(level string) int {
	resultLevel := DebugLevel
	levelLower := strings.ToLower(level)
	switch levelLower {
	case "debug":
		resultLevel = DebugLevel
	case "trace":
		resultLevel = TraceLevel
	case "notice":
		resultLevel = NoticeLevel
	case "warn":
		resultLevel = WarnLevel
	case "fatal":
		resultLevel = FatalLevel
	case "none":
		resultLevel = NoneLevel
	default:
		resultLevel = NoticeLevel
	}

	return resultLevel
}

// FileFormatFromStr 可以从含有 :: 的字串中，分解出 日志档名 和 日志格式
func FileFormatFromStr(format string) (string, string) {
	// 拆分字串
	arr := strings.Split(format, "::")

	// 正确回传
	if len(arr) == 2 {
		return arr[0], arr[1] // 如果有两个元素就两个都回传
	}
	if len(arr) == 1 {
		return "", arr[0] // 如果只有一个元素，就放在 format 的回传位置回传
	}

	// 错误回传
	return "", ""
}

func getRuntimeInfo(skip int) (function, filename string, lineno int) {
	function = "???"
	pc, filename, lineno, ok := runtime.Caller(skip)
	if ok {
		function = runtime.FuncForPC(pc).Name()
	}
	return
}

func formatLog(body *string, fields ...string) string {
	var buffer bytes.Buffer
	for _, v := range fields {
		buffer.WriteString("[")
		buffer.WriteString(v)
		buffer.WriteString("] ")
	}

	buffer.WriteString(*body)
	buffer.WriteString("\n")

	return buffer.String()
}

func formatValue(format string, a ...interface{}) (result string) {
	if len(a) == 0 {
		result = format
		return
	}

	result = fmt.Sprintf(format, a...)
	return
}

func formatLineInfo(runtime bool, functionName, filename, logText string, lineno int) string {
	var buffer bytes.Buffer
	if runtime {
		buffer.WriteString("[")
		buffer.WriteString(functionName)
		buffer.WriteString(":")

		buffer.WriteString(filename)
		buffer.WriteString(":")

		buffer.WriteString(strconv.FormatInt(int64(lineno), 10))
		buffer.WriteString("] ")
	}
	buffer.WriteString(logText)

	return buffer.String()
}

func newError(format string, a ...interface{}) error {
	err := fmt.Sprintf(format, a...)
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		return errors.New(err)
	}

	function := runtime.FuncForPC(pc).Name()
	msg := fmt.Sprintf("%s func:%s file:%s line:%d",
		err, function, file, line)
	return errors.New(msg)
}
