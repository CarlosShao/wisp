# 证据件 265-r1 — 票 265 形 ⓐ-Ⅰ 落地（晚绑定 holder，`GrantRecorder` 由 `cmd/wisp` 自己实现）

**腿**：`265-r1`（产码落地腿）
**对象票**：`.scratch/wisp/issues/265-resident-gate-has-no-session-grant-writer.md`
**裁定出处**：票面「编排者裁定（形）＝ⓐ-Ⅰ」＋台账 `A601` §4／§5、`A602`（具名解冻那枚 AST 钉）
**本件九节状态**：骨架先落盘（每格初值＝**未判**），此后每做完一格 Edit 一节并 commit 一次。

---

## 0. 本腿的运行约束与开工闸门

（本节答：开工闸门两条的现跑读数＋时刻；`git status --porcelain -- cmd internal`；在飞腿名单与本腿的归因边界；sherpa PATH 这一格怎么解决。）

- 闸门① `git status --porcelain -- cmd internal`＝**空**（15:48:44 +08，HEAD `c528035f`；本腿动笔后 17:4x 复跑＝只有本腿那三枚写面为 `M`，见 §8.1）。
- 闸门② sherpa PATH：`cmd/wisp` 的测试二进制在**装载期**就要三枚 dll（`third_party/sherpa-onnx/{onnxruntime,sherpa-onnx-c-api,sherpa-onnx-cxx-api}.dll`，本腿现跑 `find` 量到），缺则 `exit status 0xc0000135` 且**没有 `--- FAIL` 行**。本腿取数一律用仓里既有那把尺的形式：`export PATH="$PWD/third_party/sherpa-onnx:$PATH"`（逐字抄 `scripts/wisp-cli-tests.sh:103-108` 的实测结论——`pwd -W` 那形在这台机上**会** 0xc0000135，MSYS 形才对）。
- 闸门③ `winlive`：**一枚都没跑**，只跑了 `go vet -tags winlive ./cmd/wisp/`（见 §6.4）。
- 归因边界：`internal/tools`／`internal/agent/spill.go`（`174-r4`）与 `internal/config`（`257-r1`）**本腿没碰、没跑、不归因**；非本票的红一律先 `-count=3` 安静复量再判。


## 1. 形状复认（ⓐ-Ⅰ 三枚写面的开工时刻现量）

（本节答：三枚文件在 `git cat-file blob HEAD:` 那儿的原始行号与原文；写面边界——`internal/agent/approval/**` 与 `cmd/wisp/run.go` 零改动怎么自证；有没有越界。）

- 1.1 `cmd/wisp/resident_approval_windows.go`：**对得上派单**。HEAD blob md5 `7a26c7a990dbd2351bdf9898b5bdc192`（盘上同值＝起手干净）；`:219-225` 那枚 `approval.New(approval.Options{…})` 字面量起手实传 **5 枚**（`UI`/`Channels`/`Window`/`ApprovalTimeout`/`Logf`，无 `Grants`），`:185-192` 逐字含 "WHAT THIS DELIBERATELY DOES NOT CLOSE" ＋ "Options.Grants stays unset here"。
- 1.2 `cmd/wisp/resident_task_source_windows.go`：HEAD blob md5 `5fe2fc858e953aaebdfb11de0da1f5ff`；`:281`＝`src.run = run`，`:321`＝`src.submitTask(injectedText)`，两者之间只有 `src.root = …` 与 reply/console 绑定那一段 ⇒ 绑定位落在这一带是**唯一**能满足"任何卡都可能被答复之前"的位置。
- 1.3 `cmd/wisp/resident_approval_risk_256_windows_test.go`：HEAD blob md5 `f26bfbc7a2f1a10840fb007fc728f1a8`；`:445` 的 `want` 起手 5 枚，`:452-456` 起手是 `if got["Grants"] { t.Errorf(...) }` 反向断言。派单说"连带搬家／改写成反向形状"——本腿取**改写成反向形状**那一支（见 §3.2），并**没有**动同文件 `[risk]` 那两枚字段的任何断言。
- 1.4 **写面零越界**：交件时刻 `git diff --numstat -- cmd internal` 里本腿的名册只有 §8.1 那四枚（三枚 `M` ＋一枚 `??`），`internal/agent/approval/**` 与 `cmd/wisp/run.go` **不在其上**；尺另跑一遍 `git status --porcelain -- cmd internal` 复量＝盘上另有 `M internal/config/settings.go` 与 `?? internal/agent/spill_pointer_authority_ac3_174r4_test.go`，**那是 `257-r1`／`174-r4` 的在飞面，本腿没碰、没跑、不归因**（派单闸门条）。

