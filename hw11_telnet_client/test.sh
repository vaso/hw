#!/usr/bin/env bash
set -xeuo pipefail

go build -o go-telnet
nc_out_file=$(mktemp)
telnet_out_file=$(mktemp)

(echo -e "Hello\nFrom\nNC\n" && cat 2>/dev/null) | nc -l localhost 4242 > ${nc_out_file} &
NC_PID=$!

sleep 1
(echo -e "I\nam\nTELNET client\n" && cat 2>/dev/null) | ./go-telnet --timeout=5s localhost 4242 > ${telnet_out_file} &
TL_PID=$!

sleep 5
kill ${TL_PID} 2>/dev/null || true
kill ${NC_PID} 2>/dev/null || true

function fileEquals() {
  local fileData
  fileData=$(cat "$1")
  [ "${fileData}" = "${2}" ] || (echo -e "unexpected output, $1:\n${fileData}" && exit 1)
}

expected_nc_out='I
am
TELNET client'
fileEquals ${nc_out_file} "${expected_nc_out}"

expected_telnet_out='Hello
From
NC'
fileEquals ${telnet_out_file} "${expected_telnet_out}"

rm -f go-telnet
echo "PASS"
