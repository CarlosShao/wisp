# 件2 — GUARD D 两形代价表（用例级 vs 档级 OS 语义）（腿 111-a6b，HEAD `f76091cd`）

⛔ 推荐、⛔ 裁形、⛔ 翻框。这一张只列**代价**，每格带尺与内容锚（行号一律当快照，重跑尺＝`01-rulers.sh` ＋ `sed -n`／`grep -n`）。
两形共同前提：今天 GUARD D 的判定单位＝**导入路径**（`portable-tests.sh` census 那支，内容锚 `census totals: packages=`；
退码文本内容锚 `GUARD D - ` 那一块，末句 `into a named scope and update that tier's pin in the SAME commit`）。
今天 D 的**唯一调用点**＝`ci.yml` 里 `--scope=census` 那一发（`runs-on: windows-latest` 那个 job），
该步自述口径内容锚：`runs `go env GOOS` and `go list` only, never `go test`, never a build`。

## A. 形一：GUARD D 扩到**用例级**

| 代价项 | 现量 | 尺／锚 |
|---|---|---|
| 要动的文件 | **4 枚**：`scripts/portable-tests.sh`（census 支＋GUARD D 退码文本＋四枚 pin 变成用例 pin 或加一枚 case pin）、`scripts/portable-tests-selftest.sh`（载体）、`scripts/winsec-tests.sh`（它把显式包路径递给 portable-tests.sh，用例级口径要它认）、`.github/workflows/ci.yml`（census 那步的自述文案＋可能的 `-tags`） | `git grep -n -e 'portable-tests.sh' HEAD -- scripts .github` 的全部非注释命中＝**8 处**（`ci.yml` 3 处 `run:`、`winsec-tests.sh`、`portable-tests-selftest.sh` 等） |
| 用例宇宙从哪儿来 | **只有两条路，各有代价**：① `go test -list '.*'`＝**census 那一步从此要编译＋链接 28-35 枚包的测试二进制**，直接顶掉该步自述的"never a build"，并且吃 `cmd/wisp` 的加载墙（`scope=cli` 那支自述：test binary **dies at LOAD time** without sherpa DLLs，ticket 98）；② grep `^func Test` ＝**自己重当一遍构建标签求值器**（要处理 `!race`／多条件 tag／tag 行不在前 3 行的形状），本腿现量这种重当的错判面：文件名尺错 16 枚（§件1-3）、`TAG_OTHER` 那枚 `//go:build !race` 会被"含 windows 字样"的全文尺读反 | census 步自述锚 ＋ `portable-tests.sh` 里 cli 支的注释锚 ＋ `01-file-class.tsv` 的 R1/R1f 对拉 |
| 会不会让现有门变红（今天） | **会，且具名**：4 枚盲包（`agent`/`memory`/`models`/`tools`）里的 **32 枚 tagged 用例**逐枚红；再加票面 ⓑ 那 3 枚 ledger 行（`TestHelperProcess`／`TestLiveWasapiSmoke`／`TestSyncRegistryProbeLive`）——本腿现量：这 3 枚的包**已经**被 windows 档认领，所以红的不是"没人跑"，而是"跑了但 SKIP"：`TestLiveWasapiSmoke` 体内需 `WISP_LIVE_MIC=1`＋物理麦克风（ledger 原文自述），而 `tools/d22scan/runtests.sh` 的规则（`portable-tests.sh` 头部逐字：`zero top-level PASS *and* zero FAIL is fatal`，且 SKIP⛔ 记成 pass）⇒ **同一格既要求"未求值必须追捕"又把 ledger 明文批准的 SKIP 判红**，两态文案要落地腿先答 | `05-ruler-meta.txt` 第 5 节（ledger `platform=` 行全名册：windows 3 行／linux 4 行／any 若干）＋ R4 的 `t.Skip(` 命中名册 |
| 造出的新恒绿面 | **3 枚**：(1) 用例名册一改名/删掉就静默失效——现成先例只有 ledger 那条 staleness 检查（原文锚 `delete the case, and GUARD's staleness check goes red`），census 那支⛔ 有对等物；(2) `-run` 过滤错拼＝匹配 0 枚 ⇒ 靠 runtests 的零 PASS 判红兜，但**兜的是整步红⛔ 是"哪一枚没匹配"**；(3) grep 版宇宙（路②）对 `!race`/多条件 tag 的误判方向⛔ 对称，会既虚增又漏计（本腿已量到 1 枚方向反的：`leg_sink_nail_131_windows_test.go` 文件名带 `_windows` 而体内⛔ 没 tag） | 见上行 |
| 单位口径（票面 ⓒ） | 用例级**必须同批**动的三枚：`census totals: packages=` 那行、GUARD D 退码文本里那句包单位（内容锚 `33 个包里的 20 个`／`20 of 33`）、`ci.yml` census 步注释里的全⛔ 包枚数（本腿现量该注释提到 `35` 枚包）；载体逐字判据 `'census totals: packages='` 在 selftest 里命中 **5 处** ⇒ 改口径＝改这 5 处＋补种子 | `grep -n -e 'census totals' /tmp/ps.sh` ＝ :531 :562 :800 :816 :854 |
| 连带分母（件3 那一枚） | 用例级⛔ 免掉连带：要真跑那 32 枚，包必须进某档 scope，于是 **341 枚无 tag 用例一起进 windows 腿分母**（放大 11.7 倍）；即便用 `-run` 只点 32 枚，**测试二进制仍要整包编译＋加载** ⇒ 包级连带变成"编译/时长"代价⛔ 是"执行"代价，但 skip／加载墙那两枚照吃 | `04-case-preconditions.tsv`＋件1-2 表 |

