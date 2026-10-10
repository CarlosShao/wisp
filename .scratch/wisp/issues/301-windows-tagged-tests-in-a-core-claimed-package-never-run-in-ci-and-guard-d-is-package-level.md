# 票 301 — 一枚包已被某个 scope **认领**，于是它里面 `//go:build windows` 的那一半**永远不会在 CI 上被执行**，而 GUARD D 因为判定单位是**导入路径**而⛔ 看得见这件事：把 `./internal/audio/` 拉进 `windows)` 档是唯一能让它开口的形

**立票**：2026-10-10 13:1x 编排者（来路＝票 300 `AC#5` 那三跳我自己预跑完的结论，账＝`A804`；⛔ 塞进票 300 的 `AC#2` 那一批——那一批是修一枚偏移，这一枚是修一把**门的射程**，两件事的代价面完全不同）
**性质**：★**这是一枚仪器票（改 CI 分档＝改门的射程），⛔ 不是缺陷票**。盘上没有任何东西坏了；坏的是"有人以为有门在管这件事"那句**读反了的话**。
**为什么要紧**：票 300 的落地腿正要往 `internal/audio` 里交一枚 `//go:build windows` 的判据。按本票现量，那枚判据交进来以后**在 CI 上从未被执行过一次**，而 `--scope=census` 那枚门⛔ 会为此变红 ⇒ 下一个读"用例已经进仓了"的人会把一件**只有本机成立**的事当成有持续防护的事。本仓已经吃过三次同形（`winlive` 无 tag 档、真机读数只有本机有、票 33 `AC#13` 那句"②零覆盖"），这一枚是**把第三次的根因收口**。

## 现量（每条都带尺；⛔ 引用前先重跑，行号与枚数一律当快照）

- ★**门的判定单位＝导入路径，⛔ 是用例**（尺＝`sed -n '382,397p' scripts/portable-tests.sh`，HEAD `772ff880`）：NO-SCOPE 那一支只对 `$all` 里每枚**包**问"有没有任一档认领它"，`counts`（`:364` 那行 `go list -f '{{len .TestGoFiles}}/{{len .XTestGoFiles}}'`）非 `0/0` 才记进 `unclaimed`。⇒ 只要包被任一档认领，**它永远不会进 `unclaimed`**，而 `TestGoFiles` 是**按当前 `GOOS` 编得进二进制的那些文件**——windows-tagged 的在 linux job 上根本不在分母里。
- ★**`internal/audio` 已被 `core` 认领**（尺＝`grep -n 'internal/audio' scripts/portable-tests.sh` ⇒ `core` 档清单里确有 `./internal/audio/...`），而 `core` 那档跑在 `ubuntu-latest`（尺＝`sed -n '402,403p' .github/workflows/ci.yml` ⇒ `test-core:` ＋ `runs-on: ubuntu-latest`）。⇒ **本票那枚新用例落进 `internal/audio` 之后，CI 上没有任何一步会跑它。**
- **HEAD 上 `internal/audio` 的 windows-tagged 名册**（尺＝对 HEAD blob 逐枚 `git show HEAD:<path> | head -3 | grep -c 'go:build windows'`，⛔ 对工作树量——本票起手时那棵工作树里正躺着票 300 落地腿的未提交件）：**七枚 `_test.go`，其中 2 枚 tagged**＝`capturelevel_windows_test.go`／`hotplug_test.go`。⚠ 后者**文件名不带 `_windows` 而里面有 tag** ⇒ 按文件名分类的那把尺会读错这一枚（我已经踩过一次）。
- **`windows)` 档今天不收 audio**（尺＝`sed -n '248,253p' scripts/portable-tests.sh`）：十枚路径逐枚读过＝`./internal/proc/ ./internal/secret/ ./internal/config/ ./internal/risk/ ./internal/ball/ ./internal/perm/ ./internal/plugin/ ./cmd/llmrecord/ ./internal/session/ ./internal/projctx/`，⛔ audio。它的 `pinned=$win_pin`（`:254`），`win_pin` 正文在 `:193` 起、**十枚导入路径、逐枚与上面那张清单同形**。
- **调用点**（尺＝`grep -n 'scope=' .github/workflows/ci.yml`）：`--scope=census` 只有 `:744` 一枚（`:658` 那两行注释逐字写着 "the ONLY call site GUARD D has"）；`--scope=windows` 只有 `:780`，它所在的 job＝`test-windows:`、`runs-on: windows-latest`（**托管**，⛔ 机主这台开发机——这条前提我 10-08 写错过一次并已更正，见台账）。
- ⚠**仓里已经有一行是按"audio 会在 windows 档"那个形状写的**（尺＝`sed -n '593p' scripts/portable-tests.sh`）：夹具 ledger 逐字起头 `TestLiveWasapiSmoke|./internal/audio/|windows|fixture|`，其后那句理由里还带着一枚**行号锚** `hotplug_test.go:527`。⇒ 这一枚今天到底 inert 与否＝本票 `AC#0` 的一格必答，⛔ 现在当成已知。
- ⛔**"改一枚清单"不是零成本**（这句写给派单，也写给下一个读的人）：把 `./internal/audio/` 拉进 `windows)` 档，会同时把该包**其余** windows-tagged 用例（上面那 2 枚，以及票 300 正要交进来的那一枚）一起拉进 `test-windows` 那个 job 的分母——每一枚都可能带自己的前置条件（真麦克风、活设备、`WISP_LIVE_MIC=1`）。代价枚数＝`AC#0` 的产物，⛔ 由本票票面替它填数。

## 要建什么（⛔ 一进来就动清单；先量代价，再决定落不落）

