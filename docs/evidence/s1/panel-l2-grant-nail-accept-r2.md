# panel-l2-grant-nail-accept-r2 —— 非实现者验收：`9d85789…598620e` 那五枚补强有没有把 r1 造出来的洞闭上

- 时刻 / 锚点：2026-09-25 11:1x +08 开工。进场 HEAD `39b553e`，取数时 `dev` 在 **`0b95e9e`**，
  写这行时 `1f128df`（共享树，别家每分钟在提）。⇒ **本件每个数旁边都写了自己那一发的 sha**，
  任何"对不上"先比 sha，再比谁错。
- 被验物：`internal/panel/l2_grant_boundary_test.go`，本程锚定 blob
  **`3b7a2cb1d4105df1f279f70c605ccf6a71f7868a`**（工作树＝HEAD，逐枚同值，现量于 §0.1）。
  实现者自己的记录＝`docs/evidence/s1/panel-l2-grant-nail-fix-r2.md`（381 行，止于 §6.3）。
- 上一轮＝`docs/evidence/s1/panel-l2-grant-nail-accept-r1.md`（511 行，总裁 §7：
  **退回（附条件入账）**，最小闭合三件 (a)(b)(c)+F-5，另两件 (d)(e)）。
- 身份：**本程不是实现者**。未写一字节生产码、未改任何既有文件、未 push、未 `add -A`。
  可写路径只有本件一枚。
- 地界：未碰 `tools/d22scan/**`（兄弟验收在里）、`internal/observe/**`（别家在测、且它起了
  `golang:1.27` 容器）、`frontend/**`、`design/**`、`docs/PLAN.md`、`docs/specs/**`、`docs/reports/**`。
  **未跑 `go test ./...`**（只跑 `./internal/panel/` 与既门 `sh scripts/d22scan.sh`）——整树读数此刻不可归因。
- 方法：实现件每一条读数都当**待复测的断言**。变异只落在**仓库外副本**
  `D:\tmp\panel-l2-nail-accept-r2\headcopy\`（`tar` 排除 `.git`/`third_party`/`build`/`node_modules`，
  1.95 s／155 M），仓库树在任何时刻都没处于变异态。还原用 `cp`，每发之后当场
  `git hash-object` 复算 `bridge.go = d2cd6362…`。未 `checkout`/`reset`/`stash`/`clean`，**未 `rm`**
  （临时件只建不删，用完的探针 `mv` 进 `removed/` 留着）。

---

## §0 闸门 + 那枚被轮次上限掐死的实现程到底留下了什么

### §0.1 三门与两数（全部本程真跑，锚点 `0b95e9e`）

```
gofmt -l internal/panel/        ->  空                     rc=0
gofmt -l .                      ->  空                     rc=0
go vet  ./internal/panel/       ->  空                     rc=0
go test ./internal/panel/ -count=2 -v   ->  rc=1
RUN=196  TOPPASS=102  SUBPASS=92  FAIL=2  SKIP=0  ^panic:=0  distinct=98
--- FAIL: TestC21DesignTokensFourWayAgree (0.00s)     x2      （行 249 / 510）
```

- `FAIL=2` **就是**那两枚先存在留红（`A208③` 的 C21），本程未修、未跳、未放宽 ⇒ 它必须继续红。
- `^panic:=0` ⇒ 没有"一枚 panic 吞掉同包几十条读数"那种事（本包每条读数都取到了）。
- `SKIP=0`；全包 `t.Skip` 仍只有 `attachments_test.go:191` 那一枚先存在的
  `t.Skipf("filesystem cannot hold a 2 GB file here")`，不在本件地界。
- `sh scripts/d22scan.sh` -> **rc=0**；口径读数（现跑，锚点 `0b95e9e`）：
  `bans #1-5 internal/=203, cmd/=22, ban #6 frontend/=43, ban #7 internal/tools/=18,
  ban #8 design/=32, frontend/=43, internal/=407, cmd/=39`。
  ⚠ `frontend/=43` 这一枚**会长**（第三家在写 `frontend/**`），本件只把它当"我这一发的分母"，
  不拿它主张任何事；bans #1–5 **不含 `_test.go`**，所以**不得**声称"这道门会抓到测试里的墙钟/裸 goroutine"，
  而本批改的**恰好是一枚 `_test.go`** ⇒ 门对本次改动只有 ban #8 那一枚有分母。

### §0.2 名册：98 vs 90 —— 90 那一版我自己装回副本里跑出来，没有引用它的清单

派单要我"与 r1 的 `90` 做差集"。r1 的 §0.3 只给了**新增的 14 枚**逐名，**没存全名册**
（"逐名见本件末"那句落空——本件末尾是 §10 交件态，没有 90 行名单）。⇒ 我不用它的清单：
按 r1 声明的同一把尺（`grep -o '=== RUN…'` 取第 3 字段 → 去 `#N` 后缀 → `sort -u`，**含子测试**），
在**仓库外副本**里把 `d88c356` 那版钉文件（`git cat-file blob` → blob `dc4e0ccf…`）装回同一棵树跑一遍。

```
go test ./internal/panel/ -count=2 -v   （副本 + d88c356 版钉，out/before-90.txt）
RUN=180  TOPPASS=100  FAIL=2  SKIP=0  ^panic:=0  distinct=90
   ↑ 与 r1 §0.1 它自己那发"后态"逐枚同值（180/100/2/0/90）⇒ 我这把尺与它那把尺同形
```

```
distinct:  before(d88c356 版钉) = 90   →  after(HEAD 版钉 3b7a2cb1) = 98
comm -23 before after = **0 枚**   ← 名册一枚都没消失（没有"改着改着测试变少"）
comm -13 before after = **8 枚**   ← 全部来自本批，逐名如下（out/roster-added.txt）
```

