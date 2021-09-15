package server

import (
	"encoding/json"
	"fmt"
	"github.com/XiaoMi/Gaea/backend"
	"github.com/XiaoMi/Gaea/log"
	"github.com/XiaoMi/Gaea/models"
	"github.com/XiaoMi/Gaea/mysql"
	"github.com/XiaoMi/Gaea/parser/format"
	"github.com/XiaoMi/Gaea/proxy/plan"
	"github.com/XiaoMi/Gaea/util"
	"github.com/stretchr/testify/require"
	"gopkg.in/ini.v1"
	"strings"
	"testing"
)

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> 2 台 Master 4 台 Slave 数据库测试

/*

// 切片 Slice-0 内容如下

                   +----------------+        +----------------+
192.168.122.2:3310 | 从数据库 db1-0   |        | 从数据库 db1-1  | 192.168.122.2:3311
                   +---------\------+        +-------/--------+
                              --\                /---
                                 --\          /--
                                 +------------------------------+
              192.168.122.2:3309 | 主数据库 db1 (数据表 Book_0000) |
                                 +------------------------------+

数据表 Book_0000

| BookID | Isbn          | Title                                            | Author       | Publish | Category           |
| ------ | ------------- | ------------------------------------------------ | ------------ | ------- | ------------------ |
| 2      | 9789869442060 | Water Margin                                     | Shi Nai an   | 1589    | Historical fiction |
| 4      | 9789865975364 | Dream Of The Red Chamber                         | Cao Xueqin   | 1791    | Family Saga        |
| 6      | 9780835124072 | Rulin Waishi                                     | Wu Jingzi    | 1750    | Unofficial History |
| 8      | 9789571447278 | Amazing Tales Second Series                      | Ling Mengchu | 1628    | Perspective        |
| 10     | 9787540251499 | Flowers In The Mirror                            | Li Ruzhen    | 1827    | Fantasy Stories    |
| 12     | 9787101097559 | General Yue Fei                                  | Qian Cai     | 1735    | History            |
| 14     | 9789865700027 | Romance Of Sui And Tang Dynasties                | Chen Ruheng  | 1989    | History            |
| 16     | 9789574927913 | A Collection Of Shi                              | Anonymous    | 1850    | History            |
| 18     | 9787510434341 | Lamp In The Side Street                          | Li Luyuan    | 1790    | Unofficial History |
| 20     | 9789571470047 | Bizarre Happenings Eyewitnessed over Two Decades | Jianren Wu   | 1905    | Unofficial History |
| 22     | 9789861674193 | Officialdom Unmasked                             | Li Baojia    | 1903    | Unofficial History |
| 24     | 9787508067247 | Silent Operas                                    | Li Yu        | 1680    | Social Story       |
| 26     | 9787533948108 | Six Records Of A Floating Life                   | Shen Fu      | 1878    | Autobiography      |
| 28     | 9789571435473 | Nine-Tailed Turtle                               | Lu Can       | 1551    | Mythology          |

// 切片 Slice-1 内容如下

                   +----------------+        +----------------+
192.168.122.2:3313 | 从数据库 db2-0   |        | 从数据库 db2-1   | 192.168.122.2:3314
                   +---------\------+        +-------/--------+
                              --\                /---
                                 --\          /--
                                 +------------------------------+
              192.168.122.2:3312 | 主数据库 db2 (数据表 Book_0001) |
                                 +------------------------------+

数据表 Book_0001

| BookID | Isbn          | Title                              | Author                 | Publish | Category                |
| ------ | ------------- | ---------------------------------- | ---------------------- | ------- | ----------------------- |
| 1      | 9781517191276 | Romance Of The Three Kingdoms      | Luo Guanzhong          | 1522    | Historical fiction      |
| 3      | 9789575709518 | Journey To The West                | Wu Cheng en            | 1592    | Gods And Demons Fiction |
| 5      | 9780804847773 | Jin Ping Mei                       | Lanling Xiaoxiao Sheng | 1610    | Family Life             |
| 7      | 9787101064100 | Amazing Tales First Series         | Ling Mengchu           | 1628    | Perspective             |
| 9      | 9789861273129 | Investiture Of The Gods            | Lu Xixing              | 1605    | Mythology               |
| 11     | 9787508535296 | Stories Old And New                | Feng Menglong          | 1620    | Perspective             |
| 13     | 9789863381037 | The Generals Of The Yang Family    | Qi Zhonglan            | 0       | History                 |
| 15     | 9789575709242 | The Seven Heroes And Five Gallants | Shi Yukun              | 1879    | History                 |
| 17     | 9787533303396 | Dream Of The Green Chamber         | Yuda                   | 1878    | Family Saga             |
| 19     | 9789571447469 | The Travels of Lao Can             | Liu E                  | 1907    | Social Story            |
| 21     | 9787101097580 | A Flower In A Sinful Sea           | Zeng Pu                | 1904    | History                 |
| 23     | 9787805469836 | Tower For The Summer Heat          | Li Yu                  | 1680    | Unofficial History      |
| 25     | 9789573609049 | The Carnal Prayer Mat              | Li Yu                  | 1680    | Social Story            |
| 27     | 9786666141110 | Humble Words Of A Rustic Elder     | Xia Jingqu             | 1787    | Historical fiction      |
| 29     | 9789866318603 | A History Of Floral Treasures      | Chen Sen               | 1849    | Romance                 |

*/

