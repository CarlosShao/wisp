
### 编排者复跑＋裁定（2026-10-11 14:2x +08，收 **`308-v1`**＝一笔 `6b4005e5`（100 枚件、越界 0、`.go` 0），件 `.scratch/wisp/probes/308/v1/**`，裁决表 `docs/evidence/s1/308-v1-verdict.md`；台账 `A847`；⛔ 上面各节原句一字不改，本节追加）

★**本轮翻五格＝`AC#1`–`AC#5`**（尺后＝**6 勾／0 未勾**，尺＝`grep -cE '^[[:space:]]*- \[x\]'` ⇄ `- [ ]`）。⚠**本格不随即改名 `-done`**，理由具名＝三处**注释假句**（下面 §4）尚未落地成 `308-r2`；那三句⛔ 是 AC 格、但把它们留着改名＝再造一枚"看到 `-done` 就以为零残余"（`.scratch/wisp/issues/README.md:201` 那句原话）。腿交完、我复跑过、当场改名。

**1. 我自己复跑的（件＝`probes/orch/20261011/logs/orch-A847-308-teeth.txt`，脚本 `teeth-A847-308.sh` 同目录，可逐字重跑）**

`AC#2` 那三件凭据我**整条自己跑了一遍**（此刻两枚腿都交完、`git status -- internal cmd`＝0 行＝无人在飞，才敢在共享工作树里打瞬态突变）：

| 发 | 我这发的数 | 判 |
|---|---|---|
| pristine（原样） | `rc=0`、`--- PASS: TestVolumeInvariantSpellingsVolumeLess308`、`ok internal/risk 0.041s`、panic 命中 **0** | 正控在场 ✅ |
| **守卫拔掉**（按内容锚 `sed '/if len(vol) == 0 {/,+2d'`，⛔ 按行号；落地自证＝残留守卫行 **0**、与备份差 **3 行**） | `rc=1`、逐字 `panic: runtime error: slice bounds out of range [:1] with length 0 [recovered, repanicked]`、顶层 PASS=0／FAIL=1 | 有牙 ✅ |
| 还原 | `cmp` identical、`sha256` 前后同为 **`f2c0130702384bfd…`**＝腿与 `308-r1` 记的那枚**同一枚**；`git status -- internal cmd`＝**0 行** | ✅ |

其余承重读数各自复跑：五枚件（含两枚产码 `pathresolver.go`＝`c347e1cf`、`pathresolver_windows.go`＝`1dcf0640`）在 `1270e36e` blob ⇄ 工作树**逐字节 SAME**＝产码⛔ 动（`AC#1` 那条硬约束由尺坐实）；`t.Skip` 三枚件里**代码形**＝`pathresolver_junction_windows_test.go:128` 那枚**既有**（票 18 允许的 8.3 短名禁用支），新增三行全在注释里＝腿那句"0⇄0／1⇄1"我读原文复认；分母那句我在 linux 件里逐字看到（`7 volume-free spellings still evaluated [...]；2 UNC spellings dropped by name [unc-admin unc-extended]`）；`AC#5`＝五笔 `308-r1` commit 的**票面触及枚数＝0**（`--name-only` 逐笔过尺，所以"⛔ 翻票 72/300/301/302/307"这格是**平凡成立**，而 `AC#0` 那枚勾按 `git log -S` 归我 `0d2260eb`＝腿⛔ 自勾，这条我早先核过、这次由腿的第二把尺再核）；**ubuntu 那档真会跑到新件**＝我自己读了 `scripts/portable-tests.sh` 的 `core` 射程名册，逐字含 `./internal/risk/...`（新件⛔ build tag ⇒ 它从此进 linux 档），腿那句前提成立。

**2. 腿多做的两件我签收（⛔ 我这把尺⛔ 够到）**

