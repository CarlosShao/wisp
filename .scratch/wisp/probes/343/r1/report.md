# 票 343 交件 · A 半（改了什么／为什么／结论）343-r1 落地腿

读数册＋交件自证＝同目录 `report-sec-b.md`（起手三读数／成对四数／8 枚牙／门禁）与
`report-sec-c.md`（三把尺＋票框终量尺）。拆件口径＝公共块那句"每件 ≤160 行"只管我
撰写的正文件、写不下⇒拆具名分件（`322-r2` 先例）、⛔ 压正文。

## 1 靶子与成因（起手现量，逐枚⛔ 引用票面未重跑的数）
`internal/risk/blacklist.go:193` 那句写 `(aEquiv / aUnder report such a miss as a hit)`。
`aEquiv` 在本包**从未有过定义**：`git grep -c 'func aEquiv' HEAD -- internal/risk`＝
rc=1／0 命中；全仓唯一那枚命中＝这句注释自己。引入笔＝`f1033e1e`（2026-09-21，
`fix(72)` 那条"A 表锚点必须和候选路径走同一条句柄管线"），挂了近一个月。
＝一句指向不存在符号的产码注释（⛔ 行为影响、⛔ 判据影响；AC#4 那组成对四数正是这件事
的正面证明：改它前后四数逐枚同值）。

## 2 我读准的"两把尺"（派单起手第 3 步＝读准了才动笔）
票面要我说清"今天到底哪两把尺在这种缺位上报 miss"。盘上真身＝
```
aEq:271 / aUnder:273             两枚 A 档关系，本体只是一行委托 aRelates
aRelates:280 -> 281-283           if anchorMatches(...) { return true }
                                  return uncertainAnchorMiss(...)  ← 报 miss 为 hit 的是这一支
anchorMatches:292                 双层 spellings() 循环 + rel(...)，找不到就 return false
uncertainAnchorMiss:362 -> 363-369  两侧 root 都非空、且⛔ 两侧都 certain、且两 root
                                  互为祖先后代 ⇒ true（R17 兜底支）
```
⇒ **两把尺＝`anchorMatches`（真证明出来的关系）与 `uncertainAnchorMiss`（证明不出来时
的 fail-closed 兜底）**；而"把一枚缺位报成命中"这句**只属于后一把**。`anchorMatches`
单看恰恰⛔ 这么干——票 334 就是把两把拆开的（`:287` 那句 `anchorMatches is aRelates
WITHOUT R17's tie-breaker`；拆开⛔ 改变任何判级，因为 `aRelates` 仍按老顺序"先匹配、
后兜底"），拆开的全部理由写在 `:289`／`aRuleReason:327`：`aRuleReason` 要说清这句拒绝
到底是哪一把产出的（证明出来的命中 ⇒ 光打规则字面；兜底支 ⇒ 打
`[fail-closed, NOT a match on this rule]` 那一长句）。

## 3 交件句子（逐字，5 行全注释）
```
// provably expanded fails CLOSED: aEq / aUnder delegate to aRelates, which
// reports such a miss as a hit when anchorMatches found nothing AND
// uncertainAnchorMiss held, instead of quietly dropping a path from A into B.
// anchorMatches on its own does not do that - ticket 334 needs the two apart
// because aRuleReason has to say which of them produced the verdict.
```
带条件的事实句；句子里点名的每一枚符号（`aEq` `aUnder` `aRelates` `anchorMatches`
`uncertainAnchorMiss` `aRuleReason`）逐枚 `git grep -nE` 现验都有定义；★交件侧
`git grep -n aEquiv HEAD -- internal/risk`＝0 命中＝⛔ 保留不存在的名字。
★⛔ 新增 `aEquiv` 去让注释变真（那是为凑一句话造一枚符号＝本仓禁手，且会顶红票 336
那 8 枚牙——`report-sec-b.md` 第七节那 8/8 全红正是反证）。

