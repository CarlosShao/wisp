# 377-r1 落地：记账面报**有效档** —— `tools: call` 加 `mode=`，`Decision.Mode` 改读 `perm.Store`

- 写腿 `377-r1`（实现者）。基线 HEAD **`5eddc6223858e6b06a35bb28be59bad7661b6e53`**（`375-v1` 收口后）。
  本笔**只 commit ⛔ push**；票面／`Status:`／台账／`-done` **一字未动**。
- 起点＝`docs/evidence/s1/377-a1-census.md`（160 行）＋票面末尾「`AC#1` 裁形（改后）＝〔丁〕」。**未重做普查的量**：
  11 个读档点、`Decision.Mode` 的 0 消费者、`m1`/`m2` 两枚变体，全部照采（本笔只复跑 `m1`/`m2` 的结论面，见 §⑥）。

## ① 为什么选〔丁〕（引普查三形代价）

普查 §⑦ 原表：〔甲〕改 2 处（照字面 4 处）→ **破 `m1`＋`m2`**、且**丢「配的是什么」**；〔乙〕`Decision` 加字段＋5 处＋formatter 1 处 → 不破，**但得同时补读者**；〔丙〕审计 formatter 2 处 → 不破（展示面已全是有效档）。
编排者裁〔丁〕＝**把〔乙〕的"加字段"落到真取证面上**：`tools: call` 行加 `mode=`（普查 §④：三处最大取证面 **全无 mode**、`Decision.Mode` **0 读者**）＋`Decision.Mode` 顺手读有效档。
⇒ 本笔**不新增 `Decision` 字段**、**不动展示面**、**不动 `perm.New`／`perm.Set` 的记账**，因此 `m1`（严档请求被吞）与 `m2`（`MODE-DOWNGRADE` 整条消失）**从射程外**。

## ② 改了什么（逐处；行号为**落盘后**现取）

| # | `file:line` | 改动 | 必须在这层的理由 |
|---|---|---|---|
| 1 | `internal/tools/bridge.go:1201` | `tools: call` 行**末尾**加 ` mode=%s`，实参 `dec.Mode` | 该行由 `book()` 独占写，`Bridge` 是每发调用的唯一记账者；放别层就得复制一份决策 |
| 2 | `internal/tools/bridge.go:286`（新块 `276-285`） | 有效档**提到 `Execute` 顶部读一次**（`mode := b.permissionMode()`），未知工具分支 `Decision{Tool, Mode: mode}` | 原来只在 C3＋参数解码之后读：**能力/坏参数被拒的那两发根本没有可记的档**，而第 1 处每发都印 |
| 3 | `internal/tools/bridge.go:297` | `dec, why := b.checkCaps(...)` 之后 `dec.Mode = mode` | 让 `reject()`/`close()`/`run()` 三条记账路径共享**同一个**值（一次读、全链一致） |
| 4 | `internal/tools/bridge.go:327` | 删掉原 `mode := b.permissionMode()`，改 `sil := mode.Screen(verdict)`（`dec.Mode = mode` 保留） | 同一次读既喂筛选又喂记账 ⇒ 行上的 `mode=` **就是**筛选用过的那个值，不是第二次采样 |
| 5 | `internal/tools/bridge.go:1188-1199` | `book()` 文档块加 `mode=` 段的理由（有效档 vs config.toml；**追加在 `reason=` 之后**的兼容理由） | 该文档块是本仓库记账面的规格所在（`rules_hit`/`grant_id` 的边界都写在这） |
| 6 | `cmd/wisp/approval_always.go:168-174` | `permissionMode()` 改读 `rt.modes.PermissionMode()`；`rt.modes == nil` ⇒ `risk.DefaultMode()`（删掉 `rt.mgr.Config().PermissionMode()` 兜底） | 两枚卡（`runWidening:122`、`confirmLockedLoosening:294`）**只调这一个函数**；改一处即两处 |
| 7 | `cmd/wisp/config_reload.go:278-284` | 只改**注释**（"read through Manager.Config()" → "perm.Store 的有效档，其配置半边仍到 Manager.Config()"） | 行为零变化；不改就是一句会误导后来人的注释 |

