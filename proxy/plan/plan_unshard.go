// Copyright 2019 The Gaea Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package plan

import (
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/XiaoMi/Gaea/backend"
	"github.com/XiaoMi/Gaea/mysql"
	"github.com/XiaoMi/Gaea/parser/ast"
	"github.com/XiaoMi/Gaea/parser/format"
	"github.com/XiaoMi/Gaea/util"
)

// Transferred 🧚 单元测试的定义接口
type Transferred interface {
	IsTakeOver() bool // 是否被单元测试接管
	MarkTakeOver()    // 标记被单元测试接管
	UnmarkTakeOver()  // 反标记被单元测试接管
}

// MockPlanClient 🧚 单元测试数据库计划客户端
type MockPlanClient struct {
	// 单元测试相关设定值
	TakeOver bool // 现在是否由单元测试接管
	// 单元测试资料回应
	Result map[uint32]mysql.Result // 模拟数据库资料回传
}

// MarkTakeOver 函式 🧚 为 MockPlanClient 资料执行单元测试数据库直连的标记函式 (设定)
func (m *MockPlanClient) MarkTakeOver() {
	m.TakeOver = true // 单元测试之后可以直接进行接管
}

// IsTakeOver 函式 🧚 为 MockPlanClient 资料执行单元测试数据库计划的确认函式 (设定)
func (m *MockPlanClient) IsTakeOver() bool {
	// 因为不是每个函式或过程会完整初始化 Mock Client 变数，如果没有这一层保护，防止 nil 指标的错误
	if m == nil {
		return false // 回传 false ，之后单元测试不允许进行介入程式内部的运作
	}
	return m.TakeOver // 只要是回传 true ，之后单元测试就会接管整个程式
}

// UnmarkTakeOver 函式 🧚 为 MockPlanClient 资料执行单元测试数据库计划的反标记函式 (设定)
func (m *MockPlanClient) UnmarkTakeOver() {
	m.TakeOver = false // 解除单元测试的接管状态
}

// MakeResult 🧚 为 在单元测试数据库时建立计划回应资料的对应 (回应)
// 目前准备做法是 1设定 环境 2数据库名称 3SQL 指令 三个值的组合对应到 一个数据库资料回传
func (m *MockPlanClient) MakeResult(db, sql string, res mysql.Result) uint32 {
	// 把数据库和SQL字串转成单纯的数字
	h := fnv.New32a()
	h.Write([]byte(db + ";" + sql + ";")) // 所有的字串后面都要加上分号

	// 直接预先写好数据库资料回传
	m.Result[h.Sum32()] = res // 转成数值，运算速度较快
	return h.Sum32()          // 回传登记的数值
}

// UnshardPlan is the plan for unshard statement
// 此资料被单元测试函式包围
type UnshardPlan struct {
	basePlan

	db     string
	phyDBs map[string]string
	sql    string
	stmt   ast.StmtNode

	// 🧚 扩增一些单元测试的属性
	MockPlan *MockPlanClient // 单元测试数据库计划
	Trans    Transferred     // 单元测试的定义接口
}

// MarkTakeOver 函式 🧚 为 UnshardPlan 资料执行单元测试数据库计划的标记函式 (设定)
func (plan *UnshardPlan) MarkTakeOver() {
	plan.MockPlan.MarkTakeOver() // 操作底层函式
}

// IsTakeOver 函式 🧚 为 UnshardPlan 资料执行单元测试数据库计划的确认函式 (设定)
func (plan *UnshardPlan) IsTakeOver() bool {
	return plan.MockPlan.IsTakeOver() // 操作底层函式
}

// UnmarkTakeOver 函式 🧚 为 UnshardPlan 资料执行单元测试数据库计划的反标记函式 (设定)
func (plan *UnshardPlan) UnmarkTakeOver() {
	plan.MockPlan.UnmarkTakeOver() // 操作底层函式
}

