package containerTest

import (
	"errors"
)

const (
	// DefaultConfigPath 預設的容器設定路径
	// DefaultConfigPath is the default config path
	DefaultConfigPath = "util/mocks/containerTest/example"

	// 启用容器管理员服务
	// enable the containerd manager
	enableManager = false
)

// init 初始化 containerd 容器管理员服务
// init is init function of containerd manager
func init() {
	// 先开始辨判连接方式
	// check the configuration.
	if err := check(); err == nil {
		// 初始化容器管理员服务
		// init the containerd manager
		if err := setup(); err != nil {
			// 立即中断程序
			// immediately exit the program
			panic(err)
		}
	}
}

// check 检查配置文件和环境是否正确
// check is a function to check the config file and test environment are ok.
func check() error {
	// 如果容器管理员服务已经载入，则直接返回
	// if the containerd manager is loaded, then return directly.
	if manager != nil {
		return nil
	}

	// 根据设定档，确认容器管理员服务是否要启用
	// according to the config file, decide whether the containerd manager is enabled.
	iniCfg, err := ParseContainerConfigFromFile("util/mocks/containerTest/etc/containerTest.ini")
	if err != nil {
		return err
	}

	// 启用容器管理员服务
	// enable the containerd manager
	if iniCfg.ContainerTestEnable != "true" {
		manager.Enable = enableManager
		return errors.New("containerd manager is not enabled")
	}

	// 如果容器管理员服务已经启用，则直接返回
	// if the containerd manager is enabled, then return directly.
	return nil
}

// setup 初始化容器管理员服务
// setup is init function of containerd manager
func setup() error {
	// 决定容器管理员服务配置路径
	// decide the containerd manager config path
	absPath, err := absolutePath(DefaultConfigPath)
	if err != nil {
		return err
	}

	// 连接到容器管理器，设定档在 Gaea (/home/panhong/go/src/github.com/panhongrainbow/Gaea/) 下的相对路径下 util/mocks/containerTest/example
	// connect to the containerd manager. config file is in Gaea directory, relative path is util/mocks/containerTest/example.
	manager, err = NewContainderManager(absPath)
	if err != nil {
		return err
	}

	// 启用容器管理员服务
	// enable the containerd manager
	manager.Enable = true

	// 初始化容器管理员服务
	// init the containerd manager
	if err = initContainerTestXLog(); err != nil {
		return err
	}

	// 返回错误
	// return the error
	return err
}
