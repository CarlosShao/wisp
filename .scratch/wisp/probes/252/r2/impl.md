# 票 252 · R2 腿（`252-r2`）—— 给"两段包含不许并成一段"这枚硬禁造仪器

## §0 起手锚（本腿自己现取，⛔ 不引用别人给的号）

- `date -Iseconds` = `2026-10-03T10:33:43+08:00`
- `git log -1 --format=%H` = `e16e303773f59be81d2f02eee3cc638870aa1630`（分支 `dev`）
- `git status --porcelain -- internal/tools internal/risk` = **0 行**（起手即为空 ⇒ 允许开工）
  - ⚠ 一次工具侧误读要记下：第一次跑这枚尺时 shell 的工作目录是 `.scratch/wisp/issues`，
    于是相对 pathspec 打空并吐出一行
    `warning: could not open directory '.scratch/wisp/issues/internal/': No such file or directory`（rc 仍为 0）。
    在仓库根 `D:/work/workspace/projects plans/Wisp` 重跑才是那枚 0 行的读数。**结论：pathspec 必须在仓库根求值。**
- 起手名册（`go test -count=1 -v ./internal/tools/`，落
  `.scratch/wisp/probes/252/r2/baseline-v-test.log`）：
  - 全仓规矩的那枚正则 `^(--- )?(PASS|FAIL|SKIP)` 命中 **187 行**＝186 枚顶层 + 1 行裸 `PASS`；
  - 含子测试的 `^[[:space:]]*--- (PASS|FAIL|SKIP)` 命中 **257 行＝257 PASS／0 FAIL／0 SKIP**；
  - 包体行 `ok github.com/CarlosShao/wisp/internal/tools 14.166s`；`go test -count=1 ./internal/risk/` = `ok 5.221s`。
  - ⇒ 编排者派单里"257"这一枚**复量为真**，但**它不是他那条正则能数出来的数**（那条数出 187）。见 §7 第 R2-1 条。

## §1 问题：那枚硬禁今天到底有没有尺

票面 AC#2 的硬禁原文（`252-...-r2-l2.md:30`）：⛔ **绝不许**把两次包含改成一次（那是把票 107 的洞重新打开）。
非实现者验收表 `docs/evidence/s1/252-allowlist-sameform-v1.md` 量出这条零仪器（M-merge 257 全绿），
本腿的第一件事是**自己复认**，第二件事是**造出能让它当场红的尺**。

现跑的码体（`internal/tools/paths.go`，锚点 `e16e3037`，本腿一个字都不改）：

| 段 | 行 | 逐字形状 | 谁喂它键 |
|---|---|---|---|
| 词法段（LEG 1） | `:196` | `if !rootsContain(p.roots, f) { return false }`，`f := foldPath(canonical)`（`:190`） | **只有** `canonical` 这一枚入参 |
| `resolvedForm` 段（LEG 2） | `:208-209` | `rf, ok := resolvedForm(canonical)` + `if !ok \|\| !rootsContain(p.roots, foldPath(rf)) { return false }` | `resolvedForm` 的返回值 |
| 工作区段（LEG 3，票 92） | `:215` | `if p.workspace != "" && !rootsContain([]string{p.workspace}, f)` | 也是 `f`，⛔ 但它的书不是 `p.roots` |

为什么 M-merge 量不出来（本腿的机制解释，不是抄来的）：修法 `2b1a3071` 把同形那一步挪进了
`Canonicalize`（`:145-148` → `sameFormOfUnresolved`），于是**生产能递给 `InAllowlist` 的每一枚串都已经两侧同形**
⇒ 对任何这类输入 `LEG1 == LEG2` ⇒ 包内那把 `measure252`（`paths_shortname_252_probe_test.go:114-126`）
算的恒等式 `final == leg1 && leg2` 里的 `leg1` 是**测试自己重算的**，不是从生产函数里读出来的，
生产函数少一条腿它照样绿。⇒ **凡是对"同形输入"取值的尺都天生看不见 M-merge**，这不是偶然失手，是射程不覆盖。

## §2 修法选择：三把尺＋为什么不是别的形

本腿只加仪器，写面 `internal/tools/**_test.go`（＋必要时 testdata）。生产码一字未动（终态 porcelain 见 §5）。

1. **W1 词法段的分歧见证（ behavioural，`//go:build windows`）**：
   取真 8.3 短名 + 一枚还不存在的叶子，**绕过 `Canonicalize` 直接**递给 `InAllowlist`。
   这一枚输入的形状是 `LEG1=false / LEG2=true` —— 全仓**只有**这种"两段不同答"的输入能让删掉词法段变红。
   断言 `InAllowlist=false`，并附一枚"同一条 raw ask 走完整生产回合必为 true"的对照，
   说明这一枚红不是"多拦了一次真授权"，而是"少了第二段之外还少的第一段"。