## B. 形二：GUARD D 扩到**档级 OS 语义**（档认领关系携带"这档跑在哪台 OS"）

| 代价项 | 现量 | 尺／锚 |
|---|---|---|
| 要动的文件 | **2 枚**：`scripts/portable-tests.sh`（新增一张 tier→OS 的**数据**表＋census 判定式＋豁免名册）、`scripts/portable-tests-selftest.sh`（补种子）。`ci.yml` 与 `winsec-tests.sh` ⛔ 必改（判定单位仍是包） | 同上调用点尺 |
| 前置 ⓐ 的真实规模 | 本腿现量 `portable-tests.sh` 全文 **768 行**里 `ubuntu`/`windows` 字样 **25 处**，其中**注释行 15 处**；非注释命中 **10 处**，逐枚看过：`tiers='core windows cli winsec census'`、`case` 分支标签 `windows)`、pin 查找行、**ledger 的 6 行理由串**——⇒ **携带 runner-OS 的代码＝0 处**（票面那句"代码 0 处"在本 HEAD 仍成立）。`core`→ubuntu 只在散文里、`windows`→`test-windows` 只在散文里，`cli`／`winsec` 两档的 runner **连注释都没成文**（本腿是从 `ci.yml` 的 `runs-on` 反推的） | `grep -n -i -e ubuntu -e windows /tmp/pt.sh` ＋ `grep -v '^[0-9]*:[[:space:]]*#'` |
| 会不会让现有门变红（今天） | **会，同样 4 枚包红**（判定单位还是包：`agent`/`memory`/`models`/`tools`），但**⛔ 多红 ⓑ 那 3 枚**（那 3 枚的包已被 windows 档认领，包级看不见用例级才看得见的"跑了但 SKIP"） | R5 联接表 `03-blind-roster.tsv` |
| 镜像那一半要不要收 | 形二天然不对称：只写"windows 用例需一枚会跑 windows 的档"⇒ `cmd/wisp` 的 **1 枚 `!windows` 文件／3 枚用例**继续没人跑（本腿 `BLIND_POSIX` 现量）。若把定义写成对称，**今天就多 1 枚包红**，而那枚包的解法代价极高：`cmd/wisp` 进 ubuntu 档要拖 **136 枚用例**（133 无 tag＋3 POSIX）进 ubuntu 腿，还要撞 `cli` 支自述的加载墙 | R1/R5 ＋ `03-blind-roster.tsv` 里 `cmd/wisp` 那行 |
| 造出的新恒绿面 | **2 枚**：(1) tier→OS 表是**手抄的第二真相源**——今天⛔ 任何仪器对拉它与 `ci.yml` 的 `runs-on`，步骤搬 job ⇒ 表过期后**门重新恒绿**（这正是票面 ⓐ 想消的事，形二把它从散文搬成数据但同时新建了一处可腐烂面）；(2) **豁免名册是一枚新清单**，⛔ 有人查"删一行"——`portable-tests.sh` 头部对 GUARD C 的自述就是这一族的先例（内容锚 `GUARD C bites when a row is deleted from a list that still exists`）⇒ 需要反形守卫，那又是一枚文件内改动 | `05-ruler-meta.txt` 第 3 节调用点尺 |
| 单位口径（票面 ⓒ） | **⛔ 必改**：判定单位还是包 ⇒ `packages=`／`20 of 33`／`35` 三处口径不变，混单位读数自然⛔ 出现；这是形二对形一省下来的最大一块 | 上行 |
| 买到的可见性 | 只买到"这枚**包**的 windows 半边有档跑在 windows 上"；⛔ 买到"那 32 枚被求值"。本腿现量一枚具体残差：`WIN_LIVE` 类 **12 枚文件／29 枚用例**（`cmd/wisp` 7／13＋`internal/ball` 5／16）**包级看是被 windows 档认领的，用例级看从未被执行**（票 111 `AC#11` 落的是 `-tags winlive` 的**编译**门，同一支自述⛔ 主张执行覆盖）——形二对这 29 枚**结构上看不见** | `01-file-class.tsv` 的 `WIN_LIVE` 类 |
| 连带分母 | 与形一同款：豁免⛔ 成立而真收口（把包拉进 windows 档）⇒ 同样 341 枚连带。**两形在这一格上⛔ 有差别**——差别只在"⛔ 用 `-run` 能不能少拖几枚" | 件1-2 表 |

## C. 两形代价对照（只列差值，⛔ 排序、⛔ 结论）

| 维度 | 用例级 | 档级 OS 语义 |
|---|---|---|
| 动几枚文件 | 4 | 2 |
| 今天红的枚数 | 32 枚**用例**（＋ⓑ 那 3 枚要写两态文案） | 4 枚**包**（对称版再＋1 枚包＝`cmd/wisp`） |
| 要不要改单位口径（ⓒ 三处＋载体 5 处） | 要 | ⛔ 要 |
| 要不要新建第二真相源 | 要（用例 pin） | 要（tier→OS 表＋豁免名册） |
| 要不要碰编译面 | 路①要（`go test -list`＝编译＋链接，撞加载墙）；路②⛔ 要但自己重当标签求值器 | ⛔ 要（沿用 census 已有的 `go list` 过滤） |
| 结构性残差（买不到 visibility 的那一半） | ledger 批准的 SKIP 与新断言正面冲突 | `WIN_LIVE` 12 文件／29 用例；`cmd/wisp` 3 枚 POSIX；"包有档＝档跑对 OS"之外的用例级形状 |
| 连带分母 | 341 枚 | 341 枚（同款） |

rc=0  # 本件由 111-a6b 撰写；每格的尺在"尺／锚"列，原始名册在同批 `01`–`05` 件
