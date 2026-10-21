# 355-a1 只读普查腿：AC#1 三把尺逐名（part 2／2）

接 part 1：`docs/evidence/s1/355-a1-census.md`（sec 0 结论、sec 1 前提新鲜度、sec 2 尺1）。

---

## 3. 尺2：`<dataDir>\logs\wisp-*.jsonl` 的产码读取者名册

**射程声明**：三条子尺，全部对 HEAD blob。
(a) `git grep -n 'jsonl' HEAD -- cmd internal`（**含注释**、**含测试**，子串、大小写敏感）＝45 行；
(b) `git grep -n 'logFilePrefix\|logFileExt\|\.jsonl' HEAD -- cmd internal` 加 `CountLogFiles`／`logDir`／`logSinkDir` 全库点名；
(c) AST 形：`go/parser` 扫 770 枚 `.go`，取**调用表达式的渲染源串**含 `jsonl`／`LogDir`／`logFilePrefix`／`.logDir()` 的 CallExpr。
NOT 含 `.scratch/`（探针副本）、NOT 含 `frontend/`、`design/`（非 Go）。

### 3.1 名册（产码侧，逐名；全部是"写／定义／注释"，无一处读回）

| 行 | 逐字 | 性质 |
|---|---|---|
| `internal/observe/logging.go:31` | `logFilePrefix        = "wisp-"` | **常量定义**（命名规则） |
| `internal/observe/logging.go:32` | `logFileExt           = ".jsonl"` | **常量定义** |
| `internal/observe/logging.go:243` | `filepath.Join(w.dir, fmt.Sprintf("%s%s-%03d%s", logFilePrefix, w.day, w.seq, logFileExt))` | **写侧造名** |
| `internal/observe/logging.go:315` | `// named wisp-<day>-<seq>.jsonl; the day is parsed from the name first and` | **注释** |
| `internal/observe/logging.go:317-346 sweepLocked` | `os.ReadDir(w.dir)` `:318`，前缀／后缀过滤 `:325`，`os.Remove` `:341` | **删旧文件（保留期清扫），NOT 读内容** |
| `internal/observe/diagnostics.go:39` | `// LogDir is the rolling log directory (wisp-*.jsonl files).` | **注释** |
| `internal/observe/diagnostics.go:100-111` | `os.ReadDir(o.LogDir)` => `os.ReadFile(full)` `:111` | **唯一真读内容的产码点**，但它在 `BuildDiagnosticsBundle` 里（sec 4：产码调用者 0 枚）＝**不可达** |
| `cmd/wisp/panel_resident_windows.go:505` | `// wisp-20261014-001.jsonl carries exactly that sequence, and` | **注释**（讲 2026-10-14 的现场记录） |
| `cmd/wisp/resident_approval_windows.go:212` | `// put a false "recorded into this session" into <dataDir>\logs\*.jsonl - a` | **注释** |

产码侧 `jsonl` 命中 **5 行**（`panel_resident_windows.go:505`、`resident_approval_windows.go:212`、`diagnostics.go:39`、`logging.go:32`、`logging.go:315`）。
放宽到"命名规则／目录句柄"那族，另 4 枚产码点全是**写侧或删除侧**：`logging.go:31`、`:243`、`:317-346`、`cmd/wisp/logsink.go:87 logSinkDir()`
（`filepath.Join(dataDir, logDirName)`＝只拼路径，不读）。`CountLogFiles`（`logging.go:351`）只 `os.ReadDir` 数文件个数、不读内容，
且它自己的产码调用者也＝**0 枚**（见 3.3 反控）。
**没有任何一条产码路径把 `wisp-*.jsonl` 读回来解析。**

### 3.2 三枚 CLI 出口逐个排除（"会不会有别的读数口"）

- `wisp slo`：`slo_windows.go:272 observe.InitLog(...)` 是**另起一条管线写自己的日志**；全文件唯一 `os.ReadFile` 在 `:953 fileHas(path, needle)`，
  读的是 ready 标记文件（`:945-946` 先 `os.WriteFile` 写 `ready=1\npid=`）＝**不读 `wisp-*.jsonl`**。
- `wisp doctor`：全文件唯一 `os.ReadFile` 在 `:212`，读的是 `deps.toml`（`:207 p := filepath.Join(d, "deps.toml")`）；
  `git grep -n 'logs\|LogDir\|jsonl' HEAD -- cmd/wisp/doctor.go`＝**1 行注释**（`:243`，讲数据根布局）＝**doctor 零枚读日志目录**。
