# 票 226 — 交付证据件 r1（写腿 `226-w1`）：`always` 不再把整份内存快照写回 `config.toml`

- 工单：`.scratch/wisp/issues/226-always-answer-rewrites-the-whole-config-snapshot-and-hides-hand-edits.md`
- 台账出处：`docs/reports/pending-and-issues.md` 的 `A428` 第 ② 条（缺陷怎么被发现）＋ `A424` 第 1 条边界（"持久、可审计、可撤销"）
- 治理：根目录 `AGENTS.md` §1.1（`D/C/R` 契约与 SLO/golden/`thresholds.go` 一字不动）、§1.2（`tools/d22scan` 禁形）、§1.4（只 commit 不 push、commit 带显式 pathspec）
- 本件作者＝**实现者本人**，不是裁决者；按 `AGENTS.md` §0 第 3 句与 `SPEC-12 §4.3` #1/#3，**缺口审计与对抗验收必须由另一个 agent 做**，本件不构成自证。
- 定锚（起手 HEAD）：`4540d5b11abce0522f811068b10c902647c06846`（committer 时间 `2026-09-29T12:25:21+08:00`，现跑 `git log -1`）
- 起手时刻：`2026-09-29 12:26 +0800`（现跑 `date`）

## 0. 一句话结论

写路径改成**写前重读磁盘＋只替换这一节的那一枚键**（票面候选**甲**），并把 `statOwnWrite()` 的认领范围收窄到"这次写产生的内容与内存一致"那一支；`SetPermissionMode` 同形那半（AC#4）**同批改掉**，没有留残缺。

## 1. 选了哪一支，为什么（票面把这一支交给自己定并写明理由）

**选甲（写前重读＋只替换那一枚键），不选乙、不选丙。** 三条现读依据，不是口味：

1. **乙省不了那枚差异器。** 票面 AC#2 给乙的那半句是"拒绝并**说清哪几枚键冲突**"，而本包今天没有任何能把两枚 `*Config` 比出键路径的东西（现跑 `grep -rn "diffKeyPaths|diffConfig|func diff" internal/ cmd/ tools/` 只命中本件 `internal/config/schema.go:570` 那句"the diff walk treats it as a transparent map"的过期注释、零实现）。乙为要名冲突必须写这枚 `diffKeyPaths`，为不改文件又什么都不用写 ⇒ **同样的钱、更少的保护**。
2. **乙在生产里不可恢复**：轮询没接上（票 223 已证、且它在派单里排在本票之后）。手改一次 ⇒ 磁盘与内存从此分叉 ⇒ 之后每一发「一直」都被拒 ⇒ 直到重启。甲把这枚"长期允许"保住，同时把手改也保住。
3. **丙没做**：逐键决定"内存赢还是磁盘赢"每一个都是一枚新的静默行为。甲已经把丙的**前置条件**（AC#2"报告实际写了哪几枚键"）交付了，所以按票面那句"没有 AC#2 就别做丙"，丙仍然不是今天的选项。

**甲的代价，照票面要求写成明文案**（不静默）：`MarshalCanonical` 是结构体序列化，注释与手排键序本来就不保留 ⇒ 一次受守卫的写会把文件重排成 canonical 形，**值不丢、格式丢**。这一句同时写在代码注释（`internal/config/writeguard.go:45-56`）与那两行日志文本里（`internal/config/writeguard.go:169-189`，实测输出见 §5），所以它在用户的日志里也可读，不只在源码里。

**AC#3 的认领规则**（这一条不是照抄票面候选，是甲的必要配套）：`statOwnWrite()` 只在"**这次写产生的内容 == 内存持有的内容**"时才调用；一旦合并把 foreign 值带进文件，就**不认领**，让下一次轮询把文件读回去、由 `applyLocked` 对它作 D36 裁决（放宽仍是拒绝＝fail-closed）。这样既没有"我们写过一次之后什么都算我自写的"，也没有让甲变成一枚绕过 D36 的后门。

