# 300-v4 · 残余两问（ⓓ R2／ⓔ R1）＋ `AC#4` 三裁（ⓕ）

锚＝`34e1962b`。台面＝仓外导出树。件路径同 `10-ac6-adjudication.md`。

---

## ⓓ verdict：**⛔ 算 `AC#6` 未闭合——本格判据已满足；R2 是本格的**射程上界**，而"把生产那行变可观察"这一步要另一次批准**

`AC#6` 原文判据逐字只有两句：①`把常量换成任一其它值（3→1 与 3→0 两形）指名用例必须红`；②`红句具名指向那枚常量而⛔ 指向 tag 的等值比较`。**两句都⛔ 要求观察侧抵达 `wasapi_windows.go:378`。**
⇒ 把 R2 读成"本格没闭合"＝**给本格追加一枚它⛔ 承诺过的判据**，那才是真跑歪。裁语 ④ 的标题逐字写着`残余（具名，⛔ 塞进 AC#6）`——这一格照办。

### 我这一发的客观尺（覆盖尺，件 `logs/37-cov-drain-region.txt`、`logs/39-cov-func-named.txt`）

```
wasapi_windows.go:377.2,380.6 4 0        <- 装着 :378 floating := s.format.tag == waveFormatFloat 的那个块，命中数 0
wasapi_windows.go:376:  Drain            0.0%
wasapi_windows.go:409:  convertPacket    0.0%
wasapi_windows.go:200:  parseWaveFormat  81.8%
```

⇒ "**生产那一行今天零执行**"再⛔ 是读码推论，是一把覆盖尺的读数（`rc=0`，整包含 `-coverprofile`，件 `logs/34-cov-run.txt`）。

### 为什么除 C5 之外无解（我逐枚查过消费者，⛔ 假设）

- `Drain()` 体第一跳就是 `comCall(s.capture, 5, …)`（`:382`，`GetNextPacketSize`）⇒ 需要**活流的 COM 指针**；票面 `:31` 逐字 `⛔ 真设备／真麦克风那一发归编排者（会占麦、机主在场才做）；⛔ 派给腿` ⇒ 腿这条路封。
- `convertPacket(data, frames, ch, floating, silent)` 收的是 **`floating bool` 参数**（`:400` 是唯一生产调用点，`300-a2` 的 `F3` 与我这发 `grep -rn convertPacket` 同读数）⇒ 直接测 `convertPacket` 钉的是**参数语义**、⛔ 那枚常量。
- 剩下唯一一形＝把 `s.format.tag == waveFormatFloat` 抽成包内谓词（例：`func (f waveFormat) isFloating() bool`）让 `:378` 与测试**共用同一个谓词** ⇒ **形 C5＝动产码语义面**。

### ★具名写给编排者、⛔ 我替你决定

**若有人主张"必须让生产那行可观察才算闭"，那一步＝C5＝改 `internal/audio/wasapi_windows.go` 的产码语义面，本票⛔ 授权它，必须另一次批准（新格／新票），且它一旦落地就把 `AC#6` 的判据换成另一把尺。**
本腿**⛔ 做这一主张**，也⛔ 顺手起草它。建议处置＝R2 作为**已具名、已量到的残余**留在台账（配我这发 `Drain 0.0%` 那三行读数当凭据），要推就单开一票；票面 `AC#6` 那格**照现有判据翻勾⛔ 欠任何东西**。
（另留一枚便宜的替代给未来那格参考：`真机 GetMixFormat` 那发编排者自己交过了（`AC#3` 的 hexdump，四端点全 extensible＋float），它证的是"真实数据会走到那一行"，⛔ 证"那一行在测试里被执行过"——两件事⛔ 混。）

---

## ⓔ verdict：**覆盖尺我这一发就补了（⛔ 再欠）；R1 坐实＝真缺口；它⛔ 属于 `AC#6`，该另立一格，并且我这发顺带量到一枚**新的同族零尺**（R1b）**

