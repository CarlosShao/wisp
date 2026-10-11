# 309-a1 · 21 形ⓑ「写明故意」（照 ② `cmd/wisp/panel_assets.go:164`）的代价（料，⛔ 裁）

尺与三态同前（HEAD blob @`7dc753e3`；`git status --porcelain -- internal cmd`＝0 行 ⇒ HEAD blob＝工作树；⛔ 编译面命令，全部静态判定）。

## 一、会动哪几枚文件（逐枚＋最小编辑形状）

| # | 文件 | 在哪枚位置 | 最小编辑形状 | 必动？ |
|---|---|---|---|---|
| b1 | `internal/tools/bridge.go` | 文档块 `:920-923`（现量逐字只讲"per scope, hence per task"）＋守卫行 `:925` 上方 | 在文档块里**追加**一段，内容＝那三枚前提逐枚具名 ＋ 空 id 那支的**后果**（R4 少报）＋ 为什么⛔ 删这一肢（`SrcUnboundScope` 地板）＋ 谁该为此响（该具名哪枚仪器/哪枚门）。纯注释，⛔ 动 `return` 语句、⛔ 新字面量 | **必动（唯一 1 枚）** |
| b2 | `internal/tools/*_test.go`（`AC#1` 那格） | 新函数 | ⓑ 本身⛔ 要求新仪器，但票面 `AC#1` 是独立一格（恒真形必须响）：ⓑ 落地后仪器钉的是"**注释说了真话**＋**摘掉守卫/改成恒挂会红**"两件事。所以 b2 属 `AC#1`，⛔ 记在ⓑ的账上，但派单时要说清ⓑ＋`AC#1` 是**同一程的两格** | 属 `AC#1`，非ⓑ独有 |
| b3 | 票面 309 末尾 append-only bullet | — | 同 `AC#4` 允许项 | 落地程 |

**枚数小结：ⓑ＝必动 1 枚（`internal/tools/bridge.go` 的注释），⛔ 生产码语义变化 ⇒ 与 `AC#1` 天然同片。**

## 二、会不会打红既有钉（静态点名）

- **静态判定：0 枚必红。**理由不是"注释无害"，而是**逐枚查过谁会把注释读进断言**：
  - 尺＝`git grep -n "bridge.go" HEAD -- '*_test.go'`（现量 30+ 命中，逐条看过）：命中的全部是**注释里引用** `bridge.go` 的行号/形状，**没有一枚 Go 测试在运行时读 `bridge.go` 的源码文本**。真正钉文字的测试族只有 3 枚常量，且钉的是别的分支：`cmd/wisp/panel_inbound_guards_35r3_test.go:65-67`（`不是面板 composer 通路的能力入口`／`按伪造/串台拒绝`／`缺少 requestId`）。ⓑ 的编辑⛔ 碰这三句。
  - 尺＝`git grep -n "os.ReadFile\|ReadFile" HEAD -- 'internal/tools/*_test.go'`（本腿未展开全量 ⇒ 归欠读数 U-3；静态已知的是：`internal/tools` 里没有一枚"源码文本钉"点到 `assessorFor` 的文档块，因为名册尺 `git grep -n "assessorFor" HEAD -- internal cmd` 现量 7 命中里测试只有 1 枚、且是注释行 `panel_assets_143_test.go:7`）。
- **⚠ⓑ 有一条自己的红法，具名登记**：`AC#1` 要求的仪器（"改成恒挂探测器／整枚摘掉守卫时指名用例必须红"）一旦落地，**ⓑ 的那段注释就成了被钉对象**。⇒ 注释与仪器必须同片写、口径一致；注释先写、仪器后补＝`AC#1` 那格交付不了（票面明令"⛔ 只交改⇒红不交换反形也红＝那格没交"）。
- ⛔ 撞 `72-r1` 面：ⓑ 只动 `internal/tools/bridge.go`，⛔ `internal/risk/*_test.go`（票面排程那条冲突在此消失）。

## 三、★这形里"注释怎么写才算说实话"——一枚可判据

> 病灶对照：台账 `A838` §2 那族（票 40＝面板 Security 页零消费者）的病是**用现在时给还没建的东西写注释**。ⓑ 最大的失败模式就是同一枚病：注释把"这一支是故意的"写成一句心情。

**判据（三条，缺一即那格没交；每条都配一把⛔ 跑 Go 的尺）：**

