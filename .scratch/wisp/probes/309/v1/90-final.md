# 309-v1 · 90-final（非实现者终裁腿：五格裁决表＋具名欠账）

腿＝`309-v1`（⛔ 实现者）；被裁＝`309-r1` 三笔 `1cd57ec1`→`16e86bae`→（票尾 bullet）`b2dd7fb2`（＋读数补交 `8ffa61c6`／`8bbffc50`）。
台面＝HEAD `0d2260eb`；每格判语都配一把**本腿自己跑过**的尺（⛔ 一枚⛔ 抄实现腿的数）。尺的原文与单位在各件里。

## 裁决表

| 格 | 判语 | 最硬那一把尺（可逐字重跑） | 射程 | 一句话凭据 |
|---|---|---|---|---|
| **`AC#1`** | **附条件成立** | `sh .scratch/wisp/probes/309/v1/logs/mut-rig.txt` ＋ `…/mut-rig2.txt`（四形红／两形绿，逐形 `LANDS` 自证＋`RESTORE cmp: byte-identical OK`＋`sha256sum cb86eaa7…`） | 工作树 `internal/tools/bridge.go`（瞬态）⇄ 还原后＝HEAD blob | 摘肢⇒`:64` 红、不挂⇒`:54` 红、整枚守卫 off⇒`:64` 红、肢取反⇒`:54`＋`:64` 双红（票面"换反形也红"兑现）；⛔ 恒真⛔ 恒假比较（`rules=[]` 那发就是反证）；⛔ 空集换绿；**但**摘掉另外两肢（`b.prov == nil`／`b.injected`）**两发全绿**＝仪器只覆盖三枚前提里的一枚 |
| **`AC#2`** | **附条件成立** | `sh …/logs/nail-rig.txt`（三形：字面量 rc=1／声明＋赋值 rc=0／别名 rc=0）＋第三方对账 `269−24=245` | 工作树（同文件瞬态）＋`find cmd internal` | 分母⛔ 假（钉打 245＝`find` 两把尺之差，实算）；内置正控在场；腿的"加一枚即红"**可复现且加在钉真扫的文件里**（constructions 0→1、分母⛔ 动）；**但**"任一 caller"那一半⛔ 闭合，且这枚形状洞**已有先例**（票 154 验收程五形件） |
| **`AC#3`** | **附条件成立** | `git show -U0 16e86bae -- internal/tools/bridge.go \| grep -E '^[+-]'` 逐行滤 `^//`＝**非注释行 0**；两版 blob 剥注释后 `cmp`＝**IDENTICAL（768⇄768）** | HEAD blob ⇄ `16e86bae^` blob | 改动全在注释行、零代码行；6 枚新锚今日逐枚在盘上（`defer revokeAdmission()`＝`loop.go:379`、`id := newTaskID()`＝`:329`、`return l.run(ctx, newTaskID(), input)`＝`:340`、`ruleTaint`+`ctx.taint == nil`＝`rules_gateway.go:101/102`、`scopeMarks`＝`provenance.go:694` 且**正是追加更正②的口径**）；**但**新写那把尺自指、第一发重跑就对不上（注释说 1、实测 3） |
| **`AC#4`** | **成立** | `logs/gates-rc.txt`（三行 `rc_build=0`/`rc_vet=0`/`rc_d22scan=0`）＋ `logs/fmt-head.txt`/`fmt-gofumpt-worktree.txt` | 全仓 build／单包 vet／全仓 d22scan；格式名册四档 | 三枚门禁各 rc=0 且件⛔ 0 字节；eol 三把尺（`CR=0`／`i/lf w/lf attr/text eol=lf`／`.gitattributes:4 *.go text eol=lf`）**坐实那次自抓自纠**；格式名册**归 309-r1 新增＝0**（`gofmt` 工作树 0↔HEAD blob 0；宽域工作树 4 枚全＝`w/crlf` 假枚、blob 侧 0；`gofumpt`（`$(go env GOPATH)/bin/gofumpt.exe`，本腿找到了、⛔ 免检）工作树 `internal/tools`＝0、blob＝1 枚＝票 298 `AC#5` 已承接的既有债；全仓 `git ls-files '*.go' \| xargs gofmt -l`＝26＝22 `.scratch`＋4 非，与 `A846` §(c) 逐枚对得上）；越写面＝**0 枚**（五笔 `--name-only` 逐名） |
| **`AC#5`** | **成立**（附一句⛔ 阻格的措辞更正） | `logs/ac5-*.log` 四家票面逐字＋`git grep` 家族尺 | `.scratch/wisp/issues/` ⇄ 台账 `A838`/`A842`/`A844`/`A846` | 四家按**它们自己的正文**分开：307＝折叠少报（同向不同因，已 `-done`）／40＝R4 的**显示面**（"零消费者"住在 `A838` §3、⛔ 在票 40 正文，`grep` 0 命中）／197＝**产码文本**才是拒绝形出处（票 197 正文 `grep "没有宿主给的任务 id"` ＝0 命中 ⇒ 腿那句"票 197 那支的出处"要按"文件名对票号命名＋产码注释"读）／242＝`per scope, hence per task` 的前因（那句一字未动）；`A838` §3⛔ 在本票重记 ✅；`want_n`／G1b⛔ 揽进 309 ✅（`A844` §3 归票 154 族） |

