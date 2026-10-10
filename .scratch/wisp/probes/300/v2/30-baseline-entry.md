# 300-v2 — `30` 一进门就绿 ＋ 配对基线"新增红 0"（我这把独立复现）＋ 五枚基线红的归因

## 1. 一进门就绿（工作树，窄档，⛔ 引用编排者的读数）

```
$ GOFLAGS= go test ./internal/audio/ -count=1 -run 'TestParseWaveFormatSubFormatOffset300' -v
rc=0
pass-all(grep -c -- '--- PASS')=5      pass-top(grep -cE '^--- PASS')=1      fail-all=0
```

⇒ 与编排者 13:4x 那一发同形（4 枚子项＋1 枚顶层、`rc=0`）；导出树那一侧另有一发同值（`10` 件 §1 第一行，`rc=0`／5／1）。
⇒ **两形各自一发、互相独立**（工作树↔`41329475` 的 blob 已由 `10` 件 §2 两枚 `IDENTICAL` 钉成同一套字节）。

## 2. 配对基线（★我这把**总共两发全量**，逐字同一形，只换树）

命令逐字（腿的 §5 配方，我这把⛔ 改射程）：

```
cd <树> && PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./internal/audio/ ./cmd/wisp/ -count=1 -v
```

| 发 | 树 | 来源号 | `pass-all` | `pass-top` | `fail-all` | `fail-top` | `skip` | `load0xc` | `panic` | rc | `cmd/wisp` 耗时 |
|---|---|---|---|---|---|---|---|---|---|---|---|
| 改前 | treeA | `4eb29312`（＝`028529fc^`） | **432** | **320** | 5 | 5 | 4 | 0 | 0 | 1 | 416.140s |
| 改后 | treeH | `41329475`（起手锚） | **437** | **321** | 5 | 5 | 4 | 0 | 0 | 1 | 411.017s |

- ★**跑了几发，具名**：全量＝**2 发**（改前 1＋改后 1）⛔ 腿的四发（它每侧 2 发）。我这把的两发**逐枚复现它的读数**（432／320 ↔ 437／321，⛔ 一枚数不符）。
  增量＝`437-432=5`＝1 枚顶层＋4 枚子测试、`321-320=1`＝只顶层那一枚 ⇒ 增量**全部来自新用例**（另有一把：`grep -c TestParseWaveFormatSubFormatOffset300` ⇒ 改前 **0**、改后 **10**）。
  ⚠ 我⛔ 交"每侧两发的同侧对称差"那一格（腿交了，`comm -3` 零行）：我这把只有两发、每侧一发，那一格**归它的读数＋编排者复跑**，我具名标〔仅读码推的〕。
- ★**新增红作差（我这把自己的两发之间）**：

```
$ comm -13 <改前红名册> <改后红名册> | wc -l   → new-red=0
$ comm -23 <改前红名册> <改后红名册> | wc -l   → gone-red=0
$ comm -3  …                                  → sym-diff=0
```

⇒ **"新增红 0"这一格我这把独立成立**，⛔ 靠腿的 `comm`。
- ★**载入级坑⛔ 咬到我**：两发 `load0xc=0`、`panic=0`、`internal/audio` 那半逐发 `ok … 16.1s/16.5s` ⇒ 两棵树的 `cmd/wisp` **真跑了**（⛔ 死于 `0xc0000135`）；
  前置那一跳我按 `00` 件 §3 先 `mkdir -p`＋`cp` 再 `ls …/*.dll | wc -l`＝**3**（腿在 `30` 件 §4 报的那枚"目录⛔ 存在"的坑，我这把⛔ 重踩＝第一跳就带 `mkdir`）。

## 3. 五枚基线红：逐名定位（我自己的尺）＋ **归因分两档**（★这一档腿没分、编排者⛔ 量过）

```
$ grep -rn "func <名>(" --include=*_test.go cmd internal        # 逐名，treeA 里跑
TestAC14AwaitedBindingReplyReachesThePage   cmd/wisp/panel_resident_windows_test.go:812
TestAC14GoSideEvalPushReachesThePage        cmd/wisp/panel_resident_windows_test.go:855
TestAC4FocusReturnToPriorWindowGap33r5      cmd/wisp/panel_host_windows_test.go:875
TestPanelHostRealWindowHopAndLifecycle      cmd/wisp/panel_host_windows_test.go:614
TestCleanCheckoutBuilds_AC11                cmd/wisp/panel_host_gate_test.go:385
```

