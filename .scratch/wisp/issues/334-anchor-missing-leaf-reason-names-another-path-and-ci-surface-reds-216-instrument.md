# 334 — 锚的叶子缺位时那枚**不可豁免的拒绝**，理由说的是**别的路径**（fail-toward-exists × 缺位锚 两件事相乘）；★CI 的 test-windows 会因此红在票 216 那枚新仪器上

- Status: **已立，派 `334-r1`（写腿）**（2026-10-17 11:1x 编排者立；来源＝非实现者终裁腿 `216-v3` §4 第 1 条"具名上交"（判语件 `docs/evidence/s1/216-v3-verdict.md:53`）＋编排者自己复跑的三把尺（见"现量"最后一列）；账＝台账 `A1016`）
- 类型：**产码行为（文案不实）＋仪器缺世界**。射程＝`internal/risk/**`（产码＋测试面），⛔ `cmd/wisp/**`（`228-r3` 在写）、⛔ `internal/panel/**`（`333-r1` 在写）、⛔ `frontend/**`／`design/**`
- ⚠ 行号只当**读数**、⛔ 当判据；每条附可重跑的尺。

## 这是什么（人话）

wisp 判"这个文件敏不敏感"的时候，有两件事会**乘**在一起：

1. 它先用一次 `Lstat` 试探路径（`internal/risk/pathresolver.go:427-430` 的 `pathExists`）。那枚注释**逐字写着它是刻意的"偏向存在"**：只有干净的"没有这个文件"才算不存在，**其它任何错**（拒绝访问、共享冲突、**路径语法不正确**）都算"这枚名字的真实拼法我没法证明" ⇒ `partial=true`（`:364`）⇒ `certain()=false`（`:327`）。
2. 同时，某枚 A 档锚（例：`~/.aws/credentials`）的**叶子目录今天在这台机器上根本不存在**（`C:\Users\swq\.aws` 缺位）⇒ 它展开停在主目录 ⇒ `anchorForms` 的 `root` 变成 `c:\users\swq`。

两件事一乘：候选那串"我证明不了真实拼法"＋锚那串"它的 root 就是主目录" ⇒ `uncertainAnchorMiss`（`internal/risk/blacklist.go:314-322`）判"两者有祖先关系"⇒ 走 `blacklist.go:223-224` **报 A 档**，而那句 reason 打印的是 **`~/.aws/credentials`**——★**可那枚被拒的路径根本不是它**（红句里被打印的路径是临时目录下的 `id_ed25519`）。

⇒ 后果两样：**(a) 拒绝的等级是"不可豁免"**（TierA⇒Deny，`internal/tools/rules_gateway.go:77`），任何落在主目录祖先链内、名字里有 Windows 拒绝字符（`*`、`?`、ESC…）的合法无关文件都会被报成 A——**射程含 wisp 自己的数据根**（`%APPDATA%\wisp\logs` 盘上在场）；**(b) 理由指错了东西**，人（和模型）看到的那句话把一次"锚没装"说成了"这是 AWS 凭据"。

**方向是更严，不是更松**：这一形买到的是**拒绝**，⛔ 放行——所以它不是"有人能绕过闸门"，是"闸门有时会误伤，且说的理由不诚实"。

## 现量（起手逐枚复跑；锚＝`6e4abaf1`）

