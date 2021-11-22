package router

import (
	"github.com/XiaoMi/Gaea/models"
	"github.com/stretchr/testify/require"
	"testing"
)

// 参考文件 https://github.com/XiaoMi/Gaea/blob/master/docs/shard-example.md#gaea_kingshard_mod

// TestNovelRouterHashType 函式 🧚 是用来测试小說数据库的 hash 路由
func TestNovelRouterHashType(t *testing.T) {

	// >>>>> >>>>> >>>>> >>>>> >>>>> 案例1
	// 在第 1 台 Master 数据库有数据表 Book_0000
	// 在第 2 台 Master 数据库有数据表 Book_0001

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

	// 直接产生路由规则
	rule, err := parseRule(&cfgRouter)
	require.Equal(t, err, nil)

	// 检查目前的路由设定值
	require.Equal(t, rule.ruleType, "hash")
	require.Equal(t, rule.db, "novel")
	require.Equal(t, rule.table, "book")
	require.Equal(t, rule.slices, []string{"slice-0", "slice-1"})
	require.Equal(t, rule.shard.(*HashShard).ShardNum, 2)
	require.Equal(t, rule.shardingColumn, "bookid")

	// 下面的 rule.subTableIndexes 和 rule.tableToSlice 是传输函式 parseHashRuleSliceInfos 以 models.Shard 的 Locations 和 Slices 为参数，产生输出得来的
	require.Equal(t, rule.subTableIndexes, []int{0, 1})
	require.Equal(t, rule.tableToSlice, map[int]int{0: 0, 1: 1})

	require.Equal(t, len(rule.mycatDatabases), 0)
	require.Equal(t, len(rule.mycatDatabaseToTableIndexMap), 0)

	// >>>>> >>>>> >>>>> >>>>> >>>>> 案例2
	// 在第 1 台 Master 数据库有数据表 Book_0000
	// 在第 2 台 Master 数据库有数据表 Book_0001 Book_0002

	// 修改 路由规则 设定模组
	cfgRouter.Locations = []int{1, 2}

	// 直接产生路由规则
	rule, err = parseRule(&cfgRouter)
	require.Equal(t, err, nil)

	// 检查目前的路由设定值
	require.Equal(t, rule.subTableIndexes, []int{0, 1, 2})
	require.Equal(t, rule.tableToSlice, map[int]int{0: 0, 1: 1, 2: 1})
}

// TestNovelRouterModType 函式 🧚 是用来测试小說数据库的 mod 路由
func TestNovelRouterModType(t *testing.T) {

	// >>>>> >>>>> >>>>> >>>>> >>>>> 案例1
	// 在第 1 台 Master 数据库有数据表 Book_0000
	// 在第 2 台 Master 数据库有数据表 Book_0001

	// 再建立 路由规则 设定模组
	cfgRouter := models.Shard{
		DB:            "novel",
		Table:         "Book",
		ParentTable:   "",
		Type:          "mod",
		Key:           "BookID",
		Locations:     []int{1, 1},
		Slices:        []string{"slice-0", "slice-1"},
		TableRowLimit: 0,
	}

	// 直接产生路由规则
	rule, err := parseRule(&cfgRouter)
	require.Equal(t, err, nil)

	// 检查目前的路由设定值
	require.Equal(t, rule.ruleType, "mod")
	require.Equal(t, rule.db, "novel")
	require.Equal(t, rule.table, "book")
	require.Equal(t, rule.slices, []string{"slice-0", "slice-1"})
	require.Equal(t, rule.shard.(*ModShard).ShardNum, 2)
	require.Equal(t, rule.shardingColumn, "bookid")

	// 下面的 rule.subTableIndexes 和 rule.tableToSlice 是传输函式 parseHashRuleSliceInfos 以 models.Shard 的 Locations 和 Slices 为参数，产生输出得来的
	require.Equal(t, rule.subTableIndexes, []int{0, 1})
	require.Equal(t, rule.tableToSlice, map[int]int{0: 0, 1: 1})

	require.Equal(t, len(rule.mycatDatabases), 0)
	require.Equal(t, len(rule.mycatDatabaseToTableIndexMap), 0)

	// >>>>> >>>>> >>>>> >>>>> >>>>> 案例2
	// 在第 1 台 Master 数据库有数据表 Book_0000
	// 在第 2 台 Master 数据库有数据表 Book_0001 Book_0002

	// 修改 路由规则 设定模组
	cfgRouter.Locations = []int{1, 2}

	// 直接产生路由规则
	rule, err = parseRule(&cfgRouter)
	require.Equal(t, err, nil)

	// 检查目前的路由设定值
	require.Equal(t, rule.subTableIndexes, []int{0, 1, 2})
	require.Equal(t, rule.tableToSlice, map[int]int{0: 0, 1: 1, 2: 1})
}

