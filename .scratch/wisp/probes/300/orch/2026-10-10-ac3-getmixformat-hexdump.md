# 票 300 `AC#3`〔仅本机可量、归编排者〕：这台机器真机混音格式的前 40 字节

跑＝编排者本人（⛔ 腿）。时刻 `2026-10-10 13:53:37 → 13:53:38 +0800`（尺＝`date` 的 stdout，起与终各一发）。
HEAD 当时＝`849ce6e96a0fbf2fbfe38b22cda9def50b528` 的**前一枚**读数见下面「provenance」那节，逐号列在那儿。

## 0. 一句话结论

这台机器上四枚默认端点的混音格式**都是** `WAVE_FORMAT_EXTENSIBLE`(65534) ＋ SubFormat＝`KSDATAFORMAT_SUBTYPE_IEEE_FLOAT`
（完整 16 字节 GUID 逐字节对得上），而旧代码从 **26** 取 `Data1` 落进的那两个字节是 `00 00` ⇒
`tag=0` ⇒ `floating=false`。票 297 三支里的 ③（"切错字节"）到这一行拿到一枚**真机上可复现的形状**。

## 1. 射程（这发⛔ 是什么）

只读字节面。全程**没有**`IAudioClient::Initialize`、**没有**`Start`、**没有**打开麦克风、**没有**任何出声、
零窗口、零弹窗。取到的那块内存由 COM 分配、用 `CoTaskMemFree` 归还（`defer` 两处都在载具里）。

⚠ 这发**⛔ 证明**的事，逐条列清楚，免得下一程把它读大：
- ⛔ 证明采集链路能跑通／球能进干活进程（那是票 228 那一族的活）。
- ⛔ 证明 CI 现在会跑到这一族用例（票 301 `AC#5` 那枚凭据仍未闭合，见台账 `A804`）。
- ⛔ 证明别的机器／别的驱动／别的采样率也是这样——**只这一台、只这一刻、四枚端点**。
- ⛔ 证明生产 `parseWaveFormat` 改完之后整条路就对了；它只把"字节面上那 2 字节之差"钉成真事。

⚠ 载具读的是**固定 40 字节**（＝`sizeof(WAVEFORMATEXTENSIBLE)`），它**没有**先看 `cbSize` 再决定读多少。
在这台机器上这是安全的（`wFormatTag=65534` 且 `cbSize=22` ⇒ 40 字节全在块内）；
若端点回的是普通 `WAVEFORMATEX`（`cbSize=0`、块长 18），这发载具就是越界读——
**那是载具的边界，⛔ 生产的**（生产那条分支先判 `tag==extensible && cbSize>=22` 才取 GUID）。
下一程如果要把这形转成常驻仪器，这条得先修。

## 2. provenance（尺口径，逐字）

- 载具＝**仓外导出树**，不在仓里、⛔ 入库：`C:\Users\swq\AppData\Local\Temp\wisp300-hex\internal\audio\zz_orch300_mixformat_windows_test.go`
  （`//go:build windows`、`package audio`、只加这一枚 `_test.go` 到导出树）。
- 命令逐字（`cd` 到导出树根）：
  `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./internal/audio/ -run TestOrch300MixFormatHexdump -v -timeout 180s`
  ⇒ **`rc=0`**、`--- PASS`、`ok github.com/CarlosShao/wisp/internal/audio 0.095s`（那一行 `rc=N` 就在这件末尾，⛔ 0 字节）。
  ⚠ harness 那句 PATH 对本包**实测是多余保险**（`grep -rln 'sherpa\|onnxruntime' internal/audio` ⇒ 0 命中，见 `A805`），照抄只为不与既有配方分叉。
- 跑的是哪一版码（⛔ 靠"应该是新的"，靠尺）：
  `cmp` 导出树的 `internal/audio/wasapi_windows.go` ↔ `git show HEAD:internal/audio/wasapi_windows.go` ⇒ **IDENTICAL**，
  两边同一行都是 `:211 sub := *(*windows.GUID)(unsafe.Add(p, 24))` ⇒ **这发跑的是已经改到 24 的那一版**（`028529fc` 起）。
  对拉用的 HEAD 号＝起笔现量 `849ce6e96a0fbf2fbfe38b22cda9def50b528`（`git log -1 --format=%H`；导出树文件时间戳 13:41，
  其 `wasapi_windows.go` 与 HEAD 逐字节同）。
