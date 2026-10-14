# 票 133 五格复核 —— 分发跳新尺的非实现者验收（133-v 第一版）

- **身份**：非实现者。实现者＝`133-r1`（仪器 `da389ca5`；交付 `1740fd48`／`46dee9ef`）。
  本表未翻任何框、未动票面 `Status:`、零产码、只 commit 未 push。
- **锚**：起手 `51e2f5ea`（20:51:42）→ 现读 **`b75089cc`**（20:52:57，编排者"133-r1 收件"提交落地的漂移，逐报）。
  仪器文件自 `da389ca5` 起零改动（`git diff --name-only da389ca5 HEAD -- cmd/wisp/` 为空）。
- **方法**：全部变异落**仓外快照** `D:/tmp/wisp133v/base`（`git archive b75089cc`），
  `GOTMPDIR='C:\tmp\wisp133vgotmp'`；⛔ 仓内 `cmd/wisp/**` 零写。仓内只写本表＋`probes/133/v/`＋票面一条 bullet。
  逐发读数原文在 `.scratch/wisp/probes/133/v/`（≤7 件）与仓外 `D:/tmp/wisp133v/logs/`（只建不删）。
- **窗口自证**：W0 20:57:11／W2 21:25:58 无 `go|test|compile` 进程（`cargo.exe` 系子串误配、`wisp.exe`＝机主常驻球，逐条归因）；
  字面两趟 21:11:06-21:21:19 与 21:33:20-21:43:33 窗口内无第二枚 cmd/wisp 重包进程。
  另：21:2x 起 `cmd/wisp/` 出现本波另一腿的在飞件 `?? task_output_pointer_notice_174r8_test.go`（零触碰，本腿全在快照跑）。

## 结论（逐格判语）

| 格 | 判语 | 一句话凭据 |
|---|---|---|
| AC#1 归属硬线 | **PASS** | 五发（本腿自造）门关下逐发 rc=1、红名逐字 `TestAC1AC2DispatchHopGate133`；门开两份给出；N-3"不加尺"字面复核成立 |
| AC#3 自证腿 | **PASS** | 独立弄坏 `isRunnableCase133`（去掉 Recv 半边）⇒ 自证腿 rc=1 红、门 rc=0 PASS；还原双绿 |
| AC#4 静默缩短 | **PASS** | 注册表↔census 与 usage 双向差、Y2 换靶 pair 红（131 门同树 PASS）、下限 11<12 与全瞎两条红 |
| AC#5 门禁/314 | **PASS（附如实报）** | 字面 `-count=2 -v` 撞 `panic: test timed out after 10m0s`；第 2 趟 603.103s 独立复现；四数按其留件抽验＋startable 300→302 |
| AC#6 CI 落点 | **PASS（口径收窄一处）** | run/job/step/四数现读逐字对上；该 run sha 早于仪器 ⇒ 两枚新用例尚无远程 run 跑过（如实登记） |
| 面/不丢牙 | **PASS** | `da389ca5` 只两枚 cmd/wisp 测试件、-12 行全是自尺旧文案；两枚交付提交越面零；gofumpt 只点既有债；d22scan 不降 |

## §1 AC#1（核心）：五发×两份读数＋N-3 字面复核

做法：五发**本腿自造**（新命令名/新文件名/新常量名，形状与 r1 探针不重名），落仓外快照；
门关＝`-run '^TestAC1AC2DispatchHopGate133$'`（131 的门不注册），门开＝两门同跑。逐发读数原文见
`probes/133/v/01-ac1-five-shots.txt`。

| 发 | 本腿形状（落地行） | 门关 | 门开 |
|---|---|---|---|
| N-3 | `main.go` 摘掉 `os.Exit(cmdModels(...))` 一行 | rc=1，红名本尺；红点名 `cmdModels (models.go:104)`＋`modelStore.handOffModel (models.go:302)`＋`modelsEnsure (models.go:262)` 三枚孤儿＋nail-hop 一行 | rc=1 双红 |
| X4 | 早退 `if args[0]=="--hopv133"`＋新件 `hopv133v.go` | rc=1，红名本尺；`leg "--hopv133" (main.go:69) reaches installLogSink on this path: dispatch -> installLogSink`＋usage 方向 | rc=1 双红 |
| X8 | `case doctorCmdV133:`（新件 const 1 枚） | rc=1，红名本尺；非字面量标签（"resolves to the literal \"doctor\""）＋usage 方向；账本仍出全 12 行、只该行变合成腿（不是整包坏） | rc=1 双红 |
| X12 | `var hopAliasV133 = installLogSink`＋`case "hopalias133"` | rc=1，红名本尺；`dispatch -> hopAliasV133 -> installLogSink`（别名边被走图接住） | rc=1 双红 |
| X14 第一拍 | 结构体字段初值 `vHolder133.openFn=installLogSink`＋`case "vhold133"` | rc=1，红名本尺；`dispatch -> vHolder133.openFn -> installLogSink` | rc=1；**131 门 PASS**（种子无 slog＝票面原形） |
| X14 第二拍 | 字段初值与那次调用删掉、腿仍在 | rc=1，红名本尺；`covered by nothing`（无 install 义务后的第二条腿也红） | rc=1；**131 门 PASS** |

