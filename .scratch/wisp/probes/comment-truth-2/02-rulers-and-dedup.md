# comment-truth-2 / 02 — 尺的原始输出与查重证据（HEAD `7a452d0a`）

## A. 撞钉尺（逐字串 × `*_test.go`）

尺＝`git grep -nF "<串>" HEAD -- '*_test.go'`，13 条串逐一跑，**全部空输出、rc=1**（⛔ 未接 `2>/dev/null`，逐条检 rc）：

`NO production caller yet`／`zero production callers in this tree`／`not yet wired by tickets 18/19`／`This call is the caller`／
`Nothing calls it yet`／`zero assignments outside`／`had zero callers`／`no production path ever did`／`has no production caller`／
`ZERO product callers`／`TierOf had ZERO`／`Wiring still owed by ticket 12`／`Nothing else in this repository calls Handle`。

⇒ 结论：**这 15 枚没有一枚被测试逐字钉住**，改文案不会因"字面量"红。
真正的形状约束来自**扫同一文件**的仪器（非逐字钉），四枚具名见 `01` 汇总，其原文：

- `internal/panel/composer_dispatch_test.go:375 TestDispatcherSpellsNoRouteLiteralOfItsOwn` → `:382 routeLiteral := regexp.MustCompile(` + '`"panel\.[^"]*"`' + `)`，且 `:384-386` **显式跳过 `//` 开头的行**。
- `internal/panel/composer_dispatch_test.go:708 TestTheInboundHopAddsNoSwitchingCapability` → `:714 for _, banned := range []string{"checkoutBranch", "changeRepo", "repoPicker", "branchSelect", "vcs.switch"}`，判据是 `strings.Contains(string(data), banned)`——**整份文件、含注释**（`:376`/`:710` 都是 `os.ReadFile(composer_dispatch.go)`）。
- `cmd/wisp/firstrun_257_test.go:430-438` → `strings.Count(string(source), "\nfunc ") != 1` ＋ 签名逐字 `func ensureFirstRunConfig(dataDir string, stderr io.Writer) (bool, error) {`。
- `cmd/wisp/resident_hotkey_258_test.go:73-79` → `if !strings.Contains(string(src), "config.LoadFile")`。
- `cmd/wisp/resident_grant_writer_265_windows_test.go:465-470` → `parser.ParseFile(... resident_task_source_windows.go)`（AST，注释不进 AST）。
- `internal/panel/composer_dispatch_test.go:437 TestPanelHostIsAttachedAndNamesTheWindowHops` → 谓词是"树里**必须有**宿主/WebView2 通道"，读的是标识符与 go.mod，⛔ 不读注释；台账 `A720` §4 已由编排者真跑（rc=0、PASS 0.24s），本程不再跑。

## B. 归口尺

### B.1 `issues/README.md` 里的"都归票 X"句

尺＝`git grep -nE '都归.{0,6}票 ?[0-9]{3}' HEAD -- .scratch/wisp/issues/README.md` → **1 命中**：

`README.md:70` 去 markdown 加粗后的**真串**＝`⚠ 两处已知的不彻底，都归票 262，别把这行当"已经装好牙"：① 这条帽只有词面——今天没有任何仪器`

⇒ 具名：该句射程＝**票名长度门禁**（`scripts/check-path-length-budget.sh` 的 `HAT_NAME_LIMIT=100`），与本族"状态级注释"**不是一件事，⛔ 不能当去处**。
再跑宽形 `归 ?\*?\*?票 ?[0-9]{3}|一律归|统一归|都并进` → 仍只有那 1 命中。⇒ **README 里不存在"注释更正都归票 X"这类归口句**。

### B.2 台账查重（`docs/reports/pending-and-issues.md`，尺＝`git grep -nF`）

| 串 | 命中 | 是否已登记 |
|---|---|---|
| `This call is the caller` | `:14109`（A720 §3 ②） | **已登记**＝P15 |
| `not yet wired by tickets` | `:14112`（A720 §3 ⑤） | **已登记**＝P13 |
| `NO production caller` | `:14093` 等多处，其中 A720 §3 ④ 具名 `internal/panel/composer_dispatch.go:50` | **已登记**＝P01 |
| `TierOf had ZERO`／`no production path ever did`／`zero assignments outside`／`Wiring still owed` | **空输出、rc=1** | 未登记＝本程新读数（P02/P06/P07/P08/P09/P10） |

另：P12 的过期**已在票面登记**（`201-…md:73` 逐字「⛔ 票面 §10 那把尺的读数（"三枚入口生产零调用者"）已过期」），P14 与 P15 同体（A720 §3 ①②）。
⇒ **⛔ 这五枚（P01 P12 P13 P14 P15）不算本程新发现**，本程只补它们的"落到哪格"。
具名一枚**不在 15 枚名册内**的：`A720` §3 ③ 那把 `SealDir` 尺（`:14110` 一带）判的是**票 132 票面断言**（表二/表三族），
它在 HEAD 复跑仍是"定义行 1 枚、调用点 0"⇒ 表一·附那节原样成立，⛔ 与本程的 15 枚无关，别混进"已登记"计数。
⇒ 换算：`A720` §3 的**五把尺**落在本名册上＝**4 枚**（P01 P13 P14 P15），加票 201 票面自登记的 P12＝**5 枚在册**。

### B.3 票面在册与承接格现量

尺＝`git ls-tree -r --name-only HEAD .scratch/wisp/issues` ＋ 未勾框 `git show HEAD:<票> | grep -cE '^[[:space:]]*- \[ \]'`（⛔ 不用 `^- [ ]`）。

