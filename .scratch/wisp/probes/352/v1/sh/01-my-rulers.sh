#!/bin/sh
# 352-v1 终裁腿·独立尺件（可逐字重跑；凡引用数＝本件输出，⛔ 抄实现腿的数）
# 用法：sh .scratch/wisp/probes/352/v1/sh/01-my-rulers.sh
# ★CR／码位一律对入库 blob 量（git cat-file blob … | tr -cd '\r'），⛔ 对工作树用 grep -c $'\r'
cd '/d/work/workspace/projects plans/Wisp' || exit 9
P=.scratch/wisp/probes/352/r1
T=.scratch/wisp/issues/352-wrappath-errno-half-stringified-so-fs-ernotexist-cannot-punch-through-private-route.md

echo "== R0 起手锚（量的；★⛔ cut -c 截中文＝按字节切＝非法 UTF-8＝本腿第一版在此踩中并当场改掉）=="
git log -1 --format='HEAD=%h %ad' --date=format:'%Y-%m-%d %H:%M'
git log -1 --format='%s' | head -1

echo "== R1 AC#1：HEAD 的 wrapPath 函数体逐字（量的）=="
git cat-file blob HEAD:internal/winsec/winsec.go | sed -n '430,432p'
echo "  BLOB_AT_HEAD_NOW=$(git rev-parse HEAD:internal/winsec/winsec.go)  HEAD=$(git log -1 --format=%h)"
echo "  ★本趟期间 HEAD 被别腿推进＝blob 逐锚报清（量的）："
for c in bf509bb8 ecffb93b a651c086 HEAD; do
  printf '    %-9s winsec.go blob=%s\n' "$c" "$(git rev-parse $c:internal/winsec/winsec.go)"
done
printf '  那一枚格式串自 ecffb93b 起被改过否（尺＝git diff ecffb93b..HEAD 里含 wrapPath 行的增删枚数）= '
git diff ecffb93b HEAD -- internal/winsec/winsec.go | grep -E '^[+-][^+-]' | grep -c 'wrapPath\|%w: %s: %' || true
echo '  ⇐ 0＝353-r1 那 6 行漂移⛔ 碰 wrapPath（漂移本身＝另票〔乙〕＝⛔ 本票射程，具名登记）'

echo "== R2 AC#1：ecffb93b 净变动（量的）=="
git show -U0 --format='' ecffb93b -- internal/winsec/winsec.go | grep -E '^[+-][^+-]'
git show --numstat --format='' ecffb93b | sed 's/^/  numstat=/'
echo -n "  签名行改前⇄改后同否="
a=$(git cat-file blob ecffb93b~1:internal/winsec/winsec.go | sed -n '430p')
b=$(git cat-file blob ecffb93b:internal/winsec/winsec.go | sed -n '430p')
[ "$a" = "$b" ] && echo "SAME ($a)" || echo "DIFF"
echo -n "  sealError 体 435-446 改前后 symdiff 行数="
git cat-file blob ecffb93b~1:internal/winsec/winsec.go | sed -n '435,446p' > /d/tmp/wisp/20261020/352v1/scratch-seal-a.txt
git cat-file blob ecffb93b:internal/winsec/winsec.go | sed -n '435,446p' > /d/tmp/wisp/20261020/352v1/scratch-seal-b.txt
diff /d/tmp/wisp/20261020/352v1/scratch-seal-a.txt /d/tmp/wisp/20261020/352v1/scratch-seal-b.txt | grep -c '^[<>]' || true

echo "== R3 AC#1：两分类并存（量的，尺＝对入库底片 grep 逐名；★tag 里带"."＝第一版 [a-zA-Z -]+ 漏咬 R3 那一枚，本腿自抓并改成 [^:]+，改前读数照实在 logs 里具名）=="
for f in ac1-col-s2-mut-2-treeL ac1-host-s2-mut-2-treeL; do
  echo "  -- $f --"
  git cat-file blob HEAD:$P/logs/$f.raw.txt | grep -oE 'R[0-9]+ [^:]+: isErrNotExist=true isNotSealable=true' | sed 's/^/    /'
done
echo "  Windows 落地侧仍 false/true 的跳（量的＝⛔ 射程那两枚要照实点名）＝"
git cat-file blob HEAD:$P/logs/ac1-host-s2-mut-2-treeL.raw.txt | grep -oE 'R[0-9]+ [^:]+: isErrNotExist=false isNotSealable=true' | sed 's/^/    /'
echo "  同一把尺数 POSIX 落地侧 false/true（量的）＝"
git cat-file blob HEAD:$P/logs/ac1-col-s2-mut-2-treeL.raw.txt | grep -oE 'R[0-9]+ [^:]+: isErrNotExist=false isNotSealable=true' | sed 's/^/    /'

