# 94 — 私有数据目录的"封 ACL"用了 `filepath.Abs` 而不是 C26 PathResolver ⇒ **D22 门在 CI 之前把它拦下了**（票 89 的码，push 因此压住）

**Status:** ready-for-review（agent-ticket94：AC#1-AC#5 五框自勾，代码 `7910bcd`+`20b525d`，
仓外纯净快照 `sh scripts/d22scan.sh` **rc=0** ⇒ **编排者的 push 可以放开了**；⚠ 留给验收方的两点写在 Progress log 末条）
**原 Status:** in-progress（agent-ticket94 已接手；票 89 已交件 `0a3a445`，`internal/winsec/` 现在是本代理地界）
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

- [x] **AC#1** 定性：给出上面 (a)/(b) 的**判决 + 每个调用点的 file:line 证据**（"谁把路径递给 winsec、
      递之前解没解过"）。判不准就写"未证"并说明缺哪次测量——**不许两种都写成可能**。
      ⇒ **判决 (b)**，四条链的 file:line 与"推翻无环前提"的测量都写在 Progress log 末段（checkpoint 1）。
- [x] **AC#2** 修法落地后 `sh scripts/d22scan.sh` 在**仓外纯净快照**（`git archive <你的 SHA>`，目录带会话后缀）
      **rc=0**，并把逐作用域台账原样贴出来（`ban #6` 的 35 不许变少）。⚠ 快照目录必须带你的会话后缀。
      ⇒ `/tmp/wisp94gate2-94` = `git archive 20b525d`，**rc=0**；台账见下面"门禁读数"段（`ban #6` 35→**37**，只升未降）。
- [x] **AC#3** 一条**反向用例**：构造"未解析形式"的输入（junction / 8.3 短名 / 尾随点空格 任一种真实形状），
      断言 `PrivateDirAll` **要么拒、要么封到解析后的那一棵**——**不许封到它以为的那棵**。
      并做变异：把修法退回 `filepath.Abs` ⇒ 该用例必须红（锚点=承载行为那一行，同链 grep 证落地，
      还原后 `git diff --quiet` 证干净；**编译失败不算变异**）。
      ⇒ `internal/winsec/resolve_windows_test.go`：真 junction（`mklink /J`）三条子用例，变异红名与 SID 级证据在下面。
- [x] **AC#4** 流程洞（上面那个"按包门禁看不见全仓 ban"）：把它变成**每个代理都会做一次的固定动作**——
      在本票的 Rules/门禁段里加一条"收尾前必须跑 `sh scripts/d22scan.sh`"，并**实测它跑得起来**（给时长与 rc；
      慢到不可接受就报数字并说明为什么，**不许以"太慢"为由不加**）。
      ⚠ 只登记不改 `ci.yml`（CI 侧的改法是票 85 的地界，别抢）。
      ⇒ 固定动作已写进本票"## 门禁（AC#4）"段（`## Rules` 之前）；时长/ rc 在下面读数里。
- [x] **AC#5** 门禁：`gofmt -l` + `gofumpt -l`（**票 89 目前就红在这条上，见下面"附带"**）、
      `go vet ./internal/winsec/` rc=0、`go test -count=2 ./internal/winsec/` rc=0 并逐条点名 SKIP/FAIL。
      ⇒ 五项全绿，逐条数字在下面"门禁读数"段；"附带"那条 gofmt 已被票 89 的 `0a3a445` 自行带成 0 行（本票未动）。

## 附带（同一枚 push 里的第二个红，与本票同批清）

HEAD 上 `gofmt -l internal/winsec/` **不为空**（`winsec_windows.go` 未格式化）⇒ 推上去 `lint` 的
`gofmt/gofumpt` 那步会红。**这是纯格式化**（票 70 有先例：`style(70,AC#1): gofumpt … 纯格式化、无行为改动`）。
⚠ **但 `internal/winsec/**` 现在是活人的地界**（票 89 在写），所以**谁先交件谁顺手把它带成 0 行**，
本票在它之后接续。编排者不代跑（共树，A34）。

## 门禁（AC#4：按包门禁结构性看不见全仓 ban ⇒ 固定动作）

本票的教训不是 `winsec` 写错了 API，而是**票 89 把按包门禁全跑绿之后 HEAD 仍然是红的**：
`bans #1-5` 的作用域是整个 `internal/` 与 `cmd/`，而代理的门禁清单只有 `gofmt/gofumpt/vet <pkgs>/test <pkgs>`。
⇒ 从今天起本仓每个动 `internal/**` 或 `cmd/**` 的票，**收尾前必须跑一次**（写进票面门禁段，不跑不许勾 AC#5）：