- [x] **`AC#0` 只读代价普查（这一格⛔ 任何产码，⛔ 任何 go 编译面）**：把三张表交回——① **会被拉进分母的枚数与名字**：对 HEAD blob 逐枚列出 `internal/audio` 里 windows-tagged 的**顶层用例名**（尺要写清是 `grep -n '^func Test'` 还是 `go list`，后者算编译面＝本票⛔）；② **每一枚的前置条件**：它读不读 `WISP_LIVE_MIC`、要不要真设备、会不会在托管 runner 上因缺音频端点而**红**（⛔ 把"skip"当"绿"——按本仓既有口径，`skipped` 与"从未求值"是两回事）；③ **`:593` 那行 ledger 今天到底 inert 与否**（判据＝它声称的那个用例在当前档下有没有被任何档跑到）。＋**`:593` 里那枚行号锚 `hotplug_test.go:527` 在 HEAD 上指向的是不是那件事**（锚腐烂先例＝票 255 `AC#6` 那一族）。判据＝每条都带"哪把尺＋射程目录＋blob 还是工作树"，⛔ 裸数。
- [x] **`AC#1` 落地形（只在 `AC#0` 交出代价、并经编排者裁"落／不落"之后才开工）**：**同一笔 commit** 里改两处——`windows)` 档清单加 `./internal/audio/` ＋ `win_pin` 补上 `github.com/CarlosShao/wisp/internal/audio`。⚠ 那句要求的原文**跨两行**、逐字是这样（尺＝`sed -n '417,418p' scripts/portable-tests.sh`，两行各自都以 `echo "portable-tests.sh: ` 起头）：`:417` 尾部 `… Pull the package` ／ `:418` 开头 `into a named scope and update that tier's pin in the SAME commit, or` ⇒ 引用时⛔ 把它写成一行"逐字"。⛔ 拆成两笔＝第二笔会让 GUARD C 在中间态红；⛔ 顺手改别的档的清单或别的档的 pin。
- [x] **`AC#2` 让它真的开口（反恒真那一格）**：交付要包含一次**同一枚用例在两种档下各跑一发**的对照读数＝改清单前该用例在 CI 语义下从未被求值／改清单后被求值（本机可复跑的最便宜形＝`bash scripts/portable-tests.sh --scope=windows` 前后各一发，名册逐名作差，具名新增枚数）。⚠⛔ 拿 `go test ./internal/audio/` 本机跑绿当这格的凭据——那证的是本机，⛔ 是"CI 现在会跑它"。
- [x] **`AC#3` 门的盲区要不要也收口（本格只裁**形状**，⛔ 默认必做）**：`AC#1` 之后，"某包里 windows-tagged 的判进不了任何档"这一类**仍然**无人管——因为 GUARD D 天生只到包级。要不要给它一枚用例级的钉，是**比本票大一枚**的射程。交付＝一份两形代价表（甲＝不动门、把"哪些包里有 tagged 用例"钉成一枚名册断言；乙＝门扩射程），交编排者裁。⚠ 若腿自行判断"顺手把门改了"＝越界，停手上报。
- [ ] **`AC#4` 门禁与越界**：`sh scripts/d22scan.sh` rc=0；`bash scripts/portable-tests.sh --scope=census` rc=0（⚠ 它的 `STALE` 腿只列表、⛔ 计退码，判定看名册⛔ 看颜色）；`--scope=windows` 与 `--scope=core` 改前改后各一发、逐名作差**具名新增红 0 枚**（⛔ 只报枚数不报名字）；`git show --stat` 名册只含 `scripts/portable-tests.sh` ＋ `probes/301/**`（`AC#1` 那笔⛔ 碰 `.github/workflows/ci.yml`——清单在脚本里，动 yml 就是第二枚射程）；`frontend/**`／`design/**`／三枚冻结件／golden／`thresholds.go`／`tools/d22scan/allowlist.txt`／D43 表零字节；每把门禁件自落一行 `rc=N`（⛔ 0 字节＝那格没交）；⛔ 零 push、commit 必带显式 pathspec。
- [ ] **`AC#5`（2026-10-10 13:3x 编排者追加，来路＝只读腿 `301-a1` 派单必答⑥第 ④ 条；⛔ 不改上面各格原句）把 `ci.yml:400-401` 那两行过期注释改成说实话**：blob 上逐字是 `  # Windows-only packages (ball GUI, cgo speech) are out of this job's scope` ＋ `  # by platform, not skipped: they run in test-windows / slo jobs.`。★现量：对 `internal/audio` 这两句**今天两条都⛔ 成立**——它⛔ 在 `windows)` 档（`:248-253` 十枚路径 `grep -c audio`＝0），也⛔ 在 slo 那两档里跑该包。⇒ 这两行是 `A804` 那条"门只到包级"的**同一根、长在另一枚文件里**。**次序硬约束＝⛔ 现在动它**：那两行该改成什么，取决于 `AC#1` 落不落地（落了＝"runs in test-windows" 半句成真、"slo jobs" 那半句仍⛔ 真），所以**必须按在 `AC#1` 落地并核过之后**。判据＝纯注释面、零行为（`git diff -U0` 逐行看，⛔ 一行可执行 yaml 都不许多变）；⚠ **单独一笔 commit**，pathspec＝`.github/workflows/ci.yml` ＋ `probes/301/**`（⛔ 与 `AC#1` 那笔合并，那笔被 `AC#4` 逐字禁碰 yml）；文案里⛔ 写会被自己扫的现在时计数；⛔ 零 push。

## 边界与已知禁区（⛔ 派单要逐字带上）

- ⛔ **把票 300 的判据搬到 `cmd/wisp` 来"让它进 CI"**：`parseWaveFormat` 是**包内未导出**函数，判据只能长在 `audio` 包里。这条路物理上⛔ 存在。
- ⛔ **为变绿放宽任何断言**；⛔ 动 SLO 阈值／golden／`internal/observe/thresholds.go`／`tools/d22scan/allowlist.txt`／D43 转移表／C1–C32／D1–D47（改契约＝人工批准）。
- ⛔ 动 `.gitattributes` 或任何 git 配置（AGENTS.md「NEVER update git config」逐字压着）。
- ⚠ **凭据只有本机**这条要逐字进派单：票 300 那条判据的落地面本票⛔ 承诺修好；修好之前，任何"CI 会兜着"的假设都是读反了（`A804`）。
- ⚠ **在飞的票 300 落地腿握着的两枚面**：⛔ 第二枚走 Go 编译面的腿；⛔ 任何人动 `internal/audio/**`。本票 `AC#0` 只读、⛔ 编译面 ⇒ 可以在那一枚在飞时开工；`AC#1`⛔。

## 排程

- 现在可派＝`AC#0`（只读、⛔ go 编译面、⛔ 写任何跟踪文件，写面＝只新建 `.scratch/wisp/probes/301/**` 下的 `.md`）。
- `AC#1`／`AC#2` 按在＝票 300 落地腿交回、`300-v2` 裁完 `AC#2` 之后（同一枚包、同一枚写面，⛔ 并发）。
- `AC#3` 按在＝`AC#0` 的代价表到手之后由编排者裁形。
- 与票 255 的关系＝**同族⛔ 同一枚**（255 那族是"档位登记表与它的名册口径同源"，本票是"包的档归属与用例能不能被执行"）；`AC#0` 里那枚行号锚腐烂一格，归口票 255 `AC#6` 的既有形，⛔ 本票新造判据。
- ⛔ 零 push（机主从未授权）。

## 编排者裁定（2026-10-10 13:3x，来路＝只读普查腿 `301-a1`，两笔 `c2853233`→`8c29a526`；件 `.scratch/wisp/probes/301/a1/`＝三枚 `.md` 我逐枚现量＝`00` 1,179／`01` 23,307／`02` 4,215，**0 字节＝0 枚**；逐笔 `git show --name-only` 越界 **0 枚**。★下面每一条承重读数都是我在 blob 层 `772ff880` 自己复跑过的，⛔ 引腿的数）

