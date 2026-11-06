# 票 345 验收腿 345-v2 · 证据件（2026-11-06）

- 任务：非实现者对票 345（迁移-on-读挪出读路径，甲形）的对抗验收；同批核销 wave35 派单里的 `W345`。
- 本件作者＝腿 345-v2。本腿**零产码改动**：被裁的实现（`2d3ad35e9`，2026-11-02）另有实现者（345-r1），本腿与其无同人关系 ⇒ 本件可作 `SPEC-12 §4.3` 意义上的非实现者裁决输入。
- 起手锚：共享树 HEAD `d0a230547`；派单锚 `92d5a33d1` 现验为其祖先（`git merge-base --is-ancestor` rc=0），两树之间唯一漂移＝`docs/reports/pending-and-issues.md`（台账追加，零码漂移；`git diff --name-only 92d5a33d1 d0a230547` 唯一行）。
- 隔离（rule 21）：`git clone --shared` 两棵仓外树（临件只建不删，落点符合 D:\tmp\wisp\ 位置规矩）：
  - after 树＝`D:\tmp\wisp\20261106\345v2\repo-after` @ `d0a230547`（含 build/、third_party/ junction 复用）
  - before 树＝`D:\tmp\wisp\20261106\345v2\repo-before` @ `063dcd4d1`（＝`2d3ad35e9~1`，345 落地前的基线）
  - 原始读数：`cmdwisp-full.log`、`cmdwisp-fails.txt`、`testnames-before/after.txt`、`zz_probe_345v2_test.go`
- 门禁环境：Git Bash；`export PATH="$PWD/build:$PWD/third_party/sherpa-onnx:$PATH"`；`unset GOTMPDIR TMPDIR TEMP TMP`；`CGO_ENABLED` 未设（实测 unset；⛔ 不置 0）。

## 为什么 wave35 还会派 W345（陈旧分诊的证据链）

1. `A1395`（提交 `c9ef940c8`，2026-11-02）＝机主明示「345＝做（甲形：读纯读、迁移搬到显式写）」⇒ 票面 `AC#5` 人工批准门**当时已开**。
2. `2d3ad35e9`（2026-11-02 12:43）＝345-r1 落地：9 件、+357/−92（本腿 `git diff --numstat 063dcd4d1 2d3ad35e9 -- internal/config` 复算＝8 枚码面＋1 枚证据件，与 A1402 记录一致）。
3. `A1402`（提交 `ea41fcab7`）＝编排者收件通过（自跑 `go test ./internal/config` 绿）。
4. `A1409`（提交 `2806bb9f0`）＝分诊自己写明：「345 分诊记『待人』→ 实则已由 345-r1 交付、待 345-v1 验收」。
5. `A1571`（2026-11-06）的 23 枚「待人」清单把 345 照旧抄入 B 段，wave35 随之重派 `W345` ⇒ **重复派单**。本腿按验收腿处置，不重做已落地的事。

## AC#0 起手现跑（after 树 @ `d0a230547`，尺逐枚可重跑）

| 票面落点 | 现跑读数 |
|---|---|
| `loader.go:126`（读管线里的迁移） | `applyMigrations` **全包零命中**（仅 3 处历史注释提及）；`ver != SchemaVersionCurrent` 支现调 `migrateRaw(raw, ver)`（`loader.go:156`）＝**纯内存**，无 path 参数 |
| `migrate.go:93`（`winsec.PrivateFile(backup,…)`） | 现居 `migrate.go:145`，唯一调用者＝`backupOldSchema`（`:128`），而 `backupOldSchema` 唯一调用点＝`SaveFile`（`loader.go:299`）＝**显式写侧** |
| `migrate.go:96`（`atomicWrite(path, data)` 迁移写回） | 迁移链（`migrateRaw`）内**零磁盘写**；`atomicWrite` 生产调用点唯一＝`loader.go:302`（`SaveFile` 本体）；定义在 `parse.go:196` |

自加尺（写点普查）＝`git grep -nE "atomicWrite|PrivateFile|os\.(WriteFile|Create)" -- internal/config`：
**读路径堆（经 LoadFile/readConfigFile 可达）：空**。写路径堆：`PrivateFile`@migrate.go:145、`atomicWrite`@loader.go:302、`os.CreateTemp`@parse.go:201（全部只在显式写路径可达）；其余 16 处命中全为 `*_test.go` fixture（逐名见 `git grep` 原始输出，本腿留档）。

## AC#1 三形代价表

表本体在实现者件 `docs/evidence/s1/345r1-migration-off-read.md` §AC#1（甲/乙/丙三列，含每形动几枚文件、打红几枚用例）。本腿复验其**甲形**落地与表承诺一致：读路径纯读（AC#0）、迁移与备份只在 `SaveFile`、`D36` 三档与票 226 判语未动（`git diff 063dcd4d1..HEAD -- docs/PLAN.md docs/specs`：`PLAN.md` 零行；`SPEC-01` 的 +3 行属票 24 依赖白名单、机主批 `A1445`、2026-11-03，与本票无关——`git log 063dcd4d1..HEAD -- docs/specs` 唯二提交 `1c66fee02`/`05ec59ef9` 均为 24-PDF）。

