# 300-a2 · 30 — CI 可见性：这一格落在 CI 上⛔ 落（⛔ 行号引调用点，一律内容锚）

**起手锚＝`4e00357f`（2026-10-10T17:34:47+08:00）。** 全部读自 blob 层；逐字 stdout 存 `logs/ci-visibility.txt`／`logs/ci-jobs.txt`／`logs/ci-guards.txt`。

## 尺（两把，交叉）

- 尺 A＝`git show HEAD:scripts/portable-tests.sh` 取 `case` 标签名册、各档清单、四枚 pin、夹具 ledger —— 内容锚：`core)`／`windows)`／`cli)`／`census)`／`winsec)`、`core_pin='`、`win_pin='`、`name|package|platform|class|reason`
- 尺 B＝`git show HEAD:.github/workflows/ci.yml` 取 job 名册（`^    <name>:$`）、`runs-on:`、调用点 —— 内容锚逐字：`run: bash scripts/portable-tests.sh --scope=core`、`… --scope=census`、`… --scope=windows`、`run: bash scripts/wisp-cli-tests.sh`（★⛔ 用行号引调用点）
- 尺 C＝调用点那一步**有没有零求值 `if:`**：对每个调用点做 `grep -n -B8 <锚>` 取它自己的 step 级 `if:` 行
- 尺 D＝`git grep -n -i -e 'internal/audio' HEAD -- .scratch/wisp/issues/301*` 取票 301 的自述射程（只为对拉"audio 进 windows 档"这件事的**出处**，⛔ 当凭据）

## 名册（逐字读到的）

**尺 A（portable-tests.sh，blob）**
- 档标签：`core)`（blob 内 `case` 分支，内容锚命中）／`windows)`／`cli)`／`census)`／`winsec)`
- `windows)` 档清单里**确有** `./internal/audio/`（该档清单末行逐字含 `./internal/session/ ./internal/projctx/ ./internal/audio/`）⇒ **票 301 `AC#1` 那两改（清单＋`win_pin` 同笔）已落在 HEAD 上**。
- `win_pin='` 名册里**确有** `github.com/CarlosShao/wisp/internal/audio`（`core_pin` 里也有它）⇒ GUARD C 的 pin 侧已配套。
- `cli)` 档只有 `./cmd/wisp/`，且注释逐字：`# cmd/wisp's own tests. Only scripts/wisp-cli-tests.sh may call this, because`。
- 夹具 ledger 逐字一行：`TestLiveWasapiSmoke|./internal/audio/|windows|fixture|live WASAPI capture needs WISP_LIVE_MIC=1 and a physical microphone (hotplug_test.go:527, //go:build windows); a hosted runner has no audio endpoint, so the case has no subject there.`（⚠ 带一行号锚，见 20 件 N3）
- ⚠ 该脚本**⛔ 任何"某包应有 N 枚用例"的计数钉**（尺＝逐名册 grep `audio`／`expect`／`denominator`，命中的分母全是"import path 名册"与"本平台编进二进制的 `.TestGoFiles`+`.XTestGoFiles`"两件事）⇒ **新增用例⛔ 打红分母类仪器**。

**尺 B（ci.yml，blob）**
- job 名册：`lint:`／`test-core:`／`test-windows:`／`slo-smoke:`／`slo-full:`／`lint-frontend:`（`^    <name>:$` 那把尺数到 6 枚 job ＋ `on:`/`concurrency:`/`jobs:` 三枚顶层键）
- `runs-on:` 逐枚：`lint`＝`ubuntu-latest`；`test-core`＝`ubuntu-latest`；`test-windows`＝`windows-latest`；`slo-smoke`＝`windows-latest`；`slo-full`＝`[self-hosted, wisp-slo]`；`lint-frontend`＝`ubuntu-latest`
- 调用点归属（内容锚 ⇒ 档）：`--scope=core` 在 `test-core`（ubuntu）；`wisp-cli-tests.sh`、`--scope=census`、`--scope=windows` 三枚都在 `test-windows`（`runs-on: windows-latest`）；`slo-smoke`／`slo-full` 两枚 job 在本把尺的调用点名册里**零枚 `--scope=` 命中** ⇒ 它们⛔ 跑 portable 档（⛔ 说"它们测 audio"，⛔ 量过它们的 run 步骤全文）。

