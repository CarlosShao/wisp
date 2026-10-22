# 票 357 落地腿 `357-r1`（part 2）— 落地＋牙＋三读数＋门禁

起手锚 `ed1d9e3c`（裁决书存档）。前件 `docs/evidence/s1/357-r1-impl.md`（156 行、`876aa3776`）是**停手上报**，其地基照收；本件是**落地件**（新件、NOT 追加）。
约定：写作 sec N / "=> " / NOT / ZERO / (1)(2)(3)。门禁前 `unset GOTMPDIR TMPDIR TEMP TMP`、`GOFLAGS=`；未设 `CGO_ENABLED=0`；`cmd/wisp` 抬 `-timeout 30m` ＋ `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"`。每次变异都先 `grep -n` ＋ `go build` 证落地再读读数；变异与读数**串行**。

## 1 改哪几行

**(1.1) `internal/models/ads_windows.go`（新、161 行、`//go:build windows`）** 只读枚举一枚流列表。**零手写 syscall**：`GetFileInformationByHandleEx`（`zsyscall_windows.go:2432`）与 `windows.FileStreamInfo`（`types_windows.go:2765`）都是 x/sys v0.48.0 **已导出**的；手写的只有 24 B `FILE_STREAM_INFO` 布局（`unsafe.Sizeof` 当尺），与 `internal/ball/win32_windows.go:160` 起那串手映结构体同形（裁决书判可那条）。落点：`fileStreamInfo`（24 B 头）/ `streamStartBuf=4<<10` / `streamCeiling=16<<20` / `listNamedStreams`（生产入口）/ `listNamedStreamsBounded(path,start,ceiling)`（两界当实参，让撞界可驱动）/ `parseFileStreamInfo`（逐条先校验长度再 `unsafe` 取）。
**(1.2) `internal/models/ads_other.go`（新、`!windows`）** `listNamedStreams` 恒返 `(nil,nil)`；文件头**具名说清** POSIX 侧那枚用例只能是**控制腿、NOT 牙**。
**(1.3) `internal/models/ads.go`（新、跨平台）** `func (m *Manager) UnnamedStreams(id string) (dir string, streams []string, err error)`：目录走 `m.resolveInstalledDir(id)`＝`VerifyInstalled` 验的那枚 `dir`（有 `local_override` 时是 **override 目录**、NOT `store.root` —— 裁决书点名的"本票最容易做错的地方"）；**逐文件调**（`WalkDir` 里每个非目录项调一次，NOT 拿目录句柄一把梭：目录句柄返空且 `err=nil`）；符号链接跳过不跟随；流名拼 `rel+":name:$DATA"`（与 `dir /r` 同拼法）；排序后返回。
**(1.4) `internal/models/downloader.go`（改 3 行、加 18 行注释＋1 函数）** `:295` `return m.handOffCheck(entry, m.resolveInstalledDir(id))`（原 `:298-302` 三行收成一行）＋ `+resolveInstalledDir`（那三行搬进来、语义一字不动）。**实参仍是同一枚 `dir`**；`Ensure` 的 cache 分支（`:202`）**未动**（换成 override 解析会让 cache 分支指向 override，是行为变更、不在写面内）。
**(1.5) `cmd/wisp/models.go`（改 1 行、加 46 行）** `:253-254` 那句 `与已验签清单一致` **逐字不动**（含 `\n`）；`:262` 原裸 `return 0` 换成 `return announceUnnamedStreams357(io_, id, store.mgr)`；新增 `unnamedStreamsLister357` 接口 ＋ `announceUnnamedStreams357`（枚举空则 `return 0`；枚举出错则 **stderr 具名一句**且仍 `return 0`）＋ `streamNotice357`（那句的唯一拼装处）：`wisp models verify: <id> 在 <stream> 上枚举到清单未点名的备用数据流（ADS dir=<dir>）；本次一致结论只覆盖清单点名的文件`
**(1.6) 新用例 13 枚** — `internal/models/ads_357_windows_test.go`（`windows`、5 枚）`Test357AlternateStreamIsEnumeratedByName` / `Test357TreeWithoutStreamsEnumeratesEmpty` / `Test357ManyStreamsAllComeBack` / `Test357BoundedRetryNamesTheCeilingInsteadOfAnsweringEmpty` / `Test357OverrideDirectoryIsTheOneEnumerated`；`internal/models/ads_357_other_test.go`（`!windows`、1 枚）`Test357POSIXHasNoStreamsToReport`＝**控制腿 NOT 牙**；`cmd/wisp/models_verify_ads_357_windows_test.go`（6 枚）`Test357VerifyTailSaysOneSentenceAndKeepsRC0` / `...WithoutAStreamPrintsNothing` / `...ReportsTheVerifiedTreeNotTheStoreRoot` / `Test357StreamNoticeIsTheOneSentence` / `Test357VerifyTailIsStillOnThePassingPath` / `Test357StoreRootNamingIsNotWhatTheReportUses`。
**为什么今天响**：`UnnamedStreams`／那枚接口／那句／那个调用**改前全都不存在**（改前 `models.go:255` 是裸 `return 0`），故每枚都能被"摘掉本票新增那一跳"打红（sec 3 三发实测）。CLI 侧驱动 `announceUnnamedStreams357` 而 NOT `modelsVerify` 整发：后者前面站着 C29 闸门（manifest 要拿 `buildinfo` 那枚 dev 公钥对应的**私钥**签、私钥按设计在仓外），夹具签不出；接缝因此留在 `unnamedStreamsLister357` 这层 —— **枚举实现与目录解析都是真的**（真 `models.Manager`、真 NTFS、真 `os.WriteFile` 写流），只换"用哪份清单"。`Test357VerifyTailIsStillOnThePassingPath` 用 go/ast 读 `modelsVerify` 函数体，钉住"`:262` 调用还在、且在 `:253` 那句之后"—— sec 3.2 那发变异正是它抓的。

