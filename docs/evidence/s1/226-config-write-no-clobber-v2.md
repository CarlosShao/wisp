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

时刻 `2026-09-29 14:07 +08`。每格：〔成立／不成立／成立但带注〕＋作依据的原文行＋"把实现拿掉会不会红"。

- **AC#1 手改不许被吞 —— 成立，但带硬注。**
  依据：`g0` 基线整包 `ok`（`writeguard_226_test.go` 全绿）；`a1-config-mutated.txt:10` 把 `allowdirs.go` 合并拿掉后 `--- FAIL: TestAC1AllowedDirsWriteKeepsAHandEditedKey`（:74 `got 56, want 60`）、`a1-cmdwisp-mutated.txt:12` CLI 接缝同突变 `--- FAIL: TestAC1AlwaysBranchDoesNotRevertAHandEditedKey`。`a4-control-inverted.txt` 证明常驻正控 `TestAC1ControlSnapshotWriteIsWhatRevertsTheHandEdit` 非恒真（反接即 :111 "the control has lost its teeth"）。**拿掉合并＝红**，答案：会。
  ⚠ 硬注：AC#1 的字面判据只钉住"**改一枚键的值**"这一支；用户"**删一行**"这一支没有任何用例覆盖，且实测删的行会被补回成 schema 默认值（详见 §⑤）。⇒ 票面"运行期间你手改的任何其它键会被悄悄还原"这句，对"删"这一形今天仍未被 AC#1 的判据接住。

- **AC#2 冲突必须响亮 —— 成立，但带注。**
  依据：`a1-config-mutated.txt` 拿掉合并后 AC#2 三枚用例同时红（:127 未报实际写的键、:154 干净写未在 INFO 报键、:170 竟接受了不可读文件并覆盖）；`writeguard.go:168-185` 分叉走 `slog.Warn` 且 `kept_in_file_not_in_memory` 点名保留的手改、`written` 点名实改的键。**拿掉合并＝红**，答案：会。
  ⚠ 注：合并报告里"什么都不写就不动文件"的承诺，只在 **schema 已是当前版** 时成立——`b4-probe-migration-inside-read.txt` Y2、`b6-probe-nothing-written-branch.txt` W2/W3 显示：手删 `schema_version` 行或塞入旧版文件时，`readConfigFile` 的迁移-on-读会把文件重排成 canonical（bytes 变、mtime 变、操作者注释丢、落 `config.toml.bak-1`），**即便 mergeWrite 走的是 `nothing written` 早返回支也一样**。这是继承自 loader 的行为、非本票新造，但 AC#2 文案"不白丢注释"要据此收窄。

- **AC#3 `statOwnWrite` 认领范围收窄 —— 成立。**
  依据：本腿 §③① 亲跑的"无条件认领"突变，稳定红 **2 枚**（`writeguard_226_test.go:224` "adopted a stat for content this process never read; the hand edit is now invisible forever"、:259 "the hand-added [fs] entry was never judged; our write hid it (AC#3)"），与 v1 `a3-config-unconditional-adopt.txt` 逐字一致。**把"有分叉就不认领"改成无条件认领＝红**，答案：会 ⇒ **AC#3 那格不是装饰**。`d-probe-d36-adopt-bypass.txt` S3/S4 另证：分叉外键保留在文件里时不认领、下一轮读到并按 D36 rule 1 拒放宽＝**没有把甲变成绕过 D36 的后门**。

- **AC#4 同形的另一处一起处理 —— 成立。**
  依据：`a2-config-permmode-mutated.txt:14` 拿掉 `permmode.go` 合并后 `--- FAIL: TestAC4PermissionModeWriteKeepsAHandEditedKey`（:278 `got 56, want 60`）；`c3-permmode-rollback-removed.txt:7` 拿掉 mode 写失败回滚后 `--- FAIL: TestGuardedWriteFailureRollsBackMemory`（:387）。`g0` 显示 `internal/perm ok`、冻结件 `TestTicket90PersistFailureKeepsMemory` 仍 PASS。**拿掉＝红**，答案：会。`SetPermissionMode` 与 `AddAllowedDir` 同批改，未只修新入口。

- **AC#5 "可撤销"那一半要有真路径 —— 半格成立（落库那半），另半未闭。**
  依据：`writeguard_226_test.go:308` `TestAC5SetAllowedDirsPersistsARevocationWithoutAClobber`（`g0` 绿）把"零调用者零用例"的 `SetAllowedDirs` 首次钉成受守卫写（撤销落库、且不覆盖别的手改）。`a1-config-mutated.txt` 拿掉合并后此例红（:324 revocation reverted the hand edit）。**落库＝会红**。但票面 AC#5 要求"**既能落库、又能被答复语法触发**"，答复语法那半不在本票射程（属票 219/224）——实现腿 §8 自己也记为未闭。**这一格只有半边为真**，另半边推迟事项须按 `SPEC-12 §5` 五字段登记（见 §⑥）。

