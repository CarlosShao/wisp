# 票 252 · AC#1 读数 · 探针腿 `252-p1`

> 本件只交读数，不交修法，不勾 AC。判语归非实现者（`SPEC-12 §4.3` #1/#3、D22 双角色）。
> 原始逐字输出＝`.scratch/wisp/probes/252/p1/probe-run.log`（314 行探针 / 171 行 / 20,502 字节日志）。

## 0. 起手锚（只读闸门，逐字）

- `date -Iseconds`：`2026-10-02T09:37:00+08:00`
- `git log -1`：`b5b8a09e ledger(A517)：224-c2 结档——它推翻我派单那句根因（"会话授权没接进装配"判为不成立）＋新立票 252（同一路径两种拼法在允许根判定里分家 ⇒ 允许根内新建被判越界 R2 升 L2）＋票 195 我裁完写进票面`
- 分支：`dev`
- `git status --porcelain internal/tools internal/risk`：

```
$ git status --porcelain internal/tools internal/risk
（无输出）
```

**起手名册＝空集（0 枚）**。本腿终态（§6）除本腿自己具名登记的新增件外必须还是空集；三条在飞写腿（`198-r1` cmd/wisp+internal/config、`33-r8b` internal/ball、`251-r1` scripts/）都不在这两个包里，全程未碰、未还原、未提交。

## 0.1 派单抄来的尺（本腿复核结果）

| 派单断言 | 复核方式 | 读数 |
|---|---|---|
| `internal/tools/paths.go:133` 起 `InAllowlist` | `grep -n` | `133:func (p *PathCanonicalizer) InAllowlist(canonical string) bool {` 成立 |
| 第一段在 `paths.go:140` | `grep -n` | `140:	if !rootsContain(p.roots, f) {` 成立 |
| 第二段在 `paths.go:152-153` | `grep -n` | `152:	rf, ok := resolvedForm(canonical)` / `153:	if !ok \|\| !rootsContain(p.roots, foldPath(rf)) {` 成立：两种情况确实写在同一个 `if` 里 |
| `pathresolver.go:127-139` 存在⇒折长／不存在⇒`:138` 照抄 | `Read` 全函数 | `127: if final, ok := resolveHandle(p); ok {` … `136-137` 注释逐字一致 … `138: res.Canonical = p` 成立 |
| `rules_gateway.go:45` 出 R2 | `grep -n` + 读 `:32-54` 全函数 | `45:	if !ctx.canon.InAllowlist(canonical) {` 成立；其体内 `rules: []RuleID{R2}` / `level: L2`（`:47-48`） |
| `bridge.go:334` 会话授权只挂 L1 | `grep -n` | `334:	if !sil.Silenced && sil.Level == risk.L1 {` 成立（编排者说已自证，本腿只核行号） |
| 本机长短名两串 | 探针 `GetShortPathNameW` 实测 | 逐字：long=`D:\work\workspace\projects plans\Wisp`／short=`D:\work\WORKSP~1\PROJEC~1\Wisp`，`equal? : false` |

## 1. AC#1 四问读数

载具＝`internal/tools/paths_shortname_252_probe_test.go`（同包，直接用生产码的 `foldPath`/`rootsContain`/`resolvedForm`/`Roots()`，不复制逻辑、不改一行生产码）。
每发都同时打印 **三样**：`Roots()`、`foldPath(canonical)`、`resolvedForm(canonical)` 的 `(值, ok)`，再打印 `InAllowlist` 的最终布尔；并把 `final == leg1 && leg2` 作为断言（workspace 腿 `paths.go:159` 在本探针里恒惰化，因为没有任何面板窄化过）。**9 枚 `--- PASS`、0 枚 `--- FAIL`、0 枚 SKIP。**

### 1.1 问答表（先答那一句）

