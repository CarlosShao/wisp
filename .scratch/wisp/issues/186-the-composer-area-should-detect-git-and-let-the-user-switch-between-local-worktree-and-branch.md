# 186 — 聊天框附近要能**认出这棵工作区是不是 git 仓库，并在"本地／工作树／分支"之间切**（owner 09-28 明确要的功能；`Q-64` 已改判＝**做**）

- Status: **ready-for-agent（先只读普查，再落地）**。⚠ 本票是**功能票**，不是缺陷票。
- 来路：owner 09-28 14:3x 的原话（逐字见台账 `A368`）＋他给的两张截图（聊天框那一排：加附件／权限下拉／模型＋思考档位／发送，下方那一排：工作区名／"本地"／分支名／占用条）。他明确否掉了我先前的读法——**"1234 都是跟聊天框有关的，不是真让你切换到别的 git 干"** ⇒ 这不是"面板去动编队正在跑的那棵树"，是**产品要有的工作区/版本上下文切换能力**。
- 关联：**票 181**（只读那一半：Go 侧今天零 git 读面，先把"读什么、从哪儿读"定下来）· **票 145 AC#2b**（模型清单与思考档位的载体）· **票 187**（模型／档位的写入口，与本票同族不同物）· **票 92**（面板"只许显示＋发起请求"的口径，本票把它**扩**成"可发起切换请求"）· **票 102／C26**（选工作区＝授权一棵树，那一层的改写与拒绝必须复用）

## 为什么值得做（不做会怎样）

用户选完工作区目录之后，Wisp 今天**完全看不见**那是不是一棵 git 树（尺：`grep -rln 'exec.Command("git")\|\.git/HEAD\|rev-parse' --include=*.go internal/ cmd/` 排 `_test.go` ⇒ **零枚**，台账 `A360` 已量）。⇒ 没有这一维，"在哪个版本上干活"这件事只能靠用户自己在终端里记，聊天框那一排就永远只是一个输入框，成不了 harness 的驾驶位。

⚠ **一枚已经躺在仓里的半成品（别重造）**：`internal/panel/workspace.go:76 RequestWorkspaceSwitch(scope, input, audit)` **已经实现**——它做"解析输入 → 走 C26 → 检查改写/重解析 → 窄化工作区根 → 审计一行"，而且**有跟踪判据**（`internal/panel/workspace_test.go` 六枚：干净切换／重解析拒绝／被改写拒绝／失败保留原范围／范围自己动了／快照渲染）。尺：`grep -n "func RequestWorkspaceSwitch" internal/panel/workspace.go`、`grep -n "^func Test" internal/panel/workspace_test.go`。
⇒ **它今天零生产调用方**（尺：`grep -rn "RequestWorkspaceSwitch" --include=*.go internal/ cmd/ | grep -v _test.go` ⇒ 只有定义与注释，无调用点；`cmd/wisp/panel_pump.go:79-81` 那句注释自己写着"the tree still has no caller for"）。⇒ **本票"切本地目录"那一支主要是接线，不是从零写**。

## 这一票要交付哪三样（分开可验）

1. **认出 git**：工作区根（或它的任一祖先目录）里有没有 `.git`；有就报"当前分支／是否 detached／仓库根在哪"，没有就**显式报"不是 git 仓库"**（不许留空——票 92 的"宁缺毋造"同形，留空＝用户以为功能坏了）。
2. **列出可去的地方**：本地工作树清单（`git worktree list` 那一族的等价读法）＋分支清单。
3. **切过去**：在聊天框附近选一枚（工作树／本地目录／分支）⇒ 宿主把这次干活的范围换过去，**并且把"换成了什么"如实显示回来**。

## 现量（锚 `72c76d42`，09-28 14:3x 编排者本程现跑）

