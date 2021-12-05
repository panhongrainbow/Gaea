package router

import (
	"github.com/XiaoMi/Gaea/models"
	"github.com/stretchr/testify/require"
	"testing"
)

// 参考文件一 https://github.com/XiaoMi/Gaea/blob/master/docs/shard-example.md#gaea_kingshard_mod
// 参考文件二 https://github.com/XiaoMi/Gaea/blob/master/docs/shard.md

var (
	testSql = "INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(1, 9781517191276, 'Romance Of The Three Kingdoms', 'Luo Guanzhong', 1522, 'Historical fiction');"
)

// TestNovelRouterHashType 函式 🧚 是用来测试小說数据库的 hash 路由
func TestNovelRouterHashType(t *testing.T) {
	// >>>>> >>>>> >>>>> >>>>> >>>>> 设定档 1 cfgShard1
	// 在第 1 台 Master 数据库有数据表 Book_0000
	// 在第 2 台 Master 数据库有数据表 Book_0001

	// 建立 路由规则 设定模组 cfgShard1
	cfgShard1 := models.Shard{
		DB:            "novel",                        // 数据库
		Table:         "Book",                         // 数据表
		ParentTable:   "",                             // 暂不设定
		Type:          "hash",                         // hash 路由规则
		Key:           "BookID",                       // 以 BookID 栏位作为分表的依据
		Locations:     []int{1, 1},                    // 切片 slice-0 的数据表有 1 张，而 slice-1 的数据表有 1 张
		Slices:        []string{"slice-0", "slice-1"}, // 切片 Slice-0 和 Slice-1
		TableRowLimit: 0,                              // 暂不设定
	}

	// >>>>> >>>>> >>>>> >>>>> >>>>> 设定档 2 cfgShard2
	// 在第 1 台 Master 数据库有数据表 Book_0000
	// 在第 2 台 Master 数据库有数据表 Book_0001 Book_0002

	// 建立 路由规则 设定模组 cfgShard2
	cfgShard2 := models.Shard{
		DB:            "novel",                        // 数据库
		Table:         "Book",                         // 数据表
		ParentTable:   "",                             // 暂不设定
		Type:          "hash",                         // hash 路由规则
		Key:           "BookID",                       // 以 BookID 栏位作为分表的依据
		Locations:     []int{1, 2},                    // 只修改这里，代表切片 slice-0 的数据表有 1 张，而 slice-1 的数据表有 2 张
		Slices:        []string{"slice-0", "slice-1"}, // 切片 Slice-0 和 Slice-1
		TableRowLimit: 0,                              // 暂不设定
	}

	// 建立测试资料
	tests := []struct {
		cfgShard        models.Shard // 路由设定档
		shardNum        int          // 切片的数量
		subTableIndexes []int        // 在路由规则里数据表的 Index
		tableToSlice    map[int]int  // 在路由规则里切片的 Index
		insertBookID    []int        // 插入数据库的 BookID 的值
		tableIndex      []int        // 数据表的 Index
		sliceIndex      []int        // 数据表的 Index 和 切片的 Index 的对应
	}{
		{
			cfgShard:        cfgShard1,               // 路由规则变数 cfgShard1
			shardNum:        2,                       // 数据表的数量为 2
			subTableIndexes: []int{0, 1},             // 0，1 分别对应到数据表的 Book_0000，Book_0001
			tableToSlice:    map[int]int{0: 0, 1: 1}, // 数据表 Book_0000，Book_0001 分别对应到 Slice-0，Slice-1
			insertBookID:    []int{1, 2, 3},          // 在数据库分别插入 BookID 为 1，2 和 3 的资料
			tableIndex:      []int{1, 0, 1},          // BookID 为 1，2 和 3 的资料分别会插入 Book_0001，Book_0000 和 Book_0001
			sliceIndex:      []int{1, 0, 1},          // BookID 为 1，2 和 3 的资料分别会插入 slice-1，slice-0 和 slice-1
		},
		{
			cfgShard:        cfgShard2,                     // 路由规则变数 cfgShard2
			shardNum:        3,                             // 数据表的数量为 3
			subTableIndexes: []int{0, 1, 2},                // 0，1 和 2 分别对应到数据表的 Book_0000，Book_0001 和 Book_0002
			tableToSlice:    map[int]int{0: 0, 1: 1, 2: 1}, // 数据表 Book_0000，Book_0001 和 Book_0002 分别对应到 Slice-0，Slice-1 和 Slice-1
			insertBookID:    []int{1, 2, 3},                // 在数据库分别插入 BookID 为 1，2 和 3 的资料
			tableIndex:      []int{1, 2, 0},                // BookID 为 1，2 和 3 的资料分别会插入 Book_0001，Book_0002 和 Book_0000
			sliceIndex:      []int{1, 1, 0},                // BookID 为 1，2 和 3 的资料分别会插入 slice-1，slice-1 和 slice-0
		},
	}

	// 开始进行测试
	for i := 0; i < len(tests); i++ {
		// 直接产生路由规则
		rule, err := parseRule(&tests[i].cfgShard)
		require.Equal(t, err, nil)

		// 检查目前的路由设定值
		require.Equal(t, rule.ruleType, "hash")
		require.Equal(t, rule.db, "novel")
		require.Equal(t, rule.table, "book")
		require.Equal(t, rule.slices, []string{"slice-0", "slice-1"})
		require.Equal(t, rule.shard.(*HashShard).ShardNum, tests[i].shardNum)
		require.Equal(t, rule.shardingColumn, "bookid")

		// 下面的 rule.subTableIndexes 和 rule.tableToSlice 是传输函式 parseHashRuleSliceInfos 以 models.Shard 的 Locations 和 Slices 为参数，产生输出得来的
		require.Equal(t, rule.subTableIndexes, tests[i].subTableIndexes)
		require.Equal(t, rule.tableToSlice, tests[i].tableToSlice)

		require.Equal(t, len(rule.mycatDatabases), 0)
		require.Equal(t, len(rule.mycatDatabaseToTableIndexMap), 0)

		// 直接建立路由
		rt := new(Router)
		rt.rules = make(map[string]map[string]Rule)
		m := make(map[string]Rule)
		rt.rules[rule.db] = m
		rt.rules[rule.db][rule.table] = rule

		// 直接建立预设路由
		rt.defaultRule = NewDefaultRule(rule.slices[0]) // 设定第一组切片为预设路由

		// 会回传布林值显示路由规则是否存在，在路由中用一开始设定的资料库和资料表，就可以找到路由规则
		_, has := rt.GetShardRule(rule.db, rule.table)
		require.Equal(t, has, true)

		// 检查插入的 BookID 和路由规则进行组合
		for j := 0; j < len(tests[i].insertBookID); j++ {
			// 由路由推算出要插入到那一个切片的表
			tableIndex, err := rt.rules[rule.db][rule.table].FindTableIndex(tests[i].insertBookID[j])
			require.Equal(t, err, nil)
			require.Equal(t, tableIndex, tests[i].tableIndex[j]) // 检查插入的表编号
			sliceIndex := rt.rules[rule.db][rule.table].GetSliceIndexFromTableIndex(tableIndex)
			require.Equal(t, sliceIndex, tests[i].sliceIndex[j]) // 检查插入的切片编号
		}
	}
}

