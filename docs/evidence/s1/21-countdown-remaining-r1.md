# 票 21 段 2 落地腿 21-r1（纯 Go 面）：L1 倒计时的「活读数」＋派单前提复核

**时间：** 2026-11-05 11:44–12:0x +08 · **起点 HEAD：** `f4ad34ee1`（派单锚点）
**施工中树动了：** 另一腿 `8a75004ce`（29-r2，只碰 `internal/memory/**`、`internal/tools/spawn_roster_census_229_test.go`、`docs/evidence/s1/29r2-*`）在我读码期间落库；本腿 diff 的基线因此读作 `8a75004ce`，两者无交集。
**写面（本腿实际）：** `internal/agent/approval/{gate.go,ui.go}`＋新案 `ticket21_countdown_remaining_test.go`＋本文件。
**未碰：** `cmd/wisp/**`、`internal/tools/**`、`internal/panel/**`、`internal/ball/**`、`internal/statemachine/**`、`frontend/**`、`design/**`、`docs/PLAN.md`、`docs/specs/**`、`thresholds.go`、golden、票面（含 Status／勾）。

## 1 派单前提复核（两项，逐条给读数）

派单一句话是「剩下纯 Go 的活＝拆掉 loop 的预闸（pre-gate）＋把 Gate 装起来」。逐条复核，**两条在 HEAD 上都不是"剩下的活"**：

### 前提①「拆掉 loop 的预闸」— **STALE（早已改成条件放行，且剩下的那半被写面外的用例钉着）**

- 现量（非测试码）：`internal/agent/loop.go:919-924`
  ```
  case RiskL1, RiskL2:
          if l.opt.AdmitTask == nil {
                  j.decide(ctx, rowID, DecisionReject)
                  return false, false, "该操作属于 " + risk + " 级，审批通道尚未接入（ticket 21），已拒绝执行"
          }
          return true, true, ""
  ```
  ⇒ 有审批层（`Options.AdmitTask` 已接）时 **L1/L2 直接放行给桥**，loop 不再拦；只有**没接审批层的宿主**才拒。这正是票 12（`cbdea7c`）做的"删先拒＋引入 AdmitTask"，票据 `12-cli-text-path-s1-gate-done.md:163` 逐字记着。
- 生产装配确实接了：`cmd/wisp/run.go:1589` `AdmitTask: rt.admitTask`；`rt.admitTask`（`run.go:1305-1319`）第一件事就是 `rt.gate.AdmitTextTask(taskID)`（`run.go:1312`）。
  ⇒ **生产里这条 `if` 分支不可达**（`approval/doc.go:35-37` 已如实这么写）。
- 机器读数（真 loop＋真桥＋真门）：`go test -count=1 -run "TestLoopPassesDeclaredL1WriteThroughTheGate|TestLoopStillRefusesL1WhenNoGateIsRegistered" -v ./internal/tools/`
  ```
  --- PASS: TestLoopPassesDeclaredL1WriteThroughTheGate (3.27s)     ← 声明 L1 的 fs.write 今天真的穿过窗口落盘
  --- PASS: TestLoopStillRefusesL1WhenNoGateIsRegistered (0.19s)
  ok   github.com/CarlosShao/wisp/internal/tools  3.602s
  ```
  第一条 3.27s ≈ 2–3s 阻断窗 ⇒ "预闸挡死 fs.write" 这个病在 HEAD 上**已经不存在**。
- 第二条恰恰是"拆不得"的理由：同一文件 `internal/tools/loop_approval_test.go:231-258` 的注释与断言把"无审批层的宿主必须拒且文件不出现"当作**故意钉住的失败安全性质**（原话"If the early-reject had simply been deleted, the second case would fail … That is the failure mode the dispatch called worse than shipping nothing"）。
  ⇒ 物理上拆掉 `loop.go:919-924` 的 `if`＝要改 `internal/tools` 那枚用例（**不在本腿写面**），并**撤掉一条已交付的失败安全网**。**判 STOP（写面冲突＋需裁定），本腿未动 `loop.go` 一字。**
- 顺带登记（不是本腿的活）：那句拒因文案 `审批通道尚未接入（ticket 21）` 在"有门但没准入"时也不会出现（有门就放行），所以它不是用户可见的过期字串；真正的过期是 `approval/doc.go:32-37` 的 "Still owed" 清单（第 1 条 Native/Panel 已经落地，见下）——本腿**不改**它，因为改它需要同时改票面 Status，而票面不在写面。

### 前提②「把 Gate 装起来」— **DONE（生产两处装配点都在，且段 2 的出口已有生产调用者）**

