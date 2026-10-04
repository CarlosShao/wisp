# 212-r2 实现腿交付件（票 212 第 66 行「归 212-r2 的清单（两任合并）」四枚）

- 写腿：212-r2＝**实现者**（非验收者、非编排者）。票面
  `.scratch/wisp/issues/212-comments-cite-evidence-files-that-do-not-exist.md` §编排者终裁节 第 66 行。
- 两份独立判语（同一件事、两个非实现者）：`.scratch/wisp/probes/212/v1/verdict.md`（第三任满稿）与
  `.scratch/wisp/probes/212/v1/verdict-first-instance.md`（第一实例）。两任对本清单无分歧。
- 起手锚 `d934e016`；起手时刻 `2026-10-04 09:05 +0800`；分支 dev。
- **commit 分界（后续程必须按这两枚取基线）**：
  - 未修码基线＝`d934e016`（起手锚）与本腿骨架枚 `7ac8965a`（只落本件的七节标题，产码零字改）。
    这两枚之上的 `-self-test` 分母是 **36**。任何把 36 当"修后读数"的程都是把改前当基线。
  - 修码枚＝`5413f46d`（五枚产码文件：`tools/d22scan/main.go`、`tools/d22scan/selftestsamples.go`、
    `internal/agent/tools.go`、`internal/tools/bridge.go`、`internal/agent/approval/pending_read.go`）。
    从这枚起 `-self-test` 分母是 **37（20 expect-ring / 17 expect-silent）**。
  - 本件与读数档＝第三枚（`git log --oneline` 里 `5413f46d` 之上那一枚），只动 `.scratch/**`。

## §0 起手锚与三把尺的起手读数（本腿自跑，不复认编排者读数）

编排者 2026-10-04 08:5x 给的三行起手读数，本腿在动任何码之前复跑一遍：

| 尺 | 起手逐字末行／rc | 与编排者读数 |
|---|---|---|
| `cd tools/d22scan && go vet ./...` | rc=0，无输出 | 相符（编排者未给这一把的末行） |
| `go run . -self-test` | `d22scan -self-test: clean - all 36 direction checks passed (20 expect-ring, 16 expect-silent)`，rc=0 | **逐字相符** |
| `go test -count=1 ./...` | `ok  	github.com/CarlosShao/wisp/tools/d22scan	22.248s`，rc=0 | 同一判语，秒数不同（本机此刻负载；分母与 PASS 数才是要对的东西） |

全仓门禁起手一发（`sh scripts/d22scan.sh`，rc=0，档 `.scratch/wisp/probes/212/r2/logs/pre-gate-full.txt`）：

```
runtests.sh: OK - packages=[./...] top-level: PASS=34 FAIL=0 SKIP=0, === RUN=76, '[no tests to run]'=0
d22scan: examined 266 production Go files under internal/ and cmd/ of D:/work/workspace/projects plans/Wisp
d22scan: clean - no D22 ban violations; live scope work: bans #1-5 internal/=228, bans #1-5 cmd/=38,
ban #6 frontend/=85, ban #7 internal/tools/=23, ban #8 design/=39, ban #8 frontend/=85,
ban #8 internal/=498, ban #8 cmd/=97
```

- ⚠ **本腿基线取 `cmd/#8=97`，不取 A578 的 96**。第三任判语 §3 已把 96→97 具名归因于并发 258-v1
  腿 21:45 落的 untracked 测试件（ban #8 数 `_test.go` 且连 untracked 一起读＝A207 机器相依分母），
  本腿 09:05 复跑仍是 97 ⇒ 那枚件还在树上，不是漂移。
- 起手既有两枚 ban #9 样本的逐字读数（`logs/pre-selftest.txt`）：
  ring 向 `internal/probe/cites.go` 引 `docs/evidence/s1/212-citation-ruler.md` ⇒ `ring OK`；
  silent 向 `internal/probe/cites-ok.go` ⇒ `stayed silent while bans #1-5 + #8 internal/ Go walk examined 3 file(s); the fixture produced zero findings at all`。

## §1 四枚落点（逐字对应票面第 66 行）

