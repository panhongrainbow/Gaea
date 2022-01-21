# Gaea 数据库中间件连线启动说明

> 此为 store.md 的辅助说明文件，当 store.go 把读取设定文档的接口制定后，数据库中间件就可以读取设定值正常启动

## 1 读取环境说明

先略过，因为要画架构图

## 2 设定档读取方式

> 目前支援设定档的读取方式为
>
> 1. 方法一：使用 File 文档去储存设定值
> 2. 方法二：使用 Etcd V2 API 去储存设定值
> 3. 方法三：使用 Etcd V3 API 去储存设定值

### 1 File 文档储存设定值

> 当准备使用 File 文档去储存设定值时，需要修改改两个设定文档
>
> 1. 初始化设定文档，位于 Gaea/etc/gaea.ini
> 2. 命名空间设定文档，集中于目录 Gaea/etc/file/namespace/

修正初始化设定档 Gaea/etc/gaea.ini

```ini
; 这里的重点在把 config_type 值改成 file !!!!!

; config type, etcd/file, you can test gaea with file type, you shoud use etcd in production
config_type=file
;file config path, 具体配置放到file_config_path的namespace目录下，该下级目录为固定目录
file_config_path=./etc/file

; 以下略过，因为重点要把前面的 config_type 设定值改成 file
```

在目录 Gaea/etc/file/namespace 下新增一个 namespace 设定档，档名为 novel_cluster_namespace.json，关于 小说数据库丛集 的相关设定

丛集设定值内容整理成下表，命名空间的名称为 novel_cluster_namespace

| 丛集编号 |    Mater 服务器    |   Slave 服务器一   |   Slave 服务器二   |  帐号  | 密码  |
| :------: | :----------------: | :----------------: | :----------------: | :----: | :---: |
|  丛集1   | 192.168.122.2:3309 | 192.168.122.2:3310 | 192.168.122.2:3311 | xiaomi | 12345 |
|  丛集2   | 192.168.122.2:3312 | 192.168.122.2:3313 | 192.168.122.2:3314 | xiaomi | 12345 |

切片设定值内容整理成下表

| 丛集编号 | 切片名称 | 数据表名称 | 是否预设 |
| :------: | :------: | :--------: | :------: |
|  丛集1   | slice-0  | Book_0000  |    是    |
|  丛集2   | slice-1  | Book_0001  |    否    |

演算法设定值整理成下表

| 演算法设定项目 | 演算法设定值          | 说明                                                         |
| -------------- | --------------------- | ------------------------------------------------------------ |
| 用户名称       | hash                  | kingshard hash分片演算法                                     |
| 分表依据的键值 | BookID                | 会以 BookID 的数值作为分表的依据                             |
| 数据表数量     | [1,1]                 | 阵列 [1,1] 分别指出每一个切片的数据表数量，比如<br />slice-0 有 1 个 数据表，<br />slice-1 有 1 个 数据表 |
| 切片列表阵列   | ["slice-0","slice-1"] | 这个命名空间里，有两个切片，分别为 slice-0 和 slice-1        |

命名空间用户设定值整理成下表

| 用户设定项目   | 用户设定值                             | 说明                                                         |
| -------------- | -------------------------------------- | ------------------------------------------------------------ |
| 用户名称       | xiaomi                                 |                                                              |
| 用户密码       | 12345                                  |                                                              |
| 命名空间的名称 | novel_cluster_namespace                |                                                              |
| 用户读写标记   | rw_flag 为 2 ，该用户可进行读写操作    | rw_flag 为 1，只能进行 唯读 操作<br />rw_flag 为 2，可进行 读写 操作 |
| 读写分离标记   | rw_split 为 1，该用户进行 读写分离操作 | rw_split 为 0，进行 非读写 分离操作<br />rw_split 为 1，进行 读写 分离操作 |

整份设定文档内容如下

```json
{
  "name": "novel_cluster_namespace",
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
      "user_name": "xiaomi",
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
      "user_name": "xiaomi",
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
      "user_name": "xiaomi",
      "password": "12345",
      "namespace": "novel_cluster_namespace",
      "rw_flag": 2,
      "rw_split": 1,
      "other_property": 0
    }
  ],
  "default_slice": "slice-0",
  "global_sequences": null
}
```

### 2 使用 File 设定文档 启动 Gaea



### 3 Etcd API 储存设定值

目前测试的 key 有两组

/gaea/namespace/novel_cluster_namespace

/gaea_default_cluster/namespace/novel_cluster_namespace

登入 Gaea 的指令为

mysql -h 127.0.0.1 -P 13306 --protocol=TCP -u xiaomi -p

目前测试的资料如下

MySQL [novel]> INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(1, 9781517191276, 'Romance Of The Three Kingdoms', 'Luo Guanzhong', 1522, 'Hi
storical fiction'); 

MySQL [novel]> INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(2, 9789869442060, 'Water Margin', 'Shi Nai an', 1589, 'Historical fiction'); 

MySQL [novel]> INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(3, 9789575709518, 'Journey To The West', 'Wu Cheng en', 1592, 'Gods And Demon
s Fiction'); 

先简单记录，再仔细处理

### 4 使用 Etcd API 启动 Gaea























