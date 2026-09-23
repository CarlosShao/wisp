# 票 128 — 对抗验收（被验版本 = 锚定 sha `58302cc`）

**验收会话：** `acceptor-ticket128-adversarial`（本枚独立代理，非实现者）
**被验版本：** **锚定 sha `58302cc`**（代码在其父辈 `4e5d240`，AC#2 裁定入面 `7a5dab2`）· 分支 `dev` · 日期 2026-09-23（首条读数 11:52 +08:00）
**快照：** `git archive 58302cc | tar -x -C /tmp/wisp128-acc-z7`（**仓外、带会话后缀 `-z7`**；Windows 实路径 `C:\Users\swq\AppData\Local\Temp\wisp128-acc-z7`）
⛔ 五发变异（实为九发，见 §3）**全部落在上面那枚快照树**里，逐发还原；生产工作树在本段只新增本文件一枚（§7 的 `git diff` 自证）。
**范围：** 票面 AC#2（判定用例）+ AC#3（变异自证）为主；AC#1 只抽验"支撑 AC#2 的那几条"；**AC#4 门禁那一格实现者自己写明未做** ⇒ 本表判 `NOT-RUN（未被声称）`，不计退回。
**没有取任何时序/RSS 读数**（本机 CI 有 run 在飞：`35815467601` in_progress @03:43:38Z；§3-M-D′ 那 120 秒是 guard 触发，不作耗时读数）。
**读过的输入：** `.scratch/wisp/issues/128-*.md`（全票面 + 三条腿 AC）、`docs/evidence/s1/128-ac1-consequences.md`、`docs/evidence/s1/128-ac2-refusal-and-ac3-mutation.md`、`docs/reports/pending-and-issues.md` 的 `A105①⑤⑦`／`A109①②`。

## 0. 判据仪器清单（先把仪器本身摊开，再攻它）

被验的判据 = 5 枚用例 / **16 行 `--- PASS`**（第 0 列 5 + 缩进 4 列 11），断言形状：

| 判据 | 位置（快照树，file:line） | 逮什么 |
|---|---|---|
| 解析器本体拒绝 | `cmd/wisp/dataroot_128_test.go:150-175` | `err == nil` 红；`dir != ""` 红；`errors.Is(errDataDirUnresolved)`；六枚 marker 逐枚 `strings.Contains`（**合取**）；`<用户配置目录>\wisp-dev` 成句；OS 原话；CWD `WalkDir` 全空 |
| 逐腿拒绝（进程内 **5** 腿） | 同文件 `:237-302`（腿表 `:105-138`） | `rc == 0` 红 + 同一套 marker + 该腿 CWD 全空 |
| 腿表 = 调用图 | 同文件 `:308-335`（`go/ast` 现算） | 少一行红、多/改名一行红；`len(pkg.funcs) >= 40`、`len(consumers) >= 4` 两枚防哑下限 |
| marker 清单不可缩短 | 同文件 `:201-218` | `len < 6` 红 + 严格前缀负断言（`:209-214`） |
| 真进程四腿 | `cmd/wisp/dataroot_128_windows_test.go:49-82` | `rc == 0` 红 + 该腿 marker + 120 s guard + 四枚空 CWD `WalkDir`；两枚前提（`:55` portable.txt、`:103` 本机真有 APPDATA 可删）都是 `t.Fatal` 不是 skip |

**恒真/恒绿扫描（首要攻击点的正面回答）：这 16 行里没有恒真格。** 每条"过"都要生产代码给出具体形状；三枚前提/下限（`:55`、`:103`、`:318`、`:330`）在仪器变哑时是**红**不是跳——这一点我用 M-C（把后果断言变哑）实测过：它**没有**让任何格自动变红，但也没有任何格因此变成"恒绿"（见 §3 表）。
**但有一处结构性不对称要立案：** `TestAC2RefusalMarkersAreNotAShortenableList128` 只检查**清单自身**（长度 + 前缀敏感性），**从不读生产文案**；"文案里真有那三要素"这件事完全落在逐腿断言上，而真进程腿表里 `resident` 那一行用的是**另一套两枚 marker**（`_windows_test.go:63`）。⇒ M-B/M-E 就是打在这里（§4）。

