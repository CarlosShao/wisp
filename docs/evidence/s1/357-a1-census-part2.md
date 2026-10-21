# 票 357 普查腿 `357-a1` 读数件 — part 2/2（本件续 `357-a1-census.md`；主件过帽故拆 `-partN`）

## 3 `AC#1` 三件（仓外真 Go 程序，真 NTFS）

造件（`D:\tmp\wisp\20261021\357a1\ntfs\`，卷 `D:`）：主体 `model.onnx` **12 B**；ADS `model.onnx:attacker` **19 B**；
硬链接 `hard.onnx` 指向同一 inode（`New-Item -ItemType HardLink`；`fsutil hardlink list` 回列两个名字）。
没有 mock：全部 `os.*` 与真 Win32 API 打在真盘上。程序 `ntfsprobe\main.go`，读数件 `logs/ac1-ntfs.txt`。

| 读数 | 主体 `model.onnx` | ADS `model.onnx:attacker` | 硬链接 `hard.onnx` |
|---|---|---|---|
| `os.ReadDir(dir)` | **列出**（`"model.onnx"`） | **列不出**（`contains "attacker"`=false、`"model.onnx:attacker"`=false） | **列出**（`"hard.onnx"`） |
| `filepath.WalkDir`（`verifyNothingUnnamed` 的仪器） | 列出 `type=----------` | **列不出** | 列出 `type=----------` |
| `os.ReadFile` | 12 B，sha `206402cab345…b0eb` | **读得到 19 B**，sha **`dc3e65008d74…f2f4`（与主体不同）** | 12 B，sha `206402cab345…b0eb`（**与主体同 sha**） |
| `os.Lstat` | `regular=true size=12 symlink=false` | `regular=true size=19 symlink=false`（**读作一枚普通 regular file**） | `regular=true size=12 symlink=false`（**分不出是链接**） |
| `os.Stat` | 同上 | `size=19 regular=true` | `size=12 regular=true` |
| `nlink`（`GetFileInformationByHandle`，**`fs.FileInfo` 不带**） | **2** | 2（同一 inode） | **2** |
| fileIndex（volSerial,hi:lo） | `1747844702, 196608:2976803` | 同 | **同**（`identical=true`） |
| 流枚举 `FindFirstStreamW` | `["::$DATA" ":attacker:$DATA"]` | —（同一 inode） | 同主体 |
| 硬链接名枚举 `FindFirstFileNameW` | `["\tmp\wisp\…\ntfs\hard.onnx" "\tmp\wisp\…\ntfs\model.onnx"]` | — | 同 |
| `dir /r`（cmd） | `12 model.onnx` ＋ `19 model.onnx:attacker:$DATA` | 同行显示 | `12 hard.onnx` ＋ `19 hard.onnx:attacker:$DATA` |

**3.1 产品级复量（不是理论推演）**：把这套形状塞进真 store，跑真 `wisp models verify`（读数件 `logs/ac1-product2.txt`）
| 发 | 形状 | `wisp models verify` | rc |
|---|---|---|---|
| D0 | 只有主体 | `与已验签清单一致` | 0 |
| **D1** | 主体 ＋ **ADS 19 B** | `与已验签清单一致`（点名文件的哈希面照过） | **0** |
| D2 | 主体 ＋ 目录内硬链接 `hard.onnx` | `model dir holds 1 entry the signed manifest does not name: hard.onnx`（票 121 `AC#4` 那条把它当**多余目录项**抓住） | 1 |
| **D3** | 主体 ＋ **目录外**硬链接指向它 | `与已验签清单一致`（目录项集合没变，改名面看不见） | **0** |
| D4 | 主体被换成指向**外部字节**的硬链接 | `model.onnx: sha256 mismatch (pins 206402cab345..., got 950edc0ad300...)`（哈希面抓住） | 1 |

=> 硬链接**不是一律绕过**：**落在被验目录里**的被票 121 的"未点名目录项"抓（D2）；**落在目录外**的它抓不到（D3），
但那种形状改的是**同一份已验字节**，哈希面本来也无从反对。**真正静默通过的是 ADS（D1）**。

**3.2 一条比 `109-v2` 更狠一格的读数（属本腿新增）**：`dir /r`、`os.ReadDir`、`nlink` 三把尺在
"目录外种一枚硬链接、再从那枚外部名字写入"这一手面前**全都无痕**（读数件 `logs/ac1-writethrough.txt`）：
```
step 1 verify => 与已验签清单一致                                        rc=0
step 2 种外部硬链接; nlink=2; os.ReadDir => model.onnx    (无第二名的痕迹)
step 3 从外部名字写 12 B(同长) => 树内 sha 变成 a33238fcb04b…
step 4 verify => model.onnx: sha256 mismatch (pins 206402cab345..., got a33238fcb04b...)   rc=1
```
哈希面抓得到（复量那一发必红）；目录项面抓不到（第 2、3 步全程只有 1 枚条目）。

