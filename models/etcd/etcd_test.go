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

package etcdclient

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/client"
	"github.com/coreos/etcd/clientv3"
	"testing"
	"time"
)

func Test_isErrNoNode(t *testing.T) {
	err := client.Error{}
	err.Code = client.ErrorCodeKeyNotFound
	if !isErrNoNode(err) {
		t.Fatalf("test isErrNoNode failed, %v", err)
	}
	err.Code = client.ErrorCodeNotFile
	if isErrNoNode(err) {
		t.Fatalf("test isErrNoNode failed, %v", err)
	}
}

func Test_isErrNodeExists(t *testing.T) {
	err := client.Error{}
	err.Code = client.ErrorCodeNodeExist
	if !isErrNodeExists(err) {
		t.Fatalf("test isErrNodeExists failed, %v", err)
	}
	err.Code = client.ErrorCodeNotFile
	if isErrNodeExists(err) {
		t.Fatalf("test isErrNodeExists failed, %v", err)
	}
}

// 以下为新增的测试
// func New(addr string, timeout time.Duration, username, passwd, root string) (*EtcdClient, error) {
func Test_Etcd(t *testing.T) {
	remote, err := New("http://127.0.0.1:2379", 30000*10000, "", "", "")
	fmt.Println("1", err)
	ctx := context.Background()
	test, err := remote.kapi.Set(ctx, "test3", "test5", nil)
	fmt.Println("2", err)
	fmt.Println(test)
}

func Test_Etcd1(t *testing.T) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		// handle error!
	}
	defer cli.Close()

	ctx := context.Background()
	resp, err := cli.Put(ctx, "sample_key", "sample_value")
	fmt.Println(resp)
	if err != nil {
		// handle error!
	}
}

/*
cli, err := clientv3.New(clientv3.Config{
                Endpoints:   []string{"localhost:2379"},
                DialTimeout: 5 * time.Second,
        })
        if err != nil {
                // handle error!
        }
        defer cli.Close()

        ctx := context.Background()
        resp, err := cli.Put(ctx, "sample_key", "sample_value")
        fmt.Println(resp)
        if err != nil {
           // handle error!
        }
*/
