# 300-v2 — `10` `AC#2` 成不成立：五枚偏移扫描（换形也红）＋夹具等价对拉

★尺的口径先写死：**这不是"改前必红"那一发，而是"把生产偏移逐枚换成别的数，用例⛔ 恒真"**。
射程＝`/tmp/wisp300-v2/treeH`（`git archive 41329475` 导出，⛔ 共享工作树）。
`sed` 射程**只打那一行**，逐字模式串＝`sub := \*(\*windows\.GUID)(unsafe\.Add(p, 24))` ⇒ 替换成 `…, $N))`。

## 0. 被扫的两枚件与"只有第 4 枚读数在动"的自证

```
$ wc -c internal/audio/parse_wave_format_300_windows_test.go       # 工作树（＝入库 blob）
7267
$ git show 41329475:internal/audio/parse_wave_format_300_windows_test.go | wc -c
7267                                       # blob↔工作树同值（⛔ 我这一程⛔ 碰它一个字节）
```

每一枚突变之后当场量"别的偏移⛔ 被我一起改了吗"（尺＝`grep -o 'unsafe\.Add(p, [0-9]*)'` 对备份做 `diff`）：

```
OFF=20 offsetsdiff[4c4 < unsafe.Add(p, 24) --- > unsafe.Add(p, 20) ]
OFF=22 offsetsdiff[…同上，只 24→22…]
OFF=26 offsetsdiff[…同上，只 24→26…]
OFF=28 offsetsdiff[…同上，只 24→28…]
OFF=24 offsetsdiff[]                      # 未突变的那一发，差集为空
```

⇒ 每枚突变**恰好一行、恰好一个数**（`@2`／`@4`／`@14` 三枚读数⛔ 动）⇒ 派单点名的"第一发把 `sed` 射程打宽到整文件＝台件无效"那一形，我这把⛔ 踩（差集非空即作废，四枚差集都只含第 4 枚）。

## 1. 逐节 PASS/FAIL 表（★grep 射程钉在**当节自己的日志文件**内，⛔ 整份累计）

命令逐字（每枚偏移一发，⛔ 互减）：

```
cd /tmp/wisp300-v2/treeH && go test ./internal/audio/ -count=1 -run 'TestParseWaveFormatSubFormatOffset300' -v > scan-$N.txt
```

| 生产偏移 | S1 `float_at_24` | S2 `pcm_at_24` | S3 `bogus_7_at_24` | S4 `data1_shifted_to_26` | `pass-all` | `fail-all` | `pass-top` | **rc** |
|---|---|---|---|---|---|---|---|---|
| **24（入库那一枚）** | PASS | PASS | PASS | PASS | **5** | 0 | 1 | **0** |
| 20 | **PASS** | FAIL | FAIL | **FAIL** | 1 | 4 | 0 | 1 |
| 22 | FAIL | FAIL | FAIL | **PASS** | 1 | 4 | 0 | 1 |
| 26（＝改前那一枚） | FAIL | FAIL | FAIL | FAIL | **0** | 5 | 0 | 1 |
| 28 | FAIL | FAIL | FAIL | **PASS** | 1 | 4 | 0 | 1 |

尺名一律带出：`pass-all`＝`grep -c -- '--- PASS'`（含子测试），`pass-top`＝`grep -cE '^--- PASS'`（只顶层），逐节各数**当节那份文件**。
`load0xc` 我这把没在这一发量（窄档⛔ 依赖 harness，见 `00` 件 §3），但顶层 `ok … 0.030s` 与 `rc` 同向 ⇒ 用例**真跑了**（⛔ 载入期死被读成红）。

- ★**只有 24 让四形同时绿**；26 四形**全红**；20／22／28 **各只绿一枚** ⇒ 与编排者 13:0x 那把的**按节重算表**（`20→1/3、22→1/3、24→4/0、26→0/4、28→1/3`）**逐枚同值**，也与 `300-v1` 那把同向。
  ⇒ 我这把**对得上**，⛔ 报分歧。
- ★**"一把尺两种形状都放行"这一形被排除**：同一份用例在 24 全绿、在其余四枚偏移里最多绿一枚 ⇒ `AC#2` 的那枚用例**⛔ 恒真**，它有牙。
- ⚠ **判别力扛在 S2＋S3**（这两枚只绿 24），S1 单形在 20 也绿（读到 `dwChannelMask`＝3），S4 单形在 22／24／28 都绿。⇒ 少任何一枚⛔ 恒真，但**四枚里最弱的是 S4（放行两枚错偏移）**、其次是 S1。

### 1.1 红句逐字（本仓定式＝红名册⛔ 只数颜色，坏法种类要分得开）

`scan-20.txt`（S4 在 20 这一枚**红**，拿到的是 3）：

```
    parse_wave_format_300_windows_test.go:135: tag = 3, want 1 (KSDATAFORMAT_SUBTYPE_PCM (mmreg.h:2474, Data1=0x00000001) at byte 24)
    parse_wave_format_300_windows_test.go:135: tag = 3, want 7 (positive control: an obviously wrong Data1 must NOT be swallowed into PCM/FLOAT; …)
    parse_wave_format_300_windows_test.go:135: tag = 3, want 0 (positive control for offset sensitivity (the weakest face in this set, see the scan note above): malformed face with Data1 two bytes late; per the authoritative layout the value at byte 24 is 0, so an assertion that passes here would read nothing meaningful)
    parse_wave_format_300_windows_test.go:152: floating = true, want false (convertPacket would take the int16 path)
```

