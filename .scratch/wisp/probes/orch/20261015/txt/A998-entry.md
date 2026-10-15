
## A998（2026-10-15 23:0x +08，编排者）120-r3 收档＋165-r2 死腿遗产代提＋★A997 欠的「推后 CI 色」回填＋★`TestSessionListByWorkspace` 归因（⛔ 不是我那批）＋我自抓四枚坏尺

### 一、`120-r3` 收档（细节同释放记录 `599d35f7`，这里只落账）
- 尺①提交真在：`76e861a3 → 29be2b67 → 397da9b4 → 6f53583d`（另有 `7111d51b`）；尺②证据件真在：`docs/evidence/s1/120-sealhandle-fixture.md` 195 行＋读数目录 28 枚。
- 尺③＝我自己跑的：仓外导出树（`/d/tmp/orch120r3/*`，`.go` 1076 枚先证非空）里 `internal/winsec` 整包 **rc=0 / RUN=62 / PASS_TOP=40 / FAIL=0 / SKIP=0**，它七枚 `120r3` 用例逐枚点名 PASS（连耗时）。
- ★第二把尺具名并列：只数 POSIX 名册（`-run 'TestAC[1-9]POSIX'`）＝**rc=0 / RUN=36 / PASS_TOP=27 / FAIL=0**，而腿报的是 **38/38/0/0**。两把都对，**差 11 枚＝子测试**（`^--- PASS` 只数顶层 vs 含子测试）。⇒ 同一份日志两把尺，⛔ 不带尺名并排写。
- ⚠ 具名欠账：它交回的两发**去牙反控我没复跑**（只到它自己的读数）⇒ 归 `120-v2`，在那之前 ⛔ 拿它当代替我成对读数的凭据。

### 二、`165-r2` 死于 150 轮上限（零提交）⇒ 遗产三件＋读数件十二枚由我代提 `17ca21ee`
盘上遗产＝`cmd/wisp/run.go` +29 行（`for _, e := range tools.BuiltinPlanEntries()` 那处注册循环）＋两枚新测试件（`plan_wiring_165r2_test.go` 206 行、`plan_assembled_165r2_test.go` 压行后 260 行）＋`.scratch/wisp/probes/165/r2/` 十二枚读数件。
- **读数搬到我要提交的字节上重跑**（台件＝`git archive HEAD` 导出＋盖这三件）：改后 **rc=0 RUN=3 PASS=3 FAIL=0**；**正控＝物理摘掉那 4 行注册循环 ⇒ rc=1 RUN=3 PASS=0 FAIL=3**，红句逐字「接线 RED：tools.BuiltinPlanEntries 在 internal/ 与 cmd/ 的产码里调用点 0 枚，要 1 枚…命中清单＝[]」，另两枚各红约 11 s＝真走到装配那一跳才红，⛔ 红在某句文案上；还原后调用点回到 1 枚、三枚复 PASS。**这就是票面欠的"摘掉那一行必须红"**。
- 邻居尺没被打死：255 两枚＋324 两枚 **rc=0 RUN=4 PASS=4**；`go build ./...`=0、`go vet ./cmd/wisp/`=0、`gofmt -l` 三枚＝**列出 0 枚**。
- ★装配读数本身值得留：registry 10 枚含 `plan.present`、卡＝`plan.present`、`tool_call` 行 risk=L2 decision=allow、**会话授权行 0 枚**；第二发「方案卡 yes → 清单外 `fs.write` 自己再起一张卡 → 答 no 后 `IsError=true` → 目标文件未落盘 → 授权行 0 枚」。这两行是票 165 那条「批准只批这一批」第一次有装配级凭据。
- 我动的唯二两行＝压注释回 260 行预算（派单 `165r2-dispatch.txt:36` 那条款），⛔ 断言⛔ 判据，压前后各跑上面那三发、读数一字未变。
- 腿的账具名留：读数件编号 **00/01/01b/02/03/04/05/06/07/08/09/11 ＝ 缺 `10-*`**。⛔ 我猜它写过又删了，只登记盘上现量的缺号，验收腿必答「09 与 11 之间那格判据有没有落过盘」。
- **AC 框一枚不翻**：实现者＝死腿＋编排者代提，判语必须出自非实现者 ⇒ 排 `165-v1`。

