# 342 AC#4 读数件：消费者五包两平台成对 ＋ `cmd/wisp` 那一发 ＋ 门禁

- 本腿＝**读数腿**（⛔ 实现、⛔ 选形、⛔ 翻票框、⛔ push）。写面＝本件 **1 枚**；产码与测试字节 **0**（收尾现量 `git status --porcelain --untracked-files=no` 里本腿名下 0 行）。
- 锚＝`e60b76877`（派单 HEAD，16:04 现认；本机起跑 16:05、容器输入树 `git archive` 16:07，**两档的输入字节都取于 16:15 前**）。★**共树披露**：工作树 HEAD 于 16:15:58 前移到 `ed09c7376`，`git diff --name-only e60b76877..ed09c7376` ＝ **1 枚文档**（`docs/evidence/s1/220-r3-preflight.md`，+46 行，`.go`/`.sh` 0 枚）⇒ 本件全部读数在**同一份产码字节**上成立。
- 工具链：本机 `go1.27.1 windows/amd64`、`CGO_ENABLED=1`、gcc 16.2.0（MSYS2）；Linux 档＝容器 `golang:1.27`（`go1.27.1 linux/amd64`、gcc 14.2.0、**`CGO_ENABLED=1`**、`GOPROXY=off`、`GOTOOLCHAIN=local`、宿主 module 缓存只读、`TMPDIR`/`GOTMPDIR` 留在容器内 `/tmp`，⛔ 绑 Windows 卷）。
- 纪律：两平台各一发 `-count=1`；Linux 档用 `git archive HEAD` 的**仓外树**（`/d/tmp/wisp/20261105/342ac4/tree`），⛔ `-overlay`、⛔ 复用别人的树；本机跑前 `unset GOTMPDIR TMPDIR TEMP TMP`；⛔ `go test ./...`。
- ★**342 的提交面现量**：按标题锚取 `票 342 (落地腿|终裁腿|普查腿)` ＋ `342-r3:` 共 **30 笔**，逐枚 `git show --name-only` 取并集，非 `docs/`、非 `.scratch/` 的命中**只有 `internal/winsec/**`**（10 枚：`winsec.go`、`pinnedleaf_342_*` 三枚、`ancestor*_342*`／`348*` 六枚），**非 winsec 行数＝0** ⇒ 五包与 `cmd/wisp` **零字节**归 342。

## 1. 尺（四数怎么读，⛔ 与 gate 口径混用）

- `RUN` ＝ `^=== RUN` 行数（**含**子用例）。
- `PASS_TOP`／`FAIL_TOP`／`SKIP` ＝ 列 0 且测试名不含 `/` 的结果行（**只数顶层**）。
- 逐包命令：`go test -v -count=1 -timeout 900s ./internal/<pkg>`；四数用 `grep -cE` 现取。
- 本仓 `portable-tests.sh` 的四数把 `--- PASS` **全量**计入（含子用例），与本尺**不同**；第 4 节的 gate 四数照它自己的口径**原样抄**，两套数⛔ 互减。

## 2. 消费者五包：本机（Windows host，共享工作树 `e60b76877`）

| 包 | rc | RUN | PASS_TOP | FAIL_TOP | SKIP | 红/跳点名（顶层） | 时长 |
|---|---|---|---|---|---|---|---|
| `internal/risk` | **1** | 281 | 170 | **1** | 1 | FAIL `TestC26RewriteAccountIsConsumedAtEverySecurityLeg`；SKIP `TestSyncRegistryProbeLive` | 14.234s |
| `internal/secret` | **0** | 29 | 21 | 0 | 0 | — | 0.279s |
| `internal/config` | **0** | 220 | 132 | 0 | 0 | — | 2.047s |
| `internal/agent` | **0** | 186 | 153 | 0 | 0 | — | 9.051s（共树）／**8.613s（干净树复跑，四数同 186/153/0/0）** |
| `internal/memory` | **0** | 193 | 103 | 0 | 2 | SKIP `TestSubprocessCrashWriter`、`TestAppendSessionTurnCrashChild` | 35.045s |

★**共树污染披露（⛔ 圆场）**：本机面跑在共享工作树上。收尾 `git status` 有**一枚他人未提交件**：`internal/agent/approval/gate.go`（`+9/−0`，mtime **16:07:14**，非本腿）。因此 `internal/agent` 另用同一份 `git archive HEAD` 干净树（`/d/tmp/wisp/20261105/342ac4/tree`）复跑一发：**rc=0、186/153/0/0、8.613s**，与共树读数逐数相同 ⇒ 那 9 行未改结论。其余四包运行时该件不在其编译面内（`internal/agent/approval` 只被 `internal/agent` 依赖）。