| 枚 | 落点（改后逐字行号） | 改了什么 |
|---|---|---|
| ① 裁 ⓐ 修法 | `tools/d22scan/main.go:860-894`（新增 `shorthandRegionRe` 与 `shorthandPathStarts()`）、`:805`（`short := shorthandPathStarts(c.Text)`）、`:807-809`（命中即 `continue`）、`:791-800`（ban #9 块注释自述同步改准）；样本＝`tools/d22scan/selftestsamples.go:394-404`（`internal/probe/cites-shorthand.go`，wantSilent），`:31`（`glyphEllipsis = string(rune(0x2026))`，按本文件既有的"码位构字、字节保持 ASCII"定式） | `repoPathRe` 本体**一字未动**；豁免做在"区域起始偏移"这一层：区域类同 `repoPathRe` 的前缀表，但字符类多带 `…`(U+2026) 与 `*`，区域文本含 `...`／`…`／`*` 之一 ⇒ 该偏移起的 token 登记为名册③「缩写」，只规定、不问罪。为什么不能在 `repoPathRe` 里写排除：RE2 无前后瞻，而 `…`/`*` 根本不在它的字符类里——旧行为是把 `docs/evidence/s1/152…-accept-r1.md` **截断**成 `docs/evidence/s1/152` 再去 stat（比截不断更糟），把 `...` 形整吞再去 stat。 |
| ② symRefRe 自述改实话 | `tools/d22scan/main.go:814-823`（发射点内的行注释，旧 :802-807 那段）与 `:896-905`（`symRefRe` 的 doc 注释，正则本体在 `:906-907` 未动） | 旧文把 `internal/tool` 与两枚真 API 引用并列成"symRefRe 排除例"。实测：`symRefRe` = `^[a-z][a-z0-9]*(/[a-z][a-z0-9]*)*\.[A-Z]`，**要一个 `.大写` 才命中**；`internal/tool` 无点 ⇒ 不排除 ⇒ 走到 `os.Stat` 并定罪（两任判语的夹具 p2 同读数）。新文照实测写，并把"那枚 token 当年是靠改注释消红的"这一句钉进去。**没有为对上注释去动正则**（形 ⓑ 未批，口令「212 改 ⓑ」不生效）。 |
| ③ 三处半假 → 逐字对得上 | `internal/agent/tools.go:88-93`、`internal/tools/bridge.go:251-259`、`internal/agent/approval/pending_read.go:43-46` | 见本节下方三行解剖。三枚都**只改注释**，未动任何产码逻辑；三枚都引全路径（带目录），因为本票的尺就抓短引。 |
| ④ 计数口径重述 | 本件 §5 | 三个"10/11/8/9"各是什么数、怎么数、差在哪，全部用本腿自己现跑的读数重推一遍（不是抄判语）。 |

### ③ 的逐枚解剖（动手前先把 D37 真实名册量出来，见 §2）

1. **`internal/agent/tools.go:89`（改前逐字）**：`// timeout. The returned error is a host failure (error class internal-tool);`
   - 假在哪：`internal-tool` 既不是 D37 十七类之一，也不是 `internal/tool` 这条路径（后者不存在，正是 ban #9 抓的形状）。5e8748b3 把斜杠去掉、消了红，但**换了个仍不存在的名字**。
   - 实测真身（本腿自跑链路，不抄判语）：`internal/agent/loop.go:692-694` `case execErr[i] != nil:` → `log.ErrorClass = ErrorClassOfTurnError(execErr[i])`；`internal/agent/guard.go:271-277` 的 `ErrorClassOfTurnError` 先 `observe.ClassOf(err)`，**取不到类就 `return string(observe.ClassInternal)`**；`internal/observe/errors.go:230` 的 `ClassOf` 对未分类错误返回 `ok=false`（doc 在 `:227-229`），`:240` 的 `ClassOfOrInternal` 是同一回退的便利形。⇒ 落账名目＝**class `internal`**（常量 `ClassInternal`），不是任何带连字符的合成名。
   - 改后：说"落账走 `internal/agent/guard.go:271`，记的是错误自带的 D37 类，无类即回退 class `internal`（`observe.ClassInternal`）；D37 十七类里没有 `internal-tool`"。
