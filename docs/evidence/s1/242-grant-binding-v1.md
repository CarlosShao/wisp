# 242-grant-binding — 非实现者对抗验收表（腿 `242-v1`）

被验收件：`internal/agent/approval/ticket242_binding_test.go`（`4bf7e683`，四枚用例）＋
`internal/agent/approval/ticket242_panelface_test.go`（`2048d6b6`，两枚尺）。
票面：`.scratch/wisp/issues/242-the-grant-binding-layer-and-the-panel-facing-read-surface-both-have-zero-rulers.md`（AC#1／AC#2／AC#3 三格）。
身份：被验收代码＝编排者本人代笔（他同时是被验收方），本表是**攻击腿**。
判据口径：**"判据换成反形仍然全绿＝它对这件事不敏感，不许当凭据"**。
⛔ 本表不翻任何 AC 框；AC 框一枚未碰。全部读数只来自 `go test ./internal/agent/approval/ -count=1`（无跨包测量）。

---

## §0 起手锚（同发取）

| 项 | 读数 |
|---|---|
| `date` | `Sat Oct 3 09:09:49 CST 2026` |
| `git rev-parse HEAD` | `8a3790f0dd057a525f31cc5a9606193371b1eef7` |
| `git status --porcelain -- internal cmd` | **空**＝起手态（此后两发突变之间的复位基准） |
| 全仓 porcelain | 非空：`.gitignore`／`.scratch/**`／`design/assets/**` 被他人删除（在册红因族，非本地界） |
| 起手包读数 | `go test ./internal/agent/approval/ -count=1 -v` → **66 `--- PASS`／0 `--- FAIL`**，`ok 0.345s`，六枚 `TestTicket242*` 全绿（`.scratch/wisp/probes/242/v1/base-v.txt:18-29`） |

备份锚（`git cat-file blob HEAD:<path>`，md5 与编排者 commit 里自报的三枚逐枚吻合 ⇒ 他量的是同一份盘）：

```
d97267a803fcba28ce1d7913bfbc9b58  approval.go.orig
a170f3f2899251f0b826d464e3973d8f  queue.go.orig
a4cee69ef76adc54ea0154527f283604  ui.go.orig
dc99738456f72cdb0fabad1cb81e99bf  ticket242_panelface_test.go.orig
2e2d9bca0c38d23f2dafbfdf9daf5f64  ticket242_binding_test.go.orig
```

---

## §1 逐格判语

### AC#1（真·跨卡：A 的活令牌花到 B 上⇒拒且拒因指名"绑定不对"）＝**附条件成立**

**成立的那半**（有牙，实测）：
- 票面给 AC#1 指定的正控就是 **D2**（"`grantStore.spend` 不再比对 binding 摘要"）。我把这一种落实了＝M-A2（`approval.go:303` 改成
  `equalSecret(stored, bind) || len(bind) >= 0`，用法保留、编译通过、锚点在盘上量到）：**两枚新用例红**，红句见 §2。
  票面"那一层的**内容**"那行记的旧读数（"删掉比对⇒定向尺绿＋包也绿"）**已被这两枚尺改掉** ⇒ 本格欠账的 store 那一半确实还上了。
- 铸侧也有牙：M-C1（`queue.go:162` 掉 `d.TaskID`）⇒ 唯一红＝`TestTicket242QueuedItemsBindGrantsToTheirOwnDigest`。

**不成立的那半**（票面字面要求的形状，实测**零尺**）：
1. **票面要求的载具根本不在测试里**：AC#1 要"造两枚同时活着、correlation id 不同的卡，把 A 的未花令牌递给
   `Native().Allow(B的corr, A的grant)`"。六枚新用例**一枚都没调用 `Allow`／`DecideFromNative`／`allowScoped`**——
   `ticket242_binding_test.go` 全文只有 `newGrantStore()`／`q.push()`／`it.grants.spend()`（`:20-23`、`:59-64`、`:104`）。
