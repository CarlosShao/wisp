#!/bin/bash
# 353-v1 R18 = 落点位置尺（内容锚⇄行号，对入库 blob）＋⛔借两族分子尺＋门禁底片在位尺（量的）
cd '/d/work/workspace/projects plans/Wisp' || exit 9
B=/d/tmp/wisp/20261020/353v1/blob
echo "###### 1) nametopup HEAD blob：两味落点逐枚行号（尺＝grep -n 内容锚）"
grep -n 'WHAT THIS INSTRUMENT CANNOT COVER' $B/HEAD-nametopup_353_windows_test.go.txt | sed 's/^/  /'
grep -n 'NOT refused and NOT measurable here' $B/HEAD-nametopup_353_windows_test.go.txt | sed 's/^/  /'
grep -n 'AC#2 GREEN SCOPE' $B/HEAD-nametopup_353_windows_test.go.txt | sed 's/^/  /'
grep -n 'did not ENTER' $B/HEAD-nametopup_353_windows_test.go.txt | sed 's/^/  /'
grep -n 'AC#2 CONTROL' $B/HEAD-nametopup_353_windows_test.go.txt | sed 's/^/  /'
echo "-- 该件里 truncated / deleted 两族分子被引用的枚数（应 0＝⛔ 借） --"
grep -c 'truncated' $B/HEAD-nametopup_353_windows_test.go.txt || true
grep -c 'deleted=' $B/HEAD-nametopup_353_windows_test.go.txt || true
grep -n 'truncated' $B/HEAD-nametopup_353_windows_test.go.txt | sed 's/^/  /'
echo "###### 2) pinnedleaf HEAD blob：两味落点行号"
grep -n 'does NOT stop one arriving after it' $B/HEAD-pinnedleaf_342_windows.go.txt | sed 's/^/  /'
grep -n 'NOT stopped. Closing that window needs a hook inside' $B/HEAD-pinnedleaf_342_windows.go.txt | sed 's/^/  /'
grep -n 'observable behaviour change in the attack state' $B/HEAD-pinnedleaf_342_windows.go.txt | sed 's/^/  /'
grep -n 'no new exported name, no new seam' $B/HEAD-pinnedleaf_342_windows.go.txt | sed 's/^/  /'
echo "###### 3) ancestorwindow HEAD blob :405-425（既有钉两味落点）"
sed -n '405,425p' $B/HEAD-ancestorwindow_342r2_windows_test.go.txt | nl -ba -v405 | cut -c1-200
echo "###### 4) 0b75e40b 之后有没有任何 winsec .go 字节漂（决定门禁底片是否还合现势）"
git log --format='%h %ad %s' --date=format:'%H:%M' 0b75e40b..HEAD -- internal/winsec | sed 's/^/  /' || true
echo "  WINSEC_COMMITS_AFTER_0b75e40b=$(git log --format='%h' 0b75e40b..HEAD -- internal/winsec | grep -c . || true)"
echo "  ANY_GO_AFTER_0b75e40b=$(git log --format='%h' 0b75e40b..HEAD -- '*.go' | grep -c . || true)"
git log --format='%h %ad %s' --date=format:'%H:%M' 0b75e40b..HEAD -- '*.go' | cut -c1-150 | sed 's/^/  /'
echo "###### 5) 门禁底片在位＋字节（⛔ 0 字节＝假读数；尺＝wc -c）"
for f in gates.txt gate-gofumpt.txt gate-d22scan.txt gate-length.txt gate-vet-native.txt gate-vet-linux.txt gate-vet-darwin.txt build-native.txt build-linux.txt build-darwin.txt build-js.txt; do
  printf '  %-22s %s B\n' "$f" "$(wc -c < /d/tmp/wisp/20261020/353r1/logs/$f 2>/dev/null || echo MISSING)"
done
echo "-- gates.txt 逐字 --"
cat /d/tmp/wisp/20261020/353r1/logs/gates.txt
echo "-- build-js.txt 逐字（腿报 rc=1 存量形状） --"
cat /d/tmp/wisp/20261020/353r1/logs/build-js.txt
echo "###### 6) d22scan 底片里与本腿/本票相关的行（⛔ 我跑 d22scan＝内含 go run；只读底片）"
grep -n 'violation\|clean\|FAIL\|pinnedleaf\|nametopup' /d/tmp/wisp/20261020/353r1/logs/gate-d22scan.txt | tail -20 | cut -c1-200 | sed 's/^/  /'
