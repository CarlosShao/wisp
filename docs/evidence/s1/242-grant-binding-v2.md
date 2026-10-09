# 票 242 三格重裁 —— 非实现者对抗验收腿 `242-v2`（2026-10-09）

被审对象：票 242 的 **AC#1／AC#2／AC#3** 三格判据（`.scratch/wisp/issues/242-…-zero-rulers.md`，现量 63 行、三格全未勾）。
重裁理由＝票 259 已结案 `-done`（未勾＝0）并落了两族尺，票 242 自己的"载具前置"也落了。
⛔ 本件不翻任何勾、不改票面、不改台账；判语只出自本程亲手跑的读数。

起手锚与基线名册＝`.scratch/wisp/probes/242/v2/00-anchor.md`（HEAD `2bdf385d`，本腿自跑
`go test ./internal/agent/approval/ -count=1 -v` ⇒ rc=0／**90 PASS／0 FAIL**／`ok … 0.372s`）。

⚠ **基线枚数顶回**：本票 §5 与派单转述里的「66 全绿」是 10-03 快照；今天同包名册 **90 枚**。
本件每一发的"绿/红"按 90 枚现量报。

⚠ **同目录已有前腿**：`probes/242/v2/anchor.md`＋`verdict.md`（10-08 11:52 同名腿，错法全在仓外用
`go test -overlay` 造）。本派单定式＝⛔ `-overlay` 对读盘的尺结构性失明 ⇒ 那枚腿的读数**不作本腿凭据**，
本腿七发全部种在盘上并逐发还原；前腿三件**未覆盖、未删**，本腿只写新文件名。

## 0. 突变纪律（七发逐件）与门禁

每一发四件套＝种前 `git hash-object` → 种后 `sed -n 'Np'` 复量该行 → 红句逐字（带 `文件:行`）→
还原后 hash 逐字等值 ＋ `git status --porcelain -- internal/agent/approval` 空。
种前四枚产码件 hash（本腿现量）：

| 文件 | hash |
|---|---|
| `internal/agent/approval/approval.go` | `67fb146898dc7b235623bb8ad2fac0e0b2465fe1` |
| `internal/agent/approval/ui.go` | `55bcd1daf251a9ef513ef43b3c63605787252b13` |
| `internal/agent/approval/gate.go` | `78e2d28a8b676754514bd80898ace8c5a974d3ce` |
| `internal/agent/approval/queue.go` | `66fec7ae375344d7bc741270b42c4ec5050cf41f` |

收工复量＝**四枚逐字等值**，`git status --porcelain -- internal cmd` **空**，包复绿 `ok … 0.380s`。
⛔ 零 `internal/tools/**` 写（只读＋跑：`go test ./internal/tools/ -count=1 -run 'TestTicket283|Corr'` ⇒
2 PASS／`ok … 0.063s`，票 283 基线未被本腿洗）；⛔ 零 `cmd/wisp/**` 写（只读＋跑）；
⛔ 三枚冻结件零触碰；⛔ 无 `-overlay`、无 `t.Skip`、无放宽任何断言。

| 发 | 落点（种后 `sed -n` 逐字复量） | rc | PASS／FAIL |
|---|---|---|---|
| MU-A | `approval.go:569` → `if false && !equalSecret(stored, bind) {` | 1 | 86／**4** |
| MU-M | `approval.go:574` → `return denialNone // MU-M planted by leg 242-v2: membership scan gutted` | 1 | 81／**9** |
| MU-D2a | `ui.go:57`（PanelItem 新字段）→ `	Permitted     bool` | 1 | 88／**2** |
| MU-D2b | 同字段换类 → `	Permitted     int` | 1 | 未枚举 `-v`（`-count=1` 非详单）／**1** |
| MU-D2c | MU-D2b ＋ `queue.go:578` → `		Permitted:     it.grants.live(),` | 1 | 88／**2** |
| MU-E | `ui.go:171` → `	Allow(correlationID, grant string) error` ＋ `gate.go:703` → `func (p panelAPI) Allow(corr, grant string) error   { return p.q.allow(corr, grant) }` | 1 | 88／**2** |
| MU-OUT | `ui.go:130` → `	ErrBadGrant = errors.New("approval: 原生令牌无效（与本次请求不绑定）")`（定向 `-run 'Ticket259R1Outward|Ticket259R1Denial'`） | 1 | 4／**1** |

