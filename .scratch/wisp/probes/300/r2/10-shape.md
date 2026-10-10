# 300-r2 落点形状与逐枚断言（票 300 `AC#6`，形＝编排者裁的 C2）

新件＝`internal/audio/wave_format_float_300_windows_test.go`（`//go:build windows`，第 1 行逐字 `//go:build windows`，第 3 行 `package audio`）。
⚠ 本文件里"换 3→1 / 3→0"两列写的是**设计时的预期**；**实测读数在 `20-mutations.md`**，两列各带逐字红句。本件⛔ 用预期冒充读数。

## 名册尺（先证我没撞既有件）

尺＝`git ls-files internal/audio` ⇒ 既有 8 枚 `_test.go`：
`capturelevel_windows_test.go`／`captureopt_test.go`／`gate_test.go`／`hotplug_test.go`／`level_test.go`／`parse_wave_format_300_windows_test.go`／`resample_test.go`／`wavinjector_test.go`。
⇒ 我新增的 `wave_format_float_300_windows_test.go` **⛔ 与任何一枚同名**（`300-a2` 的 `F5` 那 8 枚）。

顶层用例名唯一性尺＝逐文件 blob `git show HEAD:<file> | grep -E '^func Test'`（本腿 18:4x 现跑）⇒ 42 枚顶层名里**无**
`TestWaveFormatConstantsMatchMmregAuthority300` ⇒ ⛔ 重名。
包级标识符碰撞尺＝新增的 4 枚标识符 `mmregWaveFormat*`／`r2Face`／`r2Parse`／`r2Authority*` 在 `internal/audio` 内
`grep -rn` 命中**只在我这枚新文件里**（既有件用的是 `face300` 与其 `build()` 方法，⛔ 同名，且本件⛔ 改它）。

## 八枚断言（1 枚顶层用例 ＋ 5 枚子测试 ＋ 子测试里再拆的独立判据，逐枚列）

| # | 断言名（子测试名 ⇔ 判据） | 观察侧来源 | 期望侧来源 | 设计上 3→1 | 设计上 3→0 |
|---|---|---|---|---|---|
| A1 | `ieee_float_low_word`：`waveFormatFloat ⇔ 0x0003` | 产码常量 `wasapi_windows.go` 声明的 `waveFormatFloat`（本文件 `got: waveFormatFloat`） | 测试自带字面 `mmregWaveFormatIEEEFloat = 0x0003`（注释具名抄自 `mmreg.h:2110` 那一行原文） | 红 | 红 |
| A2 | `extensible_tag`：`waveFormatExt ⇔ 0xFFFE` | 产码常量 `waveFormatExt` | 测试自带字面 `mmregWaveFormatExtensible = 0xFFFE`（`mmreg.h:2376`） | 不红（与 float 无关） | 不红 |
| A3 | `pcm_tag`：`waveFormatPCM ⇔ 1` | 产码常量 `waveFormatPCM` | 测试自带字面 `mmregWaveFormatPCM = 1`（`mmreg.h:2418`；权威那行写的就是 `1`，⛔ 是 `0x0001`，按那一行逐字） | 红（与 A4 同时） | 红（与 A4 同时） |
| A4 | `the_three_constants_map_one_to_one_onto_the_authority_table` 第一支：三枚常量⛔ 两两相等 | 三枚产码常量组成的切片（⚠ **运行期切片**，⛔ 用常量做 map 字面量的键——那会让 3→1 变成**编译错误**而⛔ 是一句红，本腿实测过这个形状再改掉，见下"造过又拆掉的一形"） | 测试自带的 `mmreg.h` 三值互不相同这一事实（字面 1／`0x0003`／`0xFFFE`） | 红（float 与 pcm 撞值） | 红（float 与两者都不等 → 走 A4 第二支） |
| A5 | 同上前缀 第二支：每枚常量的值必须落在权威三值之内且**名字对上** | 同上 | `r2Authority` 表（三行都是测试自带的字面 ⇔ 名字） | 红（float 的 `1` 在权威表里登记的名字是 `waveFormatPCM` ⇒ 名字不匹配那一句） | 红（`0` ⛔ 在权威三值里 ⇒ "none of WAVE_FORMAT_PCM/IEEE_FLOAT/EXTENSIBLE" 那一句） |
| A6 | `parseWaveFormat_expands_a_face_tagged_with_the_authority_extensible_value` | `parseWaveFormat` 返回的 `tag`（真产码解析路径） | 测试自带字面 `0x0003`；夹具的 tag 也写字面 `0xFFFE`（⛔ 写 `waveFormatExt`） | 不红（`parseWaveFormat` ⛔ 用 `waveFormatFloat`，F1 那把尺） | 不红 |
| A7 | 同上红句的用途：钉 `waveFormatExt` 的**行为面**（⛔ 只钉声明值） | 同上：改 `waveFormatExt` ⇒ `if f.tag == waveFormatExt` 那支关死 ⇒ 返回的 `tag` 停在 `0xFFFE` | 同上 | 不红 | 不红（`waveFormatExt` 挪才红） |
| A8 | `control_plain_pcm_tag_is_not_expanded`（A6/A7 的正控） | `parseWaveFormat` 返回的 `tag`，夹具 tag＝字面 `1`、`cbSize`＝字面 0 | 测试自带字面 `1` | 不红 | 不红 |

