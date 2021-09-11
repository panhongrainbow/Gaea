package backend

import (
	"fmt"
	"github.com/XiaoMi/Gaea/mysql"
)

// >>>>> >>>>> >>>>> >>>>> >>>>> 载入 29 本小说测试资料方法

// >>>>> >>>>> >>>>> >>>>> >>>>> 以下方法不会修改到测试的资料

// novelData 资料 🧚 为 是用来引出载入测资料的变数
type novelData struct {
	dbName string
}

// GetDB 函式 🧚 为 在数据库池并不会传送资料库名称到下层函式，用此函式取出数据库名称
func (n *novelData) GetDB() string {
	return n.dbName
}

// UseDB 函式 🧚 为 在数据库池并不会传送资料库名称到下层函式，用此函式指定数据库名称
func (n *novelData) UseDB(db string) error {
	n.dbName = db
	return nil
}

// IsInited 函式 🧚 为 确认是否 初始化模拟数据库
func (n *novelData) IsInited() bool {
	fakeDBInstance[n.GetDB()].Lock()
	defer fakeDBInstance[n.GetDB()].Unlock()
	return fakeDBInstance[n.GetDB()].Loaded // 回传载入资料是否完成
}

// MarkInited 函式 🧚 为 标记 初始化模拟数据库 完成
func (n *novelData) MarkInited() {
	fakeDBInstance[n.GetDB()].Loaded = true // 载入资料完成
}

// UnMarkInited 函式 🧚 为 去除 初始化模拟数据库 的标记
func (n *novelData) UnMarkInited() {
	fakeDBInstance[n.GetDB()].Loaded = false // 去除 载入资料完成 的标记
}

// EmptyData 函式 🧚 为 清空已载入的测试资料
// 在大部份的测试状况下，会先载入特定的测试资料
// 进行一连串的测试后，才会再换载入新的测试资料
// 所以已载入的测试资料就全部清除，不需要考虑一笔一笔去移除
func (n *novelData) EmptyData() error {
	// 清空载入测试资料
	fakeDBInstance[n.GetDB()].MockReAct = nil
	return nil
}

// Lock 和 UnLock
/* 函式目前只用在
   1 确认单元测试资料是否正常载入
   在函式 IsInited() 可以找到新增上解锁的机制
   2 载入单元测试资料时
   在函式 NewDirectConnection 可以找到新增上解锁的机制
*/

// Lock 函式 🧚 上锁函式
func (n *novelData) Lock() {
	fakeDBInstance[n.GetDB()].Lock()
}

// UnLock 函式 🧚 解锁函式
func (n *novelData) UnLock() {
	fakeDBInstance[n.GetDB()].Unlock()
}

// LoadData 函式 🧚 为 载入一些测试资料
func (n *novelData) InitData() error {
	// 载入测试资料
	fakeDBInstance[n.GetDB()] = new(fakeDB)
	fakeDBInstance[n.GetDB()].MockReAct = make(map[uint32]mysql.Result)
	return nil
}

// LoadData 函式 🧚 为 载入一些测试资料
func (*novelData) LoadData() error {
	// 编写测试资料
	data := subFakeDB{
		addr:     "192.168.122.2:3307",
		user:     "panhong",
		password: "12345",
		sql:      "SELECT * FROM `novel`.`Book`",
		result:   *mysql.SelectnovelResult(),
	}

	// 载入测试资料
	fakeDBInstance[data.db] = new(fakeDB)
	fakeDBInstance[data.db].MockReAct = make(map[uint32]mysql.Result)
	key := fakeDBInstance[data.db].MakeMockResult(data)

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
	// fakeDBInstance.MockReAct = make(map[uint32]mysql.Result)
	key = fakeDBInstance[data.db].MakeMockResult(data)

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
	// fakeDBInstance.MockReAct = make(map[uint32]mysql.Result)
	key = fakeDBInstance[data.db].MakeMockResult(data)

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
	// fakeDBInstance.MockReAct = make(map[uint32]mysql.Result)
	key = fakeDBInstance[data.db].MakeMockResult(data)

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
	// fakeDBInstance.MockReAct = make(map[uint32]mysql.Result)
	key = fakeDBInstance[data.db].MakeMockResult(data)

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
	// fakeDBInstance.MockReAct = make(map[uint32]mysql.Result)
	key = fakeDBInstance[data.db].MakeMockResult(data)

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
	// fakeDBInstance.MockReAct = make(map[uint32]mysql.Result)
	key = fakeDBInstance[data.db].MakeMockResult(data)

	// 显示测试资料序号并回传 nil
	fmt.Printf("\u001B[35m 载入测试资料序号 Key: %d\n", key)

	return nil
}

// >>>>> >>>>> >>>>> >>>>> >>>>> 以下方法 会 修改到测试的资料

// InsertData 函式 🧚 会新增模拟数据库的内容
func (novelData) InsertData() {
	//
}

// UpdateData 函式 🧚 会修改模拟数据库的内容
func (novelData) UpdateData() {
	//
}

// DeleteData 函式 🧚 会删除模拟数据库的内容
func (novelData) DeleteData() {
	//
}
