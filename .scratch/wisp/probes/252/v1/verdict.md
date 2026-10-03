# 票 252 · 非实现者验收腿 v1（二程，`252-v1`）—— AC#1／AC#6 结清 + 五发复认 + 恒真两问 + 调用者名册

**角色**：本腿＝非实现者验收（写修法的是 `252-r2`，落地 `70a935ce`）。本腿产码一字未改（含 `*_test.go`）；
落盘面只有本目录；生产码实验全部经 `go test -overlay` 从 `.scratch` 挂载或 cp 副本→改→还原（md5 逐发复量）。

## §0 起手锚（本腿自取，全部现跑）

- `date` = `Sat Oct  3 15:55:13 CST 2026`（起手）；收尾 `Sat Oct  3 17:16:39 CST 2026`（§9）。
- `git log -1 --format=%H` = `6c7a12b656188405e7ab40d6bb7c47066ff55dd6`（分支 `dev`）。
- 起手 `git status --porcelain internal/tools internal/risk` = **0 行**（逐字空输出）。
- 起手包级名册（`go test -count=1 -v`，本腿自跑；日志 `start-tools-test.log`／`start-risk-test.log`）：
  - `./internal/tools/`：**262 PASS／0 FAIL／0 SKIP**（`ok … 18.775s`；含子测试行计数）；
  - `./internal/risk/`：**197 PASS／1 SKIP／0 FAIL**（`ok … 4.888s`；那枚 SKIP＝`TestSyncRegistryProbeLive`，
    `syncdirs_windows_test.go:133` 的既有文档化 skip，起手与收尾各一枚、逐字同形，与本腿无关）。
  - 名册文件 `roster-start-tools.txt`（262 行）／`roster-start-risk.txt`（198 行）。
- ⛔ 未跑全仓 `go test`／`go build ./...`／`go vet ./...`；`cmd/wisp` 一枚没跑（那包有别人的脏面，见 §9）。
- 生产码锚点（本腿现跑 `find cmd internal tools scripts -name paths.go` 唯一命中 `internal/tools/paths.go`；
  `grep -n` 现取；md5 现量）：
  - `internal/tools/paths.go` md5 = `7a86da7aeba8420639491cd252a47a0c` ＝ 与 `252-r2` §3.1 留档**逐字节相等**；
    `git log --oneline 70a935ce..HEAD -- internal/tools internal/risk` = 空 ⇒ 自 `70a935ce` 起该包零漂移。
  - `InAllowlist` 函数体：**起 `paths.go:189`，讫 `:219`**。词法腿（LEG 1）＝`:196` `if !rootsContain(p.roots, f) {`
    （`f := foldPath(canonical)` 在 `:190`，拒绝支 `:197-198`）；票 107 round 2 注释块 `:199-207`；
    `resolvedForm` 腿（LEG 2）＝`:208` `rf, ok := resolvedForm(canonical)` ＋ `:209`
    `if !ok || !rootsContain(p.roots, foldPath(rf)) {`（拒绝支 `:210-211`）；工作区腿 `:215-217`（不在硬禁分母）。
  - 同形修法落点：`Canonicalize` 内 `:145-148`（`if res.Resolved { return res.Canonical }; return
    sameFormOfUnresolved(res.Canonical)`）；`sameFormOfUnresolved` 体 `:169-184`（`return rf` 在 `:183`）。
- ⚠ 目录交代：`probes/252/v1/` 在本腿进场前已存在（`6c7a12b6` 以"死腿残面"入库的前一枚同名腿的遗物，
  约 60 枚文件；无 `verdict.md`）。**本腿未把它们任何一枚当凭据**；本腿自己的件＝`verdict.md`（本件）、
  `overlay/`、`overlay.json`、`overlay-inverted.json`、`overlay-inv/`、`v1-*.log/txt`、`roster-*`、
  `*pristine-v1`、`my-*-reds*.txt`／`r2-*-reds*.txt`。
- 前人件已读：票面全文（106 行）＋ `.scratch/wisp/probes/252/r2/impl.md`（284 行）＋ 其 `mut/` 目录盘点
  （18 枚文件，含 `m-inroots-with-rulers.log` 1073 行）。

## §1 AC#1 探针读数（先失配的是哪一腿）

