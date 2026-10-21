# 359 — `sameAbsoluteness` 那条合取在 POSIX 侧**零分母**（注释性结论没有任何用例钉住）

- Status: **ready-for-agent**（`AC#0` 只读现量可先派；`AC#1` 动码须先过裁形）
- 立案：2026-10-21 19:3x +08 编排者立（从票 124 的 `AC#6` 拆出——★那格与本票三批工作**零交集**，⛔ 混为一格）。
- 来源：票 124 `AC#6` 原文（09-23 14:0x 编排者追加；来路＝票 129 验收方 `acceptor-ticket129-r1` 的 `R-129-2` + 实现者 `next=` N2）。
- 关联：票 124（三批已交：`internal/config` 33→0／`internal/tools` 38→3／`cmd/wisp` 37→0）· 票 129（`-done`）· `internal/winsec/resolve.go`

## 这是什么（人话）

`internal/winsec` 里有一条判定叫 `sameAbsoluteness`（判断两个路径的"绝对性"是否一致），它作为合取项参与某个决策。
`resolve.go:513-517` 有一段**注释**写着"这一形在非 Windows 上会怎样"——但今天：

- 在容器（POSIX）里跑 `./internal/winsec/` 得 **104/60/0/0**（全绿）；
- 而 `absoluteness` 这个词在 **POSIX 用例里命中 0 枚**。

⇒ **那段注释的结论没有任何用例当分母**：把 `sameAbsoluteness` 那条合取**摘掉**，POSIX 侧**零枚红**。
换句话说，那条合取在 POSIX 上**今天可以被删掉而无人察觉**——注释说它重要，但没有一枚用例证明它重要。

## 结案判据

- [ ] **AC#0 现读＋归因（只读，⛔ 动任何产码字节）**：把这件事量成可复跑读数——
      ①`sameAbsoluteness` 的**定义与全部调用点**逐名列出（⛔ 只给计数）；②`internal/winsec/resolve.go:513-517` 那段注释**逐字引出**；
      ③**变异**：摘掉那条合取 ⇒ 在 POSIX 容器里跑 `./internal/winsec/`，**必须复现"零枚红"**（★先证变异落地：`grep -n` 原文 ＋ `go build` rc=0）；
      ④**正控**：同一发变异在 Windows 侧（或任何**有分母**的包）必须在**某处**体现出差异，否则说明变异根本没生效。
      ★**若 ④ 也拿不到差异** ⇒ 停手上报（那意味着"这条合取今天在哪里都不起作用"，与本票前提不同）。
- [ ] **AC#1 裁形（编排者拍）**：三形各写代价＝〔甲〕给 POSIX 侧**补一枚用例**钉住该合取（最小侵入）；
      〔乙〕把那段注释**改成实话**（"这一形在非 Windows 上未有分母"），⛔ 补用例；〔丙〕两者都做。
      ★**未裁形任何实现都⛔ 开工**。
- [ ] **AC#2 落地＋牙（写腿）**：按裁定那一形做；每处改动配**一枚今天响、不修不响**的用例；
      ★反向判据＝把新加的那一跳摘掉 ⇒ 用例必须转红并**逐名报红行**；三读数（变异红／还原绿／原样绿）逐字贴 ＋ **还原核 `sha256`**。
- [ ] **AC#3 门禁**：`GOFLAGS= go build ./...`／`go vet ./internal/winsec/`／`go test -count=1 ./internal/winsec/`（四数分开报、⛔ 相减）／
      `sh scripts/d22scan.sh`／`sh scripts/check-path-length-budget.sh` 各 rc。★`CGO_ENABLED=0` 本仓必然编不过（sherpa），⛔ 当缺陷。
- [ ] **AC#4 纪律格**：⛔ 翻任何框（本票各格由非实现者裁）、⛔ 改票名／`Status:`／`-done`、⛔ push、⛔ 碰票 124／129 任何原句。

## 硬约束

- ⛔ 改 `internal/winsec` 的**判定语义**（`sameAbsoluteness` 的合取本身）；本票只**补分母**或**改注释**，⛔ 改行为。
- ⛔ 改 `docs/PLAN.md`／`docs/specs/**`（契约文字须人工批准、另步）。
- ⛔ `internal/risk/**`、`tools/d22scan/**`、`allowlist.txt`、`thresholds.go`、任何 golden、`frontend/**`、`design/**`。
- ★**POSIX 与 Windows 两侧都要报**：该合取在两边语义不同，只在一边量＝分母不完整。

## 排程

`359-a1`（只读＝`AC#0` 现读＋变异复现）→ 编排者裁 `AC#1` → `359-r1`（写腿＝`AC#2`）→ `359-v1`（非实现者裁）。

**Status 追加（2026-10-21 19:3x）**：号位＝编排者定（`358` 已占、`359` 现量空号）。⛔ 零翻框、零 push。

## Progress log

（追加式；每条 `- [日期 时刻 +08] agent=… did=… next=…`，钟点由 `date` 插值，不手打。）

- [2026-10-21 19:3x +08] agent=`编排者` did=立案（从票 124 `AC#6` 拆出；★那格与票 124 三批工作零交集，⛔ 混为一格＝`A1099`② 那族"把两件事混成一格"）next=派 `359-a1` 只读腿出 `AC#0` 现读＋变异复现。
