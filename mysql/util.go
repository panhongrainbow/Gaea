// Copyright 2016 The kingshard Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

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

package mysql

import (
	"crypto/sha1"
	"crypto/sha256"
	"math/rand"
	"time"
	"unicode/utf8"
)

var (
	dontEscape = byte(255)
	encodeMap  [256]byte
)

// CalcPassword calculate password hash
func CalcPassword(scramble, password []byte) []byte {
	// 如果客户端没有输入密码，就直接中断回传
	if len(password) == 0 {
		return nil
	}

	// 1 公式说明
	// https://dev.mysql.com/doc/internals/en/secure-password-authentication.html
	// 安全密码产生规则如下，并整理成以下公式
	// SHA1( password ) XOR SHA1( "20-bytes random data from server" <concat> SHA1( SHA1( password ) ) )
	// 修改成 SHA1( password ) XOR SHA1( scramble <直接连接> SHA1( stage1 ) )
	//     stage1 = SHA1( password )
	//     stage1Hash = SHA1( stage1 ) = SHA1( SHA1( password ) )
	//     scramble = SHA1( scramble + SHA1( stage1Hash ) ) // 第一次修改 scramble 的数值
	//     scramble = stage1 XOR scramble // 第二次修改 scramble 的数值

	// 2 假设参数传入
	// 2-1 假设 输入密码参数 password
	// 客户端的输入密码 password 为 12345
	// 2-2 假设 数数库服务器回传的乱数 scramble
	// 十进位
	// scramble 为 []uint8{81, 64, 43, 85, 76, 90, 97, 91, 34, 53, 36, 85, 93, 86, 117, 105, 49, 87, 65, 125}
	// 十六进位
	// scramble 为 []uint8{51, 40, 2B, 55, 4c, 5a, 61, 5b, 22, 35, 24, 55, 5d, 56,  75,  69, 31, 57, 41,  7d}
	// 在 Bash 的十六进位输入
	// scramble 的值为 51402B554c5A615b223524555d5675693157417d

	// 3 计算 stage1
	// stage1 = SHA1(password)
	crypt := sha1.New()
	crypt.Write(password)
	stage1 := crypt.Sum(nil)

	// 使用 Linux Bash 去验证 stage1
	// echo 12345 | sha1sum | head -c 40 # 把 密码 12345 转成 stage1
	// 8cb2237d0679ca88db6464eac60da96345513964 # 此为 stage1 的值

	// 4 计算 stage1Hash
	// inner Hash
	crypt.Reset()
	crypt.Write(stage1)
	stage1hash := crypt.Sum(nil)

	// 使用 Linux Bash 去验证 stage1Hash
	// echo -n 12345 | sha1sum | xxd -r -p | sha1sum | head -c 40
	// 也可似用以下命令求出 echo -n 8cb2237d0679ca88db6464eac60da96345513964 | xxd -r -p | sha1sum | head -c 40
	// 00a51f3f48415c7d4e8908980d443c29c69b60c9 # 此为 stage1hash 的值

	// 5 第一次重写 scramble
	// outer Hash
	crypt.Reset()
	crypt.Write(scramble)
	crypt.Write(stage1hash)
	scramble = crypt.Sum(nil)

	// 使用 Linux Bash 去验证 第一次重写 scramble
	// scramble 的值为 51402B554c5A615b223524555d5675693157417d，为连接的前半段
	// stage1Hash 的值为 00a51f3f48415c7d4e8908980d443c29c69b60c9，为连接的后半段
	// echo 51402B554c5A615b223524555d5675693157417d 00a51f3f48415c7d4e8908980d443c29c69b60c9 | xxd -r -p | sha1sum | head -c 40
	// 0ca0f764a59d1cdb10a87f0155d61aa54be1c71a # 此为第一次修改的 scramble

	// 6 第二次重写 scramble
	// token = scrambleHash XOR stage1Hash
	for i := range scramble {
		scramble[i] ^= stage1[i]
	}

	// 使用 Linux Bash 去验证 第二次重写 scramble
	// stage1=0x8cb2237d0679ca88db6464eac60da96345513964 # 之前计算出来的 stage1 的数值
	// scrambleFirst=0x0ca0f764a59d1cdb10a87f0155d61aa54be1c71a # 第一次修改 scramble 的数值
	// echo $(( $stage1^$scrambleFirst ))
	// -7792437067003134338 # 十进位的结果，转换成 十六进位为 0x 93 DB B3 C6 0E B0 FE 7E，因为 Bash 精度不够，只能显示部份数度进行检查

	// 第二次重写 scramble 正确的数值
	// scramble 为 []uint8{128, 18, 212, 25, 163, 228, 214, 83, 203, 204, 27, 235, 147, 219, 179, 198, 14, 176, 254, 126} // 十进位
	// scramble 为 []uint8{ 80, 12,  D4, 19,  A3,  E4,  D6, 53,  CB,  CC, 1B,  EB,  93,  DB,  B3,  C6, 0E,  B0,  FE,  7E} // 十六进位
	// scramble 的值为 8012D419A3E4D653CBCC1BEB93DBB3C60EB0FE7E // 十六进位

	// 最后回传 第二次重写 scramble 正确的数值
	return scramble
}

