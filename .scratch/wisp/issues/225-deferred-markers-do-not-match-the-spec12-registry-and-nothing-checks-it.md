# 225 — **`DEFERRED(...)` 代码标记与 `SPEC-12 §5` 登记表今天对不上，而且没有任何仪器在查这一条**：`AGENTS.md` §1.1 那条"1:1 双向"目前是**只靠人**的规矩

- Status: **待派（只读审计腿即可，不产码）**。来源：票 223 的普查腿 `223-c1` 顺带撞出、编排者现跑复认（台账 `A426`）。
- ⚠ **为什么现在要管**：票 223 之所以差点被写成"我们漏做了热加载"，就是因为**代码里那枚标记说这件事"由票 05 实现"，而票 05 早就结案了、热加载在生产里根本不转**。标记在骗人 ⇒ 后面的程要么重复造、要么以为有人排过。

## 现量（起手逐条复算，别信这里的行号）

| 事实 | 读数（锚 `36b46125`，2026-09-29 11:57 现跑） | 尺 |
|---|---|---|
| 代码侧标记枚数 | `grep -rn "DEFERRED(" --include=*.go internal/ cmd/ tools/ \| grep -v _test` ＝**26 处命中／22 枚去重标记** | 现跑 |
| 登记表侧枚数 | `docs/specs/SPEC-12-roadmap-governance.md` §5 那张表：`^\| DEFERRED`＝**12 行**、`^\| RESERVED`＝**7 行**、`^\| REJECTED`＝**6 行**（另有 2 行是 `~~DEFERRED~~ → 已采纳`） | 现跑 |
| 差 | **22 vs 12 ⇒ 至少 10 枚代码标记在 §5 找不到同名行**（差数只是线索，**逐枚映射要人做**：有的标记是同一件事的第二种拼法、有的可能对应 §5 用了另一个名字） | 现跑 |
| 一枚确凿的失配 | `internal/config/doc.go:14` 逐字 `DEFERRED(schema/hot-reload/migration): implemented by ticket 05. This`，而 `.scratch/wisp/issues/05-config-model-done.md` **已结案**、`CheckAndReload` 的**唯一非测试调用者仍在 `cmd/balldebug/main.go:243`**（见票 223）⇒ **"标记声称已由 05 实现"与"生产里没人调"两句同时为真** | 现读 |
| 另一枚确凿的失配 | `internal/watchdog/doc.go:18` `DEFERRED(watchdog loop/thresholds): implemented by ticket 42`；而 `grep -n "看门狗\|watchdog" docs/specs/SPEC-12-roadmap-governance.md` 只命中 S6 那行**切片表**、**不在 §5 登记表** | 现读 |
| 没有仪器 | `grep -rln "SPEC-12" --include=*.go tools/ internal/ cmd/` 只命中 `internal/agent/approval/` 两枚文件里的注释；`tools/` 里没有任何程序读 `DEFERRED(` ⇒ **`AGENTS.md` §1.1 那句"必须 1:1 双向对得上"今天没有自动检查**（`tools/d22scan` 扫的是代码形状，不是登记表） | 现跑 |

## 判据（每格都要现跑读数）

- [ ] **AC#1 逐枚映射表**：22 枚去重标记逐枚列出（文件／行／标记名／它自称由哪枚票实现／那枚票今天什么状态），并各自判定三态之一：**§5 有对应行**／**§5 没有该补一行**／**标记本身过期了（自称已实现且确已实现，该摘掉）**。⛔ **不许改动 `docs/specs/**` 任何文字**——要补行／要摘标记都写成**待人批准的落点**（改契约＝人工批准；`DEFERRED(D-xx)` 与 §5 的 1:1 本身就是契约面）。
- [ ] **AC#2 反方向也走一遍**：§5 那 12 行 DEFERRED＋7 行 RESERVED 逐枚问"代码里有没有对应标记"，缺的具名列出。⚠ 方向二是本票的主要价值——**"登记了但代码里没标记"的项最容易被后面的人当成"没这回事"**。
- [ ] **AC#3 那两枚确凿失配要有结论**：`internal/config/doc.go:14`（热加载，见票 223）与 `internal/watchdog/doc.go:18`（看门狗循环）。各写一句"该改标记还是该补登记"，并具名指出改它要不要人批准。
- [ ] **AC#4 判"要不要上仪器"，不要自己上**：本票**只出结论与形状**——如果要把这条做成常驻检查（像 `tools/d22scan` 那样），那是新增一台仪器＝新射程，**要先摆给 owner**；不许本腿自作主张写进 CI。

