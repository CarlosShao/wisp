#!/bin/sh
# 收件尺更正：严格前缀尺会漏掉主题被编码打坏（逐字 ?）的那笔 ⇒ 改用"逐枚哈希点名"尺
cd "/d/work/workspace/projects plans/Wisp" || exit 9
CS="82c4a1d4 ab49e9f9 7f805583 4280bfaa 83ffd90d 6ef91d75 e73b8c35 9c372800 17e67cca 631a04ba"
echo "=== per-commit existence + subject integrity ==="
for c in $CS; do
  s=$(git log -1 --format="%s" "$c" 2>/dev/null) || { echo "$c MISSING"; continue; }
  case "$s" in
    "?"*) flag="SUBJECT-MANGLED" ;;
    *)  flag="ok" ;;
  esac
  n=$(git show --name-only --format='' "$c" | grep -vc '^$')
  echo "$c files=$n $flag"
done
echo "=== union of paths (all 10) ==="
for c in $CS; do git show --name-only --format='' "$c"; done | grep -v '^$' | sort -u > /tmp/c22-all.txt
echo "paths=$(wc -l < /tmp/c22-all.txt)"
echo "go=$(grep -c '\.go$' /tmp/c22-all.txt || true)"
echo "=== outside its 3 write faces ==="
grep -v -E '^internal/winsec/|^\.scratch/wisp/probes/322/r2/|^\.scratch/wisp/issues/322-' /tmp/c22-all.txt
echo "outside=$(grep -vc -E '^internal/winsec/|^\.scratch/wisp/probes/322/r2/|^\.scratch/wisp/issues/322-' /tmp/c22-all.txt || true)"
echo "=== how many commits would a strict-prefix ruler see? (the blind spot) ==="
echo "prefix-ruler=$(git log --all --grep='^票 322 收尾腿 322-r2' --format='%h' | wc -l)  actual=$(echo $CS | wc -w)"
