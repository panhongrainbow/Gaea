package server

import (
	"github.com/XiaoMi/Gaea/backend"
	"github.com/XiaoMi/Gaea/models"
	"github.com/XiaoMi/Gaea/parser"
	"github.com/XiaoMi/Gaea/parser/ast"
	"github.com/XiaoMi/Gaea/proxy/router"
	"github.com/stretchr/testify/require"
	"testing"
)

// TestSlice 函式 🧚 是用来测试 NameSpace 的创造和运行
func TestNameSpace(t *testing.T) {
	// 测试整个 Slice 切片的创造、连线和读写数据库
	TestNovelSliceConnect(t)
	// 测试建立 NameSpace 的切片 Slice
	TestNovelNameSpaceSlice(t)
	// 测试建立 NameSpace 的路由 Router 规则
	TestNovelNameSpaceRouter(t)
}

// TestNovelSliceConnect 函式 🧚 是用来测试整个 Slice 切片的创造、连线和读写数据库
func TestNovelSliceConnect(t *testing.T) {
	// 初始化单元测试程式 (只要注解 Mark TakeOver() 就会使用真的资料库，不然就会跑单元测试)
	backend.MarkTakeOver() // MarkTakeOver 函式一定要放在单元测试最前面，因为可以提早启动一些 DEBUG 除错机制

	// >>>>> >>>>> >>>>> >>>>> >>>>> 建立新的切片变数
	// 先建立 models Slice 设定档
	cfg := models.Slice{
		Name:     "slice-0",
		UserName: "panhong",
		Password: "12345",
		Master:   "192.168.122.2:3309",
		Slaves: []string{
			"192.168.122.2:3310",
			"192.168.122.2:3311",
		},
		StatisticSlaves: nil,
		Capacity:        12,
		MaxCapacity:     24,
		IdleTimeout:     60,
	}

	// 用设定档开始建立新的切片
	newSlice, err := parseSlice(&cfg, "utf8mb4", 46)
	require.Equal(t, err, nil)

	// >>>>> >>>>> >>>>> >>>>> >>>>> 建立 Master 数据库的连线
	pcM, err := newSlice.GetConn(false, 0)
	require.Equal(t, err, nil)

	// >>>>> >>>>> >>>>> >>>>> >>>>> 向 Master 数据库写入资料
	// 使用 数据库
	err = pcM.UseDB("novel")
	require.Equal(t, err, nil)

	// 新增一笔资料
	result, err := pcM.Execute("INSERT INTO `novel`.`Book_0000` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (2,9789869442060,'Water Margin','Shi Nai an',1589,'Historical fiction')", 100)
	require.Equal(t, err, nil)

	// 检查数据库写入结果
	require.Equal(t, err, nil)
	require.Equal(t, result.AffectedRows, uint64(0x1))
	require.Equal(t, result.InsertID, uint64(0x0))

	// >>>>> >>>>> >>>>> >>>>> >>>>> 建立 Slave 数据库的连线
	pcS, err := newSlice.GetConn(true, 0)
	require.Equal(t, err, nil)

	// >>>>> >>>>> >>>>> >>>>> >>>>> 向 Slave 数据库读取资料
	// 查询一笔资料
	result, err = pcS.Execute("SELECT * FROM `novel`.`Book_0000`", 100)
	require.Equal(t, err, nil)

	// 检查数据库读取结果
	require.Equal(t, result.Resultset.Values[0][0].(int64), int64(2))
	require.Equal(t, result.Resultset.Values[0][1].(int64), int64(9789869442060))
	require.Equal(t, result.Resultset.Values[0][2].(string), "Water Margin")

	// 关闭关闭整个单元测试
	backend.UnmarkTakeOver()
}