**N-3"全绿不可满足"逐字复核**：同一发下把本尺用 `-run` 排除（＝"不加尺"）⇒ 131 门单跑 rc=1
（红名只有 `TestAC4EveryLegIsNailedOrRuled`，红文逐字见 01 文件）⇒ 字面第①步"先证今天全绿"**不可能成立**，
且"红名只有 131 的门"这一元事实在今日树上仍为真。本尺装上后同一发必须也红 ⇒ 归属硬线成立。
五发逐发还原（main.go sha256 回 `24efbffb…`；新文件删净；终态四用例 rc=0 全 PASS）。

## §2 AC#3：自证腿不哑（本腿独立变异）

本腿形状＝去掉 `isRunnableCase133` 的 `fd.Recv != nil` 半边（≠r1 的 `return true` 顶置）：
夹具里的同名方法 `TestPhantomMethod133` 被洗进覆盖桶 ⇒
`TestAC3SelfProof133TheDispatchReadersCanAnswerNo` rc=1（失败断言逐字在 `dispatch_hop_selfproof_133_test.go:315`：
"the coverage reader did not report gamma's only would-be driver is a method…"），
**同一棵树上本尺门 rc=0 PASS**（真目录无此形状，门不连坐）⇒ 没有这枚自证腿，该形在这一棵树上就是哑的。
还原后两枚同跑 rc=0。读数原文 `probes/133/v/02-ac3-ac4-mutations.txt`。

## §3 AC#4：三读数独立重现

1. **双向差**：预读（零变异两门）双 PASS；变异 A＝注册表一条 claim 改名（`models`→`modelsv`）⇒ rc=1
   `claims leg "modelsv" … and func main dispatches no such leg … the Y1 shape ticket 129's readings were filed for`；
   另一对（census↔usage）两向各一发：删 `case "doctor"`⇒ usage 方向红＋下限红；usage 加 `phantomcmd` 行 ⇒
   `the usage block documents command "phantomcmd", which func main does not dispatch`。
2. **Y2 实体**：`run`↔`models` 两枚 claim 的 test 名对调（2 行）——**改前两门全绿**（预读 rc=0 双 PASS）、
   **改后本尺 rc=1 且 pair 红恰 2 枚逐字点名**（"pairs the case … with the entry …, and that case never calls …"）；
   同树 **131 门单跑 PASS** ⇒ pair 红不蹭 131。
3. **下限 N=12**：y1drop（删整枚 doctor 子句）⇒ `enumerated 11 legs (12 is the floor written down next to minLegs133)`
   ＋usage 方向；y2blind（`os.Args[1:]`→`os.Args`）⇒ `func main binds no local to os.Args`＋0 腿撞下限
   ⇒ **不返回空清单当"没有腿要测"**。**131 门内下限未动**：`leg_sink_gate_131_test.go` 只有
   `len(legs)==0` 的 Fatal（:258-259）、无 `>=N` 数值下限，且该文件最后动于 `6f702aeb`、不在本票提交文件表里。

## §4 AC#5/314

- **字面命令** `go test -count=2 -v ./cmd/wisp/`（默认 10m `-timeout`）：本腿两趟都撞
  `panic: test timed out after 10m0s`（第 2 趟姿态修正后 603.103s、rc=1），逐字复现票面那一枚。
  第 1 趟红名单曾混入两种快照姿态债（无 `third_party` 的 `no native DLLs …` 一族、无 `.git` 的
  `git rev-parse --show-toplevel` 一枚）；修正后第 2 趟只剩 7 枚 r1 登记的宿主底噪＋那 1 枚 `.git` 姿态红。
  两趟 partial 数（403/272/21/5、311/228/8/5）随负载浮动，不作判语；与 r1 读数（318/235/7/3，603.108s）**同形**。
- **四数抽验（按 r1 留件）**：CI 形 step 第 2 跑 282/217/4/1（四枚本票用例逐枚 PASS，原档 18 可复读）；
  40m 跑完 824/575/23/6；step 第 1 跑 306.2s 整进程死亡（314 形）；低负载同树 306/226/7/3。
  **改前/后哪数动（本腿独立量）**：startable 分母 `go test -list` **300→302**（+2 恰为自证腿两枚；
  门与 131 门都在改前已有）⇒ 同环境整趟的 RUN/PASS 各 +2、FAIL/SKIP 零；票 123 批审批用例零触碰。

## §5 AC#6