```
rm -rf /tmp/wisp<票号>gate-<会话> && mkdir -p /tmp/wisp<票号>gate-<会话>
git archive HEAD | tar -x -C /tmp/wisp<票号>gate-<会话>
cd /tmp/wisp<票号>gate-<会话> && sh scripts/d22scan.sh   # 期望 rc=0，并把台账逐作用域贴进票面
```

**实测跑得起来**（本票两次读数）：
- 基线（HEAD=`15c649f`，独占）：**rc=1**，墙钟 **10.6s**。
- 交件（HEAD=`20b525d`，与并发的 `go test -count=2 ./internal/winsec/` 同跑，故偏慢）：**rc=0**，墙钟 **43.7s**；
  其中 `tools/d22scan` 自扫 `go test ./...` 11.3s。慢的唯一原因是并发编译 + `icacls` 争 CPU，**不是脚本本身**，
  所以这条动作按 10~45s 计价可接受。
⚠ CI 侧要不要常驻这一步是**票 85 的地界**，本票只登记不改 `.github/workflows/ci.yml`。

## 门禁读数（AC#2/AC#5 的原始数字，快照 `/tmp/wisp94gate2-94` = `git archive 20b525d`）

```
sh scripts/d22scan.sh                                   → rc=0，d22scan: clean - no D22 ban violations
d22scan: scope bans #1-5 internal/      examined 197 production Go files
d22scan: scope bans #1-5 cmd/           examined  20 production Go files
d22scan: scope ban #6 frontend/         examined  37 text files          （基线 35 ⇒ 只升未降 ✓）
d22scan: scope ban #7 internal/tools/   examined  17 production Go files
d22scan: scope ban #8 design/           examined  16 text files
d22scan: scope ban #8 frontend/         examined  37 text files
d22scan: scope ban #8 internal/         examined 335 Go files（含注释与 _test.go）
d22scan: scope ban #8 cmd/              examined  25 Go files（含注释与 _test.go）
```
`gofmt -l internal/winsec/ internal/risk/` → 0 行；`gofumpt -l` 同 scope → 0 行；
`gofmt -l ./internal ./cmd ./tools`（全仓，防我把别处弄红）→ 0 行；`go vet ./internal/winsec/` → rc=0。
`go test -count=2 -v ./internal/winsec/` → **rc=0**，`=== RUN` 44 行 = 2 × 22 个不同名（12 个顶层 + 10 个子测试）✓，
`--- SKIP` **0 条**、`--- FAIL` **0 条**（不是"没跑"：44 次 RUN 全有对应 PASS）。
连带回归（HEAD 基线同样全绿，逐包 rc=0）：`./internal/secret/ ./internal/memory/ ./internal/agent/ ./internal/risk/ ./internal/observe/`。
⚠ 这条回归本身就是形状的一部分：把未接线的默认做成"响亮拒绝"会让 `internal/secret`/`internal/memory`
两套 suite 红 **20+ 条**（它们不 link `internal/risk`），实测读数就是内置 verifier 存在的理由。

**收尾全仓读数**（`/tmp/wisp94final-94` = `git archive 73cec78`，含本票三枚 commit）：
`sh scripts/d22scan.sh` **rc=0**（墙钟 26.9s，台账与上表逐字相同：`bans #1-5 internal/=197 cmd/=20`、
`ban #6 frontend/=37`、`#7 internal/tools/=17`、`#8 design/=16 frontend/=37 internal/=335 cmd/=25`）。
`go test ./...` **rc=1，两条红，都不是本票引入的**（逐条归因）：
① `cmd/wisp` `exit status 0xc0000135`（DLL 找不到）—— 同一枚红在**基线快照** `/tmp/wisp94base-94b`
（HEAD=`15c649f`，本票动码之前）与当前工作树里**一模一样复现** ⇒ 环境/他票（`cmd/wisp/` 禁碰清单内）；
② `internal/risk` `--- FAIL: TestResolvePerCallBudget`（30.291 ms/op vs 1 ms 预算）—— 那次读数是在
与 `go test -count=2 -v ./internal/winsec/`（icacls/mklink 密集）+ 全树编译并发时取的；
静默后 `go test -count=1 -run TestResolvePerCallBudget ./internal/risk/` **三次全 ok**（1.504s / 2.441s / 1.863s，多样本全报）。
⇒ 这条是**仪器坑**（预算测试对 CPU 争用敏感），不是回归；本票没动 `pathresolver*.go` 一个字。

**AC#3 变异（退回 `filepath.Abs`）的红名与证据**：锚点 = `internal/winsec/winsec.go` 的
`dir, err := ResolvePath(path)` 那四行；同一条 `&&` 链里先 `grep -n "MUTATION-94"` 证落地、
`go vet ./internal/winsec/` rc=0 证**不是编译失败**，然后：

