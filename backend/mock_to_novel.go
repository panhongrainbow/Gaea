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
func (fdb *fakeDB) switchNovelResult(key uint32) (*mysql.Result, error) {
	switch key {
	case 4290409450:
		/*
			所对应各切片 SQL 执行字串 以及 切片相关资讯
			数据库名称: novel
			查询模拟数据库的网路位置: 192.168.122.2:3313
			数据库执行字串: INSERT INTO `novel`.`Book_0001` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (1,9781517191276,'Romance Of The Three Kingdoms','Luo Guanzhong',1522,'Historical fiction')
			数据库执行时所对应的 Key: 4290409450
		*/
		fmt.Println("切片長度", len(fakeDBInstance["novel"].MockDataInDB))
		tmp0 := mysql.MakeNovelEmptyResult()
		fakeDBInstance["novel"].MockDataInDB[0] = tmp0
		tmp1 := mysql.MakeNovelEmptyResult()
		fakeDBInstance["novel"].MockDataInDB[1] = tmp1
		_ = fakeDBInstance["novel"].MockDataInDB[1].InsertFirstNovelResult()
		fmt.Println("第一个切片", fakeDBInstance["novel"].MockDataInDB[0])
		fmt.Println("第二个切片", fakeDBInstance["novel"].MockDataInDB[1])
		fmt.Println("命中 4290409450")
		return &tmp0, nil
	case 3124618913:
		test := fakeDBInstance
		fmt.Println(test)
		fmt.Println("命中 3124618913")
		tmp := mysql.MakeNovelEmptyResult()
		return &tmp, nil
	case 1260331735:
		/*
			所对应各切片 SQL 执行字串 以及 切片相关资讯
			数据库名称: novel
			查询模拟数据库的网路位置: 192.168.122.2:3313
			数据库执行字串: SELECT *,`BookID` FROM `novel`.`Book_0001` ORDER BY `BookID`
			数据库执行时所对应的 Key: 1260331735
		*/
		test := fakeDBInstance
		fmt.Println(test)
		fmt.Println("命中 1260331735")
		tmp := mysql.MakeNovelEmptyResult()
		return &tmp, nil
	case 1196547673:
		/*
			所对应各切片 SQL 执行字串 以及 切片相关资讯
			数据库名称: novel
			查询模拟数据库的网路位置: 192.168.122.2:3310
			数据库执行字串: SELECT *,`BookID` FROM `novel`.`Book_0000` ORDER BY `BookID`
			数据库执行时所对应的 Key: 1196547673
		*/
		test := fakeDBInstance
		fmt.Println(test)
		fmt.Println("命中 1196547673")
		tmp := mysql.MakeNovelEmptyResult()
		return &tmp, nil
	case 1401931444:
		/*
			所对应各切片 SQL 执行字串 以及 切片相关资讯
			数据库名称: novel
			查询模拟数据库的网路位置: 192.168.122.2:3314
			数据库执行字串: SELECT *,`BookID` FROM `novel`.`Book_0001` ORDER BY `BookID`
			数据库执行时所对应的 Key: 1401931444
		*/
		test := fakeDBInstance
		fmt.Println(test)
		fmt.Println("命中 1401931444")
		tmp := mysql.MakeNovelEmptyResult()
		return &tmp, nil
	}
	tmp := mysql.MakeNovelEmptyResult()
	log.Fatal("没有命中模拟测试 key 为: ", key) // 中断，因为测试程式有问题
	return &tmp, nil
}
