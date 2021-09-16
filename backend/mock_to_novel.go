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
	case 1260331735:
		/*
			所对应各切片 SQL 执行字串 以及 切片相关资讯
			数据库名称: novel
			模拟数据库的网路位置: 192.168.122.2:3313
			数据库执行字串: SELECT *,`BookID` FROM `novel`.`Book_0001` ORDER BY `BookID`
			数据库执行时所对应的 Key: 1260331735
		*/
		fmt.Println("命中 1260331735")
		return mysql.SelectNovelResult(), nil
	case 1196547673:
		/*
			所对应各切片 SQL 执行字串 以及 切片相关资讯
			数据库名称: novel
			模拟数据库的网路位置: 192.168.122.2:3310
			数据库执行字串: SELECT *,`BookID` FROM `novel`.`Book_0000` ORDER BY `BookID`
			数据库执行时所对应的 Key: 1196547673
		*/
		fmt.Println("命中 1196547673")
		return &fakeDBInstance["novel"].MockDataInDB[0], nil
		// return mysql.SelectNovelResult(), nil
	case 1401931444:
		/*
			所对应各切片 SQL 执行字串 以及 切片相关资讯
			数据库名称: novel
			模拟数据库的网路位置: 192.168.122.2:3314
			数据库执行字串: SELECT *,`BookID` FROM `novel`.`Book_0001` ORDER BY `BookID`
			数据库执行时所对应的 Key: 1401931444
		*/
		fmt.Println("命中 1401931444")
		return mysql.SelectNovelResult(), nil
	}
	log.Fatal("没有命中模拟测试 key 为: ", key) // 中断，因为测试程式有问题
	return nil, nil
}