2. **W2 `resolvedForm` 段的分歧见证（portable，复用票 107b 的 `makeDirLink107b`／`dirLinkEvidence107b`）**：
   根真树、根内一枚 junction（POSIX 上是 symlink）通到根外。`LEG1=true / LEG2=false`，
   断言 `InAllowlist=false` —— 拆掉第二段当场红。
3. **S1 调用图尺（portable，`go/parser`＋`go/ast` 读 `paths.go`）**：
   只对 `InAllowlist` 的**函数体**取断言，两格：
   - 计数：体内对**授权本（`p.roots`）取用的包含调用**必须 **恰为 2 枚**，各带一条独立的拒绝支；
   - 来源：**两枚的键不许同源自同一枚表达式**——一枚必须只由入参折叠得来、另一枚必须经过一次再解析调用。
     这一格是行为尺天生量不到的那一种突变：把词法段的键悄悄换成 `foldPath(rf)`，
     对**所有**输入（含 W1、W2）取值都不变，两段却已经在实现上并成一段。
   - ⛔ 不给生产码加任何导出钩子：读文件用包内相对名，`go test` 的工作目录就是包目录。

选边理由（⛔ 不是"结构尺更好"的一句话）：**W1/W2 有牙但不完全，S1 完全但没有牙**。
删除任一段 → W1 或 W2 必红（这是票面要的"当场变红"）；把两段并成同源的"伪装合并" → 只有 S1 红。
两把互为补集，缺一把就留一个洞。编排者建议的"按调用图而不是按返回值"本腿**接受其一半**：
调用图尺保留（S1），但**不许只交它**，因为纯结构尺对"改成两段同源的两次包含"以外的情形
（例如把 `rootsContain` 内联成循环）会误响，而它自己宣称的"两次独立判定都在"最好有一枚取值证据兜着。

脆性自陈（命中面／下次合法重构会不会误响）写在 §3 末与 §6，读数一律现跑。

## §3 正控读数（五发突变名册＋三枚尺各自的正控，全部现跑）

### 3.1 突变名册（每一发都：cp 副本 → 改 → 整包 `-count=1 -v` → cp 还原）

突变件都留在 `.scratch/wisp/probes/252/r2/mut/`（`m-*.log` 是整包读数，`m-*.red.txt` 是抽出的红句）。
⛔ 终态 `internal/tools/paths.go` 的 md5 = `7a86da7aeba8420639491cd252a47a0c` ＝起手副本逐字节相等（§5 有 porcelain 证据）。

| 代号 | 改了什么（都在 `InAllowlist` 体内） | 未加本腿尺 | 加本腿尺 | 谁响 |
|---|---|---|---|---|
| **M-merge** | 删掉词法那一段（`paths.go:196-198` 三行） | **257 PASS／0 FAIL** ⇒ 复认 `252-v1` 那句"零仪器"为真 | **258 PASS／4 FAIL** | S1 断言 a ＋ W1（顶层＋B1＋B2 三行）；W2 正确地不响 |
| **M-unresolve** | 删掉 `resolvedForm` 那一段（`paths.go:208-211`，保留 `rf, ok :=` 以隔离"包含"这一半） | 254 PASS／3 FAIL（既有两枚：r1 的 N7、107b probe C） | 257 PASS／**5 FAIL** | W2 ＋ S1 断言 a ＋ 既有那两枚 |
| **M-samesource** | 两段都在、各自仍能拒，但**都问同一枚已折好的键**（`rf` 上提，两段的 arg1 同为 `foldPath(rf)`） | **257 PASS／0 FAIL** ⇒ "枚数够"不等于"没并成一段"，既有名册全瞎 | 258 PASS／**4 FAIL** | S1 断言 c ＋ 断言 d ＋ W1 两格 |
| **M-inroots** | 最写实的那一发："同形已经保证了，问一次就够" ⇒ 只留词法段，`resolvedForm` 调用留着但不再参与判定（等价于把判定改走 `paths_workspace.go:106` 的 `inRoots`） | 未单跑；**行为与 M-unresolve 同一枚函数**（都只剩"词法段＋工作区段"），故取 M-unresolve 那一发的既有名册读数 254／3——这一格是**读数转移**，不是本腿另跑的一发，具名标出 | 257 PASS／**5 FAIL** | W2 ＋ S1 断言 a ＋ 既有 N7 ＋ 既有 probe C |
| （对照）**未突变** | 一字不改 | 257 PASS／0 FAIL（起手名册） | **262 PASS／0 FAIL／0 SKIP**（新增 5 行名册） | — |

⇒ 票面要的判据达成：**拆掉任一段都当场变红**，且两形（词法段／`resolvedForm` 段）各有本腿的一枚取值尺；
再加一发"两段同源"的伪装形，只有调用图尺看得见枚数之外的部分（W1 也响，见 §7 第 R2-6 条对这句的收窄）。

