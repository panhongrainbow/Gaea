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
package util

import (
	"github.com/stretchr/testify/require"
	"testing"
)

// Benchmark_RequestContext_Get 函式为多读单写中读的基准测试
func Benchmark_RequestContext_Get(b *testing.B) {
	var rctx = NewRequestContext()
	rctx.Set("key", "value")
	for i := 0; i < b.N; i++ {
		rctx.Get("key")
	}
}

// Test_RequestContext_Set 函式是用来测试多读单写中其中写的特性
func Test_RequestContext_Set(t *testing.T) {
	// 产生一个内含多读单写物件
	var rctx = NewRequestContext()

	// 建立一个通道去记录协程完成的先后顺序
	chanRW := make(chan string, 3)

	// 准备第一个读取函式
	r0 := 0
	funcR0 := func() {
		res := rctx.Get("key")
		r0 = res.(int)
		chanRW <- "r0"
	}

	// 准备第二个读取函式
	r1 := 0
	funcR1 := func() {
		res := rctx.Get("key")
		r1 = res.(int)
		chanRW <- "r1"
	}

	// 准备第一个写入函式
	funcW0 := func(i int) {
		rctx.Set("key", i)
		chanRW <- "w0"
	}

	// 上锁有会成功和失败的状况，准备两个变数进行之后的统计
	correct := 0
	notCorrect := 0

	// 一开始 key 值归 0
	rctx.Set("key", 0)

	// 进行一万笔测试
	for i := 1; i <= 10000; i++ {
		// 准备两读一写共三笔协程
		go funcR0()
		go funcR1()
		go funcW0(i)

		// 记录第二笔完成的协程
		<-chanRW
		secondFinished := <-chanRW
		<-chanRW

		// 如果第二笔完成的协程为写的协程
		if secondFinished == "w0" {
			if r0 != r1 {
				correct++ // 因第二完成的协程为写入，所以另外二个读取的协程取出的值应要不同，这才是正确
			}
			if r0 == r1 {
				notCorrect++ // 反之，就是不正确
			}
		}
	}
	// 计算正确比率
	correctPercent := float32(correct-notCorrect) / float32(correct) * 100
	// 正确比率要超过百分之九十五
	require.Equal(t, correctPercent > 95, true)
	// 取出写入的值
	value := rctx.Get("key")
	// 检查写入的值是否正确
	require.Equal(t, value.(int), 10000)
}