| 新增名 | 来自哪一味 |
|---|---|
| `TestJSONKeyDerivationAgreesWithEncodingJSON` | (c′) 新顶层测试（`func Test` 5→6，diff 里只有 1 行 `>`、0 行 `<`） |
| `…AgreesWithEncodingJSON/i_the_AST_key_rule_matches_encoding/json's,_spelling_by_spelling` | (c′) |
| `…AgreesWithEncodingJSON/ii_the_AST_list_and_the_reflection_list_are_one_list` | (c′) |
| `…AgreesWithEncodingJSON/iii_the_decoder_is_the_third_vote_and_the_verdict_words_are_not_bindable` | (c′) |
| `TestPlantedGrantWiringGoesRedInASnapshot/E_a_guard_route_reached_through_a_var_or_a_concatenation_is_named,_not_dropped` | (a) 的植物 |
| `TestPlantedGrantWiringGoesRedInASnapshot/F_a_route_named_outside_the_guard_still_reaches_the_pool` | (b) 的植物 |
| `TestPlantedGrantWiringGoesRedInASnapshot/G_an_inbound_envelope_with_a_renamed_route_key_is_still_seeded_from_its_decode` | (c) 的植物 |
| `TestPlantedGrantWiringGoesRedInASnapshot/H_a_verdict_tagged_without_a_name_is_still_found_under_its_field_name` | (c′) 的植物 |

⇒ 8 枚**逐枚可归到某一味**，没有一枚是"来路不明的增长"，也没有一枚是**抄来的名单**
（E/F/G/H 是 facet 4 在自己造的快照树里种的植物，不是硬编码断言）。
跑完我把副本钉文件 `cp` 回 `3b7a2cb1…`（与 HEAD 同值，现量在同一发里）。

**两版"改前"我都跑了**，因为 r1 判的那枚 blob 不是 `d88c356` 而是它后面 `660ffa1`（+4/−2 行）：

```
660ffa1 版钉 = 5b618d14…(r1 与 fix-r2 记的同值)  distinct = 90
d88c356 版钉 = dc4e0ccf…                          distinct = 90   ← 同一枚名册
```
⇒ `660ffa1` 那次 4/2 行没改任何测试名，两把"改前尺"给出同一个 90，本件的 +8/−0 与谁比都成立。

⚠ 口径提醒：`d88c356..HEAD` 在 `internal/panel/**` 上**只动了这一枚文件**（§6 契约轴现量
`1068 49`，一枚路径），所以这 90↔98 是**同一棵树换一枚文件**的差集，可比性是构造出来的，不是假设的。

### §0.3 那枚死掉的实现程：留下的与没留下的（本程自己核，不看它的通知）

`docs/evidence/s1/panel-l2-grant-nail-fix-r2.md` = **381 行，最后一枚标题是 §6.3**
（`grep '^## '` 现量：§0 §1 §2 §3 §4 §5 §6 各一枚，**没有 §7/§8/§9**）。

| 派单要它做的 | 盘上有没有 | 我怎么量的 |
|---|---|---|
| 五枚 commit 各自只带自己路径 | **有** | 逐枚 `git show --name-only`（§6.1 表） |
| 码（(a)(b)(c)(c′)(d)(e)） | **有** | 生产码零字节 + 测试件 +1068/−49（§6 契约轴） |
| 改前六发种子复现（§0.3） | **有**（本程独立复跑，见 §1） | `out/seed-*.txt` |
| **改后每味的反向对照** | **没有** | 它自己最后一句"Now the reverse controls. First, add two new seeds to the bench:" 之后无 §7 |
| **总裁 / 逐格判** | **没有** | 无 §7、无 `^总裁` |
| **"本程没测什么"** | **没有** | 无 §8 |
| **零术语那一段** | **没有** | 无 §9 |

⇒ **本件 §2 就是它没做完的那格，由我自建**（不是我替它补，是它被派单要求、而它没交的东西）。

### §0.4 它自己文件里的悬空引用（这一格是本仓那个已登记缺陷的名字）

已登记的缺陷名："**预先引用尚未产出的读数的文档 = 假绿前身**"。现量 `grep -n '§7\|§9'`：

```
:4   … 工单来源 … §7.3 / §7.4          -> 指 accept-r1，那两节存在（:417 是 §7，最小集合是它下面的编号 3、4）= 有效
:73  "改完之后同样的种子必须当场红（§7）。"   -> 指本件 §7 = **不存在**  ← 悬空，而且正是它没写完那一格
:81  "…现核（读数在 §9）。"                -> 指本件 §9 = **不存在**  ← 悬空
:8   "生产码零字节（见 §9 的 git diff --numstat）" -> 指本件 §9 = **不存在**  ← 悬空（派单点名的那一枚）
:131 "验收 §7.3(b) 的原话…"                -> 指 accept-r1 = 有效
:201 "验收 §7.3(c) 的判据选择成立…"          -> 指 accept-r1 = 有效
```

⇒ **三枚自指悬空（`:8`、`:73`、`:81`），两枚外指有效**。除此之外我把其余被点名的符号也扫了一遍：
`:604`/`TestComposerMethodNamesMatchFrontend` 这类"指向不存在的测试"的旧病，本批**已把本件里那一处修掉**，
现量全仓只剩 `bridge.go:33` 一枚（先存在，见 §5.2）。

**判这一格**：缺陷成立，形状与登记的那条**同类**，且 `:73` 那一枚不是随手写错——它指向的正是
"反向对照"那一格，**读者会以为反向对照跑过了**。⇒ 定性：`panel-l2-grant-nail-fix-r2.md` 是一份
**码已落地、读数半份**的交件；它的 §0–§6 里凡本程复算过的都对得上（见 §5.1），
但它**没有总裁、没有反向对照、没有"没测什么"**，所以它**不能**被当成"这一批已经验收过"凭据。
逐格判在 §5.3 与 §7。

---

## §1 生死格：M13 现在红不红 —— **红**，且每枚种子都点名到 `file:line`

### §1.1 台子

- 变异只在 `D:\tmp\panel-l2-nail-accept-r2\headcopy`（`tar` 排除 `.git`，副本
  `bridge.go=d2cd6362…`、钉 `3b7a2cb1…` 与仓库/HEAD **逐字节同值**，四锚现量）。
- 副本基线（`out/copy-baseline.txt`，`-count=2 -v`）：`RUN=196 TOPPASS=102 FAIL=2 SKIP=0 panic=0`
  —— 与仓库那一发（§0.1）**同值** ⇒ 我这把尺在副本上读的"没红"不是副本自己坏了。
- 每一发一条命令内完成"种 → 跑 → 还原"，还原后当场复算
  `bridge.go = d2cd6362ecc6941a0cee58073a2bf8ab6669a590`（**九发九次同值**，逐行在
  `out/seed-*.txt` 旁的 `run-seed.sh` 回显里），仓库 `internal/panel/` 全程 `git status --porcelain` **空**。