**载具**：`.scratch/wisp/probes/252/v1/overlay/paths_252v1_probe_test.go`（`//go:build windows`，同包，
经 `overlay.json` 挂到 `internal/tools/paths_252v1_probe_test.go`——**树内零新文件**）。复用既有量尺
`measure252`／`legs252`／`shortName252`／`aliasOf252`／`repoRoot252`（`paths_shortname_252_probe_test.go`，
同包同 build tag）。短名一律 `GetShortPathNameW` 现取（`shortName252` 拿不到即 `t.Fatalf`，零 skip）。
跑法＝`go test -count=1 -v -overlay .scratch/wisp/probes/252/v1/overlay.json -run 'TestTicket252V1' ./internal/tools/`
，读数全文 `v1-probe-run.log`；结果 **2 枚测试 PASS／0 FAIL**，未创建任何文件（除 `t.TempDir()` 内）。

### §1.1 Roots() 与两次 Canonicalize 的逐字读数

起手自证（本机就是 owner 机）：

```
REPO LONG : "D:\\work\\workspace\\projects plans\\Wisp"
REPO SHORT: "D:\\work\\WORKSP~1\\PROJEC~1\\Wisp"
EQUAL? false
Roots() = ["d:\\work\\workspace\\projects plans\\wisp"]
```

**C1（修法的参照形：允许根长名、ask 用 8.3 别名、叶子不存在＝一次新建）**：

```
asked raw           : "D:\\work\\WORKSP~1\\PROJEC~1\\Wisp\\q252v1-never-created.txt"
Canonicalize(asked) : "D:\\work\\workspace\\projects plans\\Wisp\\q252v1-never-created.txt"
Roots()                       : ["d:\\work\\workspace\\projects plans\\wisp"]
foldPath(canonical)           : "d:\\work\\workspace\\projects plans\\wisp\\q252v1-never-created.txt"
LEG1 rootsContain(roots,folded): true
resolvedForm(canonical)        : "D:\\work\\workspace\\projects plans\\Wisp\\q252v1-never-created.txt" ok=true
LEG2 rootsContain(roots,foldRF): true
InAllowlist(canonical)         : true
FIRST FAILING LEG              : none: both legs hold
```

**C2（短拼法再深一层：短根＋`docs`＋缺失叶子）**：`Canonicalize` 同样折成长名，LEG1/LEG2 **双 true**、
`InAllowlist=true`、`FIRST FAILING LEG: none: both legs hold`（全文在 `v1-probe-run.log`）。

**C3（长名正控）**：`asked=长根\docs\q252v1-never-created.txt` ⇒ LEG1/LEG2 双 true、`InAllowlist=true`。

**D1（关键对照：同一条短拼法 ask **不过** `Canonicalize`、直递 `InAllowlist`——修后唯一还能让两腿分歧的输入形）**：

```
Roots()                       : ["d:\\work\\workspace\\projects plans\\wisp"]
foldPath(canonical)           : "d:\\work\\worksp~1\\projec~1\\wisp\\q252v1-never-created.txt"
LEG1 rootsContain(roots,folded): false
resolvedForm(canonical)        : "D:\\work\\workspace\\projects plans\\Wisp\\q252v1-never-created.txt" ok=true
foldPath(resolvedForm)         : "d:\\work\\workspace\\projects plans\\wisp\\q252v1-never-created.txt"
LEG2 rootsContain(roots,foldRF): true
InAllowlist(canonical)         : false
FIRST FAILING LEG              : LEG 1 lexical (paths.go:140 rootsContain(p.roots, foldPath(canonical)))
```

### §1.2 判语

**先失配的是词法那一腿（LEG 1，pristine 码 `paths.go:196`；`measure252` 的旧号自报 `:140` 是 p1 时的行号）**。
凭据＝D1 那一发读数：`LEG1=false／LEG2=true`、`InAllowlist=false`，量尺派生的 `FIRST FAILING LEG` 逐字
`LEG 1 lexical`。这不是"看起来两支都会失败"：第二腿在同一输入上明确答 true（宽容的是 `resolvedForm` 段）。
**修后生产回合（C1/C2/C3 经 `Canonicalize`）两腿同答且全 true**——同形步（`:145-148`→`:169-184`）把
短拼法折成长名后递入，两段不再各说各话。⇒ 与 `252-p1` 的判定**同向**，本腿以自己的载具与读数独立复认。