⚠ MU-E 的红**不是编译失败冒充**：同一次 `-v` 跑里先打印了 88 枚 `--- PASS`，红句是断言句（见 §2）。

---

## AC#1 真·跨卡：A 的未花令牌花到 B 上 ⇒ 要拒，且拒因指名"绑定不对"

**判据原文（票面 `:18` 逐字，⛔ 未换写法）**：「造**两枚同时活着、correlation id 不同**的卡，把 A 的**未花**
令牌递给 `Native().Allow(B的corr, A的grant)` ⇒ 要拒且**拒因指名"绑定不对"**（不是"这条路由没在跑"）」
＋10-08 射程注（同格内）：路由级那一半不可达、`misbound` 只有 store 级可达、归口票 259。

**载具前置（本腿现量，非抄派单）**：`cmd/wisp/subagent_selfapproval_197_test.go:116` 逐字
`			TaskID: taskID, CorrelationID: corrID,`；两枚 corr 在 `:396`／`:463` 各自吃
`childID+"-corr-1"`／`childID+"-corr-2"`；`startChildWrite197(rt *agentRuntime, taskID, corrID, path string)` 在 `:101`；
跨卡那一发在 `:473` 逐字 `rec.foreignErr = rt.gate.Native().Allow(context.Background(), card2.CorrelationID, card1.Grant)`。
⇒ **票面 `:14`／§5 `:60` 说的"两枚同名、构造不出来"在盘上已过期**（票 259 §12 那笔 `77150dcc` 拆的；
票 259 面上写的 `:109` 也已漂到 `:116`）。本腿跑该载具：
`PATH="$PWD/third_party/sherpa-onnx:$PATH" go test ./cmd/wisp/ -count=1 -run 'Test197Subagent…Assembly' -v`
⇒ **rc=0／2 PASS／0 FAIL／0 SKIP**（无 `exit status 0xc0000135`）。

**本腿实测的三件读数**：

1. **store 级指名 `misbound` 有牙**（MU-A 掏空绑定比对）⇒ 4 枚红，逐字：
   - `ticket242_binding_test.go:40: AC#1 RED: a grant bound to item one was spent with item two's digest - the binding layer did not reject the forged binding`
   - `ticket242_binding_test.go:57: AC#1 RED: empty binding spent a live grant`
   - `ticket242_binding_test.go:171: AC#1 RED: item A's live grant, spent against item B's binding digest, returned none and not misbound. …`
   - `ticket259_denial_rulers_test.go:172: AC#2 RED: a spend against a mismatched binding was classified none, want the binding cause named apart`
   ⇒ 与 §5 记录的 `242-v1` M-A「六枚全绿」**反向**：那枚恒真用例已被 `242-r3` 换成活令牌形（`ticket242_binding_test.go:158-194`
   现量：`q.grantNonce(itA)` ＋ `itA.grants.live()!=1` 前置断言 ＋ 正控 `spend(nonceA2, itA.bind)==denialNone`
   ＋ 鉴别器 `spend("never-issued-value", …)==denialSpentNonce`），恒真那一形今天会被指名红。
