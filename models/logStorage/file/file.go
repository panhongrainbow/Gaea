package file

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Client 为 单一输出日志档案
type Client struct {
	fileName string // 日记文档的名称
	path     string // 日志文档集中的目录路径
	// level    int        // 日志等级
	// skip     int        // 略过等级
	// runtime  bool       // 本机资讯
	file    *os.File // 日志文档
	errFile *os.File // 错误的日志文档
	// hostname string     // 本机名称
	// service  string     // 服务名称
	split sync.Once  // 文档切割
	mu    sync.Mutex // 上锁
}

// New 为产生 单一输出日志档案 的对象，可以各自进行设定
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

	// 决定等级
	/*level, ok := config["level"]
	if !ok {
		err = fmt.Errorf("init XFileLog failed, not found level")
		return nil, err
	}
	c.level = LevelFromStr(level)*/

	// 决定服务
	/*service, _ := config["service"]
	if len(service) > 0 {
		c.service = service
	}*/

	// 决定显示细节
	/*runtime, ok := config["runtime"]
	if !ok || runtime == "true" || runtime == "TRUE" {
		c.runtime = true
	} else {
		c.runtime = false
	}*/

	// 决定忽略等级
	/*skip, _ := config["skip"]
	if len(skip) > 0 {
		skipNum, err := strconv.Atoi(skip)
		if err == nil {
			c.skip = skipNum
		}
	}*/

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

	// 决定主机名称
	/*hostname, _ := os.Hostname()
	c.hostname = hostname*/

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

// 重复使用的函式，先复制在这里，之后再处理

func isDir(path string) (bool, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return stat.IsDir(), nil
}

const (
	SpliterDelay = 5
	CleanDays    = -3
)

// split the log file
func (c *Client) spliter() {
	preHour := time.Now().Hour()
	splitTime := time.Now().Format("2006010215")
	defer c.Close()
	for {
		time.Sleep(time.Second * SpliterDelay)
		if time.Now().Hour() != preHour {
			c.clean()
			c.rename(splitTime)
			preHour = time.Now().Hour()
			splitTime = time.Now().Format("2006010215")
		}
	}
}

func (c *Client) clean() (err error) {
	deadline := time.Now().AddDate(0, 0, CleanDays)
	var files []string
	files, err = filepath.Glob(fmt.Sprintf("%s/%s.log*", c.path, c.fileName))
	if err != nil {
		return
	}
	var fileInfo os.FileInfo
	for _, file := range files {
		if filepath.Base(file) == fmt.Sprintf("%s.log", c.fileName) {
			continue
		}
		if filepath.Base(file) == fmt.Sprintf("%s.log.wf", c.fileName) {
			continue
		}
		if fileInfo, err = os.Stat(file); err == nil {
			if fileInfo.ModTime().Before(deadline) {
				os.Remove(file)
			} else if fileInfo.Size() == 0 {
				os.Remove(file)
			}
		}
	}
	return
}

// openFile 为开档并写档程式，如果在进行单元测试时，会把日志导到通道，并写入到通道
func (c *Client) openFile(filename string) (*os.File, error) {
	// 以下保留原程式码
	file, err := os.OpenFile(filename,
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0644,
	)

	// 错误回传
	if err != nil {
		return nil, fmt.Errorf("open %s failed, err:%v\n", filename, err)
	}

	// 正确回传，不是 在单元单试的执行环境下
	return file, err
}

func delayClose(fp *os.File) {
	if fp == nil {
		return
	}
	time.Sleep(1000 * time.Millisecond)
	fp.Close()
}

func (c *Client) rename(shuffix string) (err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	defer c.ReOpen()
	if c.file == nil {
		return
	}
	var fileInfo os.FileInfo
	normalLog := c.path + "/" + c.fileName + ".log"
	warnLog := normalLog + ".wf"
	newLog := fmt.Sprintf("%s/%s.log-%s.log", c.path, c.fileName, shuffix)
	newWarnLog := fmt.Sprintf("%s/%s.log.wf-%s.log.wf", c.path, c.fileName, shuffix)
	if fileInfo, err = os.Stat(normalLog); err == nil && fileInfo.Size() == 0 {
		return
	}
	if _, err = os.Stat(newLog); err == nil {
		return
	}
	if err = os.Rename(normalLog, newLog); err != nil {
		return
	}
	if fileInfo, err = os.Stat(warnLog); err == nil && fileInfo.Size() == 0 {
		return
	}
	if _, err = os.Stat(newWarnLog); err == nil {
		return
	}
	if err = os.Rename(warnLog, newWarnLog); err != nil {
		return
	}
	return
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
