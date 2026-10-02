# 票 252 — 同一条路径的两种拼法在允许根判定里分家：短名（8.3）侧留短、长名侧折长 ⇒ 允许根内的新建操作被判成越界（R2）并升到 L2

**立票时刻**：2026-10-02 09:3x +08，锚点 HEAD `ecbb7e86`（`dev`）
**来路**：只读腿 `224-c2`（`.scratch/wisp/probes/224/c2/census.md`，607 行 / 57,124 字节，骨架 `530ba16f` → 交件 `ecbb7e86`）判死票 224 那枚 CI 红的根因；**它推翻了我写进派单的前提**（我原以为"会话授权没接进装配"），详见台账 `A517`。
**性质**：⚠ **这一枚算产品行为**，不是纯仪器——它决定的是一次真操作要不要弹审批卡、要不要等到超时。

## 现量（编排者本机自跑，2026-10-02 09:3x，全部逐字）

1. **短/长名折法分岔的那两行**（我 `sed` 读过）：
   - `internal/risk/pathresolver.go:127-139`：路径**存在** ⇒ 走 `resolveHandle(p)`，把 `res.Canonical` 折成**长名**；路径**不存在** ⇒ 逐字注释 `// Path does not exist (or handle query unsupported): lexical fallback, classification still enforced — fail-closed.` 然后 `res.Canonical = p`——**原拼法照抄**（短名进去就短名出来）。
   - `internal/tools/paths.go:133` 起的 `InAllowlist`：要求**词法包含**（`rootsContain(p.roots, f)`）**与** `resolvedForm(canonical)` 的**再解析包含**（`if !ok || !rootsContain(p.roots, foldPath(rf))`）**两次都成立**（注释自陈"Requiring both is strictly narrower"，出处＝票 107 第二轮）。
   ⇒ 两种拼法只要有一侧被折长、另一侧留短，**这两次包含就不可能同时成立** ⇒ 判越界。
2. **降级那一跳**：`internal/risk/rules_gateway.go`（腿报 `:45`，我没逐行读，**引前自跑**）出 **R2** ⇒ 等级升到 **L2** ⇒ 走审批队列。
3. **授权那一支为什么救不回来**：我读了 `internal/tools/bridge.go:334` 起那一段——逐字 `if !sil.Silenced && sil.Level == risk.L1 {`，会话授权（`sessionGrantID`）**只挂在 L1 分支里**。⇒ 等级一旦被判成 L2，**授权根本没有被查询的机会**，不是"查询没命中"。
4. **owner 这台机器（swq 的本机）现量**——决定"这枚到底影不影响他"：
   - 仓库根：`REPO_LONG=D:\work\workspace\projects plans\Wisp` ／ `REPO_SHORT=D:\work\WORKSP~1\PROJEC~1\Wisp` ⇒ **两种拼法真的不相等**（`workspace` 9 字符＋目录名里有空格 ⇒ 存在 8.3 别名；该卷的短名生成**是开着的**）。
   - 他的临时目录：`TEMP_LONG=C:\Users\swq\AppData\Local\Temp` ／ `TEMP_SHORT=` **同串** ⇒ CI runner 那一发具体的触发物（`C:\Users\RUNNER~1\…`）**在他机器上不成立**。
   - `fsutil 8dot3name query` 我跑不了（`Error 5: Access is denied`，要管理员）⇒ 短名生成状态**只用上面两串别名存在性证明，不许引 fsutil 结论**。

## 今天已知的后果（逐枚带出处，⛔ 不许合并）

- CI 新发那枚 `TestTicket224LiveGrantStopsTheAssembledBridgeFromAsking`（`cmd/wisp`，41.72s）——`224-c2` 判：**不该记在票 224 名下**（同因在基线发就红）。
- `run_test.go:378` 的 `TestComposedGateBlocksAWriteForTwoSeconds`（两发同形 301.70s／基线 301.12s）——腿具名为**同一根因的第二枚 witness**：名字说 2 秒、实跑 301 秒＝等审批超时。⚠ 这条与我此前挂的"票 123／C18 超时账"**是同一件事的两种说法**，本票结案时要把它归口过来。
- `run_mode101_test.go:506`：`auto_approve` 档下仍撞 L2 卡，因 R2 不可静默——第三枚 witness。
- ⚠ 未定性：**"其余 21 枚共红"里有多少同因**——`224-c2` §5 只登记、未逐枚归因（它的 J4）。⛔ 不许在本票里"顺手当成已解释"。

