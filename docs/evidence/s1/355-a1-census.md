# 355-a1 只读普查腿：AC#1 三把尺逐名（part 1／2）

- 腿＝`355-a1`（只读腿）。写面＝本件两枚（`.go` 零字节；`internal/**`、`cmd/**`、`.scratch/wisp/**`、`docs/reports/**` 零字节）。
- 起手锚＝`c4098fc80f68e678dc0854e87148c0fb5e32d82e`（`dev`，2026-10-21 13:06:20 +08）。
- 全部读数对 **HEAD blob** 量（`git cat-file blob HEAD:<path>`、`git grep ... HEAD`、`git cat-file -e HEAD:<path>`），
  NOT 工作树。第二形尺＝仓外 AST 探针 `D:\tmp\wisp\20261021\355a1\probe\`（`go/parser`，起于
  `git archive HEAD internal cmd tools` 导出的 `D:\tmp\wisp\20261021\355a1\head\`＝770 枚 `.go`：产码 300／测试 470）。
- 本件分两枚：part 1 ＝ sec 0 至 sec 2（结论、前提新鲜度、尺1）；part 2 ＝ sec 3 至 sec 7（尺2、尺3、对照、门禁、自检）。

---

## 0. 结论一句

**"生产侧零读者"成立**，三把尺逐名支撑，每把两形独立同号：
(1) `res.Compression` 产码命中 1 枚且是赋值（写），**读取者 0 枚**；
(2) `<dataDir>\logs\wisp-*.jsonl` 的**产码读取者 0 枚**（产码侧 5 行 `jsonl` 全是注释或常量，无一处读回）；
(3) `BuildDiagnosticsBundle` 产码**调用者 0 枚**（产码侧 1 枚＝定义本身）。

---

## 1. 前提新鲜度（票面点名落件，逐枚现跑）

尺 A＝`git cat-file -e HEAD:<path>`（在不在）；尺 B＝票面那句可判定断言对 HEAD blob 复量。
`docs/evidence/s1/355-a1-census.md` 落笔前 rc=128（本件自身，正常）。

| # | 票面点名的落件 | 尺A | 尺B（断言对 HEAD blob 复量） |
|---|---|---|---|
| 1 | 票面 `.scratch/wisp/issues/355-...-never-reads-it.md` | rc=0 | 第 8 行三项判语本轮逐项量过，见 sec 2／sec 3／sec 4（part 2） |
| 2 | `docs/evidence/s1/139-v1-verdict.md`（非实现者签名来源） | rc=0 | 第 5 件 `:18`"`cmd/wisp` 对 `Compression` 命中 0"＝**仍真**（`git grep -n 'Compression' HEAD -- cmd`＝0 行，含测试亦 0）；"`BuildDiagnosticsBundle` 非测试调用者 0"＝**仍真** |
| 3 | `.scratch/wisp/probes/139/v1/02-new-ticket-draft.md`（SEC-B 四格原文） | rc=0 | `:12` 三把尺本轮逐把现跑同号（见 sec 2／3／4 各"与草案对照"行） |
| 4 | `.scratch/wisp/probes/139/a1/census.md`（407 行，10-03 读数） | rc=0 | 该件行号**已漂**：引 `run.go:1144`＝现 `:1222`；`loop.go:399/402/404`＝现 `:406/415/417`；`compress.go:216` 与 `:240` 未漂。详见 sec 1.1 |
| 5 | `.scratch/wisp/issues/139-...-done.md` | rc=0 | 第 12 节 `:133`"`BuildDiagnosticsBundle` 非测试调用者 0"＝仍真 |
| 6 | `.scratch/wisp/issues/153-...-done.md`（痕的两形补齐） | rc=0 | `compress_trace_test.go` 现 7 枚 `func Test`（`:177/:269/:302/:366/:426/:478/:534`）＝与 `139-v1` 第 2 件报的"7 枚"同号 |
| 7 | `internal/agent/compress.go` | rc=0 | `:240 c.log().Info("agent: history compressed", attrs...)` 在场；八计数＋一布尔 `:222-231`；`"task"` 键 `:237-239` |
| 8 | `internal/agent/loop.go` | rc=0 | `:100 Compression CompressionReport`；`Need` 门 `:406`；`:415 } else if rep.Ran {` => `:416 l.replaceHistory(nh)` => `:417 res.Compression = rep` |
| 9 | `internal/observe/diagnostics.go` | rc=0 | `:62 func BuildDiagnosticsBundle(o BundleOptions) (string, error)`；`:107` 用 `logFilePrefix`/`logFileExt` 选文件；`:111 os.ReadFile(full)` |
| 10 | `internal/observe/logging.go` | rc=0 | `:31 logFilePrefix = "wisp-"`、`:32 logFileExt = ".jsonl"`；`:243` 造名；`:317-346` 保留期清扫；`:351 CountLogFiles` |
| 11 | `cmd/wisp/run.go` | rc=0 | 收尾摘要行 `:1221-1222`；`costLine` `:1272-1279` 只读 `CostMicros/Currency/Usage`；`notifyBody` `:1261-1269` 只读 `Text/Message/Status`；`res := bg.Wait()` `:1183` |
| 12 | `cmd/wisp/logsink.go` | rc=0 | `:99 func (s *logSink) logDir()`；`:122 logger()`；`:82 logSinkLevel = "info"`；`:149 observe.InitLog(...)` |
| 13 | `internal/memory/models.go`／`schema.go`（H3b 载体） | rc=0 | `models.go:58-70 TaskLog` 11 字段、`schema.go:53-65 task_log` 12 列＝都放不下压缩读数（与 139-census 同号） |
| 14 | `internal/panel/composer.go`（H5 契约面） | rc=0 | `:57-92 Snapshot` 七键（`pending/results/composer/generatedAt/instructions/tasks` 等），无压缩维 |
| 15 | `docs/specs/SPEC-08-ui-ball-panel.md` | rc=0 | `:156` 第 5.2 节"C17 PanelBridge 方法白名单"在（本腿只读不改） |
| 16 | `docs/reports/pending-and-issues.md` | rc=0 | `Q-312` 命中 5 行、`A1071` 命中 8 行（票面第 9／24 行引的两枚账都在） |

### 1.1 特别复量：票面第 12 节"管道七跳三态"（10-03 读数，会漂）

**判语：七跳三态全对（H0/H1/H2 已接、H3 差一行、H3b 与 H5 未写且属契约面、H6 未接、H7 到盘无读者）；
但该节引的 `file:line` 有五跳已漂。** 逐跳对 HEAD blob 复量：

| 跳 | 10-03 引的行 | HEAD 现值 | 三态 |
|---|---|---|---|
| H0 产出 | `compress.go:174/:216/:240` | `:174` 函数头、`:216 if rep.Ran {`、`:240` 那枚 `Info`＝**同号未漂** | 已接 |
| H1 上到 Result | `loop.go:399` => `:402` => `:404` | 现 `:406 if l.comp.Need(hist) {` => `:415 } else if rep.Ran {` => `:417 res.Compression = rep`＝**漂** | 已接 |
| H2 出 run 到宿主 | `loop.go:369/:326/:308-314`、`run.go:1106` | `loop.go:369 j := newTaskJournal(...)`（同号）、`:333 t.result = l.run(...)`（漂）、`:315 func (t *RunningTask) Wait()`（漂）；`run.go:1183 res := bg.Wait()`（**漂 77 行**） | 已接 |
| H3 控制台摘要行 | `run.go:1144-1145`、`costLine:1195-1201`、`notifyBody:1184-1192` | 现 `run.go:1221-1222`（**漂 77 行**）、`costLine:1272-1279`、`notifyBody:1261-1269` | **差一行**（成立） |
| H3b `task_log` | `run.go:1119-1140`、`models.go:58-70`、`schema.go:53-65` | 现 `run.go:1200 tl := memory.TaskLog{`（漂）；models／schema 两处**同号** | 一整块没写（契约面） |
| H5 面板快照 | `run.go:998`、`panel_pump.go:401-411`、`composer.go:57-91` | 现 `run.go:1075 Sink: consoleSink{... publish: rt.publishPanelSnapshot}`（漂）、`panel_pump.go:447 func (rt *agentRuntime) publishPanelSnapshot()`（漂）、`composer.go:57-92`（同号） | 一整块没写，且属契约面 |
| H6 诊断包 | `diagnostics.go:62/:39/:96-120` | `:62` 与 `:39` 同号；抄文件循环实为 `:99-130` | 一整块没写（接线缺失，sec 4） |
| H7 日志文件本身 | `logsink.go:144/:153/:208-214`、`logging.go:87-90/:82` | `logsink.go:144`／`:153`／`:208-214 teeHandler.Handle` 同号；`logging.go:87-90` 同号；`logging.go:82` 现是 `_ = w.Close()`（`min_level` 语义落在 `cmd/wisp/logsink.go:82 logSinkLevel = "info"` 与 `:184` 记录里的 `min_level` 键） | 已接但今天没有读者 |

结论：那句"会漂"**兑现了**——H3／H5 两跳落点各偏 77 行，H1／H2 偏 7 行。三态判语本身一跳动未变，逐跳已证。

---

## 2. 尺1：`res.Compression` 的非测试读取者名册

**射程声明**：`git grep -n 'Compression' HEAD -- cmd internal`，**含注释行**、**含测试文件**，再按 `_test` 切口分账；
子串匹配（`Compression`／`CompressionReport`／`Compress` 都算命中），因此另跑窄形 `.Compression` 校正。
NOT 含 `tools/`（本轮现量 `git grep -n 'Compression' HEAD -- tools`＝**0**，射程不含它不改变结论）；
NOT 含 `docs/`、`.scratch/`（票面与探针副本里命中很多，与产码无关）。

### 2.1 名册（产码侧，逐名）

| 行 | 逐字 | 性质 |
|---|---|---|
| `internal/agent/compress.go:42` | `// CompressionReport describes one compression pass (assertion + log surface).` | **注释** |
| `internal/agent/compress.go:61` | `type CompressionReport struct {` | **类型声明** |
| `internal/agent/compress.go:142` | `// TestCompressionTraceNeverInventsATaskID. Today the only non-test caller is` | **注释** |
| `internal/agent/compress.go:174` | `func (c *Compressor) Compress(ctx, hist) ([]llm.Message, CompressionReport, error) {` | **声明（返回类型）** |
| `internal/agent/compress.go:175` | `rep := CompressionReport{TokensBefore: c.TotalTokens(hist)}` | **构造（压缩机自己）** |
| `internal/agent/loop.go:99` | `// Compression reports the last history compression pass (zero if none).` | **注释** |
| `internal/agent/loop.go:100` | `Compression CompressionReport` | **字段声明**（`Result` 的一枚） |
| `internal/agent/loop.go:417` | `res.Compression = rep` | **写**（等号左边，生产者，NOT 读取者） |

产码侧共 8 行，按名分账＝注释 3（`:42/:142/:99`）＋声明 3（`:61/:174/:100`）＋构造 1（`:175`）＋写 1（`:417`）＝**读取者 0 枚**。

**窄形校正**（射程只认 `.Compression` 选择子，产码侧）：`git grep -n '\.Compression' HEAD -- '*.go'` 排 `_test.go`，在 `cmd`／`internal` 下
**只有 `internal/agent/loop.go:417` 那 1 行**，且是赋值。两把尺同号：**读取者 0**。

### 2.2 测试侧逐名（不计入，供判名册唯一性）

`internal/agent/compress_trace_test.go:343/:346/:349/:521`＝4 枚真读（`res.Compression.TokensBefore/.TokensAfter/.Ran`）；
`:21/:25/:300/:405` 是注释；`:177/:269/:302/:366/:426/:478/:534` 是函数名；`:344/:347/:350/:440` 是失败句里的 `Result.Compression` 字面文本。
另有 `.scratch/wisp/probes/155-accept/30-zzaccept155-parsed-ondisk_test.go:61`（仓外探针副本，不属产码树）。

### 2.3 射程内指定要查的三个面，逐个具名

- **`cmd/wisp/**`**：`git grep -n 'Compression' HEAD -- cmd`＝**0 行（含测试 0）**；`git grep -n 'Compress' HEAD -- cmd`＝**0 行**；
  `git grep -in 'compress' HEAD -- internal/panel cmd/wisp`＝**0 行**。宿主侧连"压缩"这个字都不出现。
- **面板桥**：`internal/panel` 对 `Compression`＝**0**（同尺）；`internal/panel/composer.go:57-92` 的 `Snapshot` 七键无压缩维。
  另：`cmd/wisp/panel_control_windows.go:746 json.Marshal(res)` 的 `res` 是 `c.callbinding(d)` 的返回值（`:742`），**不是** `agent.Result`
  （`run.go:1261/:1272` 是 `agent.Result` 仅有的两枚产码读者）＝不存在"通用序列化顺带写出压缩维"的暗道。
- **`internal/observe/diagnostics.go`**：`git grep -n 'Compression' HEAD -- internal/observe`＝**0**。

### 2.4 正控（同尺打已知非零样本，证明尺不瞎）

| 样本 | 同尺读数 | 说明 |
|---|---|---|
| `.Rounds`（同族 `Result` 字段，已知有真读者） | **3 行**：`cmd/wisp/run.go:1222`（读）、`internal/agent/loop.go:431`（写）、`:961`（写） | 尺能命中产码读者，并能分辨读／写 |
| `.Usage`（同族） | **54 行** | 尺在大样本上正常 |
| `Compression` 打在 `docs/` 上 | **67 行** | 同串放宽射程时打得开（0 是射程事实，NOT 正则失效） |
| AST 形（第二把）`-target Rounds -mode field` | `prod=5 test=2`，其中 `read cmd/wisp/run.go:1222 res.Rounds` | 第二形尺同样能判读／写 |

尺1 两形同号：**`res.Compression` 产码读取者 0 枚，唯一产码命中 `loop.go:417` 是写。**

（接 part 2：`docs/evidence/s1/355-a1-census-part2.md`）
