# 299-a1 · 20-registration-api-inventory — 依赖侧可用的注册 API 名册

依赖坐标（尺＝`git show HEAD:go.mod` 逐字）：`github.com/jchv/go-webview2 v0.0.0-20260205173254-56598839c808 // indirect`。
模块目录（尺＝`go env GOMODCACHE` ⇒ `D:\work\base\gopath\pkg\mod`）＝
`D:/work/base/gopath/pkg/mod/github.com/jchv/go-webview2@v0.0.0-20260205173254-56598839c808/`。

**尺声明**：下面每一枚 file:line 都由
`grep -rn 'AddWebResourceRequestedFilter\|WebResourceRequested\|CreateWebResourceResponse\|ICoreWebView2WebResourceRequestedEventArgs\|ICoreWebView2WebResourceResponse\|GetResponse\|WebResourceResolved\|put_Response' <模块目录>`
现跑得到 ⇒ **总命中 49 行**（先 `sed` 掉目录前缀落 `$TEMP/wisp-299-a1/api-hits.txt`，再逐条开窗）。所有路径都相对模块根；`pkg/edge` 的导入路径＝`github.com/jchv/go-webview2/pkg/edge`。
⛔ 大段源码没抄进本件，每枚只给行号＋逐字短片段。

## 1. 名册（逐枚带 file:line）

### 1.1 注册面（过滤器）

| API | 出处 file:line | 逐字片段（签名/声明首段） | 谁的 |
|---|---|---|---|
| `Chromium.AddWebResourceRequestedFilter` | `pkg/edge/chromium.go:292` | `func (e *Chromium) AddWebResourceRequestedFilter(filter string, ctx COREWEBVIEW2_WEB_RESOURCE_CONTEXT) {` | **导出、可直接用** |
| ↳ 它内部委托 | `pkg/edge/chromium.go:293` | `err := e.webview.AddWebResourceRequestedFilter(filter, ctx)` | 库内 |
| `ICoreWebView2.AddWebResourceRequestedFilter` | `pkg/edge/corewebview2.go:383` | `func (i *ICoreWebView2) AddWebResourceRequestedFilter(uri string, resourceContext COREWEBVIEW2_WEB_RESOURCE_CONTEXT) erro`（行尾被裁，返回值枚数本腿未读） | **导出；注册的真身＝CoreWebView2** |
| vtbl 槽位 Add / Remove 过滤器 | `pkg/edge/corewebview2.go:131`、`:132` | `AddWebResourceRequestedFilter          ComProc` / `RemoveWebResourceRequestedFilter       ComProc` | COM vtbl |
| vtbl 槽位 Add / Remove 处理器 | `pkg/edge/corewebview2.go:129`、`:130` | `AddWebResourceRequested                ComProc` / `RemoveWebResourceRequested             ComProc` | COM vtbl |
| 处理器**实际挂载点** | `pkg/edge/chromium.go:211` | `_, _, _ = e.webview.vtbl.AddWebResourceRequested.Call(` | 库自己在 init 里挂 |

★ 挂载点那一段的**位置就是判据**：`:211` 落在 `func (e *Chromium) CreateCoreWebView2ControllerCompleted(res uintptr, controller *ICoreWebView2Controller) uintptr`（`pkg/edge/chromium.go:186`）体内，同一个块里还挂着 `AddWebMessageReceived`（`:202`）与 `AddNavigationCompleted`（`:217`），函数尾部 `:225` 是 `_ = e.controller.AddAcceleratorKeyPressed(...)`、`:227` `atomic.StoreUintptr(&e.inited, 1)`。⇒ **`WebResourceRequested` 事件处理器是库无条件挂的**，应用侧⛔ 不需要（也没有 API 去）自己 `add_WebResourceRequested`。

### 1.2 回调面（应用怎么写进去）

