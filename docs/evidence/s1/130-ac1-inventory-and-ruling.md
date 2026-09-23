# 票 130 · AC#1「会掉的记录」可 grep 清单 + AC#2 裁定建议

**锚定 sha**：`ac6f31c`（分支 dev）。**本文件是本票本轮唯一新建/唯一写入的文件**；未修改任何既有文件（票面未勾框、码未动）。
**产出者**：只读审计子代理（2026-09-23 10:2x-10:4x 本机）。**性质**：AC#1 = 清单（可重跑）；AC#2 = **裁定建议，零生产码**。

---

## 0. 一句话结论

- **AC#1**：会掉的记录**今天共 10 枚站点**落在「init 期 / 换 handler 之前的窗口 / 整条腿无听众」三段时间里，另有 **21 枚站点**经「构造期抓走的 `*slog.Logger`」出声（那种式子第一遍 grep 看不见，已补）。三条已被**真二进制实测**证实今天只在 stderr、盘上零命中（§3）。
- **AC#2**：倾向 **(a) 缓冲 + 装好听众再冲刷**，但必须是 **(a-with-mirror)**：缓冲是**额外一份**，stderr 那一遍**永不关**。
- **AC#2 的关键约束（决定正解落点）**：`init()` 一定早于 `main()`，所以**在 `cmd/wisp/` 里加一枚 `init()` 也追不上**。但 `internal/risk` 依赖 `internal/observe`（`go list` 实测：risk 的 imports 含 `wisp/internal/observe` 与 `wisp/internal/winsec`），Go 保证**被依赖包先 init** ⇒ 缓冲装在 `internal/observe` 自己的 `init()` 里就**必然早于** `internal/risk/winsec_c26.go:21` 那次 `SetPathResolver`。⇒ **`internal/risk/winsec_c26.go` 一个字都不必动，解冻授权本轮用不上。**

---

## 1. 枚举方法（照抄可重跑）

三段：分母 → 站点 → 窗口。**每一段都给了式子本身和它的漏法**（漏法在 §4 汇总成一枚待做的钉）。

### 1.1 分母：哪些代码真的在二进制里

```
go list -deps ./cmd/wisp/ | grep CarlosShao          # 22 枚（含 cmd/wisp 自己）
```

闭包**外**的仓内包（同一条式子的补集）：`internal/audio`、`internal/ball`（`go list -deps ./cmd/wisp/ | grep -c internal/ball` → `0`）。它们的记录不是 `wisp` 的记录，但 `cmd/balldebug` 是另一枚 main（§2 E 组）。

### 1.2 站点：谁在出声

```
# 主式子（包级 slog 记录）
grep -rn "slog\.\(Error\|Warn\|Info\)(" --include="*.go" internal/ cmd/ | grep -v _test     # 53 行
# 三条补漏式子（缺一不可，理由见 §4）
grep -rnE "slog\.(Info|Warn|Error|Debug)Context\(" --include="*.go" . | grep -v _test      # 0 命中（今天是空的，式子必须含它）
grep -rn "slog\.Debug(" --include="*.go" internal/ cmd/ | grep -v _test                     # 1 命中（internal/winsec/resolve.go:176）
grep -rnE "\b[a-z]+\.log(ger|\(\))\.(Info|Warn|Error|Debug)\(" --include="*.go" internal/ cmd/ | grep -v _test   # 21 命中（D 组）
grep -rn "\*slog\.Logger" --include="*.go" internal/ cmd/ | grep -v _test                    # 字段/形参清单，用来核上一条有没有漏网
```

53 枚按包分布（闭包内 28 / 闭包外 25）：winsec 7、observe 4、config 4、proc 4、plugin 3、secret 1、statemachine 1、cmd/wisp 4 = **28**；audio 6、ball 19 = **25**。

### 1.3 窗口：三段"没有听众"的时间

```
# (i) init 期：谁在 main 之前就出声
grep -rnE "^[[:space:]]*func[[:space:]]+init[[:space:]]*\([[:space:]]*\)" --include="*.go" . | grep -v _test | grep -v "^./frontend"   # 4 枚
#   3 枚 llm adapter 的 init 只调 llm.RegisterProtocol（internal/llm/provider.go:192，只会 panic，不出声）
#   => 全仓 init 期唯一的记录路径 = internal/risk/winsec_c26.go:21 -> winsec.SetPathResolver
# (ii) 换 handler 之前的窗口：installLogSink 内部 InitLog 已跑、SetDefault 未跑
grep -n "installLogSink\|InitLog(\|SetDefault" cmd/wisp/logsink.go            # 144 InitLog / 148 SetDefault
# (iii) 每条腿的装 listener 行号（一律取 @ac6f31c，工作树里 cmd/wisp 有他人 WIP）
for f in run.go secret.go models.go resident_windows.go slo_windows.go providers.go doctor.go panel_assets.go; do
  git show ac6f31c:cmd/wisp/$f | grep -n "installLogSink"; done                # providers/doctor/panel_assets 命中 0
```

腿的装点在 HEAD 的实际行号：`cmd/wisp/resident_windows.go:57`、`cmd/wisp/run.go:164`、`cmd/wisp/secret.go:219`、`cmd/wisp/models.go:276`（只在 `ensure` 支）、`cmd/wisp/slo_windows.go:250 + :259`（自带 pipeline，已挂 `WISP-LEG-SINK-RULING`）。`providers.go` / `doctor.go` / `panel_assets.go` **一枚都没有**。

---

## 2. AC#1 清单（逐条：文件:行 + 今天去了哪里 + 复现命令）

**"今天去了哪里"三档口径**：`stderr` = 只有 Go 原生默认 logger（`log/slog` 的 Default 是 stderr TextHandler + Info 级）；`盘上` = `<data root>\logs\wisp-<day>-<seq>.jsonl`；`什么都没有` = 级别门或无人调用。