## 4 ★AC#0 我自加的同族尺 ⇒ 命中 2 枚（⛔ 一枚）
票面 AC#0 明令⛔ 只交 `aEquiv` 那一枚就走（票 334 那族"逐枚锚"的教训）。我加了**两把**尺，
因为第一把我自己找到了盲区：
- 尺①`census.exe`（`go/ast`，注释面驼峰标识符对撞包内名册）。名册第一版漏了
  param／field／local ⇒ `bOverrides` 被误报"无定义"，其实它是 `Gate` 的形参（真名）。
  第二版补齐。**它另有一处结构性盲区：⛔ 看见纯小写标识符。**
- 尺②`callshape.exe`＝专挑注释里写成**调用形**的名字（`foo(` / `foo()`），那是在断言
  "存在一枚可调用的东西"，对撞包内 402 枚 func/method。
结果（跑在**我起手的基线树**上＝`blacklist.go` blob `c1799cd2`，⛔ 我的字节）：
```
★同族第 1 枚  aEquiv   blacklist.go:193  驼峰形、全仓 0 定义         ＝本票、已改
★同族第 2 枚  fires()  blacklist.go:307  纯小写调用形、包内⛔ 此名    ＝⛔ 本票射程
  （`fires` 作标识符的真身在别包：internal/config/unwired.go:54 那枚 struct 字段；
   引入笔＝e052b11a 2026-10-17＝票 334 落码腿自己、就在它把 aRelates 拆两把尺那一句里）
```
⇒ **怎么办＝改第 1 枚（本票射程内）＋具名上报第 2 枚**（票面 AC#1＝"那一行注释"＝193；
动 310 行＝把本票射程扩大＝派单⛔ 那条"若还有⇒⛔ 顺手扩大本票射程"）。同一对尺扫**整包**
还捞出 2 枚同族＝`provenance.go:54` `taskScopeID`、`:427` `disposalScope`（全仓 0 枚
IDENT），同样⛔ 我写面、同样具名上报。⇒ 建议编排者拿这四枚开一张清扫票（逐枚一行、
⛔ 动可执行字节），别一次一张大票。

## 5 AC#3 结论＝**不该上通用仪器**（只出结论⛔ 自己上、⛔ 造门）
代价是我量的、⛔ 拍的：拿票面给的形状（驼峰 token 对撞**包内**名册）跑 `internal/risk`
一包＝**6 枚红、其中 3 枚假阳＝50%**（`hostPathBox`／`assessorFor`＝
`internal/tools/bridge.go` 里的真 type 和真 method；`macOS`＝散文里的平台名）；换我第二把
尺（调用形）＝**11 枚红、9 枚假阳＝82%**，假阳真身逐枚点名（Go builtin `len`、关键字
`func`、`DEFERRED(...)` 标记约定 4 枚、`Nutstore` 那家厂商名＋括号、以及 `CloseScope(id)`＝
票 160 **故意删掉**的名字、那句自己写着 "Before it, scope lifetime was TWO independent
functions — OpenScope(id) and CloseScope(id)"＝诚实的历史引用，词面上与假话⛔ 可分）。
根因⛔ 我的分词器不够好，而是**注释对一枚名字做了什么断言，无法从名字本身恢复**：
`CloseScope` 看着正像一枚悬空可调用、`DEFERRED` 看着正像一次调用。任何把它们分开的规则
都是人手写的一条规则＝一份 allowlist；而"新门的补丁是手工维护名册"正是
`tools/d22scan/allowlist.txt` 那形＝冻结件。⇒ 第一周就会被静音的门⛔ 值得立。
**若仍要上的形状（交我立新票、⛔ 顺手做）**＝只对一枚**词表**作硬检查："凡 `internal/risk`
注释里匹配 `^a[A-Z][A-Za-z]*` 的名字必须是本包 func/method"＝全枚 6（`aEq` `aUnder`
`aAnyForm` `aRelates` `aRuleReason` ＋ 那个缺陷本身）⇒ 假阳面被词表大小界住。代价如实＝
它抓⛔ `fires()`（纯小写）、⛔ 泛化到别的文件族，且仍属"契约邻接"＝人工批准。
★我⛔ 造任何门、⛔ 新测试、⛔ 碰 `tools/d22scan/**`。

