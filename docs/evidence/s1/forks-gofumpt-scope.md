# 票 410 r1 — `forks/go-webview2` gofumpt 射程裁定（甲＝格式化 fork，非乙＝exclude）

- **日期**：2026-11-08 23:3x +08
- **触发**：276 腿现跑 `attrib.sh --tracked-only` ⇒ RED rc=1，红点＝`forks/go-webview2` 8 枚（`2bb639e4e` 11-07 票65 CarlosShao，已推 origin/dev）。A1556 记录的「门今天 green」已再变——这是新事实。
- **裁定**：**甲＝把 fork 件按 gofumpt 格式化**；**不采乙**（exclude `forks/**`）；**不采丙**（单独「不格式化但可编译」档）。
- **implementer**：Wisp 写腿（本波写面＝fork 8 件＋本证据件＋票 410；`attrib.sh` **一字未动**）。

## 1. 起手三核（现跑，非转述）

1. **红名单**：`bash .scratch/wisp/probes/161/r5/attrib.sh --tracked-only` ⇒ RED rc=1，
   `TRACKED-DIRTY` 恰 8 行、全部在 `forks/go-webview2`、last commit 全＝`2bb639e4e 11-07 16:44 CarlosShao`
   （`common.go`、`pkg/edge/chromium_386.go`、`comproc_go17.go`、`comproc_go18.go`、`corewebview2.go`、
   `guid.go`、`ICoreWebView2AcceleratorKeyPressedEventHandler.go`、`webviewloader/module.go`）。与派单点名一致。
2. **exclude 现状**（`attrib.sh:560-600` 段现读）：只有一个具名排除 `A_N_EXCLUDE=':(exclude).scratch/**'`，
   带三重防滥用设计（分母三数打印／`--self-test` 名册控件种在 `.scratch` 外／362-r2 出树实测）。
   该排除的立票理由是「`.scratch` 里是**故意坏的台件**」——与 fork 的情形（真产码、被编译、被 import）**不同轴**。
3. **fork 同步状态**（`git log --oneline -- forks/go-webview2/`）：仅两笔——
   `2bb639e4e`（票65：jchv 上游整模块直搬＋本仓补 `Close()/Release()` 包装＋`go.mod replace`）
   与 `da5d559b7`（票408 r3：**在本仓作者权下**改 `pkg/edge/chromium.go` 落 nil-env 守卫，该件今天 gofumpt 干净）。
   **fork 不是「只直搬」**：已有本仓写的码。

## 2. 为什么是甲、不是乙（推翻派单的乙倾向，判据逐条）

**乙的判据是「那些文件本仓没有作者权」。这条今天事实不成立：**

- 票65（`docs/evidence/s1/65-jchv-controller-leak.md`）已在 fork 里**写了** `ICoreWebView2Controller.Close()/Release()` 包装方法——本仓作者权第一笔。
- 票408 r3（`da5d559b7`）**在本仓作者权下**改了 `pkg/edge/chromium.go`（nil-env 守卫），且按 gofumpt 干净标准落——本仓作者权第二笔，同时证明既成规范就是「fork 内的手写码守 gofumpt」。
- ⇒ exclude `forks/**` 不是「划清射程」，而是**把本仓正在持续写码的目录摘出格式门**，并削弱 attrib.sh 自述的「CI 五目录白名单 vs 全树尺」分歧发现功能。

**甲的「同步冲突」代价是伪命题：**

- `go.mod:24`＝`replace github.com/jchv/go-webview2 => ./forks/go-webview2`，pseudo-version 钉死（`v0.0.0-20260205173254-56598839c808`）；
- 票65 证据件逐字：「已钉 pseudo-version，无自动漂移风险，**上游更新＝手动同步一次**」——同步本就是手动事件，且**语义层已分叉**（票408 守卫），手动同步本来就要人眼对 hunk；
- 8 枚格式 hunk（见 §3，共 +16/−12：interface 内空行、func 间空行、import 分组、`//go:uintptrescapes` 挪到 doc comment 后）只给那一次手动同步添一点噪音，无语义风险。

**先例与生态**：票70 AC#1 对同族问题（gofumpt 红 vs 门）的裁定就是「纯格式化单独 commit，**放宽门不是修门**」。
Go 生态免格式门的只有 `vendor/`（go 工具链自管理、规范上不许手改）；本仓**无 `vendor/`**（台账 `Q-84` 现量：`test -d vendor`→NO），
`forks/` 是手改硬叉，不属豁免族。丙档（可编译不格式化）对 8 枚空行级差异是复杂买不到价值。

## 3. 落地与四数

- 改动：`gofumpt.exe -w forks/go-webview2/` ⇒ **恰 8 件，+16/−12，纯格式零语义**（commit `504858fd1`）。
- **验收（`git clone --shared` 导出树 `/d/work/tmp/wisp-exp-410`，HEAD=`504858fd1`）**：
  `bash .scratch/wisp/probes/161/r5/attrib.sh --tracked-only` ⇒ **GREEN rc=0**。
- **四数**（attrib.sh 终行逐字）：**(A) tracked-tree lines=0 files=0**；**handed to gofumpt = 1286**；
  **total tracked .go = 1657**；**subtracted by the one named exclusion = 371**。

## 4. 门禁（同一导出树，`unset GOTMPDIR TMPDIR TEMP TMP`＋`PATH=$PWD/build:$PWD/third_party/sherpa-onnx:$PATH`，未设 CGO_ENABLED=0）

| 门 | 读数 |
|---|---|
| gofumpt（改动件 `forks/go-webview2/`） | `-l` 输出 **0 行** |
| `go vet ./...` | **rc=0** |
| `go build ./...` | **rc=0** |
| `scripts/d22scan.sh` | **clean**（rc=0，ban #8 覆盖 design 14／frontend 121／internal 966／cmd 242） |
| LF | 8 件 `git ls-files --eol` 全 **i/lf w/lf** |

## 5. 防作弊自查

- 本裁定**没有放宽任何门**：`attrib.sh` 未改一字、CI 未改、thresholds/golden 未碰；红消失的原因是**码就范**（与票70 AC#1 同形状），不是尺让步。
- 与 362 的 `.scratch` 排除不同轴：那是「故意坏的台件永远红」，这是「真产码格式化后真干净」。
- 撤销口令：`git revert 504858fd1` 一枚即回到 8 枚红的现状。
