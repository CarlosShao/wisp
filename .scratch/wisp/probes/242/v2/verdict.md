# 票 242 v2 对抗验收裁决（非实现者腿）

被审对象：票 242 判据①甲形 `internal/agent/approval/ticket242_binding_test.go`（实现者＝编排者本人，代笔提交 `4bf7e683`）
＋ 判据②甲形 `ticket242_panelface_test.go`（`2048d6b6`／242-r2）。
本腿起手锚 `0c9726f9`，收工锚 `a0455bb1`（跑中途 HEAD 前进了 7 枚，见 Q1）。
⛔ 本腿不翻票面勾、不动 `internal/**`；全部错法在**仓外** `$HOME/242v2-mut/` 用 `go test -overlay` 造。
票面尺复认：`grep -cE '^[[:space]]*- \[ \]'`＝3 未勾（rc=0）、`- \[x\]`＝0 已勾（rc=1）——与派单现量一致，本腿未动。
`probes/242/v1/`（10-03 死腿遗产）**未引为凭据、未删**。

突变名册（七形，其中 M-D3/M-E/M-F/M-G/M-B 五形＝派单自报之外的本腿自造）：

| 形 | 种下的错法 | `-run TestTicket242` 红的枚 | 整包 |
|---|---|---|---|
| M-A | 绑定比对掏空（`bindingEnforced := false` ⇒ misbind 被接受） | `SpendRejectsForgedBinding`、`SpendRequiresTheExactBinding` | — |
| M-B | 铸侧写外来摘要（`issue(nonce, it.bind+"m")`） | **无（六枚全绿）** | 5 枚别票用例红（224/87/97） |
| M-C | 烧牌顺序（`delete` 挪到绑定比对之后＝拒绝不再消费 nonce） | `SpendRejectsForgedBinding`、`SpendRequiresTheExactBinding` | — |
| M-D3 | seq 彻底折出摘要（`strconv.FormatUint(seq,10)`→`(0,10)`） | **无（六枚全绿）** | **ok，0 红（`0.431s`）** |
| M-E | store 成员扫描掏空（未知 nonce 返 `denialNone`） | `SpendRejectsForgedBinding`、**`ForgedBindingCannotSpendAnotherItemsGrant`** | — |
| M-F | `PanelAPI` 长出一枚真花令牌的 `Allow` 方法 | — | 红＝`259R1PanelAPIMethodSetIsClosed`＋`259R1PanelCarrierCarriesNoAnswerVerb`；**242 两枚全绿** |
| M-G | `PanelItem` 长出 `Permitted bool`（名与类型都不含禁词）＋补进读面名单 | — | **红＝只有 `259R1PanelItemFieldCapabilityFence`**；**242 两枚全绿** |

---

## Q1 落地性 —— 成立一半；派单第二问的后半句不成立

**结论一句**：`4bf7e683` 真在 HEAD 祖先里、真只动那一枚文件、真 113 增 0 删零产码；**但那枚文件被后续提交改过**，本腿审的是 HEAD 版不是 `4bf7e683` 版。

凭哪把尺：
- `git merge-base --is-ancestor 4bf7e683 HEAD` ⇒ rc=0。
- `git show --name-status --format= 4bf7e683` ⇒ 单行 `A internal/agent/approval/ticket242_binding_test.go`，rc=0；`--numstat` ⇒ `113 0 <同一枚>`，rc=0。
- `git hash-object <file>`＝`1eaba0972097d714882fd81f9fd7a33984b4fcfc`＝`git rev-parse HEAD:<file>` ⇒ SAME，rc=0（工作树无未提交改动）。
- ⚠ `git rev-parse 4bf7e683:<file>`＝`a4932b4466d15a59dd70478177b110741a715992` **≠** HEAD blob；`git log --oneline -- <file>` 多出一枚 `b6b1d6a4`（票 259 腿 259-r1）。差异 12 增 5 删（`logs/q1-delta-4bf7-to-HEAD.txt`，rc=0），内容逐行＝`spend` 由 `bool` 迁到 `grantDenial` 的机械适配（`d == denialNone` ≡ 旧 `true`），**判定条件等价、断言文案零移动**。
- ⚠ 派单写的 HEAD `e05b8a2b` 起手实测已不是 HEAD：起手 `git log --oneline -1`＝`0c9726f9`（111-ciif1），收工＝`a0455bb1`（A708）。`e05b8a2b` 是 HEAD 祖先（`git merge-base --is-ancestor e05b8a2b HEAD` rc=0），本身不是假提交，只是派单转述的是十几分钟前的状态级读数。

