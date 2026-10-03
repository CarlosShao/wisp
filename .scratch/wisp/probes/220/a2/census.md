# 票 220 只读普查 a2 — 「`L1Windows` 这一跳今天断在哪、下一跳要动哪几枚文件」

## §0 起手（现跑）

- 采集起点时刻：`2026-10-03T10:33:43+08:00`（`date -Iseconds` 现跑）
- 起手锚点：`git log -1 --format=%H` = `e16e303773f59be81d2f02eee3cc638870aa1630`（分支 `dev`）
- **采集期间锚点前移**：第二次现跑＝`1289d8bbe0b90c2f58854c2b58e21d9f518ec62d`（`10:38:41 +0800`，`docs(probes/252-r2)` 骨架），
  本文所有行号在**两枚锚点之间被复量过一轮**，复量结果逐格写在 §6。
- 票面：`.scratch/wisp/issues/220-roster-blockedOnApproval-only-sees-L2-and-second-card-silently-drops.md`（全文读过）
- 前一枚同票普查：`.scratch/wisp/probes/220/c1/census.md`（锚 `6233dedc`，2026-09-29，**在 220-r1 之前**）；
  写腿交件 `.scratch/wisp/probes/220/r1/impl.md` §8（本腿就是去量它留下的那一跳）。
- 性质：**只读普查**。零编译零测试零写码；⛔ 本腿没跑过任何 Go 命令（闸门见 §7）。
- `frontend/**`／`design/**`：未读、未引、结论不转述其内容。所有 grep/find 均写显式搜索根
  （`cmd internal tools docs scripts`），无一根用 `.`。

## §1 脏面名册（起手 `git status --porcelain` 现抄）

起手时 `cmd/wisp` 确有一枚未提交改动，**采集期间它被 255-r2 自己提交了**，两枚读数都留在下面。

| 时刻 | 现跑读数 |
|---|---|
| 起手 `10:33:43` | `M cmd/wisp/config_reload.go` ＋ `?? cmd/wisp/restart_tier_keys_255r2_test.go`（全仓脏名册里 Go 产码只有这两枚，其余是 `.scratch/**` 证据件与 `design/**` 删除） |
| 复量 `10:39+` | `git status --porcelain -- cmd/wisp` ＝ **空**；两枚已由 `67ab595d`（`10:37:18 +0800`，`255-r2（编排者代提·死腿尾部那一格落定）`）入库，该 commit 的 stat 恰含 `cmd/wisp/config_reload.go` ＋8/−? 与 `cmd/wisp/restart_tier_keys_255r2_test.go` ＋136 |

⇒ **本腿读 `cmd/wisp` 时唯一需要标注的中间态窗口＝`10:33:43`–`10:37:18`**，那几分钟里 `config_reload.go` 是 255-r2 的工作树版本。
本腿对 `cmd/wisp` 的全部引用落在 `run.go`／`panel_pump.go`／`resident_*.go`／`panel_config_store.go` 与若干 `_test.go`，
**没有一枚落在 `config_reload.go`** ⇒ 行号不受那个窗口影响；`run.go` 在该窗口零改动（`git status` 空即可证）。

仍在飞的写腿（编排者令里点名的三枚，本腿一律未碰其写面）：`255-r2`（cmd/wisp＋internal/config）、
`171-r3`（门钉尺）、`252-r2`（internal/tools，其骨架件 `1289d8bb` 落在这几分钟）。

## §2 待验断言（编排者令里写下的每一句，全部按「先当假的」处理）

| # | 断言 | 本腿复量 |
|---|---|---|
| F1 | `ce693ce4` 交回 `window_read.go`（`L1Window` 三枚 string／零方法＋`(*Gate).LiveL1Windows()`） | §8 |
| F2 | `480b970d` 交回 `L1WindowWait`＋`PumpSources.L1Windows` 那一枚 reader 位 | §8 |
| F3 | 「`cmd/wisp` 里没人给 `PumpSources.L1Windows` 供数据」 | §8 |
| F4 | 「L1 窗口的开／关缺发布触发」 | §8 |
| F5 | 「`internal/panel/git_test.go` 一带有 `want-4` 那族锚，只认 `bridge.go` 里恰有四枚 `panel.*`」 | §8 |
| F6 | 「票 256：`approval.New` 建得太早 ⇒ 常驻腿读不到 `[risk]` 那几枚配置键」 | §8 |
| F7 | 编排者令里「`cmd/wisp` 此刻是脏面（255-r2 未提交）」 | §8（已被时间推翻，见 §1） |
| F8 | 票面「现量」那一表的行号（票自己写了「别信这里的行号」） | §8 |

## §3 供给链逐跳（问题 A）

本节记 `PumpSources` 每一枚 reader 位的**生产填充点**、`L1Windows` 缺的是哪一档，
以及 `approval.Gate.LiveL1Windows()` 到面板快照之间每一跳的三态判定。

## §4 发布触发（问题 B）

本节记 L1 窗口「开」与「关」各自的触发口在哪一枚文件、被谁调、以及触发那一刻窗口还在不在表里。

## §5 装配根现状与最小改动面（问题 C／D）

本节记 `NewPanelManager`／泵／`approval.Gate` 三者的建点与先后，
再把它落成「动哪几枚文件、每枚几行、要不要新开导出方法」的面。

## §6 邻居隔离规矩（问题 E）与量不到的格子

本节记 `cmd/wisp`／`internal/panel` 里与本跳相邻的词面型与计数型负向尺，逐枚给「落地腿怎么不误响」。

## §7 本腿闸门与自查

- ⛔ 未跑 `go test`／`go build`／`go vet`／`go list`／`go run`／`wisp`；未跑 `staticcheck`；未跑 `d22scan`。
  凡「要跑才知道」的格子一律进 §6 的「量不到」名册，不猜。
- ⛔ 未写、未还原、未 commit 任何别人的文件；本腿唯一写面＝本文件。
- ⛔ 未碰任何 `- [ ]` 勾选框；未改 `docs/specs/**`／`docs/PLAN.md`／`internal/observe/thresholds.go` 一字。
- 不 push；commit 带显式 pathspec；未用 `git add -A`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`。
