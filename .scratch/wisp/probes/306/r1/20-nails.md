# 306-r1 · 20 三对突变凭据（`AC#2`／`AC#2b`，加 `AC#1` 那一枚的牙）

台面＝**仓外导出树** `C:/Users/swq/tmp/306r1-mut-20261011-0730`（`git archive c0d5ec24 | tar -x`；⛔ 在共享工作树里做任何改产码／checkout，
真仓的 `internal/audio/wavinjector.go` 全程一字未动，见 §5）。导出树里的产码基线逐字＝`pristine-wavinjector.go`
（`cmp` 对真仓工作树 `rc_tree_vs_worktree=0`，件＝`logs/00-anchor.txt` 同型尺＋`logs/mutate-final-safety.txt`）。
驱动脚本＝`logs/mutate.sh`（可逐字重跑：`bash mutate.sh <树>`）。

三枚各一对，⛔ 用"一次夹具三枚都装上了"当三格都闭的证据：**每一枚自己那一发红**，红句逐字引在下面。
★第四发（丁）＝按派单**字面**写法的那一枚，见 §4。

每对的四发＝① 改前台面（把那枚新用例 `mv` 走、⛔ 删）整包照绿＝"没有尺"那一半；② 夹具在场、目标用例当场红；
③ `cp pristine` 还原＋`cmp` 逐字节；④ 还原后定向复绿。

## 一（`AC#2`）＝`wavinjector.go:192` 的 `0xFFFE` → `0xFFFD`

改前逐字行（`sed -n '192p' pristine-wavinjector.go`）：

```
			if fmtTag == 0xFFFE { // WAVE_FORMAT_EXTENSIBLE: real tag is SubFormat[0:2]
```

改后退码与逐字红句（`logs/mut-jia-192-fffe-to-fffd-03-red.txt`，`rc_named_cases_with_mutation=1`；
名册＝顶层 `--- FAIL` **2** 枚／`--- PASS` **0**／`--- SKIP` **0**，件＝`logs/mut-jia-192-fffe-to-fffd-SUMMARY.txt`）：

```
--- FAIL: TestParseWavExtensibleFloat32306 (0.00s)
    wavinjector_extensible_float_306_test.go:177: parseWav error = unsupported wav format: tag=65534 bits=32 (want PCM16 or float32), want nil: an extensible fmt tag with a float SubFormat has to resolve to the float32 path (tag read at wavinjector.go:192 and :196, branch chosen at :218); rate=0 chans=0
--- FAIL: TestWavInjectorExtensibleFloat32306 (0.00s)
    wavinjector_extensible_float_306_test.go:213: NewWavInjector error = wav C:\Users\swq\AppData\Local\Temp\TestWavInjectorExtensibleFloat323064246903304\001\extensible-float32-306.wav: unsupported wav format: tag=65534 bits=32 (want PCM16 or float32), want nil: the injector must accept a 40-byte extensible fmt chunk with a float32 SubFormat, not reject it as an unsupported wav format
```

改前对照（夹具 `mv` 走、产码已带突变，整包 `GOFLAGS= go test ./internal/audio/ -count=1 -v`）：
`rc_pkg_no_fixture=0`、顶层 `--- FAIL` **0**／`--- PASS` **42**／`--- SKIP` **1**（`TestLiveWasapiSmoke`，⛔ 本腿造的，见 `40-gates.md` §5），
FAIL 名册里⛔ 任何提及＝`NONE`（件＝`logs/mut-jia-192-fffe-to-fffd-02-nofixture-fullpkg.txt`）。
⇒ 这一发**逐字复现**票面 §现量 那一句"`0xFFFE`→`0xFFFD` 包级尺照绿、`rc=0`、顶层 PASS 42"（腿⛔ 再造该证据，此处＝同台面复跑并对拉）。

`cmp` 逐字节还原＝`logs/mut-jia-192-fffe-to-fffd-04-cmp.txt`，**`rc_cmp_restore=0`**（⛔ 任何差异输出）。
还原后定向复绿＝`logs/mut-jia-192-fffe-to-fffd-05-regreen.txt`，**`rc_regreen=0`**，两枚指名用例逐字 `--- PASS`。

## 二（`AC#1` 的牙）＝`wavinjector.go:218` 的 `3` → `2`

改前逐字行：

```
		case fmtTag == 3 && bits == 32:
```

改后退码与逐字红句（`logs/mut-yi-218-three-to-two-03-red.txt`，`rc_named_cases_with_mutation=1`；顶层 FAIL **2**／PASS **0**／SKIP **0**）：

