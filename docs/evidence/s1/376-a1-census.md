# 376-a1 普查：L2 卡印 `窗口 0.0s` —— 「卡上那行字」与「真等的超时」是两个来源（只读腿）

- 腿 `376-a1`（只读）。读数时 HEAD `8ced487bb`（2026-10-24 12:26 +08）。本笔**未落任何 `.go`**、未改票面/台账/`Status:`、未 push。
- 探针全在仓外 `D:\tmp\wisp\20261024\376a1\`：`zz376a1_probe_test.go`（6 用例）经 `-overlay` 注入 `cmd/wisp`；三枚变体 overlay 另把 `gate.go`/`run.go` 的**仓外副本**换进编译，仓内零字节改动。读数原文件：`head-run.txt`／`jia-run.txt`／`yi-run.txt`／`yip-run.txt`／`bing-run.txt`／`floor31-run.txt`。
- ★**复跑钉 HEAD 字节**：读数期间 `cmd/wisp/run.go` 被**票 373 的写腿**改过（工作树 mtime 12:31:51，改点在我引的行号之前）⇒ 全部读数用 `pin_*.json` 把 `gate.go`＋`run.go` **钉回 HEAD blob 副本**再跑一遍，数字与首轮逐项一致（`pin-*-run.txt`）。**行号永远对 `git show HEAD:` 取，⛔ 对工作树。**
- ★**尺的射程（落笔前五问）**：①对象＝`approval.Prompt` 的**两个字段**及其唯一印发点，⛔ 不是 Win32 的 window；②形状＝`file:line`＋读/写＋那一处的语义，行号一律 `git show HEAD:` 现取；③数的是**代码**：`Window|Deadline` 原始命中 **5569** 行（其中 3484 行在 `_test.go`，且绝大多数是 `CreateWindowExW`/`ContextWindow`/`FocusWindow`）⇒ 再用字段访问尺收窄到**生产** 16 行，两字段共 **9** 条（含 2 条「没有这个读者」）；④总数核对：生产 `.go`（非 `_test.go`）292 枚，逐条表 9 行、测试组 5 行，与收窄后的数相等；⑤秒数＝`time.Duration.Seconds()` 现算，卡面按 `%.1f` 打印、故一位小数。

## ① 两处逐字（现跑，HEAD blob）

```
$ git rev-parse --short=9 HEAD
8ced487bb
$ git show HEAD:internal/agent/approval/gate.go | sed -n '534p'
	p := g.promptFor(d, corr, 0, g.q.Timeout(), nil, nonce)
$ git show HEAD:cmd/wisp/run.go | sed -n '1563p'
	fmt.Fprintf(u.out, "  窗口 %.1fs\n", p.Window.Seconds())
