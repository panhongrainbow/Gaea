package containerdTest

import (
	"context"
	"fmt"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/oci"
	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/stretchr/testify/require"
	"syscall"
	"testing"
	"time"
)

// TestRunMariadb 测试數據庫容器服务运行 is to test database container service run.
func TestRunMariadb(t *testing.T) {
	// create a new client connected to the default socket path for containerd
	client, err := containerd.New("/run/containerd/containerd.sock")
	require.Nil(t, err)
	defer client.Close()

	// create a new context with an "example" namespace
	ctx := namespaces.WithNamespace(context.Background(), "example")

	// pull the redis image from DockerHub
	image, err := client.Pull(ctx, "docker.io/library/redis:alpine3.13", containerd.WithPullUnpack)
	require.Nil(t, err)

	ns := specs.LinuxNamespace{specs.NetworkNamespace, "/var/run/netns/gaea-mariadb-sakila"}

	// create a container
	container, err := client.NewContainer(
		ctx,
		"redis-server",
		containerd.WithImage(image),
		containerd.WithNewSnapshot("redis-server-snapshot", image),
		containerd.WithNewSpec(oci.WithImageConfig(image), oci.WithLinuxNamespace(ns)),
	)

	require.Nil(t, err)
	defer container.Delete(ctx, containerd.WithSnapshotCleanup)

	// create a task from the container
	task, err := container.NewTask(ctx, cio.NewCreator(cio.WithStdio))
	require.Nil(t, err)
	defer task.Delete(ctx)

	// make sure we wait before calling start
	exitStatusC, err := task.Wait(ctx)
	require.Nil(t, err)

	// call start on the task to execute the redis server
	err = task.Start(ctx)
	require.Nil(t, err)

	test, _ := task.Spec(ctx)
	fmt.Println(test)

	// sleep for a lil bit to see the logs
	time.Sleep(120 * time.Second)

	// kill the process and get the exit status
	err = task.Kill(ctx, syscall.SIGTERM)
	require.Nil(t, err)

	// wait for the process to fully exit and print out the exit status

	status := <-exitStatusC
	code, _, err := status.Result()
	require.Nil(t, err)
	fmt.Printf("redis-server exited with status: %d\n", code)

	require.Nil(t, err)
}