2. **恒真／路不跑的对照（§2 M-B）＝判据对此完全无感**：我把允许侧的路由整个弄哑
   （`queue.go:235` `lookupForAllowLocked` 恒返 `nil`，即"这条路由根本没在跑"），**六枚 `TestTicket242*` 全绿**，
   同时 12 枚在册路由级用例红（含 `TestGrantIsSingleUseAndBoundToItsItem`、`TestL2NativeAllowExecutesThroughTheBridge`）。
   ⇒ 票面那句"**不是这条路由没在跑**"**没有仪器**：这两件事在新尺下读数是同一个颜色。
3. **"拒因指名绑定不对"这一枚不存在**：`spend` 返回 `bool`，没有拒因可指名（`approval.go:292-306`）；
   而 API 层唯一的授权错误是 `ErrBadGrant`（`ui.go:127-130`，注释明写"the message never says which"，防探针）。
   ⇒ 票面这一句与产码的防探针设计**正面冲突**，不是测试能单方面满足的（这一刀归人，见 §4）。
4. **名字叫"跨卡"的那枚是恒真**：`TestTicket242ForgedBindingCannotSpendAnotherItemsGrant`（`:87-107`）递给 store 的
   nonce 是字面量 `"leaked-or-guessed-nonce"`，**它从没被 `issue` 进这个 store** ⇒
   `spend` 在 `:299-301` 的 nonce 查找处就 `continue` 掉、根本走不到 `:303` 的绑定比对。
   实测：**M-A2 把绑定层掏空之后这一枚仍然绿**（`--- PASS`，见 §2）。
   ⇒ 按本仓铁律，它对"跨卡绑定"这件事**不敏感，不许当凭据**；它实际测的是"没发过的令牌花不掉"（新鲜度／存在性）。
5. **票面写着"前置是本票自己的活"的那枚载具没改**：`cmd/wisp/subagent_selfapproval_197_test.go:109`
   现量仍是 `TaskID: taskID, CorrelationID: taskID,`（两卡同名那枚缺陷原样在）。
   编排者在 `4bf7e683` 的 commit 正文里单方面"撤销":109"前置"——**票面文字没有被改**（我没动票面），
   所以按票面读，AC#1 的前置未交付。

**条件（翻勾只许按这个口径）**：本格若翻，只能翻成"**铸／花两端的 store 层绑定尺已立（D2 已闭合）**"，
**不许**读成"真·跨卡路由有尺"或"拒因指名绑定"。后两半该退回（见 §4）。

### AC#2（出向读面不得带 grant；名字尺＋能力尺）＝**不成立**

- **形状对那半成立**：反射钉的对象确实是 `approval.PanelItem`，不是 `panel.ApprovalCardView`
  （`ticket242_panelface_test.go:31` 与 `:54` 两处 `reflect.TypeOf(PanelItem{})`；`ui.go:50` 是定义点）——票面"落点两候选／推荐甲"的甲形要求已执行。
  名字尺也有牙：M-D1 种下 `PanelItem.GrantToken string` ⇒ **两枚尺齐红**（未声明红＋词面两红，逐字见 §2），
  编排者自报的 M3 读数我复现成功。
- **票面明写的第二半没做**（判"不成立"的直接理由）：AC#2 原文要求"再加一发**能力侧**判据（不是词面）：
  把 `PanelAPI` 的出向读面扩到'能改变卡的状态'那一形也要红"。落地件里**零枚仪器碰 `PanelAPI`**：
  全仓 `_test.go` 里 `PanelAPI` 只出现在 `queue_test.go:136` 的一句**注释**（不是断言），
  `approval` 包内没有任何 `NumMethod`／方法名集尺（`NumMethod` 的在册用法只在 `internal/tools` 三包测试里，射程不含 `PanelAPI`）。
- **词面名单可被改名绕过（实测 M-D2）**：`PanelItem.Permitted bool` ＋ 把 `"Permitted"` 补进 `panelItemReadFace`
  ⇒ **包内 66 枚全绿、0 红**（`.scratch/wisp/probes/242/v1/m-d2.txt`）。
  字段语义就是"这卡被允许了"，两枚尺都看不见它。⇒ 名单尺是**词面**的，且名单与被它管的字段**由同一只手维护**（镜像检查自声明）。
