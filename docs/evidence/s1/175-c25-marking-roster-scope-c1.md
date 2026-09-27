# 175-c1 — 只读裁决：C25 污染源标记名册是穷举还是举例

- 派单：`.scratch/wisp/dispatches/2026-09-27-201x-readonly-175-c1-is-the-c25-marking-roster-exhaustive-or-illustrative.md`（`2b2a01de`）
- 角色：只读裁决位。本程唯一写面＝本文件。产码／测试／尺／票面一字节未动。
- 裁决时刻：2026-09-27 20:2x +08｜锚点：`git rev-parse HEAD` = `2b2a01dec950cba92ca754ce883cc800af86d7d1`，分支 `dev`。

## 1. step-0 五件（可复制命令与读数）

| 件 | 命令 | 读数 |
|---|---|---|
| 时间 | `date` | `Sun Sep 27 20:20:31 CST 2026` |
| 分支 | `git rev-parse --abbrev-ref HEAD` | `dev` ✓ |
| 锚点 | `git rev-parse HEAD` | `2b2a01de…` ✓ |
| 干净 | `git status --porcelain -- internal/ cmd/ docs/specs/` | 空 ✓ |
| 名字在不在规格 | `grep -c 'task\.output' docs/specs/SPEC-06-security-gatekeeping.md docs/PLAN.md` | **SPEC-06：0｜PLAN.md：1**（唯一命中 `PLAN.md:2564`，D34 表 task.output 行）；`grep -n 'task\.' SPEC-06` 只有 `:108` 的 `task.list`/`task.cancel`。派单 19:5x 现量前提**属实**。 |

## 2. 本程没读／没测什么

- 没读 `frontend/**`、`design/**`（他人领地，未读未引）。
- 没跑任何 `go test`（本程无需读数，一条台件都没建，`.scratch/wisp/probes/175/**` 不存在）。
- 没核 `Q-59`／票 174 的原文（AC#5 分账只按票 175 面与 `164-v1` §7.3 已载事实对账，见 §8）。
- 没独立复测「`RunAsync`/`Record` 生产零写者」（引票 175 面与 `docs/evidence/s1/164-task-output-accept-r1.md:496` 的现量，二者一致）。

## 3. §1 四段原文（逐字＋行号）

### 3.1 SPEC-06 §5（`docs/specs/SPEC-06-security-gatekeeping.md`）

```
70	## 5. C25 Provenance 污染追踪（F4，六个外泄通道）
72	- **敏感源**（读取即打 `sensitive` 标记）：`fs.read`（含越界尝试）· `search.content` 结果 ·
73	  `clipboard.read` · `sysinfo` 的当前窗口标题（高价值泄漏源）· `web.fetch` 响应 · `doc.read` ·
74	  `screen.capture` 图像。
75	- **外泄通道**（六个，不只 HTTP body）：`web.search` query 串 · `notify` 文本与 URL ·
76	  **TTS 播报**（物理出网）· `clipboard.write` · **`fs.write` 到同步盘目录**（OneDrive/Dropbox/
77	  坚果云自动上传，最隐蔽）· HTTP body。
```

敏感性名单的写法计量：`:72-74` 共 **7 枚**，以「。」收尾；**既无**「等／例如／包括但不限于」，**也无**「仅限／只有／以下七种」——单看这段，穷举与举例两读都不违文字。裁决不能停在这一段，见 §3.2。

### 3.2 PLAN.md 的三处决定性文字（同族规格互证）

```
PLAN.md:2470-2474（D33/F4，与 SPEC-06 §5 同源重写）：
**修复 = C25 `Provenance`（污染追踪）**，统一模型：
- **敏感源**（读取即打 `sensitive` 标记）：`fs.read` 越界尝试 · `search.content` 结果 ·
  `clipboard.read` · **`sysinfo` 的当前窗口标题**（标题常含文件名/URL/聊天对象名，
  这是一个容易被忽略的高价值泄漏源）· `web.fetch` 响应 · `doc.read` 内容 · `screen.capture` 图像
```

