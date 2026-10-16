# 票 317 验收腿 317-v1 判语（非实现者）——测试临时根的分隔符拼法与 junction 假红

- 腿身份＝**非实现者**（实现者＝`317-r1`／笔 `78805129`；本腿零产码，判断只落本件＋`.scratch/wisp/probes/317/v1/**`＋票尾一条 bullet）。
  ⛔ 未翻框、⛔ 未动 `Status:`／票名／`-done`、⛔ 未 push、⛔ 未动 `scripts/**`／`tools/**`／任何测试件／`.github/**`。
- 起手锚（现读）＝`git rev-parse --short HEAD` = **55e11b2f**；被检笔 `78805129` 真在（`git log --oneline --all --since='2026-10-14 00:00' -- scripts` 命中）；
  读数台＝**仓外导出树** `git archive 8235e683 | tar -x -C D:\tmp\wisp317v1\tree`（六件套 `go.mod go.sum internal cmd tools scripts` 逐枚 `[ -e ]` 验在；
  `find … -name '*.go' | wc -l` ＝ **1087 枚 .go 文件**＝非空尺后再信其 0）。
- 台面（三枚别腿在飞）：`git status --porcelain -- internal` 现读＝`internal/session/manage_test.go` ＋ `internal/winsec/{winsec.go,winsec_other.go,winsec_windows.go}` ＋ 3 枚未跟踪 winsec 测试件。
  本腿要读的 `./internal/risk/` **未受影响**（导出树与共享树都编得过：`go build ./...` rc=0、`go vet ./internal/risk/` rc=0），⇒ 无需另归别腿红。
- 本腿⛔ 跑：整包 `go test ./cmd/wisp`（票 314 死亡区）／`go mod tidy`／`go get`／完整 `--scope=windows` 八包与 `--scope=cli`（判据不需要，⛔ 扩射程）。

## 0. 结论先行（判语只三种）

| 格 | 判语 | 一句话凭据（尺＋射程在本件 §1–§4） |
|---|---|---|
| **AC#1 复现** | **成立** | 红落在那四行（两行字面 `want`／两行真 `mklink /J` 目标），⛔ 不是建目录失败；实现腿那对读数同一台面但**叶名也不同**（欠洁一处），我用**同目录两拼法**把那发补成严格对照：`C:\tmp\wisp317v1pair` ⇒ 顶层用例 `PASS=4 FAIL=0` rc=0 ⇄ `C:/tmp/wisp317v1pair` ⇒ `PASS=0 FAIL=4` rc=1；四名今天逐枚 `-list` 自名命中（5/5）；"原意"三条我独立复核到凭据 |
| **AC#2 修** | **成立** | 三处组成式现读两处位置都原生且运行时值逐字单反斜杠；override 仍原样采用**今天仍成立**（posix override ⇒ 拼法未被改写、仍 4 红＝复现通道未拆）；4 处 `wisp_311_ref_path` "不修"这个**处置**对（折形永不流向 GOTMPDIR＋读取器自己再折一次）；★**恒真检查两向都红**：退回原状 4 红、**反形（只把尾接写成 `/`）也 4 红**、出厂形全绿 ⇒ 尺对分隔符在两个位置都敏感，这修是承重的 |
| 第四处同形（`tools/d22scan/runtests.sh:92`） | 判语见 §4 | 今天＝**混形**（基座反斜杠＋尾接正斜杠）；**形状会爆**（同 runner 同根指四枚名 ⇒ 4 红 rc=1），今天不爆只因为**射程**（CI 三格都不覆盖那四枚）⇒ 建议**新立一票**、一行修 |

一句话处置建议（编排者拍）：**立一枚新票**把 `tools/d22scan/runtests.sh:92` 的尾接 `/d22scan-` 改成 `\\d22scan-`，并把票 317 面上"同形、今天未爆"那句重述为"混形，射程遮蔽中的活坑"。

## 1. AC#1 判据复核——三问逐问

尺：`go test ./internal/risk/ -count=1 -v -run <四名>`（导出树；四数＝`grep -c '^--- PASS' / '^--- FAIL' / '^=== RUN'`＝**顶层用例**数）。