// TestNovelRouterModType 函式 🧚 是用来测试小說数据库的 mod 路由
func TestNovelRouterModType(t *testing.T) {
	// >>>>> >>>>> >>>>> >>>>> >>>>> 设定档 1 cfgShard1
	// 在第 1 台 Master 数据库有数据表 Book_0000
	// 在第 2 台 Master 数据库有数据表 Book_0001

	// 建立 路由规则 设定模组 cfgShard1
	cfgShard1 := models.Shard{
		DB:            "novel",                        // 数据库
		Table:         "Book",                         // 数据表
		ParentTable:   "",                             // 暂不设定
		Type:          "mod",                          // mod 路由规则
		Key:           "BookID",                       // 以 BookID 栏位作为分表的依据
		Locations:     []int{1, 1},                    // 切片 slice-0 的数据表有 1 张，而 slice-1 的数据表有 1 张
		Slices:        []string{"slice-0", "slice-1"}, // 切片 Slice-0 和 Slice-1
		TableRowLimit: 0,                              // 暂不设定
	}

	// >>>>> >>>>> >>>>> >>>>> >>>>> 设定档 2 cfgShard2
	// 在第 1 台 Master 数据库有数据表 Book_0000
	// 在第 2 台 Master 数据库有数据表 Book_0001 Book_0002

	// 建立 路由规则 设定模组 cfgShard2
	cfgShard2 := models.Shard{
		DB:            "novel",                        // 数据库
		Table:         "Book",                         // 数据表
		ParentTable:   "",                             // 暂不设定
		Type:          "mod",                          // mod 路由规则
		Key:           "BookID",                       // 以 BookID 栏位作为分表的依据
		Locations:     []int{1, 2},                    // 只修改这里，代表切片 slice-0 的数据表有 1 张，而 slice-1 的数据表有 2 张
		Slices:        []string{"slice-0", "slice-1"}, // 切片 Slice-0 和 Slice-1
		TableRowLimit: 0,                              // 暂不设定
	}

	// 建立测试资料
	tests := []struct {
		cfgShard        models.Shard // 路由设定档
		shardNum        int          // 切片的数量
		subTableIndexes []int        // 在路由规则里数据表的 Index
		tableToSlice    map[int]int  // 在路由规则里切片的 Index
		insertBookID    []int        // 插入数据库的 BookID 的值
		tableIndex      []int        // 数据表的 Index
		sliceIndex      []int        // 数据表的 Index 和 切片的 Index 的对应
	}{
		// 插入的 BookID 为正值
		{
			cfgShard:        cfgShard1,               // 路由规则变数 cfgShard1
			shardNum:        2,                       // 数据表的数量为 2
			subTableIndexes: []int{0, 1},             // 0，1 分别对应到数据表的 Book_0000，Book_0001
			tableToSlice:    map[int]int{0: 0, 1: 1}, // 数据表 Book_0000，Book_0001 分别对应到 Slice-0，Slice-1
			insertBookID:    []int{1, 2, 3},          // 在数据库分别插入 BookID 为 1，2 和 3 的资料
			tableIndex:      []int{1, 0, 1},          // BookID 为 1，2 和 3 的资料分别会插入 Book_0001，Book_0000 和 Book_0001
			sliceIndex:      []int{1, 0, 1},          // BookID 为 1，2 和 3 的资料分别会插入 slice-1，slice-0 和 slice-1
		},
		{ // 插入的 BookID 为正值
			cfgShard:        cfgShard2,                     // 路由规则变数 cfgShard2
			shardNum:        3,                             // 数据表的数量为 3
			subTableIndexes: []int{0, 1, 2},                // 0，1 和 2 分别对应到数据表的 Book_0000，Book_0001 和 Book_0002
			tableToSlice:    map[int]int{0: 0, 1: 1, 2: 1}, // 数据表 Book_0000，Book_0001 和 Book_0002 分别对应到 Slice-0，Slice-1 和 Slice-1
			insertBookID:    []int{1, 2, 3},                // 在数据库分别插入 BookID 为 1，2 和 3 的资料
			tableIndex:      []int{1, 2, 0},                // BookID 为 1，2 和 3 的资料分别会插入 Book_0001，Book_0002 和 Book_0000
			sliceIndex:      []int{1, 1, 0},                // BookID 为 1，2 和 3 的资料分别会插入 slice-1，slice-1 和 slice-0
		},
		// 插入的 BookID 为负值 (路由规则 mod 可以处理负值的插入键值)
		{
			cfgShard:        cfgShard1,               // 路由规则变数 cfgShard1
			shardNum:        2,                       // 数据表的数量为 2
			subTableIndexes: []int{0, 1},             // 0，1 分别对应到数据表的 Book_0000，Book_0001
			tableToSlice:    map[int]int{0: 0, 1: 1}, // 数据表 Book_0000，Book_0001 分别对应到 Slice-0，Slice-1
			insertBookID:    []int{-1, -2, -3},       // 在数据库分别插入 BookID 为 -1，-2 和 -3 的资料
			tableIndex:      []int{1, 0, 1},          // BookID 为 -1，-2 和 -3 的资料分别会插入 Book_0001，Book_0000 和 Book_0001
			sliceIndex:      []int{1, 0, 1},          // BookID 为 -1，-2 和 -3 的资料分别会插入 slice-1，slice-0 和 slice-1
		},
		{ // 插入的 BookID 为负值 (路由规则 mod 可以处理负值的插入键值)
			cfgShard:        cfgShard2,                     // 路由规则变数 cfgShard2
			shardNum:        3,                             // 数据表的数量为 3
			subTableIndexes: []int{0, 1, 2},                // 0，1 和 2 分别对应到数据表的 Book_0000，Book_0001 和 Book_0002
			tableToSlice:    map[int]int{0: 0, 1: 1, 2: 1}, // 数据表 Book_0000，Book_0001 和 Book_0002 分别对应到 Slice-0，Slice-1 和 Slice-1
			insertBookID:    []int{-1, -2, -3},             // 在数据库分别插入 BookID 为 -1，-2 和 -3 的资料
			tableIndex:      []int{1, 2, 0},                // BookID 为 -1，-2 和 -3 的资料分别会插入 Book_0001，Book_0002 和 Book_0000
			sliceIndex:      []int{1, 1, 0},                // BookID 为 -1，-2 和 -3 的资料分别会插入 slice-1，slice-1 和 slice-0
		},
	}

	// 开始进行测试
	for i := 0; i < len(tests); i++ {
		// 直接产生路由规则
		rule, err := parseRule(&tests[i].cfgShard)
		require.Equal(t, err, nil)

		// 检查目前的路由设定值
		require.Equal(t, rule.ruleType, "mod")
		require.Equal(t, rule.db, "novel")
		require.Equal(t, rule.table, "book")
		require.Equal(t, rule.slices, []string{"slice-0", "slice-1"})
		require.Equal(t, rule.shard.(*ModShard).ShardNum, tests[i].shardNum)
		require.Equal(t, rule.shardingColumn, "bookid")

		// 下面的 rule.subTableIndexes 和 rule.tableToSlice 是传输函式 parseHashRuleSliceInfos 以 models.Shard 的 Locations 和 Slices 为参数，产生输出得来的
		require.Equal(t, rule.subTableIndexes, tests[i].subTableIndexes)
		require.Equal(t, rule.tableToSlice, tests[i].tableToSlice)

		require.Equal(t, len(rule.mycatDatabases), 0)
		require.Equal(t, len(rule.mycatDatabaseToTableIndexMap), 0)

		// 直接建立路由
		rt := new(Router)
		rt.rules = make(map[string]map[string]Rule)
		m := make(map[string]Rule)
		rt.rules[rule.db] = m
		rt.rules[rule.db][rule.table] = rule

		// 直接建立预设路由
		rt.defaultRule = NewDefaultRule(rule.slices[0]) // 设定第一组切片为预设路由

		// 会回传布林值显示路由规则是否存在，在路由中用一开始设定的资料库和资料表，就可以找到路由规则
		_, has := rt.GetShardRule(rule.db, rule.table)
		require.Equal(t, has, true)

		// 检查插入的 BookID 和路由规则进行组合
		for j := 0; j < len(tests[i].insertBookID); j++ {
			// 由路由推算出要插入到那一个切片的表
			tableIndex, err := rt.rules[rule.db][rule.table].FindTableIndex(tests[i].insertBookID[j])
			require.Equal(t, err, nil)
			require.Equal(t, tableIndex, tests[i].tableIndex[j]) // 检查插入的表编号
			sliceIndex := rt.rules[rule.db][rule.table].GetSliceIndexFromTableIndex(tableIndex)
			require.Equal(t, sliceIndex, tests[i].sliceIndex[j]) // 检查插入的切片编号
		}
	}
}

