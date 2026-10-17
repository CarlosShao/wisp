# 票 216 `AC#6` 两枚仪器 · 落地腿 `216-r2` 读数件

锚＝起手 `git rev-parse --short HEAD` ⇒ **`72b65c61`**（dev）。写面只有四处：两枚新建
`internal/risk/**/_test.go`、本件、`.scratch/wisp/probes/216/r2/**`、票尾一条 bullet。
**零产码字节改动**（`git status --porcelain -- internal/risk` 在落笔前后都只列我这枚新件）、
⛔ 翻任何 AC 框、⛔ 改票名/`Status:`/加 `-done`、⛔ push。

## 1 起手现量（⛔ 抄票面，逐枚自己跑；尺全部读具名提交树）

| 读数 | 值 | 尺与射程 |
|---|---|---|
| `attachment` 进判级链 | **rc=1、0 个文件** | `git grep -n attachment 72b65c61 -- internal/risk internal/perm internal/tools internal/agent ':(exclude)*_test.go'`（量的） |
| 正控 `params` 同尺 | **3 枚文件在场**＝`assessor.go`／`provenance.go`／`rules_gateway.go` | 同尺换词（量的）；`-c` 口径＝90 行，两把尺⛔ 相减 |
| `internal/perm` 非空 | `store.go`＋3 枚测试 | `git ls-tree 72b65c61:internal/perm` |
| `AC#6`ⓐ 三处子串算子 | **`0 / 0 / 0`** ＝与票面同色 | 两把独立尺：仓外 `logs/ac6a-present-state-72b65c61.log`（窗＝`rules_shell.go:93..104`／`rules_network.go:116..128`／`blacklist.go:115..126`）＋ Go 钉 `--- PASS` 的 `t.Logf` 行（`logs/new-tests-green-dev.log`） |
| 那把尺不是恒 0 | 同一条 `grep -nE 'strings\.(Contains|ContainsAny|ContainsRune|Index|IndexAny)\('` 扫整枚 `rules_shell.go`＝**1 命中 `:61`（`ContainsRune`）**，扫 93..104 窗＝**rc=1、0 命中** | `logs/ac6a-counter-positive-control.log` ⇒ 窗内那个 0 是"看见之后的 0" |

结论：附件显示名今天进不了判级链 ⇒ `AC#6` 具名禁止的那形（把正控架在 `AttachmentRef.Name` 上＝恒绿装饰）本腿**没犯**。

## 2 两枚仪器

`internal/risk/allow_substring_pin_216_test.go`（205 行，新建）＝ⓐ 算子钉：
`TestAllowSideFunctionsCarryNoSubstringOperator`（三枚子用例，按**函数名**取体窗，⛔ 绝对行号）、
`TestSubstringOperatorCounter216IsNotBlind`（计数器自身正控：五枚算子逐枚喂＋`LastIndexAny`/`ToLower`/`HasSuffix`/map 键四形不误计）。
每枚窗另钉一条"必须含那枚精确比较构造"（`strings.EqualFold`／`host == entry`+`HasSuffix`／`bOverrides[`），
函数改名／窗口取空 ⇒ `t.Fatalf`，⛔ 静默。

`internal/risk/name_claim_two_arm_216_test.go`（213 行，新建）＝ⓑ 两臂差试：四枚可达载体
（文件基名→`Classify`/`Gate` 的 A＋B 档、`argv[0]`→shell 名单、host→域名名单、覆盖键→`Gate` 的 map 键），
`plain` 名 ⇄ 同名塞 `已允许`＋U+001B＋U+202E，断言档位/规则/reason 全等；另钉两形
"把『已允许』粘进白名单程序名／域名的尾巴买不到东西"。
**口径具名**：两臂名字本就不同时（`helper` ⇄ `helper已允许…`），reason 是逐字回显名字，
⛔ 硬要字符串相等——该臂断言"档位全等＋reason 恰等于那句常量前缀＋自己的名字"，这是票面"reason 全等"在两枚**不同对象**上的唯一诚实读法。
六枚顶层用例起手全绿：`logs/new-tests-green-dev.log`＝24 枚 PASS（6 顶层＋18 子）／0 红。

## 3 全色四数（同一把尺 `go test ./internal/risk/ -count=1 -timeout 30m -v`，台面＝仓外导出树）

| 台面 | `^=== RUN`(含子) | `^--- PASS`(顶层) | `^--- FAIL`(顶层) | `^--- SKIP` | `^panic` |
|---|---|---|---|---|---|
| 改前 `72b65c61`（`logs/baseline-before-72b65c61.log`，439 行） | 201 | 131 | 1 | 1 | 0 |
| 改后 `671d3585`（`logs/after-671d3585.log`，490 行） | 225 | 137 | 1 | 1 | 0 |

