# 132-c2 — DB 备份副本「不经封条落盘」的落点与代价普查（只读腿）

腿：`132-c2`（只读普查，零产码，零编译）。票：`.scratch/wisp/issues/132-log-files-are-never-sealed-sealdirdir-has-zero-production-callers.md` §「编排者增量（09-29 17:4x）」。
锚点：本腿**取数时** `git rev-parse --short HEAD` = `795ed767`；交付时 HEAD 已被并行的产码腿推进到 `5759d7fb`。⚠ 两者之间的 8 枚文件（`git diff --name-only 795ed767..5759d7fb`＝两枚 issue 件、两枚 probe 件、两枚 `docs/reports/` 件、`internal/tools/subagent_197_test.go`、`subagent_222_test.go`）**不含本件引用的任何一枚源文件**（`internal/memory/open.go`、`internal/winsec/*`、`cmd/wisp/*`、`internal/observe/*`、`internal/models/*`、`internal/config/*`、`internal/secret/*`、`internal/agent/*`）⇒ 本件全部行号在新锚点上仍逐一对得上。票面增量段自称的起手锚点是 `a8f3c020`，本腿不依赖任何历史行号，全部现读。
取数方式：只用 Read / grep / ls / find / git；**没有跑过 `go test`／`go build`／`go vet`／任何编译**（产码腿 `235-r1` 独占读数期）。本腿在树里的唯一足迹＝新建 `.scratch/wisp/probes/132/c2/`（`git status --porcelain .scratch/wisp/probes/132/` ＝ `?? .scratch/wisp/probes/132/`），未 commit、未 push、未改动任何既有件。
本文件只做三件事：把「要接哪一跳」摊平、把两形（实测为**三形**）的代价具名、把「哪一形会撞已有契约句」指回原文行。⛔ 不含裁决，不含实现选择。

## §1 现量那一跳：`internal/memory/open.go` 备份路径逐行重抽（542 行文件）

尺：`Read internal/memory/open.go`（整段现读）+ `grep -n` 逐句核。下表左列＝**本腿现读行号**，右列＝票面增量段写的行号。

| 本腿现读 | 内容（逐字） | 票面写的 | 漂移 |
|---|---|---|---|
| `:493` `func (s *Store) backupPath(from, to int) string` | 名字渲染 | 未标 | — |
| `:494` `return filepath.Join(s.backupDir, fmt.Sprintf("%s.bak-%d-%d", dbFileName, from, to))` | `dbFileName="wisp.db"`（`:39`） | 未标 | — |
| `:497`–`:500` 函数 doc 首四行 | 「An existing backup is kept … must not destroy it」 | 票面 `:497-499` | **少一行**（承诺句实际跨 498–500） |
| `:501` `func (s *Store) backupDatabase(from, to int) error {` | 函数体起 | 票面 `:501`（§72 未标，`:533` 有标） | 一致 |
| `:502` `dst := s.backupPath(from, to)` | 目标名 | — | — |
| **`:503`** `if err := winsec.PrivateDirAll(s.backupDir, 0o700); err != nil {` | **目录加锁那一行** | 票面 `:503` | 一致 |
| `:506` `if _, err := os.Stat(dst); err == nil {` | **只比基名**的那一步 | 票面 `:505` | **+1 漂移** |
| `:507`–`:509` Warn「migration backup already exists, keeping it」+ `return nil` | 早退：基名存在则**整段不写** | 票面未提早退的范围 | ⚠ 见 §4 |
| `:511` `for _, suffix := range []string{"", "-wal", "-shm"} {` | **侧车命名规则**（唯一定义处） | 票面「dst+"-wal"／dst+"-shm"」 | 一致（措辞层级不同） |
| `:512` `src := s.dbPath + suffix` | 源名拼法 | — | — |
| `:514` `if _, err := os.Stat(src); err != nil { continue }` | 侧车**只在源存在时**才拷 | — | — |
| `:518` `if err := copyFile(src, dst+suffix); err != nil {` | 调用点 | 票面 `:518` | 一致 |
| `:527` `func copyFile(src, dst string) error {` | 体起 | 票面 `:527` | 一致 |
| `:528` `in, err := os.Open(src)` | 读端 | — | — |
| **`:533`** `out, err := os.Create(dst)` | **造出文件那一句**（`os.Create` = `O_WRONLY|O_CREATE|O_TRUNC`，mode 走 `0666 &^umask`） | 票面 `:533` | **一致** |
| `:537` `if _, err := io.Copy(out, in); err != nil {` | 流式拷贝 | 票面「`copyFile` 是 `io.Copy` 流式的」 | 一致 |
| `:541` `return out.Close()` | 无 sync/无 SealFile | — | — |

小结（只述事实）：票面**两处行号漂了**（承诺注释 497–499 → 实际 497–**500**；`os.Stat(dst)` `:505` → 实际 **`:506`**）；**`os.Create` 那一跳 `:533` 与目录加锁那一行 `:503` 都还在原位**。尺本身（票面第一行那把 grep）现跑仍命中**恰好 1 处**：

```
$ grep -rn 'os\.Create(\|os\.OpenFile(' --include='*.go' internal/memory internal/secret internal/agent | grep -v '_test.go'
internal/memory/open.go:533:	out, err := os.Create(dst)          # 计数 = 1
```

