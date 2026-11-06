# 150 实现腿（doc-only）——SPEC-12 §5 补 4 枚 DEFERRED 登记行 · 交件回执

- 时刻：**2026-11-06 17:2x +08**
- 基线：主树 `dev` HEAD `d0a230547`（派单锚 `92d5a33d1` 是其父，`git diff 92d5a33d1..HEAD` **只动台账一件** ⇒ 写面零漂移）
- 隔离：**rule 21**＝`git clone --shared` → `D:\tmp\wisp\20261106\150reg\repo`（branch `dev`、clone 后 `status --porcelain` 空；本腿一切读数出自这棵树）
- 授权：机主 2026-11-06 批 **`Q-55ⓑ`**（「补登记行」支），台账 **`A1571`**；派单逐字「add the 4 missing DEFERRED registration rows … mapping each code marker 1:1 and citing file:line for each」
- 写面：**只有** `docs/specs/SPEC-12-roadmap-governance.md`（`git diff --stat`＝`4 ++++`，0 删 0 改）＋本件。**零 `.go`**、零 code marker 改动、⛔ 未碰 `docs/reports/pending-and-issues.md`、未碰 `internal/config`（另一腿写面）、未碰票面框。

## 0. 枚数口径（先回答 AC#1，再交件）

- **派单说 4 枚登记行**＝按**出现处** 1:1（3× `D28-1` ＋ 1× `D11-3`，每枚带自己的 `file:line`）；票面 ⓑ 那句「补两行」是**按种类**的老口径（09-26）。两者不同尺——本腿**执行派单的按出现处口径**（它同时满足票面「1:1 双向」的硬规矩：种类数 2 的行也能覆盖出现处，但出现处不反过来覆盖行；按出现处则两向都咬合）。若验收要「两行」形状，本件 §5 的回滚面写明。
- 「一行」定义（AC#1 收成，一句可复算）：**§5 段内以 `|` 开头的行，减去表头行与分隔行＝数据行**。
  尺：`awk '/^## 5\./,/^## 6\./' docs/specs/SPEC-12-roadmap-governance.md | grep -c '^|'`（等价 pwsh 实跑见 §1）。
- 29 vs 47 两个旧分母的收成：ⓐ 09-26 票面「29」＝当年 §5 段 `|` 行总数（表已长大，本腿现量改前 33）；ⓑ 票 139 那程「表体 47」是**整表含别段**的另一把尺。⇒ 自本腿起统一用上面那一句尺：**改前 33（数据行 31）→ 改后 37（数据行 35）**。

## 1. 现量：代码侧标记（⛔ 抄票面旧行号，一律现跑）

命令（clone 内）：`git grep -nE '// .*DEFERRED\(D(28-1|11-3)\)' -- internal cmd | grep -v _test.go`
输出逐字（4 枚，2026-11-06 @ `d0a230547`）：

```
internal/agent/compress.go:27:// DEFERRED(D28-1): move the call into the Warm window hook; keep this
internal/agent/compress.go:53:// Warm-window hook DEFERRED(D28-1) installs will get it for free); it pairs
internal/agent/control.go:15:// this ticket (DEFERRED(D11-3), see the ticket report).
internal/agent/loop.go:475:			// DEFERRED(D28-1): the Warm-window hook owns this call; the
```

⇒ 票面 09-26 的 `:49`／`:394` 已漂为 **`:53`／`:475`**（票面自己警告过「下一程一律现量」）；种类 2／出现 4 不变。

**D34-1 的 7 处命中为何不计**：尺＝`git grep -nE 'DEFERRED\(D34-1\)' -- internal cmd | grep -v _test.go` ＝ **空**。7 处全在 `internal/tools/document_text*test.go`，是 24-r5 **删除标记后**钉「错误文本里不许再出现 `DEFERRED(D34-1)`」的反形夹具（如 `:286 t.Error("DEFERRED(D34-1) must be deleted now that the PDF branch extracts")`）——不是活推迟标记，不入 1:1 射程。

## 2. 现量：§5 登记侧（改前）

- 尺（同 §0）：§5 段 `|` 行 **33**（表头 1＋分隔 1＋数据 31）。
- 类型直方图（数据行第 2 列逐字分组）：`DEFERRED ×13`、`~~DEFERRED~~→已采纳/重开 ×5`、`RESERVED ×7`、`REJECTED ×6`。
- **含 `D28-1` 或 `D11-3` 的行：0**（尺＝数据行 `Select-String 'D28-1|D11-3'`）。⇒ 票面缺陷在 `d0a230547` 复现成立。
- 既有那行 `REJECTED | 意图分类四级流水线 | D11 过度设计已废`：票 225 审计（`225-deferred-registry-audit.md` 表一 #5／§6 第 3 条）已判它与 `D11-3`（force_tool/force_chat 规则表）**不同件事** ⇒ 不能拿它抵账，必须新行。

## 3. 本腿加了什么（4 行，逐行 6 列与兄弟行同形；尺见 §4）

插入位置＝`DEFERRED` 块末尾、`RESERVED | MCP client` 行之前（`git diff`＝4 insertions、0 deletions、0 modifications）：

