package containerdTest

import (
	"context"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
	"github.com/stretchr/testify/require"
	"testing"
)

// TestRunMariadb 测试數據庫容器服务运行
func TestRunMariadb(t *testing.T) {
	client, err := containerd.New("/run/containerd/containerd.sock")
	require.Nil(t, err)
	defer func() {
		_ = client.Close()
	}()

	// db := new(mariaDB)
	ctx := namespaces.WithNamespace(context.Background(), "example")
	// err = db.Pull(client, ctx, "docker.io/library/mariadb:latest")
	_, err = client.Pull(ctx, "docker.io/library/mariadb:latest", containerd.WithPullUnpack)
	require.Nil(t, err)
}
