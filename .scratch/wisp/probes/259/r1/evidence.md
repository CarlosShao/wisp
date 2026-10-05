# 票 259 落地腿 `259-r1` — AC#2 拒因可指名 + AC#3 三枚能力尺

本件＝写码腿 `259-r1` 的自证表。派单口径：**只两格**（AC#2 拒因可指名／AC#3 三枚能力尺）；
AC#1（ⓐ具名降级／ⓑ做出牙 二选一）与 AC#4（`cmd/wisp` 载具前置）**本轮不落**，理由与归口见 §6。
写面＝`internal/agent/approval/**` ＋本证据件目录。`cmd/wisp`／`internal/panel`／`internal/tools`
正被 `257-r2`／`181-r3` 取整包红名册，那三枚包本腿连跑都没跑（§5 逐枚列实跑命令）。

---

## §0 起手锚（同发取，2026-10-05 11:3x +08）

| 项 | 读数 |
|---|---|
| `date` | `2026-10-05 11:35:48 +0800` |
| `git log -1 --format=%h` | `d7236abc` |
| 分支 | `dev` |
| `git status --porcelain -- cmd internal tools` | 见下（**起手态，逐字**） |

```
 M internal/panel/git.go
 M internal/panel/git_test.go
 M internal/panel/instructions_200.go
 M internal/panel/workspace.go
 M internal/panel/workspace_test.go
 M internal/tools/paths.go
 M internal/tools/paths_workspace.go
 M internal/tools/paths_workspace_test.go
 M internal/tools/task_pointer_authority_ac3_174r4_windows_test.go
?? internal/panel/workspace_account_181r3_test.go
?? internal/tools/paths_workspace_account_181r3_test.go
```

⚠ 登记不归因：以上全部落在 `internal/panel`／`internal/tools`＝派单点名的**别的腿的脏面**
（`181-r3`）。`internal/agent/approval/` 起手零命中；`cmd/wisp` 起手零命中（`257-r2` 尚未落笔或未进该面）。
本腿⛔ 未 `stash`／未 `checkout`／未 `clean`，一个字节都没清。

起手备份锚（`git cat-file blob HEAD:<path> > <path>.orig`，本腿突变还原的**唯一**比对基准）：见 §4 表末列的逐枚 md5，
每枚文件起手值＝`git cat-file blob HEAD:` 抽出的那份，终值＝交付前的 `md5sum`。

---

## §1 起手名册 ＋ 本腿会撞的既有断言原文（逐枚读过，不是 grep 新符号）

### §1.1 起手逐名读数（⛔ 只跑 `internal/agent/approval/` 这一枚包）

| 命令 | 起／止时刻 | 读数 |
|---|---|---|
| `go test -count=1 ./internal/agent/approval/` | 11:36:10 → 11:36:18 | `ok github.com/CarlosShao/wisp/internal/agent/approval 0.618s`（rc=0） |
| `go test -count=1 -v ./internal/agent/approval/` | 11:36:23 → 11:36:29 | **77 枚 `--- PASS`／0 枚 `--- FAIL`**，`ok 0.611s`；顶层名 **58 枚**（含子测试共 77 行 PASS） |

名册落盘：`.scratch/wisp/probes/259/r1/base-v.txt`（逐字输出）＋`base-names.txt`（58 枚顶层名排序）。
⛔ 起手态**零红名册**——本包全绿，所以本腿交付后如果出现红，只能是本腿种出来的。

对照前人读数：`docs/evidence/s1/242-grant-binding-v1.md` §0 写的是 **66 PASS**（10-03 锚）。
今天 58 枚顶层名／77 行 PASS＝242-v1 之后又落了票 220／260 系那几批用例，**不是同一名册**；
引用时按今天的读数，不要拿 66 当常量。

### §1.2 会撞我新返回形状的既有断言（逐枚原文，逐枚判怎么办）

`spend` 的返回形状从 `bool` 改成带拒因的形状，仓里**只有 6 枚调用者**（尺＝`grep -rn "\.spend("`）：

1. `internal/agent/approval/queue.go:376` 逐字 `spent := it.grants.spend(nonce, it.bind)`——**产码唯一调用点**，本腿要改的就是这一跳。
2. `ticket242_binding_test.go:23` 逐字 `if s.spend("nonce-live", "digest-of-item-two") {` → 红句 `AC#1 RED: a grant bound to item one was spent with item two's digest - the binding layer did not reject the forged binding`
3. `ticket242_binding_test.go:29` 逐字 `if s.spend("nonce-live", "digest-of-item-one") {` → 红句 `AC#1 RED: the rejected nonce was still spendable afterwards - single-use was broken by the forged attempt`
4. `ticket242_binding_test.go:40` 逐字 `if s.spend("nonce-e", "") {` → 红句 `AC#1 RED: empty binding spent a live grant`
5. `ticket242_binding_test.go:50` 逐字 `if !happy.spend("nonce-a", want) {` → 红句 `AC#1 RED: the exact minted digest failed to spend its own grant - the mint/spend round trip is broken`
6. `ticket242_binding_test.go:104` 逐字 `if itA.grants.spend("leaked-or-guessed-nonce", itB.bind) {` → 红句 `AC#1 RED: item B's binding spent item A's grant - cross-item binding is broken`

**处置**：2–6 这五枚按**逐字等价**改判据形状（`spend` 返 `denialNone` ⇔ 旧 `true` ⇔ 旧"允许放行"那一支），
红句文本一字不改、失败条件不放宽。等价式写死在本件 §2。⛔ 不删任何一枚、⛔ 不把 `t.Fatal` 降级成 `t.Log`。

### §1.3 会撞我新审计行的既有钉（读到的原文，逐枚）

- `cmd/wisp/approval_reply_201_test.go:418` 钉逐字 `"approval: FORGED-OR-STALE allow rejected corr=t201-panel-corr"`，
  匹配方式＝`strings.Contains(audit, want) || strings.Contains(sentences, want)`（`:424-427` 那圈）。
- `cmd/wisp/subagent_selfapproval_197_test.go:336` 钉逐字 `"approval: FORGED-OR-STALE allow rejected corr=" + rec.corr1`，
  同包同一枚 `strings.Contains(audit, want)`（`:340-343`），同发还钉 `"approval: ANSWER-ALLOW corr=" + rec.corr1 + " tool=fs.write route=native decision=allow"`。
  ⇒ 结论：`queue.go:379` 那枚 `FORGED-OR-STALE` 行**必须一字不动**（含那句 `missing/spent/misbound` 合并口径），
  新增的指名行只能**另起一行**。枚数/计数式断言在两枚钉里都不存在（是 Contains），本腿另在 §5 复核"没有一枚断言数审计行的条数"。
