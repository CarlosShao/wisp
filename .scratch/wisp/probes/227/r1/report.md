# 票 227 落地腿 227-r1 report（写代理，2026-10-18 第二波）

起手锚＝我自己 `git log -1` 现认 `feac8f77`（336-r1 那笔）。★本腿期间 HEAD 漂到 `17e67cca`（322-r2 收尾），
已现验 `internal/config` 两枚目标件 blob 在两号之间未变（`writeguard.go`=`eab383bf`、
`writeguard_226_test.go`=`29dbd968`）⇒ 基线树有效。票 227 现量尺＝6 未勾／0 已勾（`grep -c '^[[:space:]]*- \[ \]'`），
**本腿一枚框没翻**（⛔ 实现者翻自己的框）。

## 提交名册（逐枚 `git show --name-only`＝写面自证，AC#4）

| 提交 | 主题 | name-only（应只含被允许写面） |
|---|---|---|
| `b72d7391` | AC#1③＋AC#2甲·三句日志字面串说实话 | `internal/config/writeguard.go`（`8/3`） |
| `3b8c114e` | AC#1＋AC#2·常驻用例 | `internal/config/writeguard_handedit_227_test.go`（`227/0`） |
| `dd0e489d` | 定稿用例件（编排者收回 ≤160 帽对 `*_test.go` 的适用） | 同上（`52/84`） |
| `6da21d4a` | report（本件）＋八枚 logs 入库 | `.scratch/wisp/probes/227/r1/**`（9 files＝本件＋c0…c7） |
| `32a70e2f` | 票尾一条 bullet | `.scratch/wisp/issues/227-*.md`（`2/0`＝只追加） |
| `9f5b0022` | 名册尺更正＋交件侧终量门禁 | 本件＋`logs/c8…c11`（5 files） |
| `3cda0aef` | AC#4 终尺 c12 入库（＝c8 过期，具名更正） | 本件＋`logs/c12`（2 files） |

★**还有一枚＝本行所在这一笔，号⛔ 在这里预先写死**（写这行时它还不存在＝尚未存在的号⛔ 当凭据）；
收件侧现跑尺＝`git log --format="%h %s" feac8f77..HEAD` 只留 subject 以「票 227 落地腿 227-r1·」开头的枚。
名册尺终量＝`logs/c12-roster-six-commits.txt`（它自身也⛔ 覆盖把自己写进去的那一笔＝同一枚自指道理，件头写明）：
逐笔 `internal/config/**` ∪ `.scratch/wisp/probes/227/r1/**` ∪ 票 227 **之外＝零字节**，合计 0，两枚正控在件内。
- ⚠ **票尾那条 bullet 里的"四枚提交"已过期**＝它写于本腿第 4 枚（`6da21d4a`）落盘之后，同一枚自指道理
  （一条 bullet 装不下后续笔）。派单只许票尾**一条**追加 bullet ⇒ 本腿⛔ 为补枚数去追加第二条（那＝扩写面）。
  **终名册以本件这张表＋`logs/c12`＋收件侧现跑那把 subject 尺为准**；票尾那条只当"本腿来过的路标"，不当名册。

## AC#3 交付主体：四种手改形状矩阵（分母＝`226-v1` 的 `b2` P1–P6，⛔ 重造）

三态＝〔有用例〕常驻 `*_test.go`／〔只有台件〕只活在 `226-v1` 一次性台件／〔什么都没有〕零证。
"本腿后"列＝本腿动过的格；⛔ 本腿没动的格一律**照实报未覆盖**（AC#3 要的就是这张表的诚实）。

| # | 手改形状 | 226-v1 台件读数 | 本腿前 | 本腿后 | 载体（`file:line`） |
|---|---|---|---|---|---|
| P0 | **改值**（`size = 56`→`60`） | 未测（票 226 自己量的那一支） | 有用例 | 有用例（未动） | `writeguard_226_test.go:63`/`:102`/`:117`/`:187`/`:269` |
| P1 | **删一枚非默认值标量行**（`size = 60`） | 补回＝默认 56、内存仍 60、注释丢 | **什么都没有** | **有用例** | `writeguard_handedit_227_test.go:23` |
| P2 | 删一枚**值本就是默认**的标量行（`size = 56`） | 补回＝56，走 Info 支 | 只有台件 | ⚠ 只有台件（本腿未补，见"欠账"） | `probes/226/v1/b2-probe-deleted-shapes.txt:12-18` |
| P3 | **删整节**（`[ball]` 连表头） | 表头＋`size` 都补回、注释丢 | 什么都没有 | ⚠ 仍未补（见"欠账"） | `b2…txt:20-27` |
| P4 | **删 map 条目**（`[llm.providers.deepseek]`） | **不补回**（`providers=0`） | 只有台件 | **有用例**（反向断言） | `writeguard_handedit_227_test.go:70` |
| P5 | 删 `schema_version` 行（其余与内存分叉⇒走**写**支） | Info `wrote merged change`、补回 `= 2` | 只有台件 | ⚠ 只有台件（未补） | `b2…txt:38-44` |
| P5′ | 删 `schema_version` 行（文件**已带**该键⇒走 **nothing written** 支） | **Y2：日志说 `nothing written`、bytes 变、mtime 变、注释丢、落 `.bak-1`** | 只有台件 | **有用例** | `writeguard_handedit_227_test.go:106` ⇄ `probes/226/v1/b4-probe-migration-inside-read.txt:12-19` |
| P6 | 删 provider `base_url` 行（preset 补） | 补回 `base_url = ''`（preset 只在内存补、文件留空） | 什么都没有 | ⚠ 仍未补 | `b2…txt:46-52` |

