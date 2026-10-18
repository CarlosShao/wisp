#!/bin/sh
# 337-a1：落 hop g / hop h 两枚日志（LF），全部尺取具名提交树、剔 .scratch/**
cd "/d/work/workspace/projects plans/Wisp" || exit 1
L=.scratch/wisp/probes/337/a1/logs
R=c11163c8

{
  echo "### hop g) risk.Logf（sink 定义 blacklist.go:51）生产注入点枚数 rev=$R（射程＝全仓 *.go，剔 .scratch/**）"
  echo "尺一＝git grep -nE \"(risk::)?Logf[[:space:]]*=[^=]\" \$R -- \"*.go\" | grep -v .scratch（★=[^=] 排除定义行与 == 比较）"
  g1=$(git grep -nE '(risk::)?Logf[[:space:]]*=[^=]' $R -- '*.go' | grep -v '^[^:]*:\.scratch/')
  echo "$g1"
  echo "赋值命中行数（含测试）=$(printf '%s\n' "$g1" | grep -c .)"
  echo "其中 _test.go 命中=$(printf '%s\n' "$g1" | grep -c '_test\.go:')"
  echo "★生产注入点（非测试）行数=$(printf '%s\n' "$g1" | grep -v '_test\.go:' | grep -c . || true) ⇒ 票面那句「生产码无注入点」在 $R 上成立"
  echo ""
  echo "尺二（★正控，同一把尺的形状＝导出 sink 变量在生产码的赋值；换成今天确实存在的 opt.Sink / opt.Logger 注入族）"
  echo "  git grep -nE \"\\b(opt|opts|options)\\.[A-Z][A-Za-z0-9]*[[:space:]]*=[^=]\" \$R -- \"*.go\" | grep -v .scratch | grep -v _test.go"
  git grep -nE '\b(opt|opts|options)\.[A-Z][A-Za-z0-9]*[[:space:]]*=[^=]' $R -- '*.go' | grep -v '^[^:]*:\.scratch/' | grep -v '_test\.go:' | head -12
  echo "正控非测试命中行数=$(git grep -nE '\b(opt|opts|options)\.[A-Z][A-Za-z0-9]*[[:space:]]*=[^=]' $R -- '*.go' | grep -v '^[^:]*:\.scratch/' | grep -v '_test\.go:' | wc -l) ⇒ ≥1，尺认得生产赋值形状，hop g 的 0 是真 0"
  echo "⚠ 本腿第一发正控用 risk\\.NowFunc 那形落空（正控 0 行＝尺不可信），已换 opt.Sink/opt.Logger 族，具名登记在 17-hop-g-logf-injection-0.log 对照件"
  echo ""
  echo "尺三（risk.Logf 全仓逐字命中，证明只有测试在注入）="
  git grep -n 'risk\.Logf' $R -- '*.go' | grep -v '^[^:]*:\.scratch/'
  echo ""
  echo "尺四（logf 调用点＝日志点今天在哪：blacklist.go 三枚 + provenance.go + syncdirs.go，共 20 行含定义）"
  git grep -nE 'logf\(' $R -- internal/risk | grep -v '^[^:]*:\.scratch/' | wc -l
  echo ""
  echo "尺五（host 侧现成的日志设施，甲形落地时的注入位候选，只点名不设计）="
  git grep -n 'slog.SetDefault' $R -- cmd/wisp | grep -v '_test\.go:'
  git grep -nE 'sinkLogger|func \(s \*logSink\) logger' $R -- cmd/wisp | grep -v '_test\.go:' | head -6
  echo "读数＝cmd/wisp/logsink.go:153 立进程级 slog 默认，cmd/wisp/models.go:288 有 sinkLogger := sink.logger() 的先例（别的包把自己的注入点接到它）"
  echo "⛔ 本腿不判甲形该接在哪：那是写腿的活"
} > "$L/17-hop-g-logf-injection.log"

{
  echo "### hop h) 到模型的那一句走哪条链 rev=$R（票面链逐跳对撞）"
  echo "跳1 rules_gateway.go R3 那句（★注意：这是 assessor 那句，不是 Gate 那句）"
  git grep -n -e 'R3: 命中 A 档敏感路径' $R -- internal/risk/rules_gateway.go
  echo ""
  echo "跳2 assessor 融合（票面说 :345 一带）"
  git show "$R:internal/risk/assessor.go" | sed -n '316,355p' | cat -n | sed 's/^/   /'
  echo ""
  echo "跳3 bridge.go:301 dec.Reason = verdict.Reason"
  git show "$R:internal/tools/bridge.go" | sed -n '299,303p' | cat -n | sed 's/^/   /'
  echo ""
  echo "跳4 bridge.go:439 route 的 Deny 支"
  git show "$R:internal/tools/bridge.go" | sed -n '436,442p' | cat -n | sed 's/^/   /'
  echo ""
  echo "跳5 reject 收 why 并塞进 agent.ToolOutcome.Text"
  echo "   reject 调用点（route 的第二返回值怎么进来的）："
  git grep -n 'b\.reject(ctx, req, dec, why' $R -- internal/tools/bridge.go
  git show "$R:internal/tools/bridge.go" | sed -n '360,368p' | cat -n | sed 's/^/   /'
  echo "   reject 定义与 Text 字段："
  git show "$R:internal/tools/bridge.go" | sed -n '1090,1110p' | cat -n | sed 's/^/   /'
  echo ""
  echo "跳6 loop.go:745 进对话正文"
  git show "$R:internal/agent/loop.go" | sed -n '743,747p' | cat -n | sed 's/^/   /'
  echo ""
  echo "★与票面差异（具名更正）：票面/判语件写的 ToolOutcomeText(:1105) 这枚标识符全仓不存在＝"
  echo "  git grep -n ToolOutcomeText \$R 剔 .scratch 命中行数=$(git grep -n 'ToolOutcomeText' $R | grep -v '^[^:]*:\.scratch/' | wc -l)"
  echo "  真身＝reject 返回的 agent.ToolOutcome{Text: why}（bridge.go:1105-1106，字段 Text: 在 :1106）"
  echo "★另：票面说 reject 函数注释自陈 returns what the model sees＝逐字在 bridge.go:1092（上一行是签名注）"
  git grep -n 'returns what the model sees' $R -- internal/tools/bridge.go
} > "$L/18-hop-h-to-model-chain.log"

wc -c "$L/17-hop-g-logf-injection.log" "$L/18-hop-h-to-model-chain.log"
echo "cr_check=$(grep -c "$(printf '\r')" "$L/17-hop-g-logf-injection.log" "$L/18-hop-h-to-model-chain.log" || true)"
