# 346-a1 普查件（票 346 第一趟＝只读；★本腿 0 枚 `go test`、0 个 `.go` 字节）

- 腿＝`346-a1`（普查）；起手锚现认＝`3abbb6391a8a497c0332e5363f24c8a250a7eb78`（`git log -1` 本腿第一笔）；分支 `dev`。
- 拆开列表＝同根 `logs/01-shape-120.md`（甲枚）⇄ `logs/02-shape-322.md`（乙枚）；
  取日志尺⇄归一⇄四数＝`logs/00-ci-fetch.md`；blob 同一性⇄射程＝`logs/04-blob-and-scope.md`；
  另外两味＝`logs/03-other-causes.md`。本件只放结论，⛔ 重复摊表。

## (i) 弃权形状（要点，逐格凭据在两份表里）
两枚都是 `//go:build !windows` ⇒ 射程＝**只有 Linux 档（CI `test-core` ⇄ docker 容器）**；
"本机 Windows"那一轴对这两枚**结构上不存在**（尺＝`orch227-b/full.log` 1544211B 与 `win.norm.log`
683457B 各 0 命中；正控＝同尺在 `core.norm.log` 命中且该文件有 7 枚顶层 FAIL＝⛔ 尺失明）。

| 形状 | 甲 `TestPOSIX120SealWindowLeafSwapIsClosed` | 乙 `TestAC1POSIX…322r1` |
|---|---|---|
| 红句原文（弃权响） | **观测 1 发**＝39124122121 `refused=1998, sealed-own=0`（B1319461..B1321464） | **观测 1 发**＝同 run `answered 0 … landed 6 flips`（B1276161..B1277160） |
| 绿（窗口到达） | 2 发＝39244425666（`refused-at-act=19`，B631380..B633837）⇄39101060753（`=5`，B638019..B640466） | 2 发＝同两 run（B587695..B588782 ⇄ B594096..B595087） |
| `SKIP` | **0 命中**（五份 CI 日志＋文件内 `t.Skip`＝0，正控同行） | **0 命中**（同尺） |
| 无 `--- FAIL` 的整包死 | **0 命中**（`setup failed|build failed|creating work dir|0xc0000135`＝0，同文件有顶层 FAIL＝尺灵） | **0 命中**（同尺） |
| 建目录即死 | **未观测**；但甲枚有第三死法＝`WISP_120_BASELINE=1` 读者档把弃权**降级成 `t.Logf`＝绿**（`:231-237`，观测＝`probes/120/v2/logs/v2-prebase-reader.log:17`） | **未观测**，且乙枚⛔ 任何降级出口（`:205-209` 直 `t.Errorf`） |
| 容器（`/gotmp` 绑 D 卷） | 绿但慢 460 倍（0.08s→37.15s）＝同档 `os.SameFile` 失真（`342-v1-verdict.md:15`） | 绿但慢 250 倍（0.42s→104.53s，同件） |
| 容器（⛔ bind，安静） | 绿 3/3（`342/a1/logs/24-container-swing-nails-3x.log` 506B） | 绿 3/3（同件） |
| CI⇄本机 × 长⇄短根 × D⇄C 卷 | **未知**＝两枚⛔ 在 Windows 侧编译；A1053 长根那发红册 29 名里两枚名 0 命中 | **未知**（同尺） |

★同 blob 两枚**同 run 同色**（三锚 blob 逐字同＝`git rev-parse`：`dfaddac8…`／`0b0e5173…`；
`winsec.go` 在 `a5080612→6f808a83` 之间**变过**＝`67783ae0`→`07300bff`，但两枚仪器自身⛔ 变）
⇒ 它们⛔ 是两枚独立随机事件，是**同一件事的两个读数**。

## (ii) 结论与形状（⛔ 本腿动手）
判＝**盘上确实有三味，但对"本票这两枚"只有两味成立——味乙（临时根拼法）对本票两枚被本腿★排除**：
- 味甲＝**仪器形状**（这两枚的主因，成立）：弃权那一支在两枚里⛔ 同形（甲有 `WISP_120_BASELINE=1`
  读者档把弃权藏进绿、乙⛔ 任何降级出口）＝同一判据下"弃权"有两种可观测形状，红⇄绿的翻转⛔ 可归因。
