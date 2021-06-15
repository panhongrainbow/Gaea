package server

import (
	"fmt"
	"github.com/XiaoMi/Gaea/backend"
	"github.com/XiaoMi/Gaea/models"
	"github.com/XiaoMi/Gaea/mysql"
	"github.com/XiaoMi/Gaea/parser"
	"github.com/XiaoMi/Gaea/parser/format"
    // 先忽略 stats 包，因为这个测试会跟其他测试发生冲突
	// "github.com/XiaoMi/Gaea/stats"
	"github.com/XiaoMi/Gaea/util"
	"github.com/XiaoMi/Gaea/util/sync2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/require"
	"gopkg.in/ini.v1"
	"net/http"
	"strings"
	"testing"
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
		{"cluster_name", "gaea_default_cluster",},
		{"slow_sql_time", "100",},
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
	tmp3, _ := NewNamespace(modelsNameSpace)
	nsMgr.namespaces[modelsNameSpace.Name] = tmp3
	serverManager.namespaces[current] = nsMgr

	// user, _ := CreateUserManager(namespaceMap)
	user := new(UserManager)
	user.users = make(map[string][]string, 64)
	user.userNamespaces = make(map[string]string, 64)
	// user.addNamespaceUsers(modelsNameSpace)
	user.userNamespaces["root" + ":" + "12345"] = "env1_namespace_1"
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
			"", // 期望的 SQL 字串关键字
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
 
