#!/bin/sh
# 352-r1 AC#2 逐名作差尺（第二版＝★行区间切片；第一版按 awk 缓冲归属＝包归属算错＝本腿自抓、改在入库前）。
# 定尺＝go test -v 多包输出里 `^(ok|FAIL)<TAB><pkg><TAB>时间` 汇总行的行号＝该包区间右边界，
#        上一枚汇总行＋1＝左边界 ⇒ 逐枚包只在自己那段区间里数四数＋抽红名。
# ★⛔ 拿整包 RUN/PASS_TOP 的差当"我的字节改了 N 枚测试"＝只认逐名作差（comm -13）。
# ★四数只在 ^=== RUN／^--- PASS／^--- FAIL／^--- SKIP 取、⛔ 混用⛔ 相减。
set -u
ROOT=/d/tmp/wisp/20261020/352r1/logs
OUT="$ROOT/perpkg-diff.txt"
: > "$OUT"

NAMED='TestAC5FailedSealRefusesTheWrite
TestAC5FailureIsNotSwallowedByTheHappyPath
TestAC2OverwritePromiseHoldsForEveryProductionCaller322r1
TestAC3WindowsLeafGuardRefusesLinkNamesAndAcceptsRealOnes322r1
TestAC3WindowsPrivateFileRouteHasTheLeafGuardSpeaking322r1
TestAC1POSIXOpenThroughASwappedLeafLinkNeverTruncatesAForeignFile322r1
TestAC1POSIXSealFileThroughASymlinkRefusesAndLeavesTheForeignTreeAlone
TestAC1POSIXSealFileRefusesALinkAncestorAtEveryDepth
TestAC1POSIXSealDirThroughASymlinkRefuses
TestAC1POSIXPrivateFileThroughASymlinkRefusesAndWritesNothing
TestAC1POSIXSealFileThroughABackslashNamedLink
TestAC3POSIXSealFileStillNarrowsAPlainFileInsideTheNamedTree
TestAC3POSIXSealDoesNotFoldABackslashIntoASeparator'

# 汇总行表：包名 \t 起行 \t 止行
bounds() {
  # ★`ok` 行是 "ok<3空格><TAB>github…"、`FAIL` 行是 "FAIL<TAB>github…" ⇒ 分隔符必用 [[:space:]]+（枚一版用单个
  #   [ \t]＝只咬到 FAIL 那两枚⇒ok 那四枚 BOUNDS_MISSING＝本腿自抓的第二把钝尺，改在入库前）。
  grep -nE '^(ok|FAIL)[[:space:]]+github\.' "$1" | awk -F: -v prev=0 '
    { line=$1+0; sub(/^[0-9]+:/, "", $0); split($0, a, /[ \t]+/); pkg=a[2];
      printf "%s\t%d\t%d\n", pkg, prev+1, line; prev=line }
  '
}

getbound() { # $1=boundsfile $2=pkgfullname  -> "start end"
  awk -F'\t' -v q="$2" '$1==q {print $2" "$3; exit}' "$1"
}

