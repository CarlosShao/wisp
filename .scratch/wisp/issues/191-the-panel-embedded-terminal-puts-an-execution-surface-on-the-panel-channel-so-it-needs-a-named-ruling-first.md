# 191 — 面板"内置终端"＝**执行面进面板通道**，是这次补功能里最大的一枚雷（先要 owner 逐枚批准，再谈实现）

- Status: **ready-for-agent（AC#0 已由编排者按甲裁结·09-28 16:4x，见 Progress log 末段与台账 `A379`）**。⚠ **功能票**（原"批准项"那一半已解除：owner 给的是**做法授权**，不是又一次待批）。
- 来路：owner 09-28 15:2x 放权（`A373`：「缺什么就给我补」）＋截图里"终端／在会话工作区里跑命令（Ctrl `）"。普查出处＝`docs/evidence/s1/180-182-panel-fields-census-c1.md` §5／§7（S5 面板终端载体＝**没票认领**；雷区判定＝**最大雷：执行面进面板通道，撞 `PLAN.md:2027` ①③**）。
- 关联：票 163（常驻终端会话＝工具侧地基，**不同物**：那是模型用的，这是人用的）· 票 146／`Q-49`（面板侧批准那一族刚钉过五枚门）· R20／票 92（面板只许显示＋发起请求）

## 为什么它不能像其他几堆那样直接补
"终端"意味着**在面板通道里执行任意命令**。这仓已定案的三条边界都被它压在下面：
1. `PLAN.md:2027` ①③（面板侧来源不得成为执行／批准的入口）；
2. R20 口径"面板里的档位与工作区是**权限输入口**，只许显示＋发起请求"；
3. D4"用户侧只能调严、调松要显式勾选"。
⇒ 现成的合规形状**有**：终端不执行命令，它**把命令交给既有的工具腿（`shell.*`／票 163 的常驻会话）走同一套门控与审计**，面板只是显示输入口与输出流。这一支不需要放宽任何边界。另一支（面板直接起 pty 自己跑）＝放宽，**必须 owner 点名批准**。

## AC
- [x] **AC#0 ⚠ 待人一句话（未答之前本票不许派任何落地腿）**：面板终端走哪一支？**(甲)** 复用既有工具腿＋门控＋审计（不放宽任何边界，代价＝命令集受工具面白名单约束、启动慢一点）；**(乙)** 面板侧独立 pty 直接执行（体验最接近截图，代价＝在面板通道上开一条执行面，要动 `PLAN.md:2027` ①③ 与 R20 两条已定案口径）。**我的推荐＝甲**。**不答的代价＝本票一直挂着，那一栏永远是示意数据。** 口令：「面板终端走甲」／「面板终端走乙」。
- [x] **AC#1 只读设计核（甲支）**：现有工具面里能承载"跑一条命令并流式回显"的是哪一枚（票 163 的常驻会话在不在名册级？尺：`grep -rn "shell\." --include=*.go internal/tools/ | grep -v _test.go` 现跑），缺什么（流式输出？超时？取消？）。
- [ ] **AC#2 载体与快照**：输入口／输出流怎么进快照（票 145 的载体扩张），⚠ **输出流里的外泄面**（命令输出会带绝对路径、token、环境变量）——必须过 C25 盖戳与外泄检测那套管子（票 175／177／183 同一条管子，别绕）。
- [ ] **AC#3 反向判据（本票的牙）**：交付后必须有一枚常驻判据钉住"面板侧送进来的命令**走的是同一套门控**，没有任何一条路径能绕过 `risk` 分级直接执行"。⚠ 若判据在"甲支"上今天就能绿，它是**不许弄坏的守卫**，不许充当本票新增的牙——牙要长在"绕过门控那一形会红"上（造变异现量）。
- [x] **AC#4 取消与超时语义**：长命令跑到一半用户关掉面板／按取消，今天这套语义是什么（票 176 的取消语义同族），逐形给读数。
- [x] **AC#5 契约轴**：`PLAN.md:2027` 那三条**不许由 agent 单方面改写**；D34 表、`allowlist.txt`、`thresholds.go`／golden、审批超时常量一字节不动；C17 若要加方法＝先落 `A##`。
- [x] **AC#6 门禁**：`./internal/tools/ ./internal/panel/ ./internal/agent/`＋`d22scan`（基线 433）＋`gate-clauses` 比名册（只 `G6neg`）；⚠ `internal/panel` 2 枚已知红照实记不修。

