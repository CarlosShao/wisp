# 371-a1 普查腿（只读）—— `R20`／`M4` 是「契约」还是「机主定案」＋那一发读数＋`F4` 的 panic

- 腿 `371-a1`（**只读**，⛔ 实现者）。HEAD **`127784149`**，量于 2026-10-24 09:3x–09:5x +08。
- 写面：本文件 ＋ 仓外 `D:\tmp\wisp\20261024\371a1\**`（3 支探针／3 份 overlay／`leg1..4.out`／`build|vet|d22scan|pathlen`.out）。
  **仓内 `.go` 改动 0 字节**；`git status --porcelain internal` **空**（现跑）。
- ★**本件所有行号一律取自 HEAD blob**（`git show HEAD:<path>`）。理由＝⑥.3：工作树此刻**有别人在飞**（票 114 AC#3 正在改 `cmd/wisp/panel_inbound.go`／`cmd/wisp/run.go`）。

## ⓪ 判（六句，逐条落点见后）

1. `M4`／`R20` 在 `docs/PLAN.md` ＋ **14** 份 `docs/specs/*.md` 里 **0 处**（`grep -n` 0 行、rc=1）⇒ **⛔ `D` 级契约文字**
   （`SPEC-12 §4.1` 枚举的受保护体是 `C1–C32`／`D1–D47`；`R1–R9` 是另一族）。
2. 但它们**具名在台账**，且逐字是**机主本人的定案**：`docs/reports/pending-and-issues.md:3019`（「裁定 R20（…**owner 对 M1–M5 的定案**）」）
   ＋ 同文件 `:3031` 的 `| **M4** |` 行；落点＝**票 90 AC#3 已勾**。
3. **判语**：「扩 `M4` 让所有变宽出卡」**⛔ 改契约文字（字面）**、**是反机主原话（实质）** ⇒ 我判 **须摆机主**（＝铁律三类① 的实质）；
   只按字面（受保护体只算 `C1–C32`／`D1–D47`）才轮到编排者拍。★**这一步未定即停、我不自行落**：派单给的二分框（`plan/specs` 有／只有代码）**⛔ 覆盖本件实况**（第三处：台账）。
4. ★**一枚词面漂移（新发现）**：代码拒绝句的括注写「R20/M4：**只有变宽要确认**」，机主原话是「**只有全自动要确认**」；
   前者在 `PLAN.md`／`specs`／台账 **0 命中**（只在代码 1 处 ＋ 证据件若干）⇒ 那枚括注把**方向判据**说成了**产品规则**。
5. 那一发我自己量（真 `config.toml`、真 `Confirm`）：`AE→AH` ⇒ **`confirm_calls=0`、`applied`**；正控 `AE→AA` ⇒ **`1`**。⇒ `367-r1` 的读数**成立**。
6. `F4` **我自己造出来了**：谓词一翻 ⇒ **`store.go:224` panic**（`rc=1`、`signal 0xc0000005`）。归属＝**票 114（`12d8e28c7`）**、⛔ 367 造成 ⇒ 票面的归属更正**成立**。

## ① `R20`／`M4` 的逐字出处（现跑）

### ①.1 契约文档：**0 处**（尺＋分母都核过）

| 尺（UTF-8 脚本文件跑，⛔ PS 传参） | 读数 |
|---|---|
| `grep -n 'M4' docs/PLAN.md docs/specs/*.md` | **0 行**，rc=1 |
| `grep -n 'R20' docs/PLAN.md docs/specs/*.md` | **0 行**，rc=1 |
| `ls docs/specs/*.md \| wc -l`（分母） | **14** |
| `grep -ohE '(^\|[^A-Za-z0-9_])M[0-9]+' … \| sort \| uniq -c` | **只有 `M0` 15 处**，其它 M# **0** |
| 同式数 R-token | `R1 15 · R2 2 · R3 2 · R4 9 · R5 2 · R6 4 · R7 7 · R8 10 · R9 11` ⇒ **最大 `R9`，无 `R20`** |
| 语义反查（怕只是没写名字）：`全自动`／`只有全自动`／`只有切到`／`每步都问`／`只问高危`／`auto_approve`／`ask_every_step` | `PLAN.md`＋specs **全 0 命中** |

