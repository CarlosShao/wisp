# 票 242 落地腿 `242-r3`（写手）交件说明 — 把没牙的跨卡用例改成真的

- 起手锚 `git log --oneline -1` ＝ **`e6dc79ed`**（派单写的 `a0455bb1` 起手已过期；跑中途 HEAD 又前进两次，收工现量 **`03899fd8`**，见 §7①）。
- 起手 `ls -la .scratch/wisp/probes/242/r3/` ⇒ `No such file or directory`（rc=2）＝**r3 不存在、无死腿**；`probes/242/` 现量确有 `precheck.md`／`r1/`／`v1/`／`v2/`。
  ⚠ `v1/` 全程未读数、未引用、未删（10-03 死腿遗产，无 verdict）。本腿关于"改前照绿"的复现**不取自 v1**，是仓外 overlay 重打的两发（§3 的 OLD2-plain／OLD2-MA）。
- 写面**只有一枚**：`internal/agent/approval/ticket242_binding_test.go`（笔 `df1b962e`，`1 file changed, 104 insertions(+), 8 deletions(-)`，零产码）。
- ⛔ 未动产码、未动票面、未动禁区；分辨两种 denial **不需要**动产码（store 级就能拿到 `misbound`，见 §2／§6），派单 §3 那个"必须走你"的例外**没有被触发**，我也没自行扩射程去改路由级语义。

---

## ① 改前为什么没牙（读码＋"改前照绿"这一发我自己重打）

读码链（现量，非转述）：

| 事实 | 出处 |
|---|---|
| `push` 只建**空** store：`grants: newGrantStore()` | `internal/agent/approval/queue.go:173` |
| `issue` 只在答牌那一路发生：`it.grants.issue(nonce, it.bind)` | `queue.go:358`（`Queue.grantNonce` 体内） |
| `spend` 的绑定比对**在 for 循环体内**：命中行才 `delete` ＋ `equalSecret(stored, bind)` ⇒ `denialMisbound`；扫完全程不命中 ⇒ `return denialSpentNonce` | `approval.go:560-572` |
| ⇒ 递一枚**从未 issue 过**的 nonce 时，绑定那一行**一次都不执行** | 同上的控制流，非注释推断 |
| 改前那一枚（`df1b962e^` 版 `:111`）递的正是字面量 `"leaked-or-guessed-nonce"`，全程不调 `grantNonce`／`allow` | `git show df1b962e^:…` ＝ 120 行那一版 |

复现"改前照绿"（仓外 overlay，`logs/mut-summary-r3b.txt`）：

- **OLD2-plain**（改前那份文件、产码不突变）⇒ `go_test_rc=0`，六枚 TestTicket242 全 PASS，且名册里**没有** seq 那枚（sanity 计数 0 ⇒ overlay 确实把旧件换了进来，不是拿我的工作树凑数）。
- **OLD2-MA**（改前那份文件 ＋ 绑定比对掏空）⇒ `go_test_rc=1`，红的只有 `SpendRejectsForgedBinding`（`:31`）与 `SpendRequiresTheExactBinding`（`:48`）；
  **`--- PASS: TestTicket242ForgedBindingCannotSpendAnotherItemsGrant` 逐名在 pass 名单里** ⇒ 票面 AC#1 唯一声称盯"绑错卡"的那一枚，把绑定层整个掏空**照绿**。复认成功，且与派单／`242-v2` 的判定同数。

一句话：它绿在"这枚令牌不存在"上，不绿在"绑错卡"上；而票面要防的是后者。

---

## ② 改后用例全文要点（`ticket242_binding_test.go:130-196`，行号已逐枚 `grep -n` 复尺）

用例名保持原样（`TestTicket242ForgedBindingCannotSpendAnotherItemsGrant`，函数体起 `:130`），内容整枚重写。造法与断言：

