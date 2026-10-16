
---

## 台面配方更正（编排者 12:5x 追加，只追加、⛔ 改你已读过的段落）

★**七件套不够，要八件套**：`go.mod go.sum internal cmd tools scripts frontend` ＋ **`models`**。
成因具名＝`internal/models/manifest.go:317 ResolveManifestPath` 在运行时读 `models/manifest.json`（跟踪件 2 枚：`manifest.json` 与 `manifest.json.minisig`）⇒ 导出树少这一枚就长出 `committed models/manifest.json not found` 那族红（今天 `internal/models` 三枚红正是它，⛔ 计入判据）。这条是我的清单写漏、⛔ 你的错；发现者＝在飞的 `220-v1`。

★另两枚同时更正（都在我给你的原文里，追加说明、不改原文）：
1. **取"改前"那一发的锚＝具名提交号，⛔ `HEAD~n`**（共享分支上别人正往中间落 commit，`HEAD~2` 会取到已含产码的树）。
2. **⛔ `go test ./...`（全树）**：那形会把 `cmd/wisp` 整包重跑一遍，不带 `-timeout` 时默认 10 分钟帽 `panic: test timed out after 10m0s`＝台面故障长得像"包红了"，⛔ 计入判据。
