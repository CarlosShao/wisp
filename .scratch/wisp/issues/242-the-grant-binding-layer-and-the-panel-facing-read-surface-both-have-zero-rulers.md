# 242 — 一次性令牌的**绑定那一层全仓零尺**，而"面板那面不带令牌"这句话**只由三行注释守着**：两形我都有读数，都不是 AC#5 的洞，但下一枚腿照它们写就会签错字

- Status: **已立，未派**（09-30 15:3x，编排者立；来源＝非实现者验收腿 `197-v1` §4.1 第 3 条＋§1.3 第 2 点＋§1.4，裁决表 `docs/evidence/s1/197-ac5-selfapproval-v1.md`，52,134 字节）。
- ⛔ **为什么不塞回票 197**：那两格**不在票 197 AC#5 的射程里**（被测文件从没宣称扫"跨卡绑定"或"出向读面"）。把它们并回去＝让一格判语承担两枚它没测的东西，而这正是 `197-v1` §3.3 拿来判我那把尺的那句话（"拿一把只量得出 R1 的尺去宣布满足这条裁定"）。
- ⚠ **这不是"本机已被攻破"，照这三行读**：① **现象在哪**＝**测试仪器与注释**，不是生产行为——今天 `PanelItem` 结构上确实没有 grant 字段，`bindDigest` 也确实比对；② **有没有本机被入侵的证据**＝**没有**；③ **最坏后果是什么形状**＝**将来某一版把它改坏而 CI 不响**，以及**下一枚验收腿照现有脚注跑会误判"这格没牙"**。⇒ 这是**防回归的仪器缺口**，不是正在漏的洞。

## 现量（09-30 15:3x 编排者自己跑，别信行号、自己复算）

| 事实 | 读数 | 尺 |
|---|---|---|
| **绑定那一层今天零尺** | `bindDigest`／"不绑定"在全仓 `_test.go` 里**只出现在被验收那一枚文件自身**（`cmd/wisp/subagent_selfapproval_197_test.go`）＝**1 枚文件** | `grep -rln "bindDigest" --include=*_test.go internal cmd`（我现跑） |
| 那一层的**内容** | `grantStore.spend` 比对 `bindDigest(corr, taskID, tool, level, seq, args)`（`internal/agent/approval/approval.go:253` 铸、`:292` 花） | 〔`197-v1` D2 读数：删掉比对 ⇒ 定向尺绿＋`go test ./internal/agent/approval/ -count=1` 也绿（`ok 0.392s`）〕 |
| **"面板那面不带令牌"只有注释** | 三处逐字：`internal/agent/approval/approval.go:223`「The panel-facing surface (PanelItem) has no grant」、`internal/agent/approval/pending_read.go:9`「PanelItem deliberately carries no Params and no grant」、`cmd/wisp/approval_reply.go:28`「PanelItem has no grant field and PanelAPI has no Allow method」。**仪器＝0 枚** | 三处 `sed` 我现跑复认；"零仪器"的读数来自 `197-v1` §1.4：给 `PanelItem` 加 grant 并在 `viewLocked` 填真令牌（M1a）后**全仓 `go test ./...` 无一枚因此变红** |
| ⚠ **造不出"两张不同名"的卡**（现成的载具缺陷） | `cmd/wisp/subagent_selfapproval_197_test.go:109` 逐字 `TaskID: taskID, CorrelationID: taskID` ⇒ 同一次跑里两张卡**共用同一个 correlation id** | `197-v1` §1.3 第 2 点：13 发读数里三枚 id 逐次全等（例 `child=corr1=corr2=45ce09d5-…`）；⇒ 那条 `:466` 的"跨卡借证"实际测的是**新鲜度** |

## 要建什么（两条判据，同一块石头：令牌与"哪件事"的绑定关系今天没人守）

