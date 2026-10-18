#!/bin/bash
# 227-v1 AC#6 read-only gauges (committed tree only; NO go test in cmd/wisp — that needs the lane)
cd '/d/work/workspace/projects plans/Wisp'
out='.scratch/wisp/probes/227/v1/logs/07-ac6-readings.txt'
{
echo "gauge = git show HEAD:<path> + grep (算的: 命中枚数; 读的: 逐字句). tree = HEAD ($(git rev-parse --short HEAD))"
echo
echo "[1] the sentence under AC#6: always_write_no_clobber_226_test.go:43-44 verbatim"
git show HEAD:cmd/wisp/always_write_no_clobber_226_test.go | sed -n '43,44p'
echo "hits of 'nothing polls' across the line break (content anchor, not line number):"
git show HEAD:cmd/wisp/always_write_no_clobber_226_test.go | tr '\n' ' ' | tr -s ' \t' ' ' | grep -c 'running and nothing polls:'
echo "negative control (same gauge, phrase that must NOT exist): $(git show HEAD:cmd/wisp/always_write_no_clobber_226_test.go | tr '\n' ' ' | tr -s ' \t' ' ' | grep -c 'ZZZ not a phrase') (rc=1 => zero)"
echo
echo "[2] is the claim 'nothing polls' true in production bytes at HEAD?"
git grep -n 'configReloadPollInterval = ' HEAD -- cmd/wisp/config_reload.go
git grep -n 'time.NewTicker(configReloadPollInterval)' HEAD -- cmd/wisp/config_reload.go
git grep -n 'startConfigReload' HEAD -- cmd/wisp | sed -n '1,8p'
echo
echo "[3] which commit wired it (票面 says 248095d1) - is 248095d1 an ancestor of HEAD?"
git merge-base --is-ancestor 248095d191ad6b97b1d4af4066be0f44be2598b5 HEAD && echo "ANCESTOR=yes" || echo "ANCESTOR=no rc=$?"
git show -s --format='%h %ad %s' --date=format:'%m-%d %H:%M' 248095d191ad6b97b1d4af4066be0f44be2598b5 | cut -c1-120
echo
echo "[4] AC#6's other half: does ANY resident *_test.go in cmd/wisp pin 'hand edit still visible under the real 1s tick'?"
echo "  (question: a permanent case that lets the production ticker run and asserts a hand-edited key is seen)"
git grep -n -E 'CheckAndReload' HEAD -- 'cmd/wisp/*_test.go' | cut -c1-160
echo "  ---- how those call it: manual in-test calls vs the wired ticker ----"
git grep -n -E 'startConfigReload|reloadHandle|configReloadTick' HEAD -- 'cmd/wisp/*_test.go' | cut -c1-160
echo
echo "[5] the :102 assertion direction (票面 AC#6 第二条要求) verbatim"
git show HEAD:cmd/wisp/always_write_no_clobber_226_test.go | sed -n '97,105p'
echo
echo "[6] the 09-29 correction the ticket carries: 4.82s/1.31s readings are NOT mine (cited as 票面 text) - gauge on ticket bytes:"
grep -c '4.82' .scratch/wisp/issues/227-guarded-write-delete-shapes-and-migration-on-read-are-unasserted.md
grep -n '4.82' .scratch/wisp/issues/227-guarded-write-delete-shapes-and-migration-on-read-are-unasserted.md | cut -c1-80
} > "$out" 2>&1
echo "written $out bytes=$(wc -c < "$out")"