**尺 C（step 级 `if:`）**
- 四处调用点**各自**的紧邻上行逐字都是 `if: ${{ !cancelled() }}` ⇒ 这四处**有求值**（⛔ 落进"零求值步级 `if:`"那一族）。
- ⚠ 口径限定：全文件 `if:` 计数＝**56 枚**（尺＝`git show HEAD:.github/workflows/ci.yml | grep -c -e 'if:'`）；本腿**⛔ 复跑**台账里那句"零求值 26 枚"，⛔ 用它做任何判断。

## 答复（⛔ 裁形、⛔ 推荐）

1. **这一格落在 CI 上。** 依据＝F1（`waveFormatFloat` 只在 `//go:build windows` 文件里 ⇒ 任何钉它的判据⛔ 得 tagged）× 尺 A（`internal/audio` 已在 `windows)` 档清单**且** `win_pin` 已配套）× 尺 B（`--scope=windows` 调用点在 `test-windows`，`runs-on: windows-latest`）× 尺 C（该 step `if: ${{ !cancelled() }}` ⇒ 真求值）。
2. **求值它的档与 OS＝`test-windows` 那一档、`windows-latest`（GitHub 托管镜像）**；同包里那枚 tagged 用例的**唯一另一条 CI 路径**是 `--scope=census`（同 job），而它判的是**包级**认领（`go list -f '{{len .TestGoFiles}}/{{len .XTestGoFiles}}'` 那把探针就是"编进本平台二进制"的枚数）⇒ census⛔ 执行用例，只数名册。
3. **ubuntu 那档（`test-core`）永远⛔ 求值它**——包在 `core)` 清单里但 tagged 文件编不进 linux 二进制；这不是推断，是脚本自己写着的话（GUARD A 那两行注释逐字：`.TestGoFiles`+`.XTestGoFiles` are the files that COMPILE INTO THIS PLATFORM's test binary）。⇒ **推论（只有尺，⛔ 我裁）：硬约束③要的"换形必红"与"两台 OS 都看得见"二者⛔ 能同时买到**；要 ubuntu 也响，唯一形是把判据搬到⛔ 引用该常量的东西上（C4 的 injector 那支是 linux 编得进的，但它⛔ 能因该常量换形而红——见 10 件 C4）。
4. **缺哪几发读数（具名，⛔ 静默；全部⛔ 本腿可取）**：
   - **① CI 色本身**：`test-windows` 那一档今天跑起来，新判据是 PASS／FAIL／SKIP —— 只有 push 之后 CI 给（本腿⛔ push；票面 `300-v3` 那节的 next 已写"跟票 303 修复后那批推送一起取，⛔ 为它单推"）。**归编排者。**
   - **② 托管 windows 镜像上那枚 SDK 头文件在不在**（C3 的生死判）：尺＝`git grep -n -i -e 'Windows Kits' -e 'mmreg.h' HEAD -- scripts tools .github` ⇒ **仓内 0 命中**（⛔ 准备动作、⛔ 读取先例）⇒ 盘上事实**只有 CI 那台机器能答**。**归编排者（那一发 ① 顺带答）。**
   - **③ `internal/audio` 在 `windows)` 档的当前基线色**：票 301 `AC#2` 自称要交"改清单前后各一发、名册逐名作差"，而本格⛔ 由 301 的腿交回与否**本腿⛔ 查**（`probes/301/**` 在我的授权名册之外，只按票面引用）。⚠ 引用票 301 任何"已落地"字样前先现跑 `git show HEAD:scripts/portable-tests.sh` 那两枚名册（本腿已跑，见上）。**归编排者／301 的腿。**
   - **④ `TestLiveWasapiSmoke` 在托管 runner 上到底是登记过的 SKIP 还是红**：ledger 那行的理由句自称"a hosted runner has no audio endpoint"，⛔ 实测凭据在盘上。**归那一发 ①。**
   - **⑤ 56 枚 `if:` 里到底哪几枚是零求值**：本腿只把四处调用点各自的 `if:` 读到（都是 `!cancelled()`），⛔ 复跑"零求值 N 枚"那句。**归编排者（那句原写在台账里）。**
5. **一处与本票派单措辞的对拉（⛔ 改票面一字）**：派单说"★调用点⛔ 能用行号引"、"`internal/audio` 自 `0a0f62ef` 起已在 `windows)` 档里"——我这把尺在 **blob 名册层面**证实了后半句（内容锚命中），但**⛔ 能证实它是 `0a0f62ef` 那笔带进来的**（本腿没跑 blame／bisect 面）；引用那笔号时请现跑 `git show --name-only 0a0f62ef`。⛔ 矛盾，只是我的尺⛔ 到那一层。

rc=0
