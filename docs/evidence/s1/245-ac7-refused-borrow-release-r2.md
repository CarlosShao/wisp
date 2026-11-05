# 245 · AC#7 落地自述（借被拒那一支的记账归位）— 实现腿 `245-r2`

> ⚠ **这份不是裁决表**。按 `AGENTS.md §1.5`／`SPEC-12 §4.3`，`docs/evidence/s1/` 的判语必须出自**非实现者**；
> 本文件与 `245-...-r1.md` 同性质，是〔实现方自述＋台件〕，**一枚框都不翻**。裁决归下一程非实现者腿。
>
> 派单读数（逐字）：`Status: 已立，未派，排在队列最前`；派单只给 **AC#7**，AC#8／AC#9 需真机另归。
> 本腿**只做 AC#7**，未碰 AC#4／#6／#8／#9。

- 起手 HEAD：`2f7f9de1632d0bca88e7a4fa9cf30cd96010ce2f`（＝派单里的 `2f7f9de16`）
- 写面：`internal/ball/ball_windows.go` ＋ 新测试 `internal/ball/hotkey_refused_borrow_release_245_test.go` ＋ 本文件。
  未动 `cmd/wisp/**`、`PLAN.md`、`docs/specs/**`、`internal/observe/thresholds.go`、golden、`frontend/**`、`design/**`，票面本身也一字未改（不属本腿写面）。

## 1. 三核前提（每一条都是我自己现跑，不转抄票面）

### (a) 票面点名的字符串／行还在不在

| 票面点名的形状 | HEAD 现读 | 结论 |
|---|---|---|
| `takeEsc` 失败分支**只写收据、不设 `escTakenOver`** | `internal/ball/ball_windows.go:892-904`：`if err := takeEscBorrow(...); err != nil { … withCancel(cancelFailedLineFor(borrow, err)); return }`，其后才是 `b.escTakenOver = true` | **在**（票面写 `:878-899`，见 (b)） |
| `ReleaseEscAfterSession` **首行早退** | `:911-933`（改前 `:913-915`）：`if !b.escTakenOver { return }` | **在** |
| 报告因此永久停在 `HotkeyError` | `HotkeyStatus` 的失败档 → `Problems()` `:381-404`（`"hotkey %s = %q was not registered: %v"` 在 `:400-401`）；`AllLive()` `:358` 对 `HotkeyError` 返回 false | **在** |
| `escTakenOver` 字段 | `:125` | **在** |
| 报告 cancel 行由 `withCancel` 唯一替换 | `hotkey_windows.go:654-665` | **在** |

⇒ 缺陷形状与票面逐字一致：**Win32 什么都没绑，报告却说"注册失败"，且这句话没有归位路径**。

### (b) 票面引用的 `file:line` 对不对

- **`ball_windows.go:878-899` 已过期 +14 行**：那一段现在是 `:892-904`（`TakeEscForCancel` 整体 `:881-903`）。漂移来源＝票 260 在 `TakeEscForCancel` 上方追加的 `resolveCancelBorrow`/形ⓐ 注释块，**不是**本票前有人改过这一支。
- `ReleaseEscAfterSession` 早退在改前是 `:913-915`（票面未给行号）。
- 票面 AC#7 的判据句（"报告必须说稳态三枚、cancel 未绑"）与 HEAD 上 `HotkeyStandby` 的文案 `hotkey_windows.go:281-285` 对得上。

### (c) 有没有别的腿的在飞改动

- `git status --porcelain internal/ball cmd/wisp` ⇒ 只有 `?? cmd/wisp/zz_probe_temp_test.go.txt`（一个 `.txt` 探针文件，不参与编译）与我自己新增的测试文件；**`internal/ball` 无他人未提交改动**。
- 起手基线：`go test ./internal/ball -count=1` ⇒ `ok ... 0.241s`（**全绿**，票 80 记的那枚 `TestC21TableColourRowsMatchTokensCSS` 现已不红）。

### 写面／Blocked-by（从票面自己的文字读）

- 票面「禁区」只禁 `PLAN.md`／`docs/specs/**`／`thresholds.go`／golden／三枚冻结件／`frontend/**`／`design/**`／新增词面型仪器；AC#7 的落点 `internal/ball` **合法**。
- AC#7 **无 Blocked-by**（掉在 `internal/ball`，不需要真机、不需要 `cmd/wisp`、不需要配置通路）。

## 2. AC 可达性（本腿只做 AC#7，其余逐条说为什么不碰）

