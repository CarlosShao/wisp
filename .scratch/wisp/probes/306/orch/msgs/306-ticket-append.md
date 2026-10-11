
## 编排者收 `306-v1`（2026-10-11 08:2x，五笔 `ec87076f`→`a458114b`→`ffc41214`→`56634d23`→`1eaa91f6`，件 `probes/306/v1/**`＝五份正文＋`logs/**`；⛔ 本节目＝我自己复跑的读数）⇒ **五格全翻**＋★**两枚 append-only 更正落在我账上**＋★**裁决表 1:1 已建**（`docs/evidence/s1/306-inline-literals-v1.md`）

### 0. 台面与越界（⛔ 抄腿自陈）

- 起手锚＝`6284a489`（派单给的号），腿自己**三次重锚** `e4740e35`→`05b37976`→`17dd4a6f`，我复核＝三枚之间 `internal/audio` **0 字节**（`git diff --quiet e4740e35 17dd4a6f -- internal/audio` rc=0）⇒ 漂移⛔ 影响本票任何一行号与判据。**定式再钉一遍：派单里的 HEAD 号只算派发时刻快照。**
- 逐笔越界尺我在 **13 笔**（`306-r1` 八笔＋`306-v1` 五笔）上自己重跑 `git show --name-only --format=%H`：每笔**写面外＝0**、`.out`＝**0**；并集尺＝`5480434f..HEAD` 里碰 `internal/audio` 的**只有那 2 枚**（新测试文件＋`AC#3` 那枚注释文件）。件＝`probes/306/orch/logs/orch-03-roster-306-commits.txt`。
- ⚠ **`306-v1` 第 2 笔 `a458114b` 把两枚 HEAD blob 副本命名 `.go` 落进了 `logs/`**（那一笔 新go＝2），第 5 笔 `1eaa91f6` 走 `git mv` 改成 `.txt`、内容 `cmp rc=0` 两把都过 ⇒ **把 `.go` 的分母还给别家的尺**（`gofumpt`／`d22scan` 的射程按 `.go` 数），⛔ 删、⛔ 追认成缺陷。这一枚我复核过：现 HEAD 上 `git ls-files .scratch/wisp/probes/306/v1 | grep -c '\.go$'`＝**0**。

### 1. 承重读数我自己在盘上现跑（⛔ 抄腿；每条带尺；台面＝共享工作树 HEAD `1eaa91f6`＋一枚仓外导出树 `C:/Users/swq/tmp/wisp-306-orch-20261011-0821`）

