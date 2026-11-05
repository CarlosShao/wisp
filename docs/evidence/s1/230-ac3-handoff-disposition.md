# 票 230 `AC#3` 处置记录：票 80 → 票 21 那句交接＝**乙（作废＋两格接收在本票名下）**

- 腿：`230-ac3`（实现腿·只改票面·**零码**）。落刀时刻 2026-11-05 10:3x +08。
- 起手锚 `HEAD`＝`ba0ea5121c2e23b4a48b0636205735e618a0ddde`（`git rev-parse HEAD` 现读）。
- 本件是**读数与处置记录**，⛔ 裁决文书：`AC#3` 的框仍归编排者收件后翻（本票 `AC#1`／`AC#2` 的勾都是编排者翻的；`230-r1` 当年 8 笔一枚框没翻）。
- 尺一律**现跑**（⛔ 引票面旧读数；每数带命令）。落档＝本件＋票面 `AC#3` 名下一段＋票尾一条 Progress log。

## 0. 射程与硬约束（写面）

| 项 | 本腿 | 读数/说明 |
|---|---|---|
| 可写 | 票 230 票面（`.scratch/wisp/issues/230-four-cells-left-unfinished-inside-closed-tickets.md`）＋本件 | 派单原文：「only edit the 230 ticket file」 |
| ⛔ 写 | `.go`（含 `_test.go`）／票 80 票面／票 21 票面／`docs/PLAN.md`／`docs/specs/**`／`thresholds.go`／golden／`frontend/**`／`design/**` | 派单硬闸＋本票 §禁区 |
| 零码证明 | `git status --porcelain -- "*.go"` 中**已跟踪**脏件全部属别的在飞腿（`internal/tools/bridge.go`／`tool.go`／`system.go`／`system_23_test.go`／`cmd/wisp/panel_models_320_test.go`，`git status` 现读＝5 枚 ` M`），本腿自己**0 枚** | 见 §5.2 |

## 1. 结论：走**乙**

`AC#3` 二选一 → **乙＝把票 80 那两格转进本票＋具名作废原交接句**。甲（把该格写进票 21 的判据）**不选**，三条理由都能回指盘上原文，⛔ 一条是本腿的偏好：

1. **甲＝把一条已被否决的通道重新挂到能开工的票上。** 该格的字面行为是"同一份配置里写 override ⇒ 被点名文件从 B 变放行"＝**静态预授权**。
   - `A53②` 逐字：「**(A) 接静态预授权** = 契约从未承诺它生效，且票面 AC#3(a) 的形状恰是 `PLAN.md:2399` 禁止的那类（没人点就放行）⇒ 需 D22」。
   - `docs/PLAN.md:2407` 逐字：「**但豁免必须按文件、必须由人点、必须留痕。**」
   - owner 已答（`A448②` 逐字）：「**黑名单这个按你推荐吧**」；编排者照答裁：「**明确不做**（⛔ 推迟到"以后加固"）」。
2. **票 21 不是它的接收方。** 判据面 5 枚（`21-approval-gates-minimal.md:44/:46/:49/:50/:51`），治 L1 窗口／L2 卡／queue／batch／cancel／状态机；`blacklist_overrides` 命中 **0**（§2 表）。该格的裁决面（改契约＋owner 点）**不在票 21**。
3. **交接句自己认了"已取消"。** `80-blacklist-overrides-never-wired-to-gate-done.md:60-63` 逐字：「AC#3/AC#4 的接手方 = **票 21**（审批队列），裁决走 A53② 的选项 (C) ⇒ **本票永远不勾这两框**，它们是**被裁决取消**、不是**没做完**」。

⇒ 乙的落法：两格**接收在本票 `AC#3` 名下**，落地门＝本票 `AC#4`（owner 点／改契约）；原交接句**具名作废**，接收方＝本票。

## 2. THREE-CORE 前提复核（strings 在不在／file:line 对不对／别腿有没有动过）

