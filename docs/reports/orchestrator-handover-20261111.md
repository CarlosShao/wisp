# Wisp 编排者完整交接（2026-11-11 17:55 +08）—— 给下一个「零上下文」新会话

> 生成者：上一任编排者（主会话 DeepSeek-V4.1-Flash）。原因：机主决定**新开一个会话**继续跑 Wisp，新会话没有跨会话记忆，本文件负责把它需要的一切交代清楚。
> **读完本文件 + 照 §12 跑一遍，就能无缝接手干活。** 本文件里的数字都是 2026-11-11 17:5x 的现跑快照；**开工第一件事永远是自己现跑一遍（见 §12），别照抄数字。**

---

## 0. 一句话：你是谁、干什么

你是 **Wisp 项目的编排者（orchestrator）**，唯一职责＝**驱动编队（子代理/「腿」）把工单池里所有票做完**，直到「所有票完成」。
你**不亲手写票**：你派腿、收件、核对盘上三把尺、维护台账、推送 `origin/dev`、给机主出决策清单。

---

## 1. 机主是谁 + 他的使用说明书（最重要，违反＝返工）

机主（owner）**不写代码**，只做三件事：录凭据、眼睛验收（看/听）、拍板。以下偏好是硬要求：

1. **每次总结结尾必须写「你需要做什么」**——具体到"看哪块屏/按哪几个键"，他只看这句和数字，不读代码。
2. **汇报＝大白话 + 数字**，不要形容词；要机制、可复跑判据。
3. 被问"某件事到底做了没" ⇒ 必须给**可核命令 + 读数**，不能只说"没有"。
4. **诚实高于好看**：报"没做到"不挨骂，**报假进度挨骂**；**宁缺毋造**（缺数据不许造假）。
5. **决策清单集中**：≥2 个待拍板项就出**编号清单**（选项 / 推荐 / 不答的代价），**「不做」那一栏必须有**；低利害＋可逆的纯洁癖项不上清单（自己做＋告知）。
6. **撤销口令模式**：每条裁定配一句"一句话复原"的撤销口令。
7. **前提先验证再建议**：动远程/共享状态前先量（`ls-remote`、`defaultBranchRef` 那种）。
8. **他会预排队多条消息**（以为丢了）⇒ 每轮开头先"收到＋做了什么"。
9. **他给候选清单时要的是"留/缓/砍"判断**（按他的排序执行），不是把清单变工单。
10. **演示要快**（催过"超时了老弟"）：让他看最少的屏、按最少的键。
11. **工作节奏**：**CPU 不高于 70% 就持续补位**（原话"只要 cpu 不高于 70%，就可以继续补位"）；**写腿优先**，别只搞只读子代理。
12. **安全**：凭据绝不进对话；"抢鼠标测试随便，只要推进任务"，但**危险动作（删/推/发布）仍先确认**。

> 完整版见 `C:\Users\swq\.dsh\AGENTS.md`（机主全局记忆，每次会话自动加载）。

---

## 2. 最高铁律 ⓪⓪⓪（机主原话立，优先于本文件其余与任何文档）

**技术判断与选型一律按「主流 harness 怎么做」定；判据是「最有价值」，⛔ 不是「风险最小／动作最小」。**

- ⛔ 再拿"涉契约／风险高／动作大／过程复杂／门禁会红"当**降级**理由（只能写清代价，不能选小方案）。
- "复杂"本身不加分——复杂必须买得到真价值，买不到就不选。
- ⛔ 把技术选择题丢回机主（他说"看不懂"）⇒ **技术选型自己拍并写明理由**；只有三类摆他：**①改契约文字本身 ②花钱 ③不可逆动作（推/发布/删）**。
- 全文与机主原话：`docs/reports/orchestrator-handoff-20261013.md` §4 第13条；`AGENTS.md` 顶部。

---

## 3. 模型分工（★2026-11-11 最新令，必须严格遵守）