| 腿的读数 | 我这发的独立复跑 |
|---|---|
| `AC#1` 覆盖块图三发并排：夹具 `mv` 走⇒五块全 0，夹具在与只跑那两枚新用例⇒五块全 1，PCM16 支与 default 支恒 0 | 我**⛔ 复跑块图**，用的是**等价的那一把**：发 A（pristine 导出树、具名两枚）＝`rc_A=0`／`--- PASS`=**2**／FAIL=0／SKIP=0（这一发同时是**台面自检**——具名用例真发出了颜色，⛔ 拿退码当"跑过了"）；发 D（把 `:218` 那一行的 `bits == 32` 改成 `16`）＝`rc=1`／FAIL **恰好那 2 枚**、红句逐字 `unsupported wav format: tag=3 bits=32 (want PCM16 or float32)`。⚠ **口径差别具名**：票面那枚形是 `fmtTag` 的 **`3`→`2`**（腿 `sens-driver.sh:21` 跑的就是那一形，红句同为 `tag=3 bits=32`），我这发改的是同一行的 `bits` 侧 ⇒ 本格凭据以**腿那一形**为主、我这发作对拉 |
| ★"⛔ 断被调用过"那一问＝拔掉 `err` 门后仍红且红句是解析值；另九枚逐枚翻转全红 | 我另起一把更窄的**等号倒置**尺（⛔ 动产码、只动用例）：**F1**＝`if rate != 48000`→`==` ⇒ `rc=1`／FAIL 1、红句逐字 `rate = 48000, want 48000`；**F2**＝`if len(inj.samples) != …`→`==` ⇒ `rc=1`／FAIL 1、红句 `injector holds 8 samples, want 8` ⇒ 断言**逐枚 Binding、⛔ 恒真装饰**（腿那九枚覆盖的是全族，我这两枚抽验的是"翻转必响"这一形） |
| `AC#2` `:192` `0xFFFE`→`0xFFFD`⇒两枚指名用例当场红、无夹具正控整包照绿 | 发 B 同形＝`rc_B=1`／FAIL **2**、红句逐字（`…306_test.go:177` 与 `:213`）`tag=65534 bits=32`；**发 E**＝把新用例 `mv` 走＋同一坏值＋**整包**＝`rc_E=0`／顶层 PASS **42**／FAIL **0**／SKIP 1 ⇒ 与腿的 `D2` 同形同读数；还原＝`git show HEAD:…` 回写后 `cmp -s`⇒`identical rc=0`，随后发 Z＝`rc=0`／PASS 2 |
| `AC#2b` 甲形（上下界同移）⇒ `tag=0` 的正常判词；乙形（派单字面）⇒ `panic` 且第二枚根本没跑到 | 发 C＝`data[body+24 : body+26]`→`data[body+26 : body+28]`＝`rc_C=1`／FAIL **2**／红句逐字 `unsupported wav format: tag=0 bits=32 …`。**两枚腿（落地腿＋验收腿）各自独立量到"字面那一发⛔ 是断言红而是 `panic`"，我今天⛔ 复跑乙形**（它⛔ 是本格凭据，只算可达性正控，腿件里完整留档）⇒ **裁：本格凭据取甲形**，见下面第 2 节那条追加更正 |
| `AC#3` 只动注释：三把"非注释行"尺全零命中、`numstat` 16/5；SDK 行号逐枚核中 | 我自己跑那把逐字尺＝**零命中**（`rc_grep_expect_zero=1`），`numstat` 逐字 **16 加／5 删**；★**SDK 行号我自己去本机那份头文件里逐枚 `sed`**（路径＝`C:/Program Files (x86)/Windows Kits/10/Include/10.0.26100.0/shared/`）：`mmreg.h:2110`＝`#define WAVE_FORMAT_IEEE_FLOAT 0x0003`✓、`:2376`＝`#define WAVE_FORMAT_EXTENSIBLE 0xFFFE`✓、`ksmedia.h:850`＝`#if defined(_INC_MMREG)`✓、`:854`＝`DEFINE_GUIDSTRUCT("00000003-…", KSDATAFORMAT_SUBTYPE_IEEE_FLOAT);`✓、`:856`＝`#endif` ⇒ **腿写进注释的号是对的，票面／派单那句 `851-855` 的守卫起点晚一行** |
| `AC#4` 门禁（腿判**附条件成立**，两枚条件⛔ 在它车道） | **条件①我这发补上了**：共享工作树 `GOFLAGS= go build ./...`＝**rc=0**（件 `g1-build-all-worktree.txt`）；**条件②也补上了**：`bash scripts/portable-tests.sh --scope=census`＝**rc=0**、totals 行逐字 `packages=35 with-zero-compiled-tests=7 claimed-by-no-scope=7 unclaimed-with-tests=0`，与落地腿件 `probes/306/r1/logs/40-census-raw.txt:38` **逐字相等**（⚠ 台面⛔ 同一枚：腿那发在它自己的台面、我这发在共享工作树 HEAD `1eaa91f6` ⇒ 只作对拉、⛔ 相减）。另我自己跑 `GOFLAGS= go vet ./internal/audio/`＝**rc=0**、`sh scripts/d22scan.sh`＝**rc=0** clean（ban #8 `internal/`=**527**／`cmd/`=**119**，仪器自证行 `scan_test.go:1261: verified ban #8 internal/: 527` 同数）、整包 `go test ./internal/audio/ -count=1 -v`＝`rc=0`／顶层 PASS **44**／FAIL **0**／SKIP **1**（唯一那枚 `TestLiveWasapiSmoke`，世代＝`hotplug_test.go`，⛔ 本票所造）⇒ **`AC#4` 判成立** |
| 产码一字未动 | 我这把换成**对象层尺**（比 `cmp` 硬）：`git rev-parse 5480434f:internal/audio/wavinjector.go` ＝ `git rev-parse HEAD:internal/audio/wavinjector.go` ＝ **`871f8eb9b3b27830d5ae7c0062268b0a5a3161e2`**；三枚行内权威值在导出树里逐字仍在原位（`:192` `if fmtTag == 0xFFFE {`／`:196` `data[body+24 : body+26]`／`:218` `case fmtTag == 3 && bits == 32:`）；scoped porcelain `-- internal cmd`＝**0 行** |
| 腿的⛔ 做到的四条＋越格自报那一枚 | 我复核成立。那两枚被 `rm` 的 gofmt 计数中间件＝**动作违规、后果零枚读数受损**（读数逐字在 `r1/logs/40-gofmt.txt`，腿另起三把尺独立复跑同一问）⇒ **⛔ 撤回那枚自报**；罪名定准＝**在共享仓根造未跟踪件**（仓根现量别家未跟踪件 47 枚，⛔ 我名下、⛔ 我删）。⚠ 另：仓根那枚 0 字节、文件名是一个减号 `-` 的件（`17` 个硬链接、mtime `2026-10-03 10:20`）＝某次重定向写歪的产物，腿具名报了两次 mtime 落在它两批发之间但⛔ 能归因 ⇒ **我⛔ 删**（别人的件、共享工作树），只在这里记它存在 |

