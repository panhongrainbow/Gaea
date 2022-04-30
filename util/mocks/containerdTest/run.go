package containerdTest

import (
	"context"
	"github.com/containerd/containerd"
)

const (
	defaultSock = "/run/containerd/containerd.sock" // default sock path. 默认的 sock 路径
)

// 容器的执行状态 Containerd's run status.
const (
	containerdStatusInit            = iota // 初始化状态 init status
	containerdStatusLoadImage              // 加载镜像 load image
	containerdStatusCreateContainer        // 创建容器 create container
	containerdStatusCreateTask             // 创建任务 create task
	containerdStatusStartTask              // 启动任务 start task
	containerdStatusRunning                // 容器运行中 running
	containerdStatusStopped                // 容器停止 stopped
	containerdStatusKilled                 // 容器被杀死 killed
	containerdStatusError                  // 容器服务错误 containerd error
)

// Run 接口会容器對象，在這裡要可以直接操作容器 Run is an interface for containerd client to implement.
type Run interface {
	Pull(client *containerd.Client, ctx context.Context, imageUrl string) (containerd.Image, error)
	Create(client *containerd.Client, ctx context.Context, containerName string, networkNS string, imagePulled containerd.Image, snapShot string) (containerd.Container, error)
	Task(container containerd.Container, ctx context.Context) (containerd.Task, error)
	Start(task containerd.Task, ctx context.Context) error
	Interrupt(task containerd.Task, ctx context.Context) error
	Delete(task containerd.Task, container containerd.Container, ctx context.Context) error
}
