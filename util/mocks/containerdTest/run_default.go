package containerdTest

import (
	"context"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/oci"
	"github.com/opencontainers/runtime-spec/specs-go"
	"syscall"
	"time"
)

// defined data 数据定义
type defaults struct {
	debianVersion string // debain version. debian 版本
}

// defined interface 约定的接口

// Pull is to pull image from registry. 为容器拉取镜像
func (d *defaults) Pull(client *containerd.Client, ctx context.Context, imageUrl string) (containerd.Image, error) {

	// pull image from registry. 从注册中心拉取镜像
	download := func(client *containerd.Client, ctx context.Context, imageUrl string) (containerd.Image, error) {
		image, err := client.Pull(ctx, imageUrl, containerd.WithPullUnpack)
		if err != nil {
			return nil, err
		}
		return image, nil
	}

	// message. 消息
	type message struct {
		image containerd.Image
		err   error
	}

	// channel. 通道
	chMessage := make(chan message)

RETRY:
	// goroutine to pull image from registry. 开启一个goroutine从注册中心拉取镜像
	go func(client *containerd.Client, ctx context.Context, imageUrl string) {
		image, err := download(client, ctx, imageUrl)
		chMessage <- message{image: image, err: err} // return the message. 回传消息
	}(client, ctx, imageUrl)

	// wait for image pull. 等待镜像拉取
	for {
		select {
		case <-ctx.Done():
			// stop downloading the image because the context is canceled. 逾时停止下载镜像
			return nil, ctx.Err()
		case downloadMsg := <-chMessage:
			// download the image failed. 下载镜像失败
			if downloadMsg.err != nil {
				goto RETRY
			}
			// download the image successfully. 下载镜像成功
			return downloadMsg.image, downloadMsg.err
		default:
			// wait for a while. 等待一段时间
			time.Sleep(1 * time.Second)
		}
	}
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

	// create a task from the container. 建立新的容器工作
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
			// wait for a while. 等待一段时间
			time.Sleep(1 * time.Second)
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
func (d *defaults) Delete(task containerd.Task, container containerd.Container, ctx context.Context) error {

	// delete the task. 刪除容器工作
	if task != nil { // check task exist. 确认容器工作存在
		if _, err := task.Delete(ctx); err != nil {
			return err
		}
	}

	// delete the container. 刪除容器
	if container != nil { // check container exist. 确认容器存在
		if err := container.Delete(ctx, containerd.WithSnapshotCleanup); err != nil {
			return err
		}
	}

	// delete success. 刪除成功
	return nil
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
