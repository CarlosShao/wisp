# 票 135 裁决（非实现者 / 只读腿）— AC#1–AC#7 逐条

- 裁决腿：`135-verdict`＝**非实现者**（实现腿 `135-r1`／`135-r2`；本腿只读、只写本件）
- 时刻：2026-11-05 21:0x–21:2x +08
- 现跑 HEAD：`b914db12e`（派单写 `7eb64d5ce`）。`git diff --stat 7eb64d5ce b914db12e`＝**3 件且全非产码**（`docs/evidence/s1/51-verdict.md`、`docs/reports/pending-and-issues.md`、`.scratch/wisp/issues/51-tier2-goja.md`→`-done.md`）；`-- '*.go'` **0 行** ⇒ 代码面就是 `7eb64d5ce` 的树，下读数对两枚锚都成立。
- 环境：`go1.27.1 windows/amd64`；`PATH="$PWD/build:$PWD/third_party/sherpa-onnx:$PATH"`、`unset GOTMPDIR TMPDIR TEMP TMP`、`CGO_ENABLED=1`（**未**置 0）；`gofumpt v0.12.0 (go1.27.1)`（宿主现成 binary，未装新东西）。
- 台面（诚实记账）：取读数时 CPU **94%**、在飞 `go` 进程 **3 枚**（编队跑动中）⇒ 整包重活本腿**未采**，缺口具名在 §6。
- 方法：**所有变异都在仓外快照**（`/d/tmp/wisp135v/<树>`＝`git archive HEAD` ＋ 三枚 sherpa DLL）里做，**工作树一个字节没动**：末了 `git status --porcelain -uno` ＝ **0 行**，起手与末了两量。
- 判据：**盘上现跑**，不采信任何证据件的自述读数；`docs/evidence/s1/135-ac1-ac6-impl.md` 只当「尺在哪」的索引与「该复算什么」的清单。
- 未改代码、未改工单、未改别腿证据件；本会话唯一写入＝本文件。

## §1 七格结论（先给可结案的判断）