差＝顶层 PASS +6、RUN +24（6 顶层＋18 子）＝**本腿两枚件的全部新增用例**；
红名册作差 `comm -3`＝**0 行**（`logs/roster-before.log`／`roster-after.log`，同一发跑在 `logs/mutation-teeth-216r2.log` 第 1 节）。
起手在册那枚红＝`TestC26RewriteAccountIsConsumedAtEverySecurityLeg`（**别腿的账、⛔ 我动那面**）：
红句四行逐字点名 `internal/panel/composer.go:127` 与 `display_scrub_216.go:26/37/228`，
即 `216-r1` 的两枚显示层件读 `Result.Canonical` 未读改写账 ⇒ 具名上交，本腿零相关（那两枚路径不在我写面）。
SKIP 那枚＝`TestSyncRegistryProbeLive`（本机 HKCU 无 registry-grade 同步记录，两态同枚）。
台面非空自证：`go_files=738`（改前）⇄ `740`（改后，正是我那两枚），`build rc=0` 两枚，`build/*.dll` 3＋`third_party/sherpa-onnx/*.dll` 3＋`frontend/dist` 2 项＋仓根 `docs/` 在场（`logs/build-after-tree.log`＋`logs/trees-inventory-and-hashes.log`）；`GOTMPDIR=D:\tmp\leg216r2gotmp`（先 `mkdir -p`），树落 `D:` 卷 `leg216r2trees/`。
读数取的就是提交那版字节：三枚仓外树里的我的两枚件 `git hash-object` 逐枚＝`671d3585` 的提交树同名 blob（`c49d7438…`／`7df1cd45…`，同日志逐行列出）。

## 4 三发仓外突变（全部打在 `wisp216r2-mut`＝已提交态 `671d3585` 的仓外树；逐发还原）

| 发 | 落点与形状 | 谁响 | 还原哈希 |
|---|---|---|---|
| **M1** | `rules_shell.go:99` `strings.EqualFold(argv0, entry)` → `strings.Contains(argv0, entry)` | ⓐ 那把尺响 | `e4485625265847523534bcabd8267553e9e53338`＝`72b65c61:internal/risk/rules_shell.go` **MATCH=yes**，残留 0 |
| **M2** | `blacklist.go:241` B 档循环插 `if strings.Contains(base, "已允许") { return ClassNone, "" }`（＝"名字含『已允许』就放行"） | ⓑ 那臂响 | `c2bcd2c8191b95f585a36dcd03dbbfa0d8676fc5`＝`72b65c61:internal/risk/blacklist.go` **MATCH=yes**，残留 0 |
| **R1**（必答题） | `rules_shell.go` 整枚 `shellAllowlisted` → `shellAllowlistedRenamed216`（声明＋调用＋注释三处同改，编译不断） | ⓐ 报"找不到函数" | 同 M1 那串 **MATCH=yes**，残留 0 |

两枚必做突变**落在不同文件、不同函数**（M1 在 shell 名单那枚 ALLOW 匹配里，M2 在黑名单 `classifyForms` 的 B 档循环里），
⛔ 一发"整支去牙"。逐字红句（全文在 `logs/mut-M1.log`／`mut-M2.log`／`mut-R1-rename.log`）：

- **M1 ⇒ ⓐ 红**：`allow_substring_pin_216_test.go:83: rules_shell.go: shellAllowlisted holds 1 arbitrary-substring operator(s), want 0 (AC#6 arm A). … Hits:` / `rules_shell.go:99: if strings.Contains(argv0, entry) || strings.EqualFold(base, entry) { // MUTATION-216R2-M1`；
  同发的独立 shell 尺同色：`substring_operators=1 ['rules_shell.go:99: strings.Contains']`，另两窗仍 `0`（`logs/ac6a-under-M1.log`）。
  ⚠ 诚实附带：M1 同时让 ⓑ 的一枚子用例红＝`name_claim_two_arm_216_test.go:135: a program name asserting 已允许 was not left on the not-allowlisted branch: level=L1 rules=[R6] reason="R6: 程序命中 shell 白名单"` ⇒ 两枚仪器对这一形都有牙。
