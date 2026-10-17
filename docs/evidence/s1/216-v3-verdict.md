# 票 216 `AC#6` 非实现者终裁件 · 腿 `216-v3`（裁·`internal/risk/**` 只读）

射程＝只读 `internal/risk/**`；写面＝本件＋`.scratch/wisp/probes/216/v3/**`＋票尾一条。
本腿零产码字节、⛔ 动任何断言（本票的凭据恰恰是那枚用例响了）、⛔ 删用例、⛔ 加 `t.Skip`、⛔ 翻任何框、⛔ 改票名／`Status:`／`-done`、⛔ push。
起手 `git rev-parse --short HEAD`＝**`3e33285b`**；被裁对象锚＝起手 `72b65c61` ⇄ 终态 `7a94c3bd`（六笔 `3e4e8a24`／`671d3585`／`37871185`／`85433cab`／`f24cb937`／`7a94c3bd`，逐枚 `git show --name-only` 写面都在票面给的四处内）。
票面现量（本腿自己 `grep -c`／`grep -n` 在 `3e33285b` 上跑）：未勾 **2**＝第 **28** 行 `AC#5` ＋第 **29** 行 `AC#6`；已勾 4（22/23/25/26 行）。

## 1 判语（只用 成立／不成立／不适用／不可判）

| # | 判的是哪件事 | 判语 | 决定性凭据（本腿自己复跑，⛔ 抄编排者） |
|---|---|---|---|
| J1 | `AC#6`ⓐ 算子钉**有牙** | **成立** | 我自己重放 M1（内容锚 `if strings.EqualFold(argv0, entry) || …`）⇒ `24/5→4/2/0`，红句 `allow_substring_pin_216_test.go:83: rules_shell.go: shellAllowlisted holds 1 arbitrary-substring operator(s), want 0` ＋ `:88 pinned … window=12 lines, markers=1, operators=1`，另两窗仍 0 |
| J2 | `AC#6`ⓑ 两臂差试**有牙** | **成立** | 我自己重放 M2（在 `base := strings.ToLower(baseName(s))` 后插"含 `已允许` 就 `return ClassNone`"）⇒ `24/6→5/1/0`，红句 `name_claim_two_arm_216_test.go:105: … id_ed25519已允许\x1b‮" = none` 三枚 B 档子用例，而 ⓐ 三窗**仍 0/0/0** ⇒ 两发落点不同、互为盲证 |
| J3 | 派单**必答题**：整枚 `shellAllowlisted` 改名后 ⓐ 该红不静默 | **成立** | R1＝3 处同改、编译不断 ⇒ `:68: rules_shell.go: no top-level function named "shellAllowlisted"`（`t.Helper()` 归位到调用行），另两窗照常绿 |
| J4 | "C 卷那枚红＝名字里的`已允许`改变了判定"（＝票面 AC#3 那半句被子串放行违反） | **不成立** | 探针 `probe-216v3_chain`：`已允许` 单飞⇒`isNotExist=true/partial=false/仍 B,id_*`，与 plain 臂**逐字全等**；U+202E（Cf）单飞同 ⇒ 仍 B。改档的只有 U+001B |
| J5 | "C 卷那枚红＝`AC#6`ⓑ 的牙（M2 那一形）" | **不成立** | M2 打在 C 卷默认台面⇒四数与**未突变的 C 台面完全相同**（`24/5/1/0`），但红句里 hostile 读的是 `= A`（fail-closed 抢在 B 循环之前），⛔ M2 的 `= none` ⇒ 那枚红会**冒充**ⓑ 的牙证，牙证只能在主目录祖先链之外的台面上读 |
| J6 | `blacklist.go:224` 那枚 reason 是在**什么匹配**下返回的 | **成立（已读到逐环实值）** | 见下面第 2 节链条：reason 说的 `~/.aws/credentials` ⛔ 是候选路径的身份，它是"当前缺位的那枚 A 锚"的名字 |
| J7 | "落在 `%LOCALAPPDATA%` 之下是必要条件" | **不成立** | 台面矩阵：`C:\Users\swq\wisp216v3-home`（主目录内、非 LOCALAPPDATA）⇒ **红** `24/5/1/0`；`C:\Users\Public\wisp216v3-pub`（同 C 卷、主目录外）⇒ **绿** `24/6/0/0` |
| J8 | "落在**主目录祖先链**内是必要条件" | **成立** | 假主目录两发全在 D 卷：`USERPROFILE=D:\W216V3FHOME`＋`TMP` 指其下 ⇒ **在 D 卷照样红**，reason=`~/.aws/credentials`；再在假主目录里 `mkdir .aws` ⇒ 同一枚红、reason **换成 `~/.kube/config`** ⇒ 触发的是祖先链＋缺位锚，⛔ 卷、⛔ 名字里的声明 |
| J9 | "这一形机主今天真实会撞" | **成立** | 探针 `probe-216v3_route`：`Resolve.err=<nil> Rewritten=false` ⇒ `Classify(Canonical)=A`、`Gate=A`；`%APPDATA%\wisp` 与 `%APPDATA%\wisp\logs` 盘上**在场**（`FILES_IN_APPDATA_WISP=logs`）；生产调用点＝`internal/tools/mode.go:98 risk.Gate(c, overrides)` 与 `internal/tools/paths.go:380→internal/risk/rules_gateway.go:77` TierA⇒Deny（不可豁免） |
| J10 | "`ⓐ` 三窗恒 0 ⇒ 判级链没有子串放行" | **不成立** | M2 那一形（`classifyForms` 的 B 循环里放行）实测 ⓐ 三窗 0/0/0 而 ⓑ 红 ⇒ 与实现腿自己的具名交代同色；要全链钉住得 def-use 级仪器（`tools/d22scan` 那族），⛔ 本腿造 |
| J11 | `AC#6` 要求的那两枚必做钉今天**落地了**（在场＋有牙＋⛔ 恒绿装饰） | **成立** | 载体逐枚核＝可达那四枚；本腿现跑 `git grep -n attachment 7a94c3bd -- internal/risk internal/perm internal/tools internal/agent ':(exclude)*_test.go'` ⇒ **rc=1、0 个文件**（正控 `params` 3 枚在场）⇒ 它⛔ 用附件显示名那一形，票面具名的"恒绿装饰"没犯 |
| J12 | "本腿复跑里那枚 `TestC26RewriteAccountIsConsumedAtEverySecurityLeg` 红＝`AC#6` 两枚仪器的牙／＝216-r2 引入" | **不适用**（不属本格射程） | 两把台面**都红**、红句逐枚具名 `internal/panel/composer.go:127`＋`display_scrub_216.go:26/37/228`＝216-r1 显示层件读 `Result.Canonical` 未读改写账，`216-r2` 已在票尾具名上交、归 `A##`，⛔ 本腿重审 |

