# gw-1 味③＝票 353 `AC#3` winsec 两平台成对色（★同一把尺复量·⛔ 353 身份＝本腿是⛔ 353 的独立腿码）

尺＝`sh/05-run-pair.sh`（定向 `./internal/winsec/` `-count=1 -v -timeout 30m`）＋`sh/15-winsec-remeasure.sh`（撞红三发＋并集）＋`sh/16-five-verdicts.sh`（具名逐枚色）。
两对锚（★量的＝`git rev-parse <锚>:<path>` ⇄ 树非空尺）＝
- **对甲**＝基线侧 `gwW`＝锚 `3a509e50`（＝票 353 乙本体落地前那一笔＝笔①仪器已入库、乙⛔ 落）⇄ 落地侧 `gwL`＝锚 `a19d5bf2`（＝今天 HEAD）。
- **对乙**＝基线侧 `gwB`＝锚 `0fb399c3`（＝`ecffb93b^`，＝票 352 那一行的父，⛔ 353 仪器也未入库）⇄ 同 `gwL`。
`gwW .go=778` ⇄ `gwB .go=777` ⇄ `gwL .go=781`；`nametopup_353_windows_test.go` 在位尺＝`gwW PRESENT`／`gwB ABSENT`／`gwL PRESENT`。
★⛔ 整包 `./cmd/wisp`、⛔ `./...`、⛔ 共享工作树跑 `go test`；全部 8＋6 发 `SHAPE=OK`／`VOIDKIND=none`（＝台损 0 枚）。

## 一、四数（`RUN/PASS_TOP/FAIL_TOP/SKIP`，★⛔ 混用⛔ 相减）

对甲（本味判据那一对）：

| 发 | 基线侧 `gwW` | 落地侧 `gwL` |
|---|---|---|
| host 首发 s1/s2 | RC=1 `121/67/3/0` | RC=1 `121/67/3/0` |
| host 复量 1 | RC=1 `121/67/3/0` | RC=1 `121/68/2/0` |
| host 复量 2 | RC=1 `121/67/3/0` | RC=1 `121/68/2/0` |
| host 复量 3 | RC=1 `121/67/3/0` | RC=1 `121/68/2/0` |
| container s1/s2 | RC=0 `74/46/0/0` | RC=0 `74/46/0/0` |

对乙：host `119/65/3/0` ⇄ `121/69/1/0`；container `74/46/0/0` ⇄ `74/46/0/0`。
★`RUN 119 ⇄ 121` 那 2 枚⛔ 是判据＝名集 symdiff 尺（`^=== RUN` 名集，量的）＝`ONLY_LAND=2`／`ONLY_BASE=0`＝
`TestAC2WindowsNameTopUpBeforeTheReCheckLeavesTheTopUpObjectNamed353r1` ⇄ `TestAC2WindowsRemovalStillRunsWhenNobodyTopsUpTheName353r1`
＝正是 353 那枚件的两支仪器（`gwB` 里该件 ABSENT）⇒ 名集作差解释得动枚数差、⛔ 拿四数差当"改了 N 枚测试"。

## 二、逐名红名册 ⇄ 并集作差（★判据＝并集 comm，⛔ 单发换位）

四发 host 并集（＝`p353` 首发 + `ws3` 三发）：

| 名 | 基线侧 `gwW` | 落地侧 `gwL` | 定性 |
|---|---|---|---|
| `TestAC2WindowsNameTopUpBeforeTheReCheckLeavesTheTopUpObjectNamed353r1` | **RED（4/4 发都红）** | **绿（4/4 发都绿）** | ★票 353 的仪器＝修前红⇒落地绿＝这一味今天真在（本腿独立复量到） |
| `TestAC2WindowsRemovalStillRunsWhenNobodyTopsUpTheName353r1`（正控） | 绿 | 绿 | ★两侧都绿＝"删除照旧执行"那把牙没被摘（⛔ 借它圆判） |
| `TestAC1WindowsPrivateFileThroughASwappedAncestorNeverTruncatesAForeignFile342r1` | RED（**3/4** 发） | RED（**4/4** 发） | 票 342 族既有红＝两侧都在册（★本腿⛔ 据"落地那侧更多"说任何方向＝四发逐名摆着交裁腿） |
| `TestAC2WindowsLeafIdentityReCheckRefusesAForeignNameAndNeverDeletesIt342r1` | RED（**2/4** 发） | RED（**1/4** 发） | 342/348 racer＝公共块在册摆动族 |
| `TestAC1WindowsLeafRecheckAnswersOnTheBareSiblingForm342r2` | RED（**2/4** 发） | RED（**2/4** 发） | 同上（摆动） |
| `TestAC2WindowsPrivateFileThroughASwappedAncestorNeverDeletesAForeignFile342r1` | RED（**1/4** 发） | RED（**2/4** 发） | 同上（摆动） |

⇒ 名⇄发红/绿逐枚底片＝四枚红名册文件 `logs/p353-s1-gwW-winsec.rednames.txt`⇄`logs/ws3-{1,2,3}-gwW.raw.txt.rednames.txt` ⇄ 落地侧同四名（发序＝首发、复量 1、2、3；本表那 eight 个分数＝逐枚数 `grep -c` 量的，⛔ 算的）。

