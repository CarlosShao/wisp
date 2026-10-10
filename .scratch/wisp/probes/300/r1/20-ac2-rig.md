# 300-r1 — `20` `AC#2` 落地那一发（ⓒ 入库用例＝**一进门就绿**，四形各自读数）

落点＝`internal/audio/parse_wave_format_300_windows_test.go`（新增，`//go:build windows`，包内 `audio`）。
形状＝**逐形沿用** `300-a1r` 交付、`300-v1` 重跑并扫过五枚偏移的那套夹具（`probes/300/a1r/20-ac1-rig.md` §1 代码块），
⛔ 另造字节布局；只有头顶注释从"⛔ 入库的仓外台件"改写为"随修入库"，并把 `300-v1` §6.2 那句"S4 是四枚里最弱的一枚"写进了注释。

## 1. 夹具的权威形状（与 `AC#0` 逐格对得上，⛔ 来自本仓任何一枚实现）

`b[0..1]=0xFFFE` · `b[2]=channels` · `b[4]=rate` · `b[8]=avgBytes` · `b[12]=blockAlign` · `b[14]=bits` ·
**`b[16]=22`（cbSize）** · `b[18..19]=bits`（写进 union `Samples`）· **`b[20..23]=dwChannelMask`** ·
`SubFormat.Data1` 由 `data1At` 决定写在哪（S1/S2/S3＝**24**，S4＝故意 **26**）· 面长 **40**。
尺＝我这一棵树上 `grep -n "le.PutUint" internal/audio/parse_wave_format_300_windows_test.go`（逐字输出见 `logs/rig-writes.txt`）。

四形断的都是**解析出来的那个值**（`tag`，另加 `header` 三枚与派生的 `floating`），⛔ 断"被调用过"；
⛔ 真设备、⛔ 真麦、⛔ `Initialize`、⛔ `Start`。两把尺（逐字输出＝`logs/rig-writes.txt`）：

```
$ grep -n "Initialize\|Start\|GetMixFormat\|comCall\|wasapiOpener" internal/audio/parse_wave_format_300_windows_test.go
31:// Initialize, no Start.

$ grep -nE "(comCall|CoCreateInstance|GetMixFormat|wasapiOpener|\.init\(|deviceStream)" internal/audio/parse_wave_format_300_windows_test.go | wc -l
0
```

⚠ **第一把我自己踩的坑，具名留痕不删**：先前我在这一格写"读数 **0**"，用的尺是
`grep -c "Initialize\|Start\|GetMixFormat\|deviceStream"` ⇒ 实际给的是 **1**，命中＝第 31 行那**句注释**（"No / Initialize, no Start." 被折行写在注释里）。
⇒ 名字类尺会把**散文**算进去，"没用设备"这一判只能数**调用形状**（第二把，读数 **0**，命中面＝函数调用与标识符，⛔ 注释）。
这条与 `tools/d22scan` 的"注释豁免／字符串不豁免"是同一族仪器事。

## 1.1 "⛔ 另造一套字节布局"的正面凭据（与 `300-a1r` 那枚台件对拉）

台原件＝`/tmp/wisp300-a1r/tree/internal/audio/parse_wave_format_300_windows_test.go`（**6172 字节**，腿留着没删）。
尺＝两枚文件各剥掉行尾注释（`sed 's/[[:space:]]*\/\/.*$//'`）再 `diff`，整块逐字＝`logs/diff-vs-a1r-nocomment.txt`，结果只剩三类：

```
26a27,32      （块注释被剥成空行 ⇒ 位置差，⛔ 代码）
73a80,85      （同上）
103c115
< 			why:     "positive control for offset sensitivity: malformed face …",
> 			why:     "positive control for offset sensitivity (the weakest face in this set, see the scan note above): malformed face …",
123a136       （同上）
rc-diff-nocomment=1
```

⇒ **十处 `le.PutUint*` 写入、四枚 `face300{…}` 字面量、三处断言全部逐字节同**；
唯一的语义差＝S4 那句 `why` 里我按 `300-v1` §6.2 加的"(the weakest face in this set…)"，是**散文**。
不剥注释那一把（`logs/diff-vs-a1r-rig.txt`）另外多出的 12 行全是 `le.PutUint*` 那一段的**注释列宽**差
（`300-a1r` 那枚台件本来⛔ 过 `gofmt`；我这枚过了 ⇒ 列对齐动了，代码没动）。
头顶那段"⛔ 入库的仓外台件"注释被我改写成"随修入库"，那本来就该改。

