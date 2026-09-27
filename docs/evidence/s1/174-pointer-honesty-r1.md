# 174-r1 证据件：`task.output` 指针诚实化（丙形，两形）——回执在模型看得见的那段文本里说实话

- 派单＝`.scratch/wisp/dispatches/2026-09-27-211x-impl-174-r1-pointer-honesty-both-shapes-no-allowlist-widening.md`（提交 `51c32ef9`）
- 量取时刻：2026-09-27 21:16 → 22:0x +08｜分支 `dev`
- 产码提交＝`b5a07961`＋`f55ddd3d`（文案微调）｜台件提交＝`9ae0be78`（新判据）＋本文件同枚的 probes 提交
- ⚠ **测量层＝接缝级**：真 `tools.Bridge` ＋ 真 `task.output` ＋ 真 C26 `PathCanonicalizer` ＋ 真 `fs.read` ＋
  真 `agent.Spiller` 落盘，跑在 `go test` 里。**端到端 `wisp run` 今天仍然没跑过**（原因见 §2，是本机环境事实，不是本程偷懒）。

## 0. 起手五件（逐字，全是本程现量）

| 件 | 读数 |
|---|---|
| `date` | `Sun Sep 27 21:16:31 CST 2026` |
| `git rev-parse --abbrev-ref HEAD` | `dev`（＝要求值，未停手） |
| `git rev-parse HEAD` | `51c32ef9fa4597ccde1a41b857c5df058dcf3fb0`（＝派单自身那枚，未抄别人的号） |
| `git status --porcelain -- internal/tools/ cmd/wisp/ .scratch/wisp/probes/174/` | **空** |
| 改前基线 `go test -count=1 ./internal/tools/` | `ok github.com/CarlosShao/wisp/internal/tools 15.557s`；顶层 PASS **116**／FAIL **0**／SKIP **0**／包 **1 枚 ok** |

- `internal/tools/bridge.go` 起手干净：`git status --porcelain -- internal/tools/bridge.go` ＝ **空**，且
  `git log --oneline -1 -- internal/tools/bridge.go` ＝ `06eb4efb 回退 175-r1 的产码修法（bridge.go 逐字节回到改前）＋正控用例转入台件留档`
  ⇒ 派单那句"刚被回退过"**成立**。本程对它**零字节**（见 §7）。
- 名册口径：改前 `internal/tools` 顶层用例名 **116** 枚（`git grep -h '^func Test' 51c32ef9 -- internal/tools/`）。

## 1. 本程没测什么（先说清射程，再看数）

1. **端到端没测**。起手 `go build ./cmd/wisp` 与 `go test ./cmd/wisp/` 我没重跑，但 174-c1 那程现量的是
   `wisp-probe.exe` rc=127、`go test ./cmd/wisp/` 退码 `0xc0000135`（缺 DLL）。**同一台宿主、同一枚 exe 形状** ⇒
   我按"这条腿今天依旧不通"处理，没有把"真机上回执长这样"写进任何结论。要补需换宿主，属另一枚账。
2. **真宿主没接线路径判定者**：`cmd/wisp/run.go:362` 今天构造的是 `tools.TaskDeps{Roster: rt.tasks}`，
   **没有** `Paths` 字段（本程禁改 `cmd/**`）。所以端到端真跑起来，回执会是**未接线那一支的 fail-closed 文案**
   （"无法核实，按读不到处理"），**不是** §4 里那两句精确到"不在授权范围内／并不存在"的说明。
   这不是本格的缺陷（派单明写"拿不到判定者＝fail-closed 说清楚"），但**结题前必须补那一行接线**，见 §11 `next=`。
3. **`spill.go`（回路侧那枚桩）没动也没测**：本票 AC#2 原文把 `task.output` 与 `spill.go` 并列，派单把这一格射程
   写成 `task.output` 一支（写面清单里也没有 `internal/agent/**`）。⇒ **`agent.Spiller` 那枚桩今天依旧静默**，
   本格没修它，也没假装修了。
4. 没测并发下的判定者热换（`Paths` 在结构体构造时定格），没测 `[fs]` 之外的任何配置面，没测 256 KiB 帽与偏移那两枚
   （它们分别是票 164 的 canary 与 `Q-59`，本程一字节未动）。

