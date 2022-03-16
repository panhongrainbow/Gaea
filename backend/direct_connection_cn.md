# backend 包的直连 direct_connection



## 代码说明

### 第一步 初始交握，传送讯息方向为 MariaDB 至 Gaea

参考官方文档 https://mariadb.com/kb/en/connection/ ，有以下内容

<img src="/home/panhong/go/src/github.com/panhongrainbow/note/typora-user-images/image-20220315221559157.png" alt="image-20220315221559157" style="zoom:100%;" /> 

根据官方文档，使用范例说明

| 内容                            | 演示范例                                                     |
| ------------------------------- | ------------------------------------------------------------ |
| int<1> protocol version         | 协定 Protocol 版本为<br />10                                 |
| string<NUL> server version      | 数据库的版本号 version 为<br /><br />[]uint8{<br />53, 46, 53, 46, 53,<br />45, 49, 48, 46, 53,<br />46, 49, 50, 45, 77,<br />97, 114, 105, 97, 68,<br />66, 45, 108, 111, 103<br />}<br /><br />对照 ASCII 表为<br />5.5.5-10.5.12-MariaDB-log |
| int<4> connection id            | 连接编号为<br /><br />[]uint8{16, 0, 0, 0}<br />先反向排列为 []uint8{0, 0, 0, 16}<br /><br />最后求得的连接编号为 uint32(16) |
| string<8> scramble 1st part     | 第一部份的 Scramble，Scramble 总共需要组成 20 bytes，<br />第一个部份共 8 bytes，其值为 []uint8{81, 64, 43, 85, 76, 90, 97, 91} |
| string<1> reserved byte         | 数值为 0                                                     |
| int<2> server capabilities      | 第一部份的功能标志 capability，数值为 []uint8{254, 247}      |
| int<1> server default collation | 數據庫編碼 charset 为 33，经<br />以下文文档查询 [character-sets-and-collations](https://mariadb.com/kb/en/supported-character-sets-and-collations/)<br />或者是 命令 SHOW CHARACTER SET LIKE 'utf8'; 查询，<br />charset 的数值为 utf8_general_ci |
| int<2> status flags             | 服务器状态为 []uint8{2, 0}<br />进行反向排列为[]uint8{0, 2}，再转成二进制为 []uint16{2}<br />对照 Gaea/mysql/constants.go 后，得知目前服务器的状况为<br />Autocommit (ServerStatusAutocommit) |
| int<2> server capabilities      | 延伸的功能标志 capability，数值为 uint16[255, 129]           |

先对 功能标志 capability 进行计算

```
先把所有的功能标志 capability 的数据收集起来，包含延伸部份

数值分别为 []uint8{254, 247, 255, 129}
并反向排列
数值分别为 []uint8{129, 255, 247, 254}
全部 十进制 转成 二进制，为 []uint8{10000001, 11111111, 11110111, 11111110} (转成十进制数值为 2181036030)

再用文档 https://mariadb.com/kb/en/connection/ 进行对照
比如，功能标志 capability 的第一个值为 0，意思为 CLIENT_MYSQL 值为 0，代表是由服务器发出的讯息
```

接续上表

| 项目 | 内容                                                         |
| ---- | ------------------------------------------------------------ |
| 公式 | if (server_capabilities & PLUGIN_AUTH)<br/>        int<1> plugin data length <br/>    else<br/>        int<1> 0x00 |
| 范例 | 跳过 1 个 byte                                               |

接续上表

| 项目 | 内容             |
| ---- | ---------------- |
| 公式 | string<6> filler |
| 范例 | 跳过 6 个 bytes  |

接续上表

| 项目 | 内容                                                         |
| ---- | ------------------------------------------------------------ |
| 公式 | if (server_capabilities & CLIENT_MYSQL)<br/>        string<4> filler <br/>    else<br/>        int<4> server capabilities 3rd part .<br />        MariaDB specific flags /* MariaDB 10.2 or later */ |
| 范例 | 跳过 4 个 bytes                                              |

接续上表

| 项目 | 内容                                                         |
| ---- | ------------------------------------------------------------ |
| 公式 | if (server_capabilities & CLIENT_SECURE_CONNECTION)<br/>        string<n> scramble 2nd part . Length = max(12, plugin data length - 9)<br/>        string<1> reserved byte |
| 范例 | scramble 一共要 20 个 bytes，第一部份共 8 bytes，所以第二部份共有 20 - 8 = 12 bytes，该值为 []uint8{34, 53, 36, 85, 93, 86, 117, 105, 49, 87, 65, 125} |

接续上表

| 项目 | 内容                                                         |
| ---- | ------------------------------------------------------------ |
| 公式 | if (server_capabilities & PLUGIN_AUTH)<br/>        string<NUL> authentication plugin name |
| 范例 | 之后的资料都不使用                                           |

合拼所有 Scramble 的资料

```
第一部份 Scramble 为 []uint8{81, 64, 43, 85, 76, 90, 97, 91}
第二部份 Scramble 为 []uint8{34, 53, 36, 85, 93, 86, 117, 105, 49, 87, 65, 125}

两部份 Scramble 合拼后为 []uint8{81, 64, 43, 85, 76, 90, 97, 91, 34, 53, 36, 85, 93, 86, 117, 105, 49, 87, 65, 125}
```



### 第二步 回应交握，传送讯息方向为 Gaea 至 MariaDB



### 第三步 交握完成，传送讯息方向为 MariaDB 至 Gaea



## 测试说明

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
