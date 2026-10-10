# 300-r2 突变凭据（票 300 `AC#6`：成对两发 3→1／3→0，各带逐字红句＋`rc=1`，还原 `cmp`⇒绿）

## 我用的是哪一条形（⛔ 在共享工作树里做任何 checkout／改产码）

**仓外导出树**＝`git archive HEAD | tar -x -C /tmp/wisp300-r2/tree`（锚＝我自己的 `9a442923`，含本腿新增那枚用例）。
⚠ 派单给的两条可选形里我选的是**"仓外导出树"那一支，⛔ 是 `git clone` 那一支**，理由＝导出树⛔ 带 `.git`，
从根上⛔ 可能碰共享工作树的任何 ref；两处已知 handicap 我现量了、并具名它们对本射程 **inert**：
`ls third_party/sherpa-onnx/*.dll` 在该树里 **0 枚**（目录⛔ 存在；编排者 14:2x 已量 `grep -rln 'sherpa\|onnxruntime' internal/audio` ⇒ 0 命中 ⇒ audio 与那三枚 DLL 无关），
`ls -d .git` ⇒ **0**（audio 那 42/43 枚顶层用例里⛔ 一枚读 git 元数据；`TestCleanCheckoutBuilds_AC11` 在 `cmd/wisp`，本腿⛔ 跑那枚包）。
**台面同一性**（`300-v2` 之后立的那条规矩）＝两发突变＋改前改后的整包对照**全在同一枚台面上取**：
突变四发都在导出树（同一路径、同一棵、只改那行再改回来）；整包两发（门③）都在这台机器的**共享工作树**，
且 `cmp` 证导出树里那枚用例与产码文件与工作树**逐字节相同**（读数＝`logs/mut-rig-identity.txt`）。
⚠ 我没有拿"一发母仓、一发 clone"去比名册。

## 四发读数（逐字件都在 `logs/`）

| 发 | 产码那行 | `rc` | 顶层 | 子测试 | 件 |
|---|---|---|---|---|---|
| M0 pristine（导出树，未突变） | `98:	waveFormatFloat = 3` | **0** | `--- PASS` 1 | `--- PASS` 6、`--- FAIL` 0 | `logs/mut-m0-pristine-exporttree.txt` |
| M1 突变 3→1 | LANDS＝`98:	waveFormatFloat = 1` | **1** | `--- FAIL: TestWaveFormatConstantsMatchMmregAuthority300` | FAIL 2（`ieee_float_low_word`／`the_three_constants_map_one_to_one_onto_the_authority_table`），PASS 4 | `logs/mut-m1-float-3to1.txt` |
| M2 突变 3→0 | LANDS＝`98:	waveFormatFloat = 0` | **1** | `--- FAIL: TestWaveFormatConstantsMatchMmregAuthority300` | FAIL 2（同上两枚），PASS 4 | `logs/mut-m2-float-3to0.txt` |
| M3 还原后复跑 | `98:	waveFormatFloat = 3` | **0** | `--- PASS` 1 | `--- PASS` 6 | `logs/mut-m3-restored-exporttree.txt` |

还原自证＝**同一条命令内** `cp` 备份 → 跑 → `cp` 写回 → `cmp` ⇒ `restore_cmp=IDENTICAL`，且锚尺
`grep -n 'waveFormatFloat = 3'` 回到 `98:	waveFormatFloat = 3`（M1、M2 各一次，两次都印在 stdout 里）。

## 逐字红句（**这些是"红句"：`t.Errorf` 那三枚位置，句里点名那枚常量**）

M1（3→1）三句：

```
    wave_format_float_300_windows_test.go:191: waveFormatFloat = 0x0001, want 0x0003; mmreg.h:2110 WAVE_FORMAT_IEEE_FLOAT 0x0003, copied out of the SDK header, not derived from this package. waveFormatFloat is the constant wasapi_windows.go uses for floating := s.format.tag == waveFormatFloat, so its value decides whether 32-bit float frames are sliced as float32 or as int16; at 1 it would fire for PCM mix formats and never for float ones, at 0 for nothing at all
    wave_format_float_300_windows_test.go:214: waveFormatPCM and waveFormatFloat both = 0x0001, but mmreg.h gives WAVE_FORMAT_PCM 1, WAVE_FORMAT_IEEE_FLOAT 0x0003 and WAVE_FORMAT_EXTENSIBLE 0xFFFE as three distinct values, so one of those two constants has moved off the value copied from mmreg.h:2418/2110/2376
    wave_format_float_300_windows_test.go:229: waveFormatFloat = 0x0001, but mmreg.h reserves 0x0001 for WAVE_FORMAT_PCM (mmreg.h:2418), which this package declares as waveFormatPCM; the constant names and the authority values no longer map one to one
```

M2（3→0）两枚红句（`Errorf`）：

