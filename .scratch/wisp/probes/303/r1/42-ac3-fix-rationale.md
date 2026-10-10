# 303-r1 `AC#3` 修复说明（门那一侧，cmd/wisp/panel_host_windows.go）

改动只落 `panelPostMessageForwardInit`（`installPanelTransport` 里 `w.Init` 注入的那枚 JS 钩子）：
新增一枚 `window.external.invoke` 包装，在库桩真正**发出 RPC 帧**的那一刻把 `inside` 置真，
使覆盖层把"库自己发的帧"逐字节交回 native 出口，而**只把页面自己直接 postMessage 的信封**折进门一次。

## 为什么这⛔ 是"把判据改软"
- 判据本体三句（`no report … within 15s (what DID arrive at the door: nothing at all)`）一字未动（改动只落 `panel_host_windows.go`，测试文件一枚没碰）。
- ⛔ 动 SLO 阈值／golden／`internal/observe/thresholds.go`／D32 预算面。
- ⛔ 放宽任何断言、⛔ 把 t.Fatalf 换 t.Skip、⛔ 造 `--- SKIP`。
- ⛔ 动 `frontend/src/**`／`design/**`（读页面源码一律对象层）。
- 纯**加**：现有 postMessage-origin 折入路径与 `inside`/`native`/`try-finally` 全部原样保留；
  新包装只在 direct-wispDispatch 入口把 inside 抬起（原本该入口 inside 恒假→帧被折回→门收到双包 RPC→Handle 路由前丢弃→nothing at all）。
- 不改绑定名、不新绑名、不动 C17 名册、不重解析/不重序列化任何 message（逐字节透传性质保留）。

## 承重的"没变松"证据（非实现者可攻）
- `cmd/wisp/panel_transport_35r1_test.go` 与 `cmd/wisp/panel_transport_35r2_test.go` 的**行为夹具逐字执行本枚编辑后的 `panelPostMessageForwardInit` 常量本身**（35r2:25-26 "read as a value from the shipping constant, never copied"），含其自带的再入帽与 `var native`/`try-finally` 支的突变（mutant）用例。
- 件 `40-transport-after.txt`：`go test -run 'Transport|ShapeA3|LegacySubShape|Unforwarded|ForwardingHook' ./cmd/wisp/` ⇒ rc=0、0 枚 `--- FAIL` ⇒ 编辑后的钩子在**同一把判据尺**下（①旧形仍报 re-entered 9/native 0；③ postMessage-origin 仍 native 1/door 1；门缺失形仍逐字节交 native）全绿。
- ⚠ 这一步"⛔ 是放宽"由本腿**自陈理由、不盖章**；最终裁语归 `303-v1`（非实现者）。

## 反证成对（见 50-/51- 件）
- (a) 撤掉本枚 invoke 包装 ⇒ 三枚指名用例当场按预期红（逐字红句）。
- (b) `cmp` 逐字节还原本文件到 HEAD blob ⇒ 三枚复绿。
