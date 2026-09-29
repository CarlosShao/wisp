# 226 对抗验收 v2 —— 「always 答复整份覆盖 config.toml」缺陷（裁决表）

> 本件＝票 226 非实现者对抗验收腿 `226-v2` 的裁决表。
> 前一腿 `226-v1` 把 34 发攻击跑完、源码逐字节还原，但死在轮次上限、一个字都没写，本件接手把读数收成表。
> 主料＝`.scratch/wisp/probes/226/v1/`（26 顶层文件＋7 枚 `probe*/main.go`）。
> ⛔ 本腿不重跑 v1 跑过的东西；只自跑第 3 节点名的两发（① AC#3 无条件认领突变、② 删行一问，若 v1 读数不足）。
> 实现腿 `226-w1` 自述件＝`docs/evidence/s1/226-config-write-no-clobber-r1.md`，本件一个字都不预认。
> 判据 AC#1..AC#6 与禁区见工单
> `.scratch/wisp/issues/226-always-answer-rewrites-the-whole-config-snapshot-and-hides-hand-edits.md`。

---

## ① 锚点与四枚被验收 commit 复核

时刻 `2026-09-29 13:58 +08`（本腿起手现跑 `date`）。起手锚点 `git rev-parse --short HEAD`＝**`c039bc4e`**。

四枚被验收 commit 的 `git show --name-status` 现读复核（本腿亲跑，非抄实现腿自述）：

- **`78988901`（唯一产码枚）**：
  - `A internal/config/writeguard.go`（新文件，`mergeWrite` 主体，见下）
  - `A internal/config/writeguard_226_test.go`（本票新增用例）
  - `A cmd/wisp/always_write_no_clobber_226_test.go`（本票新增用例，走 CLI 接缝）
  - `M internal/config/allowdirs.go`（`writeAllowedDirs` 改走 `mergeWrite`）
  - `M internal/config/loader.go`
  - `M internal/config/permmode.go`（`SetPermissionMode` 同形那支一起处理，AC#4）
  - `A .scratch/wisp/probes/226/r1/*`（实现腿自己那 7 枚台件）＋`A docs/evidence/s1/226-config-write-no-clobber-r1.md`
- **`361b314a`**：`A .scratch/wisp/probes/226/r1/ac6-full-suite-start.txt`／`ac6-full-suite.txt` ＋ `M` 证据件。⛔ 代码一件未改。
- **`d90b9aec`**：仅 `M` 证据件 §10 回填。代码未动。
- **`c31239a9`**：仅 `M` 证据件 §3 行 7 ＋ §11 commit 台账。代码未动。

⇒ **产码只在 `78988901` 一枚**，后三枚都是文档/台件补全。本腿读源码（`internal/config/writeguard.go` 现读 337 行）确认 AC#3 的收窄就在该文件末段：`mergeWrite` case 3 里 `diverged := diffKeyPaths(m.cur, base)`，`if len(diverged) == 0 { … m.statOwnWrite(); return nil }`（writeguard.go:166–177），非空则走 `slog.Warn` 且不认领（181–186）。这正是 §3① 突变要改的那一格。

四枚 commit 的 `--name-status` 均只碰 `internal/config/**`、`cmd/wisp/**`、`.scratch/wisp/probes/226/r1/**`、`docs/evidence/s1/226-config-write-no-clobber-r1.md`，**未见越界**（无 `frontend/**`／`design/**`／`PLAN.md`／`thresholds.go`／golden／三枚冻结件）。

---

## ② v1 34 发逐枚对应表（发＝改了什么／原始读数／红绿第几例）

时刻 `2026-09-29 13:58 +08`。逐枚打开 `.scratch/wisp/probes/226/v1/` 对号，原始红／绿句抄自该目录文件本身（行号引用＝该 .txt 的行）。

> ⚠ 枚数口径：v1 目录顶层 26 枚 .txt，其中 `head-writeguard_226_test.go` 是被测测试文件的逐字副本、不是一发攻击（只作对照原文），故**实际发＝25 枚顶层读数＋7 枚 `probe*/main.go`**（后者是 `b-*`／`d-*` 那些小实验程序的真身）。票面记的"34 发"与本腿现量的"25＋7＝32"差 2 枚，见 §⑦。

