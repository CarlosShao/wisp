# 票 306 · AC 编号 1:1 裁决表（六行：`AC#0`／`AC#1`／`AC#2`／`AC#2b`／`AC#3`／`AC#4`）

**这张表是什么**：工单池规矩 6（`.scratch/wisp/issues/README.md`，机主 2026-09-20 批准）要求的"与 AC 编号 1:1 的裁决表"，
缺行即 FAIL。**判语出处＝非实现者验收腿 `306-v1`**（件＝`.scratch/wisp/probes/306/v1/30-judgments.md`，实现者是 `306-r1`，
裁决者≠实现者＝`SPEC-12 §4.3`）；本文件由编排者**汇编**，并把**编排者自己现跑的对拉尺**逐行标出来。
**编排者自己现跑的**那一发绝不冒充腿的判语，腿的判语也⛔ 被我改成更好看的形。

- 台账：`A825`（立票）／`A827`（收 `306-a1`）／`A830`（收 `306-r1`）／`A831`（收 `306-v1`＋本表）
- 腿件：`306-a1`（只读三张表，`probes/306/a1/**`）／`306-r1`（落地腿八笔 `4e73d0a5`→`e4740e35`，`probes/306/r1/**`）／`306-v1`（非实现者五笔 `ec87076f`→`1eaa91f6`，`probes/306/v1/**`）
- 编排者现跑件：`probes/306/orch/logs/**`（`orch-A/B/C/D/E/F1/F2/Z`＋`orch-G2/G3/G4`＋`g1-build-all-worktree`＋`g6-census-worktree`＋`orch-03-roster-306-commits.txt`）
- 台面（编排者 08:2x 现量）＝HEAD `1eaa91f6`；`git status --porcelain -- internal cmd`＝**0 行**；产码 `internal/audio/wavinjector.go` 的 blob `git rev-parse 5480434f:…` ＝ `git rev-parse HEAD:…` ＝ **`871f8eb9b3b27830d5ae7c0062268b0a5a3161e2`**（对象层尺，比 `cmp` 硬）

**出处三态**（每行必带）：〔编排者现跑〕＝我自己这一把；〔非实现者腿现跑〕＝`306-v1`；〔落地腿自陈〕＝`306-r1`（**只算待验断言**）；〔读码推到〕＝⛔ 任何一格凭据。

---

## `AC#0` 只读普查（决定落点）＝**成立**（`A827` 已翻勾，本行只补尺）

| 要 | 内容 |
|---|---|
| 判据 | 交回三张表：三枚行内权威值名册／覆盖块图（哪几块 hit=0）／落点候选 |
| 凭据 | 编排者自己在同台面复跑过覆盖尺与块级零命中六块（`A827` 正文逐字）：`193.5,193.18`／`194.6,195.1`／`196.5,196.65`／`219.3,221.26`／`222.4,223.1`／`224.3,224.43` |
| 出处 | 〔编排者现跑〕（`A827`）＋〔非实现者腿现跑〕（`306-v1` `logs/sens/c2-nofixture-fullpkg.cover`：夹具 `mv` 走后我关心的五块**全 0**） |
| ★第三枚的来历 | `:196` 的 SubFormat 起点偏移 `body+24` 是 `306-a1` 表① 现量**新增**的一枚（立票时⛔ 数到）⇒ 编排者把它入成**新格 `AC#2b`**，并写死"一枚夹具装三枚牙⛔ 等于三格都闭" |

## `AC#1` 决定性一发（让 `0xFFFE`＋`bits==32` 那一支第一次有执行者）＝**成立**