| 票 | 文件名带 `-done`？ | 行数 | 未勾框 |
|---|---|---|---|
| 12 | 否 | 187 | 3 |
| 33 | 否（`33-panel-host-c27.md`） | 394 | 13 |
| 35 | 否 | 425 | 6 |
| 114 | 否 | 114 | 9 |
| 132 | 否 | 89 | 5 |
| 186 | 否 | 73 | 9 |
| 187 | 否 | 39 | 9 |
| 198 | **是 `-done`** | — | **0** |
| 201 | 否 | 73 | 5 |
| 223 | 否 | 70 | 2 |
| 246 | 否 | 97 | 1（只剩 AC#8；**AC#7 是 `[x]`**） |
| 255 | 否 | 167 | 2 |
| 259 | 否 | 59 | 5 |
| 262 | 否 | 130 | 2 |
| 18 / 19 / 90 | **均 `-done`** | — | — |

未勾框里含"注释/comment/更正/stale"字样的**只有 2 处**（尺＝框行再筛词）：`33:AC#13`（冷启动往返探测，非注释格）、`114:AC#10`（`panel.ts` 注释自陈两种装法，⛔ 前端文件，本程禁读不判）。
⇒ **除 P01/P03（票 33 `AC#9`）之外，没有任何一张票的未勾框是为"更正过期注释"留的格。**

## C. 被引对象存在性尺（〔引用假凭据〕检查）

尺＝`git ls-tree -r --name-only HEAD \| grep <名>` ＋ `git grep -n "func <测试名>" HEAD -- '*_test.go'` ＋ `git cat-file -t <短 sha>`。

- `cmd/wisp/panel_inbound_33_test.go` **存在**（434 行）；`:183 inboundHandleCallSites33(t)` **存在**，其定义 `:278` 起，射程是 **`cmd/wisp` 目录**的非测试文件。
- `docs/evidence/s1/33-minimal-inbound-hop-r1.md`、`docs/evidence/s1/33-inbound-hop-design-a1.md`、`docs/evidence/s1/219-approval-reply-surface-c1.md`、`docs/evidence/s1/181-186-git-detection-census-c1.md` **四份都存在**。
- `.scratch/wisp/probes/223/c1/census.md` **存在**，且 `:12/:13/:23/:24` 逐字带着 P07/P08/P09 那三条 grep 尺（⛔ 不是假凭据）。
- 测试名 `TestGitDimensionHasNoModelCallableTool`（`internal/panel/git.go:19` 点名）**存在**＝`internal/panel/git_test.go`。
- `cmd/wisp/resident_task_source_246_windows_test.go:389-391` **存在且逐字对得上**（`:389 if _, err := os.Stat(filepath.Join(dataDir, "config.toml")); !errors.Is(err, os.ErrNotExist) {`／`:390 t.Errorf("AC#7 RED: the leg created a config.toml it was never asked for (err %v)", err)`）。
- P10 块内的其它锚：`internal/config/defaults.go:58` ＝ `func NewDefaults() *Config {` **对**；`internal/config/writeguard.go:125-131` 区间内 `:127` ＝ `// legitimate base and writing it creates what first-run did not.` **对**；`loader.go:238` **错**（现量 `:238` ＝ `if ref := c.Voice.Realtime.APIKeyRef; ref != "" {`，`func SaveFile` 在 `:252`）。
- 枚数引用：P06/P07 的块头 `:7-8` 逐字 `Three of config.Manager's hooks had no production / assignment point before this file (measured at HEAD 7ffa9520, census .scratch/wisp/probes/223/c1/census.md)`——枚数 **3** 与子弹数 **3** 对得上；锚 commit `7ffa9520`／`8a3790f0`／`dd92bb92` 三枚 `git cat-file -t` 均回 **commit**。

## D. 调用形状尺（HEAD 复跑，射程逐枚写明）

统一射程＝`HEAD -- '*.go' ':(exclude).scratch' ':(exclude)*_test.go'`。

- `[A-Za-z0-9_]\.AskOnTaskRoot\(` → **空、rc=1**；`AskOnTaskRoot` 全形状只出现在 `resident_approval_windows.go:688`（注释）／`:697`（定义）、`resident_task_source_windows.go:14`（注释）、`resident_windows.go:248`（注释）。
- `askConfirmation` → `:641` 定义、`:700` 包内调用（在 `AskOnTaskRoot` 体内）、外加两条注释行（`resident_task_source_windows.go:14`、`resident_windows.go:248`）。
- `ParseComposerRequest\(` → `internal/panel/bridge.go:126` 定义、`internal/panel/composer_dispatch.go:155` 调用。
- `modeWrites` → `run.go:331/:337`（注释＋字段定义）、`run.go:674`（赋值）、`panel_inbound.go:248/:272`（**同名局部变量**）⇒ 字段**读侧 0**。
- `.Handle(` 在 `cmd/wisp` 的 `*panel.ComposerDispatch` 命中＝`panel_host_windows.go:821`、`panel_inbound.go:163`（其余 30 余行是 `windows.Handle`／`observe`／`logsink` 的**同名不同型**，按类型剥掉＝第 107/119 条）。
- `cmd/wisp/main.go:115 case "panel-inbound":` 现存；`:47` 的帮助文本逐字写着 `wisp panel-inbound  inbound panel dispatch leg (ticket 33 AC#9)`。

## E. 与本仓既有"尺"的重复检查（查重结论）

本程用的三把尺——① 逐字串 vs `*_test.go`、② `README` 归口句、③ 票面未勾框承接格——**都不是新判据**：①是 `A720` §1 采的"按这句话在说什么扫"的下游用法、②③是派单点名的既有检索式。
⛔ 本程**不新建判据、不翻框、不改票**；具名一条防重结论：`A720` §3 ⑤ 明写"这族…落哪一票待我先查重再裁（第 123 条）"，本文件 `01` 的"归口"列就是那一次查重的**读数**，裁定权在编排者。
