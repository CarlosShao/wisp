# 票 231 · 对抗验收腿 `231-v1`（非实现者终裁）读数件

- 腿＝`231-v1`（裁决腿，D22 双角色：写码腿 `231-r1` 不能自判）。本件只读数与判语，⛔ 零产码净变化。
- ⛔ 票面 `.scratch/wisp/issues/231-*.md` 的五枚 `- [ ]` 框本腿一枚未翻（现量 `grep -c "^- \[ \]"`＝**5**，全程如此）。

## §0 起手锚（本腿现量，落笔前）

- HEAD＝`691470bf60b65283297610a040fcdad8da270f5b`（`git rev-parse HEAD`，`docs(evidence-close-4 §0–§3)`）；分支＝`dev`。
- 现量时刻＝`2026-10-06 09:46 +0800`。
- scoped porcelain（`git status --porcelain -- cmd/wisp internal tools scripts docs/`）＝**0 行**（`wc -l` 现量）。
- 全包工作树（与本腿写面无关，⛔ 零读零写）：` M .gitignore`／
  ` M .scratch/wisp/issues/231-*.md`（＝r1 的 Progress log 8 行 + 表头，未由编排者入库，
  `git diff --stat`＝`8 insertions(+), 0 deletions(-)`，⛔ 本腿不翻框不动它）／
  一批 `probes/**` 在飞件（`111/r4`、`268/v1`、`152`、`161/r6`）／` D design/**` 与 `frontend/**` 老脏面。
- 受验四笔产码（每枚 `git log -1 --oneline` 坐实＋`--numstat` 现量）：
  - `290dca87 票 231 AC#2 · 给"版本更高"一条自己的出口（cause=newer-build）`
    — numstat `28 0  cmd/wisp/config_reload.go`＝**纯插入零删改**。
  - `a16d1ff7 票 231 AC#3＋同批名册 · 常驻钉加一行，cause=newer-build 进两张互斥名册`
    — numstat `4 1  cmd/wisp/config_reload_223_test.go`／`4 0  cmd/wisp/config_reload_perm_223_windows_test.go`／`14 0  cmd/wisp/config_sentences_223r2_test.go`（三枚文件名由 `git show --name-only a16d1ff7` 现坐实）。
  - `c5040ea7 票 231 AC#2 收尾 · 注释不再抄既有句的原文（保住那把"零枚测试钉"的尺）`
    — numstat `1 1  cmd/wisp/config_reload.go`。
  - `42e89896 票 231 AC#3 收尾 · 测试侧注释也别抄 cause=invalid 的原句`
    — numstat `4 4  cmd/wisp/config_sentences_223r2_test.go`。
  - 写面复核＝`cmd/wisp/config_reload.go`＋三枚同包测试件；⛔ 未越界：
    `git diff --name-only 290dca87^..HEAD -- internal tools`＝**0 行**；
    `git diff --numstat 290dca87^ -- internal/config/loader.go`＝**零输出**（loader 那枚文件从未进 diff）。
- 证据件＝`.scratch/wisp/probes/231/r1/evidence.md`：`wc -l -c`＝**257 行 / 25,367 字节**（与派单口径逐字一致）；
  编排者代提＝`fd8c396e probes(231-r1 落地腿遗产代提)`；同目录 38 枚台件（本腿 `ls` 现量＝28 项＋`mutation/` 3 枚）。
- 占位尺（并集 `待\[填\]|填写\[中\]|未判|待验`，`grep -nE` 全文逐节亲读）＝**0 命中**（`RC=1`）⇒ 派单那把尺的读数坐实；
  ⚠ 但**逐节亲读**另量出两枚该把尺管不到的「待补」字样（见 §8 上报 #2）——编排者 `fd8c396e` 提交信息里的"占位 0"
  只对那把并集尺成立，不对"全文无待补"成立。
- 台件枚数复量（派单说"同目录 38 枚"）：`find probes/231/r1 probes/231/a1 -type f | wc -l`＝**38**
  （`r1` 自身＝**29 枚**、`a1`＝**9 枚**；`git show --stat fd8c396e`＝**37 files / 2,273 insertions**，
  差的第 38 枚＝`a1/census.md`，它由更早的 `2048edcb` 入库）⇒ 三处口径互不矛盾，本腿逐枚坐实。
- ★新分支位置（本腿 `grep -n` 现取）：
  - `cmd/wisp/config_reload.go:396` = `case strings.Contains(d, "was written by a newer build"):`
  - `cmd/wisp/config_reload.go:420` = `return "cause=newer-build detail=\"" +`（句体 `:421-423`）
  - `cmd/wisp/config_reload.go:424` = `case strings.HasPrefix(d, "config.toml:"):`（`:425-427` 回 `cause=invalid`）
  - 两处同在那枚内层 `switch { ... }`（`errors.As(err, &oe)` 之后，`:368` 起），**case 按书写次序求值** ⇒ 判定见 §1。
- 出口枚数尺：`grep -c "cause=" cmd/wisp/config_reload.go`＝**8**（r1 起手 7 ＋本票 1；票面"四句"＝枚数不实，与 r1 同判）。
- 影响面包数：`GOFLAGS= go list ./cmd/wisp ./internal/...`＝**31 枚包**（`cmd/wisp`＋30 枚 `internal/...`）。
- 起手 CPU／MEM＝**37.6% / 71.3%**（`Get-Counter`＋`Win32_OperatingSystem` 现读，脚本 `load.ps1`）。
  复量（等待 45s 后）＝**34.5% / 71.6%**。
  ⚠ MEM 那一条本腿复算＝**环境常驻、非在飞测试**：`TotalGB=31.9 FreeGB=9.1`；
  前八名常驻 `Memory Compression 3174MB`／`vmmemWSL 1484MB`／`DeepSeek Harness 808MB`／`Qoder CN 759+741MB`／
  `explorer 540MB`／`Weixin 485MB`／`Everything 480MB`；
  `Get-Process go,wisp,test`＝**0 枚在飞**（`INFLIGHT_GO_OR_WISP=0`）⇒ `cmd/wisp` 运行窗归本腿独占成立。
  ⇒ 本腿按派单口径把 **CPU < 70%** 放行，MEM ≥70% 具名交回（与 r1 §7 同形，非本票新增）。

