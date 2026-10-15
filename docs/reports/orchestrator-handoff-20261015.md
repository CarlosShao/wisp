# 编排者交接文档（2026-10-15 19:2x +08）——给接手本编队的下一位 agent

> **这份是什么**：本会话（ZCode，编排者）把 Wisp 流水线的活交出去。**新接手者只需要读这一份＋它指的三处权威文件**，不需要本会话的任何历史。
> **生成时刻**：2026-10-15 19:2x +08（`date` 现读）。**交接锚**：`HEAD=origin/dev=626fcc7b`（未推 0）。

---

## 0. 三十秒版（先看这个）

- **你是谁**：Wisp 项目的**编排者**——不写生产码（除机主点名/死腿代提），负责**派腿、收件、复跑、翻格、记账、向机主汇报**。仓库 `D:\work\workspace\projects plans\Wisp`（Windows＋Git Bash，路径带空格要引号），分支 `dev`。
- **现在的状态**：**在飞 0 腿**；工作树只有"别家 agent"的既有噪声（见 §5）；机上跑着机主的 dev 实例 **PID 21892**（今晚新编的 `build/wisp.exe`，⛔ 别 kill、别关它的面板窗）。
- **机主刚做完什么**：在模型页把 StepFun 配好了。**其计费口径（务必对齐）**：**按量（API）行只承担 realtime**，其余（ASR/TTS/对话）**一律走套餐行的模型**（台账 `A988`）。
- **下一步队列（具名，按此顺序）**：`321-r1` → `165` 的 `cmd/wisp` 装配半格 → `120-r1` 修复轮 → `270`/`287` → `171`/`86`；另两条**待机主拍板**（§7）。
- **进门前先跑一遍对账**（§4），别信这份文档里的数字——它们会腐烂。

---

## 1. 权威状态源（引用前**先现跑**，本文件只是指路）

1. `docs/reports/HANDOVER.md` **§4 最新那一节**（现最新 `4.0dq`＋`4.0dp` 追记；最新是本会话 19:0x 那条）。它是"从哪儿接"。
2. `docs/reports/pending-and-issues.md` **台账**（`A##` 账与裁定／`R##` 退回／`Q##` 待人拍板；**只追加**）。现最新 `A988`。**本会话的关键条目**：`A972`（票 320 五裁）、`A975`（111 收口）、`A976`（319-v 判语＋D1/D2）、`A978`（抬前台裁定）、`A980`（319 收口）、`A982`~`A988`（机主编新 exe／前端三轮实报修复／删除能力／机主计费口径）。
3. `AGENTS.md`（仓根薄索引：禁止清单／未定义即停／每片后五件事／D1–D47 一行摘要）＋ `docs/specs/SPEC-12-roadmap-governance.md`（治理硬规矩）。
4. `docs/reports/orchestrator-handoff-20261013.md`（**前任编排者 403 行全套交接**：角色/铁律/机主偏好/队列/派单模板/陷阱）——本文件不重复它的全部，只补**今天新增/变化**的部分。

⚠ 老 `docs/reports/next-session-bootstrap.md` 是 10-15 上午那版（写的是 `319-v`／`320-a1` 那一批）——**已过期**，别照它做。

---

## 2. 铁律速记（**违反＝判失败**，全文见上一节第 3/4 条）

**推送**：**只有编排者能推**，**只推 `origin/dev`**（机主授权 `A813`）。⛔ 裸 `git push`（会去 cnb 远端）、⛔ force。**子代理"腿"只 commit、不 push**。
**Git 纪律**：`add`／`commit` 两处**显式 pathspec**；⛔ `add -A`／`add .`／目录；⛔ `--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`；改名用 `git mv` 后**记得补一笔"旧名删除"**（本仓 `git mv` 分两笔落，先例：133/111/319 三次都这样）。
**面互斥（本仓最常被撞的一条）**：`cmd/wisp/**` 与 `internal/panel/**`＝**共享编译面，每波只许一枚写腿**；`go test ./cmd/wisp` 重包含"整进程死亡"死区（票 314）⇒ 同一时刻只跑一个。其余包（`internal/tools`／`internal/winsec`／`scripts`／`internal/config` 等）各自独立面，可并行。
**派单规格**（机主 10-14 明令、10-15 又强调"多派"）：**写腿 ≥3／总 ≤7／写 ≤4**；CPU ≤70% 就补位；与面互斥冲突时**以面/CPU 为准并具名报账**。
**派单先存档再发**：`.scratch/wisp/probes/orch/YYYYMMDD/txt/<腿名>-dispatch.txt`＋发送时 `<腿名>-send.txt`；**派前查重** `ls .scratch/wisp/probes/<票>/`。
**收件三把尺**：① `git log -1 <sha>` 提交真在 ② `ls`/`wc -c` 件真在盘 ③ **自己复跑**关键读数。**完成通知一律当未验证**；空回执＝幽灵，按"轻量催收"（让它自己补 commit/件/票尾 bullet）。
**每收一条腿立即 `TaskStop`**（防面板残影）。
**报告结尾必带「你需要做什么」栏**（没有就写"没有"）——机主要求，缺了算格式错。
**别碰**：`frontend/**`、`design/**`（归机主的**外部前端 agent**；本会话的例外见 §6）、别家腿半成品（现盘：`.scratch/wisp/probes/152/my152.py`、`268/v1/evidence.md`、`design/**` 的删除态）、`docs/PLAN.md`／`docs/specs/**`／阈值/golden／`tools/d22scan/allowlist.txt`／`.github/**`（改 `.github` 先登记）。

