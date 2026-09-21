# 88 — 把 `ban #6` 从"豁免"翻成"真覆盖"（`frontend/` 已经进树了，豁免就变成谎话）

**Status:** open
**Type:** 门禁完整性（A26 / A44① / A54② / A56④ 的同族，**这一次是轮到 ban #6 自己**）
**Blocks:** **编排者的 push**——票 77 的 `frontend/` 已落在 `63ef895`（未推），
             现在推上去会让 `lint` 的 D22 两步红，从而**把刚拿到的 `go vet` 绿证又挡成 skipped**。
             本票闭掉之前 push 压着不走。
**Blocked by:** nothing
**Packages:** `tools/d22scan/main.go`（`declaredScopes()` 那一条 + 必要的 `walkText` 过滤器）、
              `tools/d22scan/scan_test.go`（受影响的那几条台账用例）。
              **禁改**：`allowlist.txt`（**必须仍是 5 行非注释**）、任何 ban 的**文本**（D22：不是代理能缩的）、
              "空作用域致命"这条规则本身、`frontend/**`（票 77 的地界）、`internal/**`。

## 背景（编排者已经试过一次，把账留在这是给你少走弯路）

票 67 当年把 `ban #6 frontend/` 登记成 **exempt-with-a-note**（`live:false, absentOK:true`），
理由写得很好："目录不存在 ⇒ 这条只能走 0 个文件，留着就是假装在扫"。
并配了 **drift guard**：目录一旦冒出来而没人翻牌，扫描器**致命报错**（fail-closed）。
票 71 的 AC#4 又把"0 覆盖"从沉默变成明打字的 `NOT COVERED`。

**现在票 77 把 `frontend/` 建出来了** ⇒ 豁免的前提已经不成立。编排者在 `main.go:306-314`
试把 `live: false, absentOK: true` 改成 `live: true`，`cd tools/d22scan && go test ./...` **当场 5 条红**：

```
--- FAIL: TestVerdictGreenOnFullyLiveFixture
--- FAIL: TestVerdictRedOnEmptyBan7Scope
--- FAIL: TestExemptScopeCannotOutliveItsAbsentTree
--- FAIL: TestRealRepoLedgerIsHonest
--- FAIL: TestBuiltBinaryGoesRedEndToEnd
      ├─ /seeded_violation_exits_1
      ├─ /exempt_scope_whose_tree_appeared_exits_2
      └─ /fully_live_fixture_exits_0
报错原文：d22scan: scope ban #6 frontend/ examined 0 files but is declared live
          in declaredScopes() - an empty instrument is not a verdict (ticket 71 AC#4)
```

⇒ **这套台账是在"ban #6 永远豁免"的假设下写的**，翻牌要把它们**逐条重新表述**，不是删掉。
⚠ 那条"live 作用域若 examine 0 个文件 ⇒ 致命"**是对的，是你最不该碰的东西**：
它存在的意义就是"空仪器不算判据"（票 78 的代理当年也是按这条拒绝我的省事修法）。

## 先回答一个问题再动手（AC#1）

**真实仓里 `ban #6` 翻牌之后到底能扫到几个文件？** 编排者的实验里 fixture 是 0，
真树里 `frontend/` 有 38 个非 `node_modules` 文件（`.tsx`/`.ts`/`.css`/`.mjs`/`.md`）。
⇒ 如果 `walkText` 的后缀过滤器**不收 `.tsx`/`.ts`**，那真树也会量到 0 ⇒
**那才是本票真正要修的洞**：一条"看起来已覆盖、实际永远 0"的门。
先给出真实读数（`examined N text files`）再决定改哪里，**不许**为了让 N>0 而把 ban 的语义改宽或改窄。

## AC（1:1，裁决表 `docs/evidence/s1/88-*.md`）

- [ ] **AC#1** 真实读数：在 `git archive HEAD` 的**仓外纯净快照**（目录带你的会话后缀）里跑
      `sh scripts/d22scan.sh`，报出 `ban #6 frontend/ examined N text files` 的 **N**，
      并说明过滤器现在收哪些后缀。若 N=0 ⇒ 明写"这条门翻牌后仍是空仪器"，**这就是本票的缺陷本体**。
- [ ] **AC#2** 翻牌落地：`declaredScopes()` 里 ban #6 变 `live: true`（去掉 `absentOK`），
      `note` 改写成"谁在哪一批翻的、为什么"（保留 D22 那句"ban 文本不是代理能缩的"）。
      **ban 的判定文本一字不改**。
- [ ] **AC#3** 上面 5 条红**逐条重新表述并真绿**：
      涉及豁免语义的（`TestExemptScopeCannotOutliveItsAbsentTree`、
      `/exempt_scope_whose_tree_appeared_exits_2`）要**保留 drift guard 的可测性**——
      办法是**用一个合成的 exempt 作用域来测这条规则**，而不是靠"ban #6 恰好豁免"这个巧合；
      **规则本身必须仍然有牙**（不许把它测成永真）。
- [ ] **AC#4** 阳性对照：`/seeded_violation_exits_1` 用的就是
      `frontend/src/app.js` 里那句 `approval.decide(...)`（`scan_test.go:80`）⇒
      **翻牌后必须还能因这条种子违规红**，并且红的时候打印出 `ban #6` 命中。
      这是本票"门真的有牙"的证明，**不许用"改种子文件名/后缀"绕过 AC#1 的过滤器问题**。
- [ ] **AC#5** 全套绿之后跑 `sh scripts/d22scan.sh`（**逐字同 CI**）于纯净树：rc=0，
      且台账里 `ban #6` 那行的 N>0；同时确认 `allowlist.txt` 仍 5 行非注释、
      `git diff` 里**没有任何 ban 文本被改动**。
- [ ] **AC#6** 如果 `frontend/` 里的**真实代码**（票 77 vendored 的 Approval Card）命中 `ban #6`：
      **不要动扫描器也不要加豁免**——把命中点登记给编排者，由**票 77 那侧改名**。
      本票只负责让门有牙，不负责替别人消音。

## Rules（本仓固定）

15 次工具调用内交第一枚 checkpoint commit；每完成一组**就把结论追加进票面**（别攒，本仓代理死在轮数上限上过）；
`git commit -q -F - -- <显式路径>` + 带引号 heredoc；禁 `git add -A`/`.`；commit 前核对
`git diff --cached --name-only`（**别把活着代理的票面 add 进去**，编排者今天栽过一次）；
禁 `--amend`/`reset`/`rebase`/`stash`/`checkout .`；**不 push**；不在仓内建 worktree（A38④）；
四种假绿逐跑点名；数红/绿用全量输出仪器；票面 append-only，**要改的那行先读再替换**。

## Progress log（append-only）

- 2026-09-21（编排者）：建票。我做过一次 `live:true` 的试验并**已 `git restore` 还原**（工作树干净），
  5 条红的名单与报错原文在上面。我的判断：本票的价值不是"翻一个布尔"，
  而是**回答 AC#1 那个 N**——如果过滤器不收 `.tsx`，那我们今天所有关于"ban #6 已覆盖面板"的说法都是空的。