## 2. 产码：holder 的形状与晚绑定接线

（本节答：holder 类型逐字设计；未绑定窗口为什么必须走 error 路径（硬约束第 1 条）；`run.session == nil`（mint 失败）落在哪一支；为什么只能晚绑定（硬约束第 2 条）；账本唯一构造点没被动（硬约束第 3 条）。）

- 2.1 **holder**（`cmd/wisp/resident_approval_windows.go:190-290` 一段）：
  - 类型 `residentGrantHolder{ mu sync.Mutex; target approval.GrantRecorder }`，零值＝未绑定 ⇒ "从没接过的形"天然 fail-closed。
  - `Record`（`:253`）：**未绑定即 `return 0, errors.New(...)`**，句子中文、点名"还没绑上本进程的会话账本"；⛔ 不造 id、⛔ 不 `return nil`。落点＝`internal/agent/approval/gate.go:689-690` 的 `GRANT-RECORD-FAILED`（本腿没改那一行，读来的）。
  - `bind`（`:229`）与 `bindSessionLedger`（`:242`）各带 **typed-nil 守卫**：`nil` recorder 与 `nil *session.Ledger` 都不写进 `target`。这一条不是洁癖——`gate.go` 那边只看 `g.grants == nil`，一枚"装着 nil 指针的非 nil 接口"会让门以为有账本，然后第一发答复就在 nil 接收者上炸掉；同型陷阱在本仓有名字，`cmd/wisp/run.go:570` 那句 typed-nil 守卫是同一个东西的读侧版本。
  - `bound()`（`:268`）**只给判据用**，产码零调用者（尺＝`git grep -n "\.bound()" -- cmd`：命中全在本腿那枚新测试件里）。
  - `bindResidentGrantLedger`（`:283`）＝ `ra.grants.bindSessionLedger(l)` 的薄封装，是唯一的绑定向导。
- 2.2 **实传 5 枚 → 6 枚**：`:365` `ra.grants = &residentGrantHolder{}`（在建门之前），`:374` `Grants: ra.grants,`。派单说的"实传 6 枚"到位；尺＝§3 的 AST 钉（`TestTicket265GateLiteralCarriesTheResidentHolder` 断值是 `ra.grants`，不只是键在场）。
- 2.3 **绑定时序**：`cmd/wisp/resident_task_source_windows.go:294` `src.run = run` → `:309` `ra.bindResidentGrantLedger(run.session)` → `:351` `src.submitTask(injectedText)`。⚠ **本腿没有在绑定位新加任何打印**：mint 失败那一路 `run.go:475-476`／`:483-484` 已经在 boot 印过 `SESSION-MINT-FAILED`／`SESSION-LEDGER-FAILED`，再补一句就是同件事说两遍；而 boot 状态句 `residentStatusLine` 有票 256 名册的零漂移钉（普查件 §4 N#9 具名），本腿不去碰它。**判读了**：`run.session == nil` ⇒ holder 保持未绑定 ⇒ 那一发的答复落 `GRANT-RECORD-FAILED`，是真话不是静默。
- 2.4 **`:185-192` 那段注释改写后的原文**：段名从 "WHAT THIS DELIBERATELY DOES NOT CLOSE" 换成三节 ——「WHAT THIS CLOSES SINCE TICKET 265, AND WITH WHAT SHAPE」＋「WHAT STAYS OPEN ON PURPOSE…」（写清未绑定窗口与为什么那是硬约束第 1 条）＋「WHY THE GATE IS NOT JUST BUILT AFTER THE LEDGER EXISTS」（无控制台即 return ⇒ 载体不存在）。⛔ 原文那句 "Options.Grants stays unset here" 已不在盘上（尺＝`git grep -n "Grants stays unset" -- cmd` 交件时刻 0 命中）。上一段 "Two of the ten fields are now passed" 也**就地限定为 "for [risk]"**，否则它会与新第六枚字段互相打脸。
- 2.5 **顺带同批改释两枚**（编排者 `A601` §5 点到的那两枚，就在本腿写面里）：`resident_approval_windows.go:32-33` 与 `resident_task_source_windows.go:52-53`（现量行号有位移，见 §8.1）那句 "this process links no WebView2 host" 改成真话两面：**宿主确实链了**（票 33／`panel_host_windows.go`／`resident_windows.go` 真建真起），**缺的是 Go→页面的通道**（`internal/panel/panel_pump.go:12-21` 逐字），所以"卡片在屏上"仍然不是这两枚文件能声称的。⚠ 后者原先那段"printed at boot rather than smoothed over"跟着普查件 §1.4① 一起更正——**boot 从来没印过那件事**，本腿把它写成"这一段从前声称的东西不存在"。第三枚 `internal/agent/approval/approval_always.go:165` **不在写面**，见 §7.1。


