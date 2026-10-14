# 279/r1 — P39 注释落地＋票 255 名册七处 run.go 读数刷新（随行账）

- 腿：`279-r1`（落地腿，实现者）。起手 HEAD＝`ade60808`；取数/落盘期间编排者一腿落 `1d19b60b`（12:24，面＝orch 探针件＋台账，**未碰 cmd/wisp**）；裁定diff 以 `1d19b60b` 为基。
- 票：`.scratch/wisp/issues/279-p39-comment-and-the-255-roster-citation-ride-together.md`（本腿落 AC#1；AC#2/AC#3 框不翻，判语归非实现者）。
- 触发（AC#2）：`run.go` 最近一次真改＝`bfc74ae9`（315-r1，2026-10-14 11:27:28 +0800，其 run.go numstat＝`34 5`）；`git diff --numstat bfc74ae9..HEAD -- cmd/wisp/run.go` 空 ⇒ 随行锚唯一（`logs/18-trigger.txt`，746 字节）。
- 料文：`.scratch/wisp/probes/comment-fix-prep-1/01-ready-to-apply.md` P39 块（段首 :327；原文 :332-337；替换 :341-349）。

## 1. 锚复验（贴前；内容锚，票面行号只当找法）

尺＝python 现读比对（脚本留盘 `C:/tmp/wisp279r1gotmp/check_p39_v2.py`，可对 git 逐字重跑）：

| 检查 | 读数 |
|---|---|
| A 料文原文 :332-337 == `HEAD:cmd/wisp/run.go` 行 331-336（贴前） | **True** |
| B 料文替换 :341-349 == 工作树 `run.go:331-339`（贴后） | **True** |
| C 块行数 | 原 6 行 → 替 9 行 |

- 判别位现读（贴前）`:333` 逐字 `// Nothing calls it yet - the WebView2 "event -> ParseComposerRequest" hop does` 与料文一致 ⇒ 该处未被别人改过。
- 证明留 `logs/19-anchor-verbatim.txt`（A/B 两行 True）。

## 2. 贴入（逐字，未擅改）

- `cmd/wisp/run.go` 现块＝**331-339**（9 行，前导 1 tab 保留）；numstat **+6/−3**（净 +3；首两行与末行被 diff 当共有上下文）。
- 全注释：该件 diff 正文 9 行，R1 尺 `^[-+][[:space:]]*//` 命中 **9/9**，非注释 **0**。
- 贴后 `modeWrites` 命中点＝注释 :331／定义 :340／赋值 :677，**读侧仍 0** ⇒ 替换文本「nothing attaches this one」与现盘一致（`logs/18-trigger.txt`）。
- ⚠ **具名顶回（料文内嵌读数已漂）**：替换块第 4 行内嵌 `composer_dispatch.go:155`。料基 `9e8477ed` 时 `ParseComposerRequest(` 确在 `internal/panel/composer_dispatch.go:155`；现读 **:161**（漂 +6；该文件自料基后被 `d7a3d62c`、`341ca42e`(186-r1) 等改动）。派单要求「逐字照料文／别自己改写」⇒ **原样贴入未擅改**；是否随行刷新到 :161 待编排者裁（本腿未做）。

## 3. 名册读数刷新（`cmd/wisp/config_readers_255.go`，7 处）

新读数＝贴后现读重跑（`sed -n '<N>p' cmd/wisp/run.go` 逐行验过）：

| 处（新行号） | 引用对象/token | 旧记录 | 贴前实读 | 新记录 |
|---|---|---|---|---|
| :56 | `run.go:423-424` takes `cfg := mgr.Config()` | 423-424 | 423-424 ✓ | **426-427** |
| :58 | execute() reads that copy `(run.go:991` | 991 | 1020（已漂 +29） | **1023** |
| :107 | `[res := llm.NewResolver(cfg, st)]` | 435 | 435 ✓ | **438** |
| :113 | `[rt.cfg = cfg]` | 424 | 424 ✓ | **427** |
| :114 | `[cfg := rt.cfg]` | 991 | 1020 | **1023** |
| :115 | `[PerToolTimeout: time.Duration(cfg.Agent.PerToolTimeoutMS)]` | 1014 | 1043 | **1046** |
| :117 | **名册证据引用**（`"agent"` 行）`[cfg := rt.cfg]` | 991 | 1020 | **1023** |