| 发 | GOTMPDIR 逐字形状（`printf '%s' \| od -c` 已进日志） | 四数 | rc | 读数件 |
|---|---|---|---|---|
| R2a | `C:\tmp\wisp317v1pair`（原生） | `RUN=4 PASS=4 FAIL=0 SKIP=0` | 0 | `probes/317/v1/logs/r2a-native.txt` |
| R2b | `C:/tmp/wisp317v1pair`（**同一目录**，只差分隔符） | `RUN=4 PASS=0 FAIL=4 SKIP=0` | 1 | `logs/r2b-posix.txt` |

① **红在哪一行**（判"是不是别的建目录失败"）＝四行逐字，与 CI 同形：
- `syncdirs_ancestor_actable_leg_116_test.go:182: precondition: deepestExistingAncestor("C:\\tmp\\wisp317v1pair\\…") = "C:\\tmp\\wisp317v1pair\\…", want "C:/tmp/wisp317v1pair\\…"`（**want 混合分隔符**＝票面那条症状）
- `syncdirs_ancestor_actable_leg_116_test.go:286: precondition: the ~ spelling must resolve under the relocated home "C:/tmp/wisp317v1pair\\…"`
- `syncdirs_ancestor_reparse_ticket105_windows_test.go:59` 与 `:116`：`mklink /J C:\tmp\wisp317v1pair\…\002\proj C:/tmp/wisp317v1pair\…\001 failed: exit status 1: Invalid switch - "tmp".`（**真 mklink 目标被 cmd 当开关**＝CI 的 `Invalid switch - "Users".` 同句尾）
⇒ ⛔ 任何一发红来自 `mkdir`／建根失败；harness 的 `DEGRADED` 行未出现。实现腿 `logs/10` 那四行与我的四行**同行同号**＝可当凭据。

② **成对是否只差分隔符**＝实现腿那对是同一台面（同树同命令、都直测 `./internal/risk/`）但叶名不同（`…/posix` ⇄ `…\nat`）：不是"改前跑共享树、改后跑导出树"那种混台面（⛔ 此病），但留下一个混淆项。R2a/R2b 把叶名钉成同一目录后同四枚名同四行复现 ⇒ 混淆项排除，原读数加固。**另：⛔ 我用 `-overlay` 抵账**——本腿全部台件都是**仓外物理树**（尺运行时读盘）。

③ **四名今天存在**＝逐枚 `go test ./internal/risk/ -run <名> -list <名>`，命中的判据＝该名字自己被打出来：四枚 116/105 ＋ `TestPathResolverJunctionWindows` ＝ **5/5 命中**（⛔ 用不存在的名数 0 这条先例我照着验了名）。

**"先读 310／311 的原意"这半**：三条我都能独立指到——① `.scratch/wisp/probes/310/r1/10-impl.md:23-24` 把 `${LOCALAPPDATA//\\//}` 组成式与"本机现量形 `C:/Users/swq/…`"并列登记（＝转录风格）；② `scripts/portable-tests.sh:113-151` 那四条长注释里 `fold|slash|separator|posix` 只命中 ticket 317 自己加的两行 ⇒ 折法从未被写成需求；③ `grep -rn GOTMPDIR --include=*.go internal cmd tools` 只命中**别腿在飞件的一行注释**（`internal/winsec/openleaf_guard_322_windows_test.go:95`）、0 处真读取 ⇒ "无测试依赖其拼法"成立；反向依赖（116/105 拿 `t.TempDir()` 字面拼法做串比与不经 `filepath.Join` 的 mklink 目标）由 R2b 直接量到。

## 2. AC#2 四支逐枚独立复核

