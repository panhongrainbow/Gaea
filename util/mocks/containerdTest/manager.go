package containerdTest

import (
	"flag"
	"github.com/containerd/containerd"
	"io/ioutil"
)

var (
	cmConfigFile = flag.String("cm", "./etc/containerd.ini", "containerd manager 配置")
	// 預設的容器設定路徑
	defaultConfigPath = "./etc/containerd/"
)

// init 初始化 containerd 容器管理员服务
func init() {
	/*flag.Parse()
	if err := NewContainderdManager(*cmConfigFile); err != nil {
		log.Error("init containerd manager failed, %v", err)
		panic(err)
	}*/
}

// 先進行 Parse 操作

// ContainerD 为 Containerd 容器服务所要 Parse 对应的内容
type ContainerD struct {
	Type      string `json:"type"`      // 容器服务类型
	Name      string `json:"name"`      // 容器服务名称
	Image     string `json:"image"`     // 容器服务镜像
	Container string `json:"container"` // 容器服务容器名称
	Task      string `json:"task"`      // 容器服务任务名称
	IP        string `json:"ip"`        // 容器服务 IP
	Schema    string `json:"schema"`    // 容器服务 Schema，用於數據庫設定
	User      string `json:"user"`      // 容器服务用户名
	Password  string `json:"password"`  // 容器服务密码
}

func (c *ContainerD) NewClient() *containerd.Client {
	return nil
}

// ContainderdManager 容器服务管理員
type ContainderdManager struct {
	run        *Run
	task       chan *ContainerdTask
	configPath string
}

// ParseCMconfig 載入容器服务管理员設定
func ParseCMconfig(configPath string) ([]byte, error) {
	return ioutil.ReadFile(configPath)
}

// NewContainderdManager 新建容器服务管理員
func NewContainderdManager(path string) (*ContainderdManager, error) {
	/*if strings.TrimSpace(path) == "" {
		path = defaultConfigPath
	}
	if err := checkDir(path); err != nil {
		log.Warn("check file config directory failed, %v", err)
		return nil, err
	}
	return &ContainderdManager{configPath: path}, nil*/
	return nil, nil
}
