# 212-a2 — 注释引用仓内路径的全仓普查（票 212 AC#1 的分母/分子名册，a1 的补全腿）

> 本件是**只读普查腿 `212-a2`** 的交件，对应票 `.scratch/wisp/issues/212-comments-cite-evidence-files-that-do-not-exist.md` 的 **AC#1**。
> 前程 `212-a1` 的骨架（`.scratch/wisp/probes/212/a1/census.md`）只有 §0 锚与脏名册是有效产出；本件 §0 独立自取锚，§1–§8 全部为本腿现量。**零产码、零 Go 命令**，全部读数只来自 `grep`/`sed`/`test -f`/`git ls-files` 这一类盘上尺，每格带推导式原文。
> ⛔ 本腿不读不引 `frontend/**` 与 `design/**`（两层禁令：不打开、结论不引其内容）。

---

## §0 起手锚 + 起手脏名册（产码子集）

### 0.1 锚点（本腿自取）

| 项 | 值 | 尺 |
|---|---|---|
| 起手时刻 | `2026-10-03T16:15:50+08:00` | `date -Iseconds` |
| HEAD | `1ba22d25914600fb0156a67b414b13aee3fd9c9a` | `git log -1 --format=%H` |
| 分支 | `dev` | `git rev-parse --abbrev-ref HEAD` |
| 起手全树脏条目总数 | 482 条（`git status --porcelain \| grep -c "^"`，a1 同窗读数） | |
| 跟踪中的 `.go` | **847** 枚（`git ls-files '*.go' \| wc -l`；含 `design/` 1 枚、`frontend/` 0 枚） | |
| 跟踪中的 `.md` | **891** 枚（`git ls-files '*.md' \| wc -l`；含 `design/` 2 枚、`frontend/` 0 枚） | |
| 跟踪中的 `.sh/.ps1/.py` | 135+6+34（见 §4 分母表） | |

> 本件全部读数的唯一锚点＝上面这枚 HEAD。HEAD 今天在动；任何不带锚的读数一律无效。

### 0.2 起手脏名册（产码子集：`cmd internal tools scripts`）

尺：`git status --porcelain -- cmd internal tools scripts`（2026-10-03T16:15:50+08:00，HEAD `1ba22d25…`）

```
 M internal/llm/enabled_reach_261_test.go
 M internal/llm/resolver.go
?? internal/llm/enabled_gate_261_r1_test.go
```

**对本腿的意义**：`internal/llm` 有 `261-r1` 写腿在飞的未提交修改。本腿读到的这三个文件是**未提交中间态**；`internal/llm/enabled_reach_261_test.go` 与 `resolver.go` 的注释引用按**中间态**计入（读数时注明），`enabled_gate_261_r1_test.go` 是未跟踪件、不在 `git ls-files` 分母里，其注释引用**量不到**（§7 登记）。本腿对脏面上的内容只登记、不评判、不引用其具体注释文本。

（其余脏面：`cmd/wisp/panel_host_windows.go`/`panel_resident_windows.go` 在 a1 同窗的脏名册里、本窗已不在——写腿 `255-r3` 已收；`.scratch/**`、`design/**`、仓根垃圾件 30+ 枚零字节 heredoc 产物照 a1 §0.2 登记，本腿不处置。）

---

## §1 尺的校准（4 枚已知样本逐枚贴）

**校准目的**：本机 Windows＋Git Bash，路径形状有 `D:\…`/`D:/…` 两写法、8.3 短名、反引号包裹形。尺没在已知样本上校准就全仓跑＝白跑。

**抽取尺（候选）**：从跟踪产码的**注释行**里抽"看起来像仓内路径"的 token，正则
`[A-Za-z0-9_][A-Za-z0-9_./+-]*\.(md|go|sh|ps1|py|json|sse|txt|log|toml|sql|yml|yaml|csv|tsv|html|css|js|ts|png|jpg|svg|exe|db|mod|sum|golden)`
配合行首 `//`（Go）或 `#`（sh/py）判定注释行。**为什么这么抽**：票面 AC#1 点名的是"注释里引用的仓内路径"，路径必有扩展名才谈得上"那枚文件盘上存在与否"；纯目录名（无扩展名）"存在与否"歧义太大，且票面例子全是文件引用。