### 3.2 红句逐字（⛔ 不改写、不缩写；换行是终端造成的，字面照抄）

**形一 · 词法段被拆（M-merge）** —— W1 的 B1 格：

```
paths_twocontainments_252_r2_windows_test.go:89: TICKET 252 BAN RED (lexical containment missing / two containments merged into one): InAllowlist("D:\\work\\WORKSP~1\\PROJEC~1\\Wisp\\q252r2-lexical-leg-only.md") = true (the 8.3 spelling of the allowed root, new file under it) although the lexical leg refuses it: foldPath(canonical)="d:\\work\\worksp~1\\projec~1\\wisp\\q252r2-lexical-leg-only.md" sits under no root in ["d:\\work\\workspace\\projects plans\\wisp"], and only resolvedForm="D:\\work\\workspace\\projects plans\\Wisp\\q252r2-lexical-leg-only.md" (ok=true) would. AC#2's "绝不许把两次包含改成一次" is what keeps that leg in the code; ticket 107 round 2's own note (paths.go:199-207, "Requiring both is strictly narrower") is the reason it may not be called redundant and removed.
```

同发的 S1（结构侧，枚数）：

```
paths_twocontainments_252_r2_test.go:301: TICKET 252 BAN RED (assertion a: two containments, not one): InAllowlist's body asks the authorization book 1 times, want 2 (found: [paths.go:208 key=foldPath(rf)]). AC#2's "绝不许把两次包含改成一次" is a hard ban, and the shape this fires on is exactly the one no value test can see: after 2b1a3071 every string production hands this function is already one form, so the two legs agree on all of them and a missing leg is invisible from outside. If this deletion was NOT intended, restore both legs at paths.go:196 and paths.go:209. If the code is stale and the ruler is right, the ruler is the thing that needs a human approval, not a green test run.
```

**形二 · `resolvedForm` 段被拆（M-unresolve）** —— W2：

```
paths_twocontainments_252_r2_test.go:158: TICKET 252 BAN RED (resolvedForm leg missing / two containments merged into one): InAllowlist("C:\\Users\\swq\\AppData\\Local\\Temp\\TestTicket252R2ResolvedLegRefusesWhatOnlyItRefuses1234709264\\001\\proj\\esc\\marker.txt") = true while the lexical leg is the only one that vouched for it (leg1=true leg2=false, resolvedForm="" ok=false, roots=["c:\\users\\swq\\appdata\\local\\temp\\testticket252r2resolvedlegrefuseswhatonlyitrefuses1234709264\\001\\proj"]): that path leaves the allowed root "...\001\proj" through the link "...\001\proj\esc", which the OS reports as "...\001\outside". Either half of the second leg - "it will not answer" or "its answer is elsewhere" - is the refusal, and both are the same containment in paths.go. Ticket 107's hole re-opened: the pair looks redundant on aligned input by design, and AC#2's "绝不许把两次包含改成一次" is what keeps it that way.
```

（同一条路径的 `...` 是本腿为省版面折的**日志内容**，不是尺的字面；尺的字面全文在 `mut/m-unresolve-with-rulers.red.txt`。）

同发的 S1：`... asks the authorization book 1 times, want 2 (found: [paths.go:196 key=f]) ...`（其余同上）。

**形三 · 两段同源（M-samesource）** —— 只有 c/d 两格看得见（断言 a 数出 2 枚，通过）：

```
paths_twocontainments_252_r2_test.go:321: TICKET 252 BAN RED (assertion c: the two containments must not share one source): both keys are built by the same chain of calls ([foldPath resolvedForm]) - two calls, one check. Measured at anchor e16e3037 this shape passes all 257 cases of this package (M-samesource), which is why the count alone is not enough.
paths_twocontainments_252_r2_test.go:335: TICKET 252 BAN RED (assertion d: one lexical leg and one re-resolving leg, exactly): found clean=0 tainted=2 (keys "foldPath(rf)" / "foldPath(rf)", tainted-name set [ok rf], guards [f foldPath ok p rf roots rootsContain string workspace]). The pair exists because the lexical spelling and the resolved spelling can name different trees; if both legs read the same one, only one tree is ever being asked about.
```

### 3.3 三枚尺各自的正控（未突变码上，逐字读数）

**W1**（`TestTicket252R2LexicalLegRefusesWhatOnlyItRefuses`，2 格全 PASS）：

