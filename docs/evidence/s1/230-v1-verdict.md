# 票 230 AC#2 终裁判语（230-v1＝裁决者·≠实现者·只读）

## 0. 本裁射程与硬约束
- 只裁 `AC#2`。`AC#1` 归在飞的 `115-v1` 出表趟，本腿⛔ 代它裁；`AC#3`／`AC#4`／`AC#5` ⛔ 碰。
- 本腿零产码零测试字节：⛔ 动任何 `.go`、⛔ `go test`／`go build`（三枚写腿 `352-r1`／`109-r1`／`321-r1` 正在
  winsec/models/cmd-wisp 落笔）。尺一律走 `git cat-file blob HEAD:`＋`git grep HEAD`＋既有日志＋只读形脚本。
- 被裁对象＝实现腿 `230-r1` 的 8 笔：`c489eacc`① `c73df208`② `495ba902`③ `1e497544`④ `19d2f92c`⑤
  `dcaf3995`⑥ `c1349524`⑦ `ca0ad6bd`⑧。名册尺＝`git log --grep='^票 230 落地腿 230-r1'`（锚定行首）＝**8 笔**〔量〕，
  锚⇄非锚同值＝8⇄8（名册尺灵）。非锚定子串尺⛔ 用（公共块 A1045 同族）。
- 实现腿自报一律当未验证；下文每枚数带【尺名｜射程｜量/算】。落档＝`.scratch/wisp/probes/230/v1/`
  （`00-index.md` 尺册／`01-bytes-blobs.txt`／`02-truth-table.txt`／`03-scope-discipline.txt`）。

## 1. 总裁
**`AC#2` 成立**。票面判据两味＝「改完注释与盘上形状一致」＋「门绿」：
- 前一味＝本腿逐枚现读＋自算真值表，**复现成立**（第 2 节、第 3 节）。
- 后一味**分档**：格式两把＋长度门＝本腿亲跑复现（第 4 节）；build／d22scan／vet＋定向色＝本腿无授权复跑，
  腿报 rc 全 0 与 `212/168/0/1⇄212/168/0/1` ＝**未验证**（第 5.4 节）。若编排者要按"门绿"字面收齐，缺的就是这一味。
- 三处自报与盘上不符／过位＝具名见 第 5.1 节／第 5.2 节／第 5.3 节，均不破 第 2 节 的实质判据。本腿**⛔ 翻框**。

## 2. AC#2 的实质＝那族"唯一读者"陈述今天是不是实话（不是"格式改了"）
定位尺＝内容锚 `only reader of q.alias`（⛔ 绝对行号；票面 `:247` 已漂）：现量＝字段文档 `:86-:93`（声明 `:94`）
⇄ 函数文档 `:256-:275`（函数体 `:276`）〔量＝`git cat-file blob HEAD:internal/agent/approval/queue.go` 逐行现读〕。

`alias` 真值名册（本腿自算，尺＝blob 内 `q.alias` 逐枚命中＋enclosing func 用 awk 现推）：
- `indexLocked` 定义 `:213`，读 `:215`＋写 `:218`；被 `push`（定义 `:159`）体内 `:183` 调用。
- `unindexLocked` 定义 `:226`，读 `:228`＋删 `:234`；被 `dropLocked`（定义 `:305`）体内 `:313` 调用。
- `lookupForRefusalLocked` 体内 `:283`＝全仓**唯一**按"回答的名字"检索该 map 的读点（`:280` 先读 `byID`＝键先赢）。
- 调用侧：`Queue.reject` 定义 `:502`、体内 `:504`＝唯一非注释调用点（尺＝`git grep 'lookupForRefusalLocked' HEAD` 逐枚
  判注释性；命中另有 `:276` 定义行⇄`.scratch/wisp/probes/259/r2/before/queue.go` 两枚死副本＝⛔ 编译⛔ 包内，具名剔）。
- 全包／全仓其余＝0（字段非导出；`internal`＋`cmd` 之外零命中〔量〕）。`:118` 是 `NewQueue` 的初始化写（分 map），
  非读者，名册不必列它，列了也不错——腿没列，本腿只登记。

