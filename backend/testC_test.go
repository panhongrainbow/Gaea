package backend

import (
	"fmt"
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

func TestC2(t *testing.T) {
	res := mysql.Result{}
	res.Status = 34
	res.InsertID = 0
	res.AffectedRows = 0

	field := mysql.Field{}
	field.Data = mysql.FieldData{3, 100, 101, 102, 7, 76, 105, 98, 114, 97, 114, 121, 4, 66, 111, 111, 107, 4, 66, 111, 111, 107, 6, 66, 111, 111, 107, 73, 68, 6, 66, 111, 111, 107, 73, 68, 12, 63, 0, 11, 0, 0, 0, 3, 3, 80, 0, 0, 0}
	field.Schema = []uint8{76, 105, 98, 114, 97, 114, 121}

	resultset := mysql.Resultset{}
	res.Resultset = &resultset
	res.Resultset.Fields = []*mysql.Field{}
	res.Resultset.Fields = append(res.Resultset.Fields, &field)

	fmt.Println(res)
}
