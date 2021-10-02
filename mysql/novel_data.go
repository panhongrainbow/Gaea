package mysql

import (
	"strconv"
)

// >>>>> >>>>> >>>>> >>>>> >>>>> 向数据库查询 29 本小说的回传结果

// MakeNovelEmptyResult0 回传没有任何一本小说库存
func MakeNovelEmptyResult0() (*Result, error) {
	tmp := new(Result)
	tmp.Status = 34
	tmp.InsertID = 0
	tmp.AffectedRows = 0

	fdTest0 := fieldTestData{}
	fdTest0.def = "def"
	fdTest0.schema = "novel"
	fdTest0.table = "Book_0000"
	fdTest0.orgTable = "Book_0000"
	fdTest0.name = "BookID"
	fdTest0.orgName = "BookID"
	fdTest0.charset = 63
	fdTest0.columnLength = 11
	fdTest0.fieldtype = 3
	fdTest0.flag = 20483

	field0 := Field{}
	field0.ConvertFieldTest2Field(fdTest0)

	/*
		field0.Data = FieldData{3, 100, 101, 102, 5, 110, 111, 118, 101, 108, 9, 66, 111, 111, 107, 95, 48, 48, 48, 48, 9, 66, 111, 111, 107, 95, 48, 48, 48, 48, 6, 66, 111, 111, 107, 73, 68, 6, 66, 111, 111, 107, 73, 68, 12, 63, 0, 11, 0, 0, 0, 3, 3, 80, 0, 0, 0}
	*/

	fdTest1 := fieldTestData{}
	fdTest1.def = "def"
	fdTest1.schema = "novel"
	fdTest1.table = "Book_0000"
	fdTest1.orgTable = "Book_0000"
	fdTest1.name = "Isbn"
	fdTest1.orgName = "Isbn"
	fdTest1.charset = 63
	fdTest1.columnLength = 50
	fdTest1.fieldtype = 8
	fdTest1.flag = 4097

	field1 := Field{}
	field1.ConvertFieldTest2Field(fdTest1)

	/*
		field1.Data = FieldData{3, 100, 101, 102, 5, 110, 111, 118, 101, 108, 9, 66, 111, 111, 107, 95, 48, 48, 48, 48, 9, 66, 111, 111, 107, 95, 48, 48, 48, 48, 4, 73, 115, 98, 110, 4, 73, 115, 98, 110, 12, 63, 0, 50, 0, 0, 0, 8, 1, 16, 0, 0, 0}
	*/

	fdTest2 := fieldTestData{}
	fdTest2.def = "def"
	fdTest2.schema = "novel"
	fdTest2.table = "Book_0000"
	fdTest2.orgTable = "Book_0000"
	fdTest2.name = "Title"
	fdTest2.orgName = "Title"
	fdTest2.charset = 33
	fdTest2.columnLength = 300
	fdTest2.fieldtype = 253
	fdTest2.flag = 4097

	field2 := Field{}
	field2.ConvertFieldTest2Field(fdTest2)

	/*field2.Data = FieldData{3, 100, 101, 102, 5, 110, 111, 118, 101, 108, 9, 66, 111, 111, 107, 95, 48, 48, 48, 48, 9, 66, 111, 111, 107, 95, 48, 48, 48, 48, 5, 84, 105, 116, 108, 101, 5, 84, 105, 116, 108, 101, 12, 33, 0, 44, 1, 0, 0, 253, 1, 16, 0, 0, 0}
	field2.Schema = []uint8{110, 111, 118, 101, 108}
	field2.Table = []uint8{66, 111, 111, 107, 95, 48, 48, 48, 48}
	field2.OrgTable = []uint8{66, 111, 111, 107, 95, 48, 48, 48, 48}
	field2.Name = []uint8{84, 105, 116, 108, 101}
	field2.OrgName = []uint8{84, 105, 116, 108, 101}
	field2.Charset = 33
	field2.ColumnLength = 300
	field2.Type = 253
	field2.Flag = 4097
	field2.Decimal = 0
	field2.DefaultValueLength = 0
	field2.DefaultValue = nil*/

	field3 := Field{}
	field3.Data = FieldData{3, 100, 101, 102, 5, 110, 111, 118, 101, 108, 9, 66, 111, 111, 107, 95, 48, 48, 48, 48, 9, 66, 111, 111, 107, 95, 48, 48, 48, 48, 6, 65, 117, 116, 104, 111, 114, 6, 65, 117, 116, 104, 111, 114, 12, 33, 0, 90, 0, 0, 0, 253, 0, 0, 0, 0, 0}
	field3.Schema = []uint8{110, 111, 118, 101, 108}
	field3.Table = []uint8{66, 111, 111, 107, 95, 48, 48, 48, 48}
	field3.OrgTable = []uint8{66, 111, 111, 107, 95, 48, 48, 48, 48}
	field3.Name = []uint8{65, 117, 116, 104, 111, 114}
	field3.OrgName = []uint8{65, 117, 116, 104, 111, 114}
	field3.Charset = 33
	field3.ColumnLength = 90
	field3.Type = 253
	field3.Flag = 0
	field3.Decimal = 0
	field3.DefaultValueLength = 0
	field3.DefaultValue = nil

	field4 := Field{}
	field4.Data = FieldData{3, 100, 101, 102, 5, 110, 111, 118, 101, 108, 9, 66, 111, 111, 107, 95, 48, 48, 48, 48, 9, 66, 111, 111, 107, 95, 48, 48, 48, 48, 7, 80, 117, 98, 108, 105, 115, 104, 7, 80, 117, 98, 108, 105, 115, 104, 12, 63, 0, 4, 0, 0, 0, 3, 0, 0, 0, 0, 0}
	field4.Schema = []uint8{110, 111, 118, 101, 108}
	field4.Table = []uint8{66, 111, 111, 107, 95, 48, 48, 48, 48}
	field4.OrgTable = []uint8{66, 111, 111, 107, 95, 48, 48, 48, 48}
	field4.Name = []uint8{80, 117, 98, 108, 105, 115, 104}
	field4.OrgName = []uint8{80, 117, 98, 108, 105, 115, 104}
	field4.Charset = 63
	field4.ColumnLength = 4
	field4.Type = 3
	field4.Flag = 0
	field4.Decimal = 0
	field4.DefaultValueLength = 0
	field4.DefaultValue = nil

	field5 := Field{}
	field5.Data = FieldData{3, 100, 101, 102, 5, 110, 111, 118, 101, 108, 9, 66, 111, 111, 107, 95, 48, 48, 48, 48, 9, 66, 111, 111, 107, 95, 48, 48, 48, 48, 8, 67, 97, 116, 101, 103, 111, 114, 121, 8, 67, 97, 116, 101, 103, 111, 114, 121, 12, 33, 0, 90, 0, 0, 0, 253, 1, 16, 0, 0, 0}
	field5.Schema = []uint8{110, 111, 118, 101, 108}
	field5.Table = []uint8{66, 111, 111, 107, 95, 48, 48, 48, 48}
	field5.OrgTable = []uint8{66, 111, 111, 107, 95, 48, 48, 48, 48}
	field5.Name = []uint8{67, 97, 116, 101, 103, 111, 114, 121}
	field5.OrgName = []uint8{67, 97, 116, 101, 103, 111, 114, 121}
	field5.Charset = 33
	field5.ColumnLength = 90
	field5.Type = 253
	field5.Flag = 4097
	field5.Decimal = 0
	field5.DefaultValueLength = 0
	field5.DefaultValue = nil

	field6 := Field{}
	field6.Data = FieldData{3, 100, 101, 102, 5, 110, 111, 118, 101, 108, 9, 66, 111, 111, 107, 95, 48, 48, 48, 48, 9, 66, 111, 111, 107, 95, 48, 48, 48, 48, 6, 66, 111, 111, 107, 73, 68, 6, 66, 111, 111, 107, 73, 68, 12, 63, 0, 11, 0, 0, 0, 3, 3, 80, 0, 0, 0}
	field6.Schema = []uint8{110, 111, 118, 101, 108}
	field6.Table = []uint8{66, 111, 111, 107, 95, 48, 48, 48, 48}
	field6.OrgTable = []uint8{66, 111, 111, 107, 95, 48, 48, 48, 48}
	field6.Name = []uint8{66, 111, 111, 107, 73, 68}
	field6.OrgName = []uint8{66, 111, 111, 107, 73, 68}
	field6.Charset = 63
	field6.ColumnLength = 11
	field6.Type = 3
	field6.Flag = 20483
	field6.Decimal = 0
	field6.DefaultValueLength = 0
	field6.DefaultValue = nil

	resultset := Resultset{}
	tmp.Resultset = &resultset
	tmp.Resultset.Fields = []*Field{}
	tmp.Resultset.Fields = append(tmp.Resultset.Fields, &field0)
	tmp.Resultset.Fields = append(tmp.Resultset.Fields, &field1)
	tmp.Resultset.Fields = append(tmp.Resultset.Fields, &field2)
	tmp.Resultset.Fields = append(tmp.Resultset.Fields, &field3)
	tmp.Resultset.Fields = append(tmp.Resultset.Fields, &field4)
	tmp.Resultset.Fields = append(tmp.Resultset.Fields, &field5)
	tmp.Resultset.Fields = append(tmp.Resultset.Fields, &field6)

	tmp.FieldNames = make(map[string]int)
	tmp.FieldNames["BookID"] = 6
	tmp.FieldNames["Isbn"] = 1
	tmp.FieldNames["Title"] = 2
	tmp.FieldNames["Author"] = 3
	tmp.FieldNames["Publish"] = 4
	tmp.FieldNames["Category"] = 5

	tmp.Values = make([][]interface{}, 0, 29)
	tmp.RowDatas = make([]RowData, 0, 29)

	return tmp, nil
}

