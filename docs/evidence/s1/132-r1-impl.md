# 票 132 落地腿 `132-r1` —— 日志树与滚动日志文件**在产品里真的被 seal**

- **腿**：`132-r1`＝**写腿（实现者）**。非实现者裁决表另出（`docs/evidence/s1/132-*.md`）。
- **起手锚**：`96f1af9c9`；**本件读数全在本腿自己的树上现跑**，命令逐条写在表里。
- **写面**：`internal/observe/**`（＋ `internal/winsec/**` 只读未改）＋ `docs/evidence/s1/132-r1-*.md`。
  `cmd/**`、`internal/{tools,config,agent,session,memory,risk,panel,ball,proc,audio,watchdog}/**`、
  `scripts/**`、`.github/**`、`go.mod`/`go.sum`、`docs/PLAN.md`／`docs/specs/**` **一字未动**。
- **一句话**：改前 `SealDir` 生产调用者 **0**、真二进制落下的日志文件**只有 `(I)` 继承 ACE**（`Authenticated Users:(I)(M)`、`Users:(I)(RX)`）；
  落地后同一个真二进制、同一条落盘路径 ⇒ 目录与文件都是**显式窄描述符**，`/findsid *S-1-5-11` 与 `*S-1-5-32-545` 双双 `No files with a matching SID was found`，而同账户读得出、记录完好。

---

## §1 改前：三个读数（票面前提逐条现跑，⛔ 未沿用）

| # | 量 | 票面说 | **我现量** | 尺（可逐字重跑） |
|---|---|---|---|---|
| ① | `SealDir` **非测试调用点** | 0 | **0**（非测试命中 1 行＝**定义行** `internal/winsec/winsec.go:301`） | `grep -rn "SealDir(" --include="*.go" . \| grep -v "^\./\.scratch"` ⇒ 17 行；再去 `_test.go` ⇒ 1 行 |
| ② | `SealDir` **测试引用数** | 16 | **16** | 同上，`grep "_test\.go" \| wc -l` |
| ③ | 票面行号 | `winsec.go:222` | **`:301`** | 票面腐坏第 ≥18 例；本件一律按尺不按行号 |

**日志链（写 → 落盘 → 权限），改前逐环**：

| 环 | 位置 | 改前形状 |
|---|---|---|
| 安装 | `cmd/wisp/logsink.go:149` `installLogSink` → `observe.InitLog` | 生产四腿都调（见 §2） |
| 目录 | `internal/observe/logging.go:72`（改前） | `os.MkdirAll(cfg.Dir, 0o755)` —— **不封** |
| 文件 | `internal/observe/logging.go:244`（改前） | `os.OpenFile(name, O_CREATE\|O_WRONLY\|O_APPEND, 0o644)` —— **不封** |
| 封 | 全仓对 `<data>\logs\**` | **零处**（①的读数即证） |

**真二进制读数（改前）**：`wisp-before.exe secret set 132r1probe --from-stdin`（仓外临时根 `data-before2`，1015 字节记录）⇒
`icacls` 文件：`BUILTIN\Administrators:(I)(F)` / `NT AUTHORITY\SYSTEM:(I)(F)` / `NT AUTHORITY\Authenticated Users:(I)(M)` / `BUILTIN\Users:(I)(RX)`；
`icacls <file> /findsid *S-1-5-11` ⇒ **`SID Found`**（`*S-1-5-32-545` 同）。

---

## §2 落地改了什么（逐行）＋ 为什么那是"产品真会走的那一跳"

**只动一个生产文件** `internal/observe/logging.go`（＋ 4 枚测试；`internal/winsec` **一字未改**）：

| 行（改后） | 改前 | 改后 | 为什么 |
|---|---|---|---|
| `:99` | `os.MkdirAll(cfg.Dir, 0o755)` | `winsec.PrivateDirAll(cfg.Dir, 0o700)` | 封条**在层出现之前**就位（父先于子），且既有树走修复分支 |
| `:280` | `os.OpenFile(..., 0o644)` | `os.OpenFile(..., 0o600)` | POSIX 侧创建即不带 group/other |
| `:299` | （无） | `if err := sealLogFile(name); err != nil { close; return Wrap(...) }` | create 之后、**本进程首个内容字节之前**封；失败 ⇒ 关句柄＋返回错，**不静默降级** |
| `:206` | （无） | `var sealLogFile = winsec.SealFile` | 包内 var，只为让"拒绝方向"可测（同 winsec 的 `applyDescriptor`／`removeRefusedLeaf` 形状）；**零新增导出名** |

**为什么这一跳"产品真会走"（⛔ 不是没人会走的一跳）** —— `installLogSink` 的生产调用点**逐枚**：