- ★**`AC#0` 我翻勾**，⚠ 具名射程＝本格判据是"三张表交回并带尺"，**表② 里"这八枚求值并 PASS"那一半⛔ 在本格、归 `AC#2` 的真跑**（腿全部〔读码推的〕，它自己第 ④-2 条具名报了这一点，我认）。
- **名册（我复跑＝对上了）**：`git ls-tree --name-only 772ff880 internal/audio/` ＋ 逐枚 `git show | head -3 | grep -c 'go:build windows'` ＋ `grep -c '^func Test'` ⇒ **7 枚 `_test.go`／2 枚 windows-tagged／顶层用例 9 名**（`capturelevel_windows_test.go` 1 ＋ `hotplug_test.go` 8）。★**腿顶我第 ⑥-1 条成立、我按它改读法**：本票"现量"节里那句"上面那 2 枚"指的是**文件**，谁按**用例**读就把代价少算 **7 枚**——原句⛔ 改，以本节为准（同一枚毛病＝`A791` 名册口径那一族：⛔ 写清"文件还是用例"的尺就会造出一枚假小的代价数）。
- **锚稳定（＝我在 `A806` 欠自己的那一格，腿先做了、我又复一遍）**：`git diff --stat 772ff880..8c29a526 -- internal/audio scripts tools/d22scan internal/proc internal/memory` 输出**空**（0 行）⇒ 两锚下我量过的每枚 blob 逐字相同；腿在 13:2x 于新锚重跑表① 得同一张（7／2／9），与 13:1x 那张对得上。⚠ 两枚在飞面（` M internal/audio/wasapi_windows.go`、`?? internal/audio/parse_wave_format_300_windows_test.go`）13:30 我现量仍躺原处、⛔ 被任何人的 commit 顺走、⛔ 计入本票任何名册（票 300 那枚新用例⛔ 在 blob 上）。
- ★**`:593` 那行 ledger 我裁＝⛔ 整行 inert（腿第 ⑥-3 条顶回我的预设，成立）**，并把机制读全：那行有**两条腿**——① `-skip` 那条今天**打空**（`:694` 实调只喂 `"${scope[@]}"`，⛔ `ledger_pkgs`；且 `:639-641` 的 `any | "$goos")` 一支使 windows 档条目在 ubuntu 上被 `continue`）；② **陈旧性那条是活的**（`:651` `listed=$(go test -list "^${name}\$" "$pkg" ...)` 按 **ledger 行自己的包**跑、与 scope 无关，`:652-653` 打不中就进 `stale`）。⇒ 删那行⛔ 只是"少一行注释"，会**摘掉一枚活着的钉**（改了名或删了用例，`test-windows` 不再红）。旁证三枚同表锚我逐枚读到那一行：`internal/memory/concurrent_test.go:186`、`internal/proc/jobscope_windows_test.go:87`、`hotplug_test.go:527` **全部指向 `t.Skip(` 那一行**（`func TestLiveWasapiSmoke` 在 `:525`）⇒ 惯例成立、那枚锚**没漂**、本票⛔ 供一枚腐烂锚（我派单里按"可能漂"给的靶，记我）。
- ★**我自己新量一枚、腿与我都⛔ 写过（它决定 `AC#1` 落地后的真颜色）**：`skip_pattern` 是**从 ledger 现造的**——`:680` 空档时落到 `^portable_tests_ledger_is_empty_on_this_platform$`、`:686` 否则 `^(A|B|…)$`，`:688` 会把样式打进日志。⇒ 把 `./internal/audio/` 拉进 `windows)` 档之后，第 9 枚 `TestLiveWasapiSmoke` 走的是**被 `-skip` 排除＝从未求值**那一形，⛔ 是它自己 `t.Skip` 出 `--- SKIP`；而 `runtests.sh` 那两枚判红我逐行读到（`:98` `if [ "$skipped" -ne 0 ]` ⇒ `:102 exit 1`，文案含 "SKIP is not a pass (ticket 71 AC#3)" 与 "do not relax this script"；`:104-108` 还有一枚"退 0 而零顶层结果＝`-run` 打空"的判红）。⇒ **`AC#2` 的作差应当＝新增被求值 8 枚、第 9 枚以"未求值"出现而门⛔ 红**；而**只要谁在那一笔里顺手删了 `:593`，第 9 枚就真跑出 `--- SKIP` 并把 `test-windows` 判红**＝腿第 ⑥-3 那条警告的形状，我复跑机制后确认成立 ⇒ **`AC#1` 派单必须逐字带"⛔ 动 `scripts/portable-tests.sh:593` 那一行"**。
- **另两枚数（我复跑）**：`WISP_LIVE_MIC` 在 `internal/audio` blob 上 **3 处命中、全在 `hotplug_test.go`**，在 `ci.yml` blob 上 **0**（CI 从不设它）⇒ 第 9 枚在 CI 上必走 skip 支，与其 ledger 理由那句"hosted runner 无音频端点"同向；`-short` 在 `portable-tests.sh`／`ci.yml`／`runtests.sh` 三份 blob 里 **0／0／0**，而 `TestPinnedThreadStable10s` 的出口是 `testing.Short()`（`hotplug_test.go:444-446`）⇒ 那一枚的 skip 支今天永不触发，进档就实打实占 **10 秒**（代价＝时长，⛔ 缺硬件）。⚠ 腿第 ⑥-2 条也成立：`TestEndpointsPairQueryable` 名字带 `Endpoints` 而 `:428-437` 全用 `fakeWatcher/fakeStream/fakeOpener` ⇒ **派单里按用例名猜前置条件⛔ 是取证**（这一条靶是我给的，记我）。
- ★**裁 `AC#1`＝落**（判据＝本票要买回的可见性只有这一条形，代价已从"⛔ 多少枚"读到"9 枚、枚枚具名、其中 8 枚设备无关"）；**⛔ 现在开工**：`AC#1`／`AC#2` 都要跑 `--scope=windows` 的改前改后名册，而那一跑会把在飞的 `300-r1` **未提交的两枚面**吞进分母（同一枚包！）⇒ **排程硬约束＝`300-r1` 交回 → `300-v2` 裁票 300 `AC#2` → 才轮到 `301-r1`**；撤销口令＝**「301 不落 audio 档」**（落／不落由我裁，⛔ 由腿裁）。
- **`AC#3` 我裁＝本格⛔ 动门**（腿第 ④-5 条自报的越界风险我认它的处置）：甲／乙两形代价表**此刻⛔ 派腿**——编队我先只留 `300-r1` 一枚，等它交回再排。
- **腿第 ⑥-4 条交给我的一枚过期注释，我裁成 `AC#5`**（⛔ 塞进 `AC#1` 那笔：`AC#4` 逐字禁那一笔碰 `.github/workflows/ci.yml`；而我量完才敢说它多过期：`ci.yml` blob 上 `internal/audio` **0 命中**、两枚 slo job 各自只跑 `scripts/slo-check.ps1`（`:825` smoke／`:888` full），而那份脚本里 `go test` 与 `audio` **各只有 1 处命中、两处都在注释里**——`:336` 那枚 `go test` 说的是"别的腿正在跑"、`:15` 逐字写着 `(memory/handle subset, no audio)` ⇒ slo 那两档对本包零执行、脚本自己也这么声明 ⇒ "they run in test-windows / slo jobs" 那半句对本包**两条都假**）。
- next＝等 `300-r1`；`301` 这一列往下＝`301-r1`（`AC#1`＋`AC#2` 同批，带上面三句逐字禁区：⛔ 动 `:593`／⛔ 拆成两笔／⛔ 碰 `ci.yml`）→ `301-v1` 非实现者裁 → `AC#5` 单独一笔。我自己名下仍欠三发（真机 `GetMixFormat` hexdump、照 `293-v1` 含 `cmd/wisp` 重跑坐实 677、`AC#2` 那一发我复核腿的名册作差），⛔ 一枚都没跑。⛔ 零 push。

