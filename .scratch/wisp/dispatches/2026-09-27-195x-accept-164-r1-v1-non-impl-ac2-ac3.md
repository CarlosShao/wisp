# 派单 164-v1（**非实现者验收位，只裁不改**）＝裁票 164 的 **AC#2＋AC#3** 两格：`164-r1` 交了三枚提交（`bde17aca`／`24d1c09b`／`c84beb37`），它说"改前那一发做不到"、说"截断带的指针能找回全文"、还说"你的判断没被推翻"——**这三句都要由没写过这行码的人重量一遍**，尤其是一枚**永远会过的"改前"用例**能不能算判据。

- 派单时刻：09-27 19:5x｜锚点＝**你自己 step 0 现量**（别抄我的号；我这份派单写于 HEAD=`22313e33` 之上，但你起手仍要自己 `git rev-parse HEAD`）。
- 票面＝`.scratch/wisp/issues/164-background-jobs-cannot-be-read-fix-the-roster-then-add-the-output-leg.md`（**AC#2／AC#3 两格连行号＋票面那五条"编排者定案"抄进你的表**，定案①②③④⑤ 是判据的一部分，不是背景）。
- 实现者的表＝`docs/evidence/s1/164-task-output-impl-r1-ac2-ac3.md`（224 行）——**它是被告的陈述，不是证据**。
- ⚠ **预算硬顶 ≤70 次工具调用**；**每裁完一格 commit 一次**；到点停手→写 `next=`→把手上未提交增量用显式 pathspec 提掉。**不许勾任何 AC 框**（勾归我）。

## 0. 起手五件
`date`｜`git rev-parse --abbrev-ref HEAD`（**必须 `dev`**，不是＝停手上报，一个字不写；全程禁 `switch`／`checkout <分支>`／`merge`／`rebase`／`reset`／`stash`／`worktree`／`clean`）｜`git rev-parse HEAD`｜`git status --porcelain -- internal/ cmd/ docs/evidence/s1/`（**必须空**，非空报回别动）｜基线：`go test -count=1 ./internal/tools/ ./internal/agent/`（**逐包，禁全仓 `./...`**）四数＋名册口径（顶层 PASS vs 含子测试 RUN，别混着报）。

## 1. 四格主裁（每格给「成立／成立·带条件／不成立」＋可复制命令）
- **AC#2 那格先判"判据本身"**：⚠ **这是我这一单最想让你打的一枚**——`internal/tools/task_output_ac2_before_test.go` 是一枚**改完之后仍然会过**的用例（它自己搭注册表、不放 `task.output`，于是永远能拿到"未知工具"）。请裁：**这算不算"未修码读数"的判据**？如果算，它的**证伪条件**是什么（换句话说：产码怎么改会让它红？答不出来＝它是一枚恒真）。⇒ 若你判它恒真，**给出替代形状**（我倾向：改前那一发应由**归档锚 `f206e9f` 上的真树**产出——`git grep`/`git show` 那形，或在 `.scratch/wisp/probes/164/v1/**` 落一份 `f206e9f` 的仓外副本跑同一发请求；两形选一形并说为什么）。
- **AC#3 正向那一半**："全文见某路径"＋头尾摘要＋总长，**逐字节读得回来**吗？它声称 `TestLongOutputPointerRecoversEveryByte` 断言"桩里那条路径读回来的全文与实际全文逐字节相同、且宣布的总长是真全文的长"。**复算它**（读测试码，别读注释）。
- **AC#3 反向那一半（这一发必须在未修码上响）**：它自报两枚变异——摘掉 `，全文见 %s` ⇒ 那枚读回用例红；`refTaskSpillHeadTokens` 500→400 ⇒ 全库只 `TestTruncationShapeIsTheD15Triple` 红。**自己重跑其中至少一枚**，⚠ **变异必须先证落地**（改完先 `grep -n` 那一行贴出来，再跑测试，别拿"我改了"当"它生效了"）。**还原只允许这一形**：`git cat-file blob HEAD:<那枚路径> > <那枚路径>`（**禁** `checkout`/`reset`/`stash`），还原后 `git status --porcelain -- internal/tools/` 必须空并把这一行贴出来。
- **能力类 AC 必问生产调用者**：`task.output` 现在**注册在 `cmd/wisp/run.go`（生产），但名册零写者**（`RunAsync` 仍零调用点——自己现量一次）。⇒ 裁："读一个后台任务的输出"这件事**在真机上今天能不能做到**？如果不能，AC#3 该记〔成立〕还是〔成立·带条件〕还是〔不成立〕？⚠ **不许用"判据已过"直接顶掉这一问**——本仓有前例：安全性只因"线没接"才成立，后来被独立验收打成完整绕过（`Q-49`）。