```
B1 roots=LONG 直喂 SHORT+缺失叶子：Roots()=["d:\work\workspace\projects plans\wisp"]
   foldPath(canonical)="d:\work\worksp~1\projec~1\wisp\q252r2-lexical-leg-only.md" LEG1=false
   resolvedForm="D:\work\workspace\projects plans\Wisp\q252r2-lexical-leg-only.md" ok=true LEG2=true
   InAllowlist=false            ← 词法段单独在拒
   同一 raw ask 走 Canonicalize 后：LEG1=true LEG2=true InAllowlist=true   ← 正控，回合可用
B2 直喂 SHORT 的真目录本身：   LEG1=false LEG2=true InAllowlist=false
   走 Canonicalize 后：         LEG1=true LEG2=true InAllowlist=true
```

⇒ 复认票面 AC#1 那格的方向（响在第一段），并复认 `252-p1` 的"分岔变量是存在性"一句**对 `InAllowlist` 这一层不成立**：
B2 里整条路径都存在、也拿得到短名，两段照样一真一假 ⇒ 真正的变量是**这一次调用有没有先过 `Canonicalize`**。见 §7 R2-4。

**W2**（`TestTicket252R2ResolvedLegRefusesWhatOnlyItRefuses`，PASS）：

```
LEG1 lexical        : true
resolvedForm        : "" ok=false            ← Windows 对穿过 junction 的路径报"答不了"，不是"答在根外"
LEG2 resolved       : false
InAllowlist         : false
link target written : ...\outside（mklink /J 成功：Junction created for ...\proj\esc <<===>> ...\outside）
```

⚠ 一条现量纠我自己写在 §2 的话：**Windows 上 W2 的 LEG2 假是 `!ok` 那一半给的**（`EvalSymlinks` 走 junction 直接失败），
POSIX 上才是"答在根外"那一半。两半在同一条 `if` 里（`paths.go:209`），所以这一枚尺对两平台都有牙，
但它的红句不许被读成"解析器看到了根外"——消息里已按"Either half"写。

**S1**（`TestTicket252R2BothContainmentsReachable`，PASS） census 逐字：

```
S1 scanned 23 non-test .go files of this package; InAllowlist is declared at paths.go:189
S1 carrier: receiver="p" parameter="canonical" decl=paths.go:189
containment paths.go:196 key=f            calls=[foldPath]           tainted=false refuses=true
containment paths.go:209 key=foldPath(rf) calls=[foldPath resolvedForm] tainted=true  refuses=true
S1 re-resolution names: [ok rf]
S1 names an if-with-return-false consults: [f foldPath ok p rf roots rootsContain string workspace]
```

⇒ 工作区段（`paths.go:215`，arg0 是 `[]string{p.workspace}`）与 `paths_workspace.go:113` 那枚（不在本函数体内）
都被自动排除在枚数之外，实测 2 枚、不是 3 枚。

## §4 恒真自查（判据换反形／恒真化，两向都实测）

尺件都从 `.scratch/wisp/probes/252/r2/mut/r2-*.go.pristine` 还原（md5 `28b7fafc…`／`57256d99…`，终态已复认）。

1. **反形（三枚判据各换成相反断言，⛔ 未突变码）** —— 读数 `mut/inverted-predicates-pristine-code.log`：
   `0 PASS／5 FAIL`（这一发只用 `-run TestTicket252R2` 选了三枚尺，5 行＝W2 顶层＋S1 顶层＋W1 顶层＋B1＋B2）。
   逐名：`TestTicket252R2ResolvedLegRefusesWhatOnlyItRefuses`、`TestTicket252R2BothContainmentsReachable`、
   `TestTicket252R2LexicalLegRefusesWhatOnlyItRefuses`（＋两格子名）。
   三枚改法：W1 `if l.final` → `if !l.final`；W2 `if got` → `if !got`；S1 断言 a `len(book) != 2` → `len(book) == 2`。
   ⇒ **三枚都不是恒真式**：判据反过来，一字未改的生产码立刻红（S1 的反形红句甚至自己打印出
   "asks the authorization book 2 times, want 2"＝自证它响的是判据而不是世界）。
2. **恒真化（把两枚尺的判据换成永假＝尺失效，再重跑 M-merge）** —— 读数 `mut/m-merge-tautologized-rulers.log`：
   **`262 PASS／0 FAIL`**。
   ⇒ M-merge 的那 4 枚红**只由 W1 与 S1 供**，不是包内别的东西顺手拦的（同一条突变在未加尺时是 257／0）。
3. **负向尺配正控**：W1 与 W2 每一格都自带"同一条 ask 走完整生产回合必须 true"的正控
   （W2 另加"根内普通文件"那一格），所以它们不会因为判官整体变严、或因为根被丢弃而假绿。

## §5 门禁读数（终态，全部现跑；⛔ 无一行引用空凭据）