- 对外文案钉：`cmd/wisp/approval_reply.go:227`／`:272` 逐字 `原生令牌无效（缺失/已用/与本次请求不绑定）`，
  `cmd/wisp/approval_reply_201_test.go:387` 钉 `"原生令牌无效"`；`ui.go:130` 是这三份副本里的本包那一份。
  ⇒ 编排者已裁「对外维持合并」：`ui.go:127-130`（含 `// The message never says which, so the API cannot be used to probe` 那句）
  **一字不改**，`ErrBadGrant` 的文案与错误值**一字不改**。撤销口令＝「259 对外也区分」，本腿没碰。
- `internal/agent/approval/queue_test.go:132`／`:166` 用 `errors.Is(err, approval.ErrBadGrant)` 钉拒答的错误值；
  `ticket87_veto_l2_test.go:119`、`ticket97_alias_direction_test.go:87` 同。⇒ 四种拒因对外**仍然**折成 `ErrBadGrant` 才不红，这正是裁定的形状。
- `internal/agent/approval/fakes_test.go:276-280`：`newGate` 的 `t.Cleanup` 断言 `g.Queue().Depth() != 0` 即报红
  ⇒ 本腿所有走路由的新用例**必须把卡收干净**（答复或拒绝），否则会被这枚既有钉打死。
- `internal/agent/approval/ticket242_panelface_test.go:21-28` 的 `panelItemReadFace` 是 M-D2 那发绕过的那份名单
  （`"Permitted"` 一加就绿）。⇒ 本腿⛔ 不改这枚既有尺的名单与判据（改了就是"把被绕过的钉拔下来换新形"，
  对抗验收会读成洗读数）；新尺另立一枚文件，**按能力**判，与它并存。

### §1.4 已知正控 `197-v1` 的 D2（本腿新尺不许把它做成常红或常绿）

出处（本腿自己检索并读到原文，⛔ 不转述前人结论）：
`docs/evidence/s1/197-ac5-selfapproval-v1.md:75` 那行逐字把 **D2** 定义成
"`approval.go:292 grantStore.spend` 不再比对 binding 摘要，活令牌可花在本机任何一张卡上"，
当时的读数＝**⛔ 全绿，两把尺都绿**（同文件 `:228` 那节标题逐字 "`bindDigest` 那一层今天零尺（D2，跨条发现）"）。
第二次落实是同一种的 `docs/evidence/s1/242-grant-binding-v1.md` §2 的 **M-A2**（读数＝两枚 `TestTicket242Spend*` 红）。
⚠ 那两行行号是他们那次的现量：今天 `spend` 在 `approval.go:420`、比对在 `:431`（本腿改形后在 `:486` 一带）
⇒ **引"形"（花侧不比摘要），不引行号**。

**本腿的处置＝§4 的 PC-D2 一发**：种下同形 ⇒ 指名"绑定不对"那一枚新用例**必须红**（不许常绿）；
不种时本包 89 `--- PASS`／0 `--- FAIL`（不许常红）。两向都有读数。
对照 197-v1 那次"**全绿、零尺**"：同一种今天响 **三枚**（两枚在册 242 钉＋一枚本腿新钉）——
这一格是 259 相对 197-v1 唯一被翻过来的读数。⛔ 但它翻的仍是 **store 面**那一层：
跨卡路由那一句（现量 1／2 的恒等式）本腿一字未改，见 §6 第 1、2 条。
---

## §2 AC#2 改形：逐枚 before／after（一字逐字，含等价式）

### §2.1 `internal/agent/approval/approval.go`（新增形状＋`spend` 换返回）

before（起手态，逐字 15 行）：

```go
// spend validates and consumes. It reports why a rejection happened in
// user-safe terms: never "close, you got 3 bytes right".
func (s *grantStore) spend(nonce, bind string) bool {
	if nonce == "" {
		return false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for v, stored := range s.values {
		if !equalSecret(v, nonce) {
			continue
		}
		delete(s.values, v) // consumed whether or not the binding matched
		return equalSecret(stored, bind)
	}
	return false
}
```

after：`grantDenial`（未导出，`uint8`）五枚值 `denialNone`／`denialMissingNonce`／
`denialSpentNonce`／`denialMisbound`／`denialNotPending` ＋ `label()` 渲染审计词元；
`spend` 返 `grantDenial`，控制流逐支等价：

```go
func (s *grantStore) spend(nonce, bind string) grantDenial {
	if nonce == "" {
		return denialMissingNonce
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for v, stored := range s.values {
		if !equalSecret(v, nonce) {
			continue
		}
		delete(s.values, v) // consumed whether or not the binding matched
		if !equalSecret(stored, bind) {
			return denialMisbound
		}
		return denialNone
	}
	return denialSpentNonce
}
```

**等价式（本腿全部的"没放宽"都挂在这一行上）**：
`旧 true` ⇔ `新 denialNone`；`旧 空值 false` ⇔ `denialMissingNonce`；
`旧 未命中 store 的 false` ⇔ `denialSpentNonce`；`旧 摘要不等的 false` ⇔ `denialMisbound`。
`delete` 的位置、`equalSecret` 的用法、循环结构一字未动 ⇒ 改形**没有**给任何一支添牙，
只让它说得出自己是谁。（现量 1／2 那句恒等式本腿未修，见 §6 第 1 条。）

### §2.2 `internal/agent/approval/queue.go`（`allowScoped` 两处分岔写指名）

before 逐字：

```go
	if it.state != statePending {
		q.mu.Unlock()
		return ErrNotPending
	}
	...
	spent := it.grants.spend(nonce, it.bind)
	q.mu.Unlock()
	if !spent {
		q.logf("approval: FORGED-OR-STALE allow rejected corr=%s (native grant missing/spent/misbound)", corr)
		return ErrBadGrant
	}
	if !q.deliver(it, answer{a: tools.AnswerAllow, why: "用户在原生侧批准了本次操作"}) {
		return ErrNotPending
	}
```

after：`denial := it.grants.spend(nonce, it.bind)`；拒因不为 None 时**先**写
`approval: GRANT-DENY corr=%s tool=%s denial=%s`、**再**原样写那句 `FORGED-OR-STALE`
（合并口径一字未动，含括号里那三个并列词），仍返 `ErrBadGrant`；
状态那道关与 `deliver` 落败那一支各写一行 `denial=card-not-pending`。
⇒ 审计现场从"三条并列词"变成"指名一枚词元"，对外仍是同一枚错误值同一句话。

