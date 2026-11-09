# 票 323 腿 `323-a1` — 五格台件只读普查（零产码；2026-11-09 18:5x–19:2x +08）

- **腿性质**：**只读普查**（写面＝本件 ＋ `.scratch/wisp/probes/323/**` 日志）。**零产码**：`git status --porcelain -- '*.go'` 收尾为空，一枚 `.go` 未建未改。
- **起手锚**：`git rev-parse --short HEAD` = **`18ce9f5bc`**；分支 `dev`。
- **门禁口径**：★所有测试读数一律在 **`git clone --shared` 的 LF 导出树**里跑（`-c core.autocrlf=false -c core.eol=lf`），路径 `D:/work/tmp/wisp323a1/clean`，**⛔ 不在共享工作树跑**（前腿 `internal/ball` 假红就是 autocrlf 造成的）。导出树起手 `git status --porcelain` 空，收尾空。

## ★★ 起手三核（现跑，⛔ 信转述）

### 核①：派单点名的串／行还在不在

| 派单点名的东西 | 现跑读数 | 结论 |
|---|---|---|
| 票面 `323-*.md` 未勾 AC | `grep -n '^- \[' .scratch/wisp/issues/323-*.md` ⇒ **AC#0–AC#7 八枚全 `[ ]` 未勾** | 在，且**八枚全未勾** |
| 六条残余的原文 | 票面 `## 为什么立案` 六条逐枚在（第 9–14 行） | 在 |
| `isProtocol` | `internal/config/settings_320.go:142` `func isProtocol` | 在 |
| 预算常数 | `cmd/wisp/panel_models_320.go:50` `var modelsFetchBudget320 = 8 * time.Second` | 在 |
| 错误句 | `cmd/wisp/panel_models_320.go:120` 只带 `observe.ClassOfOrInternal(err)`、**⛔ 不带 `err.Error()`** | 已按 `A1483` 裁④收口 |
| 上游首行入 `HTTPError.Message` | `internal/llm/discover.go:71-77` `Message: firstLine(body)` | 在（**但已不进页面**，见 AC#4） |

### 核②：五格现在各自缺什么台件（现跑 grep 现有测试台件）

★**这是本腿最重要的读数：派单前提「五格都缺台件」已经过期四格。** 现树 `cmd/wisp/panel_models_320_test.go`（**已跟踪**，655 行）里已有四枚 AC 牙：

| 现树已有的牙 | 行 | 对应票面 AC |
|---|---|---|
| `Test320FetchLegWritesTheKeyNowhereOnDisk` | `:366` | **AC#1** |
| `Test320FetchLegTimesOutWithinItsBudget` | `:471` | **AC#2** |
| `Test320FetchLegBudgetIsTheOneThePageWaitsOn` | `:535` | **AC#2**（生产常数钉） |
| `Test320FetchLegClassifiesNon200WithoutRelayingTheBody` | `:553` | **AC#3** |
| `internal/config/settings_320_test.go:154` + `:178` | — | **AC#5** |

溯源（`git log`）：`eda34703f`（腿 `323-a0-ac1-ac2`，AC#1/#2/#3）与 `d1b1a2d22`（腿 `323-ac5`，AC#5）**均已是 HEAD 祖先**（`git merge-base --is-ancestor` 逐枚 YES）。⇒ 派单说"先派 a1 只读普查台件该怎么造"这句话，**对 AC#1/#2/#3/#5 四格已经过期**。

### 核③：写面（`git status --porcelain`）

- 起手 `git status --porcelain -- '*.go'` = **空**。
- 在飞件（**非本腿**，具名不遮）：`?? cmd/wisp/panel_resident_323_windows_test.go`（330 行）＋`?? cmd/wisp/panel_resident_assembly_323_test.go`（131 行）＝**另一条腿正在造 AC#6 的台件**，尚未 commit。本腿**一枚未碰**。
- 收尾 `git status --porcelain -- '*.go'` = **空**（本腿零产码自证）。

