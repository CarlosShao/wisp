# 台件 265-a1b — 我只补了两格：普查件 §4 撞钉名册 ＋ §5 顺带上报

**腿**：`265-a1b`（只读接续腿，票 265 AC#0 的最后一枚交付物）
**开工锚点** HEAD `c528035f`（2026-10-04 15:34:03 +0800）／**交件锚点** HEAD `89aa72d5`（09:58:26Z＝17:58 +0800）
**唯一写面**：`.scratch/wisp/probes/265/a1/census.md` 的 §4/§5 正文 ＋ 本目录 `.scratch/wisp/probes/265/a1b/**`
**⛔ 未做**：没动票面（一枚 `- [ ]`／`- [x]` 都没碰）、没动 §7、没动前人 §0–§3/§6 任何一行、没改产码、没跑任何 go 命令。

## 0. 交了什么／每格多少行

| 件 | 内容 | 行数（`awk` 区间实量） |
|---|---|---|
| `a1/census.md` §4 | 撞钉名册：族⓪/(a)/(b)/(c)/(d)/(e)/(f) ＋ 汇总，立案 **N1–N21** | **124 行**（`### 4.A` → `## 5.` 之前；含 4.A 锚点节） |
| `a1/census.md` §5 | 顺带上报：(a) Options 字段对表／(b) 两句 record 尺的有无／(c) 7 处跨票冲突／(d) 3 件档外 | **50 行**（`### 5.A` → `## 6.` 之前） |
| 本件 | 自对抗三节＋读数索引 | — |

**完整性两把尺的读数（编排者点名的那两把）**：
- `grep -Ec "待填|填写中" a1/census.md` ＝ **1** ——⚠ **那唯一 1 处在 `a1/census.md:493`，是编排者自己那节标注里引占位尺的原话**（前人/编排者件，⛔ 不许我改）。⇒ §4/§5 内**零枚**占位。
- `grep -c "未判" a1/census.md` ＝ **0**（我自己骨架里的"未判"全部被正文吃掉）。
- `grep -c "本节答" a1/census.md` ＝ **6** ＝ 意图句 3 枚（`:254` §4／`:384` §5／`:479` §7）＋ 编排者标注引用 3 枚（`:491/:494/:495`）。⇒ **§4/§5 的意图句还在（按派单不许删），但它们下面各已有 124／50 行正文**，不是"只剩意图句"那一形。

## 1. 读数档索引（原样命令在每枚文件首行之后）

`logs/00-open-ruler.txt`（起手尺 5 命中＋glob 射程复认）·
`01-family-a-ast.txt`／`01b-counts.txt`（族 a：19 枚 AST 文件／字段名字面量只 2 命中）·
`02-families-bcde.txt`（族 b 42 命中／族 d 反向钉／族 e 4 命中／族 c 面）·
`03-family-c-evidence.txt`（evidence 带行号引用的分布面，含 780 那一发）·
`04-family-c-mech.txt`（evidence 行号有无机器尺：`docs/evidence` 132 命中全注释＋255 名册那两发）·`05-family-c-narrow.txt`（写面两枚文件的带行号引用：docs／.scratch 370）·
`06-session-reply.txt`（`AllowSession` 产码/测试调用者＋读产码文本的尺）·
`07-source-gates.txt`／`08-gate-entrypoints.txt`／`09-ruling-markers.txt`／`10-leg-keys.txt`（N19/N20 那两枚包级台账）·
`11-write-face-anchors.txt`（写面真身行号＋盘脏读数）·
`12-head-vs-inflight.txt`（HEAD vs 在飞 ⓐ-Ⅰ 工作副本）·
`13-sec5-ab.txt`（Options 10 枚字段／GRANT-* 计数）·
`14-sec5a-fields.txt`（4 枚留空字段的兜底与尺）·
`15-queue-and-drift.txt`（NewQueue 默认／在飞件是否碰状态句）·
`16-clock-status.txt`（`gate.go:125-165` 全兜底段＋状态句真身＋256 六臂行号）·
`17-last-checks.txt`（260r4 单源钉／`PanelAllowSession` 零枚／l2_grant_boundary 断言）·
`18-sec5c-conflicts.txt`（四枚票面逐字句＋`GRANT-RECORD-FAILED` 的 HEAD 尺）·
`19-final-status.txt`（污染面终态）。