2. **真防线（store 按 item 分本＝成员检查）有牙，且牙在路由级**（MU-M 掏空成员扫描）⇒ 9 枚红，其中
   **走路由**的三枚逐字：
   - `queue_test.go:167: err=<nil>，期望 ErrBadGrant`（`Native().Allow` 那条路）
   - `ticket259_denial_rulers_test.go:139: AC#2 RED: the never-issued-proof route returned <nil>, want ErrBadGrant`
   - `ticket259_panel_capability_rulers_test.go:301: AC#3 RED: a value read out of the panel face spent as a native grant and allowed the card; …`
     （同发 `:308` 二次失败＝正控被吞，说明这一发不是空跑）
   ⇒ 成员检查若被拆，路由级立刻红＝票 259 §9 选形 ⓐ 那句"真防线是 store 成员检查＋state 检查"**今天有仪器**。
3. **"拒因指名绑定不对"在路由级仍然只能指到 `spent-or-never-live-nonce`**（结构性事实，本腿自己复跑再判）：
   `queue.go:414` 逐字 `denial := it.grants.spend(nonce, it.bind)`，而 `it.bind` 与 `grantNonce` 写进**同一本**
   store 的是同一个串 ⇒ 经路由的两端恒等；A 的 nonce 不在 B 的本里，扫描落空在 `approval.go:574`
   `return denialSpentNonce`（MU-M 正是把这一行换成 `denialNone` 才让路由级红起来的）。
   `approval.go:477` 那句 `//     property the ErrBadGrant comment states, and ticket 259 does not move it.`
   **本腿复认＝仍成立**（对外确实仍是一句合并 `ErrBadGrant`，`ui.go:130`；且这一形被正向钉住，见 4）。
4. **"对外不许分"这一形有钉**（MU-OUT 把合并句改成只列绑定那一因）⇒
   `ticket259_denial_rulers_test.go:252: AC#2 RED: the merged outward sentence moved ("approval: 原生令牌无效（与本次请求不绑定）"); ticket 259 keeps it merged, the split is internal-only`，
   同发四枚 `Ticket259R1DenialNames*` 全绿＝对内四枚拒因与对外合并是**两枚独立的钉**。

**判语＝成立，但判据字面需要缩。** 票面那句 `Native().Allow(B的corr, A的grant)` ⇒ "拒因指名'绑定不对'"
在现行形状下（票 259 AC#1 选形 ⓐ，`A562`／`A619`）**按字面永远不可能满足**，任何腿都不许为闭合它去改合并文案
（票面射程注原话）。该缩成的字面（交编排者，⛔ 本腿不改票面一字）：

> 判据＝造两枚同时活着、correlation id 不同的卡，把 A 的**未花**令牌递给 `Native().Allow(B的corr, A的grant)`
> ⇒ **要拒**，且**审计面（`GRANT-DENY … denial=`）指名是哪一枚拒因**：经路由的这一发指名
> `denial=spent-or-never-live-nonce`（store 一本一卡，A 的 nonce 不在 B 的本里）；`misbound` 那一枚指名
> **只在 store 级断言**（直递不匹配摘要进 `grantStore.spend`）。**对外维持一句合并的 `ErrBadGrant`，
> 且"不许分"本身有正向钉。** ⛔ 本格不许把"路由能回 `misbound`"当判据（那属形 ⓑ，要动答复面形状才谈，
> 撤销口令「259 改形 ⓑ」）。

⛔ **本腿判不动的一发（具名）**：路由级"两枚活卡"那一发的 **denial 名字**没有任何断言面能读出来——
`internal/agent/approval` 里唯一带两枚活卡走 `Native().Allow` 的断言在 `queue_test.go:161`
`g.Native().Allow(…, "corr-B", pA.Grant)`，而 `corr-B` **从未 push** ⇒ 它拿的是 `ErrUnknownCorrelation`，
不是 grant denial（票面"不是'这条路由没在跑'"那一维在这一枚上是**由缺失的卡**给出的，不是由成员检查给出的）；
真两枚活卡那发只在 `cmd/wisp/…197_test.go:473`，它断的是 `errors.Is(err, ErrBadGrant/ErrPanelAllow)`
（`:252`、`:262` 一带），而该测试**不捕获 gate 的 Logf**——尺＝本腿 `grep -c "GRANT-DENY\|FORGED-OR-STALE" /tmp/t197.txt` ⇒ **0 命中**（同发 rc=0／2 PASS）。
⇒ 缺的读数原文命令＝`PATH="$PWD/third_party/sherpa-onnx:$PATH" go test ./internal/agent/approval/ -count=1 -v`
里**没有**任何一枚用例断言"两枚活卡跨卡 ⇒ `denial=spent-or-never-live-nonce`"；本腿只能按缩后的字面裁，
⛔ 不许按票面原字面判绿，也不许自己造那一枚尺（造尺归实现腿）。

