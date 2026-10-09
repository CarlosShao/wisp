# 票 285 · 腿 `285-r1` · AC#3 不放宽既有的钉 ＋ AC#2 只裁形状

锚＝`00-anchor.md`。本文件两格：AC#3＝`242-v2` 具名的四发红句**逐枚复跑**（全部种在盘上，⛔ 零 `-overlay`）；
AC#2＝②③④三形各给"留／补尺／归票"，⛔ 本程一律不动手。

⛔ 零枚既有断言被放宽、零枚 `t.Skip`、⛔ 三枚冻结件（`internal/panel/tokens_fourway_test.go`／
`internal/panel/l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go`）零触碰。
四枚产码件（`approval.go`／`gate.go`／`queue.go`／`ui.go`）在下列每一发之后都还原成起手 hash，
逐字等值（起手表＝`00-anchor.md`，与 `242-v2` 表 `:25-28` 也逐字等值）。

## AC#3 四发红句逐枚复跑

| 发 | 种的内容（本腿 `sed -n 'Np'` 复量到的逐字行） | rc／枚数 | 判 |
|---|---|---|---|
| **MU-A** | `approval.go:569` ⇒ `		if false && !equalSecret(stored, bind) {` | 1／86 PASS／**4 FAIL** | **仍在**，枚数与 242-v2 一致 |
| **MU-M** | `approval.go:574` ⇒ `	return denialNone // …membership scan gutted` | 1／81 PASS／**9 FAIL** | **仍在**，枚数与 242-v2 一致（本腿新文件在场时＝10 枚，见 `10-ac1-ruler.md` 发 4） |
| **MU-D2a** | `ui.go:57` 插入 `	Permitted     bool` | 1／88 PASS／**2 FAIL** | **仍在** |
| **MU-D2c** | `ui.go:57` 插入 `	Permitted     int` ＋ `queue.go:578` 插入 `		Permitted:     it.grants.live(),` | 1／88 PASS／**2 FAIL** | **仍在** |
| **MU-E** | `ui.go:171` 插入 `	Allow(correlationID, grant string) error` ＋ `gate.go:703` 插入 `func (p panelAPI) Allow(corr, grant string) error   { return p.q.allow(corr, grant) }` | 1／88 PASS／**2 FAIL**（4 条红句） | **仍在** |
| **MU-OUT** | `ui.go:130` ⇒ `	ErrBadGrant = errors.New("approval: 原生令牌无效（与本次请求不绑定）")`，定向 `-run 'Ticket259R1Outward\|Ticket259R1Denial'` | 1／4 PASS／**1 FAIL** | **仍在** |

⇒ 票面 `:18` 那句"⛔ 任何一枚变绿＝本票退回"＝**不触发**；四发全部见红，红句逐字如下。

### MU-A 的四条红句（逐字，含 `file:line`；行号＝本腿同一次取数）

```
    ticket242_binding_test.go:40: AC#1 RED: a grant bound to item one was spent with item two's digest - the binding layer did not reject the forged binding
    ticket242_binding_test.go:57: AC#1 RED: empty binding spent a live grant
    ticket242_binding_test.go:171: AC#1 RED: item A's live grant, spent against item B's binding digest, returned none and not misbound. spent-or-never-live-nonce means the comparison never ran because the row is not in A's store (the shape this case had before, zero power over bindings); missing-nonce means nothing was presented at all; none means the cross-card grant was accepted.
    ticket259_denial_rulers_test.go:172: AC#2 RED: a spend against a mismatched binding was classified none, want the binding cause named apart
```
四枚行号（`:40`／`:57`／`:171`／`:172`）与 `242-v2` 的 `:66-70` 记录**逐一对齐＝未漂**。

### MU-M 的九枚红名册＋票面点名的三枚"走路由"红句

红名册（`--- FAIL` 逐行）＝`TestTicket242SpendRejectsForgedBindingAndConsumesTheNonce`／
`TestTicket242ForgedBindingCannotSpendAnotherItemsGrant`／`TestTicket259R1DenialNamesSpentOrNeverLiveNonce`／
`TestTicket259R1DenialNamesMisbound`／`TestTicket259R1OutwardAnswerStaysMergedAcrossDenials`／
`TestTicket259R1PanelReadFaceCannotBeLaunderedIntoAnAnswer`／`TestPanelSourcedAllowIsRejectedOnEveryForgeableAxis`／
`TestGrantIsSingleUseAndBoundToItsItem`／`TestTicket224ForgedSessionAnswerRecordsNothing`（共 9 枚）。

