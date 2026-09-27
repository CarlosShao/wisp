# 163-a1 只读设计核裁定表 — 常驻 shell 会话的落点与污染源撞格

- 时刻：2026-09-27 21:16–21:2x +08｜锚点＝step-0 现量 `51c32ef9fa4597ccde1a41b857c5df058dcf3fb0`｜分支 `dev`
- 派单＝`.scratch/wisp/dispatches/2026-09-27-211x-readonly-163-a1-persistent-shell-landing-and-taint-collision.md`
- 本程零产码、零判据改动、零勾选；`internal/tools/task.go` 只读过、未动过。

## 1. step-0 五件（全部现跑）

| # | 读数 | 命令 |
|---|---|---|
| 日期 | `Sun Sep 27 21:16:34 CST 2026` | `date` |
| 分支 | `dev`（符合硬门） | `git rev-parse --abbrev-ref HEAD` |
| HEAD | `51c32ef9` | `git rev-parse HEAD` |
| 脏栏 | **空**——起手与两枚台件跑完后重量均空 ⇒ **「此刻 174-r1 正在写 internal/tools/task.go」这一前提在本程现场未成立**，原样报回；本程仍全程一字未进 `internal/**` | `git status --porcelain -- internal/tools/ cmd/wisp/` |
| 基线四数 | 包级 `ok github.com/CarlosShao/wisp/internal/tools 14.139s`；`-v` 下 **PASS=116 FAIL=0 SKIP=0**（顶层+子测试合计） | `go test -count=1 ./internal/tools/`；`go test -count=1 -v ./internal/tools/ \| awk …` |

## 2. 票面 AC#1…AC#5 连行号逐字（`.scratch/wisp/issues/163-…-one-shot.md`）

- `:3` Status 原句：「**立而不派**（要扩 `shell.exec` 的语义＝动 `D34`/`D46` 那条线，**人工批准**；且它和票 161 抢同一批审批判据）」
- `:4` Status 更正（18:4x）：「上面那句……**这道门已经开过了**——owner 09-27 09:30 那句『那六张票全做……』（台账 `A322`）已把新增工具名落进冻结文本：`PLAN.md:2568`（`shell.session`，L2）＋`SPEC-07:75` 镜像（提交 `ced72f8`……）⇒ **本票现在是 ready-for-implementation**；唯一还没做的是“契约文本已改、代码里零实现”」
- `:26` **AC#1** 先把“第二条命令不知道第一条的状态”量成读数（未修码上一发先 `cd` 再 `pwd` 的两步调用）。量不到＝前提不成立，停手上报。
- `:27` **AC#2** 认退出码那一发不许靠时间：埋标记取回退出码（照 dsh 那形），并答“摘掉标记这一味，是否存在一发超时从此判错”。
- `:28` **AC#3** 会话寿命绑任务：任务被取消／结束时会话必须跟着收；拿一发“任务取消后会话还在”的变异证明这条有牙。
- `:29` **AC#4** 审批形状不许松：常驻会话不许变成“授权一次就一直放行”——批准只批这一次（三家＋`Q-49` 同向）。不许动 `internal/agent/approval/**`。
- `:30` **AC#5** 契约轴：只许动 `PLAN.md` 的 `shell.exec` 那一行（需 owner 单独批准）＋ `internal/tools/**`／`internal/proc/**` 实现与测试＋证据件；`frontend/**`、`design/**`、`thresholds.go`、golden、`allowlist.txt`、`tools/d22scan/**` 零字节。

冻结文本本程亲验（〔现跑〕`grep -n shell docs/PLAN.md`）：`PLAN.md:2567` `shell.exec`（L2、默认禁用、argv 强制）与 `:2568` `shell.session`（L2；「一次性外观、常驻内核；会话寿命绑任务作用域；超时＝真 deadline（禁墙钟差）；每发命令仍按 R6 单独判；不引 PTY」）逐字在册。

## 3. 现状三发读数（台件 `.scratch/wisp/probes/163/a1/`，日志落台件自身目录 `logs/`）

### 3.1 今天“第二条命令不知道第一条”的形状（AC#1，真 cmd.exe，非镜像非假进程）

台件 `cd-pwd.sh`（`logs/cd-pwd-2nd.log`）与 `main.go`（`logs/go-probe.log`）逐字：

```
## 发1(进程A) cmd.exe /c 'cd /d C:\Windows & cd' => "C:\\Windows" rc=0
## 发2(进程B,全新) cmd.exe /c 'cd' => "D:\\work\\workspace\\projects plans\\Wisp" rc=0
## 读数: 发2 ≠发1的目录 ⇒ 两次生发之间没有任何存活载体
```