## 2 改前 / 改后 D1 那一发（逐字，含 rc）

台子 `D:\tmp\wisp\20261021\357r1b\`（真 NTFS、零 mock）：`dataG` 只有主体；`dataD1`＝主体＋`model.onnx:attacker` 20 B；`dataOVR`＝`local_override` 指 `override-a1`（ADS 在 override、store 干净）。两枚 exe 由同一工作树改前/改后各 `go build ./cmd/wisp` 一次得到；跑法 `WISP_ENV=test` ＋ `WISP_MODELS_MANIFEST=<仓外 fixture 清单>` ＋ `WISP_TEST_DATA_DIR=<夹具>`；stdout 与 rc 分抓。

```
改前 dataG : wisp models verify: fixture-a1 与已验签清单一致（store=D:\tmp\wisp\20261021\357r1b\dataG\models，0.00s）   rc=0
改前 dataD1: wisp models verify: fixture-a1 与已验签清单一致（store=D:\tmp\wisp\20261021\357r1b\dataD1\models，0.00s）  rc=0
改后 dataG : wisp models verify: fixture-a1 与已验签清单一致（store=D:\tmp\wisp\20261021\357r1b\dataG\models，0.00s）   rc=0
改后 dataD1: wisp models verify: fixture-a1 与已验签清单一致（store=D:\tmp\wisp\20261021\357r1b\dataD1\models，0.00s）
            wisp models verify: fixture-a1 在 model.onnx:attacker:$DATA 上枚举到清单未点名的备用数据流（ADS dir=D:\tmp\wisp\20261021\357r1b\dataD1\models\fixture-a1）；本次一致结论只覆盖清单点名的文件
            rc=0
改后 dataOVR（ADS 只在 override 里）:
            wisp models verify: fixture-a1 与已验签清单一致（store=D:\tmp\wisp\20261021\357r1b\dataOVR\models，0.00s）
            wisp models verify: fixture-a1 在 model.onnx:attacker:$DATA 上枚举到清单未点名的备用数据流（ADS dir=D:\tmp\wisp\20261021\357r1b\override-a1）；本次一致结论只覆盖清单点名的文件
            rc=0
```

=> 〔丁〕四条逐条成立：**D1 多一句点名**（含流名＋清单未点名）／**rc 仍 0**／**G 逐字不变**／`dataOVR` 那句 `ADS dir=` 打的是 `override-a1` 而 NOT 第一行的 `store=…\dataOVR\models`（裁决书第 3 条那个坑落在码面上）。逐字节亲核（`ReadAllBytes`＋`SequenceEqual`）：`before-dataG-stdout.bin` 与 `after2-dataG-stdout.bin` **identical=True**，两枚 `sha256`＝`B6B2D4AA248B1FEAC3C139AE39270A2561CE1E06E1C52BD72295C941EA081E45`（116 B、末字节 `0x0a`、CR=0）；D1 首发那行与改前**逐字相同**（只差 `store=` 路径段）。

## 3 三读数逐字 ＋ 还原核 sha256

**3.1 变异 1 — 把新加那句渲染成空**（`streamNotice357` 体改 `return ""`）。落地 `go build ./... rc=0`，`sha256(models.go)=F3F256AD86C33B1BE5B9CC8F81C406FBB44C6507E6AD508F17EC06E1979CD6B5`。读数（变异红，`go test -count=1 -timeout 30m -run 'Test357' ./cmd/wisp/`，rc=1）**逐名红行**（`<temp>` 只压掉临时目录前缀）：

```
--- FAIL: Test357VerifyTailSaysOneSentenceAndKeepsRC0 (0.02s)
    models_verify_ads_357_windows_test.go:165: ticket 357: the passing path did not name the stream "model.onnx:attacker:$DATA", so the silent pass is still silent: ""