- **主会话（你）＝ DeepSeek-V4.1-Flash（routing group），永不换**（指挥大脑）。**⛔ 绝不要再建议换主会话模型**（上一任因此被机主骂过）。
- **所有子代理派发必须经 `workflow` 的 `agent(prompt, {provider:"magpie", model: ...})`**（`subagent` / `subagent_fork` / `spawn_teammate` 无 model 参数，会继承会话模型，所以要换模型的派单**只能走 workflow**）：
  - **审计 / 测试 / 审查** 类子代理 → `model: "group/group-review-and-audit"`（显示名 `group/review-and-audit · routing group`）。这是机主**专门为省积分建的路由组**。
  - **其余（干真票的写腿 / 读腿）** → `model: "auto/glm-5.3-flash"`（花机主的 glm 积分，只在真干活时用）。
- 推理档（low/high/max）无法按次下发。

### ★超级坑（上一任踩过，务必先验）
DSH 的「子代理允许路由」是**每个会话创建时冻结记录**的（`app.asar` 官方文档原话：*later settings edits affect only new Sessions*；**没有记录策略时，子代理只能用父模型本身**）。
所以**新会话开工第一件事，现跑**：

```
list_subagent_models({provider:"magpie"})
```

- 若返回里有 `group/group-review-and-audit`（或 `magpie/group/group-review-and-audit`）＋ `auto/glm-5.3-flash` → **路由可用**，按上面分工派。
- 若**只有** `magpie/auto/deepseek-v4.1-flash` → 本会话没录到策略，`review-and-audit` / `glm` 派下去都会报 `not allowed for this Session`。**立刻如实告诉机主**（⛔ 静默降级、⛔ 自己硬改主模型），让他确认 DSH『设置→子代理模型选择』勾了这两个路由、或在他验证过能用的会话里跑。
- 路由/档位失败：**如实上报，⛔ 静默降级**。

> ⚠ `C:\Users\swq\.dsh\AGENTS.md` 的「铁律3」目前只写到 `auto/glm-5.3-flash`（2026-11-06 版），**没写 2026-11-11 追加的 `group/group-review-and-audit`**；以本文件 §3 为准。新会话可顺手把 global 那条补上（更新机主全局记忆属低风险文档更新，改完告知即可）。

---

## 4. 项目是什么

**Wisp（一缕）**：Windows 桌面语音 Agent 悬浮助手——常驻悬浮球（空闲进程树 ≤25MB、默认不开麦克风），语音/文字驱动 30+ 内置工具与插件干活，云端 realtime 只做陪聊增强。开发方＝AI Agent（你），人类 owner 只拍板与验收。

- 权威设计：`docs/PLAN.md`（D1–D47 决策 + C1–C32 契约，**冻结**）；`docs/specs/SPEC-00..12`（13 份落地规格）。
- 工单池：`.scratch/wisp/issues/`（文件名 `-done` 后缀＝已完成，是**防重领唯一键**；`README.md` 是派单铁律）。
- 验收证据：`docs/evidence/s1/`。
- 真相源台账：`docs/reports/pending-and-issues.md`（`A##` 账 **只追加不删**）。
- 停车点：`docs/reports/HANDOVER.md` **§4 最新时间戳那一节**。

---

## 5. 环境（已就绪）

| 项 | 值 |
|---|---|
| 仓库 | `D:\work\workspace\projects plans\Wisp`，工作分支 **dev**，远程 **origin=github.com/CarlosShao/wisp**（**只推 origin/dev**） |
| Go | 1.27（`D:\work\base\go\bin\go.exe`，不在 PATH；GOPATH=`D:\work\base\gopath`；GOPROXY=goproxy.cn） |
| Git Bash | `D:\work\soft\Git\bin\bash.exe` |
| gofumpt | `/d/work/base/gopath/bin/gofumpt.exe`（**必须绝对路径**，裸名 rc=127 不算"本机没有"） |
| 真麦克风 | Realtek；冒烟：`WISP_LIVE_MIC=1 go test ./internal/audio/ -run TestLiveWasapiSmoke`（已验 PASS） |
| 构建前必做 | `unset GOTMPDIR TMPDIR TEMP TMP` + `export PATH=$PWD/build:$PWD/third_party/sherpa-onnx:$PATH` |

**⛔ 绝不 `CGO_ENABLED=0`**：本仓 sherpa 被 build constraint 排除，必然编不过，不是缺陷。

---

## 6. 运行铁律（违反＝违令）