| 文件 | 这一发射的是什么 | 原始读数（红／绿、具名例） |
|---|---|---|
| `g0-config-perm-v.txt` | 基线：`go test ./internal/config ./internal/perm -v`（未改码） | **全绿**：行298 `ok internal/config 0.777s`、行335 `ok internal/perm 0.231s`，PASS |
| `a1-config-mutated.txt` | 突变①：把 `allowdirs.go` 的写前重读／合并拿掉（`writeAllowedDirs` 退回交整份快照给 `SaveFile`） | **红 8 例**：AC#1 `TestAC1AllowedDirsWriteKeepsAHandEditedKey`（:74 got 56 want 60）、`TestAC2GuardedWriteReportsTheKeysItWroteAndTheOnesItKept`、`TestAC2CleanWriteReportsOneKeyAtInfoLevel`、`TestAC2UnreadableFileIsRefusedAndNothingIsWritten`、AC#3 `TestAC3AdoptionClaimsOnlyWhatThisWriteProduced`、AC#3 `TestAC3HandAddedEntryIsPreservedButNotAppliedToMemory`、AC#5 `TestAC5SetAllowedDirsPersistsARevocationWithoutAClobber`、`TestGuardedWriteDoesNotRewriteAFileThatAlreadySaysIt`。AC#4 不响（未动 permmode） |
| `a1-cmdwisp-mutated.txt` | 同一枚 a1 突变、经 CLI 接缝 `cmd/wisp` 触发 | **红 1 例**：`TestAC1AlwaysBranchDoesNotRevertAHandEditedKey`（`always_write_no_clobber_226_test.go:87` app.theme="dark" 被还原、want light） |
| `a1b-cmdwisp-control.txt` | a1 的正控：`cmd/wisp` **未改码**跑同名三例 | **全绿**：`TestAC1AlwaysBranchDoesNotRevertAHandEditedKey`／`TestAlwaysBranchStoresItsRuleOnlyAfterASecondL2Card`／`TestAlwaysBranchStoresNothingWhenTheSecondCardIsRefused` 皆 PASS，行38 `ok cmd/wisp 3.105s` |
| `a2-config-permmode-mutated.txt` | 突变②：把 `permmode.go`（`SetPermissionMode`）的合并拿掉（AC#4 那支同形） | **红 1 例**：`TestAC4PermissionModeWriteKeepsAHandEditedKey`（:278 got 56 want 60）。余绿 |
| `a3-config-unconditional-adopt.txt` | **突变③（本票主攻）**：把 `writeguard.go:168` "仅 `diverged` 为空才 `statOwnWrite()`" 改成无条件认领 | **红 2 例**：AC#3 `TestAC3AdoptionClaimsOnlyWhatThisWriteProduced`（:224 "the merged write adopted a stat for content this process never read; the hand edit is now invisible forever: \<nil\>"）、AC#3 `TestAC3HandAddedEntryIsPreservedButNotAppliedToMemory`（:259 "the hand-added [fs] entry was never judged; our write hid it"）。⇒ AC#3 那格**不是装饰**，本腿 §③ 复跑坐实 |
| `a4-control-as-shipped.txt` | 常驻正控未改态 | **绿**：`TestAC1ControlSnapshotWriteIsWhatRevertsTheHandEdit` PASS |
| `a4-control-inverted.txt` | 把该正控反接（令快照写不再还原）验它有没有牙 | **红**：同例 FAIL（:111 "INVERTED-PROBE: the control has lost its teeth"）⇒ 常驻正控**非恒真**，AC#1 的绿可归因给合并 |
| `b-probe-deleted-lines.txt` | 实验：手删一枚标量行（B1）／删整节（B2）后跑受守卫写 | B1/B2 皆 `RESULT ... size line re-added to the file : true`、`ball.size on disk now : 56`、`memory ball.size : 60`、comment 未存活、手排顺序未保留 ⇒ **删的行被补回、且补成 schema 默认值 56（非内存那 60）** |
| `b2-probe-deleted-shapes.txt` | 六种删除形状矩阵 P1–P6 | P1 非默认标量删→补回 56；P2 默认标量删→补回；P3 整节删→header+行都补回；**P4 provider 子表（map 项）删→header 未补回（`providers=0`）**、其 `base_url` 行补回 `''`；P5 `schema_version` 行删→补回 2；P6 `base_url` 删→补回。**标量／节＝补回；map 项＝不补回** |
| `b3-probe-adopt-and-migrate.txt` | 实验：删除后 `statOwnWrite` 认领与迁移的耦合 | X1 删默认值行→**认领**（reload `report=<nil>`，删除对 reload 路径不可见，但系默认==默认之无操作）；X3 删非默认 60→**不认领**、`report=&{Hot:[ball]...}`；X2 删 `schema_version`→认领并落 `.bak-1` |
| `b4-probe-migration-inside-read.txt` | 实验：`readConfigFile` 内部迁移是否越权改写 | Y1（已带该键＋schema 在场）→ `nothing written`、字节不变、mtime 不变、comment 存活；**Y2（手删 `schema_version`）→ 仍记 `nothing written` 但 `bytes changed:false→真变了`、mtime 变、size 104→92、comment 丢失**、多出 `config.toml.bak-1` ⇒ **迁移-on-读即便在"什么都不写"支也会重排文件** |
| `b5-probe-refused-write-migrates.txt` | 实验：不可解析／未知键时是否仍迁移改写 | Z1/Z2 皆 `AddAllowedDir returned: ... refused to write ... cannot be read back`、`bytes changed by the REFUSED write : false` ⇒ 读失败先于迁移，**fail-closed、不动文件** |
| `b6-probe-nothing-written-branch.txt` | 实验："什么都不写"承诺在迁移情形下是否兑现 | W1（当前 schema、已带键）→ 字节同、mtime 同、comment 保留（真·什么都不写）；W2（删 `schema_version`）／W3（塞入 `schema_version = 1` 旧文件）→ **迁移把文件重排成 canonical v2、comment 丢失**、落 `.bak-1` ⇒ "nothing written" 支只对 schema 已对的键成立 |
| `c1-relative-hop-check-removed.txt` | 突变：拿掉含 `..` 相对跳拒收 | **红 1 例**：`TestAddAllowedDirStillRefusesARelativeHop`（:338 accepted `D:\data\..\escape`）。`TestAC2Unreadable...`／`TestGuardedWriteFailureRollsBackMemory` 仍绿 |
| `c2-rollback-and-sethop-removed.txt` | 突变：拿掉回滚内存＋`SetAllowedDirs` 的 hop 拒收 | **红 2 例**：`TestAC2UnreadableFileIsRefusedAndNothingIsWritten`（:182 拒写未回滚内存，got `[D:\data]`）、`TestAddAllowedDirStillRefusesARelativeHop`（:351/:354 `SetAllowedDirs` 收下 hop 且先改了内存） |
| `c3-permmode-rollback-removed.txt` | 突变：拿掉 `SetPermissionMode` 写失败回滚 | **红 1 例**：`TestGuardedWriteFailureRollsBackMemory`（:387 a failed mode write left the new mode in memory）；`internal/perm` `TestTicket90PersistFailureKeepsMemory` 仍 PASS（行13） |
| `c-seed-relative-hop-removed.txt` | c1 的正控：种"收 hop"必响 | **红 1 例**：`TestAddAllowedDirStillRefusesARelativeHop` FAIL（:338）⇒ hop 钉有牙 |
| `c-seed-rollback-removed.txt` | c2 的正控：种"不回滚"必响 | **红 2 例**：`TestAC2UnreadableFileIsRefusedAndNothingIsWritten`＋`TestGuardedWriteFailureRollsBackMemory` |
| `c-seed-permmode-rollback-removed.txt` | c3 的正控：种"mode 不回滚"必响 | **红 1 例**：`TestGuardedWriteFailureRollsBackMemory`（:387）；`internal/perm` `TestTicket90...` 两例仍 PASS |
| `d-probe-d36-adopt-bypass.txt` | 问：不认领收窄后会不会绕过 D36 判级 | S1（手松 `risk.permission_mode`＝本写自己的键）→认领、无残留外键；S2（`net.block_private_ranges` 未知键）→**拒写**；S3（手松 `plugins.tier2_enabled`）／S4（手加 `fs.allowed_dirs` 一项）→ **NOT claimed、reload 判 `locked loosening rejected`（fail-closed D36 rule 1）** ⇒ **不绕过 D36**：被保留的分叉外键不认领、下一轮读到并按 D36 拒 |
| `f-full-suite.txt` | AC#6 整包第一发（as shipped，安静期） | **24 包 ok／FAIL 只有 5 例**＝`internal/ball` `TestC21TableColourRowsMatchTokensCSS`＋`internal/panel` 四例（`TestApprovalCardViewJSONKeysMatchFrontendTypes`／`TestComposerContractTypesMatchFrontend`／`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`／`TestC21DesignTokensFourWayAgree`）；`internal/config ok 1.454s`、`internal/perm ok 0.332s`、`cmd/wisp ok 86.122s`、`internal/risk ok 5.754s` |
| `f2-full-suite-terminal.txt` | AC#6 第二发 | **同集合、同 5 枚红名**，逐包名册与 f 一致；`internal/risk ok 5.624s` |
| `f3-full-suite-final.txt` | AC#6 第三发 | **同集合、同 5 枚红名**；`internal/config ok 1.106s`、`internal/risk ok 6.023s` |
| `f-d22scan.txt` | 仪器：d22scan 自检＋对当前树跑 | 全 PASS（含 `TestScanCleanRepoIsGreen`、种子违规端到端 `exits_1` 正控、`TestBuiltBinaryGoesRedEndToEnd` 各族），末行 `d22scan: clean - no D22 ban violations`；`runtests.sh: OK ... PASS=34 FAIL=0` |
| `head-writeguard_226_test.go` | 非攻击：被测测试文件逐字副本 | 不作一发计数；供 §②／§④ 引例名与断言行号 |

