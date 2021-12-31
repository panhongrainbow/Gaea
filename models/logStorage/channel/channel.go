package channel

// Client 为 模拟用的日志通道
type Client struct {
	fileName    string             // 将要被模拟的档案数量
	chanSize    string             // 模拟用途的双向通道一次可以记录日志的笔数
	mockChannel *MockMultiXLogFile // 模拟用的通道
	// 模拟用途的双向通道并不用建立目录路径
}

// New 为产生 模拟用的双向通道 的对象，可以各自进行设定
func New(fileName, chanSize string) (*Client, error) {
	// 先组成设定值
	config := make(map[string]string, 2)
	config["filename"] = fileName // 文档输出设定，多档输出用逗号隔开 "log1,log2"
	config["chanSize"] = chanSize // 每一个模拟用的通道可以接收的笔日志资料，比如 5 笔就填入 "5"

	// 初始化模拟用的通道
	mf := new(MockMultiXLogFile) // 建立物件
	err := mf.Init(config)       // 初始化物件
	if err != nil {
		return nil, err // 错误回传
	}

	// 进行开启档案
	// (表面上是重新开启档案，但实际上是检查双向通道是否有全部被开启) // 先暂不使用 !
	/* err = mf.Open()
	if err != nil {
		return nil, err
	} */

	// 设定全域的模拟通道
	// SetGlobalMockMultiXLogFile(mf) // 先暂不使用 !

	// 正确回传
	return &Client{ // 回传客户端
			fileName:    fileName,
			chanSize:    chanSize,
			mockChannel: mf,
		},
		nil // 没错误发生
}

// Open 为开启 模拟用途的双向通道 的输出
func (c *Client) Open() error {
	return c.mockChannel.ReOpen() // 直接交由模拟通道物件执行开启
}

// Close 为关闭 模拟用途的双向通道 的输出
func (c *Client) Close() error {
	return c.mockChannel.Close() // 直接交由模拟通道物件执行开启
}
