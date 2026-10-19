#!/bin/sh
set -u
REPO='/d/work/workspace/projects plans/Wisp'
cd "$REPO" || exit 9
V=.scratch/wisp/probes/349/v1
T=/d/tmp/wisp/20261019/349v1
echo '########## 349-v1 三把尺自证（快形＝单次 %H%x09%s 导出＋内存 index==1，⛔ 逐枚 git log）##########'
git log --all --format='%H%x09%s' > "$T/logs/export2.txt"
awk -F'\t' 'index($2,"票 349 终裁腿 349-v1")==1 {print $1}' "$T/logs/export2.txt" > "$T/logs/mine-commits.txt"
MINE=$(cat "$T/logs/mine-commits.txt")
echo "MY_COMMITS=$(printf '%s\n' "$MINE" | wc -l)"
printf '%s\n' "$MINE" | while read -r h; do echo "  $(git log -1 --format='%h %ad %s' --date=format:'%H:%M' "$h")"; done
echo '### 尺① 逐枚 name-only ⊆ 三处写面；.go 命中 0'
printf '%s\n' "$MINE" | while read -r h; do
  echo "--- $(git log -1 --format='%h' "$h") name-only ---"
  git show --name-only --format='' "$h" | grep -v '^$' | while read -r p; do
    case "$p" in
      docs/evidence/s1/349-v1-verdict.md) echo "  INSCOPE $p";;
      .scratch/wisp/probes/349/v1/*) echo "  INSCOPE $p";;
      .scratch/wisp/issues/349-*) echo "  INSCOPE-bulletfile $p";;
      *) echo "  *** OUT_OF_SCOPE $p";;
    esac
  done
  printf '  go_hits=%s\n' "$(git show --name-only --format='' "$h" | grep -c '\.go$' || true)"
done
echo '### 正控（尺灵）＝同一把尺打 349-r1 两把锚并集 7 笔 ⇒ .go 恰 3'
awk -F'\t' 'index($2,"票 349 落地腿 349-r1")==1 || index($2,"349-r1")==1 {print $1}' "$T/logs/export2.txt" | sort -u > "$T/logs/union2.txt"
echo "  union=$(wc -l < "$T/logs/union2.txt")"
for h in $(cat "$T/logs/union2.txt"); do git show --name-only --format='' "$h"; done | grep '\.go$' | sort -u | sed 's/^/  GO /'
for h in $(cat "$T/logs/union2.txt"); do git show --name-only --format='' "$h"; done | grep '\.go$' | sort -u | wc -l | sed 's/^/  POSITIVE_CONTROL_go_unique=/;s/$/ (want 3)/'
echo '### 尺② 件：逐枚 bytes/lines/U+00A7，非 0 字节，⛔ .out'
F=docs/evidence/s1/349-v1-verdict.md
printf 'verdict bytes=%s lines=%s U00A7=%s U00B7=%s\n' "$(wc -c < "$F")" "$(wc -l < "$F")" \
  "$(perl -CSD -0ne 'print scalar(() = /\x{00A7}/g)' "$F")" "$(perl -CSD -0ne 'print scalar(() = /\x{00B7}/g)' "$F")"
find "$V" -type f | sort | while read -r f; do
  printf '  %-34s bytes=%-6s lines=%-4s U00A7=%s\n' "${f#"$V"/}" "$(wc -c < "$f")" "$(wc -l < "$f")" "$(perl -CSD -0ne 'print scalar(() = /\x{00A7}/g)' "$f")"
done
printf 'zero_byte=%s out=%s\n' "$(find "$V" -type f -size 0 | wc -l)" "$(find "$V" -name '*.out' | wc -l)"
echo '  占位符尺（本腿件 want 0）='$(grep -rniE 'TODO|FIXME|placeholder|待补|XXXX' "$V" "$F" | wc -l)
echo '  占位符尺正控（同尺打 tools/d22scan 与 HANDOVER，want 非 0）='$(grep -rniE 'TODO' tools/d22scan/main.go docs/reports/HANDOVER.md 2>/dev/null | wc -l)
printf '  HANDOVER_U00A7_POSITIVE_CONTROL=%s\n' "$(perl -CSD -0ne 'print scalar(() = /\x{00A7}/g)' docs/reports/HANDOVER.md)"
echo '### 尺③ 票框终量尺＋本腿对票面纯追加＋票名／Status 未动'
TICKET=.scratch/wisp/issues/349-blacklist-provenance-comments-name-undefined-symbols-fires-taskscope-cleanup.md
printf 'boxes_unchecked=%s boxes_checked=%s\n' "$(grep -c '^- \[ \] \*\*AC#' "$TICKET")" "$(grep -c '^- \[x\] \*\*AC#' "$TICKET" || true)"
printf '%s\n' "$MINE" | while read -r h; do
  d=$(git diff-tree --no-commit-id --numstat -r "$h" -- "$TICKET" | tr -s '\t' ' ')
  [ -n "$d" ] && printf '  ticket-numstat %s [%s]\n' "$(git log -1 --format='%h' "$h")" "$d"
done
printf '  headline_or_status_lines_changed_in_my_commits=%s\n' "$(for h in $MINE; do git diff-tree --no-commit-id -p -U0 "$h" -- "$TICKET"; done | grep -cE '^[-+](# 349|- Status)')"
echo '### 尺④ 三枚产码 .go：交件锚 ⇄ 现 HEAD blob 同一（本腿⛔ 动 .go 字节）'
for f in internal/risk/blacklist.go internal/risk/provenance.go internal/risk/syncdirs.go; do
  printf '  %-32s ed8db1ea=%s HEAD=%s %s\n' "$(basename "$f")" "$(git rev-parse ed8db1ea:$f)" "$(git rev-parse HEAD:$f)" \
  "$([ "$(git rev-parse ed8db1ea:$f)" = "$(git rev-parse HEAD:$f)" ] && echo UNCHANGED || echo CHANGED)"
done
echo '### 尺⑤ winsec 车道（我只读具名快照、⛔ 读在飞工作树、⛔ 跑该包）'
printf '  winsec.go ed8db1ea=%s HEAD=%s\n' "$(git rev-parse ed8db1ea:internal/winsec/winsec.go)" "$(git rev-parse HEAD:internal/winsec/winsec.go)"
echo RC=$?
