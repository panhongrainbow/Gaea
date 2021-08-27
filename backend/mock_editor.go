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

// MarkLoaded 函式 🧚 为 标记载入资料完成
func (basicLoad) MarkLoaded() {
	FakeDBInstance.Loaded = true // 载入资料完成
}

// LoadData 函式 🧚 为 载入一些测试资料
func (basicLoad) LoadData() {
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
	fmt.Println("载入测试资料序号", key)
}
