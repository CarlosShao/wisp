# 331 — `internal/observe/errors.go` 那句注释说"UI 和日志已经把这个键嵌进去了"，而**生产嵌键 0 枚**（过期文字族；本票只改注释、⛔ 动产码一字）

- Status: **待派**（料齐，凭据在下面两枚尺，本票起手仍要落地腿自己复跑一遍）。
- 来路＝非实现者裁决腿 `216-v2` §5 第 5 条（`docs/evidence/s1/216-v2-verdict.md:65`）＋编排者 10-17 自己复跑的那把尺；账 `A1010`。
- ⚠ **本票不是"错误文案体系"那一层**：那层已登记推迟＝`docs/specs/SPEC-12-roadmap-governance.md:77` 逐字 `| DEFERRED | 完整错误文案体系 | D23 | 所有用户可见错误有人话文案+建议（含 D37 全 17 类） | S6 | 部分错误暴露技术细节 |`。**推迟的是"做那层文案"，不是"注释可以替它说假话"**。本票只把那句不实注释改成说实话；⛔ 在本票里实现文案层、⛔ 新加导出名、⛔ 新增 C17 方法名。

## 现量（起手逐枚复跑；锚＝`0e944d93`）

| 事实 | 读数 | 尺 |
|---|---|---|
| 注释自称 | `internal/observe/errors.go:160-162` 三行逐字："MessageKey returns the stable user-visible copy key for the class. The copy text itself is deferred (D23 error-copy work, S6); **UI and logs embed the key** so copy can be added later without re-touching error sites." | 现读 |
| 生产嵌键 | **0 枚**——`MessageKey` 在射程 `internal cmd tools`（⛔`*_test.go`）只剩 `:160` 自家注释＋`:163` 声明两行 | `git grep -n "MessageKey" HEAD -- internal cmd tools ':(exclude)*_test.go'` |
| 消费者在哪 | 只在测试里：`internal/observe/errors_test.go:191-192`、`internal/llm/matrix_14_2_test.go:175`／`:249` | 同尺去掉 `:(exclude)` |
| 正控（尺不是恒 0） | 同尺含测试＝**4 处在场**（上面那四行）；同族尺在 `internal/risk` 出 **37** 枚 | 见 `216-v2` 台件 `M1`／`M2` |
| 那句"将来加上"的前提成立吗 | **不成立**：`MessageKey()` 返回 `"error." + string(c)`，全仓没有任何 UI/日志调用它 ⇒ 注释里"so copy can be added later without re-touching error sites"许诺的是一条**今天没有消费者、因此也无所谓"不重碰"**的路 | 现读码 |

## 判据

- [x] **AC#0（料格）**：起手把上面四把尺**逐枚重跑**（一律取具名提交树，⛔ 读工作树在飞字节），枚枚带 `rc=`；与表里读数不同 ⇒ 当场具名更正并把差异写进证据件，⛔ 照抄票面。
- [x] **AC#1（改口）**：把 `internal/observe/errors.go:160-162` 那三行注释改成**说实话**：① 键是为那层文案预备的、② **今天生产没有任何地方嵌它**（消费者只有测试，逐枚具名）、③ 那层文案仍归 `D23`／`S6` 的 `DEFERRED`（`SPEC-12 §5` 第 77 行）。⛔ 写成"将来也永远拦得住"式许诺；⛔ 顺手加"我们已经有 UI 用键"的任何同族句。
- [x] **AC#2（产码字节零改动那把尺）**：滤掉**整行注释**后，锚 blob ⇄ 工作树两侧**逐字节同串**（`sha256` 两侧同值＋`diff` 空、`rc=0`）；★**正控必做**＝只在副本的**非注释行**改一枚字面量 ⇒ 同一把尺必须回**非空、rc=1**（否则这把尺两种形状都放行＝恒真，不许拿它的"空"当证据）。逐行核＝改动枚数、其中非整行注释枚数**必为 0**。
- [x] **AC#3（行为不变）**：`go test ./internal/observe/ -count=1 -timeout 30m -v` 改前／改后各一发，台面＝`git archive` 的**仓外导出树**（⛔ 共享工作树跑"改前那发"），四数逐枚点名（`^=== RUN` 含子测试 ⇄ `^--- PASS` 只数顶层、⛔ 混用）＋ `^panic` 枚数；终态红名册与起手**作差＝0**。
- [x] **AC#4（门禁四数，只在共享工作树跑）**：`gofumpt.exe -l`（真身 `/d/work/base/gopath/bin/gofumpt.exe`，PATH 里可能 rc=127）／`sh scripts/d22scan.sh`（★它自己声明 `design/`＋`frontend/` 射程，指到缺 `design/` 的树会 `exit 2`＝GUARD 拒绝，⛔ 读成红）／`sh scripts/check-path-length-budget.sh`（要 `.git`）／`go vet ./internal/observe/`，逐枚落 `rc=`；红句含 `walk`／`build failed`／`no test files`／`setup failed`／`0xc0000135` ⇒ 判**台面故障**、具名写是哪一行，⛔ 计入判据。
- [x] **AC#5（纪律格）**：⛔ 翻任何框（本票六枚 `- [ ]` 由**非实现者**裁）、⛔ 改 `Status:`／票名／加 `-done`；写面只五处＝`internal/observe/errors.go` 的**整行注释**面＋`.scratch/wisp/probes/331/r1/**`＋`docs/evidence/s1/**` 一枚件（短名、⛔ 带 `§`）＋本票尾**一条** bullet；⛔ `git add -A`／`.`、`git add` 与 `git commit` 两处逐枚显式 pathspec；只 commit、⛔ push；⛔ `--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；临时件只建不删。

