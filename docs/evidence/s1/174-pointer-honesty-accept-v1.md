# 174-v1 验收表（非实现者）：票 174 **AC#2 那一格**——八枚自写判据里有没有一枚是恒真的，"安静正例"是不是装饰

- 派单＝`.scratch/wisp/dispatches/2026-09-27-215x-accept-174-r1-v1-eight-criteria-any-tautological.md`
- 量取时刻：2026-09-27 22:02 → 22:2x +08｜分支 `dev`｜起手 HEAD＝`748d027e0f27a84b3b9da9d609f058ad1a042c44`（我自己 step 0 现量，未抄别人的号）
- 被告＝`174-r1`（已完成）：产码 `b5a07961`＋`f55ddd3d`｜判据 `9ae0be78`｜读数件＋台件＋票面 `bb3eae02`／`ad363d58`｜表 `docs/evidence/s1/174-pointer-honesty-r1.md`
- 本程写面＝本文件＋`.scratch/wisp/probes/174/v1/**`＋票 174 Progress log 一枚追加。**产码零字节**（六枚变异全部走还原协议，末次 `git status --porcelain -- internal/tools/` ＝ **0 行**，见 §4）。

## 0. 起手五件（逐字现量）

| 件 | 读数 |
|---|---|
| `date` | `Sun Sep 27 22:02:58 CST 2026` |
| `git rev-parse --abbrev-ref HEAD` | `dev`（＝要求值，未停手） |
| `git rev-parse HEAD` | `748d027e…1a042c44` |
| `git status --porcelain -- internal/tools/ cmd/wisp/` | **空** |
| 基线 `go test -count=1 ./internal/tools/` | `ok github.com/CarlosShao/wisp/internal/tools 17.951s`（名册口径见下） |

- 名册口径分开报（`-v` 那一发，`.scratch/wisp/probes/174/v1/final-clean.txt` 是全部还原之后的复算）：
  顶层 `--- PASS`＝**124**｜含子测试 `=== RUN`＝**173**｜`--- FAIL`＝**0**｜`--- SKIP`＝**0**｜`^ok  ` 行＝**1**。
  命令：`go test -count=1 -v ./internal/tools/ > final-clean.txt 2>&1; grep -c '^--- PASS' final-clean.txt; grep -c '^=== RUN' final-clean.txt`
  ⇒ 派单那句"PASS 124"**成立**（顶层口径），但同一包在"含子测试"口径下是 173——两把尺不同名。派单 21:44:53 的 `ok 17.899s` 我这边三次分别 `17.951s`／`15.873s`／`16.216s`（时长非断言，不列为不符）。
- 八枚＝真：`git grep -h '^func Test' 51c32ef9 -- internal/tools/` 116 枚 → 工作树 124 枚，`comm -3` **左栏空**、右栏正好 8 枚（台件 `roster-before.txt`／`roster-after.txt`）。

## 1. 本格票面原文（三块，逐字抄，不改一字）

- **AC#2 原句**（票 174 第 31 行）：`- [x] **AC#2 指针不许许诺一条读不回来的路（⚠ 这一发在未修码上必须"不响"）**：构造"回执里的路径**不在授权根内**"那一发 ⇒ 今天 `task.output`（与 `spill.go` 的桩）**静默把路径写进回执**、零提示＝**不响**＝本格成立；修完之后回执必须**外部可见**地说明"这条路径你现在读不到，需要用户先授权哪个目录"（措辞归实现程，但**不许**只在日志里说、不许只在 Go error 里说——模型看到的那一段才算数）。⚠ **判据不许写成"以后再说"**。`
- **20:1x 追加第二形块**（第 32 行）：要点＝`task.go:227` 只有 `rec.ArtifactPath != ""`、没有任何 stat ⇒ 不存在／是枚目录那一形走"有指针"那一支、`不可找回` 不响；射程明写两形 **(a) 授权根外** ＋ **(b) 不存在／不是文件**，两形未修码都必须不响、修完都必须响。
- **21:4x 翻勾块**（第 33 行）：判语＝〔成立·带条件〕，两形都响、三枚变异自证、`C26 边界守住`、四枚门禁重跑，两条条件各立新格（端到端只会说"未接线"／`spill.go` 仍静默），并登记被告自报的两处不实、驳回 `--amend` 手段。
- ⚠ **AC#2b／AC#2c 是尚未做的两格**（票面第 35／36 行），按派单**不计入本格缺陷**；本表只在 §8 说明它们与本表发现的接口。

