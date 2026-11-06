# 150 非实现者裁决表 — SPEC-12 §5 补 4 枚 DEFERRED 登记行（D28-1 ×3 ＋ D11-3 ×1）

- 角色：**非实现者裁决腿**（裁决者≠实现者，`SPEC-12 §4.3` #1/#3）
- 时刻：2026-11-06 18:xx +08
- 裁决对象：commit `9b3cb56d7`（现于主树 `dev`，其上再一枚 `218ed2ee9` 与本票无关面）
- 本腿 HEAD 现量：`218ed2ee9`（`9b3cb56d7` 为其父）；实现腿派单锚与授权：owner 2026-11-06 批 `Q-55ⓑ`、台账 `A1571`（现量复核见 §3①）
- 实现腿回执：`docs/evidence/s1/150-deferred-registry-rows-impl.md`（本裁决逐条独立取证，**不照单采纳**）
- 本腿写面：**仅本件**；票面、代码、规格零改动；git 只 commit 不 push，显式 pathspec

---

## 1. 独立现量（全部本腿现跑，Git Bash @ 主树）

### ① 代码侧活标记（尺：`grep -rnE 'DEFERRED\(D[0-9]+(-[0-9]+)?\)' --include='*.go' internal cmd | grep -v _test.go`）

```
internal/agent/compress.go:27:// DEFERRED(D28-1): move the call into the Warm window hook; keep this
internal/agent/compress.go:53:// Warm-window hook DEFERRED(D28-1) installs will get it for free); it pairs
internal/agent/control.go:15:// this ticket (DEFERRED(D11-3), see the ticket report).
internal/agent/loop.go:475:			// DEFERRED(D28-1): the Warm-window hook owns this call; the
```

⇒ **4 枚活标记全部在实现腿所报的 file:line 上**（`:27`／`:53`／`:475`／`:15`，与回执 §1 逐字一致）；种类 2、出现 4。票面 09-26 的 `:49`／`:394` 旧行号确已漂移，票面自己警告过"现量"，成立。

### ② commit 形状（尺：`git show 9b3cb56d7 --stat` + 显式 pathspec 差集）

- `git show 9b3cb56d7 --stat` ⇒ **仅 2 文件、89 插入 0 删**：`docs/evidence/s1/150-deferred-registry-rows-impl.md`（新件 +85）＋ `docs/specs/SPEC-12-roadmap-governance.md`（+4）。
- `git diff --name-only 9b3cb56d7^..9b3cb56d7 -- internal cmd` ⇒ **空**（0 文件）。
- `git diff --name-only 9b3cb56d7^..9b3cb56d7 -- docs/reports docs/PLAN.md docs/specs/* thresholds.go tools scripts` ⇒ **仅 SPEC-12 一件**；台账 `pending-and-issues.md` 未被本 commit 触碰。
- SPEC-12 的 diff 逐字读过：**恰 4 行 `| DEFERRED |` 插入**（落点 §5 段内：`## 5.` 在 :59、`## 6.` 在 :105、4 行在 :84–:87），0 删 0 改，插在 `RESERVED | MCP client` 行之前——与回执 §3 声称的插入位置一致。
- `git diff 9b3cb56d7..HEAD -- docs/specs/SPEC-12-roadmap-governance.md internal/agent` ⇒ **空**（本票面自落地后零再漂移）；工作树对这两处 `status --porcelain` 干净。

### ③ 双向差集（尺：两侧名册各 `sort` 后 `comm -3`；**两向各为空 = 过**）

- 尺 A（代码→表）：§1① 的 4 枚 `file:line`。
- 尺 B（表→代码）：`awk '/^## 5\./,/^## 6\./' … | grep '^| DEFERRED'` 行内 `internal/…go:\d+` 引用（剔 `_test`）→ `sort -u`。
- 实跑两侧名册**逐枚相同**：`compress.go:27`／`compress.go:53`／`control.go:15`／`loop.go:475`；`comm -3` 输出**空**。
- **按出现处**逐枚对号（不止集合相等）：行①标题=出现处①·包级落位句 ↔ `compress.go:27` 注释（"move the call into the Warm window hook"）；行②=出现处②·叙述性引用 ↔ `compress.go:53`（"hook … will get it for free"）；行③=出现处③·行为面调用点 ↔ `loop.go:475`（"the Warm-window hook owns this call"——真调用侧）；行④=出现处① ↔ `control.go:15`。语义对号无错位。

### ④ 表形与计数

