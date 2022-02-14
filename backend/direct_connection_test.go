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
package backend

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/XiaoMi/Gaea/mysql"
	"github.com/stretchr/testify/require"
	"io"
	"net"
	"sync"
	"testing"
)

func TestAppendSetVariable(t *testing.T) {
	var buf bytes.Buffer
	appendSetVariable(&buf, "charset", "utf8")
	t.Log(buf.String())
	appendSetVariable(&buf, "autocommit", 1)
	t.Log(buf.String())
	appendSetVariableToDefault(&buf, "sql_mode")
	t.Log(buf.String())
}

func TestAppendSetVariable2(t *testing.T) {
	var buf bytes.Buffer
	appendSetCharset(&buf, "utf8", "utf8_general_ci")
	t.Log(buf.String())
	appendSetVariable(&buf, "autocommit", 1)
	t.Log(buf.String())
	appendSetVariableToDefault(&buf, "sql_mode")
	t.Log(buf.String())
}

var (
	// 准备资料库的回应资料
	mysqlInitHandShakeResponse = []uint8{
		// 资料长度
		93, 0, 0,
		// 自增序列号码
		0,
		// 以下 93 笔数据
		// 数据库的版本号 version
		10, 53, 46, 53, 46,
		53, 45, 49, 48, 46,
		53, 46, 49, 50, 45,
		77, 97, 114, 105, 97,
		68, 66, 45, 108, 111,
		103,
		// 数据库的版本结尾
		0,
		// 连线编号 connection id
		16, 0, 0, 0,
		// Salt
		81, 64, 43, 85, 76, 90, 97, 91,
		// filter
		0,
		// 取得功能标志 capability
		254, 247,
		// 數據庫編碼 charset
		33, // 可以用 SHOW CHARACTER SET LIKE 'utf8'; 查询
		// 服务器状态，在 Gaea/mysql/constants.go 的 Server information
		2, 0,
		// 延伸的功能标志 capability
		255, 129,
		// Auth 资料和保留值
		21, 0, 0, 0, 0, 0, 0, 15, 0, 0, 0,
		// 延伸的 Salt
		34, 53, 36, 85,
		93, 86, 117, 105, 49,
		87, 65, 125,
		// 其他未用到的资料
		0, 109, 121, 115, 113, 108, 95, 110, 97, 116,
		105, 118, 101, 95, 112, 97, 115, 115, 119, 111,
		114, 100, 0,
	}
)

// replyFuncType 回应函式的型态
type replyFuncType func([]uint8) []uint8

// testReplyFunc 在这里会处理常接收到什么讯息，要将下来跟着回应什么讯息
// 目前此函式只是在测试验证流程，回应讯息为接收讯息加 1
//     比如 当接收值为 1，就会回传值为 2 给对方
//     比如 当接收值为 2，就会回传值为 3 给对方
func testReplyFunc(data []uint8) []uint8 {
	return []uint8{data[0] + 1} // 回应讯息为接收讯息加 1
}

// dcMocker 用来模拟数据库服务器的读取和回应的物件
type dcMocker struct {
	t         *testing.T      // 单元测试的物件
	bufReader *bufio.Reader   // 服务器的读取
	bufWriter *bufio.Writer   // 服务器的回应
	connRead  net.Conn        // pipe 的读取连线
	connWrite net.Conn        // pipe 的写入连线
	wg        *sync.WaitGroup // 流程的操作边界
	replyFunc replyFuncType   // 设定相对应的回应函式
	err       error           // 错误
}

// newDcServerClient 产生直连 DC 模拟双方，包含客户端和服务端
func newDcServerClient(t *testing.T, reply replyFuncType) (mockClient *dcMocker, mockServer *dcMocker) {
	// 先产生两组 Pipe
	read0, write0 := net.Pipe() // 第一组 Pipe
	read1, write1 := net.Pipe() // 第二组 Pipe

	// 产生客户端和服务端双方，分别为 mockClient 和 mockServer
	mockClient = newDcMocker(t, read0, write1, reply)
	mockServer = newDcMocker(t, read1, write0, reply)

	// 结束
	return
}