## 禁区

不动 `PLAN.md`／`docs/specs/**`／`thresholds.go`／golden／`allowlist.txt`；不动三枚冻结件（`internal/panel/tokens_fourway_test.go`／`internal/panel/l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go`）；不新增导出名；`frontend/**`／`design/**` 零读零写零转述；⚠ **本票与票 223 同源但射程不同**：223 管"接线"，225 管"账目对得上"，**不许把两票的合成一票做**。

> ✅ **09-29 12:2x 已派已收（腿 `225-a1`，`docs/evidence/s1/225-deferred-registry-audit.md`，46,705 字节，提交 `dd7c447f`＋`3397dbd0`）；编排者现读复验，台账 `A429`。三格结论落定＋两枚比原票面更硬的补充：**
> - **表一**：〔§5 有对应行〕**3 枚**／〔该补一行〕**3 枚**（`queue`／`audio, ticket 09`／`update/rollback`）／〔该摘或改写〕**12 枚**（14 处，占 54%）。⚠ **"1:1"这个措辞本身今天不成立**——`macOS/Linux`＋`P12-macos` 两枚标记共用 §5 同一行＝**N:1**。
> - **表二（反向）**：12 行 DEFERRED 里**只有 2 行**有代码标记 ⇒ **10 行"登记了而代码查无此人"**，最危险 5 行＝剪贴板历史／`doc.read` xlsx-OCR／`system.eject`／完整错误文案／i18n（编排者用自己那份 22 枚名册交叉核过：这五件事一枚标记都没有，而相邻代码是活的）。
> - ⛔ **补充①（编排者现跑）：这条规矩规定的载体从来没建过**——`SPEC-12:94` 逐字要求标记写成 `// DEFERRED(D-xx): … → docs/DEFERRED.md#锚点`，而 `ls docs/DEFERRED.md` ＝ **不存在**，且 22 枚标记**没有任何一枚**带那一段后缀。⇒ **补建 `docs/DEFERRED.md`＝动 `SPEC-12 §6` 交付物清单＝契约面＝要人批**；**默认动作＝不建**，只把这一条摆在这张票的待人批清单里。
> - ⛔ **补充②（编排者现跑复现）：一处类型词与登记表相反**——`internal/tools/registry.go:55` 逐字 `{Kind: KindMCP, Ticket: "REJECTED: D13/16.9#7, interface slot only"}`，而 §5 把 MCP client 登记为 **RESERVED**；两者给后续人的指令相反（`REJECTED`＝不得被修复，`RESERVED`＝只留接口位）。⇒ **不由编排者也不由腿裁**，进待人批清单。
> - **表三归口**：`config/doc.go:14` 该改标记、落点归**票 223**；`watchdog/doc.go:18` 该改标记，但"watchdog 建一圈 vs `cmd/wisp` 自己起"是拓扑选择、牵 D25/D38 ⇒ **归编排者＋owner**。
> - **AC#4 按原意按住**：正向检查要先有一枚人工"别名字典"＝新真相源 ⇒ **本轮不建仪器、不进 CI**；并记一枚实现坑：按行 `grep "implemented by ticket"` 只命中 9 枚、**注释折行漏 5 枚**。
> - ⚠ **本腿留下一枚共享索引事故（详见 `A429`）**：第一次 `git commit` **没带 pathspec** ⇒ `dd7c447f` 把在飞腿 `222-r1` 已 `add` 未提交的 `internal/tools/subagent_222_test.go` 一起收走（盘上量得零丢失，不改写历史；`dd7c447f` **不作** `internal/tools` 的"已核过绿"锚点）。⇒ **今后派单的 Git 纪律改成把 pathspec 写进 `commit` 本身**，并**报路径必须报到包级**（本腿正文写 `registry.go:55` 没写包名，真位置 `internal/tools/registry.go:55`）。