## 编排者裁定（2026-10-10 14:0x，`AC#3` 那一格裁完；来路＝只读普查腿 `301-a2`，两笔 `4eb29312`→`849ce6e9`，件 9 枚＝`probes/301/a2/**`（最大 62,192 字节、**0 字节＝0 枚**、逐笔越界 0 枚）；★它照做了 `A807` 那条**双锚**规矩＝起手 `86d89478`／porcelain 2 行 ＋ 交回 `41329475`／porcelain 0 行，**两把读数都留**）

- **裁＝本格⛔ 落任何一形，整枚射程搬到票 111 新加的那格 `AC#12`**（撤销口令**「301 AC#3 改回本票落」**）。三条具名理由：
  1. **甲形"认领断言"那一子形会当场造出 4 枚红**（`AC#1` 落地后＝`internal/agent`／`internal/memory`／`internal/models`／`internal/tools`，15 枚文件／32 枚用例），而把这四枚塞进 windows 档**必须先逐枚量它们 tagged 用例的前置条件**（要不要真设备／会不会因缺音频端点在托管 runner 上红）＝**另一枚普查**，⛔ 本票 `AC#1` 那句"⛔ 顺手改别的档的清单或别的档的 pin"允许的形状。
  2. **甲形"登记钉"那一子形今天 0 红**（新增用例不登记才咬＝只咬未来那笔）⇒ 它**⛔ 管那 32 枚**；把这一形记成"盲区已收口"＝造一枚在跑而⛔ 管事的灯（`feedback-evidence-teeth` 第 8 条那一族）。它作为**第一步候选**留在 111 `AC#12`，票面必须逐字带着"它今天 0 红"这句。
  3. **乙形⛔ 是"包级→用例级"，而是"包级→用例级＋档级 OS 语义"**：硬缺的输入＝**tier→runner-OS 映射今天⛔ 是数据、只是注释**（腿尺＝`portable-tests.sh` 里 12 处 OS 字样命中全在注释与 ledger 理由串、代码 0 处；我⛔ 复跑，标〔仅腿量〕）。它一落地同时撞三处——**GUARD A**（`:565-583`；同一件事两枚门给相反颜色：A 绿／乙红）、**ledger 的 `fixture`／`reexec`**（`:585-601`，"从未求值"在本仓既是**要被追捕的事**又是**明文批准的状态**，与 `tools/d22scan/runtests.sh:98`→`:102` 的"SKIP 即红"正面相遇）、**票 111 那句包单位**（`:413` ＋ `:406` totals 行 ＋ `ci.yml:688` 的 35，改它要连载体案例 15/16 一起）⇒ **归属＝GUARD D 的出处票＝111**，⛔ 让本票长出第二座真相源。
- ★**复跑后⛔ 采与采的**：⛔ 采纳表① 第 3 列那枚预设（"完全⛔ 在任何档＝0 枚"确实⛔ 有产出，但⛔ 是本票的账）——我自己 `sed -n '193,205p'` 逐枚读 `win_pin`（10 枚：`cmd/llmrecord`/`ball`/`config`/`perm`/`plugin`/`proc`/`projctx`/`risk`/`secret`/`session`）＋`sed -n '244,256p'` 读 `windows)` 档清单（10 枚）⇒ **audio 确实只在 core 一侧**，本票要修的那枚成立、`AC#1` 落地形一字⛔ 改。采的＝★真盲区在"只被 ubuntu 档认领"那一列＝**5 枚包／17 文件／41 用例**（audio 出局后 4／15／32）。⚠**未复跑具名**＝census 的 `unclaimed` 读数（0 枚）与"只有散文"那把尺，两枚都要 `go list`＝编译面，`300-v2` 在飞 ⇒ 一律标〔仅腿量，编排者未复跑〕，归 111 `AC#12` 的落地腿自跑。
- ⚠**派单措辞错在我（⛔ 记腿）**：我让腿去 GUARD C 块找"同一笔 commit"那两句，而它们长在 **GUARD D 的退码文本**里（`:407-423`，`:409` 起头就是 `GUARD D - $guardd package(s)…`；GUARD C 自己的同义句在 `:556-558`）。票面 `AC#1` 引的那两行逐字**我 `sed -n '416,419p'` 复跑对上**（`:417` 尾＝`Pull the package`／`:418` 头＝`into a named scope and update that tier's pin in the SAME commit, or`，两行各以 `echo "portable-tests.sh: ` 起）⇒ **票面一字⛔ 改**；`AC#1` 那句"⛔ 拆两笔＝第二笔会让 GUARD C 在中间态红"也⛔ 被推翻（**引文归属≠判红者**：清单与 pin 分家确实由 C 判红）。
- **`300-r1` 那枚新用例算不算本票 `AC#0` 名册的第 10 枚——裁＝算**（尺＝`41329475` 上 `internal/audio` ＝3 枚 tagged 文件／10 枚顶层用例，多的那枚＝`parse_wave_format_300_windows_test.go`）。⇒ **`AC#1` 落地腿做 `AC#2` 前后作差时分母按 10 枚**，⛔ 拿已裁的 9 枚旧数去对（"名册口径"那一族已第五次：`A791`／`A799`／`A803`／`A805`／本条）。
- **排程（改）**：⛔ 现在派 111 `AC#12`。硬序＝`300-v2` → `301-r1`（`AC#1`＋`AC#2` **同一笔**：档清单加 `./internal/audio/` ＋ `win_pin` 补 `internal/audio`；⛔ 拆两笔、⛔ 碰 `ci.yml`、⛔ 动 `portable-tests.sh:593`）→ `301-v1` → 票 300 `AC#5` → **111 `AC#12` 的普查腿**（先量那 4 枚盲包各自 tagged 用例的前置条件）→ 111 `AC#12` 落地腿。⛔ 零 push（机主从未授权）。

## Progress log