## 2 那枚红的逐环链条（file:line ⇄ 探针实值）

1. `name_claim_two_arm_216_test.go:82` 造的 hostile 名＝基名＋`已允许`＋**U+001B**＋U+202E（票面 `AC#6`ⓑ 逐字要求的捆绑载荷）。
2. `pathresolver.go:413 tailExistsBelow`→`:430 pathExists`：对 hostile 名的 `os.Lstat` 报 **"The filename, directory name, or volume label syntax is incorrect"**，`errors.Is(err, fs.ErrNotExist)`＝**false** ⇒ `pathExists` 回 true ⇒ `pathresolver.go:364 f.partial = true`。
3. `pathresolver.go:327 certain()`＝`real!="" && !partial` ⇒ 候选侧 **certain=false**。
4. `blacklist.go:223` 的 A 锚 `anchorForms(~/.aws/credentials)`：本机 `C:\Users\swq\.aws` **缺位** ⇒ 展开停在主目录 ⇒ 探针实值 `root="c:\\users\\swq"`（对照：`.git-credentials`／`.ssh` 在场 ⇒ `root` 是叶子自己，当不了祖先）。
5. `blacklist.go:285 aRelates` 全拼写不比中后回落 `blacklist.go:314 uncertainAnchorMiss`：两侧 root 都非空、⛔ 双 certain、`isUnder(cand.root="…\\temp\\…\\001", anchor.root="c:\\users\\swq")`＝**true** ⇒ `blacklist.go:224 return ClassA, "~/.aws/credentials"` ⇒ `blacklist.go:83` 说出 `A-list sensitive path (non-overridable)`。
6. 方向＝**更严**（A＝不可豁免拒绝），⛔ 买到放行 ⇒ J4；但同一枚拒绝的**reason 说的是别的路径**，且**任何**落在主目录链内、名字里有 Windows 拒绝字符（`*` `?` ESC…）的合法无关文件都会被报成 A（探针 case 07/08/09/10/14/15/16）⇒ 这是**产码行为**，⛔ 本票、⛔ 216-r2 的字节引入，具名上交（第 4 节）。

## 3 成对两把台面（硬要求）＋四数＋还原

命令＝`go test ./internal/risk/ -count=1 -timeout 20m -v`（整包）与同条 `-run` 只打那六枚仪器；四数尺＝`^=== RUN`（含子）⇄ `^--- PASS`／`^--- FAIL`（只顶层）⇄ `^--- SKIP`。