★读法：**票 226 AC#1/AC#2 的射程边界就在这张表上**——它今天只在 P0 一行为真；本腿把 P1／P4／P5′ 三行钉成常驻，
P2／P3／P5／P6 四行**如实报未覆盖**（AC#4 不许顺手扩大射程＝本腿只补票面 AC#1/AC#2 点名的三行）。

## AC#1 三断（票面逐字）＋成对正控

用例 `TestHandDeletedScalarLineComesBackAtSchemaDefault`（`writeguard_handedit_227_test.go:23`）：
- ①那一行**回来了**，值＝**schema 默认 56**，⛔ 内存那枚 60（反向断：内存仍 60，受守卫写没把回落塞进内存）；
- ②`kept_in_file_not_in_memory=[ball.size]` **有**把这枚算进去；
- ③那行注释**没了**，且**日志那句人话**在场（`a key line the operator deleted by hand comes back at its schema
  default value`，`b72d7391` 落 `writeguard.go:174-175` 与 `:188-189` 两处）＝票面"不许只在源码注释里"。
- 成对正控＝**未改动树上同一发绿**：本腿五枚在交件树全 PASS（`logs/c1-final-delivery.txt`），⛔ 恒红。

## AC#2 走甲（编排者裁 (A)：授权扩到 `:152`，边界四条全守）

1. **先证指名那枚用例真存在**：`git grep -n 'func TestGuardedWriteDoesNotRewriteAFileThatAlreadySaysIt'`
   → `internal/config/writeguard_226_test.go:393`（票面行号零漂移）；正控＝同尺＋`XYZ` 后缀 → rc=1 零命中（尺活）。
2. **种一发让它红**（仓外导出树、内容锚、串行）：fixture 里删掉 `schema_version = 2\n` ⇒
   ```
   writeguard_226_test.go:405: the file already carried the entry and was rewritten anyway:
   writeguard_226_test.go:408: the file's mtime moved on a write that wrote nothing
   --- FAIL: TestGuardedWriteDoesNotRewriteAFileThatAlreadySaysIt (0.03s)   exit=1
   ```
   逐字见 `logs/c6-mu3-393-schema-red.txt`；还原 sha256=`49960DA6…EF9B8`＝pristine。
   ⇒ **它今天恒绿的原因被暴露**＝fixture 的 schema 恰为当前版、迁移没机会进场。⛔ 没把断言改成"接受新字节"。
3. **文案说实话（三句字面串，⛔ 任何分支／条件／签名／等级／`"key", ownedKey` 实参／`diverged` 算法）**：
   `:152` 由 `already carries this key; nothing written` 收窄为带条件句（bytes/mtime 不动**只在 schema_version
   已是当前版时**成立；落后版文件早已被这次重读里的迁移重排成 canonical、注释与自选键序不保）；
   两句 `re-laid…` 各追加 AC#1③ 那句删行人话。
4. **动前三把读者尺**（`git grep -n`，含 `scripts` `tools` `.github` `*_test.go`）：
   `the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved` →
   判据面只有 `writeguard.go:183`（它自己）＋`docs/evidence/s1/226-…-r1.md:121`（叙述引用，非计数）；
   `re-laid`／`hand-chosen key order` 在 `scripts tools .github` → **rc=1 零命中**；
   `already carries this key` 与 `nothing written` 在 `*.go` → 只有 `writeguard.go:96`（注释）＋`:152`（它自己）；
   全仓唯一被当判据的是**属性名** `kept_in_file_not_in_memory`（`writeguard_226_test.go:132`/`:156`）＝一字未碰。
   ⇒ 零枚用例按这三句字面串计数／断言 ⇒ 未触发"改它＝改判据"那道停手闸。
