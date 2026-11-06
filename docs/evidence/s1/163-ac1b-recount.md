# 163 AC#1b 重计裁定表 — 名册差集（D34/SPEC-07 §3 ↔ 生产注册表）＋ 残余 AC 决策表

- 时刻：2026-11-06 20:2x +08｜腿名 `163-ac1b`（READ/ANALYSIS，零产码、零票改、零 push）
- 锚：派单说 `1a1f6b344`；**开工现 HEAD ＝ `8bad595370f217c668162bb8359702e9cc54db31`**（branch `dev`，`1a1f6b344` 为其祖先、其上仅 1 笔＝票 65 收口文书），**读数全取现 HEAD**（264-ac0 / 274-ac0 同一定式）
- 隔离（规矩 21）：全部读数出自共享克隆 `D:\tmp\wisp\20261106\163ac1b\repo`（`git clone --shared` 自主树，branch `dev`，起手 `git status --porcelain` ＝ **0 行**；主树在飞未跟踪件一律不在克隆里；alternates 指回主树——**主树 gc 前勿 prune**）。主树本腿零字节改动。
- 派单环境令照办：`export PATH=$PWD/build:$PWD/third_party/sherpa-onnx:$PATH`；`unset GOTMPDIR TMPDIR TEMP TMP`；**未设 `CGO_ENABLED`**（台件自证输出 `# CGO_ENABLED 未设置（默认 1）`）；零 GUI、零 `balldebug`/`wisp` exe、零 `*winlive*`、零 `pending-and-issues.md`、零 `frontend/**`/`design/**`。

## 0. 开工三核（现跑）

| 核 | 读数 | 命令 |
|---|---|---|
| HEAD | `8bad595370f2…` ＝ 派单锚的后一笔 | `git rev-parse HEAD`；`git merge-base --is-ancestor 1a1f6b344 HEAD` ⇒ rc 0 |
| 分支 | `dev` | `git rev-parse --abbrev-ref HEAD` |
| 脏栏 | **0 行**（克隆起手） | `git status --porcelain \| wc -l` |

## 1. 派单要的三个数（全部现跑，尺可逐字复跑）

| # | 数 | 读数 | 命令（克隆内逐字） |
|---|---|---|---|
| 1 | **在册名字（冻结文本）** | **2 枚**：`shell.exec` ＋ `shell.session`（PLAN 与 SPEC-07 镜像各 1 次） | `awk 'NR>=2535 && NR<=2575' docs/PLAN.md \| grep -oE 'shell\.[a-z]+' \| sort \| uniq -c` ⇒ `1 shell.exec / 1 shell.session`；`sed -n '74,75p' docs/specs/SPEC-07-tools-and-plugins.md \| grep -oE 'shell\.[a-z]+' \| sort \| uniq -c` ⇒ 同 |
| 2 | **生产注册表** | **0/2**：两枚 `Lookup` 均 `registered=false`（真 `internal/tools.Registry`，21 枚实名册逐名列出） | `sh .scratch/wisp/probes/163/ac1b/ac1b.sh`（台件 §3） |
| 3 | **代码命中** | `shell.session` ＝ **0**（`internal/`＋`cmd/`＋`tools/`，含测试与注释）；`shell.exec` ＝ **13**（非测试 **6**：注释 2＋`unwired.go` 文案串 4；测试 **7**；**零枚是工具实现**） | `grep -rn 'shell\.session' --include=*.go internal/ cmd/ tools/ \| wc -l`；`grep -rn 'shell\.exec' --include=*.go internal/ cmd/ tools/ \| wc -l` |

**结论一句话：地基仍是 a1 量到的"在册＋零实现"那一形，40 天后一格没动。** 变的只有两处：①注册表从 a1 时代的 7 枚（FS+Task）长到 **21 枚**（今日 run.go 还装 System/Subagent/Plan 三族）——所以本表台件照**今日** `run.go:624-1082` 装配形状补齐，否则名册会漏枚；②行号漂移：`shell.exec`／`shell.session` 两行从 `PLAN.md:2567/2568`（a1 锚）漂到 **`:2574/2575`**——**行文本逐字未变**（`git show 51c32ef9:docs/PLAN.md | sed -n '2568p'` 与 HEAD `:2575` diff ＝ 空）。

