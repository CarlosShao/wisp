# 135 AC#1–AC#6 实现侧台账（在 `45371ca0b` 上重量的读数）

**本文件是什么**：票 135 的实现侧（landing leg）读数件，锚 `45371ca0b3ade1136f7c89950e4405c7a4fb3d3d`（2026-11-05 20:1x–20:4x +08）。
**本文件不是什么**：不是裁决表。`SPEC-12 §4.3` #1/#3 与 `issues/README` 末条要求裁决表出自**非实现者**；
本文件里的每一格都是实现方自己的读数，**一格都不翻框**、不动 `Status:`、不动票名。

## §0 这一程为什么写这枚件（先摆盘上的事实）

票面 `## Progress log` 末条（2026-10-15 16:2x，`agent=135-r1`）写的是「next=剩 AC#6 的整包 cmd/wisp 四数、
以及非实现者的终裁」，但盘上实际落得比那条 log 多：`135-r1` 之后又落了两笔（`8f2f863ec` AC#6 门禁、`9b0f7d555` AC#6 整包读数更正）。

本程第 0 步先量「AC#1–AC#5、AC#7 在今天的树上还在不在」：

| AC | 落在盘上的东西（逐枚现读） | 引入 commit |
|---|---|---|
| AC#1 | `cmd/wisp/dataroot_128_test.go:506` `TestAC1TheStartUpDirectoryCheckHasItsOwnTeeth128`＋`:426` `redSink128`／`:438` `recorder128` | `b46a867b1` |
| AC#2 | 同文件 `:600` `TestAC2ARefusalWithATreeInTheStartUpDirectoryIsRed128` | `b46a867b1` |
| AC#3 | 同文件 `:545` `TestAC3TheDoctorLegLeavesTheStartUpDirectoryUntouched128` | `b53b79bda` |
| AC#4 | 同文件 `:238–270`（`TestAC2RefusalMarkersAreNotAShortenableList128` 里读生产 `err.Error()` 的三条成句断言） | `144b038f5` |
| AC#5 | `internal/proc/envfork_test.go:351` `TestDefaultLayoutRefusesWithNoUserConfigDir` | `bf25809a3` |
| AC#7 | `cmd/wisp/leg_sink_gate_131_test.go:595` `TestAC7ARulingParkedInAFileBacksNothing131`＋`:466` `rulingAnchors131` | `2781ec380` |
| AC#8 | 已由非实现者终裁翻勾（`docs/evidence/s1/135-ac8-r1-acceptance.md`） | 七笔 `a025892…29deecb` |

⇒ **本程没有改动任何生产码或测试码**（`git diff --name-only` 全程 0 行，见 §6）。本程做的是三件事：
① 在今天的树上把上述用例跑绿；② 逐格用**真变异**重证「它坏了会响」（§4）；③ 把 AC#6 那笔门禁账补成可复算的读数（§5）。
⚠ 其中 ② 现量出一处**票面存证与今天树上事实不符**的更正，见 §4.3。

## §1 台面

- `go version go1.27.1 windows/amd64`，`GOOS=windows`、`GOARCH=amd64`、`CGO_ENABLED=1`（未改）。
- 所有 `cmd/wisp` 读数都在 `PATH="$PWD/build:$PWD/third_party/sherpa-onnx:$PATH"` 下取（票 98：不注入就是加载期 `0xc0000135`，那是"用例根本没跑"）。
- `unset GOTMPDIR TMPDIR TEMP TMP`（票 310/311 守卫）。
- 树＝`45371ca0b` 的干净工作树：`git diff --name-only` = 0 行（起手与末了各量一次）。
- 跑前离场读数：`Get-CimInstance Win32_Processor` 平均负载 37%，`Get-Process go` = 0（无别的腿在编队里跑）。

## §2 起手绿（本票用例，今天逐枚在）