---

## AC#2 出向读面不得带 grant／不得能被答复（反射侧＋能力侧）

**判据原文（票面 `:21` 逐字）**：「反射扫 `approval.PanelItem` 的字段名枚数与名字 ⇒ 出现 `Grant`／可答复类字段
**就红**；再加一发**能力侧**判据（不是词面）：把 `PanelAPI` 的出向读面扩到"能改变卡的状态"那一形也要红。」

`242-v1` 当时实测：M-D2（加 `Permitted bool` ＋把 `"Permitted"` 补进出向读面名单）⇒ 66 全绿；
M-E（`PanelAPI` 加 `Allow(correlationID, grant)` 并实现成真花令牌）⇒ 本包仍 66 PASS。
⇒ 票面自己点名的这两形＝本腿各重跑一发（种在产码，⛔ 不碰测试件名单）：

- **MU-D2a**（`ui.go` PanelItem 新字段 `Permitted bool`）⇒ 2 枚红，逐字：
  - `ticket242_panelface_test.go:39: AC#2 RED: PanelItem.Permitted (bool) is visible to the panel but not declared in panelItemReadFace …`
  - `ticket259_panel_capability_rulers_test.go:174: AC#3 RED: PanelItem.Permitted has type bool - a bool on the card a panel reads IS a verdict bit (allowed / not allowed); … the field name is free, the list is updatable by the same hand that adds the field, so the fence has to be about what the field CAN carry (ticket 259 AC#3)`
  ⇒ **`242-v1` 的"改个字段名＋补名单就绕过词面尺"那一形今天不成立了**：第二枚红句**不看名字只看 Kind**，
  名单同手补不改它判的东西。（第一枚红仍只按名＝它仍可被同手扩名单绕过，见下面"仍全绿的坏形状"。）
- **MU-D2b**（同字段换成 `Permitted int`，即 Kind 栅栏放行的类型）⇒ **只 1 枚红**
  （`ticket242_panelface_test.go:39` 那枚词面尺），`:174` 与 `:199` 与 `:229` 三枚能力尺**全绿**。
- **MU-D2c**（MU-D2b ＋ 把这枚 int 在投影里填成 `it.grants.live()`＝真的把裁决态读到面板读面）⇒ 2 枚红，逐字：
  - `ticket259_panel_capability_rulers_test.go:229: AC#3 RED: PanelItem.Permitted changed when the card's grants were revoked; the panel read face must not reflect grant or decision state (name-blind: this is the capability half of M-D2)`
  - ＋同发的词面尺 `:39`
  ⇒ 逃生口被**第二枚 Kind-盲的尺**堵上：只要该字段真的携带状态，值比对（`reflect.DeepEqual` 前后两投）就红，
  与类型、名字都无关。