- 按派单口径："只给名字不给能力的形要判不成立" ⇒ **AC#2＝不成立**。

### AC#3（本格不设产码要求：两枚判据是不是同一块石头；甲形够不够）＝**不成立（＝不许勾）**

- **是不是同一块石头：不是**，值两格。AC#1 守的是"令牌与哪件事的绑定在花的那一刻成不成立"（`queue.go:376` 那一跳），
  AC#2 守的是"面板能看见什么形状"（`ui.go:50` 的字段集与 `ui.go:167` 的方法集）。
  两者共用的心虚处是"grant 别漏到不受信的面"，但机制不同：M-A2 只红 AC#1 那两枚、M-D1／M-D2 只动 AC#2 那两枚，
  **互相零敏感**（实测：M-D1 时四枚 binding 用例全 `--- PASS`；M-A2 时两枚 panelface 用例全 `--- PASS`）⇒ 两格各自独立成立，没有重复计。
- **甲形够不够：不够**。M-D2 就是"甲形够不够"的反证。
  但**不够不必靠乙**：能力尺可以在甲内补（同包，`PanelAPI` 是本包导出接口，`ui.go:167-171`），
  例如"方法名集恰好＝`{Reject, Head, View}`＋任何别名／新方法出现即红"，再加一枚"没有任何产码函数把 `PanelItem` 的字段读回队列状态"的尺。
- 票面对 AC#3 的处置文字："若判'要乙'，本格不许勾，写清'等 owner 那句解冻'"。
  我的判断更靠近中间态并**取严**：**不许勾，但理由不是"等解冻"，而是"甲内缺能力尺"**（解冻 `internal/panel/l2_grant_boundary_test.go` 对本格并非必需）。

### 生产接线（第 3 问的单列）＝**已接，但绑定比对在花侧是恒等式**

见 §3：`go build`/`go vet` 的颜色我一律不作证据（本包有一枚零调用者的函数正是这类盲区，见 §3 末）。

---

## §2 突变自证表

纪律：每发跑前 `grep -n` 锚点证突变真落盘；跑完立刻 `git cat-file blob HEAD:<path> > <path>` 还原并 `md5sum` 回锚；
**两发之间**量 `git status --porcelain -- internal cmd` 回到起手态（＝空）。全部读数取自
`go test ./internal/agent/approval/ -count=1 -v`，判绿只认 `--- PASS`／`--- FAIL` 行。
输出文件在 `.scratch/wisp/probes/242/v1/`。