⇒ **AC#1 量到了**：两发独立 spawn 之间 cwd 丢失，读数非空。（首版 `cd-pwd.sh` 被 MSYS 把 `/c` 改写成横幅坏读数，旧日志 `logs/cd-pwd.log` 按规矩留存未删。）

⚠ **形状具名**：`shell.exec` 与 `shell.session` **今天都不在注册表里**（〔现跑〕真 Registry，形状照 `cmd/wisp/run.go:333-366`）：

```
## 注册表实名册: fs.edit(L2) fs.list(L0) fs.move(L1) fs.read(L0) fs.trash(L1) fs.write(L1) task.output(L0)
## Lookup("shell.exec") => registered=false
## Lookup("shell.session") => registered=false
```

全仓非测试代码 `os/exec` 零命中（〔现跑〕`grep -rln "os/exec\|exec\.Cmd\|CommandContext" internal/ cmd/ --include=*.go | grep -v _test` ⇒ 只有 `internal/llm/adaptertest/mockllm.go` 测试替身）。⇒ 票 163 `:46` 那句「`shell.exec` 本体今天名存实无」**为真**，写手腿是**同时新造两枚**，不是“扩一枚既有工具”。

### 3.2 退出码今天从哪儿来（AC#2 的形状裁定）

〔现读+现跑〕生产码里**一枚退出码取回点都没有**：唯一带 spawn 的基元是 `internal/proc/jobscope_windows.go:110 StartInJob(cmd *exec.Cmd) (*os.Process, error)`，其生产调用点只有 SLO 台架 `cmd/wisp/slo_windows.go:500`；`cmd/balldebug/diff_windows.go:549-568` 读 `ProcessState.ExitCode()` 但那是调试件。**既不是“取回的字段”也不是“从输出文本里猜”——是两者皆无**（派单前提“今天从文本里猜”不成立，按未验证断言报回）。

真 shell 上两向读数（`logs/go-probe.log`）：

```
## exit 7: errors.As(*exec.ExitError) 命中, ExitCode()=7 —— Wait 状态【取回的字段】
## 真值=0 的进程打印了 "exit 1"：若『从输出文本里猜』会把 1 读成 1；字段取回 rc=0 err=true
```

⇒ **摘掉取回机制后是否仍有路拿退出码：有——只剩“从输出文本里猜”（dsh 埋标记那形），而上面第二发就是它判错的活样本**（打印 `exit 1` 的 0 退码进程）。裁定形状：退出码主通道必须是 Wait 状态的字段；文本标记只许作会话形下的辅助，不许成唯一来源。取回发生在 `Wait()` 返回处，全程无墙钟差（d22scan ban #4 之形未出现，门禁 §6 现量绿）。

### 3.3 本程没测什么

- **`gate-clauses.sh` 未执行**——原因见 §6。引用的 G5/G6/G7 基线枚数一律是该脚本**文本声明**（09-27 别的程在锚上现量登记的 want_n），不是本程读数。
- 票 175-r1 的两枚红测复算：未跑，引用台账 `A346`／票 177 `:3` 的记载〔读码，非现跑〕。
- 完整 `bridge.Execute` 审批链路（需接 approval UI）未构造；R4→L2→拒的那条链以逐行读码＋`gate.go:144` 的字面为准。
- PTY/ConPTY、`SPEC-07:75` 镜像行号、PLAN 其余 46 枚决策：未量（非本单相干）。

## 4. 三格落点裁定

### 格 A：会话寿命绑任务（AC#3）——落 **桥的 `OpenTask↔CloseTask` 那一族**，收尾写进今天已存在的那枚任务结束边界本体

