# 306-a1 · 10 表①：行内权威字面量全名册（`mmreg.h` 那族）

## 射程（本仓铁律：枚数不带尺名＝不算读数）

- 尺＝`grep -rnE --include=*.go '0xFFFE|0xfffe|65534|fmtTag == [0-9x]+|0x0003|0x00000003|0x00000001|data1: *[0-9]|bits == 3[2x]|waveFormat(Ext|Float|PCM)|KSDATAFORMAT|IEEE_FLOAT' internal cmd tools`
- 射程目录＝**`internal/ cmd/ tools/` 全部 `*.go`**（`tests/` 目录不存在＝尺报 `No such file or directory`，已把该词从射程里摘掉重跑，原始读数件＝`.scratch/wisp/probes/306/a1/logs/30-roster-worktree.txt`，72 行含自落的一行 `rc_roster=0/1`）。
- 量的是**工作树**；本锚点 `f994f95` 上 `git status --porcelain -- internal cmd tools | wc -l` ＝ **0**，且 `git show HEAD:internal/audio/wavinjector.go | cmp - internal/audio/wavinjector.go` ＝ **IDENTICAL** ⇒ **blob 与工作树两个名册在此锚点逐字相同**（两把尺并排＝同一份行号，这是读数不是假设）。
- 族定义（＝票 300 `AC#6` 钉的那族）：`WAVE_FORMAT_PCM 1`（`mmreg.h:2418`）· `WAVE_FORMAT_IEEE_FLOAT 0x0003`（`mmreg.h:2110`）· `WAVE_FORMAT_EXTENSIBLE 0xFFFE`（`mmreg.h:2376`）· GUID 低字 `0x00000003`/`0x00000001`（`KSDATAFORMAT_SUBTYPE_IEEE_FLOAT` / `_PCM`）。

## 甲类＝行内字面量 · ⛔ 引常量 · **且零执行者**（本票要装牙的就是这一类）

| 枚位 | 逐字（内容锚） | 牙（尺＝本腿复跑的覆盖率块图 `logs/40-cover.txt`） |
|---|---|---|
| `internal/audio/wavinjector.go:192` | `			if fmtTag == 0xFFFE { // WAVE_FORMAT_EXTENSIBLE: real tag is SubFormat[0:2]` | 块 `193.5,193.18`／`194.6,195.1`／`196.5,196.65` 全 **hit=0** ⇒ 改值⛔ 红 |
| `internal/audio/wavinjector.go:218` | `	case fmtTag == 3 && bits == 32:` | 块 `219.3,221.26`／`222.4,223.1`／`224.3,224.43` 全 **hit=0** ⇒ 改值⛔ 红 |

⇒ **枚数＝2**（与票面/派单给的数一致；尺名如上）。原始读数件＝`logs/31-roster-wavinjector.txt`。

**附赠第三枚（同一把尺量到，票面未列，AC#0 ① 的职责就是把它说话）**：

| `internal/audio/wavinjector.go:196` | `				fmtTag = binary.LittleEndian.Uint16(data[body+24 : body+26])` | **hit=0**（同上块 `196.5,196.65`） |

`:196` 不是那 4 枚"值"而是**权威偏移 24**（`WAVEFORMATEXTENSIBLE.SubFormat` 起点，`mmreg.h` 那份布局；`wasapi_windows.go:188-193` 注释里同一件事被钉过，`wavinjector.go` 这一份⛔ 牙）。⇒ 一枚 extensible 夹具能**同时**给 `:192`／`:196`／`:218` 三处装牙，本腿把它登记给编排者裁（票的 `AC#1`/`AC#2` 只承诺前两枚）。

## 乙类＝行内字面量 · **有**执行者（今天就有牙，⛔ 本票射程要修）

| 枚位 | 逐字 | 谁在读它 |
|---|---|---|
| `internal/audio/wavinjector.go:211` | `	case fmtTag == 1 && bits == 16:` | 块 `212.3,214.26`／`215.4,216.1`／`217.3,217.35` **hit=1** ⇒ 10+ 枚 PCM16 夹具读着它（表②） |
| `internal/audio/wavinjector.go:226` | `		return nil, 0, 0, fmt.Errorf("unsupported wav format: tag=%d bits=%d ...", fmtTag, bits)` | 块 `226.3,226.111` **hit=1** ⇒ default 支有执行者（`TestWavInjectorBadFile` 的 8-bit 支） |
| `internal/audio/parse_wave_format_300_windows_test.go:59` | `	le.PutUint16(b[0:], 0xFFFE)                                       // WAVE_FORMAT_EXTENSIBLE (mmreg.h:2376)` | 同文件自己的断言（`//go:build windows`，块级：`parseWaveFormat` 覆盖 **81.8%**，`logs/41-cover-func.txt`） |
| `internal/audio/parse_wave_format_300_windows_test.go:99/105/111/117` | `face300{channels: 2, rate: 48000, bits: 32, data1: 3, data1At: 24}` ／ `bits: 16, data1: 1` ／ `bits: 32, data1: 7` ／ `bits: 32, data1: 3, data1At: 26` | 同上文件钉着（GUID 低字 3/1/7 三形＋一枚"错偏移"负控） |
| `internal/audio/wave_format_float_300_windows_test.go:70/74/78` | `mmregWaveFormatPCM uint16 = 1` ／ `mmregWaveFormatIEEEFloat uint16 = 0x0003` ／ `mmregWaveFormatExtensible uint16 = 0xFFFE` | ★票 300 `AC#6` 的期望侧，`TestWaveFormatConstantsMatchMmregAuthority300`（`:147`）正读着 ⇒ **改这三枚会红**（＝那枚大钉装好的形状，本票要复制的就是它） |

