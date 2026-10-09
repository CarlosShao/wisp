# reconcile-10｜只读对账腿：把台账 A751–A771 的承重句往前复跑（2026-10-09）

- 起手锚 `2415bbb92b9673e4e823604ff1c0701e4d0d41ce`（`dev`，`2026-10-09 12:01 +0800`）；**收尾锚已漂到 `fc359af8`**（别的腿在本程期间 commit）⇒ 一切"名册/计数"读数都在行内记了自己那枚锚。
- 起手快照 `git status --porcelain -- docs .scratch/wisp/issues`＝**空**（当时无人在飞）。跟踪件一律走对象层 `git show <ref>:<path>`；判"在不在"只认 `git cat-file -e`／`git ls-tree -r --name-only`。
- **自报**：本程**零 `go` 子命令**（无 `go build/vet/test/gofmt/gofumpt`、零门禁脚本）。唯一越界一格＝为复跑 A756 的"工作树 CR 列"对 7 枚 `.go` 跑了 `tr -cd '\r' | wc -c`（只读）。零翻框、零 `A##`、零 push、写点只有本件。

## ① 抽了哪些承重句

编号＝`A##:台账行号`（行号按 `git show HEAD:` 现读，A751 起于 `:14509`、A771 止于 `:14751`，全件 14751 行）。
射程内共抽 **84 句**：翻框/Status/名册类 21、`文件:行` 锚 39、commit↔文件形状（`--numstat`／`wc -l`）14、全称式句子 6、门禁与主机侧事实 4。跳过＝纯叙述与判语。

## ② 三档清单

### ✅ 对上（58 枚，逐条带锚＋尺＋读数）

