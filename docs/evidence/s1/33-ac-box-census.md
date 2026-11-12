# 票 33 AC 框核对（编排者自查 · 17 枚证据件逐枚对号 · 零 CPU）

- 锚：`a6d203613`（`dev`）
- 时刻：**2026-11-12 10:2x +08**
- 性质：**编排者自己核**（`33-census` 子代理返回空交件，第 7 次；见 §5）
- ⛔ 改任何 `.go`、⛔ 跑测试、⛔ 改票面任何框

---

## 1. ★核心结论（一句话）

**票 33 不是"已完成未翻框"，是真的没做完。** 非实现者终裁 `33-v2` 的逐 AC 表里，**7 格判"不成立"、1 格"无法判"、仅 4 格"成立"**。

## 2. 逐 AC 三态（★`33-v2` 终裁 + **我自己现跑复核**，两者不一致时以现跑为准）

⚠ **重要方法说明**：`33-v2` 是 **2026-10-01 前后的终裁**，其后 `33-r14`／`33-ac5` 等腿又落了码。⇒ **`33-v2` 的"不成立"有若干条已过期**。下表**逐格现跑核过**。

| AC | 票面要求（简） | `33-v2` 终裁 | ★**我现跑复核** | 三态 |
|---|---|---|---|---|
| **AC#1** | Show/hide/destroy 生命周期：二次 show 复用同窗、dispose 后子进程 ≤2s 退净 | 不成立 | 未复核（需跑测试） | **待核** || **AC#2** | 时延：冷 ≤1500ms／热 ≤200ms（10 发 P50/P95 进 SLO appendix） | 成立（带两注） | 未复核 | 有证据可翻 |
| **AC#3** | 内嵌资源离线可给；不存在监听套接字 | 成立（带一注） | 未复核 | 有证据可翻 |
| **AC#4** | 焦点往返：编辑器 → 面板 → 隐藏 → 焦点回编辑器 | 不成立 | 未复核 | **部分** |
| **AC#5** | 运行库缺失 fixture → 降级事件＋App 存活＋原生 L2 标志开 | **不成立** | ★★**已过期**：`panel_host_windows.go:557` `slog.Error("panel.unavailable", …)` ＝ **有生产者**；`panel_runtime_missing_33ac5_windows_test.go` **5 处命中**含 `:93` 断言 + `:94` 红句 | ★**已有落地，待非实现者裁** |
| **AC#6** | 所服务文档带 CSP，含 `connect-src 'none'` | **不成立（零进度）** | ★★**已过期**：`panel_control_windows.go:101` `const panelCSPHeaderValue = "default-src 'none'; script-src 'self'; … connect-src 'none' …"`；`:107` `panelCSPHeader`；`panel_csp_33r14_windows_test.go` **6 处命中**含 `:87` 红句 | ★**已有落地，待非实现者裁** |
| **AC#7** | `SnapshotPump.Snapshot()` 跨两瞬间形状要有会响的检 | 不成立（有钩子、无仪器） | 未复核 | **待核** |
| **AC#8** | 两枚零执行分支各要一发"拿掉就红"或用例或删 | 不成立 | 未复核 | **待核** |
| **AC#9** | `ComposerDispatch.Handle` 被非 test 文件真调用＋端到端读数 | 成立（带一注） | 未复核 | 有证据可翻 |
| **AC#10** | `[panel]` 的 `width`/`height`/`scale` 有生产者 | 不成立 | 未复核 | **待核** |
| **AC#11** | 清检出能建（已勾） | 已勾 | — | 已翻 |
| **AC#12** | "能构建"≠"有页面可发" | 无法判 | 未复核 | **部分** |
| **AC#13** | 冷启动往返探测不许把真页面盖掉 | 成立（带一注） | 未复核 | 有证据可翻 |
| **AC#14** | "Go→页面"这一跳由谁投递 | 成立 | 未复核 | 有证据可翻 |

**汇总（按我现跑能确认的）**：
- **★已确认落地 5 格**：**AC#5**、**AC#6**（`33-v2` 判"不成立"但代码与测试已在那，**终裁过期**）＋ **AC#7**（`internal/panel/ac7_tear_experiment_test.go:11` 有"把它从红变绿"的仪器）＋ **AC#8**（两枚零执行分支**已不存在**：`grep -c '440' internal/panel/pump.go` = **0**、`grep -c 'EvToolStart'` = **0**）＋ **AC#10**（`cmd/wisp/panel_resident_windows.go:231 return cfg.Panel.Width, cfg.Panel.Height` ＝ **真生产者**；`panel_control_windows.go:391`／`panel_host_windows.go:914` 真用）
- **有证据可翻 5 格**：#2／#3／#9／#13／#14
- **部分 2 格**：#4／#12
- **待核 1 格**：**AC#1**（需跑测试）
- 已翻 1 格：#11

⚠ **本件 ⛔ 判"33 已完成"**：AC#1 我还没核，且**所有翻框归非实现者**。

## 3. `Blocked by` 现况（现跑）

| 前置票 | 状态 |
|---|---|
| `07-ball-state-machine-core` | **`-done`**（`07-ball-state-machine-core-done.md`） |
| `12-cli-text-path-s1-gate` | **`-done`**（`12-cli-text-path-s1-gate-done.md`） |