```
    queue_test.go:167: err=<nil>，期望 ErrBadGrant
    ticket259_denial_rulers_test.go:139: AC#2 RED: the never-issued-proof route returned <nil>, want ErrBadGrant
    ticket259_panel_capability_rulers_test.go:301: AC#3 RED: a value read out of the panel face spent as a native grant and allowed the card; the outward read face and the answer face are not the same channel, and this is the case that says so
    ticket259_panel_capability_rulers_test.go:308: AC#3 RED: the panel-observable values were refused, but so was the live native grant (approval: correlation_id 无对应待审批项) - the card stopped being answerable before the control ran, which makes the pass above meaningless
```
三枚行号 `:167`／`:139`／`:301`（＋同发 `:308`）与 `242-v2` 的 `:76-79` **逐一对齐＝未漂**。

### MU-D2a／D2c／MU-E／MU-OUT 的红句

```
ticket242_panelface_test.go:39: AC#2 RED: PanelItem.Permitted (bool) is visible to the panel but not declared in panelItemReadFace (…)[MU-D2a]
ticket259_panel_capability_rulers_test.go:174: AC#3 RED: PanelItem.Permitted has type bool - a bool on the card a panel reads IS a verdict bit (…)[MU-D2a]
ticket242_panelface_test.go:39: AC#2 RED: PanelItem.Permitted (int) is visible to the panel but not declared in panelItemReadFace (…)[MU-D2c]
ticket259_panel_capability_rulers_test.go:229: AC#3 RED: PanelItem.Permitted changed when the card's grants were revoked; the panel read face must not reflect grant or decision state (name-blind: …)[MU-D2c]
ticket259_panel_capability_rulers_test.go:73: AC#3 RED: PanelAPI gained method "Allow"; the panel surface is closed at Reject/Head/View, …[MU-E]
ticket259_panel_capability_rulers_test.go:82: AC#3 RED: PanelAPI has 4 methods, the closed set has 3; …[MU-E]
ticket259_panel_capability_rulers_test.go:101: AC#3 RED: the value Gate.Panel() returns (approval.panelAPI) carries method "Allow", …[MU-E]
ticket259_panel_capability_rulers_test.go:111: AC#3 RED: the panel carrier satisfies an Allow-shaped verb; …[MU-E]
ticket259_denial_rulers_test.go:252: AC#2 RED: the merged outward sentence moved ("approval: 原生令牌无效（与本次请求不绑定）"); ticket 259 keeps it merged, the split is internal-only[MU-OUT]
```
行号 `:39`／`:174`／`:229`／`:73`／`:82`／`:101`／`:111`／`:252` 与 `242-v2` 记录**全部对齐＝未漂**。

### 顺带复认的结构性事实（本腿自己读盘，不是转抄）

- `queue_test.go:161` 现量逐字
  `if err := g.Native().Allow(context.Background(), "corr-B", pA.Grant); !errors.Is(err, approval.ErrUnknownCorrelation) {`
  ⇒ 242-v2 `:103-105` 的"corr-B 从未 push ⇒ 拿的是 `ErrUnknownCorrelation`"**仍成立**，票面 `:8` 的 `:161` **未漂**。
  本腿的 `TestTicket285R1UnknownCorrelationBooksNoDenialLine` 就是把这一形的"零 `GRANT-DENY` 行"钉住的（见 `10-ac1-ruler.md`）。
- `queue.go:414` 现量逐字 `denial := it.grants.spend(nonce, it.bind)`；`approval.go:574` 现量 `	return denialSpentNonce`；
  `approval.go:477` 那句 `//     property the ErrBadGrant comment states, and ticket 259 does not move it.` **仍在**。
  ⇒ 票面 `:12` 的结构性事实（路由两端恒等、`misbound` 只 store 级可达）**复认成立**。
- `ui.go:130` 现量逐字 `	ErrBadGrant = errors.New("approval: 原生令牌无效（缺失/已用/与本次请求不绑定）")`。
- `cmd/wisp/approval_reply.go:229`／`:274` 两份副本现量都以 `原生令牌无效（缺失/已用/与本次请求不绑定）：%w；` 开头
  ⇒ `242-v2` 修正后的 `:229`／`:274` **未漂**（票面 `:11` 原写的 `:227`／`:272` 已过期）。