⇒ **新增红（只在落地侧、并集尺）＝0 枚**；★消失红＝**1 枚**＝正是 353 自己那枚仪器（`3a509e50` 侧红⇄`gwL` 侧全绿）＝方向＝落地更干净、且它是票 353 `AC#2` 要买的那一味，⛔ 本腿翻框、⛔ 据此盖章。
⇒ container 档两侧逐名红册＝**各 0 枚**（`FAIL_TOP=0` 两侧同、名集 `74/74` symdiff 0）＝POSIX 半边干净。

## 三、摆动／弃权名册（照公共块逐族点名＋本腿现量命中）

| 族 | 本腿射程内现量 | 处置 |
|---|---|---|
| `internal/winsec` 两枚 `//go:build !windows` 仪器（`openleaf_322_other_test.go` ⇄ `placement_symlink_113_other_test.go`） | host 档两侧 **ABS**（⛔ 编译）＝本腿尺读出 `ABSENT`；container 档两侧 **逐枚 PASS** | ★既⛔ 记绿也⛔ 记弃权＝票 322 `AC#1` 逐字那形；两档合起来每枚都在它该编译的那档量到色 |
| winsec 名集 host⇄container 面差（同一把尺 `^=== RUN` 名集） | host 121 枚 ⇄ container 74 枚；`only_container=52`／`only_host=99` | 平台档结构性差（tag 面）＝⛔ 任何一侧的"枚数大"当"跑得多" |
| `…116` 族 4 枚 | 具名尺打在 winsec（host ⇄ container 落地发）＝**0 命中**；同尺正控读 1＝尺⛔ 瞎 | ★真身＝`internal/risk/syncdirs_ancestor_actable_leg_116_test.go` 四枚（`grep` 现量＝`TestSyncAncestorActableLegFailsClosedOnMovedAncestor116` 等 4 名）＝**⛔ 本腿三味的射程**（⛔ 跑 `./internal/risk/`）＝缺哪一味具名报，⛔ 圆判 |
| `TestTicket223ModeLooseningChangesTheRunningModeAfterAllow` | 命中 0（同尺正控先行） | ★真身＝`cmd/wisp/config_reload_223_test.go:535`＝住在**派单明令⛔ 跑的整包 `./cmd/wisp`** 里⇒本腿⛔ 量其色＝票 353 那味（＝"已定性仪器形状"）由编排者那窗取；本腿只登记它不在我射程 |
| 机主常驻 `wisp.exe` 占 `Ctrl+Alt+Q/M/P` 热键那一族 | `Get-Process -Id 21892`＝**在场**（`ProcessName=wisp`，StartTime 2026-10-15 19:03:14，量的）；winsec 两档红名册里**无**任何热键族名（尺＝逐名册 `logs/*.rednames.txt` 逐枚列全） | ★热键族＝弃权、⛔ 归因给字节；本腿⛔ 拿它任何一枚的绿当凭据 |

## 四、`private_fail_test.go:37/50/61/72/75` 五枚判定（★对入库 blob 逐枚取那一行，量的）

五枚行号尺＝`git cat-file blob HEAD:internal/winsec/private_fail_test.go \| sed -n '37p;50p;61p;72p;75p'`＝
`:37`／`:50`＝`if !errors.Is(err, ErrNotSealable) {`、`:61`＝`} else if !errors.Is(err, ErrNotSealable) {`、`:72`＝`if err := SealFile(p); !errors.Is(err, ErrNotSealable) {`、`:75`＝`if err := SealDir(dir); !errors.Is(err, ErrNotSealable) {`；
同尺打在 `0fb399c3` 的 blob＝**逐枚同串**＝两侧测试字节同一（`private_fail_test.go` blob `e95d48aa` 三锚同串）。逐枚两侧色＝**八发全 PASS**（含 4 枚子测试名 ⇄ 顶层），详见 `02` 第四节（★具名一处⛔ 圆判＝票面算五枚判定、仪器只印四枚子测试名，`:72`⇄`:75` 同住 `/seal_file_and_dir_direct` 一支）。

## 五、味③ 锚⇄时点＋台面

对甲⇄对乙 全部 8 发起手锚＝`a0124803`（23:25）⇄ 该批结束 `4fb6c088`（23:30 交件尺）；复量 6 发起手⇄结束＝**同串 `4fb6c088`**（未漂）。
★期间 HEAD 漂了 3 笔（139-v1/353-v1 的判语⇄尺册件）⇒ 本腿重跑尺＝量的：`git diff --name-only a19d5bf2..HEAD -- "*.go"`＝**0 枚**，且 `internal/winsec` 的 **tree hash** 在 `a19d5bf2` ⇄ `4fb6c088` 同串（`68e765de…`＝`SAME=YES`）⇒ ★读数所依字节与我锚定的那次 HEAD 逐目录同一＝⛔ 任何一发需要因漂移重跑（这味＝今天两枚裁腿栽的那味，本腿守住）。
台面：`0xc0000135`=0／`build failed`=0／`walk`=0／`setup failed`=0／`no test files`=0／`SHAPE=VOID`=0；台根＝`%LOCALAPPDATA%\wisp-test-gotmp\gw1-*`（★⛔ `D:\tmp\wisp`；`icacls` 现量该根＝三条 `(I)(OI)(CI)(F)`、无 `BUILTIN\Users` 继承）；DLL 段＝`cygpath -w` 两枚前缀进 `PATH`。
