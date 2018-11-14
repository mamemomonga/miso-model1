# go-rqh-webserver

* HTTP Request Headerを表示するシンプルなウェブサーバです。
* dep, packr, html/templateの習作

# 実行例

	$ ./rqh-webserver-darwin-amd64 -listen 0.0.0.0:8000

# 開発

## ビルド環境

事前に必要なもの

* make
* go

導入されてない場合導入されるもの
		
* dep
* packr

### 準備

* GOPATHが正しく設定されている必要があります。
* make deps を実行すると、dep eusure が実行されます。dep, packrが導入されていない場合は導入されます。
* linux-armの場合、go getでdepを取得します。

コマンド

	$ mkdir -p $GOPATH/src/github.com/mamemomonga/go-rqh-webserver
	$ git clone https://github.com/mamemomonga/go-rqh-webserver $GOPATH/src/github.com/mamemomonga/go-rqh-webserver
	$ make deps

実行

	$ make run

ビルド

	$ make

実行

	$ bin/rqh-webserver -listen 0.0.0.0:8000

### リリース向けビルド

* 公開用のバイナリが生成されます。
* Dockerが必要です。

コマンド

	$ make release

# 別の実行方法

depによる依存関係が無視されて導入されます

	$ go get -u -v github.com/mamemomonga/go-rqh-webserver/src/rqh-webserver
	$ rqh-webserver -listen 0.0.0.0:8000