## §2 谁真的调到它：生产调用者 1 处，且**在真跑路径上**（不是「没人调」）

尺：`grep -rn "backupDatabase" --include="*.go" .` ⇒ **全仓 3 行命中**（`:328` 调用、`:497` 注释、`:501` 定义）⇒ **调用点＝恰好 1 枚**。

链路（逐跳现读）：

| 跳 | 位置 | 条件 |
|---|---|---|
| 1 | `internal/memory/open.go:328` `if err := s.backupDatabase(step.from, step.to); err != nil {` | 在 `migrate()` 的 `for _, step := range steps` 循环里（`:324`） |
| 2 | `internal/memory/open.go:191` `if err := s.migrate(); err != nil {` | `Open()` 无条件调用 |
| 3 | `cmd/wisp/run.go:384` `mem, err := memory.Open(s.dataDir)` | 在 `assembleRuntime()`（`run.go:329`）内，由 `runTextTask`（`run.go:154`，调用点 `:207`）驱动，入口＝`cmd/wisp/main.go:153` `cmdRun` |
| 4 | `cmd/wisp/providers.go:180` `store, err := memory.Open(io_.dataDir)` | `wisp providers probe <provider/model>` 的一条腿（`:160`–`:186` 现读：端点就绪之后才打开） |

会不会真跑起来（这条决定本票是「安全洞」还是「未接线」）：

- `SchemaVersionTarget = 2`（`internal/memory/schema.go:129`），链＝`{0→1, 1→2}`（`schema.go:168`–`:171`）。
- `migrate()` 在 `version == s.target` 时**早退**（`open.go:303`–`:306`）⇒ **稳态运行不写备份**。
- 全新数据根：`readSchemaVersion` 回 `(0, fresh=true, nil)`（`open.go:430`–`:432`）⇒ `planMigration` 出 2 步 ⇒ **第一次 `memory.Open` 就写两枚备份**（`wisp.db.bak-0-1`、`wisp.db.bak-1-2`）。
- 升级（v1 库跑 v2 二进制）：写 1 枚（`wisp.db.bak-1-2`）。
- 旁证（**不是本腿跑的，是既有测试的断言原文**）：`internal/memory/migrate_v1seed_test.go:110` 与 `internal/memory/schema_test.go:350` 把备份目录钉成**恰好** `["wisp.db.bak-0-1","wisp.db.bak-1-2"]` ⇒ 这两枚文件在真链路上确实会被造出来。
- 对照读数：`memory.Open` 的非测试调用者＝2 枚（`grep -rn "memory\.Open(" --include="*.go" . | grep -v _test.go`；另有 5 行命中在 `.scratch/wisp/probes/**` 的快照件，非产码）；测试侧引用它的文件＝17 枚。

⇒ **照实写：`backupDatabase` 有生产调用者，落在 `wisp run` / `wisp providers probe` 的真实落盘路径上**；它和 `winsec.SealDir` 那一族「写好但零调用者」不是同一枚形状。频率特征也照实写：**每枚安装／每级升级各触发一次，稳态零触发**。

## §3 落点形状的代价：票面写「两形」，现读是**三形**（丙形零新增契约面）

票面硬约束 1 说「落点只有两支」。**这条在本锚点上不成立**：仓里已经存在一枚「先封、后流式写」的**生产现役写法**，它只用现有导出名，不需要新契约面。三形并排：

| 形 | 形状 | 现读出处 | 代价（具名） |
|---|---|---|---|
| **甲** 整档进内存 | `os.ReadFile(src)` → `winsec.PrivateFile(dst, data, 0o600)` | 写入口＝`internal/winsec/winsec.go:86` `func PrivateFile(path string, data []byte, perm fs.FileMode) error`；体内 `:122` `if _, err := f.Write(data); err != nil {`＝**一次 Write 全量** | 峰值内存＝**单枚文件的全量字节**。`backupDatabase` 一次调用最多拷 3 枚（`:511` 的 `""`/`-wal`/`-shm`），逐枚串行 ⇒ 峰值≈`max(wisp.db, wisp.db-wal, wisp.db-shm)`，**而 `wisp.db` 没有任何字节上限**（下段）。语义上等价今日：`PrivateFile` 用 `O_WRONLY|O_CREATE|O_TRUNC`（`:87`），与 `os.Create` 同 flags ⇒ **不改用户可见行为** |
| **乙** 新导出流式入口 | 例：`winsec.CreatePrivate(path) (*os.File, error)`（先造空档、封、再交句柄给调用方 `io.Copy`） | 现无此名。`winsec` 导出面现量＝**9 个函数 + 3 个 sentinel error + 3 个类型**（`grep -rn "^func [A-Z]\|^type [A-Z]\|^var [A-Z]" internal/winsec/ --include="*.go" \| grep -v _test.go`＝**15 行**，其中函数 9 枚：`SetPathResolver`/`PathResolverInstalled`/`ResolvePath`/`PrivateFile`/`PrivateFileExclusive`/`SealFile`/`PrivateDirAll`/`SealDir`/`RemoveUnlinked`） | **契约级**：票面硬约束 1 已明写「实现腿不许自己造，要造就先回编排者、由我落 `A##` 并摆给 owner」⇒ 本腿只登记射程，不裁可行性。底层机制其实**已存在**（`sealHandle`＝`winsec_windows.go:328`／`winsec_other.go:45`，注释逐字「seals a file that is open but still empty」），只是**没导出** |
| **丙** 现役 idiom：`os.Create` → `SealFile` → 流式写 | 造空档（零字节）→ `winsec.SealFile(dst)` → `io.Copy` → `Close` | **同仓生产先例**：`internal/config/parse.go:201` `os.CreateTemp` → `:215` `winsec.SealFile(tmpName)` → `:219` `tmp.Write(data)`，注释逐字「Seal before a single content byte, not after」；另一枚 `internal/secret/migrate.go:174` `winsec.SealFile(backupPath)`（对**已存在**文件的修腿） | **零新增导出名**。代价是把 `winsec` 内部的两条保证搬到调用方：① 封失败必须**不写字节并删档**（`winsec.go:117`–`:121` 的 `f.Close(); os.Remove(resolved)` 现在要手写）；② 名字二次解析的 TOCTOU 由 `SealFile` 自己做（`winsec.go:137` `resolveString`，`winsec_windows.go:300`–`:302` 逐字「SetNamedSecurityInfo takes the object by name, so it also validates that the path still resolves to the object we meant」）。**空档窗口本身不漏**：宽的空文件无内容，config 的注释已把这句话写成先例 |

