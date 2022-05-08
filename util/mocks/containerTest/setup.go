package containerTest

// init 初始化 containerd 容器管理员服务 init is init function of containerd manager
func init() {
	// 先开始辨判连接方式
	if err := check(); err == nil {
		if err := setup(); err != nil {
			panic(err)
		}
	}
}

// Setup 连线测试的到取得资料库版本就为正确
func setup() error {
	/*var err error
	// 连接到容器管理器 connect to the containerd manager.
	Manager, err = NewContainderManager("./example/")
	return err*/
	return nil
}

// check 检查配置文件和环境是否正确 check is a function to check the config file and test environment are ok
func check() error {
	return nil
}
