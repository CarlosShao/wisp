# 票 259 非实现者终裁（腿 `259-v1`）— 读数与突变四件套

- 现量时刻：2026-10-08 18:44 +0800；HEAD `7f2ecf35810dd659500bbad880bb899fbd686847`（dev）。
- 起手锚＝`00-anchor.md`（commit `0af6594d`）。
- 判"具备/不成立"一律锚调用形状与现跑读数；行号与被引短语同一发取。
- 基线（干净树，注 sherpa/onnxruntime PATH）：`go test . -run 'TestTicket259' -count=1 -v` → **12 PASS / 0 FAIL，rc=0**（approval 包，0.031s）。还原后复跑同 12 枚 **PASS，rc=0**。

## 台账已收件（照实引用，非新发现）

- `A619`（10-05 12:1x）：收 `259-r1`（八笔，末 `77a23fa3`；AC#2 四拒因各有名＋AC#3 三枚能力尺；件 `.scratch/wisp/probes/259/r1/evidence.md` 483 行）；落具名裁定 **AC#1 选 ⓐ＝具名降级**（撤销口令「259 AC#1 改 ⓑ」/票面「259 改形 ⓑ」）。
- `A621`（10-05 12:5x–13:0x）：收 `259-r2`（五笔：`549cb218`→`6a021830`→`789a02e2`→`0359d168`→`7de53ae8`；ⓐ 进树，剥注释后两文件逐字节相同）；行号漂移已登记（旧 `approval.go:287/:298-305` → `:415/:485-502`）⇒ 本轮全部按新号现取。
- 10-08 两笔 `df1b962e`/`6ad29787` 属票 242 的 `242-r3`，只动 `ticket242_binding_test.go`＋台账，未碰 `spend` 产码 —— 已核，见 AC#5。

## AC#1（ⓐ 具名降级）— 判语：成立

**今天盘上的口径（现取，逐字或具名）**：
- `internal/agent/approval/approval.go:361-362`："(2) is the one a cross-card attempt lands on: A's nonce is not a row in B's store"（真防线＝store 成员扫描，命中处 `:574 return denialSpentNonce`）。
- `approval.go:364-368`："The binding digest is NOT a fourth gate, and it stops nothing across cards on its own... the two sides of the comparison inside spend are one string by construction - that branch cannot come out unequal on any path the gate can reach."
- `approval.go:405-415`（bindDigest 注）："the wording on this function used to claim that it was, and ticket 259 AC#1 form (a) downgraded it."
- `approval.go:436-446`（grantStore 注）：per-item split 是真检查；companion guard＝队列 state 拒绝，"not a cross-card attempt, because the other card is sitting there pending"。
- `approval.go:528-557`（spend 注）：花侧唯一产码调用者＝`Queue.allowScoped` 传 `it.bind`（"value against itself"），"ZERO independent power"；残余风险与 ⓑ 未设计（答复面无 tool/args/level/seq ⇒ 动它＝扩 C17 入向面）逐句在册。
- `queue.go:44-62`（qitem.bind）：guard 1＝store per item / guard 2＝队列 state 拒绝 分开写；`queue.go:341-347`（grantNonce）与 `:404-413`（花点）同口径。
- 事实核：`it.bind` 全仓**唯一写点**＝`queue.go:177 it.bind = bindDigest(...)`；读点只有 `:358 issue(nonce, it.bind)` 与 `:414 spend(nonce, it.bind)`（其余为注释）⇒ 路由上两端同值、比对恒等成立。

**改前那句话说没说过"能挡跨卡" —— 说过，逐枚（来源＝r2 存档 `.scratch/wisp/probes/259/r2/before/`，HEAD 的 git 历史同一内容）**：
- `before/queue.go:44-47`："so a nonce cannot be replayed onto a different request or onto a later re-issue of the same one."（这就是那句要撤的"能挡跨卡"）
- `before/approval.go:378`："bindDigest is the domain separator that ties a grant to ONE pending item."
- `before/approval.go:356-358`：authority (3)＝"the grant's binding digest matches the item it is spent on."

