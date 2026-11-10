# 413 — `[ball] size` 有配置项、有范围校验，唯独没有生产读者：滑块落地要补的那一条喂线

- **Status**：ready-for-agent（编排者 2026-11-09 22:5x 立；起手锚点＝本票现读 `git rev-parse --short HEAD`）
- 来路＝机主裁定 `A1657②`（2026-11-09 21:5x 逐字：「**给滑块这个好，可以，可以搞个滑块，然后设置项自己调整，这个设计好**」）
- 关联：票 65（球样式族）· `A1636`（上限 72→160 经机主批准）· `D36`（配置三档生效级别）· 前端配置页 `config-screen.tsx:105`（归**外部 agent**）

## 现量的形状（2026-11-09 22:5x，编排者现跑；落地腿起手必须逐条复跑）

| 层 | 现状 | 读数 |
|---|---|---|
| 配置 schema | **早有** | `internal/config/schema.go:178` `BallSection.Size int \`toml:"size" default:"56"\``，注释：「range 44-160 (SPEC-03 sec 3; ceiling widened 72 -> 160 with owner approval A1636, default stays 56)」 |
| 范围校验 | **早有** | `internal/config/validate.go:64` `if c.Ball.Size < 44 \|\| c.Ball.Size > 160`，错误串 `config.toml: ball.size %d out of range [44, 160]` |
| Ball 内部 option | **早有** | `internal/ball/ball_windows.go:64` `SizePx int // configured orb size 44..160 (0 = default 56)`；夹取 `:155-163`（0→56／<44→44／>160→160） |
| **生产读者** | **缺** | `cmd/wisp/resident_ball_windows.go:396` `ball.New(ball.Options{...})` **不传 `SizePx`** ⇒ 永远走默认 56，配置改了没人读 |
| 前端配置页 | **画了但标 blocked** | `frontend/src/components/config-screen.tsx:105` `{ group: "ball", name: "size", desc: "悬浮球直径（44-72）", blocked: "那颗球是原生绘制的（D29），页面改不到" }` —— ★两处过期：范围写旧 `44-72`（现 44-160）；`blocked` 那句是**假前提**（不是"改不到"，是缺一条读+写回路由） |

## 要建什么

- [x] **AC#1 喂线（本票的核心）**：`cmd/wisp` 装配根把 `cfg.Ball.Size` 喂进 `ball.Options.SizePx`。★**生效档位按 `D36` 三档判**：size 是窗口几何（`WindowEdgePx` 依赖它建窗）⇒ **大概率重启档**，落地腿必须**具名写清是热改还是重启生效**，⛔ 含糊。
- [x] **AC#2 Go 侧读写接口**：`internal/config` 出 `BallSizeView`，照 `settings_65.go` 的 `BallStyleView` **同形**（读 + 写 + 校验 + 越界拒绝）。⛔ 只出读者不出写者（那就是半条路）。
- [x] **AC#3 反向判据（防装饰品）**：造突变「把 `SizePx` 喂线删掉」⇒ **必须有测试变红**；还原 ⇒ 变绿。**两个读数都贴**。★这条是本票唯一能证明"喂线真被见证"的判据。
- [x] **AC#4 往返测试**：配置写 90 ⇒ 读回 90 ⇒ 球建窗用 90；写 43／161 ⇒ 越界拒绝并给出 `validate.go` 那句原文。
- [x] **AC#5 门禁**：gofumpt（绝对路径 `/d/work/base/gopath/bin/gofumpt.exe`）空／`go vet` rc0／`go build` rc0／`go test ./internal/config ./internal/ball` 绿／`d22scan` clean／`bash scripts/check-path-length-budget.sh` GREEN。
- [x] **AC#6 前端契约（⛔ 不写前端码）**：把前端要的那条读+写接口**具名写进证据件**（方法名／入参／返回／错误码），交**外部 agent** 接线。★本票写面**不含** `frontend/**`。

## 禁区

- 写面＝`cmd/wisp/**` ＋ `internal/config/**` ＋ `internal/ball/**`（**同一时刻只许一枚写腿**，这三个包是共享编译面）。
- ⛔ 碰 `frontend/**`／`design/**`（机主外部 agent）。
- ⛔ 动 `thresholds.go`／golden／`allowlist.txt`／`PLAN.md`／`specs/**`／`.github/**`。
- ⛔ 改本票 AC 原文（只许翻勾＋追加 Progress log）。
- ⛔ push（子代理只 commit）。

## Progress log (append-only, newest last)

- [2026-11-09 22:5x +08] 编排者立票（来路 `A1657②`；五层现状逐条现量；机主已拍"做滑块"）。
- [2026-11-10 +08] AC#1 落地腿（commit c8dbd2288）：`cmd/wisp/ball_size_413_windows.go` 种子（ballSizeSeed413 读 c.Ball.Size，无字面量无自夹取）→ `withBallSizePx` hook → `resident_ball_windows.go:423 SizePx: hs.sizePx`；反射正/反测试（fed=88/unfed=56）。AC#2/4/6 未动。
- [2026-11-10 19:2x +08] 补齐腿：**AC#2**＝`internal/config/settings_413.go`（`SetBallSize` 守卫写＋`BallSizeOptions` reader＋`BallSizeView{Size,Min,Max}`，照 65 同形；band 常量与拒绝句单源进 `validate.go` 的 `ballSizeMinPx/MaxPx/ballSizeRangeError`，句子原文不变——写闸与加载闸共用一句，永不分叉）。**AC#4**＝config 侧三测（写 90→收据 ball.size→盘上 size=90→读回 90；43/161/0 拒绝且 Detail 逐字＝validate 句、文件字节不变）＋cmd/wisp 全程往返两测（写 90→读回 90→种子 90→live boot 球持 90，球窗真起非 skip；43/161 拒绝贴原文）。**AC#3 反向两读数**：①删喂线（摘 `SizePx: hs.sizePx`）⇒ `Test413FeedReachesBallOptionsSizePx` FAIL（SizePx=56 want 88），还原 PASS 0.91s；②删写校验（摘 SetBallSize 守卫）⇒ `TestSetBallSizeRefusesOutsideRange` FAIL（拒绝句退化成整验包装句；writeOneKey 的 validate 闸仍拒写——守卫非唯一防线，但具名句契约破即红），还原 PASS。**AC#6**＝契约证据件 `docs/evidence/s1/413-ball-size-slider-contract.md`（读 `BallSizeOptions`/写 `SetBallSize`/错误句逐字/档位具名＝内存热+窗口重启/快照键 ballSize 与 config.set ball_size 的接线点具名交外部 agent）。**门禁**：gofumpt 空／vet rc0／build rc0（三包）／`go test ./internal/config ./internal/ball` 绿（2.661s/0.690s）／`go test ./cmd/wisp -run Test413 -timeout 30m` 绿 3.198s 五枚 PASS／d22scan clean／path-length GREEN。**没做到（诚实账）**：面板两条路由未接线（composer/pump/carrier/roster 四件在 `internal/panel/**`，超本腿写面——今天页面拿不到 ballSize 快照键、config.set ball_size 仍被门拒）；CLI 未扩 ball.size 键（config_cmd_65 名册闭合，留与面板接线同片）；config.set 现量证据＝internal/panel/config_write_65.go:39-46 名册只有 ball_style、panel_config_store.go ApplySetting 无 ball.size case ⇒ 写面不存在，故另造 setter 合规。前端 config-screen.tsx:105 两处过期归外部 agent。