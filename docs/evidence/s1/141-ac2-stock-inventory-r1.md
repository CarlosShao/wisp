# 141 AC#2 —— ban #8 缺口段存量盘点（只读程 r1）

**锚点**：`99263cc`（本程开工自量：`git rev-parse --short HEAD` → `99263cc`，分支 `dev`）
**agent**：`auditor-ticket141-inventory-r1`
**性质**：**只读盘点**。本程**没有运行任何仪器**：零次 `go build` / `go test` / `go vet` / `gofmt` / `gofumpt` / `docker` / `wisp` / `wisp slo` / `scripts/slo-check.ps1`，零次 `gh`。全部读数是 `git grep` / `git show` / `git blame` / `sed -n` / 自写的**文本分类脚本**（在 `D:\tmp\wisp141inv\`，不在仓内）。同机 `acceptor-ticket136-ac14-r1` 正在取真读数，本文件不与任何读数争时段。
**只做 AC#2**（纯盘点，不依赖 `Q-46`）。`AC#1` 的选支**不由本程决定，本程也不建议**——见文末「我不替 owner 选支」。

> **本节末尾回执的约定**（写在最前面，免得后一位以为缺了）：共树不得 `--amend`，
> 所以**每一节正文末尾贴的是「上一枚 commit」的回执**；本节自己的回执由下一节带来，
> 最后一节的回执在文末「commit 流水」汇总。所有回执都是 `git log --oneline -1` ＋
> `git show --name-only HEAD` 的原样输出。

---

## 0. 复核：简报给的四个数，与两张要补的账

### 0.1 复核实到（票面用 `git grep -P '[\x{2190}-\x{25FF}]' -- 'internal/**.go' 'cmd/**.go'`，我逐条重跑）

| 简报的数 | 我的复量 | 判定 |
|---|---|---|
| 41 枚 `.go` 文件 | **41 枚**（`git grep -lP` 于 `99263cc`，实测见 §5 命令） | ✅ 对得上 |
| 116 行 | **116 行** | ✅ |
| 121 处（逐字符） | **121 处** | ✅ |
| 注释 78 行 / 非注释 38 行 | **78 / 38**，但**判据被换掉了**（见 0.2） | ✅ 数值对，方法不采信 |
| 非注释散在 11 枚文件 | **11 枚** | ✅ |
| **12 枚在禁改清单** | **41 枚里只有 11 枚**（`internal/risk/` 7 ＋ `internal/winsec/` 4）。第 12 枚 `tools/d22scan/scan_test.go` **根本不在 41 枚的取数 pathspec 里**，而且**不在 ban #8 的 walk scope 里** ⇒ 三支选项里**没有任何一支会因它变红** | ❌ **推翻** |
| 今日新增 0 | `git diff -U0 182daed..HEAD -- '*.go'` 加号行 gap 命中 **0**；`-- 'frontend/'` 亦 **0**（`182daed` 存在，09-24 15:04） | ✅ |

**票面字形分布少列了两枚**：票面「`②`32／`①`20／`⇒`18／`→`16／`③`15／`⑤`7／`⑥`6／`⑦`2／`≥`2／`≤`1」加起来是 **119**，不是它自己写的 121。实测 12 枚码点、总和 121，漏的是 `U+2467` ×1（在 `internal/risk/syncdirs_other_test.go`）与 `U+2229` ×1（在 `internal/risk/pathresolver_anchor_spelling_windows_test.go`）。**总数没错，分布表缺两行。**

### 0.2 推翻的方法：注释／非注释不是「前面有没有 `//`」

票面写「我按"命中位置前有没有 `//`"切的」。那个切法在两处会错，本程不沿用：

1. **它认不出块注释**，也会把**串内的 `//`**（如 `"https://…"`）当成注释起点；
2. **它认不出 raw string（反引号串）里的 `--` SQL 注释**——那在 Go 语法里是字符串，按票面切法会被算成"非注释"，但它在 SQL 语义里**正是注释**。

本程改用一枚**Go 迷你词法器**（`D:\tmp\wisp141inv\classify.pl`，状态机：`code`／`c-line`／`c-block`／`s-dq`／`s-raw`／`s-rune`，处理转义与跨行 raw string），逐 rune 记录命中时刻的状态。取数方式：`git cat-file blob 99263cc:<file>` 把 41 枚文件抽到 `D:\tmp\wisp141inv\tree\`（**不在仓内**）后分类，所以读的是**锚点那版**、不是工作树。

词法器输出 116 行／121 处，与票面**数值全等**，状态分布为：

| 词法状态 | 行数 | 含义 |
|---|---|---|
| `c-line` | **78** | `//` 行注释 |
| `s-dq` | **36** | 双引号字符串字面量内 |
| `s-raw` | **2** | 反引号 raw string 内（**都是 `internal/memory/schema.go` 的 `ddlV1` DDL 文本里的 SQL `--` 注释**） |
| `c-block` / `code`（标识符） | **0** | —— |

⇒ **两个结论**（都影响 owner 的代价估算）：
- **38 枚非注释命中全部在字符串字面量里，一枚都不在标识符里**。简报说「测试名 `func Test…` 里的箭头」这一档**在本仓实测为 0**：Go 标识符不能含这些字形，票面上那 20 枚 `→`/18 枚 `⇒` 没有一枚在函数名位。
- 那 2 枚 raw string 命中是**「Go 眼里是字符串、SQL 眼里是注释」**。若 `Q-46` 走 (c)「注释豁免」，这 2 行**落在哪一边需要 owner 单独认一句**——词法上讲它们在字符串内（不豁免），语义上讲它们是 SQL 注释（该豁免）。本程只把这个岔口标出来。

### 0.3 补的第二张账：票面的 pathspec 只圈了 ban #8 射程的一半

`emojiScopes()`（`tools/d22scan/main.go:480-487`，`git blame` 该函数体所在版：`5e8f87b` 09-21 15:36 是这枚文件最近一次改动）声明 ban #8 **实际走四棵树**：`design/`、`frontend/`（`everyFile`）、`internal/`、`cmd/`（后两棵 `goOnly`）。票面的取数 pathspec 只有 `internal/**.go` ＋ `cmd/**.go`，**漏了 `frontend/`**——而 `frontend/` 恰是这条禁令**动机上真正要管的那棵树**（面板就是 UI）。实测（同一 `99263cc`）：

| 树 | 文件 | 行 | 处 | 是否 ban #8 walk scope |
|---|---|---|---|---|
| `internal/` + `cmd/`（`.go`） | 41 | 116 | 121 | 是（`goOnly`） |
| `frontend/` | **10** | **25** | **921** | **是**（`everyFile`） |
| `design/` | 0 | 0 | 0 | 是——PLAN.md:3574 真点名的那棵树，**今天确实是 0** |
| `tools/`（含 `tools/d22scan/`） | 1 | 1 | 1 | **否**（`emojiScopes()` 里没有 `tools/`） |

⇒ **补宽后的真实爆炸半径是 51 枚文件／141 行／1042 处**，不是 41／116／121。差额里 `frontend/` 一枚就贡献 921 处（96% 的字符量），且 `frontend/` 的 25 行里有 **6 行是渲染文本**（`≤`、`−`），不是注释——详见 §1 与表 B。

---

## 1. 小结① —— 非注释那 38 行的两档

### 1.1 判据（**不是**"看起来会不会显示"）

一条命令、一张 pattern 表，可复跑：

```
perl /d/tmp/wisp141inv/sink.pl <文件> <行号>
```

它先沿**字符串拼接续行**（上一行以 `+` 结尾）爬回**打开这枚字面量的那条语句**，
再按下面的 pattern 顺序给它贴 sink 标签。**分档规则 = 这枚字面量是不是「要被写出去的东西」**：

| 标签 | pattern（语句首行） | 档 |
|---|---|---|
| `test-print` | `t\.(Fatalf\|Errorf\|Logf\|Fatal\|Error\|Log\|Skipf\|Skip)\(` | **(i)** 进 `go test -v` 输出＝CI 可见 |
| `render-field` | `(\w*(Text\|Reason\|reason))\s*:\s*fmt\.(Sprintf\|Sprint)` | **(i)** 进面板/回执 |
| `audit-ledger` | `se\.record\(` | **(i)** 进步骤账本（渲染证据见 1.3） |
| `struct-field` | 仅 `(\w*(Text\|Reason\|reason))\s*:`（右边不是 `fmt.`） | **(ii)** 输入 fixture 或表驱动期望值 |
| `parse-template` | `Sscanf\(` | **(ii)** 解析模板 |
| `NONE` | 上面都不沾（raw string 里的 SQL DDL 文本） | **(ii)** |

### 1.2 结果：**(i) 31 行／7 枚文件**，**(ii) 7 行／4 枚文件**（合 38／11，与 §0 全等）

**(i) 会进可见输出 —— 31 行**

| 文件 | 行 | sink | 谁把它写出去 |
|---|---|---|---|
| `internal/ball/tokens_test.go` | 203,206,216,220,224,228,234,237,247,278（**10 行**） | `t.Errorf` | `go test -v` |
| `internal/tools/fs_staging_windows_test.go` | 66,73,157,235,246,337,365,371（**8 行**） | `t.Fatalf` | 同上 |
| `internal/winsec/reparse_windows_test.go` | 121,126,236,245（**4 行**） | `t.Logf` | 同上（`-v` 下**通过时也打印**，比 Fatalf 更外露） |
| `internal/tools/fs_write.go` | 346,480（`se.record`）、485,567（`Text: fmt.Sprintf`）（**4 行**） | 生产渲染 | 见 1.3 |
| `internal/agent/approval/ticket97_alias_direction_test.go` | 119,132（**2 行**） | `t.Fatalf` | `go test -v` |
| `internal/tools/bridge_a18_kill_windows_test.go` | 220,222（**2 行**，同一枚 `t.Fatalf` 的拼接续行，语句首行在 218） | `t.Fatalf` | `go test -v` |
| `internal/risk/rules_scale.go` | 24（**1 行**） | 生产渲染 | 见 1.3 |

**(ii) 只在源码字符串里 —— 7 行**

