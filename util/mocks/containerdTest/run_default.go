package containerdTest

import (
	"context"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/oci"
	"github.com/opencontainers/runtime-spec/specs-go"
)

// defined data 数据定义
type defaults struct {
	debianVersion string // debain version. debian 版本
}

// defined interface 约定的接口

// Pull is to pull image from registry. 为容器拉取镜像
func (d *defaults) Pull(client *containerd.Client, ctx context.Context, imageUrl string) (containerd.Image, error) {
	// 建立新的容器的连接客户端 create a new client connected to the default socket path for containerd
	return client.Pull(ctx, imageUrl, containerd.WithPullUnpack)
}

// Create is to create container. 为容器创建
func (d *defaults) Create(client *containerd.Client, ctx context.Context, containerName string, networkNS string, imagePulled containerd.Image, snapShot string) (containerd.Container, error) {
	// 连接到网路环境 gaea-default connection to the default network environment
	defaultNS := specs.LinuxNamespace{Type: specs.NetworkNamespace, Path: networkNS}

	// 建立一个新的预设容器 create a default container
	return client.NewContainer(
		ctx,
		containerName,
		containerd.WithImage(imagePulled),
		containerd.WithNewSnapshot(snapShot, imagePulled),
		containerd.WithNewSpec(oci.WithImageConfig(imagePulled), oci.WithLinuxNamespace(defaultNS)),
	)
}

// Task is to create task. 为容器任务创建
func (d *defaults) Task(container string) (string, error) {
	return "", nil
}

// Start is to start task. 为容器任务启动
func (d *defaults) Start(container string) error {
	return nil
}

// Stop is to stop task. 为容器任务停止
func (d *defaults) Stop(container string) error {
	return nil
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