**⛔ 顺带没有做的**：没有给 `CheckAndReload` 接轮询、没有给 `ConfirmLocked` 找赋值点（那是票 223 的判据表），合并读盘**只用于写、不用于把值装进内存**——`TestAC1AllowedDirsWriteKeepsAHandEditedKey` 与 `cmd/wisp/always_write_no_clobber_226_test.go` 各有一枚断言钉住这一点（内存仍是开机那枚 `56` / `dark`）。

## 2. 改了哪几枚文件（包级路径＋现跑 `grep -n` 行号，锚点＝提交 `78988901`）

| 文件（包级） | 行 | 是什么 |
|---|---|---|
| `internal/config/writeguard.go`（**新**，336 行） | `:111` | `mergeWrite`：重读→只在文件的实际内容上替换自己那一枚键→写→报告 `wrote` 与 `kept_in_file_not_in_memory`→按 AC#3 决定是否认领 stat |
| 同上 | `:70-71` | 两个包内键路径常量 `keyFSAllowedDirs` / `keyRiskPermissionMode`（**未导出**，不是 C17 面） |
| 同上 | `:201` `:216` `:241` `:287` `:330` | `diffKeyPaths` / `diffStruct` / `configFieldName` / `diffMap` / `truncateKeys`：按 `toml` tag 走，报的点是文件里真存在的键路径；`schema_version` 与 `toml:"-"` 跳过；`Plugins` 走 `plugins.<id>`（`marshalPlugins` 的形） |
| `internal/config/loader.go` | `:39` `:62` | `LoadFile` 拆成 `readConfigFile`（sec 4.1 管线，止于 `validate`、不做 `api_key_ref` 解析）＋ `resolveRefs`；**管线顺序与错误文案一字未动**，`LoadFile` 对外行为不变。拆的理由写在 `:51-61`：重读不该被 DPAPI 失败拖成"配置写不了" |
| `internal/config/allowdirs.go` | `:62` `:109` `:139` `:142` `:157` | `writeAllowedDirs` 交 `mergeWrite`；`AddAllowedDir` 的追加改成"接在**文件**现有名单后面"（手加的那行既不重复也不丢）；`SetAllowedDirs` 明确"只对 `fs.allowed_dirs` 这一枚键权威"；`statOwnWrite` 的注释改成"只在一处被调用、且那一处已证内容与内存一致" |
| `internal/config/permmode.go` | `:70` `:79` | AC#4 同形那半**同批**改：`SetPermissionMode` 走 `mergeWrite`，删掉它自己那四行认领 |
| `internal/config/writeguard_226_test.go`（**新**，456 行） | `:63` 起 13 枚用例 | 这枚写路径的**第一份直接单测**（此前 `grep "^func Test" internal/config/*_test.go \| grep -i allow` 为空，只被 `cmd/wisp/approval_always_201_test.go` 间接覆盖） |
| `cmd/wisp/always_write_no_clobber_226_test.go`（**新**，106 行） | `:39` | AC#1 的"起进程"形：真实 `always` 分支＋真实第二张 L2 卡＋真实 `Manager.AddAllowedDir` |
| `.scratch/wisp/probes/226/r1/`（**新**，7 枚读数件） | — | 门、逐名绿册、两发现修前红突变、d22scan |

**没动的**：`internal/perm/ticket90_persist_test.go`、`internal/panel/tokens_fourway_test.go`、`internal/panel/l2_grant_boundary_test.go`、`PLAN.md`、`docs/specs/**`、`thresholds.go`、golden、`allowlist.txt`、`scripts/d22scan.sh`、`internal/tools/**`、`frontend/**`、`design/**`（由 `git show --stat 78988901` 的 14 枚清单可核）。

## 3. 门与读数（全部本人现跑，原始输出在 `.scratch/wisp/probes/226/r1/`）

每发之前都跑了 `export PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"`（少这步 `cmd/wisp` 会以 `0xc0000135` 假失败且没有任何 `--- FAIL` 行）。