**3.3 ADS 的耐久性（决定"该不该加流枚举"的成本读数）**，读数件 `logs/ac1-durability.txt`
| 动作 | 目标文件上的流 |
|---|---|
| `os.Rename` 同卷移动（产品 `install()` 主路 `:771`） | `["::$DATA" ":attacker:$DATA"]` => **流跟着走** |
| Go `os.Open`+`os.Create`+`io.Copy`（产品 `copyFile` `:781`，rename 失败兜底） | `["::$DATA"]` => **流丢掉** |
| Win32 `copy /b`（`CopyFileEx`，非本仓仪器，仅作对照） | 流**跟着走** |
- 故 ADS 能否进 store **取决于 `install()` 走 rename 还是 copyFile**：同卷 rename 会带进 store（D1 那条路可达）。
  `os.Remove` 删主体时流一起消失（V4）。
- **外部硬链接的先决条件已量**：安装目录 ACL 逐字 `Authenticated Users:(I)(M)`（`logs/ac1-writethrough.txt` 末段），
  `NumberOfLinks` 与 `FindFirstFileNameW` 给出两个名字 => "同卷 + 有写权"即可种链，D3 不是纸面形状。

**3.4 哪一层该拦（本腿判语）**
- **哈希面（再验一次哈希）**：对**硬链接写入**有效（3.2 第 4 步红），对 **ADS 原理性无效**（验的是主体那条流）。
- **目录项面（`WalkDir` + 未点名判定，票 121 已落地）**：对目录内硬链接有效（D2），对 **ADS 与目录外硬链接无效**（D1/D3）。
- **读取面**：今天**没有这一层**（读取者不存在，见 part1 sec 2.5）。
=> **判语：味二该拦的一层是"读取面契约"**（读的人只许开清单点名的那条流、那个名字），不是哈希面也不是目录项面。
理由：哈希面已证明对 ADS 无效；目录项面受 NTFS 语义限制（流不是目录项）；真要把 ADS 挡在门外的机制
（`FindFirstStreamW` 枚举 ＋ `nlink>1` 判定）是**新增的 Windows-only 平台件**，而它保护的那个读取者今天不存在。
=> **今天可达性 = 0**（不是"很小"）：`wisp models verify`/`ensure` 只重新哈希已存在的树，从不 open 交还的目录
（`bridge.go:44` 丢弃返回值；`cmd/` 与 `internal/` 对 store 字节的 open **0 处**；`internal/speech` 仍只 `doc.go`；
`internal/engines/` 不存在；`go list -deps ./cmd/wisp` 里 `sherpa_onnx` 只用到 `GetVersion`/`GetOnnxruntimeVersion`
两枚**版本号**符号）。

## 4 每把尺的射程声明 ＋ 正控

**4.1 行号尺**：`git show HEAD:<path> | grep -n` / `sed -n`。射程＝**HEAD blob 的 1-based 行号**，含注释与 `_test.go`
（本件只引非测试文件）；不含工作树 CRLF 差异、不含跨行子串（单行正则）。正控：同尺命中 `:293`
`func (m *Manager) VerifyInstalled`（part1 sec 1.5 表）；再打一枚必不命中的 `VerifyInstalledXXX` => 0 命中。

**4.2 "0 处读取者"那把尺（读 0，已换第二把独立形复读）**
- 尺 A（文本）：`git grep -n -E "(DataDir|store\.root|ms\.root|installDir)" HEAD -- cmd/ internal/ | grep -v _test.go | grep -E "Open|ReadFile|OpenFile"`
  => **0 命中**。射程：HEAD blob 的 `.go`、含注释（注释里出现这些词也会命中 => 只会**高估**、不会漏）、
  不含跨行子串、不含 `internal/tools` 与 `frontend/`。正控：同尺把关键词换成 `Ensure` => 命中 2 处
  （`bridge.go:44`、`models.go:305` 一带），尺能命中真调用。
- 尺 B（语义，独立仪器）：`go list -deps ./cmd/wisp | grep -c wisp/internal/speech` => **0**；
  `...grep -c wisp/internal/models` => **1**；`git ls-tree -r --name-only HEAD internal | grep -E "^internal/(engines|engine|loader|onnx)/"`
  => **0**；`internal/speech` 树 = 仅 `doc.go`。射程：**包粒度**，不含"某包里某函数没被调用"这种函数级判断
  （只能证明"装载器包不存在"，不能证明"没有别的读法"）。正控：同尺 `grep -c wisp/internal/models` => 1（数得到真包）。
- 两把尺给同一答案（0）=> "今天可达性 0"不是单尺结论。

**4.3 `cat-file -e` 尺**：射程＝**该 path 是否在 HEAD tree 里**；不解引用子模块、不判内容。正控见 part1 sec 1.1。