| 编号 | 种什么（file:line ＋逐字突变形） | 哪枚必须红 | 实测红句**逐字**（含 file:line） | 其余读数 | 还原复跑终值 |
|---|---|---|---|---|---|
| **M-A2** 花侧绑定层掏空（＝票面正控 D2） | `approval.go:303` `return equalSecret(stored, bind)` → `return equalSecret(stored, bind) \|\| len(bind) >= 0`（用法保留，非 unused-var） | `TestTicket242SpendRejectsForgedBindingAndConsumesTheNonce`、`TestTicket242SpendRequiresTheExactBinding` | `ticket242_binding_test.go:24: AC#1 RED: a grant bound to item one was spent with item two's digest - the binding layer did not reject the forged binding` ／ `ticket242_binding_test.go:41: AC#1 RED: empty binding spent a live grant` | 全包红数＝**2**；`TestTicket242ForgedBindingCannotSpendAnotherItemsGrant`＝`--- PASS`；**12 枚路由级用例全绿**⇒ 绑定层从路由不可见 | md5 回 `d97267a8…`，porcelain 空 |
| **M-B** 恒真对照：让路都不跑 | `queue.go:231-236` `return q.byID[corr]` → `return nil // 242-v1 M-B: the allow route never finds any card` | 票面 AC#1 要求这里**必须红**（"不是这条路由没在跑"） | **无红句可引——六枚 `TestTicket242*` 逐枚 `--- PASS`**（`m-b.txt`） | 同发 12 枚在册用例红：`TestGrantIsSingleUseAndBoundToItsItem`／`TestL2NativeAllowExecutesThroughTheBridge`／`TestPanelSourcedAllowIsRejectedOnEveryForgeableAxis`／`TestReplayRedisplaysUnderAFreshGrant`／`TestAnAliasCanNeverBuyAnAllow`／`TestTicket224*`×4／`TestR7SizedBatchEscalatesToL2Unaggregated`／`TestL2QueueAutoRejectsAt300s…`／`TestAVetoOnAnUnknownNameLeavesTheCardAllowableByItsOwnGrant` ⇒ 路确实被弄哑了，新尺确实看不见 | md5 回 `a170f3f2…`，porcelain 空 |
| **M-C1** 铸侧掉一个身份字段 | `queue.go:162` `bindDigest(corr, d.TaskID, …)` → `bindDigest(corr, "", d.Tool, d.LevelString(), q.seq, d.Args)` | `TestTicket242QueuedItemsBindGrantsToTheirOwnDigest` | `ticket242_binding_test.go:73: AC#1 RED: the queued item's binding does not cover its own identity fields (seq/task/tool/level/args)` | 全包红数＝**1**（`m-c1.txt`）⇒ 铸侧"逐字段覆盖"有牙 | md5 回 `a170f3f2…`，porcelain 空 |
| **M-C2** 铸侧 seq 折入变常量 | `queue.go:162` `…, q.seq, d.Args)` → `…, 0, d.Args)` | 票面口径下该由"重放"那一枚红 | 实际红在**同一行**：`ticket242_binding_test.go:73: AC#1 RED: the queued item's binding does not cover its own identity fields (seq/task/tool/level/args)` | 红数＝**1**，红因是 `:73` 的覆盖断言（测试用 `itA.Seq` 复算），**不是** `:82-84` 的重放断言 ⇒ 见 M-C3 | md5 回 `a170f3f2…`，porcelain 空 |
| **M-C3** **seq 从摘要里彻底消失**（digest 内掏空） | `approval.go:255` `strconv.FormatUint(seq, 10)` → `strconv.FormatUint(seq-seq, 10)`（值恒 `"0"`，用法保留、编译通过） | 票面 AC#1 主张"重放同一身份因 seq 折入必不同 digest" ⇒ 该红 | **无红句**：六枚 `TestTicket242*` 逐枚 `--- PASS`，**全包 66 PASS／0 FAIL，`ok 0.345s`**（`m-c3.txt`） | ⇒ "seq 折入"这半句**反形全绿＝不敏感**，不许当凭据。根因两处：① 测试的 `want` 用**同一个** `bindDigest` 复算（`ticket242_binding_test.go:71` 自指），函数内部怎么改都追不上；② 那两枚"重放"项 `corr` 本来就不同（`queue.go:150` 由 `seq` 生成 `approval-<seq>`），`corr` 一项就把它们分开了，**同身份重放在这个载具里构造不出来** | md5 回 `d97267a8…`，porcelain 空 |
| **M-D1** 种下带味道的字段（＝编排者自报的 M3） | `ui.go:57` `PanelItem` 末尾加 `GrantToken    string` | 两枚 panelface 尺 | `ticket242_panelface_test.go:39: AC#2 RED: PanelItem.GrantToken (string) is visible to the panel but not declared in panelItemReadFace (ticket 242: a new field on the panel read face must be declared, not slipped in)` ／ `ticket242_panelface_test.go:61: AC#2 RED: PanelItem.GrantToken (string) smells like an allow channel ("grant"); …` ／ 同 `:61` 第二行 `… ("token"); …` | 红数＝**2**；四枚 binding 用例全 `--- PASS`（两族互不敏感） | md5 回 `a4cee69e…`，porcelain 空 |
| **M-D2** **改个字段名绕过词面名单**（能力侧反证） | `ui.go:57` 加 `Permitted     bool` **且** `ticket242_panelface_test.go:28` 把 `"Permitted"` 补进 `panelItemReadFace`（＝同一只手补名单的真实场景） | 票面 AC#2 要求的能力侧那一形 ⇒ 该红 | **无红句**：两枚尺逐枚 `--- PASS`，**全包 66 PASS／0 FAIL，`ok 0.347s`**（`m-d2.txt`） | ⇒ 一个语义即"这张卡被允许了"的字段能带着名单一起静默过审；名单尺只认词，不认能力 | md5 回 `a4cee69e…`＋`dc997384…`，porcelain 空 |