## 3. 判据：新建的钉与解冻后的钉

（本节答：普查件 §5.2 具名的"常驻腿那一格零枚用例守"补了哪几枚；每枚断什么、正控形是什么；256 那枚 AST 钉解冻后的期望集与反向断言搬到哪儿去了。）

- 3.1 新建判据件与用例名册：**判读成立**（20:0x，对 commit `a04a095f` 的树现读＋起手档实跑）。
  判据件＝`cmd/wisp/resident_grant_writer_265_windows_test.go`（707 行，`//go:build windows`，`package main`），
  **8 枚用例**，逐枚（名＋断什么＋正控形）：
  1. `TestTicket265ResidentGrantHolderUnboundFailsLoudly`——**A601 §4 硬约束 1 的钉**。断：全新 holder `Record` 必返非 nil err 且 id==0（静默 `return 0,nil` 或造 id 两形都红）；err 句非空（`GRANT-RECORD-FAILED` 要插值它）；`bound()` 初值 false；`bindSessionLedger(nil)`／`bind(nil)` 都不得把 holder 变 bound，nil 后再 Record 仍返错。正控形＝种"未绑定返 nil err"即此名红（M1）。
  2. `TestTicket265ResidentGrantHolderBoundDelegatesEveryField`——断：bind 后 Record 逐字段透传（ctx 标记位 sawCtx、tool、pattern），id 是账本回的（fixture 从 20 起、首呼返 21，防小造 id 混过）；两次调用计两次。
  3. `TestTicket265LedgerErrorReachesTheOperatorThroughTheHolder`——断：账本自己报错时 holder 透传错误原句（`账本写不进去`），不吞。
  4. `TestTicket265UnboundHolderAnswersThroughTheGateWithoutClaimingARow`——**过真门的全链 POSTURE 钉**：真 `approval.New`（fixture 带真 holder 走 `Grants:`）＋真 L2 卡（`AdmitTextTask`→`PendingApproval`→UI Prompt）＋`Native().AllowSession`。断四向：答案仍放行（`AnswerAllow`，不能收回用户已给的答复）；审计**有** `GRANT-RECORD-FAILED`；**无** `GRANT-RECORDED`／`grant_id=`；**无** `GRANT-DROPPED`（防"Grants 根本没传"的旧形混绿）。
  5. `TestTicket265BoundHolderRecordsEveryPathTheAnsweredCardNamed`——AC#1 两句兑现的正控：bind 后同一张卡 AllowSession ⇒ recorder 收到**逐路径**调用（len==len(p.Paths)、tool/pattern 逐枚对上卡片），审计逐枚有 `GRANT-RECORDED corr=… grant_id=21+i`，且 `GRANT-RECORD-FAILED`／`GRANT-DROPPED`／`grant_id=0 ` 三禁。
  6. `TestTicket265GateLiteralCarriesTheResidentHolder`——AST 读产码字面量：`Options.Grants` 的值表达式**逐字**必须是 `ra.grants`（per-call 字面量＝绑定位永远够不着＝进程终身 GRANT-RECORD-FAILED 的形状，专为此红）。
  7. `TestTicket265BindSiteRunsAfterTheAssemblyAndBeforeAnyTask`——AST 读 `startResidentTaskSource`：绑定调用必须在场、实参逐字 `run.session`（ⓐ-Ⅲ 形红）、位置在 `src.run = run` 之后（早了＝账本还不存在＝装饰行）、且在每枚 `src.submitTask` 之前（晚了＝首卡可被答时仍未绑）。摘掉绑定位＝此名红（M3）。
  8. `TestTicket265SessionLedgerStillHasOneProductionConstructionSite`——ⓐ-Ⅲ 拒绝形的钉：扫包内全部非测试 .go，`session.Mint(`/`session.NewLedger(` 的命中必须全部前缀 `run.go:`（M4 的判据）；一个没有＝判据失去对象。
  9. `TestTicket265ResidentApprovalConstructsAnUnboundHolder`——跑**生产构造器** `newResidentApproval()`：holder 非 nil（`Grants:` 不是 nil 接口）、初值未绑定、未绑 Record 返错、`bindResidentGrantLedger(nil)` 不动它、bind 后透传。
  （件头自述 NOT claimed 两格：真双击进程无答复入口、无 Go→页面通道——与本腿 §2.5 改释一致。）
  **起手档实跑**：`export PATH="$PWD/third_party/sherpa-onnx:$PATH"` 后 `go test -count=1 ./cmd/wisp/`＝**ok 201.894s**（19:53:30–19:57:24，HEAD `a04a095f`）＝既有红名册**空**，本节所有用例起手全绿。
