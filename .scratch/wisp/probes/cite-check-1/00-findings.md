# cite-check-1 · 凭据引用对账（只读腿，⛔ 零产码改动、⛔ 未动票面一字）

时刻：2026-10-09（读数全部取自 `git show HEAD:<path>`，对象层＝LF；HEAD＝`fec4f5e7f70819d9623c408e8ad4858d251009f9`）
尺：`git show HEAD:<文件> | awk -v n=<行号> -v q=<短语或正则> 'NR==n{ATLINE} $0~q{记录真行号}'`
⇒ **行号与短语在同一次取数里拿到**（⛔ 不分两次，两次之间会有腿在飞）。
判三形：**对上**／**漂了**（短语还在文件里，行号不同）／**内容不符**（那一行根本不是那句，或那个符号已不在文件里）。
显示截断到 180-190 字节只影响我眼看，**判定按整行内容**（`$0~q` / `index()` 都在整行上算）。

## 0. 总数

| 类 | 处数 |
|---|---|
| 核过的 `文件:行号` 引用出现数（跨票重复按出现计） | **173** |
| 对上 | **143** |
| 漂了（票面未自纠） | **12** |
| 内容不符（票面未自纠，最严重） | **2** |
| 票面已就地打旧、更正后的行号复查**全部落准** | **13 组** |
| 判不动 | **3** |

按票的分布：242＝对上 33／漂 12／不符 2／自纠 1／判不动 2 · 277＝14＋2自纠 · 278＝18＋4自纠 · 279＝21＋1自纠 · 280＝12＋1判不动 · 281＝7 · 282＝16＋1自纠 · 283＝7 · 285＝15＋4自纠。
另：票面引的 **23 枚 commit sha 全部存在**（`git cat-file -e <sha>^{commit}`），含 `7789f953`／`4d9ac0d8`／`ed8a7358`／`2a2dad22`／`c9348a99`／`039ec93c`／`2bdf385d`／`351e5a5e`／`6c969dc8`／`37f2a0b`／`77150dcc`／`bd124b2a`／`c39e2853`／`623a4d81`／`a1b19873`／`aa6c1881`／`613606c0`／`133b1bfa`／`0abe217c`／`3fe377e3`→`4df747d8` 四枚。
裁决表/探针件的**行段引用全对**：`docs/evidence/s1/242-grant-binding-v2.md`（200 行）`:49-110`／`:91-100`（"判据＝造两枚同时活着"实测在 `:95`）／`:102-110`（"判不动"标题在 `:102`）／`:114-168`（`MU-D2a` 在 `:123`）／`:170-200`（"互不敏感"在 `:175`/`:181`）；`probes/280/r1b/10-census.md`（91 行）`:17-53`／`:55-73`／`:75-80`／`:85`／`:86` 逐段命中；`probes/281/r1/10-reconcile.md`（75 行）`:16-27`／`:32-41`／`:43-46`／`:68` 逐段命中；`probes/279/c1/00-ac3-and-anchor-check.md`＝**135 行**（票面那句"135 行"逐字对上）；`probes/comment-fix-prep-1/01-ready-to-apply.md` P39 段首 `:327`、原文 `:332-337`（判别位在 `:334`）、替换 `:341-349`（`Nothing reads this field yet:` 在 `:343`）；`probes/orch-277/00-own-rerun.txt`、`probes/orch-283r1/0{0,1,2,3}-*.txt`、`probes/done-class-b-1/13-summary.md`、`probes/282/v1/{10,20}-*.md`、`probes/stale-claim-1/` 全部在盘。

## 1. 漂了／不符（逐条，按严重度；票面原文引的 path:line → 该行现读 → 真行号）

**S 级＝下一枚腿会照着动手的落点／突变目标**