## 2. 名册差集裁定表（AC#1b 主件；尺同票 164 AC#1）

| 名字 | 在册（D34＋镜像） | 生产注册表 | 代码命中（非测试） | 处置 |
|---|---|---|---|---|
| `shell.exec` | ✅ `PLAN.md:2574`＋`SPEC-07:74`（L2、白名单命中→L1、默认禁用、argv 强制） | ❌ `Lookup => false` | 6 枚，**全不是实现**：`internal/agent/spill.go:26`（注释）、`internal/config/unwired.go:63/64/70/76`（三枚 unwired 键的 missing/lands 文案串）、`internal/risk/rules_shell.go:14`（注释） | **留名不摘**；主人＝本票实现腿（AC#2 起） |
| `shell.session` | ✅ `PLAN.md:2575`＋`SPEC-07:75`（L2、寿命绑任务、真 deadline、不引 PTY；owner 批准新增，台账 `A322`） | ❌ `Lookup => false` | **0 枚** | **留名不摘**；主人＝本票实现腿（AC#2 起） |

**判据执行**：差集＝2 枚全空、**非空** ⇒ 按票面 `:28` 不许进实现格——**本表就是那份裁定**（形态照 164 AC#1：编排者归口落）。逐枚处置没有一条是"悄悄改名"：

1. **不摘名、不标 DEFERRED**——名字是 owner `A322` 批准后落进冻结文本的（落地笔 `ced72f8f0`，09-27 09:34，`docs/PLAN.md`＋`docs/specs/SPEC-07` 各 +1 行、删除列 0；本表现跑 `git show ced72f8f0 --stat` 复核）。
2. **⛔ "顺手先把 `shell.exec` 注册上"照票面 `:28` 禁止**——名册差集由本表闭格；实现腿开工时两枚名字**同批**由同一枚工具家族注册（它们本就是同一批"命令执行面"，164 AC#1 那次裁定也是这么落 `task.*` 的）。
3. **硬账（照 164 AC#1 那笔原样抄给本票）**：实现腿落地那天若任一枚仍零实现，**它必须变成五字段 DEFERRED 行（同 `PLAN.md §7:1531` 那族形状）或由 owner 批准摘名**——不许继续"在册＋无实现＋无人认领期限"（`Q-58` 前例）。今天 `PLAN.md` 里 shell 家族 DEFERRED 行＝**0**（`grep -n DEFERRED docs/PLAN.md \| grep -i shell` ⇒ 空）。
4. ⚠ **本表没测什么（照 164 AC#1 的自我警示抄）**：只量了"名字在不在册"这一把尺（真注册表 `Lookup`＋注册调用普查）；**没查**别的拼法。但今日比 164 当时多一把现成尺：组合根 `RegisterProvider` 调用＝**0 枚**（`grep -n RegisterProvider cmd/wisp/run.go`），动态注册形状今天结构性不存在，误判风险比当时更低。

## 3. 台件与逐把尺（全部可复跑）

**台件**：`.scratch/wisp/probes/163/ac1b/`（`ac1b.sh` ＋ `main.go` ＋ `logs/{ac1b.out,go-probe-ac1b.log}`）。`main.go` 照 `cmd/wisp/run.go:624-1082` 今日装配顺序组成**真注册表**（`BuiltinFSEntries`→`BuiltinTaskEntries`→`BuiltinSystemEntries`→`BuiltinSubagentEntries`→`BuiltinPlanEntries` 五族；deps 形状逐字段对齐 run.go 同名调用），不用任何手抄名册。日志原文：