// TestNovelRouterRangeType 函式 🧚 是用来测试小說数据库的 range 路由
func TestNovelRouterRangeType(t *testing.T) {

	// >>>>> >>>>> >>>>> >>>>> >>>>> 案例1
	// 在第 1 台 Master 数据库有数据表 Book_0000
	// 在第 2 台 Master 数据库有数据表 Book_0001

	// 再建立 路由规则 设定模组
	cfgRouter := models.Shard{
		DB:            "novel",
		Table:         "Book",
		ParentTable:   "",
		Type:          "range",
		Key:           "BookID",
		Locations:     []int{1, 1},
		Slices:        []string{"slice-0", "slice-1"},
		TableRowLimit: 3, // Book_ID 0 至 2  放在第一个切片，Book_ID 3 至 5  放在第二个切片，Book_ID 6 至 8  放在第三个切片，以前类推
	}

	// 直接产生路由规则
	rule, err := parseRule(&cfgRouter)
	require.Equal(t, err, nil)

	// 检查目前的路由设定值
	require.Equal(t, rule.ruleType, "range")
	require.Equal(t, rule.db, "novel")
	require.Equal(t, rule.table, "book")
	require.Equal(t, rule.slices, []string{"slice-0", "slice-1"})

	// 观察分表是否是依照主键 BookID
	require.Equal(t, rule.shardingColumn, "bookid")

	// 检查在切片中分表编号的上下界限的范围
	require.Equal(t, len(rule.shard.(*NumRangeShard).Shards), 2)
	require.Equal(t, rule.shard.(*NumRangeShard).Shards[0].Start, int64(0))
	require.Equal(t, rule.shard.(*NumRangeShard).Shards[0].End, int64(3))
	require.Equal(t, rule.shard.(*NumRangeShard).Shards[1].Start, int64(3))
	require.Equal(t, rule.shard.(*NumRangeShard).Shards[1].End, int64(6))

	// 下面的 rule.subTableIndexes 和 rule.tableToSlice 是传输函式 parseHashRuleSliceInfos 以 models.Shard 的 Locations 和 Slices 为参数，产生输出得来的
	require.Equal(t, rule.subTableIndexes, []int{0, 1})
	require.Equal(t, rule.tableToSlice, map[int]int{0: 0, 1: 1})

	require.Equal(t, len(rule.mycatDatabases), 0)
	require.Equal(t, len(rule.mycatDatabaseToTableIndexMap), 0)

	// >>>>> >>>>> >>>>> >>>>> >>>>> 案例2
	// 在第 1 台 Master 数据库有数据表 Book_0000
	// 在第 2 台 Master 数据库有数据表 Book_0001 Book_0002

	// 修改 路由规则 设定模组
	cfgRouter.Locations = []int{1, 2}

	// 直接产生路由规则
	rule, err = parseRule(&cfgRouter)
	require.Equal(t, err, nil)

	// 检查目前的路由设定值
	require.Equal(t, rule.subTableIndexes, []int{0, 1, 2})
	require.Equal(t, rule.tableToSlice, map[int]int{0: 0, 1: 1, 2: 1})

	// 检查在切片中分表编号的上下界限的范围
	require.Equal(t, len(rule.shard.(*NumRangeShard).Shards), 3)
	require.Equal(t, rule.shard.(*NumRangeShard).Shards[0].Start, int64(0))
	require.Equal(t, rule.shard.(*NumRangeShard).Shards[0].End, int64(3))
	require.Equal(t, rule.shard.(*NumRangeShard).Shards[1].Start, int64(3))
	require.Equal(t, rule.shard.(*NumRangeShard).Shards[1].End, int64(6))
	require.Equal(t, rule.shard.(*NumRangeShard).Shards[2].Start, int64(6))
	require.Equal(t, rule.shard.(*NumRangeShard).Shards[2].End, int64(9))
}

