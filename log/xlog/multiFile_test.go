package xlog

import (
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

// testMultiFileSuite 建立测试环境
func testMultiFileSuite(config map[string]string) error {
	// 初始化模拟文件
	mf := new(MockMultiXLogFile)
	err := mf.Init(config) // 在这里进行初始化模拟文件
	if err != nil {
		return err // 错误回传
	}

	// 进行重新开启档案
	// 表面上是重新开启档案，但实际上是检查双向通道是否有全部被开启
	err = mf.Open()
	if err != nil {
		return err // 错误回传
	}

	// 设定全域的双向通道
	SetGlobalMockMultiXLogFile(mf)

	// 正确回传
	return nil
}

// TestMultiFileLogLevel 为 用来测试日志分流后的日志分级
func TestMultiFileLogLevel(t *testing.T) {
	// >>>>> >>>>> >>>>> >>>>> >>>>> 准备双向通道的测试环境

	// 准备摸拟设定值
	cfgSuit := make(map[string]string, 1)
	cfgSuit["filename"] = "log1,log2" // 将要用双向通道摸拟 log1.log 和 log2.log 两份日志档案
	cfgSuit["chanSize"] = "3"

	// 开始建立模拟环境
	err := testMultiFileSuite(cfgSuit) // 内含初始化的操作
	require.Equal(t, err, nil)

	// >>>>> >>>>> >>>>> >>>>> >>>>> 准备和初始化日志分流物件

	// 产生日志分流物件
	ps := new(XMultiFileLog)

	// 准备日志分流设定值
	cfg := make(map[string]string)
	cfg["path"] = "/home/panhong/go/src/github.com/panhong/demo/logs"
	cfg["filename"] = "log1,log2"
	cfg["level"] = "Notice,Debug"
	cfg["service"] = "result-get"
	cfg["skip"] = "5"

	// 初始化日志分流物件
	err = ps.Init(cfg)
	require.Equal(t, err, nil)

	// >>>>> >>>>> >>>>> >>>>> >>>>> 开始写入日志的操作

	// 日志格式为 档名::日志内容
	// 如果没有指定档名，将会写入预设档案
	// 预设档案为 fileName 设定值的第一个元素

	// 第一笔 Notice 日志
	err = ps.Notice("record1") // 此日志 会 写入档案内 (V)
	require.Equal(t, err, nil)

	// 第二笔 Notice 日志
	err = ps.Notice("log1::record2") // 此日志 会 写入档案内 (V)
	require.Equal(t, err, nil)

	// 第三笔 Notice 日志
	err = ps.Notice("log2::record3") // 此日志 会 写入档案内 (V)
	require.Equal(t, err, nil)

	// 第四笔 Debug 日志
	err = ps.Debug("log1::record4") // log1 日志档案 log1 等级为 Notice，会把 Debug 的日志讯息忽略掉，所以此日志 不会 写入档案内 (X)
	require.Equal(t, err, nil)

	// 第五笔 Debug 日志
	err = ps.Debug("log2::record5") // log1 日志档案 log2 等级为 Debug，会显示整个 Debug 的日志讯息，所以此日志 会 写入档案内 (V)
	require.Equal(t, err, nil)

	// 第六笔 Notice 日志
	// 在设定中，没有 log6 这个日志档案，这时又不能把这个日志抛弃，所以就先存放在预设的日志中
	//
	err = ps.Notice("log6::record6") // 此日志 会 写入档案内 (V)
	require.Equal(t, err, nil)

	// >>>>> >>>>> >>>>> >>>>> >>>>> 最后检查日志的写入结果

	ret := PrintMockChanMsg()

	// 数量检查
	require.Equal(t, len(ret), 5) // 写入五笔日志，但被忽略一笔，最后只会被写入四笔日志

	// 细节鐱查

	require.Equal(t,
		strings.Contains(ret[0], "log1") && // 在日志档案 log1 里，写入一笔记录为 record1
		strings.Contains(ret[0], "record1"),
		true)

	require.Equal(t,
		strings.Contains(ret[1], "log1") && // 在日志档案 log1 里，写入一笔记录为 record2
			strings.Contains(ret[1], "record2"),
		true)

	require.Equal(t,
		strings.Contains(ret[2], "log1") && // 在日志档案 log1 里，写入一笔记录为 record6
			strings.Contains(ret[2], "record6"),
		true)

	require.Equal(t,
		strings.Contains(ret[3], "log2") && // 在日志档案 log1 里，写入一笔记录为 record6
			strings.Contains(ret[3], "record5"),
		true)

	require.Equal(t,
		strings.Contains(ret[4], "log2") && // 在日志档案 log1 里，写入一笔记录为 record3
			strings.Contains(ret[4], "record3"),
		true)
}