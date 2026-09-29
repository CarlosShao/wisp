# 226 — **答复"以后这类都别问"时，程序把整份内存里的配置快照写回 `config.toml`**：运行期间你手改的**任何其它键**会被悄悄还原，而且因为文件时间被"认领"，程序此后**永远看不见那次手改**

- Status: **产码已交（09-29 13:04，四枚 commit `78988901`／`361b314a`／`d90b9aec`／`c31239a9`），⛔ 六格今天一枚没勾**——非实现者对抗验收腿 **`226-v1`** 正在飞（表＝`docs/evidence/s1/226-config-write-no-clobber-v1.md`；实现腿自报 AC#1/2/3/4/6 已闭、AC#5 半格，**我一个字都不预认**，理由＝`SPEC-12 §4.3` 与本仓 `A428` 定死的口径"裁决腿的判语不自动等于翻勾、字面判据优先"）。选支＝票面候选**甲**（写前重读＋只替换那一枚键），实现腿给了三条现读理由、我在台账 `A432` 逐行复认了代码形状。**排程仍照原话：本票必须先于票 223**（223 接上轮询之后，本票那条"不认领"的路径才真被走到；今天它已由用例钉住）。
- ⚠ **我在读 `internal/config/writeguard.go:111` 时起了一枚连实现腿都没量过的疑问，已作为主攻点交给 `226-v1`（结论不由我猜）**：落盘走的是**整个结构体序列化**，所以**用户"删掉一行"（而不是改一个值）会不会被程序按内存里的值补回文件**——"删"也是一次手改，若会被补回＝**下面 AC#1 只钉住了"改值"这一支**。
  - ⚠⚠ **上面这两行是被 09-29 14:3x 那一发取代关系里的"原句"，我第一版编辑把它们整行删掉了、这一发逐字补回**（`git diff --numstat` 的删除列说真话；我自己写"原句不抹"却抹了＝这就是那枚"拿下一条目的标题当锚点却不复带它"的老坑换了个形态）。**保留在此不删，下面那四行才是现行口径。**
- Status: **终裁已下（09-29 14:2x，编排者）：六格里 4 格勾（AC#1／AC#2／AC#3／AC#4，两格带注）、AC#6 勾（我自己复跑到终态）、AC#5 不勾（字面只满足半边，残余转票 227／219／224）⇒ 本票不翻 `-done`。** 裁决表＝**`docs/evidence/s1/226-config-write-no-clobber-v2.md`（31,097 字节，六枚 commit `ee40bbd0`→`172715d6`）**；实现腿自报的"AC#1/2/3/4/6 已闭"我**没有一句预认**，全部换成非实现者读数＋我自己复跑的尺（见下面"收 `226-v2`"那节）。
- ⛔ **上面这行原来写的两件事都已作废，就地更正（原句不抹）**：**①表不是 `…-v1.md`——那枚文件从来没存在过**（`226-v1` 死在写表之前，见台账 `A436`；真正入库的是 **`-v2.md`**，`226-v2` §⑦ 最后一条顶出这一格，我 `ls` 复认）；**②"本票必须先于票 223"是排程陈述、不是盘上事实**——`226-v2` §⑥ 条 8 现读证实：**本票没有给 `CheckAndReload` 接生产轮询**（那是票 223 的射程），所以 AC#3 那条"下一次轮询要能看见"今天**只由测试里手动调 `CheckAndReload()` 钉住**。准确说法＝**"本票的写路径先于票 223 的接线"**。
- ⚠ **我在读 `internal/config/writeguard.go:111` 时起的那枚疑问（"用户删掉一行会不会被按内存值补回文件"）——今天有答案了，不是悬案**：**会被补回，而且补成 schema 默认值、不是内存那枚值**（`226-v2` §⑤，据 `226-v1` 的 `b-*`／`b2-*` 台件读数；map 条目那一支**不**补回）。⇒ **它推翻了我自己写下的 AC#1 覆盖面**（那一格只钉住了"改值"支），残余两格转成**票 227**（见下面"收 `226-v2`"第 3 节）。原句保留在此不抹。
- 来源：非实现者裁决腿 `201-v1`（`docs/evidence/s1/201-reply-listener-r2-verdict.md`，45,385 字节，提交 `21d9a310`）第 ② 条＋**编排者自己逐行读码复认**（台账 `A428`）。
- ⚠ **这不是安全漏洞，请照这三行读**：①**现象出现在哪**＝本机 `%APPDATA%\wisp\config.toml` 这个文件的内容；②**有没有本机被入侵的证据**＝**没有**，这是我们自己写文件的逻辑缺陷，不涉及任何外部攻击面；③**最坏后果是什么形状**＝**用户自己改的设置被静默还原成开机时的样子，且没有任何一句告诉他**（丢的是配置，不是数据盘里的文档）。