```
--- FAIL: TestAC3JunctionInputIsRefusedNotSealed/existing_directory_behind_the_link
--- FAIL: TestAC3JunctionInputIsRefusedNotSealed/missing_directory_under_the_link
--- FAIL: TestAC3JunctionInputIsRefusedNotSealed/with_only_the_built-in_verifier
PrivateDirAll returned success for a path whose tree lives behind a junction: …\data\link\artifacts
WRONG TREE SEALED: …\someone-elses-tree\artifacts
  went from [S-1-1-0(Everyone) S-1-5-32-544 S-1-5-18 S-1-5-21-…-1001]
        to  [S-1-5-18 ×2 S-1-5-32-544 ×2 S-1-5-21-…-1001 ×2]     ← 外来读授权被我们的封存抹掉
```
⇒ 旧写法**确实**把别人的 DACL 重写了（不是"可能"，是 icacls 读数）。还原后
`git diff --quiet -- internal/winsec/winsec.go` rc=0（共树，全局 diff 不可能干净，故按路径 scope），
重跑 `go test -run TestAC3JunctionInputIsRefusedNotSealed ./internal/winsec/` → ok。
内置 verifier 那一腿（`SetPathResolver(nil)`）在同一条用例里独立断 `ErrUnresolvedPath`，
所以"接线被删掉"也仍然红；`TestC26PipelineIsWiredIntoWinsec` 单独钉住 `risk` 的 `init` 装上了真管道。

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

- 2026-09-21 17:1x（agent-ticket94，交件）：五框自勾，代码两枚 commit：`7910bcd`（形状落地）+ `20b525d`
  （变异复现要打印"封到了哪一棵"）。改动清单：新增 `internal/winsec/resolve.go`（`C26Resolver` seam +
  `ResolvedPath` 类型 + 内置 verifier）、`internal/winsec/placement_windows.go`（Windows 只验不写的放置检查）、
  `internal/winsec/resolve_windows_test.go`（AC#3 + 接线 pin）、`internal/risk/winsec_c26.go`（`init()` 把
  `Resolve(input, nil)` 装进 seam）；修改 `internal/winsec/winsec.go`（五个封存入口动手前先解析，
  `PrivateDirAll` 的 `filepath.Abs` 消失，封存核心改收 `ResolvedPath`）、`winsec_other.go`（POSIX 腿写明与
  risk 的 DEFERRED 对齐）。**allowlist 一行未动（仍 5 行非注释）**、`tools/d22scan/**` 未动、
  `pathresolver*.go`/`assessor.go`/`rules_gateway.go`/`docs/PLAN.md`/`docs/specs/*.md` 一字未改、
  `frontend/**`/`internal/panel/`/`cmd/wisp/`/`internal/config/`/`internal/agent/` 未碰、`ci.yml` 未碰。
  需要谁协调：**目前为零** —— 我刻意保住 `PrivateDirAll(string)` 的签名，`internal/memory/open.go:176,188,503`、
  `internal/secret/store.go:49`、`internal/agent/spill.go:111` 三个人的文件一行没改就能编过，
  它们的输入现在由 winsec 自己解析/拒绝（读侧 C26 的接线仍是票 18 的活，`internal/memory/open.go`
  那条 allowlist 也因此**没有**变短：本票不消费它）。
  留给验收方的两点（我自认的弱处，别当成已证）：
  ① 内置 verifier 只在 `internal/risk` 没被 link 时才是主角，生产二进制永远走 C26 ⇒ 那条腿的 CI 覆盖
     靠 `resolve_windows_test.go` 里 `SetPathResolver(nil)` 那一腿撑着；
  ② `RemoveUnlinked` 明确不解析（理由写在函数注释），但它因此仍能沿着一个 junction 走到"它以为在自己的
     artifacts 树里"的位置去 unlink —— 那棵树是 `internal/memory/artifacts.go:171,181` 递进来的，
     要收这个口得先动 memory（票 18/79 的地界），本票没动。
  票面 numstat 提醒：本次 72 加 6 删 = 5 条 AC 勾框 + 1 条 Status 转换，段落插入部分删除列为 0（标题已重抄）。
  next= 交给验收方出 `docs/evidence/s1/94-*.md`；编排者可放开 push（HEAD 的 D22 已 rc=0）。
  若验收方要我把 seam 换成"调用方递已解析类型"的强形状，那需要票 18 或 79 在 memory/secret/agent 三处
  同步改签名，本票单独做会让别人编译失败。