**甲·决定性一句（D46，owner 已批决策，`PLAN.md:2661`）**：

```
| **输出必须经声明式提取器** | `extract = {kind: "json", path: ".data.items[].summary"}` 或 regex；**原始 stdout 不直接进 LLM 上下文**（截断上限 ≤32KB） | 上下文爆炸；也让 C25 能给输出打 taint |
```

`command` 类插件的输出＝外部子进程 stdout，**这个名字不在 3.1/3.2 的任何一张枚举里**，而规格原文要求「让 C25 能给输出打 taint」。若枚举是穷集，这一行自相矛盾；若枚举是当时已知来源的快照，这一行自然。

**乙·同型前例（D43 转移表，`PLAN.md:3070`）**：

```
| 11 | `Listening` | VAD 判停且语音时长 ≥300ms | `Thinking` | 停采集；标点恢复；输入打 taint 源标记 |
```

正是这条催生了名册第 8 枚 `asr.transcript`（`provenance.go:80-82` 注明「promised to ticket 15」）。也就是说：**本仓已经发生过一次「SPEC-06 §5 名单之外补进一枚盖戳来源」**，当时无人按契约追加流程走。

**丙·C25 契约行本体（`PLAN.md:1375`）**：`Provenance` ＝「**工具结果的来源标记**（`{tool, origin, sensitive}`）＋出网前的粗粒度污染匹配……D30① 组合闸门的实现契约」——措辞是「工具结果」，不带名单限定。

### 3.3 反证，如实报（读 B 支的人唯一能抓的文字）

```
PLAN.md:2543	| `search.content` | 全文检索 | L0 | `fs.read` | S3 | 结果打 taint（C25） |
PLAN.md:2564	| **`task.output`** | **读一个后台任务吐了什么**（含"太长怎么续读"） | **L0** | — | **S7** | 票 164（2026-09-27 owner 批准新增）的输出腿。截断按 **D15「单个工具结果」那一行的既有规矩**……
```

D34 表的备注列给 `search.content` 明写「结果打 taint（C25）」，给 `task.output` **没写**。若把备注列读成「逐工具盖戳指令的唯一清单」，则 task.output 不打戳＝规格本意。**但此读法不成立**：备注列若穷举盖戳指令，它连 `fs.read`/`web.fetch` 都没标（那些靠 §5 名单），而 `:3070`/`:2661` 两行恰恰证明盖戳指令在表外也存在。备注列是提醒性注记，缺席是沉默，不是否定。

### 3.4 代码侧被审的转述（`internal/risk/provenance.go`）

```
80	// Sensitive source tools whose output is marked (SPEC-06 §5). system.get is
81	// the sysinfo focused-window-title source; asr.transcript is the user
82	// transcript hook promised to ticket 15.
83	const (
84		SrcFSRead        = "fs.read" // incl. out-of-allowlist attempts
85		SrcSearchContent = "search.content"
86		SrcClipboardRead = "clipboard.read"
87		SrcSystemGet     = "system.get" // focused-window title (high-value leak)
88		SrcWebFetch      = "web.fetch"
89		SrcDocRead       = "doc.read"
90		SrcScreenCapture = "screen.capture"
91		SrcTranscript    = "asr.transcript"
92	)
94	// sensitiveSourceTools is the SPEC-06 §5 set (marker-side only; a Mark() for
95	// any other tool is still honored fail-closed, just logged).
```

`:94` 那句 `the SPEC-06 §5 set` **逐字不实**：SPEC-06 §5 写的是 **7** 枚（且第 4 枚叫 `sysinfo`，不叫 `system.get`——名字映射属实），代码名册是 **8** 枚（多出 `asr.transcript`）。这张 8 名册在任何规格文件里都**逐字不存在**——它是「SPEC-06 7 枚＋D43 表 1 枚」的拼装快照。