// MakeNovelEmptyResult1 回传没有任何一本小说库存
func MakeNovelEmptyResult1() (*Result, error) {
	tmp := new(Result)
	tmp.Status = 34
	tmp.InsertID = 0
	tmp.AffectedRows = 0

	field0 := Field{}
	field0.Data = FieldData{3, 100, 101, 102, 5, 110, 111, 118, 101, 108, 9, 66, 111, 111, 107, 95, 48, 48, 48, 49, 9, 66, 111, 111, 107, 95, 48, 48, 48, 48, 6, 66, 111, 111, 107, 73, 68, 6, 66, 111, 111, 107, 73, 68, 12, 63, 0, 11, 0, 0, 0, 3, 3, 80, 0, 0, 0}
	field0.Schema = []uint8{110, 111, 118, 101, 108}
	field0.Table = []uint8{66, 111, 111, 107, 95, 48, 48, 48, 49}
	field0.OrgTable = []uint8{66, 111, 111, 107, 95, 48, 48, 48, 49}
	field0.Name = []uint8{66, 111, 111, 107, 73, 68}
	field0.OrgName = []uint8{66, 111, 111, 107, 73, 68}
	field0.Charset = 63
	field0.ColumnLength = 11
	field0.Type = 3
	field0.Flag = 20483
	field0.Decimal = 0
	field0.DefaultValueLength = 0
	field0.DefaultValue = nil

	field1 := Field{}
	field1.Data = FieldData{3, 100, 101, 102, 5, 110, 111, 118, 101, 108, 9, 66, 111, 111, 107, 95, 48, 48, 48, 49, 9, 66, 111, 111, 107, 95, 48, 48, 48, 48, 4, 73, 115, 98, 110, 4, 73, 115, 98, 110, 12, 63, 0, 50, 0, 0, 0, 8, 1, 16, 0, 0, 0}
	field1.Schema = []uint8{110, 111, 118, 101, 108}
	field1.Table = []uint8{66, 111, 111, 107, 95, 48, 48, 48, 49}
	field1.OrgTable = []uint8{66, 111, 111, 107, 95, 48, 48, 48, 49}
	field1.Name = []uint8{73, 115, 98, 110}
	field1.OrgName = []uint8{73, 115, 98, 110}
	field1.Charset = 63
	field1.ColumnLength = 50
	field1.Type = 8
	field1.Flag = 4097
	field1.Decimal = 0
	field1.DefaultValueLength = 0
	field1.DefaultValue = nil

	field2 := Field{}
	field2.Data = FieldData{3, 100, 101, 102, 5, 110, 111, 118, 101, 108, 9, 66, 111, 111, 107, 95, 48, 48, 48, 49, 9, 66, 111, 111, 107, 95, 48, 48, 48, 48, 5, 84, 105, 116, 108, 101, 5, 84, 105, 116, 108, 101, 12, 33, 0, 44, 1, 0, 0, 253, 1, 16, 0, 0, 0}
	field2.Schema = []uint8{110, 111, 118, 101, 108}
	field2.Table = []uint8{66, 111, 111, 107, 95, 48, 48, 48, 49}
	field2.OrgTable = []uint8{66, 111, 111, 107, 95, 48, 48, 48, 49}
	field2.Name = []uint8{84, 105, 116, 108, 101}
	field2.OrgName = []uint8{84, 105, 116, 108, 101}
	field2.Charset = 33
	field2.ColumnLength = 300
	field2.Type = 253
	field2.Flag = 4097
	field2.Decimal = 0
	field2.DefaultValueLength = 0
	field2.DefaultValue = nil

	field3 := Field{}
	field3.Data = FieldData{3, 100, 101, 102, 5, 110, 111, 118, 101, 108, 9, 66, 111, 111, 107, 95, 48, 48, 48, 49, 9, 66, 111, 111, 107, 95, 48, 48, 48, 48, 6, 65, 117, 116, 104, 111, 114, 6, 65, 117, 116, 104, 111, 114, 12, 33, 0, 90, 0, 0, 0, 253, 0, 0, 0, 0, 0}
	field3.Schema = []uint8{110, 111, 118, 101, 108}
	field3.Table = []uint8{66, 111, 111, 107, 95, 48, 48, 48, 49}
	field3.OrgTable = []uint8{66, 111, 111, 107, 95, 48, 48, 48, 49}
	field3.Name = []uint8{65, 117, 116, 104, 111, 114}
	field3.OrgName = []uint8{65, 117, 116, 104, 111, 114}
	field3.Charset = 33
	field3.ColumnLength = 90
	field3.Type = 253
	field3.Flag = 0
	field3.Decimal = 0
	field3.DefaultValueLength = 0
	field3.DefaultValue = nil

	field4 := Field{}
	field4.Data = FieldData{3, 100, 101, 102, 5, 110, 111, 118, 101, 108, 9, 66, 111, 111, 107, 95, 48, 48, 48, 49, 9, 66, 111, 111, 107, 95, 48, 48, 48, 48, 7, 80, 117, 98, 108, 105, 115, 104, 7, 80, 117, 98, 108, 105, 115, 104, 12, 63, 0, 4, 0, 0, 0, 3, 0, 0, 0, 0, 0}
	field4.Schema = []uint8{110, 111, 118, 101, 108}
	field4.Table = []uint8{66, 111, 111, 107, 95, 48, 48, 48, 49}
	field4.OrgTable = []uint8{66, 111, 111, 107, 95, 48, 48, 48, 49}
	field4.Name = []uint8{80, 117, 98, 108, 105, 115, 104}
	field4.OrgName = []uint8{80, 117, 98, 108, 105, 115, 104}
	field4.Charset = 63
	field4.ColumnLength = 4
	field4.Type = 3
	field4.Flag = 0
	field4.Decimal = 0
	field4.DefaultValueLength = 0
	field4.DefaultValue = nil

	field5 := Field{}
	field5.Data = FieldData{3, 100, 101, 102, 5, 110, 111, 118, 101, 108, 9, 66, 111, 111, 107, 95, 48, 48, 48, 49, 9, 66, 111, 111, 107, 95, 48, 48, 48, 48, 8, 67, 97, 116, 101, 103, 111, 114, 121, 8, 67, 97, 116, 101, 103, 111, 114, 121, 12, 33, 0, 90, 0, 0, 0, 253, 1, 16, 0, 0, 0}
	field5.Schema = []uint8{110, 111, 118, 101, 108}
	field5.Table = []uint8{66, 111, 111, 107, 95, 48, 48, 48, 49}
	field5.OrgTable = []uint8{66, 111, 111, 107, 95, 48, 48, 48, 49}
	field5.Name = []uint8{67, 97, 116, 101, 103, 111, 114, 121}
	field5.OrgName = []uint8{67, 97, 116, 101, 103, 111, 114, 121}
	field5.Charset = 33
	field5.ColumnLength = 90
	field5.Type = 253
	field5.Flag = 4097
	field5.Decimal = 0
	field5.DefaultValueLength = 0
	field5.DefaultValue = nil

	field6 := Field{}
	field6.Data = FieldData{3, 100, 101, 102, 5, 110, 111, 118, 101, 108, 9, 66, 111, 111, 107, 95, 48, 48, 48, 49, 9, 66, 111, 111, 107, 95, 48, 48, 48, 48, 6, 66, 111, 111, 107, 73, 68, 6, 66, 111, 111, 107, 73, 68, 12, 63, 0, 11, 0, 0, 0, 3, 3, 80, 0, 0, 0}
	field6.Schema = []uint8{110, 111, 118, 101, 108}
	field6.Table = []uint8{66, 111, 111, 107, 95, 48, 48, 48, 49}
	field6.OrgTable = []uint8{66, 111, 111, 107, 95, 48, 48, 48, 49}
	field6.Name = []uint8{66, 111, 111, 107, 73, 68}
	field6.OrgName = []uint8{66, 111, 111, 107, 73, 68}
	field6.Charset = 63
	field6.ColumnLength = 11
	field6.Type = 3
	field6.Flag = 20483
	field6.Decimal = 0
	field6.DefaultValueLength = 0
	field6.DefaultValue = nil

	resultset := Resultset{}
	tmp.Resultset = &resultset
	tmp.Resultset.Fields = []*Field{}
	tmp.Resultset.Fields = append(tmp.Resultset.Fields, &field0)
	tmp.Resultset.Fields = append(tmp.Resultset.Fields, &field1)
	tmp.Resultset.Fields = append(tmp.Resultset.Fields, &field2)
	tmp.Resultset.Fields = append(tmp.Resultset.Fields, &field3)
	tmp.Resultset.Fields = append(tmp.Resultset.Fields, &field4)
	tmp.Resultset.Fields = append(tmp.Resultset.Fields, &field5)
	tmp.Resultset.Fields = append(tmp.Resultset.Fields, &field6)

	tmp.FieldNames = make(map[string]int)
	tmp.FieldNames["BookID"] = 6
	tmp.FieldNames["Isbn"] = 1
	tmp.FieldNames["Title"] = 2
	tmp.FieldNames["Author"] = 3
	tmp.FieldNames["Publish"] = 4
	tmp.FieldNames["Category"] = 5

	tmp.Values = make([][]interface{}, 0, 29)
	tmp.RowDatas = make([]RowData, 0, 29)

	return tmp, nil
}

