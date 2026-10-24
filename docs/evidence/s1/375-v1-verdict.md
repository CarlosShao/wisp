# 票 375 终裁腿 375-v1：〔乙〕动作层达标｜**`AC#2` 能翻**

- 腿 `375-v1`（**非实现者**）。开工/收口 HEAD `b63d564130d`（本笔已提交）。本件只读主体：**未落 `.go`**、未改票名／`Status:`／`-done`、未碰台账／`internal/**`／`cmd/**`、未 push、未翻框。
- ★**写面**：本文件 ＋ 仓外 `D:\tmp\wisp\20261024\375v1\**`（7 支探针/overlay、31 份日志、变异树）。★**落点先 `ls`**：`375-v1-verdict.md` 与 `D:\tmp\...\375v1\` 均不存在 ⇒ 用原号，**未改号**。
- ★**尺的射程（落笔前五问）**：①对象＝**四个装配根各自**“此刻答得出吗”与**同一发 `fs.write`** 的路由；⛔“`Confirm` 是否 nil”（372 的尺）。②形状＝`V1P375 read=…` 一行 ＋ `[audit]` 原句逐字 ＋ `disk=absent|unchanged|present`。③数的是**代码**（`.go` 改动，⛔ 注释）。④总数核对：`.go` 加总 **980** ＋ 件 **110** ＝ **1090** ＝ `--stat` 的 `+1090`。⑤谁算的：整数由 `git show --numstat`；sha256 由 `sha256sum`（64 位）——⛔ 与 git blob（SHA-1）互比。
- ★**改前树怎么来的**：`b63d56413^` 的 10 枚改动 `.go` 导到仓外、3 枚新增映射为空串，`-overlay` 注入 ⇒ **改前树＝真 `HEAD^`**（`go vet` rc=0 先立住）。★现树＝`HEAD`（写面 `git status` 仅 `.scratch/**` 与他人台账脏，`.go` 0 枚）。

## ① 十一项逐项判语

### 1) 票面前提新鲜度（对 `HEAD` 现跑行号）

| 票面/普查断言 | 现量 |
|---|---|
| `store.go:201` 只判 `confirm==nil` | ★**仍真**：`gate.go` 未碰此面；`internal/perm` 零字节改动 |
| A 的腿 `run.go:668-671/:672/:674` | ★**行号已腐**（373 的 `taskLifecycleJournal` 推后）：现为闸门 `run.go:700`、那枚 `Confirm` `:743-746`→`perm.New :749`、`Modes` 包装 `:864`；**结构未变** |
| B `resident_task_source_windows.go:271-284` | ★**对上**：`spec := runSpec` `:271`、`gate` `:275`、`answer` `:283` |
| C `panel_inbound.go:236 nil`、`:272/:277` | ★**对上**：`:236` nil、`perm.New :272`／`Confirm: storeConfirmFromPanel :277` |
| D 腿 `:224`、球 `:298/:309`（早 85 行） | ★**逐字对上**（309−224=85） |
| 闸门守卫在 R7 升级之后、开窗之前 | ★**对上**：R7 `:283-299` → 守卫 `:315-323` → `openWindow :327` |
| 四条装配、无第五条 | ★**复核为真**：生产 `approval.New(` **2 枚**（`run.go:700`／`resident_approval_windows.go:395`），生产 `Confirm:` 4 枚 |

### 2) 核心判点一 `READ-F`（★我自己造；⛔ 采信）

| | 改前 `HEAD^` | 改后 `HEAD` |
|---|---|---|
| 读数 | `cards=0 executed=true elapsed=29ms disk=present:375-v1 byte` | `cards=0 executed=false elapsed=34ms **disk=absent**` |
| 伴随行 | `MODE-SILENCE mode=auto_approve assessed=L1 effective=L0` | `MODE-UNANSWERABLE configured=auto_approve effective=ask_every_step` ＋ `L1-UNANSWERABLE … decision=reject`（`MODE-SILENCE` **0** 行） |
| stdout | 无卡 | 无卡 |

⇒ ★★**判据成立**：**写入成功 → 拒**，且**文件不在盘上**。⛔ 只测“日志出现降级行”。

### 3) 核心判点二 `READ-S`（最严档 ＋ 无答复源）

| | 改前 | 改后 |
|---|---|---|
| 读数 | `cards=1 executed=true elapsed=2031ms disk=present`、`ANSWER-EXPIRED … decision=timeout->execute after=2s`、stdout 有 `[确认 L1 ` | `cards=0 executed=false elapsed=27ms disk=absent`、`L1-UNANSWERABLE … decision=reject`、`ANSWER-EXPIRED` **0** 行 |

⇒ ★★“到点即执行”的**前提**被补上：无人可答时窗口根本不开。

### 4) 核心判点三 `READ-W`（L2）＋`MODE-REDLINE` 漂移

| | 改前 | 改后 |
|---|---|---|
| 读数 | `cards=1 executed=false 31017ms disk=unchanged`；`timeout->reject after=31s`；`MODE-REDLINE … still asks` **n=1** | `cards=1 executed=false 31018ms`；**同一条** `timeout->reject`；`MODE-REDLINE` **n=0** |

★★**漂移判语：可接受**。第三发把它夹住了：**有答复源**的 `READ-PL2`（`reply=pipe`、`auto_approve`、覆盖写）改后 **`MODE-REDLINE` n=1 仍在** ⇒ 漂移**只**发生在“有效档已降到最严档”的场景，而那时红线豁免**确实不再是理由**；`configured=auto_approve` 由 `MODE-UNANSWERABLE` 携带。★L2 的**问与拒**一字未动。

### 5) 核心判点四：四个装配**逐装配各造一发**（⛔ 一个读数概括四个）

| # | 我的发 | 改前 | 改后 | 判语 |
|---|---|---|---|---|
| **A** | 控制台 ＋ `fs.write` 新文件，两档 | `auto_approve`：0 卡写入成功；`ask_every_step`：1 卡 2s 后写入 | **两档都拒**（0 卡、`L1-UNANSWERABLE`、文件不在盘） | ✅ |
| **B** | 常驻管线：注入 `ra.gate`、`auto_approve`、**无球**、`fs.write` 新文件 | `cards=0 executed=true 29ms disk=present`；`MODE-SILENCE n=1` ⇒ ★**根不到闸门**（`errResidentNoBall` 这条“恰好诚实”的面**根本没被问**） | `cards=0 executed=false 32ms disk=absent`；`MODE-UNANSWERABLE` ＋ `L1-UNANSWERABLE`（读自闸门自己的 `residentAuditf`／slog） | ✅ **洞补上** |
| **C** | 无闸门／无 bridge | 372 的 nil-leg 降级 | **零字节改动**（`panel_inbound.go` 不在 14 文件内） | ✅ 无需动 |
| **D** | ① 腿在球之前交出的那个事实 ② 谓词活性 ③ 强制面 | `hasBall=false`（`:224` 时刻） | ① `gate_src_nil=false src=false` ② `attach→true`／`release→false` ③ 包装器 `noBall=ask_every_step`（1 行 `MODE-UNANSWERABLE`）／`withBall=auto_approve` | ✅ 活谓词覆盖“早 85 行” |

★**D 的一处收紧**：D **自己那条路**（`panelModeConfirm(ra.confirmModeSwitch)`→`PendingApproval`，**L2**）改前／改后**逐字相同**：`err="L2 强确认未通过（reject）：审批界面不可达，已 fail-closed 拒绝执行"`、`0ms`。⇒ 375 对 D 的贡献是**同一枚闸门现在答得出**（①②③），**⛔ 是“D 的路由变了”**；`375-r1` §⑦ D 行“球不在＝拒”读起来像新行为，**措辞过头、效果为真**（见 §⑬ 残留 1）。

### 6) 核心判点五：正控（有答复源时一字不变）

| 读数 | 改前 | 改后 |
|---|---|---|
| `READ-PL1`（L1 ＋ `io.Pipe` 无人打字） | `cards=1 executed=true 2035ms`、`timeout->execute` | `cards=1 executed=true 2046ms`、**同句**、`L1-UNANSWERABLE` 0 行 |
| `READ-PL2`（L2 未答） | `cards=1 executed=false 31032ms`、`timeout->reject`、`MODE-REDLINE n=1` | `31030ms`、**同句**、`MODE-REDLINE n=1` |
| `READ-PL2A`（L2 答 yes） | `executed=true 21ms`、`ANSWER-ALLOW route=native` | `executed=true 29ms`、**同句** |

⇒ ★★**正常路径没被改**（毫秒级差异是测试开销，形状与判决句逐字相同）。

### 7) 核心判点六：那 2 处“为何被拒”的句子

| 文件 | 只动了什么 | 结果断言 | 判语 |
|---|---|---|---|
| `resident_approval_246_windows_test.go` | `为什么拒`：`确认界面不可达`→`无人可答`（1 行） | `ans != tools.AnswerReject`、`AwaitingHuman()`、`displayedCards()==0` **一字未动**；两树都 PASS | ✅ 可接受 |
| `resident_task_source_246_windows_test.go` | 读数从 **L1** 挪到 **L2 面**（先把目标写成已存在文件 ⇒ R8）、哨兵改成从 **闸门自己的 slog** 读 | `!out.IsError` 未动；“文件不落盘”改成**内容未变**（L2 覆盖写必须先存在）；`审批界面不可达` 仍断言；`errResidentNoBall` 已在 `gateLog` 里 | ✅ 可接受 |

★★**“L1 覆盖有没有丢？”**：★**装配级“L1 ＋ 真实常驻面”这一组合没人测了**（现在无球时 L1 在问面前就拒，有球又造不出真 orb）；★**但覆盖不是零**：①闸门→`ra.ui` 的**指针身份**仍由同文件 `ar.gate == ra.gate` 钉住；②`g.ui` 是**一个字段**，L1/L2 共用；③L1 的“可答⇒开窗⇒到点执行”由 `approval` 包 `Test375L1AnswerableKeepsTimeoutMeansExecute`（假 UI）接管。⇒ **判：可接受，具名残留**（§⑬ 残留 2）。

### 8) 核心判点七：`nil ≠ false`

★**我的发**：`approval.New` **不带** `Answer`（＝生产外的旧宿主形状）＋ L1 决策 ⇒ 改前／改后**同形**：`prompts=1 answer="timeout" why="" elapsed=2001ms`（卡照样打、到点照样执行）。
⇒ ★★**判“未声明＝保持今天行为”对**：①生产里 **nil 不可达**——2 枚生产闸门**都**声明了（`run.go:711 Answer: answer`／`resident_approval:410 Answer: ra.ui.hasBall`），且注入闸门缺谓词时 `assembleRuntime` **结构性拒绝装配**（`run.go:676` 缺谓词／`:689` 闸门自己没声明）；②若把 `nil` 读成 `false`，会连 12 枚既有 fake 闸门一起判死——那是**清洗**，不是修。★主流形状（presence 判据 ＋ 动作层 `unavailable`）正好是这一分工。

### 9) 核心判点八：反向判据（★摘跳；从最终树导出）

| 变异（仓外导出 ＋ `-overlay`） | `READ-F` | `READ-S` |
|---|---|---|
| **M1 摘模式读跳**（`Modes: rt.modes`） | `cards=0 executed=true 34ms disk=present`（回写入成功，`MODE-SILENCE n=1`） | `cards=0 executed=false 31ms`（守卫单独够用） |
| **M2 摘闸门跳**（删 `!g.answerable()` 整块） | ★★`cards=1 executed=true **2043ms** disk=present` ⇒ **回到「2s 后写入成功」** | `cards=1 executed=true 2063ms disk=present` |
| **M3 两跳都摘** | `cards=0 executed=true 76ms disk=present` | `cards=1 executed=true 2041ms disk=present` |
| **还原核（★最后一发变异 M3 之后）** | `cards=0 executed=false 29ms disk=absent` | `cards=0 executed=false 33ms`（`READ-W` 31020ms 仍 `timeout->reject`） |

★**还原核对账**：`sha256sum` 五枚（`run.go a69ef890…`／`gate.go cab0caa7…`／`reach.go be110e95…`／`resident_task_source e8c6eb83…`／`resident_approval 9d246ed1…`）与 `git show HEAD:` 同文件**逐一相等**（64 位对 64 位）⇒ 上面各行同出一个最终树状态；★也证明**写面零未提交脏**。

### 10) 四数复跑（四包分开、⛔ 相减；★我的尺，⛔ 与它的 69/83 相减）

| 尺（我现造） | 改前 | 改后 |
|---|---|---|
| `cmd/wisp` 定向 **41 名**（票面 9 个文件里全部 `Test*` 顶层名） | **37 跑：36 pass / 1 fail**（12 子案 0 fail） | **41 跑：40 pass / 1 fail**（14 子案 0 fail） |
| `cmd/wisp` 超集 **77 名**（＋ 所有 `approval.New(`／`assembleRuntime(` 文件的用例） | **73 跑：72 / 1** | **77 跑：76 / 1** |
| `internal/agent/approval` | **78 top ＋ 19 sub ＝ 97 / 0** | **84 top ＋ 24 sub ＝ 108 / 0** |
| `internal/tools`（零字节改动） | —（同尺） | **250 top ＋ 87 sub ＝ 337 / 0** |

★**唯一红**＝`TestAC246DevLegIgnoresTheTestTaskInjection`（**两树都红**、`:37s`，失败句逐字 `wisp: another instance is running in this session`＝**桌面态、非本笔**）⇒ **新增红 0 枚**。

### 11) 边界逐枚

- 14 文件＝**13 `.go` ＋ 1 件**；`+1090/−16`。★**逐枚**：新增 3（`reach.go 162`／`reach_375_test.go 274`／`approval_unanswerable_375_test.go 308`）；改 10（`run.go 84/1`、`gate.go 52/3`、`resident_approval_windows 9/0`、`resident_task_source_windows 6/0`、4 枚 `cmd/wisp` 测试、`resident_task_source_246 45/10`、`resident_approval_246 9/2`）。
- ⛔ **`internal/perm`／`internal/risk`／`internal/tools` 零文件**（`git show --name-only` 逐枚核）。
- ★**`run.go` 的 `Prompt` 命中数 ＝ 0**：`git diff b63d56413^ b63d56413 -- cmd/wisp/run.go | grep -c '^[+-].*Prompt'` ⇒ **0**（与票 376 落点不重叠）。
- ★**动作层可答性只有两个读者**（`grep` 现跑）：`Gate.answerable()` 调用点**1**（`gate.go:315`，L1）；`UnsilencedMode` 生产装配点**1**（`run.go:864`）⇒ 没有第三条暗路。
- ★**边界：R7 升级仍走 L2**（我的发：`Answer: func() bool {return false}` ＋ 60 条 `Paths` ⇒ `answer="reject" why="审批界面不可达，已 fail-closed 拒绝执行"`，**⛔ “无人可答”**）⇒ 守卫**没有**吃进被升级的 L2 问句。

### 12) 门禁 rc（逐条现跑、`cmd > out 2>&1; rc=$?`）

| 门 | rc | 读数 |
|---|---|---|
| `go build ./...` | **0** | 空 |
| `go vet ./cmd/wisp ./internal/agent/approval ./internal/tools` | **0** | 空 |
| `tools/d22scan/d22scan.exe -root ../../` | **0** | `clean - no D22 ban violations`（ban #8：`internal/` **645**、`cmd/` **152**） |
| 格式**甲尺** `gofumpt -l` **13 枚** | **0** | 输出**空**（⛔ `gofumpt -l .`：`.scratch` 探针会 `rc=2`） |
| 格式**乙尺** `gofmt -l` 同 13 枚 | **0** | 输出**空**（判归属用 `[\\/]`，无需：输出为空） |

### 13) 提交号 ＋ 逐笔越界 ＋ 我抓到的两处**记账**缺陷

- 被裁对象：**`b63d56413`**（主题逐字 `票 375 落地腿 375-r1：…`），14 文件、`+1090/−16`。**逐笔越界：0**。
- ★**缺陷 A（记账）**：`375-r1` 件 §⑫ 写“本笔 **12** 枚 `.go`：**9** 枚改动（`+209/−12`）＋ 3 新增”。现量 **13** 枚 `.go`：**10** 枚改动（**`+236/−16`**）＋ 3 新增；`236+744=980`＋件 110＝**1090** 才与 `--stat` 对上。★它自己的**提交正文**写的是 **13** ⇒ 件 §⑫ 是笔误（不影响交付，但两处数打架，⛔ 当“两数自洽”）。
- ★**缺陷 B（记账，⛔ 逻辑）**：§⑪ 写“超集 **83 名 → 83 pass / 0 fail**”。若那 83 真含它自己点名的 69（其中 1 枚桌面态红），则**不可能 0 fail**。我的 **77 名超集**在**两棵树**上都保留**同一枚**红 ⇒ 该 83/0 要么集合不嵌套、要么两发桌面态不同。**这处我不采信、也不据此扣分**（本笔绿否由定向尺判）。
- ★**测试适配 9 处**复核：**7 处给答复源**（`runFixture.reply` 字段＋`TestComposedGate…`／`panel_l1_blocked_220`／`run_mode101(t101host.start)`／`ticket224_assembly`／`resident_task_source_246` 两条注入支）＋**2 处句子**（§7）。**逐处读过 diff：无一处放宽断言**。

## ② `AC#2` 能翻／⛔翻 ＋ 缺什么

★★★**能翻**。★**验收判据唯一一条**：`READ-F` 从“写入成功”变“拒”。★我的两发：`executed=true 29ms / disk=present` → `executed=false 34ms / disk=absent`（＋`MODE-UNANSWERABLE`＋`L1-UNANSWERABLE`）。★**缺什么**：**不缺**。★展示面／记账面仍报配置值（`375-r1` §⑬1 自述、我是“它报得对不对＝对”：`internal/perm` 零字节、`MODE-UNANSWERABLE` 两值并列可见）⇒ 归**票 376（已结）／票 377**，**⛔ 阻塞本票**。

## ⑬ 残留（具名、⛔ 藏）

1. ★**D 的措辞过头**：D 自己那条 **L2** 路改前改后逐字同（`0ms` 拒）⇒ 375 对它**零行为增量**；增的是**共享闸门的可答性**（§5 D）。
2. ★**“L1 ＋ 真实常驻面”组合无人测**（§7）：身份由 `ar.gate==ra.gate` ＋ L1 可答路由 `approval` 包假 UI 覆盖。
3. ★**谓词是“存在答复流”而非“有人正在打字”**：`strings.NewReader("")`（7 处适配用的就是它）在**EOF 之后**仍算可答 ⇒ 该边角**保持改前行为**（到点执行）。普查 A 行已划此界，**非新增洞**。
4. ★**A 的 L1 否决通道仍 0 枚**（`run.go` 给闸门空 `Channels`）；正控要求它不变，故本笔未动（`375-r1` §⑬3）。