### A 组 · `init()` 期（听众在任何意义上都不存在；先于 `main`）

唯一入口：`internal/risk/winsec_c26.go:21` → `winsec.SetPathResolver`。`SetPathResolver` 的**非测试调用者全仓只有这一枚**（`grep -rn "SetPathResolver(" --include="*.go" . | grep -v _test` → 1 命中），所以每进程恰好命中下面一条、且**只在 stderr**。

| # | 文件:行 | 级 | 消息（原文首句） | 生产可达性 | 今天去了哪里 | 复现命令 |
|---|---|---|---|---|---|---|
| A1 | `internal/winsec/resolve.go:182` | Info | `winsec: sealing path resolver installed`（`resolver=… probes_passed=N`） | 每进程一次（安装成功支） | **stderr**；盘上零命中（**实测**，§3.1） | 最便宜的一发（0.035s，无需 exe）：`go test ./internal/risk/ -run 'TestNothingMatches130Audit' -v` ⇒ **stderr 首行**就是这条（本机实测）。装听众的腿那发：`cd /tmp/wisp130-audit && WISP_ENV=test MSYS_NO_PATHCONV=1 WISP_TEST_DATA_DIR='C:\tmp\wisp130-audit\data' ./wisp.exe secret list`，随后 `grep -rl "resolver installed" data/` → `NONE` |
| A2 | `internal/winsec/resolve.go:158` | Error | `winsec: refusing to install a path resolver into the sealing seam` | 条件：conformance probe 不过（票面容器读数即此条；A 组在 linux 侧同样发声，见 §5.3 末条） | **stderr**；盘上零命中（缺口与级别无关） | 本机 Windows 造不出（probe 全过，实测 `probes_passed=2`）。已在 linux 容器**真跑**过同一入口（`GOOS=linux go test -c ./internal/risk/` + alpine 容器，§3.3）：那一次走的是 **A1 支**（`probes_passed=1`），A2 支没出现 ⇒ **本机侧未验证**（§6）；票面那轮的 A2 读数平台归属未核（§3.3 末句已标明它是推断）。 |
| A3 | `internal/winsec/resolve.go:166` | Error | `winsec: refusing to replace the already installed sealing path resolver` | **生产不可达**：要第二次、不同实例的 `SetPathResolver` | 若被造出则 stderr | `grep -rn "SetPathResolver(" --include="*.go" . \| grep -v _test` 只 1 命中 ⇒ 无第二次 |
| A4 | `internal/winsec/resolve.go:151` | Error | `winsec: refusing to release the sealing path resolver` | **生产不可达**：要 `SetPathResolver(nil)`，非测试调用者 0（票 108 已把这道门移出生产） | 若被造出则 stderr | 同上 + `internal/winsec/resolve_windows_test.go:113` 那句记载 |
| A5 | `internal/winsec/resolve.go:147` | Info | `winsec: sealing path resolver left at the built-in floor` | **生产不可达**（同 A4 的前置） | 若被造出则 stderr | 同上 |
| A6 | `internal/winsec/resolve.go:176` | Debug | `winsec: sealing path resolver installed again identically` | 生产不可达（同 A3） | **什么都没有**：既无听众、又被 `Info` 级门挡（原生默认 = Info） | `grep -n "slog.Debug" internal/winsec/resolve.go` |

> A1/A2 是**互斥**的两支：一次进程里要么装上（A1）要么拒装（A2）。两支都掉。**这一组是本票的主缺口。**

### B 组 · 换 handler 之前的窗口（在 `installLogSink` 内部：`InitLog` 已跑、`slog.SetDefault` 未跑）

| # | 文件:行 | 级 | 消息 | 今天去了哪里 | 复现命令 |
|---|---|---|---|---|---|
| B1 | `internal/observe/logging.go:204` | Warn | `observe: log retention sweep failed` | **stderr**（听众自己的失败是哑的；票 127 已手证） | 触发条件 = 启动 sweep 的 `os.ReadDir(<data>\logs)` 失败，而它前一毫秒刚被 `MkdirAll` 成功 ⇒ 本机未能造出（只读权限/DACL 形状待 AC#3 设计）⇒ **未验证** |
| B2 | `internal/observe/goroutine.go:271` | Warn | `goroutine outside the D38 roster (leak symptom)` | **stderr**（若触发）：`InitLog` 在 `logging.go:96` 就 `reg.Spawn("log-flusher", …)`，而 `SetDefault` 在 `cmd/wisp/logsink.go:148` | 今天 `log-flusher` 在名册内 ⇒ 实测未触发（§3.1 的 stderr 里没有这条）。要红它只需把名册条目摘掉（AC#3 的变异素材） |

**这一组的形状**：`installLogSink` 的 4 行里，第 1 行已经会出声，第 2 行才装听众 ⇒ **B 组永远是空的才该是终态**（§4 的钉 #2）。

### C 组 · 各腿"装听众之前"自己就能出声的记录

| # | 文件:行 | 级 | 消息 | 哪条腿在装之前 | 今天去了哪里 | 复现命令 |
|---|---|---|---|---|---|---|
| C1 | `internal/winsec/winsec_windows.go:143`（`noticeNarrowed`） | Warn | `winsec: seal cleared principals that stood on this object` | **`models ensure` 腿**：`cmd/wisp/models.go:260 loadModelConfig` / `:265 openModelStore` 都在 `:276 installLogSink` **之前**，而 `config/parse.go:215`（`atomicWrite`→`SealFile`）与 `config/migrate.go:93`（`PrivateFile`）都是封印点 | **stderr**（该腿那一次） | `./wisp.exe models ensure <id>`（需配置）；今天无凭据 ⇒ 见 §6 未验证。等价形状已实测：`./wisp.exe providers bogus` 的 stderr 第 2 行就是这条（§3.2） |
| C1' | 同上一枚站点 | Warn | 同上 | `secret` / `run` / `resident` / `slo` 四腿：装点在**任何封印之前** ⇒ **同一条记录那四腿落盘** | **盘上**（对照项，不是缺口） | §3.1 第二条 JSONL 记录就是它（时间戳晚于 install 那条 6ms） |