## 现量（起手逐条复算，别信这里的行号）

| 事实 | 读数 | 尺 |
|---|---|---|
| `always` 那一支走哪条写路径 | `cmd/wisp/approval_always.go:134` 一带 `s.rt.mgr.AddAllowedDir(dir)` | 现读 |
| 落盘时写的是**整份内存快照** | `internal/config/allowdirs.go` 的 `writeAllowedDirs`：`m.cur.FS.AllowedDirs = next` 之后 **`SaveFile(m.path, m.cur)`** ⇒ 交出去的是**整个 `*Config`**，不是那一节 | 现读 |
| `SaveFile` 从不回读文件 | `internal/config/loader.go` 里 `SaveFile` 逐字 `cp := deepCopyConfig(c)` → `cp.SchemaVersion = SchemaVersionCurrent` → `MarshalCanonical(cp)` → `atomicWrite(path, data)`——**没有任何一步把磁盘上的当前内容读进来合并** | 现读 |
| 写还把"这是我自己写的"这件事认领掉 | `writeAllowedDirs` 成功后调 `statOwnWrite()`，把新 `mtime+size` 记成"已见过" ⇒ **此后 `CheckAndReload` 认为文件没变过**（`allowdirs.go` 里那两句自陈逐字可见：`adopt the file's new stat so the program's own write is not read back as an unconfirmed hand edit`） | 现读 |
| ⚠ **同一枚形状不是新造的，是继承的** | `SetPermissionMode` 那条路早就这么写（裁决腿原话："`SetPermissionMode` 同形，属继承但**暴露面 1→2**"）——**新增的只是"聊天里点一下就能触发"这个入口** | 〔腿报＋编排者复认方向；枚数请续腿自己现跑〕 |

## 后果（为什么这是"owner 第一天就会撞上"的那一枚）

1. 他**开着程序**，手改 `config.toml` 里任何一个与审批无关的键（模型名、音量、快捷键……）；
2. 接着在对话里点一次**"以后这类都别问"**；
3. ⇒ **刚才那次手改被开机时的旧值覆盖**，而且**程序把新文件的 mtime 认领成"我写的"** ⇒ 既不会重读、也不会提示冲突 ⇒ **他只会看到"我改的东西怎么自己回去了"**，且**查不到是谁干的**。
4. 这条**正好打在票 219 与批准记录 `A424` 的安全轨上**：那里写着"长期＝往 `allowed_dirs` 加一行、**持久、可审计、可撤销**"——**今天的写路径不满足"可审计"**（它顺带改写了整份文件），也**不满足"可撤销"**（`SetAllowedDirs` 全仓**零调用者、零用例**，我现跑 `grep -rn "SetAllowedDirs" --include=*.go .` 只剩注释与定义行）。

## 判据（每格都要现跑读数；不许用"字段在场"充当"功能通了"）