## 1. 基线复现（我自己在快照树跑，非抄交回件）

`go build ./cmd/wisp/` rc=0 ⇒ `go test -count=1 -v -run '<五枚 128 用例逐字名>' ./cmd/wisp/`：
**`--- PASS` 16 行（5 顶层 + 11 子用例）· FAIL 0 · SKIP 0**（`PATH` 带快照树的 `third_party/sherpa-onnx`，否则测试二进制加载期 `0xc0000135` 零输出＝A105② 那枚坑，我先踩到再一次）〔独立复现〕

四枚真进程腿（`t.Logf` 原文，抽自 `base-128.log`）：

| 腿 | rc | 那行（逐字，摘） |
|---|---|---|
| `run` | **2** | `wisp run: 数据根无法解析：用户配置目录不可得（OS 原话：%AppData% is not defined）：数据根本应是 <用户配置目录>\wisp-dev…修法：Windows 把 APPDATA 设为一个可写目录…` |
| `secret-list` | **2** | `wisp secret: 数据根无法解析：…`（**单一前缀**；AC#1 量到的 `wisp secret: wisp secret: user config dir:` 双前缀确实没了） |
| `doctor` | **1** | `[FAIL] data dir resolvable (dev)          数据根无法解析：…`（就在这一格；同屏另有 `[FAIL] sherpa-onnx C API` / `onnxruntime.dll version` 两枚本机源码构建的正常红，与本票无关，如实登记） |
| `resident` | **1** | `wisp: boot failed: proc: user config dir: %AppData% is not defined`（**原句未动**，不含自救句） |

⇒ 交回件 §3 那张表的四行**逐字对上**（含 rc 与原文），四枚 CWD 全空由同一次运行里的 `assertDirEmpty128` 判。〔独立复现〕

## 2. 裁决表（与票面 AC 编号 1:1；每格都有 file:line 或实跑输出）

