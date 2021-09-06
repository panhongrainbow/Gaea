package server

import (
	"fmt"
	"github.com/XiaoMi/Gaea/util"

	"github.com/XiaoMi/Gaea/backend"
	"github.com/XiaoMi/Gaea/models"
	"github.com/XiaoMi/Gaea/mysql"
	"github.com/XiaoMi/Gaea/parser"
	"github.com/XiaoMi/Gaea/parser/format"
	"github.com/XiaoMi/Gaea/proxy/plan"
	"github.com/XiaoMi/Gaea/proxy/router"
	"github.com/XiaoMi/Gaea/proxy/sequence"

	// 先忽略 stats 包，因为这个测试会跟其他测试发生冲突
	// "github.com/XiaoMi/Gaea/stats"

	"strings"
	"testing"

	"github.com/XiaoMi/Gaea/util/cache"
	"github.com/stretchr/testify/require"
	"gopkg.in/ini.v1"
)

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> 单台独立数据库测试

// TestB1 设定档 key value 细部测试
func TestB1(t *testing.T) {
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

	tests := []struct {
		key   string
		value string
	}{
		{"cluster_name", "gaea_default_cluster"},
		{"slow_sql_time", "100"},
		// 之后再扩充
	}

	cfg, err := ini.Load([]byte(proxyCfg))
	require.Equal(t, err, nil)

	for _, test := range tests {
		cluster, err := cfg.Sections()[0].GetKey(test.key)
		require.Equal(t, cluster.String(), test.value)
		require.Equal(t, err, nil)
	}
}

