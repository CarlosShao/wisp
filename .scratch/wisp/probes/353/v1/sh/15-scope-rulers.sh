#!/bin/bash
# 353-v1 R15 = 射程尺（内容锚，对入库 blob；量的）：乙本体是否＝机主那一形、⛔ 扩射程
cd '/d/work/workspace/projects plans/Wisp' || exit 9
B=/d/tmp/wisp/20261020/353v1/blob
echo "###### 0) 每笔触碰 pinnedleaf 的 blob 串（逐笔）"
for c in a59f5630 3a509e50 fa395b33 110ed00b ffc3c62a HEAD; do
  printf '%-10s %s\n' "$c" "$(git rev-parse $c:internal/winsec/pinnedleaf_342_windows.go)"
done
echo "###### 1) :207-211 自陈段内容锚尺（BASE ⇄ 每笔）＝取段+sha256"
seg() { # <rev> -> 该段（从 'The error is still discarded' 到 'removal itself is now safe.'）
  git cat-file blob "$1:internal/winsec/pinnedleaf_342_windows.go" \
    | awk '/^\/\/ The error is still discarded/{f=1} f{print} /the removal itself is now safe\.$/{if(f)exit}'
}
for c in a59f5630 fa395b33 ffc3c62a HEAD; do
  h=$(seg "$c" | sha256sum | cut -d' ' -f1)
  n=$(seg "$c" | wc -l)
  printf '%-10s SELFREPORT_SEGMENT sha256=%s lines=%s\n' "$c" "$(echo $h|cut -c1-12)" "$n"
done
echo "-- 'now verifies identity first' 那串在每笔 blob 命中数（正控=HEAD 的乙段里确有 verifies identity 词） --"
for c in a59f5630 fa395b33 ffc3c62a HEAD; do
  printf '%-10s verifiesIdentityFirst=%s  CONTROL_identityWords=%s\n' "$c" \
   "$(git cat-file blob "$c:internal/winsec/pinnedleaf_342_windows.go" | grep -cF 'now verifies identity first' || true)" \
   "$(git cat-file blob "$c:internal/winsec/pinnedleaf_342_windows.go" | grep -c 'identity' || true)"
done
echo "###### 2) 删除原语逐字照旧（Root.Remove(base) 命中数 + 上下文）"
for c in a59f5630 HEAD; do
  printf '%-10s RemoveBase=%s\n' "$c" "$(git cat-file blob "$c:internal/winsec/pinnedleaf_342_windows.go" | grep -cF '_ = p.root.Remove(base)' || true)"
done
git cat-file blob "HEAD:internal/winsec/pinnedleaf_342_windows.go" | grep -nF 'p.root.Remove(base)' | sed 's/^/   /'
echo "###### 3) 导出名尺：BASE ⇄ HEAD 的 func 行集合作差（★只应见同名的参数/返回变化）"
for c in a59f5630 HEAD; do
  git cat-file blob "$c:internal/winsec/pinnedleaf_342_windows.go" | grep -E '^func ' | sed "s|^|$c\	|" > /d/tmp/wisp/20261020/353v1/func-$c.txt
done
echo "-- pinnedleaf func 名集作差（只名） --"
cut -f2 /d/tmp/wisp/20261020/353v1/func-a59f5630.txt | sed -E 's/^func \(([a-zA-Z0-9_ ]*\) )?([A-Za-z0-9_]+).*/\2/; s/^func ([A-Za-z0-9_]+).*/\1/' | sort -u > /d/tmp/wisp/20261020/353v1/fn-B.txt
cut -f2 /d/tmp/wisp/20261020/353v1/func-HEAD.txt       | sed -E 's/^func \(([a-zA-Z0-9_ ]*\) )?([A-Za-z0-9_]+).*/\2/; s/^func ([A-Za-z0-9_]+).*/\1/' | sort -u > /d/tmp/wisp/20261020/353v1/fn-L.txt
echo "  NEW=$(comm -13 /d/tmp/wisp/20261020/353v1/fn-B.txt /d/tmp/wisp/20261020/353v1/fn-L.txt | grep -c . || true)  GONE=$(comm -23 /d/tmp/wisp/20261020/353v1/fn-B.txt /d/tmp/wisp/20261020/353v1/fn-L.txt | grep -c . || true)"
comm -13 /d/tmp/wisp/20261020/353v1/fn-B.txt /d/tmp/wisp/20261020/353v1/fn-L.txt | sed 's/^/   NEWFN: /'
comm -23 /d/tmp/wisp/20261020/353v1/fn-B.txt /d/tmp/wisp/20261020/353v1/fn-L.txt | sed 's/^/   GONEFN: /'
echo "###### 4) 大写字母开头的导出名（包级 var/func/type）在 winsec.go+3枚pinnedleaf 的作差"
for c in a59f5630 HEAD; do
  git cat-file blob "$c:internal/winsec/winsec.go" > $B/x.txt
  cat $B/x.txt > $B/agg-$c.go.txt
  for f in pinnedleaf_342_windows.go pinnedleaf_342_unix.go pinnedleaf_342_fallback.go; do
    git cat-file blob "$c:internal/winsec/$f" >> $B/agg-$c.go.txt
  done
  grep -hoE '^(func|var|type|const) [A-Z][A-Za-z0-9_]*' $B/agg-$c.go.txt | sort -u > /d/tmp/wisp/20261020/353v1/exp-$c.txt