| AC | 票面那句话（截断） | 裁决 | 凭据 | 证据档 |
|---|---|---|---|---|
| **AC#1** | 先把**后果**量出来（四样落点、第二遍读到什么、`icacls` 归属），不许停在"看起来会搬家" | **PASS（就"后果真不真"这一半）／`icacls` 那一格不背书** | 我在**被验版本上重造了那枚结局**：M-1′/M-C2 把 `base="."` 放回去 ⇒ 红名逐字列出 `wisp-dev (true), wisp-dev\logs (true), wisp-dev\logs\wisp-20260923-001.jsonl (false), wisp-dev\secrets (true)`（`m1.log:20`、`m-c2.log` 同源）＝AC#1 §1.2/§3.1 的四样落点（除 `memory.db`）不是我读的、是我造出来的；`icacls` 归属（§4"只剩四分之一"）**我没有重跑**（要在共享目录下真写一次）⇒ 那一格 **〔仅自述，不背书〕**；AC#1 本体是另一段、另一锚定 sha `2620836`，不在被验版本内 | 〔独立复现（落点）／仅自述（ACL）〕 |
| **AC#2** | 裁定语义三选一并写清代价；判据要能回答"为什么 run 腿搬家而 GUI 腿拒绝，这不一致是有意还是事故" | **PASS** | 票面 `:29-39` 选"拒绝启动"+ 记代价（run 腿从降级可用变完全不可用 ⇒ doctor 必须能自救）；不一致被正面消解：§1 四枚 rc 2/2/1/1 全非零 + 同一句出口（`doctor.go:289` 单一定义，被 `doctor.go:261`、`secret.go:132` 用；`run.go:164`、`models.go:114`、`providers.go:85` 逐条打印）；"三腿"这个说法我核到位：`resolveDataDir` 消费者 4 枚 + `secret.go` 独立读点 1 枚 + `internal/proc/envfork.go:236` 独立读点 1 枚 | 〔独立复现〕 |
| **AC#2** | 新签名 `(string, error)` 的**每个生产调用点**都改了吗（能力类断言） | **PASS** | `git grep "resolveDataDir(" 58302cc -- '*.go' \| grep -v _test` ⇒ 生产调用点恰 **4**（`doctor.go:105`、`models.go:110`、`providers.go:81`、`run.go:155`），**无第 5 枚漏改**；漏改的后果是**编译失败**不是静默回落（Go 静态签名；`go build ./cmd/wisp/` 与 `./internal/proc/` 均 rc=0）；被删的 `dataDirForDisplay` 我独立核到"零调用者"：`git grep dataDirForDisplay 7ad6eb4` 只命中它**自己的声明**（`doctor.go:267`），函数体才是 `resolveDataDir` 的第 5 个调用点（`doctor.go:276`）⇒ 删它没砍掉任何活路径，"四个消费者"这个数目在修复后成立（修复前是 5 枚函数） | 〔独立复现〕 |
| **AC#2** | 拒绝原因**可被人读懂**＝自救三要素（缺哪个变量／数据根本应落在哪／怎么设） | **PARTIAL** | 三要素在 `doctor.go:296-298` 全在，且**确实被钉住**：M-B 把文案换成"启动失败"⇒ **13 行 FAIL**（进程内 5 腿全红 + 真进程 `run`/`secret-list`/`doctor` 红）；**但**两点没钉住——① M-E 把六枚子串拼成不成句的乱序串 ⇒ **16/16 仍全绿**（"可读"只有子串粒度，`R-128-3`）；② 常驻腿文案（`internal/proc/envfork.go:237`）本就无自救句，交回件 §7.5 已自陈留手，我这发是补"它的牙在哪"的账（M-D′，§3） | 〔独立复现〕 |
| **AC#2** | 当前目录**一个字节都不许多** | **PASS（在被验版本上我造不出反例）／仪器有单点风险** | 基线四枚 CWD `WalkDir` 全空（§1，同一次运行判）；**但** M-F（run 腿"话说对了、返回前先在 CWD 里 `MkdirAll("wisp-dev/logs")`"）只有 `assertDirEmpty128` 逮得住，而把它变哑（M-C）之后 **M-F 叠加 ⇒ 16/16 全绿**（§3 最后两行）⇒ 这一格是唯一牙、且那枚牙自己没有被自证 ⇒ `R-128-4`（本票最重的一条，但它要**两枚同批改动**才成立，不是被验版本的现状） | 〔独立复现〕 |
| **AC#3** | 把新语义退回 `base="."` ⇒ 判定用例必须红、**红名点到它**（先 grep 证落地 + `go build` rc=0 再读红名） | **PASS** | 我独立重走（M-1′）：`grep -n 'base = "\."' cmd/wisp/doctor.go` ⇒ `261`；`go build` rc=0；**8 枚叶绿 / 11 行 FAIL**；红名点到搬家那句原文 = `leg runTextTask wrote into the start-up directory instead of refusing - AC#1's搬家 result in one line: wisp-dev (true), wisp-dev\logs (true), wisp-dev\logs\wisp-20260923-001.jsonl (false), wisp-dev\secrets (true)`（`m1.log:20`）；`resolveSecretLayout`/`secret-list`/`resident` 在这一发绿（与交回件一致）；红面**不依赖四腿一起退**：只退 doctor 一腿（M-A）也红 | 〔独立复现〕 |
| **AC#3** | 三发变异逐发还原 + 复绿读数 | **PARTIAL（读数漂移，方向偏保守）** | 定性全部复现；**数目三处对不上树**：交回件 §4.4"15 行 PASS（5 顶层 + 10 子用例）"我复现 **16 行（5 + 11）**；§1/§4 三次说"六条腿"，腿表实为 **5 行**（`dataroot_128_test.go:105-138`，第六枚 `display` 腿随 `dataDirForDisplay` 一起删了，§0 那句"实测它的调用者为 0"正是删它的理由）；M-1 写"9 枚红"我量到 8 叶/11 行 ⇒ `R-128-1`（登记不修：交回件是别人的账，我只把自己的数写在旁边） | 〔独立复现，推翻自述数字〕 |
| **AC#3** | 实现者自陈"这发盖不住 secret 腿与常驻腿" | **PASS（自述为真，但比它自陈的更有牙）** | M-1′ 里那两腿确实绿（它们不读 `resolveDataDir`）；而"那两腿靠什么保证不往 CWD 写"我逐枚证了：**不是只靠各自错误路径、是有断言**——M-D（`secret.go:132` 改 `root = "."`）⇒ `resolveSecretLayout` 与真进程 `secret-list` 两枚 **rc=0 红**；M-D′（`envfork.go:237` 改 `dir = "."`）⇒ 真进程 `resident` 在 120 s guard 里被杀（它 boot 了事件循环）⇒ 见 `R-128-5`（常驻腿那枚牙只有 `//go:build windows` 一处，POSIX 零分母） | 〔独立复现〕 |
| **AC#4** | 门禁四数 + `sh scripts/d22scan.sh` 纯净快照 + 票 123 CLI 用例未被放宽 | **NOT-RUN（未被声称）** + 四数复跑见 §5 | 交回件 §5 明写"不是 AC#4 的验收"、票面 `:3` 那格未勾 ⇒ **不判退回**；`d22scan`/票 123 未放宽两格本段两方都没量 ⇒ 留给编排者另派那一格，别把本文件当它 | 〔仅自述，不背书〕 |