- ⚠ 我自己踩过一次"假种"：第一版 `mutate.py` 里 `%q` 不是合法 Python 格式符，
  M6/M7/M8/M13 那四发**根本没落地**。如果我只看"跑出来几枚红"，会把"种失败"读成"钉失灵"。
  现在脚本对每发断言 `delta_bytes` 且 `out == SRC` 时报 `BYTE-IDENTICAL-TO-CLEAN (plant failed!)`，
  并逐发 `gofmt` 过一遍语法。⇒ 这条写出来给下一位：**验收方自己的台件也要有正控**。

### §1.2 六发种子（`-count=1 -v`，锚点 `0b95e9e`，原文 `out/seed-<名>.txt`）

| 发 | 形状（都在副本的 `bridge.go` 里） | 本件红的钉 | 关键红句（逐字） | 判 |
|---|---|---|---|---|
| `M2` | 守卫 case 里放一枚包级 **var** `dynRouteM2 = "panel.review.allow"` | 4 枚 | `:1133`/`:1235`/`:1438` 三处同句 `the inbound route guard under … has 1 case label(s) this instrument cannot resolve to a string:` + `:1624` 快照预检 | **红** |
| `M3` | 守卫 case 里放拼接串 `"panel.review." + routeTailM3`（两半各自都不像路由） | 4 枚 | 同上，且点名 `bridge.go:102: "panel.review." + routeTailM3` | **红** |
| `M4` | **第二枚函数** `extraGateM4` 自己的 switch 答 `"panel.review.allow"`，守卫末尾 `return extraGateM4(m)` | 2 枚 | `:1179` `knownComposerMethod answers "panel.review.allow" (written at bridge.go:108) but the guard's own case list … INCOMPLETE, not clean` + `:1182` 审批形状那一判 | **红** |
| `M6` | 新入站封套 `{Cmd json:"cmd"; Outcome json:"outcome"}` + 真 `json.Unmarshal` | 3 枚 | `:1237` `bridge.go:122: inbound envelope acceptM6Envelope can bind the JSON key "outcome"` + `:1444` `decode destination acceptM6Envelope has no reflection twin in inboundTypeRegistry` | **红** |
| `M7` | 新封套**绑 `method`** 但 `Outcome string json:",omitempty"` | 3 枚 | `:1237` `… acceptM7Envelope can bind the JSON key "Outcome" … Outcome string `+"`"+`json:",omitempty"`+"`"+`` | **红**（r1 时零枚红） |
| `M8` | 同一形状放进**正主** `ComposerRequest` | 4 枚 | `:1218`（反射半）`ComposerRequest.Outcome binds "Outcome"` **加** `:1237`（AST 半）`bridge.go:70: inbound envelope ComposerRequest can bind the JSON key "Outcome"` | **红**，且**两台仪器不再分家** |
| **`M13`** | **M4 的第二条链 + `cmd`/`outcome` 封套合体** | **4 枚** | `:1179` + `:1182`（路由半）与 `:1237 bridge.go:131: inbound envelope acceptM13Envelope can bind the JSON key "outcome"`（封套半）**同发同红**，另 `:1444` 与 `:1625` 两条引信叫 | **红** |

⇒ **r1 的决定性一发（全包零反应）现在四处叫**：一条链把 `panel.review.allow` 从"另一枚函数"里答出来，
被 (b) 的字面量池点名；一条 `cmd` 键的封套把 `outcome` 收进来，被 (c) 的 decode 种子点名。

### §1.3 M13 的**运行期真值**（我打印的，不是我推理的）

探针 = 副本里一枚临时 `_test.go`（`goSourceFiles` 排除 `_test.go`，不进字面量池、不进 AST），
跑完 `mv` 进 `removed/`（**只建不删**）。

```
负控（干净树 + 探针，out/probe-clean.txt）
  PROBE-CLEAN ParseComposerRequest err=panel: composer request refused: 方法 "panel.review.allow"
              不是面板 composer 通路的能力入口 accepted=false method="panel.review.allow"
  PROBE-CLEAN knownComposerMethod("panel.review.allow")=false
  PROBE-CLEAN knownComposerMethod("panel.review.grant")=false
  PROBE-CLEAN knownComposerMethod("panel.approval.request")=false
  PROBE-CLEAN knownComposerMethod("panel.allow")=false
正控（M13 树 + 探针，out/probe-M13.txt）
  PROBE ParseComposerRequest err=<nil> accepted=true method="panel.review.allow"
  PROBE knownComposerMethod(panel.review.allow)=true
  PROBE envelope cmd="panel.review.allow" outcome="grant"
  PROBE VERDICT: a panel-supplied allow reached Go on an answered route and the
                 verdict came off the wire
```

⇒ 我**逐字复现了 r1 §6 那三行决定性读数**（`accepted=true`、`method="panel.review.allow"`、
`outcome:"grant"` 真的从线上被读进 Go），差别只有一条：**同这棵树上的这枚钉，四处叫了**。
两向都有 ⇒ "钉没响"不能怪探针没生效，"钉响了"也不能怪树本来就在响。

### §1.4 §1 总裁：**这一格成立**

派单的判据是"任何一发还绿 → 这一格退回，明说"。**没有一发还绿**：
六发全部至少 2 枚红，且每枚红都点名到 `bridge.go:<行>` 或"读不懂的那枚表达式"。
最硬的一枚证据不是"红了多少"，是 **M3**：它的名字 `panel.review.allow`
**从来没有以一枚 route-shaped 字面量的形式写进包里**（`"panel.review."` 尾段为空、`"allow"` 只有一段），
所以 (b) 的池看不见它 —— 红它的是 (a) 自己。⇒ (a) 不是 (b) 的影子（见 §2）。

## §2 反向对照（那枚死掉的程唯一没做完的一格，本程自建）

### §2.1 怎么"摘掉一味"

`denail.py` 对**钉文件本身**做单点摘除（每处摘除都断言锚点、产出零变化就拒绝报告），
`mutate.py` 对副本的 `bridge.go` 种种子，两向都还完当场复算 hash。摘除口径：