### §2.3 `ui.go` 零字节

`git status --porcelain -- internal/agent/approval` 交付态只有本腿点名的五枚文件，
`ui.go` 不在其中；它的 md5 起手＝终值＝`a4cee69ef76adc54ea0154527f283604`
（与 `242-v1` §0 那枚锚**同一把值**，跨两枚腿可比）。`ui.go:127-130` 那句
`The message never says which, so the API cannot be used to probe for valid nonces`
与 `ErrBadGrant` 的文案都在原位；撤销口令「259 对外也区分」未被本腿使用。

### §2.4 五枚既有调用点的 before／after 逐对（红句文本零改动）

| # | before | after |
|---|---|---|
| 1 | `if s.spend("nonce-live", "digest-of-item-two") {` | `if d := s.spend("nonce-live", "digest-of-item-two"); d == denialNone {` |
| 2 | `if s.spend("nonce-live", "digest-of-item-one") {` | `if d := s.spend("nonce-live", "digest-of-item-one"); d == denialNone {` |
| 3 | `if s.spend("nonce-e", "") {` | `if d := s.spend("nonce-e", ""); d == denialNone {` |
| 4 | `if !happy.spend("nonce-a", want) {` | `if d := happy.spend("nonce-a", want); d != denialNone {` |
| 5 | `if itA.grants.spend("leaked-or-guessed-nonce", itB.bind) {` | `if d := itA.grants.spend("leaked-or-guessed-nonce", itB.bind); d == denialNone {` |

五枚 `t.Fatal` 的文本逐字未改（PC-D2 那发把 1、3 两枚打回红，读数见 §4）。
⛔ 无一枚删去、无一枚从 `t.Fatal` 降级为 `t.Errorf`／`t.Log`、无一枚 `t.Skip`。

### §2.5 四枚拒因各自的句子（本腿新用例逐字）

| 拒因 | 用例名 | 走哪条面 | 指名红句（逐字） |
|---|---|---|---|
| 缺失 | `TestTicket259R1DenialNamesMissingNonce` | 路由（`Native().Allow`，空值） | `AC#2 RED: the audit does not name the missing-proof cause (%s)\naudit:\n%s` |
| 已用 | `TestTicket259R1DenialNamesSpentOrNeverLiveNonce` | 路由（陌生值）＋ store（同一枚活牌第二次递进来） | `AC#2 RED: the audit does not name the not-live-proof cause (%s)\naudit:\n%s` ／ `AC#2 RED: a spent nonce was classified %s, want the spent-or-never-live cause named apart` |
| 绑定不对 | `TestTicket259R1DenialNamesMisbound` | store（路由今天不可达，§6 第 2 条） | `AC#2 RED: a spend against a mismatched binding was classified %s, want the binding cause named apart` |
| 状态不是 pending | `TestTicket259R1DenialNamesCardNotPending` | 路由（那道防御关，手设态）＋ `deliver` 落败那一支（路由可达） | `AC#2 RED: the audit does not name the settled-card cause (%s)\naudit:\n%s` |
| 反折叠 | `TestTicket259R1FourDenialLabelsArePairwiseDistinct` | 形状 | `AC#2 RED: denials %d and %d share the label %q; four causes must not fold into fewer names` |
| 对外维持合并 | `TestTicket259R1OutwardAnswerStaysMergedAcrossDenials` | 路由＋错误值 | `AC#2 RED: the merged outward sentence moved (%q); ticket 259 keeps it merged, the split is internal-only` ／ `AC#2 RED: denial %q reuses an existing refusal sentence (%q); each cause carries its own words` |

**没复用现有文案**那半：`ErrPanelAllow`／`ErrNotAdmitted`／`ErrChannelUnavailable`／
`ErrNotPending`／`ErrUnknownCorrelation` 五句被拉进一枚 `strings.Contains` 对照圈，
任何一枚拒因词元若借用（或被塞进）那五句之一即红——这就是"两件事各有各的话"的尺，
不是我对文案的一句自述。

---

## §3 AC#3 三枚能力尺的形状（⛔ 只加尺，未解冻任何冻结件）

### ① 方法名封闭集尺（`ticket259_panel_capability_rulers_test.go`）

- `TestTicket259R1PanelAPIMethodSetIsClosed`：`reflect.TypeOf((*PanelAPI)(nil)).Elem()` 的
  方法名集必须**恰好** `{Reject, Head, View}`，双向钉（多一枚红、少一枚也红），再加
  `NumMethod() != 3` 的红——枚数那一问是给"别名换名混进来"留的。
- `TestTicket259R1PanelCarrierCarriesNoAnswerVerb`：接口名单看不见的另一半——
  `Gate.Panel()` 交出去的**具体载体**（`approval.panelAPI`）可能带比接口更多的方法。
  这枚尺两问：载体方法名集必须等于同一份封闭名单；再用 Go 自己的类型系统探四形
  （`Allow(corr,grant)`／`AllowSession(corr,grant)`／`Approve(corr)`／`Settle(corr,allow bool)`），
  任一枚被满足即红。**探针问的是"能不能答这张卡"，不是叫什么名字**。

### ② 字段封名单按能力判（同一枚文件）

- `TestTicket259R1PanelItemFieldCapabilityFence`（KIND 尺）：`r1PanelFieldIsDisplayData`
  按 `reflect.Kind` 判，**一次都没读字段名**——
  `bool`＝判断位、`func`＝动词、`chan`／`ptr`／`unsafe.Pointer`＝活句柄、`interface`＝装得下任何东西、
  `map`＝无界载荷、`array`／`struct`＝嵌套把面撑宽、`slice` 的元素不是 `string`＝原始字节形状（令牌就是字节）。
  白名单只剩：`string`、各宽度整数、`[]string`。
- `TestTicket259R1PanelItemCarriesNoCredentialValue`（VALUE 尺）：拿一枚**真活着、真签发过令牌**的卡，
  从 `Panel().Head()`／`Panel().View()` 读出投影，逐字段取全部可见字符串，
  与**当场那枚活 nonce** 和**当场那枚绑定摘要**比相等／比包含；红句只点字段名，
  ⛔ 不打印值（凭据值不进日志这条规矩）。
- `TestTicket259R1PanelItemProjectionIgnoresGrantState`（不变量尺）：同一枚卡，
  `revokeGrants` 前／后各读一次投影，`DeepEqual` 必须成立——
  一枚会随授权状态变的字段就是在把判断状态往界面上抄。

