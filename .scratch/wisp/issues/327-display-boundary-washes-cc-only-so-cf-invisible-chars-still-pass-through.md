# 327 — 显示边界只洗 `Cc`（C0/C1），**`Cf` 那批看不见字符今天照样穿过进菜单**：`U+202E`（RTL OVERRIDE）正是票 216 票面写的"盖掉后半段名字"那一形

- Status: **已立，待派**（2026-10-16 12:4x 编排者立；来源＝落地腿 `216-r1` 交回的具名残余 ②（件 §7）＋台账 `A1006`；票 216 的 `AC#2` 判据原文写的是"控制字符／可选可见转义"，`Cf` 不在那句射程里 ⇒ ⛔ 塞进 216 已勾的格子＝改判据，另立本票）
- 类型：**行为缺陷**（显示边界字符集窄了一档）。射程＝`internal/panel/**`，⛔ 动 `cmd/wisp`（共享编译面）／`frontend/**`／`design/**`（归机主的外部前端 agent）
- ⚠ 行号只当**读数**、⛔ 当判据（本仓行号锚会随别人加行腐烂）；下面每条都附可复跑的尺。

## 这是什么（人话）

票 216 已经装好一道"显示前洗一遍"的边界：文件名、目录名里夹着**控制字符**（换行、`ESC` 这些）时，界面上看到的不是原样，而是被写成看得见的 `\x1b` 那种转义——所以"照片␛已允许.png"不会再看起来像"已允许"。

但有一类字符它**不管**：**看不见的排版字符**（`Cf` 类）。其中最要紧的一枚叫 **RTL OVERRIDE**（`U+202E`）：它不是字母也不是符号，它的作用是"从这儿起，后面这段文字**倒着**排"。文件名里夹了它，屏幕上就会出现**前半段是真名、后半段被它盖掉或反着读**的样子——这正是票 216 票面举的那个形状，今天**仍然能穿过**。另外两枚 `U+200B`（零宽空格）与 `U+FEFF`（BOM）更安静：它们在，界面上一个像素都不显示，两个看起来一模一样的名字其实不是同一个名字。

## 现量（起手逐条复算，⛔ 信这里的号）

| 事实 | 读数（编排者 2026-10-16 12:4x 现跑，尺附在旁边） |
|---|---|
| 洗的那道只认 `Cc` | `internal/panel/display_scrub_216.go` 的 `scrubForDisplay` 里唯一那枚类别判据是 `case unicode.IsControl(r)`；`unicode.IsControl` **＝类别 `Cc`**（C0 `U+0000–U+001F` ＋ C1 `U+007F–U+009F`），`Cf` 不在里面 |
| 快路径会让"只含 `Cf` 的名字"一字不动地回去 | 同件 `needsDisplayWash` 的判据与写入侧是**同一枚字符集**（`IsControl` ＋ 非法 UTF-8 字节）⇒ 名字里只有 `U+202E` 时它报"不用洗"，`scrubForDisplay` 直接 `return s` |
| 这处窄是**写件里自己承认的**，不是漏 | 同件 `:52–54` 的注释逐字写着"it does not wash Cf (format) characters such as U+202E RTL OVERRIDE or U+200B ZERO WIDTH SPACE … Widening to Cf is a real question"（`216-r1` 交回，我 12:4x `sed` 读回原文） |
| 双向都是 0：全树 `Go` 源里**没有任何一处**认得那几枚双向控制符 | `git grep -nE '0x20(2E|66|67|68|69)|\\u202[eE]|\\u206[6-9]' -- internal cmd tools`（滤掉 `internal/risk/taintmatch.go` 那枚不相干的行）**命中 0 枚文件** ⇒ `U+202E` 与 `U+2066–U+2069`（隔离符一族）今天既没被洗、也没被判定、也没被任何用例喂过 |
| 但**同类问题在判定层已经有一份名单**（可抄的先例，⛔ 复用同一份函数） | `internal/risk/taintmatch.go:40` `isIgnorableTaintRune`＝`0x00AD, 0x200B, 0x200C, 0x200D, 0x2060, 0xFEFF` ＋ `unicode.Mn`，用途是"伪授权里的看不见字符不许顶掉词形"（对抗报告 `16.9#1/M-6` 的产物）。★它**不含** `U+202E` 一族，且它是**判定用**的（问"这算不算一个字符"），⛔ 拿它当显示洗的名单——两件事的取向相反：判定要"忽略后还认得"，显示要"显出来让人看见" |

## 判据（AC）

