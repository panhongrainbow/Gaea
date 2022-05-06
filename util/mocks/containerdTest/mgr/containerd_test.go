package mgr

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

var (
	Version = "" // 数据库版本 MariaDB Version
)

func TestMain(m *testing.M) {
	// 想辨法精简这里
	setup()
	/*Version = "10"
	fmt.Println(Version)*/

	// 退出
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestSoarTest(t *testing.T) {
	require.Equal(t, true, true)
}
