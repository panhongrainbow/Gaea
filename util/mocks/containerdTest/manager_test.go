package containerdTest

import (
	"fmt"
	"testing"
	"time"
)

func TestManager(t *testing.T) {
	//
	mgr, _ := NewContainderManager("./example/")
	err := mgr.Builder["mariadb-server"].Build(60 * time.Second)
	fmt.Println(err)

	time.Sleep(time.Second * 30)

	_ = mgr.Builder["mariadb-server"].TearDown(60 * time.Second)
}
