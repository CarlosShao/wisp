# 183-r1 — 先立在 HEAD 上就红的包级裁判，再逐族现量：F1（前程推荐的挖跨度）翻不动裁判，F2（按声明串自身的窗口值扣）翻得动，落的是 F2

- 程：`183-r1`（写码腿）｜派单 `.scratch/wisp/dispatches/2026-09-28-123x-impl-183-r1-referee-first-then-pick-the-family-that-actually-flips-it.md`
- 票面：`.scratch/wisp/issues/183-the-host-minted-pointer-exemption-does-not-hold-on-the-real-cli-so-r4-refuses-the-models-own-reread-and-the-spilled-output-still-cannot-be-read-back.md`（Progress log 只追加，**AC 框一枚没勾**）
- 前程表：`docs/evidence/s1/183-exemption-why-not-live-a1.md`
- 台件与读数：`.scratch/wisp/probes/183/r1/**`（`probes/183/a1/**` 一字未动，只读抄）

---

## 1. step-0 五件（起手 12:34:07 现跑）

| 件 | 读数 |
|---|---|
| `date "+%Y-%m-%d %H:%M:%S"` | `2026-09-28 12:34:07` |
| `git rev-parse --abbrev-ref HEAD` | `dev` |
| `git rev-parse --short HEAD` | **`b277e1e1`** ≠ 派单起手锚 `58124185` ⇒ 登记差异并按 HEAD 做。尺：`git merge-base --is-ancestor 58124185 HEAD` = yes；`git log --oneline 58124185..HEAD -- internal/risk/` = **空**（中间只多编排者那枚 A363 台账＋立票 184＋派单本身，没碰 `internal/risk`）⇒ 我的码基线与 `58124185` 同一份 |
| `git status --porcelain -- internal/ cmd/` | **空**（起手与交件各复量一次，见 §13） |
| 两支冻结文字 md5（起手） | `provenance.go:468-473` = `858e45116383caa3e7c1dd4b0924fad1`（与派单基线一致）；`taintmatch.go:11-15` = `5680ddd18e2d2ec2a85e485b54f4c12e`（一致） |
| 门禁基线 | ⚠ **偏离派单**：起手我只复量了 md5 与 `internal/ cmd/` 空，**没有先跑** d22scan／gate-clauses 的基线；两枚门禁改取在交件侧（§11），比名册时对照的是派单 §0.5 那枚在册读数（`ban #8 internal/ 431`／红腿名册只有 `G6neg`） |

## 2. AC#0 裁判判据在未修码上的红读数（逐字）

件：`internal/risk/pointer_183_test.go`（尺：`grep -n "^func Test" internal/risk/pointer_183_test.go`）。形状＝派单 §2 指定的那一形：**同一枚 mark，宿主路径只出现一次，正文别处含一条与路径同款的 ≥8 片段**；断言＝`Inspect("task-1","fs.read",{path})` 不命中 R4。尺（可重跑）：

```
go test -count=1 -run 'Pointer183' -v ./internal/risk/        # 未修码＝HEAD 的 internal/risk
```

逐字（`probes/183/r1/logs/referee-on-HEAD.txt`）：

```
--- FAIL: TestPointer183RefereeTwinFragmentOutsideSpanStaysExempted (0.01s)
    pointer_183_test.go:156: 183 referee: the host's own pointer must stay re-readable even when the body spells one of its windows elsewhere, got hit taint-hit channel=unclassified.exfil-scan source=fs.read task.output fragment_len=8
--- PASS: TestPointer183RefereeControlBodyWithoutTwinIsClean (0.01s)
--- FAIL: TestPointer183PathAppearingTwiceInBodyStaysExempted (0.00s)
--- PASS: TestPointer183BodySecretOutsideSpanStillHits (0.00s)
--- FAIL: TestPointer183ForeignMarkCarryingSamePathStillHits (0.00s)
    pointer_183_test.go:252: AC#3: hit must be attributed to the foreign mark, got "fs.read"
--- PASS: TestPointer183AlmostPathInSameHostMarkStillHits (0.01s)
```

