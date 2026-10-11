# 309-r2 · 90 收尾账（写腿，⛔ 裁决者）

起手 HEAD＝`202ed4a5`（派单给的锚，现跑一致）；交回 HEAD＝本腿最后一笔（`git rev-parse --short HEAD` 现跑，见 `logs/r359-14-handin.log`＝`d75db9ee`＋本证据笔）。
其间 `dev` 上落过 `308-r2` 的六笔（`6f389389`…`daa2f9be`），其中两笔改的是 `internal/risk/*_test.go` ⇒ 本腿所有整包名册/导出树读数与它们**同面⛔ 同树**，具名在 `logs/behaviour-vs-anchor.stdout` 的 V1/V1b。

## 1. 五件各自的一句话＋承重的尺

| 件 | 做了什么 | 证明它的尺（读数） |
|---|---|---|
| ① roots | `tools` 进、`scripts` 也进（具名答复） | `logs/roots-rulers.log`：三枚 nested module 前缀全在 `github.com/CarlosShao/wisp/` 之内；活例尺 1 命中＝`tools/signmodels/main.go:34` |
| ② 落点 | 排除支→"合计 1 枚＋具名 `internal/agent/loop.go`" | 钉自打印 `构造点=1（==1…）分母=293（cmd 40／internal 229／tools 13／scripts 11）`⇄第三把 `find` 尺 293 对得上；RIG-4b 往 agent 里加一条⇒`rc=1`、名册 2 枚、分母 294 |
| ③ 三形 | 补"声明＋赋值"与"同文件别名"；裸名按 `package` 子句放行 | `nail-rig-1/2/3.log` 各 `rc=1`（旧钉 2/3 形⛔ 绿）；六枚 control 双向在场（含"别的包也叫 ToolRequest"断 0）；`nail-rig-5-elided.log`／`nail-rig-6-crosspkg-alias.log` 各 `rc=0`＝具名漏形的读数 |
| ④ 注释 | 三处说实话＋具名 roots 判据＋漏形清单 | `logs/r359-10-final-rulers.log`＝71／69／2，剥行首 `//`＝1；`logs/r359-9-filterdemo.log`＝同枚 fixture 两把口径 2⇄1；R5＝`sed -n '211,213p'` 逐字 |
| ⑤ prov 肢 | ⛔ 造仪器，注释里另记一段（族／谁响／⛔ 具名仪器） | `logs/prov-rigs.stdout`：MUT-F（摘 `injected`）整包 `rc=0`；MUT-E（摘 `prov == nil`）整包 `rc=1`、顶层 16 枚红＝管道噪声；R359-4＝`Assessor:` 普查 2 行⇒生产写点 0 |

行为逐字节不变＝`git diff -U0 202ed4a5..d75db9ee -- internal/tools/bridge.go` 增 44／删 9 **全为 `//` 行**、非注释 0；两枚 blob 剥整行注释后 `cmp` **IDENTICAL**（sha `496ed6d6…`，854⇄854）；守卫那行 1⇄1 逐字同一；CR＝0、`i/lf w/lf attr/text eol=lf`。
门禁＝`rc_build=0`／`rc_vet=0`／`rc_d22scan=0`（两轮红⛔ 别腿：ban #9 phantom-citation，我自己写的 `cmd/internal/tools/scripts` 与一枚未落盘的件名）／`rc_pathlen=0`。
整包名册（两棵导出树，同面可比）＝顶层 210/1/0 ⇄ 210/1/0、含子测试 278/1/0 ⇄ 278/1/0、`comm` 作差 **added 0／removed 0**、`SKIP` 两向 0；两枚 236 空 id 载体两棵树照绿。

## 2. 新增欠读数 6 枚（⛔ 我盖章哪格该翻，判语归 `309-v2`／编排者）