（追加式；每条 `- [YYYY-MM-DD HH:MM:SS +08] agent=… did=… next=…`，钟点由 `date` 的 stdout 插值，不手打。）
- [2026-10-10 13:22:37 +0800] agent=编排者 did=立票 301 并派 `301-a1`（只读、⛔ go 编译面、写面＝只新建 `probes/301/a1/**`）⇒ 本票 `AC#0` 开工；票面"现量"节全部尺子逐条重跑在 HEAD `772ff880`（⛔ 引工作树——此刻工作树里躺着票 300 落地腿的未提交改动 `wasapi_windows.go` ＋未跟踪新件 `parse_wave_format_300_windows_test.go`，按工作树量会把在飞的活当既有事实）＝`ci.yml:402-403`／`:658`／`:744`／`:780`、`portable-tests.sh:241`／`:248-253`／`:193`／`:364`／`:382-397`／`:413`／`:417-418`／`:593` 逐枚 sed 读到那一行；两枚数＝HEAD 上 `internal/audio` 七枚 `_test.go` 里 2 枚 windows-tagged（`capturelevel_windows_test.go`／`hotplug_test.go`★文件名不带 `_windows` 而内有 tag）、`windows)` 档清单与 `win_pin` 今天各 10 枚且逐枚同形；★`AC#1` 那句"同一笔 commit 改两处"的原文要求跨两行，票面已按两行抄并具名标⛔ 写成一行"逐字"（今天刚为同类事裁过一把）；`AC#3` 只裁形状、⛔ 默认必做，腿若自判"顺手把门改了"＝越界停手；归口关系＝与票 255 那族**同族⛔ 同一枚**，`:593` 里那枚行号锚 `hotplug_test.go:527` 我⛔ 量、写成 `AC#0` 必答 next=`301-a1` 三张表交回后我复跑并对拉尺子口径，再裁 `AC#1` 落不落；⛔ 翻任何框（本票 `AC#0`..`AC#4` 五格全未勾）；⛔ 零 push
- [2026-10-10 13:34:25 +0800] agent=编排者 did=收 301-a1（两笔 c2853233→8c29a526；三枚 .md 逐枚现量 1,179／23,307／4,215、0 字节=0 枚、逐笔越界 0 枚）⇒ **翻本票 AC#0**（凭据＝承重读数我全部在 blob 772ff880 自己复跑：7 枚 _test.go／2 枚 windows-tagged／9 枚顶层用例，尺＝git ls-tree ＋逐枚 git show｜head -3 grep -c go:build windows ＋ grep -c ^func Test；runtests.sh:98 if [ "$skipped" -ne 0 ] ⇒ :102 exit 1 逐行读到；skip_pattern 从 ledger 现造＝:680／:686／:688；WISP_LIVE_MIC 在 internal/audio blob 上 3 处全在 hotplug_test.go、在 ci.yml blob 上 0；-short 三份脚本 0／0／0 而 TestPinnedThreadStable10s 出口是 testing.Short()（hotplug_test.go:444-446）⇒ 进档实占 10 秒）；★本格射程具名＝表② 那半句"八枚求值并 PASS"⛔ 算已证、归 AC#2 真跑；★**腿顶我三处我全认**：①票面"上面那 2 枚"指文件、按用例读会少算 7 枚（原句⛔ 改，以裁定节为准，定式入 A807④"几枚必带单位名"）；②派单表② 把 TestEndpointsPairQueryable 当要硬件的靶是**我给的假靶**（blob:427-438 用 fakeWatcher＋fakeOpener、⛔ 真枚举）；③派单表③ 那句"判整行 inert"射程不足——`:593` 那行有两条腿，`-skip` 那条今天打空（:694 只喂 scope、:639-641 在 ubuntu 上 continue）而 `go test -list` 陈旧性那条**活着**（:651-653 按行自己的包跑），删它会摘掉一枚监管"改名/删用例"的钉（A807①）；★我自己新量一枚补进裁定节：包进了档后第 9 枚走的是**被 pattern 排除＝从未求值**、⛔ 是自己 t.Skip 出 `--- SKIP`（后者才判红）⇒ AC#2 作差应读成"新增求值 8 枚＋1 枚仍未求值而门⛔ 红"；三枚同表锚（memory:186／proc:87／audio:527）我逐枚读到全是 `t.Skip(` 行⇒ 惯例成立、本票⛔ 供一枚腐烂锚、我派单那枚"可能漂"的靶记我；裁 **AC#1＝落、⛔ 现在开工**（撤销口令「301 不落 audio 档」），硬序＝300-r1→300-v2→301-r1，起因＝--scope=windows 那一跑会把在飞 300-r1 未提交的两枚面吞进同枚包分母（13:3x 现量：` M internal/audio/wasapi_windows.go` ＋ `?? internal/audio/parse_wave_format_300_windows_test.go` 仍在、⛔ 被顺走、⛔ 计入本票名册）；AC#3 本格⛔ 动门⛔ 派腿；腿第 ④ 条那枚过期注释我裁成**新格 AC#5**（ci.yml:400-401 逐字称 windows-only 包"they run in test-windows / slo jobs"，对本包两条都假：ci.yml blob 上 internal/audio 0 命中、两枚 slo job 只跑 slo-check.ps1（:825／:888）而那脚本里 go test 与 audio 各只 1 处命中且两处都在注释（:336／:15 逐字 "(memory/handle subset, no audio)"）；次序硬约束＝⛔ 现在动，那两行该写什么取决于 AC#1 落不落地，且**单独一笔**（AC#4 逐字禁 AC#1 那笔碰 ci.yml）） next=等 300-r1；⛔ 其余四格（AC#1..AC#5）一枚不碰；⛔ 零 push
- [2026-10-10 14:11:38 +0800] agent=编排者 did=裁完 `AC#3`＝**本格⛔ 落任何一形、整枚射程搬去票 111 新格 `AC#12`**（撤销口令「301 AC#3 改回本票落」）＋翻 `AC#3` 那一格（凭据＝两形代价表交回、`probes/301/a2/**` 9 枚件、0 字节＝0 枚、逐笔越界 0 枚；★该腿照做了 `A807` 那条双锚规矩＝起手 `86d89478`／porcelain 2 行＋交回 `41329475`／porcelain 0 行、两把读数都留）；三条理由＝① 甲形"认领断言"当场造 4 枚红（`AC#1` 后 `agent`／`memory`／`models`／`tools`＝15 文件／32 用例）而那四枚要先逐枚量前置条件＝另一枚普查、撞 `AC#1` 那句"⛔ 顺手改别的档的清单或 pin"② 甲形"登记钉"今天 0 红＝它⛔ 管那 32 枚，记成"已收口"＝造一枚在跑而⛔ 管事的灯 ③ 乙形的真实射程＝包级→用例级**＋档级 OS 语义**，一落地同时撞 GUARD A（相反颜色）／ledger `fixture`+`reexec`（"从未求值"既是追捕对象又是批准状态，与 `runtests.sh:98`→`:102` 正面相遇）／票 111 那句包单位（`:413`＋`:406`＋`ci.yml:688` 的 35，改它要连载体案例 15/16）⇒ 归属＝GUARD D 的出处票＝111，⛔ 本票长出第二座真相源；★另裁一格＝`300-r1` 新那枚 `TestParseWaveFormatSubFormatOffset300` **算** `AC#0` 名册的第 10 枚用例（尺＝`41329475` 上 audio＝3 枚 tagged 文件／10 枚用例）⇒ `AC#1` 落地腿做 `AC#2` 前后作差时分母按 10⛔ 9（名册口径那一族第五次）；⚠ 派单措辞错在我⛔ 记腿＝我把"同一笔 commit"那两句标成 GUARD C，`sed -n '405,423p'` 现量它们在 **GUARD D 的退码文本**（块 `:407-423`），GUARD C 的同义句在 `:556-558` ⇒ 票面 `AC#1` 那两行逐字我 `sed -n '416,419p'` 复跑对上、⛔ 改一字，"⛔ 拆两笔会让 GUARD C 在中间态红"⛔ 被推翻（引文归属≠判红者）；我自己复跑坐实的三枚＝`win_pin` 10 枚逐枚读⛔ audio、`windows)` 档 10 枚⛔ audio、`c26_seam_posix_125_test.go:1` 逐字 `//go:build !windows`；⚠ 未复跑两枚具名＝census `unclaimed` 读数 0 枚与"tier→runner-OS 只在注释"那把尺（都要 `go list`＝编译面，`300-v2` 在飞）⇒ 一律标〔仅腿量，编排者未复跑〕归 `AC#12` 落地腿自跑；⛔ 动门、⛔ 动别的档、⛔ 零 push next=`300-v2` → `301-r1`（`AC#1`＋`AC#2` 同一笔，分母按 10）→ `301-v1` → 票 300 `AC#5` → 票 111 `AC#12` 的普查腿