## 本票**不**解决
不做浏览器／产物预览（＝票 192，那里另有 `AGENTS §2` 未定义即停项）；不做树上写动作与 git 提交推送（票 189／190 的雷区，同样待批）。

## Progress log
- 09-28 15:2x 编排者立票：来路＝owner 放权（`A373`）＋`180-c1＋182-c1` 名册"没票认领"七处之一，且它是七处里唯一**压在已定案安全口径上**的一枚 ⇒ 立成"批准项＋实现项"两层，AC#0 未答不派。未派。
- 09-28 16:4x **owner 当场给决定（原话逐字）⇒ AC#0 由编排者裁结＝甲**（台账 `A379`）：「**终端你反正得做出那个效果，我管你怎么做的呢，就是点右侧边栏，点击终端，那就得出默认终端啊，比如 powershell 啊，有问题吗？**」
  - **读法（这次不再摆回去问他）**：他那句"我管你怎么做的"＝**做法授权交给我**，不是又一次待批；他要的是**效果**——右栏点"终端"就出一个能敲命令的默认终端（PowerShell 那一类）。⇒ 按票面三条候选里**唯一不新增执行口**的那一支走：**甲＝复用已过门的工具腿＋门控＋审计**。乙（面板侧独立 pty）不因"他给了效果要求"而被顺带批准——那要新增 `C17` 方法＝契约变更，需要他点名。
  - **效果写死进验收（防"做了但没生效"）**：① 点右栏那一格 ⇒ 真起一个终端会话（不是占位文案）；② 命令输出**必须过 C25 盖戳与外泄检测那套管子**（票 175／177／183 同一条管子，**别绕**）——⚠ 终端输出天生带绝对路径／token／环境变量，这是本票最大的外泄面；③ 关闭那一格要能取消在跑的会话（票 176 的取消语义）。
  - **排程**：本票落地腿排在 `33-r1`（入向那一跳）**之后**——没有"网页点一下→Go 收到"这一跳，"点击终端"这一下今天无处落；⚠ 且它要新增 `internal/panel/**` 的输入口，与 `33-r1` 同包 ⇒ **串行，不并行**。

## 编排者现量（2026-11-20，只读核，件 `docs/evidence/s1/191-readonly-survey-r1.md`）

**现状读数（已复核项标"复核"）**：能跑命令的工具面 **2 枚**（`shell.exec`／`shell.session`，`internal/tools/shell.go:235-243`，且只在 `[risk] shell_enabled=true` 注册）；**能流式回显的 0 枚**（输出整段缓冲，只有"已启动"一句 onUpdate）。输出进面板只有"工具结果"一条一次性管道。取消＝有（任务级），超时＝真 deadline，**关面板不取消**，面板按取消有但缺一跳。
**建议（编排者拍）**：AC#1 甲支＝先给 `shell.session` 做**增量交付**（读协程与 marker 抽取已有，缺"边读边发"），不新造第三枚工具；`shell.exec` 保持一次性。实现顺序：先 AC#1 增量交付＋AC#4 取消语义补齐，再谈 AC#2/AC#3。
**���解耦**：现量发现 shell 输出不在 C25 敏感源名册（`internal/risk/provenance.go:113-117`，与票 175 的洞同型）⇒ 单独立票、且属 SPEC-06 §5 契约文字需机主批（见台账 Q）。
## 编排者收件·AC#1＋AC#4（2026-11-20 21:1x，件 `docs/evidence/s1/191-readonly-survey-r1.md`）

