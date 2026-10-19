#!/bin/sh
# 342-v2 必答(v)：错误分类差消费者侧现读名册。
# 尺形：① 全仓 git grep 'fs.ErrNotExist' 逐枚点名，按〔产码/测试/别腿散件〕三分；
#      ② 只看真依赖 winsec 返回错误的：先取 PrivateFile/PrivateFileExclusive/SealFile 调用者名册，
#         再逐枚现读其错误处理块里有无 errors.Is(..., fs.ErrNotExist) / os.IsNotExist。
cd '/d/work/workspace/projects plans/Wisp' || exit 9
OUT=/d/tmp/wisp/20261019/342v2
LOG=$OUT/logs/v01-errnotexist-consumers.txt
mkdir -p "$OUT/logs"
{
echo '### A. 全仓 fs.ErrNotExist 命中（尺＝git grep -n "fs.ErrNotExist" -- "*.go"）'
echo '###    分类：PROD=非 _test.go 且不在 .scratch/ ；TEST=*_test.go；SCRATCH=.scratch/**（别腿散件，⛔ 生产消费者）'
git grep -n 'fs\.ErrNotExist' -- '*.go' > "$OUT/errnotexist-raw.txt"
echo "TOTAL_HITS=$(wc -l < "$OUT/errnotexist-raw.txt")"
awk -F: '{print $1}' "$OUT/errnotexist-raw.txt" | sort | uniq -c | sort -rn > "$OUT/errnotexist-byfile.txt"
grep -v '_test\.go:' "$OUT/errnotexist-raw.txt" | grep -v '^\.scratch/' > "$OUT/errnotexist-prod.txt"
echo "PROD_HITS=$(wc -l < "$OUT/errnotexist-prod.txt")"
grep -c '^\.scratch/' "$OUT/errnotexist-raw.txt" | sed 's/^/SCRATCH_HITS=/' || true
echo '--- PROD 名册逐枚（file:line + 该函数名）---'
cat "$OUT/errnotexist-prod.txt"
echo
echo '### B. winsec 三枚入口的生产调用者名册（尺＝git grep -n "winsec\.\(PrivateFile\|PrivateFileExclusive\|SealFile\)" -- "*.go"）'
git grep -n 'winsec\.\(PrivateFile\|PrivateFileExclusive\|SealFile\|RemoveUnlinked\)' -- '*.go' | grep -v '^\.scratch/' | grep -v '_test\.go:' > "$OUT/winsec-callers-prod.txt"
cat "$OUT/winsec-callers-prod.txt"
echo
echo '### C. 逐枚调用者的错误处理块：有无对 fs.ErrNotExist / os.IsNotExist 的判定（±14 行窗口）'
while IFS= read -r line; do
  f=${line%%:*}; rest=${line#*:}; ln=${rest%%:*}
  echo "===== CALLER $f:$ln ====="
  awk -v s=$((ln-3)) -v e=$((ln+16)) 'NR>=s&&NR<=e{printf "%d\t%s\n", NR, $0}' "$f"
done < "$OUT/winsec-callers-prod.txt"
} > "$LOG" 2>&1
# ★净化：写面⛔ 出现被禁字面字符（码位 U+00A7）。逐字 grep 输出里若带那枚字符（例：抄到的源码注释
# 指向 SPEC 的分节），一律换成其码位名；其余字节一字不动。正控＝同一枚 U+00B7（腿主题常用间隔号）⛔ 被动。
perl -i -pe 's/\xc2\xa7/U+00A7/g' "$LOG"
echo "SANITIZED_CTA7_HITS=$(grep -cP "\xc2\xa7" "$LOG" || true)"
wc -l < "$LOG" | sed 's/^/LOG_LINES=/'
