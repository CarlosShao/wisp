# 票 223 写腿 `223-r2` —— AC#5 间歇红改确定性 ＋ AC#4 两句归正

- 工单：`.scratch/wisp/issues/223-checkandreload-has-zero-production-callers-hot-reload-never-runs.md`
- 前案裁决表（必读，本腿的射程由它划出）：`docs/evidence/s1/223-hot-reload-wiring-v1.md`（56,522 字节／265 行）
- 本腿起手锚点：`git rev-parse --short HEAD` = **`c774f8da`**，分支 `dev`
- 本腿性质：**写腿**，只做派单的两件事（AC#5 轮询化／AC#4 句子归正＋常驻用例），⛔ 不扩射程
- 台件目录：`.scratch/wisp/probes/223/r2/`（原始输出全量落盘，不接 `| head`／`| tail`）
- 入库 pathspec（只此四枚）：`cmd/wisp/**`、`internal/config/**`、本表、`.scratch/wisp/probes/223/r2/`、票 223 工单面（只追加）
- ⛔ 零读零写零转述：`frontend/**`、`design/**`；⛔ 不碰：`PLAN.md`、`docs/specs/**`、`thresholds.go`、golden、`allowlist.txt`、三枚冻结件
- ⛔ 票 223 不翻 `-done`、不勾任何 AC 框（翻勾由编排者按非实现者表＋自己复跑裁）

## 起手读数（本腿现跑）

- 起手 `date`：**Tue Sep 29 16:27:20 CST 2026**；骨架提交＝`28475620`。
- `git diff --cached --name-only` 起手＝**空**（行数 0，索引干净，没有别人的活在飞）。
- 在飞检查：`ps -W | grep -i -E "go\.exe|go-test|test\.exe"` 起手＝**空**。
- 别人地界的脏改动（起手 `git status --porcelain` 现读到名字、本腿一律不碰其内容）：`.gitignore`、`design/**`、`.scratch/wisp/probes/152/**`、`.scratch/wisp/probes/161/r6/logs/**`、`docs/evidence/s1/152-*.md`、`docs/reports/HANDOVER.md`、`docs/reports/pending-and-issues.md`、若干未跟踪新票面（如 `230-*.md`）。

## 任务一：`TestTicket223RestartTierSaysItWillNotApply` 间歇红 → 有界轮询

### 1.1 复现（先确认红的是读时序，不是产品句子缺失）

尺（逐字）：`export PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"`，再 `go test ./cmd/wisp -count=1 -run 'TestTicket223RestartTierSaysItWillNotApply'` 连跑 **25 发**（台件 `.scratch/wisp/probes/223/r2/shots_223r2.sh`，原始件 `pre-fix-repro.txt`，不接 `| head`／`| tail`）。

- **读数＝25 发 1 红（第 3 发，`=== shot 3/25 EXIT=1`）**。与历史命中率同向：v1 腿 1/5、编排者 1/15、本腿 1/25（合计 3/45）。**不是必红、是间歇红**，所以判据按"20 连全绿"定。
- 红句逐字（`config_reload_223_test.go:453`）：

  > `the operator is not told the edit will not land this run; stdout:`

  后面贴出的 stdout 只有"答复监听已接入……"＋"配置热加载已接管……"两行通用文案，**没有重启档那句**。
- **归因复认（本腿自己现读产码，不转述）**：`cmd/wisp/config_reload.go` 的 `reportRestartPending` 先写两行审计（`:282`、`:284`，stderr）、**后**写 stdout 句子（`:288`）；用例等的是 stderr（`:438`、`:442` 两次 `awaitAudit`），接着在 `:443` 对 stdout **只读一次**。审计到位、句子未到的窗口是真实存在的 ⇒ **红的是读时序，不是产品行为**。同发里两句 restart 审计都已落盘（原始件可见），产品句子本身没缺。
- ⛔ 按派单：**没有改产品文案迁就测试**。

### 1.2 修法

改动只有一处形状（`cmd/wisp/config_reload_223_test.go`）：

