# 派单 175-r2（**写手腿·一格一单**）＝把 `task.output` 送进 C25 那八枚名册，并**第一次在真桥上**量出"盖了戳之外，宿主自己在桩里写下的那条路径仍要被豁免、而外来说到同一条路径仍要命中 R4"——⚠ 票 177 的豁免今天**是惰性的**，这一枚是它**第一次真生效**

- 派单时刻：09-28 00:2x｜锚点＝**你自己 step 0 现跑那三行**（**不许手敲我的号**，我今天手敲错过一次）。
- 票面＝`.scratch/wisp/issues/175-…`（`ls .scratch/wisp/issues | grep '^175-'` 取真名；**AC#1/AC#2/AC#3 已勾**，你要做的是 AC#3 那句"判据要能防下一程自动通"的落地半，＋AC#4 的自证）。邻居票面：`177-…`（读 22:5x 与 23:4x 两段追加——**豁免已落地但惰性**）、`176-…`（读 00:1x 那段——**起跑口仍不存在，别顺手做**）。
- ⚠ **预算硬顶 ≤45 次工具调用；到顶即停手回禀并说第几次超**。本单**只做上面那一格**：不修票 174 的 AC#2b／AC#2c／AC#2d，不做票 176 的起跑口，不动 `cmd/**`。

## 0. 起手四件（缺一件停手）

`date`｜`git rev-parse --abbrev-ref HEAD`（**必须 `dev`**）｜`git rev-parse HEAD` ＋ `git log --oneline -3`（逐字贴进你的表）｜`git status --porcelain -- internal/ cmd/`（**必须空**）。⚠ **别用 `&&` 串可能无命中的 `grep`**（链一断后面全没跑，我今天这样写过一枚假读数进台账＝`A356`／作废见 `A357`）；串命令一律 `;` 分隔或 `| cat`。

## 1. 基线里有一枚**争用红**，你要认识它、不要修它

起手请跑**两遍**并都留档：`go test -count=1 -run TestResolvePerCallBudget ./internal/risk/`。
已知读数：`23:51`／`23:54` 两枚程各见 `1.098`／`1.199 ms/call` 对 `1ms` 预算；**编队空后（00:0x 我自己跑）＝`--- PASS (1.20s)`**。⇒ **若你起手遇到它红：复跑一次并在表里写"争用/复跑绿"或"复跑仍红"**；**复跑仍红就停手上报**，⚠ **`thresholds.go` 与那枚判据一字节都不许动**（AGENTS §1.1），也不许为了绿去调并发或加 sleep。

## 2. 落地（两枚文件，分两枚提交）

1. **`internal/risk/provenance.go`**——把 `task.output` 送进 `sensitiveSourceTools`（现读 `:96-101`，八枚；行号自己复算）。⚠ 来路已裁＝那张名册是**举例定义**（票 175 AC#1 已勾、`175-c1` 裁的），所以**这是实现缺陷修复、不是契约变更、不需要 owner 批准**；但**只加这一枚名字**，别顺手调顺序、别改注释里那句"穷举/举例"的措辞。
2. **`internal/tools/bridge.go`**——让 `task.output` 那条**宿主自己写下的路径**（票 177 落地的 `hostPathBox` 通道，现读 `:457`／`:579-635` 一带）**真的送达** `MarkWithHostPath`。⚠ **177 只把载具修好了、没人往里放东西**：`177-r1` 交件逐字说"今日生产行为逐字不变（`task.output` 未进名册）"。⇒ 你这一枚是**第一次让它生效**，所以**判据必须落在桥级，不能只落在 risk 级**。

## 3. 三发常驻判据（本单真正的货，落 `internal/tools/`）

| 编号 | 要钉什么 | 未修码（起手锚）上响不响——**两向都要贴读数** |
|---|---|---|
| **J-1 正向** | `task.output` 的成功结果**必须有来源标记**（票 175 AC#3 的 canary 本体在 `.scratch/wisp/probes/175/r1/canary-…ticket175_test.go.txt`，9301 字节，可直接搬回） | **今天不响＝它正是缺失本身**（`task.output` 不在名册） |
| **J-2 反向（最值钱）** | **外来内容里出现同一条 artifacts 路径、模型随后去 `fs.read` 它 ⇒ 仍然命中 R4**（票 177 AC#2 那一发的**桥级版**；夹具形状见 `.scratch/wisp/probes/177/c1/fixture-reverse-criterion.md`） | 未修码上**不可能响**（今天压根不盖戳），落地后必须**会响**——⚠ 任何修法让它变安静＝**洗戳**，本单判不通过 |
| **J-3 豁免真生效** | 宿主自己在桩里写下的那条路径**续读不再被 R4 拦**：⚠ **这就是票 164 已勾的那两枚续读腿**（`TestLongOutputPointerRecoversEveryByte`／`TestPointerPast256KiBIsNotFullyReadable`，`internal/tools/task_output_leg_test.go`）——**你加完名册之后它们必须仍然绿**；若红＝177 的豁免在真桥上覆盖不住，**停手回禀、不许改那两枚判据**（它们已被非实现者裁过） | 今天绿（因为不盖戳）；落地后**必须仍绿**＝豁免生效的证明 |

⚠ **每一枚都要自证非恒真**：给出"把它防的那种做宽法做进去→它红"的变异读数（先 `grep -n` 证落地，再还原）。**一枚不会红的判据＝装饰**。
⚠ **别拿 `provenance_test.go:136`／`:139` 当哨**：`177-m1` 已现量它们对"路径豁免宽过头"**天生不敏感**（值不含 `/`）；那两枚仍不许碰，但**不是你的警报器**。真正的警报器是 J-2。