## 3. 消费者五包：Linux 档（容器，`git archive HEAD` 仓外树）

| 包 | rc | RUN | PASS_TOP | FAIL_TOP | SKIP | 红/跳点名（顶层） | 时长 |
|---|---|---|---|---|---|---|---|
| `internal/risk` | **1** | 255 | 152 | **1** | 0 | FAIL `TestC26RewriteAccountIsConsumedAtEverySecurityLeg` | 17.843s |
| `internal/secret` | **0** | 4 | 4 | 0 | 0 | — | 0.011s |
| `internal/config` | **1** | 219 | 128 | **1** | 0 | FAIL `TestHandDeletedSchemaVersionOnTheWriteBranchIsForcedBackAndUnclaimed` | 1.384s |
| `internal/agent` | **0** | 185 | 152 | 0 | 0 | — | 5.016s |
| `internal/memory` | **1** | 192 | 100 | **2** | 2 | FAIL `TestDsn196r3BusinessPoolsKeepTheirOwnForeignKeys`、`TestMigration31V5ToV6KeepsOldData`；SKIP `TestSubprocessCrashWriter`、`TestAppendSessionTurnCrashChild` | 22.573s |

`RUN` 数两平台不同是**预期**（`//go:build` 按平台选测试件：`secret` 29 对 4、`agent` 186 对 185），成对判据按**包级 rc ＋ 顶层红名**读，⛔ 按四数相等读。

## 4. `cmd/wisp` 那一发

**先给口径**：`cmd/wisp` 里**没有任何一枚测试名带 342**（尺：`grep -rhoE '^func Test[A-Za-z0-9_]*342[A-Za-z0-9_]*' cmd/wisp/*_test.go | wc -l` ＝ **0**）。AC#4 点名的"`cmd/wisp` 那一发"＝**整包门**，故按仓里的门跑：

- 门＝`bash scripts/wisp-cli-tests.sh`（转 `portable-tests.sh --scope=cli`，最终 `runtests.sh ./cmd/wisp/ -count=1 -skip <ledger>`，强制 `-v -count=1`）。
- 读数：**rc=1**；门自己的四数＝`=== RUN=271  --- PASS=203  --- FAIL=0  --- SKIP=0`；包行＝`FAIL github.com/CarlosShao/wisp/cmd/wisp 259.253s`。
- ★**红不是断言红**：最后一枚 `=== RUN TestAC13BringUpSurvivesAReusedThreadQuit` **没有结果行**，全日志**零 `--- FAIL`、无 panic、无 `test timed out`**、进程在 259.253s 终止。
- ★**归因（既有、非 342）**：台账 `A1177` 记过**同一形状**——"整包 `cmd/wisp` 本机跑不完（既有、非本笔）…… `HEAD` 原样与改后两发都在 ~301–310 s 整进程终止、零 `--- FAIL`、无 panic、无 `test timed out`；`-skip TestAC14` 只把终止点挪到 `TestAC13`（310.752 s）⇒ 本机环境上限"（`docs/reports/pending-and-issues.md:23199`）。另台账同处记过该包**有时能达终态**（275s；`A1540` 记 220.9s 绿）⇒ 这是**本机环境级摆动**，⛔ 判 342 顶红。
- 门自己的前置守卫都过：ticket 310 的 `GOTMPDIR` 落在 `%LOCALAPPDATA%\wisp-test-gotmp\cli-20261105-162243-21`（同卷），311 的跨卷守卫**未触发**；DLL 名册 `onnxruntime.dll sherpa-onnx-c-api.dll sherpa-onnx-cxx-api.dll` 在场。

**342 邻接的具名子集**（因整包门在本机不可判，补一发可判读数）：`TestTicket198AC3CreatedFileLandsPrivate`、`TestAC3EarlyReplayKeepsTheSinkInsideTheDataRoot`、`TestAC3SecretLegBooksItsAuditRecordsOnDisk`、`TestAC3LogSinkLandsInsideTheEnvDataRoot`、`TestAC3EmptyDataRootIsARefusalNotAFallback`、`TestAC2SealNoticeLandsInTheRunLegLogFile`、`TestAC2RunLegInstallsTheSinkBeforeItsFirstSealingSite`、`TestSecretEndToEndConfigRefResolvesAtRequestTime`、`TestSecretOverwriteIsAnnounced`（＝CLI 上真走 `winsec` 封私密文件/路径解析的那些腿）。

| 命令 | rc | 读数 |
|---|---|---|
| `go test ./cmd/wisp -count=1 -timeout 420s -v -run '^(上述 9 名)$'` | **0** | **9/9 PASS**、`ok … 6.191s`、FAIL 0、SKIP 0 |