- 端点取法判给**生产自己**的写法（见下一节）。
- 现量机器状态：`Get-Process wisp`＝0、`Get-Process balldebug`＝0；CPU 37%→39%、MEM 约 54.7%→56%（跑前跑后各一发，都在加派线以下）。

## 3. 我自己踩的三形（具名留痕，⛔ 抹）

前两形都**⛔ 是**有效读数，原因如下——写下来是因为每一形都可能被下一程照抄：

1. **槽位判错**：我第一／第二形把 `IMMDeviceEnumerator::GetDefaultAudioEndpoint` 按 **slot 3** 调。
   实测 `(flow, 0, &out)` 回 `0x80070057`（E_INVALIDARG），`(flow, 1|2, &out)` 回 `S_OK` 而 `out` 为空
   ——后者的语义正是 **`EnumAudioEndpoints` 的 `stateMask`**（0 非法／1=ACTIVE），也就是我拿 `EnumAudioEndpoints` 的参数形状去调了另一枚方法。
   ⇒ **判给生产写法**：本仓 `internal/audio/mmdevice_windows.go:85` 逐字 `comCall(enum, 4, ...)`，注释写着 `(slot 4)`。
   第三形照抄 slot 4 ⇒ 四枚端点 `hr=0x0`、`dev!=nil` 全取到。
   ★定式并回：**槽位这种"文档里没有的东西"，以本仓生产代码自己的写法为权威，⛔ 以我的记忆为权威。**
2. **写文件工具把载具规范化到了另一棵树**：第二形我写的路径落到了 `C:\tmp\wisp300-hex\internal\audio\...`（2,725 字节，13:46），
   而导出树真身在 `C:\Users\swq\AppData\Local\Temp\wisp300-hex\...`（3,110 字节，13:50）。
   那枚 `C:\tmp` 里的 `.../go.mod` 根本不存在 ⇒ 第二形**从未被编译过**，谈不上读数。
   ★定式并回：写完**立刻 `ls -la` 那一条路径**，⛔ 相信回执（`A` 账里"交付判据是 `wc -c`，不是回执"同族，这次踩的是**绝对路径被改写**那一支）。
3. **打印循环未夹长度**：`b[i:i+16]` 在 `i=32` 时 slice 越界 panic ⇒ 第三形加了 `hi := i+16; if hi > len(b) { hi = len(b) }`。

## 4. 原始读数（逐字，四枚端点各一节；节内三行 hexdump 四枚端点**逐字节相同**）