| AC | 判 | 承重读数（本腿现跑，不是引用） | 一句话可结案 |
|---|---|---|---|
| **#1** 后果断言装自证腿 | **PASS** | `dataroot_128_test.go:506` 三形（树→响且点名／空→静／不存在的目录→响）；**M-C 弄哑 `:493`** ⇒ 七名切片 `19/11/8/0 rc=1`，红的正是 AC#1＋AC#2A；**同变异下摘掉 AC#1** ⇒ `18/11/7/0 rc=1`（AC#2A 六腿仍红） | 牙有两名看守，且检查器被弄哑时**第一名看守自己先响**——判据两条都当场成立 |
| **#2** M-F 做成常备用例 | **PASS** | `:600` 六腿（`runTextTask`/`cmdModels`/`cmdProviders`/`cmdSessions`/`cmdDoctor`/`resolveSecretLayout`）各植 `wisp-dev/logs/*.jsonl`＋`config.toml`＋`secrets/key.dpapi`；M-C 下六腿全红、红名逐字 `:621` 点「目录非空」并**点出 `wisp-dev`** | 红名点的是「空检查一声不响」，不是 rc——票面那句「不要只点到 rc」当场满足 |
| **#3** doctor 腿独立 CWD 牙 | **PASS** | `:545` 用**生产自己的探针名** `doctor-write-probe.tmp` 做哨兵，并用 `runtime.Caller`＋读 `doctor.go` 自证耦合；**M-A2b**（拒绝句一字不动＋另起 `probeWritable(".")`）⇒ 唯一红＝AC#3，红名 `:577` 点「探针删掉自己的文件」 | marker 之外有一枚独立读数；且那一发是**今天唯一能复现的形** |
| **#4** 自证用例真读生产文案 | **PASS** | `:219` 先 `failConfigDir128` 再 `resolveDataDir("dev")` 取**生产** `err.Error()`，三条**成句**（连接词＋句末标点）＋次序＋`HasSuffix("然后重试。")`；**M-E** ⇒ `cmd/wisp` 唯一红＝AC#4（`:260` 三条＋`:269` 句末） | 六枚子串全在而句子拆光时它必红；同变异下六子串合取检查**仍全过**（AC2EveryLeg PASS）⇒ 新牙与旧清单相互独立 |
| **#5** POSIX 有分母（选“补上”） | **PASS** | `internal/proc/envfork_test.go:351`（**无 GOOS 后缀**，`GOOS=linux go vet` rc=0 证它编得进 linux）；先自证前提（`os.UserConfigDir()` 必须失败）；**linux 容器**（`golang:1.27`、`GOPROXY=off`＋宿主 modcache）单发 **PASS**；**M-D′**（`DefaultLayout` 回落 `"."`）下 linux 单发 **FAIL**，红名含 `DataDir:wisp-dev` | POSIX 有分母**且会响**，不是「只在 windows 有证人」；全件无一处写「两腿一致」 |
| **#6** 门禁 | **PASS（一处具名缺口）** | build rc=0；`gofmt -l`／`gofumpt -l` 零输出；`go vet` windows 两包＋`GOOS=linux ./internal/proc/` **三条 rc=0**；`sh scripts/d22scan.sh` **rc=0 clean**（正控 `PASS=35 FAIL=0 SKIP=0`、`=== RUN=77`，scope `363/58/116/48/39/116/905/211`）；本票切片 `-count=2 -v` **38/38/0/0 rc=0**；`internal/proc` 整包 `-count=2 -v` **106/104/0/2 rc=0**；票 123 那批 **rc=0 全绿**（未放宽） | 能绿的门**逐条现跑绿**；缺口＝整包 `./cmd/wisp/ -count=2` 四数未取（§6.1，两程皆未取，且该包本就 44 枚旧红） |
| **#7** 131 第二把尺收到句段粒度 | **PASS** | 尺仍是**原始文本**扫描（`rulingAnchorLines131:477`），自证腿 `:595` 喂真 `rulingCrossChecks131`；**两拍本腿亲自重跑**（下表）：**(a)** 放宽谓词＋真裁决句改名 ⇒ 不修树红（`in none of those files` 1 命中）、修后仍红（同句 1 命中）；**(b)** 再加裸标记到 `slo_other.go` 末尾 ⇒ **不修树（文件粒度）绿 rc=0**、修后**红**且点到粒度（`carries the marker at line(s) [43]`），`in none of those files` **0 命中**；反向那枚「annexing to line 43 … prose」同现；**把尺退回文件粒度 ⇒ 自证腿自己红**（`:618`） | 两拍都复现、且红名点的是「粒度」不是「缺标记」；这一格此前只有 r1 的自报读数，现在有非实现者的现跑 |

## §2 台面读数（所有结论的底数）

| # | 命令（仓外快照 `/d/tmp/wisp135v/base`＝`git archive HEAD`） | 读数 |
|---|---|---|
| 1 | 本票七名 `-count=1 -v` | **7 PASS、`ok 1.315s`、rc=0**（`=== RUN` 19＝7 顶层＋12 子腿） |
| 2 | `go test ./internal/proc/ -run '^TestDefaultLayoutRefusesWithNoUserConfigDir$'` | **PASS、`ok 0.104s`** |
| 3 | 七名 `-count=2 -v` | **RUN=38 / PASS=38 / FAIL=0 / SKIP=0，rc=0，0.944s** |
| 4 | `./internal/proc/ -count=2 -v`（整包） | **RUN=106 / PASS=104 / FAIL=0 / SKIP=2，rc=0，2.829s**（两枚 SKIP＝`TestHelperProcess`，与 r2 同形） |
| 5 | `go build ./cmd/... ./internal/... ./tools/...` | **rc=0**（零输出） |

⚠ 计数口径：本件 PASS 把**子测试**计进（38＝2×19）；r2 的 `56/20/0/0` 是**十名顶层**且 PASS 只数顶层（20＝2×10）、RUN 含子测试。两套数**不是矛盾**，是我在读数时就换了尺——下一程引数请连口径一起引。

## §3 变异逐发（每发：`grep` 证落地＋`go build` rc=0 再读红名）

