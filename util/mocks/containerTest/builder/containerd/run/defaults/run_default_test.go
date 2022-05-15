package defaults

import (
	"context"
	"github.com/XiaoMi/Gaea/util/mocks/containerTest/builder/containerd/run"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/oci"
	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"syscall"
	"testing"
	"time"
)

// TestRunContainerdEnv 测试容器的运行环境 test containerd env
func TestRunContainerdEnv(t *testing.T) {
	return
	// 创建新的容器的连接客户端 create a new client connected to the default socket path for containerd
	client, err := containerd.New(run.DefaultSock)
	require.Nil(t, err)
	defer func() {
		_ = client.Close()
	}()

	// 测立一个新的命名空间 create a new context with a "default" namespace
	ctx := namespaces.WithNamespace(context.Background(), "default")

	// 拉取预设的测试印象档 pull the default test image from DockerHub
	image, err := client.Pull(ctx, "docker.io/library/debian:latest", containerd.WithPullUnpack)
	require.Nil(t, err)

	// 连接到网路环境 gaea-default connection to the default network environment
	defaultNS := specs.LinuxNamespace{Type: specs.NetworkNamespace, Path: "/var/run/netns/gaea-default"}

	// 创建一个新的预设容器 create a default container (离开后移除 remove after leaving the test)
	container, err := client.NewContainer(
		ctx,
		"default-server",
		containerd.WithImage(image),
		containerd.WithNewSnapshot("default-server-snapshot", image),
		containerd.WithNewSpec(oci.WithImageConfig(image), oci.WithLinuxNamespace(defaultNS)), // 在这里添加网路环境 add the network environment
	)
	require.Nil(t, err)
	defer func() {
		_ = container.Delete(ctx, containerd.WithSnapshotCleanup)
	}()

	// 创建新的容器工作 create a task from the container (离开后移除 remove after leaving the test)
	task, err := container.NewTask(ctx, cio.NewCreator(cio.WithStdio))
	require.Nil(t, err)
	defer func() {
		_, err = task.Delete(ctx)
	}()

	// 等待容器工作完成 make sure we wait before calling start
	exitStatusC, err := task.Wait(ctx)
	require.Nil(t, err)

	// 开始执行容器工作 call start on the task to execute the server
	err = task.Start(ctx)
	require.Nil(t, err)

	// 等待 3 秒 sleep for 3 seconds
	time.Sleep(30 * 60 * time.Second)

	// 删除容器和获得离开讯息 kill the process and get the exit status
	err = task.Kill(ctx, syscall.SIGKILL)
	require.Nil(t, err)

	// 每待容器全部删除完成 wait for the process to fully exit and print out the exit status
	status := <-exitStatusC
	code, _, err := status.Result()
	require.Nil(t, err)
	require.Equal(t, code, uint32(0x89))
}

// TestDefaultContainerd 使用约定的接口进行测试 test default containerd interface
func TestDefaultContainerd(t *testing.T) {
	return
	// 测试约定的接口 test the default interface
	t.Run("test default interface", func(t *testing.T) {
		// 创建新的容器的连接客户端 create a new client connected to the default socket path for containerd
		client, err := containerd.New(run.DefaultSock)
		require.Nil(t, err)
		defer func() {
			_ = client.Close()
		}()

		// 测立一个新的命名空间 create a new context with a "default" namespace
		ctx := namespaces.WithNamespace(context.Background(), "default")

		// 创建测试对象 create a test object
		d := Defaults{
			debianVersion: "latest",
		}

		// 拉取预设的测试印象档 pull the default test image from DockerHub
		img, err := d.Pull(client, ctx, "docker.io/library/debian:latest")
		assert.Nil(t, err)

		// 创建一个新的预设容器 create a default container
		c, err := d.Create(client, ctx, "default-server", "/var/run/netns/gaea-default", img, "default-server-snapshot")
		assert.Nil(t, err)

		// 创建新的容器工作 create a task from the container
		tsk, err := d.Task(c, ctx)
		assert.Nil(t, err)

		// start the task. 开始执行容器工作
		err = d.Start(tsk, ctx)
		assert.Nil(t, err)

		// interrupt the task. 强制终止容器工作
		err = d.Interrupt(tsk, ctx)
		assert.Nil(t, err)

		// 删除容器和获得离开讯息 kill the process and get the exit status
		err = d.Delete(tsk, c, ctx)
		require.Nil(t, err)
	})
	// 测试非约定的函数 test the non-default function
	t.Run("test non-default function", func(t *testing.T) {
		// 先创建接口去进行储存和操作 create a new interface to store and operate
		// store := new(ContainerdClient)

		// 创建一个待储存的对象 create a test object
		/*d := &defaults{
			debianVersion: "9",
		}
		store.Run = d*/

		// 进行储存和操作 store and operate
		/*d.Set("10")
		require.Equal(t, store.Run.(*defaults).Version(), "10")
		store.Run.(*defaults).Set("11")
		require.Equal(t, d.Version(), "11")*/
	})
}