- 3.2 256 AST 钉的解冻落地（A602 边界内）：**判读成立，diff 逐块核过**。`git diff a04a095f^ a04a095f -- cmd/wisp/resident_approval_risk_256_windows_test.go` 全部改动＝四处：① 头注释块改写（"Grants stays unset" 叙述换成 ⓐ-Ⅰ 交付叙述，纯注释）；② ruler ④ 名下注释改写（仍保留原名不改名——件内自述理由＝A602 边界禁重排/改名、验收腿逐字比红名）；③ `:455`（现量行号）`want := []string{…, "Grants"}`＝5→6 枚；④ 原 `if got["Grants"] { t.Errorf("the resident leg now passes Options.Grants …") }` 反向断言整块替换为 `if !got["Grants"]` 正向断言（红句点名 A601/A602、"revert A602 first"）。**A602 边界四条逐条对**：只动期望集＋那一条反向断言＝是；不含同包其它断言＝是；`[risk]` 两枚字段（`Window`／`ApprovalTimeout`）判据一字未动＝是（diff 里 `:445` 一带的 Window/ApprovalTimeout 断言循环零改动）；未重排/改名文件＝是。**越界枚数＝0**。
- 3.3 硬约束第 1 条的钉（holder 未绑定 ⇒ 指名用例红）：**在册**＝3.1#1（对象级）＋3.1#4（门级全链）。红兑现见 §4 M1。
- 3.4 AC#3 ⓐ 支的正控（摘掉接线 ⇒ 指名用例红）：**在册**＝摘 `Grants:` 实传→3.1#4 的 GRANT-DROPPED 禁＋3.1#6+3.1#7 AST 钉＋256 ruler ④ 翻正后同红；摘绑定位→3.1#7。红兑现见 §4 M2／M3。

## 4. 突变自证（每一发：种下必红 ＋ 还原 ＋ md5 三读数）

（本节答：每一发突变的原样命令、红句**逐字**、还原用的副本出处（必须是 `git cat-file blob HEAD:<path>`）、起手＝还原后＝HEAD blob 三枚 md5。）