// TestNovelRouterRangeType 函式 🧚 是用来测试小說数据库的 range 路由
func TestNovelRouterRangeType(t *testing.T) {
	// >>>>> >>>>> >>>>> >>>>> >>>>> 设定档 1 cfgShard1
	// 在第 1 台 Master 数据库有数据表 Book_0000
	// 在第 2 台 Master 数据库有数据表 Book_0001

	// 建立 路由规则 设定模组 cfgShard1
	cfgShard1 := models.Shard{
		DB:            "novel",                        // 数据库
		Table:         "Book",                         // 数据表
		ParentTable:   "",                             // 暂不设定
		Type:          "range",                        // range 路由规则
		Key:           "BookID",                       // 以 BookID 栏位作为分表的依据
		Locations:     []int{1, 1},                    // 切片 slice-0 的数据表有 1 张，而 slice-1 的数据表有 1 张
		Slices:        []string{"slice-0", "slice-1"}, // 切片 Slice-0 和 Slice-1
		TableRowLimit: 3,                              // 暂不设定
	}

	// >>>>> >>>>> >>>>> >>>>> >>>>> 设定档 2 cfgShard2
	// 在第 1 台 Master 数据库有数据表 Book_0000
	// 在第 2 台 Master 数据库有数据表 Book_0001 Book_0002

	// 建立 路由规则 设定模组 cfgShard2
	cfgShard2 := models.Shard{
		DB:            "novel",                        // 数据库
		Table:         "Book",                         // 数据表
		ParentTable:   "",                             // 暂不设定
		Type:          "range",                        // range 路由规则
		Key:           "BookID",                       // 以 BookID 栏位作为分表的依据
		Locations:     []int{1, 2},                    // 只修改这里，代表切片 slice-0 的数据表有 1 张，而 slice-1 的数据表有 2 张
		Slices:        []string{"slice-0", "slice-1"}, // 切片 Slice-0 和 Slice-1
		TableRowLimit: 3,                              // 暂不设定
	}

	// 建立测试资料
	tests := []struct {
		cfgShard        models.Shard   // 路由设定档
		shardNum        int            // 切片的数量
		subTableIndexes []int          // 在路由规则里数据表的 Index
		tableToSlice    map[int]int    // 在路由规则里切片的 Index
		shardsStartEnd  map[int][2]int // 数据表资料的上下界限范围
		insertBookID    []int          // 插入数据库的 BookID 的值
		tableIndex      []int          // 数据表的 Index
		sliceIndex      []int          // 数据表的 Index 和 切片的 Index 的对应
	}{
		{
			cfgShard:        cfgShard1,                            // 路由规则变数 cfgShard1
			shardNum:        2,                                    // 数据表的数量为 2
			subTableIndexes: []int{0, 1},                          // 0，1 分别对应到数据表的 Book_0000，Book_0001
			tableToSlice:    map[int]int{0: 0, 1: 1},              // 数据表 Book_0000，Book_0001 分别对应到 Slice-0，Slice-1
			shardsStartEnd:  map[int][2]int{0: {0, 3}, 1: {3, 6}}, // 数据表 book_0000 所被写入 bookID 的最小值为 0，最大值为 3
			insertBookID:    []int{0, 1, 2, 3, 4, 5, 6},           // 在数据库分别插入 BookID 为 0，1 和 2 等等 的资料
			tableIndex:      []int{0, 0, 0, 1, 1, 1, -1},          // BookID 为 0，1，2 和 3 等等 的资料分别会插入 Book_0000，Book_0000，Book_0000 和 Book_0001
			sliceIndex:      []int{0, 0, 0, 1, 1, 1, -1},          // BookID 为 0，1，2 和 3 等等 的资料分别会插入 slice-0，slice-0，slice-0 和 slice-1
		},
		{
			cfgShard:        cfgShard2,                                       // 路由规则变数 cfgShard2
			shardNum:        3,                                               // 数据表的数量为 3
			subTableIndexes: []int{0, 1, 2},                                  // 0，1 和 2 分别对应到数据表的 Book_0000，Book_0001 和 Book_0002
			tableToSlice:    map[int]int{0: 0, 1: 1, 2: 1},                   // 数据表 Book_0000，Book_0001 和 Book_0002 分别对应到 Slice-0，Slice-1 和 Slice-1
			shardsStartEnd:  map[int][2]int{0: {0, 3}, 1: {3, 6}, 2: {6, 9}}, // 数据表 book_0000 所被写入 bookID 的最小值为 0，最大值为 3
			insertBookID:    []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},             // 在数据库分别插入 BookID 为 1，2 和 3 等等 的资料
			tableIndex:      []int{0, 0, 0, 1, 1, 1, 2, 2, 2, -1},            // BookID 为 0，1，2 和 3 等等 的资料分别会插入 Book_0000，Book_0000，Book_0000 和 Book_0001
			sliceIndex:      []int{0, 0, 0, 1, 1, 1, 1, 1, 1, -1},            // BookID 为 0，1，2 和 3 等等 的资料分别会插入 slice-0，slice-0，slice-0 和 slice-1
		},
	}

	// 开始进行测试
	for i := 0; i < len(tests); i++ {
		// 直接产生路由规则
		rule, err := parseRule(&tests[i].cfgShard)
		require.Equal(t, err, nil)

		// 检查目前的路由设定值
		require.Equal(t, rule.ruleType, "range")
		require.Equal(t, rule.db, "novel")
		require.Equal(t, rule.table, "book")
		require.Equal(t, rule.slices, []string{"slice-0", "slice-1"})
		require.Equal(t, len(rule.shard.(*NumRangeShard).Shards), tests[i].shardNum)
		require.Equal(t, rule.shardingColumn, "bookid")

		// 检查在切片中分表编号的上下界限的范围，数据表被写入 bookID 的最小值、最值值 (上下界)
		// shardIndex 为数据表的编号，比如 0 为数据表 Book_0000，1 为数据表 Book_0001 等等
		// shardRange 为数据表被写入 bookID 的最小值、最大值 (上下界) 的范围，比如 数据表 Book_0000 的最小值为 0、最值值为 3
		for shardIndex, shardRange := range tests[i].shardsStartEnd {
			require.Equal(t, rule.shard.(*NumRangeShard).Shards[shardIndex].Start, int64(shardRange[0])) // 最小值 (上界)
			require.Equal(t, rule.shard.(*NumRangeShard).Shards[shardIndex].End, int64(shardRange[1]))   // 最大值 (下界)
		}

		// 下面的 rule.subTableIndexes 和 rule.tableToSlice 是传输函式 parseHashRuleSliceInfos 以 models.Shard 的 Locations 和 Slices 为参数，产生输出得来的
		require.Equal(t, rule.subTableIndexes, tests[i].subTableIndexes)
		require.Equal(t, rule.tableToSlice, tests[i].tableToSlice)

		require.Equal(t, len(rule.mycatDatabases), 0)
		require.Equal(t, len(rule.mycatDatabaseToTableIndexMap), 0)

		// 直接建立路由
		rt := new(Router)
		rt.rules = make(map[string]map[string]Rule)
		m := make(map[string]Rule)
		rt.rules[rule.db] = m
		rt.rules[rule.db][rule.table] = rule

		// 直接建立预设路由
		rt.defaultRule = NewDefaultRule(rule.slices[0]) // 设定第一组切片为预设路由

		// 会回传布林值显示路由规则是否存在，在路由中用一开始设定的资料库和资料表，就可以找到路由规则
		_, has := rt.GetShardRule(rule.db, rule.table)
		require.Equal(t, has, true)

		// 检查插入的 BookID 和路由规则进行组合
		for j := 0; j < len(tests[i].insertBookID); j++ {
			// 由路由推算出要插入到那一个切片的表
			tableIndex, _ := rt.rules[rule.db][rule.table].FindTableIndex(tests[i].insertBookID[j])
			// require.Equal(t, err, nil)
			// 当插入的 BookID 超过数据表的界限时，就会发生错误，比如 Book_0000 的范围为 0 至 3，Book_0001 的范围为 3 至 6，没有 Book_0002 这张数据表
			// 当插入 BookID 为 6 时，就会发生错误
			// 但在 BookID 为 0 至 5 时，就不会发生错误
			// 所以这里不能进行测试
			require.Equal(t, tableIndex, tests[i].tableIndex[j]) // 检查插入的表编号
			sliceIndex := rt.rules[rule.db][rule.table].GetSliceIndexFromTableIndex(tableIndex)
			require.Equal(t, sliceIndex, tests[i].sliceIndex[j]) // 检查插入的切片编号
		}
	}
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

	// 直接建立路由
	rt := new(Router)
	rt.rules = make(map[string]map[string]Rule)
	m := make(map[string]Rule)
	rt.rules[rule.db] = m
	rt.rules[rule.db][rule.table] = rule

	// 由路由推算出要插入到那一个切片的表
	insertTableIndex, err := rt.rules[rule.db][rule.table].FindTableIndex("1500") // 数值 1500 是值 SQL 字串中的 publish 为 1500，这是经由 parser 传入的值
	require.Equal(t, err, nil)
	require.Equal(t, insertTableIndex, 1500) // 数据表的 index 為 1500

	// 由多个数据表组成切片，所以可以由 数据表的 index 转成 切片的 index
	insertSliceIndex := rt.rules[rule.db][rule.table].GetSliceIndexFromTableIndex(insertTableIndex)
	require.Equal(t, insertSliceIndex, 0) // 数据表的 index 為 1500 所对应的切片 index 为 0，插入的数据表为 Book_0000

	insertTableIndex, err = rt.rules[rule.db][rule.table].FindTableIndex("1601") // 数值 1601 是值 SQL 字串中的 publish 为 1601，这是经由 parser 传入的值
	require.Equal(t, err, nil)
	require.Equal(t, insertTableIndex, 1601) // 数据表的 index 為 1500

	// 由多个数据表组成切片，所以可以由 数据表的 index 转成 切片的 index
	insertSliceIndex = rt.rules[rule.db][rule.table].GetSliceIndexFromTableIndex(insertTableIndex)
	require.Equal(t, insertSliceIndex, 1) // 数据表的 index 為 1601 所对应的切片 index 为 1，插入的数据表为 Book_0001
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

