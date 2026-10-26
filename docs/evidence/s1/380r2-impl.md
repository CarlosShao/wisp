# 380r2 落地腿（写腿）—— `C″`：78 行一行不删 ＋ 两张 2 列对照表 ＋ 四条替代判据

主题：`票 380 落地腿 380r2`。父锚点 `413653b96`，现跑 2026-10-26 11:50–12:2x +08。
机主批的形＝**`C″`**（逐字：「C″ 不删行、加两张对照表（推荐）」）。
写面 4 件：`docs/evidence/s1/c21-native-tokens.md`（只加节）· `internal/panel/tokens_fourway_test.go` · `docs/PLAN.md`（两行）· 本文件。
临时件（只建不删）：`D:\tmp\wisp\20261026\380r2\{m1,m2,m3,m4,m5,m6,logs}` —— **6 棵仓外导出树**，每棵一发变异，逐发可复跑。
开工前核索引：`git diff --cached --name-status` **空**。

## 1 九项读数（逐条）

| # | 读数 | 结果 |
|---|---|---|
| 1 | 78 行零删 | `git diff --numstat` ＝ **73 插入 / 0 删除**；`## 明色系` **40→40**、`## 亮色系` **38→38**（§2） |
| 2 | `internal/ball` 零影响 | 提交只含 4 个文件、`internal/ball/` 命中 **0**；`go test -run TestC21 ./internal/ball` 改前 **rc 0** / 改后 **rc 0**（§3） |
| 3 | `P4` 带色 token ＝ 34 | 自数 **34**，**43 − 9**（6 条 `--shadow-*` ＋ `--panel-alpha` ＋ 2 条字体栈），逐名见 §4；**11 ＋ 23 ＝ 34** |
| 4 | 四条判据各自的牙 | m1→②③④／m2→②钉死清单／m3→③／m4→①（§5） |
| 5 | 假绿自证 | **m5**：`--inset` 暗值 `#1f2022`→`#1f2023` ⇒ **rc 1、1 条 error**，**名字集合一个不动**（86 条声明、同哈希、diff 仅 1 行）（§6） |
| 6 | 极性修正的牙 | **m6**：极性装回旧读法 ⇒ 当场 **`:620 no :root block parsed`** fatal，其余 **0 条**（§7） |
| 7 | `PLAN.md` | **只改 `:3308`／`:3604` 两行**（各 `1 - / 1 +`），总行数 **3611→3611**；`:1377` **未改**（§8） |
| 8 | 改后门禁 | `./internal/panel` **rc 1**（4 条既存红，**无新增**）／`./internal/ball` **rc 0**／`-run TestC21 ./internal/ball` **rc 0**／`go build ./...` **rc 0**／`go vet ./internal/panel` **rc 0**（§9） |
| 9 | 格式甲乙两尺 | `gofmt -l internal/panel` **空**；`gofumpt -l internal/panel`（绝对路径）**空**；归属判用 `[\\/]tokens_fourway_test\.go$` 命中 **0**（§9） |

## 2 78 行一行不删的证明

- `git diff --numstat -- docs/evidence/s1/c21-native-tokens.md` ＝ **`73  0`**（73 行加、**0 行删**）。
- `git diff -U0` 里以 `-` 开头的行（排除 `---` 文件头）＝ **0 条**。
- 行数/行数改写前后的整表计数（同一把尺，改前取 `git show HEAD:…`）：`## 明色系` **40 = 40**、`## 亮色系` **38 = 38**、合计 **78 = 78**；文件 **298 → 371 行**（＋73）。
- 新增的只有两节，插在「亮色系」与「几何／动效」之间，两节都是 **2 列**：`## 第四代槽位对照` **11 行**、`## 第四代未覆盖` **23 行**。
- ⇒ **11 ＋ 23 ＝ 34 ＝ `P4` 带色 token 数**，`git` 逐字核过：两节名字并集与 `P4` 的 34 名**完全相等**（无缺无多）。

## 3 `internal/ball` 零影响的证明

| 尺 | 改前 | 改后 |
|---|---|---|
| `go test -count=1 ./internal/ball` | rc **0** | rc **0** |
| `go test -count=1 -run TestC21 ./internal/ball` | rc **0**（5 条 PASS） | rc **0**（5 条 PASS） |
| 表插入两节**之后**单跑一发 | — | rc **0**（`logs/after-doc-ball-c21.txt`） |
| 提交的 `--name-only` | — | **0 条**命中 `internal/ball/` |