| 摘除 | 具体动作 | 是不是"只动这一味" |
|---|---|---|
| `noA` (a) | 删 `instrumentBlindnessProblem` 里 `if len(pkg.dropped) > 0 {…}` 那一支 | 是，−5 行 |
| `noB` (b) | 删 facet 1 里从 `pool := routeNamePool(pkg)` 到 `answeredOutsidePool` 整块 | 是，−15 行（`routeNamePool`/`poolJudgedByRealGuard` 函数体留着，facet 4 的植物 F 仍走它们） |
| `noC` (c) 局部 | `inboundEnvelopes` 的种子退回"绑 `method` 键"，并删两枚 decode 引信 | 半摘：`classifyDecodes` 仍在填 `pkg.inboundSeeds` |
| **`noC_full` (c) 整味** | 上面那些 **+** `pkg.inboundSeeds[d.TypeName] = true` 停填 | **是**——(c′) 那枚防腐引信读的就是 (c) 喂的数据 |
| `noCp_rule` (c′) 规矩 | `jsonOr` 退回 `if tagged { return jsonName }`；反射侧 `if tag == "-"` 整枚比法退回 | 是 |
| `noCp_test` (c′) 测试 | 整枚 `TestJSONKeyDerivationAgreesWithEncodingJSON`（含 `:1444` 防腐 `t.Fatalf`）删掉 | 是，−113 行 |
| `noABCp` | 四味全摘 | 用来验"我能不能回到 r1 那发零反应" |

**先给两枚台件自己的失败读数**（不藏）：
① `noCp_test` 第一版我拿 `"// TestJSONKeyDerivationAgreesWithEncodingJSON"` 当起点，命中了
`:1303` 那枚**别的**注释（`inboundTypeRegistry` 的文档注释也点名它），切错了边界 ⇒ 四发全是
`builderr=1 nail_red=0`。**如果我按"零枚红"读数，就会把"我把包切坏了"报成"这味是镜子"。**
② `mutate2.py` 的 MEXT 那发我留了个 `assert "acceptMEXTDecision"`（函数后来改名 `acceptMEXTDelay`），
它**在写文件之前**抛异常，于是那一行 `MEXT` 的读数其实是 MDEC 的重播。两处都由**后验断言**
（`builderr` / `assert` 在写盘前）当场抓住，重跑后的数在下面。

### §2.2 摘一味 → 哪一发种子不再红（`-count=1 -v`，锚点 `0b95e9e`，原文 `out/rc-*.txt`、`out/full-noC-*.txt`）

`nail_red` = 本件 6 枚顶层测试里红的枚数（C21 已排除）；"只剩植物"= 只有
`TestPlantedGrantWiringGoesRedInASnapshot`（它测的是快照树，不是真树）红 ⇒ **真树全绿**。

| 摘除 | CLEAN | M2 | M3 | M4 | M6 | M7 | M8 | M13 |
|---|---|---|---|---|---|---|---|---|
| **无（交付态）** | 0 | 4 | 4 | 2 | 3 | 3 | 4 | **4** |
| `noA` (a) | 只剩植物 E | 2 | **只剩植物** | 2 | 3 | 3 | 4 | 4 |
| `noB` (b) | 0 | 4 | 4 | **只剩植物** | 3 | 3 | 4 | 3 |
| `noC` (c) 局部 | 只剩植物 G | 4 | 4 | 2 | 2 | 3 | 4 | 3 |
| **`noC_full` (c) 整味** | 只剩植物 G | — | — | — | **只剩植物** | 2 | 4 | 2 |
| `noCp_rule` (c′)规矩 | 3 | — | — | — | 3 | **2（ban 那枚绿）** | 4 | 4 |
| `noCp_test` (c′)测试 | 0 | — | — | — | 2 | 2 | 3 | 3 |
| `noCp_both` (c′)全摘 | 只剩植物 H | — | — | — | — | **只剩植物** | 3 | 3 |
| **`noABCp` 四味全摘** | — | **只剩植物** | — | **只剩植物** | **只剩植物** | — | — | **只剩植物** |

### §2.3 逐味判（本仓那句定义 + 那句"公式会骗人"的补语）

> 定义：摘掉任意一味，都存在一发变异从此打不红 ⇒ 那味承重。反面：摘掉它在**任何输入上**都不改变任何
> 可观察读数 ⇒ 那是镜子，要说破，不许它当钉交付。
> 补语（兄弟程刚立的）：公式可能给误导答案——还要问"摘掉它到底有没有任何**外部可见**读数变过"。

- **(a) 承重，且是 M3 的唯一目击者。** `noA` 下 M3 真树全绿（只剩植物 E），
  而 M3 的名字 `panel.review.allow` **从未以 route-shaped 字面量写进包里**（`"panel.review."` 尾段为空、
  `"allow"` 只有一段）⇒ (b) 的池结构上看不见它。M2 在 `noA` 下仍被 (b) 咬住 ⇒ (a)(b) 对 M2 冗余、
  对 M3 不冗余。**不是镜子**：摘掉它有一发真树变异由红转绿。
- **(b) 承重，且是 M4 的唯一目击者。** `noB` 下 M4 真树全绿；M13 在 `noB` 下仍红（(c) 兜住封套半），
  但 `:1179`/`:1182` 那两句路由侧的红**只有 (b) 会报**。⇒ 对 M13 两味各够；这不是缺陷，是要写清的冗余。
- **(c) 承重——但只有在"整味摘"时才看得见。** `noC`（局部摘）下 M6 仍红，因为 (c′) 的防腐引信
  `:1444` 读的是 (c) 喂的 `pkg.inboundSeeds`；把那味数据也停填（`noC_full`），**M6 真树全绿**。
  ⇒ 结论要带这句：**"摘掉 (c) 还有人兜"是一个半摘伪象**，两味之间有单向依赖（(c′) 的引信靠 (c) 的数据活着），
  记进 §7 的缺陷面：将来有人只删 (c) 的判据而留着 `classifyDecodes`，会以为"没人管也不碍事"。
  另一发同向证据：按 `:1444` 那句话**把新类型注册进去**之后（`REG+`），摘掉 (c) 仍只剩
  `TestNoInboundEnvelope` 一枚由绿转红↔红转绿之差（`REG+NONE|M6`＝3 枚红含 ban 那枚，
  `REG+noC|M6`＝2 枚红不含 ban 那枚）⇒ (c) 的红不是靠"你忘了注册"这一脚。