字段文档逐句（`:86-:93`）：S1 reject-direction 索引⇄`otherNames` `:192` 取 `d.CorrelationID`／`d.TaskID` 剥空与同号 ✓；
S2「Of the readers, only lookupForRefusalLocked resolves a reply into an item」⇄上面名册 ✓；
S3「indexLocked／unindexLocked touch the map too, but each only moves the one item handed to it」⇄`:213-:221`⇄`:226-:236`
只碰传入 `it` 自己的桶 ✓；S4「The allow side goes through lookupForAllowLocked, which cannot see this field」⇄
函数体逐字 `if corr == "" {return nil}` ＋ `return q.byID[corr]`（`:249-:254`，体内只 `byID`）✓；
S5「the index can only ever make a REFUSAL land」⇄`:283` 唯一出口是 `reject`，`reject` 体内方向字面量唯一＝
`tools.AnswerReject`（`:517`）✓。

函数文档逐句（`:256-:275`）：F2「the only reader of q.alias **that turns a reply into an item**」＝限定后的主张 ✓；
F3「NOT the only reader in the file - indexLocked (run by push) ⇄ unindexLocked (run by dropLocked)」✓（两个宿主名
本腿现推 enclosing 对得上）；F4「neither resolves a reply and neither can settle anything」＝二者均不调 `deliver`
亦不返回条目 ✓；F5「Queue.reject is the only caller of this function」＝按上面调用点名册 ✓（票面旧句这半句
本来为真，腿把它保住并写了"holds as written"，诚实）；F6「lookupForAllowLocked, whose body reads q.byID and
nothing else」✓；F7「pinned by TestAnAliasCanNeverBuyAnAllow」＝**注释里指名的用例与盘上真名一字不差、且在位**
〔量：注释 `:275` ⇄ `git grep -n 'func TestAnAliasCanNeverBuyAnAllow' HEAD`＝
`ticket97_alias_direction_test.go:75`；腿件 `00-report.md:90` 那枚名也逐字核过＝同串 ⇒ ⛔ 在注释里指名一枚
不存在的测试当凭据那一味也过〕，该用例的**绿**＝既有日志⇄CI 旁证（本腿⛔ 重跑）。
F8「Caller holds q.mu」⇄`:503` ✓。

⇒ **实质判定：那两句注释今天逐句为实话**，走的是票面第二支（把其余读者的存在具名解释掉＋限定"唯一"的射程），
票面第一支的字面行号（`:197`／`:210`）已漂，腿用函数名替代＝合规（派单 15:2x 逐字⛔ 用绝对行号）。

## 3. 可执行字节变动＝0（★本腿另起两把独立尺，⛔ 抄腿的那把）
尺 V1-A〔量，`git show -U0 <commit>` 逐变动行 awk 判 `^\s*//`，对象＝入库 diff〕：
- 笔①＝变动 14 枚（加 11／删 3），非注释 **0**；笔②＝变动 7 枚（加 5／删 2），非注释 **0**；累计 **21 枚／0 枚**。
  与腿报 14⇄7⇄21 同值——本腿自己算的，独立复现。
尺 V1-B〔量，一律 `git cat-file blob <rev>:<path>`，⛔ 工作树〕：
- 剥整行注释（`grep -vE '^[[:space:]]*//'`）后正文：`c16ff47e`⇄`c73df208`⇄`HEAD` ＝ **421⇄421⇄421 行**，
  `diff` RC=0／0 字节，sha256 两侧同串 **`fadaf8d1dd30abb511569184c48bb63c5b04f2ceecaf2eaa29a8d8e6dba5089f`**（＝腿报
  `fadaf8d1…5089f` 全串复现）。再剥空行（更严的一把）＝382⇄382、sha256 两侧同 `c8645d14…`。
