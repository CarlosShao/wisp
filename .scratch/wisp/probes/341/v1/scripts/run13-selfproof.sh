#!/usr/bin/env bash
# 341-v1 交件自证三把尺（全部我自己现跑）
R="/d/work/workspace/projects plans/Wisp"; cd "$R" || exit 9
MINE="95b1015d 67596f91 7d91f761 f5d5165c"
LEG="6948eb08 fd82f786 f710d104 920537ee 46c920c8 cef2a6c9"
V="docs/evidence/s1/341-v1-verdict.md"

echo "########## 尺一：逐枚提交号 + 逐枚 git show --name-only（应⊆ 三处写面），.go 命中 ===== 0 ##########"
for h in $MINE; do
  echo "--- $h ---"
  git show --name-only --format='   subject=%s' "$h" | cut -c1-110 | sed 's/^/  /'
done
echo "-- my 4 commits: .go hits (want 0) --"
git show --name-only --format='' $MINE | grep -c '\.go$' || true
echo "-- POSITIVE CONTROL: same ruler on 341-r1's 6 commits (want 1) --"
git show --name-only --format='' $LEG | sort -u | grep '\.go$' | sed 's/^/     carrier: /'
echo "-- my commits ⊆ allowed面? list anything else (want empty) --"
git show --name-only --format='' $MINE | sort -u | grep -vE '^$' \
 | grep -vE '^docs/evidence/s1/341-v1-verdict\.md$|^\.scratch/wisp/probes/341/v1/|^\.scratch/wisp/issues/341-' | sed 's/^/   OUTSIDE: /'
echo "-- roster scale (strict prefix per-hash, want 4) --"
git log --all --grep='票 341 终裁腿 341-v1·' --format='%h %ad %s' --date=format:'%H:%M' | cut -c1-70
echo "-- count = $(git log --all --grep='票 341 终裁腿 341-v1·' --format='%H' | wc -l)"

echo
echo "########## 尺二：判语件绝对路径 + wc -c + 行数 + 占位符尺 0 命中（正控非零）+ 逐枚非 0 字节 ##########"
echo "ABS = $(cygpath -w "$R/$V")"
echo "bytes=$(wc -c < "$V") lines=$(wc -l < "$V")"
echo "placeholder ruler (TODO/FIXME/XXX) hits in verdict = $(grep -c 'TODO\|FIXME\|XXX' "$V" || true)"
echo "   POSITIVE CONTROL same ruler on HANDOVER.md = $(grep -c 'TODO\|FIXME\|XXX' docs/reports/HANDOVER.md || true) (XXX 1 hit = ruler alive)"
echo "banned-char ruler hits in verdict = $(grep -cP '\x{00A7}' "$V" || true)"
echo "   POSITIVE CONTROL same ruler on HANDOVER.md files = $(grep -rlP '\x{00A7}' docs/reports/HANDOVER.md | wc -l)"
echo "CR in verdict = $(grep -c "$(printf '\r')" "$V" || true)"
echo "-- every file in my write面: zero-byte count (want 0) --"
find .scratch/wisp/probes/341/v1 -type f -size 0 | wc -l
find .scratch/wisp/probes/341/v1 -type f | wc -l
echo "-- .out extension count (want 0) --"; find .scratch/wisp/probes/341/v1 -name '*.out' | wc -l
echo "-- per-file bytes --"; find .scratch/wisp/probes/341/v1 -type f | sort | while read x; do printf "   %-72s %s\n" "$x" "$(wc -c < "$x")"; done

echo
echo "########## 尺三：决定性读数汇总（四格成对红句逐字 + 还原哈希 + 四数 + 门禁 rc + 票框终量尺）##########"
echo "-- AC#4 成对四数（我的两棵从零导出树，两把临时根） --"
grep -A4 '^###' .scratch/wisp/probes/341/v1/logs/03-color-base.txt | grep -E 'tree=|RUN=|PASS_TOP=|FAIL_TOP=|SKIP=' | sed 's/^/   /'
grep -E 'tree=|^RUN=|^PASS_TOP=|^FAIL_TOP=|^SKIP=' .scratch/wisp/probes/341/v1/logs/04-color-deliv.txt | sed 's/^/   /'
echo "-- 四发去牙：四数 + 目标⇄正控逐名 + 还原哈希 --"
for x in p2 p3 p5 p6; do
  echo "   ===== $x ====="
  grep -E '^RUN=|^PASS_TOP=|^FAIL_TOP=|^SKIP=|^-- old anchor|^-- new anchor|^-- survival|^-- restored hash' .scratch/wisp/probes/341/v1/logs/10-mu-$x.txt | sed 's/^/     /'
  grep -A2 '^-- per-case verdict' .scratch/wisp/probes/341/v1/logs/10-mu-$x.txt | tail -1 | sed 's/^/     /'
  sed -n '/^-- per-case verdict for the five/,/^-- verbatim red/p' .scratch/wisp/probes/341/v1/logs/10-mu-$x.txt | grep -E 'Test' | sed 's/^/     /'
  echo "     red roster: $(tr '\n' ' ' < /d/tmp/wisp/20261018/341v1/log-$x/rednames.txt)"
done
echo "-- 门禁四数（共享工作树，跑前 unset GOTMPDIR TMPDIR TEMP TMP） --"
grep -E 'GOFUMPPT_RC=|VET_RC=|D22_RC=|PATHLEN_RC=' .scratch/wisp/probes/341/v1/logs/08-gates-and-ac5.txt | sed 's/^/   /'
echo "-- 禁区三枚 blob（我的 rev-parse ⇄ 腿报 ⇄ 我导出树物理字节） --"
grep -E 'internal/config/writeguard.go HEAD|internal/config/loader.go HEAD|internal/config/migrate.go HEAD' .scratch/wisp/probes/341/v1/logs/08-gates-and-ac5.txt | sed 's/^/   /'
echo "-- 票框终量尺（七格一枚没翻） --"
T=".scratch/wisp/issues/341-four-hand-edit-shapes-still-unnailed-and-two-226-sentences-overclaim.md"
echo "   unchecked=$(grep -c '^- \[ \] \*\*AC#' "$T" || true)  checked=$(grep -c '^- \[x\]' "$T" || true)  lines=$(wc -l < "$T")"
echo "   title  = $(head -1 "$T" | cut -c1-26)..."
echo "   Status = $(sed -n 3p "$T" | cut -c1-26)..."
echo "   my bullet numstat vs parent:"; git show --numstat --format='' f5d5165c | sed 's/^/     /'
echo "   leg bullet still present = $(grep -c '落地腿 341-r1 交件' "$T" || true) (未抹)"
echo
echo "-- worktree/index hygiene --"
echo "   index count = $(git diff --cached --name-only | wc -l)"
echo "   my lanes status = [$(git status --porcelain -- internal/config cmd/wisp .scratch/wisp/probes/341/v1 docs/evidence/s1/341-v1-verdict.md "$T")]"
echo "   HEAD = $(git log -1 --format='%h %s' | cut -c1-50)"
echo "   internal/risk (343-r1 lane, untouched by me) = [$(git status --porcelain -- internal/risk)]"
