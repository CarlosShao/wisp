# 253-r1 · 起手锚（00-anchor）

- 腿：`253-r1`／时刻 `2026-10-08 14:47 +0800`（`date '+%Y-%m-%d %H:%M %z'` 现跑）
- 锚点：分支 `dev`，HEAD `6547fd30`（`git log --oneline -1`；该笔正文 §4 逐字把本腿列进本波，出处＝台账 `A717` §5 收窄）
- CWD＝仓库根 `D:\work\workspace\projects plans\Wisp`（含空格，命令一律加引号）

## 起手锚读数

1. `git status --porcelain -- cmd/wisp` ＝ **空**（rc=0）⇒ cmd/wisp 无别人的在飞改动，可起手。
2. 票 253＝`.scratch/wisp/issues/253-panel-inbound-three-ruler-holes.md`，`wc -l` ＝ **53**；两把框尺现量：未勾 `- [ ]` ＝ **4**、已勾 `- [x]` ＝ **0**（改前后各复跑一次，见 99-final）。
3. 槽位复量＝`ls .scratch/wisp/probes/253/` → `p4 / r1 / r3 / r5`。`r1` **目录已存在但为空**（`ls -la r1/` 仅 `.`/`..`，mtime `Oct 2 16:54`，无 `00-anchor.md` 无任何件）。`r1` 即编排者派给本腿的槽位（`A717` §5／`A718` §4 具名 `253-r1`）⇒ **本件为该空槽的首枚**，不撞别人已交的件（非"同槽覆盖"）。⚠ 该空目录早于今日派单存在，按薄索引规矩只具名不擅动其它槽（`r3`=precheck.md 属 253-r3、`r5`=inbound_roster 属 253-r5/AC#2、`p4`=253-p4）。

## 射程复跑（我自己 grep -rn，不照抄普查）

尺＝`grep -rn "panelDispatchBinding" --include=*.go cmd internal` ＋ `grep -rnE '\.Bind\(' cmd/wisp/*.go | grep -v _test.go`。

- 产码命中（`panelDispatchBinding`）＝**只有三处**，与普查 `probes/ruler-dedup-1/ruler-dedup-1.md:24` R5 行号**逐一对上、零漂**：
  - `cmd/wisp/panel_host_windows.go:80`　`panelDispatchBinding = "wispDispatch"`（名册枚＝常量定义）
  - `cmd/wisp/panel_host_windows.go:792`　`})();`, panelDispatchBinding)`（JS 转发半，`panelPostMessageForwardInit` 的 `fmt.Sprintf` 实参 → 页面侧 `window.wispDispatch(...)`，`:781–792`）
  - `cmd/wisp/panel_host_windows.go:802`　`if err := w.Bind(panelDispatchBinding, func(raw string) string {`（绑定半＝派发处，在 `installPanelTransport` 内，`:801–810`；`:808` `w.Init(panelPostMessageForwardInit)`）
- 整包 `.Bind(` 产码点＝**两枚**：`:802`（传输门）与 `:852` `w.Bind("wispProbeRT", ...)`（`firstRoundTrip` 里的一次性探针门，⛔ 不属传输名册，本尺射程具名排除）。
- **没有任何尺断言"这三处绑定字面量 ↔ 传输名册一致"**——复认普查 R5 那句"没有任何尺断言产码含这枚 `Bind`"。现存读 `installPanelTransport`/`panelDispatchBinding` 的四枚测试（`panel_transport_35r1_test.go:195→:203` 直调、`panel_inbound_guards_35r3_test.go`、`panel_transport_35r2_test.go`、`panel_transport_live_35v2_windows_test.go`）**全是能力形**（真调 `installPanelTransport`），不读源、不断绑定字面量与名册的锁步关系。⇒ 普查结论成立，本腿照 A717 §5 收窄后动手。

## 落点 / 形状定案（依派单形状要求）

- 新建 `cmd/wisp/panel_dispatch_binding_roster_253r1_windows_test.go`，`//go:build windows`、`package main`；⛔ 不改任何产码、⛔ 不碰 `internal/panel/**`、⛔ 不碰 `internal/agent/approval/**`（该包此刻有一格被按住，与本腿无关，我不进）。
- 射程＝**整包产码**：抄现成两把 `panel_locked_naming_33r11_windows_test.go:174`（`os.ReadDir(".")`＋逐文件 `os.ReadFile`＋零枚即 `t.Fatal`）与 `panel_geometry_255r6_range_windows_test.go:40`（剔 `_test.go`＋排序＋零枚即 Fatal）。名册为空 / 找不到 `installPanelTransport` / 传输内零绑定门 ⇒ **`t.Fatal`，不绿**。
- 路径基准具名＝**包目录**（`os.ReadDir(".")`，`go test` 的 cwd＝包目录），文件注释写明读的是**工作树**；⛔ 不用相对仓库根的字符串路径（`A716` §4 的假阴形状）。
- 枚具自证＝**盘上种＋`git hash-object` 拉平还原**（`A713`/`A716` 定式）；⛔ 不给 `go test -overlay` 的绿背书。现量行尾符尺：`tr -cd '\r'` 工作树 `panel_host_windows.go`＝**0**、`git show HEAD:`＝**0**、`git hash-object` 工作树＝HEAD blob＝`58e2b155dbcc8d380793e90bffc5514ab00743f1` ⇒ 该文件盘上即 LF，还原用 `git show HEAD:<file> > <file>`（纯读+写、非任何 git 变更命令）后逐字相等。
- 正控夹具＝本文件自带源码串（overlay 可见那一支，仅证明匹配逻辑会咬、不误咬；**不能**替盘上种背书，见 `probes/ruler-dedup-1/ruler-dedup-1.md:65` 33r11 注）。

## 判据与票面 AC#1 的关系（先写死，回报再对一句）

票面 AC#1（`:16`）原判据＝"把可达性裁决从词面尺迁到能力形＋补一发只 `postMessage` 不 `wispDispatch` 的正控"；§8 `:42` 已把真洞**收窄到窄义**＝"没有任何词面/AST 尺断言产码里 `w.Bind("wispDispatch")` 存在"。`A717` §5 再收窄＝能力形那半今天已由 35r1/35r2/35r3 落地，本腿**只补绑定字面量那枚词面/AST 形**。⇒ 出入点：票面"正控（只 postMessage 不 wispDispatch 的假腿必须红）"属**能力形激励**那一维，已由现存 35r1 覆盖，不在本枚词面尺射程；本枚的 M1 正控是**词面/AST 层的"绑定字面量与名册不一致"**。这一点在回报里具名，不擅改票面判据、⛔ 不翻任何框。