| 问 | roots 拼法 | asked 拼法／叶子 | 第一段（词法） | 第二段（`resolvedForm`） | `InAllowlist` | **先失配的腿** |
|---|---|---|---|---|---|---|
| Q1 | 长名 | 短名，叶子不存在 | **false** | **true**（`ok=true`） | **false** | **第一段**（`paths.go:140`） |
| Q2 | 短名 | 长名，叶子不存在 | true | true | **true** | 无（两段都过） |
| Q3 | 长名 | 长名，叶子不存在（正控） | true | true | true | 无 |
| Q3-补 | 短名 | 短名，叶子不存在 | false | true | false | 第一段 |
| Q4 | 长名（=仓库根） | 短名，`docs`／`internal\tools` 下新建文件名 | false | true | false | 第一段 |
| Q4 | 长名（=仓库根） | 长名，同上 | true | true | true | 无 |

**AC#1 的那一句：`InAllowlist` 里先失配的是第一段（词法那一腿，`paths.go:140`），不是 `resolvedForm` 那一腿。**
短名 asked ＋ 叶子不存在这一发里，`resolvedForm` **答得上来而且答案是长名**（`ok=true`，第二段 true）⇒ 票面 现量 1 末尾那句"这两次包含就不可能同时成立"（见 §4 R-1）**不成立**：实测是**一次成立、一次不成立**，两段判定对同一条物理路径**给出互相矛盾的答复**，而生产码按顺序在第一段就 `return false`。

### 1.2 Q1 逐字（roots 长名／asked 短名＋叶子不存在）

```
CASE Q1 roots=LONG asked=SHORT+missing-leaf
  roots as configured : ["D:\\work\\workspace\\projects plans\\Wisp"]
  asked raw           : "D:\\work\\WORKSP~1\\PROJEC~1\\Wisp\\q252p1-never-created.txt"
  risk.Resolve        : err=<nil> Canonical="D:\\work\\WORKSP~1\\PROJEC~1\\Wisp\\q252p1-never-created.txt" Resolved=false Rewritten=false Rewrites=[]
  Canonicalize        : "D:\\work\\WORKSP~1\\PROJEC~1\\Wisp\\q252p1-never-created.txt"
  Roots()                       : ["d:\\work\\workspace\\projects plans\\wisp"]
  foldPath(canonical)           : "d:\\work\\worksp~1\\projec~1\\wisp\\q252p1-never-created.txt"
  LEG1 rootsContain(roots,folded): false
  resolvedForm(canonical)        : "D:\\work\\workspace\\projects plans\\Wisp\\q252p1-never-created.txt" ok=true
  foldPath(resolvedForm)         : "d:\\work\\workspace\\projects plans\\wisp\\q252p1-never-created.txt"
  LEG2 rootsContain(roots,foldRF): true
  InAllowlist(canonical)         : false
  FIRST FAILING LEG              : LEG 1 lexical (paths.go:140 rootsContain(p.roots, foldPath(canonical)))
```

三条派单没预料的读数：`risk.Resolve` 的 `Resolved=false` 且 `Canonical` 与 `asked raw` **逐字相等**（票面 现量 1 第一枚的 `res.Canonical = p` 照抄，实测成立）；`resolvedForm` **成功**返回长名（`ok=true`）；`Roots()` 是**小写折过的长名**（`foldPath` 在建根时就吃进 `paths.go:93`）。

### 1.3 Q2 逐字（roots 短名／asked 长名＋叶子不存在）

```
CASE Q2 roots=SHORT asked=LONG+missing-leaf
  roots as configured : ["D:\\work\\WORKSP~1\\PROJEC~1\\Wisp"]
  asked raw           : "D:\\work\\workspace\\projects plans\\Wisp\\q252p1-never-created.txt"
  Canonicalize        : "D:\\work\\workspace\\projects plans\\Wisp\\q252p1-never-created.txt"
  Roots()                       : ["d:\\work\\workspace\\projects plans\\wisp"]
  foldPath(canonical)           : "d:\\work\\workspace\\projects plans\\wisp\\q252p1-never-created.txt"
  LEG1 rootsContain(roots,folded): true
  resolvedForm(canonical)        : "D:\\work\\workspace\\projects plans\\Wisp\\q252p1-never-created.txt" ok=true
  LEG2 rootsContain(roots,foldRF): true
  InAllowlist(canonical)         : true
  FIRST FAILING LEG              : none: both legs hold
```