| 发 | 造了什么（仓外快照） | 落地证明 | 读数 | 判 |
|---|---|---|---|---|
| **M-C** | `dataroot_128_test.go:493` `if len(found) > 0` ⇒ `if false && …` | `grep -n "MUTATION M-C"` 1 处、build rc=0 | **19/11/8/0 rc=1**；红＝AC#1＋AC#2A（6 子腿），其余 5 绿 | 牙活着 |
| **M-C 摘 AC#1** | 同上，`-run` 去掉 AC#1 | —— | **18/11/7/0 rc=1**（AC#2A 六腿仍红） | AC#1 自己的结案判据成立 |
| **M-A2**（票面存证的形） | `doctor.go:107` `probeWritable(dir)`⇒`probeWritable(".")` | build rc=0 | **19/19/0/0 rc=0**（全绿） | **死码**：票 128 后 `resolveDataDir` 一律拒绝，`else if` 走不到 ⇒ 照抄该形今天测不到任何东西（r2 的更正**成立**） |
| **M-A2b**（忠实的那形） | 拒绝句一字不动，另起 `_ = probeWritable(".")` | build rc=0 | **19/18/1/0 rc=1**，唯一红＝AC#3 | AC#3 的牙活着，与 marker 无关 |
| **M-E** | `envfork.go:273` 换成「六 marker＋路径＋OS 原话**全在**、但无成句/无句末标点」的一串 | build rc=0（两包） | `cmd/wisp` **19/18/1/0 rc=1**，唯一红＝AC#4（`:260`×3＋`:269`）；`internal/proc` **1/0/1/0**，红＝AC#5（`:407`×3＋`:416`） | 新牙独立于旧六子串：同变异下 `AC2EveryLeg` **PASS**，六子串检查**零** `never says`／`cannot be acted on` |
| **M-D′**（AC#5 的 POSIX 分母） | `DefaultLayout` 在 `os.UserConfigDir` 失败时 `dir, err = ".", nil` | `grep` 1 处、build rc=0 | **linux 容器**单发 **FAIL**，红名 `DefaultLayout(dev) answered {… DataDir:wisp-dev …} with no error` | POSIX 侧分母**会响**，不止「有分母」 |

**AC#7 两拍（本腿亲自重跑，这是 r2 具名没做的那一项）**——变异＝`leg_sink_gate_131_test.go:1070` 谓词放宽成 `… || fi.name == "cmdSLO"` ＋ `slo_windows.go:184` 真裁决句改名（`WISP-LEG-SINK-RULING-OLD135:`）：

| 拍 | 尺 | 读数 | 红名逐字（节选） |
|---|---|---|---|
| **(a)** 放宽＋真裁决句没了，**无**裸标记 | 不修树（文件粒度） | **rc=1**，`in none of those files` **1 命中** | `leg "slo" is booked as ruled at [slo_other.go:34 \| slo_windows.go:204], but the marker sentence "WISP-LEG-SINK-RULING:" is in none of those files` |
| **(a)** 同上 | 修后（今天这把尺） | **rc=1**，同句 **1 命中** | 同上（逐字相同） |
| **(b)** 再追加裸标记到 `slo_other.go` 末尾 | 不修树（文件粒度） | **rc=0 全绿** | ——（**这就是票面要立的那一格：两枚变异叠起来把红洗掉**） |
| **(b)** 同上 | 修后（今天这把尺） | **rc=1**，`in none of those files` **0 命中** | `… is in one of those files but not at sentence-segment granularity: slo_other.go carries the marker at line(s) [43], none of them a site this row names`；反向：`slo_other.go carries a "WISP-LEG-SINK-RULING:" sentence anchoring to line 43, and no ruled ledger row's site names that line: … is prose, not a …` |
| **自证腿的哑法** | 把正向谓词悄悄退回文件粒度 | **rc=1**，红的正是 `TestAC7ARulingParkedInAFileBacksNothing131` | `AC#7 RED (the instrument, not the code): a sentence parked at the end of the file … produced []. … Anchors computed for this fixture: [6]` |

⇒ AC#7 的两拍在**今天的树上逐拍复现**，且红名点的是**粒度**；自证腿也确实**绑在那条粒度规则上**（规则一退回文件粒度它自己就红）。

## §4 落点与「只碰测试面」的现读

