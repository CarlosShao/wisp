# 352-r1 · AC#1 落地那一行 ⇄ 改后读数（判据＝两分类并存）

落地笔＝`ecffb93be03bd86f647c1de468783bcef35a9fc9`（主题起头逐字＝「票 352 落地腿 352-r1」）；
`winsec.go` blob＝改前 `07300bfffc3e5b14c017029d6f84fb5a7ede3db8` → 改后 `f3c1b68ee60d8d5d42b816e47633d31f1ba191f5`；
仓外树 sha256＝基线 `9086f2bb…f2da6f2331` ⇄ 落地 `2f2887ab…7df69abd`（＝`treeL` 由 `git archive ecffb93b` 现建，`.go` 枚数 777 非空尺，`FMT_V=0`⇄`FMT_W=1` 内容锚现量）。
那一行的 `git show -U0` 全貌（＝⛔ 别枚字节）：

```
@@ -431 +431 @@ func wrapPath(path string, err error) error {
-	return fmt.Errorf("%w: %s: %v", ErrNotSealable, path, err)
+	return fmt.Errorf("%w: %s: %w", ErrNotSealable, path, err)
```

尺＝与 AC#0 **同一枚**探针件、同一枚驱动、同一对树槽（`sh/02-errclass-pair.sh`：slot1＝treeB 基线，slot2＝treeL 落地；两枚独立 mut 拷贝＝⛔ 后一发覆盖前一发）。
★判据＝改后同一跳 `errors.Is(err, fs.ErrNotExist)`＝**true** 且 `errors.Is(err, ErrNotSealable)` **照旧 true**＝两分类并存。

## 表 · 改前 ⇄ 改后四值对照（`fs.ErrNotExist` / `ErrNotSealable`）

| 平台 | 跳（尺名） | 改前 ErrNotExist / NotSealable | 改后 ErrNotExist / NotSealable | 判据 | 发 |
|---|---|---|---|---|---|
| POSIX（容器） | R1 `PrivateFile` 父目录不存在 | **false / true** | **true / true** | ★并存 | 2 发逐字同 |
| POSIX（容器） | R2 `PrivateFileExclusive` 同跳 | false / true | **true / true** | ★并存 | 2 发 |
| POSIX（容器） | R3 `wrapPath` 本体直打可穿 sentinel | false / true | **true / true** | ★并存 | 2 发 |
| POSIX（容器） | R4 `SealFile` 撞缺失叶子 | false / true | **true / true** | ★并存 | 2 发 |
| POSIX（容器） | R5 `SealFile` 父目录不存在 | false / true | **true / true** | ★并存 | 2 发 |
| POSIX（容器） | R6 `SealDir` 父目录不存在 | false / true | **true / true** | ★并存 | 2 发 |
| POSIX（容器） | R13 注入 seal 失败（⛔ 可穿 errno） | false / true | **false / true** | ★未被改成第三种分类 | 2 发 |
| Windows（本机） | R3 同上直打 | false / true | **true / true** | ★并存 | 2 发 |
| Windows（本机） | R4 `SealFile` 撞缺失叶子 | false / true | **true / true** | ★并存 | 2 发 |
| Windows（本机） | R5 `SealFile` 父目录不存在 | false / true | **true / true** | ★并存 | 2 发 |
| Windows（本机） | R6 `SealDir` 父目录不存在 | false / true | **true / true** | ★并存 | 2 发 |
| Windows（本机） | R1⇄R2 `PrivateFile*` 父目录不存在 | false / true（`isUnresolved=true`） | **逐字不变** false / true（`isUnresolved=true`） | ★⛔ 本票射程（走 `pinRefusal`→`preOpenRefusal`，自带双 sentinel） | 2 发 |
| Windows（本机） | R13 注入 seal 失败 | false / true | **false / true** | ★未换成第三种分类 | 2 发 |

⇒ **AC#1 判据成立**：Windows 半边靠 R3⇄R4⇄R5⇄R6 那四枚真经 `wrapPath` 的跳量到"两分类并存"（★票 352 现量第 10 行逐字已写"Windows 同跳 落地另带 `isUnresolvedPath=true`"＝同一形状，本腿独立复现）；
R1⇄R2 在 Windows 那一跳⛔ 经 `wrapPath`＝两侧逐字同＝**改前⇄改后没被本票动过**（＝⛔ 第二处、⛔ 顺手扩）。
POSIX 半边 R1..R6 全绿到 true/true。

