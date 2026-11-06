# 98 — AC#1 事实层新读数 ＋ AC#4 收敛内容 ＋ AC#5 门禁读数（第二腿证据；非裁决表）

**票**：`.scratch/wisp/issues/98-cmd-wisp-tests-never-run-on-this-host.md`
**Status**：**实现证据，不是裁决表**。票面写死裁决表由验收方出；台账 `docs/reports/pending-and-issues.md` 归编排者（本腿一字未碰，AC#4 交付的是**收敛内容文本**，见 §6）。
**树**：派单锚 `abf2a35a6`；本腿读数时主树 HEAD `8879c3d7b1`（`A1567收wave30`）。主树裸读数与门禁 @`8879c3d7b1`；隔离普查 clone pin `8879c3d7b1fcf670d89766cde1fb777bb2423bb8`（`/d/tmp/98-leg/census-20261106-122916`，`git clone --shared`，起手 `git status --porcelain` 空＝0）。
**时间**：2026-11-06 12:15–13:0x +08（普查窗 12:33:00–12:52:31）。
**环境**：本机 Git Bash；`unset GOTMPDIR TMPDIR TEMP TMP`；`CGO_ENABLED` 保持默认 **1**（从未设 0）；PATH 注入只用派单给的 `$PWD/build:$PWD/third_party/sherpa-onnx`。
**写面**：**仅本文件**。`scripts/**` 零改动——AC#3 已由 `4a682f698` 落地（守卫在 `scripts/wisp-cli-tests.sh:222-315`），AC#2 的 documented 命令已住在同脚本头部与票面 09-28 追加，本腿没有要补的脚本缺口，不做无收益动作。票面（`.scratch/wisp/issues/98-*.md`）未动，Progress log 收敛由编排者执行。

## 0. 一句话

四件事：① **AC#1(i)(ii) 在当前树上重新实测**（裸命令红＝`0xc0000135`、0 条测试；加载器逐字点名 `sherpa-onnx-c-api.dll`；分母＝489 条可跑用例/156 枚测试文件）；② **AC#1(iii)** file:line 已在（`ci.yml:886/:906`），**真 run 步级结论＝E 尾，具名留档不冒充**；③ **AC#4** 收敛内容文本（§6，给编排者的笔）；④ **AC#5** 门禁全 rc 干净 ＋ 全仓**首次** `go test -count=2 ./cmd/wisp/` 全量普查：**PASS=874 FAIL=96 SKIP=8（count=2，每轮 437/48/4，489×2=978 逐行对账），rc=1（真红非加载失败），守卫行 0 次**（§8.1）。另有一个**不属于本票的意外发现**（§7）：**HEAD 纯净态编译不出 `go test ./cmd/wisp/`**——sweeper `abf2a35a6` 提交了引用、没提交定义，属于第 21 处的镜像形状，登记不追。

## 1. AC#1(i) — 裸命令本机原文（@`8879c3d7b1`，无 PATH 注入，实测）

```
$ go test -count=1 ./cmd/wisp/ ; echo rc=$?
exit status 0xc0000135
FAIL	github.com/CarlosShao/wisp/cmd/wisp	0.059s
FAIL
bare_rc=1
RUN=0 FAIL=0 PASS=0 SKIP=0
```

`=== RUN`/`--- PASS`/`--- FAIL`/`--- SKIP` **全 0** ⇒ 一条测试都没执行，正是票 98 钉死的"空仪器"形状。
加载器自己点名哪条 DLL（同一编译产物、从临时目录起、PATH 剥掉注入）：

```
$ cd /d/tmp/98-leg/rundir && PATH="/c/Windows/System32:/c/Windows" ../wisp-cli.test -test.list '.*'
list_temp_rc=127
D:/tmp/98-leg/wisp-cli.test: error while loading shared libraries: sherpa-onnx-c-api.dll: cannot open shared object file: No such file or directory
```

对照（机制＝`98-load-guard.md` §2.3 的"探针绿、门瞎"形，本树复测同读数）：同一二进制**从仓根直起**因仓根未跟踪野 DLL 得 rc=0——**直起能救、`go test` 救不了**（测试二进制从 TEMP 目录起）。

## 2. AC#1(ii) — 分母：它能跑的话会跑掉多少条（@`8879c3d7b1` 实测）

