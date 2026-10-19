#!/bin/sh
# 342-v2 必答(v) 第二把尺：每一枚 fs.ErrNotExist / os.IsNotExist 判定点的 err **来源**是谁。
# 只判一件事＝有没有一枚真消费者把 winsec.PrivateFile / PrivateFileExclusive / SealFile 的返回错误
# 喂进 errors.Is(err, fs.ErrNotExist)（或 os.IsNotExist）。逐枚点名＋±6 行窗口，⛔ 结论在本件里。
cd '/d/work/workspace/projects plans/Wisp' || exit 9
OUT=/d/tmp/wisp/20261019/342v2
LOG=$OUT/logs/v03-classify-decision-points.txt
{
echo '### 尺＝git grep -n "errors.Is(.*fs.ErrNotExist\|os.IsNotExist" -- 产码（⛔ _test、⛔ .scratch）'
git grep -n 'errors\.Is([^)]*fs\.ErrNotExist\|os\.IsNotExist' -- '*.go' ':!**/*_test.go' ':!.scratch/**' > "$OUT/decision-points.txt"
echo "DECISION_POINTS=$(wc -l < "$OUT/decision-points.txt")"
cat "$OUT/decision-points.txt"
echo
echo '### 逐枚窗口（±6 行）＝看 err 从哪来'
while IFS= read -r line; do
  f=${line%%:*}; rest=${line#*:}; ln=${rest%%:*}
  echo "===== POINT $f:$ln ====="
  awk -v s=$((ln-6)) -v e=$((ln+4)) 'NR>=s&&NR<=e{printf "%d\t%s\n", NR, $0}' "$f"
done < "$OUT/decision-points.txt"
echo
echo '### 谁把 winsec 入口的返回错误**存下来**再判（尺＝winsec 调用赋值行 ±12 行窗口内有无 NotExist 判定）'
for f in $(git grep -l 'winsec\.\(PrivateFile\|PrivateFileExclusive\|SealFile\)' -- '*.go' ':!.scratch/**'); do
  echo "## FILE $f"
  grep -n 'fs.ErrNotExist\|os.IsNotExist' "$f" || echo "   (no NotExist decision in this file)"
done
} > "$LOG" 2>&1
echo "RC=$? LINES=$(wc -l < "$LOG")"