**终态复跑（全部还原之后，同一次 `-count=1 -v`）**：
`PASS=66 / FAIL=0`，`ok github.com/CarlosShao/wisp/internal/agent/approval 0.338s`，
五枚被改文件 md5 逐枚回到 §0 锚（`diff` 备份表＝`backups_match_HEAD=YES`），`git status --porcelain -- internal cmd`＝**空**。
证据件：`m-a2.txt`／`m-b.txt`／`m-c1.txt`／`m-c2.txt`／`m-c3.txt`／`m-d1.txt`／`m-d2.txt`／`base-v.txt`／`final-v.txt`／`backup/`。

---

## §3 生产调用者现量（bindDigest 铸／花）

尺＝`grep -rn` 全仓 `internal cmd --include=*.go`（`bindDigest`／`.issue(`／`.spend(`／`it.bind` 逐符号），**不是** `go build` 的颜色。

| 端 | 调用点枚数 | file:line（逐枚） |
|---|---|---|
| `bindDigest` 定义 | 1 | `internal/agent/approval/approval.go:253` |
| `bindDigest` **产码铸点** | **1** | `internal/agent/approval/queue.go:162`（`push` 里 `it.bind = bindDigest(corr, d.TaskID, d.Tool, d.LevelString(), q.seq, d.Args)`） |
| `it.bind` 写入点 | **1** | 同上 `queue.go:162`（此后再无第二次赋值 ⇒ 生命周期内常量） |
| `grantStore.issue` | 1 | `internal/agent/approval/queue.go:336` `it.grants.issue(nonce, it.bind)` |
| `grantStore.spend` | 1 | `internal/agent/approval/queue.go:376` `spent := it.grants.spend(nonce, it.bind)`（`allowScoped`，`:361`） |
| `grantStore.revoke` | 2 | `queue.go:288`、`queue.go:426`（`revokeGrants`） |
| 测试侧调用（对照用） | 3 | `ticket242_binding_test.go:39`、`:48`、`:71` 直接调 `bindDigest`；`:23`、`:40`、`:50`、`:104` 调 `spend` |

**接线判定**：**已接**（铸 1 枚、花 1 枚，`Native().Allow` → `allowScoped` → `spend` 这一跳链上有真调用点，不是仪器绿／产码空）。

**但花侧的绑定比对在生产里是恒等式**（三条同时成立才有这句）：
① `issue` 存的是 `it.bind`（`:336`）；② `spend` 比的是**同一个 `it` 的** `it.bind`（`:376`，`it` 由 `lookupForAllowLocked(corr)` 得到）；
③ `it.bind` 自 `:162` 之后没有任何第二次写（枚到 1）。
⇒ `approval.go:303` 的 `equalSecret(stored, bind)` 两端**永远读同一个值**，"misbound" 分支从路由不可达。
实测同条：M-A2 掏空它之后**唯一变红的只有两枚直接手搓 store 的新用例**，12 枚走 `Allow` 的在册路由级用例全绿（§2 M-A2 行）。
⇒ 今天真正挡住"拿 A 的令牌花到 B"的是**per-item store 的成员检查**（`:299-301`），不是绑定摘要。
所以 AC#1 那半句"跨卡绑定"的**安全效果**由 store 隔离承担；**摘要那一层本身在生产里没在做事**——
票面把它写成"bindDigest 也确实比对"（`:253` 铸、`:292` 花）在"有没有调用点"意义上对，在"有没有起作用"意义上需要补这句。