| 发 | 台件（⛔ overlay；全部物理树） | 运行时 GOTMPDIR 逐字 | 四数（顶层用例） | rc |
|---|---|---|---|---|
| R3 | 出厂形 `tree`（＝HEAD 8235e683） | `C:\Users\swq\AppData\Local\wisp-test-gotmp\portable-20261016-090055-1512` | `RUN=200 PASS=132 FAIL=0 SKIP=0` | 0 |
| R4 | **M1**＝`:155-156` 逐字退回 `78805129^`（`diff` 只这两行） | `C:/Users/swq/AppData/Local/wisp-test-gotmp/portable-20261016-090221-552` | `RUN=200 PASS=128 FAIL=4` | 1 |
| R5 | **M2 反形**＝只把 `:156` 尾接 `\\` 写成 `/`（`diff` 只一行） | `C:\Users\swq\AppData\Local\wisp-test-gotmp/portable-20261016-090256-813` | `RUN=200 PASS=128 FAIL=4` | 1 |
| R6 | 出厂形＋`WISP_TEST_GOTMP_ROOT='C:/tmp/wisp317v1/ovr'`（`od -c` 进日志） | `C:/tmp/wisp317v1/ovr\portable-20261016-090416-462` | `RUN=200 PASS=128 FAIL=4` | 1 |

读数件＝`logs/r3-harness-asshipped.txt`／`r4-harness-m1-prefix-fold.txt`／`r5-harness-m2-mixed-tail.txt`／`r6-posix-override.txt`。三发红的都是**同四枚名、同四行号**（R4/R5/R6 的 `--- FAIL:` 逐枚名列在日志里）。

- **① 三处 GOTMPDIR 赋值（现读，⛔ 抄它的行号）**：`grep -n GOTMPDIR scripts/portable-tests.sh` ⇒ 组成式在 `:155`（基座 `${WISP_TEST_GOTMP_ROOT:-"$LOCALAPPDATA\\wisp-test-gotmp"}`）＋`:156`（尾接 `"$gotmp_base\\portable-…"`）；`scripts/winsec-tests.sh:69`／`scripts/wisp-cli-tests.sh:82` 各一行同形。票面那句"元凶行 `portable-tests.sh:144`"今天是一行注释（⇒ 派单更正 (a) 复核＝成立）。R3 的运行时值是**单反斜杠** ⇒ `${X:-"…\\…"}` 吃一层这条今天仍成立、写法正确。
- **② "override 仍原样采用"防的是什么**＝防这把尺被拆：R6 实测操作员拼法**未被改写**（打印行逐字 `C:/tmp/wisp317v1/ovr\portable-…`）且四枚仍红 ⇒ 复现通道在用。今天这一支**仍成立**：四枚入口（`portable:155`／`winsec:69`／`cli:82`／`runtests:92`）都写成 `${WISP_TEST_GOTMP_ROOT:-…}` 原样采用，⛔ 发现别腿改写过这一支。
- **③ "不修的那几处"判对没有**＝**处置对、一条理由无仪器凭据**。现量是 **4 处**赋值（派单更正 (b) 复核＝成立）：`scripts/portable-tests.sh:818`／`scripts/winsec-tests.sh:174`／`scripts/wisp-cli-tests.sh:186`／`tools/d22scan/runtests.sh:155`。四处只喂两样：卷读取器（其第一行自己再折一次：`wisp_311_arg=${1//\\//}`，四枚入口各有＝`:811/:167/:179/:148`）与 INCONCLUSIVE／GUARD 显示文本；`wisp_311_ref_path` 从不流向 GOTMPDIR。⇒ **不修**对。但实现腿那句"改它会动 311 已钉的显示文案"按现量**没有仪器**钉它：`grep -rn wisp_311_ref_path --include=*.go --include=*.yml --include=*.ps1 .`（排除 `.scratch`）＝**0 命中**（正控＝同串在 `*.sh` 命中 **16 行**，先跑正控再报 0）。
  **"差的那一处"归谁**＝`runtests.sh:155` 在它写面外（`tools/**`），它的尺名册也**确实只数了 `scripts/`**（`probes/317/r1/logs/24-no-fold-left.txt` 第 1 行标题逐字"any LOCALAPPDATA fold left in scripts/"）⇒ "三处"是**射程限定的计数**、不是"仓里只有三处"的断言；但它的交付里第 4 处**从未被具名复核过**（§3ⓒ 只点了同文件 `:92` 那行），这一处判定由本腿现在补上（结论＝同样不修）。
