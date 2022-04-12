package containerdTest

import (
	"context"
	"github.com/containerd/containerd"
)

// defined data 数据定义
type etcd struct {
	Version string // database version. 数据库版本
}

// defined interface 约定的接口

// Pull is to pull image from registry. 为容器拉取镜像
func (m *etcd) Pull(client *containerd.Client, ctx context.Context, imageUrl string) (containerd.Image, error) {
	return nil, nil
}

// Create is to create container. 为容器创建
func (m *etcd) Create(client *containerd.Client, ctx context.Context, containerName string, networkNS string, imagePulled containerd.Image, snapShot string) (containerd.Container, error) {
	return nil, nil
}

// Task is to create task. 为容器任务创建
func (m *etcd) Task(container string) (string, error) {
	return "", nil
}

// Start is to start task. 为容器任务启动
func (m *etcd) Start(container string) error {
	return nil
}

// Stop is to stop task. 为容器任务停止
func (m *etcd) Stop(container string) error {
	return nil
}

// Delete is to delete task. 为容器任务停止
func (m *etcd) Delete(container containerd.Container, ctx context.Context) error {
	return container.Delete(ctx, containerd.WithSnapshotCleanup)
}
