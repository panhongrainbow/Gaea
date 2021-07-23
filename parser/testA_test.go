package parser

import (
	"strings"
	"testing"

	. "github.com/XiaoMi/Gaea/parser/format"
	"github.com/stretchr/testify/require"
)

// TestA 这个测试主要是用来测试 Parser 的每个细节，并有图文对照
func TestA(t *testing.T) {
	var sb strings.Builder
	parser := New()

	tests := []testCase{
		// 针对 Cluster db0 db0-0 db0-1 或其他的丛集 常用的 SQL 指令进行测试
		{"SELECT * FROM Library.Book", true, "SELECT * FROM `Library`.`Book`"},

		// SQL 指令 SELECT
		{"SELECT /*+ MAX_EXECUTION_TIME(1000) */ * FROM db.t", true, "SELECT /*+ MAX_EXECUTION_TIME(1000)*/ * FROM `db`.`t`"},            // (1) 使用优化器
		{"SELECT DISTINCT v FROM db.t", true, "SELECT DISTINCT `v` FROM `db`.`t`"},                                                       // (2) 过滤重复出现的纪录值
		{"SELECT DISTINCTROW v FROM db.t", true, "SELECT DISTINCT `v` FROM `db`.`t`"},                                                    // (3) 过滤重复出现的纪录值。在网路上查资料，说明 DISTINCTROW 是针对一整行数据。但是 misc.go 视 DISTINCT 和 DISTINCTROW 相同，都为 DISTINCT
		{"SELECT ALL v FROM db.t", true, "SELECT `v` FROM `db`.`t`"},                                                                     // (4) 如果是 ALL 的话，DistinctOpt 在 parser.y 就会直接被设为 false，ALL 就会被忽略。
		{"SELECT HIGH_PRIORITY v FROM db.t", true, "SELECT HIGH_PRIORITY `v` FROM `db`.`t`"},                                             // (5) 高优先权
		{"SELECT LOW_PRIORITY v FROM db.t", true, "SELECT LOW_PRIORITY `v` FROM `db`.`t`"},                                               // (6) 低优先权
		{"SELECT DELAYED v FROM db.t", true, "SELECT DELAYED `v` FROM `db`.`t`"},                                                         // (7) 延迟写入，先回报客户 SQL 完成。大多是用在 INSERT ，看看用在 SQL 会不会有问题
		{"SELECT SQL_CACHE v FROM db.t", true, "SELECT `v` FROM `db`.`t`"},                                                               // (8) 资料库查询时使用快取
		{"SELECT SQL_NO_CACHE v FROM db.t", true, "SELECT SQL_NO_CACHE `v` FROM `db`.`t`"},                                               // (9) 资料库查询时不使用快取，直接对资料库进行性能测试
		{"SELECT SQL_CALC_FOUND_ROWS v FROM db.t", true, "SELECT `v` FROM `db`.`t`"},                                                     // (10) func (n *SelectStmt) Restore(ctx *format.RestoreCtx) error 函式不处理 SQL_CALC_FOUND_ROWS 字串
		{"SELECT * FROM t1 INNER JOIN t2 WHERE t1.id=t2.id", true, "SELECT * FROM `t1` JOIN `t2` WHERE `t1`.`id`=`t2`.`id`"},             // (11) 执行 INNER JOIN
		{"SELECT * FROM t1 STRAIGHT_JOIN t2 WHERE t1.id=t2.id", true, "SELECT * FROM `t1` STRAIGHT_JOIN `t2` WHERE `t1`.`id`=`t2`.`id`"}, // (12) 执行 STRAIGHT JOIN

		// 新增用的 SQL 指令 INSERT
		{"INSERT IGNORE INTO t VALUES (1), (2), (3)", true, "INSERT IGNORE INTO `t` VALUES (1),(2),(3)"}, // (13) 插入时有重复记录就会忽略
		{"INSERT t VALUES (1), (2), (3)", true, "INSERT INTO `t` VALUES (1),(2),(3)"},                    // (14) 这次没有使用 INTO 这个 TOKEN
	}

	for _, test := range tests {
		stmts, _, err := parser.Parse(test.src, "", "")
		require.Equal(t, err, nil)
		restoreSQLs := ""
		for _, stmt := range stmts {
			sb.Reset()
			err = stmt.Restore(NewRestoreCtx(DefaultRestoreFlags, &sb))
			require.Equal(t, err, nil)
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

		require.Equal(t, restoreSQLs, test.restore)
	}
}
