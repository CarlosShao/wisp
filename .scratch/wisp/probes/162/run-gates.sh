#!/usr/bin/env bash
# 162-r1 台件：票 162 门禁跑法（名册差集＝比"名字集合"，不比包级 rc）。
#
# 为什么要有这一份：一条用例 panic 会吞掉同包其余读数，包级 ok/fail 看不出来。
# 基线锚＝5e83808（本程进场时的 HEAD，未含票 162 的任何改动）。
#
# 用法：bash .scratch/wisp/probes/162/run-gates.sh
set -uo pipefail
cd "$(git rev-parse --show-toplevel)"
export PATH="$PATH:$(go env GOPATH)/bin"
OUT=.scratch/wisp/probes/162
BASE=5e83808

echo "### 1) d22scan 自测（tools/d22scan 模块全量）"
bash tools/d22scan/runtests.sh -C tools/d22scan ./... | tail -2
echo "rc=$?"

echo "### 2) 全仓 d22scan 门"
sh scripts/d22scan.sh | tail -1
echo "rc=$?"

echo "### 3) go test ./internal/tools/ 名册（-json，取 run 名）"
go test -count=1 -json ./internal/tools/ > "$OUT/package-run.jsonl" 2> "$OUT/package-run.stderr"
echo "rc=$? (stderr 行数 $(wc -l < "$OUT/package-run.stderr"))"
grep -o '"Test":"[^"]*"' "$OUT/package-run.jsonl" | sed 's/.*:"//;s/"$//' | grep -v '/' | sort -u > "$OUT/roster-now.txt"
echo "现在的名册：$(wc -l < "$OUT/roster-now.txt") 枚顶层用例"

echo "### 4) 基线名册（$BASE 的源码级名册）"
git grep -h -E '^func (Test|Benchmark)[A-Z]' "$BASE" -- internal/tools \
  | sed -E 's/^func ([A-Za-z0-9_]+).*/\1/' | sort -u > "$OUT/roster-base.txt"
echo "基线名册：$(wc -l < "$OUT/roster-base.txt") 枚"

echo "### 5) 名册两向差集（左＝基线有而现在没了，右＝现在多了）"
comm -23 "$OUT/roster-base.txt" "$OUT/roster-now.txt" | sed 's/^/  丢了: /'
comm -13 "$OUT/roster-base.txt" "$OUT/roster-now.txt" | sed 's/^/  新增: /'

echo "### 6) 失败与 panic 读数"
grep -c '"Action":"fail"' "$OUT/package-run.jsonl" | sed 's/^/  fail 行数: /'
grep -c '"Action":"panic"' "$OUT/package-run.jsonl" | sed 's/^/  panic 行数: /' || true

echo "### 7) vet / gofumpt 甲乙两形"
go vet ./internal/tools/ && echo "  vet: ok"
git ls-files "*.go" > "$OUT/tracked-go.txt"
xargs -a "$OUT/tracked-go.txt" gofumpt -l > "$OUT/gofumpt-A-tracked.txt"; echo "  甲（已跟踪 $(wc -l < "$OUT/tracked-go.txt") 个 .go）不干净行数: $(wc -l < "$OUT/gofumpt-A-tracked.txt")"
git ls-files --others --exclude-standard "*.go" > "$OUT/untracked-go.txt"
cat "$OUT/tracked-go.txt" "$OUT/untracked-go.txt" | xargs gofumpt -l > "$OUT/gofumpt-B-all.txt"
echo "  乙（已跟踪＋未跟踪 $(cat "$OUT/tracked-go.txt" "$OUT/untracked-go.txt" | wc -l) 个）不干净清单:"
sed 's/^/    /' "$OUT/gofumpt-B-all.txt"
gofumpt --version | sed 's/^/  gofumpt 版本: /'
