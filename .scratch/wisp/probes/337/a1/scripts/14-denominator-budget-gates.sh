#!/bin/sh
# 337-a1：hop a 分母件（16）＋三形落点行的现长度尺（23）＋适用门禁 rc（90）
cd "/d/work/workspace/projects plans/Wisp" || exit 1
L=.scratch/wisp/probes/337/a1/logs
R=c11163c8

{
  echo "### hop a 分母核对 rev=$R（尺＝git grep -nE 'return ClassA, ' 单文件射程）"
  echo "全部 return ClassA 行＝$(git grep -nE 'return ClassA, ' $R -- internal/risk/blacklist.go | wc -l)"
  git grep -nE 'return ClassA, ' $R -- internal/risk/blacklist.go
  echo "走 aRuleReason＝$(git grep -nE 'return ClassA, aRuleReason' $R -- internal/risk/blacklist.go | wc -l)"
  echo "走裸字面串＝$(git grep -nE 'return ClassA, "' $R -- internal/risk/blacklist.go | wc -l)"
  echo "★差异具名：判语件 logs/80-chain-reach.log 第 2 行写「九枚 A case 走 aRuleReason」；现跑＝9 枚 A 档 case 里 7 枚走 aRuleReason、2 枚走裸字面串（blacklist.go:223 .git/config、:234 browser credential store）。该九枚是真分母，走那把拼装尺的是七枚。"
} > "$L/16-hop-a-denominator.log"

{
  echo "### AC#1 三形落点的现长度尺 rev=$R（★量＝每行字节数；⛔ 估改动量：本腿零写码，不报“要动几字节”）"
  echo "尺＝git show \$R:<path> | sed -n '<N>p' | wc -c（逐枚量）"
  m() { printf '%-58s %s bytes\n' "$1" "$(git show "$R:$2" | sed -n "$3p" | wc -c)"; }
  echo "-- 甲（本机日志）候选点位现长度 --"
  m "risk/blacklist.go:51 (var Logf sink 定义)" internal/risk/blacklist.go 51
  m "tools/bridge.go:73 (Options.Logf 字段)" internal/tools/bridge.go 73
  m "cmd/wisp/run.go:787 (Logf: rt.auditf 既有注入形)" cmd/wisp/run.go 787
  echo "甲-既有注入枚数＝cmd/wisp 里 Logf: 生产注入 $(git grep -nE '^[[:space:]]*Logf:' $R -- cmd/wisp | grep -v '_test\.go:' | wc -l) 枚（尺见 22）"
  echo "-- 乙（面板）候选点位现长度 --"
  m "panel/approval.go:48 (json:\"reason\" 键)" internal/panel/approval.go 48
  m "panel/approval.go:83 (Reason: decision.Reason)" internal/panel/approval.go 83
  m "cmd/wisp/panel_pump.go:75 (Reason: it.Decision.Reason)" cmd/wisp/panel_pump.go 75
  m "tools/mode.go:111 (n.Reason = d.Reason 写入)" internal/tools/mode.go 111
  echo "-- 丙（模型正文）候选点位现长度 --"
  m "risk/rules_gateway.go:82 (R3 那句)" internal/risk/rules_gateway.go 82
  m "tools/bridge.go:439 (Deny 支拼句)" internal/tools/bridge.go 439
  m "tools/bridge.go:1106 (Text: why)" internal/tools/bridge.go 1106
  m "agent/loop.go:745 (进对话正文)" internal/agent/loop.go 745
  echo "-- 那句实话自身的长度（拼装后的样本，从测试自陈里取形状）--"
  git show "$R:internal/risk/blacklist.go" | sed -n '329,331p' | wc -c
} > "$L/23-ac1-line-budget.log"

{
  echo "### 门禁（本腿零 .go：四数里两枚对 md/log 件不适用，具名说明；只在共享工作树跑）"
  t() { sh -c "$1" >/dev/null 2>&1; echo "$2 rc=$?"; }
  t "sh scripts/d22scan.sh" "d22scan.sh（★它内部 go run 的是 tools/d22scan 自己的扫描模块＝guard 本身，派单明文要求跑；⛔ 不碰根模块）"
  t "sh scripts/check-path-length-budget.sh" "check-path-length-budget.sh"
  echo "go vet ./<本腿的包>/ rc=不适用（具名：本腿零 .go 字节、零自己的包；跑 go vet 会在共享工作树真编译 internal/risk＝只读腿禁区，且本波该面写腿＝336-r1）"
  echo "gofumpt -l 对 md/log/sh：本腿非 .go 件不受 gofumpt 射程；★派单那条「你的新字节必须 LF」＋⛔ 顺手重排别腿文件，尺＝下面的 CR 计数"
  echo "★新字节 LF 尺（工作树侧 CR 计数，逐枚）："
  for f in .scratch/wisp/probes/337/a1/logs/*.log .scratch/wisp/probes/337/a1/scripts/*.sh; do
    printf '%-70s CR=%s\n' "$f" "$(grep -c "$(printf '\r')" "$f" || true)"
  done
  echo "★blob 侧（入库归一）CR 计数（core.autocrlf=true ⇒ blob 应恒 0）："
  git ls-files -s -- .scratch/wisp/probes/337/a1 | while read mode sha stage p; do
    printf '%-70s blob_CR=%s\n' "$p" "$(git cat-file blob "$sha" | tr -cd '\r' | wc -c)"
  done
} > "$L/90-gates-and-lf.log"

wc -c "$L/16-hop-a-denominator.log" "$L/23-ac1-line-budget.log" "$L/90-gates-and-lf.log"
