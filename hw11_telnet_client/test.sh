#!/usr/bin/env bash
set -xeuo pipefail

go build -o go-telnet
nc_out=nc.out
telnet_out=telnet.out
port=$(shuf -i 40000-50000 -n 1)

rm -f ${nc_out}
rm -f ${telnet_out}

(echo -e "Hello\nFrom\nNC\n" && cat 2>/dev/null) | nc -l localhost ${port} > ${nc_out} &
NC_PID=$!

sleep 1
(echo -e "I\nam\nTELNET client\n" && cat 2>/dev/null) | ./go-telnet --timeout=5s localhost ${port} > ${telnet_out} &
TL_PID=$!

sleep 5
kill ${TL_PID} 2>/dev/null || true
kill ${NC_PID} 2>/dev/null || true

function fileEquals() {
  local fileData
  fileData=$(cat "$1")
  [ "${fileData}" = "${2}" ] || (echo -e "unexpected output, $1:\n${fileData}:\n${2}" && ls -lsa $1 && cat "$1" && exit 1)
}

expected_nc_out=$'I\nam\nTELNET client'
fileEquals ${nc_out} "${expected_nc_out}"

expected_telnet_out=$'Hello\nFrom\nNC'
fileEquals ${telnet_out} "${expected_telnet_out}"

rm -f go-telnet
rm -f ${nc_out}
rm -f ${telnet_out}

echo "PASS"
