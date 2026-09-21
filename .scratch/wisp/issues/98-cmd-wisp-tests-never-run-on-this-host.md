# 98 — `go test ./cmd/wisp/` 在这台机器上**永远测不到任何东西**（加载期 `sherpa-onnx-c-api.dll not found`）⇒ 宿主包的门禁是一架空仪器

**Status:** open（2026-09-21 15:4x 编排者建；根因由 `acceptor-ticket87` 挖出，它把它当"环境事实"登记了，
我判断**它本身就是一张票**：见下面"为什么这不是环境问题"）
**Type:** 门禁完整性（票 71 / A44① / 票 96 同族：**"没跑过"与"跑过且没问题"长得一样**）
· **Blocks:** nothing（但它让**四条已登记的"不追"变成无人还的债**）· **Blocked by:** nothing
**Packages:** `cmd/wisp/`（测试怎么跑起来）、`Makefile`/`scripts/`（DLL 路径）、
              `.github/workflows/ci.yml`（**只读核对**：CI 上这个包到底跑没跑过——要改 `ci.yml` 的部分归票 85）。
              **禁改**：`docs/PLAN.md`、`docs/specs/*.md`、`internal/risk/**`、`tools/d22scan/**` 与 `allowlist.txt`、
              `internal/models/**` 与 sherpa 的绑定本体（那是 C13/C29 的地界）。

## 现场（三张票各自撞过、都写了"不追"）

- 票 87 的代理与 `acceptor-ticket87` 都实测：`go test ./cmd/wisp/` 在本机 **rc=1**，
  纯净树 `git archive 63ef895` **同读数** ⇒ 先于票 87 存在（**这条判归因的结论已经站住**）。
- 验收代理把根因挖到了：**加载期 `sherpa-onnx-c-api.dll not found`**（Windows 的 `0xc0000135` = DLL 缺失，
  发生在进程启动、**测试还没开始**）。
- 我自己在 A64 的门检里也遇到同一现象，并在票 89/90/92/94/96/97 的简报里反复写"不要去追这条红"。

## 为什么这不是"环境问题"而是**仪器问题**

`0xc0000135` 是**加载期**失败 ⇒ 这个包里的**所有** `*_test.go` 在开发机上**一条都不会执行**。
于是今天凡是"落点在宿主 `cmd/wisp/`"的判据，**在本地都不可证伪**——
它不是"暂时红"，是**那架仪器在这个平台上永远不产出结论**（票 71 AC#4 给"0 覆盖"立的规矩，正对应这里）。
⚠ 更要紧的一半：**CI 上到底跑没跑过这个包？** 如果 CI 的 `test-windows` 清单里没有 `./cmd/wisp/`，
那"宿主层有测试"这句话就**两侧都没有来源**。**先把这个事实量出来**（AC#1），再谈修。

## AC（1:1，裁决表 `docs/evidence/s1/98-*.md` 由验收方出）

- [ ] **AC#1** **事实层**：贴出 (i) `go test ./cmd/wisp/` 本机原文（含 `0xc0000135` 与那句 DLL 名）、
      (ii) `cmd/wisp` 里现有测试文件与用例名的**清单**（说明"如果它能跑，会跑掉多少条"）、
      (iii) `ci.yml` 里**是否**有跑到 `./cmd/wisp/` 的那一步——给 file:line 与（若可取）一次真 run 的步级结论。
      ⚠ 不许只写"CI 会兜"：说不出 run id 就当它不存在（A44/A45 的固定判据）。
- [ ] **AC#2** 让本机**至少能跑起来一次**：给一条** documented 的命令/脚本**（例：把 sherpa 的 `.dll` 路径
      在测试前注入 `PATH`，或提供 `-tags` 的 no-cgo 构建），**贴出它跑出来的真实 PASS/FAIL 数**。
      ⚠ **不许**用"把测试搬出 `cmd/wisp`"或"给整包加 build tag"来变绿——那是把被检对象从门禁里删掉
      （票 78 的代理当年就是按这条拒绝过我，那次它是对的）。
- [ ] **AC#3** 一条**永不 skip** 的守卫：如果 DLL 缺失是**平台现实**，那就要有一条明确报"此平台该包未测"的**响亮失败**
      （像 `tools/d22scan/runtests.sh` 拒 SKIP 那样），而不是 rc=1 的加载崩溃被后人当成"环境问题不追"。
      判据：把 AC#2 的注入拿掉 ⇒ 守卫必须红并**说清是哪条 DLL**。
- [ ] **AC#4** 台账更正：把散在票 89/90/92/94/96/97 简报里那句"`cmd/wisp` 的红不要追"**收敛成一条引用**
      （引用本票票号），并在本票里逐条列出"哪些 AC 的判据此前**实际依赖过** `cmd/wisp` 的测试"（如果有，它们要重判）。
- [ ] **AC#5** 门禁：`gofmt -l`/`gofumpt -l` 空、`go vet ./cmd/wisp/` rc=0、`go test -count=2` 按 AC#2 的命令跑并点名
      每一条 SKIP/FAIL。**收尾前必跑 `sh scripts/d22scan.sh`**（A64②）。

## Rules（本仓固定）

`git commit -q -F - -- <显式路径> <<'MSGEOF'`（引号 heredoc）；提交前 `git diff --cached --name-only`；
禁 `git add -A`/`.`、`--amend`/`reset`/`rebase`/`stash`/`checkout .`（A34）；**不 push**；
不在仓内建 worktree（A38④，快照目录带会话后缀）；票面 append-only（改行前先读；标题前插段落要重抄标题，
`git diff --numstat` 删除列必须 0）；四种假绿逐条点名；数字不达标写 FAIL 附数字。
15 次工具调用内交回第一枚 checkpoint；每次提交同步 Status + 勾框 + `next=`；接近轮数上限主动收尾留断点。
⚠ 共树：`internal/agent/`+`internal/config/`（票 90）、`internal/winsec/`（票 94）、`tools/d22scan/`（票 96）、
`frontend/`+`internal/panel/`（票 77 已交件但票 92 会接）——**别碰**；`ci.yml` 的**改法**归票 85，本票只读核对。

## Progress log（append-only）

- 2026-09-21 15:4x（编排者）：建票。**这条我全天以"环境事实"的名义转述了 6 次**（票 87 的简报、票 89 的验收简报、
  票 90/92/94/96/97 的 Rules 段），每次都是"不要去追它"。
  `acceptor-ticket87` 把根因挖到 `sherpa-onnx-c-api.dll` 之后我才看清：**"不追"累积起来就是"这个包的判据没人能证伪"**。
  ⚠ 本票最怕的修法我在 AC#2 与 AC#3 里各钉了一句：
  **把测试搬走 / 加 build tag = 把被检对象从门禁里删掉**（票 78 的代理拒绝过我的同款省事修法）；
  允许的方向是"让它真跑起来"或"跑不起来就响亮失败"，**两者都不是把清单改短**。
  next= 可派（它不碰任何在飞代理的文件；`cmd/wisp/` 此刻无人写）。但**先等本轮 push 出去**——
  它可能需要看 CI 上 `test-windows` 的步级结论，而那要一次真 run（AC#1(iii)）。
