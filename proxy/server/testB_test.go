package server

import (
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
	// 再扩充
} 