| 框 | 本腿 | 理由（票面原文为准） |
|---|---|---|
| **AC#7** | **做了** | 见下 §3–§6 |
| AC#4 | 不做 | 票面自陈「阻塞于票 228」；本腿不造第二条配置通路 |
| AC#6① | 不做 | 落点是 `cmd/wisp/main.go`，**出本腿写面**（且票面把它判归 228） |
| AC#6② | 不做 | 那四枚"按枚数锚"的钉不在本腿 AC；改钉会牵动票 245 既有判据面，属裁决/编排者口径 |
| AC#8／AC#9 | 不做 | 需真机／需第二条真 GUI 进程台件，派单明说另归 |

## 3. 实现形状（一行话＋一张状态表）

新增一枚 UI 线程私有字段 `escBorrowRefused bool`（`ball_windows.go:126`）：**"这次会话里借过、且 Win32 拒了"**。
`ReleaseEscAfterSession` 从"只看 `escTakenOver`"改成"**两种结束方式都归位**"：借到过就 `releaseEsc` 丢键，只有被拒就只把 cancel 行写回 idle 真值。
`RebindHotkeys` 同批清这枚旗（新的一轮注册本来就是 idle 名册）。

| 时刻 | `escTakenOver` | `escBorrowRefused` | cancel 行 | `Live()` | `AllLive()` | `Problems()` |
|---|---|---|---|---|---|---|
| 稳态（boot／rebind） | false | false | standby | 3 | true | 无 |
| 借成功（Confirming 中） | **true** | false | live | 4 | true | 无 |
| **借被拒（卡片挂着）** | false | **true** | **HotkeyError**（票 260 AC#2 的收据，**保留**） | 3 | false | 1 条"was not registered" |
| **借被拒＋会话结束（本腿新增）** | false | false | **standby** | 3 | **true** | **无** |
| 借成功＋会话结束 | false | false | standby | 3 | true | 无 |

⇒ 判据两侧都收：**卡片挂着时照旧喊出来**（零症状那一形没被治好成另一种零症状），**会话结束后报告必须说桌面真值**。

## 4. 红前／绿后（逐字）

新测试三发（默认层，`//go:build windows`，**无 `t.Skip`**，不依赖真桌面；借用的是 `260-r5` 已量过的"无窗口 ui-sta 门"＋非法 `HWND` 让 `RegisterHotKey` 真拒绝，cancel 绑定刻意取 `Ctrl+Alt+F8` 这枚**带修饰键**的组合，全程不碰裸 Esc）：

```
# RED（改 ball_windows.go 之前，HEAD 原样）
--- FAIL: TestRefusedBorrowReleaseReturnsTheIdleRoster245
    after the session ended the cancel line reads "registration failed" (err=Invalid window handle.),
      want "not bound while idle (cancel is borrowed only during Confirming)"
    the idle cancel line carries a stale error after the session ended: Invalid window handle.
    Problems() after the session ended = ["hotkey cancel = \"Ctrl+Alt+F8\" was not registered: Invalid window handle."], want none
    AllLive() = false after a refused borrow was returned
      (lines: summon="live" mute="live" cancel="registration failed" panel="live")
--- PASS: TestRefusedBorrowStaysNamedWhileTheCardWaits245   （正控①，改前即绿）
--- PASS: TestHeldBorrowReleaseStillReturnsStandby245        （正控②，改前即绿）
FAIL	github.com/CarlosShao/wisp/internal/ball	0.044s
```

```
# GREEN（改后）
--- PASS: TestRefusedBorrowReleaseReturnsTheIdleRoster245 (0.01s)
    245 AC#7: refused line before release = {Name:cancel ID:3 Binding:Ctrl+Alt+F8 Status:registration failed
      Acc:{Mods:16387 VK:119} Err:Invalid window handle. Note:}
      / ["hotkey cancel = \"Ctrl+Alt+F8\" was not registered: Invalid window handle."];
      after release = {Name:cancel ID:3 Binding:Ctrl+Alt+F8 Status:not bound while idle
      (cancel is borrowed only during Confirming) Acc:{Mods:0 VK:0} Err:<nil> Note:} / problems=[] live=map[1:{16387 81} 2:{16387 77} 4:{16387 80}]
--- PASS: TestRefusedBorrowStaysNamedWhileTheCardWaits245 (0.00s)
--- PASS: TestHeldBorrowReleaseStillReturnsStandby245 (0.01s)
ok  	github.com/CarlosShao/wisp/internal/ball	0.050s
```

**三发各自钉的是什么**

