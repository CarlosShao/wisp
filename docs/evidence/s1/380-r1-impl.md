# 380-r1 落地腿（写腿）—— **停手上报**：机主批的 A′「删 50 行」与球侧的防删门直接对撞

主题：`票 380 落地腿 380-r1`。锚点 `7ac64ff13`，现跑 2026-10-26 11:35–11:5x +08。
写面：本文件（仓内**唯一**改动）。仓内代码与 `PLAN.md` **零字节改动**，理由见 §1。
临时件（只建不删）：`D:\tmp\wisp\20261026\380r1\{tree-a,tree-b,tree-c,tree-d,tree-e,logs}`，五棵仓外导出树，逐发可复跑。

## 1. ★★★ 停手点（`未定义即停` 命中，不是①②③④里的任何一条字面项）

★**A′ 的实体动作「表缩到 ≤28 行 ⇒ 删那 50 行」不是面板一枚测试的事**：那张表被 **`internal/ball` 五条 C21 测试当契约读**，
其中两条是**防删门**，而 `internal/ball/tokens_table_test.go` **不在我的写面**（写面只给了 `internal/ball/testdata/**`）。
⇒ ★**实测（不是推理，§8）**：把表删到 28 行 ⇒ `internal/ball` **rc 1／30 条 error**；删到 5 行 ⇒ **39 条**。
⇒ ★**两条互斥的路，我无权选**：要么**越写面**改球侧那两道门，要么**让球侧由绿变红**。裁形是 `AC#1`（编排者/机主）的事，⛔ 我不自己换形。
⇒ 附带的第二个理由：机主批的是「**缩契约到 ≤28 行**」。**删掉 26 个仍在绘制路径上的 `Palette` 字段的唯一书面记录**是否仍在授权范围内，属契约面（`SPEC-12 §4.1`），须您一句明示。

## 2. ★★★ 第 0 步答案：那 78 行写在哪

★★**表＝`docs/evidence/s1/c21-native-tokens.md`**（**298 行**，`git` 跟踪）。逐段：
- `:46` `## 明色系（:root，默认）` ＋ 行 **`:50-89`＝40 行**；`:91` `## 亮色系（[data-theme="light"]…）` ＋ 行 **`:95-132`＝38 行**。40＋38＝**78**。
- ★**每行恰好 4 格**：`| CSS 变量 | tokens.css 值 | Go 字段 | Go 值 |`。这个「恰好 4」在 §8 是硬约束。
- ★**不是 `design/assets/tokens.css` 本体**（该文件 `72b65c614` 已删；盘上 `design/old/…` 是 `??` 未跟踪）；★**也不是 `C1–C32` 文本**（`PLAN.md` 全 3,611 行 **0 处**命中 `c21-native`）。**契约文字在 `PLAN.md:1377`／`:3308`／`:3604`，表是那张契约的书面镜像。**

## 3. 四条核（`c21-a1` 的结论我逐条自己跑，⛔ 采信）

| # | 它的结论 | 我的尺 | 结果 |
|---|---|---|---|
| ① | P4 `:root` 43 条、`[data-theme="dark"]` 43 条、`light` 0 条 | 对 `frontend/src/styles/tokens.generated.css` 按声明正则实数 | ✅ **43／43／0**；＋5 条派生后 = 测试眼里的 **48** |
| ② | `PLAN.md` 里 `c21-native` 0 命中 | `Select-String -SimpleMatch` | ✅ **0**（文件 3,611 行） |
| ③ | `checked < 70` 在 `:546` | `Select-String 'checked < 70'` | ✅ **`:546`**（fatal 文本在 `:547`） |
| ④ | S2 那发确实 `0/78` | 仓外树里把 P1 指到夹具 ＋ 28 行改写成第四代，跑**真测试** | ✅ **`:547 only 0 of 78 colour rows … not doing its job`**（§6） |

⇒ ★**四条全对得上 ⇒ 不触发停手②。**（另外白捡两条对账：我 S1 的 `:477` 恰好 **39** 条、`126` 条 P1 独有 —— 与它的「39／126」逐数相同，见 §6。）

## 4. ★★ 改前四数（整包，⛔ 相减）＋ 格式甲乙两尺