2. **`internal/tools/bridge.go:252`（改前逐字）**：`// The returned error is reserved for host-internal faults (D37 class` / `// internal-provider), because the loop treats non-nil as "the provider broke"` / `// (loop.go:654-657) and books a class the model cannot self-correct against.`
   - 假在哪有两处：(a) `internal-provider` 同形问题（不在 §2 名册）；(b) 它引的 `loop.go:654-657` 是**短引＋错锚**——那五行是 `timeout := guard.PerToolTimeout()` 与 `req := ToolRequest{...}`，非-nil 语义真正落在 `:692-694`。
   - 改后：类名同 ①.1 的实测链路；引用改成 `internal/agent/loop.go:692-694`（全路径＋核过的行）；并写清"a provider fault 是 class `provider`（`observe.ClassProvider`），只有错误被按该类包装时才进得了那一行"。
   - ⚠ **本枚被仪器咬了一口，具名记录**：第一版我把 `internal-provider` 加了引号写进注释 ⇒ `sh scripts/d22scan.sh` 的 `TestScannerSelfScanOfRealRepoIsGreen` 当场红：
     `internal/tools/bridge.go:256: [internal-artifact-tool] host-internal artifact write looks implemented as a gated tool name (D34 note 2/D22)`。
     成因＝ban #7 的 `artifactToolRe`（`main.go:158`）是 `(?i)"(spill|internal[._-][a-z0-9_.-]+)"`，**带引号的 `internal-xxx` 就命中**，而 ban #7 走 `walkText(..., goOnly=true)` 且**注释不豁免**（＝selftestsamples.go 里钉着的 temperament A1）。⇒ 修法＝**改措辞**（把该名写成不带引号的散文），**不动判据**；改后三把尺全绿（§3）。这一口正是"注释诚实化"与本仓仪器的相互作用，值得留在案上。
3. **`internal/agent/approval/pending_read.go:43`（改前逐字）**：`// verbatim, ticket 17's frozen-contract note (enforced by tools/d22scan's pathresolver-bypass ban since the scanner landed), copied per`
   - 票面第 66 行那句"真身是 `internal/tools/gate.go:13` 只差一个前缀"——**本腿读那两个文件后判：这句成立**：
     `internal/tools/gate.go:13` 逐字＝`// renders RulesHit and Reason VERBATIM (ticket 17's frozen-contract note), so`，
     而 5e8748b3 改前的原引是 `tools/gate.go`（少 `internal/` 前缀，ban #9 在 pre-image 里就是抓这一枚，见 §5 的 8 枚真②名册）。
   - 但改后那句**把话题换成了执法**：`pathresolver-bypass` ban 管的是 `risk.PathResolver` 之外的 `filepath.Clean/Abs`
     （`tools/d22scan/main.go` 的 ban #2 发射点），它**不执行**票 17 的"逐字显示"契约 ⇒ 两任判语同判"语义嫁接＋路标丢失"。
   - 改后：路标接回 `internal/tools/gate.go:13`，并明写"这不是 tools/d22scan 执法的东西"（同一句里把两个 ban 的射程差别说清），不新增任何不成立的断言。

## §2 D37 十七类真实名册（现量，逐字）

推导式：`grep -n "accept these 17 values\|^const (\|^// The 17 classes\|^var allClasses" internal/observe/errors.go`
→ 契约句在 `:13`（"accept these 17 values and nothing else"），常量块＝`:17-36`（`:17` 是那句
"// The 17 classes, in D37 table order."），权威集合 `allClasses` 的 `var` 在 `:39`（注释在 `:38`，自称
"the authoritative set, in D37 table order"），校验器 `ValidateErrorClass` 在 `:61`。17 枚逐字：

```
config  auth  network  rate_limit  provider  model  audio_device  asr  tool
permission_denied  user_rejected  cancelled  budget  loop  injection  internal  resource
```

- **`internal-tool` 与 `internal-provider` 都不在这 17 枚里**（两任判语同判，本腿现量复认）。
  文件里唯一的两个连字符合成名，就是 §1 ③.1/③.2 改掉的那两处；改后全仓仅剩我写的那两句负断言
  （"…is not one of D37's 17 classes"），不再是凭据形。
- 计数校验：上面两行 9＋8＝**17**，与 `internal/observe/errors.go:13` 的注释
  （"accept these 17 values and nothing else"）自洽；`ValidateErrorClass`（`:61`）对表外字符串直接报错，
  所以"造一个类名"在运行时就进不了 `task_log.error_class`——这也是为什么注释里造名比代码里造名更隐蔽。

## §3 门禁读数（三把尺逐字末行＋新分母；AC#5「不许伤」的本腿自证）

改后全部重跑，档在 `.scratch/wisp/probes/212/r2/logs/`：

