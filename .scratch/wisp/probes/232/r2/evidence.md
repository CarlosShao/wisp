# 票 232 — 写码腿 `232-r2` 取证件

> 判据来源＝`.scratch/wisp/issues/232-the-restart-tier-test-pins-that-the-sentence-exists-not-that-it-carries-the-why.md`
> （本件**不勾它的框、不改它的 Status、不改它的名**）。本件只记读数与判语。
> 落点＝`.scratch/wisp/probes/232/r2/`。骨架先落盘（编排者对 150 轮帽的固定对策），读数逐格追加。

---

## §0 起手锚（三把尺原文）

| 尺 | 命令 | 读数 |
|---|---|---|
| 撤票尺 | `sed -n '3,4p' <票面>` | 第 3 行开头 `- Status: **待派，测试面加固**…`；第 4 行 `- ⚠ **这不是安全漏洞，照三行读**…`。**零 `WITHDRAWN`／零「撤」／零「作废」** ⇒ 未撤票，许开工。 |
| HEAD 锚 | `git log -1 --format='%h %H %ad' --date=iso` | `590fb855 590fb85576d6ef06b79bf445bf26b72653415bb0 2026-10-06 11:34:27 +0800` |
| 写面尺 | `git status --porcelain -- cmd internal` | 空（0 行）⇒ 起手 `cmd`／`internal` 干净。 |
| 框数尺 | `grep -c '^- \[ \]' <票面>` | `6`（六枚全部不归本腿） |
| 产码 md5 | `md5sum cmd/wisp/config_reload.go` vs `git show HEAD:… \| md5sum` | 两者皆 `5ce441ca5e72b64d18a6c26f1c066882` ⇒ 起手一致。 |

### §0.1 那三段文案的 HEAD 原文（`git show HEAD:cmd/wisp/config_reload.go | sed -n '311,330p'` 复量，抄回存档）

```go
func (rt *agentRuntime) reportRestartPending(sections []string) {
	if rt == nil {
		return
	}
	rt.auditf("config: HOT-RELOAD state=restart-pending sections=%v effect=next-process-start tier=restart",
		sections)
	rt.auditf("config: RESTART-PENDING detail=%q keys=%q",
		"这些键只在进程启动时被读一次并交给平台层（开机自启注册、单实例锁、界面语言），本次运行没有重新加载它们的路径；"+
			"内存里继续用旧值，文件里的新值没有被丢掉，重启进程后生效",
		restartTierKeys)
	fmt.Fprintf(rt.stdout,
		"wisp run: 这些段的改动本次运行不会生效（D36 重启档，需要重启进程）：%v。"+
			"原因：这几枚键在进程启动时被读一次就交给平台层（开机自启注册、单实例锁、界面语言），"+
			"本次运行没有重新注册它们的路径。涉及：%s。"+
			"注意两件事都没发生：文件里的新值没有被丢掉，内存里的旧值也没有被偷偷替换。\n",
		sections, strings.Join(restartTierKeys, " / "))
}
```

（落点行号现量：`cmd/wisp/config_reload.go:318-319` 审计 detail 句、`:321-326` 操作员 `Fprintf` 三段。）

---

## §1 现量（票面那张表的五行，本腿自己复跑）