| # | 命令 | 时刻（现跑 `date`） | 原始读数 |
|---|---|---|---|
| 0 | `go build ./...` **动手前**（基线） | 12:27 | 空输出、`BUILD_EXIT=0` |
| 0 | `go test ./internal/config ./internal/perm ./cmd/wisp -count=1` **动手前**（基线） | 12:27→12:28 | `ok internal/config 0.711s`／`ok internal/perm 0.296s`／`ok cmd/wisp 79.856s` ⇒ **包级预检今天全绿** |
| 1 | `go build ./...` 交件前 | 12:50 | `BUILD_OK exit=0` |
| 2 | `go test ./internal/config ./internal/perm ./cmd/wisp -count=1` | 12:50→12:51 | `gate-config-perm-cmdwisp.txt`：`ok internal/config 0.990s`／`ok internal/perm 0.434s`／`ok cmd/wisp 84.582s`（`EXIT=0`） |
| 3 | `go test ./internal/config ./internal/perm -count=1 -v` 逐名册 | 12:51 | `config-perm-named-green.txt`：`--- PASS` **83 枚**、`--- FAIL`/`--- SKIP` **0 枚** |
| 4 | `gofumpt.exe -l`（`"$GOPATH/bin"`, v0.12.0）我改的五枚文件 | 12:51 | 空 |
| 5 | `sh scripts/d22scan.sh`（不改这枚脚本） | 12:52 | `d22scan.txt`：步 1 正向对照 `runtests.sh: OK - packages=[./...] top-level: PASS=34 FAIL=0 SKIP=0`；步 2 `d22scan: clean - no D22 ban violations`，实扫 `internal/` 219＋`cmd/` 28 枚生产 Go 文件（`--self-test` 那枚 seed 对照在场，所以"clean"不是"仪器瞎了"） |
| 6 | AC#6 整包终态 `go test ./cmd/wisp ./internal/... -count=1` | 见 §7 | 逐名红册见 §7 |

### 3.1 派单里预检的那四枚用例名（今天真绿，逐名从 `config-perm-named-green.txt` 抄）

- `--- PASS: TestSaveFileRoundTripsThroughLoad (0.00s)`（`internal/config/loader_test.go:278`）
- `--- PASS: TestManagerReloadTierEmitsEvent (0.02s)`（`internal/config/manager_test.go:113`）
- `--- PASS: TestManagerVoiceHotKeysDoNotEmitReload (0.01s)`（`internal/config/manager_test.go:136`）
- `--- PASS: TestAC2SaveFileLandsPrivate (0.08s)`（`internal/config/private_acl_windows_test.go:107`）
- 设了 `ConfirmLocked` 的五枚（`manager_test.go:158/:188/:206/:241` 所属）：`TestManagerLockedLooseningRejectedKeepsOld`／`TestManagerLockedLooseningApprovedApplies`／`TestManagerLockedTighteningHotAppliesWithoutHook`／`TestManagerLockedNilConfirmHookDenies`／`TestManagerLockedBothDirectionsLogged` = **全 PASS**
- `TestAC2BaselineModeIsDecorative`／`TestAC2MigrationBackupLandsPrivate` 也 PASS（同一枚文件，顺手记）

### 3.2 ⛔ 冻结件没有被打红（这是停手格，所以单列）

`internal/perm/ticket90_persist_test.go` 我一字未动，`-v` 册里它的**十二枚**用例在改后仍全 PASS，其中包含派单点名的那三处 `ConfirmLocked`（`:342`／`:386` 在 `TestTicket90HandEditLooseningGoesThroughD36` 内，`:329` 是它的文档注释行）：

```
--- PASS: TestTicket90HandEditLooseningGoesThroughD36 (0.08s)
--- PASS: TestTicket90PersistFailureKeepsMemory (0.00s)
--- PASS: TestTicket90ManualSwitchSurvivesRestart (0.02s)
--- PASS: TestTicket90ConfigKeyChangesWhatTheChainAsks (0.02s)
--- PASS: TestTicket90SessionGrantDoesNotSurviveRestart (0.08s)
（其余七枚见 config-perm-named-green.txt）
```