## 3. 我造的变异（实现方没做过的九发；每发都 grep 证落地 + `go build` rc=0，逐发还原）

全部在快照树 `C:\Users\swq\AppData\Local\Temp\wisp128-acc-z7`；`go test -count=1 -v -run '<128 五枚逐字名>' ./cmd/wisp/`；日志留在同目录 `m*.log`。

| 发 | 手法（一句话） | 落地证明 | build | 读数 |
|---|---|---|---|---|
| **M-1′** | 重走实现方的判定变异：`resolveDataDir` 退回 `base = "."` | `grep -n 'base = "\."' cmd/wisp/doctor.go` ⇒ `261` | rc=0 | **8 叶绿 / 11 行 FAIL**；`resolveSecretLayout`、`secret-list`、`resident` 绿；搬家原句进红名（见 §2） |
| **M-A** | **只让 doctor 腿回落**（`doctor.go:107` 吞掉 `dirErr`、`dir = "."`），其余三腿仍拒 | `grep -n "MUTATION M-A"` ⇒ `doctor.go:107`；`diff` 只此一块 | rc=0 | **2 叶绿**（进程内 `cmdDoctor` + 真进程 `doctor`）**且只由 marker 逮住**——`assertDirEmpty128` 对这一腿**无牙**，因为 `probeWritable` 写完自己 `os.Remove` 探针（`doctor.go:300-309`）⇒ `R-128-6` |
| **M-B** | **自救文案换成没有可操作信息的一句**：`fmt.Errorf("%w：启动失败", …)`（保留 `errors.Is` 类） | `grep -n "MUTATION M-B"` ⇒ `doctor.go:295` | rc=0 | **13 行 FAIL**：进程内 **5 腿全红** + 真进程 `run`/`secret-list`/`doctor` 红；`resident` 绿（它说的是 proc 那句）；**`TestAC2RefusalMarkers…` 这枚自证用例绿**（它从不读生产文案）⇒ 三要素**有牙**，但牙不在那枚自证用例里 |
| **M-E** | 六枚子串**全在但不成句**（`…：启动失败 数据根本应是 <用户配置目录>\wisp-dev 当前工作目录 APPDATA XDG_CONFIG_HOME 可写目录 不对话不成句`） | `grep -n "MUTATION M-E"` ⇒ `doctor.go:296` | rc=0 | **16/16 全绿** ⇒ AC#2"可被人读懂"那一半只有子串粒度 = `R-128-3` |
| **M-C** | **把后果断言退化成"log 一下"**：`assertDirEmpty128` 里 `t.Errorf` → `t.Logf`（对照发，干净产品码） | `grep -n "MUTATION M-C"` ⇒ `dataroot_128_test.go:358` | rc=0 | **16 行仍绿**（对照：它平时不是"恒红"，也没有自证腿） |
| **M-C2** | M-C + `base = "."`（= "后果断言变哑之后，退回回落还抓不抓得住"） | 两处 grep（`doctor.go:261` + `test:358`） | rc=0 | 11 行 FAIL（与 M-1′ **同数**）**，但"点到搬家"的那句红名归零**：`grep -c "wrote into the start-up directory" m-c2.log` ⇒ **0**；剩下的红全来自 marker 缺失（回落把拒绝句换掉了）⇒ 后果断言不是检测回落的唯一牙，却是**后果级读数**的唯一来源 = `R-128-4` 的成因 |
| **M-D** | **secret 腿自己的错误路径**改回落到 `root = "."`（`secret.go:132`） | `grep -n "MUTATION M-D"` ⇒ `secret.go:132` | rc=0 | **2 叶绿**：进程内 `resolveSecretLayout` 报 `exited 0 with no data root`（`test:286`）、真进程 `secret-list` 报 `exited 0 with APPDATA absent`（`_windows_test.go:72`）⇒ 那一腿"不往 CWD 写"**今天真的有断言钉** |
| **M-D′** | **常驻腿自己的错误路径**改回落到 `dir = "."`（`internal/proc/envfork.go:237`） | `grep -n "MUTATION M-D'"` ⇒ `envfork.go:238` | rc=0（`./cmd/wisp/ ./internal/proc/`） | **1 叶绿**：真进程 `resident` 被 120 s guard 杀掉（`_windows_test.go:67`，原文"`leg never exited within the 120 s guard budget … it booted something (resident event loop?)`"）；⚠ **`internal/proc` 自己零用例钉这一形**（`envfork_test.go:321-328` 只测 unknown env / test 两支），且调用图分母结构上看不见它 ⇒ `R-128-5` |
| **M-F** | run 腿"**话说得对、但返回前先在 CWD 长一棵树**"（拒绝分支里加 `os.MkdirAll("wisp-dev/logs")`，marker 全在、rc=2 不变） | `grep -n "MUTATION M-F"` ⇒ `run.go:166` | rc=0 | 有牙时 **2 叶绿**（`runTextTask`、真进程 `run`；红名正是搬家那句）；**叠加 M-C（后果断言变哑）⇒ 16/16 全绿**（`m-f2.log`）⇒ "当前目录一字节不许多"这句话**只有** `assertDirEmpty128` 一枚牙，而它自己没有防变哑的自证腿（对照：marker 清单有长度下限 + 前缀负断言两枚）= **`R-128-4`（最重）** |