| 尺 | 读数 |
|---|---|
| `go build ./...` | **rc=0**（终态 11:0x）。⚠ 同一条尺在 10:4x 复跑为 **rc=1**，报错全在 `cmd/wisp`（`panel_host_windows.go:239 declared and not used: gh` 等），那一刻 `git status --porcelain -- cmd/wisp internal/config` 列出 2 枚 ` M cmd/wisp/*.go`＝`255-r2` 在飞的脏面，与本腿无关；两枚数都记下，⛔ 不许只留好看那一枚 |
| `go build ./internal/...` | rc=0（本腿射程内的产码面） |
| `gofumpt -l internal/tools` | **空输出，rc=0**（`gofumpt v0.12.0 (go1.27.1)`，取自 `$(go env GOPATH)/bin/gofumpt.exe`；两枚新件先 `-l` 命中、`-w` 修形后再 `-l` 为空） |
| `sh scripts/d22scan.sh` | **rc=0，clean**。分母逐字：bans #1-5 internal/=228、cmd/=37；ban #6 frontend/=85；ban #7 internal/tools/=23；ban #8 design/=39、frontend/=85、internal/=494（含 `_test.go`）、cmd/=91。全文 `final-d22scan.log`（前一次跑在 `d22scan.log`，同为 clean） |
| `go test -count=1 ./internal/tools/` | 终态 **262 PASS／0 FAIL／0 SKIP**，`ok … 15.276s`，rc=0 |
| `go test -count=1 ./internal/risk/` | `ok … 5.367s`（本腿一枚没跑、一字未碰；起手为 `ok 5.221s`） |
| 跨平台可编 | `GOOS=linux GOARCH=amd64 go vet ./internal/tools/` rc=0；`GOOS=darwin GOARCH=arm64` rc=0 ⇒ portable 那枚尺（W2＋S1）在非 Windows 构建里也在 |
| 名册逐名 comm | 起手 257 行（`roster-baseline.txt`）→ 终态 262 行（`roster-final.txt`）；`comm -23` **丢失 0 名**；`comm -13` 新增 5 名（`roster-new.txt`）：`TestTicket252R2ResolvedLegRefusesWhatOnlyItRefuses`、`TestTicket252R2BothContainmentsReachable`、`TestTicket252R2LexicalLegRefusesWhatOnlyItRefuses`（＋其 `/B1_…`、`/B2_…` 两格） |
| 还原 | 起手 `git status --porcelain -- internal/tools internal/risk`＝**0 行**；终态（尺件已入库后）同一把尺＝**0 行**。产码突变全还原：`internal/tools/paths.go` md5 `7a86da7aeba8420639491cd252a47a0c` ＝起手副本逐字节相等；两枚尺件 md5 亦与 `.scratch` 快照相等 |
| 锚点与 commit | 起手 `git log -1`＝`e16e3037…`；骨架 `1289d8bb`；尺件 `70a935ce`。⚠ 共享树 HEAD 在跑期间从 `e16e3037` 漂到 `b368380a`→`17cfc416`→…，本腿自证没被吞：`git merge-base --is-ancestor 1289d8bb HEAD` rc=0，且 `git diff --stat e16e3037 HEAD -- internal/tools internal/risk`＝**空**（同包无第二枚写腿） |
| 占位符自查 | 0 命中（命令与字符类等价式见 §6.5，⛔ 尺的字面文本没写进被扫文件） |

**收尾复跑（字都写完之后再跑一遍全套，读数不变才算交件）**：`go build ./...` rc=**0**、
`gofumpt -l internal/tools` **0 行**、`sh scripts/d22scan.sh` rc=**0**（clean）、
`go test -count=1 -v ./internal/tools/` **262 PASS／0 FAIL／0 SKIP**（`ok 14.892s`）、
名册 `comm -23` 对起手 257 枚 **丢失 0 名**、与上一轮终态名册 `diff` **0 行**（＝字写完整包没被改动带歪）；
`git status --porcelain -- internal/tools internal/risk`＝**0 行**（此刻脏的只有 `cmd/wisp` 四枚，`255-r2` 的面）。
件：`final2-build-all.log`／`final2-gofumpt.txt`／`final2-d22scan.log`／`final2-v-test.log`／`roster-final2.txt`（未进仓，同上）。

**证据件的落点（哪枚进了仓、哪枚只在本地）**
- 已随本腿 commit 进仓：本件 `impl.md`、`placeholder-hits.txt`（0 字节＝0 命中的那一次输出）、
  `roster-baseline.txt`（257 名）／`roster-final.txt`（262 名）／`roster-lost.txt`（0 字节）／`roster-new.txt`（5 名）、
  `final-d22scan.log`、`mut/*.red.txt`（五发抽出的红句全文，逐字未折行处比本件 §3.2 更全）。