⇒ ★三档权限模式这套东西在 `PLAN.md`／`specs` 里**连语义都不存在**（`PLAN.md:122` 的 D4 讲的是 L0/L1/L2 **风险级**，另一根轴；
`PLAN.md:2138` 唯一一处「权限模式」是在讲 **Claude Code** 的借鉴表）。

### ①.2 台账：具名，且逐字是**机主定案**（`docs/reports/pending-and-issues.md`）

- `:3019`（节标题，逐字）：`## 裁定 R20（2026-09-21 17:1x，**owner 对 M1–M5 的定案**）：面向用户的权限模式 = **三档 / 默认最严 / 档位持久化 / 只有全自动要确认 / 显示在主面板输入框**`
- `:3031`（M4 那一行，逐字）：`| **M4** | **只有**切到"全自动"要确认（一次 L2 强确认）；**但审计日志三档切换全部写** | 一致（我把"审计"当默认，他把"确认"收窄） | 票 90 AC#3 |`
- `:2312`（A63①，逐字半句）：`只有切到"全自动"要 L2 确认但**审计三档都写**`
- `:1271`（A84⑥，逐字）：`R20（三档 / 默认最严 / 档位持久化 / 只有全自动要 L2 / 显示在主面板输入框 / 不换 git 分支）`
- **落点**（`.scratch/wisp/issues/90-user-facing-permission-modes-done.md:55`，`[x]` 已勾，逐字）：
  `- [x] **AC#3** 切到"全自动"档要 L2 强确认 + 写审计（**审计三档都写**，只有全自动那一次要确认）。`

### ①.3 判语（按铁律三类①，逐字给理由）

- **字面**：`M4`／`R20` 不在 `PLAN.md`／`specs` ⇒ **⛔ `D` 级契约文字** ⇒ 按派单二分框会落到「实现决定、编排者拍」。
- **实质**：`:3019` 的标题自己写着「**owner 对 M1–M5 的定案**」、`:3031` 的 M4 格是「**他**的定案」⇒ 扩 `M4` ＝ **改机主拍下的产品决定**，
  且该决定**已落成已勾的 `AC#3`** ⇒ 属铁律三类①「**改契约文字本身**」的实质面 ⇒ **摆机主**。
- ⇒ **我判「须摆机主」**；★这也是派单「未定义即停 ①」的邻格：**派单只叫我查 `plan`／`specs` 两处，实况在第三处（台账）**，我不自行把「台账里的机主定案」降级成「实现细节」。**给编排者的一句话问法**（供抄）：「面板里把档位从『每步都问』切到『只问高危』，要不要弹一张卡？（今天不弹；弹＝改您 9-21 的 M4）」

## ② `run_mode101_test.go:251` 那句断言引的出处

- `:232`（注释，逐字）：`// R20/M4 is asymmetric: only the switch INTO auto_approve costs a card.`
- `:251`（断言，逐字）：`t.Errorf("L2 confirmations asked = %v, want exactly [auto_approve] (R20/M4)", gotAsked)`
- 同文件第二枚引用 `:558`（逐字）：`t.Errorf("cards shown = %d, want exactly one L2 card per switch (R20/M4)", cards)`
- **它引的 `R20`／`M4` 在契约文档里有没有出处 ⇒ 没有**（①② 已量：`PLAN.md`＋14 份 specs **0 处**）。
  名字的唯一落脚处＝**台账 `:3019`／`:3031`** ＋ 全仓代码注释（`internal/perm/store.go:17/175/216`、`internal/panel/composer_handlers.go:27/82/86/132/137`、
  `cmd/wisp/run.go:121/871/896`、`cmd/wisp/run_mode101_test.go:232/251/558` 等）。
- ⇒ **判**：这条断言钉的**⛔ 契约**，钉的是**机主 9-21 的定案**（①.2）。它**照样是硬的**：改它＝改产品行为，且 `cmd/wisp/**` 是票 371 的写面外。

## ③ 那一发：`AE→AH` ＋ 真 `Confirm`（我自己量，⛔ 抄 `367-v1`）

