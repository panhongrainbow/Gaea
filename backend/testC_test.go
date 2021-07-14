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

	// 先在这里直接中断
	return

	// 向 MariaDB 数据库直接写入资料
	sql := "INSERT INTO `Library`.`Book` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (1,9781517191276,'Romance Of The Three Kingdoms','Luo Guanzhong',1522,'Historical fiction')"
	err := dc.connect()
	require.Equal(t, nil, err) // 检查连线是否有问题
	res, err := dc.Execute(sql, 1)
	require.Equal(t, nil, err)                      // 检查执行 SQL 是否有问题
	require.Equal(t, res.AffectedRows, uint64(0x1)) // 检查新增数据是否为 1 笔
}

func TestC3(t *testing.T) {
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

	// 向 MariaDB 数据库直接查询资料
	// sql := "SELECT * FROM `Library`.`Book`;"

	// 初始化單元測試程式
	dc.Mclient = new(MockDcClient)
	dc.Trans = new(MockDcClient)
	// 單元測試接管
	dc.Mclient.MarkTakeOver()

	// 进行连线
	err := dc.connect()
	require.Equal(t, nil, err) // 检查连线是否有问题

	// res, err := dc.Execute(sql, 50)
	// require.Equal(t, nil, err)                      // 检查执行 SQL 是否有问题
	// require.Equal(t, res.AffectedRows, uint64(0x0)) // 检查新增数据是否为 1 笔
}

func TestC4(t *testing.T) {
	client := new(MockDcClient)
	client.Result = make(map[uint32]mysql.Result)
	number := client.MakeResult("Library", "SELECT * FROM Book;", mysql.SelectLibrayResult())
	require.Equal(t, uint32(2125487740), number)
}