| 整包 | 改前 rc | 明细 |
|---|---|---|
| `go test -count=1 ./internal/panel` | **1** | **5 个 FAIL**：`TestC21DesignTokensFourWayAgree`（本票）＋**4 个既存红**（`TestApprovalCardViewJSONKeysMatchFrontendTypes`：`Go Snapshot emits [instructions tasks] that interface PanelSnapshot does not declare`；`TestComposerContractTypesMatchFrontend`；`TestStreamLogFloodBelowKeyBoundIsNotBounded35r8`；`TestStreamLogDroppedNamingLedgerIsNotBounded35r8`）★**这 4 条与 C21 无关、但我没动它们** |
| `go test -count=1 ./internal/ball` | **0** | —— |
| `go test -count=1 -run TestC21 -v ./internal/ball` | **0** | 5 条**全 PASS**（`ColourRowsMatchCode`／`GeometryRowsMatchCodeConstants`／`TokenConsumerReport`／`CodeTokensAreTabledOrExempt`／`ColourRowsMatchTokensCSS`） |
| `go test -count=1 -run TestC21DesignTokensFourWayAgree ./internal/panel` | **1** | **`:441` 读不到文件** ⇒ 一条断言都没跑（S0） |

改后：**N/A**（仓内零改动）。★格式甲乙两尺（工作树零改动 ⇒ 两把都空，仍报）：`gofmt -l internal/panel` ＝ **空**；`gofumpt -l internal/panel`（绝对路径 `/d/work/base/gopath/bin/gofumpt.exe`）＝ **空**；`go vet ./internal/panel` **rc 0**。

## 5. ★★★ 六道 fatal 门逐条（改前：**六道里五道够不到**）

| # | 站点 | 判据 | 改前状态 |
|---|---|---|---|
| g1 | `:139` | 解析出 `dark` 块才算数（`no :root block parsed`） | ★**不可达** |
| g2 | `:211` | 表行数 `< 70` ⇒ fatal（`the table shape changed and this check would be vacuous`） | ★**不可达**（§8 实测：**5 行**的表仍死在 `:441`） |
| g3 | `:284` | 两个 palette 都要解析出字段 | ★**不可达** |
| g4 | `:441` | `readRepoFile(design/assets/tokens.css)` | ★★**唯一可达、且第一个响** ⇒ 后面全不执行 |
| g5 | `:456` | 某主题 `len(want)==0` ⇒ `this check would be vacuous` | ★**不可达** |
| g6 | `:546` | `checked < 70` | ★**不可达** |

⇒ ★★**这张网今天零牙**：不是"红"，是**一条都没跑到**。「六道门」是设计里的六道，实测**六道全哑**（g4 替其他五道挡了枪）。

## 6. ★★ S1／S2 复现（真测试，不是移植）

| 发 | 树的状态 | 读数 |
|---|---|---|
| **S1** | P1 指到 `internal/ball/testdata/tokens_c21.css`（14,392 B，`c21-r1` 那枚） | rc 1；`:446 parties:` ⇒ P1 **135** dark decls、P4 **48**；`:461`×**197**（＝P1 独有 **126** ＋ light 腿 **71**）、`:477`×**39**（P4 独有）、`:471`×9、`:481`×2、`:481`；**`:547 only 0 of 78`** |
| **S2** | S1 ＋ 28 行值改写成第四代 | rc 1；`:501`×**28**（leg2 由 0 红变 28 红）、`:532` 73→**50**、`:461`×197、`:477`×39；**`:547 only 0 of 78`** |

⇒ ★**S1、S2 都是 `0/78`** ⇒ 「只改表」买不到绿，与 `c21-a1` 同数（197＝126＋71 这个分解是我这次才算出来的，它对上了它的 126／39）。

## 7. ★★★ 极性的牙（A/B 两发，这是「修正真在起作用」的证明）

两发都让 **P1＝P4**（旧 `tokens.css` 已删，这是唯一能读到的真相源），只差极性：

| 发 | 极性 | 读数 |
|---|---|---|
| **d** | **修好**（`:root`→light、`[data-theme="dark"]`→dark） | rc 1；`:483 leg 1 theme dark: 48 tokens reconciled` ＋ **`leg 1 theme light: 48 tokens reconciled`**；★**78 行真的跑到了**（`:494`×68 ＋ `:503`×10）；`:549 only 0 of 78` |
| **e** | **原样**（`:root`→dark、只认 `[data-theme="light"]`） | rc 1；**`:458 theme light: … parsed no declarations, this check would be vacuous` ⇒ 当场 fatal**，★**78 行一条都没跑到** |

