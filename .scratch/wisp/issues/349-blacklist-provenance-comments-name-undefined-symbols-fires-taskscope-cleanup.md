# 349 — `blacklist.go`／`provenance.go` 的注释里点名**包里不存在的符号**（`fires()`／`taskScopeID`／`disposalScope`）＝票 343 那枚 `aEquiv` 的同族残留，⛔ 主人（来路＝`343-r1` 用两把改正尺普查抓到）

- Status: **已立，未派**（2026-10-19 09:5x 编排者立；来源＝票 343 落地腿 `343-r1` 必答 (ii) 具名上报（`.scratch/wisp/probes/343/r1/report*`＋`logs/03`）＋编排者自己现读坐实三枚锚（账＝台账 `A1045`））
- 类型：**注释说实话**（⛔ 任何行为改动、⛔ 任何判据改动、⛔ 动可执行字节）。射程＝`internal/risk/blacklist.go`＋`internal/risk/provenance.go` 的**注释面**（逐枚一行）
- ★**与票 343 同族但分票**：343 只管 `blacklist.go:193` 那**一行**（票面 `AC#1` 射程＝那一行注释），`343-r1` 明写"动 `:310`＝扩本票射程、⛔ 我动"⇒ 它上报、我另立本票接同族残留（⛔ 让 `343-r1` 顺手扩大自己的写面）。

## 现量（编排者 2026-10-19 09:5x 现读，锚＝当时 HEAD；引用前重跑）

| 事实 | 读数 | 我这把尺 |
|---|---|---|
| `fires()` 注释点名、⛔ 定义 | `internal/risk/blacklist.go:310` 现读逐字起头 `// AC#1), and it reads nothing but the reason: fires()/aRelates above de…`＝**注释行**；`git grep -c 'func .*fires' HEAD -- internal/risk`＝**0 命中**＝包内无此定义；真身 `fires` 作标识符在**别包** `internal/config/unwired.go:54`（struct 字段，路径 `cat-file -e` 现验存在） | `git show HEAD:<f>`＋`git grep -c` |
| 引入者＋年龄 | 来路＝`e052b11a`（2026-10-17＝票 334 落码腿自己）、就在它把 `aRelates` 拆成两把尺那一句注释里 ⇒ **只 1 天历史**（⛔ 等下一次月度腐烂） | `343-r1` 报＋编排者复跑 `git log -S` |
| `taskScopeID` 注释点名、⛔ 定义 | `internal/risk/provenance.go:54` 现读 `// assessor.WithTaintDetector(prov.Detector(taskScopeID)) yields the…`＝**注释行**；`343-r1` 现量＝全仓 0 枚 IDENT、从未作标识符定义 | `git show HEAD:<f>`＋尺① |
| `disposalScope` 注释点名、⛔ 定义 | `internal/risk/provenance.go:427` 现读 `// disposalScope.Defer(func() { _ = scope.Close() }) — the handle, not…`＝**注释行**；同上 0 定义 | `git show HEAD:<f>`＋尺① |
| ★两把改正尺 | 尺①`census.exe`（`go/ast` 驼峰 token 对撞包内名册，第二版补齐"形参/字段/局部⛔ 当包符号"盲区＝它自纠一处）⇄尺②`callshape.exe`（专挑注释里的**调用形** `foo(`/`foo()` 对撞 402 枚 func/method＝尺①结构性看不见纯小写调用形） | `343-r1` 交付台件（`.scratch/wisp/probes/343/r1/sh/`） |

## 判据（逐枚可机读）