**「甲在最大 DB 下会把多少字节一次性读进内存」——现量答案：这个数没有上限，也没有读数**：

1. **没有配置上限**：全仓唯一与「体积」相关的硬数是 artifacts 目录的 LRU 配额 `ArtifactsQuotaBytes = 500 * 1024 * 1024`（`internal/memory/retention.go:38`），它管 `artifacts\`，**不管 `wisp.db`**；`retention.go:32`–`:36` 只有行级 TTL（`TaskLogTTL`/`ToolCallTTL` 30 天、`CostTTL` 400 天、`GrantAuditTTL` 30 天），是**行数窗口不是字节窗口**。
2. **没有仪器**：包内两个本该用来问体积的常量是死的——`walPageSizeQuery`（`open.go:48`）与 `walCheckpointPassiveOpts`（`open.go:49`）**非测试零引用**（`grep -rn "walPageSizeQuery" --include="*.go" internal/` 只命中定义行；`_test.go` 半边也 0 命中）⇒ 没有任何代码知道 `wisp.db` 有多大。
3. **没有本机读数**：`%APPDATA%\wisp` 与 `%APPDATA%\wisp-dev` 两枚数据根**都是空目录**（`ls -la`：只有 `.`/`..`；`find … -name "*.db*"`＝0 行），仓内也无 `wisp.db*` 任何副本（`find . -name "wisp.db*" -not -path "./.git/*"`＝0 行）⇒ **「实测 DB 体积」这一格今天取不到数**，只能取到结构性上下界。
4. 能给的最硬的数值上/下界，只覆盖**侧车**不覆盖主文件：`-wal` 受 `_pragma=wal_autocheckpoint(1000)`（`open.go:151`，常量 `pragmaWALAutocheckpoint = 1000` 在 `:45`）约束在约 1000 页内，页大小未在本机实测（见 2）⇒ 按 SQLite 默认 4 KiB 页算是**约 4 MB 量级**，按 64 KiB 页算约 64 MB 量级；`wisp.db` 本身**无界**。
5. 与 RSS 预算对表（⛔ 本腿不改阈值，只指认已有尺）：`internal/observe/thresholds.go:19` `memCapSleeping = 25<<20`、`:20` `memCapArmed = 90<<20`、`:24` `memCapWorkPeak = 700<<20 // TARGET (gate=false)`。迁移发生在**进程启动的组装阶段**（§2 链路），正落在 sleeping/armed 这两档的门下；`memVerdict`（`:117`–`:132`）判的是 `MemMedianBytes`（中位数），一枚一次性尖峰**未必**被它抓到 ⇒ 「甲会把峰值顶到门上面」这句话今天**只能说到形状，不能说已证**。

⇒ 本段只摊代价不作选择：**甲＝零新契约面 + 无界峰值内存**；**乙＝有界内存 + 新增导出名（契约级，须回编排者）**；**丙＝有界内存 + 零新增契约面 + 把 winsec 的两条内部保证（失败删档、命名二次校验）搬到 `internal/memory` 自己肩上**。票面「落点只有两支」这句需要按三形改写。

## §4 `PrivateFileExclusive` 那一支会不会撞旁注文件：会，但撞的面比票面写的窄，且它确实是行为变更

**承诺原文逐字**（`internal/memory/open.go:497`–`:500`，本腿现读；票面写成 `:497-499`）：

```
// backupDatabase snapshots wisp.db (plus -wal/-shm when present) into
// backup\wisp.db.bak-<from>-<to>. An existing backup is kept (the oldest
// pre-migration snapshot is the most valuable one; a retry after a failed
// migration must not destroy it).
```

**侧车命名规则逐字**（唯一定义处，`open.go:511`–`:518`）：

