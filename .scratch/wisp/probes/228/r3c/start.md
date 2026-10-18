# 228-r3c 起手件（落点更正＋起手三把尺现量）2026-10-18 第一波

腿＝`228-r3c`（票 228 突变复跑腿，只仓外导出树）。任务＝把 `R-228-3b-1`（台账 `A1019` 第③节）那两发突变
在导出树上**串行**重做成彼此独立的两发读数。派单存档＝
`D:\work\workspace\projects plans\Wisp\.scratch\wisp\probes\orch\20261018\txt\228r3c-dispatch.txt`（56 行，已 Read 全文）。

## 1. 起手锚（现跑 `git log -1`，⛔ 抄旧数）

- `git log -1 --format='%H %ad %s'` ＝ `c11163c896b09e7b6863d347dadbcb65702585cc  Sun Oct 18 10:02:56 2026 +0800`（branch dev）。
- 派单写"编排者 09:3x 见 `f60d50d6`"、会话头写"10:0x 见 `c11163c8`" ⇒ 本腿起手锚＝**我自己现认的 `c11163c8`**（派单第 2 行明写"起手锚＝你自己 `git log -1` 现认"）。
- 起手侧跟踪面尺＝`git status --porcelain` 剥 `??` 后 **0 枚**（未跟踪历史临时件按"只建不删"留着，共 7 枚 `??` 里含 `.go` 的都在 `.scratch/**` 下，非跟踪）。

## 2. ★落点更正（机主 2026-10-18 10:1x 亲口新铁律，覆盖派单 §台面配方的盘符根形状）

