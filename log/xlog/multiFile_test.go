package xlog

import (
	"github.com/XiaoMi/Gaea/models/logStorage/channel"
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
	// >>>>> >>>>> >>>>> >>>>> >>>>> 准备和初始化日志分流物件

	// 产生日志分流物件
	ps := new(XMultiFileLog)

	// 准备日志分流设定值
	cfg := make(map[string]string)
	cfg["path"] = "/home/panhong/go/src/github.com/panhong/demo/logs"
	cfg["filename"] = "log1,log2"
	cfg["level"] = "Notice"
	cfg["service"] = "svc1"
	cfg["skip"] = "5"
	cfg["storage"] = "channel"
	cfg["chan_size"] = "10"

	// 初始化日志分流物件
	err := ps.Init(cfg)
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

	// >>>>> >>>>> >>>>> >>>>> >>>>> 最后检查日志 log1 的写入结果

	// 把整个模拟双向通道的内容取出，会依照写入排列去排列，排列如下
	// [
	// [2022-01-05 15:39:54] [svc1] [debian5] [NOTICE] [1000000001] [testing.tRunner:testing.go:1410] record1 // record1 日志会写到预设的 log1 文档
	// [2022-01-05 15:39:54] [svc1] [debian5] [NOTICE] [1000000001] [testing.tRunner:testing.go:1410] record2 // record2 日志会写到指定的 log1 文档
	// [2022-01-05 15:39:54] [svc1] [debian5] [DEBUG]  [1000000001] [testing.tRunner:testing.go:1410] record4 // 因为目前的等级为 Notice ， Debug 等服的日志会被忽略
	// [2022-01-05 15:39:54] [svc1] [debian5] [NOTICE] [1000000001] [testing.tRunner:testing.go:1410] record6 // log6 日志文档不存在，record1 日志写到 log1 文档
	// ]

	// 把整个模拟双向通道的内容取出，会依照写入排列去排列，排列如下

	// 把整个模拟双向通道的内容取出
	ret1 := ps.multi["log1"].storage.client.(*channel.Client).Read("log1")

	// 数量检查
	require.Equal(t, len(ret1), 3) // 写入四笔日志，但被忽略一笔，最后只会被写入三笔日志

	// 检查由模拟的双向通道中取出资料内容
	require.Equal(t,
		// 在日志档案 log1 里，服务的名称要为svc1，等级为 NOTICE，写入一笔记录为 record1
		strings.Contains(ret1[0], "svc1") &&
			strings.Contains(ret1[0], "NOTICE") &&
			strings.Contains(ret1[0], "record1"),
		true)

	require.Equal(t,
		// 在日志档案 log1 里，服务的名称要为svc1，等级为 NOTICE，写入一笔记录为 record2
		strings.Contains(ret1[1], "svc1") &&
			strings.Contains(ret1[1], "NOTICE") &&
			strings.Contains(ret1[1], "record2"),
		true)

	require.Equal(t,
		// 在日志档案 log1 里，服务的名称要为svc1，等级为 NOTICE，写入一笔记录为 record6
		strings.Contains(ret1[2], "svc1") &&
			strings.Contains(ret1[2], "NOTICE") &&
			strings.Contains(ret1[2], "record6"),
		true)

	// >>>>> >>>>> >>>>> >>>>> >>>>> 最后检查日志 log2 的写入结果

	// 把整个模拟双向通道的内容取出，会依照写入排列去排列，排列如下
	// [
	// [2022-01-05 16:54:36] [svc1] [debian5] [NOTICE] [1000000001] [testing.tRunner:testing.go:1410] record3 // record3 日志会写到预设的 log2 文档
	// [2022-01-05 16:54:36] [svc1] [debian5] [DEBUG]  [1000000001] [testing.tRunner:testing.go:1410] record5 // log2 日志文档 等级为 Notice，会忽略 Debug 等服的日志
	// ]

	// 把整个模拟双向通道的内容取出，会依照写入排列去排列，排列如下

	// 把整个模拟双向通道的内容取出
	ret2 := ps.multi["log2"].storage.client.(*channel.Client).Read("log2")

	// 检查由模拟的双向通道中取出资料内容
	require.Equal(t,
		// 在日志档案 log2 里，服务的名称要为svc1，等级为 NOTICE，写入一笔记录为 record3
		strings.Contains(ret2[0], "svc1") &&
			strings.Contains(ret2[0], "NOTICE") &&
			strings.Contains(ret2[0], "record3"),
		true)

	// 关闭所有模拟用的双向通道
	ps.Close()
}

