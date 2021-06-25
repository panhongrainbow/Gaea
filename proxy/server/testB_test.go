package server

import (
	"fmt"

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

// 设定档 key value 细部测试
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

func TestB5(t *testing.T) {
	// 组成管理员的 NameSpace
	managerNamespace := &Namespace{} // 管理员的 NameSpace
	defer managerNamespace.Close(false)
	managerNamespace.openGeneralLog = false             // 记录sql查询的访问日志，说明 https://github.com/XiaoMi/Gaea/issues/109
	managerNamespace.name = "env1_namespace_1"          // namespace 为划分工作业务的最基本单位，一个 namespace 可以有多个使用者
	managerNamespace.allowedDBs = make(map[string]bool) // 数据库列表
	managerNamespace.allowedDBs[strings.TrimSpace("Library")] = true
	defaultPhyDBs := make(map[string]string, 0) // 预设数据库列表
	// defaultPhyDBs[strings.TrimSpace(db)] = strings.TrimSpace(phyDB) // 再指定
	managerNamespace.defaultPhyDBs, _ = parseDefaultPhyDB(defaultPhyDBs, managerNamespace.allowedDBs)
	managerNamespace.slowSQLTime = 1000 // 慢sql阈值，单位: 毫秒
	managerNamespace.sqls = make(map[string]string, 16)
	managerNamespace.sqls = parseBlackSqls([]string{})       // 黑名单sql
	managerNamespace.allowips, _ = parseAllowIps([]string{}) // 白名单IP

	// >>>>> 组织主从物理实例
	var dbSlices []*models.Slice // 一主多从的物理实例，slice里map的具体字段可参照slice配置
	dbSlice := models.Slice{}
	dbSlice.Name = "slice-0"           // 分片名称，自动、有序生成
	dbSlice.UserName = "root"          // 连接后端mysql所需要的用户名称
	dbSlice.Password = "12345"         // 连接后端mysql所需要的用户密码
	dbSlice.Master = "172.17.0.2:3306" // 主实例地址
	dbSlice.Slaves = []string{}        // 从实例地址列表
	dbSlice.StatisticSlaves = nil      // 统计型从实例地址列表
	dbSlice.Capacity = 12              // gaea_proxy与每个实例的连接池大小
	dbSlice.MaxCapacity = 24           // gaea_proxy与每个实例的连接池最大大小
	dbSlice.IdleTimeout = 60           // gaea_proxy与后端mysql空闲连接存活时间，单位:秒
	dbSlices = append(dbSlices, &dbSlice)

	// 管理员的 NameSpace 載入 组织主从物理实例
	managerNamespace.slices, _ = parseSlices(dbSlices, "utf8mb4", managerNamespace.defaultCollationID) // 一主多从的物理实例，slice里map的具体字段可参照slice配置

	// >>>>> 组织 NS 用户
	var nsUsers []*models.User // 应用端连接gaea所需要的用户配置，具体字段可参照users配置
	nsUser := models.User{}
	nsUser.UserName = "root"              // 用户名
	nsUser.Password = "12345"             // 用户密码
	nsUser.Namespace = "env1_namespace_1" // 对应的命名空间
	nsUser.RWFlag = 2                     // 读写标识
	nsUser.RWSplit = 1                    // 是否读写分离
	nsUser.OtherProperty = 0              // 其他属性，目前用来标识是否走统计从实例
	nsUsers = append(nsUsers, &nsUser)

	// 管理员的 NameSpace 載入 NS 用户
	managerNamespace.userProperties = make(map[string]*UserProperty, 2)
	for _, user := range nsUsers {
		up := &UserProperty{RWFlag: user.RWFlag, RWSplit: user.RWSplit, OtherProperty: user.OtherProperty}
		managerNamespace.userProperties[user.UserName] = up
	}

	// >>>>> 建立 全局唯一序列号
	sequences := sequence.NewSequenceManager() // 生成全局唯一序列号的配置, 具体字段可参考全局序列号配置
	for _, v := range []*models.GlobalSequence{} {
		globalSequenceSlice, _ := managerNamespace.slices[v.SliceName]
		seqName := strings.ToUpper(v.DB) + "." + strings.ToUpper(v.Table)
		seq := sequence.NewMySQLSequence(globalSequenceSlice, seqName, v.PKName)
		sequences.SetSequence(v.DB, v.Table, seq)
	}

	// 管理员的 NameSpace 載入 全局唯一序列号
	managerNamespace.sequences = sequences

	// 管理员的 NameSpace 載入 预设选项
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
	namespaceManager.namespaces["env1_namespace_1"] = managerNamespace
	serverManager.namespaces[current] = namespaceManager

	// 指定服务器用户
	serverUser := new(UserManager)
	serverUser.users = make(map[string][]string, 64)
	serverUser.userNamespaces = make(map[string]string, 64)
	serverUser.userNamespaces["root"+":"+"12345"] = "env1_namespace_1"
	serverUser.users["root"] = append(serverUser.users["root"], "12345")

	// 把 服务器用户资料 写回 服务器管理员
	serverManager.users[current] = serverUser

	// 产生执行者资料
	executor := SessionExecutor{}
	executor.manager = serverManager
	executor.namespace = "env1_namespace_1" // 用 executor.namespace 和 current(切换标签)去取出NS值

	// 取出 NS
	ns := executor.GetNamespace()

	// 檢查 NS
	testNS := []struct { // NameSpace 底下有 Master 和 Slave 数据库，检查是否有正确载入
		db    string
		allow bool
	}{
		{ // 允许的数据库列表
			"Library", // SQL 字串內容
			true,      // 允许
		},
	}
	for _, test := range testNS { // 进行检查
		_, ok := ns.allowedDBs[test.db]
		require.Equal(t, test.allow, ok)
	}

	// 取出 Router
	rt := ns.GetRouter() // 会取出预设的 Slice 为 Slice-0
	fmt.Println(rt)

	// 补齐执行者资料
	executor.sessionVariables = mysql.NewSessionVariables()
	executor.txConns = make(map[string]backend.PooledConnect)
	executor.stmts = make(map[uint32]*Stmt)
	executor.parser = parser.New()
	executor.status = initClientConnStatus
	executor.user = "root"
	collationID := 33 // "utf8"
	executor.collation = mysql.CollationID(collationID)
	executor.charset = "utf8"
	executor.db = "Library"

	// 检查数据库的 Parser
	testParser := []struct {
		sql    string
		expect string
	}{
		{
			"INSERT INTO Library.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(1, 9789865975364, 'Dream of the Red Chamber', 'Cao Xueqin', 1791, 'Family Saga');",       // 原始的 SQL 字串
			"INSERT INTO `Library`.`Book` (`BookID`,`Isbn`,`Title`,`Author`,`Publish`,`Category`) VALUES (1,9789865975364,'Dream of the Red Chamber','Cao Xueqin',1791,'Family Saga')", // Parser 后的 SQL 字串
		},
	}
	for _, test := range testParser {
		n, _, _ := executor.parser.Parse(test.sql, "", "")
		s := &strings.Builder{}
		ctx := format.NewRestoreCtx(format.EscapeRestoreFlags, s)
		_ = n[0].Restore(ctx)
		require.Equal(t, test.expect, s.String())

		// 获得 计划
		db := executor.db
		seq := ns.GetSequences()
		phyDBs := ns.GetPhysicalDBs()
		p, _ := plan.BuildPlan(n[0], phyDBs, db, s.String(), rt, seq)
		fmt.Println(p)
	}
}