## 2. 主菜：八枚判据逐枚恒真性（每一枚都要答"产码怎么改会让它红"）

判据源文件＝`internal/tools/task_output_pointer_notice_test.go`；被测产码＝`internal/tools/task.go` 的 `pointerNotice`（`258-328` 区）与指针那一支（`241-252`）。
"我动不动了它"＝本程自己的变异（`MUT-V1…V8`，全在 §3/§4），不是被告那三枚的复算。

| # | 判据（行号） | 它钉什么（断言级凭据） | 产码怎么改会让它红 | 我动不动了它 | 恒真？ |
|---|---|---|---|---|---|
| 1 | `TestPointerOutsideAuthorizedRootSpeaks`（:105） | 形状 (a) 出 `注意：`＋`读不到`（:122）＋逐字"不在你被授权的目录范围内"（:125）；指针 `全文见` 仍在；**同一枚桥的真 `fs.read` 必须真的落 L2**（:135）＋`user_rejected` | 把 `InAllowlist` 换成永远 true／摘掉 :315-317 那支／改掉那句措辞 | **动了**：MUT-V1（`paths.go:134` 硬置 true）→ 它红＋#6 红（共 15 枚红，含 13 枚别人的旧用例） | 否 |
| 2 | `TestPointerToMissingCopyFileSpeaks`（:146） | 形状 (b) 第一支：出"并不存在"（:163）＋**反向钉**：这条路径是在根内的，不许顺嘴说"不在你被授权的目录范围内"（:169） | 摘掉 `task.go:319-320` 那条 append | **动了**：MUT-V7（那一支换成 `_ = st`）→ 它红＋#6 红（reds=2） | 否 |
| 3 | `TestPointerToNonRegularPathSpeaks`（:178） | 形状 (b) 第二支：出"不是一般文件"（:191）＋指针仍在＋工具不得整体报错（:188） | 把 :321-322 的文案与 :320 合并成同一句（＝两形糊一形） | **动了**：MUT-V2（:322 文案改成"并不存在"）→ **恰 1 枚红**（就是它） | 否 |
| 4 | `TestHealthyPointerStaysSilent`（:203） | 六枚禁字（`注意：`/读不到/不存在/不是一般文件/未接线/不可找回，:218）一个都不许出现 ＋ **安静必须是真的**：同桥 `fs.read` 必须 `isError=false level=L0` 且全文读回（:225） | 让健康臂也喷说明（且喷的字落在禁表里） | **动了**：MUT-V8（`len(notes)==0` 那一支改成立即返回恒定说明）→ 它红＋#8 红（reds=2） | 否 |
| 5 | `TestUnwiredJudgeFailsClosedInReply`（:234） | 判定者没接线＝回执里必须出现"未接线"＋"按读不到处理"（:253）＋**不许替它下任何一个裁决**（:256）＋指针仍给 | 把 `task.go:307-309` 换成 `return ""`（＝fail-open） | **动了**：MUT-V4（`return ""`）→ **恰 1 枚红**（就是它） | 否 |
| 6 | `TestBothPointerDefectsAreReportedSeparately`（:267） | 既在根外又不存在那一发：两句各在场（:283/:286）＋`注意：` 计数**必须等于 2**（:289） | 任一支被另一支遮蔽（`switch`/`else if` 糊成一形）、或摘掉任一条 append | **动了两次**：MUT-V1、MUT-V7 各让它红一次 | 否 |
| 7 | `TestPointerNoticeKeepsTheD153StubShape`（:299） | 说明不得吃掉 D15(3) 形状：头 2000／尾 800 字节（:316/:319）、头尾内容逐字（:322）、`总长 N 字节`（:325）、`省略` 计数（:328）、`全文见 `（:331）、不得出现 `不可找回`（:334）、说明不得被塞进 `onUpdate` 进度回调（:337-341） | 动截断窗口／动那些数字／把说明搬进回调 | **动了**：MUT-V6（`headTokens()-1`）→ 它红（连同票 164 的两枚旧用例，reds=3），**而 #8 仍绿** | 否 |
| 8 | `TestHealthyReplyStillMatchesThePreFixTemplate`（:349） | 健康回执与"修码前那一整句"**逐字节**相等（:359-365） | 往模板里加任何不在禁表里的宿主话术 | **动了两次**：MUT-V3（模板里加"（副本已由宿主登记）"）→ **恰 1 枚红**（就是它）；MUT-V8 让它与 #4 同红 | 否 |