- blob 尺：`git rev-parse c489eacc^:queue.go`＝`66fec7ae…`＝`git rev-parse c16ff47e:queue.go`（起手锚⇄笔①父同字节，
  中间别腿的笔⛔ 碰过这枚文件）；笔①＝`8751a40f…`；笔②起＝**`dcef3fa8…`⇄③④⑤⑥⑦⑧⇄HEAD 八枚同值**；
  `git hash-object`（对 `cat-file` 落的仓外副本）＝`dcef3fa8…`＝`rev-parse HEAD:` 同值 ⇒ 续跑四笔对 `.go` 零字节〔量〕。
- **CR／码位一律对入库 blob 量**：本机 `git config core.autocrlf`＝`true`〔量〕＝工作树随时变 CRLF，
  对它量会冒假警报（编排者今天栽过的同形）；本腿全程⛔ 碰工作树。

## 4. 门禁：本腿能跑的部分（⛔ 编译类）
- gofmt⇄gofumpt 并排两把〔量，对象＝仓外副本 `blob/queue.go`＝HEAD blob 逐字节〕：`gofmt -l` 列 **0** 枚 rc=0 ⇄
  `/d/work/base/gopath/bin/gofumpt.exe -l` 列 **0** 枚 rc=0。**尺灵正控**＝故意写坏的 `badctl.go` 被两把都列出〔量〕。
  起手锚侧同 blob 两把也列 0 ⇒ 改前改后同形（并排＝腿报"并排两把 rc=0 列 0"复现）。
- 长度门〔量，`sh scripts/check-path-length-budget.sh` 裸形＝纯 `git ls-files`＋awk，⛔ 调 go，只读；跑前
  `unset GOTMPDIR TMPDIR TEMP TMP`〕：**rc=0 VERDICT GREEN**，denominator=17219、over-budget=60＝roster=60、
  not-in-roster=0 ⇒ 腿报 rc=0 GREEN 复现。
- 裸 `gofumpt` 不在 PATH〔量：`command -v gofumpt` rc=1〕⇒ 腿报的 rc=127 形状属实、⛔ 当"跳过"它自己也没当。
- 码位带一把〔算，python 按 `tools/d22scan` 源码声明的 band 表扫 `queue.go` blob〕：命中 **0**。
  d22scan 射程（`emojiScopes`⇄`walkGo` 源码〔量〕）＝`design/`＋`frontend/`＋`internal/`(go)＋`cmd/`(go)
  ⇒ `.scratch/**` 与本判语⛔ 仪器射程＝门料不进来（本腿正文件的码位形状据此起）。

## 5. 自报⇄盘上不符处（具名，⛔ 替它圆）
1. **"全射程"措辞过位**：票尾 bullet 写「★AC#2 全射程复核…同族扫描 `sole reader|only user|唯一读者|唯一调用者`＝
   0 命中⇒无第三处假话残留」。本腿把它自己落档的尺挖出来对质（`sh/01-metric-exec-bytes.sh` 第 68 行，只读旁证）：
   那把尺的射程＝**只 `internal/agent/approval/queue.go` 一枚**；同射程本腿复跑＝0 命中 rc=1 ✓（数对）。
   但本腿的更宽射程〔量〕：包级 `-- internal/agent/approval` ＝ **1 命中**＝`gate.go:595`（"its only user verb is veto"）；
   产码树 `-- internal cmd`（词界修正、剔 `con|sole reader` 子串噪声 4 枚）＝同这 1 枚；全仓任意文本＝134 枚
   （多为腿件/台账转述）。⇒ **"0 命中"那枚数对、"全射程"那句名不副实**。实质结论另判：`gate.go:595` 那枚
   今天**为真**（第 6.2 节），所以"无第三处假话残留"这件事**成立**，成立的理由不是那个 0，是本腿的宽射程＋逐枚真值。
