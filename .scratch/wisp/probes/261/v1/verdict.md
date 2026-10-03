# 票 261 · 腿 261-v1 裁决表（非实现者验收：攻 261-r1 修法的门有没有牙）

代理：`261-v1`（**非实现者**验收子代理；实现者＝`261-r1`，其产码三枚 commit `2bbe7086`→`e2bfd7c8`→`9d0cd518` 已入库）。
票面：`.scratch/wisp/issues/261-the-model-enabled-key-promises-removal-from-discovery-and-selection-but-no-production-code-reads-it-while-the-default-true-tag-is-inert-for-map-entries.md`
实现件：`.scratch/wisp/probes/261/r1/impl.md`；前程件：`.scratch/wisp/probes/261/p1/verdict.md`。
本腿定位：**攻防**——不复述 r1/p1 已留档读数当自己的凭据；变异复认两形各至少 1 发自跑、恒真两问、三条入口现跑、生产调用者问句、旁格核对、AC 格判语。

---

## §0 起手锚与起手名册

- 起手时刻：`2026-10-03T16:25:51+08:00`（`date -Iseconds` 自取）
- 起手 HEAD：`89c863f168cf27c76bfae165666f34a5dcd3fea9`（`git log -1 --format=%H` 自取），分支 `dev`
  （即 r1 三枚 commit 之后的锚点：`9d0cd518` 在 HEAD 历史内）
- 起手 `git status --porcelain internal/llm internal/config`：**空**（rc=0，零行）——两包起点干净，无他人脏面
- 起手 md5（跑 overlay 前 pre 记录，跑后逐枚复归核对）：
  - `internal/llm/resolver.go` = `e8a2cbc85e2393e6b9ca9c8db849dc1b`
  - `internal/llm/enabled_reach_261_test.go` = `dede3bb0863830bed0aa5f0e1d89e6de`
  - `internal/llm/enabled_gate_261_r1_test.go` = `3631fe757cba6ce7c00083df1683d580`
  - `internal/config/enabled_261_test.go` = `f2eef15374a39ecdd747e74298f72b2c`
- 起手绿名册：`go test -count=1 -v ./internal/llm/ ./internal/config/`（背景任务 b3b8phyjt）→
  名册抽取与 comm 对账在 §4 收尾处给（剥时延后缀后 `comm` 双向，零丢名才收）。
  （读数待填于此行上方，占位见 §4。）

## §1 变异复认（两形各至少 1 发自跑，红句逐字）

（待填）

## §2 恒真两问

（待填）

## §3 三条入口现跑读数

（待填）

## §4 生产调用者名册＋「行为变没变」判语

（待填）

## §5 旁格核对（三处逐字）

（待填）

## §6 AC 格判语（AC#0ⓐ/AC#0ⓑ/AC#1/AC#2）

（待填）

## §7 推翻清单

（待填）

## §8 判不动／量不到

（待填）