## 1.2 vet 尺（`unsafe.Pointer` 那枚用法）

```
$ go vet ./internal/audio/
rc-vet=0        （stdout⛔ 空，原始件＝logs/vet.txt，0 字节）
```

⇒ 过 `unsafeptr` 检查，与仓里 `wasapi_windows.go` 头顶那句 `// unsafe discipline (vet-clean, same as internal/ball)` 同形。
⚠ `-gcflags=-m` 那枚"escapes to heap"尺我⛔ 重跑（`300-a1r` 在同形夹具上跑过、命中 16 枚），欠账具名在 `90-unrun-rulers.md` §1。

## 2. 一进门就绿（工作树两发成对，命令逐字）

```
$ PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" go test ./internal/audio/ -count=1 -run 'TestParseWaveFormatSubFormatOffset300' -v
entry-1 rc=0
entry-2 rc=0
```

尺（两发各一把，⛔ 互减）：

```
entry-1: pass-all=5 fail-all=0 pass-top=1 fail-top=0 load0xc=0
entry-2: pass-all=5 fail-all=0 pass-top=1 fail-top=0 load0xc=0
```

- `pass-all` 尺＝`grep -c -- '--- PASS'`（**含子测试**）＝5 ＝ 1 枚顶层 ＋ 4 枚子测试；
  `pass-top` 尺＝`grep -cE '^--- PASS'`（**只有顶层**）＝1。两把差 4 是本套件的结构（4 形），⛔ 拿一把减另一把。
- `load0xc` 尺＝`grep -c '0xc0000135'` ⇒ **0** ⇒ 用例**真跑了**（⛔ "载入期死"被读成红／绿）。

`logs/entry-1.txt` 整块逐字。两发的**成对尺**＝`diff logs/entry-1.txt logs/entry-2.txt` ⇒ 只差第 12 行那一枚耗时
（`0.030s` vs `0.025s`），用例名册逐行同序同色（`rc-diff=1` 是 `diff` 见到差异的正常退码）：

```
12c12
< ok  	github.com/CarlosShao/wisp/internal/audio	0.030s
---
> ok  	github.com/CarlosShao/wisp/internal/audio	0.025s
```

```
=== RUN   TestParseWaveFormatSubFormatOffset300
=== RUN   TestParseWaveFormatSubFormatOffset300/S1_float_at_24
=== RUN   TestParseWaveFormatSubFormatOffset300/S2_pcm_at_24
=== RUN   TestParseWaveFormatSubFormatOffset300/S3_bogus_subtype_7_at_24
=== RUN   TestParseWaveFormatSubFormatOffset300/S4_data1_shifted_to_26
--- PASS: TestParseWaveFormatSubFormatOffset300 (0.00s)
    --- PASS: TestParseWaveFormatSubFormatOffset300/S1_float_at_24 (0.00s)
    --- PASS: TestParseWaveFormatSubFormatOffset300/S2_pcm_at_24 (0.00s)
    --- PASS: TestParseWaveFormatSubFormatOffset300/S3_bogus_subtype_7_at_24 (0.00s)
    --- PASS: TestParseWaveFormatSubFormatOffset300/S4_data1_shifted_to_26 (0.00s)
PASS
ok  	github.com/CarlosShao/wisp/internal/audio	0.030s
```

## 3. 四形各自断的是什么值（★回答"ⓑ 那一形在你的 24 上断的是什么值"）

| 形 | 字节面 | 断言的**值** | 我这把实测 | 颜色 |
|---|---|---|---|---|
| S1 `float_at_24` | `Data1=3` ＠ **24**、bits 32 | `tag == 3` 且 `floating == true`（`convertPacket` 走 float32 那条） | `tag=3`／`floating=true` | PASS |
| S2 `pcm_at_24` | `Data1=1` ＠ **24**、bits 16 | `tag == 1` 且 `floating == false` | `tag=1` | PASS |
| S3 `bogus_subtype_7_at_24` | `Data1=7` ＠ **24** | `tag == 7`（正控·**值**：一枚明显错的 Data1 ⛔ 被吞成 PCM/FLOAT） | `tag=7` | PASS |
| **S4 `data1_shifted_to_26`** | `Data1=3` ＠ **26**（坏面） | **`tag == 0`** 且 `floating == false` | `tag=0` | PASS |