## §1 新分支的位置：是不是死代码

现量三枚行号（本腿 `grep -n cmd/wisp/config_reload.go`，非 r1 转述）：

| 行 | 内容 |
|---|---|
| `:396` | `case strings.Contains(d, "was written by a newer build"):` |
| `:420` | `return "cause=newer-build detail=\"" +`（句体 `:421-423`） |
| `:424` | `case strings.HasPrefix(d, "config.toml:"):` |
| `:425` | `return "cause=invalid detail=\"" +`（句体 `:426-427`） |

判定依据（三条独立读数，⛔ 不采信 r1 自述）：

1. **静态**：`:396` 与 `:424` 在同一枚内层 `switch {…}`（`config_reload.go:368` 起，`errors.As(err, &oe)` 之后），
   Go 的表达式 `switch` 按**书写次序**逐个求值 case、首个为真即执行 ⇒ 认领串
   `Contains(d,"was written by a newer build")` 的 case（`:396`）**早于**吞它的 `HasPrefix(d,"config.toml:")`（`:424`）
   ⇒ **非死代码**。
2. **实跑正证**（当前工作树，`09:59`＋`10:00` 两发）：
   - 常驻钉 `TestTicket223R2FailureSentenceRouting/声明未来版_正文解析得开_归更高版本自己那句` 绿
     （`logs-ac3-green-on-worktree-v1.txt`，`--- PASS` 枚数＝10＝父＋9 子，`--- FAIL`＝0）；
   - AC#1 探针（本腿自己的台件，`logs-ac1-now-v1.txt`）操作员那行逐字 `cause=newer-build` 那句。
3. **实跑反证（本腿新造的一发，r1 没做过）**：把那一支**整块搬到** `HasPrefix` 那条**之后**
   （overlay 件 `mutation/config_reload_branch_after_prefix.go`＋`overlay-branch-after.json`，
   自检尺＝搬完后 `case ...newer build...` 的字符偏移必须 **大于** `case strings.HasPrefix(d, "config.toml:"):`，
   脚本 `make-branch-after.ps1` 内置不满足即 throw）跑同一枚用例 ⇒ 见 §1.1 读数。

⇒ 结论：**位置承重成立**，本票最硬那一格判「**新分支非死代码，且位置由正控与反证双向钉住**」。

### §1.1 反证读数（分支搬到 HasPrefix 之后必须红）

- 突变件 `mutation/config_reload_branch_after_prefix.go`（437 行／24,323 字节，`make-branch-after.ps1` 生成）：
  把 `:396-423` 那一整块（case＋注释＋四行 return）**原样搬到** `HasPrefix` 那三支之后，
  自检尺内置＝搬完后 `case …newer build…` 的字符偏移必须 **大于** `case strings.HasPrefix(d, "config.toml:"):`，否则 throw。
  现量该件 `:396`＝`case strings.HasPrefix(d, "config.toml:"):`、`:400`＝`case strings.Contains(d, "was written by a newer build"):`
  ⇒ 次序确实颠倒（`gofmt -l` 把它列出来＝只读为格式差异，本腿不格式化突变件；解析本身通过，见下面实跑）。
- 实跑（`2026-10-06 10:16:54 → 10:16:55 +0800`，HEAD＝`78072249`，日志 `logs-ac3-mutation-branch-after-v1.txt`）：
  `go test -overlay .scratch/wisp/probes/231/v1/overlay-branch-after.json -run TestTicket223R2FailureSentenceRouting -count=1 -v ./cmd/wisp`
  ⇒ `RC=1`、`--- FAIL` 计数＝**2**（父＋新那一枚子形）、三把尺同响（`:125`／`:129`／`:133`），
  booked 值逐字＝`cause=invalid detail="…内容被校验拒绝（值不合法或引用解不开）…"`
  ⇒ **位置一颠倒，那一支立刻变成死代码并被前缀吞掉**，与 §2 的"整支删掉"正控给出同一枚红句。
- 跑完自证：`git status --porcelain -- cmd internal`＝**0 行**。

⇒ 这一发是本腿自己加的**反证**：票面与 r1 只用"删掉必须红"证那一支有牙；
本腿另证"位置写晚＝同一枚红"，把 §1 的"非死代码"从静态读序升级为**双向实跑**。

### §1.2 AC#2 的认领面自证（本腿复算，⛔ 不采信 r1 转述）

- 认领短语 `was written by a newer build` 在**非测产码**全仓＝**2 枚命中**：
  `internal/config/loader.go:122`（唯一生产者）＋`cmd/wisp/config_reload.go:396`（消费者本身，`:412` 是注释）。
  ⇒ 除它自己那枚 case 之外**零枚**其它形状会被这一支认领 ⇒ 抢不走别的句子（r1 说的"internal/ 唯一一枚"复算成立）。
- 反向：新支**严格窄于**它竞速的那根前缀（`internal/config` 非测产码里 `"config.toml: ` 起手的 detail＝**48 枚**，本腿复尺同值），
  故 unknown-key／值不合法／引用解不开／迁移产物不合法任一形都不可能被抢到——
  并且这四形的常驻钉在 §2 正控那发里**全部仍 PASS**（8 枚子形零红）。
- 票面 AC#2 要求的三项内容逐字现读（`config_reload.go:420-423`）：
  ①由更新版本写出＝`config.toml 是由一个更新的 Wisp 写出来的（它声明的 schema_version 比这份程序懂得的高；`；
  ②继续用旧配置＝`本次运行继续用内存里的旧配置；`；③升级或恢复备份才读它＝`升级 Wisp 或恢复备份才会读它`；
  另含一支"说清自己不是什么"（`这一条不说语法错，也不说值不合法，因为它还没走到校验`），
  与票面"第二句还是假的"那一格的现象**正面对齐**（本腿 §4 复现的改前那句确实说了"值不合法或引用解不开"，而 `loader.go:120` 那支从未走到校验）。
- ⛔ 未改 loader 原文（§3 有字节尺）、⛔ 未把 `:424` 那条 `HasPrefix` 改宽（`:424` 现读逐字与 `290dca87^` 同值）、
  ⛔ 未新增导出名／C17 方法名（新支只是同一枚 `describeReloadFailure` 里多一条 `case`）。

