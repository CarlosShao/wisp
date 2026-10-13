# 票 299 接续腿 `299-r1b` · 第 2 笔：`AC#2` ⓑ 反控（摘注册⇒白屏/红）＋复原证明

时刻 `2026-10-13 09:1x +08`。锚＝HEAD `91c1d8ab`。被反控的实现件＝`cmd/wisp/panel_control_windows.go`（`AddWebResourceRequestedFilter` 注册那一跳在 `:447`）。

**判语（一句话）**：注册那一跳被摘掉后，**页面用自己的话报回它落在了 Chromium 错误页**（`href=chrome-error://chromewebdata/ … root=no roottext=-1 resjs=0 rescss=0`），指名用例 `TestPanelPageFetchesItsOwnAssets299` **红**（20.66s／20.68s，两形各一发，rc=1）；复原**逐字节相等**（两处 sha256 对拉见 §3）＋同件正发绿（§2 先）。⇒ AC#2 ⓑ 的"摘掉必须红"成立，且实现件里**没有**变异形状残留。

---

## 1. 反控两形（同一摘法，两种驱动）

变异体＝`mut/m2-panel_control_windows.go.txt`：把 `c.browser.AddWebResourceRequestedFilter(panelAssetFilterPattern, edge.COREWEBVIEW2_WEB_RESOURCE_CONTEXT_ALL)` 这一行摘掉（回调赋值保留）⇒ 没有任何 pattern 匹配 ⇒ 对 `https://wisp.panel/index.html` 的导航**无人应答**，Chromium 走到网络（`.panel` 无 DNS）、落错误页。

- **形一（`-overlay`，仓内文件零字节未动）**：件 `logs/run-m2-ac2-filter-drop.txt`，overlay＝`mut/m2-overlay.json`。读数逐字：

```
    panel_assets_299_windows_test.go:218: 299 AC#2 page report (head 91c1d8ab, 10 attempt(s)): "href=chrome-error://chromewebdata/;base=chrome-error://chromewebdata/;ready=complete;sheets=4;rules=264;root=no;roottext=-1;resjs=0;rescss=0;"
    panel_assets_299_windows_test.go:221: the page is not on the panel's own synthetic origin: it reports href="chrome-error://chromewebdata/", want "https://wisp.panel/index.html" ...
    panel_assets_299_windows_test.go:229: the page reports #root present="no" with -1 character(s) of its own text: ... (resjs=0 rescss=0)
--- FAIL: TestPanelPageFetchesItsOwnAssets299 (20.66s)
rc=1
```

  两条读数说清：① 报告是**页面自己发的**（`window.wispDispatch` 通道在两形下都活着）；② `sheets=4 rules=264` 是 Chromium **错误页自带**的样式表，不是本页资源（本页资源两形下 `resjs=0 rescss=0`）=**白屏**那一形。

- **形二（物理摘除，按派单词面；随后复原）**：件 `logs/run-m2-ac2-filter-drop-physical.txt`。读数逐字（与前形同形）：

```
    panel_assets_299_windows_test.go:218: 299 AC#2 page report (head 91c1d8ab, 10 attempt(s)): "href=chrome-error://chromewebdata/;base=chrome-error://chromewebdata/;ready=complete;sheets=4;rules=264;root=no;roottext=-1;resjs=0;rescss=0;"
--- FAIL: TestPanelPageFetchesItsOwnAssets299 (20.68s)
rc=1
```

## 2. 正发（先/后各一发，页面自报）

- 摘掉**之前**：`3b806876` 自带的正发读数在 `probes/299/r1/logs/run-ac2-live.txt`（r1 跑，逐字 `sheets=1;rules=97;root=yes;roottext=100;resjs=1;rescss=1;`，PASS）。
- 复原**之后**（本腿现跑，`logs/run-ac2-positive-restored.txt`，rc=0）：

```
    panel_assets_299_windows_test.go:217: 299 AC#2 page report (head 91c1d8ab, 1 attempt(s)): "href=https://wisp.panel/index.html;base=https://wisp.panel/index.html;ready=complete;sheets=1;rules=97;root=yes;roottext=100;resjs=1;rescss=1;"
--- PASS: TestPanelPageFetchesItsOwnAssets299 (1.11s)
```

同一发在最终字节上再跑仍绿（`logs/run-final-green.txt`，1 attempt）。

## 3. 复原证明（sha256 对拉＋diff 逐行）

- `-overlay` 那一形**从未写过**仓内文件：脚注逐字挂在件尾（"cmd/wisp/panel_control_windows.go was NEVER written for it"），件头的 `sha256-before`（腿一开始、任何实验之前记录的）＝`7d257021bdc544a3fccb7971bb60ccacdd419da8f474748e0f844da34d8fdee7`。
- 物理那一形的复原脚注（件尾逐字）：`sha256 before mutation = ae714aee…`、`sha256 after restore = ae714aee…`、`RESTORE-BYTE-EQUAL yes`；复原源＝`mut/ctl-pristine-r1b.go.bak`（现跑 `sha256sum` 复核＝`ae714aeeafe2f24c584f2cb4eb14a36c3a75f252fc3340cbc969d88087de1e46`，三处相等）。
- 与锚件 `3b806876` 的 `git diff`：`cmd/wisp/panel_control_windows.go` 相对该锚**只有注释行**（`git diff 3b806876 -- cmd/wisp/`＝16 加/12 删，全部 `//`；本腿申报的另一笔）。**变异形状不在树里**：`grep -c "M2 reverse control" cmd/wisp/panel_control_windows.go`＝**0**；注册那一跳在（`:447`，`grep -c "AddWebResourceRequestedFilter(panelAssetFilterPattern" ＝1）。
- 复原后现跑正发绿（§2 尾）＝行为级复原证明。

## 4. 射程

- 反控的是"**注册那一跳**"（过滤器 pattern 的登记），不是回调赋值；这正是 `AC#2 ⓑ` 字面（"把注册那一行摘掉"）。回调仍在而没有 pattern ⇒ 匹配不到任何请求 ⇒ 无人应答。⛔ 不证明"若换一种摘法（连回调一起摘）页面会怎样"——那不是本格。
- 两发反控都让页面**说出了自己的处境**（错误页/白屏），不是拿 Go 侧计时/计数当判据——沿用票 33 `AC#14` 的 P2 定式。
- 变异体与 overlay 留着（`mut/m2-*`），审计者可逐字重跑；⛔ 未提交任何变异形状的**产品件**。
