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
	"bytes"
	"github.com/stretchr/testify/require"
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

// >>>>> >>>>> >>>>> >>>>> >>>>> 以下为新增的直连 DC 单元测试，目的是用来了解直连 DC 的基本动作

// TestDc 函式 🧚 是用来測試所有的直連 DC 的基本動作
func TestDc(t *testing.T) {
	// 直连 DC 的单元测试是否能正常启动
	TestDcTakeOver(t)
	// 直连 DC 的新建连线
	TestDcNewDirectConnection(t)
	// 重新建立直连 DC 连线
	TestDcReCreateConnection(t)
}

// TestDcTakeOver 函式 🧚 是用来测试直连 DC 的单元测试是否能正常启动
func TestDcTakeOver(t *testing.T) {
	// 启动单元测试的开关
	MarkTakeOver()

	// 确认单元测试的开关是否正常开启
	require.Equal(t, IsTakeOver(), true)

	// 关闭单元测试的开关
	UnmarkTakeOver()

	// 确认单元测试的开关是否正常关闭
	require.Equal(t, IsTakeOver(), false)
}

// TestDcNewDirectConnection 函式 🧚 是用来测试直连 DC 的新建连线
func TestDcNewDirectConnection(t *testing.T) {
	// 启动单元测试的开关
	MarkTakeOver()

	// 直接在這裡建立新的直连 DC 連線
	dcConn, err := NewDirectConnection(
		"192.168.122.2:3309",
		"panhong",
		"12345",
		"novel",
		"utf8mb4",
		46,
	)

	// 檢查測試直連 DC 的連線是否成功建立
	require.Equal(t, err, nil)

	// 用于测试直连 DC 的所有基本动作
	err = dcConn.Ping()

	// 检查连线测试是否正常
	require.Equal(t, err, nil)

	// 关闭单元测试的开关
	UnmarkTakeOver()
}

// TestDcReCreateConnection 函式 🧚 是用来测试重新建立直连 DC 连线
/*
參考以下程式碼，在重建直連 DC 連線時，會先關閉連線 Close，再重新建立連線
func (pc *pooledConnectImpl) Reconnect() error {
	pc.directConnection.Close() // 第 1 步，先关闭连线
	newConn, err := NewDirectConnection(pc.pool.addr, pc.pool.user, pc.pool.password, pc.pool.db, pc.pool.charset, pc.pool.collationID) // 第 2 步，再重新建立建立连线
	if err != nil {
		return err
	}
	pc.directConnection = newConn
	return nil
}
*/
func TestDcReCreateConnection(t *testing.T) {
	// 启动单元测试的开关
	MarkTakeOver()

	// 直接在這裡建立新的直连 DC 連線
	dcConn, err := NewDirectConnection(
		"192.168.122.2:3309",
		"panhong",
		"12345",
		"novel",
		"utf8mb4",
		46,
	)

	// 檢查測試直連 DC 的連線是否成功建立
	require.Equal(t, err, nil)

	// 第 1 步，先关闭连线
	dcConn.Close()

	// 檢查連線是否已經關閉，應要為 True
	require.Equal(t, dcConn.IsClosed(), true)
	require.Equal(t, dcConn.closed.Get(), true)

	// 第 2 步，再重新建立建立连线
	dcConn, err = NewDirectConnection(
		"192.168.122.2:3309",
		"panhong",
		"12345",
		"novel",
		"utf8mb4",
		46,
	)

	// 檢查測試直連 DC 的連線是否成功建立
	require.Equal(t, err, nil)

	// 檢查連線是否已經關閉，應要為 False
	require.Equal(t, dcConn.IsClosed(), false)
	require.Equal(t, dcConn.closed.Get(), false)

	// 关闭单元测试的开关
	UnmarkTakeOver()
}
