# 票 289 / 腿 289-v1 — 判语表（裁决者＝本腿，被验物＝`289-r1`）＋交付名册＋上报清单

## 1. 四格判语

| 格 | 判语 | 凭据（本腿现跑，⛔ 不是转述） |
|---|---|---|
| **AC#1 四句订正成说实话** | **成立** | `git show 73c30dfe` 的 `+` 行逐字读（`10-ac1-four-comments-verdict.md` §2）；每一枚锚本腿自跑命中：`loop.go:603 func callCorr(taskID, callID string, index int) string {`／`loop.go:676 TaskID: taskID, CorrelationID: callCorr(taskID, p.call.ID, i),`／函数体两支 `fmt.Sprintf("%s#call-%d", taskID, index)` 与 `taskID + "#" + callID`（永不等 taskID）／`internal/tools/subagent_197.go:262 parentID := TaskID(ctx)` 且 `:429 ParentTaskID: parentID`／`cmd/wisp/subagent_carrier_197_test.go:173 TaskID: rootID, CorrelationID: rootID`／join 的建表 `subagent_roster_197.go:212-221` 与查询 `:236`；写面＝4 枚件 4 个 hunk，无第五处、无"顺手改其它注释"；`bd124b2a` 现跑存在（rc=0）且提交信息就是那一族事实 |
| **AC#2 行中性自证（本票的牙）** | **成立** | 尺A：`git show --numstat --format="" 73c30dfe` ＝ 4/4、2/2、3/3、5/5，合计 `4 files changed, 14 insertions(+), 14 deletions(-)`；尺B：`git show -U0` 全量非 `+++`/`---`＝28、匹配 `^[+-][[:space:]]*//`＝28 ⇒ **非注释行 0**；第三把尺（比 AC#2 要求的更硬）：剥整行注释后四枚件 `fd692f7e` vs `73c30dfe` **md5 逐枚相同**，且指令形注释＝0；腐烂尺：`:127`／`:190`／`:205`／`:210`／`:236` 两态 **SAME**，夹具 `carrier_197_test.go:173` 未漂 ⇒ **票 220/197 的票面锚没被本票推动** |
| **AC#3 没伤到任何用例** | **成立，带三处具名残余** | 本腿自跑改后全量＝**390/11/2 rc=1**、`0xc0000135`＝0，红名册与 r1 改后名册 **diff 0 行**；对 r1 两份 raw 用本腿的尺复算＝pre 389/12/2、post 390/11/2；新增红 **0 枚**（`comm -13` 空）；残余＝① 改前那一发本腿**无法复跑**（复跑要动工作树，硬约束禁止），只做到"独立复算"；② 那一枚转绿〔原因未定〕，本腿 21 发复跑（含 `-cpu=1`、组串跑、5 进程并发抢 CPU）**零复现**，r1 的"同包串跑时序"归因**不收下**，只收"非本票所致"（凭语义面零的 md5 尺）；③ 同一棵改后树的**第二发整包串跑**翻出 `TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink` 一枚前三发都绿的红（终止码 `0xc000013a`，⛔ 不是 `0xc0000135`）⇒ 名册尺带 ±1 枚噪声，本格判语靠 §3 那把语义面尺立、不靠单发名册 |
| **AC#4 门禁与越界** | **成立** | `GOFLAGS= go build ./...` rc=0；`gofmt -l <4 枚件>` rc=0 空；裸 `gofumpt` rc=127（复现派单那一处，⛔ 没当跳过）＋`$(go env GOPATH)/bin/gofumpt.exe -l <4 枚件>` rc=0 空；`sh scripts/d22scan.sh` rc=0 `clean`，ban#8 射程含 `internal/ 524`＋`cmd/ 116`（注释与 `_test.go` 计入）⇒ 改的 4 句在射程内未触雷；反面 `go run ./tools/d22scan` rc=1（独立模块，正解是脚本那一发）；逐笔名册尺（`git diff-tree -r 898d1f89 73c30dfe 913161ae` 过禁列 grep）＝**0 命中 rc=1**；三枚冻结件／`thresholds.go`／`allowlist.txt`／`testdata/**`／`internal/agent/loop.go` 相对 HEAD **零字节**；脏树具名（770 行在飞、16 枚 ` D design/`、`.gitignore` 为 ` M`）⛔ 未还原、未据它判绿 |

## 2. 上报清单（给编排者，⛔ 本腿一条都没动手改）

