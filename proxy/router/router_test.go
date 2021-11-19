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