```go
for _, suffix := range []string{"", "-wal", "-shm"} {      // :511
        src := s.dbPath + suffix                            // :512
        if suffix != "" {
                if _, err := os.Stat(src); err != nil {
                        continue // no WAL/SHM sidecar: nothing to copy  // :515
                }
        }
        if err := copyFile(src, dst+suffix); err != nil {   // :518
```

基名渲染＝`open.go:494`（`wisp.db.bak-<from>-<to>`）。⇒ 侧车名＝**基名 + `-wal`／`-shm`**，循环里没有任何「已存在就保留」的 guard。

**但票面增量段那句「二次迁移重跑时侧车副本是被覆盖的」少算了一步**：`:506` 的 `os.Stat(dst)` 命中就 `return nil`（`:507`–`:509`），**早退发生在进 `:511` 循环之前**，所以基名在 ⇒ 侧车**也不会被重写**。逐状态摊平（左＝今日，右＝换成 `PrivateFileExclusive`）：

| 状态 | 今日（`os.Create`＝`O_TRUNC`） | 排他创建 | 是否变更 |
|---|---|---|---|
| A 基名在（**票面说的「二次重跑」**） | `:506` 早退，**整套一枚都不写** | 同：早退在写之前，排他根本不会被问到 | **无变更** |
| B 基名不在、侧车在（上一次崩在循环中途且基名被删／从未写成） | `copyFile` 把 `dst-wal`／`dst-shm` **覆盖**，继续迁移 | `PrivateFileExclusive` 返回 `fs.ErrExist`（`winsec.go:90`–`:95` doc 逐字「fails with fs.ErrExist if any entry already occupies the name」）⇒ `:518` 出错 ⇒ `backupDatabase` 返错 ⇒ `migrate()` `:329` 上抛 ⇒ `Open()` `:191` 上抛 ⇒ `run.go:384`–`:387` 打「`wisp run: 存储不可用：…`」并 **rc=2** | **有变更：静默覆盖 → 拒绝启动** |
| C 全新（正常首跑／正常升级） | 写基名（+ 若在位的侧车） | 名字都不在 ⇒ 排他与非排他同果 | **无变更** |

**这三支里侧车到底会不会出现**（现读，⛔ 未跑码）：调用点是 `open.go:325` `boot.Close()` **紧接** `:328` `backupDatabase`——最后一枚 SQLite 连接关闭时做 checkpoint 并删 `-wal`/`-shm`，所以正常路径下 `:514` 的 Stat 落空、侧车那两圈直接 `continue`。既有断言正把这件事钉成了读数：`migrate_v1seed_test.go:110` 与 `schema_test.go:350` 都要求备份目录**恰好**两枚、名字就是 `wisp.db.bak-0-1`／`wisp.db.bak-1-2`（无侧车名）。⇒ **B 态今天没有实测发生过的读数**，它是「代码允许、仪器没见过」的形状。

**判定（照题面要的那句）**：排他创建**会改变已有的用户可见行为**，但只在 B 态，且方向是「把静默覆盖换成拒绝启动」⇒ 这一支**不是纯洁癖**，是一枚低频、但拿可用性换隐私性的行为变更；值不值由编排者裁。若目的只是「先封后写」，甲／丙两支**都不动这张表里任何一格**（`PrivateFile` 的 flags 与 `os.Create` 同为 `O_WRONLY|O_CREATE|O_TRUNC`，`winsec.go:87`）。

**顺带把票面硬约束 2 的两个处置摊平**（不裁）：① 注释改成实话＝只动 `:497`–`:500` 措辞，零行为影响，但那是把「must not destroy」在侧车上的缺口**合法化**；② 把「已存在就保留」扩到侧车＝给 `:511` 循环补一枚 `os.Stat(dst+suffix)` 逐枚早退，行为变更方向＝今天覆盖的变成保留，且 `schema_test.go:350`／`migrate_v1seed_test.go:110` 的「恰好 2 枚」断言**继续为真**（侧车仍不出现）。两支都不许顺手做轮转／上限／删旧（票面已写死：临时件只建不删，删备份＝数据损失）。

## §5 同类站点普查：「封了目录却没锁住随后新建的那枚文件」全仓现量

尺（两把，都现跑）：

```
$ grep -rn 'PrivateDirAll(\|SealDir(' --include='*.go' internal/ cmd/ | grep -v '_test.go'
internal/agent/spill.go:111 / internal/memory/open.go:176 / :188 / :503 / internal/secret/store.go:49
internal/winsec/winsec.go:165:func PrivateDirAll(...) / internal/winsec/winsec.go:222:func SealDir(...)
=> 目录封条生产调用点 5 枚 + 定义 2 枚；SealDir 生产调用者 = 0（票面原读数未变，17 行命中里 16 行在 _test.go）

$ grep -rn 'os\.Create(\|os\.CreateTemp(\|os\.OpenFile(\|os\.WriteFile(' --include='*.go' internal/ cmd/ tools/ | grep -v '_test.go' | wc -l
23
```

**逐枚判「封了目录之后，自己有没有把随后新建的那枚文件也锁上」**（5 枚目录封条全数清点）：