## 5. 门禁（HEAD 现跑；跑前 `unset GOTMPDIR TMPDIR TEMP TMP`）

| 门 | 命令 | rc | 读数 |
|---|---|---|---|
| 格式化 | `/d/work/base/gopath/bin/gofumpt.exe -l internal/risk internal/secret internal/config internal/agent internal/memory internal/winsec` | **0** | 列出 **1 枚**：`internal/agent/approval/pending_read.go`（**既有** CRLF，`A1540` 已具名、非本波）；其余 0 |
| vet（本机） | `go vet <五包>` | **0** | 无输出 |
| vet（`GOOS=linux`） | `GOOS=linux go vet <五包>` | **0** | 无输出（交叉档 cgo 自动关，五包纯 Go） |
| vet（`GOOS=darwin`） | `GOOS=darwin go vet <五包>` | **0** | 无输出 |
| vet（`GOOS=linux ./cmd/wisp/`） | `GOOS=linux go vet ./cmd/wisp/` | **1** | `imports github.com/k2-fsa/sherpa-onnx-go-linux: build constraints exclude all Go files` ⇒ **AC#4 逐字点名的那枚交叉档红**（Windows 宿主交叉＝cgo 关、linux 侧 sherpa 全件被排除）；⛔ 与本机绿成对判红绿、⛔ 记新红 |
| D22 | `sh scripts/d22scan.sh` | **0** | `clean - no D22 ban violations`；分母 `bans#1-5 internal=361 cmd=57`、`ban#6 frontend=116`、`ban#7 tools=48`、`ban#8 design=39 frontend=116 internal=893 cmd=208` |
| 长度门 | `sh scripts/check-path-length-budget.sh` | **1** | `over-budget=65 / roster=64 / not in roster=1` ⇒ **RED**；点名 `.scratch/wisp/issues/179-loop-deciderisk-…-done.md`（相对 175、名 154、runner 219 > 165）＋同名"roster 漂移"一条。与 `342-r3` 第 4.3 节**逐数同**（既有、非本波） |

## 6. 逐红归因：没有一枚是 342 推的