**读错了最坏会放行什么形状的真 bug**：把「HEAD 版＝`4bf7e683` 交付版」当成既定，则 259-r1 那次改写里的任何语义漂移（例如把 `d != denialNone` 写成 `d == denialNotPending`，把三种拒因混判）都会被这枚裁决当作「与自报一致」放过；而票面上写的尺（113/0）也核不出第 5 行以外的改动。

---

## Q2 这枚测试有没有牙 —— **不成立**：文件对"绑错卡"零敏感，对"烧牌顺序"有牙；且它命名的跨卡枚根本没造出"活令牌"

**结论一句**：`TestTicket242ForgedBindingCannotSpendAnotherItemsGrant`（`ticket242_binding_test.go:94-114`）递的是**从未 `issue` 过的字面量 nonce**，A 的 grant store 在那个时刻是**空的** ⇒ 它测的是"不存在的令牌也不能开卡"，不是票面那句「A 的**活**令牌花到 B 上」；把绑定层整个掏空（M-A）它**仍绿**。

凭哪把尺（本腿自跑，非复述派单）：
- **①活令牌这一维没被造出来**：`q.push` 只 `grants: newGrantStore()`（`internal/agent/approval/queue.go:173`），`issue` 只出现在 `grantNonce`（`queue.go:358`），而跨卡枚全程不调 `grantNonce`／`allow`／`allowScoped` ⇒ `itA.grants.values` 长度 0。`spend` 对空 map 走完循环落到 `return denialSpentNonce`（`approval.go:572`），**绑定比对那一行（`approval.go:567`）一次都没执行**。尺：M-A（`logs/mut-MA-test.txt`）红枚只有两枚 spend 用例，跨卡枚在 pass 名单里（rc=0）。
- **它唯一的牙是成员扫描**：M-E（`return denialNone` 顶掉 `denialSpentNonce`，＝"这条路由没在跑"那一形被反过来）⇒ 跨卡枚**红**（`logs/mut-ME-test.txt`）。⇒ 该枚对"未知 nonce 被放行"有牙，对"绑错卡"无牙。
- **②摘掉绑定校验那道条件**：M-A ⇒ `SpendRejectsForgedBinding`、`SpendRequiresTheExactBinding` 红，跨卡枚绿。这印证产码自己的话：`approval.go:545`「The binding line has ZERO independent power over a grant aimed at the wrong card; the guard inside this function is that scan」，`queue.go:343-346`「Which card a nonce is live on is decided by WHICH store the row went into, not by the digest matching」。**派单问我"防的是绑错卡还是烧牌顺序"：答＝烧牌顺序（＋成员扫描），不是绑错卡。** 票面说的形状与本腿钉住的形状不是同一件事 ⇒ 按派单令，**直接判不成立**。
- **③时序那一形**：M-C（把 `delete` 移到绑定比对之后＝被拒的 caller 可以拿同一枚 nonce 重试）⇒ 用例 1 与用例 2 双双红（`logs/mut-MC-test.txt`）。派单转述的 `precheck.md` 那句"真阻断其实是时序烧牌"本腿**独立复认**，且与 `probes/242/v1/` 零关系（没引它）。
- **④本腿新撞出来的两处零仪器（比"绑错卡"更难看，因为派单自报里把它算成有牙）**：
  M-D3＝`bindDigest` 里把 `strconv.FormatUint(seq, 10)` 换成常量 `0`（`approval.go:418`）⇒ **`-run TestTicket242` 六枚全绿**，且**整包 `-count=1` 也全绿**（`ok 0.431s`，`logs/mut-MD3-fullpkg.txt`）。⇒ 派单自报的「重放同身份因 seq 折入必得不同 digest」这一维**今天全包零尺**。
  根因具名：`TestTicket242QueuedItemsBindGrantsToTheirOwnDigest` 的作重放那半（`:85-91`）两次 `q.push(d)` 用的是同一枚 `d`，而 `d.CorrelationID` 为空 ⇒ `queue.go:163-165` 自动生 `approval-1`／`approval-2` ⇒ **两枚 corr 本就不同** ⇒ 即便 seq 完全折出摘要，`itB.bind != itA.bind` 仍成立 ⇒ 那句 `the sequence number stopped separating items` 指不到它自己命名的那件事＝**文案与判据不同物**。（本腿第一发 M-D 写成 `strconv` 未使用＝build failed，已按"编译失败不许冒充红"作废，`logs/mut-MD-test.txt`；M-D3 才是有效读数。）
  M-B＝铸侧往 store 里写外来摘要（`issue(nonce, it.bind+"m")`）⇒ **242 六枚全绿**；整包由 `ticket224_reply_grant_test.go:278/:303`、`ticket87_veto_l2_test.go:226`、`ticket97_alias_direction_test.go:119` 红 ⇒ 性质有守手，但**不在 242 名下**。
