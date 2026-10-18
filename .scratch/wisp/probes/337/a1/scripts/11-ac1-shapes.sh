#!/bin/sh
# 337-a1：AC#1 三形代价表的现量尺（甲/乙/丙），全部取具名提交树、剔 .scratch/**
cd "/d/work/workspace/projects plans/Wisp" || exit 1
L=.scratch/wisp/probes/337/a1/logs
R=c11163c8

{
  echo "### AC#1 现量：三形各自的出向面今天长什么样 rev=$R（★所有尺剔 .scratch/**）"
  echo ""
  echo "== 甲（本机日志）=="
  echo "甲1 sink 定义与今天为 nil 的形状：internal/risk/blacklist.go:51 var Logf func(...)（见 17-hop-g log 尺一）"
  echo "甲2 注入位候选（host 侧现成设施，只点名）："
  echo "  cmd/wisp/logsink.go:153 slog.SetDefault / :122 logSink.logger()"
  echo "  先例＝别的包把自己的注入点接上去：cmd/wisp/models.go:288-290 sinkLogger := sink.logger()"
  echo "甲3 现有 A/B 日志行枚数（＝今天已经在 logf 里走的内容形状）："
  git grep -nE 'logf\("risk' $R -- internal/risk | grep -v '_test\.go:' | wc -l
  echo "甲4 host 是否已在别处把 risk 的串往出向面写（尺＝git grep -n risk.Logf 剔测试）＝0（见 17 尺三）"
  echo ""
  echo "== 乙（面板）：C17 现有方法面 + 确认卡今天的字段面 =="
  echo "乙1 C17 PanelBridge 接口方法名册（★不许新名，先点名）："
  git grep -nE 'type PanelBridge interface' $R -- internal
  b1=$(git grep -l -e 'type PanelBridge interface' $R -- internal | sed 's/^[^:]*://')
  echo "  接口文件＝$b1"
  git show "$R:$b1" | grep -nE '^\s+[A-Z][A-Za-z0-9]*\(' | head -40
  echo "乙2 确认卡今天携带的 reason 字段（panel 侧）："
  git grep -nE 'Reason|ReasonKnown' $R -- internal/panel/approval.go | head -20
  echo "乙3 panel 侧 Blacklist 三枚切片是否已进快照（尺＝git grep -n Blacklist 于 internal/panel）："
  pb=$(git grep -n 'Blacklist' $R -- internal/panel | grep -v '_test\.go:' | wc -l)
  echo "  internal/panel 非测试命中行数=$pb"
  git grep -n 'Blacklist' $R -- internal/panel | grep -v '_test\.go:' | head -10
  echo "乙4 前端是否已在渲染 approval 的 reason（★是否触碰 frontend/**）："
  fb=$(git grep -n -iE 'reason' $R -- frontend | wc -l)
  echo "  frontend 命中 reason 的行数（含 dist 等）=$fb"
  git grep -n -iE '"reason"|reason\b' $R -- frontend/src frontend/*.ts frontend/*.js 2>/dev/null | head -12
  echo "乙5 出向 JSON 键今天的形状（approval 卡片序列化处）："
  git grep -nE 'json:"reason|json:"reasonKnown' $R -- internal/panel | head -10
  echo ""
  echo "== 丙（模型正文）=="
  echo "丙1 今天到模型的拒绝串是哪一句＋在哪一跳拼的（见 18-hop-h：bridge.go:439 → reject → Text）"
  echo "丙2 该句里今天含不含绝对路径（R3 那句带 %s canonical）："
  git grep -n 'R3: 命中 A 档敏感路径' $R -- internal/risk/rules_gateway.go
  echo "丙3 真 reason（aRuleReason）里带的字节形状＝cand.raw / cand.root / anchorSpelling / anchor.root（%q 四枚）"
  echo "丙4 ★先例尺：模型上下文里的绝对路径／控制字符，本仓今天有哪几处治理："
  git grep -nE 'display_scrub|scrub|controlChar|\\\\r|\\\\n' $R -- internal/panel/display_scrub_216.go | head -10
  echo "  174/264/216 三族件点名（只列在树里真实存在的）："
  git ls-tree -r --name-only $R | grep -iE 'scrub|174|264|216' | grep -v '^\.scratch' | head -20
} > "$L/20-ac1-three-shapes.log"

wc -c "$L/20-ac1-three-shapes.log"
echo "cr=$(grep -c "$(printf '\r')" "$L/20-ac1-three-shapes.log" || true)"
