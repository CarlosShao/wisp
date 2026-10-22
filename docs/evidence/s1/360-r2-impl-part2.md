# 票 360 补强腿 `360-r2` 交付件 — part2/2（红名册与同树对照／写面尺／未做边界／汇总）

- 本件接 `360-r2-impl.md`（§1–§7：开工前／补了哪一行前置断／`M1＋M3b` 判据／其余变异发／三读数与 `sha256`／零产码／门禁）。
- ★结论仍是：**补强成立**（`M1＋M3b` 与更狠的 `M1＋M3a＋M3b` 两发都红在 `:84`）。

## 8. 红名册（★取主检出，⛔ 取导出树）

**本票写面 `internal/config`：`ok … 25.229s`、`RUN=196 PASS=115 FAIL=0 SKIP=0` ⇒ 新增红 0。**

**全仓主检出 `go test -count=1 ./...` ⇒ `MAIN_RC=1`，12 枚逐名**：

1. `TestTenOpsInOneToolCallGetOneConfirm`
2. `TestC21TableColourRowsMatchTokensCSS`
3. `TestApprovalCardViewJSONKeysMatchFrontendTypes`
4. `TestStreamLogFloodBelowKeyBoundIsNotBounded35r8`
5. `TestStreamLogDroppedNamingLedgerIsNotBounded35r8`
6. `TestComposerContractTypesMatchFrontend`
7. `TestC21DesignTokensFourWayAgree`
8. `TestResolvePerCallBudget`
9. `TestC26RewriteAccountIsConsumedAtEverySecurityLeg`
10. `TestAC1WindowsPrivateFileThroughASwappedAncestorNeverTruncatesAForeignFile342r1`
11. `TestAC2WindowsLeafIdentityReCheckRefusesAForeignNameAndNeverDeletesIt342r1`
12. `TestAC1WindowsLeafRecheckAnswersOnTheBareSiblingForm342r2`

### 8.1 ★★同树对照（我另起的一把尺，⛔ 派单要求但必需）

只说"这些红不在 `internal/config`"不够——**必须证明我的改动不是它们的原因**。所以我做**同树对照**：
**同一棵主检出、同一时刻**，只把 `internal/config/writeguard_adopt_360_test.go` **换回 HEAD blob**（`sha256 481fb162…62c7e`、93 行），跑同一发 `go test -count=1 ./...`，然后**立刻换回我改后的件**（`sha256 52a1aaee…eeb8`、107 行，逐字核过）。

| 发 | 台 | 红枚数 | 名册 |
|---|---|---|---|
| **我改后（第 1 发）** | 主检出、本腿件在位 | **12** | 见上 12 枚（`logs/17-red-roster.txt`） |
| ★**同树对照（HEAD 件）** | **同一棵主检出**、件换回 93 行 HEAD blob | **9** | `TestAC1WindowsPrivateFileThroughASwappedAncestor…342r1`／`TestApprovalCardViewJSONKeysMatchFrontendTypes`／`TestC21DesignTokensFourWayAgree`／`TestC21TableColourRowsMatchTokensCSS`／`TestC26RewriteAccountIsConsumedAtEverySecurityLeg`／`TestComposerContractTypesMatchFrontend`／`TestResolvePerCallBudget`／`TestStreamLogDroppedNamingLedgerIsNotBounded35r8`／`TestStreamLogFloodBelowKeyBoundIsNotBounded35r8`（`logs/25-fullbase-roster.txt`） |
| **我改后（第 2 发）** | 主检出、本腿件在位 | **12** | `logs/27-mine-run2-roster.txt` |

★**两发"只在我有"的 3 枚，两发之间自己就换人了**：

| 发 | 只在"我改后"出现的 3 枚 |
|---|---|
| 第 1 发 | `TestAC1WindowsPrivateFileThroughASwappedAncestor…342r1`（★这一枚**在对照组里也红**）／`TestAC2WindowsLeafIdentityReCheck…342r1`／`TestTenOpsInOneToolCallGetOneConfirm` |
| 第 2 发 | `TestAC2WindowsLeafIdentityReCheck…342r1`／`TestAC2WindowsPrivateFileThroughASwappedAncestorNeverDeletesAForeignFile342r1`／`TestTenOpsInOneToolCallGetOneConfirm` |

★**三把独立的尺都说这 3 枚与我无关**：
1. **漂移**：同一棵树上跑两次，"多出来的 3 枚"**不是同一组**（第 1 发的 `…NeverTruncates…` 第 2 发不红了、第 2 发换成 `…NeverDeletes…`）⇒ 是**抖**，⛔ 我的改动引入的确定红。
2. **单发复跑全绿**：把那 3 枚各自 `-run` 单跑 ⇒ `…LeafRecheck…342r2` **`rc=0`**、`…LeafIdentity…342r1` **`rc=0`**、`TestTenOpsInOneToolCallGetOneConfirm` **`rc=0`**（件 `logs/` 末尾三发）。
3. ★**结构上不可能**：`writeguard_adopt_360_test.go` 是 **`package config` 的包内测试件**（`go list -f '{{.TestGoFiles}}'` 里可见），**只会被编进 `internal/config` 的测试二进制**，⛔ 被 `internal/winsec`／`internal/agent/approval` 链接；且 `go list -deps ./internal/winsec/` 里**根本没有 `internal/config`**。
   ⇒ ★**本腿这一枚测试件的字节，物理上不可能改变那两个包的任何一枚用例的判**。