⇒ ★★**牙在哪**：极性装反的代价不是"数字难看"，是 **light 腿空跑、g5 直接 fatal、整张表 78 行全部不执行**；修好后**两个主题各自 48 条被真量到**。
⇒ ★**附带发现（这条对裁形很硬）**：极性修好后 **g1（`:139 len(out["dark"])==0`）反过来绊住旧夹具** —— 旧 `tokens.css` 只有 `:root` ＋ `[data-theme="light"]`，**没有 `[data-theme="dark"]`** ⇒ 实测 `:443 no :root block parsed`。
⇒⇒ ★★★**所以「极性修正」与「P1＝旧夹具」不可共存**：极性一修，P1 只能指 P4，而 **P1＝P4 时 leg1 的 P1→P4 方向变成 P4↔P4 恒等** ⇒ **leg1 那条腿在 A′ 里必然死掉**，底座必须换（`c21-a1` 说的"表和四条腿底座一起换"就是这个）。

## 8. ★★★ 删表实验：28 行 / 5 行（对撞的正身）

表照第四代真值删到 **28 行＝14 名×2 主题**（留：`--bg-inset`／`--bg-overlay`／`--border-soft`／`--fg-primary`／`--fg-secondary`／`--fg-tertiary`／`--accent`／`--accent-line`／`--danger`／`--danger-line`／`--warn`／`--warn-line`／`--success`／`--success-line`），跑 `go test -count=1 -run TestC21 ./internal/ball`：

| 表 | 球侧 | 逐条 |
|---|---|---|
| **28 行** | **rc 1／30 条** | `FAIL TestC21TableColourRowsMatchCode`；**`:678`×2**（`dark 14 rows, want at least 40` ／ `light 14 rows, want at least 38`）；**`:761`×26**（`Palette.X exists in tokens.go but the dark (:root) table declares no row for it`）；`:805`／`:780` 各 1 |
| **5 行** | **rc 1／39 条** | `:678`×2（其中 `light 0 rows, want at least 38`）；**`:761`×35**；余同上 |

★**同一张 5 行的表拿来跑面板**：仍然 **`:441` 读不到文件** ⇒ ★**面板侧自己的防删门 g2（`:211`）永远够不到**，删表在面板侧**今天连红都不会红**（这正是不许只降门槛的理由）。
★★**硬约束补一条**：球侧解析器是 `if len(cells) != 4 { continue }`（`tokens_table_test.go:180`）⇒ ★**给表加第 5 列会让球侧把整张颜色表解析成 0 行**，g2 同类 fatal 立刻响。⇒ ★**「加一列记第四代槽位」这条路在写面内也是死的。**

## 9. ★★ 逐行清单：留 28／删 50／丢 26 个字段的记录

- **留 14 名（×2＝28 行）**：见 §8 名单。理由：第四代有同名或一对一的槽位（能 12 ＋ 换算 16）。
- **删 50 行（26 暗＋24 亮）**：玻璃 `--glass-ring/-hi/-hi-soft` **6** ＋ `--info`／`--info-line` **4** ＋ `--warm`／`--warm-line` **4** ＋ `--orb-*` **10** ＋ `--ball-halo`＋`--on-solid` **2** ＋ `--tint-*`（12 名）**24**。理由：第四代**整族没有这个概念**（`gen-tokens.mjs:4-7` 极简风；`--tint-*` 是表自己记的零消费者死色）。
- ★★**代价（这是停手的正身）**：删除让 **26 个 `Palette` 字段**（`GlassRing GlassHi GlassHiSoft Info InfoLine Warm WarmLine OrbHi OrbBody OrbBody2 OrbRim OrbShadow BallHalo OnSolid TintAccentHi TintSuccessHi TintSuccessLo TintWarmHi TintWarmLo TintWarnHi TintWarnLo TintDangerHi TintDangerLo TintNeutralHi TintNeutralMid TintNeutralLo`）**失去唯一一行书面记录** —— 它们**仍在绘制路径上**（球体玻璃/球体/液面），不是死代码。

## 10. ★★ 替代判据设计（若裁形换，这是我准备落的形；⛔ 不是"把 70 改成 22"）