// SelectNovelResult 回传所有 29 本小说数据库测试资料
func SelectNovelResult() (*Result, error) {
	res := Result{}
	res.Status = 34
	res.InsertID = 0
	res.AffectedRows = 0

	field0 := Field{}
	field0.Data = FieldData{3, 100, 101, 102, 7, 76, 105, 98, 114, 97, 114, 121, 4, 66, 111, 111, 107, 4, 66, 111, 111, 107, 6, 66, 111, 111, 107, 73, 68, 6, 66, 111, 111, 107, 73, 68, 12, 63, 0, 11, 0, 0, 0, 3, 3, 80, 0, 0, 0}
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

	field1 := Field{}
	field1.Data = FieldData{3, 100, 101, 102, 7, 76, 105, 98, 114, 97, 114, 121, 4, 66, 111, 111, 107, 4, 66, 111, 111, 107, 4, 73, 115, 98, 110, 4, 73, 115, 98, 110, 12, 63, 0, 50, 0, 0, 0, 8, 1, 16, 0, 0, 0}
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

	field2 := Field{}
	field2.Data = FieldData{3, 100, 101, 102, 7, 76, 105, 98, 114, 97, 114, 121, 4, 66, 111, 111, 107, 4, 66, 111, 111, 107, 5, 84, 105, 116, 108, 101, 5, 84, 105, 116, 108, 101, 12, 33, 0, 44, 1, 0, 0, 253, 1, 16, 0, 0, 0}
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

	field3 := Field{}
	field3.Data = FieldData{3, 100, 101, 102, 7, 76, 105, 98, 114, 97, 114, 121, 4, 66, 111, 111, 107, 4, 66, 111, 111, 107, 6, 65, 117, 116, 104, 111, 114, 6, 65, 117, 116, 104, 111, 114, 12, 33, 0, 90, 0, 0, 0, 253, 0, 0, 0, 0, 0}
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

	field4 := Field{}
	field4.Data = FieldData{3, 100, 101, 102, 7, 76, 105, 98, 114, 97, 114, 121, 4, 66, 111, 111, 107, 4, 66, 111, 111, 107, 7, 80, 117, 98, 108, 105, 115, 104, 7, 80, 117, 98, 108, 105, 115, 104, 12, 63, 0, 4, 0, 0, 0, 3, 0, 0, 0, 0, 0}
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

	field5 := Field{}
	field5.Data = FieldData{3, 100, 101, 102, 7, 76, 105, 98, 114, 97, 114, 121, 4, 66, 111, 111, 107, 4, 66, 111, 111, 107, 8, 67, 97, 116, 101, 103, 111, 114, 121, 8, 67, 97, 116, 101, 103, 111, 114, 121, 12, 33, 0, 90, 0, 0, 0, 253, 1, 16, 0, 0, 0}
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

	resultset := Resultset{}
	res.Resultset = &resultset
	res.Resultset.Fields = []*Field{}
	res.Resultset.Fields = append(res.Resultset.Fields, &field0)
	res.Resultset.Fields = append(res.Resultset.Fields, &field1)
	res.Resultset.Fields = append(res.Resultset.Fields, &field2)
	res.Resultset.Fields = append(res.Resultset.Fields, &field3)
	res.Resultset.Fields = append(res.Resultset.Fields, &field4)
	res.Resultset.Fields = append(res.Resultset.Fields, &field5)
	res.Resultset.Fields = append(res.Resultset.Fields, &field0)

	res.FieldNames = make(map[string]int)
	res.FieldNames["BookID"] = 6
	res.FieldNames["Isbn"] = 1
	res.FieldNames["Title"] = 2
	res.FieldNames["Author"] = 3
	res.FieldNames["Publish"] = 4
	res.FieldNames["Category"] = 5

	res.Values = make([][]interface{}, 29)
	res.Values[0] = make([]interface{}, 6)
	res.Values[0][0] = int64(1)
	res.Values[0][1] = int64(9781517191276)
	res.Values[0][2] = "Romance Of The Three Kingdoms"
	res.Values[0][3] = "Luo Guanzhong"
	res.Values[0][4] = 1522
	res.Values[0][5] = "Historical fiction"

	res.Values[1] = make([]interface{}, 6)
	res.Values[1][0] = int64(2)
	res.Values[1][1] = int64(9789869442060)
	res.Values[1][2] = "Water Margin"
	res.Values[1][3] = "Shi Nai an"
	res.Values[1][4] = 1589
	res.Values[1][5] = "Historical fiction"

	res.Values[2] = make([]interface{}, 6)
	res.Values[2][0] = int64(3)
	res.Values[2][1] = int64(9789575709518)
	res.Values[2][2] = "Journey To The West"
	res.Values[2][3] = "Wu Cheng en"
	res.Values[2][4] = 1592
	res.Values[2][5] = "Gods And Demons Fiction"

	res.Values[3] = make([]interface{}, 6)
	res.Values[3][0] = int64(4)
	res.Values[3][1] = int64(9789865975364)
	res.Values[3][2] = "Dream Of The Red Chamber"
	res.Values[3][3] = "Cao Xueqin"
	res.Values[3][4] = 1791
	res.Values[3][5] = "Family Saga"

	res.Values[4] = make([]interface{}, 6)
	res.Values[4][0] = int64(5)
	res.Values[4][1] = int64(9780804847773)
	res.Values[4][2] = "Jin Ping Mei"
	res.Values[4][3] = "Lanling Xiaoxiao Sheng"
	res.Values[4][4] = 1610
	res.Values[4][5] = "Family Saga"

	res.Values[5] = make([]interface{}, 6)
	res.Values[5][0] = int64(6)
	res.Values[5][1] = int64(9780835124072)
	res.Values[5][2] = "Rulin Waishi"
	res.Values[5][3] = "Wu Jingzi"
	res.Values[5][4] = 1750
	res.Values[5][5] = "Unofficial History"

	res.Values[6] = make([]interface{}, 6)
	res.Values[6][0] = int64(7)
	res.Values[6][1] = int64(9787101064100)
	res.Values[6][2] = "Amazing Tales First Series"
	res.Values[6][3] = "Ling Mengchu"
	res.Values[6][4] = 1628
	res.Values[6][5] = "Perspective"

	res.Values[7] = make([]interface{}, 6)
	res.Values[7][0] = int64(8)
	res.Values[7][1] = int64(9789571447278)
	res.Values[7][2] = "Amazing Tales Second Series"
	res.Values[7][3] = "Ling Mengchu"
	res.Values[7][4] = 1628
	res.Values[7][5] = "Perspective"

	res.Values[8] = make([]interface{}, 6)
	res.Values[8][0] = int64(9)
	res.Values[8][1] = int64(9789861273129)
	res.Values[8][2] = "Investiture Of The Gods"
	res.Values[8][3] = "Lu Xixing"
	res.Values[8][4] = 1605
	res.Values[8][5] = "Mythology"

	res.Values[9] = make([]interface{}, 6)
	res.Values[9][0] = int64(10)
	res.Values[9][1] = int64(9787540251499)
	res.Values[9][2] = "Flowers In The Mirror"
	res.Values[9][3] = "Li Ruzhen"
	res.Values[9][4] = 1827
	res.Values[9][5] = "Fantasy Stories"

	res.Values[10] = make([]interface{}, 6)
	res.Values[10][0] = int64(11)
	res.Values[10][1] = int64(9787508535296)
	res.Values[10][2] = "Stories Old And New"
	res.Values[10][3] = "Feng Menglong"
	res.Values[10][4] = 1620
	res.Values[10][5] = "Perspective"

	res.Values[11] = make([]interface{}, 6)
	res.Values[11][0] = int64(12)
	res.Values[11][1] = int64(9787101097559)
	res.Values[11][2] = "General Yue Fei"
	res.Values[11][3] = "Qian Cai"
	res.Values[11][4] = 1735
	res.Values[11][5] = "History"

	res.Values[12] = make([]interface{}, 6)
	res.Values[12][0] = int64(13)
	res.Values[12][1] = int64(9789863381037)
	res.Values[12][2] = "The Generals Of The Yang Family"
	res.Values[12][3] = "Qi Zhonglan"
	res.Values[12][4] = 0
	res.Values[12][5] = "History"

	res.Values[13] = make([]interface{}, 6)
	res.Values[13][0] = int64(14)
	res.Values[13][1] = int64(9789865700027)
	res.Values[13][2] = "Romance Of Sui And Tang Dynasties"
	res.Values[13][3] = "Chen Ruheng"
	res.Values[13][4] = 1989
	res.Values[13][5] = "History"

	res.Values[14] = make([]interface{}, 6)
	res.Values[14][0] = int64(15)
	res.Values[14][1] = int64(9789575709242)
	res.Values[14][2] = "The Seven Heroes And Five Gallants"
	res.Values[14][3] = "Shi Yukun"
	res.Values[14][4] = 1879
	res.Values[14][5] = "History"

	res.Values[15] = make([]interface{}, 6)
	res.Values[15][0] = int64(16)
	res.Values[15][1] = int64(9789574927913)
	res.Values[15][2] = "A Collection Of Shi"
	res.Values[15][3] = "Anonymous"
	res.Values[15][4] = 1850
	res.Values[15][5] = "History"

	res.Values[16] = make([]interface{}, 6)
	res.Values[16][0] = int64(17)
	res.Values[16][1] = int64(9787533303396)
	res.Values[16][2] = "Dream Of The Green Chamber"
	res.Values[16][3] = "Yuda"
	res.Values[16][4] = 1878
	res.Values[16][5] = "Family Saga"

	res.Values[17] = make([]interface{}, 6)
	res.Values[17][0] = int64(18)
	res.Values[17][1] = int64(9787510434341)
	res.Values[17][2] = "Lamp In The Side Street"
	res.Values[17][3] = "Li Luyuan"
	res.Values[17][4] = 1790
	res.Values[17][5] = "Unofficial History"

	res.Values[18] = make([]interface{}, 6)
	res.Values[18][0] = int64(19)
	res.Values[18][1] = int64(9789571447469)
	res.Values[18][2] = "The Travels of Lao Can"
	res.Values[18][3] = "Liu E"
	res.Values[18][4] = 1907
	res.Values[18][5] = "Social Story"

	res.Values[19] = make([]interface{}, 6)
	res.Values[19][0] = int64(20)
	res.Values[19][1] = int64(9789571470047)
	res.Values[19][2] = "Bizarre Happenings Eyewitnessed over Two Decades"
	res.Values[19][3] = "Jianren Wu"
	res.Values[19][4] = 1905
	res.Values[19][5] = "Unofficial History"

	res.Values[20] = make([]interface{}, 6)
	res.Values[20][0] = int64(21)
	res.Values[20][1] = int64(9787101097580)
	res.Values[20][2] = "A Flower In A Sinful Sea"
	res.Values[20][3] = "Zeng Pu"
	res.Values[20][4] = 1904
	res.Values[20][5] = "History"

	res.Values[21] = make([]interface{}, 6)
	res.Values[21][0] = int64(22)
	res.Values[21][1] = int64(9789861674193)
	res.Values[21][2] = "Officialdom Unmasked"
	res.Values[21][3] = "Li Baojia"
	res.Values[21][4] = 1903
	res.Values[21][5] = "Unofficial History"

	res.Values[22] = make([]interface{}, 6)
	res.Values[22][0] = int64(23)
	res.Values[22][1] = int64(9787805469836)
	res.Values[22][2] = "Tower For The Summer Heat"
	res.Values[22][3] = "Li Yu"
	res.Values[22][4] = 1680
	res.Values[22][5] = "Unofficial History"

	res.Values[23] = make([]interface{}, 6)
	res.Values[23][0] = int64(24)
	res.Values[23][1] = int64(9787508067247)
	res.Values[23][2] = "Silent Operas"
	res.Values[23][3] = "Li Yu"
	res.Values[23][4] = 1680
	res.Values[23][5] = "Unofficial History"

	res.Values[24] = make([]interface{}, 6)
	res.Values[24][0] = int64(25)
	res.Values[24][1] = int64(9789573609049)
	res.Values[24][2] = "The Carnal Prayer Mat"
	res.Values[24][3] = "Li Yu"
	res.Values[24][4] = 1680
	res.Values[24][5] = "Unofficial History"

	res.Values[25] = make([]interface{}, 6)
	res.Values[25][0] = int64(26)
	res.Values[25][1] = int64(9787533948108)
	res.Values[25][2] = "Six Records Of A Floating Life"
	res.Values[25][3] = "Shen Fu"
	res.Values[25][4] = 1878
	res.Values[25][5] = "Autobiography"

	res.Values[26] = make([]interface{}, 6)
	res.Values[26][0] = int64(27)
	res.Values[26][1] = int64(9786666141110)
	res.Values[26][2] = "Humble Words Of A Rustic Elder"
	res.Values[26][3] = "Xia Jingqu"
	res.Values[26][4] = 1787
	res.Values[26][5] = "Historical fiction"

	res.Values[27] = make([]interface{}, 6)
	res.Values[27][0] = int64(28)
	res.Values[27][1] = int64(9789571435473)
	res.Values[27][2] = "Nine-Tailed Turtle"
	res.Values[27][3] = "Lu Can"
	res.Values[27][4] = 1551
	res.Values[27][5] = "Mythology"

	res.Values[28] = make([]interface{}, 6)
	res.Values[28][0] = int64(29)
	res.Values[28][1] = int64(9789866318603)
	res.Values[28][2] = "A History Of Floral Treasures"
	res.Values[28][3] = "Chen Sen"
	res.Values[28][4] = 1849
	res.Values[28][5] = "Romance"

	res.RowDatas = make([]RowData, 29)
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
	res.RowDatas[15] = []uint8{2, 49, 54, 13, 57, 55, 56, 57, 53, 55, 52, 57, 50, 55, 57, 49, 51, 19, 65, 32, 67, 111, 108, 108, 101, 99, 116, 105, 111, 110, 32, 79, 102, 32, 83, 104, 105, 9, 65, 110, 111, 110, 121, 109, 111, 117, 115, 4, 49, 56, 53, 48, 7, 72, 105, 115, 116, 111, 114, 121}
	res.RowDatas[16] = []uint8{2, 49, 55, 13, 57, 55, 56, 55, 53, 51, 51, 51, 48, 51, 51, 57, 54, 26, 68, 114, 101, 97, 109, 32, 79, 102, 32, 84, 104, 101, 32, 71, 114, 101, 101, 110, 32, 67, 104, 97, 109, 98, 101, 114, 4, 89, 117, 100, 97, 4, 49, 56, 55, 56, 11, 70, 97, 109, 105, 108, 121, 32, 83, 97, 103, 97}
	res.RowDatas[17] = []uint8{2, 49, 56, 13, 57, 55, 56, 55, 53, 49, 48, 52, 51, 52, 51, 52, 49, 23, 76, 97, 109, 112, 32, 73, 110, 32, 84, 104, 101, 32, 83, 105, 100, 101, 32, 83, 116, 114, 101, 101, 116, 9, 76, 105, 32, 76, 117, 121, 117, 97, 110, 4, 49, 55, 57, 48, 18, 85, 110, 111, 102, 102, 105, 99, 105, 97, 108, 32, 72, 105, 115, 116, 111, 114, 121}
	res.RowDatas[18] = []uint8{2, 49, 57, 13, 57, 55, 56, 57, 53, 55, 49, 52, 52, 55, 52, 54, 57, 22, 84, 104, 101, 32, 84, 114, 97, 118, 101, 108, 115, 32, 111, 102, 32, 76, 97, 111, 32, 67, 97, 110, 5, 76, 105, 117, 32, 69, 4, 49, 57, 48, 55, 12, 83, 111, 99, 105, 97, 108, 32, 83, 116, 111, 114, 121}
	res.RowDatas[19] = []uint8{2, 50, 48, 13, 57, 55, 56, 57, 53, 55, 49, 52, 55, 48, 48, 52, 55, 48, 66, 105, 122, 97, 114, 114, 101, 32, 72, 97, 112, 112, 101, 110, 105, 110, 103, 115, 32, 69, 121, 101, 119, 105, 116, 110, 101, 115, 115, 101, 100, 32, 111, 118, 101, 114, 32, 84, 119, 111, 32, 68, 101, 99, 97, 100, 101, 115, 10, 74, 105, 97, 110, 114, 101, 110, 32, 87, 117, 4, 49, 57, 48, 53, 18, 85, 110, 111, 102, 102, 105, 99, 105, 97, 108, 32, 72, 105, 115, 116, 111, 114, 121}
	res.RowDatas[20] = []uint8{2, 50, 49, 13, 57, 55, 56, 55, 49, 48, 49, 48, 57, 55, 53, 56, 48, 24, 65, 32, 70, 108, 111, 119, 101, 114, 32, 73, 110, 32, 65, 32, 83, 105, 110, 102, 117, 108, 32, 83, 101, 97, 7, 90, 101, 110, 103, 32, 80, 117, 4, 49, 57, 48, 52, 7, 72, 105, 115, 116, 111, 114, 121}
	res.RowDatas[21] = []uint8{2, 50, 50, 13, 57, 55, 56, 57, 56, 54, 49, 54, 55, 52, 49, 57, 51, 20, 79, 102, 102, 105, 99, 105, 97, 108, 100, 111, 109, 32, 85, 110, 109, 97, 115, 107, 101, 100, 9, 76, 105, 32, 66, 97, 111, 106, 105, 97, 4, 49, 57, 48, 51, 18, 85, 110, 111, 102, 102, 105, 99, 105, 97, 108, 32, 72, 105, 115, 116, 111, 114, 121}
	res.RowDatas[22] = []uint8{2, 50, 51, 13, 57, 55, 56, 55, 56, 48, 53, 52, 54, 57, 56, 51, 54, 25, 84, 111, 119, 101, 114, 32, 70, 111, 114, 32, 84, 104, 101, 32, 83, 117, 109, 109, 101, 114, 32, 72, 101, 97, 116, 5, 76, 105, 32, 89, 117, 4, 49, 54, 56, 48, 18, 85, 110, 111, 102, 102, 105, 99, 105, 97, 108, 32, 72, 105, 115, 116, 111, 114, 121}
	res.RowDatas[23] = []uint8{2, 50, 52, 13, 57, 55, 56, 55, 53, 48, 56, 48, 54, 55, 50, 52, 55, 13, 83, 105, 108, 101, 110, 116, 32, 79, 112, 101, 114, 97, 115, 5, 76, 105, 32, 89, 117, 4, 49, 54, 56, 48, 12, 83, 111, 99, 105, 97, 108, 32, 83, 116, 111, 114, 121}
	res.RowDatas[24] = []uint8{2, 50, 53, 13, 57, 55, 56, 57, 53, 55, 51, 54, 48, 57, 48, 52, 57, 21, 84, 104, 101, 32, 67, 97, 114, 110, 97, 108, 32, 80, 114, 97, 121, 101, 114, 32, 77, 97, 116, 5, 76, 105, 32, 89, 117, 4, 49, 54, 56, 48, 12, 83, 111, 99, 105, 97, 108, 32, 83, 116, 111, 114, 121}
	res.RowDatas[25] = []uint8{2, 50, 54, 13, 57, 55, 56, 55, 53, 51, 51, 57, 52, 56, 49, 48, 56, 30, 83, 105, 120, 32, 82, 101, 99, 111, 114, 100, 115, 32, 79, 102, 32, 65, 32, 70, 108, 111, 97, 116, 105, 110, 103, 32, 76, 105, 102, 101, 7, 83, 104, 101, 110, 32, 70, 117, 4, 49, 56, 55, 56, 13, 65, 117, 116, 111, 98, 105, 111, 103, 114, 97, 112, 104, 121}
	res.RowDatas[26] = []uint8{2, 50, 55, 13, 57, 55, 56, 54, 54, 54, 54, 49, 52, 49, 49, 49, 48, 30, 72, 117, 109, 98, 108, 101, 32, 87, 111, 114, 100, 115, 32, 79, 102, 32, 65, 32, 82, 117, 115, 116, 105, 99, 32, 69, 108, 100, 101, 114, 10, 88, 105, 97, 32, 74, 105, 110, 103, 113, 117, 4, 49, 55, 56, 55, 18, 72, 105, 115, 116, 111, 114, 105, 99, 97, 108, 32, 102, 105, 99, 116, 105, 111, 110}
	res.RowDatas[27] = []uint8{2, 50, 56, 13, 57, 55, 56, 57, 53, 55, 49, 52, 51, 53, 52, 55, 51, 18, 78, 105, 110, 101, 45, 84, 97, 105, 108, 101, 100, 32, 84, 117, 114, 116, 108, 101, 6, 76, 117, 32, 67, 97, 110, 4, 49, 53, 53, 49, 9, 77, 121, 116, 104, 111, 108, 111, 103, 121}
	res.RowDatas[28] = []uint8{2, 50, 57, 13, 57, 55, 56, 57, 56, 54, 54, 51, 49, 56, 54, 48, 51, 29, 65, 32, 72, 105, 115, 116, 111, 114, 121, 32, 79, 102, 32, 70, 108, 111, 114, 97, 108, 32, 84, 114, 101, 97, 115, 117, 114, 101, 115, 8, 67, 104, 101, 110, 32, 83, 101, 110, 4, 49, 56, 52, 57, 7, 82, 111, 109, 97, 110, 99, 101}

	return &res, nil
}