## 4. 裁语：**A 支＝举例／凡外部内容都要盖戳**（漏项是实现缺陷）

判据链（每一环都指向上面具名行号，不用「注释说是 §5」当证据）：

1. 枚举本身两读都不违文字（§3.1 的计量），所以裁决必须由规格族其余文字定；
2. `PLAN.md:2661`（D46，已批）要求名单外的外部内容「让 C25 能打 taint」，`PLAN.md:3070`（D43，已批）要求转写盖戳——**只有 A 读法能让规格族自洽**；B 读法使这两行与 SPEC-06 互相矛盾，其中 asr.transcript 入册在 B 读法下已是违规行为，而它已落地且无人退回；
3. 引擎语义与常驻测试都按 A 建：`Mark()` 对名册外工具照样记录、只多一条日志（`provenance.go:489-491`），`provenance_test.go:135-140` 逐字钉着 *"unknown-source marks must still be recorded (fail-closed)"* 与 *"taint from an off-list source must still gate"*——设计者明示名册不是安全边界；
4. 真正的缺口是 `bridge.go:552` 的调用方早退把「按定义返回外部内容」的门整个跳过盖戳，与契约（C25＝工具结果的来源标记，`PLAN.md:1375`）方向相反。

⇒ **补上 task.output 的盖戳＝修 bug，不需 owner 批**；票 163 写手的实现程**可以派**。
代价与边界（A 支自己要吃下的）：
- 修法不得动 `docs/PLAN.md`/`docs/specs/**` 一字（含往 `:2564` 行补备注——那是文字契约，虽非必需，别顺手做）；
- `provenance.go:94` 那句不实转述应改（注释级修正，任何下一次合法触碰该文件的程捎带即可，不单开票）；
- 名册外延此后归实现判断，**必须由 AC#3 的常驻正控用例兜住**（见 §6），否则 A 支等于把安全结论再次寄存于「下次有人记得」。

## 5. 20:1x 落点收窄复核：**对**（三条全中，附本程独立读数）

- `Mark()` 接受任何工具名：`provenance.go:474` 签名无名单检查，`:489-491` 名册外仅 `logf` 后照记；常驻测试 `provenance_test.go:135-140` 把这条钉成行为契约。✓
- 唯一闸门在调用方：`internal/tools/bridge.go:552` `if b.prov == nil || dec.TaskID == "" || !risk.IsSensitiveSource(dec.Tool) { return }`——行号逐字如此。✓
- ⇒ 最小修法（tools 侧调用方判定形，或桥侧自持一枚「外部内容工具」局部名单）住在 `internal/tools/**`，`internal/risk/**` 一字节不必动。✓ 唯一保留：若实现程选择改 `sensitiveSourceTools` 本身，在 A 裁语下那是改实现清单（8 名册逐字不在任何规格里，见 §3.4），仍不违例——但不是最小形。

## 6. 排程裁语：**不会红**（没有具名文件＋步骤可报，故按派单规则报"不会红"）

票 163 把写者接上之日，逐枚排过的候选全绿：

- `internal/tools/task_output_leg_test.go`／`task_output_ac2_before_test.go`：断言的是可达性、错误形状、D15 截断三件套与「未接线 fail-closed」，**没有一枚断言 task.output 的内容被盖戳**（`grep -rn 'Mark' internal/tools/*_test.go` 命中仅 `bridge_scope_open_ticket158_test.go` 对 `fs.read` 的 OpenTask 正控）。taskBridge（`task_output_leg_test.go:35-51`）虽接了真 prov，但 `bridge.go:552` 早退使它恒为 no-op。
- `internal/risk/provenance_test.go:146-161`（`TestSensitiveSourceSetComplete`）：只钉 8 枚在册＋`fs.list` 不在册，写者落地不触它。
- CI：跑这些包的是 `.github/workflows/ci.yml` step **"Portable package tests (core scope)"**（`scripts/portable-tests.sh --scope=core`，清单含 `./internal/tools/...`、`./internal/risk/...`，`portable-tests.sh:176-177`）——上述测试在该步全绿即过。`gate-clauses.sh`（`.scratch/wisp/probes/154/`）**在 ci.yml 里零引用**（`grep -n 'gate-clauses|probes' .github/workflows/ci.yml` 无命中），它红了也没有流水线会拦。
⇒ **这条「换个门读同一份外部内容就不盖戳」的路会自动通，且无任何既有用例／CI 门禁会红。**