- **⑤票面 AC#1 后半"拒因指名绑定不对"**：`grantDenial`/`denialMisbound` 已落（`approval.go:481-507`），但其注释逐字写着「No routed call can produce this value —— Queue.allowScoped presents the item's own stored digest」（`approval.go:496-499`）。⇒ 路由级 misbound 今天**不可达**，任何用例都拿不到它。要拿到就是票面"乙形"之外的另一件事（从答复请求重算摘要），本格不许由本腿代为勾。
- ⛔ 全程未放宽任何断言；七形全是把产码写坏。
- 仓未被写过（尺在 Q5）。

**读错了最坏会放行什么形状的真 bug**：
(a) 以后有人把 `grantStore` 从"每卡一枚"改成"队列共用一枚"（＝真·跨卡可达那一形），或把 `allowScoped` 的 `statePending` 前置拆掉——**两枚都会让"A 的活令牌花到 B 上"变成放行**，而 242 名下唯一声称盯这件事的枚**全绿不动**，因为它手里从来没有活令牌；
(b) `bindDigest` 的 seq 被日后某次"摘要只用身份字段"的重构拿掉 ⇒ 同一 corr 重放到同一枚摘要上，配合 `replay`（`queue.go:607-618`）就能让旧卡答复落在新卡上 ⇒ 今天的包对此**完全无感**；
(c) 铸/花两端同源于 `it.bind`（`queue.go:358` 与 `allowScoped`），若有人把花侧改成传请求带来的字段而铸侧仍是 `it.bind`，`denialMisbound` 那一支才是唯一拦手——它今天只在两枚手搓摘要的用例里跑过，一次路由都没跑过。

---

## Q3 AC#2 那一格今天到底有没有对象 —— **有对象，但 242 名下那两枚不敷用；票面后半的尺今天已存在、只是记在票 259 名下**

**结论一句**：AC#2 的反射判据今天真实存在（`internal/agent/approval/ticket242_panelface_test.go`，落地提交是 `2048d6b6`＝242-r2，**不是** `4bf7e683`），反射对象逐字钉 `approval.PanelItem`、**没碰 `panel.ApprovalCardView`**＝禁区合规；但本腿种了 `Permitted bool` 之后，242 那两枚**双双全绿**，真正咬住的是票 259 的能力尺。

