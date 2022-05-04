package containerdTest

import (
	"errors"
	"github.com/XiaoMi/Gaea/log"
	"github.com/containerd/containerd"
	"os"
	"strings"
	"sync"
)

// 先進行 Parse 操作

// ContainerD 为 Containerd 容器服务所要 Parse 对应的内容
type ContainerD struct {
	Sock      string `json:"sock"`      // 容器服务在 Linux 的连接位置
	Type      string `json:"type"`      // 容器服务类型
	Name      string `json:"name"`      // 容器服务名称
	NameSpace string `json:"namespace"` // 容器服务命名空间
	Image     string `json:"image"`     // 容器服务镜像
	Task      string `json:"task"`      // 容器服务任务名称
	NetworkNs string `json:"networkNs"` // 容器服务网络
	IP        string `json:"ip"`        // 容器服务 IP
	SnapShot  string `json:"snapshot"`  // 容器服务快照
	Schema    string `json:"schema"`    // 容器服务 Schema，用於數據庫設定
	User      string `json:"user"`      // 容器服务用户名
	Password  string `json:"password"`  // 容器服务密码
}

func (c *ContainerD) NewClient() *containerd.Client {
	return nil
}

// ContainerManager 容器服务管理員
type ContainerManager struct {
	configPath    string
	containerList map[string]*ContainerList
}

type ContainerList struct {
	networkLock sync.Locker
	cfg         ContainerD
	builder     Builder
	user        string
	status      int
}

// NewContainderManager 新建容器服务管理員
func NewContainderManager(path string) (*ContainerManager, error) {
	if strings.TrimSpace(path) == "" {
		path = defaultConfigPath
	}
	if err := checkDir(path); err != nil {
		log.Warn("check file config directory failed, %v", err)
		return nil, err
	}
	r := Load{prefix: "./example/", client: nil}
	configs, err := r.loadAllContainerD()
	if err != nil {
		log.Warn("load containerd config failed, %v", err)
		return nil, err
	}
	containerList := make(map[string]*ContainerList)
	for container, config := range configs {
		builder, err := NewBuilder(config)
		if err != nil {
			log.Warn("make containerd client failed, %v", err)
			return nil, err
		}
		containerList[container] = &ContainerList{
			user:        "",                   // 获取函数名称 register's function
			status:      containerdStatusInit, // 容器服务状态为占用 containerd is occupied
			cfg:         config,               // 容器服务配置 containerd config
			builder:     builder,              // 容器服务构建器 containerd builder
			networkLock: &sync.Mutex{},        // 容器服务网络锁 containerd network lock
		}
	}

	return &ContainerManager{configPath: path, containerList: containerList}, nil
}

// RegisterFunc 注册函式名称 register's function
// Register containerd service
type RegisterFunc func() string

// GetBuilder 获取容器服务构建器 GetBuilder is used to get containerd builder
func (cm *ContainerManager) GetBuilder(containerName string, regfunc RegisterFunc) (Builder, error) {
	// 如果没有配置，则返回错误
	if _, ok := cm.containerList[containerName]; !ok {
		return nil, errors.New("invalid config container name")
	}

	// 如果可以进行占用，则续续以下操作
	// if we can occupy, then continue to do the following operation.
	cm.containerList[containerName].networkLock.Lock()                // 加锁 lock
	cm.containerList[containerName].user = regfunc()                  // 获取函数名称 register's function
	cm.containerList[containerName].status = containerdStatusOccupied // 容器服务状态为占用 containerd is occupied

	// 正常返回容器服务构建器 return containerd builder
	return cm.containerList[containerName].builder, nil
}

// ReturnBuilder 适放容器服务构建器 ReturnBuilder is used to release containerd builder.
func (cm *ContainerManager) ReturnBuilder(containerName string) error {
	// 如果没有配置，则返回错误
	if _, ok := cm.containerList[containerName]; !ok {
		return errors.New("invalid config container name")
	}

	// 如果可以进行占用，则续续以下操作
	// if we can occupy, then continue to do the following operation.
	cm.containerList[containerName].networkLock.Unlock()              // 解锁 unlock
	cm.containerList[containerName].user = ""                         // 获取函数名称 register's function
	cm.containerList[containerName].status = containerdStatusReturned // 容器服务状态为被适放 containerd is released

	// 正常适放容器服务构建器 release containerd builder
	return nil
}

// SetContainerManagerStatus 设置容器服务状态 SetStatus is used to set containerd status.
func SetContainerManagerStatus(containerName string, status int) error {
	// 如果没有配置，则返回错误
	// if we can occupy, then continue to do the following operation.
	if _, ok := Manager.containerList[containerName]; !ok {
		return errors.New("invalid config container name")
	}

	// 如果可以进行设定状态，则续续以下操作
	// if we can set status, then continue to do the following operation.
	Manager.containerList[containerName].status = status // 设置容器服务状态 set containerd status

	// 正常返回 return
	return nil
}

// checkDir 检查目录是否存在 checkDir is used to check directory is existed.
func checkDir(path string) error {
	if strings.TrimSpace(path) == "" {
		return errors.New("invalid path")
	}
	stat, err := os.Stat(path)
	if err != nil {
		return err
	}

	if !stat.IsDir() {
		return errors.New("invalid path, should be a directory")
	}

	return nil
}
