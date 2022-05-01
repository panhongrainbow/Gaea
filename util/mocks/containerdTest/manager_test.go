package containerdTest

import (
	"fmt"
	"testing"
	"time"
)

func TestManager(t *testing.T) {
	//
	mgr, _ := NewContainderManager("./example/")
	builder, _ := mgr.GetBuilder("mariadb-server")

	err := builder.Build(60 * time.Second)
	fmt.Println(err)

	time.Sleep(time.Second * 30)

	_ = builder.TearDown(60 * time.Second)
}
