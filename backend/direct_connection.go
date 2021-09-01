// Copyright 2019 The Gaea Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package backend

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/fnv"
	"net"
	"strings"

	sqlerr "github.com/XiaoMi/Gaea/core/errors"
	"github.com/XiaoMi/Gaea/log"
	"github.com/XiaoMi/Gaea/mysql"
	"github.com/XiaoMi/Gaea/util/sync2"
)

// MockDcClient 🧚 单元测试数据库直连客户端
type MockDcClient struct {
	MockKey uint32 // 识别要 Mock 资料的关键 Key 值
	// 单元测试相关设定值
	// TakeOver bool // 现在是否由单元测试接管 (TakeOver 变数移到全域变数)
	// 单元测试资料回应
	// MockResult map[uint32]mysql.Result // (MockResult 变数移到全域变数)
}

// MarkTakeOver 函式 🧚 为 MockDcClient 资料执行单元测试数据库直连的标记函式 (TakeOver 变数移到全域变数后废除)
/*func (m *MockDcClient) MarkTakeOver() {
	m.TakeOver = true // 单元测试之后可以直接进行接管
}*/

// MarkTakeOver 函式 🧚 为 MockDcClient 资料执行单元测试数据库直连的标记函式 (设定)
func MarkTakeOver() {
	TakeOver = true // 单元测试之后可以直接进行接管
}

// IsTakeOver 函式 🧚 为 MockDcClient 资料执行单元测试数据库直连的确认函式 (TakeOver 变数移到全域变数后废除)
/*func (m *MockDcClient) IsTakeOver() bool {
	// 因为不是每个函式或过程会完整初始化 Mock Client 变数，如果没有这一层保护，防止 nil 指标的错误
	if m == nil {
		return false // 回传 false ，之后单元测试不允许进行介入程式内部的运作
	}
	return m.TakeOver // 只要是回传 true ，之后单元测试就会接管整个程式
}*/

// IsTakeOver 函式 🧚 为 MockDcClient 资料执行单元测试数据库直连的确认函式 (设定)
func IsTakeOver() bool {
	return TakeOver // 只要是回传 true ，之后单元测试就会接管整个程式，如果回传 false 则反之
}

// UnmarkTakeOver 函式 🧚 为 MockDcClient 资料执行单元测试数据库直连的反标记函式 (TakeOver 变数移到全域变数)
/*func (m *MockDcClient) UnmarkTakeOver() {
	m.TakeOver = false // 解除单元测试的接管状态 (移到全域变数)
}*/

// UnmarkTakeOver 函式 🧚 为 MockDcClient 资料执行单元测试数据库直连的反标记函式 (设定)
func UnmarkTakeOver() {
	TakeOver = false // 解除单元测试的接管状态
}

// MakeMockResult 函式 🧚 为 在单元测试数据库时建立直连回应资料的对应 (回应)
// 目前准备做法是 1设定 环境 2数据库名称 3SQL 指令 三个值的组合对应到 一个数据库资料回传
func (m *MockDcClient) MakeMockResult(addr, sql string, res mysql.Result) uint32 {
	// 把数据库网路位置和SQL字串转成单纯的数字
	h := fnv.New32a()
	h.Write([]byte(addr + ";" + sql + ";")) // 所有的字串后面都要加上分号

	// 直接预先写好数据库资料回传
	FakeDBInstance.MockResult[h.Sum32()] = res // 转成数值，运算速度较快
	return h.Sum32()                           // 回传登记的数值
}

// DirectConnection means connection to backend mysql
// 此资料被单元测试函式包围
type DirectConnection struct {
	conn *mysql.Conn

	addr     string
	user     string
	password string
	db       string

	capability uint32

	sessionVariables *mysql.SessionVariables

	status uint16

	collation mysql.CollationID
	charset   string
	salt      []byte

	defaultCollation mysql.CollationID
	defaultCharset   string

	pkgErr error
	closed sync2.AtomicBool

	// 🧚 扩增一些单元测试的属性
	MockDC *MockDcClient // 单元测试数据库直连客户端
	Trans  Transferred   // 单元测试的定义接口
}

// MarkTakeOver 函式 🧚 为 DirectConnection 资料执行单元测试数据库直连的标记函式 (TakeOver 变数移到全域变数后废除)
/*func (dc *DirectConnection) MarkTakeOver() {
	dc.MockDC.MarkTakeOver() // 操作底层函式
}*/