- [x] **AC#1 手改不许被吞（核心正控）**：起进程 → **手改一枚与 `[fs]` 无关的键** → 触发一次 `always` 落盘 → **那枚手改必须仍在文件里**。⚠ 正控方向：**把"写前重读／合并"那一步拿掉，这发必须红**；如果它现在就是绿的，先证明它绿是因为"合并真做了"，而不是因为测试根本没改到第二个键。
  - **09-29 终裁＝成立（勾）**。作依据：`a1-cmdwisp-mutated.txt` 进程级接缝红 `TestAC1AlwaysBranchDoesNotRevertAHandEditedKey`（`cmd/wisp/always_write_no_clobber_226_test.go:87`）；`a1-config-mutated.txt` 包级红 `:74 got 56, want 60`；常驻正控 `TestAC1ControlSnapshotWriteIsWhatRevertsTheHandEdit` 反接即红（`a4-control-inverted.txt:111` "the control has lost its teeth"）⇒ **绿只能来自合并，这一格不是恒真**。**⛔ 带硬注**：字面判据只钉住"**改一枚键的值**"这一支——现跑 `internal/config/writeguard_226_test.go` 13 枚 `^func Test`（我自己 `grep -c`）里"删一行"的形状＝**0 枚**（`grep -ic 'delet|remov.*line|missing key'` 我自己跑，读 0），而实测删的行会被补回成 **schema 默认值**（`226-v2` §⑤）⇒ **本票第 3 条后果对"删"这一形今天仍未被接住，残余转票 227 AC#1。**
- [x] **AC#2 冲突必须响亮**：磁盘内容比内存快照新（或与之分叉）时，写路径只能二选一——**要么拒绝并说清"哪几枚键冲突、你要不要覆盖"**，**要么合并并报告实际写了哪几枚键**。⛔ **不许"写了才说"，也不许静默覆盖**（`AGENTS.md` §1.2／D37 错误模型那一族：错误要能分类、能传播、能显示成状态）。
  - **09-29 终裁＝成立（勾），选的是"合并并报告"那一支**。作依据：`a1` 拿掉合并后 AC#2 三枚用例同时红（`:127` 未报实写的键／`:154` 干净写未在 INFO 报键／`:170` 竟接受了不可读文件并覆盖）；码侧我**自己 `sed` 读到**（`internal/config/writeguard.go:181-186`）：分叉走 `slog.Warn`＋`kept_in_file_not_in_memory` 点名保留的手改、`wrote` 点名实改的键。**⚠ 带注（这句是我证据件里没有的）**：AC#2 那句"不白重排、不白丢注释"**只在 schema 已是当前版时成立**——`b4-probe-migration-inside-read.txt` Y2／`b6-probe-nothing-written-branch.txt` W2、W3 实测：**手删 `schema_version` 或塞入旧版文件时，即便走 `nothing written` 早返回支，`readConfigFile` 内部的迁移照样把文件重排成 canonical、注释丢、mtime 变、落 `config.toml.bak-1`**。属继承自 loader 的行为、非本票新造 ⇒ **收窄与常驻用例转票 227 AC#2**。
- [x] **AC#3 `statOwnWrite` 的认领范围要收窄**：只允许认领"**本次确实由这次写产生的那份内容**"；**不许把"我们写过一次"变成"这之后文件的任何变化都算我自写的"**。判据要能证：写完之后再手改一枚键 ⇒ 下一次轮询**必须**看见它（否则就是本票第 3 条后果没被修掉）。
  - **09-29 终裁＝成立（勾），且这一格是本票最硬的一格**。**两枚腿独立复现同一发红**：`226-v1` 的 `a3-config-unconditional-adopt.txt` 与我重派的 `226-v2` 亲跑（`226-v2` §③①，落盘件 `.scratch/wisp/probes/226/v2/rerun1-unconditional-adopt.txt`）——把"仅有分叉才不认领"改成无条件认领 ⇒ **稳定红 2 枚**（`writeguard_226_test.go:224` "adopted a stat for content this process never read; the hand edit is now invisible forever"、`:259` "the hand-added [fs] entry was never judged"），且**只红这两枚**⇒ **AC#3 不是装饰判据**。码侧我**自己 `sed -n '158,190p'` 读到**：`:166 diverged := diffKeyPaths(m.cur, base)` → `:168 if len(diverged) == 0 {` → `:177 m.statOwnWrite()` → `return nil`，分叉支不认领。`d-probe-d36-adopt-bypass.txt` S3/S4 另证**甲没有变成绕过 D36 的后门**（分叉外键不认领、下一轮读到并按 D36 rule 1 拒放宽）。⚠ **但"下一轮"今天在生产里不存在**（没人轮询＝票 223）——这一格的"下一次轮询必须看见"目前是**测试里手动调 `CheckAndReload()`** 钉住的。
