# 146 — `Queue.LiveApprovals()` 的注释说"返回副本"，实现是**按值拷贝一个含引用字段的结构**：那张 map 与那三枚切片都和队列共用底层

- Status: ready-for-agent（**未派**）
- 来源：票 35 快照泵收尾程 `docs/evidence/s1/35-panel-snapshot-pump-fix-r1.md` 留给下一程的空格 ②；
  它**没改行为、也没在两个方向任一处下断言**，所以这一格今天仍是空的（原话："这格今天仍是空的"）。
  ⚠ **本票面把它写大了**：它只点了两枚切片，编排者 18:2x 逐字段读被拷贝的结构后是**一枚 map ＋ 三枚切片**（见下）。
- 关联：票 35（`LiveApprovals` 就是它那一批里新增的出口）、票 33（泵将来要往页面送这份东西）、
  票 145 `AC#3`（"当前屏"那一维的归属裁定在这里会被引用到）

## 现量的形状（三条命令可复算）

- `internal/agent/approval/pending_read.go:19-20` 的注释逐字写着
  **"It returns copies of verdicts the assessor already reached"**（"copies" 那一个字在本行行首）。
- 被返回的结构是 `LiveApproval{CorrelationID string; Decision tools.Decision; Position int}`
  （`pending_read.go:30-34`），实现里那一行是 **`Decision: it.Dec`——按值拷贝整个 `tools.Decision`**。
  ⇒ 但 `tools.Decision`（`internal/tools/gate.go:15-24` 起）**自己就带引用类型字段**：
  | 字段 | 类型 | 按值拷贝之后 |
  |---|---|---|
  | **`Params`** | **`map[string]any`**（`gate.go:18`） | **同一个 map**——写它会直接改掉队列里那条记录，**这比切片更硬** |
  | `RulesHit` | `[]risk.RuleID`（`gate.go:21`） | 共用底层数组 |
  | `Args` | `json.RawMessage`（`gate.go:19`，本质是 `[]byte`） | 共用底层数组 |
  | `Paths` | `[]string`（`gate.go:23`） | 共用底层数组 |
  | `Capabilities` | `[]Capability`（`gate.go:24`） | 共用底层数组（元素本身还有没有引用字段，本票没量） |
  ⇒ **"copies" 许诺的是深拷贝，实现给的是"浅一层"**。
- **今天读它的是谁（全量现量，含测试）**：`LiveApprovals(` 共 **11 处**＝1 处定义
  （`pending_read.go:42`）＋ **1 处生产调用**（`cmd/wisp/panel_pump.go:61`）＋ 9 处本包测试（`pending_read_test.go`）。
  而 `grep -rn "RulesHit =\|RulesHit\["`（剔 `_test.go`）在全仓只命中 3 处，都在别的地方构造/搬运 verdict、
  **没有一处是"拿到 `LiveApprovals()` 返回值之后就地写它"**。
  ⇒ **所以今天不是缺陷，是一个还没被触发的形状**——这句要写死在结论里，别读成"有个 bug"。
  ⚠ 上面那三条是我 18:2x 在本机工作树量的；**下一程一律现量，别抄我的数**。

## 为什么现在不许顺手"改成深拷贝"就完事

这一格有两种合法收法，**方向相反、代价不同**，本票要你**先判再改**：

- **ⓐ 改实现**（把那枚 map 与那三枚切片真拷一层）——买的是"注释说的字成立"；代价＝每次泵动一下都多一张 map＋三份底层数组，
  而**这张表是要频繁重算的**（`wisp run` 每有面板状态变动就现算一份），所以它**撞 `PLAN.md` 的 D32 资源预算那一族**。
  ⇒ 选 ⓐ 必须**同时量一次**：多拷这些东西，`wisp slo` 的 WorkPeak 那一档数字动没动。没量过就不许写"代价可忽略"。
- **ⓑ 改注释**（把 "copies" 换成"逐条复制结构体本身；**`Params`／`RulesHit`／`Args`／`Paths`／`Capabilities` 与队列共享底层，调用方不许就地写**"）
  ——买的是"注释不再许诺它没做的事"；代价＝把纪律放在调用方身上，**而调用方里有 `internal/panel` 与将来的页面**。