// IsTakeOver 函式 🧚 为 DirectConnection 资料执行单元测试数据库直连的确认函式 (TakeOver 变数移到全域变数后废除)
/*func (dc *DirectConnection) IsTakeOver() bool {
	return dc.MockDC.IsTakeOver() // 操作底层函式
}*/

// UnmarkTakeOver 函式 🧚 为 DirectConnection 资料执行单元测试数据库直连的反标记函式 (TakeOver 变数移到全域变数后废除)
/*func (dc *DirectConnection) UnmarkTakeOver() {
	dc.MockDC.UnmarkTakeOver() // 操作底层函式
}*/

// NewDirectConnection return direct and authorised connection to mysql with real net connection
func NewDirectConnection(addr string, user string, password string, db string, charset string, collationID mysql.CollationID) (*DirectConnection, error) {
	dc := &DirectConnection{
		addr:             addr,
		user:             user,
		password:         password,
		db:               db,
		charset:          charset,
		collation:        collationID,
		defaultCharset:   charset,
		defaultCollation: collationID,
		closed:           sync2.NewAtomicBool(false),
		sessionVariables: mysql.NewSessionVariables(),
	}

	// 🧚 指定要载入测试的方法
	if IsTakeOver() {
		// >>>>> >>>>> >>>>> >>>>> >>>>> 先决定要使用假资料的方法

		// 将来要抽换制造假资料的方法，就直接在这里抽换就好，这是唯一要修改的地方
		dc.Trans = new(basicLoad) // 目前是使用最简单的测试资料载入方法，做测试用

		// >>>>> >>>>> >>>>> >>>>> >>>>> 开始载入资料
		if dc.Trans.IsLoaded() == false { // 如果之前没载入测试资料
			if err := dc.Trans.LoadData(); err == nil {
				dc.Trans.MarkLoaded() // 标记单元测试资料载入成功
			} else {
				// 做成对应到 网路位置、帐号、密码等相关资料，会回传 SQL 的执行结果
				// 如果在执行单元测试过程中，没有 1 命中单元测试的测试资料 或者是 2 测试资料载入失败，就使用 Fatal 中止
				// Fatal 中止 只有在单元测试的环境下才会执行，不会影响到主程式，还算安全
				log.Fatal("单元测试载入测试资料失败 %s", err.Error()) // 直接中断
			}
		}

	}

	err := dc.connect()
	return dc, err
}

// connect means real connection to backend mysql after authorization
func (dc *DirectConnection) connect() error {
	// 🧚 直接由单元测试接管
	// if dc.MockDC.IsTakeOver() { // TakeOver 变数移到全域变数后修正
	if IsTakeOver() {
		return nil // 立刻中断
	}

	// 以下保持原有程式
	if dc.conn != nil {
		dc.conn.Close()
	}

	typ := "tcp"
	if strings.Contains(dc.addr, "/") {
		typ = "unix"
	}

	netConn, err := net.Dial(typ, dc.addr)
	if err != nil {
		return err
	}

	tcpConn := netConn.(*net.TCPConn)
	// SetNoDelay controls whether the operating system should delay packet transmission
	// in hopes of sending fewer packets (Nagle's algorithm).
	// The default is true (no delay),
	// meaning that data is sent as soon as possible after a Write.
	tcpConn.SetNoDelay(true)
	tcpConn.SetKeepAlive(true)
	dc.conn = mysql.NewConn(tcpConn)

	// step1: read handshake requirements
	if err := dc.readInitialHandshake(); err != nil {
		dc.conn.Close()
		return err
	}

	// step2: write handshake response
	if err := dc.writeHandshakeResponse41(); err != nil {
		dc.conn.Close()

		return err
	}

	response, err := dc.readPacket()
	if err != nil {
		dc.conn.Close()
		return err
	}

	switch response[0] {
	case mysql.OKHeader:
	default:
		return errors.New("dc connection handshake failed with mysql")
	}

	// we must always use autocommit
	if !dc.IsAutoCommit() {
		if _, err := dc.exec("set autocommit = 1", 0); err != nil {
			dc.conn.Close()

			return err
		}
	}

	return nil
}

// Close close connection to backend mysql and reset conn structure
func (dc *DirectConnection) Close() {
	if dc.conn != nil {
		dc.conn.Close()
	}

	dc.conn = nil
	dc.salt = nil
	dc.pkgErr = nil
	dc.closed.Set(true)

	return
}