```
=== RUN   TestOrch300MixFormatHexdump
    zz_orch300_mixformat_windows_test.go:77: === render-console === GetDefaultAudioEndpoint(slot4) hr=0x0 dev=true
    zz_orch300_mixformat_windows_test.go:79:   parseWaveFormat -> tag=3 channels=2 rate=48000 bits=32 err=<nil>
    zz_orch300_mixformat_windows_test.go:79:   byte 00-15: fe ff 02 00 80 bb 00 00 00 dc 05 00 08 00 20 00
    zz_orch300_mixformat_windows_test.go:79:   byte 16-31: 16 00 20 00 03 00 00 00 03 00 00 00 00 00 10 00
    zz_orch300_mixformat_windows_test.go:79:   byte 32-39: 80 00 00 aa 00 38 9b 71
    zz_orch300_mixformat_windows_test.go:79:   tag@0=65534 ch@2=2 rate@4=48000 bits@14=32 cbSize@16=22 Samples@18=32 mask@20=3
    zz_orch300_mixformat_windows_test.go:79:   Data1@24=3 | lowWORD@24=3 lowWORD@26=0
    zz_orch300_mixformat_windows_test.go:79:   floating?(tag==3)=true  旧 26 那形会得 tag=0
    zz_orch300_mixformat_windows_test.go:77: === capture-console === GetDefaultAudioEndpoint(slot4) hr=0x0 dev=true
    zz_orch300_mixformat_windows_test.go:79:   parseWaveFormat -> tag=3 channels=2 rate=48000 bits=32 err=<nil>
    zz_orch300_mixformat_windows_test.go:79:   byte 00-15: fe ff 02 00 80 bb 00 00 00 dc 05 00 08 00 20 00
    zz_orch300_mixformat_windows_test.go:79:   byte 16-31: 16 00 20 00 03 00 00 00 03 00 00 00 00 00 10 00
    zz_orch300_mixformat_windows_test.go:79:   byte 32-39: 80 00 00 aa 00 38 9b 71
    zz_orch300_mixformat_windows_test.go:79:   tag@0=65534 ch@2=2 rate@4=48000 bits@14=32 cbSize@16=22 Samples@18=32 mask@20=3
    zz_orch300_mixformat_windows_test.go:79:   Data1@24=3 | lowWORD@24=3 lowWORD@26=0
    zz_orch300_mixformat_windows_test.go:79:   floating?(tag==3)=true  旧 26 那形会得 tag=0
    zz_orch300_mixformat_windows_test.go:77: === render-multimedia === GetDefaultAudioEndpoint(slot4) hr=0x0 dev=true
    zz_orch300_mixformat_windows_test.go:79:   parseWaveFormat -> tag=3 channels=2 rate=48000 bits=32 err=<nil>
    zz_orch300_mixformat_windows_test.go:79:   byte 00-15: fe ff 02 00 80 bb 00 00 00 dc 05 00 08 00 20 00
    zz_orch300_mixformat_windows_test.go:79:   byte 16-31: 16 00 20 00 03 00 00 00 03 00 00 00 00 00 10 00
    zz_orch300_mixformat_windows_test.go:79:   byte 32-39: 80 00 00 aa 00 38 9b 71
    zz_orch300_mixformat_windows_test.go:79:   tag@0=65534 ch@2=2 rate@4=48000 bits@14=32 cbSize@16=22 Samples@18=32 mask@20=3
    zz_orch300_mixformat_windows_test.go:79:   Data1@24=3 | lowWORD@24=3 lowWORD@26=0
    zz_orch300_mixformat_windows_test.go:79:   floating?(tag==3)=true  旧 26 那形会得 tag=0
    zz_orch300_mixformat_windows_test.go:77: === capture-multimedia === GetDefaultAudioEndpoint(slot4) hr=0x0 dev=true
    zz_orch300_mixformat_windows_test.go:79:   parseWaveFormat -> tag=3 channels=2 rate=48000 bits=32 err=<nil>
    zz_orch300_mixformat_windows_test.go:79:   byte 00-15: fe ff 02 00 80 bb 00 00 00 dc 05 00 08 00 20 00
    zz_orch300_mixformat_windows_test.go:79:   byte 16-31: 16 00 20 00 03 00 00 00 03 00 00 00 00 00 10 00
    zz_orch300_mixformat_windows_test.go:79:   byte 32-39: 80 00 00 aa 00 38 9b 71
    zz_orch300_mixformat_windows_test.go:79:   tag@0=65534 ch@2=2 rate=48000 bits=32 err=<nil>
    zz_orch300_mixformat_windows_test.go:79:   tag@0=65534 ch@2=2 rate@4=48000 bits@14=32 cbSize@16=22 Samples@18=32 mask@20=3
    zz_orch300_mixformat_windows_test.go:79:   Data1@24=3 | lowWORD@24=3 lowWORD@26=0
    zz_orch300_mixformat_windows_test.go:79:   floating?(tag==3)=true  旧 26 那形会得 tag=0
--- PASS: TestOrch300MixFormatHexdump (0.04s)
PASS
ok  	github.com/CarlosShao/wisp/internal/audio	0.095s
rc=0
```

⚠ **上面那一块我⛔ 当作"纯逐字"卖**：第四节是从终端回显转写的，四枚端点里**前三枚**的逐行形式我按回显照抄，
第四枚（`capture-multimedia`）我转写时把 `parseWaveFormat ->` 那行误排了两遍（第二遍其实是 `tag@0=...` 那行）。
⇒ **本件的凭据是下面那张解码表＋四枚端点三行 hexdump 逐字节相同这一条**（这条我在回显里逐字节看过：`fe ff 02 00 …` 三行四组全同）；
要一份⛔ 带转写瑕疵的原始件，就按第 2 节那条命令重跑一发、把 stdout 直接 `>` 到文件（我这次是先看了屏再落笔，⛔ 让"已经看过"冒充"已经存了"——这正是台账里"交付判据是 `wc -c`"那一族的反面）。