```
## 注册表实名册（21 枚）: clipboard.read(L0) clipboard.write(L1) fs.edit(L2) fs.list(L0) fs.move(L1)
   fs.read(L0) fs.trash(L1) fs.write(L1) input.type(L2) media.control(L1) plan.present(L2)
   screen.capture(L1) system.get(L0) system.power(L2) system.set(L1) task.cancel(L1)
   task.output(L0) task.spawn(L0) window.focus(L1) window.list(L0) window.manage(L1)
## Lookup("shell.exec") => registered=false
## Lookup("shell.session") => registered=false
## Lookup("task.output") => registered=true（对照组）
## IsSensitiveSource(shell.exec)=false IsSensitiveSource(shell.session)=false
   IsSensitiveSource(fs.read)=true IsSensitiveSource(task.output)=true
```

对照组两发都真：`task.output`（164 那枚）在册 ⇒ 台件不是"只会答 false"；mark 侧 `IsSensitiveSource` 两枚＝false ⇒ **C25 标记闸（`bridge.go:552` 那条）不会替 shell 家族打标**，与 SPEC-06 §5 名单一致（shell 输出的敏感性走审批面不走 taint 面——这个形状落地时要在判据里写明）。

| 尺 | 读数 | 命令 |
|---|---|---|
| R1 冻结文本 | 2＋2 行在册；a1 锚↔HEAD 行文本 diff＝空 | `grep -n 'shell\.exec\|shell\.session' docs/PLAN.md docs/specs/SPEC-07-tools-and-plugins.md`（命名行 `:2574/:2575/:74/:75`；另 9 枚 `shell.exec` 命中是 D33/D46/§15 论述句，16x-c1 §4.3⑤ 已分类为"规矩断言 4＋论证复述 1＋本体 1＋表格镜像 1"，**一字不用动**） |
| R2 真注册表 | 21 枚实名册；两枚 Lookup false | `sh .scratch/wisp/probes/163/ac1b/ac1b.sh` |
| R3 名字普查 | 13/0（`internal/`+`cmd/`+`tools/`）；全树 22/8（差值全在 `.scratch` 台件） | `grep -rn 'shell\.exec' --include=*.go internal/ cmd/ tools/ \| wc -l` 等 |
| R4 os/exec 面 | 生产 spawn 基元＝`proc.JobScope.StartInJob`（`jobscope_windows.go:110`）；今日生产调用者＝`internal/tools/launcher_{other,windows}.go`（`app.launch`）、`cmd/wisp/slo_windows.go:500`（SLO 台架）、`internal/llm/adaptertest/childjob_windows.go`——**工具面一枚 shell 形都没有** | `grep -rn StartInJob --include=*.go internal/ cmd/ \| grep -v _test` |
| R5 已建好的半边 | **R6 判级器已在册且已接线**：`internal/risk/rules_shell.go`（元字符集＋白名单→L1）在 `internal/risk/assessor.go:271` 组装序列里；`Facts.ShellArgv/ShellString/ShellAllowlist`（`assessor.go:115-121`）；C3 capability `shell`（`PLAN.md:1359`） | `grep -n ruleShell internal/risk/assessor.go` |
| R6 unwired 守卫 | 三枚 `risk.shell_*` 键今天被加载期响亮拒绝，理由逐字＝`no shell.exec tool is registered…`（`internal/config/unwired.go:60-78`）；旁边钉着 `TestEveryLockedSectionKeyIsAccountedFor`（`unwired_test.go:311`） | `sed -n '60,78p' internal/config/unwired.go` |
| R7 DEFERRED | shell 家族＝0 行 | `grep -n DEFERRED docs/PLAN.md \| grep -i shell` |
| R8 门禁 | `gofumpt -l .scratch/wisp/probes/163/ac1b/main.go` ⇒ **空**（`gofumpt v0.12.0`，`/d/work/base/gopath/bin/gofumpt.exe`）；`sh scripts/d22scan.sh` ⇒ **rc 0 clean**（分母 `bans#1-5 internal/=366 cmd/=58、ban#8 internal/=912 cmd/=217、frontend/=116、design/=14`）；`sh scripts/check-path-length-budget.sh` ⇒ **rc 1 RED＝既有腐坏，非本腿**（见 §4） | 台件目录内逐字复跑 |