// IsClosed check if connection closed
func (dc *DirectConnection) IsClosed() bool {
	return dc.closed.Get()
}

// readPacket doesn't use EphemeralBuffer
func (dc *DirectConnection) readPacket() ([]byte, error) {
	data, err := dc.conn.ReadPacket()
	dc.pkgErr = err
	return data, err
}

// writePacket doesn't use EphemeralBuffer
func (dc *DirectConnection) writePacket(data []byte) error {
	err := dc.conn.WritePacket(data)
	if err != nil && strings.Contains(err.Error(), "broken pipe") {
		// retry 3 times, close dc's conn、reset dc's stats and reconnect
		for i := 0; i < 3; i++ {
			dc.Close()
			e := dc.connect()
			if e == nil { // no need to write data again
				break
			}
		}

	}
	return err
}

// writeEphemeralPacket
func (dc *DirectConnection) writeEphemeralPacket() error {
	err := dc.conn.WriteEphemeralPacket()
	if err != nil && strings.Contains(err.Error(), "broken pipe") {
		// retry 3 times, close dc's conn、reset dc's stats and reconnect
		for i := 0; i < 3; i++ {
			dc.Close()
			e := dc.connect()
			if e == nil { // no need to write data again and ephemeral buffer is recycled
				break
			}
		}
	}
	return err
}

func (dc *DirectConnection) readInitialHandshake() error {
	data, err := dc.readPacket()
	if err != nil {
		return err
	}

	if data[0] == mysql.ErrHeader {
		return errors.New("read initial handshake error")
	}

	if data[0] < mysql.MinProtocolVersion {
		return fmt.Errorf("invalid protocol version %d, must >= 10", data[0])
	}

	//skip mysql version
	//mysql version end with 0x00
	pos := 1 + bytes.IndexByte(data[1:], 0x00) + 1

	// get connection id
	dc.conn.ConnectionID = binary.LittleEndian.Uint32(data[pos : pos+4])

	pos += 4

	dc.salt = append(dc.salt, data[pos:pos+8]...)

	//skip filter
	pos += 8 + 1

	//capability lower 2 bytes
	dc.capability = uint32(binary.LittleEndian.Uint16(data[pos : pos+2]))

	pos += 2

	if len(data) > pos {
		//skip server charset
		//c.charset = data[pos]
		pos++

		dc.status = binary.LittleEndian.Uint16(data[pos : pos+2])
		pos += 2

		dc.capability = uint32(binary.LittleEndian.Uint16(data[pos:pos+2]))<<16 | dc.capability

		pos += 2

		//skip auth data len or [00]
		//skip reserved (all [00])
		pos += 10 + 1

		// The documentation is ambiguous about the length.
		// The official Python library uses the fixed length 12
		// mysql-proxy also use 12
		// which is not documented but seems to work.
		dc.salt = append(dc.salt, data[pos:pos+12]...)
	}

	return nil
}

// writeHandshakeResponse41 writes the handshake response.
func (dc *DirectConnection) writeHandshakeResponse41() error {
	// Adjust client capability flags based on server support
	capability := mysql.ClientProtocol41 | mysql.ClientSecureConnection |
		mysql.ClientLongPassword | mysql.ClientTransactions | mysql.ClientLongFlag
	capability &= dc.capability

	//we only support secure connection
	auth := mysql.CalcPassword(dc.salt, []byte(dc.password))

	length := 4 + // Client capability flags
		4 + // Max-packet size.
		1 + // Character set.
		23 + // Reserved.
		mysql.LenNullString(dc.user) + // user
		1 +
		len(auth)

	if len(dc.db) > 0 {
		capability |= mysql.ClientConnectWithDB
		length += mysql.LenNullString(dc.db)
	}

	dc.capability = capability

	data := make([]byte, length, length)
	pos := 0

	// Client capability flags.
	pos = mysql.WriteUint32(data, pos, capability)

	// Max-packet size, always 0. See doc.go.
	pos = mysql.WriteZeroes(data, pos, 4)

	// Character set.
	pos = mysql.WriteByte(data, pos, byte(dc.collation))

	// 23 reserved bytes, all 0.
	pos = mysql.WriteZeroes(data, pos, 23)

	// user type: null terminated string
	pos = mysql.WriteNullString(data, pos, dc.user)

	// auth [length encoded integer]
	data[pos] = byte(len(auth))
	pos++
	pos += copy(data[pos:], auth)

	// db type: null terminated string
	if len(dc.db) > 0 {
		pos = mysql.WriteNullString(data, pos, dc.db)
	}

	if err := dc.writePacket(data); err != nil {
		return err
	}

	return nil
}

