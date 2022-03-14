# Direct Connection in Backend package



## Code



## Testing

> I will describe what I consider about in Unit Test below.

### Considering about Anonymous Function

There is a code below whose name is "Response after Handshake," containing an anonymous function.

The variables in the anonymous function inside the code will take the address of other variables and bring them inside.

I consider about it again and again. It seems correct.

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

## Check

When I use the Linux command to calculate Sha1sum or Use tools on other websites, the result is Hexadecimal.
However, the IDE tools show the result in Decimal.

### Linux command

I am using the Linux command to generate the Sha1shum.

<img src="./assets/image-20220314214316673.png" alt="image-20220314214316673" style="zoom:80%;" /> 

### Website

I am using the tools on the website https://coding.tools/tw/sha1 to calculate the Sha1shum.

<img src="./assets/image-20220314215924425.png" alt="image-20220314215924425" style="zoom:80%;" /> 

### Broken Point

I am using broken point to take a look at the Stage1 variable.

<img src="./assets/image-20220314220921338.png" alt="image-20220314220921338" style="zoom:100%;" /> 

### Comparison

I am using the table below to compare the two results.

One comes from taking broken point, and the other comes from tools on the website.

Thus I  am sure the result is correct.

| Position |  Binary  | Decimal | Hexadecimal |
| :------: | :------: | :-----: | :---------: |
|    0     | 10001100 |   140   |     8c      |
|    1     | 10110010 |   178   |     b2      |
|    2     | 00100011 |   35    |     23      |
|    3     | 01111101 |   125   |     7d      |
|    4     | 00000110 |    6    |     06      |
|    5     | 01111001 |   121   |     79      |
|    6     | 11001010 |   202   |     ca      |
|    7     | 10001000 |   136   |     88      |
|    8     | 11011011 |   219   |     db      |
|    9     | 01100100 |   100   |     64      |
|    10    | 01100100 |   100   |     64      |
|    11    | 11101010 |   234   |     ea      |
|    12    | 11000110 |   198   |     c6      |
|    13    | 00001101 |   13    |     0d      |
|    14    | 10101001 |   169   |     a9      |
|    15    | 01100011 |   99    |     63      |
|    16    | 01000101 |   69    |     45      |
|    17    | 01010001 |   81    |     51      |
|    18    | 00111001 |   57    |     39      |
|    19    | 01100100 |   100   |     64      |