// TestB2 内含把 29 本小说写入数据库的程式码，这时还没有进行数据库读写分离
func TestB2(t *testing.T) {
	// 组成管理员的 NameSpace (会写入变数 SessionExecutor -> namespaces[0] -> namespaces["key"])
	managerNamespace := &Namespace{} // 管理员的 NameSpace
	defer managerNamespace.Close(false)
	managerNamespace.openGeneralLog = false             // 记录sql查询的访问日志，说明 https://github.com/XiaoMi/Gaea/issues/109
	managerNamespace.name = "db0_namespace"             // namespace 为划分工作业务的最基本单位，一个 namespace 可以有多个使用者
	managerNamespace.allowedDBs = make(map[string]bool) // 数据库列表
	managerNamespace.allowedDBs[strings.TrimSpace("novel")] = true
	defaultPhyDBs := make(map[string]string, 0) // 预设数据库列表
	// defaultPhyDBs[strings.TrimSpace(db)] = strings.TrimSpace(phyDB) // 再指定
	managerNamespace.defaultPhyDBs, _ = parseDefaultPhyDB(defaultPhyDBs, managerNamespace.allowedDBs)
	managerNamespace.slowSQLTime = 1000                      // 慢sql阈值，单位: 毫秒
	managerNamespace.sqls = make(map[string]string, 50)      // 有 29 本小说，先暂定 50 好了
	managerNamespace.sqls = parseBlackSqls([]string{})       // 黑名单sql
	managerNamespace.allowips, _ = parseAllowIps([]string{}) // 白名单IP

	// >>>>> 组织主从物理实例
	var dbSlices []*models.Slice // 一主多从的物理实例，slice里map的具体字段可参照slice配置
	dbSlice := models.Slice{}
	dbSlice.Name = "slice-0"              // 分片名称，自动、有序生成
	dbSlice.UserName = "panhong"          // 连接后端mysql所需要的用户名称
	dbSlice.Password = "12345"            // 连接后端mysql所需要的用户密码
	dbSlice.Master = "192.168.122.2:3306" // 主实例地址 (db0 192.168.1.2:3350)
	dbSlice.Slaves = []string{}           // 从实例地址列表
	dbSlice.StatisticSlaves = nil         // 统计型从实例地址列表
	dbSlice.Capacity = 12                 // gaea_proxy与每个实例的连接池大小
	dbSlice.MaxCapacity = 24              // gaea_proxy与每个实例的连接池最大大小
	dbSlice.IdleTimeout = 60              // gaea_proxy与后端mysql空闲连接存活时间，单位:秒
	dbSlices = append(dbSlices, &dbSlice)

	// 管理员的 NameSpace 載入 组织主从物理实例 (会写入变数 SessionExecutor -> namespaces[0] -> slices["key"])
	managerNamespace.slices, _ = parseSlices(dbSlices, "utf8mb4", managerNamespace.defaultCollationID) // 一主多从的物理实例，slice里map的具体字段可参照slice配置

	// >>>>> 组织 NS 用户 (会写入变数 SessionExecutor -> namespaces[0] -> userProperties["key"])
	var nsUsers []*models.User // 应用端连接gaea所需要的用户配置，具体字段可参照users配置
	nsUser := models.User{}
	nsUser.UserName = "panhong"        // 用户名
	nsUser.Password = "12345"          // 用户密码
	nsUser.Namespace = "db0_namespace" // 对应的命名空间
	nsUser.RWFlag = 2                  // 读写标识
	nsUser.RWSplit = 1                 // 是否读写分离
	nsUser.OtherProperty = 0           // 其他属性，目前用来标识是否走统计从实例
	nsUsers = append(nsUsers, &nsUser)

	// 管理员的 NameSpace 載入 NS 用户 (会写入变数 SessionExecutor -> namespaces[0] -> userProperties["key"])
	managerNamespace.userProperties = make(map[string]*UserProperty, 2)
	for _, user := range nsUsers {
		up := &UserProperty{RWFlag: user.RWFlag, RWSplit: user.RWSplit, OtherProperty: user.OtherProperty}
		managerNamespace.userProperties[user.UserName] = up
	}

	// >>>>> 建立 全局唯一序列号 (会写入变数 SessionExecutor -> namespaces[0] -> sequence["key"])
	sequences := sequence.NewSequenceManager() // 生成全局唯一序列号的配置, 具体字段可参考全局序列号配置
	for _, v := range []*models.GlobalSequence{} {
		globalSequenceSlice, _ := managerNamespace.slices[v.SliceName]
		seqName := strings.ToUpper(v.DB) + "." + strings.ToUpper(v.Table)
		seq := sequence.NewMySQLSequence(globalSequenceSlice, seqName, v.PKName)
		sequences.SetSequence(v.DB, v.Table, seq)
	}

	// 管理员的 NameSpace 載入 全局唯一序列号 (会写入变数 SessionExecutor -> namespaces[0] -> sequence["key"])
	managerNamespace.sequences = sequences

	// 管理员的 NameSpace 載入 预设选项
	/*
		补齐所有的 managerNamespace 变数，最后先组成 managerNamespace (物品)
		namespaceManager(map) 的人 指定 db0_cluster_namespace 字串 map 到前面的 managerNamespace (物品)
		合拼到服務器管理員 serverManager := new(Manager)
		再合拼到 Session 执行者 executor.manager = serverManager
	*/
	managerNamespace.defaultCharset, managerNamespace.defaultCollationID, _ = parseCharset("", "") // 用于指定数据集如何排序，以及字符串的比对规则 & 用于指定数据集如何排序，以及字符串的比对规则
	managerNamespace.maxSqlExecuteTime = defaultMaxSqlExecuteTime                                  // sql最大执行时间，大于该时间，进行熔断
	managerNamespace.maxSqlResultSize = defaultMaxSqlResultSize                                    // 限制单分片返回结果集大小不超过max_select_rows

	// 管理员的 NameSpace 載入 延伸选择
	managerNamespace.slowSQLCache = cache.NewLRUCache(defaultSQLCacheCapacity)
	managerNamespace.errorSQLCache = cache.NewLRUCache(defaultSQLCacheCapacity)
	managerNamespace.backendSlowSQLCache = cache.NewLRUCache(defaultSQLCacheCapacity)
	managerNamespace.backendErrorSQLCache = cache.NewLRUCache(defaultSQLCacheCapacity)
	managerNamespace.planCache = cache.NewLRUCache(defaultPlanCacheCapacity)

	// 建立路由
	modelsNameSpace := new(models.Namespace)
	modelsNameSpace.Slices = []*models.Slice{} // 一主多从的物理实例，slice里map的具体字段可参照slice配置
	modelsNameSpace.Slices = append(modelsNameSpace.Slices, &dbSlice)
	modelsNameSpace.DefaultSlice = "slice-0" // 预设分片名称
	modelsNameSpace.ShardRules = nil         // 分库、分表、特殊表的配置内容，具体字段可参照shard配置 (载入设定档)
	managerNamespace.router, _ = router.NewRouter(modelsNameSpace)

	// 初始化管理员 Manager
	serverManager := new(Manager)             // 服务器管理员
	current := 0                              // 切换标签
	namespaceManager := new(NamespaceManager) // NameSpace 管理员
	namespaceManager.namespaces = make(map[string]*Namespace, 64)
	namespaceManager.namespaces["db0_namespace"] = managerNamespace
	serverManager.namespaces[current] = namespaceManager

	// 指定服务器用户
	serverUser := new(UserManager)
	serverUser.users = make(map[string][]string, 64)
	serverUser.userNamespaces = make(map[string]string, 64)
	serverUser.userNamespaces["panhong"+":"+"12345"] = "db0_namespace"
	serverUser.users["panhong"] = append(serverUser.users["panhong"], "12345")

	// 把 服务器用户资料 写回 服务器管理员
	serverManager.users[current] = serverUser

	// 产生执行者资料
	executor := SessionExecutor{}
	executor.manager = serverManager
	executor.namespace = "db0_namespace" // 用 executor.namespace 和 current(切换标签)去取出NS值

	// 取出 NS
	ns := executor.GetNamespace()

	// 檢查 NS
	testNS := []struct { // NameSpace 底下有 Master 和 Slave 数据库，检查是否有正确载入
		db    string
		allow bool
	}{
		{ // 允许的数据库列表
			"novel", // SQL 字串內容
			true,    // 允许
		},
	}
	for _, test := range testNS { // 进行检查
		_, ok := ns.allowedDBs[test.db]
		require.Equal(t, test.allow, ok)
	}

	// 取出 Router
	rt := ns.GetRouter() // 会取出预设的 Slice 为 Slice-0

	// 补齐执行者资料
	executor.sessionVariables = mysql.NewSessionVariables()
	executor.txConns = make(map[string]backend.PooledConnect)
	executor.stmts = make(map[uint32]*Stmt)
	executor.parser = parser.New()
	executor.status = initClientConnStatus
	executor.user = "panhong"
	collationID := 33 // "utf8"
	executor.collation = mysql.CollationID(collationID)
	executor.charset = "utf8"
	executor.db = "novel"

	// 补齐 Session 执行者资料
	sessionExecutor := newSessionExecutor(serverManager)
	sessionExecutor.user = "panhong"
	collationID = 33 // "utf8"
	sessionExecutor.SetCollationID(mysql.CollationID(collationID))
	sessionExecutor.SetCharset("utf8")
	sessionExecutor.SetDatabase("novel") // set database
	sessionExecutor.namespace = "db0_namespace"

	// 检查数据库的 Parser
	testParser := []struct {
		sql    string
		expect string
	}{
		// 第一本小说 三国演义
		{
			"INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(1, 9781517191276, 'Romance Of The Three Kingdoms', 'Luo Guanzhong', 1522, 'Historical fiction');",       // 原始的 SQL 字串
			"INSERT INTO `novel`.`Book` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (1,9781517191276,'Romance Of The Three Kingdoms','Luo Guanzhong',1522,'Historical fiction')", // Parser 后的 SQL 字串
		},
		// 第二本小说 水浒传
		{
			"INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(2, 9789869442060, 'Water Margin', 'Shi Nai an', 1589, 'Historical fiction');",       // 原始的 SQL 字串
			"INSERT INTO `novel`.`Book` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (2,9789869442060,'Water Margin','Shi Nai an',1589,'Historical fiction')", // Parser 后的 SQL 字串
		},
		// 第三本小说 西游记
		{
			"INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(3, 9789575709518, 'Journey To The West', 'Wu Cheng en', 1592, 'Gods And Demons Fiction');",       // 原始的 SQL 字串
			"INSERT INTO `novel`.`Book` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (3,9789575709518,'Journey To The West','Wu Cheng en',1592,'Gods And Demons Fiction')", // Parser 后的 SQL 字串
		},
		// 第四本小说 红楼梦
		{
			"INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(4, 9789865975364, 'Dream Of The Red Chamber', 'Cao Xueqin', 1791, 'Family Saga');",       // 原始的 SQL 字串
			"INSERT INTO `novel`.`Book` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (4,9789865975364,'Dream Of The Red Chamber','Cao Xueqin',1791,'Family Saga')", // Parser 后的 SQL 字串
		},
		// 第五本小说 金瓶梅
		{
			"INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(5, 9780804847773, 'Jin Ping Mei', 'Lanling Xiaoxiao Sheng', 1610, 'Family Life');",       // 原始的 SQL 字串
			"INSERT INTO `novel`.`Book` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (5,9780804847773,'Jin Ping Mei','Lanling Xiaoxiao Sheng',1610,'Family Life')", // Parser 后的 SQL 字串
		},
		// 第六本小说 儒林外史
		{
			"INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(6, 9780835124072, 'Rulin Waishi', 'Wu Jingzi', 1750, 'Unofficial History');",       // 原始的 SQL 字串
			"INSERT INTO `novel`.`Book` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (6,9780835124072,'Rulin Waishi','Wu Jingzi',1750,'Unofficial History')", // Parser 后的 SQL 字串
		},
		// 第七本小说 初刻拍案惊奇
		{
			"INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(7, 9787101064100, 'Amazing Tales First Series', 'Ling Mengchu', 1628, 'Perspective');",       // 原始的 SQL 字串
			"INSERT INTO `novel`.`Book` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (7,9787101064100,'Amazing Tales First Series','Ling Mengchu',1628,'Perspective')", // Parser 后的 SQL 字串
		},
		// 第八本小说 二刻拍案惊奇
		{
			"INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(8, 9789571447278, 'Amazing Tales Second Series', 'Ling Mengchu', 1628, 'Perspective');",       // 原始的 SQL 字串
			"INSERT INTO `novel`.`Book` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (8,9789571447278,'Amazing Tales Second Series','Ling Mengchu',1628,'Perspective')", // Parser 后的 SQL 字串
		},
		// 第九本小说 封神演义
		{
			"INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(9, 9789861273129, 'Investiture Of The Gods', 'Lu Xixing', 1605, 'Mythology');",       // 原始的 SQL 字串
			"INSERT INTO `novel`.`Book` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (9,9789861273129,'Investiture Of The Gods','Lu Xixing',1605,'Mythology')", // Parser 后的 SQL 字串
		},
		// 第十本小说 镜花缘
		{
			"INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(10, 9787540251499, 'Flowers In The Mirror', 'Li Ruzhen', 1827, 'Fantasy Stories');",       // 原始的 SQL 字串
			"INSERT INTO `novel`.`Book` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (10,9787540251499,'Flowers In The Mirror','Li Ruzhen',1827,'Fantasy Stories')", // Parser 后的 SQL 字串
		},
		// 第十一本小说 镜花缘
		{
			"INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(11, 9787508535296, 'Stories Old And New', 'Feng Menglong', 1620, 'Perspective');",       // 原始的 SQL 字串
			"INSERT INTO `novel`.`Book` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (11,9787508535296,'Stories Old And New','Feng Menglong',1620,'Perspective')", // Parser 后的 SQL 字串
		},
		// 第十二本小说 说岳全传
		{
			"INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(12, 9787101097559, 'General Yue Fei', 'Qian Cai', 1735, 'History');",       // 原始的 SQL 字串
			"INSERT INTO `novel`.`Book` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (12,9787101097559,'General Yue Fei','Qian Cai',1735,'History')", // Parser 后的 SQL 字串
		},
		// 第十三本小说 杨家将
		{
			"INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(13, 9789863381037, 'The Generals Of The Yang Family', 'Qi Zhonglan', 0, 'History');",       // 原始的 SQL 字串
			"INSERT INTO `novel`.`Book` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (13,9789863381037,'The Generals Of The Yang Family','Qi Zhonglan',0,'History')", // Parser 后的 SQL 字串
		},
		// 第十四本小说 说唐
		{
			"INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(14, 9789865700027, 'Romance Of Sui And Tang Dynasties', 'Chen Ruheng', 1989, 'History');",       // 原始的 SQL 字串
			"INSERT INTO `novel`.`Book` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (14,9789865700027,'Romance Of Sui And Tang Dynasties','Chen Ruheng',1989,'History')", // Parser 后的 SQL 字串
		},
		// 第十五本小说 七侠五义
		{
			"INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(15, 9789575709242, 'The Seven Heroes And Five Gallants', 'Shi Yukun', 1879, 'History');",       // 原始的 SQL 字串
			"INSERT INTO `novel`.`Book` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (15,9789575709242,'The Seven Heroes And Five Gallants','Shi Yukun',1879,'History')", // Parser 后的 SQL 字串
		},
		// 第十六本小说 施公案
		{
			"INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(16, 9789574927913, 'A Collection Of Shi', 'Anonymous', 1850, 'History');",       // 原始的 SQL 字串
			"INSERT INTO `novel`.`Book` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (16,9789574927913,'A Collection Of Shi','Anonymous',1850,'History')", // Parser 后的 SQL 字串
		},
		// 第十七本小说 青楼梦
		{
			"INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(17, 9787533303396, 'Dream Of The Green Chamber', 'Yuda', 1878, 'Family Saga');",       // 原始的 SQL 字串
			"INSERT INTO `novel`.`Book` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (17,9787533303396,'Dream Of The Green Chamber','Yuda',1878,'Family Saga')", // Parser 后的 SQL 字串
		},
		// 第十八本小说 歧路灯
		{
			"INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(18, 9787510434341, 'Lamp In The Side Street', 'Li Luyuan', 1790, 'Unofficial History');",       // 原始的 SQL 字串
			"INSERT INTO `novel`.`Book` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (18,9787510434341,'Lamp In The Side Street','Li Luyuan',1790,'Unofficial History')", // Parser 后的 SQL 字串
		},
		// 第十九本小说 老残游记
		{
			"INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(19, 9789571447469, 'The Travels of Lao Can', 'Liu E', 1907, 'Social Story');",       // 原始的 SQL 字串
			"INSERT INTO `novel`.`Book` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (19,9789571447469,'The Travels of Lao Can','Liu E',1907,'Social Story')", // Parser 后的 SQL 字串
		},
		// 第二十本小说 二十年目睹之怪现状
		{
			"INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(20, 9789571470047, 'Bizarre Happenings Eyewitnessed over Two Decades', 'Jianren Wu', 1905, 'Unofficial History');",       // 原始的 SQL 字串
			"INSERT INTO `novel`.`Book` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (20,9789571470047,'Bizarre Happenings Eyewitnessed over Two Decades','Jianren Wu',1905,'Unofficial History')", // Parser 后的 SQL 字串
		},
		// 第二十一本小说 孽海花
		{
			"INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(21, 9787101097580, 'A Flower In A Sinful Sea', 'Zeng Pu', 1904, 'History');",       // 原始的 SQL 字串
			"INSERT INTO `novel`.`Book` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (21,9787101097580,'A Flower In A Sinful Sea','Zeng Pu',1904,'History')", // Parser 后的 SQL 字串
		},
		// 第二十二本小说 官场现形记
		{
			"INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(22, 9789861674193, 'Officialdom Unmasked', 'Li Baojia', 1903, 'Unofficial History');",       // 原始的 SQL 字串
			"INSERT INTO `novel`.`Book` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (22,9789861674193,'Officialdom Unmasked','Li Baojia',1903,'Unofficial History')", // Parser 后的 SQL 字串
		},
		// 第二十三本小说 觉世名言十二楼
		{
			"INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(23, 9787805469836, 'Tower For The Summer Heat', 'Li Yu', 1680, 'Unofficial History');",       // 原始的 SQL 字串
			"INSERT INTO `novel`.`Book` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (23,9787805469836,'Tower For The Summer Heat','Li Yu',1680,'Unofficial History')", // Parser 后的 SQL 字串
		},
		// 第二十四本小说 无声戏
		{
			"INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(24, 9787508067247, 'Silent Operas', 'Li Yu', 1680, 'Social Story');",       // 原始的 SQL 字串
			"INSERT INTO `novel`.`Book` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (24,9787508067247,'Silent Operas','Li Yu',1680,'Social Story')", // Parser 后的 SQL 字串
		},
		// 第二十五本小说 肉蒲团
		{
			"INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(25, 9789573609049, 'The Carnal Prayer Mat', 'Li Yu', 1680, 'Social Story');",       // 原始的 SQL 字串
			"INSERT INTO `novel`.`Book` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (25,9789573609049,'The Carnal Prayer Mat','Li Yu',1680,'Social Story')", // Parser 后的 SQL 字串
		},
		// 第二十六本小说 浮生六记
		{
			"INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(26, 9787533948108, 'Six Records Of A Floating Life', 'Shen Fu', 1878, 'Autobiography');",       // 原始的 SQL 字串
			"INSERT INTO `novel`.`Book` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (26,9787533948108,'Six Records Of A Floating Life','Shen Fu',1878,'Autobiography')", // Parser 后的 SQL 字串
		},
		// 第二十六本小说 浮生六记
		{
			"INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(26, 9787533948108, 'Six Records Of A Floating Life', 'Shen Fu', 1878, 'Autobiography');",       // 原始的 SQL 字串
			"INSERT INTO `novel`.`Book` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (26,9787533948108,'Six Records Of A Floating Life','Shen Fu',1878,'Autobiography')", // Parser 后的 SQL 字串
		},
		// 第二十七本小说 野叟曝言
		{
			"INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(27, 9786666141110, 'Humble Words Of A Rustic Elder', 'Xia Jingqu', 1787, 'Historical fiction');",       // 原始的 SQL 字串
			"INSERT INTO `novel`.`Book` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (27,9786666141110,'Humble Words Of A Rustic Elder','Xia Jingqu',1787,'Historical fiction')", // Parser 后的 SQL 字串
		},
		// 第二十八本小说 九尾龟
		{
			"INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(28, 9789571435473, 'Nine-Tailed Turtle', 'Lu Can', 1551, 'Mythology');",       // 原始的 SQL 字串
			"INSERT INTO `novel`.`Book` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (28,9789571435473,'Nine-Tailed Turtle','Lu Can',1551,'Mythology')", // Parser 后的 SQL 字串
		},
		// 第二十九本小说 品花宝鉴
		{
			"INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(29, 9789866318603, 'A History Of Floral Treasures', 'Chen Sen', 1849, 'Romance');",       // 原始的 SQL 字串
			"INSERT INTO `novel`.`Book` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (29,9789866318603,'A History Of Floral Treasures','Chen Sen',1849,'Romance')", // Parser 后的 SQL 字串
		},
	}

	// 数据库新增资料测试
	for _, test := range testParser {
		// 先 Parser
		node, warns, err := executor.parser.Parse(test.sql, "utf8", "utf8_general_ci")
		require.Equal(t, []error(nil), warns)
		require.Equal(t, nil, err)

		// 再 Restore
		s := &strings.Builder{}
		ctx := format.NewRestoreCtx(format.EscapeRestoreFlags, s)
		err = node[0].Restore(ctx)
		require.Equal(t, nil, err)

		// 检查 Restore 后的 SQL 字串
		require.Equal(t, test.expect, s.String())

		// 获得 计划
		_, err = plan.BuildPlan(node[0], ns.GetPhysicalDBs(), executor.db, s.String(), rt, ns.GetSequences())
		require.Equal(t, nil, err)

		// BuildPlan 函式会产生计划 p 变数，
		// 再执行 sessionExecutor.ExecuteSQL(util.NewRequestContext(), "slice-0", "novel", s.String())，
		// ret, err := p.ExecuteIn(util.NewRequestContext(), sessionExecutor)
		// 把SQL 字串写入数据库，但因会触发统计功能，暂时先不使用

		// 以下操作会真的写入数据库，在此先中断
		return

		// 对数据库进行直连操作
		dc, err := executor.manager.GetNamespace(executor.namespace).GetSlice("slice-0").GetDirectConn("192.168.122.2:3306")
		require.Equal(t, nil, err)
		dc.Execute(s.String(), 50)
		fmt.Println("目前执行的数据库指令为 ", s.String())
	}
}