| 样本 | 注释原文（所在行） | 引用 token | 盘上复量 | 期望类 | 尺应判 |
|---|---|---|---|---|---|
| A：`internal/tools/subagent_197_test.go` 顶部注释块 | 「Reading is in docs/evidence/s1/197-subagent-entity-r1b.md」（注释行，前缀 `//`） | `docs/evidence/s1/197-subagent-entity-r1b.md` | `ls` rc=0，24790 字节 | ① | ① 存在 |
| B：`cmd/wisp/subagent_stream_key_197_test.go:16-18` | 「Reading is in docs/evidence/s1/197-subagent-entity-r1b.md.」（句号紧贴扩展名后） | `docs/evidence/s1/197-gent…r1b.md`（token 抽取时需剥尾标点） | rc=0，同上 | ① | ① 存在 |
| C：`cmd/wisp/slo_report_144_windows_test.go:851` | 「recorded the deletion at docs/evidence/s1/152-...-accept-r1.md §2.2」 | `docs/evidence/s1/152-...-accept-r1.md` | 字面 rc=2（不存在）；全名 `docs/evidence/s1/152-subject-death-never-measured-r1-accept-r1.md` rc=0 | ③ | ③ 不可机读（**不是②**） |
| D：`cmd/wisp/config_reload.go:334` | 「TestTicket255RestartTierKeysAreBackedByATest in restart_tier_keys_255r2_test.go」 | `restart_tier_keys_255r2_test.go`（**裸文件名，无目录**） | `ls cmd/wisp/restart_tier_keys_255r2_test.go` rc=0 | ① | ① 存在（**前提：全仓唯一**；若多枚同名⇒归③） |

**校准结论（尺的三条修正，全部由样本 C/D 教出）**：
1. token 尾部标点（`.`、`、`、`，`、`)`、`」`）要剥掉——样本 B 的句号紧贴扩展名，不剥会把存在判成不存在；
2. 含 `...`/`…` 的路径**不进②的名册**、进③（样本 C：省略号是手写缩写，不是"路径不存在"）；同理，**裸文件名**（无 `/`）能被全仓唯一定位的判①，不能唯一定位的判③（样本 D）；
3. 路径含 `:` 行号（如 `foo.go:361`）要剥 `:行号` 再 `test -f`；`…`/`...`/`xx-…-yy` 缩写形、`*` 通配形、反引号内的形一并进③候选。

**排除什么、为什么**：
- 非"注释引用"的形状一律不算分子：出现在**字符串字面量**里的路径（不是注释，不在 AC#1 射程）、出现在**代码**里的路径（是逻辑不是路标）、出现在**行内注释**里的路径（算！——AC#1 说"注释"，行内注释也是注释，保留）。
- URL 形（`http(s)://…`）不算仓内路径，排除。
- 指向**本注释文件自身**的路径存在与否依实判（保留）。

---

## §2 四类名册（每类：计数 + 推导式 + ②③④ 逐枚表）

> 总口径（分母定义）见 §4；本节计数全部基于 §4 表 A 的分母。

### 2.1 ① 引用的路径存在

- **计数**：〔见 §2.4 汇总〕
- **推导式**：对 §1 抽取尺抽出的每个引用 token（剥尾标点、剥 `:行号`、裸文件名做全仓唯一定位）`test -f` 判存在；命中即①。
- 逐枚表：存在类枚数大（数百），不逐枚列名（票面只要求②③逐枚点名）；抽样复核 4 枚见 §1。

### 2.2 ② 引用的路径盘上不存在

- **计数**：〔见 §2.4 汇总〕
- **推导式**：同上抽取，`test -f` rc≠0 且**路径形态完整**（无省略号、无通配、有目录、全仓能唯一定位）⇒②。
- 逐枚表：〔见 §2.4〕

### 2.3 ③ 写法不可机读

