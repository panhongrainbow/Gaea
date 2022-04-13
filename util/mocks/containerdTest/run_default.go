package containerdTest

import (
	"context"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/oci"
	"github.com/opencontainers/runtime-spec/specs-go"
	"syscall"
)

// defined data 数据定义
type defaults struct {
	debianVersion string // debain version. debian 版本
}

// defined interface 约定的接口

// Pull is to pull image from registry. 为容器拉取镜像
func (d *defaults) Pull(client *containerd.Client, ctx context.Context, imageUrl string) (containerd.Image, error) {
	// create a new client connected to the default socket path for containerd. 建立新的容器的连接客户端
	return client.Pull(ctx, imageUrl, containerd.WithPullUnpack)
}

// Create is to create container. 为容器创建
func (d *defaults) Create(client *containerd.Client, ctx context.Context, containerName string, networkNS string, imagePulled containerd.Image, snapShot string) (containerd.Container, error) {
	// gaea-default connection to the default network environment. 连接到网路环境
	defaultNS := specs.LinuxNamespace{Type: specs.NetworkNamespace, Path: networkNS}

	// create a default container. 建立一个新的预设容器
	return client.NewContainer(
		ctx,
		containerName,
		containerd.WithImage(imagePulled),
		containerd.WithNewSnapshot(snapShot, imagePulled),
		containerd.WithNewSpec(oci.WithImageConfig(imagePulled), oci.WithLinuxNamespace(defaultNS)),
	)
}

// Task is to create task. 为容器任务创建
func (d *defaults) Task(container containerd.Container, ctx context.Context) (containerd.Task, error) {

	// create a task from the container (离开后移除 remove after leaving the test). 建立新的容器工作
	return container.NewTask(ctx, cio.NewCreator(cio.WithStdio))
}

// Start is to start task. 为容器任务启动
func (d *defaults) Start(task containerd.Task, ctx context.Context) error {

	// call start on the task to execute the server. 开始执行容器工作
	if err := task.Start(ctx); err != nil {
		return err
	}

	// check the status of the task. 检查容器任务状态
LOOP:
	for {
		select {
		case <-ctx.Done():
			// stop the task because the context is canceled. 逾时停止容器工作
			return ctx.Err()
		default:
			// start listening for the container work 開始監聽容器工作
			status, err := task.Status(ctx)
			if err != nil {
				// container monitor failed. 容器工作监听失败
				return err
			}
			if status.Status == containerd.Running {
				// container is not running. 容器工作不为正在运行
				break LOOP
			}
		}
	}

	// 容器工作执行成功 container work executed successfully
	return nil
}

// Interrupt is to stop task immediately. 为立刻停止容器任务
func (d *defaults) Interrupt(task containerd.Task, ctx context.Context) error {
	// kill the process work. 删除容器工作
	return task.Kill(ctx, syscall.SIGKILL)
}

// Delete is to delete task. 为容器任务停止
func (d *defaults) Delete(container containerd.Container, ctx context.Context) error {
	return container.Delete(ctx, containerd.WithSnapshotCleanup)
}

// non-defined interface 非约定的函数

// Set is to set current version. 获取当前版本
func (d *defaults) Set(version string) {
	d.debianVersion = version
}

// Version is to get current version. 获取当前版本
func (d *defaults) Version() string {
	return d.debianVersion
}