```
go test ./cmd/wisp/ -count=1 -timeout 420s -v -run '^(TestAC1TheStartUpDirectoryCheckHasItsOwnTeeth128|TestAC2ARefusalWithATreeInTheStartUpDirectoryIsRed128|TestAC3TheDoctorLegLeavesTheStartUpDirectoryUntouched128|TestAC2RefusalMarkersAreNotAShortenableList128|TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128|TestAC5RefusalSentenceHasOneWordingPoint128|TestAC7ARulingParkedInAFileBacksNothing131)$'
```
读数：七枚 `--- PASS`、`ok github.com/CarlosShao/wisp/cmd/wisp 1.253s`、rc=0（日志 `00-base-targeted.log`）。
`go test ./internal/proc/ -run '^TestDefaultLayoutRefusesWithNoUserConfigDir$' -count=1 -v`：`--- PASS`，`ok … 0.070s`，rc=0。

## §3 每一格的红名口径（票面原文）

| AC | 票面要求的红名 |
|---|---|
| AC#1 | 弄哑 `assertDirEmpty128`（或喂不存在的目录）⇒ **它必须响**；摘掉这枚新用例 ⇒ 其余不会全绿 |
| AC#2 | 拒绝路径上先长一棵 CWD 树 ⇒ 必须红，**且红名点到"目录非空"而不是只点到 rc** |
| AC#3 | `doctor` 腿要有**独立**的 CWD-empty 断言（marker 一撤仍有牙） |
| AC#4 | 自证用例**真读生产文案**，三要素做成**成句断言**（六枚子串全在但不成句必须红） |
| AC#5 | 常驻腿 POSIX 有分母（`internal/proc` 本包用例） |
| AC#6 | 门禁：四数（只能从 `-v` 量）＋`gofmt`/`gofumpt`＋`go vet` 双 GOOS＋`d22scan` rc=0 且台账各 scope 不降 |

## §4 变异（RED 在变异下、GREEN 在还原后）—— 逐发 `grep` 证落地 + `go build` rc=0 再读红名

### §4.1 M-C：把"目录空"那一句弄哑（AC#1 / AC#2 的正反面）

- 变异落地：`cmd/wisp/dataroot_128_test.go:493` `if len(found) > 0 {` ⇒ `if false && len(found) > 0 { // MUTATION M-C 135-r2`（`grep -n` 现命中 1 处；`go build ./cmd/wisp/` rc=0）。
- 读数：**7 用例里 2 红**——`--- FAIL: TestAC1TheStartUpDirectoryCheckHasItsOwnTeeth128`、`--- FAIL: TestAC2ARefusalWithATreeInTheStartUpDirectoryIsRed128`（6 个子腿全红），其余 5 绿，rc=1，`1.167s`。
- 红名逐字（AC#1）：`dataroot_128_test.go:512: AC#1 RED (the instrument, not the code): the start-up-directory check said nothing about a tree that is provably there (…)`。
- 红名逐字（AC#2）：`dataroot_128_test.go:621: AC#2 RED: leg "runTextTask" refused correctly with a wisp-dev tree sitting in the start-up directory and the emptiness check said nothing. That is M-F exactly, and it is the shape that read green once this check was muted.`（六腿各一枚）
- **AC#1 自己的结案判据**（"摘掉这枚新用例 ⇒ 其余不会全绿"）：同变异下把 AC#1 那枚从 `-run` 里摘掉 ⇒ 仍 `--- FAIL: TestAC2ARefusalWithATreeInTheStartUpDirectoryIsRed128`、其余 5 绿、rc=1、`0.605s` ⇒ 牙有两名看守，不是单点。

### §4.2 M-A2（**照票面存证那一形复现，今天它是死码**）

- 存证 `02-ac3-doctor-tooth.txt` 写的形是「拒绝句一字不动 ＋ 仍 `probeWritable(".")`」。按它把 `cmd/wisp/doctor.go:107` 的 `probeWritable(dir)` 改成 `probeWritable(".")`：
  落地、`go build` rc=0，但 **7 枚用例全绿 rc=0**。
- **为什么**（现读）：`doctor.go:105` 先判 `dirErr != nil`，票 128 之后 `resolveDataDir` 在 OS 读失败时**一律拒绝**，所以 `else if probeWritable(...)` 这一支**今天根本走不到** ⇒ 该变异是死码，它测不到任何东西。
  ⚠ 这不推翻 AC#3 的用例，而是推翻**那一发变异**的可复算性：照抄它读到的"修后唯一红"，今天复现不出来。