**谁改的**：`6a021830`（r2 第 2 步；`approval.go` +79/-7、`queue.go` +36/-4，**只动注释**；A621 复认"剥注释后逐字节相同"）＋`b6b1d6a4`（r1，含拒因形状）＋`789a02e2`（收尾注）。
**HEAD 残留扫描**：`git grep -nE '能挡|different request|ties a grant' HEAD -- '*.go'` ⇒ 编译树内零命中"能挡跨卡"（唯一 `queue.go:41` 命中是 correlation `names` 的"不是给另一个请求的回执"，与绑定无关；另一类是 `before/` 存档）；`git grep -nE '跨卡|cross-card' HEAD -- '*.go'`（非测试）⇒ 3 命中全在 `approval.go:361/:443`、`queue.go:59`，语义全是"真防线是 guard 1/2"。
**测试有没有把新文案逐字钉住**：没有。`git grep -nF` 对 `stops nothing across cards`／`zero independent power`／`NOT a fourth gate`／`-i fourth gate`（射程 `*_test.go`）⇒ 四把 rc=1。本区唯一逐字钉住的是**对外合并句**（`ticket259_denial_rulers_test.go:250`、`cmd/wisp/approval_reply_201_test.go:387`），该句 r2 一字未动（按裁定保持合并）。
**测试面同族口径不跟降级**＝编排者已裁（`A621` §6②），`ticket242_binding_test.go` 的 `...CannotSpendAnotherItemsGrant` 名字保留；其 `:166/:171` 注释本身写的已是"per-item split 才是真拒绝者"。
本轮**不种突变**（文案格，读数即凭据）。

## AC#2（拒因可指名）— 判语：成立

- 形状（现取）：`approval.go:558 func (s *grantStore) spend(nonce, bind string) grantDenial`（**已不是 `bool`**）。四因常量：`denialMissingNonce :488`／`denialSpentNonce :493`／`denialMisbound :502`／`denialNotPending :506`（第四因由 queue 产、不是 store 产，`:503-505` 明写）；label 逐枚不同（`:511-526`）。
- 四种各自可断言：拒因尺 6 枚（`ticket259_denial_rulers_test.go` :115 缺失／:134 已用或从未发放／:165 绑定不对／:192 状态离 pending／:226 四 label 两两不同／:249 对外仍合并）；审计行 `GRANT-DENY ... denial=%s`（`queue.go:397/:423/:428`）＋合并句 `FORGED-OR-STALE`（`:424`）保持；对外仍单值 `ErrBadGrant`（`ui.go:130`）＋`cmd/wisp/approval_reply.go:227/:272` 两份合并文案仍在。
- **最小形状差异（今天缺什么）＝0**（在"内部＋审计区分／对外合并"裁定范围内已成立）。唯一登记事实：`denialMisbound` 在**路由级不可产生**（只有包内直调 store 两枚手递摘要能看到）——这是 AC#1 ⓐ 已登记的恒等式结构事实，不是 AC#2 缺口。⛔ 不写实现方案。
- 基线 12/12 PASS 覆盖以上全部。

## AC#3（三枚能力尺）— 判语：成立（三枚在盘，三发突变各自见红且逐枚可归因）

在库核对：两枚尺件只在 `internal/agent/approval/`（`git ls-tree -r --name-only HEAD -- internal/agent/approval/ internal/panel/ | grep -i ticket259` ⇒ 2 中 in approval、panel 0 中），各 6 枚 `func Test`，落笔 `b6b1d6a4`。映射：①＝`...PanelAPIMethodSetIsClosed` `:66`＋`...PanelCarrierCarriesNoAnswerVerb` `:93`；②＝`...PanelItemFieldCapabilityFence` `:169`（＋值级 `:184`／不变性 `:216`）；③＝`...ReadFaceCannotBeLaunderedIntoAnAnswer` `:283`。

