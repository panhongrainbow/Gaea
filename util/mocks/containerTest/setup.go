package containerTest

const (
	// DefaultConfigPath 預設的容器設定路径
	// DefaultConfigPath is the default config path
	DefaultConfigPath = "util/mocks/containerTest/example"

	// 启用容器管理员服务
	// enable the containerd manager
	enableManager = true
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

// setup 初始化容器管理员服务
// setup is init function of containerd manager
func setup() error {

	// 如果容器管理员服务已经启用，则直接返回
	// if the containerd manager is enabled, then return directly.
	if Manager != nil {
		return nil
	}

	// 初始化容器管理员服务
	// init the containerd manager
	err := initContainerTestXLog()
	if err != nil {
		return err
	}

	// test, _ := ParseContainerConfigFromFile("/home/panhong/go/src/github.com/panhongrainbow/Gaea/etc/containerTest.ini")
	// fmt.Println(test.ContainerTestEnable)

	// 获得容器管理员服务配置
	// get the containerd manager config
	/*path, err := os.Getwd()
	if err != nil {
		return err
	}*/

	// 决定容器管理员服务配置路径
	// decide the containerd manager config path
	/*absolutePath := ""
	if strings.Contains(path, "Gaea") {
		absolutePath = filepath.Join(strings.Split(path, "Gaea")[0], "Gaea", DefaultConfigPath)
	} else {
		return errors.New("invalid config path")
	}*/

	// 决定容器管理员服务配置路径
	// decide the containerd manager config path
	absPath, err := absolutePath(DefaultConfigPath)
	if err != nil {
		return err
	}

	// 连接到容器管理器，设定档在 Gaea (/home/panhong/go/src/github.com/panhongrainbow/Gaea/) 下的相对路径下 util/mocks/containerTest/example
	// connect to the containerd manager. config file is in Gaea directory, relative path is util/mocks/containerTest/example.
	Manager, err = NewContainderManager(absPath)

	// 启用容器管理员服务
	// enable the containerd manager
	Manager.Enable = enableManager

	// 返回错误
	// return the error
	return err
}

// check 检查配置文件和环境是否正确
// check is a function to check the config file and test environment are ok
func check() error {
	return nil
}