1. 票 242 `AC#1` 正控 D2（第 20 行）→ `internal/agent/approval/approval.go:292`「`grantStore.spend` 不再比对 binding 摘要」→ 现读 `}`（空的花括号行）→ **真身 `approval.go:558` `func (s *grantStore) spend(nonce, bind string) grantDenial {`**＝漂 −266。种在 `:292` 什么也改不动 ⇒ D2 会"全绿"被读成"这格没牙"。
2. 票 242 ★结构事实（第 66 行）→ `queue.go:376` 声称**逐字** `spent := it.grants.spend(nonce, it.bind)` → `:376` 现读注释 `// a recorded fact and start being an inferred one.` → **真身 `queue.go:414`，且字面已变：`denial := it.grants.spend(nonce, it.bind)`**。＝行号漂 −38 **＋内容不符**（变量名 `spent`→`denial`，票面那句"逐字"已不成立）。
3. 票 242 ★（第 66 行）→ `approval.go:303` 声称逐字 `equalSecret(stored, bind)` → `:303` 现读空行 → **真身 `approval.go:569`**＝漂 −266。
4. 票 242 `AC#1`（第 62 行）→ `approval.go:298-305`「消费即删」→ 该块现读属别处 → **真身 `approval.go:558-575`（`delete(s.values, v) // consumed whether or not the binding matched` 在 `:568`）**＝漂 −266。
5. 票 242 ★＋现量表 → `queue.go:162`「`it.bind` 全仓只有一处写（铸造那一次）」→ `:162` 现读 `q.seq++`（同一段的上一枚语句）→ **真身 `queue.go:177` `it.bind = bindDigest(corr, d.TaskID, d.Tool, d.LevelString(), q.seq, d.Args)`**＝漂 −15。
6. 票 242 ★ → `queue.go:336`「`issue` 存进 store 的也是同一个值」→ `:336` 现读 `}`（一个 `select` 的收尾）→ **真身 `queue.go:358` `it.grants.issue(nonce, it.bind)`**＝漂 −22。
7. 票 242 ★ → `queue.go:368-371`「state 检查（`statePending`）」→ 该段现读不在此 → **真身 `queue.go:390` `if it.state != statePending {`**（`allow` 函数头现在 `:365`）＝漂 −22。
8. 票 242 现量表（第 12 行）→ `approval.go:253`「`bindDigest(...)` 铸」→ `:253` 现读 `Channel Channel`（结构体字段）→ **真身：定义在 `approval.go:416`，调用（＝铸那一次）在 `queue.go:177`**＝漂 −163。同一行还引 `:292` 花，见第 1 条。
9. 票 242 现量表（第 13 行）→ `approval.go:223`「The panel-facing surface (PanelItem) has no grant」→ `:223` 现读是另一条中文注释（`// 「取消」能否决 L1 | AEC | Confirming 明示…`）→ **真身 `approval.go:351`**＝漂 −128。⚠ 票面把这三处叫"**三处逐字**"，其中两处的行号已不可用（第三处 `pending_read.go:9` 对上）。
10. 票 242 现量表（第 13 行）→ `cmd/wisp/approval_reply.go:28`「PanelItem has no grant field and PanelAPI has no Allow method」→ `:28` 现读是裸 `//` → **真身现在跨 `:30`-`:31` 两行**（`:30` 末 `...has no grant field`／`:31` 首 `and PanelAPI has no Allow method`）＝漂 −2 **且短语已不成单行**；按"整行逐字"去搜这一句会 0 命中。
11. 票 242 `AC#1` 第 62 行 ②（历史突变落点）→ `queue.go:235`「让允许侧那条路都不跑」→ `:235` 现读空行 → **允许侧那条路现在 `queue.go:365`（`func (q *Queue) allow(...)`）**＝漂 −130。（这条是 10-03 腿的读数，不是今天的落点，故排最后；但票面没标它过期。）
12. 票 242 禁区（第 41 行）→ `internal/tools/ticket90_test.go:431`「那枚反射钉」→ `:431` 现读 `func TestTicket90ModeCarriesNoAllowAuthority(t *testing.T) {`（函数头）→ **反射那枚断言在 `:449`**（`strings.Contains(low, "grant") || ...`）＝漂 −18，落的是同一个用例（低危，但票面说的是"那枚钉"）。
13. 票 242 ★旁证（第 66 行）→ `ticket242_binding_test.go:111`「那枚 `sameDigest` 零调用者而包仍 ok」→ `:111` 现读 `// prompt makes - and A and B sit pending at the same time (asserted via` → **该行不是那句，且 `sameDigest` 在该文件 HEAD 上 0 命中（符号已不存在）**＝**内容不符**。它是 10-03 读数，被 10-08/10-09 的改写清掉了；票面仍用现在式把它当盘上事实陈述。
14. 票 242 `AC#1` 第 62 行 ④ → `ticket242_binding_test.go:71`「测试用同一枚 `bindDigest` 自指复算 `:71`」→ `:71` 现读 `}`（块收尾）→ 自指复算那句在 **`:70`**＝漂 +1 方向为 −1（低危）。