---

## 编排者补量（2026-10-10 14:3x，`301-r1` 在飞期间；尺与逐字读数见 `.scratch/wisp/probes/301/orch/2026-10-10-baseline-and-delta-prediction.md`）

⛔ 改上面任何一句原话。这一节只做三件事：把两枚我写错的/写小了的地方就地更正、补一枚本票⛔ 有过的现量、把"腿交件时我拿哪三行去对"写死在**它交件之前**。

1. **★更正 `AC#1` 的代价面（那句话是我自己写的，记我）**：现量节末条写的是"拉进 `windows)` 档会把该包**其余 windows-tagged 用例**一起拉进分母"。**分档的单位是包**，所以进分母的是 `./internal/audio/` 的**全部 42 枚顶层用例**（尺＝HEAD blob 逐枚 `grep -c '^func Test'`，rc=0：8 枚文件＝tagged 3 枚文件 10 用例／无 tag 5 枚文件 32 用例）。⇒ 本票真正在赌的是那 **32 枚无 tag 的**——它们在 CI 上到今天只⛔ 在 ubuntu 的 `core` 档里被求值过，**在 windows 平台上从未被求值过一枚**。`AC#0` 那三张表里"前置条件"那一列射程窄了 3 倍，后续引用⛔ 再按"只有 tagged 有前置"读。
2. **★新现量：`-skip` 掉的用例⛔ 产 `--- SKIP` 行 ⇒ `AC#1` 落地⛔ 会因为 ledger `:593` 那行把 `--scope=windows` 打红**。证据＝两枚归档 CI 日志（射程＝仓外 `.scratch/ci-logs/`，托管 windows runner 真跑出来的字节）：run `37158259050`（10-03）windows 档 `RUN=488 PASS=335 FAIL=12 SKIP=1`，run `37166458550`（10-04）`RUN=492 PASS=339 FAIL=12 SKIP=1`——两发的 `-skip` 模式串都含 7 枚名字，其中 `TestHelperProcess`（`./internal/proc/`）与 `TestSyncRegistryProbeLive`（`./internal/risk/`）**枚枚都在那一档的 scope 里**；若过滤会产 SKIP，这两发至少该 `SKIP=3`。实测那一枚 SKIP 的名字与出处逐字＝`--- SKIP: TestSyncRedTeamRealOneDrive` ＋ `syncdirs_redteam_windows_test.go:220`（它⛔ 在 ledger 里，是自己 `t.Skip` 的）。⇒ `TestLiveWasapiSmoke` 进档后走的是**静默过滤**那一条，与 `runtests.sh:99-104` 那句"SKIP is not a pass"⛔ 相干。
3. **落地后要盯的两枚具名计时面**（⛔ 判定、只是"红了先看哪两枚"）：`hotplug_test.go:443` `TestPinnedThreadStable10s` 只在 `-short` 下跳，而 `portable-tests.sh:694` 那条命令从⛔ 传 `-short`（尺＝`grep -n -- '-short' scripts/portable-tests.sh` ⇒ 0 命中）⇒ **落地一次＝那一档多 10 秒墙钟**，且它判的是 `LockOSThread` 独占性（托管 runner＝共享机器）；同族＝`capturelevel_windows_test.go` 的 `waitForLevels(..., 5*time.Second)`。⚠ 载入形风险（`cmd/wisp` 缺 sherpa DLL 那枚 `0xc0000135` 旧例）在这里⛔ 适用：尺＝8 枚测试文件的 import 块逐枚读，rc=0，全纯 Go、⛔ cgo、⛔ 设备依赖，唯一包外依赖＝`internal/observe`。
4. **我在腿交件之前写死的 Δ 预测**（对 `AC#2`／`AC#4` 那两发本机 `--scope=windows`；⛔ 读数、是判据）：**Δ `=== RUN` ＝ +41**（42 减被过滤的那枚 `TestLiveWasapiSmoke`；报 +9/+10 ＝ 量的是 tagged 那一半，⛔ 是这一档的真实分母；报 +42 ＝ 那一发里 `-skip` 根本没生效，具名上报）；**Δ `--- SKIP` ＝ 0**（若 Δ SKIP≥1 且名字来自 `internal/audio`，则第 2 条被推翻，归因＝"SKIP 是红"那一形，⛔ "audio 的判据坏"）；**Δ `--- PASS` ＝ 41 −（逐名报出的新增红）**，⛔ 枚数一致当"没问题"。
5. **一枚对上了的复跑**：`301-a2` 表①里 `internal/audio` 那行的 total `_test.go`＝8、tagged＝3，我这一枚用同一把 `git ls-tree`＋`head -3` 尺独立复认（8 枚文件／3 枚 tagged／42 枚用例）⇒ **这一格从〔仅腿量，编排者未复跑〕销账**。仍未销的两格⛔ 变：census 真读数、`--scope=core`／`--scope=windows` 的本机颜色（都要 go 面，此刻归 `301-r1`）。
- [2026-10-10 14:3x +0800] agent=编排者 did=在飞补量一节（本件最上那节「编排者补量」）：更正 `AC#1` 代价面＝分母粒度是**包**（audio 全部 42 枚顶层用例，其中 32 枚无 tag 从未在 windows 上被求值过）＋新现量 `-skip` ⛔ 产 `--- SKIP`（两发归档 CI 字节 488/492、SKIP=1 且那枚 SKIP 有名有 file:line）⇒ 落地⛔ 因 ledger `:593` 变红＋写死 Δ 预测三行 next=核 `301-r1` 交件的三行 Δ；⛔ 零 push；⛔ 零翻勾
- [2026-10-10 14:5x +0800] agent=编排者 did=收 `301-r1`（六笔 `74eb032c`→`3a5771ce`）＋纯文本尺复跑对上：windows A-evaluated 399→**440**（＋41＝我 14:3x 写死的预测逐字命中）、新增 41 枚枚枚 ∈ audio 那 42 枚名册、42 里唯一⛔ 被新求值的＝`TestLiveWasapiSmoke`（静默过滤，两个独立方向各证一次）、消失 0 枚、core 两发名册与 four numbers 逐字相同、census 差异只一行（audio `core`→`corewindows`）、九枚 rc 件全非 0 字节、六笔里 `ci.yml`／`internal/**` 各 0 枚、票框⛔ 一枚被它自勾。★我自抓两枚：Δ 预测式少写"摘除"那一侧（实测 A-pass ＋42 因 `TestResolvePerCallBudget` 红→绿＝墙钟噪声、⛔ 读成成绩）；`A811④` 那枚归属机制盘上判不了（两说产同名册）⇒ 降级成纵深防御、AGENTS.md 规则 1.4 那枚违反记在腿自己头上。⛔ 翻勾（`AC#1`／`AC#2`／`AC#4` 全交 `301-v1`，含 GUARD C 反形那发）；⛔ 零 push