1. **必须含"少报"这个词，且被一把 grep 尺点到。**尺＝`git grep -n "少报" HEAD -- internal/tools/bridge.go` 命中 ≥1，且该行必须落在 `assessorFor` 的文档块范围内（现量范围＝`:920-923`＋追加行；尺＝`sed -n '/func (b \*Bridge) assessorFor/=;920,940p'`）。理由：这三枚前提里唯一**用户可感**的后果就是 R4 少报；写不出"少报"两个字＝没说出代价。
   最小可核形状示例（⛔ 本腿落码，只给形）：
   ```go
   // taskID == "" is a DELIBERATE under-report, not an oversight: with no task
   // there is no C25 scope to bind, so R4 stays dormant (internal/risk/
   // rules_gateway.go:102 `if ctx.taint == nil` short-circuits). The limb is
   // kept because removing it does NOT fail closed: provenance.go:694-710
   // reports SrcUnboundScope for an unregistered scope whenever any other
   // scope holds marks, i.e. a false-positive storm instead of a missing hit.
   ```
2. **每句必须此刻可核**（文件:行、`git grep` 命中的字面量、或常量名三选一），⛔ 出现任何"when/if the panel grows …""以后接上 X 就会"式的**未来时页面描述**。尺＝对新增行做反向 grep：`git grep -nE "(将会|以后|未来|once the|when the panel|will be wired|TODO)" HEAD -- internal/tools/bridge.go` 在 `assessorFor` 块内命中＝0。这条直接对着票 40 那族病。
3. **必须点名"谁该为此响"**（一枚门或一枚仪器的路径），⛔ 以"注释即文档"收尾。②的成色基准就是现量 `panel_assets.go:171-173` 那句逐字 `"Nobody ever closed this one" stays G5's job in .scratch/wisp/probes/154/gate-clauses.sh, not the compiler's`。⇒ ⓑ 合格版要写出等价的一句：要么指 `AC#1` 那枚新仪器（`internal/tools/<file>_test.go` 的具体函数名），要么指 `gate-clauses.sh` 的某枚 G 条款。尺＝`git grep -n "gate-clauses.sh" HEAD -- internal/tools/bridge.go` 命中 ≥1 **或** 注释里出现一枚 `_test.go` 具名函数。

## 四、用户侧表现

- **完全零变化**：ⓑ ⛔ 改判级、⛔ 改分支、⛔ 改文案 ⇒ 今天的判定与昨天逐字节相同（生产构造点 1 枚恒非空，本路径今天不可达；②那形本来就是 CLI 探针，⛔ 过桥）。
- 谁看到：**只有读代码的人**看到"这一支少报、为什么留着、谁该响"。**模型/用户/审计行都⛔ 变化**——这正是ⓑ的代价面：**它对空 id 的实际防御为零**，票面 `AC#2` 那枚"结构保证"在ⓑ里必须**另找一处**（见 `30-structure-guarantee.md`，ⓑ＋`AC#1` 仪器是ⓑ形里唯一能让"未来接线就响"落地的组合）。
- 与①的**不一致残留**：ⓑ 保留"三处三种处置"这个事实（①②静默但都写明、③④拒绝）；票面 `AC#0` 要的"对齐"在ⓑ里变成**只对齐注释成色**、⛔ 对齐行为。这是ⓑ的**待裁代价**，写清供裁形。

## 五、⛔ 顺手动到的东西

1. ⛔ 顺手改 `:925` 的**布尔式**（把 `taskID == ""` 挪走、合并前提、或"从严"删肢——删肢撞 `SrcUnboundScope` 地板，票面边界）。
2. ⛔ 顺手改 `:869-883` 那段票 154 的归属注释（那是 `AC#3` 的射程：只做**行号→内容锚**的替换＋枚数重述，⛔ 改其中任何**契约声称**文字）。ⓑ 的写面是 `:920-923`＋守卫上方，⛔ 越到 `:869-891`。
3. ⛔ 顺手补 `AC#1` 的仪器（那是另一格，且它的判据是"反形必须红"，混进ⓑ那格就两格都算没交）。
4. ⛔ 把注释写成"今天没有生产 caller"的**永久保证**（票面「边界」明令：⛔ 把这类现读句当永久保证；ⓑ 若引用枚数，必须带尺、带"现量"二字）。
5. ⛔ 动 `internal/risk/**`／三枚冻结件／SLO／golden／D43／C1–C32／D1–D47；⛔ 放宽任何断言；⛔ 造 SKIP。
6. ⛔ 翻框、⛔ 动台账/`docs/**`/票面正文（票尾 append-only 1 条除外，那是落地程的允许项）。
