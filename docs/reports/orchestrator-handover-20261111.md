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
13. **连着两次"按你的"之后转告知式**（"我按 X 做了＋边界＋撤销口令"），不再每件事都问。
14. **他拿别的 AI 的产出回来要你审**：逐条独立取证、引用层与判断层分开报、**不照单采纳、不自动执行、落档**。
15. **验收类交互**：用 `ask_user_question` **阻塞式弹选项**、一门一门来（他在对话里直接点/听/看，如音频验收）。
16. **汇报零术语**（有术语配翻译对照）；**定性词不许重于证据**；涉安全/攻击字样自带三行（现象在哪／有无本机被入侵证据／最坏后果形状），不够就写"未定性"。

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
>
> ⚠ 若在旧 goal 文本／旧文档／旧台账里看到 `auto/review-and-audit`，那是**笔误**；正确 id＝`group/group-review-and-audit`（显示名 `group/review-and-audit · routing group`）。

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
| cgo 工具链 | `E:\work\base\msys64\mingw64\bin`（gcc 16.2.0）；`make` 在 `E:\work\base\msys64\usr\bin` |
| 一键构建 | `scripts/build.ps1 [-Env dev|prod]`（fetch-deps→build→DLL 同目录→SHA256SUMS）；`docs/BUILD.md` 已冻结 |
| 自托管 runner | `wisp-selfhosted-01`（标签 self-hosted/Windows/X64/wisp-slo；登录自启 HKCU Run → `E:\work\base\start-wisp-runner.vbs`） |
| Docker | Desktop 可用；compose 的 model-mirror（18081）＋fixtures（18082/18083）＋mock-llm（18080） |
| minisign dev 密钥 | `E:\work\base\wisp-minisign\`（盘外，公钥硬编码 internal/buildinfo；生产仪式归 S8） |
| 模型缓存 | `third_party/spike-models/`（git 忽略；URL/SHA256 见 docs/evidence/s0/json/） |

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
- ⚠ **`go build ./...` 当前 rc=1 来自未跟踪/被忽略的 `build/mut_rvw.go`**（别人遗留的 91KB 碎片）；`go build ./cmd/... ./internal/...` rc=0。⛔ 别删/别改 `build/mut_rvw.go`（不在你账上）。
- `go test ./cmd/wisp` 重包**同一时刻只跑一个**；跑前先 `tasklist` 看 wisp.exe/balldebug.exe 为 0。
- **推送前做凭据扫描**（仓是 **PUBLIC**，⛔ 无外泄；只引用变量名不引值）。

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
| `docs/reports/injection-timeline.md` | 伪授权形状的世代记录（判"是不是授权"看这份的判据） |
| `C:\Users\swq\.qoder-cn\projects\D--work-workspace-projects-plans-Wisp\memory\` ＋ `C:\Users\swq\.qoder-cn\memory\` | 历史 memory 文件（**次要源**，可当线索；权威永远是盘上现状＋台账＋停车点＋现跑读数） |

---

## 12. 新会话开工第一步（照做，别跳）

1. **读**：本文件 → `docs/reports/HANDOVER.md` 最新 §4 → `docs/reports/orchestrator-handoff-20261013.md` 关键节 → `.scratch/wisp/issues/README.md`。
2. **现跑地面真相**（别信本文数字）：
   - `git log --oneline -3`；`git log origin/dev..HEAD --oneline | wc -l`；`git status --porcelain`
   - 工单：`ls .scratch/wisp/issues/*.md | wc -l` / `ls .scratch/wisp/issues/*-done.md | wc -l` / 未结＝两者差
   - 台账最新：`grep -oE '^## A[0-9]+' docs/reports/pending-and-issues.md | tail`
   - 在飞：`list_agents({})` ＋ `D:\tmp\wisp\YYYYMMDD\` 临时件 mtime ＋ `tasklist | grep go.exe`
3. **验子代理路由**：`list_subagent_models({provider:"magpie"})` —— 见 §3 超级坑。
4. **建 goal（常驻，不是单轮清单）**：用 `create_goal` 建目标「把 Wisp 编队跑到所有票完成」，`max_goal_rounds` 设大（如 400）；正文参照本文件 §2/§6/§9/§10/§15（⛔ 记快照数字，权威源永远盘上）。每轮收件后继续自动续跑，直到所有票完成。
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


---

## 15. 硬规则全表（AGENTS.md / 工单 README 抄录，新会话必守）

> 本节是对 `AGENTS.md`（项目根）、`C:\Users\swq\.dsh\AGENTS.md`（机主全局）、`.scratch/wisp/issues/README.md`、`docs/reports/orchestrator-handoff-20261013.md` §3 的**压缩抄录**；与原文件冲突时以原文为准。

### 15.1 进场三句（项目 AGENTS.md §0）
1. **未定义即停**（D22 闸门③）：方案没覆盖的情况 → 停下来问，**不得自行假设**。
2. **改契约＝人工批准**：`D1–D47` / `C1–C32` / `R1–R9` / D43 转移表；agent 单方面改＝跑歪模式 #1，**对抗验收直接判失败**。
3. **裁决者≠实现者**：缺口审计与对抗验收必须由另一个 agent 做。

### 15.2 禁止清单（AGENTS.md §1，CI 自动扫 tools/d22scan）
- 不许改：`D1–D47`/`C1–C32`/`R1–R9`/D43 转移表（改契约须人工批准）；**SLO 阈值 / golden / `internal/observe/thresholds.go` 一字节不许动**；`tools/d22scan/allowlist.txt`、`.github/**`（含 ci.yml）⛔ 动。
- **三枚冻结测试件**：`internal/panel/tokens_fourway_test.go`、`internal/panel/l2_grant_boundary_test.go`、`internal/perm/ticket90_persist_test.go`。
- **禁止代码形状**（CI 自动扫）：裸 `go func(` 无 owner/recover；在 `risk.PathResolver` 之外用 `filepath.Clean|Abs` 做文件系统决策；明文 API 密钥；用"墙钟时间差"实现超时；从镜像站取哈希；由面板侧来源的 L2「允许」；把宿主内部 artifacts 写入实现成受门控的 Tool；**UI 代码任何 emoji**（仪器扫 `U+1F000–1FAFF`/`2200–22FF`/`2600–27BF`/`2B00–2BFF`/`FE0F`，**注释豁免、字符串不豁免**；`PLAN.md` 的"绝对禁止"清单比仪器更严（含 `→`），**判文案取严的那份**）。
- **`DEFERRED(D-xx)` 代码标记必须与 `SPEC-12 §5` 登记表 1:1 双向对得上**。
- **`cmd/wisp/config_readers_255.go` 按行号钉住 11 枚被引文件**（含 3 枚 `internal/**`）：被引文件"注释变高"就会打红（红句 `the evidence drifted, so re-adjudicate this row`）；**往被引件追加任何行前先跑它**。

### 15.3 测试注入面（AGENTS.md §1.3）
- 只在接缝注入：`C8 AudioSource`（wav）· `C5 LlmProvider`（golden SSE）· `C17 PanelBridge` · CLI `wisp run`。
- **不许"用 mock 代替真的"来假报完成**；**never weaken SLO thresholds**。

### 15.4 票与角色纪律（orchestrator-handoff §3.3 + issues README 状态机）
- 票状态机：`ready-for-agent` → `in-progress`（Claimed; ≥1 Progress-log）→ `review`（做完待对抗验收）→ `done`（全框勾 + 改名 `-done` + 标题 `(DONE ✅)` + 更新索引）；缺决定就 `blocked`（把开放问题写进日志）。
- **开工前**：置 `Status: in-progress`、填 `Claimed by`/`Last update`、追加 1 条 Progress-log，**同一步 commit**（子代理只 commit 不 push）。
- **翻框/改状态只有编排者能做**：`- [ ]`→`- [x]`、改 `Status:`、`git mv` 加 `-done` 后缀，**腿一律禁止**（派单里每次都写）。
- **`in-progress` 超 24h 无日志**：下一任先按盘上状态核对日志再续，**绝不重置已勾框**，追加 `resumed` 行。
- **置 done 前，验收报告必须含一张与 AC 编号 1:1 的裁决表**（每行＝一条 AC 的 PASS/FAIL/PARTIAL＋判定依据 `file:line` 或实跑输出）；**缺行即 FAIL**。
- **腿的边界**：⛔ 翻 AC 框 / ⛔ 动 `Status:` / ⛔ 加 `-done` / ⛔ 改票面原句；只许**在票尾追加 1 条** Progress-log bullet（格式 `- [YYYY-MM-DD HH:MM:SS +08] agent=… did=… next=…`，钟点用 `date` stdout 插值）。
- **结案＝`git mv` 改名加 `-done`**，commit pathspec **同时含旧名＋新名**。
- **`-done` = 零残余证明**：结案前逐格核完；历史上有"已 -done 却仍有未勾框"的乙类问题——**别批量撤名、逐张核**。
- **每格验收靠非实现者**：验收腿产出裁决表 → `docs/evidence/s1/`（逐格判语＋每句带可复跑尺）→ 编排者复跑关键读数后翻框。
- **验收腿的活是攻"判据真假"**：恒真判据／空集通过／"诚实的降级检测器"都算没牙；每枚新仪器要自带**红形演示**（把判据改坏必须红）。
- **新立票**：入台账（`A##` 指名新票号），票面写清 AC＋禁列＋出处；**票名 ≤100 字符**（开完新票、提第一笔 commit 前跑 `sh scripts/check-path-length-budget.sh`）。
- **同票 ≤2 子代理**（按票内 `Parallel slots`）；需安静测量的票独占 1 个代理；只读检索/审计可 10+；>1 个写码代理必须 worktree 隔离（`third_party` 用 junction 复用）。

### 15.5 每片完成后的五件事（AGENTS.md §3 / SPEC-12 §4.3）
1. 重跑**缺口审计**（审计者 ≠ 实现者）
2. **失败预演**（SPEC-10 §7）
3. **对抗验收报告**（SPEC-10 §8 七项检查）
4. **场景↔切片双向核对**（SPEC-00 §4 矩阵、SPEC-12 §3）
5. **推迟项登记更新**（新增推迟必须五字段齐全，SPEC-12 §5 + M0 要求①）

### 15.6 未定义即停：已知待定案项（AGENTS.md §2）
碰到下列任一，**停手上报、不要自己填**：`AwaitingApproval` 遇系统挂起（倾向作废判拒绝，S7 切片卡定案）；`web.search` 实现路径（抓结果页 vs API，S3 前）；Go module 组织名 + P10 命名核查（S1 建仓前，**阻塞**）；便携模式 × DPAPI（SPEC-02 §6 建议，S1 定案）；`C24 GojaHostAPI` 初始集与 `C17` 方法白名单定稿（S7/S5 切片卡批准）；WebView2 冷拉起 >2s → 重评 L2 卡是否回原生（P11，S0）。

### 15.7 前端与 design 归属（⛔ 最常踩，orchestrator-handoff §3.4）
- `frontend/**`、`design/**` 归**机主的外部 agent**；本编队 **⛔ 写页面、⛔ 替它改、⛔ 还原/提交 `design/` 的删除态**。判据**按目录不按票号**：要写的文件只要有一枚落在 `frontend/` 下就**整段停手上报**（⛔ "只动一小角"）。
- 界面需求**唯一出口**：写进**票面/证据件**，由**机主自己**带给他的 agent——你只告诉他"件在哪"。
- **⛔ 用 builtin 跨会话工具**（`list_chat_sessions`/`read_chat_session`/`send_message_to_chat_session`/…）联系、查看、指令**任何**会话；**⛔ 在汇报里提"要不要我转过去"——连问句都不许出现**（机主为此发火三次，原话「你转给谁？？？说了八百回了都」）。
- 门禁"读" frontend 不违规（扫描器 walk 它）；**写**才违规。判"我们动没动过前端/design"：逐笔 `git show --name-only <commit>` 看**路径**（**⛔ 按 author**——所有会话共用同一 git 身份）。
- `internal/panel/**`、`cmd/wisp/**` 的 **Go 侧接线不冻结**，属本编队。

### 15.8 凭据与安全（orchestrator-handoff §3.5）
- API key／凭据 **⛔ 进对话框**；只引用变量名，或写"由界面录入"（机主在界面自助录入；他明确否了 CLI 录 key）。
- 配置写路径有白名单（`config.set` 7 枚字段全在 `[llm]`＋票 313 第二张表 6 枚）；**`credentialValue` 是只写键**；设置写**不能新建行**（必须先有段）。
- 自动模式下改生产凭据会被宿主硬拦（走业务 API 或让机主本人操作）。

### 15.9 编队与环境补充（orchestrator-handoff §3.7）
- **一波里 `cmd/wisp/**` 只交一枚**；起手机械三问：① `git status --porcelain -- cmd/wisp` ② 该包"颜色"归谁 ③ 非空 ⇒ 停手报回。
- `GOTMPDIR`/`WISP_TEST_GOT_ROOT` 若要设**只许 C: 卷**（换卷会让既有门变色——票 311；守卫"换卷即拒跑"已在装）。
- 跑读数前：`tasklist //FI "IMAGENAME eq wisp.exe"` 与 `balldebug.exe` ==0；⛔ 裸跑 `sqlite3`（挂死）；⛔ 读二进制文件；grep 一律加 `--include` 与 `| head` 设界。
- **件纪律**：⛔ 件名带 `§`；⛔ 超长名；⛔ `.out`（静默不收）；证据件每枚**自落 `rc=N` 行**（0 字节＝没交，除非同格同伴件里有那行）；⛔ 在 `probes/**` 下建 `.go`。
- **测量台面**：改前⇄改后必须**同台面**；"改前那发"⛔ 用共享工作树跑 ⇒ `git archive` 仓外导出树，六件套 `go.mod go.sum internal cmd tools scripts` 一枚⛔ 少，导出后先数 `.go` 枚数证非空；一树一 `mktemp -d`。

### 15.10 证据纪律（防漂移核心，orchestrator-handoff §3.8）
- 三把尺收件：① `git log -1 <sha>` 提交真在 ② `ls` 证据件真在盘上 ③ **你自己复跑**关键读数。腿自报的数≠数。
- **完成通知一律当未验证**：通知里的提交号可以整枚不存在（幽灵交件七形：假提交号／回执"创建成功"但盘上没有（判据 `wc -c`）／注释里点名不存在的测试／"我已写"其实没写…）。
- 引任何**行号**前现跑复认（行号会腐烂）；优先**内容锚**（可 grep 的句＋"偏离即重锚"自证）。
- 写腿在飞时**树会呼吸**（它能临时换回文件做对照）⇒ 锚点一律引 **HEAD 对象层**（`git grep … HEAD`），⛔ 现读工作树当锚。
- 名册对拉先写"**这把尺的射程**"；`grep -c` 命中 0 ⇒ rc=1 会截断 `&&`（写 `N=$(… || true)`）；一把尺报 0 前**先看 stderr**（假 0 陷阱）。
- 变异台件（改坏判据必须红）：用 `-overlay` 或仓外备份＋`sha256sum` 还原证明；**⛔ 裸改跟踪文件不还原**；突变后必须跑"还原→回绿"证明。

### 15.11 盘符根目录零污染（issues README 规矩 10，⛔ 机主点名批评过）
- ⛔ 任何导出树／`GOTMPDIR`／`TMP`/`TEMP` 落点／探针工作目录／缓存／快照，**一律不许建在盘符根目录下**（`D:\W334R1TreeAfter`、`D:\wisp228r3b` 那一族＝违规）。
- **唯一允许位置＝`D:\tmp\wisp\<YYYYMMDD>\<腿名短码>\...`**。
- 旧前提"必须放盘符根才能生成 8.3 短名"**已作废**（`D:\tmp\wisp\...` 同样生成短名）。
- 本条与"临时件只建不删"**叠加**：新落点要规整，但**谁都不许删存量**（D 根旧树由机主处置，编排者只出清单与证据引用尺）。
- **规矩 8 细节**：删除类命令（`rm`/`rmdir`/`Remove-Item`）每发一次就弹一次授权窗——所以收尾**不要清理自己的工作目录**，把清单（路径＋体积＋有没有被证据文件引用）写进交件；**被证据引用的快照目录任何人都不许删**；批量清理由编排者一次做完。

### 15.12 高危陷阱补充（orchestrator-handoff §8）
- **票 314：`cmd/wisp` 测试进程死亡区**——先认它，再怀疑自己的码。
- **行号钉与"注释变高"**：`cmd/wisp/config_readers_255.go` 会因被引件追加行而红；追加前先跑。
- **"通知 ≠ 交件"**：幽灵七形，见 §15.10。
- **树会呼吸**：写腿在飞时不要拿工作树当锚。
- **名册尺的射程与假 0**：见 §15.10。
- **GOTMPDIR 换卷变色**：只许 C: 卷。
- **共享编译面与前端相关门**：写面互斥，见 §6。
- **宿主弹窗与"系统发生异常"**：主动解释"谁渲染的/谁请求的/动的是不是仓库、拒绝了会怎样"；**弹窗数量常是编队自己造的**。

### 15.13 授权与撤销口令全表（orchestrator-handoff §10）
| 事项 | 状态 | 边界/口令 |
|---|---|---|
| 推送 `origin/dev` | **已授权**（A813, durable） | ⛔ force；⛔ 裸推（会去 cnb）；推前核 fast-forward `git rev-list --count HEAD..origin/dev`==0；仓 **PUBLIC**；推送自启本机 CI（抢 CPU） |
| 合并前端工作树 | 已用（A873，"能合并就合并"） | 再有改动要再并一次；归因按 pathspec |
| 清 C 盘遗留 | 已用（A870） | 只删可重建；证据面不动 |
| 294/295 串行 | 现行（面对空后可改并发） | 撤销口令「**294 295 改并发**」 |
| 票 145 `internal/panel` 局部解冻 | 在效（09-26） | 「**145 别动**」；三枚冻结测试件仍冻 |
| 299 换供给形状 | 已落地 | 「**回退 299 的换供给形状**」 |
| 313 落地前四裁定 | 已落（A877） | 语音链两枚可写字段开账（实时段走β）；显式建行方向开账；默认 stepfun＝界面预选＋只做读面；⛔ 新 C17；⛔ 编模型名 |
| 票 312（扫本地配置） | 只立不派 | 顺序：299→前端树→313→**312 最后** |
| 前端归属 | 现行 | ⛔ 本编队写页面／⛔ 跨会话工具／⛔ "要不要我转过去" |
| 面板窗口快捷键 | 现行 | `Ctrl+Alt+P` 叫出/收起（给机主的演示键） |
- 其余裁定/撤销口令见台账 `A##`/`Q##`/`R##`。

### 15.14 工作模式 / 停车点 / 门禁链（orchestrator-handoff §1）
- **编排者＝调度者**：生产码由落地腿写；你只在两类场景亲自动手：⑴**死腿遗产代提**（腿死了、码/证据在盘上没提交，你替它 commit 并标〔未验证半成品〕）⑵机主直接点名由你做的活（如清 C 盘遗留、构建出厂 exe）。
- **子代理一律后台派发**，绝不前台等（09-29 前台等待卡死主会话十小时，机主发过火）；派发后不要轮询，完成通知会自动来（通知只是"有东西落地了"，数字以盘上现跑为准）。
- **波次制**：派 N 枚 → 收件（三把尺）→ 翻勾/记账/落停车点 → 按 CPU 补位 → 循环。
- **腿的死亡形状**：150 轮帽死亡／断线死亡／基础设施故障死亡；死腿按进度重派 `r1b`/`r1c`（**带加固规则**：针对上一程死因补硬约束；零遗产才直接重派，有遗产先代提）。
- **停车点格式**：`docs/reports/HANDOVER.md` §4 加新节 `### 4.0<下个字母> 停车点（**YYYY-MM-DD HH:MM +08**）｜编排者｜<一句话标题>`，正文大白话，结尾必有「你需要做什么」（没有就写"没有"）。
- **派单存档**：`.scratch/wisp/probes/orch/YYYYMMDD/txt/<腿名>-dispatch.txt`（主文）＋ `<腿名>-send.txt`（发送时刻补记：只记"从写派单到发送之间变了的前提"，冲突时以 send 为准）；**先存档再发**。
- **门禁链**：`go build ./...`；`sh scripts/d22scan.sh`（ban 扫描含零 emoji）；`$(go env GOPATH)/bin/gofumpt.exe -l`；`go test ./...`（真窗测试要桌面）；`sh scripts/check-path-length-budget.sh`；CI＝`.github/workflows/ci.yml`（push 触发，跑在本机 self-hosted runner 上，会抢 CPU）。

### 15.15 Git 附加纪律（orchestrator-handoff §3.1）
- ⛔ `--no-verify`；**已推送的历史不改写**，要更正就追加新 commit；⛔ 在仓库目录内建 worktree 或 checkout。
- **共享工作树**：别人的脏件不是你的账（当前 2 件见 §9），⛔ 顺手提交、⛔ 还原。
- 凡"我已写／已提交／已发出"，**落笔前后各跑一把尺**（`git log -1 <sha>`／`ls -l <件>`／现跑读数）。

---

## 16. 新会话启动提示词（机主直接粘贴给新会话）

```text
你是 Wisp 项目的编排者（orchestrator）。这是接手会话，前一任已把所有交接写在磁盘上；你没有历史记忆，先读文档再动手，不许凭空猜。

【第一步（必做，按序）】
1. 读 D:\work\workspace\projects plans\Wisp\docs\reports\orchestrator-handover-20261111.md（完整交接，16 节，必读全；§15 是硬规则全表，务必逐条读）。
2. 读 docs/reports/HANDOVER.md 最新 §4.0ga 那一节。
3. 读 docs/reports/orchestrator-handoff-20261013.md（前任全套交接：角色/铁律/派单模板/陷阱）。
4. 读 .scratch/wisp/issues/README.md（派单铁律）。
5. 现跑地面真相（⛔ 别信交接文档里的数字，一律现跑）：git log --oneline -3；git log origin/dev..HEAD --oneline | wc -l；git status --porcelain；工单总数/done/未结；grep -oE '^## A[0-9]+' docs/reports/pending-and-issues.md | tail；list_agents({})；D:\tmp\wisp\ 临时件 mtime。
6. 验子代理路由：list_subagent_models({provider:"magpie"}) —— 必须出现 group/group-review-and-audit 和 auto/glm-5.3-flash 才能按分工派；如果只有 deepseek，立刻如实告诉机主（别静默降级、别改主模型）。
7. 用 create_goal 建常驻目标：「以 Wisp 编排者身份把编队跑到所有票完成」（正文照抄交接文档 §2/§6/§9/§10/§15 要点；⛔ 记快照数字）。goal 是常驻目标（max_goal_rounds 设大，如 400），每轮继续自动续跑，直到所有票完成。

【你是谁 / 铁律（违反＝违令）】
- 你是编排者：不亲手写票，派「腿」（子代理）、收件核盘上三把尺、维护台账、只推 origin/dev。
- 主会话模型＝DeepSeek-V4.1-Flash，永远不换（指挥大脑）。
- 模型分工：审计/测试/审查类子代理 → workflow 的 agent(prompt,{provider:"magpie", model:"group/group-review-and-audit"})；其余干票子代理 → provider:"magpie", model:"auto/glm-5.3-flash"。⛔ 用 subagent/spawn_teammate 换模型（它们无 model 参数，只能走 workflow）。
- ⛔ 绝不阻塞主会话：所有 workflow 一律 run_in_background: true，拿 job id 立刻返回；⛔ job_output(wait:true) 长等、⛔ wait_agent 干等、⛔ sleep 轮询。
- 常驻运行：写腿 ≥3 在飞、总派发 ≤7、读腿不占写额、CPU 三发全 ≤70% 即补位；每收一腿落账同轮必须放下一批。
- 子代理只 commit 不 push；git add/commit 两处显式 pathspec；⛔ add -A/--amend/reset/rebase/stash/checkout ./clean；rename 提交新旧名都点名。
- 收件只认盘上三把尺（提交真在/件真在/自己复跑），完成通知一律当未验证；翻框须非实现者。
- ⛔ 动 thresholds/golden；写面互斥（同包一写腿，cmd/wisp 每波一枚）；派单先存档再发；临时件只建不删；清理类动作删前机主点头。
- 筛票：派单前逐格读那格正文（排除归口/不适用/等票/withdrawn）；落地腿起手三核。
- 最高铁律 ⓪⓪⓪：技术选型按「主流 harness/最有价值」，⛔ 用「风险小/动作小」降级；技术选择自己拍+写理由，只有①改契约文字②花钱③不可逆动作(推/发布/删)才摆机主。
- 禁止清单／未定义即停／每片五件事／前端归属／凭据白名单／盘根零污染：见交接文档 §15 硬规则全表，逐条遵守。

【机主偏好（他是机主，不写代码）】
- 每次总结结尾必须写「你需要做什么」：具体到看哪块屏/按哪几个键。
- 大白话+数字；诚实高于好看；宁缺毋造；问「做了没」给可核命令+读数。
- ≥2 待拍板项出编号清单（含「不做」栏+撤销口令），集中一次摆。
- 写腿优先，别只搞只读子代理。

【当前状态（交接时快照，动手前必须自己重跑）】
HEAD 8ab61ce52；未推 0；工单总 410 / done 326 / 未结 84；台账最新 A1702；在飞 0；goal 已暂停（前一任 rev57）。
第一优先：①音频麦克风 4 张 247/290/291/297 真机仪器读数（cmd/wisp 空闲时跑 WISP_LIVE_MIC=1 go test ... -run TestAC247LiveMicrophoneLevelsReachTheBallSeam，每波只一枚 cmd/wisp）②163 裁决+追认 N/Q ③388 可关 ④423 两笔 ⑤未结 84 票逐格读铺满写腿≥3 ⑥D 盘清理（报告 .scratch/wisp/probes/orch/20261111/cleanup-inventory.md，可回收≈48.7GB，删前机主点头）。

【开工】读完文档+跑完地面真相+验完路由后，立即建 goal、派第一批 ≥3 写腿，然后持续常驻干活，别等我。
```

