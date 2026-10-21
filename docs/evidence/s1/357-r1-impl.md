# 票 357 落地腿 `357-r1` — 停手上报（写腿）

起手锚 `650fac07`。**本腿没有落地任何产码**：`AC#3` 的第一块砖撞上派单自己写的停手条款，
本件把"撞在哪、我量到了什么、请编排者裁什么"逐条落盘，供下一腿**一轮开工**。

写作约定：禁码位写作 sec N / "=> " / NOT / (1)(2)(3)。跑门禁前 `unset GOTMPDIR TMPDIR TEMP TMP`、`GOFLAGS=`，未设 `CGO_ENABLED=0`。

## 1 停在哪：派单"未定义即停"第 (2) 条，逐字成立

派单 `:70-71` 逐字：**(2)若你发现『枚举 ADS 必须动 `internal/winsec`』=> 停手上报**；
派单 `:24` 逐字：**"若不在 => 停手上报（NOT 自己加依赖）"**，条件挂的是 `FindFirstStreamW`/`FindNextStreamW`
是否已有封装。这两条我都没踩（`internal/winsec` **不需要动**、`golang.org/x/sys` **已在** `go.mod:12` v0.48.0）。
**真正撞上的是选型本身**：`x/sys/windows` **没有** `FindFirstStreamW`/`FindNextStreamW` 的封装。

### 1.1 尺与读数（同一把尺，带正控）

台子：`go env GOMODCACHE` = `D:\work\base\gopath\pkg\mod`，目录 `golang.org\x\sys@v0.48.0\windows`
（`go list -m golang.org/x/sys` 现读 = `golang.org/x/sys v0.48.0`，`go.mod:12` 逐字 v0.48.0）。
尺＝对该目录 `*.go` 全文 `Select-String -SimpleMatch`，逐符号计数。

| 符号 | 命中 | 判 |
|---|---|---|
| `FindFirstStreamW` | **0** | 派单点名的 API，x/sys **无封装** |
| `FindNextStreamW` | **0** | 同上 |
| `WIN32_FIND_STREAM_DATA` | **0** | 配套结构体也没有 |
| `FILE_STREAM_INFO` | **0** | （名称写法，见下行） |
| **正控** `func CreateFile` | **2** | 尺活 |
| **正控** `func FindClose` | **2** | 尺活 |
| **正控** `func GetFileInformationByHandleEx` | **1** | 尺活（见 sec 2） |
| **正控** `FileStreamInfo`（常量） | **1** | 尺活（见 sec 2） |

旁证两把独立尺，同一答案：
- 全体 27 枚 `*.go` 里 `-SimpleMatch 'streamw'` => 命中 **1**，逐字是 `zerrors_windows.go` 的
  `WINCODEC_ERR_STREAMWRITE`（错误码常量，NOT API 封装）。
- GOROOT `D:\work\base\go\src\syscall`（`zsyscall_windows.go` + 全部 `*windows*.go`）：
  正控 `procCreateFile` **4**、`procFindClose` **2**（尺活）；目标 `FindFirstStreamW` **0**、
  `FIND_STREAM_DATA` **0**（标准库也没有）。

=> 要用 `FindFirstStreamW`，唯一办法是**自己在 `internal/models` 里手写 `NewLazySystemDLL`+`NewProc`
+手映 `WIN32_FIND_STREAM_DATA`**。派单 `:23-26` 把它当作"用 `golang.org/x/sys/windows` 里**已有的**
`FindFirstStreamW`/`FindNextStreamW`"来派的 —— **这个前提今天不成立**，它不是"换个写法"，
而是**新增手写 syscall 面**。派单 `:70` 的停手条款正是为这种情形写的，所以我停手，NOT 自行改选型。

## 2 我已经把"另一条不需要手写 syscall 的路"量通了（供裁决，NOT 已落地）

停手不该只递一张"走不通"。我用 `x/sys/windows` **已经导出**的那枚函数做了仓外实测：
`GetFileInformationByHandleEx(h, windows.FileStreamInfo, ...)`（`zsyscall_windows.go:2432` 逐字有封装，
常量 `FileStreamInfo = 7` 在 `types_windows.go:2765` 逐字有）。这条路**只手映一枚数据布局**
（`FILE_STREAM_INFO` 24 字节头），**不写任何 syscall**，与本仓既有做法同形
（`internal/ball/win32_windows.go:160` 起 `wndClassEx`/`msg`/`point` 等一串手映结构体 + `unsafe.Sizeof`，
正控：`internal/ball/tray_windows.go:33` 等 8 处 `unsafe.Sizeof`）。