**A 级＝引用有效但词面/位置需注记（⛔ 不算错，写下来免得下一枚腿误读）**

15. 票 285 `AC#1`（第 44 行）→ `ticket285_route_denial_name_rulers_test.go:125` 标为"尺句落点"⇒ **`:125` 确是那条 `t.Fatalf("AC#1 RED: ... (want %q)...")`＝对上**；但票面把词面 `denial=spent-or-never-live-nonce` 也系在这一行——`want` 常量在 **`:51`**，`:125` 只是 `%s` 打印它。红句报 `:125` 是对的。
16. 票 285 `AC#3`（第 56-58 行）→ `ticket242_binding_test.go:40/:57/:171`、`queue_test.go:167`、`ticket259_denial_rulers_test.go:139/:172/:252`、`ticket259_panel_capability_rulers_test.go:73/:82/:101/:111/:301/:308`、`ticket242_panelface_test.go:39` **全部逐枚对上**，且每条现读都确是 `t.Fatal*/t.Error` 红句本身（⛔ 不是同名符号凑分母：我按 `Fatalf|Errorf` 匹配并要求命中行＝被引行）。⇒ 票 285 那句"行号未漂"**成立**。
17. 票 285 `AC#2`（第 64-67 行）→ `ticket242_panelface_test.go:21` `var panelItemReadFace = []string{`（名单确由测试件手写＝对上）；`approval_reply.go:397` `func panelItemLine(item approval.PanelItem) string {`（`PanelItem` 唯一包外消费者＝对上）；`approval_reply_201_test.go:387` 逐字 `{"yes " + corr, "原生令牌无效"},`＝对上；`:229` 不变性轴＝对上（见 16 那条能力尺的 `:229`）。
18. 票 285 现量 ① → `queue_test.go:161`（`corr-B` 从未 push）＝对上；`subagent_selfapproval_197_test.go:473`（断的是 `Native().Allow(...)`）＝对上；`ticket87_veto_l2_test.go:173` 逐字 `if len(cards) != 2 || cards[0].CorrelationID == cards[1].CorrelationID {`＝对上；`queue.go:414` 两端恒等＝对上；`queue.go:423` `denial=%s`（285 自己的突变目标）＝对上；`approval.go:518`（`label()` 那枚 name-mapping switch 里 `:512-522`，两枚名 `:518`/`:520` 并列）＝对上。

## 2. 票面已自纠的 13 组：更正后的行号我全部复量**落准**