5. **双向钉**：收窄后的承诺还配了正控 `TestNothingWrittenLeavesFileUntouchedWhenSchemaIsCurrent`
   （schema 当前版＝字节/mtime/注释全保、真不动文件）⇒ 这枚边界⛔ 只咬"变严"那一形。
   代价（具名）：日志变长、那句不再是短句；行为一字未变。
6. **⛔ 乙**（真让 `nothing written` 那支不动文件）＝要动 `readConfigFile` 的迁移-on-读管线顺序
   （最小缺口具名：`internal/config/loader.go:126` → `internal/config/migrate.go:93`＋`:96`），
   量级属票面「落点候选」乙，**本腿没做**（派单：走乙＝停手上报）。

## AC#4 不改判级、不扩射程

- 动过的产码＝**只有字符串字面**（`git show b72d7391` 三处 hunk，条件表达式/函数签名/`statOwnWrite` 调用位置/
  `diverged` 算法/`allowed_dirs` 语义/`readConfigFile` 迁移行为**全部未动**）。
- 新增导出名＝0；新增 `panel.*` 方法名＝0；新增 JSON 键＝0；⛔ 动 SLO 阈值／golden／`thresholds.go`／
  `allowlist.txt`／三枚冻结件／`go.mod`／`go.sum`／`.gitignore`／`ci.yml`／D1–D47／C1–C32／R1–R9。
- 名册尺逐笔＝`logs/c7-roster-per-commit.txt`。

## 颜色与门禁（交件侧终量）

台面＝仓外导出树（`git archive` 八件套＋`docs`；`.go` 数＝762 基线⇄763 交件，⛔ "0"是读数）；
`-count=1 -timeout 30m`；四数⛔ 混用⛔ 相减。

| 那把尺 | RUN | PASS_TOP | FAIL_TOP | SKIP |
|---|---|---|---|---|
| 基线（锚 `feac8f77` 原样） | 185 | 104 | 0 | 0 |
| 交件（原码＋我的字面串＋我的用例） | 190 | 109 | 0 | 0 |
| 未改产码 ⇄ 只换我的用例（成对） | 190 | 107 | **2** | 0 |
| MU-1 `diverged :=` 掏空（报告侧） | 190 | 105 | **4** | 0 |
| MU-2 合并基准改内存快照（写侧） | 190 | 100 | **9** | 0 |

- 红名册作差：基线⇄交件＝**零枚新增红**（交件恰＝基线＋本腿 5 枚全绿）；未改产码那一发只红本腿两枚**文案**断言
  （`writeguard_handedit_227_test.go:64`／`:129`）＝证"用例盯的是文案、⛔ 恒红"。
- ★**"未修码响不响"逐件成对**（＝派单必答 (iv)，尺＝未改产码⇄只换我的用例，件 `logs/c2-final-unmodcode.txt`）：
  | 本腿用例 | 钉的事实 | 未修码响不响 | 成对读数出处 |
  |---|---|---|---|
  | `TestHandDeletedScalarLineComesBackAtSchemaDefault:23` | ①补回＝默认 56／②`kept_…=[ball.size]`／③注释丢＋日志说实话 | ①②③-注释＝**不响**（既有码本就如此，属被钉的既有代价）；③-日志＝**响**（红句 `:64`） | c2 ⇄ c1 |
  | `TestHandDeletedMapEntryIsNotRestored:70` | P4 删 map 条目不补回 | **不响**（既有行为，本枚把它从台件升级成常驻） | c2 ⇄ c1 |
  | `TestNothingWrittenPromiseOnlyHoldsWhenSchemaIsCurrent:106` | 边界＋收窄句 | 边界半边**不响**（迁移确实重排＝既有码）；收窄句**响**（红句 `:129`） | c2 ⇄ c1 |
  | `TestNothingWrittenLeavesFileUntouchedWhenSchemaIsCurrent:150` | 正控（schema 当前版真不动文件） | **不响**（＝成对正控，证上两枚⛔ 恒红） | c2 ⇄ c1 |
  | `TestDeletedLineRestorationIsNotClaimedAsAWrittenKey:181` | AC#4 删行形不认领别键 | **不响**（既有形状） | c2 ⇄ c1 |
  ⇒ 未修码合计只红 **2** 枚，且两枚红句逐字都指向**文案**（⛔ 指向字节或行为）＝AC#2 走甲
  "只补仪器＋让文案说实话、⛔ 改行为"那句的机读证据。★AC#2 的"改码才响"那一维另有一发＝
  `logs/c6-mu3-393-schema-red.txt`：未改产码⇄**把 `:393` 那枚的 fixture 塞进 schema 落后形** ⇒
  既有码红（`:405`＋`:408`）＝证那枚用例今天对这一形**失明**，不是它测不了。