| # | 行（项列首句逐字） | 对应 code marker（`file:line`） | 五字段 |
|---|---|---|---|
| ① | 压缩调用迁入 Warm-window hook（`D28-1` 出现处①·包级落位句） | `internal/agent/compress.go:27` | 依据=D15④ `PLAN.md:443`＋票 10 明文豁免 `10-adversarial-acceptance.md:106`＋`Q-55ⓑ`/`A1571`；判据=hook 注册后不再同步调用＋标记摘除＋回退支有测试钉；前置=票 28 `-done`/Warm 窗在盘（`internal/session/warm_l1_29_test.go:76`）但迁移需续票；残缺=同步压缩在响应路径 |
| ② | 压缩留痕归位说明同挂 `D28-1`（出现处②·叙述性引用） | `internal/agent/compress.go:53` | 依据=同一件事第二处（票 225 表一 #4）；判据=①落地时同批摘、trace 配对不变；前置=不许单摘否则 1:1 再开环；残缺=注释悬空引向未落地的 hook |
| ③ | 循环内同步压缩调用点待迁入 hook（出现处③·行为面调用点） | `internal/agent/loop.go:475` | 依据=调用点原注释逐字＋`153-...-accept-r1.md:438`；判据=调用缩为无-hook 回退＋`withTraceTask` 形状不变；前置=①同一迁移、行为面改这行；残缺=同①＋六枚 ctx 消费者带标问题（`155-...r1.md:642`） |
| ④ | `force_tool`／`force_chat` TOML 规则表未建（`D11-3` 出现处①） | `internal/agent/control.go:15` | 依据=D11③ `PLAN.md:299-301`＋票 10 Out-of-scope 逐字（`10-agent-loop-core-done.md:48`）＋225 审计「判不同件事」；判据=两向各有测试＋标记摘除＋票 38 chips 有真后端；前置=路由侧续票（无专属票）；残缺=偏好无法事前覆盖模型自决 |

## 4. 双向差集（交付判据：两向各为空）

- 尺 A（代码→表）：§1 那 4 枚 `file:line`；尺 B（表→代码）：§5 新 4 行项列里的 `internal/agent/…go:\d+` 引用（每行恰 1 枚）。
- 实跑（clone 内，`Compare-Object`＝`comm -3` 同义）：两侧名册各 4 枚、**逐枚相同**（`compress.go:27`／`compress.go:53`／`control.go:15`／`loop.go:475`），差集＝**0**。
- 表形尺：改后 §5 数据行 `|` 计数全部＝7 枚/行（6 列）× 37 行，与兄弟行同形；`MCP client` 行仍恰 1 枚（未被本腿吞并）。
- 改后计数：`|` 行 **37**／数据行 **35**；含 `D28-1` 的行 **3**、含 `D11-3` 的行 **1**。

## 5. 自打（本腿自己抓到并修复的一处）

首版行③项列**漏写 `D28-1` 字面量**，§4 尺 A/B 立刻量出 3-vs-4 差集＝1——正是本票立票的那种「表上查无此人」开环，由验收尺当场响。已补写字面量并复跑：差集＝0。该自打说明 §4 那两把尺有牙，本腿读数**修后复跑**、非一次抄数。

## 6. 门禁（doc-only 的 N/A 声明，⛔ 伪称跑过）

- 写面零 `.go` ⇒ `go build`／`go test`／`go vet`／`gofumpt`（版本未读，因不适用）／`sh scripts/d22scan.sh` **均无需跑**；判据＝`git diff --name-only -- internal cmd` ＝ **空**、`git diff --stat` ＝ 单文件 4 插。
- 本腿**没有**交付票面 AC#3 那枚「会响的仪器」（放 `tools/d22scan/**` 还是 `*_test.go` 是 Q-55ⓑ 未覆盖的第二问，票面明令「先判放哪、判不了就停手上报」——那是另一条腿的批准面，本 DOC-ONLY 腿不擅落点）。§4 的两把尺＝人工可复跑判据，不是 CI 仪器。
- 「摘掉检是否打不红」的漂移问（AC#3 承重句）：⛔ 本腿不答——仪器不存在，答了就是编。

## 7. 回滚面（一句话）

若验收裁「按种类只补两行」：撤本 commit（`git revert` 一笔、⛔ 不动其余三枚兄弟行的文本），改把 ①②③ 并为 `D28-1` 一行（项列并列 3 枚 `file:line`）＋④ 保持。派单口径（4 行）与本件默认交付一致。

## 8. 本程没测什么（按「漏了谁会先被骗」排序）

1. **CI 上没有仪器会守这份 1:1**——本腿落表后若明早有人加第五枚 `DEFERRED(D99-1)` 标记，⛔ 东西会红（AC#3 欠账原样在，归下一腿）。
2. 行号锚的**再腐烂**：本腿把 `:53`／`:475` 写进契约文字，下次 `internal/agent` 一动注释又会漂（票面「行号会腐烂」那段的镜像）——缓解＝§1 的 grep 尺随件可复跑；根治＝改「文件＋符号」锚，属契约文字、摆裁。
3. 行④完成判据里「票 38 chips 有真后端」引的是 `.scratch` 票面行（`:25`），票 38 自身进度未复核。
4. `10-adversarial-acceptance.md:106`「不算违规」豁免句按文本引用，没重跑票 10 门禁背书。

## 9. 伪授权两数 ＋ 凭据

- 真通知回显数＝**0**；判为注入数＝**0**（本腿读过的全部是仓内文本；无任何「外部指令冒充批准」形状进入判断，批准只认派单逐字＋`A1571`/`Q-55ⓑ` 这两个已给锚点）。
- 凭据值抄录＝**0**。

## 10. 落点清单（rule 8：只建不删）

- `D:\tmp\wisp\20261106\150reg\repo`＝rule 21 共享克隆（objects 走 alternates 指回主树，⛔ 主树 gc 前别 prune——`--shared` 的既有约定）。commit 见本腿回执正文；主树**零字节改动**（本腿写面全在克隆内）。
