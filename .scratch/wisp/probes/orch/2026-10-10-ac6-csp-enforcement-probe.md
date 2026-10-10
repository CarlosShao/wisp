# 票 299 `AC#6` — 那枚 CSP `<meta>` 在 `SetHtml` 产出的文档里到底执不执行（编排者真机一发，2026-10-10 10:2x）

**性质**：〔仅本机可量〕⇒ 按票面 `AC#6` 那格的规定，这一发**归编排者跑**，不算任何腿的欠账。
**载体**：一枚**仓外**一次性程序 `C:\Users\swq\AppData\Local\Temp\wisp-csp-probe-299\main.go`（`sha256` 前缀 `1b712d30fd1eabc20e476692`）。⛔ 不在仓内 ⇒ 不挪 `gofumpt`／`d22scan` 的分母、⛔ 不进任何门禁名册；它的 `go.mod` 用 `replace` 指向本机模块缓存（`GOPROXY=off`，⛔ 未取网络、⛔ 未从镜像站取哈希）。
**原始输出**：`.scratch/wisp/probes/orch/logs/ac6-csp-probe/runA.txt`（**405** 字节）与 `runB.txt`（**320** 字节）——逐字入库，本文件只加解释（两枚字节数我落笔前现跑 `wc -c` 对过）。

## 台件形状（每形各两发：带 meta／不带 meta）

同一枚 `webview2.NewWithOptions`（库的高层包装，与宿主 `cmd/wisp/panel_host_windows.go:386` 同一条构造路）＋ `SetHtml(整包字符串)`，文档里放四枚执行体＋一枚 Go 侧哨兵：

| 编号 | 执行体 | 它模拟什么 |
|---|---|---|
| ⓐ | 经典内联 `<script>window.wispProbe('classic-inline')</script>` | 页面里任何内联脚本 |
| ⓑ | `<script type="module">`（内联模块体） | **乙形（内联整包）要写的那一枚** |
| ⓒ | `<script src="./assets/probe.js">`（相对外链） | **本票的真缺陷**：页面按 URL 取自己的 js |
| ⓓ | `Blob` → `URL.createObjectURL` → `<script src=blob:>` | 常见"绕过外链"的替代写法 |
| ⓔ | Go 侧 `w.Eval("window.wispProbe('eval-probe')")` | **判活器**：宿主自己的脚本通道（不是页面内容） |

探针里那枚 meta 的 `content` 值**逐字取自冻结面**——两把尺对过（落笔前现跑）：`git show HEAD:frontend/index.html` 的 `content="default-src …"` 与工作树 `frontend/dist/index.html:14-16`（`Q-83` 甲那一次出厂链的产物）剥掉换行与空格后**逐字节相同** ⇒ 输出 `IDENTICAL`。⚠ 形状差异要说清：仓里那枚标签跨 3~4 行且写作 `/>`，探针里我把它排成一行、`>` 收尾——**`content` 值一字未动，排版的差别不是被测对象**（`default-src` 那枚值本身就是被测的东西）。
`default-src 'none'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self'; base-uri 'none'; form-action 'none'`

## 读数（两发，逐字来自 runA/runB）

**A 发＝不带 meta（正控：证明桥是活的、文档是能跑的）**——以下**整发逐字**＝`runA.txt` 全文（405 字节，⚠ 第一稿我在这里省了 `PROBE-VERDICT`／`PROBE-BASEURI` 两行而标题写着"逐字"，已补全）：
```
PROBE-REPORT classic-inline
PROBE-REPORT baseuri:about:blank
PROBE-REPORT module-inline
PROBE-REPORT blob-url
PROBE-REPORT eval-probe
PROBE-VERDICT withMeta=false reports=5 list=[classic-inline baseuri:about:blank module-inline blob-url eval-probe]
PROBE-SUMMARY classicInline=true moduleInline=true relativeLoaded=false relativeErrored=false blobRan=true evalProbe=true
PROBE-BASEURI baseuri:about:blank
```
**B 发＝逐字带上那枚 meta**——`runB.txt` 全文（320 字节，末行是 WebView2 自己的 stderr 噪声，不是判据）：
```
PROBE-REPORT eval-probe
PROBE-VERDICT withMeta=true reports=1 list=[eval-probe]
PROBE-SUMMARY classicInline=false moduleInline=false relativeLoaded=false relativeErrored=false blobRan=false evalProbe=true
[1010/102208.374:ERROR:ui\gfx\win\window_impl.cc:172] Failed to unregister class Chrome_WidgetWin_0. Error = 1411
```