// TestMultiFileLogContentCase2 为 用来测试日志分流后的写入内容 (2 个档案，2 个等级，1 个服务)
func TestMultiFileLogContentCase2(t *testing.T) {
	// >>>>> >>>>> >>>>> >>>>> >>>>> 准备和初始化日志分流物件

	// 产生日志分流物件
	ps := new(XMultiFileLog)

	// 准备日志分流设定值
	cfg := make(map[string]string)
	cfg["path"] = "/home/panhong/go/src/github.com/panhong/demo/logs"
	cfg["filename"] = "log1,log2"
	cfg["level"] = "Notice,Debug"
	cfg["service"] = "svc1"
	cfg["skip"] = "5"
	cfg["storage"] = "channel"
	cfg["chan_size"] = "10"

	// 初始化日志分流物件
	err := ps.Init(cfg)
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

	// >>>>> >>>>> >>>>> >>>>> >>>>> 最后检查日志 log2 的写入结果

	// 把整个模拟双向通道的内容取出，会依照写入排列去排列，排列如下
	// [
	// [2022-01-05 15:39:54] [svc1] [debian5] [NOTICE] [1000000001] [testing.tRunner:testing.go:1410] record1 // record1 日志会写到预设的 log1 文档
	// [2022-01-05 15:39:54] [svc1] [debian5] [NOTICE] [1000000001] [testing.tRunner:testing.go:1410] record2 // record2 日志会写到指定的 log1 文档
	// [2022-01-05 15:39:54] [svc1] [debian5] [DEBUG]  [1000000001] [testing.tRunner:testing.go:1410] record4 // 因为目前的等级为 Notice ， Debug 等服的日志会被忽略
	// [2022-01-05 15:39:54] [svc1] [debian5] [NOTICE] [1000000001] [testing.tRunner:testing.go:1410] record6 // log6 日志文档不存在，record1 日志写到 log1 文档
	// ]

	// 把整个模拟双向通道的内容取出，会依照写入排列去排列，排列如下

	// 把整个模拟双向通道的内容取出
	ret1 := ps.multi["log1"].storage.client.(*channel.Client).Read("log1")

	// 数量检查
	require.Equal(t, len(ret1), 3) // 写入五笔日志，但被忽略一笔，最后只会被写入四笔日志

	// 检查由模拟的双向通道中取出资料内容
	require.Equal(t,
		// 在日志档案 log1 里，服务的名称要为svc1，等级为 NOTICE，写入一笔记录为 record1
		strings.Contains(ret1[0], "svc1") &&
			strings.Contains(ret1[0], "NOTICE") &&
			strings.Contains(ret1[0], "record1"),
		true)

	require.Equal(t,
		// 在日志档案 log1 里，服务的名称要为svc1，等级为 NOTICE，写入一笔记录为 record2
		strings.Contains(ret1[1], "svc1") &&
			strings.Contains(ret1[1], "NOTICE") &&
			strings.Contains(ret1[1], "record2"),
		true)

	require.Equal(t,
		// 在日志档案 log1 里，服务的名称要为svc1，等级为 NOTICE，写入一笔记录为 record6
		strings.Contains(ret1[2], "svc1") &&
			strings.Contains(ret1[2], "NOTICE") &&
			strings.Contains(ret1[2], "record6"),
		true)

	// >>>>> >>>>> >>>>> >>>>> >>>>> 最后检查日志 log2 的写入结果

	// 把整个模拟双向通道的内容取出，会依照写入排列去排列，排列如下
	// [
	// [2022-01-05 16:54:36] [svc1] [debian5] [NOTICE] [1000000001] [testing.tRunner:testing.go:1410] record3 // record3 日志会写到预设的 log2 文档
	// [2022-01-05 16:54:36] [svc1] [debian5] [DEBUG]  [1000000001] [testing.tRunner:testing.go:1410] record5 // record5 日志会写到预设的 log2 文档
	// ]

	// 把整个模拟双向通道的内容取出，会依照写入排列去排列，排列如下

	// 把整个模拟双向通道的内容取出
	ret2 := ps.multi["log2"].storage.client.(*channel.Client).Read("log2")

	// 检查由模拟的双向通道中取出资料内容
	require.Equal(t,
		// 在日志档案 log2 里，服务的名称要为svc1，等级为 NOTICE，写入一笔记录为 record3
		strings.Contains(ret2[0], "svc1") &&
			strings.Contains(ret2[0], "NOTICE") &&
			strings.Contains(ret2[0], "record3"),
		true)

	// 检查由模拟的双向通道中取出资料内容
	require.Equal(t,
		// 在日志档案 log2 里，服务的名称要为svc1，等级为 DEBUG，写入一笔记录为 record6
		strings.Contains(ret2[1], "svc1") &&
			strings.Contains(ret2[1], "DEBUG") &&
			strings.Contains(ret2[1], "record5"),
		true)

	// 关闭所有模拟用的双向通道
	ps.Close()
}