**反向那一发根本不红。** 原因看得见：`NewPathCanonicalizer`（`paths.go:56` 建根、`paths.go:93` 存根）把 roots 也过一次 C26，而根**存在** ⇒ 必被折成长名 ⇒ `Roots()` 与配置的拼法无关。票面 现量 1 把分岔写成"短名侧留短、长名侧折长"，实测的分岔变量是**存在性**，不是拼法（§4 R-2）。

### 1.4 Q3 正控逐字（三条都绿，载具可信）

```
CASE P1 roots=LONG asked=LONG+missing-leaf   -> LEG1 true  LEG2 true  InAllowlist true
CASE P2 roots=LONG asked=SHORT of an EXISTING dir
  risk.Resolve : Canonical="D:\\work\\workspace\\projects plans\\Wisp" Resolved=true
                 -> LEG1 true LEG2 true InAllowlist true
CASE P3 roots=SHORT asked=SHORT of an EXISTING dir
  Roots()=["d:\\work\\workspace\\projects plans\\wisp"]  -> InAllowlist true
```

⇒ **同一形（长/长，或 asked 指向真存在的树）全部 true**，负控的 false 不是载具坏。
但"**两侧配置同形**"不是正控：`Q3b`（roots 配短、asked 短＋叶子缺失）实测 `LEG1 false / LEG2 true / InAllowlist false`——配置里都写短名照样红，因为根被折长、缺失叶子留短。本腿按派单纪律先怀疑载具，怀疑的方式是把 `Roots()` 打出来（上表逐字）：`Roots()` 已是长名，所以"两边同形"这条前提在比较界面上不成立，不是探针坏了。

### 1.5 Q4＝AC#6 那一格：本机可达性

载具＝同一枚，asked 换成**这台开发机上真实存在的工作区目录**下的新建文件名（不创建任何文件）：

```
workspace dir="D:\\work\\workspace\\projects plans\\Wisp\\docs"                 short="D:\\work\\WORKSP~1\\PROJEC~1\\Wisp\\docs"
workspace dir="D:\\work\\workspace\\projects plans\\Wisp\\internal\\tools"      short="D:\\work\\WORKSP~1\\PROJEC~1\\Wisp\\internal\\tools"
W-long  docs          -> LEG1 true  LEG2 true  InAllowlist true
W-short docs          -> LEG1 false LEG2 true  InAllowlist false   FIRST FAILING LEG = LEG 1
W-long  internal\tools-> LEG1 true  LEG2 true  InAllowlist true
W-short internal\tools-> LEG1 false LEG2 true  InAllowlist false   FIRST FAILING LEG = LEG 1
```

**形状上，owner 本人完全撞得到**：允许根＝仓库根、新建文件在 `docs` 或 `internal\tools` 底下、asked 以短名拼法进来 ⇒ 同一条物理路径判 false，第一段响。

但"撞不撞"还差第二个条件——**有没有东西把 asked 拼成短名**。本机环境拼法普查（逐字，`GetShortPathNameW` 实测，未引 fsutil）：

```
env USERPROFILE  : value="C:\\Users\\swq"                              alias=same -> same-as-long (no usable alias)
env APPDATA      : value="C:\\Users\\swq\\AppData\\Roaming"            -> same-as-long (no usable alias)
env LOCALAPPDATA : value="C:\\Users\\swq\\AppData\\Local"              -> same-as-long (no usable alias)
env TEMP         : value="C:\\Users\\swq\\AppData\\Local\\Temp"        -> same-as-long (no usable alias)
env TMP          : value="C:\\Users\\swq\\AppData\\Local\\Temp"        -> same-as-long (no usable alias)
env CD           : (unset)
env HOME         : value="C:\\Users\\swq"                               -> same-as-long (no usable alias)
process cwd     : value="D:\\work\\workspace\\projects plans\\Wisp\\internal\\tools" alias="D:\\work\\WORKSP~1\\PROJEC~1\\Wisp\\internal\\tools" -> DIFFERENT spelling exists
os.TempDir()    : value="C:\\Users\\swq\\AppData\\Local\\Temp"         -> same-as-long (no usable alias)
os.UserHomeDir(): value="C:\\Users\\swq" err=<nil>                     -> same-as-long (no usable alias)
```

