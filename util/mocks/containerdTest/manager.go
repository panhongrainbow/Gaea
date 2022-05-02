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

// ContainderManager 容器服务管理員
type ContainderManager struct {
	configPath     string
	containderList map[string]ContainderList
}

type ContainderList struct {
	networkLock sync.Locker
	cfg         ContainerD
	builder     Builder
}

// GetBuilder 获取容器服务构建器
func (cm *ContainderManager) GetBuilder(name string) (Builder, error) {
	// 如果没有配置，则返回错误
	if _, ok := cm.containderList[name]; !ok {
		return nil, errors.New("invalid config name")
	}
	cm.containderList[name].networkLock.Lock()
	return cm.containderList[name].builder, nil
}

// ReturnBuilder 归还容器服务构建器
func (cm *ContainderManager) ReturnBuilder(name string) error {
	cm.containderList[name].networkLock.Unlock()
	return nil
}

// NewContainderManager 新建容器服务管理員
func NewContainderManager(path string) (*ContainderManager, error) {
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
	containerList := make(map[string]ContainderList)
	for container, config := range configs {
		builder, err := NewBuilder(config)
		if err != nil {
			log.Warn("make containerd client failed, %v", err)
			return nil, err
		}
		containerList[container] = ContainderList{
			cfg:         config,
			builder:     builder,
			networkLock: &sync.Mutex{},
		}
	}

	return &ContainderManager{configPath: path, containderList: containerList}, nil
}

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