## 承重那句（必须跟着这张表一起读）
形 Ⓑ 之后**对空 id 的实际防御＝零**（三把尺坐实：代码行 `cmp` 逐字节同一／两枚 236 空 id 载体三张台面全绿／整包名册纯新增 2、删 0、SKIP 0）。
盘上有的只是**两枚跳线**：一枚定向仪器（只认"挂没挂探测器"这一件事的两向）＋一枚枚数钉（只认 `agent.ToolRequest{` 字面量拼法、`internal/agent` 之外、`cmd|internal` 射程之内）。
⛔ 一张格⛔ 读成"以后有人递空 id 会有人管"。overclaim 清单＝`60-overclaim.md` §1 第 3／6／10 行（三处措辞在**产码注释**与**测试头注释**里，改它们的写面⛔ 在我）。

## 具名欠账（`Q-309-v1-*`，归口逐枚写清）
| 号 | 是什么（带读数） | 归口 | ⛔ 阻哪格 |
|---|---|---|---|
| `Q-309-v1-prov-nil-limb-blind` | 摘 `b.prov == nil` 肢⇒`mut-e.log` rc=0；摘 `b.injected` 肢⇒`mut-f.log` rc=0。`AC#1` 标题"三枚前提"只装了一肢；`b.injected` 票面已另记，`b.prov == nil` **⛔ 另记也⛔ 有仪器** | 编排者（照 `b.injected` 那行另记，或派一发夹具） | ⛔ 阻 `AC#1` 判据（判据已过），阻"三枚前提"那句标题兑现 |
| `Q-309-v1-who-rings-overclaim` | `bridge.go:947` "**Who rings** if this limb's shape changes"＝全称，被上一行两发绿打穿；该句射程只到 `:948-949` 那两向 | 编排者（产码注释措辞） | ⛔ 阻格，属overclaim |
| `Q-309-v1-other-gates-uninstrumented` | `bridge.go:935` "Every other gate still runs"＝结构真（`:211-213`）、断言⛔（仪器两腿只问 R4） | 编排者（要证就得再来一枚夹具，⛔ 我造） | ⛔ 阻格 |
| `Q-309-v1-nail-misses-decl-and-alias-shapes` | 形 2／形 3 实测 rc=0（`nail-2-decl.log`／`nail-3-alias.log`）＝假阴性；**先例＝`.scratch/wisp/probes/154-accept/ac1-g1-blindspot-synthetic.txt` 五形**，⛔ 另立新账，引用即可 | 编排者（归口票 154 族／Q-56；"任一 caller"那句⛔ 闭合） | **阻 `AC#2` 的"任一"**，⛔ 阻"最小事"那一半 |
| `Q-309-v1-nail-excludes-whole-agent-package` | `strings.Contains(rel,"internal/agent/")` ⇒ 今天 24 枚非测试源一枚⛔ 解析；重试包装器最自然的落点就在里面；且 `loop.go:675` 是**裸名** `ToolRequest{`，matcher 本就⛔ 认 | 编排者（同上归口） | ⛔ 阻格（裁语就写了"在 `internal/agent` 之外"） |
| `Q-309-v1-nail-roots-miss-tools-dir` | 同 module 的 `tools/` 还有 13 枚非测试源在射程外（今日⛔ import `internal/agent`＝命中 0，尺＝`grep -rn "wisp/internal/agent" tools/`） | 编排者 | ⛔ 阻格 |
| `Q-309-v1-nail-control-not-exercisable` | 内置正控（`==1` 否则 Fatalf）在场、逻辑 4 行可核；**它自己的红⛔ 可能测**（要改腿写的测试件＝⛔ 我写面）⇒ 本腿改用扫描面三形代证其区分度 | 编排者签收这个替代（若嫌⛔ 够，派一发在导出树里改 matcher） | ⛔ 阻格 |
| `Q-309-v1-selfref-ruler` | `bridge.go:879-881` 那把尺把注释自己两行数进去了：自称 1、`grep` 实测 **3**（再剥行首 `//`＝1） | 编排者（产码注释） | **阻 `AC#3` 的"带尺重述"这一半**（本格因此⛔ 判"成立"） |
| `Q-309-v1-stale-count-comment` | 钉头注释 `:9-11` 复述的"60／59／1"今日同一把尺＝**69／66／3**（增量就是这两枚 309 件自己） | 编排者 | ⛔ 阻格（钉的现量 0／245 ⛔ 受影响） |
| `Q-309-v1-anchor-debt-154` | `bridge.go` 还剩 **5 处**绝对 `loop.go:N`（⛔ 属 `AC#3` 点名段，腿在 `90-final` ⑤ 具名停手）；本腿逐枚现读：`:230` 说 `783-806`＝`decideRisk`（真身 `:813`）、`:1084` 说 `:585`（真身 `rowID int64`）、`:113`/`:252` 也已移 | **⛔ 本票账**＝票 154/151 那一族（编排者归口） | ⛔ 阻 309 任何格 |