### ③ 语义侧尺（回填方向）

`TestTicket259R1PanelReadFaceCannotBeLaunderedIntoAnAnswer`：
`r1EverythingThePanelCanRead` 把出向读面的**每一个字符串**（含 `[]string` 的逐元素）收成一份，
然后两份都拿去答真卡——原生路由 `Native().Allow(corr, v)` 与面板路由
`DecideFromPanel({Allow:true, Grant:v, Source:"native"})`——任何一发返回 `nil` 即红。
收尾一发**正向对照**：拿真令牌再去开同一张卡，必须成功；对照失败本身也红
（"上一发的绿"如果只是卡死了，这枚尺不许它蒙过去）。同一张队里两枚卡＝`r1CardIn`，
两枚各自 `New` 一枚 Gate 不是同一条路由（这一条我起手就写错了，§8 第 1 条）。

**⛔ 三枚尺都没新增 C17 入向方法名、没解冻 `internal/panel/tokens_fourway_test.go`／
`l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go` 任一枚、没碰 D43 转移表。**
尺对"新增第四个名字"的答复是红，不是把门调松——那属契约面，`Q-76` 还挂在机主那儿。

---

## §4 突变名册：红句逐字／还原凭据／窗口起止（到秒）

纪律：还原源＝**突变前**从 `git cat-file blob HEAD:<path>` 抽的副本
（`.scratch/wisp/probes/259/r1/backup/*.orig`，五枚），每发跑完立刻 `git cat-file blob HEAD:` 覆盖回去并 `md5sum`；
每一发都量到"终值＝起手值"，⛔ 没有一发是"改回来再跑一遍"当证明的。
⛔ 全程只跑 `./internal/agent/approval/` 这一枚包（那三枚别人在跑的包一次都没跑）。

起手＝终值 md5（＝HEAD `b6b1d6a4` 那笔里各文件的 blob）：

```
580bddb9b7991c8bfccdb3ef52c73136  approval.go
eace17e8a8b51e3507fb85481a06644e  queue.go
a4cee69ef76adc54ea0154527f283604  ui.go            （与 242-v1 §0 同值，本腿未改）
8d93476e34e5a884b6268d9e64626618  gate.go          （与 242-v1 §2 M-E 还原锚同值）
dc99738456f72cdb0fabad1cb81e99bf  ticket242_panelface_test.go（同上）
```

| 编号 | 种什么（逐字形） | 哪枚必须红 | 实测红句逐字 | 同发读数 | 窗口（起→终） |
|---|---|---|---|---|---|
| **M-E**（票面点名的必抓形） | `ui.go` `PanelAPI` 加 `Allow(correlationID, grant string) error` ＋ `gate.go` `func (p panelAPI) Allow(corr, grant string) error { return p.q.allow(corr, grant) }`（真花令牌、真改卡状态） | 尺① | `ticket259_panel_capability_rulers_test.go:73: AC#3 RED: PanelAPI gained method "Allow"; the panel surface is closed at Reject/Head/View, and a fourth name is a new capability on the untrusted face (ticket 242 M-E shape: an allow verb that spends a real grant). Refusing a card is safe; settling one as allowed is not, and no spelling of a new verb makes it the same act` ／ `:82: AC#3 RED: PanelAPI has 4 methods, the closed set has 3; a count match is part of the fence, otherwise an alias for an existing verb passes silently` ／ `:101: AC#3 RED: the value Gate.Panel() returns (approval.panelAPI) carries method "Allow", which the PanelAPI list does not name; a panel host that asserts a wider interface calls it` ／ `:111: AC#3 RED: the panel carrier satisfies an Allow-shaped verb; a host holding Panel() can settle a card as allowed (ticket 242 M-E)` | **2 FAIL／87 PASS**；同发 242 那六枚逐枚 `--- PASS`（尺①两枚＋四枚 binding 全在册；逐枚行在 `m-e.txt`，计数尺＝`grep -cE "^--- PASS: TestTicket242" m-e.txt`＝6）＝242-v1 的"本包 66 全绿"这一形被复现，只是这次新尺响了 | 起 **11:52:59**（打卡；`ui.go` 那半枚在同一次窗口的上一步落盘，**未单独打卡**＝本表唯一的时刻缺口，见 §6 第 8 条）→ 跑毕 11:53:23 → 终 **11:53:31** |
| **M-D2**（票面点名的必抓形） | `ui.go` `PanelItem` 末尾加 `Permitted bool` **且** `ticket242_panelface_test.go` 把 `"Permitted"` 补进 `panelItemReadFace`（同一只手补名单的真实场景） | 尺②（KIND 那问） | `ticket259_panel_capability_rulers_test.go:174: AC#3 RED: PanelItem.Permitted has type bool - a bool on the card a panel reads IS a verdict bit (allowed / not allowed); the panel read face is plain display data and has no decisions to report. This is the M-D2 shape: the field name is free, the list is updatable by the same hand that adds the field, so the fence has to be about what the field CAN carry (ticket 259 AC#3)` | **1 FAIL／88 PASS**；两枚 242 词面尺逐枚 `--- PASS`＝"改名＋补名单静默过审"这一形**仍然静默过审**，被抓到的是新尺；VALUE／不变量两问不响（bool 不是串）＝三问互相独立，不是同一块石头 | 起 **11:53:39** → 跑毕 11:54:08 → 终 **11:54:15** |
| **M-G**（尺③的必响形） | `approval.go` 加 `func (s *grantStore) peekAny() string`（回一枚活 nonce）＋ `queue.go` `viewLocked` 把 `Reason: it.Dec.Reason` 换成 `Reason: it.grants.peekAny()`——**不加字段、不改名、不改类型**，凭据从既有已声明的展示字段里漏出去 | 尺③（＋尺②的 VALUE／不变量两问） | `ticket259_panel_capability_rulers_test.go:301: AC#3 RED: a value read out of the panel face spent as a native grant and allowed the card; the outward read face and the answer face are not the same channel, and this is the case that says so` ／ `:308: AC#3 RED: the panel-observable values were refused, but so was the live native grant (approval: correlation_id 无对应待审批项) - the card stopped being answerable before the control ran, which makes the pass above meaningless` ／ `:199: AC#3 RED: Head renders a live native grant in PanelItem.Reason; the value a panel can read is the value that opens the card` ／ 同 `:199` 第二行 `View renders ...` ／ `:229: AC#3 RED: PanelItem.Reason changed when the card's grants were revoked; the panel read face must not reflect grant or decision state (name-blind: this is the capability half of M-D2)` | **4 FAIL／85 PASS**；第四枚红是**在册旧尺**：`ticket146_liveapprovals_backing_test.go:349: AC#2 RED: PanelItem 非引用字段也被牵动：tool="shell.run" reason=""`——它钉的是"投影的 Reason 必须等于 Decision 的 Reason"，⛔ 只钉那一枚字段名；我的②／③按值与按方向判，换个字段名（`CorrelationID`、或将来的新字段）它就不响、这两问照响 ⇒ 尺③不是那枚在册钉的换名版 | 起 **11:54:15** → 跑毕 11:54:46 → 终 **11:54:54** |
| **PC-D2**（197-v1 的 D2 正控；242-v1 §2 M-A2 同形） | `approval.go` `spend` 里 `if !equalSecret(stored, bind) { return denialMisbound }` 换成 `if !equalSecret(stored, bind) || len(bind) >= 0 { return denialNone }`＝**花侧不比摘要**，用法保留、编译通过 | 本腿那枚指名"绑定不对"的新用例**必须红**（不许常绿）；242 那两枚在册钉也必须红 | `ticket259_denial_rulers_test.go:172: AC#2 RED: a spend against a mismatched binding was classified none, want the binding cause named apart` ／ `ticket242_binding_test.go:31: AC#1 RED: a grant bound to item one was spent with item two's digest - the binding layer did not reject the forged binding` ／ `ticket242_binding_test.go:48: AC#1 RED: empty binding spent a live grant` | **3 FAIL／86 PASS**。反向一半：不种任何突变时本包 89 PASS／0 FAIL（§5 逐名）＝新尺**不常红**；`TestTicket242ForgedBindingCannotSpendAnotherItemsGrant` 在 PC-D2 下仍 `--- PASS`＝242-v1 那条"名字叫跨卡那枚不敏感"的读数不变，本腿没替它加牙（那属 AC#1） | 起 **11:55:18** → 跑毕 11:55:37 → 终 **11:55:55** |
| **M-H**（AC#2 审计侧的必响形） | `queue.go` 拒因那一行的 `denial.label()` 换成常量 `denialMissingNonce.label()`＝**指名退回并列词**、`FORGED-OR-STALE` 那枚合并句留在原地 | 走路由的两枚拒因用例必须红（缺失那一枚按定义仍绿） | `ticket259_denial_rulers_test.go:143: AC#2 RED: the audit does not name the not-live-proof cause (approval: GRANT-DENY corr=corr-259-spent tool=shell.run denial=spent-or-never-live-nonce)` ／ `ticket259_denial_rulers_test.go:268: AC#2 RED: the audit stopped naming the two causes it merged outward` | **2 FAIL／87 PASS** ⇒ "四句不许折成一句"这一格有牙不是自述：把两枚词元折成同一枚，红两枚用例。⛔ 它没种在 `deliver` 那一支（要竞态载具、会把窗口挂长），那一支的指名读数由 hand-set 载具那枚用例承担——这条差别我没糊成"两支都被 M-H 抓到" | 起 **11:55:55** → 跑毕 11:56:34 → 终 **11:56:34** |