// TestNovelNameSpaceSlice 函式 🧚 是用来测试建立 NameSpace 的切片 Slice
func TestNovelNameSpaceSlice(t *testing.T) {
	// 初始化单元测试程式 (只要注解 Mark TakeOver() 就会使用真的资料库，不然就会跑单元测试)
	backend.MarkTakeOver() // MarkTakeOver 函式一定要放在单元测试最前面，因为可以提早启动一些 DEBUG 除错机制

	// 先建立一个空的 NameSpace
	ns := new(Namespace)

	// 再建立一群切片的设定值

	// 先建立对于切片 Slice-0 的 models Slice 设定档
	cfgSlice0 := models.Slice{
		Name:     "slice-0",
		UserName: "panhong",
		Password: "12345",
		Master:   "192.168.122.2:3309",
		Slaves: []string{
			"192.168.122.2:3310",
			"192.168.122.2:3311",
		},
		StatisticSlaves: nil,
		Capacity:        12,
		MaxCapacity:     24,
		IdleTimeout:     60,
	}

	// 先建立对于切片 Slice-1 的 models Slice 设定档
	cfgSlice1 := models.Slice{
		Name:     "slice-1",
		UserName: "panhong",
		Password: "12345",
		Master:   "192.168.122.2:3312",
		Slaves: []string{
			"192.168.122.2:3313",
			"192.168.122.2:3314",
		},
		StatisticSlaves: nil,
		Capacity:        12,
		MaxCapacity:     24,
		IdleTimeout:     60,
	}

	// 组成 models slice 阵列
	cfgSliceS := make([]*models.Slice, 0, 2)
	cfgSliceS = append(cfgSliceS, &cfgSlice0)
	cfgSliceS = append(cfgSliceS, &cfgSlice1)

	// 建立 NameSpace 的切片阵列
	tmp, err := parseSlices(cfgSliceS, "utf8mb4", 46)
	require.Equal(t, err, nil)
	ns.slices = tmp

	// 建立 Master 资料库连线
	pcM, err := ns.GetSlice("slice-0").GetConn(false, 0)
	require.Equal(t, err, nil)
	require.Equal(t, pcM.GetAddr(), "192.168.122.2:3309")

	// 利用 Master 数据库连线去读写资料
	result, err := pcM.Execute("INSERT INTO `novel`.`Book_0000` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (2,9789869442060,'Water Margin','Shi Nai an',1589,'Historical fiction')", 100)
	require.Equal(t, err, nil)

	// 检查数据库写入结果
	require.Equal(t, err, nil)
	require.Equal(t, result.AffectedRows, uint64(0x1))
	require.Equal(t, result.InsertID, uint64(0x0))

	// 立刻读取数据库写入的结果
	result, err = pcM.Execute("SELECT * FROM `novel`.`Book_0000`", 100)
	require.Equal(t, err, nil)

	// 检查数据库读取结果
	require.Equal(t, result.Resultset.Values[0][0].(int64), int64(2))
	require.Equal(t, result.Resultset.Values[0][1].(int64), int64(9789869442060))
	require.Equal(t, result.Resultset.Values[0][2].(string), "Water Margin")

	// 关闭单元测试
	backend.UnmarkTakeOver()
}

