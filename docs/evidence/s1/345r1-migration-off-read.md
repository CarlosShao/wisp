# 票 345 落地腿 345-r1 · 证据件（2026-11-02）

- 任务：把"迁移-on-读"挪出读路径——读＝纯读不写盘；迁移只在显式写时发生（＝票面 **甲形**，机主已明示＝做）。
- 派单存档：`.scratch/wisp/probes/orch/20261102/txt/345r1-dispatch.txt`
- 本件作者＝写腿 345-r1（实现者自证；终裁须非实现者，按 `SPEC-12 §4.3`）。
- 原始读数文件（仓外临时件，只建不删）：`D:\tmp\wisp\20261102\345r1\00..09-*.txt`

## 起手三核（读数）

| 尺 | 读数 |
|---|---|
| HEAD | `c9ef940c853d417d528e0c47e1b0887fcac6cf14`（≥ 锚 `c9ef940c8` ✓，且恰为该枚） |
| 索引（`git diff --cached`） | 空＝staged 零 ✓ |
| `git status --porcelain -- internal/config`（开工时） | 空＝无他人脏 ✓（工作树他处脏属 337r2 的 `cmd/wisp/run.go`，未碰） |

## AC#0 三落点现跑（具名提交树 c9ef940c8）

尺＝`git show c9ef940c8:<file>` 取行。票面行号按规矩只当读数：

| 票面落点 | 现跑读数 |
|---|---|
| `loader.go:126` | **漂移**。:126 现为 parse 错误分支一行；真调用点＝**`loader.go:150`**：`raw, err = applyMigrations(path, raw, ver)`（逐字仍在读管线 `readConfigFile` 内，`if ver != SchemaVersionCurrent` 支） |
| `migrate.go:93` | 逐字命中：`if err := winsec.PrivateFile(backup, raw, 0o600); err != nil {` ← 真写 `.bak` |
| `migrate.go:96` | 逐字命中：`if err := atomicWrite(path, data); err != nil {` ← 迁移写回 |

自加尺（"读路径还有没有别的写盘点"）＝
`git grep -nE "atomicWrite|PrivateFile|os\.(WriteFile|Create)" internal/config`，逐枚分堆点名：

**读路径堆（经 LoadFile/readConfigFile 可达，改前＝真写盘）**
- `migrate.go:93` `winsec.PrivateFile`（备份写）——本票消除；
- `migrate.go:96` `atomicWrite`（迁移写回）——本票消除。
- 读路径其余磁盘调用仅 `os.ReadFile`（loader.go:89，纯读）与 stat/peek（纯读）⇒ **除这两枚外读路径无别的写盘点**（本尺结论）。

**写路径堆（显式写专属，保留）**
- `loader.go:286`（改前读数）`SaveFile` 内 `atomicWrite`＝显式写本体 ✓；
- `parse.go:196-215` `atomicWrite` 实现（`os.CreateTemp`＋`winsec.SealFile`＋rename）＝写侧 helper；
- `permmode.go`：只指路（注释），磁盘写经 `mergeWrite→SaveFile`，无自有写点；
- 测试文件 16 处匹配全部是 fixture 写（`*_test.go`），不计生产写面。
- `winsec.SealDir`（SaveFile 封目录，票 124）在显式写路径上，读路径不可达 ✓。

## AC#1 三形代价表（⛔ 未出表未动码；表先落 `D:\tmp\wisp\20261102\345r1\AC-1-reference-cost-table.txt`，本段为存档版）

| 维度 | **甲＝迁移只在显式写（读本票所办）** | 乙＝读仍迁移但文本级保注释/键序 | 丙＝维持现状＋响亮报告 |
|---|---|---|---|
| 行为语义 | 读＝纯读不写盘；迁移与 `.bak` 在 `SaveFile` 落地 | 读仍写盘；迁移改注释保留式文本 patch | 读仍写盘；只把话说响 |
| 买到的价值 | "读一次配置就重排你的文件"这一根因**移除**；手改文件（含注释/键序）跨读永久存活 | 部分缓解症状（注释活）；**根因不除**（读仍写盘） | 已部分到位（票 227 甲形＝日志说实话） |
| 动几枚生产文件 | 3（`loader.go`/`migrate.go`/`writeguard.go` 句子） | ≥3 且 `migrateV1toV2` 重写为文本 patch 引擎（go-toml/v2 无注释保留 marshal，需新工件） | 0-1（只加日志） |
| 要打红/重钉的既有用例 | 重钉 5 枚（见 AC#2 逐名）；2 枚只修注释；其余全绿 | `migrate_test.go` 全部＋`writeguard_226/227/341` 里所有"canonical re-lay 丢注释"断言**整体翻案**＋补丁引擎需新常驻用例 | 零枚 |
| 破 `D36` 三档生效级别？ | 不破（迁移是加载流程，不是 section 生效档） | 不破 | 不破 |
| 破票 226 已勾判语？ | 226 的 merge/保值/报告判语全绿；仅 227 P5′ 那条"nothing-written 边界带 schema 条件"的判语按票 345 换钉（机主明示射程内） | 226 丙形"文本级 patch"量级被票面点名（`migrateV1toV2` 的重排＝其"管线顺序承诺"要连读路径一起重议） | 无 |
| 用户可见代价（量在 AC#3） | 注释存活到**第一次显式写**；写时仍丢（＝226 甲形既有代价，本票不扩） | 注释写时也保住 | 注释继续在读时丢 |

