# 300-r2 门禁八把（票 300 `AC#6`，落点 C2；每把都在自己的件里落一行 `rc=N`，⛔ 0 字节）

终锚读数在 `40-final.md`；本节的"改后"全部落在**同一枚台面**（这台机器的共享工作树），
"改前"那发取自我新增那枚文件**还没进树**的时刻（18:39，工作树＝锚 `05db4bc6` 干净态）——⇒ ⛔ 一发母仓、一发 clone。

| # | 门禁 | 尺（逐字可重跑） | 读数 | 件 | rc |
|---|---|---|---|---|---|
| 1 | `go vet ./internal/audio/` | `go vet ./internal/audio/` | stdout＋stderr **0 字节**；⚠ 0 字节与 `rc=0` 是两件事，⛔ 混写：这里是"输出为空**且** rc=0" | `logs/gate1-vet.txt` | **0** |
| 2 | 定向用例（pristine） | `go test ./internal/audio/ -count=1 -v -run 'TestWaveFormatConstantsMatchMmregAuthority300'` | 逐字绿：顶层 `--- PASS` **1**、子测试 `--- PASS` **6**、`=== RUN` **7**、FAIL/SKIP **0**、`ok ... 0.061s` | `logs/gate2-targeted-pristine.txt` | **0** |
| 3 | 整包两发取交集 | `go test ./internal/audio/ -count=1 -v`（改前 18:39／改后 18:4x，同一枚台面） | 名册尺＝`grep -E '^--- (FAIL\|SKIP): '` **剥时长后排序**（⚠ 顶层名那一把，⛔ 与含子测试那把混）；改前名册 1 行、改后名册 1 行，**逐字同一枚 `--- SKIP: TestLiveWasapiSmoke`**；`comm -13`＝**空**（新增红 0）、`comm -23`＝**空**（消失红 0）、`comm -3`＝**空** | `logs/gate3-full-package-before.txt`／`logs/gate3-full-package-after.txt` | 两发都 **0** |
| 4 | 格式两把并排 | 尺逐字见件；射程目录照写＝`internal/audio`；三把＝工作树／HEAD blob（改后）／锚 `05db4bc6` blob（改前） | 三处名册**各 0 枚** ⇒ 本腿新增脏行 **0**；⚠ 全仓 `gofmt -l cmd internal scripts tools` 有 **5 枚**残留（`cmd/wisp/models.go`、`internal/agent/approval/pending_read.go`、`internal/agent/tools.go`、`internal/risk/provenance.go`、`internal/tools/bridge.go`），**⛔ 一枚在本腿射程、⛔ 顺手修，具名上报** | `logs/gate4-format.txt` | **0**（三把都是空名册） |
| 5 | d22scan | `sh scripts/d22scan.sh` | `clean - no D22 ban violations`；ban #8 扫到 `internal/` **526** 枚 Go 文件、**comments and `_test.go` included** ⇒ 我这枚新用例在仪器射程内且**⛔ 命中**；⛔ 动 `tools/` 下任何文件 | `logs/gate5-d22scan.txt` | **0** |
| 6 | census totals 逐字未变 | `bash scripts/portable-tests.sh --scope=census` | totals 行逐字＝`packages=35 with-zero-compiled-tests=7 claimed-by-no-scope=7 unclaimed-with-tests=0`（⛔ 变一个字符）；⚠ **同族读数**：`internal/audio` 那一行是 `9/0  corewindows`——比改前**多 1 枚编进本平台二进制的测试文件**（本腿那枚），分母⛔ 动、认领档⛔ 变 | `logs/gate6-census.txt` | **0** |
| 7 | 名册差集 | **逐笔** `git show --name-only --format= <每笔自己的 commit>`（⛔ 区间尺 `git diff A..B`）⇒ 与"允许动的写面"名册作差 | 见 `40-final.md`；本腿三枚写面＝`internal/audio/wave_format_float_300_windows_test.go`／`.scratch/wisp/probes/300/r2/**`／票面追加一节 ⇒ 越界 **∅** | `logs/gate7-*`（在 `40-final.md` 里逐笔贴全） | — |
| 8 | 顶层用例枚数 | 逐文件 `git show <ref>:<file> \| grep -c '^func Test'`，射程目录＝`internal/audio`（⛔ 按文件名猜 tag） | 锚 `05db4bc6`＝**42**（8 枚文件）→ HEAD＝**43**（9 枚文件）；⚠ 具名口径＝**新增一枚顶层用例**（`TestWaveFormatConstantsMatchMmregAuthority300`）**＋它内部 6 枚子测试**，⛔ "只加子测试"；对拉 `300-a2 F5` 那把＝同一锚位重跑仍得 42 ⇒ 两把尺同源可复现 | `logs/gate8-test-counts.txt` | — |

## ⛔ SKIP 那一格要说清（这是全仓最容易被读反的一处）

- 本腿**没有**造任何 `--- SKIP`，也**没有**把 `t.Fatalf` 换成 `t.Skip`；`tools/d22scan/runtests.sh` 把任何 `^--- SKIP` 判红这条我照读。
- 改前那发（18:39，本腿的文件还⛔ 在树里）盘上**已经有**一枚 `--- SKIP: TestLiveWasapiSmoke` ⇒
  那是**仓／机既有**的一枚，且它是 `scripts/portable-tests.sh` 夹具台账里**已登记**的 fixture 例外（票 300 票面 `AC#5` 那节 12:5x 现读过 `:593` 那行）。
  ⇒ 两发名册里它**逐字同形**（`comm -3` 空）⇒ 它⛔ 是本腿的产物，也⛔ 因本腿而消失。**归因写死：既有、已登记、非本腿。**

## 我这枚用例进不进 CI 执行面（`AC#5` 那一族的事实层，⛔ 本格判语，只报读数）

- `portable-tests.sh --scope=census` 在**本机（Windows）**给 `internal/audio` 的计数从 `8/0` 变 `9/0` ⇒ 我这枚文件确实**编进了本平台测试二进制**。
- `301-r1`（`0a0f62ef`）之后 `./internal/audio/` 已在 `windows)` 档 ⇒ 新增那枚**顶层用例**名就是 `test-windows` 红名册作差面上多出来的那一枚：
  `TestWaveFormatConstantsMatchMmregAuthority300`（裁语 ④ 要具名的就是它）。
- ⚠ 本腿**没有**跑托管 CI、⛔ 拿本机读数冒充 CI 那一发色；`AC#5` 的"CI 到底求值过没有"仍欠（编排者 17:3x 已写：与票 303／111 那批同一枚推送一起取）。

## 本腿⛔ 碰的东西（越界自检，逐条）

产码（含 `internal/audio/wasapi_windows.go`、`wavinjector.go`）⛔／`parse_wave_format_300_windows_test.go` ⛔／`wavinjector_test.go` ⛔／
`scripts/**` ⛔／`.github/**` ⛔／`frontend/**` ⛔／`design/**` ⛔／`tools/**` ⛔／SLO 阈值·golden·`internal/observe/thresholds.go`·
`tools/d22scan/allowlist.txt`·D43 转移表·C1–C32·D1–D47 ⛔／`third_party/**` ⛔。
读页面源码一律对象层（`git show <ref>:<path>`）；⛔ 在仓内建 worktree／checkout；⛔ `add -A`／`amend`／`reset`／`rebase`／`stash`／`clean`／`--no-verify`；⛔ push。
`.scratch` 里⛔ 任何 `.go`（本腿只新建 `.md`／`.txt`），证据件⛔ 叫 `.out`（根 `.gitignore` 第 8 行是全仓 `*.out`）。