## 2. 自对抗（一）我写错的尺与读数（逐条，含重跑，⛔ 不删原句）

1. **★一把读数自相矛盾的坏尺（已弃用）**：`logs/01b-counts.txt` 里我写
   `A3 total approval.Options in *_test.go (incl .scratch)`＝**34**，
   紧接着一把"去掉 `.scratch`"的尺＝`git grep -n 'approval\.Options' -- cmd internal tools ':!.scratch' -- '*_test.go'` 读 **36**。
   **36 > 34 不可能**（排除之后反而变多）⇒ 那把尺的形状是错的（`--` 出现第二次之后路径语义变了）。
   ⇒ **处置：这个 36 我不用、也不进任何结论**；族 (a) 的"只有一枚实例"是靠 `01-family-a-ast.txt` §A2（字段名字面量 2 命中）立的，与本条无关。⛔ 我没重跑修正版——修正它要再造一把尺，而我这轮的射程里 A2 那把已经够承重。**归口：下一任若要把"approval.Options 在测试里出现几次"当数用，请重拉。**
2. **同一份日志里锚点行没展开**：`logs/01b-counts.txt` 第一行逐字是
   `# R01b counts only, HEAD=$(git rev-parse --short HEAD)`——**`$( )` 没被求值**（我把它放进了单引号语境）。
   ⇒ 那枚档**缺自己的锚点**，它的读数只能靠"它排在 01 与 02 之间"来定位时刻。**登记为缺陷，不抹。**（其余 20 枚读数档的锚点行都正常展开了。）
3. **差点把"全仓引用面"当名册报出去**：`logs/03-family-c-mech.txt` 第一发我用 4 枚模式合起来打 `docs`＋`docs/reports` ⇒ **780 命中**。
   我最初的写法会是"evidence 钉 780 枚"——**错**：那 780 里绝大多数指的是 `run.go`/`resident_windows.go`/别的票的文件。
   ⇒ 重跑窄尺（`05-family-c-narrow.txt`：只打写面那两枚文件）后才落 §4.E。**原句留在读数档里不删。**
4. **族 (c) 第一发把"零命中"读成"没有行号引用"**：`logs/03` §C7 用 `'resident_approval_windows\.go:[0-9]|resident_windows\.go:[0-9]|run\.go:[0-9]'` 打 255 那两张名册 ⇒ **0 命中**；
   但同档 §C10 换 `'\.go:[0-9]'` ⇒ **17 命中**。⇒ **正解是"255 名册有 17 枚带行号引用，但没有一枚指向 ⓐ-Ⅰ 的写面文件"**（它钉的是 `tiers_255_test.go`／`manager.go`）。
   §4.E-2 那句更正按这个窄写法写；**如果我只有 C7 那一发就下判，我会把"255 没有行号引用"这句假话签进名册**。两发读数都在档里。
5. **§5.B 的枚数第一版被在飞腿污染**：我先跑的是**工作副本**尺 `git grep -n 'GRANT-RECORD-FAILED' -- cmd internal tools` ⇒ 7 命中／测试 **2** 命中；
   第二发换成带锚点的 `git grep -n 'GRANT-RECORD-FAILED' HEAD -- '*_test.go'` ⇒ **恰好 1**。
   ⇒ **那第 2 枚测试命中（`resident_approval_risk_256_windows_test.go:473`）属 `265-r1` 正在改的工作副本，不属 HEAD**；而它还是**错误文本、不是断言**。§5.B 用的是第二发（HEAD 级），**第一发的 2 我不删，就记在这里**。
6. **§4.A 的锚点声明写错一次**：骨架里我写"本节行号一律现量于 `c528035f`"——**第二把尺起 HEAD 就一直在动**（`40aae961`→`9f1e89bf`→`e3574b91`→`2d6aed92`→`5be8821e`→`ba900551`→`862d1736`→`f9dc152c`）。
   ⇒ 已把 4.A 改写成"锚点是散动的、每把尺自带锚点"。**这句骨架原文在本件的 git 历史里（`cc369ce1`），不抹。**
