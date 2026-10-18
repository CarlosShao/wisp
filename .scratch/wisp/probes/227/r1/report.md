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
| 本笔 | report＋logs 入库 | `.scratch/wisp/probes/227/r1/**` |
| 末笔 | 票尾一条 bullet | `.scratch/wisp/issues/227-*.md` |

`internal/config/**` 与 `.scratch/wisp/probes/227/r1/**` 与票 227 之外＝**零字节**（尺见 logs/`c7-roster-per-commit.txt`）。

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
- 成对突变落点不同（报告侧⇄写侧）⇄ 还原 sha256＝pristine `EFE12661…D47BD`（跑前后各验一次）。
- 门禁四数（只在共享工作树，★跑前 `unset GOTMPDIR TMPDIR`，逐枚 `rc=`）：`gofumpt -l <本腿两枚>`＝空输出 rc=0
  ／`sh scripts/d22scan.sh` rc=0 clean／`sh scripts/check-path-length-budget.sh` rc=0 VERDICT GREEN／
  `go vet ./internal/config/` rc=0。★本机 6 枚 CRLF 存量脏＝**别腿件，本腿一枚没重排**（具名登记，⛔ 动）。
- ⛔ 突变那两发未碰任何跟踪文件；`c6` 那发在仓外树上改 fixture 后**当场还原**。

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