五枚落地 commit 逐枚 `git show --name-only`：

- `b46a867b1`／`b53b79bda`／`144b038f5`：`cmd/wisp/dataroot_128_test.go`（＋票面／`probes/135/*.txt`）
- `bf25809a3`：`internal/proc/envfork_test.go`（＋同上）
- `2781ec380`：`cmd/wisp/leg_sink_gate_131_test.go`（＋同上）

⇒ **产码 0 字节**（本票是「判据仪器缺陷」，交付面本来就该只有测试）——这同时是「票 123 那批不许放宽换绿」的**最强证明**：那些文件一个字节没被本票碰过。补一手实跑（不是只靠文件表）：`TestTicket101*` 五枚＋`TestComposedGateBlocksAWriteForTwoSeconds` ⇒ **9 RUN 全 PASS、rc=0、139.651s**（后者 3.95s，仍是「两秒」那条闸）。

r2 那枚 commit `7eb64d5ce` 只带**票面＋本件**（零码），与它「零改码」的自述一致。

## §5 逐格与票面文字的对照（易被读松的地方）

| 票面文字 | 今天树上的形状 | 判 |
|---|---|---|
| AC#1「主动把 `assertDirEmpty128` 弄哑（**或喂它一个不存在的目录**）」 | 两样都有：`recorder128` 走三形，第三形就是「不存在的目录必须响」（walk 错不得被读成「空」） | 满足（比票面多给了「空目录必须静」这一反向腿） |
| AC#2「红名点到“目录非空”**而不是只点到 rc**」 | `:621` 的 Fatalf 在 rc 已非 0 之后才发；`:626` 还要求点名 `wisp-dev` 与 `wisp-20260923-001.jsonl` | 满足 |
| AC#3「marker 一撤就无牙」 | 用例与 marker 零耦合（M-A2b 下 marker 全过、只有它响） | 满足 |
| AC#4「成句断言而不是六段子串」 | 三条子句带连接词与句末标点＋`at` 次序检查＋`HasSuffix` | 满足（**加断言，未改判据**：六子串合取与「清单不可缩短」负腿原样在） |
| AC#5「二选一，不许沉默」 | 选「补上」，且全件无「两腿一致」（票面 3 处命中全是**禁令／引用**：`:25`／`:38`／`:220`） | 满足 |
| AC#6「`gofmt`/`gofumpt` 真跑」 | 宿主现成 `gofumpt v0.12.0`，`-l cmd/wisp internal/proc` 零输出；`gofmt -l` 同 | 满足 |
| AC#7「② 装一枚主动弄哑自己的自证腿」 | 尺未折回「只信 AST」：`rulingAnchorLines131` 是纯文本行扫，与 AST doc 走读零共享代码 | 满足 |

## §6 本件具名**没**做到的（不许被读成已做）

### §6.1 AC#6 的整包 `./cmd/wisp/ -count=2` 四数：**仍未取**（两程皆未取）

- r2 在 `135-ac1-ac6-impl.md` §7.3 具名它没取；本腿取读数时 CPU 94%、在飞 `go` 进程 3 枚（编队跑动），按本仓「不在编队跑动时采读数」**没有**再加第 4 枚重活。
- 它**不可能翻本格结论**：该包在 `45371ca0b` 上就是 `640/439/44/2 rc=1`（44 红全在语音/模型、窗口焦点、常驻汇三族，零枚在本票面，与 `9b0f7d555` 同形）⇒ 整包的 rc 由别族决定，不由本票字节决定。
- 谁取谁写，并请在件里连**计数口径**（是否含子测试）一起写。

### §6.2 r2 证据件 §5 里那条 `GOOS=linux go vet ./cmd/wisp/` 的**归因句今天不复现**（结论不变）

- 本腿现量：`GOOS=linux go vet ./cmd/wisp/` **rc=1**，但首故障链是 **`runtime/cgo` 交叉编译**（mingw 的 `gcc_*` 找不到 `sys/mman.h`、`pthread_unix.c` 的 `-Werror`），**不是** r2 写的 `sherpa-onnx-go-linux … build constraints exclude all Go files`。
- 两件事都不变：**rc=1**、且**零条诊断落在本仓源码**（`grep -c "cmd/wisp"` ＝ 0）。
- ⇒ 判：AC#6 的「go vet 双 GOOS」在本机成立（windows 两包＋linux `internal/proc` 三条 rc=0；linux `cmd/wisp` 的 rc=1 归第三方/工具链），**但 r2 那句归因是环境相关的，下一程别照抄**。

