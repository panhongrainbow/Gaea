package containerTest

import (
	"errors"
	"fmt"
	"github.com/XiaoMi/Gaea/log"
	"github.com/XiaoMi/Gaea/util/mocks/containerTest/builder"
	"github.com/XiaoMi/Gaea/util/mocks/containerTest/builder/containerd"
	"github.com/XiaoMi/Gaea/util/mocks/containerTest/builder/containerd/run"
	"os"
	"runtime"
	"strings"
	"sync"
)

import (
	"flag"
)

var (
	cmConfigFile = flag.String("cm", "./etc/containerd.ini", "containerd manager 配置")
	// DefaultConfigPath 預設的容器設定路徑
	DefaultConfigPath = "./etc/containerd/"
)

var Manager *ContainerManager

// ContainerManager 容器服务管理員
// ContainerManager is used to manage Containerd.
type ContainerManager struct {
	Enable        bool
	ConfigPath    string
	ContainerList map[string]*ContainerList
}

// ContainerList 容器服务列表
// ContainerList is used to list containerd clients.
type ContainerList struct {
	NetworkLock sync.Locker
	Cfg         containerd.ContainerD
	Builder     builder.Builder
	User        string
	Status      int
}

// SetContainerManagerStatus 设置容器服务状态 SetStatus is used to set containerd status.
func SetContainerManagerStatus(containerName string, status int) error {
	// 如果没有配置，则返回错误
	// if we can occupy, then continue to do the following operation.
	if _, ok := Manager.ContainerList[containerName]; !ok {
		return errors.New("invalid config container name")
	}

	// 如果可以进行设定状态，则续续以下操作
	// if we can set status, then continue to do the following operation.
	Manager.ContainerList[containerName].Status = status // 设置容器服务状态 set containerd status

	// 正常返回 return
	return nil
}

// NewContainderManager 新建容器服务管理員
func NewContainderManager(path string) (*ContainerManager, error) {
	if strings.TrimSpace(path) == "" {
		path = DefaultConfigPath
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
		builder, err := containerd.NewBuilder(config)
		if err != nil {
			log.Warn("make containerd client failed, %v", err)
			return nil, err
		}
		containerList[container] = &ContainerList{
			User:        "",                       // 获取函数名称 register's function
			Status:      run.ContainerdStatusInit, // 容器服务状态为占用 containerd is occupied
			Cfg:         config,                   // 容器服务配置 containerd config
			Builder:     builder,                  // 容器服务构建器 containerd builder
			NetworkLock: &sync.Mutex{},            // 容器服务网络锁 containerd network lock
		}
	}

	return &ContainerManager{ConfigPath: path, ContainerList: containerList}, nil
}

// registerFunc 注册函式名称 register's function
// Register containerd service
type registerFunc func() string

func defaultFunction() string {
	counter, _, _, success := runtime.Caller(2)
	if !success {
		return "unknown"
	}
	return runtime.FuncForPC(counter).Name()
}

// AppendCurrentFunction 返回当前函数名
// AppendCurrentFunction returns the current function name
// appendStr 是用来判别各个协程
// appendStr is to identify each goroutine's function name
func AppendCurrentFunction(layerNumber int, appendStr string) string {
	counter, _, _, success := runtime.Caller(layerNumber)
	if !success {
		return "unknown"
	}
	return runtime.FuncForPC(counter).Name() + appendStr
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

// GetBuilder 获取容器服务构建器 GetBuilder is used to get containerd builder
func (cm *ContainerManager) GetBuilder(containerName string, regFunc registerFunc) (builder.Builder, error) {
	if regFunc == nil {
		regFunc = defaultFunction
	}

	// 如果没有配置，则返回错误
	if _, ok := cm.ContainerList[containerName]; !ok {
		return nil, errors.New("invalid config container name")
	}

	// 如果可以进行占用，则续续以下操作
	// if we can occupy, then continue to do the following operation.
	cm.ContainerList[containerName].NetworkLock.Lock() // 加锁 lock
	cm.ContainerList[containerName].User = regFunc()   // 获取函数名称 register's function
	fmt.Println("user get >>> ", cm.ContainerList[containerName].User)
	// log.SetGlobalLogger()
	log.Notice("user get >>> ", cm.ContainerList[containerName].User)
	cm.ContainerList[containerName].Status = run.ContainerdStatusOccupied // 容器服务状态为占用 containerd is occupied

	// 正常返回容器服务构建器
	// return containerd builder
	return cm.ContainerList[containerName].Builder, nil
}

// ReturnBuilder 适放容器服务构建器 ReturnBuilder is used to release containerd builder.
func (cm *ContainerManager) ReturnBuilder(containerName string, regFunc registerFunc) error {
	if regFunc == nil {
		regFunc = defaultFunction
	}
	fmt.Println("user return >>> ", regFunc())
	log.Notice("user return >>> ", cm.ContainerList[containerName].User)

	// 如果没有配置，则返回错误
	if _, ok := cm.ContainerList[containerName]; !ok {
		return errors.New("invalid config container name")
	}

	// 如果可以进行占用，则续续以下操作
	// if we can occupy, then continue to do the following operation.
	cm.ContainerList[containerName].NetworkLock.Unlock()                  // 解锁 unlock
	cm.ContainerList[containerName].User = ""                             // 获取函数名称 register's function
	cm.ContainerList[containerName].Status = run.ContainerdStatusReturned // 容器服务状态为被适放 containerd is released

	// 正常适放容器服务构建器 release containerd builder
	return nil
}
