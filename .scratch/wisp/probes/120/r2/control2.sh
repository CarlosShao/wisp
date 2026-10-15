#!/bin/sh
# 120-r2 control, second cut: gate2/control.sh aborted on its log copy (the snapshot
# is a git archive, and .scratch is gitignored, so /work/.scratch does not exist),
# and set -e skipped the restore. This script restores first, then re-does the
# mutation with the log redirected on the host side, then restores again and proves
# the hash. Nothing here touches the repository's own working tree.
set -e
S=/d/tmp/probe120/r2/base-fixed/internal/winsec
F="$S/winsec_other.go"
BAK="$S/winsec_other.go.bak-120r2-control"
OUT="D:/work/workspace/projects plans/Wisp/.scratch/wisp/probes/120/r2"

echo "== state at entry =="
echo "sha256 current =$(sha256sum "$F" | cut -d' ' -f1)"
echo "sha256 bak     =$(sha256sum "$BAK" | cut -d' ' -f1)"
cp "$BAK" "$F"
CLEAN=$(sha256sum "$F" | cut -d' ' -f1)
echo "sha256 restored=$CLEAN"

echo "== re-mutate: sameObjectPinned neutered on disk =="
python - <<'PYEOF'
import io
p = r'D:/tmp/probe120/r2/base-fixed/internal/winsec/winsec_other.go'
s = io.open(p, encoding='utf-8').read()
start = s.index('func sameObjectPinned(')
end = s.index('\n}\n', start) + 3
body = s[start:end]
a = 'if now.Mode()&os.ModeSymlink != 0 {'
b = 'if !os.SameFile(pinned, now) {'
assert body.count(a) == 1 and body.count(b) == 1
body = body.replace(a, 'if now.Mode()&os.ModeSymlink != 0 && false {')
body = body.replace(b, 'if !os.SameFile(pinned, now) && false {')
io.open(p, 'w', encoding='utf-8', newline='\n').write(s[:start] + body + s[end:])
print('neutered')
PYEOF
MUT=$(sha256sum "$F" | cut -d' ' -f1)
echo "sha256 mutated =$MUT"

echo "== control run against those bytes (strict default, no env) =="
set +e
MSYS_NO_PATHCONV=1 docker run --rm -e CGO_ENABLED=0 -v /d/tmp/probe120/r2/base-fixed:/work -v probe120build:/root/.cache/go-build -w /work golang:1.27 \
	sh -c 'go test -count=1 -v -run "TestPOSIX120" ./internal/winsec/' > "$OUT/120-control-gate-neutered.log" 2>&1
echo "control rc=$?"
set -e
echo "failing names:"
grep -E '^--- FAIL|^    --- FAIL' "$OUT/120-control-gate-neutered.log" || true
echo "red rows quoted with the文案:"
grep -E 'AC#1 RED' "$OUT/120-control-gate-neutered.log" || true

echo "== restore and prove =="
cp "$BAK" "$F"
AFTER=$(sha256sum "$F" | cut -d' ' -f1)
echo "sha256 after restore=$AFTER"
[ "$CLEAN" = "$AFTER" ] && echo "RESTORE MATCHES" || echo "RESTORE MISMATCH"
echo "== the repository's own file: never part of this control =="
cd "D:/work/workspace/projects plans/Wisp"
echo "repo sha256=$(sha256sum internal/winsec/winsec_other.go | cut -d' ' -f1)"
echo "snapshot-clean sha256=$CLEAN"
git status --porcelain -- internal/winsec/winsec_other.go
