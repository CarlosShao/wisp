# HANDOVER — Wisp 实施交接文档（给新会话的新 agent）

> 写于 2026-09-19/20（+0800）。**2026-09-20 23:20 增量更新（§3/§4 已刷新）**。
> ⚠ 原来那句"用户平台配额 2026-09-20 09:00 到期"**已作废**：那是 ZCode 免费套餐时代的约束，
> 用户现在是 **Qoder 付费**，瓶颈已不在平台配额（车队上限见偏好记忆：写码稳态 3、只读 10+、测量类独占 1）。
> 阅读顺序：§1 使命与铁律 → §2 环境 → §3 完成度 → §4 在途状态 → §5 下一步 → §6 技术事实 → §7 事故手册 → §8 待人项。

## 0. 项目是什么 + 前因后果在哪里（新 agent 先读这节）

**Wisp（一缕）**：Windows 桌面语音 Agent 悬浮助手——常驻悬浮球（空闲进程树 ≤25MB、默认不开麦克风），
语音/文字驱动 30+ 内置工具与插件干活，云端 realtime 增强陪聊。开发方是 AI Agent（你），
人类 owner 只做拍板与验收。

**前因后果的文档分层**（HANDOVER 只带操作状态；"为什么"按需查下表，越左越稳定）：

| 层 | 文件 | 回答什么问题 |
|---|---|---|
| 为什么这样设计 | `docs/PLAN.md`（D1–D47 决策 + §1 Context + 打磨过程） | 为什么 Go 不用 Rust？为什么半双工？为什么 realtime 只做陪聊增强？每条决策含取舍与被否方案 |
| 契约（不可擅改） | `docs/PLAN.md` §3 的 C1–C32 + `docs/specs/` | 界面长什么样、数据结构、门控规则——改动需人工批准（D22） |
| 怎么做 | `docs/specs/SPEC-00..12` | 每个领域的落地规格（骨架/存储/语音/安全/UI/测试/构建） |
| 做到哪了 | `.scratch/wisp/issues/`（61 票，`-done` 后缀=已完成） | 本票做什么/不做什么/验收标准/进度日志 |
| 做过的证明 | `docs/evidence/` | 每票对抗验收报告 + 截图/JSON 证据 |
| 事故与待人 | `docs/reports/pending-and-issues.md` | 平台故障记录、需要 owner 的事项（H1–H6）、裁定（R1–R2） |

**30 分钟上手路径**：读本文件（10 分钟）→ 读 `.scratch/wisp/issues/README.md` 索引与铁律（5 分钟）→
按 §5 行动清单打开你要做的票 + 它 Spec refs 指向的 spec 小节（10 分钟）→ `git log --oneline -20`
扫一遍近期演进（5 分钟）→ 开工。**不要通读 PLAN.md**（3542 行）——它是字典不是小说，
按票据里的（Dxx/Cxx）引用按需查。

## 1. 使命与治理铁律

- **目标**：按 `.scratch/wisp/issues/README.md` 的票据顺序完成 Wisp（Windows 语音 Agent 悬浮助手）全部 61 张票。
- **权威文档**：`docs/PLAN.md`（D1–D47 决策 + C1–C32 契约，冻结不得擅改）→ `docs/specs/SPEC-*.md`（13 份落地规格）→ `.scratch/wisp/issues/`（61 张原子票）。
- **票据铁律**（`.scratch/wisp/issues/README.md`）：开工先置 `Status: in-progress` + 填 Claimed by + Progress log 追加（newest last）+ commit push；**每完成一个最小单元立即 commit+push**；完成 → 全验收框打勾 → Status: done → **文件重命名加 `-done` 后缀 + 标题加 (DONE ✅)** → 更新索引。
- **对抗验收**：每票实现完成后由**非实现者**验收（子代理或编排者亲自——独立性=实现者≠验收者）；FAIL → 最小修复集发回原实现代理 → 复核 PASS 才置 done。报告存 `docs/evidence/`。
- **commit 纪律**：**永远显式列路径，禁止 `git add -A`**（两次吞并行 WIP 事故教训）；Conventional Commits；双远程 push（origin 偶发 TLS 失败重试≤5 次，cnb 稳定）。
- **并发纪律**：子代理并发 2 为底线、3 为试验位（3 曾触发平台 captcha/quota 故障，回落 2）；子代理闪断（Captcha timed out / exceed quota）→ 直接重派续传，**不丢弃不降级**。
- **D22 未定义即停**：方案未覆盖 → 停下来问用户；DEFERRED 项不得擅改。