func CalcCachingSha2Password(salt []byte, password string) []byte {
	if len(password) == 0 {
		return nil
	}
	// XOR(SHA256(password), SHA256(SHA256(SHA256(password)), salt))
	crypt := sha256.New()
	crypt.Write([]byte(password))
	message1 := crypt.Sum(nil)

	crypt.Reset()
	crypt.Write(message1)
	message1Hash := crypt.Sum(nil)

	crypt.Reset()
	crypt.Write(message1Hash)
	crypt.Write(salt)
	message2 := crypt.Sum(nil)

	for i := range message1 {
		message1[i] ^= message2[i]
	}

	return message1
}

// RandomBuf return random salt, seed must be in the range of ascii
func RandomBuf(size int) ([]byte, error) {
	buf := make([]byte, size)
	rand.Seed(time.Now().UTC().UnixNano())
	min, max := 30, 127
	for i := 0; i < size; i++ {
		buf[i] = byte(min + rand.Intn(max-min))
	}
	return buf, nil
}

// Escape remove exceptional character
func Escape(sql string) string {
	dest := make([]byte, 0, 2*len(sql))

	for i, w := 0, 0; i < len(sql); i += w {
		runeValue, width := utf8.DecodeRuneInString(sql[i:])
		if c := encodeMap[byte(runeValue)]; c == dontEscape {
			dest = append(dest, sql[i:i+width]...)
		} else {
			dest = append(dest, '\\', c)
		}
		w = width
	}

	return string(dest)
}

var encodeRef = map[byte]byte{
	'\x00': '0',
	'\'':   '\'',
	'"':    '"',
	'\b':   'b',
	'\n':   'n',
	'\r':   'r',
	'\t':   't',
	26:     'Z', // ctl-Z
	'\\':   '\\',
}

type lengthAndDecimal struct {
	length  int
	decimal int
}

// defaultLengthAndDecimal provides default Flen and Decimal for fields
// from CREATE TABLE when they are unspecified.
var defaultLengthAndDecimal = map[byte]lengthAndDecimal{
	TypeBit:        {1, 0},
	TypeTiny:       {4, 0},
	TypeShort:      {6, 0},
	TypeInt24:      {9, 0},
	TypeLong:       {11, 0},
	TypeLonglong:   {20, 0},
	TypeDouble:     {22, -1},
	TypeFloat:      {12, -1},
	TypeNewDecimal: {11, 0},
	TypeDuration:   {10, 0},
	TypeDate:       {10, 0},
	TypeTimestamp:  {19, 0},
	TypeDatetime:   {19, 0},
	TypeYear:       {4, 0},
	TypeString:     {1, 0},
	TypeVarchar:    {5, 0},
	TypeVarString:  {5, 0},
	TypeTinyBlob:   {255, 0},
	TypeBlob:       {65535, 0},
	TypeMediumBlob: {16777215, 0},
	TypeLongBlob:   {4294967295, 0},
	TypeJSON:       {4294967295, 0},
	TypeNull:       {0, 0},
	TypeSet:        {-1, 0},
	TypeEnum:       {-1, 0},
}

// IsIntegerType indicate whether tp is an integer type.
func IsIntegerType(tp byte) bool {
	switch tp {
	case TypeTiny, TypeShort, TypeInt24, TypeLong, TypeLonglong:
		return true
	}
	return false
}

// GetDefaultFieldLengthAndDecimal returns the default display length (flen) and decimal length for column.
// Call this when no Flen assigned in ddl.
// or column value is calculated from an expression.
// For example: "select count(*) from t;", the column type is int64 and Flen in ResultField will be 21.
// See https://dev.mysql.com/doc/refman/5.7/en/storage-requirements.html
func GetDefaultFieldLengthAndDecimal(tp byte) (flen int, decimal int) {
	val, ok := defaultLengthAndDecimal[tp]
	if ok {
		return val.length, val.decimal
	}
	return -1, -1
}

// defaultLengthAndDecimal provides default Flen and Decimal for fields
// from CAST when they are unspecified.
var defaultLengthAndDecimalForCast = map[byte]lengthAndDecimal{
	TypeString:     {0, -1}, // Flen & Decimal differs.
	TypeDate:       {10, 0},
	TypeDatetime:   {19, 0},
	TypeNewDecimal: {11, 0},
	TypeDuration:   {10, 0},
	TypeLonglong:   {22, 0},
	TypeJSON:       {4194304, 0}, // Flen differs.
}

// GetDefaultFieldLengthAndDecimalForCast returns the default display length (flen) and decimal length for casted column
// when flen or decimal is not specified.
func GetDefaultFieldLengthAndDecimalForCast(tp byte) (flen int, decimal int) {
	val, ok := defaultLengthAndDecimalForCast[tp]
	if ok {
		return val.length, val.decimal
	}
	return -1, -1
}

func init() {
	for i := range encodeMap {
		encodeMap[i] = dontEscape
	}
	for i := range encodeMap {
		if to, ok := encodeRef[byte(i)]; ok {
			encodeMap[byte(i)] = to
		}
	}
}
