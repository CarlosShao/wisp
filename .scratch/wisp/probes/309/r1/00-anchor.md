# 309-r1 · 锚件（落地腿，形 Ⓑ 已裁）

## 进场三读数（现跑）
- `date`：`Sun Oct 11 12:51:25 CST 2026`
- `git rev-parse HEAD`：`24929768dc41eb32715e4823deaee297aa4e95cc`
- `git status --porcelain -- internal cmd`：**0 行**（台面 `internal`/`cmd` 干净）

## 形状（编排者已裁，见票面 `AC#0` 裁语 :53；本轮⛔ 重开、⛔ 改判语）
裁＝**形 Ⓑ（写明故意）＋ `AC#2` 枚数钉**；⛔ 形 ⓐ（入口拒绝）。撤销口令＝「309 改形 ⓐ」（⛔ 我引）。
⇒ 本轮⛔ 改任何行为、⛔ 在桥入口加拒绝分支；`b.assess` 返回路径逐字节不变。

## 动笔前逐枚 `sed` 复核的现场事实（行号＝HEAD blob，⛔ 用绝对行号作判据锚）
- `bridge.go:924` 函数签名 / `:925` 守卫 `if b.injected || b.prov == nil || taskID == "" {` / `:926` `return b.assess` / `:931` `WithTaintDetector(b.prov.Detector(taskID))`。文档块 `:920-923`。⇒ **未漂**，与票面一致。
- `bridge.go:205-216` `else` 支 `b.assess = NewRiskAssessor().WithCanonicalizer(o.Paths).WithSensitiveClassifier(b.classifier)`，**无 `WithTaintDetector`**。⇒ 探测器没挂＝结构事实。
- `internal/risk/rules_gateway.go:102` 逐字 `if ctx.taint == nil {`（`ruleTaint`）⇒ 空 id 那次判定门照走、R4 整条不发生。
- `internal/risk/provenance.go` `scopeMarks`：`reg != nil && len(marks)>0`→(marks,false)；`reg != nil`空→(nil,false)；`reg == nil` 且**另有 scope 带 marks**→(nil,true)。消费者 `Inspect`（`if unbound` 现量 `:727`，字面量 `:729`）、`CheckText`（`if unbound` 现量 `:827`，字面量 `:830`）。⇒ 注释写"看情形转误报"，⛔ 写"无条件风暴"。
- 生产构造点：`grep -rn 'ToolRequest{' internal cmd --include=*.go` 非 `_test.go` ⇒ **仅 `internal/agent/loop.go:675`**（`req := ToolRequest{`）；`reflect.TypeOf(agent.ToolRequest{})` 2 枚（`grant_test.go:509`/`ticket90_test.go:454`，均在测试里）。全仓 `ToolRequest{` **60** 枚。⇒ `internal/agent` 之外生产构造点 **0**。

## AC#3 复核（现读，U-5）
- `bridge.go:861` 引 `internal/agent/loop.go:366` 说 TEXT loop "already defers" AdmitTask 的 revoke ⇒ **已漂**：现量 `defer revokeAdmission()` 在 `loop.go:379`（`:366` 落在 journal 块）。⇒ 换内容锚。
- `bridge.go:871` 引 `loop.go:332/:321`（RunAsync/Run→newTaskID）⇒ **已漂**：现量 `RunAsync` 起 `:328`（铸造 `:329` `id := newTaskID()`）、`Run` 起 `:339`（`:340` `return l.run(ctx, newTaskID(), input)`）。⇒ 换内容锚（认两次 `newTaskID()` 调用本身）。
- `bridge.go:230`（`loop.go:783-806`）、`:252`（`loop.go:692-694`）⇒ AC#3 射程之外，⛔ 动。

## 指名的两枚新用例（注释里点名的 = 真写出的这枚）
- `AC#1` 恒真形仪器：`Test309AssessorForEmptyTaskIDDropsTaintDetector`（`internal/tools/assessorfor_taint_instrument_309_test.go`）
- `AC#2` 枚数钉：`Test309ToolRequestProductionConstructionNail`（`internal/tools/toolrequest_nail_309_test.go`）
名册尺：`grep -rn 'func Test309' internal --include=*.go` ⇒ 起手 0 枚（无冲突）。

## 写面声明（越面即停手上报）
只写：`internal/tools/**` ＋ `.scratch/wisp/probes/309/r1/**` ＋ 票 309 末尾 1 条 append-only bullet（`--numstat` 必 1 加 0 删）。
⛔ 碰：`internal/risk/**`、`agent.ToolRequest` 加字段、`lit236*` 常数值、`bridge.go:925` 布尔式/那一肢、`b.injected`/`Options.Assessor` 肢、`bridge.go:869-891` 里⛔ 行号/枚数的字、`thresholds.go`/`golden`/`allowlist.txt`/D43/C1-C32/D1-D47/三枚冻结件/`frontend`/`design`/`PLAN.md`/`docs/specs`/`docs/reports/**`。
⛔ 碰别的腿：`probes/orch/20261011/writeface-1/**`（另一枚腿 writeface-1 正在写）、`308-a1`、`298-v1`（已交件）；`probes/309/a1/**`（普查腿，只读）。
Git：⛔ push/amend/reset/rebase/stash/checkout ./clean/--no-verify/add -A/add .；每笔显式 pathspec；`.gitignore`+`go.mod` 未提交改动⛔ 我名下（⛔ 动/提交/revert）。
