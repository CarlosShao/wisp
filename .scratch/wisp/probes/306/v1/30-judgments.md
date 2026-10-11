# 306-v1 · 30 五格判语（逐格独立取证，每格带本腿自己跑过的尺）

判语只⛔ 三形：`成立`／`附条件成立`／`不成立`。⛔ 本件翻任何一枚 `AC` 框、⛔ 写"翻勾"、⛔ 动 `-done`——
框由编排者核过之后自己翻。每格那句"尺"里的每一枚读数都在 `10-sensitivity.md`／`20-scope-and-gates.md`，原始件在 `logs/**`。

台面（三枚权威值的内容锚，blob＝工作树 `rc=0`）：`:192` `if fmtTag == 0xFFFE {`／`:196` `data[body+24 : body+26]`／`:218` `case fmtTag == 3 && bits == 32:`。

---

## `AC#1` 决定性一发（让那一支第一次被执行）＝**成立**

**本格要的凭据**＝"让 `0xFFFE`＋`bits==32` 那一支真的被执行"＋"断言的是解析出来的那个值，⛔ 断'被调用过'"。

**本腿那把尺（三把并排，各自射程＝仓外导出树 `git archive e4740e35`，量的是那棵树里的工作面）**

1. **覆盖块图 before/after**（`-coverprofile`，set 模式，块图逐枚贴＝`10-sensitivity.md` §1.5；
   件＝`logs/sens/c1-head-fullpkg.cover`／`c2-nofixture-fullpkg.cover`／`c3-newtests-only.cover`）：
   夹具 `mv` 走 ⇒ `193.5,193.18`／`196.5,196.65`／`219.3,221.26`／`222.4,223.1`／`224.3,224.43` **五块全 0**（＝票面 §现量 那六块的本腿复跑）；
   夹具在 ⇒ 五块**全 1**；**只跑那两枚新用例**⇒ 五块**仍全 1**，且 `212.3,214.26`（PCM16 支）与 `226.3,226.111`（default 错误支）**保持 0**。
2. **"⛔ 断被调用过"那一问**＝拔掉 `err` 门之后仍红（`VR` 组：`rc=1`，红句逐字 `rate = 0, want 48000`／`chans = 0, want 2`／`len(samples) = 0, want 8`）。
3. **每一枚断言都在Binding**＝逐枚翻转判据（`FLIP` 组：9 枚逐枚 `!=`→`==`，9/9 `rc=1` 且红句是那一枚自己的行号）。

**判语理由**：三处行内权威值第一次有执行者这一句，本腿⛔ 是读腿的名册而是自己取块图取到的；
"同一支"这一问也钉住了（`C3` 那一列：`:192`→`:196`→`:218`→float32 三支一起从 0 变 1，PCM16/default 两支仍 0）。

⚠ **具名残余（⛔ 扣本格，⛔ 本票射程）**：`194.6,195.1`（`size < 40` 的错误 return）在三发里**恒 0**
＝新夹具⛔ 造"extensible 但 size<40"那一枚世界。文件第 61-66 行自己写死了⛔ 装它，票面 §边界第 4 条也写着⛔ 顺手补别的分支
⇒ 这一枚是**设计上的⛔ 买**，⛔ `AC#1` 的缺口；谁下一程要闭它要另立射程。裸 `tag=3`（非 extensible）那一形同理。

---

## `AC#2` 同族第二枚的牙（`:192` `0xFFFE`→`0xFFFD` 必须让指名用例红）＝**成立**

**本格要的凭据**＝一对：改⇒红／还原⇒绿。

**本腿那把尺**：`logs/sens/P1-mut192.txt`——产码那行改成 `if fmtTag == 0xFFFD {`（本腿自己的 `perl` 取坏值，⛔ 复用腿的 `mutate.sh`）
⇒ `rc=1`、`--- FAIL` **2 枚**（`TestParseWavExtensibleFloat32306`＋`TestWavInjectorExtensibleFloat32306`逐枚）、
红句逐字 `parseWav error = unsupported wav format: tag=65534 bits=32 …`；
整包名册那一发（`logs/sens/B4-mut192-fullpkg-v.txt`）＝`rc=1`／FAIL 恰好那 2 枚／其余顶层 PASS **42** 逐字绿（⛔ 连带红）；
"改⇒红"的**反面**（＝这一枚牙只属于新夹具）＝`logs/D2-nofixture-mut192.txt`：同一枚取坏值、把新用例 `mv` 走 ⇒ `rc=0`／PASS 42／FAIL 0。
还原⇒绿＝`cmp -s` 逐字节（`rc_final_cmp_inj=0`，件＝`logs/sens/ZZ-summary.txt` 末段）之后，同一棵树整包 `rc=0`、PASS 44（件＝`logs/sens/B*-fullpkg-shot*.txt`）。