## 要建什么

- [ ] **AC#1 判死到底哪一次包含先失配**（腿的 J1，本票的第一格必须是读数不是码）：`InAllowlist` 里先失败的是**词法那一腿**（`paths.go:133` 起第一段）还是 **`resolvedForm` 那一腿**？判据＝一发探针打印 `Roots()` 与两次 `Canonicalize` 的结果，逐字进证据件。⛔ 不许用"看起来两支都会失败"代替读数。
- [ ] **AC#2 双形比较的落点与形状**（修法候选，⛔ 不是"放宽"）：允许根判定对**同一物理路径的两种拼法**必须给出**同一个答案**。候选形（至少两形、附代价）：ⓐ 进 `InAllowlist` 之前把 asked 与 roots **两侧都规范化到同一形**（长名优先、缺失叶子按父目录折长后重拼）；ⓑ 沿用 `pathresolver.go:283-300` 里那条**已经写下的** runner 差异处理，把它从 A/B 黑名单侧**扩到 `[fs] allowed_dirs` 这一腿**（`224-c2` 报："那次修复只给了 A/B 黑名单侧"）。⛔ **绝不许**把两次包含改成一次（那是把票 107 的洞重新打开）、⛔ 绝不许把"解析不到＝未授权"那条 fail-closed 改软。
- [ ] **AC#3 本机可判的载具**（⛔ 不许挂〔仅 CI 可量〕）：`224-c2` §3 给的尺 C 形——在本机用**真短名**当临时根复刻 runner 形状。前提＝那枚目录真拿得到短名（本机仓库根已证可拿到，见现量 4）；⚠ 本仓既有规矩：**拿不到短名就"必须红、不许 skip"**（出处 `bridge_junction_windows_test.go:77-88` 那一族的定式）。
- [ ] **AC#4 三枚 witness 逐枚结线**：`TestTicket224…`／`run_test.go:378`／`run_mode101_test.go:506` 各配"改前必红、改后必绿"的同机对照读数。⚠ 结案时**必须明说** `TestComposedGateBlocksAWriteForTwoSeconds` 那 301 秒的账**从"C18 超时既有一笔"改归本票**，别留两笔。
- [ ] **AC#5 审计现场那一行别再被自己滤掉**（仪器的牙）：`224-c2` 报 `grantLinesOf`（`cmd/wisp/ticket224_assembly_test.go:206-217`，**行号待验**）把自己的 `tools: call … risk=… in_allowlist_scope=…` 审计行**过滤掉了**，所以 CI 红名册天生看不见"等级"这一维 ⇒ 判据＝载具**必须 dump 那行等级**；这样下次同因红一眼可读，不用再派一枚取证腿。
- [ ] **AC#6 owner 本机可达性那一格要有答案**：按现量 4，他机器上仓库根**有**短名别名、临时目录**没有** ⇒ 本票修完之后，必须有一发读数回答**"在他这台机器上，一次允许根内的新建操作到底撞不撞"**（撞＝产品缺陷、影响他；不撞＝CI 形状为主，但仍要修）。⛔ 不许按"他大概撞不到"结案。

## 禁区