⇒ **本腿对 AC#6 的答案是半的，且缺的那半写明在 §5 甲-1**：C26 第 1 步能替换的每一个环境锚点在这台机器上**都没有可用的短名别名**（编排者量的 `TEMP_LONG==TEMP_SHORT` 得到独立复证，且不止 Temp：USERPROFILE/APPDATA/LOCALAPPDATA/HOME 全无别名），进程 cwd 也是以长名到达的。所以 runner 那一族的**唯一已知喂入物在本机不存在**；本机能把短名喂进 `InAllowlist` 的只剩"asked 字符串本身就带短名"这一条外部来源（模型/CLI/面板递进来的参数），而那条**不是本探针能判的**（§5 甲-1）。

## 2. 短名从哪来（生产链路）

本腿能定位到的**只有喂入接口，没有本机的喂入者**：

1. **C26 第 1 步的环境替换**＝runner 那一族的入口。`internal/risk/pathresolver.go:159-202` `expandAccounted`：`%VAR%`/`$VAR`/前导 `~` 用 `os.Getenv` 的值就地替换；替换完的值**原样进后续管道**。若 `USERPROFILE`/`TEMP` 本身是短名（GitHub windows-latest 的 `C:\Users\RUNNER~1`），短名就从这里进管道。仓里对该形状的文字记录：`internal/risk/pathresolver.go:291`（逐字 `// runner USERPROFILE is the 8.3 short form (C:\Users\RUNNER~1) while`，该句起于 `:290`）与 `internal/risk/blacklist.go:64`（"a short-form USERPROFILE, say"）。
   本机实测（§1.5 普查）：这六个锚点**全部 same-as-long** ⇒ 该入口在本机不产出短名。
2. **相对 asked 的绝对化**＝`internal/risk/pathresolver.go:206-214` `lexCanonical`（全仓唯一被赦免的 `filepath.Abs`+`Clean`），把相对路径锚到**进程 cwd**。实测本机 cwd 以长名到达（其短名别名存在但没被用作值）。
3. **缺失叶子的照抄**＝`internal/risk/pathresolver.go:136-138`：`resolveHandle` 失败⇒`res.Canonical = p`。**这一条不生产短名，它只是不消灭短名**——上面 1/2 递进来的拼法如果是短的，就一路带到 `InAllowlist`（实测 `Resolved=false`、`Canonical` 与 raw 逐字相等）。
4. **asked 字符串本身**（第三条以外的一切）：`cmd/wisp/run.go:501` 只喂 roots，asked 走 `internal/tools/bridge.go` 的工具参数（模型/CLI 给的 JSON 字符串）。仓内**没有任何生产码调用 `GetShortPathNameW`**（`grep -rn GetShortPathName cmd internal tools scripts docs .scratch` 的全部命中都在 `*_test.go` 与注释里）⇒ 短名不可能由本程序生成，只能由外部递进来或从环境继承。
5. **roots 侧不是喂入者**：`paths.go:56/93` 把 roots 也过一次 C26，根存在⇒折长（实测 `Q2`/`P3` 的 `Roots()` 恒为长名）。

**结论（本腿口径）**：生产链路里"谁把短名喂进来"＝**环境拼法（第 1 条）或调用方给的字符串（第 4 条）**；本机可证的第一条为空，第四条超出本探针射程（§5 甲-1）。**未定位到本机任何一个真在生产码里跑、且会把短名送进 `InAllowlist` 的具体调用点。**

## 3. 门禁读数