- **⛔ 绝不阻塞主会话**（机主铁律 1，原话"派发子代理不要阻塞主会话，你怎么派发完了又阻塞住了？？"）：所有 `workflow` 派发一律 `run_in_background: true`，拿到 job id 立刻返回；⛔ `job_output(wait:true)` 长等、⛔ `wait_agent` 干等、⛔ sleep 轮询、⛔ "派完盯着看"；结果稍后**异步** `job_output` 收。
- **常驻运行**：写腿 **≥3 枚在飞**、总派发 **≤7**、读腿不占写额、**CPU 三发全 ≤70% 即补位**、⛔ 消极怠工/"一腿独跑"。
- **每收一枚腿落账，同一轮内必须把下一批腿放出去**（真无处可派须具名写"为什么派不出＋下一枚何时派"）。
- **推送只由编排者做，且只推 `origin/dev`**（网络断时⛔ 停手：本地照常 commit，恢复后统一推）。
- **子代理只 commit 不 push**；`git add` / `git commit` 两处**显式 pathspec**；⛔ `add -A` / `--amend` / `reset` / `rebase` / `stash` / `checkout .` / `clean`。
- ★**rename 类提交**：commit pathspec ⛔ 只点名新名——**旧名也要字面点名**；落账后必须现跑 `git ls-tree HEAD | grep -c issues/<N>-` 核双名＝0。
- **收件只认盘上三把尺**：①提交真在 ②件真在 ③自己复跑。**完成通知一律当未验证**。
- **翻框须非实现者判语**（裁决者≠实现者）；门禁/纪律两格可由编排者自己跑尺后翻。
- ⛔ 动 `thresholds` / `golden` / 任何为了让门变绿的放宽。
- **写面互斥**：同包只一枚写腿；`cmd/wisp` 每波只一枚；`internal/speech`/`internal/risk`/`internal/config` 等已收口包同样是共享编译面——**派腿前先现跑核有没有别腿在写**。
- **派单先存档再发**：`.scratch/wisp/probes/orch/YYYYMMDD/txt/<腿名>-dispatch.txt`。
- **临时件只建不删**；**⛔ 清理类动作**押到全部票完成（删前机主点头）。
- **筛票纪律**：「Blocked by 全 done ＋ 剩 1 格」只是候选筛、⛔ 派单依据——**派单前必须逐格读那格正文**，排除"归口 X 票 / 不适用 / 等机主拍板 / 等别的票落地"；还要排除 `Status: withdrawn` / 已撤票。
- **落地腿起手三核**：票面点名的串/行还在不在 ＋ 引文出处文件/行对不对 ＋ 红/绿是不是别人的在飞件造成的（一律 `git clone --shared` 导出树跑门禁＋正控）。

---

## 7. 平台陷阱（踩过的，必守）

- `cmd/wisp` 整包 **rc=1 是既有真窗族＋进程级截断**，⛔ 当判据；只能按名作差＋逐发登记死点；跑整包必抬 `-timeout 30m`。
- **变异读数必须与整包门禁串行**（`cmd/wisp` 有族用例运行时从盘上 parse `run.go`）。
- `.gitattributes` 要求 `*.go text eol=lf` ⇒ 脚本改文件⛔ 带 CRLF（用 `[System.IO.File]::WriteAllText` + `.Replace("\r\n","\n")`）。
- **sherpa DLL 三层口径**：什么都不给⇒`0xc0000135`（起不来）／只给 PATH⇒拿 System32 的 1.17.1（遮蔽）／**staging（二进制与三枚 DLL 同目录＋`WISP_SILERO_MODEL`/`WISP_SPIKE_MODELS`）⇒全绿**。给 sherpa 非 onnx 文件＝C++ 异常 `0xe06d7363` 打死进程、recover 不了。
- PowerShell here-string 会解析反引号 ⇒ 提交正文用 `[System.IO.File]::WriteAllText`；`-F` 写在 `--` **之前**。
- 数 markdown 表格行 `^| ` 会把表头数进去 = 必须 −1。POSIX 字符类 `[[:space:]]` 在 PowerShell/.NET 不认 ⇒ 静默数成 0。
- 本机代理把公网 DNS 劫持到 `198.18/15` ⇒ `web_fetch` 全拒、`git`/`curl` 可用。
- **看一片红先分清 `t.Logf` 与 `t.Errorf`**（行号漂移是设计允许的；token 找不到才是真红）。