- **MU-E**（票面"能力侧"那一形逐字种：`ui.go:171` 接口加 `Allow(correlationID, grant string) error`
  ＋ `gate.go:703` 把它实现成真花令牌 `p.q.allow(corr, grant)`）⇒ 2 枚红／88 PASS（真断言句，非编译失败），逐字：
  - `ticket259_panel_capability_rulers_test.go:73: AC#3 RED: PanelAPI gained method "Allow"; the panel surface is closed at Reject/Head/View, and a fourth name is a new capability on the untrusted face (ticket 242 M-E shape: an allow verb that spends a real grant). …`
  - `ticket259_panel_capability_rulers_test.go:82: AC#3 RED: PanelAPI has 4 methods, the closed set has 3; a count match is part of the fence, otherwise an alias for an existing verb passes silently`
  - `ticket259_panel_capability_rulers_test.go:101: AC#3 RED: the value Gate.Panel() returns (approval.panelAPI) carries method "Allow", which the PanelAPI list does not name; a panel host that asserts a wider interface calls it`
  - `ticket259_panel_capability_rulers_test.go:111: AC#3 RED: the panel carrier satisfies an Allow-shaped verb; a host holding Panel() can settle a card as allowed (ticket 242 M-E)`
  ⇒ 票面那句"能力侧判据（不是词面）"**今天有仪器**：封闭集尺（名＋枚数）＋载体行为尺（Go 类型系统四枚断言，
  `Allow`/`AllowSession`/`Approve`/`Settle` 四种写法各自一条，改名不通吃）。
- **可答复性那一半（语义尺）也在跑**：MU-M 一发里 `:301`／`:308` 两条红句证明
  `TestTicket259R1PanelReadFaceCannotBeLaunderedIntoAnAnswer` 既试"读出来的一切回递给 `Native().Allow`"
  也试回递给 `DecideFromPanel`，并带一条**必须最后跑的正控**（真令牌仍能开卡），⇒ 它绿不是"卡悄悄死了"。

**判语＝成立。** 两半各有本腿亲手种的发与逐字红句；票面点名的两形（加字段／加 allow 方法）
从 `242-v1` 的"全绿"变成今天的 2 红／2 红。⛔ 不许按"尺在文件里"判，本件按的是种进去红没红。

**仍然全绿的坏形状（这行最重要）**：
1. `Permitted int`（或任何 display-scalar：`string`／各宽度 `int*`／`uint*`）**且其值不由 grant／decision 态派生**
   ⇒ 三枚能力尺全绿（MU-D2b 实测：只有词面尺 `:39` 红）。这条形"字段存在但没人填"＝零能力，
   但**它证明词面尺与 Kind 尺的组合仍有缝**：一枚被填了**常量**、日后由别的分支填真值的字段可以先进来。
2. 上面那一形要**完整复现 `242-v1` 的 M-D2**（同手把 `"Permitted"` 补进 `ticket242_panelface_test.go:21`
   的 `panelItemReadFace` 名单 ⇒ `:39` 与 `:44` 两枚词面尺一起绿）——**本腿不能种**：派单硬区＝
   突变只准落 `internal/agent/approval/**` 的非测试产码。⇒ 具名写"这一发本腿未测、缺口在测试件名单同手扩"，
   ⛔ 不许读成"有牙"。
3. **合并文案的三份副本只按子串钉**：对外那半句在 `ui.go:130`、`cmd/wisp/approval_reply.go:229`、`:274`
   有三份（票 259 §9 写的 `:227`／`:272` 已漂到现在这两行）；`cmd/wisp` 侧的钉是子串匹配
   （`cmd/wisp/approval_reply_201_test.go:387` 逐字 `{"yes " + corr, "原生令牌无效"}`）⇒ 从合并句里
   **删掉"缺失/已用"两因、只留"原生令牌无效"前缀**那一形今天能让 `cmd/wisp` 那枚钉不动声色。
   ⛔ 本腿没种这一发（种它要写 `cmd/wisp`，硬区），只按尺的形状登记这条缝。
4. 审计面 `queue.go:424` 那句合并的 `FORGED-OR-STALE …(native grant missing/spent/misbound)`
   与 `:423` 的具名 `GRANT-DENY` 是**两行**：`242-v1` 说的"只有五处注释守着"今天不成立了（`:423` 是被
   `Ticket259R1DenialNames*` 四枚按 `denial=` token 钉的），但**把 `:424` 那半句删掉**今天只有
   `queue_test`／`ticket259_denial_rulers_test.go:127` 两处按 `FORGED-OR-STALE` 子串钉——本腿未种，登记不判。