| 目录封条 | 随后在此目录里造文件的站点 | 造出来时封不封 | 判定 |
|---|---|---|---|
| `internal/memory/open.go:176` `PrivateDirAll(abs)`（数据根） | SQLite 自造 `wisp.db`／`-wal`／`-shm` | 不经 winsec（**刻意**：`:173`–`:175` 注释逐字「sealed with inheritance, because wisp.db-wal / wisp.db-shm are created by SQLite and nothing in this repository gets to seal them afterwards」） | **规则内**（由 `:503` 同族机制兜，`internal/winsec/acl_windows_test.go:332` 是这条规则的具名尺） |
| 同上 | `internal/observe/logging.go:72` `os.MkdirAll(cfg.Dir, 0o755)` + `:244` `os.OpenFile(name, O_CREATE\|O_WRONLY\|O_APPEND, 0o644)`（`<data>\logs\*.jsonl`） | **不封**，且 `logs` 目录是**裸 MkdirAll 造的**、`internal/observe/` 全目录**零 winsec 引用** | **同类第 2 枚**（＝票 132 的主题本身，仍在原地未接线） |
| 同上 | `internal/models/downloader.go:424` `os.OpenFile(dst, …, 0o644)`、`:726` `os.Create(dst)`、`archive.go:81` `os.OpenFile(…, 0o644)`（模型树在数据根下：`cmd/wisp/models.go:196` `filepath.Join(dataDir, cfg.Models.Dir)`） | **不封** | **已由裁决豁免**：`internal/models/no_seal_ruling_test.go:12`–`:15` 逐字「Every write site in this package … is NOT passed through winsec. That is a decision, not an omission」，两条前提＝`:17`–`:22`（字节公开：签名 manifest hash-pin；读时逐档 `VerifyDir`），`:27`–`:28`＋`:64` 写明「前提塌了 ⇒ class is void and it must be re-ruled」 |
| 同上 | `internal/ball/position.go:73` `os.MkdirAll(…, 0o755)` + `:77` `os.WriteFile(tmp, b, 0o644)`（球位置） | 不封 | 同类形状，内容非「用户跟 agent 说过什么」那一族；本腿不判它该不该封 |
| `internal/memory/open.go:188` `PrivateDirAll(s.artifactsDir)` | `internal/memory/artifacts.go:75` `winsec.PrivateFileExclusive(...)` | **封** | 正控 |
| `internal/memory/open.go:503` `PrivateDirAll(s.backupDir)` | `internal/memory/open.go:533` `os.Create(dst)`（经 `:518 copyFile`） | **不封** | **＝本票那一枚** |
| `internal/agent/spill.go:111` `PrivateDirAll(s.dir)` | `internal/agent/spill.go:125` `PrivateFile(tmp, []byte(capped), 0o600)`、`:254` `PrivateFileExclusive(path, data)` | **封** | 正控 |
| `internal/secret/store.go:49` `PrivateDirAll(s.dir)` | `internal/secret/store.go:79` `PrivateFile(s.blobPath(id), blob, 0o600)`；另有 `internal/secret/migrate.go:169` `PrivateFile(backupPath…)`、`:174` `SealFile(backupPath)`（对已存在件的修腿）、`:189` `PrivateFile(tmpPath…)`；同族 `internal/config/migrate.go:93` `PrivateFile(backup, raw, 0o600)`、`internal/config/parse.go:201`+`:215` `CreateTemp`→`SealFile` | **封** | 正控（票面第一把尺的正向半边） |

**尺外但同形的另外两枚**（不在票面 `internal/memory internal/secret internal/agent` 那把尺的射程里，本腿按题面「别只报票里那一枚」登记，⛔ 不并账）：

- `internal/observe/diagnostics.go:73` `f, err := os.Create(o.OutPath)` —— 造的是**打包了日志的 zip**（`:66` `MaxLogBytes` 默认 `10 << 20`）。但 `BuildDiagnosticsBundle` 的**生产调用者＝0**（`grep -rn "BuildDiagnosticsBundle" --include="*.go" .` 只命中定义与 doc，两行）⇒ 这一枚的形状是**未接线**，与 §2 那枚相反，别混着写。
- `cmd/wisp/slo_windows.go:946` `os.WriteFile(path, []byte(body), 0o600)`（就绪标记，临时树，`RemoveAll` 清理）与 `:1002` `os.WriteFile(out, data, 0o644)`（SLO 报告，落点由用户给）；`cmd/wisp/doctor.go:301`+`:305` 的 `MkdirAll`+`WriteFile(probe, "probe", 0o644)` 是可写性探针。三枚都**不是私有内容**，登记只为把「裸创建 23 枚」这笔账合上，不建议并入 132 的分母。
- 工具／调试件另 5 枚（`cmd/balldebug/{diff_windows.go:220,:706,shot_windows.go:116,main.go:391}`、`cmd/llmrecord/main.go:137`、`tools/d22scan/selftest.go:406`、`tools/signmodels/main.go:99/:104/:139`）——不在数据根语义内。

