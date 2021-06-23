package server

import (
	"fmt"

	"github.com/XiaoMi/Gaea/backend"
	"github.com/XiaoMi/Gaea/models"
	"github.com/XiaoMi/Gaea/mysql"
	"github.com/XiaoMi/Gaea/parser"
	"github.com/XiaoMi/Gaea/parser/format"
	"github.com/XiaoMi/Gaea/proxy/router"
	"github.com/XiaoMi/Gaea/proxy/sequence"

	// 先忽略 stats 包，因为这个测试会跟其他测试发生冲突
	// "github.com/XiaoMi/Gaea/stats"
	"net/http"
	"strings"
	"testing"

	"github.com/XiaoMi/Gaea/util"
	"github.com/XiaoMi/Gaea/util/cache"
	"github.com/XiaoMi/Gaea/util/sync2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

func TestB2(t *testing.T) {
	proxy := new(models.Proxy)
	proxy.ConfigType = "file"
	proxy.FileConfigPath = "./etc/file"
	proxy.CoordinatorAddr = "http://127.0.0.1:2379"
	proxy.CoordinatorRoot = ""
	proxy.UserName = "root"
	proxy.Password = "root"
	proxy.Environ = "local"
	proxy.Service = "gaea_proxy"
	proxy.Cluster = "gaea_default_cluster"
	proxy.LogPath = "./logs"
	proxy.LogLevel = "Notice"
	proxy.LogFileName = "gaea"
	proxy.LogOutput = "file"
	proxy.ProtoType = "tcp4"
	proxy.ProxyAddr = "0.0.0.0:13306"
	proxy.AdminAddr = "0.0.0.0:13307"
	proxy.AdminUser = "admin"
	proxy.AdminPassword = "admin"
	proxy.SlowSQLTime = 100
	proxy.SessionTimeout = 3600
	proxy.StatsEnabled = "true"
	proxy.StatsInterval = 10
	proxy.EncryptKey = "1234abcd5678efg*"

	modelsNameSpace := new(models.Namespace)
	modelsNameSpace.OpenGeneralLog = false
	modelsNameSpace.IsEncrypt = false
	modelsNameSpace.Name = "env1_namespace_1"
	modelsNameSpace.Online = true
	modelsNameSpace.ReadOnly = false
	modelsNameSpace.AllowedDBS = make(map[string]bool)
	modelsNameSpace.AllowedDBS["Library"] = true
	modelsNameSpace.DefaultPhyDBS = nil
	modelsNameSpace.SlowSQLTime = "1000"
	modelsNameSpace.BlackSQL = []string{}
	modelsNameSpace.AllowedIP = nil
	modelsNameSpace.Slices = []*models.Slice{}
	tmp1 := models.Slice{}
	tmp1.Name = "slice-0"
	tmp1.UserName = "root"
	tmp1.Password = "12345"
	tmp1.Master = "172.17.0.2:3306"
	tmp1.Slaves = []string{}
	tmp1.StatisticSlaves = nil
	tmp1.Capacity = 12
	tmp1.MaxCapacity = 24
	tmp1.IdleTimeout = 60
	modelsNameSpace.Slices = append(modelsNameSpace.Slices, &tmp1)
	modelsNameSpace.ShardRules = nil
	modelsNameSpace.Users = []*models.User{}
	tmp2 := models.User{}
	tmp2.UserName = "root"
	tmp2.Password = "12345"
	tmp2.Namespace = "env1_namespace_1"
	tmp2.RWFlag = 2
	tmp2.RWSplit = 1
	tmp2.OtherProperty = 0
	modelsNameSpace.Users = append(modelsNameSpace.Users, &tmp2)
	modelsNameSpace.DefaultSlice = "slice-0"
	modelsNameSpace.GlobalSequences = nil
	modelsNameSpace.DefaultCharset = ""
	modelsNameSpace.DefaultCollation = ""
	modelsNameSpace.MaxSqlExecuteTime = 0
	modelsNameSpace.MaxSqlResultSize = 0

	serverManager := new(Manager)
	// current, _, _ := serverManager.switchIndex.Get()
	current := 0
	// namespaceMap := map[string]*models.Namespace{"env1_namespace_1": modelsNameSpace}

	// serverManager.namespaces[current] = CreateNamespaceManager(namespaceMap)
	nsMgr := new(NamespaceManager)
	nsMgr.namespaces = make(map[string]*Namespace, 64)

	tmp3 := &Namespace{}
	tmp3.sqls = make(map[string]string, 16)
	tmp3.userProperties = make(map[string]*UserProperty, 2)
	tmp3.openGeneralLog = modelsNameSpace.OpenGeneralLog
	tmp3.slowSQLCache = cache.NewLRUCache(defaultSQLCacheCapacity)
	tmp3.errorSQLCache = cache.NewLRUCache(defaultSQLCacheCapacity)
	tmp3.backendSlowSQLCache = cache.NewLRUCache(defaultSQLCacheCapacity)
	tmp3.backendErrorSQLCache = cache.NewLRUCache(defaultSQLCacheCapacity)
	tmp3.planCache = cache.NewLRUCache(defaultPlanCacheCapacity)
	defer tmp3.Close(false)
	tmp3.sqls = parseBlackSqls(modelsNameSpace.BlackSQL)
	tmp3.slowSQLTime, _ = parseSlowSQLTime(modelsNameSpace.SlowSQLTime)

	if modelsNameSpace.MaxSqlExecuteTime <= 0 {
		tmp3.maxSqlExecuteTime = defaultMaxSqlExecuteTime
	} else {
		tmp3.maxSqlExecuteTime = modelsNameSpace.MaxSqlExecuteTime
	}

	if modelsNameSpace.MaxSqlResultSize <= 0 {
		tmp3.maxSqlResultSize = defaultMaxSqlResultSize
	} else {
		tmp3.maxSqlResultSize = modelsNameSpace.MaxSqlResultSize
	}

	allowDBs := make(map[string]bool, len(modelsNameSpace.AllowedDBS))
	for db, allowed := range modelsNameSpace.AllowedDBS {
		allowDBs[strings.TrimSpace(db)] = allowed
	}
	tmp3.allowedDBs = allowDBs

	defaultPhyDBs := make(map[string]string, len(modelsNameSpace.DefaultPhyDBS))
	for db, phyDB := range modelsNameSpace.DefaultPhyDBS {
		defaultPhyDBs[strings.TrimSpace(db)] = strings.TrimSpace(phyDB)
	}

	tmp3.defaultPhyDBs, _ = parseDefaultPhyDB(defaultPhyDBs, allowDBs)

	allowips, _ := parseAllowIps(modelsNameSpace.AllowedIP)
	tmp3.allowips = allowips

	tmp3.defaultCharset, tmp3.defaultCollationID, _ = parseCharset(modelsNameSpace.DefaultCharset, modelsNameSpace.DefaultCollation)

	for _, user := range modelsNameSpace.Users {
		up := &UserProperty{RWFlag: user.RWFlag, RWSplit: user.RWSplit, OtherProperty: user.OtherProperty}
		tmp3.userProperties[user.UserName] = up
	}

	tmp3.slices, _ = parseSlices(modelsNameSpace.Slices, tmp3.defaultCharset, tmp3.defaultCollationID)

	tmp3.router, _ = router.NewRouter(modelsNameSpace)

	sequences := sequence.NewSequenceManager()
	for _, v := range modelsNameSpace.GlobalSequences {
		globalSequenceSlice, _ := tmp3.slices[v.SliceName]
		seqName := strings.ToUpper(v.DB) + "." + strings.ToUpper(v.Table)
		seq := sequence.NewMySQLSequence(globalSequenceSlice, seqName, v.PKName)
		sequences.SetSequence(v.DB, v.Table, seq)
	}
	tmp3.sequences = sequences

	nsMgr.namespaces[modelsNameSpace.Name] = tmp3
	serverManager.namespaces[current] = nsMgr

	// user, _ := CreateUserManager(namespaceMap)
	user := new(UserManager)
	user.users = make(map[string][]string, 64)
	user.userNamespaces = make(map[string]string, 64)
	// user.addNamespaceUsers(modelsNameSpace)
	user.userNamespaces["root"+":"+"12345"] = "env1_namespace_1"
	user.users["root"] = append(user.users["root"], "12345")
	serverManager.users[current] = user

	statisticManager := StatisticManager{}
	statisticManager.manager = new(Manager)
	statisticManager.manager.reloadPrepared = sync2.NewAtomicBool(false)
	tmp4 := util.BoolIndex{}
	tmp4.Set(true)
	statisticManager.manager.switchIndex = tmp4
	tmp5 := NamespaceManager{}
	tmp6 := NamespaceManager{}
	statisticManager.manager.namespaces[0] = &tmp5
	statisticManager.manager.namespaces[1] = &tmp6
	tmp7 := UserManager{}
	tmp8 := UserManager{}
	statisticManager.manager.users[0] = &tmp7
	statisticManager.manager.users[1] = &tmp8
	statisticManager.manager.statistics = nil // 里面有更多内容
	statisticManager.clusterName = "gaea_default_cluster"
	statisticManager.statsType = ""
	statisticManager.handlers = make(map[string]http.Handler)
	statisticManager.handlers["/metrics"] = promhttp.Handler() // 暂时，未确认
	tmp9, _ := initGeneralLogger(proxy)
	statisticManager.generalLogger = tmp9
	// 在这里和其他的测试会发生冲突，在这里先忽略，因为发现把其他在/proxy/server 里面的测试档删除，整体测试就通过，可能是因为在统计时有publish 函式，其他函式已经publish 一次了，这个测试在publish 一次，就会发生冲突，到时看看要不要跟其他测试合拼，这里先注解
	/*statisticManager.sqlTimings = stats.NewMultiTimings("SqlTimings",
		"gaea proxy sql sqlTimings", []string{statsLabelCluster, statsLabelNamespace, statsLabelOperation})
	statisticManager.sqlFingerprintSlowCounts = stats.NewCountersWithMultiLabels("SqlFingerprintSlowCounts",
		"gaea proxy sql fingerprint slow counts", []string{statsLabelCluster, statsLabelNamespace, statsLabelFingerprint})
	statisticManager.sqlErrorCounts = stats.NewCountersWithMultiLabels("SqlErrorCounts",
		"gaea proxy sql error counts per error type", []string{statsLabelCluster, statsLabelNamespace, statsLabelOperation})
	statisticManager.sqlFingerprintErrorCounts = stats.NewCountersWithMultiLabels("SqlFingerprintErrorCounts",
		"gaea proxy sql fingerprint error counts", []string{statsLabelCluster, statsLabelNamespace, statsLabelFingerprint})
	statisticManager.sqlForbidenCounts = stats.NewCountersWithMultiLabels("SqlForbiddenCounts",
		"gaea proxy sql error counts per error type", []string{statsLabelCluster, statsLabelNamespace, statsLabelFingerprint})
	statisticManager.flowCounts  = stats.NewCountersWithMultiLabels("FlowCounts",
		"gaea proxy flow counts", []string{statsLabelCluster, statsLabelNamespace, statsLabelFlowDirection})
	statisticManager.sessionCounts = stats.NewGaugesWithMultiLabels("SessionCounts",
		"gaea proxy session counts", []string{statsLabelCluster, statsLabelNamespace})
	statisticManager.backendSQLTimings = stats.NewMultiTimings("BackendSqlTimings",
		"gaea proxy backend sql sqlTimings", []string{statsLabelCluster, statsLabelNamespace, statsLabelOperation})
	statisticManager.backendSQLFingerprintSlowCounts = stats.NewCountersWithMultiLabels("BackendSqlFingerprintSlowCounts",
		"gaea proxy backend sql fingerprint slow counts", []string{statsLabelCluster, statsLabelNamespace, statsLabelFingerprint})
	statisticManager.backendSQLErrorCounts = stats.NewCountersWithMultiLabels("BackendSqlErrorCounts",
		"gaea proxy backend sql error counts per error type", []string{statsLabelCluster, statsLabelNamespace, statsLabelOperation})
	statisticManager.backendSQLFingerprintErrorCounts = stats.NewCountersWithMultiLabels("BackendSqlFingerprintErrorCounts",
		"gaea proxy backend sql fingerprint error counts", []string{statsLabelCluster, statsLabelNamespace, statsLabelFingerprint})
	statisticManager.backendConnectPoolIdleCounts = stats.NewGaugesWithMultiLabels("backendConnectPoolIdleCounts",
		"gaea proxy backend idle connect counts", []string{statsLabelCluster, statsLabelNamespace, statsLabelSlice, statsLabelIPAddr})
	statisticManager.backendConnectPoolInUseCounts = stats.NewGaugesWithMultiLabels("backendConnectPoolInUseCounts",
			"gaea proxy backend in-use connect counts", []string{statsLabelCluster, statsLabelNamespace, statsLabelSlice, statsLabelIPAddr})
	statisticManager.backendConnectPoolWaitCounts = stats.NewGaugesWithMultiLabels("backendConnectPoolWaitCounts",
			"gaea proxy backend wait connect counts", []string{statsLabelCluster, statsLabelNamespace, statsLabelSlice, statsLabelIPAddr})*/
	statisticManager.slowSQLTime = 100
	statisticManager.closeChan = make(chan bool, 0)
	executor := SessionExecutor{}
	executor.sessionVariables = mysql.NewSessionVariables()
	executor.txConns = make(map[string]backend.PooledConnect)
	executor.stmts = make(map[uint32]*Stmt)
	executor.parser = parser.New()
	executor.status = initClientConnStatus
	executor.manager = serverManager
	executor.user = "root"
	collationID := 33 // "utf8"
	executor.collation = mysql.CollationID(collationID)
	executor.charset = "utf8"
	executor.db = "Library"
	executor.namespace = "env1_namespace_1"

	// 开始检查和数据库的沟通
	tests := []struct {
		sql    string
		expect string
	}{
		{ // 执行 SQL
			"INSERT t.* VALUES (1), (2), (3)", // SQL 字串内容
			"",                                // 期望的 SQL 字串关键字
		},
	}

	// 执行测试
	for _, test := range tests {
		result, _, _ := executor.parser.Parse(test.sql, "", "")
		s := &strings.Builder{}
		ctx := format.NewRestoreCtx(format.EscapeRestoreFlags, s)
		_ = result[0].Restore(ctx)
		fmt.Println(s.String())

		serverNameSpace := executor.GetNamespace()
		router := serverNameSpace.GetRouter()
		fmt.Println(router)
	}
}

func TestB4(t *testing.T) {
	// 初始化测定档 (Load0)
	modelsNameSpace := new(models.Namespace)

	// >>>>> 载入NS设定档 (Load1)
	modelsNameSpace.OpenGeneralLog = false             // 记录sql查询的访问日志，说明 https://github.com/XiaoMi/Gaea/issues/109
	modelsNameSpace.IsEncrypt = false                  // true: 加密存储 false: 非加密存储，目前加密Slice、User中的用户名、密码
	modelsNameSpace.Name = "env1_namespace_1"          // namespace 为划分工作业务的最基本单位，一个 namespace 可以有多个使用者
	modelsNameSpace.Online = true                      // 是否在线，逻辑上下线使用
	modelsNameSpace.ReadOnly = false                   // 是否只读，namespace级别
	modelsNameSpace.AllowedDBS = make(map[string]bool) // 数据库列表
	modelsNameSpace.AllowedDBS["Library"] = true       // 数据库列表
	modelsNameSpace.DefaultPhyDBS = nil                // 预设数据库列表
	modelsNameSpace.SlowSQLTime = "1000"               // 慢sql阈值，单位: 毫秒
	modelsNameSpace.BlackSQL = []string{}              // 黑名单sql
	modelsNameSpace.AllowedIP = nil                    // 白名单IP

	// >>>>> 组织主从物理实例
	modelsNameSpace.Slices = []*models.Slice{} // 一主多从的物理实例，slice里map的具体字段可参照slice配置
	slicePiece := models.Slice{}
	slicePiece.Name = "slice-0"           // 分片名称，自动、有序生成
	slicePiece.UserName = "root"          // 连接后端mysql所需要的用户名称
	slicePiece.Password = "12345"         // 连接后端mysql所需要的用户密码
	slicePiece.Master = "172.17.0.2:3306" // 主实例地址
	slicePiece.Slaves = []string{}        // 从实例地址列表
	slicePiece.StatisticSlaves = nil      // 统计型从实例地址列表
	slicePiece.Capacity = 12              // gaea_proxy与每个实例的连接池大小
	slicePiece.MaxCapacity = 24           // gaea_proxy与每个实例的连接池最大大小
	slicePiece.IdleTimeout = 60           // gaea_proxy与后端mysql空闲连接存活时间，单位:秒

	// >>>>> 载入物理设定档 (Load2)
	modelsNameSpace.Slices = append(modelsNameSpace.Slices, &slicePiece)

	// >>>>> 载入分片设定档 (Load3)
	modelsNameSpace.ShardRules = nil // 分库、分表、特殊表的配置内容，具体字段可参照shard配置 (载入设定档)

	// >>>>> 组织用户配置
	modelsNameSpace.Users = []*models.User{} // 应用端连接gaea所需要的用户配置，具体字段可参照users配置
	userPiece := models.User{}
	userPiece.UserName = "root"              // 用户名
	userPiece.Password = "12345"             // 用户密码
	userPiece.Namespace = "env1_namespace_1" // 对应的命名空间
	userPiece.RWFlag = 2                     // 读写标识
	userPiece.RWSplit = 1                    // 是否读写分离
	userPiece.OtherProperty = 0              // 其他属性，目前用来标识是否走统计从实例

	// >>>>> 载入用户设定档 (Load4)
	modelsNameSpace.Users = append(modelsNameSpace.Users, &userPiece)

	// >>>>> 载入预设值设定档 (Load5)
	modelsNameSpace.DefaultSlice = "slice-0" // 预设分片名称
	modelsNameSpace.GlobalSequences = nil    // 生成全局唯一序列号的配置, 具体字段可参考全局序列号配置
	modelsNameSpace.DefaultCharset = ""      // 用于指定数据集如何排序，以及字符串的比对规则
	modelsNameSpace.DefaultCollation = ""    // 用于指定数据集如何排序，以及字符串的比对规则
	modelsNameSpace.MaxSqlExecuteTime = 0    // sql最大执行时间，大于该时间，进行熔断
	modelsNameSpace.MaxSqlResultSize = 0     // 限制单分片返回结果集大小不超过max_select_rows

	// 初始化管理员 (Manager0)
	serverManager := new(Manager)             // 服务器管理员
	current := 0                              // 切换标签
	namespaceManager := new(NamespaceManager) // NameSpace 管理员
	namespaceManager.namespaces = make(map[string]*Namespace, 64)
	managerNamespace := &Namespace{} // 管理员的 NameSpace
	defer managerNamespace.Close(false)

	// 组成管理员的 NameSpace (对应到 Load1)
	managerNamespace.openGeneralLog = modelsNameSpace.OpenGeneralLog                     // 记录sql查询的访问日志，说明 https://github.com/XiaoMi/Gaea/issues/109
	managerNamespace.name = "env1_namespace_1"                                           // namespace 为划分工作业务的最基本单位，一个 namespace 可以有多个使用者
	managerNamespace.allowedDBs = make(map[string]bool, len(modelsNameSpace.AllowedDBS)) // 数据库列表
	for db, allowed := range modelsNameSpace.AllowedDBS {
		managerNamespace.allowedDBs[strings.TrimSpace(db)] = allowed
	}
	defaultPhyDBs := make(map[string]string, len(modelsNameSpace.DefaultPhyDBS)) // 预设数据库列表
	for db, phyDB := range modelsNameSpace.DefaultPhyDBS {
		defaultPhyDBs[strings.TrimSpace(db)] = strings.TrimSpace(phyDB)
	}
	managerNamespace.defaultPhyDBs, _ = parseDefaultPhyDB(defaultPhyDBs, managerNamespace.allowedDBs)
	managerNamespace.slowSQLTime, _ = parseSlowSQLTime(modelsNameSpace.SlowSQLTime) // 慢sql阈值，单位: 毫秒
	managerNamespace.sqls = make(map[string]string, 16)
	managerNamespace.sqls = parseBlackSqls(modelsNameSpace.BlackSQL)        // 黑名单sql
	managerNamespace.allowips, _ = parseAllowIps(modelsNameSpace.AllowedIP) // 白名单IP

	// 组成管理员的 NameSpace (对应到 Load2)
	managerNamespace.slices, _ = parseSlices(modelsNameSpace.Slices, managerNamespace.defaultCharset, managerNamespace.defaultCollationID) // 一主多从的物理实例，slice里map的具体字段可参照slice配置

	// 组成管理员的 NameSpace (对应到 Load4)
	managerNamespace.userProperties = make(map[string]*UserProperty, 2)
	for _, user := range modelsNameSpace.Users {
		up := &UserProperty{RWFlag: user.RWFlag, RWSplit: user.RWSplit, OtherProperty: user.OtherProperty}
		managerNamespace.userProperties[user.UserName] = up
	}

	// 组成管理员的 NameSpace (对应到 Load5)
	sequences := sequence.NewSequenceManager() // 生成全局唯一序列号的配置, 具体字段可参考全局序列号配置
	for _, v := range modelsNameSpace.GlobalSequences {
		globalSequenceSlice, _ := managerNamespace.slices[v.SliceName]
		seqName := strings.ToUpper(v.DB) + "." + strings.ToUpper(v.Table)
		seq := sequence.NewMySQLSequence(globalSequenceSlice, seqName, v.PKName)
		sequences.SetSequence(v.DB, v.Table, seq)
	}
	managerNamespace.sequences = sequences
	managerNamespace.defaultCharset, managerNamespace.defaultCollationID, _ = parseCharset(modelsNameSpace.DefaultCharset, modelsNameSpace.DefaultCollation) // 用于指定数据集如何排序，以及字符串的比对规则 & 用于指定数据集如何排序，以及字符串的比对规则
	if modelsNameSpace.MaxSqlExecuteTime <= 0 {                                                                                                              // sql最大执行时间，大于该时间，进行熔断
		managerNamespace.maxSqlExecuteTime = defaultMaxSqlExecuteTime
	} else {
		managerNamespace.maxSqlExecuteTime = modelsNameSpace.MaxSqlExecuteTime
	}
	if modelsNameSpace.MaxSqlResultSize <= 0 { // 限制单分片返回结果集大小不超过max_select_rows
		managerNamespace.maxSqlResultSize = defaultMaxSqlResultSize
	} else {
		managerNamespace.maxSqlResultSize = modelsNameSpace.MaxSqlResultSize
	}

	// 组成管理员的 NameSpace (延伸部份)
	managerNamespace.slowSQLCache = cache.NewLRUCache(defaultSQLCacheCapacity)
	managerNamespace.errorSQLCache = cache.NewLRUCache(defaultSQLCacheCapacity)
	managerNamespace.backendSlowSQLCache = cache.NewLRUCache(defaultSQLCacheCapacity)
	managerNamespace.backendErrorSQLCache = cache.NewLRUCache(defaultSQLCacheCapacity)
	managerNamespace.planCache = cache.NewLRUCache(defaultPlanCacheCapacity)

	// 组成管理员的 NameSpace (建立路由)
	managerNamespace.router, _ = router.NewRouter(modelsNameSpace)

	// 把 管理員之NS 分别写回 NS 和 服务器管理员
	namespaceManager.namespaces[modelsNameSpace.Name] = managerNamespace
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
	executor.namespace = "env1_namespace_1"

	// 进行测试
	serverNameSpace := executor.GetNamespace()
	fmt.Println(serverNameSpace)
}
