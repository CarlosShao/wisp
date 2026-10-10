# 300-v1 — `20` ② `AC#1` 那套台件**是不是循环论证**（我自己重跑＋我自己扫偏移＋还原自证）

锚＝`efd85155`；被审树＝`/tmp/wisp300-a1r/tree`（`git archive b2933dc9` 导出，腿留着没删）。
`b2933dc9..efd85155` 之间 `internal/audio/**` **零改动**（尺见 `00-anchor.md`）⇒ 我这棵树上跑的就是被审的那枚产码。

## 判语：**成立**——台件⛔ 循环。夹具的 10 处写入**逐处落在我自己从 `mmreg.h` 算出来的那张偏移表上**（`10` 件 §4），
## 而我**在看任何读数之前**就能只用"权威＋夹具"预言锚上四形的读数：S1/S2/S3 读 **0**、S4 读 **3**。实测**逐字命中**。

## 1. 夹具逐处 → 权威格（我自己逐字节读 `build()`，尺＝`Read` 全文 144 行；⛔ 引用腿的映射）

树上那枚文件＝`/tmp/wisp300-a1r/tree/internal/audio/parse_wave_format_300_windows_test.go`，**6172 字节**
（尺＝`ls -l`；与被审件报的字节数**同值**），内容与 `20-ac1-rig.md` 里的代码块**逐字一致**（我逐行对过，含注释）。

| `build()` 里的写入（逐字） | 落在哪一格 | 我这把尺的权威出处 |
|---|---|---|
| `le.PutUint16(b[0:], 0xFFFE)` | `Format.wFormatTag` @0–1 | mmreg.h:2376 ＋ 我的偏移表 |
| `le.PutUint16(b[2:], f.channels)` | `Format.nChannels` @2–3 | WAVEFORMATEX 抽段 |
| `le.PutUint32(b[4:], f.rate)` | `Format.nSamplesPerSec` @4–7 | 同上 |
| `le.PutUint32(b[8:], f.rate*ch*(bits/8))` | `Format.nAvgBytesPerSec` @8–11 | 同上 |
| `le.PutUint16(b[12:], f.channels*(f.bits/8))` | `Format.nBlockAlign` @12–13 | 同上 |
| `le.PutUint16(b[14:], f.bits)` | `Format.wBitsPerSample` @14–15 | 同上 |
| `le.PutUint16(b[16:], 22)` | `Format.cbSize` @16–17 ＝ **22** | mmreg.h:2540/2550 两处逐字 `/* Format.cbSize = 22 */` |
| `le.PutUint16(b[18:], f.bits)` | `Samples.wValidBitsPerSample` @18–19（**2 字节 union**） | WAVEFORMATEXTENSIBLE 抽段：union 三成员皆 WORD |
| `le.PutUint32(b[20:], uint32((1<<f.channels)-1))` | `dwChannelMask` @20–23 | 同上（union **外面**的 DWORD） |
| `le.PutUint32(b[f.data1At:], f.data1)` | `SubFormat.Data1` @24–27（S4 故意写 26） | 偏移表 ⇒ **24** ＋ guiddef.h `Data1`＝`unsigned long` |

★**这十处里没有任何一处的位置来自仓里任一枚实现**——`wasapi_windows.go` 的 `@18/@22/@26` 那三个数**一个都没被用**
（用了 `18` 那处的宽度是 2 而⛔ 注释暗示的 4，用了 `22` 的是 `cbSize`@16 算出来的 18＋4 而⛔ 注释）。

## 2. 反循环的那一发硬证：**读数在我看到它之前就被权威算出来了**

只用 §1 那张表 ＋ 夹具取值，我自己手算四形在**两种世界**里的读数（尺＝算术，⛔ 跑之前算的）：

| 形 | 面内字节（我自己按 `build()` 逐格填） | 权威世界（读 24）**该读到** | 锚上世界（读 26）**该读到** |
|---|---|---|---|
| S1 `bits=32, ch=2, Data1=3@24` | `…b20-23=03 00 00 00, b24-27=03 00 00 00` | 3 ⇒ **绿** | `b26-29=00 00 00 00` ⇒ 0 ⇒ **红** |
| S2 `bits=16, ch=2, Data1=1@24` | 同上，`b24-27=01 00 00 00` | 1 ⇒ **绿** | 0 ⇒ **红** |
| S3 `bits=32, Data1=7@24` | `b24-27=07 00 00 00` | 7 ⇒ **绿** | 0 ⇒ **红** |
| S4 `bits=32, Data1=3@26`（坏面） | `b24-25=00 00, b26-29=03 00 00 00` | `b24-27=00 00 03 00` ⇒ 低字 **0** ⇒ **绿** | `Data1=3` ⇒ 3 ⇒ **红** |

⇒ 我的预言＝"锚上四枚**全红**，且红句里的数字是 `tag = 0 / 0 / 0 / 3`"。
**实测逐字命中**（下面 §3）。
★这一发为什么把"循环"这条路堵死：**如果夹具是按产码的假设造的，那锚上的码就会读到它自己那枚布局的值（1/3/7 ⇒ 全绿）**；
它读到的是 0，说明夹具站在权威那一侧、与产码相对。