echo "== R4 名册两把独立尺（量的）=="
echo "  锚定行首 --grep='^票 352 落地腿 352-r1' = $(git log --grep='^票 352 落地腿 352-r1' --format='%h' | grep -c .)"
git log --format='%H%x09%s' > /d/tmp/wisp/20261020/352v1/scratch-subj.txt
awk -F'\t' 'index($2,"票 352 落地腿 352-r1")==1{c++} END{print "  内存内 index==1 = " c+0}' /d/tmp/wisp/20261020/352v1/scratch-subj.txt
echo "  非锚定形 = $(git log --grep='票 352 落地腿 352-r1' --format='%h' | grep -c .)"
echo -n "  正控（票 321 同一对尺，两把都须 >0）= "
git log --grep='^票 321 落地腿 321-r1' --format='%h' | grep -c . | tr '\n' ' '
awk -F'\t' 'index($2,"票 321 落地腿 321-r1")==1{c++} END{print c+0}' /d/tmp/wisp/20261020/352v1/scratch-subj.txt

echo "== R5 写面⇄.go 枚数（算的，尺＝逐枚 git show --name-only）=="
for c in 0fb399c3 ecffb93b aacde632 c16ff47e 06aa2abf 3ecdea30 5c9402f8 1c787282 8a10a7e3 e74762b5 a651c086; do
  git show --name-only --format='' $c
done | grep '\.go$' | sort | uniq -c | sed 's/^/  /'
echo -n "  越界枚数（⛔ 三处写面之外）= "
for c in 0fb399c3 ecffb93b aacde632 c16ff47e 06aa2abf 3ecdea30 5c9402f8 1c787282 8a10a7e3 e74762b5 a651c086; do
  git show --name-only --format='' $c
done | grep -v -E '^(internal/winsec/winsec\.go|\.scratch/wisp/probes/352/r1/|\.scratch/wisp/issues/352-)' | grep -c . || true
echo -n "  禁区 pinnedleaf_342_windows.go 命中笔数 = "
for c in 0fb399c3 ecffb93b aacde632 c16ff47e 06aa2abf 3ecdea30 5c9402f8 1c787282 8a10a7e3 e74762b5 a651c086; do
  git show --name-only --format='' $c
done | grep -c 'pinnedleaf_342_windows.go' || true

echo "== R6 票面那味留账：wrapPath 调用点（量的，尺＝git grep 现算）=="
git grep -n 'wrapPath(' HEAD -- 'internal/winsec/*.go' | sed 's/^HEAD://' > /d/tmp/wisp/20261020/352v1/scratch-cs.txt
echo "  全仓命中（含定义）= $(grep -c . /d/tmp/wisp/20261020/352v1/scratch-cs.txt)；定义 = $(grep -c 'func wrapPath' /d/tmp/wisp/20261020/352v1/scratch-cs.txt)；调用点 = $(( $(grep -c . /d/tmp/wisp/20261020/352v1/scratch-cs.txt) - $(grep -c 'func wrapPath' /d/tmp/wisp/20261020/352v1/scratch-cs.txt) ))"
for f in winsec.go winsec_other.go winsec_windows.go pinnedleaf_342_windows.go pinnedleaf_342_fallback.go; do
  printf '  %-32s 行号=' "$f"
  git cat-file blob HEAD:internal/winsec/$f | grep -n 'wrapPath(' | cut -d: -f1 | tr '\n' ',' 
  echo
done
echo -n "  起手锚 bf509bb8 上同尺逐枚枚数="
for f in winsec.go winsec_other.go winsec_windows.go pinnedleaf_342_windows.go pinnedleaf_342_fallback.go; do
  printf '%s:%s ' "$f" "$(git cat-file blob bf509bb8:internal/winsec/$f | grep -c 'wrapPath(' || true)"
done
echo
echo -n "  _test.go 命中枚数 = "; git grep -c 'wrapPath(' HEAD -- 'internal/winsec/*_test.go' | grep -c . || true

echo "== R7 票面尺（量的）=="
echo "  行数=$(git cat-file blob HEAD:$T | wc -l)  Status 所在行=$(git cat-file blob HEAD:$T | grep -n '^\*\*Status:' | cut -d: -f1)"
echo "  未勾=$(git cat-file blob HEAD:$T | grep -cE '^- \[ \]')  已勾=$(git cat-file blob HEAD:$T | grep -cE '^- \[x\]' || true)  总框=$(git cat-file blob HEAD:$T | grep -cE '^- \[[ x]\]')"
echo -n "  正控票 348 同三把尺 = "
git cat-file blob HEAD:.scratch/wisp/issues/348-refusal-path-remove-follows-an-ancestor-link-and-deletes-a-foreign-file.md | grep -cE '^- \[x\] ' | tr '\n' '/'
git cat-file blob HEAD:.scratch/wisp/issues/348-refusal-path-remove-follows-an-ancestor-link-and-deletes-a-foreign-file.md | grep -cE '^- \[ \] ' | tr '\n' '/'
git cat-file blob HEAD:.scratch/wisp/issues/348-refusal-path-remove-follows-an-ancestor-link-and-deletes-a-foreign-file.md | grep -cE '^- \[[ x]\] '

