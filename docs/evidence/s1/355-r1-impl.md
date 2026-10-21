# 355-r1 写腿（part 1／2）：选型、闸门、改了哪几行

- 腿＝`355-r1`（写腿／实现者）。票面＝`.scratch/wisp/issues/355-...-never-reads-it.md`（`AC#2`/`AC#3` 是本腿范围）。
- 起手锚 `ac2847d5` => 落笔 `5b54c1b4`（漂移账见 sec 1）。**只 commit、零 push。**
- **本件按编排者裁决 `a7f26079` 落笔**：走主 flow（给 `CompressionReport` 加只读字段），NOT 我中途提的 slog tap。
  裁决原文＝`.scratch/wisp/probes/orch/20261021/txt/355r1-ruling-three-flavours.txt`。
- 本件分两枚：part 1 ＝ 头部、sec 0 选型、sec 1 闸门、sec 2 改了哪几行；
  part 2 ＝ sec 3 新用例、sec 4 三读数、sec 5 门禁、sec 6 未覆盖（`355-r1-impl-part2.md`）。

---

## 0. 选型一句（AC#2 要求的"为什么选它"）

**选 H3（`cmd/wisp/run.go` 收尾摘要那一带）＋主 flow（`CompressionReport` 加只读 `HistoryChanged`），不选 H6（诊断包）。**

1. **三味必须先"答得出"才谈出口。** 现量：`history_changed` **不是** `CompressionReport` 的字段（七枚字段表无此味），
   只在 `compress.go:247` 那枚 slog 的 attrs 里现算。`cmd/wisp` 侧还原不出迭代次数
   （`len(out) = len(hist) - CompressedMsgs + I`，I 不可见）=> 单看 `Result.Compression` **只能答两味**。
   故"加一枚只读字段"是前置，不是附加项。
2. **H3 是唯一今天就有产码读者的出口。** 现量：`BuildDiagnosticsBundle` 产码调用者 **0 枚**、
   `internal/observe/diagnostics.go` 对 `Compression` **0 命中**。选 H6 要同时新建一条 CLI 动词/接线＝"新出口"，
   超出票面"把已有的那条痕变成可答"。H3 的收尾行每跑一次 `wisp run` 就过一次。
3. **主流 harness 的形状**（AGENTS.md 铁律 000）：返回对象带结构化读数，NOT 从日志反捞状态。
   `cmd/wisp` 全包 `Logger:` 现量 **0 处** => 压缩机走 `slog.Default()`。我那版 slog tap 捞得出三味，
   但读数是"最后一条日志行"这个**时序假设**（级别一改／handler 链一改／被采样就静默变空），
   且分不出"没读数"与"读数真是 0" = 本仓反复咬人的假绿。**裁决否决它，我照办。**

**代价（写清落点，NOT 当降级理由）**：
- (a) struct 多一枚字段 => 位置式字面量构造者会受影响。现量：`internal/agent` 内 `CompressionReport{` 只有
  `compress.go:175`（键值式）与测试，**无位置式**，故零影响。
- (b) 收尾行**多一行**（仅真压过历史时）=> 会撞"行数/结尾行"断言。现量：`cmd/wisp` 测试无一处断言 stdout 最后一行
  （`git grep 'lines\[len(lines)-1\]'`＝0），故零影响。
- (c) 常驻腿（resident）走同一个 `execute()` => 同样受益，不需第二条出口。

---

## 1. 开工前闸门与前提新鲜度（尺 v4 8）

| # | 检查 | 现跑读数 |
|---|---|---|
| 1 | `git status --porcelain -- cmd/ internal/` | **空**（起手与落笔各一次）=> `cmd/wisp/**` 无别人未提交改动，**通过** |
| 2 | 锚点漂移 | `ac2847d5` => `5b54c1b4`（中途另有 `6633c45d`）。期间 `cmd/`／`internal/` 的改动**零字节**（裁决 `a7f26079` (4) 同判）=> 读数不作废 |
| 3 | 派单引的 `:1221-1222` | **仍准**：收尾 `Fprintf` 在 **`:1221-1222`**（参数行 `:1222`） |
| 4 | `costLine` | **仍准**：`:1272-1279` 只读 `CostMicros`／`Currency`／`Usage` => 三味**塞不进** |
| 5 | `notifyBody` | **仍准**：`:1261-1269` 只读 `Text`／`Message`／`Status` => 三味**塞不进** => 须**新起一段** |
| 6 | `git grep -n 'Compression' HEAD -- cmd` | **0 行**（含测试）=> 票面"零读取"成立 |
| 7 | `BuildDiagnosticsBundle` 产码调用者 | **0 枚**：全仓 `.go` 命中 6 行＝decl 1（`diagnostics.go:62`）＋注释 2（`:35`/`:61`）＋**测试 call 3**（`diagnostics_test.go:35`/`:125`/`:140`） |
| 8 | `internal/observe` 对 `Compression` | **0** => H6 今天只是潜力，NOT 通路 |
| 9 | `CompressionReport` 字段表 | **七枚**：`Ran`/`CompressedMsgs`/`KeptRawRounds`/`PreservedIDs`/`RawIDsInsideSummary`/`TokensBefore`/`TokensAfter` => **无 `history_changed`** |
| 10 | `cmd/wisp` 设 `Logger:` 处数 | **0** => 压缩机走 `slog.Default()`（tap 路因此可行，但裁决否决） |

**派单/普查件的一处空洞（本腿报缺陷，NOT 我发明第三出口）**：派单把三味当成 `Result.Compression` 里现成的东西，
但 `history_changed` **不在** struct 上 => "三味可见"必须先补这一枚只读字段（裁决已批，见 sec 0）。

---

## 2. 改了哪几行（两枚文件，逐处）

**(1) `internal/agent/compress.go`（+18/−1；裁决只许两处）**
- `:71-81` `CompressionReport` 加只读字段 `HistoryChanged bool`（注释写明：它是读数、NOT `Ran` 的副本、两者可矛盾）。
- `:238` `rep.HistoryChanged = rep.TokensAfter != rep.TokensBefore || len(out) != len(hist)`
  => **判定语义一字未改**（同一表达式），只是把值**同时**存进 rep。
- `:247` attrs 改读 `"history_changed", rep.HistoryChanged`（原来把 `:238` 那串表达式**重算一遍**）。
  => 记录与调用方**同一处取值**，此后不可能漂。**attrs 的键名/键序/其它值一字未改。**

**(2) `cmd/wisp/run.go`（+46/−0）**
- `:1223-1229` 收尾行之后**新起一段**：`printCompressionTrace(rt.stdout, res.Compression)`。
- `:1288-1301` `printCompressionTrace(w io.Writer, rep) bool`：写一行并报"是否写了"。抽成函数＝
  让**测试驱动产码那一条路径**，NOT 测试自拼字符串。
- `:1303-1325` `compressionLine(rep) (string, bool)`：`rep.Ran == false => ("", false)`（静默），否则渲染
  `压前 N tokens => 压后 N tokens（history_changed=<bool>，M 条消息被折进摘要，保留 K 轮原文）`。

**零越界**：NOT 前端事件／NOT 新 `panel.*` 方法／NOT 新 JSON 键（输出是**控制台文本行**，不是 wire 形状）；
NOT `Q-312` 那根链（`KeepRawRounds`／`Need()` 零触碰）；NOT `internal/agent` 其它文件（零字节）。
