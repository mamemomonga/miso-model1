# MISO MODEL1

Raspberry Pi をベースとしたハードウェア専用です。

# /boot/config.txtの設定

以下の内容を /boot/config.txt に追記してください。

	dtoverlay=gpio-poweroff,gpiopin=21,active_low="y"
	dtoverlay=gpio-shutdown,gpio_pin=20
	dtparam=audio=on
	gpu_mem=16
	dtparam=spi=on