- ⛔ 没进仓（体量大，留在 `.scratch/wisp/probes/252/r2/`，只建不删）：10 枚整包 `-v` 日志
  （`baseline-v-test.log` 116,114 字节、`final-v-test.log` 124,323、`with-rulers-pristine.log` 124,340、
  `mut/m-*.log` 各 116–127 KB、`mut/m-merge-tautologized-rulers.log` 124,187）与三份 `.pristine` 副本
  （`paths.go.pristine` 16,954／两枚尺件副本）。§5 表里凡引用这些文件名的地方都同时抄了读数本体，
  ⇒ 核销不需要它们，⛔ 别因为看不见日志就把某一枚数当空凭据。

## §6 判不动的地方（射程边界＋脆性自陈，⛔ 不藏）

**6.1 S1 会假响的合法重构（命中面，实测形状非推测）**
- 把两段包含**下沉进一枚 helper、`InAllowlist` 里只留一次调用**：枚数 1 → 红。这一发**算不算假响要人来判**：
  它确实只剩"对授权本的一次取用"，与硬禁的字面冲突，红是有理的；但若有人主张"两段还在、只是搬走"，
  那就得改尺（改尺＝人工批准）。
- 把 `rootsContain` **内联成循环**（不再对 `p.roots` 取用）：枚数 0 → 红＝**真假响**。
- 把授权本字段 `roots` **改名**：`readsBook252r2` 认不到 → 枚数 0 → 红＝**真假响**。
- 把 `InAllowlist` **挪去别的包**：`allowlistJudge252r2` 直接 `t.Fatalf`（消息明说"尺必须跟着人搬，⛔ 不许悄悄删"）＝**响亮假响**。
- 红句里的 `paths.go:196 / paths.go:209` 是**给人看的提示**、不是判据（判据用 `fset` 现取的行号，`found:` 列里就是真身）。
- ⛔ **不会**假响的（刻意做成这样的）：`rootsContain`／`foldPath`／`resolvedForm` 三枚函数改名（尺不认函数名，
  `calls` 集合只用于比较两枚键是否同源）；把两段并成 `if !A || !B { return false }`（仍是两次包含、各自能拒，
  **不算违禁**，尺正确地绿）；接收者 `p` 与入参 `canonical` 改名（从 AST 现读）；工作区段那枚
  （arg0 是 `p.workspace`，天然不计入）。

**6.2 W1 的边界**：只在 `windows` 构建里存在（短名只有 Windows 有）。拿不到可用别名时它 **FailF**
（沿用 `shortName252` 的"不许 skip、不许手打 `~1`"规矩），所以：在关掉 8.3 生成的卷上它是**红**、不是绿也不是 skip。
⇒ 编排者读红名册时要分得清"前置缺失"与"违禁命中"：前者首句逐字 `前置条件缺失`，后者首句逐字 `TICKET 252 BAN RED`。
这也是本腿**保留 S1** 的硬理由：S1 不需要任何 OS 合作，跨平台都在（§5 那两枚 vet 就是它的存在证明）。

**6.3 W2 的边界**：要能建目录链接（Windows `mklink /J` 不需管理员，107b 先例；POSIX `os.Symlink`）。
建不成 → `makeDirLink107b` 自己 FailF；链接建了但本平台解析器**既不拒也不指到根外** → W2 自己
`carrier moved` FailF。⇒ 它不会"静默地不测"。

**6.4 三枚尺都不覆盖的（具名，⛔ 别当已完成）**
- 票面 **AC#4／AC#5／AC#6**：三枚 witness 在 `cmd/wisp` 面（本腿一枚没跑，那是 `255-r2` 的脏面）；
  AC#5 的 `grantLinesOf` 属 `cmd/wisp`；AC#6 owner 真机那一发要 `[fs] allowed_dirs` 配置＋真 `wisp run` 的审计行。
- **`!ok` 语义被改软**（"解析不到＝未授权"那一支）：属票 102／107 地界，既有 N7／probe C 有牙
  （M-unresolve 那一发就是它俩先响的），本腿没动它、也不自称接管它。
- **折叠本身被改坏**（比如让 `foldPath` 把 `~1` 视为普通字符、或 `rootsContain` 改成不认组件边界）：
  那是"一次包含准不准"的尺，不是"两次包含在不在"的尺，⛔ 不在本票硬禁射程（边界属票 107／174）。
- `paths_workspace.go:106` 的 `inRoots` **只有一条腿**，但它是票 92 的"候选工作区必须在根内"那一问，
  **不在本硬禁的字面射程内**（本禁说的是允许根判定那两次包含）。⇒ 只有当生产**判定**改走它时才算越线，
  那一发本腿实量过＝M-inroots（257／5 红）。
- **谁在 `InAllowlist` 之外另造一条授权路径**（例如 `bridge.go:1145` / `task.go:838` 改成自己判根）：
  三枚尺都看不见——它们钉的是 C19 那一枚判定函数的形状。这一格要钉住得另派一枚消费者侧尺，⛔ 本腿不冒充覆盖。

