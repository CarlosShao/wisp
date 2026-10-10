# 303-a1 `13-ac0-repro.md` — 票 303 `AC#0` 复现钉死（⛔ 归因，⛔ 修复；本格只交读数）

腿＝`303-a1`。这一件覆盖票面 `AC#0` 的 ①②③④ 四小格。所有读数都是**本腿现跑**，件路径都具名。

## 0. 台面（三发读数的口径）

| 发 | 台面 | 被测树 | 起法 |
|---|---|---|---|
| A | 母仓 `D:/work/workspace/projects plans/Wisp` | 工作树（起手与收尾 `git status --porcelain -- cmd internal scripts .github docs` 均 **0 行** ⇒ 工作树≡HEAD；那一发时 HEAD＝`b3f9c6d7`，我此后三笔只动 `.scratch/wisp/probes/303/a1/**` ⇒ `cmd/**`／`internal/**` 面与现 HEAD 逐字节同） | `go test -count=1 …`＋`-count=3 …` |
| B | **仓外 clone** `$HOME/wisp-303-bisect`（＝`C:/Users/swq/wisp-303-bisect`，`git clone` 母仓，无网络） | `git checkout cc315261` | 同一条尺（台件 `bisect-step.sh`） |
| C | 同一枚 clone | `git checkout bcd0a543` | 同一条尺 |

clone 的**干净度现量**（这台面为什么算"干净"）：`git ls-files frontend/dist` → 只有 `frontend/dist/.gitkeep`；clone 里 `ls -d third_party build` → **两枚都不存在**（根 `.gitignore` :12 `third_party/`、:14 `build/`）。
⇒ clone 里 `frontend/dist` 的产物字节＝**0 枚**，而 `go:embed` 的目录仍在（`.gitkeep` 顶着，票面 `:17` 与 `.gitignore` :18-23 的注释就是这个意思）。

## ① 起点／终点具名读数（两枚终点分开留）

| 发 | rev | 色 | 逐字 | 件 |
|---|---|---|---|---|
| B 起点 | `cc315261`（10-06 那台 CI 的同码，本机 clone） | **PASS** | `--- PASS: TestAC14AwaitedBindingReplyReachesThePage (1.59s)`；同发 :825 逐字 `AC#14 nail 1 (reply hop), page's own words: "REPLIED,REPLIED,REPLIED" (Go's handler was reached by 3 of the 3 real requests)`；go rc=0 | `$HOME/wisp-303-steps/step001-cc315261….txt` |
| C 终点甲 | `bcd0a543` | **FAIL** | `--- FAIL: TestAC14AwaitedBindingReplyReachesThePage (20.02s)`；`panel_resident_windows_test.go:821: no report "ac14r-0" from the page within 15s (what DID arrive at the door: nothing at all)…`；go rc=1 | `step002-bcd0a543….txt` |
| A 终点乙 | 现 HEAD（`cmd/**` 面≡`b3f9c6d7`） | **FAIL** | 同一句红，`--- FAIL (20.02s)`，件内 :7 命中判据句，`rc=1` | `10-repro-head-count1.txt` |

⇒ `bcd0a543`（票面的"改后那发 CI"）与当前 HEAD **各留了一发**，两枚都红、红句逐字同；起点 `cc315261` **在同一台机器、同一个干净 clone 里**照旧拿到页面回执（这条比 CI 归档件更强：⛔ "镜像漂"、⛔ "我那台机器的毛病"两支都⛔ 依据）。

## ② 最小复现集＝一枚用例名，跑法逐字

```sh
# 母仓台面（A）
cd "D:/work/workspace/projects plans/Wisp"
export PATH="/d/work/workspace/projects plans/Wisp/third_party/sherpa-onnx:/d/work/workspace/projects plans/Wisp/build:$PATH"
go test -count=1 -timeout 420s -v -run 'TestAC14AwaitedBindingReplyReachesThePage' ./cmd/wisp/
```

```sh
# 仓外 clone 台面（B／C）——PATH 铺的是【母仓】那两枚目录的绝对路径，clone 自己没有 DLL
cd "$HOME/wisp-303-bisect"
export PATH="/d/work/workspace/projects plans/Wisp/third_party/sherpa-onnx:/d/work/workspace/projects plans/Wisp/build:$PATH"
go test -count=1 -timeout 420s -v -run 'TestAC14AwaitedBindingReplyReachesThePage' ./cmd/wisp/
```

**PATH 那两枚目录（逐字，本腿每一发都铺的就是这两枚）**
1. `/d/work/workspace/projects plans/Wisp/third_party/sherpa-onnx`（`onnxruntime.dll`／`sherpa-onnx-c-api.dll`／`sherpa-onnx-cxx-api.dll`）
2. `/d/work/workspace/projects plans/Wisp/build`（同三枚 DLL 的副本）

⚠ 形状：用 **shell 形式**（`/d/…`），⛔ `D:/…` 形式。出处＝`scripts/wisp-cli-tests.sh:101-109` 的实测注释（`D:/` 那形＝`0xc0000135`）。
⚠ 缺 DLL 的形状＝`exit status 0xc0000135` 且**零 `--- FAIL`**＝用例根本没跑（`AC#1` 里这叫**无效步**，⛔ 好⛔ 坏）。