## 2. 三枚我点名的攻击（各自单独一节，答"可接受／不可接受＋理由"）
1. **同源拷贝**：`internal/tools/task.go:280-307` 的 `takeHeadTokens`/`takeTailTokens` 是 `internal/agent/spill.go` 里 `takeTokens`/`takeTokensLast` 的**本地拷贝**（实现程自陈"故意"）。**逐字节比两族实现**（含 rune 边界处理），裁：语义有没有已经不同？两枚拷贝各归谁、以后谁改谁？⚠ 本仓为"同源拷贝逐枚没归属"返工过（票 141 那起）。
2. **它的理由里有一枚半伪约束**：注释说"不导出方法给 `agent` 是因为 `gate-clauses.sh` 的 G3 腿"。**现读那条腿的 pattern**（`:366-368`，我这份派单写的时候它是 `^func \(l \*Loop\) [A-Z][A-Za-z0-9]*\(.*taskID`，行号你自己量）⇒ 导出**自由函数**（`agent.TakeTokens(s, budget)`）撞不撞那条腿？如果不撞，**"为了不响 G3"就不能当"必须拷贝"的理由**——裁这句注释是不是**给一个设计选择编了一个仪器理由**（那在本仓比"拷贝"本身更糟：它会误导下一程）。
3. **戳的那一枚（不归本票，但你要判真伪）**：我本轮现读 `bridge.go:551-553`（`mark` 在 `!risk.IsSensitiveSource(dec.Tool)` 时直接 return）与 `internal/risk/provenance.go:94-99`（8 个名字的名册，`task.output` 不在内）⇒ 我的读法是"`task.output` 把外部内容读回上下文时不打 C25 污染源标记"。**复算我这两处读数**（⚠ 包括量我有没有看漏：`Origin` 字段、`res.Origin` 是否由产码自己填、还有别的盖戳路径吗）。**若成立，它归已立的票 175，你不要动 `internal/risk/**`、也不要把它算进本票 AC#3 的账**——只在报告里写"这枚是不是真破口、最坏后果什么形状"。

## 3. 门禁（你自己重跑一遍，别引用它的数）
`go test -count=1 ./internal/tools/ ./internal/agent/`｜`scripts/d22scan.sh`（预期 rc=0）｜`bash tools/d22scan/runtests.sh -C tools/d22scan ./...`（基线自己现量；四数之外名册两向 `comm` 比差集）｜`sh .scratch/wisp/probes/154/gate-clauses.sh`（**我这轮 19:40:17 现量 rc=0、`腿数＝14 声明与实测不符＝0`、G5/G6/G7 一枚没动**——你若跑出不同的数，**以你的为准并报回**，那是我给实现程的"没造出未成对收尾"这个结论被推翻）｜`gofumpt --version` 现跑＋对你名下的 `.go` 文件 `-l`（gofumpt 在 `$(go env GOPATH)/bin`，不在 PATH 就先 `export PATH="$PATH:$(go env GOPATH)/bin"`）。⚠ **`probes/161/r6/flip-declaration.sh` 默认别跑**（它会弄脏跟踪的 `logs/flip-*.txt`，共享树里那是别人的现场）；非跑不可就在报告里列出被改脏的路径。

## 4. 写面（超出即越权）
`.scratch/wisp/probes/164/v1/**`（台件，**不进主树**）｜`docs/evidence/s1/164-task-output-accept-r1.md`（新建，你的表）｜票面 164 的 Progress log（**只追加、不勾任何框、留到最后一步**；写之前 `git status --porcelain -- 那枚票面` 必须为空）。
**禁改**：`internal/**`、`cmd/**`（第 2 节那枚变异除外，且必须按还原协议收尾）、`tools/d22scan/**`、`allowlist.txt`、`.github/**`、`docs/PLAN.md`、`docs/specs/**`、`docs/reports/**`（台账与停车点归我）、`frontend/**`、`design/**`、别人的票面与证据件、`probes/**`（只读）。⚠ 现场躺着 `design/**` 里 owner 未提交的删除、`.gitignore`（` M`）、`probes/152/my152.py`（` M`）、`probes/161/r6/logs/flip-*.txt`（` M`）⇒ **不提交、不还原、不补完、不评论**；**`git add -A`／`git add .`／`git commit -a` 一律禁**；提交一步式 `git commit -q -F - -- <显式路径>`，heredoc **必须加引号**（`<<'EOF'`，正文里的反引号在壳里会被真执行）。**只 commit，绝不 push。**

## 5. 交件报告（顺序固定）
1. step-0 五件（分支名逐次）；2. **本程没测什么**；3. **AC#2 那格判据本身**（恒真否＋证伪条件＋替代形状）；4. AC#3 正向读回复算；5. AC#3 反向变异（**落地证明＋未修不响／修了响两向**＋还原后那行空 status）；6. **生产调用者那一问的裁定**（成立／带条件／不成立＋为什么）；7. 三枚攻击各自一节；8. 门禁四数＋名册差集＋`gofumpt --version`；9. 被拒／没成功的调用（发生在取数之前还是之后）；10. 有没有跑过删除命令；11. 伪授权两栏（各带出处）；12. 凭据值零抄录；13. `next=`。
⚠ **每个"几枚"旁边附可复制命令**；**先测→再写→再提交**（表里不许预先引用还没跑出来的读数）。我这段话里每条前提（含"改前那一发是未知工具""rc=0""那两枚变异会红""8 个名字里没它"）**都是未验证断言**：不符就报回并继续做做得动的部分，**不许为了对我那句话去改判据或改测试**。
