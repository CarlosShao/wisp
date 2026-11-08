# 票 408 — `cmd/wisp` 整包测试跑不完（30/40/60 分钟级超时）：根因＝所有进程共用一个 WebView2 用户数据目录（UDF）；修法＝UDF 按进程唯一化（orbDataPath `<base>/wisp/orb/<pid>-<8hex>` ＋ 锁文件启动清扫）

**立票**：2026-11-08 15:31:58 +0800，落地腿（编排者派单「修 cmd/wisp 整包测试跑不完的结构缺陷」，技术裁定⓪⓪⓪＝按进程唯一化 UDF、⛔ 靠串行 test 绕过）
**性质**：★**结构缺陷票（测试台面稳性）**。与票 314（真窗区整进程无 panic 死亡、未定性）**同族不同根**：314 管的是「窗建起后整进程消失」，本票管的是「多进程共享 UDF ⇒ WebView2 fork 的 `EnvironmentCompleted` 拿回 `S_OK`＋NULL env 空指针崩 ⇒ recover 吞 ⇒ STA 线程死 ⇒ `ball.New` 无限等 ⇒ 整包必超时」。

## 根因（两腿独立现跑撞出＋编排者复核；本腿三核逐字吻合）

- 文件与行：`internal/ball/orb_windows.go:411-417`（修前）：
  `func orbDataPath() string { base := os.Getenv("LOCALAPPDATA"); if base == "" { base = os.TempDir() }; return filepath.Join(base, "wisp", "orb") }`
- ⇒ **所有进程共用一个 UDF 目录**。`go test ./cmd/wisp` 整包连跑时多个球窗用例（跨进程）同时起 WebView2，自己挤自己。
- **现量证据**：329 腿整包 `panic: test timed out after 30m0s`（1802s）；220 腿三发全被掐（40m／60m／60m＋全新 LOCALAPPDATA），死在三枚**不同**球窗用例——换全新 LOCALAPPDATA 仍死 ⇒ UDF 撞点不止一处（见 §残余共享态）。

## 修法（编排者裁定方向内，选型由本腿拍并写理据）

**选乙：`<base>/wisp/orb/<pid>-<8hex>`，`sync.Once` 进程内稳定。** 理据：
- 甲（纯 pid）：Windows pid 激进复用，同号残留目录会把陈锁/缓存态继承给新进程——正是要消掉的碰撞类。
- 丙（exe+pid）：exe 名对活进程零唯一性增量（pid 本就唯一），垃圾判别交给锁文件更准，且拉长路径（WebView2 对深路径历史敏感）。
- 乙：pid 保证活进程间唯一；8hex 随机段把 pid 复用风险压到 2⁻³²；`sync.Once` 必须——附体两次调 `orbDataPath()`（`browser.DataPath` ＝orb_windows.go:305、`siteDir`＝:363），site 目录必须是浏览器的目录。
- **垃圾界（⛔ 无限留垃圾）＝Chromium 式锁文件启动清扫**：锁**先**建、目录**后**建（清扫的唯一活性信号＝被独占持有的锁文件）；扫 sibling 时锁能独占打开 ⇒ 拥有者已死 ⇒ 连夹带锁删；锁忙 ⇒ 活主（含本进程）⇒ 跳；无锁目录与外来形状（`EBWebView`/`site` 等不匹配 `<pid>-<8hex>(.lock)?` 的旧布局遗物）**永不碰**。**不查 pid** ⇒ pid 复用不再能钉死垃圾。扫不动就留给下一跑；任意时刻垃圾界＝「活进程的目录」。修前共享目录本体（`%LocalAppData%\wisp\orb` 的 `EBWebView`/`site` 遗物）是**一次性遗留**，不被本清扫认领，可手删。

## 落地（本腿，e7f420a6c + c656c0ed8）

- `internal/ball/orb_windows.go`：`orbDataPath()` 重写（`orbDataBase()`／`orbDataDirFor()`／`acquireOrbDataDir()`／`orbSweepStaleDirs()`／`orbDataLeafRe`／`orbInstanceSuffix()`）。
- `internal/ball/orb_datapath_windows_test.go`（新件）四牙：①进程内幂等＋叶形＋含本 pid；②子进程 re-exec 异夹且**子进程自己的 pid**（读 `cmd.Process.Pid`）出现在子叶里；③LOCALAPPDATA 空 ⇒ 落 `os.TempDir()`；④清扫三态（死主删＋锁删／活主留／无锁与外来不碰）。
- **突变④**（导出树 `408-udf-tree` 现跑）：摘 pid 段 ⇒ 牙②红（`child's path leaf "bccc42d5" does not carry the child's own pid 20144`）；整叶唯一性去掉 ⇒ 牙②红（`child process returned THIS process's folder ...\orb\shared-leaf`）。还原后四牙全 PASS。

## 门禁读数（导出树 `/d/tmp/wisp/20261108/408-udf-tree`，sha c656c0ed8）