- 裁判**红在 HEAD**，且归因就是票 183 那一发（`source=fs.read task.output`＝宿主自己那枚 mark，`fragment_len=8`）。
- **归因控制腿** `…ControlBodyWithoutTwinIsClean` 在未修码上就绿 ⇒ 红不是"fs.read/path 被扫"造成的，是"正文别处的同款片段"造成的（＝M-D1 的反向，做成了常驻腿）。
- 第二形（同一条路径出现两次，`runeIndexOf` 只回第一处）也红，按派单要求**单独一枚**，不当裁判。
- `…ForeignMarkCarryingSamePathStillHits` 在未修码上**也红**，但红在归因：命中来自宿主那枚 mark（因为豁免本身就漏了），不是外来那枚。修法之后它绿且归因换成 `web.fetch`＝"没洗戳"的正证（见 §4）。

## 3. 逐族现量：翻没翻裁判

尺（两族都只用 `-overlay`，跟踪件零改动；台件 `probes/183/r1/mut-f1`、`mut-f2`＋`overlay-f1.json`／`overlay-f2.json`）：

```
go test -count=1 -overlay=.scratch/wisp/probes/183/r1/overlay-f1.json -run 'Pointer183|ShapeA|177' -v ./internal/risk/
go test -count=1 -overlay=.scratch/wisp/probes/183/r1/overlay-f2.json -run 'Pointer183|ShapeA|177' -v ./internal/risk/
```

| 族 | 是什么 | 裁判翻没翻 | 读数 |
|---|---|---|---|
| **F1**＝`183-a1` §6 推荐：`provenance.go` 把豁免跨度从喂给 `newFragmentIndex` 的文本里物理挖掉（`rp[:lo]+rp[hi:]`） | 只作用于"这一处位置" | **没翻**（裁判仍红、第二形仍红） | `logs/f1-referee.txt`：`--- FAIL: TestPointer183RefereeTwinFragmentOutsideSpanStaysExempted`、`--- FAIL: TestPointer183PathAppearingTwiceInBodyStaysExempted`、`--- FAIL: TestPointer183ForeignMarkCarryingSamePathStillHits`，八枚 ShapeA 全绿（⇒ 不是把别的判据改坏了） |
| **F2**＝把"这条被声明路径**自身的窗口值**"从这枚 mark 的证据里扣掉（不分位置） | 位置＋值，仍只在这枚 mark 内 | **翻了**（六枚全绿，八枚 ShapeA 一枚没红） | `logs/f2-referee.txt`：`ok github.com/CarlosShao/wisp/internal/risk` |
| **F3**＝读码后我另找的形 | — | 没有更优的形 | 三条路都被现有契约堵着：① 参数侧按值放行＝票 183 AC#5① 明禁；② per-scope 路径名册＝`taintmatch.go` 那段与 `177-c1 §2`/`A353` 明令不批；③ 让 `contains()` 只在"跨度外的正文"里复核——那枚 `-output-` 本来就在跨度外，复核照样命中＝等于 F1。 |

**对编排者预测的裁定**：派单 §2 〔编排者预测·待推翻或坐实〕说 F1 翻不动真机那一发——**我的读数是坐实，不是推翻**（F1 三枚红、八枚 ShapeA 绿）。前程 `183-a1` 判对了根因（(d) 位置 vs 字符串集合），但它推荐的落地点治不到它自己那一发，本表 §3 的 F1 行就是它的反证。

## 4. 落地的族与落地形状（为什么不是 F2 的索引侧那版）

落的是 **F2**，但**不是**我用 overlay 量过的那版实现（索引侧 map），而是它的**等价复核侧形态**：