1. 新增 helper `awaitStdout(t, needle)`（紧随 `awaitAudit` 之后）：有界轮询 `r.h.out.String()`，每 20ms 一发，期限 `time.NewTimer(reloadCaseBudget)`（40s）。**形状逐字仿 `awaitAudit`**——用的是 monotonic 的 `time.Timer`/`time.NewTicker`，⛔ 无 `.Sub(time.Now())`、⛔ 无"同一行 `time.Now().Unix*()`＋timeout 词"（d22scan ban #4 两形都不沾）；⛔ 无 `go func(`（ban #1）；轮询与既有 `awaitAudit`/`awaitCard`/`awaitLive` 同构，那些已在今日的 d22scan 读数里证明干净。
2. `TestTicket223RestartTierSaysItWillNotApply` 里把 `out := r.h.out.String()` 换成 `out := r.awaitStdout(t, "本次运行不会生效")`——循环读到出现该句或期限到为止，`t.Fatalf` 兜底。
3. ⛔ **一条断言都没放宽**：`:452` 的 `本次运行不会生效` 断言、`:444` 的三针循环、`:449` 的"重启档不得中途生效"、`:457` 的"不得冒充立即档"、`:460` 的"零卡"全部逐字保留；helper 是"等"，不是"或超时算过"——句子真不到，用例照样红（在 awaitStdout 处红，且 1.4 正控证明了这一点）。

### 1.3 二十连跑判据（同一把尺）

尺（逐字，与 1.1 同一条、台件同一个脚本）：`export PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"`，再 `go test ./cmd/wisp -count=1 -run 'TestTicket223RestartTierSaysItWillNotApply'` **20 发**，原始件 `.scratch/wisp/probes/223/r2/flake-20.txt`（未接 `| head`／`| tail`；跑前 `ps -W` 现查＝空，无别的 `go test` 在飞；shot 1 起手 16:38:55、shot 20 起手 16:40:15、约 16:40:3x 收尾）。

**逐发 exit 码：20/20 全 0**——

| 发 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 10 | 11 | 12 | 13 | 14 | 15 | 16 | 17 | 18 | 19 | 20 |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| EXIT | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 | 0 |

`TOTAL=20 FAILS=0`；每发各有一行 `ok github.com/CarlosShao/wisp/cmd/wisp`（耗时 1.957s～2.285s 逐发在案；尺不带 `-v`，逐名 PASS 行不在这把尺的输出里——20 枚 `ok` 行在原始件内逐发可数）。对照改前 25 发 1 红（1.1）与历史 2/20：判据"同一把尺连跑 20 发全绿"**满足**。⛔ 未出现票 227 射程的 `always_write_no_clobber_226_test.go`——尺只 `-run` 本用例，且整包单发里它也没有红（见三把门节 `final-packages.txt`）。

### 1.4 正控（拿掉那句产品文案必须红）

突变形状＝**改产码**（比改期望串更硬：证的是整条链"产品句子→轮询→断言"）：把 `cmd/wisp/config_reload.go` 里 `reportRestartPending` 的 stdout `Fprintf` 整段注释掉（两行审计 `:282`/`:284` 保留），跑一发 `go test ./cmd/wisp -count=1 -v -run 'TestTicket223RestartTierSaysItWillNotApply'`，原始件 `.scratch/wisp/probes/223/r2/positive-control-m1.txt`（16:41:57→16:42:43）。

- **读数＝红，`EXIT=1`，`--- FAIL: TestTicket223RestartTierSaysItWillNotApply (42.18s)`**（40s 轮询期限＋运行时间＝形状正确：不是秒过，是等到了期限）。
- 红句逐字：

  > `config_reload_223_test.go:481: stdout never carried "本次运行不会生效" within 40s; full stdout:`

  贴出的 stdout 里只有"答复监听已接入……"＋"配置热加载已接管……"两行通用文案；同一发的 stderr 里 `state=restart-pending`、`RESTART-PENDING detail=`、`state=applied` 三行审计都在——**正是"轮询非恒真"的直接证明：句子被拿掉时用例必红，句子在场时（1.3）用例必绿。**