**本格真正的牙在 A1＋A4＋A5**：三枚都只靠"测试自带字面 ⇔ 经常量拿到的值"，两侧同枚常量＝本格作废这一条**没被踩**（期望侧那三行⛔ 出现 `waveFormatFloat`/`waveFormatExt`/`waveFormatPCM` 任何一枚，尺＝本腿对 `mmregWaveFormat` 那三枚 `const` 的赋值行 `grep` ⇒ 右侧只有字面量）。

## 造过又拆掉的一形（具名留痕，⛔ 抹）

A4 的第一版是 `map[uint16]string{waveFormatPCM: …, waveFormatFloat: …, waveFormatExt: …}`。
⇒ **枚枚键都是常量**，把 `waveFormatFloat` 改成 `1` 之后 `waveFormatPCM` 与 `waveFormatFloat` 是**同一枚常量键** ⇒ Go 前端报 duplicate key **编译错误**，`go test` 拿 `rc=1` 但**一句 `--- FAIL` 都没有、更没有点名那枚常量的红句**。那不是凭据，那是一枚⛔ 响的仪器。
⇒ 改成运行期切片＋两两比较（即上表 A4/A5 的形状），并先跑一发 `go vet ./internal/audio/` 证它编得过。
⚠ 这一形的教训＝**"⛔ 动既有件"之外还要"⛔ 把突变跑成编译错误"**：突变攻击要的读数是一句具名红，⛔ 是一个非零 rc。

## 观察侧到不了生产那一行（具名，⛔ 造假日为面）

`waveFormatFloat` 在产码里唯一的消费者＝`internal/audio/wasapi_windows.go` 的
`floating := s.format.tag == waveFormatFloat`，它在 `(*wasapiStream).Drain()` 体内，**要一条活流＝真设备**；
票 300 票面第 31 行把真设备／真麦克风那一发归编排者、⛔ 归腿。
`convertPacket` 收的是 `floating bool` **参数**，自己去调它证的是参数、⛔ 是那枚常量 ⇒ 本腿**没有**造那一形。
⇒ 本格能买到的最硬一句＝**那枚常量的值 ⇔ 测试自带的 `mmreg.h` 字面**；"生产真的这么判"仍⛔ 可观察＝编排者裁语里具名的残余 **R2**（要买它得走 C5＝改产码语义面，本票⛔ 授权）。

## ⛔ 删 `waveFormatPCM`

`F2` 那句"声明后零使用"是**词面尺**读数。本件按裁语 ③ 把它钉上（A3），**⛔ 删它、⛔ 改它的值、⛔ 顺手清理**。
"零使用的那枚该不该删"是另一格、归后续验收裁。

## 权威字面量的出处（⛔ 混写两枚对象）

- 三枚 **WORD 值**＝`mmreg.h:2110`（`0x0003`）／`:2376`（`0xFFFE`）／`:2418`（`1`），逐字原文＝编排者本机读数存档
  `.scratch/wisp/probes/300/orch/2026-10-10-mmreg-authority.txt`。⚠ 那是**本机** SDK 路径，托管 runner 上⛔ 保证有该头文件 ⇒ 本用例**⛔ 读盘、⛔ `t.Skip`**，权威值只以**注释＋字面量**的形式抄进测试，来源在那段注释里具名为"copied out of the SDK header"。
- `KSDATAFORMAT_SUBTYPE_IEEE_FLOAT` 那枚 **GUID 在 `ksmedia.h`、⛔ 在 `mmreg.h`**。本件注释里⛔ 两处提到它，都没把 GUID 与 WORD 值写成同一句出处；
  `AC#0` 用 `mmreg.h` 钉的是**偏移**，本格钉的是**值**，同源⛔ 同物。

## 分母影响（裁语 ④ 要具名那一枚）

本件新增**一枚顶层用例**（`--- PASS` 那把尺 1 枚；源码 `^func Test` 那把尺也是 1 枚），内部 5 枚子测试。
顶层名＝`TestWaveFormatConstantsMatchMmregAuthority300` ⇒ audio 进 `windows)` 档之后（`0a0f62ef`），
`test-windows` 那一档的红名册作差面多出的就是**这一枚**。⚠ 本腿**没有**跑托管 CI，本机两发整包对照的读数在 `30-gates.md` 门③。
