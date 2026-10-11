# 309-r1 · 90-final（落地腿收口；六格中我这条腿交 AC#1/AC#2/AC#3/AC#4 的实现面 + Ⓑ 注释，翻框归 `309-v1`）

台面：起手 HEAD `24929768`；写期间 `308-r1`/`writeface-1` 在 `dev` 上追加了 `510c95ca`/`1ad48333`（共享工作树，正常）。我已核：这两笔**未触**我写面交集（`bridge.go`/`*_309_*`/`tools/` 相交＝空），我的 `bridge.go` diff 对当前 HEAD 仍 `32 3`。所有门禁 rc 都在最终 HEAD 上现跑。

## ① 交了几笔 commit（写面外枚数应 0）
| 笔 | 短 sha | 内容 | `git show --name-only` 文件数 | 越写面枚数 |
|---|---|---|---|---|
| commit1 | `1cd57ec1` | 锚件（commit-first） | 1 | 0 |
| commit2 | `16e86bae` | Ⓑ 注释 + AC#3 重述 + AC#1 仪器 + AC#2 枚数钉（`bridge.go`＋两 `_test.go`） | 3 | 0 |
| commit3 | (本笔) | 证据件 `10/20/30/40/90.md` + `logs/*` + 票 309 末尾 1 bullet | 证据件+logs，`git diff --cached --name-only` 提交前数并落 `logs/staged-evidence.txt` | 0（票 bullet 属写面内许可项） |
尺＝`git show --name-only --format= <sha> | grep -c .` ⇄ 滤 `^internal/tools/|^\.scratch/wisp/probes/309/r1/|^\.scratch/wisp/issues/309-` 取反计数。

## ② 五格逐格答了什么
- `AC#0`：编排者已裁（形 Ⓑ），我⛔ 重开、只按裁语落地。
- `AC#1`：交了恒真形仪器 `Test309AssessorForEmptyTaskIDDropsTaintDetector`，两向都有牙（突变 Rig-A/Rig-B 各红，见 `10-*.md`）。**闭合。**
- `AC#2`：交了枚数钉 `Test309ToolRequestProductionConstructionNail`（自带分母 245、内置正控、盘上可核），Rig 证"新构造点即红"。**闭合。**
- `AC#3`：三处文案（`:861`/`:871`/`877`）换内容锚 + 带尺现量重述，另现读并修正 `:861` 的 `loop.go:366`→`defer revokeAdmission()`（U-5 已答：漂了）。**闭合。**
- `AC#4`：三枚门禁 rc=0、格式名册两尺 `internal/tools` 0↔0、`cmd/wisp` 既有债 2 枚未动、逐名作差纯新增。**闭合。**
- `AC#5`：见本件末节。**闭合。**
⇒ 我这条腿五格全交；**⛔ 翻任何框**（`AC#0` 已由编排者翻，其余五格归非实现者验收腿 `309-v1`）。

## ③ 欠读数清单（我⛔ 代答，写明谁用什么尺答）
1. **`assessor.go:165-171` 这枚既有文档锚是否仍准**：`bridge.go:922` 现引它指 C19 `TaintDetector` seam。我目测 seam 在 ~164-169，但 AC#3 射程⛔ 含它、我没改。尺＝`grep -n 'type TaintDetector' internal/risk/assessor.go` → 归普查腿/`309-v1`。
2. **`gofumpt -l cmd/wisp` 现量 2**（`models.go`/`panel_transport_35r2_test.go`）＝既有债：归**编排者面**决定另立与否，⛔ 本票修、⛔ 因它判红。
3. **窄变体 ⓐ′**（只对会用探测器的写类工具在空 id 时拒）＝票面登记的"未计价第三形"：⛔ 本轮采纳，要真走得先来一枚普查（必动文件/必红枚数/与 `internal/risk` 判级接缝），归未来程。
4. **`①` 生产构造点 `loop.go` 的 taskID 恒非空**：我按票面读为 `newTaskID()` 恒 36 字符、熵失败退化成时钟派生仍非空；本票⛔ 复核 `newTaskID()` 实现细节（那属票 242/154 族），尺＝`sed -n '1153,+20p' internal/agent/loop.go` → 归 `309-v1` 或普查腿。

