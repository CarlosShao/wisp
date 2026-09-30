# 238-c1 — 只读普查腿：首启数据根的"谁造的"两形＋常驻尺＋POSIX 半边

- 腿：`238-c1`（**只读普查**，不产码、不跑测试、不占桌面）
- 起手锚点：`git log -1` 现读＝`bf26dfc4d027f85067549018e0d15a89aa4900e2`（Wed Sep 30 11:37:56 2026 +0800）
- 落笔钟：`date` 现读＝`2026-09-30 11:44:47 +0800`（§0 复跑时刻）／末次更正见各行内联时间
- 工单面：`.scratch/wisp/issues/238-first-run-on-a-fresh-machine-creates-the-whole-data-root-with-no-seal-to-inherit.md`
- 本节任何"今天红/绿"的判断**一律来自读断言**（本轮硬红线禁止 `go test`），逐条已具名标注〔未跑包，此条来自读断言〕。

---

## 0. 现量口径（我跑了哪些尺）

⛔ 本轮硬红线：零 `go test`、零容器、零产码。下面每一把尺都是 `git grep` / `sed -n` / `Read`；
`rc` 一律取 **git grep 自己**的退码（不接 `| cat`，接了读到的是 `cat` 的 0——本轮一条这类坑已实测：
`git grep ... | cat; echo $?` 恒等于 0，会"证明"零命中那把尺跑成功）。

### 0.1 票面"现量"三行复跑

**行 1（票面 `:11`）——复跑＝成立，但其中半句按另一种读法是假的，具名更正。**

| 尺（逐字） | 本机读数 |
|---|---|
| `git grep -nE '"github.com/CarlosShao/wisp/internal/winsec"' -- internal/proc internal/buildinfo` | **零命中，rc=1** ⇒ 两包的 winsec **导入语句数＝0**：成立 |
| `git grep -n "internal/winsec" -- internal/proc` | 5 行命中，**全在注释里**（`envfork.go:113/116/148`＋`envfork_test.go:115/147`）⇒ "文字提到、代码不依赖" |
| `git grep -nE "winsec\|Mkdir" -- internal/buildinfo` | **零命中，rc=1** |
| `git grep -n "internal/buildinfo" -- internal/proc` | **5 行命中（rc=0）**，其中产码 2 枚：`boot_windows.go:17`、`envfork.go:10` |

⚠ **具名更正**：票面那一格写「`internal/proc` 对 `internal/winsec` 导入数＝0；**`internal/buildinfo` 同样**」。
这句有两种读法，**只有一种成立**：
- 读法甲＝"`internal/buildinfo` 对 winsec 的导入数也是 0"⇒ **为真**（上面第二／三把尺，连 `winsec` 这个词都在该包零命中）。
- 读法乙＝"`internal/proc` 对 `internal/buildinfo` 的导入数也是 0"⇒ **为假**：`internal/proc` 的两枚产码文件都导入 `internal/buildinfo`。
出处＝`132-c3` §5.3 的原尺是 `grep -rn "MkdirAll\|os.Mkdir" ... internal/proc/ internal/buildinfo/`（比的是**建目录动作**，不是导入数），
票面把它压缩成"同样"两个字时把比较对象丢了。**落地腿按甲读，别照乙写。**

**行 2（票面 `:12`）——复跑＝成立，名册逐枚。** `PrivateDirAll` 生产点＝**恰好 5 枚**（尺：逐文件 `grep -Hn "PrivateDirAll("` over `git ls-files internal cmd | grep -v _test.go`）：

| 落点 | 它封的是 |
|---|---|
| `internal/secret/store.go:49` | `<data>\secrets` |
| `internal/agent/spill.go:111` | spill 目录（`s.dir`） |
| `internal/memory/open.go:176` | **`<data>` 本身**（`abs`＝数据根，注释 `:173`–`:175` 逐字「The data root is sealed before anything is written into it, and sealed **with inheritance**, because wisp.db-wal / wisp.db-shm are created by SQLite and nothing in this repository gets to seal them afterwards」） |
| `internal/memory/open.go:188` | `<data>\artifacts` |
| `internal/memory/open.go:503` | `<data>\backup` |

⇒ **正控形状确实存在，但"数据根自己"只有 `open.go:176` 一枚**；`memory.Open` 的生产调用者**只有 2 枚**
（`cmd/wisp/run.go:384`、`cmd/wisp/providers.go:180`；尺：同法 `grep -Hn "memory\.Open("`）。
⇒ 本票的真正射程比票面标题更宽，见 §0.4。