| 文件:行 | sink | 为什么落 (ii) |
|---|---|---|
| `internal/llm/anthropic/cache_test.go:21,22,23` | `struct-field:Text` | `llm.TextPart{Text: "①identity…"}` 是**喂给 mock 的请求输入**，不是消息文案 |
| `internal/risk/assessor_test.go:154` | `struct-field:Reason` | 表驱动 `want:` **期望值**，与生产串比对用 |
| `internal/tools/fs_write_test.go:521` | `parse-template` | `fmt.Sscanf` 的**解析模板**，且整个 `if` 被 `; false` 守卫短路 |
| `internal/memory/schema.go:29,118` | `NONE`（raw string） | `ddlV1` 的 **SQL DDL 文本**，送进 SQLite，不是任何输出流 |

### 1.3 生产那 5 行的外露证据（逐条点名，不用推测）

- `internal/tools/fs_write.go:346`／`:480` 的 `se.record(...)` → `sideEffect.steps`；`fs_write.go` **最近一版 `ed74595`（09-21 10:58）** 的 `record` 函数自带注释逐字写着「**the report renders these verbatim under 「已执行」**」（`internal/tools/fs_write.go:70-72`，同上一版）。
- `internal/tools/fs_write.go:485`／`:567` 的 `Text: fmt.Sprintf(...)` → `tools.Result.Text`，定义处 `internal/tools/tool.go:40-43` **最近一版 `0986d63`（09-20 19:16）** 注释逐字：**「Text is the payload the host puts back into the context (the loop caps and spills it, D15(3) …)」** ⇒ 回灌上下文＝会进对话呈现。
- `internal/risk/rules_scale.go:24`（最近一版 `67ffbd8` 09-20 08:22）的 `reason: fmt.Sprintf("…（`U+2265`%d）")` → `risk.Decision.Reason` → `internal/panel/approval.go:83` `Reason: decision.Reason`；该文件**最近一版 `63ef895`（09-21 13:26）**，`:47` 注释逐字：**「Reason is Decision.Reason verbatim」** ⇒ **审批卡原样渲染这枚 `U+2265`**。这是 38 行里唯一一枚**已经在真面板上外露**的缺口段字形。

### 1.4 两档的**边界声明**（owner 若要换判据，数会怎么变）

`struct-field` 那 4 行**并非绝对不外露**：
- `internal/risk/assessor_test.go:154` 的期望值，在断言不等时由 `:178-181` 的 `t.Fatalf("%s:\n  want %+v\n  got  %+v", tc.name, tc.want, got)` **整块打印**（`internal/risk/assessor_test.go` 最近一版 `67ffbd8` 09-20 08:22）；
- `internal/llm/anthropic/cache_test.go:21-23` 的输入 fixture，在解析失败时被 `:54`／`:65` 的 `t.Fatalf("…\n%s", body)` / `("system field: %v\n%s", …)` **连带请求体整块打印**（该文件最近一版 `f342413` 09-21 08:25）。

⇒ 若 owner 把 (i) 定义为「**任何可能到达终端的字面量**」，这 4 行改判 (i)，分布变成 **(i) 35 ／ (ii) 3**（只剩 Sscanf 模板 1 ＋ SQL raw 2）。
本程**按 sink 位取 (i) 31／(ii) 7**，并把改判路径写清楚——**不替 owner 决定用哪条判据**。

### 1.5 顺手量的两枚**耦合**（清存量时会互相拽一下）

- `internal/tools/fs_write.go:485`／`:567` 与 `internal/tools/fs_write_test.go:521` **看着像镜像、其实不是**：生产侧回收站文案（`fs_write.go:417`，`ed74595` 09-21 10:58）写的是「已放入回收站：%s（还原记录 %s，可在回收站还原）」——**不含缺口段字形**；测试侧那枚模板却写着「…（条目数 %d `U+2192` %d」。⇒ 这是一枚**陈旧模板**，与生产串不同步。清理它不需要与生产同批改，但**动它就在动一枚断言的输入**（AC#4 的"放水"判据会看这个）。
- `internal/risk/rules_scale.go:24`（生产）与 `internal/risk/assessor_test.go:154`（期望值）是**同一枚 `U+2265` 的两半**：改生产不改期望＝断言红。**必须同批**。两枚文件都在禁改清单（§3）。

---

<!-- 回执 ade897c：本节开头贴上一节那枚 commit 的原样输出 -->
```
$ git log --oneline -1
ade897c docs(141 AC#2 盘点 r1 第1节): 复核实到 41/116/121 全等, 推翻两处—— "12 枚在冻结路径"实为 11(第12枚 tools/d22scan/scan_test.go 不在取数 pathspec 也不在 ban#8 walk scope); 票面 pathspec 漏了 frontend/(ban#8 真射程), 补宽爆炸半径是 51 枚/141 行/1042 处。注释切法改用 Go 词法器, 38 枚非注释 全在字符串字面量内、0 枚在标识符。
$ git show --name-only HEAD
    全在字符串字面量内、0 枚在标识符。

docs/evidence/s1/141-ac2-stock-inventory-r1.md
```

---

## 2. 小结② —— 按码点段：owner 要选「补宽到哪一段」，这是各段的价

### 2.1 先说一件会改变这张表形状的事：**四段不够，缺第五段**

简报要的四段是「箭头 `U+2190–U+21FF`／带圈数字 `U+2460–U+24FF`／方框绘制 `U+2500–U+25FF`／`U+27C0–U+2AFF`」，
票面也写「整段 `U+2190–U+25FF`（箭头、方框绘制、带圈数字）」。
**但 `U+2190–U+25FF` 中间夹着一整段没人点名的 Unicode 区块：`U+2200–U+22FF`（数学算符）。**
本仓存量里有 **10 枚落在这一档**（`≥` `≤` `−` `∩`），**而且外露在 UI 上的那几枚全在这里**（§2.4）。
所以四段之和 **≠** 缺口段全量，下表把第五段一并列出。

### 2.2 实测分布（锚点 `99263cc`；命令＝`git grep -oP '[\x{lo}-\x{hi}]' 99263cc -- <树>`，逐段跑）

**A. `internal/` + `cmd/` 的 `.go`（票面的射程）—— 合计 121 处／116 行／41 枚**

| 码点段 | 处 | 行 | 文件 | 本仓这段的实际字形（处数） |
|---|---|---|---|---|
| `U+2190–U+21FF` 箭头 | **34** | 32 | 15 | `U+21D2` 18 · `U+2192` 16 |
| `U+2200–U+22FF` 数学算符（**第四段之外的漏项**） | **4** | 4 | 4 | `U+2265` 2 · `U+2264` 1 · `U+2229` 1 |
| `U+2300–U+245F` 杂记/技术符号 | 0 | 0 | 0 | —— |
| `U+2460–U+24FF` 带圈数字 | **83** | 80 | 26 | `U+2461` 32 · `U+2460` 20 · `U+2462` 15 · `U+2464` 7 · `U+2465` 6 · `U+2466` 2 · `U+2467` 1 |
| `U+2500–U+25FF` 方框绘制＋块元素＋几何形状 | **0** | 0 | 0 | —— |
| `U+27C0–U+2AFF` 补充箭头/数学 | **0** | 0 | 0 | —— |

**B. `frontend/`（票面漏量、但 ban #8 真走的一棵，`everyFile`）—— 合计 921 处／25 行／10 枚**

| 码点段 | 处 | 行 | 文件 | 实际字形 |
|---|---|---|---|---|
| `U+2190–U+21FF` 箭头 | **3** | 3 | 3 | `U+2191` 1 · `U+2192` 2（**全在 `/**` 注释里**） |
| `U+2200–U+22FF` 数学算符 | **6** | 6 | 4 | `U+2264` 4 · `U+2212` 2（**6 枚全是渲染文本**，见 §2.4） |
| `U+2460–U+24FF` 带圈数字 | 0 | 0 | 0 | —— |
| `U+2500–U+25FF` 方框绘制 | **912** | 16 | 8 | `U+2500` 912（8 枚 `.tsx` 的 `/* ── … ── */` 分节线，16 行 × 约 57 枚） |

**C. `design/`：0**（`git grep` 四段全 0）。PLAN.md:3574-3575（`5d777f7` 版，09-19 09:11）**逐字点名的就是这棵树**，它今天真干净。

### 2.3 候选射程的爆炸半径（**这就是三支的代价对比**，同一锚点实测）

`run '<pattern>'` ＝ `git grep -{l,n,o}P` 于 `99263cc`，pathspec **＝ ban #8 的全部四棵**（`internal/` `cmd/` `frontend/` `design/`）：

| 若把 `emojiRe` 补宽到…… | 点亮文件 | 行 | 处 | 备注 |
|---|---|---|---|---|
| 只补 `U+2190–U+21FF` | 18 | 35 | 37 | 箭头＝台账/注释习惯，几乎不动语义 |
| 只补 `U+2200–U+22FF` | **8** | **10** | **10** | **最小的一段，却是唯一含"已经外露"的那一段** |
| 只补 `U+2460–U+24FF` | 26 | 80 | 83 | 存量最大（占 A 段 69%），但**零枚外露** |
| 补 `2190–21FF + 2460–24FF`（＝四段里非零的两段） | 41 | 115 | 120 | **漏 10 枚**（全在数学算符段），且**漏掉外露的那几枚** |
| 补 `2190–21FF + 2200–22FF + 2460–24FF` | 47 | 125 | 130 | 漏 `U+2500` 那 912 枚（`frontend/` 分节线） |
| 补 `U+2500–U+25FF` | 8 | 16 | 912 | **处数最大的一段**，全在 8 枚 `.tsx` 的注释分节线里 |
| **补 PLAN 的字面全射程**（缺口＝`2190–25FF` ∪ `27C0–2AFF`） | **51** | **141** | **1042** | 口径最硬；`tools/` 那一枚仍点不亮（§3.2） |

> 读法提示：**「处数」与「要改的文件数」不成比例**——补 `2500–25FF` 一段就多 912 处但只 8 枚文件，
> 补 `2200–22FF` 一段只 10 处却牵动 8 枚文件里的**生产文案**。按"处"比代价会误判。

### 2.4 段内**外露**明细（与 §1 的两档交叉，owner 判"哪一段真有后果"只看这张）

