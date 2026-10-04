# 212-v3 非实现者验收件（攻 `212-r3` 的三笔改动：4ea6411e 骨架 → d03d166f 落地 → 0cb71a58 满稿）

- 本程：验收腿 **212-v3**（裁决者≠实现者；实现＝`212-r3`，ⓐ修法裁定＝编排者）。
- 起手锚：`dev` HEAD＝`57a33804`（本程第一条命令现跑，时刻 `2026-10-04 10:45:03 +08`）。
- 本程只读产码、只在 `tools/d22scan/**` 里做突变；AC 框一枚不碰；不放宽任何既有断言。
- 每一节的判语只认本程自己跑的尺；`212-r3` 自述里的每一句都按待验断言处理。

## ① 进场、基线与台件纪律自证

打算答什么：起手 HEAD/脏度、被验三笔的逐枚 commit 与写面名册、四枚受检文件的 md5 起手值，
以及"工作树 CRLF ≠ 进树 blob"这一枚尺坑的具名处理（还原自证用哪把尺、为什么）。

## ② 格1：AC#4 那把正控到底还有没有牙

打算答什么：自己在盘外真树种一发"注释引用盘上不存在的证据件"⇒ 必响，删掉 ⇒ 必不响，两发红/绿句逐字入账；
再读 `main.go` 喂给 ban #9 的节点集合，裁"只读 `f.Comments`、字符串字面量永远扫不到"这句话成不成立，
**以及它成不成立会不会把 AC#4 的正控做瞎**（要分开答：对仪器射程的影响 vs 对正控的影响）。

## ③ 格2：ⓐ 豁免的真实射程（本票最要紧的一格）

打算答什么：种一发全拼、不含任何省略号、盘上不存在的 ⇒ 必红；种一发带 `...` 的全拼不存在 ⇒ 静默（设计）；
两向都有逐字读数。再拿真仓现量回答"有没有一种真实写法把谎报伪装成缩写从而溜过①"：
逐枚列出真仓注释里**当前被豁免吃掉**的 token，具名判其中有没有"完整幻影＋尾标记"活体；
找不到也要列我试了哪几种形（`*` markdown 强调、尾部 `...`、CJK 粘连、换行、空格）。

## ④ 格3：新增那枚前缀同步钉是不是同义反复

打算答什么：确认 `repoPathRe`/`shorthandRegionRe` 是不是同一份常量的两个别名（读源码＋读钉的实现），
然后用**我自己的刀**做一发 m2 同形突变（换不同的前缀、换不同的那一枚表），
逐字入账指名 FAIL 的行号，并复跑 `-self-test` 证它仍全绿（＝钉只在 `go test` 那一侧有牙，这算不算缺口要判）。

## ⑤ 格4：AC#5「八枚 ban 读数逐名不变」与两枚 +1 的归因

打算答什么：自己跑 `sh scripts/d22scan.sh` 取终态八枚 scope 行逐字；
用 `git log --name-status`／`git ls-tree` 在 212-r3 三枚 commit 的前后锚点上复算 ban #8 分母口径（含 `_test.go`），
给出"哪一枚新增 `_test.go`、随哪一枚 commit 进树"的结论，并明确判它是不是 `212-r3` 的伤。

## ⑥ 格5：(a) 那处"诚实化"引用的每一枚 `file:line` 真不真

打算答什么：把 `pending_read.go:40-54` 新注释里出现的每一处 `file:line`（`internal/tools/gate.go:12-13`／
`cmd/wisp/panel_pump.go:44-46`／`internal/panel/approval.go` 那一段）逐枚打开对行号与内容，
判"两枚 vs 四枚"的归因是否已经指对；引错一枚即本条判失败。

## ⑦ 格6：恒真判据自查——三枚新样本是不是装饰

打算答什么：自己跑一发摘豁免的突变，逐枚看三枚新样本里有没有**摘掉豁免仍不红**的（那就是装饰），
并判每枚样本响的**那一枚 token** 是否正是它声称钉住的那一形（防止"响在别的 token 上"的假钉）。

## ⑧ 门禁读数（终态）

本程自己跑的四把尺，全部重定向到 `logs/`，rc 用 `$?` 直取（⛔ 没有用 `cmd | head; echo $?` 那种取到 `head` 退码的形状）。