`TestTicket90PersistFailureKeepsMemory` 与 `TestTicket90HandEditLooseningGoesThroughD36` 分别正是"写失败回滚内存"与"D36 那圈不能被写路径绕过"的两枚钉——它们绿＝这两条既有行为没被本票削弱（不是我的自述，是可复跑的读数）。

## 4. 修复前红／修复后绿（票面 AC#1 要求的正控方向）

我没有"信"任何历史读数：**把合并那一步拿掉**是用真突变跑的，两枚突变各自只改一处，跑完立刻还原并比对哈希。

| 发 | 突变（只此一处） | 结果 | 件 |
|---|---|---|---|
| 修复前红 ① | `internal/config/allowdirs.go:142` 的 `m.mergeWrite(...)` 换回旧形 `SaveFile(m.path, m.cur)`＋`m.statOwnWrite()` | **5 枚 FAIL**：`TestAC1AllowedDirsWriteKeepsAHandEditedKey`／`TestAC2GuardedWriteReportsTheKeysItWroteAndTheOnesItKept`／`TestAC3AdoptionClaimsOnlyWhatThisWriteProduced`／`TestAC3HandAddedEntryIsPreservedButNotAppliedToMemory`／`TestAC5SetAllowedDirsPersistsARevocationWithoutAClobber`（AC#4 那枚这一发仍绿，因为只动了 allowed_dirs 那一支）；失败原文：`the operator's [ball] size edit was reverted by the allowed_dirs write: got 56, want 60` | `mut-allowdirs-no-merge.txt` |
| 修复前红 ② | `internal/config/permmode.go:79` 同样换回旧形 | `TestAC4PermissionModeWriteKeepsAHandEditedKey` = **FAIL**，原文：`the mode switch reverted the operator's [ball] size edit: got 56, want 60` | `mut-permmode-no-merge.txt` |
| 还原 | `python` 逐串换回，`certutil -hashfile … SHA256` 比对 | 两枚文件与动手前**逐字节相同**（`ALLOWDIRS_RESTORED_IDENTICAL`／`PERMMODE_RESTORED_IDENTICAL`） | 本件 §4 |
| 修复后绿 | 无突变 | `internal/config` 整包 `ok 0.990s`（新 13 枚用例全 PASS）＋ `cmd/wisp` 的 `TestAC1AlwaysBranchDoesNotRevertAHandEditedKey` `--- PASS (1.24s)` | `config-perm-named-green.txt`／`cmd-wisp-ac1-green.txt` |

**一枚常驻正控**（不依赖谁记得跑突变）：`TestAC1ControlSnapshotWriteIsWhatRevertsTheHandEdit`（`internal/config/writeguard_226_test.go:102`）直接把内存快照交给 `SaveFile`，断言**这一形确实会把 `60` 改回 `56`**；若哪天它不再还原，用例自己 `t.Fatal("the control has lost its teeth…")`。这正面回答票面那句"如果它现在就是绿的，先证明它绿是因为合并真做了，而不是因为测试根本没改到第二个键"。

## 4. 逐格判据

〔待填〕

## 5. 逐格判据（票面 AC#1..AC#6；每格给"用例名＋它到底断言了什么"，不是"字段在场"）

