#!/bin/sh
# 120-r2 control run (反控): the identity gate in winsec_other.go is turned into a
#恒真 by physical rewrite on disk (NOT -overlay, which only replaces compile-time
# inputs and is blind to anything that reads the tree at runtime), the named cases
# are run against those bytes, then the file is restored and both hashes are printed.
# The .bak is kept, not deleted (临时件只建不删).
set -e
S=/d/tmp/probe120/r2/base-fixed/internal/winsec
F="$S/winsec_other.go"
OUT="D:/work/workspace/projects plans/Wisp/.scratch/wisp/probes/120/r2"

BEFORE=$(sha256sum "$F" | cut -d' ' -f1)
cp "$F" "$S/winsec_other.go.bak-120r2-control"
echo "sha256 before      =$BEFORE"
echo "sha256 of the .bak =$(sha256sum "$S/winsec_other.go.bak-120r2-control" | cut -d' ' -f1)"

python - <<'PYEOF'
import io
p = r'D:/tmp/probe120/r2/base-fixed/internal/winsec/winsec_other.go'
s = io.open(p, encoding='utf-8').read()
# Scope the rewrite to sameObjectPinned only: that is the by-name identity gate the
# three cases drive. sealHandle's copy is left standing (and is reported as having
# no named case, which is this leg's honest gap).
start = s.index('func sameObjectPinned(')
end = s.index('\n}\n', start) + 3
body = s[start:end]
a = 'if now.Mode()&os.ModeSymlink != 0 {'
b = 'if !os.SameFile(pinned, now) {'
assert body.count(a) == 1 and body.count(b) == 1
body = body.replace(a, 'if now.Mode()&os.ModeSymlink != 0 && false {')
body = body.replace(b, 'if !os.SameFile(pinned, now) && false {')
s = s[:start] + body + s[end:]
io.open(p, 'w', encoding='utf-8', newline='\n').write(s)
print('sameObjectPinned neutered on disk (sealHandle leg untouched)')
PYEOF

MUT=$(sha256sum "$F" | cut -d' ' -f1)
echo "sha256 mutated     =$MUT"
grep -n '&& false' "$F"

echo "== control run: the same three cases against the neutered bytes =="
MSYS_NO_PATHCONV=1 docker run --rm -e CGO_ENABLED=0 -v /d/tmp/probe120/r2/base-fixed:/work -v probe120build:/root/.cache/go-build -w /work golang:1.27 \
	sh -c 'go test -count=1 -v -run "TestPOSIX120" ./internal/winsec/ > /tmp/control.log 2>&1; echo "control rc=$?"; grep -c "^--- FAIL" /tmp/control.log; sed -n "s/^--- FAIL: \(Test[A-Za-z0-9_]*\).*/FAILING: \1/p" /tmp/control.log | sort | uniq; grep -E "AC#1 RED sample" /tmp/control.log; cp /tmp/control.log /work/.scratch/wisp/probes/120/r2/120-control-gate-neutered.log'

cp /tmp/none /dev/null 2>/dev/null || true
echo "== restoring =="cp "$S/winsec_other.go.bak-120r2-control" "$F"
AFTER=$(sha256sum "$F" | cut -d' ' -f1)
echo "sha256 after restore=$AFTER"
if [ "$BEFORE" = "$AFTER" ]; then echo "RESTORE MATCHES (snapshot byte-identical to the pre-control file)"; else echo "RESTORE MISMATCH"; fi
echo "== the repository's own file was never touched by this control =="
cd "D:/work/workspace/projects plans/Wisp"
git status --porcelain -- internal/winsec/winsec_other.go
sha256sum internal/winsec/winsec_other.go | cut -d' ' -f1