⇒ **结论（主问）：八枚里没有一枚是恒真的。** 每一枚我都指得出一发具体的产码改动，并且**八枚全部被我自己动红过至少一次**（MUT-V1…V8，共八发、覆盖八枚），不是复算被告那三枚。
⇒ 八枚的分布也不是"同一件事写八遍"：V2/V7 各只红一枚、V3 只红 #8、V4 只红 #5、V6 只红 #7 —— 这五发是"每枚各守一块地"的正证。

**但覆盖不全（本格第一枚新条件）**：`pointerNotice` 有**五**支会说话（`未接线`/`规范化没通过`/`不在授权范围`/`并不存在`/`不是一般文件`），八枚只钉了**四**支。`task.go:313-314` 那支（C26 `Canonicalize` 返错）——
`MUT-V5`：把它的说明文案整条摘掉（`grep -n 'MUT-V5'` 落地证明见 §4），`go test -count=1 ./internal/tools/` ＝ **rc=0、FAIL 0**（124 枚全绿）。
⇒ 全仓 124 枚用例里没有任何一枚摸过这一支（`grep -n '规范化\|没通过' internal/tools/*_test.go` 只命中 `fs.go` 的注释与无关用例）。这一支还会把 `err.Error()` **原样拼进模型看得见的那段文本**，是本格唯一一处"文案不由判据背书"的开口。**不算恒真判据，算判据少一枚。**

## 3. 安静正例那一问：它到底钉了哪一半

派单问："它是否**同时**钉了'该安静的不许多嘴'和'修码前后模板逐字节相同'？" ⇒ **不成立（作为单枚测试的说法）**：这两半分别在 #4 与 #8 两枚用例里，**没有一枚同时钉**。被告表 §4 把它写成"三重钉"归在"安静正例"一个标题下，读起来像一枚测试干两件事——实际是两枚。这不削弱覆盖面（两枚都在八枚名册里），但**被告的文字会让人高估单枚的强度**。

- **"逐字节等于修码前模板"这一句我是真比过的**：`git show 51c32ef9:internal/tools/task.go` 第 229 行（剥行首缩进后）与 `task_output_pointer_notice_test.go:360` 的 `want` 字面量 md5 同为 `8bd433707c1fd9039ac2103da08e7d22`（台件 `tpl-prefix.txt`／`tpl-test8.txt`）⇒ **基线不是编的**，那一枚钉的是真·修码前句子。
- **但"逐字节"的射程要划清**：#8 的 `want` 里 `head`/`tail`/`省略字符数` 是从 `out.Text` 自己身上切下来的（`headOf/tailOf`），**自回声**。证据＝MUT-V6（把头窗口缩短一枚 token）：#7 红、票 164 两枚红，而 **#8 仍绿**。⇒ 数字与几何是 #7 与票 164 的已勾用例在守，#8 守的是**固定文案＋路径＋token 计数**。两枚合起来才完整，任何一枚单独都不完整。
- **#4 单独也不完整**：MUT-V3 证明"用不在禁表里的字眼往健康回执上喷话术"这一假修法**#4 抓不到、只有 #8 抓得到**。⇒ **安静正例不是装饰**，它是两枚一组、彼此补洞；被告 §4 的结论方向对，措辞（"三重钉"挤在一个标题下）偏乐观。

**我自己造的那一发**（`.scratch/wisp/probes/174/v1/zz174v1_pos_windows_test.go`，经 `-overlay` 编进 `package tools`，`internal/tools/` 落盘零字节）：根内嵌套子目录、真实普通文件、30000 字节长输出，两腿只差宿主登记路径的写法：

```
--- PASS: TestProbe174V1HealthyPointerStaysSilentAndReadable/canonical-recorded-path (0.01s)
--- PASS: TestProbe174V1HealthyPointerStaysSilentAndReadable/raw-recorded-path (0.01s)
174-v1 positive leg …: bytes=30000 noticeCount=0 fs.read level=L0
```
命令：`go test -count=1 -overlay=.scratch/wisp/probes/174/v1/overlay.json -run 'TestProbe174V1' -v ./internal/tools/`
⇒ 两腿都：整段回执里 `注意：` 出现 **0 次**、头尾窗口未被扰动、同桥 `fs.read` `L0` 且 30000 字节全文读回。**修码的"安静"这一半我复现成立，且对 raw 写法不敏感。**

