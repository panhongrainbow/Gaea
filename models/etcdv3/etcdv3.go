// Copyright 2016 CodisLabs. All Rights Reserved.
// Licensed under the MIT (MIT-LICENSE.txt) license.

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

package etcdclientv3

import (
	"context"
	"errors"
	"github.com/coreos/etcd/client"
	"github.com/coreos/etcd/clientv3"
	"strings"
	"sync"
	"time"
)

// ErrClosedEtcdClient means etcd client closed
var ErrClosedEtcdClient = errors.New("use of closed etcd client")

const (
	defaultEtcdPrefix = "/gaea"
)

// EtcdClientV3 为新版的 etcd client
type EtcdClientV3 struct {
	sync.Mutex
	kapi clientv3.Client // 就这里改成新版的 API

	closed  bool
	timeout time.Duration
	Prefix  string
}

// New 建立新的 etcd v3 的客户端
func New(addr string, timeout time.Duration, username, passwd, root string) (*EtcdClientV3, error) {
	endpoints := strings.Split(addr, ",")
	for i, s := range endpoints {
		if s != "" && !strings.HasPrefix(s, "http://") {
			endpoints[i] = "http://" + s
		}
	}
	config := clientv3.Config{
		Endpoints:            endpoints,
		Username:             username,
		Password:             passwd,
		DialTimeout:          timeout * time.Second, // 只设定第一次连线时间的逾时，之后不用太担心连线，连线失败后，会自动重连
		DialKeepAliveTimeout: timeout * time.Second, // 之后维持 etcd 连线的逾时
	}
	c, err := clientv3.New(config)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(root) == "" {
		root = defaultEtcdPrefix
	}
	return &EtcdClientV3{
		kapi:    *c,
		timeout: timeout,
		Prefix:  root,
	}, nil
}

// Close 关闭新的 etcd v3 的客户端
func (c *EtcdClientV3) Close() error {
	c.Lock()
	defer c.Unlock()
	if c.closed {
		return nil
	}
	c.closed = true
	return nil
}

func (c *EtcdClientV3) contextWithTimeout() (context.Context, context.CancelFunc) {
	if c.timeout == 0 {
		return context.Background(), func() {}
	}
	return context.WithTimeout(context.Background(), c.timeout)
}

// isErrNodeExists (直接借用 v2 的函式)
func isErrNoNode(err error) bool {
	if err != nil {
		if e, ok := err.(client.Error); ok {
			return e.Code == client.ErrorCodeKeyNotFound
		}
	}
	return false
}

// isErrNodeExists (直接借用 v2 的函式)
func isErrNodeExists(err error) bool {
	if err != nil {
		if e, ok := err.(client.Error); ok {
			return e.Code == client.ErrorCodeNodeExist
		}
	}
	return false
}

// Mkdir create directory
func (c *EtcdClientV3) Mkdir(dir string) error {
	c.Lock()
	defer c.Unlock()
	if c.closed {
		return ErrClosedEtcdClient
	}
	return c.mkdir(dir)
}

func (c *EtcdClientV3) mkdir(dir string) error {
	if dir == "" || dir == "/" {
		return nil
	}
	cntx, canceller := c.contextWithTimeout()
	defer canceller()
	_, err := c.kapi.Put(cntx, dir, "", clientv3.WithPrevKV())
	if err != nil {
		if isErrNodeExists(err) {
			return nil
		}
		return err
	}
	return nil
}
