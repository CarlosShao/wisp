# 裁决表：census GUARD D 剩余 4 枚包级洞的逐枚分类与 claim（census-guard-d 腿，2026-11-08）

来源＝票 409 speech-runner 腿的证据件 `speech-runner-routing.md` §4：speech 闭合后 census 的 GUARD D 仍剩
`cmd/balldebug`／`internal/agent/scheduler`／`internal/command`／`internal/plugin/manifestslot` 四枚存量洞。
本腿按派单逐枚分类并 claim，**只 spend claim，不动 census/tiers=/任何 GUARD 的逻辑一字**。

## 1. 三核（全部现跑）

1. **洞还在**：`bash scripts/portable-tests.sh --scope=census`（工作树）现跑，四行逐字
   `NO-SCOPE <-UNCLAIMED-HAS-TESTS`（balldebug 1/0、scheduler 2/0、command 1/0、manifestslot 7/0），
   `census totals: ... unclaimed-with-tests=4`，GUARD D rc1。
2. **引文出处**：GUARD D 判据现读于 `scripts/portable-tests.sh` census 段——无 pin 命中且
   `{{len .TestGoFiles}}/{{len .XTestGoFiles}}` 非 `0/0` 即记洞（不可读按洞算，default-deny）；
   票 111 `AC#12` 归属理由逐字「本票就是 GUARD D 的出处」。
3. **红/绿隔离**：两棵 `git clone --shared` 导出树——tree1 @`25ad83737`（首 commit）、
   tree2 @`c121390a3`（fix commit），untracked 资产 `build/`＋`third_party/spike-models`＋
   `third_party/sherpa-onnx` 逐件拷入；门禁全部在导出树跑。

## 2. 逐枚分类：**四枚全甲（有测试、无档认领），零乙、零丙**

| 包 | windows 编译测试 | 引入票（测试文件头现读） | 判词 |
|---|---|---|---|
| `cmd/balldebug` | 1/0 | 票 71（`diffgate_windows_test.go` 头注「Ticket 71 AC#1/AC#2」） | 甲。测试 `//go:build windows`——但 census 本来就跑在 GOOS=windows 上（CI 唯一调用点＝windows 腿），编译数非 0，GUARD D 咬的正是它；tag 语义归票 111 `AC#12` 用例级那格，本腿不碰 |
| `internal/agent/scheduler` | 2/0 | 票 47（`scheduler_test.go` 头注「Ticket 47 acceptance legs」） | 甲。零 //go:build／t.Skip／os.Getenv／exec.Command（grep 实测）；票 111 Progress log 登记过它是「实现票落地时必须把包放回 scope」的 DEFERRED 桩——测试已落地，本 claim 是**既登记义务的执行**，不是新立法 |
| `internal/command` | 1/0 | 票 213（`roster_test.go` 头注） | 甲。同上量杆全零；GOOS=linux 测试二进制编译＋`go vet` 双过 |
| `internal/plugin/manifestslot` | 7/0 | 票 50/51/52/138（七个 `_test.go` 的头注） | 甲。测试 exec `tasklist`、探 `SystemRoot`/`LOCALAPPDATA`——core 档「零 exec.Command」量杆恰是它不过的那条，读数归供得起主体的腿 |

## 3. 落刀（`scripts/portable-tests.sh`，八处、零逻辑行）

- `core_pin`＋core scope：`internal/agent/scheduler`、`internal/command`（pin 与 scope 同 commit，GUARD C 两向）。
- `win_pin`＋windows scope：`cmd/balldebug`、`internal/plugin/manifestslot`（同上）。
- 注释段追加 `UPDATE 2026-11-08` 一段：逐枚量杆与两档分配理由（scheduler/command 走 present/watchdog 杠；
  balldebug 进 core 会被 ubuntu 腿的 GUARD A 当天拒＝linux 上 0 编译测试；manifestslot 撞 core 零 exec 杠）。
- ⛔ census 段／`tiers=`／GUARD A–D 本体：一字未动（diff 可验）。

## 4. 验收读数（判据＝census 自检 rc0）