func ConvertNovelData2byte(value []interface{}) []byte {
	bookid := value[0].(int64)
	bookidStr := strconv.FormatInt(bookid, 10)

	isbn := value[1].(int64)
	isbnStr := strconv.FormatInt(isbn, 10)

	titleStr := value[2].(string)

	authorStr := value[3].(string)

	publish := value[4].(int)
	publishStr := strconv.Itoa(publish)

	categoryStr := value[5].(string)

	ret := string(uint8(len(bookidStr))) +
		bookidStr +
		string(uint8(len(isbnStr))) +
		isbnStr +
		string(uint8(len(titleStr))) +
		titleStr +
		string(uint8(len(authorStr))) +
		authorStr +
		string(uint8(len(publishStr))) +
		publishStr +
		string(uint8(len(categoryStr))) +
		categoryStr +
		string(uint8(len(bookidStr))) +
		bookidStr

	return []byte(ret)
}

// >>>>> >>>>> >>>>> >>>>> >>>>> 向数据库写入 29 本小说的回传结果 (目前是一本本写入)

// InsertFirstNovelResult 函式 🧚 为插入第一本小说到模拟数据库 三国演义 (会分配到 Slice-1)
func (r *Result) InsertFirstNovelResult() (*Result, error) {
	// 新增 Value
	tmp := make([]interface{}, 7)
	tmp[0] = int64(1)
	tmp[1] = int64(9781517191276)
	tmp[2] = "Romance Of The Three Kingdoms"
	tmp[3] = "Luo Guanzhong"
	tmp[4] = 1522
	tmp[5] = "Historical fiction"
	tmp[6] = int64(1)
	r.Values = append(r.Values, tmp)

	// 新增 RowDatas
	// r.RowDatas = append(r.RowDatas, []uint8{1, 49, 13, 57, 55, 56, 49, 53, 49, 55, 49, 57, 49, 50, 55, 54, 29, 82, 111, 109, 97, 110, 99, 101, 32, 79, 102, 32, 84, 104, 101, 32, 84, 104, 114, 101, 101, 32, 75, 105, 110, 103, 100, 111, 109, 115, 13, 76, 117, 111, 32, 71, 117, 97, 110, 122, 104, 111, 110, 103, 4, 49, 53, 50, 50, 18, 72, 105, 115, 116, 111, 114, 105, 99, 97, 108, 32, 102, 105, 99, 116, 105, 111, 110, 1, 49})
	r.RowDatas = append(r.RowDatas, ConvertNovelData2byte(tmp))

	// 回传成功写入的结果
	ret := new(Result)
	ret.Status = 2
	ret.InsertID = 0
	ret.AffectedRows = 1

	// 正常回传
	return ret, nil
}