## 5. 解码对照（这 40 字节按权威布局读）

权威＝盘上 SDK 头 `C:\Program Files (x86)\Windows Kits\10\Include\10.0.26100.0\shared\mmreg.h`（该段被 `pshpack1.h` 包着 ⇒ 1 字节对齐，无填充）。
`WAVEFORMATEXTENSIBLE` ＝ `WAVEFORMATEX`(18) ＋ `Samples@18`(WORD union) ＋ `dwChannelMask@20`(DWORD) ＋ `SubFormat@24`(GUID 16) ＝ 40。

| 偏移 | 字节（四枚端点同） | 字段 | 读出来的值 |
|---|---|---|---|
| 0 | `fe ff` | `wFormatTag` | 65534 ＝ `WAVE_FORMAT_EXTENSIBLE` |
| 2 | `02 00` | `nChannels` | 2 |
| 4 | `80 bb 00 00` | `nSamplesPerSec` | 48000 |
| 8 | `00 dc 05 00` | `nAvgBytesPerSec` | 384000（＝48000×2×4，自洽） |
| 12 | `08 00` | `nBlockAlign` | 8（＝2×4，自洽） |
| 14 | `20 00` | `wBitsPerSample` | 32 |
| 16 | `16 00` | `cbSize` | 22 |
| 18 | `20 00` | `Samples.wValidBitsPerSample`（union） | 32 |
| 20 | `03 00 00 00` | `dwChannelMask` | 3 ＝ `SPEAKER_FRONT_LEFT｜SPEAKER_FRONT_RIGHT` |
| **24** | `03 00 00 00` | `SubFormat.Data1` | **3 ＝ `KSDATAFORMAT_SUBTYPE_IEEE_FLOAT`** |
| 26 | `00 00` | （旧代码读的那两个字节） | **0** ⇒ 旧形 `tag=0` ⇒ `floating=false` |
| 28–39 | `00 00 10 00 80 00 00 aa 00 38 9b 71` | `Data2`／`Data3`／`Data4` | 完整 GUID ＝ `{00000003-0000-0010-8000-00AA00389B71}` ＝ IEEE_FLOAT |

⇒ **这次⛔ 止于 `Data1`**：16 字节 GUID 与权威常量逐字节相等，比"只看头 4 字节"那一枚断言更硬。
入库用例 `internal/audio/parse_wave_format_300_windows_test.go` 现在钉的是 `Data1`（四形），全 GUID 那一枚⛔ 在本包——那是"要不要再加一枚"的活，⛔ 本发替它背书。

## 6. 这一发买到什么（对着票 300 的格子说）

- `AC#3` 那句"真机那枚混音格式到底是不是 extensible＋float"＝**是**，四枚端点（render/capture × console/multimedia）全是，同一份 40 字节。
- 与 `A805` 那发仓外五枚偏移扫描对拉：本机在 24 上四形全绿、在 26 上四形全红；
  这发真机读数**独立给出同一条分界**（24 处是真 GUID 的 3，26 处是 0）⇒ 两把⛔ 同一把尺（一把＝合成夹具、一把＝真机字节面），
  这正好是 `feedback-evidence-teeth` 第 8 条要的"一把尺两种形状都放行"的反面。
- ⛔ 翻框：`AC#3` 那一格仍归 `300-v2` 之后的非实现者裁；本件只是把〔仅本机可量〕那半枚欠的读数交上盘。
  票面上 `AC#3` 还挂着第二件事（与票 297 的关系），那一半⛔ 靠这发闭合。

## 7. 后续（记在这儿，免得丢）

1. 载具的**固定 40 字节、先看形状再决定读长**那一枚边界（第 1 节末尾）如果将来要转成常驻仪器必须先修。
2. `scripts/portable-tests.sh:593` 那枚 ledger 与票 301 的关系：本包 windows-tagged 用例今天进不了任何一档 CI（`A804`），
   ⇒ **这发真机读数的可复现性只有本机**，⛔ 任何"CI 会跑它"的暗示。

rc=0
