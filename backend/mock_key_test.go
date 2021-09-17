package backend

import (
	"github.com/stretchr/testify/require"
	"testing"
)

// TestMockNovelKey 函式 🧚 测试 SQL 字串所对应到的 Key
func TestMockNovelKey(t *testing.T) {
	/*
		所对应各切片 SQL 执行字串 以及 切片相关资讯
		数据库名称: novel
		模拟数据库的网路位置: 192.168.122.2:3313
		数据库执行字串: SELECT *,`BookID` FROM `novel`.`Book_0001` ORDER BY `BookID`
		数据库执行时所对应的 Key: 1260331735
	*/
	tests := []struct {
		addr     string // 网路地址
		user     string // 用户名称
		password string // 用户密码
		db       string // 数据库名称
		sql      string // SQL 执行字串
		key      int    // 所对应的 key 值
	}{
		// >>>>> >>>>> >>>>> >>>>> >>>>> 对第二个丛集进行查询
		// 192.168.122.2:3309 192.168.122.2:3310 192.168.122.2:3311 使用第二个切片 Slice-0
		// 数据库名称为 Book_0000
		{
			// 第一台丛集主伺服器
			"192.168.122.2:3309", // 网路地址
			"panhong",            // 用户名称
			"12345",              // 用户密码
			"novel",              // 数据库名称
			"SELECT *,`BookID` FROM `novel`.`Book_0000` ORDER BY `BookID`", // SQL 执行字串
			3717314451, // 所对应的 key 值
		},
		{
			// 第一台丛集从伺服器
			"192.168.122.2:3310", // 网路地址
			"panhong",            // 用户名称
			"12345",              // 用户密码
			"novel",              // 数据库名称
			"SELECT *,`BookID` FROM `novel`.`Book_0000` ORDER BY `BookID`", // SQL 执行字串
			1196547673, // 所对应的 key 值
		},
		{
			// 第二台丛集从伺服器
			"192.168.122.2:3311", // 网路地址
			"panhong",            // 用户名称
			"12345",              // 用户密码
			"novel",              // 数据库名称
			"SELECT *,`BookID` FROM `novel`.`Book_0000` ORDER BY `BookID`", // SQL 执行字串
			4270781616, // 所对应的 key 值
		},
		// >>>>> >>>>> >>>>> >>>>> >>>>> 对第二个丛集进行资料新增
		{
			// 新增第二本小说 水浒传 至第一台丛集主伺服器
			"192.168.122.2:3309", // 网路地址
			"panhong",            // 用户名称
			"12345",              // 用户密码
			"novel",              // 数据库名称
			"INSERT INTO `novel`.`Book_0000` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (2,9789869442060,'Water Margin','Shi Nai an',1589,'Historical fiction')", // SQL 执行字串
			618120042, // 所对应的 key 值
		},
		// >>>>> >>>>> >>>>> >>>>> >>>>> 对第三个丛集进行查询
		// 192.168.122.2:3312 192.168.122.2:3313 192.168.122.2:3314 使用第二个切片 Slice-1
		// 数据库名称为 Book_0001
		{
			// 第一台丛集主伺服器
			"192.168.122.2:3312", // 网路地址
			"panhong",            // 用户名称
			"12345",              // 用户密码
			"novel",              // 数据库名称
			"SELECT *,`BookID` FROM `novel`.`Book_0000` ORDER BY `BookID`", // SQL 执行字串
			2864051087, // 所对应的 key 值
		},
		{
			// 第一台丛集从伺服器
			"192.168.122.2:3313", // 网路地址
			"panhong",            // 用户名称
			"12345",              // 用户密码
			"novel",              // 数据库名称
			"SELECT *,`BookID` FROM `novel`.`Book_0000` ORDER BY `BookID`", // SQL 执行字串
			575021710, // 所对应的 key 值
		},
		{
			// 第二台丛集从伺服器
			"192.168.122.2:3314", // 网路地址
			"panhong",            // 用户名称
			"12345",              // 用户密码
			"novel",              // 数据库名称
			"SELECT *,`BookID` FROM `novel`.`Book_0001` ORDER BY `BookID`", // SQL 执行字串
			1401931444, // 所对应的 key 值
		},
		// >>>>> >>>>> >>>>> >>>>> >>>>> 对第三个丛集进行资料新增
		{
			// 新增第一本小说 三国演义 至第一台丛集主伺服器
			"192.168.122.2:3309", // 网路地址
			"panhong",            // 用户名称
			"12345",              // 用户密码
			"novel",              // 数据库名称
			"INSERT INTO `novel`.`Book` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (1,9781517191276,'Romance Of The Three Kingdoms','Luo Guanzhong',1522,'Historical fiction')", // SQL 执行字串
			3261299321, // 所对应的 key 值
		},
	}
	for i := 0; i < len(tests); i++ {
		// 把资料转换成 key
		dc := DirectConnection{ // 组成 DC 资料
			addr:     tests[i].addr,
			user:     tests[i].user,
			password: tests[i].password,
			db:       tests[i].db,
		}
		// 进行最后比对
		key := dc.MakeMockKey(tests[i].sql)      // 组成 key 资料
		require.Equal(t, int(key), tests[i].key) // 要符合 1260331735
	}
}