## 本腿越格自报（⛔ 瞒）
1. 为拿"改前那发"的颜色，我**建了两棵仓外导出树**（`mktemp -d`，`git archive` 只读），在里面对 `internal/tools` 跑了整包 `go test` ⇒ 只读仓内文件、⛔ 写仓内、⛔ 建 worktree。
2. 瞬态突变只打在 `internal/tools/bridge.go`（派单点名的许可面），六形各一次；每次 `trap` 还原＋`cmp`＋`sha256sum`；收尾 `git status --porcelain -- internal cmd`＝**0 行**。**⛔ 动过任何一枚测试件**。
3. `AC#2` 的 Rig 我⛔ 照抄腿的做法新建瞬态 `.go` 文件，改成往**已被扫的** `bridge.go` 末尾追加函数（三种形状）⇒ 这是**我的选择**，票面⛔ 要求；Rig 脚本存成 `.txt` 以守"⛔ 在 `probes/**` 建 `.go`"。
4. 我**⛔ 跑** `bash .scratch/wisp/probes/154/gate-clauses.sh`（`A844` 具名说整脚本会写脏 tracked 日志＋抢别腿资源），只按其排除项手工复现 G1 的尺。
5. 一次 `Write` 我把路径打歪成 `Wisp.scratch_placeholder__`（工具报成功，但盘上**查无此物**，仓内外都⛔ 留下任何文件；正件已按对的路径重落）。一次 `find /` 越界扫了全盘被超时中止（⛔ 写任何东西，只读）。两条都在这里报，⛔ 藏着。
6. ⛔ 翻框／⛔ 改 Status／⛔ push／⛔ amend／⛔ add -A；写面只有三张（本 probes ⇄ `docs/evidence/s1/309-v1-verdict.md` ⇄ 票 309 末尾 1 行）。

## 我认为编排者记录里需要跟着动的三处（⛔ 我代改）
- **`A846` §(b) 那条派单模板**："注释可以用 `⛔⇒①`，字符串/文案⛔ 能用 `≈ ✓ ≤`"——**⛔ 全**。`tools/d22scan/main.go:186` 的带里含 `U+2600–U+27BF`，`⛔` 自己＝**U+26D4** ⇒ 落进字符串一样红（`:1150-1152` 逐字"comments are blanked before matching, **strings are not**"）；反过来 `⇒`(U+21D2)／`①`(U+2460) 落在**刻意留的空隙**里（AGENTS §1.2 已写），⛔ 只在注释里安全、在哪都⛔ 被扫。模板该改成"⛔／≈／✓／≤ 之类落带字符⛔ 进字符串；⇒／①／→ 仪器根本⛔ 扫"。
- **`A842` 标题仍写"三处三种处置"**：票 309 的追加更正①已定"四处同形、两拒一自命名一静默"（本腿复跑四处逐字对上：`bridge.go:954`／`panel_assets.go` 的 "a CLI probe has no task, so it names one"／`subagent_197.go:263`／`task.go:709`）。台账只追加⛔ 改写，但**单引 `A842` 的下一位会拿到过期数** ⇒ 建议在下一笔台账里具名"以票 309 追加节为准"。
- **`AC#2` 的"任一"⛔ 是本轮才暴露的**：票 154 验收程早在合成树上量过同一族形状洞（`probes/154-accept/ac1-g1-blindspot-synthetic.txt`）。裁语把 `AC#2` 落成一枚枚数钉时⛔ 带上这份先例 ⇒ 若把 309 读成"以后手搭 `ToolRequest` 会响"，就是把票 154 已经量过的洞重新当新发现。归口由编排者定（⛔ 我派单）。