凭哪把尺：
- 存在性尺：`grep -rn "reflect\." --include=*.go internal/agent/approval/ internal/panel/ cmd/wisp/ | grep -i "PanelItem|StructOf|NumField|FieldName|TypeOf"` ⇒ rc=0，命中 `ticket242_panelface_test.go:31` 与 `:54` 逐字 `reflect.TypeOf(PanelItem{})`；`ticket259_panel_capability_rulers_test.go:67/:170/:189/:334` 另有六枚。`git log --oneline -- ticket242_panelface_test.go` ⇒ 唯一提交 `2048d6b6`，rc=0。
- 禁区尺：该文件 `:16-20` 逐字写「⛔ Deliberately NOT panel.ApprovalCardView」，且本腿在 `internal/agent/approval/` 范围内未见任何对 `panel.ApprovalCardView` 的反射替代扫描（`internal/panel/approval_test.go:176` 那枚反射的对象是它自己的 view 型，不是 242 的尺）。
- 牙齿尺 M-G（`logs/mut-MG-test.txt`，rc=1）：给 `PanelItem` 加 `Permitted bool`（名字与类型都不含 `grant/allow/approve/nonce/token` 五枚禁词）＋把 `"Permitted"` 补进 `panelItemReadFace` ⇒ `TestTicket242PanelItemReadFaceIsFullyDeclared` **PASS**、`TestTicket242PanelItemStaysGrantFree` **PASS**（＝"名单与结构互为镜像"那句被同时改名单绕过），只有 `TestTicket259R1PanelItemFieldCapabilityFence` **FAIL**。
- 牙齿尺 M-F（`logs/mut-MF-fullpkg.txt`，rc=1，68 PASS／2 FAIL）：照票面那一形把 `PanelAPI` 扩到"能改变卡的状态"（加 `Allow(correlationID, grant string) error` 并实现成真花令牌 `p.q.allow(...)`）⇒ 红的是 `TestTicket259R1PanelAPIMethodSetIsClosed` 与 `TestTicket259R1PanelCarrierCarriesNoAnswerVerb`；**242 名下两枚全绿**（它们只看 `PanelItem`，不看 `PanelAPI`）。
- ⚠ **顶回台账一处过期**：票面 `:56` 那句「M-E ⇒ 本包仍 66 PASS／0 FAIL＝票面那句能力侧判据零仪器」是**在 `b6b1d6a4`（259-r1 落三枚能力尺）之前**取的数，今天不成立：本腿同形（M-F）在 HEAD 上**两枚红**。按"状态级读数不许转述只能各腿现量"那条定式，这一句要就地更新。

**读错了最坏会放行什么形状的真 bug**：若把这格读成"AC#2 已由 242 交付"，则任何人往出向读面加一枚**换了名字的**可答复字段（`Permitted`／`Approved`／`Settled`／`Yes`）都能同时过 242 的两枚尺；一旦面板侧出现 type-assert 到更宽接面或日后把 `panelAPI` 具体型直接交给宿主，"面板来源的 allow 在结构上不可能"这条 SPEC-06 §9 layer 3 的守卫就只剩 259 那三枚尺挡着——而它们记的是**另一张票**的名字，按票翻勾时会以为这格没人守。

---

## Q4 AC#1 与 AC#2 是不是同一块石头 —— **不是；该两格，但不该合成一枚**

**结论一句**：两枚判据互相零敏感，值两格；AC#3 那句「甲形够不够」＝**甲形够（不必解冻乙形）**，缺的三枚尺里有两枚已由票 259 落在同一个包里。