`resident_approval_windows.go:885→:886`（`:886` 现读正是那句 `fmt.Printf("wisp: 卡片挂起：%s %s（编号 %s）\n"...)`；`:885` 现读正是 `"channels", channelRosterText(p.Channels))`）· `:880→:881`（`:881` `slog.Info("approval: 常驻进程显示一张确认卡片",`）· `:697→:698`（`:698` `func (ra *residentApproval) AskOnTaskRoot(c residentCard) ...`）· `:700→:701`（`:701` `return ra.askConfirmation(taskCtx, c)`）· `resident_windows.go:248→:249`（`:249` `...This call is the caller: it reads`；`:248` 装的正是票面说的半句 `// the way out - and AskOnTaskRoot / askConfirmation still had zero product`）· `subagent_selfapproval_197_test.go:109→:116`（`:116` 逐字 `TaskID: taskID, CorrelationID: corrID,`；`:396`/`:463` 两枚 `-corr-1`/`-corr-2`；`:109` 现读 `return w`）· `ci.yml:654→:655`（`:655` `run: go vet -tags winlive ./cmd/wisp/ ./internal/ball/`；`:606` 步名逐字 `winlive compile gate (go vet -tags winlive, ticket 111 AC#11)`）· `config_readers_255.go:109→:110`（`:110` 名册行逐字带 `run.go:991 [cfg := rt.cfg]`；`:109` 现读另一条注释）· `loop_approval_test.go:215→:219-222`（`:219` `if r.CorrelationID == "" || r.CorrelationID == res.TaskID ||`、`:220` `!strings.HasPrefix(r.CorrelationID, res.TaskID) {`、`:222` 打印实参；`:215` 是注释行）· `gate.go:749→:755`（`:755` `func (g *Gate) Replay(...)`，注释块 `:751-754` 逐字含 "C18 一键重放…requires the task to be admitted again (D47 still applies to a replay), and it is NOT an answer"）· `approval_reply.go:227/:272→:229/:274`（两枚合并句副本确在 `:229`/`:274`；`:227`/`:272` 是 `s.record("REFUSED", ...)` 行）· 票 279 `run.go:991` 及随行四处散文 `:100`/`:106`/`:107`/`:108`（现读分别含 `run.go:435`/`run.go:424`/`run.go:991`/`run.go:1014`＝票面"只改 `:110` 会让这四处静默过期 3 行"**成立**）· 票 279 尺件面：`config_readers_255.go:229` `var evidenceCite = regexp.MustCompile(...)`＝对上，`config_receipt_255_test.go:179` 函数声明＝对上，断言体 `:207-211`（`got := strings.TrimSpace(lines[n-1])` → `Contains` → `Errorf`）＝对上，`citedRowsFloor = 8` 在 `:220`＝对上。