其余 24 枚闭包内站点（`internal/config/manager.go:244/247/254/260`、`internal/proc/shutdown.go:130/145/147/172`、`internal/proc/boot_windows.go:127`、`internal/observe/goroutine.go:168`、`internal/observe/logging.go:307`、`internal/plugin/disposal.go:208/332/334`、`internal/statemachine/machine.go:139`、`internal/config` 其余、`cmd/wisp/secret.go:385/530/532`、`cmd/wisp/logsink.go:165`）经核对**今天都在装之后**：`config.NewManager` 的唯一生产调用点 `cmd/wisp/run.go:252`、`memory`/`plugin` 的装配点均在 `installLogSink` 之后（run 腿 `:164` 早于 `:287 memory.Open`）。`observe/logging.go:307` 是 500ms tick 之后 ⇒ 也在装之后。

### D 组 · 经"构造期抓走的 logger"出声的一形（第一遍 grep 看不见；21 枚站点）

`internal/memory/open.go:118` 与 `internal/agent/loop.go:208` 都把 **`slog.Default()` 的当时值存进字段**，之后每次记录都用那枚被抓住的 logger。⇒ 这条缺口**与"记录发生在哪一刻"无关，只与"对象在哪一刻被构造"有关**：谁把 `memory.Open` 挪到 `installLogSink` 之前，那 14 枚记录就**永久**走 stderr，听众已经在也一样。

| # | 站点 | 级 | 今天去了哪里 | 复现命令 |
|---|---|---|---|---|
| D1 | `internal/memory/artifacts.go:102/271/321`、`internal/memory/open.go:236/240/343/507/522`、`internal/memory/retention.go:109/120/141/167/174/180`（14 枚，Info/Warn/Error 混） | Info/Warn/Error | **`providers` 腿：stderr**（`cmd/wisp/providers.go:172 memory.Open`，而该腿全生无听众）；`run` 腿：盘上（构造在 `run.go:164` 之后） | `./wisp.exe providers probe <p>/<m>`（需凭据，未真跑 ⇒ §6）；静态式子：`git show ac6f31c:cmd/wisp/providers.go \| grep -n installLogSink` → 0 命中 |
| D2 | `internal/agent/loop.go:373/398/588/660/701/741/933`（7 枚，全 Warn） | Warn | 盘上（Loop 在 run 腿装之后构造）；** latent 缺口**，同上形状 | 同上（构造点 `internal/agent/loop.go:208`） |

### E 组 · 整条腿没有听众（不是"早"，是"never"）

| 腿 | 状态 | 会掉的站点 | 今天去了哪里 |
|---|---|---|---|
| `wisp providers`（HEAD `cmd/wisp/providers.go:72`） | **无 `installLogSink`、无 `WISP-LEG-SINK-RULING`** | C1（`:84 secret.NewStore` 触发的封印通知）+ D1 的 14 枚 + `plugin/disposal.go:208/332/334` | **stderr**（实测，§3.2） |
| `wisp models list` / `verify`（`:201/:230/:235` 三支） | `installLogSink` 只在 `modelsEnsure`（`:276`） | 今天该三支自身无记录站点（`internal/models` 零 `slog.*`） | 若将来加一条即掉 |
| `wisp doctor` / `panel-assets` / `version` / `help` | 无听众 | 0 枚（grep 证明） | 无义务 |
| `cmd/balldebug`（另一枚 main，`cmd/balldebug/main.go:118`） | 只把默认 logger 装成 **stderr TextHandler**，无 JSONL | `internal/ball` 19 枚 + `internal/config` 4 枚 + A1 | **stderr**（dev 工具，本票范围外但同形） |
| `internal/audio`(6) / `internal/ball`(19) 对 `wisp` | **不在 `cmd/wisp` 闭包**（票 07 未落地） | — | 不发声（不是"掉"，是"没接线"） |

### F 组 · 枚举顺带抓到的一枚"能力完成、无人调用"

| # | 文件:行 | 级 | 今天去了哪里 | 证据 |
|---|---|---|---|---|
| F1 | `internal/secret/migrate.go:198`（`MigratePlaintext` 的 D33 迁移通知） | Warn | **什么都没有**：`grep -rn "MigratePlaintext(" --include="*.go" . \| grep -v _test` → 只有它自己的定义 ⇒ 生产零调用者 | `cmd/wisp/logsink.go:8-9` 的建文件理由里点名 "the D33 … credential migration"（字面跨两行，`grep -F 'D33 credential'` 打不到，核它要按行号读）是"已经在写给一个 logger 的记录"之一；它的 `internal/config/migrate.go:85` 注释也点名这枚函数（`secret.MigratePlaintext exists precisely because`）。**结论：那半句承诺今天没有任何生产路径能兑现**（属 `A105④` 同族，不属本票顺序缺口） |

**条数合计**：A 6 + B 2 + C 1（同一站点在 models 腿；对照项 C1' 不计缺口）+ D 21 + E 组按腿计（providers 3 类站点、balldebug 23 站点、无听众但零站点 3 腿）+ F 1 = **本票范围内 10 枚站点（A+B+C+F）+ 21 枚 D 组站点 = 31 枚**；E 组的 balldebug/audio/ball 属另一枚 main 与未接线的包，登记但不计入 `wisp` 的缺口。

---