echo "== R8 票面⛔ 改写历史（量的）=="
echo -n "  bf509bb8⇄HEAD 删除行数 = "; git diff bf509bb8 HEAD -- $T | grep -c '^-[^-]' || true
git diff --numstat bf509bb8 HEAD -- $T | sed 's/^/  numstat=/'
echo -n "  票尾 bullet 枚数 = "; git cat-file blob HEAD:$T | grep -c '^- \[2026-10-20 16:5x +08\] agent=`352-r1`'
echo -n "  前 24 行逐字节同否 = "
git cat-file blob bf509bb8:$T | head -24 > /d/tmp/wisp/20261020/352v1/scratch-tk-a.txt
git cat-file blob HEAD:$T | head -24 > /d/tmp/wisp/20261020/352v1/scratch-tk-b.txt
diff -q /d/tmp/wisp/20261020/352v1/scratch-tk-a.txt /d/tmp/wisp/20261020/352v1/scratch-tk-b.txt >/dev/null && echo IDENTICAL || echo DIFF

echo "== R9 AC#2 消费者面独立复算（量的，尺＝对入库底片整枚文件级四数＋逐名 comm；★本腿⛔ 跑 go test）=="
for f in ac2-host-pair-s1-treeB ac2-host-pair-s2-treeL ac2-col-s1-treeB ac2-col-s2-treeL; do
  git cat-file blob HEAD:$P/logs/$f.raw.txt > /d/tmp/wisp/20261020/352v1/scratch-$f.txt
  printf '  %-26s RUN=%s PASS_TOP=%s FAIL_TOP=%s SKIP=%s\n' "$f" \
    "$(grep -c '^=== RUN' /d/tmp/wisp/20261020/352v1/scratch-$f.txt || true)" \
    "$(grep -c '^--- PASS' /d/tmp/wisp/20261020/352v1/scratch-$f.txt || true)" \
    "$(grep -c '^--- FAIL' /d/tmp/wisp/20261020/352v1/scratch-$f.txt || true)" \
    "$(grep -c '^--- SKIP' /d/tmp/wisp/20261020/352v1/scratch-$f.txt || true)"
done
grep '^--- FAIL' /d/tmp/wisp/20261020/352v1/scratch-ac2-host-pair-s1-treeB.txt | sed 's/^--- FAIL: //;s/ .*//' | sort > /d/tmp/wisp/20261020/352v1/scratch-hb.txt
grep '^--- FAIL' /d/tmp/wisp/20261020/352v1/scratch-ac2-host-pair-s2-treeL.txt | sed 's/^--- FAIL: //;s/ .*//' | sort > /d/tmp/wisp/20261020/352v1/scratch-hl.txt
echo "  本机档逐名作差：新增红（只在落地侧）＝$(comm -13 /d/tmp/wisp/20261020/352v1/scratch-hb.txt /d/tmp/wisp/20261020/352v1/scratch-hl.txt | grep -c . || true) 枚；消失红＝$(comm -23 /d/tmp/wisp/20261020/352v1/scratch-hb.txt /d/tmp/wisp/20261020/352v1/scratch-hl.txt | grep -c . || true) 枚"
grep '^--- FAIL' /d/tmp/wisp/20261020/352v1/scratch-ac2-col-s1-treeB.txt | sed 's/^--- FAIL: //;s/ .*//' | sort > /d/tmp/wisp/20261020/352v1/scratch-cb.txt
grep '^--- FAIL' /d/tmp/wisp/20261020/352v1/scratch-ac2-col-s2-treeL.txt | sed 's/^--- FAIL: //;s/ .*//' | sort > /d/tmp/wisp/20261020/352v1/scratch-cl.txt
echo "  容器档逐名作差：新增红＝$(comm -13 /d/tmp/wisp/20261020/352v1/scratch-cb.txt /d/tmp/wisp/20261020/352v1/scratch-cl.txt | grep -c . || true) 枚；消失红＝$(comm -23 /d/tmp/wisp/20261020/352v1/scratch-cb.txt /d/tmp/wisp/20261020/352v1/scratch-cl.txt | grep -c . || true) 枚"
echo "  ★那枚换位名不在三发复量并集里（量的）："
for i in 1 2 3; do git cat-file blob HEAD:$P/logs/ac2r$i-winsec-s2-treeL.rednames.txt; done | sort -u | sed 's/^/    L_union: /'
echo "== R9b winsec 竞态族四发并集作差（量的，尺＝首发 raw 的 ^--- FAIL ＋复量三枚 .rednames.txt 各自并集）=="
{ grep '^--- FAIL' /d/tmp/wisp/20261020/352v1/scratch-ac2-host-pair-s1-treeB.txt | sed 's/^--- FAIL: //;s/ .*//' | grep 342r1;
  for i in 1 2 3; do git cat-file blob HEAD:$P/logs/ac2r$i-winsec-s1-treeB.rednames.txt | grep -v '^#'; done; } | sort -u > /d/tmp/wisp/20261020/352v1/scratch-wub.txt
{ grep '^--- FAIL' /d/tmp/wisp/20261020/352v1/scratch-ac2-host-pair-s2-treeL.txt | sed 's/^--- FAIL: //;s/ .*//' | grep 342r1;
  for i in 1 2 3; do git cat-file blob HEAD:$P/logs/ac2r$i-winsec-s2-treeL.rednames.txt | grep -v '^#'; done; } | sort -u > /d/tmp/wisp/20261020/352v1/scratch-wul.txt