2. **"锚 0 命中＝已改"证据力过位**：单行锚对 HEAD＝0 命中 rc=1〔量〕、对 `c16ff47e`＝1 命中（`:262`）rc=0〔量＝正控，
   尺灵〕——两向都复现。但本腿新起一把诚实尺〔量〕：把行首 `//` 标记剥掉、按行拼接后，HEAD 上仍能拼出
   `only reader of q.alias that turns a reply into an item`（1 命中）。⇒ 那个 0 **部分是换行伪证**，
   "谎话被删"的真凭据是 第 2 节 的限定语真值表，不是这枚 0。（起手锚侧拼接尺同样复现原文，两向自洽。）
3. **腿件⇄正文两处小漂**：bullet 报的同族正则与其 `sh/01` 实际正则差一枚 `sole (reader|user)`；四枚腿件内
   **没有**同族扫描的逐字输出落档（只有 `FAMILY_RC=1`＋结论句），射程要从仓外临时根脚本才读得见——
   盘上可核性弱了一档（本腿已替它把逐字射程钉死如上）。
4. **未验证残留**：build／d22scan／vet／定向色四把本腿⛔ 复跑（硬约束）。旁证只有：单变量前提成立
   〔量：`git diff --name-only 3eeaa2e2 c73df208` ＝ **只剩 `queue.go` 一枚**；`c16ff47e`⇄`3eeaa2e2` 的 `queue.go`
   blob 同哈希 `66fec7ae…`＝base 与起手所见字节逐字同，腿 00-report 第 6 节那两枚前提复现〕；本机 ci-logs
   不覆盖这 8 笔（⛔ push＝无远程 run）。⇒ "门绿"里这三把的绿**仍以腿报为据**，编排者翻框时请按此分档。
5. 一条更正（⛔ 记在腿账上、记在本腿账上）：本腿初稿 第 5.5 节 曾写"腿件 `00-report.md:90` 用例名多抄一枚 In"——
   本腿逐字重跑（`00-report.md:90` 原文 ⇄ `queue.go:275` ⇄ `git grep 'func TestAnAliasCanNeverBuyAnAllow'`）
   ＝三处**一字不差**，那枚是**本判语自己的抄写滑笔**，已就地改正、具名登记（判语件也走"只追加＋更正注记"）。
   盘上无不符。

## 6. 腿具名请裁的两件
### 6.1 票面无 `## Progress log` 节⇒不新造、改照同形先例落一条自带 `did=` 的 bullet＝算守"票尾只追加一条"吗
**裁＝算，守住了**。量证：票 230 `grep -c '^## Progress log'`＝**0**〔量〕（池内 242 枚有该节，它没造）；
先例枚数＝本腿逐字重跑腿的 `sh/17` 那把尺〔量，射程＝`.scratch/wisp/issues/*.md` 工作树现读〕＝
TOTAL_NOLOG_WITH_BULLET=**12**，减去票 230 自己那条（腿跑尺时它还没落）＝**11**＝腿报数复现。
追加纪律尺〔量〕：笔⑦ numstat＝**3/0**（纯追加）；`c1349524^`⇄`HEAD` 对票 230 逐行＝**ADDED=3 DELETED=0**
（票面原 40 句一字未动）；票名⇄`Status:` 前 5 行 `c16ff47e`⇄`HEAD` diff 空；票尾 `^- \[2026` 条数＝**1**（唯一）。
笔⑧＝1/1，只重写它自己刚追加的那一条（走追加一笔、⛔ amend⛔ 改写入库历史）＝合规，追加枚数仍＝1。
时间戳形〔量，三把〕：池内 `[YYYY-MM-DD HH:[0-9]x +08] agent=` 形＝**39** 枚文件；光 `HH:x` 形＝**0**；
真实分钟＝10 ⇒ 腿自抓「`18:x` 是独形」为真、改回主流形方向正确。
### 6.2 射程外两枚同族陈述，腿判"今天为真"且⛔ 动字节＝那判断成不成立
**裁＝两枚判断都成立，且"⛔ 扩写面"守对了**（本腿没扩任何写面去"帮"它们落字）。
- `gate.go:661-662`「The alias index is not in this path, so a refusal alias cannot be used to point a grant at
  some other call」真值链〔量〕：`allowSession`（`:658`）→ `sessionSubject`（`:663`，体内 `:468` 一跳＝
  `lookupForAllowLocked`）＋ `allowScoped`（`:670`，体内 `:396` 同跳）；`lookupForAllowLocked` 函数体逐字只读
  `q.byID`；`q.alias` 全部读点＝`:215/:218/:228/:234/:283`，无一在此路；字段非导出、包外零读。→ 为真。