- [ ] **AC#0 起手现跑**：上面四行逐枚重跑＋★**自加一把"同族到底还有几枚"的尺**＝对 `internal/risk/**` 注释里点名的每一个包内标识符逐枚问"它在包里有定义吗"（⛔ 抄 `343-r1` 的三枚名册＝票 334/343 那族"一枚抓到⛔ 同族只有这一枚"的教训；命中逐枚 `file:line`，0 命中带正控）。与表不符⇒当场具名更正。
- [ ] **AC#1 逐枚改准（每枚一行注释）**：改成**带条件的事实句**＝用**盘上真名**说清"这句注释今天到底指向什么"（别包的字段名要写全包路径、被删的历史名要写"这是 X 之前的名字"、⛔ 保留指向不存在符号的裸名、⛔ 为凑话新造一枚符号）。★完成判据＝交件侧逐枚 `git grep -n '<那枚名>' HEAD -- internal/risk` 只剩注释且**注释已改成不误导**＋`git diff -U0` 里 `^[+-]` 行逐枚以 `//` 开头（非注释行＝0）＋★**剥掉注释后可执行字节 sha256 两侧相同**（＝票 343 用过的结构尺，比"看起来只有空格"硬）。
- [ ] **AC#2 ⛔ 动任何可执行字节**：⛔ 新增 `fires`/`taskScopeID`/`disposalScope` 函数去"让注释变真"（为凑一句话造符号＝本仓禁手）；⛔ 动判级档位／回落支／豁免键语义／两文件任何非注释字节；⛔ 动票 336 的 `anchor_fallback_teeth_336_*`、票 339 的 `anchor_miss_reason_339_other_test.go`、票 334/343 的仪器件一字。
- [ ] **AC#3 ★"注释点名的符号不存在"要不要上通用仪器＝只出结论、⛔ 自己上**：`343-r1` 必答 (iii) 判**否**并量了误伤面（驼峰形 6 红/3 假阳＝50%、调用形 11 红/9 假阳＝82%；假阳真身＝Go builtin／关键字／`DEFERRED(...)` 标记约定／厂商名／**票 160 故意删掉的历史名 `CloseScope`（词面与假话不可分）**；根因＝"注释对名字做了什么断言无法从名字恢复"⇒任何区分规则退化成手工 allowlist）。★本票**复核这个结论**（⛔ 抄）：若仍判⛔ 上，具名写"哪一类词表能有界"（腿给的候选＝只对 `^a[A-Z]…` 那 6 枚词表作硬检查、但⛔ 抓 `fires()`、⛔ 泛化、仍属契约邻接＝人工批准）；⛔ 在本票新造门。
- [ ] **AC#4 颜色与门禁**：`./internal/risk/` 成对（基线＝起手锚仓外树 ⇄ 交件、同一台面、`-count=1 -timeout 30m`、四数逐枚点名⛔ 混用⛔ 相减）＋★**票 336 那 8 枚牙仍 8/8 全红**（改注释⛔ 动它们、但这是本票唯一"改了这一包"的活⇒复跑是便宜保险，台件现成＝`.scratch/wisp/probes/336/v1/sh/run-ac2-sweep.sh` 一字⛔ 动、env 指自己仓外重导的树）；门禁四数只在共享工作树、★跑前 `unset GOTMPDIR TMPDIR TEMP TMP`（票 311 shape B GUARD＝rc=2 射程拒绝⛔ 读成红）。⛔ `./cmd/wisp` 整包、⛔ `go test ./...`。
- [ ] **AC#5 禁区**：⛔ 动 `blacklist.go`／`provenance.go` 任何非注释字节；三枚冻结件／`internal/observe/thresholds.go`／golden／`allowlist.txt`／`go.mod`／`go.sum`／`.gitignore`／`.github/**`／`ci.yml`／D1–D47／C1–C32／R1–R9／D43 转移表；⛔ `frontend/**`／`design/**`；★⛔ 回头改票 343／336／334／339 判据文字（都 `-done` 或将 `-done`，要更正只追加）。
- [ ] **AC#6 纪律格**：⛔ 翻任何框（本票七格（AC#0…AC#6；盘上现量 7 枚框；原句"六格"＝编排者枚数笔误、见 A1047）由非实现者裁）、⛔ 改票名／`Status:`／加 `-done`；写面＝`internal/risk/blacklist.go`＋`internal/risk/provenance.go` 注释面＋`.scratch/wisp/probes/349/r1/**`（短名、⛔ 写出被禁的字面字符本身、逐枚非 0 字节、⛔ `.out`）＋本票末尾**一条**只追加 bullet；`add` 与 `commit` 两处 pathspec 逐枚手写字面路径（⛔ 从共享 index 反推＝`A1027`/`A1029`/`A1040`）；中文正文 Write 成文件再 `git commit -F`、★`-F` 路径必须 **Windows 形**（⛔ bash 形 `/d/...`＝pwsh 里 git 读不到、整笔 rc=128＝`A1042`/`A1043`）；名册尺锚定行首⛔ 用正文命中的 `--grep`（＝`A1045`）；⛔ push。

