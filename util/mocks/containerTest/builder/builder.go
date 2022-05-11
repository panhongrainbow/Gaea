package builder

import (
	"time"
)

// Builder 是操作一个新的测试容器环境的接口
// Builder is an interface for making a new container test environment.
type Builder interface {
	Build(t time.Duration) error
	OnService(t time.Duration) error
	TearDown(t time.Duration) error
}