### §4.3 M-A2b：忠实的那一形（AC#3 的牙）

- 变异落地：`doctor.go:104–115`（拒绝句与 `fail(...)` 一字不动，**另起**一段 `probeWritable(".")` 并 `_ = dir`）——即"拒绝了、但可写性探针仍在启动目录跑"。
- 读数：**只有一枚红** `--- FAIL: TestAC3TheDoctorLegLeavesTheStartUpDirectoryUntouched128 (0.08s)`，其余 6 绿，rc=1。
- 红名逐字：`dataroot_128_test.go:577: AC#3 RED: the doctor leg resolved the start-up directory as its data root and its writability probe removed the sentinel (…\doctor-write-probe.tmp), so the leg fell back to "." instead of refusing. This is the shape the emptiness check cannot see - the probe deletes its own file.`
- ⇒ AC#3 的牙**活着**，且是 marker 之外的一枚独立读数（这一发下 marker 全过、AC#3 唯一红）。票面 `R-128-6` 的更正口径照旧成立。

### §4.4 M-E：六枚子串全在、但不成句（AC#4 / AC#5 的牙）

- 变异落地：`internal/proc/envfork.go:273`（那一句 `fmt.Errorf` 换成 `"%w（OS 原话：%v）当前工作目录 APPDATA XDG_CONFIG_HOME 可写目录 数据根本应是 %s"`）——六枚 marker、`<用户配置目录>\wisp-dev`、OS 原话**全在**，但三条成句、次序与句末标点全没了；`go build ./cmd/wisp/ ./internal/proc/` rc=0。
- 读数（`cmd/wisp`）：**只有 AC#4 那枚红** `--- FAIL: TestAC2RefusalMarkersAreNotAShortenableList128`，其余 6 绿，rc=1；红名逐字含 `:260: AC#4 RED: …the clause "数据根无法解析：用户配置目录不可得（OS 原话：…）：" is not in it` 三条＋`:269: AC#4 RED: …does not end with the fix and its full stop`。
  ⇒ 六枚子串的合取检查在 M-E 上**仍全过**，新牙与旧清单**独立**（这正是 `R-128-2` 要的那一层）。
- 读数（`internal/proc`）：`--- FAIL: TestDefaultLayoutRefusesWithNoUserConfigDir`，红名 `envfork_test.go:407/416` 同款成句断言两条＋句末标点一条，rc=1。

### §4.5 AC#7 的牙

本程**未重跑** AC#7 的两拍 (a)/(b)。原因写在这里而不是留白：它的牙是**用例自带**的 fixture 腿（`TestAC7ARulingParkedInAFileBacksNothing131` 把裸标记与诚实锚两形喂真扫描），起手绿已在 §2；票面的 (a)/(b) 需要另改 `ruled` 谓词＋在 `slo_other.go` 末尾追加裸句，属**下一程可补的读数**，本程没做，不写成"已复核"。

## §5 AC#6 门禁（本程现跑；数字各尺独立 `grep -c`，不相减）

