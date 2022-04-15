package containerdTest

import (
	"context"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// TestDefaultContainerd 使用约定的接口进行测试 test default containerd interface
func TestMariaDBContainerd(t *testing.T) {
	// 测试约定的接口 test the default interface
	t.Run("test mariadb interface", func(t *testing.T) {
		// 建立新的容器的连接客户端 create a new client connected to the default socket path for containerd
		client, err := containerd.New("/run/containerd/containerd.sock")
		require.Nil(t, err)
		defer func() {
			_ = client.Close()
		}()

		// 测立一个新的命名空间 create a new context with a "mariadb" namespace
		ctx := namespaces.WithNamespace(context.Background(), "mariadb")

		// 建立测试对象 create a test object
		m := mariaDB{}

		// 拉取预设的测试印象档 pull the default test image from DockerHub
		img, err := m.Pull(client, ctx, "docker.io/library/mysql:8.0.28-debian")
		assert.Nil(t, err)

		// 建立一个新的预设容器 create a default container
		c, err := m.Create(client, ctx, "mariadb-server", "/var/run/netns/gaea-mariadb-sakila", img, "mariadb-server-snapshot")
		assert.Nil(t, err)

		// 建立新的容器工作 create a task from the container
		tsk, err := m.Task(c, ctx)
		assert.Nil(t, err)

		// start the task. 開始执行容器工作
		err = m.Start(tsk, ctx)
		assert.Nil(t, err)

		time.Sleep(time.Second * 60 * 10)

		// interrupt the task. 强制终止容器工作
		err = m.Interrupt(tsk, ctx)
		assert.Nil(t, err)

		// 删除容器和获得离开讯息 kill the process and get the exit status
		err = m.Delete(tsk, c, ctx)
		require.Nil(t, err)
	})
}