run `38554642217`（sha `72db0f39…`，dev，failure）／job `test-windows` `115820641622`／step #7
**"cmd/wisp CLI tests (needs the sherpa DLLs staged above, ticket 111 AC#4)"**（ci.yml:533/:601 现读）；
step 自己的四数逐字 **RUN=410 PASS=297 FAIL=2 SKIP=1**；该 step 2 枚红都不是本票用例。
链＝step→`scripts/wisp-cli-tests.sh`(:223-224)→`portable-tests.sh --scope=cli`→`scope=(./cmd/wisp/)`（:341）。
**口径收窄（如实登记，不是它的破口）**：该 run sha 早于 `da389ca5` ⇒ 自证腿/逐枚可区分两枚
**没有远程 run 跑过**；"仪器自己进 CI"结构上成立（`_test.go` 在 cmd/wisp、step 跑整包），
但"上次真跑过的那两枚用例的 run"今天给不出。⛔ `.github/**` 零改动。

## §6 面核与门禁

`da389ca5` ＝恰好两枚 cmd/wisp 测试件（+424／+66-12），-12 行逐行看全是该尺自己的旧文案/常量/旧 claim 行，
**无断言放宽、无第二枚生产件**；`1740fd48`／`46dee9ef` ＝票面＋`probes/133/r1/**`，三枚合扫
`^\.github/|thresholds|golden` 零命中。gofmt/gofumpt 复量只点既有债（`models.go` CRLF 伪影、
`panel_transport_35r2_test.go`＝`18d43f8e` 自认）；仪器两枚文件 clean；`go vet ./cmd/wisp/` rc=0；
`sh scripts/d22scan.sh` rc=0、各 scope 相对同日控制**不降**（活树 vs base 快照的 +1/+1 与新落件同量级）。
细节原文见 `probes/133/v/04-face-and-undone.md`。

## §7 分歧/偏差登记（不许悄悄改数的地方）

1. **行号漂移（信息级）**：r1 探针日志引的行号（:203/:209/:219/:185）取自 a7b88997 上的未提交版本，
   与已提交文件（:216/:222/:241/:199/:202）差若干注释行；红名文本与判据逐字相同。本腿读数一律取提交后树。
2. **r1 票面"CI 逐字同形 gofumpt … 空"缺原档**：其 12 号原档只留 `gofumpt -l cmd/wisp/`；
   本腿复量 CI 形原命令**非空**（HEAD 快照点出 `.scratch/**` 已入库样品＋`panel_transport_35r2_test.go`），
   与仓已定的 `ci.yml:189-196`（161 AC#7：raw walk 结构性不可满足，真门槛＝tracked-set 步）一致。
   该句不当凭据用；本票面无增债。
3. **AC#6 口径收窄**（见 §5）：远程 run 早于仪器，两枚新用例无远程跑读。
4. **本腿自身两处环境修正**：快照无 `.git`／无 `third_party` 造成的一过性红名单，已逐枚归因并修正后重跑
   （第 1 趟数作废理由是这三族姿态债，明写在 03）。
5. **X14 种子差异**：票面原形（无 slog）在 131 门上 PASS——本腿两拍都复现这一"独占目击"；
   r1 主发种子带 slog 所以 131 也红，二者并列不冲突。

## §8 具名没做

- "两门都摘"的全包读数未重跑（只做了 N-3 的 131 单跑字面复核；全包是本机死区）。
- 方法值 `v.M(...)` 支未造、跨包 install 支不接（票面已划出地界）。
- linux 侧分母未取；4 枚 `_windows_test.go` 注册表钉的平台豁免未在此复核。
- 远程"两枚新用例真跑过"的 run 不存在（见 §7.3）；另一条同 step 的 `38551683523` 未复核。
- `go vet` 只跑 windows 形（rc=0）；`GOOS=linux` 形的 cgo 加载限制照 r1 登记，不重复。
- ⛔ 未动 `docs/reports/pending-and-issues.md`、`HANDOVER.md`、`.github/**`、131 的门、阈值/golden。

## §9 交件

- 本表；`probes/133/v/`：`00-anchor-env-windows.md`、`01-ac1-five-shots.txt`、
  `02-ac3-ac4-mutations.txt`、`03-ac5-ac6-gates.txt`、`04-face-and-undone.md`（5 枚，皆 ≤200 行）；
  票面 1 条 bullet（append-only）；仓外原档 `D:/tmp/wisp133v/logs/`（只建不删）。⛔ 未 push。

## §10 收笔后漂移补记（append-only，2026-10-14 21:5x）

- 本腿提交 `c1975a7c` 的父提交是 **`8616d71a`**（"编排者：171s 遗产代提＋174-r9 半程"），而非快照锚 `b75089cc`：
  `b75089cc → 8616d71a` 之间 `cmd/wisp/` 唯一变化＝新落件 `cmd/wisp/task_output_pointer_notice_174r8_test.go`
  （本波另一条腿，已提交；本腿零触碰）。
- 因此在新父上重跑了本票四枚用例（快照 `D:/tmp/wisp133v/head2`＝`git archive 8616d71a`）：
  `rc=0` 四枚全 PASS、账本仍 `12 legs, 4 claims` ⇒ **新落件未扰动本尺；本表全部读数对 `8616d71a` 依旧成立**
  （原档 `D:/tmp/wisp133v/logs/17-head2-green.log`）。