1. **两枚同时活着的卡**：`q.push(A)`、`q.push(B)`（身份不同：`shell.run/task-1` vs `fs.write/task-2`），并**断言** `q.Depth() == 2`（`:148`）、两卡摘要非空且互异（相异那一发在 `:154`）。最后这条是"绑错卡"这件事**成立的前提**——两卡同摘要时 misbound 无从谈起（NEW-SAME 那一发就是被它咬住的，红句在 `:155`）。
2. **真的 issue 过 A 的活令牌**：`nonceA, err := q.grantNonce(itA)`（`:158`＝答牌那一路同一枚构造函数，不是手搓摘要），并且**把"活"本身当断言**：`itA.grants.live() == 1`（`:162`）与 `itB.grants.live() == 0`。这一对是把"改前那个形状"钉死的桩——以后谁再把 issue 摘掉，红句直接指名"这就是 242-v2 退回的那个没牙形"，不会再静默变绿。
3. **A 的活令牌配 B 的绑定摘要 ⇒ 必须红，且必须指名那一支**：
   `d := itA.grants.spend(nonceA, itB.bind)`（`:169`）；断言 `d != denialMisbound ⇒ Fatal`（红句落 `:171`），红句里把**另外两支各自的含义逐字写开**：
   > `AC#1 RED: item A's live grant, spent against item B's binding digest, returned none and not misbound. spent-or-never-live-nonce means the comparison never ran because the row is not in A's store (the shape this case had before, zero power over bindings); missing-nonce means nothing was presented at all; none means the cross-card grant was accepted.`

   ⇒ "绑错卡"与"令牌不存在／已花掉"分辨得开，不是只判 `d != denialNone`。
4. **拒了还要烧牌**：`itA.grants.live() == 0`（`:174` 那发 Fatalf），被拒的 caller 不能拿同一枚 nonce 重试。
5. **同形正控**：再 `grantNonce(itA)` 铸第二枚（`:179`），`spend(nonceA2, itA.bind)` **必须** `denialNone`，且花后 `live()==0`。红句逐字（NEW-MB 就是被它咬住的，落 `:184`）：
   > `AC#1 RED: the positive control was refused too (misbound) - item A's own live grant against item A's own binding digest must spend, otherwise the cross-card refusal above could come from any cause`
6. **判别桩**（`:192`）：`spend("never-issued-value", itA.bind)` **必须** `denialSpentNonce`——把"未知令牌"那一支单独钉一枚，这样本用例永远不可能悄悄退回改前那个形状（NEW-ME 红句在此，见 §3）。

新增一枚独立用例：`TestTicket242BindDigestSeparatesItemsBySequenceNumber`（函数体起 `:203`），见 §4。

---

## ③ 突变表（全部仓外 `go test -overlay`；`C:/Users/swq/242r3-mut/`；读数 `logs/mut-summary*.txt` ＋ 逐名 `logs/mut-*.txt`）

