package parser

import (
	. "github.com/XiaoMi/Gaea/parser/format"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

// 这个测试主要是用来测试 Parser 的每个细节，并有图文对照
func TestA(t *testing.T) {
	var sb strings.Builder
	parser := New()

	tables := []testCase{
		// SQL 指令 SELECT
		{"SELECT /*+ MAX_EXECUTION_TIME(1000) */ * FROM db.t", true, "SELECT /*+ MAX_EXECUTION_TIME(1000)*/ * FROM `db`.`t`"}, // 使用优化器
	}

	for _ ,table := range tables {
		stmts, _, err := parser.Parse(table.src, "", "")
		require.Equal(t, err, nil)
		restoreSQLs := ""
		for _, stmt := range stmts {
			sb.Reset()
			err = stmt.Restore(NewRestoreCtx(DefaultRestoreFlags, &sb))
			restoreSQL := sb.String()
			restoreStmt, err := parser.ParseOneStmt(restoreSQL, "", "")
			require.Equal(t, err, nil)
			CleanNodeText(stmt)
			CleanNodeText(restoreStmt)
			if restoreSQLs != "" {
				restoreSQLs += "; "
			}
			restoreSQLs += restoreSQL
		}

		require.Equal(t, restoreSQLs, table.restore)
	}
}
