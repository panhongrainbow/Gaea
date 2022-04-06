package containerdTest

import (
	"context"
	"github.com/containerd/containerd"
)

// init 初始化 Conainerd 容器服务
func init() {
	// 先开始辨判连接方式
	// err := setup()
}

// Run is an interface for containerd client to implement. 容器對象，在這裡要可以直接操作容器
type Run interface {
	Pull(client *containerd.Client, ctx context.Context, image string) error
	Create(container string, image string) error
	Task(container string) (string, error)
	Start(container string) error
	Stop(container string) error
	Delete(container string) error
}
