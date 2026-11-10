# 422 — `[ball] size` 滑块的**面板接线**：快照键 `ballSize` 与 `config.set ball_size` 路由未接，页面拿不到读口/写被门拒

- **Status**：ready-for-agent（编排者 2026-11-10 19:2x 立，从票 413 交付腿的诚实边界反推）
- 来路：票 413 落地（`c8dbd2288` 喂线＋`baefa152b` 写接口/往返/契约证据）的**诚实边界**——413 是后端 + 契约证据（AC#6 明说"交外部 agent 接线"），但**面板两条路由**在 `internal/panel/**`（413 写面外）未接。
- 关联：票 413（契约件 `docs/evidence/s1/413-ball-size-slider-contract.md`）· 票 65（ball_style 的 composer/pump/carrier + config.set roster 链，照它复刻）· 票 253（inbound roster）· 票 255（config_readers roster）

## 现量（⛔ 引用前先重跑）

- 413 已建：读口 `BallSizeOptions`/`BallSizeView`、写口 `SetBallSize`（`internal/config/settings_413.go`）。
- **未接**：`grep -rn "ballSize\|ball_size" internal/panel/ cmd/wisp/ | grep -v _test` 应见零（除 413 写面）；composer Snapshot 无 `ballSize` 键；`config.set` 名册（`internal/panel/config_write_65.go:39-46`）只有 `ball_style`，无 `ball_size` case ⇒ 页面写 ball.size 被 `unlisted-field` 门拒。
- 照 65 链：`ball_style` 的 composer/pump/carrier + config.set roster 四件（读 65 那几张票/码定位）。

## AC

- [ ] **AC#1 快照读键**：`Snapshot` 增 `ballSize` 键（值＝`BallSizeOptions().Size`），composer/pump/carrier 三处照 `ballStyle` 同形扩；判据＝快照 JSON 含 `ballSize` 且值==config 现值。
- [ ] **AC#2 写路由**：`config.set ball_size` 进 `config_write_65.go` 名册 + `panel_config_store.go` ApplySetting 加 `ball.size` case（调 `SetBallSize`，越界拒绝透出 validate 句）；判据＝经 config.set 写 90 读到 90、写 43/161 拒绝贴原文。
- [ ] **AC#3 反向判据**：删 ballSize 快照键 ⇒ 读判据红；删 ball_size 写 case ⇒ 写判据红。两读数贴。
- [ ] **AC#4 门禁**：gofumpt 绝对空／go vet ./internal/panel ./cmd/wisp rc0／go build 同包 rc0／go test ./internal/panel -count=1 绿 ＋ go test ./cmd/wisp -run '<新用例>' -count=1 -timeout 30m 绿／d22scan clean／path-length GREEN。
- [ ] **AC#5 禁碰**：⛔ frontend/**（外部 agent）；⛔ thresholds/golden/allowlist；⛔ 改 413 契约件原文。

## 撤销口令
- 「**422 撤回**」。