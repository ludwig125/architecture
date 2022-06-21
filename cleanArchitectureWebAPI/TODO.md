# TODO

- 環境変数から切り替えできるように
- packageを分ける?
- mysql からsqliteに置き換えてみる
- sqliteとmysqlはテストは一緒、セットアップ部分は分ける


[~] $curl -X POST -H "Content-Type: application/json" -d '{"id": 4, "name":"Jackson", "age":64}' 'http://localhost:8080/update' -w '%{http_code}\n'