判语：**AC#2 成立**（位置承重由正反两向实跑钉住；认领面自证窄；三项齐全；三条禁改项零触碰）。

## §2 正控重放：拿掉新分支必须红（overlay 法）

本腿**不采信** r1 的 `mutation/` 件，自己重做一遍（同法、同形状、不同物理件，全在 `probes/231/v1/` 下）：

- 切断脚本 `make-no-branch-v1.ps1`（逐字照 r1 的 `make-no-branch.ps1` 形状，写面改指 v1 目录）：
  从**当前** `cmd/wisp/config_reload.go` 现切 `mutation/config_reload_no_branch.go`，
  自检＝切完文件里不得再出现 `was written by a newer build`，否则 throw。
  本腿量到 `409 行 / 22,481 字节`。
- **三向 identical 校验**（这把尺比 r1 自证更硬）：
  - `diff v1/mutation/config_reload_no_branch.go  r1/mutation/config_reload_no_branch.go`＝**0 行**（两腿独立切断，逐字节同物）；
  - `diff v1/mutation/config_reload_no_branch.go  <290dca87^ 的 config_reload.go>`＝**0 行**
    ⇒ 切掉的那 28 行**正是** `290dca87` 插入的那 28 行，没有多切、没有少切。
- overlay 件 `overlay-no-branch.json`（形状照 r1，键值指向 v1 目录）：
  `Replace: cmd/wisp/config_reload.go → probes/231/v1/mutation/config_reload_no_branch.go`。
- 命令逐字（`2026-10-06 09:59:01 → 09:59:10 +0800`，HEAD＝`691470bf`，日志 `logs-ac3-mutation-red-v1.txt`）：
  `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" GOFLAGS= go test -overlay .scratch/wisp/probes/231/v1/overlay-no-branch.json -run TestTicket223R2FailureSentenceRouting -count=1 -v ./cmd/wisp`
  ⇒ `RC=1`（⛔ 未加 `-race`／`-cover*`／`-short`；`-overlay` 路径写法在本机 Go 1.27.1 windows/amd64 直接可用，**正控形可重放**）。

**★红句逐字**（尺在 `cmd/wisp/config_sentences_223r2_test.go:125`，本腿这一发日志第 39 行）：

```
config_sentences_223r2_test.go:125: shape 声明未来版_正文解析得开_归更高版本自己那句 is booked "cause=invalid detail=\"config.toml 语法没问题，但内容被校验拒绝（值不合法或引用解不开）。本次运行继续用内存里的旧配置\"", want it to open with "cause=newer-build" - the two sentences swapped again
```

同发另外两把尺也同时响（`r1` §3 说"三把尺全响"，本腿复现）：

```
config_sentences_223r2_test.go:129: shape 声明未来版_正文解析得开_归更高版本自己那句 borrows "cause=invalid", which is a different pipeline's answer:
config_sentences_223r2_test.go:133: shape 声明未来版_正文解析得开_归更高版本自己那句 says "语法没问题", which contradicts its own cause:
```

- 计数尺：本发 `^--- FAIL`＋`^    --- FAIL`＝**2**（父用例＋新那一枚子形）、`^    --- PASS`＝**8**
  ⇒ 这颗牙**只咬新那一形**，其余 8 枚子形照绿。
- **红句逐字节比对**：本腿这发第 39 行剥缩进后＝**354 字节**，与 r1 证据件 `:106` 那行（剥缩进、同尺）**354 字节**，
  `diff`＝**0 行** ⇒ `RED_SENTENCE_BYTE_IDENTICAL_TO_R1_RECORD`（红句不是转述，是同一枚输出）。
- 跑完自证跟踪文件零改动：`git status --porcelain -- cmd internal`＝**0 行**（现量，日志件 §2 末）。

判语：**AC#3 的正控成立**——拿掉新分支这一枚常驻钉必红，红句逐字与 r1 记录逐字节同；
overlay 法在本机可重放，⛔ 未改用"直接编辑源码"做突变。

### §2.1 推论复现（⛔ 无常驻钉时拿掉分支不会红＝今天那颗牙确实来自"加行"）

- 双突变 overlay `overlay-no-branch-and-no-pin.json`：
  键 1 把 `config_reload.go` 换成 §2 的无分支件，键 2 把 `config_sentences_223r2_test.go` 换成
  `mutation/config_sentences_223r2_pre_ac3_test.go`（＝`git show 290dca87^:…` 现取，124 行、表里**没有**那枚新行）。
- 实跑（`2026-10-06 10:18:51 → 10:18:52 +0800`，日志 `logs-ac3-no-branch-no-pin-v1.txt`）
  ⇒ `RC=0`、`--- PASS` 子形 **8 枚**、`FAIL` 命中 **0**：那一形根本没人种，删掉分支不响任何一枚。
- 对照 §2（同一枚无分支件、**只有**加了行的表）＝`RC=1` 三把尺全响
  ⇒ 本腿独立坐实派单第 3 条的推论：**这一格的牙是"加行"给的**，不是既有 (a) 那枚 J1 冒充的
  （J1 在两发里都 PASS，`--- FAIL` 枚数始终＝1/2 枚且名字只有新那一形）。
- 常驻钉在位读数：`config_sentences_223r2_test.go:97-100` 现量该行＝
  `"声明未来版_正文解析得开_归更高版本自己那句", "schema_version = 99\n\n[ball]\nsize = 64\n"` ／
  `"cause=newer-build", "cause=invalid", "语法没问题"`（三把尺复用现成件，未新增尺）；
  `a16d1ff7` 对该文件的 numstat＝`14 0`＝**纯插入**，`42e89896`＝`4 4` 只改注释四行（hunk 现读：
  删掉的两行注释曾整段抄 `内容被校验拒绝（值不合法或引用解不开）`，改成本腿自己的话），
  ⛔ 未删任何断言、未改任何 needle。
- 两张互斥名册（票派单第 5(i) 条）本腿现量枚数：
  `config_reload_223_test.go:572-576` 的 `all`＝**9 枚**（含 `"cause=newer-build"`，`:575`）；
  `config_reload_perm_223_windows_test.go:112-118` 的名册＝**7 枚**（含 `"cause=newer-build"`，`:118`）。
  ⇒ r1 说的"8→9／6→7"复认。