// TestNovelRouterModDateYear 函式 🧚 是用来测试小說数据库的 date year 路由
func TestNovelRouterModDateYear(t *testing.T) {

	// >>>>> >>>>> >>>>> >>>>> >>>>> 案例1
	// 在第 1 台 Master 数据库有数据表 Book_0000
	// 在第 2 台 Master 数据库有数据表 Book_0001

	// 再建立 路由规则 设定模组
	cfgRouter := models.Shard{
		DB:          "novel",
		Table:       "Book",
		ParentTable: "",
		Type:        "date_year",
		Key:         "Publish",
		// Locations:     []int{1, 1}, // 路由规则模式 date_year 不使用 Locations
		Slices:        []string{"slice-0", "slice-1"},
		DateRange:     []string{"1500-1600", "1601-1700"}, // 路由规则模式 DateRange 使用 date_year
		TableRowLimit: 0,
	}

	// 直接产生路由规则
	rule, err := parseRule(&cfgRouter)
	require.Equal(t, err, nil)

	// 检查目前的路由设定值
	require.Equal(t, rule.ruleType, "date_year")
	require.Equal(t, rule.db, "novel")
	require.Equal(t, rule.table, "book")
	require.Equal(t, rule.slices, []string{"slice-0", "slice-1"})
	// require.Equal(t, rule.shard.(*DateYearShard), rule.shard.(*DateYearShard))
	require.Equal(t, rule.shardingColumn, "publish")

	// 下面的 rule.subTableIndexes 和 rule.tableToSlice 是传输函式 parseHashRuleSliceInfos 以 models.Shard 的 Locations 和 Slices 为参数，产生输出得来的
	require.Equal(t, rule.subTableIndexes[0], 1500)
	require.Equal(t, rule.subTableIndexes[101], 1601) // 都同时加上 101
	require.Equal(t, rule.tableToSlice[1500], 0)      // 第一个范围的开头
	require.Equal(t, rule.tableToSlice[1600], 0)      // 第一個范围的结尾
	require.Equal(t, rule.tableToSlice[1601], 1)      // 加上 101 之后，进入下一个切片

	require.Equal(t, len(rule.mycatDatabases), 0)
	require.Equal(t, len(rule.mycatDatabaseToTableIndexMap), 0)
}

// TestNovelRouterModDateMonth 函式 🧚 是用来测试小說数据库的 date month 路由
func TestNovelRouterModDateMonth(t *testing.T) {

	// >>>>> >>>>> >>>>> >>>>> >>>>> 案例1
	// 在第 1 台 Master 数据库有数据表 Book_0000
	// 在第 2 台 Master 数据库有数据表 Book_0001

	// 再建立 路由规则 设定模组
	cfgRouter := models.Shard{
		DB:          "novel",
		Table:       "Book",
		ParentTable: "",
		Type:        "date_month",
		Key:         "Publish",
		// Locations:     []int{1, 1}, // 路由规则模式 date_month 不使用 Locations
		Slices:        []string{"slice-0", "slice-1"},
		DateRange:     []string{"150001-160012", "160101-170012"}, // 路由规则模式 DateRange 使用 date_month
		TableRowLimit: 0,
	}

	// 直接产生路由规则
	rule, err := parseRule(&cfgRouter)
	require.Equal(t, err, nil)

	// 检查目前的路由设定值
	require.Equal(t, rule.ruleType, "date_month")
	require.Equal(t, rule.db, "novel")
	require.Equal(t, rule.table, "book")
	require.Equal(t, rule.slices, []string{"slice-0", "slice-1"})
	// require.Equal(t, rule.shard.(*DateMonthShard), rule.shard.(*DateMonthShard))
	require.Equal(t, rule.shardingColumn, "publish")

	// 下面的 rule.subTableIndexes 和 rule.tableToSlice 是传输函式 parseHashRuleSliceInfos 以 models.Shard 的 Locations 和 Slices 为参数，产生输出得来的
	require.Equal(t, rule.subTableIndexes[0], 150001)
	require.Equal(t, rule.subTableIndexes[100*12], 160001) // 都同时加上 100 年有 12 个月
	require.Equal(t, rule.tableToSlice[150001], 0)         // 第一个范围的开头
	require.Equal(t, rule.tableToSlice[160100], 0)         // 第一個范围的结尾
	require.Equal(t, rule.tableToSlice[160101], 1)         // 加上 101 年之后，进入下一个切片

	require.Equal(t, len(rule.mycatDatabases), 0)
	require.Equal(t, len(rule.mycatDatabaseToTableIndexMap), 0)
}