**日志那一跳的继承账（本腿顺手量到的、和备份同命的那条）**：`installLogSink` 在 `cmd/wisp/run.go:195` 被调，**早于** `assembleRuntime`（`:207`）里的 `memory.Open`→`PrivateDirAll(数据根)`（票面 `logsink.go:28`–`:31` 自己就写了这个先后）。⇒ 数据根／`logs` 是**先以宽状态被造出来**、后被根的封条收窄。Windows 半边成立（继承 ACE 不落子档、随父重算，`internal/winsec/acl_windows_test.go:386`–`:388` 逐字「an inherited ACE is not stored in the child - it is recomputed from the parent the moment the parent's DACL changes」）；**POSIX 半边不成立**：`sealDir` 在非 Windows 就是一次 chmod（`internal/winsec/winsec_other.go:170` `func sealDir(path string) error { return applyDescriptor(path, true) }`，`applyDescriptorPOSIX` `:24`–`:38` 只 `os.Chmod(path, 0o700)`，**不递归**），而 `privateDirAll` 的下降只覆盖**它自己创建的那些层**（`winsec.go:181`–`:217`，`missing` 表）⇒ 已存在的 `<data>\logs` 会停在 `0755`、里面的 `*.jsonl` 停在 umask 定的 mode。⇒ 「继承兜得住」这句**只在 Windows 侧是既证事实**，跨平台侧它是一枚**待量**。（本票射程＝Windows，D7；这条登记为对照，不建议改射程。）

**补一枚比「没锁文件」更宽的形（本腿顺手核到，⛔ 不在票面六件的字面里，但与 §5 第 2 行同一条腿）**：`installLogSink` 的生产调用者共 **4 枚**（`grep -rn "installLogSink(" --include="*.go" cmd/ | grep -v _test.go`）：`cmd/wisp/run.go:195`、`cmd/wisp/models.go:284`、`cmd/wisp/secret.go:226`、`cmd/wisp/resident_windows.go:57`。其中**常驻腿（resident）在本进程内一枚封条都不下**：`grep -n "memory\.Open\|secret\.NewStore\|PrivateDirAll\|winsec\." cmd/wisp/resident_windows.go` ＝ 只命中两行**注释**（`:48`、`:54`，后者逐字「internal/proc has zero winsec imports」），`internal/proc/` 全目录对 `winsec` 的**导入数＝0**（`grep -rn "CarlosShao/wisp/internal/winsec" --include="*.go" internal/proc/ | grep -v _test.go`＝空），`internal/proc/boot_windows.go` 亦无 `MkdirAll`/`PrivateDirAll`。⇒ 若一台**新机器**先跑常驻腿，`<data>` 与 `<data>\logs` 都由 `observe` 的裸 `MkdirAll(0o755)` 造出、此后**没有任何一层下过封条**（不是"靠继承兜住"，是**没有可继承的封条**）；`models` 腿同理（`models.go:284` 之后只走 `openModelStore`，`internal/models` 零 winsec 导入）。`secret` 腿会经 `store.go:49` 把数据根**在创建时**封掉（`PrivateDirAll` 逐层造＋封），但日志目录是它之前/之后由 `observe` 裸造的，仍只靠继承。⇒ 这条不改 §5 表格的判定，但它把票 132 AC#2 那句「日志目录在**创建时**就 `SealDir`」的**必要性**从"锁不锁文件"抬到了"哪条腿先跑"。

## §6 尺：今天会因「备份文件没加锁」变红的测试＝**零枚**（具名清单＋读数）

**6.1 直接尺：0 枚。** 候选全数点过：

- `internal/memory/` 里带 `winsec`/ACL 字样的测试件只有 2 枚（`grep -rln "icacls\|assertPrivate\|PrivateFile\|SealFile\|winsec\." --include="*_test.go" internal/memory/`）：`artifacts_junction_tripwire_windows_test.go`、`tempdir_resolved_124_test.go` ⇒ **都不判备份件**。
- 静态门也无此规则：`tools/d22scan/main.go` 的规则正则是这 9 条（`:143`–`:151`）加 `:164` 的 `emojiRe`，**没有一条关于 `os.Create`／seal**。
- ⇒ 「把 `:533` 换成不封条的写法」或「保持不封」在今天的门禁里**没有任何一格会红**。

**6.2 间接近尺 3 枚（都判得到别的东西，判不到这一跳）：**

