# 12-r1 预筛裁决表（票 12 三格 · 只读读腿，2026-10-31）

预筛腿 12-r1。零写入本表之外的工作树文件；测试全部在 `git archive HEAD`（`dcf6e910f`）导出的
仓外 staging（`D:\tmp\wisp\20261031\12r1\tree`）跑，未碰在飞 41a2-r1 的 `cmd/wisp` 脏件。

## 0. 开工三核
HEAD=`dcf6e910f` 在 dev（非 detached）；索引零已暂存条目；`cmd/wisp` 脏 = `resident_audio_windows.go`、
`resident_voice_windows.go`、`resident_voice_41a2r1_windows_test.go` —— 归 41a2-r1，未碰、未归因改动内容。

## 1. 格一 AC#2（Idle SLO gate PASS 入 docs/SLO.md 附录）—— 现在就能清，不缺测量
- 解除条件已成立：票面裁定「本框继续未勾，解除条件 = 票 66 的 AC#1–AC#4 全绿」（票 12 `:80-81`）。
  票 66 现为 **done（7/7 勾）**，1:1 裁决表 `docs/evidence/s1/66-adversarial-acceptance.md` 在案。
- 字面判据全部已 recorded，无需新跑：
  - numbers+machine+date：附录 C（机器 DESKTOP-LVS7839、2026-09-20、逐样本逐跑列表）。
  - tree-private ≤25MB：C.1 五连跑全 exit=0、主体私有 WS 4.21–4.39MB、九门 9/9。
  - settle ≤10s + FreeOSMemory>0：C.5（261ms 回 cap、free_os_memory_count=2、pass=true）；A.3 3/3 同量级。
  - CPU 行：树外口径当门（票 66 书面裁定），C.4 smoke 5/5 主体 CPU 0.0000%；树内数字降为记录项。
- `docs/SLO.md:77` 仍写「标点 未测 → 待票 27（S4）」——票 27 未收口，但该行与 AC#2（idle gate）无耦合，不挡本格。
- 动作（谁）：对抗验收腿按 SPEC-10 §8 复核附录 C 证据链并出 1:1 裁决表 → 编排者勾框。
  若对验收要求新鲜跑：桌面独占 `slo-check.ps1 -Subset smoke` 一轮（票 66 修好的仪器、零 keeper）。
- 本腿未勾框（只读）。**判定：能做，今天可清。**

## 2. 格二 AC#8（球状态人工目验，user signs off screenshots — D29）—— 要人=机主；建议排在 R15 答复后
- D29 原文（PLAN:986–1058）无「截图签收」程序句；它定的是「UI 质量升级为与 SLO 同级的验收维度」
  （PLAN:1002）＋悬浮球必须原生。签收程序与形态是票面自带的；**user = 机主**（D23 自用）。
- ⛔ 无替代可由 agent 发明：机审/录屏只能做**签收材料**（截图 + 逐态说明，可代办），不能替代 user 签收本身。
- 时机关键事实（票 68 实证）：默认构建现在画的仍是 12px 微点（`prototypeVisuals=false`）；
  owner 2026-09-20 签收的 44px 玻璃体只在 `balldebug`。票 68 AC#2（翻默认）被 **R15 三项挡住**
  （#3 Settling alpha 双乘 ≈0.33、#4 停靠露出 ≈82% vs 承诺 42%、#5 命中半径 21px < 可见外沿 26px），
  owner 未答前**不得翻默认值**。
- ⇒ 最省机主时间的路径：机主答 R15 三项 → 票 68 翻默认 + 复测（AC#2/AC#3）→ 默认构建 == 签收样子
  → 机主对 7 态截图**一次**签收（同时清票 12 AC#8 与票 68 的目验面）。目验操作照票 68 AC#3 命令集：
  显示器 100% 缩放（96 DPI 基准，A28/A29）、`balldebug -diff` 七态 + `-frozen` 对照。
- 谁做/何时：机主签收（agent 备料）；**等人，是三格唯一纯外部依赖**。

