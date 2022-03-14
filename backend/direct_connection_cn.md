# backend 包的直连 direct_connection



## 代码



## 测试

> 以下会说明在写测试时考量的点

### 匿名函数的考量

在以下代码内有一个 测试函数 t.Run("测试数据库后端连线初始交握后的回应", func(t *testing.T)

此测试函数内含 匿名函数 customFunc

以下代码，匿名函数 customFunc 内的变量将会取 dc 对象的內存位置，考量后，觉得可以这样写

```go
	// 交握第二步 Step2
	t.Run("测试数据库后端连线初始交握后的回应", func(t *testing.T) {
		var connForSengingMsgToMariadb = mysql.NewConn(mockGaea.GetConnWrite())
		dc.conn = connForSengingMsgToMariadb
		dc.conn.StartWriterBuffering()
        
		customFunc := func() {
			err := dc.writeHandshakeResponse41()
			require.Equal(t, err, nil)
			err = dc.conn.Flush()
			require.Equal(t, err, nil)
			err = mockGaea.GetConnWrite().Close()
			require.Equal(t, err, nil)
		}

		fmt.Println(mockGaea.CustomSend(customFunc).ArrivedMsg(mockMariaDB))
	})
```

## 验证

使用 Linux 命令 或者是 网站 去计算 Sha1sum 时，计算出来的结果为 16 进位，只不过 IDE 工具在取中断点时会显示为 10 进位，以下使用 mysql 包里的 CalcPassword 函式中的 stage1 變量为例

### 使用工具和网站把密码转换成验证码

使用 Linux Command 去产生 stage1 的 Sha1sum 验证码

<img src="./assets/image-20220314214316673.png" alt="image-20220314214316673" style="zoom:80%;" /> 

使用网站 https://coding.tools/tw/sha1 去计算 stage1 的 Sha1sum 验证码

<img src="./assets/image-20220314215924425.png" alt="image-20220314215924425" style="zoom:80%;" /> 

### 使用中断点去观察 stage1 变量

使用中断点去取出相对应 stage1 的 Sha1sum 验证码

<img src="./assets/image-20220314220921338.png" alt="image-20220314220921338" style="zoom:100%;" /> 

### 确认 Sha1sum 验证码的数值

使用下表去对照检查 "中断点去取出相对应 stage1" 和 "Linux Command 去产生 stage1" 的验证码，确认其值为正确的

| 数组位置 |  二进位  | 十进位 | 十六进位 |
| :------: | :------: | :----: | :------: |
|    0     | 10001100 |  140   |    8c    |
|    1     | 10110010 |  178   |    b2    |
|    2     | 00100011 |   35   |    23    |
|    3     | 01111101 |  125   |    7d    |
|    4     | 00000110 |   6    |    06    |
|    5     | 01111001 |  121   |    79    |
|    6     | 11001010 |  202   |    ca    |
|    7     | 10001000 |  136   |    88    |
|    8     | 11011011 |  219   |    db    |
|    9     | 01100100 |  100   |    64    |
|    10    | 01100100 |  100   |    64    |
|    11    | 11101010 |  234   |    ea    |
|    12    | 11000110 |  198   |    c6    |
|    13    | 00001101 |   13   |    0d    |
|    14    | 10101001 |  169   |    a9    |
|    15    | 01100011 |   99   |    63    |
|    16    | 01000101 |   69   |    45    |
|    17    | 01010001 |   81   |    51    |
|    18    | 00111001 |   57   |    39    |
|    19    | 01100100 |  100   |    64    |
