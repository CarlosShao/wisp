# `305-a1` · `00` 起手锚＋台面＋Go 命令自陈

## 起手锚（现跑，`logs/l1-anchor.txt`）
- `git rev-parse --short HEAD` = **`a81c2980`**（全号 `a81c2980` 分支 `dev`）⇒ **与派单写的锚一致**，无需报差。
- `git status --porcelain -- internal cmd docs scripts .github` = **0 行**（件里逐字可见：该节下无任何输出行）。
  ⇒ 本腿引用的两枚文件（`cmd/wisp/panel_host_windows.go`、`cmd/wisp/panel_resident_windows_test.go`）
  **工作树＝`a81c2980` blob**，所以"读工作树"与"读 `git show HEAD:`"在这两枚上同一物；
  但**所有行号锚我仍按派单要求走对象层**——尺＝`git show a81c2980:<path> | sed -n '<N>p'` 逐枚比对该行含预期符号，
  **28 枚全 MATCH**，逐字名册见 `logs/l7-anchor-verify-head-blob.txt`。
- 工作树本来就有别人的脏面（`design/**`、`.gitignore`、别家 `probes/**`）⇒ 本腿⛔ 动、⛔ 暂存；只往 `.scratch/wisp/probes/305/a1/**` 写。

## 台面（哪一把台面，量的是什么）
本腿**⛔ 跑任何用例**，所以本件里没有一枚"新颜色"。所有结论建立在：
1. **编排者已有的四件读数**（母仓改前／母仓改后／clone `f718e9b6`／clone `807497c1`）＋托管 CI 那件 `r9` —— 逐字读回、作差，见 `20-`。
2. **HEAD 对象层的源码＋依赖库源码**（静态读码，零执行）—— 见 `10-`／`30-`。
3. **`docs/evidence/s1/` 里 10-01 那两发的既有读数**（票 33 的 r6/r7 件）—— 见 `40-` 与本件下面的"与票面冲突"。

## 与票面/派单冲突处（具名，先摆在最前面）
★**票面第一句标题与其前提⛔ 盘上证据**（尺＝`logs/l4-earlier-green-readings.txt`，`grep -rn -E "PASS: TestAC13ColdStart|PASS: TestAC14GoSideEvalPush" docs/evidence/s1/`）：
- `docs/evidence/s1/33-panel-host-c27-r7.md` 的整包名册里**逐字同时有**
  `--- PASS: TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe (0.97s)`（该件 `:201`，另两发 `:241`/`:277` 亦 PASS）
  与 `--- PASS: TestAC14GoSideEvalPushReachesThePage (0.50s)`（该件 `:207`）。
  该件起手锚（`:10`）HEAD＝**`416d9d56`**；那一发是**整包 rc=0／PASS=240／FAIL=0／SKIP=0**，
  ⇒ SKIP=0 意味着 AC13 没跳 ⇒ **那一发是在一份带产物 bundle 的台面上跑的**（干净 clone 里它必 SKIP，见 `r3`/`r4` 的 `:317` 具名跳过句）。
⇒ **后果**：票面标题那句"`Eval` 推不进文档那一枚**早在回归之前就红、从来没绿过**"与票面 判语① 的隐含读法（"从没做通过的功能"）**⛔ 成立**：
nail2 与 AC13 都在 10-01 的带 bundle 台面上**绿过**，两枚今天都红 ⇒ 两枚都落在"某时绿过、后来红"的形态里。
票面"⛔ 已知改前绿过"那一行（第 18 行）只在**票 303 那次改前那一发**（`cf46c24`）的范围内真，⛔ 覆盖到 10-01。
（本腿只交料：**`AC#2` 要不要按"缺功能"降级＝编排者裁**。）

★**票面 `:660` 那一族行号在另一枚文件里**：CI 那四枚 `TestPanelHostRealWindowHopAndLifecycle` 的 `:660`/`:662`/`:665` 属
`cmd/wisp/panel_host_windows_test.go`（该枚用例从 `:605` 起），⛔ 属 `panel_resident_windows_test.go`；
本腿在 `l7` 里按后者去核 `:660`，那一条对不上（逐字打印在件里）。⇒ 这是**引用面写错文件**，⛔ 行号腐烂那一族；记下来免得下一腿再踩。

## Go 命令自陈（⛔ 编译面；本腿真跑了什么）
| 命令 | 用途 | 改任何码？ |
|---|---|---|
| `go list -m -f '{{.Dir}}' github.com/jchv/go-webview2` | 把依赖库源码目录名解析出来，好只读它的 `webview.go`／`pkg/edge/chromium.go` | **⛔**（零写入，退码 0，件 `logs/l3-library-hops.txt`） |
| ⛔ `go test` / ⛔ `go build` / ⛔ `go vet` / ⛔ 任何 `cmd/wisp` 真跑 | —— | **一枚没跑**，所以本件⛔ 任何新颜色 |

派单许可的就是这两行（`go env`／`go list` 可以跑）；本腿只用了 `go list`，**没用 `go env`**（⛔ 需要）。

## 本腿写面
只 `.scratch/wisp/probes/305/a1/**`，全部 `.md`／`.txt`；**⛔ 新建 `.go`**（`find` 自查见 `90-final.md`）。