- M1 holder 未绑定 ⇒ 让 `Record` 静默 `return 0, nil`：**红兑现**（19:59:10）。种法＝`resident_approval_windows.go` 的 `Record` 未绑定支整段替换为 `return 0, nil // MUT-265-M1`（还原副本出处＝`git cat-file blob a04a095f:cmd/wisp/resident_approval_windows.go`，md5 `b3b1bc5610cbff62d69528022af6ba43`）。指名红兑现 **3 枚用例**（红句逐字）：
  - `TestTicket265ResidentGrantHolderUnboundFailsLoudly`：`an UNBOUND grant holder returned a nil error. That is the shape A601 §4 forbids: …` ＋ `after bindSessionLedger(nil) the holder still claimed success`；
  - `TestTicket265UnboundHolderAnswersThroughTheGateWithoutClaimingARow`：`audit is missing gate.go's GRANT-RECORD-FAILED line…`（审计实况里出现 `GRANT-RECORDED corr=265-corr-1 grant_id=0 …`，即 gate 把静默成功当真记账）＋ `audit claims a stored row while the holder was unbound`；
  - `TestTicket265ResidentApprovalConstructsAnUnboundHolder`：`the shipped holder recorded before any bind…`。
  还原后 md5：盘上＝副本＝blob `b3b1bc5610cbff62d69528022af6ba43` 三枚全等（19:59:16），`git status --porcelain -- cmd/wisp/`＝空。
- M2 摘掉 `Grants:` 那一枚实传：**红兑现**（19:59:47）。种法＝`newResidentApprovalWithConfig` 字面量删 `Grants: ra.grants,` 一行（还原副本同上出处）。指名红兑现（红句逐字）：
  - `TestTicket256ResidentGateOptionsFieldSetIsTheFiveItClaims`（256 ruler ④ 翻正后的正向臂）：`the resident leg does NOT pass Options.Grants to approval.New…` ＋ `the resident leg passes no Options.Grants, so 「本会话内允许」 on this gate books GRANT-RECORD-FAILED forever and the two sentences at approval_reply.go:277-281 stay false…revert A602 first.` ＋ `resident leg passes 5 of Options' fields, want 6`；
  - `TestTicket265GateLiteralCarriesTheResidentHolder`：`no Options.Grants in …resident_approval_windows.go: the field this ruler reads is gone…`。
  还原后 md5 三枚全等 `b3b1bc56…`，porcelain 空。
- M3 摘掉绑定位（`:281` 之后那一发）：**红兑现**（20:00:14）。种法＝`resident_task_source_windows.go` 的 `ra.bindResidentGrantLedger(run.session)` 替换为空转占位（还原副本出处＝`git cat-file blob a04a095f:cmd/wisp/resident_task_source_windows.go`，md5 `b8f6ed28642f8a5e1eb1638dd860fe69`）。指名红兑现 **1 枚**：`TestTicket265BindSiteRunsAfterTheAssemblyAndBeforeAnyTask`：`…resident_task_source_windows.go never calls ra.bindResidentGrantLedger inside startResidentTaskSource: … exactly what AC#3's 「摘掉那处接线」 arm asks to see red.`（AST 钉对"整行摘除"敏感；对"移到 submitTask 之后"的敏感由用例自身断言 `bindCall >= s` 分支覆盖，本发未单测该形——具名，不冒充）。
  还原后 md5 三枚全等 `b8f6ed28…`，porcelain 空。
- M4 种"第二枚 mint"形（证本腿没走 ⓐ-Ⅲ）：**红兑现**（20:01:11）。种法＝在**本腿写面内**的 `resident_approval_windows.go` 落一枚 `func mut265SecondMintSite265() { _, _ = session.Mint(); _, _ = session.NewLedger(session.LedgerOptions{}) }`（第一版 `_ = session.Mint()` 编译不过——Mint/NewLedger 各返两值——当场改 `_, _ =` 形再跑；判据本身零改动）。指名红兑现 **1 枚用例两行**：`TestTicket265SessionLedgerStillHasOneProductionConstructionSite`：`a second production construction site appeared: resident_approval_windows.go:252. …refuses form ⓐ-Ⅲ by name…`（:253 同句第二行）。选写面内文件而非 run.go 的理由＝种突变是测量动作但仍在物理改产码文件，落在 `cmd/wisp/resident_approval_windows.go`（本腿四枚写面之一）不越包界；判据扫的是全包非测试 .go，写面内种照样判到。
  还原后 md5 三枚全等 `b3b1bc56…`，porcelain 空（20:01:17）。