- **④ ★恒真检查（首要攻击点＝判据仪器自己）**＝**两向都红 ⇒ 不是空尺**。"改前必红"有 M1（4 红 rc=1，GOTMPDIR 逐字 `C:/Users/…`）；派单追问的"**换成反形会不会也红**"我也跑了：**M2 也 4 红**，红句逐字
  `syncdirs_ancestor_reparse_ticket105_windows_test.go:59: mklink … C:\Users\swq\AppData\Local\wisp-test-gotmp/portable-20261016-090256-813\… failed: exit status 1: Invalid switch - "portable-20261016-090256-813\TestSy…"`
  ＋ `syncdirs_ancestor_actable_leg_116_test.go:286: precondition: the ~ spelling must resolve under the relocated home "C:\\Users\\swq\\AppData\\Local\\wisp-test-gotmp/portable-…"`。
  ⇒ 这把尺对**基座与尾接两个位置**的分隔符都敏感，317-r1 改的两处都是承重的；"造不出坏世界"这一支（缺世界）不适用，⛔ 需要写"证据不足"。

## 3. 门禁与台面（读数件 `logs/r7-gates-sharedtree.txt`，台面＝共享工作树）

- `bash -n`：`scripts/portable-tests.sh`／`scripts/winsec-tests.sh`／`scripts/wisp-cli-tests.sh`／`tools/d22scan/runtests.sh` ＝各 **rc=0**（脚本那把尺）。
- `go build ./...` **rc=0**；`go vet ./internal/risk/` **rc=0**（⇒ `322-r1`／`166-r2` 此刻的在飞态没有把本腿的射程打红；若红我会另节具名、⛔ 计入本票）。
- `sh scripts/d22scan.sh` **rc=0** clean：`runtests.sh: OK - packages=[./...] top-level: PASS=35 FAIL=0 SKIP=0, === RUN=77, '[no tests to run]'=0`；扫描分母 `ban #8 internal/ = examined 571 Go files`（＝含别腿新件，⛔ 归本票、⛔ 当红）。
- `sh scripts/check-path-length-budget.sh --with-self-test` **rc=0** `VERDICT GREEN`（`tracked paths=13497 over-budget=60 covered=60 not-in-roster=0`）；仓根跑＝本件与读数件的路径名也在这把尺里过。
- `gofmt -l internal/risk scripts tools`（**比门禁那把窄**，本机 `gofumpt` ⛔ 装＝`command not found`／rc=127 ⇒ ⛔ 冒充）：退码恒 0，**列出 1 枚文件**＝`internal\risk\provenance.go`。
  ★归类＝**台面条件、⛔ 计入本票判据**：`git status --porcelain -- internal/risk` 空、`git ls-files --eol` 读 `i/lf w/crlf attr/text eol=lf`、字节级 `od -c` 显示工作树副本是 `\r\n`；**同一枚文件在 8235e683 导出树里 `gofmt -l` 列出 0 枚** ⇒ 共享树那枚是 CRLF 工作树副本（内容＝blob），与票 317 无关。
- 归类面故障扫描：本腿所有红句里 ⛔ `walk`／`build failed`／`no tests to run`／`0xc0000135`／`RUN=0`（`RUN=200`／`RUN=4` 都是非零读数）。
- ⚠ 我自己踩到的一把假尺：第一次用 `grep -rn 'LOCALAPPDATA//\\\\//}'` 数折法命中＝**0（假 0**，正则把反斜杠吃掉一层），改 `-F` 固定串＋正控（`*.sh` 命中 16 行）后才用。

## 4. 第四处同形——`tools/d22scan/runtests.sh:92`（本腿⛔ 动，只报判语＋建议）

现量形状（`sed -n '92p' | od -c`）＝基座 `"$LOCALAPPDATA\\wisp-test-gotmp"`（原生）＋尾接 `/d22scan-…`（正斜杠）＝**混形**（派单更正 (c) 复核＝成立）。
出身＝编排者两笔只镜像了基座：`581dd252`（折形基座 → `"$LOCALAPPDATA\wisp-test-gotmp"`，尾接 `/` 原样留着）＋`76b9e214`（单反斜杠 → 双，"字节级验证"）。
⇒ 实现腿报的那句"仍折形"按今天现量必须重述为"基座已原生、**尾接仍正斜杠**"。

会不会爆（实测，导出树、`env -u GOTMPDIR -u WISP_TEST_GOTMP_ROOT`；读数件 `logs/r8-runtests-mixed.txt`）：