```
    wave_format_float_300_windows_test.go:191: waveFormatFloat = 0x0000, want 0x0003; mmreg.h:2110 WAVE_FORMAT_IEEE_FLOAT 0x0003, copied out of the SDK header, not derived from this package. waveFormatFloat is the constant wasapi_windows.go uses for floating := s.format.tag == waveFormatFloat, so its value decides whether 32-bit float frames are sliced as float32 or as int16; at 1 it would fire for PCM mix formats and never for float ones, at 0 for nothing at all
    wave_format_float_300_windows_test.go:235: waveFormatFloat = 0x0000, which is none of WAVE_FORMAT_PCM 1, WAVE_FORMAT_IEEE_FLOAT 0x0003 or WAVE_FORMAT_EXTENSIBLE 0xFFFE as copied from mmreg.h:2418/2110/2376
```

⚠ **口径具名，⛔ 把两形混成一把尺**：M2 的 grep 还捞出三枚 `wave_format_float_300_windows_test.go:194:` 行——那是
`t.Logf` 的**观察记录**（每枚常量一行"observed … against test-carried …"），⛔ 是断言失败；`--- FAIL` 只挂在
上面那两枚 `Errorf` 所在的子测试上。红句计数按 `Errorf` 那把尺＝M1 三句／M2 两句。

**票面 `AC#6` 原文那句"红句具名指向那枚常量而⛔ 指向 `tag` 的等值比较"——上面五句逐字都点名 `waveFormatFloat`（或点名撞值的两枚常量），没有一句写成 `tag = X, want Y` 那种形状。**
（`tag` 那一句形状在 A6/A7/A8 里存在，它们钉的是 `waveFormatExt` 与正控，两枚突变⛔ 碰到它们——见 `10-shape.md` 那张表最后三行。）

## ★ 一枚**⛔ 算凭据**的形状，我先踩了再拆掉（具名，⛔ 抹）

`300-v3`／`300-v4` 之后任何一位引这格时要知道：**把三枚常量当 map 字面量的键，突变就⛔ 红而**是编译错误**。
我在仓外用一枚最小件证了这条语言行为（⛔ 依赖本仓任何产码）：

```
$ cat /tmp/dupkey300/main.go   # const a = 1; const b = 1; map[uint16]string{a: "x", b: "y", 3: "z"}
$ go build ./...
.\main.go:9:33: duplicate key 1 in map literal
rc=1
```

⇒ `rc=1` 而**零句 `--- FAIL`、零句点名常量的红**。那⛔ 是牙，那是一枚⛔ 响的仪器。
⇒ 交付里 A4/A5 用的是**运行期切片＋两两比较**，突变后照样编译、照样红（M1/M2 的 `rc=1` 各带 2–3 句 `Errorf`）。

## 正控两发（裁语与派单各自要的那两枚）

1. **pristine ⇒ 绿**：门②（共享工作树，`logs/gate2-targeted-pristine.txt`，`rc=0`）与 M0（导出树，`rc=0`）**两枚台面各自都绿**，
   两形的件字节数都⛔ 是 0。
2. **观察侧到不了生产那一行**（具名，⛔ 造假的行为面）：`waveFormatFloat` 唯一消费者＝
   `internal/audio/wasapi_windows.go` 的 `floating := s.format.tag == waveFormatFloat`，在 `(*wasapiStream).Drain()` 体内，
   要一条活流＝真设备；票 300 票面第 31 行把真设备归编排者、⛔ 归腿。`convertPacket` 收 `floating` 为 **bool 参数**，
   自己调它证的是参数⛔ 是那枚常量 ⇒ 本腿**没有**交那一形。残余＝编排者裁语里具名的 **R2**（要买它得走 C5＝改产码语义面，本票⛔ 授权）。

## 本腿⛔ 做／⛔ 该被读成做到的事（写死，免得下一程读大）

- ⛔ 证明"生产真的按 3 判浮点"（R2 未闭，本机⛔ 到）。
- ⛔ 给 `injector` 那两支（`fmtTag == 3`／`0xFFFE`）买到覆盖面＝**没有**。那是裁语里的**残余 R1**，要买得用覆盖尺（`go test -coverprofile`），
  本腿⛔ 跑覆盖尺；本件⛔ 调 `parseWav`（⛔ 把 C4 当本格裁判＝编排者已裁⛔）。
- ⛔ 证明"值⇔外部权威"在 CI 上可核：权威只以**抄进测试的字面量＋注释**存在，读盘那一发只在 `logs/authority-mmreg.txt`（编排者裁语 C3 的处置＝证据进探针件、⛔ 进 CI 分母；他收件时要自己再跑一遍那把尺）。
- ⛔ 删／改 `waveFormatPCM`（只钉值，"零使用该不该删"归另一格）。
