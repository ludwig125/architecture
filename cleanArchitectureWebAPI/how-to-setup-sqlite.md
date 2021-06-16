## sqliteのインストール

以下、Ubuntuの場合
```
sudo apt -y install sqlite3

$sqlite3 -version
3.31.1 2020-01-27 19:55:54 3bfa9cc97da10598521b342961df8f5f68c7388fa117345eeb516eaa837balt1
```

## sqlite table作成まで

#### db作成

注：cleanArchitectureWebAPIディレクトリ直下で行う

```
sqlite3 arch_db
```

#### db確認とtable 作成
```
$sqlite3 arch_db
SQLite version 3.31.1 2020-01-27 19:55:54
Enter ".help" for usage hints.
sqlite> .databases
main: /home/ludwig125/go/src/github.com/ludwig125/architecture/cleanArchitectureWebAPI/arch_db
sqlite>
```

以下でtable作成
```
CREATE TABLE actor(id INTEGER PRIMARY KEY ASC, name TEXT, age INTEGER);
```

確認

```
sqlite> .tables
actor

sqlite> .schema actor
CREATE TABLE actor(id INTEGER PRIMARY KEY ASC, name TEXT, age INTEGER);
```

#### insert test data

dataの例
```
INSERT INTO actor(name, age) VALUES("Portman", 32);
INSERT INTO actor(name, age) values("Knightley", 35);
INSERT INTO actor(name, age) values("Hopkins", 56);
```

確認
```
sqlite> select * from actor;
1|Portman|32
2|Knightley|35
3|Hopkins|56
sqlite>
```

#### おまけ
dbとtableの作成とデータのInsertは以下のように一気にすることもできる

```
sqlite3 arch_db 'CREATE TABLE actor(id INTEGER PRIMARY KEY ASC, name TEXT, age INTEGER);'
sqlite3 arch_db 'INSERT INTO actor(name, age) VALUES("Portman", 32);'
sqlite3 arch_db 'INSERT INTO actor(name, age) values("Knightley", 35);'
sqlite3 arch_db 'INSERT INTO actor(name, age) values("Hopkins", 56);'
sqlite3 arch_db 'INSERT INTO actor(name, age) values("Depp", 54);'
sqlite3 arch_db 'INSERT INTO actor(name, age) values("Watson", 24);'
```