`scan-26.txt`（＝改前那一枚，四形全红的两类坏法）：

```
S1: tag = 0, want 3 (KSDATAFORMAT_SUBTYPE_IEEE_FLOAT (mmreg.h:2483, Data1=0x00000003) at byte 24) ＋ floating = false, want true (convertPacket would take the float32 path)
S2: tag = 0, want 1 (KSDATAFORMAT_SUBTYPE_PCM (mmreg.h:2474, Data1=0x00000001) at byte 24)
S3: tag = 0, want 7 (positive control: …)
S4: tag = 3, want 0 (positive control for offset sensitivity (the weakest face in this set, …)) ＋ floating = true, want false
--- FAIL: TestParseWaveFormatSubFormatOffset300 (0.00s)
    --- FAIL: …/S1_float_at_24  --- FAIL: …/S2_pcm_at_24  --- FAIL: …/S3_bogus_subtype_7_at_24  --- FAIL: …/S4_data1_shifted_to_26
```

⇒ 两类坏法**名字分得开**：读早了（20）＝**恒读到 3**（`dwChannelMask`），读晚了（26）＝**好面读到 0、坏面读到 3**。⇒ 这枚尺对**方向**也敏感。

## 2. 还原自证（⛔ 只在仓外突变，跑完必证还原）

```
$ cp /tmp/wisp300-v2/wasapi.bak internal/audio/wasapi_windows.go && cmp …
restore_cmp=IDENTICAL
$ grep -c 'unsafe.Add(p, 24)' → 1     # 回到入库那一枚
  grep -c 'unsafe.Add(p, 26)' → 0
  grep -c 'unsafe.Add(p, 20)' → 0
  grep -c 'unsafe.Add(p, 22)' → 0
  grep -c 'unsafe.Add(p, 28)' → 0
```

★**第二把还原尺（更硬的形）**：扫描全部跑完之后，我把那棵树的两个文件与**起手锚的 blob** 逐字节 `cmp`：

```
cmp-treeH-code-vs-HEADblob=IDENTICAL
cmp-treeH-test-vs-HEADblob=IDENTICAL
```

⇒ 那棵树现在就是 `41329475` 的形态，`full-postH-1` 那一发跑的就是**入库字节**（⛔ 带着我突变过的字节）。

## 3. 入库用例 ↔ `300-a1r` 仓外台件＝代码逐字节等价（★"落地时偷偷改夹具"这一形）

台原件还在：`/tmp/wisp300-a1r/tree/internal/audio/parse_wave_format_300_windows_test.go`（**6172 字节**，⛔ 我造成）。
尺＝两枚文件各剥行尾注释再 `diff`（⛔ 引腿的 `logs/diff-vs-a1r-nocomment.txt`，我自己重跑）：

```
$ sed 's/[[:space:]]*\/\/.*$//' a1r原件 > a ; 同尺 landed > b ; diff a b
26a27,33   （+ 6 枚空行：块注释被剥成空 ⇒ 位置差，⛔ 代码）
73a81,86   （同上）
103c116    why: "positive control for offset sensitivity:" ⇒ "positive control for offset sensitivity (the weakest face in this set, see the scan note above):"
123a137    （空行）
rc-diff=1
$ grep -c "le.PutUint"  a1r原件=10  landed=10
$ grep -c "face300{"    a1r原件=4   landed=4
```

⇒ **十处写入、四枚字面量、三处断言逐字节同**，唯一语义差＝S4 那句 `why` 的**散文**（腿在 `20` 件 §1.1 报的形状与我这把**逐条对上**，含 `rc-diff=1`）。
⇒ 那一套夹具正是被 `300-a1r`／`300-v1`／编排者三把尺扫过的同一套字节；入库⛔ 换形、⛔ 另造布局。

## 4. 判语（`AC#2`）

**成立（就 `AC#2` 这一格所要的四件事而言）**，各配我自己的尺：

1. **偏移落在那一枚权威数上**：盘上 `211:		sub := *(*windows.GUID)(unsafe.Add(p, 24))`〔盘上现量〕。
2. **用例⛔ 恒真**：§1 那五枚偏移里只有 24 四形全绿、26 全红；§3 证明入库夹具与被扫过的那套逐字节同〔盘上现量〕。
3. **一进门就绿**：工作树与导出树两形各自 `rc=0`、`pass-all=5`／`pass-top=1`（`0` 件＋§1 第一行）〔盘上现量〕。
4. **同批把注释改成实话**：改写落地了、且**⛔ 夹带**任何一枚别人的文件（名册尺见 `50` 件）；那句"坏掉的只有两枚数"我与腿同判，与 `300-a1r` 那句"三个数全是坏的"**⛔ 同判**（裁决与尺在 `20` 件）〔盘上现量〕。

⚠ **我这一格⛔ 判到的东西（具名）**：`AC#2` 这一格⛔ 等于本票可结。
`AC#3`（真机 `GetMixFormat` 前 40 字节 hexdump）与 `AC#5`（凭据只有本机）两格⛔ 在我这关门；
"这枚缺陷今天在这台机器上响不响"这一句**仍只有编排者那一发能答**〔仅本机可量、归编排者〕。
