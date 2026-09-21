# 94 — 私有数据目录的"封 ACL"用了 `filepath.Abs` 而不是 C26 PathResolver ⇒ **D22 门在 CI 之前把它拦下了**（票 89 的码，push 因此压住）

**Status:** in-progress（agent-ticket94 已接手；票 89 已交件 `0a3a445`，`internal/winsec/` 现在是本代理地界）
⚠ 第一枚 checkpoint 已落盘（见 Progress log 末条）：AC#1 判明为 **(b)**，并**推翻**本票"无循环依赖"的前提。
**Type:** 安全边界（D22 ban #2 `pathresolver-bypass`）+ 一个**没被回答的设计问题**：拿没解析过的路径去决定"给哪棵树封权限"，本身是不是缺陷
**Blocks:** **编排者的 push**（HEAD 上 `sh scripts/d22scan.sh` 现在就是红的，见下面读数 ⇒ 推上去 `lint` 会红在 D22 两步）
· **Blocked by:** 票 **89**（它的代码引入了这个调用点；本票不能和它同时在 `internal/winsec/` 里写）
**Packages:** `internal/winsec/`（那个调用点及其调用方）。
              **禁改**：`tools/d22scan/**` 与 `allowlist.txt`（**allowlist 现在 5 行非注释，只许变短或持平 ⇒
              "把 winsec.go 加进 allowlist"这条省事修法在本票里是被禁的**）、`internal/risk/pathresolver*.go`（冻结）、
              `docs/PLAN.md`、`docs/specs/*.md`。

## 实测读数（编排者亲自跑的，可复现）

```
rm -rf /tmp/wisp-pushgate90 && mkdir -p /tmp/wisp-pushgate90
git archive HEAD | tar -x -C /tmp/wisp-pushgate90
cd /tmp/wisp-pushgate90 && sh scripts/d22scan.sh
```

```
--- FAIL: TestScannerSelfScanOfRealRepoIsGreen (4.29s)
    scan_test.go:269: repo HEAD violates: internal/winsec/winsec.go:126:
        [pathresolver-bypass] filepath.Abs outside the C26 PathResolver is banned (D22);
        see tools/d22scan/allowlist.txt for the sanctioned exceptions
--- FAIL: TestRealRepoLedgerIsHonest (1.33s)   （同一条命中，包在 rc 里）
```
台账逐作用域自报（**这次翻过牌的 `ban #6` 也在里面**）：
`bans #1-5 internal/=190`、`cmd/=20`、`ban #6 frontend/=35 text files`、`#7 internal/tools/=16`、
`#8 design/=16`、`#8 internal/=324`、`#8 cmd/=25`。

⇒ **这是那道门第一次在它该拦的时候拦下活着代理的码**（票 70-d 那两次卡点都是"注释里的字形"，
这次是真判据）。它同时也暴露一个**流程洞**：票 89 的门禁是**按包 scope** 的
（`gofmt`/`gofumpt`/`go vet <pkgs>`/`go test <pkgs>`），而 `ban #1-5` 的作用域是**整个 `internal/`**
⇒ **按包门禁结构性看不见全仓 ban**，代理就算把自家门禁跑绿也照样能把 HEAD 弄红。
（本票不解决这个问题，但解决办法写在下面 AC#4，别再让它靠编排者手气发现。）

## 背景：为什么"顺手用 `filepath.Abs`"在这里不是格式问题