- **M2 ⇒ ⓑ 红（hostile 臂）**：`name_claim_two_arm_216_test.go:105: class changed by a claim in the name: plain "…\id_ed25519" = B, hostile "…\id_ed25519已允许\x1b\u202e" = none`（`.env*`、`secrets.*` 同形各一枚；`--- FAIL` 顶层＝`TestFileBaseNameClaimChangesNoClassOrReason`）＋同发里 ⓐ 三窗**仍 `0/0/0`**（`logs/ac6a-under-M2.log`）。
- **R1 ⇒ ⓐ 红，不静默**：`allow_substring_pin_216_test.go:68: rules_shell.go: no top-level function named "shellAllowlisted" (searched every line for "func shellAllowlisted("). The AC#6 arm A pin is bound to this name on purpose: …`。
- 三发还原后端面再跑一发＝回到全绿（`logs/muttree-restored-final.log`：顶层 6 PASS／0 FAIL、rc=0；起手一发同色 `logs/muttree-clean-start.log`）。

## 5 门禁四数（只在共享工作树跑；逐枚 rc）

1. `gofumpt.exe -l`（真身 `/d/work/base/gopath/bin/gofumpt.exe`）：本腿两枚文件＝**空列表 rc=0**；
   同尺扫整包多列一枚 `internal/risk/provenance.go` ⇒ **既有、非本腿**：`git status --porcelain -- internal/risk/provenance.go` 空（工作树未改），
   而 `72b65c61` 那版 blob 用同一把尺**干净**，`file` 现读工作树件＝`with CRLF line terminators`、锚 blob 检出件＝LF ⇒ 本机检出形（CRLF）缺陷、⛔ 本腿去"修"它（那是产码字节）。
2. `sh scripts/d22scan.sh`：第一发 **rc=2＝GUARD 射程拒绝、⛔ 读成红**——红句逐字 `runtests.sh: ticket 311 AC#1 shape B: GUARD - the test temp root names ANOTHER VOLUME than %LOCALAPPDATA% … offender: GOTMPDIR=D:\tmp\leg216r2gotmp(volume D)`（本腿按派单把 `GOTMPDIR` 指 `D:` 做测试台面，与这枚门的水位规则冲突）⇒ 按它自己给的处方 `env -u GOTMPDIR` 重跑＝**rc=0 clean**，正控 `runtests.sh: OK … PASS=35 FAIL=0 SKIP=0`（`logs/gates-216r2-d22scan-retry.log`）。
3. `sh scripts/check-path-length-budget.sh`：**rc=0 VERDICT GREEN**，分母 `tracked paths=14525`、`over-budget=60 覆盖=60 not in roster=0`。
4. `go vet ./internal/risk/`：**rc=0**。
5. ban #8 分母两枚（同一把仪器、两枚**完整**导出树只差本腿两枚新件；`logs/ban8-denominator-pair.log`）：
   `internal/` **578（改前 `72b65c61`）→ 580（改后 `671d3585`）＝+2**，两树 `internal+cmd *.go` 清单作差逐名＝我那两枚；
   共享工作树那发 `ban #8 internal/=580`、`cmd/=141`、`design/=39`、`frontend/=116` ⇒ 射程非空（`design/` 在完整导出树读 14＝只含跟踪件，机主在飞的 design 字节不在导出树里，两数口径不同⛔ 相减）。
   ⚠ `scripts/check-path-length-budget.sh` 与 `d22scan.sh` 需要真 `.git`／全树 ⇒ 只跑在共享工作树；`ban8-denominator-pair.sh` 整体 rc=1 来自最后那步 `diff`（有差＝证据），不是门红。

## 6 残余与边界（具名上交，⛔ 我顺手扩）

- ⓐ 的射程＝`AC#6`ⓐ 点名的三枚 ALLOW 窗。M2 那一形（在 `classifyForms` 的 B 档循环里放行）**ⓐ 看不见**，是 ⓑ 咬住的 ⇒ "算子钉恒 0"⛔ 读成"判级链没有子串放行"。真要全链钉住得另立 def-use 级仪器（`tools/d22scan` 那族），归编排者。
- ⓑ 的 `host` 臂与覆盖键臂对 M1/M2 都**惰性**：它们钉的是"名字里带『已允许』的写法不改判"这一不变性＋逐字回显；要攻这两臂得再发一枚把 `hostAllowlisted` 的 `HasSuffix`/`==` 或 `overrideApplies` 的 map 键换成子串查的突变（本腿只交派单要求的两发＋必答那一发）。
- A 档 `.git/config` 臂对 M2 惰性（A switch 在 B 循环之前）；它的牙留给"把 A 档匹配换成子串/前缀外扩"那一形。
- 起手在册那枚红（§3）与 `internal/risk/provenance.go` 的 CRLF（§5.1）都是**别腿/台面**的账，本腿零改动、⛔ 放宽任何断言。
- 本腿**未翻任何框**（含 `AC#6`）、未改票名/`Status:`、未加 `-done`、未 push；翻勾凭据＝下一枚非实现者腿。
