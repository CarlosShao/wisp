# 票 260 · AC#2 落地腿 `260-r5` — 装一发**会响**的尺（非 winlive 看见"借了但没成键"）

**派单时刻锚**：HEAD `c43a791a`（`dev`）；本腿写面＝只新增 `internal/ball/*_260r5_test.go`＋本证据件目录。
**本格只有一格**：票 260 的 **AC#2**（零症状那一格要有声）。⛔ 票面任何 AC 框本腿一枚没碰（翻勾归编排者）。
**读码／现量时刻**：2026-10-04 14:0x–14:1x +0800，本机（win32／go1.27.1 windows/amd64）。

---

## 1. 可达性那一问的答案（先答这一问，再动手）

**答：可达，但不是编排者猜的那一形。** 三句话：

1. **失败只可能来自真 Win32**：`Ball.TakeEscForCancel`（`internal/ball/ball_windows.go:881-903`）里
   `takeEscBorrow(b.hwnd, borrow)`（`:892`）→ `hotkey_windows.go:631-633` → `takeEscWithAcc(hotkeyUnregisterer(hwnd), hotkeyRegisterer(hwnd), b)`，
   而 `hotkeyRegisterer`（`hotkey_windows.go:412-424`）是**具体函数不是变量**，`registerFn` 那道接缝在借用这条路上**根本没有被参数化**。
   ⇒ 除"真调用被拒"之外没有第三条产码可达的非 nil 路径（既有的 `takeEscWithAcc(unreg, reg, b)` 可注入形只被测试直接调用，不经 `TakeEscForCancel`）。
2. **编排者给的"hwnd 为 0 自然失败"这一形实测为假**（见 §2 的探针读数）：`RegisterHotKey` 的 `hWnd` 是 `_In__opt_`，
   **NULL 成功**（rc=1，注册成"线程热键"，随后 `UnregisterHotKey(NULL, id)` 亦 rc=1）⇒ 拿 0 当种子不但量不到失败，还会**真把组合键借到手**。
3. **可达的那一形＝无效（非零）窗口句柄**：实测 `rc=0 / errno=1400 ERROR_INVALID_WINDOW_HANDLE "Invalid window handle."`，
   三发不同形状的坏句柄读数一致。句柄校验发生在登记之前 ⇒ **什么都没注册**，机主的桌面不会被这枚用例碰到一个键。

而"进了那段闭包要不要真窗"＝**不要**：`Ball.uiRun`（`ball_windows.go:741-752`）的第一支是
`windows.GetCurrentThreadId() == b.sta.threadID()` ⇒ **原地跑**；`staThread.start` 在 `create` 之前就把 `s.tid` 记成本线程
（`sta_windows.go:75-77`），并且 `create` 返错时**根本不进泵、也不建窗**（`:84-90` 那条门，ticket 33 的
`sta_release_windows_test.go:194-310` 已经在非 winlive 档跑着同一枚门）。⇒ 本腿的台架＝**同一扇门**，
零真窗、零 `activeBall`、零窗口类，只把 `b.hwnd` 摆成"这个进程不持有的句柄"。

**结论＝能装。** 若这一支仍判做不到，需要的就不是接缝，而是"允许真开窗跑 winlive"那个词——本腿用不到它。

## 2. 现量凭据（本节记录的每一发都是本腿自己跑的）

- `grep -rn "TakeEscForCancel" --include=*_test.go internal cmd` 的读数与本腿要证的射程 ⇒ `logs/ruler-takeesc-testcallers.txt`。
- 非 winlive 那批测试文件的 build tag ⇒ `logs/ruler-build-tags.txt`。
- 失败回执写入真身位置（⚠ 文件名是 `ball_windows.go` 不是 `hotkey_windows.go`）⇒ `logs/ruler-receipt-writesite.txt`。
- `RegisterHotKey` 句柄形状探针（仓库外临时件，stdlib，从不建窗，成功即同进程释放）⇒ `logs/probe-registerhotkey-handle-refusal.txt`。
- 基线：改动前 `go test ./internal/ball/` ⇒ `logs/baseline-gotest-names.txt`（唯一红＝既有环境红 `TestC21TableColourRowsMatchTokensCSS`）。

## 3. 新用例：文件名／用例名／每一枚断言钉哪一行产码

本节答"这枚尺到底断什么"：新增文件 `internal/ball/hotkey_borrow_refused_260r5_test.go`（tag＝`//go:build windows`，**不带** `winlive`），
用例逐枚具名、断言逐条指回 `ball_windows.go:892-901` 与 `hotkey_windows.go:609-617`/`654-664`/`695-700` 的产码行，
并写明"台架里没有一枚回执是测试自己塞进去的"这一句怎么被钉住（种子回执走 `registerAllWith`＋假登记器，
cancel 槽按票 245 是 standby、假登记器对它**一次调用都没发生**，失败行只能由产码 `:895` 写出来）。

## 4. 判据①：不依赖真窗的证明

本节答"它真的被执行到了"：文件 tag（`windows`，无 `winlive`）＋普通 `go test ./internal/ball/`（不带任何 tag）里的
`--- PASS` 行逐字，以及"台架里 `IsWindow(种子句柄)==0`"这枚前提钉——前提塌了用例是 **红**，不是 SKIP。
⛔ 本腿全程没有 `t.Skip`，也没有跑 `-tags winlive`。

## 5. 判据②：正控（临时摘掉那处失败回执写入）

本节答"尺会响"：突变台件 `mutate260r5.py`（还原一律用**突变前从 HEAD 抽的副本** `ball_windows.go.head-blob`，
不是"finally 里重读当前文件"那一形），每一发当场打印 md5 与 `matches_start`，终值必须等于起手值；
红句**逐字**入本节，并点名是哪一条断言红、哪几条仍绿（仍绿的为什么不该红）。

## 6. 判据③：反控（防"恒真判据"）

本节答"这枚尺敏不敏感"：两发——（甲）把 `:895` 那行改成**无条件写**（去掉 `if err != nil` 的门）⇒ 必须仍能区分；
（乙）"成功形状"那一腿（`escTakenOver` 已置位、回执仍 live、`Problems()` 为空）⇒ 若回执是句永远成立的话，这一腿就得红。
区分不出来的那一部分本腿具名写"我这枚尺在这一维不敏感"，⛔ 不糊。

## 7. 门禁读数（逐条带时刻）

本节记录：`sh scripts/d22scan.sh`、`sh scripts/check-path-length-budget.sh --with-self-test`、
`"$(go env GOPATH)/bin/gofumpt.exe" -l <本腿名册内文件>`、`go vet ./internal/ball/`，
每发一条命令＋取数时刻＋读数档；既有环境红那一枚按"不归本腿、不修、不读成我改红的"处理并写明。

## 8. 名册与还原证明

本节记录：本腿每一枚 commit 的 `git show --name-only --format=`，`git status --porcelain -- cmd internal tools scripts .github docs`
终态（必须空），以及 `ball_windows.go` 的 md5 起手值＝终值那一行证据。

## 9. 判不动的地方／这格今天仍然欠的

本节答"还剩什么没验到"：⛔ 本腿不拿"接线层有尺了"去抵"端到端"那一跳——
`cmd/wisp` 把回执念成人看的那句话、以及真桌面被别人占着同一枚键那一形（`ERROR_HOTKEY_ALREADY_REGISTERED`＝1409），
今天仍要真开窗／真登记，归 winlive 与机主那个词；本节逐条具名并把"能不能不用真窗"分开写。