| API | 出处 file:line | 逐字片段 |
|---|---|---|
| 导出字段（唯一的注入口） | `pkg/edge/chromium.go:42` | `WebResourceRequestedCallback func(request *ICoreWebView2WebResourceRequest, args *ICoreWebView2WebResourceRequestedEventArgs)` |
| 处理器实例的创建 | `pkg/edge/chromium.go:64` | `e.webResourceRequested = newICoreWebView2WebResourceRequestedEventHandler(e)` |
| 处理器字段声明 | `pkg/edge/chromium.go:27` | `webResourceRequested  *iCoreWebView2WebResourceRequestedEventHandler` |
| 库侧回调蹦床 | `pkg/edge/chromium.go:281` | `func (e *Chromium) WebResourceRequested(sender *ICoreWebView2, args *ICoreWebView2WebResourceRequestedEventArgs) uintptr {` |
| ↳ 判空后转调 | `pkg/edge/chromium.go:286-287` | `if e.WebResourceRequestedCallback != nil {` ⇒ `e.WebResourceRequestedCallback(req, args)` |
| COM impl 接口 | `pkg/edge/ICoreWebView2WebResourceRequestedEventHandler.go:29-31` | `type _ICoreWebView2WebResourceRequestedEventHandlerImpl interface {` ⇒ `WebResourceRequested(sender *ICoreWebView2, args *ICoreWebView2WebResourceRequestedEventArgs) uintptr` |
| Invoke 蹦床 | 同文件 `:25-26` | `func _ICoreWebView2WebResourceRequestedEventHandlerInvoke(...) uintptr {` ⇒ `return this.impl.WebResourceRequested(sender, args)` |
| vtbl 组装 | 同文件 `:36-40` | `NewComProc(_ICoreWebView2WebResourceRequestedEventHandlerIUnknownQueryInterface),` … `NewComProc(_ICoreWebView2WebResourceRequestedEventHandlerInvoke),` |

### 1.3 事件参数面（读请求／写响应）

| API | 出处 file:line | 逐字片段 | 备注 |
|---|---|---|---|
| 类型 | `pkg/edge/ICoreWebView2WebResourceRequestedEventArgs.go:18` | `type ICoreWebView2WebResourceRequestedEventArgs struct {` | 导出 |
| **`PutResponse`** | `pkg/edge/ICoreWebView2WebResourceRequestedEventArgs.go:27` | `func (i *ICoreWebView2WebResourceRequestedEventArgs) PutResponse(response *ICoreWebView2WebResourceResponse)` | ⇒ 派单里那枚 `put_...` 在本绑定里**真名是 `PutResponse`**（无下划线小写形式） |
| `GetRequest` | 同文件 `:40` | `func (i *ICoreWebView2WebResourceRequestedEventArgs) GetRequest() (*ICoreWebView2WebResourceRequest` | |
| vtbl `GetResponse` 槽位 | 同文件 `:12` | `GetResponse        ComProc` | ★ **只有 vtbl 槽位，没有导出 Go 包装**：尺＝`grep -rn 'func .*) GetResponse' <模块目录>` ⇒ **0 命中** |
| 读 URL 的唯一出口 | `pkg/edge/ICoreWebView2WebResourceRequest.go:29` | `func (i *ICoreWebView2WebResourceRequest) GetUri() (string, error) {` | 回调里想知道"页面要哪枚路径"只能靠它 |
| 响应类型 | `pkg/edge/ICoreWebView2WebResourceResponse.go:14` | `type ICoreWebView2WebResourceResponse struct {` | 导出；本绑定里只看到 `AddRef`（`:18`） |

### 1.4 造响应面

| API | 出处 file:line | 逐字片段 |
|---|---|---|
| `Environment.CreateWebResourceResponse` | `pkg/edge/corewebview2.go:169` | `func (e *ICoreWebView2Environment) CreateWebResourceResponse(content []byte, statusCode int, reasonPhrase string, headers` （行尾裁切；参数是 4 枚起，headers 之后本腿未再读） |
| vtbl 槽位 | `pkg/edge/corewebview2.go:159` | `CreateWebResourceResponse        ComProc` |
| 实际 Call | `pkg/edge/corewebview2.go:192` | `_, _, err = e.vtbl.CreateWebResourceResponse.Call(` |
| 返回类型 | `pkg/edge/corewebview2.go:191` | `var response *ICoreWebView2WebResourceResponse` |
| 拿到 Environment 的导出口 | `pkg/edge/chromium.go:299` | `func (e *Chromium) Environment() *ICoreWebView2Environment {` |

### 1.5 名册里**不存在**的（派单点名要逐枚判的几枚）

尺＝上面同一发 grep 的按符号计数（`grep -c`）：

| 名字 | 命中 | 判语 |
|---|---|---|
| `add_WebResourceRequested` | **0** | 本绑定不导出这个名；事件挂载由库自己在 `chromium.go:211` 完成（见 §1.1★） |
| `remove_WebResourceRequested` | **0** | 同上；只有 vtbl 槽位 `RemoveWebResourceRequested`（`corewebview2.go:130`） |
| `WebResourceRequestedEventData` | **0** | 那枚是 .NET/WinRT 封装里的形状，**不是这个纯 Go COM 绑定**的东西；写进方案会落空 |
| `get_Response` / `put_Response` | **0 / 0** | 小写下划线形式不存在；`GetResponse` 仅 vtbl 槽位（无导出包装），`put` 的可用名＝`PutResponse` |
| `WebResourceResolved` | **0** | 本绑定没带 `ICoreWebView2_8`/`WebResourceResolved` 那一族 ⇒ ⛔ 别按"用 Resolved 做审计点"设计 |
| `COREWEBVIEW2_WEB_RESOURCE_CONTEXT` | 有独立文件 | `pkg/edge/COREWEBVIEW2_WEB_RESOURCE_CONTEXT.go`（枚举文件存在，枚数本腿未展开读） |