- `wisp panel-inbound`／`models`／`secret`／`providers`／`panel-assets`／`version`：`git grep -n 'jsonl' HEAD -- cmd/wisp/main.go`＝0；
  子命令名册（`main.go:89-124`）十个 `case`，无 `diagnostics`／`logs` 一类；`LogDir` 在 `cmd` 树产码侧只出现在 `diagnostics.go`（不属 cmd）。

### 3.3 反控：这把尺读 0，换第二把独立形复读

| 形 | 读数 |
|---|---|
| 形1（grep，产码加测试全射程）`jsonl` | 45 行＝**测试 40／产码 5**（产码 5 行见 3.1，全为注释或常量） |
| 形2（AST CallExpr 渲染源串含 `jsonl`） | `prod=0 test=36`＝**产码 0 个调用点** |
| 形3（AST，含 `o.LogDir`／`w.dir` 的调用） | `o.LogDir`＝`prod=2`（均在 `diagnostics.go:100/:110`）；`w.dir`＝`prod=4`（`logging.go:243/:318/:341` 两处）＝全在写／删／清扫侧 |
| 形4（`.logDir()` 的调用者） | `prod=0`，`test=2`（同一行 `cmd/wisp/logsink_test.go:158`） |

**形1 与形2 同号于"产码读取者 0"**：形1 那 5 行产码命中逐行过目全是注释／常量；形2 的 AST 判据直接不把注释算调用。
另**同族反控**：`CountLogFiles`（`logging.go:351`）grep 形产码 **0**、AST 形 `prod=1`（仅 `decl`，6 枚 `call` 全在测试里）＝两形同号。

### 3.4 正控（报 0 之前先证明尺能命中非零）

| 正控样本 | 同尺读数 | 结论 |
|---|---|---|
| `os.ReadDir`（AST 形2 同款，全树） | **prod=14 test=47** | 同一把 AST 尺能命中产码调用点，非零 |
| `os.ReadDir` 在 `internal/observe` 产码侧（grep 形） | **3** | 尺在"本包内"也打得开 |
| `logFilePrefix` AST 调用形 | **prod=7 test=5** | 同一把 AST 尺在"日志文件名"族上非零 |
| `wisp-*.jsonl` 字面（全树任意射程） | **37 行** | 字面在树里存在，尺没打错目标 |
| `filepath.Glob`（另一枚"找文件"形状） | **prod=0 test=17** | 这条本身是 0，故**不**当正控；它作**旁证**：连另一形状的遍历也只在测试侧 |

**尺2 成立：`<dataDir>\logs\wisp-*.jsonl` 的产码读取者＝0 枚**（唯一会读内容的 `diagnostics.go:111` 不可达，且它自己也零调用者）。

---

## 4. 尺3：`BuildDiagnosticsBundle` 的产码调用者名册

**射程声明**：`git grep -n 'BuildDiagnosticsBundle' HEAD -- .`（**全仓，含注释、含测试、含 `.scratch/`**，子串）；
再以 AST 形（`CallExpr` 的 callee 名等于 `BuildDiagnosticsBundle`，与 `FuncDecl` 分开计数）复核。
**逐名区分"声明／注释"与"真调用"，NOT 混算。**

### 4.1 名册（唯一那枚文件，逐名）

| 行 | 逐字 | 性质 |
|---|---|---|
| `internal/observe/diagnostics.go:35` | `// BundleOptions configures BuildDiagnosticsBundle.` | **注释** |
| `internal/observe/diagnostics.go:61` | `// BuildDiagnosticsBundle writes the zip and returns its path.` | **注释** |
| `internal/observe/diagnostics.go:62` | `func BuildDiagnosticsBundle(o BundleOptions) (string, error) {` | **声明（定义本身）** |
| `internal/observe/diagnostics_test.go:35` | `path, err := BuildDiagnosticsBundle(BundleOptions{` | **真调用（测试）** |
| `internal/observe/diagnostics_test.go:125` | `_, err := BuildDiagnosticsBundle(BundleOptions{` | **真调用（测试）** |
| `internal/observe/diagnostics_test.go:140` | `if _, err := BuildDiagnosticsBundle(BundleOptions{OutPath: out}); err != nil {` | **真调用（测试）** |

`cmd`／`internal` 树内命中 6 枚＝声明 1 ＋ 注释 2 ＋ 测试真调用 3；**产码真调用 0 枚**。
全仓 `-- .` 另命中 `.scratch/` 与 `docs/` 的历史叙述（如 `docs/evidence/s1/139-v1-verdict.md`、`.scratch/wisp/probes/167/...`）；
那些是**文字**，NOT 调用，不计入名册。