- 还原（⛔ 用 `git cat-file`，未碰 `checkout`/`restore`/`reset`/`stash`/`clean`）：`git cat-file blob HEAD:cmd/wisp/config_reload.go > cmd/wisp/config_reload.go`。
- **certutil SHA256 两哈希逐字**（突变前基线＝还原后，`diff` 输出空＋`git status --porcelain -- cmd/wisp/config_reload.go` 空）：

  `542f705ab7e3c7501ecd479c4096ffdda8a2192858ee7ce3557248a0d34f4670` ＝ `542f705ab7e3c7501ecd479c4096ffdda8a2192858ee7ce3557248a0d34f4670`

  （原件落盘：`hash-config_reload-baseline.txt`／`hash-config_reload-restored.txt`）

## 任务二：AC#4 两句会说反的话归位

### 2.1 现读归因（loader.go 的＋72/−6 与句子分支）

现读三处（本腿逐行读，不是转述）：

- `internal/config/loader.go` 的 `readConfigFile`：`peekSchemaVersion` 只在 `ver == 0 && peekErr != nil`（**连版本都读不出**）时报 `config.toml parse`；`peekErr != nil && ver != 0` 时——版本读得出、正文解析不开——**这一支整段漏掉**，文件直接掉进 `decodeStrict`。
- `internal/config/parse.go` 的 `formatDecodeError`：`toml.DecodeError` 分支产 `"config.toml: line N, col M: expected character ="`。
- `cmd/wisp/config_reload.go` 的 `describeReloadFailure`：`HasPrefix(d, "config.toml:")` ⇒ `cause=invalid` ⇒ 那句中文逐字是「config.toml **语法没问题**，但内容被校验拒绝（值不合法或引用解不开）」。**而这批文件的毛病正是语法错** ⇒ 两句归反。
- v1 腿 overlay 17 发读数（`.scratch/wisp/probes/223/v1/ac4-probe.txt`，本腿直接引用、未重造）：命中这一错形的正是 **C1（注释行以 `[` 开头）／E1（CRLF）／G1（`schema_version=2` 无空格）／L1（缩进版本行）** 四发，`declaredSchemaVersion -> value=2 ok=true`、`peekErrNil=false`，生产句子全是 `cause=invalid`。另 **J1（声明 99＋正文坏）** 同落 invalid，v1 判"半说反"（newer-build 那半句诚实、"语法没问题"那半句假）。**A2（声明 1＋坏表头）今天读回 `cause=migration`——那是迁移管线该拿的，不许动**（`migrate_test.go:123` 钉的就是它）。
- 既有两支不需碰的原因现读：F1（BOM）/I1/K1/D1/D2/B2 走 branch 1（读不出版本），M1/H1 走原有正确路径 ⇒ **改一处路由即可，四形各归其句不需要四各的补丁**。

### 2.2 修法

改在 **`internal/config/loader.go` 的 `readConfigFile`**（不是在 `describeReloadFailure` 里猜错误文本——那会逼句子函数重新解析 go-toml 的行号，等于把路由错误又搬一层）：在 branch 1 之后新增一支——

```go
if peekErr != nil && ver >= SchemaVersionCurrent {
    return nil, observe.Wrap(observe.ClassConfig, peekErr, "config.toml parse")
}
```

三条形由此全部成立：

