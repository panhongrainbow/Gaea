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

// XLogger declares method that log instance should implement.
type XLogger interface {
	// Init logger
	Init(config map[string]string) error

	// ReOpen logger
	ReOpen() error

	// SetLevel 为设置日志级别, 级别如下: "Debug", "Trace", "Notice", "Warn", "Fatal", "None"
	SetLevel(level string)

	// SetSkip 为 set skip
	SetSkip(skip int)

	// Debug 为打印 Debug 日志. 当日志级别大于 Debug 时, 不会输出任何日志
	Debug(format string, a ...interface{}) error

	// Trace 为打印 Trace 日志. 当日志级别大于 Trace 时, 不会输出任何日志
	Trace(format string, a ...interface{}) error

	// Notice 为打印 Notice 日志. 当日志级别大于 Notice 时, 不会输出任何日志
	Notice(format string, a ...interface{}) error

	// Warn 为打印 Warn 日志. 当日志级别大于 Warn 时, 不会输出任何日志
	Warn(format string, a ...interface{}) error

	// Fatal 为打印 Fatal 日志. 当日志级别大于 Fatal 时, 不会输出任何日志
	Fatal(format string, a ...interface{}) error

	// Debugx 为打印 Debug 日志, 需要传入 logID. 当日志级别大于 Debug 时, 不会输出任何日志
	Debugx(logID, format string, a ...interface{}) error

	// Tracex 为打印 Trace 日志, 需要传入 logID. 当日志级别大于 Trace 时, 不会输出任何日志
	Tracex(logID, format string, a ...interface{}) error

	// Noticex 为打 印Notice 日志, 需要传入 logID. 当日志级别大于 Notice 时, 不会输出任何日志
	Noticex(logID, format string, a ...interface{}) error

	// Warnx 为打印 Warn 日志, 需要传入 logID. 当日志级别大于 Warn 时, 不会输出任何日志
	Warnx(logID, format string, a ...interface{}) error

	// Fatalx 为打印 Fatal 日志, 需要传入 logID. 当日志级别大于 Fatal 时, 不会输出任何日志
	Fatalx(logID, format string, a ...interface{}) error

	// Close 为关闭日志库. 注意: 如果没有调用 Close() 关闭日志库的话, 将会造成文件句柄泄露
	Close()
}
