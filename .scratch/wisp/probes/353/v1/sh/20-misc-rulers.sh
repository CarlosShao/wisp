#!/bin/bash
# 353-v1 R20 = 杂项硬尺：⛔racer ⇄ ErrExist ⇄ 笔④纯注释 ⇄ deleted 全 0 ⇄ ⛔跑 cmd/wisp ⇄ 名册三把 ⇄ 脏尺
cd '/d/work/workspace/projects plans/Wisp' || exit 9
B=/d/tmp/wisp/20261020/353v1/blob
L=/d/tmp/wisp/20261020/353r1/logs
echo "###### 1) 仪器件 ⛔ racer ⇄ ⛔ 计时（量的；正控＝同尺在 leafrace_342r2 命中）"
for w in 'go func' 'time.Sleep' 'sync.WaitGroup' 'atomic.'; do
  printf '  nametopup[%-16s]=%s   leafrace正控=%s\n' "$w" \
   "$(grep -cF "$w" $B/HEAD-nametopup_353_windows_test.go.txt || true)" \
   "$(git cat-file blob HEAD:internal/winsec/leafrace_342r2_windows_test.go | grep -cF "$w" || true)"
done
echo "###### 2) fs.ErrExist 那一支两侧同串（量的）"
for c in a59f5630 HEAD; do
  printf '  %s fs.ErrExist in winsec.go=%s in winsec_windows.go=%s\n' "$c" \
   "$(git cat-file blob "$c:internal/winsec/winsec.go" | grep -cF 'fs.ErrExist' || true)" \
   "$(git cat-file blob "$c:internal/winsec/winsec_windows.go" | grep -cF 'fs.ErrExist' || true)"
done
echo "  ERR_EXIST_BRANCH_LINES_IDENTICAL=$(diff -u0 <(git cat-file blob a59f5630:internal/winsec/winsec.go | grep -nF 'fs.ErrExist') <(git cat-file blob HEAD:internal/winsec/winsec.go | grep -nF 'fs.ErrExist') >/dev/null && echo YES || echo NO)"
echo "###### 3) 笔④⇄笔② 之间＝pinnedleaf 是否纯注释（尺＝剥行首注释后的可执行行集合作差）"
strip_exec() { git cat-file blob "$1:internal/winsec/$2" | grep -vE '^[[:space:]]*//' | grep -vE '^[[:space:]]*$'; }
for f in pinnedleaf_342_windows.go pinnedleaf_342_unix.go; do
  diff -u0 <(strip_exec fa395b33 "$f") <(strip_exec HEAD "$f") > /dev/null && echo "  $f EXEC_LINES_FA395_vs_HEAD=IDENTICAL" || { echo "  $f EXEC DIFF:"; diff -u0 <(strip_exec fa395b33 "$f") <(strip_exec HEAD "$f") | cut -c1-120; }
done
echo "  正控：同一把尺打在 fa395b33 ⇄ a59f5630（乙本体）上必⛔ 同＝"
diff -u0 <(strip_exec fa395b33 pinnedleaf_342_windows.go) <(strip_exec a59f5630 pinnedleaf_342_windows.go) | grep -Ec '^[+-][^+-]' | sed 's/^/   execdiff_B vs 乙=/ | sed 's/=$/ 行/'
echo "###### 4) deleted 非零（正控先行）"
printf '  正控串存在性 deleted= 总命中=%s\n' "$(grep -rho 'deleted=' $L | wc -l)"
printf '  deleted 非零命中=%s\n' "$(grep -rhoE 'deleted=[1-9][0-9]*' $L | wc -l)"
grep -rhoE 'deleted=[1-9][0-9]*' $L | sort | uniq -c | sed 's/^/   /'
echo "###### 5) ⛔ 跑 ./cmd/wisp（尺＝底片里的包名＋尺脚本本体）"
grep -rl 'cmd/wisp' /d/tmp/wisp/20261020/353r1/sh | sed 's/^/   尺脚本点名 cmd\/wisp: /'
grep -rho 'github.com/CarlosShao/wisp/cmd/wisp' $L | wc -l | sed 's/^/   底片 cmd\/wisp 包名命中=/'
printf '  正控（底片里 winsec 包名命中枚数）=%s\n' "$(grep -rho 'github.com/CarlosShao/wisp/internal/winsec' $L | wc -l)"
echo "###### 6) 名册三把尺（锚定 %s ⇄ 非锚定子串 ⇄ 单次导出 index==1）"
git log --all --format='%s' > /d/tmp/wisp/20261020/353v1/subj-v1.txt
git log --all --format='%H%x09%s' > /d/tmp/wisp/20261020/353v1/subjhash-v1.txt
printf '  ruler1 anchored=%s\n' "$(grep -c '^票 353 落地腿 353-r1' /d/tmp/wisp/20261020/353v1/subj-v1.txt || true)"
printf '  ruler2 nonanchored[票 353 落地腿 353-r1 anywhere]=%s\n' "$(grep -c '票 353 落地腿 353-r1' /d/tmp/wisp/20261020/353v1/subj-v1.txt || true)"
printf '  ruler3 single-export index==1=%s\n' "$(awk -F'\t' 'index($2,"票 353 落地腿 353-r1")==1{c++} END{print c+0}' /d/tmp/wisp/20261020/353v1/subjhash-v1.txt)"
grep '^票 353 落地腿 353-r1' /d/tmp/wisp/20261020/353v1/subj-v1.txt | cut -c1-58 | sed 's/^/   /'
printf '  正控 异锚 ^票 353 普查腿 353-a1=%s  ^票 353 终裁腿 353-v1=%s\n' "$(grep -c '^票 353 普查腿 353-a1' /d/tmp/wisp/20261020/353v1/subj-v1.txt||true)" "$(grep -c '^票 353 终裁腿 353-v1' /d/tmp/wisp/20261020/353v1/subj-v1.txt||true)"
echo "###### 7) tracked 脏逐枚列名＋mtime（量的）"
git status --porcelain | while read -r st p; do
  f=$(echo "$p" | sed 's/^"//; s/"$//')
  printf '   [%s] %s  mtime=%s\n' "$(echo $st)" "$f" "$(date -r "$(date +%s)" 2>/dev/null >/dev/null; stat -c '%y' "$f" 2>/dev/null | cut -c1-19)"
done
echo "###### 8) 那 11 枚盯§C 的件在位尺（⛔ 指不存在的件当凭据；尺＝cat-file -e 逐枚）"
for f in ancestoropen_342_windows_test.go ancestoropen_342_other_test.go ancestorunlink_348_windows_test.go ancestorunlink_348_other_test.go ancestorworld_342_348_windows_test.go ancestorworld_342_348_other_test.go leafrace_342r2_windows_test.go openleaf_overwrite_322_windows_test.go openleaf_322_other_test.go posix_seal_window_120_test.go sealhandle_privateroute_120_test.go; do
  git cat-file -e "HEAD:internal/winsec/$f" 2>/dev/null && echo "   EXISTS $f" || echo "   ABSENT $f"
done
echo "-- 名册尺：internal/winsec 下真正含这些族名的件逐枚（防我指错名） --"
git ls-tree --name-only HEAD internal/winsec/ | sed 's|internal/winsec/||' | grep -E 'ancestor|leafrace|openleaf|posix_seal|sealhandle' | sed 's/^/   /'
