#!/bin/sh
# 337-a1：AC#1 补充现量（三形的 frontend 触面 / 216 名册覆盖 / 264 模型侧治理）＋ 必答(iv) 最远字节
cd "/d/work/workspace/projects plans/Wisp" || exit 1
L=.scratch/wisp/probes/337/a1/logs
R=c11163c8
{
  echo "### AC#1 补充现量 rev=$R（★剔 .scratch/**，前端只读码不判定）"
  echo ""
  echo "== 乙（面板）：卡片今天是否已把 reason 画出来 =="
  git grep -nE 'view\.reason|reasonKnown' $R -- frontend/src/components/l2-approval-card.tsx
  echo "读数＝前端 l2-approval-card.tsx:207/:216 已在渲染 view.reason ⇒ 乙形的显示侧今天就在；缺的是那一串没被塞进去（面板的 Reason 来自 risk.Decision 融合句，见 approval.go:83）"
  echo ""
  echo "== 乙2：面板出向键今天有没有 Blacklist 形状（新 JSON 键＝禁区，先量现状）=="
  echo "internal/panel 非测试命中 Blacklist 的行数=$(git grep -n 'Blacklist' $R -- internal/panel | grep -v '_test\.go:' | grep -vc '^[^:]*:\.scratch/')"
  echo "⇒ 0 命中＝面板出向面今天完全不带 BlacklistNote；乙形若要送出去，只有两条既有路：(a) 复用现有 json:\"reason\"（approval.go:48，会改到面板那条拒绝文案本身），(b) 新键＝票面 AC#5 禁区"
  echo ""
  echo "== 乙3：216 显示底座的清洗名册覆盖谁（★具名读数）=="
  echo "scrubForDisplay 生产调用点＝internal/panel/composer.go:145 一处（scrubComposerForDisplay 内洗 Workspace.Canonical + Attachments[].Name）"
  echo "display_scrub_216.go:58 自陈（逐字）："
  git show "$R:internal/panel/display_scrub_216.go" | sed -n '58p'
  echo "⇒ 名册今天不含 ApprovalCardView.Reason ⇒ 乙形落地面若把带控制字符的真 reason 交给卡片，216 底座不洗它（读码推的，尺＝上面两处调用点枚数）"
  echo ""
  echo "== 丙（模型正文）：模型侧今天有没有绝对路径/控制字符治理 =="
  echo "sanitize/redact/scrub 在 internal/tools + internal/agent 的非测试命中（★只有 forensics 列一处）："
  git grep -nE 'sanitize|redact|scrub|maskPath' $R -- internal/tools internal/agent | grep -v '_test\.go:' | grep -v '^[^:]*:\.scratch/'
  echo "== 三族先例的真实落地件名（尺＝git -c core.quotePath=false ls-tree，★quotePath 必关：开着会把中文名的八进制转义里的 264 当命中＝假命中，本腿抓到并已具名）=="
  git -c core.quotePath=false ls-tree -r --name-only $R | grep -v '^\.scratch' | grep -E '174|216|264' | head -24
  echo ""
  echo "== 必答(iv)：这句实话今天最远流到哪一枚字节 =="
  echo "1) internal/risk/blacklist.go:83 logf(...) 把 rule 串当实参算出：sink 为 nil ⇒ 该分支丢弃（17-hop-g：生产注入点 0 枚）"
  echo "2) internal/tools/mode.go:111 n.Reason = d.Reason ＝写进 BlacklistNote.Reason（写，不是读）"
  echo "3) ★最远一枚＝internal/agent/approval/pending_read.go:67 out := d（cloneDecision 的整枚 struct 值拷贝，把 BlacklistNote 连 Reason 一起搬进 LiveApproval.Decision）——搬运型接触，不指名读该字段"
  git show "$R:internal/agent/approval/pending_read.go" | sed -n '66,69p' | cat -n
  echo "⇒ 之后无任何读者：判定型 0／打印型 0／面板出向 0（internal/panel Blacklist 命中 0）"
  echo "★对票面的一处细化（具名更正）：票面说「读者 0 枚」——严格读数＝判定/打印型读者 0 枚；搬运型有 1 枚（pending_read.go:67 的 struct copy，同函数 :73-75 那三枚切片是它的显式逐枚克隆）。这把尺两种形状不放行：正控 Decision.Reason 同形命中 $(git grep -nE '(dec|it\.Decision|decision)\.Reason[^A-Za-z]' $R -- '*.go' | grep -v '^[^:]*:\.scratch/' | grep -v '_test\.go:' | wc -l) 枚（见 12-landing 尺三）。"
} > "$L/21-ac1-panel-and-farthest-byte.log"
wc -c "$L/21-ac1-panel-and-farthest-byte.log"
echo "cr=$(grep -c "$(printf '\r')" "$L/21-ac1-panel-and-farthest-byte.log" || true)"