| 尺 | 逐字末行 | rc |
|---|---|---|
| `go run . -self-test` | `d22scan -self-test: clean - all 37 direction checks passed (20 expect-ring, 17 expect-silent)` | 0 |
| `go test -count=1 ./...` | `ok  	github.com/CarlosShao/wisp/tools/d22scan	18.889s` | 0 |
| `go vet ./...`（tools/d22scan 模块） | 无输出 | 0 |
| `go vet ./internal/agent/ ./internal/tools/ ./internal/agent/approval/` | 无输出 | 0 |
| `go build ./internal/...` | 无输出 | 0 |
| `sh scripts/d22scan.sh`（全仓门禁） | `d22scan: clean - no D22 ban violations; live scope work: bans #1-5 internal/=228, bans #1-5 cmd/=38, ban #6 frontend/=85, ban #7 internal/tools/=23, ban #8 design/=39, ban #8 frontend/=85, ban #8 internal/=498, ban #8 cmd/=97; ...` | 0 |

**新分母怎么来的（写清楚，别只写一个数）**：`selfCases` 表加了一枚 wantSilent ⇒ 36→**37**；
expect-ring 仍 **20**（我没加 ring 样本，票面也只要求补 silent 一枚）；expect-silent 16→**17**。
`-self-test` 末行的两个数由 `countWant()` 现数表得出，不是硬写。
`runtests.sh` 的顶层计数**不变**（PASS=34 FAIL=0 SKIP=0 / RUN=76）——新样本进的是
`TestSelfTestEntryPassesEveryCase` 的执行表，不是新的测试函数，所以"76 RUN"这一枚不该动，也确实没动
（`logs/post1-regex-gate-full.txt` 与 `logs/pre-gate-full.txt` 的该行逐字相等）。

**八枚 ban 逐名零漂移**：`diff logs/pre-scope-lines.txt logs/post3-scope-lines.txt` → **空**（八行 scope
＋`clean` 行＋`examined 266` 行逐字相等）。中间那一发失败的（`logs/post2-comments-gotest.txt` rc=1）
是 §1 ③.2 记的 ban #7 一口，不是分母漂移；改措辞后的 `post3` 一发全绿。

**gofumpt 卫生**：`gofumpt` 找得到＝`%GOPATH%/bin/gofumpt.exe`，版本 `v0.12.0 (go1.27.1)`。
- `gofumpt -l tools/d22scan/main.go tools/d22scan/selftestsamples.go` → 空（已 `-w` 过一次）。
- 三枚 `internal/**` 文件在**工作树**里被 `-l` 列出，但那是行尾：`.gitattributes` 写 `*.go text eol=lf`，
  本机 checkout 出来是 CRLF（git 自己也在 commit 时警告 `CRLF will be replaced by LF`）。
  把 LF 化后的字节（＝将入库的字节）单独喂给 gofumpt ⇒ **三枚全部 EMPTY**（脚本读数见 `logs/`）。
  ⇒ 判：我的改动 gofumpt 干净；工作树的 `-l` 命中是机器面行尾，不是代码。
- ⚠ **具名一处过期项（非本腿造成、本腿未顺手改）**：`tools/d22scan/selftest.go` 连 LF 字节都不干净——
  `git show HEAD:tools/d22scan/selftest.go | gofumpt -l` 命中，差异是 `selfSkeleton()` 里
  `"docs/readings.md"` / `"internal/probe/roster.md"` 两行的对齐（gofumpt v0.12.0 要去掉多出来的空格）。
  本枚一行没碰这文件（`git diff --stat` 里它不在）。它会不会让 CI 的 `gofmt (gofumpt)` 那一步红，
  取决于 CI 用的 gofumpt 版本；这属"格式门的版本相依"，不在票 212 的四枚里 ⇒ 只上报，不并入。

## §4 变异自证（判据有牙吗）

问题：① 的排除逻辑是不是真在起作用，还是一句写在注释里的希望。两台突变都由
`.scratch/wisp/probes/212/r2/mut_class3_exemption.py`（自测台）与
`.scratch/wisp/probes/212/r2/mut_class3_preimage.py`（pre-image 树）执行：
只把发射点那一行 `short := shorthandPathStarts(c.Text)` 换成 `short := map[int]bool{}`（＝摘回 2026-10-03 的形状），
`try/finally` 还原，跑完再量 md5。**全程不改 git、只动跟踪文件的一次瞬时内容，还原后 `git diff --stat tools/d22scan` 为空。**

| 发 | 目标 | 逐字末行 | rc | md5 |
|---|---|---|---|---|
| 红发（摘掉排除） | `go run . -self-test` | `d22scan -self-test: 1 direction(s) failed, 36/37 passed - the gate does not see what it claims` | **1** | — |
| 绿发（还原） | `go run . -self-test` | `d22scan -self-test: clean - all 37 direction checks passed (20 expect-ring, 17 expect-silent)` | **0** | before＝after＝`e6d1745996e8e02c82ab691ac2cbadf8`，`restored identical = True` |