**6.5 占位符自查的跑法与读数**：⛔ 这一节**不写**那把尺的字面量（写进来就命中自己——票面 `:47` 那格编排者自己踩过一次，
本腿按描述形写）。跑法＝在命令行上以**字符类拼法**组一枚正则（四枚英文占位大写词各把末字母放进方括号、
两枚中文残句词各拆成两段中间夹通配），对 `impl.md` ＋两枚尺件跑一次 `grep -nE`，再 `| wc -l` 收尾。
⇒ 终态读数 **0 命中**（件名 `placeholder-hits.txt`，同一条命令的空白括号一类**不计**在内：那枚模式会命中
Go 源码里的 `[]string{…}` 与 `x()`，第一次跑就是被它误报 39 行，已按"只数占位标记"重跑）。
另外三处自查同跑，读数逐枚具名：
- 两枚尺件里 `t\.Skip\(`（**调用形**）＝ **0／0**；⚠ 但按字面串 `t.Skip` 数是 `0`（portable）与 **`1`**（windows），
  那一枚命中是**注释里的英语句子** "Zero t.Skip calls in this file"（`_windows_test.go:32`）——
  与 `AGENTS §1.2` 记的"注释豁免、字符串不豁免"是同一类差别，⛔ 别把两种数法混用。
- 跟踪文件里 `MUTATION`／`EXPERIMENT` 字样 ＝ `paths.go` **0**、两枚尺件各 **0**（突变与反形只留在 `.scratch` 副本）。
- 宽字符带（`U+2190–21FF`／`U+2200–22FF`／`U+2300–23FF`／`U+2460–24FF`／`U+2600–27BF`／`U+2B00–2BFF`／`U+FE0F`／
  `U+1F000–1FAFF`）在三枚文件里 ＝ **0 命中**（⛔ 连"仪器抓不到的箭头"都没用）。
- 两枚尺件里 `t.Fatal`／`t.Fatalf` 共 **16** 枚（portable 12＋windows 4），全部落在"载体不动了就响"那一类。

## §7 推翻清单（编排者派单每一句都是待验断言；复跑结果逐条具名）

约定：**复真**＝本腿自己跑到同一读数；**推翻**＝不成立，写出真身；**不成立但无害**＝字面错、结论不受影响。
票面与派单原话⛔ 不改一字，这里只就地打旧。