凭哪把尺（具名到文件行）：
- M-A（把 `approval.go:567` 的绑定比对掏空）＝AC#1 那一侧坏 ⇒ `ticket242_panelface_test.go:30` 与 `:53` 两枚**全绿**（`logs/mut-MA-test.txt` 的 PASS 名单里逐名在列）。
- M-G／M-F（把 `ui.go:49-56` 的 `PanelItem` 或 `ui.go:167-171` 的 `PanelAPI` 弄坏）＝AC#2 那一侧坏 ⇒ `ticket242_binding_test.go` 四枚**全绿**（M-G/M-F 的读数里 `TestTicket242Spend*/Queued*/Forged*` 无一红）。
- 两两互不敏感＝两枚缺陷，不是一枚。判据＝本仓铁律"判据换成反形还全绿＝它对这件事不敏感，不许当凭据"。
- 另有一处**必须点破的不同物**：AC#1 拦的是"令牌花到错误的卡上"（写侧/花费侧），AC#2 拦的是"能力出现在面板可见的数据上"（出向读侧）。今天真正的守手也不同：前者＝`grantStore` 的每卡一枚（`approval.go:436-450`）＋队列状态拒（`allowScoped` 的 `statePending`），后者＝`PanelAPI`/`PanelItem` 的封闭集（259 的尺）。合成一枚会把"两枚不同守手"压成一句，日后任一守手被拆都不知该找哪格。

**读错了最坏会放行什么形状的真 bug**：若判"同一块石头"并只留一格，则要么读面尺（242-panelface）随绑定尺一起被撤，面板拿到 `Permitted`；要么绑定尺随读面尺一起被撤，`it.bind` 那行比对变成装饰——两形今天都**只由剩下那一格守着**，合格即同时失明。

---

## Q5 不新增红＋门禁 —— 全绿，零新增红，仓未被本腿写过

**结论一句**：`internal/agent/approval/` 前后同数同册，vet／d22scan／gofmt 三把 rc=0，`pending_read.go` 那枚是 autocrlf 幻影非红；`cmd/wisp` 整包按约束不跑。

凭哪把尺（每条自己落 rc，测 rc 那句前面无管道）：
- `go test -count=1 ./internal/agent/approval/`：前 rc=0、后 rc=0（末行 `ok … 0.396s`）；`go test -count=1 -v`：前 **70 PASS／0 FAIL**，后同数；逐名作差 `diff <前 PASS 册> <后 PASS 册>` ⇒ **0 行，rc=0**＝新增红 0、名册无漂移（`logs/q5-name-diff.txt`）。
  ⚠ 与派单自报的 45/47 PASS 不同＝其间别票（259-r1 等）往同包加了用例，不是本腿读数。
- `go vet ./internal/agent/approval/` ⇒ **rc=0，零输出**（`logs/q5-vet.txt`）。
- `sh scripts/d22scan.sh` ⇒ **rc=0**，末行 `d22scan: clean - no D22 ban violations`（`logs/q5-d22scan.txt`；按约束未在仓根跑 `go run ./tools/d22scan`）。
- `gofmt -l internal/agent/approval/` ⇒ rc=0，名册一枚 `internal\agent\approval\pending_read.go`。行尾符尺 `tr -cd '\r' | wc -c`：`pending_read.go` **worktreeCR=131／HEADblobCR=0** ⇒ 与派单现量一致，autocrlf 幻影非红；`ticket242_binding_test.go`、`ticket242_panelface_test.go` **CR 双 0**＝本腿审的两枚件不带幻影。（未用 `grep -c $'\r'`。）
- 仓未被写过（overlay 证明）：`git hash-object` 对 `git rev-parse HEAD:<f>`，五枚逐名 SAME，rc=0：
  `approval.go 67fb1468…`／`queue.go 66fec7ae…`／`ui.go`／`ticket242_binding_test.go 1eaba097…`／`ticket242_panelface_test.go 819002515f…`。