1. **`Q-309-r2-word-ruler-still-self-referential`** — `ToolRequest{` 词面尺今天 71／69／2；注释⛔ 消灭这个咬尾，只是把它**写进尺的口径**（自带"按文件名字段滤＋剥行首 `//`"两支，并指解析器尺）。⇒ 残余＝"任何新注释都可能再挪一次分子"，⛔ 由注释负责，靠解析器尺兜。
2. **`Q-309-r2-nail-blind-shapes-unrigged`** — 钉头具名七条漏形，本腿只给两条装了实测读数（省类型嵌套／跨包别名）。其余五条（函数体内别名、改名 import 与 dot-import、`new()`／`make()`／反射／解码、定义型＋转换、别处一枚叫 `agent` 的包）＝**具名⛔ 实测**。⛔ 装成"已量过"。
3. **`Q-309-r2-pinned-total-is-constant-1`** — 新钉把"合计 1 枚"当常数钉住。将来若有**正当**的第二条生产构造点（带非空 id 的面板/球宿主），钉会红且红句给的是"接上非空 id 或进仪器入口"，⛔ 自动授权第二枚。归口＝下一次有人真要写它时由编排者裁（同 `AC#0` 形 ⓑ 的口径），⛔ 本票先放宽。
4. **`Q-309-r2-scripts-roots-admitted-by-prefix-only`** — `scripts` 进 roots 的根据只有"前缀放行"，它自己今天⛔ import `internal/`（尺＝1 命中全在 `tools/signmodels`）。⇒ 若将来有人把 `scripts/spike` 拆成 `wisp` 前缀**之外**的 module，`internal/agent` 就⛔ 可达，本钉仍会扫它（⛔ 错，只是射程变宽）；反之若把 `tools`/`scripts` 挪出前缀⛔ 更新注释，注释里那句"三枚都在前缀之内"就过期——它是**现读句**，⛔ 是保证。
5. **`Q-309-r2-tree-rigs-mixed-face`** — RIG-4b／6 的树＝`git archive 6f389389` ＋ `cp` 本腿⛔ commit 的两枚文件 ⇒ **混面**（钉是新码、其余是 `6f389389`）。commit 之后可逐字重跑（`git archive d75db9ee` 就⛔ 需要 cp），那一发本腿⛔ 补跑（时间轴在收尾），⛔ 读数算已固化。
6. **`Q-309-r2-prov-limb-colour-discrepancy`** — 派单 §5/票面 §2 那句"摘掉 `b.prov == nil` 或 `b.injected` 两发**各 rc=0**"是 `-run` 只打 309 仪器那把尺；**整包**那把尺下 `prov == nil` 摘掉＝`rc=1`。两把尺⛔ 矛盾，但"哪把尺支撑哪句判语"必须说清：支撑"`⛔ 具名该肢的仪器`"＝仪器那把；支撑"这肢有⛔ 任何反馈"＝整包那把，**结论⛔ 同**。⇒ 交 `309-v2` 定一句口径；本腿注释按"⛔ 具名＋整包会噪声地响"两半都写。

## 3. 越格自报（⛔ 瞒）

- **写面**：代码笔 `git show --name-only d75db9ee`＝`internal/tools/bridge.go` ＋ `internal/tools/toolrequest_nail_309_test.go`（2 枚，都在授予的 `internal/tools/**` 内）；其余笔只落 `.scratch/wisp/probes/309/r2/**` 与本票末尾**一条** bullet（1 加／0 删）。⛔ 翻框／⛔ 动 `Status:`／⛔ 动别家票／⛔ 动 `docs/**`。
- **多做的一件**（超出五件清单，⛔ 产品语义）：`assessorFor` 文档块里那句 `Test309... pins agent.ToolRequest{ construction points outside internal/agent to 0` 在 job 2 之后**会变成假话** ⇒ 一并改成"合计 1 枚、具名 `loop.go`、旧版整包跳过"的新口径（纯注释，落在 B1/B2 那两把尺里：44/9 全 `//`、剥注释 blob 同一）。
- **多跑的 go 命令**（都留了件）：`go test ./internal/tools/ -run 'Test309…'` 若干发；`-run 'Test309'` 一发；整包 `-count=1 -v` 三发（工作树 1 ⇄ 导出树 2）；突变台内整包 2 发（MUT-E/F，当场还原）；`go build ./...` 4 发、`go vet ./internal/tools/` 4 发、`sh scripts/d22scan.sh` 4 发（含两轮我自己的红）、`check-path-length-budget.sh` 3 发。**⛔ 跑过 `cmd/wisp` 那包的测试**（sherpa DLL 前置归编排者）。
- **多写的文件**：`logs/` 下 24 枚（脚本 `.txt`＋读数 `.log`／`.stdout`＋两枚 `.stripped`＋三枚 `.go.txt` 快照＋`msg-*.txt` commit 消息＋`tree-*.dir`）。⛔ 仓库根、⛔ `.out` 名（注：`logs/r359-*` 等一律 `.log`／`.txt`）；⛔ 在 `probes/**` 留 `.go`——中途落过一枚 `logs/bridge-head.go`，**同一次运行内**改名 `.go.txt` 并把脚本改成导到 `/tmp`，现 `find .scratch/wisp/probes/309/r2 -name '*.go'`＝**0**（尺见 `logs/r359-14-handin.log`）。
- **⛔ 删过任何临时件**：突变全部是"还原⇄⛔ 删除"（trap ＋ `cmp` ＋ `sha256sum`）；⛔ `git add -A`／`.`；⛔ `--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`／`--no-verify`；⛔ push。

## 4. 我认为编排者那份裁定里⛔ 成立的一处（顶回来请复跑再判）

§4 第 4 件与 §2 签收的 `Q-309-v1-prov-nil-limb-blind` 用的是"**两肢各 `rc=0`**"这个读数；整包那把尺下只有 `b.injected` 成立，`b.prov == nil` 是 `rc=1`（16 枚噪声红）。
⇒ **处置⛔ 变**（照样⛔ 造仪器、照样注释另记），变的是那句"同族同色"的说法：本腿注释写成"一肢全盲／一肢只有噪声式反馈、⛔ 具名仪器"。若编排者复跑 `logs/prov-rigs.txt` 那两发对得上，请把 §4 第 4 件的措辞按两把尺分口径记进台账；对不上就撤我这处。

## 5. 追加（交回之后本腿又跑的那一发，⛔ 改上面任何一句）