| 尺 | 命令（逐字） | rc | 终态逐字（关键行） | 档 |
|---|---|---|---|---|
| 门禁本体 | `sh scripts/d22scan.sh` | **0** | `runtests.sh: OK - packages=[./...] top-level: PASS=35 FAIL=0 SKIP=0, === RUN=77, '[no tests to run]'=0` ＋ `d22scan: examined 266 production Go files under internal/ and cmd/ of D:/work/workspace/projects plans/Wisp` ＋ 八枚 scope 行 `bans #1-5 internal/=228, bans #1-5 cmd/=38, ban #6 frontend/=85, ban #7 internal/tools/=23, ban #8 design/=39, ban #8 frontend/=85, ban #8 internal/=501, ban #8 cmd/=98` ＋ `d22scan: clean - no D22 ban violations` | `logs/gate-v3-final.txt`／八行摘出 `logs/gate-v3-scoperows.txt` |
| 仪器自检 | `cd tools/d22scan && go run . -self-test` | **0** | 终行 `d22scan -self-test: clean - all 40 direction checks passed (20 expect-ring, 20 expect-silent)`；名册行 `... = 9 numbered ban(s) [1 bare-goroutine, 2 pathresolver-bypass, 3 plaintext-key, 4 wallclock-timeout, 5 mirror-hash, 6 panel-approval, 7 internal-artifact-tool, 8 emoji, 9 phantom-citation] + 1 finding type(s); 40 cases, 10 tag(s) covered, both directions required per tag` | `logs/m1-selftest-restored.txt` |
| 模块内 vet | `cd tools/d22scan && go vet ./...` | **0** | 零输出 | `logs/vet-tools.txt` |
| 新钉在门禁里真跑过 | `grep -n "PASS: TestBan9PrefixTablesAreTheSameList" logs/gate-v3-final.txt` | 命中 | `--- PASS: TestBan9PrefixTablesAreTheSameList (0.00s)`（第 224 行）＝这枚钉**在 `sh scripts/d22scan.sh` 的第一步里执行**，不在 `-self-test` 里 | `logs/gate-v3-final.txt:224` |

**八枚 scope 行与编排者 10:40 那发、与 `212-r3` §7 终值逐字相同**（我 this 程 10:53 自跑，非引用）。

**突变还原尺（起手 vs 终态）**：
- `md5sum -c logs/md5-start.txt` → 四枚全 `OK`：`tools/d22scan/main.go 442d46cadef8c641c54f7ac00984a470`·`selftestsamples.go f9e10163199a89c2614daf1ab1ad7d52`·`scan_test.go 1bcf66c1bc7da3e7ff0601b66da13504`·`internal/agent/approval/pending_read.go acef5b1739e9446cd7118bbffa6a0a3b`（档 `logs/md5-final-check.txt`）。
- `git diff -- tools/d22scan internal/agent/approval/pending_read.go | md5sum` = `d41d8cd98f00b204e9800998ecf8427e` ＝ **空 diff 的 md5**（档 `logs/diff-md5-end.txt`）⇒ 两发突变零残留。
- ⚠ **尺坑具名**（本程起手就撞到并改刀）：`pending_read.go` 的**工作树 md5 ≠ `git show HEAD:<path> | md5sum`**（工作树 `acef5b1739e9446cd7118bbffa6a0a3b` vs HEAD blob `219a0465f27c469217227ef7a48b907d`，档 `logs/md5-head-vs-work.txt`），`file` 报该文件 `with CRLF line terminators` 而 `.gitattributes:4` 写 `*.go text eol=lf` ⇒ 差异**只来自行尾**，`git status --porcelain` 对该文件为空。所以本程对**这一枚**文件的还原尺用 md5（对起手工作树）**＋ `git diff` 为空**两把，不用 md5(HEAD) 比 md5(工作树)。`tools/d22scan/*.go` 三枚是 LF，两把尺同值。
- ⚠ 编队对账：终态 `git status --porcelain -- tools internal cmd scripts docs .github` 有 **2 枚**（` M cmd/wisp/resident_approval_windows.go`·` M internal/agent/approval/approval.go`，档 `logs/porcelain-end.txt`），逐枚核对＝并发腿 **260-r3** 声明的写面，**本程一枚未碰**；我 10:53 那发门禁正是扫着这两枚在飞的修改跑的，它们是对**已存在文件**的修改（不增减文件枚数），故八枚分母不受影响——这一点与 §⑤ 的归因分开讲，不混。
- 卫生（非本票六格，具名不修）：`gofmt -l tools/d22scan` 报 `selftest.go`；本程把 **HEAD 版**该文件落到临时处再 `gofmt -l` **同样被报**（`logs/gofmt-head-selftest.txt`）＝起手既有，`212-r3` 未碰该文件，判定与腿 §7 一致。