| 句（A##:行） | 尺（逐字） | 锚 | 读数 |
|---|---|---|---|
| A751:14511 翻框恰 10／零反向 | `git diff a99f9a39 58a4b2fe -- .scratch/wisp/issues/ \| grep -c '^-- \[ \]'`（＋`grep -c '^+- \[x\]'`） | 两 ref | 10／10，`+[]`与`-[]`反向＝0 |
| A751:14511 逐对同文 AC#3×2/#4×2/#5×4/#6×1/#7×1 | `grep -oE '^\+- \[x\] \*\*AC#[0-9]' \| sort \| uniq -c` | 两 ref | 2/2/4/1/1＝10 ✓ |
| A751:14512 Status 逐票 92/97/104 1→1、105/110 1→2、113 BASE 已是 2 | `git show <ref>:<票> \| grep -cE '\*\*Status'` | a99f9a39↔58a4b2fe | 六枚逐一对上 |
| A751:14512 `07`/`115` 一字未动 | `git diff --name-only a99f9a39 58a4b2fe -- .scratch/wisp/issues` | 两 ref | 恰 6 文件（92/97/104/105/110/113） |
| A751:14511 `-done` 名册仍 98 | `git ls-tree -r --name-only 58a4b2fe .scratch/wisp/issues \| grep -cE '\-done\.md$'` | 58a4b2fe（腿的收口） | **98** ✓（朴素 `grep -c -- '-done'`＝**99**＝281 自己 slug 里带 `-done` ⇒ 尺形警告） |
| A751/A756 乙类 8 张 未勾 14→4／已勾 26→36 | 逐票 `grep -cE '^[[:space:]]*- \[ \]'`／`- \[x\]` 求和，名册＝07/92/97/104/105/110/113/115 | c39e2853↔542d7661↔HEAD | 14/26 → 4/36 → 4/36 ✓（⇒ A751 那句"已勾 14→24"确为错，见 ❌7） |
| A752:14516 尺件 227 insertions 零删除 | `git show --numstat --format='' 039ec93c` | 039ec93c | `…identity_rulers_test.go 227/0` ✓ |
| A752:14517 种 A 位 `cancel.go:69` | `git show HEAD:internal/tools/cancel.go \| sed -n 69p` | HEAD | `return h.taskID` ✓ |
| A752:14518 种 C 位 `task.go:708` | 同尺 | HEAD | `caller := TaskID(ctx)` ✓ |
| A752:14520 `归口` 在 README＝0（rc=1） | `git show bf9974c2:.scratch/wisp/issues/README.md \| grep -c 归口` | bf9974c2 | 0／rc=1 ✓ |
| A752:14520 `TaskID(ctx)` 测试面唯一本尺 | `git grep -l "TaskID(ctx)" 039ec93c -- 'internal/tools/*_test.go'` | 039ec93c | 恰 1 文件 ✓ |
| A753:14529 exe 旧 31h23m | `ls -l build/wisp.exe` ＋ `git log -1 --pretty=%ad aa6c1881` | 工作树/HEAD | 10-07 11:57 ↔ 10-08 19:20 ✓ 逐字 |
| A753:14529 `%APPDATA%\wisp` 只有 `logs/` | `ls "$APPDATA/wisp"` | 本机 | `logs` 一枚 ✓ |
| A753:14529 `runConsoleLoop` 零台件 | `git grep -c "runConsoleLoop" HEAD -- '*_test.go'` | HEAD | rc=1 空 ✓ |
| A753/A760 命令形状 `:40-41` 无 `-count=1` | `sed -n '39,41p'` 该台件 | HEAD | 逐字两行注释、确无 `-count=1` ✓ |
| A754:14534 AskOnTaskRoot 定义 `:698`／体内唯一调用 `:701` | `sed -n '697,702p'` | HEAD | ✓✓ |
| A754:14534 三枚 winlive 用例 `:115/:220/:311` | 逐行 `sed -n Np` | HEAD | 三处皆 `a, w := ra.AskOnTaskRoot(residentCard{` ✓ |
| A754:14535 `Gate.Replay` 定义 `gate.go:755`＋C18 注释 `:751-754` | `sed -n '751,756p'` | HEAD | ✓（"C18 一键重放…NOT an answer"逐字在位） |
| A754:14535 台账旧锚 `:11139` 写 `gate.go:749` | `sed -n 11139p \| grep -o 'gate\.go:[0-9]*'` | HEAD | `gate.go:749` ✓（漂 +6 成立） |
| A754:14535 `PLAN.md:1368` 含"一键重放／不取消" | `git grep -n 一键重放 HEAD -- docs/PLAN.md` | HEAD | 命中 `:1368`（C18 行）✓ |
| A754:14535 两把零调用者尺 | `git grep -nE "[A-Za-z0-9_]\.AskOnTaskRoot\(" HEAD -- ':(exclude).scratch' '*.go' ':!*_test.go'`（`.Replay(` 同形） | HEAD | 皆 **rc=1 空** ✓ |
| A754/A757/A760 `ui.go:143-162`／`:146` 三参／`:167-171` | 真身 `internal/agent/approval/ui.go`，逐行取 | HEAD | NativeAPI 143-162 ✓／`Allow(ctx, correlationID, grant)` ✓／PanelAPI 167-171 无 Allow ✓ |
| A754:14536 boot 报告 `resident_windows.go:269-270` | `sed -n '268,271p'` | HEAD | `任务来源：%s` ✓ |
| A754:14537 `approval_reply.go:566-584` 无 `replay` | `awk 560-590 \| grep -ci replay` | HEAD | 0 ✓ |
| A755:14550 仓库根无跟踪 `probes/` | `git ls-files probes \| wc -l` | HEAD | 0 ✓ |
| A755:14551 258 relative 165／261 relative 184／全树最长 | `git ls-files \| awk 最长` | HEAD | 165 ✓／184 ✓／全树 max＝261 那枚 184 ✓ |
| A755:14552 `R[...]` 名册 57 枚＋理由串同族 | `grep -cE '^\s*R\[' scripts/check-path-length-budget.sh` | HEAD | 57 ✓（`doubles as the -done anti-double-claim key` 20 行 ✓） |
| A756:14561 七枚表"工作树 CR／HEAD blob CR"两列 | `tr -cd '\r' <工作树> \| wc -c` vs `git show HEAD:<f> \| tr -cd '\r' \| wc -c` | HEAD | 334/131/225/1115/1284/0/0 ↔ 全 0 ✓ **逐格复现** |
| A756/A758 `bridge.go` CR 1284＝`wc -l` | 同尺＋`git show HEAD:… \| wc -l` | HEAD | 1284＝1284 ✓ |
| A756:14569 六票 numstat 92 4/2、97 2/1、104 4/2、105 4/1、110 6/2、113 4/2 | `git diff --numstat a99f9a39 58a4b2fe` | 两 ref | 六枚逐一 ✓ |
| A756:14571 票 97"已勾 3"三 ref 一致 | 锚定式尺 @c39e2853／226d3717^／6d22dd6c | 三 ref | 3／3／3 ✓ |
| A756:14566 ci.yml 七枚行号 | 逐行 `sed -n` | HEAD | 172 步名／173 `if: !cancelled()`／175 `go install …gofumpt`／176 `gofumpt -l . tools/d22scan tools/mockllm`／189 步名／239 `if`／240 `attrib.sh --tracked-only` ✓ |
| A757:14576 `7789f953`＝2 件 307（200＋107） | `git show --numstat`＋`wc -l` | 该 commit | ✓ |
| A757:14585 四枚漂后锚 | 逐行 `sed -n` | HEAD | `subagent_selfapproval_197_test.go:116`／`approval_reply.go:229`／`:274`／`gate.go:755` ✓ |
| A757/A760/A762 八枚裁决锚 | 逐行 `sed -n` | HEAD | `queue.go:161/162/177/414/423`、`queue_test.go:167`、`ticket259_denial_rulers:139/252`、`ticket259_panel_capability:301/308`、`ticket242_binding:171`、`approval.go:292/303/558/569` ✓ |
| A757:14587 `l2_grant_boundary_test.go` 在位 | `git cat-file -e HEAD:internal/panel/l2_grant_boundary_test.go` | HEAD | 存在 ✓ |
| A757 `sameDigest` 整文件不存在 | `git show HEAD:…ticket242_binding_test.go \| grep -c sameDigest` | HEAD | 0 ✓（＝内容不符而非漂） |
| A758:14593 三枚件 13/148/33 行 | `wc -l` | HEAD | ✓✓✓ |
| A758:14601 `internal/agent/subagent_197.go` 不存在／tools 版 `:262` | `git cat-file -e`＋`sed -n 262p` | HEAD | ABSENT ✓／`parentID := TaskID(ctx)` ✓ |
| A758:14598 现量 `:219-222` 是断言 | `sed -n '212,222p'` | HEAD | ✓（非空＋`== res.TaskID` 红＋前缀） |
| A759:14605 件 135 行＋commit 形状 | `git show --numstat ed8a7358` | ed8a7358 | 135/0 ✓ |
| A759:14606 `run.go` 最后提交 `613606c0`（10-02）／`a1b19873^..HEAD` 空 | `git log -1 -- cmd/wisp/run.go`＋`git diff \| wc -l` | HEAD | ✓／0 行 ✓ |
| A759:14606 `run.go:991`＝`cfg := rt.cfg`／名册行 `config_readers_255.go:110`／用例 `:179` 是函数声明行 | 逐行 `sed -n` | HEAD | ✓✓✓ |
| A760:14616 新尺 185 行＋`61e13e92` | `--numstat`＋`wc -l` | 该 commit | ✓ |
| A760 F1（`:14511` 逐字仍在） | `sed -n 14511p \| grep -oE '已勾 \*\*[0-9]+→[0-9]+\*\*'` | HEAD | `已勾 **14→24**` **仍在原位** ✓（＝F1 判语成立） |
| A760 F3 票 284 名长 66 | `ls-tree \| grep /284- \| awk length` | HEAD | 66 ✓（A754 写 62＝错，见 ❌8） |
| A760 F4 `residentCard` 5 文件／24 行 | `git grep -c residentCard HEAD -- '*.go'`（剥 `.scratch`） | HEAD | 4/6/6/2/6＝**24 行／5 文件** ✓ |
| A760 F5 全仓 9 件／22 行 | `git grep -lE "整链尺\|回退形无仪器" <ref>`＋`-n` | bf9974c2／6c969dc8 | 9/22 ✓（同尺在 ce18b3d6 起＝11/28＝后来加节所致） |
| A760/A761 五枚节头↔提交时刻 | `git log --date=format:%H:%M -- docs/reports/pending-and-issues.md` | HEAD | A756 10:22／A757 10:29／A758 10:42／A759 10:52／A760 11:01 全对 ✓ |
| A762:14646 件 77 行＋`871d53e0` | `--numstat`＋`wc -l` | HEAD | ✓ |
| A763:14659 `:151`／`:218` 两行逐字 | `sed -n` | HEAD | `newResidentPanelManager` ✓／`return panel.RequestToggle(via)` ✓ |
| A763:14656 未勾 64／已勾 33／行数 1393／14 枚全未 `-done` | 逐票求和（名册由 `chain-map.md` §1 抽出＝33/145/185/189/190/220/246/247/248/252/279/282/284/285） | **0091b1b9** | 64／33／**1393**／0 ✓（腿的 1433 复现不出＝采编排者） |
| A764:14668 六张表全长 92b347/97 51/104 306/105 232/110 43/113 134 | `wc -l` | HEAD | 六枚 ✓ |
| A764:14673 `92-:147`＝FAIL 那行／115 不在 s1／承载表 :16 表头 :17 分隔 :18 首行 | `sed -n 147p`／`ls-tree docs/evidence/s1 \| grep -c '/115-'`／`awk 16-18` | HEAD | ✓／0 ✓／✓ |
| A764/A766 承重表六枚 sha | `git log -1 --pretty=%h <sha>` ×6 | HEAD | 六枚全存在 ✓ |
| A765:14679 两枚新尺 143／113 | `wc -l` | HEAD | ✓（＝采现量那侧） |
| A765:14681 既有尺早作差 `corr_percall_242_test.go:99-105` | `awk 99-105 \| grep -c 'c1\|c2'` | HEAD | 5 行命中 c1/c2 ✓ |
| A766:14690 件 137 行（`89080fcc`）／`f8b738f3` 68 行 | `--numstat`＋`wc -l` | HEAD | ✓✓ |
| A766:14689 **`internal/statemachine/table.go:97`** | `sed -n 97p` | HEAD | `D43: 14, From: StateListening, Event: EvAudioDeviceLost, To: StateError` ✓ |
| A767:14703 A746 名册尺＝9 枚（97/170/228/246/256/258/259/278/284） | `git grep -lE 'AskOnTaskRoot\|\bReplay\b' HEAD -- .scratch/wisp/issues` | HEAD | 9 枚、成员逐字 ✓ |
| A768:14718 数行尺 new=12（末 12.）／copy=9（末 21.）／腿件标题"誊抄（10 组）" | 腿件那条 awk＋`grep -oE '誊抄（[0-9]+ 组）'` | HEAD | ✓✓✓（＝"标题 10 ≠ 行 9"确证） |
| A769:14729 腿的字样读数（465/22/21/14/17/4/11） | 尺形＝**行**：`git grep -c <词> c2b422a4` 求和 | **c2b422a4**（腿自己的锚） | 逐枚精确复现 ✓（今天第一例零漂移成立） |
| A769:14731 四枚锚 | `sed -n` | HEAD | `config_reload.go:191`／`config_handlers.go:431`＋`:435`／`panel_config_store.go:228`＋`:275`（`res.Tier = panel.EffectiveRestart` 两枚硬填）✓ |
| A769:14733 登记表行 `config_readers_255.go:103`（`hotClaimConsumed`） | `sed -n` | HEAD | ✓ |
| A770:14741 断言真身 `loop_golden_test.go:69`、`:70` 是 `t.Errorf` | `sed -n '66,72p'` | **d106fc25**（该节自称锚） | ✓✓（现 HEAD 已被 286 搬成正向形状，`:69` 成注释＝后续合法位移，非错） |
| A770:14742 `bd124b2a` subject 逐字含"`loop_approval_test.go:213` 唯一等值断言具名重判"／author `2026-10-08 19:56:27 +0800`／父＝`52501e22`／`ae23da67 loop.go:599`＝`CorrelationID: taskID` | `git log -1`＋`rev-parse ^`＋`sed -n 599p` | 各该 commit | 四枚 ✓ |
| A770:14743 名册并排八枚（distinct:119／loop_approval:219／283:222／rows:72／bridge:269／models:167／golden:341→现 :352） | `git grep -nE 'CorrelationID (!=\|==) ' 0b0ab517 -- '*.go'` | 0b0ab517 | 六枚 ✓；`净动作面＝1 枚`在 d106fc25 上 ✓ |
| A770:14745 282/r1 三件 78/134/91／两枚被种行回原形／282-done 86 字符／287 名 87 字符 | `wc -l`＋`sed -n`＋`awk length` | HEAD | ✓✓✓✓ |
| A770:14745 全池 `-done`＝**106** | `git ls-tree -r --name-only 0b0ab517 .scratch/wisp/issues \| grep -c -- '-done'` | **0b0ab517**（写该节的 commit） | **106** ✓（现 HEAD＝107＝`a09c3ca4` 结案票 286 之后，不是错） |
| A770:14741 `d106fc25`＝2 件 146 行（67＋79） | `--numstat`＋`wc -l` | 该 commit | ✓ |
| A771:14748 件 116 行＝`1e1017a5` 115＋`22812cda` 补 1 | `--numstat`×2 | 两 commit | ✓（与"各只含这一枚"同形） |
| A771:14748 回执两枚读者锚＋配对尺 | `sed -n`＋`git grep -nE 'NewRequestID\(' HEAD … ':!*_test.go'` | HEAD | `panel_inbound.go:164`✓／`panel_host_windows.go:803`✓／`bridge.go:80` 只有定义＝**非 test 调用者 0**✓ |
| A771:14748 两把尺并排（腿 27 行/4 文件 vs 编排者 38 行/4 文件） | `git grep -nE 'RequestID\|requestId' HEAD -- internal/panel ':!*_test.go'` | HEAD | 我＝**38 行／4 文件**（复现编排者那把，文件数与腿一致；行数不同＝尺形不同，两把并记）✓ |
| A771:14749 `internal/ball/table.go` 在 HEAD 不存在 | `git cat-file -e HEAD:internal/ball/table.go`＋`ls-tree -r internal/ball \| grep -i table` | HEAD | `does not exist` ✓／带 table 的只有 `tokens_table_test.go` ✓ |
| A771:14749 `SetAudioLevel` 定义 `:42`＋注释 `:40-41`、`EvAudioDeviceLost` 定义 `events.go:30` | `sed -n` | HEAD | ✓✓ |
| A771:14750 三枚默认档 `schema.go:245/197/277`＋`manager.go:380`、`SpawnCapture` `audio.go:180` | `sed -n` | HEAD | ✓（`default:"true"`／`"false"`／`MicMutedDefault … default:"true"`） |
| A771:14750 `SetAudioLevel` 非 test 调用 2 枚全在 `cmd/balldebug`（`:418`／`:477`）；`MicMutedDefault` 产码零读取处；`internal/audio` 非 test importer＝0 | `git grep -nE '…' HEAD -- ':(exclude).scratch' '*.go' ':!*_test.go'`（importer 尺用带引号的模块串） | HEAD | 2 枚 ✓／只剩定义＋两行注释 ✓／0 ✓（松散尺给 1＝`cmd/wisp/config_readers_255.go` 的一句注释） |

