# ticket-181-status-1 · 线索复跑 + AC#7 现状核对（2/2 件）

起手 HEAD＝`7a452d0aca1bf5157a97fb5329ca6eb38b9f0694`（10-08 16:10）；现量时刻 `date`＝`2026-10-08 16:15 +0800`。
锚数据与票面四把尺在第一件 `anchor.md`。本程⛔ 零 Go 命令、⛔ 零突变、⛔ 零翻框（`253-r1` 正在 `cmd/wisp` 上跑）。

## §1 五条线索逐条复跑（命令原文 + 输出摘要；裸文件名一律先定包再量行数）

包/行数先定（防本仓 `bridge.go` 两枚同名的老坑）：
`wc -l internal/panel/workspace.go internal/tools/paths_workspace.go cmd/wisp/panel_pump.go internal/panel/instructions_200.go`
⇒ `182` / `160` / `420` / `215`（rc=0）。四枚都是不同包、无重名歧义。

| # | 派单转述的线索 | 我跑的命令 | 读数 | 判定 |
|---|---|---|---|---|
| 1 | `internal/panel/workspace.go:72` 注释逐字 `This is the producer side of WorkspaceView.Rewritten. Before this ticket the` | `grep -n "producer side of WorkspaceView.Rewritten" internal/panel/workspace.go` | `72:` 该行逐字命中（rc=0） | **成立** |
| 2 | 同文件 `:91`＝`Reparse: res.Reparse, Rewritten: res.Rewritten,`；`:94` 起是 `if res.Rewritten {` | `grep -n "Rewritten" internal/panel/workspace.go` | `91:` 与 `94:` 逐字命中；另 `:141`／`:144` 是 `RequestWorkspaceSwitch` 里同形的第二处 | **成立** |
| 3 | `internal/tools/paths_workspace.go:48` 具名 `WHY THIS RETURNS AN ACCOUNT AND NOT A STRING. WorkspaceView.Rewritten is read`；`:61-63` 说 `Rewritten == true` 从此可达 | `grep -n "WHY THIS RETURNS AN ACCOUNT\|Rewritten" internal/tools/paths_workspace.go` | `48:`／`61:`／`63:` 命中；`WorkspaceRoot()` 签名现量＝`func (p *PathCanonicalizer) WorkspaceRoot() risk.Result`（`:67`） | **成立**（但"可达"那句只到**函数层**，见 §4） |
| 4 | 生产调用者＝`cmd/wisp/panel_pump.go:87` 一行 `return panel.WorkspaceViewFromRoot(rt.paths.WorkspaceRoot())` | `grep -n "WorkspaceViewFromRoot" cmd/wisp/panel_pump.go` ＋ `git grep -n "WorkspaceViewFromRoot(" HEAD -- '*.go' \| grep -v _test.go` | `:87` 逐字命中 ⇒ **成立但不全**：HEAD 上非测试调用点共 **两枚**＝`cmd/wisp/panel_pump.go:87` 与 `internal/panel/workspace.go:165`；第三处命中是 `internal/panel/git.go:262` 的**注释**，第四处是他人证据件 `.scratch/wisp/probes/197/r3b/pre/panel_pump.go:87`（非产码） | **成立＋漏计一枚**（按调用形状锚，未用裸符号名） |
| 5 | `internal/panel/instructions_200.go:109-116` 具名 `Rewritten=true IS reachable through the snapshot path`，并点名测试 `TestRewrittenWorkspaceViewYieldsNoDirectoryAtAll` | `grep -n "reachable through the snapshot path" …` ＋ `git grep -n TestRewrittenWorkspaceViewYieldsNoDirectoryAtAll HEAD` ＋ worktree `grep -rn --include='*.go'` | `:109` 逐字命中；注释块实为 `:107-118`（含"200-r2 那句诚实说明已过期"的自我打旧） | **成立** |

## §2 本程最重那一格：被点名的测试**存在**，不是假凭据

- `git grep -n TestRewrittenWorkspaceViewYieldsNoDirectoryAtAll HEAD -- '*.go'` ⇒ 两类命中：
  `HEAD:internal/panel/instructions_200.go:116`（引用它的注释）＋
  **`HEAD:internal/panel/instructions_200_test.go:60:func TestRewrittenWorkspaceViewYieldsNoDirectoryAtAll(t *testing.T) {`**（**定义在盘上**，rc=0）。
