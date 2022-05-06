package mgr

import (
	"flag"
)

var (
	cmConfigFile = flag.String("cm", "./etc/containerd.ini", "containerd manager 配置")
	// 預設的容器設定路徑
	defaultConfigPath = "./etc/containerd/"
)

var Manager *ContainerManager

// init 初始化 containerd 容器管理员服务 init is init function of containerd manager
func init() {
	// 先开始辨判连接方式
	if err := check(); err == nil {
		if err := setup(); err != nil {
			panic(err)
		}
	}
}

/*flag.Parse()
if err := NewContainderdManager(*cmConfigFile); err != nil {
	log.Error("init containerd manager failed, %v", err)
	panic(err)
}*/

// setup 连线测试的到取得资料库版本就为正确
func setup() error {
	var err error
	// 连接到容器管理器 connect to the containerd manager.
	Manager, err = NewContainderManager("./example/")
	return err
}

/*func setup() error {
	typ := "tcp"
	if strings.Contains("192.168.122.2:3309", "/") {
		typ = "unix"
	}

	netConn, err := net.Dial(typ, "192.168.122.2:3309")
	if err != nil {
		return err
	}

	// 先随意测试
	test := make([]byte, 20)
	netConn.Read(test)
	fmt.Println(test)

	_ = netConn.Close()

	return nil
}*/

// check 检查配置文件和环境是否正确 check is a function to check the config file and test environment are ok
func check() error {
	return nil
}
