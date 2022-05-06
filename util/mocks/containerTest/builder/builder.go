package builder

import (
	"time"
)

// Builder 是操作一个新的测试环境的接口 interface for making a new test environment.
type Builder interface {
	Build(t time.Duration) error
	OnService(t time.Duration) error
	TearDown(t time.Duration) error
}