- **AC#1 勾**：名册级候选 **2 枚**（`shell.exec`／`shell.session`，`internal/tools/shell.go:235-243`，只在 `[risk] shell_enabled=true` 注册）；**能流式回显 0 枚**（`shell_exec.go:112-116` 整段缓冲、`shell_session.go:296-352` 见 marker 才整段取、`shell_exec.go:121-123` 唯一 onUpdate 只发"已启动"、`shell_session.go:431` 零 onUpdate）。结论：甲支可行，只缺"边读边发"一跳。
- **AC#4 勾**（逐形读数）：任务取消＝有（`internal/tools/cancel.go:21-94`＋`bridge.go:720-730`＋`cmd/wisp/run.go:1441`）；超时＝真 deadline（`shell.go:145-152` 30m→`bridge.go:836-841` `context.WithTimeout`）；**关面板不取消**（`panel_host_windows.go:1287-1329/1345-1351/1368-1379` 三处不碰 task ctx，真取消在进程退出 `resident_windows.go:533-535`）；面板按取消＝有但缺一跳（`panel/bridge.go:266` `roster.cancel`→`composer_dispatch.go:493-497`→`panel_roster_cancel_wave2.go:119-149`）。
- 仍开：AC#2（实现时输出流必须过 C25 盖戳——现量发现 shell 不在名册，已挂 `Q-457` 待机主）、AC#3（牙要落在"绕过门控会红"）、AC#5（契约不动）、AC#6（门禁）。
## 进度·第一波交付（2026-11-20 21:3x，`6512d5a73`，件 `docs/evidence/s1/191-stream-tools-r1.md`）

工具侧增量交付已落：`SessionOptions{OnChunk func(chunk string)}`，`shell_session.run()` 在写命令行前装 sink、`defer` 清；`guardedText` 加只读第二面 `setSink`（`takeAll`/缓冲语义一字未变，marker 剥离仍在最终 Result）；`Execute` 把 onUpdate 接缝接进去，**签名零改动 ⇒ 回调 nil 时与今天逐字节一致**；`shell_exec.go` 一字未动；增量帽＝`sessionPendingCap` 64MiB 同源同值；零新增裸 `go func(`。
**我复跑**：`go test ./internal/tools -count=1` **ok 27.408s**、`d22scan` clean。腿报四数 `RUN=601 PASS=600 FAIL=0 SKIP=1`（SKIP＝`Test23LiveClipboard` 真机存量）。
**牙四处**：①nil 回调逐字节一致 PASS ②分片按序拼接 == 整段 PASS ③"读协程只发一次"⇒红（已证）④"去掉上限"⇒红（已证）。
**腿自报的两处存量未修（只报不动）**：①重置时旧读协程对重建 `pending` 的窄竞态（163 起存量）；②桥侧 onUpdate 闭包恒非 nil、宿主未接线时增量被吞（零可见变化）。
**未完**：AC#2（面板载体＋C25 盖戳，后者卡 `Q-457`）、AC#3（"绕过门控会红"的牙）、AC#5（契约轴）、AC#6（门禁三包）。
## 进度·第二波 a 交付（2026-11-20 22:0x，`8513bc2e3`，件 `docs/evidence/s1/191-panel-carrier-r1.md`）