---

## 3. 本会话（10-15）产出了什么

**收口三票**：**133**（`-done`，AC#6 翻勾，CI `38582016056` 的 PASS 读数）· **111**（`-done`；死腿 `111-r7` 活完整，我**另取读数**亲跑五格复核后翻；顺带校正 `ci.yml` 三处过期数字＝纯注释）· **319**（`-done`，全链：`319-r2` 补牙 → `319-v` 抓出 D1"已可见面板被藏窗" → `319-r3` 修 → `319-r4` 补"抬前台" → `319-v2` 判四形全成立）＝**机主"界面打不开"这件事闭环**。
**票 320**：三写入（计费模式／自定义行／模型条目）＋**取件端口 `config.models.list`**（我亲做）；六格已翻，**只剩 AC#4 的"非实现者复核"**（代码与真腿测试齐、机主真机验收过）。
**新立**：**321**（托盘/热键执行器的接线无钉——`319-v2` 交回的具名缺口）。
**机主驱动的前端改造**（★这一节接手者最需要知道）：机主今晚连报四次，我按实报就地改**前端**（§6 有边界说明）：
  1. **拉取端口**：`panel.ts` 的 `requestProviderModels`＋标记行 `wisp-models-list:`（Go 半＝`panel_models_320.go`）。
  2. **新增供应商表单重做**：`自定义接入` 页签有 行名/协议/基址/**计费模式** 四字段；**"常用第二行"**一键入口（stepfun ⇒ `stepfun-api` 按量口，预填 `https://api.stepfun.com/v1`）。
  3. **行身份读回**：面板载荷加 `baseURL/protocol/billing`（**有意移动** `config_313_test.go` 的键钉）⇒ 行上直接显示真实基址+计费标签+协议；计费按钮高亮当前值。
  4. **加完就拉**：新行建成后自动拉该行清单。
  5. **★删除能力**（此前宿主**根本没有**）：`config` 加 `RemoveModelRow`／`RemoveProviderRow`（**被 chat 角色或语音链引用则按名拒删**）＋面板两字段 `model_row_delete`／`provider_row_delete`＋两枚真腿测试；前端模型 chip 有 `×` 单删、整行删除两步确认。
  6. 修了一枚**真崩溃**：`onChange` 里在 `setState` 更新函数内读 `e.currentTarget.checked`（React 推迟求值时该引用已置空）⇒ 已改为派发期间取值。

---

## 4. 进门前对账清单（现跑，别抄本文档的数字）

```bash
cd "D:/work/workspace/projects plans/Wisp"
date "+%Y-%m-%d %H:%M %z"          # 所有你要写的时刻都用它插值（⛔ 别手打）
git log --oneline -8                # 最新提交（应含 626fcc7b A988 或更新）
git rev-parse --short HEAD; git rev-parse --short refs/remotes/origin/dev
git log --oneline refs/remotes/origin/dev..HEAD | wc -l   # 未推笔数（应为 0）
git status --porcelain | grep -vE "^\?\?" | head -20       # 谁在工作树里动（正常噪声见 §5）
tasklist | grep -i wisp.exe         # 机主实例在不在（PID 会变）
tail -30 docs/reports/pending-and-issues.md                # 台账尾
grep -n "^### 4\.0d" docs/reports/HANDOVER.md | tail -3     # 最新停车点
```
**判据**：如果 `git log` 顶部不是 `626fcc7b` 或更新 ⇒ 说明**有人在本会话之后动过**，先把那些提交读一遍再动手。

---

## 5. 工作树噪声（**是正常的，别"修"**）

- `design/**` 一片 `D`（删除态）＋`design/doubao/**` 的 `M`：**机主外部前端 agent 的面**。
- `.gitignore` 的 `.worktrees/` 段、`go.mod`、`frontend/src/styles/tokens.generated.css`：别家 agent。
- `.scratch/wisp/probes/152/my152.py`、`268/v1/evidence.md`：别家腿半成品。
- `?? -`、`?? .scratch/.scratch/` 等：历史遗留杂项，**别管**。
⇒ **一律不碰**；要提交只提交你自己的路径。

---

## 6. 机主与面边界（本会话新增的实况）

- 机主对**页面**的期望不是"能写进去"，而是"**一眼看懂、能改能删、像主流 harness**"。他明确说过：**「票不是最完美的，需求一直在变，别死板执行」** ⇒ 遇到"票写到的层级不够用"，先按他的意图把**用户能用的那一层**补齐，再把口径写进票面/台账。
- **`frontend/**` 的边界**：仍归机主的外部前端 agent（本编队通常零写面）。**但今晚这几次是机主明确要求我改**（原话见 `A986`）⇒ 若你也要动前端，**先问机主一句**是他改还是你改，别默认。留给前端 agent 的两条约定：① "常用第二行"入口不依赖下拉选择；② 事件对象的值在**派发期间**取，别挪进 `setState` 更新函数。
- **机主的计费口径**（`A988`，务必对齐）：StepFun **按量行只挂 realtime**，其余全走**套餐行**。
- 机主会**盯着截图骂**——这是他的工作方式；他骂到的每一条都当"实报缺陷"处理，别辩解，改完给他 30 秒能验的路径。

---

## 7. 队列与待人项（接手后的第一件事）

**可派（按序，面已空）**：
1. **`321-r1`**（票 321：托盘/热键执行器接线的常驻判据＋变异；面＝`cmd/wisp` 319 面）——`319-v2` 交回的具名缺口。
2. **`165` 的装配半格**：`internal/tools/plan_present.go` 已实现并有五格测试（腿 `165-r1` 已交），**但 `cmd/wisp` 的装配点没接线**（`BuiltinPlanEntries()` 一行即可，需另一腿或编排者排面）；票 165 还需**非实现者对抗验收**。
3. **`120-r1` 修复轮**（票 120）：仪器已以"基线形态"入库（含它自报的两处缺陷），AC#2 判据已答＝**"逐级 openat+O_NOFOLLOW 完全体在写面内做不到"**，替代形态＝fd 锚定＋`SameFile` 闸门；`next=` 五条写在 `probes/120/r1/index.md`。
4. **`270`**（`task.spawn` nil 防护）／**`287`**（深度判定夹具）——与 165 同抢 `internal/tools` 文件族，**一次一枚**。
5. **`171`**（三枚尺洞；前置＝"internal/tools 无写腿在飞"）／**`86`**（1ms 墙钟预算假红；条件＝无 docker/并发 go test）。
6. **票 320 的 AC#4 非实现者复核**（Go 半已落、真腿测试齐、机主真机验收过；只差一张**非实现者**的判语表，然后翻格；票 320 就只剩这个）。

**待机主拍板（`120-r1` 交回，我不替机主决定）**：① `privateFile` 的开（`winsec.go`，竞争窗口内仍可能跟随叶子链接 `O_TRUNC` 外来文件）要不要并进票 120 下一轮；② POSIX 密封改 fd 锚定后，**r 权限类（`0w0`/`000` 文件、缺 `r` 目录）由"可修"变"拒绝"**是否接受。

**冻结中**：机主 10-15 17:0x 令"额度将尽 ⇒ 先别继续派发，不要再搞半成品"。**新 agent 接手前先问机主一句"额度恢复了吗、可以派腿了吗"**。

---

## 8. 机械操作手册（今晚现学现卖，很省钱）

**编新版给机主用**（他喜欢"你说完我立刻能试"）：
1. `taskkill /F /PID <当前 wisp.exe>`（先掐，否则 exe 被占）。
2. `mv frontend/dist /d/tmp/<留底目录>` —— **必需**：产物名守卫要求"本轮 npm 必须产出新名字"，而若 dist 已是当前源码编的，重量产同名 ⇒ 守卫 FATAL。移动（留底、别删）后真重建即可过。
3. `powershell -NoProfile -ExecutionPolicy Bypass -File scripts/build.ps1 -Env dev`（跑 npm build + go build + `wisp.exe doctor`；末尾必须看到 `wisp doctor: PASS`）。
4. **拉起＝裸跑 `wisp.exe`（无参数）**——`wisp run` 是"带任务文本"的 CLI 腿，裸跑会立即退出。
   `powershell -NoProfile -Command "Start-Process -FilePath '…\build\wisp.exe' -WorkingDirectory '…\build'"`
5. 验内容：`grep -a -c "<你新加的字符串>" build/wisp.exe`（页面字符串与 Go 字符串都能这样查）。

**推送**：`git push origin dev`，**用重试环**（`github.com:443` 时通时断：`timeout 90 git push …` ×≤6、间隔 20s；本会话一批曾 6 试才成、另一批 1 试就成）。推前跑 `sh scripts/check-path-length-budget.sh --with-self-test`。`gh` 可用（CI 读数用 `gh run view/list`、`gh api`）。

**门禁清单**：`go build ./...`／`gofumpt -l <你的件>`（宿主现成 `D:\work\base\gopath\bin\gofumpt` v0.12.0）／`go vet`（winlive 面加 `-tags winlive`）／`sh scripts/d22scan.sh`（**注释豁免、字符串不豁免**，零 emoji）／`sh scripts/check-path-length-budget.sh`。
**环境**：`PATH="/d/work/workspace/projects plans/Wisp/third_party/sherpa-onnx:$PATH"`（缺它 `0xc0000135`）；`GOTMPDIR='C:\Users\swq\AppData\Local\wisp-test-gotmp\<腿名>gotmp'`（**C: 卷＋只反斜杠形；⛔ 别用 `C:\tmp\`**——会让 `TestAC2SealNoticeLandsInTheRunLegLogFile` 假红）。

---

## 9. CI 现况与两个"别被骗"的口径

- **主 CI 恒红 4 处（存量，非新债）**：`lint`(gofumpt×2＋staticcheck)／`test-core`(core scope)／`test-windows`(`cmd/wisp` CLI 2 枚真窗族＋`internal/panel` 那族)／`portable windows tests`。
- ⚠ **`cmd/wisp` 那一步会"早停"**（整进程死在开窗段，签名 `MARK299 …wndproc`＝票 314）。**读它的四数前先证该趟跑到终点锚 `TestZZ315A1*`**；没跑到就只记"该趟早停、四数不作分母"（口径出处 `A976`）。
- ⚠ 有一枚**在册间歇**：`TestTicket223HandEditedFsLooseningCostsAnL2Card`（同码两跑一绿一红，台账四处先例）⇒ 见它红先查是否 lone-red。
- 本地跑 `internal/panel` 全包会看到 **5 枚既有红**（`TestApprovalCardView…`／`TestStreamLog*35r8`×2／`TestComposerContractTypes…`／`TestC21DesignTokens…`）＝外部面族，**你的改动只要不增加它们就是干净的**（用"与基线逐枚同名"这个判法）。

---

## 10. 我这边的两处自纠（接手者别再犯）

1. **台账时刻手打**：`A975`–`A987` 的时刻字段我手打了，先后偏早约 2h／1.5h（`A982`、`A988` 已具名自纠）。**规矩：一律 `date` 插值**。
2. **把入口藏进不可达状态**：我给页面加的"切到自定义并预填"挂在"选中某厂商之后"，而那个下拉把已添加的预设过滤掉了 ⇒ 机主看不到入口（`A986` 已认）。**教训：给机主的路径必须"进页面就看得见"**，不许依赖某个可能为空的选择。
3. **★事故（本会话末，`A989` 有全套取证）**：写这份交接文档时，我给 `python -c "…"` 的**双引号**串里带了**反引号** ⇒ bash 命令替换**把 .md 当脚本执行** ⇒ 仓根被造 20 个中文名垃圾文件、**2113 个文件被抓进暂存区**（含前端 agent 的未跟踪件）。**零文件被删、零内容被改**；已用 bare `git reset`（只退索引）+ 删掉 20 个垃圾件恢复并验证。**教训**：凡进 shell 的正文**一律单引号或加引号的 heredoc**，⛔ 双引号串里不许有反引号/`$`（本仓此族第 5 次，前四次都是 heredoc 未加引号）。⇒ 接手者若发现暂存区忽然多出成百上千条，**先 `git diff --cached --name-status | grep -c "^D"` 看有没有删除**，再按本节的办法退索引（bare `reset` 不碰文件，安全；⛔ `--hard` 绝对不许）。

---

## 11. 你需要做什么（接手清单）

1. 读 §1 的三处权威文件（＋前任的 `orchestrator-handoff-20261013.md` 若有时间）。
2. 跑 §4 对账；`git log` 顶部若比 `626fcc7b` 新，先读那些提交。
3. **问机主一句**：额度是否恢复、可以派腿了吗（他 10-15 17:0x 令了冻结）。
4. 按 §7 的顺序派（记得：写腿 ≥3／总 ≤7／面互斥／先存档再发／每条汇报结尾带「你需要做什么」）。
5. 别碰 §5 的噪声、别关 §0 的机主实例面板窗。