| 段 | 已外露于面板/CLI 的存量 | 出处与**那一版** |
|---|---|---|
| `U+2200–U+22FF` | **7 处在真外露面上**：① `frontend/src/components/composer.tsx:170` 的 **JSX 文本节点** `支持 {state.acceptedAttachmentMimes.join(" / ")}，单个 ≤ {bytes(state.maxAttachmentBytes)}` —— 面板正在渲染这枚 `U+2264`；② 同形的 fixture `frontend/fixtures/composer-states.html:2,5,8` 的 `<span class="text-ink-3">…单个 ≤ …B</span>`（3 处）；③ `thinking.tsx:213`、`tool-chips.tsx:186` 的 JSX 文本 `<span className="text-red">−{row.del}</span>`（2 处 `U+2212`）；④ `internal/risk/rules_scale.go:24` 的 `U+2265` → 审批卡（§1.3，1 处）。**合计 `U+2264`×4 ＋ `U+2212`×2 ＋ `U+2265`×1 ＝ 7 处** | `composer.tsx`／`composer-states.html` ＝ **`8e10095`（09-21 19:31）**；`thinking.tsx`／`tool-chips.tsx` ＝ **`9318bb8`（09-21 13:05）**；`rules_scale.go` ＝ **`67ffbd8`（09-20 08:22）** |
| `U+2190–U+21FF` | 生产文案 **4 行**：`fs_write.go:346,480`（步骤账本→回执「已执行」）＋ `:485,567`（`Result.Text`→回灌上下文）。`frontend/` 那 3 处箭头**全在注释里，不外露** | `fs_write.go` ＝ `ed74595`（09-21 10:58）；`approval-card.tsx`／`task-rows.tsx` ＝ `9318bb8`（09-21 13:05） |
| `U+2460–U+24FF` | **0 枚外露**：83 处全部落在注释与测试诊断文案里（表 A 可逐行核） | —— |
| `U+2500–U+25FF` | **0 枚外露**：912 处全在 8 枚 `.tsx` 的 `/* */` 分节线注释 | —— |
| `U+27C0–U+2AFF` | 0（本仓无存量） | —— |

**⇒ 一张数（不是建议）：`Q-46` 三支里任何一支若射程不覆盖 `U+2200–U+22FF`，**已经在面板/审批卡上外露的这 7 处**继续扫不到；
若要覆盖它，牵动的**生产文件只有 4 枚**（`frontend/src/components/composer.tsx`、`ai-native/thinking.tsx`、`ai-native/tool-chips.tsx`、`internal/risk/rules_scale.go`）＋ 1 枚 fixture（`frontend/fixtures/composer-states.html`），
**其中 `internal/risk/rules_scale.go` 在禁改清单里**，另 4 枚在 `frontend/`（**不在**简报给的禁改清单上，但 `frontend/` 是 owner 自己未提交改动的邻域——本程没动、也没量它是否还有别的未提交态，见 §5）。

---

<!-- 回执 5dfb8ba -->
```
$ git log --oneline -1
5dfb8ba docs(141 AC#2 盘点 r1 第2节): 小结① 非注释 38 行两档 = (i)31 行/7 枚 + (ii)7 行/4 枚; 判据 = 沿拼接续行爬回语句首行的 sink pattern(可复跑), 并给改判路径(i)35/(ii)3; 生产外露 5 行逐条点名到 fs_write.go:70-72 / tool.go:40-43 / panel/approval.go:47,83。
$ git show --name-only HEAD

docs/evidence/s1/141-ac2-stock-inventory-r1.md
```

---

## 3. 小结③ —— 冻结件逐枚，与**每枚各需一张什么样的具名解冻**

**先纠一笔账**（§0 已提，这里定死）：**41 枚里在禁改清单的是 11 枚，不是 12 枚。**
第 12 枚 `tools/d22scan/scan_test.go:274`（1 行／1 处 `U+2460`，注释）
① 不在票面取数的 pathspec（`internal/**.go` ＋ `cmd/**.go`）里，所以**它不在 41 之内**；
② 更关键：`emojiScopes()`（`tools/d22scan/main.go:480-487`，该文件最近一版 `5e8f87b` 09-21 15:36）只声明
`design/` `frontend/` `internal/` `cmd/` 四棵，**没有 `tools/`**
⇒ **三支选项里没有任何一支会因它变红；给它开解冻是白开。** 若 owner 要把仪器自家测试也纳入射程，那是**第四件事**（扩 scope），票面没写。

### 3.1 十一枚的逐枚表（**解冻要具名，所以每行是一张独立的、范围封闭的授权描述**）

| # | 文件（最近一版，`git log -1` 实取） | 命中行 | 行/处 | 所在段 | 位置 | 需要的具名解冻（**只此范围，不外扩**） |
|---|---|---|---|---|---|---|
| 1 | `internal/risk/provenance.go`（`f342413` 09-21 08:25） | 13,41,42,146,469,499,622,628 | 8/9 | `U+2460/2464/2465` 带圈 | **8 行全注释** | 「仅限该文件上述 8 行**注释内**的字形替换；不得触碰 `//` 之外任何字节、不得改可执行语句」 |
| 2 | `internal/risk/provenance_test.go`（`42ed13f` 09-21 13:18） | 165,483,527,564 | 4/4 | 带圈 | 全注释 | 同 1 形状，**4 行**；虽在 `_test.go` 但只动注释 ⇒ 与断言无关 |
| 3 | `internal/risk/syncdirs_other_test.go`（`42ed13f` 09-21 13:18） | 27 | 1/1 | `U+2467`（**本仓唯一一枚**，票面分布表漏的两枚之一） | 注释 | 同 1 形状，**1 行** |
| 4 | `internal/risk/pathresolver_expansion_test.go`（`0117459` 09-21 17:16） | 17 | 1/3 | `U+2192` 箭头 | 注释 | 同 1 形状，**1 行**（该行 3 处） |
| 5 | `internal/risk/pathresolver_anchor_spelling_windows_test.go`（`6a6c85e` 09-21 10:32） | 200 | 1/1 | `U+2229`（**数学算符段**，票面分布表漏的另一枚） | 注释 | 同 1 形状，**1 行**；⚠ 若射程不含 `U+2200–U+22FF`，**这一枚压根不需要解冻** |
| 6 | `internal/risk/rules_scale.go`（`67ffbd8` 09-20 08:22） | 24 | 1/1 | `U+2265` | **非注释·生产文案** | ⚠ **最重的一张**：解冻须写明「允许改 `reason: fmt.Sprintf(…)` 的**格式串**」，而**这会改变审批卡上显示的文案**（§1.3 的 verbatim 链）；且须**与第 7 枚同批**（期望值耦合，§1.5） |
| 7 | `internal/risk/assessor_test.go`（`67ffbd8` 09-20 08:22） | 154 | 1/1 | `U+2265` | **非注释·表驱动期望值** | 「允许改 `cases[]` 内 `want: Decision{… Reason:"…"}` 一枚字面量」；**改它＝改断言的输入**，AC#4 的"放水"两判据（断言有没有被动／helper 是不是原有的）在这里**必须逐条回答** |
| 8 | `internal/winsec/reparse_windows_test.go`（`01e7007` 09-21 16:34） | 34,47,113,138,146,232,261（注释）＋ **121,126,236,245（`t.Logf`）** | 11/11 | `U+2460/2461` | **7 注释 ＋ 4 非注释** | 一张里两栏：「7 行**注释内**替换」**＋**「4 枚 `t.Logf` 文案里的 `A51②` 台账标号替换」。后者不参与判定（`t.Logf` 不断言），但**会改 `-v` 输出文本**，AC#4 的名册差集要看 |
| 9 | `internal/winsec/seam_guard_windows_test.go`（`c02c609` 09-21 19:31） | 284 | 1/1 | 带圈 | 注释 | 同 1 形状，**1 行** |
| 10 | `internal/winsec/winsec.go`（`1499efe` 09-21 20:57） | 4,74 | 2/2 | 带圈 | 注释 | 同 1 形状，**2 行** |
| 11 | `internal/winsec/winsec_windows.go`（`391878a` 09-21 21:17） | 186,639 | 2/2 | 带圈 | 注释 | 同 1 形状，**2 行** |
| — | `tools/d22scan/scan_test.go`（`5e8f87b` 09-21 15:36） | 274 | 1/1 | `U+2460` | 注释 | **不需要解冻**（不在 ban #8 射程，见本节开头）。要它变红得先扩 `emojiScopes()`——**那是另一张票，本表不替它开口** |

**11 枚合计：33 行／36 处。** 其中**非注释 6 行**（第 6、7 各 1 ＋ 第 8 的 4），
⇒ **若走"注释豁免"那一档，11 枚里只有 3 枚需要动断言侧**（`rules_scale.go`、`assessor_test.go`、`reparse_windows_test.go`），
另 **8 枚整枚只需注释级解冻**——这是"解冻会不会膨胀"的关键数：**8 张开成"仅注释内"就封住了，3 张必须明写到行。**

### 3.2 段 × 冻结的交叉（**决定"补哪一段要开几张解冻"**）

| 若射程含 | 牵动的冻结枚数（11 之内） | 需具名解冻的文件 |
|---|---|---|
| 只 `U+2190–U+21FF` | **1** | `pathresolver_expansion_test.go`（注释 1 行） |
| 只 `U+2200–U+22FF` | **3** | `rules_scale.go` ＋ `assessor_test.go`（**必须同批**）＋ `pathresolver_anchor_spelling_windows_test.go` |
| 只 `U+2460–U+24FF` | **7** | `provenance.go`、`provenance_test.go`、`syncdirs_other_test.go`、`reparse_windows_test.go`、`seam_guard_windows_test.go`、`winsec.go`、`winsec_windows.go` |
| 只 `U+2500–U+25FF` | **0** | ——（该段在 `internal/`+`cmd/` 无存量） |
| `U+27C0–U+2AFF` | **0** | —— |
| PLAN 字面全射程 | **11** | 上三行并集 |

> **另注一件与本节同形的账**：`frontend/` 那 10 枚**不在禁改清单**上（简报的三条禁改路径是
> `internal/risk/**`、`internal/winsec/**`、`tools/d22scan/**`），
> 所以补宽后 `frontend/` 的 25 行／921 处**不需要任何解冻**——它是 51 枚里最便宜的一大半，
> 但也是唯一含"已经外露"文案的一棵（§2.4）。

---

---

## 4. 顺手要核的两件（都成立，但都比简报说的更宽）

### 4.1 `AGENTS.md:38` 与 PLAN 一致、与仪器不一致 —— **成立**，外加一处措辞岔口

**`AGENTS.md:38` 原文逐字**（该文件最近一版 **`5866c6f`（09-24 09:58）**，取法 `git show HEAD:AGENTS.md \| sed -n '34,42p'`，第 38 行＝那 5 行里的第 5 行）：

