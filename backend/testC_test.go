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

	field0 := mysql.Field{}
	field0.Data = mysql.FieldData{3, 100, 101, 102, 7, 76, 105, 98, 114, 97, 114, 121, 4, 66, 111, 111, 107, 4, 66, 111, 111, 107, 6, 66, 111, 111, 107, 73, 68, 6, 66, 111, 111, 107, 73, 68, 12, 63, 0, 11, 0, 0, 0, 3, 3, 80, 0, 0, 0}
	field0.Schema = []uint8{76, 105, 98, 114, 97, 114, 121}
	field0.Table = []uint8{66, 111, 111, 107}
	field0.OrgTable = []uint8{66, 111, 111, 107}
	field0.Name = []uint8{66, 111, 111, 107, 73, 68}
	field0.OrgName = []uint8{66, 111, 111, 107, 73, 68}
	field0.Charset = 63
	field0.ColumnLength = 11
	field0.Type = 3
	field0.Flag = 20483
	field0.Decimal = 0
	field0.DefaultValueLength = 0
	field0.DefaultValue = nil

	field1 := mysql.Field{}
	field1.Data = mysql.FieldData{3, 100, 101, 102, 7, 76, 105, 98, 114, 97, 114, 121, 4, 66, 111, 111, 107, 4, 66, 111, 111, 107, 4, 73, 115, 98, 110, 4, 73, 115, 98, 110, 12, 63, 0, 50, 0, 0, 0, 8, 1, 16, 0, 0, 0}
	field1.Schema = []uint8{76, 105, 98, 114, 97, 114, 121}
	field1.Table = []uint8{66, 111, 111, 107}
	field1.OrgTable = []uint8{66, 111, 111, 107}
	field1.Name = []uint8{73, 115, 98, 110}
	field1.OrgName = []uint8{73, 115, 98, 110}
	field1.Charset = 63
	field1.ColumnLength = 50
	field1.Type = 8
	field1.Flag = 4097
	field1.Decimal = 0
	field1.DefaultValueLength = 0
	field1.DefaultValue = nil

	field2 := mysql.Field{}
	field2.Data = mysql.FieldData{3, 100, 101, 102, 7, 76, 105, 98, 114, 97, 114, 121, 4, 66, 111, 111, 107, 4, 66, 111, 111, 107, 5, 84, 105, 116, 108, 101, 5, 84, 105, 116, 108, 101, 12, 33, 0, 44, 1, 0, 0, 253, 1, 16, 0, 0, 0}
	field2.Schema = []uint8{76, 105, 98, 114, 97, 114, 121}
	field2.Table = []uint8{66, 111, 111, 107}
	field2.OrgTable = []uint8{66, 111, 111, 107}
	field2.Name = []uint8{84, 105, 116, 108, 101}
	field2.OrgName = []uint8{84, 105, 116, 108, 101}
	field2.Charset = 33
	field2.ColumnLength = 300
	field2.Type = 253
	field2.Flag = 4097
	field2.Decimal = 0
	field2.DefaultValueLength = 0
	field2.DefaultValue = nil

	field3 := mysql.Field{}
	field3.Data = mysql.FieldData{3, 100, 101, 102, 7, 76, 105, 98, 114, 97, 114, 121, 4, 66, 111, 111, 107, 4, 66, 111, 111, 107, 6, 65, 117, 116, 104, 111, 114, 6, 65, 117, 116, 104, 111, 114, 12, 33, 0, 90, 0, 0, 0, 253, 0, 0, 0, 0, 0}
	field3.Schema = []uint8{76, 105, 98, 114, 97, 114, 121}
	field3.Table = []uint8{66, 111, 111, 107}
	field3.OrgTable = []uint8{66, 111, 111, 107}
	field3.Name = []uint8{65, 117, 116, 104, 111, 114}
	field3.OrgName = []uint8{65, 117, 116, 104, 111, 114}
	field3.Charset = 33
	field3.ColumnLength = 90
	field3.Type = 253
	field3.Flag = 0
	field3.Decimal = 0
	field3.DefaultValueLength = 0
	field3.DefaultValue = nil

	field4 := mysql.Field{}
	field4.Data = mysql.FieldData{3, 100, 101, 102, 7, 76, 105, 98, 114, 97, 114, 121, 4, 66, 111, 111, 107, 4, 66, 111, 111, 107, 7, 80, 117, 98, 108, 105, 115, 104, 7, 80, 117, 98, 108, 105, 115, 104, 12, 63, 0, 4, 0, 0, 0, 3, 0, 0, 0, 0, 0}
	field4.Schema = []uint8{76, 105, 98, 114, 97, 114, 121}
	field4.Table = []uint8{66, 111, 111, 107}
	field4.OrgTable = []uint8{66, 111, 111, 107}
	field4.Name = []uint8{80, 117, 98, 108, 105, 115, 104}
	field4.OrgName = []uint8{80, 117, 98, 108, 105, 115, 104}
	field4.Charset = 63
	field4.ColumnLength = 4
	field4.Type = 3
	field4.Flag = 0
	field4.Decimal = 0
	field4.DefaultValueLength = 0
	field4.DefaultValue = nil

	field5 := mysql.Field{}
	field5.Data = mysql.FieldData{3, 100, 101, 102, 7, 76, 105, 98, 114, 97, 114, 121, 4, 66, 111, 111, 107, 4, 66, 111, 111, 107, 8, 67, 97, 116, 101, 103, 111, 114, 121, 8, 67, 97, 116, 101, 103, 111, 114, 121, 12, 33, 0, 90, 0, 0, 0, 253, 1, 16, 0, 0, 0}
	field5.Schema = []uint8{76, 105, 98, 114, 97, 114, 121}
	field5.Table = []uint8{66, 111, 111, 107}
	field5.OrgTable = []uint8{66, 111, 111, 107}
	field5.Name = []uint8{67, 97, 116, 101, 103, 111, 114, 121}
	field5.OrgName = []uint8{67, 97, 116, 101, 103, 111, 114, 121}
	field5.Charset = 33
	field5.ColumnLength = 90
	field5.Type = 253
	field5.Flag = 4097
	field5.Decimal = 0
	field5.DefaultValueLength = 0
	field5.DefaultValue = nil

	resultset := mysql.Resultset{}
	res.Resultset = &resultset
	res.Resultset.Fields = []*mysql.Field{}
	res.Resultset.Fields = append(res.Resultset.Fields, &field0)
	res.Resultset.Fields = append(res.Resultset.Fields, &field1)
	res.Resultset.Fields = append(res.Resultset.Fields, &field2)
	res.Resultset.Fields = append(res.Resultset.Fields, &field3)
	res.Resultset.Fields = append(res.Resultset.Fields, &field4)
	res.Resultset.Fields = append(res.Resultset.Fields, &field5)

	fmt.Println(res)
}