- [x] **AC#4 同形的另一处一起处理**：`SetPermissionMode` 那条写路径要么**同批改**，要么**具名登记为已知残缺**（写清"哪枚文件哪一行为什么今天不动"），⛔ **不许只修新入口就宣称这一类问题已解决**。
  - **09-29 终裁＝成立（勾，同批改、未留残缺）**。作依据：`a2-config-permmode-mutated.txt:14` 拿掉 `permmode.go` 合并后红 `TestAC4PermissionModeWriteKeepsAHandEditedKey`（`:278 got 56, want 60`）；`c3-permmode-rollback-removed.txt:7` 拿掉 mode 写失败回滚后红 `TestGuardedWriteFailureRollsBackMemory`（`:387`）⇒ 两形都有牙。⚠ **既有安全行为没被削弱**：冻结件 `internal/perm/ticket90_persist_test.go` 一字未动、`TestTicket90PersistFailureKeepsMemory` 仍 PASS（`g0` 基线整包 `internal/perm ok`）。
- [ ] **AC#5 "可撤销"那一半要有真路径**：撤销一条长期规则必须**既能落库、又能被答复语法触发**（今天 `SetAllowedDirs` 零调用者零用例）。⚠ 这一格与**票 219 的撤销格同源**——**两边谁先做都行，但只做一边就是把 `A424` 那句"可撤销"写在纸上**（记 `deferred-work-must-be-registered` 那规）。
  - **09-29 终裁＝⛔ 不勾（字面要求两半，只交付一半）**。**成立的那半**：`SetAllowedDirs` 从此**有了它的第一枚用例**（`writeguard_226_test.go:313`、`:350`；改前全仓零调用者零用例，我 09-29 现跑复认"零生产调用者"仍成立、"零用例"已过期）⇒ 落库半边会红、非恒真。**不成立的那半**：答复语法（用户说什么话能触发一次撤销）今天**零入口**——`SetAllowedDirs` 生产调用者＝**0 枚**（我自己跑，见台账 `A437`）。⇒ **这一格按字面不翻，残余具名转两处**：①撤销动词属**票 219 的答复语法面**；②"会话档授权"属**票 224**（`A428` 已把票 201 AC#5 整格转过去）。**⚠ 按 `AGENTS.md` §1.1 与 `SPEC-12 §5`：这条推迟今天还没有落进登记表（实现腿按规矩没有代码写面去挂 `DEFERRED` 标记），登记责任由我接＝见 `A437` 的"悬空账"那一格。**
- [x] **AC#6 整包终态**：`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./cmd/wisp ./internal/... -count=1` 到终态＋**逐名比红名册**（历史在册的 `internal/ball` 1 例与 `internal/panel` 4 例属别人地界，不算本票新增、不许顺手修；`internal/risk` 的 `TestResolvePerCallBudget` 是**争用型假红**，安静复量为准，见台账 `A426`）。
  - **09-29 终裁＝成立（勾），依据取自我自己那一发、不是转述**：起手 `14:15:40`／终态 `14:17:15`，件＝`.scratch/wisp/probes/226/v2/orchestrator-full-suite.log`（全量落盘、未接 `| head`／`| tail`）＝**23 包 ok／3 包 FAIL 共 6 例**，逐名册：`TestC21TableColourRowsMatchTokensCSS`（ball）＋`TestApprovalCardViewJSONKeysMatchFrontendTypes`／`TestComposerContractTypesMatchFrontend`／`TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`／`TestC21DesignTokensFourwayAgree`（panel）＋`TestResolvePerCallBudget`（risk，1.73s）。**第六例归因＝争用**：编队清空后我安静跑 `go test ./internal/risk -count=3 -run 'TestResolvePerCallBudget'` ＝ **`ok 3.642s`** 转绿，`thresholds.go` 一字节未动（同 `A432` 那一族的第二例）。⇒ **零新增红成立**：与历史在册五例逐名一致，且与 `226-v1` 三发安静期读数（24 包 ok／FAIL 恰这 5 例）对得上。