### 2. ★两枚 append-only 更正（原句⛔ 改、⛔ 删，本节追加）

1. **`AC#2b` 的凭据形**：票面那句"只改那一处偏移、语句其余字不动"在盘上产出⛔ 一句断言红——它给的是 `panic: index out of range [1] with length 0`（合法空切片 `data[body+26 : body+26]`），**炸掉整个测试二进制**：顶层 PASS 从 **44 掉到 36**、两枚 306 用例里只有第一枚 `=== RUN` 过 ⇒ 第二枚**根本没跑到**；无夹具那一发 `rc=0`、`panic` 0 枚 ⇒ 它证的是"**有人走到**"，⛔"有人钉"，而且**对界内错偏移根本不响**。⇒ **本格凭据取甲形**＝同一枚切片的**上下界一起移**（`body+24 : body+26` → `body+26 : body+28`，改的仍只有"SubFormat 起点偏移"这一枚权威值，红句是一句正常的关于值的判词 `tag=0`）；字面那一发**降为可达性正控**、完整留档⛔ 丢。同盘先例＝姊妹钉 `parse_wave_format_300_windows_test.go:116-119` 把同一形状（`data1At: 26`）钉成一枚带 `wantTag: 0` 的**子测试**（腿现量它 PASS）。
2. **`AC#3` 的 SDK 号段**：票面／派单写 `ksmedia.h:851-855`；本机 `10.0.26100.0` 那份盘上真形＝守卫起 **`:850`**、`DEFINE_GUIDSTRUCT` 在 **`:854`**、外层 `#endif` 在 `:856` ⇒ **存在性判断⛔ 错、号段粒度写偏**，落地腿按盘上写 `:850-855` 是**对的**。今后引这一格一律按盘上号。

### 3. ★我自己两处判据措辞被两枚腿否证（记我，⛔ 改原句，指回本节）

1. 派单与票面 `AC#2b` 那句字面判据＝**一枚⛔ 可满足的判据**（我写的时候⛔ 知道它给的是 panic）。**先例已记过两次**（`A826` 票 303 `AC#3` 那半句、本轮 `AC#2b`）⇒ **定式＝派单里任何"只改 X 一处、其余字⛔ 动"的突变判据，落笔前先自己在那枚切片／表达式上算一遍它会产生什么形态的失败**（断言红／panic／空操作）。本轮两枚腿各自量到、结论一致，我⛔ 复跑第二发（它⛔ 本格凭据、⛔ 需要我复跑才能定罪它）。
2. 派单里 `ksmedia.h:851-855` 那一枚号＝我从票面抄来的**过期号**（两枚腿与我自己三把尺都给 `:850`）。⇒ **行号一律当快照**这条我自己也要照做，抄票面行号前先 `sed` 一次。

### 4. 欠着的读数（⛔ 算任何一枚腿的欠账，具名归编排者）

- **core 档（ubuntu）那一步的真 CI 日志**：票面 §5 那句"新文件在 CI 两档可见"现前只有**间接**读数（`portable-tests.sh:242` 的 core scope 逐字含 audio＋新文件第一行逐字 `package audio`＋linux 档 `go vet` rc=0）。坐实＝**push 之后取 CI 日志**＝我车道 ⇒ **随下一次推送一起取**。
- 覆盖块图那一把我⛔ 自己复跑（用等价的突变＋红句尺对拉）；三枚 `.cover` 件在 `306-v1/logs/sens/` 里逐块可查，下一位要复核就点那三枚。

### 5. 现态与 next

