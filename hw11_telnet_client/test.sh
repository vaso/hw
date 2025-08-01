#!/usr/bin/env bash
set -xeuo pipefail

go build -o go-telnet
rm -f nc.out
rm -f telnet.out

(echo -e "Hello\nFrom\nNC\n" && cat 2>/dev/null) | nc -l localhost 4242 >nc.out &
NC_PID=$!

sleep 1
(echo -e "I\nam\nTELNET client\n" && cat 2>/dev/null) | ./go-telnet --timeout=5s localhost 4242 >telnet.out &
TL_PID=$!

sleep 5
kill ${TL_PID} 2>/dev/null || true
kill ${NC_PID} 2>/dev/null || true

function fileEquals() {
  local fileData
  fileData=$(cat "$1")
  [ "${fileData}" = "${2}" ] || (echo -e "unexpected output, $1:\n${fileData}:\n${2}" && ls -lsa && cat "$1" && exit 1)
}

expected_nc_out=$'I\nam\nTELNET client'
fileEquals nc.out "${expected_nc_out}"

expected_telnet_out=$'Hello\nFrom\nNC'
fileEquals telnet.out "${expected_telnet_out}"

rm -f go-telnet
echo "PASS"