## ⑨ 判不动的地方与六格三档判语

### ⑨.1 三档判语（六格一次裁，每格指回本程自己的尺）

| 格 | 被验对象 | 判语 | 一句话凭据（本程现跑） |
|---|---|---|---|
| 1 | **AC#4 正控还有没有牙** ＋ "ban #9 只读 `f.Comments`" 这句话成不成立、会不会把正控做瞎 | **成立**（两半都成立，且**不构成互相抵消**） | 种一发注释引用盘上不存在的件 ⇒ `rc=1` ＋逐字 `comment cites repo path "docs/evidence/s1/212-v3-plant-phantom.md" ...`；删掉载体 ⇒ `rc=0` 零 phantom（§② P1/P10）。`f.Comments` 那句**为真**（§② S1 字符串字面量 `rc=0`），但它**做不瞎 AC#4**——票面 AC#4 种的正是"注释"，块注释也响（S2）；被它做瞎的是**另一格**：`tools/**` 扩射程那一格（S4 已证 out of range，且 S1 证明扩了射程也看不见字符串里的旧幻影） |
| 2 | **ⓐ 豁免的真实射程**：全拼必红／`...` 全拼静默；有没有"谎报伪装成缩写"的真实写法 | **带条件成立**（条件＝豁免面比"省略号中段"宽，本程**在盘外真树造出两形真实写法并复现静默**，而**真仓当前活体 0 枚**） | 双向都有：P1 全拼 `rc=1`；P2 `docs/evidence/s1/212-v3-...-phantom.md` `rc=0`。溜过①的伪装形：P9 尾部 `...`（"等等"式）、P3 尾部 `…`、**P4 markdown 加粗 `**docs/evidence/s1/212-v3-bold-phantom.md**`**（＝把谎报写成"强调"）、P5 CJK 粘连、P12 `*.md` 通指——全部 `rc=0`；边界同时钉住：P6 粘连无标记 `rc=1`、P8 标记换行 `rc=1`、P11 标记前有空格 `rc=1`。真仓现量（我自己写的 ban#9 同语义普查，`logs/census-real-repo.txt`）：RING=**0**／EXEMPT=**1**（`cmd/wisp/panel_host_windows.go:29` 的 `docs/evidence/s1`，逐字是真缩写不是谎报）／SYMREF=3 ⇒ 三形的**活体仍是 0 枚**，编排者"不收窄"那道裁定在本程读数下**不被推翻** |
| 3 | **前缀同步钉是不是同义反复** | **成立（有牙，不是同义反复）** | 两枚正则是**两份独立字面量**（`main.go:891`·`:905`），钉用 `regexPrefixAlternation` 各读各的 `.String()`；**我自己的刀**＝给 `shorthandRegionRe` 加 `plugins\|`（腿加的是另一枚表），⇒ `scan_test.go:2441` 逐字 FAIL 且**指向第二臂**（`only in shorthandRegionRe=[plugins] ... widens the exemption without an owner's approval`），同棵树 `go run . -self-test` 仍 `all 40 direction checks passed`；还原 md5 全等 ⇒ 真实漂移看得见。**条件具名**：这枚钉只在 `go test` 那一侧有牙——我已复认它**确实在门禁第一步与 CI 的 `D22 scanner positive control` 步里执行**（`scripts/d22scan.sh` 调 `tools/d22scan/runtests.sh -C tools/d22scan ./...`；`ci.yml:74-81` 同一条；§⑧ 那发 `--- PASS` 在档），所以"只有本地跑测试才看得见"这一支**不成立** |
| 4 | **AC#5 八枚 ban 读数逐名不变**＋两枚 +1 归因 | **成立**（八枚逐名与编排者 10:40、与腿 §7 终值逐字相同；两枚 +1 **不是 `212-r3` 的伤**，且**腿对 `internal/` 那一枚的归因文件指错了，具名更正**） | `ban #8 internal/` 500→501 ＝ `internal/ball/hotkey_cancel_borrow_expect_260r2_test.go`，随 **`e3e19e8e`（260-r2，10:36）** 新增（尺＝`git log --name-status --diff-filter=A b1e59d63..HEAD -- internal cmd`＋`git ls-tree -r --name-only` 两锚点作差，档 `logs/added-go-between-anchors.txt`）；⚠ 腿 §7 猜的是 `internal/ball/hotkey_live_test.go`（那枚是 **M 不是 A**，不改变分母）⇒ 结论方向对、**点名文件错一枚**。`ban #8 cmd/` 97→98 ＝ `cmd/wisp/resident_approval_risk_256_windows_test.go`，随 **`2a10134e`（256-r1，10:20）** 新增，与腿一致。两枚都是 `_test.go`＝只有 ban #8 含、bans #1-5 不含，所以 `228/38` 未动。另：`212-r3` 在 `internal\|cmd` 内**只 M 了 `pending_read.go` 一枚、零新增文件**（同一把尺），结构上就不可能移动分母 |
| 5 | **(a) 注释"诚实化"引用的每一枚 `file:line` 真不真** | **成立**（四枚逐枚打开对全中，一处未把本票的缺陷再犯一遍） | `internal/tools/gate.go:12-13` 逐字只点名 `RulesHit and Reason VERBATIM`＝**两枚**✓；`cmd/wisp/panel_pump.go:44-46` 逐字 `Level/RulesHit/Reason/SessionOverrideBlocked`＋`a copy, not a judgement`＝**四枚**✓；`internal/panel/approval.go:72-88`＝`CardViewFromDecision` 函数体（`:72` 起、`:88` 收到结构字面量右括号），四枚分别在 `:81/:82/:83/:85` 上视图✓；`internal/agent/approval/ticket146_liveapprovals_backing_test.go`（同段旧句）`find` 命中存在✓。新句另引 `tools/d22scan`／`risk.PathResolver`／`filepath.Clean/Abs` 均为符号不是路径✓。**没有一枚是幻影** |
| 6 | **三枚新样本是不是恒真／装饰** | **成立**（三枚都有牙，且各响在自己声称的那一形） | 我自己的 m1（`short := map[int]bool{}`）⇒ `4 direction(s) failed, 36/40 passed - the gate does not see what it claims`，三枚新样本**逐枚翻红且 token 各不重叠**：`cites-tail-ellipsis.go:4 ... "docs/evidence/s1/probe-tail-ellipsis.md"`／`cites-tail-star.go:4 ... "docs/evidence/s1/probe-tail-star.md"`／`cites-cjk-attached.go:5 ... "docs/evidence/s1/probe-swallowed.md见docs/evidence/s1"`，第四枚是 r2 那枚 `cites-shorthand.go`（三 token）；还原后 `all 40 ... (20 expect-ring, 20 expect-silent)` rc=0 ⇒ **没有一枚摘掉豁免仍不红** |

