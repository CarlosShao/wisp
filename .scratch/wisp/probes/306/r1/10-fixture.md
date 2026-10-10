# 306-r1 · 10 夹具设计（`AC#1`，乙形）

台面＝工作树（`internal/audio` 起手 clean，见 `00-anchor.md`）。新文件一枚＝
`internal/audio/wavinjector_extensible_float_306_test.go`（package `audio`，**⛔ `//go:build windows`**，
⛔ 任何二进制 `.wav` 夹具、⛔ `testdata/`）。

## 1. 选形＝表③ 的乙形（⛔ 甲形）

| 判据 | 乙形（本腿走的） | 甲形（⛔ 本腿没走） |
|---|---|---|
| 被跟踪的 `.wav` 枚数 | 0（不变） | 会开第一枚先例（`git ls-files \| grep -icE '\.wav$'`＝0，306-a1 表③ 与编排者复跑同数） |
| `internal/audio/testdata/` | ⛔ 新增（audio 包今天没有） | 会开第二处先例 |
| 夹具字节可否逐字节注释 | 可（写在代码里，`sed` 可指认） | 二进制，"改期望侧要再生文件" ⇒ 要多一把尺 |
| CI 可见面 | **两档**（core=ubuntu `portable-tests.sh:242` ＋ windows `:253`，⛔ 动 pin：`core_pin:168`／`win_pin:195` 都是**包名**） | 同（票面 `:25` 只写了 windows 档，编排者第 5 节已把射程放宽成两档） |
| 标签 | ⛔ 无标签 ⇒ `GOOS=linux go vet ./internal/audio/` 那一面连测试侧都吃得到 | 同 |

入口两形都用了（派单允许"入口两形都行"，本腿**同时**装两枚，因为突变凭据要"指名用例当场红"，两枚各给一句独立红句）：

- `TestParseWavExtensibleFloat32306`＝同包直调 `parseWav(raw)`，先例逐字＝`internal/audio/hotplug_test.go:581`。
- `TestWavInjectorExtensibleFloat32306`＝`NewWavInjector(path)`，面写进 `t.TempDir()`（与 `wavinjector_test.go` 家族的习惯一致）。
  这一枚刻意用 **16k mono**，使构造期 `MonoDownmix(samples,1)` 与 `rate != TargetRate` 那两步都是 no-op
  （`wavinjector.go:65-68`、`resample.go:102-105`）⇒ 断的是解析出来的那条流本身，⛔ 掺进重采样与混叠的算术。

## 2. 那 40 字节的 `fmt` body 怎么写（偏移都相对 body 起点，即 `parseWav` 里的 `body`）

| 偏移 | 字节 | 写什么 | 为什么 |
|---|---|---|---|
| +0 | u16 | **`0xFFFE`** | `WAVE_FORMAT_EXTENSIBLE`＝`:192` 那枚行内权威值的期望侧 |
| +2 | u16 | `chans` | 断言里要检查 `chans`（48k 那枚面＝2） |
| +4 | u32 | `rate` | 断言里要检查 `rate`（48000／16000） |
| +8 | u32 | `rate*blockAlign` | 布局真实性，`parseWav` ⛔ 读 |
| +12 | u16 | `chans*4` | 同上（float32 ⇒ 4 字节一通道） |
| +14 | u16 | **`32`** | `bits==32` 是 `:218` 那枚 case 的第二半 |
| +16 | u16 | **`22`** | `cbSize`＝`WAVEFORMATEXTENSIBLE` 在 18 字节 `WAVEFORMATEX` 之后带的扩展字节数（`Samples` 2 ＋ `dwChannelMask` 4 ＋ `SubFormat` 16 ＝ 22）⇒ **18+22＝40**，与 `:193` 那句 `size < 40` 正好对齐：面⛔ 报错，又⛔ 大到能被逐字节注释 |
| +18 | u16 | `32` | `Samples.wValidBitsPerSample`，布局真实性 |
| +20 | u32 | `dwChannelMask`（立体声＝3，单声道＝0x4） | `parseWav` ⛔ 读；写进注释里说破，免得下一枚腿以为它被钉着 |
| **+24** | u32 | **`0x00000003`** | `KSDATAFORMAT_SUBTYPE_IEEE_FLOAT` 的 `Data1` ⇒ `:196` 读的正是 `body+24 : body+26` 这**两个低字节**＝`03 00`＝`3` ★这一枚是三处牙的共同支点：偏移若被写歪（+26），读出来是 `00 00`＝`0` |
| +28..39 | 12 B | GUID 尾部 `Data2/Data3/Data4`（真值 0x0000／0x0010／`80 00 00 AA 00 38 9B 71`） | 让面是一枚**逐字节诚实**的 `WAVEFORMATEXTENSIBLE`，⛔ 编造 |
| `data` 块 | u32 size ＋ N×4 字节 | `math.Float32bits` 逐枚小端 | 样本编码面；⛔ 任何真文件、⛔ 任何设备 |

`RIFF` 那枚 size 字段按 `4 + (8+40) + (8+dataLen)` 写正确值（`parseWav:166` 只验 `"RIFF"`/`"WAVE"` 两枚 magic、⛔ 读它，注释里写破了）。
`fmt` 块 size=40 与 `data` 块 size=32 都是偶数 ⇒ 词对齐补位（`:202` 的 `size & 1`）＝0，⛔ 引入第二处不确定性。

