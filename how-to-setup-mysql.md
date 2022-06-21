
## mysqlのインストール

以下、Ubuntuの場合
```bash
$ sudo apt install mysql-server mysql-client

$mysql --version
mysql  Ver 14.14 Distrib 5.7.30, for Linux (x86_64) using  EditLine wrapper
```

起動に使用するmysqlユーザーのホームディレクトリが存在しないとmysql serverを立ち上げられないので/etc/passwdに以下を追加

```bash
$ sudo usermod -d /var/lib/mysql mysql
```

/etc/passwdに以下が追加されている

```bash
mysql:x:111:115:MySQL Server,,,:/var/lib/mysql:/bin/false
```

mysql server の起動

```bash
$ sudo service mysql start
```

ubuntu18.04でデフォルトのmysql5.7ではroot権限でないと接続できない
これだと個人開発環境では不便なので、sudoいらなくさせる

```bash
mysql > ALTER USER 'root'@'localhost' IDENTIFIED WITH mysql_native_password BY '';
mysql > FLUSH PRIVILEGES;
```

これでsudoもpasswordも不要になる


## mysql でのテーブルの作成まで

#### 接続

```
$mysql -u root --host 127.0.0.1 --port 3306
Welcome to the MySQL monitor.  Commands end with ; or \g.
Your MySQL connection id is 4
Server version: 5.7.33-0ubuntu0.18.04.1 (Ubuntu)

Copyright (c) 2000, 2021, Oracle and/or its affiliates.

Oracle is a registered trademark of Oracle Corporation and/or its
affiliates. Other names may be trademarks of their respective
owners.

Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.
mysql>
```

#### database確認

```
mysql> SHOW databases;
+--------------------+
| Database           |
+--------------------+
| information_schema |
| mysql              |
| performance_schema |
| sys                |
+--------------------+
5 rows in set (0.00 sec)

mysql>
```

#### database作成

```
CREATE DATABASE IF NOT EXISTS arch_db;
```


#### table作成

```
CREATE TABLE IF NOT EXISTS arch_db.actor(
		id INT AUTO_INCREMENT,
        name VARCHAR(30),
        age INT,
        PRIMARY KEY (id)
);
```

確認

```
mysql> use arch_db;
Reading table information for completion of table and column names
You can turn off this feature to get a quicker startup with -A

Database changed
mysql> show tables;
+-----------------+
| Tables_in_arch_db |
+-----------------+
| actor       |
+-----------------+
1 row in set (0.00 sec)

mysql>
mysql> desc actor;
+-------+-------------+------+-----+---------+----------------+
| Field | Type        | Null | Key | Default | Extra          |
+-------+-------------+------+-----+---------+----------------+
| id    | int(11)     | NO   | PRI | NULL    | auto_increment |
| name  | varchar(30) | YES  |     | NULL    |                |
| age   | int(11)     | YES  |     | NULL    |                |
+-------+-------------+------+-----+---------+----------------+
3 rows in set (0.01 sec)

mysql>
```

#### insert test data
```
INSERT INTO actor(name, age) VALUES("Johannsen", 30);
INSERT INTO actor(name, age) values("Williams", 53);
INSERT INTO actor(name, age) values("Streep", 65);
```

#### 後で削除するときはDROPで

```
DROP DATABASE arch_db;
```