`internal/winsec/winsec.go:126` 在 `PrivateDirAll(path, perm)` 里，注释写得很清楚它要做的事：
**逐级创建并封存"它自己创建的每一级"**，而"封存哪棵树"这个决定**就是安全决定**
（C26 PathResolver 存在的理由：reparse point/junction、8.3 短名、`\\?\` 扩展前缀都会让
"我以为我在封 A"变成"我实际在动 B"）。`filepath.Abs` 只做**词法**拼接——它**不解析 reparse、不展开 8.3**。
⇒ 有两种可能，**本票要你判明是哪一种，而不是先修**：
- **(a) 只是"该用哪个 API"的问题**：调用方传进来的路径**已经**过了 C26 ⇒ 修法是把"已解析"变成
  **类型上可证明**的（收一个 `risk.Result`/已解析路径的新类型，而不是收 `string`），让 ban 自动消失且**下次不可能退回**。
  ⚠ 这条要给出**每一个**生产调用点的证据：`internal/agent/spill.go`、`internal/memory/artifacts.go`、
  `internal/memory/open.go`、`internal/secret/store.go`（`internal/memory/open.go` **自己就在 allowlist 上**，
  所以它的"已解析"是**声明**而非证明）。
- **(b) 封权限这条路真的能吃进未解析路径** ⇒ 那这是个**独立缺陷**，比 ban 报错严重：
  要么改成走 `risk.Resolve`（注意 `internal/risk` **不** import `internal/winsec`，**无环**，编排者已查），
  要么在入口**拒绝**未解析输入并**响亮失败**。⚠ **绝不许**"为了过门"塞进 allowlist（那是把红变成没人读的清单）。

## AC（1:1，裁决表 `docs/evidence/s1/94-*.md` 由验收方出）

- [ ] **AC#1** 定性：给出上面 (a)/(b) 的**判决 + 每个调用点的 file:line 证据**（"谁把路径递给 winsec、
      递之前解没解过"）。判不准就写"未证"并说明缺哪次测量——**不许两种都写成可能**。
- [ ] **AC#2** 修法落地后 `sh scripts/d22scan.sh` 在**仓外纯净快照**（`git archive <你的 SHA>`，目录带会话后缀）
      **rc=0**，并把逐作用域台账原样贴出来（`ban #6` 的 35 不许变少）。⚠ 快照目录必须带你的会话后缀。
- [ ] **AC#3** 一条**反向用例**：构造"未解析形式"的输入（junction / 8.3 短名 / 尾随点空格 任一种真实形状），
      断言 `PrivateDirAll` **要么拒、要么封到解析后的那一棵**——**不许封到它以为的那棵**。
      并做变异：把修法退回 `filepath.Abs` ⇒ 该用例必须红（锚点=承载行为那一行，同链 grep 证落地，
      还原后 `git diff --quiet` 证干净；**编译失败不算变异**）。
- [ ] **AC#4** 流程洞（上面那个"按包门禁看不见全仓 ban"）：把它变成**每个代理都会做一次的固定动作**——
      在本票的 Rules/门禁段里加一条"收尾前必须跑 `sh scripts/d22scan.sh`"，并**实测它跑得起来**（给时长与 rc；
      慢到不可接受就报数字并说明为什么，**不许以"太慢"为由不加**）。
      ⚠ 只登记不改 `ci.yml`（CI 侧的改法是票 85 的地界，别抢）。
- [ ] **AC#5** 门禁：`gofmt -l` + `gofumpt -l`（**票 89 目前就红在这条上，见下面"附带"**）、
      `go vet ./internal/winsec/` rc=0、`go test -count=2 ./internal/winsec/` rc=0 并逐条点名 SKIP/FAIL。

## 附带（同一枚 push 里的第二个红，与本票同批清）

HEAD 上 `gofmt -l internal/winsec/` **不为空**（`winsec_windows.go` 未格式化）⇒ 推上去 `lint` 的
`gofmt/gofumpt` 那步会红。**这是纯格式化**（票 70 有先例：`style(70,AC#1): gofumpt … 纯格式化、无行为改动`）。
⚠ **但 `internal/winsec/**` 现在是活人的地界**（票 89 在写），所以**谁先交件谁顺手把它带成 0 行**，
本票在它之后接续。编排者不代跑（共树，A34）。

## Rules（本仓固定）