| 具名尺 | 射程 | 为什么判不到本跳 |
|---|---|---|
| `internal/winsec/production_windows_test.go:18` `TestAC3ProductionDataRootIsPrivateEndToEnd`（＋`:65 sweepPrivate`） | 走真实 `memory.Open`，**递归**扫数据根每一枚目录与文件（含 `backup\wisp.db.bak-*`），逐档 icacls 判「无私外来 SID」 | 它判的是**结果 ACL**不是**落盘程序**。备份件靠 `:503` 目录的 PROTECTED DACL **继承**拿到窄集 ⇒ 断言今天即为真。绿的读数（既有腿日志，非本腿跑）：`.scratch/wisp/probes/gate-rerun-1/m5a-red.txt:698` 「sweep checked 8 entries under …\root\data, all private」，紧跟 `:699`「`--- PASS: TestAC3ProductionDataRootIsPrivateEndToEnd (0.29s)`」（同读数第二枚：`m2c-red.txt:782`），同件 `:29` 「data root contains: [artifacts backup wisp.db wisp.db-shm wisp.db-wal]」。⚠ 更要紧的是**它对「这一跳的修法」几乎全盲**：把 `:503` 换成裸 `os.MkdirAll(s.backupDir, 0o700)`（照常造目录、只去掉封条）它**照绿**——`backup` 作为已封数据根（`:176`）的子级仍继承窄集；给备份文件下显式封条或不下，它也**分不出来**（两者结果 ACL 相同）。要它红得连数据根那层一起拆——票 89 的 `PROTECTED_DACL` 变异正是这么红的（`docs/evidence/s1/89-adversarial-acceptance.md:152`–`:160`，读数「11 红」里点名 `TestAC3ProductionDataRootIsPrivateEndToEnd`）。⛔ 注：直接**删掉** `:503` 会红，但红因是「目录没被造出来 ⇒ `os.Create` 失败 ⇒ `memory.Open` 报错 ⇒ `t.Fatalf`」，与 ACL 无关，别把它当这跳的尺。 |
| `internal/memory/migrate_v1seed_test.go:18`（`:103`–`:118`） | 断言 `backup\wisp.db.bak-1-2` 存在 ＋ 目录**恰好** `["wisp.db.bak-0-1","wisp.db.bak-1-2"]` | 对 ACL 无话可说，但**会咬实现形状**：写腿若引入临时件＋rename、或让侧车真的落盘，`len(entries)!=2` 直接红（`:111`/`:116`）。这是本票最可能被误伤的既有尺，具名登记 |
| `internal/memory/schema_test.go:324` `TestMigrateFreshDatabaseCreatesCurrentVersion`（`:346`–`:356`） | 同上（重开不得再写备份） | 同上 |

**6.3 同族「备份件级 ACL 具名尺」在别的包是存在的**（＝证明这种尺写得出来，只是 `internal/memory` 没有）：
`internal/config/private_acl_windows_test.go:57`–`:67`（基线形状：对 `config.toml.bak-1` 跑 `icacls` 判 `BUILTIN\Users` 在不在）与 `:85`–`:99`（接线后 `assertPrivate(t, backup)`）；`internal/winsec/migrate_windows_test.go:72` `TestAC2MigrationBackupIsPrivate`、`:113` `TestAC2PreExistingMigrationBackupIsRepaired`；`internal/agent/spill_acl_windows_test.go`。⇒ 计数：**memory 包 0 枚，config／winsec／agent 各有**。

**6.4 若要立负向判据，正控的形状（⛔ 本腿不写测试，只描述形状）**：

- **先说陷阱，这个陷阱仓里已经有文字钉着**：`internal/winsec/acl_windows_test.go:383`–`:392` 逐字「every child above holds its foreign grant by **inheritance only**, and an inherited ACE is not stored in the child - it is recomputed from the parent the moment the parent's DACL changes. So all of them come out private even if sealDir's propagation walk is deleted」。⇒ 天真地写一句 `assertPrivateACL(备份件)` 就是**测空气**：修复删掉它照绿。本腿在 6.2 第一行给的「单删 `:503` 也照绿」是同一件事的具体化。
- **让「修复前」真的宽，现役两法**：(a) 给该文件下**显式**外来 ACE ——`migrate_windows_test.go:122` `run(t, "icacls", backup, "/grant", "*"+everyoneSID+":(RX)")`，再 `:124 assertCarriesForeign`（`:57`–`:66`，逐字「a criterion that asserts "no foreign principal" on an object that never had one is measuring air, so the wide state is checked first and the test refuses to run otherwise」）；(b) 宽父——`acl_windows_test.go:155`–`:170 wideParent` 给父目录挂 `Everyone:(OI)(CI)(RX)`。
- **但本跳真正可判的差不是「宽不宽」，是「谁给的窄」**：现量对照读数——被 winsec 直接封过的件，icacls 打印**不带 `(I)`** 的 ACE（`m5a-red.txt:29`–`:31`：`…\readable-by-inheritance.txt NT AUTHORITY\SYSTEM:(F)` / `BUILTIN\Administrators:(F)` / `swq:(F)`）；只靠继承的件打印**全部带 `(I)`**（同件 `:23`–`:25`、`:181`–`:182`）。⇒ **「种 X 必响」的正控形状**＝对 `backup\wisp.db.bak-*` 的 icacls 文本要求存在至少一条**非 `(I)`** 的 ACE（可判性现读：够用的帮手是 `acl_windows_test.go:50 icaclsRaw`，它返回含 `(I)`/`(OI)(CI)` 标记的原文；⚠ **`aclSIDs`（`:86`）取的是 `:(` 之前那一段、把 flags 丢掉了，判据不能沿用它**）。X＝把新加的封条注释掉 ⇒ 文件回到 inherited-only ⇒ 立刻红；不种 X 时它今天**本来就红**（现状即 inherited-only），所以这一枚同时是负向尺与"改动生效"探针，**方向自洽**。
- **AC#4 那枚「枚举门」的现成脚手架**：`internal/memory/artifacts_path_invariant_test.go:226`–`:232` 已经写着 `mutators` map（含 `"Create"`、`"CreateTemp"`、`"OpenFile"`、`"WriteFile"`、`"MkdirAll"`、`"Rename"`…）＋ `:259`–`:262` 的「路径必须派生自某个受控根」判定 ＋ `:270` 的计数 `t.Logf`。⚠ 但它只 parse `artifacts.go`（`:164 src := inv76ParseFile(t, "artifacts.go")`），**`open.go` 不在射程**；包级 parse 帮手 `inv76PackageFiles`（`:767`，逐字「parses every non-test .go file in the package directory」）已在位 ⇒ 把 `open.go` 纳入是**扩一根现有尺**，不是新造仪器。这条只登记形状，实现归写腿。