## 3. 实测读数（真二进制 + 一枚 linux 容器；一律仓外临时根。§3.1/§3.2 是 Windows 原生，不需要 Docker）

**仪器自证**：`go build -o /tmp/wisp130-audit/wisp.exe ./cmd/wisp/` → rc=0；三枚 DLL 从 `build/` **复制**到该临时根（仓内 `build/` 未动）；数据根一律 `WISP_ENV=test` + `WISP_TEST_DATA_DIR=C:\tmp\wisp130-audit\…`（仓外、每次 `rm -rf` 重建）。
**版本口径**：exe 构建自**脏工作树**；已 `git diff ac6f31c -- cmd/wisp/` 逐行核对——他人 WIP 改的是 `resolveDataDir` 的失败措辞与数据根解析（`providers.go`/`secret.go`/`models.go`/`run.go`/`doctor.go` 共 7+11+10+15 行），**`installLogSink` 与 `InitLog`/`SetDefault` 的相对顺序一行未动**（HEAD 的装点在 `run.go:164`/`secret.go:219`/`models.go:276`/`resident_windows.go:57`，工作树是 `:177`/`:226`/`:284`/`:57`，差的是上面那些行的位移）。`internal/` 全目录 `git status --short internal/` → 干净 ⇒ §2 的 `internal/**` 行号就是 HEAD 的行号。

### 3.1 `wisp secret list`（装听众的腿）

```
$ WISP_ENV=test WISP_TEST_DATA_DIR='C:\tmp\wisp130-audit\data' ./wisp.exe secret list
2026-09-23 10:24:50 INFO winsec: sealing path resolver installed resolver=risk.c26Pipeline probes_passed=2
time=2026-09-23T10:24:50.860+08:00 level=INFO msg="wisp: persistent log sink installed" dir=C:\tmp\wisp130-audit\data\logs min_level=info
time=2026-09-23T10:24:50.866+08:00 level=WARN msg="winsec: seal cleared principals that stood on this object" path=…\data\secrets kind=inherited …
```

盘上 `C:\tmp\wisp130-audit\data\logs\wisp-20260923-001.jsonl` 全文 **2 行**（install + seal notice，`time` 是带 `+08:00` 的 ISO），`grep -rl "resolver installed" data/` → **NONE**。stderr 首行的时间格式是 `2026-09-23 10:24:50`（Go 原生默认 logger），与盘上两行**不同形** ⇒ 与票 127 手跑的读数同形。
本机是 `probes_passed=2`（A1 支），票面容器是 A2 支（`probes_passed=1` 那次是被换打的探针）⇒ **两支都会掉，与级别无关**这一条在本机侧得到独立复现（A1 是 INFO 级、成功支，照样零命中）。

### 3.2 `wisp providers bogus`（无听众的腿）

```
$ WISP_ENV=test WISP_TEST_DATA_DIR='C:\tmp\wisp130-audit\datap' ./wisp.exe providers bogus
2026-09-23 10:29:38 INFO winsec: sealing path resolver installed resolver=risk.c26Pipeline probes_passed=2
2026-09-23 10:29:38 WARN winsec: seal cleared principals that stood on this object path=C:\tmp\wisp130-audit\datap\secrets kind=inherited …
wisp providers: 配置未就绪（Unconfigured）：… rc=2
```

`find datap -type f` → **0 个文件**（连 `logs\` 都没有）⇒ E 组那一档不是"早了一步"，是"从来没有"。同一条 `winsec_windows.go:143` 记录在 3.1 里落盘、在 3.2 里不落盘 —— **这就是本票缺口最干净的一发对照**，也是 AC#3 判据应当长的样子。

### 3.3 容器真跑（POSIX 半边，A 组在 linux 上确实发声；挂载已自证）

```
$ GOOS=linux go test -c -o /tmp/wisp130-audit/risk-linux.test ./internal/risk/   # rc=0
$ MSYS_NO_PATHCONV=1 docker run --rm -v /d/tmp/wisp130-posix:/work alpine:3.19 sh -c \
    'ls -l /work/risk-linux.test && /work/risk-linux.test -test.v -test.run TestNothingMatches130Audit'