## 编排者追加（2026-10-04 15:4x，凭据＝票 265 的普查件 §3.3＋`A601` §4；⛔ 本票 AC 框一枚没碰）

- **一格归口到本票**：票 265（常驻腿的门没有会话授权记账位）原先有三形候选，其中 **ⓒ＝按 `SPEC-12 §5` 把它登记成 DEFERRED**。编排者**没有选它**，而把它**归口到本票**——理由逐条：
  1. `docs/specs/SPEC-12-roadmap-governance.md:95-96` 逐字规定代码内标记的形状要带 **`→ docs/DEFERRED.md#锚点`**，而 **`docs/DEFERRED.md` 在盘上不存在**（尺＝`git ls-files docs | grep -E "DEFERRED|DECISIONS"`＝**空**；`AGENTS.md` §4 末段也具名说过这批交付物"截至锚点 `4e66817` 在仓里不存在"）。
  2. ⇒ **照规格那句写出来的 ⓒ 标记，会被 `tools/d22scan` 的 ban #9（phantom-citation，`main.go:799-873`／尺路径正则 `:890-891`）直接判红**——它判的正是"产码注释引一条盘上不存在的仓内路径"。⚠〔仅读码，**未跑** d22scan 复认这一发；真跑属票 265 的 U6，⛔ 不许被下一任读成实测〕。
  3. 今天仓里的 `DEFERRED(...)` 标记**没有一处带 `docs/` 引用**⇒ 这枚撞是** ⓒ 一形新造的**，不是既有状况。⚠ **枚数两组成员，不写成"新旧版"（第 75 条）**：普查腿 `265-a1` 报 **28 处**〔腿报，它的尺名册 R111，口径我没能对上〕；编排者 **16:39:14 自己现跑**＝`git grep -o "DEFERRED(" -- cmd internal tools | wc -l`＝**30 处**，其中带 `D-<数字>` 那种规格拼写的＝**0 处**（⇒ 现存的标记用的是别的实参形状，正是本票"标记与登记表对不上"那颗钉子），同行含 `docs/` 的**只有 1 行**、而那一行是 `internal/risk/syncdirs_other.go:8` 的一句散文（提到 `docs/PRECHECK.md`），⛔ 不是一枚指向 `docs/` 的 DEFERRED 标记 ⇒ **"零枚标记带 docs/ 引用"这一句由我复认成立**，⚠ 而 **28 与 30 差 2 枚这件事没结**，归本票 AC#0 那格自己拉名册时对（⛔ 不许照抄任何一枚数）。
- **⇒ 本票要答的那一格因此变大**：不只是"标记 ↔ `SPEC-12 §5` 登记表 1:1 双向对账"，还包括**"规格要求的那个锚点指向的文件到底建不建"**——建＝新真相源（`docs/specs/**` 或根目录交付物，属编排者／人工批准面）；不建＝`SPEC-12:95-96` 那句形状规定**与盘上不可执行**，要改规格文字（＝**契约变更，人工批准**）。⛔ 两支都不许由产码腿自行选。
- ⚠ **本票仍不是任何落地腿的顺手活**：票 265 走的是 ⓐ-Ⅰ（接上写侧），**不需要 ⓒ**；这一格留在本票是因为**它本来就还没人管**（`A601` §4 具名登记），⛔ 不许被读成"票 265 已把 DEFERRED 那套理顺了"。