| 要 | 内容 |
|---|---|
| 判据（票面逐字） | 让那一支**真的被执行**，断言的是**解析出来的那个值**、⛔ 断"被调用过"；★凭据形状⛔ 变：把 `:218` 的 `3` 换成 `2` 之后**必须红** |
| 非实现者那把尺 | 覆盖块图三发并排（`-coverprofile`，set 模式）：夹具 `mv` 走 ⇒ 五块全 0；夹具在 ⇒ 五块全 1；**只跑那两枚新用例**⇒ 五块**仍全 1**，而 PCM16 支 `212.3,214.26` 与 default 错误支 `226.3,226.111` **保持 0**。件＝`logs/sens/c1-head-fullpkg.cover`／`c2-nofixture-fullpkg.cover`／`c3-newtests-only.cover` |
| 编排者对拉（⛔ 复跑覆盖块图，用等价的那一把） | **发 A**＝导出树 pristine、`-run '^(TestParseWavExtensibleFloat32306|TestWavInjectorExtensibleFloat32306)$'` ⇒ `rc_A=0`／`--- PASS`=**2**／`--- FAIL`=0／`--- SKIP`=0（这一发同时是**台面自检**：具名用例真发出了颜色，⛔ 拿退码当"跑过了"）。<br>**发 D**＝`case fmtTag == 3 && bits == 32:` → `bits == 16` ⇒ `rc=1`／FAIL **恰好那 2 枚**，红句逐字 `parseWav error = unsupported wav format: tag=3 bits=32 (want PCM16 or float32) …`。<br>⚠ 口径差别具名：票面那枚形是 **`fmtTag` 的 `3`→`2`**（`306-v1` 的 `sens-driver.sh:21` 跑的就是这一形，红句同为 `tag=3 bits=32`），我自己跑的是**同一行的 `bits` 侧**——两形都证"该支由夹具执行、值一坏就响"，本格凭据以**腿那一形**为主、我这发作对拉。<br>**"⛔ 断被调用过"那一问**＝腿的 `VR` 组：拔掉 `err` 门只留值断言⇒仍红且红句是解析值（`rate = 0, want 48000`／`len(samples) = 0, want 8`）；编排者另跑 **FLIP 两枚**（⛔ 动产码、只把用例的等号倒过来）⇒ `F1 rc=1` 红句 `rate = 48000, want 48000`、`F2 rc=1` 红句 `injector holds 8 samples, want 8` ⇒ 断言**逐枚Binding、⛔ 恒真装饰** |
| 具名残余（⛔ 扣本格、⛔ 本票射程） | 块 `194.6,195.1`（`size < 40` 的错误 return）在三发里**恒 0**＝新夹具⛔ 造那一枚世界（文件 `:61-66` 自己写死⛔ 装它，票面 §边界第 4 条也写着⛔ 顺手补别的分支）⇒ **设计上的⛔ 买**；谁要闭它要另立射程。裸 `tag=3`（非 extensible）那一形同理 |
| 出处 | 〔非实现者腿现跑〕＋〔编排者现跑〕（`orch-A`／`orch-D`／`orch-F1`／`orch-F2`） |

## `AC#2` 同族第二枚的牙（`:192` 的 `0xFFFE`→`0xFFFD`）＝**成立**

| 要 | 内容 |
|---|---|
| 判据 | 一对：改⇒指名用例红／`cmp` 逐字节还原⇒绿 |
| 非实现者那把尺 | `logs/sens/P1-mut192.txt`（自己的 `perl` 取坏值，⛔ 复用腿的 `mutate.sh`）⇒ `rc=1`／FAIL 2 枚／红句逐字 `tag=65534 bits=32`；整包那一发 `logs/sens/B4-mut192-fullpkg-v.txt`＝FAIL 恰好那 2 枚、其余顶层 PASS **42** 逐字绿；**反面**（这枚牙只属于新夹具）＝`logs/D2-nofixture-mut192.txt`：同一坏值＋新用例 `mv` 走 ⇒ `rc=0`／PASS 42／FAIL 0 |
| 编排者对拉 | **发 B**＝同形 mutation ⇒ `rc_B=1`／FAIL **2**／PASS 0，红句逐字（`wavinjector_extensible_float_306_test.go:177` 与 `:213`）`unsupported wav format: tag=65534 bits=32 (want PCM16 or float32)`；**发 E**＝把新用例 `mv` 走＋同一坏值＋**整包** ⇒ `rc_E=0`／顶层 PASS **42**／FAIL **0**／SKIP 1 ⇒ 与腿的 `D2` 同形同读数；**还原**＝`git show HEAD:…wavinjector.go` 回写后 `cmp -s` ⇒ `identical rc=0`，随后 **发 Z**＝`rc=0`／PASS 2 |
| 出处 | 〔非实现者腿现跑〕＋〔编排者现跑〕（`orch-B`／`orch-E`） |

## `AC#2b` 第三枚行内权威值的牙（`:196` 的 `body+24`）＝**附条件成立（条件已由编排者落成票面追加更正，见下面第 2 节）**