1. **读得出声明版本（当前版或更新）而正文语法坏 ⇒ 语法错那句**：新支把 peekErr 包成 `config.toml parse`，`describeReloadFailure` 现有分支自动归 `cause=syntax`，句子逐字「config.toml 读到了但解析不了：这一行不是合法 TOML 语法（不是权限、不是缺失）……」。**零新句子、零阈值、零 golden 改动。**
2. **读不出声明版本 ⇒ `config.toml parse`**：branch 1 原样保留（新支在它之后，抢不走）。
3. **只有版本确实低于当前版才交给迁移管线**：`ver < SchemaVersionCurrent` 不进新支（条件是 `>=`），照旧落到 `applyMigrations`——`migrate_test.go:123` 的形状、错误文本（`cannot migrate from schema version 1: ... the file was left untouched`）与断言**一字未动**，其绿由整包 `internal/config` 复跑证明（读数见三把门一节）。
- 配套（注释归真，零行为改动）：`describeReloadFailure` 迁移支的 r1 边界注释（"anything with a readable declared version"——已被 overlay 实测推翻的那句自述）改写为 r2 收窄后的真实规则；`config_reload_223_test.go` 语法错子例的同类过期注释一并改真；`loader.go` branch 1 尾注里"DOES declare a version belongs to the migration pipeline"改为"声明**更旧**版本才归迁移管线"。
- ⛔ J1（声明 99＋正文坏）按条形 1 归语法错——这是**现读裁**：99 的 newer-build 半句诚实，但"语法没问题"半句照样假，且声明未来版的坏文件根本没有迁移可跑。合法解析的未来版文件（`peekErr == nil`）不进新支，newer-build 句子原样保住——`loader_test.go:267 TestLoadFileNewerSchemaVersionRejected` 用的正是那形（`"schema_version = 99\n"` 干净一行），整包绿覆盖它。

### 2.3 新增常驻用例

`cmd/wisp/config_sentences_223r2_test.go :: TestTicket223R2FailureSentenceRouting`（新文件、8 枚子例、0.53s、同步驱动，不走 tick）：种文件 → `config.SaveFile`  canon → `config.NewManager` → 手改坏体 → `mgr.CheckAndReload()` → 把错误原样交**生产分类器** `describeReloadFailure`（与 `reloadOnce` 的调用逐字同路，`config_reload.go:153→155`）。每例断三件事：开的头是该句的 `cause=`、**不许借**另一支的 `cause=`、语法错支不许带「语法没问题」半句。表格双向钉：

| 形状 | wantCause | notCause（借到即红） |
|---|---|---|
| C1/E1/G1/L1（声明当前版＋正文坏，四形） | `cause=syntax` | `cause=invalid` |
| J1（声明 99＋正文坏） | `cause=syntax` | `cause=invalid` |
| A1（读不出版本） | `cause=syntax` | `cause=invalid` |
| A2（声明 1＋坏表头＝迁移管线的） | `cause=migration` | `cause=syntax` |
| 解析得开＋未知键 | `cause=unknown-key` | `cause=syntax` |

后两行是**反向牙**：把修法做成"一切坏文件都报语法错"会在 A2/未知键两例转红。首跑读数：8/8 PASS（原始件 `.scratch/wisp/probes/223/r2/routing-test-first-run.txt`，逐例 `PRODUCTION LINE` 原文在案）。

### 2.4 牙（把修法拿掉必须红）

拿掉方式＝把 `internal/config/loader.go` 整枚还原成**改前锚点的 blob**：`git cat-file blob c774f8da:internal/config/loader.go > internal/config/loader.go`（r2 新支随之消失；⛔ 未用 `checkout`/`restore`/`reset`/`stash`/`clean`）。跑 `go test ./cmd/wisp -count=1 -v -run 'TestTicket223R2FailureSentenceRouting'`，原始件 `.scratch/wisp/probes/223/r2/teeth-loader-removed.txt`（约 16:43，`EXIT=1`）：

- **读数＝恰红 5 发、其余 3 发照绿**——红的正是被修法接管的那 5 形，对照形（A1 读不出版本／A2 声明旧版／解析得开＋未知键）不动：

  `--- FAIL: .../声明当前版_注释以方括号开头_C1`、`.../声明当前版_CRLF_E1`、`.../声明当前版_无空格_G1`、`.../声明当前版_缩进版本行_L1`、`.../声明未来版_正文语法坏_J1`（5 枚）
  `--- PASS: .../读不出版本_A1_基线不动`、`.../声明旧版_坏表头_A2_仍归迁移`、`.../解析得开_未知键_不抢语法错`（3 枚）

- 红句逐字（5 枚同形，G1 为例）：

  > `config_sentences_223r2_test.go:111: shape 声明当前版_无空格_G1 is booked "cause=invalid detail=\"config.toml 语法没问题，但内容被校验拒绝（值不合法或引用解不开）。本次运行继续用内存里的旧配置\"", want it to open with "cause=syntax" - the two sentences swapped again`

  ——这就是 223-v1 推翻清单第 1 条实测到的那个错句，用例把它原文抄回。**用例有牙：拿掉修法必红。**