| 发 | 射程 | 运行时 GOTMPDIR 逐字 | 结果 |
|---|---|---|---|
| S1 | `runtests.sh ./internal/risk/ -run TestPathResolverJunctionWindows`＝`ci.yml:831` 真用那一格 | `C:\Users\swq\AppData\Local\wisp-test-gotmp/d22scan-20261016-090656-522`（混形） | `top-level: PASS=1 FAIL=0 SKIP=0, === RUN=1` rc=0（该件的 link/target 走 `filepath.Join`＝Clean 归一化，吃得下混形） |
| S2 | 同一枚 runner、同一枚混形根，`-run` 换成四枚 116/105 名 | `…wisp-test-gotmp/d22scan-20261016-090659-565` | **四枚全红** rc=1；`mklink … failed: exit status 1: Parameter format not correct - "d22scan-…\TestSy…"`（cmd 措辞与 M2 的 `Invalid switch` 不同，**同一因**） |

⇒ **判语：形状会爆，今天不爆只因为射程。** CI 直呼 runtests.sh 的三格是 `ci.yml:82`（`-C tools/d22scan ./...`，本腿 `sh scripts/d22scan.sh` rc=0 已覆盖）、`:462`（`internal/proc TestLayoutForTestEnv`，该件把 `t.TempDir()` 原样传回比较、不比拼法）、`:831`（S1）；
那四枚 116/105 一名都不在射程里。**建议＝新立一枚票**（一行：`/d22scan-` → `\\d22scan-`；判据＝改后 S2 那发应转绿、S1 仍绿），
并把票 317 上"同形、今天未爆"这句改成"混形，射程遮蔽中的活坑"——登记措辞若留着，下一枚把 `runtests.sh` 指向 `./internal/risk/` 的腿会把它当新 bug 重查一遍。

## 5. 交件自证三把尺

1. **提交**＝见本件所在笔（`git show --stat HEAD` 的路径清单只含 `docs/evidence/s1/317-gotmp-separator.md` ＋ `.scratch/wisp/probes/317/v1/**` ＋票 317 一行）；前批＝`252f6767`／`4e43d934`。
2. **证据件与读数件的绝对路径＋字节数**＝本件路径＝`D:\work\workspace\projects plans\Wisp\docs\evidence\s1\317-gotmp-separator.md`；
   读数件全部在 `D:\work\workspace\projects plans\Wisp\.scratch\wisp\probes\317\v1\`（`00-readings.md` ＋ `logs\r2a-native.txt`／`r2b-posix.txt`／`r3-harness-asshipped.txt`／`r4-harness-m1-prefix-fold.txt`／`r5-harness-m2-mixed-tail.txt`／`r6-posix-override.txt`／`r7-gates-sharedtree.txt`／`r8-runtests-mixed.txt`＋`.s1` ＋四枚仪器 `r2-pair.sh`／`r3-harness.sh`／`r6-override.sh`／`r7-gates.sh`／`r8-runtests.sh`）。
   ⛔ 在本件里自指一个还没写出的号：逐枚**字节数**由票尾那条 bullet 给出（本件提交后用 `git cat-file -s` 在盘上现量，一次成文）。
3. **两格判语**＝**AC#1 成立**；**AC#2 成立**；第四处同形＝形状会爆／射程遮蔽，**建议立新票**（`tools/**` 不在本腿写面，⛔ 未动）。

## 6. 没裁完的格（具名＋缺哪一发读数）

- **完整 `--scope=windows` 八包**与 **`--scope=cli` 全量**＝⛔ 裁。派单本票只把判据钉在 `./internal/risk/` 与第四处同形，这两发的读数⛔ 由本腿取（八包里 `cmd/wisp` 那档＝票 314 死亡区，跑它⛔ 被允许）。
- 缺哪一发：如果要把 AC#2 判据②"本地能以 harness 的 windows 腿跑通"读成**整条 windows 腿**（八包），需要一发 `bash scripts/portable-tests.sh --scope=windows` 的八包四数——本腿交的是**该 scope 内本票射程那一包**（`./internal/risk/`＝`RUN=200 PASS=132 FAIL=0 SKIP=0` rc=0，R3）＋改前后两发的对照（R4/R5）。
