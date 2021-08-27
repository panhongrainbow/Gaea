package backend

import (
	"github.com/XiaoMi/Gaea/mysql"
	"hash/fnv"
)

// TakeOver >>>>> >>>>> >>>>> >>>>> >>>>> 单元测试的指示灯
var TakeOver bool // 现在是否由单元测试接管

// FakeDB >>>>> >>>>> >>>>> >>>>> >>>>> 数据库模擬

// FakeDB 資料是用來模擬一台假的数据库
type FakeDB struct {
	Loaded     bool
	MockResult map[uint32]mysql.Result
}

var FakeDBInstance FakeDB // 启动一个模拟的数据库实例

// Transferred 🧚 单元测试的测试资料载入定义接口
type Transferred interface {
	IsLoaded() bool  // 是否载入资料完成
	MarkLoaded()     // 标记载入资料完成
	LoadData() error // 进行测试资的载入资料
	// IsTakeOver() bool // 是否被单元测试接管
	// MarkTakeOver()    // 标记被单元测试接管
	// UnmarkTakeOver()  // 反标记被单元测试接管
}

// SubFakeDB 为模拟數據庫的部份资料
type SubFakeDB struct {
	addr     string       // 网路位置
	user     string       // 帐户
	password string       // 密码
	sql      string       // 执行字串
	result   mysql.Result // 數據庫回傳資料
}

// >>>>> >>>>> >>>>> >>>>> >>>>> 以下为 Key 值相关函式

// MakeMockKey 函式 🧚 为 在单元测试数据库时建立识明资料，主要是用来判别要回传的测试资料 (组成对应的 key)
// 从直连物件取出相关资料，组成新的 key 值，并回传
// 网路位置 192.168.122.2:3307
// 帐户 panhong
// 密码 12345
// 执行字串作为参数 SELECT * FROM `Library`.`Book`
func (dc *DirectConnection) MakeMockKey(sql string) uint32 {
	// 把相关的资料转成单纯的 key 值数字
	h := fnv.New32a()
	h.Write([]byte(dc.addr + ";" + dc.user + ";" + dc.password + ";" + sql)) // 所有的字串后面都要加上分号

	return h.Sum32() // 回传登记的数值
}

// MakeMockResult 函式 🧚 目前准备做法是 以下相关资料组合 对应到 一个数据库资料回传
// 从直连物件取出相关资料，组成新的 key 值
// 网路位置 192.168.122.2:3307
// 帐户 panhong
// 密码 12345
// 执行字串作为参数 SELECT * FROM `Library`.`Book`
func (fdb *FakeDB) MakeMockResult(data SubFakeDB) uint32 {
	// 把相关的资料转成单纯的 key 值数字
	h := fnv.New32a()
	h.Write([]byte(data.addr + ";" + data.user + ";" + data.password + ";" + data.sql)) // 所有的字串后面都要加上分号

	// 直接预先写好数据库资料回传
	FakeDBInstance.MockResult[h.Sum32()] = data.result // 转成数值，运算速度较快

	return h.Sum32() // 回传登记的数值
}