| 格 | 状态 | 用例（`internal/config/writeguard_226_test.go` 除注明外） | 断言内容 |
|---|---|---|---|
| **AC#1 手改不许被吞** | **已闭** | `TestAC1AllowedDirsWriteKeepsAHandEditedKey:63`；`TestAC1AlwaysBranchDoesNotRevertAHandEditedKey`（`cmd/wisp/always_write_no_clobber_226_test.go:39`）；正控 `TestAC1ControlSnapshotWriteIsWhatRevertsTheHandEdit:102`；突变两发见 §4 | 包级：手改 `[ball] size = 60` → `AddAllowedDir` → 读回文件 `size` 仍是 `60`、`fs.allowed_dirs` 真含那枚目录、内存仍 `56`。进程级（票面"起进程"形）：真实 `always` 分支＋第二张 L2 卡落盘后，`[app] theme = 'light'` 仍在文件里、长期规则也在、内存仍是 `dark`。正控断"快照直写确实会把它改回 56"，所以绿只能来自合并 |
| **AC#2 冲突必须响亮** | **已闭（选"合并并报告实际写了哪几枚键"那一支）** | `TestAC2GuardedWriteReportsTheKeysItWroteAndTheOnesItKept:117`；`TestAC2CleanWriteReportsOneKeyAtInfoLevel:146`；`TestAC2UnreadableFileIsRefusedAndNothingIsWritten:164` | 分叉时日志必含 `wrote=[fs.allowed_dirs]` **且** `kept_in_file_not_in_memory=[ball.size]` **且** `level=WARN`，并反向断言**不许**出现 `wrote=[ball.size…`（不许把没改的键报成自己写的）。干净时只报一枚键、`level=INFO`、不许出现 `kept_…`。**读不回的文件＝拒绝写**：错误须含 `cannot be read back`＋`keeping the previous allowlist in memory`，文件字节不变、内存不变（这一支就是票面乙那半句"拒绝并说清"，只是它只在"无法计算合并"时出现） |
| **AC#3 认领范围收窄** | **已闭** | `TestAC3AdoptionClaimsOnlyWhatThisWriteProduced:187`；`TestAC3HandAddedEntryIsPreservedButNotAppliedToMemory:234` | (a) 无分叉：写完 `CheckAndReload()` **返回 nil**（认领只覆盖这次写的内容），**之后再手改一枚键 ⇒ 下一次轮询必须看见**（`rep.Hot` 含 `ball`）。(b) 有分叉：写完**不认领** ⇒ 紧接着的 `CheckAndReload()` 必须看见那枚手改（看不见就是票面后果 3 没修掉）。(c) 手加的那枚 `[fs]` 条目：文件里两枚都在，内存里只有本次确认过的那枚，且**不认领**使下一次轮询真去作 D36 裁决（`ConfirmLocked=false` ⇒ 判拒绝、内存不变严不放松）——这条是防"甲把合并变成绕过 D36 的后门" |
| **AC#4 同形的另一处一起处理** | **已闭（同批改，未登记残缺）** | `TestAC4PermissionModeWriteKeepsAHandEditedKey:269`；`TestGuardedWriteFailureRollsBackMemory:360`；突变红发 §4 行② | `SetPermissionMode` 走同一枚 `mergeWrite`：手改的 `ball.size` 在模式切换后仍在文件里、`risk.permission_mode` 真落盘、内存不被偷偷升级；无分叉时切换仍被认领（`CheckAndReload()` 返回 nil，保住 R20/M3 的持久语义）。`internal/config/permmode.go:70/:79` |
| **AC#5 "可撤销"要有真路径** | **半格：落库那半已闭，答复语法那半未闭（见 §8）** | `TestAC5SetAllowedDirsPersistsARevocationWithoutAClobber:308`；`TestGuardedWriteDoesNotRewriteAFileThatAlreadySaysIt:393`；`TestAddAllowedDirStillRefusesARelativeHop:330` | `SetAllowedDirs` 从此**有用例**（此前全仓零调用者零用例）：撤销一行确实写回文件、内存与文件不分叉、别的手改不被顺带覆盖。附带一枚：文件里已经有这枚键要的值时**一个字节都不写**（不白重排、不白丢注释、不动 mtime）。⚠ **没有**新增 C17 方法名，也没有新增答复语法（那是 219/224 的面） |
| **AC#6 整包终态** | 已跑，见 §7 | — | 逐名红册＋分口径归因 |