- [ ] **AC#0 现量普查（⛔ 没名册不许进实现，本仓加尺的固定前置）**：逐枚列出今天真能带着 `Cf` 走到界面上那几枚字段（起点＝票 216 `AC#1` 的名册：`AttachmentRef.Name`、`WorkspaceView.Canonical`；⛔ 顺手扩到"以后可能有"的字段），并对每枚给**三形读数**：① 只含 `U+202E` 的名字进到快照里是什么字节、② 只含 `U+200B` 的、③ 两个只差一枚 `U+200B` 的名字在快照里是不是**逐字节相同**（那才是"看起来一模一样其实不是同一个"的证明）。尺附在件里、可重跑。
- [ ] **AC#1 洗的范围要成文**：`Cf` 里**洗哪些、不洗哪些**必须一张表写死（建议起点：`U+202A–U+202E`、`U+2066–U+2069` 这九枚方向控制＋`U+200B/U+200C/U+200D/U+2060/U+FEFF/U+00AD` 这六枚零宽与软连字符；`U+2061–U+2064` 数学作用符**具名裁要不要**）。⛔ 一句"洗所有 Cf"含糊过关——`Cf` 整类还有 `U+061C`、`U+180E` 等，逐枚要有名字。
- [ ] **AC#2 ★两把尺不许分家（本票最容易自己骗自己的一格）**：写入侧（`scrubForDisplay` 的 `switch`）与判据侧（`needsDisplayWash`）**必须是同一枚字符集**。完成判据＝**定向突变**：把 `needsDisplayWash` 里新增的那一支摘掉、⛔ 动写入侧 ⇒ 那枚"只含 `U+202E`"的用例**必须红**（它证明快路径真在把关）；反向（只摘写入侧、留判据侧）也必须红。红名册与红句逐字进件。⛔ 把这条读成"两边都加就完事了"——今天这两支是同一个人同一批写的，**正因为长得像才会漏**。
- [ ] **AC#3 洗＝转义成看得见的形，⛔ 删**：与票 216 同一纪律（删会把两半拼成一句新的话，`abc`＋`U+202E`＋`def` 删完变 `abcdef`＝造出一个原本不存在的名字）。⇒ 输出必须**可逆读出**（`\u202e` 那种形），且**幂等**（洗过的串再洗一次逐字节不变，这条要有用例钉，因为它就是 216 那件注释里"同一个调用坐在快照边界上不会双重转义"的同一枚性质）。
- [ ] **AC#4 判定不许被名字里的字改变**：接票 216 `AC#3` 那一半的同一形状——把文件名做成 `photo\u202e已允许.png`，审批**判级/判定必须与不含 `U+202E` 的同名文件逐字相同**（真 `risk.NewRiskAssessor()`，⛔ 桩）。★这半枚要自带**攻它的那一发**：把判定改成"名字里含'允许'就放行"，用例必须红（票 216 那半我 12:3x 只跑到"存在且绿"、没攻过，本票⛔ 重犯）。
- [ ] **AC#5 门禁四数＋颜色**：`gofumpt`／`tools/d22scan`／长度门／`go vet ./internal/panel/` 逐枚落 `rc=`；`./internal/panel` 整包**改前⇄改后同一台面**成对、显式 `-timeout 30m`、四数逐枚点名（`^=== RUN` 含子测试 ⇄ `^--- PASS` 只数顶层，⛔ 混用），**终态红名册与起手红名册作差＝0**（起手名册＝本波 `216-r1` 交完后的 `151 PASS／5 FAIL／RUN 246／SKIP 0`，编排者 12:3x 现跑；⛔ 抄，起手自己重跑）。⚠ 界面代码零 emoji：新增的注释与文案⛔ 落进 ban #8 的射程（仪器实际扫 `U+1F000–U+1FAFF`·`U+2200–U+22FF`·`U+2600–U+27BF`·`U+2B00–U+2BFF`·`U+FE0F`，注释**豁免**、字符串**不豁免**）——`≤` 与 `✓` ⛔ 出现在字符串里。

## 边界与禁区（⛔ 混进来直接退回）

- ⛔ **动判定层的 `isIgnorableTaintRune` 去迁就显示名单**（`internal/risk/taintmatch.go:40`，那是对抗报告 `16.9#1/M-6` 钉下来的形状，动它要先有它自己那张批准账）；本票的名单**长在 `internal/panel/**` 显示侧**，两份名单各写各的、由 AC#0/AC#1 的名册说清差异。
- ⛔ 新增加进 Snapshot 名册的键（撞票 313 那枚键钉）、⛔ 新 C17 方法名、⛔ 新导出 `panel.*` 名 ⇒ 需要就停手上报（`216-r1` 的残余 ①③ 就是这个形状，我没批）。
- ⛔ 为变绿放宽任何既有断言：票 216 那 7 枚新用例（`display_scrub_216_test.go`）里 `two_escs_keep_both` 等字符类子用例是**近邻同形**，扩集合时最容易顺手改它——要动先停手上报。
- ⛔ 把 `frontend/**` 拖进来（页面那侧的渲染归机主的外部前端 agent）。本票只做 **Go 出向边界**。

## next