### §6.3 本件**不**翻任何框

未动票面 `Status:`／勾框／票名／台账／HANDOVER／`frontend/**`／`design/**`；AC#8 已由另一张非实现者表（`135-ac8-r1-acceptance.md`）结案，不在本件射程。

## §7 复算命令（想自己重量的，照这个跑）

```bash
# 1) 仓外快照（工作树不动）
B=/d/tmp/wisp135v; mkdir -p $B/base && git archive HEAD | tar -x -C $B/base
mkdir -p $B/base/third_party/sherpa-onnx && cp third_party/sherpa-onnx/*.dll $B/base/third_party/sherpa-onnx/
# 2) 每枚变异：grep 证落地 -> go build rc=0 -> 读红名
#    M-C   : cmd/wisp/dataroot_128_test.go:493   if len(found) > 0  ->  if false && len(found) > 0
#    M-A2  : cmd/wisp/doctor.go:107               probeWritable(dir) -> probeWritable(".")
#    M-A2b : cmd/wisp/doctor.go:105-107           拒绝句不动，另起 _ = probeWritable(".")
#    M-E   : internal/proc/envfork.go:273-276     换成六 marker 齐全但无成句的一串
#    M-D'  : internal/proc/envfork.go:285-287     return Layout{}, …  ->  dir, err = ".", nil
#    AC#7  : leg_sink_gate_131_test.go:1070       … || fi.name == "cmdSLO"
#            slo_windows.go:184                   真裁决句改名
#            slo_other.go 末尾                     追加一行裸标记（拍 b 才加）
# 3) 门禁
export PATH="$PWD/build:$PWD/third_party/sherpa-onnx:$PATH"; unset GOTMPDIR TMPDIR TEMP TMP
go build ./cmd/... ./internal/... ./tools/...
/d/work/base/gopath/bin/gofumpt.exe -l cmd/wisp internal/proc
go vet ./cmd/wisp/ && go vet ./internal/proc/ && GOOS=linux go vet ./internal/proc/
sh scripts/d22scan.sh               # 正控 PASS=35 FAIL=0 SKIP=0、RUN=77、rc=0 clean
# 4) POSIX 分母（AC#5）——CGO 不关
docker run --rm -v "D:/tmp/wisp135v/<树>:/src" -v "D:/work/base/gopath/pkg/mod:/gomodcache:ro" \
  -w /src -e GOMODCACHE=/gomodcache -e GOFLAGS=-mod=mod -e GOPROXY=off -e CGO_ENABLED=1 \
  golang:1.27 go test ./internal/proc/ -run '^TestDefaultLayoutRefusesWithNoUserConfigDir$' -count=1 -v
```

## §8 一句话总判

**七格 AC#1–AC#7 全部判 PASS、零格退回**：AC#1/AC#2（同一枚牙的正反面）、AC#3、AC#4、AC#5、AC#7 的「坏了会响」都由本腿**在今天的树上亲手造出来**（8 形变异、12 发读数：M-C／M-C 摘腿／M-A2／M-A2b／M-E 两包／M-D′／AC#7 (a) 两拍／(b) 两拍／自证腿退回文件粒度），逐发 `grep` 证落地、`go build` rc=0 之后再读红名；**变异全在仓外快照里做，工作树起手与末了 `git status --porcelain -uno` 均 0 行**；
落地面**只有测试与证据**（产码 0 字节，票 123 那批一个字节没碰且实跑全绿）；
**唯一具名缺口**＝AC#6 里整包 `./cmd/wisp/ -count=2` 的四数仍未取（两程皆未取，且该包本就 44 枚旧红、不由本票字节决定），另有一句环境相关的 vet 归因需下一程改写（§6.2）。
票 135 可结；下一程若要为「零缺口」再落一笔，只需补 §6.1 那一发四数并回写本件 §6.1。
