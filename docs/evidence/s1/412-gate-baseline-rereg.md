# 412 门禁基线重登 — 验收读数

## 1. 重登后门禁 rc=0（编排者亲跑，2026-11-09 12:5x，8m32s）
- `bash .scratch/wisp/probes/154/gate-clauses.sh` ⇒ **rc=0**
- `基线过期枚数＝0`
- `聚合退码＝0`（每腿的响与不响都和登记的声明一致，名册每腿在场）

## 2. 防作弊正控 rc=3（导出合成树，注入只开不关新形状）
- 导出树 `git clone --shared`，注入 `internal/tools/mutation_412_openonly.go`（OpenTask 无 CloseTask），commit `f0b0664ab`
- 同一门 ⇒ **rc=3**，逐腿 BAD：
| 腿 | 声明 | 基线 | 实测 | 因 |
|---|---|---|---|---|
| G2 | ring | 2 | 5 | 新增命中（实测>基线） |
| G6 | quiet | 0 | 1 | 空/非空与声明不符 |
| G6neg | ring | 1 | 2 | 新增未成对（实测>基线）＝**正是注入形状** |
- ⇒ **干净树 rc=0 ⇄ 注入树 rc=3**：门没放松，新违规照样被抓。

## 3. 改动面（`3bf4a4e60`，+48/−11）
- 只 `want/want_n` ＋ G1 一味 `-w`；一枚腿没删、一条射程没减、名册本体逐行原样。
- G1b quiet→ring 0→4／G2 1→2／G5 1→2／G5neg 8→14／G6neg quiet→ring piles 0→1／G7 1→6／G7neg 2→7
- G5pos quiet/0 不变＋排除 `cmd/wisp/resident_present_windows.go`