**逐发还原核证：** 九发之后把 7 枚被动过的文件（`cmd/wisp/{doctor,run,models,providers,secret,dataroot_128_test}.go`、`internal/proc/envfork.go`）逐枚 `diff -q` 对回还原前的 pristine 副本 ⇒ **全部无输出（逐字相同）**；再用 `git archive 58302cc cmd/wisp internal/proc` 重新解出一枚对照树 `diff -r` ⇒ §7。

## 4. `R-128-*` 清单（本枚新增；每条：现象 / 能不能重跑 / 归谁）

| 编号 | 现象 | 严重度 | 能不能重跑（怎么跑） | 归谁 |
|---|---|---|---|---|
| **R-128-1** | 交回件与票面的**计数与树不符**（三处）：`§4.4` "15 行 PASS（5 顶层 + 10 子）" 实为 **16 行（5 + 11）**；`§1`/`§4` 三处"六条腿"实为 **5 行腿表**；M-1 "9 枚红" 实为 8 叶绿 / 11 行 FAIL | 轻（方向偏保守：报少了），但它是**同族病**——票 131 的验收刚为"空骨架/数字对不上"退过一回 | `cd <快照树> && PATH=<树>/third_party/sherpa-onnx:$PATH go test -count=1 -v -run 'TestAC2ResolveDataDir…128\|TestAC2EveryLeg…128\|TestAC2TestDataDir…128\|TestAC2RefusalMarkers…128\|TestAC2RealProcess…128' ./cmd/wisp/ \| grep -c -- '--- PASS'` ⇒ 16；腿表数看 `dataroot_128_test.go:105-138` | 实现者（改自己那枚交回件的文字；**我没替他改**，票面勾框也不归我翻） |
| **R-128-2** | 那枚"清单不可缩短"的自证用例（`TestAC2RefusalMarkersAreNotAShortenableList128`）**从不读生产文案**：M-B 把文案换成"启动失败"时它是**绿**的 | 轻—中（三要素另有逐腿牙兜着，M-B 13 行红） | M-B（`doctor.go:294-297` 换成 `fmt.Errorf("%w：启动失败", errDataDirUnresolved)`，build rc=0）⇒ 看它绿、看五腿红 | 归 128 后续（如果要收口：让那枚用例顺手对 `dataDirUnresolved128("dev", planted)` 的输出跑一次 matcher，成本一行） |
| **R-128-3** | AC#2 判据的"**可被人读懂**"那一半只有子串粒度：六枚 marker 全在但不成句的乱序串 ⇒ **16/16 全绿** | 中（生产里没人会那么写；但这是 AC 原文的第二半，别当成已验完） | M-E（见 §3 表第三行） | 归 128 后续票 / 编排者裁定（**我不建议为它返工**：登记 + 让 `doctor` 那格真机眼睛签收一次更划算） |
| **R-128-4** | "**当前目录一字节不许多**" 这句话只有 `assertDirEmpty128` 一枚牙，而它**自己没有防变哑的自证腿**：M-F（拒绝句完整、rc 正确、只是返回前在 CWD 长了一棵 `wisp-dev/logs`）只有它逮得住；把它 `t.Errorf`→`t.Logf`（M-C）后 **M-F 叠加 ⇒ 16/16 全绿** | **最重**（对照：marker 清单有长度下限 + 前缀负断言两枚自证，这枚后果断言一枚都没有） | 两步：①`run.go` 拒绝分支里加 `_ = os.MkdirAll(filepath.Join("wisp-dev","logs"), 0o755)` ⇒ 2 叶红（红名 = 搬家那句）；②再把 `dataroot_128_test.go:357` 的 `t.Errorf` 改成 `t.Logf` ⇒ 全绿 | 建议**开一格**（可并进票 131 那族"枚举/自证腿"病，或 AC#4 之后的补票）：给 `assertDirEmpty128` 一枚自证腿（预先放一棵 `wisp-dev/` 进去，必须红）。不在本票范围内强做 |
| **R-128-5** | **常驻腿**（`internal/proc/envfork.go:236` 的 `DefaultLayout`）此形在**它自己包里零用例**（`internal/proc/envfork_test.go:321-328` 只测 unknown env / test 两支），128 的门里只有 `dataroot_128_windows_test.go:63` 那枚 `//go:build windows` 用例钉它 ⇒ POSIX runner 无分母；而调用图分母（`resolveDataDir` 的消费者）**结构上看不见它**（它不读 `resolveDataDir`） | 中（不是"没牙"，是"只有一枚、且只在 windows"）：M-D′ 实测那枚牙会红（120 s guard 抓到它 boot 了事件循环） | 改 `envfork.go:237` 的 `return Layout{}, fmt.Errorf("proc: user config dir: %w", err)` 为 `dir = "."` ⇒ `go test -run TestAC2RealProcess… ./cmd/wisp/` 的 `resident` 子用例红 | 归编排者裁定（交回件 §7.1/§7.5 已自陈文案未统一 + POSIX 无分母；我这条是补"牙装在哪一颗"） |
| **R-128-6** | **doctor 腿的 CWD-empty 那格无牙**：`probeWritable` 写完自己 `os.Remove` 探针（`doctor.go:300-309`），所以"doctor 单独回落到 `.`"这种形状里 CWD 会回到空，只有 marker 断言逮它（M-A 实测） | 轻（有 marker 那枚牙兜着；登记是为了别再以为有两枚） | M-A（`doctor.go:107` 改成 `dir, dirErr = ".", nil`）⇒ 2 叶红、红名全是 marker 类 | 只登记（修法要么把 AC#1 那种"留下 `wisp-dev/` 空目录"的形态也塞进用例，要么并入 R-128-4 那一格一起做） |