- 改前（`9b3cb56d7^`）§5 `|` 行 **33**、`D28-1` 命中 **0**、`D11-3` 命中 **0**（票面缺陷在基线复现成立）。
- 改后 §5 `|` 行 **37**（数据行 35）；类型直方图：`DEFERRED ×17`、`REJECTED ×6`、`RESERVED ×7`、`~~DEFERRED~~→已采纳/重开 ×5` ＝ 35 ✓（与回执 §2"改前 13+5+7+6=31"自洽：+4 全进 DEFERRED）。
- 含 `D28-1` 的行 **恰 3**、含 `D11-3` 的行 **恰 1**；4 行每行 **7 枚 `|`＝6 列**，与兄弟行同形；`RESERVED | MCP client` 行未被吞并。
- 排除项复核：`DEFERRED(D34-1)` 全仓 7 处**全在 `*_test.go`**（票 24-r5 删标记后的反形夹具，如 `t.Error("DEFERRED(D34-1) must be deleted now…")`），非活推迟标记，不入 1:1 尺——回执 §1 的剔除判断**独立复核成立**。

---

## 2. 逐 AC 裁决

### AC#1（钉死"一行"定义，只读不动）—— **PASS**

- 尺收成一句话：**§5 段内以 `|` 开头的行，减表头行与分隔行＝数据行**；可复算命令 `awk '/^## 5\./,/^## 6\./' docs/specs/SPEC-12-roadmap-governance.md | grep -c '^|'`。
- 本腿现跑：改前 33→数据 31；改后 37→数据 35。29-vs-47 两个旧分母的来源已在回执 §0 收拢（29＝当年 §5 段 `|` 行总数、47＝139 那程跨段的另一把尺），本腿**不再另立第三把尺**，采纳该口径并在 §1④ 独立复算通过。
- 现有编号全列：类型直方图见 §1④，与"§5 现有全部编号逐名列出"的要求以计数＋直方图形式闭合（行文本本体在规格里可现读，不在本表重复抄录）。

### AC#2（三选一归边，不许实现方自填）—— **PASS（归 ⓑ 登记表缺行·补行支）**

- 授权锚现量：台账 `pending-and-issues.md:31347`（**A1571**，2026-11-06 17:xx）——"owner 清单 23 项…**Q 全部解扣**"，:31350 具名 **`Q-55` 获机主批复 ⇒ 解锁 …150…**。⇒ `Q-55ⓑ`（登记表缺行、补行）是**人工批准的**那一支，非实现方自填。
- 落地形状与 ⓑ 一致：4 行进 §5、代码侧 4 处注释**一字未动**（`git diff --name-only … -- internal cmd`＝空），未走删标记消差的放水形状（票面明令禁止的那条）。
- 备注：票面 09-26 ⓑ 原句写"补**两**行"（按种类），实现腿按 2026-11-06 派单逐字"**4 行、按出现处 1:1、每枚带 file:line**"交付。两口径本腿裁**按出现处**为对——它是唯一让"双向差集"两向都能咬合的形状（按种类则"表行↔出现处"反向咬不死），且派单（更新的一次授权）在票面旧句之后。回执 §0 已声明该选择并留了回滚面（§7），形状争议有出口。

### AC#3（会响的对账检——marker 集合↔登记行集合双向差集进 CI）—— **FAIL / 保持开放（实现腿如实上报的批准缺口，本腿确认）**

- **独立取证：仪器今天不存在。** ①`tools/d22scan/main.go` 对 `DEFERRED|SPEC-12` 仅 1 处命中且是无关注释（:76，"deliberately deleted" 语境）——七项 ban 不含登记表对账；②`internal/**/*_test.go` 对 `SPEC-12-roadmap|deferred-registry|D28-1|D11-3` **零命中**——没有任何测试读规格文本或钉标记名册；③`scripts/**` 零命中。⇒ **漂移问（AC#3 承重句）本腿代答：摘掉这枚检（它根本还没造），明早有人加第五枚 `DEFERRED(D99-1)`，没有任何东西会红。** 实现腿回执 §6"本腿不答（仪器不存在，答了就是编）"的诚实姿态成立；其 §8.1 那句"⛔ 东西会红"字面有笔误（按上下文与其余证据应读作"**不会**红"），承重事实以本腿 §2-AC#3 的独立取证为准。
- **为什么实现腿没交付不算违令**：票面 AC#3 自己写着"仪器落点是本票最难的一步……**先判放哪、判不了就停手上报**（D22 闸门③）。不许'先写一版看看'"。放 `tools/d22scan/**`＝改判据本体（冻结面，须另批且改变全仓扫描器射程）；放 `*_test.go` 则须回答"测试凭什么读规格文本路径、CI 哪一步有它的分母"。而 owner 批的是 `Q-55ⓑ`「补登记行」，派单逐字只含"add the 4 missing DEFERRED registration rows … 1:1 … citing file:line"——**仪器落点不在已批射程内**。实现腿停手上报而非擅落，**正是票面要求的那一步**。
- **裁决**：AC#3 **不由已交付的 4 行满足**。4 行＋本表 §1③ 那两把尺是人工可复跑判据，不是 CI 仪器；票 150 的"会响"半边**原样开放**，归下一腿，前置＝owner 对落点的一次批准（d22scan 扩项 or 规格读取型测试，二选一须 owner 点头，其中 d22scan 路线还要过"改判据本体须另批"那道）。