done
echo "  EXPORTED BASE=$(grep -c . /d/tmp/wisp/20261020/353v1/exp-a59f5630.txt) HEAD=$(grep -c . /d/tmp/wisp/20261020/353v1/exp-HEAD.txt)"
echo "  NEWEXPORT=$(comm -13 /d/tmp/wisp/20261020/353v1/exp-a59f5630.txt /d/tmp/wisp/20261020/353v1/exp-HEAD.txt | grep -c . || true)"
comm -13 /d/tmp/wisp/20261020/353v1/exp-a59f5630.txt /d/tmp/wisp/20261020/353v1/exp-HEAD.txt | sed 's/^/   /'
echo "###### 5) 旗标行 ⇄ fs.ErrExist 那一支：winsec.go 两处逐字（PrivateFile/PrivateFileExclusive）"
for c in a59f5630 HEAD; do
  git cat-file blob "$c:internal/winsec/winsec.go" > $B/w-$c.txt
done
echo "-- O_WRONLY|O_CREATE|O_TRUNC 行两侧 --"
grep -nF 'os.O_WRONLY|os.O_CREATE|os.O_TRUNC' $B/w-a59f5630.txt | sed 's/^/  B /'
grep -nF 'os.O_WRONLY|os.O_CREATE|os.O_TRUNC' $B/w-HEAD.txt | sed 's/^/  L /'
grep -nF 'os.O_EXCL' $B/w-a59f5630.txt | sed 's/^/  B /'
grep -nF 'os.O_EXCL' $B/w-HEAD.txt | sed 's/^/  L /'
echo "-- winsec.go BASE..HEAD 可执行变化行数（尺＝diff 里以+/-开头且非注释） --"
diff -u0 $B/w-a59f5630.txt $B/w-HEAD.txt | grep -E '^[+-][^+-]' | grep -vcE '^[+-][[:space:]]*//' | sed 's/^/  NONCOMMENT_DIFF_LINES=/'
diff -u0 $B/w-a59f5630.txt $B/w-HEAD.txt | grep -E '^[+-][^+-]' | wc -l | sed 's/^/  ALL_DIFF_LINES=/'
echo "###### 6) privateFile 返回形状：签名⇄return 语句两侧作差"
grep -nE 'func privateFile|return sealError|return err$|return nil$' $B/w-a59f5630.txt | sed 's/^/  B /' > /d/tmp/wisp/20261020/353v1/ret-B.txt
grep -nE 'func privateFile|return sealError|return err$|return nil$' $B/w-HEAD.txt   | sed 's/^/  L /' > /d/tmp/wisp/20261020/353v1/ret-L.txt
sed 's/^  [BL] //' /d/tmp/wisp/20261020/353v1/ret-B.txt > /d/tmp/wisp/20261020/353v1/ret-B2.txt
sed 's/^  [BL] //' /d/tmp/wisp/20261020/353v1/ret-L.txt > /d/tmp/wisp/20261020/353v1/ret-L2.txt
diff -u0 /d/tmp/wisp/20261020/353v1/ret-B2.txt /d/tmp/wisp/20261020/353v1/ret-L2.txt && echo "  RETURNS_IDENTICAL=YES"
echo "###### 7) :11-20 段两侧逐字（该段应完全未动）"
sed -n '11,20p' $B/BASE-pinnedleaf_342_windows.go.txt > /d/tmp/wisp/20261020/353v1/h1120-B.txt
sed -n '11,20p' $B/HEAD-pinnedleaf_342_windows.go.txt > /d/tmp/wisp/20261020/353v1/h1120-L.txt
diff -u0 /d/tmp/wisp/20261020/353v1/h1120-B.txt /d/tmp/wisp/20261020/353v1/h1120-L.txt && echo "  LINES_11_20_IDENTICAL=YES"
echo "###### 8) 新接缝尺：包级 var 集合两侧作差（applyDescriptor 等）"
for c in a59f5630 HEAD; do
  git grep -hnoE '^var [a-zA-Z].*=|^	[a-zA-Z][A-Za-z0-9_]* [a-zA-Z].*=' "$c" -- internal/winsec/*.go >/dev/null 2>&1
done
echo "-- 用 blob 直接取包级 var 行 --"
for c in a59f5630 HEAD; do
  { git cat-file blob "$c:internal/winsec/winsec.go"; git cat-file blob "$c:internal/winsec/pinnedleaf_342_windows.go"; } | grep -E '^var ' | sort -u > /d/tmp/wisp/20261020/353v1/var-$c.txt
  printf '  %s pkgvar=%s\n' "$c" "$(grep -c . /d/tmp/wisp/20261020/353v1/var-$c.txt||true)"
done
diff -u0 /d/tmp/wisp/20261020/353v1/var-a59f5630.txt /d/tmp/wisp/20261020/353v1/var-HEAD.txt && echo "  PKGVARS_IDENTICAL=YES"