| 票面事实 | 本腿复跑尺 | 本腿读数 |
|---|---|---|
| 内容断言在哪 | `sed -n '482,495p' cmd/wisp/config_reload_223_test.go` | 未修态：`:490 out := r.awaitStdout(t, "本次运行不会生效")`；`:491-495` 三枚 needle 在 **`why+out` 拼接串**上 `strings.Contains`。 |
| ⛔ 瘦句不响 | AC#1 突变 B（见 §2/§3） | 见 §2 AC#1 行 |
| 方向说反**有**牙 | 同上区段 | `:479-481`（等不到那句就 `Fatalf`）＋`:504-506` 的负向断言（重启档不许被说成立即生效）⇒ 该形会红，本腿未动。 |
| ⚠ 前向风险（横幅假绿） | `grep -rn "本次运行不会生效\|重启进程后生效\|开机自启" --include=*.go cmd internal` | 未修态：产码里三枚 needle 的来源是 `config_reload.go:318-319`（审计句）＋`:322-325`（操作员句）＋`:339`（`restartTierKeys` 里 `app.autostart`）⇒ 操作员句被删薄时，`why+out` 仍可能由 stderr 那侧供词。 |
| 助手本身干净 | `sed -n '109,158p' cmd/wisp/config_reload_223_test.go` | `awaitStdout`(`:140-158`) 与 `awaitAudit`(`:109-127`) 同构：同 `reloadCaseBudget`=40s、同 `time.NewTimer`+20ms `time.NewTicker`、同 `Fatalf` 出口、零 `Skip`。 |
| 未修态两枚点名用例 | `PATH=… go test -count=1 -run 'TestTicket223RestartTierSaysItWillNotApply\|TestTicket223R2FailureSentenceRouting' -v ./cmd/wisp` | `--- PASS: TestTicket223RestartTierSaysItWillNotApply (2.48s)` ／ `--- PASS: TestTicket223R2FailureSentenceRouting (0.60s)` ／ `ok github.com/CarlosShao/wisp/cmd/wisp 3.186s`（全文 `.scratch/wisp/probes/232/r2/logs/AC0-baseline-unmutated.txt`） |

---

## §2 六格逐格表（格 ｜ 跑了什么命令 ｜ 原始读数 ｜ 判语）