- 派单第 5(ii) 条那颗暗坑（perm 用例把 `state=not-applied cause=invalid` 当**控制流探测器**，
  现读 `config_reload_perm_223_windows_test.go:142-146`）：本腿在 §6 那发整包里
  `cmd/wisp`＝**ok 464.811s** 复认它没撞（⛔ 零枚 `three plants all landed…`／
  `neither the permission sentence nor an adoption arrived…`），
  并另跑两发独立确认（下条）。
- perm 用例独立复跑（防"整包绿但隔离跑红"这一形）：
  - `10:17:0x` 前后各一发 `go test -run TestTicket223PermissionDeniedSitsInItsOwnSentence -count=1 -v ./cmd/wisp`
    ⇒ 见 `logs-perm-repeat-v1.txt`（两发各 PASS，`t.Fatal` 那两枚形状零命中）。


## §3 四枚既有句的字节尺（⛔ 不拿"用例全绿"当凭据）

先复认普查腿那把尺（⛔ 不许拿"用例全绿"当凭据的前提坐实）：

- 派单给的尺逐字＝`grep -rn "内容被校验拒绝" cmd internal tools --include=*.go` ⇒ 现量命中 **1 枚**＝
  `cmd/wisp/config_reload.go:426`（产码自己的那行）⇒ **测试侧零枚正向钉这四句中文全文**，复认成立。
- 四句的中文片段在 `_test.go` 里的命中（本腿逐枚扫）：
  `读不到：文件不存在`→1 枚、`没有读它的权限`→1 枚、`读到了但解析不了`→**0 枚**、`本宿主没有接 config.toml 轮询`→**0 枚**。
  ⚠ 前两枚命中**都不是正向钉**：`firstrun_198_test.go:245` 把 `config.toml 读不到：文件不存在`
  放进**"首启回执不许借用热加载措辞"的否定名册**（`if strings.Contains(...) { t.Errorf }`）；
  `config_reload_perm_223_windows_test.go:126` 钉的是 `没有读它的权限` 这半句（权限那支的另一枚半句），
  不是缺失句 ⇒ **票面 AC#2 那"四句一字不改"今天仍然没有门**，与普查腿 `231-a1` 同判。
- 逐枚句的"命中行号＋前后同值"读数（尺＝片段 `grep -n` 在 `290dca87^` 版本与工作树各取一次）：

| 句 | `290dca87^` 命中行 | 现工作树命中行 | 该行逐字节 |
|---|---|---|---|
| 缺失（`:357-359` 支体） | `:358` | `:358` | **同址同值** |
| 权限不够（`:361-363`） | `:362` | `:362` | **同址同值** |
| 语法错（`:370-372`） | `:371` | `:371` | **同址同值** |
| 热加载被禁用（常量 `:97-99`） | `:98` | `:98` | **同址同值** |

- 组尺（比逐枚更硬的一把）：把七组受保护文本（上表四枚＋`:374-376` unknown-key＋`:393-395` migration＋
  invalid 那三支 `@397-399`/`@425-427`）按"标签行＋内容行"抽出，⛔ 只比内容行：
  - `protected-lines-before-290dca87.txt` 与 `protected-lines-now.txt` 的**内容行拼接**＝各 **1,818 字节**、
    `diff`＝**0 行**、`sha256`＝**同一枚** `b8a2430987fb4ac40cc0b588766858f95de6c1bdd4143a4a7c78192c97ef25ad`
    ⇒ 四句（含另外三支）**逐字节零变化**坐实，invalid 那支只是**整体位移 ＋28**（`397→425`）。
- 独立叠尺＝`git show --numstat 290dca87`＝`28 0 cmd/wisp/config_reload.go`（**零删除** ⇒ 既有句所在行不可能被改）、
  `git show --numstat c5040ea7`＝`1 1`（本腿复查那一枚 hunk 见下）。
- `c5040ea7` 那一笔 `1 1` 到底改了哪一行（⛔ 不能只看 numstat 就说"没动产品字符串"）：
  hunk＝`@@ -400 +400 @@`，唯一改的那一行是**注释**（`:400`），
  改前注释里整段引了 `cause=invalid` 的原句、改后换成本腿自己的话 ⇒
  产品字符串那七组（`:97-99`／`:357-359`／`:361-363`／`:370-372`／`:374-376`／`:393-395`／`:425-427`）**零碰**，
  上面那把 1,818 字节尺就是它的凭据。
- `internal/config/loader.go:120-124` 英文原文复认（现量＋与 `290dca87^` 逐字节比）：
  `git diff 290dca87^ -- internal/config/loader.go`＝**0 行**；`:122` 逐字＝
  `config.toml: schema_version %d was written by a newer build (this build understands %d); upgrade Wisp or restore a backup`
  两枚钉手现读未改：`loader_test.go:273`＝`if !strings.Contains(err.Error(), "99")`、
  `migrate_test.go:160`＝`if err == nil || !strings.Contains(err.Error(), "newer build")`。

判语：**AC#2 的"四句一字不改"以字节尺成立**（同址同值四枚＋1,818 字节组尺 sha256 同值＋零删除 numstat）；
⚠ 同时**具名上报**：这四句今天**没有任何一枚测试正向钉**，字节尺只对本票这一笔有效，
"零枚门"这一条 r1 在 §4 登记 #2 已写、本腿复认，归编排者决定是否升格。

## §4 AC#1 现复现：操作员实际看到的那句原文

- 本腿**不用** r1 的探针件，自建同形状最小版
  `zz_231v1_ac1_probe_test.go`（⛔ 零断言、只读数；经 `overlay-ac1.json` 映射进 `cmd/wisp`，工作树从未存在）：
  走票 223 真宿主 `newReloadRun223 → live → awaitAudit(state=armed) → writeOver → awaitAudit(state=not-applied cause=)`，
  种 `schema_version = 99\n\n[ball]\nsize = 64\n`（正文解析得开＝唯一能走到 `loader.go:120` 那支的形状）。
- **现态（工作树＝四笔都在）**读数 `logs-ac1-now-v1.txt`，取数 `2026-10-06 10:00:47 → 10:00:57 +0800`，RC=0：

```
[audit] config: HOT-RELOAD state=not-applied cause=newer-build detail="config.toml 是由一个更新的 Wisp 写出来的（它声明的 schema_version 比这份程序懂得的高；这一条不说语法错，也不说值不合法，因为它还没走到校验）。本次运行继续用内存里的旧配置；升级 Wisp 或恢复备份才会读它"
```