1. **派单落点错（具名）**：派单写第五枚在册件＝`internal/panel/ticket283_corr_identity_rulers_test.go:7`，现跑真身＝**`internal/tools/ticket283_corr_identity_rulers_test.go:7`**（`internal/panel/` 下无此件）。r1 报的文件名与落点都对，且"那句不是假话"本腿独立核成立（`internal/tools/bridge.go:269-271` 空 corr 回落成 `req.TaskID`、`:528 corr: orDefault(req.CorrelationID, req.TaskID)` 机械坐实那句说的是既有用例的 fixture；同文件 `:13-14` 自己写明真生产链是 `callCorr`）。
2. **票面枚数错，方向＝少算**：票 289 `AC#3` 写 `internal/panel` 现量"5 枚具名红"，真值**6 枚**——漏了 `TestApprovalCardViewJSONKeysMatchFrontendTypes`（`internal/panel/approval_test.go`）。本腿两条独立尺都指到 6：r1 的 `raw-pre.md` 复算＝panel 6／cmd/wisp 6，本腿自跑 `raw-v1-post.md` 逐名归因＝panel 6（`TestApprovalCardView…`／`TestC21…`／`TestComposer…`／`TestPanelColour…`／`TestStreamLog…35r8` ×2）／cmd/wisp 5（`TestPanelHostRealWindowHopAndLifecycle`／`TestAC4FocusReturn…Gap33r5`／`TestAC13ColdStart…`／`TestAC14AwaitedBinding…`／`TestAC14GoSideEvalPush…`，`TestTicket223…` 那枚已绿）。⛔ 票面一字未改。
3. **票面"四句"这一族今天算五句（新发现）**：`cmd/wisp/subagent_carrier_197_test.go:598-600` 仍在主张 `The join needs a card whose correlation id equals a CHILD task id`，今天已过期——反例现跑在同树里：`internal/panel/subagent_blocked_220_test.go:151-152`（窗口 `CorrelationID: "corr-host-9"`／`TaskID: "child-1"` ⇒ `child-1` 行读 true）。它在四枚件之内，而 AC#1 明文禁改"这 4 枚件里的**其它**注释"⇒ r1 不动是**守规**，本格不因此判失败；要处置请另立一格/另票（同族引文形状的另一处 `cmd/wisp/subagent_blocked_197_test.go:7-10` 是"引用＋当场推翻"，⛔ 不是假话，别一起算进账）。
4. **r1 证据件两处小账**：① `00-anchor-and-baseline.md:48`／`10-four-comments.md:56` 把行侧查询写成 `subagent_roster_197.go:235`，现读是 **`:236`**（`:235` 逐字 `			StreamKey:         row.StreamKey,`）；② `raw-pre.md` 末行**没有 rc**（`grep -n "^rc=" raw-pre.md` 零命中，只有 `raw-post.md:2426 rc=1`）⇒ 改前那一发的 rc 靠叙述而非件内自落。两处都不改判语，只是纪律面。
5. **票 289 没有 `## Progress log` 那一节**（票面 31 行，结构缺项）。按「未定义即停」本腿**没有自建**，翻框也不做，全部留给编排者。
6. **两处先于本票就漂了的票面锚**（⛔ 不是票 289 造成的，本腿用两态同行比对证的）：票 220 `:11` 的 `subagent_roster_197.go:190-210` 今天住 `:212-221`（建表）与 `:236`（查询），漂因＝票 220 自己的落地；票 197 `:66` 的"四枚在跑的态 `:99-104`"——那枚锚原文指的是 **`internal/tools/subagent_197.go`**（不是 panel 那枚同名件），四枚态今天住 **`:98-102`**（`:98 Thinking`／`:99 Settling`／`:100-101 Muted`／`:102 Error`），`99-104` 少了 `Thinking`；该件 `git diff --name-only fd692f7e..HEAD` ＝ 0 行（本票没碰）。
7. **方法论一条（写给今后所有"改前/改后名册作差"的派单）**：本腿把改后态跑了**两发整包串跑**，`cmd/wisp` 的红名册在两发之间自己就变了（第二发多出 `TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink`，`exit status 0xc000013a`，那一枚在 r1 改前、r1 改后、本腿改后第一发都是绿的）。⇒ **单发名册不是可靠尺**；今后要么每态 ≥2 发取交集，要么像本票这样另备一把"语义面零"的硬尺（剥注释后 md5 相同）来承担结论。票 289 恰好两把都有，所以 AC#3 判成立；⛔ 那条新红不算本票账。

## 3. 本腿交付名册（写点唯一＝`.scratch/wisp/probes/289/v1/`，件名一律 `.md`／尺单 `.txt`，⛔ 零 `.out`）

| 件 | 是什么 |
|---|---|
| `00-anchor-and-skeleton.md` | 起手骨架（第 1 笔 `e7af1796`） |
| `10-ac1-four-comments-verdict.md` | AC#1 四句逐字判真伪＋每枚锚现跑＋两把名册尺 |
| `20-ac2-line-neutrality-and-rot.md` | AC#2 尺A/尺B/第三把 md5 尺＋票面锚腐烂尺 |
| `30-ac3-test-comparison.md` | AC#3 三数并排、逐名作差、(a)(b)(c) 三条独立攻击 |
| `40-ac4-gates-and-boundary.md` | AC#4 门禁逐把 rc＋逐笔名册＋禁区零字节＋脏树具名 |
| `raw-v1-diff-u0.md` | `git show -U0 --format="" 73c30dfe` 全量原样 |
| `raw-ac4-gates.md` | 九把尺逐节原样（每节一行 rc） |
| `raw-d22scan.md`／`raw-d22scan-badroute.md` | 正解那一发与反面那一发（独立模块）原样 |
| `raw-v1-post.md` | 本腿改后全量 `-v` 原样（2426 行，末行 rc） |
| `raw-v1-post-run2.md` | 本腿第二发整包串跑（in-situ 复现尺） |
| `raw-ticket223.md`／`raw-ticket223-probes.md`／`probe-conc-1..5.md` | 那一枚用例的 21 发靶向复跑，逐发 rc |
| `roster-fail-pre.txt`／`roster-fail-post.txt`／`roster-fail-v1post.txt` | 三份红名册（作差的输入） |

⛔ 零 push、零源码改动、零 `docs/**` 改动、零 AC 翻框；`frontend/**`／`design/**` 零读零写（只走 `git cat-file -e`／`git ls-tree` 的对象层存在性）。