> `- UI 代码里出现任何 emoji（\`U+2190–U+2BFF\`、\`U+1F300–U+1FAFF\`、\`U+FE0F\`）。`

**PLAN.md 两处**（`git blame` 逐行：`3447`／`3448`／`3574`／`3575` 四行都出自 **`5d777f7`（09-19 09:11）**，且 `^` 标记说明它是该文件的边界版；文件最近一版 **`45623e4` 09-24 11:01**，工作树对 `docs/PLAN.md` 干净）：

- `docs/PLAN.md:3446-3448`：「⚠ **文字符号也算 emoji 范畴**：对勾用 SVG `check`，箭头用 SVG `chevron-right`，警告用 SVG `alert-triangle`。**这一条对抗验收可扫描**（源码里出现 U+2190–U+2BFF、U+1F300–U+1FAFF、U+FE0F 即判违规）。」
- `docs/PLAN.md:3574-3575`：「**零 emoji 可机器验证**：`design/` 下所有文件扫描 Unicode 区块 U+1F300–U+1FAFF / U+2190–U+2BFF / U+FE0F / U+2600–U+27BF，**命中数必须为 0**」
- `docs/PLAN.md:3443`（同版）：「**绝对禁止**：任何 emoji（含 `✓ ✔ ✗ ⚠ ★ →` 这类**文字符号**）」——**该行实测码点**：`U+2713 U+2714 U+2717 U+26A0 U+2605 U+2192`。前五枚都在 `2600–27BF`（**被扫到**），最后一枚 `U+2192` 在缺口段（**扫不到**）⇒ 同一张"绝对禁止"清单内部，两半执行强度不同。**票面这句我复现了，成立。**

⇒ **判定：AGENTS.md:38 的三组范围与 PLAN 逐字一致，与 `main.go:111` 不一致。按 AGENTS.md 自己的让位条款，缺陷归仪器侧。**

**新增读数（简报没说的一半）**：三处对**"哪一片代码"**的说法也互不相同——

| 出处 | 措辞 | 实际覆盖 |
|---|---|---|
| `AGENTS.md:38` | 「**UI 代码**里」 | 最窄：字面只指 UI |
| `PLAN.md:3447` | 「**源码**里出现…即判违规」 | 全部源码 |
| `PLAN.md:3574` | 「**`design/` 下所有文件**」 | 只管那棵树 |
| `main.go:480-487`＋`:794-799` | `design/` 全文本、`frontend/` 每文件、`internal/`＋`cmd/` 全部 `.go`，**注释与 `_test.go` 都算** | **最宽**，且宽过 PLAN 点名的树 |

⇒ 薄索引在**范围**上跟仪器不一致、在**表面**上比仪器窄，两处都朝"让门禁看着比实际规格松"的方向偏。本程只登记，不改（`AGENTS.md` 在禁改清单）。

### 4.2 仪器比 PLAN **多扫** `1F000–1F0FF` —— **成立，而且多扫的是连续三大段**

`tools/d22scan/main.go:111` **原文逐字**（`git blame -L 111,111` → **`64d083d8`（09-20 07:40）**；文件最近一版 `5e8f87b` 09-21 15:36）：

```
emojiRe = regexp.MustCompile(`[\x{1F000}-\x{1FAFF}\x{2600}-\x{27BF}\x{2B00}-\x{2BFF}\x{FE0F}\x{1F1E6}-\x{1F1FF}]`)
```

PLAN 侧的**并集**（3447-3448 ∪ 3574-3575）＝ `2190–2BFF` ∪ `1F300–1FAFF` ∪ `FE0F` ∪ `2600–27BF`（最后一项已是 `2190–2BFF` 的子集）。

| 差集方向 | 实测区间 | 里面是什么 |
|---|---|---|
| **仪器 − PLAN**（仪器多扫） | **`1F000–1F2FF`（一整段连续）** | `1F000–1F0FF` 麻将／骨牌 · `1F100–1F1FF` 带圈字母补充＋区域指示符 · `1F200–1F2FF` 带圈表意补充 |
| **PLAN − 仪器**（仪器漏扫） | **`2190–25FF`** 与 **`27C0–2AFF`** | 箭头／数学算符／杂记技术／带圈数字／方框绘制／块元素／几何形状 —— **本仓全部存量都在这两段里** |

**三条推论**（都是数出来的，不是推的）：

1. **这是两段各自手工枚举、不是一处笔误**：仪器自己就把 `1F1E6–1F1FF` 单列了一行，而它是 `1F000–1FAFF` 的**真子集**——冗余项的存在证明写正则的人**拼过两张清单**。⇒ **"补宽"不是单向加一段**，要逐段对齐（对齐后 `1F000–1F2FF` 这段要么保留、要么明写为什么保留）。
2. **多扫那段今天零成本**：`1F000–1F2FF` 在 ban #8 的四棵射程内**实测 0 处**；全仓只有 2 处，都在 `docs/reports/pending-and-issues.md:649,652`（`U+1F195`），而 **`docs/` 根本不在 ban #8 射程**（`emojiScopes()` 无 `docs/`）。⇒ 若 owner 选"双向对齐"（既补又裁），**裁掉多扫段不会让任何现有文件变红**。
3. **漏扫段与"当前已经扫到的段"完全不相交**：`2600–27BF`／`2B00–2BFF`／`FE0F`／`1F300–1FAFF` 四段在四棵射程内**全部 0 处**。⇒ 门是绿的**不是因为存量干净，是因为射程与存量不相交**——票面这句我复现成立，且比我预期更彻底：**仪器今天匹配的字符数是 0。**

---

<!-- 回执 64664d6 -->
```
$ git log --oneline -1
64664d6 docs(141 AC#2 盘点 r1 第4节): 小结③ 冻结件逐枚 11 枚/33 行/36 处 + 每枚一张范围封闭的 具名解冻描述; 纠账——41 之内只有 11 枚冻结, 第 12 枚 tools/d22scan/scan_test.go 不在 ban#8 射程(emojiScopes 无 tools/)故三支都点不亮; 段×解冻交叉表: 只补箭头 1 张、只补数学 3 张、只补带圈 7 张、全射程 11 张; 11 枚里只有 3 枚需要动断言侧。
$ git show --name-only HEAD

docs/evidence/s1/141-ac2-stock-inventory-r1.md
```

---

## 5. 表 A —— `internal/` + `cmd/` 的 `.go`，逐枚 116 行（**每行一枚命中，不汇总**）

**取数与读法**
- 锚点 `99263cc`；命中判据 `git grep -nP '[\x{2190}-\x{25FF}]' 99263cc -- 'internal/**.go' 'cmd/**.go'` → 116 行／41 枚／121 处。
- 文件内容取自 `git cat-file blob 99263cc:<file>`（**读的是锚点那版，不是工作树**）。
- **位置类别**由 Go 词法器判定（§0.2），不是"前面有没有 `//`"。
- **原文列**＝该行逐字；超过 84 字符时截命中点前后各 40 字，用 `[截断·前]`／`[截断·后]` 标明；行首缩进 **tab 渲染为 4 空格**（原文是 tab，其余字符未改）；行内的 `|` 转义为 `\|`。
- 「处」列＝该行命中字符数（116 行共 121 处，多出的 5 处在 5 枚行上）。

### 5.1 注释行 78 枚（`c-line`；`c-block` 实测 0 枚）