- **计数**：〔见 §2.4 汇总〕
- **推导式**：引用 token 含 `...`/`…`/`*`/`<…>`/「那一族」类通指词，或**裸文件名**不能全仓唯一定位，或**只给票号不给路径**的"读数在某票"形 ⇒③。
- 逐枚表：〔见 §2.4〕

### 2.4 汇总

〔占位：本节骨架先落，逐类计数在量毕后回填，无"(待填)"字样留给回填轮〕

---

## §3 与票面现量（②=0、③=1）的对账

〔本节骨架先落；对账结论在 §2 量毕后给出，若推翻票面逐枚点名。〕

---

## §4 分母口径的取舍说明（probes/commit-msg 进不进、为什么）

**票面 AC#1 的字面**：「全仓**注释**里引用的**仓内路径**」。注释只存在于代码文件里。据此本腿分母两档：

**表 A（主分母）**：跟踪中的 `*.go`（847 枚，剥 `design/` 1 枚＝不读不引）＋跟踪中的 `*.sh/*.ps1/*.py` 里**产码侧**（`scripts/**` 14 枚、`tools/d22scan/runtests.sh` 1 枚）＋跟踪中的 `*.md`（891 枚，剥 `design/` 2 枚）。
- `.go` 用 `//` 判注释；`.sh/.ps1/.py` 用 `#` 判注释；`.md` **没有"注释"语法的边界**，按"正文引用"计入并单列（票面例子里 `docs/evidence/s1/**` 的引用大多源头在 Go 注释，md 里的是复述层；md 引用是否算"注释"是口径裁点，本腿**算入并单列**，理由：md 里同样有"读数在某枚文件里"的路标句子，撒谎危害同形）。
- ⚠ `scripts/**` 的 `.ps1` 有一枚 `scripts/spike/run.ps1`，spike 件，照算（它在跟踪面里）。

**表 B（不进分母，逐条给理由）**：
1. **`.scratch/wisp/probes/**` 的脚本与探针件**（135 枚 `.sh`/34 枚 `.py` 的绝大多数）：一次性实验件，其中大量是**变异体/门禁复制**，其"引用"自己票面的探针件，污染分子。票面点名的是"注释"、是产码下游的导航路标；probes 件不是下游路标。**不进主分母；补一张副表口径见下**。
2. **`.scratch/commit-msg-*.txt` / `.scratch/wisp/probes/**/msg-*.txt`**（约 1900 枚 `.txt`）：commit message 草稿，**不是注释**（注释只存在于代码文件里）；票面例子全部是 Go 注释。**不进任何分母**。
3. **`frontend/**`、`design/**`**：派单硬禁令两层（不打开、不引内容）。`git ls-files` 显示 `.go` 含 design 1 枚、`.md` 含 design 2 枚，这 3 枚**从分母剥掉**且不读。
4. **`.sse/.json/.log/.png`** 等非代码件：无注释语法，不进。
5. **未跟踪件**（如 `internal/llm/enabled_gate_261_r1_test.go`）：不在 `git ls-files` 名册里；其注释引用**量不到**，§7 登记。

---

## §5 装牙的代价一句话预告（给 AC#2 裁，本腿不裁）

〔本节骨架先落；量毕后按四类计数给一句话代价预告。〕

---

## §6 推翻清单（票面、212-a1 骨架、本派单的待验断言逐枚验）

〔本节骨架先落；每一条"票面说 X／实量 Y／推翻与否"在量毕后回填。〕

## §7 量不到的格子 + 复量法

- `internal/llm/enabled_gate_261_r1_test.go`（未跟踪，`261-r1` 在飞）：其注释引用不在分母。复量法：待该腿 commit 后重跑本件 §2 的尺。
- `frontend/**`、`design/**`：派单硬禁令，不可量。复量法：解除禁令后由下一程补。
- 字符串字面量里的路径引用：不在 AC#1 射程（"注释"），未量。复量法：把 §1 抽取尺去掉 `//`/`#` 前缀限定重跑。
- 跨仓库/工作树外绝对路径的**判定**依赖本机挂载盘形（`D:\`、`\\`、8.3 短名），④类尺在 §2.4 汇总时给。

## §8 自查

〔本节骨架先落；交件前自查表。〕
