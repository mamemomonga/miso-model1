#!/bin/bash
set -eu
BASEDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

COMMANDS="install uninstall"

function do_install {
    echo "*** INSTALL ***"

	cat > /etc/systemd/system/power-controller-startup.service << EOS
[Unit]
Description=Power Controller startup
DefaultDependencies=no
After=rc-local.service plymouth-start.service systemd-user-sessions.service

[Service]
Type=oneshot
ExecStart=$BASEDIR/power-controller -action startup -pulse

[Install]
WantedBy=rc-local.service plymouth-start.service systemd-user-sessions.service
EOS

	cat > /etc/systemd/system/power-controller-shutdown.service << EOS
[Unit]
Description=Power Controller shutdown
DefaultDependencies=no
Before=systemd-poweroff.service

[Service]
Type=oneshot
ExecStart=$BASEDIR/power-controller -action shutdown

[Install]
WantedBy=systemd-poweroff.service
EOS

	cat > /etc/systemd/system/misoapp.service << EOS
[Unit]
Description=MisoApp
DefaultDependencies=no
After=power-controller-startup.service

[Service]
Type=simple
ExecStart=$BASEDIR/misoapp

[Install]
WantedBy=power-controller-startup.service
EOS

	set -x
	systemctl enable power-controller-startup.service
	systemctl enable power-controller-shutdown.service
	systemctl enable misoapp.service
	set +x

}

function do_uninstall {
    echo "*** UNINSTALL ***"
	set -x
	systemctl disable power-controller-startup.service
	systemctl disable power-controller-shutdown.service
	systemctl disable misoapp.service

	rm -f /etc/systemd/system/power-controller-startup.service
	rm -f /etc/systemd/system/power-controller-shutdown.service
	rm -f /etc/systemd/system/misoapp.service
	set +x
}

function run {

    if [ "$(id -u)" != "0" ]; then
    exec sudo $0 $@

    fi
    for i in $COMMANDS; do
    if [ "$i" == "${1:-}" ]; then
        shift
        do_$i $@
        exit 0
    fi
    done
    echo "USAGE: $( basename $0 ) COMMAND"
    echo "COMMANDS:"
    for i in $COMMANDS; do
    echo "   $i"
    done
    exit 1
}

run $@
