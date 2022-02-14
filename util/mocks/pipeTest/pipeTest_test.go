package pipeTest

import (
	"github.com/stretchr/testify/require"
	"testing"
)

// TestPipeTestWithoutDB 为使用直连函式去测试数据库的连线流程，以下测试不使用 MariaDB 的服务器，只是单纯的单元测试
func TestPipeTestWithoutDB(t *testing.T) {
	t.Run("此为 DC 测试的验证测试，主要是用来确认整个测试流程没有问题", func(t *testing.T) {
		// 开始模拟物件
		mockClient, mockServer := NewDcServerClient(t, TestReplyFunc) // 产生 Client 和 mockServer 模拟物件

		// 产生一开始的讯息和预期讯息
		msg0 := []uint8{0}  // 起始传送讯息
		correct := uint8(0) // 预期的正确讯息

		// 产生一连串的接收和回应的操作
		for i := 0; i < 5; i++ {
			msg1 := mockClient.SendOrReceive(msg0).Reply(mockServer) // 接收和回应
			correct++                                                // 每经过一个接收和回应的操作时，回应讯息会加1
			require.Equal(t, msg1[0], correct)
			msg0 = mockServer.SendOrReceive(msg1).Reply(mockClient) // 接收和回应
			correct++                                               // 每经过一个接收和回应的操作时，回应讯息会加1
			require.Equal(t, msg0[0], correct)
		}
	})
}