// writeComInitDB changes the default database to use.
// Client -> Server.DirectConnection
// Returns SQLError(CRServerGone) if it can't.
func (dc *DirectConnection) writeComInitDB(db string) error {
	dc.conn.SetSequence(0)
	data := make([]byte, len(db)+1, len(db)+1)
	data[0] = mysql.ComInitDB
	copy(data[1:], db)
	if err := dc.writePacket(data); err != nil {
		return err
	}
	return nil
}

// writeComQuery send ComQuery request use EphemeralBuffer
func (dc *DirectConnection) writeComQuery(sql string) error {
	dc.conn.SetSequence(0)
	data := dc.conn.StartEphemeralPacket(len(sql) + 1)
	data[0] = mysql.ComQuery
	copy(data[1:], sql)
	if err := dc.writeEphemeralPacket(); err != nil {
		return err
	}
	return nil
}

func (dc *DirectConnection) writeComFieldList(table string, wildcard string) error {
	dc.conn.SetSequence(0)
	length := 1 +
		mysql.LenNullString(table) +
		mysql.LenNullString(wildcard)

	data := make([]byte, length, length)
	pos := 0

	pos = mysql.WriteByte(data, 0, mysql.ComFieldList)
	pos = mysql.WriteNullString(data, pos, table)
	pos = mysql.WriteNullString(data, pos, wildcard)

	if err := dc.writePacket(data); err != nil {
		return err
	}

	return nil
}

// Ping implements mysql ping command.
func (dc *DirectConnection) Ping() error {
	dc.conn.SetSequence(0)
	if err := dc.writePacket([]byte{mysql.ComPing}); err != nil {
		return err
	}
	data, err := dc.readPacket()
	if err != nil {
		return err
	}
	switch data[0] {
	case mysql.OKHeader:
		return nil
	case mysql.ErrHeader:
		return errors.New("dc connection ping failed")
	}
	return fmt.Errorf("unexpected packet type: %d", data[0])
}

// UseDB send ComInitDB to backend mysql
func (dc *DirectConnection) UseDB(dbName string) error {

	// 🧚 直接由单元测试接管
	if IsTakeOver() {
		return nil // 立刻中断
	}

	dc.conn.SetSequence(0)
	if dc.db == dbName || len(dbName) == 0 {
		return nil
	}

	if err := dc.writeComInitDB(dbName); err != nil {
		return err
	}

	if r, err := dc.readPacket(); err != nil {
		return err
	} else if !mysql.IsOKPacket(r) {
		return errors.New("dc connection use db failed")
	}

	dc.db = dbName
	return nil
}

// GetDB return database name
func (dc *DirectConnection) GetDB() string {
	return dc.db
}

// GetAddr return addr of backend mysql
func (dc *DirectConnection) GetAddr() string {
	return dc.addr
}

// Execute send ComQuery or ComStmtPrepare/ComStmtExecute/ComStmtClose to backend mysql
func (dc *DirectConnection) Execute(sql string, maxRows int) (*mysql.Result, error) {
	// 🧚 只要这个物件 dc *DirectConnection 一初始化时，假资料的产生方式在函式 NewDirectConnection 就决定了
	// 并在 dc.Trans.MarkLoaded() 这一行完成测试资料载入

	// 🧚 直接由单元测试接管
	if IsTakeOver() {
		dc.MockDC = new(MockDcClient)
		dc.MockDC.MockKey = dc.MakeMockKey(sql) // 3652007921

		if tmp, ok := FakeDBInstance.MockResult[dc.MockDC.MockKey]; ok {
			fmt.Printf("\u001B[35m 命中测试资料序号 Key: %d\n", dc.MockDC.MockKey)
			return &tmp, nil // 立刻中斷
		} else {
			// 做成对应到 网路位置、帐号、密码等相关资料，会回传 SQL 的执行结果
			// 如果在执行单元测试过程中，没有 1 命中单元测试的测试资料 或者是 2 测试资料载入失败，就使用 Fatal 中止
			// Fatal 中止 只有在单元测试的环境下才会执行，不会影响到主程式，还算安全
			log.Fatal("单元测试没有命中测试资料序号 %d", dc.MockDC.MockKey) // 直接中断
		}
	}

	// 以下保持原有程式
	return dc.exec(sql, maxRows)
}