- 味丙＝**时序／争用**（这两枚，成立但只钉住半条）：三发 `test-core` 同作业同批包、两枚**同 run 同色**
  ＝同一件事的两个读数；★runner 负载维＝**未知**（本腿⛔ 任何负载读数）⇒ 这味"该记"、⛔ 该被写成已定案。
- 味乙＝**临时根拼法**（对这两枚**⛔ 成立**）：两枚 `//go:build !windows` ⇒ 造出拼法红的那味
  在 Linux 档⛔ 存在；尺＝`orch227-b/full.log` 那 29 枚红册里两枚名 **0 命中**（件03 第一节末行）。
  ⇒ 拼法那味**真**、但它是**那 9 枚 junction 一族**的成因（已在册＝`227-v2` Ⓒ 段更正＋票 311），
  ⛔ 属本票两枚的弃权成因。★本腿把这条列出来是因为派单把三味一起喂给我"供对回⛔ 采信"＝对回结果＝**其中一味⛔ 属于本票射程**。
理由（各带尺）：① 弃权那一支**在两枚仪器里形状⛔ 同**（甲有读者档降级＝`:231-237`、乙⛔ 降级出口＝`:205-209`，
均读码）＝改形状前⛔ 谈别的；② 拼法那味已被 `227-v2` 成对两发定案（`2/2/0/0` vs `2/0/2/0`）⇒ 归配方、⛔ 本票新账；
③ 争用那味本腿**⛔ 可证负载维**（无 runner 读数）＝写未知、⛔ 圆场。

**只给形状（新增常驻件一枚）**：
- 件＝`scripts/abstention-census.sh`（新；⛔ 动 `.github/**`＝人工批准，⛔ 做成 CI 步骤）；
- 判据＝输入若干份 `-v` 日志（CI 的或本机的），对**每一枚**"弃权形仪器"打印一行
  `名 | 颜色 | 该枚自己的那几枚到达数（attempts／flips／refused-at-act／sealed-own／leafLegRefusals）| 判定`，
  判定三值＝`REACHED`（到达数⊥0 ⇒ 这一发作数）⇄ `ABSTAINED-RED`（弃权红已响）⇄ `ABSTAINED-HIDDEN`
  （仪器绿**但**到达数全 0 ⇒ ＝被读者档／形状缺陷藏起来的弃权）；★取行尺＝按 `=== RUN` 段切、
  段内先取正文再取该段的 `--- (PASS|FAIL|SKIP)`（`go test` 正文印在结果行**之前**，⛔ 事后 grep）；
- 弃权形仪器名册（本腿现量尺＝`t.Errorf("AC#n RED instrument`／`t.Fatalf("AC#n RED instrument`，
  命中逐名 `file:line`＝12 处／8 个文件）＝`openleaf_322_other_test.go:206`、`posix_seal_window_120_other_test.go:225/:233`、
  `ancestoropen_342_unix_test.go:58`＋`_windows_test.go:56`、`ancestorunlink_348_unix_test.go:108/:111/:153`＋`_windows_test.go:45/:75`、
  `openleaf_guard_322_windows_test.go:241`、`sealhandle_privateroute_120_other_test.go:212`；
  ★⛔ 只交票面那两枚（票 346 `AC#0` 明令的同族尺）；另两处 `t.Skipf(前置条件缺失…)`＝
  `paths_workspace_test.go:205`／`task_output_pointer_notice_test.go:408` 属**第三种弃权＝SKIP 形**，
  在册（票 310 名册那两行）⇒ 本件对它只报"SKIP 那一支不参与本票"。

