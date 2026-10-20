#!/bin/sh
# gw-1 交件自证三把尺（公共块"名册尺"＋"写面自证"＋"tracked 脏逐枚列名＋mtime"）。
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 9
R=/d/tmp/wisp/20261020/gw1
LOGS="$R/logs"
PFX='gw-1 量色腿 gw-1'
echo "== HEAD =="; git rev-parse HEAD | sed 's|^|HEAD_FULL=|'
echo
echo "== 尺1＝锚定行首打在 --format='%s' 的导出（射程＝只有主题）=="
git log --all --format='%s' > "$LOGS/roster-subject.txt"
grep -c "^$PFX" "$LOGS/roster-subject.txt" | sed 's|^|R1_ANCHORED_SUBJECT=|' || true
echo
echo "== 尺2＝git log --grep 锚定 ^（射程＝整条 message 的首行）=="
git log --all --grep="^$PFX" --format='%h %ad %s' --date=format:'%m-%d %H:%M' > "$LOGS/roster-grep.txt"
grep -c . "$LOGS/roster-grep.txt" | sed 's|^|R2_GREP_ANCHORED=|' || true
grep . "$LOGS/roster-grep.txt" | sed 's|^|  R2 |'
echo
echo "== 尺3＝单次 '%H%x09%s' 导出＋内存 index(\$2,PFX)==1 =="
git log --all --format='%H%x09%s' > "$LOGS/roster-pair.txt"
awk -F'\t' -v p="$PFX" 'index($2,p)==1 {print substr($1,1,10)}' "$LOGS/roster-pair.txt" > "$LOGS/roster-mem.txt"
grep -c . "$LOGS/roster-mem.txt" | sed 's|^|R3_MEM_INDEX=|' || true
grep . "$LOGS/roster-mem.txt" | sed 's|^|  R3 |'
echo
echo "== 非锚定子串对照（★⛔ 当名册＝公共块那形：会吃进正文引用）=="
grep -c "$PFX" "$LOGS/roster-subject.txt" | sed 's|^|R4_UNANCHORED_SUBJECT_LINE=|' || true
git log --all --grep="$PFX" --format='%h' | wc -l | sed 's|^|R5_UNANCHORED_FULLMSG=|'
echo
echo "== ★正控（证明锚定尺⛔ 瞎）＝同一条尺打一枚在册的别腿前缀 =="
for Q in '^票 139 落地腿 139-r1' '^票 321 终裁腿 321-v1' '^票 353 落地腿 353-r1'; do
  printf '  POSCTRL %s = %s\n' "$Q" "$(grep -c "$Q" "$LOGS/roster-subject.txt" || true)"
done
echo
echo "== 写面自证＝本腿名册逐枚 git show --name-only ＋ .go 枚数（★正控＝同尺打一枚含 .go 的别腿笔）=="
: > "$LOGS/writeface.txt"
while IFS= read -r h; do
  [ -n "$h" ] || continue
  git show --name-only --format='' "$h" | grep -v '^$' > "$LOGS/wf.body"
  TOTAL=$(grep -c . "$LOGS/wf.body" || true)
  GO=$(grep -c '\.go$' "$LOGS/wf.body" || true)
  OUTSIDE=$(grep -v '^\.scratch/wisp/probes/gw/1/' "$LOGS/wf.body" | grep -c . || true)
  printf '%s total=%s go=%s out_of_face=%s\n' "$h" "$TOTAL" "$GO" "$OUTSIDE" >> "$LOGS/writeface.txt"
  grep -v '^\.scratch/wisp/probes/gw/1/' "$LOGS/wf.body" | sed "s|^|  $h OOR |"
done < "$LOGS/roster-mem.txt"
grep . "$LOGS/writeface.txt" | sed 's|^|  WF |'
echo "  （正控＝同尺打在 353-r1 笔① 的 .go 枚数）"
git show --name-only --format='' 3a509e50 | grep -v '^$' > "$LOGS/wf.pos"
printf '  POSCTRL_353pen1 total=%s go=%s\n' "$(grep -c . "$LOGS/wf.pos")" "$(grep -c '\.go$' "$LOGS/wf.pos" || true)"
echo
echo "== tracked 脏＝逐枚列名＋mtime（★⛔ 只报计数）=="
git status --porcelain -uno > "$LOGS/dirty.raw"
echo "DIRTY_TRACKED_LINES=$(grep -c . "$LOGS/dirty.raw" || true)"
while IFS= read -r line; do
  [ -n "$line" ] || continue
  p=$(printf '%s' "$line" | cut -c4-)
  printf '  DIRTY [%s] %s mtime=%s\n' "$(printf '%s' "$line" | cut -c1-2)" "$p" "$(stat -c '%y' "$p" 2>/dev/null | cut -d. -f1)"
done < "$LOGS/dirty.raw"
echo "  .go 脏枚数＝$(git status --porcelain -uno -- '*.go' | wc -l)"
echo
echo "== 仓根冻结件监视（A1069③ 那形＝⛔ 0 字节）=="
for f in go.mod go.sum .gitignore .github/workflows/ci.yml; do
  printf '  FROZEN %s bytes=%s mtime=%s blob_cr=%s\n' "$f" "$(wc -c < "$f")" "$(stat -c '%y' "$f" | cut -d. -f1)" "$(git cat-file blob "HEAD:$f" | tr -cd '\r' | wc -c)"
done
echo "  （正控＝tr 尺对自造 2 枚 CR 件读出 $(tr -cd '\r' < "$R/posctrl/cr2.txt" | wc -c)）"
