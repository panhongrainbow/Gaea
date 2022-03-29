package soarTest

import (
	"fmt"
	"net"
	"strings"
)

func init() {
	//
}

// setup 连线测试的到取得资料库版本就为正确
func setup() error {
	typ := "tcp"
	if strings.Contains("192.168.122.2:3309", "/") {
		typ = "unix"
	}

	netConn, err := net.Dial(typ, "192.168.122.2:3309")
	if err != nil {
		return err
	}

	// 先随意测试
	test := make([]byte, 20)
	netConn.Read(test)
	fmt.Println(test)

	_ = netConn.Close()

	return nil
}