| 形 | 改哪一形（产码怎么坏） | 哪枚红 | 红句逐字（关键那一发） | go_test_rc | 还原后 blob 对拉 |
|---|---|---|---|---|---|
| **OLD2-plain** | 不坏产码，只把测试件换回 `df1b962e^` 那一版 | 无（六枚全绿） | — | 0 | approval.go `67fb1468…`／queue.go `66fec7ae…`／测试件 `8a7b5d2e…` 三枚 **SAME** |
| **OLD2-MA** | `spend` 里加 `bindingEnforced := false`，绑定比对永不执行（＝M-A） | 只有两枚 spend 用例；**跨卡枚 PASS＝没牙复现** | `ticket242_binding_test.go:31`／`:48` | 1 | 三枚 **SAME** |
| **NEW-MA** | 同上 M-A | 两枚 spend ＋ **改后的跨卡枚（核心）** | `:171`（全文见 §2 第 3 条，`returned none and not misbound`） | 1 | 三枚 **SAME** |
| **NEW-MA-FULL** | 同上，整包 | 上述三枚 ＋ `TestTicket259R1DenialNamesMisbound`；其余 68 枚 PASS | 同上 | 1 | 三枚 **SAME** |
| **NEW-MD3b** | `bindDigest` 里 `strconv.FormatUint(seq,10)` → `FormatUint(0,10)`（seq 折出摘要） | **只有** seq 那枚新尺红；`QueuedItemsBind…` 照绿（诚实：它量的不是 seq） | `:208: AC#1 RED: two bind digests identical in every field except the sequence number came out equal - seq is not folded into bindDigest, so a replayed correlation id can land on the earlier card's digest` | 1 | 三枚 **SAME** |
| **NEW-MD3b-FULL** | 同上，整包 | **整包 FAIL＝1 枚**（`242-v2` 现量是 0 枚全绿 ⇒ 这一形从今天有尺） | 同上 | 1 | 三枚 **SAME** |
| **NEW-SAME** | `bindDigest` 丢掉全部身份字段与 args ⇒ **A、B 绑定相同** | `QueuedItems`（`:104`）＋**跨卡枚（`:155`）**＋seq 枚（`:208`） | `:155: AC#1 RED: two distinct pending items share one binding digest - bindDigest stopped separating cards` | 1 | 三枚 **SAME** |
| **NEW-ME** | 扫完全程不命中的 `return denialSpentNonce` → `return denialNone`（未知令牌被放行） | `SpendRejects…`（`:46`）＋**跨卡枚的判别桩（`:193`）** | `:193: AC#1 RED: a value that was never issued on A returned none instead of spent-or-never-live-nonce - this case can no longer tell "bound to the wrong card" from "no such proof", so the denial it names is not evidence` | 1 | 三枚 **SAME** |
| **OLD2-MB** | `grantNonce` 铸侧写外来摘要 `issue(nonce, it.bind+"m")` | 改前那份：**无（六枚全绿）**＝复认 `242-v2` 的 M-B 读数 | — | 0 | queue.go／测试件 **SAME** |
| **NEW-MB** | 同上 | **跨卡枚红在正控那一发**（`:184`）＝v2 记的"242 名下零仪器"这一形，今天在 242 名下有尺 | `:184: AC#1 RED: the positive control was refused too (misbound) …` | 1 | queue.go／测试件 **SAME** |

⚠ 两发**作废并如实记**：
1. 第一台本的 `OLDplain`／`OLD-MA` 用 `HEAD~1` 取"改前版"，而 HEAD 在我跑的过程中被别腿推进（`38919189`／`165cdaa1`）⇒ `HEAD~1` 当时＝我自己的笔 `df1b962e`，那两发实为"新件"读数（`mut-summary.txt` 里它们与 NEW-MA 同数、行号 `:171` 可证）。第二台本改用绝对引用 `df1b962e^` 重打 ⇒ 上表的 OLD2-* 才是有效读数。定式：**取历史版次不许写 `HEAD~1`，要写绝对笔号**（共享工作树里 HEAD 会走）。
2. 第一台本的 `NEW-MD3`／`NEW-MD3-FULL` ＝ **build failed 不是红**：`C:\Users\swq\242r3-mut\approval-MD3.go:418:166: syntax error: unexpected newline in composite literal`（我把 `//` 注释写进了单行 composite literal）。按"编译失败不许冒充红"作废，`mut-NEW-MD3.txt` 原样留着；MD3b 是修好之后的有效读数。

overlay 全程没写过仓：每一发后面都跑 `git hash-object` 对 `git rev-parse HEAD:<f>`，`approval.go 67fb146898dc7b235623bb8ad2fac0e0b2465fe1`／`queue.go 66fec7ae375344d7bc741270b42c4ec5050cf41f`／测试件 `8a7b5d2e9d80675443a98608c1bc7627378ac5a6` **逐枚 SAME**（三台本共 12 组、无一 DIFF）。

---

## ④ 那句恒真句怎么处置的 —— **两头都做：钉住 ＋ 改写说实话**（二选一的"删掉"我没选，因为能钉）