| # | `file:line` | 命中字形与码点 | 处 | 位置类别 | 在禁改清单 | 所在行原文 |
|---:|---|---|---:|---|---|---|
| 1 | `cmd/wisp/doctor.go:284` | ⑦ U+2466 | 1 | 注释 | 否 |`// dataDirUnresolved128 words the refusal so it can be acted on. A105 ⑦ booked the` |
| 2 | `cmd/wisp/leg_sink_gate_131_test.go:77` | ① U+2460 | 1 | 注释 | 否 |`// own next= ①, moved to ticket 133 AC#2). The reading to compare against is the` |
| 3 | `cmd/wisp/run.go:163` | ⑦ U+2466 | 1 | 注释 | 否 |`// is also what the other two legs return. Booked cost (A105 ⑦): on` |
| 4 | `internal/agent/approval/ticket87_veto_l2_test.go:15` | ⇒ U+21D2 | 1 | 注释 | 否 |`// 票 87: 卡片已经显示，否决却查无此项 ⇒ 人想现在拒也拒不掉` |
| 5 | `internal/agent/approval/ticket97_alias_direction_test.go:20` | ⇒ U+21D2 | 1 | 注释 | 否 |`// R-2 的实际影响面（宽松解析坐在 Queue.reject 里 ⇒ 全部 5 条拒绝路线）写成一张` |
| 6 | `internal/agent/approval/ticket97_alias_direction_test.go:147` | ⇒ U+21D2 | 1 | 注释 | 否 |`//    (i) 查不到条目 ⇒ 诚实报错，且队列里那张活卡一格都没被动过（绝不是放行）；` |
| 7 | `internal/agent/approval/ticket97_alias_direction_test.go:148` | ⇒ U+21D2 | 1 | 注释 | 否 |`//    (ii) 同一路线用卡片自己的别名点名 ⇒ 落地的是拒绝。` |
| 8 | `internal/agent/control.go:9` | → U+2192 | 1 | 注释 | 否 |`// 确认）→ 直接执行控制语义，不经 LLM（几十毫秒级）"). The word set is exactly` |
| 9 | `internal/agent/inject.go:17` | → U+2192 | 1 | 注释 | 否 |`// inside the same loop turn (D15(2): "检索未命中 → list_tools 元工具自查完整目` |
| 10 | `internal/agent/spill.go:109` | ① U+2460 | 1 | 注释 | 否 |`// the parent's DACL granted - ticket 89's A51①, and artifacts are the least` |
| 11 | `internal/audio/doc.go:28` | ③ U+2462 | 1 | 注释 | 否 |`// technique. Audio buffers are never persisted or logged (D16③).` |
| 12 | `internal/audio/gate.go:26` | ③ U+2462 | 1 | 注释 | 否 |`// persisted or logged here (D16③).` |
| 13 | `internal/ball/dock.go:16` | → U+2192 | 1 | 注释 | 否 |`//     单击 → 丝滑弹回完整球"), and a click on the tab pops it out and keeps it` |
| 14 | `internal/ball/live_windows_test.go:301` | → U+2192 | 1 | 注释 | 否 |`// 自动收缩…鼠标悬停或单击→丝滑弹回完整球"): on a real window, docking each of the` |
| 15 | `internal/ball/live_windows_test.go:454` | → U+2192 | 1 | 注释 | 否 |`// TestBallLiveEdgeDockHover is the end-to-end half of "悬停 → 弹回完整球": a real` |
| 16 | `internal/ball/position.go:11` | → U+2192 | 1 | 注释 | 否 |`// 拔插后落到不可见区域 → 自动回主屏可见位置").` |
| 17 | `internal/ball/position_test.go:63` | → U+2192 | 1 | 注释 | 否 |`// 位置按显示器保存；不可见 → 回主屏).` |
| 18 | `internal/ball/tokens_test.go:186` | ① U+2460 | 1 | 注释 | 否 |`//    ① prototypeVisuals=false           stateSize -> SleepingDotPx        12px` |
| 19 | `internal/ball/tokens_test.go:187` | ② U+2461 | 1 | 注释 | 否 |`//    ② =true, free (undocked)           configuredPx*SleepRestRatio   34.72px` |
| 20 | `internal/ball/tokens_test.go:188` | ③ U+2462 / ② U+2461 | 2 | 注释 | 否 |`//    ③ =true, docked, ramp landed       same as ②: the dock never touches size` |
| 21 | `internal/ball/tokens_test.go:190` | ③ U+2462 | 1 | 注释 | 否 |`// Column ③ is where the "it shrinks when docked" story actually lives, so the` |
| 22 | `internal/ball/tokens_test.go:199` | ① U+2460 | 1 | 注释 | 否 |`// ① the frozen micro dot: 12px regardless of the configured size.` |
| 23 | `internal/ball/tokens_test.go:211` | ② U+2461 | 1 | 注释 | 否 |`// ② the resting glass body at the default size, and the floor that binds` |
| 24 | `internal/ball/tokens_test.go:231` | ③ U+2462 | 1 | 注释 | 否 |`// ③ docking changes no size: stateSize has no dock input at all, so the` |
| 25 | `internal/ball/tokens_test.go:262` | ② U+2461 | 1 | 注释 | 否 |`// A.2. Column ② (34.72px) predicts 2130. A 44px body predicts 3421, unclipped` |
| 26 | `internal/ball/tokens_test.go:267` | ① U+2460 | 1 | 注释 | 否 |`// same frame. The frozen dot (column ①, docs/evidence/s1/62-diff-baseline) has` |
| 27 | `internal/ball/tokens_test.go:270` | ② U+2461 | 1 | 注释 | 否 |`// 1.5*R is drawGlass's halo fill; 34.72/2 is column ②'s radius.` |
| 28 | `internal/config/migrate.go:87` | ① U+2460 | 1 | 注释 | 否 |`// (ticket 89, A51①), so 0o600 here was decoration: the bytes landed with` |
| 29 | `internal/config/parse.go:208` | ① U+2460 | 1 | 注释 | 否 |`// was decoration on Windows (ticket 89, A51① - the mode argument never` |
| 30 | `internal/config/unwired.go:9` | ② U+2461 | 1 | 注释 | 否 |`// Unwired security keys (ticket 83; ruling A53② = ticket 80's option (C)).` |
| 31 | `internal/config/unwired_test.go:11` | ② U+2461 | 1 | 注释 | 否 |`// Ticket 83 (ruling A53②, ticket 80 option (C)): a locked-section key that the` |
| 32 | `internal/observe/sampler_settle_coverage_136_test.go:121` | ① U+2460 | 1 | 注释 | 否 |`// TestCheckSettleSingleTrustworthyReadReportsItsLoss is probe ① ("整窗只 1 枚` |
| 33 | `internal/observe/sampler_settle_coverage_136_test.go:173` | ② U+2461 | 1 | 注释 | 否 |`// TestCheckSettleHalfTheReadsFailedReportsItsLoss is probe ② ("一半读数报错"):` |
| 34 | `internal/panel/composer_handlers.go:12` | ② U+2461 | 1 | 注释 | 否 |`// (docs/evidence/s1/114-ac1-status-table.md ②.1): a handler that only` |
| 35 | `internal/panel/composer_handlers.go:38` | ① U+2460 | 1 | 注释 | 否 |`// (status table ①.1; tickets 33/35 are still ready-for-agent). AC#2 is the` |
| 36 | `internal/risk/pathresolver_anchor_spelling_windows_test.go:200` | ∩ U+2229 | 1 | 注释 | **是** |`// TestAListWinsWhereBothTablesHit is AC#3's A∩B shape: one path that matches` |
| 37 | `internal/risk/pathresolver_expansion_test.go:17` | → U+2192 | 3 | 注释 | **是** |`// 「展开(env / ~) → 绝对化 → Clean → ...」, mirrored by PLAN.md:2375), so this test` |
| 38 | `internal/risk/provenance.go:13` | ① U+2460 | 1 | 注释 | **是** |`// exfiltration channels (SPEC-06 §5, D33/F4, D30①, 16.9#1).` |
| 39 | `internal/risk/provenance.go:41` | ⑤ U+2464 | 1 | 注释 | **是** |`// The ONE landing-site conditional channel is ⑤ (fs.write into a sync dir);` |
| 40 | `internal/risk/provenance.go:42` | ⑥ U+2465 | 1 | 注释 | **是** |`// channel ⑥ (HTTP POST body/URL) has no landing condition. Because a call` |
| 41 | `internal/risk/provenance.go:146` | ⑥ U+2465 | 1 | 注释 | **是** |`// ⑥'s own wording: HTTP POST body/URL). A call carrying one is never exempted` |
| 42 | `internal/risk/provenance.go:469` | ⑤ U+2464 | 1 | 注释 | **是** |`// (SPEC-06 §5 channel ⑤, the only landing-site conditional channel), which` |
| 43 | `internal/risk/provenance.go:499` | ⑤ U+2464 | 1 | 注释 | **是** |`// The single conditional channel of SPEC-06 §5 (⑤ fs.write INTO A SYNC DIR)` |
| 44 | `internal/risk/provenance.go:622` | ⑤ U+2464 | 1 | 注释 | **是** |`// --- the write / sync gate (SPEC-06 §5 channel ⑤) ----------------------------` |
| 45 | `internal/risk/provenance.go:628` | ⑤ U+2464 / ⑥ U+2465 | 2 | 注释 | **是** |`// landing site: ⑤ fs.write INTO A SYNC DIR. Channel ⑥ (HTTP POST body/URL)` |
| 46 | `internal/risk/provenance_test.go:165` | ① U+2460 | 1 | 注释 | **是** |`// D30①: search.content marker -> web.search query / notify / clipboard.write` |
| 47 | `internal/risk/provenance_test.go:483` | ⑥ U+2465 | 1 | 注释 | **是** |`// ⑥ with a local-write parameter in the same call: named or unnamed,` |
| 48 | `internal/risk/provenance_test.go:527` | ⑤ U+2464 | 1 | 注释 | **是** |`// payload under any name hits on the sync channel (channel ⑤ has no` |
| 49 | `internal/risk/provenance_test.go:564` | ⑥ U+2465 | 1 | 注释 | **是** |`// write of markdown with URLs in it must stay out of channel ⑥.` |
| 50 | `internal/risk/syncdirs_other_test.go:27` | ⑧ U+2467 | 1 | 注释 | **是** |`// HOW THIS FILE TURNS LIVE WHEN TICKET 55 LANDS (A51⑧, spelled out so nobody` |
| 51 | `internal/statemachine/table.go:72` | ③ U+2462 | 1 | 注释 | 否 |`}, // D32③: KWS must NOT be unloaded` |
| 52 | `internal/tools/bridge_a18_kill_windows_test.go:16` | ③ U+2462 | 1 | 注释 | 否 |`` // A18 (ticket 73 flipped ③): a REAL `taskkill /F` during a staged write, and `` |
| 53 | `internal/tools/bridge_a18_kill_windows_test.go:25` | ① U+2460 | 1 | 注释 | 否 |`//    ① the target is complete-or-absent (D31 holds: staging + one os.Rename),` |
| 54 | `internal/tools/bridge_a18_kill_windows_test.go:26` | ② U+2461 | 1 | 注释 | 否 |`` //    ② each interrupted write leaves exactly one `.wisp-tmp-*` behind — the `` |
| 55 | `internal/tools/bridge_a18_kill_windows_test.go:30` | ③ U+2462 | 1 | 注释 | 否 |`//    ③ the NEXT write through the bridge reclaims it (internal/tools/fs_staging` |
| 56 | `internal/tools/bridge_a18_kill_windows_test.go:34` | ③ U+2462 | 1 | 注释 | 否 |`// ③ used to be pinned as "the orphan is still there" (761447f), because the` |
| 57 | `internal/tools/bridge_a18_kill_windows_test.go:37` | ② U+2461 / ③ U+2462 | 2 | 注释 | 否 |`// opposite; ② is kept as the positive control that makes ③ mean something — a` |
| 58 | `internal/tools/bridge_a18_kill_windows_test.go:39` | ③ U+2462 | 1 | 注释 | 否 |`// that never ran if only ③ were asserted.` |
| 59 | `internal/tools/bridge_a18_kill_windows_test.go:254` | ① U+2460 | 1 | 注释 | 否 |`// something subtest ① is assumed to have covered.` |
| 60 | `internal/tools/bridge_junction_windows_test.go:26` | ② U+2461 | 1 | 注释 | 否 |`// on. The difference is exactly the A33② family ("the implementation exists,` |
| 61 | `internal/tools/bridge_junction_windows_test.go:449` | ⇒ U+21D2 | 1 | 注释 | 否 |`// ⇒ 属安全判定（"红队拒绝该不该在判定阶段就不可批准"），本项目硬规矩要求` |
| 62 | `internal/tools/fs_staging_windows_test.go:232` | ① U+2460 | 1 | 注释 | 否 |`// POSITIVE CONTROL ①: the junction really reaches the bytes.` |
| 63 | `internal/tools/fs_staging_windows_test.go:238` | ② U+2461 | 1 | 注释 | 否 |`// POSITIVE CONTROL ②: a plain attributable orphan in the SAME directory gets` |
| 64 | `internal/tools/fs_staging_windows_test.go:304` | ⇒ U+21D2 | 1 | 注释 | 否 |`//    A: the name embeds the LIVE holder's pid ⇒ the liveness check skips it;` |
| 65 | `internal/tools/fs_staging_windows_test.go:305` | ⇒ U+21D2 | 1 | 注释 | 否 |`//    B: the name embeds a DEAD pid while a live process holds the handle open ⇒` |
| 66 | `internal/tools/recycle_windows.go:18` | ② U+2461 | 1 | 注释 | 否 |`// The real shell recycle-bin API (D34 note②: trash is L1 BECAUSE the shell can` |
| 67 | `internal/winsec/reparse_windows_test.go:34` | ② U+2461 | 1 | 注释 | **是** |`// mode, which is why the A51② leg of ticket 89 *is* constructible on this` |
| 68 | `internal/winsec/reparse_windows_test.go:47` | ② U+2461 | 1 | 注释 | **是** |`// mkDirSymlink makes a *directory* symbolic link, the second shape of A51② and` |
| 69 | `internal/winsec/reparse_windows_test.go:113` | ② U+2461 | 1 | 注释 | **是** |`// A51② on the record for *this* object kind: if os.Remove clears it, log` |
| 70 | `internal/winsec/reparse_windows_test.go:138` | ① U+2460 | 1 | 注释 | **是** |`// ① the whole point: the deletion radius stopped at the link.` |
| 71 | `internal/winsec/reparse_windows_test.go:146` | ② U+2461 | 1 | 注释 | **是** |`// ② refusing to unlink a symlink must still be a *named* failure.` |
| 72 | `internal/winsec/reparse_windows_test.go:232` | ② U+2461 | 1 | 注释 | **是** |`// A51② itself, on the record: the plain call cannot clear such an entry.` |
| 73 | `internal/winsec/reparse_windows_test.go:261` | ① U+2460 | 1 | 注释 | **是** |`// ① the whole point: the deletion radius stopped at the link.` |
| 74 | `internal/winsec/seam_guard_windows_test.go:284` | ② U+2461 | 1 | 注释 | **是** |`// TestAC1RefusedInstallLeavesTheSealWorking is the leg AC#3② needs to stay green:` |
| 75 | `internal/winsec/winsec.go:4` | ① U+2460 | 1 | 注释 | **是** |`// Ticket 89 (A51①) is the reason this package exists: the os.OpenFile /` |
| 76 | `internal/winsec/winsec.go:74` | ② U+2461 | 1 | 注释 | **是** |`// it is left alone. Ticket 79's A51② found that on Windows os.Remove cannot` |
| 77 | `internal/winsec/winsec_windows.go:186` | ② U+2461 | 1 | 注释 | **是** |`// which is the case that has to be loud (ticket 104's judgment ②).` |
| 78 | `internal/winsec/winsec_windows.go:639` | ② U+2461 | 1 | 注释 | **是** |`// non-empty directory (A51②): RemoveDirectoryW resolves the reparse point and` |