- ⛔ **不许为了让测试变绿去放宽任何一次包含检查**、不许动 `thresholds.go`／SLO／golden／C18 审批超时常量一字节（`AGENTS §1.1`）。
- ⛔ 不许把 L2 卡的判级"顺手"改成 L1 或静默——等级判定属票 102／票 90／D4 地界，本票射程＝**拼法一致性**，不是判级策略。
- ⛔ 不许把票 224 的 AC 框动到一格（那一票的地界归它自己的验收腿）。本票只**归口**那枚 CI 红。
- `frontend/**`／`design/**` **既不读也不引**；grep/find 必写显式搜索根（`cmd internal tools docs scripts .scratch`），⛔ 不用 `.` 当根。
- Git 纪律：只 commit 不 push；commit 必带**显式 pathspec**；⛔ `git add -A`／`.`／`-a`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；临时件只建不删。
- 共享树里此刻有别的写腿：`cmd/wisp`＋`internal/config`（`198-r1`）、`internal/ball`（`33-r8b`）、`scripts/`（`251-r1`）⇒ **本票任何腿起手先 `git status --porcelain cmd/wisp internal/tools internal/risk` 并具名报回**；撞面就停手上报，⛔ 不许并行写。

## 交件要求（腿只交读数，判语归非实现者）

起手 `date -Iseconds`＋`git log -1` 取锚；**第一轮内先落骨架并 commit**；跑到第 100 轮前必须把「门禁读数」与「判不动的地方」两节写满并 commit（空着不算交件）。占位符自查（尺的字面文本不落在本票面里，含义：抓未填写的骨架残句与 TBD/FIXIT 一类标记，用字符类等价式在命令行上跑）必须 0（⛔ 别把尺的字面文本写进被扫文件）；计数尺一律 `| wc -l` 收尾。⛔ 表里任何一行不许引用一节空的凭据。写明**它推翻编排者上面哪一句**（上面所有行号、枚数与"两次包含"的因果都是我写的待验断言）。

## 归口

- 票 224 的会话授权实现（H1–H9 接线）由 `224-c2` 判为**已齐**⇒ `224-r3` 那一格从"接授权"改成**"随本票结线"**；票 224 自己的 6 枚未勾格不在本票射程。
- `ci-delta-1` §6 B 里那枚 `TestComposedGateBlocksAWriteForTwoSeconds` 的 300s 账**由本票接管**（AC#4）；结线时在那件证据件上具名指回本票，⛔ 别删原句。

## AC#1 读数（腿 `252-p1`，2026-10-02 09:4x +08，只追加本节，不碰任何勾选框）

证据件＝`.scratch/wisp/probes/252/p1/verdict.md`（212 行 / 23,702 字节）＋ `probe-run.log`（171 行 / 20,502 字节）。
载具＝`internal/tools/paths_shortname_252_probe_test.go`（314 行，`//go:build windows`，同包直调生产的 `foldPath`/`rootsContain`/`resolvedForm`/`Roots()`；短名一律 `GetShortPathNameW` 现取，拿不到即 `t.Fatalf`，0 枚 `t.Skip`）。
跑法＝`go test -count=1 -v -run 'TestTicket252P1' ./internal/tools/`，一遍：9 枚 `--- PASS`／0 枚 `--- FAIL`／0 枚 SKIP。未跑整包、未碰 `cmd/wisp`、未跑 `internal/risk` 任何测试。

1. **先失配的是第一段（词法，`paths.go:140`），不是 `resolvedForm` 那一腿。** roots 长名／asked 短名且叶子不存在 ⇒ `InAllowlist=false`，逐字三样：
   `Roots()=["d:\work\workspace\projects plans\wisp"]`、`foldPath(canonical)="d:\work\worksp~1\projec~1\wisp\q252p1-never-created.txt"`、
   `resolvedForm(canonical)="D:\work\workspace\projects plans\Wisp\q252p1-never-created.txt" ok=true` ⇒ **第一段 false、第二段 true**。