- **真在 linux 上跑了一遍**（`docker run --rm golang:1.27`，容器自报 `linux amd64 go1.27.1`＝与本机同版）。我先核这枚⛔ 是编的：`docker` 二进制确在本机（`C:\Program Files\Docker\...`），件 `logs/v26-docker-linux-full.txt`＝**26178 字节**、我数过`^--- PASS`＝**108**、`=== RUN`＝**169**、终态 `ok github.com/CarlosShao/wisp/internal/risk 2.331s`，且里面就带着上面那句 `volume-less input "/tmp/wisp-308-no-volume/…"`。PRE 树同尺＝168／107＝**名册作差恰好＋1 枚**（那枚新用例）⇒ linux 档⛔ 少跑任何东西。**这一发把我派单里那句"拿不到 CI 色就登记欠账"从"推理"顶成了"同版本真跑"，但它仍⛔ 是 `ubuntu-latest` 那台机器的色** ⇒ 见 §3。
- **归因单裁**：我留在 §3 那句"`TMPDIR` 成因"由它单独裁＝**⛔ 立、两把尺**——①读**工具链源码**（`D:/work/base/go/src/os/file_windows.go:284-304`＝本仓 `go1.27.1 windows/amd64` 的 `tempDir()` 全文，只问 `GetTempPath`，整段⛔ 出现 `TMPDIR`；对照尺＝`grep -rn TMPDIR $GOROOT/src/os/` 只命中 plan9／unix 两支）；②**当场跑**（`env -u TMP -u TEMP TMPDIR=/tmp` ⇒ 探针看到 `C:\Users\swq\AppData\Local\Temp`＝系统 temp 兜底）。⇒ **比我派单那句更硬**：我原来只说"Windows 档 Go ⛔ 读 `TMPDIR`"，这把给出机制＋反证。**现象⛔ 否**（那次 panic 的原始件在盘、我逐字读过），我这边按"形状可达已被原始件钉住、本机⛔ 可复现"继续成立。

**3. 五枚欠账逐枚签收**

- **收 `Q-308-v1-A`（CI 跑机色）**＝吸收掉 `Q-308-r1-A/B` 两枚，**归我**、与 `A813` 那批同批取；`AC#3` 本格判据（⛔ 用本机顶替 CI 色／拿不到就具名登记）在**登记**这一支上闭合，所以那格翻勾⛔ 等于这笔账清了——两条分开写，⛔ 混一格（这条＝我自己定式"降级是对的 vs 降级是允许的"的第三次兑现）。
- **收 `Q-308-v1-B`（缺世界那一枚）**：两枚 `//go:build windows` 用例里"去掉 `unc-*`"那一支**在任何可构造台面上到达数＝0**（哨兵句命中：工作树 0／导出树 0／别名 0／linux 0）。⇒ 准确口径＝守卫的牙**只由输入面那枚新用例钉住**，⛔ 是"那两枚用例真在⛔ 卷标台面上跑过一回"。这⛔ 推翻我裁语 ⓑ（我裁语写的就是"只能走输入面"），但它把我裁语里没写透的那半句写透了⇒ **`AC#1` 那句判据原文"这枚用例以可解释的颜色结束"是按 ⓑ 的替换载具闭合的**，口径以本节为准。
- **收 `Q-308-v1-C`（历史成因仍未定）**＝永久欠读数，⛔ 再猜。
- **收 `Q-308-v1-D`（两把⛔ 能证其事的尺）**：`308-a1` 的 `probe-sweep-*` 里有空白／0 字节件、`308-r1` 用 `ls 307*` 当"已 `-done`"的凭据。⇒ 处置＝**⛔ 删**（临时件只建不删），但**下一位⛔ 引它们当凭据**，这条写在台账里；记我＝我盖 A846 的章时**没逐枚看那两把尺**。
- **收 `Q-308-v1-E`（"single place the drive letter is taken"缺射程）**：求值位普查只剩 2 枚（产码 `pathresolver.go:263` 带守卫＋本 helper），所以那句"唯一一处"只在**测试面**为真 ⇒ 见 §4。

**4. ★三处注释假句＝`308-r2`（一枚小落地腿，⛔ 动语义、⛔ 动产码、⛔ 动框）**