## ④ 我跑过的编译面命令自报（⛔ 瞒）
- `go build ./...`：3 发（含一次因同批 `rm` 导出树污染 Go temp 的假 rc=1，已重跑 rc=0）。
- `go vet ./internal/tools/`：6 发（其中 1 发红＝`toolrequest_nail` 误 import `os` 未用，已删，非放宽断言）。
- `go test ./internal/tools/`：≈13 发（基线 2＋逐枚新用例多发＋突变台件 3＋枚数钉 Rig 2＋改后全量 1＋门禁全量 1）。
- `gofmt -l/-d`：≈5 发；`gofumpt -l`：4 发；`sh scripts/d22scan.sh`：3 发（1 红→抓到 `≈` 字面量→2 绿）。
- ⛔ 跑过任何会删/移别人读数的命令；临时件落 `probes/309/r1/logs/`；⛔ 证据件叫 `.out`（根 gitignore 第 8 行全仓 `*.out`）。唯一 `rm`＝我自用后 `rm` 掉了**我自己刚建的** `logs/headfmt/` 临时导出树（非任何 读数/别人件，两尺读数已落本文）。

## ⑤ ⛔ 顺手动到的东西（含射程边上有想法但没做）
- ⛔ 动 `:925` 布尔式、⛔ 删 `taskID == ""` 肢、⛔ 动 `b.injected`/`Options.Assessor` 肢（装配族另裁）——只在其上加/改注释。
- ⛔ 动 `internal/risk/**`（`scopeMarks`/`SrcUnboundScope`/`rules_gateway` 只读引用）。
- ⛔ 给 `agent.ToolRequest` 加字段；⛔ 改 `lit236*` 两枚常数值。
- `bridge.go:113/114/230/252/1084` 的其它绝对 `loop.go:N` 锚**可能也有漂**：有"顺手全扫"的想法，⛔ 越 AC#3 射程（只 `:861`/`:871` 是 AC#3 对象），未动。
- ⛔ 通过 `Bridge.Execute`/`factsFor` 端到端驱动 `AC#1`（会牵别的判定支）：改在包内直接驱动 `assessorFor+Assess`，⛔ 新造生产可达路径。
- ⛔ 修 `cmd/wisp` gofumpt 债；⛔ 翻票 309 任何框。

## ⑥ 改前⇄改后逐名作差结论（原始＝`logs/baseline-*.log` ⇄ `logs/after-1.log`；详见 `40-*.md`）
- 顶层 `--- PASS` 209→211（+2＝我两枚顶层用例）；子测试 75→75；`--- FAIL` 0→0；`--- SKIP` 0→0。名集合差**纯新增**（只两条 `> Test309…`，⛔ 删除线）⇒ **新增红 0 枚、⛔ 覆盖搬家式消失**。
- 两枚真喂空 id 进桥的 236 载体：`Test236R2TaskSpawnRefusesWhenHostGaveNoTaskID` ＋ `Test236R2TaskCancelRefusesWhenHostGaveNoCallerID` **均照绿 PASS**（Ⓑ 行为逐字节不变的直接佐证）。
- 突变三件齐：AC#1 Rig-A/Rig-B 均"改前绿→突变红→还原逐字节一致"；AC#2 Rig"加构造点即红→还原 SCRATCH_REMOVED"。三件读数在 `logs/mutation-rigs.log`/`logs/nail-rig.log`。

## `AC#5` 三本账对账（⛔ 同一件事记两遍）
本票病灶＝**桥 `assessorFor` 的 `taskID == ""` 肢静默摘 R4 探测器**（缺结构保证＋缺仪器＋四处同形处置不一致中的"静默"那一枚）。分工：
- **票 307**（已 `-done`）：折法不对称、成因⛔ 同 ⇒ 与本票⛔ 一件事。
- **票 40**：面板 Security 页零消费者，归口 `A838` §3。与 309 ⛔ 同路径；309 的Ⓑ注释专防票 40 那族"未来时"病（反向 grep 未来时字样命中 0，见 `10-*.md` 尺②）。
- **票 197**：`subagent_197.go:263` 拒绝形（"派生要有父"）＝四处同形之一（③）；本票⛔ 采此形落地，只在 `AC#0` 裁语与注释里具名对照。
- **票 242**：correlation id 改每调用一枚 ⇒ 正是 `bridge.go:920-923` 那句"per scope, hence per task"的**前因**；本票保留该句、只在其后追加 Ⓑ 说明。
- 台账 `A838` §3 那枚欠读数**已由 `A842` 结掉**，⛔ 在本票再记。