## 3. 我自己重跑的那一发（改前；命令逐字＝派单写死那一发，先核树上那行仍是 `26` 再跑）

```
cd /tmp/wisp300-a1r/tree && PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH" \
  go test ./internal/audio/ -count=1 -run TestParseWaveFormatSubFormatOffset300 -v
```

前置核对（逐字）：`grep -c "unsafe.Add(p, 26)" …/wasapi_windows.go` ⇒ **1**；
`grep -n "sub := " ` ⇒ `196:		sub := *(*windows.GUID)(unsafe.Add(p, 26))`。

stdout **整块逐字**（全文＝`logs/ac1-v1-pristine.txt`，2323 字节，尺＝`ls -l`）：

```
=== RUN   TestParseWaveFormatSubFormatOffset300
=== RUN   TestParseWaveFormatSubFormatOffset300/S1_float_at_24
    parse_wave_format_300_windows_test.go:122: tag = 0, want 3 (KSDATAFORMAT_SUBTYPE_IEEE_FLOAT (mmreg.h:2483, Data1=0x00000003) at byte 24)
    parse_wave_format_300_windows_test.go:138: floating = false, want true (convertPacket would take the float32 path)
=== RUN   TestParseWaveFormatSubFormatOffset300/S2_pcm_at_24
    parse_wave_format_300_windows_test.go:122: tag = 0, want 1 (KSDATAFORMAT_SUBTYPE_PCM (mmreg.h:2474, Data1=0x00000001) at byte 24)
=== RUN   TestParseWaveFormatSubFormatOffset300/S3_bogus_subtype_7_at_24
    parse_wave_format_300_windows_test.go:122: tag = 0, want 7 (positive control: an obviously wrong Data1 must NOT be swallowed into PCM/FLOAT; the parser has to report the value it actually read, so this face may not go green on a constant)
=== RUN   TestParseWaveFormatSubFormatOffset300/S4_data1_shifted_to_26
    parse_wave_format_300_windows_test.go:122: tag = 3, want 0 (positive control for offset sensitivity: malformed face with Data1 two bytes late; per the authoritative layout the value at byte 24 is 0, so an assertion that passes here would read nothing meaningful)
    parse_wave_format_300_windows_test.go:138: floating = true, want false (convertPacket would take the int16 path)
--- FAIL: TestParseWaveFormatSubFormatOffset300 (0.00s)
    --- FAIL: TestParseWaveFormatSubFormatOffset300/S1_float_at_24 (0.00s)
    --- FAIL: TestParseWaveFormatSubFormatOffset300/S2_pcm_at_24 (0.00s)
    --- FAIL: TestParseWaveFormatSubFormatOffset300/S3_bogus_subtype_7_at_24 (0.00s)
    --- FAIL: TestParseWaveFormatSubFormatOffset300/S4_data1_shifted_to_26 (0.00s)
FAIL
FAIL	github.com/CarlosShao/wisp/internal/audio	0.027s
FAIL
```

- **`rc-gotest=1`**（取法＝`go test` 之后立刻 `${PIPESTATUS[0]}`，那一发**中间没有管道** ⇒ 就是 `go test` 自己的退码；⛔ 从 `tee`/`tail` 取）。
- 计数两把尺（⛔ 互减）：顶层 `grep -cE '^--- (PASS|FAIL)'`＝**1**（那一枚顶层 FAIL）；含子测试行＝**5** 行 FAIL（1 顶层＋4 子）。
- ★**没有 `exit status 0xc0000135`** ⇒ 用例真跑了（⛔ "没跑成"被读成红）。这一条我另有一把独立尺，见 `30` 件 ④。

## 4. 我这把尺比腿多扫了四枚偏移：**这枚套件的绿，只在 24 这一形出现**

突变台件（⛔ 只改一枚数字，逐次从备份还原后再 sed，⛔ 连环 sed 打空）：

| 树上偏移 | S1 | S2 | S3 | S4 | 套件颜色 | rc |
|---|---|---|---|---|---|---|
| **20**（`dwChannelMask` 起点） | **PASS**（读到 3） | FAIL `tag = 3, want 1` | FAIL `tag = 3, want 7` | FAIL `tag = 3, want 0` | **红** | 1 |
| **22** | FAIL `tag = 0` | FAIL | FAIL | **PASS**（0） | **红** | 1 |
| **24**（权威） | PASS | PASS | PASS | PASS | **全绿** | **0** |
| **26**（锚上原样） | FAIL `tag = 0` | FAIL | FAIL | FAIL `tag = 3, want 0` | **全红** | 1 |
| **28** | FAIL `tag = 0` | FAIL | FAIL | **PASS**（0） | **红** | 1 |