- **改前那一形**（本腿自己复现，⛔ 不改跟踪文件）：同一枚探针 + 同一份种子，
  另叠一层 overlay 把 `config_reload.go` 换成 §2 的无分支件
  （`overlay-ac1-against-no-branch.json`，两键 Replace）⇒ `logs-ac1-before-shape-v1.txt`，
  取数 `2026-10-06 10:01:44 → 10:01:52 +0800`，RC=0，操作员那行逐字：

```
[audit] config: HOT-RELOAD state=not-applied cause=invalid detail="config.toml 语法没问题，但内容被校验拒绝（值不合法或引用解不开）。本次运行继续用内存里的旧配置"
```

  ⇒ 这一行与 r1 `logs-ac1-before-fix.txt`（HEAD `b5418b6c` 实跑）那句**逐字同**、与本腿 §2 红句里的 booked 值**逐字同**；
  与票面现量表那张表的 `cause=invalid` 那一条也同句——但本腿的读数是**实跑**，不是票面文本。
- 通道核验（⛔ 未新增 stdout）：两发的操作员那行都在 `state=not-applied` 审计 trail 里，
  来源＝`config_reload.go:155` 那一枚 `rt.auditf("config: HOT-RELOAD state=not-applied %s", describeReloadFailure(err))`（现读 `:154-155`）
  ⇒ 新分支走**既有那条审计通道**，票面"操作员听到的那句"就是它。
- 跑完自证：`git status --porcelain -- cmd internal`＝**0 行**（两发各测一次）。

判语：**AC#1 成立**（改前那句由本腿独立复现并与 r1 记录同句；改后同一形状走自己的出口）。

## §5 AC#4 前缀别再身兼两职：只登记不实现

票面 AC#4 **允许只登记不实现**，但要具名"哪枚文件哪一行为什么今天不动"，且落一条 DEFERRED 登记。

- **"哪枚文件哪一行"具名核对**（r1 §4 那五行，本腿逐枚现读到行）：

| r1 具名 | 本腿现量 | 判 |
|---|---|---|
| `cmd/wisp/config_reload.go:424`（`HasPrefix(d,"config.toml:")`） | `:424` 该行逐字存在，且 `:424-427` 那支 return 未改 | **行号真、理由成立** |
| `internal/config/loader.go:120-124` | 该五行是 `ver > SchemaVersionCurrent` 那支，英文原文 `:122`，钉手两枚在位 | **行号真、理由成立** |
| `internal/config/migrate.go:71-73`／`:78-80` | 现读两枚都是 `"config.toml: migration to schema version %d produced an invalid config (%v); the file was left untouched"`（`:72`／`:79` 起句），两支不含 `cannot migrate…`／`no migration registered` 任一串 ⇒ 今天照样被 `:424` 吞成 `cause=invalid` | **行号真、吞点复算成立** |
| `cmd/wisp/config_reload.go` 函数头注 `:342-353` | 现读注释仍写 "four different things … so the four stay four"（`:345`／`:349`），函数实际 `grep -c "cause="`＝**8** ⇒ 注释过期属实 | **行号真、读数属实** |
| `DEFERRED(...)` 代码标记一枚未加 | `git grep DEFERRED` 于四枚写面文件：唯一命中＝`config_reload.go:31`，且该句**在 `290dca87^` 就存在**（本腿 `git show 290dca87^:cmd/wisp/config_reload.go \| sed -n '28,33p'` 复现同文）⇒ 与本票无关、非本腿新增 | **属实** |

  另外三条缺尺的具名读数也复认了：`config.toml: ` 前缀在 `internal/config` 非测产码里
  `grep -rn '"config\.toml: ' internal/config/ --include=*.go \| grep -v _test.go \| wc -l`＝**48**（与 r1 同值）；
  `observe.New/Wrap(observe.ClassConfig …)` 在 `internal/config` 非测源码＝**85 枚**（与 r1 同值）；
  票 268 那枚 AST 尺＝`cmd/wisp/resident_approval_risk_268_windows_test.go:224`
  `TestTicket268RefusedBranchClassifiesBySentinelNotByErrorWords`（文件 293 行 ⇒ r1 写的 `:224-293` 段在射程内），
  `resident_approval_windows.go:447-456` 那段"不能读 error prose 所以故意不细分"注释现读在位（`:449-453`）。

- ⚠ **一处具名不实（本腿量到，须上报）**：r1 §4 登记 #4 写
  "`observe.Error.ProviderCode` 是已导出字段、**全仓非测产码零枚赋值**"。
  本腿尺＝`grep -rn "ProviderCode" --include=*.go internal cmd tools | grep -v _test.go | grep -E "ProviderCode\s*=\|ProviderCode:"`
  ⇒ **命中 5 枚赋值点**，全在 `internal/llm/errors.go`（`:42`／`:66`／`:239`／`:241`／`:248`），
  `internal/audio/device.go:131` 另有一枚**读取**比较。⇒ 该句在**限定范围**上不实：
  真话是「`internal/config` 非测产码零枚赋值」（本腿复尺＝`grep -rn ProviderCode internal/config --include=*.go`＝**0**），
  而"全仓零枚赋值"错。
  ⛔ **不影响本票任何一格**：那枚字段本腿未填、`Error()` 渲染面（`internal/observe/errors.go:194-196` 现读逐字
  `if e.ProviderCode != "" { fmt.Fprintf(&b, " (%s)", e.ProviderCode) }`）零触碰，它只是**登记件里的一处引用误差** ⇒
  本腿判"AC#4 交付形状成立、该句措辞须更正"，交编排者（⛔ 本腿不代改 r1 的登记文字，那枚文件不是本腿写面）。

- **DEFERRED 双向 1:1 核**：本票**没有**代码标记 ⇒
  按 AGENTS.md §1.1 与 `SPEC-12 §5`，"标记↔登记表"这一对 today 是**空转**（本腿旁证：台账/表侧不属本腿写面，未核表行数）；
  r1 把登记落在证据件 §4，⛔ 那是**证据件不是代码标记** ⇒ 与"不许单方面造 `DEFERRED(D-xx)`"一致，
  本腿判这一条**处理正确**（是否升格 `A##`／`DEFERRED` 归编排者拍）。