- 收尾四道门禁在**已 commit 的面**上又各跑一发： rc=0、 rc=0、d22scan.sh: positive control - runtests.sh -C tools/d22scan ./...
=== RUN   TestScanDetectsAllSeededViolations
--- PASS: TestScanDetectsAllSeededViolations (0.06s)
=== RUN   TestScanCleanRepoIsGreen
--- PASS: TestScanCleanRepoIsGreen (0.04s)
=== RUN   TestAllowlistSuppressesOnlyListedPaths
--- PASS: TestAllowlistSuppressesOnlyListedPaths (0.05s)
=== RUN   TestCheckRootRejectsBlindRoots
=== RUN   TestCheckRootRejectsBlindRoots/empty_dir
=== RUN   TestCheckRootRejectsBlindRoots/go.mod_only
=== RUN   TestCheckRootRejectsBlindRoots/the_d22scan_module_itself
=== RUN   TestCheckRootRejectsBlindRoots/repo_skeleton_with_two_go_files
--- PASS: TestCheckRootRejectsBlindRoots (0.05s)
    --- PASS: TestCheckRootRejectsBlindRoots/empty_dir (0.00s)
    --- PASS: TestCheckRootRejectsBlindRoots/go.mod_only (0.00s)
    --- PASS: TestCheckRootRejectsBlindRoots/the_d22scan_module_itself (0.01s)
    --- PASS: TestCheckRootRejectsBlindRoots/repo_skeleton_with_two_go_files (0.04s)
=== RUN   TestScanAloneIsNotAFalsifier
--- PASS: TestScanAloneIsNotAFalsifier (0.01s)
=== RUN   TestCheckRootAcceptsRealRepo
    scan_test.go:250: real repo production Go files in scope: 268
--- PASS: TestCheckRootAcceptsRealRepo (0.32s)
=== RUN   TestScannerSelfScanOfRealRepoIsGreen
--- PASS: TestScannerSelfScanOfRealRepoIsGreen (1.90s)
=== RUN   TestEmojiBanCoversGoSourcesNotJustDesign
--- PASS: TestEmojiBanCoversGoSourcesNotJustDesign (0.05s)
=== RUN   TestBan8MathBandAndRemainingGaps
=== RUN   TestBan8MathBandAndRemainingGaps/U+2265_greater-or-equal,_band_added_by_Q-46(c)
=== RUN   TestBan8MathBandAndRemainingGaps/U+2264_less-or-equal,_band_added_by_Q-46(c)
=== RUN   TestBan8MathBandAndRemainingGaps/U+2229_intersection,_band_added_by_Q-46(c)
=== RUN   TestBan8MathBandAndRemainingGaps/U+2212_minus_sign,_band_added_by_Q-46(c)
=== RUN   TestBan8MathBandAndRemainingGaps/U+2713_check,_pre-existing_U+2600-U+27BF_band
=== RUN   TestBan8MathBandAndRemainingGaps/U+1F600_emoji,_pre-existing_U+1F000-U+1FAFF_band
=== RUN   TestBan8MathBandAndRemainingGaps/U+FE0F_variation_selector,_pre-existing
=== RUN   TestBan8MathBandAndRemainingGaps/U+2192_right_arrow,_gap_kept
=== RUN   TestBan8MathBandAndRemainingGaps/U+21D2_rightdouble_arrow,_gap_kept
=== RUN   TestBan8MathBandAndRemainingGaps/U+2500_box_drawing,_gap_kept
=== RUN   TestBan8MathBandAndRemainingGaps/U+2460_circled_one,_gap_kept
--- PASS: TestBan8MathBandAndRemainingGaps (0.52s)
    --- PASS: TestBan8MathBandAndRemainingGaps/U+2265_greater-or-equal,_band_added_by_Q-46(c) (0.04s)
    --- PASS: TestBan8MathBandAndRemainingGaps/U+2264_less-or-equal,_band_added_by_Q-46(c) (0.06s)
    --- PASS: TestBan8MathBandAndRemainingGaps/U+2229_intersection,_band_added_by_Q-46(c) (0.05s)
    --- PASS: TestBan8MathBandAndRemainingGaps/U+2212_minus_sign,_band_added_by_Q-46(c) (0.04s)
    --- PASS: TestBan8MathBandAndRemainingGaps/U+2713_check,_pre-existing_U+2600-U+27BF_band (0.06s)
    --- PASS: TestBan8MathBandAndRemainingGaps/U+1F600_emoji,_pre-existing_U+1F000-U+1FAFF_band (0.04s)
    --- PASS: TestBan8MathBandAndRemainingGaps/U+FE0F_variation_selector,_pre-existing (0.05s)
    --- PASS: TestBan8MathBandAndRemainingGaps/U+2192_right_arrow,_gap_kept (0.04s)
    --- PASS: TestBan8MathBandAndRemainingGaps/U+21D2_rightdouble_arrow,_gap_kept (0.04s)
    --- PASS: TestBan8MathBandAndRemainingGaps/U+2500_box_drawing,_gap_kept (0.04s)
    --- PASS: TestBan8MathBandAndRemainingGaps/U+2460_circled_one,_gap_kept (0.05s)