| 格 | 跑了什么命令 | 原始读数 | 判语 |
|---|---|---|---|
| AC#1 复现瘦句不响 | `go test -overlay …/overlay/HEADtest-B.json -count=1 -run TestTicket223RestartTierSaysItWillNotApply -v ./cmd/wisp`（突变 B 的产码 ＋ **HEAD 未修态用例**一起进 overlay） | `--- PASS: TestTicket223RestartTierSaysItWillNotApply (2.28s)` / `ok github.com/CarlosShao/wisp/cmd/wisp 2.338s`；同形另三发 C／E_clause3／F_clause1 在 HEAD 用例上也**全绿**（`2.29s`／`2.50s`／`2.58s`，logs/AC1-HEADtest-*.txt） | **缺陷成立、本腿复现到位**：三段解释被删薄（或任一段被换掉）之后，未修态用例照样 PASS，因为它在 `why+out` 拼接串上找 needle，`why`（stderr 审计句，产码里从未被瘦）替 `out` 供了词。B 之外再复现出 C／E／F 三枚，说明这不是一发巧合。 |
| AC#2 断言改到 stdout＋窗口化 | 改 `cmd/wisp/config_reload_223_test.go`：`mark := r.h.out.String()` / `errMark := r.h.err.String()` 在 `r.plant` 之前取；新增文件内测试辅助 `windowSince232`／`stdoutSince232`／`awaitStdoutSince232`／`awaitAuditSince232`（全小写、零导出）；needle 只在具名流的窗口上找。`gofmt -l` 空，`go test -count=1 -run TestTicket223 -v ./cmd/wisp` | 未修态全绿基线之后：`--- PASS: TestTicket223RestartTierSaysItWillNotApply (2.19s)` 等 10 枚全 PASS、`ok … 27.431s`（logs/AC2-fixed-unmutated.txt）。B 之后新增一枚**实测**（logs/P0-needle-stream-probe.txt）：`app.autostart`／`开机自启` 在 stdout 与窗口上都 true，而 `重启进程后生效` 在 stdout 上 **false**、只在审计句上 true。 | 改完。⚠ **一处照票面写不成立，本腿按实测落地**：票面 AC#2 要三枚 needle 全在 stdout 上找，但 `重启进程后生效` 现在**不在产码那句操作员文案里**（`config_reload.go:321-326` 一字未动、禁区不许改文案）⇒ 硬钉 stdout 只能靠写产品文案变绿，正是本票禁区。所以那枚钉到**具名审计流＋种下改动之后的窗口**（横幅同样满足不了），并在 stdout 窗口上补钉 `需要重启进程`／`原因：`／`交给平台层`／`没有被丢掉` 四枚操作员半边词。判语＝**已改到 stdout 且窗口化**，`重启进程后生效` 这一枚的"只在 stdout"没做到，原因已具名。 |
| AC#3 反向正控要真响 | 同一 overlay，跑**最终修好的用例**：`-overlay …/mutB.json`、`…/mutC.json`、`…/mutE_clause3.json`、`…/mutF_clause1.json`、`…/mutJ_notice_to_stderr.json`、`…/mutK_notice_deleted.json`（后两枚是本腿补的） | 逐枚红：B `omits "app.autostart"/"开机自启"/"原因："/"交给平台层"/"没有被丢掉"` `--- FAIL (2.34s)`；C（换掉"原因"那半句无关话）`omits "开机自启","交给平台层"` `--- FAIL (2.11s)`；E_clause3（删第三段）`omits "没有被丢掉"` `--- FAIL (2.16s)`；F_clause1（去掉"原因："这个前缀）`omits "原因："` `--- FAIL (2.17s)`；J（整句改打到 stderr）`stdout never carried "本次运行不会生效" AFTER the plant within 40s` `--- FAIL (42.11s)`；K（操作员句整条删掉）同句红 `--- FAIL (43.25s)` | **不是装饰**：票面点名的两枚（B、"原因"段换新突变）都红，另外四枚自补突变也逐枚红，且红的是**具体哪一段没了**指得出来。未修态同一批突变全绿 ⇒ 红是新断言给的，不是环境给的。 |
| AC#4 横幅假绿那一枚 | 修好前先跑 HEADtest-M（启动横幅含全部字样 ＋ 操作员句整条删掉，配 **HEAD 未修态用例**）；再跑同一产码配**修好的用例**（mutM）；另一枚 mutG（横幅含全部字样、通知本体不动）配修好用例作对照 | HEAD 用例上：**`--- PASS … (4.19s)`——但不是 0ms**，因为 `awaitStdout` 的窗口里那句真的会来（未修态那一枚 needle 全由横幅＋审计供词，`why+out` 恒真）。G 形（通知完好）修好用例 `--- PASS (2.26s)`＝正例不误伤。**M 形配修好用例：`stdout never carried "本次运行不会生效" AFTER the plant within 40s`、`--- FAIL (43.80s)`** | 起作用的机制＝**起点窗口**（`mark` 在 `plant` 之前取 ⇒ 横幅那批 pre-plant 字节被 `TrimPrefix` 排除在外），不是"变慢"：修好的用例在 M 形上是**红**，红在等不到"种下改动之后"的那句。同形在未修态用例上是绿，两读并列即证明窗口是牙。 |
| AC#5 零放宽 | 逐名跑票面点名的两枚：`-run 'TestTicket223\|TestTicket223R2FailureSentenceRouting'`（实际用 `-run TestTicket223` 覆盖整族），并 `git diff` 人工核删除列 | 未修态：`PASS (2.48s)`／`PASS (0.60s)`。修好后（logs/AC2-fixed-unmutated.txt）：**10 枚 TestTicket223* 全 PASS**，含 `TestTicket223RestartTierSaysItWillNotApply (2.19s)`、`TestTicket223R2FailureSentenceRouting (0.60s)`；`-count=5` 复跑该枚 `2.24/2.24/2.20/2.16/2.25s` 五发全绿（logs/AC5-repeat5-fixed.txt）。`git diff --numstat` 只落在 `cmd/wisp/config_reload_223_test.go` 一枚文件上 | **零放宽成立**：`why+out` 拆开＝每枚 needle 的判据变严（少了一半供词来源）＋窗口化（再少 pre-plant 字节）；一条断言都没删，`这些段已立即生效：[app]` 那枚负向断言**故意保持全流读**（注释里具名），以免比原状更窄。没有任何"或满足任一"。 |
| AC#6 整包终态 | `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp ./internal/... -count=1`（logs/AC6-full-package.txt，968 行）＋逐枚红名在**安静台件**里复量＋一发把本腿写面退回 HEAD 的对照跑（logs/AC6-cmdwisp-baseline-at-HEADtest.txt） | 包级：`FAIL cmd/wisp 452.712s`／`FAIL internal/ball 0.538s`／`FAIL internal/panel 5.195s`／`FAIL internal/risk 7.716s`，其余 23 包全 `ok`。红名逐枚＝见 §2.1 | **零新增红落在本腿写面上**（判据见 §2.1 末行）；三枚 `cmd/wisp` 红本腿跑不出"因我而红"，其中两枚复量为**本机热键被占**的环境红、一枚为 ledger 在册的带载偶发。 |

