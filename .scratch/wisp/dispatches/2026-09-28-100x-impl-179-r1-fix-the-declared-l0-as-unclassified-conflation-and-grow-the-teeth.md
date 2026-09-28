# 派单 179-r1（**写手腿**）＝把 `Loop.decideRisk` 混判"声明 L0"与"真未分级"那一支修回 D4，并给"未分级仍 fail-closed"补上常驻的牙——**这一枚修的是门，所以正反两形缺一形就等于把门改松**

- 派单时刻：09-28 09-5x（`date` 现量于起手第一步，**别抄我的号**）。
- 票面＝`.scratch/wisp/issues/179-loop-deciderisk-treats-a-declared-l0-as-unclassified-so-every-readonly-builtin-is-refused-on-the-real-cli-and-no-criterion-holds-that-sentence.md`（**票面里已有八把尺与现量读数，先读票再动手**；AC 共八格，AC#8 由编排者终判，你只做到 AC#7＋把 AC#8 的读数跑出来）。
- ⚠ **预算硬顶 ≤55 次工具调用**；**到顶即停手回禀**（不许"再两下就好"）；**超支必须在报告里自报枚数**（`176-r1` 是本轮第一枚自报的，按它办）。
- ⚠ **每裁完一格 commit 一次**；**只 commit、绝不 push**；**不许勾任何不属于你的框**（票 179 的 AC 框：AC#1..AC#7 你可勾自己真做完并真测过的；AC#8 留给我）。

## 0. 起手五件（缺一件就别往下写）

`date`｜`git rev-parse --abbrev-ref HEAD`（**必须 `dev`**；不是＝停手上报，一个字不写。全程禁 `switch`／`checkout <分支>`／`merge`／`rebase`／`reset`／`stash`／`worktree`／`clean`）｜`git rev-parse HEAD`（我这边 `cf527f32`，你看到的更新就以你的为准并报回）｜`git status --porcelain -- internal/ cmd/`（**必须空**）｜基线：**单包**跑 `go test -count=1 ./internal/agent/`＋`./internal/risk/`＋`./internal/tools/`（**⚠ `./internal/risk/` 必须单跑**：四包并发时 `TestResolvePerCallBudget` 会因 CPU 争用红成 `1.214 ms/op`，同腿单跑 `0.512 ms/op` PASS；这是仪器吃争用，**不是回退、不许动 `thresholds.go` 一字节**）。⚠ **禁全仓 `./...`**；⚠ **`./cmd/wisp/` 在本机测不到任何东西**（票 98 仍 open，加载期缺 `sherpa-onnx-c-api.dll`；`go test -count=1 -run XXX_NONE_PKG ./cmd/wisp/` 同样 `0xc0000135`）⇒ CLI 那一面**只能走 `-overlay` 台件**（样板：`.scratch/wisp/probes/176/r1/overlay-e2e.json`）。

## 1. 五问（每问给〔成立／成立·带条件／不成立〕＋可复制命令＋两向读数）