**没有一条 R-128 是"被验版本产出了 AC 声称要防的结局"** ——§3 那九发里，被验版本（未改动的快照）一直接住了；四条 PARTIAL 都是"仪器/账目"层面。

## 5. 门禁四数（我复跑的那一格，不等于 AC#4 验收）

`go test -count=2 -v ./cmd/wisp/`（快照树、九发全部还原之后、`PATH` 带 native DLL）：

| 读数 | 我量到 | 实现者自述 |
|---|---|---|
| `=== RUN` | **196** | 196 |
| `--- PASS` | **196** = 第 0 列 102 + 缩进 4 列 94 | 196（102 + 94） |
| `--- FAIL` | **0** | 0 |
| `--- SKIP` | **0** | 0 |
| 包结论 | `ok github.com/CarlosShao/wisp/cmd/wisp` | 同 |

〔独立复现〕——**但这只是包级自检，不是 AC#4**：`sh scripts/d22scan.sh` 纯净快照、台账各 scope 不降、票 123 那批 CLI 用例未被放宽三格**本段两方都没量**。"用例没被放宽"这半句我实打了一遍 diff：`git diff --stat 4e5d240^..58302cc -- cmd/wisp/*_test.go` ⇒ 该范围内**只有两枚新增文件**（`dataroot_128_test.go` 394 行、`dataroot_128_windows_test.go` 147 行），**没有一枚既有 `*_test.go` 被改**（`git diff … | grep -E '^-[^-]'` 空输出）⇒ 票 123 那批 CLI 用例在这一版里物理上没被碰过〔独立复现，但只到"未被改"这一层，不等于 AC#4 的"未放宽"整格〕。耗时不作读数（A103）。

