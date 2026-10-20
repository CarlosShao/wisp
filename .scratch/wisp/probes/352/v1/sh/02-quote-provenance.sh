#!/bin/sh
# 352-v1 终裁腿·出处对撞尺：报告件标称"逐字"的引文 ⇄ 它标注的那枚入库底片（逐字节）
# 用法：sh .scratch/wisp/probes/352/v1/sh/02-quote-provenance.sh
cd '/d/work/workspace/projects plans/Wisp' || exit 9
P=.scratch/wisp/probes/352/r1
D=/d/tmp/wisp/20261020/352v1/q2
mkdir -p $D

echo '== A) AC#0 报告件标称"逐字"的六行引文，逐枚在底片里逐字节找（尺＝grep -qF 整行）=='
echo '   报告件＝readings-ac0-before.md :53-56（POSIX 档）⇄ :61-64（Windows 档）'
echo "   射程＝HEAD 下 probes/352/** 全件（含 md／sh／tsv，⛔ 只扫 logs）＝$(git ls-tree -r --name-only HEAD -- $P | grep -c .) 枚"
git cat-file blob HEAD:$P/readings-ac0-before.md > $D/report.txt
for n in 54 55 56 62 63 64; do
  sed -n "${n}p" $D/report.txt | sed 's/^ *//' > $D/q.txt
  hits=''
  for f in $(git ls-tree -r --name-only HEAD -- $P | grep -v 'readings-ac0-before.md$'); do
    if git cat-file blob HEAD:$f | grep -qF -f $D/q.txt; then hits="$hits ${f##*/}"; fi
  done
  tag=$(sed -n "${n}p" $D/report.txt | grep -oE 'R[0-9]+' | head -1)
  printf '  report:%s [%s] 命中件=%s\n' "$n" "$tag" "${hits:-  (零命中＝该行不在任何入库件里)}"
done

echo '== B) 拼接坐实：被引那枚 Windows R1 行的 run 根⇄测试目录名分属两枚不同底片 =='
for f in ac0-host-base-s1-mut-1-treeB ac0-host-base-s2-mut-2-treeB ac0-host-base1-mut-treeB; do
  printf '  %-32s ' "$f"
  git cat-file blob HEAD:$P/logs/$f.raw.txt | grep -m1 'R1 parent-missing' \
    | grep -oE '352r1-[0-9]+-[0-9]+-[0-9]+|TestP352R1ErrClassify[0-9]+' | sort -u | tr '\n' ' '
  echo
done
printf '  %-32s ' "report:62 所引那一行"
sed -n '62p' $D/report.txt | sed 's/^ *//' > $D/q62.txt
grep -oE '352r1-[0-9]+-[0-9]+-[0-9]+|TestP352R1ErrClassify[0-9]+' $D/q62.txt | sort -u | tr '\n' ' '; echo '(两个形各自只在一枚底片里出现＝拼接)'

echo '== C) 那三枚底片里 R1⇄R4 的分类位（＝判据承载位，逐枚现量）=='
bits() { grep -oE 'isErrNotExist=[a-z]+|isNotSealable=[a-z]+|isErrExist=[a-z]+|isUnresolved=[a-z]+' | paste -sd' ' -; }
for f in ac0-host-base-s1-mut-1-treeB ac0-host-base-s2-mut-2-treeB ac0-host-base1-mut-treeB; do
  printf '  %-32s ' "$f R1"; git cat-file blob HEAD:$P/logs/$f.raw.txt | grep -m1 'R1 parent-missing' | bits
  printf '  %-32s ' "   同底片 R4"; git cat-file blob HEAD:$P/logs/$f.raw.txt | grep -m1 'R4 SealFile missing-leaf' | bits
done
printf '  %-32s ' "report:62 所引那一行 R1"; sed -n '62p' $D/report.txt | bits
printf '  %-32s ' "report:64 所引那一行 R4"; sed -n '64p' $D/report.txt | bits

echo '== D) AC#1 报告件标称"逐字"的六行（:50-52 POSIX ⇄ :58-60 Windows）同尺对撞 =='
git cat-file blob HEAD:$P/readings-ac1-landed.md > $D/r1report.txt
for n in 50 51 52; do
  sed -n "${n}p" $D/r1report.txt | sed 's/^ *//' > $D/q.txt
  git cat-file blob HEAD:$P/logs/ac1-col-s2-mut-2-treeL.raw.txt | grep -qF -f $D/q.txt && echo "  report:$n 逐字节在 ac1-col-s2-mut-2-treeL ＝YES" || echo "  report:$n ＝NO"