| 门禁 | 命令（逐字） | 读数 |
|---|---|---|
| D22 禁形扫描 | `sh scripts/d22scan.sh > .scratch/wisp/probes/252/p1/d22scan.log 2>&1` | rc=`0`；`d22scan: clean - no D22 ban violations`；正控那一段 `runtests.sh: OK - packages=[./...] top-level: PASS=34 FAIL=0 SKIP=0, === RUN=76, '[no tests to run]'=0`（该模块＝`tools/d22scan` 独立 module，**没有**执行 `internal/tools`／`internal/risk` 的任何测试） |
| 本探针在不在门禁射程内 | 同上，逐字命中行 | 在：`ban #8 internal/ examined 482 Go files, comments and _test.go included`；`scope bans #1-5 internal/ examined 226 production Go files` ⇒ **emoji 那一枚覆盖 `_test.go`**，而裸 `filepath.Clean/Abs` 那一枚只看 production 文件。本探针两者都没踩（用的全是 `filepath.Dir/Join/VolumeName` + OS 句柄调用）。 |
| 类型/静态检查 | `go vet ./internal/tools/` | 无输出（rc=0） |
| 格式 | `gofmt -l internal/tools/paths_shortname_252_probe_test.go \| wc -l` | `0` |
| 跑测试的窄过滤纪律 | `go test -count=1 -v -run 'TestTicket252P1' ./internal/tools/` | rc=0；`--- FAIL` 计数＝`0`；SKIP 计数＝`0`；`ok github.com/CarlosShao/wisp/internal/tools 0.045s`。**只跑过这一遍。** 未跑整包、未跑 `./...`、未跑 `cmd/wisp`、**未跑 `internal/risk` 的 `TestResolvePerCallBudget`**（本腿没在 risk 侧新建/运行任何探针，见下）。 |
| 写面闸门 | `git status --porcelain internal/tools internal/risk`（终态） | 与起手名册一致：除 §6 具名的本腿新增件（已提交）外无其他条目；`internal/risk` 全程 0 改动。 |

risk 侧探针：**本腿没有新增**。理由＝等级那一跳在 risk 侧只有两种说法需要验，而两者都不需要运行 risk 包：
（甲）`rules_gateway.go:32-54` 全函数已逐字读过，`:45` 的 `if !ctx.canon.InAllowlist(canonical) {` 体内就是 `rules: []RuleID{R2}` + `level: L2`（`:47-48`）——票面"出 R2 ⇒ 升 L2"这一跳的**码面形状成立**；但"从 `InAllowlist=false` 到一次真调用最终弹卡并等到超时"的**端到端读数不在 AC#1 射程**（属票面 AC#4 的三枚 witness，且 `cmd/wisp` 面此刻是 `198-r1` 的脏面），本腿不碰，登记在 §5 甲-2。
（乙）`resolvedForm` 失败语义（`paths.go:153` 那个把 `!ok` 和"不在根下"合并的 `if`）在 §1.2 已被实测拆开：这一发里 `ok=true`。是否要拆开那个 `if` 属票 102／107 地界，**不是本腿能定的**（§5 乙-1）。

## 4. 本腿推翻编排者哪一句