-rwxrwxrwx 1 root root 6879195 Sep 23 02:52 /work/risk-linux.test          <- 挂载自证（字节数与宿主一致，非空挂载）
2026-09-23 02:52:34 INFO winsec: sealing path resolver installed resolver=risk.c26Pipeline probes_passed=1
testing: warning: no tests to run / PASS
```

**这条读数比票面的更硬**：票 130 说"换打 `INFO … resolver installed probes_passed=1` 同样零命中"——那是**改过码之后**才有的句子；本机不改一行码，在 linux 容器里**天然**打出的就是同一句（`probes_passed=1`，因为 `resolverProbeShapes()` 在非 windows 只有 1 枚形状；windows 是 2）。⇒ "缺口与级别无关、是顺序"这一条现在有**两支平台读数 + 零伪造**。
A2（`refusing to install`，ERROR 支）在这轮容器里**没有**出现（探针全过）⇒ 它需要一台让 conformance probe 失败的主机，本机与本机容器都造不出 ⇒ 仍列 §6 未验证；票面那轮的 stderr 首行是 A2，说明那轮跑的东西与这轮不同（多半是一枚自带 fake resolver 的 `internal/winsec` 测试二进制），**这一句是推断，不是读数**。

### 3.4 owner 真实数据目录（开工前 / 开工后各读一次）

- 前：`%APPDATA%\wisp` 与 `%APPDATA%\wisp-dev` 均存在、`find -type f | wc -l` = **0 / 0**（2026-09-19 14:49 与 09-20 07:06 的目录 mtime）。
- 后：见 §7（同一条命令再读，要求仍是 0/0 且 mtime 未变）。

---

## 4. 枚举式本身会不会静默漏掉新增？会。已核出的三种漏法 + 要做成的两枚钉

**直接回答**：会。§2 是**人肉枚举**，新增一条 `init()` 期的 `slog.Warn`（例如在 `SetPathResolver` 里加一条新分支的记录、或在 `internal/risk/winsec_c26.go:20` 那枚 `init()` 里直接打一行）⇒ 清单不会自己变红，§3 的实测也照样读不出来（那条新记录可能根本不触发）。**所以本票要的不是更多清单，是下面这枚钉**（按边界，我只登记形状，不实现）。

已核出并已补的三种漏法（每种都给了式子，见 §1.2）：
1. **`slog.Debug`**：主式子不含 Debug ⇒ 漏 `internal/winsec/resolve.go:176`。补法：把 `Debug` 与 **`InfoContext/WarnContext/ErrorContext/DebugContext`** 一起进式子（后者今天 0 命中，但今天 0 命中不等于式子能漏 —— 这正是 `cmd/wisp/leg_sink_gate_131_test.go` 自己 `recordEmittingSelectors` 里的四个名字，它连 `*Context` 也没含）。
2. **限定式之外还有"捕获式"**：`slog.Default()` 在构造期被抓进字段（`internal/memory/open.go:118`、`internal/agent/loop.go:208`），之后全部经 `s.logger.X` / `l.log().X` 出声 ⇒ 21 枚站点对 `grep 'slog\.' 完全隐形。**现存门禁同样看不见**：`leg_sink_gate_131_test.go` 的注释明写"Qualified form only (slog.Info, not x.Info)" ⇒ 它把 `providers` 腿判成"既不装也不出声、无义务"那一档，而实测该腿 `memory.Open`+`NewStore` 是会出声的（§3.2）。这是**现有门禁的一枚真漏**，本票必须点名。
3. **闭包**：不在 `go list -deps ./cmd/wisp/` 里的包（`internal/audio`、`internal/ball`）不是 `wisp` 的记录；反之，只按目录 grep `internal/` 会把它们误计进缺口。⇒ 站点集合必须与闭包求交，式子在 §1.1。

**钉 #1（本票 AC#3/后续动码者实现，一枚就够；写进报告不实现）**
- 断言（正向、覆盖新增）：对 `go list -deps ./cmd/wisp/` 的每枚包（非测试文件），计算「从每个 `func init()` 可达的函数闭包」内的**记录站点集合 S**（口径 = 上面 1+2 两条合起来：限定式 `slog.{Info,Warn,Error,Debug}[*Context]` ∪ 经 `*slog.Logger` 字段/返回值的方法调用），要求 **S 里每一枚站点所在的包，其传递依赖里必须含"缓冲的宿主包"**（裁定里是 `internal/observe`）。违反即红，红在**具名文件:行**。
- 为什么这形比"静态白名单"强：白名单要人改清单（新增即漏、还永远绿）；这一枚问的是**结构不变式**（"任何 init 期出声者必须在缓冲之后"），新增一条 init 期记录要么它所在包已依赖 observe（自动被缓冲接住，绿且**正确**），要么没依赖（红，且红得具体）。同形先例就在本仓：`cmd/wisp/leg_sink_gate_131_test.go`（go/ast + 名字传递闭包 + "分母不足即红"的非恒绿自证，见其 `len(pkg.funcs) < 40` 那类守卫）与 `cmd/wisp/dataroot_128_test.go:308 resolveDataDirConsumers128`。
- 非恒绿自证（AC#3 的那一发变异）：给 `internal/risk/winsec_c26.go` 的 `init()` 里加一行 `slog.Warn("probe")` ⇒ 钉 #1 必须红（`risk` 依赖 `observe`，所以**这条会绿**？）—— 不，红点在于**新增记录点必须被登记**，所以钉 #1 要同时含一枚"清单一致性"分量：S 的**大小与成员**必须等于签入的那份清单（今天 = A 组 6 枚），改结构不等于免登记。⇒ 钉 #1 = 两分量：**（i）成员白名单**（新增即红）+ **（ii）依赖方向**（缓冲宿主必须在 emit 之前 init）。只写 (ii) 会放过"缓冲装不上的那条腿"，只写 (i) 会退化成清单。
- 钉 #2（便宜、今天就该先记账）：`cmd/wisp/logsink.go` 的 `observe.InitLog(...)` 与 `slog.SetDefault(...)` 之间**不得存在任何记录调用**（B 组永远为空）。今天不为空（`InitLog` 内部会走到 `logging.go:204`、`goroutine.go:271`），所以这枚钉应当**今天红**，等 (a) 落地（缓冲在 `observe` 的 `init()` 里就装好）才转绿 —— 这正是"判据先于修法"的形状。

## 5. AC#2 裁定建议（只要裁定；本文件不含任何生产码改动）

**裁定：(a) 内存缓冲 + 装好听众再冲刷，但取 (a-with-mirror) 这一形：缓冲是"额外一份"，stderr 那一份从 `init()` 起就永不摘。**

一句理由：本票要保的是**安全裁决的可审计性**，而 (a-with-mirror) 的失败模式严格包含今天的行为（最坏 = 和今天一样只有 stderr），(b) 的失败模式是"封印 seam 被拒/被装"这条记录**永久**不可审计 —— 而 owner 真正在用的 resident 腿**没有终端**（`cmd/wisp/logsink.go:11-14` 自己写了这一点）。

### 5.1 ① 正解落在哪一枚文件（要不要申请动 `winsec_c26.go`）

**落点两枚，都不在冻结清单、都不需要授权**：