done
for n in 58 59 60; do
  sed -n "${n}p" $D/r1report.txt | sed 's/^ *//' > $D/q.txt
  git cat-file blob HEAD:$P/logs/ac1-host-s2-mut-2-treeL.raw.txt | grep -qF -f $D/q.txt && echo "  report:$n 逐字节在 ac1-host-s2-mut-2-treeL ＝YES" || echo "  report:$n ＝NO"
done

echo '== E) probe-copy-of-stage.txt：件内首行自我声明 ⇄ 它在 HEAD 的跟踪状态 =='
git cat-file blob HEAD:$P/probe-copy-of-stage.txt | sed -n '1p' | sed 's/^/  首行＝/'
git ls-files --error-unmatch $P/probe-copy-of-stage.txt >/dev/null 2>&1 && echo "  git ls-files 跟踪状态＝TRACKED（与首行声明相反）" || echo "  跟踪状态＝未跟踪"
git cat-file blob HEAD:$P/probe-copy-of-stage.txt > $D/pc.txt
git cat-file blob HEAD:$P/probe/probe-352r1-errclass_test.go.txt > $D/pd.txt
echo "  与 probe/ 那枚的逐字节差行数＝$(diff $D/pc.txt $D/pd.txt | grep -c '^[<>]' || true)（0＝同一件的两份拷贝）"
echo "  入库笔数尺＝$(git log --format='%h' --diff-filter=A -- $P/probe-copy-of-stage.txt | grep -c .) 笔首次加入"

echo '== F) 笔⑥"8 行 CRLF＝原始读数照实保留"的层面坐实（★CR 一律对入库 blob 量）=='
for f in ac2-host-pair-s1-treeB ac2-host-fiv-s1-treeB ac2-col-s1-treeB; do
  blobcr=$(git cat-file blob HEAD:$P/logs/$f.raw.txt | tr -cd '\r' | wc -c)
  wrcr=$(tr -cd '\r' < $P/logs/$f.raw.txt | wc -c)
  printf '  %-26s 入库blob CR=%s  工作树 CR=%s\n' "$f" "$blobcr" "$wrcr"
done
echo "  git core.autocrlf＝$(git config --get core.autocrlf)  ⇄ .gitattributes 那味＝$(git cat-file blob HEAD:.gitattributes | grep -E '^\.gitignore|text=auto|^\* ' | head -1 | tr -d '\n')"
echo "  结论尺＝那 8 枚 CR 只活在导出树⇄工作树那一层；入库 blob＝0＝被 git 归一化吞掉（⛔ 腿的净化，也⛔"照实保留"）"
echo "  判据承载性＝那 8 行是否 ^--- FAIL／^=== RUN／^--- PASS 形：$(while IFS= read -r l; do case $l in '--- FAIL'*|'=== RUN'*|'--- PASS'*) echo x;; esac; done < <(grep 'Successfully processed' $P/logs/ac2-host-fiv-s1-treeB.raw.txt) | grep -c x || true) 枚（0＝⛔ 判据行）"

echo '== G) 名册尺随 HEAD 漂（＝笔⑦/笔⑨那类点时读数的独立复量）=='
for c in 0fb399c3 ecffb93b aacde632 c16ff47e 06aa2abf 3ecdea30 5c9402f8 1c787282 8a10a7e3 e74762b5 a651c086; do
  printf '  到 %s 为止锚定枚数=%s 路径并集=%s\n' "$c" \
    "$(git log --grep='^票 352 落地腿 352-r1' $c --format='%h' | grep -c .)" \
    "$(for k in $(git log --grep='^票 352 落地腿 352-r1' $c --format='%h'); do git show --name-only --format='' $k; done | grep -v '^$' | sort -u | grep -c .)"
done
echo '  ⇄ 入库件 readings-ac3-discipline.md 当时所记＝'
git cat-file blob HEAD:$P/readings-ac3-discipline.md | grep -nE '5 笔|65 枚' | sed 's/^/    /' | cut -c1-120
echo '  ⇄ discipline.txt 在笔⑦那一笔里的值＝'
git cat-file blob 5c9402f8:$P/logs/discipline.txt | grep -E 'ANCHORED_COUNT|PATHS_TOTAL' | sed 's/^/    /'
echo 'DONE'