判语：**AC#4 成立**（"只登记不实现"允许的形状；五行"哪枚文件哪一行为什么不动"逐枚行号本腿复现属实、
三条读数尺（48／85／268 那枚 AST 尺）复算同值）；⚠ 附带**一处具名不实上报**（`ProviderCode` "全仓零枚赋值"＝错，
真值＝`internal/config` 零枚／全仓非测 5 枚），不改判语成立性、交编排者处置。

## §6 整包终态：逐名红名册＋有没有第六枚

命令逐字（本腿自己跑的那一发，⛔ 无 `-race`／`-cover*`／`-short`）：

```
PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" GOFLAGS= go test -count=1 ./cmd/wisp ./internal/...
```

- 取数窗口＝`2026-10-06 10:07:09 → 10:15:03 +0800`（**474 秒**），HEAD＝`78072249`，
  `FULL_RC=1`（整包非零是预期：名册里有历史在册红），日志 `logs-ac5-full-v1.txt`（110 行），
  环境件 `logs-ac5-full-v1.env`（`go version go1.27.1 windows/amd64`＋两目录 DLL 现量各 3 枚）。
- 起手 `INFLIGHT_GO_OR_WISP=0`＋CPU 49.5% ⇒ `cmd/wisp` 运行窗本腿独占（在飞三枚腿全程零 go 命令，未撞）。
- 包级算术（防"把没跑读成绿"的第一把尺）：`^ok`＝**24 枚**、`^FAIL <pkg>`＝**3 枚**、`[no test files]`＝**4 枚**
  ⇒ 24＋3＋4＝**31 枚**＝起手 `go list` 枚数，**逐枚对得上、零枚失踪**。
  `[no test files]` 四枚逐名＝`internal/agent/scheduler`／`internal/speech`／`internal/streamkey`／`internal/watchdog`
  （与 r1 §5 那四枚**逐字同名单**，⛔ 非本票新增）。
- 第二把尺＝`grep -c "0xc0000135"`＝**0**、`grep -c "SKIP"`＝**0**（两发比口一致），
  `exit status` 命中＝0 ⇒ 没有一枚是"用例根本没跑"的形状。

**逐名红名册（本腿这一发 `^--- FAIL`＝6 枚，缩进子形红＝0 枚）**：

| # | 包 | 红用例（逐名） | 本腿现读的那一句 | 归属判定 |
|---|---|---|---|---|
| 1 | `internal/ball` | `TestC21TableColourRowsMatchTokensCSS` | `tokens_table_test.go:1468: read design/assets/tokens.css: … cannot find the path specified - the CSS leg of this check must never skip` | 历史在册（`design/assets/tokens.css` 工作树是 ` D` 老脏面） |
| 2 | `internal/panel` | `TestApprovalCardViewJSONKeysMatchFrontendTypes` | `approval_test.go:129: Go Snapshot emits [instructions tasks] that interface PanelSnapshot does not declare` | 历史在册（`frontend/**` 地界） |
| 3 | `internal/panel` | `TestComposerContractTypesMatchFrontend` | `composer_test.go:74: … emits [git currentModel modelKnown credentialState credentialKnown] that interface ComposerState does not declare` | 历史在册（同上） |
| 4 | `internal/panel` | `TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme` | `frontend_hygiene_test.go:216: a second style source … right-rail.tsx:89` | 历史在册（同上） |
| 5 | `internal/panel` | `TestC21DesignTokensFourWayAgree` | `tokens_fourway_test.go:441: read design/assets/tokens.css: … cannot find the path specified` | 历史在册（同 #1 那枚缺失文件） |
| 6 | `internal/risk` | `TestResolvePerCallBudget` | `pathresolver_budget_norace_test.go:37: C26 budget breach: Resolve averages 5.521666ms per call, budget 1ms`（`:34` 同发读数 279 samples） | **争用型假红**，非本票新增（判据见下） |

- **①`cmd/wisp` 与 `internal/config` 必须 ok**：现量 `ok cmd/wisp 464.811s`（日志第 4 行，本腿自己那发）＋
  `ok internal/config 4.933s`（第 55 行）⇒ 成立。
  ⛔ 本票射程内**零枚红**：`cmd/wisp` 那一发含票 223 全族＋本票新那一枚子形，整包绿。
- **②与票面历史在册逐名比对 ⇒ "有没有第六枚红"的明确回答**：
  票面在册＝`internal/ball` 1＋`internal/panel` 4（含 `TestApprovalCardViewJSONKeysMatchFrontendTypes`、
  `TestComposerContractTypesMatchFrontend`）＝5 枚。本腿这一发的**前五枚逐名、逐句对上**（#2/#3 的名字与
  票面截断名 `TestComposerContractTypesM…` 全名一致，本腿 `grep "^--- FAIL"` 现坐实）。
  ⇒ 多出来的那一枚名字＝`TestResolvePerCallBudget`，它**不是第六枚新增红**：
  它是票面自己点名的争用型假红，且 r1 的两发（`logs-ac5-full.txt`、`logs-ac5-terminal.txt`）里它**同样红**、
  r1 的第三发（`logs-ac5-final.txt`）里它**绿**——本腿复尺那三份件，红名册逐枚如下：
  `full`＝6 枚（同本腿名单逐字同）、`terminal`＝6 枚（同）、`final`＝5 枚（少 risk 那一枚）。
  ⇒ **跨腿红名册差集＝只有 `TestResolvePerCallBudget` 一枚，且它是计时型**。
  **回答：没有第六枚属票 231 的新增红；本票射程内零枚红；名册外零枚。**
  ⛔ 前五枚属别人地界，本腿**一枚没顺手修**（`git status --porcelain -- cmd internal`＝0 行即其凭据）。