## 落点候选（≥2 支，各挂代价；本票不替实现者选）

- **甲：写前重读＋只替换那一节**——最贴"D36 的三档生效级别"，代价是要处理"节内格式注释丢失"（`MarshalCanonical` 是结构体序列化，**注释本来就不保留**；要确认这点对用户是不是可接受的残缺，写成一句明确的文案）。
- **乙：写前比对 mtime/size，分叉就拒绝并提示**——最小改动、最不容易写错，代价是"用户必须重做一次答复"。
- **丙：整份合并（读盘＋内存三方合并）**——效果最好，代价最大，且**最容易造出新的静默行为**（⛔ 没有 AC#2 那条"报告实际写了哪几枚键"就别做这支）。

## 禁区

不动 `PLAN.md`／`docs/specs/**`／`thresholds.go`／golden／`allowlist.txt`；不动三枚冻结件（`internal/panel/tokens_fourway_test.go`／`l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go`）；⛔ **`AddAllowedDir` 现有的两条安全行为一条都不许削弱**：含 `..` 的条目**拒收**、写失败**回滚内存**（我读到 `allowdirs.go:60-77`，那是 fail-closed，本票只动"覆盖范围"这一条）；不新增 C17 方法名；新增导出名要先落 `A##`；`frontend/**`／`design/**` 零读零写零转述；⚠ **与票 201／222／223／224 同撞 `cmd/wisp`＋`internal/config` ⇒ 串行**，且**必须先于票 223 落地或与之同批**（否则 223 一接上轮询，本票的覆盖会被"读回旧值"再放大一次）。

## 收 `226-v2`（编排者终裁，09-29 14:2x，起手锚点 `c039bc4e`）

**1. 三把我自己的尺（先跑尺再写条目，这条规矩的出处＝`feedback-phantom-subagent-results`）**

| 尺 | 我的读数 |
|---|---|
| 六枚 commit 在不在 | `git log -1` 逐枚复跑：`ee40bbd0`(13:59 骨架)→`29daf745`(§1+§2)→`1d4f6ff0`(§3)→`721e3864`(§4+§5)→`a7d70ebf`(§6+§7)→`172715d6`(入库清单) ＝ **全部存在** |
| 表在不在 | `wc -c docs/evidence/s1/226-config-write-no-clobber-v2.md` ＝ **31,097 字节**（与它自报同值） |
| 有没有把源码留在突变中间态 | `git status --short -- internal/config cmd/wisp internal/tools internal/perm` ＝ **空**；`git diff --numstat c039bc4e..HEAD` ＝ **只有三枚文档件（表 +204／HANDOVER +19／台账 +13），零枚 `.go`** ⇒ **产码一件未动**，`writeguard.go` 逐字节复原（它的 `certutil` 前后同值我复认了"工作树干净"这一半） |

**2. 我自己复跑/复读、不依赖它表内自述的三处**：① `sed -n '158,190p' internal/config/writeguard.go` 亲眼读到 AC#3 的收窄形状（`:168` 那个 `if len(diverged) == 0` 与 `:177` 的 `m.statOwnWrite()`）；② `grep -c '^func Test' internal/config/writeguard_226_test.go` ＝ **13**、"删行"形状 ＝ **0 枚**（这是票 227 AC#1 的分母）；③ `sed -n '393,410p'` 读到 `TestGuardedWriteDoesNotRewriteAFileThatAlreadySaysIt` 断的是"字节相同＋mtime 相同"、fixture 里 `schema_version` 正是当前版 ⇒ **票 227 AC#2 那句"种一发删 schema 行必须红"是有靶子的，不是空想判据**。

**3. 它那 10 条推翻清单的处置（逐条，不整体接受也不整体驳回）**