- 现态＝**票 306 六格全勾**（`AC#0`…`AC#4` 含 `AC#2b`），裁决表 1:1 已建（`docs/evidence/s1/306-inline-literals-v1.md`，README 规则 6 那道闸过），**⛔ 改 `-done` 由我按住到"新文件在 CI core 档那一发的颜色回来"**（理由见第 4 节第一枚；⛔ 收了再红）。`Status:` 行⛔ 写 `done`。
- 残余两件⛔ 属本票射程：① 块 `194.6,195.1`（`size<40` 错误支）与裸 `tag=3`（非 extensible）那一形**仍零执行者**＝票面 §边界明写⛔ 买，谁要闭另立射程；② 上面那枚 CI 色。
- next＝票 306 收口（改名 `-done`＋更新索引）排在**推送＋core 档色回来**之后；本程在飞＝`307-a1`（票 307 `AC#0` 只读）∥ `303-v2`（补票 303 那张缺两行的 1:1 裁决表）；我自己那六发在 `cmd/wisp` 车道（票 305 判死那一族，件 `probes/305/a2/30-deciding-run-recipe.md`）排在**两枚腿交完、编译面空**之后。
- Progress log：
  - [2026-10-11 08:2x +0800] agent=编排者 did=收 `306-v1`（五笔，13 笔越界逐笔我自己重跑＝写面外 **0**、`.out` **0**；`.go` 分母那一枚它自己 `git mv` 还回去了、现量 `probes/306/v1` 下 tracked `.go`＝0）＋**承重读数我自己现跑**（一枚仓外导出树 `C:/Users/swq/tmp/wisp-306-orch-20261011-0821`，⛔ 入库、⛔ 删）：发 A `rc=0`/PASS **2**（＝台面自检，具名用例真发出颜色）／发 B `rc=1`/FAIL **2** 红句 `tag=65534`／发 C `rc=1`/FAIL **2** 红句 `tag=0`／发 D `rc=1`/FAIL **2** 红句 `tag=3`／发 E（`mv` 走新用例＋坏值＋整包）`rc=0`/顶层 PASS **42**/FAIL **0**／F1·F2 等号倒置各 `rc=1`（`rate = 48000, want 48000`／`injector holds 8 samples, want 8`）／两枚文件对 HEAD blob `cmp` 各 `identical rc=0`／发 Z `rc=0`；★**产码零动用对象层尺**（`5480434f` 与 HEAD 的 `wavinjector.go` blob 同为 `871f8eb9…`）；门禁四把我自己跑＝`go build ./...` **rc=0**（共享工作树，＝腿的条件①）、`--scope=census` **rc=0** 且 totals 行与腿件 `40-census-raw.txt:38` **逐字相等**（条件②；⚠ 两把台面⛔ 同一枚，只作对拉）、`go vet ./internal/audio/` **rc=0**、`d22scan` **rc=0** clean（ban #8 `internal/`=527／`cmd/`=119）、整包 `rc=0`／PASS **44**／FAIL **0**／SKIP **1**；SDK 五枚号我自己去 `10.0.26100.0` 那份里 `sed`＝`mmreg.h:2110`/`:2376`✓、`ksmedia.h:850`/`:854`/`:856` ⇒ **腿对、票面与我派单错**；裁决表 1:1 建好＝`docs/evidence/s1/306-inline-literals-v1.md`（六行、每行带尺＋射程＋出处三态）。⇒ **翻 `AC#1`/`AC#2`/`AC#2b`/`AC#3`/`AC#4` 五格**（判语出自非实现者腿 `306-v1`，⛔ 落地腿自勾）；★**两枚 append-only 更正落我账上**＝`AC#2b` 字面判据⛔ 可满足（盘上给的是 `panic`、第二枚指名用例根本没跑到 ⇒ 本格凭据取**甲形**上下界同移、字面那发降为可达性正控留档）＋`AC#3` 的 `ksmedia.h:851-855` 真形 `:850-856`；★**我自己那枚措辞坑第二次踩同形**（派单里"只改 X 一处"的突变判据落笔前⛔ 自己算一遍它产什么形态的失败——第一次＝`A826` 票 303 `AC#3`）；⛔ 改 `-done`＝**按到"新测试文件在 CI core 档那一发的颜色回来"**（现前只有间接读数：core scope 逐字含 audio＋新文件首行逐字 `package audio`＋linux `go vet` rc=0；⛔ 收了再红）；腿⛔ 做到的四条与越格自报一枚复核成立（**⛔ 撤回那枚自报**，罪名定准＝在共享仓根造未跟踪件⛔"删"，读数⛔ 受损；仓根那枚文件名是 `-` 的 0 字节件 mtime `10-03 10:20`⛔ 能归因⇒**我⛔ 删**、只记它存在）。现态＝票 306 **六格全勾、⛔ `-done`**；本程在飞＝`307-a1`∥`303-v2`；next＝腿交完→我跑票 305 那六发（`cmd/wisp` 车道）→推送＋取 core 档色→票 306 收口；⛔ push（照 `A830` §6 两条理由：推送自启 `slo-full` 抢机主 CPU＋票 300/303/111 那批 CI 色要同批取）
