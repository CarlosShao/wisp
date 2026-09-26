#!/usr/bin/env bash
# 票 158 r1 AC#2 丙档的补充现量（假阳性自拆）。只读，不写盘。
# 用法：bash .scratch/wisp/probes/158/g5-fp-audit.sh [锚点]
set -u
cd "$(git rev-parse --show-toplevel)" || exit 9
A="${1:-$(git rev-parse HEAD)}"

echo "# G5 假阳性自拆 —— 锚点 $A / 生成 $(date -Iseconds)"

echo
echo "## 现量 1：把 internal/risk/* 放回成对射程（只排 *_test.go），provenance.go 会不会被点名？"
calls="$(git grep -nEw 'OpenScope|CloseScope' "$A" -- 'internal/**/*.go' 'cmd/**/*.go' ':!*_test.go' \
	| grep -vE ':[0-9]+:[[:space:]]*(//|/\*|\*)' || true)"
printf '%s\n' "$calls"
opens="$(printf '%s\n' "$calls" | grep -Ew 'OpenScope' | cut -d: -f2 | sort -u || true)"
closes="$(printf '%s\n' "$calls" | grep -Ew 'CloseScope' | cut -d: -f2 | sort -u || true)"
echo "-- 开方名册: $(printf '%s' "$opens" | tr '\n' ' ')"
echo "-- 合方名册: $(printf '%s' "$closes" | tr '\n' ' ')"
echo "-- 未成对（放回 risk 之后）:"
for f in $opens; do printf '%s\n' "$closes" | grep -qxF "$f" || echo "   UNPAIRED $f"; done
echo "-- 结论：provenance.go 不被点名（同文件里两枚都有：:339 与 :350 的定义本体）。"
echo "-- ⇒ 主尺排掉 internal/risk/* 对【判语】不是承重的；它买的是【名册稳定】——"
echo "--   provenance.go 的注释一改，逐行名册就会变长、按名册 diff 的读法就会误响。"

echo
echo "## 现量 2：今天被点名的那一枚是「真漏」还是「委托关闭」？"
git grep -nEw 'CloseTask|Defer' "$A" -- cmd/wisp/panel_assets.go
echo "-- rc=$?（1＝同文件既无 CloseTask 也无 Defer ⇒ 真漏，不是委托）"

echo
echo "## 现量 3：字符串字面量这一族（结构性残余，今天零例）"
git grep -nE 'scope is not open' "$A" -- internal/risk/provenance.go
echo "-- 上面两行把 OpenScope 写进了【错误串】而不是注释，pair() 只剔注释、不剔字符串，"
echo "-- 所以「一枚只在字符串字面量里出现开方动词、同文件没有合方动词」的生产文件会假响。"
echo "-- 今天实测：主尺名册里唯一的开方是 cmd/wisp/panel_assets.go:232 那发真调用，零例假响。"