## 4. 门禁 RED 的归属（第 21 处定式：红归谁）

`check-path-length-budget.sh` 报 **RED：roster entry that is not an over-budget tracked path any more** ＝ `.scratch/wisp/issues/185-every-successful-reread-mints-…-c25-mark.md`。归属两发现跑：

1. **把本腿台件整个挪出克隆（parked）→ porcelain 0 行的纯 HEAD 上重跑 ⇒ 仍 rc 1** ⇒ RED 与本腿零关系；
2. 根因＝**roster 漂移**：那枚超长票名 10-04 改短时加了 `-done` 后缀（`git ls-files '.scratch/wisp/issues/185-*c25-mark.md'` ⇒ 0 枚、`185-*` ⇒ `…-c25-mark-done.md`），roster 里还记着旧名 ⇒ 归 **票 262/台账**域，本腿不修（改 roster 不在本腿写面）。本腿自己件名 `.scratch/wisp/probes/163/ac1b/**` 全部远低于帽。

## 5. 残余 AC 决策表（哪些现在能做、哪些等 owner、写面各是什么）

| AC | 内容一句话 | **现在能实现吗** | 卡点／红线 | **写面（具体到落点）** |
|---|---|---|---|---|
| AC#2（`:29`） | 埋标记取回退出码，不许靠时间；答"摘掉标记是否存在一发超时从此判错" | **能**——零契约缺口、零批准缺口 | ①主通道必须是 `Wait` 状态的字段（a1 已量：`errors.As(*exec.ExitError)` 命中 `ExitCode()=7`；"文本猜"有活样本会判错）；②超时＝真 deadline（ban #4：禁墙钟差，走 `context` 取消/Deadline 形）；③摘标记变异两向都要真跑 | 新建 `internal/tools/shell_*.go`（＋`*_test.go`）；会话进程载体落 `internal/proc/**`（**复用自家 `proc.JobScope.StartInJob` 形状**——a1 那句"别照抄 Windows ConPTY 那一半"的仓内正解就是 Job 管控这一族）；**risk 侧零改**（R6 判级器已在册接线，§3 R5） |
| AC#3（`:30`） | 会话寿命绑任务作用域：任务取消/结束 ⇒ 会话必须跟着收，拿"取消后会话还在"变异证有牙 | **能** | ①goroutine 只许 `observe.Registry.Spawn`（`internal/observe/goroutine.go:262`，裸 `go func(` 是 ban #1）；②句柄缝现成＝`TaskRoster.AttachCancel`（`internal/tools/task.go:420`；今日生产调用点仅 `subagent_197.go:374` 一枚）＋`TaskRoster.Cancel`（`task.go:447`）；③变异证明由实现腿自造两向 | 同 AC#2 一条腿；判据里点名"取消⇒会话进程/Job 句柄收到"的可观察读数 |
| AC#4（`:31`） | 审批不许松：常驻会话不许变成"授权一次就一直放行" | **能，但带一条硬红线** | 红线＝**不许动 `internal/agent/approval/**`**（票面逐字，要动先回来报）。形状上每发命令各判一次走既有 C19/C18 桥路（`internal/tools/bridge.go:339` 每调用 `assessorFor(req.TaskID).Assess`）——D45 会话授权（`internal/tools/grant.go` 的 `GrantSource`）**不许**被会话拿来当"一次批准、终身放行"的载体 | 零字节 `approval/**`；只动 tools/proc 侧；判据写"同会话第二发仍过 C19" |
| AC#5（`:32`） | 契约轴 | **实现本身＝零契约改动**（A322 已落 `shell.session` 行；`shell.exec` 行与 D33 `:2515`、D46 `:2674` "保持不变"句一字不用动） | **一格空白需要 owner 一句话**：`internal/config/unwired.go` 三行 missing 文案（`risk.shell_enabled`／`allow_shell_string`／`shell_allowlist`）在工具落地当天**开始撒谎**（"no shell.exec tool is registered"将成假话），而 `internal/config/**` **既不在票面允许写面、也不在零字节名单**（16x-c1 §4.3⑥② 就点名过这格）＋旁边钉着 `TestEveryLockedSectionKeyIsAccountedFor`，新增 🔒 键不登记就红 | 待 owner 批（16x-c1 §4.3⑦ 原句照用）：**"放开 `internal/config/unwired.go` 里那三行 `shell_*` 的 `missing` 文案（只改文案、不删守卫、不动阈值）"**；不批＝要么复用 `shell_enabled` 当会话总开关（同样要改那三行）、要么实现腿到这一步**停手上报** |
| AC#1b（`:28`） | 名册差集裁定表 | **本表即裁定，本格可闭** | —（差集非空的判据已执行：裁定先于实现） | 即本件 |