| 要 | 内容 |
|---|---|
| 判据（票面原句） | `:196` 的 SubFormat 偏移 `body+24`→`body+26` 必须让**同一枚**指名用例红；⛔ 与 `AC#1`/`AC#2` 共用一发读数；成对两发缺一⛔ 算闭 |
| ★盘上否证（两枚腿各自量到、编排者复跑） | **字面那一发（只动上界起点 ⇒ `data[body+26 : body+26]`＝合法空切片）产出⛔ 一句断言红**，而是 `panic: runtime error: index out of range [1] with length 0`（栈里逐字 `encoding/binary.littleEndian.Uint16(...) binary.go:70` ＋ `wavinjector.go:196`），它**炸掉整个测试二进制**：顶层 PASS 从 **44 掉到 36**、两枚 306 用例里只有第一枚 `=== RUN` 过 ⇒ 第二枚**根本没跑到**。无夹具那一发（`logs/D5-nofixture-mut196literal.txt`）`rc=0`、`panic` 0 枚 ⇒ 崩溃只在"夹具执行到 `:196`"时发生 ⇒ **它证明"有人走到"，⛔ 证明"有人钉"**，且对界内错偏移根本不响 |
| 本格凭据（**甲形**＝同一枚切片的上下界一起移，语句其余字⛔ 动） | 〔非实现者腿现跑〕`logs/sens/P3-mut196-boundsshift.txt` ⇒ `rc=1`／FAIL 2 枚／红句 `tag=0 bits=32`＝**一句关于那枚偏移读出来的值的判断**；整包 `logs/D6-…-fullpkg.txt`＝FAIL 恰好 2、其余 42 绿；无夹具正控 `logs/D4-…` ⇒ `rc=0`／PASS 42。<br>〔编排者现跑〕**发 C**＝`data[body+24 : body+26]` → `data[body+26 : body+28]` ⇒ `rc_C=1`／FAIL **2**／红句逐字 `unsupported wav format: tag=0 bits=32 …`；还原 `cmp` ⇒ `identical rc=0` |
| 同盘先例（支持甲形） | 姊妹钉 `internal/audio/parse_wave_format_300_windows_test.go:116-119` 把同一形状（`data1At: 26`）钉成一枚带 `wantTag: 0` 的**子测试**，`306-v1` 现量它 PASS（件 `logs/Y3-ticket300-pins.txt`），⛔ 一次崩溃 |
| 处置 | 票面 `AC#2b` 的字面判据⛔ 可满足 ⇒ 编排者**append-only 追加更正**（原句⛔ 改、框由我翻）；字面那一发降为**可达性正控**完整留档 |

## `AC#3` 只改注释那一格＝**成立**

| 要 | 内容 |
|---|---|
| 判据 | 把 `wave_format_float_300_windows_test.go:38-41` 那句改成实话；**硬边界**＝⛔ 动任何一行断言／任何一枚期望值；⛔ 删 `waveFormatPCM`；⛔ 顺手改别的注释 |
| 非实现者那把尺（四把并排） | ① 腿那把逐字重跑（`git show c0d5ec24 -- <file> \| grep -E '^[+-]' \| grep -vE '^(\+\+\+\|---)\|^[+-]//'`）⇒ **零命中**；② 换区间口径（`git diff 5480434f HEAD -- <file>`）⇒ 零命中；③ 更严一把（也豁免缩进 `//`）⇒ 仍零命中；`numstat` **16 加／5 删**；④ 票 300 那两枚钉的形状自跑 ⇒ 顶层 1＋子测试 6、顶层 1＋子测试 4，`rc=0` 逐枚 `--- PASS`（件 `logs/Y3-ticket300-pins.txt`）；`waveFormatPCM` ⛔ 删：`grep -c`＝`wasapi_windows.go` **1**／该测试件 **7** |
| 编排者对拉 | ①②③同一把尺我自己跑 ⇒ **非注释行命中 0**（`rc_grep_expect_zero=1`＝零命中的退码），`numstat` 逐字 **16／5**；<br>★SDK 行号我**自己去本机那份头文件里逐枚核**（路径＝`C:/Program Files (x86)/Windows Kits/10/Include/10.0.26100.0/shared/`）：`mmreg.h:2110` 逐字 `#define WAVE_FORMAT_IEEE_FLOAT 0x0003`、`mmreg.h:2376` 逐字 `#define WAVE_FORMAT_EXTENSIBLE 0xFFFE`、`ksmedia.h:850` 逐字 `#if defined(_INC_MMREG)`、`:854` 逐字 `DEFINE_GUIDSTRUCT("00000003-…", KSDATAFORMAT_SUBTYPE_IEEE_FLOAT);`、`:856` 逐字 `#endif` ⇒ **腿写进注释的号＝对的**；票面／派单那句 `ksmedia.h:851-855` 的**守卫起点晚一行**（`:851` 是里层 `#if !defined( STATIC_… )`）——存在性判断⛔ 错、号段粒度写偏 ⇒ append-only 追加更正（第 2 节） |
| 出处 | 〔非实现者腿现跑〕＋〔编排者现跑〕（`orch-02-ac3-comment-only.txt`＋上面那五枚 `sed`） |