```

- ★**行号未漂**：`gate.go` blob `78e2d28a8…`、`run.go` blob `7b9c285db…` 在 `48d193298`（374-a1 基线）与 HEAD **逐字节相同**（`git rev-parse <rev>:<path>` 两两相等；`git log 48d193298..HEAD -- <两文件>`＝**0 行**，区间 11 个提交）。票面 534／1563 两枚号对 HEAD 成立。
- 相邻逐字（同一函数族）：`gate.go:597` 签名 `func (g *Gate) promptFor(d tools.Decision, corr string, win, deadline time.Duration, bv *BatchView, grant string) Prompt {`；`:609` `Window: win,`；`:610` `Deadline: deadline,`；L1 那一路 `:293` `p := g.promptFor(d, corr, g.window, 0, bv, "")`。

## ② `Prompt.Window` 与 `Prompt.Deadline` 的**全部**读写点

声明：`internal/agent/approval/ui.go:38` `Window time.Duration`／`:39` `Deadline time.Duration`（两字段无逐字段注释）。

| # | `file:line` | 字段 | 读/写 | 那一处的语义 |
|---|---|---|---|---|
| 1 | `internal/agent/approval/gate.go:609` | `Prompt.Window` | **写** | 写入 `promptFor` 的第 3 参数 `win` |
| 2 | `internal/agent/approval/gate.go:610` | `Prompt.Deadline` | **写** | 写入第 4 参数 `deadline` |
| 3 | `internal/agent/approval/gate.go:293` | 两者（实参） | **写** | L1 那一路：`win=g.window`（2–3s 块长）、`deadline=0` |
| 4 | `internal/agent/approval/gate.go:534` | 两者（实参） | **写** | L2 那一路：`win=`**字面 0**、`deadline=g.q.Timeout()`（5m0s） |
| 5 | `internal/agent/approval/replies.go:181` | `ReplyCard.Window` | 读 `p.Window`／写 | 逐字拷贝（`replies.go:83` 注释：L1 块长、L2 为 0） |
| 6 | `internal/agent/approval/replies.go:182` | `ReplyCard.Deadline` | 读 `p.Deadline`／写 | 逐字拷贝（C18 上限） |
| 7 | `cmd/wisp/run.go:1563` | `Prompt.Window` | **读** | ★**全仓唯一生产读者**：打印 `窗口 %.1fs` |
| 8 | **（无）** | `Prompt.Deadline` | **读** | ★**生产 0 个读者**（只在测试里被断言） |
| 9 | **（无）** | `ReplyCard.Window`／`.Deadline` | **读** | ★0 个读者：`#5/#6` 拷进账本后**再没人读**（死拷贝） |

测试侧的读/写（5 行，全列，⛔ 抽样）：`internal/agent/approval/window_test.go:30-31` 读 `p.Window`（L1 断言 `==3s`）；`internal/tools/wiring_test.go:133-134` 读 `cards[0].Window`（L1 断言 `>=2s`）；`internal/agent/approval/batch_test.go:105-106` 读 `p.Deadline`（L2 `==300s`）；`internal/agent/approval/queue_test.go:48-49` 读 `p.Deadline`（`==300s`）；`internal/agent/approval/replies_201_test.go:32-33` 写 `ReplyCard{Window,Deadline}` 字面量。**没有一枚读 L2 的 `Prompt.Window`。**

★**同名不同义（逐条核过，⛔ 混）**：`approval.Options.Window`（`gate.go:26-29,143`，被 `gate.go:143-151` 夹到 `[2s,3s]`；写点 `cmd/wisp/run.go:633`、`cmd/wisp/resident_approval_windows.go:398`）＝**配置进来的 L1 块长**，与 `Prompt.Window` 同族、⛔ 不同字段；`Gate.Window()`（`gate.go:176`）＝它的读面（`resident_approval_windows.go:410` 与 4 枚测试读）；`tools.GateFuncs.Window`（`internal/tools/gate.go:122,128,131`）＝**L1 回调函数字段**，只是撞名。

★★**还有没有别处也印 `Window`？——没有。** 全仓卡片印发点**仅** `run.go:1563`（`窗口` 在 `internal/`＋`cmd/` 的 `.go` 里共 57 行命中、含测试，只有这一行是打印时长）。逐面核过：
- 球卡 `ballCardUI.Prompt`（`cmd/wisp/resident_approval_windows.go:890-950`）打印 level/tool/编号/orb 状态/取消通道，**不印时长**；
- 面板卡 `internal/panel/approval.go:39-59` `ApprovalCardView`（以及 `ui.go:50-57` `PanelItem`）**根本没有时长字段** ⇒ 页面印不出它没收到的数；`L1WindowWait`（`internal/panel/pump.go:127-130`）只有两个 id；
- 唯一同时说出两个数的是 **boot 审计行** `resident_approval_windows.go:410`：`gate_window=… gate_queue_timeout=…`（slog，非卡，且**两个都对**）；
- `frontend/` 1367 文件里 `窗口` 11 命中，与审批卡倒计时**无一相关**。

## ③ 从配置到卡面：逐跳

1. `internal/config/schema.go:460` `ConfirmTimeoutSec`，标签 `toml:"confirm_timeout_sec"` `default:"300"`（`l1_window_sec` 默认 **2**，`:463`）
2. `internal/config/validate.go:127-130` band `[31,3600]`，`:145-149` 越界**拒绝加载**（不夹、不降级）
3. `cmd/wisp/run.go:634` `ApprovalTimeout: time.Duration(cfg.Risk.ConfirmTimeoutSec) * time.Second`（控制台腿）；常驻腿同式 `resident_approval_windows.go:541-542`
4. `internal/agent/approval/gate.go:158` `q: NewQueue(o.ApprovalTimeout, …)`
5. `internal/agent/approval/queue.go:102-108` 零值→`DefaultApprovalTimeout`(300s)，warn 失效→`DefaultApprovalWarning`(30s)
6. `internal/agent/approval/queue.go:144` `Timeout()` 返回 `q.timeout` ← ★**乙源**
7. `gate.go:526` `deadline := g.clock.After(g.q.Timeout())` ← ★**真正起表的 deadline**
8. `gate.go:571-574` `case <-deadline: a := g.q.expire(it)` ← ★**真正在等的两行**
9. `gate.go:534` 同一函数里 `win=`**0**、`deadline=g.q.Timeout()` ← ★**分叉点（第 3 参数 vs 第 4 参数）**
10. `gate.go:609-610` `Prompt{Window:0, Deadline:5m0s}`
11. `cmd/wisp/run.go:1563` `p.Window.Seconds()` ⇒ 印 `窗口 0.0s` ← ★**甲源**

★★**两个来源，具名**：**甲源**＝卡面那行字＝`Prompt.Window` ← `promptFor` 第 3 参数 ← `gate.go:534` 的字面 `0`；**乙源**＝进程真正等的＝`q.Timeout()` ← `Queue.timeout` ← `[risk].confirm_timeout_sec`（同一个值也作为第 4 参数进了 `Prompt.Deadline`，但**生产没人读**）。同一次调用，第 3／第 4 参数分叉，这就是 `0.0s` 与 5 分钟的全部来路。

**现跑读数（HEAD，真实 printer＋真实 `approval.Gate`）**：
- 默认配置：`CONFIG_DEFAULT confirm_timeout_sec=300 l1_window_sec=2` ⇒ `QUEUE_TIMEOUT=5m0s (300.0s)` ⇒ `PROMPT level=L2 window=0s(0.0s) deadline=5m0s(300.0s)` ⇒ `CARD_LINE="  窗口 0.0s"`。
- 带底 31（`ticket 267` 下限）：`FLOOR_QUEUE_TIMEOUT=31s` ⇒ 卡仍 `"  窗口 0.0s"` ⇒ 实测 `elapsed=31.004s`、`answer=timeout`、`why="审批超时（31 秒未确认），C18 一律判拒绝，已自动拒绝"`；`[warning] 审批将在 30 秒后自动拒绝` 在 **1s** 时打出（31−30）——★**卡上写 0.0 秒、真等 31 秒；默认配置下差 5 分钟**。
- 与 374-a1 复核：默认 300／`5m0s`／带底 31s 三口**逐项对上**；其 `elapsed=31.001s` 与我 `31.004s` 差 3ms（都是墙钟）。其 `result=refused-confirm` 是 `perm.ResultRefused`（`internal/perm/store.go:70`）——那条 L2 请求是**切档**；我这条是工具腿，桥把 L2 超时映射成拒（`internal/tools/bridge.go:484-488` `AnswerTimeout`→`DecisionTimeout`、返回 false）。两处都是「拒、⛔ 允许」，超时链同一条。

## ④ 三形代价（逐形实测，⛔ 推演）

| 形 | 改哪 | `git diff --numstat` | L2 卡 | L1 卡 | 会不会破别处 |
|---|---|---|---|---|---|
| **甲** | `gate.go:534` 的 `0` → `g.q.Timeout()` | **1/1** | `300.0s` ✓ | `2.0s` ✓ | 不改任何**读者**（`Prompt.Window` 只有 `run.go:1563` 一个）；但 `replies.go:83` 的文档「zero on an L2 card」**变假** ⇒ 还得改 1 行注释；`ReplyCard.Window` 会带 300s（**0 读者**，无行为影响） |
| **乙（字面）** | `run.go:1563` `p.Window` → `p.Deadline` | **1/1** | `300.0s` ✓ | ★**`0.0s` ✗** | ★**破 L1**：L1 的 `deadline` 实参是 0（`gate.go:293`）⇒ 把 L2 的镜像缺陷搬到 L1 卡上（实测） |
| **乙′（按 level）** | `run.go:1563` 1 行→5 行（`shownWindow := p.Window; if p.Level=="L2" { shownWindow = p.Deadline }`） | **5/1** | `300.0s` ✓ | `2.0s` ✓ | 生产者不动、字段语义不动；改的正是唯一读者自己 |
| **丙** | 甲＋乙′ | **6/2** | `300.0s` ✓ | `2.0s` ✓ | 不破；但若「`Window` 语义改成**剩余时间**」按字面做，需要一个**活的倒计时源**——今天没有（`internal/panel/pump.go:50-53` 具名 `approval.EventTick`「declared and emitted nowhere」）⇒ 丙≠1 行 |

★**判（按「最有价值＋主流 harness」，⛔ 按「动作最小」）**：**乙′**。主流审批卡的倒计时是**渲染层按 level 选源**（L1 显示块长、L2 显示 C18 截止），生产者语义不动；甲只省一个分支，代价是把一个**已写进文档的不变量**（L2 上 `Window==0`）改成随 level 变化的两义字段，且省下的分支将来每加一个读者都要重新判 level。丙在今天拿不到「剩余时间」，因为倒计时生产者不存在。

## ⑤ 这行字有没有测试钉着 —— **没有**

- 尺：`窗口|Window` 扫 `internal/**`＋`cmd/**` 的 `*_test.go` ＝ **3484** 行原始命中；其中 `窗口` **27** 行。
- ★**卡面那行 `窗口 %.1fs`：0 枚**。没有任何用例从控制台输出里读 `窗口`；跨面卡片尺 `cmd/wisp/panel_pump_test.go:444` 的正则只截到 `规则 …` 就停，**故意不覆盖** `窗口` 那行；`config_reload_223_test.go:435`、`approval_always_201_test.go:201` 只认 `[确认 L2 …]`／`卡片编号：` 前缀。
- `Prompt.Window` 两枚、**都是 L1**（`window_test.go:30-31` `==3s`；`wiring_test.go:133-134` `>=2s`）；`Prompt.Deadline` 两枚、都是 L2 `==300s`（`batch_test.go:105-106`、`queue_test.go:48-49`）。
- ★**L2 的 `Prompt.Window==0` 没有任何用例钉着** ⇒ 这条缺陷**无牙**：改甲/乙′/丙**不会碰红任何现有用例**（四枚变体 overlay 全 PASS，见 ⑦）。反过来说：AC#2 要的「今天响、不修不响」用例在仓里没有现货，必须新写。

## ⑥ 正控：真实 `win` 传进去 ⇒ 卡印的秒数 == `Queue().Timeout()`

- **甲 overlay**（仓外副本把 `gate.go:534` 的 `0` 换成 `g.q.Timeout()`）：`CARD_LINE="  窗口 300.0s"`、`QUEUE_TIMEOUT=5m0s` ⇒ `PRINTED_SECONDS=300.0 QUEUE_TIMEOUT_SECONDS=300.0 EQUAL=true`；同一次运行 L1 卡仍 `"  窗口 2.0s"`（`L1_ANSWER=timeout elapsed=2.001s`）。
- **乙／乙′／丙 overlay** 同 `EQUAL=true`（L2 均 `300.0s`）。
- ★**带底 31s 的甲 overlay**：卡印 `窗口 31.0s`，而同一次**实测等了** `elapsed=31.001s` ⇒ 印出的秒数**等于真实等待**（不只是等于 `Queue().Timeout()`）。HEAD 同一发：印 `0.0s`、等 `31.005s`。
- **反控（HEAD）**：`PRINTED_SECONDS=0.0 EQUAL=false`。
⇒ ★尺读得出「对的值」：同一台 printer，`win` 带真值就印 `300.0s`，且**逐位等于** `Queue().Timeout()`；印出的数不是另一套换算，就是 `p.Window.Seconds()`。

## ⑦ 门禁 rc

- `go build ./...` ⇒ **rc=0**；`go vet ./cmd/wisp ./internal/agent/approval` ⇒ **rc=0**；`d22scan`（`tools/d22scan` 自成 module，`-root ../..`）⇒ **rc=0**（`examined 291 production Go files under internal/ and cmd/`，`clean`）。
- `go test -overlay <变体> ./cmd/wisp -run 'Test376a1' -count=1`（`PATH="$PWD/third_party/sherpa-onnx:$PATH"`、`GOFLAGS=`、未 unset `TMPDIR/TEMP/TMP`）：`head`／`jia`／`yi`／`yip`／`bing`（首轮 overlay）与 `pin_head`／`pin_jia`／`pin_yi`／`pin_yip`／`pin_bing`（**钉 HEAD 字节**复跑）各 **rc=0**；`floor31` 与 `pin_floor31` 各 **rc=0**（31.0–31.2s）。输出：`D:\tmp\wisp\20261024\376a1\*-run.txt` 与 `pin-*-run.txt`。
- 格式尺：本笔 `.go` 零字节 ⇒ `gofumpt` 无事可判（⛔ 用 `gofumpt -l .` 判 HEAD 的 CRLF：会扫 `.scratch` 里故意写坏的探针）。两条被改的仓外副本与 HEAD blob 的 diff 用 `git diff --no-index --numstat` 量的（1/1、1/1、5/1），⛔ 不是 `grep '^[+-][^+-]'`。

## ⑧ 提交与 `.go` 计数

- 本件提交：`git log -1 --format=%h -- docs/evidence/s1/376-a1-census.md` 读回（⛔ amend，故不在正文写死哈希）；读数基线 `8ced487bb`。
- **本笔 `.go` 改动＝0 枚**（`git show --name-only --format= <本件提交>` 应只有本文件）；6 枚探针用例全在 `D:\tmp\wisp\20261024\376a1\`，经 `-overlay` 注入。
