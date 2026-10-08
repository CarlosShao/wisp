# 180-c1 起手闸门 · 锚件（只读普查腿）

现量时刻：`2026-10-02`（本程为**第三枚** `180-c` 编号腿；票内 09-28 14:4x 已有一枚同名 `180-c1` 的记录，见票 §Progress log 第 38 行——**编号复用，不是本程伪造**，本程以本文件为唯一自称）。

## 1. 起手 HEAD / 分支 / 工作树

尺与原文：

```
git log -1 --format=%H            -> 601c2b1872c29f2ec222574cd9f9fbcb2f214eb0
git rev-parse --abbrev-ref HEAD   -> dev
git status --porcelain | wc -l    -> 753
```

`git log -1 --format=%H` rc=0 ／ `git rev-parse` rc=0 ／ `git status --porcelain | wc -l` rc=0。

起手 HEAD **等于**派单给的地板 `601c2b18`（按派单口径：以现量为准，非"必须等于"）。
工作树 753 行为别人的在飞改动，本程一字不动（见硬约束 3）。

## 2. 票面 AC#1 那句前提的原文（逐字抄）

出处 `.scratch/wisp/issues/180-panel-width-is-a-config-field-with-zero-production-readers-so-changing-config-toml-has-no-visible-effect.md:22`（AC#1 那一整格的行内文，本程不截断）：

> - [x] **AC#1 先判归属，再谈修法**（09-28 17:5x 编排者复跑对格后勾：三问答完且我今天自己重跑过——**没有任何人**决定面板尺寸，因为这棵树里**根本没有面板窗口**：`internal/panel/pump.go:15-17` 逐字"no WebView2 host"、`grep -rn "NewWindow\|SetBounds\|Rect{\|Width:" --include=*.go internal/panel/ | grep -v _test`＝**0 命中**、那五枚 `[panel]` 字段的生产读取方各 **0**；规格要它生效＝`PLAN.md:2743` 逐字 `hot`；唯一真窗口是球（`internal/ball/ball_windows.go:233`，与 `[panel] width` 无关）＝账 `A391`）

⇒ 本程要重验的可证伪断言有三条，分开数：
- **(P1)**「这棵树里根本没有面板窗口」
- **(P2)** `internal/panel/pump.go:15-17` 逐字写着 "no WebView2 host"（以及票 09-28/17:5x 那几处把它读成"ticket 33 is unclaimed"）
- **(P3)** 尺 `grep -rn "NewWindow\|SetBounds\|Rect{\|Width:" --include=*.go internal/panel/ | grep -v _test` = 0 命中

## 3. AC#4 入口尺：`default:` 的现量（我自己复跑）

```
grep -c 'default:' internal/config/schema.go   -> 70        rc=0
```

**70** —— 与派单给的入口尺 70 一致。

⚠ 但与票内 **09-28 `180-c1`**（第 47 行：「带 `default:` 的行 **69**」）和 **10-02 `180-a1` 第二程**（第 91 行：「带 `default:` 标签 **69** 枚（与 09-28 `180-c1` 独立对上）」）**冲突**。
两说都要现验：`grep -c` 数的是**含该子串的行数**，不等于**标签枚数**（一行可能两枚、注释里可能出现 `default:`）。
本程后续会分别给出：行数／标签枚数／去重字段名数 三个数，并把与 69 的差具名定位到行。
**本锚件以现量 70 行为准，并具名报回该冲突。**

## 4. 本程边界（自检）

- 零 go 命令（不 `build`/`vet`/`test`）；`cmd/wisp` 测试面归另一枚在飞腿 `198-v1`。
- 只读产码：`internal/**` `cmd/**` `frontend/**` `docs/PLAN.md` `docs/specs/**` 零字节改动；写点＝票面追加一节 + 本目录。
- 不碰 `docs/reports/pending-and-issues.md`（台账归编排者）。
- 不开真窗／不起 GUI／不抢鼠标。
- 证据件只用 `.md`/`.txt`/`.tsv`，每把尺自落 `rc=N`，无 0 字节件，无 `.out`。
