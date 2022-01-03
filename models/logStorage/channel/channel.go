package channel

// Client 为 模拟用的日志通道
type Client struct {
	fileName    string             // 将要被模拟的档案数量
	chanSize    string             // 模拟用途的双向通道一次可以记录日志的笔数
	mockChannel *MockMultiXLogFile // 模拟用的通道
	// 模拟用途的双向通道并不用建立目录路径
}

// New 为产生 模拟用的双向通道 的对象，可以各自进行设定
func New(config map[string]string) (*Client, error) {
	// 先组成设定值
	// config := make(map[string]string, 2)
	// config["filename"] = fileName // 文档输出设定，多档输出用逗号隔开 "log1,log2"
	// config["chanSize"] = chanSize // 每一个模拟用的通道可以接收的笔日志资料，比如 5 笔就填入 "5"

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
			fileName:    config["filename"],
			chanSize:    config["chanSize"],
			mockChannel: mf,
		},
		nil // 没错误发生
}

// ReOpen 为开启 模拟用途的双向通道 的输出
func (c *Client) ReOpen() error {
	return c.mockChannel.ReOpen() // 直接交由模拟通道物件执行开启
}

// Close 为关闭 模拟用途的双向通道 的输出
func (c *Client) Close() error {
	return c.mockChannel.Close() // 直接交由模拟通道物件执行开启
}

// Write 为最后的写入函式，会把日志写入通道里
func (c *Client) Write(logByte []byte) error {
	// 正式写入日志
	// _, err := c.mockChannel.mockFile[c.fileName].Write(logByte) // 执行写入日志 (写法一)
	return c.mockChannel.Write(c.fileName, logByte)

	// 正确或错误回传
	// return err
}

// WriteErr 为最后的写入函式，会把错误日志写入通道里
func (c *Client) WriteErr(logByte []byte) error {
	// 正式写入日志
	// _, err := c.mockChannel.mockFile[c.fileName].Write(logByte) // 执行写入日志 (写法一)
	return c.mockChannel.WriteErr(c.fileName, logByte)

	// 正确或错误回传
	// return err
}
