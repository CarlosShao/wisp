# 352-r1 · AC#0 改前读数（现状钉死）＝票 352 AC#1 的对照物

★本文件只登记**改前**读数＝落地那一行之前共享树 `internal/winsec/winsec.go` 的字节。
起手现认锚（本腿第一笔 `git log -1`）＝`bf509bb837aa1328617320608be38e4bca5a5d35`；
`winsec.go` 在该锚的 blob＝`07300bfffc3e5b14c017029d6f84fb5a7ede3db8`；工作树字节＝同值（`git hash-object` 现量）；
仓外导出树 `treeB` 里那枚文件的 sha256＝`9086f2bbbf531823ed5c4a615391785f5af749542df15c046b5fb4f2da6f2331`。

尺（同一枚探针件逐字打进两棵树的两枚独立拷贝）＝`.scratch/wisp/probes/352/r1/probe/probe-352r1-errclass_test.go.txt`
＋驱动 `.scratch/wisp/probes/352/r1/sh/02-errclass-pair.sh`。
⛔ 任何断言、只 `t.Logf`＝一把尺，不是一枚钉；件⛔ 入库（只存在于 `D:\tmp\wisp\20261020\352r1\mut-*` 的拷贝里）。

## 台面

- 基线侧＝`treeB`＝`git archive bf509bb8`（八件套＋`docs/`＋手工补 `build/*.dll`⇄`third_party/sherpa-onnx/*.dll`⇄`frontend/dist`）；`.go` 枚数＝777 非空尺。
- POSIX（容器）＝`golang:1.27`＋`CGO_ENABLED=0`＋`GOPROXY=off`＋宿主 module 缓存只读；`TMPDIR`⇄`GOTMPDIR`＝容器自 `/tmp`（★⛔ 绑 Windows 卷）。
- Windows（本机）＝`go1.27.1 windows/amd64`；`GOTMPDIR`⇄`TMPDIR`＝反斜杠形且落 `%LOCALAPPDATA%` 同卷短根（逐字见每发的 `GOTMPDIR_USED=`）；★⛔ `cygpath -wp` 转整条 PATH。
- 逐字发件＝`logs/ac0-host-base-s{1,2}-*.raw.txt` ⇄ `logs/ac0-col-base-s{1,2}-*.raw.txt`；四数＝`logs/four-counts-errclass.tsv`。
- 首发两枚**作废发**也入库＝`logs/ac0-host-base1-*.raw.txt`（探针件没拷进仓外暂存 ⇒ `RUN=0`＝台损）⇄ `logs/ac0-col-base1-*.raw.txt`（msys 把 `/src` 换成 `E:\work\base\msys64\src` ⇒ `RC=125`／`RUN=0`＝台损）。★两枚都⛔ 计入判据，只当"这味会真咬人"的具名记录。

## 表 · 点名那一跳（父目录不存在）改前两个真值

尺＝`errors.Is`；每格＝`fs.ErrNotExist` ⇄ `ErrNotSealable`。

| 平台 | 跳（探针尺名） | `errors.Is(err, fs.ErrNotExist)` | `errors.Is(err, ErrNotSealable)` | 发 |
|---|---|---|---|---|
| POSIX（容器） | `PrivateFile`＝R1 | **false** | **true** | 2/2 逐字同 |
| POSIX（容器） | `PrivateFileExclusive`＝R2 | **false** | **true** | 2/2 逐字同 |
| POSIX（容器） | `SealFile` 撞缺失叶子＝R4 | **false** | **true** | 2/2 |
| POSIX（容器） | `SealFile` 父目录不存在＝R5 | **false** | **true** | 2/2 |
| POSIX（容器） | `SealDir` 父目录不存在＝R6 | **false** | **true** | 2/2 |
| POSIX（容器） | `wrapPath` 本体直打一枚可穿 sentinel＝R3 | **false** | **true** | 2/2 |
| Windows（本机） | `PrivateFile`＝R1 | false | true（★此跳在本平台⛔ 经 `wrapPath`，见下节） | 2/2 |
| Windows（本机） | `PrivateFileExclusive`＝R2 | false | true（同上） | 2/2 |
| Windows（本机） | `SealFile` 撞缺失叶子＝R4 | **false** | **true** | 2/2 |
| Windows（本机） | `SealFile` 父目录不存在＝R5 | **false** | **true** | 2/2 |
| Windows（本机） | `SealDir` 父目录不存在＝R6 | **false** | **true** | 2/2 |
| Windows（本机） | `wrapPath` 本体直打＝R3 | **false** | **true** | 2/2 |

⇒ **AC#0 成立**：两平台上凡经 `wrapPath` 的那几跳，`fs.ErrNotExist` 全＝**false**、`ErrNotSealable` 全＝**true**。
★AC#1 的判据形状由此钉死＝同一枚尺改后必须读到 **两分类并存**（`fs.ErrNotExist` true **且** `ErrNotSealable` true）。