1. **要不要现在补一把覆盖尺＝要，且我已补**（本腿⛔ 停在词面尺上）。尺＝`go test ./internal/audio/ -count=1 -coverprofile=<仓外>`，`rc=0`；`parseWav` **70.0%**（件 `logs/39-cov-func-named.txt`）。
2. **读数（件 `logs/36-cov-injector-tag-branches.txt`，行尾那一列＝命中数）**：

```
wavinjector.go:193.5,193.18 1 0     <- fmtTag == 0xFFFE 的展开体（real tag is SubFormat[0:2]）
wavinjector.go:194.6,195.1 1 0
wavinjector.go:196.5,196.65 1 0     <- 从 body+24 读真 tag 那一行
wavinjector.go:219.3,221.26 3 0     <- case fmtTag == 3 && bits == 32 的整个 float32 支
wavinjector.go:222.4,223.1 1 0
wavinjector.go:224.3,224.43 1 0
wavinjector.go:226.3,226.111 1 1    <- 唯一被走到的：unsupported wav format 拒绝路径
```

⇒ 编排者 17:5x 那句"两把词面尺（`float|FFFE|extensib` 只命中 `sineI16At` 的 `float64` 生成器；`float32|Float32` 0 命中）"**⛔ 覆盖尺"**——今天升级为覆盖尺：**两支俱为 0 命中**，测试面**从未**给 injector 喂过 float32 载荷而⛔ extensible 容器。`R1` **成立**。
3. **补给谁＝⛔ 塞进 `AC#6`，另立一格**。判据形状已经⛔ 是"常量⇔权威"，而是"**夹具⇔生产分支**"：要一枚新的**测试侧**件（`writeWav` 能写 `IEEE_FLOAT` 32-bit／能写 `0xFFFE`＋`cbSize=22` 容器），⛔ 动产码、⛔ 需真设备 ⇒ 落点＝`internal/audio` 的 injector 测试面，射程窄、可测性现成（`wavinjector_test.go:18 writeWav` 就在那儿）。⇒ 派给**下一枚落地腿**⛔ 派给本腿，也⛔ 挂到 `AC#6` 的翻勾条件上（那格⛔ 欠它）。
4. ★**R1b（本格顺带量到的新洞，建议具名进台账）**：`wavinjector.go` 把同源两枚权威值写成**行内字面量**（`:192 0xFFFE`、`:218 3`），⛔ 引用 `waveFormatExt`/`waveFormatFloat`（平台中立件⛔ 能引 `//go:build windows` 里的常量）。我这发把它从"推断"变"实测"（件 `logs/70-50-*.txt`、`logs/70-51-*.txt`；`cp`→`sed`→跑→`cmp`＝IDENTICAL）：

| 突变 | rc | 名册尺（`^--- FAIL`＋`^    --- FAIL`） |
|---|---|---|
| `:218 fmtTag == 3` → `== 2` | **0** | **0 行＝整包照绿** |
| `:192 fmtTag == 0xFFFE` → `== 0xFFFD` | **0** | **0 行＝整包照绿** |

⇒ `AC#6` 治的那枚病（**改值无人喊红**）在 `wavinjector.go` 里**原样还有一枚**，而 `AC#6` 原文与裁语 ③ 都⛔ 覆盖它。**这是新格，⛔ 是本格的漏**；把它写成票时可直接引本件那两行读数当"改前必红"凭据的对偶（今天⛔ 红＝缺口的存在证明）。
⚠ 一处必须说破：`R1b` 与 `R1` 是同一次落地**⛔ 同一次裁决**——补 injector 的 float32/extensible 夹具会**同时**买到 R1 与 R1b 的牙，因为那两支一旦被执行，行内 `3`/`0xFFFE` 就被断言钉住。⇒ 建议**一票两洞**，别开两票。

---

## ⓕ `AC#4` 三裁

先把"票面 `AC#4` 逐字要求 ⇔ 谁交了什么"摊平（票面 `:25` 全文；台面＝仓外导出树＝HEAD blob 同字节，工作树 `git status --porcelain internal/audio`＝**0 行**）：