- gofumpt `-l internal/ball`：空（rc0）。
- `go vet ./internal/ball`：rc0。
- `go build ./...`：rc0。
- `go test ./internal/ball -count=1`：`ok ... 0.416s`。
- ★整包 `go test ./cmd/wisp/ -count=1 -timeout 40m`：见 Progress log（本腿最重要读数）。
- `-tags winlive ./internal/ball/ -run 'TestBallLiveOrb'`：见 Progress log。
- `scripts/d22scan.sh`：见 Progress log。

## 残余共享态（如实具名；本腿写面禁碰 cmd/wisp，未修）

- **cmd/wisp 面板测试件的固定名 profile 目录**（同一碰撞类）：`panel_host_windows_test.go:620,919`／`panel_host_windows_live_test.go:37` 共用 `os.TempDir()/wisp-33r1-panel-profile`；另有 `wisp-299r1-panel-profile`（:161）、`wisp-319r2-live-profile`（panel_death_replay_319_winlive_test.go:75）、`wisp-35v2-panel-profile`、`wisp-33r5/33r6/33r7` 等——**固定名**意味着跨进程跨跑次共享 UDF，正是本票在 orb 侧修掉的那类碰撞。若整包仍超时，**这是下一个具名嫌疑**。
- 产码面板 UDF＝`panel_resident_windows.go:305` `dataDir/panel-webview2`（按数据根一份；产码进程拓扑 D2 单常驻，不构成测试碰撞，但测试 fork 出的 resident 进程若共用同一数据根则同撞——归属票 314/128 族的数据根定型）。
- `scripts/spike/webview2-latency/main.go:151`（探针，非测试面）。

## AC

- [ ] AC#1 整包跑完：`go test ./cmd/wisp/ -count=1 -timeout 40m`（导出树）跑到终态（允许既有环境红，⛔ 再超时）。
- [ ] AC#2 winlive 全绿：`go test -tags winlive ./internal/ball/ -run 'TestBallLiveOrb' -count=1 -timeout 30m` 全 PASS。
- [ ] AC#3 门禁四件：gofumpt 空／vet rc0／build rc0／`go test ./internal/ball` 绿／d22scan clean。
- [ ] AC#4 若整包仍超时：如实报新死点并另立票（残余共享态 §是第一嫌疑清单）。

## Progress log

（追加式；每条 `- [YYYY-MM-DD HH:MM:SS +08] agent=… did=… next=…`，钟点由 `date` 的 stdout 插值，不手打。）

- [2026-11-08 15:31:58 +0800] agent=408落地腿 did=三核（根因串逐字在:411-417／调用点:305+:363 需进程内稳定／无 t.Parallel、activeOrb 单指针／票 314 非同根故新建本票）＋选型乙＋锁文件清扫落码 e7f420a6c/c656c0ed8＋四牙绿＋突变甲乙双红＋快速门禁（gofumpt 空/vet rc0/build rc0/单包 ok 0.416s）＋发现导出树缺未跟踪 DLL（0xc0000135，已从主树拷 sherpa-onnx+build 三 DLL）与 clone 伪红（autocrlf=true 把 embed 资产翻 CRLF，导出树内强制 LF 重检出后绿） next=整包 40m 读数＋winlive＋d22scan，读数入 Progress log
- [2026-11-08 18:52:39 +0800] agent=408r2腿 did=三核（两死点单跑皆绿 1.3-1.5s；fork Embed 无界泵+EnvironmentCompleted S_OK+NULL env panic 双门逐字对上）＋复现实验（同进程第二枚球 7.89s=单跑 6 倍；r2 前整包 393s 无声硬死于面板区 TestBookedSnapshot329 单跑绿）＋根因裁定＝进程内顺序多枚 orb 共享同枚 per-pid UDF⇒WebView2 同夹共享浏览器进程 JOIN 濒死树⇒controller 完成回调不触发/NULL env panic⇒close(started) 永不执行⇒ball.New 无限等（internal/ball 默认档零枚真 orb=该路径零覆盖）＋修复①orb per-boot 叶 boot-<n>（orb_windows.go，四牙零改动全绿，垃圾界随父夹清扫）＋修复②startResidentBallBounded 有界启动等待 2m+全量栈诊断（新测试件，接全部 8 处默认档调用点，断言/skip 口径零改动）＋验收整包 40m=931.402s 到终态 RUN=708/PASS=497/FAIL=54/SKIP=2（两死点变 120.01s 有界响红，活体栈抓到 EnvironmentCompleted NULL env 同步回调=与夹无关的 loader 出生即死行为；红名册=存量6+语音族30(onnxruntime System32 遮挡)+他腿落树欠账2(256Options/247prototypeVisuals)+重执行族6+两枚有界红）＋门禁 gofumpt 空/vet rc0/按域 build rc0/internal/ball ok/d22scan clean＋证据 docs/evidence/s1/cmdwisp-live-window-hang.md next=报编排者：fork nil-env 守卫（forks/go-webview2，落了则两枚红按设计退化 D2D 回退绿）＋waitStarted 有界化（产线语义须裁）＋存量红与语音 DLL 环境归各票
