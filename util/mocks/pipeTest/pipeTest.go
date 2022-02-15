package pipeTest

import (
	"bufio"
	"github.com/stretchr/testify/require"
	"net"
	"sync"
	"testing"
)

// ReplyFuncType 回应函式的型态
type ReplyFuncType func([]uint8) []uint8

// TestReplyFunc 在这里会处理常接收到什么讯息，要将下来跟着回应什么讯息
// 目前此函式只是在测试验证流程，回应讯息为接收讯息加 1
//     比如 当接收值为 1，就会回传值为 2 给对方
//     比如 当接收值为 2，就会回传值为 3 给对方
func TestReplyFunc(data []uint8) []uint8 {
	return []uint8{data[0] + 1} // 回应讯息为接收讯息加 1
}

// DcMocker 用来模拟数据库服务器的读取和回应的物件
type DcMocker struct {
	t         *testing.T      // 单元测试的物件
	bufReader *bufio.Reader   // 服务器的读取
	bufWriter *bufio.Writer   // 服务器的回应
	connRead  net.Conn        // pipe 的读取连线
	connWrite net.Conn        // pipe 的写入连线
	wg        *sync.WaitGroup // 流程的操作边界
	replyFunc ReplyFuncType   // 设定相对应的回应函式
	err       error           // 错误
}

// NewDcServerClient 产生直连 DC 模拟双方，包含客户端和服务端
func NewDcServerClient(t *testing.T, reply ReplyFuncType) (mockClient *DcMocker, mockServer *DcMocker) {
	// 先产生两组 Pipe
	read0, write0 := net.Pipe() // 第一组 Pipe
	read1, write1 := net.Pipe() // 第二组 Pipe

	// 产生客户端和服务端双方，分别为 mockClient 和 mockServer
	mockClient = NewDcMocker(t, read0, write1, reply)
	mockServer = NewDcMocker(t, read1, write0, reply)

	// 结束
	return
}

// NewDcMocker 产生新的 dc 模拟物件
func NewDcMocker(t *testing.T, connRead, connWrite net.Conn, reply ReplyFuncType) *DcMocker {
	return &DcMocker{
		t:         t,                          // 单元测试的物件
		bufReader: bufio.NewReader(connRead),  // 服务器的读取 (实现缓存)
		bufWriter: bufio.NewWriter(connWrite), // 服务器的回应 (实现缓存)
		connRead:  connRead,                   // pipe 的读取连线
		connWrite: connWrite,                  // pipe 的写入连线
		wg:        &sync.WaitGroup{},          // 流程的操作边界
		replyFunc: reply,                      // 回应函式
	}
}

// GetConnRead 为获得直连 dc 模拟物件的读取连线
func (dcM *DcMocker) GetConnRead() net.Conn {
	return dcM.connRead
}

// GetConnWrite 为获得直连 dc 模拟物件的写入连线
func (dcM *DcMocker) GetConnWrite() net.Conn {
	return dcM.connWrite
}

// GetBufReader 为获得直连 dc 模拟物件的读取缓存
func (dcM *DcMocker) GetBufReader() *bufio.Reader {
	return dcM.bufReader
}

// GetBufWriter 为获得直连 dc 模拟物件的写入缓存
func (dcM *DcMocker) GetBufWriter() *bufio.Writer {
	return dcM.bufWriter
}

// PutConnRead 为临时修改直连 dc 模拟物件的读取连线
func (dcM *DcMocker) PutConnRead(connRead net.Conn, bufReader *bufio.Reader) error {
	// 先进行修改
	if connRead != nil {
		dcM.connRead = connRead // 修改连线
	}
	if bufReader != nil {
		dcM.bufReader = bufReader // 修改缓存
	}

	// 正确回传
	return nil
}

// PutConnWrite 为临时修改直连 dc 模拟物件的写入连线
func (dcM *DcMocker) PutConnWrite(connWrite net.Conn, bufWriter *bufio.Writer) error {
	// 先进行修改
	if connWrite != nil {
		dcM.connWrite = connWrite // 修改连线
	}
	if bufWriter != nil {
		dcM.bufWriter = bufWriter // 修改缓存
	}

	// 正确回传
	return nil
}

// ResetDcMockers 为重置单一连线方向的直连 dc 模拟物件
func (dcM *DcMocker) ResetDcMockers(otherSide *DcMocker) error {
	// 重新建立全新两组 Pipe
	newRead, newWrite := net.Pipe() // 第一组 Pipe

	// 单方向的状况为 dcM 写入 Pipe，otherSide 读取 Pipe

	// 先重置 发送讯息的那一方 部份
	dcM.bufWriter = bufio.NewWriter(newWrite) // 服务器的回应 (实现缓存)
	dcM.connWrite = newWrite                  // pipe 的写入连线
	dcM.wg = &sync.WaitGroup{}                // 流程的操作边界

	// 先重置 mockServer 部份
	otherSide.bufReader = bufio.NewReader(newRead) // 服务器的读取 (实现缓存)
	otherSide.connRead = newRead                   // pipe 的读取连线

	// 正常回传
	return nil
}

// SendOrReceive 为直连 dc 用来模拟接收或传入讯息
func (dcM *DcMocker) SendOrReceive(data []uint8) *DcMocker {
	// dc 模拟开始
	dcM.wg.Add(1) // 只要等待直到确认资料有写入 pipe

	// 在这里执行 1传送讯息 或者是 2接收讯息
	go func() {
		// 执行写入工作
		_, err := dcM.bufWriter.Write(data) // 写入资料到 pipe
		err = dcM.bufWriter.Flush()         // 把缓存资料写进 pipe
		require.Equal(dcM.t, err, nil)
		err = dcM.connWrite.Close() // 资料写入完成，终结连线
		require.Equal(dcM.t, err, nil)

		// 写入工作完成
		dcM.wg.Done()
	}()

	// 重复使用物件
	return dcM
}

// Reply 为直连 dc 用来模拟 dc 回应数据
func (dcM *DcMocker) Reply(otherSide *DcMocker) (msg []uint8) {
	// 读取传送过来的讯息
	b, _, err := otherSide.bufReader.ReadLine() // 由另一方接收传来的讯息
	require.Equal(dcM.t, err, nil)

	// 等待和确认资料已经写入 pipe
	dcM.wg.Wait()

	// 重置模拟物件
	err = dcM.ResetDcMockers(otherSide)
	require.Equal(dcM.t, err, nil)

	// 回传回应讯息
	if otherSide.replyFunc != nil {
		msg = otherSide.replyFunc(b)
	}

	// 结束
	return
}

// WaitAndReset 为直连 dc 用来等待在 Pipe 的整个数据读写操作完成
func (dcM *DcMocker) WaitAndReset(otherSide *DcMocker) error {
	// 先等待整个数据读写操作完成
	dcM.wg.Wait()

	// 单方向完成 Pipe 的连线重置
	err := dcM.ResetDcMockers(otherSide)
	require.Equal(dcM.t, err, nil)

	// 正确回传
	return nil
}