## 禁区

⛔ 动 `MessageKey()` 函数体或签名一字（本票零产码改动，AC#2 就是钉这件事）；⛔ 动 `internal/observe/thresholds.go`、golden、`tools/d22scan/allowlist.txt`、`go.mod`／`go.sum`／`.gitignore`／`ci.yml`／`.github/**`、三枚冻结件、`frontend/**`、`design/**`、`PLAN.md`、`docs/specs/**`；⛔ 为"让注释说满"去补一条生产调用（那是 `D23`／`S6` 那层 DEFERRED 的活，不归本票）；⛔ 跑整包 `./cmd/wisp`、⛔ `go test ./...`。

- **331-r1（实现者·只改注释）交件，2026-10-17**：起手锚 `72b65c61`（票面 `0e944d93` 为开票日旧锚，本腿三把尺现跑：生产嵌键 0、测试调用 4 枚在场、改前三行注释逐字已档，行号未漂）。提交四枚＝`e09187e1`（`internal/observe/errors.go` 整行注释面，AC#1 三件事实：键为文案层预备／今天生产 0 嵌、消费者只有 `errors_test.go` `TestMessageKeys` 与 `matrix_14_2_test.go`／文案层仍归 D23·S6 DEFERRED，末句明写"非此处许诺"）、`5cbec52d`（探针＋台件 `zero-change-check.sh`）、`adcbe8f3`＋`7ab51937`（证据件 `docs/evidence/s1/331-r1-evidence.md` 90 行＋占位自尺追加更正）。AC#2 尺两侧同 sha `f51f944c…` rc=0、正控非注释改一枚字面量回非空 rc=1（一发因 sed 行号锚旧空跑、已具名重做）；反向写实验（注释改口"已嵌键"）此尺**不红 rc=0**＝其射程只测整行注释外字节、不测注释真伪，答案全文与还原哈希 `ef90aa58…` 在证据件 §4。AC#3 两发仓外导出树同形 71/71/0/0、红名册作差 0；AC#4 四门共享工作树 rc=0。六枚 AC 框与 `Status:` 未动，归非实现者裁；⛔ push。