```
--- FAIL: TestParseWavExtensibleFloat32306 (0.00s)
    wavinjector_extensible_float_306_test.go:177: parseWav error = unsupported wav format: tag=3 bits=32 (want PCM16 or float32), want nil: an extensible fmt tag with a float SubFormat has to resolve to the float32 path (tag read at wavinjector.go:192 and :196, branch chosen at :218); rate=0 chans=0
--- FAIL: TestWavInjectorExtensibleFloat32306 (0.00s)
    wavinjector_extensible_float_306_test.go:213: NewWavInjector error = wav C:\Users\swq\AppData\Local\Temp\TestWavInjectorExtensibleFloat323062182935132\001\extensible-float32-306.wav: unsupported wav format: tag=3 bits=32 (want PCM16 or float32), want nil: ... not reject it as an unsupported wav format
```

★红句的判别力＝`tag=3`：`fmtTag` 被 `:196` 正确重写成了权威值 `3`，而 case 那侧被换成了 `2` ⇒ 两枚值对不上 ⇒ 落 `:226`。
这一句与那一枚（`tag=65534`）、下面那一枚（`tag=0`）**三枚互不相同**，⛔ 一句红被复用成三格的读数。

改前对照：`rc_pkg_no_fixture=0`、FAIL **0**／PASS **42**／SKIP **1**、`NONE`（`logs/mut-yi-218-three-to-two-02-nofixture-fullpkg.txt`）。
`cmp` 还原＝`rc_cmp_restore=0`（`logs/mut-yi-218-three-to-two-04-cmp.txt`）。复绿＝`rc_regreen=0`（`…-05-regreen.txt`，两枚逐字 `--- PASS`）。

## 三（`AC#2b`）＝`wavinjector.go:196` 的 SubFormat 偏移 `body+24` → `body+26`

改前逐字行（票面收腿节 3 的内容锚逐字＝`fmtTag = binary.LittleEndian.Uint16(data[body+24 : body+26])`）：

```
				fmtTag = binary.LittleEndian.Uint16(data[body+24 : body+26])
```

改后退码与逐字红句（`logs/mut-bing-196-off24-to-26-03-red.txt`，`rc_named_cases_with_mutation=1`；顶层 FAIL **2**／PASS **0**／SKIP **0**）：

```
--- FAIL: TestParseWavExtensibleFloat32306 (0.00s)
    wavinjector_extensible_float_306_test.go:177: parseWav error = unsupported wav format: tag=0 bits=32 (want PCM16 or float32), want nil: an extensible fmt tag with a float SubFormat has to resolve to the float32 path (tag read at wavinjector.go:192 and :196, branch chosen at :218); rate=0 chans=0
--- FAIL: TestWavInjectorExtensibleFloat32306 (0.00s)
    wavinjector_extensible_float_306_test.go:213: NewWavInjector error = wav ...: unsupported wav format: tag=0 bits=32 ... not reject it as an unsupported wav format
```

★这一枚改的是**读哪儿**⛔ 是读出来的值：`body+26..28` 落在 GUID 里 `Data1` 的高两个字节＝`00 00` ⇒ `fmtTag` 变成 `0` ⇒ 落 `:226`。
红句里的 `tag=0` 就是"偏移被写歪"这件事在盘上的**唯一可见形状**（票 300 的 S4 负控在 `wasapi_windows.go` 那侧量的是同一件事，两枚⛔ 同一把尺）。

★与派单措辞的一枚**具名偏差**（先报，再说本腿怎么补）：派单写"只改那一处偏移，语句其余字不动"。
本腿按字面只动起点那一枚（＝§四那一发）在盘上⛔ 是一句断言红，而是 `panic: runtime error: index out of range [1] with length 0`
——`data[body+26 : body+26]` 是合法的空切片，`binary.LittleEndian.Uint16` 拿到 0 字节才炸（`encoding/binary/binary.go:70`）。
⇒ 本格的主凭据（本节）取**把同一枚切片的上下界一起移**的写法：`data[body+24 : body+26]` → `data[body+26 : body+28]`，
语句其余字（`fmtTag =`／`binary.LittleEndian.Uint16`／行号）⛔ 动，改的仍然只有"SubFormat 起点偏移"这一枚权威值；
字面那一枚⛔ 丢掉，作为第四发交在下面，两发各自带 `cmp` 还原与复绿。**以盘上为准。**