// TestNovelRouterModMyCat 函式 🧚 是用来测试小說数据库的 MyCat 路由
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

// TestNovelRouterModMyCatLong 函式 🧚 是用来测试小說数据库的 MyCat Long 路由 (固定hash分片算法)
func TestNovelRouterModMyCatLong(t *testing.T) {

	// >>>>> >>>>> >>>>> >>>>> >>>>> 案例1
	// 在第 1 台 Master 数据库有数据表 Book_0000
	// 在第 2 台 Master 数据库有数据表 Book_0001

	// 再建立 路由规则 设定模组
	cfgRouter := models.Shard{
		DB:              "novel",
		Table:           "Book",
		ParentTable:     "",
		Type:            "mycat_long",
		Key:             "bookid",
		Locations:       []int{1, 1},
		Slices:          []string{"slice-0", "slice-1"},
		Databases:       []string{"db_mycat_[0-1]"},
		TableRowLimit:   0,
		PartitionCount:  "2",   // 此值为 Locations 阵列里的 1+1
		PartitionLength: "512", // 此值为 1024 / 2 = 512
	}

	// 直接产生路由规则
	rule, err := parseRule(&cfgRouter)
	require.Equal(t, err, nil)

	// 检查目前的路由设定值
	require.Equal(t, rule.ruleType, "mycat_long")
	require.Equal(t, rule.db, "novel")
	require.Equal(t, rule.table, "book")
	require.Equal(t, rule.slices, []string{"slice-0", "slice-1"})
	require.Equal(t, rule.shard.(*MycatPartitionLongShard).shardNum, 2)
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
	cfgRouter.Locations = []int{2, 2}
	cfgRouter.Databases = []string{"db_mycat_[0-3]"}
	cfgRouter.PartitionCount = "4"    // 此值为 Locations 阵列里的 2+2
	cfgRouter.PartitionLength = "256" // 此值为 1024 / 4 = 256
	// cfgRouter.Locations = []int{1, 2} 这种设定不存在，因为 1024 不能被 3 整除

	// 直接产生路由规则
	rule, err = parseRule(&cfgRouter)
	require.Equal(t, err, nil)

	// 检查目前的路由设定值
	require.Equal(t, rule.subTableIndexes, []int{0, 1, 2, 3})
	require.Equal(t, rule.tableToSlice, map[int]int{0: 0, 1: 0, 2: 1, 3: 1})

	// 检查 MyCat 的路由设定值
	require.Equal(t, rule.mycatDatabases, []string{"db_mycat_0", "db_mycat_1", "db_mycat_2", "db_mycat_3"})
	require.Equal(t, rule.mycatDatabaseToTableIndexMap, map[string]int{"db_mycat_0": 0, "db_mycat_1": 1, "db_mycat_2": 2, "db_mycat_3": 3})
}

