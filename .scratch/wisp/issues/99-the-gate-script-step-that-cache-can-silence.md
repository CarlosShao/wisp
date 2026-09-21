# 99 — `scripts/d22scan.sh` 第一步是**裸 `go test ./...`** ⇒ 本机测试缓存能把一次真违规"端"成绿

**Status:** open（2026-09-21 16:0x 编排者建；来源=`acceptor-ticket88` 交件时点名的残留，它明写"非本票账"）
**Type:** 门禁完整性（票 71 / 票 96 同族：**仪器自己会偷偷不跑**）
**Blocks:** nothing · **Blocked by:** nothing
**Packages:** `scripts/d22scan.sh`（就那个 `-count=1` 的差别）。**禁改**：`tools/d22scan/**`（扫描器本体）、
              `allowlist.txt`（**5 行非注释，只许变短**）、`.github/workflows/ci.yml`（改法归票 85）、
              任何 ban 的文本、任何测试的断言。

## 现场

票 88 的验收代理实测（它的话："另登记 `scripts/d22scan.sh` 第一步裸 `go test ./...` 可被本机缓存端过去"）：
脚本第一步跑 `tools/d22scan` 的**台账/阳性对照测试**，而 `tools/d22scan/runtests.sh:75` 那一侧是**强制 `-count=1`** 的。
⇒ 同一个仓里两套调用方式，一处严一处松（这与票 93 的"SKIP 一侧拒、portable 不拒"是**同一种结构缺陷**）。

**为什么这不是"理论风险"**：`TestScannerSelfScanOfRealRepoIsGreen` 是在**测试运行时**读仓库树的。
Go 的测试缓存**看不见运行时读的文件**——它在缓存判定里只记 import/构建输入。
⇒ 往 `internal/` 或 `frontend/` 加一个真违规（比如一个 `approval.decide`），
若**包与测试源码本身没变**，`go test` 完全可能直接回放上一次的 `ok`。
CI 是干净 runner 所以每轮真跑；**本地开发与代理自检拿到的可能就是回放**——
而我们现在要求每张票收尾前都跑这个脚本（A64② 新立的规矩），**等于让一条会缓存的步骤承担门禁职责**。

## AC（1:1，裁决表 `docs/evidence/s1/99-*.md` 由验收方出）

- [ ] **AC#1** 先把"缓存真能端过去"**量出来**，不许只论证：在**仓外快照**里
      ① 跑一次脚本让它绿；② 往树里种一个已知违规；③ **再跑一次** ⇒ 报第二次的 rc 与结论原文。
      若第二次仍是绿 ⇒ 缺陷成立并记下这条序列；若第二次就红 ⇒ **如实写"未复现"**并说明差异（Go 版本/缓存策略），
      然后**仍然**做 AC#2——因为判据不能依赖"缓存恰好没命中"。
- [ ] **AC#2** 修法：脚本第一步与 `runtests.sh` **对齐成同一档**（显式 `-count=1`，或干脆改调 `runtests.sh`，
      两处调用不再各持一套规矩）。⚠ **不许**用 `GOFLAGS=-count=1` 之外的方式"全局加快速"，
      也**不许**顺手给脚本加 `-failfast`/`-short`/`-run` 之类**改变覆盖面**的开关。
- [ ] **AC#3** 变异双向：(i) 把我加的 `-count=1` 去掉 ⇒ AC#1 那条序列必须重新出现"种了违规还绿"；
      (ii) 阳性对照：种一个 `frontend/` 里的 `approval.decide` ⇒ 脚本 rc=1 并点名（票 88 已证明扫描器本身有牙，
      本票要证的是**这条调用路径**也有牙）。锚点=承载行为那一行，同链 grep 证落地，还原后 `git diff --quiet` 证干净，
      **编译失败不算变异**。
- [ ] **AC#4** 门禁：`bash -n scripts/d22scan.sh` rc=0；在纯净快照 `sh scripts/d22scan.sh` 连跑**两次**都 rc=0
      （⚠ **HEAD 上现在它是 rc=1**，唯一命中是 `internal/winsec/winsec.go:126`，那是**票 94** 的账 ⇒
      若你开工时它还没修完，你的判据是"**除那一条已知命中之外没有新命中**"，并如实登记，**不要去碰 `internal/winsec/`**）。
      **收尾前必跑 `sh scripts/d22scan.sh`**（A64②）。

## Rules（本仓固定）

`git commit -q -F - -- <显式路径> <<'MSGEOF'`（引号 heredoc）；提交前 `git diff --cached --name-only`；
禁 `git add -A`/`.`、`--amend`/`reset`/`rebase`/`stash`/`checkout .`（A34）；**不 push**；
不在仓内建 worktree（A38④，快照目录带会话后缀）；票面 append-only（改行前先读；标题前插段落要重抄标题，
`git diff --numstat` 删除列必须 0）；四种假绿逐条点名；数字不达标写 FAIL 附数字。
15 次工具调用内交回第一枚 checkpoint；每次提交同步 Status + 勾框 + `next=`；接近轮数上限主动收尾留断点。
⚠ 共树：`internal/winsec/`（票 94）、`internal/config`+`internal/agent`（票 90）、`tools/d22scan/`（票 96）都有人在写，
**本票只碰 `scripts/d22scan.sh` 这一个文件**，别扩界。

## Progress log（append-only）

- 2026-09-21 16:0x（编排者）：建票。票 88 的验收代理交件时把它判为"非本票账"并交回给我 ⇒ 按规矩落成票，不留口头。
  我把它排在很后：它**只在本地自检路径上有害**（CI 是干净 runner），
  但**恰恰是我们现在要求每个代理都走的那条路**（A64② 刚立的"收尾必跑此脚本"）⇒
  一条会回放的步骤当门禁，等于给每张新票发一台可能说谎的自检仪。
  ⚠ AC#1 我特意写成"**量出来，不许只论证**；量不到也照修"——因为这类"缓存会不会命中"的说法，
  本仓已经抓到过好几次"听起来对、实测没复现"，判据不能建在猜测上。
  next= 排队（它小、独立，可与票 93 同批做：那两张都是"两套仪器一严一松"的形状）。