=== RUN   TestBan8CommentExemptionInGoSources
=== RUN   TestBan8CommentExemptionInGoSources/line_comment
=== RUN   TestBan8CommentExemptionInGoSources/doc_comment_above_a_declaration
=== RUN   TestBan8CommentExemptionInGoSources/one-line_block_comment
=== RUN   TestBan8CommentExemptionInGoSources/multi-line_block_comment,_glyph_on_a_middle_line
=== RUN   TestBan8CommentExemptionInGoSources/trailing_comment_after_clean_code
=== RUN   TestBan8CommentExemptionInGoSources/string_on_the_same_line_as_a_trailing_comment
=== RUN   TestBan8CommentExemptionInGoSources/raw_string_that_reads_as_a_SQL_comment
=== RUN   TestBan8CommentExemptionInGoSources/interpreted_string_containing_a_comment_marker
=== RUN   TestBan8CommentExemptionInGoSources/glyph_inside_a_comment_AND_inside_code_in_one_block_comment
=== RUN   TestBan8CommentExemptionInGoSources/unparseable_file_gets_no_exemption_at_all
--- PASS: TestBan8CommentExemptionInGoSources (0.44s)
    --- PASS: TestBan8CommentExemptionInGoSources/line_comment (0.04s)
    --- PASS: TestBan8CommentExemptionInGoSources/doc_comment_above_a_declaration (0.04s)
    --- PASS: TestBan8CommentExemptionInGoSources/one-line_block_comment (0.04s)
    --- PASS: TestBan8CommentExemptionInGoSources/multi-line_block_comment,_glyph_on_a_middle_line (0.05s)
    --- PASS: TestBan8CommentExemptionInGoSources/trailing_comment_after_clean_code (0.04s)
    --- PASS: TestBan8CommentExemptionInGoSources/string_on_the_same_line_as_a_trailing_comment (0.04s)
    --- PASS: TestBan8CommentExemptionInGoSources/raw_string_that_reads_as_a_SQL_comment (0.04s)
    --- PASS: TestBan8CommentExemptionInGoSources/interpreted_string_containing_a_comment_marker (0.05s)
    --- PASS: TestBan8CommentExemptionInGoSources/glyph_inside_a_comment_AND_inside_code_in_one_block_comment (0.05s)
    --- PASS: TestBan8CommentExemptionInGoSources/unparseable_file_gets_no_exemption_at_all (0.04s)
=== RUN   TestBan8CommentExemptionInTextScopes
=== RUN   TestBan8CommentExemptionInTextScopes/tsx_line_comment
=== RUN   TestBan8CommentExemptionInTextScopes/tsx_block_divider,_glyph_on_the_closing_line
=== RUN   TestBan8CommentExemptionInTextScopes/html_comment
=== RUN   TestBan8CommentExemptionInTextScopes/css_comment
=== RUN   TestBan8CommentExemptionInTextScopes/JSX_text_node_with_a_math_glyph
=== RUN   TestBan8CommentExemptionInTextScopes/code_after_a_closed_block_comment_on_the_same_line
=== RUN   TestBan8CommentExemptionInTextScopes/markdown_heading_is_not_a_comment
=== RUN   TestBan8CommentExemptionInTextScopes/SQL-style_dash_is_not_a_comment_in_a_text_scope
=== RUN   TestBan8CommentExemptionInTextScopes/extension-less_file,_plain_text_line
--- PASS: TestBan8CommentExemptionInTextScopes (0.42s)
    --- PASS: TestBan8CommentExemptionInTextScopes/tsx_line_comment (0.05s)
    --- PASS: TestBan8CommentExemptionInTextScopes/tsx_block_divider,_glyph_on_the_closing_line (0.05s)
    --- PASS: TestBan8CommentExemptionInTextScopes/html_comment (0.05s)
    --- PASS: TestBan8CommentExemptionInTextScopes/css_comment (0.05s)
    --- PASS: TestBan8CommentExemptionInTextScopes/JSX_text_node_with_a_math_glyph (0.04s)
    --- PASS: TestBan8CommentExemptionInTextScopes/code_after_a_closed_block_comment_on_the_same_line (0.05s)
    --- PASS: TestBan8CommentExemptionInTextScopes/markdown_heading_is_not_a_comment (0.04s)
    --- PASS: TestBan8CommentExemptionInTextScopes/SQL-style_dash_is_not_a_comment_in_a_text_scope (0.05s)
    --- PASS: TestBan8CommentExemptionInTextScopes/extension-less_file,_plain_text_line (0.05s)
=== RUN   TestDeclaredEmojiScopeCannotWalkZeroFiles
--- PASS: TestDeclaredEmojiScopeCannotWalkZeroFiles (0.11s)
=== RUN   TestScopeReportMatchesRealCoverage
--- PASS: TestScopeReportMatchesRealCoverage (2.24s)
=== RUN   TestVerdictGreenOnFullyLiveFixture
--- PASS: TestVerdictGreenOnFullyLiveFixture (0.06s)
=== RUN   TestVerdictRedOnEmptyBan7Scope
--- PASS: TestVerdictRedOnEmptyBan7Scope (0.07s)
=== RUN   TestVerdictRedOnGoScopeWithOnlyTestFiles
--- PASS: TestVerdictRedOnGoScopeWithOnlyTestFiles (0.06s)
=== RUN   TestExemptScopeCannotOutliveItsAbsentTree
--- PASS: TestExemptScopeCannotOutliveItsAbsentTree (0.09s)
=== RUN   TestBan6ScopeIsNotNarrowedByAnExtensionFilter
    scan_test.go:1002: ban #6 examined 7 and fired in all 7 extension classes
