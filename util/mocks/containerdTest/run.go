package containerdTest

import (
	"context"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
)

const (
	defaultSock = "/run/containerd/containerd.sock" // default sock path. 默认的 sock 路径
)

const (
	containerdStatusInit            = iota // init status. 初始化状态
	containerdStatusLoadImage              // load image. 加载镜像
	containerdStatusCreateContainer        // create container. 创建容器
	containerdStatusCreateTask             // create task. 创建任务
	containerdStatusStartTask              // start task. 启动任务
	containerdStatusRunning                // running. 容器运行中
	containerdStatusStopped                // stopped. 容器停止
	containerdStatusKilled                 // killed. 容器被杀死
	containerdStatusError                  // containerd error. 容器服务错误
)

// ContainerdClient 容器服务的客户端
type ContainerdClient struct {
	// create in NewContainerdClient. 在 NewContainerdClient 中创建
	Status    int                // 容器服务的状态
	Conn      *containerd.Client // 容器服务的连接
	Type      string             // 容器类型，是 etcd 或者 mariaDB 等等
	IP        string             // 容器服务的 IP
	Container ClientContainerd   // 容器服务的容器
	Schema    ClientSchema       // 容器服务的 Schema
	Running   *ClientRunning     // 容器服务是否运行中

	// create in Distinguish. 在容器区分时候创建
	Run Run // 容器服务的运行的接口
}

// ClientContainerd 客戶端的容器服务設定
type ClientContainerd struct {
	Name      string
	NameSpace string
	Image     string
	SnapShot  string
	NetworkNS string
	Container string
	Task      string
}

// ClientSchema 客戶端的容器服务設定，之后有新类型的容器服务，需要在此增加
type ClientSchema struct {
	User     string
	Password string
	Schema   string
}

type ClientRunning struct {
	ctx context.Context
	img containerd.Image
	tsk containerd.Task
	c   containerd.Container
}

// NewContainerdClient is a function to create a new containerd client. 新建容器服务的客户端
func NewContainerdClient(cfg ContainerD) (*ContainerdClient, error) {
	currentSock := new(containerd.Client) // it's a connection to containerd 新的连接
	var err error                         // error message 报错信息
	var usedSock = ""                     // usedSock is the user defined sock path 客戶端的 sock 指定路徑

	// if socketPath is empty, use default path. 如果没有配置，则使用默认的路径

	if cfg.Sock == "" { // use default path. 使用默认的路径
		usedSock = defaultSock
	}
	if cfg.Sock != "" { // use the specified path. 使用指定的路径
		usedSock = cfg.Sock
	}

	// create a new containerd connection. 建立容器服务的客户端连接
	currentSock, err = containerd.New(usedSock)
	if err != nil {
		return nil, err
	}

	// create a new containerd client. 建立容器服务的客户端
	client := &ContainerdClient{
		// create in NewContainerdClient. 在 NewContainerdClient 中创建

		// client's config. 客户端的配置
		Status: containerdStatusInit, // It's init status. 现在为初始化状态
		Conn:   currentSock,          // It's a connection to containerd. 容器服务的连接
		Type:   cfg.Type,             // It's a containerd type. 容器服务的类型
		IP:     cfg.IP,               // It's a containerd IP. 容器服务的 IP
		// container's config. 容器的配置
		Container: ClientContainerd{
			NameSpace: cfg.NameSpace, // It's a container name. 容器服务的名称
			Image:     cfg.Image,     // It's a container image. 容器服务的镜像
		},
		// Schema's config. Schema的配置
		Schema: ClientSchema{
			User: cfg.User, // It's a containerd user. 容器服务的用户
		},
	}

	// create in Distinguish. 在容器的区分时候创建
	err = Distinguish(client) // implement in Distinguish. 实现在 Distinguish 中
	if err != nil {
		return nil, err
	}

	// create in Execution. 在容器執行时候创建
	// client.Ctx = nil // context. 容器服务的上下文

	// return the new containerd client. 返回新的容器服务的客户端
	return client, nil
}

// Distinguish is a function to distinguish the containerd client. 对容器服务进行区分，判断容器服务的类型，给容器服务所需要的功能
func Distinguish(client *ContainerdClient) error {
	// distinguish the containerd client. 对容器服务进行区分
	switch client.Type {
	case "etcd":
		client.Run = new(etcd) // use etcd. 容器服务为 etcd
		return nil             // return nil. 返回 nil
	case "mariadb": // use mariaDB. 容器服务为 mariaDB
		client.Run = new(MariaDB) // return mariaDB. 返回 mariaDB
	default:
		client.Run = new(defaults) // use defaults. 容器服务为 defaults
	}

	// return the error. 返回错误
	return nil
}

//
func (cc *ContainerdClient) Build() error {
	// create Running object. 创建 Running 对象
	cc.Running = new(ClientRunning)

	// 错误信息 error message.
	var err error

	// 测立一个新的命名空间 create a new context with a "mariadb" namespace
	cc.Running.ctx = namespaces.WithNamespace(context.Background(), cc.Container.NameSpace)

	// 拉取预设的测试印象档 pull the default test image from DockerHub
	// example: "docker.io/panhongrainbow/mariadb:testing" OR "localhost/mariadb:latest"
	cc.Running.img, err = cc.Run.Pull(cc.Conn, cc.Running.ctx, cc.Container.Image)
	if err != nil {
		return err
	}

	// 建立一个新的容器 create a new container
	cc.Running.c, err = cc.Run.Create(cc.Conn, cc.Running.ctx, cc.Container.Name, cc.Container.NetworkNS, cc.Running.img, cc.Container.SnapShot)
	if err != nil {
		return err
	}

	// 建立新的容器工作 create a task from the container
	cc.Running.tsk, err = cc.Run.Task(cc.Running.c, cc.Running.ctx)
	if err != nil {
		return err
	}

	// 開始执行容器工作 start the task.
	err = cc.Run.Start(cc.Running.tsk, cc.Running.ctx)
	if err != nil {
		return err
	}

	// 建立容器環境成功 Build the container environment successfully.
	return nil
}

//
func (cc *ContainerdClient) TearDown() error {
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

	// 删除容器環境成功 delete the container environment successfully.
	return nil
}