**两条既有安全行为（派单"一条都不许松"）单独钉住**：
- 含 `..` 拒收：`internal/config/allowdirs.go:67-73`（`AddAllowedDir`）与 `:116-119`（`SetAllowedDirs`）的 `fail-closed` 一字未动，`TestAddAllowedDirStillRefusesARelativeHop:330` 现在是它的**第一枚直接单测**（断两枚入口、断文件字节不变、断内存不变、断 Set 那一支不许先改内存再拒）。
- 写失败回滚内存：`internal/config/allowdirs.go:139-147` 与 `internal/config/permmode.go:77-85` 的 `old` 还原与错误文案原样保留；`TestGuardedWriteFailureRollsBackMemory:360` 用真写入失败（父目录不存在 ⇒ `atomicWrite` 建不了临时件）断两枚入口都回滚；冻结件 `TestTicket90PersistFailureKeepsMemory` 仍绿。

实测的 AC#2 那行日志（从 `cmd-wisp-ac1-first-run.txt` 里进程自己的 stderr 抄，逐字）：

```
level=WARN msg="config: wrote merged change into config.toml and kept hand edits this process does not hold (ticket 226); the file was NOT claimed as our own write, so the next reload reads it back and a loosening there is denied until it is confirmed (D36 rule 1). Only the keys listed under wrote changed value; the file is re-laid in canonical form, so comments and hand-chosen key order are not preserved" key=fs.allowed_dirs wrote=[fs.allowed_dirs] kept_in_file_not_in_memory=[app.theme]
```

## 6. Git 纪律（照字面做的，附可核证据）

- 起手与每次 commit 之前都跑了 `git diff --cached --name-only`：**起手为空、commit 前索引里只有我自己的 14 枚文件**（实际输出见 §9）。
- 提交方式：`git add -- <逐枚显式列>` 之后 `git commit -F - -- <同一批逐枚显式列>`，**pathspec 写进 commit 本身**；没用 `git add -A`/`.`/`-a`。
- 没用过 `--amend`／`reset`／`rebase`／`stash`／`checkout .`／`restore`／`clean`／`merge`／`worktree`；**没有 push**；**没有删除任何东西**（突变那两发是"改后逐串换回＋`certutil` SHA256 比对"，不是 `restore`）。
- 不在 `internal/tools/**` 里 commit：`git show --stat 78988901` 的 14 枚清单里没有 `internal/tools`。工作树里别人的脏件（`internal/tools/bridge.go`、`subagent_197.go`、`subagent_222_test.go`、`.gitignore`、`design/**` 的 16 枚删除、`docs/evidence/s1/152-…`）**一律未碰、未 staged、未提交**。

## 7. AC#6 整包终态（票面口径：`go test ./cmd/wisp ./internal/... -count=1`）

起手 `2026-09-29 12:54 +0800`、终态 `12:56 +0800`（现跑 `date`，两枚时刻记在 `ac6-full-suite-start.txt`），原始输出 `ac6-full-suite.txt`，**枚举命令没接 `| head`/`| tail`**：

- **24 枚包 `ok`**、**2 枚包 `FAIL`（共 5 例）**、**5 枚包无测试文件**（`internal/agent/scheduler`／`internal/session`／`internal/speech`／`internal/streamkey`／`internal/watchdog`）。
- 与本票相关的那几枚：`cmd/wisp ok 97.157s`、`internal/config ok 1.523s`、`internal/perm ok 0.602s`、`internal/agent/approval ok 0.549s`、`internal/risk ok 8.731s`、`internal/tools ok 16.934s`。

**逐名红册＋分口径**（一枚都不顺手修）：