| **M-I**（反折叠那枚尺的必响形） | `approval.go` 的 `label()` 里 `case denialSpentNonce:` 从 `return "spent-or-never-live-nonce"` 改成 `return "missing-nonce"`＝**两枚词元在渲染器里折成一枚**（M-H 折的是调用点、这一发折的是词元本身） | 反折叠尺＋走路由的"已用"那一枚＋维持对外合并那一枚 | `ticket259_denial_rulers_test.go:235: AC#2 RED: denials 1 and 2 share the label "missing-nonce"; four causes must not fold into fewer names` ／ `:240: AC#2 RED: 4 causes produced 3 distinct labels; the four sentences are the deliverable` ／ `:143: AC#2 RED: the audit does not name the not-live-proof cause (approval: GRANT-DENY corr=corr-259-spent tool=shell.run denial=spent-or-never-live-nonce)` ／ `:268: AC#2 RED: the audit stopped naming the two causes it merged outward` | **3 FAIL／86 PASS**；缺失那一枚按定义仍绿。⇒ §8 第 6 条那句"反折叠尺的牙由哪一发量出来"就此从"待量"变成有读数，⛔ 不再挂在 M-H 身上 | 起 **12:04:56** → 跑毕 12:05:15 → 终 **12:05:16** |

**盘上无残留突变**（五枚被改文件逐枚 md5 回到起手值——六发突变压在同样这五枚上，见上表末段；另尺＝
`grep -rn "259-r1 M-E probe|259-r1 M-D2 probe|259-r1 M-G probe|259-r1 PC-D2 probe|259-r1 M-H probe|259-r1 M-I probe|peekAny" --include=*.go internal/agent/approval/`
在 11:56:42 与 12:05:16 各量一次 **0 命中**）。
全部突变输出留在
`.scratch/wisp/probes/259/r1/{m-e,m-d2,m-g,pc-d2,m-h,m-i}.txt`（`-count=1 -v` 逐字），
六发的同发对照读数逐枚复核过并写进上表（例：PC-D2 下发 `TestTicket242ForgedBindingCannotSpendAnotherItemsGrant`
逐枚 `--- PASS`＝242-v1 那条"名字叫跨卡那枚不敏感"未被本腿改动；M-H 下发
`TestTicket259R1FourDenialLabelsArePairwiseDistinct` 逐枚 `--- PASS`＝那枚反折叠尺的牙不由 M-H 作证、由 M-I 作证）。


---

## §5 门禁读数（带时刻，本腿自己跑的）

