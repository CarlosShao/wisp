# 223 — **`CheckAndReload` 在生产里一次都没被调用过**：配置改了不重读、D33 那条"放宽必须 L2 复确认"接的是一根没通电的线

- Status: **待派，且它是票 219（三枚答复按钮）与票 222 之后第一批的前置之一**——⚠ **不修它，任何"放宽 `[fs]` 要重新确认"的判据都只能在测试里绿**。
- 来源：非实现者设计复核 `219-v0`（`docs/evidence/s1/219-approval-reply-design-adversarial.md`，27,833 字节）第 5 条＋编排者自己复跑。

## 现量（起手逐条复算，别信这里的行号）

| 事实 | 读数 | 尺 |
|---|---|---|
| 那枚轮询函数存在 | `internal/config/manager.go:113` 一带 `CheckAndReload`，注释自陈是 SPEC-03 §4.3 的"幂等轮询（watchdog 1s tick）" | 现读 |
| ⛔ **非测试调用者只有一枚，且在旁支程序里** | **`cmd/balldebug/main.go:243`**：`bridge.Refresh = func() error { _, err := mgr.CheckAndReload(); return err }` | `grep -rn 'CheckAndReload()' --include='*.go' . \| grep -v _test` ⇒ 只剩这一行＋`internal/ball/hotkey_reload.go:21` 的**注释** |
| 生产装配不接它 | `cmd/wisp/run.go` 只在**装配时**建 canonicalizer（`grep -n 'Paths\|canonicaliz' cmd/wisp/run.go` 现跑），**没有任何 tick 调 `CheckAndReload`** | 现读 |
| 冻结要求 | `docs/PLAN.md:1644-1645` 逐字「**（新增）D33 配置提权**：热加载放宽 `[risk]`/`[fs]`/`[net]`/`[plugins]` → **必须触发 L2 级重新确认，不得静默生效**」 | `sed -n '1643,1646p'` |

## 后果（为什么这不是"少一个便利功能"）

1. **D36 的三档生效级别（立即／重启后／下次会话）今天没有任何一条走"立即"那档的实现路径**——没人轮询，就没有"热加载"。
2. **D33 的安全轨接不上**：想"放宽可读写范围时强制再确认一次"，接的是一根不转的线；于是"长期允许"这一支今天**要么做不到、要么做出来就是静默生效**（后者直接违反上面那句逐字要求）。
3. **owner 第一次用就会得出"这是个 stub"的结论**：改了 `config.toml` 里的东西，产品不重启不生效，而且**没有任何一句告诉他这件事**。

## 判据（每格都要现跑读数）

- [ ] **AC#1 生产里真有人在轮询**：`CheckAndReload`（或等价的显式刷新入口）在 `wisp run` 的装配路径上有**生产调用者**（现跑点数：改前非测试命中＝`cmd/balldebug` 那一枚、改后 `cmd/wisp` 里 ≥1 枚），并具名写清触发形状（watchdog tick／文件通知／显式命令）与**为什么不是"用墙钟时间差实现超时"**（`AGENTS.md` §1.2 硬禁）。
- [ ] **AC#2 三档生效级别各有读数**：立即生效的那一档要能观测到"改了→不重启→行为变了"；重启后生效的那一档要能观测到"改了→不重启→**明确告知未生效**"。⛔ **不许用"静默不生效"充当第二档**。
- [ ] **AC#3 放宽必带 L2 复确认（D33 正控）**：种一条把 `[fs]` 放宽的改动 ⇒ **必须产生一张 L2 卡**；把那次复确认拿掉 ⇒ 判据**必须红**。⚠ 收紧（变严）那一向**不许**也弹卡（否则把安全轨变成骚扰）。
- [ ] **AC#4 不生效与读不到是两句话**（同票 216 的纪律）：配置文件缺失／语法错／权限不够／热加载被禁用，**四种各一句**，不许合成一句"配置未生效"。
- [ ] **AC#5 整包终态读数**：`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp ./internal/...` 到终态＋逐名比红名集合（历史在册的 `internal/panel` 那几枚红**不算本票新增**，也不许顺手修）。

## 禁区

不动 `PLAN.md`／`docs/specs/**`／`thresholds.go`／golden／`allowlist.txt`；不新增 C17 方法名；不新增导出名（要新增先落 `A##`）；`internal/panel/l2_grant_boundary_test.go`／`tokens_fourway_test.go`／`internal/perm/ticket90_persist_test.go` 三枚冻结件一字不动；`frontend/**`／`design/**` 零写零转述；⚠ 与票 201／222 同撞 `cmd/wisp` ⇒ **串行**。