## 2. 环境（全部已就绪并实测）

| 项 | 位置/值 |
|---|---|
| 仓库 | `D:\work\workspace\projects plans\Wisp`（git，dev 分支为工作分支） |
| 远程 | origin=github.com/CarlosShao/wisp，cnb=cnb.cool/CarlosShao/wisp（凭据已存，gh CLI 已登录） |
| Go | `D:\work\base\go\bin\go.exe`（1.27.1，不在 PATH，自行 export）；GOPATH=`D:\work\base\gopath`；GOPROXY=goproxy.cn |
| cgo 工具链 | `E:\work\base\msys64\mingw64\bin`（gcc 16.2.0）；make 在 `E:\work\base\msys64\usr\bin` |
| 构建 | `scripts/build.ps1 [-Env dev|prod]`（一键：fetch-deps→build→DLL 同目录→SHA256SUMS）；`docs/BUILD.md` 已冻结 |
| 自托管 runner | wisp-selfhosted-01（标签 self-hosted/Windows/X64/wisp-slo），登录自启（HKCU Run → E:\work\base\start-wisp-runner.vbs） |
| Docker | Desktop 可用；compose 的 model-mirror（18081）+ fixtures（18082/18083）票 14 已建，mock-llm（18080）票 08 接线 |
| minisign dev 密钥 | `E:\work\base\wisp-minisign\`（盘外，公钥硬编码 internal/buildinfo；生产密钥仪式归 S8） |
| 真麦克风 | Realtek，LIVE 冒烟：`WISP_LIVE_MIC=1 go test ./internal/audio/ -run TestLiveWasapiSmoke` |
| 模型缓存 | `third_party/spike-models/`（git 忽略；URL/SHA256 见 docs/evidence/s0/json/） |

## 3. 完成度（下表是 **2026-09-20 08:4x 的快照，仅 13/61**；**今天 23:19 已是 18/65**，新增哪四张见 §4）

| 票 | 内容 | 证据 |
|---|---|---|
| 01 ✅ | 构建链 + BUILD.md 冻结 | docs/evidence/s0/01-* |
| 02 ✅ | S0 spike（Path Y 判定、五基线、goja/WebView2/模型常驻实测） | docs/evidence/s0/02-*、docs/SLO.md |
| 03 ✅ | 19 包骨架 + D37 错误模型 + goroutine 注册表 + C11/C30 | docs/evidence/s0/03-* |
| 04 ✅ | SQLite 核心（8 表 DDL、WAL、db-writer、保留期）+ schema v2（provider_health，票 09） | docs/evidence/s0/04-* |
| 05 ✅ | config 全量模型（catalog v2、三档热加载、🔒方向引擎、迁移） | docs/evidence/s1/05-* |
| 06 ✅ | SecretStore(DPAPI) + WISP_ENV 分叉 + 便携/P13 | docs/evidence/s0/06-* |
| 07 ✅ | 悬浮球窗口 + D2D 渲染 + 20 态状态机（D43 全 42 行）+ 热键/托盘/多显示器 | docs/evidence/s1/07-*、ball-states/ 截图 |
| 09 ✅ | LLM seam（C5/C6/C7 全集）+ OpenAI Chat adapter + golden 两跑器 + mockllm + text_chain/令牌桶 + provider_health 原语 | docs/evidence/s1/09-* |
| 13 ✅ | 音频采集（WASAPI 真机验证、热插拔、半双工门、重采样） | docs/evidence/s2/13-* |
| 14 ✅ | 模型分发（C29 签名清单 6 条目、minisign 手写验证、断点续传、compose model-mirror；matcha P3 阻分发不阻自用——裁定见 reports R1） | docs/evidence/s2/14-* |
| 08 ✅ | 可观测性（六态×七指标采样器——**门禁口径=私有工作集**（commit 假红裁定）、slo-check.ps1 泄漏 fixture、CI 五 job 矩阵 lint/test-core/test-windows/slo-smoke/slo-full(wisp-slo)、compose mock-llm 接线） | docs/evidence/s1/08-* |
| 17 ✅ | C19 RiskAssessor（R1–R9 全实现 + R9 fail-closed + 融合 + R4 会话覆盖阻断；R2/R3/R4 接口待 18/19/20 接线） | docs/evidence/s1/17-* |
| 18 ✅ | C26 PathResolver（junction/8.3/UNC/`\?\` 红队四连真产物全拒）+ A/B 黑名单 | docs/evidence/s1/18-* |

**（23:20 更新）** S0 完成；S2 前置 13✅+14✅。**10 ✅**、**19 ✅**；
**票 12 = S1 gate，仍在 review**：AC#1/3/4/5/6 已勾，**未勾的四条是 S1 的全部残口**——
**AC#2** 空闲 SLO（CPU 行按 A15 判"仪器未定义"，等票 66）、
**AC#7** emoji/tokens（A22/A23/A24 把前提证伪了，等票 67 AC#3 + 票 68 AC#1，另加 A24-D4＝票 69）、
**AC#8** owner 视觉签收、
**新加的「3s 回落」条**（A25：S1 四项判据里**唯一从未被量过**的一条）。
**票 20** review 中（我今天补回被删的 `Needs` 逐条断言＝A17，另 A18 真 kill 残留待办）；
**票 21** 段 1 已对抗验收（判定内核干净、**无生产应答路径**＝A19，按 R14 降为"S1 可接受、S3 必做"），
**段 2 未开工**。S3 安全地基 17✅+18✅ 不变。

## 4. 在途状态（**2026-09-21 09:38 版，取代下面 23:5x 版**；旧版整段保留在文末供追溯）

**done 计数：20 张**（判据永远是 `ls .scratch/wisp/issues/*-done.md | wc -l`，09:49 实测 = 20）。
本轮（09:0x–09:50）**归档两张票**：**69**、**67**（4/4，AC#4 按票面要求在 AC#3 落地后**重跑过**；
整仓 `go vet` 我故意没跑，理由写在裁决表里），外加**两处跨票解锁**：
票 20 的第 5 框（桥层真 junction/8.3）与票 67 留下的这条——
票 12 的 AC#7（它的两个解除条件是票 67 AC#3 + 票 68 AC#1，现均成立）已勾选，
措辞里明写"这一框只覆盖 Go 面的字面量与注释，**非 Go 的 UI 文案没有门**（`frontend/` 树根本不存在）"。
本轮（09:0x–09:38）收口的三张：

- **票 69 已 `-done`**：C21 表的双向机器检查。裁决表 `docs/evidence/s1/69-adversarial-acceptance.md`。
  **AC#2 我没有采信代理日志**，自己在独立树里把 `tokens.go:361 DockTriggerPx 16→17` 重做了一遍——
  唯一变红的是那条几何用例且它点名表行 155。四条"本票没证明的"全部带票号转给票 74。
- **票 20 第 5 框已勾**（桥层真 junction / 真 8.3，`d63bc49`）+ **A18 特征化落地**（`761447f`）。
  本票仍有 **2 框未勾**（含第 7 框：它要的 `ErrReparseDenied` 原因枚举**不存在**，补它要碰冻结的
  `internal/risk/pathresolver*.go` ⇒ **D22，不许代理自作**）。**保持 in-progress 是诚实状态。**
- **票 67 的 AC#3 判据②**（`bcf44d6`）：ban #8 覆盖面真扩到 `internal/`+`cmd/`，`frontend/` 那个
  **实际走 0 个文件的死作用域**删掉了。

**新增两张票（09:33，都是代理交回的真缺陷，不是重构）**：
**票 73** 孤儿暂存文件清扫（真 `taskkill /F` 每次留 1 个 `.wisp-tmp-*` 且没人扫 → owner 项 **Q-16**）；
**票 74** C21 表必须对齐**真正驱动像素的常量**（`tokens.go` 声明 `SwimLevelGain`/`SpinLevelGain` 而
零消费者，真驱动是 `liquid.go` 的 `Spin*RadPerS` = **A33 同族**）。

**在飞**：票 70-c（`internal/observe`、`tools/d22scan`、`go.mod`）、票 67b（票 67 票面）、
票 73（`internal/tools`）、票 74（`internal/ball` + C21 表）。

**⚠ 一条操作纪律，是今天用一次自我纠错换来的**：**不要在仓库根建嵌套 worktree 做变异**。
它会并发代理的文件遍历（`d22scan -root .` 走 `internal/`+`cmd/`）**多算一整份源码**，
正打在票 70 那条"生产文件数下限否则扫描空转"的守卫上；而且 `git worktree remove --force` 可能
`Permission denied`（句柄未放），中间会留一段污染窗口。变异照代理的做法：**改→跑→立刻还原→grep 证还原**。

**待 owner（本轮新增两项，已并入 R15 编号清单）**：**Q-16** 孤儿暂存文件的处置（票 73 的判据依赖它）、
**Q-17** 批准卡"无法规范化"却让人批 L2 写盘＝**单层防御 + 人看不见**（细节在 **A38②**）。

### 4-old（2026-09-20 23:5x 版，取代上面 23:19 版与更早的 08:1x 版）

⚠ **先读这条，它会推翻本文档其余部分的一个隐含前提**：**CI 从来没绿过**（A27）。
`gh run view 35517463335 --json jobs` 实测 `lint`/`test-core`/`test-windows`/`slo-smoke`/`slo-full`
**5/5 failure**。所以 §3 票 08 那行"CI 五 job 矩阵"是**建好了但一次都没通过**，
凡是"CI 会拦住 X"的论证都缺前提。修它的是**票 70**（见下）。

- **done 计数：18 张**（判据不是这句话，是 `ls .scratch/wisp/issues/*-done.md | wc -l`，23:56 实测 = 18）。
  ⚠ **和 23:19 版的"18"是同一个数字、不同的成员**：那一版把**票 62 算进去了**（我今天撤销的假 done），
  又**没算票 66**（活全交了但文件一直没改名、且缺 1:1 裁决表）。两个都纠正后：**-62 +66 = 18**。
  今天新收口：**10**（agent 循环核心，我复验：整体退回旧实现→终止测试 10.00s 变红）、
  **11**（LLM REST 适配器，对抗验收 AC#6 PARTIAL）、
  **63**（`wisp secret` CLI，AC#6 转票 12 并已闭环）、
  **66**（SLO 判据仪器返工：**7 条 AC 全交**，5 次连续无 keeper 跑 `exit=0`，我用真 exit code 独立复验，
  泄漏自检仍能翻红 ⇒ 门没被调松；`scripts/slo-check.ps1` 的第三处仪器缺陷 `.Count` on 零命中已修 `cf0c050`。
  裁决表 `docs/evidence/s1/66-adversarial-acceptance.md`：**7/7 PASS、0 BLOCKER、0 MAJOR**，
  每行标了"编排者独立复现／代理日志＋归档我抽验／仅代理自述"三档之一，**AC#1 的变异检验落在第三档**
  ⇒ 票面残口④明写了"票 70 落地后把 `WalkSystemProcesses` 顺序退回旧实现、确认
  `TestParseSystemProcessesKeepsLastSnapshotEntry` 变红"这条必做复核）。

- ⚠ **票 62 已从 `-done` 改回 `review`**（提交 `fc33532`，registry **A30①**）：我在**八个 AC 框一个都没勾**、
  且 AC#8 要求的 `62-adversarial-acceptance.md` **根本不存在**的情况下给它加了 `-done` 后缀，
  把 owner 对**票 65** 的降级放行（「算是赝品…先勉强用吧」）当成了票 62 的验收。
  **新会话既不要重跑它的实现、也不要以为它验收过**——票面现在有一张八行的"每条 AC 现在归谁"表。
- **⚠ 23:45 机器休眠 8 小时，把当时的三个在途代理全部冻死**（08:25 醒：`tasklist` 无任何
  `go.exe`/wisp/balldebug 进程、无完成通知、无未跟踪产物）。醒后的处置：
  - **票 70 的代理死前已把 gofumpt 全仓 sweep 落到工作树（67 个 `.go` 未提交）** ⇒ 我**审计后入库为 `f342413`**：
    `gofmt -l internal cmd tools` 空、`go build ./...` rc=0、`go vet ./...` rc=0、
    `git diff -w` 只剩 composite literal 展开与参数尾逗号 ⇒ **零语义改动**。入库用的是
    `git diff --name-only -- '*.go'` 的**显式清单**（不是 `git add -A`）。
    ⇒ **AC#1 的"本地落盘"半已闭，"CI lint 那一步真绿"仍待实跑**，接续代理已按断点派出。
  - 票 64 的 winlive 代理**什么都没写出就死了** ⇒ 原样重派（判据不变：交互四项任一 SKIP 就退回 AC#4）。
  - 票 62 契约草案代理**第一次连接即失败**（3 次工具调用） ⇒ 已重派。
- **正在跑（三个写码代理，00:52 起）**：**票 70** CI 逐 job 转绿（已自跑并入库 `8bfd47d`，正在办 R16 的
  `cmd/balldebug` 三处与 `tools/d22scan` 覆盖面）、**票 69** C21 表 61 行双向机器检查 + 变异检验
  （`internal/ball/tokens*`，只读 `tokens.go` 的值）、**票 20 第 5 框** 桥层真 junction / 8.3 拒绝用例
  （`internal/tools`，明令"建不出真产物就大声失败并写明缺什么权限，不许用字符串假装"）。
  ✅ **票 62 的契约草案代理已交付**：`docs/evidence/s1/62-visual-spec-draft.md`（`48f0cfd`，193 行），
  我核过 **20 态逐行 / 22 条与 SPEC-08 的冲突 / 15 个待 owner 定的问题 / 一格都没编数字**。
  ⚠ 桌面此刻**空闲**（票 64 的 winlive 已完成并入库 `c26abca`）。

- **工作树**：**干净**（`f342413` 之后 dirty=0），origin 与 cnb 均已推平（`ahead=0/0`）。
  ⚠ 但**下次读到"树脏"不要再当成我的疏忽**：今天两次是**代理死在批量写之后、提交之前**
  （票 66 前两批分别在 102 次调用与更早，票 70 这次是 67 个文件）。判据：**`ls -t` 看最新 mtime + `tasklist` 看有没有 go 进程**，
  两个都静默才认定死亡；恢复动作是**先审计再入库**（格式化与逻辑改动分两个 commit），
  而不是 `git checkout .`——那会直接把别人的工作抹掉。

- **review 中（勿重开已交付的 AC）**：**票 67**（AC#1/2/4 已闭并被我独立复跑；**AC#3 压住**：
  要改 `cmd/wisp/providers.go` 的 `PASS`/`FAIL` 字形，与票 70 同树）、
  **票 68**（AC#1 三列表 + AC#4 已交付；**AC#2/AC#3 blocked-on-owner**，R15 #3/#4/#5 未答前
  **不得翻 `prototypeVisuals` 默认值**）。
- **blocked-on-owner**：**票 65**（缺图1/图2——到今天为止没有任何代理见过参考图）、票 68 的 AC#2/3。
- **ready 但被显式压住**：**票 69**（C21 表 61 行无机器检查；与票 68 同动 `internal/ball` ⇒ 串行，
  而此刻票 64 的 winlive 正占着那个包）、**票 20 残口**（未勾的三框按**位置**称：第 5 框 junction/8.3 桥层、
  第 6 框 artifacts spill、新补的**第 7 框** `allowed_dirs` 首次询问流；另有一条不在框里的 **A18** 真 `taskkill` 残留）
  —— ⚠ **格式化阻塞已于 `f342413` 解除**，现在压住它们的只剩"包内正有代理在做"这一条；
  **票 71**（每道门自报工作量：`blocked-on-70` 现在**只剩 `cmd/balldebug`**，因为 R16 的三处
  `Registry.Spawn` 改动归票 70 的接续代理）。
  ⚠ **别给票 20 的框编 AC 号**：它的框原本无编号，我今天就先写错成 "AC#8" 又 grep 回来。

- **registry 从 A13 一路涨到 A33 + 裁定到 R16**；`docs/reports/pending-and-issues.md` 是唯一真相源，
  **新会话先读它的「待 owner 拍板 R15」——现在 11 项**（第 10 项＝批票 62 草案的 Q-1…Q-15，
  第 11 项＝`SPEC-08:42` 引用了不存在的 `docs/evidence/s1/68-*` 该怎么修）。
  **不答第 10 项，票 65 与票 68 的视觉半全部动不了。**
  ⚠ 另有两条**账实核对**结论供信任用：**A32**（18/18 张 done 票都有 1:1 裁决表；仅 07/11/63 有未勾框且
  三条都有书面接手方 ⇒ 不是假 done）、**A31**（我自己用未加引号的 heredoc 写 commit message，
  把反引号当命令执行、在仓库根造出 16 个垃圾文件，并把别人预先 staged 的重命名一起提交进去——
  **固定动作已改成：`git add` 后、`commit` 前必须 `git diff --cached --name-only`**）。




## 5. 下一步（新会话的行动清单，按序）

1. **票 10（agent 循环核心，2–4h 大票）——新会话第一票**：被 05✅+09✅ 解锁。参考：Pi 骨架（SPEC-05 §2）、golden 回放（票 09 的 `internal/llm/golden` + mockllm 控制端点）、failToolCallsFromTruncatedMessage、D15 阈值随 context_window 缩放。
2. **票 12（S1 gate）**：被 04✅+06✅+07✅+08✅+09✅+10 解锁。
   - **票 19（taint 引擎）+ 20（host bridge+fs 工具）**已部分就绪：17/18 完成后即为自然后续，建议新会话按 10 → 19∥20 → 12 推进。含 SLO 实测（回落调优义务——T02 发现 Settle 残留 ~31MB>25MB，`debug.SetMemoryLimit` 调优是票 12/15 的明文义务）+ 打字→回复→通知端到端 + 人工视觉项（H1）。
3. **S2 收口**：票 15（TTS 引擎——**TTS 常驻+ASR 按需为默认策略**，spike 实测串行 4.7s>1.6s 预算）、票 16（S2 验收：CER 双门禁 6%/15% + 真机冒烟三条物理场景）。
4. **S3 批量**：17–25 九张票（安全+能力面），两并发按 17∥18 → 19/20 → 21–24 → 25 推进；每片完成跑缺口审计+对抗验收。
5. **Key 类等待项**：真 LLM Key 到位后先跑 `cmd/llmrecord` 补录真 provider 黄金流（票 09 的 H2 DEFER）。

## 6. 关键技术事实（省得新会话重踩）

- **spike 结论（实测，已回填 docs/SLO.md）**：Path **Y**（单进程 cgo，空闲 16.5MB≤25）；**TTS 常驻+ASR 按需**（串行切换 4.7s>1.6s）；WebView2 冷 880–1126ms（P11 过，L2 卡不回退原生）；goja **无 ES modules**（Tier2=ES2020-minus-modules，Interrupt 精度 OK）；句柄口径 <600（D2D 窗口栈实测 ~420）；**Settle 残留义务**（Sleeping ~31MB>25MB → SetMemoryLimit 调优归票 12/15）；matcha 无官方 int8（TTS fp32 174MB）。
- **intra_op_num_threads=1 强制**；音频帧 16k/mono/int16/512 样本；`WISP_ENV` 三环境分叉（prod/dev/test 目录与互斥名全隔离）；config schema_version=**2**（catalog v2）；SQLite schema v2=9 表（provider_health）。
- **golden 格式**：`# wisp golden sse v1` + `# @response/@latency` 指令——单测回放器与 mockllm 子进程两跑器字节一致（契约）。
- **WASAPI 勘误**（真机发现）：IAudioClient vtable Start=10/Stop=11/SetEventHandle=13/GetService=14；共享模式 BUFFER_SIZE_NOT_ALIGNED 对齐重试——`internal/audio/wasapi_windows.go` 有注释留痕。
- **待修技术债**：cmd/wisp 的 sherpa cgo 使 `CGO_ENABLED=0 go build ./...` 在 cmd 失败（internal 全绿——历棒基线，票 08 CI 已按此设计）。

## 7. 事故手册（已验证的应对）

| 症状 | 应对 |
|---|---|
| 子代理派发即死（Captcha timed out / exceed quota） | 直接重派续传（随机性，重试存活率高）；持续失败 → 降并发或编排者亲自 |
| origin push TLS/EOF | 重试≤5 次×12s；cnb 稳定可先推 |
| 子代理死在中途 | 票据 Progress log 是断点：新代理读日志 `next=` 字段接续；盘上未跟踪 WIP 属实且完好 |
| 平台容量（Start Plan busy） | 回落 2 并发；验收可由编排者亲自（独立性=实现者≠验收者） |
| 截止窗口（配额到期） | 提前发警戒给在途代理（08:30 硬停规则：干净单元边界 deadline-pause）；大票不派发；HANDOVER 提前写并随收口刷新 |

## 8. 待人项（docs/reports/pending-and-issues.md 同步维护）

- **H1** 悬浮球 20 态视觉签收（看 docs/evidence/s1/ball-states/*.png）
- **H2** LLM API Key（阻塞真 provider 黄金录制；框架已就绪，`cmd/llmrecord` 一条命令补录）
- **H3** P10 命名残余核查（阻塞票 56）
- **H4** web.search 实现路径拍板（票 22 以接口先行，不返工）
- **H5** 默认 TTS 选型 —— **registry 权威编号下 H5 = TTS 选型，已裁定 R5（云 TTS，供应商 StepFun），
  已移入 pending-and-issues.md「已解决」**；本项不再是待人项。（本地 matcha 条目保持 `blocked-p3` 不变。）
- **（旧 H5 误标）生产 minisign 密钥仪式** —— registry 未给它 H 号，登记为
  「[H5-旧编号/HANDOVER 的 minisign 项]」，按 R6 降级为发布前才需要；dev 密钥
  （`E:\work\base\wisp-minisign\`）继续用于本地更新/模型链验证。
- **（旧 H6 误标）票 16 真机三场景**（物理拔插麦克风/隐私开关/独占占用——人工操作）——
  registry **无 H6 条目**；该人工件按 [B3] 补登于 pending-and-issues.md 的「完成度审计」节。
  引用 H 号时以 registry 文本为准（本文件此前把 minisign 记作 H5、把票 16 记作 H6，与 registry 冲突，已按 registry 更正）。

## 9. 进度节奏参考（新会话据此排程）

实测：2 并发下 ~9 票/24h（含平台故障开销），单票中位 1–2h、大票 2–4h、验收 0.5–1h。
剩余量级：S1 尾（10+12）≈1 天 → S2 收口（15/16）≈半天 → S3 九票 ≈3–4 天 → S4 ≈3 天 → S5 ≈4 天 → S6 ≈2 天 → S7 ≈4 天 → S8 另计。**全项目 ≈3–4 周连续双代理运行。**
