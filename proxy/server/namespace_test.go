package server

import (
	"fmt"
	"github.com/XiaoMi/Gaea/models"
	"testing"
)

// TestNameSpace 函式 🧚 是用来测试 NameSpace 的创造和运行
func TestNameSpace(t *testing.T) {
	// 直连 DC 的单元测试是否能正常启动
	TestNameSpaceCreate(t)
}

// TestNameSpaceCreate 函式 🧚 是用来测试 NameSpace 的创造
func TestNameSpaceCreate(t *testing.T) {
	// 先建立 models 设定档
	cfg := models.Slice{
		Name: "slice-0",
	}
	fmt.Println(cfg)
}
