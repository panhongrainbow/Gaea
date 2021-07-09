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

	res.FieldNames = make(map[string]int)
	res.FieldNames["BookID"] = 0
	res.FieldNames["Isbn"] = 1
	res.FieldNames["Title"] = 2
	res.FieldNames["Author"] = 3
	res.FieldNames["Publish"] = 4
	res.FieldNames["Category"] = 5

	// https://segmentfault.com/q/1010000000198391
	res.Values = make([][]interface{}, 18)
	res.Values[0] = make([]interface{}, 6)
	res.Values[0][0] = 1
	res.Values[0][1] = 9781517191276
	res.Values[0][2] = "Romance Of The Three Kingdoms"
	res.Values[0][3] = "Luo Guanzhong"
	res.Values[0][4] = 1522
	res.Values[0][5] = "Historical fiction"

	res.Values[1] = make([]interface{}, 6)
	res.Values[1][0] = 2
	res.Values[1][1] = 9789869442060
	res.Values[1][2] = "Water Margin"
	res.Values[1][3] = "Shi Nai an"
	res.Values[1][4] = 1589
	res.Values[1][5] = "Historical fiction"

	res.Values[2] = make([]interface{}, 6)
	res.Values[2][0] = 3
	res.Values[2][1] = 9789575709518
	res.Values[2][2] = "Journey To The West"
	res.Values[2][3] = "Wu Cheng en"
	res.Values[2][4] = 1592
	res.Values[2][5] = "Gods And Demons Fiction"

	res.Values[3] = make([]interface{}, 6)
	res.Values[3][0] = 4
	res.Values[3][1] = 9789865975364
	res.Values[3][2] = "Dream Of The Red Chamber"
	res.Values[3][3] = "Cao Xueqin"
	res.Values[3][4] = 1791
	res.Values[3][5] = "Family Saga"

	res.Values[4] = make([]interface{}, 6)
	res.Values[4][0] = 5
	res.Values[4][1] = 9780804847773
	res.Values[4][2] = "Jin Ping Mei"
	res.Values[4][3] = "Lanling Xiaoxiao Sheng"
	res.Values[4][4] = 1610
	res.Values[4][5] = "Family Saga"

	res.Values[5] = make([]interface{}, 6)
	res.Values[5][0] = 6
	res.Values[5][1] = 9780835124072
	res.Values[5][2] = "Rulin Waishi"
	res.Values[5][3] = "Wu Jingzi"
	res.Values[5][4] = 1750
	res.Values[5][5] = "Unofficial History"

	res.Values[6] = make([]interface{}, 6)
	res.Values[6][0] = 7
	res.Values[6][1] = 9787101064100
	res.Values[6][2] = "Amazing Tales First Series"
	res.Values[6][3] = "Ling Mengchu"
	res.Values[6][4] = 1628
	res.Values[6][5] = "Perspective"

	res.Values[7] = make([]interface{}, 6)
	res.Values[7][0] = 8
	res.Values[7][1] = 9789571447278
	res.Values[7][2] = "Amazing Tales Second Series"
	res.Values[7][3] = "Ling Mengchu"
	res.Values[7][4] = 1628
	res.Values[7][5] = "Perspective"

	res.Values[8] = make([]interface{}, 6)
	res.Values[8][0] = 9
	res.Values[8][1] = 9789861273129
	res.Values[8][2] = "Investiture Of The Gods"
	res.Values[8][3] = "Lu Xixing"
	res.Values[8][4] = 1605
	res.Values[8][5] = "Mythology"

	res.Values[9] = make([]interface{}, 6)
	res.Values[9][0] = 10
	res.Values[9][1] = 9787540251499
	res.Values[9][2] = "Flowers In The Mirror"
	res.Values[9][3] = "Li Ruzhen"
	res.Values[9][4] = 1827
	res.Values[9][5] = "Fantasy Stories"

	res.Values[10] = make([]interface{}, 6)
	res.Values[10][0] = 11
	res.Values[10][1] = 9787508535296
	res.Values[10][2] = "Stories Old And New"
	res.Values[10][3] = "Feng Menglong"
	res.Values[10][4] = 1620
	res.Values[10][5] = "Perspective"

	res.Values[11] = make([]interface{}, 6)
	res.Values[11][0] = 12
	res.Values[11][1] = 9787101097559
	res.Values[11][2] = "General Yue Fei"
	res.Values[11][3] = "Qian Cai"
	res.Values[11][4] = 1735
	res.Values[11][5] = "History"

	res.Values[12] = make([]interface{}, 6)
	res.Values[12][0] = 13
	res.Values[12][1] = 9789863381037
	res.Values[12][2] = "The Generals Of The Yang Family"
	res.Values[12][3] = "Qi Zhonglan"
	res.Values[12][4] = 0
	res.Values[12][5] = "History"

	res.Values[13] = make([]interface{}, 6)
	res.Values[13][0] = 14
	res.Values[13][1] = 9789865700027
	res.Values[13][2] = "Romance Of Sui And Tang Dynasties"
	res.Values[13][3] = "Chen Ruheng"
	res.Values[13][4] = 1989
	res.Values[13][5] = "History"

	res.Values[14] = make([]interface{}, 6)
	res.Values[14][0] = 15
	res.Values[14][1] = 9789575709242
	res.Values[14][2] = "The Seven Heroes And Five Gallants"
	res.Values[14][3] = "Shi Yukun"
	res.Values[14][4] = 1879
	res.Values[14][5] = "History"

	res.Values[15] = make([]interface{}, 6)
	res.Values[15][0] = 16
	res.Values[15][1] = 9789575709242
	res.Values[15][2] = "A Collection Of Shi"
	res.Values[15][3] = "Anonymous"
	res.Values[15][4] = 1850
	res.Values[15][5] = "History"

	res.Values[16] = make([]interface{}, 6)
	res.Values[16][0] = 17
	res.Values[16][1] = 9787533303396
	res.Values[16][2] = "Dream Of The Green Chamber"
	res.Values[16][3] = "Yuda"
	res.Values[16][4] = 1878
	res.Values[16][5] = "Family Saga"

	res.Values[17] = make([]interface{}, 6)
	res.Values[17][0] = 18
	res.Values[17][1] = 9787510434341
	res.Values[17][2] = "Lamp In The Side Street"
	res.Values[17][3] = "Li Luyuan"
	res.Values[17][4] = 1790
	res.Values[17][5] = "Unofficial History"

	res.RowDatas = make([]mysql.RowData, 18)
	res.RowDatas[0] = []uint8{1, 49, 13, 57, 55, 56, 49, 53, 49, 55, 49, 57, 49, 50, 55, 54, 29, 82, 111, 109, 97, 110, 99, 101, 32, 79, 102, 32, 84, 104, 101, 32, 84, 104, 114, 101, 101, 32, 75, 105, 110, 103, 100, 111, 109, 115, 13, 76, 117, 111, 32, 71, 117, 97, 110, 122, 104, 111, 110, 103, 4, 49, 53, 50, 50, 18, 72, 105, 115, 116, 111, 114, 105, 99, 97, 108, 32, 102, 105, 99, 116, 105, 111, 110}
	res.RowDatas[1] = []uint8{1, 50, 13, 57, 55, 56, 57, 56, 54, 57, 52, 52, 50, 48, 54, 48, 12, 87, 97, 116, 101, 114, 32, 77, 97, 114, 103, 105, 110, 10, 83, 104, 105, 32, 78, 97, 105, 32, 97, 110, 4, 49, 53, 56, 57, 18, 72, 105, 115, 116, 111, 114, 105, 99, 97, 108, 32, 102, 105, 99, 116, 105, 111, 110}
	res.RowDatas[2] = []uint8{1, 51, 13, 57, 55, 56, 57, 53, 55, 53, 55, 48, 57, 53, 49, 56, 19, 74, 111, 117, 114, 110, 101, 121, 32, 84, 111, 32, 84, 104, 101, 32, 87, 101, 115, 116, 11, 87, 117, 32, 67, 104, 101, 110, 103, 32, 101, 110, 4, 49, 53, 57, 50, 23, 71, 111, 100, 115, 32, 65, 110, 100, 32, 68, 101, 109, 111, 110, 115, 32, 70, 105, 99, 116, 105, 111, 110}
	res.RowDatas[3] = []uint8{1, 52, 13, 57, 55, 56, 57, 56, 54, 53, 57, 55, 53, 51, 54, 52, 24, 68, 114, 101, 97, 109, 32, 79, 102, 32, 84, 104, 101, 32, 82, 101, 100, 32, 67, 104, 97, 109, 98, 101, 114, 10, 67, 97, 111, 32, 88, 117, 101, 113, 105, 110, 4, 49, 55, 57, 49, 11, 70, 97, 109, 105, 108, 121, 32, 83, 97, 103, 97}
	res.RowDatas[4] = []uint8{1, 53, 13, 57, 55, 56, 48, 56, 48, 52, 56, 52, 55, 55, 55, 51, 12, 74, 105, 110, 32, 80, 105, 110, 103, 32, 77, 101, 105, 22, 76, 97, 110, 108, 105, 110, 103, 32, 88, 105, 97, 111, 120, 105, 97, 111, 32, 83, 104, 101, 110, 103, 4, 49, 54, 49, 48, 11, 70, 97, 109, 105, 108, 121, 32, 76, 105, 102, 101}
	res.RowDatas[5] = []uint8{1, 54, 13, 57, 55, 56, 48, 56, 51, 53, 49, 50, 52, 48, 55, 50, 12, 82, 117, 108, 105, 110, 32, 87, 97, 105, 115, 104, 105, 9, 87, 117, 32, 74, 105, 110, 103, 122, 105, 4, 49, 55, 53, 48, 18, 85, 110, 111, 102, 102, 105, 99, 105, 97, 108, 32, 72, 105, 115, 116, 111, 114, 121}
	res.RowDatas[6] = []uint8{1, 55, 13, 57, 55, 56, 55, 49, 48, 49, 48, 54, 52, 49, 48, 48, 26, 65, 109, 97, 122, 105, 110, 103, 32, 84, 97, 108, 101, 115, 32, 70, 105, 114, 115, 116, 32, 83, 101, 114, 105, 101, 115, 12, 76, 105, 110, 103, 32, 77, 101, 110, 103, 99, 104, 117, 4, 49, 54, 50, 56, 11, 80, 101, 114, 115, 112, 101, 99, 116, 105, 118, 101}
	res.RowDatas[7] = []uint8{1, 56, 13, 57, 55, 56, 57, 53, 55, 49, 52, 52, 55, 50, 55, 56, 27, 65, 109, 97, 122, 105, 110, 103, 32, 84, 97, 108, 101, 115, 32, 83, 101, 99, 111, 110, 100, 32, 83, 101, 114, 105, 101, 115, 12, 76, 105, 110, 103, 32, 77, 101, 110, 103, 99, 104, 117, 4, 49, 54, 50, 56, 11, 80, 101, 114, 115, 112, 101, 99, 116, 105, 118, 101}
	res.RowDatas[8] = []uint8{1, 57, 13, 57, 55, 56, 57, 56, 54, 49, 50, 55, 51, 49, 50, 57, 23, 73, 110, 118, 101, 115, 116, 105, 116, 117, 114, 101, 32, 79, 102, 32, 84, 104, 101, 32, 71, 111, 100, 115, 9, 76, 117, 32, 88, 105, 120, 105, 110, 103, 4, 49, 54, 48, 53, 9, 77, 121, 116, 104, 111, 108, 111, 103, 121}
	res.RowDatas[9] = []uint8{2, 49, 48, 13, 57, 55, 56, 55, 53, 52, 48, 50, 53, 49, 52, 57, 57, 21, 70, 108, 111, 119, 101, 114, 115, 32, 73, 110, 32, 84, 104, 101, 32, 77, 105, 114, 114, 111, 114, 9, 76, 105, 32, 82, 117, 122, 104, 101, 110, 4, 49, 56, 50, 55, 15, 70, 97, 110, 116, 97, 115, 121, 32, 83, 116, 111, 114, 105, 101, 115}
	res.RowDatas[10] = []uint8{2, 49, 49, 13, 57, 55, 56, 55, 53, 48, 56, 53, 51, 53, 50, 57, 54, 19, 83, 116, 111, 114, 105, 101, 115, 32, 79, 108, 100, 32, 65, 110, 100, 32, 78, 101, 119, 13, 70, 101, 110, 103, 32, 77, 101, 110, 103, 108, 111, 110, 103, 4, 49, 54, 50, 48, 11, 80, 101, 114, 115, 112, 101, 99, 116, 105, 118, 101}
	res.RowDatas[11] = []uint8{2, 49, 50, 13, 57, 55, 56, 55, 49, 48, 49, 48, 57, 55, 53, 53, 57, 15, 71, 101, 110, 101, 114, 97, 108, 32, 89, 117, 101, 32, 70, 101, 105, 8, 81, 105, 97, 110, 32, 67, 97, 105, 4, 49, 55, 51, 53, 7, 72, 105, 115, 116, 111, 114, 121}
	res.RowDatas[12] = []uint8{2, 49, 51, 13, 57, 55, 56, 57, 56, 54, 51, 51, 56, 49, 48, 51, 55, 31, 84, 104, 101, 32, 71, 101, 110, 101, 114, 97, 108, 115, 32, 79, 102, 32, 84, 104, 101, 32, 89, 97, 110, 103, 32, 70, 97, 109, 105, 108, 121, 11, 81, 105, 32, 90, 104, 111, 110, 103, 108, 97, 110, 1, 48, 7, 72, 105, 115, 116, 111, 114, 121}
	res.RowDatas[13] = []uint8{2, 49, 52, 13, 57, 55, 56, 57, 56, 54, 53, 55, 48, 48, 48, 50, 55, 33, 82, 111, 109, 97, 110, 99, 101, 32, 79, 102, 32, 83, 117, 105, 32, 65, 110, 100, 32, 84, 97, 110, 103, 32, 68, 121, 110, 97, 115, 116, 105, 101, 115, 11, 67, 104, 101, 110, 32, 82, 117, 104, 101, 110, 103, 4, 49, 57, 56, 57, 7, 72, 105, 115, 116, 111, 114, 121}
	res.RowDatas[14] = []uint8{2, 49, 53, 13, 57, 55, 56, 57, 53, 55, 53, 55, 48, 57, 50, 52, 50, 34, 84, 104, 101, 32, 83, 101, 118, 101, 110, 32, 72, 101, 114, 111, 101, 115, 32, 65, 110, 100, 32, 70, 105, 118, 101, 32, 71, 97, 108, 108, 97, 110, 116, 115, 9, 83, 104, 105, 32, 89, 117, 107, 117, 110, 4, 49, 56, 55, 57, 7, 72, 105, 115, 116, 111, 114, 121}
	res.RowDatas[15] = []uint8{2, 49, 54, 13, 57, 55, 56, 57, 53, 55, 53, 55, 48, 57, 50, 52, 50, 19, 65, 32, 67, 111, 108, 108, 101, 99, 116, 105, 111, 110, 32, 79, 102, 32, 83, 104, 105, 9, 65, 110, 111, 110, 121, 109, 111, 117, 115, 4, 49, 56, 53, 48, 7, 72, 105, 115, 116, 111, 114, 121}
	res.RowDatas[16] = []uint8{2, 49, 55, 13, 57, 55, 56, 55, 53, 51, 51, 51, 48, 51, 51, 57, 54, 26, 68, 114, 101, 97, 109, 32, 79, 102, 32, 84, 104, 101, 32, 71, 114, 101, 101, 110, 32, 67, 104, 97, 109, 98, 101, 114, 4, 89, 117, 100, 97, 4, 49, 56, 55, 56, 11, 70, 97, 109, 105, 108, 121, 32, 83, 97, 103, 97}
	res.RowDatas[17] = []uint8{2, 49, 56, 13, 57, 55, 56, 55, 53, 49, 48, 52, 51, 52, 51, 52, 49, 23, 76, 97, 109, 112, 32, 73, 110, 32, 84, 104, 101, 32, 83, 105, 100, 101, 32, 83, 116, 114, 101, 101, 116, 9, 76, 105, 32, 76, 117, 121, 117, 97, 110, 4, 49, 55, 57, 48, 18, 85, 110, 111, 102, 102, 105, 99, 105, 97, 108, 32, 72, 105, 115, 116, 111, 114, 121}

	fmt.Println(res)
}
