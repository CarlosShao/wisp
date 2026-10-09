# 282-r1 起手锚（独占 Go 编译/测试面）

- 锚 sha（`git rev-parse HEAD`）：`6041d7f9369c7b0963e23d9587728ead63bc3871`
- 分支：`dev`
- 时刻：`2026-10-09 11:32 +0800`（`date '+%F %H:%M %z'` 现量）
- 本腿射程：票 282 的 `AC#1`（行为尺／种刀牙口）与 `AC#3`（深度判定前后同形状对拉）。
  两格今天因串行铁律按住，⛔ 不是任何人的欠账。

## `git log --oneline -5`（原文）

```
6041d7f9 A765/A766/A767：收 285-r2/247-a3/tick-drill-2/reconcile-9 → 票 285 结案、票 247 AC#0 裁完翻勾、新立票 286
aee14886 reconcile-9 落件：对 A730-A750 的承重数逐笔复跑（复跑 129／对上 114／对不上 8／复跑不了 7）
f8b738f3 281 drill-2: 票面 AC#1 注-承载表行号对账 + 6 枚 commit 逐枚验(存在/subject/票面 全对, ②零不符)
89080fcc docs(247-a3): census for ticket 247 question 3 (audio degradation: busy/denied/no-device)
9204f858 票 285 AC#4 腿 285-r2 读数落件：三发突变的红句逐字（MU-1 同发下票 283 的尺二 PASS＝作差是新增承重件／MU-2 只有甲的单元尺响／MU-3 作差轴静默）＋顶回三处（corr_percall_242_test.go:99-105 今天已在作差、尺①产码枚数现量 4、internal/agent 锚上就红 TestGoldenSingleToolCall）；⛔ 判据句与四格框未动、未翻勾
```

## 起手工作树快照（`git status --porcelain` 前 40 行原文）

```
 M .gitignore
 M .scratch/wisp/probes/152/my152.py
 M .scratch/wisp/probes/161/r6/logs/flip-1.txt
 M .scratch/wisp/probes/161/r6/logs/flip-2.txt
 M .scratch/wisp/probes/161/r6/logs/flip-3.txt
 M .scratch/wisp/probes/161/r6/logs/flip-4.txt
 M .scratch/wisp/probes/161/r6/logs/flip-5.txt
 M .scratch/wisp/probes/161/r6/logs/flip-6.txt
 M .scratch/wisp/probes/161/r6/logs/flip-baseline.txt
 M .scratch/wisp/probes/161/r6/logs/flip-restored.txt
 M .scratch/wisp/probes/242/r3/logs/probe-routed.txt
 M .scratch/wisp/probes/268/v1/evidence.md
 D design/assets/base.css
 D design/assets/icons.js
 D design/assets/theme.js
 D design/assets/tokens.css
 M design/doubao/README.md
 M design/doubao/demo/app.js
 M design/doubao/demo/index.html
 M design/doubao/demo/styles.css
 D design/index.html
 D design/screens/approval.html
 D design/screens/ball.html
 D design/screens/chat.html
 D design/screens/config.html
 D design/screens/cost.html
 D design/screens/firstrun.html
 D design/screens/palette.html
 D design/screens/privacy.html
 D design/screens/security.html
 D design/screens/states.html
 D design/screens/tasks.html
?? -
?? .scratch/.scratch/
?? .scratch/ci-logs/run-37158259050.err
?? .scratch/ci-logs/run-37158259050.log
?? .scratch/ci-logs/run-37166458550-failed.err
?? .scratch/commit-msg-167-ac1.txt
?? .scratch/commit-msg-167a2-s01.txt
?? .scratch/commit-msg-167a2-s23.txt
```

## 在飞的别人字节（本腿一律不碰、不复述内容）

- `M .gitignore`
- `design/**` 的删除与修改（⛔ 零读零写零转述）
- 大量 `probes/**` 日志
- 仓根一枚名叫 `-` 的未跟踪文件
- 其余 `.scratch/**` 在飞件

## 本腿承诺的不变式

- ⛔ 零产码字节改动：种刀只落盘、每发当回合还原，收工时起手 hash 必须逐字等值。
- 收工尺：`git status --porcelain -- internal cmd` 必须为空。
- ⛔ 不翻 AC 框、⛔ 不改判据句、⛔ 不 push、⛔ 无 `-overlay`、⛔ 无 `t.Skip` 当通过。
- 既有红 `internal/agent` `--- FAIL: TestGoldenSingleToolCall`（`internal/agent/loop_golden_test.go:70`）
  起手即在，归票 286，本腿不修不 Skip，只在三数里具名并报。
