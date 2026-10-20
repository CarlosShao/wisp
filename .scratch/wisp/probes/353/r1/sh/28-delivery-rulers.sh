#!/bin/sh
# 353-r1 交件尺包＝★入库 blob 上量 CR（⛔ 对工作树 grep -c $'\r'＝本机钝尺，台账 A1066③）⇄
#   名册三把尺（锚定行首 ⇄ 非锚定子串 ⇄ 单次导出＋内存 index==1）⇄ 逐枚写面自证 ⇄ 票框终量 ⇄ 票尾 numstat
#   ⇄ 残余那味的最终 file:line（★随笔④动过字节⇒行号重取，⛔ 用旧行号）。
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 9
L=/d/tmp/wisp/20261020/353r1/logs
{
  echo "===== A. 入库 blob：CR／字节／行数（尺＝git cat-file blob | tr -cd '\r' | wc -c）====="
  for f in internal/winsec/nametopup_353_windows_test.go internal/winsec/pinnedleaf_342_windows.go \
           internal/winsec/pinnedleaf_342_unix.go internal/winsec/pinnedleaf_342_fallback.go \
           internal/winsec/winsec.go internal/winsec/ancestorwindow_342r2_windows_test.go; do
    git cat-file blob "HEAD:$f" > "$L/b.txt"
    echo "  $f CR=$(tr -cd '\r' < "$L/b.txt" | wc -c) BYTES=$(wc -c < "$L/b.txt") LINES=$(wc -l < "$L/b.txt") BLOB=$(git rev-parse "HEAD:$f" | cut -c1-12)"
  done
  echo "  ★正控（同一条尺对一枚已知含 CR 的件）＝$(git cat-file blob HEAD:internal/winsec/pinnedleaf_342_windows.go | tr -cd '\r' | wc -c)（期望 0＝本仓 .go 全 LF；钝尺对照另见下）"
  echo "  钝尺对照（⛔ 用它判读＝本机报行计数）：对工作树 grep -c \$'\\r' 那形的命中＝$(grep -c $'\r' internal/winsec/pinnedleaf_342_windows.go || true)"
  echo
  echo "===== B. 名册三把尺（★锚定行首、⛔ 假设分隔符、⛔ 非锚定子串）====="
  echo "  尺1 锚定 ^票 353 落地腿 353-r1 ：$(git log --all --grep='^票 353 落地腿 353-r1' --format='%h' | wc -l)"
  echo "  尺2 非锚定 票 353 落地腿 353-r1 ：$(git log --all --grep='票 353 落地腿 353-r1' --format='%h' | wc -l)"
  git log --all --format='%s' > "$L/subjects.txt"
  echo "  尺3 单次导出＋内存 index==1 ：$(awk 'BEGIN{n=0} index($0,"票 353 落地腿 353-r1")==1{n++} END{print n}' "$L/subjects.txt")"
  echo "  逐枚主题（只 %s，⛔ 正文）:"
  git log --all --grep='^票 353 落地腿 353-r1' --format='    %h %ad %s' --date=format:'%H:%M' | sed 's/\(.\{130\}\).*/\1 …/'
  echo
  echo "===== C. 逐枚写面自证（name-only）＋.go 枚数 ======"
  for c in $(git log --all --grep='^票 353 落地腿 353-r1' --format='%h'); do
    echo "  commit $c"
    git show --name-only --format='' "$c" | sed '/^$/d' | sed 's/^/     /'
    echo "     GFILE_COUNT=$(git show --name-only --format='' "$c" | grep -c '\.go$' || true)"
  done
  echo
  echo "===== D. 越界自证（本腿名册里⛔ 允许的写面）====="
  git log --all --grep='^票 353 落地腿 353-r1' --format='%h' | while read -r c; do
    git show --name-only --format='' "$c" | sed '/^$/d'
  done | sort -u | grep -vE '^(\.scratch/wisp/probes/353/r1/|internal/winsec/(winsec\.go|pinnedleaf_342_windows\.go|pinnedleaf_342_unix\.go|pinnedleaf_342_fallback\.go|nametopup_353_windows_test\.go|ancestorwindow_342r2_windows_test\.go))$' | sed 's/^/  OUT_OF_RANGE: /' || true
  echo "  （上面⛔ 任何一行＝越界）"
  echo
  echo "===== E. 票 353 框终量＋票尾 numstat＋票名／Status ======"
  T='.scratch/wisp/issues/353-pinned-name-swap-after-seal-deletes-an-object-the-caller-never-named.md'
  echo "  UNCHECKED=$(git cat-file blob "HEAD:$T" | grep -cE '^- \[ \]' || true) CHECKED=$(git cat-file blob "HEAD:$T" | grep -cE '^- \[x\]' || true)"
  echo "  本腿对票面的改动枚数＝$(git log --all --grep='^票 353 落地腿 353-r1' --format='x' | wc -l) 笔里含票面的笔数＝$(for c in $(git log --all --grep='^票 353 落地腿 353-r1' --format='%h'); do git show --name-only --format='' "$c" | grep -c 'issues/353' || true; done | paste -sd+ | tr -d ' ')"
  echo "  numstat（本腿那些笔对票面）＝$(git log --all --grep='^票 353 落地腿 353-r1' --format='%h' | while read -r c; do git show --numstat --format='' "$c" | grep 'issues/353' || true; done | wc -l) 行（期望 0＝⛔ 动票面）"
  echo
  echo "===== F. ★"窗内那一支测不到"写进判据的最终 file:line（随笔④后重取）====="
  grep -n 'CANNOT COVER\|NOT refused and NOT measurable\|GREEN SCOPE\|did not ENTER\|positive control' internal/winsec/nametopup_353_windows_test.go | sed 's/^/  /'
  echo "  ---- 另一处＝产码 ⇄ 既有钉 ----"
  grep -n 'does NOT stop\|NOT covered and NOT claimed\|no instrument here can reach\|may not be read as covering' internal/winsec/pinnedleaf_342_windows.go internal/winsec/ancestorwindow_342r2_windows_test.go | sed 's/^/  /'
} 2>&1