- **(c′) 承重，但两半互相顶班；且它的"M8 那一半"是白盒。**
  `noCp_both` 下 M7 真树全绿 ⇒ 整味承重。`noCp_rule` 下 **M7 的 ban 断言（`TestNoInboundEnvelope`）转绿**、
  只剩那枚三向测试红 ⇒ "AST 算对键名"是 (c′) 规矩半的行为后果，量出来了。
  `noCp_test` 下**没有任何一发**转绿、CLEAN 也全绿 ⇒ **那枚新测试自己不拦任何今天的变异**：
  它拦的是"jsonOr 以后再被改坏"（`noCp_rule|CLEAN` 三枚红就是它在工作）与"新增入站类型不许悄悄不进反射对照"。
  ⇒ 判：**测试半 = 承重但白盒（instrument-integrity，不是 ban-catching）**，措辞只许这么写；
  它**不是镜子**，因为我量到了它的触发输入（摘规矩 ⇒ 它响）。而 **M8 这一发**：摘掉 (c′) 规矩后包体颜色不变
  （反射半一直红）⇒ (c′) 对 M8 买的是"**两台仪器不再分家**"（红句从只有 `:1218` 变成 `:1218`＋
  `bridge.go:70` 那条带行号的 AST 诊断），**不是新增一次拦截**。这句必须写在明处，否则 (c′) 就领了它没挣的红。
- **四味合起来＝r1 那发的完整复现。** `noABCp` 下 M2/M4/M6/M13 全部只剩植物、真树全绿 ⇒
  这四味**恰好**是闭上 r1 那个洞的东西，一枚不多一枚不少（(d)(e) 是文档，不在这个判据里）。

### §2.4 顺手钉出来的一枚新缺陷（不是我替谁加条件）

`classifyDecodes` 把"解不到同包 struct"一律算 `holes` ⇒ **任何**入站外的合法 decode 都会把
`requireReadableInstrument` 踩响。我用两发**合法形状**量了代价：

```
MDEC（`var m map[string]any` + json.Unmarshal，最常见的"先看一眼再说"写法）
  -> 本件 4 枚红：:1133 / :1235 / :1438 / 快照，全同一句
     "a JSON decode destination in 1 place(s) cannot be enumerated: bridge.go:123:
      decodes into &m (type \"map[string]any\"), which is not a same-package struct: ...
      Judge it in a test that can see that type, or route the bytes through a same-package
      struct - do not let this file report the boundary clean"
MEXT（`json.RawMessage`，"把字节留着稍后再解"的正规写法，零新依赖）
  -> 同样 4 枚红，红句 bridge.go:125 decodes into &d (type "json.RawMessage")
```

⇒ 方向要说准：这是 **fail-closed（宁可吵也不放行）**，消息自解释、给的出路可执行，
不是"误伤好人"的那种红；但它意味着 **`internal/panel` 里今天不能出现第二枚 decode**，
除非它落到同包 struct。这条约束是本批新加上生产码的，**派单与 r1 都没写过它**，
r2 的文件头 `DOES NOT COVER` 也没写它。⇒ 记为**文档级缺陷**（见 §7 F-R2-1），不改判。

---

## §3 新尺子的误伤审计（本程认为最值钱的一格）

固定问法：**判据一收紧，就必须答"哪一类**今天本来就合法**的数据今天会被拒、回退路径是什么"。**

### §3.1 (b) 今天到底扫进多少东西 —— 分母是我量出来的，不是推的

同一棵副本树上放一枚 `_test.go` 审计探针（`_test.go` 被 `goSourceFiles` 排除，
探针自己进不了它正在量的池），跑完 `mv` 进 `removed/`。原文 `out/audit-pool.txt`、`out/audit-filter.txt`。

```
池与判定（真树，HEAD 版钉）
AUDIT answered=4 pool=5 literals=5 consts=7 structs=18 decodes=1
POOL index.html                  answeredByRealGuard=false  sites=[assets.go:29 const EntryFile]
POOL panel.attachment.add        answeredByRealGuard=true   sites=[bridge.go:37 const MethodAttachmentAdd]
POOL panel.message.send          answeredByRealGuard=true   sites=[bridge.go:38 const MethodMessageSend]
POOL panel.mode.request          answeredByRealGuard=true   sites=[bridge.go:35 const MethodModeRequest]
POOL panel.workspace.request     answeredByRealGuard=true   sites=[bridge.go:36 const MethodWorkspaceRequest]
```

```
过滤器自己的分母
FILTER files=8 stringLiterals=295 withDot=16 routeShapedKept=5 rejected=11
REJECTED approval.go:21  "github.com/CarlosShao/wisp/internal/risk"
REJECTED assets.go:25    "github.com/CarlosShao/wisp/frontend"
REJECTED composer.go:37 / composer_handlers.go:47 / workspace.go:31  同上那条 risk 导入路径
REJECTED assets.go:79   "./"      assets.go:80 "."    assets.go:83 ".."
REJECTED assets.go:105  "."       assets.go:166 "."   attachments.go:260 "."
```

⇒ **今天被 (b) 纳进射程的"合法但非路由"数据 = 恰好一枚：`index.html`**
（`internal/panel/assets.go:29` 的 `const EntryFile`，一个真正的文件名）。
它过 `routeShapedName` 是因为 `词.词` 这条形状**天生分不清"路由"与"带扩展名的文件名"**——
`panel.ts`、`wisp.db`、`manifest.json`、`golden.sse` 只要将来被写成一枚生产码字面量就会进池。

**它被拒了吗？没有。** 判语只有一句：
`if !knownComposerMethod(name) { continue }`（`:1024`）——**在池里不等于被起诉**，
被守卫应答才成立，而守卫应答的判据与它是不是字面量无关。⇒ 这条判据的"回退路径"今天用不上，
因为它**不拒任何东西**；代价是"多算几次哈希"。这一条我在 §3.3 用一发反证钉死。

