package backend

import (
	"context"
	"fmt"
	"github.com/XiaoMi/Gaea/models"
	"github.com/stretchr/testify/require"
	"testing"
)

// TestPooledConnect 函式 🧚 测试 是用测试 连接池 的连接
func TestPooledConnect(t *testing.T) {
	// 开启单元测试
	MarkTakeOver()

	// 载入设定档
	s := new(Slice)
	s.Cfg = models.Slice{
		Name:     "slice-0",
		UserName: "panhong",
		Password: "12345",
		Master:   "192.168.122.2:3309",
		Slaves: []string{
			"192.168.122.2:3310",
			"192.168.122.2:3311",
		},
		Capacity:    12,
		MaxCapacity: 24,
		IdleTimeout: 60,
	}
	s.charset = "utf8mb4"
	s.collationID = 46

	// 先 Parse Master
	err := s.ParseMaster(s.Cfg.Master)
	require.Equal(t, err, nil)
	require.Equal(t, s.Master.Capacity(), int64(12))

	// 检查 Parse Master 的结果
	require.Equal(t, s.Master.Addr(), "192.168.122.2:3309")

	// 再 Parse Slave
	err = s.ParseSlave(s.Cfg.Slaves)
	require.Equal(t, err, nil)

	// 检查 Parse Slave 的结果
	require.Equal(t, s.Slave[0].Addr(), "192.168.122.2:3310")
	require.Equal(t, s.Slave[1].Addr(), "192.168.122.2:3311")

	// 再 Parse StatisticSlave
	err = s.ParseStatisticSlave(s.Cfg.StatisticSlaves)

	// 检查 Parse Slave 的平衡器 的结果
	require.Equal(t, s.slaveBalancer.total, 2)
	require.Equal(t, s.slaveBalancer.lastIndex, 0)
	require.Equal(t, s.slaveBalancer.roundRobinQ, []int{0, 1})

	// 检查 Parse StatisticSlave 的结果
	require.Equal(t, s.StatisticSlave, []ConnectionPool(nil))

	// 建立 Ctx
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 建立 Master Pool 的 Connection
	pcM, err0 := s.Master.Get(ctx)
	require.Equal(t, err0, nil)

	// 检查 Pool Connection
	require.Equal(t, pcM.GetAddr(), "192.168.122.2:3309")

	// 建立 Slave0 Pool 的 Connection
	pcS0, err1 := s.Slave[0].Get(ctx)
	require.Equal(t, err1, nil)

	// 检查 Slave0 Pool Connection
	require.Equal(t, pcS0.GetAddr(), "192.168.122.2:3310")

	// 建立 Slave1 Pool 的 Connection
	pcS1, err2 := s.Slave[1].Get(ctx)
	require.Equal(t, err2, nil)

	// 检查 Slave1 Pool Connection
	require.Equal(t, pcS1.GetAddr(), "192.168.122.2:3311")

	// pcM.Close()

	fmt.Println(pcM.IsClosed())

	err3 := pcM.Reconnect()
	fmt.Println(err3)

	pcM.Recycle()

	// 关闭单元测试
	UnmarkTakeOver()
}