7. **起手尺 glob 的越界风险，我自己加了一发复认**：派单给的 `git grep -n "…" -- '*_test.go'` 是**全仓 glob**（不限派单许可的六个目录）。我没有原样照用其"全仓"含义：补跑 `git ls-files '*_test.go'` 的一级目录分布（`.scratch 59／cmd 62／internal 276／tools 5`）⇒ **本仓 `_test.go` 只活在这四个目录内，全都在许可射程里**；且那 59 枚 `.scratch` 探针件对本尺零命中。⇒ 结论未受影响，但**"我用了派单没限定的 glob"这一格该被记下来**。
8. **`| head -N` 的使用**：`02-families-bcde.txt` 的 §D1、`06`、`07`、`10`、`18` 里有几发用了 `head`（限宽打印）。⇒ **每一处的"枚数"我都在同一次调用里另发了 `wc -l`（或不带 head 的同尺一发）**，名册里落的是 `wc -l` 那个数（族 b 42、族 c 780/370、族 e 4、`GRANT-RECORDED` 14/12）。**没有一处拿截断行数当枚数**——但"用了 head"这件事本身列在这里。

## 3. 自对抗（二）我判不动／量不到的地方（具名＋归口）

⛔ 这一节是本件最该被读的一节：**下面每一格我都没有凭据，只有读码预测。**

1. **"会不会翻"的 12 枚"不翻"判语，全部没有凭据**（§4.B N2–N6、§4.D N7–N10、§4.F N11–N13、§4.G N15–N18）。我没跑过任何一枚用例。
   ⇒ **归口＝U14**（我在 §4.I 新立的）：`go test -count=1 ./cmd/wisp/ ./internal/agent/approval/ ./internal/panel/` 一发即可逐枚证伪。
2. **N19（`cmd/wisp/leg_sink_gate_131_test.go:249` 包级台账）我量不到**：我读到 leg 来自 `main()` 的 argv 分派（枚举器 `:814-817`）、裁定句全仓只在 `cmd/wisp/slo_windows.go:184`，
   ⛔ **但我没能读出"常驻那一枚 leg 今天的 state 是 nailed 还是 ruled"**——那枚门的 ledger 是它自己跑起来才打印的。
   我 `git grep 'registerLegNail131("resident'` ⇒ **0 命中**（`logs/10`），那不等于"resident 没被钉"（leg 的 key 可能不是那个字符串）。
   ⇒ 判语写作"**可能翻＝量不到**"，⛔ 不许被下一任读成"不翻"。**归口＝U13**。
3. **N20（`leg_dispatch_gate_133_test.go`）对"落地腿新写的那枚测试"具体咬在哪一支，我判不了**：我只读到 `:1403`（nail 必须是本目录里一枚 top-level `func TestXxx`）与 `:1651`（不编译的文件里的 case＝no witness）。
   ⇒ 归 ⓐ-Ⅰ 落地腿自己跑。**这一格是 ⓐ-Ⅰ 唯一一处"写测试本身可能被判红"的形状**，编排者的三形代价表里没有它。
4. **`GRANT-RECORD-FAILED` 那格突变（§5.B-3：种"ledger 已接但 Record 恒 error"⇒ 我判 `cmd/wisp` 侧零枚红）我量不到**：这是**突变判语**，只有一发 `go test` 能证。⛔ 我按派单一枚都没跑。**归口＝票 265 AC#3。**
5. **§4.C a-1（"字面量必须恰好 1 枚"）在 ⓐ-Ⅰ 之下到底几枚，是在飞件的事，我不裁**：我只量到 HEAD 是 1 枚、工作副本 `resident_approval_windows.go:368` 仍是 1 枚（`Grants` 进了同一枚字面量）。⇒ 若落地腿为了拿 ledger 又造一枚，那是它的判据面，不是我的。
6. **`.scratch/**` 里那 59 枚 `_test.go` "不在任何分母里"这一句我只读到一半**：我的凭据是 go 工具不进点号目录这一条**通用语义**，⛔ 我**没有跑 `go list ./...` 复认**。⇒ 该句按"读码级"用，别当实测。
7. **前人 §6.2 那 7 格（I1–I7）我一律未复判**（按派单"直接引用不重做"），尤其 **I3（A435 第 1 条是硬约束还是描述）**——票面自己写了"我没裁 A435 原文，所以也不许任何腿拿它当选 ⓐ-Ⅲ 的理由"，我不碰。
8. **U1–U12（前人 §6.1 的没跑尺名册）我一枚都没补跑**（派单硬约束）。⇒ **§4 全部"必翻"判语与它们同一条因果线**：前人 U3/U11 写了"这正是会洗掉 260-v1 读数的动作"，我这轮同样不能跑。
9. **265-r1 的落地形状我读到的是**突变窗口内的工作副本**（17:0x—17:5x +0800 之间它仍是 `M`）。⇒ §4.A／§5.D-1／§5.D-2 三处引用它的地方，**在它自己上盘之后都可能过期**，读时按 `git log -1 -- cmd/wisp/resident_approval_windows.go` 的锚点重核。
10. **我没有裁"这两句假话该由 ⓐ-Ⅰ 的哪一发处理"**（§4.F N14／§5.C-5）：票面 AC#1 ★ 条已把它归给 AC#2 落地，我量到的只是"今天零枚尺"。**归口＝编排者＋落地腿。**

