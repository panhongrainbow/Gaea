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
	case 3124618913:
		fmt.Println("命中 3124618913")
		return mysql.SelectnovelResult(), nil
	case 1260331735:
		fmt.Println("命中 1260331735")
		return mysql.SelectnovelResult(), nil
	case 1196547673:
		fmt.Println("命中 1196547673")
		return mysql.SelectnovelResult(), nil
	}
	log.Fatal("没有命中模拟测试 key 为: ", key) // 中断，因为测试程式有问题
	return &mysql.Result{}, nil
}