台子＝仓外 `D:\tmp\wisp\20261021\357r1\probe\`（真 Go 程序，真 NTFS 卷 `D:`，零 mock）。
真件：`dataG` 只有主体；`dataD1`＝主体 + `model.onnx:attacker` **19 B**（真 ADS）。

| 发 | 形状 | 读数（逐字） |
|---|---|---|
| G | 只有主体 | `all=["::$DATA"] named=[] err=<nil>` => **0** 枚 |
| D1 | 主体 + ADS 19 B | `all=["::$DATA" ":attacker:$DATA"] named=[":attacker:$DATA"] err=<nil>` => **1** 枚，流名**点了名** |

边界两面都量过（这两条决定"能不能直接抄去用"）：
- **单次缓冲不够**：4096 B 缓冲在 **200 条**短流时返回 `err = More data is available.`（`all=[]`）。
  加"长大重试"后 200 条**全出**（`err=<nil>`，`named` 200 枚）；40 条与 200 字符长流名都 `<nil>`。
  => 落地时**必须**带 `ERROR_MORE_DATA` 重试，否则"流多"会退化成"看不见"，即**假阴**。
- **目录句柄**返回空列表且 `err=<nil>`（不是错误）=> 只能逐文件调，不能拿目录一把梭。

**成本读数**：这是 `internal/models` 内一枚只读函数 + 一枚 24 字节结构体，**零新依赖、零手写 syscall、
零 `internal/winsec` 改动**、非 Windows 侧 `//go:build !windows` 恒返回空（派单 `:26` 那一条照做）。

## 3 停手前我把"开工前"三项与改前读数都做完了（下一腿可直接用）

### 3.1 `git status --porcelain`：写面内零别人的改动
`git status --porcelain -- internal/models cmd/wisp docs/evidence/s1` => **空**（零行）。
（全仓那 60+ 行未跟踪件是历史 `.scratch/commit-msg-*` 等，**不在写面内**。）

### 3.2 票面前提新鲜度（派单 `:37-40` 三条，逐条现跑，三处腐坏都成立）
- 票面 `:9` 写"签名＝`109-v2` 判语**第 2 条**" => 实为 **sec 5(c)**。现读
  `HEAD:docs/evidence/s1/109-v2-verdict.md:109` 逐字
  `**(c) 建议：并进"第一个模型装载器"那张票的 AC，NOT 另立一票、也 NOT 并本票。**` => **腐坏成立**。
- 票面 `:6` 的 `grep -rln '模型装载器'`"只命中 109 自己" => 今天命中 **2 枚**
  （`issues/109-…md`、`issues/357-…md`）=> **腐坏成立**（只可当立案时刻读数引）。
- 票面 `:11` 说票 342/348 是 `-done` => **不成立**：两枚文件名**无** `-done` 后缀，`Status:` 逐字均
  `**已立，未派**`（342 `:3`、348 `:3`）=> **腐坏成立**。

### 3.3 普查腿两条关键读数复跑（本机今晚重跑，非引用）
`docs/evidence/s1/357-a1-census.md` 的 D1 = **主体 + ADS => `与已验签清单一致` rc=0（静默通过）**。
本腿用**今晚新编的 `wisp.exe`**（`go build ./cmd/wisp` rc=0，落在 `357r1/bin/`）在**新台子**上重跑，逐字：

```
===== dataG : wisp models verify fixture-a1 =====
  | wisp models verify: fixture-a1 与已验签清单一致（store=D:\tmp\wisp\20261021\357r1\dataG\models，0.00s）
rc=0

===== dataD1 : wisp models verify fixture-a1 =====
  | wisp models verify: fixture-a1 与已验签清单一致（store=D:\tmp\wisp\20261021\357r1\dataD1\models，0.00s）
rc=0
```

=> **D1 静默通过复量成立**（`dataD1` 里 `dir /r` 逐字有 `19 model.onnx:attacker:$DATA`，
而 `os.ReadDir` 只见 1 枚 `model.onnx`）。两条读数同时给正例 G：`dataG` 与 `dataD1`
**输出逐字相同**（差别只有 `store=` 路径段），这就是〔丁〕第 3 点要保的那条基线。

