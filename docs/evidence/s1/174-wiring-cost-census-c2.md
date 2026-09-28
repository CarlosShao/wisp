# 174-c2（只读代价普查·零产码）——"把路径判定者接到生产"这一跳要动哪几层、会把谁顶出去

- 程：`174-c2`｜性质：**只读**（AC 框一枚不许勾、产码零字节不动、不落地、不选方案）
- 派单：`.scratch/wisp/dispatches/2026-09-28-145x-readonly-174-c2-wiring-cost-census.md`
- 工单：`.scratch/wisp/issues/174-...-spill-artifacts-are-outside-fs-allowed-dirs.md`
- 骨架落盘时刻：`2026-09-28 14:42+0800`（现跑 `date '+%Y-%m-%d %H:%M%z'`）

> 这份表**只报代价与形状，不替编排者决定接不接、什么时候接**。
> 每一节里凡是"本程没跑"的，一律明写"本程未验"，不拿推理充数。

---

## 1. 起手锚＋写面闸门

（待填：起手 HEAD／分支／两条 `git status --porcelain -- internal/ cmd/`／写面清单）

## 2. Q1 接线那一跳的完整形状（逐环：谁构造·谁持有·谁调用）

（待填：逐环 `grep -n` 尺＋文件:号；要动几枚文件；每动一处哪几枚既有判据会红，**逐枚点名**）

## 3. Q2 顶出去效应（接上前／后两发对照）

（待填：判定者语义现读；只读对照台件 `probes/174/c2/`，走 `-overlay`；两向读数逐字；收益与退让一起报）

## 4. Q3 与票 183／185 那根管子的关系

（待填：接上之后 `task.output` 回执那条产物路径是否自动读得回；给现跑或明写"本程未验"）

## 5. AC#2d 那一支今天到底有没有判据钉着

（待填：现量尺与逐字读数；没有就明说没有）

## 6. AC#3 三条禁区自证（只读能答的部分）

（待填：(i) 受门控 Tool／(ii) `risk.PathResolver` 之外的 `filepath.Clean|Abs`／(iii) `artifacts` 进 `allowed_dirs` 默认值——逐条"我这程没碰"的证据：`git status` ＋ `git diff --numstat` 删除列逐枚 0）

## 7. 本程没测什么（逐名）

（待填）

## 8. 门禁终态

（待填：`sh scripts/d22scan.sh`（基线 ban #8 internal/ examined=433）／`probes/154/gate-clauses.sh` 红腿名册（在册只有 `G6neg`）／`go test -count=1 ./internal/tools/`；取在最后一枚 commit 之后）

## 9. 被拒调用＋零删除自证＋工具调用终值

（待填）

## 10. next：落地腿派之前还缺什么（含要不要人先批准哪一格）

（待填）