`perm.Store` 的 `from/to`（`store.go:254`）与启动 `MODE-DOWNGRADE` 落行点（`store.go:158`）**一字未动**（§⑦ 硬边界）。`cmd/wisp/run.go` **零改动**（`consoleApprovalUI.Prompt` 未碰）；`internal/agent/approval/reach.go`／`gate.go` **零改动**。

## ③ ★`AC#2` 那发（判据：`tools: call` 行必须读出 `mode=ask_every_step`）

装配＝`wisp panel-inbound` 形状，**真对象**：磁盘 `config.toml`（`auto_approve`）＋`config.NewManager`＋`perm.New{Confirm: nil}`＋`perm.Set(ask_high_risk)`＋真 `BuiltinFSEntries`/`C26`/`Bridge`。用例 `Test377PanelInboundBookingNamesTheEffectiveMode`。

前提三格（用例内先钉死，先红的是前提不是断言）：
`perm.New` ⇒ `MODE-DOWNGRADE`；`Set` 后 `config.toml = ask_high_risk`；`store.PermissionMode() = ask_every_step`（切档前后都不动）。
读数（逐字）：

```
tools: call kind=success task=task-377 corr=corr-377 tool=fs.write risk=L1 decision=allow outcome=success rules_hit=[R1] in_allowlist_scope=true grant_id=0 refusal=none reason="R1: 工具声明为下界（L1）" mode=ask_every_step
```

⇒ 行里那个值**既不是** `auto_approve`（开工时文件）**也不是** `ask_high_risk`（切档后文件）：读配置档的两种写法都会印出其中之一，用例把这两个**错误答案逐名断言为红**（`mode=auto_approve`／`mode=ask_high_risk`）。所以这格**不是**"字段存在"，也**不是**只测 `Decision.Mode`。

## ④ 两发对照（普查已给，本笔复跑；只差有没有腿）

| 发 | `Decision.Mode` | `tools: call` 的 `mode=` | gate | 审计旁证 |
|---|---|---|---|---|
| `auto_approve`＋`confirm=nil`＋一次 L1 | `ask_every_step` | **`mode=ask_every_step`** | **1 个 L1 窗**、执行成功 | 无 `MODE-SILENCE` |
| `auto_approve`＋`confirm=leg`＋一次 L1 | `auto_approve` | **`mode=auto_approve`** | 0 窗、被静默 | `tools: MODE-SILENCE mode=auto_approve assessed=L1 effective=L0` |

用例把两张表**互相钉住**：`MODE-SILENCE` 的存在性必须与本发 `mode=` 的值同向（`Test377ConfiguredVersusEffectiveTwoShots` 的 `wantSilence`），所以行上的 token 不能靠硬编码过。

## ⑤ `tools: call` 行改前／改后（逐发）

- **改前**（HEAD `5eddc62` 的 `bridge.go`，`-overlay` 换回，逐字）：`... refusal=none reason="R1: 工具声明为下界（L1）"` ⇒ **行尾没有 `mode`**（`strings.Contains(line,"mode=")==false`）。
- **改后**：`... reason="R1: 工具声明为下界（L1）" mode=ask_every_step`（§③）／`mode=auto_approve`（§④ 第二发）。
- 字段位置：**行尾**。理由：仓里至少两处把行内字段当**连续子串**匹配（`internal/tools/ticket363_test.go:258` 的 `"tool=fs.read risk=L0 decision=allow outcome=success"`、`internal/tools/fs_edit_ac5_gate_r4_test.go:206` 的 `"risk=L2"`,`"R2"`,`"in_allowlist_scope=false"`），插在中段会破不该被本票碰的读者。

## ⑥ 正控：`wisp run`（`Confirm` 永远非 `nil`）那条 seam 仍报**配置档**

真装配读数（仓外探针，`-overlay` 注入，不进交付面；`runTextTask` ＋ 挂了回复管的有腿跑，`config=auto_approve`，一次 L1 写）：

```
[audit] tools: call kind=success task=p377-run corr=p377-run-corr tool=fs.write risk=L1 decision=allow outcome=success rules_hit=[R1] in_allowlist_scope=true grant_id=0 refusal=none reason="R1: 工具声明为下界（L1）" mode=auto_approve
```