**4.4 `d22scan` 尺**（"禁码位"那条的仪器）：射程＝`U+1F000–1FAFF`、`U+2200–22FF`、`U+2600–27BF`、`U+2B00–2BFF`、
`U+FE0F`、`U+1F1E6–1F1FF`；**注释豁免、字符串不豁免**；**箭头 `U+2190–21FF` 与带圈数字 `U+2460–24FF` 刻意不扫**。
=> 本件与提交正文仍**人工**避开全部禁码位（不靠仪器兜底）。

**4.5 `wisp models verify` 尺**：射程＝**清单点名文件的 sha256+size ＋ 目录里未点名条目**；**不含**同一目录项上的第二数据流
（ADS）、**不含**目录之外的任何名字。正控：D4 => `rc=1` 点名 `sha256 mismatch (pins 206402cab345..., got 950edc0ad300...)`；
D2 => `rc=1` 点名 `hard.onnx`。两发都红 => 尺不是恒绿。

## 5 门禁（本次现跑；四数分开报，不相减）

| 把 | rc | 末行逐字 |
|---|---|---|
| `GOFLAGS= go build ./...` | **0** | （无输出） |
| `sh scripts/d22scan.sh` | **0** | `d22scan: clean - no D22 ban violations; live scope work: ... ban #8 internal/=608, ban #8 cmd/=144; ...` |

`go env` 现读：`GOOS=windows GOARCH=amd64 CGO_ENABLED=1`（未设 `CGO_ENABLED=0`，本仓在此配置下编得过）。

## 6 写面与纪律格

- **提交号：`a8f22b99`**（`a8f22b99eb964c9d82bf7825e90764348abc38e9`），主题逐字以 `票 357 普查腿 357-a1` 开头，
  正文经 `/tmp/msg-357a1.txt` ＋ `-F` 提交。`add`/`commit` 两处都用显式 pathspec（两枚 `docs/evidence/s1/357-a1-census*`）。
- 仓内**只落两枚读数件**（`357-a1-census.md` ＋ `357-a1-census-part2.md`，222 行合计，单件 101/121 行，均未过帽）；
  仓外 `D:\tmp\wisp\20261021\357a1\**`（`bin/` `probe/` `ntfsprobe/` `copyshape/` `manifest/` `data{A,B,C,D}/`
  `ntfs/` `durab/` `logs/` `shots.ps1`）只建不删。
- 逐枚 `.go`（本笔 `a8f22b99`）：**0 枚**（`git show --name-only --pretty= HEAD | grep -c "\.go$"` => 0）。
  越界枚数：**0**（`internal/`/`cmd/`/`scripts/`/`.github/`/`.scratch/wisp/`/`docs/reports/` 前缀逐条 grep 合计 0）。
  本笔文件全集只有那两枚 `docs/evidence/s1/` 件。
- 正控尺（同尺打一枚已知非零的样本）：`git show --name-only --pretty= a7dfca91` => 7 文件 / **7 `.go`**（尺活）；
  打本笔 => 2 文件 / 0 `.go`。
- 禁码位尺（本件＋提交正文）：射程见 sec 4.4，**注释不豁免于我的自扫**（我按派单逐字避开）。
  正控：同尺打一行含 U+00A7 与 U+2192 的探针串 => 命中 2 枚（尺活）；打本件两枚与提交正文 => 各 **0 行**。
- 临时件只建不删；未产 `.out`；未 `add -A`/`.`；未 amend/reset/rebase/stash/checkout ./clean；**未 push**
  （`git log origin/dev..HEAD` 现读仍含本笔，即它只在本机）；未翻任何框、未改票面/台账/`Status:`。

## 7 两件必须摆给编排者的

1. **味一的判据要改写一格**：今天不是"守卫验错了目录"，而是"**守卫验的是 override、CLI 说的是 store，而交还值没人接**"。
   part1 sec 2.5 的两把独立尺把"读者 open 错目录"这一后果压到**可达性 0**。若仍要按"守卫验的目录 == 读者 open 的目录"
   立判据，落点应是**给未来读取者一条强制 handing 契约**（拿 `Ensure` 的返回值，别自己拼 `store/_root/id`），
   而不是改 `VerifyDir`。
2. **味二的 A 档代价（`AC#2` 要用）**：流枚举（`FindFirstStreamW`，本机已跑通）＋`nlink` 判定
   （`GetFileInformationByHandle`，`fs.FileInfo` 不带）都是 **Windows-only 新平台件**，且它们保护的对象今天不存在
   => 按本仓既有口径"别给没人走的路加第二把锁"，A 档应明确写"今天零收益"。**比它们便宜得多的一格是 D1**：
   `wisp models verify` 报"与清单一致"时，连"目录里有第二数据流"这句话都不说。
