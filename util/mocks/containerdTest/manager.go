package containerdTest

import (
	"flag"
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

// ContainderdManager 容器服务管理員
type ContainderdManager struct {
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