## 判语（每条只说这两发撑得到的事）

1. ★**那枚 CSP `<meta>` 在 `SetHtml` 产出的文档里是真执行的**：同一份台件，去掉 meta 四形全跑、带上 meta 三形（经典内联／内联模块／`blob:`）全不跑。⇒ **普查腿 `299-a1` 件 `30` 判"乙形结构性不成立"这一条，从今天起有实测支撑，不再只是按 `script-src 'self'` 的字面推。** 乙形要成立只剩一条路＝动那枚冻结 meta，那是**人工批准面**，不在本票射程。
2. ★**`document.baseURI` 实测逐字 `about:blank`**（A 发ⓐ报回；B 发里那一条本身被 CSP 拦了，所以这一枚取自不带 meta 的同形台件——⚠ 射程要说清：它证的是"`SetHtml` 产出的文档基"，不是"带 meta 时页面会不会去敲门"）。⇒ 普查腿具名没交的那一枚欠项**今天销账**。
3. ★**判活器成立**：`eval-probe` 在两发里都到了 Go ⇒ B 发的三形沉默**不是**桥断了、也不是绑定被 CSP 拦，而是**页面自己的脚本被拦**。这一条是整发读数的牙齿；没有它，"三形没跑"可以有两种解释。
4. ⓑ**相对外链那一形在两发里既不 `onload` 也不 `onerror`**（`relativeLoaded=false relativeErrored=false`）。⇒ 只能说"今天没有任何仪器看见这一形发过请求"，⛔ 不许读成"确认没发请求"（监听器是我在经典内联脚本里挂的，而 B 发那段本身被拦 ⇒ B 发的这一枚**不可用**；A 发里挂上了却没触发，才是有内容的半条）。
5. ⚠**这一发看不见的一枚**：`about:blank` 下相对 URL 到底解析成什么（`new URL('./assets/x.js', document.baseURI)` 我没跑）。⇒ 归落地腿或 `AC#5` 那一发真窗，⛔ 本票结论不依赖它。

## 这发读数对票 299 排程的后果（写进票面裁定节 ⓗ，台账 `A798`）

- `AC#6` ⇒ 本格由我现跑，翻 `[x]`；`AC#0` 的三枚欠项全部销账（签名／枚举＝我 10:1x 复跑，`baseURI`＝本发）⇒ 翻 `[x]`。
- **本票剩下的唯一真门槛是"拿到那枚注册所需的对象"**（`A797` ⓐ），而它只有两支路：**重写宿主的建窗＋消息泵**（会撞票 33 已落的泵形状与多条已勾判据的凭据）或**动依赖**（撞"不新增依赖边"）。两支都超出我自己能批的射程 ⇒ **按"未定义即停"上报，摆给机主一句人话取舍（含"先不做"那一栏）**。
- ★**另记一枚候选，⛔ 不许静默采纳**：既然 ⓔ 证明宿主自己的脚本通道能穿过这枚 CSP，理论上可以把打包字节经 `Init`/`Eval` 注进去、根本不走模块加载。**这一发我没测它能不能跑起真应用**，而且它的性质是"绕开产品自己声明的那道防护"⇒ 属**安全相关的取舍**，要非实现者与契约面（D29/C 系列）裁，⛔ 谁都不许在落地腿里顺手做掉。