// SelectLastInsertIDPlan is the plan for SELECT LAST_INSERT_ID()
// TODO: fix below
// https://dev.mysql.com/doc/refman/5.6/en/information-functions.html#function_last-insert-id
// The value of LAST_INSERT_ID() is not changed if you set the AUTO_INCREMENT column of a row
// to a non-“magic” value (that is, a value that is not NULL and not 0).
type SelectLastInsertIDPlan struct {
	basePlan
}

// IsSelectLastInsertIDStmt check if the statement is SELECT LAST_INSERT_ID()
func IsSelectLastInsertIDStmt(stmt ast.StmtNode) bool {
	s, ok := stmt.(*ast.SelectStmt)
	if !ok {
		return false
	}

	if len(s.Fields.Fields) != 1 {
		return false
	}

	if s.From != nil || s.Where != nil || s.GroupBy != nil || s.Having != nil || s.OrderBy != nil || s.Limit != nil {
		return false
	}

	f, ok := s.Fields.Fields[0].Expr.(*ast.FuncCallExpr)
	if !ok {
		return false
	}

	return f.FnName.L == "last_insert_id"
}

// CreateUnshardPlan constructor of UnshardPlan
func CreateUnshardPlan(stmt ast.StmtNode, phyDBs map[string]string, db string, tableNames []*ast.TableName) (*UnshardPlan, error) {
	p := &UnshardPlan{
		db:     db,
		phyDBs: phyDBs,
		stmt:   stmt,
	}
	rewriteUnshardTableName(phyDBs, tableNames)
	rsql, err := generateUnshardingSQL(stmt)
	if err != nil {
		return nil, fmt.Errorf("generate unshardPlan SQL error: %v", err)
	}
	p.sql = rsql
	return p, nil
}

func rewriteUnshardTableName(phyDBs map[string]string, tableNames []*ast.TableName) {
	for _, tableName := range tableNames {
		if phyDB, ok := phyDBs[tableName.Schema.String()]; ok {
			tableName.Schema.O = phyDB
			tableName.Schema.L = strings.ToLower(phyDB)
		}
	}
}

func generateUnshardingSQL(stmt ast.StmtNode) (string, error) {
	s := &strings.Builder{}
	ctx := format.NewRestoreCtx(format.EscapeRestoreFlags, s)
	_ = stmt.Restore(ctx)
	return s.String(), nil
}

// CreateSelectLastInsertIDPlan constructor of SelectLastInsertIDPlan
func CreateSelectLastInsertIDPlan() *SelectLastInsertIDPlan {
	return &SelectLastInsertIDPlan{}
}

// ExecuteIn implement Plan
func (p *UnshardPlan) ExecuteIn(reqCtx *util.RequestContext, se Executor) (*mysql.Result, error) {
	// 🧚 单元测试接管
	if p.MockPlan.IsTakeOver() {
		return mysql.SelectLibrayResult(), nil // 立刻中斷
	}

	r, err := se.ExecuteSQL(reqCtx, backend.DefaultSlice, p.db, p.sql)
	if err != nil {
		return nil, err
	}

	// set last insert id to session
	if _, ok := p.stmt.(*ast.InsertStmt); ok {
		if r.InsertID != 0 {
			se.SetLastInsertID(r.InsertID)
		}
	}

	return r, nil
}

// ExecuteIn implement Plan
func (p *SelectLastInsertIDPlan) ExecuteIn(reqCtx *util.RequestContext, se Executor) (*mysql.Result, error) {
	r := createLastInsertIDResult(se.GetLastInsertID())
	return r, nil
}

func createLastInsertIDResult(lastInsertID uint64) *mysql.Result {
	name := "last_insert_id()"
	var column = 1
	var rows [][]string
	var names = []string{
		name,
	}

	var t = fmt.Sprintf("%d", lastInsertID)
	rows = append(rows, []string{t})

	r := new(mysql.Resultset)

	var values = make([][]interface{}, len(rows))
	for i := range rows {
		values[i] = make([]interface{}, column)
		for j := range rows[i] {
			values[i][j] = rows[i][j]
		}
	}

	r, _ = mysql.BuildResultset(nil, names, values)
	ret := &mysql.Result{
		Resultset: r,
	}

	return ret
}