## `AC#4` 门禁＋越界＝**附条件成立 → 两枚条件都由编排者在自己车道闭掉，本格判成立**

| 项 | 读数 | 出处 |
|---|---|---|
| `GOFLAGS= go build ./...`（**条件①，编排者补的那发**） | **rc=0**，射程＝**共享工作树**（腿只取到导出树那一发） | 〔编排者现跑〕`g1-build-all-worktree.txt` |
| `bash scripts/portable-tests.sh --scope=census`（**条件②，编排者补的那发**） | **rc=0**，totals 行逐字 `portable-tests.sh: census totals: packages=35 with-zero-compiled-tests=7 claimed-by-no-scope=7 unclaimed-with-tests=0` ⇒ 与落地腿件 `probes/306/r1/logs/40-census-raw.txt:38` **逐字相等**（同数≠同台面：腿那发在它自己的台面，我这发在 HEAD `1eaa91f6` 的共享工作树，本票⛔ 跨台面减数，只作对拉） | 〔编排者现跑〕`g6-census-worktree.txt` |
| `GOFLAGS= go vet ./internal/audio/` | **rc=0**（共享工作树，⛔ 输出行） | 〔编排者现跑〕`orch-G2-vet-worktree.txt`（腿那一发同 rc=0） |
| `sh scripts/d22scan.sh` | **rc=0**，逐字 `clean - no D22 ban violations`；ban #8 射程 `internal/`=**527**、`cmd/`=**119**，仪器自证行 `scan_test.go:1261: verified ban #8 internal/: 527` 同数 | 〔编排者现跑〕`orch-G3-d22scan-worktree.txt` |
| 整包名册（⛔ 用交集尺） | `GOFLAGS= go test ./internal/audio/ -count=1 -v`（共享工作树）⇒ `rc=0`／顶层 `--- PASS`=**44**／`--- FAIL`=**0**／`--- SKIP`=**1**（唯一那枚 `TestLiveWasapiSmoke`，其**世代钉死**＝`5480434f:internal/audio/hotplug_test.go` 里改前就有，⛔ 本票所造）。腿另取 **5 发**逐对 `diff` rc=0 ⇒ 名册逐字相同；**"只红过一次被藏住"那一形在这 5 发里没出现** ⇒ 成论只写"稳定新增红 0 枚"，并具名承认交集尺的结构性盲区 | 〔编排者现跑〕`orch-G4-fullpkg-roster-head.txt`＋〔非实现者腿现跑〕`logs/sens/B1..B5` |
| 改前／改后**同台面** | 改前＝同树把新用例 `mv` 走 ⇒ 顶层 PASS **42**（腿 `B3`/`D1`，编排者 `orch-E` 同数）；改后＝**44** ⇒ "＋2" 这一枚 delta 是**同树**量出来的，⛔ 两树相减；旁证一把＝`^func Test` 计数 `5480434f`=43 → HEAD=**45**。落地腿另有两枚仓外导出树（`git archive 5480434f`／`git archive fbd150e3`）各 ≥2 发＋`comm -13`／`comm -23` **双向 0 行**（五对，件＝`probes/306/r1/40-gates.md:70-74`） | 〔编排者现跑〕＋〔落地腿自陈〕（那五对我⛔ 逐对复跑，颜色由上面两把独立尺覆盖） |
| 三数带尺名 | `^--- PASS` 顶层＝**44**；`^    --- PASS` 含子测试＝**12**；两把⛔ 相加（票面老规矩） | 〔编排者现跑〕 |
| 格式名册**两把并排** | 工作树 `gofmt -l internal/audio`＝**0 枚**；blob 那一把（`git archive <sha> internal/audio` 落仓外再 `gofmt -l`）＝`5480434f`（21 枚 `.go`）**0 枚**／HEAD（22 枚）**0 枚** ⇒ 新增未格式 0 枚（分母 21→22＝那枚新文件）。⚠ 腿具名报过：全仓那一把它量到 `rc=2`／44 行而 `internal/audio` 命中 0 ⇒ **⛔ 写"全仓 0 枚"** | 〔非实现者腿现跑〕`g4*`／`g5-*` 件 |
| 越界**逐笔**（⛔ 区间尺） | 编排者自己在 **13 笔**（`306-r1` 八笔＋`306-v1` 五笔）上重跑 `git show --name-only --format=%H`：每笔**写面外＝0**、`.out`＝**0**；新建 `.go` 只有两笔（`1b34e475` 那枚新测试文件、`c0d5ec24` 那枚注释文件，⛔ 在 `internal/audio/**` ＝票面写面内）；`306-v1` 曾把两枚 HEAD blob 副本命名 `.go` 落进 `logs/`（`a458114b` 那笔 新go=2），第 5 笔 `1eaa91f6` 用 `git mv` 改成 `.txt`（内容 `cmp rc=0` 两把都过）⇒ **把 `.go` 分母还给别家的尺**，处置正确。并集尺＝`5480434f..HEAD` 碰 `internal/audio` 的**只有那 2 枚** | 〔编排者现跑〕`orch-03-roster-306-commits.txt` |
| 禁改清单 | 7 枚逐枚 `git diff --quiet 5480434f HEAD -- <f>`：6 枚 rc=0，第 7 枚 rc=1＝它**确实变了**，变的只有 `AC#3` 那 16/5 枚注释行；`frontend design third_party` 区间命中 0 枚；被跟踪的 `.wav`＝**0**（`git ls-files \| grep -icE '\.wav$'`＝0）⇒ ⛔ 开"往仓里塞二进制夹具"那处先例 | 〔落地腿自陈〕＋〔编排者现跑〕（`git status --porcelain -- internal cmd`＝0 行＋blob hash 那把） |
| 票面纪律 | `306-r1` 与 `306-v1` 各只追加一枚 Progress log bullet（append-only；`306-v1` 那笔 numstat **1 加 0 删**）；框普查全程 **1 勾／5 未勾**＝⛔ 任何腿翻过任何东西；⛔ 零 push | 〔非实现者腿现跑〕`logs/BF-preappend.txt`／`BG-postappend.txt`＋〔编排者现跑〕 |

