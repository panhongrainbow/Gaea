package server

import (
	"github.com/XiaoMi/Gaea/models"
	// "github.com/XiaoMi/Gaea/stats"
	"github.com/XiaoMi/Gaea/util"
	"github.com/XiaoMi/Gaea/util/sync2"
	// "github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/require"
	"gopkg.in/ini.v1"
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

	namespaceConfig := new(models.Namespace)
	namespaceConfig.OpenGeneralLog = false
	namespaceConfig.IsEncrypt = false
	namespaceConfig.Name = "env1_namespace_1"
	namespaceConfig.Online = true
	namespaceConfig.ReadOnly = false
	namespaceConfig.AllowedDBS = make(map[string]bool)
	namespaceConfig.AllowedDBS["Library"] = true
	namespaceConfig.DefaultPhyDBS = nil
	namespaceConfig.SlowSQLTime = "1000"
	namespaceConfig.BlackSQL = []string{}
	namespaceConfig.AllowedIP = nil
	namespaceConfig.Slices = []*models.Slice{}
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
	namespaceConfig.Slices = append(namespaceConfig.Slices, &tmp1)
	namespaceConfig.ShardRules = nil
	namespaceConfig.Users = []*models.User{}
	tmp2 := models.User{}
	tmp2.UserName = "root"
	tmp2.Password = "12345"
	tmp2.Namespace = "env1_namespace_1"
	tmp2.RWFlag = 2
	tmp2.RWSplit = 1
	tmp2.OtherProperty = 0
	namespaceConfig.Users = append(namespaceConfig.Users, &tmp2)
	namespaceConfig.DefaultSlice = "slice-0"
	namespaceConfig.GlobalSequences = nil
	namespaceConfig.DefaultCharset = ""
	namespaceConfig.DefaultCollation = ""
	namespaceConfig.MaxSqlExecuteTime = 0
	namespaceConfig.MaxSqlResultSize = 0

	statisticManager := StatisticManager{}
	statisticManager.manager = new(Manager)
	statisticManager.manager.reloadPrepared = sync2.NewAtomicBool(false)
	tmp3 := util.BoolIndex{}
	tmp3.Set(true)
	statisticManager.manager.switchIndex = tmp3
	tmp4 := NamespaceManager{}
	tmp5 := NamespaceManager{}
	statisticManager.manager.namespaces[0] = &tmp4
	statisticManager.manager.namespaces[1] = &tmp5
	tmp6 := UserManager{}
	tmp7 := UserManager{}
	statisticManager.manager.users[0] = &tmp6
	statisticManager.manager.users[1] = &tmp7
	statisticManager.manager.statistics = nil // 里面有更多内容
	statisticManager.clusterName = "gaea_default_cluster"
	statisticManager.statsType = ""
	/*statisticManager.handlers["/metrics"] = promhttp.Handler() // 暂时，未确认
	tmp8, _ := initGeneralLogger(proxy)
	statisticManager.generalLogger = tmp8 // 以下先使用函式，之后在分析
	statisticManager.sqlTimings = stats.NewMultiTimings("SqlTimings",
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
	// statisticManager.backendSQLErrorCounts
	// statisticManager.backendSQLFingerprintErrorCounts
	// statisticManager.backendConnectPoolIdleCounts
	// statisticManager.backendConnectPoolInUseCounts
	statisticManager.backendConnectPoolWaitCounts = stats.NewGaugesWithMultiLabels("backendConnectPoolWaitCounts",
			"gaea proxy backend wait connect counts", []string{statsLabelCluster, statsLabelNamespace, statsLabelSlice, statsLabelIPAddr})
	statisticManager.slowSQLTime = 100
	statisticManager.closeChan = make(chan bool, 0)
	*/
}
 