// InsertSecondNovelResult 函式 🧚 为插入第二本小说到模拟数据库 水浒传 (会分配到 Slice-0)
func (r *Result) InsertSecondNovelResult() (*Result, error) {
	// 新增 Value
	tmp := make([]interface{}, 7)
	tmp[0] = int64(2)
	tmp[1] = int64(9789869442060)
	tmp[2] = "Water Margin"
	tmp[3] = "Shi Nai an"
	tmp[4] = 1589
	tmp[5] = "Historical fiction"
	tmp[6] = int64(2)
	r.Values = append(r.Values, tmp)

	// 新增 RowDatas
	// r.RowDatas = append(r.RowDatas, []uint8{1, 50, 13, 57, 55, 56, 57, 56, 54, 57, 52, 52, 50, 48, 54, 48, 12, 87, 97, 116, 101, 114, 32, 77, 97, 114, 103, 105, 110, 10, 83, 104, 105, 32, 78, 97, 105, 32, 97, 110, 4, 49, 53, 56, 57, 18, 72, 105, 115, 116, 111, 114, 105, 99, 97, 108, 32, 102, 105, 99, 116, 105, 111, 110, 1, 50})
	r.RowDatas = append(r.RowDatas, ConvertNovelData2byte(tmp))

	// 回传成功写入的结果
	ret := new(Result)
	ret.Status = 2
	ret.InsertID = 0
	ret.AffectedRows = 1

	// 正常回传
	return ret, nil
}

// InsertThirdNovelResult 函式 🧚 为插入第三本小说到模拟数据库 西游记 (会分配到 Slice-1)
func (r *Result) InsertThirdNovelResult() (*Result, error) {
	// 新增 Value
	tmp := make([]interface{}, 7)
	tmp[0] = int64(3)
	tmp[1] = int64(9789575709518)
	tmp[2] = "Journey To The West"
	tmp[3] = "Wu Cheng en"
	tmp[4] = 1592
	tmp[5] = "Gods And Demons Fiction"
	tmp[6] = int64(3)
	r.Values = append(r.Values, tmp)

	// 新增 RowDatas
	// r.RowDatas = append(r.RowDatas, []uint8{1, 51, 13, 57, 55, 56, 57, 53, 55, 53, 55, 48, 57, 53, 49, 56, 19, 74, 111, 117, 114, 110, 101, 121, 32, 84, 111, 32, 84, 104, 101, 32, 87, 101, 115, 116, 11, 87, 117, 32, 67, 104, 101, 110, 103, 32, 101, 110, 4, 49, 53, 57, 50, 23, 71, 111, 100, 115, 32, 65, 110, 100, 32, 68, 101, 109, 111, 110, 115, 32, 70, 105, 99, 116, 105, 111, 110, 1, 51})
	r.RowDatas = append(r.RowDatas, ConvertNovelData2byte(tmp))

	// 回传成功写入的结果
	ret := new(Result)
	ret.Status = 2
	ret.InsertID = 0
	ret.AffectedRows = 1

	// 正常回传
	return ret, nil
}

// InsertFourthNovelResult 函式 🧚 为插入第四本小说到模拟数据库 红楼梦 (会分配到 Slice-0)
func (r *Result) InsertFourthNovelResult() (*Result, error) {
	// 新增 Value
	tmp := make([]interface{}, 7)
	tmp[0] = int64(4)
	tmp[1] = int64(9789865975364)
	tmp[2] = "Dream Of The Red Chamber"
	tmp[3] = "Cao Xueqin"
	tmp[4] = 1791
	tmp[5] = "Family Saga"
	tmp[6] = int64(4)
	r.Values = append(r.Values, tmp)

	// 新增 RowDatas
	// r.RowDatas = append(r.RowDatas, []uint8{1, 52, 13, 57, 55, 56, 57, 56, 54, 53, 57, 55, 53, 51, 54, 52, 24, 68, 114, 101, 97, 109, 32, 79, 102, 32, 84, 104, 101, 32, 82, 101, 100, 32, 67, 104, 97, 109, 98, 101, 114, 10, 67, 97, 111, 32, 88, 117, 101, 113, 105, 110, 4, 49, 55, 57, 49, 11, 70, 97, 109, 105, 108, 121, 32, 83, 97, 103, 97, 1, 52})
	r.RowDatas = append(r.RowDatas, ConvertNovelData2byte(tmp))

	// 回传成功写入的结果
	ret := new(Result)
	ret.Status = 2
	ret.InsertID = 0
	ret.AffectedRows = 1

	// 正常回传
	return ret, nil
}

// InsertFifthNovelResult 函式 🧚 为插入第五本小说到模拟数据库 金瓶梅 (会分配到 Slice-1)
func (r *Result) InsertFifthNovelResult() (*Result, error) {
	// 新增 Value
	tmp := make([]interface{}, 7)
	tmp[0] = int64(5)
	tmp[1] = int64(9780804847773)
	tmp[2] = "Jin Ping Mei"
	tmp[3] = "Lanling Xiaoxiao Sheng"
	tmp[4] = 1610
	tmp[5] = "Family Saga"
	tmp[6] = int64(5)
	r.Values = append(r.Values, tmp)

	// 新增 RowDatas
	// r.RowDatas = append(r.RowDatas, []uint8{1, 53, 13, 57, 55, 56, 48, 56, 48, 52, 56, 52, 55, 55, 55, 51, 12, 74, 105, 110, 32, 80, 105, 110, 103, 32, 77, 101, 105, 22, 76, 97, 110, 108, 105, 110, 103, 32, 88, 105, 97, 111, 120, 105, 97, 111, 32, 83, 104, 101, 110, 103, 4, 49, 54, 49, 48, 11, 70, 97, 109, 105, 108, 121, 32, 76, 105, 102, 101, 1, 53})
	r.RowDatas = append(r.RowDatas, ConvertNovelData2byte(tmp))

	// 回传成功写入的结果
	ret := new(Result)
	ret.Status = 2
	ret.InsertID = 0
	ret.AffectedRows = 1

	// 正常回传
	return ret, nil
}

// InsertSixthNovelResult 函式 🧚 为插入第六本小说到模拟数据库 儒林外史 (会分配到 Slice-0)
func (r *Result) InsertSixthNovelResult() (*Result, error) {
	// 新增 Value
	tmp := make([]interface{}, 7)
	tmp[0] = int64(6)
	tmp[1] = int64(9780835124072)
	tmp[2] = "Rulin Waishi"
	tmp[3] = "Wu Jingzi"
	tmp[4] = 1750
	tmp[5] = "Unofficial History"
	tmp[6] = int64(6)
	r.Values = append(r.Values, tmp)

	// 新增 RowDatas
	// r.RowDatas = append(r.RowDatas, []uint8{1, 54, 13, 57, 55, 56, 48, 56, 51, 53, 49, 50, 52, 48, 55, 50, 12, 82, 117, 108, 105, 110, 32, 87, 97, 105, 115, 104, 105, 9, 87, 117, 32, 74, 105, 110, 103, 122, 105, 4, 49, 55, 53, 48, 18, 85, 110, 111, 102, 102, 105, 99, 105, 97, 108, 32, 72, 105, 115, 116, 111, 114, 121, 1, 54})
	r.RowDatas = append(r.RowDatas, ConvertNovelData2byte(tmp))

	// 回传成功写入的结果
	ret := new(Result)
	ret.Status = 2
	ret.InsertID = 0
	ret.AffectedRows = 1

	// 正常回传
	return ret, nil
}

// InsertSeventhNovelResult 函式 🧚 为插入第七本小说到模拟数据库 初刻拍案惊奇 (会分配到 Slice-1)
func (r *Result) InsertSeventhNovelResult() (*Result, error) {
	// 新增 Value
	tmp := make([]interface{}, 7)
	tmp[0] = int64(7)
	tmp[1] = int64(9787101064100)
	tmp[2] = "Amazing Tales First Series"
	tmp[3] = "Ling Mengchu"
	tmp[4] = 1628
	tmp[5] = "Perspective"
	tmp[6] = int64(7)
	r.Values = append(r.Values, tmp)

	// 新增 RowDatas
	// r.RowDatas = append(r.RowDatas, []uint8{1, 55, 13, 57, 55, 56, 55, 49, 48, 49, 48, 54, 52, 49, 48, 48, 26, 65, 109, 97, 122, 105, 110, 103, 32, 84, 97, 108, 101, 115, 32, 70, 105, 114, 115, 116, 32, 83, 101, 114, 105, 101, 115, 12, 76, 105, 110, 103, 32, 77, 101, 110, 103, 99, 104, 117, 4, 49, 54, 50, 56, 11, 80, 101, 114, 115, 112, 101, 99, 116, 105, 118, 101, 1, 55})
	r.RowDatas = append(r.RowDatas, ConvertNovelData2byte(tmp))

	// 回传成功写入的结果
	ret := new(Result)
	ret.Status = 2
	ret.InsertID = 0
	ret.AffectedRows = 1

	// 正常回传
	return ret, nil
}

// InsertEighthNovelResult 函式 🧚 为插入第八本小说到模拟数据库 二刻拍案惊奇 (会分配到 Slice-0)
func (r *Result) InsertEighthNovelResult() (*Result, error) {
	// 新增 Value
	tmp := make([]interface{}, 7)
	tmp[0] = int64(8)
	tmp[1] = int64(9789571447278)
	tmp[2] = "Amazing Tales Second Series"
	tmp[3] = "Ling Mengchu"
	tmp[4] = 1628
	tmp[5] = "Perspective"
	tmp[6] = int64(8)
	r.Values = append(r.Values, tmp)

	// 新增 RowDatas
	// r.RowDatas = append(r.RowDatas, []uint8{1, 56, 13, 57, 55, 56, 57, 53, 55, 49, 52, 52, 55, 50, 55, 56, 27, 65, 109, 97, 122, 105, 110, 103, 32, 84, 97, 108, 101, 115, 32, 83, 101, 99, 111, 110, 100, 32, 83, 101, 114, 105, 101, 115, 12, 76, 105, 110, 103, 32, 77, 101, 110, 103, 99, 104, 117, 4, 49, 54, 50, 56, 11, 80, 101, 114, 115, 112, 101, 99, 116, 105, 118, 101, 1, 56})
	r.RowDatas = append(r.RowDatas, ConvertNovelData2byte(tmp))

	// 回传成功写入的结果
	ret := new(Result)
	ret.Status = 2
	ret.InsertID = 0
	ret.AffectedRows = 1

	// 正常回传
	return ret, nil
}

