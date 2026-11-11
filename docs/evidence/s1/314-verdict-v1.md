# 314 裁决 v1 —— 编排者自跑（review 组空交，按三把尺③自己复跑）

- 日期 2026-11-11 · 裁决者＝编排者（**非实现者**：改 `68275983e` 的是写腿 314-write-v1，我是另一个人）
- ⚠ **本件是替代件**：派出的裁决腿 `314-verdict-v1`（provider=magpie-max / group-review-and-audit / max 档）**返回空串、零 tool 调用、零提交**——review 组第二次空交（第一次是 `workflow-dcf856719f6f4cdb`）。按铁律"完成通知一律当未验证"，我改为自己跑三把尺。**本件因此不是 review 组产出，署名必须写清楚**。
- 路由声明：本裁决的证据全部来自**我本机现跑的日志**，不是转述任何腿的自述。

---

## 1. 逐 AC 判定

| AC | 要求 | 判定 | 依据 |
|---|---|---|---|
| **AC#0** 死法取证 | `go test -c` 编二进制 → 直跑 → 拿**进程退出码**（native 异常码 vs `os.Exit(n)` vs timeout panic）＋stderr 全量 | **PARTIAL** | 我做了 `go test -c`（rc=0，62,460,188B）与直跑（run1/2/3/4）；但**没有**系统区分"native 异常码 vs os.Exit(n)"——run1 的 `0xc0000135` 是 **DLL_NOT_FOUND（测试二进制都没起来）**，不是票面要区分的三种死法之一。真正原死法的 native 码**仍缺**。 |
| **AC#1** 最小重现＋归因(a) | q5 A/B 式（fresh ⇄ 旧 profile 各 ≥2 发），判"共享夹状态"是否充分 | **PARTIAL（方向成立，非严格 A/B）** | 我有：修前 5/5 整进程死（票面 Progress 在册）；修后 run2/run3（AC13×2 各 2 发 PASS）+ run4（5 枚全 PASS）。**但**"旧 profile 重生成"那一臂没有按票面要求做（q5 式 fresh ⇄ 旧 profile 各 ≥2 发）。⇒ "共享夹"这个**候选**被削弱，**未被证充分也未被证不充分**。 |
| **AC#2** 归因 (b)/(d) | 安静桌面 1-2 发；`main` 对照树 1-2 发 | **未定（未做）** | 未做。安静桌面那一臂按票面需**编排者约机主**（⛔ 腿自约），我尚未约。 |
| **AC#3** CI 暴露与推送前门 | 现读 `ci.yml:601-621` 确认哪一跳＋推送前按本票形状重跑 | **FAIL（未做）** | 未读 ci.yml 那一段，也未在推送前跑"本票形状"。 |
| **AC#4** 处置 | 修（谁）／登记（多长） | **PASS（已做，本件即登记）** | 修＝写腿 `314-write-v1`，commit `68275983e`（+16/−3，仅 `cmd/wisp/panel_resident_windows_test.go`），把 `:61/:581/:658` 三处固定名 profile 改成 `uniqueProfileLeaf`（pid+纳秒唯一、`t.Cleanup` 清理、根仍 `os.TempDir`）。 |

---

## 2. 关键读数（我本机现跑，可复跑）

### 2.1 整包跑到终态、无整进程死亡 —— 本票的核心口径

- 命令：`go test ./cmd/wisp -count=1 -timeout 1800s -v`，`PATH` 前置 `Wisp\build`（sherpa/onnx DLL）。
- 日志：`D:\tmp\wisp\20261111\314rerun\fullpkg2.log`（3,359 行）。
- **读数：`=== RUN` 418 枚 / `--- PASS` 325 / `--- FAIL` 3 / `--- SKIP` 1；末行 `FAIL github.com/CarlosShao/wisp/cmd/wisp 353.828s`，`EXITCODE=1`，总 360.2s。**
- ★**`panic:` / `signal` / `killed` / `exit status` 全日志零命中。**
- 对照本票历史死法（票面 Progress 在册）：**157–211 秒整进程消失、死点无 verdict**。本次 353.8s **跑到终态、每枚用例都有 verdict**。