**过滤器不滤掉什么**（说破比留白好）：`词.词`、`词.词.词`、段内允许 `_`、数字、**非行首**的 `-`。
所以 `a.b-c`、`photo.png`、`shell.run`、`config.toml` 全过；被滤掉的 11 枚今天**全是因为带 `/`、
`\`、`%` 或整枚就是一个点**——也就是说它滤的是**路径**，不是**名字**。
⚠ 两处今天没被扫到、但**换个写法就会进池**的东西：`shell.run`（`approval.go:31`）与
`"C:\dir\photo.png"`／`"photo.png"`（`attachments.go:315-316`）——它们在**注释**里，AST 看不见，
`noB`/`noA` 系列也没报错。⇒ **注释天然豁免**这条在 (b) 上成立，方向是"少扫"，与 ban #8 的注释豁免同族。

### §3.2 (c) 今天的射程，以及"加第二枚 decode"会变成什么

```
SEEDS（decode 派生）      [ComposerRequest]
INBOUND（向下闭包）        [AttachmentPayload AttachmentRef ComposerRequest]
OLD-criterion seeds（旧判据）[ComposerRequest]        ← 两把判据今天外延相同
DECODE bridge.go:79 dst=&r type=ComposerRequest       ← 全包唯一一枚 decode
FIELD ApprovalCardView.SessionOverrideBlocked json="sessionOverrideBlocked" at approval.go:54
FIELD ApprovalCardView.DecidedBy              json="decidedBy"               at approval.go:58
```

- `decodes=1` ⇒ 换尺今天**零成本**（fix-r2 §4.1 那句"外延完全相同"我独立复现，两把判据同一个格）。
- **r1 警告的"别用查所有 struct 那版"我复算成立，并且第二枚反例是真的**：
  全包今天恰好 **2 枚**字段名带审批词（`sessionOverrideBlocked` 命中 `override`、`decidedBy` 命中 `decide`），
  两枚都在 `ApprovalCardView`（**出站**渲染视图，既不是 decode 目的、也不绑 `method`）⇒
  现在的两把判据都够不到它们；若判据换成"包里所有 struct 都查"，**今天就是两枚红**。
  r1 只点了 `decidedBy` 一枚，fix-r2 §4 补了第二枚；我数到的仍是这两枚，没有第三枚。
- **加第二枚 decode 会发生什么**（我造了两发合法形状，见 §2.4）：
  *新同包 struct* ⇒ (i) 进 `inboundSeeds`；(ii) 若带审批键 ⇒ `TestNoInboundEnvelope` 点名 `file:line`；
  (iii) 若忘了登记 ⇒ `:1444` 那枚 `t.Fatalf` 先叫，**红句直接把出路写出来**
  （"drop it from inboundTypeRegistry instead of dropping the check" / "has no reflection twin"）。
  *非同一包的目的地（`map[string]any`、`json.RawMessage`）* ⇒ 三枚 ban 测试全 `t.Fatalf`，
  方向 **fail-closed**。⇒ 存量合法数据被"拒"的唯一一类是**第二枚 decode 本身**，
  代价是"必须落到同包 struct 或另找一把尺"，**不报错的形状不存在**（它是响，不是忍）。

### §3.3 "池里有名字但要命的是别的东西" —— 一发证明池本身无罪

我把 11 枚被滤掉的东西里最像路由的那枚（`shell.run`）从注释搬进代码，看会不会凭空造红：

⇒ 这一发我没跑（`mutate2.py` 只支持两枚种子；跑它要再改台件），**记在 §8"没测什么"里**，
不当已证。我手上等效的证据是两条盘上事实：(i) `index.html` 今天在池里且 `answeredByRealGuard=false`、
包体全绿（§3.1 现量）；(ii) `poolJudgedByRealGuard` 第一句就是 `if !knownComposerMethod(name) { continue }`。
⇒ 判："进池不产生代价"这一条我给**代码事实＋一枚真实驻留样本**，不给穷尽证明。

### §3.4 这一格总裁

**没有误伤，但有一条新的、没人声明过的生产码约束**（§2.4 的 F-R2-1）。
判：**成立（附一条文档级要求）**——(b) 与 (c) 今天的存量代价都是零，
而它们对**将来**的合法形状施加的约束（第二枚 decode 必须落到同包 struct / 新入站类型必须进
`inboundTypeRegistry`）都**响在明处且自带出路**。这不是"放宽断言"，是"把新约束写进文件头"的事，
所以是文档级、不退回。

---

## §4 (c) 与 (c′) 有没有把功劳记混

实现程自己纠过一次（`cc6dc16` 把 M7 从 (c) 的收益里划给 (c′)）。我用**自己的种子**复核这个更正，
并回答派单那一问："(c) 单独有没有任何行为后果？"

| 种子 | 今天谁让它红 | 摘掉 (c) 整味 | 摘掉 (c′) 两半 | 归属（本程判） |
|---|---|---|---|---|
| `M6`（`cmd` 路由键 + 显式 `outcome` 标签） | `TestNoInboundEnvelope` 点名 `bridge.go:122` | **真树全绿**（只剩植物 G） | 仍 3 枚红（ban 那枚在） | **只属于 (c)** |
| `M7`（绑 `method`，但 `Outcome json:",omitempty"`） | `TestNoInboundEnvelope` 点名 `bridge.go:122` "Outcome" | 仍红（旧判据也看得见它，它绑 method） | **真树全绿**（只剩植物 H） | **只属于 (c′)** ⇒ `cc6dc16` 那半句更正**成立** |
| `M8`（同形状放进正主 `ComposerRequest`） | 反射半 `:1218`（**这半先存在就红**）＋ AST 半 `:1237`（新增）＋ 三向 `:1492` 两条 | 4 枚红，不变 | **颜色不变**（反射半一直兜着），变的是"AST 不再 0 findings" | **不是任何一发新增拦截**，是两台仪器不再分家 |
| `M13` 封套半 | `:1237 bridge.go:131` | 转绿 | 仍红（M13 用显式 `json:"outcome"`，不需要 (c′)） | **(c)** |

⇒ **(c) 单独有没有行为后果**：在**今天的真树**上没有（`SEEDS`/`INBOUND` 与旧判据同一枚格，§3.2 现量），
但在**我的台件上有一发唯一目击**：**M6**（以及 M13 的封套半）。
⇒ 所以 (c) **不是镜子**，判**承重**；措辞要精确到这句：
**(c) 今天不拦任何已存在的写法，它拦的是"第二枚 decode"——而第二枚 decode 正是票 33/35 接线那一步必然写下的东西。**
"defense-in-depth 而非 today-caught mutation" 这个说法对 (c) 只对一半：它对 *clean tree* 成立，
对 *M6* 不成立。我按现量给结论，不给那句顺口话。

⇒ **(c′) 的两半不许混着领功**：规矩半（`jsonOr` 回落 Go 字段名）是 M7 的唯一 ban 级目击者；
测试半（三向互印 `TestJSONKeyDerivationAgreesWithEncodingJSON`）**在我全部 6 发种子 + CLEAN 上
都不改变任何一枚的颜色**（`noCp_test` 行：CLEAN=0 红、M6/M7=2、M8=3、M13=3，ban 那枚始终在）。
它改变的是"**摘掉 jsonOr 的回落之后有没有人响**"（`noCp_rule|CLEAN` 三枚红，其中 `TestJSONKeyDerivation`
是那枚**唯一**响的 ban 无关者）。⇒ 判：**承重、但白盒（instrument-integrity）**，
不写"M8 因为它才被拦下"（M8 今天被反射半拦下，与它无关）。

---

## §5 (d)(e) 两格 + 那枚悬空的 §9

### §5.1 (d)："每个断言配一枚可复算命令"——我把三条都跑了

| 文件里的断言 | 它的命令（逐字照抄去跑） | 我量到的 | 判 |
|---|---|---|---|---|
| `ParseComposerRequest` 无生产调用者 | `grep -rn ParseComposerRequest --include=*.go .` | 非测试命中 6 处：`bridge.go:77`(定义) + `run.go:225`、`bridge.go:16/74/106`、`composer_handlers.go:37` **五处全 `//` 注释** | **结论成立**（生产调用者=0），**明细两处对不上**：文件写 `l2_grant_boundary_test.go(6)`，盘上是 **9**；文件写"其余四处全是注释"，盘上是**五处** |
| 主模块收不到 postMessage | `grep -c webview go.mod` | `0`（`go.sum` 也 0；仓内另有 3 枚独立模块 `scripts/spike`、`tools/d22scan`、`tools/mockllm`，与主模块无关） | 成立 |
| 那跳不存在是生产码自己说的 | `sed -n '225,227p' cmd/wisp/run.go` | 逐字复现 `Nothing calls it yet - the WebView2 "event -> ParseComposerRequest" hop does not exist in this tree (tickets 33/35)` | 成立 |
| 头段的 "NOTHING IS WIRED TODAY" | 同上三条 | 三条独立复算全中；另加我本程自己跑的：`PanelBridge` 在非测试生产码里**只有 3 处注释命中**（`internal/agent/approval/gate.go:595`"a future PanelBridge"、`doc.go:2/:9`），**零实现零调用者** | 成立，且比它自己声明的更硬 |