// newDcMocker 产生新的 dc 模拟物件
func newDcMocker(t *testing.T, connRead, connWrite net.Conn, reply replyFuncType) *dcMocker {
	return &dcMocker{
		t:         t,                          // 单元测试的物件
		bufReader: bufio.NewReader(connRead),  // 服务器的读取 (实现缓存)
		bufWriter: bufio.NewWriter(connWrite), // 服务器的回应 (实现缓存)
		connRead:  connRead,                   // pipe 的读取连线
		connWrite: connWrite,                  // pipe 的写入连线
		wg:        &sync.WaitGroup{},          // 流程的操作边界
		replyFunc: reply,                      // 回应函式
	}
}

// resetDcMockers 为重置单一连线方向的直连 dc 模拟物件
func (dcM *dcMocker) resetDcMockers(otherSide *dcMocker) error {
	// 重新建立全新两组 Pipe
	newRead, newWrite := net.Pipe() // 第一组 Pipe

	// 单方向的状况为 dcM 写入 Pipe，otherSide 读取 Pipe

	// 先重置 发送讯息的那一方 部份
	dcM.bufWriter = bufio.NewWriter(newWrite) // 服务器的回应 (实现缓存)
	dcM.connWrite = newWrite                  // pipe 的写入连线
	dcM.wg = &sync.WaitGroup{}                // 流程的操作边界

	// 先重置 mockServer 部份
	otherSide.bufReader = bufio.NewReader(newRead) // 服务器的读取 (实现缓存)
	otherSide.connRead = newRead                   // pipe 的读取连线

	// 正常回传
	return nil
}

// sendOrReceive 为直连 dc 用来模拟接收或传入讯息
func (dcM *dcMocker) sendOrReceive(data []uint8) *dcMocker {
	// dc 模拟开始
	dcM.wg.Add(1) // 只要等待直到确认资料有写入 pipe

	// 在这里执行 1传送讯息 或者是 2接收讯息
	go func() {
		// 执行写入工作
		_, err := dcM.bufWriter.Write(data) // 写入资料到 pipe
		err = dcM.bufWriter.Flush()         // 把缓存资料写进 pipe
		require.Equal(dcM.t, err, nil)
		err = dcM.connWrite.Close() // 资料写入完成，终结连线
		require.Equal(dcM.t, err, nil)

		// 写入工作完成
		dcM.wg.Done()
	}()

	// 重复使用物件
	return dcM
}

// reply 为直连 dc 用来模拟 dc 回应数据
func (dcM *dcMocker) reply(otherSide *dcMocker) (msg []uint8) {
	// 读取传送过来的讯息
	b, _, err := otherSide.bufReader.ReadLine() // 由另一方接收传来的讯息
	require.Equal(dcM.t, err, nil)

	// 等待和确认资料已经写入 pipe
	dcM.wg.Wait()

	// 重置模拟物件
	err = dcM.resetDcMockers(otherSide)
	require.Equal(dcM.t, err, nil)

	// 回传回应讯息
	if otherSide.replyFunc != nil {
		msg = otherSide.replyFunc(b)
	}

	// 结束
	return
}

// waitAndReset 为直连 dc 用来等待在 Pipe 的整个数据读写操作完成
func (dcM *dcMocker) waitAndReset(otherSide *dcMocker) error {
	// 先等待整个数据读写操作完成
	dcM.wg.Wait()

	// 单方向完成 Pipe 的连线重置
	err := dcM.resetDcMockers(otherSide)
	require.Equal(dcM.t, err, nil)

	// 正确回传
	return nil
}

