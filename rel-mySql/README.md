## mysqlサーバー起動
```shell
mysql.server start
```
## ユーザー名をgo_example、ポート番号3306にして、ログイン(パスワード: 12345!)
```mysql 
mysql -h localhost --port 3306 -u go_example -p12345!;
```
## ログアウト
```mysql
exit
```
## mysqlサーバー停止
```shell
mysql.server stop
```