⇒ 该 seam 上 `confirm = rt.confirmModeSwitch`（方法值、恒非 nil）⇒ 372 的读侧塌陷不触发 ⇒ 配置档＝有效档，行仍报 `auto_approve`，**本腿没把它改成"一律 ask_every_step"**。旁证：`cmd/wisp` 的两枚单测（`…KeepsTheConfiguredModeWhenAChannelExists` 绿、`…WithoutAStoreFailsClosed` 绿）。
边界说明（本腿**不改**）：`wisp run` 给**桥**的不是裸 store 而是 `approval.UnsilencedMode(rt.modes, answer, …)`（`run.go:864`）；无人可答时桥按 `ask_every_step` 筛（375），而 `permissionMode()` 读的是 store 的 372 有效面 ⇒ 无人可答的 console 跑里，**卡上**的 Mode 会比**桥上**宽一格。这属 375 的谓词面（要么把 `answer` 引进 runtime、要么让卡也用包装源），不在〔丁〕的裁形里，本笔只记录不擅动。

## ⑦ 每处改动配的用例（逐名＋为什么今天响）

新文件 2 枚（`internal/tools/booked_mode_377_test.go`／`cmd/wisp/permission_mode_effective_377_test.go`），共 6 个用例：

| 用例 | 配哪处改动 | 今天为什么响（改前红行逐字） |
|---|---|---|
| `Test377PanelInboundBookingNamesTheEffectiveMode` | ②#1/#2/#3/#4 | `the booked line does not name the effective档 "ask_every_step": … reason="…"`（行尾无 mode） |
| `Test377ConfiguredVersusEffectiveTwoShots/{nil,leg}` | ②#1/#4 | `mode= is not "ask_every_step"`／`mode= is not "auto_approve"` |
| `Test377UnknownToolBookingStillNamesAMode` | ②#2（**提读**那一跳） | `the pre-C3 refusal booked no effective档: … reason="…"` |
| `Test377BookedModeReadsTheEffectivePermissionMode` | ②#6 | `the widening card would book auto_approve for a process that screens at ask_every_step` |
| `Test377BookedModeKeepsTheConfiguredModeWhenAChannelExists` | ②#6 的**正控**（改前也绿，刻意） | —（防"改成硬编码最严档"） |
| `Test377BookedModeWithoutAStoreFailsClosed` | ②#6 的 fail-closed 半 | —（改前也绿，刻意） |

`cmd/wisp/config_reload.go` 那处**只改注释**（②#7）：无行为，故不配用例。

## ⑧ 反向判据：三读数 ＋ 还原核（时机在**最后一发变异之后**）