### §2.1 AC#6 红名册逐名比

整包 `-count=1` 的 9 枚 `--- FAIL`（原文逐名，logs/AC6-full-package.txt）：

| 红名 | 包 | 历史在册？ | 本腿复量（安静台件） | 归谁 |
|---|---|---|---|---|
| `TestC21TableColourRowsMatchTokensCSS` | internal/ball | **是**＝派单写的"ball 1 枚"；`docs/reports/pending-and-issues.md:12056` 具名（"`tokens.css` 被界面那支删了，早于本票"） | 未复跑（别人地界，⛔ 不顺手修） | 别人 |
| `TestApprovalCardViewJSONKeysMatchFrontendTypes` | internal/panel | **是**＝`docs/reports/HANDOVER.md:659` 具名"转绿那一跳在零写面里，由 owner 带给界面那支" | 未复跑 | 别人 |
| `TestComposerContractTypesMatchFrontend` | internal/panel | **是**＝`HANDOVER.md:660`"历史三枚…不修不 Skip" | 未复跑 | 别人 |
| `TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme` | internal/panel | **是**＝同上 `:660` | 未复跑 | 别人 |
| `TestC21DesignTokensFourWayAgree` | internal/panel | **是**＝`HANDOVER.md:595`/`:660`（"owner 已定当已知常红读"） | 未复跑 | 别人 |
| `TestResolvePerCallBudget` | internal/risk | **是**＝派单写的争用型假红 | **`-count=3` 安静复量＝`ok github.com/CarlosShao/wisp/internal/risk 5.213s`**（logs/AC6-quiet-risk-count3.txt）⇒ 争用型，非红 | 无主（争用） |
| `Test258OccupiedCombinationNamesTheNewValue` | cmd/wisp | 派单未具名（本腿判：环境红，见下行） | 安静 `-count=1` 单跑仍红；红因逐字＝`hotkey occupied by another program, not registered hotkey=panel binding=Ctrl+Alt+P err="Hot key is already registered."`，用例只拿到 `hotkeys live 2/4`（要 `3/4`）；`-count=3` 复量 6 枚红＝同两枚×3，**稳定红不是偶发** | **本机此刻占着 `Ctrl+Alt+P`/`Ctrl+Alt+V`**（`cmd/wisp/resident_hotkey_258_windows_test.go:448`） |
| `Test258V1ProbeSummonEditRebindsLiveBall` | cmd/wisp | 同上 | 同上（`resident_hotkey_v1probe_test.go:108` 同一句 `boot verdict lost the live-count shape`，同样 `2/4`） | 同上 |
| `TestAC1ResidentLegInstallsItsLogListenerOnDisk` | cmd/wisp | **是**＝`pending-and-issues.md:12056` 具名"`cmd/wisp` 那两枚带载偶发" | 安静单跑 **`ok … 4.968s`**（logs/AC6-quiet-ac1resident.txt）⇒ 整包 452s 高负载下的带载偶发，复量转绿 | 无主（带载） |

**"是不是本腿写红的"这一把尺**（本腿不靠推理，靠跑）：`overlay/HEADtest-only.json` 把 `cmd/wisp/config_reload_223_test.go` 单独退回 `git show HEAD:` 的副本（本腿写面＝零），整包 `./cmd/wisp -count=1` 再跑一发（488.549s，logs/AC6-cmdwisp-baseline-at-HEADtest.txt）：