| 红例（逐名） | 所在包 | 归因口径 |
|---|---|---|
| `TestC21TableColourRowsMatchTokensCSS` | `internal/ball` | **历史在册 1 例**，台账 `A427`/`A428` 已记：另一队在工作树里删了 `design/assets/tokens.css`（`git status` 现读仍是 ` D design/assets/**` 那 16 枚未提交删除），别人地界 |
| `TestApprovalCardViewJSONKeysMatchFrontendTypes` | `internal/panel` | **历史在册 4 例**之一，同上＋面板契约字段差，别人地界 |
| `TestComposerContractTypesMatchFrontend` | `internal/panel` | 同上 |
| `TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme` | `internal/panel` | 同上 |
| `TestC21DesignTokensFourWayAgree` | `internal/panel` | 同上 |

⚠ 我**没有**读 `design/**` 与 `frontend/**` 的任何内容（派单禁读），上表的归因口径来自台账 `A427`/`A428` 的在册记录＋`git status` 的删除标记，不是我去比对了那些文件。

**基线对照（证明"零新增红"不是推断）**：动手前 `12:27` 跑过 `./internal/config ./internal/perm ./cmd/wisp` = 三包全 `ok`；交件前 `12:50→12:51` 同样三包仍全 `ok`；整包 5 枚红例的名册与 `A427`/`A428` 记的历史在册名册**逐名一致、枚数一致（1＋4）**。`internal/risk` 的 `TestResolvePerCallBudget` 这一发在 `ok 8.731s` 里绿（`A426` 记它是争用型假红；本轮没为它放宽任何断言，`thresholds.go` 一字节没动）。

## 8. 没做完的格、卡在哪、以及要新增却没敢造的名字

1. **AC#5 只闭了"落库"半格**。"撤销必须**能被答复语法触发**"没做：答复语法与 `always`/`yes`/`no` 那一族同属 `internal/agent/approval/replies.go`＋`cmd/wisp` 的监听器面，票 219（理由框/撤销）与票 224（会话档）都点名同一格，且 `A428` 已把票 201 AC#5 整格转给 224——同一格三方各做一次会撞。本票不抢那半格；`A424` 第 1 条"持久、可审计、可撤销"里的**可撤销**要等答复语法那半落地才算真闭合。**这就是需要登记的那一条推迟**：事项＝"一条长期规则的撤销入口（答复语法）"，落点＝票 219／224 的答复语法面，本票已交付的部分＝`Manager.SetAllowedDirs` 有真写路径且有单测（此前全仓零调用者零用例），未交付＝触发它的那枚入口，阻塞原因＝同一格三方重叠、新增答复动词属别的票的判据表。⚠ 按 `AGENTS.md` §1.1，`DEFERRED(D-xx)` 标记与 `SPEC-12 §5` 的 1:1 登记要动 `docs/specs/**`（本票禁动），所以这一条**没有**落成代码里的 `DEFERRED` 标记，只落在本节与交回正文里——请编排者决定是否入 `A##`（`deferred-work-must-be-registered` 那规要的是台账，不是我的证据件）。
2. **没有做票 223 的活**（具名以澄清边界）：没给 `CheckAndReload` 接生产轮询、没给 `ConfirmLocked` 找赋值点。合并读盘只用于**写**，不用于把磁盘值装进内存；包级与进程级各有一枚用例钉住这条边界。⚠ 由此留下一个**已知、非本票**的形状：一次带分叉的写会把文件的新 mtime 留在"未被认领"状态，票 223 接上轮询后那一发轮询会读回它、并对那枚手改作 D36 裁决（放宽＝拒绝、收紧/中性＝热应用）。这正是"226 必须先于 223"要产生的效果，不是缺陷。
3. **想报但没有造的名字**：本票没新增任何导出标识符、没新增 C17 方法名。`mergeWrite`／`diffKeyPaths`／`diffStruct`／`diffMap`／`configFieldName`／`truncateKeys`／`readConfigFile`／`keyFSAllowedDirs`／`keyRiskPermissionMode` **全部是包内未导出**。如果后续票想把"这次写实际改了哪几枚键"变成**跨包可读的返回值**（例如让 `cmd/wisp` 把那枚键名册直接打到 stdout、而不只进日志），那要新增的是一枚导出方法或导出结构体——**本票不自造，摆在这里请批**。现用的报告通道是 `slog`（与 `manager.go` 里 D36 那几行同一条；`cmd/wisp/logsink.go:153` 把它 tee 到 stderr 与日志件），AC#2 的两枚用例断的就是这行日志。
4. **一处甲的固有残缺，明写在此**：受守卫的写仍会把文件重排成 canonical 形，**注释与手排键序不保留**（`MarshalCanonical` 是结构体序列化，票面甲那一格自己就点了这一代价）。值不丢，格式丢；这一句同时写进那两行日志，所以用户从日志读得到，不只从源码。要连格式一起保住（文本级补丁式改写）是丙的另一个量级，本票没做。

