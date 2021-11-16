package router

import (
	"github.com/XiaoMi/Gaea/models"
	"github.com/stretchr/testify/require"
	"testing"
)

// TestNovelRouter 函式 🧚 是用来测试小說資料庫的路由
func TestNovelRouter(t *testing.T) {

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

	// 直接建立路由规则
	rule, err := parseRule(&cfgRouter)
	require.Equal(t, err, nil)

	// 直接检查路由的设定值
	require.Equal(t, rule.ruleType, "hash")
	require.Equal(t, rule.db, "novel")
	require.Equal(t, rule.table, "book")
	require.Equal(t, rule.slices, []string{"slice-0", "slice-1"})
	require.Equal(t, rule.shard.(*HashShard).ShardNum, 2)
	require.Equal(t, rule.shardingColumn, "bookid")
	require.Equal(t, rule.subTableIndexes, []int{0, 1})
	require.Equal(t, rule.tableToSlice, map[int]int{0: 0, 1: 1})
	require.Equal(t, len(rule.mycatDatabases), 0)
	require.Equal(t, len(rule.mycatDatabaseToTableIndexMap), 0)
}
