# 票65 · jchv WebView2 包装层缺 `ICoreWebView2Controller::Close` —— 根治证据件

- 时间：2026-11-07 16:3x +08
- 执行腿：Wisp 编队写腿（本件由实现腿现跑取证；对抗验收另出）
- 缺陷判定来源：前腿决定性读数（`9a11eff81`／`1cab5ed6a`／`79fdc5a8b`）——修前全套 `count=10` 实测 39/40，像素非黑率逐轮坍塌 `40.6%×7 → 27.2 → 7.7 → 0.000×3`，根因＝`destroy()` 从不关 controller，每轮泄漏一组 `msedgewebview2` 进程积 GPU 压力；测试侧收割器已落地（40/40 绿）但产码级泄漏仍在。

---

## §1 步1 现跑读数：jchv 绑定层的位置与形状（2026-11-07 现跑，非旧账）

### 1.1 位置

```
$ go list -m all | grep -i webview2
github.com/jchv/go-webview2 v0.0.0-20260205173254-56598839c808
```

- module cache：`D:\work\base\gopath\pkg\mod\github.com\jchv\go-webview2@v0.0.0-20260205173254-56598839c808`
- 体积：565K，40 个 `.go` 文件；无 go.work；Wisp 无 `vendor/` 目录。
- 仓内真实使用者（排除 `.scratch` 克隆）：
  - `internal/ball/orb_windows.go`（本票写面）
  - `cmd/wisp/panel_host_windows.go`、`panel_resident_windows.go`、`panel_control_windows.go` 等（panel 家族，⛔ 本腿不碰，replace 后自动共享 fork）
  - `internal/panel/composer_dispatch_test.go`

### 1.2 形状：vtable 已映射 Close 槽，缺的只是 Go 包装方法

`pkg/edge/ICoreWebView2Controller.go` 的 vtbl 结构（11–34 行）：

```go
type _ICoreWebView2ControllerVtbl struct {
	_IUnknownVtbl
	GetIsVisible                      ComProc
	...
	NotifyParentWindowPositionChanged ComProc   // 行32
	Close                             ComProc   // 行33 ★槽已映射
	GetCoreWebView2                   ComProc   // 行34
}
```

- 槽序与 WebView2 官方头文件一致：`IUnknown`(3) + 20 个 Controller 方法后，`Close` 是 IUnknown 之后第 22 个方法（紧跟 `NotifyParentWindowPositionChanged`、在 `GetCoreWebView2` 之前）——**票面预估"第 6 号槽左右"不成立，实际在尾部**；jchv 的槽序没有排错。
- `ICoreWebView2Controller2`（`{c979903e-d4ca-4228-92eb-47ee3fa96eab}`）的 vtbl 同样已含 `Close ComProc`，其后再挂 `Get/PutDefaultBackgroundColor`——与官方 `ICoreWebView2Controller2` 一致。
- **全模块 grep `func.*Close|\.Close\.Call` = 0 命中**：`ICoreWebView2Controller` 的 Go 结构只有 `AddRef/GetBounds/PutBounds/AddAcceleratorKeyPressed/PutIsVisible/GetICoreWebView2Controller2/NotifyParentWindowPositionChanged/MoveFocus` 8 个包装方法，**没有 `Close()`，也没有 `Release()`**。
- `vtbl` 字段与 vtbl 结构类型均不导出 ⇒ **Wisp 侧代码无法直接触到 Close 槽**，包装方法只能落在 `edge` 包内。

### 1.3 controller 的可达路径（不用碰私有字段）

`pkg/edge/chromium.go`：

- 行20：`controller *ICoreWebView2Controller`（私有字段，Embed 完成回调里填充）
- 行328：`func (e *Chromium) GetController() *ICoreWebView2Controller`（**公开访问器，已存在**）

⇒ `o.browser.GetController()` 直接拿到 controller 指针，无需 `QueryInterface`、无需自己写 COM 调用。

### 1.4 泄漏现场（产码现状）

`internal/ball/orb_windows.go` `destroy()`（本腿起手时行470–486，wave62 在飞腿未动此函数）：

```go
func (o *orbWindow) destroy() {
	if activeOrb.Load() == o {
		activeOrb.Store(nil)
	}
	o.browser = nil            // ← 只丢 Go 引用，COM 世界无人关 controller
	if o.hwnd != 0 {
		pDestroyWindow...
	}
}
```

`o.browser = nil` 只是丢掉 Go 侧指针；controller 的 COM 对象与整棵 `msedgewebview2` 浏览器进程树无人通知关闭。与测试侧收割器的实测读数（每轮 reap 1 组、count=10 泄 10 组）自洽。

### 1.5 结论（修正票面预估）

缺陷比票面预估**浅一档**：不需要"补 vtable 槽"（槽已在），需要的是给 `edge` 包补 `Close()`（顺带 `Release()`）两个各 7 行的包装方法。而这两个方法只能落进 `edge` 包内部 ⇒ 途径只剩 fork/vendor 或上游 PR（上游 PR 不在本腿射程）。

---

## §2 步2 裁定：选甲（fork 到仓内 `third_party/go-webview2` + `replace`）

按铁律 ⓪⓪⓪（判据＝最有价值＋最符合主流 harness，⛔ 不是动作最小）：

**选甲。** Go 生态对"依赖缺一个小方法"的主流做法就是 vendor/fork 补丁（`go mod vendor` 后就地补，或 `replace` 指向仓内 fork 副本）；本例 fork 副本仅 565K/40 文件，补丁是两个 7 行包装方法，成本极低。

- **甲（选）**：全模块拷入 `third_party/go-webview2`（⛔ 不动上游 module cache），`go.mod` 加 `replace github.com/jchv/go-webview2 => ./third_party/go-webview2`，fork 内补 `Close()`/`Release()`。
  - 为什么最有价值：①类型安全——Close 走的是 jchv 自己已映射好的 vtbl 槽，不新增 unsafe 声明面；②修复面共享——`cmd/wisp` panel 家族与 `internal/ball` 全部 jchv 使用者同享，panel 关窗路径将来也能用同一方法；③未来再撞绑定缺口（`ICoreWebView2_4+` 等）可就地补，不必再开票等待上游；④已钉 pseudo-version，无自动漂移风险，上游更新＝手动同步一次。
  - 代价（如实）：仓内多一份 565K 依赖副本；上游演进需手动跟。
- **乙（不选）**：在 Wisp 产码里用 unsafe 自声明一份 25 槽 COM vtbl 布局、对 `GetController()` 返回的指针做转型调槽。
  - 具体代价：产码里复制第三方模块的内部布局知识＝双份真相源，25 个槽位的重复声明只为够到一个槽；jchv 一旦改布局即静默腐坏（错槽＝调错函数，不是报错）；unsafe 转型面在安全敏感的产码里扩大。复杂度花掉了，但没买到 fork 买不到的任何东西——按铁律⓪⓪⓪第3条，不选。
- **丙（不选）**：只文档记账"已知泄漏"。
  - 具体代价：⓪⓪⓪ 逐字⛔"动作最小"路线；且泄漏是实测像素逐轮坍塌到全黑的根因，收割器只是拆掉测试可见后果，不修产码＝留着病根。

---

## §3 步3 落地读数

（步3 完成后回填）

## §4 步4 门禁读数

（步4 完成后回填）
