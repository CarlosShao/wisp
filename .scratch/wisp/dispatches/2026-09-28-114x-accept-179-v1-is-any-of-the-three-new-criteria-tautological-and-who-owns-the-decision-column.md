# 派单 179-v1（**非实现者验收位，只裁不改**）＝裁票 179 的 AC#1..AC#7：那枚"声明 L0 被当成未分级"的修复**是不是真修好了、三枚新判据里有没有一枚是恒真的**，以及我（编排者）自己替它做的那格红腿名册差**可不可信**

- 派单时刻：09-28 11:4x（`date` 现量于你起手第一步，**别抄我的号**）。
- 票面＝`.scratch/wisp/issues/179-…-no-criterion-holds-that-sentence.md`（AC#1..AC#7 归你裁；**AC#8 我已判"未结案"并转票 183**——你要裁的是"我这个判法对不对"，不是替我结案）。
- 被告的两份陈述＝`docs/evidence/s1/179-declared-l0-refused-r1.md`（18276 字节，实现者自己写的）＋`docs/evidence/s1/179-cli-reread-r2.md`（20739 字节，第二枚程）。⚠ **两份都是被告陈述、不是证据**。
- 提交：`ec520750`（修法＋三枚判据）｜`b6c34d34`（两枚注释指针）｜`91a6c926`（表＋台件）｜`391822a2`＋`4fbda6a5`（第二枚程，零产码）。
- ⚠ **预算硬顶 ≤50 次工具调用**；每裁完一格 commit 一次；**不许勾任何框、不许改产码**（还原协议见 §3）。

## 0. 起手五件

`date`｜`git rev-parse --abbrev-ref HEAD`（**必须 `dev`**；不是＝停手上报。全程禁 `switch`／`checkout <分支>`／`merge`／`rebase`／`reset`／`stash`／`worktree`／`clean`／`restore`）｜`git rev-parse HEAD`（我这边 `95218c4c`）｜`git status --porcelain -- internal/ cmd/`（**必须空**）｜基线：`go test -count=1 ./internal/agent/`＋`./internal/tools/`＋**`./internal/risk/` 单跑**（四包并发会把 `TestResolvePerCallBudget` 打成假红，见台账 `A359`；跑出不同数以你为准并报回）。

## 1. 四问（每问〔成立／成立·带条件／不成立〕＋可复制命令）

1. **三枚新判据逐枚验恒真性**（本单主菜）。对每一枚回答：**产码怎么改会让它红？** 答不出＝恒真。⚠ **至少自己动两枚变异**（别只复算被告报过的那枚）：候选——把 `case RiskL0:` 整支摘掉、把 `j.decide(...DecisionAllow)` 改成不记账、把 `default:` 那支的 fail-closed 条件反过来。每枚变异**先证落地**（`grep -n` 那一行贴出来再跑）。
2. **"守卫那一枚"是不是装饰**：`TestUnclassifiedCallStillRefusedWhenPassThroughIsOff179` 与 `TestDeclaredL1L2RoutingUnchanged179` **在未修码上今天也绿**（实现者自己这么报）——裁清：作为"不许弄坏"的守卫它们合格吗？有没有哪一枚**其实今天就能被造红**（若判"能"，具名给那一发变异；若判"不能"，说清它买到什么）。⚠ 特别核：`AdmitTask == nil` 那一支与"非 nil 时不记账"那一支，是不是各有一枚独立断言，还是共用一枚。
3. **它推翻我的两条前提，你独立复算**：① 桥具默认用**无参** `NewEchoProvider()`、其目录**不写 `RiskLevel`**（`internal/agent/tools.go:120-131`＋`harness_test.go:128`）⇒ "声明 L0 这一形从来没被任何既有用例跑过"；② `cmd/wisp` 本机**不是**测不到东西，带 `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"` 就绿（`scripts/wisp-cli-tests.sh:20`、CI `ci.yml:475`）。⚠ 这两条我都判"实现者对"——**你要独立验，不要沿用我的判**。第 ② 条若为真，票 98 的"可结案候选"就成立；若不为真，报回来。
4. **⚠ 我最不放心的次生面**：把声明 L0 改成"loop 直接放行并记账 `allow`"——**journal 那一行的 `decision` 从此由谁拥有**？现读 `internal/agent/journal.go` 与 `loop.go` 里那段"ok came back for a declared L1/L2 call the loop booked no decision"的注释，裁一句：**L0 记 `allow` 会不会与"门拥有那一列"的既有分工打架**（这是票 12 的 A13 裁定那一族），以及**审计/成本读数上会不会因此出现"宿主替用户批准了只读操作"这种误读**。⚠ 不许动 `internal/risk/**`、不许动 `bridge.go`；若你判它变了，报回归票 177／175，**别自己开第二支修法**。