## 为什么现在不派

`internal/risk/**` 此刻无人写（`343-r1` 已交件），但★**⛔ 与票 343 的裁决腿 `343-v1` 同波动同一枚 `blacklist.go`**（`343-v1` 要复跑 `343-r1` 对 `blacklist.go` 的注释改动＋8 枚牙那一发；本票也动 `blacklist.go`＋跑同一批牙＝两腿同刻跑同包会互相洗脏四数、`TestResolvePerCallBudget` 争用假红在册）⇒ **排 `343-v1` 交件并收件之后**。

**撤销口令**：「撤 349」（⇒ 本票退回未派，两文件注释退回原状＝本票⛔ 动过码）。

## Progress log (append-only, newest last)

- [2026-10-19 09:5x +08] 编排者立票：来路＝`343-r1` 同族普查上报（它抓到 `fires()`＋整包两枚、⛔ 扩自己射程、具名交回立票权给我＝做对了）。★我现读坐实三枚锚（`blacklist:310`／`provenance:54`／`provenance:427` 均在注释行、`func fires` 全仓 0 定义、他包真身路径 `cat-file -e` 现验）＋四枚点名路径逐枚 `cat-file -e HEAD:<path>` rc=0。名 83 字符／rel 104（≤帽 165，★实测尺＝`printf '%s' | wc -c`；我起草时写"96"＝未现跑的抄数、现改真值）、★开票前跑两门现量＝长度门 rc=0 `VERDICT GREEN`（349 名未进 over-budget 输出＝在帽内）／d22scan rc=0 `clean`。★本票 `AC#1` 沿用票 343 那把"剥注释后可执行字节 sha256 相同"的结构尺（`343-r1` 交付、我复现过＝比"看着只有空格"硬）。

> ★**编排者追加（2026-10-19 11:3x，据票 343 终裁腿 `343-v1` 必答 (iii) 交回的一枚真阳；来源＝台账 `A1047`）**：本票名册**扩第三枚**——`internal/risk/syncdirs.go:453` 注释写 `Complete()==false`、真名＝`SyncDetectionComplete:472`（`343-v1` 现量、编排者逐枚 `cat-file -e`＋现读两行复验）＝同族**纯小写调用形**（只有"调用形尺"看得见、驼峰尺结构上抓不到），★在原本两枚射程文件（`blacklist.go`／`provenance.go`）**之外** ⇒ 落地腿写面加一枚 `internal/risk/syncdirs.go` 注释面。⇒ 本票真数＝**四枚**（`fires()`／`taskScopeID`／`disposalScope`／`Complete`）。★另两处本票票面缺陷同笔改正（⛔ 抹原句、就地标）：`AC#4` 曾把台件路径写成 `scripts/probes/336/v1/…`＝**该路径 HEAD 不存在**（`cat-file -e` rc=128）＝正是本仓那族"指名不存在的凭据路径"、由我犯、现改为真身 `.scratch/wisp/probes/336/v1/sh/run-ac2-sweep.sh`（rc=0 现验）；`AC#6` 曾自写"本票七格（AC#0…AC#6；盘上现量 7 枚框；原句"六格"＝编排者枚数笔误、见 A1047）"而盘上现量 7 枚框＝枚数笔误第十三次、现改七格。