---

## AC#3 本格不设产码要求：① 是不是同一块石头 ② 甲形够不够

**判据原文（票面 `:28` 逐字）**：「判 AC#1／AC#2 两枚判据**是不是同一块石头**（值不值两格）、以及**甲形够不够**
——若判"够"，本格按"甲落地"翻勾；若判"要乙"，本格不许勾，写清"等 owner 那句解冻"。」

**① 不是同一块石头（本腿复用 `242-v1` 的互不敏感方法自证，读数是本程的）**：
- MU-A／MU-M（绑定族）⇒ 红的全在 `ticket242_binding_test.go`／`ticket259_denial_rulers_test.go`／
  `queue_test.go`／`ticket224_reply_grant_test.go`，**没有一枚 `PanelItem`／`PanelAPI` 尺因它红**（MU-A 那发
  242 panelface 两枚、259 panel 六枚全绿）。
- MU-D2a／b／c／MU-E（出向读面族）⇒ 红的只有 `ticket242_panelface_test.go:39` 与
  `ticket259_panel_capability_rulers_test.go`，**绑定族七枚用例（`TestTicket242*`）全绿**。
⇒ 两族对彼此的错法**互不敏感**＝两个判据、值两格。与 `242-v1` 的判语一致（这次是本腿自己量的，不是转抄）。

**② 甲形够——三枚尺逐枚给尺（都在 `internal/agent/approval` 本包，⛔ 三枚冻结件一字未动）**：

| 尺（票 259 AC#3 的①②③） | 用例名 | 本腿种的那一发 | 红句落点 |
|---|---|---|---|
| ① `PanelAPI` 方法名**封闭集**（名＋枚数）＋载体行为 | `TestTicket259R1PanelAPIMethodSetIsClosed`／`TestTicket259R1PanelCarrierCarriesNoAnswerVerb` | MU-E（真花令牌的 `Allow`） | `:73`／`:82`／`:101`／`:111` |
| ② `PanelItem` 字段**按能力**封名单（Kind 盲于名）＋值侧＋**状态不变性** | `TestTicket259R1PanelItemFieldCapabilityFence`／`…CarriesNoCredentialValue`／`…ProjectionIgnoresGrantState` | MU-D2a（bool）／MU-D2b（int）／MU-D2c（int 填 `grants.live()`） | `:174`（bool）／无（int 常量形）／`:229`（int 带状态） |
| ③ 出向读面**不可回填成答复**（语义） | `TestTicket259R1PanelReadFaceCannotBeLaunderedIntoAnAnswer` | MU-M（成员扫描掏空 ⇒ 读面物竟能开卡） | `:301`＋`:308`（带正控失败，证该发不空跑） |

另附 AC#2 反射侧本就在位的两枚（票 242 自己的甲形，`242-v1` 判过"词面可绕"）：
`TestTicket242PanelItemReadFaceIsFullyDeclared`（`:39`／`:44`，本腿 MU-D2a/b/c 三次都红）＋
`TestTicket242PanelItemStaysGrantFree`（禁词表；本腿三形都不靠它＝它只对"名字带 grant/allow 字样"敏感，
**这一形今天仍全绿**：`Permitted`／`CanProceed` 之类名字它看不见，靠的是 ②）。

**判语＝甲形够** ⇒ 按票面自己那句，本格"按甲落地"可翻（⛔ 本腿不翻，翻框归编排者）；
⛔ **不必解冻乙形**：`internal/panel/l2_grant_boundary_test.go` 本腿零读零写零射程改动，
缺的三枚零件已在 `internal/agent/approval` 本包补齐并各自见红。
⚠ 具名保留一条**不属乙形**的过期口径：票 242 §5 那句"AC#2 不成立"的理由（M-D2 改名即绕过／M-E 零仪器）
今天已被上面四发推翻，编排者翻 AC#3 时该同时把 §5 那句就地打旧（⛔ 本腿不改台件一字）。