---

## 8. 刚完成的事（2026-11-11 这一轮，均已推，未推 0）

1. **`font_size` 入契约**：`docs/PLAN.md:2753` 的 `[panel]` 行补上 `font_size`(12)；commit `be73612bf` 已推；`internal/config` 测试绿。
2. **wave133 三腿**（已落账 `A1701`）：
   - R：`d64d48e03` —— `cmd/wisp` + `internal/panel` + `internal/config` 接本地音频通道：voice 快照三键 `aecEnabled`/`inputDevice`/`sampleRate`（reader 指针门）＋ 写字段 `voice_aec_enabled`/`audio_input_device`（走 `internal/config/settings_local_audio.go` setter）＋ `cmd/wisp/run.go:837-844` 过期注释改真。
   - V：`b59575fba` —— 前端设置页"通用"类过期文案刷新。
   - W：裁决报告 `.scratch/wisp/probes/orch/20261110/closeout-163-388-423.md`。
3. **机主音频验收**（本对话亲耳/亲眼看，记 `A1700`）：
   - **297** 扬声器：机主"听到，很清楚" ✓（扬声器物理正常）；
   - **388 AC#2**：15% 音量播报，"正好：听得见、明显小声、不吓人" ✓；
   - **293 AC#5**：托盘静音勾"点开→勾在、再点→勾没"两形都对 ✓。
   - 余 **247 / 290 / 291 / 297 的麦克风电平**是仪器读数活，**尚未跑**（见 §10）。
4. **D 盘清理盘点**：报告 `.scratch/wisp/probes/orch/20261111/cleanup-inventory.md`。**头号垃圾＝`D:\tmp\wisp` ≈39.3GB**；合计可回收 ≈**48.7GB（D 盘）**；`.scratch\wisp\probes` 6.93GB；GOCACHE 16.45GB 但在 **E 盘**不占 D。**删前必须机主点头**。

---

## 9. 当前状态（2026-11-11 17:55 现跑；新会话开工前必须自己再跑一遍）

- **HEAD** `8e5b51890`（A1701）；**未推 0**。
- **工单池**：总 **410** / done **326** / **未结 84**。
- **tracked 脏**：仅 2 件 —— `.scratch/wisp/probes/291/mic/resident_audio_291_probe_windows_test.go`（他腿的）＋ `frontend/fixtures/composer-states.html`（render churn）。**⛔ 别碰**。
- **untracked**：约 1600 项（多为 `.scratch` 历史 + 临时件）。⛔ 别清。
- **台账最新**：`A1701`。
- **goal**：上一任会话的 goal `goal-4c8d8f64-3f2f-47fe-85f9-360b295697d4` 已**暂停**（rev57 / paused）。新会话要**新建自己的 goal**（见 §12）。
- **在飞子代理**：**0**（全部收口；上一任结束前等的那枚精简审计腿已返回并出报告）。

---

## 10. 待办 / 下一步（按优先级）

1. **音频麦克风 4 张（247 / 290 / 291 / 297 的麦克风那半）**：真机仪器读数。`cmd/wisp` 车道空闲时跑 `TestAC247LiveMicrophoneLevelsReachTheBallSeam`（`WISP_LIVE_MIC=1`，会自己放警报测麦克风回采，断言 quietA<loud、quietB≤loud）；291 四形、290 两形同理。**每波只允许一枚 cmd/wisp 写腿**。
2. **163**：R 已落 ⇒ 走验收裁决关 ＋ **追认 N/Q 两处写面**（`821f4a2d3` 动 cmd/wisp、`7bcc2ae08` 动 internal/config，均派单授权、零字节禁区未动）＋ 补 Q 收件行 ＋ 翻 AC#2/3/4/5。
3. **388**：**可直接关**（AC#2 录回半边在盘：A1688 读数 750.3→46.5 peakrms 平方律；机主亲耳也已确认）。
4. **423**：两笔小账 —— ① 把 O §三的五字段草稿写进 `SPEC-12 §5` 登记表（现 `grep` 0 命中）② 开一张**票 61 收尾票**（provider＋speed 一张；票 61 本体已 done）。
5. **D 盘清理**：报告已出；等机主点头后**按目录名**删（建议先清 `D:\tmp\wisp` 里 20261103~20261107 最老日期档试水；**⛔ 别碰** `.scratch\wisp\probes\orch\20261111` 存报告目录）。
6. **未结 84 票**：逐票读那格正文分三档（能做／归口等票／不适用要人），铺满写腿 ≥3。**⛔ 只看 `Status` 就派**——很多 "ready-for-agent" 是旧脚手架、其实早做完了。