**它⛔ 依赖 `frontend/dist`——现量回答（两路各一发）**
- 路①（直接验，票面 `:17` 要的那一发）：**在干净 clone 里**（`frontend/dist` 只有 `.gitkeep`、产物字节 0 枚）那一枚用例
  起点 `cc315261` **PASS 并拿到 `REPLIED,REPLIED,REPLIED`**、终点 `bcd0a543` **FAIL 且红句就是判据句**（⛔ 走具名跳过）。
  ⇒ 它**无对象可跳**：dist 为空时它照旧跑、照旧判（件 `step001`／`step002`）。
- 路②（码层，同尺可重跑）：尺＝`grep -nE 't\.Skip|dist|embed' cmd/wisp/panel_resident_windows_test.go`
  那枚具名跳过 `AC#13 has no subject in this tree: the embed resolves no entry` 在 **:317**，属 `TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe`；
  `TestAC14AwaitedBindingReplyReachesThePage` 起于 **:812**，体内⛔ embed/dist 跳过，输入只有 `startPanelForTest`／`showAndWait`／`evalOnPanelThread`／`awaitReport`。
- ⚠ 反面对照（母仓台面）：母仓 `frontend/dist` **有旧产物字节**，所以 `TestAC13…` 在母仓有对象、跟着一起红（`probes/302/orch/g4-local-three-cases.txt` 里它给的是 `AC#13 probes from the resolved entry (1044 bytes): 1 id(s) [root]`）；
  同一枚用例在干净 clone 里会走 `:317` 的具名跳过 ⇒ **⛔ 它是最小复现集**（票面 `:17` 的陷阱本腿认了，并用路①那发把它钉死）。
  这一发的逐字件＝`15-ac13-skips-in-clean-clone.txt`（bisect 跑完之后补，⛔ 并发跑两发 go）。

## ③ 复现率＝同一条命令 `-count=3` 的三色

尺＝A 台面，逐字＝上面 ② 那条命令把 `-count=1` 换成 `-count=3`（件 `11-repro-head-count3.txt`，`rc=1`，2026-10-10 17:07）：

| 色 | 尺 | 读数 |
|---|---|---|
| FAIL | `grep -cE '^--- FAIL: TestAC14AwaitedBindingReplyReachesThePage'` | **3** |
| PASS | `grep -cE '^--- PASS: TestAC14AwaitedBindingReplyReachesThePage'` | **0** |
| SKIP | `grep -cE '^--- SKIP'` | **0** |
| 判据句 | `grep -cF 'no report "ac14r-0" from the page within 15s'` | **3** |

⇒ 复现率＝**3/3 红**（同机、同树、连着三发，每发都建了真窗）。
⚠ 口径：这是 `-count=3` 同一**进程**里连跑三发；⛔ 三次独立进程。派单说"⛔ 单发当恒红"——这里交的就是三发的色，⛔ 由本腿判"恒红"。

## ④ "这⛔ 是我这台机器的毛病"的对照＝同码在 CI 也红那一枚

**⛔ 为这一句推 CI**（票面 `AC#0` ④ 与派单同令）。引用盘上归档件：
- 改后那发（`bcd0a543`，托管 runner `windows-2025-vs2026`／agent `20260925.250.1`）：
  `.scratch/wisp/probes/301/orch/logs/ci-after-cli-block.txt` **:1360** 逐字
  `no report "ac14r-0" from the page within 15s (what DID arrive at the door: nothing at all)`（另有 :1335 `"ac13-probe"`、:1367 `"ac14-push"`）
- 基线那发（`cc315261`，同镜像同 agent 版本）：同目录 `ci-baseline-cli-block.txt` **:1293**／**:1300**
  （本腿尺＝`grep -n -e REPLIED -e PUSHED-33R5-OK`，读数另件 `tmp-baseline-grep.txt`，rc=0、2 命中）

⇒ CI 与本机**同码同红**、且基线**同码绿** ⇒ "环境差"那一支在本票⛔ 依据（票面 `:13` 已就地作废那句旧话）。

## 5. 本格的件清单与 rc

| 件 | 内容 | rc 行 |
|---|---|---|
| `10-repro-head-count1.txt` | 终点乙（HEAD）`-count=1` 整发原文 | `rc=1`（文件末段自落） |
| `11-repro-head-count3.txt` | 复现率 `-count=3` 整发原文 | `rc=1` |
| `tmp-baseline-grep.txt` | 基线两枚绿句的尺输出 | rc=0（本件 §2 记） |
| `bisect-step.sh` | R4 台件（GOOD/BAD/INVALID 三分，⛔ 别的红当 bad） | 台件，⛔ 读数 |
| clone 侧原文（⛔ 入库，仓外） | `step001`…`step0NN`＋`index.tsv`＋`bisect-log.txt` | 逐行 tsv 带 go rc |