---

## 编排者两枚 append-only 更正（落进票面，原句⛔ 删）

1. **`AC#2b` 的凭据形**：票面那句"只改那一处偏移、语句其余字不动"在盘上产出⛔ 一句断言红（它给的是 `panic`＋第二枚指名用例根本没跑到）。
   ⇒ 本格凭据取**甲形**＝同一枚切片的**上下界一起移**（`data[body+24 : body+26]` → `data[body+26 : body+28]`，改的仍只有"SubFormat 起点偏移"这一枚权威值），
   字面那一发**降为可达性正控**、完整留档（⛔ 丢）。同盘先例＝姊妹钉把同形钉成带 `wantTag: 0` 的子测试。
2. **`AC#3` 的 SDK 号段**：票面／派单写 `ksmedia.h:851-855`；本机 `10.0.26100.0` 那份盘上真形＝守卫起于 **`:850`**、`DEFINE_GUIDSTRUCT` 在 **`:854`**、外层 `#endif` 在 `:856`
   ⇒ 存在性判断⛔ 错、号段粒度写偏；落地腿按盘上写 `:850-855` 是**对的**。

## 欠着的读数（⛔ 算任何一枚腿的欠账，具名归编排者）

- **core 档（ubuntu）那一步的真日志**：票面 §5 那句"新文件在 CI 两档可见"里，core 那一档现前只有间接读数（`portable-tests.sh:242` 的 scope 逐字含 audio＋新文件第一行逐字 `package audio`）。
  坐实它＝**push 之后取 CI 日志**，那是编排者车道（腿自陈"⛔ 我车道"）。⇒ 这一发**随下一次推送一起取**。
- 覆盖块图那一把（`-coverprofile` 三发并排）编排者**⛔ 自己复跑**，我用的是等价的突变＋红句那一把。
  块图那三份 `.cover` 在 `306-v1` 的件里逐块可查，下一位要复核就点那三枚件。

## 现态

票 306＝**六格全勾**（`AC#0`…`AC#4`，含 `AC#2b`），⛔ `-done` 由本程收口动作执行（改名＋更新索引），Status 行照规矩写。
残余两件⛔ 属本票：① `194.6,195.1`（`size<40` 错误支）与裸 `tag=3` 那一形零执行者＝票面 §边界明写⛔ 买；② core 档 CI 日志＝上面那枚欠账。