## 3. 格三 AC#9（A25：回复后 3s 回落实测记录）—— 写腿活，等 41a2-r1 收口
- 判据原文（票 12 `:149-158`）：①最后 token → 进 `Settling` ≤3000ms、单调钟（`observe.Timeout` 路）；
  ②≥5 样本全报、坏尾部不许平均抹掉；③机器可判形态优先，若 D43 回落行不暴露「最后 token 时间」，
  **加事件钩子/接缝 = 代码改动，写进本票**。
- grep 实证（本腿做的核）：`internal/statemachine/events.go` 全事件面无任何「最后 token 时间」载体
  （无 EvAfterReply / LastToken / tokenAt；S1 判据原文在 `PLAN.md:1430`，票面引 `:1424` 是行号漂移，内容同一条）；
  `cmd/wisp/run.go` 无 Settling/3s 计时（CLI text path 跑完即退）。3s `Settling` 窗的现成缝在
  resident/session 侧（`resident_voice_windows.go:265`、`resident_voice_engines_windows.go:15` 有 Settling window 概念）——
  写腿优先在缝上接单调钟两点测量，能不新发明事件就不发明。
- ⇒ **现在做不了**：首先是产码（判据③自己这么定），而 `cmd/wisp` 每波只许一枚写腿、41a2-r1 在飞。
  谁做/何时：41a2-r1 收口后的下一枚 cmd/wisp 写腿（接缝 + ≥5 样本实测，实测跑需错峰/桌面独占纪律）。

## 4. 票 16 解锁（16 blocked by 12＋15）
- 票 15 已收口（`15-speech-engines-cer-harness-done.md`）⇒ **唯一剩余阻塞 = 票 12 翻 done**。
- 三格清法与时序：
  1. AC#2 —— 对验收腿/编排者，**今天**（证据全在案，零新测量必要）。
  2. AC#9 —— 41a2-r1 收口后一枚 cmd/wisp 写腿（关键路径之一）。
  3. AC#8 —— 机主答 R15 → 票 68 翻默认 → 机主签收截图（关键路径之二，等人）。
- ⇒ **16 的关键路径 = max(AC#9 写腿, AC#8 等机主)**；AC#2 不在关键路径上。

## 5. `Status: review` 的含义与收口流程（具名出处）
- `issues/README.md:14`：`review` = Work done, awaiting adversarial acceptance（SPEC-10 §8）。
- 票面本身没有「review 后谁收口」字段；流程具名在 `issues/README.md` 规则 4/6：
  **非实现者**出与 AC 编号 1:1 的裁决表（每框 PASS/FAIL/PARTIAL，落 `docs/evidence/s1/`）
  → 全框勾 → `Status: done` → 改名 `-done` → 更新索引 → commit+push（推送归编排者）。
  裁决者≠实现者（D22 双角色、SPEC-12 §4.3）。
- 票 12 现状 8 框中 5 已勾、3 未勾（AC#2/AC#8/AC#9）⇒ review 的实质缺口就是这三格。

## 6. 零回归基线（staging，错峰）
- `git archive HEAD`（`dcf6e910f`）导出仓外 staging；十几个 `.scratch` 探针日志因文件名特殊字符
  tar 解不开（全部非代码件，零影响）。
- `go test -count=1 -run 'TestSLO' ./cmd/wisp/` → **20/20 PASS、0 FAIL、ok 0.773s、exit=0**
  （SLO144×8 + SLO147×3 + SLO149×3 + SLO152×1 + SLO156×5 + SLO159×1）。
- 跑前确认无并发 go 进程；`third_party/sherpa-onnx` 未被 git 跟踪，PATH 指向真树该目录（只读引用）。

## 7. 结论一句话
票 16 只等票 12；票 12 三格里 **AC#2 今天可清（对验收+编排者勾框）、AC#9 等 41a2-r1 后一枚写腿、
AC#8 等机主（先答 R15 再对默认构建一次性签收）** —— 编排者今天该做的动作只有两件：
派 AC#2 的对抗验收、催机主答 R15 三项。