for PAIR in "ac2-host-pair-s1-treeB:ac2-host-pair-s2-treeL" "ac2-col-s1-treeB:ac2-col-s2-treeL"; do
  B="${PAIR%%:*}"; L="${PAIR##*:}"
  echo "##### PAIR ${B}.raw.txt  <->  ${L}.raw.txt" >> "$OUT"
  bounds "$ROOT/${B}.raw.txt" > "$ROOT/boundsB.tsv"
  bounds "$ROOT/${L}.raw.txt" > "$ROOT/boundsL.tsv"
  for P in internal/secret internal/config internal/memory internal/agent internal/panel internal/winsec; do
    Q="github.com/CarlosShao/wisp/$P"
    S=$(printf '%s' "$P" | tr '/' '_')
    bb=$(getbound "$ROOT/boundsB.tsv" "$Q")
    ll=$(getbound "$ROOT/boundsL.tsv" "$Q")
    if [ -z "$bb" ] || [ -z "$ll" ]; then echo "   $P BOUNDS_MISSING b=[$bb] l=[$ll]" >> "$OUT"; continue; fi
    bs=${bb%% *}; be=${bb##* }; ls_=${ll%% *}; le=${ll##* }
    sed -n "${bs},${be}p" "$ROOT/${B}.raw.txt" > "$ROOT/sliceB.$S.txt"
    sed -n "${ls_},${le}p" "$ROOT/${L}.raw.txt" > "$ROOT/sliceL.$S.txt"
    brun=$(grep -c '^=== RUN' "$ROOT/sliceB.$S.txt" || true)
    bpt=$(grep -c '^--- PASS' "$ROOT/sliceB.$S.txt" || true)
    bft=$(grep -c '^--- FAIL' "$ROOT/sliceB.$S.txt" || true)
    bsk=$(grep -c '^--- SKIP' "$ROOT/sliceB.$S.txt" || true)
    lrun=$(grep -c '^=== RUN' "$ROOT/sliceL.$S.txt" || true)
    lpt=$(grep -c '^--- PASS' "$ROOT/sliceL.$S.txt" || true)
    lft=$(grep -c '^--- FAIL' "$ROOT/sliceL.$S.txt" || true)
    lsk=$(grep -c '^--- SKIP' "$ROOT/sliceL.$S.txt" || true)
    grep '^--- FAIL' "$ROOT/sliceB.$S.txt" | sed 's/^--- FAIL: //; s/ *([0-9.]*s)//' | sort > "$ROOT/rB.$S.txt"
    grep '^--- FAIL' "$ROOT/sliceL.$S.txt" | sed 's/^--- FAIL: //; s/ *([0-9.]*s)//' | sort > "$ROOT/rL.$S.txt"
    comm -13 "$ROOT/rB.$S.txt" "$ROOT/rL.$S.txt" > "$ROOT/newred.$S.txt"
    comm -23 "$ROOT/rB.$S.txt" "$ROOT/rL.$S.txt" > "$ROOT/gone.$S.txt"
    NC=$(grep -c . "$ROOT/newred.$S.txt" || true); GC=$(grep -c . "$ROOT/gone.$S.txt" || true)
    printf '%-18s B RUN=%-4s PASS_TOP=%-4s FAIL_TOP=%-3s SKIP=%-3s | L RUN=%-4s PASS_TOP=%-4s FAIL_TOP=%-3s SKIP=%-3s | NEW_RED=%s GONE=%s\n' \
      "$P" "$brun" "$bpt" "$bft" "$bsk" "$lrun" "$lpt" "$lft" "$lsk" "$NC" "$GC" >> "$OUT"
    [ "$NC" != 0 ] && { echo "     NEW_RED_NAMES: $(tr '\n' ' ' < "$ROOT/newred.$S.txt")" >> "$OUT"; }
    [ "$GC" != 0 ] && { echo "     GONE_NAMES: $(tr '\n' ' ' < "$ROOT/gone.$S.txt")" >> "$OUT"; }
    echo "     B_reds=[$(tr '\n' ' ' < "$ROOT/rB.$S.txt")]" >> "$OUT"
    echo "     L_reds=[$(tr '\n' ' ' < "$ROOT/rL.$S.txt")]" >> "$OUT"
  done
  echo "   == 票面具名连带逐枚（纯字符串尺＝整枚 raw，⛔ 依赖包切块）==" >> "$OUT"
  echo "$NAMED" | while IFS= read -r T; do
    [ -z "$T" ] && continue
    bs=$(grep -m1 -oE "^--- (PASS|FAIL|SKIP): $T" "$ROOT/${B}.raw.txt" | sed -E 's/^--- //; s/: .*//')
    lst=$(grep -m1 -oE "^--- (PASS|FAIL|SKIP): $T" "$ROOT/${L}.raw.txt" | sed -E 's/^--- //; s/: .*//')
    [ -z "$bs" ] && bs=ABSENT; [ -z "$lst" ] && lst=ABSENT
    [ "$bs" = "$lst" ] && SAME=SAME || SAME=DIFF
    printf '     %-70s B=%-7s L=%-7s %s\n' "$T" "$bs" "$lst" "$SAME" >> "$OUT"
  done
  echo "   == private_fail_test.go:37/50/61/72/75 五枚 ErrNotSealable 判定的宿主＝四枚子测试逐枚 ==" >> "$OUT"
  for SUB in exclusive_artifact replacement_write directory_chain seal_file_and_dir_direct; do
    bs=$(grep -m1 -oE "^ +--- (PASS|FAIL|SKIP): TestAC5FailedSealRefusesTheWrite/$SUB" "$ROOT/${B}.raw.txt" | sed -E 's/^ +--- //; s/: .*//')
    lst=$(grep -m1 -oE "^ +--- (PASS|FAIL|SKIP): TestAC5FailedSealRefusesTheWrite/$SUB" "$ROOT/${L}.raw.txt" | sed -E 's/^ +--- //; s/: .*//')
    [ -z "$bs" ] && bs=ABSENT; [ -z "$lst" ] && lst=ABSENT
    [ "$bs" = "$lst" ] && SAME=SAME || SAME=DIFF
    printf '     %-70s B=%-7s L=%-7s %s\n' "TestAC5FailedSealRefusesTheWrite/$SUB" "$bs" "$lst" "$SAME" >> "$OUT"
  done
  echo >> "$OUT"
done
wc -l "$OUT"