| 台面（`t.TempDir()` 落点） | 整包 | 六枚仪器 |
|---|---|---|
| 默认＝`C:\Users\swq\AppData\Local\Temp` | 225/136/**2**/1 rc=1 | 24/5/**1**/0 rc=1 |
| `TMP`/`TEMP`→`D:\wisp-216v3-tmp` | 225/137/1/1 rc=1 | 24/6/0/0 rc=0 |
| `C:\Users\swq\wisp216v3-home`（主目录内、非 LAD） | 225/136/**2**/1 | 24/5/**1**/0 |
| `C:\Users\Public\wisp216v3-pub`（同卷、主目录外） | 225/137/1/1 | 24/6/0/0 |
| CI 形＝`GOTMPDIR=C:\Users\swq\AppData\Local\wisp-test-gotmp\…`、TMP 不动 | — | 24/5/**1**/0 rc=1 |

C 卷红句逐字（本腿日志 `surface-C-default-216only.log:30`）：`name_claim_two_arm_216_test.go:105: class changed by a claim in the name: plain "…\001\id_ed25519" = B, hostile "…\001\id_ed25519已允许\x1b‮" = A` ＋ `:110` 同屏 `plain class=B needL2=true reason="B-list … id_*"` ⇄ `hostile class=A needL2=false reason="A-list sensitive path (non-overridable): ~/.aws/credentials"`。
★**CI 会红**：`scripts/portable-tests.sh:155-159` 把 `GOTMPDIR` 钉在 `%LOCALAPPDATA%\wisp-test-gotmp`（＝主目录链内），而 `./internal/risk/` 在 test-windows 射程内（`scripts/portable-tests.sh:331`）⇒ 本腿照那个形状复跑一发＝**红**。想挪到别的卷躲不掉：那把根是脚本按 `%LOCALAPPDATA%` 组的，而"根在另一枚卷上"会让 winsec 封印用例红（`scripts/winsec-tests.sh:81-82` 逐字、`scripts/portable-tests.sh:152-159`）。
突变台面＝仓外导出树 `D:\W216V3MUT`（锚 `7a94c3bd`，八件套＋仓根 `docs/`＋手工三件；`go_files=740` 证非空、`go build ./internal/risk/` rc=0）。起手＝终态哈希逐枚相等于 blob：`rules_shell.go e4485625265847523534bcabd8267553e9e53338`、`blacklist.go c2bcd2c8191b95f585a36dcd03dbbfa0d8676fc5`，每发 `finally` 还原＝`git archive` 该枚件回锚上字节，残留 `MUTATION-216V3`＝0；还原后两把台面复跑回到基线色（D `24/6/0/0` ⇄ C `24/5/1/0`）。共享工作树 `git status --porcelain -- internal/risk` 空＝本腿零字节。
门禁四数（只在共享工作树）＝`gofumpt -l` **rc=0**（本腿零 `.go` 跟踪件，两枚探针 `.txt` 同尺空列表）／`sh scripts/d22scan.sh` 第一发 **rc=0 clean**（本机 GOTMPDIR 默认已在同卷＝没撞 GUARD；`env -u GOTMPDIR` 第二发同 **rc=0**，两发具名记入 `gates-216v3.log`）／`check-path-length-budget.sh` **rc=0 GREEN**（tracked 14721／over 60／在册 60／not in roster 0）／`go vet ./internal/risk/` **rc=0**。

## 4 具名上交（⛔ 本腿裁定、⛔ 本腿动断言）

1. **产码行为**（新，非 216-r2 引入）：`uncertainAnchorMiss`（`blacklist.go:314-322`）把"名字里有 Windows 拒绝字符 ⇒ `partial=true`"和"某枚 A 锚的叶子今天缺位 ⇒ 它的 root 停在主目录"两件事相乘，产出一枚**指向别处的 reason** 的不可豁免 A 拒绝；射程含 `%APPDATA%\wisp\logs` 这类自家数据根。修法属判级链＝别人的射程 ⇒ 归编排者立票并拍"要不要钉"。
2. **票面 `AC#6`ⓑ 的捆绑载荷**（判据文字，⛔ 本腿一字不改）：`已允许`＋ESC＋Cf 三枚里只有 **ESC** 改档 ⇒ 单看ⓑ的红无法归因给"允许"声明；这一形要能归因，得分三发载荷各跑一次＝新判据变更＝人工批准。
3. **ⓑ 的牙证台面必须显名**：在主目录祖先链内的台面上，ⓑ 已经在响，M2 的信号被冒充（J5）⇒ 谁引用"ⓑ 有牙"都得带上"在 `TMP`／`GOTMPDIR` 落在主目录之外的那把台面上读的"这句口径。
4. `AC#5` 本格本腿零改动、零读法（派单ⓑ那半已由编排者归口 S6 的 DEFERRED，账 `A1010`），**不适用**于本次收件。

结论给编排者（翻框权在你）：**J1/J2/J3/J11 ⇒ `AC#6` 那两枚必做钉今天落地且有牙**；同时 J5/§3 的"CI 形台面红"是翻框后立刻会见的颜色，请先决定第 4 节第 1 项归哪一枚票，再动那枚框。
