# 票 360 终裁腿 `360-v1` 裁决件 — part2/2（门禁／红名册／写面／`attrib` 分解／汇总）

- 本件接 `360-v1-verdict.md`（§1–§3：前提新鲜度／`AC#2` 判语／零产码复核）。
- ★**结论仍是：`AC#2` ⛔ 翻**（机理与补法见 part1 §2.3／§2.4）。

## 4. 门禁复跑（我自己现跑、⛔ 抄；`GOFLAGS=`／unset `GOTMPDIR TMPDIR TEMP TMP`／未设 `CGO_ENABLED`）

| 门 | 命令 | rc | 末行逐字 |
|---|---|---|---|
| build | `GOFLAGS= go build ./...` | **0** | （无输出） |
| vet | `go vet ./internal/config/` | **0** | （无输出） |
| test | `go test -count=1 -v ./internal/config/` | **0** | `ok  github.com/CarlosShao/wisp/internal/config  2.085s` |
| d22scan | `sh scripts/d22scan.sh` | **0** | `d22scan: clean - no D22 ban violations; live scope work: bans #1-5 internal/=247, bans #1-5 cmd/=42, …` |
| pathlen | `sh scripts/check-path-length-budget.sh` | **0** | `check-path-length-budget.sh: VERDICT GREEN - every over-budget tracked path is rostered by name with a reason, and the roster equals the tree` |

★**四数分开报（⛔ 相减）**：`RUN=196`／`PASS=115`／`FAIL=0`／`SKIP=0`，另 `^panic:` **0**。
★pathlen 分母：`tracked paths=17572  over-budget=61  covered by roster=61  not in roster=0`。

### 格式门禁：甲乙两尺都报（裁定件 `pipe-format-gate-reads-tracked-blob-not-worktree.md`）

- **甲＝入库 blob 尺（★判据）**：`git archive HEAD` 到仓外 ⇒ **1156 枚 `.go`** 跑 `gofumpt v0.12.0 (go1.27.1) -l` ⇒ **列 20 枚**。
  - ★**本腿新增件 0 命中**（`grep -c writeguard_adopt_360_test.go` ⇒ `0`）。
  - ★**20 枚全在 `.scratch/wisp/probes/**`**（`probes_count=20`、`not_probes=0`）。
  - ★**对照父提交 `d68ddfae^`＝`1cfe94bb` 导出树 ⇒ 同样 20 枚、`diff` 逐字相同 ⇒ IDENTICAL** ⇒ **本腿新增 0 枚**。
- **乙＝工作树尺（★只作诊断）**：`git ls-files '*.go'` 1156 枚 ⇒ **列 35 枚**（⛔ 它报的 55——★我实测 **35**；55 那个数它取自**文件系统遍历**，见 §8）。
  - 单文件直跑 `gofumpt -l internal/config/writeguard_adopt_360_test.go` ⇒ **空、`rc=0`**；本腿新增件 **0 命中**。
  - `comm` 两向：**只在乙＝15、只在甲＝0** ⇒ 甲绿乙红的差**全是**检出副本假红（§7 逐枚证）。

## 5. 红名册逐名（⛔ 只报枚数）

**本票写面 `internal/config` 内：新增红＝0。** 整包 `ok … 14.304s`（主检出全跑）、`RUN=196 PASS=115 FAIL=0 SKIP=0`。

**全仓（主检出 `go test -count=1 ./...`，`MAIN_RC=1`）**——11 枚逐名，★**无一枚属本票写面、无一枚被本腿区间碰过**：

1. `TestAC1WindowsPrivateFileThroughASwappedAncestorNeverTruncatesAForeignFile342r1`
2. `TestAC1WindowsPrivateFileThroughASwappedAncestorNeverTruncatesAForeignFile342r1/symlink-dir`
3. `TestAC2WindowsPrivateFileThroughASwappedAncestorNeverDeletesAForeignFile342r1`
4. `TestApprovalCardViewJSONKeysMatchFrontendTypes`
5. `TestC21DesignTokensFourWayAgree`
6. `TestC21TableColourRowsMatchTokensCSS`
7. `TestC26RewriteAccountIsConsumedAtEverySecurityLeg`
8. `TestComposerContractTypesMatchFrontend`
9. `TestResolvePerCallBudget`
10. `TestStreamLogDroppedNamingLedgerIsNotBounded35r8`
11. `TestStreamLogFloodBelowKeyBoundIsNotBounded35r8`