## 一处具名照实（⛔ 扩射程、⛔ 自行假设）

Windows 侧 R1⇄R2（`PrivateFile`/`PrivateFileExclusive` 的"父目录不存在"）**不落在 `wrapPath`**：
`pinnedParent` 的 Windows 支在 `cur.OpenRoot(comp)` 那跳失败 ⇒ 走 `pinRefusal`→`placementRefusalWindows`→`preOpenRefusal`，
自带双 sentinel，逐字 `isUnresolved=true isNotSealable=true`。＝票 352 现量第 11 行那句"那支已自带双 sentinel＝⛔ 第二处"的同一形状，
本腿⛔ 动它、⛔ 把它算进本票射程。**⇒ Windows 半边的 AC#1 判据取 R3⇄R4⇄R5⇄R6 这四枚真经 `wrapPath` 的跳**（POSIX 半边 R1..R6 全经 `wrapPath`）。

## 逐字红/绿句（改前基线，两平台各一枚代表行；完整 13 行见 `logs/`）

POSIX 容器档（`logs/ac0-col-base-s1-mut-1-treeB.raw.txt`，逐字）：

```
    probe_352r1_err_test.go:51: R1 parent-missing PrivateFile: isErrNotExist=false isNotSealable=true isErrExist=false isUnresolved=false isReparse=false PathError{Op:<no *os.PathError> errno:<n/a>} verbatim="winsec: cannot apply a private security descriptor: /tmp/TestP352R1ErrClassify2275833704/001/no-such-dir: openat no-such-dir: no such file or directory"
    probe_352r1_err_test.go:54: R3 wrapPath direct fs.ErrNotExist: isErrNotExist=false isNotSealable=true isErrExist=false isUnresolved=false isReparse=false PathError{Op:<no *os.PathError> errno:<n/a>} verbatim="winsec: cannot apply a private security descriptor: the/spelling: file does not exist"
    probe_352r1_err_test.go:56: R4 SealFile missing-leaf in existing dir: isErrNotExist=false isNotSealable=true isErrExist=false isUnresolved=false isReparse=false PathError{Op:<no *os.PathError> errno:<n/a>} verbatim="winsec: cannot apply a private security descriptor: /tmp/TestP352R1ErrClassify2275833704/001/no-such-leaf.bin: openat no-such-leaf.bin: no such file or directory"
```

Windows 本机档（`logs/ac0-host-base-s1-mut-1-treeB.raw.txt`，逐字）：

```
    probe_352r1_err_test.go:51: R1 parent-missing PrivateFile: isErrNotExist=false isNotSealable=true isErrExist=false isUnresolved=true isReparse=false PathError{Op:<no *os.PathError> errno:<n/a>} verbatim="winsec: path is not provably resolved, refusing to seal: C:\\Users\\swq\\AppData\\Local\\wisp-test-gotmp\\352r1-20261020-160844-1453\\TestP352R1ErrClassify1725796622\\001\\no-such-dir\\leaf.bin: refusing to open for a private write: pinned object at no-such-dir is not what the name carries: the walk could not pin that component through the handle of its parent: The system cannot find the path specified. (winsec: cannot apply a private security descriptor)"
    probe_352r1_err_test.go:54: R3 wrapPath direct fs.ErrNotExist: isErrNotExist=false isNotSealable=true isErrExist=false isUnresolved=false isReparse=false PathError{Op:<no *os.PathError> errno:<n/a>} verbatim="winsec: cannot apply a private security descriptor: the/spelling: file does not exist"
    probe_352r1_err_test.go:56: R4 SealFile missing-leaf in existing dir: isErrNotExist=false isNotSealable=true isErrExist=false isUnresolved=false isReparse=false PathError{Op:<no *os.PathError> errno:<n/a>} verbatim="winsec: cannot apply a private security descriptor: C:\\Users\\swq\\AppData\\Local\\wisp-test-gotmp\\352r1-20261020-160844-1453\\TestP352R1ErrClassify1725796622\\001\\no-such-leaf.bin: The system cannot find the file specified."
```

## 控制跳（两棵树之间⛔ 该被本票改动＝AC#1 落地后须逐字同）

R7 覆盖写 err=nil／R8 独占撞已存在＝`isErrExist=true`＋`PathError{Op:open errno:file exists}`／R9 空串＝`isUnresolved=true`／
R10 相对拼写：POSIX＝`isUnresolved=true`（`is not absolute`）⇄ Windows＝`isUnresolved=true isNotSealable=true`（走 pin 支）／
R11 `SealFile` 命中存在文件＝err=nil／R12 叶子＝符号链接：POSIX＝`isUnresolved=true`（resolve 面先答）⇄ Windows＝`risk: path traverses a reparse point …`（resolver 抢先，⛔ sentinel）／
R13 注入 seal 失败＝`isNotSealable=true` 且其余全 false＝★落地后仍须如此＝证明"⛔ 换成另一种不可穿、也⛔ 造出第三种分类"。