| # | 派单／票面原话（摘要） | 复跑结论 | 真身／修正 |
|---|---|---|---|
| R2-1 | "internal/tools 整包 **257 PASS／0 FAIL**"，并给了抄名册尺 `^(--- )?(PASS|FAIL|SKIP)` | **数复真、尺不匹配（推翻其可操作性）** | 257 只在含子测试时成立（`^[[:space:]]*--- ` 计数 257）；派单那枚正则同一份日志数出 **187**＝186 顶层＋1 行裸 `PASS`。本腿两枚都记：名册比对用 257→262 那枚 |
| R2-2 | "非实现者验收表 … 154 行 … 最末 commit `e047940b`" | **复真** | `wc -l docs/evidence/s1/252-allowlist-sameform-v1.md` = 154；`git log --oneline -- <该文件>` 最末＝`e047940b docs(evidence/252-v1): 渲染修正…` |
| R2-3 | "修法 `2b1a3071`／载具 `34d89191`／证据件 `ebd5da2b`" | **复真** | 三枚都在 `dev` 上，标题逐字对得上（09:25 载具 → 09:28 修法 → 09:36 证据件） |
| R2-4 | "`internal/risk` 亦 `ok`"（当作 M-merge 无牙的第二处证据） | **数复真、**推理**推翻** | `ok 5.221s`（M-merge 下亦 `ok 4.050s`）。但 risk 包**不实现** `PathCanonicalizer`：`InAllowlist` 的实现全仓只有一枚（`internal/tools/paths.go:189`；`assessor.go:141` 是接口声明，`rules_test.go:29` 是测试桩）⇒ risk 绿**不提供任何**关于这两段包含的证据，删掉词法段它本来也不可能响 |
| R2-5 | "原因＝折叠保证两侧同形之后，包内那把 `measure252` 的 `final == leg1 && leg2` 恒等式**失去鉴别力**" | **成立但不完全（补一枚必要条件）** | 恒等式失鉴别力只是**一半**：另一半是**包内没有一枚用例把"两侧不同形的串"递给 `InAllowlist`**——252-r1／107／107b／102 每一枚取值用例都从 `Canonicalize`（或 `judge107b`）出发。⇒ 修法不是去修那枚恒等式（恒等式左右都是测试自己算的，怎么改都自证），而是**造分歧输入**（W1）＋一把不按取值取证的尺（S1）。⛔ 别把这句读成"包内一枚直调用例都没有"：`bridge_test.go:637` 就直调两次（`InAllowlist(root)`／`InAllowlist(outside)`），但那两枚串**两侧同答**（root 两段都真、outside 两段都假），构不成分歧，所以它抓不到 M-merge（本腿现量） |
| R2-6 | "两形各一枚…拆掉任一段（词法段 或 `resolvedForm` 段）都必须当场变红"＋"它真的把词法那一条腿删掉…今天全仓零仪器" | **"零仪器"仅对词法段成立（半推翻）** | `resolvedForm` 段**本来就有牙**：M-unresolve 在未加本腿尺时是 254 PASS／**3 FAIL**（`TestTicket252R1AlignmentAddsNoAuthorization/N7` 的 `AC#2 RED (widening)` ＋ `TestTicket107bProbeCLinkInsideAllowedRootStaysOutside` 的 `AC#3 RED`）。⇒ 本腿三枚尺里，**W1＋S1 是补缺**（词法段此前零仪器），**W2 是加固**（把该段的独立性从"根外链接"一枚扩成"名册里一枚具名 witness"）。⚠ 顺带收窄我自己 §3.1 那句：M-samesource 也被 W1 抓到（同源后词法段不再独立拒绝），S1 的独特贡献是**枚数之外那一层**与**不依赖 OS 合作**，不是"只有它看得见这一发" |
| R2-7 | 票面 AC#1／裁定节引的 `paths.go:133 / :140 / :152 / :153` | **过期（推翻，就地打旧）** | 现量：`InAllowlist` = `:189`、词法段 = `:196`、注释块 `:199-207`、`resolvedForm` 段 = `:208-209`、工作区段 = `:215`、`rootsContain` = `:224`、`sameFormOfUnresolved` = `:169-184`。原因＝`2b1a3071` 在 `Canonicalize` 里插了同形那一步（`:127-148`），整段下移。票 252 面与裁定节里的旧号⛔ 不改字，本节即真身 |
| R2-8 | "编排者给的形＝按调用图而不是按返回值"（并把"断言对 `rootsContain` 的直接调用枚数"当例） | **接受一半、推翻"只按调用图"（半推翻）** | 本腿交的是**三枚并存**：调用图尺 S1 只按 (receiver 类型, 方法名) 认函数、⛔ 不硬编码 `rootsContain`／`foldPath`／`resolvedForm` 三枚函数名（改了名也照绿），"再解析步"用**多返回值绑定**这一形状认。理由：纯结构尺对 M-merge／M-unresolve 只能报"枚数不对"，报不出"那一段真有独立拒绝能力"（这一半只有 W1／W2 的分歧输入能证）；而它有实测到的假响面（§6.1）。"直接调用枚数"这一枚本腿**改写成"对授权本取用的包含调用枚数"**（`X.roots`，或绑到它的一跳局部量），工作区那枚与 `paths_workspace.go:113` 那枚都不计入——这条改写是让尺不响错的关键，见 §3.3 的 2 枚读数 |
| R2-9 | "`→` 可以出现在界面文案里（仪器抓不到）"一类对仪器的过期描述（AGENTS §1.2），与本腿无直接关系但影响写码 | **复真（对尺件本身无约束）** | 两枚新尺件的注释⛔ 一个 emoji、⛔ U+2190–U+2BFF／U+2200–U+22FF 带内的字符（消息里一律 ASCII `->`）；`sh scripts/d22scan.sh` rc=0，ban #8 覆盖 internal/ 494 枚含 `_test.go`。⇒ 本腿没有"靠豁免过关"的地方 |
| R2-10 | "`design/**` 此刻在工作树里是被别人删了未 staged，`internal/panel` 那 2 枚 colour/token 红与此有关，不归你" | **前半复真、后半不背书** | 起手 `git status --porcelain` 确实列出 ` D design/assets/base.css`、` D design/screens/*.html` 等一批未 staged 删除（本腿一枚 `internal/panel` 都没跑，⛔ 也不引 `design/**`／`frontend/**` 内容）。"2 枚红与此有关"这半句本腿**既不支持也不推翻**：无读数 |
| R2-11 | "本腿只加仪器；写面 `internal/tools/**`；生产码一字不许动" | **复真（自证见 §5 还原行）** | `paths.go` md5 与起手副本相等、`git status --porcelain -- internal/tools internal/risk` 终态 0 行；新增只有两枚 `*_test.go`（commit `70a935ce`）＋本 probe 件 |

一句总结给编排者：**票面那句硬禁从今天起有牙了**——三枚尺、五发读数、两形各自红；
但"零仪器"这顶帽子只该扣在**词法段**上，`resolvedForm` 段此前已有两枚（N7／probe C），
下一程写结案时⛔ 别把两段的账合并成"本票新造三枚"。