### 3.4 落点已定位（下一腿不用再找）
- 那句话的**出口**：`cmd/wisp/models.go:253-255`（`git show HEAD:cmd/wisp/models.go` 逐字
  `fmt.Fprintf(io_.out(), "wisp models verify: %s 与已验签清单一致（store=%s，%.2fs）\n", ...)` 后 `return 0`）。
  这就是 rc=0 那条路径，**新句子只能加在这里或它上游**。
- **枚举要打的目录**：`internal/models/downloader.go:298-301` 逐字
  `dir := filepath.Join(m.opts.DataDir, id)` / `if override, ok := m.opts.LocalOverride[id]; ok { dir = override }`
  => **有 override 时必须枚举 override 那个目录**（与 `:302 handOffCheck(entry, dir)` 同一枚 `dir`），
  NOT `store.root`。这条正是味一留下的坑，落地腿**必须**沿 `VerifyInstalled` 的 `dir` 走。
- 现状：`VerifyDir`（`downloader.go:628`）只做"哈希面 + `verifyNothingUnnamed`"两跳，
  **没有**任何流枚举 => 派单说的"今天静默"在码面成立。

## 4 我为什么不自作主张选 sec 2 那条路（按 `AGENTS.md` 铁律第 0 条第 4 项的边界）

铁律第 4 条要 agent **自己拍技术选型**，但派单 `:22` 逐字**"技术选型（我拍，NOT 你自己换）"**，
且 `:70` 把"需 syscall 手写"单列成**停手上报**项 => 这里是**派单明确收回**的那一格，
NOT "铁律与派单冲突、按铁律自己拍"的那一格。**故我停手，NOT 换形。**

我倾向 sec 2 那条路（`GetFileInformationByHandleEx` + 24 字节布局），理由按铁律第 2/3 条：
(1) **最有价值**：它把 D1 的"静默通过"变可观测，与派单 `:27` 拍〔丁〕的理由**同一句**；
(2) **复杂买得到真价值**：复杂度只多一枚 24 字节布局 + 一个重试循环，**不新增依赖、不手写 syscall**；
(3) 它比派单原选型**更便宜**（派单那枚 API 在 x/sys 里根本不存在）。
**但这属于改选型 => 你拍。**

## 5 门禁（停手前现跑，四数分开报，NOT 相减）

| 把 | rc | 末行逐字 |
|---|---|---|
| `GOFLAGS= go build ./...` | **0** | （无输出） |
| `go vet ./internal/models/` | **0** | （无输出） |
| `go test -count=1 ./internal/models/` | **0** | `ok  github.com/CarlosShao/wisp/internal/models  13.419s` |

（`./...` 与 `internal/models` 两包在**零产码改动**的工作树上跑；`cmd/wisp` 整包因本腿未落地故未跑，
留待落地腿按派单 `:67` 抬 `-timeout 30m` 现跑。）

## 6 写面与纪律格

- 仓内 `.go` 改动 **0 枚**、`internal/**`/`cmd/**`/`scripts/**/.github/**` 改动 **0 字节**
  （`git status --porcelain` 现读，写面内零行）。本件是**唯一**仓内落点（`docs/evidence/s1/357-r1-impl.md`）。
- 仓外只建不删：`D:\tmp\wisp\20261021\357r1\**`（`bin/` `logs/` `probe/` `dataG/` `dataD1/`
  `many/` `wide/` `longname/`）。未产 `.out`。
- 未 `add -A`/`.`；未 amend/reset/rebase/stash/checkout ./clean；**未 commit、未 push**（本腿零提交）。
- 未翻任何框、未改票面 `Status:` / 名字 / `-done`、未改台账、未碰票 109 任何原句。

## 7 请裁一件（一件，NOT 三件）

**[Q] 派单 `:23-26` 的技术选型前提不成立**（`x/sys v0.48.0` 无 `FindFirstStreamW`/`FindNextStreamW` 封装）。
请二选一：

- **(甲)** 批准改用 **`GetFileInformationByHandleEx` + `windows.FileStreamInfo`**（sec 2 已量通：
  ZERO 新依赖、零手写 syscall、G/D1 两发读数分得开、200 条流带重试全出）。=> 我按〔丁〕原判据落地。
- **(乙)** 批准**手写** `FindFirstStreamW`（我手映 `WIN32_FIND_STREAM_DATA` + `NewLazySystemDLL`）。
  => 也能落地，但比 (甲) 多一片手写 syscall 面，按铁律第 3 条我**不推荐**。

裁完我按原派单一次做完 `AC#3`（落地 + 今天响的用例 + 反向判据转红 + 三读数 + 正例 G）＋ `AC#4`。