- 票面/派单列 5 处（1 证据引用＋4 散文）；本腿整件重扫 `grep -n "run\.go:"` 另抓到 2 处同形散文（:56/:58，票面未列）⇒ 一并对齐；不如此则贴后静默过期（派单「行号一律现读重跑」「别让它们静默过期」）。
- token 真身核验（贴后，`grep -c`）：`cfg := rt.cfg`＝1（:1023）／`rt.cfg = cfg`＝1（:427）／`res := llm.NewResolver(cfg, st)`＝1（:438）／`PerToolTimeout: time.Duration(cfg.Agent.PerToolTimeoutMS)`＝1（:1046）／`cfg := mgr.Config()`＝1（:426）——名册尺「token 唯一命中当判据」不受影响。
- ⛔ 未放宽任何断言；`cmd/wisp/config_receipt_255_test.go` 一字未动。

## 4.「改动全为注释」自证（279 AC#3 那把尺）

R0 基数＝`git diff HEAD -- <两枚写面件>` 正文行（`^[+-]` 排 `+++`/`---`）：

| 量法 | 读数 |
|---|---|
| 正文行合计 | **23** |
| R1 整行注释 `^[-+][[:space:]]*//` | **21/23** |
| R2 形状探针 `/*|*/|*` | 0 |
| R3 形状探针（产码行带 `//`） | 0 |
| **非注释** | **2**＝名册 `"agent"` 引用字符串的旧/新那一对（`:991`→`:1023`）：AC#1 明文要求改的那枚**字符串读数**，字符串不是注释、按构造必现 2 行 |
| `run.go` 单独非注释 | **0** |

明细：`logs/16b-diff-ruler-vs-HEAD.txt`＋`logs/16b-diffbody-vs-HEAD.txt`（早先同尺对同内容的读数亦存 `logs/16-diff-ruler.txt`）。

## 5. 门禁读数（全绿）

| 门 | 命令 | rc | 读数 | 日志（字节） |
|---|---|---|---|---|
| 构建 | `GOFLAGS= go build ./...` | 0 | 空输出 | `10-build.txt`（0） |
| 格式 | `gofumpt v0.12.0 -l` 两枚／`gofmt -l` 两枚 | 0 | 均空 | `11-gofumpt.txt`（0）／`12-gofmt.txt`（0） |
| 255 名册 | `go test ./cmd/wisp/ -run 'TestTicket255' -count=1 -v` | 0 | 顶层 20 PASS／0 FAIL，无 SKIP；`RosterEvidenceLinesStillSayWhatTheyClaim` PASS 且无 drifted Logf；ok 18.931s | `13-test255.txt`（24378） |
| D22 扫描 | `sh scripts/d22scan.sh` | 0 | clean；cmd/ examined 125 | `14-d22scan.txt`（22968） |
| 长度门 | `sh scripts/check-path-length-budget.sh` | 0 | VERDICT GREEN；over=59／roster=59（59=59）；**暂存后复跑同绿** | `15-length-gate.txt`（1457）／`15b-length-gate-staged.txt`（1457） |

## 6. 面与纪律

- 写面＝`cmd/wisp/run.go`（仅注释块 331-339）／`cmd/wisp/config_readers_255.go`（读数/引用）／本 probes 目录／票 279 尾 1 条 bullet。
- ⛔ 未碰：315-r1 的 `OpenTask` 语义行（`admitTask` 一带）、`.github/**`、`scripts/**`、三枚冻结件、`thresholds.go`、golden、`allowlist.txt`、`PLAN.md`、`specs/**`、`frontend/**`、`design/**`。
- 只 commit、⛔ 未 push；`git add` 均显式 pathspec（⛔ 无 `add -A`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`／`--no-verify`）。
- 具名没做：① `composer_dispatch.go:155`→`:161` 未擅改（见 §2，待裁）；② `docs/reports/pending-and-issues.md` 台账未动（不在写面）；③ 票 279 三框未翻（判语归非实现者）。