## 1. 五格 × 五问（本腿现跑读数）

判据口径：**"缺什么台件"＝现树已有台件对该形是否敏感**（本腿自己造突变体实测，⛔ 抄腿的号）。全部突变体在**导出树**内物理改写、跑完即还原（`git status` 空）。

| 格 | ①缺什么台件（现量） | ②最小方案 | ③落哪个包 | ④是否撞禁区 | ⑤需几轮 | 实测读数（本腿现跑） |
|---|---|---|---|---|---|---|
| **AC#1 明文落盘敏感尺** | **半缺**：配置根＋工作树两处已有两把尺（文件集/逐 sha、含密钥新文件）；**配置根与工作树之外零射程** | 把 `filesCarrying320` 的根集合从 `{dir, wd}` 扩到 `{dir, wd, os.TempDir(), 用户主目录}`；或改判据为"进程内所有 `os.WriteFile` 目标都不得含密钥"（拦截式） | `cmd/wisp`（已有台件原地扩根，**零产码**） | **不撞**（只动测试件，`cmd/wisp` 测试面本腿无写权但落地腿有） | **1 轮** | ★**实测漏洞**：m6＝把密钥写进 `os.TempDir()`（真落到 `C:\Users\swq\wisp323a1-tmpkey.txt`，39 字节）⇒ **rc=0 全绿**（`logs/m6b.log`）。对照 m1＝写进工作树 ⇒ 红，`panel_models_320_test.go:424`（`logs/m1.log`）⇒ **尺只在两个根内有眼** |
| **AC#2 超时夹具** | **不缺**（两枚牙已在）；**但一枚有偶发假红** | 无需造新台件；只需把 `attempted==1` 那条断言改成"至少 1"或加同步（见右） | `cmd/wisp` | 不撞 | **0 轮**（若要修假红 1 轮） | m3（`8s`→`8h`）⇒ **rc=1 红**，`:543` 逐字 `the fetch budget is 8h0m0s; A972 ruling 3 caps…`（`logs/m3full1.log`）；m7（**完全删掉预算**）⇒ 红三句 `:515`/`:518`/`:524`（`logs/m7.log`）⇒ 牙对"没预算"有牙。★**但**：`Test320FetchLegTimesOutWithinItsBudget` 单独跑 15 趟 **1 趟假红**（`flk3.log`：`:521 the slow endpoint was called 0 time(s)`）——夹具在 120ms 预算下**偶发先超时于请求发出之前**，这条断言是**竞态**、非缺陷 |
| **AC#3 上游错误体** | **不缺** | 无 | — | — | **0 轮** | m4（把 `err.Error()` 拼回页面句）⇒ **rc=1 红**，两形各一句 `:599`，401 与 500 **各带自己的 D37 分类**（`auth` / `provider`）且 `UPSTREAM_LEAK_320XYZ` 被逮（`logs/m4.log`）⇒ 两形都有牙 |
| **AC#4 宿主装配** | **缺，但已有在飞台件** | 即那两枚未 commit 的 `panel_resident_*_323_test.go`；本腿**只读确认其形状**：`newResidentPanelManager`（生产装配函数，`panel_resident_windows.go:342`）＋`resident_windows.go:303` 同一调用 | `cmd/wisp` | ⚠ **撞在飞腿**（同一写面，同 `cmd/wisp`）⇒ **必须等那条腿交件或收编**，⛔ 并行写 | **0 轮（等交件）** | 现树**已跟踪**测试里**无**任何一枚跑 `newResidentPanelManager` 端到端（`grep -rln` 只命中 4 处：1 处注释、1 处 AST 名单、1 处 AST 读参、1 处 pump 测试用其返回值但只读快照包）。在飞件 `Test323ResidentHostAnswersThePageRows` 正是补这一格 |
| **AC#5 `protocol` 预检** | **不缺** | 无 | — | — | **0 轮** | m5（`isProtocol`→恒真）⇒ **rc=1 红**，`internal/config/settings_320_test.go:178` 逐字 `the refusal came from the whole-file gate instead of the field's own pre-check`（`logs/m5.log`）⇒ `323-ac5` 的牙是真的 |

