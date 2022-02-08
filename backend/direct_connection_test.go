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
	"github.com/XiaoMi/Gaea/mysql"
	"github.com/stretchr/testify/require"
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

// TestDCWithoutDB 为使用直连函式去测试数据库的连线流程，以下测试不使用 MariaDB 的服务器，只是单纯的单元测试
func TestDCWithoutDB(t *testing.T) {
	// 准备资料库的回应资料
	mysqlResponse := []uint8{
		// 资料长度
		93, 0, 0,
		// 自增序列号码
		0,
		// 以下 93 笔数据
		// 数据库的版本号
		10, 53, 46, 53, 46,
		53, 45, 49, 48, 46,
		53, 46, 49, 50, 45,
		77, 97, 114, 105, 97,
		68, 66, 45, 108, 111,
		103,
		// 数据库的版本结尾
		0,
		// 连线编号
		16, 0, 0, 0,
		// Salt
		81, 64, 43, 85, 76, 90, 97, 91,
		// filter
		0,
		// 取得功能标志
		254, 247,

		// 之后再处理
		33, 2, 0, 255, 129,
		21, 0, 0, 0, 0,
		0, 0, 15, 0, 0,
		0, 34, 53, 36, 85,
		93, 86, 117, 105, 49,
		87, 65, 125, 0, 109,
		121, 115, 113, 108, 95,
		110, 97, 116, 105, 118,
		101, 95, 112, 97, 115,
		115, 119, 111, 114, 100,
		0}

	// 测试开始
	t.Run("测试数据库连线的直連流程", func(t *testing.T) {
		// 先建立 pipe
		read, write := net.Pipe()

		// 实现缓存
		// reader := bufio.NewReaderSize(read, connBufferSize) // 用来模拟 Gaea 读取数据
		writer := bufio.NewWriter(write) // 用来模拟数据库回传数据

		// 进行等待作业
		wg := sync.WaitGroup{}
		wg.Add(1)

		// 启动数据库
		go func() {
			_, err := writer.Write(mysqlResponse) // 回传给客户端
			require.Equal(t, err, nil)
			err = writer.Flush() // 把缓存资料写进 pipe
			require.Equal(t, err, nil)
			err = write.Close() // 资料写入完成，终结连线
			require.Equal(t, err, nil)

			// 工作完成
			wg.Done()
		}()

		// 產生 Conn 物件
		var dc DirectConnection
		var mysqlConn = mysql.NewConn(read)
		dc.conn = mysqlConn
		dc.readInitialHandshake()

		// 等待中
		wg.Wait()

	})

}
