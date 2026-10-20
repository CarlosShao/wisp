# 票 353 `AC#0`＝三形代价表（`353-a1`＝只读普查腿·续跑·⛔ 选形 ⛔ 动产码/测试字节 ⛔ 翻框 ⛔ push）

腿＝票 353 普查腿 `353-a1`（同身份续跑：断腿 17:10 落 1 笔 `eaa12987`＋`logs/02` 落盘未入库；本腿 17:3x 起在同一棵共享工作树续做）。
读面＝⛔ 任何 `go test`／⛔ 任何 `go build`（`352-r1`⇄`230-r1` 两枚写腿正在 `internal/winsec`⇄`internal/agent/approval` 落笔＝跑任何包都可能读到半保存态假红）。凭据只用 `git cat-file blob HEAD:<path>`＋`git grep -n`＋`342-v3` 既有 logs。
锚＝本腿续跑第一笔 `git log -1` 现认 `96a2eeb7`（17:34）；三枚承重 blob 与断腿 15:3x 锚 `c16ff47e` **逐枚同串**（`winsec.go f3c1b68e`／`pinnedleaf_342_windows.go db077e30`／`winsec_windows.go 63714a83`〔量的，尺＝`git rev-parse <锚>:<path>`，射程＝`internal/winsec` 三枚〕）＝本表与本腿复算读的是同一份字节。
散件＝`.scratch/wisp/probes/353/a1/logs/01|02|03|04`（01＝断腿读面件；02＝断腿自算 TOCTOU 窗；**03＝本腿独立复算＋两处具名更正**；04＝`AC#2` 仪器形状草案）。

## 0. 这一跳是什么（⛔ 本文选形，只给代价）
`privateFile` 的拒绝路径删的是**钉目录里的那枚名字**（`pinnedleaf_342_windows.go:216` `_ = p.root.Remove(base)`），而落地句柄共享掩码带 `FILE_SHARE_DELETE`（`internal/syscall/windows/at_windows.go:151`）⇒ seal 判语之后那枚名字可被腾空＋顶进一枚调用者⛔ 点名的对象＝顶进来的那枚被删。票面"现量"两枚去牙（〔量的，`342-v3` 射程＝它的仓外导出树，**本腿⛔ 复跑**〕）：摘删除那一步 ⇒ damage `3→0`；换按名开 ⇒ `entered 3→0`。
★三形都在"删除这一跳"上动；**覆盖写⇄独占那两条既有承诺三形都⛔ 改**（见第 4 节）。