红发那一行的完整罪句（`logs/mut-01-class3-exemption-removed.txt`，一枚样本响三次）：

```
d22scan -self-test: ban #9 phantom-citation    silent FAIL   CLASS 3 AS A PATH: "...", "*" and U+2026 inside a repo-relative token stay silent (ticket 212 裁 ⓐ) | rang on a sample that must stay silent:
  internal/probe/cites-shorthand.go:4: [phantom-citation] comment cites repo path "docs/evidence/s1/152-...-accept-r2.md" ...
  internal/probe/cites-shorthand.go:5: [phantom-citation] comment cites repo path "docs/evidence/s1/212-" ...
  internal/probe/cites-shorthand.go:5: [phantom-citation] comment cites repo path "docs/evidence/s1/152" ...
```

三枚 token 各对应一种缩写拼法（`...` 整吞、`*` 截断成 `212-`、`…` 截断成 `152`）——
这正是 ① 修的两条不同失败路径，摘掉排除后**三条都响**，所以"silent 样本"不是只练了一种形状。
同一发里 ban #1-#8 与 ban #9 的 ring 向全绿（`logs/mut-01-class3-exemption-removed.txt` 全量在档）。

**牙的方向也验了**（不是只验"会响"）：pre-image 树两发对比，见 §5 的"口径 C"——
排除在＝8 枚 distinct token，排除摘＝9 枚，**差集恰为 `docs/evidence/s1/62-`**，
其余 8 枚一枚不因①而变；同时 3 枚 API 形在两发里都静默（`symRefRe` 未被本腿触动）。

## §5 计数口径重述（票面"10 枚"与树不符的那一格：11 token／8 真②／9 hunk 各是怎么数出来的）

本腿不复用两任判语的数，全部现推。四个口径：

**口径 A｜「10 枚」＝ r1 死腿那一发的 finding 数（票面/commit 用的就是它）**
- 尺：`grep -o 'phantom-citation] comment cites repo path "[^"]*"' .scratch/wisp/probes/212/r1/post-scan-full.txt | sort -u | wc -l` → **10**。
- 那份档里 20 行＝同 10 枚各出现两次（一次在 `TestScannerSelfScanOfRealRepoIsGreen` 的 `repo HEAD violates:`，一次在 `TestRealRepoLedgerIsHonest` 的扫描 stdout），
  且**每 token 每扫只报一次**（发射点的 `citationReported` 去重）⇒ `internal/agent/approval/queue.go` 两处同 token 只有一枚进数。
- 关键：那一发的扫描器还**没有 `symRefRe`**，所以 10 枚里含着 3 枚 API 形
  （`internal/tools.Result.AppliedSteps`、`internal/buildinfo.Name`、`internal/proc.WithRegistry`），
  也还**没有 `scripts/` 前缀**，所以 AC#1 名册点名的 `internal/risk/pathresolver.go:28` 那一枚根本不在表上。
- ⇒ 5e8748b3 的 commit message 说"首发 10 枚分类"＋"10 枚真②全部改注释"：**第一个 10 是这份档的数（口径合法），第二个 10 是假的**——那 10 枚里只有 7 枚是真②。

**口径 B｜「11 token」＝ 已发行扫描器在 `5e8748b3^` 工作树上的 distinct phantom token**
- 尺（两任判语各跑一次、同数）：**11 ＝ 8 真② ＋ 3 API 形**；API 形被 `symRefRe` 排除后静默，所以"11"是普查值、不是当日红数。
- 8 真②名册：`scripts/check-pathclean-ban.sh`、`tools/gate.go`、`tools/bridge.go`（两处站点＝1 token）、
  `internal/tool`、`internal/engines`、`tools/agent/cmd`、`internal/provider`、`docs/evidence/s1/62-`。
- 本腿**没直接复现 11**（我复现的是口径 C，差集能对上）：11 与 C 的差＝`fullpack2.log` 那枚方法学伪影（见 C 的⚠）与"工作树 vs 归档树"的 untracked 面。

**口径 C｜本腿现推：归档的 pre-image 树（`git archive 5e8748b3^` 解到 OS 临时目录，仓外）**
- 尺：`cd tools/d22scan && go run . -root <归档树>`；两发（排除在／排除摘）读数在
  `logs/preimage-tree-scan-fixed.txt` 与 `logs/preimage-tree-scan-mutated-exemption-removed.txt`。