★**归因（我另起尺）**：`git diff --name-only d68ddfae^..7c42145c -- internal/ball internal/panel internal/risk internal/winsec cmd/wisp` ⇒ **全部 untouched**。
★`TestResolvePerCallBudget` 是**抖**（全跑红、单跑 `-run` **3/3 `RC=0`**），文件 `internal/risk/pathresolver_budget_norace_test.go` 在父提交与 HEAD **同存、未被本腿动**。
★**导出树上的红名册⛔ 当判据**：两棵导出树都没有 `.git`，`internal/tools` 那 8 枚 `TestTicket252*`／`TestReadGitOnThisRepositoryIsSelfConsistent` 在导出树上**假红**、在主检出 **`rc=0`** ⇒ ★**必须用主检出读数**（这本身是一条可复用的教训）。

## 6. 写面与纪律 —— 成立

逐笔 `git show --name-only`：

| 提交 | 文件数 | `.go` | 越界 |
|---|---|---|---|
| `d68ddfae` | 1（`internal/config/writeguard_adopt_360_test.go`） | **1** | 0 |
| `fc7c85cc` | 1（同件） | **1** | 0 |
| `0d907103` | 1（`docs/evidence/s1/360-r1-impl.md`） | 0 | 0 |
| `7c42145c` | 1（同件） | 0 | 0 |

- 区间并集 `git diff --name-only d68ddfae^..7c42145c` ⇒ **仅 2 枚**（用例＋交付件）；越界判据（`internal/config/*_test.go` ∨ `docs/evidence/s1/*`）⇒ **OUT_OF_SCOPE=0**。
- **正控尺**（同尺打含 `.go` 的真提交 `ecebc7a9`）⇒ 文件 **4**、`.go` **2**（`cmd/wisp/always_write_no_clobber_{226,227}_test.go`）⇒ **尺不空转**。
- **负控** `5ae6716c` ⇒ 文件 2、`.go` **0**。
- 编排者 `1cfe94bb` ⇒ 1 枚、`.go` **0**（⛔ 算本腿）。
- **⛔ 改票面**：`git log --all -- '.scratch/wisp/issues/360-*.md'` ⇒ 4 笔**全为编排者**、实现腿 **0 枚**；`- Status:` 在 HEAD 与工作树**同为** `ready-for-agent`、无 `-done`。
- **⛔ push**：`origin/dev..HEAD` ＝ **398**（非空）⇒ 全在本地。
- ★本腿自身：`add`／`commit` 两处**显式 pathspec**、⛔ `add -A`／`.`、⛔ amend／reset／rebase／stash／checkout ./clean、⛔ 翻任何框、⛔ 落任何 `.go`。

## 7. ★`attrib.sh --tracked-only` 35 枚分解的独立复核 —— 成立（分解对，且我补了三处它没证的）

`sh .scratch/wisp/probes/161/r5/attrib.sh --tracked-only` ⇒ **`rc=1`**、末行逐字：
`attrib.sh: RED (tracked-only) - tracked-dirty=35 files_dirty=35. The CI gofmt step and this step read the same denominator; this one additionally refuses to answer when handed 0 files.`
其中 `(A) tracked tree … tracked .go files handed to it: 1156`。

**我的独立分解（三把尺逐枚，⛔ 采信它）**：

