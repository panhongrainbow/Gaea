package file

import (
	"fmt"
	"os"
	"sync"
)

const (
	SpliterDelay = 5
	CleanDays    = -3
)

// Client 为 单一输出日志档案
type Client struct {
	fileName string     // 日记文档的名称
	path     string     // 日志文档集中的目录路径
	file     *os.File   // 日志文档
	errFile  *os.File   // 错误的日志文档
	split    sync.Once  // 文档切割
	mu       sync.Mutex // 上锁
}

// New 为产生一个新档案储存的介面
func New(config map[string]string) (*Client, error) {

	// 产生回传物件和错误回传变数
	var c Client
	var err error

	// 决定目录
	path, ok := config["path"]
	if !ok {
		err = fmt.Errorf("init XFileLog failed, not found path")
		return nil, err
	}

	// 决定文档
	fileName, ok := config["filename"]
	if !ok {
		err = fmt.Errorf("init XFileLog failed, not found filename")
		return nil, err
	}

	// 建立目录
	isDir, err := isDir(path)
	if err != nil || !isDir {
		err = os.MkdirAll(path, 0755)
		if err != nil {
			return nil, fmt.Errorf("Mkdir failed, err:%v\n", err)
		}
	}

	// 设定路径和档名
	c.path = path
	c.fileName = fileName

	// 决定分档
	body := func() {
		go c.spliter()
	}
	doSplit, ok := config["dosplit"]
	if !ok {
		doSplit = "true"
	}
	if doSplit == "true" {
		c.split.Do(body)
	}

	// 成功回传物件
	return &c, nil
}

// >>>>> >>>>> >>>>> >>>>> >>>>> 以下为公开的介面函式

// ReOpen 为开启 单一输出日志档案 的输出 (为最初代码的 ReOpen 函式)
func (c *Client) ReOpen() error {
	go delayClose(c.file)
	go delayClose(c.errFile)

	normalLog := c.path + "/" + c.fileName + ".log"
	file, err := c.openFile(normalLog)
	if err != nil {
		return err
	}

	c.file = file
	warnLog := normalLog + ".wf"
	c.errFile, err = c.openFile(warnLog)
	if err != nil {
		c.file.Close()
		c.file = nil
		return err
	}

	return nil
}

// Close 为关闭 单一输出日志档案 的输出
func (c *Client) Close() error {
	c.file.Close()
	c.errFile.Close()
	return nil
}

// Write 为最后的写入函式，会把日志写入档案
func (c *Client) Write(logByte []byte) error {
	// 正式写入日志
	_, err := c.file.Write(logByte) // 执行写入日志

	// 正确或错误回传
	return err
}

// WriteErr 为最后的写入函式，会把错误日志写入档案
func (c *Client) WriteErr(logByte []byte) error {
	// 正式写入日志
	_, err := c.errFile.Write(logByte) // 执行写入日志

	// 正确或错误回传
	return err
}