⚠ **本票的默认方向是 ⓑ，但 ⓑ 单独存在不够**：一句"调用方不许就地写"是**没有仪器的纪律**，
而这一族东西一旦被谁忘一次，症状是**队列里那条已决断记录被悄悄改掉**——那正好落在
"面板侧不许变成授权来源"那条硬禁线的旁边。⇒ **无论选 ⓐ 还是 ⓑ，都必须有一枚会响的检**（AC#2）。

## AC

- [ ] **AC#1 先判 ⓐ／ⓑ，判不了就停手上报**（D22 闸门③，不许自行假设）。判据是**可核的一句**：
      今天有没有**任何**一枚调用方（含测试）在拿到返回值之后对那枚 map 或那三枚切片做 改元素／append／就地写／取地址后传出去？
      现量命令自己写进证据件。⇒ 全为"无"时 ⓑ 是诚实且零代价的一支；只要有**一枚**为"有"，ⓑ 就不许单独收，必须 ⓐ。
- [ ] **AC#2 一枚会响的检（本票的硬核心）**：造一发**只有"就地改掉返回值里的 `Params` 或那三枚切片"才能触发**的用例，
      断言队列里那条原始记录**没有**跟着变。判"承重"按本仓操作定义答一句：
      **把 AC#1 选定的修法摘掉，这一发是不是从此打不红？** 答"仍不红"＝这枚检是装饰，本票作废并登记为什么。
      ⚠ **不许**用"永远绿"的那种断言凑数（本仓那族缺陷有名字：**恒真判据是一类新假绿**）。
- [ ] **AC#3 若走 ⓐ，资源那一问不许空着**：给出"多拷那枚 map 与那三枚切片"在 WorkPeak 档的前后读数
      （`wisp slo` 或 `scripts/slo-check.ps1` 的同形发法），**阈值与 golden 一字节不许动**，
      数字没动要写"未观察到差异"而**不是**写"无代价"。
- [ ] **AC#4 契约轴**：`internal/risk/**`、`rules_gateway.go`、`thresholds.go`、golden、`allowlist.txt`、
      `docs/PLAN.md`、`docs/specs/**`、`frontend/**`、`design/**`、`tools/d22scan/**`、
      `internal/panel/l2_grant_boundary_test.go` **一律零字节**。
      ⚠ 这枚函数**不在 `PanelAPI`／C17 名册上**、不能回答从页面来的请求（`pending_read.go:14-18` 自己就这么写）⇒
      **本票不许把它变成一条入站路由**；那一变要重新走契约批准（`Q-51` 那一族）。
- [ ] **AC#5 门禁（一律逐包单跑）**：`go test ./internal/agent/approval/` 改前改后各一次
      （四数之外**点名册差集**；本目录前八枚测试文件是 `package approval_test`，
      而 `7040493` 那枚是 `package approval` **白盒**——两种都有分母，别把它们当成同一棵尺）。
      ⚠ **两包合跑会凭空造红**（`TestAC1ResidentLegBooksItsShutdownBeforeClosingTheSink` 报 `0xc000013a`、
      单跑就过、那文件与本票无关），这条是本仓 09-25 现量的仪器形状。`gofmt -l` 空、`go vet` 带你动过的包、
      `sh scripts/d22scan.sh` rc=0（独立 module；`set -eu`，**第一步是正控，正控红则真扫描根本不跑**）。

## 两条已知读数，别当结论用（原样交给下一程）

1. 收尾程在**锚点归档副本**里读到两枚 `TestSecret*` 红，**工作树单跑是 PASS**——它**没归因**，
   并明确写"别当'改前就红'引用"。⇒ 你若在副本里再见到它们，先做 `git diff <锚>..HEAD -- <那个包>` 对照再判。
2. 行号会漂：那份简报与验收件写 `cmd/wisp/run.go:741`，现树同一句 `changed = true` 在 **740**。
   ⇒ 本票引用的任何行号都**现量**，别从上面两枚文件抄。

## 规矩

显式 pathspec、禁 `add -A`/`.`/`--amend`/`reset`/`stash`/`checkout .`/`clean`、禁 push、
**临时件只建不删**（变异落仓外副本）、台账与勾归编排者、简报里任何一句前提复算不符 ⇒ 写进证据件并报回，
**不要按我认为的样子改**。
