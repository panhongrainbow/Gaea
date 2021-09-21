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
		{
			// 新增第四本小说 红楼梦 至第一台丛集主伺服器
			"192.168.122.2:3309", // 网路地址
			"panhong",            // 用户名称
			"12345",              // 用户密码
			"novel",              // 数据库名称
			"INSERT INTO `novel`.`Book_0000` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (4,9789865975364,'Dream Of The Red Chamber','Cao Xueqin',1791,'Family Saga')", // SQL 执行字串
			4273731942, // 所对应的 key 值
		},
		{
			// 新增第六本小说 儒林外史 至第一台丛集主伺服器
			"192.168.122.2:3309", // 网路地址
			"panhong",            // 用户名称
			"12345",              // 用户密码
			"novel",              // 数据库名称
			"INSERT INTO `novel`.`Book_0000` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (6,9780835124072,'Rulin Waishi','Wu Jingzi',1750,'Unofficial History')", // SQL 执行字串
			1926088204, // 所对应的 key 值
		},
		{
			// 新增第八本小说 二刻拍案惊奇 至第一台丛集主伺服器
			"192.168.122.2:3309", // 网路地址
			"panhong",            // 用户名称
			"12345",              // 用户密码
			"novel",              // 数据库名称
			"INSERT INTO `novel`.`Book_0000` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (8,9789571447278,'Amazing Tales Second Series','Ling Mengchu',1628,'Perspective')", // SQL 执行字串
			1708424148, // 所对应的 key 值
		},
		{
			// 新增第十本小说 镜花缘 至第一台丛集主伺服器
			"192.168.122.2:3309", // 网路地址
			"panhong",            // 用户名称
			"12345",              // 用户密码
			"novel",              // 数据库名称
			"INSERT INTO `novel`.`Book_0000` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (10,9787540251499,'Flowers In The Mirror','Li Ruzhen',1827,'Fantasy Stories')", // SQL 执行字串
			3303343655, // 所对应的 key 值
		},
		{
			// 新增第十二本小说 说岳全传 至第一台丛集主伺服器
			"192.168.122.2:3309", // 网路地址
			"panhong",            // 用户名称
			"12345",              // 用户密码
			"novel",              // 数据库名称
			"INSERT INTO `novel`.`Book_0000` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (12,9787101097559,'General Yue Fei','Qian Cai',1735,'History')", // SQL 执行字串
			600352469, // 所对应的 key 值
		},
		{
			// 新增第十四本小说 说唐 至第一台丛集主伺服器
			"192.168.122.2:3309", // 网路地址
			"panhong",            // 用户名称
			"12345",              // 用户密码
			"novel",              // 数据库名称
			"INSERT INTO `novel`.`Book_0000` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (14,9789865700027,'Romance Of Sui And Tang Dynasties','Chen Ruheng',1989,'History')", // SQL 执行字串
			1226676578, // 所对应的 key 值
		},
		{
			// 新增第十六本小说 施公案 至第一台丛集主伺服器
			"192.168.122.2:3309", // 网路地址
			"panhong",            // 用户名称
			"12345",              // 用户密码
			"novel",              // 数据库名称
			"INSERT INTO `novel`.`Book_0000` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (16,9789574927913,'A Collection Of Shi','Anonymous',1850,'History')", // SQL 执行字串
			3585696861, // 所对应的 key 值
		},
		{
			// 新增第十八本小说 歧路灯 至第一台丛集主伺服器
			"192.168.122.2:3309", // 网路地址
			"panhong",            // 用户名称
			"12345",              // 用户密码
			"novel",              // 数据库名称
			"INSERT INTO `novel`.`Book_0000` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (18,9787510434341,'Lamp In The Side Street','Li Luyuan',1790,'Unofficial History')", // SQL 执行字串
			1792929480, // 所对应的 key 值
		},
		{
			// 新增第二十本小说 二十年目睹之怪现状 至第一台丛集主伺服器
			"192.168.122.2:3309", // 网路地址
			"panhong",            // 用户名称
			"12345",              // 用户密码
			"novel",              // 数据库名称
			"INSERT INTO `novel`.`Book_0000` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (20,9789571470047,'Bizarre Happenings Eyewitnessed over Two Decades','Jianren Wu',1905,'Unofficial History')", // SQL 执行字串
			1187323765, // 所对应的 key 值
		},
		{
			// 新增第二十二本小说 官场现形记 至第一台丛集主伺服器
			"192.168.122.2:3309", // 网路地址
			"panhong",            // 用户名称
			"12345",              // 用户密码
			"novel",              // 数据库名称
			"INSERT INTO `novel`.`Book_0000` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (22,9789861674193,'Officialdom Unmasked','Li Baojia',1903,'Unofficial History')", // SQL 执行字串
			2032678570, // 所对应的 key 值
		},
		{
			// 新增第二十四本小说 无声戏 至第一台丛集主伺服器
			"192.168.122.2:3309", // 网路地址
			"panhong",            // 用户名称
			"12345",              // 用户密码
			"novel",              // 数据库名称
			"INSERT INTO `novel`.`Book_0000` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (24,9787508067247,'Silent Operas','Li Yu',1680,'Social Story')", // SQL 执行字串
			2457093340, // 所对应的 key 值
		},
		{
			// 新增第二十六本小说 浮生六记 至第一台丛集主伺服器
			"192.168.122.2:3309", // 网路地址
			"panhong",            // 用户名称
			"12345",              // 用户密码
			"novel",              // 数据库名称
			"INSERT INTO `novel`.`Book_0000` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (26,9787533948108,'Six Records Of A Floating Life','Shen Fu',1878,'Autobiography')", // SQL 执行字串
			4020693348, // 所对应的 key 值
		},
		{
			// 新增第二十八本小说 九尾龟 至第一台丛集主伺服器
			"192.168.122.2:3309", // 网路地址
			"panhong",            // 用户名称
			"12345",              // 用户密码
			"novel",              // 数据库名称
			"INSERT INTO `novel`.`Book_0000` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (28,9789571435473,'Nine-Tailed Turtle','Lu Can',1551,'Mythology')", // SQL 执行字串
			1776512190, // 所对应的 key 值
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
			"SELECT *,`BookID` FROM `novel`.`Book_0001` ORDER BY `BookID`", // SQL 执行字串
			2403537350, // 所对应的 key 值
		},
		{
			// 第一台丛集从伺服器
			"192.168.122.2:3313", // 网路地址
			"panhong",            // 用户名称
			"12345",              // 用户密码
			"novel",              // 数据库名称
			"SELECT *,`BookID` FROM `novel`.`Book_0001` ORDER BY `BookID`", // SQL 执行字串
			1260331735, // 所对应的 key 值
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
			"192.168.122.2:3312", // 网路地址
			"panhong",            // 用户名称
			"12345",              // 用户密码
			"novel",              // 数据库名称
			"INSERT INTO `novel`.`Book_0001` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (1,9781517191276,'Romance Of The Three Kingdoms','Luo Guanzhong',1522,'Historical fiction')", // SQL 执行字串
			1389454267, // 所对应的 key 值
		},
		{
			// 新增第三本小说 西游记 至第一台丛集主伺服器
			"192.168.122.2:3312", // 网路地址
			"panhong",            // 用户名称
			"12345",              // 用户密码
			"novel",              // 数据库名称
			"INSERT INTO `novel`.`Book_0001` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (3,9789575709518,'Journey To The West','Wu Cheng en',1592,'Gods And Demons Fiction')", // SQL 执行字串
			514659115, // 所对应的 key 值
		},
		{
			// 新增第五本小说 金瓶梅 至第一台丛集主伺服器
			"192.168.122.2:3312", // 网路地址
			"panhong",            // 用户名称
			"12345",              // 用户密码
			"novel",              // 数据库名称
			"INSERT INTO `novel`.`Book_0001` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (5,9780804847773,'Jin Ping Mei','Lanling Xiaoxiao Sheng',1610,'Family Life')", // SQL 执行字串
			4076192191, // 所对应的 key 值
		},
		{
			// 新增第七本小说 初刻拍案惊奇 至第一台丛集主伺服器
			"192.168.122.2:3312", // 网路地址
			"panhong",            // 用户名称
			"12345",              // 用户密码
			"novel",              // 数据库名称
			"INSERT INTO `novel`.`Book_0001` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (7,9787101064100,'Amazing Tales First Series','Ling Mengchu',1628,'Perspective')", // SQL 执行字串
			1572904758, // 所对应的 key 值
		},
		{
			// 新增第九本小说 封神演义 至第一台丛集主伺服器
			"192.168.122.2:3312", // 网路地址
			"panhong",            // 用户名称
			"12345",              // 用户密码
			"novel",              // 数据库名称
			"INSERT INTO `novel`.`Book_0001` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (9,9789861273129,'Investiture Of The Gods','Lu Xixing',1605,'Mythology')", // SQL 执行字串
			3188314210, // 所对应的 key 值
		},
		{
			// 新增第十一本小说 喻世明言 至第一台丛集主伺服器
			"192.168.122.2:3312", // 网路地址
			"panhong",            // 用户名称
			"12345",              // 用户密码
			"novel",              // 数据库名称
			"INSERT INTO `novel`.`Book_0001` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (11,9787508535296,'Stories Old And New','Feng Menglong',1620,'Perspective')", // SQL 执行字串
			3599615497, // 所对应的 key 值
		},
		{
			// 新增第十三本小说 杨家将 至第一台丛集主伺服器
			"192.168.122.2:3312", // 网路地址
			"panhong",            // 用户名称
			"12345",              // 用户密码
			"novel",              // 数据库名称
			"INSERT INTO `novel`.`Book_0001` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (13,9789863381037,'The Generals Of The Yang Family','Qi Zhonglan',0,'History')", // SQL 执行字串
			709958148, // 所对应的 key 值
		},
		{
			// 新增第十五本小说 七侠五义 至第一台丛集主伺服器
			"192.168.122.2:3312", // 网路地址
			"panhong",            // 用户名称
			"12345",              // 用户密码
			"novel",              // 数据库名称
			"INSERT INTO `novel`.`Book_0001` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (15,9789575709242,'The Seven Heroes And Five Gallants','Shi Yukun',1879,'History')", // SQL 执行字串
			56203336, // 所对应的 key 值
		},
		{
			// 新增第十七本小说 第十七 至第一台丛集主伺服器
			"192.168.122.2:3312", // 网路地址
			"panhong",            // 用户名称
			"12345",              // 用户密码
			"novel",              // 数据库名称
			"INSERT INTO `novel`.`Book_0001` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (17,9787533303396,'Dream Of The Green Chamber','Yuda',1878,'Family Saga')", // SQL 执行字串
			3821388015, // 所对应的 key 值
		},
		{
			// 新增第十九本小说 老残游记 至第一台丛集主伺服器
			"192.168.122.2:3312", // 网路地址
			"panhong",            // 用户名称
			"12345",              // 用户密码
			"novel",              // 数据库名称
			"INSERT INTO `novel`.`Book_0001` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (19,9789571447469,'The Travels of Lao Can','Liu E',1907,'Social Story')", // SQL 执行字串
			398747927, // 所对应的 key 值
		},
		{
			// 新增第二十一本小说 孽海花 至第一台丛集主伺服器
			"192.168.122.2:3312", // 网路地址
			"panhong",            // 用户名称
			"12345",              // 用户密码
			"novel",              // 数据库名称
			"INSERT INTO `novel`.`Book_0001` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (21,9787101097580,'A Flower In A Sinful Sea','Zeng Pu',1904,'History')", // SQL 执行字串
			1498815330, // 所对应的 key 值
		},
		{
			// 新增第二十三本小说 觉世名言十二楼 至第一台丛集主伺服器
			"192.168.122.2:3312", // 网路地址
			"panhong",            // 用户名称
			"12345",              // 用户密码
			"novel",              // 数据库名称
			"INSERT INTO `novel`.`Book_0001` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (23,9787805469836,'Tower For The Summer Heat','Li Yu',1680,'Unofficial History')", // SQL 执行字串
			2614046017, // 所对应的 key 值
		},
		{
			// 新增第二十五本小说 肉蒲团 至第一台丛集主伺服器
			"192.168.122.2:3312", // 网路地址
			"panhong",            // 用户名称
			"12345",              // 用户密码
			"novel",              // 数据库名称
			"INSERT INTO `novel`.`Book_0001` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (25,9789573609049,'The Carnal Prayer Mat','Li Yu',1680,'Social Story')", // SQL 执行字串
			238477972, // 所对应的 key 值
		},
		{
			// 新增第二十七本小说 野叟曝言 至第一台丛集主伺服器
			"192.168.122.2:3312", // 网路地址
			"panhong",            // 用户名称
			"12345",              // 用户密码
			"novel",              // 数据库名称
			"INSERT INTO `novel`.`Book_0001` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (27,9786666141110,'Humble Words Of A Rustic Elder','Xia Jingqu',1787,'Historical fiction')", // SQL 执行字串
			2745523730, // 所对应的 key 值
		},
		{
			// 新增第二十九本小说 品花宝鉴 至第一台丛集主伺服器
			"192.168.122.2:3312", // 网路地址
			"panhong",            // 用户名称
			"12345",              // 用户密码
			"novel",              // 数据库名称
			"INSERT INTO `novel`.`Book_0001` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (29,9789866318603,'A History Of Floral Treasures','Chen Sen',1849,'Romance')", // SQL 执行字串
			424563096, // 所对应的 key 值
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