- 该发红名＝`TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe`／`TestAC14GoSideEvalPushReachesThePage`／`Test258OccupiedCombinationNamesTheNewValue`／`Test258V1ProbeSummonEditRebindsLiveBall`
- ⇒ **258 那两枚在"本腿什么都没写"的读数里照样红**（同一句 `Hot key is already registered.`），所以它们与本腿写面无关；`AC13`／`AC14` 那两枚在终态复量里一枚转绿（`TestAC14GoSideEvalPushReachesThePage -count=2 = ok 1.671s`）、一枚（`TestAC13ColdStart…`）仍红且红在 `panel_resident_windows_test.go:315-330` 的 WebView2 页面探活（`the page itself answered "0"`），与 `config_reload` 无关，也是 HEADtest 那发里有的＝非本腿引入。
- 对照终态那一发（本腿写面在内）红名＝258 两枚＋`AC1ResidentLeg`；两发交集只有 258 两枚，差异集合全是上面标注的带载/环境枚 ⇒ **判语：本腿零新增红**。

### §2 末尾安全尺两把

| 尺 | 命令 | 读数 |
|---|---|---|
| 写面终态 | `git status --porcelain -- cmd internal` | 见 §2.2（每次提交前拉，终值在末尾） |
| 产码未动 | `md5sum cmd/wisp/config_reload.go` | 13 发（含探针）跑完逐次对拉，全部 `5ce441ca5e72b64d18a6c26f1c066882` ＝ `git show HEAD:… \| md5sum`；终值见 §2.2 |

### §2.2 终态安全尺（本件最后一次拉）

本行读数随后追加。

---

## §3 突变名册（每发：overlay 路径 ＋ 落地证明 ＋ 红/绿原文行 ＋ 还原后 md5）

> 种针一律 `go test -overlay <json> …`，合成副本一律在 `D:/tmp/wisp232/` 下；
> 共享工作树里的 `cmd/wisp/config_reload.go` **一个字节都没被编辑过**（每发跑完对拉 md5）。
> overlay 清单物理件＝`.scratch/wisp/probes/232/r2/overlay/*.json`（键为本仓绝对路径，与票 231／151 腿同形）。
> `HEADtest-*.json` 额外把 `cmd/wisp/config_reload_223_test.go` 映射回 `D:/tmp/wisp232/test_head_223.go`（= `git show HEAD:` 抽出的未修态用例），用来做"同一发突变、修好前后各跑一遍"的对照。