## 1. 甲＝删除改按"句柄实际站着的对象"走
- **动哪几枚**：`pinnedleaf_342_windows.go:212-217`（`unlinkPinnedLeaf` 换形，且要在 `f.Close()` **之前**拿到那个动作＝`winsec.go:176-177`／`:181-182` 的 close⇄unlink 顺序要重排）＋`winsec.go:141-191` 路由两处删除支。原语可达性〔量的，尺＝`grep -rn 'func ReOpenFile'`，射程＝`$GOROOT/src/syscall` ⇄ `x/sys@v0.48.0/windows`〕＝**两处都⛔ 导出 `ReOpenFile`**（真身只在 stdlib 自家 `internal/syscall/windows/zsyscall_windows.go:478`＝本包够⛔ 着）⇒ 甲＝**新写一枚句柄面绑定**（本包已有先例形状＝`winsec_windows.go:754-793`：`NewLazySystemDLL` 绑 `SetFileInformationByHandle`＋`FileDispositionInfoEx`，现码就在用）。另一形 `FILE_FLAG_DELETE_ON_CLOSE`（开时就挂删除）＝**成功路径也会删**＝⛔ 是甲要的那味，本文不采它当甲。
- **覆盖写⇄独占**：⛔ 动旗标行（`winsec.go:96`／`:106`）⇄⛔ 动 `fs.ErrExist` 那一支。删除只在失败路径＝成功路径⛔ 经过 `unlinkPinnedLeaf`〔算的，见第 4 节〕。
- **C26 入口**：⛔ 动（`resolveString`／`ResolvePath`／`RemoveUnlinked` 一字不碰）。返回错误照旧丢弃＝⛔ 碰 `:207-211` 自陈那句"改 privateFile 的返回＝C26 契约变更"。
- **342/348 已批落地要不要重述**＝**要，且是最重的一味**：票 348 `AC#1` 已批文字（票面 `:44`）把落地形写成"两处删除都换成 `unlinkat(父fd, base, 0)` ⇄ `Root.Remove(base)`"＝甲**改的正是这批已批落地的删除形**；POSIX 面**没有**按句柄删的单发原语（`unlinkat` 只有按名那一形）⇒ 甲在 unix／fallback 两枚件⛔ 可implementable＝`AC#3` 那句"两平台成对"必须改成"Windows 修／unix 具名挂着"。
- **顶红谁现有的钉**〔算的，逐枚读断言行〕：★`ancestorwindow_342r2_windows_test.go:371-374`（`if kName == "present" ⇒ AC#5 RED: AC#5's half-made file survived the refusal`）在顶名世界里**必红**＝甲删的是句柄那枚对象、调用者那枚名字此时已属顶进来的对象 ⇒ 名字存活。同件 `:381-388` 那条"RESIDUAL OBSERVED"分支随之作废。⛔ 顶红＝`private_fail_test.go:26/:99`（无顶入的世界照旧删得动）、`openleaf_overwrite_322_test.go:36`（成功路径）、票 348 `ancestorunlink_348_windows_test.go:58-62`（`deleted` 只可能更小）。⇒ 甲落地必须**同时重述那枚仪器的判据文字**＝写腿活、⛔ 本腿权。
- **★要不要第二道产码语义批准＝要。** 必答① 的结论：**甲确实是"语义变化候选"**——删除从此只删句柄站着的对象＝改的是**删除这个动作的语义**（"拒绝路径清掉调用者那枚名字"变成"拒绝路径清掉我造的那枚对象"）。⇒ **这是第二道批准、归机主**，本腿⛔ 替机主点这个头，也⛔ 在表里写"建议选甲"。

## 2. 乙＝保留按名删＋删除前做一次"名字⇄句柄对象"同一性复验
- **动哪几枚**：`openPinnedLeaf`（`:171-199`）现在把两枚身份读数（`pinned`/`now`）**算完就丢** ⇒ 乙要带出来＝**函数形状改动**（返回一枚身份＋`unlinkPinnedLeaf` 多收一枚参数；`leafPin` 按值传 `:50-53` 装⛔ 下额外状态）；`winsec.go:177/:182` 两处传参。⛔ 新导出名、⛔ 新接缝＝⛔ 需要 `342-v3` 那枚"⛔ 该批的第二道接缝"批准（那道是**新产码接缝**，乙是**改内部函数形状**）。
- **覆盖写⇄独占**：⛔ 动（同甲：删除只在失败路径）。
- **C26 入口**：⛔ 动——**前提是复验拒了之后仍丢弃**（照 `_ =` 的形）。若乙要"我把没删这件事报告出来"⇒ 改的就是 `privateFile` 的返回＝**C26 入口契约**＝另批（`:207-211` 逐字写着这件事⛔ 归本落地）。
- **★乙那"同一性复验"本身有没有 TOCTOU 窗口＝有。窗口多大〔算的，本腿自己重读现码复算，⛔ 跑测试〕**：
  起点＝复验那发 `Root.Lstat` 的 `NtCreateFile`（`at_windows.go:144/151`）返回它所见 entry 之后；终点＝`Root.Remove`→`rootRemove`（`root_openat.go:210`）→`removeat`（`root_windows.go:352-353`）→`Deleteat`（`at_windows.go:236`）**第一枚 `NtOpenFile`（`:247`）按名解析并绑定那一刻**。中间逐枚＝`GetFileType`＋`GetFileInformationByHandle`＋`CloseHandle`〔3 枚句柄面〕＋`GetFullPathName`（`splitPathInRoot`→`rootCleanPath`→`syscall.FullPath`，`root_windows.go:59`）〔1 枚路径面〕＋`os.SameFile`／`incref`〔用户态 2 段；`SameFile`＝纯用户态 0 枚内核面，`types_windows.go:74-77` 注释逐字写明置空 `path` 就是为了⛔ 再查一次〕⇒ **内核面 4 枚、用户态 2 段**；`STATUS_ACCESS_DENIED` 重试支再 **+1** 枚在窗内的 `NtOpenFile`。量级＝**个位到十几微秒**〔算的，非量：本腿⛔ 跑任何计时〕。
  对照**现码敞口**〔算的〕：起点＝`openPinnedLeaf` 的 SameFile 应答（`:194`）、终点＝同一枚绑定；中间＝`sealHandle`→`applyDescriptorWindows`（`winsec_windows.go:291-320`：`foreignPrincipals`→`readDACL`＝**GetNamedSecurityInfo**＋`privateSetSIDSet` → **SetNamedSecurityInfo** → `verifyPrivate`＝**再一次 GetNamedSecurityInfo**）＋`f.Write`＋`f.Close` ⇒ **内核面 ≥8 枚、含 3 枚按名的安全描述符重操作（两读一写）**＝百微秒到毫秒级。⇒ **乙把敞口从"整个 seal＋write 段"缩到"删除前一跳"＝短 1～2 个数量级〔算的〕、⛔ 归零。**
  ★残余射程〔算的〕：乙对"顶入已先于复验"的世界**确定性拦住**（复验必不匹配⇒不删⇒damage 0）；对"复验之后、绑定之前才顶入"的世界⛔ 拦不住。而 `342-v3` 的跨进程那一发（〔量的，`probes/342/v3/logs/03`，射程＝它自己那棵台架〕）里攻击者**⛔ 赢这扇微秒窗也已成伤**＝plant 先于复验。