**其余两笔（不在六格里，本程顺手量到，具名入账）**：
- **(c) 没有变成新幻影**：`ls -l docs/evidence/s1/212-comments-phantom-citation-v2.md`＝**33,578 字节、mtime Oct 4 09:57**，复认编排者③；`grep -rn "212-citation-ruler" --include=*.go cmd internal tools` 只命中 `selftestsamples.go:381` 的 **`note:` 串**（那一句自己写明"该件从未存在"＝历史说明，不是路标），**产码注释零命中**⇒ 复认编排者④。
- **行为零变化这一条我另有一把更硬的尺**：`git show d03d166f -- tools/d22scan/main.go` 与 `-- internal/agent/approval/pending_read.go` 的**全部改动行都是注释行**（`grep -E '^[+-]' | grep -v '^[+-][[:space:]]*//'` 两枚文件各 **0 行**，档 `logs/diff-main-noncomment.txt`／`logs/diff-pendingread-noncomment.txt`）；`scan_test.go` 那一枚 hunk **86 增 0 删**（未放宽、未删除任何既有断言）；`selftestsamples.go` 删除行只有 3 枚，逐枚都是 (c) 那一句的 `src`／`summary`／`note`。⇒ 腿自述"行为零变化"**成立**，且比它自己给的 `--numstat` 版更强。

### ⑨.2 判不动／没跑到的地方（具名＋归口，⛔ 不写"应该没问题"）