## 4. 自对抗（三）我对盘的污染面自证

终态读数在 `logs/19-final-status.txt`（取数 **2026-10-04 09:58:26Z＝17:58:26 +0800**，锚点 `89aa72d5`）。

`git status --porcelain -- cmd internal tools scripts docs .scratch` 终态：
- **跟踪件脏 17 枚**，⛔ **其中没有一枚是我写的**：
  `.scratch/wisp/issues/263-…md`、`.scratch/wisp/issues/266-…md`（**两枚票面被别人改着，我没碰**）、
  `.scratch/wisp/probes/152/my152.py`、`.scratch/wisp/probes/161/r6/logs/flip-*.txt`（9 枚）、
  `.scratch/wisp/probes/265/r1/evidence.md`（落地腿自己的件）、
  **产码 5 枚＝`cmd/wisp/resident_approval_risk_256_windows_test.go`／`cmd/wisp/resident_approval_windows.go`／`cmd/wisp/resident_task_source_windows.go`（＝在飞的 `265-r1` ⓐ-Ⅰ）＋`internal/config/loader.go`／`internal/config/settings.go`（别的腿）**。
  ⇒ **我全程只 `Read`/`grep`/`git show HEAD:` 它们，一次写都没有**；这也正是"我据此下的每一句判语都带〔仅读码，未跑〕"的原因。
- **未跟踪 564 枚**：全是仓里既有的探针/提交说明堆（前人 §6.3 已具名说过"只建不删"），**我不清、不删、不改归属**。
- **我自己写的东西**＝①`.scratch/wisp/probes/265/a1/census.md` 的 §4/§5 正文（已提交：`cc369ce1` 骨架 → `06b899ad` §4 → `89aa72d5` §5）；②`.scratch/wisp/probes/265/a1b/**`（本件＋`msg-01-skeleton.txt`＋21 枚 `logs/*.txt`，其中 00–17 随 §4 那发提交，18/19 与本件随本发）。
- **git 纪律**：三发提交都带**显式 pathspec 且写在 commit 上**（`git commit -F/-m … -- <paths>`）；⛔ 无 `add -A`/`add .`，无 `--amend`/`reset`/`rebase`/`stash`/`checkout .`/`clean`/`restore`，⛔ **未 push**；仓内**零删除**（新加的只有本目录）。
- **一次失败提交的记录（不藏）**：我第一次执行 `git commit -F <msg> -- <census> <a1b 新件>` 时把**未跟踪文件**写进了 pathspec ⇒ `error: pathspec … did not match any file(s) known to git`、**那发提交整体未发生**（`logs` 之外无一枚被动）。⇒ 修正＝先 `git add -- <我的显式路径>` 再 `git commit … -- <同一路径>`。**这次教训是"新文件必须先 add，pathspec 仍要在 commit 上再写一遍"**。
- **路径长度**：本腿新建的最长路径是 `.scratch/wisp/probes/265/a1b/logs/18-sec5c-conflicts.txt`，**远低于** issues 目录那枚 100／121 的帽子（前人 §6.3 量过帽子读数）。⇒ 这一格**是估的、不是量的**：我**没跑** `scripts/check-path-length-budget.sh`（纯 shell 本可跑，但前人 §6.3 已跑过一发、我不需要重复取数）。⛔ 若编排者要 AC#4 那四门的现量，归 ⓐ-Ⅰ 落地腿。

