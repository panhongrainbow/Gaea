package containerdTest

// init 初始化 Conainerd 容器服务
func init() {
	// 先开始辨判连接方式
	// err := setup()
}

// Client 容器對象，在這裡要可以直接操作容器
type Client interface {
	Pull(image string) error
}