〔现读〕今天的任务结束边界＝`cmd/wisp/run.go:563` 的 `rt.bridge.CloseTask(taskID)`，它挂在 `agent.Options.AdmitTask` 的 revoke 上，而该 revoke 被文本环逐字 defer（〔现读〕`internal/agent/loop.go:360-366` `defer revokeAdmission()`；`internal/tools/bridge.go:651-658` 注释同述）。裁定：**会话名册的 Close 步骤落进 `CloseTask` 同路径同文件（`cmd/wisp/run.go`），不新造第二族收尾动词。**
- 为什么：这是全仓**今天唯一有 owner 的任务边界**；`plugin.DisposalScope.Defer` 那一族被排除——〔现读〕`internal/risk/provenance.go:62-65` 明写「no *plugin.DisposalScope reaches tools/agent/cmd in production today (160-c1 §2.2)」，把管子铺到 tools/agent 要动 `internal/agent/**`，超出本票射程。
- 落在哪条尺腿射程里：`gate-clauses.sh` **G6**（`pair 'OpenTask' 'CloseTask'`，差集判据＝「开过却没在**同一文件**关过」，脚本 `:189-204` diffsets 本体、`:450-453` 腿声明）。⚠「同文件成对」是**这把尺自己的形状**（派单该前提成立），且 G6 今天反向在册 1 枚＝`cmd/wisp/run.go`（文本声明 `want_n 1 ring rev`，`:446`「正向 0 枚／反向 1 枚」）⇒ 写手腿若把 Close 落到别的文件而该文件不出现 OpenTask，就是给反向基线**加枚**——G6 的响＝实测>基线。若选 DisposalScope 形，落 G7（今天正向 0 枚，任何新文件单独出现 `Defer` 即响）。
- 未修码上响不响：**不响**（shell.session 未注册，G6/G7 都在声明基线；本程未执行该尺，见 §6——此句为“按其声明形状推理”，标〔读码〕）。

### 格 B：审批形状不许松（AC#4）——**每发命令过一遍完整 C19 门**，会话存在本身不携带任何授权

冻结文本已把这句写死（〔现读〕`PLAN.md:2568`「每发命令仍按 R6 单独判」）；R6 判据本体 `internal/risk/rules_shell.go:14`（SPEC-06 §10 强制 argv）。与两枚冻结件共处：`Q-49` 口径（台账 `:1083`）禁的是**面板侧来源的 L2「允许」**——本裁定不依赖任何来源的“允许”，会话内每条命令的 L2 卡仍只能原生侧出；`SPEC-06-security-gatekeeping.md:37` 逐字「R4 …→ L2，且**不受 D45 会话授权覆盖**」⇒ 即使未来 D45 梯度授权批了某条命令，**R4 命中面永不进梯度**。本程**未提出任何豁免/白名单**（那是冻结面，且 `Q-61` 未决）。未修码上不响（无新码）。

### 格 C：来源戳——**163 的落地腿与票 175/176/177 撞同一批文件与同一枚语义，具名如下，不开第二支修法**

同文件：`internal/tools/bridge.go`（`mark`，`:551-564`）、`internal/risk/provenance.go`（`Mark`/`Inspect`/`matchText`）；同语义：盖戳×续读互斥（票 177 本体）。见 §5。

## 5. 污染源那一问的正面回答（含逐行引用）

**问：若常驻会话的输出被盖戳，后续一次把“会话输出里出现过的路径”当参数的 `fs.read` 会不会命中 R4？**

**答：会。** 链条逐枚（全部〔现读〕本程亲读）：

1. **今天盖不发生**：`bridge.go:552` `if b.prov == nil || dec.TaskID == "" || !risk.IsSensitiveSource(dec.Tool) { return }`——名闸按工具名放行，回退后现状与派单该前提**相符**；而 `shell.exec`/`shell.session` **不在名册**（`provenance.go:84-99` 八枚 `sensitiveSourceTools` 无 shell；〔现跑〕探针 `IsSensitiveSource(shell.session)=false`，`logs/go-probe.log`）。这正是票 175 命名的洞（`175-mark-gate-r1.md`／台账 `A345` 裁定「补戳是实现缺陷修复」）。
2. **一旦补上名闸（票 175 方向）即刻命中**：且 `provenance.go:489-491` 逐字「a Mark() for any other tool is still honored fail-closed」——直调 `Mark` 今天就按敏感记录。
3. **进索引**：`Mark` 把整段输出归一化后建 ≥8 字符连续碎片索引（`provenance.go:22-23` 合同形状、`:474-494` 本体）；绝对路径归一化后必超 8 字符。
4. **路径参数不豁免**：`Inspect` 的 `fs.read`（无冻结通道名）走 `provenance.go:647-674` else 支；`:662` `if !gateOpen && !isPathKey(k) { continue }`——**只有非路径键被跳过**，`:666` 注释逐字「The write target itself is still scanned.」；`isPathKey` 定义在 `:820-822`。
5. **匹配→R4**：`matchText` `:721-724` `m.idx.contains(norm)` 命中；桥把该 scope 的 Detector 绑进 C19（`bridge.go:712-723`，`TaintHit` 即 `Inspect(scope,"",params)`，`provenance.go:740-746`）；`ruleTaint` `internal/risk/rules_gateway.go:101-114` 出 `R4 + L2 + blockSessionAuth: true`；`assessor.go:18` 与 `SPEC-06-security-gatekeeping.md:37`：不受 D45 会话授权覆盖。
6. **无审批通道的后果**：`internal/tools/gate.go:144` 逐字「`L2 审批通道尚未接入（票 21），已拒绝执行`」——`wisp run` 下直接拒。