// prepareDb0db1NamespaceManagerForCluster 函式 产生针对 Cluster (db1 db1-0 db1-1) (db2 db2-0 db2-1) 的设定档
func prepareDb0db1NamespaceManagerForCluster() (*Manager, error) {
	// 服务器设定档
	proxyCfg := `
; config type, etcd/file, you can test gaea with file type, you shoud use etcd in production
config_type=file
;file config path, 具体配置放到file_config_path的namespace目录下，该下级目录为固定目录
file_config_path=./etc/file

;coordinator addr
coordinator_addr=http://127.0.0.1:2379
;etcd user config
username=root
password=root

;environ
environ=local
;service name
service_name=gaea_proxy
;gaea_proxy cluster name
cluster_name=gaea_default_cluster

;log config
log_path=./logs
log_level=Notice
log_filename=gaea
log_output=file

;admin addr
admin_addr=0.0.0.0:13307
; basic auth
admin_user=admin
admin_password=admin

;proxy addr
proto_type=tcp4
proxy_addr=0.0.0.0:13306
proxy_charset=utf8
;slow sql time, when execute time is higher than this, log it, unit: ms
slow_sql_time=100
;close session after session timeout, unit: seconds
session_timeout=3600

;stats conf
stats_enabled=true
;stats interval
stats_interval=10

;encrypt key
encrypt_key=1234abcd5678efg*
`

	// 针对 db0 db0-0 db0-1 丛集的设定档
	nsCfg := `
{
  "name": "db0_db1_cluster_namespace",
  "online": true,
  "read_only": false,
  "allowed_dbs": {
    "novel": true
  },
  "slow_sql_time": "1000",
  "black_sql": [
    ""
  ],
  "allowed_ip": null,
  "slices": [
    {
      "name": "slice-0",
      "user_name": "panhong",
      "password": "12345",
      "master": "192.168.122.2:3309",
      "slaves": ["192.168.122.2:3310", "192.168.122.2:3311"],
      "statistic_slaves": null,
      "capacity": 12,
      "max_capacity": 24,
      "idle_timeout": 60
    },
	{
      "name": "slice-1",
      "user_name": "panhong",
      "password": "12345",
      "master": "192.168.122.2:3312",
      "slaves": ["192.168.122.2:3313", "192.168.122.2:3314"],
      "statistic_slaves": null,
      "capacity": 12,
      "max_capacity": 24,
      "idle_timeout": 60
    }
  ],
  "shard_rules": [
	{
      "db": "novel",
      "table": "Book",
      "type": "hash",
      "key": "BookID",
      "locations": [
        1,
        1
      ],
      "slices": [
        "slice-0",
        "slice-1"
      ]
    }
  ],
  "users": [
    {
      "user_name": "panhong",
      "password": "12345",
      "namespace": "db0_db1_cluster_namespace",
      "rw_flag": 2,
      "rw_split": 1,
      "other_property": 0
    }
  ],
  "default_slice": "slice-0",
  "global_sequences": null
}`

	// 把设定档载入到变数

	// 加载 proxy 配置
	var proxy = &models.Proxy{} // 把 cfg map 到 models.Proxy
	cfg, err := ini.Load([]byte(proxyCfg))
	if err != nil {
		return nil, err
	}
	if err = cfg.MapTo(proxy); err != nil { // 把 cfg map 到 models.Proxy
		return nil, err
	}

	// 加载 namespace 配置
	namespaceName := "db0_db1_cluster_namespace"
	namespaceConfig := &models.Namespace{}
	// namespaceConfig Unmarshal 到 nsCfg
	if err := json.Unmarshal([]byte(nsCfg), namespaceConfig); err != nil {
		return nil, err
	}

	// 载入 管理员
	m := NewManager()
	// 初始化 statistics
	statisticManager, err := CreateStatisticManager(proxy, m)
	if err != nil {
		log.Warn("init stats manager failed, %v", err)
		return nil, err
	}
	m.statistics = statisticManager

	// 初始化 namespace
	current, _, _ := m.switchIndex.Get()
	namespaceConfigs := map[string]*models.Namespace{namespaceName: namespaceConfig}
	m.namespaces[current] = CreateNamespaceManager(namespaceConfigs)
	user, err := CreateUserManager(namespaceConfigs)
	if err != nil {
		return nil, err
	}
	m.users[current] = user
	return m, nil
}