- 装配点（非测试码）：`cmd/wisp/run.go:820` 与 `cmd/wisp/resident_approval_windows.go:401`，各 `approval.New(approval.Options{…})` 一处；门交给桥（`tools.Options.Gate`，见 `cmd/wisp/run.go` 的桥装配段），准入交给 loop（`AdmitTask` 同上）。`approval/doc.go:26-31` 也逐条记着这条已落地的接线。
- 段 2 的「Native/Panel 交出去」也已不是零调用者：`cmd/wisp/approval_reply.go`（`DecideFromNative`／`DecideFromPanel`／`Native()`／`Panel()`／`Veto` 的生产路由）。
  ⇒ "assemble the Gate" 没有剩下的纯 Go 活；剩下的段 2 项是**卡片渲染/球视觉/状态机驱动**，全在 `cmd/wisp/**`＋`internal/ball/**`＋`internal/panel/**`（本腿写面外）。

### 复核结论

派单描述的"剩余纯 Go 活"在 HEAD 上是 **0 项**（一项已落地，一项被写面外的钉住用例否决）。按 STEP0 的"blocked／写面冲突即停"，这两项**均判 STOP，不施工**。
本腿因此改做**票 21 段 2 里第一件写面内、无人钉住、且真的是"今天没人做"的格：AC#1 的倒计时数据面**（下节）。这不是派单原话的活，**具名登记**在此，判语归非实现者。

## 2 本腿施工：L1 倒计时的 `Event.Remaining` 从"整窗常数"改成"活读数"

**病（改前现量，不是推理）：** L1 路由每次发 `EventWarning` 都写 `Remaining: g.window`（改前 `gate.go:387`、`:518`）——**窗口已过 2.5s 仍说"还剩 3s"**。`Event.Remaining`（`ui.go:71-76`）是球条/卡片画倒计时环要用的那个数；照它画的环会**满格不动然后一跳到底**。`internal/panel/pump.go:50-54` 把这件事具名登记成"`approval.remainingMs` 无源，`Remaining` 发的都是静态窗长"——那枚"无源"的**生产者一半就在本包**。

**改法（小而直）：**
- `gate.go:116-125` 给内部 `window` 加 `start`（构造时用层内**同一枚** `Clock.Now()` 打点，只读、不加锁）；**到期权威仍是那条 `clock.After` 通道**，注释写明"start 只供显示，不能结束/延长/回答窗口"。
- `gate.go:461-487` 新 `remainingIn(w)`：`left = window − (Now − start)`，**下夹 0、上夹整窗**（越界读数不许报负数/超窗）。它**只报数**，没有任何缩短/延长/放行的路。
- 两处 `EventWarning` 改用 `g.remainingIn(w)`：`gate.go:396`（窗口 select 里的未加载通道否决）、`gate.go:558`（`Veto` 里未加载通道）。控制台/球两处 `Update` 只打印 `Text`（`cmd/wisp/run.go:2267-2268`、`resident_approval_windows.go:976-985`）⇒ **`Remaining` 的变化对现有两个界面零行为改变**（也就没有写面外回归）。
- `ui.go:64-72` 把 `EventTick` 的"声明的、没人发"写成实话（并说明为什么不在这里补发：控制台 `Update` 逐帧打印，tick 会往操作员日志里刷行；平滑环的节拍归宿主）；`ui.go:77-84` 把 `Event.Remaining` 的语义钉成"活读数、绝不等于配置窗长、0 表示不是倒计时"。

**RED-before（改前真红，不是编译红）：**
```
$ go test -count=1 -run TestTicket21UnavailableChannelWarningCarriesLiveRemaining ./internal/agent/approval/
--- FAIL: TestTicket21UnavailableChannelWarningCarriesLiveRemaining (0.00s)
    ticket21_countdown_remaining_test.go:50: 提示事件的 Remaining = 3s, want 500ms (窗口已过 2.5s): 报整窗 3s 的倒计时环会骗人
FAIL
```
**GREEN-after：**
```
$ go test -count=1 -run "TestTicket21" ./internal/agent/approval/     → ok
$ go test -count=1 -timeout 200s ./internal/agent/approval/          → ok 0.825s（89 枚顶层用例，含既有 87 枚）
```

**牙（两向变异，逐枚现跑，还原后 sha 逐位相同）：**
| 变异 | 改动 | 读数 |
|---|---|---|
| M1 退回旧形 | `left := g.window` | **红**：`Remaining = 3s, want 500ms`（`ticket21…:50`）；起点那枚照绿（它要的正是 3s） |
| M2 方向弄反 | `left := Now−start`（已过时间当剩余） | **红两枚**：`:50 → 2.5s, want 500ms`；`:84 → 0s, want 3s` |
| 还原 | `cp` 备份回写 | `sha256 ceafff88…0192` 与改后基线**逐位相同**，包复绿 `ok 0.850s` |