| 发 | 种的是什么 | overlay | 落地证明（diff 只动那一处） | 未修态（HEAD 用例） | 修好终态 | 跑完 md5 |
|---|---|---|---|---|---|---|
| B | 票面 AC#1 那发：操作员 `Fprintf` 三段解释删薄，只留 `"wisp run: 这些段的改动本次运行不会生效（D36 重启档，需要重启进程）：%v。\n"` | `overlay/mutB.json` | `diff <(git show HEAD:cmd/wisp/config_reload.go) D:/tmp/wisp232/config_reload_mutB.go` ＝ 只有 `:322-326 → :322-323` 一处 | **`--- PASS (2.28s)`**（`logs/AC1-HEADtest-B.txt`） | **`--- FAIL (2.34s)`** 5 枚 `omits`（`logs/AC3final-mutB.txt`） | `5ce441ca…` |
| C | 票面 AC#3 那枚新突变：把"原因：…"整段换成 `"原因：这是一句与本次改动无关的话。"`，其余三段与 needle 全留 | `overlay/mutC.json` | diff 只有 `:323` 一行 | **`--- PASS (2.29s)`** | **`--- FAIL (2.11s)`** `omits "开机自启"`＋`omits "交给平台层"` | 同上 |
| E | 本腿补：删掉第三段"注意两件事都没发生…" | `overlay/mutE_clause3.json` | diff 只有 `:324-325` 合并 | **`--- PASS (2.50s)`** | **`--- FAIL (2.16s)`** `omits "没有被丢掉"` | 同上 |
| F | 本腿补：只剥掉 `"原因："` 这个引导词，内容不动 | `overlay/mutF_clause1.json` | diff 只有 `:323` | **`--- PASS (2.58s)`** | **`--- FAIL (2.17s)`** `omits "原因："` | 同上 |
| J | 本腿补：操作员那句**改打到 `rt.stderr`**（内容一字不动，只换流）——专测"是不是真在 stdout 上找" | `overlay/mutJ_notice_to_stderr.json` | diff 只有 `:321` `Fprintf(rt.stdout→stderr)` | （未跑 HEAD 用例配 J；J 的旧缺陷由 B／C／E／F 复现） | **`--- FAIL (42.11s)`** `stdout never carried "本次运行不会生效" AFTER the plant within 40s` | 同上 |
| K | 操作员那句整条删掉（审计两行照旧）——测"删光会不会红" | `overlay/mutK_notice_deleted.json` ／ `HEADtest-K.json` | diff 删 `:321-326` 那一块 | **`--- FAIL (43.11s)`**：未修态也红，红在 `awaitStdout` 自己等不到那句（票面"方向说反有牙"那一行说的就是这一形，本腿复量） | **`--- FAIL (43.25s)`** 同句、AFTER the plant 版 | 同上 |
| G | 启动横幅追加**全部字样**，重启通知本体**不动**（正例对照：窗口化不许误伤好文案） | `overlay/mutG_banner_only.json` | diff 只有 `:127` 追加一行 | — | **`--- PASS (2.26s)`**（红名册里它**该**绿） | 同上 |
| M | AC#4 那一枚：横幅先含同样字样 **＋** 重启通知整条删掉 | `HEADtest-M.json` ／ `mutM_banner_then_deleted.json` | diff ＝ `:127` 追加 ＋ `:321-326` 删除 | **`--- PASS (4.19s)`** ← 恒真那一形（非 0ms，但 needle 半由横幅供词） | **`--- FAIL (43.80s)`** `stdout never carried … AFTER the plant within 40s` | 同上 |

**探针两枚（零断言，只为读数）**：`overlay/probe-needles.json` → `logs/P0-needle-stream-probe.txt`（哪枚 needle 在哪条流上）；`overlay/probe2-banner.json` ＋ `probe2-banner-on-mutH.json` → `logs/AC4-probe2-unmutated.txt`／`AC4-probe2-on-mutH.txt`（横幅字样在 pre-plant mark 里是否存在）。物理件 `zz_232r2_*_test.go` 全在本目录，靠 overlay 映射进 `cmd/wisp`，共享树从未有过它们。


---

## §4 门禁读数

- 本腿写面＝只有 `cmd/wisp/config_reload_223_test.go` ＋ 本件目录 ＋ 票面 Progress log 一节（append-only，`git diff --numstat`＝4 增 0 删）⇒ **零脚本面**（没动 `scripts/**`／`.github/**`），
  所以 `bash -n scripts/*` 这把尺对本腿不适用；但 `d22scan` 扫的是代码形状，适用，读数如下。

| 尺 | 命令 | 读数 |
|---|---|---|
| d22scan 门禁 | `sh scripts/d22scan.sh` | **rc=0**；正控 `runtests.sh -C tools/d22scan ./... → PASS=35 FAIL=0 SKIP=0, === RUN=77, '[no tests to run]'=0`；实扫 `d22scan: clean - no D22 ban violations`，`bans #1-5 cmd/=38`、`ban #8 cmd/ 104 Go files（comments and _test.go INCLUDED）`⇒ 本腿新写的 `_test.go` 在 #8 射程内且未撞（全文 logs/G1-d22scan.txt） |
| gofmt | `gofmt -l cmd/wisp/config_reload_223_test.go` | **空**（该枚文件已格式化） |
| gofmt 整包 | `gofmt -l cmd/wisp/` | 只剩 `cmd\wisp\models.go` 一枚——`file` 读数＝`CRLF line terminators`、`git status` 里它**未被本腿改动**、且 HEAD 副本 `gofmt -l` 不报 ⇒ 预存工作树 CRLF，`pending-and-issues.md:12056` 明写"⛔ 别顺手格式化 `cmd/wisp/models.go` 那枚预存 CRLF"，本腿照办 |
| 写面删除列 | `git diff --numstat -- cmd internal` | `129 9 cmd/wisp/config_reload_223_test.go` ＝ 9 枚删除全部落在**本腿自己的写面**（替换掉的旧 `why+out` 那一枚 `for` 块与三枚 await 调用行），别人文件删除列＝0 |