- **③`internal/risk` 定它那一发（本腿自己另跑的安静 `-count=3`）**：
  `go test -run TestResolvePerCallBudget -count=3 -v ./internal/risk`
  取数 `2026-10-06 10:17:22 → 10:17:30 +0800`（起手 CPU **29.6%**／MEM 68.7%、`INFLIGHT_GO_OR_WISP=0`），
  HEAD＝`4b15869d`（其间落进 `78072249`／`d12a6784`／`4b15869d` 三笔，本腿尺
  `git diff --name-only 78072249..4b15869d -- cmd internal`＝**0 行** ⇒ Go 码同一态，读数可比），日志
  `logs-risk-quiet-count3-v1.txt` ⇒ `RC=0`、
  `--- PASS` **3 枚**（1.62s／1.69s／3.02s）、`ok internal/risk 6.357s`、`budget breach` 命中＝**0**。
  ⇒ 判定＝**争用型假红坐实**（整包那发它红在 `10:07:36`＝包级并行＋其余 30 枚包同时编译的窗口；
  安静窗隔离跑三枚全绿），与票面口径、r1 §5 判定一致，⛔ 非本票新增。
  ⛔ 本腿未动 `thresholds.go`、未动那枚 1ms 预算、未放宽任何断言。
- 与 r1 的耗时对照（本腿读数比它慢，具名不解读）：`cmd/wisp` 本腿 464.811s vs r1 三发 405.208/427.669/431.636s；
  `internal/config` 本腿 4.933s vs r1 三发 2.718/3.098/5.298s ⇒ 同一枚包级 ok，差异属机器窗口。

判语：**AC#5 成立**——整包到终态（`rc=1` 全因历史在册＋一枚计时型假红）、
`cmd/wisp`／`internal/config` 双绿、逐名红名册比票面**零枚新增**、`internal/risk` 以安静 `-count=3` 判为假红、
`[no test files]`＝4 枚（同 r1 名单）、`SKIP`＝0、`0xc0000135`＝0。

## §7 门禁四门

| 门 | 本腿现跑读数 | 时刻（+08） |
|---|---|---|
| `sh scripts/d22scan.sh` | **RC=0**；正控 `runtests.sh: OK - packages=[./...] top-level: PASS=35 FAIL=0 SKIP=0, === RUN=77, '[no tests to run]'=0`；实扫 `d22scan: clean - no D22 ban violations`；bans #1-5 internal=228／cmd=38，#6 frontend=85，#7 tools=23，#8 design=39／frontend=85／internal=512／cmd=104；另报 `skipped as git-ignored: 1 file(s) under frontend/dist/assets/` | `10:03:14`＋复采 `10:04:47`（同值） |
| `sh scripts/check-path-length-budget.sh --with-self-test` | **RC=0**；`positive control PASSED`；`denominator: tracked paths=6213  over-budget=57  covered by roster=57  not in roster=0`；`VERDICT GREEN`；`unit cross-check: … SAME 57 paths` | `10:03:54`＋复采 `10:05:24`（同值） |
| `GOFLAGS= go vet ./cmd/wisp/` | **RC=0**（零输出） | `10:04:21`＋复采 `10:05:33`（同值） |
| `"D:/work/base/gopath/bin/gofumpt.exe" -l cmd/wisp` | 输出**只剩 1 枚**＝`cmd\wisp\models.go`；⛔ 非本票（`git log --oneline 290dca87^..42e89896 -- cmd/wisp/models.go`＝**0 行**；它上一次被碰是 `5e8748b3` 那笔 212/258 代笔），本腿**未格式化它**，具名即可 | `10:05:35` |

- 瞎尺对照（派单第 2 枚尺的负证）：`ls ~/GOPATH/bin/gofumpt.exe`＝**No such file**
  ⇒ 用那把尺会得到"零枚脏"的假读数，具名坐派单警告成立。
- 四门合并日志＝`logs-gates-v1.txt`（38 行／3,271 字节，含每发 RC 与时刻）。

⚠ **门禁与整包之间的 HEAD 漂移**（本腿量到，具名上报）：起手锚＝`691470bf`，
门禁复采时 HEAD＝`78072249`（其间落进 `9a8a3c5a`／`27870302`／`e22522ad`／`78072249` 四笔）。
本腿尺＝`git diff --name-only 691470bf..HEAD -- cmd internal`＝**0 行** ⇒ 漂的四笔全是
`probes/**`＋`docs/evidence/s1/**`＋`docs/reports/pending-and-issues.md`，
⛔ 未触及任何 Go 代码，四门与整包读数的被验对象仍是同一枚产码态。

## §8 判不动的地方／本腿没做的／须编排者处置

**判语汇总（本腿五格）**：AC#1 成立（§4）／AC#2 成立（§1＋§1.2）／AC#3 成立（§2＋§2.1）／
AC#4 成立（§5，形状＝只登记不实现）／AC#5 成立（§6）——**票 231 判「通过」**，附带下面 5 处上报。

1. **一处具名不实（须更正文字，不影响判语）**：r1 §4 登记 #4 写「`observe.Error.ProviderCode` …
   全仓非测产码零枚赋值」。本腿尺＝非测产码赋值点 **5 枚**，全在 `internal/llm/errors.go`
   （`:42`／`:66`／`:239`／`:241`／`:248`），另 `internal/audio/device.go:131` 是一枚读取比较。
   真值＝「`internal/config` 非测产码零枚赋值」（本腿复尺＝0）。⛔ 本腿不代改 r1 的登记件（那不是本腿写面）。
2. **一把尺的口径（正是派单警告的那一形）**：并集尺 `待\[填\]|填写\[中\]|未判|待验` 对 r1 证据件＝0 命中（`RC=1`），
   但**逐节亲读**另量出 2 枚管不到的字样——`:30` 节状态表「§5 … 待读数落地后补齐」、`:248`「读数见 §5（待补）」。
   两处指的是同一件事，且那份读数**今天已由本腿 §6 独立补齐**（不是 r1 补的）⇒
   `fd8c396e` 提交信息里的「占位 0」对那把尺成立、对「全文无待补」不成立。⚠ 请后续腿把并集尺扩到 `待补|待读数`。
