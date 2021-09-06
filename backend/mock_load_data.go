package backend

import (
	"fmt"
	"github.com/XiaoMi/Gaea/mysql"
)

// >>>>> >>>>> >>>>> >>>>> >>>>> 第一个简单的载入测试资料方法

// basicLoad 资料 🧚 为 是用来引出载入测资料的变数
type basicLoad struct {
	// 空资料
}

// IsLoaded 函式 🧚 为 确认是否载入资料完成
func (basicLoad) IsLoaded() bool {
	fakeDBInstance.Lock()
	defer fakeDBInstance.Unlock()
	return fakeDBInstance.Loaded // 回传载入资料是否完成
}

// MarkLoaded 函式 🧚 为 标记载入资料完成
func (basicLoad) MarkLoaded() {
	fakeDBInstance.Loaded = true // 载入资料完成
}

// UnMarkLoaded 函式 🧚 为 去除 载入资料完成 的标记
func (basicLoad) UnMarkLoaded() {
	fakeDBInstance.Loaded = false // 去除 载入资料完成 的标记
}

// LoadData 函式 🧚 为 载入一些测试资料
func (basicLoad) LoadData() error {
	// 编写测试资料
	data := subFakeDB{
		addr:     "192.168.122.2:3307",
		user:     "panhong",
		password: "12345",
		sql:      "SELECT * FROM `novel`.`Book`",
		result:   *mysql.SelectnovelResult(),
	}

	// 载入测试资料
	fakeDBInstance.MockResult = make(map[uint32]mysql.Result)
	key := fakeDBInstance.MakeMockResult(data)

	// 显示测试资料序号并回传 nil
	fmt.Printf("\u001B[35m 载入测试资料序号 Key: %d\n", key)

	// 编写测试资料
	data = subFakeDB{
		addr:     "192.168.122.2:3306",
		user:     "panhong",
		password: "12345",
		sql:      "SELECT * FROM `novel`.`Book_0000`",
		result:   *mysql.SelectnovelResult(),
	}

	// 载入测试资料
	// fakeDBInstance.MockResult = make(map[uint32]mysql.Result)
	key = fakeDBInstance.MakeMockResult(data)

	// 显示测试资料序号并回传 nil
	fmt.Printf("\u001B[35m 载入测试资料序号 Key: %d\n", key)

	// 编写测试资料
	data = subFakeDB{
		addr:     "192.168.122.2:3307",
		user:     "panhong",
		password: "12345",
		sql:      "SELECT * FROM `novel`.`Book_0000`",
		result:   *mysql.SelectnovelResult(),
	}

	// 载入测试资料
	// fakeDBInstance.MockResult = make(map[uint32]mysql.Result)
	key = fakeDBInstance.MakeMockResult(data)

	// 显示测试资料序号并回传 nil
	fmt.Printf("\u001B[35m 载入测试资料序号 Key: %d\n", key)

	// 编写测试资料
	data = subFakeDB{
		addr:     "192.168.122.2:3308",
		user:     "panhong",
		password: "12345",
		sql:      "SELECT * FROM `novel`.`Book_0000`",
		result:   *mysql.SelectnovelResult(),
	}

	// 载入测试资料
	// fakeDBInstance.MockResult = make(map[uint32]mysql.Result)
	key = fakeDBInstance.MakeMockResult(data)

	// 显示测试资料序号并回传 nil
	fmt.Printf("\u001B[35m 载入测试资料序号 Key: %d\n", key)

	// 编写测试资料
	data = subFakeDB{
		addr:     "192.168.122.2:3309",
		user:     "panhong",
		password: "12345",
		sql:      "SELECT * FROM `novel`.`Book_0001`",
		result:   *mysql.SelectnovelResult(),
	}

	// 载入测试资料
	// fakeDBInstance.MockResult = make(map[uint32]mysql.Result)
	key = fakeDBInstance.MakeMockResult(data)

	// 显示测试资料序号并回传 nil
	fmt.Printf("\u001B[35m 载入测试资料序号 Key: %d\n", key)

	// 编写测试资料
	data = subFakeDB{
		addr:     "192.168.122.2:3310",
		user:     "panhong",
		password: "12345",
		sql:      "SELECT * FROM `novel`.`Book_0001`",
		result:   *mysql.SelectnovelResult(),
	}

	// 载入测试资料
	// fakeDBInstance.MockResult = make(map[uint32]mysql.Result)
	key = fakeDBInstance.MakeMockResult(data)

	// 显示测试资料序号并回传 nil
	fmt.Printf("\u001B[35m 载入测试资料序号 Key: %d\n", key)

	// 编写测试资料
	data = subFakeDB{
		addr:     "192.168.122.2:3311",
		user:     "panhong",
		password: "12345",
		sql:      "SELECT * FROM `novel`.`Book_0001`",
		result:   *mysql.SelectnovelResult(),
	}

	// 载入测试资料
	// fakeDBInstance.MockResult = make(map[uint32]mysql.Result)
	key = fakeDBInstance.MakeMockResult(data)

	// 显示测试资料序号并回传 nil
	fmt.Printf("\u001B[35m 载入测试资料序号 Key: %d\n", key)

	return nil
}

// EmptyData 函式 🧚 为 清空已载入的测试资料
// 在大部份的测试状况下，会先载入特定的测试资料
// 进行一连串的测试后，才会再换载入新的测试资料
// 所以已载入的测试资料就全部清除，不需要考虑一笔一笔去移除
func (basicLoad) EmptyData() error {
	// 清空载入测试资料
	fakeDBInstance.MockResult = nil
	return nil
}

// Lock 和 UnLock
/* 函式目前只用在
   1 确认单元测试资料是否正常载入
   在函式 IsLoaded() 可以找到新增上解锁的机制
   2 载入单元测试资料时
   在函式 NewDirectConnection 可以找到新增上解锁的机制
*/
// Lock 函式 🧚 上锁函式
func (basicLoad) Lock() {
	fakeDBInstance.Lock()
}

// UnLock 函式 🧚 解锁函式
func (basicLoad) UnLock() {
	fakeDBInstance.Unlock()
}