| 门 | 命令 | 起→止 | 读数 |
|---|---|---|---|
| D22 静态门 | `sh scripts/d22scan.sh` | 11:57:25 → 11:57:57 | rc=0，`clean - no D22 ban violations`；`examined 266 production Go files under internal/ and cmd/`；ban #8 覆盖 `internal/ 512 Go files, comments and _test.go included`／`cmd/ 103` ⇒ **本腿两枚新 `_test.go` 在射程内且零命中**（含 emoji 那把） |
| 同上，SKIP 面 | 同一次输出 | — | 一行 SKIP：`skipped as git-ignored: 1 file(s) under 1 ignored director(ies) [frontend/dist/assets/]`。按"SKIP 不当通过"的口径登记：它是 `frontend/.gitignore` 决定的**既有仓库状态**（同句在 242-v1 §2 与前人多张表里逐字在册），⛔ 不在本腿写面、不归本票 |
| 路径长度帽 | `sh scripts/check-path-length-budget.sh --with-self-test` | 11:58:05 → 11:58:11 | rc=0，`positive control PASSED`；`denominator: tracked paths=5839 over-budget=57 covered by roster=57 not in roster=0`；`VERDICT GREEN` ⇒ 本腿新增两枚文件（52／54 字符名）未进名单、也未撞帽 |
| `go vet` | `go vet ./internal/agent/approval/` | 11:58:11 → 11:58:12 | rc=0，零输出。（起手 11:45:22 那发同样 rc=0） |
| gofumpt | `"$(go env GOPATH)/bin/gofumpt.exe" -l <八枚点名文件>` | 11:57:10 → 11:57:10；**终值复跑 12:05:58 → 12:06:26** | `-l` **空清单**（rc=0）。八枚＝本腿碰过的五枚产码／在册件＋两枚新尺＋`ticket242_panelface_test.go`。**尺确实在跑的负控**：12:04:47 在同一次调用里对一枚故意写坏的样本跑了同一把尺——`gofumpt -l .scratch/wisp/probes/259/r1/gofumpt-negctl/probe.go` 逐字点出该文件（`v0.12.0 (go1.27.1)`，`go env GOPATH`＝`D:\work\base\gopath`，⛔ 没用 `~/GOPATH/bin` 那把瞎尺）⇒ 空清单是"这些文件确实合式"，不是"命令没跑到" |
| 本包全量 | `go test -count=1 -v ./internal/agent/approval/` | 11:56:42 → 11:56:44（五发突变还原后）；**六发全还原后复跑 12:05:42 → 12:05:44** | **89 `--- PASS`／0 `--- FAIL`**，`ok ... 0.435s`；终值输出＝`final-v.txt`，逐枚 md5 名册＝`final-md5.txt` |
| 四门终值复跑 | `sh scripts/d22scan.sh`／`sh scripts/check-path-length-budget.sh --with-self-test`／`go vet`／`gofumpt -l` | 12:05:58 → 12:06:26 | 四把**全部 rc=0**：d22scan `clean - no D22 ban violations`（`internal/` 512 枚含 `_test.go` 在 ban #8 射程内）；路径帽 `VERDICT GREEN`（`over-budget=57`、`not in roster=0`）；`vet` 零输出；`gofumpt -l` 空清单 |

⛔ 本腿一次都没跑 `cmd/wisp`／`internal/panel`／`internal/tools`（那三枚包的整包红名册此刻属 `257-r2`／`181-r3`）。
所有空读数都配了正证：`-l` 空清单同发由 rc=0 与"文件确实存在"两问撑住；grep 0 命中同发由"突变词元确实出现在 `.txt` 输出里"撑住。

---

## §6 判不动的格／量不到的读数（具名归口，⛔ 本腿不替谁选边）

1. **AC#1 的两支（ⓐ具名降级／ⓑ做出牙）本腿两支都没落。**
   现量 1／2 那句恒等式（`queue.go:376` 拿 `it.bind` 去比 `approval.go` 里存进 store 的同一个 `it.bind`）**原样在盘上**——
   本腿只把该比的分支改了名，没换任何一端的实参。
   结构上分不分得开，具名到行：`queue.go:376`（本腿改后为 `denial := it.grants.spend(nonce, it.bind)`）与
   `approval.go` 的 `issue`（`s.values[nonce] = bind`）之间那一跳要**换成答复侧递进来的原料**才谈得上ⓑ，
   而答复面只有 `Request{CorrelationID, Allow, Grant, Reason, Source}`（`gate.go:715-721`）与
   `Allow(ctx, corr, grant)`（`ui.go:146`）——ⓑ需要动的是**入向形状**，不是本腿能单方面挑的形。
   ⇒ 归口：编排者落一枚具名 `A##` 选边；本腿⛔ 未删那枚恒等比对、⛔ 未把"绑定层能挡跨卡"这句话写进任何注释。
2. **绑定不对这一枚拒因今天从路由不可达**（现量 1-3 的直接后果）。
   所以 §2.5 那格诚实的形状是：**三枚走路由（缺失／已用／状态不是 pending）＋一枚走 store（绑定不对）**，
   ⛔ 不是"四枚全走路由"。要它走路由可达＝上面第 1 条的ⓑ。
3. **状态不是 pending 这一枚的两道分岔可达性不同**，本腿分开钉、没混成一句：
   `allowScoped` 开头那道 `it.state != statePending` 的守卫，两个 `it.state` 写点都在
   `dropLocked`（它 `delete(q.byID, ...)`）的同一临界区内 ⇒ 经路由走不到它的不等分支，
   我的用例用 hand-set 态把那道守卫本身钉住；`deliver` 落败那一支是路由可达的（并发答复的败方），
   同一条 `denial=card-not-pending` 词元在那里也写。
4. **"已用"与"从未签发"在 store 里折成一枚词元**（`spent-or-never-live-nonce`）。
   分得开要在 `spend` 消费后仍留一份死牌名册（新的状态面＋无界增长），
   AC#2 的判据（四种各自可断言）不要求它，本腿没擅自加，登记在此供 AC#1 选形时一并看。
5. **跨包读数一律未量**（派单禁跑那三枚包）：
   ① `M-E` 那一形若真落在盘上，`cmd/wisp`／`internal/panel` 里任何 `PanelAPI` 的包外实现者会**编译失败**——
   本腿没跑那两枚包，所以"尺①抓到 M-E"这句的**射程只到本包**（引法同 `242-v1` §5.8）；
   ② 新增的 `GRANT-DENY` 审计行会不会撞上 `cmd/wisp` 里某枚"数审计行条数"的尺，未量
   （本包内我逐枚看过：`FORGED-OR-STALE` 与那五枚钉全是 `strings.Contains` 形状，见 §1.3，无计数式断言）；
   ③ `internal/panel`／`cmd/wisp` 是否已另有针对面板侧 allow 入口的能力尺，未量。
   ⇒ 这三格请 259-v1（或 253／257 系腿）在 `cmd/wisp` 空窗时补一发。
6. **三枚冻结件的射程判断（具名，非内容引用；本腿没读它们、也没跑它们）**：
   `internal/panel/tokens_fourway_test.go`＝C21 token 表四方对账，射程不含 grant 拒因与 `PanelAPI` 方法集；
   `internal/panel/l2_grant_boundary_test.go`＝面板侧不得允许那道界，与本腿尺①**同题不同包**——
   若将来解冻它，两枚方法集尺要同批复看（谁也不许被对方"我已经管了"糊过去）；
   `internal/perm/ticket90_persist_test.go`＝perm 落盘持久化，射程不含 approval 的 grant 层。