**排序建议（给编排者）**：AC#1b 闭格 ⇒ **AC#2＋AC#3＋AC#4 同一条实现腿**派（同一族文件、判据互相咬，拆两腿会在同一批文件上互相踩——README 规矩 2 的"同票多写手互踩"条款）；**AC#5 那格 unwired 文案先摆 owner**（三选一：批文案放开／批复用 `shell_enabled`／维持现状＝实现腿到注册那步停手上报），**不必等它**才能开工——工具实现与测试全在允许写面内，只有"注册进组合根"那一步撞三行文案，实现腿可以把那一格留成最后一笔。派单里照 164 先例写两条：①写腿开工前现跑本表 §1 三数对账（名字若有第三者动了，以它现量为準回来推翻本裁定）；②**顺手修票面过期行号**（a1 表引 `:26-:29`、本表引 `:28-:32`、`PLAN.md:2567/2568`→`:2574/2575`）归编排者，不归写腿。

## 6. 本表没测什么（逐格照实）

- `bridge.Execute` 审批链路的端到端（需接 approval UI）未构造——AC#4 的"每发各判"只给了**缝的存在**（`bridge.go:339` 逐行读码），没给实跑读数；
- `gate-clauses.sh`（`.scratch/wisp/probes/154/`）未执行——G5/G6/G7 基线只做了文本核对，未跑尺（同 a1 §3.3 的口径）；
- `internal/proc` 里有没有可直接复用的"长活进程＋收尸"现成件，只做了 **StartInJob 调用者普查**（§3 R4），没逐文件读 `boot_windows.go`／`shutdown.go` 的收尾链——实现腿开工时第一件就是读它，别按本表推断直接写；
- 退出码/超时/取消三判据的变异两向——实现腿的活，本腿零测码。

## 7. 临时件清单（规矩 8：只建不删）

| 件 | 说明 |
|---|---|
| `D:\tmp\wisp\20261106\163ac1b\repo` | 规矩 21 共享克隆（alternates 指回主树，**主树 gc 前勿 prune**）；本腿一切读数出自这棵树；提交见 §8 |
| `D:\tmp\wisp\20261106\163ac1b\parked\`（空） | 归属两发用的暂挪点，已空 |
| `/tmp/row_a1.txt`、`/tmp/row_head.txt` | 行文本 diff 的两枚碎片（几行文本，未引用即弃） |

仓内件＝`.scratch/wisp/probes/163/ac1b/{ac1b.sh,main.go,logs/*}`（台件，随本腿提交）＋ `docs/evidence/s1/163-ac1b-recount.md`（本件）。

## 8. 撤销口令

**「撤 163-ac1b 名册裁定表」**＝把克隆里那笔提交 revert 掉即可：diff 恰好 5 枚文件（本件＋台件 4 枚），产码零字节、票面零字节、主树零字节（提交在克隆，由编排者收件后按规矩并入 dev）。
