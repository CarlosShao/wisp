# 352-r1 续跑（cont）· 票面现量第 9 行欠账＝`wrapPath` 调用点逐枚现数并点名

★这味＝票 352 现量第 9 行逐字留的账：「`git cat-file blob HEAD:internal/winsec/winsec.go | grep -c 'wrapPath('`＝**2**（含定义行；实际调用点少，**落地腿逐枚现数并点名**）」。
原 9 笔（`0fb399c3`..`8a10a7e3`）盘上⛔ 有这一味（尺＝`grep -rln '调用点' .scratch/wisp/probes/352/r1`＝命中 0 枚，见 `logs/10-blob-cr.txt` F 段）＝射程缺口；本续跑腿补齐。
★本件只登记现数尺＋点名，⛔ 任何 AC 盖章（四格归非实现者裁）、⛔ 动任何产码字节（那一行已落＝`ecffb93b`）。

## 尺① 票面那把（winsec.go 单文件）逐字重跑（量的）

```
git cat-file blob HEAD:internal/winsec/winsec.go | grep -c 'wrapPath('   → 2
git cat-file blob bf509bb8:internal/winsec/winsec.go | grep -c 'wrapPath(' → 2   （⛔ 随落地漂＝那一行没增删调用点）
```
逐字两枚命中＝`:430 func wrapPath(...)`（定义）＋`:445 return wrapPath(path, err)`（`sealError` 体内唯一调用）。
⇒ 票面「2（含定义行）」对回＝winsec.go 内**调用点＝1 枚**（`sealError`）。

## 尺② 全仓调用点逐枚点名（尺＝`git grep -n 'wrapPath(' HEAD -- 'internal/winsec/*.go'`，剥定义行；量的）

定义＝1 枚（`winsec.go:430`）；调用点＝**24 枚**，逐枚 file:line：

| 文件 | 调用枚数 | 逐枚行号 |
|---|---|---|
| `winsec.go` | 1 | 445 |
| `winsec_other.go`（POSIX 支） | 11 | 167, 173, 179, 192, 198, 217, 257, 261, 264, 294, 298 |
| `winsec_windows.go`（Windows 支） | 7 | 298, 306, 540, 543, 547, 571, 574 |
| `pinnedleaf_342_windows.go` | 4 | 79, 91, 103, 130 |
| `pinnedleaf_342_fallback.go` | 1 | 40 |

⇒ `CALLSITES_PROD=24`（剥定义 1 枚）；`_test.go` 里命中＝**0**（尺＝`git grep -c 'wrapPath(' HEAD -- 'internal/winsec/*_test.go'`＝0 命中＝分母不含测试件）。起手锚 `bf509bb8` 上同尺＝`{1,11,7,4,1}` 逐枚同数＝调用点⛔ 被本票改动（只改了函数体内那一枚格式串的第二个动词）。

## 尺③ 这 24 枚的归类（＝票面现量第 10/11 行「七跳两侧逐字同」+「⛔ 第二处」的可核化）

- **经 `wrapPath` 的跳**＝POSIX 半边（`winsec_other.go` 11 枚调用）＋Windows 半边 `SealFile`/`SealDir`/`pinnedleaf` 里真落到 `wrapPath` 的那几枚 ⇒ AC#1 落地后这些跳都读出 `fs.ErrNotExist=true` 且 `ErrNotSealable=true`（两分类并存，逐字见 `readings-ac1-landed.md`）。
- **不经 `wrapPath` 的跳**＝Windows `PrivateFile`/`PrivateFileExclusive` 的"父目录不存在"（R1⇄R2）：走 `pinRefusal`→`placementRefusalWindows`→`preOpenRefusal` 自带双 sentinel（票面第 11 行逐字"⛔ 第二处"）⇒ 落地侧仍 `false/true`＋`isUnresolved=true` 逐字不变（尺＝`logs/11-ac1-pername.txt`，Windows 落地那 4 枚 false/true＝R1/R2/R10/R13，⛔ 一枚经改动而翻转分类）。

## 现量结论（一句话给编排者收件）

票面现量第 9 行的欠账（调用点逐枚现数并点名）＝**24 枚调用点＋1 枚定义**，逐枚已点名（尺②表）；起手锚与落地锚同数＝那一枚格式串改动⛔ 增删任何调用点＝射程守住。本件＝原 9 笔缺的那一味，⛔ 翻任何 AC 框。