- **342/348 已批落地要不要重述**＝**轻**：删除原语照旧＝票 348 那批落地形一字⛔ 改；要重述的只有 `pinnedleaf_342_windows.go:165-167` 那句"the removal … cannot delete a foreign name either"（乙之后它变成"复验拦得住先到的顶入、拦不住窗内的"）。
- **顶红谁现有的钉**〔算的〕：同甲——`ancestorwindow_342r2_windows_test.go:371-374` 在"先顶入"世界里也必红（名字存活＝复验拒删）。同件 `:381-388` 的 residual 分支要改写。`openPinnedLeaf` 那两半门的**错误文字**⛔ 改 ⇒ `leafrace_342r2_windows_test.go:84-85`⇄`ancestorworld_342_348_windows_test.go:86` 的归因名册照绿（它们按 `why` 子串作差）。
- **要不要第二道批准**＝⛔ 是甲那道（删除语义⛔ 变）；**但**"复验拒了就留着那枚名字"是**攻击态可观测行为变化**，且它使那枚既有钉红 ⇒ 本文按"要人看一眼、⛔ 属第二道**产码语义**批准"报送，裁与摆归编排者＋机主。

## 3. 丙＝登记已知残缺＋常驻仪器把 damage 钉成可见红
- **动哪几枚**＝⛔ 动产码字节；新增一枚常驻仪器（形状＝`AC#2` 草案 `logs/04`）＋台账/票面具名登记＋`pinnedleaf_342_windows.go:11-20`／`:158-170` 那两处机制归因（★票面"禁区"已写明＝**另批**、⛔ 本腿权）。
- **覆盖写⇄独占／C26／342-348 已批落地**＝全⛔ 动（这正是它的代价：damage 照旧发生）。
- **顶红谁现有的钉**＝★**丙把新造的那枚红做成常驻红**：CI `test-windows`（`runs-on: windows-latest`）那一步逐字＝`bash scripts/winsec-tests.sh`＝本包整档都在这条道上，且 ci.yml 里 `continue-on-error` 现量＝**17 枚命中全在注释行、可执行行 0 枚**〔量的，尺＝`grep -n 'continue-on-error' .github/workflows/ci.yml` ⇄ 同串再过一遍 `grep -v '#'`＝0；射程＝整个 ci.yml〕⇒ 常驻红**直接顶住合并道**；另与公共块那句"唯一合法存量红＝`TestC26RewriteAccountIsConsumedAtEverySecurityLeg`（住 `internal/risk`）"相冲＝**新增一枚在册常驻红＝人工批准**（⛔ 本腿自裁）。丙⛔ 动任何既有钉的绿。
- **要不要第二道批准**＝⛔ 产码语义批准（语义一字⛔ 变）；要的是"允许一张永红的格"这道治理批准，归机主／编排者。

