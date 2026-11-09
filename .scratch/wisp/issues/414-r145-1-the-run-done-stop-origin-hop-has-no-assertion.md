# 414 — `R-145-1`：`run.done.stop` 的起源跳（loop→Event.Stop）全仓无断言

- **Status**：ready-for-agent（编排者 2026-11-09 22:5x 立；起手锚点＝本票现读 `git rev-parse --short HEAD`）
- 来路＝票 145 验收腿 `145-accept-r1` 登记的 `R-145-1`（**实测非推断**）；145 已于 `ff6b8a2c5` `-done`，本枚**不阻断**那张票，另立结清。
- 关联：票 145（面板快照载体）· 票 200（给载体装生产者 `a90677e92`）· `run.done.stop` 是 AC#6 六枚新键之一

## 现量的形状（票 145 验收腿原文，落地腿起手必须复跑）

- **缺口**：`run.done.stop` 的**起源跳**（`loop.go` → `Event.Stop`）**全仓无断言**。值本身能读到（AC#6 反射探针抽到该键），但**没人断言它来自 loop 那一跳**。
- **★实测证据（非推断）**：`acc-r1-mut`（只掐 `loop.go:1134` 的 `Stop: stop`）下，全量 `./cmd/wisp` **FAIL 名册与干净树逐名相等（63 vs 63）** ⇒ **把起源掐掉，零枚用例变红** ⇒ 现状确实是"装饰品"。
- **★已证明可断言**：验收腿明说「6 行就能补」。

## 要建什么

- [ ] **AC#1 补断言（6 行量级）**：写一条用例断言 `run.done.stop` 的值**来自 loop 那一跳**，而不是碰巧非 nil。
- [ ] **AC#2 反向判据（本票的灵魂）**：造突变「掐 `loop.go` 的 `Stop: stop`」⇒ **新用例必须变红**；还原 ⇒ 变绿。★这条直接复现验收腿的 `acc-r1-mut` 读数（改前 63 vs 63 逐名相等）。**两个读数都贴**。
- [ ] **AC#3 门禁**：gofumpt（绝对路径）空／`go vet` rc0／`go build` rc0／`go test ./cmd/wisp -run '<新用例名>'` 绿／`d22scan` clean／path-length GREEN。★跑 `./cmd/wisp` 整包必抬 `-timeout 30m`（默认 10 分钟会杀在 600.189s 造假红）。
- [ ] **AC#4 别把缺口换个地方**：⛔ 用"值非 nil"当断言（那就是现状的装饰品）；断言必须能区分"loop 给的"与"别处给的"。

## 禁区

- 写面＝`cmd/wisp/**`（**每波只许一枚写腿**）＋本票 Progress log。
- ⛔ 碰 `internal/panel`／`internal/config`／`internal/ball`／`frontend/**`／`design/**`。
- ⛔ 动 `thresholds.go`／golden／`PLAN.md`／`specs/**`／`.github/**`。
- ⛔ 改本票 AC 原文；⛔ push。

## Progress log (append-only, newest last)

- [2026-11-09 22:5x +08] 编排者立票（来路 `R-145-1`；缺口与 `acc-r1-mut` 读数逐字引；已证明 6 行可补）。