# 85 — lint 的两把工具**从来没有产生过判据**：`gofumpt@latest` 没钉版本、`staticcheck` 在 go1.27 上跑不动

**Status:** open（**排队**：本票要改 `ci.yml` 与**全仓**多个包，必须等写码槽位空出来，
否则会与 `agent-ticket81`/`83`/`77` 撞同一片文件）
**Type:** 治理/门禁完整性（A26、A44 的同族第 3、4 起）
**Blocks:** 票 70 的 AC#6（`lint` job 在四步里只有两步产出过判据）· **Blocked by:** nothing
**Packages:** `.github/workflows/ci.yml`（版本 pin 与步骤名）+ findings 清算涉及
`cmd/wisp`、`internal/agent`（含 `approval`）、`internal/llm` 族、`internal/audio`、`internal/tools`、
`internal/memory`、`internal/observe`、`internal/ball`、`internal/config`、`internal/models`、`internal/proc`。
**禁改**：`docs/PLAN.md`、`docs/specs/*.md`、`internal/risk/assessor.go`、`internal/risk/pathresolver*.go`、
`rules_gateway.go`、`tools/d22scan/**`、`tools/d22scan/allowlist.txt`（**只许变短或不变**）。

## 两条实测事实（来自票 70 的接续代理，2026-09-21）

1. **`ci.yml:72` 写的是 `go install mvdan.cc/gofumpt@latest`——没钉版本。**
   票面 AC#1 记的"用 CI 钉的 v0.7.0 复跑"是**错的**（代理纠正了它：CI 日志实际下载的是 **v0.12.0**）。
   它用 v0.7.0 与 v0.12.0 **双版本**在纯净树上复核，当前判定一致（0 行）。
   ⇒ **风险是随机性**：哪天上游改版，会**连带历史 commit 一起变红**，而我们会把它读成"有人改坏了代码"。
2. **`staticcheck` 自始至终没有产出过一条判据。** CI 钉的 `2025.1.1` 在 go1.27 上直接跑不动
   （`export data version 4 > 2`，rc=1、**0 条真实 finding**）。升到 `@latest`（2026.2.1）**能跑**，
   但立刻报 **35 条 finding**，散在 15 个包：
   `27×U1000`（未使用）、`3×SA1019`（弃用 API）、`3×S1011`（**nil 解引用**）、`1×SA4006`、`1×SA4000`（比较自己）。
   ⚠ `internal/risk` **0 条**。

- [ ] **AC#6（2026-09-21 由票 75 的第二验收会话挖出，编排者追加）**
  `ci.yml:136-142` 的 `Environment fork assertion (WISP_ENV=test data dir)` 那一步
  **排在 `Portable package tests` 之后且没有 `if: always()`** ⇒ **只要测试那步红，它就永不执行**。
  也就是说：**这条门在红 run 里从来没有产出过一个结论**——和 A44③ 的 D22 扫描一模一样的病，只是换了个位置。
  判据：给它加 `if: always()`（或等价保证），然后**在 CI 上跑出一次"test-core 仍红、而这一步打印出自己的
  PASS/FAIL 计数"的 run**，把 run id 与该步结论贴进本票 log。
  ⚠ **不许**反过来把这一步删掉或挪走；`if: always()` 是**加严**，不是放宽。


## 编排者已裁定（照此执行，不要再问）

- **两把工具都钉成具体版本**（`gofumpt@vX.Y.Z`、`staticcheck@2026.2.1`），**禁止 `@latest`**。
  理由同上：门不能有"哪天上游自己变了"的随机性。
- **先升 pin、让 lint 红在真问题上，再逐条清账**；**不许**反过来用"保持跑不动"来维持绿色。
  一个跑不动的步骤和一个红的步骤，只有前者是丑闻。
- **`S1011`（nil 解引用）3 条优先**，它们是潜在 panic，不是风格问题；
  `SA4000`（拿变量和自己比）2 条同批——那种写法基本就是写错了。
- `U1000` 27 条**逐条判，不许一把删**：要么是真死码（删，并附一次"删了之后测试仍绿"的证据），
  要么是"等着接线"（那就**登记到该接线的票号**并在代码里留一行指名，不许默默留着）。
  ⚠ 参照 A53③：`internal/config` 那批"被解析、零消费者"的键正是这类东西的下场。

## AC（1:1，裁决表 `docs/evidence/s1/85-*.md`）

- [ ] **AC#1** `ci.yml` 里两把工具都钉成具体版本，**并给出一次真 run 的日志行**证明它下载的就是那个版本
      （`gh run view <id> --log` 里 grep `downloading` / `go version`）。**不许**只给本地读数。
