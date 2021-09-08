package server

import (
	"encoding/json"
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

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> 1 台 Master 2 台 Slave 数据库测试

/*
                   +----------------+        +----------------+
192.168.122.2:3307 | 从数据库 db0-0   |        | 从数据库 db0-1   | 192.168.122.2:3308
                   +---------\------+        +-------/--------+
                              --\                /---
                                 --\          /--
                                 +--------------+
              192.168.122.2:3306 | 主数据库 db0   |
                                 +--------------+

   第一本小说 三国演义
   第二本小说 水浒传
   第三本小说 西游记
   第四本小说 红楼梦
   第五本小说 金瓶梅
   第六本小说 儒林外史
   第七本小说 初刻拍案惊奇
   第八本小说 二刻拍案惊奇
   第九本小说 封神演义
   第十本小说 镜花缘
   第十一本小说 喻世明言
   第十二本小说 说岳全传
   第十三本小说 杨家将
   第十四本小说 说唐
   第十五本小说 七侠五义
   第十六本小说 施公案
   第十七本小说 青楼梦
   第十八本小说 歧路灯
   第十九本小说 老残游记
   第二十本小说 二十年目睹之怪现状
   第二十一本小说 孽海花
   第二十二本小说 官场现形记
   第二十三本小说 觉世名言十二楼
   第二十四本小说 无声戏
   第二十五本小说 肉蒲团
   第二十六本小说 浮生六记
   第二十七本小说 野叟曝言
   第二十八本小说 九尾龟
   第二十九本小说 品花宝鉴
*/

// prepareDb0NamespaceManagerForCluster 函式 产生针对 Cluster (db0 db0-0 db0-1) 的设定档
func prepareDb0NamespaceManagerForCluster() (*Manager, error) {
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
  "name": "db0_cluster_namespace",
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
      "master": "192.168.122.2:3306",
      "slaves": ["192.168.122.2:3307", "192.168.122.2:3308"],
      "statistic_slaves": null,
      "capacity": 12,
      "max_capacity": 24,
      "idle_timeout": 60
    }
  ],
  "shard_rules": null,
  "users": [
    {
      "user_name": "panhong",
      "password": "12345",
      "namespace": "db0_cluster_namespace",
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
	namespaceName := "db0_cluster_namespace"
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

// prepareDb0PlanSessionExecutorForCluster 函式 产生针对 Cluster (db0 db0-0 db0-1) 的 Plan Session
func prepareDb0PlanSessionExecutorForCluster() (*SessionExecutor, error) {
	var userName = "panhong"
	var namespaceName = "db0_cluster_namespace"
	var database = "novel"

	m, err := prepareDb0NamespaceManagerForCluster()
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

// TestDb0PlanExecuteIn 函式 为向 Cluster (db0 db0-0 db0-1) 图书馆数据库查询 29 本小说
// 测试分二版，分别为连到数据库的版本和不连到数据库的版本，此版本会连到数据库
func TestDb0PlanExecuteIn(t *testing.T) {
	// 载入 Session Executor
	se, err := prepareDb0PlanSessionExecutorForCluster()
	require.Equal(t, err, nil)
	db, err := se.GetNamespace().GetDefaultPhyDB("novel")
	require.Equal(t, err, nil)
	require.Equal(t, db, "novel") // 检查 SessionExecutor 是否正确载入

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

		// 初始化单元测试程式
		backend.MarkTakeOver()

		// 执行数据库分库指令
		res, err := p.ExecuteIn(reqCtx, se)
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