全部变异走 `-overlay`（仓外 `D:\tmp\wisp\20261024\377r1\`），**工作树零字节突变窗口**。

1. **变异红 M1**＝把 `book()` 新加那一跳（`mode=%s`＋`dec.Mode` 实参）摘掉 ⇒ `internal/tools -run Test377` **rc=1**，三枚用例逐名红（红行见 §⑦ 第一、二、三行；行尾回到"没有 mode"）。
2. **变异红 M2**＝让 `perm.Store.PermissionMode()` 回到"配置档"（`return s.configuredMode()`，overlay，仓内不改）⇒
   `Test377PanelInboundBookingNamesTheEffectiveMode`: `premise broken: effective档 = auto_approve, want ask_every_step`；
   `Test377ConfiguredVersusEffectiveTwoShots/nil`: `windows = 0, want 1` ＋ `mode= is not "ask_every_step": … mode=auto_approve`（**行上印的正是配置档**）。
3. **改前红 M0**＝`-overlay` 把 `internal/tools/bridge.go`＋`cmd/wisp/approval_always.go` 换回 HEAD ⇒ `internal/tools -run Test377` rc=1（三红）、`cmd/wisp -run Test377` rc=1（`the widening card would book auto_approve …`）。
4. **还原绿**（最后一发变异之后）＝裸工作树 `internal/tools -run Test377` **rc=0/0.179s**、`cmd/wisp -run Test377` **rc=0/0.143s**。
5. **原样绿**＝全包／定向读数见 §⑨。
6. **还原核**（`2026-10-24T17:31:35+08:00`，**在 M0 之后**量；`17:37:15` 在**全部跑完之后**（含仓外探针）**复量一次**，两次数值逐字相同）：五枚文件 sha256 `bridge.go 952380bc…`、`approval_always.go e55cad27…`、`config_reload.go c22708f8…`、两枚测试 `48ea5fcc…`／`fcc69adc…`；`git diff --numstat`＝`26/6`、`5/2`、`37/9`（与变异前同）；`git status --porcelain -- internal/tools cmd/wisp` 只有本腿 3 改 2 新，overlay/变异件**全在仓外**。

## ⑨ 门禁 rc（改前／改后**分开**报，⛔ 不相减）

| 命令 | 改后（工作树） | 改前（`git archive HEAD` 导出树） |
|---|---|---|
| `go build ./...` | **0** | 0（导出树 `go test` 已含编译） |
| `go vet ./internal/tools ./internal/perm ./internal/agent/approval ./cmd/wisp` | **0**（无输出） | —（同左） |
| `go test -count=1 ./internal/tools` | **0** `ok 21.326s` | **0** `ok 20.492s` |
| `go test -count=1 ./internal/agent/approval` | **0** `ok 0.727s` | **0** `ok 0.678s` |
| `go test -count=1 ./internal/perm`（未改，旁证） | **0** `ok 0.571s` | **0** `ok 0.544s` |
| `go test ./cmd/wisp -run '<9 名roster>'`（带 `PATH=$PWD/third_party/sherpa-onnx`，只定向） | **0** `ok 207.246s` | **0** `ok 192.763s` |
| `go test ./cmd/wisp -run TestTicket224`（另一枚读者：`tools: call ` 作滤针） | **0** `ok 11.994s` | 未跑（同族，见注） |
| `go run . -root <repo>`（`tools/d22scan`，独立 module） | **0** clean；292 枚生产 `.go`（`ban #1–7` 违例 0、`#8` 违例 0） | — |
| 格式**甲尺** `/d/work/base/gopath/bin/gofumpt.exe -l <5 枚>` | **rc=0、输出 0 字节** | — |
| 格式**乙尺** `gofmt -l <5 枚>` | **rc=0、输出 0 字节** | — |
| 长度门：本件 | **≤160 行**（件尾自报） | — |

roster＝`TestAlwaysBranch|TestAC1AlwaysBranch|TestTicket223|TestTicket114r2|TestTicket101|Test375|Test376R1|Test377|TestReplyListenerReject`。
注①：导出树**第一次**跑 `internal/tools` 红 9 枚（8 枚 252 短名前提＋1 枚 309 走 `../../scripts`）——**全是导出树缺件/路径形状**，补 `scripts`＋`frontend` 并换到带 8.3 别名的路径后 **rc=0**（`before-tools-clean2.txt`）。⛔ 拿那 9 枚当 HEAD 缺陷。
注②：`gofumpt -l .` 未用（会扫 `.scratch` 故意写坏的探针 ⇒ `rc=2`）；EOL 另量：五枚 `0x0D`＝**0**。

## ⑩ 硬边界逐条自检

1. ⛔ 动 `perm.Store.Set` 的 `from/to`（`store.go:254`）：**未动**（`git status --porcelain -- internal/perm` 空）。⛔ 动 `store.go:158`：**未动**。`m1`/`m2` 两枚变体只在核心结论面复跑（M2 即"配置档"形），**不落地**。
2. ⛔ 改机主 2026-09-21 `R20/M4` 定案：**未动**（`RequiresCard`／`Set` 的确认规则零字节变化）。
3. ⛔ 动 `internal/{winsec,memory,panel}/**`：`git status` 在此三树**空**。
4. ⛔ 改 `cmd/wisp/run.go` 的 `consoleApprovalUI.Prompt`、⛔ 动 `internal/agent/approval/reach.go`／`gate.go`：**三者零改动**（`git diff --numstat` 只列 §② 三枚）。
5. 提交面：`internal/tools/bridge.go`、`internal/tools/booked_mode_377_test.go`、`cmd/wisp/approval_always.go`、`cmd/wisp/config_reload.go`、`cmd/wisp/permission_mode_effective_377_test.go`、本件。⛔ `-A`／⛔ glob／⛔ amend/reset/rebase/stash/checkout。
