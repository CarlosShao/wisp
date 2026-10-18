#!/bin/sh
set -e
cd "/d/work/workspace/projects plans/Wisp/.scratch/wisp/probes/orch/20261018/txt"
cat 336v1-head.txt common-20261018.txt 336v1-specific.txt > 336v1-dispatch.txt
echo "336v1-dispatch = $(wc -l < 336v1-dispatch.txt) lines / $(wc -c < 336v1-dispatch.txt) bytes"
mkdir -p "/d/work/workspace/projects plans/Wisp/.scratch/wisp/probes/336/v1/logs"
echo "dir-ok"
