# 票 64 只读取证裁定 — 球缺陷（热键/交互）今日现量

- **日期**: 2026-11-07T13:5xZ（机主本地 2026-11-07 21:5x +08）
- **腿**: wave 64-evidence（只读取证，零产码；`internal/**`、`cmd/**` 一字节未改）
- **方法**: 全部结论出自现跑 grep/sed/git + 一次非 winlive 单测复跑；**未跑 winlive**（桌面独占面，本波 `internal/ball` 写面归票 62 腿，不占桌面）。
- **门禁**: 零 Go 改动 ⇒ gofumpt/d22scan **无对象可扫**（如实登记，非跳过）。

## 起手三核（现跑记录）

1. **git status --porcelain**：`internal/ball`、`cmd/wisp`、`internal/config` **无未提交跟踪改动**（仅他人未跟踪临时件，与本面无关）；
   `tasklist` 无 balldebug/wisp/go 进程 ⇒ 本件读数**无在飞件污染**。
2. **票面锚点现量**：多数已过期（下表逐条给现址），票面所述缺陷按今天的盘上形状逐条复裁。
3. **引文出处**：见各条「落点」列，全部为现跑行号。

## 逐条裁定

| # | 票面条目 | 裁定 | 今日事实（落点 = 现跑行号） |
|---|---|---|---|
| A1 | `Ball.RebindHotkeys`（票锚 `ball_windows.go:693`）全仓零调用者，改 `[hotkey]` 不重注册 | **已修** | def 现在 `internal/ball/ball_windows.go:862`。**生产调用者已存在**：`cmd/wisp/resident_ball_windows.go:444` 经 `ball.NewHotkeyReloader` 接桥，`bridge.Run(ctx, time.Second)` 1s 轮询、observe 注册表 spawn、停机 join；`cmd/balldebug/main.go:269` 同款。球窗本身由 `cmd/wisp/resident_ball_windows.go:396` 建（票面 R12「`cmd/wisp` 不建球」已失效）。重载源 = `hotkeyReloadSource296`（`cmd/wisp/resident_windows.go:49`）：每 tick 新读 `config.LoadFile` + `ApplyHotkeyDefaults` 合并空槽（票 296 守卫）。⚠ 形状与 R12 处方（`config.NewManager`+watchdog）不同：走每 tick 重读而非 Manager 事件——承诺等价（diff 前零 Win32 churn），票 64 无需动作。 |
| A1b | 注册失败两类语义（1409 占用 vs 未尝试）不分开 | **已修** | 四族枚举 `HotkeyDisabled/HotkeyUnparsable/HotkeyTaken/HotkeyError`（`hotkey_windows.go:245-253`），`Attempted()` 分界 `:296`，`Problems()` 每族独立话术。复跑 `TestRegisterAllSplitsFailureFamilies` ×2 **PASS**（ok 0.062s）。 |
| A1c | 默认唤起键 `Ctrl+Alt+W` 被占 → 改 `Ctrl+Alt+Q`；`Ctrl+Alt+Space` 备选写文档 | **已修** | `DefaultHotkeys()`（`hotkey_windows.go:75`）Summon=`Ctrl+Alt+Q`（注释引 R10/84 组合实扫）；钉 = `TestDefaultHotkeys`（`hotkey_test.go:63`，票锚 `:55` 漂移）断 Q、反钉 W/Space、AltSummonSpace 可解析。文档半 = `docs/HOTKEYS.md`（在盘，3082B：四键表 + 为什么不是 W + Space 为何不当默认）。复跑 ×2 **PASS**。 |
| A2① | 单击球 Sleeping→Listening（含隐形边框偏移 (8,31) 缺陷） | **已修** | `WS_POPUP` 落地（`ball_windows.go:208-220`，client==window，注释点名票 64 A2）；`HitTest`（`hit.go:50`）现 DPI 感知、可点圆 = max(orb/2+ClickTolerance, 0.75·orb)（R15#5 可见即命中），四角 HTTRANSPARENT 穿透。live 测 `TestLiveClickSummonsAndDragDoesNot`（`interaction_live_test.go:53`）在盘；2026-09-21 双跑 PASS×2 无 SKIP（票面 Progress log 有逐跑记录），本腿未重跑 winlive。⚠ 几何已随票 62/65/68 变形（`RingMarginPx` 8→22、`BallSizeDefaultPx=56`、overlay orb 窗叠加），交互面现形状见 A2④。 |
| A2② | Confirming 用 Esc/单击取消 | **已修** | `TestLiveConfirmingCancelAndEscReturned`（`interaction_live_test.go:133`）在盘，且被 245-r3（`9e84ae42e`）与 260-r2（`e3e19e8e0`）**加严过**：稳态 cancel 改 Standby、Esc 仅 Confirming 借走（借到的键吃 `[hotkey] cancel` 配置值）。 |
| A2③ | 不抢焦点（真焦点窗对照） | **已修** | 球窗 ex-style = `WS_EX_LAYERED\|TOPMOST\|TOOLWINDOW\|NOACTIVATE`（`ball_windows.go:207`），`WM_MOUSEACTIVATE` 回 `MA_NOACTIVATE`（`:610`）；live 测 `TestLiveNeverStealsFocus`（`:282`）在盘。 |
| A2④ | 透明区点击穿透（`WindowFromPoint` 能否命中） | **已修** | 命中链今天双层：球窗在可点圆外由 `WM_NCHITTEST` 回 `HTTRANSPARENT`（`hit.go` + `ball_windows.go` 消费点 :657 一带）；票 65 新增的 orb overlay 窗**纯穿透**——`WS_EX_TRANSPARENT` + wndproc 对 `WM_NCHITTEST` 恒回 `HTTRANSPARENT`（`orb_windows.go:92-96`），且注释留有实测：**非 layered 的 `WS_EX_TRANSPARENT` 窗不被 `WindowFromPoint` 跳过**（`:486-487`），故双层兜底是必要的而非冗余。live 测 `TestLiveTransparentCornerFallsThrough`（`interaction_live_test.go:351`）在盘。 |
| A3 | Sleeping 零定时器句柄实测（winlive、不许 skip） | **已修（测试在盘）** | `TestLiveSleepingZeroTimerHandles`（`hotkey_live_test.go:449`）在盘（数真 `WM_TIMER`，带探针反证）。2026-09-21 双跑 PASS×2、全文 `SKIP` 0 次（票面 Progress log）。CI 半边转 nightly（票面已记残留），本腿未重跑。 |
| A4 | 多显示器真拖到第二屏 | **仍存在（按票面规则保持未勾）** | 票面 AC 框仍 `- [ ]`（未代翻）。`enumMonitors()`（`monitors_windows.go:40`）仍是普通函数，**注入接缝仍没开**（票面预告过的前置条件依旧缺席）；本机仍单屏。硬件到位批次再开。 |
| A16 | `hotkey_windows.go` 按下标 zip（票锚 `:382-386`）配对错位风险 | **仍存在（已登记降级，非现行缺陷）** | zip 形状未消：`hkNames`（`hotkey_windows.go:44`）× `bindings[i]`（`:481-492`）。配对钉今天在 `hotkey_status_test.go:486/:502`（票锚 `:225/:241` 漂移）。按 2026-09-21 裁定：修法=具名单表消 zip，优先级挂票 39 第五热键同批——维持原挂账，本腿不改。 |