- `internal/risk/taintmatch.go`（**动了它，理由如下**）：`fragmentIndex` 新增一个**被读**的字段 `declaredPath []rune`；`contains()` 在算哈希之前先问"这枚候选窗口是否逐字落在这条被声明的路径里"，是则从**这一枚 mark** 的证据里跳过；新增 `attachDeclaredPath`（mark 仍在本函数手里、尚未发布给 scope 时调用，类型文档的"建完即不可变"对读者仍然成立）与 `spellsDeclaredPath`（**逐 rune 比较，不查哈希** ⇒ 64 位碰撞不可能多扣一枚窗口）。
- `internal/risk/provenance.go`：`MarkWithHostPath` 里 span 定位成功后把归一化的声明串挂到这枚索引上；文档块**只追加**一段"交付语义更正"（§5）。
- 为什么选复核侧：① 语义与索引侧那版等价（命中要成立，候选窗口必须在哈希集里；两版都让"声明串自己的窗口"不可能成为证据），而票 183 的 AC#3/AC#4 那两枚反向腿在两种形态下都仍然命中；② 索引侧那版要给每枚声明过的 mark 常驻一张 map（D32 索引内存预算那一面），复核侧只存一条路径字符串；③ 复核侧把"扣"发生在**逐字比较**那一面，绕开"按哈希集合减法＝碰撞会变成漏判"这一支；④ 纯插入，跟踪件**删除列 0**（§11 的 numstat）。
- 与 AC#5① "参数侧按值放行"的区别（写死在注释里，也写在这里）：不比对名册、不比对参数、不跨 mark；豁免集合从**这枚 mark 自己声明的那条串**导出，随 `Scope.Close` 一起消失。裁判腿的兄弟腿（同一条路径在**外来** mark 里、改一个字符的近邻路径、正文里那条与路径无关的秘密）全绿＝边界还在。

**票 183 AC#3 的 `marks=2` 那一格（派单点名"不许假设"）**：现量在 §6 的 CLI 读数里——`fs.read` 读回来的正文随后也被盖戳（`source=fs.read <产物路径>`），那枚 mark **没有声明**，所以它的 `-output-` 窗口全部照常入索引；第二发续读就是被这枚 mark 判中的＝**没有洗戳**的直接证据（同形常驻腿 `TestPointer183ForeignMarkCarryingSamePathStillHits` 归因 `web.fetch` 也绿）。桥级那对腿（真桥具、非包级 stub）在落地码上仍然成对绿：`go test -count=1 ./internal/tools/` = `ok 15.339s`，含 `TestHostMintedPointerRereadStaysCleanOnRealBridge175r2` 与 `TestForeignMentionOfHostPathStillHitsR4OnRealBridge175r2`。

## 5. 批准面复核（§4）

- 那 25 行文档块（`provenance.go:482-506`）：**原句一字未动**，机检证据是它仍然逐字节等于派单基线——`sed -n '482,506p' internal/risk/provenance.go | md5sum` = `f89e891e5eee3f3ea2b4f89d921072c4`（起手/交件两向同值），我的更正段追加在 `:506` 之后、`func` 行之前。尺：`git diff -- internal/risk/provenance.go | grep -aE "^-" | grep -av "^---"` = **空**（无任何被删／被改写的原句）。
- 两支冻结文字：`provenance.go:468-473`、`taintmatch.go:11-15` 两向 md5 见 §1 与 §11，**逐字节不变**。
- 追加段落里已具名标〔orchestrator-approved strengthening／出处 ledger A363〕，并写明那句 `only the normalized window that hostPath occupies inside THIS mark is kept out of the fragment index` 现在只是豁免的**位置性那一半**、另一半是值性；同时写了"为什么这不等于复活被否决的逐 token 追踪"（不沿 token 传播污点、不引入按 scope 的名册、只忽略逐字拼出这条声明串的窗口）。
- 没动：`docs/PLAN.md`、`docs/specs/**`、`internal/memory/schema.go`／`models.go`、`thresholds.go`／golden／审批超时常量／`allowlist.txt`、`internal/agent/**`、`internal/tools/**`、`cmd/**`、`frontend/**`、`design/**`、别人的票面与证据件、`probes/183/a1/**`、票 184 面。

## 6. 票 183 AC#2 在真机 CLI 上的端到端复测（背景正文一字未换）

台件：`probes/183/r1/zz183r1_e2e_test.go` ＝ `probes/183/a1/zz183a1_e2e_test.go` 的副本，**只**改两处：① 载荷钉回 179-r2 的原版背景正文 `strings.Repeat("WISP179R2-background-output-line.", 700)[:20000]`（它的 `-output-` 就是路径同款片段；a1 那份盘上是 M-D1 的纯 CJK 载荷，派单 §3 明令"换正文换来的绿不算数"）；② 读数目录 `probes/183/a1/logs` → `probes/183/r1/logs`（`runtime.Caller` 那条护栏照 a1 版原样保留）。断言、名册比对、分岔清点一枚没改。

