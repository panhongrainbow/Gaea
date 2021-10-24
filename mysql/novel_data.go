package mysql

import (
	"strconv"
)

// >>>>> >>>>> >>>>> >>>>> >>>>> 向数据库查询 29 本小说的回传结果

// MakeNovelFieldData 函式 🧚 是用来产生初始回传数据库的栏位资料
func MakeNovelFieldData(tableSlice string) (*Result, error) {
	// 初始化数据库回传数值
	ret := new(Result)
	ret.Status = 34
	ret.InsertID = 0
	ret.AffectedRows = 0

	fdTest := make([]fieldTestData, 6)
	field := make([]Field, 7)

	fdTest[0] = fieldTestData{
		name:         "BookID",
		orgName:      "BookID",
		charset:      63,
		columnLength: 11,
		fieldType:    3,
		flag:         20483,
	}
	fdTest[1] = fieldTestData{
		name:         "Isbn",
		orgName:      "Isbn",
		charset:      63,
		columnLength: 50,
		fieldType:    8,
		flag:         4097,
	}
	fdTest[2] = fieldTestData{
		name:         "Title",
		orgName:      "Title",
		charset:      33,
		columnLength: 300,
		fieldType:    253,
		flag:         4097,
	}
	fdTest[3] = fieldTestData{
		name:         "Author",
		orgName:      "Author",
		charset:      33,
		columnLength: 90,
		fieldType:    253,
		flag:         0,
	}
	fdTest[4] = fieldTestData{
		name:         "Publish",
		orgName:      "Publish",
		charset:      63,
		columnLength: 4,
		fieldType:    3,
		flag:         0,
	}
	fdTest[5] = fieldTestData{
		name:         "Category",
		orgName:      "Category",
		charset:      33,
		columnLength: 90,
		fieldType:    253,
		flag:         4097,
	}

	resultset := Resultset{}
	ret.Resultset = &resultset
	ret.Resultset.Fields = []*Field{}

	for i := 0; i < 6; i++ {
		fdTest[i].def = "def"
		fdTest[i].schema = "novel"
		fdTest[i].table = tableSlice
		fdTest[i].orgTable = tableSlice
		field[i].ConvertFieldTest2Field(fdTest[i])
		ret.Resultset.Fields = append(ret.Resultset.Fields, &field[i])
	}
	field[6].ConvertFieldTest2Field(fdTest[0])
	ret.Resultset.Fields = append(ret.Resultset.Fields, &field[6])

	ret.FieldNames = make(map[string]int)
	ret.FieldNames["BookID"] = 6
	ret.FieldNames["Isbn"] = 1
	ret.FieldNames["Title"] = 2
	ret.FieldNames["Author"] = 3
	ret.FieldNames["Publish"] = 4
	ret.FieldNames["Category"] = 5

	ret.Values = make([][]interface{}, 0, 29)
	ret.RowDatas = make([]RowData, 0, 29)

	return ret, nil
}

// MakeNovelFieldDataTmp 函式 🧚 是用来产生初始回传数据库的暂存栏位资料
func MakeNovelFieldDataTmp() (*Result, error) {
	// MakeNovelFieldData 函式内容精简
	ret := new(Result)
	resultset := Resultset{}
	ret.Resultset = &resultset
	ret.Values = make([][]interface{}, 0, 29)
	ret.RowDatas = make([]RowData, 0, 29)
	return ret, nil
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