★**S4 在 24 上断的值＝`0`**（不是"读到 3"、⛔ "被调用过"）。它的分量全在"断的是 0"这一句上：
按权威布局，那枚坏面的第 24 字节起是 `dwChannelMask` 的高两字节（`00 00`）＋零字节 ⇒ 读 24 得到 **0**、读 26 得到 **3**。
⇒ 修完之后 S4 是**绿的、且断的是 0**；锚上（读 26）同一枚 S4 是**红的、拿到 3**（凭据＝§5）。
⇒ 这一枚把"偏移往后挪两字节"这个方向钉住了：它⛔ 是一把只会朝一个方向响的尺。
⚠ 同时按 `300-v1` 的扫描它**是四枚里最弱的一枚**（单形在 22／24／28 三形都放行），这句我抄进了用例注释，防下一位只留它。

## 4. 这枚用例的凭据面（★只有本机，⛔ "CI 会兜"）

票面 `AC#5` 那三／四跳编排者已在 2026-10-10 12:5x 自己逐行量完并写进本票"编排者现量"一节，结论：
`internal/audio` 虽被 `core` 档在名册上认领，但那一档 `runs-on: ubuntu-latest` ⇒ windows-tagged 文件编不进测试二进制；
`windows)` 档清单（`portable-tests.sh:248-253`）⛔ `./internal/audio/`；`cli)` 档只有 `./cmd/wisp/`；
`--scope=census`（GUARD D）的判定单位是**导入路径** ⇒ `internal/audio` 永远进不了 `unclaimed`。
⇒ **本枚用例在 CI 无一档真跑；凭据只有本机（这台 Windows 开发机 + 成对导出树）。交付件就这么写。**
⛔ 我为了迁就门禁把落点搬到别处——派单与票面 `AC#5` 双重写死（且 `parseWaveFormat` 包内未导出，"搬去 `cmd/wisp`"这条路物理上不存在）。
⛔ 我自己重跑那四跳（编排者写死"你别再花时间找"）；欠账具名在 `90-unrun-rulers.md`。
★**归口已成形**：编排者在 `6d81638b`（我这把尺现读 HEAD）已把这半件事立成**独立仪器票 301**
（`.scratch/wisp/issues/301-windows-tagged-tests-in-a-core-claimed-package-never-run-in-ci-and-guard-d-is-package-level.md`，
13:1x 立，票面开头就写着"票 300 的落地腿正要往 `internal/audio` 里交一枚 `//go:build windows` 的判据"）⇒
我这枚用例的可见性问题**记在 301**，⛔ 塞进本批（⛔ 动 `scripts/portable-tests.sh`／`.github/workflows/ci.yml`，⛔ 我射程）。

## 5. "改前必红"那半截凭据归谁（⛔ 我重跑、⛔ 我另造）

本程⛔ 需要重证"该不该修"：三把尺已独立量过——`300-a1r`（仓外成对树，改前 4/4 FAIL、改后 4/4 PASS、`cmp` 还原同值）、
`300-v1`（非实现者重跑＋**五枚偏移扫描 20/22/24/26/28 ⇒ 全绿只在 24**）、编排者自己那一把（同一天 13:02 重跑同一扫描，
件＝`probes/300/orch/2026-10-10-orch-five-offset-scan.txt`，逐字：`===== MUT offset=24 =====` 那节 `rc-gotest=0` 四形全 PASS，
`===== MUT offset=26 =====` 那节 `rc-gotest=1` 四形全 FAIL）。
⇒ 我这把只补一件事：**入库的这枚用例在 24 上一进门就绿**（§2），恒红测试⛔ 入库这条定式没被破坏。
⚠ 编排者那件里末行自己登记过一枚仪器缺陷：那三行 `SUBTEST tally` 是整份日志的**累计值**（grep 打在全文而⛔ 只打当节）＝坏尺；
有效读数＝逐节 `rc-gotest`。我引用时只取逐节那几枚，⛔ 取 tally 那三行。