## 控制跳两侧逐字同（＝⛔ 动任何其他跳分类）

R7 覆盖写＝`err=nil`（两侧）／R8 独占撞已存在＝`isErrExist=true`＋`PathError{Op:open errno:file exists}`（两侧）／
R9 空串＝`isUnresolved=true`、`isNotSealable=false`（两侧）／R10 相对拼写＝POSIX 两侧逐字同、Windows 两侧逐字同／
R11 `SealFile` 命中存在文件＝`err=nil`（两侧）／R12 叶子＝符号链接＝POSIX 两侧 `isUnresolved=true`、Windows 两侧 `risk: path traverses a reparse point …`（⛔ sentinel）／
R13 注入 seal 失败＝两侧都 `isNotSealable=true` 且其余 false。
⇒ ★`sealError` 幂等判定⛔ 动（R13 那枚"已带 `ErrNotSealable` 的原样返回"两棵树逐字同）。

## 逐字落地行（POSIX 容器档，`logs/ac1-col-s2-mut-2-treeL.raw.txt`，未折未删）

```
    probe_352r1_err_test.go:51: R1 parent-missing PrivateFile: isErrNotExist=true isNotSealable=true isErrExist=false isUnresolved=false isReparse=false PathError{Op:openat errno:no such file or directory} verbatim="winsec: cannot apply a private security descriptor: /tmp/TestP352R1ErrClassify2347502323/001/no-such-dir: openat no-such-dir: no such file or directory"
    probe_352r1_err_test.go:54: R3 wrapPath direct fs.ErrNotExist: isErrNotExist=true isNotSealable=true isErrExist=false isUnresolved=false isReparse=false PathError{Op:<no *os.PathError> errno:<n/a>} verbatim="winsec: cannot apply a private security descriptor: the/spelling: file does not exist"
    probe_352r1_err_test.go:56: R4 SealFile missing-leaf in existing dir: isErrNotExist=true isNotSealable=true isErrExist=false isUnresolved=false isReparse=false PathError{Op:openat errno:no such file or directory} verbatim="winsec: cannot apply a private security descriptor: /tmp/TestP352R1ErrClassify2347502323/001/no-such-leaf.bin: openat no-such-leaf.bin: no such file or directory"
```

## 逐字落地行（Windows 本机档，`logs/ac1-host-s2-mut-2-treeL.raw.txt`，未折未删）

```
    probe_352r1_err_test.go:54: R3 wrapPath direct fs.ErrNotExist: isErrNotExist=true isNotSealable=true isErrExist=false isUnresolved=false isReparse=false PathError{Op:<no *os.PathError> errno:<n/a>} verbatim="winsec: cannot apply a private security descriptor: the/spelling: file does not exist"
    probe_352r1_err_test.go:56: R4 SealFile missing-leaf in existing dir: isErrNotExist=true isNotSealable=true isErrExist=false isUnresolved=false isReparse=false PathError{Op:<no *os.PathError> errno:<n/a>} verbatim="winsec: cannot apply a private security descriptor: C:\\Users\\swq\\AppData\\Local\\wisp-test-gotmp\\352r1-20261020-161608-201\\TestP352R1ErrClassify3243477728\\001\\no-such-leaf.bin: The system cannot find the file specified."
    probe_352r1_err_test.go:57: R5 SealFile parent-missing: isErrNotExist=true isNotSealable=true isErrExist=false isUnresolved=false isReparse=false PathError{Op:<no *os.PathError> errno:<n/a>} verbatim="winsec: cannot apply a private security descriptor: C:\\Users\\swq\\AppData\\Local\\wisp-test-gotmp\\352r1-20261020-161608-201\\TestP352R1ErrClassify3243477728\\001\\no-such-dir\\leaf.bin: The system cannot find the path specified."
```

★第二发（复量）＝`logs/ac1-col2-*.raw.txt` ⇄ `logs/ac1-host2-*.raw.txt`，同一枚尺名逐值同（★竞态读数⛔ 要求路径逐字同，只认量级与形状＝此处连分类位都逐字同）。

★⛔ 突变＝本票不需要去牙：同一把尺在两棵树上读出**不同**的分类（`false/true` ⇄ `true/true`）＝尺本身有牙；★成对两发落点⛔ 同＝slot1/slot2 各自独立 mut 树，且每发跑后 `RESTORE=OK`（sha256 前后对回）。