- 成对突变落点不同（报告侧⇄写侧）⇄ 还原 sha256＝pristine `EFE12661…D47BD`（跑前后各验一次）。
- 门禁四数（只在共享工作树，★跑前 `unset GOTMPDIR TMPDIR`，逐枚 `rc=`）：`gofumpt -l <本腿两枚>`＝空输出 rc=0
  ／`sh scripts/d22scan.sh` rc=0 clean／`sh scripts/check-path-length-budget.sh` rc=0 VERDICT GREEN／
  `go vet ./internal/config/` rc=0。★本机 6 枚 CRLF 存量脏＝**别腿件，本腿一枚没重排**（具名登记，⛔ 动）。
- ⛔ 突变那两发未碰任何跟踪文件；`c6` 那发在仓外树上改 fixture 后**当场还原**。
- 交件侧**终量**门禁（本腿最后一枚件定稿后复跑）＝`logs/c9-gates-final.txt`（d22scan rc=0 clean＋gofumpt 空输出 rc=0）
  ／`logs/c10-pathlen-final.txt`（分母 15570 tracked、over-budget 60 全在册、not in roster 0、rc=0）
  ／`logs/c11-vet-final.txt`（rc=0）。⚠ 读 `c9` 别把里面的 `1 finding(s)` 当红：那六行 rc=1／rc=2 全部来自
  **d22scan 自带的 self-test**（故意 seeded violation，件内正读＝`--- PASS: TestBuiltBinaryGoesRedEndToEnd`），
  真判据是末行逐字 `d22scan: clean - no D22 ban violations`＋`d22scan rc=0`。
- ★名册尺更正：`logs/c8-roster-five-commits.txt` 头四行写明——第一发用 `git log --grep='227-r1'` 量到的名册
  **不干净**（把别腿正文里提到本腿名的笔也吃进来＝336-r1／编排者收工笔），那把尺**作废**（⛔ 删，只建不删），
  改用严格 subject-prefix 尺 ⇒ 逐枚「允许写面之外字节」＝**0**，两枚正控证这把尺⛔ 空转。
  ⚠ **c8 自己又过期了一次**：它量到 5 枚，本腿随后又多两笔（`9f5b0022` 与它自己）⇒
  **终尺＝`logs/c12-roster-six-commits.txt`**（6 枚全量、逐枚之外字节合计＝**0**；正控一＝同尺打票 336 那枚前缀
  读到非零，证⛔ 空转；正控二＝把 `6da21d4a` 只按 `internal/config/**` 量到 9 枚之外字节，证尺能开火）。
  ⛔ 任何一枚名册件自指覆盖"把自己写进去的那一笔"（写件时它还不存在）⇒ **收件侧请现跑 c12 头两行那把尺**。
- 票面框尺（交件侧现量，`c12` 末行）＝未勾 **6**／已勾 **0**＝与起手读数逐字相同 ⇒ **本腿一枚没翻**。

## 欠账（具名交回编排者，本腿⛔ 自行扩面）

1. **AC#5／AC#6**＝要 `cmd/wisp` 整包（车道归写腿，现 `338-r1`）⇒ **具名交回**，本腿零跑零判。
   `cmd/wisp/always_write_no_clobber_226_test.go:43-44` 那句 `nothing polls` 自 `248095d1` 起为假（票面 AC#6 已记）
   ＝**归你**（派单公共块：`cmd/wisp/**` 是共享编译面）。
2. **`245-v1` 那枚挂在票面的欠账**（"`cmd/wisp` 默认层读到 1 红 2 绿、用例名丢了"）＝**归你排程**；
   本腿⛔ 跑那包，也⛔ 与 `322-r2`／`338-r1` 同跑。
3. **矩阵里四行仍未钉**（P2 值本就是默认的删行／P3 删整节／P5 删 `schema_version` 且走写支／P6 删 `base_url`）
   ＝票面 AC#3 要求的"逐格标三态"已如实标〔只有台件〕／〔什么都没有〕；补不补＝**另判**，
   本腿按 AC#4 只补票面 AC#1/AC#2 点名的三行。
4. **`docs/evidence/s1/226-config-write-no-clobber-r1.md` §1/§5 那句"不白重排、不白丢注释"**仍说过头
   ＝票面「谁会说谎」那一行点名要更正的那句，且它落在**裁决件**（本腿写面⛔ `docs/evidence/**`）⇒ 归非实现者。