## 2. 形 (a)：路径在授权根外 —— 两向逐字

**未修码（不响）**——`.scratch/wisp/probes/174/r1/` BEFORE 支，`allowed_dirs=[]`、真 spill：

```
[A-outside] isError=false truncated=true announcement="[…输出已落文件：省略 48833 字符，总长 51633 字节 / 约 12908 token，全文见 C:\Users\swq\AppData\Local\Temp\TestProbe174R1PointerHonestyBefore398114459\001\artifacts\tool-output-call-1.txt…]"
[A-outside] 提到读不到=false 提到不存在=false 提到不是一般文件=false 提到未接线=false 全文见=true 不可找回=false
[A-outside] the fact behind the words, fs.read on that path: isError=true level=L2 class="user_rejected" text="L2 审批通道尚未接入（票 21），已拒绝执行"
```

⇒ 派单那句"今天不响"**成立**，且 174-c1 的拒绝形状复现（`isError=true level=L2 class="user_rejected"`）。

**修完之后（响）**——同目录 AFTER 支（`TaskDeps` 带上同一枚判定者）：

```
[after A-outside] isError=false truncated=true announcement="[…输出已落文件：省略 48833 字符，总长 51633 字节 / 约 12908 token；注意：这条路径现在读不到，它不在你被授权的目录范围内，fs.read 会被拒；要用户先把所属目录加进 [fs] allowed_dirs 才读得回来，全文见 C:\Users\swq\AppData\Local\Temp\TestProbe174R1PointerHonestyAfter4174769956\001\artifacts\tool-output-call-1.txt…]"
[after A-outside] 提到读不到=true 提到不存在=false 提到不是一般文件=false 提到未接线=false 全文见=true 不可找回=false
[after A-outside] fs.read of that same path: isError=true level=L2 class="user_rejected"
```

⇒ 指针照旧在（`全文见=true`），说明在模型看得见的那段 `Result.Text` 里，且这段话的真伪由**同一条桥、同一枚判定者**
的真 `fs.read` 读数兜住（L2 拒）。判据：`TestPointerOutsideAuthorizedRootSpeaks`。

## 3. 形 (b)：路径不存在／是枚目录 —— 两向逐字

**未修码（两支都不响）**：

```
[B-ghost] path-said="[…，全文见 …\artifacts\tool-output-does-not-exist.txt…]"
[B-ghost] 提到读不到=false 提到不存在=false 提到不是一般文件=false 提到未接线=false 全文见=true 不可找回=false
[B-dir] path-said="[…，全文见 …\artifacts…]"
[B-dir] 提到读不到=false 提到不存在=false 提到不是一般文件=false 提到未接线=false 全文见=true 不可找回=false
```

⇒ 派单那句"那一支只看 `ArtifactPath != ""`、没有 stat"**成立**；不存在的那条与一枚目录都走"有指针"那一支。
（口径提醒：`不可找回=false` 是本程 `speakFlags` 自己打印的旗标；174-c1 那句"`不可找回` 出现 0 次"是 `grep -c` 口径，两把尺不同名同义。）

**修完之后（两支都响）**：

```
[after B-ghost] announcement="[…输出已落文件：省略 48833 字符，总长 51633 字节 / 约 12908 token；注意：这条路径现在读不到，宿主登记的那份副本文件并不存在，全文见 C:\...\artifacts\tool-output-does-not-exist.txt…]"
[after B-ghost] 提到读不到=true 提到不存在=true 提到不是一般文件=false 提到未接线=false 全文见=true
[after B-dir] announcement="[…输出已落文件：省略 48833 字符，总长 51633 字节 / 约 12908 token；注意：这条路径现在读不到，它存在但不是一般文件（是目录或别的形状），全文见 C:\...\artifacts…]"
[after B-dir] 提到读不到=true 提到不存在=false 提到不是一般文件=true 提到未接线=false 全文见=true
```

⇒ 判据：`TestPointerToMissingCopyFileSpeaks`、`TestPointerToNonRegularPathSpeaks`。
两形是**两枚独立事实**、不是 if/else 的一支：`TestBothPointerDefectsAreReportedSeparately` 造一枚"既在根外又不存在"的路径，
回执里 `注意：` 出现 **2 次**（`strings.Count(...) != 2` 即红）。

**fail-closed 那一支（判定者没接线）**：