```
cmd/wisp/run.go:222:            sink, sinkErr := installLogSink(s.dataDir)
cmd/wisp/resident_windows.go:109: sink, sinkErr := installLogSink(rt.Layout.DataDir)
cmd/wisp/models.go:339:         if sink, sinkErr := installLogSink(io_.dataDir); sinkErr != nil {
cmd/wisp/secret.go:226:         if sink, sinkErr := installLogSink(layout.DataDir); sinkErr != nil {
```

再加一条直连：`cmd/wisp/slo_windows.go:272` `observe.InitLog(...)`。⇒ `run`／常驻／`models`／`secret`／`slo`
五条腿，**每一条启动都在 `InitLogWithRegistry` 里穿过 `:99` 与 `:299`**。§3 的读数不是"函数被调了"，
是这五条腿里两条真跑出来的盘上对象。

---

## §3 "产品面真的封了"：逐字 DACL 读数

**第一组＝受控对照**：同一条腿（`wisp slo -settle -seconds 1`）、同一条命令、两枚真二进制、各起一棵**仓外新根**
（`WISP_ENV=test` ＋ `WISP_TEST_DATA_DIR`），根在跑之前**不存在**：

| 对象 | 改前 `wisp-before.exe` | 改后 `wisp-after.exe` |
|---|---|---|
| 数据根 | `BUILTIN\Administrators:(I)(F)` … `NT AUTHORITY\Authenticated Users:(I)(M)` … | `NT AUTHORITY\SYSTEM:(F)`／`BUILTIN\Administrators:(F)`／`DESKTOP-LVS7839\swq:(F)`（全无 `(I)`） |
| `logs\` | 全 `(I)`（同上四条宽 ACE） | 同上六条**显式** ACE（含 `(OI)(CI)(IO)` 三对） |
| `wisp-20261023-001.jsonl` | `BUILTIN\Administrators:(I)(F)` / `NT AUTHORITY\SYSTEM:(I)(F)` / **`NT AUTHORITY\Authenticated Users:(I)(M)`** / **`BUILTIN\Users:(I)(RX)`** | `NT AUTHORITY\SYSTEM:(F)` / `BUILTIN\Administrators:(F)` / `DESKTOP-LVS7839\swq:(F)` |

改后 `logs\` 逐字（创建该目录的那一跳的产物）：
```
D:\tmp\wisp\20261023\132r1\data-fresh-after\logs NT AUTHORITY\SYSTEM:(F)
                                                 NT AUTHORITY\SYSTEM:(OI)(CI)(IO)(F)
                                                 BUILTIN\Administrators:(F)
                                                 BUILTIN\Administrators:(OI)(CI)(IO)(F)
                                                 DESKTOP-LVS7839\swq:(F)
                                                 DESKTOP-LVS7839\swq:(OI)(CI)(IO)(F)