1. `internal/observe/logging.go`：新增缓冲 handler + 一枚 `func init()`，把自己装成默认 logger（tee = 缓冲 + stderr text）。
2. `cmd/wisp/logsink.go`：在 `slog.SetDefault(teeHandler{…})`（现 :148）之后、`slog.Info("wisp: persistent log sink installed")`（现 :165）之前，把缓冲按时间序冲进 pipeline；并在无听众腿的退出路径上 drain 一次。

**明说：不要动 `internal/risk/winsec_c26.go`。owner 那枚具名解冻本轮可以不用。**

为什么不需要它：Go 的包初始化顺序是「被依赖包先 init」，而 `internal/risk` 的 imports 实测含 `wisp/internal/observe` 与 `wisp/internal/winsec`（`go list -f '{{join .Imports "\n"}}' ./internal/risk/`）⇒ **`observe` 的 `init()` 必然早于 `risk` 的 `init()`**，也就是早于 `internal/risk/winsec_c26.go:21` 那次 `SetPathResolver`。那枚 `init()` 只是**发声方**，不是**装配方**；把它改成"延迟安装/自己攒着"反而会把 C26 的封印装配与日志缓冲搅在同一枚文件里（该文件 :5-19 的注释正在防"第二套归一化管道"这件事）。
反过来一条硬事实要写清楚，因为它否掉了最直觉的修法：**`cmd/wisp/` 里加 `init()` 追不上**（main 包的 init 最后跑），所以正解只能是"更早的依赖包" —— `observe` 恰好是 `risk` 的依赖，且**不必新增包**、**不必给 `internal/winsec` 加 import**（`internal/winsec` 目前零仓内 import、零 `init()`，实测）。

### 5.2 ② 缓冲的上界、溢出策略与"带进坟"的问题

- **上界**：条数 64 + 字节 64 KiB 双限（与 `internal/observe/logging.go:34` 的 `bufferedLogBytes = 64 << 10` 同形）。量级依据是实测：本机一次 `wisp secret list` 在装听众之前只出 **1 条**（§3.1），A 组每进程最多 1 条 + B 组条件 2 条 + C 组那一次。64 条对"将来多枚包 init 都出声"仍够，且 ≤64 KiB 对 boot 态 RSS 的可见增量应为 0 —— 但**不许口头说很小**，AC#4 要用 `wisp slo` 的 boot 态复核一次（D32 阈值、`A105⑤` 的"冷读单独成格"）。
- **溢出策略**：满后**停止追加 + 计数**（不滚动丢旧、也不丢新）。冲刷时写两条：缓冲里的 N 条 + 一条 `observe: early log buffer overflow` WARN，属性带 `dropped=N`、`first_kept=<msg>`、`capacity=64`。理由：启动顺序里**最早**那条往往是原因、最新那条只是结果；而"丢了 N 条"必须成为盘上的一条事实，不能是一个静默的减法。
- **"启动期被卡住/崩溃会不会把缓冲带进坟"**：不会比今天更糟 —— 因为 **mirror 那一遍 stderr 从不关**（这就是我不选"纯 (a)"的原因）。两条路都核过：(i) 装听众之前 panic / 被 Job Object 杀 / 卡死在 init ⇒ 缓冲随内存消失，但 stderr 里那几条本来就在（= 今天的状态）；(ii) 一条**永不装听众**的腿（`providers` / `doctor` / `models list`+`verify` / `balldebug`）⇒ 由 `main` 退出路径上的 drain 收口，且 drain 幂等（stderr 已镜像过，drain 只为"将来那半已落盘"兜底，不重复打）。⇒ **(a-with-mirror) 的行为集合 ⊇ 今天的 + 盘上多一条**，这是"绝不倒退"的那道保险。
- 不推荐的更硬版本：早期记录 append 时就 best-effort 直接写盘 —— 它要求在 `init()` 里解析数据根（启动期新增文件系统活动，会先撞上 D32 冷读预算）。登记以便 owner 否决时有据。
- 若 owner 仍选 **(b)**：文案**不能**只写在 `docs/specs/**` 或 `docs/PLAN.md`（两者禁改）。钉得住的形状只有：`cmd/wisp/logsink.go` 顶部一段 + 新标记 `WISP-INIT-NO-LISTENER:`（与 `WISP-LEG-SINK-RULING` 同族，先例 `cmd/wisp/slo_windows.go:171`），门禁 = §4 钉 #1。**代价**：这枚门禁只钉"注释与代码不脱节"，钉不住"记录消失"本身 —— 这一支我登记但不推荐。

### 5.3 ③ 判据必须能红（先想清楚这一发变异长什么样，才交裁定）

**判据（今天就是红的，不需要伪造任何记录）** —— 新用例形状，复用现成仪器（`cmd/wisp/resident_sink_nail_127_windows_test.go` 的 `bootResidentLeg` / `readResidentSink` / `pollUntil127` / `jsonlFilesUnder`；`cmd/wisp/logsink_test.go` 的 `readSink`）：

1. 驱动 shipped exe：`WISP_ENV=test` + `WISP_TEST_DATA_DIR=<仓外临时根>`，命令取**装听众那条腿**（`wisp secret list`）。A1（`internal/winsec/resolve.go:182`）必然在 init 期出声。
2. 读 `<临时根>\logs\wisp-*.jsonl`（按 `observe.CountLogFiles` 认文件，改名即红，不同 127 那枚读法）。
3. 断言**存在**一条 `msg == "winsec: sealing path resolver installed"`（字面量抄进用例，同 127 的 `residentInstallMsg` 口径），且 `resolver` / `probes_passed` 两个属性在。
4. **顺序断言**（防"装完之后重打一遍"糊过去）：它的下标 **<** install 那条的下标，且它的 `time` 早于 install 那条的 `time`。
5. **反向断言**（防 (b) 用"关掉 stderr"蒙过 (a)）：子进程 stderr 里 `"sealing path resolver installed"` 仍**恰好 1 次**。
6. 非恒绿：§3.1 已实测盘上只有 2 行、`grep -rl "resolver installed" data/` → `NONE` ⇒ **断言 3 今天必红**；绿只能靠真把那条搬进文件，断言 4 再堵"搬进来但不是那一条/那个时刻"。