其余对上（抽读，⛔ 不逐条复述）：`cmd/wisp/main.go:66` `runResident()` · `resident_windows.go:260` `src := startResidentTaskSource(rt, ra)` · `:269-270` boot 报告含 `任务来源：%s` · `:250-252` 注释 · `resident_approval_windows.go:690-693`（`voice chain, ticket 228's config wiring`）· `:864` 属"一带"（`ballCardUI.Prompt` 函数头在 `:865`，注释块 `:854-864`）· `resident_approval_live_246_windows_test.go:115/220/311` 三枚 `ra.AskOnTaskRoot(residentCard{` · `resident_task_source_live_246_windows_test.go:145` `t.Logf("AC#7 LIVE steps 1-2: ...")` · `internal/tools/bridge.go:424` `switch sil.Level {`／`:443` `PendingWindow`／`:458` `PendingApproval` · `gate.go:294/:536` 两处 `g.ui.Prompt(ctx, p)`（**恰好两枚**）· `queue.go:598` `func (q *Queue) replay(`／`:69` `replayOf`·`ui.go:130` `ErrBadGrant = ...（缺失/已用/与本次请求不绑定）`、`:127-130`、`type NativeAPI interface`＝`:143`、`type PanelAPI interface`＝`:167`（票面 `:143-162`／`:167-171` 成立）· `gate.go:702` 属"一带"（`panelAPI` 方法簇 `:698-704`）· `pending_read.go:9` 逐字对上 · `approval.go:477` 逐字含 "ticket 259 does not move it" · `queue_test.go:283` `g.Replay(...)` · `PLAN.md:1368`＝C18 那一行（含"拒绝后**任务 root ctx**…"）· 票 282：`cancel.go:67` 逐字 `func TaskID(ctx context.Context) string {`、`cancel.go:69` `return h.taskID`（＝回退那一支，突变目标可用）、`subagent_197.go:262` `parentID := TaskID(ctx)`、`task.go:708` `caller := TaskID(ctx)`、`loop.go:369` 逐字 `newTaskJournal(l.opt.Journal, taskID, taskID)`、`loop.go:603`/`:676`/`:367`（`:676` 确是**唯一调用点**，`callCorr` 其余命中 `:594`/`:367` 都是注释）、`journal.go:59`/`:81`、`run.go:1003` `Journal: nil,`/`:750` `Journal:  mem,`（⚠ 双空格，字面非逐字）、`.gitattributes:4` `*.go text eol=lf`、`bd124b2a~1:internal/tools/loop_approval_test.go:13` `"strings"`（＝"非新加"成立）· 票 283：`ticket283_corr_identity_rulers_test.go:13`（`callCorr` 只出现在注释）、`:40`、`:194`、`:206`、`:209`（两枚都是 `t.Errorf` 断言行）、`:222` 逐枚断"非空／不等于 taskID／前缀"＝票 285 `AC#4` 那句"⛔ 没把两枚互相作差"在 `:222` 上**逐字成立**（该 `for` 里没有两两作差）· 票 280：`ci.yml:172` `- name: gofmt (gofumpt)`→`:173` `if: ${{ !cancelled() }}`→`:175` `go install mvdan.cc/gofumpt@latest`→`:176` `OUT="$(gofumpt -l . tools/d22scan tools/mockllm)"`、`:189` 步名（tracked set/denominator）、`:239` `if: ${{ !cancelled() }}`→`:240` `attrib.sh --tracked-only` 全部逐条对上 · 票 281：`13-summary.md:14` 现读 `| 97 | 228 / 1 / 4 |`＝"那枚 4 出自此处"成立；105/110 两张的 `Status` 行现读＝票面 `AC#2` 抄的两行（`accepted-done（附条件已清）`／`accepted-done（附条件）`）对上；`probes/281/r1/10-reconcile.md:68` 确是 `pending-human-review.md` 不存在那一条 · 票 279 的 `panel_inbound.go:8/:9`（含 code span 的注释行）、`:285`/`:289`（`panelInboundUsage` 起止）、P03 块 `:11-15` 互不重叠＝对上；`internal/agent/approval/doc.go:27` 存在且是注释行（`623a4d81` 动的正是该文件＋台账）。

## 3. 判不动（具名缺哪一发读数）

1. 票 242 现量表第 12 行的括注〔`197-v1` D2 读数：删掉比对 ⇒ 定向尺绿＋`go test ./internal/agent/approval/` 也绿〕——这是**运行结果**，我按纪律零 `go test`，只能判它的**行号**（见 S 级第 1 条：`approval.go:292` 已不是那枚函数）。缺：今天重跑 D2 的红绿读数。
2. 票 242 第 62 行 ③ 的块引用 `ticket242_binding_test.go:87-107`（"递的是从未 `issue` 过的字面 nonce"）——我只取了 `:71`/`:111`/`:40`/`:57`/`:171`，没取 `87-107` 整块，判不了那块现在还是不是"字面 nonce"那枚用例。缺：该块 21 行名册。
3. 票 280 现量的数字面（CR 字节 334/131/225/1115/1284、`gofmt -d`／`gofumpt -d` 行数 14/14/101/0/0、7 枚/2 枚名册）——**枚数归 `reconcile-8`**，我只核了它引的 CI 步号与腿件行段（全对）。缺：不是我缺的，是分工。

