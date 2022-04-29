package containerdTest

// MariaDB 为数据定义 MariaDB defined data.
type MariaDB struct {
	defaults
}

// 约定的接口 defined interface.

// Pull 为容器拉取镜像 Pull is to pull image from registry.
/*func (m *MariaDB) Pull(client *containerd.Client, ctx context.Context, imageUrl string) (containerd.Image, error) {
	return nil, nil
}*/

// Create 为容器创建 Create is to create container.
/*func (m *MariaDB) Create(client *containerd.Client, ctx context.Context, containerName string, networkNS string, imagePulled containerd.Image, snapShot string) (containerd.Container, error) {
	// gaea-default connection to the default network environment. 连接到网路环境
	mariadbNS := specs.LinuxNamespace{Type: specs.NetworkNamespace, Path: networkNS}

	// 建立一个新的数据库容器 create a default container.
	return client.NewContainer(
		ctx,
		containerName,
		containerd.WithImage(imagePulled),
		containerd.WithNewSnapshot(snapShot, imagePulled),
		containerd.WithNewSpec(
			oci.WithImageConfig(imagePulled),
			oci.WithLinuxNamespace(mariadbNS),
			// 之后可以在这里加环境变数 add environment variables later.
			// oci.WithEnv([]string{"MYSQL_ROOT_PASSWORD=12345", "MYSQL_USER=xiaomi", "MYSQL_PASSWORD=12345"}),
		),
	)
}*/

// Task 为容器任务创建 Task is to create task.
/*func (m *MariaDB) Task(container containerd.Container, ctx context.Context) (containerd.Task, error) {
	return nil, nil
}*/

// Start 为容器任务启动 Start is to start task.
/*func (m *MariaDB) Start(task containerd.Task, ctx context.Context) error {
	return nil
}*/

// Interrupt 为立刻停止容器任务 Interrupt is to stop task immediately.
/*func (m *MariaDB) Interrupt(task containerd.Task, ctx context.Context) error {
	// kill the process work. 删除容器工作
	return task.Kill(ctx, syscall.SIGKILL)
}*/

// Delete 为容器任务停止 Delete is to delete task.
/*func (m *MariaDB) Delete(task containerd.Task, container containerd.Container, ctx context.Context) error {
	return container.Delete(ctx, containerd.WithSnapshotCleanup)
}*/
