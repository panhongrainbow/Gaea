package backend

import (
	"testing"

	"github.com/XiaoMi/Gaea/mysql"
	"github.com/stretchr/testify/require"
)

func TestC1(t *testing.T) {
	// 建立直接连线
	var dc DirectConnection
	dc.addr = "192.168.1.2:3350"
	dc.user = "docker"
	dc.password = "12345"
	dc.db = "Library"
	dc.capability = 0
	dc.sessionVariables = mysql.NewSessionVariables()
	dc.status = 0
	dc.collation = 0
	dc.charset = "utf8mb4"
	dc.salt = nil
	dc.defaultCollation = 0
	dc.defaultCharset = "utf8mb4"
	dc.pkgErr = nil
	dc.closed.Set(false)

	// 以下会直接连线到实体数据库，先在这里中断
	return

	// 向 MariaDB 数据库直接写入资料
	sql := "INSERT INTO `Library`.`Book` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (1,9781517191276,'Romance Of The Three Kingdoms','Luo Guanzhong',1522,'Historical fiction')"
	err := dc.connect()
	require.Equal(t, nil, err) // 检查连线是否有问题
	res, err := dc.Execute(sql, 1)
	require.Equal(t, nil, err)                      // 检查执行 SQL 是否有问题
	require.Equal(t, res.AffectedRows, uint64(0x1)) // 检查新增数据是否为 1 笔
}

// TestC 主要是用于测试数据库直连
func TestC2(t *testing.T) {
	// 建立直接连线
	var dc DirectConnection
	dc.addr = "192.168.122.2:3306"
	dc.user = "panhong"
	dc.password = "12345"
	dc.db = "Library"
	dc.capability = 0
	dc.sessionVariables = mysql.NewSessionVariables()
	dc.status = 0
	dc.collation = 0
	dc.charset = "utf8mb4"
	dc.salt = nil
	dc.defaultCollation = 0
	dc.defaultCharset = "utf8mb4"
	dc.pkgErr = nil
	dc.closed.Set(false)

	// 初始化單元測試程式
	dc.MockDC = new(MockDcClient)
	// dc.Trans = new(MockDcClient)

	// 單元測試接管
	MarkTakeOver()

	// 进行连线
	err := dc.connect()
	require.Equal(t, nil, err) // 检查连线是否有问题

	// 准备 SQL 查询字串
	sql := "SELECT * FROM `Library`.`Book`;"

	// 向 MariaDB 数据库直接查询资料
	res, err := dc.Execute(sql, 50)
	require.Equal(t, nil, err)                      // 检查执行 SQL 是否有问题
	require.Equal(t, res.AffectedRows, uint64(0x0)) // 因为只是查询，所以数据库不会有资料被修改

	// 检查数据库回传第 1 本书的资料
	require.Equal(t, res.Values[0][0].(int64), int64(1))
	require.Equal(t, res.Values[0][1].(int64), int64(9781517191276))
	require.Equal(t, res.Values[0][2].(string), "Romance Of The Three Kingdoms")

	// 检查数据库回传第 28 本书的资料
	require.Equal(t, res.Values[28][0].(int64), int64(29))
	require.Equal(t, res.Values[28][1].(int64), int64(9789866318603))
	require.Equal(t, res.Values[28][2].(string), "A History Of Floral Treasures")
}

func TestC3(t *testing.T) {
	client := new(MockDcClient)
	client.MockResult = make(map[uint32]mysql.Result)
	number := client.MakeMockResult("Library", "SELECT * FROM Book;", *mysql.SelectLibrayResult())
	require.Equal(t, uint32(2125487740), number)
}