- 还原：`git cat-file blob HEAD:internal/config/loader.go > internal/config/loader.go`。**certutil SHA256 两哈希逐字**（改法后基线＝还原后，`diff` 空＋`git status --porcelain -- internal/config cmd/wisp` 空）：

  `14251d9fa271e4eb93c7d58961fa49a494480b3fb71f2904242f9dabb28b0d30` ＝ `14251d9fa271e4eb93c7d58961fa49a494480b3fb71f2904242f9dabb28b0d30`

  （原件落盘：`hash-loader-baseline.txt`／`hash-loader-restored.txt`）
- 还原后的整包单发（`final-packages.txt`，16:44:00→16:45:53）：`ok internal/config 0.787s`＋`ok cmd/wisp 110.376s`——**含 `migrate_test.go:123` 那枚既有断言在内全绿**（`internal/config` 是整包 `-count=1` 跑的），且 226 那枚靠时序过关的用例这一发也没有红（读数在案，归因照旧是票 227 的账，本腿不动）。

## 三把门（原始结论落盘）

| 门 | 读数 | 原始件 |
|---|---|---|
| `go build ./...` | **`BUILD_EXIT=0`**（16:46:02） | `.scratch/wisp/probes/223/r2/gates-final.txt` |
| `gofumpt -l <本腿动过的四枚文件>` | **零输出＝四枚全净**（新用例文件初稿被点名一次，`gofumpt -w` 修平后复检净） | 同上 |
| `sh scripts/d22scan.sh` | **`d22scan: clean - no D22 ban violations`**，`D22SCAN_EXIT=0`；口径现抄：bans #1-5 扫 248 枚生产 Go 件（internal 219＋cmd 29），ban #8 含注释与 `_test.go` 扫 internal 460＋cmd 63 枚（**本腿新增/修改的四枚全在射程内且未报**）；自检 `runtests.sh: OK - PASS=34 FAIL=0 SKIP=0` | 同上 |
| 附加：整包单发 `./internal/config ./cmd/wisp` | 两包 `ok`（`EXIT=0`） | `final-packages.txt` |

⛔ `thresholds.go`、golden、`allowlist.txt`、三枚冻结件：全程零接触（本腿连打开都没有打开过）；`tools/d22scan` 一字未动。

## 定性（涉及"会说反的话"字样，按派单要求自带三行）

① **现象出现在哪**：本机 `wisp run` 进程手改 `config.toml` 后重读失败时、写给操作员看的那句中文的**归句选择**（测试腿红则是用例的读时序）——纯文案路由与测试写法层面。
② **有没有本机被入侵的证据**：**没有**。零外部输入、零网络、零权限变化；所有失败路径的行为面从头到尾都是"本次重读不生效、内存继续用旧配置"（fail-kept），没有放宽过任何东西、没有绕过任何一张卡、SLO/阈值零接触。**这不是安全事件，是一次会说反的提示语＋一枚读早了的断言。**
③ **最坏后果是什么形状**：改前——操作员把语法改坏，却被告诉"语法没问题，是内容不合法"，**排查方向被指错**（该看那一行 TOML 写法的人去查值），当次热加载不生效、旧配置继续跑；改后——这句归位，且由常驻用例双向钉住。测试侧最坏是"拿掉产品句子用例照绿"的装饰风险，1.4 正控已把它排除。

## 没做完／判不了（具名）