原因（不是运气）：球侧解析器 `internal/ball/tokens_table_test.go:181` 是 `if len(cells) != 4 { continue }`，且 `:154-168` 只认
「明色系／亮色系／色板／几何／豁免」五个 `#` 节名 —— 两节都是 **2 列**、节名都不含这五个词 ⇒ 球侧**整节跳过**。
球侧的 40/38 下限（`:673-674`）与「每个 `Palette` 字段必须有暗行」（`:760-761`）**一行都没动，也就一行都没少**。

## 4 `P4` 带色 token ＝ **34**（自数，逐名）

`frontend/src/styles/tokens.generated.css`：`:root` **43** 条声明（＝ light）、`[data-theme="dark"]` **43** 条（＝ dark）；
两块的**名字集合相等**。判「带色」＝ 值本身是 hex/rgb/rgba、或 `var()` 别名指向一个带色 token、或 `color-mix(...)`。
非带色 **9** 条：`--shadow-hairline` `--shadow-btn` `--shadow-card` `--shadow-raised` `--shadow-overlay` `--shadow-inset-field`
`--panel-alpha` `--font-sans-stack` `--font-mono-stack` ⇒ **43 − 9 ＝ 34**。
⇒ 与票面 `c21-a1` 的「第四代 43 个 token（带色 34）」**逐数相同**，不触发停手①。

**11 名（表里已有可核对槽位，`## 第四代槽位对照`）**：`--inset` `--surface` `--line` `--ink` `--ink-2` `--ink-3` `--accent` `--danger` `--warn` `--warn-line` `--success`
**23 名（原生世代不覆盖，`## 第四代未覆盖`，逐名一句理由）**：`--page` `--canvas` `--hover` `--hover-2` `--field` `--stripe` `--stripe-bg` `--line-strong` `--accent-ink` `--accent-tint` `--green` `--green-tint` `--orange` `--orange-tint` `--red` `--red-tint` `--tooltip-bg` `--tooltip-fg` `--tooltip-muted` `--tooltip-border` `--warn-soft` `--stop` `--success-soft`
**11 ＋ 23 ＝ 34** ✔（下文 `§5` 的等价对账读数：`check 2: … declares 34 coloured tokens = 11 mapped + 23 waived`）

## 5 四条判据 ＋ 各自的牙（逐条造变异、逐条报红）

落地位置：`TestC21DesignTokensFourWayAgree`（`internal/panel/tokens_fourway_test.go`），`checked < 70` **已删**，**不是**改成 `checked < 22` 了事：
`70` 那道门**没有任何对象可量**（旧世代表 78 行 vs 第四代生成件 43 条），四条判据是**换测量对象**，`22` 只是第 4 条的**下界且是天花板**（11 名 × 2 主题）。

| # | 判据 | 落地 | 改后绿读数 |
|---|---|---|---|
| ① | 表引用的名必须在 `P4` **两块都声明** | 11 个槽位 × 2 主题逐个查 | 22/22 通过 |
| ② | `P4` 的 **34** 带色 token ⊆ 表 ∪ 豁免表；**豁免集合 ＝ 代码里钉死的 23 名字面量清单** | `c21FourthGenUncovered`（钉死）＋「映射必须有配色行兜底」的连接语义 | `34 = 11 + 23` |
| ③ | 两块名集合相等、各主题 reconciled>0、**行数 == 2×名字数**、每名两主题各一行 | 11 名 × 2 ＝ **22 行** | 40/38 字段、22 行齐 |
| ④ | **下界 22**（不是 28）：`--accent-line`／`--danger-line`／`--success-line` 在 `P4` **无对象可校** ⇒ 14 − 3 ＝ 11 名 ⇒ 22 | `c21FourthGenMappedPairs = 22` ＋ `len(slots) != 22/2` 报警 | `22 of 22` |