1. **先复现红，再动手**：在**未修码**上把你选的判据形状先写出来跑一次，贴出"它现在就是红的"那一行（尺与读数都要）。⚠ 如果你的用例在未修码上**就是绿的**，那它没在钉这一枚缺陷——**重写，别提交**。
2. **修法只许动那一支**：`internal/agent/loop.go:773-790` 的 `switch`（**以及它上面 `:755-771` 那段注释里"没写声明 L0"的那一格**）。判据＝`RiskL0` 显式成为一支；`""`（`RiskUnclassified`，`internal/agent/tools.go:37`）仍走 fail-closed。⚠ **同文件里已有正确写法可对照**：`riskLabel`（`loop.go:1084-1089`）按 `r == RiskUnclassified` 判"未分级"——**照这个形状问，别自己发明第四种**。
3. **⚠ 两支毒修法一律禁，判"该用哪一支"就停手上报**：**(i)** 在组合根（`cmd/wisp/run.go` 的 `agent.Config`，`:606-613` 一带）把 `PassThroughUnclassifiedRisk` 设成 `true`——那枚开关的语义是"给**未分级**用的宿主级政策开关、不是裁决"（`tools.go:21-25`、`loop.go:767-770` 两处注释逐字），设成 `true`＝连真未分级一起放过＝**放宽门**；**(ii)** 把 `default:` 改成一律放行／或把 `RiskL0` 并进"直接 ok 且不记账"之外的任何第三种状态。⇒ **这两种都不许做，也不许"先做了再说回头收紧"。**
4. **"修完会把谁顶出去"要你自己答一遍**（不许抄票面）：把 `Loop` 放过声明 L0 之后，**(a)** 一枚真 L2 的工具调用仍然必须走到审批门——给出凭据（现读 `Bridge.Execute`：C3 `bridge.go:255-256` → C19 `:272` → 模式筛 `:286` → 门），并**动一发变异**证明这条链不是装饰；**(b)** 名字不在目录里的调用（`loop.go:610 info := byName[c.Name]` 取到零值）必须**仍被拒**——这是 AC#3 那发，它**在未修码上今天也绿**，所以它是"不许弄坏"的守卫、**不许拿来充当本票新增的牙**；**(c)** `AdmitTask == nil` 时 L1/L2 仍拒（`loop.go:777-779` 那句逐字文案），`!= nil` 时仍"放行但不记账"。
5. **端到端那一发必须真跑出来（本票存在的唯一理由）**：复用 `176-r1` 的 `probes/176/r1/zz176r1_e2e_test.go` 形状（⚠ **不改它**，你自己的台件放 `probes/179/r1/**`，`logdir` 取脚本自身目录），拿到**后半截**读数：模型续读宿主自己写下的指针 ⇒ 在 CLI 上**不再出现 `风险未分级`**、产物逐字节读回。⚠ **分岔要报准**：如果新读数是 `…L2 级…已拒绝执行` 或带 R4 的拒绝，那是**另一枚**缺陷，**照实写"本票的 L0 混判已修、端到端仍不通，卡在哪一句逐字"**，不许写成已结案、不许顺手去改 `internal/risk/**` 或豁免那一族。

## 2. 门禁（你自己重跑，别引任何别人的数）

`go test -count=1 ./internal/agent/`（判据主战场）｜`./internal/risk/`**单跑**｜`./internal/tools/`｜`sh scripts/d22scan.sh`（应 rc=0；`ban #8 internal/` 我 09:59 现量 **430**、`ban #7 internal/tools/` **21**，你动了文件就会变，**逐名说清涨在哪枚文件**）｜`bash tools/d22scan/runtests.sh -C tools/d22scan ./...`（＋名册两向 `comm` 差集）｜`sh .scratch/wisp/probes/154/gate-clauses.sh`（⚠ **今天退码本就是 1**：唯一红腿 `G6neg 声明=ring 基线=1枚 实测=2枚`＝票 178 在册那枚，**不是你的锅，也不许为它改仪器**；你要报的是**红腿集合有没有新增**）｜`gofumpt --version` 现跑＋对你名下两枚 `.go` 的 `-l`。⚠ **`probes/161/r6/flip-declaration.sh` 别跑**（跑一次就脏跟踪日志）。⚠ **最终那一次门禁读数必须在你的最后一枚 commit 之后取并贴出**（本轮已有两枚程用过期读数冒充终态）。⚠ **G3 陷阱**：不许给 `(*Loop)` 加收任何带 `taskID` 的**导出**方法（普查尺会把安静腿打红）。

## 3. 写面（超出即越权）