| 尺 | 命令 | 读数 |
|---|---|---|
| build | `go build ./cmd/... ./internal/... ./tools/...` | **rc=0**（零输出） |
| 本票切片四数（`-count=2 -v`） | `go test ./cmd/wisp/ -count=2 -timeout 420s -v -run '<本票十枚>'` | **RUN=56 / PASS=20 / FAIL=0 / SKIP=0，rc=0，13.110s** |
| `internal/proc` 整包四数（`-count=2 -v`） | `go test ./internal/proc/ -count=2 -timeout 420s -v` | **RUN=106 / PASS=88 / FAIL=0 / SKIP=2，rc=0，3.529s**（两枚 SKIP ＝ `TestHelperProcess` 帮助进程形） |
| `cmd/wisp` 整包四数（`-count=1 -v`） | `go test ./cmd/wisp/ -count=1 -timeout 1500s -v` | **RUN=640 / PASS=439 / FAIL=44 / SKIP=2，rc=1，662.377s**（见下"不作分母"） |
| `gofmt`/`gofumpt` | `gofumpt.exe -l cmd/wisp internal/proc`（宿主现成 `v0.12.0 (go1.27.1)`，未装任何新东西） | **零输出 rc=0**；对本腿面六件单列同样零输出 |
| `go vet` 双 GOOS | `go vet ./cmd/wisp/`／`go vet ./internal/proc/`／`GOOS=linux go vet ./internal/proc/` | **三条全 rc=0** |
| `go vet` 的 linux `cmd/wisp` | `GOOS=linux go vet ./cmd/wisp/` | **rc=1**，唯一诊断链 ＝ `imports github.com/k2-fsa/sherpa-onnx-go-linux … build constraints exclude all Go files`（第三方模块，与 `R-131r3-5` 同形；零条落在本仓源码） |
| `d22scan` | `sh scripts/d22scan.sh` | **rc=0 clean**；正控 `runtests.sh: OK - packages=[./...] top-level: PASS=35 FAIL=0 SKIP=0, === RUN=77, '[no tests to run]'=0` |
| `d22scan` 台账各 scope | 同上，末行 | 今量 `bans #1-5 internal/=363, cmd/=58, ban #6 frontend/=116, ban #7 internal/tools/=48, ban #8 design/=39, frontend/=116, internal/=905, cmd/=211` |

**整包 44 红为什么"不作分母"**：本程 `git diff --name-only` = 0 行（没改一个字节），所以这 44 枚是 `45371ca0b` **自己的**读数，不是本程造的。
逐名归族（44 枚，无一落在本票面）：语音/模型族（`Test15e*`、`Test16r1*`、`Test26P*`、`Test26R3*`、`Test381r2*`、`Test382WR2*`、`Test390r1*`、`Test392*`、`Test394R1*`、`Test395r1*`、`Test41a2r1*`、`TestWireR1*`，红因逐字含
`speech: onnxruntime 1.17.1 is older than the 1.28 the linked sherpa-onnx C API requires; a second onnxruntime.dll is shadowing the one on PATH`）＋窗口/焦点族（`TestAC4FocusReturnToPriorWindowGap33r5`）＋常驻汇族（`TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink`）＋另两枚（`TestTicket256ResidentGateOptionsFieldSetIsTheFiveItClaims`、`TestAC247HandingTheLevelToTheSeamIsNotVisibility`）。
⇒ 与 `9b0f7d555` 那笔（RUN 452/PASS 441/FAIL 7/SKIP 3，七红逐名、无一在面内）**同形**：整包在本机不由本票字节决定。
⚠ **`-count=2` 的整包四数本程未取**（一发 662s，两发 ~22min，且红名册是同一批台面族）；票面 AC#6 这一条因此**仍差一发**，谁取谁写。

## §6 还原与洁净（末了逐枚现量）

| 件 | 变异前/还原后 sha256 | 判 |
|---|---|---|
| `cmd/wisp/dataroot_128_test.go` | `3c11ecda1964bb258d90cdeb9a6d25cb3c9c4d7500c734a9377b2c2de1a1d2a4` | 与变异前逐字相同 |
| `cmd/wisp/doctor.go` | `d7e1f1fbb274cd56b6c68eebaefd98f987e24e71332b9eb71f892abe33a91643` | 与 `git show HEAD:` 逐字相同 |
| `internal/proc/envfork.go` | `a37c4ecc0e867dad1720405feade0fa8fffa4aab080635a4290dca949651bc3a` | 与变异前逐字相同 |
| 全树 | `git diff --name-only` | **0 行** |
| 五发变异的落地证明 | 每发 `grep -n "MUTATION …"` 命中 1 处＋`go build` rc=0 | 逐发在案 |

## §7 这一程**没**做什么（留给非实现者与下一程）

1. **一格未翻框**、未动 `Status:`、未动票名、未碰台账 `pending-and-issues.md`／`HANDOVER.md`。
2. **AC#7 的 (a)/(b) 两拍未重跑**（§4.5 具名）。
3. **整包 `-count=2` 四数未取**（§5 末段具名）。
4. 裁决表仍须**非实现者**出（`SPEC-12 §4.3` #1/#3、`issues/README`「Hard global constraints」末条、D22 双角色）。