### 4.2 第二把独立形（AST）复核

```
decl      internal/observe/diagnostics.go:62        func BuildDiagnosticsBundle
call      internal/observe/diagnostics_test.go:125  BuildDiagnosticsBundle(...)
call      internal/observe/diagnostics_test.go:140  BuildDiagnosticsBundle(...)
call      internal/observe/diagnostics_test.go:35   BuildDiagnosticsBundle(...)
# TOTAL hits: prod=1 test=3
```
与 grep 形同号：**产码侧 1 枚＝decl；call 0 枚。**

### 4.3 正控（同尺打已知非零样本）

| 样本 | 同尺读数 | 结论 |
|---|---|---|
| `observe.InitLog`（同包，已知**有**产码调用者） | grep 形：`cmd/wisp/logsink.go:149` 加 `slo_windows.go:272`＝**2 枚产码真调用**；AST 形：`prod=3`（2 call 加 1 decl） | 尺能命中产码真调用，并分开 decl／call |
| `logSinkDir`（同包兄弟） | grep 形产码 4 行（`:49` 注释、`:84` 注释、`:87` 定义、`:148` **1 枚真调用**） | 尺能分辨"注释里出现"与"真调用" |

**尺3 成立：`BuildDiagnosticsBundle` 产码调用者 0 枚。**

---

## 5. 与票面三项判语的逐条对照（票面第 8 行那句）

| 票面原话 | 本腿现量 | 判 |
|---|---|---|
| "`cmd/wisp` 对 `Compression` 零读取" | `git grep -n 'Compression' HEAD -- cmd`＝**0 行（含测试）**；`compress` 大小写不敏感＝**0** | **成立，且比票面更强**（连出现都没有） |
| "`BuildDiagnosticsBundle` 产码调用者 0" | sec 4：decl 1／comment 2／test-call 3／**产码 call 0** | **成立** |
| "`.logDir()` 产码调用者 0" | `git grep -n '\.logDir()' HEAD -- .`＝产码 **0**、测试 1（`cmd/wisp/logsink_test.go:158`）；AST 形 `prod=0 test=2` | **成立**（两形同号） |

票面第 8 行与第 22 行之外，本腿另量到一枚票面没写、但 AC#1 第 1 项射程点名的事实：
`internal/observe/diagnostics.go` 对 `Compression` 命中 **0**，且它对日志内容的读（`:111`）位于一枚**不可达**的函数里——
"诊断包一旦接上就顺带带走这条痕"这个说法**今天只是潜力，不是通路**（与 139-census H6 那句同判）。

---

## 6. 门禁（本腿现跑，rc 自己拿）

| 门 | 命令 | rc | 说明 |
|---|---|---|---|
| 构建 | `GOFLAGS= go build ./...` | **0** | 12s。同一命令加 `CGO_ENABLED=0` 时 rc=**1**（`sherpa-onnx-go-windows` 全文件被 build constraint 排除），故本表 rc=0 是**未改 CGO 环境**下的读数 |
| D22 | `sh scripts/d22scan.sh` | **0** | 67.5s。逐字判语：`d22scan: clean - no D22 ban violations; live scope work: bans #1-5 internal/=244, bans #1-5 cmd/=42, ban #6 frontend/=116, ban #7 internal/tools/=25, ban #8 design/=39, ban #8 frontend/=116, ban #8 internal/=608, ban #8 cmd/=143`；测试侧 `runtests.sh: OK - packages=[./...] top-level: PASS=35 FAIL=0 SKIP=0` |

两发日志原件在仓外：`D:\tmp\wisp\20261021\355a1\gate-gobuild.log`、`gate-d22scan.log`。

---

## 7. 交件自检

- 逐枚 `.go`＝**0**（本腿对仓内任何 `.go` 零字节；AST 探针与导出树全在 `D:\tmp\wisp\20261021\355a1\`）。
- 越界枚数＝**0**：`internal/**`、`cmd/**`、`scripts/**`、`.github/**`、`.scratch/wisp/**`、`docs/reports/**` 各 0 字节。
- 写面＝本件两枚（part 1 加 part 2）。正控尺＝sec 2.4、3.4、4.3 三张表（每把尺各有非零样本）。
- 反控＝sec 3.3（读 0 那把用四形复读，两形同号）加 sec 3.3 末尾的 `CountLogFiles` 同族反控。
- 未覆盖／未假设：`AC#2` 的落点选型（H3 与 H6 之间）**本腿不选**（属 `355-r1`）；`Q-312` 契约面本腿零触碰。