7. **AC#4 够不着**：`cmd/wisp/subagent_selfapproval_197_test.go:109` 那枚 `TaskID == CorrelationID`
   不在本腿写面（`cmd/wisp` 此刻是 `257-r2` 的整包面），本腿⛔ 未跑该包、未读那枚行的新值，
   框也⛔ 未勾 ⇒ 具名归口：**AC#4 仍挂在本票面上，等 `cmd/wisp` 空窗另派**。
   本腿的四枚拒因用例不依赖那枚载具（它们要的是"同一张队里两枚卡"，`r1CardIn` 已给）。
8. **M-E 那发的窗口起点有一枚时刻缺口**（如实登记，不补一个像样的数进去）：那发是两枚文件、
   两次落盘，我只在**第二枚**（`gate.go`）落盘后打了 `11:52:59` 这一发；`ui.go` 那半枚落在同一次窗口的
   上一步工具调用里、**没单独打卡** ⇒ 真起点在 `11:52:4x` 与 `11:52:59` 之间，界上界是 11:53:31。
   其余五发的起止都是同一把命令里打的、无缺口。事后若要精确窗口，只有重跑那一发（本腿⛔ 没为一枚
   时刻数挂第二次突变）。
9. **`ui.go` 起手＝交付＝`a4cee69ef76adc54ea0154527f283604`**，与 `242-v1` §0 那枚锚同值 ⇒
   "M-D2 那形今天仍只被新尺抓到"这句可以跨两枚腿对表；⛔ 但同一把值也说明本腿对它做的只有临时突变、
   一次都没留在盘上（六枚 `.txt` 之外无残留）。

---

## §7 commit 名册（只 commit，⛔ 未 push；每笔都带显式 pathspec 写在 commit 命令上）

| 笔 | 内容 |
|---|---|
| `ea522444` | 证据件 §0／§1 ＋ `base-v.txt`／`base-names.txt` |
| `b6b1d6a4` | AC#2 改形（`approval.go`／`queue.go`）＋ 五枚调用点等价改判据（`ticket242_binding_test.go`）＋ 两枚新尺文件（12 枚用例）＋ `after-code-v.txt`／`after-code-names.txt` |
| `50e0299c` | 证据件 §2—§9 写满 ＋ 六枚突变输出（`m-e`／`m-d2`／`m-g`／`pc-d2`／`m-h`／`m-i`）＋ 五枚还原源 `backup/*.orig` ＋ `final-v.txt`／`final-md5.txt` ＋ gofumpt 负控样本 |
| `e8a1bec3` | §1.4／§4／§6／§7 的四处准确性回头改（M-E 窗口起点缺口登记进 §6 第 8 条、"六枚文件"改回"五枚文件六发突变"、D2 出处改引 `197-ac5-selfapproval-v1.md:75` 原文、反折叠尺的牙从 M-H 改记到 M-I） |
| 本笔 | §7 名册补全（含本笔自身的点名，写在下一格里——一枚 commit 装不进自己的哈希，⛔ 我没有为了凑一个数去 amend）＋ §8 第 9 条（`msg-s3.txt` 正文里那枚 `11:52:47` 是未打卡的估计值，读数以下表为准）＋ 交件面尺一行 |