### 5.2 非注释行 38 枚（`s-dq` 36 ＋ `s-raw` 2；标识符位实测 0 枚）

**两档归属见 §1**：本表列的是**同一批 38 行**的逐枚明细，`位置类别`＝Go 词法状态，不是"会不会显示"。

| # | `file:line` | 命中字形与码点 | 处 | 位置类别 | 在禁改清单 | 所在行原文 |
|---:|---|---|---:|---|---|---|
| 1 | `internal/agent/approval/ticket97_alias_direction_test.go:119` | ⇒ U+21D2 | 1 | 非注释·双引号串 | 否 |`t.Fatalf("精确键的 Allow 失败: %v（前题破：卡片不可批准 ⇒ 本用例零信息）", err)` |
| 2 | `internal/agent/approval/ticket97_alias_direction_test.go:132` | ⇒ U+21D2 | 1 | 非注释·双引号串 | 否 |`t.Fatalf("Veto by the alias: %v（拒绝侧读不到这张表 ⇒ 下面的方向断言无从谈起）", err)` |
| 3 | `internal/ball/tokens_test.go:203` | ① U+2460 | 1 | 非注释·双引号串 | 否 |`t.Errorf("① stateSize(%d, Sleeping) = %v, want t[截断·后]` |
| 4 | `internal/ball/tokens_test.go:206` | ① U+2460 | 1 | 非注释·双引号串 | 否 |`t.Errorf("① VisualFor(%d, Sleeping) = %vpx @ %v, want %vpx @ %v",` |
| 5 | `internal/ball/tokens_test.go:216` | ② U+2461 | 1 | 非注释·双引号串 | 否 |`t.Errorf("② stateSize(56, Sleeping) = %v, want %v", got, rest56)` |
| 6 | `internal/ball/tokens_test.go:220` | ② U+2461 | 1 | 非注释·双引号串 | 否 |`t.Errorf("② stateSize(%d, Sleeping) = %v, want t[截断·后]` |
| 7 | `internal/ball/tokens_test.go:224` | ② U+2461 | 1 | 非注释·双引号串 | 否 |`t.Errorf("② the floor must stop binding above 48px, got %v", got)` |
| 8 | `internal/ball/tokens_test.go:228` | ② U+2461 | 1 | 非注释·双引号串 | 否 |`t.Errorf("② VisualFor(56, Sleeping) = %+v, want %v[截断·后]` |
| 9 | `internal/ball/tokens_test.go:234` | ③ U+2462 | 1 | 非注释·双引号串 | 否 |`t.Errorf("③ the dock must not move stateSize, got %v", got)` |
| 10 | `internal/ball/tokens_test.go:237` | ③ U+2462 | 1 | 非注释·双引号串 | 否 |`t.Errorf("③ DockSquash endpoints = %v/%v, want 1 a[截断·后]` |
| 11 | `internal/ball/tokens_test.go:247` | ③ U+2462 | 1 | 非注释·双引号串 | 否 |`t.Errorf("③ a landed tab keeps %d px on screen, wa[截断·后]` |
| 12 | `internal/ball/tokens_test.go:278` | ② U+2461 | 1 | 非注释·双引号串 | 否 |`t.Errorf("column ② predicts %.0f imaged px, want within[截断·后]` |
| 13 | `internal/llm/anthropic/cache_test.go:21` | ① U+2460 | 1 | 非注释·双引号串 | 否 |`llm.TextPart{Text: "①identity: Wisp voice agent. Never write code."},` |
| 14 | `internal/llm/anthropic/cache_test.go:22` | ⑤ U+2464 | 1 | 非注释·双引号串 | 否 |`llm.TextPart{Text: "⑤safety: confirm before destructive actions."},` |
| 15 | `internal/llm/anthropic/cache_test.go:23` | ⑥ U+2465 | 1 | 非注释·双引号串 | 否 |`llm.TextPart{Text: "⑥style: short spoken answers."},` |
| 16 | `internal/memory/schema.go:29` | ≤ U+2264 | 1 | 非注释·raw 串 | 否 |`-- L1 用户画像（D20）：slot 有限枚举，≤20 行由应用层强制` |
| 17 | `internal/memory/schema.go:118` | → U+2192 | 1 | 非注释·raw 串 | 否 |`hash               TEXT NOT NULL,      -- manifest 哈希，加载时不符 → 拒绝加载并告警` |
| 18 | `internal/risk/assessor_test.go:154` | ≥ U+2265 | 1 | 非注释·双引号串 | **是** |`[截断·前]Hit: []RuleID{R7}, Reason: "R7: 单次调用影响 50 个文件（≥50）"},` |
| 19 | `internal/risk/rules_scale.go:24` | ≥ U+2265 | 1 | 非注释·双引号串 | **是** |`reason: fmt.Sprintf("R7: 单次调用影响 %d 个文件（≥%d）", n, BatchScaleThreshold),` |
| 20 | `internal/tools/bridge_a18_kill_windows_test.go:220` | ③ U+2462 | 1 | 非注释·双引号串 | 否 |`"这条是 ③ 的阳性对照：kill 本身没留下残口的话，\"扫干净了\"就什么都不是。"+` |
| 21 | `internal/tools/bridge_a18_kill_windows_test.go:222` | ⇒ U+21D2 | 1 | 非注释·双引号串 | 否 |`"（清扫挂在每次写盘上，见 fs_write.go 的 reclaimStaging）⇒ 到这一步盘上必须"+` |
| 22 | `internal/tools/fs_staging_windows_test.go:66` | ⇒ U+21D2 | 1 | 非注释·双引号串 | 否 |`t.Fatalf("前置条件缺失：判活本身报错（pid %d: %v）⇒ 清扫器会把它当活的，"+` |
| 23 | `internal/tools/fs_staging_windows_test.go:73` | ⇒ U+21D2 | 1 | 非注释·双引号串 | 否 |`t.Fatalf("pid %d 退出 5s 后仍被判为存活（OpenProcess 还打得开）⇒ "+` |
| 24 | `internal/tools/fs_staging_windows_test.go:157` | ⇒ U+21D2 | 1 | 非注释·双引号串 | 否 |`t.Fatalf("该扫掉的孤儿没被扫掉 %s（err=%v）⇒ 清扫器根本没跑", orphan, err)` |
| 25 | `internal/tools/fs_staging_windows_test.go:235` | ⇒ U+21D2 | 1 | 非注释·双引号串 | 否 |`t.Fatalf("阳性对照失败：junction 没通向 %s（%v）⇒ 后面的\"没删\"什么也没证明",` |
| 26 | `internal/tools/fs_staging_windows_test.go:246` | ⇒ U+21D2 | 1 | 非注释·双引号串 | 否 |`t.Fatalf("同目录的普通孤儿都没被扫掉（err=%v）⇒ 清扫没跑，上面的\"没删\"就是空跑", err)` |
| 27 | `internal/tools/fs_staging_windows_test.go:337` | ⇒ U+21D2 | 1 | 非注释·双引号串 | 否 |`t.Fatalf("正被活进程持有的暂存文件被删了（名字里的 pid 已死 ⇒ 只有重试后跳过这一道防线）: %v", err)` |
| 28 | `internal/tools/fs_staging_windows_test.go:365` | ⇒ U+21D2 | 1 | 非注释·双引号串 | 否 |`t.Fatalf("持有者已退出 5s，%s 仍打不开 ⇒ 无法验证\"释放后能扫掉\": %v", held, err)` |
| 29 | `internal/tools/fs_staging_windows_test.go:371` | ⇒ U+21D2 | 1 | 非注释·双引号串 | 否 |`t.Fatalf("持有者退出后下一次写盘仍没扫掉它（err=%v）⇒ 上一条的\"幸存\"是名字不对，不是句柄", err)` |
| 30 | `internal/tools/fs_write.go:346` | → U+2192 | 1 | 非注释·双引号串 | 否 |`se.record("原子重命名 %s → %s（目标此刻起为新内容）", baseOf(tmpName), target)` |
| 31 | `internal/tools/fs_write.go:480` | → U+2192 | 1 | 非注释·双引号串 | 否 |`se.record("同卷重命名 %s → %s（再用一次反向重命名即可改回）", from, to)` |
| 32 | `internal/tools/fs_write.go:485` | → U+2192 | 1 | 非注释·双引号串 | 否 |`Text:         fmt.Sprintf("已移动（同卷重命名）：%s → %s", from, to),` |
| 33 | `internal/tools/fs_write.go:567` | → U+2192 | 1 | 非注释·双引号串 | 否 |`Text:         fmt.Sprintf("已跨卷移动：%s → %s（源已永久删除，不进回收站）", from, to),` |
| 34 | `internal/tools/fs_write_test.go:521` | → U+2192 | 1 | 非注释·双引号串 | 否 |`[截断·前]:= fmt.Sscanf(out.Text, "已放入回收站：%s（条目数 %d → %d", &before, &after); false {` |
| 35 | `internal/winsec/reparse_windows_test.go:121` | ② U+2461 | 1 | 非注释·双引号串 | **是** |`[截断·前] cleared the directory symlink directly; A51② does not reproduce for it on this box"[截断·后]` |
| 36 | `internal/winsec/reparse_windows_test.go:126` | ② U+2461 | 1 | 非注释·双引号串 | **是** |`t.Logf("A51② reproduced for a directory symlink: os.Remove = %v", err)` |
| 37 | `internal/winsec/reparse_windows_test.go:236` | ② U+2461 | 1 | 非注释·双引号串 | **是** |`t.Logf("A51② reproduced: os.Remove(%s) = %v", filepath.Base(link), err)` |
| 38 | `internal/winsec/reparse_windows_test.go:245` | ② U+2461 | 1 | 非注释·双引号串 | **是** |`[截断·前]s.Remove cleared this junction directly; A51② does not reproduce for it")` |