尺＝3 支自建探针，`go test -overlay=` 注入（探针 `.go` 全在仓外，模板/夹具用包内既有 `writeConfig`／`openManager`，**真 `config.toml`**）：

| # | 形状 | 读数（逐字取自 `leg1.out`／`leg2.out`） | rc |
|---|---|---|---|
| 1 | store 层 `AE→AH` ＋ **真** `Confirm` | `confirm_calls=0` `asked_for=[]` `set_err=<nil>` `mode_after=ask_high_risk` `last_result=applied` | 0 |
| 2 | **正控** store 层 `AE→AA` ＋ 同一真 `Confirm` | `confirm_calls=1` `asked_for=[auto_approve]` `set_err=<nil>` `mode_after=auto_approve` | 0 |
| 3 | 组合层（真 `ModeWriteHandler` ＋ 真 `Store`，**两层都接腿**）`AE→AH` | `handler_confirm_calls=0` `store_confirm_calls=0` `handle_err=<nil>` `mode_after=ask_high_risk` `last_result=applied` | 0 |
| 4 | 组合层正控 `AE→AA` | `handler_confirm_calls=0` `store_confirm_calls=1` `handle_err=<nil>` `mode_after=auto_approve` | 0 |
| 5 | **无腿**两发（`AE→AH`／`AE→AA`） | 两发都 `refused-confirm`、`mode_after=ask_every_step`，句＝`perm: 把档位从 … 放宽到 … 需要一条确认通道，但本机没有接入（fail-closed 未切）` | 0 |
| 6 | 组合层**产线装配**（两层都 nil）`AE→AH` | `handle_err`＝`panel: 确认腿未接入，档位不得变宽: 请求 "a1-none-ae-ah" …（R20/M4：只有变宽要确认）`、`last_result=no-change` | 0 |

⇒ ★**尺两侧都读得出**（0 与 1 都能读）⇒ 读数 1 不是「仪器读不到」的假零。★**「问没问」的判据在 `Set` 里没有**：
`store.go:216-221` 的注释逐字 `// R20/M4's asymmetry stands on top of that gate: of the widenings, only` / `// ... There is no nil test here on purpose - ...`
⇒ 内层**只有** `:222 if to == risk.ModeAutoApprove {` ＋ `:224 if err := s.confirm(ctx, sw); err != nil {`，**没有任何"问过没问过"的判据**，
门（`:209 if IsWidening(from, to) && s.confirm == nil {`）判的确实是**通道的 nil-ness**。

## ④ `F4` 的 panic 实证（我自己造，⛔ 读码推断）

- **变异**：`-overlay` 把 `internal/perm/direction.go` 换成 `func IsWidening(from, to risk.Mode) bool { return to < from }`（仓内 0 字节动）
  ⇒ `AE→AA` 不再算变宽。**同一支探针**在两种谓词下各跑一发（单变量 A/B）：
- **对照**（HEAD 谓词）：`A1|F4|predicate_says_widening=true (HEAD predicate)` ⇒ `survived`、`last_result=refused-confirm` ⇒ **不 panic**（`leg3.out`，rc=0）。
- **变异**（翻后谓词）：`A1|F4|predicate_says_widening=false (flipped overlay is in force)` ⇒ 随后**崩**（`leg4.out`，**rc=1**），逐字：

```
panic: runtime error: invalid memory address or nil pointer dereference [recovered, repanicked]
[signal 0xc0000005 code=0x0 addr=0x0 pc=0x7ff75f7ad9b3]
github.com/CarlosShao/wisp/internal/perm.(*Store).Set(...)
	D:/work/workspace/projects plans/Wisp/internal/perm/store.go:224 +0x333
```

⇒ **判**：`F4` **真会 panic**，落点逐字＝`store.go:224`（`:222` 的 auto 分支调 `s.confirm`，而 `Confirm: nil`）。
⇒ 这决定「**要不要补回内层 `nil` 兜底**」：`:216-221` 那枚注释今天**靠"谓词永不错"承重**（`There is no nil test here on purpose`）⇒ 兜底是**免费的**，但属**实现决定**、⛔ 契约。

## ⑤ 归属：`composer_handlers.go` 那道门什么时候进来的

