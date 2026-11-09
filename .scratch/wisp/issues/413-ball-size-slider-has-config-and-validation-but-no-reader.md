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

- [ ] **AC#1 喂线（本票的核心）**：`cmd/wisp` 装配根把 `cfg.Ball.Size` 喂进 `ball.Options.SizePx`。★**生效档位按 `D36` 三档判**：size 是窗口几何（`WindowEdgePx` 依赖它建窗）⇒ **大概率重启档**，落地腿必须**具名写清是热改还是重启生效**，⛔ 含糊。
- [ ] **AC#2 Go 侧读写接口**：`internal/config` 出 `BallSizeView`，照 `settings_65.go` 的 `BallStyleView` **同形**（读 + 写 + 校验 + 越界拒绝）。⛔ 只出读者不出写者（那就是半条路）。
- [ ] **AC#3 反向判据（防装饰品）**：造突变「把 `SizePx` 喂线删掉」⇒ **必须有测试变红**；还原 ⇒ 变绿。**两个读数都贴**。★这条是本票唯一能证明"喂线真被见证"的判据。
- [ ] **AC#4 往返测试**：配置写 90 ⇒ 读回 90 ⇒ 球建窗用 90；写 43／161 ⇒ 越界拒绝并给出 `validate.go` 那句原文。
- [ ] **AC#5 门禁**：gofumpt（绝对路径 `/d/work/base/gopath/bin/gofumpt.exe`）空／`go vet` rc0／`go build` rc0／`go test ./internal/config ./internal/ball` 绿／`d22scan` clean／`bash scripts/check-path-length-budget.sh` GREEN。
- [ ] **AC#6 前端契约（⛔ 不写前端码）**：把前端要的那条读+写接口**具名写进证据件**（方法名／入参／返回／错误码），交**外部 agent** 接线。★本票写面**不含** `frontend/**`。

## 禁区

- 写面＝`cmd/wisp/**` ＋ `internal/config/**` ＋ `internal/ball/**`（**同一时刻只许一枚写腿**，这三个包是共享编译面）。
- ⛔ 碰 `frontend/**`／`design/**`（机主外部 agent）。
- ⛔ 动 `thresholds.go`／golden／`allowlist.txt`／`PLAN.md`／`specs/**`／`.github/**`。
- ⛔ 改本票 AC 原文（只许翻勾＋追加 Progress log）。
- ⛔ push（子代理只 commit）。

## Progress log (append-only, newest last)

- [2026-11-09 22:5x +08] 编排者立票（来路 `A1657②`；五层现状逐条现量；机主已拍"做滑块"）。