---

## 编排者收 `301-v1`（2026-10-10 15:3x，非实现者判语表＝`.scratch/wisp/probes/301/v1/50-verdict.md`，凭据＝`probes/301/v1/{logs,rosters}/**`）

判语来源＝`301-v1`（⛔ 实现者、⛔ 本票产码者）。我这边**只采我自己复跑上的那一半**，其余标〔腿裁，编排者未复跑〕。

### 翻两格、留一格

- **`AC#1` ⇒ 勾**。凭据两半：本体（`0a0f62ef` 那一笔两枚 token，我在 blob 层复跑过＝`win_pin` 11 枚／档清单含 `./internal/audio/`）＋★**反形那一发**（把 `win_pin` 里 audio 那枚摘掉 ⇒ `--scope=windows` 当场 `exit 1`，红句逐字含 `Pinned: 10, resolved: 11` 与 `> github.com/CarlosShao/wisp/internal/audio`，且 stdout＝0 字节＝压根⛔ 走到跑测试；还原后正控再跑⇒ `GUARD C` 命中 0、名册与首发逐字相同、`cmp` vs `git show HEAD:` rc=0 ＋ `porcelain -- scripts` 四次全 0 行）。⇒ 这两枚 token **是枚被活门看着的钉**，⛔ 一句写在注释里的愿望。〔反形＝腿跑，编排者未复跑；我复跑的是还原后的树干净＝0 行〕
- **`AC#2` ⇒ 勾**。判语＝按票面逐字预授权的"本机档差"形结，⛔ 把 111 `AC#11`（那格票面自含"CI 真跑过一次"）搬过来。腿用自己的两发复现了改前（A-evaluated **399**）与改后（**440**），`comm -13`＝41／`comm -23`＝0，41 枚枚 ∈ 它独立抽的 audio 42 枚名册，唯一未求值＝`TestLiveWasapiSmoke`（本机复认"静默过滤"）。**★这与我 14:3x 在腿交件前写死的 Δ 预测逐字对上**（＋41／SKIP 0／A-evaluated 440）。
  - ⚠ **本格结掉⛔ 等于"CI 上已经有颜色"**。v1 具名欠我一发：`test-windows`（`ci.yml:516-517`，调用点 `:780`）在**托管** runner 上的首跑；红时先看 `TestPinnedThreadStable10s` 与 `waitForLevels(5*time.Second)` 那族计时。**这发欠账已于 15:2x 由推送启动**（见下"推送"一节），读数回来前本格凭据射程＝本机。
- **`AC#4` ⇒ ⛔ 勾，理由具名**。本格列的子句里"commit 必带显式 pathspec"被实现方**自己违反并自报**（v1 判＝缺陷成立、归写腿、只能靠追加闭）。⛔ `--amend`（AGENTS.md 硬禁）⇒ 追溯⛔ 可能，而"追加闭是否**满足**这格的子句"是一枚判语、⛔ 我给自己盖章能给的东西。⇒ 处理形＝随 `AC#5` 的验收腿（`301-v3`）带一枚必答裁它；在那之前这格留空。其余子句**都已交**：`d22scan` rc=0、`census` rc=0、`windows`／`core` 改前改后逐名作差**具名新增红 0 枚**、六笔名册尺越界 0 枚（相对本票射程）、三枚冻结件／golden／`thresholds.go`／`allowlist.txt` 零字节、九枚（腿）＋12 枚（v1）rc 件全非 0 字节、⛔ 零 push（腿与验收腿都⛔ 推）。

### ★更正本票 `AC#4` 里我写的一句机制话（⛔ 改原句，只追加这一节）

原句逐字＝"⚠ 它的 `STALE` 腿只列表、⛔ 计退码，判定看名册⛔ 看颜色"。**盘上两处都⛔ 对**，我 15:3x 自己在 HEAD blob 上复跑确认：

1. 那枚 stale 腿**长在档路径里、⛔ 长在 census 里**：`--scope=census` 那一支从 `:287` 到 `:425` 打完 totals 就 `exit 0`（GUARD D 命中时 `exit 1`），**永远⛔ 进到** `:627-675` 那个 ledger 循环 ⇒ **"census 的 STALE 腿"⛔ 存在这枚东西**。
2. 活着那支的行为与"⛔ 计退码"**相反**：`:652` `listed=$(go test -list "^${name}\$" "$pkg")` 按 ledger 行自己的包问，`:653-655` 不中则进 `stale`，`:667-674` **`exit 1`＝计退码**。
- ⇒ 后果（写给下一个读这格的人）：`--scope=windows`／`--scope=core` 今天**带着**那枚 ledger 行改名／删用例就当场红的活钉。⛔ 因此放宽 `AC#4` 任何子句，⛔ 因为"看名册⛔ 看颜色"这句话省得一发。⚠ 另一枚连带：`AC#4` 原文那句"⛔ 动 `:594` 那行"（原 `:593`）的理由现在比立票时更硬——**删那行⛔ 只是少一行注释，会一并摘掉那枚活钉**（v1 独立读到同一段码）。
- 归属：这句误述⛔ 是本票造成的（立票时就那样写），⛔ 影响任何一格判定，⛔ 谁的欠账。

### 一枚既有红转绿＝⛔ 是成绩（v1 的决定性读数，我复跑不了它那一发，标〔腿裁〕）

`TestResolvePerCallBudget`（`internal/risk/pathresolver_budget_norace_test.go`，1 ms 预算属⛔ 动面、两枚腿都⛔ 动）：v1 量到**同树同机一分钟之内**`--scope=core` 那发 `--- FAIL (3.90s)`、`--scope=windows` 那发 `--- PASS`，且装回 before-blob 的复现里它本来就是绿的 ⇒ 判语＝"第 5 枚负载敏感读数"，登记为**待归因红（负载敏感计时判据）**，归编排者＋那枚 1 ms 预算的契约面（`internal/risk`／C26）。⛔ 许任何人把它读成"这次改动让它变绿"。

### 推送（2026-10-10 15:2x，机主当场授权）

- 授权原话入台账 `A813`。执行＝`git push origin dev`，`cc315261..bcd0a543`（**741 枚**），fast-forward 已核（`HEAD..origin/dev`＝0），终判据已核（远端 tip 逐字＝本地 HEAD）。⛔ `cnb`（第二远端）＝⛔ 动，仍差 921 枚。推送前预检＝密钥形状 `ghp_`／`AKIA`／`xox*`／`-----BEGIN` 各 0 命中，7 枚 `sk-` 命中全是票名 slug 里的词片段（`-task-scope-…`、`-risk-assessor…`）。⚠ 事实一句：`CarlosShao/wisp` 在 GitHub 上⛔ 是私有，**是公开仓库**（现量 `gh repo view` ⇒ visibility＝PUBLIC）。
- ★这一推顺带把上面那枚 `test-windows` 首跑欠账**启动**了（run `38034689386`）。