- `approval.go:496-501`「No routed call can produce this value（denialMisbound）… the only callers that reach it
  are inside this package and hand spend two different digests by hand: TestTicket242…And TestTicket259…」真值链
  〔量〕：`spend` 非测试调用点唯一＝`queue.go:425`，传 `it.bind`＝值对自身；`it.bind` 唯一赋值＝push `:180`；
  `grantNonce` `:369` 把同一枚 `it.bind` 签进同store；`allowScoped` 的 routed 调用者恰两枚（`queue.go:377` 传 false／
  `gate.go:670` 传 true）；以手写异 digest 真打到 `denialMisbound` 的调用点＝`ticket242_binding_test.go:169-170`⇄
  `ticket259_denial_rulers_test.go:170-171`（各带 `!= denialMisbound` 断言），两文件 package 行＝`package approval`
  ⇒ "inside this package" 字面有据（in-package 测试件、与 `grantStore.spend` 定义同包）。→ 为真。
- 零字节复核〔量〕：8 笔逐枚 `git show --format='' --name-only` 名册里 `gate.go`／`approval.go` 出现 **0** 次（第 7.3 节）。

## 7. 写面自证＋框终量（本腿的，逐枚）
1. 本腿落笔＝一枚判语＋四枚尺档＋票 230 尾一条只追加 bullet，每落一枚即一笔 commit（枚数＝件数，⛔ 攒）；
   逐枚 `git show --format='' --name-only` 的 verbatim 名册落在各笔正文里，**每笔 `.go` 枚数＝0**、
   写面只两枚根（`docs/evidence/s1/230-v1-verdict.md` ⇄ `.scratch/wisp/probes/230/v1/**`）＋票尾一条。
2. 实现腿 8 笔名册〔量〕＝①② 各只含 `internal/agent/approval/queue.go` 一枚；③④⑤⑥ 各只含自己那枚
   `.scratch/wisp/probes/230/r1/*`；⑦⑧ 各只含票 230——越界 **0**、`.go` 只在 ①②、`*.out` 枚数 **0**
   （`.gitignore` 第 8 行 `*.out` 在场〔量〕＝腿用 `.txt` 的理由成立）。
3. 票 230 框终量＝**OPEN(未勾)=5／DONE(勾)=0**〔量，`git cat-file blob HEAD:…| grep -cE '^- \[ \]'`＝5；
   `grep -cE '^- \[x\]'`＝0 rc=1〕＝本腿一枚没翻（起手侧 5/0 同值，`AC#2` 的翻勾归编排者凭本判语做）。
4. ⛔ push／⛔ 改票名与 `Status:`／⛔ `add -A`／⛔ amend/reset/stash/checkout/clean／两处显式 pathspec／
   `-F` 用 Windows 形正文／临时件只建不删（`D:\tmp\wisp\20261020\230v1\`，尺与原始输出全在）。
5. 本判语自身 ≤160 行、非 0 字节、⛔ `.out` 扩展名。

## 8. 交件必答题（编排者收件用一句话版）
`AC#2`＝**成立**；缺的一味＝build/d22scan/vet/色的绿⛔ 第二人复跑（本机硬约束＋⛔ push 致 CI 无覆盖），
按 第 5.4 节 分档翻框；两处措辞（第 5.1 节 全射程、第 5.2 节 换行伪证）建议编排者落台账时把"实质成立凭据"记成本腿 第 2 节/第 6.2 节
的逐枚真值，⛔ 记成那两个 0。腿请裁的两件（第 6.1 节／第 6.2 节）本腿都裁＝合规／成立。