## §2 AC#6 本机可达性（撞／不撞）

**判语（按票面问句"本票修完之后，一次允许根内的新建操作到底撞不撞"）：修后不撞；修前的缺陷对这台机器形状上真实可达，但仓内生产链今天零喂入者——"撞"只可能由外部递进的短拼法触发。**

凭据四发（`v1-probe-run.log`，载具同 §1）：

1. **O1（门级主读数）**：允许根＝真仓库根，ask＝`D:\work\WORKSP~1\PROJEC~1\Wisp\q252v1-owner-reach-probe.txt`
   （8.3 拼法的新建叶），走**生产回合**（`Canonicalize`→`InAllowlist`，即 `rules_gateway.go:37→45` 的判定序）：
   ```
   Canonicalize        : "D:\\work\\workspace\\projects plans\\Wisp\\q252v1-owner-reach-probe.txt"
   LEG1=true  LEG2=true  InAllowlist=true   FIRST FAILING LEG: none: both legs hold
   ```
   ⇒ 修后这个 ask 判 **true**：不弹 R2/L2 卡。同一 ask 在修前码上会撞——凭据＝D1（同形输入不过同形步时
   `LEG1=false`→`InAllowlist=false`→`rules_gateway.go:45` 出 R2→L2）；这正是票 252 立票的缺陷本体。
2. **O2（长名对照）**：同叶长拼法 ⇒ 双 true。⇒ 修后两种拼法在 owner 机上同答（票面要的"同一个答案"）。
3. **O3（负控，不许变宽）**：根外 ask（`t.TempDir()` 下）⇒ `LEG1=false LEG2=false InAllowlist=false`。
   修法没把任何根外路径放进来。
4. **O4（OS 级等价＋本机形状，假根在 `t.TempDir()` 里，零真实仓库根写入）**：新建长名目录
   `...\q252v1-scratch-root` **在本机拿到了 8.3 别名** `C:\Users\swq\AppData\Local\Temp\TESTTI~1\001\Q252V1~1`
   （`aliasOf252` 判 `DIFFERENT spelling exists`）；经短拼法写入的文件从长拼法读回同一内容
   （size=11, body=`ticket252v1`）⇒ **两种拼法在这台机器上指同一个文件**。

**喂入者半边**（决定"结构性必撞"还是"只有外部输入才撞"）：生产码 `GetShortPathName` 调用者普查
（`grep -rn "GetShortPathName" cmd internal tools scripts --include="*.go"` 排除 `_test.go`）＝**0 枚命中**
（rc=1＝没找到＝好消息，与 `252-p1` 的普查同向）；环境串（`USERPROFILE/TEMP/…`）无可用别名（p1 现量，
本腿未重跑其八枚环境数）。⇒ owner 的管线今天不会**自己**生成短拼法 ask；O4 说明本卷 8.3 生成是开着的
（新建长名也拿别名），"形状可达"成立，但触发需要有人把短拼法字符串递进工具参数（用户手打／外部工具回显）。

## §3 五发突变复认（重跑的两发＋附带的第三发：红句逐字＋还原证明）

**选择**：派单要"挑它写得最含糊的两发"。选了 **M-samesource**（impl.md 形三，"只有 S1 红"那发的红句
从未被逐字引过）与 **M-inroots**（impl.md §3.1 自称"未单跑"的那一行——**但盘上存有 1073 行的
`m-inroots-with-rulers.log`**，见 §7 R-2）。另有 **M-merge 形**与 **M-unresolve 形**作为对照各跑一发
（见 §3.3/§3.4），以及 E4（对齐回退，§6）。

**实验纪律**：每发＝`cp paths.go.pristine-v1 → Edit 突变 → go test -count=1 ./internal/tools/ →
cp 还原 → md5 复量 7a86da7a…`。突变只动 `InAllowlist`／`sameFormOfUnresolved` 体内；r2 的两枚尺件
（tracked）始终在包内＝所有突变跑法都是"加尺"侧；**"未加尺"那一侧本腿没有跑**（那需要删 tracked 尺件
＝改跟踪文件，禁区）——该侧读数沿用 r2 留档，本腿不冒充。终态 porcelain＝0 行（§9）。

