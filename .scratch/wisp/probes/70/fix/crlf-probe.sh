#!/usr/bin/env bash
# Ticket 70 probe: where does scripts/portable-tests.sh carry CRLF? Used to explain
# a standing-carrier reading on this (autocrlf=true) checkout, not to change bytes.
set -u
f=$1
total=$(wc -l <"$f")
crlf=$(awk '/\r$/ {n++} END {printf "%d", n+0}' "$f")
printf '%s: total=%s crlf_lines=%s\n' "$f" "$total" "$crlf"
awk '/\r$/ {print NR}' "$f" | head -5 | tr '\n' ' '
printf '\n'