## 票面过期锚点清单（供下一位写腿，勿再按票面行号找）

| 票面所引 | 今日现量 |
|---|---|
| `RebindHotkeys` @ `ball_windows.go:693` | 现在在 `:862` |
| `len(live)==4`（`hotkey_live_test.go:101-109` 引文） | **已过期**：245-r1 起 cancel 稳态 Standby，测试现断 `len(live)==3`（`hotkey_live_test.go:234` 一带）；步骤 5 物理注入新键 R 唤起/旧键 Q 不响**原样保留** |
| `hotkey_test.go:55`（默认钉） | 现在在 `:63` |
| `hotkey_status_test.go:225/:241`（配对钉） | 现在在 `:486/:502` |
| `manager.go:203`（rest 表整体覆盖 `cur.Hotkey`） | 机制仍在，行移 `:306`（`{"hotkey", &cur.Hotkey, ...}`）；但生产重载已不走 Manager 事件，走每 tick `LoadFile`（`hotkeyReloadSource296`） |
| `cmd/wisp/run.go:189` / `providers.go:90`（一次性 LoadFile、不建球） | 两行内容已变；「不建球」整句失效（`resident_ball_windows.go:396` 建球、`:444` 接桥） |
| `balldebug/main.go:171` | 现 `:221`（ball.New）、`:269`（桥） |
| `SKIP-LOUD` 8 处 | 现 **9** 处（245/260 波新增） |
| `Cancel` 默认四键全注册 | 已被 245 改为「三键稳态 + Esc 仅 Confirming 借还」，`docs/HOTKEYS.md` 已按新形状书写 |

## 门禁与红绿归属

- 本腿写面仅 `docs/evidence/s1/64-ball-defects-verdict.md` + 票 64 Progress log 一行（append-only；AC 勾与 Status 零触碰）。
- 零 Go 改动 ⇒ gofumpt/d22scan **无对象可扫**。
- 复跑读数（非 winlive、count=2）：`TestDefaultHotkeys` / `TestRegisterAllSplitsFailureFamilies` / `TestDefaultHotkeysIdlePassHoldsNoEsc` 全 **PASS**，`ok internal/ball 0.062s`。
- 红绿归属：取证开始时 `internal/ball`、`cmd/wisp` 无跟踪面未提交改动、无 go/wisp/balldebug 进程 ⇒ 上述读数无他腿在飞件混淆。
- 未做（如实）：winlive 交互四项与零定时器**今天未重跑**（桌面独占 + 62 腿本波独占 `internal/ball` 写面）；其"已修"判定依据 = 盘上测试体 + 2026-09-21 逐跑记录，标注为"测试在盘、最近已证跑"，非本腿今日实测。