## (iii) 三味该不该各立续票（盘上现读，⛔ 重复立案）
| 味 | 判 | 凭据（现读） |
|---|---|---|
| 甲＝仪器形状（弃权三支不对称：一枚有读者档、一枚没有） | **⛔ 另立新票**＝归**本票 346 的 `AC#1` 丙形射程**（票面逐字写着"弃权照旧红但同屏打出可归因读数"＝同一味），只是 (ii) 那枚形状要写进本票的三形表；★本票 346 现量 0 勾、未结案 ⇒ 新票＝重复立案 | 票 346 `AC#1` 丙形原文＋`AC#5` 禁区（⛔ 动 winsec 字节）现读 |
| 乙＝台面拼法（`GOTMPDIR` 正斜杠 ⇒ 9 枚 junction＋SID 那族红） | **⛔ 新票**＝已在两处在册：①票 311（六格齐、⛔ `-done`、残余具名＝卷形⇄`Q-311-r1b-1`）②`227-v2` 件Ⓒ 已把 `A1053`⑵ 的根因更正为"拼法"并给了定式（`cygpath -w`／反斜杠形）；★补一句：票 311 现文**⛔ 覆盖"拼法"这一维**（尺＝`Select-String 'cygpath|正斜杠|反斜杠|拼法|Invalid switch|mklink'` 于票 311 文件＝**0 命中**＝量的）⇒ 该味**归本票派单公共块（(iv) 配方）＋`A1053`/`A1054` 那两枚台账更正**，⛔ 起票 | 票 311 现文尺 0 命中＋`227-v2/logs/03-…md` Ⓒ 段＋`probes/236/v2/raw/baseline-tools-v1-forward-slash-gotmpdir.txt`（同形早现） |
| 丙＝调度争用／CI 摆动那枚 `TestTicket223ModeLooseningChangesTheRunningModeAfterAllow` | **该立一枚，但⛔ 本票立**＝它是 `cmd/wisp` 包"1s tick ⇄ 等待预算"那一族的摆动（`config_reload_223_test.go:31/:98`），与本票两枚的"窗口未到达"**⛔ 同判据族**；同码三处两绿一红（件03 逐名 5 发）＋⛔ 任何未结案票具名它（尺＝未结案票池 grep 该名＝**0 命中**，量的）⇒ 形状＝**另立一枚 `cmd/wisp` 侧摆动票**（射程＝`reloadCaseBudget` 与 tick 的到达判据，⛔ 放宽断言、⛔ 加 `t.Skip`）；⛔ 塞进 346＝会让 346 的"仪器形状"结论被一包无关摆动稀释 | 件03 的颜色表＋blob 同串尺（`f1cac524…` 两锚同）＋`git log a5080612..6f808a83 -- cmd/wisp`＝0 枚 |
| 附：那枚 `TestTicket223PermissionDeniedSitsInItsOwnSentence` | ⛔ 任何新账＝四发 CI 全红＝**稳定红在册**（票 303／305／314 一族，票 346 现量表第 5 行已逐名归位） | 件03 表末行 |