- **引入提交**：**`12d8e28c783f529f762c4cb8d9130f55d51a2f5d`**（2026-09-23 10:47:20 +0800）
  `feat(114,AC#2): 面板档位写入的处理器在没有 L2 确认腿时拒，且拒在调用写入腿之前`（`+514/−3`，新件 `composer_handlers.go` 165 行）。
- **当时的门（该提交 blob 逐字）**：`:107 func modeIsWidening(from, to risk.Mode) bool { return to > from }` ＋
  `:133 if modeIsWidening(from, to) && h.Confirm == nil {` ＋ `:141 if err := h.Modes.Set(ctx, to, PanelModeOrigin, h.actor()); err != nil {`
- ⇒ ★**票面归属更正成立且更强**：`internal/panel` **自票 114 起**就①**只读 `nil-ness`**、②**已经按"方向"判**（私有谓词，与 `perm` 各一份）。
  ⇒ 今天的形状⛔ 367 造成：367 只做了两件事 —— 删掉 panel 的私有谓词、换成 `perm.IsWidening`（`git show 7fabccf1b` 的 hunk 逐字
  `-func modeIsWidening(...)` / `+if perm.IsWidening(from, to) && h.Confirm == nil {`），并把**同一判据**搬进 `store.Set:209`。
- 旁证（367 自己改写掉的旧件头，逐字 `-` 行）：`// WHY A GATE HERE WHEN perm.Store ALREADY HAS ONE. internal/perm's Set refuses` /
  `// the switch INTO auto_approve when it was handed no confirmation leg` / `// (internal/perm/store.go's \`if s.confirm == nil\`). That covers one档 of three`
  ⇒ **pre-367 的两层是两种判据**（store 判目标、panel 判方向），但**panel 那层判方向＋只读 nil-ness 从 114 就是这么写的**。

## ⑥ 票面前提新鲜度（对 HEAD blob 复量，含行号）

| 断言（票面／`367-v1`） | HEAD `127784149` 读数 |
|---|---|
| `func (s *Store) Set` `store.go:189` | **189** ✔ 未动 |
| 新门 `if IsWidening(from, to) && s.confirm == nil` `store.go:209` | **209** ✔ |
| auto 分支 `:222` ＋ 调 `s.confirm` `:224`（`F4` 落点） | **222 / 224** ✔（④ 实证） |
| `composer_handlers.go` 门 `:131`、`h.Modes.Set` `:141`、括注 `:132` | **131 / 141 / 132** ✔（`367-v1` 的 131/141 仍准） |
| `run_mode101_test.go` 注释 `:232`、断言 `:251`、第二枚 `:558` | **232 / 251 / 558** ✔ 未动 |
| 票面 `:15` 的关联行 `cmd/wisp/run.go:687` | **687**＝`rt.modeWrites = &panel.ModeWriteHandler{` ✔（`:663 confirm := s.modeConfirm`、`:667 perm.New`） |
| 产线两条入口的腿 | `panel_inbound.go`（HEAD）：`:240` 4 参签名、`:254 Confirm: nil`（store）、`:264 Confirm: nil`（handler）、`:260` 本地 `modeWrites`、`:296` 塞进 dispatch ⇒ **唯一可达的档位处理器＝两层都 nil** |
| ★**「有通道的装配」今天可达吗** | **⛔ 不可达**：`git grep -n modeWrites HEAD -- 'cmd/**/*.go'` ⇒ 仅 5 处（`panel_inbound.go:260/296`、`run.go:331/340/687`）⇒ **`rt.modeWrites` 全仓无读者、是死字段** ⇒ 票面的形状今天**潜伏**、一接线就活 |

### ⑥.3 ★工作树有别人在飞（我的行号因此一律取 HEAD blob）

- `cmd/wisp/panel_inbound.go`：`git status` **` M`**、mtime **09:42:31**、`+53/−10`；内容逐字是**票 114 AC#3**（把 `panel.ModeConfirm` 作参数穿进
  `newComposerDispatchChain`，给**常驻面板**递真腿，`storeConfirmFromPanel()` 收窄成 `perm.ConfirmFunc`）。