```
PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" \
  go test -count=1 -v -overlay=.scratch/wisp/probes/183/r1/overlay-cli.json -run 'TestRereadHostPointerOnCLISeam183a1' ./cmd/wisp/
```

两形（派单 §6）：
- **不带** PATH 前缀：`exit status 0xc0000135`、`FAIL github.com/CarlosShao/wisp/cmd/wisp 0.052s`（票 98 那枚加载期坑，`logs/cli-no-prefix.txt`）。
- **带** PATH 前缀：`logs/cli-reread-landed.txt`，跑 60.11s，**测试仍然 FAIL**，但 FAIL 的位置换了：

```
51:[audit] tools: call kind=success task=9d3b802c-… tool=task.output risk=L0 decision=allow outcome=success rules_hit=[] in_allowlist_scope=true
55:[audit] tools: call kind=success task=9d3b802c-… tool=fs.read   risk=L0 decision=allow outcome=success rules_hit=[] in_allowlist_scope=true
60:[audit] tools: call kind=refused task=9d3b802c-… tool=fs.read   risk=L2 decision=reject outcome=error rules_hit=[R4] in_allowlist_scope=true reason="R4: 包含来自 fs.read C:\Users\swq\AppData\Local\Temp\…\artifacts\tool-output-agent-task-d9f7e6be-….txt 的内容"
2:run A 退出 0 / run B 退出 1 / 确认卡 0 -> 1
```

- **票 183 的那一发已经翻过来了**：模型照宿主指针发的**第一发** `fs.read` 从 179-r2 的 `kind=refused rules_hit=[R4] source=task.output / 读回 0 字节` 变成 `kind=success risk=L0 rules_hit=[]`，`logs/e2e-readings.txt` 第 22 行起那枚 tool 行带着**整份 20000 字节**正文与宿主为它新落的产物指针（`总长 20000 字节 … 全文见 …tool-output-call-183a1-full.txt`）。归因不再是 `task.output`。
- **还差第二发**（`max_bytes=10000` 那发，tool 行 **0 字节**，sha256 `e3b0c442…`＝空）：判中它的是 `fs.read` **读回来那 20000 字节的正文**所盖的戳——那枚 mark 没有声明、正文里的 `-output-` 是真·被污内容，按 AC#3/177 W-2 的边界**本就必须命中**。所以这不再是"豁免没生效"，而是"宿主自己的文件命名约定（每条产物路径都含 `-output-`）与任意一段被污正文撞车"的新形：任何一次读回来的正文只要含 `-output-`，模型就再也无法用路径寻址任何产物。**这一形要动的是桥（`internal/tools/bridge.go` 的 `b.mark(…, hostPaths.get())` 对 `fs.read` 从不填 box）或 D15 复读再落盘那一层**——两枚都在我的写面之外（派单 §5：判"必须动桥"＝停手上报）。⇒ 我**没**动、**没**自己修，具名报在 §12 `next=`。
- 顺带一格（正交、非本票阻断者）：`e2e-readings.txt` 里 `task.output` 的回执仍带票 174 AC#2b 那句"路径授权判定者未接线（fail-closed：C26 没接进来…）"。

## 7. AC#4 / 多判残留 / AC#6