- **排除在＝8 枚 distinct token；排除摘＝9 枚；差集＝`docs/evidence/s1/62-`**。
  摘掉排除后的 9 枚＝8 真②＋1 伪影，与口径 B 的"8 真②"逐名相等 ⇒ 两任的 8 我这旁证住了。
- ⚠ **这 1 枚伪影具名（我的方法造成的，不是仓里的猎物）**：`cmd/wisp/panel_host_windows.go` 引
  `.scratch/wisp/probes/33/r6/fullpack2.log`——该件**在盘上存在但从未被 git 跟踪**
  （`git ls-files --error-unmatch` → "Did you forget to 'git add'?"），所以 `git archive` 里没有它，
  stat 落空即红。**真仓扫（工作树）同一枚是静默的**（§3 的 rc=0 clean 就是证据）。
  ⇒ 教训写进口径：**用归档树复算 phantom-citation 会把"untracked 的证据件"算成幻影**；
  后续程若要比分母，必须说清扫的是工作树还是归档树（A207 那一类机器相依分母的另一个头）。

**口径 D｜「9 hunk／8 文件」＝更正 diff 本身**
- 尺：`git show 5e8748b3 -- <212 的八枚文件> | grep -c "^@@"` → **9**；逐枚 `git show 5e8748b3 -- $f` 得
  `queue.go`=2 hunk（同一 token、两处站点），其余 7 枚各 1 ⇒ **8 文件 9 hunk**。
- 同一 commit 整枚是 **35 hunk**（`git show 5e8748b3 | grep -c "^@@"`），因为 212-r1 与 258-r1 两枚死腿的
  产码被编排者并进了同一枚代笔 commit ⇒ 数 hunk 不圈 pathspec 就会把 258 的料算进 212 的账上
  （第二任骨架的 §4 表就犯过一次，第一实例 §6 #8 具名纠正）。
- **三个数各自回答的问题不同**：A 回答"当日扫描器打印了几枚"、B/C 回答"pre-image 里有几枚幻影 token"、
  D 回答"改了几处注释"。票面第 66 行说"10 枚与树不符"，指的是拿 A 去当 D 用；真正的缺陷是
  commit message 把 A 的 10 说成"10 枚真②"，把 3 枚 API 形也报成了猎物。

## §6 判不动／量不到（本轮按住的、只上报不动手的）

1. **`cmd/wisp/models.go:20` 这一格：没按，而且不需要按**。派单指令假设"如果判语落在 `cmd/wisp/models.go` 就按住归 212-r3"。
   现量两任判语：第三任 §4 #5 判 **真**（`internal/` 下确无 `engines`），第一实例 §4 #9 也判 **真**。
   ⇒ 它不在"三处半假"里（三处＝`pending_read.go:43`、`internal/agent/tools.go:89`、`internal/tools/bridge.go:252`，两任同指，与本腿 §1 ③ 逐枚相符），
   本腿对该文件**零触碰**（`git diff --name-only d934e016..5413f46d | grep '^cmd/'` → **0 行**）。258-r2 的写面也没碰到。
2. **`internal/risk/provenance.go:63`（第 7 处更正）两任判语有分歧、且不在票面四枚清单里**：第三任判"真（带名注）"，
   第一实例判"仍存疑"（理由＝`the plugin agent command` 不可机读锚定）。票面第 66 行只点三枚 ⇒ **本腿不自作主张加一枚**，
   原样保留、具名上报，是否再改归编排者裁。（改它也不在我的写面上：`internal/risk/**` 本轮未获授权。）
3. **我这三枚注释会把下游行号路标整体推下去，逐枚量过、只上报**：
   - `internal/agent/tools.go` 净增 4 行、`internal/tools/bridge.go` 净增 6 行（`pending_read.go` 净增 3 行）。
   - 现量下游引用：`internal/tools/task_output_ac2_before_test.go:23` 引 `bridge.go:247-253` 指"Lookup 失败分支"，
     而 HEAD 上该分支在 `:266`、我改后在 `:272` ⇒ **它在 HEAD 就已经错锚 13～19 行，我的 +6 只是把误差推大，不是新造错误**。
   - `internal/session/grants.go:256` 引 `internal/tools/bridge.go:874`（`git show HEAD` 那一行是光杆 `//`）——同一形状：
     文件仍在（ban #9 不红），**行号对不对没有任何仪器在管**。⇒ 判：ban #9 只钉"文件存在性"，
     **不钉行号**；"行号路标"是本票射程外的一整族，值不值得立第五形/第十枚 ban 归编排者开票，本腿不猜。
   - 另三处引 `tools.go:82-85` 的（`internal/agent/guard_test.go:204`、`internal/agent/loop.go:746`、
     `internal/agent/testtools_test.go:25`）在 HEAD 上指向的其实是 `Tools()` 的 doc（Execute 契约句在 `:88-90`），
     同样**改前就错、且都在我写面之外**。
