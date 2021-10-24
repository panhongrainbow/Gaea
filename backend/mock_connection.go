package backend

import (
	"fmt"
	"github.com/XiaoMi/Gaea/mysql"
	"hash/fnv"
	"sync"
)

const (
	transactionAutoCommit int = iota // 解析出的数值为 0
	transactionBegin                 // 解析出的数值为 1
	transactionCommit                // 解析出的数值为 2
	transactionRollback              // 解析出的数值为 3
)

// TakeOver >>>>> >>>>> >>>>> >>>>> >>>>> 单元测试的指示灯
var TakeOver bool // 现在是否由单元测试接管

// FakeDB >>>>> >>>>> >>>>> >>>>> >>>>> 数据库模擬

// FakeDB 資料是用來模擬一台假的数据库
type fakeDB struct {
	sync.Mutex
	Loaded       bool
	MockReAct    map[uint32]mysql.Result // 模拟数据库的回应，当有一串 SQL 传入时，先计算出特征值，再对应出应有的回应
	MockDataInDB []fakeSlice             // 模拟在数据库里的资料
}

// 進行交易 Transaction
// 当执行交易时，会暂存在 transaction 变数里，当 Commit 之后，才会移到 result
type fakeSlice struct {
	transaction int
	result      *mysql.Result // 已存入数据库的资料
	resultTmp   *mysql.Result // 进行交易 Transaction
}

var fakeDBInstance = make(map[string]*fakeDB) // 启动一个模拟的数据库实例

// Transferred 🧚 单元测试的测试资料载入定义接口
type Transferred interface {
	// IsInited 至 EmptyData 以下为 不会修改资料的 基本操作函式
	IsInited() bool         // 是否载入资料完成
	MarkInited()            // 标记载入资料完成
	UnMarkInited()          // 去除 载入资料完成 的标记
	EmptyData() error       // 清空已载入的测试资料
	GetFakeDB() string      // 取得登记后的数据库名称
	UseFakeDB(string) error // 登记将要使用的数据库名称
	// Lock 至 UnLock 上锁相关函式另外独立成函式
	// 因为频繁的上锁和解锁会影响效能，而且上锁和解锁的间隔可能会创造资料被改写的机会
	Lock()   // 上锁
	UnLock() // 解锁
	// InsertData 至 DeleteData 以下为 会修改资料的 基本操作函式
	InsertData() // 新增资料
	UpdateData() // 修改资料
	DeleteData() // 删除资料
}

// subFakeDB 为模拟數據庫的部份资料
type subFakeDB struct {
	addr     string       // 网路位置
	user     string       // 帐户
	password string       // 密码
	db       string       // 数据库名称
	sql      string       // 执行字串
	result   mysql.Result // 數據庫回傳資料
}

// >>>>> >>>>> >>>>> >>>>> >>>>> 以下为 Key 值相关函式

// MakeMockKey 函式 🧚 为 在单元测试数据库时建立识明资料，主要是用来判别要回传的测试资料 (组成对应的 key)
// 从直连物件取出相关资料，组成新的 key 值，并回传
// 网路位置 192.168.122.2:3307
// 帐户 panhong
// 密码 12345
// 执行字串作为参数 SELECT * FROM `novel`.`Book`
func (dc *DirectConnection) MakeMockKey(sql string) uint32 {
	// 把相关的资料转成单纯的 key 值数字
	h := fnv.New32a()
	h.Write([]byte(dc.addr + ";" + dc.user + ";" + dc.password + ";" + dc.db + ";" + sql)) // 所有的字串后面都要加上分号

	return h.Sum32() // 回传登记的数值
}

// MakeMockResult 函式 🧚 目前准备做法是 以下相关资料组合 对应到 一个数据库资料回传
// 从直连物件取出相关资料，组成新的 key 值
// 网路位置 192.168.122.2:3307
// 帐户 panhong
// 密码 12345
// 执行字串作为参数 SELECT * FROM `novel`.`Book`
func (fdb *fakeDB) MakeMockResult(data subFakeDB) uint32 {
	// 把相关的资料转成单纯的 key 值数字
	h := fnv.New32a()
	h.Write([]byte(data.addr + ";" + data.user + ";" + data.password + ";" + data.db + ";" + data.sql)) // 所有的字串后面都要加上分号

	// 直接预先写好数据库资料回传
	fakeDBInstance[data.db].MockReAct[h.Sum32()] = data.result // 转成数值，运算速度较快

	return h.Sum32() // 回传登记的数值
}