- [ ] **AC#1 真·跨卡：A 的活令牌花到 B 上 ⇒ 必须红**。判据＝造**两枚同时活着、correlation id 不同**的卡，把 A 的**未花**令牌递给 `Native().Allow(B的corr, A的grant)` ⇒ 要拒且**拒因指名"绑定不对"**（不是"这条路由没在跑"）。〔10-08 15:6x 编排者**就地加射程注**、⛔ 未翻勾、原判据一字未改：本格**"递到路由＋拒因指名绑定不对"那一半，在现行形状下路由级不可达**——`242-r3` 现量路由级只能回 `spent-or-never-live-nonce`（route 按 corr 取到的是 **B 自己那本 store**，A 的 nonce 本就不在其中），`misbound` **只有 store 级可达**。⇒ 那一半**归口票 259**：它 `:12` 逐字记着"对外只有一个 `ErrBadGrant`、三种原因并成一句"，`internal/agent/approval/approval.go:477` 逐字写着 **"ticket 259 does not move it"**，而**选边（ⓐ具名降级／ⓑ补牙）是票 259 `AC#1` 明写了归编排者的那一格**。⛔ 任何腿不许为了闭合本格去自己改合并文案。本格**今天已闭合**的是：store 级跨卡用例**有牙**（`M-A` 掏空绑定比对即红，红句 `returned none and not misbound`）＋同形正控（重铸 A 的令牌能花成功）＋`seq` 那句恒真注释已改成说实话。台账＝`A712` §2/§3、`A714`。〕
  ⛔ **前置是本票自己的活**：先把载具的 `CorrelationID` 从 `taskID` 改成真正的 corr（`subagent_selfapproval_197_test.go:109`），否则两枚卡同名、这一发**构造不出来**（上面那行读数就是证据）。
  **正控**＝`197-v1` 的 **D2**（`approval.go:292 grantStore.spend` 不再比对 binding 摘要）：改前该判据**必须绿**、改后**必须红**；⚠ D2 今天**两把尺都绿**＝这就是本格的欠账本身。
- [ ] **AC#2 出向读面不得带 grant**：反射扫 `approval.PanelItem` 的字段名枚数与名字 ⇒ 出现 `Grant`／可答复类字段**就红**；再加一发**能力侧**判据（不是词面）：把 `PanelAPI` 的出向读面扩到"能改变卡的状态"那一形也要红。
  **正控**＝M1a 那一形（给 `PanelItem` 加 grant 并在 `viewLocked` 填真令牌）今天**全仓不红** ⇒ 本格落地后**同一形必须红**。
  ⚠ **落点两候选，这一刀归人**：**甲**＝新钉放 `internal/agent/approval` 本包测试里（**不碰任何冻结件**，编排者推荐）；**乙**＝把 `internal/panel/l2_grant_boundary_test.go` 那族从"入向 envelope"扩到"出向读面"——⛔ **那枚是三枚冻结件之一**，动它必须先有 owner 一句话并在台账落 `A##`（登记编号 `Q-74`，见票面末节）。
  ⚠ 腿原话建议的是"扩那族的射程"，**我照抄就会把一枚冻结件派给别人动**——这一处错记在我自己头上（`A471`）。

## 验收这票要谁裁（非实现者）

- [ ] **AC#3 本格不设产码要求**：判 AC#1／AC#2 两枚判据**是不是同一块石头**（值不值两格）、以及**甲形够不够**——若判"够"，本格按"甲落地"翻勾；若判"要乙"，本格不许勾，写清"等 owner 那句解冻"。

## 禁区

- `frontend/**`／`design/**` 零读零写零转述；`PLAN.md`／`docs/specs/**`／`docs/SLO.md`／`internal/observe/thresholds.go`／golden／`allowlist.txt` 一字节不动。
- ⛔ 三枚冻结件一字不动：`internal/panel/tokens_fourway_test.go`／`internal/panel/l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go`（**AC#2 若走乙形＝本票今天派不出去**，须先解冻）。
- ⛔ 不许**放宽任何断言换绿**；不许新造 D43 之外的状态词；不许把 `SessionID`／grant 塞进 `tools.Decision` 或 `agent.ToolRequest`（`internal/tools/ticket90_test.go:431` 那枚反射钉的射程，票 224 也撞过同一枚）。
- ⛔ **不许顺手把票 197 AC#5 的四条附条件当已完成**（C1／C2／C3 三枚尺的修法归 `197-r3`，见票 197 AC#5 那一格）。

