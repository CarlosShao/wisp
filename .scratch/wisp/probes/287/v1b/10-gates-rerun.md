# 287-v1b §1 门禁复跑读数（派单第 0 步第 1 项：五件＋格式两把，全部本腿自己跑）

时钟：2026-10-16 10:0x +08 起手｜腿＝287-v1b（验收腿·裁决者·⛔ 实现者）
HEAD 现量＝**五件门禁跑在 `1a651618d6a5b71e06cb216f1a309304becfd8aa`**（＝本腿第 0 步第 3 项那笔 287-v1 遗产代提笔，量 `git rev-parse HEAD`于跑前）；
写本件时盘上 HEAD 已推进到 `37291a1d4de5760d4b923f3f133cf477eaf900dd`（＝编排者一笔，`git show --stat` 现量只动三枚 `.md`）⇒ **读数面不变**，本件按"跑的那一发"记账。
⚠ 具名过期：派单第 2 行写的"当前 HEAD＝`46007856`"在本腿起手时已不是 HEAD ⇒ 本件所有读数都在上面那一发上量。
★可比的证据：`git log --oneline 46007856..HEAD | wc -l`＝**16**（含本腿那笔遗产），而
`git diff --stat 46007856..HEAD -- internal/tools`＝**空**、夹具那枚 blob 在两发 HEAD 同哈希
（`git rev-parse 46007856:internal/tools/subagent_depth_fixture_287_test.go`＝`…HEAD:同路径`＝`a2681ead…`，量）
⇒ 287-r1 的落地面（夹具＋`subagent_197.go`）一字未动，本腿读数与前任 `dce9eb90` 那一发同面可比。

台面起手（量，`git status --porcelain -- <面>`）：`internal/winsec`＝0 行／`tools/d22scan`＝0 行／
`internal/tools`＝0 行／`cmd`＝0 行 ⇒ 同批两枚写腿（`322-r1b`／`325-r1`）在**本腿起手时已收**，
派单第 0 步第 2 项那枚雷（`go build ./...` 被别腿在飞件弄红）**没有响** ⇒ 本腿门禁一律在**仓内工作树**跑，
只有突变台件跑在仓外导出树（§4）。

## 1.1 五件门禁：本腿自己跑的数

| 尺 | 命令（全部仓根，GOFLAGS 置空＝门禁名义那把） | rc | 本腿读数 | 与 287-v1 那发比 |
|---|---|---|---|---|
| 编译 | `GOFLAGS= go build ./...` | **0** | 日志 **0 字节**（v1b-01） | **同**（v1＝rc 0／空） |
| 包级 | `go test ./internal/tools/ -count=1 -v` | **0** | 顶层 `^--- PASS`＝**242**、`^--- FAIL`＝**0**、`^--- SKIP`＝**0**；`^=== RUN`＝**325**；子测试 `^    --- PASS`＝**83**；`^ok  `＝**1**（v1b-07，157,003 字节） | **五把尺逐枚同数**（v1＝242/0/0、RUN 325、子 PASS 83）⇒ 两发独立复跑对上 |
| 静态 | `go vet ./internal/tools/` | **0** | 日志 **0 字节** | **同** |
| D22 | `sh scripts/d22scan.sh` | **0** | `clean - no D22 ban violations`；ban #8 `internal/`=**572**、`cmd/`=**138**；bans #1-5 `internal/`=**239**、`cmd/`=**41**；ban #7 `internal/tools/`=**25**（v1b-03，23,749 字节） | **⛔ 同**：v1＝`internal/` **571** ⇒ 本腿 **572**（+1）。成因＝两发 HEAD 之间 `322-r1b` 落了 `internal/winsec/openleaf_322_other_test.go` 拆分（`git diff --stat 46007856..HEAD` 现量＝winsec 三枚＋d22scan 一枚，`internal/tools` 空）⇒ **别腿入库件令分母上涨，⛔ 归 287、⛔ 当红**（ban #8 那一支含 `_test.go`，见 §1.3） |
| 长度门 | `sh scripts/check-path-length-budget.sh --with-self-test` | **0** | `VERDICT GREEN`、`tracked paths=**13688**`、`over-budget=60`、`covered by roster=60`、`not in roster=0`；正控 **control 1/3 ok／2/3 ok／3/3 ok**（三档全 ok，v1b-04，2,535 字节） | `over-budget=60`／`roster=60`＝**同**；`tracked paths` v1＝**13455** ⇒ 本腿 **13688**（+233）。成因＝同批发出的件与两笔遗产入库（⛔ 287 的病）。⚠ 措辞差：v1 写"控制组 2/3 ok"，本腿日志逐字是**三行各带 `control N/3 ok`**（＝三档全过）⇒ v1 那句读起来像"只过了 2/3"，**本腿日志否证那一读**，属前任措辞腐坏、⛔ 判据差 |
| 格式（票 287 两枚件） | `gofmt -l internal/tools/subagent_depth_fixture_287_test.go internal/tools/subagent_197.go` | 0 | **列出 0 枚件名**（`gofmt -l` 退码恒 0＝空尺，只能这样写；v1b-05＝0 字节） | **同** |
| 格式（对照，整目录） | `gofmt -l internal/tools` | 0 | **列出 1 枚**：`internal\tools\task_output_pointer_notice_errtext_174r7_test.go`（v1b-06，64 字节） | **同**（v1 亦 1 枚、同一件）⇒ 存量债，来路 `67a7ab41`（票 174-r7），⛔ 票 287 写面、⛔ 记到 287-r1 头上 |
| 格式（门禁名义那把） | `gofumpt -l` | **127** | `command not found`：`command -v gofumpt` 空（本腿现量，正控＝同一把尺对 `gofmt` 回 `/d/work/base/go/bin/gofmt`） | **同**（v1 亦 rc=127）⇒ 台面没装。本格只能由 `gofmt -l`（射程更窄）代交，⛔ 冒充"门禁那把绿" |

### 台面故障归类（⛔ 计入任何判据）
- `gofumpt` rc=127 一枚（工具缺失）。
- 包级日志里故障形状扫描：尺＝`grep -inE 'walk|build failed|no tests to run|connection refused|0xc0000135|panic|test timed out'`
  命中 **2 行**，逐枚读＝同一枚**既有子测试的名字** `TestTicket224GrantSourceThatCannotAnswerMeansAsking/panics_mid-lookup`
  （`:373` RUN／`:376` PASS）⇒ **零真红、零台面故障**（量，射程＝v1b-07 全日志）。

### 三枚夹具用例在本腿整包跑里的颜色（逐枚，`grep '^--- PASS: Test287'`）
- `Test287SubagentCallerRefusedAtTheRosterDepthCheckThroughTheRealLoop`（主用例，件 `:155-200`）＝ **PASS 0.11s**
- `Test287RootCallerWithTheSameShapeStillDerives`（反形正控，`:203-230`）＝ **PASS 0.19s**
- `Test287DepthRefusalWinsOverThePoolRefusal`（次序尺，`:233-260`）＝ **PASS 0.11s**

### 1.3 复用的两把尺（本腿自己跑，⛔ 抄前任的射程结论）
- `sh scripts/d22scan.sh` 日志逐字写着 `ban #8 internal/ … comments and _test.go included`（＝**含** `_test.go`），
  而 bans #1-5 那一支的分母 239 < 572 ⇒ 两支射程不同这一条本腿**独立复量成立**。
- `internal/` 的 Go 文件枚数两把尺：ban #8＝**572**（含 `_test.go`）／bans #1-5＝**239**（不含）——差 333 枚＝测试件与豁免件。