```
[after-nojudge 复跑 BEFORE 支，即今天 run.go 的接线形状]
[A-outside] …；注意：路径授权判定者未接线（fail-closed：C26 没接进来，这条路径是否还读得回来无法核实，按读不到处理），全文见 …
[POS-ctrl]  同一枚真能读回的文件也吃到这句未接线（旗标：提到未接线=true）
```

⇒ 形状参考同文件 `Execute` 里那句 `任务名册未接线（fail-closed：…）`，没有静默按"能读"处理。判据：`TestUnwiredJudgeFailsClosedInReply`。

## 4. 安静正例（反向那一发，不是恒真）

在授权根内、且是一枚真普通文件——**修完之后仍须逐字闭嘴**：

```
[after POS-ctrl] announcement="[…输出已落文件：省略 48833 字符，总长 51633 字节 / 约 12908 token，全文见 C:\...\artifacts\inside-root-real-file.txt…]"
[after POS-ctrl] 提到读不到=false 提到不存在=false 提到不是一般文件=false 提到未接线=false 全文见=true 不可找回=false
[after POS-ctrl] fs.read of the control path: isError=false level=L0 bytes=51633 (whole file is 51633)
```

三重钉，缺一重就挡不住"把回执写成永远提示"这类假修法：
1. `TestHealthyPointerStaysSilent`：`注意：`／`读不到`／`不存在`／`不是一般文件`／`未接线`／`不可找回` 六枚旗标全不存在，
   **并且**同一条桥的真 `fs.read` 必须 `isError=false level=L0` 且 51633→20000 字节全文读回（安静得说是真的）。
2. `TestHealthyReplyStillMatchesThePreFixTemplate`：把这枚健康回执与本票 §2/§3 里**修码前那一句模板**逐字节比对，
   差一个字节即红。⇒ "没顺手放宽成噪音"是字节级命题，不是感觉。
3. `TestPointerNoticeKeepsTheD153StubShape`：head 仍 2000 字节、tail 仍 800 字节、`省略`/`总长` 数字仍是全文口径、
   指针句 `全文见 ` 原样在场，且说明**没被塞进 `onUpdate` 进度回调**（回调里出现 `注意：` 即红）。

票 164 已勾判据的两枚文件（`task_output_leg_test.go`／`task_output_ac2_before_test.go`）**一字节未动**，
且它们里面那枚把 `ArtifactPath` 填成目录的用例现在会多出诚实说明——`headOf/tailOf` 的取法只按**首个** `…]\n`
切尾，所以我把说明插进括号行内部（`token` 与 `，全文见` 之间），**既有断言全部原样绿**：全仓名册只增不减（§6）。

## 5. 变异自证（每枚先证落地，再跑，再还原）

| 变异 | 落地证明（`grep -n`） | 结果 | 还原 |
|---|---|---|---|
| M1 整段说明摘掉：`notice := "" // MUT-1` | `internal/tools/task.go:243` | `TestPointerOutsideAuthorizedRootSpeaks`／`TestPointerToMissingCopyFileSpeaks`／`TestPointerToNonRegularPathSpeaks`／`TestUnwiredJudgeFailsClosedInReply`／`TestBothPointerDefectsAreReportedSeparately` **五枚红** | `git status --porcelain -- internal/tools/task.go` ＝ **0 行** |
| M2 授权判定换成死条件 `case canon == "": // MUT-2b` | `internal/tools/task.go:315` | `TestPointerOutsideAuthorizedRootSpeaks`／`TestBothPointerDefectsAreReportedSeparately` **红**；存在性两枚＋安静正例**仍绿**（＝它俩真由 `InAllowlist` 那支驱动） | **0 行** |
| M3 存在性 stat 那支摘掉 `statErr != nil && false // MUT-3` | `internal/tools/task.go:319` | `TestPointerToMissingCopyFileSpeaks` 红，**但红法是 panic**（`invalid memory address`，`task.go:321`）＝我这枚变异把 `st` 留在了 nil 分支上，**是变异的产物不是产码的缺陷**；方向仍对（那一支确实被执行到），但**别把它读成干净的断言红**——所以我补了 M3b | **0 行** |
| M3b 同一支的文案摘字（把两枚存在性说明换成 `（M3b）`） | `internal/tools/task.go:320`＋`:322` | `TestPointerToMissingCopyFileSpeaks`／`TestPointerToNonRegularPathSpeaks`／`TestBothPointerDefectsAreReportedSeparately` **三枚干净断言红**；授权腿 `TestPointerOutside…` 与安静正例 **仍绿** ⇒ 存在性两形真由 `os.Stat` 那支驱动，与授权支互不遮蔽 | **0 行** |
| M4（把说明改成恒挂）**本程没有效读数**：我那次 `sed` 把函数签名行改坏 ⇒ `FAIL [build failed]`，是变异脚本自己的缺陷，不是变异生效；预算到顶后我没重跑第二版。**如实记成"未做的变异"，不记成做过的**。 | `grep -n 'M4 恒挂'` 只命中被我写坏的那一行 | 无红无绿可用 | 还原后 `git status`＝**0 行**，随后全仓 `go test ./internal/tools/` ＝ `ok 17.197s` |