1. **钉住**：新增 `TestTicket242BindDigestSeparatesItemsBySequenceNumber`（说明起 `:197`、函数体起 `:203`）——`corr/task/tool/level/args` 五者全等、**只动 seq**：
   `bindDigest("approval-7","task-1","shell.run","L2",1,args) != bindDigest("approval-7","task-1","shell.run","L2",2,args)`。
   牙齿实测：**MD3b 那一发只有它红**（`:208` 红句逐字在 §3），queue 级那枚照绿 ⇒ 它指的就是它命名的那件事。派单要求的突变②"把 seq 那一味摘掉 ⇒ 指名用例必须红"＝**已满足**；派单要求的另一支"让 A、B 绑定相同"我也打了（NEW-SAME），三枚红、无需具名豁免。
2. **删掉假断言**：`TestTicket242QueuedItemsBindGrantsToTheirOwnDigest` 重放那半的注释（`:92` 起）与红句（`:104`）都改成盘上真形——今天分离两枚卡的是 `Queue.push` 每发新铸的 correlation id（`queue.go:163-165` 把空 CorrelationID 变成 `approval-1`／`approval-2`），**不是 seq**；红句从 `the sequence number stopped separating items` 改成
   `… the digest stopped covering the fields that make the two items distinct (on this path: the correlation id the queue issued)`。
3. **文件头注释里同一句假话也改了**（头块 `:9-27`，"Wording note" 那一段起 `:20`）：派单 §2(b) 只点了 `bindDigest` 那条，但"重放必因 seq 折入而不同"这半句**同时写在用例文件头的判据清单里**（原第 16-18 行）；两处一起改成"seq 确实折在里面，但没有任何 queue 级用例能察觉它被摘掉，seq 由新那枚单独量"。

---

## ⑤ 门禁逐条 rc ＋ 名册作差（全部现跑；`logs/gates.txt`／`logs/mut-*.txt`）

| 尺 | rc | 读数 |
|---|---|---|
| `go test -count=1 -v ./internal/agent/approval/`（改前基线，锚 `e6dc79ed`） | **0** | `base-v.txt`＝89 枚 PASS（含子测试）／**70 个顶层名**／0 FAIL |
| 同上（改后） | **0** | `post-v.txt`＝90 枚 PASS／**71 个顶层名**／0 FAIL |
| 逐名作差 `diff base-names.txt post-names.txt` | **1**（有差＝新增名） | 差集**只两行**：`41a42` ＋ `> --- PASS: TestTicket242BindDigestSeparatesItemsBySequenceNumber` ⇒ **新增红 0、摘尺 0**（`name-diff.txt`） |
| 收尾复跑 `go test -count=1 -v ./internal/agent/approval/` | **0** | `final-v.txt`＝71 顶层 PASS／`grep -c '^--- FAIL'`＝**0**（该 grep 自身 rc=1） |
| `go vet ./internal/agent/approval/` | **0** | 输出 **0 字节**（`gate-vet.txt` bytes=0）＝成功即静默，rc 已单独落行 |
| `sh scripts/d22scan.sh`（⛔ 没用仓根 `go run ./tools/d22scan`） | **0** | 末行 `d22scan: clean - no D22 ban violations; live scope work: … internal/=228 cmd/=38 …` |
| `gofmt -l internal/agent/approval/` | **0** | 名册一枚 `internal\agent\approval\pending_read.go` |
| 行尾符尺 `tr -cd '\r' \| wc -c`（⛔ 未用 `grep -c $'\r'`） | — | `pending_read.go` worktree=**131**／HEAD blob=**0** ⇒ **autocrlf 幻影、非红、非我地界**；`ticket242_binding_test.go` worktree=**0**／blob=**0** ⇒ 我这枚不带幻影 |
| `git status --porcelain`（起手／收尾） | — | 起手 **779** 行、收尾 **756** 行（`git-status-start.txt`／`git-status-end.txt`）；差集里 `grep ticket242\|agent/approval` ＝ **0 命中（rc=1，`git-status-delta-mine.txt` 0 字节）** ⇒ 别腿在飞只登记、未提交、未还原 |
| `git show --stat b6b1d6a4 -- <那一枚测试件>` | **0** | `12 insertions(+), 5 deletions(-)`＝`17 +-` ⇒ 派单"17 行"这一处**我复尺成立**；`4bf7e683` 那 113 行不是现状，本腿审的是 HEAD 版 |
| ⛔ `cmd/wisp`／`./internal/...` 整包 | 未跑 | 按派单约束（两枚腿正在那上面做 blob 对拉） |