## 4. 三形共用的两行（⛔ 各形自己再辩一遍）
- **覆盖写⇄独占**：三形都只动**拒绝路径**那两处删除；成功路径根本⛔ 经过 `unlinkPinnedLeaf`〔算的〕⇒ 票 342 `AC#2`⇄票 322 家族（`openleaf_overwrite_322_test.go:36` 的四枚调用者子用例＋`private_fail_test.go:99` 的 `fs.ErrExist` 那一支）三形都⛔ 动。**但**甲／乙会让"拒绝后那枚名字已消失"这句**在顶名世界**失守＝`AC#5` 的那枚既有钉（见第 1／第 2 节）。
- **`342-v3` 那句"甲作废"**＝编排者 12:2x 收件笔已采纳它（⛔ 该批**第二道接缝**，理由＝批了也拿不到那枚"祖先在 open 之后翻"的世界）。⇒ 本表的"甲"与那枚被否的接缝**⛔ 是同一件事**：本票要买的世界（`logs/03` 顶名）今天已确定性量到、跨进程天然可达。若编排者仍把甲否掉，本票只剩乙／丙两枚可选。

## 5. 三枚必答（逐条作答）
1. **甲是⛔ 是"语义变化候选"＝是。** 删除从此只删句柄站着的对象＝改删除这个动作本身 ⇒ **这是第二道批准、归机主**（本文⛔ 选形、⛔ 替机主点头）。
2. **乙的复验有⛔ 有 TOCTOU 窗＝有。** 窗＝**内核面 4 枚＋用户态 2 段**（重试支＋1 枚）＝**个位到十几微秒**〔算的，尺＝第 2 节那串 `file:line`，射程＝`internal/winsec` HEAD blob ＋ `$GOROOT/src/os` ⇄ `internal/syscall/windows`（go1.27.1 windows/amd64）＋ `x/sys v0.48.0`〕；对照现码敞口＝**≥8 枚内核面、含 3 枚按名 SD 重操作**＝百微秒到毫秒级 ⇒ 乙＝**短 1～2 个数量级、⛔ 归零**〔全算的，本腿⛔ 跑任何计时〕。
3. **三形能否被一枚确定性仪器证伪**＝甲：**能**（damage 归 0 ＋ 无名对象存活，两枚都可确定性读）；乙：**部分能**（先顶入的世界可确定性证伪；窗内那一支⛔ 可证伪＝要在 `Deleteat` 的 `NtOpenFile` **之内**插钩子＝平台内部）；丙：**⛔ 证伪任何代码行为**（它把可证伪性搬到"读数常驻可见"，仪器恒红＝对"红存在"最好、对"行为对"零）。⇒ `AC#2` 的形状草案＝`logs/04`，★其中"乙窗内那一支测不到"必须逐字写进判据、⛔ 让仪器在一个它⛔ 覆盖的世界上报绿。

## 6. 本腿纪律自证（三把尺）
- 框终量＝票 353 **未勾 5／已勾 0**（尺＝`grep -cE '^- \[ \]'`⇄`'^- \[x\]'` 现量，射程＝该票面）＝`AC#0`..`AC#4` 一枚⛔ 翻；票名／`Status:`／`-done` 一字⛔ 动。
- 只读自证＝本趟 `go test` **0 发**、`go build` **0 发**、`internal/winsec/**` 字节 **0 枚**（写面只有 `probes/353/a1/logs/03|04`＋本件＋票尾一条追加）。
- ⛔ 把 `342-v3` 报的数当自己量过：本表凡引它的读数都带"〔量的，`342-v3` 射程＝它自己的仓外台架／本腿⛔ 复跑〕"；本腿自己的新数只有两类＝**blob 同串**（量的）⇄**窗口枚数与量级**（算的）。