`git add -- <点名文件> && git commit -F .scratch/wisp/probes/259/r1/msg-sN.txt -- <同一批点名文件>`，
中间不停顿；⛔ 无 `add -A`／`.`、无裸 `git commit`、无 `--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；
别人的脏面（§0 那十一枚 `internal/panel`／`internal/tools`）一个字节没动。
写面自证（本腿逐枚回看，含本笔）：`git show --name-only` 对上表每一笔回出的文件清单
**只含** `internal/agent/approval/**` 与 `.scratch/wisp/probes/259/r1/**` 两棵子树 ⇒
票面 AC#5（越界检查）那七条路径一枚都没出现；三枚冻结件⛔ 未读未跑未改（射程判断在 §6 第 6 条）；
票面与 `docs/reports/pending-and-issues.md` ⛔ 一个字节未改（AC 框一枚未勾，翻勾归编排者）。

临时件只建不删（本腿收尾不清理）：`.scratch/wisp/probes/259/r1/` 下
`base-v.txt`／`base-names.txt`／`after-code-v.txt`／`after-code-names.txt`／`final-v.txt`／`final-md5.txt`／
`m-e.txt`／`m-d2.txt`／`m-g.txt`／`pc-d2.txt`／`m-h.txt`／`m-i.txt`（12 枚读数）＋
`msg-s1.txt`—`msg-s5.txt`（5 枚 commit 正文，逐枚＝上表那一笔的 `-F` 源）＋ `backup/` 五枚还原源 ＋ `gofumpt-negctl/probe.go` 一枚负控样本。

**交件面尺（编排者收件可直接对表）**：本件在 12:14:14（`c854a8e0`）之后由本腿最后一次 `wc` 量得
**475 行／44,447 字节**；⛔ 本腿⛔ 不再为凑一个整数改本件，验收腿若要往这里追加段落，
那枚数会随之一动——所以真正当凭据的是下面三条不变量，不是那两枚数。
占位尺 `grep -nE "待填|填写中|未判|placeholder|TBD|TODO|XXX" .scratch/wisp/probes/259/r1/evidence.md`＝**0 命中**
（12:14 之后复量同一条尺）；
本包终值 `go test -count=1 -v`（12:05:42→12:05:44）＝**89 `--- PASS`／0 `--- FAIL`**，
其中本腿新立的 12 枚逐枚 `--- PASS`（名册在 `final-v.txt`）；
起手 58 枚顶层名／77 行 PASS → 交付 70 枚／89 行，`diff base-names.txt after-code-names.txt` **只增不减**（0 枚删除行）。


---

## §8 自我对抗（本腿真改掉的东西，不是清单装饰）

1. **两枚 fixture 各自 `New` 了一枚 Gate ⇒ 两枚卡根本不在同一条路由上。**
   尺③起手就是错的载具，读数把它打回了：
   `ticket259_panel_capability_rulers_test.go:290: AC#3 fixture: card corr-259-launder-y is not readable while pending`（11:50:19 那一发）。
   改成 `r1Card` 只造队首、`r1CardIn` 往**同一枚** `Queue` 里加第二枚，才真的是"一张队、两张卡"。
   ⛔ 我没把它糊成"面板面读不到"然后跳过。
2. **我在 `TestTicket242SpendRequiresTheExactBinding` 那枚红句里顺手加了 `(denial=%s)`。**
   这与我自己在同一次提交里写的"红句文本一字未改"正面矛盾，也动了别人那枚在册钉的文本 ⇒ 撤回头版文本，
   新增读数只在**我自己**的用例里加。
3. **`r1Card` 的 cleanup 原来断言整枚队列 `Depth()==0`。**
   一张队两枚卡时它会互相误报（先收干净的那枚被别人那张卡打死）⇒ 改成逐枚 `it.state`，
   并留"不许留 pending"这一问本身。
4. **`r1FieldStrings` 原来带一枚没用到的形参 `fieldName`。**
   在一枚"按能力不按名字"的尺里留一个字段名形参，正是这格被判不成立的老形状（且 `reflect.StructField.Name` 一路带到调用点）
   ⇒ 删形参，两枚值尺的调用点不再持有任何字段名信息。
5. **我把"197-v1 的 D2"与"242-v1 的 M-A2"混着引了一次。**
   前者是**票面的正控编号**、后者是**另一条腿对同一种的实跑表**，且我最初按 242-v1 的旧行号（`:292`／`:303`）引用；
   今天 `spend` 在 `approval.go:420`、比对在 `:431`（改形后为 `:488` 附近），
   ⇒ §1.4 与 §4 都改成"形"的引法（不比摘要）＋逐字实跑红句，不引旧行号当现值。
6. **反折叠那枚尺我最初写成 `for _, name := range []string{...}` 并把下标当名字用**——
   那样它恒绿（下标永不重复），恰好是"看起来有牙、反形全绿"那一形。换成按 `grantDenial` 枚举＋`seen` 映射，
   并补一发 **M-I**（在 `label()` 里把两枚词元折成同一枚）量它的牙：**3 FAIL**，
   红句逐字在 §4 的 M-I 行（`:235` 那行直接说"denials 1 and 2 share the label"）。
   ⛔ 我原来打算让 M-H 顶这一格，那是错的——M-H 折的是**调用点**、这枚尺问的是**词元之间**，
   M-H 那一发它确实没响（11:56:34 读数只有另两枚红）。这一格从"我以为它有牙"变成"M-I 量过"。
7. **走路由的用例不吹成"整条路由"。**
   `PendingApproval` 那一跳（铸造＋显示＋等答复）没进本腿用例：进就得要么起 goroutine（`_test.go` 里那把
   `bare go func(` 的尺要过）、要么用真 300s 钟。我用了 `push`＋`grantNonce`＋导出答复面这条**同一段产码**的路，
   并把这一格留在 §8 第 7 条而不是在 §9 判语里写"显示侧已钉住"。
8. **窗口长度我量过自己的最差一发**：最长的一发 M-G 是 11:54:15→11:54:54＝**39 秒**（其中测试跑 31 秒），
   ⛔ 没有任何一发挂着突变去跑别的东西；每发跑完**同一条命令里**就还原＋md5 回锚。
9. **我在 `msg-s3.txt`（commit `50e0299c` 的正文）里写了一枚 `11:52:47`——那不是读数。**
   M-E 那发是两枚文件、我只在第二枚落盘后打了 `11:52:59` 一发，`ui.go` 那半没单独打卡；
   commit 正文已入库（⛔ 禁 `--amend`，共享工作树里更不改写），所以这一格的处理是：
   **在 §4 表里把该发的起点改写成"界内估计"、把缺口具名进 §6 第 8 条、在这里留一笔账**，
   而不是回头把那条 msg 改成一个看起来量过的数。同一格另有一处同族失误：
   `msg-s3.txt` 首行原本写"七枚突变输出"、正文写"六枚文件逐枚 md5"（突变六发、被改文件五枚），
   首行我在提交前一次编辑里已改成"六枚"，正文那处落在 `e8a1bec3`——
   两枚都是"把发数当文件数"的同一枚口误，读数本身没受影响（`wc -l` 与 md5 名册都是实跑）。

---

## §9 交件判语（本腿只对自己那两格负责）

- **AC#2（拒因可指名）＝交付，且实测有牙。**
  形状：内部未导出的 `grantDenial` ＋审计行 `GRANT-DENY ... denial=<词元>`；对外一字未动
  （`ErrBadGrant` 那句合并文案、`ui.go:127-130` 那句"故意不区分"、`FORGED-OR-STALE` 那整句都在原位，
  `ui.go` 的 md5 跨两枚腿可比）。四枚拒因各自的指名红句在 §2.5；
  "四句不许折成一句"由两发分别量到：**M-H**（调用点折成常量 ⇒ 两枚用例红）与
  **M-I**（词元本身折成同一枚 ⇒ 三枚用例红，含反折叠那枚尺），
  "不许复用别层文案"由 `borrowed` 那一圈量到。
  诚实的限制：绑定不对那一枚今天只能走 store（§6 第 2 条），三枚走路由。
- **AC#3（三枚能力尺）＝三枚都装上，且各有一发必响。**
  尺①抓 M-E（四行逐字红句，含"载体带第四枚方法"与"满足 Allow 形接口"两问）；
  尺②抓 M-D2（同发两枚 242 词面尺仍绿＝这格上次被判不成立的根因被正面覆盖）；
  尺③抓 M-G（**一枚真令牌从出向读面回填、真把卡打开了**，那一发同时让②的 VALUE／不变量两问响）。
  三枚尺都⛔ 不读字段名／方法名之外的任何字符串味道，⛔ 未解冻冻结件，⛔ 未新增 C17 入向方法名。
- **AC#1／AC#4＝未动，具名归口在 §6 第 1、7 条**（ⓐ／ⓑ 两支都没落，本腿连"哪一支更省"都没写）。
- **名册**：本包顶层名 58 → 70（只增不减，diff 无删除行），PASS 行 77 → 89，起手与交付都是 **0 FAIL**。
- **四门**：d22scan clean（且新 `_test.go` 在 ban #8 射程内）／路径帽 GREEN／`go vet` rc=0／gofumpt `-l` 空。
- **本腿没资格读的格子**：票面 AC 框一枚未勾（翻勾归编排者）；
  `docs/reports/pending-and-issues.md` 一个字节未改；`cmd/wisp` 与 `internal/panel`／`internal/tools` 的颜色本腿**没跑也没推断**。