echo "  四发 B 并集枚数=$(grep -c . /d/tmp/wisp/20261020/352v1/scratch-wub.txt) ⇄ L 并集枚数=$(grep -c . /d/tmp/wisp/20261020/352v1/scratch-wul.txt)"
echo "  ★并集作差：新增红=$(comm -13 /d/tmp/wisp/20261020/352v1/scratch-wub.txt /d/tmp/wisp/20261020/352v1/scratch-wul.txt | grep -c . || true) 消失红=$(comm -23 /d/tmp/wisp/20261020/352v1/scratch-wub.txt /d/tmp/wisp/20261020/352v1/scratch-wul.txt | grep -c . || true)"
echo "  （那三枚名的产地尺＝git grep -l 现算）"; git grep -l 'NeverDeletesAForeignFile342r1' HEAD -- 'internal/winsec/*.go' | sed 's/^HEAD:/    /'
echo "== R9c config 容器档那族：四发并集作差（量的，★⛔ 把名册件那行 # RED_NAMES=0 说明行当名＝先 grep -v 注释）=="
{ grep '^--- FAIL' /d/tmp/wisp/20261020/352v1/scratch-ac2-col-s1-treeB.txt | sed 's/^--- FAIL: //;s/ .*//' | grep -E 'HandDeletedSchema|NothingWritten';
  for i in 1 2 3; do git cat-file blob HEAD:$P/logs/ac2r$i-col-config-s1-treeB.rednames.txt | grep -v '^#'; done; } | sort -u > /d/tmp/wisp/20261020/352v1/scratch-cub.txt
{ grep '^--- FAIL' /d/tmp/wisp/20261020/352v1/scratch-ac2-col-s2-treeL.txt | sed 's/^--- FAIL: //;s/ .*//' | grep -E 'HandDeletedSchema|NothingWritten';
  for i in 1 2 3; do git cat-file blob HEAD:$P/logs/ac2r$i-col-config-s2-treeL.rednames.txt | grep -v '^#'; done; } | sort -u > /d/tmp/wisp/20261020/352v1/scratch-cul.txt
echo "  B 并集=$(grep -c . /d/tmp/wisp/20261020/352v1/scratch-cub.txt) ⇄ L 并集=$(grep -c . /d/tmp/wisp/20261020/352v1/scratch-cul.txt) 新增红=$(comm -13 /d/tmp/wisp/20261020/352v1/scratch-cub.txt /d/tmp/wisp/20261020/352v1/scratch-cul.txt | grep -c . || true) 消失红=$(comm -23 /d/tmp/wisp/20261020/352v1/scratch-cub.txt /d/tmp/wisp/20261020/352v1/scratch-cul.txt | grep -c . || true)"
echo "== R9d 底片⇄字节的两根绳（量的，尺＝sha256 对回 blob）：tsv 里 WINSEC_SHA256 那列＝哪棵树的哪个 blob =="
git cat-file blob bf509bb8:internal/winsec/winsec.go | sha256sum | cut -c1-24 | sed 's/^/  treeB 真身 sha256 前 24＝/'
git cat-file blob ecffb93b:internal/winsec/winsec.go | sha256sum | cut -c1-24 | sed 's/^/  treeL 真身 sha256 前 24＝/'
git cat-file blob HEAD:$P/logs/four-counts-errclass.tsv | grep -oE 'WINSEC_SHA256=[0-9a-f]{24}' | sort -u | sed 's/^/  底片所记＝/g'
echo "== R10 具名连带五枚（量的）=="; git cat-file blob HEAD:internal/winsec/private_fail_test.go | grep -n 'errors.Is(err, ErrNotSealable)' | sed 's/^/  /'