--- PASS: TestBan6ScopeIsNotNarrowedByAnExtensionFilter (0.07s)
=== RUN   TestVerdictRedOnUndeclaredCounter
--- PASS: TestVerdictRedOnUndeclaredCounter (0.05s)
=== RUN   TestRealRepoBan8CoversFrontendTreeAtBan6sCount
--- PASS: TestRealRepoBan8CoversFrontendTreeAtBan6sCount (2.22s)
=== RUN   TestBan8FrontendScopeIsNotNarrowedByAnExtensionFilter
    scan_test.go:1155: ban #8 examined 7 and fired in all 7 frontend/ file classes, incl. .mjs, .md and extension-less
--- PASS: TestBan8FrontendScopeIsNotNarrowedByAnExtensionFilter (0.07s)
=== RUN   TestLedgerCountsMatchAnIndependentWalk
    scan_test.go:1261: verified bans #1-5 internal/: 229 files
    scan_test.go:1261: verified bans #1-5 cmd/: 39 files
    scan_test.go:1261: verified ban #7 internal/tools/: 23 files
    scan_test.go:1261: verified ban #8 internal/: 531 files
    scan_test.go:1261: verified ban #8 cmd/: 119 files
    scan_test.go:1261: verified ban #6 frontend/: 85 files
    scan_test.go:1261: verified ban #8 frontend/: 85 files
--- PASS: TestLedgerCountsMatchAnIndependentWalk (2.16s)
=== RUN   TestWalksSkipGitIgnoredPaths
    scan_test.go:1375: note: d22scan: skipped as git-ignored: 3 file(s) under 1 ignored director(ies) [internal/build/], decided by .gitignore (2 path(s)), frontend/.gitignore (2 path(s))
--- PASS: TestWalksSkipGitIgnoredPaths (0.37s)
=== RUN   TestTrackedPathsAreNeverSkippedByTheIgnoreFilter
--- PASS: TestTrackedPathsAreNeverSkippedByTheIgnoreFilter (0.89s)
=== RUN   TestFullTrackedListCoversWhatTheNarrowListCannot
--- PASS: TestFullTrackedListCoversWhatTheNarrowListCannot (1.05s)
=== RUN   TestIndexChangingBetweenTheTwoGitReadsAppliesNoRules
--- PASS: TestIndexChangingBetweenTheTwoGitReadsAppliesNoRules (1.19s)
=== RUN   TestIgnoreRulesAreOffWhenTheIndexCannotBeAsked
--- PASS: TestIgnoreRulesAreOffWhenTheIndexCannotBeAsked (0.11s)
=== RUN   TestGitIgnoreRuleSemantics
--- PASS: TestGitIgnoreRuleSemantics (0.01s)
=== RUN   TestRealRepoLedgerIsHonest
    scan_test.go:2207: ban #6 examined 85 frontend/ text files in the real repo