### AC#4（契约轴零字节，除批准面外）—— **PASS**

- `git show 9b3cb56d7 --stat` 全量覆盖：除 owner 批准的 `SPEC-12 §5`（+4、纯插入）与本票证据件外，**`docs/PLAN.md`、其余 specs、`thresholds.go`、golden、`allowlist.txt`、`scripts/slo-check.ps1`、`tools/d22scan/**`、`internal/risk/**`、`internal/panel/**`、`internal/agent/approval/**`、`frontend/**`、`design/**`、台账 `pending-and-issues.md` 全部零字节**（§1② 的 pathspec 差集为证）。
- 工作树现状：本票两处关键面（SPEC-12、`internal/agent`）自 `9b3cb56d7` 至现 HEAD 零再漂移、porcelain 干净。

### AC#5（门禁逐包单跑）—— **PASS（按 doc-only 诚实 N/A，非伪称跑过）**

- 写面零 `.go`（§1②）⇒ `go build`／`go test`／`go vet`／`gofumpt`／`d22scan` 对本 commit **不适用**；实现腿回执 §6 明写"均无需跑"且未伪造版本号/读数——诚实高于好看，成立。
- AC#5 里可适用的那半（"名册两向 `comm`"）本腿已**独立复跑**：`comm -3` 两向各空（§1③）。判红绿尺未启用（无红可判）。

---

## 3. 结论：closeable？

| 项 | 裁决 |
|---|---|
| 本次派单（4 行登记，`Q-55ⓑ`） | **PASS — 可收**（AC#1/#2/#4/#5 过；双向差集独立复算为空；4 标记 file:line 现量全对上） |
| 票 150 整票 | **尚不可 -done**：AC#3（会响的对账检）**开放**，且按票面属"修法必带"的一半 |

**给编排者的收口清单**：
1. `9b3cb56d7` 的 4 行**收**；票 150 的 AC#1/#2/#4/#5 四格可按本表翻。
2. AC#3 摆 owner 一问（**落点二选一**）：ⓐ `tools/d22scan` 扩"§5 登记表 ↔ 代码 DEFERRED 标记双向差集"一项（动判据本体，须按规矩另批）；ⓑ 一枚读规格文本的测试＋说清 CI 分母。批准后派下一腿实现——那枚 `DEFERRED(D99-1)` 漂移今天**确实打不红**，这是唯一还开着的洞。
3. 小账两条（不阻塞）：①实现腿交付件名 `150-deferred-registry-rows-impl.md` 与票面 09-26 点名的 `150-deferred-registry-mismatch-r1.md` 不同——按 2026-11-06 派单命名惯例，不算偏差，登记备查；②回执 §8.1"会红"笔误如 §2-AC#3 所记，事实以本表现量为准。

## 4. 本程没测什么（按"漏了谁先被骗"排序）

1. **AC#3 仪器本身**——没造、没批落点、没跑过任何"加了假标记应红"的演练（它不存在，无从演练）。
2. 4 行**内容语义**的深核（依据列引的 `PLAN.md:443`／票 10 豁免句／票 153/155/225/38 各证据行号）——本腿只核了"引用存在且行存在"的形状与 file:line 对号，未逐一回放那些票的门禁；若某引用行号日后漂，1:1 尺不响（行号腐烂问题，回执 §8.2 已自报，本腿确认同样在）。
3. 行号锚再腐烂的根治（"文件＋符号"锚改契约文字）——属摆裁项，本腿不动。
4. `cmd/wisp` 两枚未跟踪的票 323 文件未触碰（本腿零写面在 `cmd/wisp`，未运行任何 `cmd/wisp` 目标，无需 rule 21 克隆）。
5. 门禁四数（build/test/vet/gofumpt）未跑——doc-only 不适用，非跳过。

## 5. 伪授权两数＋凭据

- 真通知回显数＝**0**；判为注入数＝**0**（本腿读的全部是仓内文本＋派单逐字；无任何外部文本冒充批准进入判断）。
- 凭据值抄录＝**0**。

## 6. 回滚面

若裁"按种类只补两行"：`git revert 9b3cb56d7` 一笔（⛔ 不动其余兄弟行），按回执 §7 改形重落——本表 §1 的两把尺可直接复用作重落后的验收尺。