## 4. 我自己动的八枚变异（每枚先证落地再跑；还原只用 `git cat-file blob HEAD:<路径> > <路径>`）

| 变异 | 落地证明（`grep -n` 逐字） | 读数 | 还原 |
|---|---|---|---|
| MUT-V1 `InAllowlist` 永远 true | `internal/tools/paths.go:134:	return true // MUT-V1: the authorization judge says yes to everything` | `rc=1`，顶层 FAIL **15** 枚：本格 #1＋#6，另 13 枚是既有 C26 用例（`TestEmptyAllowlistAuthorizesNothing`／`TestTicket107*`／`TestBridgeRefusesARealJunctionOnTheReadRoute` 等）；**#2/#3/#4 仍绿**＝存在性两支与安静支确由别的fact驱动 | `git status --porcelain -- internal/tools/` ＝ **0 行** |
| MUT-V2 两形文案合并（:322 用"并不存在"） | `internal/tools/task.go:322:		notes = append(notes, "注意：这条路径现在读不到，宿主登记的那份副本文件并不存在") // MUT-V2: both (b) arms share one sentence` | `rc=1`，FAIL **恰 1** 枚＝#3 `TestPointerToNonRegularPathSpeaks`；#2/#6 仍绿 | 0 行 |
| MUT-V3 健康臂也被喷不在禁表里的话术 | `internal/tools/task.go:245:			"%s\n[…输出已落文件：省略 %d 字符，总长 %d 字节 / 约 %d token%s（副本已由宿主登记），全文见 %s…]\n%s", // MUT-V3: unconditional reassurance on every stub`（`grep -c '副本已由宿主登记'`=1） | `rc=1`，FAIL **恰 1** 枚＝#8；#4/#7 仍绿 | 0 行 |
| MUT-V4 未接线＝fail-open | `internal/tools/task.go:308:		return "" // MUT-V4: no judge wired => stay silent (fail-OPEN)` | `rc=1`，FAIL **恰 1** 枚＝#5 | 0 行 |
| MUT-V5 摘掉 `Canonicalize` 出错那一支的文案 | `internal/tools/task.go:313:	case err != nil: // MUT-V5: notice text for this arm deleted` | `rc=0`，FAIL **0**，124 枚全绿 ⇒ **零覆盖**（§2 末条件） | 0 行 |
| MUT-V6 头窗口短一枚 token | `internal/tools/task.go:227:	head := takeHeadTokens(full, t.d.headTokens()-1) // MUT-V6: head window one token short` | `rc=1`，FAIL **3** 枚＝#7＋票 164 的 `TestTruncationShapeIsTheD15Triple`／`TestTruncationBudgetsAreRead`；**#8 仍绿**（自回声，见 §3） | 0 行 |
| MUT-V7 摘掉"并不存在"那一支 | `internal/tools/task.go:320:		_ = st // MUT-V7: the missing-copy-file notice is deleted` | `rc=1`，FAIL **2** 枚＝#2＋#6 | 0 行 |
| MUT-V8 健康返回换成恒定说明（fail-open 的镜像：喷） | `internal/tools/task.go:325:		return "；注意：这条路径现在读不到（MUT-V8：健康臂也喷说明）"` | `rc=1`，FAIL **2** 枚＝#4＋#8 | 0 行 |

**末次现场**（全部还原之后）：`grep -c 'MUT-' internal/tools/task.go`＝**0**、`grep -c 'MUT-' internal/tools/paths.go`＝**0**、
`git status --porcelain -- internal/tools/`＝**0 行**、`go test -count=1 -v ./internal/tools/`＝`ok 15.873s`（PASS 124／FAIL 0／SKIP 0）。

## 5. 两形独立性那一问（派单 §1.3）