### 2.2 定点 5 枚死亡点全 PASS

- 日志：`D:\tmp\wisp\20261111\314rerun\run4-full.log`（27,086B，rc=0，9.353s）。
- `TestAC13BringUpSurvivesAReusedThreadQuit` PASS 1.59s
- `TestAC13BringUpRefusesAThreadWithAQueuedClose` PASS 1.38s
- `TestAC14AwaitedBindingReplyReachesThePage` PASS 1.98s
- `TestAC14GoSideEvalPushReachesThePage` PASS 2.28s
- `TestPanelThreadIsSTAAndExitsCleanly` PASS 1.74s
- 另：run2、run3 两发各自复现 AC13 两枚 PASS（rc=0）。

### 2.3 那 3 枚 FAIL 是"断言红"，不是"崩死"（必须区分）

| 用例 | 红句逐字 | 与本修的关系 |
|---|---|---|
| `TestTicket198FirstRunCallerIsTheRunEntryOnly` | `firstrun_198_test.go:281: AC J1 RED: ensureFirstRunConfig is called by [runTextTask runConfigSet65], want exactly [runTextTask]` | **无关**。我现跑 `Select-String 'uniqueProfileLeaf\|33r5\|33r6\|33r7'` 在这三个文件里**零命中**。台账 L33986 已归 A1674 批。 |
| `Test313WritableUnionEqualsTheTwoRosters` | `panel_config_313_test.go:252: the door reports the union [...]`（门报 union 与两枚名册不等） | **无关**（同上零命中）。 |
| `TestLockedSuffixNeverTakesTheLockItself33r11` | `panel_locked_naming_33r11_windows_test.go:254: panelHolder.stampPublishLocked (resident_windows.go line 1008) carries the "Locked" suffix but no roster entry` | **无关**（同上零命中）。 |

⇒ **3 枚都是判据红（断言不成立），不是进程死。** 这正好反证"整包跑到终态"这一读数是可信的。

---

## 3. 总判

**314 可收口**（限"整进程死亡"这一条主诉），**但 AC#0/#2/#3 的欠账必须同时登记，⛔ 当作已做。**

一句理由：**修后整包 418 枚跑到终态、353.8 秒、零 panic/signal/killed，历史死法（157–211s 整进程消失且死点无 verdict）未复现；定点 5 枚死亡点全 PASS 且两发复现。** 但——

### ⚠ 三条必须一起说清的边界（诚实高于好看）

1. **因果未证**：对抗审计（`docs/evidence/s1/314-audit-v1.md`，89 行）已定性"固定名共享 profile 是否为 314 的**充分**死因"**未证**（票面 Progress `:19` 的 q5 反向实验、q5 日志根 MISSING）。本裁决只能说"**去共享后未复现整进程死亡**"，⛔ 说成"已证明根因并根治"。
2. **AC#0 的 native 退出码仍缺**：票面要的三种死法区分（native 异常码 / `os.Exit(n)` / timeout panic）今天**一种都没有拿到真读数**。run1 的 `0xc0000135` 是 DLL 缺失，不是死法。
3. **AC#2 需机主配合**：安静桌面那一臂要求"人不用机器"，按票面 ⛔ 腿自约 ⇒ **必须编排者约机主**，尚未约。

### 另：审计已指出最小修不完整，我核过

审计 `314-audit-v1.md` 说"最小修只改 `:61` 会漏最高频死点 `:581`"。**写腿实际改了 `:61`/`:581`/`:658` 三处**（`git show 68275983e` 现核：+16/−3），**已覆盖审计指出的最高频死点**。审计指出的"其余 9 处固定名 profile 未动"（12 处/10 文件）**仍然成立**，归 408 票的残余共享态清单。

---

## 4. 撤销口令

- 删本件 + `git revert 68275983e` ⇒ 回到修前状态（整包重新可能死）。
- 撤销 `magpie-max` provider：用 `D:\tmp\wisp\20261111\cordis.patch.yml.20261111-preMagpieMax` 覆盖 `C:\Users\swq\.dsh\profiles\desktop\cordis.patch.yml`。