## 排程与串行（共享工作树）

- 本票动 `internal/agent/approval`＋`cmd/wisp` ⇒ ⛔ **与 `197-r3`（同动 `cmd/wisp/subagent_selfapproval_197_test.go`）串行**；与 **票 224-r3**（`internal/tools/bridge.go`＋`cmd/wisp/run.go`）、**票 228-r1**（`cmd/wisp`）串行；两枚跑突变的腿**绝不并发**（互相洗读数）。
- 派单前必做**撞钉预检**（这条规矩本仓已付过两次学费）：先跑 `go test ./internal/agent/approval/ ./cmd/wisp/ -count=1` 把**今天绿的用例名**抄进派单，再逐枚读带 `grant`／`Panel`／`Allow` 字样的断言；⛔ 光 `grep` 新符号名不算预检。
- 起跑名册（我现跑，派单时重取）：`git status --porcelain -- internal cmd`＝**空**；`internal/panel` 4 枚与 `internal/ball` 1 枚红**在册**（成因＝别人在工作树里删了 `design/assets/**`），⛔ **不是**任一腿的地界，只记归因、不判不改。

## 待人拍板（`Q-74`，编排者立，**不阻塞本票甲形**）

- **要拍的那一刀**：AC#2 走**甲**（新钉在 `internal/agent/approval` 本包，⛔ 不碰冻结件）还是**乙**（把 `internal/panel/l2_grant_boundary_test.go` 的射程从入向扩到出向）。
- **推荐＝甲**。理由：甲一步就能让 M1a 那一形变红，且不动任何冻结件。
- **不答的代价**：默认走甲；**乙永远开不了**——除非他给一句"解冻那枚文件"。⛔ 按甲做完，出向读面仍有一支没被那族更大的仪器覆盖（新增 json 键那一形今天有 `TestApprovalCardViewJSONKeysMatchFrontendTypes`，而那把尺**本来就红**在册）。
- **撤销口令**：「242 改乙」（要改回乙只需说这四个字，届时先落解冻记录再派腿）。

## 5. 收非实现者对抗验收腿 `242-v1`＋编排者处置（2026-10-03 09:3x；台账 `A559`；**三格一枚不翻**）

交件凭据（盘上尺复量）：`docs/evidence/s1/242-grant-binding-v1.md` **252 行／29,758 字节**，四枚 commit `3fe377e3`→`8fe080d1`→`850a76b4`→`4df747d8` 逐枚 `git log -1` 存在；九发突变全部还原（`git status --porcelain -- internal/agent/approval`＝**空**，我 09:29 复量），终态 66 PASS／0 FAIL；AC 框未碰、未 push。