// prepareDb0db1PlanSessionExecutorForCluster 函式 产生针对 Cluster (db1 db1-0 db1-1) (db2 db2-0 db2-1) 的 Plan Session
func prepareDb0db1PlanSessionExecutorForCluster() (*SessionExecutor, error) {
	var userName = "panhong"
	var namespaceName = "db0_db1_cluster_namespace"
	var database = "novel"

	m, err := prepareDb0db1NamespaceManagerForCluster()
	if err != nil {
		return nil, err
	}
	executor := newSessionExecutor(m)
	executor.user = userName

	collationID := 33 // "utf8"
	executor.SetCollationID(mysql.CollationID(collationID))
	executor.SetCharset("utf8")
	executor.SetDatabase(database) // set database
	executor.namespace = namespaceName
	return executor, nil
}

// TestDb0db1PlanExecuteInSuite 函式是为了要让以下的单元测试有顺序性
func TestDb0db1PlanExecuteInSuite(t *testing.T) {
	// TestDb0db1PlanExecuteInWrite(t) // 先寫入 29 本小說
	TestDb0db1PlanExecuteInRead(t)
}

// TestDb0db1PlanExecuteInWrite 函式 为向 Cluster (db1 db1-0 db1-1) (db2 db2-0 db2-1) 图书馆数据库写入 29 本小说
// 開始會分庫，但是事先要知道資料會跑到那一個資料庫
func TestDb0db1PlanExecuteInWrite(t *testing.T) {
	// 初始化单元测试程式 (只要注解 Mark TakeOver() 就会使用真的资料库，不然就会跑单元测试)
	backend.MarkTakeOver() // MarkTakeOver 函式一定要放在单元测试最前面，因为可以提早启动一些 DEBUG 除错机制

	// 载入 Session Executor
	se, err := prepareDb0db1PlanSessionExecutorForCluster()
	require.Equal(t, err, nil)
	db, err := se.GetNamespace().GetDefaultPhyDB("novel")
	require.Equal(t, err, nil)
	require.Equal(t, db, "novel") // 检查 SessionExecutor 是否正确载入

	// 开始检查和数据库的沟通
	tests := []struct {
		sql        string
		expect     string
		shardIndex int
	}{
		// 切片规则
		// 会分配到那一个切片，设定档有指定依据 BookId 分配
		// 比如 BookId 为 0 时，会分配到 Slice-0
		// BookId 为 1 时，会分配到 Slice-1

		// 第一本小说 三国演义 (会分配到 Slice-1)
		{
			"INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(1, 9781517191276, 'Romance Of The Three Kingdoms', 'Luo Guanzhong', 1522, 'Historical fiction');",       // 原始的 SQL 字串
			"INSERT INTO `novel`.`Book` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (1,9781517191276,'Romance Of The Three Kingdoms','Luo Guanzhong',1522,'Historical fiction')", // Parser 后的 SQL 字串
			1, // 分配到 Slice-1
		},
		// 第二本小说 水浒传 (会分配到 Slice-0)
		{
			"INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(2, 9789869442060, 'Water Margin', 'Shi Nai an', 1589, 'Historical fiction');",       // 原始的 SQL 字串
			"INSERT INTO `novel`.`Book` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (2,9789869442060,'Water Margin','Shi Nai an',1589,'Historical fiction')", // Parser 后的 SQL 字串
			0, // 分配到 Slice-0
		},
		// 第三本小说 西游记 (会分配到 Slice-1)
		{
			"INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(3, 9789575709518, 'Journey To The West', 'Wu Cheng en', 1592, 'Gods And Demons Fiction');",       // 原始的 SQL 字串
			"INSERT INTO `novel`.`Book` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (3,9789575709518,'Journey To The West','Wu Cheng en',1592,'Gods And Demons Fiction')", // Parser 后的 SQL 字串
			1, // 分配到 Slice-1
		},
		// 第四本小说 红楼梦 (会分配到 Slice-0)
		{
			"INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(4, 9789865975364, 'Dream Of The Red Chamber', 'Cao Xueqin', 1791, 'Family Saga');",       // 原始的 SQL 字串
			"INSERT INTO `novel`.`Book` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (4,9789865975364,'Dream Of The Red Chamber','Cao Xueqin',1791,'Family Saga')", // Parser 后的 SQL 字串
			0, // 分配到 Slice-0
		},

		// 第五本小说 金瓶梅 (会分配到 Slice-1)
		{
			"INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(5, 9780804847773, 'Jin Ping Mei', 'Lanling Xiaoxiao Sheng', 1610, 'Family Life');",       // 原始的 SQL 字串
			"INSERT INTO `novel`.`Book` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (5,9780804847773,'Jin Ping Mei','Lanling Xiaoxiao Sheng',1610,'Family Life')", // Parser 后的 SQL 字串
			1, // 分配到 Slice-1
		},
		// 第六本小说 儒林外史 (会分配到 Slice-0)
		{
			"INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(6, 9780835124072, 'Rulin Waishi', 'Wu Jingzi', 1750, 'Unofficial History');",       // 原始的 SQL 字串
			"INSERT INTO `novel`.`Book` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (6,9780835124072,'Rulin Waishi','Wu Jingzi',1750,'Unofficial History')", // Parser 后的 SQL 字串
			0, // 分配到 Slice-0
		},
		// 第七本小说 初刻拍案惊奇 (会分配到 Slice-1)
		{
			"INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(7, 9787101064100, 'Amazing Tales First Series', 'Ling Mengchu', 1628, 'Perspective');",       // 原始的 SQL 字串
			"INSERT INTO `novel`.`Book` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (7,9787101064100,'Amazing Tales First Series','Ling Mengchu',1628,'Perspective')", // Parser 后的 SQL 字串
			1, // 分配到 Slice-1
		},
		// 第八本小说 二刻拍案惊奇 (会分配到 Slice-0)
		{
			"INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(8, 9789571447278, 'Amazing Tales Second Series', 'Ling Mengchu', 1628, 'Perspective');",       // 原始的 SQL 字串
			"INSERT INTO `novel`.`Book` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (8,9789571447278,'Amazing Tales Second Series','Ling Mengchu',1628,'Perspective')", // Parser 后的 SQL 字串
			0, // 分配到 Slice-0
		},
		// 第九本小说 封神演义 (会分配到 Slice-1)
		{
			"INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(9, 9789861273129, 'Investiture Of The Gods', 'Lu Xixing', 1605, 'Mythology');",       // 原始的 SQL 字串
			"INSERT INTO `novel`.`Book` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (9,9789861273129,'Investiture Of The Gods','Lu Xixing',1605,'Mythology')", // Parser 后的 SQL 字串
			1, // 分配到 Slice-1
		},
		// 第十本小说 镜花缘 (会分配到 Slice-0)
		{
			"INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(10, 9787540251499, 'Flowers In The Mirror', 'Li Ruzhen', 1827, 'Fantasy Stories');",       // 原始的 SQL 字串
			"INSERT INTO `novel`.`Book` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (10,9787540251499,'Flowers In The Mirror','Li Ruzhen',1827,'Fantasy Stories')", // Parser 后的 SQL 字串
			0, // 分配到 Slice-0
		},
		// 第十一本小说 镜花缘 (会分配到 Slice-1)
		{
			"INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(11, 9787508535296, 'Stories Old And New', 'Feng Menglong', 1620, 'Perspective');",       // 原始的 SQL 字串
			"INSERT INTO `novel`.`Book` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (11,9787508535296,'Stories Old And New','Feng Menglong',1620,'Perspective')", // Parser 后的 SQL 字串
			1, // 分配到 Slice-1
		},
		// 第十二本小说 说岳全传 (会分配到 Slice-0)
		{
			"INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(12, 9787101097559, 'General Yue Fei', 'Qian Cai', 1735, 'History');",       // 原始的 SQL 字串
			"INSERT INTO `novel`.`Book` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (12,9787101097559,'General Yue Fei','Qian Cai',1735,'History')", // Parser 后的 SQL 字串
			0, // 分配到 Slice-0
		},
		// 第十三本小说 杨家将 (会分配到 Slice-1)
		{
			"INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(13, 9789863381037, 'The Generals Of The Yang Family', 'Qi Zhonglan', 0, 'History');",       // 原始的 SQL 字串
			"INSERT INTO `novel`.`Book` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (13,9789863381037,'The Generals Of The Yang Family','Qi Zhonglan',0,'History')", // Parser 后的 SQL 字串
			1, // 分配到 Slice-1
		},
		// 第十四本小说 说唐 (会分配到 Slice-0)
		{
			"INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(14, 9789865700027, 'Romance Of Sui And Tang Dynasties', 'Chen Ruheng', 1989, 'History');",       // 原始的 SQL 字串
			"INSERT INTO `novel`.`Book` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (14,9789865700027,'Romance Of Sui And Tang Dynasties','Chen Ruheng',1989,'History')", // Parser 后的 SQL 字串
			0, // 分配到 Slice-0
		},
		// 第十五本小说 七侠五义 (会分配到 Slice-1)
		{
			"INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(15, 9789575709242, 'The Seven Heroes And Five Gallants', 'Shi Yukun', 1879, 'History');",       // 原始的 SQL 字串
			"INSERT INTO `novel`.`Book` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (15,9789575709242,'The Seven Heroes And Five Gallants','Shi Yukun',1879,'History')", // Parser 后的 SQL 字串
			1, // 分配到 Slice-1
		},
		// 第十六本小说 施公案 (会分配到 Slice-0)
		{
			"INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(16, 9789574927913, 'A Collection Of Shi', 'Anonymous', 1850, 'History');",       // 原始的 SQL 字串
			"INSERT INTO `novel`.`Book` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (16,9789574927913,'A Collection Of Shi','Anonymous',1850,'History')", // Parser 后的 SQL 字串
			0, // 分配到 Slice-0
		},
		// 第十七本小说 青楼梦 (会分配到 Slice-1)
		{
			"INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(17, 9787533303396, 'Dream Of The Green Chamber', 'Yuda', 1878, 'Family Saga');",       // 原始的 SQL 字串
			"INSERT INTO `novel`.`Book` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (17,9787533303396,'Dream Of The Green Chamber','Yuda',1878,'Family Saga')", // Parser 后的 SQL 字串
			1, // 分配到 Slice-1
		},
		// 第十八本小说 歧路灯 (会分配到 Slice-0)
		{
			"INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(18, 9787510434341, 'Lamp In The Side Street', 'Li Luyuan', 1790, 'Unofficial History');",       // 原始的 SQL 字串
			"INSERT INTO `novel`.`Book` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (18,9787510434341,'Lamp In The Side Street','Li Luyuan',1790,'Unofficial History')", // Parser 后的 SQL 字串
			0, // 分配到 Slice-0
		},
		// 第十九本小说 老残游记 (会分配到 Slice-1)
		{
			"INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(19, 9789571447469, 'The Travels of Lao Can', 'Liu E', 1907, 'Social Story');",       // 原始的 SQL 字串
			"INSERT INTO `novel`.`Book` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (19,9789571447469,'The Travels of Lao Can','Liu E',1907,'Social Story')", // Parser 后的 SQL 字串
			1, // 分配到 Slice-1
		},
		// 第二十本小说 二十年目睹之怪现状 (会分配到 Slice-0)
		{
			"INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(20, 9789571470047, 'Bizarre Happenings Eyewitnessed over Two Decades', 'Jianren Wu', 1905, 'Unofficial History');",       // 原始的 SQL 字串
			"INSERT INTO `novel`.`Book` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (20,9789571470047,'Bizarre Happenings Eyewitnessed over Two Decades','Jianren Wu',1905,'Unofficial History')", // Parser 后的 SQL 字串
			0, // 分配到 Slice-0
		},
		// 第二十一本小说 孽海花 (会分配到 Slice-1)
		{
			"INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(21, 9787101097580, 'A Flower In A Sinful Sea', 'Zeng Pu', 1904, 'History');",       // 原始的 SQL 字串
			"INSERT INTO `novel`.`Book` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (21,9787101097580,'A Flower In A Sinful Sea','Zeng Pu',1904,'History')", // Parser 后的 SQL 字串
			1, // 分配到 Slice-1
		},
		// 第二十二本小说 官场现形记 (会分配到 Slice-0)
		{
			"INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(22, 9789861674193, 'Officialdom Unmasked', 'Li Baojia', 1903, 'Unofficial History');",       // 原始的 SQL 字串
			"INSERT INTO `novel`.`Book` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (22,9789861674193,'Officialdom Unmasked','Li Baojia',1903,'Unofficial History')", // Parser 后的 SQL 字串
			0, // 分配到 Slice-0
		},
		// 第二十三本小说 觉世名言十二楼 (会分配到 Slice-1)
		{
			"INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(23, 9787805469836, 'Tower For The Summer Heat', 'Li Yu', 1680, 'Unofficial History');",       // 原始的 SQL 字串
			"INSERT INTO `novel`.`Book` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (23,9787805469836,'Tower For The Summer Heat','Li Yu',1680,'Unofficial History')", // Parser 后的 SQL 字串
			1, // 分配到 Slice-1
		},
		// 第二十四本小说 无声戏 (会分配到 Slice-0)
		{
			"INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(24, 9787508067247, 'Silent Operas', 'Li Yu', 1680, 'Social Story');",       // 原始的 SQL 字串
			"INSERT INTO `novel`.`Book` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (24,9787508067247,'Silent Operas','Li Yu',1680,'Social Story')", // Parser 后的 SQL 字串
			0, // 分配到 Slice-0
		},
		// 第二十五本小说 肉蒲团 (会分配到 Slice-1)
		{
			"INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(25, 9789573609049, 'The Carnal Prayer Mat', 'Li Yu', 1680, 'Social Story');",       // 原始的 SQL 字串
			"INSERT INTO `novel`.`Book` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (25,9789573609049,'The Carnal Prayer Mat','Li Yu',1680,'Social Story')", // Parser 后的 SQL 字串
			1, // 分配到 Slice-1
		},
		// 第二十六本小说 浮生六记 (会分配到 Slice-0)
		{
			"INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(26, 9787533948108, 'Six Records Of A Floating Life', 'Shen Fu', 1878, 'Autobiography');",       // 原始的 SQL 字串
			"INSERT INTO `novel`.`Book` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (26,9787533948108,'Six Records Of A Floating Life','Shen Fu',1878,'Autobiography')", // Parser 后的 SQL 字串
			0, // 分配到 Slice-0
		},
		// 第二十七本小说 野叟曝言 (会分配到 Slice-1)
		{
			"INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(27, 9786666141110, 'Humble Words Of A Rustic Elder', 'Xia Jingqu', 1787, 'Historical fiction');",       // 原始的 SQL 字串
			"INSERT INTO `novel`.`Book` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (27,9786666141110,'Humble Words Of A Rustic Elder','Xia Jingqu',1787,'Historical fiction')", // Parser 后的 SQL 字串
			1, // 分配到 Slice-1
		},
		// 第二十八本小说 九尾龟 (会分配到 Slice-0)
		{
			"INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(28, 9789571435473, 'Nine-Tailed Turtle', 'Lu Can', 1551, 'Mythology');",       // 原始的 SQL 字串
			"INSERT INTO `novel`.`Book` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (28,9789571435473,'Nine-Tailed Turtle','Lu Can',1551,'Mythology')", // Parser 后的 SQL 字串
			0, // 分配到 Slice-0
		},
		// 第二十九本小说 品花宝鉴 (会分配到 Slice-1)
		{
			"INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(29, 9789866318603, 'A History Of Floral Treasures', 'Chen Sen', 1849, 'Romance');",       // 原始的 SQL 字串
			"INSERT INTO `novel`.`Book` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (29,9789866318603,'A History Of Floral Treasures','Chen Sen',1849,'Romance')", // Parser 后的 SQL 字串
			1, // 分配到 Slice-1
		},
	}

	// 执行 Sql 字串
	for _, test := range tests {
		// 執行 SQL Parser
		ns := se.GetNamespace()
		stmts, err := se.Parse(test.sql)
		require.Equal(t, err, nil)

		// 检查 Parser 后的 SQL 字串
		var sb strings.Builder
		err = stmts.Restore(format.NewRestoreCtx(format.DefaultRestoreFlags, &sb))
		require.Equal(t, err, nil)
		require.Equal(t, sb.String(), test.expect)

		// 建立 SQL 查寻计划
		rt := ns.GetRouter()
		seq := ns.GetSequences()
		phyDBs := ns.GetPhysicalDBs()
		p, err := plan.BuildPlan(stmts, phyDBs, db, test.sql, rt, seq)
		require.Equal(t, err, nil)

		// 执行 Parser 后的 SQL 指令
		reqCtx := util.NewRequestContext()
		// reqCtx.Set(util.FromSlave, 1) // 在这里设定读取时从 Slave 节点，达到读写分离的效果 (因为是 Insert ，一定会用主数据库，所以此行注解)

		// 执行数据库分库指令
		res, err := p.ExecuteIn(reqCtx, se)

		// 单元测试进行最后检查
		require.Equal(t, err, nil)
		require.Equal(t, p.(*plan.InsertPlan).GetRouteResult().GetShardIndexes(), []int{test.shardIndex})
		require.Equal(t, res.AffectedRows, uint64(0x1))
	}
}

