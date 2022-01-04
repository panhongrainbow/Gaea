package xlog

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

const (
	defaultSize = 10 // 模拟文件的双向通道的预设尺寸
)

// >>>>> >>>>> >>>>> >>>>> >>>>> xWrite 抽换输出位罝的介面
// xWrite 为 XMultiFileLog，XFileLog 建立单元测试前的约定，

// 设定全域模拟档案的双向通道
var xWrite XWrite

// XWrite 功能是要处理写入介面的抽换，这和 XLogger 介面重复
// 就是把会写入 XMultiFileLog 资料整个会写入档案的介面抽掉，用双向通道替换
type XWrite interface {
	// Init 以后的函式为 XMultiFileLog，XFileLog 和 XWrite 共用的部份
	Init(config map[string]string) error // 进行初始化工作，部份取代 XLogger 的功能
	Open() error                         // 重新开启日志输出，全部取代 XLogger 的功能
	Close()                              // 关闭日志输出，全部取代 XLogger 的功能
	// GetChan 为 XMultiFileLog，XFileLog 和 XWrite 三者要新增的函式，目的是要把日志传到应要传送的双向通道
	GetChan(string) (chan string, error) // XMultiFileLog，XFileLog 和 XWrite 三者都缺这个介面，XMultiFileLog 要补
}

// getChan 为所有的 GetChan 函式都会指到这一只 getChan 函式去执行
func getChan(fileName string) (chan string, error) {
	// 如果对应到的资料类型为 MockMultiXLogFile，就进行相对应的处理
	if mock, ok := xWrite.(*MockMultiXLogFile); ok == true {
		// 正确回传
		return mock.mockFile[fileName], nil
	}

	// 错误回传
	return nil, fmt.Errorf("not found xwrite channel [%s]", fileName)
}

// PrintMockChanMsg 为印出所有双向通道的讯息
func (mf *MockMultiXLogFile) PrintMockChanMsg() (ret []string) {
	// 把所有的双向通道一个个列出来，并取出在双向通道内的讯息
	for fileName, mockChan := range mf.mockFile {
	LOOP:
		for { // 无限回圈取出讯息
			select {
			case msg := <-mockChan:
				ret = append(ret, fileName+"::"+msg)
			default:
				break LOOP
			}
		}
	}

	// 排序，这样单元测试的答案才会统一
	sort.Strings(ret)

	// 正确回传
	return
}

// >>>>> >>>>> >>>>> >>>>> >>>>> MockMultiXLogFile 的功能为把原本要写入档案的日志移转到双向通道

// MockMultiXLogFile 使用通道 Channel 去取代档案 file 去进行单元测试
type MockMultiXLogFile struct {
	mockFile map[string]chan string // 分别用一个字串对应到一个双向通道，用双向通道取代档案
	fileName []string               // 模拟档案名称的阵列
	chanSize int                    // 每一份双向通道的尺寸大小
}

// Init 初始化模拟档案的双向通道
func (mf *MockMultiXLogFile) Init(config map[string]string) error {
	// 取得要被模拟档案名称的 fileName 阵列
	var fileName []string
	fStr, ok := config["filename"] // 先确认 fileName 设定值是否存在
	if ok {                        // 如果 fileName 值 存在
		fileName = strings.Split(fStr, ",") // 以逗点分隔开来
	}
	if !ok { // 如果 fileName 值 不 存在
		err := fmt.Errorf("init XFileLog failed, not found filename")
		return err
	}
	mf.fileName = fileName // 把 fileName 阵列存放在物件中

	// 先取得双向通通道的尺寸 chanSize 设定值
	size, ok := config["chan_size"] // 先确认 chanSize 设定值是否存在
	if ok == true {                 // 如果 chanSize 设定值 存在
		if len(size) > 0 { // 如果 chanSize 值 存在
			sizeNum, err := strconv.Atoi(size)
			if err == nil { // 如果设定值没有错误，就使用设定值
				mf.chanSize = sizeNum
			}
			if err != nil { // 如果设定值有错误，就使用预设值
				mf.chanSize = defaultSize
			}
		}
	}

	// 初始化建立双向通道
	mf.mockFile = make(map[string]chan string, len(fileName))

	// 正确回传
	return nil
}

// Open 在档案中的 ReOpen 函式会真的先关档再重新开档，但在使用双向通道去模拟档案时，只检查那些双向通道没有开启，没开启就立刻开启
func (mf *MockMultiXLogFile) Open() error {
	// 对每一个档名进行检查
	for i := 0; i < len(mf.fileName); i++ {
		mockChan, ok := mf.mockFile[mf.fileName[i]]
		if ok == false || mockChan == nil { // 只要任一键值或者是双向通道不存在，就立刻建立双向通道
			mf.mockFile[mf.fileName[i]] = make(chan string, mf.chanSize) // 建立双向通道
		}
	}

	// 正确回传
	return nil
}

// Close 为关闭所有双向通道
func (mf *MockMultiXLogFile) Close() {
	// 对每一个档名进行关闭
	for i := 0; i < len(mf.fileName); i++ {
		mockChan, _ := mf.mockFile[mf.fileName[i]]
		if mockChan != nil { // 只要双向通道存在，就立刻关闭双向通道
			close(mf.mockFile[mf.fileName[i]]) // 关闭双向通道
		}
	}
}

// SetGlobalMockMultiXLogFile 如果全域模拟档案的双向通道存在就关闭，重新在赋值
func SetGlobalMockMultiXLogFile(xw XWrite) {
	// 把双向通道，全部关闭
	if xWrite != nil {
		xWrite.Close()
	}

	// 重新赋值
	xWrite = xw
}

// CloseGlobalMockMultiXLogFile 关闭所有的全域模拟档案的双向通道
func CloseGlobalMockMultiXLogFile() {
	// 把双向通道，全部关闭
	if xWrite != nil {
		xWrite.Close()
	}

	// 重新赋空值
	xWrite = nil
}

// GetChan 为 XMultiFileLog，XFileLog 和 XWrite 三者都要新增的函式，目的是要把日志传到应要传送的双向通道
func (mf *MockMultiXLogFile) GetChan(fileName string) (chan string, error) {
	// 由 getChan 函式接手进行后续处理
	return getChan(fileName)
}

// >>>>> >>>>> >>>>> >>>>> >>>>> 测试环境的准备

// testMultiAndFileSuite 建立测试环境 (多档分流的测试环境也可以用在单档)
func testMultiAndFileSuite(config map[string]string) error {
	// 初始化模拟文件
	mf := new(MockMultiXLogFile)
	err := mf.Init(config) // 在这里进行初始化模拟文件
	if err != nil {
		return err // 错误回传
	}

	// 进行重新开启档案
	// 表面上是重新开启档案，但实际上是检查双向通道是否有全部被开启
	err = mf.Open()
	if err != nil {
		return err // 错误回传
	}

	// 设定全域的双向通道
	SetGlobalMockMultiXLogFile(mf)

	// 正确回传
	return nil
}