// TestNovelRouterModDateDay 函式 🧚 是用来测试小說数据库的 date day 路由
func TestNovelRouterModDateDay(t *testing.T) {

	// >>>>> >>>>> >>>>> >>>>> >>>>> 案例1
	// 在第 1 台 Master 数据库有数据表 Book_0000
	// 在第 2 台 Master 数据库有数据表 Book_0001

	// 再建立 路由规则 设定模组
	cfgRouter := models.Shard{
		DB:          "novel",
		Table:       "Book",
		ParentTable: "",
		Type:        "date_day",
		Key:         "Publish",
		// Locations:     []int{1, 1}, // 路由规则模式 date_day 不使用 Locations
		Slices:        []string{"slice-0", "slice-1"},
		DateRange:     []string{"15000101-16001231", "16010101-17001231"}, // 路由规则模式 DateRange 使用 date_day
		TableRowLimit: 0,
	}

	// 直接产生路由规则
	rule, err := parseRule(&cfgRouter)
	require.Equal(t, err, nil)

	// 检查目前的路由设定值
	require.Equal(t, rule.ruleType, "date_day")
	require.Equal(t, rule.db, "novel")
	require.Equal(t, rule.table, "book")
	require.Equal(t, rule.slices, []string{"slice-0", "slice-1"})
	// require.Equal(t, rule.shard.(*DateMonthShard), rule.shard.(*DateMonthShard))
	require.Equal(t, rule.shardingColumn, "publish")

	// 下面的 rule.subTableIndexes 和 rule.tableToSlice 是传输函式 parseHashRuleSliceInfos 以 models.Shard 的 Locations 和 Slices 为参数，产生输出得来的
	require.Equal(t, rule.subTableIndexes[0], 15000101)
	require.Equal(t, rule.subTableIndexes[100*365+24], 16000101) // 加上 365 年 * 100 天 + 润月 24 天
	require.Equal(t, rule.tableToSlice[15000101], 0)             // 第一个范围的开头
	require.Equal(t, rule.tableToSlice[16001231], 0)             // 第一個范围的结尾
	require.Equal(t, rule.tableToSlice[16010101], 1)             // 加上 101 年之后，进入下一个切片

	require.Equal(t, len(rule.mycatDatabases), 0)
	require.Equal(t, len(rule.mycatDatabaseToTableIndexMap), 0)
}

// TestNovelRouterModMyCat 函式 🧚 是用来测试小說数据库的 mycat 路由
func TestNovelRouterModMyCat(t *testing.T) {

	// >>>>> >>>>> >>>>> >>>>> >>>>> 案例1
	// 在第 1 台 Master 数据库有数据表 Book_0000
	// 在第 2 台 Master 数据库有数据表 Book_0001

	// 再建立 路由规则 设定模组
	cfgRouter := models.Shard{
		DB:            "novel",
		Table:         "Book",
		ParentTable:   "",
		Type:          "mycat_mod",
		Key:           "bookid",
		Locations:     []int{1, 1},
		Slices:        []string{"slice-0", "slice-1"},
		Databases:     []string{"db_mycat_[0-1]"},
		TableRowLimit: 0,
	}

	// 直接产生路由规则
	rule, err := parseRule(&cfgRouter)
	require.Equal(t, err, nil)

	// 检查目前的路由设定值
	require.Equal(t, rule.ruleType, "mycat_mod")
	require.Equal(t, rule.db, "novel")
	require.Equal(t, rule.table, "book")
	require.Equal(t, rule.slices, []string{"slice-0", "slice-1"})
	require.Equal(t, rule.shard.(*MycatPartitionModShard).ShardNum, 2)
	require.Equal(t, rule.shardingColumn, "bookid")

	// 下面的 rule.subTableIndexes 和 rule.tableToSlice 是传输函式 parseHashRuleSliceInfos 以 models.Shard 的 Locations 和 Slices 为参数，产生输出得来的
	require.Equal(t, rule.subTableIndexes, []int{0, 1})
	require.Equal(t, rule.tableToSlice, map[int]int{0: 0, 1: 1})

	// 检查 MyCat 的路由设定值
	require.Equal(t, rule.mycatDatabases, []string{"db_mycat_0", "db_mycat_1"})
	require.Equal(t, rule.mycatDatabaseToTableIndexMap, map[string]int{"db_mycat_0": 0, "db_mycat_1": 1})

	// >>>>> >>>>> >>>>> >>>>> >>>>> 案例2
	// 在第 1 台 Master 数据库有数据表 Book_0000
	// 在第 2 台 Master 数据库有数据表 Book_0001 Book_0002

	// 修改 路由规则 设定模组
	cfgRouter.Locations = []int{1, 2}
	cfgRouter.Databases = []string{"db_mycat_[0-2]"}

	// 直接产生路由规则
	rule, err = parseRule(&cfgRouter)
	require.Equal(t, err, nil)

	// 检查目前的路由设定值
	require.Equal(t, rule.subTableIndexes, []int{0, 1, 2})
	require.Equal(t, rule.tableToSlice, map[int]int{0: 0, 1: 1, 2: 1})

	// 检查 MyCat 的路由设定值
	require.Equal(t, rule.mycatDatabases, []string{"db_mycat_0", "db_mycat_1", "db_mycat_2"})
	require.Equal(t, rule.mycatDatabaseToTableIndexMap, map[string]int{"db_mycat_0": 0, "db_mycat_1": 1, "db_mycat_2": 2})
}
