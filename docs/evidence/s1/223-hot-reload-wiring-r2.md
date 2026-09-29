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

〔待填：等 20 发跑完〕

### 1.4 正控（拿掉那句产品文案必须红）

〔待填〕

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

〔待填：还原 `c774f8da` 版 loader.go 后的红句原文、还原读数与哈希〕

## 三把门（原始结论落盘）

| 门 | 读数 | 原始件 |
|---|---|---|
| `go build ./...` | 〔待填〕 | 〔待填〕 |
| `gofumpt -l <本腿动过的文件>` | 〔待填〕 | 〔待填〕 |
| `sh scripts/d22scan.sh` | 〔待填〕 | 〔待填〕 |

## 定性（涉及"会说反的话"字样，按派单要求自带三行）

〔待填：①现象出现在哪 ②有没有本机被入侵的证据 ③最坏后果是什么形状〕

## 没做完／判不了（具名）

〔待填；含：若 20 连跑里冒出 `always_write_no_clobber_226_test.go`（票 227 射程）的红——不修、具名登记〕

## 时刻表（本腿现跑）

| 时刻(`date` 现读) | 命令 | 落盘 | 结果 |
|---|---|---|---|
| 〔待填〕 | | | |

## 交回编排者的六节（大白话）

**① 两件事各自动了哪几枚文件、行数**〔待填：`git diff --numstat` 原文〕
**② 任务一 20 发逐发 exit 码与正控红句**〔待填〕
**③ 任务二改前／改后两句原文并排＋新增用例有牙读数**〔待填〕
**④ 三把门原始结论**〔待填〕
**⑤ 票 223 该不该翻勾**〔待填——本腿是写腿，不自翻勾；只报"我这边两格缺口堵没堵上"〕
**⑥ 没做完／判不了**〔待填〕
