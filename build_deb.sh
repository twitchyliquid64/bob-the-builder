#!/bin/bash
set -e

SCRIPT_BASE_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"
BUILD_DIR=${1%/}

if [ -f "${BUILD_DIR}" ]; then
	echo "Cannot build into '${BUILD_DIR}': it is a file."
	exit 1
fi

if [ -d "${BUILD_DIR}" ]; then
  rm -rfv ${BUILD_DIR}/*
fi

mkdir -pv ${BUILD_DIR}

mkdir -pv "${BUILD_DIR}/usr/share/bob-the-builder"
go build -o "${BUILD_DIR}/usr/share/bob-the-builder/bob-the-builder" *.go
cp -rv templates "${BUILD_DIR}/usr/share/bob-the-builder"
cp -rv static "${BUILD_DIR}/usr/share/bob-the-builder"

mkdir -pv "${BUILD_DIR}/DEBIAN"
cp -rv ${SCRIPT_BASE_DIR}/DEBIAN/* "${BUILD_DIR}/DEBIAN"

ARCH=`dpkg --print-architecture`
sed -i "s/ARCH/${ARCH}/g" "${BUILD_DIR}/DEBIAN/control"

mkdir -pv "${BUILD_DIR}/etc/bob-the-builder"
cat > ${BUILD_DIR}/etc/bob-the-builder/config.json << "EOF"
{
	"name": "Build Server",

	"Web": {
		"Domain": "localhost",
		"Listener": "0.0.0.0:8010"
	},

	"AWS": {
		"Enable": true,
		"AccessKey": "PUT_UR_STUFF_IN_HERE",
		"SecretKey": "PUT_UR_STUFF_IN_HERE"
	}
}
EOF

mkdir -pv "${BUILD_DIR}/lib/systemd/system/"
cat > ${BUILD_DIR}/lib/systemd/system/bob-the-builder.service << "EOF"
[Unit]
Description=Bob the builder automation service.

[Service]
Type=simple
ExecStart=/usr/share/bob-the-builder/bob-the-builder --definitions_dir /var/lib/bob-the-builder/definitions --base_dir /var/lib/bob-the-builder/base --build_dir /var/lib/bob-the-builder/build /etc/bob-the-builder/config.json
WorkingDirectory=/usr/share/bob-the-builder
KillMode=control-group
TimeoutStopSec=5s
Restart=always
RestartSec=15s
IgnoreSIGPIPE=no
StandardOutput=syslog
StandardError=syslog
SyslogIdentifier=bob-the-builder
EOF

mkdir -pv ${BUILD_DIR}/var/lib/bob-the-builder/{definitions,base,build}

dpkg-deb --build "${BUILD_DIR}" ./