| 发 | 树 `D:\tmp\wisp\20261026\380r2\` | 变异 | rc | 逐条报红 |
|---|---|---|---|---|
| **m1** | `m1` | 表删到 **5 行** | **1** | ②×**9**＋③×**18**＋④ fatal「only 5 of the 22 …」＝ **②③④ 同时红** |
| **m2** | `m2` | 删掉 11 名 × 2 主题的 **22 行** ＋ 往 `第四代未覆盖` **补 11 行豁免** | **1** | ②×**11**（全是 `is not in the pinned c21FourthGenUncovered list`）＋③×23＋④ fatal。★**计数仍自洽**：`11 mapped + 34 waived`、`0 + 34 == 34` 不报 ⇒ **只有钉死清单能抓「补豁免」** |
| **m3** | `m3` | 只留 dark 行（删 11 名的小节亮色行） | **1** | ③×**12**（11 条「no light colour row」＋1 条「11 mapped colour rows … want 22」）＋④ fatal；②**绿** |
| **m4** | `m4` | `P4` 的 `:root` 删掉映射槽位 `--inset` | **1** | ①×**1**（`:771 … the light block … does not declare it - half a decision`）＋③×1＋②×2＋值腿×2＋④ fatal |

## 6 假绿自证：**形状变异 与 值变异 正交**

★两发**正交**：`m1` 动**行数**（名字集合也会动），`m5` 只动**一位值**（名字集合一个不动）。

| 发 | 变异 | rc | 读数 |
|---|---|---|---|
| **m5** | `--inset` **暗值** `#1f2022` → `#1f2023`（`P4` 内，仅此一处） | **1** | 全程只有 **1 条 error**：`:670 … dark declares --inset = "#1f2023" …, but gen-tokens.mjs says "#1f2022" - a hand edit, not a decision`。① ② ③ ④ **全绿**（`22 of 22`、`34 = 11 + 23`） |

★名字集合不动的证据：改前/改后 `P4` 均有 **86** 条声明（43×2），排序后名字串**哈希相同**；逐行 diff **只有第 75 行**一条、只差一个字符。
⇒ **结论**：抓漂移的**不是行数**，是「名字集合相等（正反两向）＋**逐值相等**＋两主题各占一行」这三件；行数只是必要条件。
★这一发对 C″ 是**必需**的：`--inset` 在 78 行表里**没有行**（它是第四代独有），所以值相等只能由「生成器表 ↔ 生成件」这条腿提供 —— 见 §10 偏离①。

## 7 极性修正的牙

| 发 | 极性 | rc | 读数 |
|---|---|---|---|
| 改后 | `:root` → **light**、`[data-theme="dark"]` → **dark** | **0** | `value leg theme light: 43 tokens reconciled` ＋ `… dark: 43` |
| **m6** | 装回旧读法（`:root`→dark、只认 `[data-theme="light"]`） | **1** | **唯一**一条：`:620 frontend/src/styles/tokens.generated.css: no :root block parsed - :root is the LIGHT block … an empty light leg makes every check below vacuous` ⇒ 与 `380-r1` §7 的「light 腿空跑 fatal」**同一句、同一落点** |

根因（这次查清、可核）：**第一代把默认暗色放 `:root`、亮色放 `[data-theme="light"]`**（盘上幸存的那份第一代参考实现
`internal/ball/testdata` 的 `tokens_c21.css`：`:38 :root {` ＋ `:243 [data-theme="light"] {`，sha256 `402E3231…`），
**第四代把极性反过来了**（`:root` = light `--page: #fafafb`）。旧测试是按第一代的约定写的，被套到第四代上 ⇒ 装反。

## 8 `PLAN.md` 逐行 diff（只有这两行）

```
@@ -3308 +3308 @@
-    │   ├── tokens.css           ← C21 DesignTokens 的 CSS 参考实现（**唯一样式真相源**）
+    │   ├── tokens.css           ← 第一代 CSS 参考实现（**原件已在 `72b65c614` 删除**）；C21 的**唯一样式真相源**改为 `frontend/src/styles/tokens.generated.css`（第四代生成件，由 `frontend/scripts/gen-tokens.mjs` 的 `TOKENS` 表逐值产出）：原生侧 78 行对照表**原样保留**＋**两张 2 列对照表**（11 名槽位对照／23 名未覆盖逐名理由），第四代覆盖 **34** 个带色 token
@@ -3604 +3604 @@
-| **C21 `DesignTokens`** | `tokens.css` 就是它的**参考实现**。S1 的 Direct2D 侧与 S5 的 CSS 侧都以它为准。**原型阶段发现 token 不够用 → 现在就补，不要等到 S5** |
+| **C21 `DesignTokens`** | 真相源是 **`frontend/src/styles/tokens.generated.css`**（第四代生成件，由 `frontend/scripts/gen-tokens.mjs` 的 `TOKENS` 表逐值产出；`:root` = light、`[data-theme="dark"]` = dark）。`design/assets/tokens.css` 是第一代参考实现，**原件已在 `72b65c614` 删除**。S1 的 Direct2D 侧与 S5 的 CSS 侧都以生成件为准；原生侧那张 78 行对照表（`docs/evidence/s1/c21-native-tokens.md`）**原样保留**，另加**两张 2 列对照表**（11 名→第四代槽位／23 名未覆盖逐名理由）⇒ 第四代覆盖 **34** 个带色 token。**原型阶段发现 token 不够用 → 现在就补，不要等到 S5** |
```
`git diff --numstat docs/PLAN.md` ＝ **`2  2`**，文件 **3611 → 3611** 行。`:1377`（`C21` 那一行）**没改**：它只说「单一 token 源…原生侧与 CSS 侧各自实现同一份」，**不点任何路径** ⇒ 没有腐烂的字要改。

