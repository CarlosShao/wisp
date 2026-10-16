### 4.0du 停车点（**2026-10-16 11:3x +08**）｜编排者｜第五波收完（三枚验收腿）⇒ 翻框 14 枚＋票 317／325／287 改 `-done`＋台账 `A1004`＋第二枚写腿 `326-r1` 已放∥第七波三枚派单**已存档、⛔ 未放腿（CPU 76%＞70% 那把尺）**

**① 盘上状态（接手先复跑这三把尺，⛔ 信本节）**：`git log -1`＝`ce8df8b6`（A1004 收档笔）之后另有派单存档笔；本地领先 `origin/dev`＝`git rev-list --count origin/dev..HEAD` 现量 **70 枚**（⛔ 裸推会去 cnb、⛔ force；推送权只在我，机主授权 `A813`）；`docs/reports/pending-and-issues.md`＝**18722 行**（`A1001`／`A1002`／`A1003`／`A1004` 逐枚 `grep -c '^## A10'` 在盘）；本文件 §4 节头 **145 枚**、最新一节＝本节 `4.0du`（`4.0dt`／`4.0ds` 逐名仍在）。

**② 今日累计（10-16 09:1x→11:3x）**：收档 9 枚腿（`317-v1`／`322-v1`／`216-a1`／`166-v1`／`287-v1b`／`325-v1`／`165-v1` ＋ 两枚死腿遗产代提）；**翻框 25 枚**（上波 11＋本波 14）；**四枚票改 `-done`**（本波 317／325／287，⛔ 票 165 不结案＝残余具名三笔）；**新立一枚票 326**（`AppendSessionTurn` 只插子表、不刷新 `session.updated_at` ⇒ 正在聊的会话永远不冒顶）并放出落地腿 `326-r1`。

**③ 在飞（两枚写腿，⛔ 第三枚）**：`220-r2`＝本波唯一碰 `cmd/wisp/**`＋`internal/panel/**` 的写腿（现读 `M cmd/wisp/panel_pump.go`）；`326-r1`＝`internal/memory`＋`internal/session`（agentId `ageneral-purpose-f4a89d405fb3b3da`，启动补记在 `f912a5af`）。★它俩交件后我要办的两笔：票 322 的 **AC#4** 缺那枚 `cmd/wisp` 整包读数（⛔ 票 322 在此之前不许 `-done`）＋票 220 供数那一跳的非实现者验收 `220-v1`。

**④ 第七波：派单已存档并 commit，⛔ 还没放腿**（规矩＝写了派单⛔ 等于放了腿）。件在 `.scratch/wisp/probes/orch/20261016/txt/`：`270v1-dispatch.txt`／`120v2-dispatch.txt`／`165c2-dispatch.txt`／`batch1005-send.txt`。释放条件＝同一把 CPU 尺（`Get-CimInstance Win32_Processor` 的 `LoadPercentage`）复测 **≤70%**，或任一枚写腿收档；放腿后必须回头在 `batch1005-send.txt` 补 agentId。

**⑤ 队列（每枚带释放条件，⛔ 无条件待派）**：`216-r1`（写面 `internal/panel/**`）与 `220-v1` 等 `220-r2` 交件；`165-c1`（把 `plan_present_165_test.go:14-15` 那两行"−overlay 注入"幻影注释改说实话＋按 `165-c2` 裁形装"门铃"仪器）等 `internal/tools` 面空；`313` 的 `cmd/wisp` 落点包、票 247 的真电平接线（`247-r1`）都撞 `cmd/wisp` ⇒ 同批按住；`259-v1`／`181-v3`／`253-v1` 三枚陈旧前提待我逐枚复验后才谈派；`270-a1`→已结案、`135-r1` 待我重读票面。

**⑥ 机械手册（本波新增／复用的坑，接手直接用）**：ⓐ 翻框走机械件（`txt/flip-1004.sh`：只改 `- [ ]`→`- [x]`、逐枚指名 AC 标题），跑前 `git diff --numstat` 证"改的行只有框"（本波＝14 增 14 删）；ⓑ 改名 `-done` 前必量**枚数＋长度门**（三枚新名 107／116／114 字符仍在帽下，`check-path-length-budget.sh --with-self-test` 复跑 GREEN＝`over-budget=60 covered=60 not-in-roster=0`）；ⓒ 越界核查**只许逐枚 `git show --name-only`**，⛔ 用 `git log --name-only A..B`（区间会把别腿的提交扫进来＝上波那枚"139 枚越界"假指控的成因）；ⓓ `GOTMPDIR` 原生反斜杠＋**先 `mkdir -p`**（没建目录＝`creating work dir … cannot find the file specified` 长得像"红"）；ⓔ ★我自己今天抓到的新形：`grep -rnF '${LOCALAPPDATA//\\//}'` 报 **0＝假 0**（同一事实换 `grep -rn 'wisp_311_ref_path='` 得 **4**）⇒ 报 0 之前先跑正控、正控也 0 就换一把具名尺，⛔ 删坏尺；ⓕ 腿自报的 HEAD 串可能**长度都不对**（本波抓到 `287-v1b` 一枚 36 位串＋一枚不存在的短号）⇒ 一律现跑 `git log -1` 核。

**你需要做什么：没有。** 两枚写腿在跑，第七波等 CPU 落回 70% 以下就放。