⛔ 任何导出树／`GOTMPDIR`／`TMP`/`TEMP`／工作目录建在 `D:\`、`C:\` 等盘符根下。
`228-r3b` 的 `D:\wisp228r3b`（8.3＝`D:\WISP22~1`）那一形从现在起＝违规，本腿⛔ 复用、⛔ 删。

✅ 本腿唯一落点＝`D:\tmp\wisp\20261018\228r3c\`：
- `tree\`＝仓外导出树（`/d/tmp/wisp/20261018/228r3c/tree`）
- `gotmp\`＝`GOTMPDIR`（先 `mkdir -p`）
- `tdst\`＝surface d 的 `TMP/TEMP/TMPDIR`（先 `mkdir -p`）
- `logs\`＝逐枚读数件（＝`.scratch/wisp/probes/228/r3c/logs/**` 的仓外原件，入库时逐枚拷进写面）

### 短名尺（现跑三把，逐字）

1) `fsutil 8dot3name query D:` ＝ **`Failed to open the volume. Error 5: Access is denied.`**
   ⇒ 本腿非管理员且本会话审批已关＝**尺读不到，不是"短名不生成"**；⛔ 提权重试，改用下面两把替代尺。
2) `cmd /c dir /X D:\tmp\wisp\20261018\228r3c`（先在该处现建一枚长名目录 `probe-longdirname-228r3c` 再复读）＝
   吐出别名 **`PROBE-~1`** ⇒ **D: 卷在 `D:\tmp\wisp\...` 这条新落点下确实生成 8.3 短名**（不是只有盘符根才生成）。
   ⚠ 该探针目录＝临时件，按"只建不删"留着不动。
3) `cygpath -d D:\tmp\wisp\20261018\228r3c\tree` ＝ **原样返回**（`tmp`/`wisp`/`20261018`/`228r3c`/`tree` 每段都 ≤8 字符
   ⇒ 短名＝长名，路径解析不炸）；对照 `cygpath -d D:\wisp228r3b` ＝ `D:\WISP22~1`。

⇒ 短名硬需求（`internal/winsec`／`paths_shortname_252` 一族与 `cmd/wisp` 真窗件按 8.3 短名办事）**在新落点下成立**，
不必申请"专用树根"例外。已回 lead 一句确认。

## 3. pristine 基准（现算，⛔ 抄 `3e2e83fa…` 那枚旧数）

尺＝对起手锚 `c11163c8` 上的 `internal/proc/boot_windows.go`：
- `git rev-parse c11163c8:internal/proc/boot_windows.go` ＝ blob sha1 **`531060a55c080ebd535b426b8c60dcd9011785fc`**
- `git cat-file blob c11163c8:internal/proc/boot_windows.go | sha256sum` ＝ **`3e2e83fab184ce91640ad7dbc06f8dd116baffe325bd47aa46af2527aad73fe7`**
  ⚠ 第一把我用 `git show … | Out-File -NoNewline` 抽，读出 `c3890203de9f…`＝**PowerShell 把 LF 换成了 CRLF＝我那把尺坏**；
  按"两把尺矛盾先假定新那把坏"换 `git cat-file`（不经 PowerShell 编码层）后与导出树物理字节**逐字同串** ⇒ 基准尺＝后者。
- 导出树物理字节（`git archive c11163c8 … | tar -x`）`sha256sum` ＝ **`3e2e83fab184ce91640ad7dbc06f8dd116baffe325bd47aa46af2527aad73fe7`**
  ⇄ `git hash-object` ＝ `531060a55c08…` ⇄ `bytes=12916` ⇄ `lines=313`
  （该件 `.gitattributes` 是 `text eol=lf` ⇒ `git ls-files --eol`＝`i/lf w/lf`，导出无换行腐，两把尺天然对齐）
- ⚠ 诚实记账：上面那枚 `3e2e83fa…` **与 `228-r3b`/`A1019` 引用的旧串逐字相同**——那是**因为我这发导出的是同一枚 blob**
  （`git log -1 -- internal/proc/boot_windows.go` ＝ `f4d4fda9  Sat Oct 17 11:09:42 2026  票 228 AC#11 产码①（internal/proc，修法乙-1）`，
  在我起手锚 `c11163c8` 之前、之后无新笔）。数是我自己现算的（三条尺逐字输出在 `logs/hashes-pristine.txt`），⛔ 抄；相同属构造使然。
  **本腿每发的对撞尺＝`git hash-object` 那把 sha1 ＋物理字节 sha256 两把同看**，任一不符即判未还原。
- ⚠ 起手锚漂移具名：落这件时盘上 HEAD 已被别腿推到 `5db0b538`（`票 337 普查腿 337-a1·AC#0 起手…`，10:13:17）。
  本腿台面**仍按起手锚 `c11163c8` 导出的那棵树**跑，⛔ 中途换锚；跟踪面 `internal/proc/boot_windows.go` 本身在两锚之间未动（上面那把"末笔"尺证）。

## 4. 台面（派单 §台面配方，逐味）

`git archive c11163c8 go.mod go.sum internal cmd tools scripts frontend models`（八件套）＋仓根 `docs/`
＋手工三件 `build/*.dll`（树内 3 枚 dll）／`third_party/sherpa-onnx/*.dll`（3 枚）／`frontend/dist`（`assets`＋`index.html`）。
取出后 `.go` 数＝**760**（尺＝`find $TREE -name '*.go' -type f | wc -l`，非 0 ⇒ 台子非空）；`docs/PLAN.md` 在位。
全程 `PATH="$PWD/third_party/sherpa-onnx:$PWD/build:$PATH"`（防 `0xc0000135`）。
⛔ 在共享工作树跑 `go test`／`go build`；⛔ `-overlay`；⛔ `go test ./...`；⛔ 跑 `cmd/wisp` 整包（本腿只定向 `-run`）。

## 5. 两发落点的内容锚（现读 `internal/proc/boot_windows.go`，⛔ 抄行号当锚）

- **M1**＝环路 stop 检查那一行。函数体＝`func (rt *Runtime) RunEventLoop() string`；
  内容锚 token＝**`stopSignal := rt.StopSignal()`**（该赋值是 `select` 里 `case <-stopSignal:` 唯一喂料者）。
  突变形＝把它改成 `stopSignal := (<-chan struct{})(nil)` ⇒ nil 通道永不 ready ⇒ 记录下的请求再也转不成 return。
  点名用例＝**`TestAC11TrayExitTriggerWalksTheTenSteps`**（ⓑ-1）。
- **M2**＝Boot 建停机位那一行。函数体＝`func Boot(...)`；
  内容锚 token＝**`rt.stopReq = newStopRequest()`**（唯一建位处，`grep -n 'stopReq = newStopRequest' internal/proc/boot_windows.go`＝1 命中）。
  突变形＝整行注释掉 ⇒ `StopSignal()` 答 nil。
  点名用例＝**`TestAC11StopBitIsBuiltAtBoot`**（ⓐ-first）。
- 两发落点同一枚文件、不同 token ⇒ 成对两发落点不同✓。互不污染自证＝每发 logs 带三段哈希＋对整文件哈希证另一枚落点未被动。