1. **AC#5 的"整包终态"这把大尺本腿没有重跑**（`./cmd/wisp ./internal/...` 全量）：派单明令"不许跑整包门禁超过必要次数"，本腿只跑了动过的两枚整包（`internal/config` 全包＋`cmd/wisp` 全包各一发，均 `ok`）。v1 表在册的另外 6 枚红（`internal/ball` 1＋`internal/panel` 4 枚历史红、`internal/risk` 争用型 1）不属 r2 射程、本腿未碰；**"r2 之后 AC#5 该不该按终态翻勾"＝编排者复跑裁**（本表只交"本票用例自身那枚红已确定性堵上"的读数）。
2. **J1（声明未来版＋正文语法坏）归语法错是本腿的现读裁**，依据＝派单条形 1（读得出声明版本而正文语法坏 ⇒ 必须报语法错那一句）＋"声明未来版的坏文件没有迁移可跑"；若编排者认为该保 newer-build 句（v1 判它"半说反"而非"说反"），改法是给新支加 `ver == SchemaVersionCurrent` 半条件＋用例改一行期望——**我没有自作扩大成两可形状，把裁点具名交回**。
3. **`OnReload` 档在 `wisp run` 仍无生产赋值点**（v1 没做完第 3 条）：不在 r2 两格射程，本腿一字未动、复认仍在。归属仍待人拍板。
4. **票 227 那枚跨票账没有触发**：整包单发里 `always_write_no_clobber_226_test.go` 在内全绿（`final-packages.txt` 的 `ok cmd/wisp 110.376s`）；20 连跑那把尺只 `-run` 本票用例、**根本不跑 226 那枚**（不算它的读数）。pool-2 警告的"1s tick 被 20 连跑转起来"没有咬到它。**不修、不立账——本腿没观测到红。**
5. **翻勾判定不归本腿**：本腿是写腿（实现者），AC#4/AC#5 两格的成立与否按 D22 双角色只能出自非实现者新表＋编排者复跑；本表只交读数。
6. 绿名册逐名对拉（v1 没做完第 1 条）仍缺分母件口径，本腿未做（不属两格射程）。

## 时刻表（本腿现跑）

| 时刻(+0800，`date` 现读) | 命令 | 落盘 | 结果 |
|---|---|---|---|
| 16:27:20 | 起手：锚点 `c774f8da`；`git diff --cached --name-only`＝空；`ps -W`＝空 | 本表起手节 | 索引干净、无门在飞 |
| 16:2x | Write 骨架＋`git commit -F … -- <两枚 pathspec>` | 本表 | `28475620` |
| ~16:28→16:31 | **改前复现** 25 发（PATH→`go test ./cmd/wisp -count=1 -run 'TestTicket223RestartTierSaysItWillNotApply'`，台件 `shots_223r2.sh`） | `pre-fix-repro.txt` | **25 发 1 红（shot 3，红句 `:453`）** |
| 16:31→16:34 | 四枚 Go 件改动（awaitStdout／loader 新支／新用例／三处注释归真）＋`go build ./...`＋`go vet` 两包 | — | 净 |
| ~16:35 | 新用例首跑（`-count=1 -v -run TestTicket223R2FailureSentenceRouting`） | `routing-test-first-run.txt` | 8/8 PASS 0.53s |
| 16:38:55→~16:40:3x | **修后判据 20 连发**（同一把尺；跑前 `ps -W`＝空） | `flake-20.txt` | **20/20 EXIT=0** |
| ~16:41 | `certutil -hashfile … SHA256` 基线两枚 | `hash-*-baseline.txt` | 见 1.4／2.4 |
| 16:41:57→16:42:43 | **任务一正控**：产品 stdout 句子注释掉后一发 `-v` | `positive-control-m1.txt` | **FAIL 42.18s**；`git cat-file blob HEAD:…` 还原，哈希逐字相同 |
| ~16:43 | **任务二牙**：`git cat-file blob c774f8da:… > loader.go` 后一发 `-v` | `teeth-loader-removed.txt` | **恰红 C1/E1/G1/L1/J1、对照三发绿**；`HEAD` 还原，哈希逐字相同 |
| 16:44:00→16:45:53 | 整包单发 `go test ./internal/config ./cmd/wisp -count=1 -timeout 20m`（跑前 `ps -W` 现查＝空） | `final-packages.txt` | `ok 0.787s`／`ok 110.376s` |
| 16:46:02→16:46:20 | **三把门**：`go build ./...`／`gofumpt -l`（四枚）／`sh scripts/d22scan.sh` | `gates-final.txt` | 0／净／clean |
| 收工 | 票面追加一行 Status（旧两条逐字保留）；⛔ 不翻勾、不 `-done`、不 push | 票面＋本表 | 各节已齐 |