| 突变 | 种在哪（sed 复量） | rc | 红句 file:line | 还原 |
|---|---|---|---|---|
| Mu-1（M-E 形） | `ui.go` PanelAPI 加 `Allow(correlationID, grant string) error`（:169）＋`gate.go` 加 `func (p panelAPI) Allow(corr, grant string) error { return p.q.allow(corr, grant) }`（:703，真花令牌） | 1 | `ticket259_panel_capability_rulers_test.go:73`（gained method "Allow"）、`:82`（4≠3 计数）、`:101`（carrier 携带）、`:111`（Allow 形断言命中） | hash 逐字回（ui `a4cee69e`／gate `8d93476e`）＋porcelain 空 |
| Mu-2（M-D2 形） | `ui.go` PanelItem 加 `Permitted bool` | 1 | `...rulers_test.go:174`（PanelItem.Permitted has type bool）；同族值级 `:184` 与不变性 `:216` **仍绿**（归因干净） | hash 逐字回（ui `a4cee69e`）＋porcelain 空 |
| Mu-3（读面回填） | `ui.go` PanelItem 加 `Grant string`；`queue.go:578` viewLocked 填 `q.leakedNonceFor259(it)`＋`:582` 取一枚活 nonce 的 helper | 1 | `:199`（Head/View 渲染活令牌值）、`:229`（撤销后字段变化）、`:301`（读面读到的值当 native grant 花出去**成功了**）、`:308`（控制也死：牌已被洗掉）；方法/名册/字段 kind 三枚**仍绿**＝③抓到了②看不到的纯字符串泄露 | hash 逐字回（ui `a4cee69e`／queue `b92a93bc`）＋porcelain 空；复跑 12/12 绿 |

诚实登记（我没一次做对的）：Mu-1 首种 `sed` 双命中（NativeAPI 也有一行同形 `Reject(...)`），产生重复方法名会编译不过——已当场按行号删回、留 PanelAPI 一枚；Mu-3 首种双命中状态结构体既有的 `Grant string`（重复字段，build failed，非红），修成单枚后跑；还原时 Ui 行号已漂（57 非 58）、queue 多留一空行，逐枚按 diff 修回。另一中间态读数：**只加 `Grant string` 字段而不填值**时 12/12 全绿——值级三枚尺只对"值真的过界"咬人（by design）；②的 kind 封名单看不到 string 字段，正是③存在的理由。

## AC#4（载具前置）— 判语：不成立

`cmd/wisp/subagent_selfapproval_197_test.go:109`（`sed -n '109p'`＋`cat -A` 同发）逐字＝`^I^I^ITaskID: taskID, CorrelationID: taskID,`——两枚 id 仍是同一变量，未拆。该文件最后被碰＝`a818df46`（09-30 12:05，早于本票立票 10-03）⇒ 无任何后续 commit 动过它。

## AC#5（越界检查）— 判语：成立（盘上无越界）

尺＝对 15 笔票相关提交（r1 八笔 `ea522444,b6b1d6a4,50e0299c,e8a1bec3,c854a8e0,2efc8940,6fb1af11,77a23fa3`；r2 五笔 `549cb218,6a021830,789a02e2,0359d168,7de53ae8`；另计 242-r3 两笔 `df1b962e,6ad29787`）逐笔 `git show --name-only`，滤 `docs/PLAN|docs/specs|thresholds.go|golden|d22scan/allowlist|frontend/|design/` ⇒ **零命中，rc=1**。正控：同尺对 `dc694430`（动过 docs/specs 的提交）⇒ 命中 `docs/specs/SPEC-12-roadmap-governance.md`，rc=0 ⇒ 滤网活着。golden 名册含 `internal/agent/testdata/golden/*.sse`。`6ad29787` 只多动台账 `docs/reports/pending-and-issues.md`（不在禁列）。

## 结论

五格：AC#1 成立 ／ AC#2 成立 ／ AC#3 成立 ／ AC#4 不成立 ／ AC#5 成立（盘上无越界）。零翻框、零票面/台账改动、零 push；三枚冻结件（`internal/panel/tokens_fourway_test.go`／`l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go`）一字未动（未读内容）。突变已逐枚还原（hash 逐字回＋porcelain 空＋复绿 12/12）。