// TestMultiFileLogContentCase3 为 用来测试日志分流后的写入内容 (2 个档案，2 个等级，2 个服务)
func TestMultiFileLogContentCase3(t *testing.T) {
	// >>>>> >>>>> >>>>> >>>>> >>>>> 准备和初始化日志分流物件

	// 产生日志分流物件
	ps := new(XMultiFileLog)

	// 准备日志分流设定值
	cfg := make(map[string]string)
	cfg["path"] = "/home/panhong/go/src/github.com/panhong/demo/logs"
	cfg["filename"] = "log1,log2"
	cfg["level"] = "Notice,Debug"
	cfg["service"] = "svc1,svc2"
	cfg["skip"] = "5"
	cfg["storage"] = "channel"
	cfg["chan_size"] = "10"

	// 初始化日志分流物件
	err := ps.Init(cfg)
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

	// >>>>> >>>>> >>>>> >>>>> >>>>> 最后检查日志 log2 的写入结果

	// 把整个模拟双向通道的内容取出，会依照写入排列去排列，排列如下
	// [
	// [2022-01-05 15:39:54] [svc1] [debian5] [NOTICE] [1000000001] [testing.tRunner:testing.go:1410] record1 // record1 日志会写到预设的 log1 文档
	// [2022-01-05 15:39:54] [svc1] [debian5] [NOTICE] [1000000001] [testing.tRunner:testing.go:1410] record2 // record2 日志会写到指定的 log1 文档
	// [2022-01-05 15:39:54] [svc1] [debian5] [DEBUG]  [1000000001] [testing.tRunner:testing.go:1410] record4 // 因为目前的等级为 Notice ， Debug 等服的日志会被忽略
	// [2022-01-05 15:39:54] [svc1] [debian5] [NOTICE] [1000000001] [testing.tRunner:testing.go:1410] record6 // log6 日志文档不存在，record1 日志写到 log1 文档
	// ]

	// 把整个模拟双向通道的内容取出，会依照写入排列去排列，排列如下

	// 把整个模拟双向通道的内容取出
	ret1 := ps.multi["log1"].storage.client.(*channel.Client).Read("log1")

	// 数量检查
	require.Equal(t, len(ret1), 3) // 写入五笔日志，但被忽略一笔，最后只会被写入四笔日志

	// 检查由模拟的双向通道中取出资料内容
	require.Equal(t,
		// 在日志档案 log1 里，服务的名称要为svc1，等级为 NOTICE，写入一笔记录为 record1
		strings.Contains(ret1[0], "svc1") &&
			strings.Contains(ret1[0], "NOTICE") &&
			strings.Contains(ret1[0], "record1"),
		true)

	require.Equal(t,
		// 在日志档案 log1 里，服务的名称要为svc1，等级为 NOTICE，写入一笔记录为 record2
		strings.Contains(ret1[1], "svc1") &&
			strings.Contains(ret1[1], "NOTICE") &&
			strings.Contains(ret1[1], "record2"),
		true)

	require.Equal(t,
		// 在日志档案 log1 里，服务的名称要为svc1，等级为 NOTICE，写入一笔记录为 record6
		strings.Contains(ret1[2], "svc1") &&
			strings.Contains(ret1[2], "NOTICE") &&
			strings.Contains(ret1[2], "record6"),
		true)

	// >>>>> >>>>> >>>>> >>>>> >>>>> 最后检查日志 log2 的写入结果

	// 把整个模拟双向通道的内容取出，会依照写入排列去排列，排列如下
	// [
	// [2022-01-05 16:54:36] [svc2] [debian5] [NOTICE] [1000000001] [testing.tRunner:testing.go:1410] record3 // record3 日志会写到预设的 log2 文档
	// [2022-01-05 16:54:36] [svc2] [debian5] [DEBUG]  [1000000001] [testing.tRunner:testing.go:1410] record5 // record5 日志会写到预设的 log2 文档
	// ]

	// 把整个模拟双向通道的内容取出，会依照写入排列去排列，排列如下

	// 把整个模拟双向通道的内容取出
	ret2 := ps.multi["log2"].storage.client.(*channel.Client).Read("log2")

	// 检查由模拟的双向通道中取出资料内容
	require.Equal(t,
		// 在日志档案 log2 里，服务的名称要为svc2，等级为 NOTICE，写入一笔记录为 record3
		strings.Contains(ret2[0], "svc2") &&
			strings.Contains(ret2[0], "NOTICE") &&
			strings.Contains(ret2[0], "record3"),
		true)

	// 检查由模拟的双向通道中取出资料内容
	require.Equal(t,
		// 在日志档案 log2 里，服务的名称要为svc2，等级为 DEBUG，写入一笔记录为 record6
		strings.Contains(ret2[1], "svc2") &&
			strings.Contains(ret2[1], "DEBUG") &&
			strings.Contains(ret2[1], "record5"),
		true)

	// 关闭所有模拟用的双向通道
	ps.Close()
}