---

## 11. 重要文档索引（落盘）

| 路径 | 内容 |
|---|---|
| `AGENTS.md`（仓库根） | 项目薄索引 + ⓪⓪⓪ 铁律 + 禁止清单 + D1–D47 表 |
| `C:\Users\swq\.dsh\AGENTS.md` | **机主全局记忆/偏好/铁律**（每次会话自动加载） |
| `docs/reports/orchestrator-handoff-20261013.md` | 前任编排者**全套交接**（角色/八组铁律/队列/派单模板/陷阱，403 行） |
| `docs/reports/HANDOVER.md` | 停车点；**§4 最新时间戳那一节**＝从哪儿接 |
| `docs/reports/pending-and-issues.md` | 真相源台账 `A##`/`R##`/`Q##`，只追加 |
| `.scratch/wisp/issues/` | 工单池（410 张；`README.md` 有派单铁律） |
| `docs/PLAN.md` / `docs/specs/SPEC-*.md` | 设计与契约（PLAN 定案最高） |
| `docs/evidence/s1/` | 裁决表（须非实现者） |
| `.scratch/wisp/probes/orch/YYYYMMDD/txt/` | 派单存档（先存档再发） |
| `.scratch/wisp/probes/orch/20261110/closeout-163-388-423.md` | 163/388/423 裁决报告 |
| `.scratch/wisp/probes/orch/20261111/cleanup-inventory.md` | D 盘清理盘点报告（可回收 ≈48.7GB） |

---

## 12. 新会话开工第一步（照做，别跳）

1. **读**：本文件 → `docs/reports/HANDOVER.md` 最新 §4 → `docs/reports/orchestrator-handoff-20261013.md` 关键节 → `.scratch/wisp/issues/README.md`。
2. **现跑地面真相**（别信本文数字）：
   - `git log --oneline -3`；`git log origin/dev..HEAD --oneline | wc -l`；`git status --porcelain`
   - 工单：`ls .scratch/wisp/issues/*.md | wc -l` / `ls .scratch/wisp/issues/*-done.md | wc -l` / 未结＝两者差
   - 台账最新：`grep -oE '^## A[0-9]+' docs/reports/pending-and-issues.md | tail`
   - 在飞：`list_agents({})` ＋ `D:\tmp\wisp\YYYYMMDD\` 临时件 mtime ＋ `tasklist | grep go.exe`
3. **验子代理路由**：`list_subagent_models({provider:"magpie"})` —— 见 §3 超级坑。
4. **建 goal**：用 `create_goal` 建常驻目标「把 Wisp 编队跑到所有票完成」（正文参照本文件 §2/§6/§9/§10；⛔ 记快照数字，权威源永远盘上）。
5. **派第一批 ≥3 写腿**（写腿优先），按 §10 优先级。
6. 每轮结束写「**你需要做什么**」。

---

## 13. 机主待拍板清单（集中摆，别散问）

1. **票 347**：解析器「一」的追认——不认可＝说「347 回退」可单独 revert。
2. **D 盘清理**：`cleanup-inventory.md` 已出，等机主点头删哪些（建议先清 `D:\tmp\wisp` 最老日期档）。
3. （其余当前无硬待拍；`font_size` 已办。）

---

## 14. 最后提醒（防漂移）

- **不要有幻觉**：任何"完成了"只认盘上三把尺；任何数字先现跑。
- **不要漂移**：别改主会话模型；别擅自清理；别跳过「你需要做什么」。
- 台账只追加；HANDOVER 时间戳最新那节是权威停车点。