**恒挂形状真正的反证**（不是上面那枚变异）：修完的产码在 `TestHealthyPointerStaysSilent`
与 `TestHealthyReplyStillMatchesThePreFixTemplate` 上**是绿的**，而那两枚用例里写着
`strings.Contains(out.Text, "注意：")` 与逐字节模板比对——**任何"永远提示读不到"的实现都会把它们打红**。
⇒ 安静正例不是摆设；把 M4 补成一枚真变异，留给裁决者（一句话的 `sed`，`return ""` 那一行换成立即返回恒定说明）。

**一次真实红（写在本程，不抹）**：M1 之后我**漏还原过一次** `notice := ""`，随即被全仓跑拦下——
`TestLongOutputPointerRecoversEveryByte`＋`TestPointerPast256KiBIsNotFullyReadable`（票 164 已勾 AC#3）打出
`the pointer must be re-readable with fs.read: {Text:… IsError:true ErrorClass:"tool"}`，
原因是我的 `announceOf/stubAnnouncement` 取首个 `…]\n` 之后全部当尾，被残留说明污染。我**没有改判据去迁就**，
而是把 `task.go` 还原（`06eb4efb…` 不涉，本程只 `sed` 反向替换）、重跑全绿才继续。⇒ 这两枚 AC#3 用例确实在守"指针必须真能再读"。

**另一处我自己的错，同样不抹**：`TestBothPointerDefectsAreReportedSeparately` 第一次跑就红（`found 1` 不是 2 枚说明），
根因是我那枚"两形俱备"的台件把不存在的文件填在了**授权根之内**——那一发只该响存在性那一支，是台件错不是产码错。
改成 `elsewhere := tempCanonical(t)`（第二枚临时目录，落在根外）后两枚说明同时在场。

最后一次 `grep -c 'MUT-' internal/tools/task.go` ＝ **0**；`git status --porcelain -- internal/tools/task.go` ＝ **空**。

## 6. 门禁（改前／改后各一次）

| 门禁 | 改前 | 改后 |
|---|---|---|
| `go test -count=1 ./internal/tools/`（逐包，未跑全仓 `./...`） | `ok 15.557s`；PASS 116／FAIL 0／SKIP 0 | `ok 14.033s`；PASS **124**／FAIL **0**／SKIP **0** |
| `sh scripts/d22scan.sh` | rc=**0**，`d22scan: clean - no D22 ban violations`；`bans #1-5 internal/=207`、`ban #7 internal/tools/=20`、`ban #8 internal/=425` | rc=**0**，同样 clean；`bans #1-5 internal/=207`、`ban #7 internal/tools/=20`、`ban #8 internal/=**426**`（＋我这一枚新文件） |
| `bash tools/d22scan/runtests.sh -C tools/d22scan ./...` | 命令跑过，我那道 grep 过滤器把汇总行吞了 ⇒ **只看到 `1 ok`**，四数没抄全（记在这里，不假装抄到） | `runtests.sh: OK - packages=[./...] top-level: **PASS=34 FAIL=0 SKIP=0**, `=== RUN=76`, `[no tests to run]`=0`（该套件只测 `tools/d22scan` 自己，与本票 diff 无交集） |
| `sh .scratch/wisp/probes/154/gate-clauses.sh` | rc=**0**；`腿数＝14 声明与实测不符＝0`；`基线过期枚数＝0`；`聚合退码＝0` | rc=**0**；`腿数＝14 声明与实测不符＝0`；`腿数断言：名册=14 声明=14 记账=14 缺腿=0 空头声明=0`；`基线过期枚数＝0`；`聚合退码＝0` ⇒ **没造出"开了不关"的形状** |
| `gofumpt --version` | `v0.12.0 (go1.27.1)`（`$(go env GOPATH)/bin`） | `-l internal/tools/task.go internal/tools/task_output_pointer_notice_test.go .scratch/wisp/probes/174/r1/` ＝ **空** |
| 名册两向 `comm` 差集 | 116 枚（`51c32ef9`） | 124 枚（工作树）；`comm -3` **左栏为空＝一枚没掉**，右栏＝本票新增 8 枚（`TestPointerOutside…`／`TestPointerToMissing…`／`TestPointerToNonRegular…`／`TestHealthyPointerStaysSilent`／`TestHealthyReplyStill…Template`／`TestUnwiredJudgeFailsClosedInReply`／`TestBothPointerDefects…`／`TestPointerNoticeKeepsTheD153StubShape`） |

