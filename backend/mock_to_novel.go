package backend

import (
	"fmt"
	"github.com/XiaoMi/Gaea/mysql"
	"log"
)

// switchMockResult 函式 🧚 为到不同的模拟数据库去找寻回应的讯息
func (fdb *fakeDB) switchMockResult(db string, key uint32) (*mysql.Result, error) {
	switch db {
	case "novel": // 29 本小说部份
		return fdb.switchNovelResult(key) // 在小说模拟数据库时去找对应到 SQL 字串的回应讯息
	}
	log.Fatal("没有命中模拟测试数据名称为: ", db) // 中断，因为测试程式有问题
	return &mysql.Result{}, nil
}

// switchNovelResult 函式 🧚 为在小说模拟数据库时去找对应到 SQL 字串的回应讯息
// 这里可以参考位于 Gaea/backend/mock_key_test.go 的测试函式 TestMockNovelKey
func (fdb *fakeDB) switchNovelResult(key uint32) (*mysql.Result, error) {
	switch key {
	// >>>>> >>>>> >>>>> >>>>> >>>>> 向多台数据库进行查询
	case 3717314451, 1196547673, 4270781616:
		/*
			向第一个切片进行查询
			所对应各切片 SQL 执行字串 以及 切片相关资讯
			数据库名称: novel
			模拟数据库的网路位置: 192.168.122.2:3309 或 192.168.122.2:3310 或 192.168.122.2:3311
			数据库执行字串: SELECT *,`BookID` FROM `novel`.`Book_0000` ORDER BY `BookID`
			数据库执行时所对应的 Key: 3717314451 或 1196547673 或 4270781616
		*/
		fmt.Printf("\u001B[35m 命中数据库所对应的 Key: %d\n", key)
		return fdb.MockDataInDB[0], nil
	case 2403537350, 1260331735, 1401931444:
		/*
			向第二个切片进行查询
			所对应各切片 SQL 执行字串 以及 切片相关资讯
			数据库名称: novel
			模拟数据库的网路位置: 192.168.122.2:3312 或 192.168.122.2:3313 或 192.168.122.2:3314
			数据库执行字串: SELECT *,`BookID` FROM `novel`.`Book_0001` ORDER BY `BookID`
			数据库执行时所对应的 Key: 2403537350 或 1260331735 或 1401931444
		*/
		fmt.Printf("\u001B[35m 命中数据库所对应的 Key: %d\n", key)
		return fdb.MockDataInDB[1], nil
	// >>>>> >>>>> >>>>> >>>>> >>>>> 向多台数据库进行写入
	case 1389454267:
		if err := fdb.MockDataInDB[1].InsertFirstNovelResult(); err != nil { // 写入第一本小说到数据库 三国演义
			return nil, err
		}
		fmt.Printf("\u001B[35m 命中数据库所对应的 Key: %d\n", key)
		return mysql.MakeNovelEmptyResult()
	case 514659115:
		if err := fdb.MockDataInDB[1].InsertThirdNovelResult(); err != nil { // 写入第三本小说到数据库 西游记
			return nil, err
		}
		fmt.Printf("\u001B[35m 命中数据库所对应的 Key: %d\n", key)
		return mysql.MakeNovelEmptyResult()
	case 4076192191:
		if err := fdb.MockDataInDB[1].InsertFifthNovelResult(); err != nil { // 写入第五本小说到数据库 金瓶梅
			return nil, err
		}
		fmt.Printf("\u001B[35m 命中数据库所对应的 Key: %d\n", key)
		return mysql.MakeNovelEmptyResult()
	case 1572904758:
		if err := fdb.MockDataInDB[1].InsertSeventhNovelResult(); err != nil { // 写入第七本小说到数据库 初刻拍案惊奇
			return nil, err
		}
		fmt.Printf("\u001B[35m 命中数据库所对应的 Key: %d\n", key)
		return mysql.MakeNovelEmptyResult()
	case 3188314210:
		if err := fdb.MockDataInDB[1].InsertNinethNovelResult(); err != nil { // 写入第九本小说到数据库 封神演义
			return nil, err
		}
		fmt.Printf("\u001B[35m 命中数据库所对应的 Key: %d\n", key)
		return mysql.MakeNovelEmptyResult()
	case 3599615497:
		if err := fdb.MockDataInDB[1].InsertEleventhNovelResult(); err != nil { // 写入第十一本小说到数据库 喻世明言
			return nil, err
		}
		fmt.Printf("\u001B[35m 命中数据库所对应的 Key: %d\n", key)
		return mysql.MakeNovelEmptyResult()
	case 709958148:
		if err := fdb.MockDataInDB[1].InsertThirteenthNovelResult(); err != nil { // 写入第十三本小说到数据库 杨家将
			return nil, err
		}
		fmt.Printf("\u001B[35m 命中数据库所对应的 Key: %d\n", key)
		return mysql.MakeNovelEmptyResult()
	case 56203336:
		if err := fdb.MockDataInDB[1].InsertFifteenthNovelResult(); err != nil { // 写入第十五本小说到数据库 七侠五义
			return nil, err
		}
		fmt.Printf("\u001B[35m 命中数据库所对应的 Key: %d\n", key)
		return mysql.MakeNovelEmptyResult()
	case 3821388015:
		if err := fdb.MockDataInDB[1].InsertSeventeenthNovelResult(); err != nil { // 写入第十七本小说到数据库 青楼梦
			return nil, err
		}
		fmt.Printf("\u001B[35m 命中数据库所对应的 Key: %d\n", key)
		return mysql.MakeNovelEmptyResult()
	case 398747927:
		if err := fdb.MockDataInDB[1].InsertNineteenthNovelResult(); err != nil { // 写入第十九本小说到数据库 老残游记
			return nil, err
		}
		fmt.Printf("\u001B[35m 命中数据库所对应的 Key: %d\n", key)
		return mysql.MakeNovelEmptyResult()
	case 1498815330:
		if err := fdb.MockDataInDB[1].InsertTwentyFirstNovelResult(); err != nil { // 写入第二十一本小说到数据库 孽海花
			return nil, err
		}
		fmt.Printf("\u001B[35m 命中数据库所对应的 Key: %d\n", key)
		return mysql.MakeNovelEmptyResult()
	case 2614046017:
		if err := fdb.MockDataInDB[1].InsertTwentyThirdNovelResult(); err != nil { // 写入第二十三本小说到数据库 觉世名言十二楼
			return nil, err
		}
		fmt.Printf("\u001B[35m 命中数据库所对应的 Key: %d\n", key)
		return mysql.MakeNovelEmptyResult()
	case 238477972:
		if err := fdb.MockDataInDB[1].InsertTwentyFifthNovelResult(); err != nil { // 写入第二十五本小说到数据库 肉蒲团
			return nil, err
		}
		fmt.Printf("\u001B[35m 命中数据库所对应的 Key: %d\n", key)
		return mysql.MakeNovelEmptyResult()
	case 2745523730:
		if err := fdb.MockDataInDB[1].InsertTwentySeventhNovelResult(); err != nil { // 写入第二十七本小说到数据库 野叟曝言
			return nil, err
		}
		fmt.Printf("\u001B[35m 命中数据库所对应的 Key: %d\n", key)
		return mysql.MakeNovelEmptyResult()
	case 424563096:
		if err := fdb.MockDataInDB[1].InsertTwentyNinethNovelResult(); err != nil { // 写入第二十九本小说到数据库 品花宝鉴
			return nil, err
		}
		fmt.Printf("\u001B[35m 命中数据库所对应的 Key: %d\n", key)
		return mysql.MakeNovelEmptyResult()
	}

	// >>>>> >>>>> >>>>> >>>>> >>>>> 如果没有命中 key 值的时候，就直接中断整个测试
	log.Fatalf("\u001B[35m 没有命中模拟测试 Key 为: %d\n", key) // 中断，因为测试程式有问题
	return nil, nil
}