`internal/agent/loop.go`（**只 `decideRisk` 那一支＋它上面那段注释**）｜新判据件 `internal/agent/*_test.go`（命名带 `179`）｜注释两行（票 179 AC#6：`cmd/wisp/run.go:357` 那句"zero production call sites"要与 `:633` 对齐、`internal/tools/bridge.go:212` 引的 `loop.go:719-738` 漂到 `773-790`）——**只改指向，不改判断内容，一字节产码行为都不动**｜票 179 的 Progress log（**只追加**）｜`docs/evidence/s1/179-declared-l0-refused-r1.md`（你的表）｜`.scratch/wisp/probes/179/r1/**`（台件与读数）。
**禁改**：`internal/risk/**`｜`Bridge.Execute` 的裁决链与 `bridge.go` 除上面那一行注释之外的一切｜`PassThroughUnclassifiedRisk` 的语义／默认值／任何生产赋值点（AC#5）｜`internal/tools/task.go` 与 `task_backfill.go`｜`thresholds.go`／golden／审批超时常量／`allowlist.txt` 一字节｜`docs/PLAN.md`（含 `:1532` 那行 DEFERRED）｜`docs/specs/**`｜`docs/reports/**`｜`frontend/**`、`design/**`（别家归属，**不读不写不引不转述**）｜票 174/175/176/177/178 的票面与别人的证据件｜`probes/**` 既有台件（**只读**，`176-r1` 那枚 e2e 台件也不许改）。
⚠ 现场躺着 `design/**` 未提交删除、`.gitignore`／`probes/152/my152.py`／`probes/161/r6/logs/flip-*`（` M`）／`docs/evidence/s1/152-*.md`（` M`）⇒ **不提交、不还原、不补完、不评论**。
**`git add -A`／`git add .`／`git commit -a` 一律禁**；一步式 `git commit -q -F - -- <显式路径>`；heredoc 必须加引号（`<<'EOF'`）；含反引号的中文段落**走编辑工具**，别走 shell。**只 commit，绝不 push。**⚠ **更正自己写错的东西只许追加新提交——本仓禁 `--amend`**；⚠ 临时件**只建不删**；⚠ 还原被你自己改脏的跟踪文件用 `git cat-file blob HEAD:<路径> > <路径>`（**不许 `checkout`／`restore`**）。⚠ shell 串命令里可能无命中的 `grep` 一律用 `;` 或 `| cat`，**别用 `&&`**（编排者今天这样把一枚没跑成的读数写进过台账 `A357`）。

## 4. 交件报告（顺序固定，缺格＝退回）

1. step-0 五件（含你的 HEAD）；2. **本程没测什么**（逐条，别写"应该没问题"）；3. **未修码上的红读数**（AC 主体的凭证，逐字贴那一行）；4. 修法那一支改前／改后逐字对照（`git diff` 那一支的行，别整文件贴）；5. **正反两形的两向读数**（AC#2 正向、AC#3/#4 两枚守卫）；6. **"顶出去谁"那一问的 (a)(b)(c)＋你动的那发变异**（变异**先证落地**：`grep -n` 那一行贴出来再跑；还原后 `git status --porcelain -- internal/` **必须空并贴出来**）；7. **端到端那一发**：新读的 `wisp run` 尾态＋那句回执原文；如果仍是红的，红在哪一句逐字；8. 门禁全部读数（⚠ **最后一枚 commit 之后**）＋红腿集合有无新增；9. 被拒／没成功的调用（取数前还是取数后）；10. 有没有跑过删除命令；11. 工具调用**用了几次 vs 硬顶 55**（超支要自报）；12. 伪授权两栏（各带出处）；13. 凭据值零抄录；14. `next=`（含一句：**票 176 AC#3/AC#4/AC#5 现在能不能翻勾、票 177 AC#3 的端到端条件卸不卸得掉**——以你手上的读数回答，别以我这票面的措辞回答）。

⚠ **每个"几枚／几行／几个"旁边都附可复制命令**。**先测→再写→再提交。**
⚠ 我这段话里的**每一条前提**（含"`decideRisk` 只认 L1/L2""全仓零枚用例钉那句文案""三枚只读工具声明为 L0""未知名字今天仍 fail-closed""放过 L0 不会削弱真门"）**都是未验证断言**：不符就**报回并继续做做得动的部分**，**不许为了对我那句话去改判据或改测试**。
