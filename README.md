# MISO MODEL1

Raspberry Pi をベースとしたハードウェア専用ソフトウェアです。

# ディレクトリ

ディレクトリ | 概要
-------------|------------
avr    | ATTiny13A用ソースコード
bin    | ユーティリティー
etc    | マストドン接続設定および設定サンプル
go     | Goアプリケーション
script | インストーラ
sounds | 音声ファイル

# /boot/config.txtの設定

以下の内容を /boot/config.txt に追記してください。

	dtoverlay=gpio-poweroff,gpiopin=21,active_low="y"
	dtoverlay=gpio-shutdown,gpio_pin=20
	dtparam=audio=on
	gpu_mem=16
	dtparam=spi=on