选型自裁（按 ⓪⓪⓪ 铁律：最有价值＋主流 harness 方向）＝**甲**：唯一一根除"读即写盘"根因的形；乙的复杂度买不到"读纯"这一半价值且工程面更大；丙不解决问题。机主明示亦为甲。

## AC#2 改钉清单（⛔ 删用例＝零；⛔ `t.Skip`＝零；⛔"两种都算对"＝零）

**重钉新语义的 5 枚（逐名，换钉内容点名）：**

1. `migrate_test.go` `TestMigrateV1Fixture`——旧钉"LoadFile 读时迁移并落盘（.bak-1 出现＋文件就地升级）"；新钉"读后文件字节不变＋无 .bak-1（新增断言）；显式 `SaveFile` 后才出现 .bak-1＝迁移前原文逐字＋文件带当前版"。**保留未动**：v1→v2 全部值映射断言、reload 幂等、备份字节相同。
2. `migrate_test.go` `TestMigrateNoVersionKeyAssumedV1`——"无 schema_version＝当 v1"的读数**不变**；`.bak-1` 出现时机由读→写（断言随迁）。
3. `migrate_test.go` `TestMigrateUnsupportedVersionUntouched`——直测调用由 `applyMigrations(path, …)` 改钉 `migrateRaw(…)`（纯链新名）；"错误点名卡住的版本号"断言不变。
4. `private_acl_windows_test.go` `TestAC2MigrationBackupLandsPrivate`——旧钉"LoadFile（＝每次启动的读）写 .bak-1 并改写文件，两者私有"；新钉"LoadFile 后**无**备份（读纯）；`SaveFile` 落迁移备份后 `.bak-1` 与 config.toml **私有**、备份＝原文逐字"。票 95 AC#2 的隐私断言一字未削，只挪触发者。
5. `writeguard_handedit_227_test.go` `TestNothingWrittenPromiseOnlyHoldsWhenSchemaIsCurrent` → 更名重钉为 `TestNothingWrittenStaysUntouchedEvenWhenSchemaIsStale`——旧钉"nothing-written 只对 schema 已Current 的文件成立，落后版文件已在复读中被迁移重排（字节/mtime 动、注释死）"；新钉"nothing-written **无条件**真不写：字节/mtime 不动、注释活、无 .bak-1、schema_version 行仍缺"，并新增反向断言"日志⛔ 再出现旧条件子句"。配对控制 `TestNothingWrittenLeavesFileUntouchedWhenSchemaIsCurrent` 未动（仍绿）。

**随迁的生产句子（钉它的用例＝上面第 5 枚，同步改钉）：** `writeguard.go` nothing-written 分支 INFO 句子旧半句"stayed put only because its schema_version was already current; a file behind … had already been re-laid in canonical form by the re-read"（改前已不真）→ 新半句"they stay put whatever the file's schema_version is: the re-read migrates an older file in memory only…"。

**只修注释、断言零动的 2 枚文件：** `writeguard_unnailed_341_test.go`（P5 机理描述改口：备份/改写现由 `SaveFile` 落；其全部断言在新语义下逐条重验为真：mtime 由**这次写**移动、.bak-1＝写前逐字原文含 `hand note`、⛔ schema_version 不被认领）；`unwired_test.go` `TestUnwiredGuardFiresBeforeAnyFileWrite` 头注释（"迁移是唯一会写盘的加载步"→"加载管线已不写盘，no-side-effects 由构造成立"；断言全留且仍绿）。

**未动的钉旧方向候选（现跑确认仍真）：** `TestMigrateCorruptFileUntouched`（"left untouched"句子＋无备份——新语义下更强）、`TestMigrateHigherIntermediateVersionRejected`、`writeguard_226_test.go` 全系（fixture 均 current-schema，不经迁移支）。

## AC#3 用户可见代价：注释在不在（真读文件字节各一发）

尺＝仓外探针 `D:\tmp\wisp\20261102\345r1\zz_probe_345r1_test.go` 经 `-overlay` 挂进包跑真 `LoadFile/SaveFile`，逐阶段 `os.ReadFile` 数注释行/键序（原始件 `01-ac3-before.txt`/`04-ac3-after.txt`；探针未入仓，跑后 `internal/config` 无 zz_* 文件 ✓）。

种入文件＝315 字节、3 枚用户注释（`hand:` 行）、手挑键序 `[ball]@20` 先于 `[app]@59`、声明 `schema_version = 1`。