## AC#2 三形各一句（⛔ 本程不动手，只写形状）

**② `Permitted` 常量形 display-scalar ⇒ 判「补尺，另派」。**
现量：MU-D2b 那一形（`Permitted int`、无人填）今天**只**被词面尺 `ticket242_panelface_test.go:39` 看见；
三枚能力尺（`:174`／`:199`／`:229`）全绿（`242-v2` `:128-129`）。
它变成**真洞**的条件：一枚字段的名字被起成无害名词（`Rows`／`Count`）、值是"当下为常量但日后由别的分支填真值"，
且填进来的那一版**不随 `revoke()` 变化**（`:229` 那枚值不变性尺的判据就是"revoke 前后两投 `DeepEqual`"，
`grants.revoke()` 只清 nonce、不动 `it.state`，所以"在 push 时读一次状态并快照"这一形它看不见）。
形状＝在 `:229` 旁边加**第二枚不可变性轴**（手写 `it.state`／换 `it.bind`／重 `issue` 之后各投一次，全部要求逐字节相同），
⛔ 不预设要建、⛔ 本腿不建。

**③ 名单与字段名同手扩 ⇒ 判「有独立钉的形状，但仍属补尺；归本票（285）后续派单」。**
现量：名单在测试件里手写（`ticket242_panelface_test.go:21` `var panelItemReadFace = []string{…}`，六枚名字），
所以"同手把新名字补进名单"那一发**本包内无解**——任何活在本包测试件里的尺都可被同一只手改。
⛔ 这一形今天仍是**已知空洞**，不是"已验不出"（写腿的硬区＝突变只准落非测试产码，本腿也不能种它）。
独立钉的形状＝把期望集搬到**另一只手才会碰的工件**里：`PanelItem` 目前唯一的包外消费者是
`cmd/wisp/approval_reply.go:397 panelItemLine(item approval.PanelItem)`（本腿 `git grep PanelItem -- '*.go' :!internal/agent/approval`
现量＝只有它，其余命中在 `.scratch/**` 的前腿快照里），而它渲染的字段集是可 golden 的——
⇒ **"名单由产码／金样派生而不是测试件手写"这一发有形状**（`cmd/wisp` 侧那一枚渲染行的金样＝第二只手），
⛔ 只答有没有形状，本腿不建、不碰 `cmd/wisp`。
⚠ 这条形状**踩在禁区块上**：AGENTS.md §1.1＝"golden 一字节都不许动"，**新建**一枚金样属契约变更、须人工批准，
不是写腿能自决的落地方式；把它记成"要 owner 拍板的那一枚"，而不是"下一腿照着做的那一枚"。

**④ `cmd/wisp` 两份副本只被子串钉 ⇒ 判「补尺，另派；归本票（285）后续派单」，按尺形登记、⛔ 本腿未种。**
现量：`cmd/wisp/approval_reply_201_test.go:387` 逐字 `			{"yes " + corr, "原生令牌无效"},` ＝**子串**匹配
（`242-v2` `:160-161` 记的同一行，行号未漂）。⇒ 从合并句删掉"（缺失/已用/…）"只留前缀那一形今天不响。
⛔ 本腿不能种（种它要写 `cmd/wisp`，与 `10` 文件的落点裁量"甲形不抢 `cmd/wisp` 面"一致）。
形状＝把 `cmd/wisp` 侧的钉换成**等值**钉（`got == approval.ErrBadGrant.Error()`，或三份副本两两等值比对），
这样 MU-OUT 那一族在 `cmd/wisp` 侧也见红；⛔ 不预设建在哪条腿。

## 门禁与本腿写面收口

- 本腿只新建：`internal/agent/approval/ticket285_route_denial_name_rulers_test.go`（185 行）＋
  `.scratch/wisp/probes/285/r1/{00-anchor,10-ac1-ruler,20-ac3-no-relaxation}.md`。⛔ 零 `.sh`／`.ps1`／`.txt`／`.out`。
- 本腿只追加：票 285 票面 Progress log 一节「交件读数」，⛔ 未翻勾、⛔ 未改判据句一字。
- `git status --porcelain -- internal cmd` ＝**空**；四枚产码件 hash ＝起手值逐字等值。
