
## 编排者更正（2026-10-11 13:1x，来路＝只读普查腿 `writeface-1` 的逐腿文件名册与两条待裁项，件 `.scratch/wisp/probes/orch/20261011/writeface-1/**`；台账 `A846`；⛔ 上面各节原句一字不改，本节追加）

本票**⛔ 翻任何格**（六格全未勾，等 `295-r1` 与非实现者验收）；本节只做两件事＝**改掉一枚过期理由**＋**答一枚待裁项**。

- **① `:36` 那句"⛔ 不与票 294 同批（它可能动同一枚文件的 `executed` 那一支）"里的理由⛔ 成立**（腿当面顶回，我现跑复认）。尺＝`git grep -n 'executed' HEAD -- cmd/wisp`＝**19 命中**；`executed` 与 `started` **同现于一枚行＝0**；`executed` 的**唯一断言**逐字＝`resident_mute_290_windows_test.go:197` 的 `if executed {`；而 `294-r1` 的**入库**写面⛔ 含 `cmd/wisp/resident_audio_windows.go`。⇒ **两件事⛔ 碰同一枚字节**。★那条"⛔ 与票 294 同批"作为**排面**我维持（真约束＝两枚都要对同一段 `toggleMute` 做瞬态变异，锚位会互相打死；形如"盘上裸改"⛔ 是唯一走法——`-overlay` 与仓外备份＋`sha256sum` 还原两枚代偿法都在本仓有先例）。撤销口令＝**「294 295 改并发」**（与票 294 追加节同一枚，两票一起生效）。
- **② 待裁项 U-5 已答，判据⛔ 改**：`AC#4` 逐字要 `go test ./cmd/wisp/ ./internal/audio/ -count=1 -v` 改前改后各 ≥2 发。★现量＝`internal/audio` 在 `295-r1` 的写面**之外**（尺＝`git grep -cE '\bstarted\b' HEAD -- internal/audio`＝**gate.go 3／gate_test.go 2／wasapimic_windows.go 8／wavinjector.go 1**＝14 命中，**全是该包自己的 `started`**，与 `cmd/wisp` 那枚同名⛔ 同物）⇒ 第二包那一发买的⛔ 是行为回归，而是**"签名级回归"**（`cmd/wisp` 若打破 `internal/audio` 的 API，那一发才会响）。★**裁＝票面 `AC#4` 原样维持**（省机时⛔ 值到动判据；⛔ 由写腿自删票面判据）。撤销口令＝**「295 门禁去掉 audio 那包」**。
- ③ 腿给的两枚工程别扭处我认下并转进派单：`AC#4` 的三枚 `cmd/wisp` 腿共用同一枚包名册 ⇒ **每波只交一枚 `./cmd/wisp/` 色**（这条与票 294 追加节那三条同源）；名册尺⛔ 相加当一枚数（顶层⇄子项两把，本仓实测差 110 枚那族）。

- [2026-10-11 13:1x] agent=编排者 did=收 `writeface-1` 对本票的两条⇒ **`:36` 那半句过期理由具名更正**（`executed`⇄`started` 同现一枚行＝0、`294-r1` 入库写面⛔ 含 `resident_audio_windows.go` ⇒ ⛔ 同写一枚文件；"⛔ 同批"作排面维持，口令「294 295 改并发」）＋**U-5 已答⛔ 改判据**（`internal/audio` 14 命中全是它自己的 `started`、第二包买的＝签名级回归而⛔ 行为回归 ⇒ `AC#4` 原样维持，口令「295 门禁去掉 audio 那包」）；★本票六格⛔ 动、⛔ 改原句、⛔ push