| # | 内容 | 我的处置 |
|---|---|---|
| 1 | 实现腿 §4 报"红 5 枚"、v1 同形跑出**红 8 枚** | **接受、但归因不写死**：两枚读数的差（`TestAC2CleanWrite…:194`／`TestAC2Unreadable…:196`／`TestGuardedWriteDoesNotRewrite…:207`）我**没有复跑 w1 的突变**，只认"8 枚这一发是可复现的"。⇒ 记 w1 自述偏保守，**不改它的证据件正文** |
| 2 | 证据件 `226-config-write-no-clobber-r1.md` **有两个 `## 4.` 标题**、其中行 97 那节正文＝**`〔待填〕`** | ✅ **我读到原文复认**（`:97` 标题＋`:99` 逐字"〔待填〕"）＝交付缺陷，**具名入账**；处置＝见下面第 4 节（我只追加一句指针、**不抹原句**） |
| 3 | "`readConfigFile` 拆出止于 `validate`"没讲全：重读内部仍跑迁移并写备份 | **接受**（与我复读的 Y2/W2/W3 一致）⇒ 落进本票 AC#2 那一格的注，残余转票 227 AC#2 |
| 4 | 常驻正控有效、但 AC#1 覆盖面讲过了 | **接受**（我 09-29 现跑 13/0 两枚数坐实）⇒ 票面 AC#1 加硬注 |
| 5 | `schema.go:570` 那枚行号它未现核 | **不追究**（不影响任何一格；该腿自己已标"属起手自述"） |
| 6 | AC#5 的推迟**没落成代码 `DEFERRED` 标记**、也**没进 `SPEC-12 §5` 登记表** | **接受，责任判给我**（实现腿无写面）⇒ 见下面"悬空账"那一格，我在台账里按五字段接住 |
| 7 | 派单/票面的"发数"口径都不对：真数＝**25 枚顶层读数＋7 枚探针＝32**（顶层 26 里 `head-writeguard_226_test.go` 是测试件逐字副本、不是一发） | ✅ 我复跑枚数：`find -maxdepth 1 -type f`＝**26**、`find -mindepth 2 -name main.go`＝**7** ⇒ 与它的推导一致。**票面第 3 行那句"34 发"作废**（已由 `-v2.md` 取代，原句不抹） |
| 8 | "本票必须先于票 223"是排程陈述、不是盘上事实 | **接受**，已写进上面 Status 的更正② |
| 9 | AC#5 那半今天无主、要人拍 | **不用他拍**：按 A424/A428 既有口径，撤销动词归票 219 的答复语法面、会话档归票 224，两票都在队里 ⇒ 我只把**推迟登记**这枚责任接走 |
| 10 | "删行会不会被补回"这问有答案了、不该再挂"不由我猜" | **接受并销掉**：会被补回（标量／节补成 **schema 默认值**，map 条目不补）⇒ 已写进 Status |

**4. 追加到实现腿证据件的那一句（原句不抹）**：`docs/evidence/s1/226-config-write-no-clobber-r1.md:97` 那节空壳由编排者加一行指针指向同文件 `:101` 起的 `## 5.`，并具名记"两枚 `## 4.` 标题属交付缺陷、由本仓台账追认"——**见本次提交**。

**5. 悬空账（`deferred-work-must-be-registered` 那一规，责任在我）**：事项＝**"一条长期规则的撤销入口（答复语法）"**；已交付＝`Manager.SetAllowedDirs` 有真写路径＋两枚用例（此前全仓零调用者零用例）；未交付＝触发它的那枚答复动词；落点＝**票 219 的答复语法面**（撤销）＋**票 224**（会话档）；阻塞原因＝同一格三方重叠、新增答复动词属别票判据表；**未挂 `DEFERRED` 代码标记**（`internal/config` 里今天只有 `doc.go:14` 那枚属票 05 的标记，`226-v2` §⑥ 条 6 现跑复认）⇒ 因 `SPEC-12 §5` 的 1:1 要求要动 `docs/specs/**`（禁动），这半截**归口票 225（DEFERRED 与登记表双向对账）**，不在本票也不在票 227。