## 9. `git diff --cached --name-only` 的实际输出

- **起手**（`12:26`，动手前）：**空**（索引干净）。
- **commit `78988901` 之前**（`12:53`，`git add -- <14 枚显式列>` 之后的实际输出，逐字）：

```
.scratch/wisp/probes/226/r1/cmd-wisp-ac1-first-run.txt
.scratch/wisp/probes/226/r1/cmd-wisp-ac1-green.txt
.scratch/wisp/probes/226/r1/config-perm-named-green.txt
.scratch/wisp/probes/226/r1/d22scan.txt
.scratch/wisp/probes/226/r1/gate-config-perm-cmdwisp.txt
.scratch/wisp/probes/226/r1/mut-allowdirs-no-merge.txt
.scratch/wisp/probes/226/r1/mut-permmode-no-merge.txt
cmd/wisp/always_write_no_clobber_226_test.go
docs/evidence/s1/226-config-write-no-clobber-r1.md
internal/config/allowdirs.go
internal/config/loader.go
internal/config/permmode.go
internal/config/writeguard.go
internal/config/writeguard_226_test.go
```

  ⇒ 索引里**没有别人的东西**，所以这一发可以提；提交后 `git show --stat 78988901` 实测就是这 14 枚、`2049 insertions(+) / 28 deletions(-)`。
- **第二枚 commit 之前**（本证据件 §3–§9 与 AC#6 那两枚读数件）：实际索引输出与 sha 记在 §10。

## 10. 第二枚 commit 实测

- **pathspec（写进 `git commit -F - --` 本身）**：`docs/evidence/s1/226-config-write-no-clobber-r1.md`、`.scratch/wisp/probes/226/r1/ac6-full-suite.txt`、`.scratch/wisp/probes/226/r1/ac6-full-suite-start.txt`。
- **提交前 `git diff --cached --name-only` 的实际输出**：`git add` 之前**为空**；`git add -- <上面三枚>` 之后逐字为
  `.scratch/wisp/probes/226/r1/ac6-full-suite-start.txt` / `.scratch/wisp/probes/226/r1/ac6-full-suite.txt` / `docs/evidence/s1/226-config-write-no-clobber-r1.md` ⇒ 索引里没有别人的东西。
- **sha `361b314a`**（committer 时间 `2026-09-29T12:58:45+08:00`），`3 files changed, 224 insertions(+), 9 deletions(-)`。
- ⚠ 两件需具名的外部事实：① 我的两枚 commit 之间，编排者自己提交了 `6cf2cc60`（`12:57:40`，`A431` 那枚）——我现跑 `git show --stat 6cf2cc60` 核对它的文件清单＝`docs/reports/pending-and-issues.md`／`docs/reports/HANDOVER.md`／`docs/reports/missing-features-2026-09-28-v3.md`＋一枚工单，**不含本票任何一枚文件**，所以它没带走我的活、我也没带走它的；② 那一发之后盘上 HEAD 已不是本腿起手锚点 `4540d5b1`，本件所有行号与读数都取自我自己这两枚 commit 的盘上状态（`78988901` 前后），不是起手锚点。
- 本件（§10 回填）在第三枚 commit 里；那一发的 pathspec 只有这一枚 `.md`，提交前索引实测输出同样先为空。