// TestNovelRouterModMyCatMurmur 函式 🧚 是用来测试小說数据库的 MyCat Murmur 路由
// murmur 算法是将字段进行hash后分发到不同的数据库,字段类型支持int和varchar
func TestNovelRouterModMyCatMurmur(t *testing.T) {

	// >>>>> >>>>> >>>>> >>>>> >>>>> 案例1
	// 在第 1 台 Master 数据库有数据表 Book_0000
	// 在第 2 台 Master 数据库有数据表 Book_0001

	// 再建立 路由规则 设定模组
	cfgRouter := models.Shard{
		DB:                 "novel",
		Table:              "Book",
		ParentTable:        "",
		Type:               "mycat_murmur",
		Key:                "bookid",
		Locations:          []int{1, 1},
		Slices:             []string{"slice-0", "slice-1"},
		Databases:          []string{"db_mycat_[0-1]"},
		TableRowLimit:      0,
		Seed:               "0",
		VirtualBucketTimes: "160",
	}

	// 直接产生路由规则
	rule, err := parseRule(&cfgRouter)
	require.Equal(t, err, nil)

	// 检查目前的路由设定值
	require.Equal(t, rule.ruleType, "mycat_murmur")
	require.Equal(t, rule.db, "novel")
	require.Equal(t, rule.table, "book")
	require.Equal(t, rule.slices, []string{"slice-0", "slice-1"})

	// require.Equal(t, rule.shard.(*MycatPartitionMurmurHashShard), rule.shard.(*MycatPartitionMurmurHashShard))
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

// TestNovelRouterModMyCatHashString 函式 🧚 是用来测试小說数据库的 MyCat Hash String 路由
// hash string 算法为取部份字串进行 hash
func TestNovelRouterModMyCatHashString(t *testing.T) {

	// >>>>> >>>>> >>>>> >>>>> >>>>> 案例1
	// 在第 1 台 Master 数据库有数据表 Book_0000
	// 在第 2 台 Master 数据库有数据表 Book_0001

	// 再建立 路由规则 设定模组
	cfgRouter := models.Shard{
		DB:              "novel",
		Table:           "Book",
		ParentTable:     "",
		Type:            "mycat_string",
		Key:             "isbn",
		Locations:       []int{1, 1},
		Slices:          []string{"slice-0", "slice-1"},
		Databases:       []string{"db_mycat_[0-1]"},
		PartitionCount:  "2",    // 此值为 Locations 阵列里的 1+1
		PartitionLength: "512",  // 此值为 1024 / 2 = 512
		HashSlice:       "-2:0", // 取 Isbn 栏位的最后两个字作 Hash 计算
	}

	// 直接产生路由规则
	rule, err := parseRule(&cfgRouter)
	require.Equal(t, err, nil)

	// 检查目前的路由设定值
	require.Equal(t, rule.ruleType, "mycat_string")
	require.Equal(t, rule.db, "novel")
	require.Equal(t, rule.table, "book")
	require.Equal(t, rule.slices, []string{"slice-0", "slice-1"})

	require.Equal(t, rule.shard.(*MycatPartitionStringShard).shardNum, 2)
	require.Equal(t, rule.shardingColumn, "isbn")

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
	cfgRouter.Locations = []int{2, 2}
	cfgRouter.Databases = []string{"db_mycat_[0-3]"}
	cfgRouter.PartitionCount = "4"    // 此值为 Locations 阵列里的 2+2
	cfgRouter.PartitionLength = "256" // 此值为 1024 / 4 = 256
	// cfgRouter.Locations = []int{1, 2} 这种设定不存在，因为 1024 不能被 3 整除

	// 直接产生路由规则
	rule, err = parseRule(&cfgRouter)
	require.Equal(t, err, nil)

	// 检查目前的路由设定值
	require.Equal(t, rule.subTableIndexes, []int{0, 1, 2, 3})
	require.Equal(t, rule.tableToSlice, map[int]int{0: 0, 1: 0, 2: 1, 3: 1})

	// 检查 MyCat 的路由设定值
	require.Equal(t, rule.mycatDatabases, []string{"db_mycat_0", "db_mycat_1", "db_mycat_2", "db_mycat_3"})
	require.Equal(t, rule.mycatDatabaseToTableIndexMap, map[string]int{"db_mycat_0": 0, "db_mycat_1": 1, "db_mycat_2": 2, "db_mycat_3": 3})
}