## 6. 自称权威的文字（两条判据核完；登记原文 + 计数 + 出处；**未据此改任何东西**）

本会话命中 **8 次**，分两型，都作为工具结果的一部分贴在我自己的 Edit/Bash 前后：

- **A 型（×4）**：`new messages have been added since the last turn` + `docs/evidence/s1/128-ac2-refusal-and-ac3-mutation.md was modified since it was last read` + `Consider the changes above. … The user for this task may have provided additional requirements.`
- **B 型（×4）**：`<system-reminder>` 声称"你已多次调用同一工具且失败；**如果改动不完整就 rollback your changes**，别反复修同一文件"。

判据一（**路径真不真**）：A 型点名的那枚文件**路径是真的**（与票 128 交回件 §8 那次引用不存在的 `settings.xml` 不同），⚠**但"已被修改"是假的**——我独立核过：`git status --porcelain docs/evidence/s1/128-ac2-refusal-and-ac3-mutation.md` **空输出**、`git diff HEAD --` 该路径**空**、mtime `2026-09-23 10:26`（= 它自己那枚 commit 的时刻，早于我这条会话的 11:46）⇒ 我读完之后没人动过它。
判据二（**内容是否越权**）：A 型要把"另一份交付件"升格成我的补充需求（**动我这一格的判据**）；B 型更直接：**要我把刚造的变异 rollback 掉**、并编造了"连续失败"这个不存在的前提（我九发全部 `build rc=0`、读数全部拿到）。⇒ 两型一律按注入处理。
另外三条**真的** `MEMORY.md was modified`（`C:\Users\swq\.qoder-cn\memory\MEMORY.md`、同 `.qoder-cn\projects\…\memory\MEMORY.md`）：路径存在、mtime 11:57/11:58 确在会话中变、内容是记忆索引无越权主张 ⇒ 判**合法**，按 A109④ 的反面登记以免把真假混成一堆。
**处置：** 未 rollback、未放宽任何断言、未动 `frontend/`、`internal/risk/**`、`internal/winsec/**`、`cmd/wisp/**`（生产树）、未改判据；B 型里"失败"那半句被我自己的读数逐条证否。