| # | 断言 | 尺 | 读数 |
|---|---|---|---|
| 1 | Go 侧零 git 读面 | `grep -rln 'exec\.Command("git")\|\.git/HEAD\|rev-parse' --include=*.go internal/ cmd/ \| grep -v _test.go` | **零枚**（与 `A360` 一致） |
| 2 | 工作区切换处理器已存在 | `grep -n "func RequestWorkspaceSwitch" internal/panel/workspace.go` | `:76` |
| 3 | 它零生产调用方 | `grep -rn "RequestWorkspaceSwitch" --include=*.go internal/ cmd/ \| grep -v _test.go` | 定义＋两处注释，**无调用点** |
| 4 | 它已有跟踪判据 | `grep -c "^func Test" internal/panel/workspace_test.go` | **6** |
| 5 | 面板可调方法今天恰好四枚 | `grep -n '"panel\.' internal/panel/bridge.go` | `:42-45`＝mode.request／workspace.request／attachment.add／message.send |
| 6 | `panel.workspace.request` 那枚名字**已在册但没人接** | 同上 `:43` ＋第 3 行 | 名字在白名单里，**写腿没接** ⇒ 本票"切本地"那一支的落点很可能就是这一枚，不必新起名字 |
| 7 | D34 内置工具权威表里没有任何 git 行 | `awk 'NR>=2530 && NR<=2585' docs/PLAN.md \| grep -nE "git\|分支\|branch"` | **0 命中**（`A360`）⇒ **本票不做成"模型可调的 git 工具"**（那是动 D34），做成**宿主侧读面＋面板请求**（票 181 的 AC#3 同一条边界） |

## AC（每格都要答"这一发在**未修码**上响不响"；先测→再写→再提交）

