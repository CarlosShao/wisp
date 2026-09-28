# 185-c1 — 只读普查：「第二次读回该由谁给这枚声明」三支候选的现量与代价

- 派单：`.scratch/wisp/dispatches/2026-09-28-142x-readonly-185-c1-who-should-declare-the-read-back-path-and-what-does-each-route-cost.md`
- 工单：`.scratch/wisp/issues/185-every-successful-reread-mints-the-blocker-for-the-next-one-because-the-fs-read-result-becomes-an-undeclared-c25-mark.md`
- 性质：**只读普查**。本表**不勾任何 AC 框**、**不派落地腿**、**不替编排者拍板选哪一支**。
- 本表所有行号皆可漂：凡引号内是**符号名 + `grep -n` 尺**，不是死号。

---

## ① 起手锚 / 写面闸门 / 三枚受保护文字 md5（起手向）

- 起手 `date`：待填
- 起手锚 `git log -1 --format='%h %cd'`：待填
- **本票的"未修码"＝待填**（票 183 的修法已落地，HEAD 是"修完之后"；本票缺陷在 HEAD 上仍在 ⇒ 未修码＝HEAD，须具名）
- 起手写面闸门 `git status --porcelain -- internal/ cmd/`：待填
- 三枚受保护文字 md5（起手向）：
  - `internal/risk/provenance.go` 的 `Mark` 契约文字块（尺 `grep -n "func (m \*Mark)\|^// "` 定位 468-473 那一段）＝待填（基线 `858e45116383caa3e7c1dd4b0924fad1`）
  - `internal/risk/taintmatch.go:11-15` 段＝待填（基线 `5680ddd18e2d2ec2a85e485b54f4c12e`）
  - `internal/risk/provenance.go` 的 `MarkWithHostPath` 文档块原 25 行＝待填（基线 `f89e891e5eee3f3ea2b4f89d921072c4`）
  - ⚠ `taintmatch.go` 第 100-109 行一带**不在受保护名单**（"仅参考"格，随号推移），本表不据它报警。

## ② AC#2(a) 桥侧：`fs.read` 成功时也填 `hostPathBox`

- 现量取号（尺：`grep -n "hostPathBox\|withHostPathBox" internal/tools/bridge.go internal/tools/task.go`）：待填
- 落点具名函数：待填 ／ 该函数今天有没有别的调用方：待填
- 会红哪几枚既有判据（逐枚点名）：待填
- 最贵的代价：待填
- 票 177 三边界（per-scope 名册／参数侧按值放行／`allowlist`）碰到哪条：待填
- 与票 183 AC#5① 明禁形状（"按参数值放行"）是否等价：待填

## ③ AC#2(b) D15 再落盘层：续读不产生第二枚 mark

- 那一层今天在哪（尺：`grep -n "func " internal/agent/spill.go`）：待填
- "不产生 mark"在码上等价于改哪一处：待填
- 落点具名文件＋现跑尺：待填
- 会红哪几枚既有判据（逐枚点名，尺 `grep -rln "C25\|Mark(" internal/agent/*_test.go internal/risk/*_test.go internal/tools/*_test.go`）：待填
- 票 175「成功结果必须有来源标记」承重声明会不会被拆：待填
- 最贵的代价：待填
- 票 177 三边界碰到哪条：待填

## ④ AC#2(c) 不修，把"同一条路径只许读一次"写成明说文案

- 文案落点具名文件：待填
- 那句话是不是冻结件（尺：`grep -n "全文见" internal/tools/task.go docs/PLAN.md`；`PLAN.md` 的 `task.output` 桩文案）：待填
- 若落点在冻结件里 ⇒ **人工批准项**：待填
- 会红哪几枚既有判据（逐枚点名）：待填
- 最贵的代价：待填
- 票 177 三边界碰到哪条：待填

## ⑤ 有没有第四支（只作用于这一发、不放宽任何一类外来正文）

- 码上是否存在第五条路：待填（没有就明写"没有"）
- 若存在：落点＋判据形状＋现跑尺：待填

## ⑥ AC#3 覆盖面名册（逐枚带尺，名册与枚数不许从票面抄）

- 尺（现跑）：`grep -rn "MarkWithHostPath\|hostPathBox" --include=*.go internal/ cmd/ | grep -v _test.go`
- 逐枚：工具 → 回执带不带宿主路径 → 会不会造出同款阻断者 → 今天带不带声明：待填
- 现跑全量输出留档：`.scratch/wisp/probes/185/c1/logs/`：待填

## ⑦ AC#7「同名不同目录」那一形的现状

- 比对面（尺：`grep -n "spellsDeclaredPath\|attachDeclaredPath" internal/risk/*.go | grep -v _test.go`）：待填
- 验了还是没验（具名）：待填
- 若验：变异名＋`-overlay` 路径＋红的那枚判据：待填

## ⑧ 本程没测什么（逐名，不写"其余都覆盖了"）

- 待填

## ⑨ 门禁终态（取在最后一枚 commit 之后）

- `sh scripts/d22scan.sh`：rc／`ban #8 internal/` examined 数（基线 433）＝待填
- `bash .scratch/wisp/probes/154/gate-clauses.sh` 红腿名册（不比退码；在册红腿基线＝`G6neg`＝票 178）：待填
- 逐包 `go test -count=1`：`./internal/risk/`（单跑）／`./internal/tools/`／`./internal/agent/`：待填
- CLI 那一面（`-overlay` + 在册 PATH 前缀）：待填
- 不带 PATH 前缀那一形（`0xc0000135`）：待填

## ⑩ 被拒／没成功的调用（逐条）

- 待填

## ⑪ 有没有跑过删除命令 + 工具调用终值自报

- 删除命令：待填
- 工具调用终值枚数（硬顶 40，28 枚停止新探索）：待填

## ⑫ next＝落地腿派之前还缺什么（含要不要摆人工批准项）

- 待填
