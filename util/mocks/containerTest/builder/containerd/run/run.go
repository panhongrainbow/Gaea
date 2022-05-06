package run

import (
	"context"
	"github.com/containerd/containerd"
)

const (
	DefaultSock = "/run/containerd/containerd.sock" // default sock path. 默认的 sock 路径
)

// 容器的执行状态 Containerd's run status.
const (
	ContainerdStatusInit                 = iota // 初始化状态 init status
	ContainerdStatusOccupied                    // 被占领的状态 running status
	ContainerdStatusBuilding                    // 容器在创建状态 build status
	ContainerdStatusBuildPullingImage           // 下载镜像 pull image
	ContainerdStatusBuildCreateContainer        // 创建容器 create container
	ContainerdStatusBuildCreateTask             // 创建任务 create task
	ContainerdStatusBuildStartTask              // 启动任务 start task
	ContainerdStatusBuildRunning                // 容器运行中 running
	ContainerdStatusChecking                    // 检查服务状态中 checking containerd status
	ContainerdStatusCheckOnService              // 检查服务状态中 check containerd status is on service
	ContainerdStatusCheckSchema                 // 检查数据库服务 check containerd schema.
	ContainerdStatusReady                       // 数据库服务可以开始服务 containerd container is ready.
	ContainerdStatusTearingDown                 // 容器在拆除的状态 tear down status
	ContainerdStatusTearDownInterrupted         // 被中断的状态 interrupted status
	ContainerdStatusTearDownKilled              // 容器被杀死 killed
	ContainerdStatusError                       // 容器服务错误 containerd error
	ContainerdStatusReturned                    // 容器服务适放 containerd returned
)

// Run 接口会容器對象，在這裡要可以直接操作容器 Run is an interface for containerd client to implement.
type Run interface {
	// Pull to Start 创建部份 create part
	Pull(client *containerd.Client, ctx context.Context, imageUrl string) (containerd.Image, error)
	Create(client *containerd.Client, ctx context.Context, containerName string, networkNS string, imagePulled containerd.Image, snapShot string) (containerd.Container, error)
	Task(container containerd.Container, ctx context.Context) (containerd.Task, error)
	Start(task containerd.Task, ctx context.Context) error

	// CheckService to CheckSchema 检查部份 check part
	CheckService(ctx context.Context, ipAddrPort string) error
	CheckSchema(ctx context.Context, ipAddrPort string) error

	// Interrupt to Delete 销毁部份 destroy part
	Interrupt(task containerd.Task, ctx context.Context) error
	Delete(task containerd.Task, container containerd.Container, ctx context.Context) error
}
