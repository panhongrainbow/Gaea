package containerdTest

import (
	"errors"
	"flag"
	"github.com/XiaoMi/Gaea/log"
	"github.com/containerd/containerd"
	"io/ioutil"
	"os"
	"strings"
	"sync"
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

/*
cfg := containerdTest.ContainerD{
		Sock:      "",
		Type:      "mariadb",
		Name:      "mariadb-server",
		NameSpace: "mariadb",
		Image:     "docker.io/panhongrainbow/mariadb:testing",
		Task:      "mariadb-server",
		NetworkNs: "/var/run/netns/gaea-mariadb",
		IP:        "10.10.10.10:3306",
		SnapShot:  "mariadb-server-snapshot",
		Schema:    "",
		User:      "xiaomi",
		Password:  "12345",
	}
*/

func (c *ContainerD) NewClient() *containerd.Client {
	return nil
}

// ContainderManager 容器服务管理員
type ContainderManager struct {
	networkLock map[string]sync.Mutex
	cfg         map[string]ContainerD
	Builder     map[string]Builder
	configPath  string
}

// ParseCMconfig 載入容器服务管理员設定
func ParseCMconfig(configPath string) ([]byte, error) {
	return ioutil.ReadFile(configPath)
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
	builder := make(map[string]Builder)
	for container, config := range configs {
		if builder[container], err = NewBuilder(config); err != nil {
			log.Warn("make containerd client failed, %v", err)
			return nil, err
		}
	}

	return &ContainderManager{configPath: path, cfg: configs, Builder: builder}, nil
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
