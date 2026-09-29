# 225 — **`DEFERRED(...)` 代码标记与 `SPEC-12 §5` 登记表今天对不上，而且没有任何仪器在查这一条**：`AGENTS.md` §1.1 那条"1:1 双向"目前是**只靠人**的规矩

- Status: **待派（只读审计腿即可，不产码）**。来源：票 223 的普查腿 `223-c1` 顺带撞出、编排者现跑复认（台账 `A426`）。
- ⚠ **为什么现在要管**：票 223 之所以差点被写成"我们漏做了热加载"，就是因为**代码里那枚标记说这件事"由票 05 实现"，而票 05 早就结案了、热加载在生产里根本不转**。标记在骗人 ⇒ 后面的程要么重复造、要么以为有人排过。

## 现量（起手逐条复算，别信这里的行号）

| 事实 | 读数（锚 `36b46125`，2026-09-29 11:57 现跑） | 尺 |
|---|---|---|
| 代码侧标记枚数 | `grep -rn "DEFERRED(" --include=*.go internal/ cmd/ tools/ \| grep -v _test` ＝**26 处命中／22 枚去重标记** | 现跑 |
| 登记表侧枚数 | `docs/specs/SPEC-12-roadmap-governance.md` §5 那张表：`^\| DEFERRED`＝**12 行**、`^\| RESERVED`＝**7 行**、`^\| REJECTED`＝**6 行**（另有 2 行是 `~~DEFERRED~~ → 已采纳`） | 现跑 |
| 差 | **22 vs 12 ⇒ 至少 10 枚代码标记在 §5 找不到同名行**（差数只是线索，**逐枚映射要人做**：有的标记是同一件事的第二种拼法、有的可能对应 §5 用了另一个名字） | 现跑 |
| 一枚确凿的失配 | `internal/config/doc.go:14` 逐字 `DEFERRED(schema/hot-reload/migration): implemented by ticket 05. This`，而 `.scratch/wisp/issues/05-config-model-done.md` **已结案**、`CheckAndReload` 的**唯一非测试调用者仍在 `cmd/balldebug/main.go:243`**（见票 223）⇒ **"标记声称已由 05 实现"与"生产里没人调"两句同时为真** | 现读 |
| 另一枚确凿的失配 | `internal/watchdog/doc.go:18` `DEFERRED(watchdog loop/thresholds): implemented by ticket 42`；而 `grep -n "看门狗\|watchdog" docs/specs/SPEC-12-roadmap-governance.md` 只命中 S6 那行**切片表**、**不在 §5 登记表** | 现读 |
| 没有仪器 | `grep -rln "SPEC-12" --include=*.go tools/ internal/ cmd/` 只命中 `internal/agent/approval/` 两枚文件里的注释；`tools/` 里没有任何程序读 `DEFERRED(` ⇒ **`AGENTS.md` §1.1 那句"必须 1:1 双向对得上"今天没有自动检查**（`tools/d22scan` 扫的是代码形状，不是登记表） | 现跑 |

## 判据（每格都要现跑读数）

- [ ] **AC#1 逐枚映射表**：22 枚去重标记逐枚列出（文件／行／标记名／它自称由哪枚票实现／那枚票今天什么状态），并各自判定三态之一：**§5 有对应行**／**§5 没有该补一行**／**标记本身过期了（自称已实现且确已实现，该摘掉）**。⛔ **不许改动 `docs/specs/**` 任何文字**——要补行／要摘标记都写成**待人批准的落点**（改契约＝人工批准；`DEFERRED(D-xx)` 与 §5 的 1:1 本身就是契约面）。
- [ ] **AC#2 反方向也走一遍**：§5 那 12 行 DEFERRED＋7 行 RESERVED 逐枚问"代码里有没有对应标记"，缺的具名列出。⚠ 方向二是本票的主要价值——**"登记了但代码里没标记"的项最容易被后面的人当成"没这回事"**。
- [ ] **AC#3 那两枚确凿失配要有结论**：`internal/config/doc.go:14`（热加载，见票 223）与 `internal/watchdog/doc.go:18`（看门狗循环）。各写一句"该改标记还是该补登记"，并具名指出改它要不要人批准。
- [ ] **AC#4 判"要不要上仪器"，不要自己上**：本票**只出结论与形状**——如果要把这条做成常驻检查（像 `tools/d22scan` 那样），那是新增一台仪器＝新射程，**要先摆给 owner**；不许本腿自作主张写进 CI。

## 禁区

不动 `PLAN.md`／`docs/specs/**`／`thresholds.go`／golden／`allowlist.txt`；不动三枚冻结件（`internal/panel/tokens_fourway_test.go`／`internal/panel/l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go`）；不新增导出名；`frontend/**`／`design/**` 零读零写零转述；⚠ **本票与票 223 同源但射程不同**：223 管"接线"，225 管"账目对得上"，**不许把两票的合成一票做**。