期望侧值**写在测试文件里**（`w306TagExtensible`／`w306SubFormatData1`／`w306CbSize`／`w306BitsFloat32`），⛔ 引用 `wavinjector.go`
里的行内字面量、⛔ 引用 `wasapi_windows.go` 那三枚常量（形状照票 300 `AC#6`："常量做期望侧的话它⛔ 能拒绝一枚错的产码值"）。

## 3. 断言的是**解析出来的值**（⛔ 断"分支被调用过"）

⛔ 计数器、⛔ coverage 钩子、⛔ mock、⛔ 任何"某支跑过没跑过"的问句。四枚断言全是**值**：

1. `err == nil`（`t.Fatalf`，红句里点名"extensible tag ＋ float SubFormat 必须走到 float32 支"）。
2. `rate == 48000`（单声道那枚面＝`TargetRate`）、`chans == 2`。
3. **`len(samples) == 8`**——这把尺专门区分"读成 float32"与"同样 32 字节读成 16 枚 PCM16"（后者会给 16）。
4. 逐枚 `samples[i] == w306WantI16[i]`，`want` 是**手推**的：
   `1.0→32767 / -1.0→-32768 / 0.5→16384 / -0.5→-16384 / 0.25→8192 / -0.25→-8192 / 0.125→4096 / -0.125→-4096`，
   推导逐枚写在文件注释里（`resample.go:119-129` 的映射：正数 `int64(v*32767+0.5)`、负数 `int64(v*32768-0.5)`、再 `clampI16`）。
   ★手推⛔ 调用 `FloatToPCM16` ⇒ 期望侧⛔ 从本包实现里取（同 `AC#6` 的纪律）；`0／1.0／-1.0／0.25` 四枚与既有
   `resample_test.go:168-172`（`TestMonoDownmixAndFloatConvert`）逐字同数，两把尺互相印证而非同一把。

## 4. 一枚夹具怎么同时走到 `:192` → `:196` → `:218`

`fmtTag`（`:188`）＝`0xFFFE` ⇒ `:192` 条件**真** ⇒ `:193` 验 `size < 40`（面的 size＝40 ⇒ ⛔ 报错、继续）⇒
`:196` 把 `fmtTag` **重写**为 `data[body+24:body+26]`＝`3` ⇒ `switch`（`:210`）里 `:211`（`1&&16`）⛔ 中、
`:218`（`3 && bits==32`）**中** ⇒ `:219-223` 按 4 字节步长解 float32 ⇒ `:224` `FloatToPCM16` ⇒ 返回。

三枚突变各自的落点（推演＝306-a1 表③ 与编排者块读数同向；凭据在 `20-nails.md`）：

- `:192` 的 `0xFFFE`→`0xFFFD` ⇒ 条件假 ⇒ `fmtTag` 留在 `0xFFFE`＝65534 ⇒ `:211`/`:218` 全不中 ⇒ `:226` 报错。
- `:196` 的 `body+24`→`body+26`（那一处偏移，连同它的配对上界一起移，见 `20-nails.md` 丙那一发的字面写法与理由）⇒ 读到 GUID 里 `Data1` 的高两字节＝`0` ⇒ 同样落 `:226`。
- `:218` 的 `3`→`2` ⇒ `:196` 重写出的 `3` ⛔ 任何 case ⇒ 落 `:226`。

⇒ 三枚都让 `err != nil` ⇒ 第一枚断言（`t.Fatalf`）当场红，红句里带的正是那枚 `tag=` 数字，三枚读数⛔ 同。

## 5. 刻意⛔ 装的东西（本腿自陈，免得被读成"顺手补全"）

- `:193-195` 的 `size < 40` 错误支：面是 40 ⇒ `:194` 仍**零执行者**（＝票面"⛔ 顺手补 injector 其它分支"）。
- 裸 `tag=3／bits=32`（非 extensible）面：会单独给 `:218` 装一枚执行者而让 `:192`／`:196` 继续没牙 ⇒ ⛔ 装。
- PCM16 支与 default 支：既有夹具已有执行者（306-a1 表① 乙类的块 hit=1 读数）。
- ⛔ 动 `wavinjector.go` 一字（本票只买尺）；⛔ 动 `wasapi_windows.go:378`（R2＝票 300 形 C5）；⛔ 动票 300 两枚钉的断言。

## 6. 现态读数（工作树，本腿自跑，⛔ 门禁那两把尺在 `40-gates.md`）

| 尺 | 读数 | 件 |
|---|---|---|
| `gofmt -l internal/audio/wavinjector_extensible_float_306_test.go`（写后第一发） | 列出该文件＝格式⛔ 净 | `logs/10-gofmt-before.txt` `rc_gofmt_l=0` |
| 同尺，`gofmt -w` 之后 | **0 行** | 同上 `rc_gofmt_after=0`（`rc_gofmt_w=0`） |
| `GOFLAGS= go test ./internal/audio/ -count=1 -run '306' -v` | `--- PASS: TestParseWavExtensibleFloat32306`／`--- PASS: TestWavInjectorExtensibleFloat32306`／`ok … 0.026s`／**`rc_targeted=0`** | `logs/10-targeted-green.txt` |