`git commit -q -F - -- <显式路径> <<'MSGEOF'`（引号 heredoc）；提交前 `git diff --cached --name-only`；
禁 `git add -A`/`.`、`--amend`/`reset`/`rebase`/`stash`/`checkout .`（A34）；**不 push**；
不在仓内建 worktree（A38④）；票面 append-only（改行前先读；标题前插段落要把标题重抄，`git diff --numstat` 删除列必须 0）；
四种假绿逐条点名；**收尾前跑一次 `sh scripts/d22scan.sh`**（AC#4 就是这条的来历）；
不碰 `.scratch/wisp/issues/{82,87,88,90,91,92,93}-*.md`；数字不达标就 FAIL 附数字。
15 次工具调用内交回第一枚 checkpoint；接近轮数上限主动收尾留断点。

## Progress log（append-only）

- 2026-09-21 15:1x（编排者）：建票并**压住 push**。我在推之前对 HEAD 做了一次门检
  （`go build ./...` rc=0 ✓，但 D22 rc≠0 + gofmt 1 文件红），命中的正是票 89 的 `winsec.go:126`。
  我**不**自己动 `internal/winsec/`（那是活人的包），也**不**碰 allowlist（只许变短）。
  顺带把票 70-d 留在它 `next=` 里的第三个决策接住：`internal/config/unwired.go:13` 的 `🔒` **已经不红了**
  （我扫过 HEAD：`ban #8` 对 `internal/` 324 文件、`cmd/` 25 文件**零命中**）⇒ 那条**不需要放行**了，
  本票的命中点是 `pathresolver-bypass`，不是 emoji。
  next= 等票 89 交件（它的 Status 已经写 ready-for-review）⇒ 派单，与本票附带那条 gofmt 同批。