- 测试文件：`test_files_total=156`，其中 `_windows_test.go` 后缀 **70** 枚、含 `//go:build …winlive` **12** 枚（命令见 §10）。
- 可启动顶层用例：`go test -c` 产物 `-test.list '.*'` = **489** 条（上一棵树 `ff61b9dbc` 同法 482，差值来自波次落地的测试）。
- 结论一句话：**裸命令的红把 489 条全部静默归零**；注入 PATH 后同样这 489 条可启动（§4）。

## 3. AC#1(iii) — CI 那一格：file:line 已交，真 run 结论＝E 尾

`.github/workflows/ci.yml`（windows leg，`if: ${{ !cancelled() }}`）：

- **`:886`** `- name: "cmd/wisp CLI tests (needs the sherpa DLLs staged above, ticket 111 AC#4)"`
- **`:906`** `run: bash scripts/wisp-cli-tests.sh`

⇒ **接线真的存在**，且就是票 98 的守卫脚本本尊。**(iii) 的后半＝"一次真 run 的步级结论"＝E 尾**：本腿无 run id，按票面固定判据"说不出 run id 就当它不存在"，**不冒充**。留给编排者的动作：对最近一次 dev push 的 `test-windows` job 用 `gh` 取该 step 的步级结论＋run id，落台账与本票票面。

## 4. AC#2 — documented 命令与真实读数

Documented 命令（两形，都已是在册文档，非本腿新造）：

1. 一行式（票面 09-28 追加逐字）：`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test -count=1 ./cmd/wisp/`
2. 脚本式（守卫全开、CI 同款）：`bash scripts/wisp-cli-tests.sh`

真实读数：注入后 489 条用例可启动（§2）；在册 runtest 读数引前腿——票 111 树 `PASS=33 FAIL=0 SKIP=0`（脚本头 `:20-21`）、`ff61b9dbc` 树 `PASS=358 FAIL=41 SKIP=1, RUN=498, rc=1`（`98-load-guard.md` §2.2，红因在册属各票）。本腿补上**全仓首次 count=2 全量普查**（AC#5 口径）：`PASS=874 FAIL=96 SKIP=8`，rc=1，逐行对账 489×2（§8.1）。

## 5. AC#3 — 已落地（本腿零新代码）

`4a682f698`（A1551 在册）交付 PATH 可达性守卫＋加载失败兜底；本腿以普查日志做正控：**守卫行实测 0 次**（`grep -c 'ticket 98 AC#3'` = **0**）——19m31s 的真红跑里守卫没被误触发，也没把真红改色。

## 6. AC#4 — 收敛内容（给编排者的笔；本腿不写台账、不改票面）

**6.1 收敛后的那一句**（替换散在票 89/90/92/94/96/97 简报里的"cmd/wisp 的红不要追"，逐字可用）：

> `go test ./cmd/wisp/` 的裸命令红是**加载期缺 DLL**（`0xc0000135`，加载器点名 `sherpa-onnx-c-api.dll`），不是包坏；判定走票 98——documented 命令见票面 09-28 追加与 `scripts/wisp-cli-tests.sh`，永不 skip 的守卫已由 `4a682f698` 落地（AC#3）；ubuntu 侧红逐枚点名见 `docs/evidence/s1/98-attribution.md` §5，windows 侧红为在册现状、归各用例所在票。**"不要追"改为"点名＋归因"，不再整包不追。**

**6.2 更正链**（编排者已记的对齐项）：09-28 票面追加已驳"永远"标题（A361 记"CLI 只能走 -overlay"错句）；AC#3 已落地（A1551）。

**6.3 哪些 AC 的判据此前实际依赖过 cmd/wisp 的测试**（票 98 AC#4 第二半；重判归各票主人，本腿交点名依据）：

1. **票 133** 派单与票面的"CLI 只能走 `-overlay`"——判据是对 `0xc0000135` 的误诊，A361 已记账（错句已在两份派单＋三枚票面流传）；
2. **票 71 AC#5**（CI 色）——`docs/evidence/s1/triage-180.md:237` 已具名 E 尾巴；
3. **票 89/90/92/94/96/97 简报里引的 ubuntu 红计数**（19/29 时代读数）——全部应以 `98-attribution.md` §5 的 81 枚逐枚点名重判（其中 76 枚 C1＝DPAPI 桩、平台性质非包坏；`ae8136261` 已把 39 枚 C1 移绿，余 40 枚归票主）；
4. **windows 侧 41 红 @`ff61b9dbc`/本腿普查红**——在册现状，判据依赖方＝各用例所在票，不是"环境问题不追"。