**那发变异（AC#3 跑，本文件只登记形状）**：
- **M-1 拆冲刷**：删 `cmd/wisp/logsink.go` 里那行冲刷 ⇒ 判据红在断言 3，红名可读是哪条腿。
- **M-2 拆装配位**：把 `observe` 的 `init()` 换成在 `main` 里装缓冲 ⇒ 红在断言 3（证明"必须早于 init"）；把 mirror 那半摘掉 ⇒ 红在断言 5（证明"stderr 那份不许关"）。两发各自钉住一条承诺，而不是靠注释。
- **M-3 守 D 组那一形**：加一条子判据 —— `wisp providers` 腿的盘上应出现 `retention job scheduled`（`internal/memory/retention.go:109`，经构造期捕获的 logger）；今天该腿零文件（§3.2）必红 ⇒ 这样 `leg_sink_gate_131_test.go` 看不见的那 21 枚站点也有人了。
- **POSIX/容器半边（已真跑，口径与上一版不同）**：`internal/winsec/resolve.go` **无 build tag**（`head -1` → `package winsec`），`GOOS=linux go list -deps ./cmd/wisp/` 里 `internal/risk` 命中 1 ⇒ **A 组在 linux 侧同样发声**，且 §3.3 已在容器里真跑到（`probes_passed=1` 那一支 = A1）。⇒ AC#3 的分母两半都有：POSIX 半边 = 容器真跑 `go test -v ./internal/risk/ ./internal/winsec/`（挂载先 `ls -l /src/go.mod` 或 `ls -l /work/<file>` 自证、加 `MSYS_NO_PATHCONV=1`）；但 **shipped exe 的 Linux 腿没有**（`GOOS=linux go build ./cmd/wisp/` 在本机失败：`sherpa-onnx-go-linux` build constraints exclude all Go files，与 `cmd/wisp/logsink.go:49-53` 的自述一致）⇒ 别把容器里"测试二进制自带的 pipeline"说成"生产听众"。要说某条 CI 门的存废，得先说出 run id + step（`A103②` 口径）。

### 5.4 代价（选 (a) 要一并付的，如实列）

1. **会撞一枚既有判据（必须一起处理，不能事后偷偷改）**：`cmd/wisp/resident_sink_nail_127_windows_test.go` 现在钉着 `recs[0].Msg == "wisp: persistent log sink installed"`。早期记录一冲进同一文件，record 0 就变成 A1 ⇒ **那枚钉红**。新承诺应重排成"早期记录按时间序在最前、install 记录紧随其后、install 之后再无早期记录"，改法属动码轮，且**改别人的判据要单独报备**。
2. `wisp slo` 腿自带 pipeline（`cmd/wisp/slo_windows.go:250/:259`，带 `WISP-LEG-SINK-RULING`）⇒ 它要么也冲一次，要么显式声明不冲（否则两个 writer 争同一棵树，正是那条 ruling 的理由）。
3. 无听众腿（E 组）不因本修得救：缓冲在那些腿上永不冲刷 ⇒ 必须与"给 `providers` 腿装听众"同批做，否则 (a) 只治了 init 期、把 providers 腿留在原样。**这是票 131 的枚举域，不并到本票**。
4. AC#4 的面变大：受影响包 `internal/observe`、`internal/risk`、`internal/winsec`、`cmd/wisp` 各 `-count=2 -v` 四数 + gofmt/gofumpt 全路径 + `go vet` 双 GOOS + d22scan 纯净快照。
5. 时序：`cmd/wisp/**` 此刻有另一名代理的 WIP（票 128/131 那批），本修的冲刷点在 `logsink.go` ⇒ **动码要排在它交件之后**（文件级冲突判据）。


## 6. 未验证项（明写，不当作下一次结论的地基）

1. **A2（`internal/winsec/resolve.go:158`）没有任何本机读数**：本机 Windows 与本机 linux 容器都让 conformance probe 全过（实测 `probes_passed=2` / `probes_passed=1`）。票面那轮读数的**触发条件与平台归属都没核**（§3.3 末句是推断，不是读数）。⇒ AC#3 若要覆盖 A2 支，只能"造一台让探针不过的主机"或"以测试内 fake resolver 复现同一支"，两者都不在本票只读范围。
2. **B1（`internal/observe/logging.go:204`）的触发条件没造出来**：启动 sweep 的 `os.ReadDir(<data>\logs)` 在 `MkdirAll` 刚成功一毫秒之后要失败——本机 Windows 上能想到的形状（权限、DACL、符号链接）都没实测成。⇒ 这一条是**静态枚举 + 票 127 的手跑**，本轮没有新读数。
3. **B2（`internal/observe/goroutine.go:271`）今天不触发**：`log-flusher` 在 D38 名册内；实测 stderr 里没有这条。它是"窗口形状"的证据（在 `SetDefault` 之前被调用的记录点），不是当下的丢失量。
4. **C1 在 `models ensure` 腿上的那一发没真跑**：`./wisp.exe models ensure <id>` 需要配置与模型 id，本机无凭据 ⇒ 只有静态顺序证据（HEAD 的 `cmd/wisp/models.go:260/:265` 早于 `:276`）+ 等价形状在 `providers` 腿上的实测（§3.2）。
5. **D1（memory 的 14 枚 + plugin 的 3 枚）在 `providers` 腿上的落盘读数没有**：`wisp providers probe <p>/<m>` 需凭据。已实测的是同一条腿上的 `NewStore` 那半（§3.2 的 `winsec` WARN 落 stderr、`find datap -type f` → 0 文件）。⇒ "providers 腿会掉 D1"目前= 静态推断。
6. **四数没跑**（AC#4 属动码轮）：本轮只跑了两枚定向编译（`go build ./cmd/wisp/`、`GOOS=linux go test -c ./internal/risk/`）与一枚 0 用例的定向 `go test ./internal/risk/ -run 'TestNothingMatches130Audit' -v`（0.035s，`ok … [no tests to run]`）。**没有跑全仓 `go test ./...`**（同树有他人 WIP）。
7. **§2 里"其余 24 枚闭包内站点今天都在装之后"是调用点核对，不是运行时读数**：核对依据是"生产调用点唯一 + 该行在 `installLogSink` 之后"。若动码轮改动装配顺序，这句要重走。
8. **未跑 mutation**：§5.3 的 M-1/M-2/M-3 与 §4 的两枚钉都只登记形状（本票禁动码）。

