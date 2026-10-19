#!/bin/sh
# 349-v1 门禁四数（只在共享工作树；跑前 unset GOTMPDIR TMPDIR TEMP TMP）
set -u
REPO='/d/work/workspace/projects plans/Wisp'
cd "$REPO" || exit 90
unset GOTMPDIR TMPDIR TEMP TMP
export GOPATH='D:\work\base\gopath'
L=/d/tmp/wisp/20261019/349v1/logs
mkdir -p "$L"
echo '### gate 1: gofumpt -l（本腿写面＝.md/.txt/.log/.sh，⛔ .go；三枚 risk 产码件按名打，验交件字节干净）'
/d/work/base/gopath/bin/gofumpt.exe -l internal/risk/blacklist.go internal/risk/provenance.go internal/risk/syncdirs.go > "$L/g1-gofumpt-worktree.txt" 2>&1
echo "gofumpt_worktree_rc=$? offenders_below"
cat "$L/g1-gofumpt-worktree.txt"
echo '### gofumpt against blob form (LF) — my bytes verdict'
for f in blacklist provenance syncdirs; do
  git cat-file blob "HEAD:internal/risk/$f.go" > "$L/blob-$f.go"
  /d/work/base/gopath/bin/gofumpt.exe -l "$L/blob-$f.go" > "$L/g1b-$f.txt" 2>&1
  echo "blob $f rc=$? listing=$(wc -c < "$L/g1b-$f.txt")bytes"
done
echo '### gate 2: d22scan'
sh scripts/d22scan.sh > "$L/g2-d22scan.txt" 2>&1
echo "d22scan_rc=$?"
tail -6 "$L/g2-d22scan.txt"
echo '### gate 3: check-path-length-budget'
sh scripts/check-path-length-budget.sh > "$L/g3-pathlen.txt" 2>&1
echo "pathlen_rc=$?"
tail -4 "$L/g3-pathlen.txt"
echo '### gate 4: go vet ./internal/risk/'
/d/work/base/go/bin/go.exe vet ./internal/risk/ > "$L/g4-vet.txt" 2>&1
echo "vet_rc=$?"
cat "$L/g4-vet.txt"
echo '### named: worktree dirty state of internal/risk at this moment'
git status --porcelain -- internal/risk | head -10
echo "RC=0"