4. ★**且 `internal/config` 包本身在同树对照两发里都是 `ok`**（`18.725s`／`17.571s`／`25.229s`）⇒ **写面新增红 0**。

★**与 `360-v1` 那 11 枚不同集**：照实记账——全仓整跑的红名册**本来就逐发漂移**（`360-v1` 自己也把 `TestResolvePerCallBudget` 记成"抖"）。★**所以判据只能取"同树对照差"与"单发复跑"，⛔ 取"两发名册逐字相同"**。

## 9. 写面尺与纪律

- 本腿提交逐枚 `git show --name-only` 见 §11；**`.go` 越界 0**（写面＝`internal/config/writeguard_adopt_360_test.go` ＋ `docs/evidence/s1/360-r2-impl*.md`）。
- **正控尺**（同尺打一笔含 `.go` 的真提交 `ecebc7a9`）⇒ 文件 **4**、`.go` **2**（`cmd/wisp/always_write_no_clobber_{226,227}_test.go`）⇒ **尺不空转**。
- **负控**（纯文档笔 `5ae6716c`）⇒ `.go` **0**。
- ⛔ push、⛔ 翻任何框、⛔ 改票名／`Status:`／`-done`／台账、⛔ amend／reset／rebase／stash／checkout ./clean；`add`／`commit` 两处**显式 pathspec**、新文件**逐个文件名字面点名**（⛔ glob）。

## 10. 未做／边界（照实记账）

- ★**判据原文未改**：`AC#1` 裁定文字与判据原文**一字未动**（那是"改契约文字"、属摆机主的三类之一）⇒ 本腿只**补强实现**。
- ★**M3b 单独恒真仍绿**已照实登记（part1 §4），机理写清、⛔ 藏、⛔ 算成红。判据那一半由 **M1** 与 **M1＋M3b** 两发共同满足。
- ★**⛔ 动产码**：若"让前置断不可被吞"做不到才停手；**做到了**（part1 §3.1 实测：摘臂 ＋ 全部别的断言弄哑，仍红在 `:84`）⇒ 未触发派单停手上报条①。
- `cmd/wisp` **整包读数⛔ 当判据**（`A1126`／`A1127`）⇒ 本腿未用它作判据。
- ★**格式门禁第三把仪器**（`attrib.sh --tracked-only` 的既有 35 枚）**未修、未隐瞒、未算进本票**：那 35 枚无一属本票写面、本腿新增件在两尺上均 0 命中 ⇒ 归属＝票 362（20 枚真脏）／票 361（15 枚 CRLF），⛔ 本票。

## 11. 提交逐枚与越界

| 提交 | 主题 | 文件数 | `.go` | 越界 |
|---|---|---|---|---|
| `095699f2` | 票 360 续跑腿 360-r2: 把"该 Adopt 臂真的跑过"独立成一条不可被同臂断言吞掉的前置断 | **1**（`internal/config/writeguard_adopt_360_test.go`） | **1** | **0** |
| `（见 §11.1）` | 票 360 续跑腿 360-r2: 交付件（补强读数／判据红行／sha256／门禁） | **2**（`360-r2-impl.md`、`360-r2-impl-part2.md`） | **0** | **0** |

★§11.1 的提交号由本腿在提交后现填（填法同 `360-r1`：逐笔 `git show --name-only` 现核）。

- 越界判据＝`internal/config/*_test.go` ∨ `docs/evidence/s1/*` ⇒ **OUT_OF_SCOPE=0**。
- `git rev-parse 095699f2:internal/config/writeguard_adopt_360_test.go` ⇒ `7077bc2b…` ＝ `git hash-object` 工作树同串 ⇒ **提交内容与工作树逐字相同**。
- **正控尺**（同尺打 `ecebc7a9`）⇒ 文件 **4**、`.go` **2** ⇒ **尺不空转**；**负控** `5ae6716c` ⇒ `.go` **0**。
- **⛔ push**：本腿全部提交只在本地。
- **⛔ 改票面**：本腿零票面提交（`git log --all -- '.scratch/wisp/issues/360-*.md'` 里实现腿 0 枚）。

## 12. 一句话汇总

★**补强成立**：`ranAdopt` 哨兵 ＋ 一条独立前置断（原件 `:93` 的 `if`、`:94` 的 `Fatalf`），**放在臂断言 `:98` 之前**。
★**判据那一发（`M1＋M3b`：同时摘臂 ＋ 臂断言改恒真）本腿自跑 ⇒ `RC=1`、红行 `writeguard_adopt_360_test.go:84`**（★`:84` 是**摘臂后**的号：删掉 `:75-84` 那 10 行，前置断的 `Fatalf` 由原件 `:94` 上移到 `:84`）；上一枚在同一发上 `RC=0` 仍绿。
★**更狠的 `M1＋M3a＋M3b`（摘臂 ＋ 把别的断言全部弄哑）⇒ 仍红在 `:84`** ⇒ 这条前置断**不可被同臂断言吞掉**（实测、⛔ 声称）。
★**原样仍绿**：`file=60 mem=56` ⇄ `mem=60 file=60 report.Hot=[ball]`，两态读数**与补强前逐字相同**。
★**零产码**：`internal/config` 全部非测试 `.go` 逐枚 `rev-parse` 对 `hash-object` ⇒ **全 SAME**。
★**门禁**：`build 0`／`vet 0`／`go test ./internal/config/` `0`（`RUN=196 PASS=115 FAIL=0 SKIP=0`、`panic=0`）／`d22scan 0 clean`／`pathlen 0 GREEN`（`17575/61/61/0`）。