// InsertNinethNovelResult 函式 🧚 为插入第九本小说到模拟数据库 封神演义 (会分配到 Slice-1)
func (r *Result) InsertNinethNovelResult() (*Result, error) {
	// 新增 Value
	tmp := make([]interface{}, 7)
	tmp[0] = int64(9)
	tmp[1] = int64(9789861273129)
	tmp[2] = "Investiture Of The Gods"
	tmp[3] = "Lu Xixing"
	tmp[4] = 1605
	tmp[5] = "Mythology"
	tmp[6] = int64(9)
	r.Values = append(r.Values, tmp)

	// 新增 RowDatas
	// r.RowDatas = append(r.RowDatas, []uint8{1, 57, 13, 57, 55, 56, 57, 56, 54, 49, 50, 55, 51, 49, 50, 57, 23, 73, 110, 118, 101, 115, 116, 105, 116, 117, 114, 101, 32, 79, 102, 32, 84, 104, 101, 32, 71, 111, 100, 115, 9, 76, 117, 32, 88, 105, 120, 105, 110, 103, 4, 49, 54, 48, 53, 9, 77, 121, 116, 104, 111, 108, 111, 103, 121, 1, 57})
	r.RowDatas = append(r.RowDatas, ConvertNovelData2byte(tmp))

	// 回传成功写入的结果
	ret := new(Result)
	ret.Status = 2
	ret.InsertID = 0
	ret.AffectedRows = 1

	// 正常回传
	return ret, nil
}

// InsertTenthNovelResult 函式 🧚 为插入第十本小说到模拟数据库 镜花缘 (会分配到 Slice-0)
func (r *Result) InsertTenthNovelResult() (*Result, error) {
	// 新增 Value
	tmp := make([]interface{}, 7)
	tmp[0] = int64(10)
	tmp[1] = int64(9787540251499)
	tmp[2] = "Flowers In The Mirror"
	tmp[3] = "Li Ruzhen"
	tmp[4] = 1827
	tmp[5] = "Fantasy Stories"
	tmp[6] = int64(10)
	r.Values = append(r.Values, tmp)

	// 新增 RowDatas
	// r.RowDatas = append(r.RowDatas, []uint8{2, 49, 48, 13, 57, 55, 56, 55, 53, 52, 48, 50, 53, 49, 52, 57, 57, 21, 70, 108, 111, 119, 101, 114, 115, 32, 73, 110, 32, 84, 104, 101, 32, 77, 105, 114, 114, 111, 114, 9, 76, 105, 32, 82, 117, 122, 104, 101, 110, 4, 49, 56, 50, 55, 15, 70, 97, 110, 116, 97, 115, 121, 32, 83, 116, 111, 114, 105, 101, 115, 2, 49, 48})
	r.RowDatas = append(r.RowDatas, ConvertNovelData2byte(tmp))

	// 回传成功写入的结果
	ret := new(Result)
	ret.Status = 2
	ret.InsertID = 0
	ret.AffectedRows = 1

	// 正常回传
	return ret, nil
}

// InsertEleventhNovelResult 函式 🧚 为插入第十一本小说到模拟数据库 喻世明言 (会分配到 Slice-1)
func (r *Result) InsertEleventhNovelResult() (*Result, error) {
	// 新增 Value
	tmp := make([]interface{}, 7)
	tmp[0] = int64(11)
	tmp[1] = int64(9787508535296)
	tmp[2] = "Stories Old And New"
	tmp[3] = "Feng Menglong"
	tmp[4] = 1620
	tmp[5] = "Perspective"
	tmp[6] = int64(11)
	r.Values = append(r.Values, tmp)

	// 新增 RowDatas
	// r.RowDatas = append(r.RowDatas, []uint8{2, 49, 49, 13, 57, 55, 56, 55, 53, 48, 56, 53, 51, 53, 50, 57, 54, 19, 83, 116, 111, 114, 105, 101, 115, 32, 79, 108, 100, 32, 65, 110, 100, 32, 78, 101, 119, 13, 70, 101, 110, 103, 32, 77, 101, 110, 103, 108, 111, 110, 103, 4, 49, 54, 50, 48, 11, 80, 101, 114, 115, 112, 101, 99, 116, 105, 118, 101, 2, 49, 49})
	r.RowDatas = append(r.RowDatas, ConvertNovelData2byte(tmp))

	// 回传成功写入的结果
	ret := new(Result)
	ret.Status = 2
	ret.InsertID = 0
	ret.AffectedRows = 1

	// 正常回传
	return ret, nil
}

// InsertTwelfthNovelResult 函式 🧚 为插入第十二本小说到模拟数据库 说岳全传 (会分配到 Slice-0)
func (r *Result) InsertTwelfthNovelResult() (*Result, error) {
	// 新增 Value
	tmp := make([]interface{}, 7)
	tmp[0] = int64(12)
	tmp[1] = int64(9787101097559)
	tmp[2] = "General Yue Fei"
	tmp[3] = "Qian Cai"
	tmp[4] = 1735
	tmp[5] = "History"
	tmp[6] = int64(12)
	r.Values = append(r.Values, tmp)

	// 新增 RowDatas
	// r.RowDatas = append(r.RowDatas, []uint8{2, 49, 50, 13, 57, 55, 56, 55, 49, 48, 49, 48, 57, 55, 53, 53, 57, 15, 71, 101, 110, 101, 114, 97, 108, 32, 89, 117, 101, 32, 70, 101, 105, 8, 81, 105, 97, 110, 32, 67, 97, 105, 4, 49, 55, 51, 53, 7, 72, 105, 115, 116, 111, 114, 121, 2, 49, 50})
	r.RowDatas = append(r.RowDatas, ConvertNovelData2byte(tmp))

	// 回传成功写入的结果
	ret := new(Result)
	ret.Status = 2
	ret.InsertID = 0
	ret.AffectedRows = 1

	// 正常回传
	return ret, nil
}

// InsertThirteenthNovelResult 函式 🧚 为插入第十三本小说到模拟数据库 杨家将 (会分配到 Slice-1)
func (r *Result) InsertThirteenthNovelResult() (*Result, error) {
	// 新增 Value
	tmp := make([]interface{}, 7)
	tmp[0] = int64(13)
	tmp[1] = int64(9789863381037)
	tmp[2] = "The Generals Of The Yang Family"
	tmp[3] = "Qi Zhonglan"
	tmp[4] = 0
	tmp[5] = "History"
	tmp[6] = int64(13)
	r.Values = append(r.Values, tmp)

	// 新增 RowDatas
	// r.RowDatas = append(r.RowDatas, []uint8{2, 49, 51, 13, 57, 55, 56, 57, 56, 54, 51, 51, 56, 49, 48, 51, 55, 31, 84, 104, 101, 32, 71, 101, 110, 101, 114, 97, 108, 115, 32, 79, 102, 32, 84, 104, 101, 32, 89, 97, 110, 103, 32, 70, 97, 109, 105, 108, 121, 11, 81, 105, 32, 90, 104, 111, 110, 103, 108, 97, 110, 1, 48, 7, 72, 105, 115, 116, 111, 114, 121, 2, 49, 51})
	r.RowDatas = append(r.RowDatas, ConvertNovelData2byte(tmp))

	// 回传成功写入的结果
	ret := new(Result)
	ret.Status = 2
	ret.InsertID = 0
	ret.AffectedRows = 1

	// 正常回传
	return ret, nil
}

// InsertFourteenthNovelResult 函式 🧚 为插入第十四本小说到模拟数据库 说唐 (会分配到 Slice-0)
func (r *Result) InsertFourteenthNovelResult() (*Result, error) {
	// 新增 Value
	tmp := make([]interface{}, 7)
	tmp[0] = int64(14)
	tmp[1] = int64(9789865700027)
	tmp[2] = "Romance Of Sui And Tang Dynasties"
	tmp[3] = "Chen Ruheng"
	tmp[4] = 1989
	tmp[5] = "History"
	tmp[6] = int64(14)
	r.Values = append(r.Values, tmp)

	// 新增 RowDatas
	// r.RowDatas = append(r.RowDatas, []uint8{2, 49, 52, 13, 57, 55, 56, 57, 56, 54, 53, 55, 48, 48, 48, 50, 55, 33, 82, 111, 109, 97, 110, 99, 101, 32, 79, 102, 32, 83, 117, 105, 32, 65, 110, 100, 32, 84, 97, 110, 103, 32, 68, 121, 110, 97, 115, 116, 105, 101, 115, 11, 67, 104, 101, 110, 32, 82, 117, 104, 101, 110, 103, 4, 49, 57, 56, 57, 7, 72, 105, 115, 116, 111, 114, 121, 2, 49, 52})
	r.RowDatas = append(r.RowDatas, ConvertNovelData2byte(tmp))

	// 回传成功写入的结果
	ret := new(Result)
	ret.Status = 2
	ret.InsertID = 0
	ret.AffectedRows = 1

	// 正常回传
	return ret, nil
}

// InsertFifteenthNovelResult 函式 🧚 为插入第十五本小说到模拟数据库 七侠五义 (会分配到 Slice-1)
func (r *Result) InsertFifteenthNovelResult() (*Result, error) {
	// 新增 Value
	tmp := make([]interface{}, 7)
	tmp[0] = int64(15)
	tmp[1] = int64(9789575709242)
	tmp[2] = "The Seven Heroes And Five Gallants"
	tmp[3] = "Shi Yukun"
	tmp[4] = 1879
	tmp[5] = "History"
	tmp[6] = int64(15)
	r.Values = append(r.Values, tmp)

	// 新增 RowDatas
	// r.RowDatas = append(r.RowDatas, []uint8{2, 49, 53, 13, 57, 55, 56, 57, 53, 55, 53, 55, 48, 57, 50, 52, 50, 34, 84, 104, 101, 32, 83, 101, 118, 101, 110, 32, 72, 101, 114, 111, 101, 115, 32, 65, 110, 100, 32, 70, 105, 118, 101, 32, 71, 97, 108, 108, 97, 110, 116, 115, 9, 83, 104, 105, 32, 89, 117, 107, 117, 110, 4, 49, 56, 55, 57, 7, 72, 105, 115, 116, 111, 114, 121, 2, 49, 53})
	r.RowDatas = append(r.RowDatas, ConvertNovelData2byte(tmp))

	// 回传成功写入的结果
	ret := new(Result)
	ret.Status = 2
	ret.InsertID = 0
	ret.AffectedRows = 1

	// 正常回传
	return ret, nil
}