### 1.1 逐格"缺／不缺"一句话

1. **AC#1**：**半缺**。配置根与工作树有眼；**`os.TempDir()`／用户目录无眼**，而密钥写到那儿今天**全绿**。
2. **AC#2**：**不缺**（m3 与 m7 两向都红）。附一条**假红竞态**（1/15），不是缺台件。
3. **AC#3**：**不缺**（m4 两形皆红）。
4. **AC#4**：**缺**，但**有在飞台件**（未 commit），⛔ 本腿不重造。
5. **AC#5**：**不缺**（m5 红）。

## 2. ★如实：哪格"不需要台件"

**AC#2／AC#3／AC#5 三格现量不需要新台件**（各有已跟踪的牙，本腿逐枚造突变体实测有牙）。本腿**没有为凑表给这三格造方案**。

## 3. 门禁（零产码 ⇒ 无对象；两把只读门）

| 尺 | 命令（在 LF 导出树） | rc／读数 |
|---|---|---|
| `d22scan` | `sh scripts/d22scan.sh` | **rc=0 `clean`**（`internal/`=452 生产件；ban#8 internal/=972、cmd/=245、frontend/=121、design/=14） |
| 路径长度门 | `sh scripts/check-path-length-budget.sh` | **rc=0 `VERDICT GREEN`**（最长 185、worst full 229、roster 63=63 与树相等） |
| `go build` / `gofumpt` | — | **无对象**（本腿零产码，⛔ 无可格式化的件） |
| `Test320` 定向 | `go test ./cmd/wisp -run Test320 -count=1 -v` | **rc=0，13/13 PASS**（诚实树基线；`logs/final.log`） |

## 4. 面核与诚实自报

- **改动面**：本件 ＋ `.scratch/wisp/probes/323/a1/logs/`（23 个日志，只建不删）。**零 `.go`**。
- **突变体全在导出树内**（`D:/work/tmp/wisp323a1/clean`），跑完逐枚还原，收尾 `git status --porcelain` **空**；导出树本身是一次性 `--shared` 克隆，**共享工作树零字节改动**。
- **m1／m6 真写出了密钥文件**（`cmd/wisp/wisp323a1-keyleak.txt` 与 `C:\Users\swq\wisp323a1-tmpkey.txt`），跑完**已删**（这正是"尺有眼／无眼"的证据，不是泄漏）。
- **具名没做**：
  1. **AC#4 台件**：在飞腿正在写同两枚文件，本腿**只读不重造**（⛔ 撞写面）。
  2. **AC#1 扩根台件**：本腿**只指出缺口**，⛔ 未落地（零产码约束；`cmd/wisp` 测试面也不在本腿写面）。
  3. **`cmd/wisp` 整包**：一次没跑（既有真窗族＋假分母，铁律；本腿全部 `-run Test320` 定向）。
  4. **票面 AC 框**：八枚**一枚未动**（归编排者翻）。
- **撤销口令**：删本件与 `probes/323/a1/` 即复原；本腿本来零产码。

## 5. 给编排者的三句话

1. **派单前提过期四格**：AC#1/#2/#3/#5 的牙已在 HEAD（`eda34703f`／`d1b1a2d22`），本腿实测**除 AC#1 的根射程外全部有牙**。
2. **唯一真缺口＝AC#1 的扫描根**：密钥写到 `os.TempDir()`／用户目录**今天全绿**（本腿造出来了）。最小修法＝把那把尺的根从 `{配置根, 工作树}` 扩到 `{＋Temp, ＋用户目录}`，**1 轮、零产码、不撞禁区**。
3. **AC#4 别并行派**：同两枚文件正被另一条腿写（未 commit），要么等它交件，要么收编。