⇒ **与 177 的关系**：机制与 `task.output` 桩完全同格——桩文本逐字嵌宿主路径（`internal/tools/task.go:229`「…全文见 %s…」，本程亲读），177 票面 `:1` 标题与 `:11-21` 段、台账 `Q-61`（`pending-and-issues.md:1100`）载的红证两枚即此链。**163 写手腿不许在此另开第二支修法**；`Q-61` 未决前，175 的修法与 176 的起跑口按台账原话按住（同一句也覆盖 163 的盖戳腿）。⚠ 派单让本程引的 `docs/evidence/s1/177-c25-r4-path-exemption-c1.md` **在盘上不存在**（〔现跑〕`ls docs/evidence/s1/ | grep 177` 空）——与票 177 AC#1 自纠段一致，该“裁定表”从未交付；本表引用的是票 177 票面与台账本体。

## 6. 门禁两枚读数 + 为何没跑 gate-clauses.sh

- `sh scripts/d22scan.sh` ⇒ `d22scan rc=0`（ examined 230 production Go files；「clean - no D22 ban violations」）〔现跑〕
- `bash tools/d22scan/runtests.sh -C tools/d22scan ./...` ⇒ `top-level: PASS=34 FAIL=0 SKIP=0, === RUN=76, '[no tests to run]'=0`〔现跑〕
- `.scratch/wisp/probes/154/gate-clauses.sh`：**只读了文本、未执行**。理由＝此刻若有别程在 `internal/tools/**` 写码/改名册，pair 腿的“同一文件成对”差集与 want_n 基线枚数会拿进**不可归因**的读数（名册非空即 BAD，错算到本程头上）；且腿文本自述聚合退码只累计「声明与实测不符」的腿（`:519-521`），基线登记权在上一程。本表引用的枚数（G5=1、G6=1、G7=3）全部标为**脚本 `:385/:451/:489` 行的 want_n 文本声明**，非本程测量。

## 7. 被拒／没成功的调用

取数**后**两枚：① 首版 `cd-pwd.sh` 的 `cmd.exe /c` 被 MSYS 路径改写，读出交互横幅（坏读数，旧日志留存未删，修后 `cd-pwd-2nd.log` 为本表引用版）；② `sed docs/specs/SPEC-06-security-gates.md` 打错文件名（实际 `SPEC-06-security-gatekeeping.md`，改后即 §5:37 读数）。取数前零被拒。

## 8. 删除命令

**零。** 未执行任何 `rm/del/Remove-Item`；坏日志、旧 fixture 输出一律留存。

## 9. 伪授权两栏

1. **「`docs/evidence/s1/177-c25-r4-path-exemption-c1.md` 已裁出候选甲不可行」**——该文件不存在（尺：`ls docs/evidence/s1/ | grep 177` 空）；票 177 `:33` AC#1 自纠段逐字「那张表在盘上不存在……那一整段作废」。⇒ 177 至今**一枚候选都没被独立量过**，谁引用它“已裁”都是拿没有的东西当授权。
2. **「163 要人工批准的门还没开」/「163 只是扩既有 `shell.exec` 语义」**——前者已被 `PLAN.md:2568` 冻结文本＋票面 `:4` 更正关掉（本程亲验 `:2568` 在册）；后者被名册现量证伪（`shell.exec` 本就无实现，见 §3.1），写手腿若照“扩既有”起草就会给一枚不存在的工具写“改语义”，落点全错。出处＝本表 §3.1 读数与票 163 `:46`。

## 10. 凭据

本表与台件日志零凭据值抄录。

## 11. next=

163 写手腿派之前还缺：① **`Q-61` 的 owner 一句话**（175/176/177 台账原话已把“不盖戳的起跑口先通”列为禁止形；163 的常驻输出正是同一格的下一个盖戳对象，先落 163 会把 177 的墙重撞一遍）；② 排程约束继承——票 177 `AC#6`：176 不得早于 177 合入，163 的 spawn 腿与 176 同管（〔现读〕`cmd/wisp/run.go:355-356` 注释逐字：「the background spawn and its cancel are AC#4's and ticket 163's land」）；③ 格 A 的落点裁定本表已给，但写手腿合入前须在**自己的锚上**现跑 `gate-clauses.sh` 重登 G6/G7 want_n 基线（本程未跑，基线声明归上一程）。本程交付完毕，无待本程的后续。