### §3.1 第二发 M-samesource（两段同源）——红句**字节级一致**

突变形状（照 r2 census 反推并现验）：`rf, ok := resolvedForm(canonical)` 上提到锁后，词法段的键换成
`foldPath(rf)`，两段各留独立拒绝支（S1 census 证：`containment paths.go:197 key=foldPath(rf)` ＋
`containment paths.go:209 key=foldPath(rf)`，两枚都 `tainted=true`）。
读数 `v1-m-samesource-baseline.log`：**4 FAIL**（S1 顶层＋W1 顶层＋B1＋B2），红句归一化比对（排序＋
随机 temp 号归一）：

```
diff my-samesource-reds.txt r2-samesource-reds.txt  ⇒  SAMESOURCE-REDS-BYTE-IDENTICAL
```

四枚红＝S1 断言 c（`paths_twocontainments_252_r2_test.go:321`，"both keys are built by the same chain of
calls ([foldPath resolvedForm]) - two calls, one check"）＋S1 断言 d（`:335`，"found clean=0 tainted=2
(keys \"foldPath(rf)\" / \"foldPath(rf)\")"）＋W1 B1/B2（`paths_twocontainments_252_r2_windows_test.go:89`）。
与 impl.md §3.1 表"M-samesource：S1 断言 c＋断言 d＋W1 两格"**逐枚对上**。还原证明：md5 复量相等＋终态
porcelain 0 行。

### §3.2 第一发 M-inroots（词法段被拆出判定）——红句一致，但**不是同一编辑形状**

impl.md 对这发的文字是"只留词法段，`resolvedForm` 调用留着但不再参与判定…未单跑"。**盘上档案
`m-inroots-with-rulers.log` 的 S1 census 显示它实际编辑的是：词法段保留（`containment paths.go:199 key=f
calls=[foldPath] tainted=false refuses=true`）而 `resolvedForm` 段从判定里消失（`S1 re-resolution names: []`）**
——即 M-unresolve 家族（`resolvedForm` 段被拆），不是表格行字面的镜像。

本腿按表格行字面重建的镜像编辑（**删词法段、保留 `resolvedForm` 段**＝M-merge 形）读数
（`v1-m-inroots-withrulers.log`）：**258 PASS／4 FAIL**，红句＝S1 断言 a
（`"asks the authorization book 1 times, want 2 (found: [paths.go:197 key=foldPath(rf)])"`）＋W1 B1/B2。
与 r2 的 **M-merge** 档案（`m-merge-with-rulers.red.txt`，3 枚红）比对：**红句除 S1 自报的真实行号外逐字一致**
——档案 `:208`／本腿 `:197`，差异＝r2 的编辑保留了 12 行注释块、本腿连注释一起删，S1 用 `fset` 现取行号
所以两句各自为真。`diff` 全文在 §3.5 的归一化比对里（唯一差异行＝行号那一处）。
**结论：M-merge 形复认成立（红集同、红文同、谁响同）；"M-inroots"这个名目下 impl.md 的表格行字面与
盘上档案不是同一次编辑（§7 R-2 具名）。**

### §3.3 附带对照 M-unresolve 形（`resolvedForm` 段被拆、词法段保留）

`v1-m-merge-baseline.log`（命名误置，实为该形）：**5 FAIL**＝`N7_inside_the_root_but_unresolvable`＋
`TestTicket107bProbeCLinkInsideAllowedRootStaysOutside`＋W2＋S1 顶层（断言 a，`found: [paths.go:196 key=f]`）
＋r1 顶层。红句归一化比对（小写＋6 位以上随机数归一）：

```
diff my-unres-reds-norm.txt r2-unres-reds-norm.txt  ⇒  UNRES-REDS-IDENTICAL-MODULO-RANDOM+CASE
```

与 impl.md §3.1"M-unresolve：加尺 257 PASS／5 FAIL，W2＋S1 断言 a＋既有 N7＋probe C"**逐枚对上**。

### §3.4 汇总

| 发 | 形状 | 本腿读数 | r2 档案 | 红句比对 |
|---|---|---|---|---|
| M-merge 形 | 删词法段 | 258 PASS／4 FAIL | 258/4（impl.md 表） | 除 S1 自报行号（197 vs 208）外逐字一致 |
| M-samesource | 两段同源 | 4 FAIL | 258/4 | **字节级一致** |
| M-unresolve 形 | 删 resolvedForm 段 | 5 FAIL | 257/5 | 归一化后一致 |
| （M-inroots 名目） | 档案＝M-unresolve 形 | 未另跑（同形） | 盘上 1073 行 log、5 FAIL | 档案红句＝W2+S1-a（`.red.txt` 2 枚） |

⇒ **复认判语：对得上的三发（M-merge／M-samesource／M-unresolve）全部对上；"M-inroots"一格 impl.md 的
行字面与档案不符（档案是真跑、且是另一形状），不构成红句造假，构成自述不精确。**

## §4 恒真两问的答案

**① 摘掉任一段之后是哪一枚具名用例红**（本腿三发各自的红名册，具名到子测试）：

- **M-merge 形（删词法段）**：`TestTicket252R2BothContainmentsReachable`（S1 断言 a）＋
  `TestTicket252R2LexicalLegRefusesWhatOnlyItRefuses` 及其 **B1/B2 两格**（W1）。既有名册（N1–N7、
  107b A/B/C、r1 A 系列）**全绿**——只有 r2 新增的两枚尺看得见，复认 r2 §3.1"谁响"列。
- **M-samesource**：`TestTicket252R2BothContainmentsReachable`（S1 断言 **c＋d**）＋W1 B1/B2。
- **M-unresolve 形（删 resolvedForm 段）**：`TestTicket252R2ResolvedLegRefusesWhatOnlyItRefuses`（W2）＋
  S1 断言 a ＋**既有两枚**：`TestTicket252R1AlignmentAddsNoAuthorization/N7_inside_the_root_but_unresolvable`
  （`AC#2 RED (widening)`）＋`TestTicket107bProbeCLinkInsideAllowedRootStaysOutside`（`AC#3 RED`）。
  ⇒ 复认 r2 R2-6 那句"`resolvedForm` 段本来就有牙（N7＋probe C）"。

**② 判据换反形会不会也全绿**——本腿自己挂 overlay 反形重跑（未突变码上，`overlay-inv/` 两份副本＋
`overlay-inverted.json`；树内零改动；三枚尺判据全部倒置：W1 `if l.final`→`if !l.final`、
W2 `if got`→`if !got`、S1 断言 a `len(book) != 2`→`== 2`）：

```
读数 v1-inverted-predicates-full.log：0 PASS／5 FAIL
--- FAIL: TestTicket252R2ResolvedLegRefusesWhatOnlyItRefuses
--- FAIL: TestTicket252R2BothContainmentsReachable
--- FAIL: TestTicket252R2LexicalLegRefusesWhatOnlyItRefuses
    --- FAIL: .../B1_short-of-existing-tree_+_missing-leaf
    --- FAIL: .../B2_short-of-existing-tree_itself
```

S1 的反形红句**自己打印出 `asks the authorization book 2 times, want 2`**（一字未改的生产码）＝自证它响的
是判据不是世界。⇒ **三枚尺没有一枚是恒真式**：反形立刻红。这与 r2 §4 的留档（同样 0 PASS／5 FAIL、
同名五枚）一致，本腿以自己的 overlay 独立复认。第一遍只倒 W1/W2 时（S1 未倒）读数＝4 FAIL／1 PASS
（`v1-inverted-predicates-pristine.log`），S1 在未倒判据下对正确码 PASS＝又一个"不是装饰"的旁证。

## §5 生产调用者名册与"行为到底变没变"的判语

**名册**（`grep -rn "InAllowlist" cmd internal tools scripts docs --include="*.go"` 排除 `_test.go`，逐枚
接线核实——Go 不为未使用方法报错，所以每枚都追了装配链）：

| 调用点 | 角色 | 接线核实（非测试） | 修后行为变没变 |
|---|---|---|---|
| `internal/risk/rules_gateway.go:45` | **R2→L2 门**（产品行为本体） | `ctx.canon` ← `assessorFor`（bridge.go:914-922 `WithCanonicalizer(b.paths)`，bridge.go:291 调）；入参先过 `:37 Canonicalize` | **变**——仅对"短拼法根内新建 ask"这一类：修前 false→R2→L2 卡（D1＋p1），修后 true（O1） |
| `internal/tools/bridge.go:1145`（`inScope`） | 审计面（`tool_call` 行的 in_scope 列） | bridge.go:1073 调用；`dec.Paths` ← `displayPaths`（bridge.go:934 先 `Canonicalize`） | 同上一类翻转；长拼法逐字节不变 |
| `internal/tools/task.go:838`（`pointerNotice`） | 任务指针提示语 | `TaskDeps{Paths: rt.paths}`（cmd/wisp/run.go:544）；`:834` 先 `Canonicalize` | 同上 |
| `internal/agent/spill.go:231` | spill 指针提示（票 174） | 经 `PointerJudge` 接口（spill.go:68）← `loop.go:245 WithPointerJudge(opt.Config.PointerJudge)`；**`cmd/wisp/run.go:1010-1017` 的 `agent.Config` 字面量不含 `PointerJudge` ⇒ 生产装配里 judge=nil**，走 fail-closed"无法核实"文案支，**`InAllowlist` 在 `wisp run` 生产路径上根本不被调用** | 不变（该支今天摸不到这个函数） |
| `internal/risk/assessor.go:141`／`internal/agent/spill.go:70` | 接口声明 | — | — |

**喂入面**：`rawPaths` 来自工具参数原文（`bridge.go:289 pathArgs(params, entry.Decl.PathParams)`；
`fs.go:301/314`、`fs_edit.go:391`、`fs_write.go:712/727/740/754` 声明 path 参数）＝模型/CLI/面板递进来的
**原字符串**。仓内生产码**零** `GetShortPathName` 调用者（§2）；环境与 cwd 全长名（p1 现量）。

**判语**：这两段包含修完之后，**行为真的变了的路径存在，但只在"外部递进短拼法"这一输入类上**——
R2 门对它的判定从"越界→L2 卡"翻成"放行"（O1 对 D1 一眼可比）。对长拼法 ask（今天生产链唯一自己
生成的拼法），`sameFormOfUnresolved` 逐字节返回原串（E4 的 A2/A4/N 系列在"对齐被回退"下发绿＝对长名
no-op 的实证；C3/O2 直读同证）⇒ **对 owner 的日常操作行为逐字节不变**。所以准确的结案口径不是
"绿了＝缺陷没了"，也不是"本票修的只是判据可见性"，而是：**gate 判定函数对一类真实输入的输出确实翻转了
（该类今天仓内零生成者）；`spill.go` 那一条生产装配根本没接（fail-closed 假装没这调用者）。**

## §6 既有钉射程检查（含 E4：对齐本身被回退时家族里谁红）

**E4 突变**＝把 `sameFormOfUnresolved` 的 `return rf`（`:183`）改成 `return canonical`（＝票 252 的修法被
整个回退）。**第一遍跑污染作废**（E3 突变没还原就叠了 E4，读数 `v1-e4-revert-baseline.log` 混入四枚
E3 家族红——已从 pristine 重做，干净读数＝`v1-e4-clean.log`，此处只引干净的）：

**11 FAIL**（r2 尺在包内）：`A1/A3`（r1 `TestTicket252R1ShortSpellingOfNewFileIsAuthorized`，POSITIVE
CONTROL RED）＋`BothSpellingsAnswerTheSame` 三格（`AC#2 RED: InAllowlist gives two answers to one physical
path: LONG=true SHORT=false`）＋`TestTicket252R1AlignmentAddsNoAuthorization`（`:252`："the same-form step
is not in effect (leg1=false leg2=true …)"）＋`TestTicket252R1CanonicalStillNamesOneTree`（`:276` 两拼法
规范化出两串）＋W1 的 **POSITIVE CONTROL**（`_windows_test.go:109`）。**不响的**：W2、S1（c/d/a 全过）、
N7、probe C、107b、A2/A4、P2。⇒ 射程判读：

- **r1 的 A1/A3＋BothSpellings＋Canonical＋Alignment:252 五枚是"对齐本身"的行为钉，有牙**——回退修法
  当场红；**S1/W2/N7/probe C 不守对齐**（它们守两段包含与 `!ok`，与对齐正交）＝两族钉互补、无互相冒充。
- **P1 的 Q1/Q2/Q3b（`wantFinal: -1`）判不了任何东西**（只读数不判），票面/impl.md 未冒充它们是钉——点名在案。
- **计数式断言（S1"恰好 2 枚"）**在本腿全部探针下没有悄悄失效：E1/E3 枚数 1→红；E2 枚数 2 但同源→
  c/d 红；E4 枚数 2 且异源→绿（正确：E4 不违"两段在"的禁）。
- **"行为型钉只写产物不写常量名"**：r1 钉全部走 `measure252` 直读布尔，不 grep 符号名；E4 实测它们对
  **行为**（对齐后的 canonical 形状）取值，不是词面尺。
- P2（`roots=LONG asked=SHORT of an EXISTING dir`，wantFinal=1）在 E4 下**绿**：存在的路径 C26 句柄查询
  直接折长（`res.Resolved=true`，`:145-147` 短路，同形步根本不被咨询）⇒ **P2 不守对齐**，它守的是句柄
  折长那一层。点名在案，防后人把它当对齐钉。

## §7 推翻清单（票面／`252-r2` 自述／派单，每一条都待验；本腿只列**读到不符**的）

| # | 对象原话 | 判 | 真身 |
|---|---|---|---|
| R-1 | 票面现量 4／`252-p1`："他的临时目录：TEMP_SHORT=同串 ⇒ …在他机器上不成立"（把"临时目录没有短名"当机器属性） | **表述推翻（结论不翻）** | O4 实测：本机 Temp 下**新建**长名目录**确实拿到 8.3 别名**（`TESTTI~1\001\Q252V1~1`）。"TEMP 路径本身无别名"为真（其分量本就 8.3 合形），但**不能**推广成"这台机器的 Temp 侧不存在短拼法"。不翻 AC#6 结论的原因：喂入者是"谁调用 GetShortPathNameW／谁递短串"，不是"别名存不存在"——O4 反而证实 8dot3 生成全卷开着，"形状可达"更硬 |
| R-2 | impl.md §3.1 M-inroots 行："未单跑…这一格是读数转移，不是本腿另跑的一发" | **自述与盘上不符（不翻红句）** | `.scratch/wisp/probes/252/r2/mut/m-inroots-with-rulers.log`（1073 行）真是一发完整跑（5 FAIL，红句档案 `.red.txt` 2 枚），且其 S1 census（`containment paths.go:199 key=f`、`S1 re-resolution names: []`）显示该发编辑＝**删 resolvedForm 段**（M-unresolve 形），与表格行字面"只留词法段、resolvedForm 留着不参与判定"读到的形状一致、但与"未单跑"矛盾——加尺列那一发显然跑了 |
| R-3 | impl.md §3.2 M-merge 红句引 `found: [paths.go:208 key=foldPath(rf)]` | **非不符，口径点名** | 该行号是 S1 对**其突变后文件**现取的真身（它保留了注释块）；本腿同形突变删了注释 ⇒ `:197`。两句各自为真；红句其余部分逐字一致（§3.5 归一化 diff）。读红名册的人须知 `found:` 里的行号是"当时那个文件的真话"，不是 pristine 行号 |
| R-4 | 派单："它自己判不动的两格（AC#1／AC#6）" | **复真** | r2 §6.4 自己点名 AC#4/5/6 没量；AC#1 它在 §3.3 只复认了方向（W1 读数），未以非实现者身份结格——本腿 §1 的读数即为那格的独立读数 |
| R-5 | 派单："五发突变（它声称拆掉任一段当场红）" | **复真＋收窄** | 四个突变形＋对照全部有档案或本腿复跑支撑；但"M-inroots"名目下档案的形状与表格行字面不同名（R-2）；且 impl.md 自己在 R2-6 已把"只有 S1"收窄成"S1＋W1"——本腿 E2 实测确实 W1 也响，收窄正确 |
| R-6 | 票面 AC#2 勾注："绝不许把两段并成一段…今天零仪器"（A562 原话）→ A567 已收窄"半侧无仪器" | **复真（收窄后）** | 本腿 E1/E2/E3 实测：词法段侧只有 r2 新尺看得见（E1 未加尺侧本腿未跑——需删 tracked 尺件，禁区；以 r2 留档 257/0 为凭），resolvedForm 段侧 N7＋probe C 本来就响（E3 实测） |
| R-7 | 派单："AC#6 测量用临时目录里的假根" | **遵守＋超出** | O4 用了假根；本腿另加了 O1（真仓库根＋缺失叶子、零写入）——比假根更贴 owner 的真实形状，且不碰禁区 |

其余复真不占条目：`paths.go` 行号（189/190/196/199-207/208-209/215/169-184/145-148）、md5、257→262 名册、
W1/W2/S1 三枚尺的正文与红句模板、`bridge.go:334` 的 `if !sil.Silenced && sil.Level == risk.L1 {`（本腿
现跑 `grep -n` 复量＝`:334` 逐字仍在，票面现量 3 为真）、`rules_gateway.go:45` 的 R2/L2 跳（`:45-50`）。

## §8 判不动／量不到的格子（具名）

1. **AC#4 三枚 witness 逐枚结线**（`TestTicket224…`／`run_test.go:378`／`run_mode101_test.go:506`）：
   全在 `cmd/wisp` 写面——按住未放（那包此刻有 `panel_host_windows_test.go` +79 行的外来脏面），本腿
   一枚没跑。301 秒那笔账的归口仍未结。
2. **AC#5（`grantLinesOf` 审计行可见性）**：同在 `cmd/wisp`，未验，归后续程。
3. **"未加尺"侧的突变读数**（M-merge 无尺 257/0 等）：需删 tracked 尺件才能量＝禁区，本腿沿用 r2 留档，
   未独立复跑，**不算本腿凭据**。
4. **r2 §5 的跨平台可编**（GOOS=linux/darwin vet rc=0）：未复跑。本腿 overlay 只挂在 windows 跑法上、
   树内零新文件，跨平台面与本腿无关；该声明仍只算 r2 自述。
5. **`spill.go:231` 在真实 `wisp run` 下的行为**：judge=nil（§5 表），故"修后 spill 文案变没变"无从在
   生产装配上量——除非先接 `Config.PointerJudge`（那是别的票的地界）。
6. **前一枚死腿的残面**（本目录约 60 枚旧件，`6c7a12b6` 入库）：存在性记录在案（§0），内容未审计、
   未引用；若编排者要清场请明示，本腿不删（只建不删）。
7. **"其余 21 枚共红"的归因**（票 224-c2 的 J4）：本腿未跑 `cmd/wisp`，无从归因，维持未定性。

## §9 门禁与收尾

| 尺 | 读数 |
|---|---|
| 收尾复跑 `go test -count=1 -v ./internal/tools/` | **262 PASS／0 FAIL／0 SKIP**（`end-tools-test.log`；`ok` rc=0） |
| 收尾复跑 `go test -count=1 -v ./internal/risk/` | **197 PASS／1 SKIP／0 FAIL**（`end-risk-test.log`；SKIP＝`TestSyncRegistryProbeLive`，起手既有、逐字同形） |
| 名册逐名 comm | tools：起手 262 → 终态 262，`comm -23` 丢 0、`comm -13` 增 0，`diff` 全等（"IDENTICAL"）；risk：198→198，丢 0 增 0 |
| 还原 | 每发突变后 `md5sum internal/tools/paths.go` = `7a86da7aeba8420639491cd252a47a0c`（逐发复量）；两枚尺件终态 md5 `28b7fafc…`／`57256d99…` ＝ 起手值；终态 `git status --porcelain internal/tools internal/risk` = **0 行** |
| 树内新文件 | **0**（生产与测试面）；全部落盘在 `.scratch/wisp/probes/252/v1/` |
| ⚠ 共享树 | `cmd/wisp/panel_host_windows_test.go` 出现 +79/-4 的外来脏面（**非本腿所致**；本腿起手 porcelain 只扫 `internal/tools internal/risk`＝0 行，未碰 `cmd/wisp`）——如实具名，交给编排者裁 |
| 本腿 commit | 骨架 `001db38f`；本件（全文）与读数件＝收尾 commit（见交件消息） |

——本件各节全部写实。判语归属：§1.2／§2／§5 的判语由本腿（非实现者）给出，供编排者翻勾 AC#1／AC#6。