- **四发共同纪律**：两发之间零其它重活（只跑 `-run 'TestTicket265…'` 定向子集；M2 加跑 256 ruler ④ 一枚）；每发还原源＝**突变前从 `git cat-file blob a04a095f:<path>` 抽的副本**（`/tmp/m1-…`／`/tmp/m2-…`／`/tmp/m3-…`，M1/M2/M4 共用同一 blob 副本），起手＝还原后＝HEAD blob 三枚 md5 全部见各发行内；无任何"红未兑现"格——四发全响。

## 5. AC#1 那两句假话 → 真话的兑现

（本节答：`approval_reply.go:279-281`（给用户）与 `:277-278`（给审计）两句**文字一字未动**的自证；写侧接上之后两句为什么为真；判据怎么指得出"写侧没接上"这一情形。）

- 5.1 两句的文字零改动：**判读成立**（20:0x 现量）。尺①＝`git log --oneline -- cmd/wisp/approval_reply.go`：最近一枚触它的提交是 `3edc11d4`（票 246-r2），**早于本票起手锚 `de484541`**；尺②＝`git --no-pager diff --numstat de484541..HEAD -- cmd internal` 名册**只有四枚**（`resident_approval_risk_256_windows_test.go` 40/20、`resident_approval_windows.go` 163/13、`resident_grant_writer_265_windows_test.go` 706/0、`resident_task_source_windows.go` 37/8）——`approval_reply.go` 不其上＝本票全程（含前任产码＋本腿收编＋四发突变）**字节级未触**。两句现量：`:277-278` 审计句 `s.record("ANSWERED", corr, card.Tool, "native/allow-session", fmt.Sprintf("原生侧允许并记入本会话，paths=%d", len(card.Paths)))`；`:279-281` 回执句 `"已按「本会话内允许」答复 …这一发已放行，卡片上那些路径对本会话后续的 L1 询问不再重复提问…"`. `record` 落点＝`:414-419` → `s.audit(...)` → 常驻腿上即 `residentAuditf`（`resident_approval_windows.go:479-481`，`slog.Info("audit: "+line)` → JSONL）。
- 5.2 兑现链（holder → `session.Ledger.Record` → `GRANT-RECORDED`）：**论证成立＋测得**。链路＝常驻腿 AllowSession（`approval_reply.go:261` `s.live.h.AllowSession`）→ `gate.go` `allowSession`：`g.grants == nil` 才落 `GRANT-DROPPED`（现在 holder 已实传，不走这支——M2 红句反向证明：摘掉 Grants 才会掉回 `GRANT-DROPPED`）→ `g.grants.Record(...)`。写侧接上后（`Grants: ra.grants` 已传＋绑定位已落）：`run.session != nil` 的进程里 holder 已绑真账本，`Record` 逐路径写真行、返真 id ⇒ **审计句 "原生侧允许并记入本会话" 为真**（§4 M1 的反面：审计实况打出 `GRANT-RECORDED corr=… grant_id=21+…` 只在已绑形出现，§4 用例 5 断其逐枚在场）⇒ **回执句 "不再重复提问" 为真**（行已记入 ⇒ 同进程后续 `Covering` 查得到 ⇒ 不再问；读侧 `run.go:758 grantRead` 本来就接着，A601 §3 已判）。**测得凭据**＝§4 用例 `TestTicket265BoundHolderRecordsEveryPathTheAnsweredCardNamed`（起手绿）：bind 后同一张卡 AllowSession ⇒ recorder 收到逐路径调用＋审计逐枚 `GRANT-RECORDED grant_id=21+i` ＋三禁（`GRANT-RECORD-FAILED`／`GRANT-DROPPED`／`grant_id=0 `）全不出现。
- 5.3 判据对"没接上"的敏感度：**测得**。种"写侧没接上"的两形都红、且都是 §4 已实跑的指名红：① **Grants 没传**（M2）⇒ `TestTicket256ResidentGateOptionsFieldSetIsTheFiveItClaims` 红句逐字点名 "the two sentences at approval_reply.go:277-281 stay false"，同时 `TestTicket265GateLiteralCarriesTheResidentHolder` 红 `no Options.Grants…`；② **绑定位摘掉**（M3）⇒ `TestTicket265BindSiteRunsAfterTheAssemblyAndBeforeAnyTask` 红；③ **绑了但撒谎**（M1 静默 `return 0,nil`）⇒ `TestTicket265ResidentGrantHolderUnboundFailsLoudly`＋`TestTicket265UnboundHolderAnswersThroughTheGateWithoutClaimingARow` 红（后者红句里审计实况直接打脸：`GRANT-RECORDED grant_id=0`）。⇒ 三形各有一枚指名尺，无"判据对没接上不敏感"格。
- **未绑定时走 `GRANT-RECORD-FAILED` 不撒谎的实测**＝§4 用例 `TestTicket265UnboundHolderAnswersThroughTheGateWithoutClaimingARow`（起手绿）：未绑 holder 上真门放行真卡，审计**有** `GRANT-RECORD-FAILED`、**无** `GRANT-RECORDED`／`grant_id=`／`GRANT-DROPPED` ⇒ 两句假话在未绑定窗口里**一句都不出**（`err != nil` 路径走 `REFUSED` 分支，`approval_reply.go:273-275`，不产 ANSWERED 行）。