## (iv) 可复用台面配方（原文，可直接粘进派单公共块）
> **★台面配方（2026-10-19 票 346 普查腿 `346-a1` 交；本波编队认可，⛔ 各自另编）**：
> 下面四味分别钉死**四类台面受害者**，缺对应一味就假红⇄假弃权（★每味自带射程，⛔ 一把通吃）——
> **①（射程＝Windows 档的 9 枚真 junction 族⇄SID 那族）临时根形＝Windows 反斜杠拼法**：
>  `GOTMPDIR` 一律经 `cygpath -w` 出来、⛔ 把 POSIX 形（`/d/…`、`C:/…`）原样 export 进去。
> 　凭据＝`GOTMPDIR` 前缀被 `os.MkdirTemp` **原样带进** `t.TempDir()`，junction 那族的 `mklink /J <link> <target>` 里 target
> 　一旦是 `X:/…` 正斜杠形，cmd 把它当开关吃掉＝`Invalid switch - "tmp"`（/`"Users"`）⇒ 红句逐字点名的"前置条件缺失"是**台损**。
> 　★与"长⇄短""D 卷⇄C 卷"无关（同树同包只差拼法＝反斜杠 `2/2/0/0` vs 正斜杠 `2/0/2/0`，凭据 `probes/227/v2/logs/03-unqualified-and-gates.md` 甲⇄乙两发）。
> 　本机 Windows 档＝`WISP_TEST_GOTMP_ROOT` 落 `%LOCALAPPDATA%` 同卷短根（票 311 认可形）＋反斜杠；
> 　容器档＝★⛔ 把 `GOTMPDIR` bind 到宿主 Windows 卷（`-v $S/gotmp:/gotmp` 那形＝`os.SameFile` 失真＋慢两个数量级），
> 　要指就指容器自己那层可写根、或干脆⛔ 指（让 `t.TempDir()` 落容器 `/tmp`）。
> 　★**本票那两枚仪器⛔ 吃这一味**（`//go:build !windows`＝Linux 档无 `mklink`；尺＝长根那发红册 29 名里两枚名 0 命中）。
> **②（射程＝`cmd/wisp` 那包 ⇄ 任何取 DLL 的 Windows `.exe`）`PATH` 形＝整条 Windows 形＋分号分隔**：`PATH="$(cygpath -wp "$PWD/third_party/sherpa-onnx:$PWD/build:$PATH")"`
> 　（★整条一起转、⛔ 半 POSIX 冒号形接半 Windows 分号形＝会造出"零枚 `--- FAIL` 的整包死"，账 `A1042`④／`A1053`④）。
> **③（射程＝本票两枚弃权仪器 ⇄ 一切"靠到达数作数"的仪器）`-timeout` 档＝`-timeout 30m` 或更长 ⇄ 必带 `-count=1`**（★票 346 `AC#4` 的"至少 3 发独立复跑"＝
>  复跑要落在**安静档**：同一步⛔ 与整包并发、⛔ 两测量脚本同写一棵导出树，＝`A1053`④ⓐ 那形）；
> 　★本腿⛔ 观测到任何**已具名的超时 kill** ⇒ ⛔ 许诺"`-timeout` 档能救这两枚"；配方里这一味的作用＝
>  **⛔ 让超时冒充判据**（`panic: test timed out`＝无 `--- FAIL` 的整包死＝台损），⛔ 改判据分母；
> 　判据只认"到达数"：`sealed-own>0`／`refused-at-act>0`／`flips>0`⇄`leafLegRefusals>0`，
> 　★总 `refused` 大⛔ 代表到达了——观测里红发 `refused=1998` > 绿发 `1768`，方向与直觉相反。
> 　`cmd/wisp` 整包那味真窗区整进程死亡照旧具名（⛔ 拿整包当全绿）。
> **④（射程＝一切 CI 日志读数，⛔ 测试台面本身）取 CI 日志＝先 `gh run view <run> --json jobs` 取**数字** `databaseId`、再 `--log-failed --job <数字ID>`；
> 　★喂作业名／带错 flag ⇒ 静默 0 字节＝"看起来⛔ 红"；GitHub 日志每行带 `job⇄step⇄时间戳` 前缀，
> 　四数必须先把 payload 顶到行首（剥两个 TAB＋时间戳＋一个空格）再在 `^=== RUN|^--- PASS|^--- FAIL|^--- SKIP` 上取。
> ★**配方的诚实边界（本腿必答的那句）**：
> ①那味对 **9 枚里的 4 枚**（`internal/risk` 那族）在 CI 已可证有效＝CI 自己发的 `GOTMPDIR` 就是反斜杠形
> （`win.norm.log` L7 逐字 `GOTMPDIR=C:\Users\runneradmin\AppData\Local\wisp-test-gotmp\winsec-…`）＋那 4 枚逐名 `--- PASS`（L4056/4060/4064/4066）。
> ★**另 5 枚（`internal/tools` 那族：bridge×6／sweep／20r2）在 CI Windows 档⛔ 现身＝⛔ 在该档射程内**
> （尺＝逐名在 `win.norm.log` 各 0 命中＋该文件包结果行只有 proc/secret/config/risk/ball/perm/plugin/llmrecord/session/projctx/audio＋winsec＋cmd/wisp）
> ⇒ 那 5 枚的"CI 两处都不弃权"＝**未知**（⛔ 许），只有**本机**（甲发反斜杠 `2/2/0/0` 绿）可证。
> ★**本票那两枚仪器⛔ 有任何一味台面配方能让它们"不弃权"**：弃权那一支是票 322 `AC#1` 逐字要求的（票 346「归因边界」第 1 条），配方只能让它们**把弃权说得出口**（③的到达数）＋
> 把"本机⇄CI"的**射程**讲清（两枚 `!windows`＝在 Windows 档⛔ 编译＝既⛔ 绿也⛔ 弃权，⛔ 拿"本机没红"当绿）。

## 写面自证尺（本腿自己现跑）
- 正文件 6 枚＝本件（`docs/evidence/s1/346-a1-census.md`）＋`logs/00..04`；逐枚 `wc -c` 非 0、行数≤160 帽。
- 被禁字面字符＝**码位 U+00A7**（分节符；★⛔ 与间隔号 U+00B7 混）：尺＝按码位在 UTF-8 文本里计数，
  六枚命中逐枚＝**0**（同时 U+00B7 也＝0，⛔ 两码位混）；CR 字节＝0（新字节全 LF）。
- `.out` 扩展名枚数＝**0**（尺＝`Get-ChildItem .scratch/wisp/probes/346 -Recurse | 扩展名=='.out'`；
  证据一律 `.md`，本腿⛔ 建任何零字节件）。