**6.4 台账条目草稿**（A## 形，逐字给编排者）：

> 98（本腿 8879c3d7b1）：AC#1(i)(ii) 新读数＝裸命令 `0xc0000135`、0 用例、加载器点名 `sherpa-onnx-c-api.dll`、分母 489/156；AC#1(iii) file:line（ci.yml:886/:906）已交、真 run 步级结论挂 **E 尾**等 gh 窗口；AC#4 收敛句与依赖清单见 `docs/evidence/s1/98-ac1-ac5-readings.md` §6（台账更正由编排者执行）；AC#5 门禁 build/gofumpt/vet/d22scan 全净＋首次 count=2 普查（读数见同件 §8）。另登记：`abf2a35a6` 纯净态编译不出 `go test ./cmd/wisp/`（§7，归票 123/W123g）。

## 7. 意外发现（登记不追；第 21 处的镜像形状）

**HEAD（`abf2a35a6` 与 `8879c3d7b1`）纯净态编译不出 `go test ./cmd/wisp/`**：
已提交的 `cmd/wisp/run_mode101_test.go`/`run_test.go` 引用 `t123script` / `t123dirSpellingForms` / `t123replyScript`，而定义只住在**未跟踪**的 `cmd/wisp/ticket123_replyleg_{,other_,windows_}test.go`（`git ls-tree` 对两枚提交均无此三文件；隔离 clone 构建逐字 `cmd\wisp\run_mode101_test.go:95:10: undefined: t123replyScript` 等 8 行，`FAIL … [build failed]`）。主树能编译只因这三枚文件恰好在盘上；sweeper 提交信息"go vet编译通过"只在盘上有该文件的形态下成立。
**归属**：票 123（wave31 的 W123g 绿验正是抓这一形的尺）＋编排者；本腿零触碰 `cmd/wisp/**`。处置：普查 clone 按"HEAD＋五枚具名在飞测试文件"重跑（§8 读数按此口径注记），**不把编译缺口记到任何人的账上**，只交事实。

## 8. AC#5 门禁读数（本腿现跑，不用别人的数）

| 门 | 树 | 读数 |
|---|---|---|
| `go build ./cmd/... ./internal/... ./tools/...` | clone（tracked HEAD `8879c3d7b1`） | **rc=0** |
| `gofumpt -l cmd internal tools scripts` | clone（tracked HEAD） | **空**（rc=0，0 行输出） |
| `go vet ./cmd/wisp/` | clone（HEAD＋§7 五枚具名在飞文件） | **rc=0** |
| `go test -count=2 -timeout 45m -v ./cmd/wisp/` 全量普查 | clone（同上口径） | **PASS=874 FAIL=96 SKIP=8, RUN=1288**, pkg `FAIL … 1166.026s`, **rc=1**；对账 `489×2=978`＝874+96+8；唯一 FAIL 名 **52**（两轮皆红 **44**＝确定性红、单轮红 **8**＝窗口/焦点时序类），唯一 SKIP 名 **4**（全部具名非静默，见 §8.1） |
| `sh scripts/d22scan.sh` | 主树 @`8879c3d7b1` | **rc=0**，尾行 `clean - no D22 ban violations`（live scope: internal/=363, cmd/=58, frontend/=116, design/=39, internal/tools/=48） |
| 裸 `go test ./cmd/wisp/`（无注入） | 主树 @`8879c3d7b1` | rc=1 `0xc0000135`、0 用例——这是票 98 钉的形状本身，**不是门禁红** |

### 8.1 count=2 普查读数（本腿现跑，全仓首次全量 count=2；2026-11-06 12:33:00–12:52:31 +08，19m31s）

```
census_rc=1
PASS=874 FAIL=96 SKIP=8 RUN=1288
FAIL	github.com/CarlosShao/wisp/cmd/wisp	1166.026s
guard 行（'ticket 98 AC#3'）：0 次
```