⇒ **5/5 全在 `cmd/wisp`**（`panel_resident_windows_test.go`×2／`panel_host_windows_test.go`×2／`panel_host_gate_test.go`×1），⛔ 一枚落 `internal/audio` ⇒ **与腿那枚名册逐枚同值**（含文件与枚数分布）。⛔ 一枚归因给本票（作差两侧同名单、且本票射程⛔ 含 `cmd/wisp`）。

红句逐字（⛔ 只数颜色，坏法种类分开，全取自我这把 `full-preA-1.txt`）：

- `panel_resident_windows_test.go:821: no report "ac14r-0" from the page within 15s (what DID arrive at the door: nothing at all). AC#14's reply hop cannot be decided without the page's own answer, a…`
- `panel_host_windows_test.go:948: the panel did not take the foreground on Show: foreground 0x50102, panel hwnd 0x36a0786 (AC#4 says only the panel takes focus when shown)`
  ＋ `:956: Hide did not hand focus back to the window that held it before Show: expected 0xd80b8e, foreground is 0x50102`
- `panel_host_windows_test.go:660: cold bring-up measured on this box: -1.000 ms (HEAD HEAD-unknown at read time, 2026-10-10T13:51:34+08:00)` ＋ `:662: cold bring-up did not produce a browser round trip (got -1.000); the message channel did not come up on the real window`
- `panel_host_gate_test.go:395: git rev-parse --show-toplevel: exit status 128` ← ★**这一枚下面单列**

### 3.1 ★归因分两档：其中**一枚是测量装置自己造的**（本票新账，⛔ 记在腿头上）

- **档 A（装置红，1 枚）＝`TestCleanCheckoutBuilds_AC11`**：尺＝读源码 `cmd/wisp/panel_host_gate_test.go:385-400` ⇒ 它先 `repoRootForTest(t)`，而 `git archive` 的导出树**没有 `.git`** ⇒ `git rev-parse --show-toplevel` 退 **128** ⇒ `t.Fatalf`。
  **我这把在真实工作树里单发复跑（⛔ 开窗、⛔ 载具：该用例只 `go build ./...` 到 temp 目录）**：

```
$ PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" GOFLAGS= go test ./cmd/wisp/ -count=1 -run 'TestCleanCheckoutBuilds_AC11' -v
rc=0
--- PASS: TestCleanCheckoutBuilds_AC11 (26.10s)
```

  ⇒ **它在工作树里是绿的、只在导出树里红** ⇒ "基线 5 枚红"这一句在**方法层**该写成：**4 枚（真窗口／WebView2 族）＋1 枚（导出树⛔ 有 `.git`＝装置红）**。
  ⚠ **对本次判语的影响＝零**：它在**两侧同色**（两发都红）⇒ `comm` 三把皆 0 行⛔ 变；我用的是对称差、⛔ 单发红名当结论（派单点名的那条）。
  ★并回＝**`A802` 那条"导出树⛔ 自带 handicap"的账⛔ 只算到 DLL 一枚**；`cmd/wisp` 里**按 `.git` 存在性**取数的用例（`repoRootForTest`／`HEAD-unknown`）是**同一枚坑的第二面**，今后凡拿导出树当"基线红名册"，**先问一句"这枚红依⛔ 依赖 `.git`"**。
- **档 B（真窗口／WebView2 族，4 枚）**：四枚红句逐枚都在说"page 没回话／面板⛔ 拿到前台／真窗口冷起＝-1.000ms"，且 `TestPanelHostRealWindowHop…` 那行还印着 `HEAD-unknown at read time` ⇒ 同族里**至少这枚的读数面也沾装置**。
  ⇒ **我这把⛔ 能判它们在工作树里绿**：跑它们＝在机主机器上真开窗／抢输入 ⇒ 归口"仅本机可量、⛔ 归腿⛔ 归我"＝**归编排者**（同 `AC#3` 那一发）。
  判据层面⛔ 需要它：两侧同名单同数已经足够把"新增红 0"钉住。
- ⚠ 派单点名的那枚本机噪声（`design/assets/*` 被工作树删 ⇒ `TestC21DesignTokensFourWayAgree` 一类只在本地红）：我这把**⛔ 撞上**（尺＝`grep -c TestC21` 两份全量日志 ⇒ **0**）⇒ 它⛔ 在 `./internal/audio/`＋`./cmd/wisp/` 这两枚包的分母里，`D design/*` 那 16 枚对本程判据**无影响**（我这把⛔ 需要拿它给任何一枚红开脱）。〔盘上现量〕