## AC#2 删用例普查（⛔ 删/裸 Skip/两可断言）

- 测试函数名普查（`git grep -h "^func Test" 063dcd4d1|2d3ad35e9 -- internal/config`）：**124 → 124**。removed＝1 枚、added＝1 枚，1:1 换钉具名：`TestNothingWrittenPromiseOnlyHoldsWhenSchemaIsCurrent` → `TestNothingWrittenStaysUntouchedEvenWhenSchemaIsStale`（与实现者件 §AC#2 逐字一致）。
- `t.Skip` 普查：before＝4 处（`c26_seam_posix_125_test.go` 3 处 Skipf＋`tempdir_resolved_124_test.go` 1 处注释）＝after 同 4 处，**零新增、零删除**（符号链接夹具的既有 Skipf，与 345 无关）。
- 逐名 `-v` 复跑 after 树（rc=0，9 枚全 PASS）：`TestMigrateV1Fixture`／`TestMigrateNoVersionKeyAssumedV1`／`TestMigrateCorruptFileUntouched`／`TestMigrateUnsupportedVersionUntouched`／`TestAC2MigrationBackupLandsPrivate`／`TestUnwiredGuardFiresBeforeAnyFileWrite`／`TestNothingWrittenStaysUntouchedEvenWhenSchemaIsStale`／`TestNothingWrittenLeavesFileUntouchedWhenSchemaIsCurrent`（配对控制）／`TestHandDeletedSchemaVersionOnTheWriteBranchIsForcedBackAndUnclaimed`。

## AC#3 用户可见代价：注释在不在（真读文件字节，两态各一发）

探针＝`zz_probe_345v2_test.go`（纯记录仪，两树同源同跑；种入 352 字节、2 枚行首注释＋1 枚行尾注释、手挑键序 `[ball]` 在 `[app]` 前、声明 `schema_version = 1`）。

| 阶段 | before 树（`063dcd4d1`＝旧码，迁移在读路径） | after 树（`d0a230547`＝甲形） |
|---|---|---|
| 一次 `LoadFile` 纯读后 | **bytes 352→409（被重排）**；行首注释 2→**0**；键序倒（ball_before_app=false）；schema 被钉走；**`.bak-1` 凭空出现** | **bytes 352→352 逐字相同**；注释 2/2 全活（行尾注释随字节同一性一并存活）；键序不变；仍声明 v1；**无 `.bak-1`** |
| 显式 `SaveFile` 后 | 3386B canonical，注释 0 | 3386B canonical，注释 0（＝票 226 甲形既有写时代价，本票未扩未修）；`.bak-1` 出现且 **=352B 原文逐字**（含全部 3 枚注释） |
| 写后再读 | 稳定 | 字节不动 |

⇒ 票面那句话现在对任意 schema 版本的文件字面成立："程序没写任何东西"不再带"仅当你的 schema 已是最新"的暗条件；代价面（写时丢注释）如实入表。

## 出牙（双向，本腿自跑）

- **红移正控**：把 HEAD 的 4 枚读纯钉（`migrate_test.go`/`private_acl_windows_test.go`/`writeguard_handedit_227_test.go` 的相关用例）放进 before 树（一枚临时编译 shim `zz_migrateshim_345v2.go` 让 `migrateRaw` 名字可解析——before 码里该符号本不存在，这正是被裁的变更本体）⇒ **4/4 FAIL**，报错句逐字点名："a read must not rewrite config.toml - migration-on-read is gone (ticket 345)"、`migrate_test.go:77`；"a read must not create a backup"、`:156`；另 `TestAC2MigrationBackupLandsPrivate`、`TestNothingWrittenStaysUntouchedEvenWhenSchemaIsStale` 同红。同 4 枚在 after 树 ⇒ 全 PASS。⇒ 钉子对"改前态真红/改后态真绿"双向咬合，不是装饰。
- M1/M2 overlay 变异＝实现者件 §变异出牙 已载（M1 打红 4 枚写侧钉、M2 打红 4 枚读纯钉）；本腿未重做，仅引用不背书。

## AC#4 门禁（四发，after 树；红名册逐名作差）

| # | 尺 | 读数 |
|---|---|---|
| 1 | `go test ./internal/config -count=1 -timeout 420s` | `ok … 2.546s`，rc=0 |
| 2 | `go build ./...` | rc=0 |
| 3 | `gofumpt -l internal/config`（`/d/work/base/gopath/bin/gofumpt.exe`） | **空**（rc=0）。附：全仓扫描仅 `.scratch/wisp/probes/**` 下 20 枚存量探针被点名（含刻意的不合格负控 `gofumpt-negctl/probe.go`），非构建面，非本票 |
| 4 | `sh scripts/d22scan.sh` | 正控 `PASS=35 FAIL=0 SKIP=0`；实扫 **421 枚生产件 clean**（bans #1-8 全 examined 非零），rc=0 |

