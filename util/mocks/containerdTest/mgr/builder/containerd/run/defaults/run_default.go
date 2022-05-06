package defaults

import (
	"context"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/oci"
	_ "github.com/go-sql-driver/mysql"
	"github.com/opencontainers/runtime-spec/specs-go"
	"syscall"
	"time"
)

// Defined data 数据定义
type Defaults struct {
	debianVersion string // debain version. debian 版本
}

// defined interface 约定的接口

// >>>>> >>>>> >>>>> 创建部分

// Pull is to pull image from registry. 为容器拉取镜像
func (d *Defaults) Pull(client *containerd.Client, ctx context.Context, imageUrl string) (containerd.Image, error) {

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
func (d *Defaults) Create(client *containerd.Client, ctx context.Context, containerName string, networkNS string, imagePulled containerd.Image, snapShot string) (containerd.Container, error) {
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

// Task 为容器任务创建 Task is to create task.
func (d *Defaults) Task(container containerd.Container, ctx context.Context) (containerd.Task, error) {
	// 建立新的容器工作 create a task from the container.
	return container.NewTask(ctx, cio.NewCreator(cio.WithStdio))
}

// Start 为容器任务启动 Start is to start task.
func (d *Defaults) Start(task containerd.Task, ctx context.Context) error {
	// 开始执行容器工作 call start on the task to execute the server.
	if err := task.Start(ctx); err != nil {
		return err
	}

	// 检查容器任务状态 check the status of the task.
LOOP:
	for {
		select {
		case <-ctx.Done():
			// 逾时停止容器工作 stop the task because the context is canceled.
			return ctx.Err()
		default:
			// 開始監聽容器工作 start listening for the container work.
			status, err := task.Status(ctx)
			if err != nil {
				// 容器工作监听失败 container monitor failed.
				return err
			}
			if status.Status == containerd.Running {
				// 容器工作正在运行 container work is running.
				break LOOP
			}

			// 等待一秒 wait for one second.
			time.Sleep(1 * time.Second)
		}
	}

	// 容器工作执行成功 container work executed successfully.
	return nil
}

// >>>>> >>>>> >>>>> 检查部分

// CheckService 为检查容器服务是否上线 CheckService is to check container service.
func (d *Defaults) CheckService(ctx context.Context, ipAddrPort string) error {
	// 預設容器沒有服務，所以不需要檢查.
	// Default container has no service, so no need to check.

	// wait for one second. 等待一秒
	time.Sleep(1 * time.Second)

	// 正常返回. return successfully.
	return nil
}

// CheckSchema 为检查容器资料是否存在 CheckService is to check container data exists.
func (d *Defaults) CheckSchema(ctx context.Context, ipAddrPort string) error {
	// 預設容器沒有服務，所以不需要檢查.
	// Default container has no service, so no need to check.

	// wait for one second. 等待一秒
	time.Sleep(1 * time.Second)

	// 正常返回. return successfully.
	return nil
}

// >>>>> >>>>> >>>>> 删除部分

// Interrupt is to stop task immediately. 为立刻停止容器任务
func (d *Defaults) Interrupt(task containerd.Task, ctx context.Context) error {
	// stop task immediately. 停止任务
LOOP:
	for {
		// stop the task. 停止任务
		_ = task.Kill(ctx, syscall.SIGKILL)

		// start listening for the container status. 開始監聽容器状态
		status, err := task.Status(ctx)
		if err != nil {
			// container monitor failed. 容器工作监听失败
			return err
		}
		if status.Status != containerd.Running {
			// container is not running. 容器工作不为正在运行
			break LOOP
		}
		// wait for a while. 等待一段时间
		time.Sleep(1 * time.Second)
	}

	// stop task successfully. 容器工作中止成功
	return nil
}

// Delete is to delete task. 为容器任务停止
func (d *Defaults) Delete(task containerd.Task, container containerd.Container, ctx context.Context) error {

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
func (d *Defaults) Set(version string) {
	d.debianVersion = version
}

// Version is to get current version. 获取当前版本
func (d *Defaults) Version() string {
	return d.debianVersion
}