--- FAIL: Test357VerifyTailReportsTheVerifiedTreeNotTheStoreRoot (0.04s)
    models_verify_ads_357_windows_test.go:216: ticket 357: the stream in the local_override directory was not reported: ""
--- FAIL: Test357StreamNoticeIsTheOneSentence (0.00s)
    models_verify_ads_357_windows_test.go:232: ticket 357: the sentence changed.
--- FAIL: Test357StoreRootNamingIsNotWhatTheReportUses (0.03s)
    models_verify_ads_357_windows_test.go:316: ticket 357: the report does not name the override directory <temp>\001\override-dir: ""
FAIL	github.com/CarlosShao/wisp/cmd/wisp	0.315s
```

**3.2 变异 2 — 把调用摘掉**（`modelsVerify:262` 换 `return 0`、渲染函数留着）。落地 `grep -n 'announceUnnamedStreams357'` 只剩 `:274` 注释与 `:288` 定义、`go build ./... rc=0`。读数（rc=1）逐名红行：

```
--- FAIL: Test357VerifyTailIsStillOnThePassingPath (0.00s)
    models_verify_ads_357_windows_test.go:286: ticket 357: modelsVerify (models.go:232) no longer calls announceUnnamedStreams357, so the command has gone back to printing 与已验签清单一致 and saying nothing about a second data stream. That is the exact silence this ticket exists to end, and every unit case in this file would stay green through it
FAIL	github.com/CarlosShao/wisp/cmd/wisp	0.272s
```

**这一发证明"先证落地再读数"的价值：只有它红、其余五枚全绿** —— 没有它，把钩子摘掉整包会全绿而命令已退回静默。

**3.3 读数 2 — 还原绿**：`sha256(models.go)=458B03D2BD22BE5C153807C2676D6A78EB09D980DF4CCC587C176751C0B39356`（与 `models.go.pristine` **同串**）=> `ok github.com/CarlosShao/wisp/cmd/wisp 0.221s`（rc=0）。
**3.4 变异 3 — 去掉 `ERROR_MORE_DATA` 重试**（`ads_windows.go` 改单次调用＋出错返 `nil,nil`）：落地后函数体只剩单次调用、`go build ./internal/models/ rc=0`、`sha256=41421F06328821A827025779DCEA30D7CA3A9DB4C863B0377018EF5230D63CA8`。读数（变异红，rc=1）：

```
--- FAIL: Test357ManyStreamsAllComeBack (0.45s)
    ads_357_windows_test.go:142: ticket 357 ERROR_MORE_DATA leg: 200 streams planted, 0 enumerated. A single 4096-byte call answers ERROR_MORE_DATA and an EMPTY list at this count (357-r1 measured it), so a number below 200 here means the retry is missing and a flooded file reads as a clean one
--- FAIL: Test357BoundedRetryNamesTheCeilingInsteadOfAnsweringEmpty (0.12s)
    ads_357_windows_test.go:171: ticket 357: a ceiling the answer cannot meet returned 0 stream(s) and no error. Form 丁 forbids reading an incomplete enumeration as an empty one, so a bounded retry that hits its bound has to be an error the caller can print
