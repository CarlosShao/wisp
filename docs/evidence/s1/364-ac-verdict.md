# 364 AC 裁决（编排者＝非实现者，2026-11-04；锚 HEAD ae2ab606d）

## AC 逐格
- **AC#0 [x]**：五尺复跑（cmd+panel 非测试）全 0＋正控成立（`memory.`=21 命中、session 包已导入）＋用户看会话路径=无。件 `364-a0.md`。
- **AC#1 [x]**：裁形=〔甲〕CLI 出口（不触 C17/gate35；主流 harness 标配）。
- **AC#2 [x]**：五动作 list/search/rename/archive/unarchive 落地（r1a `1b4b607aa`／r1b `f98cbce41`／r1c `fb27a8843`／t `ae2ab606d`）。`cmd/wisp/sessions_test.go` 320 行：`TestSessionsListSearch`/`TestSessionsRenameArchiveUndo` **绿**（编排者亲跑 `go test ./cmd/wisp -run TestSessions -count=1` = ok 0.888s rc=0）；**反向判据三读数**（注 search→仅 ListSearch 红／注 unarchive→仅 RenameArchiveUndo 红／还原→双绿＋基线绿）。
- **AC#3**：代码门绿＝`go vet ./cmd/wisp` rc=0、`d22scan` clean rc=0、`go build ./cmd/wisp ./internal/session ./internal/memory` rc=0、`go test -run TestSessions` rc=0、`gofumpt -l` 四件空。**两红均非 364**：①`go build ./...` 红在 `build/mut_rvw.go`（`build/` 0 tracked＝untracked 暂存残留，干净树无此项）；②`path-budget` 红在 `273/403 的 -done.md 长名（我先前收口所致）＋ roster 漂移(62 vs 64, 128/175 落榜)`——364 未触及 issues/ 或 roster。
- **AC#4 [x]**：腿未翻框/未改票名 Status/未 push；编排者确认。

## ★系统性问题（提机主，非 364）
`check-path-length-budget.sh` 因**收口长名票的 `-done.md`**（273、403）在 self-hosted runner 上超 165 字符而红，且 roster 未随树更新（62 vs 64；128/175 落入"roster 有但树已不超"）。⇒ 根因＝「长票名 + `-done` 后缀 + 165 路径预算」三方撞。**Q**：①给 roster 补 273/403 并剔 128/175（维护），还是②长名票收口改用短别名（改约定）？属"维护数据 vs 改约定"，摆机主定，⛔ 我擅改 gate 数据。