---

## ③ 本腿自跑的两发原始输出

时刻 `2026-09-29 14:06 +08`（本腿起手锚点 `c039bc4e` 未动，HEAD 现为我这几枚文档 commit）。

### ③① AC#3 无条件认领突变（装饰判据检测）

起跑前先 `export PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"`。

**突变**（`internal/config/writeguard.go`，在 `diverged := diffKeyPaths(m.cur, base)` 之后、`if len(diverged) == 0` 之前插一行无条件 `m.statOwnWrite()`；未删原 `if` 内那一句，两次认领幂等）。跑 `go test ./internal/config -count=1`，全量落盘到 `.scratch/wisp/probes/226/v2/rerun1-unconditional-adopt.txt`（**未接 `| head`／`| tail`**）：

原始红句（逐字，抄自该件）：
```
--- FAIL: TestAC3AdoptionClaimsOnlyWhatThisWriteProduced (0.03s)
    writeguard_226_test.go:224: the merged write adopted a stat for content this process never read; the hand edit is now invisible forever: <nil>
--- FAIL: TestAC3HandAddedEntryIsPreservedButNotAppliedToMemory (0.02s)
    writeguard_226_test.go:259: the hand-added [fs] entry was never judged; our write hid it (AC#3)
FAIL
FAIL	github.com/CarlosShao/wisp/internal/config	0.738s
```