- [ ] **AC#1 只读普查先答"三样各自怎么读、要不要起新进程"**：① 认 git——读 `.git/HEAD` 这一条文件就够，还是要 `git` 二进制？② 工作树清单——`.git/worktrees/` 目录法能不能拿到"路径＋分支"两样，拿不到要不要跑 `git worktree list --porcelain`；③ 分支清单——`refs/heads/` ＋ `packed-refs` 两堆合不合得上。**每支给一条现跑的命令与读数**（在真仓上跑，本仓就是一棵真 git 树：`git rev-parse --abbrev-ref HEAD`、`ls .git/worktrees 2>/dev/null`、`cat .git/HEAD`）。⚠ **跑外部 `git` 二进制是一枚新的宿主能力**（今天全仓零次），要么给出"不跑也能读全"的证据，要么把"要跑"具名报回并说明它归不归 `allowlist.txt`／`d22scan` 的哪一条管。
- [ ] **AC#2 定"切换请求走哪一枚方法"，并把它写成契约变更入档**：候选＝**复用已在册的 `panel.workspace.request`**（它今天零写腿，正好是"有名无实现"那一形，见票 35 的 `panel.resync` 前例）／新起 `panel.worktree.switch`／新起 `panel.git.switch`。**判据＝哪一枚改动面最小且不新增白名单条目**；若结论是"必须新增条目"，⚠ **C17 白名单定稿是契约面**——owner 09-28 这句"别老说契约冻结"是**方向授权、不逐条点名**，所以**具名报回要加哪一枚、我落 `A##` 之后再写**。
- [ ] **AC#3 安全约束做成实现约束，不许当否决理由**（⚠ 这条是我上一轮的错处，写死在这里）：切走工作区/工作树会换掉**宿主正在用的那棵树**，09-27 真出过一次（7 枚提交落到错误分支，`HANDOVER §4.0u`／`A328`）。⇒ 落地必须回答三问并各带判据：**(i) 有程在飞时这枚请求怎么摆**（拒绝并说明／排队／还是照切——**选一支并给判据**）；**(ii) 切完界面显示的那棵树与运行时真用的那棵树必须同源**（不许显示 A 跑 B——`RequestWorkspaceSwitch` 已经做了"失败保留原范围＋如实渲染"，接线不许绕过它自己写一份）；**(iii) 切换要写审计行**（现成的 `AuditFunc` 腿，不许吞）。
- [ ] **AC#4 非 git 目录那一形必须显式回答**：工作区不是一棵 git 树时（今天 `wisp run` 允许任意授权根），快照里这一维要带**一句人话理由**（票 92／145 的"宁缺毋造"），不许画一枚点了没反应的按钮，也不许留空让用户以为坏了。
- [ ] **AC#5 反向判据（防这一族走偏成"模型能自己切仓库"）**：交付里**不许出现**任何模型可调用的 `git.*` 工具、不许动 `docs/PLAN.md` 的 D34 那张表、不许把切换能力塞进 `internal/tools/` 的执行面。⇒ 补一枚常驻判据钉住"这一维只有宿主侧读面与面板请求两条路"（尺的形状照票 181 AC#3）。
- [ ] **AC#6 覆盖面：切完之后"哪些东西跟着换"要说清**——授权根（C26／`allowed_dirs`）、产物目录、会话历史、污染名册（C25 的 per-scope mark）各自在切换后的行为，逐枚答"保留／重建／作废"，各带一条现跑读数或一个具名落点。⚠ 这一格是本票最容易漏的：**换树不换名册＝上一棵树的证据会咬下一棵树的内容**（票 183/185 那一族在同一根管子上已经咬过三次）。
- [ ] **AC#7 契约轴**：`docs/PLAN.md`（D34／D36／D38／C26 都近邻）、`docs/specs/**`、`allowlist.txt`、`thresholds.go`／golden／审批超时常量一字节不许动；`internal/panel/workspace.go` 的既有语义只许**追加**不许改写（它带着票 92/102 的裁决）。⚠ C17 白名单那一处按 AC#2 的规矩**先入档再改**。
- [ ] **AC#8 门禁**：逐包 `go test -count=1 ./internal/panel/ ./internal/config/ ./internal/risk/`（⚠ `./internal/risk/` **单跑**，并发会假红＝`A359`）；CLI 那一面走 `-overlay` **且带在册 PATH 前缀**（`PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"`，出处 `scripts/wisp-cli-tests.sh:20`；不带就是 `0xc0000135`＝票 98）；`sh scripts/d22scan.sh`（基线 `ban #8 internal/` **examined=433**，多一枚少一枚都要具名解释）；`bash .scratch/wisp/probes/154/gate-clauses.sh` **比红腿名册不比退码**（今天唯一在册红腿＝`G6neg`＝票 178）；⚠ 终态读数取在**你自己最后一枚 commit 之后**。⚠ **不许跑 `probes/161/r6/flip-declaration.sh`**；⚠ **重跑任何既有台件之前先看它往哪枚路径写数**（就地覆盖型 sink 会把别家票逐行引用的读数洗掉，09-28 实测＝`A367`）。
- [ ] **AC#9 界面那一半不在本票写面**：`frontend/**`／`design/**` 归另一枚会话（owner 09-25 决定）。⇒ 本票只交 **Go 侧的读面＋请求腿＋快照字段**，并把"面板要画哪几枚控件、点下去发哪条请求"写成一张**递给前端的需求表**（`docs/evidence/s1/186-panel-ask.md`）。⚠ 编排者会把这张表转给前端会话，**本票不代写前端**。

## 本票**不**解决

- 不做"远程操作"（fetch／push／PR／克隆仓库）——那是另一族，且 owner 这次点的是**本地／工作树／分支**这一排。
- 不选模型、不改思考档位（＝票 187）。
- 不动输入框宽度（＝票 180）、不裁任务监控那一栏缺哪些堆（＝票 182）。
- 不碰"面板侧来源的 L2 允许"那一族（票 49／`Q-49` 已钉过五枚门，本票不重开）。

## Progress log

- 09-28 14:3x 编排者立票：owner 改判 `Q-64`＝**做**（原话入 `A368`），我上一轮把它读成"让面板去切宿主正在跑的树"并推荐默认不做——**那句读法作废**；风险改成 AC#3 的实现约束。上面 7 把尺本程现跑（锚 `72c76d42`）。⚠ 本票与票 181 的普查**并成一程派**（同一枚接缝、同一枚前端需求表）。未派。