| `AC#4` 那一串判据 | 落地腿 `300-r2` | 编排者复跑（19:2x） | **本腿现量** |
|---|---|---|---|
| `GOFLAGS= go build ./...` rc=0 | ⛔ 跑（派单漏列，自己具名⛔ 跑） | **rc=0** | **rc=0**（件 `logs/33-build-all.txt`） |
| `sh scripts/d22scan.sh` rc=0 | 交 | **rc=0**（ban#8 `internal/`=526 实扫） | **rc=0**（件 `logs/32-d22scan-tail.txt`） |
| 格式两把尺并排（工作树＋仓外 blob）新增 0 | 交 | 名册 0 枚 | **射程 `internal/audio`，导出树（＝blob 内容）`gofmt -l` count=0**；工作树与 blob **同字节**（`internal/audio` 脏行数 0）⇒ 两把尺在此**collapse 成同一读数 0**（件 `logs/31-gofmt-scope.txt`）。全仓那 5 枚既有残留⛔ 我射程，⛔ 动 |
| `PATH=…sherpa… go test ./internal/audio/ ./cmd/wisp/ -count=1 -v` 改前改后**各 ≥2 发**取红名交集＝新增红 0 | 门③**只跑 `./internal/audio/`、各 1 发**；`cmd/wisp` ⛔ 跑 | 整包两发红名册逐字同形、新增红 0 | **`./internal/audio/` 半边我补齐成对两发**：改后（HEAD）2 发＋改前（**HEAD 树拆掉那枚新用例**）2 发，`comm` 双向皆**空**⇒ **新增红 0 枚**；四发全 `rc=0`、`ok internal/audio 17.014/16.601/16.029/16.323s`、`--- SKIP` 两侧同为 1 枚同名、顶层 PASS **42（含新用例）⇔ 41（拆掉）**，尺名＝`^--- PASS` 顶层（件 `logs/70-*.txt`、`logs/44-*.txt`、`logs/45-*.txt`） |
| 每枚门禁件自落一行 `rc=N` | 交（`r2/logs/` **14** 枚里 **7** 枚含 `rc=`；其余是 md/名册件，⛔ 我这把尺判它＝那是腿自家件的形状，本腿⛔ 判它违规） | — | 本腿件逐把带 `rc=`（见 `90-hygiene.md`） |
| 名册只含 `internal/audio/**`＋`probes/300/**` | 自审 | 七笔并集＝越界 **0** | **落点笔 `9a442923` 名册逐字 4 行**：`.scratch/wisp/probes/300/r2/10-shape.md`、`r2/logs/gate1-vet.txt`、`r2/logs/gate2-targeted-pristine.txt`、`internal/audio/wave_format_float_300_windows_test.go` ⇒ **在允许面内** |

### ⓕ1 派单漏列那一枚＝**⛔ 算腿的欠账，记编排者**（本腿照既有定式判，且这格已双重销账）

派单（17:5x 裁语 ③④＋⑤）交给落地腿的清单里**⛔ 有** `go build ./...`；腿自己**具名报了这一枚⛔ 跑并问了一句**（票面 `:281`），⛔ 假设、⛔ 谎报、⛔ 默默缩射程——那是腿的正确行为。
⇒ 按编排者既有定式（**派单里就禁掉／漏列的资源⛔ 算腿的欠账，记编排者**）判：**这一笔⛔ 挂在 `300-r2` 名下**。且它今天已由**两枚独立 `rc=0`**（编排者 19:2x＋本腿）双重闭合，**欠账清**；建议翻勾时把"腿⛔ 跑"写成〔归编排者派单，已补〕，⛔ 写成〔腿欠〕。

### ⓕ2 `./cmd/wisp/` 那一半＝**现在⛔ 补在本腿，也⛔ 补在 `300-r2`；交给握着 `cmd/wisp` 车道的人，与 `AC#5` 的 CI 色同批取一次**