`327-a1`（只读普查腿＝AC#0 的名册与三形读数）→ 我裁 AC#1 那张表 → `327-r1`（落地腿，写面收成 `internal/panel/**` 一枚，含 AC#2 那对定向突变自证）。**释放条件**＝`internal/panel/**` 这面本波空着（`216-r1` 已于 12:12 交完）＋ CPU 连发读数 ≤70%（机主 10-13 定的补位门槛）。
- [2026-10-16 15:1x +08] **`AC#0` 只读普查腿 `327-a1` 交件**（⛔ 判语、⛔ 修法建议、⛔ 翻框、⛔ 动票名/`Status:`、⛔ 排程；三节料在三枚件里：`.scratch/wisp/probes/327/a1/census-1-roster.md` 37 行／8555 字节、`census-2-cf-and-upstream.md` 78 行／12294 字节、`logs/readings.txt` 154 行／8208 字节；锚＝具名提交 `44f51fc8`，⛔ `HEAD~n`；读数一律取自**工作树现行字节**）。头条五条读数：① **显示边界今天只碰两枚字段**——生产调用点全仓 1 枚（`internal/panel/composer.go:145`），里面只有 `display_scrub_216.go:169`（`Workspace.Canonical`）与 `:175`（`Attachments[].Name`）两处赋值 ⇒ `Snapshot` 另有 7 枚字段今天出向且不过这道边界（`pending[].tool/args/reason`、`results[].text`、`composer.workspace.reason`、`composer.attachmentError`、`tasks.rows[].label`、`instructions.files[].path`、`git.*`）。② **票 216 名册里的"附件名"今天两头都断**：入向腿在装配根即 `cmd/wisp/panel_inbound.go:301 Attachment: nil`（`:302 Message` 同形）＋泵 `pump.go:275` 传 `nil`＋`AttachmentBroker.Ingest` 生产调用者 0 枚 ⇒ 那枚字段是"洗法已就位、源与出向 0"，与票 216 `216-r1` 残余 ③ 同形。③ **`Cf` 候选 22 枚逐枚可达、无一枚"仅理论形"**：入口 A（磁盘真名 `panel_pump.go:128`→`workspace.go:119`→`display_scrub_216.go:169`）与入口 B（模型正文 `internal/llm/events.go:239`→`tools/bridge.go:382`→`panel_pump.go:68/297-312`→`approval.go:80`，及 `run.go:1331`→`pump.go:537`）都**不按字符筛**（`pathresolver.go` 里字符类拒绝 0 处），`needsDisplayWash` 对 22/22 逐枚不响、正控 `ESC`/`LF` 逐枚响。④ **AC#0 的"三形读数"里那条"逐字节相同"按盘上真值更正**：只差一枚 `U+200B` 的两枚名字**不**逐字节相同（10 字节 ⇄ 13 字节，差 `E2 80 8B`），`scrubForDisplay` 后仍不同——"屏幕上同一枚名字、名册里两枚名字"这半成立，"字节相同"那半不成立。⑤ **上游没有替我们洗掉任何一半**：`taintmatch.go:39 isIgnorableTaintRune` 的调用者全在 `internal/risk` 内（`hostpath_185.go:67`、`provenance.go:559/578/841`），`internal/panel/**`＋`cmd/wisp/**` 命中 0；`sanitizeDisplayName` 的 `TrimSpace` 对那五枚零宽一枚都不去掉；`json.Marshal`（`pump.go:344/359`）只转义 `Cc`、22/22 枚 `Cf` `escaped=false`。**必答题（编排者那一问）答"没有一环"**：`U+202E` 从模型出发的 9 环逐环 `file:line` 已在件 2 列表，每一环都逐字带过。★那 9 环里**第 8/9 环的读数改变本票"走到界面"四个字的射程**，具名上交、⛔ 本腿定性：快照字节今天**没有进 WebView 的生产路径**（`run.go:732` 的 `Out` 实现 `bookPanelSnapshot` 只存字段＋一条日志，`panel_pump.go:360-366`；读它的 `lastPanelSnapshot` 生产调用者 0 枚；入向 7 枚方法常量里没有读快照那一枚，`MethodSnapshot` 尺 rc=1；出向 Eval 那形只承载 RPC 回执），页面侧 `App.tsx:91` 的 `snapshot` 是带缺省值的 prop、`panel.ts` 只有入向 `postMessage` ⇒ 本名册的"可达"＝**到快照字节**，"屏幕上哪一行"那一列量的是**渲染器就位**、不是"像素已出现"（同屏另录两句自陈各自的射程：`pump.go:366` 那句"最后一公里"射的是 `src.Out == nil` 分支、`panel_pump.go:440` 的 "does not exist yet" 射的是推送频率规矩）。探针自报：`go env GOVERSION GOROOT`＝`go1.27.1`/`D:\work\base\go`；唯一执行过的代码＝仓外单文件纯 stdlib `D:\tmp\327a1\cfprobe\main.go`（`go run` rc=0，`needsDisplayWash`/`scrubForDisplay`/`writeControlEscape` 三段是 `display_scrub_216.go:100-147,152-161` 的**逐字重打**、⛔ import 本仓包）；⛔ `go test`/`go build`/`go vet`/`gofumpt`/门禁（编译面归 `165-v2`）。写面＝上述三枚件＋本条 bullet；⛔ 一切其它路径（`internal/**`、`cmd/wisp/**`、`docs/evidence/s1/**`、`frontend/**`、`design/**`、`go.mod`、`.gitignore`、判据文字）零字节；⛔ push。