---

## 7. 边界自证 + 待办

**只新建了一枚文件**：`docs/evidence/s1/130-ac1-inventory-and-ruling.md`。既有文件零改动（`internal/`、`cmd/`、票面、`docs/PLAN.md`、`docs/specs/**`、`rules_gateway.go`、`tools/d22scan/**`、`allowlist.txt`、`frontend/**` 全部只读）。
- `git add` 只用了显式路径；`git diff --cached --name-only` 的输出写在下面（要求：恰好 1 行）。
- 未 add 的他人 WIP（保持原样）：`cmd/wisp/{doctor,models,providers,run,secret}.go`、两枚 `cmd/wisp/dataroot_128*_test.go`（新，属票 128）、`docs/evidence/s1/131-adversarial-acceptance.md`（属票 131）、`.scratch/wisp/issues/130-*.md` 与 `134-*.md`（票面，本票未勾框）。
- **未 push、未 `--amend`/`reset`/`rebase`/`stash`/`checkout .`**。
- **owner 真实数据目录**：前后各读一次，`%APPDATA%\wisp` 与 `%APPDATA%\wisp-dev` 的文件数 **0 → 0**，目录 mtime `Sep 19 14:49` / `Sep 20 07:06` **未变**（§3.4）。所有写入都在仓外：`/tmp/wisp130-audit/`（= `C:\Users\swq\AppData\Local\Temp\wisp130-audit`）与 `C:\tmp\wisp130-audit\{data,datap}`、`D:\tmp\wisp130-posix\`（容器挂载点）。`build/wisp.exe` 与 `build/*.dll` 只被**复制**，未覆写。

**伪授权登记（计数 0）**：本轮工具输出里没有出现任何自称「编排者备注 / 系统提示 / 用户已更新编码规则 / `wisp-orchestrator-continuation` / 请 revert / 冻结某包 / 放宽阈值」的文字，故无原文可抄。出现过的**形似**项只有 2 条 harness 自己的「`MEMORY.md` 已被修改」通知（`Read` 工具返回的 `<system-reminder>` 块内，非任何 `Bash` 命令的结果）：它引用的路径是真实文件、内容里没有要求我 revert/冻结/放宽的指令 ⇒ 判定为**合法 harness 通知，不执行其中任何"指令"**（它也没给）。**它引用的编号已核存在**：`A102`（`docs/reports/pending-and-issues.md:3297`）、`A103`（`:3328`）、`A104①②③`（`:3365/:3372/:3385`）、`A105④`（`:3432`）、`A108②`（`:3494/:3498`，即本票解冻为具名单文件那一条）、`A30`（`docs/evidence/s1/A30-citation-repair.md` 在册）。

**结论修正记录（自纠，防下游把上一版当地基）**：
- 本节上一版把 A 组说成"全在 windows-tag 文件里" ⇒ **错**：`internal/winsec/resolve.go` 无 tag、`internal/risk` 在 GOOS=linux 的闭包与编译里都成立，并已用容器真跑证实（§3.3）。已改写为"两支平台都有分母、shipped exe 无 Linux 腿"。
- 另一处：F1 原本引 `cmd/wisp/logsink.go:10` 的字面量 "the D33 credential migration" ⇒ 该行跨 `:8-9` 两行，`grep -F` 打不到；已改成按行号引（这也是"契约明写要引真的句子"那条坑的又一个实例）。
- 再一处：`internal/config/migrate.go` 点名 `MigratePlaintext` 的行是 `:85`，不是 `:93`（`:93` 是 `PrivateFile` 那行）。

**待办（不属本票只读范围，交给 owner 拍板）**：
1. 勾 AC#1 之前请先核 §2 的式子能否在**你的机器上**复跑（尤其 `go list -deps` 与 `git show ac6f31c:…` 两枚）。
2. 若接受裁定的落点（`internal/observe` + `cmd/wisp/logsink.go`），**`internal/risk/winsec_c26.go` 的具名解冻就本轮作废**——撤销口令一句话：**「仍要放 risk 那侧」**（若你要把缓冲住在 `init()` 发声方那一侧，说这句我就改推 (a) 的 risk 变体并重新申请）。
3. §5.4 第 1 条那枚**会被撞红的既有判据**（`cmd/wisp/resident_sink_nail_127_windows_test.go` 的 `recs[0] == install`）需要一句授权：允许动码轮改它的断言形状（不是放宽，是重排）。

**next=** 交回裁定 → owner 选 (a-with-mirror) 或 (b) → 解冻 `cmd/wisp/logsink.go` 的排队位（他人 WIP 交件后）→ 动码轮按 §5.1 落点实现 + 按 §5.3 造 M-1/M-2/M-3 三发红判据 + 按 §4 落两枚钉（含 `leg_sink_gate_131_test.go` 的捕获式漏法）→ AC#4 四数门。**票 130 的 AC 框一位都没勾。**