// TestDb0db1PlanExecuteInRead 为向 Cluster (db1 db1-0 db1-1) (db2 db2-0 db2-1) 图书馆数据库查询 29 本小说
func TestDb0db1PlanExecuteInRead(t *testing.T) {
	// 初始化单元测试程式
	backend.MarkTakeOver() // MarkTakeOver 函式一定要放在单元测试最前面，因为可以提早启动一些 DEBUG 除错机制

	// 载入 Session Executor
	se, err := prepareDb0db1PlanSessionExecutorForCluster()
	require.Equal(t, err, nil)
	db, err := se.GetNamespace().GetDefaultPhyDB("novel")
	require.Equal(t, err, nil)
	require.Equal(t, db, "novel") // 检查 SessionExecutor 是否正确载入

	// 开始检查和数据库的沟通
	tests := []struct {
		sql    string
		expect string
	}{
		// 第一本小说 三国演义
		{
			"INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(1, 9781517191276, 'Romance Of The Three Kingdoms', 'Luo Guanzhong', 1522, 'Historical fiction');",       // 原始的 SQL 字串
			"INSERT INTO `novel`.`Book` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (1,9781517191276,'Romance Of The Three Kingdoms','Luo Guanzhong',1522,'Historical fiction')", // Parser 后的 SQL 字串
		},
		{ // 测试二，查询数据库资料
			// 有二组切片组成时，需要加 Order By，回传的数据才会有顺序性，每次回传的数据一致时，单元测试才会正常
			"SELECT * FROM novel.Book ORDER BY BookID",       // 原始的 SQL 字串
			"SELECT * FROM `novel`.`Book` ORDER BY `BookID`", // 期望 Parser 后的 SQL 字串
		},
	}

	// 执行 Sql 字串
	for _, test := range tests {
		// 執行 SQL Parser
		ns := se.GetNamespace()
		stmts, err := se.Parse(test.sql)
		require.Equal(t, err, nil)

		// 检查 Parser 后的 SQL 字串
		var sb strings.Builder
		err = stmts.Restore(format.NewRestoreCtx(format.DefaultRestoreFlags, &sb))
		require.Equal(t, err, nil)
		require.Equal(t, sb.String(), test.expect)

		// 建立 SQL 查寻计划
		rt := ns.GetRouter()
		seq := ns.GetSequences()
		phyDBs := ns.GetPhysicalDBs()
		p, err := plan.BuildPlan(stmts, phyDBs, db, test.sql, rt, seq)
		require.Equal(t, err, nil)

		// 执行 Parser 后的 SQL 指令
		reqCtx := util.NewRequestContext()
		reqCtx.Set(util.FromSlave, 1) // 在这里设定读取时从 Slave 节点，达到读写分离的效果

		// 执行数据库分库指令
		res, err := p.ExecuteIn(reqCtx, se)
		require.Equal(t, err, nil)
		fmt.Println(res)
	}
}
