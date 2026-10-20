#!/bin/bash
# 353-v1 R22 = ancestorwindow 笔③ diff 全量 + 该枚名在各批底片里的色 + container 尾行
cd '/d/work/workspace/projects plans/Wisp' || exit 9
B=/d/tmp/wisp/20261020/353v1/blob
L=/d/tmp/wisp/20261020/353r1/logs
echo "###### 1) ancestorwindow BASE..HEAD 全 diff（看可执行行改了哪些）"
git cat-file blob "a59f5630:internal/winsec/ancestorwindow_342r2_windows_test.go" > $B/BASE-ancestorwindow_342r2_windows_test.go.txt
diff -u0 $B/BASE-ancestorwindow_342r2_windows_test.go.txt $B/HEAD-ancestorwindow_342r2_windows_test.go.txt | grep -E '^[+-][^+-]' | sed -e 's/^\(.\{150\}\).*/\1/' | sed 's/^/  /'
echo "  EXEC_CHANGED_LINES=$(diff -u0 $B/BASE-ancestorwindow_342r2_windows_test.go.txt $B/HEAD-ancestorwindow_342r2_windows_test.go.txt | grep -E '^[+-][^+-]' | grep -vcE '^[+-][[:space:]]*//')"
echo "  TESTFN_SET_DIFF:"
for c in BASE HEAD; do
  grep -E '^func Test' $B/$c-ancestorwindow_342r2_windows_test.go.txt | sed 's/(t \*testing.T) {//' > $V-dummy 2>/dev/null || true
done
grep -E '^func Test' $B/BASE-ancestorwindow_342r2_windows_test.go.txt > /d/tmp/wisp/20261020/353v1/anwin-fn-B.txt
grep -E '^func Test' $B/HEAD-ancestorwindow_342r2_windows_test.go.txt > /d/tmp/wisp/20261020/353v1/anwin-fn-L.txt
diff -u0 /d/tmp/wisp/20261020/353v1/anwin-fn-B.txt /d/tmp/wisp/20261020/353v1/anwin-fn-L.txt && echo "   TESTFN_IDENTICAL=YES"
echo "###### 2) 那枚摆动名在各批底片里的逐发色（量的；尺＝grep -cF 在 raw）"
NAME='TestAC2WindowsPrivateFileThroughASwappedAncestorNeverDeletesAForeignFile342r1'
for f in "$L"/final-full-B-1.raw.txt "$L"/final-full-B-2.raw.txt "$L"/mx-treeB2-1.raw.txt "$L"/mx-treeB2-2.raw.txt "$L"/fullpair-s1-treeB.raw.txt "$L"/pername-s1-treeB2.raw.txt; do
  printf '  B-side %-32s %s\n' "$(basename $f)" "$(grep -E "^--- (FAIL|PASS|SKIP): $NAME" "$f" | sed 's/ *([0-9.]*s)$//' | tr '\n' '|')"
done
for f in "$L"/final-full-L-1.raw.txt "$L"/final-full-L-2.raw.txt "$L"/mx-treeL-1.raw.txt "$L"/mx-treeL-2.raw.txt "$L"/fullpair-s2-treeL.raw.txt "$L"/pername-s2-treeL.raw.txt; do
  printf '  L-side %-32s %s\n' "$(basename $f)" "$(grep -E "^--- (FAIL|PASS|SKIP): $NAME" "$f" | sed 's/ *([0-9.]*s)$//' | tr '\n' '|')"
done
echo "###### 2b) 另两枚 342r1/342r2 摆动名同样逐批"
for NAME in 'TestAC1WindowsLeafRecheckAnswersOnTheBareSiblingForm342r2' 'TestAC2WindowsLeafIdentityReCheckRefusesAForeignNameAndNeverDeletesIt342r1' 'TestAC1WindowsPrivateFileThroughASwappedAncestorNeverTruncatesAForeignFile342r1'; do
  echo "  -- $NAME"
  for f in "$L"/final-full-B-1.raw.txt "$L"/final-full-B-2.raw.txt "$L"/final-full-L-1.raw.txt "$L"/final-full-L-2.raw.txt "$L"/mx-treeB2-1.raw.txt "$L"/mx-treeL-1.raw.txt; do
    printf '     %-28s %s\n' "$(basename $f)" "$(grep -E "^--- (FAIL|PASS|SKIP): $NAME" "$f" | sed 's/ *([0-9.]*s)$//' | tr '\n' '|')"
  done
done
echo "###### 3) container 底片尾行（ok/FAIL 判 RC 是否自洽；尺＝tail）"
for f in final-col-s1-finalB.raw.txt final-col-s2-finalL.raw.txt; do
  echo "  -- $f"; tail -3 "$L/$f" | sed 's/^/     /'
  grep -cE '^FAIL' "$L/$f" | sed 's/^/     ^FAIL 行=/'; grep -c '^ok' "$L/$f" | sed 's/^/     ^ok 行=/'
done
echo "###### 4) 三枚去牙锚在 HEAD blob 命中数（内容锚尺；期望各 1）"
for a in '_ = p.root.Remove(base)' 'f, err := p.root.OpenFile(base, flags, perm)' 'if err != nil || !os.SameFile(pinned, now) {'; do
  printf '  anchor[%-50s]=%s\n' "$a" "$(grep -cF "$a" $B/HEAD-pinnedleaf_342_windows.go.txt || true)"
done
echo "  正控：不存在的锚命中=$(grep -cF 'ZZZ-no-such-anchor' $B/HEAD-pinnedleaf_342_windows.go.txt || true)"
echo "###### 5) 底片里有无 ANCHOR_BEFORE / MUT_ANCHOR 存档（终量批）"
grep -rl 'ANCHOR_BEFORE' "$L" | sed 's/^/   /' || echo "   (无)"
ls -1 "$L" | grep -iE 'mut' | sed 's/^/   /'