// InsertSixteenthNovelResult 函式 🧚 为插入第十六本小说到模拟数据库 施公案 (会分配到 Slice-0)
func (r *Result) InsertSixteenthNovelResult() (*Result, error) {
	// 新增 Value
	tmp := make([]interface{}, 7)
	tmp[0] = int64(16)
	tmp[1] = int64(9789574927913)
	tmp[2] = "A Collection Of Shi"
	tmp[3] = "Anonymous"
	tmp[4] = 1850
	tmp[5] = "History"
	tmp[6] = int64(16)
	r.Values = append(r.Values, tmp)

	// 新增 RowDatas
	// r.RowDatas = append(r.RowDatas, []uint8{2, 49, 54, 13, 57, 55, 56, 57, 53, 55, 52, 57, 50, 55, 57, 49, 51, 19, 65, 32, 67, 111, 108, 108, 101, 99, 116, 105, 111, 110, 32, 79, 102, 32, 83, 104, 105, 9, 65, 110, 111, 110, 121, 109, 111, 117, 115, 4, 49, 56, 53, 48, 7, 72, 105, 115, 116, 111, 114, 121, 2, 49, 54})
	r.RowDatas = append(r.RowDatas, ConvertNovelData2byte(tmp))

	// 回传成功写入的结果
	ret := new(Result)
	ret.Status = 2
	ret.InsertID = 0
	ret.AffectedRows = 1

	// 正常回传
	return ret, nil
}

// InsertSeventeenthNovelResult 函式 🧚 为插入第十七本小说到模拟数据库 青楼梦 (会分配到 Slice-1)
func (r *Result) InsertSeventeenthNovelResult() (*Result, error) {
	// 新增 Value
	tmp := make([]interface{}, 7)
	tmp[0] = int64(17)
	tmp[1] = int64(9787533303396)
	tmp[2] = "Dream Of The Green Chamber"
	tmp[3] = "Yuda"
	tmp[4] = 1878
	tmp[5] = "Family Saga"
	tmp[6] = int64(17)
	r.Values = append(r.Values, tmp)

	// 新增 RowDatas
	// r.RowDatas = append(r.RowDatas, []uint8{2, 49, 55, 13, 57, 55, 56, 55, 53, 51, 51, 51, 48, 51, 51, 57, 54, 26, 68, 114, 101, 97, 109, 32, 79, 102, 32, 84, 104, 101, 32, 71, 114, 101, 101, 110, 32, 67, 104, 97, 109, 98, 101, 114, 4, 89, 117, 100, 97, 4, 49, 56, 55, 56, 11, 70, 97, 109, 105, 108, 121, 32, 83, 97, 103, 97, 2, 49, 55})
	r.RowDatas = append(r.RowDatas, ConvertNovelData2byte(tmp))

	// 回传成功写入的结果
	ret := new(Result)
	ret.Status = 2
	ret.InsertID = 0
	ret.AffectedRows = 1

	// 正常回传
	return ret, nil
}

// InsertEighteenthNovelResult 函式 🧚 为插入第十八本小说到模拟数据库 歧路灯 (会分配到 Slice-0)
func (r *Result) InsertEighteenthNovelResult() (*Result, error) {
	// 新增 Value
	tmp := make([]interface{}, 7)
	tmp[0] = int64(18)
	tmp[1] = int64(9787510434341)
	tmp[2] = "Lamp In The Side Street"
	tmp[3] = "Li Luyuan"
	tmp[4] = 1790
	tmp[5] = "Unofficial History"
	tmp[6] = int64(18)
	r.Values = append(r.Values, tmp)

	// 新增 RowDatas
	// r.RowDatas = append(r.RowDatas, []uint8{2, 49, 56, 13, 57, 55, 56, 55, 53, 49, 48, 52, 51, 52, 51, 52, 49, 23, 76, 97, 109, 112, 32, 73, 110, 32, 84, 104, 101, 32, 83, 105, 100, 101, 32, 83, 116, 114, 101, 101, 116, 9, 76, 105, 32, 76, 117, 121, 117, 97, 110, 4, 49, 55, 57, 48, 18, 85, 110, 111, 102, 102, 105, 99, 105, 97, 108, 32, 72, 105, 115, 116, 111, 114, 121, 2, 49, 56})
	r.RowDatas = append(r.RowDatas, ConvertNovelData2byte(tmp))

	// 回传成功写入的结果
	ret := new(Result)
	ret.Status = 2
	ret.InsertID = 0
	ret.AffectedRows = 1

	// 正常回传
	return ret, nil
}

// InsertNineteenthNovelResult 函式 🧚 为插入第十九本小说到模拟数据库 老残游记 (会分配到 Slice-1)
func (r *Result) InsertNineteenthNovelResult() (*Result, error) {
	// 新增 Value
	tmp := make([]interface{}, 7)
	tmp[0] = int64(19)
	tmp[1] = int64(9789571447469)
	tmp[2] = "The Travels of Lao Can"
	tmp[3] = "Liu E"
	tmp[4] = 1907
	tmp[5] = "Social Story"
	tmp[6] = int64(19)
	r.Values = append(r.Values, tmp)

	// 新增 RowDatas
	// r.RowDatas = append(r.RowDatas, []uint8{2, 49, 57, 13, 57, 55, 56, 57, 53, 55, 49, 52, 52, 55, 52, 54, 57, 22, 84, 104, 101, 32, 84, 114, 97, 118, 101, 108, 115, 32, 111, 102, 32, 76, 97, 111, 32, 67, 97, 110, 5, 76, 105, 117, 32, 69, 4, 49, 57, 48, 55, 12, 83, 111, 99, 105, 97, 108, 32, 83, 116, 111, 114, 121, 2, 49, 57})
	r.RowDatas = append(r.RowDatas, ConvertNovelData2byte(tmp))

	// 回传成功写入的结果
	ret := new(Result)
	ret.Status = 2
	ret.InsertID = 0
	ret.AffectedRows = 1

	// 正常回传
	return ret, nil
}

// InsertTwentiethNovelResult 函式 🧚 为插入第二十本小说到模拟数据库 二十年目睹之怪现状 (会分配到 Slice-0)
func (r *Result) InsertTwentiethNovelResult() (*Result, error) {
	// 新增 Value
	tmp := make([]interface{}, 7)
	tmp[0] = int64(20)
	tmp[1] = int64(9789571470047)
	tmp[2] = "Bizarre Happenings Eyewitnessed over Two Decades"
	tmp[3] = "Jianren Wu"
	tmp[4] = 1905
	tmp[5] = "Unofficial History"
	tmp[6] = int64(20)
	r.Values = append(r.Values, tmp)

	// 新增 RowDatas
	// r.RowDatas = append(r.RowDatas, []uint8{2, 50, 48, 13, 57, 55, 56, 57, 53, 55, 49, 52, 55, 48, 48, 52, 55, 48, 66, 105, 122, 97, 114, 114, 101, 32, 72, 97, 112, 112, 101, 110, 105, 110, 103, 115, 32, 69, 121, 101, 119, 105, 116, 110, 101, 115, 115, 101, 100, 32, 111, 118, 101, 114, 32, 84, 119, 111, 32, 68, 101, 99, 97, 100, 101, 115, 10, 74, 105, 97, 110, 114, 101, 110, 32, 87, 117, 4, 49, 57, 48, 53, 18, 85, 110, 111, 102, 102, 105, 99, 105, 97, 108, 32, 72, 105, 115, 116, 111, 114, 121, 2, 50, 48})
	r.RowDatas = append(r.RowDatas, ConvertNovelData2byte(tmp))

	// 回传成功写入的结果
	ret := new(Result)
	ret.Status = 2
	ret.InsertID = 0
	ret.AffectedRows = 1

	// 正常回传
	return ret, nil
}

// InsertTwentyFirstNovelResult 函式 🧚 为插入第二十一本小说到模拟数据库 孽海花 (会分配到 Slice-1)
func (r *Result) InsertTwentyFirstNovelResult() (*Result, error) {
	// 新增 Value
	tmp := make([]interface{}, 7)
	tmp[0] = int64(21)
	tmp[1] = int64(9787101097580)
	tmp[2] = "A Flower In A Sinful Sea"
	tmp[3] = "Zeng Pu"
	tmp[4] = 1904
	tmp[5] = "History"
	tmp[6] = int64(21)
	r.Values = append(r.Values, tmp)

	// 新增 RowDatas
	// r.RowDatas = append(r.RowDatas, []uint8{2, 50, 49, 13, 57, 55, 56, 55, 49, 48, 49, 48, 57, 55, 53, 56, 48, 24, 65, 32, 70, 108, 111, 119, 101, 114, 32, 73, 110, 32, 65, 32, 83, 105, 110, 102, 117, 108, 32, 83, 101, 97, 7, 90, 101, 110, 103, 32, 80, 117, 4, 49, 57, 48, 52, 7, 72, 105, 115, 116, 111, 114, 121, 2, 50, 49})
	r.RowDatas = append(r.RowDatas, ConvertNovelData2byte(tmp))

	// 回传成功写入的结果
	ret := new(Result)
	ret.Status = 2
	ret.InsertID = 0
	ret.AffectedRows = 1

	// 正常回传
	return ret, nil
}

// InsertTwentySecondNovelResult 函式 🧚 为插入第二十二本小说到模拟数据库 官场现形记 (会分配到 Slice-0)
func (r *Result) InsertTwentySecondNovelResult() (*Result, error) {
	// 新增 Value
	tmp := make([]interface{}, 7)
	tmp[0] = int64(22)
	tmp[1] = int64(9789861674193)
	tmp[2] = "Officialdom Unmasked"
	tmp[3] = "Li Baojia"
	tmp[4] = 1903
	tmp[5] = "Unofficial History"
	tmp[6] = int64(22)
	r.Values = append(r.Values, tmp)

	// 新增 RowDatas
	// r.RowDatas = append(r.RowDatas, []uint8{2, 50, 50, 13, 57, 55, 56, 57, 56, 54, 49, 54, 55, 52, 49, 57, 51, 20, 79, 102, 102, 105, 99, 105, 97, 108, 100, 111, 109, 32, 85, 110, 109, 97, 115, 107, 101, 100, 9, 76, 105, 32, 66, 97, 111, 106, 105, 97, 4, 49, 57, 48, 51, 18, 85, 110, 111, 102, 102, 105, 99, 105, 97, 108, 32, 72, 105, 115, 116, 111, 114, 121, 2, 50, 50})
	r.RowDatas = append(r.RowDatas, ConvertNovelData2byte(tmp))

	// 回传成功写入的结果
	ret := new(Result)
	ret.Status = 2
	ret.InsertID = 0
	ret.AffectedRows = 1

	// 正常回传
	return ret, nil
}