- 起手 `git status --porcelain`＝778 行（`logs/git-status-start.txt`）；收尾＝`logs/git-status-end-beforecommit.txt`；作差 55 行（`logs/git-status-delta.txt`）**全部属于别腿**（111-ciif1 的 `probes/111/**`、`.github/workflows/ci.yml` 等），`grep -c "ticket242|approval/approval.go|approval/queue.go"` 对差集 ⇒ **0 命中，rc=1**。只登记，未提交、未还原。
- ⛔ `cmd/wisp` 整包**按约束未跑**（255-v2／35-v7 正在同面包上取数）。本腿 Q2⑤ 关于 `allowScoped` 路由级的断言取自**产码注释＋`denialMisbound` 的可达性**，不是取自 cmd/wisp 现量。〔归编排者补〕

**读错了最坏会放行什么形状的真 bug**：把幻影当红去"修"⇒ 顺手给 `pending_read.go` 改行尾符，动到不属于本票的产码哈希、搅乱同时在飞的 blob 对拉；把别腿新增的用例数当成本票读数⇒ 45/47/70 三档会被写成同一条"包一直是这个数"，掩盖 242-r2 之后又落了两批尺。

---

## 我没量到的／派单前提不成立处

**顶回派单（按"原文为准"这一条，逐具名）**：
1. **HEAD 号**：派单写 `e05b8a2b`，起手实测 `0c9726f9`，收工 `a0455bb1`。`e05b8a2b` 真在祖先里，不是假号，但派单转述的是过期状态级读数。
2. **"那枚文件没被后续改动过"**：不成立。`b6b1d6a4`（票 259-r1）改过 `ticket242_binding_test.go`（12 增 5 删）。判定条件等价，但本腿审的**不是** `4bf7e683` 那一版。
3. **"AC#2 内容锚的那把尺"不在 `4bf7e683` 里**：派单第 1 节把三格都算在被审提交名下，实际 AC#2 由 `2048d6b6` 交付。⇒ 若按派单去 `4bf7e683` 找 AC#2 的尺，会得出"AC#2 零判据"这一**错**结论（派单 Q3 正是在问会不会这样，答案是：有对象，但要认对提交）。
4. **派单自报的两枚突变（M1/M2）之外，本腿测出两处"自报里被当成有牙、实测零仪器"**：seq 折出摘要（M-D3）整包全绿；铸侧写外来摘要（M-B）242 名下全绿。第 3 条尤其要紧——`TestTicket242QueuedItemsBindGrantsToTheirOwnDigest` 的**断言文案（"the sequence number stopped separating items"）与它实际能检出的错法不同物**，属本项目反复出现的"恒真句"那一形。
5. **派单第①形预警成立**：名叫跨卡那枚只造了"两枚活卡"，没造"活令牌"，`spend` 全程不执行绑定比对。⇒ 按派单令，AC#1 的甲形**判不成立**，不许翻勾。
6. **票面 `:56` 的"M-E⇒能力侧零仪器"已过期**（HEAD 上 M-F 两枚红），归 259-r1 的 `b6b1d6a4`。这条是账目更新，不是新缺陷。

**没取到数的格**：
- `cmd/wisp` 整包读数（约束不跑，归编排者补）：Q2⑤「路由级拿不到 `denialMisbound`」目前是**产码注释级**凭据（`approval.go:496-499`、`queue.go:343-346`），未由整包现量佐证。
- M-B 只跑了整包一次（`-run` 全册），未逐枚作差到测试名（红枚已具名 4 条，够用但非全册差集）。
- 票面 AC#2 前半那句「出现 `Grant`／**可答复类**字段就红」里的"可答复类"外延，本腿只测了 `Permitted bool` 一种拼法（M-G）；其余拼法（`Yes`/`Settled`/`Decision`）未各测一发，但同一条禁词表对它们同样无感，属同一枚结论。
- `probes/242/v1/` 里那批死腿 logs（`base-v.txt`、`bindsites.txt`、`commit-*.txt`）本腿一律未读数、未引用、未删；派单要求的"自己复跑那一把"针对的是 precheck 那条时序结论，已由 M-C 独立复跑，`precheck.md` 的原始读数本身未被逐字复认。
