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

- [ ] **`AC#0` 只读代价普查（这一格⛔ 任何产码，⛔ 任何 go 编译面）**：把三张表交回——① **会被拉进分母的枚数与名字**：对 HEAD blob 逐枚列出 `internal/audio` 里 windows-tagged 的**顶层用例名**（尺要写清是 `grep -n '^func Test'` 还是 `go list`，后者算编译面＝本票⛔）；② **每一枚的前置条件**：它读不读 `WISP_LIVE_MIC`、要不要真设备、会不会在托管 runner 上因缺音频端点而**红**（⛔ 把"skip"当"绿"——按本仓既有口径，`skipped` 与"从未求值"是两回事）；③ **`:593` 那行 ledger 今天到底 inert 与否**（判据＝它声称的那个用例在当前档下有没有被任何档跑到）。＋**`:593` 里那枚行号锚 `hotplug_test.go:527` 在 HEAD 上指向的是不是那件事**（锚腐烂先例＝票 255 `AC#6` 那一族）。判据＝每条都带"哪把尺＋射程目录＋blob 还是工作树"，⛔ 裸数。
- [ ] **`AC#1` 落地形（只在 `AC#0` 交出代价、并经编排者裁"落／不落"之后才开工）**：**同一笔 commit** 里改两处——`windows)` 档清单加 `./internal/audio/` ＋ `win_pin` 补上 `github.com/CarlosShao/wisp/internal/audio`。⚠ 那句要求的原文**跨两行**、逐字是这样（尺＝`sed -n '417,418p' scripts/portable-tests.sh`，两行各自都以 `echo "portable-tests.sh: ` 起头）：`:417` 尾部 `… Pull the package` ／ `:418` 开头 `into a named scope and update that tier's pin in the SAME commit, or` ⇒ 引用时⛔ 把它写成一行"逐字"。⛔ 拆成两笔＝第二笔会让 GUARD C 在中间态红；⛔ 顺手改别的档的清单或别的档的 pin。
- [ ] **`AC#2` 让它真的开口（反恒真那一格）**：交付要包含一次**同一枚用例在两种档下各跑一发**的对照读数＝改清单前该用例在 CI 语义下从未被求值／改清单后被求值（本机可复跑的最便宜形＝`bash scripts/portable-tests.sh --scope=windows` 前后各一发，名册逐名作差，具名新增枚数）。⚠⛔ 拿 `go test ./internal/audio/` 本机跑绿当这格的凭据——那证的是本机，⛔ 是"CI 现在会跑它"。
- [ ] **`AC#3` 门的盲区要不要也收口（本格只裁**形状**，⛔ 默认必做）**：`AC#1` 之后，"某包里 windows-tagged 的判进不了任何档"这一类**仍然**无人管——因为 GUARD D 天生只到包级。要不要给它一枚用例级的钉，是**比本票大一枚**的射程。交付＝一份两形代价表（甲＝不动门、把"哪些包里有 tagged 用例"钉成一枚名册断言；乙＝门扩射程），交编排者裁。⚠ 若腿自行判断"顺手把门改了"＝越界，停手上报。
- [ ] **`AC#4` 门禁与越界**：`sh scripts/d22scan.sh` rc=0；`bash scripts/portable-tests.sh --scope=census` rc=0（⚠ 它的 `STALE` 腿只列表、⛔ 计退码，判定看名册⛔ 看颜色）；`--scope=windows` 与 `--scope=core` 改前改后各一发、逐名作差**具名新增红 0 枚**（⛔ 只报枚数不报名字）；`git show --stat` 名册只含 `scripts/portable-tests.sh` ＋ `probes/301/**`（`AC#1` 那笔⛔ 碰 `.github/workflows/ci.yml`——清单在脚本里，动 yml 就是第二枚射程）；`frontend/**`／`design/**`／三枚冻结件／golden／`thresholds.go`／`tools/d22scan/allowlist.txt`／D43 表零字节；每把门禁件自落一行 `rc=N`（⛔ 0 字节＝那格没交）；⛔ 零 push、commit 必带显式 pathspec。

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

## Progress log

（追加式；每条 `- [YYYY-MM-DD HH:MM:SS +08] agent=… did=… next=…`，钟点由 `date` 的 stdout 插值，不手打。）
- [2026-10-10 13:22:37 +0800] agent=编排者 did=立票 301 并派 `301-a1`（只读、⛔ go 编译面、写面＝只新建 `probes/301/a1/**`）⇒ 本票 `AC#0` 开工；票面"现量"节全部尺子逐条重跑在 HEAD `772ff880`（⛔ 引工作树——此刻工作树里躺着票 300 落地腿的未提交改动 `wasapi_windows.go` ＋未跟踪新件 `parse_wave_format_300_windows_test.go`，按工作树量会把在飞的活当既有事实）＝`ci.yml:402-403`／`:658`／`:744`／`:780`、`portable-tests.sh:241`／`:248-253`／`:193`／`:364`／`:382-397`／`:413`／`:417-418`／`:593` 逐枚 sed 读到那一行；两枚数＝HEAD 上 `internal/audio` 七枚 `_test.go` 里 2 枚 windows-tagged（`capturelevel_windows_test.go`／`hotplug_test.go`★文件名不带 `_windows` 而内有 tag）、`windows)` 档清单与 `win_pin` 今天各 10 枚且逐枚同形；★`AC#1` 那句"同一笔 commit 改两处"的原文要求跨两行，票面已按两行抄并具名标⛔ 写成一行"逐字"（今天刚为同类事裁过一把）；`AC#3` 只裁形状、⛔ 默认必做，腿若自判"顺手把门改了"＝越界停手；归口关系＝与票 255 那族**同族⛔ 同一枚**，`:593` 里那枚行号锚 `hotplug_test.go:527` 我⛔ 量、写成 `AC#0` 必答 next=`301-a1` 三张表交回后我复跑并对拉尺子口径，再裁 `AC#1` 落不落；⛔ 翻任何框（本票 `AC#0`..`AC#4` 五格全未勾）；⛔ 零 push
