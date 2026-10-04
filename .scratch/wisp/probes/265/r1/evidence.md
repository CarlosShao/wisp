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

- 3.1 新建判据件与用例名册：**未判**
- 3.2 256 AST 钉的解冻落地（A602 边界内）：**未判**
- 3.3 硬约束第 1 条的钉（holder 未绑定 ⇒ 指名用例红）：**未判**
- 3.4 AC#3 ⓐ 支的正控（摘掉接线 ⇒ 指名用例红）：**未判**

## 4. 突变自证（每一发：种下必红 ＋ 还原 ＋ md5 三读数）

（本节答：每一发突变的原样命令、红句**逐字**、还原用的副本出处（必须是 `git cat-file blob HEAD:<path>`）、起手＝还原后＝HEAD blob 三枚 md5。）

- M1 holder 未绑定 ⇒ 让 `Record` 静默 `return 0, nil`：**未判**
- M2 摘掉 `Grants:` 那一枚实传：**未判**
- M3 摘掉绑定位（`:281` 之后那一发）：**未判**
- M4 种"第二枚 mint"形（证本腿没走 ⓐ-Ⅲ）：**未判**

## 5. AC#1 那两句假话 → 真话的兑现

（本节答：`approval_reply.go:279-281`（给用户）与 `:277-278`（给审计）两句**文字一字未动**的自证；写侧接上之后两句为什么为真；判据怎么指得出"写侧没接上"这一情形。）

- 5.1 两句的文字零改动：**未判**
- 5.2 兑现链（holder → `session.Ledger.Record` → `GRANT-RECORDED`）：**未判**
- 5.3 判据对"没接上"的敏感度：**未判**

## 6. 门禁四数读数（带取数时刻）

（本节答：`sh scripts/d22scan.sh`、`sh scripts/check-path-length-budget.sh --with-self-test`、`gofumpt -l <本腿文件>`、`go vet ./cmd/wisp/` 四门读数＋时刻；定向 `go test` 的红名集合逐名比对，哪几枚是既有红。）

- 6.1 d22scan：**未判**
- 6.2 路径长度门：**未判**
- 6.3 gofumpt：**未判**
- 6.4 go vet（含 `-tags winlive` 只编译不跑）：**未判**
- 6.5 定向 go test 红名集合逐名比对：**未判**

## 7. 我判不动／量不到的地方（具名＋归口）

（本节答：⛔ 不许用"应该没问题"填空。每一格给：判到哪一步就停了、归口哪张票／哪枚编排者面。）

- 7.1 第三枚过期注释 `internal/agent/approval/approval_always.go:165`（不在本腿写面）：**未判**
- 7.2 其余：**未判**

## 8. 污染面自证与提交名册

（本节答：本腿动过的每一枚文件（包级路径＋numstat）；临时件只建在 `.scratch/wisp/probes/265/r1/`；commit 名册与 pathspec 逐枚；票面 AC 框零改动的自证。）

- 8.1 写面与 numstat：**未判**
- 8.2 提交名册：**未判**
- 8.3 票面 AC 框零改动：**未判**
