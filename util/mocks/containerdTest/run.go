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
	containerdStatusInit                 = iota // 初始化状态 init status
	containerdStatusOccupied                    // 被占领的状态 running status
	containerdStatusBuilding                    // 容器在创建状态 build status
	containerdStatusBuildPullingImage           // 下载镜像 pull image
	containerdStatusBuildCreateContainer        // 创建容器 create container
	containerdStatusBuildCreateTask             // 创建任务 create task
	containerdStatusBuildStartTask              // 启动任务 start task
	containerdStatusBuildRunning                // 容器运行中 running
	containerdStatusTearDown                    // 容器在拆除的状态 tear down status
	containerdStatusTearDownInterrupted         // 被中断的状态 interrupted status
	containerdStatusTearDownKilled              // 容器被杀死 killed
	containerdStatusError                       // 容器服务错误 containerd error
	containerdStatusReturned                    // 容器服务适放 containerd returned
)

// Run 接口会容器對象，在這裡要可以直接操作容器 Run is an interface for containerd client to implement.
type Run interface {
	// Pull to Start 创建部份 create part
	Pull(client *containerd.Client, ctx context.Context, imageUrl string) (containerd.Image, error)
	Create(client *containerd.Client, ctx context.Context, containerName string, networkNS string, imagePulled containerd.Image, snapShot string) (containerd.Container, error)
	Task(container containerd.Container, ctx context.Context) (containerd.Task, error)
	Start(task containerd.Task, ctx context.Context) error
	// CheckService to CheckData 检查部份 check part
	CheckService(task containerd.Task, ctx context.Context) error
	CheckData(task containerd.Task, ctx context.Context) error
	// Interrupt to Delete 销毁部份 destroy part
	Interrupt(task containerd.Task, ctx context.Context) error
	Delete(task containerd.Task, container containerd.Container, ctx context.Context) error
}