| 事实 | 读数（216-v3 交回 ⇄ 编排者复跑） | 我这把尺 |
|---|---|---|
| `pathExists` 是刻意 fail-toward-exists | 注释逐字在 `pathresolver.go:427-429`，函数在 `:430` | ★编排者 10-17 11:1x 自己 `awk 'NR>=424&&NR<=432'` 读到 427-430，**这是既有设计、⛔ 票 216 的字节引入** |
| 触发那一跳的错串 | `The filename, directory name, or volume label syntax is incorrect` ⇒ `errors.Is(err, fs.ErrNotExist)`＝false | 腿的只读探针实值；本票落地腿**自己再造一次** |
| 锚侧缺位 | `C:\Users\swq\.aws` 今天不在 ⇒ `anchorForms` 的 root 停在 `c:\users\swq`；补造 `.aws` 后 reason **换成 `~/.kube/config`** ⇒ 那句 reason 说的是"哪枚锚缺位"，⛔ 候选身份 | 腿的假主目录探针；★落地腿逐枚复跑 |
| 必要条件归因 | `%LOCALAPPDATA%` **非必要**（`C:\Users\swq\wisp216v3-home` 红 ⇄ `C:\Users\Public\wisp216v3-pub` 绿，同卷）；**主目录祖先链是必要条件** | 腿的四发台面矩阵 `.scratch/wisp/probes/216/v3/logs/**` |
| 机主今天真实会撞 | 成立：`Resolve.err=<nil>`／`Rewritten=false` ⇒ `Classify=A`、`Gate=A`；生产点 `internal/tools/mode.go:98`＋`internal/tools/paths.go:380`→`rules_gateway.go:77` | ★编排者复跑：`wc -l internal/risk/blacklist.go`＝**359** ⇒ 腿件里 `:413`/`:430` 的真身是 **pathresolver.go**（票尾 bullet 省了包名＝读数歧义，已在票 216 尾具名） |
| CI 会红 | `scripts/portable-tests.sh:152-159` 按 `%LOCALAPPDATA%` 组 `GOTMPDIR`，`./internal/risk/` 在 `:331` 的 windows 射程 ⇒ 照该形复跑＝`RUN 24／PASS_TOP 5／FAIL_TOP 1／SKIP 0` rc=1 | ★编排者自己 `sed` 核到 152-159 与 :331 两行，成立 |

## 判据（AC）

- [ ] **AC#0 起手现量**：把上面六行**逐枚重跑**（一律取具名提交树，⛔ 读工作树在飞字节），枚枚带 `rc=`；与表里读数不同 ⇒ 当场具名更正并把差异写进证据件。★再加一枚本票才有的数：**今天有多少枚真实路径会撞这一形**——只数你数得出的（探针里具名构造几枚：主目录链内含 `*`／`?`／ESC 的合法文件、`%APPDATA%\wisp\logs` 之下那一枚），⛔ 写"全部／任何"这类全称句（票 327 的 `327-a1` 已证明这种写法会被盘上读数打回）。
- [ ] **AC#1 reason 说实话（本票的产码改动就这一处）**：`uncertainAnchorMiss` 那条回落支产出的 reason，必须**指名真实候选＋真实缺位的那枚锚**，⛔ 让一句"AWS 凭据"顶替"我证明不了这枚名字的真实拼法，且主目录下某枚锚没装"。判据形状＝**定向突变**（把 reason 生成改回今天这形 ⇒ 你新增的那枚用例**必须红**，红句逐字进件）＋**成对两发落点不同**（一枚打在 reason 拼装那一行、一枚打在锚展开那一行），少任一发＝那把尺对这件事不敏感、⛔ 用它的绿当凭据。
- [ ] **AC#2 ★判级方向一字不许动（这是本票的闸门，不是可选项）**：改前⇄改后，同一枚 hostile 输入的 `class`／`needL2`／`Gate` 三值必须**逐字相同**且仍＝今天那形（A／false／A）。正控＝把 reason 换成假话以外**任何**试图把这一形降成 B 档／可豁免的改动 ⇒ 必须被你这枚正控用例咬红；`internal/tools/rules_gateway.go` 的 TierA 那一行 diff 行数**报 0**。★这一格的存在理由＝"让它变绿的最短路径就是把判级放松"，而那条**不在本票射程**（见下面"未定案那一支"）。
- [ ] **AC#3 给仪器补世界（票 216 那枚新仪器在 CI 台面上的红，只能靠这个消）**：让 `internal/risk/name_claim_two_arm_216_test.go` 在**CI 那把台面**（`GOTMPDIR` 按 `%LOCALAPPDATA%` 组、`TMP`/`TEMP` 不动，照 `scripts/portable-tests.sh:152-159` 那个形）下绿，凭据＝**造世界**——把候选放到与任何 A 锚**没有祖先关系**的位置，让"名字里的声明"这件事单独被观测。★⛔ 动任何既有断言、⛔ 删用例、⛔ `t.Skip`、⛔ 把期望值改成"红也算对"；只许**加**（加台面、加用例）。成对读数＝改前那把台面红（名册逐枚带文案）⇄ 改后绿，两发同一枚 `-run` 名册、四数逐枚点名。★主目录链内那一发的"过严形状"另存为**具名读数**（`t.Log` 或另枚只读探针件，⛔ 留成一枚恒红的断言让后续程继承）。
- [ ] **AC#4 reason 字面串的消费者普查**：`git grep` 逐枚点名今天有没有**别的**地方按那句 reason 的字面串做判定（golden、面板文案、`internal/perm`、`internal/tools`、任何 `wantReason` 常量）。命中 ⇒ 具名上报并停手，⛔ 顺手改它们（那是别的票的写面）；0 命中 ⇒ 先跑一枚已知存在的正控再报，⛔ 让"0"没有分母。
- [ ] **AC#5 门禁四数＋颜色**：`gofumpt -l`（真身 `/d/work/base/gopath/bin/gofumpt.exe`）／`sh scripts/d22scan.sh`／`sh scripts/check-path-length-budget.sh`／`go vet ./internal/risk/`，逐枚落 `rc=`，只在**共享工作树**跑；★门禁红句若指向别腿在飞件（`internal/panel/**`＝`333-r1`、`cmd/wisp/**`＋`internal/proc/**`＝`228-r3`）⇒ 具名归位、⛔ 计入本腿判据。颜色两包（`./internal/risk/`）改前⇄改后**同一台面**、★**两发都在这腿自己建的仓外导出台树里跑**（⛔ 共享工作树跑任何一发——`internal/panel` 正在被别人改写，而 panel 侧的读法会吃到你这枚 risk 字节），显式 `-count=1 -timeout 30m`、四数逐枚点名（⛔ 混 `^=== RUN` 与 `^--- PASS`）、终态红名册与起手作差逐名归位。⛔ 跑 `./cmd/wisp` 整包、⛔ `go test ./...`、⛔ 杀 `wisp.exe` PID **21892**。
- [ ] **AC#6 纪律格**：⛔ 翻任何框（本票六格由**非实现者**裁）、⛔ 改票名／`Status:`／加 `-done`、⛔ 动 `thresholds.go`／golden／`allowlist.txt`／`go.mod`／`go.sum`／`.gitignore`／`ci.yml`／三枚冻结件／D1–D47／C1–C32／R1–R9／D43 转移表。写面只三处＝`internal/risk/**` ＋ `.scratch/wisp/probes/334/r1/**` ＋本票末尾**一条**只追加 bullet；★证据件落 `.scratch/wisp/probes/334/r1/`（⛔ `docs/evidence/s1/`，那是裁决者的面）。