4. **缩写豁免的"尾部"边界量不到判定依据**：`shorthandPathStarts()` 是按**区域起点**豁免的，
   所以一枚全拼路径**后面**拖一个 `…` 或 `*`（`internal/tools/gate.go…那一族`）也会静默。
   票面 AC#1 的③类原文就写着"省略号、`…那一族`" ⇒ 本腿按票面判它属③；但这是**比 2026-10-03 更宽的方向**，
   收窄它（比如要求省略号必须落在路径中段）＝再一次射程变更 ⇒ 归人工批准，本腿不自作主张。
   该边界的正面/反面都只由 §4 那一枚 silent 样本压着（三种拼法各种一枚），尾部形没有样本。
5. **`_test.go` 豁免仍是 ban #9 射程的地基**：两任判语都指出"今天全仓 clean"里省略号形靠
   `walkGo` 跳 `_test.go` 垫着（唯一活体在 `cmd/wisp/slo_report_144_windows_test.go:851`）。
   本腿的①修法让那个活体形**在产码面也不红了**，但 `_test.go` 依然是射程外，
   所以"测试件里的谎报路标"这一族今天仍然量不到（不在四枚清单里，具名登记）。
6. **NTFS 大小写不敏感的 `os.Stat` 残洞**（承两任判语 §7）：注释引 `internal/Build/x.go` 而小写件实存时会 stat 命中 ⇒ 静默。
   本腿未重跑该探针（第一实例的 sweep 交集＝0 枚），仍是量不到的已知形状。
7. **`sh scripts/d22scan.sh` 的 CI 侧行尾/版本相依**（§3 已具名）：`selftest.go` 在 gofumpt v0.12.0 下不干净，
   但 CI 用 `go install mvdan.cc/gofumpt@latest`，两把尺的结论可能不同 ⇒ 本腿无法在这台机器上判 CI 那一发会不会红。

## §7 没做成与为什么／写面纪律自证

**没做成（逐条给因）**
1. 票面 AC 框一枚未勾（`- [ ]` 全归编排者）。自证：`git status --short .scratch/wisp/issues/` 里没有票 212 那一枚（本腿对该文件零写）。
   Progress log 也没补行——派单说"可以在 Progress log 追加一行"，但也说"碰框我不追认"；本腿的进展已在本件 §0 的 commit 分界里落清，
   不再去票面添字，避免与编排者的终裁节混写。
2. 形 ⓑ（改正则去对上旧注释）**没做**：口令「212 改 ⓑ」未生效 ⇒ `symRefRe` 与 `repoPathRe` 两个 regex 本体一字未动，
   只动了注释与"区域起点"这一层新豁免。
3. 第二枚 ring 样本没加（票面只点 silent 一枚），所以 expect-ring 仍是 20。
4. `docs/evidence/s1/212-citation-ruler.md` **没有被本腿"补做出来"**：它是 ring 样本在夹具里故意不 seed 的那枚，
   真仓里不存在也不产生任何红（§3 全仓 clean）；把它建成真件反而会改变样本语义（夹具的存在性是按被扫根判的，
   建了不响、但不建才是它想钉的东西）。⇒ 判：不属于"用改文件消红"的场合，不动。
5. 整仓 `go test ./...` 与 `./cmd/...` 没跑（派单⛔＋258-r2 在飞）。本腿跑的包级测试只有
   `tools/d22scan` 那一枚模块（三把尺之一）与三枚被改包的 `go vet`/`go build`。
6. 三处之外的一处（`provenance.go:63`）没按判语分歧去改（§6 #2）；`pending_read.go` 上文第 36 行那句
   `ticket146_liveapprovals_backing_test.go`（短引、无目录）**没动**——它不在三枚清单里，且 `internal/agent/approval/`
   本轮只授权改我核出来的那一句；顺带具名：该测试件**真身存在**（`ls` 在案），只是写法不可机读＝③类。