**行 3（票面 `:13`）——复跑＝成立，行号漂一处，具名更正。**
`BuildDiagnosticsBundle` 生产调用者＝**0 枚**（尺：`git grep -n "BuildDiagnosticsBundle" -- '*.go'` ⇒ 命中只有
定义 `internal/observe/diagnostics.go:61`（注释）／`:62`（func）＋同包 `diagnostics_test.go:35/:125/:140` 三枚测试）。
⚠ **"把日志打进 zip"那段不在 `:73`**：现读在 `internal/observe/diagnostics.go:99`–`:129`
（`:100 entries, lerr := os.ReadDir(o.LogDir)`、`:122 writeZipEntry(zw, "logs/"+name, ...)`）；
`:73` 现读是 `os.Create(o.OutPath)` 那一行。

### 0.2 票面正文引用的两行——逐字命中，无漂

- `cmd/wisp/resident_windows.go:57`：`sink, sinkErr := installLogSink(rt.Layout.DataDir)`（`:53`–`:56` 的时序注释逐字含「proc.Boot is the only thing that ran before this, and it seals nothing - internal/proc has zero winsec imports」）。
- `internal/observe/logging.go:72`：`if err := os.MkdirAll(cfg.Dir, 0o755); err != nil {`。
- 落盘第二跳：`internal/observe/logging.go:244`：`f, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)`（滚动日志文件自己**既不经 `MkdirAll` 也不经任何封条**）。

### 0.3 票面没写、但直接决定 AC#1 的一把尺：`SealDir` 的名册

尺：`git grep -n "SealDir" -- internal cmd tools ':!*/*_test.go'`
读数＝**6 行，全部在 `internal/winsec` 包内，且只有 `winsec.go:222` 是函数本身**：

```
internal/winsec/winsec.go:174  注释
internal/winsec/winsec.go:190  注释
internal/winsec/winsec.go:220  注释（doc）
internal/winsec/winsec.go:222  func SealDir(path string) error {     ← 定义
internal/winsec/winsec_windows.go:52  注释
internal/winsec/winsec_windows.go:313 注释
```

⇒ 票 132 标题那句「`SealDir` 有 16 处测试引用、**0 个生产调用者**」**到今天仍未变**（测试侧命中：`git grep -n "SealDir(" -- 'internal/winsec/*_test.go'` 逐枚可数，本腿数到 16 枚调用＋2 枚 POSIX 用例名注释）。
⇒ **`winsec` 包里"封一棵树"这个动作，产码侧唯一的入口是 `PrivateDirAll`，而它只封"自己创建的那些层＋请求点名的那一层"**（`winsec.go:152`–`:157` 逐字：「Ancestors above path are never touched」）。

### 0.4 本腿新量到的四条（票面没有，是 AC#1 的输入）

1. **造根的腿不止常驻腿一枚。** 尺：逐文件 `grep -Hn "installLogSink\|observe.InitLog\|os.MkdirAll"`（排 `_test.go`）。读数：
   - 经 `installLogSink`→`observe.InitLog`→`logging.go:72 MkdirAll` 建树的产码腿＝**4 枚**：
     `cmd/wisp/resident_windows.go:57`、`cmd/wisp/run.go:195`、`cmd/wisp/models.go:284`、`cmd/wisp/secret.go:226`；
   - **绕过 `installLogSink` 直接 `observe.InitLog`** 的腿＝1 枚：`cmd/wisp/slo_windows.go:272`（`slo_other.go:14` 说明这条腿不进 Boot）；
   - 不经任何 sink 就建根的腿＝1 枚：`cmd/wisp/doctor.go:301` `probeWritable` 的 `os.MkdirAll(dir, 0o755)`（调用点 `doctor.go:108`）。
   ⇒ **六枚腿里只有 `run` 那一枚会在后面某一步把根收窄**（`run.go:384 memory.Open` → `open.go:176`）。
   ⚠ 这条把票面标题的"先跑常驻腿"改写得更准确：**"没有一条腿把'根'当成自己的责任"**才是形状；常驻腿只是唯一一条**今天连 `memory.Open` 都不会路过**的腿。
2. **`internal/observe` 对 winsec 的导入数＝0**（尺：`git grep -n "internal/winsec" -- internal/observe` ⇒ 2 行命中，全在 `logging.go:374/:386` 注释里）。
   而 `:386`–`:387` 禁的是**反方向**那枚边（逐字「internal/winsec **must not import this package**（the graph runs observe -> secret -> winsec, so the edge would close a cycle）」）⇒ **"observe 能不能导入 winsec"这句话票面、注释、仓里都没裁过**，见 §4 第 1 条。