## 6. 门禁四数读数（带取数时刻）

（本节答：`sh scripts/d22scan.sh`、`sh scripts/check-path-length-budget.sh --with-self-test`、`gofumpt -l <本腿文件>`、`go vet ./cmd/wisp/` 四门读数＋时刻；定向 `go test` 的红名集合逐名比对，哪几枚是既有红。）

- 6.1 d22scan：**绿**（20:04:57）。先修已知红——`resident_approval_windows.go:38` 注释引 `internal/panel/panel_pump.go:12-21`＝盘上不存在（`ls internal/panel/` 无 panel_pump.go；A605 §2 预点名真身 `internal/panel/pump.go`，现读其 `:12-21` 逐字含 "no Go -> page channel … The last mile is NOT here"＝句子真、文件名错）⇒ 只改注释路径为 `internal/panel/pump.go:12-21`，语义零改。修后 `sh scripts/d22scan.sh` **rc=0、零 finding**（`clean - no D22 ban violations`，examin 266 production Go files；修前那枚红名＝ban #9 phantom-citation 会点 `cmd/wisp/resident_approval_windows.go`，修后名册消失）。同发复检：同文件另一处注释 `panel_pump.go:12-21`（`resident_task_source_windows.go:63`）引的相对名在 `cmd/wisp` 邻近语境有真身 `cmd/wisp/panel_pump.go`（`ls` 实存）＝非 phantom，不动。
- 6.2 路径长度门：**绿**（20:05:14）。`sh scripts/check-path-length-budget.sh --with-self-test` rc=0，`VERDICT GREEN - every over-budget tracked path is rostered by name`（denominator 5692 tracked paths、over-budget 57 全在名册、wall interval 0）。
- 6.3 gofumpt：**本腿写面空**（20:05:29，`"$(go env GOPATH)/bin/gofumpt.exe" -l cmd/wisp/`＝只点 `cmd\wisp\models.go` 一枚）。已知既有名册（A603 §6）＝`models.go`／`pending_read.go`／`queue.go` 三枚、票 212/258 账户；今天只点 `models.go` 一枚（另两枚本轮未点），**不归本腿、不修**（零 diff 复证：`git diff HEAD --stat -- cmd/wisp/models.go`＝空，最后触它＝`5e8748b3` 票 212/258 收尾）。本腿四枚写面＋证据件**零命中**。
- 6.4 go vet：**双绿**（20:06:14）。`go vet ./cmd/wisp/` rc=0；`go vet -tags winlive ./cmd/wisp/` rc=0（只编译，不跑任何 winlive 测试）。
- 6.5 定向 go test 红名集合逐名比对：**既有红名册＝空**（起手档 19:53:30–19:57:24，HEAD `a04a095f`，`go test -count=1 ./cmd/wisp/`＝`ok 201.894s` 零 FAIL）；终态复跑见 §9（收尾全量）。期间定向跑（M1–M4 四发）红名逐枚见 §4，全部为**种下即红、还原即清**的瞬时红，不在终态。`winlive` 标签测试**一枚没跑**（§0 闸门③照守）。