### Progress log（本节）

- [2026-10-10 15:3x +0800] agent=编排者 did=收 `301-v1` 判语表 ⇒ 翻 `AC#1`／`AC#2` 两格、`AC#4` 具名留空（子句"显式 pathspec"被实现方自报违反、⛔ 追溯修、裁"追加闭算⛔ 算满足"归下一枚非实现者）＋追加更正本票 `AC#4` 那两句 stale 机制话（我在 HEAD blob 上复跑确认：census `:287→:425` 就 exit、stale 腿在 `:627-675` 且 `:674 exit 1`＝计退码）＋推送执行（`A813`）next=`AC#5` 落地腿（v1 已给七条形状）→ 其验收腿带 `AC#4` 那枚残留必答；`test-windows` 首跑读数
- [2026-10-10 15:5x +0800] agent=编排者 did=收 `301-r2`（`4be5168c`→`b761b584`）＋自己复跑四把尺（名册 7 枚／非注释 `+/-` 行 0 枚／可执行键新增 0 枚／`yaml.safe_load` 正规化 HEAD~2≡HEAD＝True／三发门禁 rc=0 且 census totals 逐字未变）＋★落 `A816`：它顶回我派单 §1 那句"audio ⛔ 在 ubuntu 的可求值集合里"（包级⛔ 对：`core)` 数组含 `./internal/audio/...`、linux 解析 `9/5`；真话只到文件级＝3 枚 windows-tagged 文件）——这是本波第三枚"粒度"错、定式补"说包级还是文件级"；`AC#5` ⛔ 翻（等 `301-v3`，带 `AC#4` 那枚残留必答）next=推 `b761b584` 让 CI 跑这枚 yml；`test-windows` 首跑颜色三次抓取遇 GitHub API TLS timeout、待回填；⛔ 零翻勾

## 编排者 CI 色回填（2026-10-10 16:0x，`AC#2` 那句"欠一枚具名 CI 色"到此销账；正文＝`.scratch/wisp/probes/301/orch/2026-10-10-ci-color-backfill.md`，账＝`A817`）

- ★**这格买到什么**：`AC#2` 已勾靠的是本机名册作差；从今天起**同一台托管 runner 自己开口**——`--scope=windows` 被求值面 A-evaluated **399→440（＋41）**，与我在 `301-r1` 交件前写死的那行 Δ 预测**同数**，而**新增红 0 枚**（两发 `--- FAIL` 名集合 `diff` 空、与 audio 名册 `comm -12`＝0 枚）。census totals 逐字未变、六枚 job 色与基线逐枚相同 ⇒ ⛔ 一枚门的颜色被这次改动挪过。
- ⚠**载具口径（引这节时必须带上）**：origin/dev 停在 `cc315261`（10-06），今天首推一次带 **742 枚** ⇒ 盘上**不存在**"`AC#1` 前一笔"的 CI 读数。本节证的是①「落地后该档真被求值了多少」；⛔ 证②「那 12 枚红属 742 枚里哪一笔」（那要逐笔 bisect，⛔ 本票射程、⛔ 共享工作树该做的动作）。基线那发是「本批之前」，⛔ 是「本票之前」。
- **那 12 枚红⛔ 属本票**（尺逐枚：`--scope=windows` 红名 ∩ `probes/301/r1/rosters/audio-toplevel-names.txt`＝0；`AC#1` 那笔只动 `scripts/portable-tests.sh` 两行，逐笔名册尺见 `301-v1` §3-⑥）。5 枚红句逐字点名 runner TEMP 的 8.3 别名 `C:\Users\RUNNER~1`，2 枚同文件同族但红句不含该串（机制**未证**），5 枚是 `syncdirs_test.go` 的 under-profile fallback 家族（机制**未证**）⇒ **逐枚具名记进 `A817` 待归因，⛔ 塞进本票、⛔ 塞进票 302**（那一族的修法在"判据读哪棵树"，先例＝票 115 `AC#2`／`AC#3`）。
- ★**顺手看见的一枚无关新红已立成票 302**：`cmd/wisp` 那档（`--scope=cli`）`FAIL` 6→8，新增 3 枚＝票 33 那族要**页面自己回话**的 AC#13/AC#14 判据，红句同形 `no report … from the page within 15s (what DID arrive at the door: nothing at all)`；⚠ 同发里窗体那一侧是**建成**的（`embed resolves 1068 entry byte(s)`／`entry bytes carried=true`）⇒ 那句"CI 开不出面板"在本发⛔ 成立，缺的是回执那一跳。**这 3 枚⛔ 是本票造成的**，本票只是看见它的人。票 302 里写死了：⛔ 为变绿放宽任何断言，"没有回执就⛔ 判"必须原样保住。
- ⛔**本格因此翻勾**：`AC#4` 整格仍等 `301-v3`（它裁的是"追加闭算⛔ 算满足那枚显式 pathspec 子句"，与颜色无关）；`AC#5` 亦等 `301-v3`。
- ⚠ **仍欠一枚读数**：`AC#5` 落地笔 `b761b584` 之后那发 run `38035842314` 的 `test-windows` 终态（16:0x 现量＝in_progress）。纯注释面**预期**不变色，⚠ 预期⛔ 等于读数。

### Progress log（本节）

- [2026-10-10 16:0x +0800] agent=编排者 did=回填本票 `AC#2` 欠的那枚具名 CI 色（正文＝`probes/301/orch/2026-10-10-ci-color-backfill.md`，账＝`A817`）：载具＝同机两发 job 日志（基线 `cc315261`／run `38004428359`／job `114069831344`，改后 `bcd0a543`／job `114162576251`，尺＝`gh api .../jobs/<id>/logs` 原始字节全留 `logs/`）；读数＝`--scope=windows` `577/386/12/1` → `624/427/12/1`（A-evaluated 399→440＝＋41，与我腿交件前写死的 Δ 预测同数）、红名集合逐字相同（`diff` rc=0）、与 audio 名册交集 0 枚、census totals 逐字未变、六枚 job 色逐枚同色；★自抓一枚方向错尺（红句明细写在 `--- FAIL:` **之前**，我第一把取了之后 20 行⇒差点把"5/12 含 `RUNNER~1`"写成"半数不含"，重跑后 5 枚含／7 枚未含且机制未证）；★顺手量到 `--scope=cli` 新增 3 枚恒红（票 33 那族 AC#13/AC#14 的"页面回执"判据被 CI 求值）⇒ **立票 302**（归口／放置票，三形含"丙＝不做"，⛔ 为变绿放宽断言），⛔ 算进本票；`internal/risk` 那 12 枚逐枚具名入 `A817` 待归因⛔ 并案；⛔ 翻任何框（`AC#4`／`AC#5` 等 `301-v3`）；⛔ 产码、⛔ 动档、⛔ 动 `:594` ledger 文案 next=`301-v3`（在飞）→ 其交件后我落 `AC#4`／`AC#5` 判语 → 票 302 `AC#0` 只读普查腿
