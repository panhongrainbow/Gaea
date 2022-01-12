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
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_EtcdV3(t *testing.T) {
	// 设定颜色输出
	colorReset := "\033[0m"
	colorRed := "\033[31m"

	_, err := New("http://127.0.0.1:2399", 3, "", "", "")
	if err != nil {
		fmt.Println(colorRed, "目前找不到可实验的 Etcd 服务器", colorReset)
	}
	if err == nil {
		//
	}
}
