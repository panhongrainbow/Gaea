# 如何把 Gaea 切换成读写分离 

> 这一篇文章主要在说明如何把 Gaea 切换成读写分离的模式，从设定档到程式码到中断点都有讨论到

## 1 目前测试环境的说明

目前的测试环境只用一个 Cluster db0 进行测试

<img src="../../../Gaea/docs/assets/panhongrainbow/image-20210720040119317.png" style="zoom:60%;" /> 

在 Cluster db0 的内部分配，db0 为可读可写 Master 数据库节点，db0-0 为唯读 Slave 数据库节点

## 2 在设定档设定读写分离

在设定档里，可以经由 rw_split 参数设定读写分离

1. 当 rw_split 参数的值为 0 时，为 "非" 读写分离的状态
2. 当 rw_split 参数的值为 1 时，为 读写分离的状态

````go
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
      "user_name": "docker",
      "password": "12345",
      "master": "192.168.1.2:3350",
      "slaves": ["192.168.1.2:3351", "192.168.1.2:3352"],
      "statistic_slaves": null,
      "capacity": 12,
      "max_capacity": 24,
      "idle_timeout": 60
    }
  ],
  "shard_rules": null,
  "users": [
    {
      "user_name": "root",
      "password": "12345",
      "namespace": "db0_cluster_namespace",
      "rw_flag": 2,
      "rw_split": 1, # <<<<<<<<<<<<<<<<<<<<<<<<<<<<<< 在这里使用参数 rw_split 去设定读写分离
      "other_property": 0
    }
  ],
  "default_slice": "slice-0",
  "global_sequences": null
}`
````

目前设定参数 rw_split 的值为 1 ，为启用 MariaDB 数据库读写分离

## 3 在程式码设定读写分离

在单元测试时，必须加入几行测试码，以达到读写分离的状态

```go
func TestB3(t *testing.T) {
	// 载入 Session Executor
	se, err := preparePlanSessionExecutorForCluster()
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
		// 执行 SQL Parser
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
		reqCtx.Set(util.FromSlave, 1) // >>> 在这里设定读取时从 Slave 节点，达到读写分离的效果
		result, err := p.ExecuteIn(reqCtx, se)
		require.Equal(t, err, nil)

		fmt.Println(result)
	}
}
```

关键的程式码为 reqCtx.Set(util.FromSlave, 1)，此程式码会告知单元测试读取资料时从 Slave 节点

## 4 设定中断点去观察直达数据库网路位置

> Gaea 因为有定义 interface 接口，会间接跳到其他函式，先在单元测试内设定中断点，再慢慢追踪到直连的网路位置

设定第 1 个中断点位置，在 result, err := p.ExecuteIn(reqCtx, se) 这一行设定中断点

<img src="../../../Gaea/docs/assets/panhongrainbow/image-20210724054323328.png" alt="image-20210724054323328" style="zoom:100%;" /> 

设定第 2 个中断点位置

- 在 Gaea/proxy/server/executor.go 档案里
- pcs, err := se.getBackendConns(sqls, getFromSlave(reqCtx)) 这一行设定中断点

<img src="../../../Gaea/docs/assets/panhongrainbow/image-20210724055217604.png" alt="image-20210724055217604" style="zoom:100%;" /> 

主要是观察 pcs 变数，pcs 变数内含数据库直接网路位置

## 5 在中断点查寻读写分离的结果

到 pcs, err := se.getBackendConns(sqls, getFromSlave(reqCtx)) 这一行中断点执行后，展开 pcs 变数

<img src="../../../Gaea/docs/assets/panhongrainbow/image-20210724060310741.png" alt="image-20210724060310741" style="zoom:100%;" />

把 pcs 变数展开，可以发现数据库直连网路位置为 192.168.1.2:3351，此为 db0-0 Slave 数据库的网路位置，证明状态已经为读写分离