**`go build`／`go vet` 证不了接线的现例**：`internal/agent/approval/ticket242_binding_test.go:111` 的
`func sameDigest(a, b string) bool` **零调用者**（全仓 `sameDigest` 只命中定义与它自己 `:109` 的注释），
它的注释还写着"keeps the failure message honest"——**没有任何断言用它**。包照样 `ok`。
⇒ 编排者 commit 里那句"build rc=0／vet rc=0"对这两枚判据件的内部有效性**零证据量**。

---

## §4 我建议编排者怎么处置

1. **AC#1：可以翻，但只能按收窄口径翻**，且我建议先落一行票面附注再翻（翻勾的口径写清"铸／花两端 store 层尺＋D2 已闭合"）。
   理由：票面指定的正控 D2 实测已红（M-A2），铸侧覆盖实测已红（M-C1）——这是真的把旧读数"删掉比对⇒全绿"翻转了。
   ⛔ **不许**把这一格说成"真·跨卡路由有尺"或"拒因指名绑定不对"（M-B 六枚全绿＝这半句零尺）。
2. **AC#1 该退回重做的两半**（零件缺失，不是判语）：
   - **载具那一半**：`cmd/wisp/subagent_selfapproval_197_test.go:109` 仍 `TaskID: taskID, CorrelationID: taskID`。
     票面写着"前置是本票自己的活"，而 `4bf7e683` 的 commit 正文把它"撤销"了——**票面文字没跟着改**。
     建议：要么在台账具名记这笔撤销并改票面（**改票面＝人工批准面，我不碰**），要么按票面把载具补上再谈这一格。
   - **"拒因指名绑定不对"那一半**：与 `ui.go:127-130` 的 `ErrBadGrant` 防探针设计**正面冲突**（错误值故意不区分缺失／已用／不绑定）。
     这一刀不是测试能自证的，**归人拍板**（倾向：票面这句话按产码现状改写，或另立一票讨论"可指名拒因的安全代价"）。
   - **"重放 seq"那一半**：M-C3 全绿＝反形不敏感。要它成立必须造"同 corr 两次入队（第一次已离队）"的尺；
     现在那两枚项 `corr` 天生不同（`queue.go:150`），构造不出同一身份。⛔ 不许用"改名"糊过去——`bindDigest` 的自指复算（`:71`）会让函数内部的任何删改都追不上。
3. **AC#2：不许翻**（不成立）。缺的零件是**一枚能力尺，可在甲内补，不必解冻乙**：
   - 最小形 A：`reflect.TypeOf(PanelAPI(nil)).NumMethod()` 与方法名集**恰好** `{Reject, Head, View}`，出现任何第四个或任何 allow-ish 同义形状即红；
   - 最小形 B（打 M-D2 那一发）：`PanelItem` 的字段集必须是**封闭名单**且**只读**——
     "没有任何产码函数以 `PanelItem` 为参、把其中字段读回队列状态"（这条能咬住 `Permitted bool`，词面尺咬不住）。
   - 顺手：M-D2 这一发应作为**常驻**自证（种下→红→还原）留在票的验收脚注里，否则下一枚腿还会读成"这格有牙"。
4. **AC#3：不许翻**，理由＝甲形现量不够（§1 AC#3），且**不是**"等 owner 解冻 `internal/panel/l2_grant_boundary_test.go`"。
   `Q-74` 那一刀可以先不问：补甲内能力尺不碰任何冻结件。
5. **两枚判据不是同一块石头**＝值两格，这点可以直接采信（§1 AC#3 的互敏感实测）。
6. **别把 `sameDigest` 留在文件里当装饰**（§3 末）：要么有一枚断言用它，要么删（它已经在教下一枚腿"这里有一条没执行的诚实性检查"）。

---

## §5 我不确定的地方

1. **我只测了 seq 从摘要里消失（M-C3），没单独测 `args`／`level` 从 `bindDigest` **函数内部**消失**。
   按 M-C3 的机理（测试自指复算 `:71`）我**推断**同类全盲，但没读数；不要把这句当读数引用。