2. **反向那一发不红**：roots 配短名／asked 长名＋叶子不存在 ⇒ `InAllowlist=true`（两段都过）。原因看得见：`NewPathCanonicalizer` 把 roots 也过一次 C26（`paths.go:56/93`），根**存在**⇒必折长，`Roots()` 与配置拼法无关。
3. **正控**：长/长＋缺失叶子、roots 长／asked 短但指向**真存在**的目录、roots 短／asked 短且真存在 ⇒ 三枚全 `true`（载具可信）。⚠ 但"两侧**配置**同形（都写短名）"**不是正控**：roots 配短＋asked 短＋叶子缺失 ⇒ `false`，仍响在第一段。
4. **AC#6 本机可达性（半答，缺的那半写明在证据件 §5 甲-1）**：允许根＝真仓库根、asked＝`docs` 与 `internal\tools` 下**真工作区目录**里一个还不存在的文件名 ⇒ 长名 `true`、短名 `false`（响在第一段）。**形状上 owner 撞得到。** 但本机喂入者不存在：`USERPROFILE`/`APPDATA`/`LOCALAPPDATA`/`TEMP`/`TMP`/`HOME`/`os.TempDir()`/`os.UserHomeDir()` 实测**全部 same-as-long（无可用别名）**，进程 cwd 以长名到达；仓内**没有任何生产码调用 `GetShortPathNameW`**（`grep -rn GetShortPathName cmd internal tools scripts docs .scratch` 的命中全在 `*_test.go` 与注释里）⇒ runner 那一族的入口（C26 第 1 步的环境替换，`pathresolver.go:159-202`）在本机为空，剩下的唯一来源是"asked 字符串本身就带短名"（模型/CLI/面板递进来），超出本探针射程。

**推翻票面之处（详见证据件 §4 全表）**：
- `:12` 的「这两次包含就不可能同时成立」＝**推翻**。实测是一次成立一次不成立，且**两段对同一条物理路径给出相反答复**（第二段是宽容的那一段）。票面 现量 1 的因果句需按读数改写。
- `:10-11` 把分岔写成"短名侧留短、长名侧折长"（拼法决定）＝**表述推翻**：真正的变量是**存在性**（存在的树一定折长，含 roots；缺失的叶子一定照抄 `pathresolver.go:138`）。
- `:30` 的行号精度：`:283-300` 是注释块（`RUNNER~1` 那句在 `:291`），实现体在 `:303-375`。
- `:31` 的「本仓既有规矩：拿不到短名就必须红、不许 skip」＝**全称推翻**：出处 `bridge_junction_windows_test.go:77/84` 为真，但同一 8.3 前置条件在 `internal/risk/pathresolver_junction_windows_test.go:128`、`internal/winsec/tree_ownership_112_windows_test.go:122`、`internal/winsec/absoluteness_attribution_129_windows_test.go:311` 三处用的是 `t.Skip`。
- `:9-11` 的行号尺与 `:16-17` 的两串长短名读数＝**全部复证为真**（`:133/:140/:152/:153/:127/:138/rules_gateway.go:45/bridge.go:334`）。
- `:23` 那枚 301 秒 witness 的归因＝**本腿既不支持也不推翻**（`cmd/wisp` 是 `198-r1` 脏面，本腿一枚没跑）；⛔ 不许拿本件当任何候选修法（AC#2 ⓐ／ⓑ）的否证，本腿一枚候选形的作用面都没跑过。


## 10-02 10:0x 编排者裁定（`252-p1` 交件后；票面上我那三句被推翻，原话不改、就地打旧）

`252-p1` 的交件凭据：探针 `internal/tools/paths_shortname_252_probe_test.go` **314 行 / 12626 字节**；`.scratch/wisp/probes/252/p1/verdict.md` **212 / 23702**、`probe-run.log` 171/20502、`d22scan.log` 249/22985；提交 `f832fcbb`→`3465dcee`→`4b293c95`；终态 `git status --porcelain internal/tools internal/risk`＝**与起手名册逐枚相等**（都是空）；读数来自单发 `go test -count=1 -v -run 'TestTicket252P1' ./internal/tools/` ⇒ **9 PASS／0 FAIL／0 SKIP**。