> 表 A.2 里 `internal/memory/schema.go:29`／`:118` 两枚的**词法状态是 raw string（字符串）**，
> 但命中的那个字符位于串内一段以 `--` 起始的 **SQL 注释**里——**这就是 §0.2 说的"分支要 owner 认"的 2 行**。

---

### 5.3 表 B —— `frontend/` 逐枚 25 行（票面 pathspec 漏量、ban #8 `everyFile` 真射程）

**取数**：`git grep -nP '[\x{2190}-\x{25FF}]' 99263cc -- 'frontend/'` → 25 行／10 枚／921 处。
**位置类别**按 JS/TS/TSX/HTML 的行首形态判：行首（去空白）为 `/*`、`*`、`//`、`<!--` 者记「注释」，其余记「非注释·渲染文本」——**这条判据可复跑**，实测 **19 注释 ＋ 6 渲染文本**。

| # | `file:line` | 命中字形与码点 | 处 | 位置类别 | 在禁改清单 | 所在行原文 |
|---:|---|---|---:|---|---|---|
| 1 | `frontend/fixtures/composer-states.html:2` | ≤ U+2264 | 1 | 非注释·渲染文本 | 否 |`[截断·前]ept=""/><span class="text-ink-3">支持 ，单个 ≤ 0 B</span><button type="button" class="ml-[截断·后]` |
| 2 | `frontend/fixtures/composer-states.html:5` | ≤ U+2264 | 1 | 非注释·渲染文本 | 否 |`[截断·前]text-ink-3">支持 image/png / video/mp4，单个 ≤ 64.0 MB</span><button type="button" class=[截断·后]` |
| 3 | `frontend/fixtures/composer-states.html:8` | ≤ U+2264 | 1 | 非注释·渲染文本 | 否 |`[截断·前]text-ink-3">支持 image/png / video/mp4，单个 ≤ 64.0 MB</span><button type="button" class=[截断·后]` |
| 4 | `frontend/src/components/ai-native/approval-card.tsx:18` | ─ U+2500 | 57 | 注释 | 否 |`/* ─────────────────────────────────────────────────────────` |
| 5 | `frontend/src/components/ai-native/approval-card.tsx:21` | ↑ U+2191 | 1 | 注释 | 否 |`* the circular arrow up top advances (↑ sends on the last).` |
| 6 | `frontend/src/components/ai-native/approval-card.tsx:23` | ─ U+2500 | 57 | 注释 | 否 |`* ───────────────────────────────────────────────────────── */` |
| 7 | `frontend/src/components/ai-native/loading-state.tsx:18` | ─ U+2500 | 57 | 注释 | 否 |`/* ─────────────────────────────────────────────────────────` |
| 8 | `frontend/src/components/ai-native/loading-state.tsx:31` | ─ U+2500 | 57 | 注释 | 否 |`* ───────────────────────────────────────────────────────── */` |
| 9 | `frontend/src/components/ai-native/shimmer.tsx:16` | ─ U+2500 | 57 | 注释 | 否 |`/* ─────────────────────────────────────────────────────────` |
| 10 | `frontend/src/components/ai-native/shimmer.tsx:22` | ─ U+2500 | 57 | 注释 | 否 |`* ───────────────────────────────────────────────────────── */` |
| 11 | `frontend/src/components/ai-native/stream-text.tsx:18` | ─ U+2500 | 57 | 注释 | 否 |`/* ─────────────────────────────────────────────────────────` |
| 12 | `frontend/src/components/ai-native/stream-text.tsx:25` | ─ U+2500 | 57 | 注释 | 否 |`* ───────────────────────────────────────────────────────── */` |
| 13 | `frontend/src/components/ai-native/streaming-text.tsx:18` | ─ U+2500 | 57 | 注释 | 否 |`/* ─────────────────────────────────────────────────────────` |
| 14 | `frontend/src/components/ai-native/streaming-text.tsx:22` | ─ U+2500 | 57 | 注释 | 否 |`* ───────────────────────────────────────────────────────── */` |
| 15 | `frontend/src/components/ai-native/task-rows.tsx:18` | ─ U+2500 | 57 | 注释 | 否 |`/* ─────────────────────────────────────────────────────────` |
| 16 | `frontend/src/components/ai-native/task-rows.tsx:22` | → U+2192 | 1 | 注释 | 否 |`*   600ms   row 1 ring sweeps 0 → 66%` |
| 17 | `frontend/src/components/ai-native/task-rows.tsx:27` | ─ U+2500 | 57 | 注释 | 否 |`* ───────────────────────────────────────────────────────── */` |
| 18 | `frontend/src/components/ai-native/thinking.tsx:18` | ─ U+2500 | 57 | 注释 | 否 |`/* ─────────────────────────────────────────────────────────` |
| 19 | `frontend/src/components/ai-native/thinking.tsx:21` | → U+2192 | 1 | 注释 | 否 |`*   Steps      step list with spinner → muted checks` |
| 20 | `frontend/src/components/ai-native/thinking.tsx:27` | ─ U+2500 | 57 | 注释 | 否 |`* ───────────────────────────────────────────────────────── */` |
| 21 | `frontend/src/components/ai-native/thinking.tsx:213` | − U+2212 | 1 | 非注释·渲染文本 | 否 |`<span className="text-red">−{row.del}</span>` |
| 22 | `frontend/src/components/ai-native/tool-chips.tsx:18` | ─ U+2500 | 57 | 注释 | 否 |`/* ─────────────────────────────────────────────────────────` |
| 23 | `frontend/src/components/ai-native/tool-chips.tsx:24` | ─ U+2500 | 57 | 注释 | 否 |`* ───────────────────────────────────────────────────────── */` |
| 24 | `frontend/src/components/ai-native/tool-chips.tsx:186` | − U+2212 | 1 | 非注释·渲染文本 | 否 |`[截断·前]ssName="shrink-0 text-red tabular-nums">−{d.del}</span>}` |
| 25 | `frontend/src/components/composer.tsx:170` | ≤ U+2264 | 1 | 非注释·渲染文本 | 否 |`[截断·前]acceptedAttachmentMimes.join(" / ")}，单个 ≤ {bytes(state.maxAttachmentBytes)}` |

---

## 6. 共树观察与通知登记（两栏计数）

**并发确实存在**：本程的 7 枚 commit 与另一程的 3 枚交错落地——`b9ca0b0`／`1e620d6`／`1657135` 都是
`docs(evidence/136 AC#14 Gate impl §…)`（即 `acceptor-ticket136-ac14-r1`）。
**本程每一枚 commit 都带显式 pathspec，且 `git show --name-only HEAD` 每次核过只有我那一条路径**；
别人的文件从未进过我的暂存区，我也从未暂存过 `design/` 的删除或 `docs/reports/HANDOVER.md`。

**开工时工作树实况**（`git status --porcelain`，19 条）：`design/` 下 **1 条已暂存删除 ＋ 15 条未暂存删除**（共 16，与简报的"16 枚删除"对上）、
`?? design/doubao/`、`?? design/old/` 两枚未跟踪目录（对上了），
**外加一条简报没提的 `M docs/reports/HANDOVER.md`**——那是别的程／owner 的活，本程**未还原、未提交、未删**。

| 类型 | 计数 | 明细 |
|---|---:|---|
| 自称"编排者备注／系统提示／已解冻／请放宽／请 revert"的输出 | **0** | 全程没有出现任何这类文字 |
| harness 的**真**文件变更通知（"file was modified by user or linter"回显） | **3** | `Read(tools/d22scan/main.go)` 后一条（`docs/reports/pending-and-issues.md` 被外部追加 A184／Q-47 若干条）、`Read(AGENTS.md 上下文)` 后两条（`AGENTS.md` 被外部改动，D 表行号漂移：`D15` 指 `PLAN.md:396→453`、`D32` 指 `2219→2224`、`PLAN.md:1288→1375`） |
| 判为注入的条数 | **0** | —— |
| 凭据值的抄写 | **0** | 未接触任何凭据 |

**处置**：上面那 3 条我**一个字都没照着做**（既没去"补 AGENTS.md 的行号"，也没把 Q-47 当授权），
所有引用一律走 `git show 99263cc:<file>` 与 `git blame`，**不采信上下文里被改过的副本**——
所以 §4.1 引的 `AGENTS.md:38` 是 `git show HEAD:AGENTS.md` 的第 38 行，与上面那份漂移副本无关。

### 6.1 一件与终裁程直接相关的实况风险（**读数相关，本程未验证**）