- [ ] **AC#2** `staticcheck` **第一次产出真实判据**：把那 35 条按"包 × 检查码 × file:line"列表进本票 log，
      逐条给处置（修 / 删 / 登记到票 N）。**目标态是 0 条**，但**中途不许**加 `//lint:ignore`、
      不许把包排除、不许把这步改成 `continue-on-error`（D22：放宽门不是修门）。
- [ ] **AC#3** 至少 3 条修掉之后，**变异回验**：把其中一条改回旧形状 ⇒ `staticcheck` 必须重新报它
      （证明"它真的在看这个包"，而不是又一场空跑）。同一条 `&&` 链里 grep 自证落地，还原后证干净。
- [ ] **AC#4** 步骤名台账：`test-windows` 里那个 `PathResolver junction placeholder (real cases tickets 18/20)`
      步骤名里的 **"placeholder" 已过期**（它现在是真用例且首次真绿）。
      改成显式指名归属的名字（建议含票 18/20/72/75），**只改名字不改内容**，并给一次真 run 的该步结论。
- [ ] **AC#5** 门禁：`gofmt -l` / `gofumpt -l . tools/d22scan tools/mockllm` 空、
      `go vet ./...` rc=0、`sh scripts/d22scan.sh` rc=0（**逐字同 CI 那一行**，且**先看到它的阳性对照红过**）、
      `go test -count=2` 只跑你碰过的包并逐条点名 SKIP/FAIL。

## Rules（本仓固定）

15 次工具调用内交第一枚 checkpoint commit；每次 commit 同步 Status + 勾框 + 末行 `next=`；
`git commit -q -F - -- <显式路径>` + **带引号** heredoc；禁 `git add -A`/`.`；commit 前核对 `git diff --cached --name-only`；
禁 `--amend`/`reset`/`rebase`/`stash`/`checkout .`（共树 A34）；**不 push**；不在仓内建 worktree（A38④）；
纯净树用 `git archive HEAD | tar -x -C /tmp/...`；票面 Progress log append-only，**要改的那行先读再替换**；
四种假绿逐跑点名；⚠ Windows 主机上 `GOOS=linux go vet ./...` 因 CGO=0 排除 sherpa 预编译包而**永远 rc=1**，
**不要**拿它当判据（A54③）。

## Progress log（append-only）

- 2026-09-21 18:1x（编排者，**CI 步级读数到手：本票的 AC#2 前提被实测背书，且新出一条同账**）：
  ① `lint :: staticcheck` 在 run **`35586044995`、`35585821747`、`35585147258`、`35581075691`**（连 4 次）
  **全是 failure**，且错误**不是我们代码的问题**，是工具与工具链版本不匹配，CI 原文逐字：
  `go: downloading honnef.co/go/tools v0.6.1` 紧跟
  ` -: internal error in importing "internal/byteorder" (cannot decode "internal/byteorder", export data version 4 is greater than maximum supported version 2); please report an issue (compile)`
  （同一批还有 `internal/cpu`、`internal/goarch`、`math/bits`、`unicode/utf8` 四条同形）
  ⇒ 本票 AC#1（两把工具钉版本）**就是为这个而存在的**，AC#2"第一次产出真实判据"的前提已确认：**staticcheck 至今一条 findings 都没产出过**。
  ② 新同账（**归本票，因为它就是 `ci.yml` 的账**）：票 99 把 `scripts/d22scan.sh` 第一步改调
  `tools/d22scan/runtests.sh -C tools/d22scan ./...` 之后，**`ci.yml:48` 与 `:68` 现在跑的是同一台仪器的同一次测试**
  ⇒ CI 多花一整遍 `tools/d22scan` 的 `-v` 测试。判据写死：**要合就合"两步合一"，不许把 `scripts/d22scan.sh` 改回裸 `go test`**
  （那会把票 99 刚堵的缓存洞重新打开）。
  ③ 与票 99/93 的对账：票 99 交件时点名"`(cached)` 在修法版日志里出现 **0** 次"、纯净快照连跑两次 rc=0；
  票 93 正在动 `ci.yml` 的 portable 步 ⇒ **本票开工前必须先读 93 交回来的是什么形状**，别把两步的修法互相覆盖。
  next= 等写码槽位（当前 `internal/winsec`+`internal/tools`+`ci.yml`+`scripts/` 四路在飞）；
  顺序建议：**本票排在票 93 交件之后**（同文件），排在票 106/107 之后（那两条挡着 CI 转绿，本票挡的是"lint 有没有判据"）。