## 丙类＝引常量（常量侧的权威，3 枚定义＋2 枚使用）

| 枚位 | 逐字 | 状态 |
|---|---|---|
| `internal/audio/wasapi_windows.go:97` | `	waveFormatPCM   = 1` | 声明；`AC#6` 附带覆盖面读着它（票 306 `AC#3` 硬边界⛔ 删它） |
| `internal/audio/wasapi_windows.go:98` | `	waveFormatFloat = 3` | 声明，`AC#6` 已钉 |
| `internal/audio/wasapi_windows.go:99` | `	waveFormatExt   = 0xFFFE` | 声明，`AC#6` 已钉 |
| `internal/audio/wasapi_windows.go:210` | `	if f.tag == waveFormatExt {` | 引常量，**有**执行者（`parseWaveFormat 81.8%`） |
| `internal/audio/wasapi_windows.go:378` | `	floating := s.format.tag == waveFormatFloat` | 引常量、**零执行者**：`Drain 0.0%`／`convertPacket 0.0%`（本腿复跑对上，`logs/41-cover-func.txt`）＝票面 **R2**，⛔ 本票射程（登记在 `A825`，要动需 owner 另批） |

⇒ **丙类"引了常量但没牙"的枚数＝1**（`:378`＝R2，本腿只报数⛔ 碰）。

## 丁类＝只在注释/文案里出现（⛔ 牙的载体，但 AC#3 那一格在这里）

`internal/audio/wavinjector.go:163-164`（doc comment 写着 "PCM16 (format 1) and IEEE float32 (format 3, incl. WAVE_FORMAT_EXTENSIBLE…)"——**注释比代码有牙，注释说的 extensible 支今天零执行者**）；`internal/audio/wasapi_windows.go:193-194`／`:212`（`KSDATAFORMAT_SUBTYPE_PCM is 1 and … IEEE_FLOAT is 3`）；`internal/audio/wave_format_float_300_windows_test.go:40`（`KSDATAFORMAT_SUBTYPE_IEEE_FLOAT is a different object and lives in ksmedia.h,` ← **`AC#3` 要改的就是这一句**，本腿只登记枚位与世代，改不改由落地腿按票面办）；`internal/audio/parse_wave_format_300_windows_test.go:101/107`（`why:` 文案里的 `Data1=0x00000003` / `0x00000001`）。

## 戊类＝同形异族（`0x0003`/`0x00000001` 字面撞车，**逐枚具名排除**，⛔ mmreg 族）

- `internal/tools/recycle_windows.go:37` ＝ `	foDelete = 0x0003 // FO_DELETE`（shell 文件操作码）
- `cmd/wisp/panel_host_windows_test.go:181` ＝ `		style:         0x0003, // CS_HREDRAW | CS_VREDRAW`（winuser 类样式）
- `cmd/wisp/notify_windows.go:27`/`:30`/`:35` ＝ `nimModify = 0x00000001` / `nifMessage = 0x00000001` / `niifInfo = 0x00000001`（Shell_NotifyIcon）

⇒ 排除 **6 枚**；否则"仓里还有多少处 `0x0003`"会被虚报成 8 枚。

## 表①结论（一行）

`mmreg.h` 那族里**"行内字面量且零执行者"＝2 枚**（`wavinjector.go:192` 的 `0xFFFE`、`:218` 的 `3`；外加 1 枚同处形状的权威偏移 `:196`），**"引常量但零执行者"＝1 枚**（`wasapi_windows.go:378`＝R2，⛔ 本票）；其余同族枚位要么已在常量侧、要么已被票 300 `AC#6` 的钉读着 —— **票 306 的洞族没有第三处散落在 `cmd/`、`tools/` 或别的 `internal/` 包里**（尺＝上面那条 grep，射程 `internal/ cmd/ tools/` 全 `*.go`，工作树＝blob）。