// Begin send ComQuery with 'begin' to backend mysql to start transaction
func (dc *DirectConnection) Begin() error {
	_, err := dc.exec("begin", 0)
	return err
}

// Commit send ComQuery with 'commit' to backend mysql to commit transaction
func (dc *DirectConnection) Commit() error {
	_, err := dc.exec("commit", 0)
	return err
}

// Rollback send ComQuery with 'rollback' to backend mysql to rollback transaction
func (dc *DirectConnection) Rollback() error {
	_, err := dc.exec("rollback", 0)
	return err
}

// SetAutoCommit trun on/off autocommit
func (dc *DirectConnection) SetAutoCommit(v uint8) error {
	if v == 0 {
		if _, err := dc.exec("set autocommit = 0", 0); err != nil {
			dc.conn.Close()

			return err
		}
	} else {
		if _, err := dc.exec("set autocommit = 1", 0); err != nil {
			dc.conn.Close()

			return err
		}
	}
	return nil
}

// SetCharset set charset of connection to backend mysql
func (dc *DirectConnection) SetCharset(charset string, collation mysql.CollationID) ( /*changed*/ bool, error) {
	charset = strings.Trim(charset, "\"'`")

	if collation == 0 {
		collation = mysql.CollationNames[mysql.Charsets[charset]]
	}

	if dc.charset == charset && dc.collation == collation {
		return false, nil
	}

	_, ok := mysql.CharsetIds[charset]
	if !ok {
		return false, fmt.Errorf("invalid charset %s", charset)
	}

	_, ok = mysql.Collations[collation]
	if !ok {
		return false, fmt.Errorf("invalid collation %d", collation)
	}

	dc.collation = collation
	dc.charset = charset
	return true, nil
}

// ResetConnection reset connection stattus, include transaction、autocommit、charset、sql_mode .etc
func (dc *DirectConnection) ResetConnection() error {
	if dc.IsInTransaction() {
		log.Debug("get transaction connection from pool, addr: %s, user: %s, db: %s, status: %d", dc.addr, dc.user, dc.db, dc.status)
		if err := dc.Rollback(); err != nil {
			log.Warn("rollback in reset connection error, addr: %s, user: %s, db: %s, status: %d, err: %v", dc.addr, dc.user, dc.db, dc.status, err)
			return err
		}
	}

	if !dc.IsAutoCommit() {

		// 🧚 直接由单元测试接管
		if IsTakeOver() {
			return nil // 立刻中斷
		}

		log.Debug("get autocommit = 0 connection from pool, addr: %s, user: %s, db: %s, status: %d", dc.addr, dc.user, dc.db, dc.status)
		if err := dc.SetAutoCommit(1); err != nil {
			log.Warn("set autocommit = 1 in reset connection error, addr: %s, user: %s, db: %s, status: %d, err: %v", dc.addr, dc.user, dc.db, dc.status, err)
			return err
		}
	}

	return nil
}

// SetSessionVariables set direction variables according to Session
func (dc *DirectConnection) SetSessionVariables(frontend *mysql.SessionVariables) (bool, error) {
	return dc.sessionVariables.SetEqualsWith(frontend)
}

// WriteSetStatement execute sql
func (dc *DirectConnection) WriteSetStatement() error {
	var setVariableSQL bytes.Buffer
	collation, ok := mysql.Collations[dc.collation]
	if !ok {
		return fmt.Errorf("invalid collationId: %v", dc.collation)
	}
	appendSetCharset(&setVariableSQL, dc.charset, collation)

	for _, v := range dc.sessionVariables.GetAll() {
		appendSetVariable(&setVariableSQL, v.Name(), v.Get())
	}

	for _, v := range dc.sessionVariables.GetUnusedAndClear() {
		appendSetVariableToDefault(&setVariableSQL, v.Name())
	}

	setSQL := setVariableSQL.String()
	if setSQL == "" {
		return nil
	}

	// 🧚 直接由单元测试接管
	if IsTakeOver() {
		return nil // 立刻中断
	}

	if _, err := dc.exec(setSQL, 0); err != nil {
		return err
	}
	return nil
}