红名册作差：四发门禁**零红** ⇒ 无名册可作差。

**附：`go test ./cmd/wisp -count=1`（HEAD 整包，环境完整时跑通 776s）＝47 红，逐名归因（原始名单 `cmdwisp-fails.txt`）：**
- **配置面（config/reload/migrat/schema/writeguard/allowed/secret 关键词扫描）：0 红**。345 邻接的四枚句子/装钉用例现跑全 PASS（22.3s）：`TestTicket223RunArmsTheReloadTick`／`TestTicket223FailureSentencesAreDistinct`（4 子例）／`TestTicket223R2FailureSentenceRouting`（9 子例）／`Test233ReloadTierEventReachesAProductionConsumer`。
- 语音族（15e/26P/26R3/381/382/390/392/394/395/41a2/WireR1/16r1/12r2 等约 40 枚）＝同一环境根因：**"onnxruntime 1.17.1 … a second onnxruntime.dll is shadowing the one on PATH (System32 before PATH)"**——本机 System32 混入旧版 onnxruntime（174-c1 时代已记过的机器环境问题），非码面。
- `TestTicket256ResidentGateOptionsFieldSetIsTheFiveItClaims`（1 枚）＝**基线前既存**：其钉的多余 `Answer:` 字段在 `063dcd4d1` 的 `resident_approval_windows.go` 里已在（票 375-r1 `b63d56413` 引入；后续 `9468fccad` 证据件已具名"pre-existing 256-roster drift"）。
- 其余 5 枚（AC1 resident log ×3／AC228／AC4 focus／AC247）＝同环境族的 live 窗口/落盘测试；其测试文件与生产面在 `063dcd4d1..HEAD` 窗口均非 345 所触（345 只触 internal/config 3 产码件＋5 测试件，逐字清单见上）。
- **判读：47 红无一可归因 345。** cmd/wisp 整包变绿不是本票门（属 wave 车道与 375/256/语音票），此处仅作逐名归因存档。

## AC#5 人工批准门

- 机主明示一：`A1395`（2026-11-02，"都按推荐"＝345 做甲形）。
- 机主批复二：`A1571`（2026-11-06，B 段含 345 全批；Q-52 等同批解扣）。
- 契约面：`PLAN.md`（含 D36 节）零字节改动（见 AC#1 段）；票 226 已勾判语未翻。

## AC#6 纪律格

- 本腿写面＝`docs/evidence/s1/345-v2-verdict.md` 一枚（`internal/config/**` **零改动**——`git status --porcelain -- internal/config` 开工与收工均为空；生产码面与 `2d3ad35e9` 逐字节一致：`git diff 2d3ad35e9 HEAD -- internal/config/loader.go internal/config/migrate.go internal/config/writeguard.go` 空）。
- ⛔ 未翻票面任何框、未改 `Status:`、未加 `-done`（属编排者/后续收口面）；⛔ 未动台账；⛔ 未 push；`add` 与 `commit` 两处显式 pathspec 手写；临时件全部在 `D:\tmp\wisp\20261106\345v2\`（两棵 clone＋探针＋日志，只建不删，清单如篇首）。

## 判定表（与票面 AC 1:1）

| AC | 判 | 依据（本件节） |
|---|---|---|
| AC#0 起手现跑 | PASS | 三落点现跑＋自加尺：读路径写点堆＝空 |
| AC#1 三形代价表 | PASS | 甲形落地与表一致；D36/226 判语零触碰 |
| AC#2 删用例换绿 | PASS | 124→124、1:1 具名换钉、Skip 零增、9 枚逐名 PASS |
| AC#3 字节凭据 | PASS | 两态真读字节三阶段表（读纯/备份逐字/写时代价） |
| AC#4 门禁四发 | PASS | 四发全绿零红；cmd/wisp 47 红逐名归因全数非 345 |
| AC#5 批准门 | PASS | A1395＋A1571 双源；契约面零字节 |
| AC#6 纪律格 | PASS | 本腿零产码、单证据件、未 push、pathspec 显式 |

## 结论

**PASS ＝ 7/7 全 PASS。票 345 可结（closeable=YES）**：行为变更已落地（`2d3ad35e9`）并经本腿非实现者双向出牙复验，四发门禁全绿。给编排者的收口动作（本腿⛔ 不代做）：翻 7 格 → `Status: done` → 文件名加 `-done` → 台账记一笔（345 由 `2d3ad35e9`＋本件收口；wave35 的 W345 属陈旧分诊重派，零活）。wave35 同波另两腿（W350/W150）不在本件射程。