`logs/` 里**四枚 0 字节件逐枚具名，没有一枚是漏交**：
- `gate-vet.txt` 0 字节＝`go vet` 成功即**零输出**，那一格真正的凭据是 `gates.txt` 里那行 `VET_RC=0 (bytes=0)`（我这枚件**故意不补字**，补了就不再是逐字输出）；
- `git-status-delta-mine.txt` 0 字节＝差集里 `ticket242|agent/approval` **零命中**（grep 自身 rc=1），"空"就是那一格的读数；
- `mut-console.txt`／`mut-console-r3b.txt`／`mut-console-r3c.txt` 起手为 0 字节＝台本不往 stdout 写，已各补一行指针说明读数落点（现量 169 字节／枚）。
其余证据件逐枚有内容（`ls -la logs/` 现量：`base-v.txt` 14,281／`final-v.txt` 14,417／`gate-d22scan.txt` 23,082／`mut-NEW-MA-FULL.txt` 15,186／`mut-NEW-MD3b-FULL.txt` 14,655／`git-status-start.txt` 35,971／`git-status-end.txt` 34,620／`before-after.diff` 8,723 等）。后缀一律 `.txt`／`.diff`／`.sh`／`.md`，**没有 `.out`**（根 `.gitignore` 有全仓 `*.out`）。

---

## ⑥ 票 242 `AC#1` 今天到底闭合没有 —— **没闭合**（我不翻，也确实不该翻）

先说闭合的那一半（本腿交付的）：**票面 AC#1 的"甲形"这一支现在真的有牙了**。
"两枚同时活着、A 的未花活令牌配 B 的绑定摘要 ⇒ 红，且红在 `misbound` 那一支"＝盘上成立并有突变凭据（NEW-MA 由绿变红；正控同场成立；判别桩独立钉住 `spent-or-never-live-nonce`）。派单 §2(a) 四条要求逐条满足。

不闭合的是**票面原文那一形**（按"原文与转述冲突以原文为准"这一条，我把原文摆在这儿）：

> `242-the-grant-binding-layer…md:18`：**AC#1 真·跨卡：A 的活令牌花到 B 上 ⇒ 必须红**。判据＝造两枚同时活着、correlation id 不同的卡，把 A 的未花令牌递给 **`Native().Allow(B的corr, A的grant)`** ⇒ 要拒且**拒因指名"绑定不对"**（不是"这条路由没在跑"）。

派单 §2(a) 把这一格转述成"拿 A 的活令牌去配 B 的绑定摘要花"＝**store 级**；原文写的载具是**路由级**。我用一枚仓外 overlay 探针（`logs/probe-routed.txt`，⛔ 没进包、`git status --porcelain -- internal/agent/approval/` 空）量了路由那一发的现量：

```
PROBE depth=2 corrA=approval-1 corrB=approval-2 bindA==bindB=false liveA=1 liveB=0
PROBE routed cross-card allow(B.corr, A's live nonce) err=approval: 原生令牌无效（缺失/已用/与本次请求不绑定） err==ErrBadGrant=true
PROBE audit line: approval: GRANT-DENY corr=approval-2 tool=fs.write denial=spent-or-never-live-nonce
PROBE store spend of a value not live on A => spent-or-never-live-nonce
PROBE store spend of a live nonce against another card's digest => misbound
```

⇒ 三句话：**路由级跨卡今天被拒（`ErrBadGrant`），但拒因落在成员扫描那一支 `spent-or-never-live-nonce`，不是票面要的 `misbound`**；`misbound` 在 store 级可达（我最后一行现量），**在路由级不可达**——`allowScoped` 把 `it.bind` 递给自己 store 里同一枚值（`queue.go:415` 那句 `it.grants.spend(nonce, it.bind)`），结构上两端恒等。这一条从今天的 `242-v2` 说的"产码注释级凭据"升级成了**现量**。

