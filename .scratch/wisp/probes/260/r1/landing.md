# 票 260 · 落地腿 260-r1 —— 借用那遍吃 `[hotkey] cancel` 的配置值（形ⓐ＋三条硬条件）

**开工时刻**：2026-10-04 09:4x +0800　**起点 HEAD（现取）**：`528bf7bd`（`dev`）
**任务书**：票面 `.scratch/wisp/issues/260-configured-cancel-hotkey-never-registers.md` 最后一节「编排者选形裁定（2026-10-04 09:2x，账 A588）」＝形ⓐ＋条件①②③；撤销口令「260 改 ⓑ」。
**本格范围**：只做 **AC#1**（借用吃配置）。**AC#2（丢借用要有声）不在本格**，编排者另排一格；AC#3 是越界检查，见 §6。
**写面**：`internal/ball/**` ＋本证据件目录 `.scratch/wisp/probes/260/r1/`。⛔ 未碰 `cmd/wisp`／`internal/panel`／`internal/agent`／`tools/d22scan`（逐项核对见 §6）。

---

## §0 起手三把现量尺（全部现跑，未照任务书行号当常量）

### 0.1 `grep -n "escBorrowAcc\|vkEscape\|RegisterHotKey" internal/ball/*.go`

| 类别 | 读数 |
|---|---|
| | |

### 0.2 `grep -rn "cancel" internal/config/schema.go | head`

| 类别 | 读数 |
|---|---|
| | |

### 0.3 `grep -rln "hotkeys live\|Binding" internal/ball/*_test.go cmd/wisp/*_test.go` ＋未修码基线包名册

| 类别 | 读数 |
|---|---|
| | |

---

## §1 产码落点（四条形ⓐ落点，file:line 为改后现量）

| # | 落点 | 干了什么 |
|---|---|---|
| | | |

---

## §2 四条常驻判据（用例名＋断言形状，全部真写进 `internal/ball/`）

| # | 条件 | 用例名 | 断言形状（是否词面尺） |
|---|---|---|---|
| | | | |

---

## §3 门禁终值读数（逐字抄）

| 门禁 | 读数 |
|---|---|
| | |

---

## §4 突变读数（摘掉新逻辑 ⇒ 指名用例必红；还原 ⇒ 绿且 md5 一致、`git diff` 为空）

| 突变 | 红句（逐字） | 还原 |
|---|---|---|
| | | |

---

## §5 没做成的／欠账（具名＋归口）

| 项 | 归口 | 说明 |
|---|---|---|
| | | |

---

## §6 越界检查与票面禁区复核

| 检查 | 读数 |
|---|---|
| | |