### ❌ 对不上（8 枚；判语归编排者，本腿只给料＋最小更正句）

1. **A761:14638「就地订正 21 处名义时刻」⇒ 现量 23 处，且枚举相加＝27。** 尺＝`git show 69c2bfd9 \| grep -E '^-[^-]' \| grep -E '[0-9]{1,2}:[0-9]x'`（＋同尺在 `29ee2e0f` 逐票）。读数：台账节头 **6**（A755–A760）＋HANDOVER `4.0ag` **1**＋票面 **16**（242 五／282 三／279 二／285 五／281 一＝枚举逐票全对）。**缺的那三枚＝句里"三处『编队现状（…现量）』行"现量零订正**：HEAD 上 A756 那行仍写 `10:3x` 而它的节头已被改成 `10:2x`。最小更正句：「订正 23 处＝台账 6 枚节头＋HANDOVER 1 枚节头＋票面 16 枚；句里那『三处编队现状行』未订正，A756 现仍 10:3x，与节头 10:2x 自相矛盾。」
2. **A770:14747「288-…-after-280.md（74 字符）」⇒ 现量 72。** 尺＝`git ls-tree -r --name-only HEAD .scratch/wisp/issues \| awk '{n=length($0); sub(/.*\//,"",b); if (b ~ /^288-/) print length(b)}'`。同句里 282-done＝86 ✓、287＝87 ✓、284＝66 ✓，只 288 偏 2。最小更正句：「288 名长 72 字符。」
3. **A760:14626「⛔ 不许放宽 `internal/tools/loop_approval_test.go:213-214` 那枚唯一等值断言」⇒ HEAD（乃至 `bd124b2a` 落地当日）`:213-214` 是注释，断言在 `:219-222`。** 尺＝`git show <ref>:internal/tools/loop_approval_test.go \| awk 211-222`。三锚读数：父锚 `52501e22` 的 `:213` ＝**断言**（故 `bd124b2a` subject 当时不算错）；`bd124b2a` 与 HEAD 的 `:213-214`＝**注释**。＝A762 那条"裸行号腐烂"在**禁令面**的第二次发作：后续程照 `:213-214` 下刀会砍在注释上。最小更正句：「那枚唯一等值断言现量＝`loop_approval_test.go:219-222`（内容短语 `r.CorrelationID == "" || r.CorrelationID == res.TaskID`）；`:213-214` 只是注释，不许再当种刀禁线。」
4. **A768:14720「顺手重跑的另一格在件 `…/30-impact-pins-producers.md`，`grep -cE '^\| D[0-9]+ \|'` 亦回 29」⇒ 该枚件现量 0。** 尺＝逐字两把（`grep -cE '^\| D[0-9]+ '`＝0；同目录 `git grep -cE '^\| D[0-9]+ \|'`＝**20-types-and-readers.md:29**）。＝读数 29 对、**归件错**（那张 D 表住在同目录另一枚件里）。本腿同尺确证 A768 其余三数全对（new 12／copy 9／标题 10）。最小更正句：「D 表在 `242/corrcensus1/20-types-and-readers.md`（29 行），不是 `30-impact-pins-producers.md`。」
5. **A757:14578「`MU-D2c`（换 `int` 并填 `it.grants.live()`，`queue.go:578`）」⇒ `:578` 现量＝`}`，在 A757 自己的锚 `37f2a0bb` 上同样＝`}`，且 `queue.go` 里 `Permitted` 字面零命中。** 尺＝`git show 37f2a0bb:internal/agent/approval/queue.go \| awk 576-579`＋`grep -n Permitted`。⇒ 那一发的种刀位**未被任何一锚证成**（其余七枚锚在同批全部逐字命中）。最小更正句：「`MU-D2c` 的落点锚改为**内容短语**（`PanelItem{…}` 构造行），`queue.go:578` 只是该字面量收尾的 `}`。」
6. **A754:14538 的三个行号（`:248 still had`／`:249 This call is the caller`／`:260 startResidentTaskSource`）行号全对，但真身住在 `cmd/wisp/resident_windows.go`，而该条整篇在裁 `resident_approval_windows.go` 且这三处未具名文件。** 尺＝`git grep -n "This call is the caller" HEAD -- '*.go'`＝`cmd/wisp/resident_windows.go:249`（`resident_approval_windows.go` 零命中；该文件 `:249` 现是空行）。＝"路径没现量"那一族的**新形：行号对、文件缺**。最小更正句：「补文件名：`cmd/wisp/resident_windows.go:248/:249/:260`（⛔ 不是 `resident_approval_windows.go`）。」
7. **A751:14511「已勾 **14→24**」＝仍在原位的错数**（真值 26→36；`14` 是未勾栏）。尺＝八枚乙类票逐票 `grep -cE '^[[:space:]]*- \[x\]'` 求和 @c39e2853↔542d7661＝**26→36**。A760 F1 已判、A761 也定过式，本轮**再确证它仍逐字在 14511**。＝料给编排者决定是否就地打旧。
8. **A754:14537「票 284 名长 62 字符」＝仍在原位的错数**（现量 66，A760 F3 已认）。尺见 ❌2 那把。附一条**尺形警告**（本腿亲踩）：数 `-done` 名册时 `grep -c -- '-done'` 会比 `grep -cE '\-done\.md$'` **多 1**（票 281 的 slug 自带 `class-b-done-tickets`）⇒ 名册尺必须锚死后缀 `\.md$`。