1. **集合相等（正向）**：表引用的每个名，必须在 P4 的 `:root` **与** `[data-theme="dark"]` **两块都声明**；少一块即红。
2. **集合相等（反向，allow-list）**：P4 的 **34 个带色 token** 必须落在「表」∪「豁免表」里；每条豁免写一句理由（指向消费方或票号）；**豁免集合必须等于代码里一枚钉死的字面量清单**。
3. **双主题各自被检查**：两块的声明名集合必须相等（生成件结构完整性）；每块 reconciled 各自 `> 0`；且 `行数 == 2 × 名字数`、每个名字在两个主题各占一行。
4. **下界＝22，理由**：`--accent-line`／`--danger-line`／`--success-line` 三枚 **P4 根本没有对应声明**（P4 只有 `--warn-line`），3 名×2＝6 行**无对象可校** ⇒ 实测天花板 **22**（不是 28）。下界取 22 ⇒ **谁再删一行就红**，且它不许无契约变更地往下调。
- ★**为什么它不许假绿**：①删行 ⇒ 判据 2 红（那 34 个里就有被删的名字）＋判据 4 红；②删行并顺手补豁免 ⇒ 判据 2 的**钉死清单**红（豁免膨胀要改代码，看得见）；③改成一堆只有 dark 行的表 ⇒ 判据 3 红；④把表删到只剩 5 行 ⇒ 判据 2／3／4 **同时**红。★**四类形状变异没有一类能靠"数字变小"混过去。**

## 11. ★★ 「⛔ 假绿」自证（设计＋今天为什么还没跑）

★**两发正交变异**，各打一类：①**形状变异**（删到 5 行）⇒ 由 §10 的 2/3/4 抓；②**值变异**（把表里 `--inset` 暗值 `#1f2022` 手改成 `#1f2023`，名字集合一个不动）⇒ 判据 1/3 全绿，**必须**由「表值 == P4 声明值（`valuesAgree`）」那条抓。
⇒ ★**结论**：行数少不等于抓不到漂移 —— 抓漂移的不是"数了多少行"，是「**名字集合相等（正反两向）** ＋ **逐个值相等** ＋ **两主题各占一行**」这三件；行数只是其中一件的必要条件，不是充分条件。
★**今天跑不了**：这一发要落在重写后的测试上，而重写与否取决于 §12 的裁形 ⇒ 我把它连同发 ② 一起留给下一腿，**不在这里假装跑过**。

## 12. ★★ 裁形三选一（⛔ 我不自己换形）

| 形 | 做什么 | 写面 | 球侧 | 价 |
|---|---|---|---|---|
| **A″ 分区** | 78 行不删，切成「第四代共有（14 名）」与「仅原生世代（50 行）」两节；**子标题仍含「明色系／亮色系」**（球侧 `Contains` 照旧数到 40／38）；面板只对前一节做四方对账 | 表 ＋ `tokens_fourway_test.go`（**都在我写面内**） | **零影响** | 表**行数不变**，缩的只是**契约覆盖面**（78→22） |
| **B″ 真删** | 照原批删 50 行 | ★**必须**再改 `internal/ball/tokens_table_test.go:673-674`（40/38 下限）与 `:760-761`（code→table）⇒ **越我写面** | **变红** | 26 个活 token 从此无书面记录 |
| **C″ 加对照表（推荐）** | 78 行原样；表里新增两张**2 列**小表：`第四代槽位对照`（11 名）＋`第四代未覆盖`（23 名＋一句理由）。球侧 `len(cells)!=4` 直接跳过 ⇒ 不受影响；面板读这两张 | 表 ＋ `tokens_fourway_test.go` | **零影响** | 多两张小表，但**映射是数据、可审计**，不是藏在 Go 里 |

★**我推荐 C″**：它是 A″ 的严格超集（分区要动 78 行的归属，C″ 一行不动），是主流 harness 的 allow-list 形，且**球侧零影响**。★**B″ 只有您（或机主）能批**，因为它的代价落在写面外＋剪掉球侧那张网。
★**已就绪的读数**：形定了之后，§7 的极性改法（`blockRe` 加 `[data-theme="dark"]`、`:root`→light、g1 改成两块各判非空）与 §10 的四条判据都是现成的，一发可落。

## 13. 复跑（逐发）

`D:\tmp\wisp\20261026\380r1\` 下：`tree-a`＝S1/S2；`tree-b`＝28 行／5 行；`tree-c`＝极性修好＋旧夹具（撞 g1）；`tree-d`＝极性修好＋P1=P4；`tree-e`＝极性原样＋P1=P4；日志 `logs\{base-*,s1-*,s2-*,c-*,tree-d,tree-e,e2-*,e2b-*}.txt`。
★环境：`GOFLAGS=`、只 `unset GOTMPDIR`；rc 一律 `cmd > out 2>&1; rc=$?`（⛔ 走管道）；**工作树零突变**（全部变异在仓外导出树里）。