**那两处对不上算不算缺陷**：算，但只算**读数级**。它们都是**明细的枚数**，不改任何一条判语的成立；
且 `l2_grant_boundary_test.go(6)` 那枚错**在它自己那一枚 commit 上就已经错了**
（`git show 598620e:internal/panel/l2_grant_boundary_test.go \| grep -c ParseComposerRequest` = **9**，
不是它写进文件的 6）⇒ 不是"后来漂了"，是**抄错**。记 **F-R2-2**（读数级，见 §7）。

### §5.2 (d) 那句是不是可证伪的，还是装饰

**可证伪。** 让这句话为假的充分条件已经写在同一行里：`grep -rn ParseComposerRequest --include=*.go .`
出现任一**非注释、非测试、非定义**的命中（也就是票 33/35 那跳被接上）⇒ 这句当场失效。
我把它当判据复算过（上表第一行 = 今天的读数是 0）。

⚠ 但有一句要收窄，不然它会被读大：这句话的**主语是 `ParseComposerRequest`**。
理论上"接上另一条不经过 `ParseComposerRequest` 的入站路"（例如新写一枚 `handleHostMessage`）
**不会**让这句话变假，而它确实违反"什么都没接线"的精神。⇒ 今天堵这个形状的不是 (d)，
是 (c)+(c′) 那两枚 decode 引信（§2.4 的 MDEC/MEXT 现量：任何第二枚 decode 都会响）。
所以我判 (d) **成立且可证伪**，并建议补半句"（这句话的主语是 `ParseComposerRequest`；
任何第二枚 decode 由 (c)/(c′) 的引信管）"——**这是建议，不是这一格的退回理由**。

### §5.3 (e) 与 `bridge.go:33`

```
grep -rn "func TestFrontendComposerRequestsMatchTheEnvelope" internal/panel/   -> 1 处，bridge_test.go:131
读全文：方法那一圈是 for _, m := range []string{MethodModeRequest, MethodWorkspaceRequest,
        MethodAttachmentAdd, MethodMessageSend}   —— 恰四枚常量，从不读守卫
```
⇒ **"它只迭代那四枚常量"这半句为真**（那枚测试还查 `requestId`/`source`/第二通道/
`bridge.postMessage` 计数，但那些都不是路由名枚举，读者的结论"第五枚直接写进 switch 的路由它看不见"成立）。
旧名 `TestComposerMethodNamesMatchFrontend` 现量全仓 **1 处命中 = `bridge.go:33` 自己**，定义 0 处
⇒ `:604` 那一处引用**确实被修掉了**，不是换了个错名字。

**`bridge.go:33` 我判：同意本批不碰、且它确实是缺陷。**
理由：改它 = 动生产文件一个字节 = 撞本批自己那条"生产码零字节"契约轴；而 r1/派单明令我等
不扩生产码地界。归属：**它属"把 `event -> ParseComposerRequest` 那一跳接上"的那一程**
（`cmd/wisp/run.go:225` 与 `composer_handlers.go:37` 两处注释都点名 **tickets 33/35**），
不是门钉这一程。最小修法就一句：把 `bridge.go:33` 那句改成点名
`TestFrontendComposerRequestsMatchTheEnvelope`（`bridge_test.go:131`）。
**不要**把它夹进任何一枚门钉 commit——那正是本批靠"零字节"守住的边界。

### §5.4 悬空的 §9 / §7（格级判）

见 §0.4 现量：**三枚自指悬空**（`:8`→§9 numstat、`:73`→§7 反向对照、`:81`→§9 staged 清单），
**两枚外指有效**（`:4`/`:131`/`:201` 指 accept-r1 §7.x，存在）。我另扫了它引用的**行号**：
§5.3 里的 `:1211`/`:1466` 在 `c4b959a` 那版逐枚对得上（我 `git show c4b959a:… \| sed -n '1211p;1466p'` 现量），
但**在 HEAD 上漂到 `:1237`/`:1492`**——漂移是它自己**下一枚 commit `598620e`**（+36/−10 行）造的。
⇒ 这不是抄错，是"**引用自己没写完的东西，又自己把它挤走了**"。

**格级判：退回（文件级，不推翻码）。** 三枚悬空自指正好落在"码写完、读数没跑完"那个断点上，
形状与本仓已登记那条（"文档预先引用尚未产出的读数 = 假绿前身"）**同类**；
而它已经被 `4deada9`/`0fe7629` 之前的两次同类登记钉过，所以处置固定：
**文件保留、不删、不许当"这批已验收"凭据引用**；要恢复它 `:8`/`:73`/`:81` 的承诺，
要么补出 §7/§9（本程的 §2/§6 就是同一批读数的独立版本），要么就地在那三行后追加
"（本程未产出，见 `panel-l2-grant-nail-accept-r2.md` §2/§6）"。**只能追加，不许改写已提交的四行。**