// InsertTwentyThirdNovelResult 函式 🧚 为插入第二十三本小说到模拟数据库 觉世名言十二楼 (会分配到 Slice-1)
func (r *Result) InsertTwentyThirdNovelResult() (*Result, error) {
	// 新增 Value
	tmp := make([]interface{}, 7)
	tmp[0] = int64(23)
	tmp[1] = int64(9787805469836)
	tmp[2] = "Tower For The Summer Heat"
	tmp[3] = "Li Yu"
	tmp[4] = 1680
	tmp[5] = "Unofficial History"
	tmp[6] = int64(23)
	r.Values = append(r.Values, tmp)

	// 新增 RowDatas
	// r.RowDatas = append(r.RowDatas, []uint8{2, 50, 51, 13, 57, 55, 56, 55, 56, 48, 53, 52, 54, 57, 56, 51, 54, 25, 84, 111, 119, 101, 114, 32, 70, 111, 114, 32, 84, 104, 101, 32, 83, 117, 109, 109, 101, 114, 32, 72, 101, 97, 116, 5, 76, 105, 32, 89, 117, 4, 49, 54, 56, 48, 18, 85, 110, 111, 102, 102, 105, 99, 105, 97, 108, 32, 72, 105, 115, 116, 111, 114, 121, 2, 50, 51})
	r.RowDatas = append(r.RowDatas, ConvertNovelData2byte(tmp))

	// 回传成功写入的结果
	ret := new(Result)
	ret.Status = 2
	ret.InsertID = 0
	ret.AffectedRows = 1

	// 正常回传
	return ret, nil
}

// InsertTwentyFourthNovelResult 函式 🧚 为插入第二十四本小说到模拟数据库 无声戏 (会分配到 Slice-0)
func (r *Result) InsertTwentyFourthNovelResult() (*Result, error) {
	// 新增 Value
	tmp := make([]interface{}, 7)
	tmp[0] = int64(24)
	tmp[1] = int64(9787508067247)
	tmp[2] = "Silent Operas"
	tmp[3] = "Li Yu"
	tmp[4] = 1680
	tmp[5] = "Unofficial History"
	tmp[6] = int64(24)
	r.Values = append(r.Values, tmp)

	// 新增 RowDatas
	// r.RowDatas = append(r.RowDatas, []uint8{2, 50, 52, 13, 57, 55, 56, 55, 53, 48, 56, 48, 54, 55, 50, 52, 55, 13, 83, 105, 108, 101, 110, 116, 32, 79, 112, 101, 114, 97, 115, 5, 76, 105, 32, 89, 117, 4, 49, 54, 56, 48, 12, 83, 111, 99, 105, 97, 108, 32, 83, 116, 111, 114, 121, 2, 50, 52})
	r.RowDatas = append(r.RowDatas, ConvertNovelData2byte(tmp))

	// 回传成功写入的结果
	ret := new(Result)
	ret.Status = 2
	ret.InsertID = 0
	ret.AffectedRows = 1

	// 正常回传
	return ret, nil
}

// InsertTwentyFifthNovelResult 函式 🧚 为插入第二十五本小说到模拟数据库 肉蒲团 (会分配到 Slice-1)
func (r *Result) InsertTwentyFifthNovelResult() (*Result, error) {
	// 新增 Value
	tmp := make([]interface{}, 7)
	tmp[0] = int64(25)
	tmp[1] = int64(9789573609049)
	tmp[2] = "The Carnal Prayer Mat"
	tmp[3] = "Li Yu"
	tmp[4] = 1680
	tmp[5] = "Unofficial History"
	tmp[6] = int64(25)
	r.Values = append(r.Values, tmp)

	// 新增 RowDatas
	// r.RowDatas = append(r.RowDatas, []uint8{2, 50, 53, 13, 57, 55, 56, 57, 53, 55, 51, 54, 48, 57, 48, 52, 57, 21, 84, 104, 101, 32, 67, 97, 114, 110, 97, 108, 32, 80, 114, 97, 121, 101, 114, 32, 77, 97, 116, 5, 76, 105, 32, 89, 117, 4, 49, 54, 56, 48, 12, 83, 111, 99, 105, 97, 108, 32, 83, 116, 111, 114, 121, 2, 50, 53})
	r.RowDatas = append(r.RowDatas, ConvertNovelData2byte(tmp))

	// 回传成功写入的结果
	ret := new(Result)
	ret.Status = 2
	ret.InsertID = 0
	ret.AffectedRows = 1

	// 正常回传
	return ret, nil
}

// InsertTwentySixthNovelResult 函式 🧚 为插入第二十六本小说到模拟数据库 浮生六记 (会分配到 Slice-0)
func (r *Result) InsertTwentySixthNovelResult() (*Result, error) {
	// 新增 Value
	tmp := make([]interface{}, 7)
	tmp[0] = int64(26)
	tmp[1] = int64(9787533948108)
	tmp[2] = "Six Records Of A Floating Life"
	tmp[3] = "Shen Fu"
	tmp[4] = 1878
	tmp[5] = "Autobiography"
	tmp[6] = int64(26)
	r.Values = append(r.Values, tmp)

	// 新增 RowDatas
	// r.RowDatas = append(r.RowDatas, []uint8{2, 50, 54, 13, 57, 55, 56, 55, 53, 51, 51, 57, 52, 56, 49, 48, 56, 30, 83, 105, 120, 32, 82, 101, 99, 111, 114, 100, 115, 32, 79, 102, 32, 65, 32, 70, 108, 111, 97, 116, 105, 110, 103, 32, 76, 105, 102, 101, 7, 83, 104, 101, 110, 32, 70, 117, 4, 49, 56, 55, 56, 13, 65, 117, 116, 111, 98, 105, 111, 103, 114, 97, 112, 104, 121, 2, 50, 54})
	r.RowDatas = append(r.RowDatas, ConvertNovelData2byte(tmp))

	// 回传成功写入的结果
	ret := new(Result)
	ret.Status = 2
	ret.InsertID = 0
	ret.AffectedRows = 1

	// 正常回传
	return ret, nil
}

// InsertTwentySeventhNovelResult 函式 🧚 为插入第二十七本小说到模拟数据库 野叟曝言 (会分配到 Slice-1)
func (r *Result) InsertTwentySeventhNovelResult() (*Result, error) {
	// 新增 Value
	tmp := make([]interface{}, 7)
	tmp[0] = int64(27)
	tmp[1] = int64(9786666141110)
	tmp[2] = "Humble Words Of A Rustic Elder"
	tmp[3] = "Xia Jingqu"
	tmp[4] = 1787
	tmp[5] = "Historical fiction"
	tmp[6] = int64(27)
	r.Values = append(r.Values, tmp)

	// 新增 RowDatas
	// r.RowDatas = append(r.RowDatas, []uint8{2, 50, 55, 13, 57, 55, 56, 54, 54, 54, 54, 49, 52, 49, 49, 49, 48, 30, 72, 117, 109, 98, 108, 101, 32, 87, 111, 114, 100, 115, 32, 79, 102, 32, 65, 32, 82, 117, 115, 116, 105, 99, 32, 69, 108, 100, 101, 114, 10, 88, 105, 97, 32, 74, 105, 110, 103, 113, 117, 4, 49, 55, 56, 55, 18, 72, 105, 115, 116, 111, 114, 105, 99, 97, 108, 32, 102, 105, 99, 116, 105, 111, 110, 2, 50, 55})
	r.RowDatas = append(r.RowDatas, ConvertNovelData2byte(tmp))

	// 回传成功写入的结果
	ret := new(Result)
	ret.Status = 2
	ret.InsertID = 0
	ret.AffectedRows = 1

	// 正常回传
	return ret, nil
}

// InsertTwentyEighthNovelResult 函式 🧚 为插入第二十八本小说到模拟数据库 九尾龟 (会分配到 Slice-0)
func (r *Result) InsertTwentyEighthNovelResult() (*Result, error) {
	// 新增 Value
	tmp := make([]interface{}, 7)
	tmp[0] = int64(28)
	tmp[1] = int64(9789571435473)
	tmp[2] = "Nine-Tailed Turtle"
	tmp[3] = "Lu Can"
	tmp[4] = 1551
	tmp[5] = "Mythology"
	tmp[6] = int64(28)
	r.Values = append(r.Values, tmp)

	// 新增 RowDatas
	// r.RowDatas = append(r.RowDatas, []uint8{2, 50, 56, 13, 57, 55, 56, 57, 53, 55, 49, 52, 51, 53, 52, 55, 51, 18, 78, 105, 110, 101, 45, 84, 97, 105, 108, 101, 100, 32, 84, 117, 114, 116, 108, 101, 6, 76, 117, 32, 67, 97, 110, 4, 49, 53, 53, 49, 9, 77, 121, 116, 104, 111, 108, 111, 103, 121, 2, 50, 56})
	r.RowDatas = append(r.RowDatas, ConvertNovelData2byte(tmp))

	// 回传成功写入的结果
	ret := new(Result)
	ret.Status = 2
	ret.InsertID = 0
	ret.AffectedRows = 1

	// 正常回传
	return ret, nil
}

// InsertTwentyNinethNovelResult 函式 🧚 为插入第二十九本小说到模拟数据库 品花宝鉴 (会分配到 Slice-1)
func (r *Result) InsertTwentyNinethNovelResult() (*Result, error) {
	// 新增 Value
	tmp := make([]interface{}, 7)
	tmp[0] = int64(29)
	tmp[1] = int64(9789866318603)
	tmp[2] = "A History Of Floral Treasures"
	tmp[3] = "Chen Sen"
	tmp[4] = 1849
	tmp[5] = "Romance"
	tmp[6] = int64(29)
	r.Values = append(r.Values, tmp)

	// 新增 RowDatas
	// r.RowDatas = append(r.RowDatas, []uint8{2, 50, 57, 13, 57, 55, 56, 57, 56, 54, 54, 51, 49, 56, 54, 48, 51, 29, 65, 32, 72, 105, 115, 116, 111, 114, 121, 32, 79, 102, 32, 70, 108, 111, 114, 97, 108, 32, 84, 114, 101, 97, 115, 117, 114, 101, 115, 8, 67, 104, 101, 110, 32, 83, 101, 110, 4, 49, 56, 52, 57, 7, 82, 111, 109, 97, 110, 99, 101, 2, 50, 57})
	r.RowDatas = append(r.RowDatas, ConvertNovelData2byte(tmp))

	// 回传成功写入的结果
	ret := new(Result)
	ret.Status = 2
	ret.InsertID = 0
	ret.AffectedRows = 1

	// 正常回传
	return ret, nil
}