// TestB4 内含把 29 本小说写入数据库的程式码，这时开始进行数据库读写分离
func TestB4(t *testing.T) {
	// 组成管理员的 NameSpace (会写入变数 SessionExecutor -> namespaces[0] -> namespaces["key"])
	managerNamespace := &Namespace{} // 管理员的 NameSpace
	defer managerNamespace.Close(false)
	managerNamespace.openGeneralLog = false             // 记录sql查询的访问日志，说明 https://github.com/XiaoMi/Gaea/issues/109
	managerNamespace.name = "db0_cluster_namespace"     // namespace 为划分工作业务的最基本单位，一个 namespace 可以有多个使用者
	managerNamespace.allowedDBs = make(map[string]bool) // 数据库列表
	managerNamespace.allowedDBs[strings.TrimSpace("novel")] = true
	defaultPhyDBs := make(map[string]string, 0) // 预设数据库列表
	// defaultPhyDBs[strings.TrimSpace("novel")] = strings.TrimSpace("novel") // 再指定
	managerNamespace.defaultPhyDBs, _ = parseDefaultPhyDB(defaultPhyDBs, managerNamespace.allowedDBs)
	managerNamespace.slowSQLTime = 1000                      // 慢sql阈值，单位: 毫秒
	managerNamespace.sqls = make(map[string]string, 50)      // 有 29 本小说，先暂定 50 好了
	managerNamespace.sqls = parseBlackSqls([]string{})       // 黑名单sql
	managerNamespace.allowips, _ = parseAllowIps([]string{}) // 白名单IP

	// >>>>> 组织主从物理实例 (会写入变数 SessionExecutor -> namespaces[0] -> slices["key"])
	var dbSlices []*models.Slice // 一主多从的物理实例，slice里map的具体字段可参照slice配置
	dbSlice := models.Slice{}
	dbSlice.Name = "slice-0"                                          // 分片名称，自动、有序生成
	dbSlice.UserName = "docker"                                       // 连接后端mysql所需要的用户名称
	dbSlice.Password = "12345"                                        // 连接后端mysql所需要的用户密码
	dbSlice.Master = "192.168.1.2:3350"                               // 主实例地址 (db0 192.168.1.2:3350)
	dbSlice.Slaves = []string{"192.168.1.2:3351", "192.168.1.2:3352"} // 从实例地址列表  (db0-0 192.168.1.2:3351; db0-1 192.168.1.2:3352)
	dbSlice.StatisticSlaves = nil                                     // 统计型从实例地址列表
	dbSlice.Capacity = 12                                             // gaea_proxy与每个实例的连接池大小
	dbSlice.MaxCapacity = 24                                          // gaea_proxy与每个实例的连接池最大大小
	dbSlice.IdleTimeout = 60                                          // gaea_proxy与后端mysql空闲连接存活时间，单位:秒
	dbSlices = append(dbSlices, &dbSlice)

	// 管理员的 NameSpace 載入 组织主从物理实例 (会写入变数 SessionExecutor -> namespaces[0] -> slices["key"])
	managerNamespace.slices, _ = parseSlices(dbSlices, "utf8mb4", managerNamespace.defaultCollationID) // 一主多从的物理实例，slice里map的具体字段可参照slice配置

	// >>>>> 组织 NS 用户 (会写入变数 SessionExecutor -> namespaces[0] -> userProperties["key"])
	var nsUsers []*models.User // 应用端连接gaea所需要的用户配置，具体字段可参照users配置
	nsUser := models.User{}
	nsUser.UserName = "docker"                 // 用户名
	nsUser.Password = "12345"                  // 用户密码
	nsUser.Namespace = "db0_cluster_namespace" // 对应的命名空间
	nsUser.RWFlag = 2                          // 读写标识
	nsUser.RWSplit = 1                         // 是否读写分离
	nsUser.OtherProperty = 0                   // 其他属性，目前用来标识是否走统计从实例
	nsUsers = append(nsUsers, &nsUser)

	// 管理员的 NameSpace 載入 NS 用户 (会写入变数 SessionExecutor -> namespaces[0] -> userProperties["key"])
	managerNamespace.userProperties = make(map[string]*UserProperty, 2)
	for _, user := range nsUsers {
		up := &UserProperty{RWFlag: user.RWFlag, RWSplit: user.RWSplit, OtherProperty: user.OtherProperty}
		managerNamespace.userProperties[user.UserName] = up
	}

	// >>>>> 建立 全局唯一序列号 (会写入变数 SessionExecutor -> namespaces[0] -> sequence["key"])
	sequences := sequence.NewSequenceManager() // 生成全局唯一序列号的配置, 具体字段可参考全局序列号配置
	for _, v := range []*models.GlobalSequence{} {
		globalSequenceSlice, _ := managerNamespace.slices[v.SliceName]
		seqName := strings.ToUpper(v.DB) + "." + strings.ToUpper(v.Table)
		seq := sequence.NewMySQLSequence(globalSequenceSlice, seqName, v.PKName)
		sequences.SetSequence(v.DB, v.Table, seq)
	}

	// 管理员的 NameSpace 載入 全局唯一序列号 (会写入变数 SessionExecutor -> namespaces[0] -> sequence["key"])
	managerNamespace.sequences = sequences

	// 管理员的 NameSpace 載入 预设选项
	/*
		补齐所有的 managerNamespace 变数，最后先组成 managerNamespace (物品)
		namespaceManager(map) 的人 指定 db0_cluster_namespace 字串 map 到前面的 managerNamespace (物品)
		合拼到服務器管理員 serverManager := new(Manager)
		再合拼到 Session 执行者 executor.manager = serverManager
	*/
	managerNamespace.defaultCharset, managerNamespace.defaultCollationID, _ = parseCharset("", "") // 用于指定数据集如何排序，以及字符串的比对规则 & 用于指定数据集如何排序，以及字符串的比对规则
	managerNamespace.maxSqlExecuteTime = defaultMaxSqlExecuteTime                                  // sql最大执行时间，大于该时间，进行熔断
	managerNamespace.maxSqlResultSize = defaultMaxSqlResultSize                                    // 限制单分片返回结果集大小不超过max_select_rows

	// 管理员的 NameSpace 載入 延伸选择
	managerNamespace.slowSQLCache = cache.NewLRUCache(defaultSQLCacheCapacity)
	managerNamespace.errorSQLCache = cache.NewLRUCache(defaultSQLCacheCapacity)
	managerNamespace.backendSlowSQLCache = cache.NewLRUCache(defaultSQLCacheCapacity)
	managerNamespace.backendErrorSQLCache = cache.NewLRUCache(defaultSQLCacheCapacity)
	managerNamespace.planCache = cache.NewLRUCache(defaultPlanCacheCapacity)

	// 建立路由
	modelsNameSpace := new(models.Namespace)
	modelsNameSpace.Slices = []*models.Slice{} // 一主多从的物理实例，slice里map的具体字段可参照slice配置
	modelsNameSpace.Slices = append(modelsNameSpace.Slices, &dbSlice)
	modelsNameSpace.DefaultSlice = "slice-0" // 预设分片名称
	modelsNameSpace.ShardRules = nil         // 分库、分表、特殊表的配置内容，具体字段可参照shard配置 (载入设定档)
	managerNamespace.router, _ = router.NewRouter(modelsNameSpace)

	// 初始化管理员 Manager
	serverManager := new(Manager)             // 服务器管理员
	current := 0                              // 切换标签
	namespaceManager := new(NamespaceManager) // NameSpace 管理员
	namespaceManager.namespaces = make(map[string]*Namespace, 64)
	namespaceManager.namespaces["db0_cluster_namespace"] = managerNamespace
	serverManager.namespaces[current] = namespaceManager

	// 指定服务器用户
	serverUser := new(UserManager)
	serverUser.users = make(map[string][]string, 64)
	serverUser.userNamespaces = make(map[string]string, 64)
	serverUser.userNamespaces["docker"+":"+"12345"] = "db0_cluster_namespace"
	serverUser.users["docker"] = append(serverUser.users["docker"], "12345")

	// 把 服务器用户资料 写回 服务器管理员
	serverManager.users[current] = serverUser

	// 产生执行者资料
	executor := SessionExecutor{}
	executor.manager = serverManager
	executor.namespace = "db0_cluster_namespace" // 用 executor.namespace 和 current(切换标签)去取出NS值

	// 取出 NS
	ns := executor.GetNamespace()

	// 檢查 NS
	testNS := []struct { // NameSpace 底下有 Master 和 Slave 数据库，检查是否有正确载入
		db    string
		allow bool
	}{
		{ // 允许的数据库列表
			"novel", // SQL 字串內容
			true,    // 允许
		},
	}
	for _, test := range testNS { // 进行检查
		_, ok := ns.allowedDBs[test.db]
		require.Equal(t, test.allow, ok)
	}

	// 取出 Router
	// rt := ns.GetRouter() // 会取出预设的 Slice 为 Slice-0

	// 补齐执行者资料
	executor.sessionVariables = mysql.NewSessionVariables()
	executor.txConns = make(map[string]backend.PooledConnect)
	executor.stmts = make(map[uint32]*Stmt)
	executor.parser = parser.New()
	executor.status = initClientConnStatus
	executor.user = "docker"
	collationID := 33 // "utf8"
	executor.collation = mysql.CollationID(collationID)
	executor.charset = "utf8"
	executor.db = "novel"

	// 补齐 Session 执行者资料
	sessionExecutor := newSessionExecutor(serverManager)
	sessionExecutor.user = "docker"
	collationID = 33 // "utf8"
	sessionExecutor.SetCollationID(mysql.CollationID(collationID))
	sessionExecutor.SetCharset("utf8")
	sessionExecutor.SetDatabase("novel") // set database
	sessionExecutor.namespace = "db0_cluster_namespace"

	// 开始检查和数据库的沟通
	tests := []struct {
		sql    string
		expect string
	}{
		{ // 测试一，查询数据库资料
			"SELECT * FROM novel.Book",     // 原始的 SQL 字串
			"SELECT * FROM `novel`.`Book`", // 期望 Parser 后的 SQL 字串
		},
	}

	// 执行 Sql 字串
	for _, test := range tests {
		// 執行 SQL Parser
		ns := executor.GetNamespace()
		stmts, err := executor.Parse(test.sql)
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
		p, err := plan.BuildPlan(stmts, phyDBs, "novel", test.sql, rt, seq)
		require.Equal(t, err, nil)

		// 以下会直接连线到实体数据库，先在这里中断
		return

		// 执行 Parser 后的 SQL 指令
		reqCtx := util.NewRequestContext()
		reqCtx.Set(util.FromSlave, 1) // 在这里设定读取时从 Slave 节点，达到读写分离的效果

		// 要针对写单元测试 se.manager.RecordBackendSQLMetrics(reqCtx, se.namespace, v, pc.GetAddr(), startTime, err)，先中断
		// return

		res, err := p.ExecuteIn(reqCtx, sessionExecutor)
		require.Equal(t, err, nil)

		// 检查数据库回传第 1 本书的资料
		require.Equal(t, res.Resultset.Values[0][0].(int64), int64(1))
		require.Equal(t, res.Resultset.Values[0][1].(int64), int64(9781517191276))
		require.Equal(t, res.Resultset.Values[0][2].(string), "Romance Of The Three Kingdoms")

		// 检查数据库回传第 28 本书的资料
		require.Equal(t, res.Resultset.Values[28][0].(int64), int64(29))
		require.Equal(t, res.Resultset.Values[28][1].(int64), int64(9789866318603))
		require.Equal(t, res.Resultset.Values[28][2].(string), "A History Of Floral Treasures")
	}
}