**写面纪律自证**
- 本腿全部改动只落在：`tools/d22scan/main.go`、`tools/d22scan/selftestsamples.go`、
  `internal/agent/tools.go`、`internal/tools/bridge.go`、`internal/agent/approval/pending_read.go`、
  以及 `.scratch/**`（本件、`logs/`、两枚突变脚本、commit message 档）。
  ⚠ **自证尺一律按枚看自己的 commit，不用 `..HEAD`**：本文写作期间共享树又落了并发腿的
  （`git log --format="%h %s" d934e016..HEAD` 现有 12 枚，其中 `cbece45e` 258-r2 写码枚动的正是
  `cmd/wisp/config_readers_255.go`），拿 `d934e016..HEAD` 当"我只改了这五枚"的尺会当场由真变假。
  合法尺＝`git show --name-only --format= 5413f46d`＝五枚产码文件（`tools/d22scan/**` 两枚＋`internal/**` 三枚，**`cmd/` 零枚**），
  `git show --name-only --format= 7ac8965a`＝只有本件；本腿的 `.scratch` 枚同理逐枚点名。
- `cmd/**` 零字节（258-r2 在飞，交集为空：它的文件名册 `cmd/wisp/*` 与本腿五枚不重名）；
  `docs/PLAN.md`、`docs/specs/**`、`internal/observe/thresholds.go`、
  任何 golden、`tools/d22scan/allowlist.txt`、`frontend/**`、`design/**`、三枚冻结测试件——全部零触碰
  （本腿三枚 commit 的文件名册里一枚都不出现；工作树里 design/ 那堆删除与 .gitignore 改动是**别人的**，本腿未 stage、未还原、未 commit）。
- commit 全部带显式 pathspec；无 `git add -A`/`.`；无 `--amend`/`reset`/`rebase`/`stash`/`checkout .`/`clean`；
  未 push；临时件只建不删（`/tmp/tmp.nRJNr83gw6/preimage` 留在系统临时区，不在仓内；仓内新增件全部入库）。
- 产码注释里引用的仓内路径**全部逐字带目录**，且都真存在：终检一把＝`cd tools/d22scan && go run . -self-test` 37/37 绿
  ＋全仓 `sh scripts/d22scan.sh` rc=0（ban #9 对产码注释面的存在性判过，无幻影）。

## §8 终态复量（本腿最后一发，交件读数；与 §3 不同处＝共享树又被并发腿推前）

三枚 commit 之后（`7ac8965a` 骨架 → `5413f46d` 修码 → `1353ce66` 证据件），本腿把三把尺＋全仓门禁
再跑一遍终值（档 `logs/final-vet.txt`、`logs/final-selftest.txt`、`logs/final-gotest.txt`、`logs/final-gate-full.txt`）：

| 尺 | 终值逐字末行 | rc |
|---|---|---|
| `go vet ./...`（tools/d22scan 模块） | 无输出 | 0 |
| `go run . -self-test` | `d22scan -self-test: clean - all 37 direction checks passed (20 expect-ring, 17 expect-silent)` | 0 |
| `go test -count=1 ./...` | `ok  	github.com/CarlosShao/wisp/tools/d22scan	17.777s` | 0 |
| `sh scripts/d22scan.sh` | `runtests.sh: OK - packages=[./...] top-level: PASS=34 FAIL=0 SKIP=0, === RUN=76, '[no tests to run]'=0` ＋ `d22scan: clean - no D22 ban violations; live scope work: bans #1-5 internal/=228, ... ban #8 cmd/=97` | 0 |

- `diff logs/pre-scope-lines.txt logs/final-scope-lines.txt` → **空**：起手（未修码）与终态（三枚之后）
  的八枚 ban scope 行＋`examined 266 production Go files` 行＋`clean` 行**逐字相等** ⇒ 票面 AC#5
  「现有 8 枚 ban 的读数逐名不变」在 212-r2 之后仍成立（本腿只加了第五形的豁免层，未动任何既有分母）。
- ⚠ 这一段区间里并发腿也往共享树落了 commit（`cbece45e` 258-r2 改的是 `cmd/wisp/**`），
  所以"266 枚产码文件""PASS=34/RUN=76"这类分母是**带此刻注的读数**，不是本腿独立造成的差；
  两处都在同一枚 `d22scan: clean` 里，本腿无需裁决，具名登记即可。
- 本腿到交件为止：**未 push**、票面 `- [ ]` 框零碰、`cmd/**` 零字节、冻结面零触碰、
  突变两台脚本跑完即还原（md5 全等在档）、临时件只建不删。