- 09:47 再看：`cmd/wisp/run.go` **也脏了**（同一笔在飞）⇒ ★**一旦落成，常驻面板那条链就是 `arm=both`**，票 371 的潜伏形**当天变活**。
  ⛔ 我的：这两枚是**别人的活**，我没碰、没清、没 checkout（本件写面只有 1 枚 docs 文件）。

## ⑦ 门禁 rc（⛔ 走管道；四数分开；**其中一枚红不是我的**）

| # | 命令 | rc | 读数 |
|---|---|---|---|
| 1 | `GOFLAGS= go build ./internal/...` | **0** | 静默 |
| 2 | `go vet ./internal/perm ./internal/panel` | **0** | 静默 |
| 3 | `git status --porcelain internal` | 0 | **空**（本腿 0 字节） |
| 4 | `sh scripts/d22scan.sh` | **0** | `clean - no D22 ban violations`（ban #8：internal 639、cmd 149） |
| 5 | `sh scripts/check-path-length-budget.sh` | **0** | `VERDICT GREEN`（over-budget 61＝roster 61、61=61 双向同集） |
| 6 | `GOFLAGS= go build ./...` | **1** | ★**⛔ 我的、⛔ HEAD 的属性**：`cmd\wisp\panel_resident_windows.go:228:77: not enough arguments in call to newComposerDispatchChain`／`want (…, panel.ModeConfirm)` ⇒ ⑥.3 那条在飞腿**未改完调用点**（`resident_windows.go:214` 同理，`go vet ./cmd/wisp` 同因 rc=1） |
| 7 | 探针四发 | **0/0/0/1** | ③ 两发 rc=0、④ 对照 rc=0、④ **变异 rc=1（panic）** |

- **格式门禁**：本笔 **0 枚 `.go`** ⇒ 甲／乙两尺**都不适用**；⛔ 用 `gofumpt -l` 反推换行。**包测试**：本腿只读、未改码 ⇒ **未跑**
  `./internal/perm`／`./internal/panel` 整包（`367-v1 ⑨` 的 19/166 是它那笔的读数，我不冒名复述）；`cmd/wisp` 整包本就会死在 `MARK299`，我一条都没跑。

## ⑧ 提交 ＋ 逐枚 `.go`（应为 0）

- 本件入库首笔：**`fd3c6ae09`**（`+158`、1 文件）；本行之后对本件的续笔见 `git log --oneline -- docs/evidence/s1/371-a1-census.md`。逐枚文件：**只有本件 1 个路径**；逐枚 `.go`：**0**。越界核对：⛔ 改
  `internal/**`／`cmd/**`／`scripts/**`／`.github/**`／`.scratch/wisp/**`／`docs/reports/**`／`docs/PLAN.md`／`docs/specs/**`；⛔ 落 `.go`（探针全在 `D:\tmp\wisp\20261024\371a1\**`）；⛔ 翻框／改票名／`Status:`／`-done`／台账；⛔ push。

## 附：落笔前五问（尺的射程）

1. **量的是哪个对象**：①`PLAN.md`／`specs` 的**文面**（⛔ 代码）；②`*perm.Store.Set` 与 `*panel.ModeWriteHandler` 的**运行行为**（真 `config.toml`）；
   ③`git` 的**历史 blob**（归属）；④HEAD blob 的行号（⛔ 工作树，因 ⑥.3）。
2. **输出形状**：探针 `A1|…|` 一行一读数、`|` 分隔；错误先 render 再 `<…>`（`<nil>` 与原文可分）。
3. **数的是代码还是注释**：①② 数**文面**；③ 数**运行读数**；④ 数**栈帧行号**；`M0`/`M#` 数的是**被 grep 到的行**（含标题行）。
4. **和总数对上**：specs **14** 份（分母）；M-token 只有 `M0`×15；R-token 最大 `R9`；探针四发 **0/0/0/1**；③ 表 **6 行**＝1 关键＋1 正控＋1 组合＋1 组合正控＋1 无腿（两发）＋1 产线装配。
5. **谁算的/算法/几位**：行数 `wc -l`；diff 行数 `git diff --numstat`（⛔ `grep '^[+-][^+-]'`）；提交号 `git rev-parse`（SHA-1、40 位），⛔ `sha256sum`。