```

**第二组＝文件里有真记录**（`secret set --from-stdin` 落 4 条）：改前 `data-before2` 那枚 1015 字节的文件仍是
`…Authenticated Users:(I)(M)` / `BUILTIN\Users:(I)(RX)`；改后 `data-after2` 那枚 1004 字节的文件只有
`NT AUTHORITY\SYSTEM:(F)` / `BUILTIN\Administrators:(F)` / `DESKTOP-LVS7839\swq:(F)`（首条记录＝
`winsec: sealing path resolver installed`，第末条＝`wisp secret: stored dpapi blob`，都在盘上）。

**反向对照（组 SID 口径，⛔ 不是第二账户登录）**：`icacls <file> /findsid *S-1-5-11` 与 `*S-1-5-32-545`
（`BUILTIN\Users`）：改前双双 **`SID Found: …`** ⇒ 改后双双 **`No files with a matching SID was found`**。
⇒ 本机其余任何**非管理员已认证账户**的访问通道已不在；`Administrators`／`SYSTEM` 仍在（与 `winsec` 既有
`privateACL` 同一张面，本腿未改其语义）。

**AC#1 要的"两种形状"（本腿现量，两形状都测）**：

| 形状 | 数据根 | `logs\` | 文件 |
|---|---|---|---|
| **A 根不存在→本调用创建**（首跑） | **封**（`PrivateDirAll` 连"缺失的每一层"一起封） | 封 | 封 |
| **B 根已存在且宽**（老机器／别人的宽根） | **保持宽**（⛔ 不碰 `path` 以上的祖先，`PrivateDirAll` 的明文契约） | **封**（修复既有树） | **封** |

⇒ 本腿的封条对 `logs` 那一棵**在两种形状下都生效**；数据根只在**它自己由本调用创建**时被顺带封住——
这是 `PrivateDirAll` 的既有语义，不是本腿新增行为，但**必须写清**（`§9.5`）。

**AC#2 两问（同账户零影响／审计照写）**：
- 同账户读：`head -c 400` 读得出，4 条记录完好。
- `wisp slo -settle` 在封过的根上**连跑两发 rc=0**（第二发走 `O_APPEND` 同文件续写，rc 仍 0）；
  `wisp doctor` 在同一根上 `[PASS] data dir writable (test)`。
  （`doctor` 整包 rc=1，FAIL 项全是 DLL 共置／pin，与改前**逐行相同**——见 §6 红名册。）

---

## §4 新用例逐名 ＋ 为什么"今天响"

| 用例 | 读的是什么 | 为什么改前红 |
|---|---|---|
| `TestLogDirIsSealedAtCreation132r1` | 让流水线**自己创建**一个不存在的 `logs` 目录，再读它的 ACL | 改前 `MkdirAll(0o755)` ⇒ 全 `(I)`，含 `Authenticated Users` ⇒ 红 |
| `TestLogFileIsSealedBeforeTheFirstContentByte132r1` | **0 字节时**读文件 ACL，要求**文件自带非 `(I)` ACE**，随后证明同一文件真吃到内容 | 改前全 `(I)`；只封目录不封文件时**仍全 `(I)`** ⇒ 这一枚专治"测空气"（见 §5 M1/M2 分叉） |
| `TestExistingWideLogTreeIsRepaired132r1` | 先把既有目录用 `icacls /grant:r *S-1-5-11:(OI)(CI)(M)` 弄宽（并自证"真宽了"），再 InitLog | 改前不封既有树 ⇒ `these trustees still hold grants: [NT AUTHORITY\Authenticated Users]` |
| `TestSealFailureRefusesTheWrite132r1` | 把 `sealLogFile` 替成必失败，要求 **InitLog 报错且盘上零字节** | 改前无这一跳 ⇒ `reported success with a seal that always fails` |
| `TestUnsealableLogDirIsRefusedLoudly132r1` | 日志目录经 **junction** 到达（C26 拒绝穿 reparse），要求报错且**链接背后零新建** | 改前 `MkdirAll` 照样建、句柄照样开 ⇒ `sealed … through a junction and reported success` |
| `TestSealEntryPointsHaveProductionCallers132r1` | **AST** 枚举 `internal/`+`cmd/` 里对 seal 家族的**调用**，钉住 `internal/observe` 两半都在 | 改前 `internal/observe` 零引用 ⇒ 红；且**按"调用"不按"提及"** |

⚠ 第 6 枚的第一版**按提及**写，M2 量出它**仍绿**（`var sealLogFile = winsec.SealFile` 这行还在提及）⇒ 当场改成跟随别名、只认调用。
这次翻车留在本件里，因为"门自己测空气"正是本票要防的那件事。

---

## §5 反向判据：三读数 ＋ 还原核

**还原核声明**：`sha256sum`（**SHA-256，64 位十六进制**）；⛔ 不与 `git hash-object`（SHA-1，40 位）混比。
pristine `internal/observe/logging.go` = `b0dc8be9b2d6e007d3f8acba50285dd3403795608cda08551f779501af036324`。

| 突变 | 摘掉哪一跳 | `go test -run 132r1` | 红名（逐名） |
|---|---|---|---|
| **M0** | 两半全摘（＝改前形状） | **rc=1，6 红** | 全部 6 枚 |
| **M1** | 只摘目录封（`PrivateDirAll`→`MkdirAll`） | **rc=1，4 红** | `…LogDirIsSealedAtCreation`／`…ExistingWideLogTreeIsRepaired`／`…UnsealableLogDirIsRefusedLoudly`／`…SealEntryPointsHaveProductionCallers`（文件那两枚**仍绿**⇒ 两半互不代偿） |
| **M2** | 只摘文件封（删 `:299` 整块） | **rc=1，3 红** | `…LogFileIsSealedBeforeTheFirstContentByte`／`…SealFailureRefusesTheWrite`／`…SealEntryPointsHaveProductionCallers` |
| **M3** | 留着调用但吞掉错误（`_ = sealLogFile(name)`） | **rc=1，1 红** | `…SealFailureRefusesTheWrite`（＝"没锁但能写"这一支被单独钉住） |

M0 的红行逐字（最能说明"读的是盘上对象"）：
`every ACE is inherited ([NT AUTHORITY\SYSTEM(inherited) BUILTIN\Administrators(inherited) DESKTOP-LVS7839\swq(inherited)]), so nothing was sealed on this object`。
每次突变后 `cp` 回 pristine 并 `sha256sum -c` ⇒ **M0/M1/M2/M3 四次全 `OK`**；随后原样重跑 ⇒ **绿**（§6 的 154 PASS）。

---


> **续件**：§6–§9（门禁 rc／红名册／提交／诚实边界）在 `132-r1-impl-part2.md`——交付帽 160 行，整件 189 行故拆。