// FieldList send ComFieldList to backend mysql
func (dc *DirectConnection) FieldList(table string, wildcard string) ([]*mysql.Field, error) {
	if err := dc.writeComFieldList(table, wildcard); err != nil {
		return nil, err
	}
	fs := make([]*mysql.Field, 0, 4)
	var f *mysql.Field
	for {
		data, err := dc.readPacket()
		if err != nil {
			return nil, err
		}

		// EOF Packet
		if dc.isEOFPacket(data) {
			return fs, nil
		}

		if data[0] == mysql.ErrHeader {
			return nil, dc.handleErrorPacket(data)
		}

		if f, err = mysql.FieldData(data).Parse(); err != nil {
			return nil, err
		}
		fs = append(fs, f)
	}
}

// execute ComQuery command
func (dc *DirectConnection) exec(query string, maxRows int) (*mysql.Result, error) {
	if err := dc.writeComQuery(query); err != nil {
		return nil, err
	}

	return dc.readResult(false, maxRows)
}

// read resultset from mysql
func (dc *DirectConnection) readResultSet(data []byte, binary bool, maxRows int) (*mysql.Result, error) {
	result := &mysql.Result{
		Status:       0,
		InsertID:     0,
		AffectedRows: 0,

		Resultset: &mysql.Resultset{},
	}

	// column count
	pos := 0
	count, pos, _, _ := mysql.ReadLenEncInt(data, pos)

	if pos-len(data) != 0 {
		return nil, mysql.ErrMalformPacket
	}

	result.Fields = make([]*mysql.Field, count)
	result.FieldNames = make(map[string]int, count)

	if err := dc.readResultColumns(result); err != nil {
		return nil, err
	}

	if err := dc.readResultRows(result, binary, maxRows); err != nil {
		return nil, err
	}

	return result, nil
}

// readResultColumns read column information
func (dc *DirectConnection) readResultColumns(result *mysql.Result) (err error) {
	var i = 0
	var data []byte

	for {
		data, err = dc.readPacket()
		if err != nil {
			return
		}

		// EOF Packet
		if dc.isEOFPacket(data) {
			if dc.capability&mysql.ClientProtocol41 > 0 {
				//result.Warnings = binary.LittleEndian.Uint16(data[1:])
				//todo add strict_mode, warning will be treat as error
				result.Status = binary.LittleEndian.Uint16(data[3:])
				dc.status = result.Status
			}

			if i != len(result.Fields) {
				err = mysql.ErrMalformPacket
			}

			return
		}

		if data[0] == mysql.ErrHeader {
			return dc.handleErrorPacket(data)
		}

		result.Fields[i], err = mysql.FieldData(data).Parse()
		if err != nil {
			return
		}

		result.FieldNames[string(result.Fields[i].Name)] = i

		i++
	}
}

// readResultRows read result rows
func (dc *DirectConnection) readResultRows(result *mysql.Result, isBinary bool, maxRows int) (err error) {
	var data []byte

	for {
		data, err = dc.readPacket()
		if err != nil {
			return
		}

		// EOF Packet
		if dc.isEOFPacket(data) {
			if dc.capability&mysql.ClientProtocol41 > 0 {
				//result.Warnings = binary.LittleEndian.Uint16(data[1:])
				//todo add strict_mode, warning will be treat as error
				result.Status = binary.LittleEndian.Uint16(data[3:])
				dc.status = result.Status
			}

			break
		}

		if data[0] == mysql.ErrHeader {
			return dc.handleErrorPacket(data)
		}

		result.RowDatas = append(result.RowDatas, data)
		if maxRows > 0 && len(result.RowDatas) >= maxRows {
			if err := dc.drainResults(); err != nil {
				return fmt.Errorf("%v %d, drain error: %v", sqlerr.ErrRowsLimitExceeded, maxRows, err)
			}
			return fmt.Errorf("%v %d", sqlerr.ErrRowsLimitExceeded, maxRows)
		}
	}

	result.Values = make([][]interface{}, len(result.RowDatas))
	for i := range result.Values {
		result.Values[i], err = result.RowDatas[i].Parse(result.Fields, isBinary)
		if err != nil {
			return err
		}
	}

	return nil
}