## 6 我这一笔的真实形状（AC#2 命门的自陈）
- 文件＝1 枚 `.go`（`internal/risk/blacklist.go`）＋ 我探针件＋票尾一条；
- 可执行字节变化＝**0**＝判据⛔ 只报"我只动了注释"，给三把尺：`git diff -U0` 7 枚
  `^[+-]` 行逐枚以 `//` 开头（非注释行 0）／`--numstat` 5 加 2 删／★两侧**剥掉注释后的
  可执行字节 sha256 逐枚相同**（`97796acc…`、265 行）＝这把比"删空格再比"那形硬，公共块
  ⓔ 那条要求的正是内容尺；
- 四枚受保护仪器（票 336 两枚牙件／票 339 仪器／票 334 仪器）base⇄HEAD 逐枚 blob 等值
  （`a429cdae`／`f5f888de`／`e7b13c39`／`5147c189`＝SAME=True，现跑）；
- 车道全程独占：`git status --porcelain -- internal/risk` 起手、交件前、交件后都＝0 行。

## 7 AC#4 一句话（全册在 `report-sec-b.md` 第六节起）
两平台成对四数各自逐枚点名、⛔ 混用⛔ 相减：Windows `247/140/1/1` ⇄ `247/140/1/1`，
Linux 容器 `195/115/1/1` ⇄ `195/115/1/1`；红名册作差 0、`RUN` 名集作差 0；两侧唯一那枚红＝
存量（`TestC26RewriteAccountIsConsumedAtEverySecurityLeg`，静态审计 `internal/panel` 的读者、
读不到本包注释一字；与 339-r1 同名对上）；`TestResolvePerCallBudget` 四发全 PASS＝争用假红
那条未触发；★票 336 那 8 枚牙在我自己重导的树上重打＝**8/8 全红、逐名对得上**、11 次还原
哈希 match=YES；门禁四数（共享树、先 `unset GOTMPDIR TMPDIR TEMP TMP`）gofumpt／d22scan／
长度门／`go vet` 全 rc=0。

## 8 未定案上报（⛔ 自行假设、⛔ 盖章）
1. ★`fires()`（`blacklist.go:310`）＋ `taskScopeID`／`disposalScope`（`provenance.go`）＝
   三枚同族名⛔ 本票射程 ⇒ 见第 4 节，要不要并一张清扫票、拍归你。
2. 我这 Windows 四数与上一腿（339-r1）报里那组**逐字相同**（247/140/1/1），Linux 我这侧
   195/115/1/1 ⇄ 它交件侧 195/115/1/1 ⇒ 分母跨腿稳定；但票面 AC#0 那句"现量"给的是
   `blacklist.go` 的行⇄名、⛔ 给四数，⛔ 需要我判它俩谁对——一律按现跑为准。
3. 长度门那枚"每枚 .log 都进 commit 会不会把仓库撑大"：我随件入了 14 枚原始 stdout
   （约 180 KB／3023 行）。前一腿（339-r1）同形先例＝它把 `linux-ac0-fullpkg.log` 等
   26 枚 `.log` 都跟踪了 ⇒ 我按先例走、⛔ 自创口径；若你认为该只留摘要，退我那半即可。
   ★反面理由＝成对读数与产生它的机器输出⛔ 拆两处，拆了编排者复跑对不上；而⛔ 为凑行数帽
   截断 `go test` 自己的输出＝篡改读数（同票 262 长度门自己那句"rewriting a reading is
   how a record stops being evidence"）。
4. ★件数帽那条我按盘上事实自陈、请你裁：我撰写的正文件＝3 枚 report ＋ logs/01–04 四枚读数件
   ＝7 枚，超派单"≤4 枚正文件"那句；走的是同一句里给的那条路（"报告写不下 ⇒ 拆具名分件、
   ⛔ 压正文"，`322-r2` 先例），且每枚都在 ≤160 行帽内（现量 123／136／46／105／47／156／105）。
   提交笔数按"每落一枚就 commit"落到 6 笔＝公共块 14:3x 已把"4–8 笔"降为参考值、⛔ 判据
   （`336-v1` 11 笔不退件那个先例）。要合成更少枚数我也能重排，⛔ 我自己判这算不算缺陷 ⇒
   摆给你裁。
5. ★⛔ 翻框（六格终量 6 未勾／0 已勾，现跑）；⛔ 改票名／`Status:`／`-done`；⛔ push。
