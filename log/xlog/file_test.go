package xlog

import (
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

// >>>>> >>>>> >>>>> >>>>> >>>>> 测试单一日志的写入内容

// TestFileLogContent 为 用来测试单一日志的写入内容 (包含所有的测试案例)
func TestFileLogContent(t *testing.T) {
	// >>>>> >>>>> >>>>> >>>>> >>>>> 只准备一个双向通道的测试环境

	// 准备摸拟设定值
	// cfgSuit := make(map[string]string, 1) // 有几份档案就要准备几个模拟用的双向通道，只写入单一日志档案状况下，只要准备一个双向通道
	// cfgSuit["filename"] = "log1"          // 将要用双向通道摸拟 log1.log 这一份日志档案
	// cfgSuit["chanSize"] = "5"             // 在每一个模拟用的双向通道可以接收 5 笔日志资料

	// 开始建立模拟环境
	// err := testMultiAndFileSuite(cfgSuit) // 内含初始化的操作
	// require.Equal(t, err, nil)

	// >>>>> >>>>> >>>>> >>>>> >>>>> 准备和初始化日志物件 (1 个档案，1 个等级，1 个服务，单一日志档案所有设定值几乎只能设定一个)

	// 产生日志物件
	p := new(XFileLog)

	// 准备单一日志分流设定值 (单一日志档案所有设定值几乎只能设定一个)
	cfg := make(map[string]string)
	cfg["path"] = "/home/panhong/go/src/github.com/panhong/demo/logs"
	cfg["filename"] = "log1"
	cfg["level"] = "Notice"
	cfg["service"] = "svc1"
	cfg["skip"] = "5"
	cfg["storage"] = "channel"

	// 初始化单一日志物件
	err := p.Init(cfg)
	require.Equal(t, err, nil)

	// >>>>> >>>>> >>>>> >>>>> >>>>> 开始写入日志的操作

	// 第一笔 Notice 日志
	err = p.Notice("record1") // 此日志 会 写入档案内 (V)
	require.Equal(t, err, nil)

	// 第二笔 Debug 日志
	err = p.Debug("record2") // log1 日志档案等级为 Notice，会把 Debug 的日志讯息忽略掉，所以此日志 不会 写入档案内 (X)
	require.Equal(t, err, nil)

	// >>>>> >>>>> >>>>> >>>>> >>>>> 最后检查日志的写入结果

	// 把整个模拟双向通道的内容取出
	ret := PrintMockChanMsg()

	// 数量检查
	require.Equal(t, len(ret), 1) // 写入二笔日志，有一笔写入，另一笔会被忽略，只剩一笔被正式写入

	// 细节鐱查
	// 在模拟的双向通道的内容如下
	// log1::[2021-12-23 17:44:17] [svc1] [debian5] [NOTICE] [900000001] [runtime.goexit:asm_amd64.s:1581] record1

	// 检查由模拟的双向通道中取出资料内容
	// 检查由模拟的双向通道中取出资料内容
	require.Equal(t,
		// 在日志档案 log1 里，服务的名称要为svc1，等级为 NOTICE，写入一笔记录为 record1
		strings.Contains(ret[0], "log1") &&
			strings.Contains(ret[0], "svc1") &&
			strings.Contains(ret[0], "NOTICE") &&
			strings.Contains(ret[0], "record1"),
		true)
}
