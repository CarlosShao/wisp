#!/bin/sh
# 337-a1：尺终版对账（行尾 inclusive ＋ 类型级正尺 ＋ 越面尺正确分母）——本腿三枚尺面缺陷的复现尺
cd "/d/work/workspace/projects plans/Wisp" || exit 1
L=.scratch/wisp/probes/337/a1/logs
R=c11163c8
{
  echo "### 尺终版对账 rev=$R（本腿自报三枚尺面缺陷并给终值；早期那几发原样留在 logs/12・13・95 作对照件，⛔ 删）"
  echo ""
  echo "== 缺陷①：收紧形 [^A-Za-z] 会漏行尾那一枚 =="
  echo "尺＝git grep -nE 'dec\\.Reason\$' \$R -- internal/tools/bridge.go"
  git grep -nE 'dec\.Reason$' $R -- internal/tools/bridge.go
  echo ""
  echo "== 缺陷②：(note|n|Blacklist)\\.Reason 里的 n 会子串命中 decision.Reason（decision 以 n 结尾）"
  echo "  ⇒ 与 turn.Reasoning 同族坑。类型级正尺＝显式 Blacklist\\.Reason／note\\.Reason／唯一 BlacklistNote 变量 n 的落点："
  echo "-- 显式 Blacklist.Reason（剔 .scratch，含测试）命中＝"
  git grep -nE 'Blacklist\.Reason' $R -- '*.go' | grep -v '^[^:]*:\.scratch/'
  echo "Blacklist_Reason_hits="$(git grep -nE 'Blacklist\.Reason' $R -- '*.go' | grep -v '^[^:]*:\.scratch/' | wc -l)
  echo "-- note.Reason（剔 .scratch）命中＝"
  git grep -nE 'note\.Reason' $R -- '*.go' | grep -v '^[^:]*:\.scratch/'
  echo "note_Reason_hits="$(git grep -nE 'note\.Reason' $R -- '*.go' | grep -v '^[^:]*:\.scratch/' | wc -l)
  echo "-- 唯一 BlacklistNote 变量 n 的声明（mode.go:86-87）与它唯一一次 .Reason 落点（:111＝写侧）--"
  git show "$R:internal/tools/mode.go" | sed -n '86,87p'
  git show "$R:internal/tools/mode.go" | sed -n '109,112p'
  echo "-- promptFor 那枚 Reason: d.Reason 的 d 类型（tools.Decision，非 Note）＝下面签名逐字 ⇒ 属正控那一族 --"
  git grep -nE 'func \(g \*Gate\) promptFor' $R -- internal/agent/approval/gate.go
  echo ""
  echo "== 终版尺（行尾 inclusive：(\$|[^A-Za-z])；剔 .scratch、剔 _test）=="
  echo "-- 被测尺逐枚（含假命中，供判堆）--"
  git grep -nE '(note|n|Blacklist)\.Reason($|[^A-Za-z])' $R -- '*.go' | grep -v '^[^:]*:\.scratch/' | grep -v '_test\.go:'
  echo "被测尺命中枚数（收紧前）＝$(git grep -nE '(note|n|Blacklist)\.Reason($|[^A-Za-z])' $R -- '*.go' | grep -v '^[^:]*:\.scratch/' | grep -v '_test\.go:' | wc -l)"
  echo "  ⇒ 逐枚判堆：mode.go:111＝写侧（唯一真指名 Note.Reason）；panel_pump.go:75／panel/approval.go:47,83,84／assessor.go:13＝子串假命中或注释，读的是 Decision.Reason"
  echo "-- 正控尺逐枚 --"
  git grep -nE '(dec|decision|it\.Decision)\.Reason($|[^A-Za-z])' $R -- '*.go' | grep -v '^[^:]*:\.scratch/' | grep -v '_test\.go:'
  echo "正控命中枚数（终值）＝$(git grep -nE '(dec|decision|it\.Decision)\.Reason($|[^A-Za-z])' $R -- '*.go' | grep -v '^[^:]*:\.scratch/' | grep -v '_test\.go:' | wc -l)"
  echo "与 logs/12 尺三那发（两条 -e＝6 枚）的差异具名：终版多 panel/approval.go:83/:84 这两枚 decision.Reason（前发未带该别名）；前发里 bridge.go:439 本就在（该行后还有引号，不属行尾）。⇒ 终值 8 枚；两发都 ≥1 ⇒ 正控成立，被测尺那些 0 不受影响。"
  echo ""
  echo "== 缺陷③：越面尺用错分母（logs/95 那枚 3）=="
  echo "错形＝git diff --name-only $R HEAD（＝两号之间全部提交的并集，含别腿字节）＝$(git diff --name-only $R HEAD | grep -cE '^(internal|cmd|frontend|design|tools|scripts)/' || true)"
  echo "对形＝只对本腿逐笔 name-only 后过滤 ^(internal|cmd|frontend|design|tools|scripts)/ ＝"
  for c in 5db0b538 1301903a 6559fb4f a24a8684 afb840d1; do git show --name-only --format='' $c; done | grep -cE '^(internal|cmd|frontend|design|tools|scripts)/' || true
  echo "⇒ 真读数＝0 枚（越面 0）。那 3 枚属别腿，尺＝git log --format='%h %s' $R..HEAD -- internal cmd frontend 主题逐枚可见（166-r1 两笔＋322-r2 一笔）"
  echo ""
  echo "== 终值一览 =="
  echo "被测：BlacklistNote.Reason 判定型／打印型／出向读者＝0 枚；搬运型＝1 枚（pending_read.go:67 out := d）；写侧＝1 枚（mode.go:111）"
  echo "正控：Decision.Reason 非测试读形＝8 枚"
  echo "其余不变：risk.Gate 生产调用点 1／blacklistedSensitive 调用点 0（正控 1）／risk.Logf 生产注入 0（正控 24）／risk.Resolve 非测试 1（含测试 7）"
} > "$L/96-ruler-final-reconciliation.log"
wc -c "$L/96-ruler-final-reconciliation.log"
echo "cr=$(grep -c "$(printf '\r')" "$L/96-ruler-final-reconciliation.log" || true)"
