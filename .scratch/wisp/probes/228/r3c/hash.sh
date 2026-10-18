#!/bin/sh
# 228-r3c: hash ruler. Print sha256 (export-tree physical bytes) and git blob
# sha1 (content, LF-normalised) of the mutated file. Read the blob sha1 from the
# INDEX via `git hash-object`, which is content-only and thus immune to whether
# the file on disk happens to carry CRLF or LF.
set -u
F="${1:-/d/tmp/wisp/20261018/228r3c/tree/internal/proc/boot_windows.go}"
echo "file=$F"
printf 'sha256bytes='
sha256sum "$F" | cut -d' ' -f1
printf 'blobsha1='
git hash-object "$F"
printf 'bytes='
wc -c <"$F"
printf 'lines='
wc -l <"$F"