3. **依赖图现读**（尺：逐文件 `grep -n "CarlosShao/wisp/internal"`，排 `_test.go`）：
   `internal/winsec` 产码侧**零内部导入**（是一枚叶子）；`internal/proc` 产码导入＝`buildinfo`＋`observe`（`boot_windows.go:17/18`、`envfork.go:10`、`externalsampler_windows.go:10`、`treemetrics_windows.go:12`）；
   `internal/observe/redact.go:10` 导入 `secret`；`internal/secret/store.go:10`＋`migrate.go:14` 导入 `winsec`。
   ⇒ **`internal/proc` 今天就已经通过 `observe -> secret -> winsec` 传递依赖 `winsec`**；给它加一条**直接**边不闭环、不新增可执行依赖，只是把一条已经在图里的边走到明面上（这一条票面 `:17` 点名要量，答案在 §1）。
4. **静态门里没有任何一格与封条有关**：`tools/d22scan/main.go:141`–`:164` 现读 9 条规则（`secretNameRe`/`wallclockRe`/`unixTimeRe`/`timeoutWordRe`/`mirrorHashRe`/`mirrorWordRe`/`approvalPanelRe`/`artifactToolRe`/`assignKeyShapeRe`）＋`emojiRe`，**没有一条认得 `os.MkdirAll`/`Seal`/`Create`**（复认 `132-c3` §6.1 的同一条尺）。⇒ 票面"今天没有任何尺会因此变红"这句**成立**。

### 0.5 已存在的真机 ACL 归属读数（本腿没跑 `icacls`，见 §5；引仓内既有读数）

`docs/evidence/s1/128-ac1-consequences.md:143`–`:153`（票 128 AC#1，本机 `swq`）摆的就是本票 AC#2 那把尺要用的**逐枚 `(I)` 归属**：

| 对象 | 读数（原文摘要） | 性质 |
|---|---|---|
| `…\wisp-dev\secrets` | 六条 ACE **全都没有 `(I)`** | 显式私有＝winsec 自己写的（`PrivateDirAll` 站点） |
| `…\wisp-dev`（**数据根自己**） | 三条**全 `(I)(OI)(CI)(F)`** | 纯继承 ⇒ 建它的不是封条 |
| `…\wisp-dev\logs` | 同上，三条 `(I)` | 纯继承 |
| `…\wisp-dev\config.toml`／`logs\wisp-20260923-001.jsonl` | 三条 `(I)(F)` | 纯继承 |

同文件 `:160`–`:165` 给了"继承源最坏能多宽"的**只读**读数（`icacls C:\Users\Public`：`NT AUTHORITY\INTERACTIVE:(RX,WD,AD)` 等）
⇒ 这是本仓现存**唯一**一组成对的（宽／窄）真机读数，**票 238 的尺能不能咬，就拿它对表**。
⚠ 口径：那棵树是**票 128 修前的"回落到 CWD"**形态（根在启动目录下），不是 `%APPDATA%` 上的首启树；本腿只借它的**归属判据形状**，不把它当"首启现量"。本机 `%APPDATA%\wisp` 现读＝**存在且为空**（`ls` 零条目），`%APPDATA%\wisp-dev` 同样在名册里 ⇒ 这台机器上**没有**一棵真实的宽树可测。

### 0.6 三行读法（本票碰隐私边界，措辞不许重于证据）

1. **现象在哪出现**＝**盘上那棵数据根树的创建时序**（哪条腿先跑、哪枚 `MkdirAll` 建了根），不在任何对外接口上。
2. **有没有本机被入侵的证据**＝**没有**。本腿全程 `git grep`／`sed -n`／`Read`，零写入、零执行被测程序、零 `icacls`；本机 `%APPDATA%\wisp` 至今是空目录。
3. **最坏后果是什么形状**＝**首启由非封条腿建的那棵树的描述符完全由父目录决定**：父目录窄则今天无人能读（本机就是这种），父目录宽（共享卷／漫游档案／`C:\Users\Public` 那种 `(RX,WD,AD)` 形状）则日志与根跟着宽，而 `docs/evidence/s1/128-ac1-consequences.md:149` 证明同一条腿上 `secrets` 是显式窄的 ⇒ **差别的存在已被量到，差别的"宽度"没被量到**。⛔ 定性＝**未证实的暴露面差异**，不写"被窃听/已越权"。

---

## 1. 两形代价表（AC#1）

（待补：形①常驻腿自封／形②装配根先封根，各自的落点文件:函数、跨包依赖边、失败方向、最坏后果形状、blast radius。）

---

## 2. 一把不恒真的常驻尺（AC#2）

（待补：判据全文＋正控设计＋"今天绿着会被打红"的名册。）

---

## 3. POSIX 半边具名（AC#3）

（待补：chmod 单档不递归的三条走法与代价＋"本机验不了"的具名口径＋要跑哪一发在哪跑。）

---

## 4. 要编排者裁的（不许为空）

（待补。）

---

## 5. 没做完的（不许为空）

（待补。）