- **AC#4 第三形**：`TestPointer183AlmostPathInSameHostMarkStillHits`（同一条路径改**一个** rune、且它就写在**同一枚宿主 mark** 的正文里）⇒ 命中、归因 `fs.read`（宿主那枚 mark 自己的非豁免窗口）。未修码与落地码**两向都绿**（它不是 red-then-green，它是防放宽的常驻钉：任何把豁免从"这条串的拼写"扩到"目录／前缀"的改法都会把它打红——共享前缀窗口被扣，改动位邻域那几枚窗口留着的还是它们没被扣、命中还在）。票 177 的 `TestShapeAW3NeverMintedSiblingUnderSameDirStillHits` 在落地码上仍然绿（§3 F2 行）。
- **残留钉**（派单 §3 要的"最坏会多出一张什么卡"，钉语义不钉枚数）：`TestPointer183WorstCaseOfTheLandedExemptionIsPinned`。落的这族**不多判**，它的代价方向相反＝**少判一枚卡**：同一枚 mark 正文里凡逐字拼出声明串 ≥8 个连续 rune 的片段（现实形状＝盘符＋`AppData\Roaming`＋`artifacts` 那段目录前缀）不再是续读的证据。三段各钉死：代价侧（落地码干净、**未修码红**——HEAD 上这一发确实命中）、边界侧（与路径无关的秘密 `SUPERSECRET-183R1-BODY` 在同一枚 mark 里**仍然**命中 ⇒ 值规则退化成整枚豁免就红）、围栏侧（同一段目录拼写出现在**外来** mark 里仍然命中且归因外来 ⇒ 豁免离开这枚 mark 就红）。
- **AC#6 `窗口 0.0s`**：本票不修、任何审批超时常量没碰。我这发真机上下了 **1** 枚 L2 卡（`确认卡 0 -> 1`，console 行 `[确认 L2 fs.read] R4: 包含来自 fs.read …`）；**卡片的秒数我没从读数里逐字取**（见 §9 没测清单）。归口照派单＝票 162 那一族（审批通道在 CLI 形状上的窗口取值），不是本票的修法面。

## 8. 真机频率那道必答题（现量，不靠推理）

尺（HEAD 未修码，`-overlay` 把 `provenance.go`＋`taintmatch.go` 钉在 `git show HEAD:` 的副本上，所以答话的是**未修的生产规则本身**，不是我手推的谓词）：

```
go test -count=1 -overlay=.scratch/wisp/probes/183/r1/overlay-freq-head.json -run TestFreq183 -v ./internal/risk/
```

形状＝宿主 stub（`省略 17200 字符，总长 20000 字节，全文见 <artifacts 路径>`）＋**真实正文**当尾部，`MarkWithHostPath` 声明该路径，然后 `Inspect(fs.read,{path})`。读数（`logs/freq-head.txt`，逐字）：

```
REAL-BODY FREQUENCY on UNFIXED code: 3 of 10 bodies make the host's own pointer unreadable
  AGENTS.md                 hit=true  \swq\app swq\appd wq\appda q\appdat \appdata appdata\ ppdata\r pdata\ro data\roa ata\roam ta\roami a\roamin \roaming roaming\
  docs/PLAN.md              hit=false
  docs/reports/HANDOVER.md  hit=false
  internal/agent/loop.go    hit=true  （同一组 profile 前缀窗口）
  internal/tools/bridge.go  hit=false
  cmd/wisp/run.go           hit=false
  scripts/d22scan.sh        hit=false
  go.mod                    hit=false
  synthetic:dir-listing     hit=true  （同一组 profile 前缀窗口）
  synthetic:plain-prose     hit=false
```