- **(b) 那形有独立用例，不是与 (a) 共用**：#2（不存在）与 #3（目录）是两枚函数、两句不同文案，而且各自都把根**配好**（`NewPathCanonicalizer([]string{dir}, nil)`），使授权那一支不可能替它响——这是"两形里只钉一形"那种病的正解。
- **`ArtifactPath` 指向一枚目录 vs 指向一条不存在路径有没有被同一枚断言糊过去**：**没有**，且不是我读出来的判断——MUT-V2 把两支文案合成一句后**恰 1 枚红**（#3），MUT-V7 摘掉存在性一支后**2 枚红**（#2＋#6）。两形各自可红＝两形各自有牙。
- 外加 #6 用 `strings.Count(announceOf(...), noticeLead) != 2` 钉住"两形俱备时不得互相遮蔽"——这一枚在派单口径上是**第三枚独立事实**（形状独立性），不是 (a)/(b) 的重复。

## 6. 次生面：这段新增文字有没有把 C25／R4 那一族弄脏（派单 §1.4）

**裁：没变。**凭据是两腿：

1. **码级（只读，一字节未动）**：`internal/tools/bridge.go:552` 的 `mark()` 第一句就是
   `if b.prov == nil || dec.TaskID == "" || !risk.IsSensitiveSource(dec.Tool) { return }`（`grep -n -B3 -A14 'func (b \*Bridge) mark' internal/tools/bridge.go`），
   而 `internal/risk/provenance.go:96-99` 的 `sensitiveSourceTools` 名单（八枚：`fs.read`/`search.content`/`clipboard.read`/`system.get`/`web.fetch`/`doc.read`/`screen.capture`/`asr.transcript`）**不含 `task.output`**
   ⇒ **新增文案根本到不了 `Mark`**，片段索引里的字符没有增减，"外来内容→R4 命中"那一族**多不到也少不到一次**。
   唯一一支会把不由常量决定的文本拼进回执的是 `task.go:314`（`err.Error()`，即 §2 的零覆盖支）——它同样过不了 `mark()` 的门，所以不构成 C25 面；它构成的是"文案未被告判据背书"那一枚条件。
2. **零字节面**：`git diff --stat 51c32ef9 HEAD -- internal/tools/bridge.go internal/risk/ internal/agent/ cmd/ docs/specs/` ＝ **空**。
   （同一段区间 `docs/PLAN.md` 有 `+1` 行，来路是编排者自己的 `e6e3fb33 ledger(A349)`，不是被告：`git log --oneline 51c32ef9..HEAD -- docs/PLAN.md` 只回这一枚。）
- 被告写面复核（越界＝零）：`git show --name-only` 五枚提交只碰 `internal/tools/task.go`（产码两枚）＋`internal/tools/task_output_pointer_notice_test.go`（判据一枚）＋票面与 `probes/174/r1` 与它自己的证据件（交付两枚）。**没有第三枚文件。**
- ⚠ 我没有重跑 `./internal/risk/` 那套（派单口径是逐包，且 `internal/risk/**` 零字节）⇒ 这一问的"没变"是**码级＋零字节**两腿，**不是**跑出来的 risk 套件读数。若编排者要跑套件级背书，那是一枚新的量，不是本格的缺陷。
- 按派单：`probes/177/c1/**` 与 `docs/evidence/s1/177-c25-r4-path-exemption-c1.md` 我**没看不评**（只在此登记"本程未观察该写面"）。

## 7. 门禁五枚（我自己重跑，未引它的数）

| 门禁 | 我的现量 | 命令 |
|---|---|---|
| `go test -count=1 ./internal/tools/` | `ok 15.873s`；PASS **124**／含子测试 RUN **173**／FAIL **0**／SKIP **0**（逐包，未跑全仓 `./...`） | 见 §4 末次现场 |
| `sh scripts/d22scan.sh` | **rc=0**，`d22scan: clean - no D22 ban violations`；`bans #1-5 internal/=207`、`ban #7 internal/tools/=20`、`ban #8 internal/=426`、`cmd/=45`、`frontend/=85`、`design/=39` | 台件 `gate-d22scan.txt` |
| `bash tools/d22scan/runtests.sh -C tools/d22scan ./...` | **rc=0**；`runtests.sh: OK - packages=[./...] top-level: PASS=34 FAIL=0 SKIP=0, === RUN=76, '[no tests to run]'=0`；名册两向 `comm -3`（工作树 34 枚 vs `51c32ef9` 34 枚）**左右皆空** | 台件 `gate-runtests.txt`、`d22-roster-now.txt`／`d22-roster-anchor.txt` |
| `sh .scratch/wisp/probes/154/gate-clauses.sh` | **rc=0**；`腿数＝14 声明与实测不符＝0`；`腿数断言：名册=14 声明=14 记账=14 缺腿=0 空头声明=0`；`基线过期枚数＝0`；`聚合退码＝0` | 台件 `gate-clauses.txt` |
| `gofumpt --version` ＋ `-l` | `v0.12.0 (go1.27.1)`（`$(go env GOPATH)/bin`）；`-l internal/tools/task.go internal/tools/task_output_pointer_notice_test.go .scratch/wisp/probes/174/v1/zz174v1_pos_windows_test.go` ＝ **空** | 见上 |