1. `pathresolver_anchor_spelling_windows_test.go:100-101` 那句 "every other spelling is **unconditional**, so the test **never quietly stops measuring**"——**比票 §现量 登记的还假**：票当时只说"在 `vol==""` 那形下为假"，而 308 落地后那两枚 UNC 拼法在**所有台面**都成了条件拼法（被同函数下面 8 行正面否定）。⇒ **★这是顶回我 A846 那句的**（我那句把射程写窄了）。替换句我用腿在 `70-overclaim-audit.md` 建议那一版（我逐字读过才采纳：它把分母 9⇄9／7、"具名＋打印才算⛔ 悄悄停"、以及"牙只在输入面"三件都写对了）。
2. 同文件 `:108-110` 的 "the remaining seven volume-free spellings still **run** … the **case** ends with a verdict"＝**helper 推到用例的外推**（那支到达数 0，见 `Q-308-v1-B`）⇒ 收窄成"helper 把七枚留在 map 里；只有选择那一半被 `TestVolumeInvariantSpellingsVolumeLess308` 钉住，windows-tagged 那一支在任何可到达台面上⛔ 夹具"。
3. `pathresolver_volumeless_spelling_308_test.go:19` 的 "**the single place** the drive letter is taken" ⇒ 补射程二字＝"the single place **in this test surface** …"。
4. 门禁照旧（三门禁 rc 自落行＋长度门 GREEN＋格式名册两把并排各写射程＋显式 pathspec＋票尾一条 bullet 1 加／0 删）。**收口＝我复跑"剥注释后 blob `cmp` 同一"那一把尺，然后 `Status:` 更新＋改名 `-done`**（⛔ 这是 AC 判语，⛔ 冒充非实现者验收）。

**5. 顶回我那两枚的落账（★都⛔ 改原句，只追加）**

- 我 §4.0ck 队列句把 `308-r1` 写成三笔（`510c95ca`→`1ad48333`→`ee0a93e4`），**漏计两笔**＝`79242df7`／`1270e36e`；票尾我代提的那枚 bullet 里五笔是齐的 ⇒ 停车点那句以票面为准。
- 我 A846 §1 那句"`TMP=/tmp` 会被补成 **`D:\tmp`**"——腿用同一把在场尺量 `TMP=/` 与 `TMP=\` 拿到的是**当前盘符**（`C:\`）。⇒ **结论⛔ 变**（这台电脑拿不到无卷标的 temp 根），**理由要按"任何可构造值都被补上某个盘符"读**，⛔ 钉死在 `D:` 那一枚上。⚠ 它还补了一句"被补⛔ 等于⛔ 能出现无卷标字符串"：`os.TempDir()` 理论上可以返回 `\\tmp` 那形，只是**这台机器给不出来**——这一条我⛔ 独立复现，标〔腿读数、编排者未复现〕。

**Status（本程更新）**：`AC#0`–`AC#5` **六格全勾**＝`AC#0` 我 13:1x 现跑翻、`AC#1`/`AC#2`/`AC#4`/`AC#5` 本程按腿判语＋我复跑翻、`AC#3` 按"档位那问已具名答完＋CI 跑机色登记为 `Q-308-v1-A`（归我）"翻。**⛔ 随即改名 `-done`**（理由＝§4 那三句注释在 `308-r2` 手里）。撤销口令＝逐枚「**308 我翻的 AC#n 撤**」，与腿自带的「**308-v1 撤 AC#1–AC#5**」同级；`AC#0` 那枚另用「**308 AC#0 撤**」。

- [2026-10-11 14:2x] agent=编排者 did=收 `308-v1`（一笔 `6b4005e5`）⇒ **翻五格**（6 勾／0 未勾）；承重凭据我自己复跑＝拔守卫那发**逐字 panic**＋pristine 绿＋还原 `sha256=f2c01307…` 与腿同枚＋`git status -- internal cmd` 收尾 0 行；产码两枚件 `1270e36e` blob⇄工作树逐字节 SAME；`t.Skip` 代码形＝既有那一枚（票 18 允许）新增三行全在注释；`core` 档射程名册逐字含 `./internal/risk/...`（新件⛔ tag⇒从此进 linux 档）；docker linux 那发我先验真（件 26178 B、108 顶层 PASS、终态 `ok … 2.331s`）才签收；归因单裁升级（GOROOT 源码＋当场跑两把尺）＝**⛔ 立、现象⛔ 否**；五枚欠账逐枚收＝A（CI 跑机色，吸收 r1-A/B，归我同批取）／B（**缺世界**：windows-tagged 那一支到达数 0，牙只在输入面）／C（历史成因＝永久欠）／D（两把⛔ 能证其事的尺，⛔ 删但⛔ 许后人引为凭据，**我盖 A846 章时没逐枚看**）／E（"single place"缺射程）；★两处顶回我 A846＝队列漏计两笔＋`TMP=/` 补的是当前盘符（结论⛔ 变、理由改）；三处注释假句落 `308-r2`，⛔ 动语义⛔ 动产码⛔ 动框，**本格⛔ 改名 `-done` 直到那三句落地并由我复跑"剥注释后 cmp 同一"**；⛔ push
