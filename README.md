[![LICENSE](https://img.shields.io/badge/license-Apache--2.0-blue.svg)](https://github.com/XiaoMi/Gaea/blob/master/LICENSE)
[![Build Status](https://travis-ci.org/XiaoMi/Gaea.svg?branch=master)](https://travis-ci.org/XiaoMi/Gaea)
[![Go Report Card](https://goreportcard.com/badge/github.com/XiaoMi/Gaea)](https://goreportcard.com/report/github.com/XiaoMi/Gaea)

## 简介

Gaea是小米中国区电商研发部研发的基于mysql协议的数据库中间件，目前在小米商城大陆和海外得到广泛使用，包括订单、社区、活动等多个业务。Gaea支持分库分表、sql路由、读写分离等基本特性，更多详细功能可以参照下面的功能列表。其中分库分表方案兼容了mycat和kingshard两个项目的路由方式。Gaea在设计、实现阶段参照了mycat、kingshard和vitess，并使用tidb parser作为内置的sql parser，在此表达诚挚感谢。为了方便使用和学习Gaea，我们也提供了详细的使用和设计文档，也欢迎大家多多参与。

## 功能列表

#### 基础功能

- 多集群
- 多租户
- SQL透明转发
- 慢SQL指纹
- 错误SQL指纹
- 注解路由
- 慢日志
- 读写分离，从库负载均衡
- 自定义SQL拦截与过滤
- 连接池
- 配置热加载
- IP/IP段白名单
- 全局序列号

#### 分库、分表功能

- 分库: 支持mycat分库方式
- 分表: 支持kingshard分表方式
- 聚合函数: 支持max、min、sum、count、group by、order by等
- join: 支持分片表和全局表的join、支持多个分片表但是路由规则相同的join

## 架构图

![gaea架构图](docs/assets/architecture.png)

## 集群部署图  

![gaea集群部署图](docs/assets/deployment.png)  

如上图所示, 部署一套gaea-cc和etcd可用来管理多套gaea集群, 负责集群内namespace配置的增删改查.
[gaea-cc的HTTP接口文档](docs/gaea-cc.md)

## 安装使用

- [快速入门](docs/quickstart.md)
- [配置说明](docs/configuration.md)
- [监控配置说明](docs/grafana.md)
- [全局序列号配置说明](docs/sequence-id.md)
- [基本概念](docs/concepts.md)
- [SQL兼容性](docs/compatibility.md)
- [FAQ](docs/faq.md)

## 设计与实现

- [整体架构](docs/architecture.md)
- [多租户的设计与实现](docs/multi-tenant.md)
- [gaea配置热加载设计与实现](docs/config-reloading.md)
- [gaea proxy后端连接池的设计与实现](docs/connection-pool.md)
- [prepare的设计与实现](docs/prepare.md)

## 开发方式

- [🏁开发当下的测试](docs/teststart.md)
- [🏁图书馆实体数据库测试环境](docs/bitnami-mariadb-novel.md)
- [🚫初始化 JetBrain GoLand IDE 工具](docs/panhongrainbow/run-goland-gaea.md)
- [🚫设定 JetBrain GoLand IDE 权限](docs/panhongrainbow/permission-goland-gaea.md)
- [🏁使用 JetBrain GoLand IDE 进行开发](docs/goland-develop.md)
- [🏁中文文件繁简互换](docs/chinese-translate.md)
- [🚫保存程式码副本](docs/panhongrainbow/preserve-data.md)
- [🚫程式码日常维护](docs/panhongrainbow/maintain-golang-gaea.md)
- [🏁程式码 GoFmt 格式化维护](docs/gofmt-golang-gaea.md)
- [🏁触发单元测试](docs/goland-gaea-unit-test.md)

## 开发进入点

记录开发的进入位罝是否要把单元测试做成一个包	

| 项目 | 位置                                                      | 说明                              |
| ---- | --------------------------------------------------------- | --------------------------------- |
| A    | github.com/panhongrainbow/Gaea/parser/testA_test.go       | Sql Parser 转换                   |
| B    | github.com/panhongrainbow/Gaea/proxy/server/testB_test.go | 把 设定档 和 SQL字串转成 直连命令 |
| C    | github.com/panhongrainbow/Gaea/backend/testC_test.go      | 和 MariaDB 之间的交界             |
| D    |                                                           |                                   |
| E    |                                                           |                                   |

## 单元测试所接管的函式

| 项目 | 档案位置                          | 资料 struct      |
| ---- | --------------------------------- | ---------------- |
| 1    | Gaea/backend/direct_connection.go | DirectConnection |
| 2    | Gaea/proxy/plan/plan_unshard.go   | UnshardPlan      |
| 3    |                                   |                  |
| 4    |                                   |                  |
| 5    |                                   |                  |

如果想要知道资料是如何被拦截的，可以对 IsTakeOver 函式去搜寻 Find Usage

## 开发日志

> 这里记录开发所考量问题并做记录，说明在进行决择的过程
> 可用以下指令去更新日志
>
> ```bash
> # 进入日志资料夹
> $ cd /home/panhong/go/src/github.com/panhongrainbow/Gaea/docs/diary
> 
> # 进行繁简转换
> $ opencc -i 20210817.t.md -o 20210817.md -c t2s.json
> ```

| 日期           | 标题                       | 连结                                |
| -------------- | -------------------------- | ----------------------------------- |
| 2021年08月17日 | 是否要把单元测试做成一个包 | [🐱日志连结](docs/diary/20210817.md) |
| 2021年09月07日 | 发现切片规则可能不会被触发 | [🐱日志连结](docs/diary/20210907.md) |
| 2021年09月08日 | 调整假资料库载入资料的方式 | [🐱日志连结](docs/diary/20210908.md) |
|                |                            |                                     |
|                |                            |                                     |

## Roadmap

- [x] 支持配置加密存储，开关
- [ ] 支持执行计划缓存
- [ ] 支持事务追踪
- [ ] 支持二级索引
- [ ] 支持分布式事务
- [ ] 支持平滑的扩容、缩容
- [ ] 后端连接池优化 (按照请求时间排队)

## 自有开发模块

- backend  
- cmd  
- log  
- models  
- proxy/plan  
- proxy/router(kingshard路由方式源自kingshard项目本身)  
- proxy/sequence
- server  

## 外部模块

- mysql(google vitess、tidb、kingshard都有引入)  
- parser(tidb)  
- stats(google vitess，打点统计)  
- util(混合)

## 社区

### gitter
[![Gitter](https://badges.gitter.im/xiaomi-b2c/Gaea.svg)](https://gitter.im/xiaomi-b2c/Gaea?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)

### 钉钉
![钉钉](docs/assets/gaea_dingtalk.png)