`probes/161/r6/flip-declaration.sh` **未跑**（派单禁；那批日志起手就带 ` M`，不是本程弄脏的）。

## 8. 裁语

**AC#2 这一格＝〔成立·带条件〕（维持 21:4x 的翻勾，不降级）。**

维持的理由（正面）：两形各自有牙（§5）、fail-closed 那一支只被一枚钉住但钉得住（MUT-V4）、安静正例不是装饰而是 #4＋#8 两枚一组（§3）、八枚**无一恒真**且八枚全被我亲手红过（§2）、我自己的健康正例两腿复现安静＋L0 全文读回（§3）、次生面 C25/R4 **没变**（§6）、门禁五枚全绿、写面零越界。

带的条件（本格新增两枚，都在 AC#2 射程内、都不必推翻翻勾）：
- **条件①（判据少一枚）**：`task.go:313-314` 的 `Canonicalize` 出错那一支**全仓零覆盖**（MUT-V5 rc=0），且它是唯一把 `err.Error()` 原样拼进模型可见文本的开口。补法＝一枚正向用例（造一发规范化必失败的路径形状，钉那句说明＋钉"不许顺嘴下别的裁决"）。**归口＝票 174 本体的缺口，不是 AC#2b/#2c。**
- **条件②（措辞比机器严，方向是安全的那一侧）**：形状 (a) 的文案说 `fs.read 会被拒；要用户先把所属目录加进 [fs] allowed_dirs 才读得回来`。我实测：**接一枚会答"允许"的门，同一条桥对根外路径的 `fs.read` 直接读回全文**——
  `174-v1 answered-card leg: out-of-root fs.read isError=false level=L2 bytes=4096`（台件 `zz174v1_pos_windows_test.go` 第二枚函数；码级旁证＝`internal/tools/fs.go:93-105` 的 `open()` 只做 `Canonicalize`，执行时**不再查** `InAllowlist`，根外的门就是那张 L2 卡）。
  ⇒ 那句话在**今天的默认形状**（审批通道未接线＝恒拒，`gate.go:135-144` 逐字"尚未接入（票 21），已拒绝执行"）下**为真**，在"接了卡且用户批了"的宿主下**偏严**：那条路是**要一张卡**，不是**走不通**。本格的标准是"指针不许许诺一条读不回来的路"，反向的"把走得通的说成走不通"是同一根轴上的小失真 ⇒ **留给编排者定文字**（要改就改文案，不改判据方向），**本程一字节未改产码**。
- **不计入本格**：端到端只会说"未接线"（＝AC#2b，未做）、`spill.go` 那枚桩仍静默（＝AC#2c，未做）。派单明写这两格不算 AC#2 的缺陷，本表照办。

`next=` 编排者：① 要么把条件① 立成票 174 里一枚新格（判据补一支，写码位），要么明写"接受该支不覆盖"并具名登记；② 条件② 只是一句措辞，随 AC#2b 接线那一程一起定文字最省（那一程本来就要端到端重跑回执）；③ 本格不降级、也不需回归票 177（C25/R4 没变，§6 两腿凭据）。

## 9. 被拒／没成功的调用（取数前 vs 取数后）

**取数之后没有一次被拒或失败的写调用。**取数过程中的不成功三处，全部不影响任何被引用的读数：
1. `awk '/sensitiveSources = /,/}/' internal/risk/provenance.go` 模式没匹配到 ⇒ **空输出**，随后改用 `sed -n '80,100p'` 取到名单原文（§6 引的是后者）。
2. 两次 `grep` 用 `|| ls` 兜底（文件名猜错：`internal/tools/fsread.go` 不存在）⇒ 未产生任何被引用的数。
3. 多次 `Edit` 工具回"文件已被上一枚命令改过"提示（因为我刚用 `git cat-file` 还原过它）⇒ **写照样成功**，落地证明都在 §4 表里逐字贴了 `grep -n`。
（被告表 §8 那两处自报不实我另在 §10 裁，不混进本程的调用账。）