面板侧载体已落：`Snapshot.Terminal`（reader-gated＋omitempty）＝`input`（只读展示，尾窗 4096 runes）＋`output`（滚动尾窗 8192 runes，`truncated` 终身伴随、`droppedRunes` 诚实记账；runes 非 bytes，防劈中文）。**到达事实**：没到过就缺席，不发空串冒充"没有输出"；有读者空态＝`"terminal":{}`。
**只读不执行**：载体文件 import 名册钉死＝恰好 `{sync}`；`TerminalFeed` 方法集钉死＝`{Input, OutputChunk, Section}`，只收文本；生产面进程构造点普查钉死＝恰 2 枚（`git_readonly_exec.go`／`tree_reader_190.go`）。**我复跑 `git grep -c exec.Command -- internal/panel`** 生产面只 `git_readonly_exec.go:2`／`tree_reader_190.go:1`（其余全在 `_test.go`）。
**台账**：`composerTSGapLedger` 登记 `"PanelSnapshot": {"terminal"}`（judged gap；页面半归前端腿，本波 frontend 零字节）。
**我复跑**：`go test ./internal/panel -count=1` **ok 16.249s**、`gofmt -l internal/panel` 空、`d22scan` clean。腿报 panel 四数 `RUN=306 PASS=306 FAIL=0 SKIP=0`（并如实报：票面写的"2 枚已知红"今日不存在，改前基线 `3766819a8` 已全绿 —— 那是 09-28 时点的过期记载）。
**腿自报的额外动作（具名）**：顺手修平 `config_handlers.go`（8 行对齐）与 `result_golden_200_36_test.go`（62 行注释缩进，内容零改动），二者在 HEAD 即 gofmt 不合规；**本机 gofumpt 未装（该半边门没跑成，只报不装）**。
**未完**：AC#2 的 C25 盖戳半边卡 `Q-457`（待机主）；AC#3（绕过门控会红的牙）；AC#5／AC#6。**第二波 b（`cmd/wisp` 接真源）在 Q-457 拍前不派**。
## 编排者收件·AC#5＋AC#6（2026-11-20 22:2x，件 `docs/evidence/s1/191-terminal-ui-receipt.md`／`191-panel-carrier-r1.md`／`191-stream-tools-r1.md`）

- **AC#5 勾（契约零字节）**：本会话全部提交 `git diff --name-only 3bfc17229..HEAD | grep -Ei "PLAN.md|thresholds|golden|allowlist|docs/specs|D43"` 只命中 `internal/panel/result_golden_200_36_test.go`（名字含 golden）；**逐行核过**：该文件的改动**全部是注释重新缩进**（`git diff -w --ignore-blank-lines` 输出为空 ⇒ 零语义），`thresholds.go`／`allowlist.txt`／`PLAN.md`／`docs/specs/**`／D43 **一字节未动**；C17 未加方法。
- **AC#6 勾（三包＋扫描）**：`go test ./internal/tools -count=1` **ok 27.408s**；`./internal/panel` **ok 18.138s**；`./internal/agent` **ok 6.358s**（158 条 PASS/FAIL/SKIP 行）；`sh scripts/d22scan.sh` **clean 0 findings**（ban #7 internal/tools=52、ban #8 internal=1073 等基线未降）。⚠ `scripts/gate-clauses.sh` **不存在**（`ls scripts/ | grep -i gate` 空），照实记。
- **仍开**：**AC#2**（输出流外泄面必须过 C25 盖戳——现量 shell 不在 `sensitiveSourceTools`，已挂 `Q-457` 待机主；另泵源接线＝wave 2b）／**AC#3**（"绕过门控会红"的牙，依赖 wave 2b 的真源）。
- 界面半已交：`panel.ts` 声明 `terminal`＋账本收口＋`terminal-region.tsx`＋`render:terminal191` 三态仪器；我复跑 typecheck rc=0、仪器 rc=0、panel ok 18.138s、d22scan clean。
## 进度·补牙：上限数值钉（2026-11-20 23:3x，`d738d4e52`，件 `docs/evidence/s1/191-cap-pin-r1.md`）

对抗验收抓到"上限存在"被钉、"上限数值"没被钉（64→128MiB 全绿）。本笔补两枚**零派生**的钉：`Test191SessionPendingCapNumberPinnedTo64MiB`（断言常量恰 `64<<20`）与 `Test191FloodRefusedPastLiteral68157440Bytes`（剂量与期望全用字面数字，断言拒因含 `67108864`）。**我复跑**：两枚 PASS、`./internal/tools` 整包 **ok 27.000s**、`d22scan` clean。腿自行验：只改常量到 `128<<20` ⇒ **两枚都红**（逐字红句在证据件 §2），还原后 SHA256 逐字节同。
腿另顺手格式化两枚**既有**漂移件（`gate.go`／`ticket49_limits_test.go`，`git status` 证非本腿引入、纯格式零语义），已在 commit 具名。