## 交回编排者的六节（大白话）

**① 两件事各自动了哪几枚文件、行数**（`git diff --numstat c774f8da..HEAD` 原文：`11 6 cmd/wisp/config_reload.go`、`51 10 cmd/wisp/config_reload_223_test.go`、`124 0 cmd/wisp/config_sentences_223r2_test.go`、`24 6 internal/config/loader.go`；代码＝`5a755c3c`，骨架＝`28475620`）——任务一只动测试件（helper＋一处读法改轮询＋注释归真，⛔ 产品文案零改动）；任务二动 `loader.go`（新支）＋新用例一枚新文件＋`config_reload.go` 纯注释归真。

**② 任务一**：改前 25 发复现出 1 红（第 3 发，红句与历史逐字同形），改后**同一把尺 20 发 exit 码全 0**（每发 `ok`、耗时 1.957～2.285s 逐发在案）。正控＝把产品那句"本次运行不会生效"注释掉再跑：用例 **42.18 秒红**，红句逐字 `stdout never carried "本次运行不会生效" within 40s`，而同一发 stderr 里两行 restart 审计仍在——**轮询不是恒真的装饰，是真在等那句话**；还原前后 SHA256 逐字相同（`542f705a…34f4670`）。

**③ 任务二改前／改后两句原文并排＋新增用例有牙读数**——同一枚文件（声明 `schema_version = 2`、正文语法坏）：
- 改前那句（错）：`cause=invalid detail="config.toml 语法没问题，但内容被校验拒绝（值不合法或引用解不开）。本次运行继续用内存里的旧配置"`
- 改后那句（对）：`cause=syntax detail="config.toml 读到了但解析不了：这一行不是合法 TOML 语法（不是权限、不是缺失）。本次运行继续用内存里的旧配置；修好之后要再出现一次新的 mtime/大小才会被重读"`

新增常驻用例 `TestTicket223R2FailureSentenceRouting`（8 形双向钉）首跑 8/8 绿；**拿掉修法必红**：还原成改前 loader 后恰红 C1/E1/G1/L1/J1 五发（红句逐字含 `the two sentences swapped again`，把错句原文抄回），对照三发（读不出版本／声明旧版／解析得开＋未知键）纹丝不动——迁移管线那格（`migrate_test.go:123` 的既有断言）一字未动、`internal/config` 整包仍绿；还原哈希逐字相同（`14251d9f…b28b0d30`）。

**④ 三把门原始结论**：`go build ./...` 退出 0；`gofumpt -l` 对本腿四枚文件零输出；`d22scan` "clean - no D22 ban violations"（口径：bans #1-5 生产 248 枚，ban #8 含注释与 `_test.go` 扫 internal 460＋cmd 63 枚——本腿四枚全在射程内未报）。附加整包单发：`internal/config` 0.787s `ok`＋`cmd/wisp` 110.376s `ok`。原始输出全在 `gates-final.txt`／`final-packages.txt`。

**⑤ 票 223 该不该翻勾**：**本腿是写腿，不自翻勾；七格一格没勾、票面不 `-done`**（派单令＋D22 双角色）。能交回的只有读数：v1 表里 AC#5 那枚"本票自己的红"已确定性堵上（20/20＋正控非恒真），AC#4 那句会说反的话已归位（有牙常驻用例双向钉住），全部现跑于含本腿改动的锚点。AC#5 的整包终态与两格翻勾＝等编排者复跑＋非实现者新表裁；对 J1 归句若有另一判法，见"没做完"第 2 条具名改法。

**⑥ 没做完／判不了**：整包终态大尺没重跑（派单限制）；J1 归语法错＝现读裁、裁点具名交回；`OnReload` 档仍无人接（不在射程，复认仍在）；票 227 的 226 时序账本腿没观测到红、不动；翻勾不归实现者；绿名册逐名对拉仍缺分母件口径。（涉及"会说反／权限"字样的三行定性见专节——这两件事都不是安全事件。）
