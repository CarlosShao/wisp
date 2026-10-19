#!/bin/bash
# 两把收尾尺：
#  A) 23 枚没跑到的名＝哪种构建门挡掉的（逐枚读文件头的 //go:build 行）
#  B) 行尾的诚实尺＝已入库的 BLOB（⛔ 拿工作树字节当凭据：本机 core.autocrlf 会做 normalization）
set -u
REPO='/d/work/workspace/projects plans/Wisp'
cd "$REPO"
echo "core.autocrlf=$(git config core.autocrlf) eol=$(git config core.eol)"
echo '== A) 23 枚的门禁头 =='
for f in cmd/wisp/secret_dataroot_119b_test.go cmd/wisp/dispatch_hop_selfproof_133_test.go \
         cmd/wisp/resident_tray_stop_live_228_windows_test.go cmd/wisp/resident_ball_live_228_windows_test.go \
         cmd/wisp/resident_approval_live_246_windows_test.go cmd/wisp/resident_task_source_live_246_windows_test.go \
         cmd/wisp/resident_hotkey_live_258_windows_test.go; do
  printf '%-62s %s\n' "$f" "$(sed -n '1,3p' "$f" | grep '^//go:build' | head -1)"
done
echo '  （另两枚：面板死亡复量族＋WinLive，逐枚头）'
grep -rl 'func TestPanelDeathReplay319UnfixedReading\|func TestTicket255RealWindowWidthFollowsTheConfig\|func TestPanelHostWebViewChildrenExitWithinTwoSeconds_WinLive\|func TestSecondLaunchActivationRaisesThePanel319' --include='*_test.go' cmd/wisp | while read -r f; do
  printf '%-62s %s\n' "$f" "$(sed -n '1,4p' "$f" | grep '^//go:build' | head -1)"
done
echo '== B) 已入库 blob 的 CR 枚数 ⇄ 工作树 CR 枚数 =='
for f in $(git ls-files .scratch/wisp/probes/227/v2); do
  b=$(git cat-file blob "HEAD:$f" | tr -cd "\r" | wc -c)
  w=$(tr -cd "\r" < "$f" | wc -c)
  printf '%-52s blob_cr_bytes=%s worktree_cr_bytes=%s\n' "$f" "$b" "$w"
done
echo '== C) 分母尺的自身缺陷复核：字符串字面量里的 func Test（⛔ 真用例）=='
# 上一把尺＝grep -h "^func Test" cmd/wisp/*_test.go（⛔ 感知 Go 词法）⇒ 反引号里的夹具源码会被数进分母。
# 这把尺＝对那 23 枚逐名问：命中的那一行是不是落在裸字符串字面量里（od 数该行之前的反引号奇偶＝粗尺）。
grep -h '^func Test' cmd/wisp/*_test.go | sed 's/^func \(Test[A-Za-z0-9_]*\).*/\1/' | sort -u > /tmp/src2.txt
grep -E '^=== RUN   ' /d/tmp/wisp/20261019/227v2/full-B.log | sed 's/^=== RUN   //; s# .*##' | sort -u > /tmp/r2.txt
comm -23 /tmp/src2.txt /tmp/r2.txt > /tmp/miss2.txt
echo "missing=$(wc -l < /tmp/miss2.txt)"
while read -r n; do
  hit=$(grep -rn "^func $n(" --include='*_test.go' cmd/wisp | head -1)
  f=${hit%%:*}; ln=$(echo "$hit" | cut -d: -f2)
  # 该行之前（同文件内）出现的裸反引号枚数＝奇数 ⇒ 该行在字符串字面量内
  ticks=$(head -n "$ln" "$f" | tr -cd '`' | wc -c)
  if [ $((ticks % 2)) -eq 1 ]; then kind="IN-STRING-LITERAL (⛔ 真用例, 分母尺缺陷)"; else kind="real top-level func"; fi
  printf '%-58s %-58s ticks_before=%s %s\n' "$n" "$f" "$ticks" "$kind"
done < /tmp/miss2.txt
echo DONE