## 2. 门禁（你自己重跑，别引任何人的数）

`go test -count=1 ./internal/agent/`＋`./internal/tools/`＋`./internal/risk/`（单跑）｜`sh scripts/d22scan.sh`（应 rc=0；`ban #8 internal/` 我 11:4x 现量 **431**）｜`bash tools/d22scan/runtests.sh -C tools/d22scan ./...`（基线现量＋名册两向 `comm` 差集）｜`sh .scratch/wisp/probes/154/gate-clauses.sh`（⚠ 今天退码本就是 1，唯一在册红腿 `G6neg 声明=ring 基线=1枚 实测=2枚`＝票 178；**你要比的是红腿名册**，做法：`grep "BAD" <输出文件> | grep -oE "腿=G[0-9a-z]+" | sort -u` 成名册再 `comm -3`。⚠ **我在 `A361` 里替 `179-r1` 做过一次这格差集＝"一字未变"，那一格是我的读数、不是它的**——你要自己重取，并裁"编排者替做凭据"这一形可不可接受）｜`gofumpt --version` 现跑（⚠ 裸名不在 PATH，用 `"$(go env GOPATH)/bin/gofumpt"`，v0.12.0）＋对它名下两枚 `.go` 的 `-l`。⚠ **`probes/161/r6/flip-declaration.sh` 别跑**（跑一次就脏跟踪日志）。⚠ **此刻 `183-a1`（只读）可能同时在飞**：它只走 `-overlay` 不改跟踪文件；若你发现 `internal/**` 出现非你的改动，**只登记、不评论、不碰**。

## 3. 写面（超出即越权）

`docs/evidence/s1/179-declared-l0-refused-accept-v1.md`（新建，你的表）｜`.scratch/wisp/probes/179/v1/**`（变异台件；⚠ `logdir` 取脚本自身目录，**不许用 `runtime.Caller` 在 `-overlay` 下取路径**——`179-r2` 刚把读数误落进 `cmd/wisp/logs/` 并留了证）｜票 179 的 Progress log（**只追加、不勾任何框、留最后一步**；写前 `git status --porcelain -- 那枚票面` 必须为空）。
**禁改**：`internal/**`（§1 的临时变异只在还原协议内做）、`cmd/**`、`tools/d22scan/**`、`allowlist.txt`、`.github/**`、`docs/PLAN.md`、`docs/specs/**`、`docs/reports/**`、`frontend/**`、`design/**`、别人的票面与证据件、`probes/**` 既有台件（只读）。
⚠ 现场躺着 `design/**` 未提交删除、`.gitignore`／`probes/152/my152.py`／`probes/161/r6/logs/flip-*`／`docs/evidence/s1/152-*.md`（` M`）⇒ **不提交、不还原、不删除、不评论**。
**`git add -A`／`git add .`／`git commit -a` 一律禁**；一步式 `git commit -q -F - -- <显式路径>`；heredoc 加引号（`<<'EOF'`）；含反引号的中文走编辑工具；串命令里可能无命中的 `grep` 用 `;` 或 `| cat`、**别用 `&&`**。**只 commit，绝不 push。**⚠ **更正自己写错的东西只许追加新提交——本仓禁 `--amend`**。

## 4. 交件报告（顺序固定）

1. step-0 五件；2. **本程没测什么**；3. 三枚判据逐枚表（每枚：它钉什么／怎么改会红／我动不动了它）；4. **我自己动的那两枚变异**（落地证明＋两向读数＋还原后空 status）；5. 守卫那一问的答案；6. **推翻我两条前提的独立复算结果**（含票 98 那格你判成什么）；7. **次生面（journal 的 `decision` 归属）那一问**；8. 门禁全部读数＋红腿名册差集＋`gofumpt --version`；9. **对编排者代做那格（`A361` 的红腿名册差）的裁定**；10. 被拒／没成功的调用（取数前还是后）；11. 有没有跑过删除命令；12. 工具调用**用了几次 vs 硬顶 50**（超支要自报）；13. 伪授权两栏（各带出处）；14. 凭据值零抄录；15. `next=`（含一句：**票 179 该勾哪几格、AC#8 我判"未结案转票 183"对不对**）。

⚠ 每个"几枚／几行"旁边附可复制命令；**先测→再写→再提交**。
⚠ 我这段话里每条前提（含"三枚判据"、"两枚守卫今天也绿"、"零产码"、"红腿名册一字未变"、"ban #8 internal/=431"）**都是未验证断言**：不符就报回并继续做做得动的部分，**不许为了对我那句话去改判据或改测试**。