`probes/161/r6/flip-declaration.sh` **没跑**（派单点名禁止；那批日志起手已脏，现场 `git status` 复现：` M` 八枚）。

## 7. 我有没有动过"判定接口／C26 之外的路径判定"（单独一节）

- **零**。`internal/tools/task.go` 里路径相关的调用只有三处，全在 `pointerNotice` 内：
  `d.Paths.Canonicalize(raw)` → `d.Paths.InAllowlist(canon)` → `os.Stat(raw)`。前两项就是 C26 那枚
  `*tools.PathCanonicalizer`（与 `FSDeps.Paths`、桥的 `Options.Paths` 同一枚实例，接口零改动）。
- **`os.Stat` 不是路径判定**：它只回答"这个名字在不在、是不是一枚普通文件"，不做任何"能不能碰"的裁决；
  且它是**只读**——本程没新增任何 `os.Remove`/`RemoveAll`/截断/移动（`grep -n "os\.[A-Z]" internal/tools/task.go`
  只有 `os.Stat` 一枚命中，另一处是注释）。"能不能读"永远仍由 C19/C26 在 `fs.read` 那侧裁。
- 本文件**零** `filepath.Clean`／`filepath.Abs`（唯一字符串命中在 §1.2 那句禁止项的**转述注释**里，
  `d22scan` rc=0 复算通过）。
- **没动授权根**：`cmd/wisp/run.go` 的 `allowed` 构造（:327-331）一字节未改；`[fs] allowed_dirs` 默认值未改；
  `SPEC-03:35` 未改；`allowlist.txt`／`internal/risk/**`／`internal/agent/**`／`internal/tools/bridge.go` 全部零字节。
  台件里那枚"配好根"的腿用的是 `t.TempDir()` 派生目录，不改任何生产默认。
- 接口新增只有 `TaskDeps.Paths *PathCanonicalizer` 这一枚字段（在 `internal/tools/**` 内，派单明写"可以加"）。

## 8. 被拒／没成功的调用（取数前 vs 取数后）

**全部发生在取数之前**（起手洁净检查之后、基线四数之后），没有一次落在读数已抄进本表之后：
1. 三次 `search_replace`/编辑被工具层退回（old_string 不匹配或文件已被 gofumpt 改写过）⇒ 重读后改对，无副作用。
2. 一次 `grep` 因 `ugrep` 引擎参数不认（`-v` 位置）退码 2 ⇒ 换成普通 `grep` 重跑，**重跑的是我的统计口径，不是被测码**。
3. 一次 `Edit` 报"文件在会话里被别的程动过"（锚点位移：`e13196bd` 那枚 163-a1 提交把我的派单提交顶成了父提交）⇒ 重新 `Read` 后再改，未覆盖别人的东西。
4. 起手时 `git status --porcelain -- internal/tools/` 一度只列出 ` M internal/tools/bridge.go`——**那是别的程的活**，本程对它零字节；我按"它脏、我不动"处理（派单 §0 那句"确认它是干净的一枚文件"在**我的 scoped 读数**里为真：`git status --porcelain -- internal/tools/bridge.go` 起手为空）。

取数之后：无被拒调用、无失败工具调用。