| # | 票面／派单的原句（逐字定位） | 实测 | 判语 |
|---|---|---|---|
| R-1 | 票面 `:12`「⇒ 两种拼法只要有一侧被折长、另一侧留短，**这两次包含就不可能同时成立** ⇒ 判越界。」 | `Q1`：第一段 false、**第二段 true（`resolvedForm` 返回长名且 `ok=true`）** | **推翻**。失配只有一处，且在词法那一腿；`resolvedForm` 那一腿在同一条输入上给出**相反**答复。票面 AC#1 想要的那一句因此有了答案，但票面 现量 1 的因果句要按读数改写。 |
| R-2 | 票面 `:10-11` 把分岔描述为"短名侧留短、长名侧折长"（拼法决定行为） | `Q2`（roots 配短／asked 长）＝**true**；`Q3b`（roots 配短／asked 短）＝**false**；`P2`（asked 短但**存在**）＝**true** | **推翻表述、保留结论**。真正的变量是**存在性**：存在的树一定折长（含 roots，`paths.go:56/93`），不存在的叶子一定照抄原拼法（`pathresolver.go:138`）。"配短的一侧留短"只在**该侧有缺失叶子**时才发生；roots 配短并不留短。 |
| R-3 | 票面 `:29` AC#1 的判据「不许用"看起来两支都会失败"代替读数」以及派单"光读码分不出 `resolvedForm` 失败与'解析出的形不在根下'" | 实测：`ok=true`、`foldPath(rf)` 在根下 ⇒ 两种情况都不是"失败" | **确认派单的必要性**（这一发确实只有跑起来才分得出），同时给出答案：这一段里 `resolvedForm` **没失败**，是词法腿先回 false。 |
| R-4 | 派单「`internal/risk/rules_gateway.go:45` 出 R2；`internal/tools/bridge.go:334`…」行号 | `grep -n`：`45:`／`334:`／`paths.go:133/140/152/153`／`pathresolver.go:127/138` **全部逐一对上** | **不推翻**（编排者抄的尺本腿复核为真）。 |
| R-5 | 票面 `:16`「`REPO_SHORT=D:\work\WORKSP~1\PROJEC~1\Wisp`（真不相等，该卷短名生成开着）」＋`:17`「TEMP 两串同串」 | 探针 `GetShortPathNameW` 实测：仓库根及其 `docs`／`internal\tools` 子目录**都有不同拼法的别名**；TEMP/APPDATA/LOCALAPPDATA/USERPROFILE/HOME/cwd 取值侧**全无可用别名** | **确认**（两枚都复证；未引 fsutil）。 |
| R-6 | 票面 `:30` AC#2 指「`pathresolver.go:283-300` 里那条**已经写下的** runner 差异处理」 | `:283-300` 是 `Comparison forms: the ANCHOR side…` 的**注释块**（含 runner `RUNNER~1` 的叙述在 `:290-294`）；实现体在 `:303-375`（`pathForms`/`certain`/`spellings`/`formsOf`） | **只改行号精度，不改结论**。 |
| R-7 | 票面 `:31` AC#3「⚠ 本仓既有规矩：拿不到短名就"必须红、不许 skip"（出处 `bridge_junction_windows_test.go:77-88` 那一族的定式）」 | 出处为真（`:77`/`:84` 两枚 `t.Fatalf("前置条件缺失…")`，本探针照此定式写）；但"本仓既有规矩"**不普遍**：`internal/risk/pathresolver_junction_windows_test.go:128`、`internal/winsec/tree_ownership_112_windows_test.go:122`、`internal/winsec/absoluteness_attribution_129_windows_test.go:311` 三处对**同一个 8.3 前置条件**用的是 `t.Skip` | **推翻"既有规矩"这个全称**。定式在 `internal/tools` 成立，在 `internal/risk`／`internal/winsec` 存在反例三枚；AC#3 若要选载具落点，这一差别有影响。 |
| R-8 | 票面 `:23` 把 `run_test.go:378` 那 301 秒归因于本票根因 | 本腿射程外（`cmd/wisp` 是 `198-r1` 的脏面，一枚测试都没跑过） | **既不支持也不推翻**＝§5 甲-2；⛔ 不许据此排除候选形。 |

## 5. 判不动的地方

### 甲（本腿判不动，别人补得上：命令＋期望读数）

- **甲-1｜"owner 到底撞不撞得到"缺的不是探针，是喂入者**（AC#6 的后半，§1.5 已把形状钉死）。
  能补的人＝**能读 owner 真机的 `[fs] allowed_dirs` 与会话审计的人**（编排者本人，或票面 AC#4 那枚 witness 的腿）。
  命令（在 owner 机器上，只读）：
  - `grep -n -A8 '\[fs\]' "$APPDATA/wisp/config.toml" | wc -l` ⇒ 期望读数：列出真配置的 roots（⛔ 别把整份 config 抄进证据件，里面有密钥）。
  - 一次真 `wisp run` 里让模型写一个允许根内的新文件，然后把审计行按 `in_allowlist_scope=` 抓出来（票面 AC#5 的 `grantLinesOf` 修法正落在这一发上）。期望读数：若 `risk=` 那行显示 L2 而路径是允许根内的新建 ⇒ **撞**；若那行显示 L1/静默 ⇒ 本机只输在"asked 被拼成短名"这一外部条件上，即**暂时撞不到**。
  - 本腿为什么判不动：asked 的字符串来源在模型/CLI 侧，探针只能证明"给它短名就红"，不能证明"真有人给它短名"。