| 步 | 读数 | 判 |
|---|---|---|
| `attrib` 名单 | **35** | — |
| 甲尺（blob 树）名单 | **20** | — |
| `comm -23`（只在 attrib） | **15** | ✅ |
| `comm -13`（只在 blob） | **0** | ✅ ⇒ 两尺在"真脏"上**完全一致** |
| 那 20 枚是否全在 `.scratch/wisp/probes/**` | `probes_count=20`／`not_probes=0` | ✅ **它那句成立** |
| 那 15 枚**blob 级**是否干净 | ★**逐枚 `git show HEAD:<p> \| gofumpt -l` ⇒ 15/15 `CLEAN_BLOB`** | ✅ **这是它没做的那把尺** |
| 那 20 枚 blob 级是否脏 | 逐枚 ⇒ **20/20 `DIRTY_BLOB`** | ✅ |
| 那 15 枚 `HEAD:` blob CR 数 | **全 0** | ✅ |
| 那 15 枚**工作树** CR 数 | `185/176/174/188/185/185/185/176/612/612/612/389/131/225/1115` | ✅ 全 >0 |
| `git ls-files --eol` 那 15 枚 | 逐枚 **`i/lf  w/crlf  attr/text eol=lf`** | ✅ ★**它没做的那把尺** |
| `git status --porcelain` 那 15 枚 | 逐枚 **CLEAN**（git 认为无改动） | ✅ |
| `core.autocrlf` | `true`（`.gitattributes` 另有 `*.go text eol=lf`） | ✅ |

⇒ ★**35 ＝ 20 枚真脏 blob ＋ 15 枚 CRLF 检出副本**：**分解成立**，且**两把尺在"真脏"上一致**（`comm -13` 空）。
★**那 20 枚全在 `.scratch/wisp/probes/**`**：**成立**（`not_probes=0`）。
★它点名的 4 枚产码（`cmd/wisp/models.go`／`internal/agent/approval/pending_read.go`／`internal/agent/tools.go`／`internal/risk/provenance.go`）⇒ **我核实在那 15 枚里、blob 级 `CLEAN`、`i/lf w/crlf`** ⇒ ★**正是 CRLF 假红、⛔ 在甲尺名单里** ⇒ **它那句也对**。
★**归属**＝票 362（20 枚真脏）／票 361（15 枚 CRLF）⇒ ⛔ 本票。★**这 35 枚无一属本票写面、父提交对照 0 增量**。

## 8. 一处与它不一致的读数（照实记账）

- 它报「乙＝工作树尺整树 **列 55 枚**」。★**我用 `git ls-files '*.go' | xargs gofumpt -l` 实测列 35 枚**（与 attrib 的 35 同集）。
- ★**差从哪来**：它的 55 取自**文件系统遍历**（`find . -name '*.go'`）。我实测同一把尺遍历得 **1558** 枚 `.go`（含**未跟踪**的 `.scratch` 台件）⇒ 列 **50** 枚。⇒ ★**55 与 50 都不是"跟踪集"读数**；**跟踪集的正确读数是 35**。
- ⇒ 这不改甲的 20、不改 35＝20＋15、也不改裁定；★但"乙尺 55"这个数**在跟踪集口径下不成立**，下一腿引用时⛔ 抄 55。

## 9. 逐项判语汇总

| # | 项 | 判 |
|---|---|---|
| 1 | 票面前提新鲜度（HEAD blob 复量） | ✅ 全部成立 |
| 2 | ★**`AC#2` 本票核心** | ★**⛔ 翻**（两态真、但正控可整臂连断言摘除仍绿；`M3 恒真仍绿` 的解释**不成立**） |
| 3 | 零产码（五枚＋我加 17 枚／全仓） | ✅ 成立 |
| 4 | 门禁五枚 rc | ✅ 全 0（`RUN=196 PASS=115 FAIL=0 SKIP=0`、`panic=0`） |
| 5 | 格式甲乙两尺 | ✅ 甲 20（本腿 0 增量、与父同）；乙 35（诊断、全 CRLF 假红） |
| 6 | 红名册 | ✅ 写面新增红 **0**；全仓 11 枚逐名、全不属本票 |
| 7 | 写面与纪律 | ✅ 越界 0、正控不空转、⛔ 改票面、⛔ push |
| 8 | `attrib.sh` 35 分解 | ✅ 35＝20＋15、20 全在 probes、15 全 CRLF（我补 blob 级与 `--eol` 两把尺） |

★**一句话**：`360-r1` 的**读数是真的、零产码是真的、门禁是真的、诚实也是真的**；★**但它那枚用例的两态正控不满足 `AC#2` 判据的"改恒真"那一半**——`M1＋M3′` 同时摘臂＋恒真**仍绿** ⇒ **`AC#2` ⛔ 翻**，补法见 part1 §2.4（仍〔甲〕形、零产码）。