// drainResults will read all packets for a result set and ignore them.
func (dc *DirectConnection) drainResults() error {
	for {
		data, err := dc.conn.ReadEphemeralPacket()
		if err != nil {
			dc.conn.RecycleReadPacket()
			return err
		}

		if dc.isEOFPacket(data) {
			dc.conn.RecycleReadPacket()
			return nil
		} else if data[0] == mysql.ErrHeader {
			err := dc.handleErrorPacket(data)
			dc.conn.RecycleReadPacket()
			return err
		}
		dc.conn.RecycleReadPacket()
	}
}

func (dc *DirectConnection) isEOFPacket(data []byte) bool {
	return data[0] == mysql.EOFHeader && len(data) <= 5
}

func (dc *DirectConnection) handleOKPacket(data []byte) (*mysql.Result, error) {
	var pos = 1

	r := new(mysql.Result)

	r.AffectedRows, pos, _, _ = mysql.ReadLenEncInt(data, pos)
	r.InsertID, pos, _, _ = mysql.ReadLenEncInt(data, pos)

	if dc.capability&mysql.ClientProtocol41 > 0 {
		r.Status = binary.LittleEndian.Uint16(data[pos:])
		dc.status = r.Status
		pos += 2

		// TODO strict_mode, check warnings as error
		// Warnings := binary.LittleEndian.Uint16(data[pos:])
		// pos += 2
	} else if dc.capability&mysql.ClientTransactions > 0 {
		r.Status = binary.LittleEndian.Uint16(data[pos:])
		dc.status = r.Status
		pos += 2
	}

	//info
	return r, nil
}

func (dc *DirectConnection) handleErrorPacket(data []byte) error {
	e := new(mysql.SQLError)

	var pos = 1

	e.Code = binary.LittleEndian.Uint16(data[pos:])
	pos += 2

	if dc.capability&mysql.ClientProtocol41 > 0 {
		// skip '#'
		pos++
		e.State = string(data[pos : pos+5])
		pos += 5
	}

	e.Message = string(data[pos:])

	return e
}

func (dc *DirectConnection) readResult(binary bool, maxRows int) (*mysql.Result, error) {
	data, err := dc.readPacket()
	if err != nil {
		return nil, err
	}
	if data[0] == mysql.OKHeader {
		return dc.handleOKPacket(data)
	} else if data[0] == mysql.ErrHeader {
		return nil, dc.handleErrorPacket(data)
	} else if data[0] == mysql.LocalInFileHeader {
		return nil, mysql.ErrMalformPacket
	}

	return dc.readResultSet(data, binary, maxRows)
}

// IsAutoCommit check if autocommit
func (dc *DirectConnection) IsAutoCommit() bool {
	return dc.status&mysql.ServerStatusAutocommit > 0
}

// IsInTransaction check if in transaction
func (dc *DirectConnection) IsInTransaction() bool {
	return dc.status&mysql.ServerStatusInTrans > 0
}

// GetCharset return charset of specific connection
func (dc *DirectConnection) GetCharset() string {
	return dc.charset
}

func appendSetCharset(buf *bytes.Buffer, charset string, collation string) {
	if buf.Len() != 0 {
		buf.WriteString(",")
	} else {
		buf.WriteString("SET NAMES '")
	}
	buf.WriteString(charset)
	buf.WriteString("' COLLATE '")
	buf.WriteString(collation)
	buf.WriteString("'")
}

func appendSetVariable(buf *bytes.Buffer, key string, value interface{}) {
	if buf.Len() != 0 {
		buf.WriteString(",")
	} else {
		buf.WriteString("SET ")
	}
	buf.WriteString(key)
	buf.WriteString(" = ")
	switch v := value.(type) {
	case string:
		if strings.ToLower(v) == mysql.KeywordDefault {
			buf.WriteString(v)
		} else {
			buf.WriteString("'")
			buf.WriteString(v)
			buf.WriteString("'")
		}
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		buf.WriteString(fmt.Sprintf("%d", v))
	default:
		buf.WriteString("'")
		buf.WriteString(fmt.Sprintf("%v", v))
		buf.WriteString("'")
	}
}

func appendSetVariableToDefault(buf *bytes.Buffer, key string) {
	if buf.Len() != 0 {
		buf.WriteString(",")
	} else {
		buf.WriteString("SET ")
	}
	buf.WriteString(key)
	buf.WriteString(" = ")
	buf.WriteString("DEFAULT")
}