### 三、★销 A997 那句「推后色待回填」（回填前任何"CI 会变绿"的话都不算说过——现在算了，因为量到了）
同触发才可比：推前基线 **38725733693**（`0ab64fe9`，push，failure）↔ 我推的 **38750040552**（`adec1db1`，push，failure）。另有 38729417779（`0ab64fe9`，**schedule**，success）⛔ 拿来当"推前是绿的"那把尺。
- **作业级六枚同色**：`test-core`／`test-windows`／`lint`／`slo-full` failure ＋ `slo-smoke`／`lint-frontend` success ⇒ 两发⛔ 变。`lint` 那发的失败步逐枚＝`gofmt (gofumpt)`、`gofmt (gofumpt) - the tracked set is the denominator`、`staticcheck`＝**推前就在的存量账**（A997 已登 2 枚内容脏件＋staticcheck 36 条那格）。
- **用例级三形，逐形分类口径不同**（这是本节唯一有信息量的地方）：
  ① **转绿**＝`TestTicket255RosterStillMatchesTheActualReadSites` 从红名册消失（我 21:5x 在 `8beaed39` 面上亲跑 PASS，CI 同形）。
  ② **首次露出**（⛔ 写成"新增红"）＝`TestAC14GoSideEvalPushReachesThePage`：推前那发的 `cmd/wisp` 步骤⛔ 跑到它——最后一条 `=== RUN` 停在 `TestPathResolverJunctionWindows`（整进程死亡＝票 314 那族），推前日志里 `grep -c 'TestAC14'` ＝ **0**。它这发第一次被跑到，⛔ 我这批让它变红。
  ③ **翻红**＝`TestSessionListByWorkspace`：推前两发 `--- PASS`（0.02s／0.35s），推后 FAIL。下一节把这枚判死＝⛔ 回归。

### 四、★`TestSessionListByWorkspace`：不是我推的那 29 笔，是一枚硬币；台账早就登记过、这次只是把成因读出来
- **差集尺**＝`git diff 0ab64fe9..adec1db1 -- internal/session internal/memory` ＝ **0 字节**；这枚用例的出身笔是 `427c62d0`（10-13，票 166-r1），⛔ 本波任何一笔碰过它或被它碰。
- **成对命中率**（同一台面、两棵仓外导出树、同一把尺 `go test ./internal/session/ -count=20 -run 'TestSessionList'`）：HEAD＝`RUN=40 / PASS_TOP=30 / FAIL=10`，其中 **ByWorkspace PASS 10／FAIL 10**；推前 `0ab64fe9`＝`RUN=40 / PASS_TOP=27 / FAIL=13`，其中 **ByWorkspace PASS 7／FAIL 13**。⇒ **改前也红、命中率约五成**＝存量间歇红，⛔ 回归。旁站同屏：同发的 `TestSessionListByTime` 两棵树都 20/20 绿。
- **机制三行读码**（内容锚，⛔ 行号当判据）：`internal/session/manage.go:93` ＝ `now := m.now().Unix()`（**秒级**）；`internal/memory/dao_session.go:125` ＝ `ORDER BY updated_at DESC, id DESC`；`internal/session/session.go:96` ＝ `rand.Read` 铸 id。⇒ 同一秒建的两条 `updated_at` 相等，次序落到**随机 hex 的倒序**。修法接缝已在这个文件里：`manage_test.go:15-22` 收 `*time.Time` 且每次前进一分钟，用它的用例（`:48`）全绿，传 `nil` 的两处（`:69`／`:84`）＝硬币。
- ⇒ **票 166 加 AC#6**（ⓐ 仪器改确定性＋改前复现读数与改后 20/20 缺一不作凭据；ⓑ 只答"生产同秒两条会话的次序机主看不看得见"，⛔ 本腿动 `dao_session.go`——tie-break 是行为变更，读完那格由我裁）。⛔ 另立新票（同一次"这枚红归谁"的事在盘上只记一处）。
- ★**自纠两笔**：① 我在本会话早些时候把这事写成"本地 10 发全红＝确定性"——**那是我把只属一发的现象写成了通则**，现量五成，那句作废；② 当时那条 `grep -c '--- PASS' = 0` 也⛔ 成立（同一日志里 `TestSessionListByTime` 明明 PASS），是我读错了那把尺的射程。原句⛔ 抹平，按本节现量为准。