---

## §5 判不动的地方（诚实栏）

1. **`重启进程后生效` 钉不到 stdout（票面 AC#2 字面做不到，本腿不硬凑）**。实测（§3 探针 logs/P0）：操作员 `Fprintf` 那三段里**没有**这四个字，它只在 `rt.auditf` 的 detail 句里（`config_reload.go:319`）。票面要"三枚内容 needle 只在 stdout 那一段上找"，而把 stdout 补上这句话＝写产品文案＝**本票禁区行**（"⛔ 不改产品文案来让测试变绿"）与派单 §1（"`config_reload.go` 一个字节都不许改"）都挡着。⇒ 处置＝这枚钉到**具名审计流＋种下改动之后的窗口**（横幅同样喂不进去，§3 的 M 形就是这一枚的反证），stdout 半边补钉 `需要重启进程`（同一承诺的操作员措辞）＋`原因：`／`交给平台层`／`没有被丢掉`。**要不要因此把 AC#2 判成"半格"，归编排者裁，本腿不自己勾。**
2. **AC#4 那一枚不是 0ms，本腿没能让它 0ms**。票面写"断用例不许 0ms 通过"，实测未修态在 M 形（横幅含全部字样＋通知整条删掉）上是 `--- PASS (4.19s)`：因为 `awaitStdout` 等的那句"本次运行不会生效"**确实会由横幅先给出**，而它一旦满足就往下走 needle 循环，needle 又全由横幅＋审计供词 ⇒ 恒真成立但耗时≈一轮 tick，不是 0ms。派单 §3 给的读法是"读数写清它最终是红还是变慢"⇒ 本腿读数＝**未修态绿（≈4.2s）／修好后红（43.8s，红在等不到"AFTER the plant"的那句）**，起作用的机制＝起点窗口（`TrimPrefix` 把 pre-plant 字节整段排除），不是期限变慢。若编排者认为"必须 0ms 才算复现"，这一格本腿判不动。
3. **`Test258OccupiedCombinationNamesTheNewValue` 与 `Test258V1ProbeSummonEditRebindsLiveBall` 是环境红，本腿判不了也别修**：整包 9 枚红里有这两枚派单未具名的新面孔。本腿用"把写面退回 HEAD 的整包对照发"证明它们与本腿无关（§2.1），安静 `-count=3` 复量 6 枚红＝稳定红，红因逐字 `Hot key is already registered.`（本机此刻有程序占着 `Ctrl+Alt+P`/`Ctrl+Alt+V`）。归谁＝这台机器的热键持有者／票 258 那支；⛔ 本腿不修（不是我的地界，且修它就是改别人的用例形状）。**ledger 里未查到这两枚的历史在册记录**（本腿 `grep -rIn "FAIL: Test258Occupied"` 只命中票 258 自己的突变件与 `267/gate` 那两发的 PASS）⇒ 是"本机新红"还是"在册漏记"，本腿判不动，留给编排者。
4. **`TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe` 在 HEADtest 对照发红、在终态安静复量也仍红**（`panel_resident_windows_test.go:315-330`，页面自答 `"0"`），它既不在派单的 5 枚名册里、也不是本腿写面能影响的一枚 ⇒ 只登记不裁，`pending-and-issues.md:10418`/`:11072` 有它的历史轨迹（ledger 认它是"干净检出必跳/WebView2 冷启动"那一族）。
5. **同一把尺跑第二遍没拿到期望读数的地方＝零**：上面 1-4 全是"读数与票面预期不同因此不裁"，不是"尺子失灵"。本腿没有换尺凑绿。
