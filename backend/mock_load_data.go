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
	return FakeDBInstance.Loaded // 回传载入资料是否完成
}

// MarkLoaded 函式 🧚 为 标记载入资料完成
func (basicLoad) MarkLoaded() {
	FakeDBInstance.Loaded = true // 载入资料完成
}

// UnMarkLoaded 函式 🧚 为 去除 载入资料完成 的标记
func (basicLoad) UnMarkLoaded() {
	FakeDBInstance.Loaded = false // 去除 载入资料完成 的标记
}

// LoadData 函式 🧚 为 载入一些测试资料
func (basicLoad) LoadData() error {
	// 编写测试资料
	data := SubFakeDB{
		addr:     "192.168.122.2:3307",
		user:     "panhong",
		password: "12345",
		sql:      "SELECT * FROM `Library`.`Book`",
		result:   *mysql.SelectLibrayResult(),
	}

	// 载入测试资料
	FakeDBInstance.MockResult = make(map[uint32]mysql.Result)
	key := FakeDBInstance.MakeMockResult(data)

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
	FakeDBInstance.MockResult = nil
	return nil
}