⚠ **与腿的那把尺的差别（⛔ 冲突，只是口径）**：腿用 `cmd | tee; echo $?` 之外的纯重定向并留了 `cmp` 逐字节件；
本腿用 `cp` 回写＋`cmp -s` 校验＋复跑绿。两形都⛔ 依赖腿的件。

---

## `AC#2b` 第三枚行内权威值的牙（`:196` 的 `body+24`→`body+26`）＝**附条件成立**

**选形＝甲**（详见 `10-sensitivity.md` §2，那两问在那里逐字答了）。本格凭据＝本腿自己重跑的这两发：

- **甲**（同一枚切片上下界同移，`data[body+26 : body+28]`）＝`logs/sens/P3-mut196-boundsshift.txt`：`rc=1`、FAIL 2 枚、
  红句 `tag=0 bits=32`＝**一句关于那枚偏移读出来的值的判断**；整包那一发 `logs/D6-…-fullpkg.txt`：FAIL 恰好 2 枚、其余 42 枚绿；
  无夹具正控 `logs/D4-nofixture-mut196boundsshift.txt`：`rc=0`／PASS 42。
- **乙**（派单字面＝只动起点，`data[body+26 : body+26]`）＝`logs/sens/P4-mut196-literalwording.txt`＋整包 `logs/N4-fixture-mut196literal-fullpkg-v.txt`：
  `rc=1`、FAIL **1 枚**、`panic: runtime error: index out of range [1] with length 0 [recovered, repanicked]`、
  顶层 PASS 从 **44 掉到 36**、两枚 306 用例里只有第一枚 `=== RUN` 过 ⇒ 第二枚指名用例**根本没跑到**。
  它⛔ 是"钉"，它是"走到"：无夹具那一发（`logs/D5-nofixture-mut196literal.txt`）`rc=0`、`panic` 0 枚 ⇒ 崩溃只在"新夹具执行到 `:196`"时发生。

**条件（这就是那句要编排者办的）**：票面 `AC#2b` 的字面判据"只改那一处偏移、语句其余字不动"在本仓盘上产出⛔ 一句断言红，
⇒ 请**追加**（append-only，⛔ 删原句、⛔ 动框）一条更正：本格凭据取"上下界同移"那一发，字面那一发降为**可达性正控**；
同盘先例＝姊妹钉 `internal/audio/parse_wave_format_300_windows_test.go:116-119` 把同一形状（`data1At: 26`）钉成
一枚带 `wantTag: 0` 的**子测试**（本腿现量它 PASS，件＝`logs/Y3-ticket300-pins.txt`），⛔ 一次崩溃。
条件未落之前，本腿这句判语的有效范围＝**甲那一发成立；字面那一发⛔ 算本格凭据**。

---

## `AC#3` 只改注释那一格＝**成立**

**本腿那把尺（四把并排，`20-scope-and-gates.md` §1.3＋§2）**

1. 腿那把逐字重跑（`git show c0d5ec24 -- <file> | grep -E '^[+-]' | grep -vE '^(\+\+\+|---)|^[+-]//'`）⇒ **零命中**（`rc=1`）；
   同一把尺换成区间口径（`git diff 5480434f HEAD -- <file>`）⇒ **零命中**；本腿另加一把更严的（也豁免缩进 `//`）⇒ **仍零命中**；`numstat` **16/5**。
2. ⛔ 动任何一枚期望值／⛔ 放宽票 300 的钉＝本腿自己跑那两枚钉的形状：顶层 1＋子测试 6、顶层 1＋子测试 4，`rc=0`，逐枚 `--- PASS`（件＝`logs/Y3-ticket300-pins.txt`）。
3. `waveFormatPCM` ⛔ 删＝`grep -c`：`wasapi_windows.go` **1**／`wave_format_float_300_windows_test.go` **7**。
4. **三枚 SDK 行号逐枚去本腿机器上那份 `10.0.26100.0` 树里核**（件＝`logs/20-sdk-lines.txt`）：
   `mmreg.h:2418`／`:2110`／`:2376`／`:2480-2484`（`DEFINE_GUIDSTRUCT` 在 `:2483`）／`ksmedia.h:850`（`#if defined(_INC_MMREG)`）／`:854`（`DEFINE_GUIDSTRUCT`）**全部逐字对上**
   ⇒ 注释里写盘上的号＝**对的**；票面／派单写的 `ksmedia.h:851-855` 是守卫起点晚一行（`:851` 是里层 `#if !defined( STATIC_… )`），
   存在性判断⛔ 错、号段粒度写偏；本腿**⛔ 因此翻⛔ 翻任何勾，只报**（票面原句⛔ 改，append-only 的话留给编排者）。
   票面那句"其兄弟件 `:101` 引 `mmreg.h:2483`"本腿逐字核中（`parse_wave_format_300_windows_test.go:101`）。

