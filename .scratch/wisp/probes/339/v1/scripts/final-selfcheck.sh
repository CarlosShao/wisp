#!/bin/bash
# 339-v1 台件⑫：交件自证三把尺（⛔ 只说"已完成"＝全部现跑、逐枚落盘）。
set -u
REPO="/d/work/workspace/projects plans/Wisp"
cd "$REPO" || exit 9
LOG=".scratch/wisp/probes/339/v1/logs/v1-final-selfcheck.txt"
{
  echo "== 尺①＝本腿逐笔 git show --name-only ⊆ 三处写面；.go 命中数（正控＝同尺打 339-r1 ⇒ 2）=="
  MINE=$(git log --all --format='%H %s' | grep -E '^[0-9a-f]{40} 票 339 终裁腿 339-v1' | awk '{print $1}' | sort -u)
  echo "my_commit_count=$(printf '%s\n' "$MINE" | grep -c . || true)"
  printf '%s\n' "$MINE" | while read -r h; do [ -z "$h" ] && continue; echo "  commit=$h"; git show --name-only --format='    %s' "$h" | sed 's/^/      /'; done
  echo "  --- 本腿全部路径去重 ---"
  ALLP=$(for h in $MINE; do git show --name-only --format='' "$h"; done | grep -v '^$' | sort -u)
  echo "$ALLP" | sed 's/^/    /'
  echo "  my_paths=$(printf '%s\n' "$ALLP" | grep -c . || true)"
  echo "  my_go_hits=$(printf '%s\n' "$ALLP" | grep -c '\.go$' || true) (必须 0)"
  echo "  my_outside_faces=$(printf '%s\n' "$ALLP" | grep -vcE '^docs/evidence/s1/339-v1-verdict\.md$|^\.scratch/wisp/probes/339/v1/|^\.scratch/wisp/issues/339-' || true) (必须 0)"
  echo "  --- 正控：同一把尺打实现腿 339-r1 的 10 笔并集 ---"
  R1=$(git log --all --format='%H %s' | grep -E '^[0-9a-f]{40} 票 339 落地腿 339-r1' | awk '{print $1}' | sort -u)
  echo "  r1_commit_count=$(printf '%s\n' "$R1" | grep -c . || true)"
  R1P=$(for h in $R1; do git show --name-only --format='' "$h"; done | grep -v '^$' | sort -u)
  echo "  r1_paths=$(printf '%s\n' "$R1P" | grep -c . || true)"
  echo "  r1_go_hits=$(printf '%s\n' "$R1P" | grep -c '\.go$' || true) (正控＝2)"
  echo
  echo "== 尺②＝判语件绝对路径＋wc -c＋行数＋占位符尺 0 命中（正控＝同尺打 HANDOVER.md 非零）=="
  VF=docs/evidence/s1/339-v1-verdict.md
  echo "  abs=$PWD/$VF bytes=$(wc -c < "$VF") lines=$(wc -l < "$VF")"
  echo "  placeholder_hits=$(grep -cE '待填|待补|待|TBD|TODO' "$VF" || true)  (★具名：裸"待"字会误伤正文里"待人拍板"这类引用，故单列)"
  echo "  placeholder_hits_narrow=$(grep -cE '待填|待补|TBD|TODO' "$VF" || true)"
  echo "  posctl_HANDOVER=$(grep -cE '待填|待补|TBD|TODO' docs/reports/HANDOVER.md || true) (尺灵)"
  echo "  u00a7_in_verdict=$(perl -0777 -ne 'BEGIN{$c=0} while(/(\xc2\xa7)/g){$c++} END{print $c+0}' "$VF") (正控 HANDOVER=$(perl -0777 -ne 'BEGIN{$c=0} while(/(\xc2\xa7)/g){$c++} END{print $c+0}' docs/reports/HANDOVER.md))"
  echo "  CR_in_verdict=$(tr -cd '\r' < "$VF" | wc -c) (必须 0＝纯 LF)"
  echo
  echo "== 尺③＝证据件逐枚非 0 字节＋行数上限＋票框终量尺 =="
  echo "  zero_byte_files=$(find .scratch/wisp/probes/339/v1 -type f -size 0 | wc -l) (必须 0)"
  echo "  file_count=$(find .scratch/wisp/probes/339/v1 -type f | wc -l)"
  echo "  over160_lines=$(for f in $(find .scratch/wisp/probes/339/v1 -type f -name '*.md' -o -type f -name '*.txt' -o -type f -name '*.log'); do n=$(wc -l < "$f"); [ "$n" -gt 160 ] && echo "$n $f"; done | wc -l) (正文件 ≤160 口径；body.log 类原始读数件除外见下)"
  for f in $(find .scratch/wisp/probes/339/v1/logs -type f | sort); do echo "    $(wc -c < "$f")B $(wc -l < "$f")L $f"; done
  TICK=.scratch/wisp/issues/339-ticket334-anchor-miss-instrument-runs-on-linux-but-its-fake-home-is-windows-only.md
  echo "  ticket339_unchecked=$(grep -cE '^- \[ \] \*\*AC#' "$TICK" || true) checked=$(grep -cE '^- \[x\]' "$TICK" || true) lines=$(wc -l < "$TICK")"
  echo "  ticket339_status_line_head=$(sed -n '3p' "$TICK" | cut -c1-40)"
  echo "  git_log_all_verdict_leg=$(git log --all --oneline --grep='票 339 终裁腿' | wc -l)"
} > "$LOG" 2>&1
cat "$LOG"
