package xlog

import (
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

// >>>>> >>>>> >>>>> >>>>> >>>>> 测试日志分流后的写入内容

// TestMultiFileLogContentCase2 为 用来测试日志分流后的写入内容 (包含所有的测试案例)
func TestMultiFileLogContent(t *testing.T) {
	// 2 个档案，1 个等级，1 个服务 的测试
	TestMultiFileLogContentCase1(t)
	// 2 个档案，2 个等级，1 个服务 的测试
	TestMultiFileLogContentCase2(t)
	// 2 个档案，2 个等级，2 个服务 的测试
	TestMultiFileLogContentCase3(t)
}

// TestMultiFileLogContentCase1 为 用来测试日志分流后的写入内容 (2 个档案，1 个等级，1 个服务)
func TestMultiFileLogContentCase1(t *testing.T) {
	// >>>>> >>>>> >>>>> >>>>> >>>>> 准备双向通道的测试环境

	// 准备摸拟设定值
	cfgSuit := make(map[string]string, 2) // 有几份档案就要准备几个模拟用的双向通道
	cfgSuit["filename"] = "log1,log2"     // 将要用双向通道摸拟 log1.log 和 log2.log 两份日志档案
	cfgSuit["chanSize"] = "5"             // 在每一个模拟用的双向通道可以接收 5 笔日志资料

	// 开始建立模拟环境
	err := testMultiAndFileSuite(cfgSuit) // 内含初始化的操作
	require.Equal(t, err, nil)

	// >>>>> >>>>> >>>>> >>>>> >>>>> 准备和初始化日志分流物件 (2 个档案，2 个等级，1 个服务)

	// 产生日志分流物件
	ps := new(XMultiFileLog)

	// 准备日志分流设定值
	cfg := make(map[string]string)
	cfg["path"] = "/home/panhong/go/src/github.com/panhong/demo/logs"
	cfg["filename"] = "log1,log2"
	cfg["level"] = "Notice"
	cfg["service"] = "svc1"
	cfg["skip"] = "5"

	// 初始化日志分流物件
	err = ps.Init(cfg)
	require.Equal(t, err, nil)

	// >>>>> >>>>> >>>>> >>>>> >>>>> 开始写入日志的操作

	// 日志格式为 档名::日志内容
	// 如果没有指定档名，将会写入预设档案
	// 预设档案为 filename 设定值的第一个元素，也就是 log1

	// 第一笔 Notice 日志
	err = ps.Notice("record1") // 此日志 会 写入档案内 (V)
	require.Equal(t, err, nil)

	// 第二笔 Notice 日志
	err = ps.Notice("log1::record2") // 此日志 会 写入档案内 (V)
	require.Equal(t, err, nil)

	// 第三笔 Notice 日志
	err = ps.Notice("log2::record3") // 此日志 会 写入档案内 (V)
	require.Equal(t, err, nil)

	// 第四笔 Notice 日志
	err = ps.Debug("log1::record4") // log1 和 log2 日志档案等级为 Notice，会把 Debug 的日志讯息忽略掉，所以此日志 不会 写入档案内 (X)
	require.Equal(t, err, nil)

	// 第五笔 Debug 日志
	err = ps.Debug("log2::record5") // log1 和 log2 日志档案等级为 Notice，会把 Debug 的日志讯息忽略掉，所以此日志 不会 写入档案内 (X)
	require.Equal(t, err, nil)

	// 第六笔 Notice 日志
	// 在设定中，没有 log6 这个日志档案，这时又不能把这个日志抛弃，所以就先存放在预设的日志中
	// 预设档案为 filename 设定值的第一个元素，也就是 log1
	err = ps.Notice("log6::record6") // 此日志 会 写入档案内 (V)
	require.Equal(t, err, nil)

	// >>>>> >>>>> >>>>> >>>>> >>>>> 最后检查日志的写入结果

	// 把整个模拟双向通道的内容取出
	ret := PrintMockChanMsg()

	// 数量检查
	require.Equal(t, len(ret), 4) // 写入五笔日志，但被忽略一笔，最后只会被写入四笔日志

	// 细节鐱查
	// 在模拟的双向通道的内容如下
	// log1::[2021-12-23 01:47:44] [svc1] [debian5] [NOTICE] [1000000001] [testing.tRunner:testing.go:1259] record1
	// log1::[2021-12-23 01:47:44] [svc1] [debian5] [NOTICE] [1000000001] [testing.tRunner:testing.go:1259] record2
	// log1::[2021-12-23 01:47:44] [svc1] [debian5] [NOTICE] [1000000001] [testing.tRunner:testing.go:1259] record6
	// log2::[2021-12-23 01:47:44] [svc1] [debian5] [NOTICE] [1000000001] [testing.tRunner:testing.go:1259] record3

	// 检查由模拟的双向通道中取出资料内容
	require.Equal(t,
		// 在日志档案 log1 里，服务的名称要为svc1，等级为 NOTICE，写入一笔记录为 record1
		strings.Contains(ret[0], "log1") &&
			strings.Contains(ret[0], "svc1") &&
			strings.Contains(ret[0], "NOTICE") &&
			strings.Contains(ret[0], "record1"),
		true)

	require.Equal(t,
		// 在日志档案 log1 里，服务的名称要为svc1，等级为 NOTICE，写入一笔记录为 record2
		strings.Contains(ret[1], "log1") &&
			strings.Contains(ret[1], "svc1") &&
			strings.Contains(ret[1], "NOTICE") &&
			strings.Contains(ret[1], "record2"),
		true)

	require.Equal(t,
		// 在日志档案 log1 里，服务的名称要为svc1，等级为 NOTICE，写入一笔记录为 record6
		strings.Contains(ret[2], "log1") &&
			strings.Contains(ret[2], "svc1") &&
			strings.Contains(ret[2], "NOTICE") &&
			strings.Contains(ret[2], "record6"),
		true)

	require.Equal(t,
		// 在日志档案 log1 里，服务的名称要为svc1，等级为 DEBUG，写入一笔记录为 record5
		strings.Contains(ret[3], "log2") &&
			strings.Contains(ret[3], "svc1") &&
			strings.Contains(ret[3], "NOTICE") &&
			strings.Contains(ret[3], "record3"),
		true)

	// 关闭所有模拟用的双向通道
	CloseGlobalMockMultiXLogFile()
}

// TestMultiFileLogContentCase2 为 用来测试日志分流后的写入内容 (2 个档案，2 个等级，1 个服务)
func TestMultiFileLogContentCase2(t *testing.T) {
	// >>>>> >>>>> >>>>> >>>>> >>>>> 准备双向通道的测试环境

	// 准备摸拟设定值
	cfgSuit := make(map[string]string, 2) // 有几份档案就要准备几个模拟用的双向通道
	cfgSuit["filename"] = "log1,log2"     // 将要用双向通道摸拟 log1.log 和 log2.log 两份日志档案
	cfgSuit["chanSize"] = "5"             // 在每一个模拟用的双向通道可以接收 5 笔日志资料

	// 开始建立模拟环境
	err := testMultiAndFileSuite(cfgSuit) // 内含初始化的操作
	require.Equal(t, err, nil)

	// >>>>> >>>>> >>>>> >>>>> >>>>> 准备和初始化日志分流物件 (2 个档案，2 个等级，1 个服务)

	// 产生日志分流物件
	ps := new(XMultiFileLog)

	// 准备日志分流设定值
	cfg := make(map[string]string)
	cfg["path"] = "/home/panhong/go/src/github.com/panhong/demo/logs"
	cfg["filename"] = "log1,log2"
	cfg["level"] = "Notice,Debug"
	cfg["service"] = "svc1"
	cfg["skip"] = "5"

	// 初始化日志分流物件
	err = ps.Init(cfg)
	require.Equal(t, err, nil)

	// >>>>> >>>>> >>>>> >>>>> >>>>> 开始写入日志的操作

	// 日志格式为 档名::日志内容
	// 如果没有指定档名，将会写入预设档案
	// 预设档案为 filename 设定值的第一个元素，也就是 log1

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
	err = ps.Debug("log1::record4") // log1 日志档案等级为 Notice，会把 Debug 的日志讯息忽略掉，所以此日志 不会 写入档案内 (X)
	require.Equal(t, err, nil)

	// 第五笔 Debug 日志
	err = ps.Debug("log2::record5") // log1 日志档案等级为 Debug，会显示整个 Debug 的日志讯息，所以此日志 会 写入档案内 (V)
	require.Equal(t, err, nil)

	// 第六笔 Notice 日志
	// 在设定中，没有 log6 这个日志档案，这时又不能把这个日志抛弃，所以就先存放在预设的日志中
	// 预设档案为 filename 设定值的第一个元素，也就是 log1
	err = ps.Notice("log6::record6") // 此日志 会 写入档案内 (V)
	require.Equal(t, err, nil)

	// >>>>> >>>>> >>>>> >>>>> >>>>> 最后检查日志的写入结果

	// 把整个模拟双向通道的内容取出
	ret := PrintMockChanMsg()

	// 数量检查
	require.Equal(t, len(ret), 5) // 写入五笔日志，但被忽略一笔，最后只会被写入四笔日志

	// 细节鐱查
	// 在模拟的双向通道的内容如下
	// log1::[2021-12-22 23:30:43] [svc1] [debian5] [NOTICE] [1000000001] [testing.tRunner:testing.go:1259] record1
	// log1::[2021-12-22 23:30:43] [svc1] [debian5] [NOTICE] [1000000001] [testing.tRunner:testing.go:1259] record2
	// log1::[2021-12-22 23:30:43] [svc1] [debian5] [NOTICE] [1000000001] [testing.tRunner:testing.go:1259] record6
	// log2::[2021-12-22 23:30:43] [svc1] [debian5] [DEBUG]  [1000000001] [testing.tRunner:testing.go:1259] record5
	// log2::[2021-12-22 23:30:43] [svc1] [debian5] [NOTICE] [1000000001] [testing.tRunner:testing.go:1259] record3

	// 检查由模拟的双向通道中取出资料内容
	require.Equal(t,
		// 在日志档案 log1 里，服务的名称要为svc1，等级为 NOTICE，写入一笔记录为 record1
		strings.Contains(ret[0], "log1") &&
			strings.Contains(ret[0], "svc1") &&
			strings.Contains(ret[0], "NOTICE") &&
			strings.Contains(ret[0], "record1"),
		true)

	require.Equal(t,
		// 在日志档案 log1 里，服务的名称要为svc1，等级为 NOTICE，写入一笔记录为 record2
		strings.Contains(ret[1], "log1") &&
			strings.Contains(ret[1], "svc1") &&
			strings.Contains(ret[1], "NOTICE") &&
			strings.Contains(ret[1], "record2"),
		true)

	require.Equal(t,
		// 在日志档案 log1 里，服务的名称要为svc1，等级为 NOTICE，写入一笔记录为 record6
		strings.Contains(ret[2], "log1") &&
			strings.Contains(ret[2], "svc1") &&
			strings.Contains(ret[2], "NOTICE") &&
			strings.Contains(ret[2], "record6"),
		true)

	require.Equal(t,
		// 在日志档案 log1 里，服务的名称要为svc1，等级为 DEBUG，写入一笔记录为 record5
		strings.Contains(ret[3], "log2") &&
			strings.Contains(ret[3], "svc1") &&
			strings.Contains(ret[3], "DEBUG") &&
			strings.Contains(ret[3], "record5"),
		true)

	require.Equal(t,
		// 在日志档案 log1 里，服务的名称要为svc1，等级为 NOTICE，写入一笔记录为 record3
		strings.Contains(ret[4], "log2") &&
			strings.Contains(ret[4], "svc1") &&
			strings.Contains(ret[4], "NOTICE") &&
			strings.Contains(ret[4], "record3"),
		true)

	// 关闭所有模拟用的双向通道
	CloseGlobalMockMultiXLogFile()
}

// TestMultiFileLogContentCase3 为 用来测试日志分流后的写入内容 (2 个档案，2 个等级，2 个服务)
func TestMultiFileLogContentCase3(t *testing.T) {
	// >>>>> >>>>> >>>>> >>>>> >>>>> 准备双向通道的测试环境

	// 准备摸拟设定值
	cfgSuit := make(map[string]string, 2) // 有几份档案就要准备几个模拟用的双向通道
	cfgSuit["filename"] = "log1,log2"     // 将要用双向通道摸拟 log1.log 和 log2.log 两份日志档案
	cfgSuit["chanSize"] = "5"             // 在每一个模拟用的双向通道可以接收 5 笔日志资料

	// 开始建立模拟环境
	err := testMultiAndFileSuite(cfgSuit) // 内含初始化的操作
	require.Equal(t, err, nil)

	// >>>>> >>>>> >>>>> >>>>> >>>>> 准备和初始化日志分流物件 (2 个档案，2 个等级，2 个服务)

	// 产生日志分流物件
	ps := new(XMultiFileLog)

	// 准备日志分流设定值
	cfg := make(map[string]string)
	cfg["path"] = "/home/panhong/go/src/github.com/panhong/demo/logs"
	cfg["filename"] = "log1,log2"
	cfg["level"] = "Notice,Debug"
	cfg["service"] = "svc1,svc2"
	cfg["skip"] = "5"

	// 初始化日志分流物件
	err = ps.Init(cfg)
	require.Equal(t, err, nil)

	// >>>>> >>>>> >>>>> >>>>> >>>>> 开始写入日志的操作

	// 日志格式为 档名::日志内容
	// 如果没有指定档名，将会写入预设档案
	// 预设档案为 filename 设定值的第一个元素，也就是 log1

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
	err = ps.Debug("log1::record4") // log1 日志档案等级为 Notice，会把 Debug 的日志讯息忽略掉，所以此日志 不会 写入档案内 (X)
	require.Equal(t, err, nil)

	// 第五笔 Debug 日志
	err = ps.Debug("log2::record5") // log1 日志档案等级为 Debug，会显示整个 Debug 的日志讯息，所以此日志 会 写入档案内 (V)
	require.Equal(t, err, nil)

	// 第六笔 Notice 日志
	// 在设定中，没有 log6 这个日志档案，这时又不能把这个日志抛弃，所以就先存放在预设的日志中
	// 预设档案为 filename 设定值的第一个元素，也就是 log1
	err = ps.Notice("log6::record6") // 此日志 会 写入档案内 (V)
	require.Equal(t, err, nil)

	// >>>>> >>>>> >>>>> >>>>> >>>>> 最后检查日志的写入结果

	// 把整个模拟双向通道的内容取出
	ret := PrintMockChanMsg()

	// 数量检查
	require.Equal(t, len(ret), 5) // 写入五笔日志，但被忽略一笔，最后只会被写入四笔日志

	// 细节鐱查
	// 在模拟的双向通道的内容如下
	// log1::[2021-12-22 23:30:43] [svc1] [debian5] [NOTICE] [1000000001] [testing.tRunner:testing.go:1259] record1
	// log1::[2021-12-22 23:30:43] [svc1] [debian5] [NOTICE] [1000000001] [testing.tRunner:testing.go:1259] record2
	// log1::[2021-12-22 23:30:43] [svc1] [debian5] [NOTICE] [1000000001] [testing.tRunner:testing.go:1259] record6
	// log2::[2021-12-22 23:30:43] [svc1] [debian5] [DEBUG]  [1000000001] [testing.tRunner:testing.go:1259] record5
	// log2::[2021-12-22 23:30:43] [svc1] [debian5] [NOTICE] [1000000001] [testing.tRunner:testing.go:1259] record3

	// 检查由模拟的双向通道中取出资料内容
	require.Equal(t,
		// 在日志档案 log1 里，服务的名称要为svc1，等级为 NOTICE，写入一笔记录为 record1
		strings.Contains(ret[0], "log1") &&
			strings.Contains(ret[0], "svc1") &&
			strings.Contains(ret[0], "NOTICE") &&
			strings.Contains(ret[0], "record1"),
		true)

	require.Equal(t,
		// 在日志档案 log1 里，服务的名称要为svc1，等级为 NOTICE，写入一笔记录为 record2
		strings.Contains(ret[1], "log1") &&
			strings.Contains(ret[1], "svc1") &&
			strings.Contains(ret[1], "NOTICE") &&
			strings.Contains(ret[1], "record2"),
		true)

	require.Equal(t,
		// 在日志档案 log1 里，服务的名称要为svc1，等级为 NOTICE，写入一笔记录为 record6
		strings.Contains(ret[2], "log1") &&
			strings.Contains(ret[2], "svc1") &&
			strings.Contains(ret[2], "NOTICE") &&
			strings.Contains(ret[2], "record6"),
		true)

	require.Equal(t,
		// 在日志档案 log1 里，服务的名称要为svc1，等级为 DEBUG，写入一笔记录为 record5
		strings.Contains(ret[3], "log2") &&
			strings.Contains(ret[3], "svc2") &&
			strings.Contains(ret[3], "DEBUG") &&
			strings.Contains(ret[3], "record5"),
		true)

	require.Equal(t,
		// 在日志档案 log1 里，服务的名称要为svc1，等级为 NOTICE，写入一笔记录为 record3
		strings.Contains(ret[4], "log2") &&
			strings.Contains(ret[4], "svc2") &&
			strings.Contains(ret[4], "NOTICE") &&
			strings.Contains(ret[4], "record3"),
		true)

	// 关闭所有模拟用的双向通道
	CloseGlobalMockMultiXLogFile()
}