改前对照：`rc_pkg_no_fixture=0`、FAIL **0**／PASS **42**／SKIP **1**、`NONE`（`logs/mut-bing-196-off24-to-26-02-nofixture-fullpkg.txt`）
⇒ 编排者复验件里那句"块 `196.5,196.65 1 0`"的**因果**这一发坐实了：那一段⛔ 执行者，所以改它⛔ 红。
`cmp` 还原＝`rc_cmp_restore=0`。复绿＝`rc_regreen=0`（两枚逐字 `--- PASS`）。

## 四（丁，派单字面那一发，⛔ 替换上面第三格的主凭据）

改后逐字行＝`				fmtTag = binary.LittleEndian.Uint16(data[body+26 : body+26])`
退码＝`rc_named_cases_with_mutation=1`（`--- FAIL: TestParseWavExtensibleFloat32306`，顶层 FAIL **1** 枚；件＝`logs/mut-ding-196-literal-wording-03-red.txt`）。
★红句形状⛔ 是断言句，逐字＝`panic: runtime error: index out of range [1] with length 0 [recovered, repanicked]`，
栈里逐字两行＝`encoding/binary.littleEndian.Uint16(...)  D:/work/base/go/src/encoding/binary/binary.go:70` 与
`github.com/CarlosShao/wisp/internal/audio.parseWav(...) C:/Users/swq/tmp/306r1-mut-20261011-0730/internal/audio/wavinjector.go:196`。
⇒ **panic 也算"当场红"，但它炸掉整个测试二进制**：那一发里 `TestWavInjectorExtensibleFloat32306` 根本没跑到（名册只有 1 枚 FAIL），
这就是本腿⛔ 拿它当主凭据的原因（主凭据要的是"一把尺红一句、其余尺照跑"）。
改前对照＝`rc_pkg_no_fixture=0`／FAIL **0**／PASS **42**；`cmp` 还原＝`rc_cmp_restore=0`；复绿＝`rc_regreen=0`。

## 5. 三枚汇总尺（一把尺一行，⛔ 相加、⛔ 混用）

| 枚 | 突变 | 改前（夹具 mv 走，整包） | 改后（夹具在场，定向 `-run 306`） | `cmp` 还原 | 复绿 |
|---|---|---|---|---|---|
| `:192` | `0xFFFE`→`0xFFFD` | `rc=0`，顶层 FAIL 0／PASS 42／SKIP 1 | **`rc=1`**，FAIL 2（红句 `tag=65534`） | `rc=0` | `rc=0` |
| `:218` | `3`→`2` | `rc=0`，FAIL 0／PASS 42／SKIP 1 | **`rc=1`**，FAIL 2（红句 `tag=3`） | `rc=0` | `rc=0` |
| `:196` | 起点 `+24`→`+26`（连同界 `+26`→`+28`） | `rc=0`，FAIL 0／PASS 42／SKIP 1 | **`rc=1`**，FAIL 2（红句 `tag=0`） | `rc=0` | `rc=0` |
| `:196` | 派单字面（只动起点） | `rc=0`，FAIL 0／PASS 42／SKIP 1 | **`rc=1`**，FAIL 1＋panic（空切片喂 `Uint16`） | `rc=0` | `rc=0` |

三枚**逐字节还原**之外，本腿另交一把终局尺：四发全跑完后 `cp pristine` 再 `cmp` 真仓工作树 ⇒
`rc_final_cmp_pristine=0`／`rc_final_cmp_vs_real_worktree=0`（件＝`logs/mutate-final-safety.txt`）⇒ **真仓的产码面⛔ 被碰过一枚字节**。

## 6. 本腿自己那把尺的一处缺陷（具名，⛔ 藏着）

驱动脚本里"往同一枚件里 `echo` 汇总、又在 `{ } >> 同一枚件` 里**直接** `grep` 那枚件"的那两处（`-03-red.txt` 与 `-05-regreen.txt`
的尾部）被 GNU grep 拒了，逐字＝`grep: input file ‘…’ is also the output`，那两行汇总因此⛔ 是名册而是报错句。
读数本身⛔ 受影响：① 原始测试输出在 grep 之前已经落件（`--- FAIL`／`--- PASS`／`ok`／`panic` 逐字俱在）；
② 每发的名册与逐字红句本腿**另起一把尺**重导进 `mut-<名>-SUMMARY.txt`（读⛔ 同一枚件、写另一枚，故⛔ 触发该拒绝）；
③ `$(grep -c …)` 那几枚计数走的是命令替换的管道、⛔ 是直接重定向，所以 `top_fail`／`top_pass`／`top_skip` 三枚数是**真读数**。
⇒ 定式一条交给裁决者：**汇总⛔ 写进它自己要 grep 的那枚件**（本腿下次先落 `-raw`、再另起一把尺导 SUMMARY）。
