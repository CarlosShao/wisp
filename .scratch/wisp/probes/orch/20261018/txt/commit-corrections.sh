#!/bin/sh
set -e
cd "/d/work/workspace/projects plans/Wisp"
T=.scratch/wisp/probes/orch/20261018/txt
F120=$(ls .scratch/wisp/issues/120-*.md)
F337=$(ls .scratch/wisp/issues/337-*.md)
echo "120=$F120 lines_before=$(wc -l < "$F120")"
echo "337=$F337 lines_before=$(wc -l < "$F337")"
cat "$T/120-tail-correction.md" >> "$F120"
cat "$T/337-tail-correction.md" >> "$F337"
echo "120 lines_after=$(wc -l < "$F120") numstat=$(git diff --numstat -- "$F120" | cut -f1,2)"
echo "337 lines_after=$(wc -l < "$F337") numstat=$(git diff --numstat -- "$F337" | cut -f1,2)"
echo "=== AC boxes must be unchanged ==="
for f in "$F120" "$F337"; do echo "$f open=$(grep -c '^[[:space:]]*- \[ \]' "$f") done=$(grep -c '^[[:space:]]*- \[x\]' "$f")"; done
printf '%s\n' "两处过期指认只追加更正（⛔ 改原文）：票 120 §2 那格里 `sealhandle_window_120_other_test.go:7/:21-23` 的指认不实（该件 unfixed／322 命中各 0、`log -S` 零枚提交）＝真过期只有 privateroute 那一处、已由 322-r2 的 ab49e9f9 改准（7+/7−、CR 0、winsec_other.go 未动）｜票 337 面 :13 那枚 `ToolOutcomeText(:1105)` ＝跟踪面零出现的幻影，真身＝internal/tools/bridge.go:1105 的 `agent.ToolOutcome{Text: why}`（:1106）＋:439 拼装那句｜★同时打我自己两枚：把腿那句"判语件也引了它"未经复跑写进 A1027（334-v1-verdict 实为 0 命中）＋起草时把路径写成不存在的 internal/agent/bridge.go ＝本仓"抄腿引文不复核"第 87 次、这回是我｜推送已授权仍未执行" > "$T/corr-msg.txt"
git add -- "$F120" "$F337" "$T/120-tail-correction.md" "$T/337-tail-correction.md" "$T/corr-msg.txt"
echo "=== staged audit (must be only mine) ==="
git diff --cached --name-only | grep -v -E "^\.scratch/wisp/(issues/120-|issues/337-|probes/orch/20261018/)" || echo "  only-mine"
git commit -F "$T/corr-msg.txt" -- "$F120" "$F337" "$T/120-tail-correction.md" "$T/337-tail-correction.md" "$T/corr-msg.txt"
git log -1 --format='%h %s' | cut -c1-60
echo "unpushed=$(git rev-list --count origin/dev..HEAD)"
