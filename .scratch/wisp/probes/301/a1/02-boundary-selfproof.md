# 301-a1 — 边界自证 ＋ 锚点在飞行中移动的记录

## ① ★HEAD 在我这一程中途动了（不是被动因，具名记录）

- 起手锚＝`772ff880`（`00-anchor.md` 逐字回显）。
- 我的第 1 笔 commit `c2853233` 落下之后，编排者自己那笔 **`6d81638b`「立票 301 … ＋A806／§4.0bg＋派 301-a1 只读」** 落在了它上面 ⇒ 现 HEAD＝`6d81638b`。
- 尺＝`git diff --name-only 772ff880..HEAD` ⇒ 4 枚文件：票面 301 `.md`、`docs/reports/HANDOVER.md`、`docs/reports/pending-and-issues.md`、**我的** `.scratch/wisp/probes/301/a1/00-anchor.md`（第 4 枚是它的父 commit `c2853233` 带的）。
- ⇒ **我量过的每一枚 blob 在两个锚下逐字相同**（尺＝`git diff --stat 772ff880..HEAD -- internal/audio`／`-- scripts`／`-- tools/d22scan`／`-- internal/proc`／`-- internal/memory` ⇒ **五把尺全回显空**）。
- ⇒ 表① 在新 HEAD 上重跑得同一读数（尺＝`git ls-tree -r --name-only HEAD -- internal/audio | grep -c '_test\.go$'` ⇒ **7**；tagged ⇒ `capturelevel_windows_test.go` `^func Test`=**1**、`hotplug_test.go` `^func Test`=**1**+**8**⇒ **9**）。
- ⚠ **`6d81638b` 没吞掉任何在飞面**：`git status --porcelain -- internal` 现仍回显那 2 枚（` M internal/audio/wasapi_windows.go`、`?? internal/audio/parse_wave_format_300_windows_test.go`），票 300 落地腿的未提交件**躺在原处、⛔ 被顺走**。

`rc=0`

## ② 只读自证

- 写面＝**只新建** `.scratch/wisp/probes/301/a1/` 下的 `.md`（本程共 2 枚：`00-anchor.md`、`01-census-ac0-cost-tables.md`、外加本件）。
- ⛔ 任何被跟踪文件被我改过：`git status --porcelain -- internal cmd scripts .github docs .scratch/wisp/issues` ⇒ **回显 2 枚，全是起手那两枚在飞面，⛔ 一枚出自我手**（尺与 §1 末条同）。
- ⛔ 任何 `.out` 命名（根 `.gitignore` 第 8 行全仓 `*.out`）。
- ⛔ 0 字节件：尺＝`find .scratch/wisp/probes/301/a1 -type f -printf '%s %p\n'` ⇒ `00-anchor.md` 1179 B、`01-census-ac0-cost-tables.md` 23044 B（本件落盘后另量）。
- ⛔ 新建 `.sh`／`.ps1`（那会挪 gofumpt/d22scan 的分母）。
- ⛔ 任何 go 编译面：全程**零次** `go build`／`go vet`／`go test`／`go run`；只用过 `git ls-tree`／`git show`／`git diff`／`git log`／`git status`／`grep`／`sed`／`cat -n`／`tr`／`sort`／`diff`／`find`／`wc`／`cut`／`mkdir`／`ls`。`go env` 也没用上（无必要）。
- ⛔ 碰 `frontend/**`／`design/**`／`%APPDATA%\wisp-dev\secrets\`；件内**只引 env 变量名（`WISP_LIVE_MIC`、`WISP_CRASH_CHILD`、`WISP_84_MEASURE`、`WISP_IT_REAL_MIRROR`），⛔ 任何值／凭据**。
- ⛔ 翻任何 `- [ ]` 框（票 301 票面 `git diff` 里我没碰它，见 §1 的 name-only 尺：票面那枚文件的改动出自编排者那笔）。
- ⛔ `git add -A`／`git add .`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`／`--no-verify`；两笔 commit 各带**显式 pathspec**。
- ⛔ **push**（机主从未授权；本程零远端写）。

`rc=0`

## ③ 本程 commit 名册（供 `git show --stat` 逐笔自证）

| commit | 名册 | 只含 `probes/301/**`？ |
|---|---|---|
| `c2853233` | `.scratch/wisp/probes/301/a1/00-anchor.md` ＋1 枚，35 行 | ✅ 是（尺＝`git show --stat --format=%H c2853233`，逐字回显 1 file changed） |
| （第 2 笔＝本程末笔，见回报） | `.scratch/wisp/probes/301/a1/01-census-ac0-cost-tables.md`、`02-boundary-selfproof.md` | ✅ 显式 pathspec 逐枚指定 |

## ④ 票面/派单里被我⛔ 采信的一处前提（复述，细节见 `01` 件末格）

票面 `现量` 与 `:15` 那句"⛔ 零成本"的**枚数**（"上面那 2 枚"）在 blob 上读作 **2 枚文件＝9 枚用例**；
我按派单 §2 表① 的口径（"顶层用例名"）交 9，并具名报回那处口径错位。
表③ 那行 ledger 的判定也⛔ 是"整行 inert"——`-skip` 腿打空、`-list` 陈旧性腿今天真在 windows job 上求值；
把整行判死会在 `AC#1` 里诱发"顺手删 ledger 行"，而删掉之后第 8 枚会以 `--- SKIP` 命中 `runtests.sh:98-102` 把 `test-windows` 判红。
