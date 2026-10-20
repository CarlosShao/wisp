#!/bin/sh
# 353-r1 交件尺包（更正版）＝
#  ⚠ 上一版 sh/28 的 D 段尺⛔ 有效：`grep -vE '^(\.scratch/…/r1/|…)$'` 把目录前缀那一支也钉了行尾⇒永⛔ 匹配⇒
#    整表被判"越界"＝⛔ 拿坏尺当凭据（公共块：两把尺矛盾时先假定新那把坏）。本版改成"逐枚前缀==0/1 计数"。
#  ⚠ A 段的 CR 正控⛔ 是正控（它量的那枚件本就 CR=0）。本版正控＝先用同一条尺在**全仓 blob** 里找出一枚 CR>0 的件。
set -u
cd '/d/work/workspace/projects plans/Wisp' || exit 9
L=/d/tmp/wisp/20261020/353r1/logs
{
  echo "===== A. CR 尺⇄★真正控（同一条 tr 尺，先找一枚已知 CR>0 的件）====="
  git cat-file blob HEAD:internal/winsec/pinnedleaf_342_windows.go > "$L/b.txt"
  echo "  本腿那枚：CR=$(tr -cd '\r' < "$L/b.txt" | wc -c) BYTES=$(wc -c < "$L/b.txt")"
  FOUND=""
  for f in $(git ls-files | head -400); do
    n=$(git cat-file blob "HEAD:$f" 2>/dev/null | tr -cd '\r' | wc -c)
    if [ "$n" -gt 0 ] 2>/dev/null; then FOUND="$f"; echo "  正控件＝$f CR=$n（★同一条尺报>0⇒尺没瞎、0 是真 0）"; break; fi
  done
  [ -n "$FOUND" ] || echo "  ⚠ 前 400 枚里没找到 CR>0 的件＝正控缺位，此 0 只算弱读数（具名上报）"
  echo "  钝尺对照（本机⛔ 用它）＝对工作树 grep -c \$'\\r' 报 $(grep -c $'\r' internal/winsec/pinnedleaf_342_windows.go || true)＝真 CR 字节 0⇒同一条坑（台账 A1066③）"
  echo
  echo "===== D. 写面越界尺（更正版＝逐枚前缀计数，⛔ 行尾锚在目录形上）====="
  git log --all --grep='^票 353 落地腿 353-r1' --format='%h' > "$L/mycommits.txt"
  BAD=0; TOTAL=0
  while read -r f; do
    [ -n "$f" ] || continue
    TOTAL=$((TOTAL+1))
    case "$f" in
      .scratch/wisp/probes/353/r1/*) : ;;
      internal/winsec/nametopup_353_windows_test.go|internal/winsec/pinnedleaf_342_windows.go|\
      internal/winsec/pinnedleaf_342_unix.go|internal/winsec/pinnedleaf_342_fallback.go|\
      internal/winsec/winsec.go|internal/winsec/ancestorwindow_342r2_windows_test.go) : ;;
      *) BAD=$((BAD+1)); echo "  OUT_OF_RANGE: $f" ;;
    esac
  done < <(for c in $(cat "$L/mycommits.txt"); do git show --name-only --format='' "$c" | sed '/^$/d'; done | sort -u)
  echo "  FILES_TOTAL=$TOTAL OUT_OF_RANGE=$BAD（★期望 BAD=0）"
  echo "  .go 名册（去重）："
  for c in $(cat "$L/mycommits.txt"); do git show --name-only --format='' "$c" | sed '/^$/d'; done | grep '\.go$' | sort -u | sed 's/^/     /'
  echo
  echo "===== E2. 票面：⛔ 动过自证＋框终量 ======"
  T='.scratch/wisp/issues/353-pinned-name-swap-after-seal-delegates.md'
  TT='.scratch/wisp/issues/353-pinned-name-swap-after-seal-deletes-an-object-the-caller-never-named.md'
  echo "  本腿 4 笔里触碰票面的笔数＝$(for c in $(cat "$L/mycommits.txt"); do git show --name-only --format='' "$c" | grep -c 'issues/353' || true; done | paste -sd+ | tr -d ' ')"
  echo "  票面现量框＝未勾 $(git cat-file blob "HEAD:$TT" | grep -cE '^- \[ \]' || true)／已勾 $(git cat-file blob "HEAD:$TT" | grep -cE '^- \[x\]' || true)"
  echo "  票名⇄Status⇄-done：$(git log --all --format='%h' -- $TT | wc -l) 笔历史上动过该文件；本腿 0 笔（见上）"
} 2>&1