// TestDCWithoutDB 为使用直连函式去测试数据库的连线流程，以下测试不使用 MariaDB 的服务器，只是单纯的单元测试
func TestDCWithoutDB(t *testing.T) {
	t.Run("此为 DC 测试的验证测试，主要是用来确认整个测试流程没有问题", func(t *testing.T) {
		// 开始模拟物件
		mockClient, mockServer := newDcServerClient(t, testReplyFunc) // 产生 Client 和 mockServer 模拟物件

		// 产生一开始的讯息和预期讯息
		msg0 := []uint8{0}  // 起始传送讯息
		correct := uint8(0) // 预期的正确讯息

		// 产生一连串的接收和回应的操作
		for i := 0; i < 5; i++ {
			msg1 := mockClient.sendOrReceive(msg0).reply(mockServer) // 接收和回应
			correct++                                                // 每经过一个接收和回应的操作时，回应讯息会加1
			require.Equal(t, msg1[0], correct)
			msg0 = mockServer.sendOrReceive(msg1).reply(mockClient) // 接收和回应
			correct++                                               // 每经过一个接收和回应的操作时，回应讯息会加1
			require.Equal(t, msg0[0], correct)
		}
	})
	// 开始正式测试
	t.Run("测试数据库连线的实际直連流程", func(t *testing.T) {
		// 开始模拟
		mockGaea, mockMariaDB := newDcServerClient(t, testReplyFunc) // 产生 Gaea 和 mockServer 模拟物件
		mockGaea.sendOrReceive(mysqlInitHandShakeResponse)           // 模拟数据库开始交握

		// 產生 Mysql dc 直連物件 (用以下内容取代 reply() 函式 !)
		var dc DirectConnection
		var mysqlConn = mysql.NewConn(mockMariaDB.connRead)
		dc.conn = mysqlConn
		err := dc.readInitialHandshake()
		require.Equal(t, err, nil)

		// 等待和确认资料已经写入 pipe 并单方向重置模拟物件
		err = mockGaea.waitAndReset(mockMariaDB)
		require.Equal(t, err, nil)

		// 开始计算

		/* 功能标志 capability 的计算
		先把所有的功能标志 capability 的数据收集起来，包含延伸部份
		数值分别为 254, 247, 255, 129
		并反向排列
		数值分别为 129, 255, 247, 254
		全部 十进制 转成 二进制
		254 的二进制为 1111 1110
		247 的二进制为 1111 0111
		255 的二进制为 1111 1111
		129 的二进制为 1000 0001
		把全部二进制的数值合并
		二进制数值分别为 1000 0001 1111 1111 1111 0111 1111 1110 (转成十进制数值为 2181036030)
		再用文档 https://mariadb.com/kb/en/connection/ 进行对照
		比如，功能标志 capability 的第一个值为 0，意思为 CLIENT_MYSQL 值为 0，代表是由服务器发出的讯息 */

		/* 连线编号 connection id 的计算
		先把所有的连线编号 connection id 的数据收集起来，包含延伸部份
		数值分别为 16, 0, 0, 0
		并反向排列
		数值分别为 0, 0, 0, 16
		全部 十进制 转成 二进制
		  0 的二进制为 0000 0000
		 16 的二进制为 0001 0000
		把全部二进制的数值合并
		二进制数值分别为 0000 0000 0001 0000 (转成十进制数值为 16) */

		// 先把所有 Salt 的数据收集起来，包含延伸部份
		// 数值分别为 81,64,43,85,76,90,97,91,34,53,36,85,93,86,117,105,49,87,65,125

		/* 服务器状态 status 的计算
		先把所有的服务器状态 的数据收集起来，包含延伸部份
		数值分别为 2, 0
		并反向排列
		数值分别为 0, 2
		全部 十进制 转成 二进制
		2 的二进制为 0000 0010
		0 的二进制为 0000 0000
		把全部二进制的数值合并
		二进制数值分别为 0000 0000 0000 0010 (转成十进制数值为 2)
		再用代码 Gaea/mysql/constants.go 里的 Server information 进行对照
		功能标志 capability 的第一个值为 0，意思为 CLIENT_MYSQL 值为 0，代表是由服务器发出的讯息 */

		// 计算后的检查
		require.Equal(t, dc.capability, uint32(2181036030))                                                                   // 检查功能标志 capability
		require.Equal(t, dc.conn.ConnectionID, uint32(16))                                                                    // 检查连线编号 connection id
		require.Equal(t, dc.salt, []uint8{81, 64, 43, 85, 76, 90, 97, 91, 34, 53, 36, 85, 93, 86, 117, 105, 49, 87, 65, 125}) // 检查 Salt
		require.Equal(t, dc.status, mysql.ServerStatusAutocommit)                                                             // 检查服务器状态

		// 以下未完成之后再处理
		mysqlConn = mysql.NewConn(mockGaea.connRead)
		mysqlConn.TestWriter(mockMariaDB.connWrite)
		dc.conn = mysqlConn
		dc.conn.SetConnectionID(uint32(16))

		go func() {
			_ = dc.writeHandshakeResponse41()
			dc.conn.Flush()
			mockMariaDB.connWrite.Close()
		}()

		var data [20]byte
		_, err = io.ReadFull(mockGaea.connRead, data[:])
		fmt.Println("data", data)
	})
}
