# Gaea 數據庫中間件連線啟動說明

> 此為 store.md 的輔助說明文件，當 store.go 把讀取設定文檔的接口制定後，數據庫中間件就可以讀取設定值正常啟動

## 1 讀取環境說明

先略過，因為要畫架構圖

## 2 設定檔讀取方式

> 目前支援設定檔的讀取方式為
>
> 1. 方法一：使用 File 文檔去儲存設定值
> 2. 方法二：使用 Etcd V2 API 去儲存設定值
> 3. 方法三：使用 Etcd V3 API 去儲存設定值

### 1 File 文檔儲存設定值

> 當準備使用 File 文檔去儲存設定值時，需要修改改兩個設定文檔
>
> 1. 初始化設定文檔，位於 Gaea/etc/gaea.ini
> 2. 命名空間設定文檔，集中於目錄 Gaea/etc/file/namespace/

修正初始化設定檔 Gaea/etc/gaea.ini

```ini
; 這裡的重點在把 config_type 值改成 file !!!!!

; config type, etcd/file, you can test gaea with file type, you shoud use etcd in production
config_type=file
;file config path, 具体配置放到file_config_path的namespace目录下，该下级目录为固定目录
file_config_path=./etc/file

; 以下略過，因為重點要把前面的 config_type 設定值改成 file
```

在目錄 Gaea/etc/file/namespace 下新增一個 namespace 設定檔，檔名為 novel_cluster_namespace.json，關於 小說數據庫叢集 的相關設定

叢集設定值內容整理成下表，命名空間的名稱為 novel_cluster_namespace

| 叢集編號 |    Mater 服務器    |   Slave 服務器一   |   Slave 服務器二   |  帳號  | 密碼  |
| :------: | :----------------: | :----------------: | :----------------: | :----: | :---: |
|  叢集1   | 192.168.122.2:3309 | 192.168.122.2:3310 | 192.168.122.2:3311 | xiaomi | 12345 |
|  叢集2   | 192.168.122.2:3312 | 192.168.122.2:3313 | 192.168.122.2:3314 | xiaomi | 12345 |

切片設定值內容整理成下表

| 叢集編號 | 切片名稱 | 數據表名稱 | 是否預設 |
| :------: | :------: | :--------: | :------: |
|  叢集1   | slice-0  | Book_0000  |    是    |
|  叢集2   | slice-1  | Book_0001  |    否    |

演算法設定值整理成下表

| 演算法設定項目 | 演算法設定值          | 說明                                                         |
| -------------- | --------------------- | ------------------------------------------------------------ |
| 用戶名稱       | hash                  | kingshard hash分片演算法                                     |
| 分表依據的鍵值 | BookID                | 會以 BookID 的數值作為分表的依據                             |
| 數據表數量     | [1,1]                 | 陣列 [1,1] 分別指出每一個切片的數據表數量，比如<br />slice-0 有 1 個 數據表，<br />slice-1 有 1 個 數據表 |
| 切片列表陣列   | ["slice-0","slice-1"] | 這個命名空間裡，有兩個切片，分別為 slice-0 和 slice-1        |

命名空間用戶設定值整理成下表

| 用戶設定項目   | 用戶設定值                             | 說明                                                         |
| -------------- | -------------------------------------- | ------------------------------------------------------------ |
| 用戶名稱       | xiaomi                                 |                                                              |
| 用戶密碼       | 12345                                  |                                                              |
| 命名空間的名稱 | novel_cluster_namespace                |                                                              |
| 用戶讀寫標記   | rw_flag 為 2 ，該用戶可進行讀寫操作    | rw_flag 為 1，只能進行 唯讀 操作<br />rw_flag 為 2，可進行 讀寫 操作 |
| 讀寫分離標記   | rw_split 為 1，該用戶進行 讀寫分離操作 | rw_split 為 0，進行 非讀寫 分離操作<br />rw_split 為 1，進行 讀寫 分離操作 |

整份設定文檔內容如下

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

### 2 使用 File 設定文檔 啟動 Gaea



### 3 Etcd API 儲存設定值

目前測試的 key 有兩組

/gaea/namespace/novel_cluster_namespace

/gaea_default_cluster/namespace/novel_cluster_namespace

登入 Gaea 的指令為

mysql -h 127.0.0.1 -P 13306 --protocol=TCP -u xiaomi -p

目前測試的資料如下

MySQL [novel]> INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(1, 9781517191276, 'Romance Of The Three Kingdoms', 'Luo Guanzhong', 1522, 'Hi
storical fiction'); 

MySQL [novel]> INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(2, 9789869442060, 'Water Margin', 'Shi Nai an', 1589, 'Historical fiction'); 

MySQL [novel]> INSERT INTO novel.Book (BookID, Isbn, Title, Author, Publish, Category) VALUES(3, 9789575709518, 'Journey To The West', 'Wu Cheng en', 1592, 'Gods And Demon
s Fiction'); 

先簡單記錄，再仔細處理

### 4 使用 Etcd API 啟動 Gaea























