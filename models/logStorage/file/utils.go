package file

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// >>>>> >>>>> >>>>> >>>>> >>>>> 以下為未公開的輔助函式

// isDir 为确认是否为目录
func isDir(path string) (bool, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return stat.IsDir(), nil
}

// split 为用来切割日志文档
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

// clean 為用來移除舊的日誌文檔
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

// delayClose 为延迟关闭日志文档
func delayClose(fp *os.File) {
	if fp == nil {
		return
	}
	time.Sleep(1000 * time.Millisecond)
	fp.Close()
}

// rename 为把日志文档重新改名
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