// Begin 函式 🧚 是用来开始进行事务操作
//    会把状态由 AutoCommit 变成 Begin
func (fdb *fakeDB) Begin() error {
	// >>>>> >>>>> >>>>> >>>>> >>>>> 先检查所有切片数据库的状态是不是 "自动 Commit"

	for i := 0; i < len(fdb.MockDataInDB); i++ {
		if fdb.MockDataInDB[i].transaction != transactionAutoCommit {
			return fmt.Errorf("cannot begin")
		}
	}

	// >>>>> >>>>> >>>>> >>>>> >>>>> 完成检查才能进行下一步

	// 修改所有的切片数据库的状态为 "交易开始"
	for i := 0; i < len(fdb.MockDataInDB); i++ {
		fdb.MockDataInDB[i].transaction = transactionBegin
	}

	// 回传无错误
	return nil
}

// Commit 函式 🧚 是用来完成进行事务操作
//    会把状态由 Begin 变成 Commit，再由 Commit 變成 AutoCommit
//    那就直接由 Begin 变成 AutoCommit
func (fdb *fakeDB) Commit() error {
	// >>>>> >>>>> >>>>> >>>>> >>>>> 先检查所有切片数据库的状态是不是 "事务开始"
	for i := 0; i < len(fdb.MockDataInDB); i++ {
		if fdb.MockDataInDB[i].transaction != transactionBegin {
			return fmt.Errorf("cannot commit")
		}
	}

	// >>>>> >>>>> >>>>> >>>>> >>>>> 完成检查后，才能搬移资料

	for i := 0; i < len(fdb.MockDataInDB); i++ { // 先复制
		for i := 0; i < len(fdb.MockDataInDB); i++ {
			for j := 0; j < len(fdb.MockDataInDB[i].resultTmp.RowDatas); j++ {
				fdb.MockDataInDB[i].result.RowDatas = append(fdb.MockDataInDB[i].result.RowDatas, fdb.MockDataInDB[i].resultTmp.RowDatas[j])
			}
			for k := 0; k < len(fdb.MockDataInDB[i].resultTmp.Values); k++ {
				fdb.MockDataInDB[i].result.Values = append(fdb.MockDataInDB[i].result.Values, fdb.MockDataInDB[i].resultTmp.Values[k])
			}
		}
	}
	for i := 0; i < len(fdb.MockDataInDB); i++ { // 再删除
		fdb.MockDataInDB[i].resultTmp = new(mysql.Result)
	}

	// >>>>> >>>>> >>>>> >>>>> >>>>> 完成资料搬移后，才能改变状态

	// 修改所有的切片数据库的状态为 "交易结束"
	for i := 0; i < len(fdb.MockDataInDB); i++ {
		fdb.MockDataInDB[i].transaction = transactionAutoCommit
	}

	// 回传无错误
	return nil
}

// Rollback 函式 🧚 是用来彻回进行事务操作
//    会把状态由 Begin 变成 Rollback，再由 Rollback 變成 AutoCommit
//    那就直接由 Begin 变成 AutoCommit
func (fdb *fakeDB) Rollback() error {
	// >>>>> >>>>> >>>>> >>>>> >>>>> 先检查所有切片数据库的状态是不是 "事务开始"
	for i := 0; i < len(fdb.MockDataInDB); i++ {
		if fdb.MockDataInDB[i].transaction != transactionBegin {
			return fmt.Errorf("cannot commit")
		}
	}

	// >>>>> >>>>> >>>>> >>>>> >>>>> 完成检查后，删除要彻回的资料
	for i := 0; i < len(fdb.MockDataInDB); i++ { // 再删除
		fdb.MockDataInDB[i].resultTmp = new(mysql.Result)
	}

	// >>>>> >>>>> >>>>> >>>>> >>>>> 完成资料搬移后，才能改变状态

	// 修改所有的切片数据库的状态为 "交易结束"
	for i := 0; i < len(fdb.MockDataInDB); i++ {
		fdb.MockDataInDB[i].transaction = transactionAutoCommit
	}

	// 回传无错误
	return nil
}