// TestNovelNameSpaceRouter 函式 🧚 是用来测试  建立 NameSpace 的路由 Router 规则
// 要组成 NameSpace 路由，需要在模组里先组成 (1)切片模组 (2)预设路由 ，和 (3)路由模组
func TestNovelNameSpaceRouter(t *testing.T) {
	// 初始化单元测试程式 (只要注解 Mark TakeOver() 就会使用真的资料库，不然就会跑单元测试)
	backend.MarkTakeOver() // MarkTakeOver 函式一定要放在单元测试最前面，因为可以提早启动一些 DEBUG 除错机制

	// 建立 NameSpace 物件
	ns := new(Namespace)

	// >>>>> >>>>> >>>>> >>>>> >>>>> 开始组合所有的设定模组

	// 先建立一个空的 NameSpace 设定模组
	cfgNs := new(models.Namespace)

	// >>>>> >>>>> >>>>> 先组成 (1)切片模组

	// 先建立对于切片 Slice-0 的 models Slice 设定档
	cfgSlice0 := models.Slice{
		Name:     "slice-0",
		UserName: "panhong",
		Password: "12345",
		Master:   "192.168.122.2:3309",
		Slaves: []string{
			"192.168.122.2:3310",
			"192.168.122.2:3311",
		},
		StatisticSlaves: nil,
		Capacity:        12,
		MaxCapacity:     24,
		IdleTimeout:     60,
	}

	// 先建立对于切片 Slice-1 的 models Slice 设定档
	cfgSlice1 := models.Slice{
		Name:     "slice-1",
		UserName: "panhong",
		Password: "12345",
		Master:   "192.168.122.2:3312",
		Slaves: []string{
			"192.168.122.2:3313",
			"192.168.122.2:3314",
		},
		StatisticSlaves: nil,
		Capacity:        12,
		MaxCapacity:     24,
		IdleTimeout:     60,
	}

	// 组成 models slice 阵列
	cfgSliceS := make([]*models.Slice, 0, 2)
	cfgSliceS = append(cfgSliceS, &cfgSlice0)
	cfgSliceS = append(cfgSliceS, &cfgSlice1)

	// >>>>> >>>>> >>>>> 再组成 (2)预设路由

	cfgNs.DefaultSlice = "slice-0"

	// >>>>> >>>>> >>>>> 再组成 (3)路由模组

	// 再建立 路由规则 设定模组
	cfgRouter := models.Shard{
		DB:            "novel",
		Table:         "Book",
		ParentTable:   "",
		Type:          "hash",
		Key:           "BookID",
		Locations:     []int{1, 1},
		Slices:        []string{"slice-0", "slice-1"},
		TableRowLimit: 0,
	}

	// NameSpace 设定模组 载入 路由规则 设定模组
	cfgNs.ShardRules = make([]*models.Shard, 0)
	cfgNs.ShardRules = append(cfgNs.ShardRules, &cfgRouter)

	// >>>>> >>>>> >>>>> 把所有的模组组合完成后，建立 NameSpace 物件的 路由 和 切片阵列

	// 建立 NameSpace 物件的 切片阵列
	tmp, err := parseSlices(cfgSliceS, "utf8mb4", 46)
	require.Equal(t, err, nil)
	ns.slices = tmp
	cfgNs.Slices = cfgSliceS

	// 建立 NameSpace 物件的 路由
	tmp1, err := router.NewRouter(cfgNs)
	require.Equal(t, err, nil)
	ns.router = tmp1

	// 开始进行路由操作
	rule := ns.router.GetRule("novel", "book")
	require.Equal(t, rule.GetDB(), "novel")
	require.Equal(t, rule.GetTable(), "book")

	// >>>>> >>>>> >>>>> 进行 Select 的 Parser 操作
	// checker := plan.NewChecker("novel", ns.router)
	newParser0 := parser.New()
	newStmts0, _, err := newParser0.Parse("SELECT MIN(Publish) FROM novel.Book;", "", "")
	require.Equal(t, err, nil)

	// >>>>> >>>>> >>>>> 检查 Select 操作的旗标
	expr := newStmts0[0].(*ast.SelectStmt).Fields.Fields[0].Expr
	require.Equal(t, expr.GetFlag(), uint64(0x18))
	// FlagHasReference 值为 8
	// FlagHasAggregateFunc 值为 16
	// 两者值相加为 8 + 16 = 24 (十进位) 等同于 18 (十六进位)

	// >>>>> >>>>> >>>>> 进行 Insert 的 Parser 操作
	newParser1 := parser.New()
	_, _, err = newParser1.Parse("INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(1, 9781517191276, 'Romance Of The Three Kingdoms', 'Luo Guanzhong', 1522, 'Historical fiction');", "", "")
	require.Equal(t, err, nil)

	// >>>>> >>>>> >>>>> 检查 Insert 操作的旗标
	// (略过) 因为 Insert 好像没有旗标

	// 关闭单元测试
	backend.UnmarkTakeOver()
}