## 4. 射程自陈（一条没核到／射程够不着）

- **票 283**：无 ✅ 标记，凭据写在 `Status` 与三格的括注里；那些引用（`cancel.go:69`、`subagent_197.go:262`、`task.go:708`、`039ec93c`、`ticket283_*:40/:194/:206/:209`、`probes/orch-283r1/*.txt`）我核了，**全对上**。票面 `:206` 在 `Status` 里被叫作"seed B 那条上游 L2 通道"，而 `AC#2` 把 `:206` 记为 seed C 的 FAIL 行——**行号本身对**（`:206` 确是 `task.cancel` 那枚 `t.Errorf`），但两种说法指同一行的用途不同，我没裁（⛔ 不属我这枚腿）。
- **票 284／票 259／票 197／票 246**：不在本腿 9 张名册里，票 278/242/285 替它们转述的行号（如票 242 第 14 行那条"票 259 面上写的载具行 `:109`"）我只在**本票引用面**核过（`:116` 为真身），没去开那几张票的正文。
- **台账引用**（`A##`/`Q##`/`R##`，如 `A745`/`A754`/`A718 §3`/`Q-74`）一律**无行号**，我这把尺不接；它们对不对＝编排者自己的账，⛔ 我未裁。
- **票 281 `AC#1` 的"每勾必须引 `docs/evidence/s1/<NN>-adversarial-acceptance.md:<行>`＋该行读数"**：我核了承载这张表的 `10-reconcile.md:16-27` 行段，**没逐枚下钻 8 张 `adversarial-acceptance.md` 的单个行号**（10 枚×各自行号＝超我调用预算）。这是本腿最大的一片空白：⚠ 若编排者要这层，需要另派一枚腿专核 281 的 10 枚补勾行。
- **票 242 第 21 行那句抄进票面的引文**（`> 判据＝造两枚同时活着…`）：我只验了词面首句在 `v2.md:95`（落在声称的 `:91-100` 内），**没做整块逐字比对**（票面自称"逐字取自"）。
- **只有真机/真跑才能验的**：票 277 的卡行是否真出现（`rc=0`／2 PASS）、票 278"链接器里留不留符号"（需 `go build` 后符号表）、票 280 的 gofumpt 发声行数、票 285 `AC#3` 的四发红绿枚数、票 282 `AC#5` 的 27 次 numstat——⛔ 我这把尺（对象层文本比对）够不着，且按纪律不跑。
- 票 285 现量 `:12` 列的 4 组自纠里，`approval_reply.go:227/:272→:229/:274` 那组**文件是靠上下文推的**（票面没写文件名）；我按 `approval_reply.go` 复量成立，但若编排者本意是别的文件，这一组的归属要他点一下。

## 5. 方法注记（避免下一枚腿重蹈）

- 取数一律 `git show HEAD:<path>`；本仓 `core.autocrlf=true`，工作树里 `internal/tools/bridge.go` 有 CR、HEAD blob 没有——**拿工作树当证据会把"行尾 artifact"读成"内容差异"**（票 282 `AC#4(b)` 已把这枚坑写死，我这轮复用它那把尺：`tr -cd '\r' | wc -c`）。
- 判"这句话在不在这一行"＝**按内容搜整行**（`$0~q` 或 `index($0,q)`），行号与短语在同一次 awk 里取；⛔ `cut -c1-150` 之后宣布读过（票 242 那几条长引用一行 400-1500 字节，截断必漏）。
- 同符号凑分母已防：每条要求**命中行＝被引行**才算对上，否则给真行号；对"某类型的某方法"类引用（`PanelAPI` 的 `Allow`、`Gate` 的 `Replay`、`grantStore` 的 `spend`、`ballCardUI` 的 `Prompt`）我都读了行上的接收者/声明类型，不是只看词。
- 本腿零产码改动、零票面改动、零台账改动、零 `go`／`gh` 调用；⛔ 未 push。