四条具名理由：
1. **因果面是零**：新增的是一枚 `package audio` 的 `_test.go`。Go 的测试文件⛔ 进任何别的包的导入图 ⇒ `cmd/wisp` 的编译单元与用例集合**逐字节⛔ 变**。我这把名册尺给了同一结论的另一面：新用例对 `./internal/audio/` 的净效应＝**顶层 PASS 41→42、FAIL 0→0、SKIP 1→1**；它对 `cmd/wisp` 的净效应是**结构上的 0**。⇒ 为它付 4 发整包（外加真窗抖动）**买⛔ 到判别力**。
2. **车道冲突**：`cmd/wisp` 整包是 `303-v1` 的**独占面**，我的派单逐字写着 `⛔ 跑 cmd/wisp 整包`。我⛔ 越。
3. **台面代价实名**：那一半要 `PATH=/d/work/…/third_party/sherpa-onnx:/d/work/…/build`（写成 `D:/…` ⇒ `0xc0000135` 而**零** `--- FAIL`＝用例根本没跑，我起手就按 shell 形路径配我的树，未踩）＋已知的本机抖动红（`TestAC4FocusReturnToPriorWindowGap33r5` 真窗焦点类 13 发 5 绿／8 红、`TestPanelHostRealWindowHopAndLifecycle`、负载敏感的 `TestResolvePerCallBudget`、缺 `frontend/dist` 时 `TestAC13…` 走具名 SKIP）。在这种基线上取"改前改后各 ≥2 发红名**交集**"是**可行的但⛔ 便宜的**，而且它⛔ 该由一枚只裁 audio 常量的腿顺手做。
4. **它有一枚免费的正替代**：`AC#5` 那一半本来就欠"新用例进 CI 之后的颜色"，而 CI 的 `windows` 档**同时**求值 `./internal/audio/` 与 `cmd/wisp`（`portable-tests.sh:253` 逐字含 `./internal/audio/`，编排者已现量）。⇒ **同一次推送那一发**同时销 `AC#4` 的 `cmd/wisp` 半边与 `AC#5` 那半格，比本地四发更硬。

⇒ **裁**：`cmd/wisp` 那半＝**归编排者的同批推送那一发取**，⛔ 派回 `300-r2`、⛔ 派给本腿；台账里请按〔已定归属：编排者／同批 CI〕记，⛔ 记成腿的欠账（同 ⓕ1 定式）。若编排者坚持要本地凭据，唯一正确形状＝**等 `303-v1` 交完并复跑后，由那一腿的台面顺带跑一次**（一枚台面、一次 `PATH` 铺设、两格共用）。

### ⓕ3 `AC#4` 终裁＝**附条件成立**

- **成立的部分**（凭据已齐，且多数有**第二人独立复跑**）：build rc=0／d22scan rc=0／格式两把（本射程 collapse 为 0）／`./internal/audio/` 成对两发**改前改后**红名交集＝新增红 0 枚／写面名册在允许面内／每枚门禁件带 `rc=`。
- **条件一枚（唯一）**：`AC#4` 原文那一串里的 **`./cmd/wisp/` 半边红名交集**尚**未交任何人的读数**；归属＝编排者那一推送同批（ⓕ2），⛔ 腿欠。
- 措辞建议（翻勾时）：**"`AC#4` 附条件成立——audio 半边三人口径齐（腿／编排者／`300-v4`），`cmd/wisp` 半边由 `AC#5` 那一发同批销；条件⛔ 落在落地腿"**。
- ⛔ 把 `AC#4` 判成"不成立"：那会把一枚**派单漏列**的判据算成实现腿的账，正撞本仓最忌的"跑歪模式 #1"反面（契约／账目⛔ 由裁决者单方面挪）。也⛔ 判成"成立"无脚注：那会抹掉 ⓕ2 那半枚仍欠的读数。

（`AC#5` 本腿⛔ 裁：它问的是**CI 那一发颜色**，只有推送能给，票面逐字欠着；我这发只证明了我这一侧⛔ 产生任何新 SKIP。）
