package containerdTest

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"sync"
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

// BenchmarkContainerdManager_Lock 为容器服务管理員加锁的性能测试 BenchmarkContainerdManager_Lock is a benchmark for ContainerManager lock
func BenchmarkContainerdManager_Lock(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := Manager.GetBuilder("mariadb-server")
		if err != nil {
			panic(err)
		}
		if err == nil {
			err = Manager.ReturnBuilder("mariadb-server")
			if err != nil {
				panic(err)
			}
		}
	}
}

// TestContainerdManager_Lock 测试容器服务管理器加锁的正确性测试 TestContainerdManager_Lock is a testing for ContainerManager lock
func TestContainerdManager_Lock(t *testing.T) {
	// it is the design of the container manager lock by using map and sync.mux
	// map and inner lock design 为使用 map 和 sync.mux 的性能测试
	t.Run("map and inner lock design", func(t *testing.T) {
		// use this design to make container manager lock
		// 使用这种设计，可以让容器管理器加锁
		type data struct {
			dataInt int
			mutex   sync.Locker
		}
		map1 := make(map[string]*data)
		map1["server"] = &data{dataInt: 0, mutex: &sync.Mutex{}}
		wg := sync.WaitGroup{}
		wg.Add(10000)
		for i := 0; i < 10000; i++ {
			go func(j int) {
				map1["server"].mutex.Lock()
				map1["server"].dataInt++
				_ = map1["server"].dataInt
				map1["server"].mutex.Unlock()
				wg.Done()
			}(i)
		}
		wg.Wait()
		require.Equal(t, 10000, map1["server"].dataInt)
	})
	// it is the design of the container manager lock by using sync.map
	// sync map design 为使用 sync.Map 的性能测试
	t.Run("sync map design", func(t *testing.T) {
		type data struct {
			dataInt int
			mutex   sync.Locker
		}
		map1 := sync.Map{}
		map1.Store("server", &data{dataInt: 0, mutex: &sync.Mutex{}})
		wg := sync.WaitGroup{}
		wg.Add(10000)
		for i := 0; i < 10000; i++ {
			go func(j int) {
				tmp, _ := map1.Load("server")
				tmp.(*data).mutex.Lock()
				tmp.(*data).dataInt++
				tmp.(*data).mutex.Unlock()
				wg.Done()
			}(i)
		}
		wg.Wait()
		tmp, _ := map1.Load("server")
		require.Equal(t, 10000, tmp.(*data).dataInt)
	})
}