| 发 | 命令（台面） | 读数 |
|---|---|---|
| census | `--scope=census`（tree2 @`c121390a3`） | `census totals: packages=42 ... claimed-by-no-scope=5 unclaimed-with-tests=0`，**rc0，GUARD D 剩 0**；四行分别 `windows/core/core/windows` 认领 |
| speech 路由未破坏 | `--scope=speech`（tree1 @`25ad83737`；tree2 @`c121390a3` 复跑） | 两发同数：**RUN=106 PASS=89 FAIL=0 SKIP=0，rc0**（含标点模型资产） |
| selftest | `portable-tests-selftest.sh`（tree1；tree2 @`c121390a3` 复跑） | **32 case / 0 断言失败，GREEN**（载体运行时从脚本字节现抽 `core_pin`，claim 增行不碰裂它） |
| core 档全量 | `--scope=core`（tree2，本机 windows dev box） | rc1＝**存量红**（见 §6 归因）；本腿两枚新增 `ok`：scheduler 0.998s、command 1.168s，零 skip |
| windows 档全量 | `--scope=windows`（tree2） | rc1＝**存量红**（同 §6）；本腿两枚新增 `ok`：balldebug 0.565s、manifestslot 15.158s，零 skip；TEMP 置回形复跑 `ok 12.862s`/`ok 0.195s` 两形全绿 |
| GUARD C 反向自证 | tree1 首跑 windows 档 | 本腿第一版 scope 行误把父包 `./internal/plugin/` 顶掉 ⇒ `GUARD C - Pinned: 13, resolved: 12` 当场点名 `internal/plugin` rc1 ⇒ 门咬合有效；修复 commit `c121390a3`（父子并存） |

## 5. 门禁（导出树 tree2 @`c121390a3`，`unset GOTMPDIR TMPDIR TEMP TMP`；PATH 带 build+sherpa）

| 门 | 读数 |
|---|---|
| `gofumpt -l scripts/`（`/d/work/base/gopath/bin/gofumpt.exe`） | **0 命中** |
| `go vet ./...` | **rc0**（CI 等价面，见 §6 第二条） |
| `go build ./...` | **rc0**（同上） |
| `sh scripts/d22scan.sh` | **clean**，no D22 ban violations |
| `bash -n`（portable-tests / selftest / speech-tests 三枚） | rc0 |
| LF | `scripts/portable-tests.sh` CRLF=0；证据件与本票行同 LF |

## 6. 如实

1. **core/windows 档全量红＝存量，非本腿造成**。判据是结构性的：`git diff --stat 2077416a5..HEAD -- "*.go"`＝
   **0 行**（本腿两 commit 只动 `scripts/portable-tests.sh`＋票 111 日志），所以红集合里的包
   （proc `TestCrossVetForLinux`、ball `TestOrb*`×3、**父包** `internal/plugin` 的 `Test52SchemaViolations`/
   `Test52CommandKindValidates`、tools `TestTicket252R1AlignmentAddsNoAuthorization` 等 FAIL=9；
   SKIP=13 含 aecspike rig 门控族等）源字节与父提交逐位相同。这些包全是 windows/core 档的**既有成员**，
   没一枚是本腿加进去的。另：父包 `internal/plugin` 自带一份 `command_52_test.go`（与 manifestslot 同名族），
   其 `Test52AbsoluteTempPathIsRefusedForTheTempReason` 在 TEMP 未设环境下 skip（`no TEMP in this environment`）——
   本机按纪律 `unset TEMP` 跑出的 SKIP=1 来自**父包**（manifestslot 无此测试，grep 验），TEMP 置回后该枚 PASS、
   另两枚转 SKIP；windows 步在 CI 上本就红（票 111 Progress log：`Portable windows tests` failure 12/12），
   归既有各票，不销本腿账。
2. **vet/build 首发红是 untracked 杂件染色**：`build/mut_rvw.go`（11-04，非本腿、非 tracked——
   `git ls-files build/` 空）两个 undefined 符号让 `go vet ./...`/`go build ./...` rc1；把该目录挪出模块后
   双 rc0。CI 等价面（tracked 树上根本无此目录）＝rc0 读数为准，杂件原样留在工作树（只建不删）。
3. **丙形（tag-gated）零枚**：四枚中唯一带 tag 的测试文件是 balldebug 的 `//go:build windows`，而 census 的
   唯一 CI 调用点就在 windows 腿上，编译非 0，GUARD D 咬得着——它不需要也不允许「写清 tag」式白名单。
   「认领必须携带 runner-OS 语义」那一格仍归票 111 `AC#12`（已翻勾、排程在案），本腿未替它立法。
4. 票池归属：Progress log 一行落**票 111**（GUARD D 出处票，append-only）；四枚测试的引入票
   （47/213/50-51-52-138/71）只作分类凭据引用，未动票面。
