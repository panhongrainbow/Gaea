package xlog

import (
	"github.com/stretchr/testify/require"
	"testing"
)

// TestPrintMockChanMsg 为 测试整个模拟环境的流程，从初始化到印出所有双向通道的讯息
func TestPrintMockChanMsg(t *testing.T) {
	// 准备设定文件
	config := make(map[string]string, 1)
	config["filename"] = "log1,log2"
	config["chan_size"] = "3"

	// 初始化模拟环境
	mf := new(MockMultiXLogFile)
	err := mf.Init(config) // 在这里进行初始化模拟文件
	require.Equal(t, err, nil)

	// 进行重新开启档案
	// 表面上是重新开启档案，但实际上是检查双向通道是否有全部被开启
	err = mf.Open()
	require.Equal(t, err, nil)
	require.Equal(t, len(mf.mockFile), 2) // 因为 filename 有指定两个设定档 log1 和 log2，所以要准备两个双向通道去模拟档案

	// 检查所有的双向通道的状况
	for fileName, mockChan := range mf.mockFile {
		require.Equal(t, (fileName == "log1") || (fileName == "log2"), true)
		require.NotEqual(t, mockChan, nil) // 所有的通道不能为空，not equal
	}

	// 设定全域的双向通道
	SetGlobalMockMultiXLogFile(mf)

	// 先传送讯息到所有的双向通道
	mf.mockFile["log1"] <- "empty1"
	mf.mockFile["log2"] <- "empty2"

	// 向所有的双向通道取值
	msg1 := <-mf.mockFile["log1"]
	require.Equal(t, msg1, "empty1")
	msg2 := <-mf.mockFile["log2"]
	require.Equal(t, msg2, "empty2")

	// 再传送讯息到所有的双向通道
	mf.mockFile["log1"] <- "empty3"
	mf.mockFile["log2"] <- "empty4"
	mf.mockFile["log2"] <- "empty5"

	// 把整个模拟双向通道的内容取出
	ret := mf.PrintMockChanMsg()

	// 检查由模拟的双向通道中取出资料内容
	require.Equal(t, ret[0], "log1::empty3")
	require.Equal(t, ret[1], "log2::empty4")
	require.Equal(t, ret[2], "log2::empty5")
}