⇒ **结论：AC#3 那格不是装饰。把"仅有分叉就不认领"改成无条件认领，稳定红 2 枚**（且只红这 2 枚 AC#3 用例，AC#1/2/4/5 全绿）。这与 v1 的 `a3-config-unconditional-adopt.txt` 逐字一致（同一 :224／:259 两句），两腿独立复现 ⇒ 该判据真在测"写完之后再手改一枚键、下一次轮询要能看见"。

**逐字换回 + certutil 比对**（本腿要求抄进表）：
- 起手基线 `certutil -hashfile internal/config/writeguard.go SHA256`＝`bf82226988b83fdce151478bd51a550beec8d66caef39e5cc4096f6f55150f16`
- 还原后同一命令＝`bf82226988b83fdce151478bd51a550beec8d66caef39e5cc4096f6f55150f16`
- ⇒ **逐字节相同**；`git diff --stat -- internal/config/writeguard.go` 输出为空（未残留任何改列）。

### ③② "删行会不会被补回"那一发

**未自跑**——本腿按派单纪律"v1 已跑过则只读其读数"，`b-*`／`probe*/main.go` 七枚留下的输出已足以定论（见 §② 六行 b 族读数 + §⑤ 结论），不消耗第二次编译额度。§⑤ 即据此作答。

---

## ④ 逐格裁决 AC#1..AC#6

〔待填〕

---

## ⑤ 「删行会不会被补回」结论

〔待填〕

---

## ⑥ 推翻清单

〔待填〕

---

## ⑦ 本腿没做完／读不动的

〔待填〕

---

## 入库清单

〔待填：wc -c〕