## §7 没查完／留给编排者（本腿射程外或取不到数）

1. **甲形的"多少字节进内存"取不到真实数**：`%APPDATA%\wisp`／`%APPDATA%\wisp-dev` 两枚数据根**都是空目录**（`ls -la`＋`find … -name "*.db*"`＝0 行），仓内也无任何 `wisp.db` 副本。⇒ 能证的只有「无上限」（§3 第 1–3 条），「实际会是几十 MB 还是几百 KB」**不可证**，除非在真机跑一次 `wisp run`（⛔ 本腿禁编译禁跑）。
2. **`-wal` 上界是推算不是实测**：`pragmaWALAutocheckpoint = 1000` 页（`open.go:45`）是现读，但页大小无人问 —— `walPageSizeQuery`（`open.go:48`）与 `walCheckpointPassiveOpts`（`:49`）**全仓零引用（含测试）**。要拿硬数得跑 `PRAGMA page_size`。
3. **POSIX 半边没有尺**：`internal/winsec/` 的 ACL 判据绝大多数带 `//go:build windows`；本腿对「POSIX 侧 mode 由 umask 定」只能从 `os.Create` 语义推，**没有容器内实跑读数**（票面 AC#5 要求 Docker 真跑那一格本腿无权限做）。§5 末段那条「`logs` 目录在 POSIX 上停在 0755」是**代码推论**，同样待实跑。
4. **B 态（基名缺／侧车在）没人造过**：全仓无测试或 fixture 构造过孤儿侧车（`grep -rn "bak-0-1-wal\|\.bak-.*-wal"` 只命中票面与本腿引用的台账／issue 文本，零代码命中）。⇒ §4 表里 B 那一格的风险**是纸面风险**，判它要不要紧需要一枚崩溃注入缝，属写腿／裁决腿。
5. **票面两处行号漂移未改票面**（`:497-499`→实际 `:497`–`:500`；`os.Stat(dst)` `:505`→实际 `:506`）：票面 append-only，改写归编排者／写腿。
6. **票面硬约束 1「落点只有两支」被本腿读出第三支（丙）**：`internal/config/parse.go:201`+`:215` 是**现役生产 idiom**（`CreateTemp` → `SealFile` → 再写），零新增导出名。⛔ 本腿**不主张丙就是正解**，两件事要编排者裁：① 票面那句「代价必须具名：只有甲／乙两支」需按三形改写；② 丙形把 `winsec.go:100`–`:133` 的失败语义（封不上就 `f.Close(); os.Remove(resolved)`，「a refused seal is not a completed write」）**搬到调用方自己肩上**——这算不算「绕开 winsec 写入口」的纪律违例，是治理判断不是普查结论。
7. **票面原表第 4 行与现读冲突，请复核措辞**：该行写「**数据根目录本身：没锁**，出处＝`winsec.SealDir` 在 `internal/winsec/winsec.go:222` 定义、非测试调用者 0」。本腿现读：数据根**已经在 `internal/memory/open.go:176` 用 `winsec.PrivateDirAll(abs, 0o700)` 封过**，且 `PrivateDirAll` 对**已存在**的目标本身也下封条（`winsec.go:155`–`:156` 逐字「An existing path *is* sealed, because that is the repair path for a tree that predates this package」）；`internal/secret/store.go:49` 那枚封的是 `<data>\secrets`，只有当数据根尚不存在、由它逐层造出来时（`winsec.go:144`–`:145`「sealing **every level it creates**」）才顺带封到数据根。`cmd/wisp/doctor.go:237` 的注释也在指认这条链（逐字「the very first sealing call on those systems - secret.NewStore -> winsec.PrivateDirAll」）。⇒ 「`SealDir` 这个名字零调用者」为真（本腿复核＝17 命中、16 在 `_test.go`、1 是定义行），但「数据根没锁」这句若按字面读**为假**；若该行指的是 `<data>\logs\` 那一层或指"名字未被直接调用"，请写明是哪一种，否则 AC#1 那张表会带着一枚错前提被抄下去。
8. **本票原始 AC#1–AC#5 的进度**：票面 5 个勾框全空；`ls docs/evidence/s1/ | grep '^132' | wc -l` ＝ **0**（该目录共 316 份裁决件，无 132 前缀件）⇒ 按「未开工」处理。增量段那一行才是本腿任务，本腿未替 AC#1 那张表现量任何 `icacls` 读数。
9. **`wisp run` 是否真能在本机跑到 `memory.Open`**：`run.go:370`–`:383` 要求文本角色／provider 先就绪，本机数据根为空 ⇒ 无 config.toml／密钥，本腿无法（也不该）验证首跑是否可达；§2 的链路是**代码可达性**结论，不是"owner 机器上已发生过"的读数。
