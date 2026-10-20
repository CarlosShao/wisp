#!/bin/bash
# 353-v1 R11 = 还原哈希/去牙三发 独立复量（尺＝sha256sum 打在入库 blob 导出字节 + grep 打在底片；量的）
R=/d/tmp/wisp/20261020/353r1
L=$R/logs
cd '/d/work/workspace/projects plans/Wisp' || exit 9
echo "###### 1) pristine 哈希独立复算（入库 blob，非工作树）"
git cat-file blob "HEAD:internal/winsec/pinnedleaf_342_windows.go"    | sha256sum | cut -d' ' -f1 | sed 's/^/HEAD-pinnedleaf   sha256=/'
git cat-file blob "a59f5630:internal/winsec/pinnedleaf_342_windows.go"| sha256sum | cut -d' ' -f1 | sed 's/^/BASE-pinnedleaf   sha256=/'
git cat-file blob "HEAD:internal/winsec/winsec.go"                    | sha256sum | cut -d' ' -f1 | sed 's/^/HEAD-winsec       sha256=/'
git cat-file blob "a59f5630:internal/winsec/winsec.go"                | sha256sum | cut -d' ' -f1 | sed 's/^/BASE-winsec       sha256=/'
git cat-file blob "HEAD:internal/winsec/nametopup_353_windows_test.go"| sha256sum | cut -d' ' -f1 | sed 's/^/HEAD-nametopup    sha256=/'
git cat-file blob "HEAD:internal/winsec/ancestorwindow_342r2_windows_test.go" | sha256sum | cut -d' ' -f1 | sed 's/^/HEAD-ancestorwin  sha256=/'
git cat-file blob "a59f5630:internal/winsec/ancestorwindow_342r2_windows_test.go" | sha256sum | cut -d' ' -f1 | sed 's/^/BASE-ancestorwin  sha256=/'
echo "###### 2) 终量两棵树在盘字节（腿的台面是否＝它声称的）"
for t in finalB finalL; do
  for f in pinnedleaf_342_windows.go winsec.go nametopup_353_windows_test.go ancestorwindow_342r2_windows_test.go; do
    printf '%-8s %-40s %s\n' "$t" "$f" "$(sha256sum "$R/$t/internal/winsec/$f" 2>/dev/null | cut -d' ' -f1 | cut -c1-12)"
  done
done
echo "###### 3) 去牙三发：四数 + 分子 + 红名 + 正控 + 还原"
for n in 1 2 3; do
  f="$L/final-mut-$n.raw.txt"
  echo "-- $f  ($(wc -c < "$f") B)"
  printf '   RUN=%s PASS=%s FAIL=%s SKIP=%s\n' "$(grep -c '^=== RUN' $f||true)" "$(grep -c '^--- PASS' $f||true)" "$(grep -c '^--- FAIL' $f||true)" "$(grep -c '^--- SKIP' $f||true)"
  grep -E 'over 3 shots' "$f" | sed 's/^/   /'
  grep -E '^--- FAIL' "$f" | sed 's/ *([0-9.]*s)$//' | sed 's/^/   REDNAME: /'
  grep -E 'AC#2 CONTROL' "$f" | sed -e 's/^\(.\{200\}\).*/\1/' | sed 's/^/   CTRL: /'
done
echo "###### 4) sh/24 里三发落点锚是否互不相同（内容锚尺）"
grep -nE 'anchor|ANCHOR|sed|perl|mut' "$R/../353r1/sh/24-final-mutations.sh" 2>/dev/null | head -5
grep -nE "^[[:space:]]*(M[0-9]|OLD|NEW|ANCH|pin|sed -|perl)" /d/tmp/wisp/20261020/353r1/sh/24-final-mutations.sh 2>/dev/null | sed -n '1,60p'