---

## §6 闸门 + 名册 + 契约轴（本程自己量，锚点旁注）

### §6.1 五枚被验 commit 的 name-only（逐枚现跑 `git show --name-only`）

| commit | 带的.path | 枚数 | 判 |
|---|---|---|---|
| `9d85789` (a) | `docs/evidence/s1/panel-l2-grant-nail-fix-r2.md`、`internal/panel/l2_grant_boundary_test.go` | 2 | 只带自己两枚 |
| `67c1d20` (b) | 同上 | 2 | 只带自己两枚 |
| `2de984c` (c) 码 | **只有** `internal/panel/l2_grant_boundary_test.go` | 1 | **只带一枚**——码与证据分两枚落，证据由 `ad42e6f` 后补；这不是缺陷，但派单那句"只带自己**两枚**路径"对它字面不成立，我按"没带别家"判成立 |
| `c4b959a` (c′) | 证据 + 测试件 | 2 | 只带自己两枚 |
| `598620e` (d)(e) | 证据 + 测试件 | 2 | 只带自己两枚 |
| 证据四枚 `3c0b630 ad42e6f cc6dc16 fa46fe3` | 各自只有 `panel-l2-grant-nail-fix-r2.md` | 1 | 只带自己一枚 |

`9d85789…598620e` 逐枚 numstat（`-- internal/panel/`）：`142/20`、`202/2`、`386/27`、`315/5`、`36/10`——
**分母里只出现那一枚 `_test.go`**。

### §6.2 契约轴（零容忍面，每条都带正控）

```
git diff --numstat d88c356..HEAD -- internal/panel/
   1068   49   internal/panel/l2_grant_boundary_test.go        ← 一枚路径，生产码零字节
git diff --numstat d88c356..HEAD -- internal/observe/thresholds.go internal/risk/rules_gateway.go \
     tools/d22scan/allowlist.txt internal/risk/ docs/PLAN.md docs/specs/ '*.sse'
   （空）                                                     ← 零命中
正控（同一把尺、同一个区间）：
git diff --numstat d88c356..HEAD -- internal/panel/l2_grant_boundary_test.go
   1068   49   …                                              ← 尺是活的
```
- 被点的三面先证明它们**存在**才敢说"零命中"：`git ls-files | grep -i 'thresholds.go|allowlist.txt|rules_gateway.go'`
  → `internal/observe/thresholds.go`、`internal/risk/rules_gateway.go`、`tools/d22scan/allowlist.txt`；
  golden `*.sse` 计数 = **52 枚**，未动。
- **断言没有少**：`git show d88c356:…` vs `HEAD:…` 的
  `t.Fatalf 21→31`、`t.Errorf 21→38`、`t.Fatal( 3→5`，合计 **45 → 74**；
  `t.Skip` **0 → 0**（本包唯一一枚 `t.Skipf` 在 `attachments_test.go:191`，先存在）；
  `func Test` **5 → 6**（diff 里只有一行 `>`、**零行 `<`**）。
- **删掉的 49 行**我逐枚看过（原文在 `out/`）：6 枚 `t.Fatalf` + 1 枚 `t.Errorf` 的**消息改写**、
  `jsonOr`/`inboundEnvelopes` 的判据替换、`routeNamesInFunc`→`routeLabelsInFunc` 的更名——
  没有一枚断言被**去掉**，只有同一枚断言换了措辞或更强的判据。
- **没有任何豁免表**：`git diff d88c356..HEAD -- internal/panel/ | grep '^+' | grep -i
  't.Skip|exempt|allowlist|//nolint|except'` = **空**。
- **helper 全是本批为钉子自造**：新增 26 枚 `func`（`collectPackageVars`/`collectDecodeSites`/
  `resolveDestination`/`routeLabelsInFunc`/`routeShapedName`/`routeNamePool`/`poolJudgedByRealGuard`/
  `answeredOutsidePool`/`inboundTypeRegistry`/`astTopKeys`/`reflectionTopKeys`/`decoderKnowsKey`… ），
  旧名 `routeNamesInFunc` 全包引用 **0 处**（彻底退休、没留一条还能被叫到的弱版）；
  numstat 只有那一枚 `_test.go` ⇒ **不可能是"改弱了一枚先存在的 helper"**，因为别的文件一字节未动。

### §6.3 三门 + 名册（终态复跑；锚点写在数字旁边）

| 量 | 锚点 | 读数 |
|---|---|---|
| `gofmt -l internal/panel/` 与 `gofmt -l .` | `0b95e9e` | 空，rc=0 |
| `go vet ./internal/panel/` | `0b95e9e` | 空，rc=0 |
| `go test ./internal/panel/ -count=2 -v` | `0b95e9e` | `RUN=196 TOPPASS=102 SUBPASS=92 FAIL=2 SKIP=0 ^panic:=0 distinct=98` |
| 同上的"改前对照"（同一棵副本，钉换回 `5b618d14`） | 副本 | `RUN=180 TOPPASS=100 FAIL=2 SKIP=0 panic=0 distinct=90` |
| `sh scripts/d22scan.sh` | `0b95e9e` | **rc=0**；`bans #1-5 internal/=203 cmd/=22`、`#6 frontend/=43`、`#7 internal/tools/=18`、`#8 design/=32 frontend/=43 internal/=407 cmd/=39` |

- `FAIL=2` = `TestC21DesignTokensFourWayAgree` ×2（`out/clean-count2.txt` 第 249／510 行）。
  **它必须是红的**：本程未修未跳未放宽，也没把它算进任何一格收益。⛔ 若下一位读到它是绿的，那是放水。
- ⚠ 关于 `d22scan.sh` 的两条口径（照抄本仓已定案的读法，不是本程新造）：
  bans #1–5 **不扫 `_test.go`**，所以本批这枚文件的墙钟/裸 goroutine 形状**不在它们分母里**；
  只有 ban #8 含测试文件与注释；墙钟那条只认 `.Sub(time.Now())` 一种拼写。
  `frontend/=43` 会长（第三家在写），本件只当"我这一发的分母"。
- ⚠ 本程**没有**跑 `go test ./...`（派单明令，另两路 agent 在飞）。