**结论：不是"只夹具响"，真内容也会响**，10 份里 3 份响。而且响的片段**不是**夹具那枚 `-output-`，是**用户 profile 前缀**（`\appdata\`、`\roaming`、`\swq\app`）——AGENTS.md 与 `loop.go` 里逐字写着 `%APPDATA%\wisp`／`C:\Users\…\AppData\Roaming\wisp` 那一族路径，任何一次真实任务输出只要把用户的 AppData 路径打印出来，宿主自己写下的产物指针就再也读不回来。这一维把票 183 的严重性从"夹具运气"抬成"自用第一天就会撞上"（`c:\users` 那支旁证同向）。
⚠ 这份量的**不是**发生率分布：语料是仓内 8 枚可读文件＋2 枚我手写的合成正文，其中 2 枚 SPEC 文件我按错文件名取、读不到就没计入（`corpus docs/specs/SPEC-06-….md unreadable`、`SPEC-08-….md`）；真实任务输出（LLM 答复正文）里同款片段的出现率**本程零读数**。

## 9. 本程没测什么（§5）

1. 起手没跑门禁基线（d22scan／gate-clauses），只在交件侧各取一次；`ban #8 internal/` 从派单的 431 变 432 是我新增一枚 `*_test.go` 的**被扫文件数**，不是违例数（违例仍 clean）。
2. 落地后的真机 CLI **只跑了一发**（12:51，pre-commit），终态 commit 之后没再重跑；commit 不改工作树，且 §11 的 `internal/ cmd/` 空＋两向 md5 证明取数代码＝提交代码。
3. AC#6 那枚 L2 卡的**窗口秒数**没逐字取（只取了 `确认卡 0 -> 1` 与 console 那一行）。
4. 频率那格没量真实 LLM 输出语料、没量跨进程/重启后的名册与 mark 对应、没量第二枚 SPEC 文件（文件名取错）。
5. 第二发有界续读那一形（§6 后半）**没测修法**——它要动桥或 D15 再落盘层，两枚都在写面外，按派单 §5 停手上报。
6. 没跑整仓 `go test ./...`、没跑 SLO／golden 那套（票面 AC#8 只点名 `internal/risk`＋`internal/tools`＋`internal/agent`，我按那一格跑）。
7. F2 的两种实现（索引侧 map／复核侧逐字）只做了**语义等价推理＋同一套判据两向跑绿**，没做逐窗口集合相等性证明（那需要新写一枚比较器，不在写面）。

## 10. 被拒／没成功的调用（§7）

- 没有一次工具调用被权限系统拒。
- 失败／作废逐条：① 第一次 staging 命令变量没加引号，路径里的空格把 `mkdir`／重定向打散（`ambiguous redirect`），同一条命令顺带撞出无关的 vet 噪声输出（`probes160v1ac1b`／`closebyid.go`，不是我建的、不是我跑的），改用引号后重跑成功；那次误触还在 `D:/work/workspace/` 下多了一枚空目录 `projects`（**只建不删**，留在盘上）。② `pointer_183_test.go` 第一次跑是**编译失败**（`undefined: runeDifferences`），补上函数才取到红读数——那枚 rc=1 不是判据红。③ 频率格两枚 SPEC 语料文件名猜错、读不到（§8 已具名）。④ 落地后 CLI 那发**测试 FAIL**（第二发被拒那一形），不是工具失败，是产品事实，逐字留在 §6。

## 11. 门禁与 git 纪律（§8 删除命令／两向 md5／终态三把尺）

| 尺 | 读数 |
|---|---|
| `go test -count=1 ./internal/risk/`（**单包跑**，A359 那枚假红） | `ok github.com/CarlosShao/wisp/internal/risk 4.078s`（含票 177 八枚 ShapeA＋本程七枚新腿全绿） |
| `./internal/tools/` · `./internal/agent/` | `ok 15.339s` · `ok 1.778s` |
| `sh scripts/d22scan.sh` | **rc=0**，`ban #8 internal/ examined 432 Go files`（基线 431＝我新增的那枚 `_test.go`），`clean - no D22 ban violations` |
| `bash .scratch/wisp/probes/154/gate-clauses.sh` | **rc=1**，BAD 行 2 枚，红腿名册 `grep "BAD" … \| grep -oE "腿=G[0-9a-z]+" \| sort -u` ＝ **只有 `腿=G6neg`** ⇒ 与派单 §0.5 在册名册（票 178 那一枚）**逐字相同，没新增红腿**。取数时刻 `2026-09-28 12:54:26`（pre-commit，代码与提交同一份）；读数 `logs/gate-clauses.txt` |
| `"$(go env GOPATH)/bin/gofumpt" -l` 三枚改动件 | 无输出（已格式化，v0.12.0） |
| 冻结文字 md5 **两向** | 起手 `858e45116383caa3e7c1dd4b0924fad1`／`5680ddd18e2d2ec2a85e485b54f4c12e`；交件（§11 末列，commit 之后复量）**同值** |
| 有没有跑过删除命令 | **没有**。全程没跑 `rm`／`git clean`／`restore`／`checkout .`／`reset`；作废的 overlay 台件与失败读数（含 `mut-f1`／`mut-f2` 两份变异副本、`cli-no-prefix.txt`）全部留在盘上 |