FAIL	github.com/CarlosShao/wisp/internal/models	0.808s
```

`200 planted -> 0 enumerated` 就是裁决书第 1 条那枚**假阴**，本机现量成立。
**3.5 读数 3 — 原样绿**：`sha256(ads_windows.go)=7D4102D4A91BD2DDDCE50C928D6610A0BAE177B0A97FDEABD9D45F1675FA675F`（与 `ads_windows.go.pristine` 同串）=> `ok …internal/models 0.850s`（rc=0、`-run 'Test357'`）；整包 `ok …internal/models 20.091s`（rc=0）。
**3.6 还原核 sha256 三处**（提交版 / 还原台 / `git show HEAD:` blob；第 3 处读数见 sec 8）：`cmd/wisp/models.go` `458B03D2…56`＝`458B03D2…56`；`internal/models/ads.go` `73350E3D…F5`＝`73350E3D…F5`；`ads_windows.go` `7D4102D4…5F`＝同串。

## 4 正例 G 仍绿（〔丁〕第 3 点的牙，三把尺）

(1) 二进制读数：`dataG` 改前/改后 stdout **逐字节相同**（`SequenceEqual=True`、两枚 sha256 同串、116 B）。(2) 用例读数：`Test357VerifyTailWithoutAStreamPrintsNothing` —— 无流目录 => 尾函数 stdout 与 stderr **都是空串**（那句"一致"由 `modelsVerify:253` 打，尾函数不许多说一个字）。(3) 反向对照：`Test357TreeWithoutStreamsEnumeratesEmpty` —— 同一棵树枚举 0 枚流、且 `VerifyInstalled` 仍绿（报告 NOT 拒绝）。

## 5 `ERROR_MORE_DATA` 重试：上界与越界行为

起点 4096 B、增长 x4、上界 16 MiB（＝4096 倍起点，远大于任何模型目录可能出现的流数，同时仍**有限**）：`4096 -> 16384 -> 65536 -> 262144 -> 1048576 -> 4194304 -> 16777216`，第 8 发撞界即返，**NOT 无限**。越界行为＝返 `nil` 流列表 ＋ **具名错误**（`named streams of <path> did not fit in <N> bytes (<err>): the list is incomplete, so this is an error rather than an empty answer`）—— **NOT 静默当空**；`UnnamedStreams` 上抛、CLI 在 **stderr** 打 `… 的备用数据流没枚举完（<err>）；放行结论不变，但这一面今天没看全`，**rc 仍 0**（报告口径、NOT 拒绝口径）。`ERROR_MORE_DATA` 以外的错（打不开/权限/重解析点）同样上抛、同样走那句、**不吞**。记录层第二道：`parseFileStreamInfo` 对每条 `StreamNameLength`/`NextEntryOffset` 先对缓冲校验、越界返错（用例＝`Test357BoundedRetryNamesTheCeilingInsteadOfAnsweringEmpty`，把 `ceiling` 压到 1024 B 逼出撞界那条路）。

## 6 门禁（现跑；四数分开报、NOT 相减）

| 把 | 命令 | rc | 读数 |
|---|---|---|---|
| 1 | `GOFLAGS= go build ./...` | **0** | 无输出 |
| 2 | `go vet ./internal/models/ ./cmd/wisp/` | **0** | 无输出 |
| 3 | `go test -count=1 ./internal/models/` | **0** | `ok …internal/models 20.091s` |
| 4 | `go test -count=1 -timeout 30m ./cmd/wisp/` | **1** | `FAIL …cmd/wisp 377.780s`；`RUN=224 / 判决=223 / PASS=222 / FAIL=1` |
| 5 | `sh scripts/d22scan.sh` | **0** | `d22scan: clean - no D22 ban violations`（`examined 289 production Go files`；种子正控 `PASS=35 FAIL=0`） |
| 6 | `sh scripts/check-path-length-budget.sh` | **1** | 见 6.2（红行与本票无关） |

**6.1 第 4 把按名作差**（`rc=1` 今天不能当判据，`A1079`／`A1100`）。红名册**逐名、1 枚**：`TestAC14GoSideEvalPushReachesThePage` —— **既有红、NOT 新增**：`docs/reports/pending-and-issues.md:20471` 逐字记着"`cmd/wisp` 整包 `rc=1` 有**既有红 1 枚** `TestAC14GoSideEvalPushReachesThePage`（与本票无关）"；本发该枚红行逐字 `title=""` want `PUSHED-33R5-OK`（真窗族、本机 WebView2 Eval push 维度）=> **新增红 = 0**。死点（RUN 无判决行、逐发登记）：`TestPanelThreadIsSTAAndExitsCleanly` 只有 `=== RUN`、无 `--- PASS/FAIL`；与 `A1100` 记的"每发恰 1 枚 RUN 无判决行、死点漂"同族，本发落在这一枚（`377.780s` 收尾、早于 `-timeout 30m`）。本票新增 6 枚 `Test357*`（`cmd/wisp`）在本发内**全绿**、逐名：`Test357StreamNoticeIsTheOneSentence` / `Test357StoreRootNamingIsNotWhatTheReportUses` / `Test357VerifyTailIsStillOnThePassingPath` / `Test357VerifyTailReportsTheVerifiedTreeNotTheStoreRoot` / `Test357VerifyTailSaysOneSentenceAndKeepsRC0` / `Test357VerifyTailWithoutAStreamPrintsNothing`。
**6.2 第 6 把 RED 的红行**：`.scratch/wisp/issues/181-…-done.md`（相对 159 字符、名 138）与 `…/275-…-done.md`（相对 126、名 105）"over budget 且不在名册" ＋ `181` 那条"名册里有、树里没了" ＋ 计数守卫（名册 60 / 树 61）。四条**全部落在 `.scratch/wisp/issues/**`**（票 181／275 文件名与名册漂移）；本腿没碰这两个路径、也没向 `.scratch/` 新增任何 tracked 路径；本腿新增 8 条路径**全部不在** over-budget 的 61 条里。

## 7 格式门禁（甲乙两尺都报、甲尺判据）

口径＝`docs/evidence/s1/pipe-format-gate-reads-tracked-blob-not-worktree.md`。**甲尺（入库 blob、判据）**：`git archive HEAD | tar -x` 到仓外后对 8 枚文件跑 `/d/work/base/gopath/bin/gofumpt.exe -l` => **0 枚**。**乙尺（工作树、只作诊断）** => **1 枚**：`cmd/wisp/models.go`（本机 `core.autocrlf=true` 的**已知假红**：该件工作树 353 CR、入库 blob 0 CR；改前同色、本腿没让它变）。本腿新加的两枚 `cmd/wisp` 件初版真被 gofumpt 咬一格（`"…tell "+` 要写成 `"…tell " +`），已按意见改掉、改后 0 命中。

## 8 提交（逐枚 ＋ 逐枚 `git show --name-only` 的 `.go` ／越界；正控＝同尺打一笔含 `.go` 的真提交）

**一枚提交：`aeea4e5b`**（`票 357 落地腿 357-r1：verify 报"与已验签清单一致"时明说清单未点名的备用数据流（rc 仍 0）`）。
`git show --name-only aeea4e5b` 逐字 8 条、**全 `.go`**：`cmd/wisp/models.go`、`cmd/wisp/models_verify_ads_357_windows_test.go`、
`internal/models/ads.go`、`internal/models/ads_357_other_test.go`、`internal/models/ads_357_windows_test.go`、
`internal/models/ads_other.go`、`internal/models/ads_windows.go`、`internal/models/downloader.go`。
**越界 0 条**（同一枚尺判 `^(internal/models/|cmd/wisp/)`：零 `OUT` 行）；`git show --stat` 读数 `8 files changed, 1016 insertions(+), 3 deletions(-)`，
创建 6 枚新件、修改 2 枚（`models.go`、`downloader.go`）。**正控（同尺打一笔含 `.go` 的真提交）**＝`7835d67f`
（票 359 落地腿）：同尺量出 total=2、`.go`=1、非 `.go`=1（`docs/evidence/s1/359-r1-impl.md`）=> **尺是活的**，不是恒报 0。
三处还原核（提交版 / 还原台 / `git show aeea4e5b:` blob）：`internal/models/ads.go` `73350E3D…F5` 三处同串、
`ads_windows.go` `7D4102D4…5F` 三处同串；`cmd/wisp/models.go` 三处同串前须把工作树 CRLF 归一成 LF（本机 `core.autocrlf=true`
的检出形态，入库 blob 是 389 LF／0 CR），归一后三处同为 `B0FE3890DAC04C4FC387EFD011A947C8B44D47EA039AB98101D882FA5366D603`。
甲尺口径与 sec 7 一致：`git archive HEAD` 到 `D:\tmp\wisp\20261021\357r1b\archive\` 后跑 gofumpt，**0 枚**（`JIA-EXIT=0`，输出零行）。

## 9 纪律格

写面只碰 `internal/models/**`（4 枚）＋`cmd/wisp/**`（2 枚）＋本件。NOT 碰 `internal/{winsec,proc,risk,agent,tools,config}/**`、`tools/d22scan/**`、`allowlist.txt`、`thresholds.go`、任何 golden、`go.mod`／`go.sum`、`.github/**`、`frontend/**`、`design/**`、`docs/{PLAN.md,specs,reports}/**`。NOT 改哈希／pin／seal 判定：`VerifyDir`／`verifyNothingUnnamed`／`verifyFileHash` **一字节未动**；`resolveInstalledDir` 是把 `VerifyInstalled` 原有三行搬位置、**实参不变**。NOT 为绿改既有测试：既有 `.go` 用例**一枚未改**（`git status` 里只有**新加**的两枚用例件）。NOT 翻框／改票名／`Status:`／`-done`／改台账／push。临时件只建不删：`D:\tmp\wisp\20261021\357r1b\**`（`bin/` `logs/` `probe3/` `dataG/` `dataD1/` `dataOVR/` `override-a1/` `*.pristine`），未产 `.out`。未 `add -A`／`.`；未 amend／reset／rebase／stash／checkout ./clean。