原始逐字（每一形的 `landed-as` 行＋两套计数）＝`logs/ac1-v1-mutations.txt`。
⚠ **我自己踩到的一发，写下来给下一位**：第一版台件我用 `for OFF in 20 22 24` 直接 `sed s/26/$OFF/` **不还原**
⇒ 第二、三发是**空操作**（文件里已⛔ 再有 `26` 可替），三发其实全是偏移 20。
靠 `landed-as` 那把尺（每发都 `grep -n "sub := "` 现读）当场抓到，重写成"每次从备份 `cp` 还原再 sed"后才是上面这张表。
⇒ 本仓那条"**突变台件要先证落地**"救了我自己一次；⛔ 落地就写读数＝假读数。

### 4.1 恒真的第三形（一把尺两形都放行）——**这枚套件⛔ 中招**

五枚偏移里**只有 24 让套件全绿**；26 全红、20/22/28 各只放行 1 枚 ⇒ 判据⛔ 是"两种形状都放行"的尺。✅

★但我给一枚腿⛔ 报过的**子集级**结论：**单形才危险**。
`S1` **单独**是一枚两形都放行的尺（**20 和 24 都绿**）；`S4` **单独**是一枚三形放行的尺（**22／24／28 都绿**）。
⇒ 判别力实际由 **S2（`Data1=1` 与 `dwChannelMask=3` 在 20 那一形撞车）＋S3（7 谁都撞不到）**扛着。
⇒ **AC#2 入库时⛔ 把 S1 或 S2 任意一枚当"决定性那一发"单独留着**；这套件的 strength 是**四枚合起来**才成立的。

### 4.2 恒真的第四形（断言所需的那个世界没被造出来）——**登记一形，⛔ 不挡本格**

- 四形**全部**写 `wFormatTag=0xFFFE` ⇒ 走的**只有** `if f.tag == waveFormatExt { … }` 那一支（尺＝我现读产码 `99: waveFormatExt = 0xFFFE` ＋ `195: if f.tag == waveFormatExt {`）。
  **非 extensible 那一支（plain PCM，`cbSize=0`、18 字节面）在这套件里根本没被造出来** ⇒
  将来谁把 `f.tag` 的第一次读取（`*(*uint16)(p)`）改坏，这套件⛔ 红得起来吗？会（四形都从 0xFFFE 起步 ⇒ 那枚改动多半让分支不进 ⇒ tag=65534，四形全红），
  但"**16 位 plain PCM 混音格式被读对**"这件事**没有一枚用例钉着**。⇒ 具名给 `AC#2`：**要不要顺手补一枚非 extensible 面**是它的取舍，⛔ 本格判 `AC#1` 的判据需要它（票面 `AC#1` 断的就是 extensible 面）。
- `parseWaveFormat` 那枚 `degenerate mix format`（`rate==0 || channels==0`）分支⛔ 被造出来——同样⛔ 属 `AC#1` 的射程。

## 5. 还原自证（⛔ 把突变留在我⛔ 拥有的那棵树里；被审那棵我一枚没留）

```
restored-identical=yes                  # cmp 我用的备份 vs 树上文件
grep -c 'unsafe.Add(p, 26)' = 1         # 锚上那一行回来了
grep -c 'unsafe.Add(p, 24)' = 0
grep -cE 'unsafe\.Add\(p, (20|22|28)\)' = 0
cmp vs /tmp/wisp300-a1r/wasapi_windows.go.orig ⇒ 逐字节相同（腿的备份也在场，我没动它）
```

逐字＝`logs/ac1-v1-mutations.txt` 末段。`/tmp/wisp300-a1r/tree` **留着⛔ 删**（派单写死）。
⚠ 我自己另建的两件东西**只在我自己的目录**：`/tmp/wisp300-v1-orig/wasapi_windows.go`（备份）、`/tmp/wisp300-v1/treeA`（④ 用的导出树）。
⇒ 我在 `git archive` 的导出树里⛔ 留任何未还原改动，仓里 `internal/audio` 我也⛔ 碰过（越界尺见 `30` 件 §4）。

## 6. 我与被审腿分歧的地方（具名，⛔ 圆场）

1. **它那句"★这枚尺两个方向都敏感"给的理由⛔ 是最硬的那条。** 它说敏感是因为 S4 在锚上读到 3。
   我这把偏移扫描尺显示：**判别力在 S2/S3 手上**（S1 单形在偏移 20 也会绿 ⇒ S1 单独⛔ 是尺）。
   ⇒ 结论没变（这套件**确实**有牙），**理由要换**：入库那枚用例的价值在**四形合起来只认 24**，⛔ 在"有一枚坏面"。
2. 它写 `S4` 的 `why` 是"positive control for offset sensitivity"。按我上面这把尺，**S4 是全套件里最弱的一枚**
   （放行 22/24/28 三形）。把它叫"正控"容易让下一位误以为可以只留它 ⇒ 建议 `AC#2` 入库时**注释里写清它是四枚里最弱的一枚**（本格只裁形状，⛔ 写用例）。
3. 它报"改前 5 行 FAIL／改后 5 行 PASS"——我这把尺**复现**到了同形（5 行、含 1 枚顶层）。**这条⛔ 分歧，是追认。**