## 12. next=

1. **要编排者裁一枚新的缺陷面**（本程最要紧的一笔）：真机上第一发续读已经 `kind=success`／20000 字节回全，第二发（有界重读同一条路径）被 **`fs.read` 结果那枚没有声明的 mark** 判 R4。落地的 F2 **按 AC#3 的边界不能**替它豁免（外来／未声明的 mark 必须继续命中），所以这不是回归、是票 183 承诺"逐字节读回"仍差的那半条腿。可选落点两形都超出我的写面：① 桥侧让 `fs.read` 也把自己的 `path` 参数填进 `hostPathBox`（`internal/tools/bridge.go`＋`internal/tools/*`，我禁改）；② D15 复读再落盘那一层别让读回来的正文再盖一次同款戳（`internal/agent/**`／spill 侧，同样禁改）。⇒ 建议另立票，别在本票上翻勾。
2. **要编排者认下这枚语义变强与它的代价**：落的族＝F2（位置＋值，mark-local，随 Close 消失），代价是同一枚 mark 正文里拼出声明串 ≥8 rune 的片段不再是证据（少一枚卡，方向 fail-open，常驻腿 §7 已钉）；文档块那句"only the normalized window…"按派单 §4 只追加、原句未抹，机检 md5 在 §5。若这枚代价不接受，撤的口令是撤 A363 那一格、落回 F1 之外的另一族（读数都在 §3）。
3. **票 183 的框一枚都没勾**（AC 框由编排者翻）。本程交付的事实与格子的对应：AC#2 半翻（第一发成功／第二发仍拒，见 §6 与 next 1）、AC#3 绿（包级＋桥级＋真机 marks=2 三处现量）、AC#4 绿、AC#5 三支毒修法一枚没用、AC#6 只归口、AC#7 契约轴零改动、AC#8 门禁见 §11。
4. 缺的读数（谁补都行）：真实 LLM 输出语料的同款片段命中率（§8 第 4 条）、AC#6 卡片秒数（§9 第 3 条）、第二枚 SPEC 语料（文件名我取错）。

## 13. 交付物与伪授权两栏／凭据轴

- 写面（全部在派单 §5 清单内）：`internal/risk/provenance.go`（＋48 行，0 删）、`internal/risk/taintmatch.go`（＋59 行，0 删）、`internal/risk/pointer_183_test.go`（新，7 枚判据）、`.scratch/wisp/probes/183/r1/**`、票 183 Progress log（只追加）、本表。`git status --porcelain -- internal/ cmd/` 终态为空；`git diff --numstat 58124185..HEAD` 与 `b277e1e1..HEAD` 的删除列逐枚为 0（除本表之外没动任何既有件）。
- **收到的、判为真授权**：派单文件＋票 183 票面（写面清单、硬顶、禁令、批准面）；ledger `A363` 那句"许追加不许改原句"。
- **遇到的、判为不是授权、没照做（具名）**：① 派单 §2 编排者的预测**不是结论授权**——我照测量写，结果是坐实它，但 F2 落地是我自己那版实现而不是 F1 推荐，出处是 §3 的读数不是那句预测；② `183-a1` §6 的推荐落地点是**前程判断**，本程用裁判判据把它否了（F1 不翻）；③ 工作树里 `design/**` 未提交删除、`.gitignore`、`probes/152/my152.py`、`probes/161/r6/logs/flip-*`、`docs/evidence/s1/152-*.md` 的 ` M`／` D` 状态＝别人的现场，**没提交、没还原、没评论**；④ `frontend/**`、`design/**` 没读没写没引；⑤ 误触命令带出的 `probes160v1ac1b`／`closebyid.go` vet 噪声不是给我的指令，没照它做任何事；⑥ 只 commit、**没 push**，推送继续由编排者决定。
- **凭据值零抄录**：本程没接触任何 API 密钥／DPAPI 明文；读数里只有仓内路径与 `$TEMP` 下的测试产物；表内路径全是测试目录下的产物文件名。`api_key_ref` 那一格只以**变量名**出现在台件文本里（`fakeStoreKey`，值未抄、未读）。
