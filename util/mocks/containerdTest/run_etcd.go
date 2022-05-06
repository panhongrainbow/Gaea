package containerdTest

// defined data 数据定义
type etcd struct {
	defaults
}

// defined interface 约定的接口

// 约定的接口 defined interface.

// >>>>> >>>>> >>>>> 创建部分

// Pull 为容器拉取镜像 Pull is to pull image from registry.
/*func (e *etcd) Pull(client *containerd.Client, ctx context.Context, imageUrl string) (containerd.Image, error) {
	return nil, nil
}*/

// Create 为容器创建 Create is to create container.
/*func (e *etcd) Create(client *containerd.Client, ctx context.Context, containerName string, networkNS string, imagePulled containerd.Image, snapShot string) (containerd.Container, error) {
	return nil, nil
}*/

// Task 为容器任务创建 Task is to create task.
/*func (e *etcd) Task(container containerd.Container, ctx context.Context) (containerd.Task, error) {
	return nil, nil
}*/

// Start 为容器任务启动 Start is to start task.
/*func (e *etcd) Start(task containerd.Task, ctx context.Context) error {
	return nil
}*/

// >>>>> >>>>> >>>>> 检查部分

// CheckService 为检查容器服务是否上线 CheckService is to check container service.
/*func (e *etcd) CheckService(task containerd.Task, ctx context.Context) error {
	return nil
}*/

// CheckSchema 为检查容器资料是否存在 CheckService is to check container data exists.
/*func (e *etcd) CheckSchema(ctx context.Context, ipAddrPort string) error {
	return nil
}*/

// >>>>> >>>>> >>>>> 删除部分

// Interrupt 为立刻停止容器任务 Interrupt is to stop task immediately.
/*func (e *etcd) Interrupt(task containerd.Task, ctx context.Context) error {
	// kill the process work. 删除容器工作
	return task.Kill(ctx, syscall.SIGKILL)
}*/

// Delete 为容器任务停止 Delete is to delete task.
/*func (e *etcd) Delete(task containerd.Task, container containerd.Container, ctx context.Context) error {
	return container.Delete(ctx, containerd.WithSnapshotCleanup)
}*/