| 阶段 | 改前（HEAD 旧码） | 改后（本腿） |
|---|---|---|
| 一次 `LoadFile` 纯读后 | **bytes 315→419（被重排）**；注释 3→**0**；键序变（`[app]@20`/`[ball]@42`）；schema 被钉成 2；**`.bak-1` 凭空出现**；mtime 动 | **bytes=315 逐字相同**；注释 **3/3 全活**；键序不变；仍声明 1；**无 `.bak-1`**；mtime 不动 |
| 显式 `SaveFile` 后 | （3381B canonical，注释 0） | 3381B canonical，注释 0（＝226 甲形既有代价，写时丢，如实报）；`.bak-1` 出现且＝315B 原文逐字、3 枚注释都在里面 |
| 写后再读 | 稳定 | 字节不动 ✓ |

⇒ 本票承诺面：*"程序没写任何东西"现在对任意版本的文件字面成立*；代价面：*注释仍会在显式写时被重排掉*（写时不丢注释＝乙/丙别形，票面 AC#1 已量，非本票所办）。

## 零回归（internal/config 包，含改钉用例）

- 改前基线：`go test -count=1 ./internal/config/` → `ok`（`00-baseline-gotest.txt`，FAIL=0）。
- 改后：`go test -count=1` → **FAIL=0**（`03`，`ok 1.732s`）；`go test -race -count=1` → **FAIL=0**（`05`，`ok 4.672s`）。
- 改钉 11 枚逐名 `-v` 全 PASS（`06`：migrate 系 5＋227 系 2＋341-P5＋acl 2＋unwired 1）。
- `go vet ./internal/config/` 净；`gofumpt -l internal/config/*.go` 空（绝对路径 `/d/work/base/gopath/bin/gofumpt.exe`）。
- 边界自认：`cmd/wisp` 整包**未追**（337r2 在写、共享编译面；本腿未动任何导出签名：`LoadFile/SaveFile/Manager` API 逐字不变，`applyMigrations→migrateRaw/backupOldSchema` 均为包内未导出改动）。`go build ./cmd/wisp` 由编排者波次门禁跑。

## 变异出牙（仓外 `-overlay`＋还原核）

- **M1**＝`SaveFile` 删 `backupOldSchema` 钩子（overlay `loader.m1.go`）→ 打红 4 枚写侧钉：`TestMigrateV1Fixture`(:122 backup must exist)、`TestMigrateNoVersionKeyAssumedV1`(:162)、`TestAC2MigrationBackupLandsPrivate`、`TestHandDeletedSchemaVersionOnTheWriteBranchIsForcedBackAndUnclaimed`(:251)（`07-m1-red.txt`）。
- **M2**＝读管线恢复"迁移即写盘"（overlay `loader.m2.go`）→ 打红 4 枚读纯钉：`TestMigrateV1Fixture`（读后字节不变/无 .bak）、`TestMigrateNoVersionKeyAssumedV1`、`TestAC2MigrationBackupLandsPrivate`（读后无备份断言）、`TestNothingWrittenStaysUntouchedEvenWhenSchemaIsStale`（字节/mtime/注释三断言）（`08-m2-red.txt`）。`TestUnwiredGuardFiresBeforeAnyFileWrite` 在 M2 下仍绿＝如实读数（该 fixture 迁移在 dry-check 即拒，写盘行不执行），它的牙是 M1/构造级的，不冒充。
- **还原核**＝去 overlay 复跑整包 `ok 1.594s`、`git status internal/config` 仅本腿 8 枚文件 ✓（`09-restore-core.txt`）。

## 给编排者的过期读数（⛔ 本腿写面，未碰）

- `cmd/wisp` 三处注释引用 `migrate_test.go:123`（`config_reload.go:440`、`config_reload_223_test.go:795`、`config_sentences_223r2_test.go:16,75`）：该用例现居 **:166**（重钉致行移；用例名与钉的句子未动，票 223 语义全绿）。注释非代码，不影响编译。
- 旧错误文本 `config.toml migration rewrite` 随迁移写回从读路径消失；全包＋全仓 grep 无用例/句子依赖（`config.toml migration backup write` 仍在写路径保留）。
- 票面行号读数更新：`applyMigrations` 调用点 c9ef940c8 树内＝`loader.go:150`（票面 :126 过期）。

## 纪律格自查

- 写面＝`internal/config/**` 8 枚＋本件 `docs/evidence/s1/345r1-migration-off-read.md` ✓；⛔ 未碰 `cmd/wisp/`、`thresholds.go`、golden、`frontend/**`、`design/**`、票面、台账、`.scratch` 正式区；⛔ 未翻本票任何框；⛔ 未 push；`git add`/`commit` 两处显式 pathspec 手写字面路径；临时件全在 `D:\tmp\wisp\20261102\345r1\`（只建不删）。