## 9 门禁 rc ＋ 格式甲乙两尺

| 门 | 改前 | 改后 | 失败名单 |
|---|---|---|---|
| `go test -count=1 ./internal/panel` | **rc 1**（**5** 红） | **rc 1**（**4** 红） | 改后 4 条＝**既存红**：`TestApprovalCardViewJSONKeysMatchFrontendTypes`／`TestComposerContractTypesMatchFrontend`／`TestStreamLogFloodBelowKeyBoundIsNotBounded35r8`／`TestStreamLogDroppedNamingLedgerIsNotBounded35r8`。★`TestC21DesignTokensFourWayAgree` **已由红转绿**，**无新增红**（未相减，逐条列名） |
| `go test -count=1 ./internal/ball` | rc 0 | **rc 0** | —— |
| `go test -count=1 -run TestC21 ./internal/ball` | rc 0 | **rc 0** | —— |
| `go build ./...` | —— | **rc 0** | —— |
| `go vet ./internal/panel` | rc 0 | **rc 0** | —— |
| `tools/d22scan`（`go run . -root ../../`） | —— | **rc 0** `clean`（ban#8 实扫 `internal/` **652** 个 Go 文件，含注释与 `_test.go`） | —— |

**格式甲乙两尺**（都报）：甲＝`gofmt -l internal/panel` ⇒ **空**；乙＝`/d/work/base/gopath/bin/gofumpt.exe -l internal/panel` ⇒ **空**（`-d` 单文件也空）。
★归属判（防"别人的文件算我头上"）用 `[\\/]`：`gofumpt -l internal/panel \| grep -E "[\\/]tokens_fourway_test\.go$" \| wc -l` ＝ **0**。
环境：`GOFLAGS=`、**只** `unset GOTMPDIR`；rc 一律 `cmd > out 2>&1; rc=$?`（**⛔ 走管道**）。**工作树零突变**：6 发变异全在仓外导出树，核过仓内 `--inset` 仍是 `#1f2022`。

## 10 三处必须申报的偏离／发现（⛔ 自行改口径）

1. ★★**`P1` 换了载体**（`design/assets/tokens.css` → `frontend/scripts/gen-tokens.mjs` 的 `TOKENS` 表）。理由不是"想换"：① 旧 `P1` **已被删**（`72b65c614`），只剩球侧夹具，而夹具是**第一代极性**，极性一修就撞 g1（`380-r1` §7 已实测）；② 值变异那一发（§6）**必须**有一个"逐值相等"的对手方，`--inset` 在 78 行表里**无行**，第四代唯一可机器的真相源就是生成器表。生成器自己已有 `--check`（需 node），本测试**不依赖 node**、只读它的 `TOKENS` 表，与「生成件 == 生成器产出」这条既有纪律**同向**。
2. ★**`第四代未覆盖` 表首列是「第四代 token（原生世代无对应）」**，不是派单括号里逐字的「原生世代 token」。理由：11 ＋ 23 必须**都数同一侧**才等于 **34**（数原生侧会是 40 − 11 ＝ 29）；照抄列名会让表自称 34 而实际 11 ＋ 29，属文档说谎。列名保留了「原生世代」四字，表意不丢。
3. ★**`C″` 不删行 ⇒ 78 行表里 50 行的"只此一份书面记录"仍在**（玻璃/orb/`--tint-*`/`--info`/`--warm`/`--ball-halo`/`--on-solid` 等 26 个 `Palette` 字段），这是选 `C″` 买到的**真价值**，不是遗留债；面板侧不再逐行核它们，但球侧 5 条 `TestC21*` 仍逐行核（§3 实跑）。
4. ★**发现（不在我写面，未动）**：`docs/PLAN.md` 还有 **7 处**把 `tokens.css` 当真相源（`:3342`／`:3394`／`:3395`／`:3398`／`:3408`／`:3422`／`:3592`）。派单只批 `:3308`／`:3604` ⇒ 我一处没碰，登记给编排者决定是否追一票。
5. ★**发现（不在我写面，未动）**：`internal/panel/frontend_hygiene_test.go:217` 的报错文案仍写「Add or reference a C21 token in design/assets/tokens.css instead」——那条文案现在指向一个已删文件；属文案腐烂，不影响断言，登记。