## 未定案那一支（不归本票，具名留给编排者⇄机主）

**要不要顺手把"判级"也放宽**（＝让这一形不再报 A，而是回到 B 档可豁免）。⛔ 本票不许做（AC#2 就是钉这件事）。理由三条，逐条可核：① 方向是更严、今天没有放行风险；② 放宽＝改判级链语义，按 `AGENTS.md §1.1`/`SPEC-12 §4.1` 属"改契约"级、要人工批准并落 `A##`；③ 误伤的真实分母今天**没人数过**（`AC#0` 那格只数得出构造出来的几枚）。默认动作＝**不动**，等机主那一句；不答的代价＝wisp 仍会偶尔用一句不诚实的理由拒绝一个其实合法的单文件操作。
- [2026-10-17 12:0x +08] 落地腿 **`334-r1`** 交件（⛔ 翻框：七格未勾数我自己现跑＝`7`，本票 AC#0…AC#6 一律等**非实现者**裁）。起手锚现认＝**`03f7e7a6`**（11:16:51；派单写的 `6e4abaf1` 与它相差 3 分钟且 `internal/risk/**` 字节全等，读数同源）。**裁甲-二 逐字执行＝只做 ①AC#1 reason 说实话 ②AC#3 给仪器补世界；判级那一支一字未动**（⛔ 降档、⛔ 任何"为了让用例绿而做的等价事"）。产码面一枚：`internal/risk/blacklist.go` **+53/−8**（九枚 A 档 case 的**条件表达式逐字未动**，只把 return 的 reason 走唯一拼装点 `aRuleReason`，并把 `aRelates` 拆成 `anchorMatches || uncertainAnchorMiss`＝同一短路顺序）；新仪器一枚 `internal/risk/anchor_miss_reason_334_test.go`（+209，世界自造：假主目录＋逐枚点名装上的锚叶子 ⇒ 链内那形在两棵树上都可判，⛔ 继承本机装了什么）；`name_claim_two_arm_216_test.go` **+20/−1**（那 1 行＝`dir := t.TempDir()` 换成 `candidateOffAnchorChain334(t)`，四对断言与 `want` 值一字未动、⛔ Skip／⛔ 删用例／⛔ 期望改"红也算对"）。提交：`e052b11a`／`4711955f`／`3a118226`（产码三枚，逐枚 name-only 只碰 `internal/risk/**`）＋`1b488dd5`／`5564b2b5`／`f06a288d`（证据）。决定性读数（全部仓外导出台树，⛔ 共享工作树跑 test）：**AC#2 三值尺＝同一枚只读探针（哈希 `e852ce63`）在两棵树上各跑一发，14 行 `class／allow／needL2` 逐字相等且仍＝今天那形（A／false／false），reason 尾巴三枚各换成 `[fail-closed, NOT a match on this rule]…`**；突变两发落点不同（M-A 打在 reason 拼装那处分支、M-B 打在锚展开那一行）⇒ 都只把 AC#1 两枚子用例咬红、AC#2 闸门格**仍绿**；**M-C 正控＝`uncertainAnchorMiss` 首行 `return false`（把这一形降回可豁免的 B）当场把闸门格咬红三句**（`class: got B, want A`／`Gate: got class=B … needL2=true`／`an A deny … was unlocked by an override key`），逐发还原哈希 `c1799cd2…`＝交件 blob；**AC#3 成对（CI 形台面＝GOTMPDIR 按 `%LOCALAPPDATA%` 组、TMP/TEMP 不动、票 216 那六枚仪器逐名并列）改前 `RUN 24／PASS_TOP 5／FAIL_TOP 1／SKIP 0` rc=1（唯一红 `TestFileBaseNameClaimChangesNoClassOrReason`，3 枚子用例）⇄ 改后 `24／6／0／0` rc=0**；整包 `./internal/risk/` 成对 `225／136／2／1` ⇄ `230／139／1／1`，残红一枚 `TestC26RewriteAccountIsConsumedAtEverySecurityLeg` 两发红句**逐字相同**且四行全指 `internal/panel/composer.go:127`、`internal/panel/display_scrub_216.go:26/37/218`＝**`333-r1` 在飞面，具名归位、⛔ 计入本腿判据**；SKIP 那枚 `TestSyncRegistryProbeLive` 两发同一枚＝本机注册表现实。主目录链内的"过严形状"按票面要求**另存为具名读数**（`SHAPE-READING334 inside-profile=class=A … || off-profile=class=B …` 一行 t.Logf ＋探针必要条件三发：链内 A ⇄ 链外 B ⇄ `%LOCALAPPDATA%` 挪走仍 A），⛔ 留恒红断言给后续程。AC#4 普查先跑正控（生产者自己命中 1 枚）再报**消费者 0 枚**：`A-list sensitive path` 全仓只在 `blacklist.go:83`；各锚字面串在测试里只当**子用例名／注释**，`docs/PLAN.md`·`SPEC-06` 是规格正文、`design/doubao/demo/screens/security.js:143` 是机主外部 mock ⇒ 无一处按字面串做判定，无需停手、⛔ 改别人写面。**更正一枚行号过期读数**（⛔ 判据，只当读数）：票面与派单写的 `internal/tools/rules_gateway.go:77` **该路径在仓里不存在**，真身＝`internal/risk/rules_gateway.go:78`（`case TierA:`）；`paths.go` 的 TierA 映射今天在 **:383**（票面 :380）。另具名上报一处方法差异：票面"补造 `.aws` ⇒ reason 换 `~/.kube/config`"我**没在机主 profile 里建目录**，改由用例自造世界复现同一读数（子用例 ⓑ）。门禁四数只在共享工作树：`gofumpt -l` rc=0／`sh scripts/d22scan.sh` rc=0（clean）／`sh scripts/check-path-length-budget.sh` rc=0（GREEN，分母 14818）／`go vet ./internal/risk/` rc=0。⛔ push。证据＝`.scratch/wisp/probes/334/r1/{ac0-ac4-readings.md,ac1-ac2-ac3-ac5-readings.md,selfcheck-answers.md}` ＋ `logs/`（含 `zero-byte-proofs.md` 给九枚零字节件名分）。**留给编排者**：`未定案那一支`（要不要把判级降回可豁免的 B）本腿一字未动、也没替它铺路——误伤的真实分母今天仍只数得出我具名构造的那几枚（链内 5/5 全 A），要裁的是"这形该不该继续报 A"，属改契约、要人工批准＋落 `A##`。
