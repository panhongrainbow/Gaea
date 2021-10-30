package backend

import (
	"context"
	"github.com/XiaoMi/Gaea/models"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// TestPooledConnect 函式 🧚 是测试 连接池 的连接
func TestPooledConnect(t *testing.T) {
	// 开启单元测试
	MarkTakeOver()

	// >>>>> >>>>> >>>>> >>>>> >>>>> 建立连接池资料
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

	// >>>>> >>>>> >>>>> >>>>> >>>>> 进行 Parse 动作
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

	// >>>>> >>>>> >>>>> >>>>> >>>>> 检查载入的设定值
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

	// >>>>> >>>>> >>>>> >>>>> >>>>> 开始写入资料到数据库
	// 使用 数据库
	err = pcM.UseDB("novel")
	require.Equal(t, err, nil)

	// 新增一笔资料
	result, err := pcM.Execute("INSERT INTO `novel`.`Book_0000` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (2,9789869442060,'Water Margin','Shi Nai an',1589,'Historical fiction')", 100)
	require.Equal(t, err1, nil)

	// 检查数据库写入结果
	require.Equal(t, err, nil)
	require.Equal(t, result.AffectedRows, uint64(0x1))
	require.Equal(t, result.InsertID, uint64(0x0))

	// 查询一笔资料
	result, err = pcS0.Execute("SELECT * FROM `novel`.`Book_0000`", 100)
	require.Equal(t, err1, nil)

	// 检查数据库读取结果
	require.Equal(t, result.Resultset.Values[0][0].(int64), int64(2))
	require.Equal(t, result.Resultset.Values[0][1].(int64), int64(9789869442060))
	require.Equal(t, result.Resultset.Values[0][2].(string), "Water Margin")

	// 删除一笔资料
	if !IsTakeOver() {
		_, err = pcM.Execute("DELETE FROM novel.Book_0000 WHERE BookID=2;", 100)
		require.Equal(t, err, nil)
	}

	// 关闭单元测试
	UnmarkTakeOver()
}

// TestPooledGetConnection 函式 🧚 是用来测试 测试获得连接池的连线
func TestPooledGetConnection(t *testing.T) {
	// 开启单元测试
	MarkTakeOver()

	// >>>>> >>>>> >>>>> >>>>> >>>>> 建立连接池资料
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

	// >>>>> >>>>> >>>>> >>>>> >>>>> 进行 Parse 动作
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

	// >>>>> >>>>> >>>>> >>>>> >>>>> 开始获得 Master 数据库连线
	pcM, err := s.GetConn(false, 0)
	require.Equal(t, err, nil)
	require.Equal(t, pcM.GetAddr(), "192.168.122.2:3309")

	// >>>>> >>>>> >>>>> >>>>> >>>>> 开始获得 Slave 数据库连线
	pcS0 := 0 // 统计 使用第一台从数据库 192.168.122.2:3310 的机率
	pcS1 := 0 // 统计 使用第二台从数据库 192.168.122.2:3311 的机率
	for i := 0; i < 10; i++ {
		pcS, err := s.GetConn(true, 0)
		require.Equal(t, err, nil)
		addr := pcS.GetAddr()
		require.Equal(t, (addr == "192.168.122.2:3310") || (addr == "192.168.122.2:3311"), true)
		if addr == "192.168.122.2:3310" {
			pcS0++
		}
		if addr == "192.168.122.2:3311" {
			pcS1++
		}
		time.Sleep(100 * time.Microsecond)
	}
	/*fmt.Println("使用第一台从数据库 192.168.122.2:3310 的机率为", pcS0*10, "%")
	fmt.Println("使用第二台从数据库 192.168.122.2:3311 的机率为", pcS1*10, "%")*/
	require.Equal(t, pcS0 > 0, true) // 要求在测试时，一定要使用到第一台从数据库 192.168.122.2:3310
	require.Equal(t, pcS1 > 0, true) // 要求在测试时，一定要使用到第二台从数据库 192.168.122.2:3310

	// 关闭单元测试
	UnmarkTakeOver()
}