## 2. 注册需要哪一枚控制器对象（具名回答派单那一问）

**要的是 `ICoreWebView2`（CoreWebView2 本体），⛔ 不是 `ICoreWebView2Controller`。** 证据：注册方法的宿主是 `pkg/edge/corewebview2.go:383` 的 `func (i *ICoreWebView2) AddWebResourceRequestedFilter(...)`，而库内部的 Call 用的也是 `e.webview.vtbl.AddWebResourceRequested`（`chromium.go:211`，`e.webview` 的声明＝`chromium.go:21` `webview               *ICoreWebView2`）；`e.controller` 那一支在本绑定里只被用来挂按键（`chromium.go:225` `e.controller.AddAcceleratorKeyPressed(...)`）。Controller 的意义是**间接**的：CoreWebView2 是 `CreateCoreWebView2ControllerCompleted(res, controller)` 这个回调的产物，所以"先有窗、先有 controller，才有能注册的 CoreWebView2"。

## 3. 本仓宿主现在手里握的是哪一枚／能不能拿到前者（具名回答）

- 握着的：`cmd/wisp/panel_host_windows.go:148` 逐字 `	w        webview2.WebView`；创建点 `:386` 逐字 `	w := webview2.NewWithOptions(webview2.WebViewOptions{`。
- `webview2.WebView` 是**顶层窄接口**（尺＝`git show`/读模块 `common.go`）：声明在 `common.go:26`，成员**只有** `Run`/`Terminate`/`Dispatch`/`Destroy`/`Window`/`SetTitle`/`SetSize`/`Navigate`/`SetHtml`/`Init`/`Eval`/`Bind`（`common.go:30‑83`）⇒ **零枚资源过滤器成员、零枚通往 `ICoreWebView2` 的访问器。**
- 底层对象在库里但**拿不出来**：`webview.go:102` `chromium := edge.NewChromium()` 造出的那枚 `*edge.Chromium` 只被塞进 `webview` 结构体（`webview.go:49` `type webview struct {`）的**非导出字段**（`webview.go:52` 逐字 `	browser    browser`），而那个字段的类型是**非导出接口** `browser`（`webview.go:38-47`：`Embed`/`Resize`/`Navigate`/`NavigateToString`/`Init`/`Eval`/`NotifyParentWindowPositionChanged`/`Focus`）⇒ **既没有导出访问器、接口本身也没有过滤器方法**。`webview` 类型也是非导出的（`webview.go:49`），⛔ 从包外做类型断言都点不到它。
- 顺带把"宿主手里的 `SetHtml` 到底是哪一发"钉住：`webview.go:393-395` 逐字 `func (w *webview) SetHtml(html string) {` ⇒ `w.browser.NavigateToString(html)`。**宿主调的 `SetHtml` 就是 `ICoreWebView2.NavigateToString`**（`chromium.go:123` `func (e *Chromium) NavigateToString(htmlContent string) {`）。
- 本仓自己有没有绕过这层的先例？有，但都在注释里而非产码：`panel_host_windows.go:42-46` 逐字承认 `AddWebResourceRequestedFilter is exported by pkg/edge` 并点名 `(*edge.Chromium).Resize()`；`go list` 侧的导入面（尺＝`git show HEAD:cmd/wisp/panel_host_windows.go` 的 import 块）只到 `webview2 "github.com/jchv/go-webview2"`（`:64`）与 `golang.org/x/sys/windows`（`:65`），**没有任何文件导入 `pkg/edge`**（尺＝`git grep -n 'jchv' HEAD` ⇒ 命中只有 5 枚产码/测试文件，全是 `webview2 "github.com/jchv/go-webview2"` 这一枚顶层包）。

**结论（本节终态）**：可用 API 是**齐的**（过滤器注册、回调字段、`GetUri`、`CreateWebResourceResponse`、`PutResponse` 五件套都有导出面，见 §1），但**齐在 `pkg/edge` 那一层**；本仓宿主握的是**顶层 `webview2.WebView`**，它⛔ 拿不到 `*edge.Chromium`、也拿不到 `ICoreWebView2`。⇒ "注册方 0 枚"（票面现量第 3 条，本腿复核为成立）不是漏写一行，而是**缺一枚能拿到注册对象的所有权改动**。代价落在件 `30`／`40`。
