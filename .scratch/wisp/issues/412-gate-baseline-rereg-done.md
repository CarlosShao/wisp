# 412 门禁基线重登＋G1 加 -w（票 411 裁定的机械落地）— done

**Status:** done（2026-11-09 13:0x；执行腿改＋编排者验证结案）

## 背景
票 411 修好 `gate-clauses.sh` 的命令行溢出后，这把尺**第一次对全树跑出真读数 rc=9**（9 腿真红）。票 411 AC#4 裁定：**甲（真缺陷）=0／丙=G1（pattern 不精确，加 `-w`）／乙=7 腿（合法新码，重登基线）**。本票＝把裁定机械落地。

## 做了什么（`3bf4a4e60`，仅 `want/want_n` ＋ G1 一味 `-w`）
- **G1＝丙**：词表加 `-w`（旧 `-E` 子串 `ToolRequest\{` 误命中 plugin 层另一类型 `Tier2ToolRequest{`，`tier2.go:46`，定义 `tier2_loader.go:80`）。
- **乙 7 腿重登**（每枚新基线严格＝实测真读数）：G1b quiet→ring 0→4／G2 1→2／G5 1→2／G5neg 8→14／G6neg quiet→ring piles 0→1／G7 1→6／G7neg 2→7。
- **G5pos＝正控**：quiet/0 不变，正控排除名单同格补 `cmd/wisp/resident_present_windows.go`（该文件进程生命周期 scope 与在册活形状同形，正控不能因它响）。★裁定表主值「0→2」与 facts §4/§9 实测 1 相抵，按括号条件执行。
- ★**一枚腿没删、一条射程没减、名册本体逐行原样**。

## 验收（编排者亲跑）
- **`bash .scratch/wisp/probes/154/gate-clauses.sh` ⇒ rc=0**，**基线过期枚数=0**（每腿与自己登记的声明一致）。
- **★防作弊正控（导出合成树 `git clone --shared`，注入一枚只开不关新形状 `internal/tools/mutation_412_openonly.go`）⇒ rc=3**，逐腿 BAD 具名：
  - `BAD 腿=G2 声明=ring 基线=2枚 实测=5枚`（新增命中）
  - `BAD 腿=G6 声明=quiet 基线=0枚 实测=1枚`（空/非空不符）
  - `BAD 腿=G6neg 声明=ring 基线=1枚 实测=2枚`（**新增未成对＝正是注入的形状**）
  - ⇒ **新违规照样被抓，门没放松。**
- 读数存档：`/d/tmp/wisp/20261109/gate-after.log`（rc=0）＋`gate-mut.log`（rc=3 正控）。

## AC
- [x] **AC#1 重登后 rc=0**：编排者亲跑 rc=0、基线过期 0。
- [x] **AC#2 判据语义不变**：只重登 want/want_n＋G1 加 -w，⛔ 删腿⛔ 减射程，名册本体逐行原样。
- [x] **AC#3 防作弊正控**：注入新只开形状 ⇒ rc=3，G2/G6/G6neg 具名 BAD；干净树 rc=0。
- [x] **AC#4 裁决者≠实现者**：执行腿 `411-baseline-rereg` 改（超时未提交）＋编排者核 diff 提交＋编排者亲跑验收。
- [x] **AC#5 每枚新基线＝实测真读数**：diff 逐枚对照 facts §0/§9 九腿对平。

## 撤销口令
`git revert 3bf4a4e60`（退回 9 腿 rc=9 真读数态）。