## 9. 有没有跑过删除命令

**没有**。全程零 `rm`/`del`/`Remove-Item`；临时件只建不删；本程新建的只有
`.scratch/wisp/probes/174/r1/{zz174r1_before_windows_test.go,zz174r1_windows_test.go,logs/readings.txt}`、
`internal/tools/task_output_pointer_notice_test.go`、本文件。
变异自证的"还原"是用 `sed` 把同一枚字符串换回去（写操作，非删除），还原后 `git status --porcelain -- internal/tools/task.go` 为空。
仓内没建 worktree、没 checkout 任何分支。

## 10. 伪授权两栏

**(a) 为了让哪一发变绿而放宽过什么？＝没有。**
- 没动 SLO 阈值／golden／`thresholds.go`／票 164 的已勾判据文件（那两枚文件在 `git log -1 --stat` 里不出现于本程任何提交）。
- 唯一一次"判据变了"是 §5 那枚台件 bug（两形俱备那一发填错目录）——**改的是台件的fixture 输入，没改任何一条断言的方向**；
  改完之后断言更强（要求 2 枚说明，不是 1 枚）。
- M1 之后残留的那次红我是**还原产码**，不是把打红的两枚用例删掉或放宽。
- 文案是我写的（派单"措辞归实现程"），但每句都锚在一个可重量上："不在你被授权的目录范围内"↔`InAllowlist` 假；
  "并不存在"↔`os.Stat` 的 ENOENT；"不是一般文件"↔`!st.Mode().IsRegular()`；"未接线"↔`d.Paths == nil`。

**(b) 有没有"用 mock 代替真的"来假报完成？＝没有。**
判定者是 `NewPathCanonicalizer` 本尊；`fs.read`、`agent.Spiller`、`Bridge`、C19 判定全真；
台件里唯一非生产件是 `NoGate{}`（票 21 未接线的现成形状，恒拒），本程**没有伪造任何一次 Allow**，
L2 腿的读数就是它自己拒绝的读数（`class="user_rejected"`）。`TaskRoster` 是本票要测的东西本身、不是替身。
另外明写一条**不成立的授权**：AC 框一枚都没勾（`AC#2` 在原位未动），勾它要另一枚程（裁决者≠实现者）。

## 11. 凭据值

零抄录。本文件与台件日志中没有任何 API 密钥／token／DPAPI Blob／配置文件内容；
出现的路径全部是 `t.TempDir()` 派生（`C:\Users\swq\AppData\Local\Temp\TestProbe…\001\…`）或模板字面量。

## 12. next=

`next=` 编排者：本格（票 174 AC#2 丙形，两形＋fail-closed）已在 `internal/tools/task.go` 落地，
产码 `b5a07961`＋`f55ddd3d`、台件 `9ae0be78`＋probes/本表一枚；请派**非实现者**做缺口审计＋对抗验收，
并顺手核 §7 那三条"我只读了没裁决"。三件要摆的账：
1. **接线那一行属 `cmd/**`，本程禁改** ⇒ `cmd/wisp/run.go:362` 目前把 `TaskDeps{Roster: rt.tasks}` 构造得没有判定者，
   端到端回执只会说"无法核实"。要不要现在补（一行 `Paths: paths`，`run.go` 里那枚 `paths` 已在 :331 附近构造好），
   还是和 175/177 的接线票并成一枚——**这是编排决定，不是本格缺口**。
2. **`spill.go` 那枚桩（回路侧）今天仍静默**：AC#2 原文点名"`task.output`（与 `spill.go` 的桩）"，
   派单只给了前者。要不要为后者另立一格，请明写，别让"AC#2 过了"被读成"两处都过了"。
3. **丙形做完之后，票 174 剩下的是甲还是乙？＝两枚都不是本票能自证的，必须摆 owner。**
   甲（把 `dataDir` 加进 `[fs] allowed_dirs`，让默认安装的续读真能读）与乙（给 `fs.read` 加偏移，即 `Q-59`）
   都是**授权/契约面**，`Q-60` owner 未答、`Q-59` 未批，本程一字未碰。丙形只买到"回执不再撒谎"，
   **没有**买到"默认能续读"——AC#3（三条禁区自证）与 AC#4/AC#5 仍在原地等它们各自的批准。
