package containerdTest

import "time"

// Builder 是操作一个新的测试环境的接口 interface for making a new test environment.
type Builder interface {
	Build(t time.Duration) error
	TearDown(t time.Duration) error
}

// NewBuilder 为建立一个新的 Builder 接口 NewBuilder is a function to create a new Builder interface.
func NewBuilder(cfg ContainerD) (Builder, error) {
	return NewContainerdClient(cfg)
}
