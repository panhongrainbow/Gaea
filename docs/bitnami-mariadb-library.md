# 准备 Bitnami-MariaDB 图书馆数据库测试环境

> 先做一个简单的记录

```bash
# 准备数据库测试环境，先在这个 Schema 进行扩充，之后还会准备其他不同类型较复杂的 Schema

# 远端登入虚机
$ ssh docker@192.168.122.2

# 登入 Master 容器的资料库，Master 是用 3306 埠
$ mysql -h 127.0.0.1 -P 3306 --protocol=TCP -u docker -p

# 创建资料库 Library，这一段可省略，因为 Bitnami-MariaDB 会自动创建资料库
$ # MariaDB [(none)]> CREATE DATABASE `Library` /*!40100 DEFAULT CHARACTER SET utf8 */;

# 使用图书馆资料库
$ MariaDB [(none)]> USE Library;

# 建立书资料表
$ CREATE TABLE `Book` (
  `BookID` int(11) NOT NULL,
  `Isbn` bigint(20) NOT NULL,
  `Title` varchar(30) NOT NULL,
  `Author` varchar(30) DEFAULT NULL,
  `Publish` int(4) DEFAULT NULL,
  `Category` varchar(30) NOT NULL,
  PRIMARY KEY (`BookID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

# 新增第一笔书的资料
$ INSERT INTO Library.Book
(BookID, Isbn, Title, Author, Publish, Category)
VALUES(1, 9789865975364, 'Dream of the Red Chamber', 'Cao Xueqin', 1791, 'Family Saga');

# 登入 Slave 容器的资料库，Slave 是用 3307 埠 或 3308 埠
$ mysql -h 127.0.0.1 -P 3307 --protocol=TCP -u docker -p

# 使用图书馆资料库
$ MariaDB [(none)]> USE Library;

# 在 Slave 容器可以正确读到 Master 写入的资料
$ SELECT * FROM Book;
#+--------+---------------+------------------------+------------+---------+-------------+
#| BookID | Isbn          | Title                  | Author     | Publish | Category    |
#+--------+---------------+------------------------+------------+---------+-------------+
#|      1 | 9789865975364 | The Red Chamber        | Cao Xueqin |    1791 | Family Saga |
#+--------+---------------+------------------------+------------+---------+-------------+
#1 row in set (0.001 sec) 
```

