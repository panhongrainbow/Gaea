package parser

import (
	"fmt"
	"testing"

	"github.com/XiaoMi/Gaea/parser/ast"
	"github.com/XiaoMi/Gaea/parser/tidb-types/parser_driver"
)

func TestNodeToString(t *testing.T) {
	tableName := "tb1"
	d := &driver.ValueExpr{}
	d.SetValue(tableName)
	s, err := NodeToStringWithoutQuote(d)
	if err != nil {
		t.Fatal(err)
	}
	if s != tableName {
		t.Errorf("table name not equal, expect: %s, actual: %s", tableName, s)
	}
}

type NodePrintVisitor struct {
	// 为了要让参观者可以记录节点进入和离开的次数，所以新增以下内容
	enterTimes int // 进入次数
	leaveTimes int // 离开次数
	enterBreak int // 中断的 进入次数
	leaveBreak int // 中断的 离开次数
	// 以下进行测试
}

func (v *NodePrintVisitor) Enter(n ast.Node) (ast.Node, bool) {
	v.enterTimes++
	if v.enterTimes == v.enterBreak {
		fmt.Println("在这里下中断点")
	}
	fmt.Printf("第 %d 次进入: %T\n", v.enterTimes, n)
	return n, false
}

func (v *NodePrintVisitor) Leave(n ast.Node) (ast.Node, bool) {
	v.leaveTimes++
	if v.leaveTimes == v.leaveBreak && v.leaveTimes != 0 {
		fmt.Println("在这里下中断点")
	}
	fmt.Printf("第 %d 次离开: %T\n", v.leaveTimes, n)
	return n, true
}

func TestASTNode(t *testing.T) {
	sql := `desc xm_order`
	n, err := ParseSQL(sql)
	if err != nil {
		t.Fatalf("parse sql error: %v", err)
	}
	v := &NodePrintVisitor{}
	n.Accept(v)
}
