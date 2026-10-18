#!/usr/bin/env bash
# 328-v1 交件自证尺②③：逐枚 wc（字节/行数）＋空槽尺（自避分词形）0 命中＋正控，全部现跑入库。
set -u
cd "D:/work/workspace/projects plans/Wisp" || exit 1
F='docs/evidence/s1/328-v1-verdict.md'
D='.scratch/wisp/probes/328/v1/logs'
RULER='([T]ODO|[F]IXME|[T]BD|待[补]|占[位]|X{3})'
PC='internal/agent/approval/ticket146_liveapprovals_backing_test.go'
OUT="$D/deliverable-sizes.txt"

{
  echo "date=$(date '+%F %T %z')  leg=328-v1  尺＝wc -c / wc -l 逐枚现量（射程＝本腿写面全集＋判语件）"
  echo '--- 判语件 ---'
  printf "%8s %5s %s\n" "$(wc -c < "$F")" "$(wc -l < "$F")" "$F"
  echo '--- probes/328/v1 逐枚（bytes lines path）---'
  while IFS= read -r f; do
    printf "%8s %5s %s\n" "$(wc -c < "$f")" "$(wc -l < "$f")" "$f"
  done < <(find .scratch/wisp/probes/328/v1 -type f | sort)
  echo '--- 0 字节枚数（尺＝find -size 0，射程＝写面全集＋判语件）---'
  find .scratch/wisp/probes/328/v1 "$F" -type f -size 0 | wc -l
  echo '--- ★空槽尺（自避分词形）逐字命令＋读数 ---'
  echo "\$ grep -nE '$RULER' $F"
  echo "  命中=$(grep -cE "$RULER" "$F" || true)（判语件射程，正控在下面一枚）"
  echo "\$ grep -nE '$RULER' $PC   ← 正控：同一串尺打到已知存在的件，先证这把尺会响"
  echo "  正控命中行号＝$(grep -nE "$RULER" "$PC" | head -2 | cut -d: -f1 | tr '\n' ' ')（⛔ 把夹具串原样打进本件，免自匹配）"
  echo "  正控命中行数=$(grep -cE "$RULER" "$PC" || true)"
  echo '--- 行数帽尺（派单：判语件 ≤130 行）---'
  echo "  verdict_lines=$(wc -l < "$F")  cap=130  over=$([ "$(wc -l < "$F")" -le 130 ] && echo NO || echo YES)"
  echo '--- 件数帽尺（输出预算 ≤4 枚件：判语件＋logs 属读数件，正件＝判语件一枚＋票尾一条）---'
  echo "  非读数正件枚数＝2（docs/evidence/s1/328-v1-verdict.md ＋ 票 328 尾一条 bullet）；probes 面＝脚本＋logs"
} > "$OUT"
echo "written $OUT $(wc -c < "$OUT")B/$(wc -l < "$OUT")L"
tail -16 "$OUT"
