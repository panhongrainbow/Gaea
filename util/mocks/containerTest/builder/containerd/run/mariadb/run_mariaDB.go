package mariadb

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/XiaoMi/Gaea/util/mocks/containerTest/builder/containerd/run/defaults"
	"net"
	"strings"
	"time"
)

// MariaDB 为数据定义 MariaDB defined data.
type MariaDB struct {
	defaults.Defaults
}

// 约定的接口 defined interface.

// >>>>> >>>>> >>>>> 创建部分

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

// >>>>> >>>>> >>>>> 检查部分

// CheckService 为检查容器服务是否上线 CheckService is to check container service.
func (m *MariaDB) CheckService(ctx context.Context, ipAddrPort string) error {
	// 检查容器连线设定值 check the container network settings.
	typ := "tcp"
	if strings.Contains(ipAddrPort, "/") {
		typ = "unix"
	}
	for {
		select {
		case <-ctx.Done():
			// 逾时停止容器工作 stop the task because the context is canceled.
			return ctx.Err()
		default:
			// 检查容器服务是否上线 check the container service.
			netConn, err := net.Dial(typ, ipAddrPort)
			if netConn != nil {
				// 如果连线成功，立即关闭连接 close the connection immediately, if it is successful.
				fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>ok0")
				_ = netConn.Close()
			}

			// 如果没有错误，立即返回成功. if there is no error, return success.
			if err == nil {
				fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>ok1")
				return nil
			}

			// 等待一秒 wait for one second.
			time.Sleep(1 * time.Second)
		}
	}
}

// CheckSchema 为检查容器资料是否存在 CheckService is to check container data exists.
func (m *MariaDB) CheckSchema(ctx context.Context, ipAddrPort string) error {
	db, err := sql.Open("mysql", "xiaomi:12345@tcp(10.10.10.10:3306)/mysql")
	if err != nil {
		fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>ok2")
		return err
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	err = db.Ping()
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>ok3")
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>", err)
	return nil
}

// >>>>> >>>>> >>>>> 删除部分

// Interrupt 为立刻停止容器任务 Interrupt is to stop task immediately.
/*func (m *MariaDB) Interrupt(task containerd.Task, ctx context.Context) error {
	// kill the process work. 删除容器工作
	return task.Kill(ctx, syscall.SIGKILL)
}*/

// Delete 为容器任务停止 Delete is to delete task.
/*func (m *MariaDB) Delete(task containerd.Task, container containerd.Container, ctx context.Context) error {
	return container.Delete(ctx, containerd.WithSnapshotCleanup)
}*/