| # | 前提 | 尺（现跑） | 读数 |
|---|---|---|---|
| P1 | 原交接句在票 80 | `grep -n 'AC#3/AC#4 的接手方' .scratch/wisp/issues/80-*.md` | 命中 **1 行**＝`:61`；上下文 `:60-63` 逐字见 §1 条 3 |
| P2 | 票 21 判据 5 枚、零提 override | `grep -c '^- \[ \]' .scratch/wisp/issues/21-*.md` ⇒ **5**；`grep -c blacklist_overrides` ⇒ **0** | 5 枚行号 `:44/:46/:49/:50/:51`（与票 80 `:49` 引的四个行号＋`:51` 同集） |
| P3 | 票 80 那两格今天的处置状态 | `sed -n '43,53p'` 现读 | 两格**已被摘框**（行首 `- ⛔` 而非 `- [ ]`），框后逐字「账在票 230——AC#3 与 AC#4 各对准一格」 |
| P4 | 响亮失败仍在（死线不沉默） | `grep -n BlacklistOverrides internal/config/unwired.go` | `:80` 路径／`:83` `fires` 逐字 `return len(c.Risk.BlacklistOverrides) > 0`；条目跨 `:79-84` |
| P5 | 别腿有没有动过票 230／80／21 | `git status --porcelain -- <三路径>` ⇒ **0 行**；`git log --oneline -1 -- <票 230>` ⇒ `da221145c`（2026-10-20 编排者收 `115-v1` 那笔） | 三张票面自那以后无人改；本腿落刀前 230 票面＝**53 行** |

### 2.1 两处 `file:line` **过期**（具名报上，⛔ 改写别人票面）

| 引用处 | 票面写的 | 现量（本腿现读） | 漂移 |
|---|---|---|---|
| `80-…:151`（D3 行）／`:166` | `docs/PLAN.md:2396-2399`／`PLAN.md:2399` | 「但豁免必须按文件、必须由人点、必须留痕。」在 **`docs/PLAN.md:2407`** | **+8** |
| `80-…:49`／`:50`（`done-fix-1` 追加） | `230-…md:21`／`:22` | `AC#3` 在 **`:27`**、`AC#4` 在 **`:28`**（落刀前） | **+6** |

⇒ 两张票面都按本仓 append-only 处理（原句不抹），本件只是把**正确锚点**记下来，供下一位按图索骥。

## 3. 乙的两半（落点逐字）

### 3.1 前半＝两格接收（不新增勾框，落点＝本票 `AC#3` 名下）

接收票 80 的两格**字面内容**（原文 `80-…:43-53`）：

- **格一（票 80 `AC#3`）**：接线落地＋一条**真跑在 `risk.Gate` 上**的用例：(a) 被点名文件由 B 变放行；(b) 同目录其它文件仍 B；(c) override 形状须 C26 交回的 canonical 形式（放行侧只认解析后的形式）。
- **格二（票 80 `AC#4`）**：变异双向——(i) 断线 ⇒ 格一必须红；(ii) 把 (b) 的断言指向 override ⇒ 必须红；两次都须在同一条 `&&` 链里先 grep 证改动落地再跑，还原后 `git diff --quiet` 证干净。

**状态＝被本票 `AC#4` 挡**（静态放行＝改契约、须人点）：

| 分叉 | 判据 |
|---|---|
| owner 不批（今天） | 两格永不落地；按 `A53②`(C) 以"响亮失败"结（`unwired.go:79-84` 保持响） |
| owner 批（＝先改契约） | 在本票实现这两格；届时才谈翻框 |

### 3.2 后半＝具名作废（原句不抹）

- 被作废句：`80-…-done.md:61` 逐字「AC#3/AC#4 的接手方 = **票 21**（审批队列）」。
- **自本条起作废**，接收方＝本票 `AC#3`（落地门＝`AC#4`）；⛔ 改派成"票 21 的新活"。
- **票 80 正文那一行的一枚内联注记本腿不改**（写面只许本票）⇒ 留给编排者一笔，建议逐字插在该行后：
  `（2026-11-05 作废：接收方＝票 230 AC#3，见 docs/evidence/s1/230-ac3-handoff-disposition.md）`
- 为什么"下一位按这句找得到 230"今天成立：`80:43`／`80:50` 两个框头的「⇒ **指向**＝票 230 `:22`」已经在盘上。