**三格判语（腿的原判）与我据此做的裁**：
- **AC#1 ＝ 附条件成立 ⇒ 我不翻勾**。有牙那半（store 层铸/花/烧牌三处：`queue.go:162` 铸、`queue.go:376` 花、`approval.go:298-305` 消费即删，M-B 之外多发红句逐字）我复认；**没牙那半是票面这句话本身要求的两维**——①"拒因**指名**绑定不对（不是"这条路由没在跑"）"：`spend` 返回 `bool`（`approval.go:292`），API 侧只有 `ErrBadGrant` 一枚（`ui.go:130` 把"缺失／已用／不绑定"三种并成一句、`ui.go:127-130` 注释**明写故意不区分**）⇒ 这一维今天**既无实现也无仪器**；②M-B（`queue.go:235` 让允许侧那条路都不跑）⇒ **我那六枚全绿**、同发 12 枚在册路由级用例红＝**六枚用例一枚都没走过 `Allow`／路由**；③名叫"跨卡"那枚（`ticket242_binding_test.go:87-107`）递的是**从未 `issue` 过的字面 nonce**，压根到不了绑定比对⇒绑定层掏空后仍绿＝**恒真**；④M-C3（seq 从摘要里彻底消掉）⇒ **全包 66 PASS**（测试用同一枚 `bindDigest` 自指复算 `:71`）。⇒ **按本仓铁律"判据换成反形还全绿＝它对这件事不敏感，不许当凭据"，这格不翻**；缺的零件转成**票 259** 的 AC#2（拒因可指名）与 AC#1（ⓐ 具名降级／ⓑ 做出牙，待我选边）。⛔ 不许把这格读成"真·跨卡路由有尺"。
- **AC#2 ＝ 不成立 ⇒ 不翻，且我原来那把词面尺要重做**。腿的两发实测：M-D2（`PanelItem.Permitted bool` ＋ 把 `"Permitted"` 补进出向读面名单）⇒ **66 全绿**＝**改个字段名就绕过词面名单**；M-E（照票面那一形给 `PanelAPI` 加 `Allow(correlationID, grant)` 并实现成真花令牌 `p.q.allow(...)`，落点 `ui.go:167`／`gate.go:702`）⇒ **本包仍 66 PASS／0 FAIL**（已验非编译失败冒充红）＝票面那句**能力侧判据零仪器**。⇒ "面板那面不带 allow"今天仍**只由五处注释守着**。归**票 259 AC#3** 三枚能力尺（`PanelAPI` 方法名封闭集／`PanelItem` 按能力判的字段封名单／语义侧"读面不可回填成答复"），⛔ 不为此解冻乙形（腿判"缺的零件可在甲内补"）。
- **AC#3 ＝ 不许勾**：腿判两枚判据**不是同一块石头**（互不敏感实测：M-A2 时 panelface 全绿、M-D1 时 binding 全绿）＝值两格；但"甲形够不够"答**不够**，且**不必解冻乙**（同上）。⇒ 本格按票面自己那句"若判'要乙'，本格不许勾"处理＝**不勾**，理由登记为"甲不够＋缺三枚尺"。

**★ 我自己复认的那条最重的结构事实**（比三格判语更要紧，单独立票）：`queue.go:376` 逐字 `spent := it.grants.spend(nonce, it.bind)`，而 `it.bind` 全仓**只有一处写**（`queue.go:162` 铸造时算出来的那一次），`issue` 存进 store 的也是同一个值（`queue.go:336`）⇒ **`approval.go:303` 那句 `equalSecret(stored, bind)` 在任何经路由的调用里两端永远同值**＝绑定比对**从结构上不可能为 false**。真挡住"A 的令牌花在 B 上"的是另外两处：**store 成员检查**（nonce 不在**这枚 item 自己**的 store 里就直接 false）＋**state 检查**（`queue.go:368-371` `statePending`）。⇒ 我此前在 `A549`／票 242 §4 写过"绑定层已能挡跨卡"**这句不准确，原话不改、就地打旧**：挡跨卡的是成员检查，绑定那一层今天**零独立拦截能力**。旁证（同一批发现）：`ticket242_binding_test.go:111` 那枚 `sameDigest` **零调用者而包仍 ok** ⇒ `build`/`vet` 的 rc=0 对内部有效性是**零证据量**。
**载具前置仍未做**：`cmd/wisp/subagent_selfapproval_197_test.go:109` 现量逐字 `TaskID: taskID, CorrelationID: taskID,`（我 09:3x `grep -n` 复认）＝票面 AC#1 那句"两枚同时活着、correlation id 不同"的前提**在盘上不成立**；该文件写面 `cmd/wisp` 此刻被 `255-r2` 占着 ⇒ 归票 259 AC#4 串行。
**门禁一条具名**：腿报 `d22scan.exe` rc=0 clean 但同发 1 枚 **SKIP**（`frontend/dist/assets/`，由 `frontend/.gitignore` 决定，前人 138/142/143/145 表逐字在册）＝按"SKIP 算红"登记、**不归本票**（那是既有的产物覆盖问题）。

**本票终态**：三格保持 `[ ]`、⛔ 不加 `-done`、⛔ 不撤票（判据本体与修法禁区仍是有效约束）。**欠账全部落在票 259**（`.scratch/wisp/issues/259-the-grant-binding-layer-is-an-identity-…md`）；`242-v1` 建议把 M-D2／M-E 两发**留作常驻自证**，我采纳＝台件留在 `.scratch/wisp/probes/242/v1/`，⛔ 不删（仓内不删东西）。