3. **r1 §5 表格的 HEAD 归属有一处不可能**：表里甲发（`09:09:59 → 09:17:23`）标 HEAD＝`c85a9b65`，
   而 `c85a9b65` 的提交时刻＝`09:19:31`（`git log --format=%h %ad` 现量），甲发起手时它还不存在；
   `logs-ac5-full.env` 里也**没有** HEAD 行（只有 start/rc/end，与 final/terminal 两件不同形）。
   同表「乙含本腿三笔，末笔 `c5040ea7`」经 `git merge-base --is-ancestor` 复尺＝**成立**
   （`c5040ea7` 是 `c85a9b65` 的祖先，提交时刻 09:19:28 早于乙发 09:19:58）。
   ⚠ 第四笔 `42e89896`（09:32:35，测试侧注释四行）**不在 r1 §5 那张表里**；
   目录里的第三发 `logs-ac5-terminal.txt`（HEAD `11b8f288`，09:33:03 → 09:40:18）经祖先核＝**覆盖它**，
   红名册＝本腿同一份六枚名单。⇒ 「四笔全入库之后的整包终态」这一格，**本腿 §6 那一发（HEAD `78072249`，
   `42e89896` 已核为其祖先）是最硬的一枚凭据**，本腿判语以此发为准，不依赖 r1 表格的那处归属误差。
4. **一处不可核转述（本腿未据此做任何判断）**：r1 Progress log `23:2x` 那行与本件 §2 引用的
   「编排者已裁 `cause=newer-build`，理由在派单第 4 条」——本腿手里的派单没有编号条款可核。
   本腿改从**票面**坐实：`.scratch/wisp/issues/231-a-*.md:26` 的 AC#2 逐字写作
   「（形如 `cause=newer-build detail="…"`）」⇒ 那枚 token 有票面根据，**无需**采信该转述。
5. **票面现量表自身两处读数已过期（与本腿无关，具名给编排者）**：
   `:10` 说 loader 那支在 `internal/config/loader.go:106-110` ⇒ 现量 **`:120-124`**；
   `:11` 说吞点在 `cmd/wisp/config_reload.go:357` ⇒ 现量 **`:424`**（漂 ＋67）。
   票面 `:6` 自己写了「起手逐条复算，别信这里的行号」⇒ 属预期漂移，⛔ 不是缺陷。
   另一处**票面 AC#3 前半句不实**本腿也复认了：`:15`/AC#3 说 `config_sentences_223r2_test.go:67-69` 钉着
   「(b)(c)(f) 三形」，本腿 `git show a16d1ff7^:… | sed -n '66,70p'` 现量那里**只有一枚** `声明未来版_正文语法坏_J1`＝(a) 形
   ⇒ 按字面「改成断新出口」＝删掉 (a) 唯一的常驻钉，与 AC#3 后半句「保留 (a) 那形」自相矛盾；
   r1 选「表里加一行＋J1 一字未动」＝**唯一同时满足前后半句的形状**，本腿复尺 J1 在正控那发仍 PASS、
   票面 (a) 那形今天仍归 `cause=syntax`（§2 日志 `--- PASS: …/声明未来版_正文语法坏_J1`）。
   ⚠ 台账侧 A441／A444 存在且标题与票面 `:3`／`:27` 的引用同形（`docs/reports/pending-and-issues.md:9452`／`:9478`），
   本腿只核到"存在＋标题一致"，未逐字读全文内容。

**本腿没做的（具名）**：

1. ⛔ 未翻票面任何一枚 `- [ ]` 框：现量 `grep -cF -- "- [ ]"`＝**5**、`grep -c "\[x\]"`＝**0**（交件时复尺同值）。
   ⚠ 顺带具名：派单/记忆里那把 `grep -c "^- \[ \]"` 尺在本仓 ERE 下会把 `[ ]` 当字符类（实测 `grep -cE "^[[:space:]]*- [ ]"`＝**0**，
   而 `grep -cF`＝5）⇒ 框数请用固定串尺，10-02 记忆里那条「换尺点不可相减」在这一枚票上再次现形。
2. ⛔ 未修名册外任何一枚红（`internal/ball` 1／`internal/panel` 4／`internal/risk` 1 全部原样交回），
   ⛔ 未动 `design/**`／`frontend/**`（零读零写零转述，唯一提及处是失败输出里的文件名字符串，那是包内产码自己读的）。
3. ⛔ 未加 `DEFERRED(D-xx)` 代码标记、⛔ 未动 `SPEC-12 §5` 表与台账、⛔ 未改 loader 那句原文、
   ⛔ 未把 `:424` 那条 `HasPrefix` 改宽或改窄、⛔ 未动 `thresholds.go`／golden／`allowlist.txt`／值域 `[31,3600]`／
   `[risk].l1_window_sec` 钳位、⛔ 未格式化 `cmd\wisp\models.go`、⛔ 未 push、仓内零删除。
4. ⛔ 未读未引三枚在飞腿的半成品目录（`probes/111/r4`、`probes/evidence-close/4`、`probes/pool-validity/2`）；
   它们的名字只出现在 `git status`／`git log` 的本腿自证输出里，⛔ 未引其内容。
   ⚠ 本腿注意到其间 HEAD 从 `691470bf` 漂到 `78072249`／`4b15869d`，
   但 `git diff --name-only 691470bf..HEAD -- cmd internal tools`＝**0 行** ⇒ 被验产码态全程同一。
5. **没有把 `cause=newer-build` 做成 stdout 行**——票面没要求，本腿复尺七句仍全走
   `config_reload.go:154-155` 那一枚 `auditf`（现读 `rt.auditf("config: HOT-RELOAD state=not-applied %s", …)`），
   与普查腿 §1 前提同形，本腿未据此加任何产品动作。
6. **没有独立重跑票 223 全套 22 枚 trail**（r1 有 `logs-ticket223-*` 七发）：本腿的凭据是整包那一发
   `ok cmd/wisp 464.811s`（含票 223 全族）＋两发 perm 用例隔离复跑（`logs-perm-repeat-v1.txt`，
   `10:20:13`／`10:20:20` 各 PASS 2.74s／2.80s，⛔ 零枚 `three plants all landed…`）。
   派单第 5(ii) 那颗暗坑（perm 用例 `:142-146` 把 `state=not-applied cause=invalid` 当控制流探测器）＝**未撞**。
7. ⛔ 未对 r1 的 `mutation/` 件做"看起来一样就用"的处理：本腿自己切断、自己写 overlay、
   再与 r1 那枚逐字节比（`diff`＝0 行）作为交叉验证。

**本腿的写面与自证**：`.scratch/wisp/probes/231/v1/**`（＋票 231 的 Progress log 行）；
每一发跑完复尺 `git status --porcelain -- cmd internal tools scripts docs/`＝**0 行**（逐发记录见各日志件）。