- 跨全文件类型再拉一遍（`git grep -n … HEAD`，rc=0）⇒ 除产码/测试两枚外，另有 30+ 枚**执行痕迹**：
  `docs/evidence/s1/200-project-instructions-r2.md:91/:190/:227`（票 200-r2 写的用例，且 §M1 那格记的是它被**突变打死过**的红句原文）、
  `.scratch/wisp/probes/181/r3/evidence.md:109`、多枚票（197/211/220/222/223/253/255/35）门禁日志里的 `--- PASS`、
  `.scratch/ci-logs/run-37166458550-failed.log:7808-7809`（10-04 那发 CI `--- PASS (0.00s)`）。
- ⇒ 判定：**〔引用真凭据〕**，不是〔引用假凭据〕。⚠ 一条纪律：CI 那枚 PASS 取在 **10-04**，
  **早于**落地枚 `16901acb`（10-05 11:36）⇒ 它只能证明用例存在且当时绿，⛔ 不能当"现状绿"的凭据（现状要跑，本程禁跑）。

## §3 是谁落的（具名 hash，⛔ 不按"最近改过这文件"认）

- `git log --oneline -- internal/panel/workspace.go` ⇒ 三条：`16901acb`／`91b5fc4e`（票 92）／`8e100950`（票 92）。
- `git log -S 'producer side of WorkspaceView.Rewritten' --format='%h %ad %s' --date=format:'%m-%d %H:%M' -- internal/panel/workspace.go internal/tools/paths_workspace.go`
  ⇒ **只一条**：`16901acb 10-05 11:36 票 181 AC#7 生产侧落地（181-r3 第 2 笔，产码＋判据）：改写账户的生产者从"留零值"换成"C26 的读数"`。
- `git log -S 'WHY THIS RETURNS AN ACCOUNT AND NOT A STRING' -- … internal/tools/paths_workspace.go` ⇒ 同一枚 `16901acb`。
- `git show -s --format='%H%n%an <%ae>%n%ad%n%s' 16901acb` ⇒ 全 hash `16901acba63ad2d76d6f99099e254884fb889113`，作者 `CarlosShao <1933942520@qq.com>`（本机共用身份，⛔ 不据此推哪条腿），`2026-10-05 11:36 +0800`。
- 改动面（`git show --stat --format='' 16901acb`，11 files／+653 −62）：
  `internal/panel/workspace.go +90`、`internal/tools/paths_workspace.go +62`、`internal/panel/instructions_200.go +22`、
  `internal/panel/git.go +24`、`internal/tools/paths.go +12`、新判据件 `internal/panel/workspace_account_181r3_test.go +251`、
  `internal/tools/paths_workspace_account_181r3_test.go +201`，另三枚既有测试件同步改写。
- **★与派单转述的冲突（具名报回）**：转述猜"很可能已经被**别的票**落地" ⇒ **不成立**。
  它属**本票自己的 AC#7 生产侧腿 `181-r3`**，四笔 commit 逐枚现验（`git log -1 --format` 各 rc=0）＝
  `b672f853`（10-05 11:24 证据件骨架）→ `16901acb`（11:36 产码＋判据）→ `9edb33e8`（11:43 证据件 §3-§10）→ `ce877685`（11:45 终态门禁，⛔ 不动产码）。
  证据件＝`.scratch/wisp/probes/181/r3/evidence.md`（286 行）。
- **另一处冲突（行号漂移）**：票面 `AC#7` 那行引的 `workspace.go:60-68`／`:85`／`:104-106` 与 `paths_workspace.go:64-66`
  在 HEAD 上已不对应（现量：构造者＝`:85-103`、`RequestWorkspaceSwitch` 成功分支的拷贝＝`:141`／判定＝`:144`、
  `paths_workspace.go` 现 160 行且 `:135` 是 `p.workspaceAccount = acc`）。⛔ 我不回改票面那些行号，只在追加节具名"按形状读"。

## §4 为什么这一格**仍不能翻勾**（缺的是哪把尺，具名）

1. **AC#7 的原判据要的是行为突变凭据，而盘上那枚凭据出自实现者之手**：`181/r3/evidence.md` §5（`:152-167`、`:284`）
   具名跑了 M1（把生产者退回"不填"）／M2（把渲染退回三枚字段），红句逐字 `the producer dropped the book again`，
   还原凭据＝仓外 `sha256sum -c` ＋ porcelain。⇒ 突变**跑过、但跑的人是写代码的那一枚**。
   `AGENTS §0` 第 3 句／`SPEC-12 §4.3` #1/#3 ⇒ 缺口审计与对抗验收必须**另一个** agent。