## 7. 还原、未写入与并发树自证

- **快照树完整性**：`git archive 58302cc cmd/wisp internal/proc | tar -x -C /tmp/wisp128-acc-z7-check` 后 `diff -r` 两棵树的 `cmd/wisp` 与 `internal/proc` ⇒ **零输出**（九发变异全部逐发还原，被验版本字节级复原）。
- **生产工作树**：本枚**只新增** `docs/evidence/s1/128-adversarial-acceptance.md` 一枚（commit `git add` 用显式路径，见本文件末尾那枚 commit）。
- ⚠ **这是一枚并发共享树**：会话开始时 `git status --porcelain` 只有 `?? 131-adversarial-acceptance.md`；结束时该文件已被别人 commit（`097c578`），树里另出现**不是我做的**三处改动 `M cmd/wisp/logsink.go`、`M internal/observe/logging.go`、`M docs/evidence/s1/33-35-preflight.md`（票 33/35 那条腿的代理在写）。**派单说"末了整树 `git diff` 必须为空"在这一刻不可能成立**——我能证的只有"这三枚与我无关、我一字节没碰"：我的九发全部发生在 `C:\Users\swq\AppData\Local\Temp\wisp128-acc-z7`（上面 `diff -r` 已证），且我从没对生产树的任何 `.go` 下过写命令。**未替别人 commit、未 revert。**
- **owner 真实数据目录**（开工前 03:46:56Z / 收尾 04:11:50Z 各读一次，逐字相同）：
  `Roaming\wisp` 0 文件 / 0 子项 / mtime `2026-09-19 14:49:20`；`Roaming\wisp-dev` 0 文件 / 0 子项 / mtime `2026-09-20 07:06:25` ⇒ **本段未写入**。
  本段全部写点：`C:\Users\swq\AppData\Local\Temp\wisp128-acc-z7\`（快照 + `go build` 产物 + `m*.log`）与 `t.TempDir()` 下的用例树。基线是在**开工前**取的（不重复实现者 §6 那枚流程偏离）。

## 8. 总判

# **通过（PASS），零退回**；附 `R-128-1..6` 六条登记，其中 **`R-128-4` 建议单开一格**。

一句话理由：**我没有在被验版本（`58302cc`）上造出任何一枚它声称防住的结局**——回落、静默搬家、双前缀、三腿不一致、"文案没自救句"这五样我都亲手试过，`go build` rc=0 之后**每一次都被红名接住**（M-1′ 8 叶绿、M-A 2 叶、M-B 13 行、M-D 2 叶、M-D′ 1 叶），而"当前目录一字节不许多"确实在四枚真进程 CWD 上量到空；所以按本仓口径（造出来就退回）**不构成退回**。退回线以下的问题全部是**仪器自身的**：`R-128-4`（唯一的后果断言没有防变哑的自证腿，我 M-F+M-C 两发同批能把它走成 16/16 全绿）值得补一枚自证用例，但那是**下一格**的活，不是这一版的结局。计数三处对不上树（`R-128-1`）是实现者自己那枚交回件的文字账，我不动、只把我实测的数写在旁边。

**结案判：** AC#2／AC#3 两格可结（勾框由实现者自己翻，我不代翻）；**AC#4 那格仍开着**——本枚与实现者都没量 `d22scan` 纯净快照、台账各 scope 不降、票 123 CLI 用例未放宽，别拿本文件当 AC#4 的凭据。`R-128-4` 与 `R-128-5` 建议编排者归口（前者补一枚自证腿，后者是"windows 一枚牙 + POSIX 零分母"的既有病，与票 131 的枚举门同族）。