### ⛔ 复跑不了（具名原因，一律不填空）

- **红／PASS／rc 类（禁跑 go）**：A752 基线 rc=0/2 PASS/0.129s、两发反例 rc=1 的红句；A757 起手锚 90 PASS/0 FAIL；A760 整包 92 PASS/0/1 与四件套红绿枚数；A765 三枚新用例全 PASS／`ok 13.6s`；A770 `internal/tools 209/0/0`、`internal/agent 84/1/0`、"起手即在的那枚红"；A762 那句〔D2 删比对后 go test 仍绿〕。
- **gofmt／gofumpt 发声列**：A756 表里 `gofmt -d`／`gofumpt -d` 两列（0/0…14/14、14/101）＝要 `gofumpt.exe`，属 Go 工具面；本腿只复了 CR 两列。
- **门禁脚本读数**：A755 的 `tracked paths=8446／over-budget 57／covered 57／not in roster 0`、A765/A770 的 `VERDICT GREEN`＝禁跑脚本。侧量：`git ls-tree -r --name-only` 在 `23924be3^`＝**8447**、在 `23924be3`＝**8449** ⇒ 8446 落在两次之间，**无法定点复跑**，不判对错。
- **需真窗／需机主**：07 的 `[H1]` 人眼签收、`A696` 第三发对照真窗、`build/wisp.exe` 重建与否。
- **归属冻结（⛔ 未开那棵工作树）**：A769 的 `frontend/src/lib/panel.ts:211`／`:140-142`＝只按台账原文登记；`internal/panel/composer.go:58-91` 无回执字段一格我只验了函数存在性、未逐行读。
- **外部库源码**：A771 的 `webview.go:147-157`（在 `GOMODCACHE`，取路径要 `go env`）＝没碰。
- **时点已过期**：A769 那句"我在 HEAD 上跑＝475/25/22/16/19/5/12"（HEAD 已推进 7+ 枚 commit，我侧＝479/27/23/18/20/6/13 ⇒ 只能证"会漂"，判不了它当时对不对）；A760 F2 的 13 枚/14 行（含**未跟踪**件的全树 `grep -rl`，本腿不读未跟踪面）；A767 三数 129/114/8/7＝腿件自数，按 A768 定式只当宣称。
- **未下钻（射程内但没做完，具名）**：A763 票 33 的 `:31` vs `:15` 两枚锚；A764 承载表末行 `:27`；A756 引用的腿件 `10-census.md:85/:86`；A758 的 `.gitattributes:4`；`2bdf385d`／`79957c67`／`4d9ac0d8`／`091fe90a`／`332e2828`／`6041d7f9` 六枚 sha 只验了存在性的一半。

## ③ 我复跑不了的、以及为什么

见上一档逐条。**总体判语留给编排者**；本腿只报三句形状观察：
1. 今天 A751–A771 的**承重句绝大部分对得上**（58/84），错处集中在**三族**：枚数（❌1、❌2）、裸行号腐烂（❌3、❌5）、**文件归属缺名/错记**（❌4、❌6）——后一族是今天新形：行号精确、文件没写，下一棒照样砍空。
2. 三枚**已认账的旧错仍逐字在原位**（❌7、❌8 与 A754 的 62）＝"另写一节更正"没让旧句消失，与 A760 F1 同形；是否就地打旧归你。
3. 本腿自报两处尺形学费：`-done` 名册尺必须锚 `\.md$`；字样族读数的尺形是**行数**（`git grep -c` 求和）不是出现次数（`git grep -o`）——后者会把 A769 那串 465/22/21 复现成 476/23/21 而误判腿数错。