--- PASS: TestRealRepoLedgerIsHonest (2.06s)
=== RUN   TestBuiltBinaryGoesRedEndToEnd
=== RUN   TestBuiltBinaryGoesRedEndToEnd/seeded_violation_exits_1
    scan_test.go:2363: rc=1 stdout=d22scan: gitignore rules NOT APPLIED - git cannot be consulted in C:/Users/swq/AppData/Local/Temp/TestBuiltBinaryGoesRedEndToEndseeded_violation_exits_12635825691/001: `git rev-parse --show-prefix`: fatal: not a git repository (or any of the parent directories): .git; every path in every scope is being scanned, so the counts below may include build output (A207's machine-dependent denominator). This is the loud direction: a scanner that cannot ask git which paths are tracked does not get to skip any.
        d22scan: examined 16 production Go files under internal/ and cmd/ of C:/Users/swq/AppData/Local/Temp/TestBuiltBinaryGoesRedEndToEndseeded_violation_exits_12635825691/001
        d22scan: scope bans #1-5 internal/      examined  15 production Go files
        d22scan: scope bans #1-5 cmd/           examined   1 production Go files
        d22scan: scope ban #6 frontend/         examined   1 text files
        d22scan: scope ban #7 internal/tools/   examined   1 production Go files
        d22scan: scope ban #8 design/           examined   1 text files
        d22scan: scope ban #8 frontend/         examined   1 text files
        d22scan: scope ban #8 internal/         examined  15 Go files, comments and _test.go included
        d22scan: scope ban #8 cmd/              examined   1 Go files, comments and _test.go included
        internal/bad/leak.go:5: [bare-goroutine] bare `go worker(...)` is banned (D22/D38b, R16: named calls count too): use observe.Registry.Spawn (named, owner, recover boundary)
        stderr=d22scan: 1 finding(s); D22 bans are not negotiable (see PLAN.md D22, tools/d22scan/allowlist.txt)
=== RUN   TestBuiltBinaryGoesRedEndToEnd/empty_live_scope_exits_2
    scan_test.go:2363: rc=2 stdout=d22scan: gitignore rules NOT APPLIED - git cannot be consulted in C:/Users/swq/AppData/Local/Temp/TestBuiltBinaryGoesRedEndToEndempty_live_scope_exits_23528534762/001: `git rev-parse --show-prefix`: fatal: not a git repository (or any of the parent directories): .git; every path in every scope is being scanned, so the counts below may include build output (A207's machine-dependent denominator). This is the loud direction: a scanner that cannot ask git which paths are tracked does not get to skip any.
        d22scan: examined 14 production Go files under internal/ and cmd/ of C:/Users/swq/AppData/Local/Temp/TestBuiltBinaryGoesRedEndToEndempty_live_scope_exits_23528534762/001
        d22scan: scope bans #1-5 internal/      examined  13 production Go files
        d22scan: scope bans #1-5 cmd/           examined   1 production Go files
        d22scan: scope ban #6 frontend/         examined   1 text files
        d22scan: scope ban #7 internal/tools/   examined   0 production Go files
        d22scan: scope ban #8 design/           examined   1 text files
        d22scan: scope ban #8 frontend/         examined   1 text files
        d22scan: scope ban #8 internal/         examined  13 Go files, comments and _test.go included
        d22scan: scope ban #8 cmd/              examined   1 Go files, comments and _test.go included
        stderr=d22scan: scope ban #7 internal/tools/ examined 0 files but is declared live in declaredScopes() - an empty instrument is not a verdict (ticket 71 AC#4)
=== RUN   TestBuiltBinaryGoesRedEndToEnd/armed_ban_6_goes_red_on_the_panel_violation_exits_1
    scan_test.go:2363: rc=1 stdout=d22scan: gitignore rules NOT APPLIED - git cannot be consulted in C:/Users/swq/AppData/Local/Temp/TestBuiltBinaryGoesRedEndToEndarmed_ban_6_goes_red_on_the_panel1647708783/001: `git rev-parse --show-prefix`: fatal: not a git repository (or any of the parent directories): .git; every path in every scope is being scanned, so the counts below may include build output (A207's machine-dependent denominator). This is the loud direction: a scanner that cannot ask git which paths are tracked does not get to skip any.
        d22scan: examined 15 production Go files under internal/ and cmd/ of C:/Users/swq/AppData/Local/Temp/TestBuiltBinaryGoesRedEndToEndarmed_ban_6_goes_red_on_the_panel1647708783/001
        d22scan: scope bans #1-5 internal/      examined  14 production Go files
        d22scan: scope bans #1-5 cmd/           examined   1 production Go files
        d22scan: scope ban #6 frontend/         examined   2 text files
        d22scan: scope ban #7 internal/tools/   examined   1 production Go files
        d22scan: scope ban #8 design/           examined   1 text files
        d22scan: scope ban #8 frontend/         examined   2 text files
        d22scan: scope ban #8 internal/         examined  14 Go files, comments and _test.go included
        d22scan: scope ban #8 cmd/              examined   1 Go files, comments and _test.go included
        frontend/src/app.js:1: [panel-approval] `approval.decide` in frontend/ is banned (D33/F2: allow decisions are native-side only)
        stderr=d22scan: 1 finding(s); D22 bans are not negotiable (see PLAN.md D22, tools/d22scan/allowlist.txt)
=== RUN   TestBuiltBinaryGoesRedEndToEnd/frontend_tree_gone_while_declared_live_exits_2
    scan_test.go:2363: rc=2 stdout=d22scan: gitignore rules NOT APPLIED - git cannot be consulted in C:/Users/swq/AppData/Local/Temp/TestBuiltBinaryGoesRedEndToEndfrontend_tree_gone_while_declared2277227758/001: `git rev-parse --show-prefix`: fatal: not a git repository (or any of the parent directories): .git; every path in every scope is being scanned, so the counts below may include build output (A207's machine-dependent denominator). This is the loud direction: a scanner that cannot ask git which paths are tracked does not get to skip any.
        d22scan: examined 15 production Go files under internal/ and cmd/ of C:/Users/swq/AppData/Local/Temp/TestBuiltBinaryGoesRedEndToEndfrontend_tree_gone_while_declared2277227758/001
        d22scan: scope bans #1-5 internal/      examined  14 production Go files
        d22scan: scope bans #1-5 cmd/           examined   1 production Go files
        d22scan: scope ban #6 frontend/         examined   0 text files
        d22scan: scope ban #7 internal/tools/   examined   1 production Go files
        d22scan: scope ban #8 design/           examined   1 text files
        d22scan: scope ban #8 frontend/         examined   0 text files
        d22scan: scope ban #8 internal/         examined  14 Go files, comments and _test.go included
        d22scan: scope ban #8 cmd/              examined   1 Go files, comments and _test.go included
        stderr=d22scan: ban #8 scope frontend/ examined 0 files - it is declared in emojiScopes() but walks nothing. Point it at a real tree or delete the entry; never leave a scope pretending to scan (ticket 71 AC#4)
=== RUN   TestBuiltBinaryGoesRedEndToEnd/ban_7_tree_gone_while_declared_live_exits_2
    scan_test.go:2363: rc=2 stdout=d22scan: gitignore rules NOT APPLIED - git cannot be consulted in C:/Users/swq/AppData/Local/Temp/TestBuiltBinaryGoesRedEndToEndban_7_tree_gone_while_declared_li3036717788/001: `git rev-parse --show-prefix`: fatal: not a git repository (or any of the parent directories): .git; every path in every scope is being scanned, so the counts below may include build output (A207's machine-dependent denominator). This is the loud direction: a scanner that cannot ask git which paths are tracked does not get to skip any.
        d22scan: examined 14 production Go files under internal/ and cmd/ of C:/Users/swq/AppData/Local/Temp/TestBuiltBinaryGoesRedEndToEndban_7_tree_gone_while_declared_li3036717788/001
        d22scan: scope bans #1-5 internal/      examined  13 production Go files
        d22scan: scope bans #1-5 cmd/           examined   1 production Go files
        d22scan: scope ban #6 frontend/         examined   1 text files
        d22scan: scope ban #7 internal/tools/   examined   0 production Go files
        d22scan: scope ban #8 design/           examined   1 text files
        d22scan: scope ban #8 frontend/         examined   1 text files
        d22scan: scope ban #8 internal/         examined  13 Go files, comments and _test.go included
        d22scan: scope ban #8 cmd/              examined   1 Go files, comments and _test.go included
        stderr=d22scan: scope ban #7 internal/tools/ examined 0 files but is declared live in declaredScopes() - an empty instrument is not a verdict (ticket 71 AC#4)
=== RUN   TestBuiltBinaryGoesRedEndToEnd/fully_live_fixture_exits_0
    scan_test.go:2363: rc=0 stdout=d22scan: gitignore rules NOT APPLIED - git cannot be consulted in C:/Users/swq/AppData/Local/Temp/TestBuiltBinaryGoesRedEndToEndfully_live_fixture_exits_0208823468/001: `git rev-parse --show-prefix`: fatal: not a git repository (or any of the parent directories): .git; every path in every scope is being scanned, so the counts below may include build output (A207's machine-dependent denominator). This is the loud direction: a scanner that cannot ask git which paths are tracked does not get to skip any.
        d22scan: examined 15 production Go files under internal/ and cmd/ of C:/Users/swq/AppData/Local/Temp/TestBuiltBinaryGoesRedEndToEndfully_live_fixture_exits_0208823468/001
        d22scan: scope bans #1-5 internal/      examined  14 production Go files
        d22scan: scope bans #1-5 cmd/           examined   1 production Go files
        d22scan: scope ban #6 frontend/         examined   1 text files
        d22scan: scope ban #7 internal/tools/   examined   1 production Go files
        d22scan: scope ban #8 design/           examined   1 text files
        d22scan: scope ban #8 frontend/         examined   1 text files
        d22scan: scope ban #8 internal/         examined  14 Go files, comments and _test.go included
        d22scan: scope ban #8 cmd/              examined   1 Go files, comments and _test.go included
        d22scan: clean - no D22 ban violations; live scope work: bans #1-5 internal/=14, bans #1-5 cmd/=1, ban #6 frontend/=1, ban #7 internal/tools/=1, ban #8 design/=1, ban #8 frontend/=1, ban #8 internal/=14, ban #8 cmd/=1; ban #8 emoji coverage: design/ 1 text files; frontend/ 1 text files; internal/ 14 Go files, comments and _test.go included; cmd/ 1 Go files, comments and _test.go included
        stderr=
--- PASS: TestBuiltBinaryGoesRedEndToEnd (2.04s)
    --- PASS: TestBuiltBinaryGoesRedEndToEnd/seeded_violation_exits_1 (0.17s)
    --- PASS: TestBuiltBinaryGoesRedEndToEnd/empty_live_scope_exits_2 (0.15s)
    --- PASS: TestBuiltBinaryGoesRedEndToEnd/armed_ban_6_goes_red_on_the_panel_violation_exits_1 (0.15s)
    --- PASS: TestBuiltBinaryGoesRedEndToEnd/frontend_tree_gone_while_declared_live_exits_2 (0.15s)
    --- PASS: TestBuiltBinaryGoesRedEndToEnd/ban_7_tree_gone_while_declared_live_exits_2 (0.14s)
    --- PASS: TestBuiltBinaryGoesRedEndToEnd/fully_live_fixture_exits_0 (0.15s)
=== RUN   TestBan9PrefixTablesAreTheSameList
--- PASS: TestBan9PrefixTablesAreTheSameList (0.00s)
=== RUN   TestSelfTestEntryPassesEveryCase
--- PASS: TestSelfTestEntryPassesEveryCase (2.43s)
=== RUN   TestSelfTestRosterAuditRejectsAHollowTable
--- PASS: TestSelfTestRosterAuditRejectsAHollowTable (0.00s)
=== RUN   TestSelfTestRefusesToAuditAnEmptyRoster
=== RUN   TestSelfTestRefusesToAuditAnEmptyRoster/not-a-roster
=== RUN   TestSelfTestRefusesToAuditAnEmptyRoster/block-without-numbers
--- PASS: TestSelfTestRefusesToAuditAnEmptyRoster (0.01s)
    --- PASS: TestSelfTestRefusesToAuditAnEmptyRoster/not-a-roster (0.00s)
    --- PASS: TestSelfTestRefusesToAuditAnEmptyRoster/block-without-numbers (0.00s)
=== RUN   TestSelfTestFlagIsWiredInTheBuiltBinary
--- PASS: TestSelfTestFlagIsWiredInTheBuiltBinary (6.46s)
PASS
ok  	github.com/CarlosShao/wisp/tools/d22scan	27.708s
runtests.sh: OK - packages=[./...] top-level: PASS=35 FAIL=0 SKIP=0, === RUN=77, '[no tests to run]'=0
d22scan.sh: scan of /d/work/workspace/projects plans/Wisp
d22scan: skipped as git-ignored: 1 file(s) under 1 ignored director(ies) [frontend/dist/assets/], decided by frontend/.gitignore (2 path(s))
d22scan: examined 268 production Go files under internal/ and cmd/ of D:/work/workspace/projects plans/Wisp
d22scan: scope bans #1-5 internal/      examined 229 production Go files
d22scan: scope bans #1-5 cmd/           examined  39 production Go files
d22scan: scope ban #6 frontend/         examined  85 text files
d22scan: scope ban #7 internal/tools/   examined  23 production Go files
d22scan: scope ban #8 design/           examined  39 text files
d22scan: scope ban #8 frontend/         examined  85 text files
d22scan: scope ban #8 internal/         examined 531 Go files, comments and _test.go included
d22scan: scope ban #8 cmd/              examined 119 Go files, comments and _test.go included
d22scan: clean - no D22 ban violations; live scope work: bans #1-5 internal/=229, bans #1-5 cmd/=39, ban #6 frontend/=85, ban #7 internal/tools/=23, ban #8 design/=39, ban #8 frontend/=85, ban #8 internal/=531, ban #8 cmd/=119; ban #8 emoji coverage: design/ 39 text files; frontend/ 85 text files; internal/ 531 Go files, comments and _test.go included; cmd/ 119 Go files, comments and _test.go included rc=0、check-path-length-budget.sh: positive control bench=/tmp/tmp.ywj9k5y6Hx (kept on disk; this project never deletes temp artifacts)
check-path-length-budget.sh: control 1/3 ok - bench baseline green, so this gate is not stuck red
check-path-length-budget.sh: control 2/3 ok - the planted over-budget tracked path is rejected and named verbatim:
RED - over budget and NOT in the roster: .scratch/wisp/issues/seeded-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx-over-the-hat.md  (relative 125 chars, name 104 chars, full path on the self-hosted runner 169 chars, budget 165)
check-path-length-budget.sh: control 3/3 ok - green again once the planted path is gone
check-path-length-budget.sh: positive control PASSED
check-path-length-budget.sh: repo=D:/work/workspace/projects plans/Wisp
check-path-length-budget.sh: unit=characters (probe: length of the 5-byte sample "a U+25FF b" = 3; LANG=C.UTF-8 LC_ALL=C.UTF-8)
check-path-length-budget.sh: hat: rule 9 name cap=100 + issues dir prefix=21 -> relative hat=121 ; worst checkout prefix=44 -> full-path budget=165
check-path-length-budget.sh: check: relative_length + 44 > 165, or a last name component > 100
check-path-length-budget.sh: denominator read: 10657 tracked paths
check-path-length-budget.sh: denominator: tracked paths=10657  over-budget=57  covered by roster=57  not in roster=0
check-path-length-budget.sh: longest=184 chars relative (.scratch/wisp/issues/261-the-model-enabled-key-promises-removal-from-discovery-and-selection-but-no-production-code-reads-it-while-the-default-true-tag-is-inert-for-map-entries-done.md)
check-path-length-budget.sh: worst full path on the self-hosted runner=228 chars (hat budget 165, wall open interval (206,217])
check-path-length-budget.sh: bands: over the hat=57  of which in the 122..180 middle=56  past the old debt line(180)=1  in the wall interval=0  roster entries=57
check-path-length-budget.sh: VERDICT GREEN - every over-budget tracked path is rostered by name with a reason, and the roster equals the tree
check-path-length-budget.sh: unit cross-check: the byte view and the character view select the SAME 57 paths on this tree, recomputed just now in both units rc=0（件＝／／／，时刻 15:1x +08）。⇒ 与 §1 那组数⛔ 矛盾，是同一件事的第二发。
- 同批发现在共享索引里：本腿  之前索引已带着 **27** 枚别家路径（、、、票 308 的  改名）。⇒ 三笔 commit 全部 ，逐笔核过：本腿证据笔里⛔ 别家路径＝**0**（尺＝），那 27 枚原样留在索引里⛔ 动。名册＝（125 枚）与 （27 枚）。
-  与先落盘的  同题（⛔ 取代、⛔ 删，见规则 8），该件顶部已自带这一句。