| 红 | 平台 | 机制（现读） | 归谁 | 342 新红？ |
|---|---|---|---|---|
| `TestC26RewriteAccountIsConsumedAtEverySecurityLeg` | **两平台同红** | 静态锁：`internal/risk/pathresolver_rewrite_account_test.go:200` 逐行点名 `cmd/wisp/panel_pump.go:283`、`:291`、`internal/panel/git_readonly_exec.go:409`、`:412` 读 `Result.Canonical` 未读改写账 | 台账 `A380`；票 189/190 的消费者面 | **0** |
| `TestHandDeletedSchemaVersionOnTheWriteBranchIsForcedBackAndUnclaimed`（`internal/config`） | 仅 Linux | `writeguard_unnailed_341_test.go:240` `mtime did not move on the write branch`；**容器内轮换**：同一棵 HEAD 树 `-count=3` ＝ **FAIL / FAIL / PASS**；修前基线树 `fde1df02` 单发 **PASS** | 票 341 测试面；`cfg125-r1-impl.md:41` 已登记"容器内已知轮换" | **0** |
| `TestDsn196r3BusinessPoolsKeepTheirOwnForeignKeys`（`internal/memory`） | 仅 Linux（Windows PASS） | `open_196r3_test.go:46` `got file:C:\Users\x\…` vs `want file:C:/Users/x/…`；POSIX 的 `filepath.ToSlash` **不折 `\`**；`-count=2` 两发全红（确定性） | 票 196（`196-r3`，文件生于 2026-10-23） | **0** |
| `TestMigration31V5ToV6KeepsOldData`（`internal/memory`） | 仅 Linux（Windows PASS） | `reminders_test.go:88` `reminders table after the 5->6 step = 0, want 1`；`dbPath := dir + "\\" + dbFileName` 在 POSIX 拼出名字里带 `\` 的另一个文件；`-count=2` 两发全红 | 票 31（`31-r2`，文件生于 2026-11-03） | **0** |

**基线对照**（容器，`git archive fde1df02` ＝ 342-r1 落地**之前**的锚）：同一把 `-run` 尺打过去，`TestC26…` **复现红**（2.225s）；`TestHandDeleted…` **PASS**（0.02s，与 HEAD 的 F/F/P 合起来坐实轮换）；两枚 memory 测试件在该锚**尚不存在**（`[no tests to run]`），其归因由生年＋夹具机制＋Windows 绿三证成立。

## 7. Verdict

**AC#4 satisfied?** **NO（按票面字面）**——票面写"在 Linux 档与本机**全 ok** 的成对读数＋`cmd/wisp` 那一发＋门禁逐枚 rc"，今天有三处不成立：五包里有 **4 枚既有红**（跨平台 1 枚、POSIX 专属 3 枚）、**长度门 RED**、**`cmd/wisp` 整包在本机不达终态**。

**但 342 归因那一半是 YES**：五包**零枚新增红**（342 的 30 笔只动 `internal/winsec/**`；基线 `fde1df02` 上 C26 同红；其余三红有各自生年与夹具机制），`cmd/wisp` 上 342 邻接的 9 枚封私密腿 **9/9 绿**。

**What remains（逐枚点名，全部在 342 写面之外）**：

1. `TestC26…`（两平台）＝ `A380`：`cmd/wisp/panel_pump.go:283/291` ＋ `internal/panel/git_readonly_exec.go:409/412` 要读改写账；归票 189/190 的消费者面。
2. `TestDsn196r3…`（POSIX）＝ 票 196 的夹具把 `C:\` 写死；归 `internal/memory` 的 196 面。
3. `TestMigration31…`（POSIX）＝ 票 31-r2 的 `"\\"` 拼路径；归 `internal/memory` 的 31 面。
4. `TestHandDeleted…`（容器轮换）＝ 票 341 测试面，已在 `cfg125-r1-impl.md:41` 登记为容器已知轮换。
5. `gofumpt -l` 一枚 `internal/agent/approval/pending_read.go`（CRLF，`A1540` 既有）。
6. 长度门 RED ＝ `.scratch/wisp/issues/179-…-done.md` 名与 roster 漂移；归 `scripts` ＋ 台账面。
7. `cmd/wisp` 整包本机不达终态（`TestAC13` 处终止，零 `--- FAIL`）＝ `A1177` 已记的**本机环境上限**；要整包绿须换一台能跑完 WebView2 窗测的宿主/CI 腿。
8. `GOOS=linux go vet ./cmd/wisp/` rc=1 ＝ 票面自己点名的交叉档必然红（cgo 关、`sherpa-onnx-go-linux` 全件排除），⛔ 当新红。

**要把 AC#4 变成字面 yes，需要**：1–3 由各自 owner 修掉（或台账给"登记既有红"的显式豁免）、6 修掉名/roster、7 在能达终态的宿主上跑一发整包绿。**这四件没有一件在 342 的写面里**，本腿不改（派单令：⛔ 改代码/票）。

## 8. 复跑命令（逐字）

```bash
cd "/d/work/workspace/projects plans/Wisp"
unset GOTMPDIR TMPDIR TEMP TMP

# 本机五包（逐包取四数）
for p in risk secret config agent memory; do
  go test -v -count=1 -timeout 900s ./internal/$p
done

# Linux 档：仓外树 + 只读 module 缓存 + 容器内 /tmp
git archive HEAD | tar -x -C /d/tmp/wisp/20261105/342ac4/tree
docker run --rm -v "D:\tmp\wisp\20261105\342ac4\tree:/src" \
  -v "D:\work\base\gopath\pkg\mod:/gomod:ro" -w /src \
  -e GOMODCACHE=/gomod -e GOPROXY=off -e GOFLAGS=-mod=mod -e GOTOOLCHAIN=local \
  -e CGO_ENABLED=1 -e TMPDIR=/tmp -e GOTMPDIR=/tmp golang:1.27 \
  bash -c 'for p in risk secret config agent memory; do go test -v -count=1 -timeout 900s ./internal/$p; done'

# cmd/wisp 门
bash scripts/wisp-cli-tests.sh; echo "gate rc=$?"

# 门禁
/d/work/base/gopath/bin/gofumpt.exe -l internal/risk internal/secret internal/config internal/agent internal/memory internal/winsec
go vet ./internal/risk ./internal/secret ./internal/config ./internal/agent ./internal/memory
GOOS=linux  go vet ./internal/risk ./internal/secret ./internal/config ./internal/agent ./internal/memory
GOOS=darwin go vet ./internal/risk ./internal/secret ./internal/config ./internal/agent ./internal/memory
GOOS=linux  go vet ./cmd/wisp/
sh scripts/d22scan.sh
sh scripts/check-path-length-budget.sh
```

**本腿没做、也不声称的**：⛔ `go test ./...`；⛔ 改任何产码/测试字节；⛔ 翻票框或改票名/Status；⛔ push；⛔ 把 `GOOS=linux go vet ./cmd/wisp/` 的 rc=1 或本机整包终止读成 342 的新红。
