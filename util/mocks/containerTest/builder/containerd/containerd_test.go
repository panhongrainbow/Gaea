package containerd

import (
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	Version = "" // 数据库版本 MariaDB Version
)

func TestMain(m *testing.M) {
	// 想辨法精简这里
	// mgr.Setup()
	/*Version = "10"
	fmt.Println(Version)*/

	// 退出
	// exitCode := m.Run()
	// os.Exit(exitCode)
}

func TestSoarTest(t *testing.T) {
	require.Equal(t, true, true)
}