---

## `AC#4` 门禁＋越界＝**附条件成立**

**成立的那一面（本腿现量，⛔ 抄腿）**

| 项 | 本腿读数 |
|---|---|
| `GOFLAGS= go vet ./internal/audio/` | **rc=0**（共享工作树，⛔ 输出行） |
| `sh scripts/d22scan.sh` | **rc=0**，逐字 `clean - no D22 ban violations`；ban #8 射程 `internal/`=**527**、`cmd/`=**119**（＋仪器自证行 `verified ban #8 internal/: 527` 同数） |
| 格式名册两把并排 | 工作树 `gofmt -l internal/audio`＝**0 枚**；blob `git archive <sha> internal/audio` 后 `gofmt -l`＝`5480434f`（21 枚 `.go`）**0 枚**／HEAD（22 枚）**0 枚** ⇒ 新增未格式 **0 枚**（分母 21→22＝那枚新文件） |
| 红名册（⛔ 交集尺） | 单发整包 `-v` **5 发**，每发 `rc=0`／`^--- FAIL`=**0**／`^--- PASS`=**44**／名册逐字同（唯一那枚 `--- SKIP: TestLiveWasapiSmoke` 在 `5480434f` 的 `hotplug_test.go` 里就有）⇒ 逐对 `diff` **rc=0**；"只红过一次被藏住"那一形在 5 发里**没出现**（⛔ 时有时无的枚） |
| 改前／改后同台面 | 改后＝5 发（PASS 44）；改前＝同树 `mv` 走新用例（`B3`/`D1` 两发，PASS **42**）⇒ 顶层 PASS 42→44＝＋2 那枚 delta 由**同树**量出（⛔ 两树相减）；旁证一把＝`^func Test` 计数 `5480434f`=43→HEAD=**45** |
| 三数带尺名 | `^--- PASS` 顶层＝44；`^    --- PASS` 含子测试＝12；两把⛔ 相加（票面那条老规矩） |
| 越界逐笔 | 8 笔 `git show --name-only` 写面外各 **0**；并集 80 枚＝`internal/audio/**` 2 ＋ `probes/306/r1/**` 77 ＋ 票 306 文件 1 ＋ 写面外 **0**（显式 SHA 列，⛔ `/tmp` 宽 glob） |
| 禁改清单 | 7 枚逐枚 `git diff --quiet 5480434f HEAD -- <f>`：6 枚 **rc=0**，第 7 枚（`wave_format_float_300_windows_test.go`）rc=1＝**它确实变了**，变的只有 `AC#3` 那 16/5 枚注释行；`frontend design third_party` 区间命中 **0 枚** |
| 票面纪律 | 腿那笔票面 diff＝**26 加 0 删**（append-only），框普查＝**1 勾／5 未勾**（＝本腿进场现量，⛔ 腿翻过任何东西） |
| 被跟踪 `.wav` | **0**（`git ls-files \| grep -icE '\.wav$'`＝0，`rc=1`）⇒ ⛔ 开甲形那处先例 |

**附的那两枚条件（欠的都归编排者车道，本腿⛔ 跑）**

1. `GOFLAGS= go build ./...` 票面那句要的是**这台机器上那枚台面**的一发；本腿只取到**导出树**那一发（`rc=0`），
   因为整包（含 `./cmd/wisp/`）在派单里被划成编排者车道 ⇒ **欠：共享工作树整包 `GOFLAGS= go build ./...` 一发 `rc`**。
2. 票面 `AC#4` 只写"红名册成对差集"，腿自愿多交了一把 `bash scripts/portable-tests.sh --scope=census`（totals 那行逐字未变）；
   本腿**⛔ 复跑 census**——那把尺会驱动跨包 `go test`（含 `cmd/wisp` 面，且缺 sherpa PATH 时退 `0xc0000135` 而零 `--- FAIL`）＝编排者车道
   ⇒ **欠：census 那一发的 totals 行＋`rc`，由编排者在自己的车道里复认**（腿的件在 `logs/40-census-raw.txt`，本腿⛔ 采信它的读数）。

⇒ 这两枚之内的所有尺本腿都跑到了 ⇒ 判 **附条件成立**（⛔ `不成立`：条件那两发⛔ 是本票射程里的"没牙"，是车道划分）。
