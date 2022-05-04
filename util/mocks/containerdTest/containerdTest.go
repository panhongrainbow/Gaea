package containerdTest

import (
	"context"
	"fmt"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
	"net"
	"strings"
	"time"
)

// ContainerdClient is core component of Containerd client. 容器服务的核心客户端
type ContainerdClient struct {
	// 在 NewContainerdClient 中创建 create in NewContainerdClient.
	Status    int                // 容器服务的状态
	Conn      *containerd.Client // 容器服务的连接
	Type      string             // 容器类型，是 etcd 或者 mariaDB 等等
	IP        string             // 容器服务的 IP
	Container ClientContainerd   // 容器服务的容器
	Schema    ClientSchema       // 容器服务的 Schema
	Running   *ClientRunning     // 容器服务是否运行中

	// 在容器区分时候创建 create in Distinguish.
	Run Run // 容器服务的运行的接口 interface for Containerd.
}

// ClientContainerd 为客戶端的容器服务設定 containerd is configured for Containerd.
type ClientContainerd struct {
	Name      string // 容器服务的名称 Name
	NameSpace string // 容器服务的命名空间 NameSpace
	Image     string // 容器服务的镜像 Image
	SnapShot  string // 容器服务的快照 SnapShot
	NetworkNS string // 容器服务的网络命名空间 NetworkNS
	Container string // 容器服务的容器 Container
	Task      string // 容器服务的任务 Task
}

// ClientSchema 客戶端的 Schema 設定 ClientSchema is the schema of the containerd client.
type ClientSchema struct {
	User     string // 容器服务的用户名 user name.
	Password string // 容器服务的密码 password.
	Schema   string // 容器服务的 schema.
}

// ClientRunning 客戶端的运行时的对象 ClientRunning is the running object of the containerd client.
type ClientRunning struct {
	ctx    context.Context      // 容器服务的上下文 context.
	img    containerd.Image     // 容器服务的镜像 image.
	c      containerd.Container // 容器服务的容器 container.
	tsk    containerd.Task      // 容器服务的任务 task.
	cancel context.CancelFunc   // 容器服务的取消函数 cancel function.
}

// NewContainerdClient 为新建容器服务的客户端 NewContainerdClient is a function to create a new containerd client.
func NewContainerdClient(cfg ContainerD) (*ContainerdClient, error) {
	// >>>>> >>>>> >>>>> 决定容器服务的连接 sock 对象 decide the sock.

	// 创建容器服务的客户端 Socket 对象 create a new sock for the containerd client.
	currentSock := new(containerd.Client) // 新的连接 it's a connection to containerd.
	var err error                         // 报错信息 error message
	var usedSock = ""                     // 客戶端的 sock 指定路徑 usedSock is the user defined sock path.

	// 如果没有配置，则使用默认的路径 if socketPath is empty, use default path.
	if cfg.Sock != "" { // 使用指定的路径 use the specified path.
		usedSock = cfg.Sock
	} else { // 使用默认的路径 use default path.
		usedSock = defaultSock
	}

	// 建立容器服务的客户端连接 create a new containerd connection.
	currentSock, err = containerd.New(usedSock)
	if err != nil {
		return nil, err
	}

	// >>>>> >>>>> >>>>> 创建容器服务的 ContainerdClient 对象 create a new ContainerdClient object

	// 建立容器服务的客户端 create a new containerd client.
	client := &ContainerdClient{
		// 在 NewContainerdClient 中创建 create in NewContainerdClient.

		// 客户端的配置 client's config.
		Status: containerdStatusInit, // 现在为初始化状态 It's init status.
		Conn:   currentSock,          // 容器服务的连接 It's a connection to containerd.
		Type:   cfg.Type,             // 容器服务的类型 It's a containerd type.
		IP:     cfg.IP,               // 容器服务的网路位置 It's a containerd IP.
		// 容器的配置 container's config.
		Container: ClientContainerd{
			Name:      cfg.Name,      // 容器的名称 It's a container's name.
			NameSpace: cfg.NameSpace, // 容器服务的名称 It's a container name.
			Image:     cfg.Image,     // 容器服务的镜像 It's a container image.
			NetworkNS: cfg.NetworkNs, // 容器服务的网络命名空间 It's a container network namespace.
			SnapShot:  cfg.SnapShot,  // 容器服务的快照 It's a container snapshot.
		},
		// Schema的配置 Schema's config.
		Schema: ClientSchema{
			User: cfg.User, // 容器服务的用户 It's a containerd user.
		},
	}

	// >>>>> >>>>> >>>>> 实现 Run 接口并回传 implement the Run interface and return

	// 实现在 Distinguish 中 implement in Distinguish.
	err = client.Distinguish()
	if err != nil {
		return nil, err
	}

	// 返回新的容器服务的客户端 return the new containerd client.
	return client, nil
}