- **对账**：每轮结果行 = 437 PASS + 48 FAIL + 4 SKIP = **489** ＝ §2 的 `-test.list` 分母，两轮 978 行分毫不差（count=2 无首轮/次轮分母漂移）。
- **确定性红（两轮皆红）44 枚**，占 FAIL 名单的绝大多数；**单轮红 8 枚**：`TestPanelThreadIsSTAAndExitsCleanly`、`TestBallPanelGesturesReachThePanelThread`、`TestAC4PriorFocusSurvivesARefusedPanelSample`、`TestPanelHostRealWindowHopAndLifecycle`、`TestAC4FocusReturnToPriorWindowGap33r5`（5/15s 级窗口/焦点时序类）、`TestAC14AwaitedBindingReplyReachesThePage`/`TestAC14GoSideEvalPushReachesThePage`、`TestAC1ResidentLegInstallsItsLogListenerOnDisk`、`Test253r6TheSilentDropBranchIsLoud`。
- **SKIP 4 枚（每轮同 4，全部具名非静默）**：`TestPanelPageFetchesItsOwnAssets299`（embed 无入口，票 33 r5 记录）、`TestAC13ColdStartEndsOnTheEmbeddedEntryNotTheProbe`（同上）、`TestAC247LiveMicrophoneLevelsReachTheBallSeam`（**真麦克风**，`WISP_LIVE_MIC=1` 才跑，`resident_audio_247_live_windows_test.go:130`）、`Test16r1PremortemMicPermissionDeniedLiveProbe`。无一枚是"加载失败被当 SKIP"。
- **归属**：红全部是**真红**（有 `--- FAIL`），不是票 98 的加载形；逐枚归票＝各用例所在票（同 `98-attribution.md` §5 的口径），本腿不替他们裁。44 枚确定性红与 `98-load-guard.md` §2.2 的 41 枚 @`ff61b9dbc` 同族（树间 ±，以本读数为准）。

## 9. 本件没有建立的事（宁缺毋造）

- **无 CI run id**（AC#1(iii) 后半＝E 尾，§3）；**未跑 gh**（等窗口，编排者面）。
- count=2 普查读数**只跟后提交**（先落盘时在飞，终值见 §8.1），绝无预填；-count=2 的单轮红（8 枚）已按两轮差分点名，未判 flaky 归因（那归各用例所在票）。
- 未做 ubuntu 重普查（`ae8136261` 后余 40 枚归属不主张）；windows 侧 FAIL 未逐枚归票（ubuntu 侧点名已档案化，windows 侧只给计数与名字清单）。
- 台账、票面、`scripts/**`、任何 Go 包：一字未动。

## 10. 复跑口令（可逐字重跑）

```bash
cd "/d/work/workspace/projects plans/Wisp"
unset GOTMPDIR TMPDIR TEMP TMP
# AC#1(i)
go test -count=1 ./cmd/wisp/ ; echo rc=$?          # rc=1, exit status 0xc0000135, 0 用例
# AC#1(ii) 分母
go test -c -o /d/tmp/wisp-cli.test ./cmd/wisp/     # rc=0（构建与 DLL 无关）
/d/tmp/wisp-cli.test -test.list '.*' | grep -cE '^(Test|Example|Benchmark)'   # 489（从仓根起，野 DLL 救直起）
ls cmd/wisp/*_test.go | wc -l                      # 156
ls cmd/wisp/*_windows_test.go | wc -l              # 70
grep -l 'go:build.*winlive' cmd/wisp/*_test.go | wc -l                        # 12
# AC#1(i) 点名 DLL
mkdir -p /d/tmp/rundir && cd /d/tmp/rundir && PATH="/c/Windows/System32:/c/Windows" /d/tmp/wisp-cli.test -test.list '.*'   # rc=127 点名 sherpa-onnx-c-api.dll
# AC#2 documented 命令
cd "/d/work/workspace/projects plans/Wisp" && export PATH="$PWD/build:$PWD/third_party/sherpa-onnx:$PATH" && go test -count=1 ./cmd/wisp/
bash scripts/wisp-cli-tests.sh
# AC#5 普查（隔离 clone，HEAD＋五枚具名在飞文件）
git clone --shared "/d/work/workspace/projects plans/Wisp" /d/tmp/98-census && cd /d/tmp/98-census
mkdir -p third_party/sherpa-onnx && cp "/d/work/workspace/projects plans/Wisp/third_party/sherpa-onnx/"*.dll third_party/sherpa-onnx/
cp -r "/d/work/workspace/projects plans/Wisp/build" build
cp "/d/work/workspace/projects plans/Wisp/cmd/wisp/"{ticket123_replyleg_test.go,ticket123_replyleg_other_test.go,ticket123_replyleg_windows_test.go,panel_resident_323_windows_test.go,panel_resident_assembly_323_test.go} cmd/wisp/
export PATH="$PWD/build:$PWD/third_party/sherpa-onnx:$PATH" && go test -count=2 -timeout 45m -v ./cmd/wisp/
```
