#!/bin/bash
cd '/d/work/workspace/projects plans/Wisp' || exit 9
echo '===== ① 0fb399c3 逐枚文件（body 自述 3 枚 ⇄ disk 现量 5 枚）====='
git show --name-only --format='' 0fb399c3 | grep '[^[:space:]]' | cat
echo '===== ③ aacde632（body 自述 8 枚 ⇄ disk 9 枚）====='
git show --name-only --format='' aacde632 | grep '[^[:space:]]' | nl
echo '===== ⑤ 06aa2abf（body 自述 1 报告+27 散件=28 ⇄ disk 41）====='
git show --name-only --format='' 06aa2abf | grep '[^[:space:]]' | awk '
  /readings-.*\.md/ {r++; next}
  /\/logs\// {l++; next}
  /\/sh\// {s++; next}
  {o++}
  END{printf "REPORT=%d LOGS=%d SH=%d OTHER=%d TOTAL=%d\n", r, l, s, o, r+l+s+o}'
echo '===== ⑦ 5c9402f8 里那枚报告件 readings-ac3 的第 45 行 ⇄ 票尾真 numstat ====='
sed -n '45p' .scratch/wisp/probes/352/r1/readings-ac3-discipline.md
echo "  ACTUAL_NUMSTAT=$(git show --numstat --format='' 1c787282 | awk '{print $1" 增 "$2" 删"}')"
echo '===== ⑥ 3ecdea30 CR 那味（利尺 tr 复量，对回 ⑥ 正文"8 行"）====='
for f in ac2-host-fiv-s1-treeB ac2-host-pair-s1-treeB; do
  echo "  $f CRBYTES=$(tr -cd '\r' < .scratch/wisp/probes/352/r1/logs/$f.raw.txt | wc -c)"
done
echo '===== 名册现量 ⇄ 笔⑨正文那三枚 8（＝当时 HEAD=8 的点时读数，非造假）====='
echo "  NOW_ANCHORED=$(git log --grep='^票 352 落地腿 352-r1' --format='%h' | grep -c . || true)"
echo "  (笔⑨ 8a10a7e3 是其时末笔⇒它自数 8；续跑加第 9 枚后=9＝handoff-number-reconcile.txt 已预告的漂)"
echo DONE