2. **编排者已经裁过同样的结论，本程只是复认**：台账 `docs/reports/pending-and-issues.md:12103`（`A618`，10-05 11:5x）逐字
   ⇒「**票 181 AC#7 此刻不翻勾**，缺的是 `181-v3`（非实现者终裁，含上面①那一格的真跑）」。
   该腿还具名交了五格判不动，其中与 AC#7 直接相关的两格：①`cmd/wisp/panel_pump_test.go` 那枚断言**只 `t.Log` 不判红**（归能跑 `cmd` 的那一枚）；
   ②面板发起的收窄那条路 `rewritten` **仍不可达**。`:12137`／`:12156`／`:12212` 三处记 `181-v3` **继续按住**（带突变的验收腿与 `cmd/wisp` 编译闭包互洗读数）。
3. **本程新取的两把尺，说明"可达"只到函数层，没到真进程**（⛔ 我未跑任何 Go，全部静态现量）：
   - `git grep -n "SetWorkspaceRoot(" HEAD -- '*.go' \| grep -v _test.go` ⇒ 生产侧唯一记账点＝`internal/panel/workspace.go:122`，
     它在 `RequestWorkspaceSwitch` 函数体内，而 `:122` 之前 `res.Actable()` 先判：`internal/risk/pathresolver.go` 的
     `func (r Result) Actable()`（`:94` 起）里 `if r.Rewritten { return "", fmt.Errorf("%w: …") }` ⇒ 任何改写仍先被拒。
     `git show --stat --format='' 16901acb -- internal/risk/` ⇒ **空输出且 rc=0** ⇒ 那两道闸一字未松（与票面原判据同向）。
   - `git grep -n "RequestWorkspaceSwitch" HEAD -- '*.go'`（按调用形状核）⇒ 非测试命中只有注释与定义；
     真调用点**全在 `_test.go`**（`cmd/wisp/instructions_200r2_test.go:124`、`cmd/wisp/panel_pump_test.go:357`、`internal/panel/workspace*_test.go`）。
     `internal/panel/composer_dispatch.go:73-77` 逐字：workspace 那枚 handler「the handler in front of it is **ticket 186's work**,
     so this file declares the socket and nothing more」⇒ 生产里今天没人调它。
   - ⇒ 综合：**字段不再是"结构上恒 false"（生产者函数已把账户拷贝进去、且有专测从账户一路驱动到消费点），
     但真进程里今天没有任何一条路把它记成 `true`**——记账点被 `Actable()` 挡着、且记账点所在函数零生产调用者（面板 handler 属票 186）。
     这恰好是 AC#7 那句"判据要能区分『填了账户』与『字段恒 false』"还没被**非实现者**结掉的原因。
4. **落点判词（本程结论，⛔ 不是"已做完"）**：
   - 已落地＝票面 `AC#7` 要做的"把 producer 侧填成真值"那半件**在盘上**，具名 `16901acb`（腿 `181-r3`，本票自己的腿），
     文件行现量＝`internal/panel/workspace.go:85-103`（含 `:91` 拷贝与 `:94-98` 文案分支）＋`internal/tools/paths_workspace.go:43-75`。
   - 仍欠＝**`181-v3`：非实现者的行为突变终裁**（含 M1/M2 的复跑＋`cmd/wisp` 那枚"只 `t.Log` 不判红"的断言真跑）。
     归属＝**编排者派给非实现者腿**（台账 `A618` 已点名腿号 `181-v3`）；时机＝`cmd/wisp` 写腿（今日在飞＝`253-r1`）退出、
     且水位与突变互斥条件满足之后；⛔ 本程一枚 Go 命令都没跑，⛔ 不代它出凭据。
   - 另一条并列具名欠账（票面既有、非本程新造）：真进程可达性要等**票 186 的 workspace handler**，
     或等 `internal/risk` 的判级语义被人工批准改动——后者是禁区，⛔ 不许任何腿"顺手"动。

## §5 票面追加自证（追加节写完后现量，见第二笔 commit）

追加方式＝只在票面真末尾挂新一节，⛔ 不回改任何已提交行；旧句以**留形写法**逐字嵌进方括号，一字不删。
四把尺（追加前后）：行数 `77 → 待量`、未勾框 `2 → 2`、已勾框 `5 → 5`、
`git diff --numstat 7a452d0aca1bf5157a97fb5329ca6eb38b9f0694 HEAD -- <T>` 只出现票 181 那一枚文件（＋本腿两枚新 `.md`）。