## 4. RED → GREEN（落刀前后两把现跑尺）

**RED（落刀前，HEAD `ba0ea5121`）**

| 尺 | 读数 |
|---|---|
| `grep -n 'AC#3/AC#4 的接手方' 80-…md` | `:61` 命中 **1 行** |
| `grep -c '^- \[ \]' 21-…md` | **5** |
| `grep -c blacklist_overrides 21-…md` | **0** |
| `wc -l 230-…md`（落刀前） | **53 行** |
| `grep -n 'AC#3' 230-…md` 的判据行 | `:27`（`- [ ]` 未勾） |
| `grep -c '接收在案' 230-…md` | **0**（票面零接收/作废记录） |

**GREEN（落刀后）**

| 尺 | 读数 |
|---|---|
| `git diff --numstat -- 230-…md` | **8 增／0 删**（纯追加尺；其中 `AC#3` 名下一段 7 行、票尾 Progress log 1 行） |
| `grep -n '处置＝乙'／'两格接收在案'／'自本条起作废'` | `:28`／`:32`／`:33`（各 1 行命中） |
| 票 230 行数 | 53 → **61**（`wc -l` 现读；与 8 增 0 删一致） |
| 票 80／票 21 `git status --porcelain` | **0 行**（`git diff` 空＝一字未动） |
| 本腿改动的 `.go` | **0**（`git status --porcelain -- "*.go"` 中无本腿路径） |

## 5. 门禁（共享树 + 纯净树两套读数，⛔ 把别人的半保存态算成本腿的绿）

### 5.1 纯净树（`git archive HEAD` 导出到仓外 `/tmp/230ac3-clean`；本仓规矩：要纯净树用 archive）

| 命令 | rc | 读数 |
|---|---|---|
| `go build ./...` | **0** | 无输出 |
| `go test ./internal/config ./internal/agent/approval -count=1 -timeout 200s` | **0** | `ok internal/config 4.326s`；`ok internal/agent/approval 1.632s` |

### 5.2 共享树现量（别人的脏件具名，⛔ 计本腿）

| 命令 | rc | 归因 |
|---|---|---|
| `go build ./...`（10:2x） | **1** | **未跟踪**的 `build/mut_rvw.go`（91,003B，mtime 2026-11-04 10:26，`git ls-files --error-unmatch` rc=1）＝别腿遗留突变样本，`undefined: voiceTtsFiles` 等 11 处 ⇒ `build/` 不是本腿写面、按"临时件只建不删"不清理 |
| `go build ./cmd/... ./internal/...`（同刻） | **0** | 绕开 `build/` 后全绿 |
| `go test ./internal/config ./internal/agent/approval`（10:3x） | **1** | 在飞写腿**半保存态**：`internal/tools/bridge.go:635:17: undefined: bridgeImages`（该文件 mtime 10:34、`git status` 显示 ` M`，本腿不碰 `internal/tools/`）⇒ 两包 `[build failed]`，⛔ 本腿缺陷 |
| `sh scripts/d22scan.sh`（10:29） | **0** | `clean - no D22 ban violations`；射程现读：`internal/` 410 枚生产 Go（bans #1-5 354）、`cmd/` 56、`frontend/` 116、`internal/tools/` 47、`design/` 39、ban #8 `internal/` 872＋`cmd/` 203 |

### 5.3 格式

本腿**零 `.go` 改动** ⇒ `gofumpt -l <changed>` 的射程是**空集**（不是"跳过"）：`git status --porcelain` 里本腿两条路径都是 `.md`。

## 6. 本件⛔ 做的事（边界，逐条）

1. ⛔ 翻 `AC#3` 框（归编排者）；⛔ 改票 80／票 21 一字（写面硬闸；票 80 那行注记建议句在 §3.2）。
2. ⛔ 动任何 `.go`／`PLAN.md`／`docs/specs/**`／`thresholds.go`／golden／`allowlist.txt`。
3. ⛔ push（子代理只 commit；推送是编排者的活）。
4. `AC#4`（owner 面）／`AC#5`（`cmd/wisp` 死亡区）**不归本腿**，照票面原样。