| 枚 | 判不动的东西 | 为什么本程判不了／没跑到 | 归口 |
|---|---|---|---|
| 1 | 把 `tools/**` 纳入 ban #9 射程到底值不值 | 本程量到了**代价的形状**（S4 证 tools/ 现在完全不扫；S1 证扩射程也扫不到字符串字面量里那枚旧幻影；粗分母另算），但**扩射程＝人工批准面**，不是验收腿能裁的 | 编排者：随票 264 同批摆给机主（`A592` §4 已写死"不单摆一次"） |
| 2 | 把 `_test.go` 纳入射程 | 同上；唯一活体缩写形 `cmd/wisp/slo_report_144_windows_test.go:851` 本程**未复跑**（不在三格内，且它与谎报形不同＝它是缩写不是完整名） | 与 1 同批 |
| 3 | 豁免面会不会被"新前缀"绕过（除枚 4 之外） | 那需要真有人往两枚表里**只加一枚前缀**；我已用 m2 复现"只加一枚表"的形状并证钉响，**"两枚都加但拼法不同"（如 `.plugins` vs `plugins`）本程没做**——钉是**字符串集合**比较，那种拼法差异它会红（同一把尺），但红句会不会指对"哪一枚该改"我没验 | 若真发生漂移由那次改动自己的验收腿判 |
| 4 | 台账两条散文过期（`docs/reports/pending-and-issues.md:11408`·`:7832` 记 `37 direction checks`） | 台账**只追加不删**，且不在本程写面；本程终值 40＝20/20（§⑧） | 编排者下次记台账时带口径与锚点（腿 §6 同判） |
| 5 | NTFS 大小写洞（P10 静默／P11 响）与 `symRefRe` 三枚排除的**语义**是否都该排除 | 不在票 212 的三格内，本程**只复算了枚数**（SYMREF=3）没逐枚判该不该排除 | 已知残余，`A592` §4 ③ 已具名归口 |
| 6 | `internal/probe/*` 夹具样本的**真实产品后果**（夹具不是真树） | 三枚 silent 样本钉的是**豁免本身**，真树里那三形的活体是 0（§② census）；"豁免会不会某天盖住真谎报"只能等活体出现，本程**造得出形、造不出需求** | 若将来抓到活体 ⇒ 新票，⛔ 不许顺手改 regex |
| 7 | 名册行/自述里"扩面＝批准面"这类**治理句**是否与 `SPEC-12 §4.1` 逐字一致 | 本程核的是**射程与行为**，没有把每一句治理措辞回抄到 spec 逐字比 | 编排者；本程不判"措辞合规" |
| 8 | 票 212 AC 框与 `-done` 改名 | AC 框与改名归编排者（⛔ 本程一枚没碰，`git log -1 -- .scratch/wisp/issues/212-*.md` 见 §①） | 编排者按 §⑨.3 办 |

### ⑨.3 票 212 现在够不够格改名 `-done`

**够格**（本程判语）。判据逐条：
- ⛔ 编排者写死的那一枚否决条件＝"只要 AC#4 那格被你自己重跑推翻，就不许 `-done`"——**本程重跑没有推翻它**：种⇒响（`rc=1`＋逐字 token）、删⇒不响（`rc=0`），两向都在**盘外真树**用**同一枚产码 binary**复现（`logs/plant-p1-*.txt`／`logs/plant-p10-*.txt`）。
- 三笔欠账逐笔有下落：**(a)** 出处指对（§格5 四枚 `file:line` 全中）、**(b)** 三形各一枚 silent 样本且**摘掉豁免逐枚翻红**（§格6）、名册行与 `main.go` 自述写出真实射程（§格2/格1，本程用 10 发种删复现，未见"读起来比实际窄"的句子残留）、**(c)** 改指真实件且**没变成新幻影**（33,578 字节现读）。
- 加的那枚钉**有牙**（§格3 我自己的刀），且**没有放宽／删除任何既有断言**（`scan_test.go` 86 增 0 删）。
- 本程唯一**必须更正腿**的一处不是伤，而是归因点错文件：`ban #8 internal/` 那枚 +1 的真身是 `hotkey_cancel_borrow_expect_260r2_test.go`（随 `e3e19e8e`），不是腿猜的 `hotkey_live_test.go`——**结论不变（不是 212-r3 的伤）**，请编排者在翻勾节带上这一处更正。
- 建议随 `-done` 一并带的两条**已知残余**（⛔ 不是退回理由）：`tools/**`／`_test.go` 是否在 ban #9 射程＝待人批准（§⑨.2 枚 1/2）；台账 `37 direction checks` 两条散文过期（枚 4）。