**1. AC#1 已答（这是本票最硬那一格）**：**先失配的是第一段（词法，`paths.go:140`），不是 `resolvedForm` 那一腿**。逐字读数：roots 长／asked 短且叶子不存在 ⇒ `InAllowlist=false`，`foldPath(canonical)`＝`d:\work\worksp~1\projec~1\wisp\q252p1-never-created.txt`（短），而 `resolvedForm(...)`＝`D:\work\workspace\projects plans\Wisp\q252p1-never-created.txt`、**`ok=true`（长）** ⇒ **LEG1 false／LEG2 true**。
- ★ **推翻我票面 §12 那句"两次包含就不可能同时成立"**：实测是**一次成立一次不成立**，且**两段对同一条物理路径给出相反答复**——宽容的那段是 `resolvedForm`。⚠ 这条比原判断更糟：不是"整体判不出"，是**两段各说各话**。
- ★ **推翻我票面 §10-11 的"短名侧留短／长名侧折长"**：分岔变量是**存在性**，不是拼法（存在 ⇒ 必折长，**roots 也被折**；缺失叶子 ⇒ 照抄原拼法，`pathresolver.go:138`）。⇒ Q2（roots 配短／asked 长／缺叶子）实测 **`true`**，原因正是 `NewPathCanonicalizer`（`paths.go:56/93`）把 roots 也过 C26。
- ⇒ **AC#2 我据此当场裁形**（⛔ 落地腿不许再选边，有异议就停手上报）：**把进 `InAllowlist` 之前那一步补齐成"两侧同形"**——asked 与 roots **一律走同一条已存在的折叠**（`resolvedForm` 已证能给缺失叶子返回长名且 `ok=true`），使第一段与第二段拿到**同一个形**再各自作差。⛔ **绝不许**把两段并成一段、⛔ 不许改 `paths.go:153` 的 `!ok` 语义（那是"解析不到＝未授权"那一支，属票 102／107 地界）、⛔ 不许动判级策略（R2→L2 那一跳归 D4／票 102）。
- 另记它复认真凭据：我抄进派单的尺**全部为真**（`paths.go:133/140/152/153`、`pathresolver.go:127/138`、`rules_gateway.go:45` 体内 `:47-48`＝R2/L2、`bridge.go:334`、仓库根两串不等）。行号精度纠一处：`:283-300` 是注释块（`RUNNER~1` 在 `:291`），实现体在 `:303-375`。

**2. AC#6 半答、剩那一半不在探针射程**：**形状上 owner 撞得到**（允许根＝真仓库根，`docs\`、`internal	ools\` 下不存在的新文件名：长名 `true`、短名 **`false`**，响在第一段）。**但本机喂入者为空**：`USERPROFILE/APPDATA/LOCALAPPDATA/TEMP/TMP/HOME/os.TempDir()/os.UserHomeDir()` 实测全 `same-as-long (no usable alias)`、进程 cwd 以长名到达，且**仓内没有任何生产码调 `GetShortPathNameW`**（命中全在 `_test.go`／注释）⇒ **短名只能由外部参数递进来**。⇒ 剩下那一半我裁给后续腿：要一发 owner 真机的 `[fs] allowed_dirs` 配置＋一次真 `wisp run` 的 `in_allowlist_scope=` 审计行才能定"他今天到底撞不撞"。⛔ 谁都不许把"形状上撞得到"读成"他已经在撞"。

**3. 一处纪律纠我自己**：我票面 §31 写的"本仓既有规矩：拿不到短名就**必须红不许 skip**"是**全称、不成立**——出处 `bridge_junction_windows_test.go:77/84` 为真，但同前置条件在 `pathresolver_junction_windows_test.go:128`、`tree_ownership_112_windows_test.go:122`、`absoluteness_attribution_129_windows_test.go:311` 三处**是 `t.Skip`**。⇒ AC#3 的判据按"**这一枚探针不许 skip**"来写，⛔ 不许再拿那句全称当既有定式派单。
- 顺带纠我自己一处：票面 §47 那行把**自查尺的字面文本写进了票面**（腿复跑命中 `:47`，正是我天天写进派单的那条坑）⇒ 已就地改成描述形。

**4. 那枚 301 秒 witness**：`252-p1` 明确**既不支持也不推翻**（`cmd/wisp` 当时是脏面，它一枚没跑），⛔ 且它**没跑过任何候选修法的作用面**⇒ **本件不构成 AC#2 ⓐ／ⓑ 的否证**（第 74 条那条规矩，腿自己守住了）。