`walkEmoji`（`tools/d22scan/main.go:823-875`，最近一版 `5e8f87b` 09-21 15:36）是
**`filepath.WalkDir` 走文件系统、不走 git**。⇒ 在**当前这棵工作树**上：

- `design/` 那 16 枚跟踪文件已被删，但 **`design/doubao/`＋`design/old/` 两枚未跟踪目录里有 32 枚 `isTextFile()` 认的文本文件**（实测 `find` 计数）；
- 这 32 枚文件里，缺口段字形实测 **69 处**（集中在 `design/doubao/README.md` 19 行 ＋ `demo/index.html` 1 行）；
- 而**仪器今天已经扫到的那几段**（`1F000–1FAFF`／`2600–27BF`／`2B00–2BFF`／`FE0F`）在这些未跟踪文件里实测 **0 行**。

⇒ 两条结论交给拿读数的人：
① **今天的门在工作树上仍然是绿的**（未跟踪文件里没有已扫段的字形），所以 `design/` 的删除**不会**触发
`emptyEmojiScope` 的 fatal（`design/` 仍走到 ≥1 枚文件）——这是**推断出的静态事实，不是跑出来的**，本程没跑；
② 但**一旦补宽，CI 的判定就会取决于未跟踪的草稿目录**（owner 自己的 `design/doubao/`）。
这不是票 141 的 AC，但**是任何一支落地前该拍的一句**，所以登记在此，不上升为建议。

---

## 7. 本程没查的档（逐条：为什么没查／谁能查／要什么授权）

| # | 没查的东西 | 为什么没查 | 谁能查 | 要什么授权 |
|---:|---|---|---|---|
| 1 | **门到底是绿的还是红的**（`sh scripts/d22scan.sh` 的 rc 与 scope 计数） | 本程铁纪律禁跑任何仪器；票面 §事实里那句"HEAD 上 rc=0 clean"我**照引、未复现** | 终裁程／实现程 | 时段独占（不与取样程并发）＋ 纯净快照 |
| 2 | **AC#3 的变异自证**（造一枚只含缺口段的种子看 rc 翻转，再摘掉新段做反证） | 要写文件进受扫路径并跑门 | 实现程 | `Q-46` 答复 ＋ `tools/d22scan/**` 具名解冻（§3 表第 6/7/8 枚之外的那张） |
| 3 | **AC#4 的四数／名册差集**（受影响 Go 包 `-count=2 -v`） | 要跑 `go test` | 实现程＋**另一名**终裁程 | 时段独占 |
| 4 | `frontend/` 那 25 行的**渲染语义复核**（例如 `composer.tsx:170` 的 `≤` 究竟渲染成什么样、面板上是否可截图核到） | 要跑面板／截图，属仪器外溢 | 面板侧验收程 | WebView2／截图授权（本程未开任何 UI） |
| 5 | **`U+2B00–U+2BFF` 与 `U+2600–U+27BF` 在全仓（含 `docs/`、`.scratch/`）的存量** | 与"补宽哪一段"无关（那两段本来就在仪器里） | 任何人 | 一条 `git grep -oP`，随时可查 |
| 6 | **`1F000–1F2FF`（仪器多扫段）在 `ban #8` 四棵树内的存量**＝我查了＝0；但**全仓**只查到 2 枚（`docs/reports/pending-and-issues.md:649,652`，`U+1F195`）——**那两枚是不是 owner 要留的台账风格，本程没判** | 判它要动 `docs/reports/**`（禁改） | owner | 一句话裁决 |
| 7 | `tools/d22scan/scan_test.go` 里**除 `:274` 之外**是否还有别的缺口段字形 | 该文件在禁改清单且不在射程，量它对选支无增量（本程已量：全文件该段字形 **1 枚／1 行**） | 任何人 | 一条 grep |
| 8 | **`AGENTS.md` 与 `PLAN.md` 的行号漂移**（§6 那 3 处）是否意味着文档已过期 | 那是 `AGENTS.md` 自己的账，本程禁改 | 文档程 | `AGENTS.md` 解冻 |
| 9 | `design/doubao/`、`design/old/` 这两枚未跟踪目录的**归属与去留**（它们会不会进 CI 视野） | 是 owner 自己的未提交件，本程不碰 | owner | 不需要授权，需要他一句话 |

---

## 8. 我不替 owner 选支

本程交的是**数**，不是判断。**我不写"建议 (a)／(b)／(c)"，也不写"应该补到哪一段"。**

三件**必须 owner 认、本程故意没有替他认**的岔口，逐条钉在明处：

1. **`Q-46` 走哪一支**——三支的代价数全在 §2.3（爆炸半径 18／8／26／41／47／51 枚）与 §3.2（解冻 1／3／7／11 张）。选支归 owner，本程一个字都不加。
2. **"注释豁免"算不算注释**——`internal/memory/schema.go:29,118` 那 2 行在 **Go 语法里是字符串**、在 **SQL 语义里是注释**（§0.2 末、§5.2 注）。任何"注释从严／从宽"的表述都盖不住这 2 行，需要一句**专门的**裁决。
3. **(i)／(ii) 的判据取哪条**——本程按 **sink 位**取 **(i) 31／(ii) 7**；若改按"任何可能到达终端"取 **(i) 35／(ii) 3**，改判路径与证据行都在 §1.4。这是**口径**问题，不是事实问题，所以也不由本程定。

**另外两件本程看见但没往里走的事**（只登记，不评论）：
`internal/tools/fs_write_test.go:521` 那枚与生产文案不同步的解析模板（§1.5），
以及 `frontend/fixtures/composer-states.html` 与 `composer.tsx:170` 的同形外露——
**都不属 AC#2 的射程**，要开票请另开。

---

## 9. commit 流水（本程落到本文的枚数；号一律现取，不写死）

| 序 | commit | 内容 |
|---:|---|---|
| 1 | `ade897c` 20:11 | §0 复核：四数全等 ＋ 推翻"12 枚冻结"与取数 pathspec 漏 `frontend/` ＋ 换掉注释判据 |
| 2 | `5dfb8ba` 20:12 | §1 小结①：非注释 38 行两档 (i)31／(ii)7 ＋ 外露证据链 ＋ 改判路径 (i)35／(ii)3 |
| 3 | `e0c4ada` 20:15 | §2 小结②：码点段分布（**第五段 `U+2200–U+22FF` 是简报漏项**）＋ 候选射程爆炸半径 6 行 |
| 4 | `64664d6` 20:18 | §3 小结③：冻结件 11 枚逐枚 ＋ 每张具名解冻的范围 ＋ 段×解冻交叉 |
| 5 | `d036929` 20:19 | §4 顺手两核：`AGENTS.md:38`（成立＋"表面"四处不一致）／`main.go:111` 多扫 **`1F000–1F2FF`**  whole 段（成立且更宽） |
| 6 | `a2d84ac` 20:19 | §5.1 表 A.1 逐枚 78 行（注释） |
| 7 | `94f4e33` 20:21 | §5.2 表 A.2 逐枚 38 行（非注释）＋ §5.3 表 B 逐枚 25 行（`frontend/`） |
| 8 | 本节这枚 | §6 共树观察与通知登记 ＋ §7 没查的档 ＋ §8 不替 owner 选支 ＋ §9 流水 |

**末节说明**：按本文开头定的约定，第 8 枚 commit 自己的回执没人带，所以列在 §10。

---

## 10. 终局回执（原样贴）

**先记一件实况**（它本身就是要留给下一位的东西）：本节第一次生成时，我把
`git log --oneline -1` 与 `git show --name-only HEAD` 直接贴了进来，
而**在我贴的那一刻 HEAD 已经不是本程的了**——共树里 `acceptor-ticket136-ac14-r1` 刚落了
`115173b`（`docs/evidence/s1/136-ac14b-impl.md`）与 `92dd40f`。
⇒ 那一份"回执"记的是**别人那枚 commit**，已作废；下面这一份是**本程自己那枚**，
取法写死成 `git log -1 --oneline 68ff486` ＋ `git show --name-only 68ff486`（**按号取，不按 HEAD 取**）。
**这正是简报里"号是读数、会漂，别抄我的"那一句的活样本**，连贴回执都会漂。

```
$ git log -1 --oneline 68ff486
68ff486 docs(141 AC#2 盘点 r1 第8节): §6 共树观察与两栏通知计数(自称授权的输出 0 条/真通知 3 条/ 判为注入 0 条)+walkEmoji 走文件系统⇒未跟踪 design/ 草稿 69 处会进 CI 视野(静态推断,未跑); §7 本程没查的档 9 条逐给"为什么/谁能查/要什么授权"; §8 我不替 owner 选支(三处必须他认的岔口); §9 commit 流水。

$ git show --name-only 68ff486 | tail -2

docs/evidence/s1/141-ac2-stock-inventory-r1.md
```

（第 8 枚之后本程还会再落 1－2 枚**只改本文这 10 行的排版与回执**的 commit；
它们的号一律用 `git log --format='%h' 99263cc..HEAD -- docs/evidence/s1/141-ac2-stock-inventory-r1.md` 现取，**不写死**。）

---

## 11. 本程的写集自证

- 本程**只写过这一枚文件**：`docs/evidence/s1/141-ac2-stock-inventory-r1.md`。
- 每一枚 commit 都带显式 pathspec `-- docs/evidence/s1/141-ac2-stock-inventory-r1.md`，
  且每次 commit 前跑过 `git diff --cached --name-only`（结果只有那一条路径）。
- 未 push、未 `--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`。
- 临时件全部在 `D:\tmp\wisp141inv\`（**只建不删**，仓内零临时件）：
  `classify.pl`（Go 词法器）· `sink.pl`（sink 回溯器）· `mktable.pl` · `fe2tsv.pl` ·
  `classified.tsv`（116 行原始分类）· `fe-classified.tsv`（25 行）· `fe-posmap.tsv` ·
  `tree/`（41 枚 `git cat-file blob 99263cc:…` 副本）· `fe/`（10 枚）· `tbl-a1.md`·`tbl-a2.md`·`tbl-b.md`·
  `raw-hits.txt`·`fe-raw.txt`·`mycommits.txt`。
- 禁改面**一个字节都没动**：`tools/d22scan/**`、`docs/PLAN.md`、`docs/specs/**`、`AGENTS.md`、
  `internal/**`、`cmd/**`、`frontend/**`、`design/**`、`.scratch/wisp/issues/**`、`docs/reports/**`。