- 2026-09-21 16:0x（agent-ticket94，checkpoint 1/AC#1 定性 + 形状选定）：
  **先复现门检**（仓外快照 `/tmp/wisp94gate-94a`，`git archive HEAD`=15c649f，`sh scripts/d22scan.sh`）：
  **rc=1**，命中 `internal/winsec/winsec.go:126: [pathresolver-bypass] filepath.Abs ...`，
  失败的两条正是票面点名的 `TestScannerSelfScanOfRealRepoIsGreen` 与 `TestRealRepoLedgerIsHonest`；
  台账逐作用域与票面读数一致：`bans #1-5 internal/=190`、`cmd/=20`、`ban #6 frontend/=35 text files`、
  `#7 internal/tools/=16`、`#8 design/=16`、`#8 internal/=324`、`#8 cmd/=25`。墙钟 **10.6s**（`real`）。
  ⚠ 附带那条 gofmt 已不需要本票带：`gofmt -l internal/winsec/` **0 行**（票 89 的 `0a3a445` 自行消除了，
  见 AC#5 段的实测）。

  **AC#1 判决：(b)** —— 封权限这条路**真的**能吃进未解析路径，四个生产调用点一个都没有过 C26

  判据不是"我看代码觉得"，而是**每一条链的源头都落到了 `filepath.Join` + 未解析的环境来源**，
  且全仓 `risk.Resolve(` 的非测试调用点（`grep -rn "risk\.Resolve(" --include=*.go`）**没有一个**在数据根这条链上：

  1. `internal/memory/open.go:176` ← `:169` `abs, err := filepath.Abs(dir)`（就是 allowlist 上那条"声明"）
     ← `cmd/wisp/providers.go:172` `memory.Open(io_.dataDir)` ← `:81` `resolveDataDir(buildinfo.EnvString())`
     ← `cmd/wisp/doctor.go:222-245`：来源只有 `os.Executable()` 目录、`os.UserConfigDir()`、`os.TempDir()`、
     `WISP_TEST_DATA_DIR`，全部只过 `filepath.Join`。⇒ **`open.go` 的"已解析"确实只是声明**。
  2. `internal/memory/open.go:188`（artifacts）与 `:503`（backup）← `:183/:182` `filepath.Join(abs, ...)`，同一个未解析根。
  3. `internal/secret/store.go:49` ← `:42` `filepath.Join(dataDir, "secrets")` ←
     `cmd/wisp/providers.go:84` / `cmd/wisp/run.go:179` / `cmd/wisp/secret.go:246`；
     `secret.go:170` 的 `layout.DataDir` ← `internal/proc/envfork.go:55-136`（`filepath.Join(userConfigRoot, "wisp")`
     / `PortableDataDirName` / `ApplyPortableOverride` 的 `filepath.Join(exeDir, name)`）⇒ **同样零解析**。
  4. `internal/agent/spill.go:111` ← `Spiller.dir` ← `internal/agent/loop.go:238` `NewSpiller(opt.Config.ArtifactsDir)`
     ← `cmd/wisp/run.go:324` `filepath.Join(rt.spec.dataDir, "artifacts")` ← `run.go:135` `resolveDataDir(...)`。

  ⇒ 也就是说 `PrivateDirAll` 的"逐级建 + 逐级封"决定的是**一条只做过词法拼接的树**。`filepath.Abs` 不解析
  reparse、不展开 8.3、不剥 `\\?\`、不尾随点空格，而 `WISP_TEST_DATA_DIR` 与 `%APPDATA%` 这两个来源
  都在这台机器上真实可控 ⇒ 这是**独立缺陷**（(b)），不是"该用哪个 API"的口味问题（(a)）。

  **⚠ 推翻本票前提：`winsec` import `risk` 会成环**

  票面（与派单语）写的"编排者已核过 `internal/risk` **不** import `internal/winsec` ⇒ 无循环依赖"
  **只对直接 import 成立**。实测闭包：
  `go list -deps ./internal/risk/` → `internal/risk` ← `internal/observe`（`redact.go:112` 调 `secret.RedactSecret`）
  ← `internal/secret`（`store.go:10` import `winsec`）← `internal/winsec`。
  ⇒ `winsec → risk` 一旦直写，链路 `risk → observe → secret → winsec → risk` 成环，`go build` 直接拒绝。
  **缺的这次测量就是 `go list -deps`**（只核 `.Imports` 不够）。

  **形状选定（AC#2/AC#3 的落点）**

  因为 (i) AC#1=(b) 必须真的解析或响亮拒绝、(ii) 四个生产调用点全在**需要协调**的文件里
  （`internal/memory/*.go`、`internal/secret/store.go`、`internal/agent/spill.go`；`cmd/wisp/`、`internal/config/` 禁碰），
  (iii) `risk → winsec` 方向是**干净**的（winsec 无内部依赖）、
  修法收在**本票地界 + 一枚反向注入的安装钩子**：
  - `internal/winsec/`：新增 `ResolvedPath` 只经 C26 管道铸造（seam 未装 ⇒ **响亮失败**，不是静默词法拼接），
    `PrivateDirAll`/`SealDir`/`SealFile`/`PrivateFile*` 的 `filepath.Abs` 位置改吃**已解析**值；
    **公开签名不收窄**（不收 `string`→换类型），以免打断票 90/18 正在写的调用方编译。
  - `internal/risk/`：**新增一个文件**（不碰 `pathresolver*.go`/`assessor.go`/`rules_gateway.go` 三个冻结件）
    把 `risk.Resolve` 装进 winsec 的 seam ⇒ C26 仍是**唯一实现**，winsec 不复制第二条管道。
  - AC#3 用例：真 junction（`mklink /J`，本包已有 `mkJunction` 且票 89 证过无特权可用）喂进 `PrivateDirAll`，
    断言"要么 `ErrReparseDenied` 拒、要么封到解析后那一棵"，**不许封到它以为的那棵**；
    变异 = 把入口退回 `filepath.Abs` ⇒ 该用例必须红（编译必须过，否则不算变异）。
  next= 落地上述形状 → 写 AC#3 用例 → 跑变异 → AC#4 门禁登记 → AC#5 全量读数 → 纯净快照 d22scan rc=0。