// Distinguish 对容器服务进行区分，判断容器服务的类型，给容器服务所需要的功能 Distinguish is a function to distinguish the containerd client.
func (cc *ContainerdClient) Distinguish() error {
	// distinguish the containerd client. 对容器服务进行区分
	switch cc.Type {
	case "etcd":
		cc.Run = new(etcd) // use etcd. 容器服务为 etcd
		return nil         // return nil. 返回 nil
	case "mariadb": // use mariaDB. 容器服务为 mariaDB
		cc.Run = new(MariaDB) // return mariaDB. 返回 mariaDB
		return nil            // return nil. 返回 nil
	default:
		cc.Run = new(defaults) // use defaults. 容器服务为 defaults
		return nil             // return nil. 返回 nil
	}
}

// Build 建立容器测试环境 Build create a new container environment for test.
func (cc *ContainerdClient) Build(t time.Duration) error {
	// 错误信息 error message.
	var err error

	// 设置容器管理器为创建状态 container manager's status is building.
	if err = SetContainerManagerStatus(cc.Container.Name, containerdStatusBuilding); err != nil {
		return err
	}

	// create Running object. 创建 Running 对象
	cc.Running = new(ClientRunning)

	// 决定之后的决行时间 decide the time duration.
	ctx := context.Background() // 创建一个上下文对象 create a context object.
	if t > 0 {
		ctx, cc.Running.cancel = context.WithTimeout(ctx, t)
	}

	// 测立一个新的命名空间 create a new context with a "mariadb" namespace
	cc.Running.ctx = namespaces.WithNamespace(ctx, cc.Container.NameSpace)

	// 设置容器管理器为下载镜像状态 container manager's status is pulling image.
	if err = SetContainerManagerStatus(cc.Container.Name, containerdStatusBuildPullingImage); err != nil {
		return err
	}

	// 拉取预设的测试印象档 pull the default test image from DockerHub
	// example: "docker.io/panhongrainbow/mariadb:testing" OR "localhost/mariadb:latest"
	cc.Running.img, err = cc.Run.Pull(cc.Conn, cc.Running.ctx, cc.Container.Image)
	if err != nil {
		return err
	}

	// 设置容器管理器为创建容器状态 container manager's status is creating container.
	if err = SetContainerManagerStatus(cc.Container.Name, containerdStatusBuildCreateContainer); err != nil {
		return err
	}

	// 建立一个新的容器 create a new container
	cc.Running.c, err = cc.Run.Create(cc.Conn, cc.Running.ctx, cc.Container.Name, cc.Container.NetworkNS, cc.Running.img, cc.Container.SnapShot)
	if err != nil {
		return err
	}

	// 设置容器管理器为创建容器任务状态 container manager's status is creating container task.
	if err = SetContainerManagerStatus(cc.Container.Name, containerdStatusBuildCreateTask); err != nil {
		return err
	}

	// 建立新的容器工作 create a task from the container
	cc.Running.tsk, err = cc.Run.Task(cc.Running.c, cc.Running.ctx)
	if err != nil {
		return err
	}

	// 设置容器管理器为启动容器任务状态 container manager's status is starting container task.
	if err = SetContainerManagerStatus(cc.Container.Name, containerdStatusBuildStartTask); err != nil {
		return err
	}

	// 開始执行容器工作 start the task.
	err = cc.Run.Start(cc.Running.tsk, cc.Running.ctx)
	if err != nil {
		return err
	}

	// 设置容器管理器为容器任务执行状态 container manager's status is running container task.
	if err = SetContainerManagerStatus(cc.Container.Name, containerdStatusBuildRunning); err != nil {
		return err
	}

	// 建立容器環境成功 Build the container environment successfully.
	return nil
}

// OnService 确认容器是否处放服务状态 OnService is making sure the container is on service.
func (cc *ContainerdClient) OnService(t time.Duration) error {
	for i := 0; i < 20; i++ {
		typ := "tcp"
		if strings.Contains("10.10.10.10:3306", "/") {
			typ = "unix"
		}

		netConn, err := net.Dial(typ, "10.10.10.10:3306")
		if err == nil {

			// 先随意测试
			test := make([]byte, 20)
			netConn.Read(test)
			fmt.Println(test)

			_ = netConn.Close()
			return nil
		}
		fmt.Println("1", err)

		// _ = netConn.Close()

		time.Sleep(time.Second)
	}

	return nil
}

// TearDown 拆除容器测试环境 TearDown is a function to tear down the container environment.
func (cc *ContainerdClient) TearDown(t time.Duration) error {
	// 决定之后的决行时间 decide the time duration.
	ctx := context.Background() // 创建一个上下文对象 create a context object.
	if t > 0 {
		ctx, cc.Running.cancel = context.WithTimeout(ctx, t)
	}

	// 测立一个新的命名空间 create a new context with a "mariadb" namespace
	cc.Running.ctx = namespaces.WithNamespace(ctx, cc.Container.NameSpace)

	// 強制中斷容器工作 interrupt the task.
	err := cc.Run.Interrupt(cc.Running.tsk, cc.Running.ctx)
	if err != nil {
		return err
	}

	// 删除容器和获得离开讯息 kill the process and get the exit status
	err = cc.Run.Delete(cc.Running.tsk, cc.Running.c, cc.Running.ctx)
	if err != nil {
		return err
	}

	defer cc.Running.cancel() // 取消上下文 cancel the context.
	cc.Running = nil          // 清空 Running 对象 clear the Running object.

	// 删除容器環境成功 delete the container environment successfully.
	return nil
}