因此**整格 AC#1 不闭合**，缺口具名三条：
1. 路由级"拒因指名绑定不对"要落地＝**必须动产码**（乙形：从答复请求重算摘要）。这既是派单 §3 划给编排者裁的例外，也和票 259 刻意保留的设计相撞——`grantDenial` 的注释逐字写着 API／UI 侧"unchanged on purpose… no caller can read a grantDenial from out here"。**我没动手，报回来给你裁。**
2. 票面自己那条前置（`cmd/wisp/subagent_selfapproval_197_test.go:109` 的 `CorrelationID: taskID`）在**我的禁区**，未做；本腿的两枚卡靠包内 `q.push` 造，与票面点名的那个载具不是同一枚。
3. 路由级之外，票面 AC#1 的"正控"是 D2（`spend` 不再比对）——这一形本腿**已按票面要求做成"改前绿、改后红"**（OLD2-MA 绿／NEW-MA 红），这一条不再欠。

顺带一条**记账更新**（不是新缺陷）：`242-v2` 的 M-B 那格原写"242 名下全绿、性质有守手但不在 242 名下"。本腿现量 **NEW-MB 让跨卡枚在正控那一发红**（`:184`）⇒ M-B 今天**在 242 名下有尺了**。

---

## ⑦ 我没量到的／要顶回你的

1. **HEAD 号又过期两次**（派单 `a0455bb1` → 起手 `e6dc79ed` → 收工 `03899fd8`）；不是我写歪，是共享工作树号在走。**派单里所有状态级读数（含 `b6b1d6a4` 那 17 行）都请当我这种各腿自取**——17 行这一把我复尺成立（`12 +/5 -`）。
2. **派单 §2(a) 与票面 `:18` 不同形**（store 级 vs `Native().Allow` 路由级），本条按"以原文为准"顶回：我照 §2(a) 交付了 store 级那一形，**但票面原文那一格仍不闭合**，判语在 §6。要不要改票面文字（＝契约变更、须人工批准）或补产码（乙形）归你裁。
3. **"分辨两种 denial 必须动产码"这一支没有被触发**：`denialMisbound` 与 `denialSpentNonce` 在 **store 级**就分得开（现量两行见 §6 探针）。需要走你的只有**路由级**那一支。
4. **`sameDigest` 那枚零调用者的死桩我没摘**：现量 `grep -rn "sameDigest("` 全仓只有定义 `:214`、零调用点（票面 :55 那条"build/vet 的 rc=0 对内部有效性是零证据量"我复认）。它不是尺、摘不摘都不咬人，本腿按"最小射程"留着没动，归你处置。
5. **没取数的格**：①`cmd/wisp`／`./internal/...` 整包读数（约束不跑）；②票 259 那三枚能力尺在 M-A／M-B 之下的完整名册作差——我只逐名确认 `TestTicket259R1DenialNamesMisbound` 在 NEW-MA-FULL 红、其余 259 枚绿，没做全册差集；③AC#2 那一格（出向读面／`PanelAPI`）本腿零触碰，`242-v2` 的 M-F／M-G 两发我**未复跑**，它的"AC#2 成立一半"判语在此程仍是腿报、非我复尺。
6. **探针那一发差点变成死腿**：第一版探针 `defer mu.Unlock()` 与后续 `q.push` 的回调撞锁 ⇒ `go test` 挂到超时（不是红、不是读数）。修掉重跑才拿到 §6 那五行；两份台本与失败描述都在 `logs/`，`.scratch/wisp/probes/242/r3/logs/probe_routed_source.go.txt` 是探针源码副本（它**不是**在册尺、不在包里）。定式建议：**overlay 注入的临时探针凡带 logger 回调，锁必须在辅助函数内成对释放**。