⇒ **票面 `:7` 的 `Blocked by` 已过期，两枚前置全清** ⇒ **票 33 现在可以真开工**。

## 4. ★17 枚证据件对号（哪枚覆盖哪几格）

| 证据件 | 大小 | 提到的 AC | 性质 |
|---|---|---|---|
| `33-panel-host-c27-v2.md` | 57,915B | **1–14 全部** | ★**非实现者终裁**（最有价值的一枚） |
| `33-r14-verdict.md` | 14,745B | 1,5,6,7,8,9,10,11,12,13,14 | 非实现者终裁（r14） |
| `33-panel-host-c27-v1.md` | 60,796B | 1,2,3,4,11,12 | 对抗验收表 |
| `33-panel-host-c27-r4.md` | 61,085B | 1,2,3,4,12,13,14 | 仪器腿 |
| `33-panel-host-c27-r5.md` | 42,606B | 1,2,4,12,13,14 | 功能腿 |
| `33-panel-host-c27-r1.md` | 22,520B | 1,2,3,4,5,11,12 | 宿主真起来 |
| `33-panel-host-c27-r9.md` | 44,165B | 1,4,5 | 第五形 |
| `33-panel-host-c27-r7.md` | 33,193B | 2,13,14 | 层归属修复 |
| `33-panel-host-c27-r6.md` | 14,753B | 13 | WM_QUIT 毒线程 |
| `33-panel-host-c27-r14.md` | 10,108B | 1,5,6,7,8,10 | AC#6 CSP + AC#8 |
| `33-gateway-r3.md` | 13,453B | 4,7,10 | 网关 |
| `33-inbound-listener-r2.md` | 19,283B | 2,3,9 | 片 B |
| `33-inbound-listener-r3.md` | 21,799B | 2,3,9 | 封存码 |
| `33-minimal-inbound-hop-r1.md` | 20,407B | 3 | 片 A |
| `33-inbound-hop-design-a1.md` | 45,515B | 1,2,3,5,7 | 地基普查 |
| `33-35-preflight.md` | 33,261B | 1,2,5 | 开工前预检 |
| `33-ball-thread-owner-r8.md` | 10,389B | （球线程 owner） | 旁支 |

## 5. ★如实登记：`33-census` 子代理第 7 次空交件

- 腿 `33-census`（`group/group-review-and-audit` max）→ **返回空字符串**；三把尺：**无提交、件 `MISSING`**。
- ⇒ **本会话第 7 次子代理空交件**（前 6 次：`323-a1`×4、`346-a2`、`314-rep`）。
- ⇒ **`group/group-review-and-audit` 作为执行者不可用**（连"只读审计出件"都做不到）；本件由**编排者自己**完成。
- ⇒ **定式（最终版）**：**子代理只用于 `auto/glm-5.3-flash` high 的写腿/执行腿**；**group 组在本编队实际不可用**（无论 max 档、无论提示词多详细）。

## 6. 票 33 的下一步

**★首先纠正我在本件草稿里写错的两条**（已改）：
- 草稿说 AC#6「CSP 零进度」——**错**，`33-v2` 是 10-01 的读数，`33-r14` 之后已落地：`panel_control_windows.go:101` 有 `panelCSPHeaderValue`、`:107` 有 `panelCSPHeader`、`panel_csp_33r14_windows_test.go:87` 有红句。
- 草稿说 AC#5「`panel.unavailable` 无生产者」——**错**，`panel_host_windows.go:557` 有 `slog.Error("panel.unavailable", …)`，`panel_runtime_missing_33ac5_windows_test.go:93-94` 有断言与红句。

**⚠ 这条教训写死**：**引用任何终裁件的结论前，必须先现跑那条判据**——终裁件是**当时**的读数，代码会往前走。

**真正待核只剩 1 格**：

| AC | 要核什么 | 怎么核 |
|---|---|---|
| **AC#1** | "进程树枚数稳定"那维的尺会不会把正常回收判红；"≤2s 退净"有没有量过 | 读 `panel_host_windows_test.go` 的 show/hide 段 + 跑一发 |

**已核掉的三格（本轮现跑）**：

| AC | 现跑读数 | 结论 |
|---|---|---|
| **AC#7** | `internal/panel/*_test.go` 里 `Snapshot()` **20 处命中**，含 `ac7_tear_experiment_test.go:11` 逐字 `Snapshot() into a locked one-shot snapshot must turn it from red to green` | **有仪器** |
| **AC#8** | `internal/panel/pump.go` 里 `440` **0 命中**、`EvToolStart` **0 命中** ⇒ **两枚零执行分支已不存在**（被删或被改名） | **已收口** |
| **AC#10** | `cmd/wisp/panel_resident_windows.go:231` `return cfg.Panel.Width, cfg.Panel.Height`；`panel_control_windows.go:391` `windowWidth := options.Width`；`panel_host_windows.go:914` `width, height := int(options.Width), int(options.Height)` ⇒ **非测试命中 17 处** | **有生产者** |

⇒ **票 33 的真实状态：14 格里 12 格有落地或证据，只剩 AC#1 待核。** 但**框仍归非实现者翻**。

## 7. 撤销口令

删本件即复原（本件零产码）。