**最窄判据该长在哪一侧（两形选一）**：**常驻用例**。放 `internal/tools/`：同一 taskBridge 形状下让 task.output 返回一段含探针串的产出，断言 `prov.ScopeTaints(scope)` 出现 `task.output` 条目（正控，0/非 0 进判据）。不选 gate-clauses 腿的理由：那一族今天不挂任何 CI 步，红了无流水线听见；而 `./internal/tools/...` 已在 core scope 里被执行，正控一红 CI 即红。**本程不落任何判据代码。**

## 7. Q-61 草稿

裁语为 A，按派单 §2 该草稿仅在 B 支时交付，**不交付**。若 owner 复核推翻本裁语改判 B，草稿要点预置：甲＝把「凡外部内容必盖戳」升格写进 SPEC-06 §5（动规格文字）；乙＝task.output 补进名册并登记为契约追加（动代码清单＋补票面）；丙＝接受现状、把豁免登记成五字段 DEFERRED。**推荐甲**（它使 `:2661`/`:3070` 两行从此有据）。撤销口令形状：`Q-61 撤销：回到裁语 A（举例），实现程按 §5 最小形派`。

## 8. AC#5 分账（票 164／174／Q-59）

- 票 164：task.output 的宿主已验收（`164-v1`，其 §7.3 即 20:1x 收窄的来路，本程 §5 独立复算为真）；本裁语不结 164 的账。
- 票 174／`Q-60`（`[fs]` 授权根）：同一根管子的另一处，**批一枚不通另一枚**；本程未读其票面，不裁。
- `Q-59`（偏移）：本程未核其原文，只在此具名登记状态为「与本裁语无涉、未复核」。
- 三处合读：本裁语只结「戳该不该盖」一处；163 写手落地判据（AC#3 硬约束：与写者同批或之前）不因本裁语放松。

## 9. 程序合规

- 被拒／没成功的调用：**无**（取数前后均无）。
- 删除命令：**没跑过**（全程只有 `date/git rev-parse/git status/grep/sed -n 读`＋本文件写入＋一次带显式 pathspec 的 `git add`/`git commit`）。
- 伪授权两栏：**①** `provenance.go:94` 注释自称 `the SPEC-06 §5 set`——规格实为 7 枚异名清单，此转述已按 §3.4 判不实，未采信；**②** 票 175 Status 行「两枚写面都属已批准射程之外」——被票面 20:1x 追加段②（`:33`）自行收窄为「先只裁早退该不该改」，本程按收窄后的形状裁，未引用收窄前旧句作授权。
- 凭据值零抄录。✓
- 脏件不碰：`design/**` 删除、`.gitignore`、`probes/152/my152.py`、`probes/161/r6/logs/flip-*.txt` 均未触碰、未入本次提交（`git status --porcelain -- internal/ cmd/ docs/specs/` 起手即空，`.scratch/**` 脏件不在写面）。

next= 编排者拿本裁语（A 支）后即可派票 175 实现腿：写面 `internal/tools/bridge.go` 的 `mark` 判定形（§5 最小形），**与票 163 写者同批**，AC#3 用 §6 具名的常驻正控用例形状；`provenance.go:94` 注释修正在该程内捎带；SPEC-06/PLAN 一字节不动。