**第二枚用例是边界另一半：** `TestTicket21WarningAtWindowStartReportsTheWholeBlock`（`ticket21_countdown_remaining_test.go:68`）钉"窗口刚开＝整窗 3s"，否则 M2 那种"已过当剩余"或"一律夹 0"的读数会单边过关。

## 3 门禁（逐条 rc 现读）

| 门 | 命令 | 读数 |
|---|---|---|
| build | `go build ./internal/agent/...` | `rc=0`（空输出） |
| vet | `go vet ./internal/agent/approval/` | `rc=0` |
| test（本包） | `go test -count=1 -timeout 200s ./internal/agent/approval/` | `ok … 0.825s` |
| test（面内整棵） | `go test -count=1 -timeout 200s ./internal/agent/...` | `rc=0`；`ok internal/agent 7.602s`＋`ok internal/agent/approval 0.825s`＋`scheduler [no test files]` |
| gofumpt | `gofumpt -l <3 枚改件>`（真身 `/d/work/base/gopath/bin/gofumpt.exe`） | **空输出，rc=0** |
| d22scan | `sh scripts/d22scan.sh` | `rc=0`，末行 `clean - no D22 ban violations`（射程：`internal/` 880 枚 Go 含注释与 `_test.go`；`cmd/` 205；`design/` 39；`frontend/` 116） |

**⚠ 一枚与本腿无关的抖红，具名留档：** 施工前（11:51，仓里另一批发正在跑）`./internal/agent/` 整包红在 `Test31DisarmStopsTheWait`（`reminder_timer_test.go:370 取消后仍发了 1 条通知`），**单独复跑两次同样红**；12:03 同一 HEAD 整棵复跑转绿（`ok 7.602s`, rc=0）。⇒ 该用例是**负载敏感**（票 31 的提醒计时器），**不是本腿引入**（本腿只碰 `approval` 包，零 `reminder` 改动）。两读数都记在这里，不覆盖。

## 4 仍欠 / 不在本腿写面（具名移交，⛔ 不顺手改）

1. **`internal/panel/pump.go:50-54` 那句话现在过期了**：它逐字写"`approval.EventTick` 全树无人发、`Remaining` 发的都是静态窗长"——**后半句本腿已改真**，前半句仍真。改它属 `internal/panel/**`（写面外）：请面板腿把该句收窄为"tick 仍无生产者；L1 的活读数走 `EventWarning.Remaining`"，并顺带决定 `approval.remainingMs` 要不要由这条读数填（Q-51 的第五键不在本腿射程）。
2. **A20（D45-1 批量聚合＝死码）仍在**，且账上已把真修法改归 S3、票 21 段 2 只欠"删除或显式隔离"。本腿复核了它的可达性（非测试码 `PathParams` 全表，最大 2 枚：`internal/tools/fs_write.go:740` 的 `from/to`；`BatchMinOps=3`，`Aggregate` 唯一调用点 `gate.go:338`）⇒ **任何在产工具单次 ≤2 路径，`Aggregate` 恒 nil**（与 A20 原判一致）。**本腿不做**：它既不是票面 AC 的任一句，删除又会让 `batch_test.go` 的 5 枚断言一起消失（听起来像"为绿删牙"），**要动请编排者先裁定"删 or 隔离"**（台账已把这件事挂在段 2 名下，但没裁定形式）。
3. **票面 5 个 AC 框、`Status:`、"segment 2" 归属：一字未动**（票面不在本腿写面）。本文件只是读数件；**勾选与终判归非实现者**。
4. **本腿⛔ 声称的**：⛔ 声称"L1 倒计时环在真机上已可见"（球的渲染属段 2 外发腿；本腿只把**数据源**改真）；⛔ 声称任何 AC 已过。

## 5 一句话判语（供验收方独立复核）

派单的两项"剩余纯 Go 活"在 HEAD 上是 **0 项**（`loop` 预闸早已是条件放行＋生产不可达；`Gate` 已在 `run.go:820`／`resident_approval_windows.go:401` 装好），拆闸更要撞写面外的钉子用例；本腿改做了写面内唯一还没人做的一格——**把 L1 倒计时的 `Remaining` 从整窗常数改成同一枚 Clock 上的活读数**，两向变异可复现，门禁 rc 全 0。
