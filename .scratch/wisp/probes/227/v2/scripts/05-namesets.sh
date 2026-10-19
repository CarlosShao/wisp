#!/bin/bash
# 名集尺：两发整包（编排者 A1053 的 full.log ⇄ 本腿 full-B.log）逐名作差。
# ⛔ 拿四数相减当结论（公共块：枚数⛔ 是单枚字节的函数）⇒ 只比名字集合＋解释枚数差。
# 四把尺各自的名册，全部锚行首（^--- X 或去缩进后固定串），⛔ 假设分隔符。
set -u
A=/d/tmp/wisp/20261019/orch227-b/full.log
B=/d/tmp/wisp/20261019/227v2/full-B.log
D=/d/tmp/wisp/20261019/227v2/namesets
mkdir -p "$D"
for pair in "orch:$A" "mine:$B"; do
  tag=${pair%%:*}; L=${pair#*:}
  grep '^=== RUN'       "$L" | sed 's/^=== RUN  *//' | sed 's/ .*//'   | sort > "$D/$tag-run.txt"
  grep '^--- PASS'      "$L" | sed 's/^--- PASS: //; s/ .*//'          | sort > "$D/$tag-pass-top.txt"
  grep '^    --- PASS'  "$L" | sed 's/^ *--- PASS: //; s/ .*//'        | sort > "$D/$tag-pass-sub.txt"
  grep '^--- FAIL'      "$L" | sed 's/^--- FAIL: //; s/ .*//'          | sort > "$D/$tag-fail-top.txt"
  grep '^    --- FAIL'  "$L" | sed 's/^ *--- FAIL: //; s/ .*//'        | sort > "$D/$tag-fail-sub.txt"
  grep '^--- SKIP'      "$L" | sed 's/^--- SKIP: //; s/ .*//'          | sort > "$D/$tag-skip-top.txt"
  grep '^    --- SKIP'  "$L" | sed 's/^ *--- SKIP: //; s/ .*//'        | sort > "$D/$tag-skip-sub.txt"
done
echo "== set sizes (lines per name file) =="
for k in run pass-top pass-sub fail-top fail-sub skip-top skip-sub; do
  printf '%-9s orch=%-6s mine=%-6s\n' "$k" "$(wc -l < "$D/orch-$k.txt")" "$(wc -l < "$D/mine-$k.txt")"
done
echo "== RUN name set: orch ⊎ mine (only-in-orch then only-in-mine) =="
comm -23 "$D/orch-run.txt" "$D/mine-run.txt" | sed 's/^/O-ONLY /'
comm -13 "$D/orch-run.txt" "$D/mine-run.txt" | sed 's/^/M-ONLY /'
echo "== top-level FAIL name set: orch ⊎ mine =="
comm -23 "$D/orch-fail-top.txt" "$D/mine-fail-top.txt" | sed 's/^/O-ONLY /'
comm -13 "$D/orch-fail-top.txt" "$D/mine-fail-top.txt" | sed 's/^/M-ONLY /'
echo "== SKIP top: orch ⊎ mine =="
comm -23 "$D/orch-skip-top.txt" "$D/mine-skip-top.txt" | sed 's/^/O-ONLY /'
comm -13 "$D/orch-skip-top.txt" "$D/mine-skip-top.txt" | sed 's/^/M-ONLY /'
echo "== duplicates inside one name set (same name run twice) =="
for t in orch mine; do
  for k in run pass-top fail-top skip-top; do
    d=$(uniq -d "$D/$t-$k.txt" | wc -l); printf '%s %s dupnames=%s\n' "$t" "$k" "$d"
  done
done