2. **"同 corr 重放"到底能不能构造，我只读到 `push` 的改名分支**（`queue.go:152-154` 冲突时 `corr#seq`）。
   若某条产码路径能在**同 corr** 上铸第二枚活令牌（`grantNonce` 对同一 `it` 再发一次，`:336`），
   那 M-A2 之后"`stored` 仍等于 `it.bind`"的恒等式结论不变，但票面 D2 的安全含义要重估。我没测这一发（会动 `cmd/wisp` 地界外的时序）。
3. **`ErrBadGrant` 的防探针取舍是不是 owner 已定案**，我没查台账（`docs/reports/pending-and-issues.md` 的 `A##` 太多，本轮没检索到对应条目）。
   若已定案，则票面 AC#1"拒因指名绑定不对"这句是**票面缺陷**而非交付缺陷——这个归因差别会影响他怎么改票，会影响我 §4 第 2 条的措辞。
4. **197-r3 那条腿是否正在改 `subagent_selfapproval_197_test.go:109`**：我只在 09:09 取了那枚行的现量（仍同名），
   `cmd/wisp` 此刻是别人的写面，我没有第二次读数，不能断言"永远没人改"。
5. **M-B 我选的是"路由整个哑"（`lookupForAllowLocked` 恒 `nil`）这一形**，
   没有试更窄的形（比如只让 `allowScoped` 在 `spend` 前 `return ErrNotPending`）。结论方向不变（新尺不碰路由），但"哪一枚在册用例红"的名册会随形变。
6. **12 枚路由级用例红**我只逐枚引了名字，没逐枚看它们**为什么**红（超时 3.00s 那批大概是等卡超时）。
   这对 M-B 的结论非必要（关键读数是"新尺全绿"），所以我没展开；若他要拿这份名册去别处引用，请自行复算红因。
7. **票面"零尺"那行旧读数我复算成"现量已非零"**：`grep -rln "bindDigest" --include=*_test.go internal cmd` 今天命中
   `cmd/wisp/subagent_selfapproval_197_test.go` 与新的 `ticket242_binding_test.go` ＝**2 枚文件**（票面写 1 枚）。
   这是票面被这次交付改写了，不是矛盾；只是别按票面的"1 枚"读现状。

---

## §6 判语

- **AC#1＝附条件成立**（铸／花 store 层有牙：M-A2 两枚红、M-C1 一枚红；真·跨卡路由与"拒因指名绑定"两半＝零尺：M-B 六枚全绿、名字叫跨卡那枚在绑定层掏空后仍绿）。
- **AC#2＝不成立**（名字尺有牙：M-D1 两枚红；能力尺＝零仪器，词面名单一发改名即全绿：M-D2 66 PASS）。
- **AC#3＝不成立／不许勾**（两枚判据不是同一块石头这点成立；"甲形够不够"答"不够"，且缺的零件不必解冻乙形）。
- **生产接线**：`bindDigest` 铸 1 点（`queue.go:162`）、花 1 点（`queue.go:376`）＝**已接**；
  但花侧比对两端读同一个 `it.bind` ⇒ **绑定摘要在生产里恒等**，实际防线是 per-item store 的成员检查。
  `go build`／`go vet` 的颜色不作证据（现例：`ticket242_binding_test.go:111` 零调用者函数，包仍 `ok`）。
- **建议**：只按 §4 第 1 条的收窄口径翻 AC#1；AC#2／AC#3 不翻；§4 第 2 条两半退回，其中"拒因指名"那一刀归人拍板。
- **纪律自证**：AC 框一枚未碰（票面三格仍 `[ ]`，本表没有写进票面一个字）；⛔ 未 push（只 commit）；
  写面只有 `internal/agent/approval/`（突变，五枚文件逐一还原并 md5 回锚）＋`.scratch/wisp/probes/242/v1/**`＋本表；
  `cmd/wisp`／`internal/config`／`internal/risk`／`internal/tools` 零写；无全仓 `go test ./...`；无 `t.Skip`；无放宽断言；
  无 SKIP 读成通过；`frontend/**`／`design/**` 零读零转述；`PLAN.md`／`docs/specs/**`／`thresholds.go`／golden／`allowlist.txt`／三枚冻结件一字未动。