## 7. 我判不动／量不到的地方（具名＋归口）

（本节答：⛔ 不许用"应该没问题"填空。每一格给：判到哪一步就停了、归口哪张票／哪枚编排者面。）

- 7.1 第三枚过期注释 `internal/agent/approval/approval_always.go:165`（不在本腿写面）：**未动，具名归口**。本腿写面＝`cmd/wisp/**` 四枚＋证据件；`internal/agent/approval/**` 零改动（复证＝§8.1 numstat 名册无它）。A601 §5 已裁它"具名留给 `265-r1` 停手上报，⛔ 不许自行扩包"——本腿照办：该注释与已修的两枚同源（"this process links no WebView2 host" 在票 33 落地后不成立），但修它要扩包进 `internal/agent/approval`，归口＝**编排者下一枚允许触 `internal/agent/approval` 的腿**（或随票 224/256 那条线的下一程）。现状仍误导，这一格本腿判不动。
- 7.2 其余：
  - **`TestAC1ResidentLegInstallsItsLogListenerOnDisk` 偶发红**（票 127 钉件，`3edc11d4` 后未再触、与本腿零 diff 关联）：**量到的形状**＝四发全量复跑里出现两次（20:22:58 一次；另一次 count=3 内），失败句逐字 `the child never reached the event loop, so the disk read above proves little.`（`resident_sink_nail_127_windows_test.go:497-499`）。**根因判读**：该断言**不 poll**——同件相邻用例 `:530` 先 `pollUntil127(20, …stdout.has(residentReachedLoop))` 再断，而 `:497` 直接读捕获缓冲；子进程 stdout 异步写、`leg.stop()` 刚过后的 stdout 里 "resident event loop running" 可能尚未进缓冲 ⇒ 竞态在读侧不在产码。**隔离复量**：单跑 `-run` 该名 ×4 全绿（4.07s/4.91s/5.76s/6.25s）；`-count=3` 复量两次＝一次 FAIL（同一句）、一次 PASS——**仅在全量长跑（机器带载）下偶发**，与票 260 那族"带载时序"同形。**归口**：判据侧竞态，属票 127/246-r2 那条线的账（非实现方本腿可修：改该钉要动 `cmd/wisp` 判据件但其归属非本票）；⛔ 本腿不修不 Skip，具名上报编排者。
  - **`TestTicket223HandEditedFsLooseningCostsAnL2Card`**（第一次全量 20:13:44 红）与 **`TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink`**（第三次全量 20:18:00 红）：各自 `-count=3` 安静复量＝**全绿**（20:13:58／20:18:49），且后续全量不再复现（第四次红的是上一行那枚、第五次全绿）＝同族带载偶发，不判产物码缺陷，归口同上（CI 若复现，按"先 -count=3 复量再判"纪律处理）。
  - **`go test -race`／`-shuffle`**：没跑（本腿门禁清单外；四发突变已占约 40 分钟包内取数）。真想做归属票 265 后续验收腿。
  - **真双击进程的 U9 尺**（真起一发常驻进程看 GRANT-DROPPED 的去向）：票面 AC#1 已具名"唯一能证伪的尺＝U9，归 ⓐ 落地腿凭据档"——本腿**没跑**：它与 `winlive` 同级（真起进程真举卡），机主未批真开窗，⛔ 不越。归口＝编排者决定要不要单派一枚带批准的测量腿。

## 8. 污染面自证与提交名册

（本节答：本腿动过的每一枚文件（包级路径＋numstat）；临时件只建在 `.scratch/wisp/probes/265/r1/`；commit 名册与 pathspec 逐枚；票面 AC 框零改动的自证。）

- 8.1 写面与 numstat：**未判**
- 8.2 提交名册：**未判**
- 8.3 票面 AC 框零改动：**未判**