- **AC#6 整包终态 —— 成立（零新增红），口径见下。**
  依据：v1 三发 `f-full-suite`／`f2-full-suite-terminal`／`f3-full-suite-final` **均为 24 包 ok、FAIL 名册恰为 5 例**＝`internal/ball` `TestC21TableColourRowsMatchTokensCSS`＋`internal/panel` 四例（`TestApprovalCardViewJSONKeysMatchFrontendTypes`／`TestComposerContractTypesMatchFrontend`／`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`／`TestC21DesignTokensFourWayAgree`），且三发 `internal/risk` 皆 `ok`（5.754／5.624／6.023s）、`internal/config ok`、`internal/perm ok`、`cmd/wisp ok`（81.8—86.1s）、`tools ok`。这 5 例＝别人地界的历史红（另一队删 `design/assets/tokens.css`＋面板契约字段差，`A427`/`A428`），**非本票新增**。
  ⚠ 读数矛盾已断：编排者 13:01:37 那一发是 **23 包 ok／3 包 FAIL 共 6 例**，多出 `internal/risk/TestResolvePerCallBudget`＝**并发争用假红**（撞在实现腿收尾的门上），编排者 13:02:57 空机器 `-count=3` 复量 0.281／0.272／0.290 ms/op（预算 1.000）三发全 PASS（台账 `A432`）。⇒ **"零新增红"的比对基线以 v1 这三发安静期 24-ok 读数为准**（全量落盘、未接 head/tail、逐名比红名册）；编排者 13:01:37 那发属争用期、其 risk 例不计入本票名册。本腿未新跑整包（非第 3 节点名两发之一），不另立第四发读数。`f-d22scan.txt` 另证仪器侧：末行 `d22scan: clean - no D22 ban violations`、`runtests.sh OK PASS=34 FAIL=0`（此处 34＝d22scan 自检包计数，与本票 34 发读数无关）。

---

## ⑤ 「删行会不会被补回」结论

时刻 `2026-09-29 14:07 +08`。⛔ 未自跑（v1 `b-*`＋`probe*/main.go` 读数已足），据以下原始输出定论：

**直接答案：会。** 用户从 `config.toml` 里**删掉一行标量键或整节**，一次受守卫的写会**把那行补回文件**——补成 **schema 默认值**（不是本进程内存里那枚值）。

- `b-probe-deleted-lines.txt` B1/B2：手删 `size` 行后受守卫写 → `RESULT ... size line re-added to the file : true`、`ball.size on disk now : 56`、`memory ball.size : 60`（补回的是默认 56，非内存的 60）。
- `b2-probe-deleted-shapes.txt`：P1 非默认标量删→补回 56；P2 默认标量删→补回；P3 整节删→header+行都补回；P5 `schema_version` 行删→补回 2；P6 `base_url` 删→补回 `''`。

**是哪段代码"允许"它补回的**：`readConfigFile(m.path)`（`writeguard.go:118`）重读磁盘时，缺失的标量键会被填成结构体默认；`mergeWrite` 以 `base = deepCopyConfig(disk)`（`:141`）为底、`SaveFile` 序列化整个结构体 ⇒ 文件里那行"缺席"与"等于默认"在结构体模型里**不可区分**，于是缺的标量行被重新落盘成默认值。

**边界（不会被补回的那一支）**：`b2` P4 删一个 **map 条目**（provider 子表）→ `[llm.providers.deepseek]` header `write put it back: false`、`loaded-after: providers=0` ⇒ map 项的"缺席"被保留、**不补回**。

**⇒ 对票面的后果**：AC#1 只钉住了"**改值**"这一支（`writeguard_226_test.go` 全部用"把 56 改成 60"作手改），**"删行"这一支既没有用例、且实测会被还原成默认**。这正是工单第 4 行编排者自己起的那枚疑问的结论——**成立**：票 226 的 AC#1 覆盖面不完整。

**建议（本腿不拍、留给编排者）**：
1. 要么新增判据（如 **AC#1-b**：手删一枚标量行 ⇒ 受守卫写不得把它按默认值补回 / 或明确"删行＝回落默认"是可接受语义并写进用户可见文案），要么另立一票；
2. 且**建议把这行为钉成常驻用例**（现有 `b-*` 只是台件、非 `*_test.go`）：至少断言"删一枚非默认标量行后跑受守卫写，文件里该行的值等于默认、且 `kept_in_file_not_in_memory` 有报（与内存分叉时不认领）"——否则该语义无回归防护；
3. 迁移-on-读在 `nothing written` 支仍重排文件（Y2/W2/W3）这一条**也建议立一枚常驻用例**钉住"schema 落后时受守卫写会重排并丢注释"，把它从"隐性副作用"变成"被断言的已知代价"。

---

## ⑥ 推翻清单

〔待填〕

---

## ⑦ 本腿没做完／读不动的

〔待填〕

---

## 入库清单

〔待填：wc -c〕
