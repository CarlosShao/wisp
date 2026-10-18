#!/bin/sh
# 337-a1：AC#1 三形的"会动哪几枚文件/行"字节级现量 ＋ AC#4 成对台面包名册 ＋ AC#3 现状同形尺
cd "/d/work/workspace/projects plans/Wisp" || exit 1
L=.scratch/wisp/probes/337/a1/logs
R=c11163c8
{
  echo "### AC#1 现量续 rev=$R（★只量现状，⛔ 选形/推荐；所有尺剔 .scratch/**）"
  echo ""
  echo "== 两枚 sink 的分工（甲形落点候选，只点名）=="
  echo "甲-1 internal/risk/blacklist.go:51 var Logf func(...) ⇒ 生产注入点=0（17-hop-g）；今天 nil ⇒ logf 分支整条丢"
  echo "甲-2 internal/tools/bridge.go:73 Logf func(...) 字段 ⇒ host 侧生产注入枚数："
  git grep -nE '^[[:space:]]*Logf:' $R -- cmd/wisp | grep -v '_test\.go:'
  echo "  甲-2 生产注入枚数=$(git grep -nE '^[[:space:]]*Logf:' $R -- cmd/wisp | grep -v '_test\.go:' | wc -l)"
  echo "甲-3 两枚 sink 之间的接线枚数（risk.Logf= 某 auditf 这种形，生产）=0（尺＝git grep -nE \"risk\\.Logf[[:space:]]*=\" \$R -- \"*.go\" 剔 .scratch/_test 空输出）"
  echo "甲-4 桥自己的 b.logf 调用点（含不含 Gate 那句）＝"
  git grep -nE 'b\.logf\(' $R -- internal/tools/bridge.go | grep -v '_test\.go:'
  echo "  读数＝只有 :1313（auditf 包装），不打印 Gate 的 PathDecision.Reason ⇒ Gate 那句今天连 tools 侧那枚有注入的 sink 都没接"
  echo ""
  echo "== readBlacklist（BlacklistNote 的产线）今天的调用点（甲/乙都绕不开它）=="
  git grep -n 'readBlacklist' $R -- '*.go' | grep -v '^[^:]*:\.scratch/' | grep -v '_test\.go:'
  echo ""
  echo "== 乙形出向面的既有承载体（★无新 JSON 键／新方法名时）=="
  echo "乙-a panel 侧 reason 键：internal/panel/approval.go:48 json:\"reason\"（Decision.Reason verbatim，见 :47 注）"
  echo "乙-b 卡片视图字段来源：internal/panel/approval.go:83 Reason: decision.Reason（risk.Decision，非 tools.Decision）"
  echo "乙-c pump 那一跳：cmd/wisp/panel_pump.go:75 Reason: it.Decision.Reason（tools.Decision 的 Reason 字段）"
  echo "乙-d ★tools.Decision.Blacklist 整枚是否已进快照出向：internal/panel 非测试命中 Blacklist 的行数=$(git grep -n 'Blacklist' $R -- internal/panel | grep -v '_test\.go:' | wc -l)（=0 ⇒ 面板字节面今天不带 BlacklistNote）"
  echo "乙-e A 档 deny 今天会不会出现在面板：refusal.go 的 Deny 支＝"
  git grep -nA2 'if level == risk.Deny' $R -- internal/tools/refusal.go
  echo "乙-f 前端渲染 reason 的位置（★是否触碰 frontend/**）＝"
  git grep -nE 'view\.reason|reasonKnown' $R -- frontend/src/components/l2-approval-card.tsx
  echo ""
  echo "== 丙形：模型正文那一跳的承载体（bridge.go reject → loop）=="
  echo "丙-a bridge.go:439 拼拒绝句（dec.Reason ＝ R3 那句，含 canonical 绝对路径）"
  echo "丙-b bridge.go:1093-1094 若 dec.Reason 为空则用 why 顶；:1106 Text: why 进 ToolOutcome"
  echo "丙-c loop.go:745 toolResultMessage(c.ID, log.Text, ...) ⇒ 进对话正文＝模型上下文"
  echo "丙-d ★模型侧今天对绝对路径/控制字符有没有治理（对照票 264 那一族）："
  git grep -nE 'func sanitizeArgs' $R -- internal/tools
  echo "  读数＝internal/tools/bridge.go:1290 sanitizeArgs 只作用在 ArgsJSON（forensics 列，见 :1201），⛔ 不作用于进模型的 Text ⇒ 丙形把真路径带进模型正文这一跳今天没有拦网（读码推的；尺＝sanitizeArgs 的调用点枚数）"
  git grep -n 'sanitizeArgs(' $R -- internal/tools | grep -v '_test\.go:'
  echo ""
  echo "== AC#4 前置料：本票落地后会落到哪几包的成对台面（★只点名列，⛔ 跑）=="
  echo "甲（本机日志）＝internal/risk（sink 侧）＋cmd/wisp（注入侧）"
  echo "乙（面板）＝internal/tools ＋ internal/panel ＋ cmd/wisp；★frontend/** 归机主外部 agent（本编队零写面）"
  echo "丙（模型正文）＝internal/tools ＋ internal/agent"
  echo "共同面＝cmd/wisp 与 internal/panel/** 本波零写腿（派单编队规则）；internal/risk/** 本波唯一写腿=336-r1 ⇒ 本票落地需按波次排他"
  echo "go test ./cmd/wisp 本波唯一持跑者=228-r3c"
  echo ""
  echo "== AC#3 前置料：门拦下文 vs 权限判定文，今天同不同形 =="
  echo "门拦下（bridge 拼的）＝'该目标被绝对禁止访问（R3 A 档，任何授权都不可豁免）：' + dec.Reason"
  echo "权限判定（R2）＝'R2: 目标路径在授权目录之外: %s'"
  git grep -n 'R2: 目标路径在授权目录之外' $R -- internal/risk/rules_gateway.go
  echo "钉'不得同串'的既有用例＝bridge_junction_windows_test.go:310（'红队拒绝不得与不在授权目录内共用措辞'）；票面点名的 242/244 族真名＝"
  git -c core.quotePath=false ls-tree -r --name-only $R | grep -v '^\.scratch' | grep -E '_242|242_|ticket242|244' | head -10
  echo "★aRuleReason 那句与 R3 那句是否同串＝不同（一为英文 fail-closed 长话，一为 R3 中文短句）；本腿只量现状⛔ 提改法"
} > "$L/22-ac1-ac3-ac4-measure.log"
wc -c "$L/22-ac1-ac3-ac4-measure.log"
echo "cr=$(grep -c "$(printf '\r')" "$L/22-ac1-ac3-ac4-measure.log" || true)"