### 五、我这几发的台面自抓（四枚坏尺，逐枚具名、⛔ 用它们的读数、件⛔ 删）
1. **`gofumpt` 本机没了**：`fmt_rc=127`＋`command not found`，`GOBIN=D:\work\base\gopath\bin` 里只有 `go`/`gofmt`。⇒ 本节的格式化读数一律写作 **`gofmt -l` 列出 0 枚**，⛔ 冒充 gofumpt 的 0；gofumpt 那把留给验收腿（或另装）。
2. **导出树里 `cmd/wisp` 测试起不来**＝`exit status 0xc0000135`＋`RUN=0`（DLL 没找到；`third_party/sherpa-onnx/*.dll` ⛔ 被 git 跟踪 ⇒ `git archive` 里没有它）。处置＝把三枚 DLL 与编译出的 test exe 放同一目录、cwd 设 `cmd/wisp` 再跑，`RUN` 才从 0 变 3。**`RUN=0` 的四个 0 ⓪次算读数**——它就是"根本没跑"。
3. **我自己的按用例名计数有两枚是死的**：那把尺列了 `TestSessionList`／`TestSessionListIncludeArchived` 两个名，`grep -c` 回 0/0，而真名册里**这两个用例不存在**（真名册＝`TestSessionListByTime`／`TestSessionListByWorkspace`／`TestSessionSearchTitleAndBody`…）。⇒ 定式同旧账：**"0" 旁边要站着一个已知存在的样本**，否则 0 是正则造的。
4. **`tee` 写进一个当时还不存在的目录**（脚本里 `mkdir -p` 在 tee 之后）⇒ 主日志没落盘；同批发 `sed -n '68p'` 指错了行（时钟那枚在 `:69`，`:68` 是函数签名行）。两件都只算手滑、⛔ 影响判语，坏脚本留盘⛔ 删。

### 六、编队／推送／待人
- **在飞**：`322-r1`（`internal/winsec`，落地腿）。⚠ 它中途工作树里 `winsec.go` 会引用尚未写出的符号（我量到 `go vet` 报 undefined）＝**共享编译面这一分钟的任何 Go 读数作废**，我这批读数全部搬到仓外导出树上跑正是为此。
- **本波放出两枚**：⛔ **未放出**＝两条派单正文已存档（`166r2-dispatch.txt`／`324v1-dispatch.txt`），放腿那一刻另发 `-send.txt` 记录。释放条件具名＝**CPU 回落到 70% 以下**（23:0x 现量 87%，占用者＝在飞的 `322-r1` 在编译）＋对应面空（`166-r2`＝`internal/session`；`324-v1`＝`cmd/wisp` 定向 `-run`）。两枚面不相交，条件一到就同批放。
- **待派（具名写释放条件，⛔ 算在飞）**：`165-v1`（票 165 五格＋AC#3① 装配判语，**释放条件＝`324-v1` 交完**，两枚同抢 `cmd/wisp` 定向跑）；`120-v2`（票 120 五格＋复跑 `120-r3` 未复的去牙反控，**释放条件＝`322-r1` 交完、winsec 面空**）；`287-v1`／`321-v1`／`323-a1` 按序在尾；票 120 的 Windows 半格（`winsec_windows.go:328` 缺同一对象闸门）另开一格。
- **推送**：本地领先远程（`origin/dev=adec1db1`）自本波起又有 `599d35f7`／`17ca21ee`＋本笔。判据⛔ 变：⛔ 把可能带红的状态推上去——`322-r1` 交完、我在它终态上重跑 `internal/winsec` 整包 `FAIL=0`＋三门禁 `rc=0`＋验 fast-forward，才推 `git push origin dev`。

### 七、你需要做什么
没有。（这一节里唯一"人话"的结论是：**那个偶尔跳出来的会话列表顺序不是被我这批改坏的，它一直是枚硬币**——硬币已经变成票上的一格，会由另一条腿去把它改成确定的。你屏幕上没有变化。）