## 10. 有没有跑过删除命令／凭据值

- **零删除命令**：全程无 `rm`/`del`/`Remove-Item`/`git clean`/`git restore`；"还原"一律是 `git cat-file blob HEAD:<路径> > <路径>`（写操作）。临时件**只建不删**（`probes/174/v1/` 现存 13 枚：基线／八枚变异读数／两台件／名册两份／门禁三份）。仓内没建 worktree、没 checkout 任何分支、没 `switch`／`merge`／`rebase`／`reset`／`stash`。
- **凭据值零抄录**：本文件与台件里出现的只有 `t.TempDir()` 派生路径与模板字面量；无任何 API 密钥／token／DPAPI Blob／配置内容，变量名也只出现 `[fs] allowed_dirs` 这一枚配置**键名**（派单要求的就是它）。
- **只 commit、绝不 push**；提交用显式 pathspec 一步式；未用 `git add -A`／`.`／`commit -a`；未 `--amend`（更正走追加）。现场那些不属于本程的脏件（`.gitignore`、`probes/152/my152.py`、`probes/161/r6/logs/flip-*`、`design/**` 的未提交删除）我一律**不提交、不还原、不评论**。

## 11. 伪授权两栏（各带出处）

**(a) 有没有"为了让哪一发变绿而放宽"？**
- 本程＝**零**：没改任何判据、任何断言、任何产码（§4 末次现场 0 行）。我造的 `v1AllowGate`（恒批）只活在我的台件里，用来问一句机器事实，**没有**参与任何 AC#2 判据，也没有被拿来说"这一格过了"。
- 被告＝我没抓到放宽的形状：五枚提交只碰两枚文件（§6 写面复核）；票 164 已勾的两枚用例文件不在其中；`internal/risk/**`、`bridge.go`、`cmd/**`、`docs/specs/**` 零字节（§6 第 2 腿）。**但两处"比声明弱"要记**：① §2 末的零覆盖支（被告表把"两形＋fail-closed"说成全覆盖，未提第五支）；② §3 的"三重钉"挤在一个标题下、且 #8 的"逐字节"射程被写宽（MUT-V6 它仍绿）。都不是放宽，是**描述偏满**。
- 被告自报的两处不实我复核：
  - **`Q-60` 未答＝假**。台账 `docs/reports/pending-and-issues.md:8057`（`A349`，09-27 21:3x）逐字＝owner "都按推荐"批完四枚、`Q-60` 走**丙**；`8061` 行另写 `Q-60` 丙＝正在跑（`174-r1`）。⇒ 证据件 §12 那句"`Q-60` owner 未答"在它自己的交付时刻就已过期；它已按规矩**追加**更正（`ad363d58`，票面第 48 行），原句留在证据件里不抹——处置正确，我照登。
  - **§8 第 4 条前半句"` M bridge.go` 一度入列"＝无法判定**（现场瞬时态，无可复算凭据，被告自己也认了）。但它**承重的那半句我复算成立**：`git diff --stat 51c32ef9 HEAD -- internal/tools/bridge.go` ＝ 空 ⇒ 桥那一枚确实一字节未动。

**(b) 有没有"用 mock 代替真的"来假报完成？**
- 被告＝**没有**：判据用的是真 `Bridge`＋真 `NewPathCanonicalizer`＋真 `fs.read`＋真落盘文件（`task_output_pointer_notice_test.go:49-89` 那三枚 helper）；`NoGate{}` 是仓里现成的**恒拒**形状（`internal/tools/gate.go:135-144`，逐字"L2 审批通道尚未接入（票 21），已拒绝执行"），**不是**伪造的批准；#1 的 L2＋`user_rejected` 读数是那条路自己拒的（:135/:137），#4 的 `L0` 全文读回是真的（:225）。我自己的台件同样全真（同一批 helper）。
- 一条**不成立的授权**照登：AC 框一枚都没勾（票 174 第 31 行那个 `[x]` 是**编排者 21:4x 翻的**，不是实现程勾的），本表也**不勾任何框**——本表只把裁语落在 `docs/evidence/s1/`，翻勾与否归编排者。