## 4. 两条禁令 + 两枚陷阱（逐条自证，别留给下一位）

- **禁区 (i)**：不给 `internal/agent.Loop` 加收"收 `taskID` 的导出方法"。⚠ **G3 安静 ≠ 合规**——那条腿按**字面 `taskID`** 抓，改名或包进 struct 就纹丝不动（`176-a1b` 现量标红的一条），**不许拿它的 rc=0 当通行证**。
- **禁区 (ii)**：不把宿主内部 artifacts 写入做成受门控的 Tool（`AGENTS §1.2` 绝对禁止项；`spill.go:29-31` 逐字 "deliberately not a gated tool"）。若你觉得"非做不可才能让 C25 认下后台产物"——**停手上报**。
- **陷阱 G2**：`gate-clauses.sh` 的 **G2 是 `want ring / want_n 2`，而其中 1 行是注释**（`cmd/wisp/run.go:564` 提到 `CloseTask`）⇒ **含 `OpenTask|CloseTask` 字面的新说明会把那条腿顶红**。本单**不动 `cmd/**`**，若你新增注释，避开这两个字面。
- **票 175 AC#4**：不许把 `task.output` 从 `L0` 改档位来"解决"（那既动 `PLAN.md:2564`＝契约、又不补上戳），不许新增配置项／目录豁免／`allowlist.txt` 条目。

## 5. 门禁（自己跑，禁全仓 `./...`）

`go test -count=1 ./internal/risk/ ./internal/tools/`｜`sh scripts/d22scan.sh`（应 rc=0）｜`sh .scratch/wisp/probes/154/gate-clauses.sh`（**应 rc=0**，逐字贴"腿数＝14 声明与实测不符＝…"；⚠ 过程若 G5/G2 因你响，**改你的腿、不许动仪器**——`177-r1` 有先例见 `e741d1cf`）｜`bash tools/d22scan/runtests.sh -C tools/d22scan ./...`（名册两向 `comm -3`：**只许多、不许少**）｜`gofumpt --version` 现跑＋对你名下 `.go` 的 `-l`。⚠ **别跑 `probes/161/r6/flip-declaration.sh`**。⚠ `internal/panel/**` 一行别碰（那枚常红属别家）。

## 6. 写面（超出即越权）

`internal/risk/provenance.go`（只改名册那一处）、`internal/tools/bridge.go`、新增判据文件 `internal/tools/*_test.go`（**新建自己的文件，别追加进 `task_output_leg_test.go`／`provenance_test.go`**）｜`docs/evidence/s1/175-task-output-stamped-r2.md`（你的表）｜`.scratch/wisp/probes/175/r2/**`（⚠ `logdir` 取脚本自身目录）｜票面 175 的 Progress log（**只追加、一枚框都不勾**）。
**禁改**：`cmd/**`（本单不放开）、`internal/agent/**`、`internal/tools/tool.go`（C1 面）、`internal/risk/taintmatch.go` 与 `provenance.go:468-473` 那两段**文字**、`internal/risk/thresholds.go`、`docs/PLAN.md`（含 `:1532`）、`docs/specs/**`、`docs/reports/**`、`frontend/**`、`design/**`、别人的票面与证据件、`probes/**` 既有台件（含 `probes/177/**`，只读）。⚠ 现场躺着 `design/**` 未提交删除、`.gitignore`／`probes/152/my152.py`／`probes/161/r6/logs/flip-*`（` M`）⇒ **不提交、不还原、不补完、不评论**。
**`git add -A`／`git add .`／`git commit -a` 一律禁**；一步式 `git commit -q -F - -- <显式路径>`（⚠ **新增文件要在同一条命令里先 `git add <显式路径>` 再 commit**，我今天这样翻车过；**产码与判据件分枚提交**）；heredoc 必须加引号（`<<'EOF'`）；**含引号或反引号的中文段走 Edit 工具**（我今天这样翻车过）。**只 commit，绝不 push。**⚠ 更正只许追加，禁 `--amend`／`reset`／`rebase`／`stash`／`checkout .`／`restore`／`clean`／`worktree`；零删除命令；绝不在仓内建 worktree 或第二枚 `.git`。

## 7. 交件报告（顺序固定）

1. 起手四件＋基线那枚 budget 红的两向读数（复跑绿／仍红）；2. **本程没测什么**；3. 两处落点逐枚（文件／行号／为什么）；4. **J-1／J-2／J-3 逐枚**：钉什么＋未修码响不响＋落地后响不响＋**我把它防的做宽法做进去之后的红读数**；5. **票 164 那两枚续读腿在加完名册后的颜色**（这是豁免真生效的正面凭据）；6. 两条禁令逐条自证（含"G3 安静不等于合规"那一句你怎么自证）；7. 门禁五枚＋名册差集＋`gofumpt`；8. 被拒／没成功的调用（取数前还是后）；9. 有没有跑过删除命令；10. **工具调用第几次停的**（未超就写"未超，共 N 次"）；11. 伪授权两栏（各带出处）；12. 凭据值零抄录；13. `next=`（含一句：票 176 起跑口那枚还缺什么，以及票 177 AC#3 能不能因此翻勾）。
⚠ **每个"几枚"旁边附可复制命令**；**先测→再写→再提交**。我这段话里每一条前提（含"177 的载具是惰性的""名册八枚""那两枚 fail-closed 不敏感""canary 9301 字节"）**都是〔别人现跑／我转述、未复算〕**：不符就报回并继续做做得动的部分，**不许为了对我那句话去改判据或改测试**。