- **331-v1（非实现者裁决腿）终裁，2026-10-17**：六格＝**AC#0 成立／AC#1 成立／AC#2 成立／AC#3 成立／AC#4 成立／AC#5 成立**，判语件 `docs/evidence/s1/331-v1-verdict.md`（160 行／12147 字节），台件与逐字读数 `.scratch/wisp/probes/331/v1/**`。**本腿全部自己重放**（⛔ 背实现腿的数）：起手四把尺逐枚复跑同色（生产嵌键 **0**、含测试 4 枚在场、三行注释逐字、★第四把我代跑＝`internal/risk` 同族尺 **37** 相符，**这一把实现腿没跑过＝它 AC#0 的唯一缺口，读数相符故 AC#0 照判成立**）；AC#2 主尺滤整行注释后两侧同串 `f51f944c…` rc=0、正控非注释改一枚字面量 `990b3297…` rc=1、逐行核改动 10 枚非整行注释 **0** 枚；AC#3 我另建两棵仓外树（各 750 枚 `.go`、`diff -rq` 只差一枚文件、D 卷 8.3 短名 same=True、GOTMPDIR 先建）两发 **71/71/0/0** rc=0、名册作差 0；AC#4 四门共享工作树 rc=0（d22scan 射程含 `design/` 39 枚＝⛔ GUARD 拒绝；7 行 `walk` 命中全是它自检用例名与 PASS 行，`--- FAIL` 尺 0）；AC#5 六笔写面逐枚 `git show --name-only` 只四处、六枚框未翻、`Status:` 未动、⛔ push（六枚都⛔ 在 `origin/dev` 上）。**★必答（AC#2 到底钉住了什么）**：那把尺的准确射程一句话＝"**errors.go 里凡不是整行注释的字节一枚不许动**"；它对产码**有牙**（非注释改字面量、非注释加 `+ ""`、代码行上加**行尾**注释、加一枚**空行**——四形全红），对注释**三环全盲**（反向写：实话件⇄说谎件、两副本全字节 cmp **不等**，尺回**同一串** rc=0；整段删掉七行注释 rc=0；锚 blob 与**它自己**比也 rc=0 ⇒ **一枚什么都没改的交件照样拿走 AC#2 的主句**）。★本腿新抓到、实现腿没量过的一形：把**离 MessageKey 很远**的 `errors.go:82`（RetryNever 那行）改掉，尺仍 rc=0 ⇒ 失明是**全文件范围**的，AC#2 的绿**连"注释只改了被裁那一块"都不钉**，那件事的凭据是另一把尺（原始 `git diff` 逐行核），⛔ 是这把。**票面与尺的射程相符**：AC#2 标题自己写着「产码字节零改动那把尺」、正文第一句「滤掉整行注释后」、正控定在**非注释行**、末句「其中非整行注释枚数必为 0」——四句同维，**我指不出任何一句超它的射程**，⛔ 改判据文字（判据文字归编排者与裁者，⛔ 派给写腿）；真正的危险只在**读者**把 AC#2 的 rc=0 当成"注释说实话了"的凭据，那是 AC#1 那一格、只能由非实现者背（本件就是）。残余三枚：R1 首行仍留 "the stable **user-visible** copy key" 这一 designation（今日无表面显示它，⛔ 判它不实）；R2 "As read on 2026-10-17" 是**会腐烂的定点读数**、没有任何测试钉它，接线那票落地时必须改写这句；R3（**归票 225、⛔ 记本票、⛔ 二次开票**）＝`errors.go` 仍 **0** 枚 `DEFERRED(` 标记（生产 `DEFERRED(` 在 `72b65c61`／`e09187e1`／今天 HEAD 三处都 **26** 枚不变），但本票给 `errors.go` 添了一枚**裸词** `DEFERRED`（生产裸词 46→47）——225 那把尺要带括号形所以 26 枚不动，⚠ 将来谁改用裸词尺就会把这行注释读成一枚标记。**票面两处不实记我（编排者）自己**：① `docs/evidence/s1/216-v2-verdict.md` 的 37 那枚正控**台件号是 `F4`**，票面 `:14` 的"尺"列写成了 `M1／M2`（数对、引用错）；② 派单与我这侧交接写的"adcbe8f3 正文 97 行实为 **90** 行"——adcbe8f3 的真实枚数是 **88** 行（`wc -l`／`grep -c ''`／`awk NR` 三把尺同值，blob 末字节 0a），90 是**它自己更正之后**（7ab51937）的枚数；d4a1da8b 写 113 实为 **114** 属实、最终件 114 行／8082 字节与交接一致。实现腿的**文件内容**类声称全部复现，只有它自己提交说明里的行数三连不准（97→实 88、"实为 90"→应 88、113→实 114），⛔ 让腿背这笔、⛔ 改写历史（更正一律追加，它这一点做对了）。另两处口径差也记我：票面 `:24` 写"写面只**五处**"而它自己只列了**四**项，腿实触四项；`DEFERRED(` 那笔账本腿⛔ 动、⛔ 为本票开第二张票。