1. `TestRefusedBorrowReleaseReturnsTheIdleRoster245`（AC#7 主尺）：借被拒 → `ReleaseEscAfterSession()` → cancel 行 `HotkeyStandby`／`Err=nil`／`Binding` 仍是配置值；`Live()` 与 idle 种子 `DeepEqual`（恰三枚、无 `hkCancel`）；`AllLive()` true；`Problems()` 空；`registeredHotkeys` 回三枚；第二次 `Release` 是 no-op（幂等）。**改前红、改后绿**。
2. `TestRefusedBorrowStaysNamedWhileTheCardWaits245`（正控①，**不许把收据吞掉**）：同一 rig，**不**调 `Release`，要求恰好 1 条问题行、`HotkeyError`、行里带 `cancel`／配置键名／`was not registered`、`AllLive()` false。它同时也是 M1 突变的红点（见 §5）。
3. `TestHeldBorrowReleaseStillReturnsStandby245`（正控②，**老路径不许被我改坏**）：借成功形态（`escTakenOver=true`＋live cancel 行）→ 归还 → standby／三枚／`Problems()` 空。

## 5. 突变（两发，均用 `-overlay`，未改工作树一个字节）

| 突变 | 改法 | 读数 |
|---|---|---|
| **M2** | `ReleaseEscAfterSession` 退回票面原形：`if !b.escTakenOver && !b.escBorrowRefused` → `if !b.escTakenOver` | **主尺红**，红句与 §4 红前逐字相同（4 条） |
| **M1** | 失败分支改成直接写 idle 行：`withCancel(cancelFailedLineFor(borrow, err))` → `withCancel(cancelIdleLine(b.cancelBinding))` | **两发红**，红句逐字：`the borrow was not refused as this rig requires: cancel line = {… Status:not bound while idle (cancel is borrowed only during Confirming) …}` ⇒ 想靠"把失败悄悄写成 standby"过主尺的路被正控①堵住 |

两发突变都是 overlay 临时文件，跑完 `git status --porcelain internal/ball` 仍只有我那两枚改动（M2/M1 未落地）。

## 6. 门禁（现跑，本机）

| 门 | 命令 | 读数 |
|---|---|---|
| 编译 | `go build ./internal/ball ./cmd/wisp` | rc=0 |
| 包测试 | `go test ./internal/ball -count=1 -timeout 200s` | `ok github.com/CarlosShao/wisp/internal/ball 0.166s` |
| 新用例单跑 | `go test ./internal/ball -run 245 -count=1 -v` | 3/3 PASS |
| 格式 | `gofumpt -l internal/ball/ball_windows.go internal/ball/hotkey_refused_borrow_release_245_test.go` | 空输出，rc=0 |
| vet | `go vet ./internal/ball`；`go vet -tags winlive ./internal/ball`；`go vet -tags winlive ./cmd/wisp` | 三发 rc=0 |
| D22 闸 | `sh scripts/d22scan.sh` | `rc=0`；`PASS=35 FAIL=0 SKIP=0`；`d22scan: clean - no D22 ban violations`（`internal/` 353 生产件＋869 Go 件含注释/测试） |

## 7. 残余与归口（具名，不藏）

1. **本腿只解"会话结束后归位"，不解"借被拒那张卡本身没有 Esc"**：卡片挂着时它仍是一条 `HotkeyError`＋一句"was not registered"，这是票 260 AC#2 要的收据，**不是**本票残留。真要让用户有一条能用的否决键（换键名／第二条通道）仍是票 245 文末五字段残余登记里的那笔账。
2. **常驻腿的调用时序我没改**（`cmd/wisp` 出写面）：`resident_approval_windows.go:924` 借、`:1002` 的 `settleOrb()` 只在"没有卡片在等人"时归还（`:995-997`），所以本腿的归位**只会发生在卡片撤下之后**；`cmd/balldebug/main.go` 的 `syncMachineToBall` 每次离开 `Confirming` 都会调一次归还，同样落在卡片之后。
3. **`RebindHotkeys` 顺手清旗**：新的一轮注册本来就是 idle 名册，若不清，下一次 `Release` 会去重写一行本来就已经是 idle 的报告（无害但记账不对）。
4. **未加 `DEFERRED(D-xx)` 代码标记**（与票 245 既定口径一致：标记须与 `SPEC-12 §5` 表 1:1 双向，本腿不动那张表）。
5. **裁决未做**：本文件是实现方自述；AC#7 的框该由非实现者腿按 §4 红/绿读数与 §5 突变独立复核后翻。

## 8. 改动清单

| 文件 | 性质 |
|---|---|
| `internal/ball/ball_windows.go` | ＋`escBorrowRefused` 字段；失败分支置位、成功分支清位；`ReleaseEscAfterSession` 两路归位；`RebindHotkeys` 清位（＋25 / −2） |
| `internal/ball/hotkey_refused_borrow_release_245_test.go` | 新增（默认层三发：AC#7 主尺＋两枚正控） |
| `docs/evidence/s1/245-ac7-refused-borrow-release-r2.md` | 本文件 |
