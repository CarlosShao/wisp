#!/bin/sh
set -e
cd "/d/work/workspace/projects plans/Wisp/.scratch/wisp/probes/orch/20261018/txt"
cat 322v2-head.txt common-20261018.txt 322v2-specific.txt > 322v2-dispatch.txt
echo "322v2-dispatch = $(wc -l < 322v2-dispatch.txt) lines / $(wc -c < 322v2-dispatch.txt) bytes"
mkdir -p "/d/work/workspace/projects plans/Wisp/.scratch/wisp/probes/322/v2/logs"
echo "dir-ok: $(ls -d "/d/work/workspace/projects plans/Wisp/.scratch/wisp/probes/322/v2/logs")"
