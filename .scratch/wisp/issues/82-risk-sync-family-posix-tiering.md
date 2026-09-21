# 82 — `internal/risk` 的 sync-root 家族在 ubuntu 上 8 条红：判据没写错，是**POSIX 侧这条探测根本没实现**

**Status:** open（**排队**：等票 70/80/81 里有写码代理交回再派，写码代理总数不超过 3）
**Type:** 门禁归属裁定（A52④）+ 测试分层落地
**Blocks:** 票 70 的 AC#6（`test-core` 在票 81/82 之前不可能绿）· **Blocked by:** nothing
**Packages:** `internal/risk/` 的**测试文件**（`syncdirs*_test.go`、写门/取证相关用例）。
**禁改**：`internal/risk/assessor.go`、`internal/risk/pathresolver*.go`、`rules_gateway.go`、
`docs/PLAN.md`、`docs/specs/*.md`、`internal/config/**`（票 80 的地界）、
`internal/agent`+`internal/memory` 的测试（票 81 的地界）、`tools/d22scan/**`、`allowlist.txt`。

## 实测证据（CI run `35558750456`，headSha `17efc2c`，`test-core` → `Portable package tests`，ubuntu）

红着的 8 条（`gh run view --log-failed` 逐条点名）：
`TestWriteGatePlainLocalWriteNotFlagged`、`TestWriteGateNotSelectedByPayloadKey`、
`TestWriteGateEveryPathTargetJudged`、`TestWriteGateAllPathsNonSyncStaysExempt`、
`TestSyncNormalNewFileWriteNotFlagged`、`TestSyncFallbackNotDisarmableByWeakRoot`、
`TestSyncEnvConfiguredRoots`、`TestSyncDotDotTailFailsClosed`。

根因不是断言写错：`syncdirs_other.go` 顶部自带 `DEFERRED(P12-macos)`，`registryProbe` 在非 Windows 直接返回 nil
⇒ POSIX 上"sync root"这个概念**尚未实现**（那是**票 55** macOS 移植的活）。
所以这 8 条在 ubuntu 上测的是一个不存在的能力。

## 编排者已裁定（照此执行，不要再问）

**不**把这些包/用例从 portable 清单里搬走（搬走 = 把覆盖面挪到没人看的地方）。
采用票 70-c 为 `internal/secret` 用过的**同一档修法：分层**——
Windows 专属判据进 `//go:build windows`，POSIX 侧留一条**断言当前真实行为**的用例。
机制票 75 已经铺好：`internal/risk/provenance_syncdirs_windows_test.go`（windows 层，7 个子用例）
与 `internal/risk/provenance_syncdirs_other_test.go`（`!windows` 层），
后者里那条由 `SyncDetectionComplete()` 门控的**休眠断言**（A51⑧）就是 POSIX 侧的占位。

## AC（1:1，裁决表 `docs/evidence/s1/82-*.md`）

- [ ] **AC#1** 上面 8 条**逐条定性**并列表：哪些属于"只在 Windows 上有意义"（进 windows 层）、
  哪些其实**两侧都该跑**（那就是真缺陷，不许搬走，就地修）。
  ⚠ 不许整族打包贴 tag —— 那正是 A49② 里票 78 代理拒绝过的"整包被静默排除"。
  **判据**：每条给一句"为什么这条在 POSIX 上没有可断言的对象"或"这条两侧都有对象、留下"。
- [ ] **AC#1b（2026-09-21 由 A57④ 追加，编排者写的判据）** 碰 `internal/risk` 的**折叠语义**时，
  必须同时检查 `internal/risk/rules_test.go:198` 是否**还在测东西**：它用
  `C:\Program Files\Git\bin\git.exe` 这条字面 Windows 路径做输入，今天两侧同字节同结论（不是问题），
  但它断言的是折叠行为——**一旦本票把折叠收窄成"只在 Windows 上折"，这条会在 ubuntu 上静默变空仪器**
  （跑过、绿、什么也没测到）。判据：给出该用例在**两侧各自的**断言命中数或红/绿读数，
  若它在 Linux 侧变成空转，就地改成平台派生的期望值（照票 81 的做法），**不许留着当装饰**。

- [ ] **AC#2b（编排者据 A49② 追加：本票最怕的失败模式）** 分层**必须用构建约束表达，不能用文件名表达**：
  `*_windows_test.go` 这个名字**什么门都不施加**——真正的门是文件里 `//go:build windows` 那一行。
  若有人只改文件名、或让整包退出构建，`go test ./internal/risk/` 会打 **`ok ... [no test files]`**：
  看着全绿，实际**整包判据被静默跳过**（票 78 的代理当初就是按这条拒绝了我那个"省事修法"）。
  ⇒ **机器可检查**：windows 层与 POSIX 层各跑一次 `-count=1 -v`，**两侧都要贴出自己的 `=== RUN` 计数与测试名清单**；
  任一侧出现 `[no test files]` / `no tests to run` / RUN=0 ⇒ **本票不过**。
- [ ] **AC#2** POSIX 层**不许是空文件、不许全是 `t.Skip`**：它必须断言 POSIX **当前**的真实行为
  （fail-closed / 尚未检出），并在注释里指名**票 55 落地时这一层要怎么转成活的**（A51⑧）。
- [ ] **AC#3** 双向证据：(i) ubuntu（docker `golang:1.27`，挂 `git archive` 快照）这 8 条**不再出现**在 FAIL 里，
  且 `internal/risk` 的**红数不增加**；(ii) Windows 本机 `go test ./internal/risk/ -count=2` 全绿，
  `=== RUN` 行数是 `-count=2` 的偶数倍，`--- SKIP` 逐条点名。
- [ ] **AC#4** 变异：把 windows 层的 tag 去掉 ⇒ 在 ubuntu 上它必须**变红**（证明"tag 不是在藏东西，
  而是这条判据在 POSIX 上确实无对象"）。锚点=承载行为的那一行，同链 grep 自证，还原后 `diff -q` 对拍。
- [ ] **AC#5** 门禁（只跑自己碰的包）：`gofmt -l` 空、`gofumpt -l <file>` 空、`go vet` rc=0、
  `GOOS=linux go vet` rc=0、`go test -count=2` rc=0。

## Rules（本仓固定）

15 次工具调用内交第一枚 checkpoint commit；每次 commit 同步 Status + 勾框 + 末行 `next=`；
`git commit -q -F - -- <显式路径>` + **带引号** heredoc；禁 `git add -A`/`.`；commit 前核对 `git diff --cached --name-only`；
禁 `--amend`/`reset`/`rebase`/`stash`/`checkout .`（共树 A34）；**不 push**；不在仓内建 worktree（A38④）；
docker 挂**快照**不挂工作树；票面 Progress log append-only，**要改的那行先读再替换**；
四种假绿逐跑点名（SKIP 当 ok / `-run` 空匹配 / `-count=N` 没核对倍数 / 步骤被静默跳过）。

## Progress log（append-only）

（空）