- **甲-2｜等级那一跳的端到端读数**（`InAllowlist=false` ⇒ R2 ⇒ L2 ⇒ 301 秒等审批）。
  能补的人＝**票 252 AC#4 的腿**（它同时结三枚 witness）。
  命令：在 `cmd/wisp` 干净之后跑 `go test -count=1 -run 'TestComposedGateBlocksAWriteForTwoSeconds' ./cmd/wisp/`，需要 `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"`（缺了是 `exit status 0xc0000135` 且没有 `--- FAIL`）。期望读数：**改前**那枚的耗时仍是 300 秒级（＝等审批超时，不是 2 秒），且审计行打出 `level=L2`；改后应回到秒级。此刻 `cmd/wisp`＝`198-r1` 脏面，本腿一行没跑。
- **甲-3｜`resolvedForm` 在"路径存在但穿 junction"那一发上到底答不答**。
  能补的人＝票 107 地界的腿（同一枚 `paths.go:153` 的 `!ok` 支）。
  命令：`mklink /J` 造一条真 junction 后单测 `resolvedForm`。期望读数：票 107 的注释声称 Windows 对 junction「reports no error and no switch」，于是 `ok=true` 且形不在根下 ⇒ 第二段 true/false 的哪一侧由那条注释决定，本探针没造 junction，**没测这一支**。
  ⚠ 本腿**没有**对"某修法做不到"下任何结论：上面三条都是"没测"，不是"测了说不行"。

### 乙（本腿明写不做）

- **乙-1｜不改 `resolvedForm` 的失败语义、不拆 `paths.go:153` 那个 `if`**：属票 102／107 地界（派单与票面 AC#2 都点名不许本腿定）。本腿只交"哪一段响"的读数。
- **乙-2｜不做 AC#2 的修法选型**（ⓐ 两侧同形化／ⓑ 把 A/B 黑名单侧的 runner 差异处理扩到 `allowed_dirs`）。本腿**没有**跑过任何候选形的作用面，因此**无权就任何一枚候选形说"做不到"**；AC#2 的腿不要拿本件当否证。
- **乙-3｜不动生产码**：`paths.go`／`pathresolver.go`／`bridge.go`／`rules_gateway.go` 一字节未改；`thresholds.go`／SLO／golden／C18 超时常量一字节未动；没放宽任何断言（本探针里的期望值只设在"同形正控"与"`final == leg1 && leg2`"两枚自证上，见 §1 首段）。
- **乙-4｜不勾任何 AC、不裁票**；票面只追加了一节读数（§6）。
- **乙-5｜不在 `internal/risk` 侧新建/运行探针**：AC#1 的三样读数在同包（`internal/tools`）里就能全部拿到，跑 risk 包反而要冒撞 `TestResolvePerCallBudget` 的风险。

## 6. 本腿写面与提交

写面（越界即缺陷；全部为**新增**，无修改）：

- `internal/tools/paths_shortname_252_probe_test.go`（314 行 / 12,626 字节，`//go:build windows`）
- `.scratch/wisp/probes/252/p1/verdict.md`（本件）
- `.scratch/wisp/probes/252/p1/probe-run.log`（171 行 / 20,502 字节）
- `.scratch/wisp/probes/252/p1/d22scan.log`（249 行 / 22,985 字节）
- `.scratch/wisp/issues/252-*.md`：仅末尾追加「AC#1 读数（腿 `252-p1`）」一节，**未碰任何 `- [ ]` 勾选框**。

未碰（他人未提交脏件，不还原不提交不动）：`cmd/wisp`＋`internal/config`（`198-r1`）、`internal/ball`（`33-r8b`）、`scripts/`（`251-r1`）。
