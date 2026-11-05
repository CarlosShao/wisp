# internal/config 去红腿 186 — `TestHandDeletedSchemaVersionOnTheWriteBranchIsForcedBackAndUnclaimed` on POSIX：**FLAKY（非确定性红）**，一字节未动，停手交件

- 起手锚＝`248041c2e9b965efcff1e5d3d2b8b36a10cec5ff`（`248041c2e`）。判据件＝`internal/config/writeguard_unnailed_341_test.go`（票 341 落地腿 `6948eb087` 交付物），本腿**一字未改**：HEAD blob 与盘上 blob 同为 `fb305668674cdba9082b4ec3f519b6176c19320b`，`git status --porcelain -- internal/config` **空**。
- 派单点名的两种情形，本腿落在**第一种**：**FLAKY ⇒ ⛔ 不"修"，只记读数、停手**。派单里的第二种（"确定性 red ⇒ 把 fixture 做成平台中立、Windows 拼法逐字相同"）**不适用**：红与路径拼写/分隔符**无关**（Windows 与 Linux 跑的是**同一串字节**、同一份 `matrixBase`，见 §6）。
- 写面＝本件（`docs/evidence/s1/config-handed-red-186.md`）。`internal/config/**` 允许写但**一个字节都没写**。

## 0. 结论先行（三行）

1. **VERDICT＝FLAKY，不是确定性 POSIX 红。** 同一枚二进制、同一台 Linux 容器、同一份源码：`-count=50` 两发＝**PASS 19 / FAIL 31** 与 **PASS 13 / FAIL 37**（合 100 发 **68 FAIL / 32 PASS**，失败率 62%–74%），失败**100% 落在同一条行**`:240`。
2. **不是产品缺陷、不是守卫失效、不是 fixture 拼写问题。** 该用例除 `:240` 之外的每一条断言（`:231` 版本号被强制回写、`:236` 字节变了、`:253` 备份逐字、`:259`/`:265`/`:269` 日志三句、`text` 内容）在 100 发里**零次**红；Windows 宿主 `-count=50`＝**50/50 PASS**。
3. **病因＝文件时间戳的量子。** 这台 Linux 的文件 mtime 只按 **≈4.00047 ms** 的台阶走（＝1/250 s 节拍，实测见 §3）；用例的两次写（`mutate()` 的直接写 → `SaveFile` 的临时文件+rename）之间只隔**2–4 ms**，于是**两次写经常落进同一个 4 ms 桶**，`statMod(...).Equal(mtimeBefore)` 就红。`:240` 想证明的"这是一次真写"由 `:236` 的**字节**比较已经在证，mtime 这一枚是**用一台比被测量事件更粗的钟去测它**。

## 1. 起手读数（容器内逐字，可复跑）

容器＝`golang:1.27`（真实 Linux，`go1.27.1 linux/amd64`，`uname -r`=`6.6.114.1-microsoft-standard-WSL2`），`CGO_ENABLED=1`，`TMPDIR=/tmp`（`/tmp` 落在容器 overlay 根上），`GOPROXY=off`＋只读挂载宿主 `GOMODCACHE`（`internal/config` 不依赖 k2-fsa，容器内无需 sherpa）：

```
docker run --rm \
  -v 'D:\work\workspace\projects plans\Wisp:/src' \
  -v 'D:\work\base\gopath\pkg\mod:/gomod:ro' \
  -e GOMODCACHE=/gomod -e GOFLAGS=-mod=mod -e GOPROXY=off -e TMPDIR=/tmp \
  -w /src golang:1.27 sh -c 'CGO_ENABLED=1 go test ./internal/config -count=1 \
     -run "TestHandDeletedSchemaVersionOnTheWriteBranchIsForcedBackAndUnclaimed" -v'
```

```
=== RUN   TestHandDeletedSchemaVersionOnTheWriteBranchIsForcedBackAndUnclaimed
    writeguard_unnailed_341_test.go:240: mtime did not move on the write branch
--- FAIL: TestHandDeletedSchemaVersionOnTheWriteBranchIsForcedBackAndUnclaimed (0.02s)
FAIL
FAIL	github.com/CarlosShao/wisp/internal/config	0.028s
```

红句出处＝`writeguard_unnailed_341_test.go:239-241`：

```go
if statMod(t, path).Equal(mtimeBefore) {
    t.Error("mtime did not move on the write branch")
}
```

（`statMod`＝`os.Stat(path).ModTime()`，`writeguard_226_test.go:449-456`；`mtimeBefore` 取自 `:220`。）

## 2. 复跑矩阵（每发都是独立事实，不相减）

| 尺 | Linux 容器 | Windows 宿主 |
|---|---|---|
| 单发 `-count=1` | **FAIL**（§1 逐字） | **PASS** `(0.04s)` |
| `-count=5` | **FAIL, PASS, PASS, PASS**（1/5） | —— |
| `-count=50`（第 1 发） | **PASS=19 FAIL=31**（62% 红） | **PASS=50 FAIL=0** |
| `-count=50`（第 2 发） | **PASS=13 FAIL=37**（74% 红） | —— |
| 整包 `go test ./internal/config -count=1` ×3 | `ok 0.912s` / `ok 1.002s` / **FAIL 1.221s**（红名**只有**本枚） | `ok 1.969s`（rc=0） |

**失败行分布（两发 `-count=50` 合并，`grep -E "_test.go:[0-9]+:" | sort | uniq -c`）**：

```
     37     writeguard_unnailed_341_test.go:240: mtime did not move on the write branch   (第 2 发)
     31     writeguard_unnailed_341_test.go:240: mtime did not move on the write branch   (第 1 发)
```

⇒ 100 发里**没有第二条红行**。这解释了转盘上的 `FAIL/FAIL/PASS`：它本来就是一个 ≈2/3 红的硬币，三发里出一次 PASS 完全在射程内。

## 3. 机理读数：这台机器的 mtime 量子＝4.00047 ms

探针（自包含 `probe.go`，在容器内对 `/tmp` 跑；源码在宿主 `%TEMP%\wisp-probe-186\probe.go`，不进仓）：

**A. 模仿 `atomicWrite` 的"直接写 → 睡 d → 临时文件+rename"，看两次 mtime 是否相同（每格 100 次）**

| 两次写之间的睡眠 | 同一 mtime 的比例 |
|---|---|
| 0 | **97/100** |
| 250 µs | 65/100 |
| 500 µs | 66/100 |
| 1 ms | 67/100 |
| 2 ms | **40/100** |
| 5 ms | **0/100** |
| 10 ms | 0/100 |

（250 µs–1 ms 那三格不是严格单调：剩余变量是"落在 4 ms 桶里哪个相位"，不是睡眠长度本身。）

**B/C. 量子本身**

```
B tight back-to-back writes: same-mtime=393/400
B smallest non-zero delta=4.00047ms  median=4.00047ms  largest=4.000471ms
C 500 tight writes -> 16 distinct mtime values, smallest observed gap between distinct values = 4000470 ns
```

⇒ **mtime 只有 4.00047 ms 的台阶**（＝1/250 s，即该内核 HZ=250 的粗粒度时间）。连续写 500 个文件只产生 **16 个不同 mtime 值**，最小差 **4000470 ns**。

**为什么用例正好落在边界上**：`mutate()` 是 `os.WriteFile`（`manager_test.go:33-38`，先睡 10 ms 再写）；此后 `SaveFile`（`loader.go:289-303`）先 `backupOldSchema` 写 `.bak-1`、再 `atomicWrite`（`parse.go:196-230`）建临时文件→写→rename。两次"真写"之间夹着 `readText`、`os.Stat`、备份读+写、Marshal、`winsec.SealFile`、`CreateTemp`——实测这段落 2–4 ms（§3A 的 2 ms 格＝40% 同桶，5 ms 格＝0%），**正好横跨 4 ms 量子**，所以它是一个**概率**，不是一个结果。

⚠ 标为**推断、未逐条实测**的部分：4 ms/16 桶是**本机内核时钟**的事实；换一台 HZ=1000 的机器量子会变 1 ms，同一断言的红率会下降但**不会归零**（只要两次写的间隔可以小于量子）。本腿只对"这台容器"给出上的读数。

## 4. 观测者效应（把间隔拖长 ⇒ 立刻全绿，反证机理）

用 `strace -f -ttt -y -e trace=write,rename,renameat` 包住同一枚测试二进制跑 8 发（`apt-get install -y strace` 在容器内）：

```
run1..run8 全部 PASS rc=0
```

strace 的逐 syscall 停顿把两次写拉开到 **17–25 ms**（同一份 `st.log` 里：`config.toml` 写 `1793870628.112776` → `.bak-1` 写 `.129205` → `renameat` `.137732`），远超 4 ms 量子 ⇒ **8/8 PASS**。**同一份源码、同一台机器，只是把事件放慢，红就消失**——这与"量子边界"模型逐字吻合，而与"守卫/合并逻辑错了"不相容。

## 5. 为什么这**不是**产品缺陷（三条独立读数）

1. `:236` `if text == before` **从未**红过：每一次跑（含 68 发红）**盘上字节都真的变了**、`schema_version = 2` 都真的被 `SaveFile` 盖回来了（`:231` 零红）——产品行为在每一发里都成立。
2. `.bak-1` 存在且逐字（`:249-255`）零红；日志三句（`:259` `wrote merged change…`、`:262` 不借"nothing-written"句、`:265` 不把 schema_version 记成自己写的、`:269` `wrote=[fs.allowed_dirs]`）零红。
3. Windows 宿主同一份字节 **50/50 PASS**：NTFS 的 100 ns 时间戳＋宿主更细的时钟让这个断言几乎恒真。**两边差的不是拼写，是钟。**

## 6. 派单里"平台中立 fixture"那条为什么不适用

派单预设的确定性红成因是 **Windows 拼写**（反斜杠/盘符/大小写）。本腿逐条对号：

- 红发生在一枚**纯时间戳比较**：`statMod` 只用 `os.Stat`，不碰路径字符串、不碰 `filepath`/`risk.PathResolver`。
- 两边跑的是**同一份 `matrixBase`**、同一份 fixture 字节（判据件 blob 见抬头）；容器内的临时目录由 `sealableTempDir124`＋`TMPDIR=/tmp` 给出 `/tmp/Test.../001/config.toml`，Windows 给 `C:\Users\...\001\config.toml`——**断言里没有任何一处依赖那个拼法**。
- ⇒ 若按"确定性红"分支去改 fixture，改的是**风马牛不相及**的东西，且必然要动 `:240` 这枚断言（＝放松）。**本腿没动。**

## 7. 门禁读数（本腿自己复跑，逐条）

| 尺 | 读数 |
|---|---|
| `go test ./internal/config -count=1`（Windows 宿主，`unset GOTMPDIR TMPDIR TEMP TMP`，`PATH=$PWD/build:$PWD/third_party/sherpa-onnx:$PATH`） | **`ok  github.com/CarlosShao/wisp/internal/config  1.969s`**（rc=0） |
| `gofumpt.exe -l internal/config` | **空**（rc=0） |
| `go build ./cmd/... ./internal/... ./tools/...` | **rc=0** |
| `go build ./...` | **rc=1**，报错**全部**落在 `build\mut_rvw.go`（`undefined: voiceTtsFiles / residentBall / residentAudio / voiceBargeinMonitor`…）。该件＝`git check-ignore` 命中 **`.gitignore:14  build/`** 的**未跟踪**本机变异台架，mtime **2026-11-04 10:26:45**（本腿开工于 **11-05 17:02**）⇒ **别腿遗产、⛔ 本腿产物**；剔除该包后 `go build` **rc=0**。同一枚噪声在 `docs/evidence/s1/C26-dered-185.md` §6「噪声 A」已按同口径记过 |
| `sh scripts/d22scan.sh` | **rc=0 clean**：`d22scan: clean - no D22 ban violations; live scope work: bans #1-5 internal/=362, bans #1-5 cmd/=57, ban #6 frontend/=116, ban #7 internal/tools/=48, ban #8 design/=39, ban #8 frontend/=116, ban #8 internal/=898, ban #8 cmd/=208 …`；`runtests.sh: OK - packages=[./...] top-level: PASS=35 FAIL=0 SKIP=0` |

## 8. 边界 / 未动清单

- `git status --porcelain -- internal/config` ＝ **空**；判据件 HEAD blob ＝ 盘上 blob ＝ `fb305668674cdba9082b4ec3f519b6176c19320b`。
- ⛔ 未动：任何 `internal/config/**`、任何 `t.Skip`、任何断言强弱（`:240` 逐字仍在）、`frontend/**`、`design/**`、`docs/PLAN.md`、`docs/specs/**`、`tools/d22scan/**`、`thresholds.go`/golden/SLO。
- ⛔ 未做：`git add -A`／amend／reset／rebase／stash／push。本腿只提交本件（显式 pathspec）。
- 探针源码与日志**不进仓**（宿主 `%TEMP%\wisp-probe-186\`，容器内 `/tmp`）；仓内零临时件。

## 9. 给下一位的三条登记（本腿**未取**，只登记，⛔ 当成已做）

1. **（记录，不改码）把 `:240` 从"钟"改成"钟之外的事实"。** 该断言的**意图**（"这是一次真写"）已由 `:236` 的字节比较证明；若哪一腿要消这枚 flake，主流做法是**不拿 mtime 当"写过"的证据**，或把"mtime 必须前进"改成"mtime 不早于 `mtimeBefore` 且字节变了"这类**不依赖量子**的形状。**代价**：那是判据件的一次改切，需另开票、按"断言不得放松"的口径审（`mtime` 前进这条**确实**想证 `SaveFile` 真的落了盘，不是纯装饰）。**本腿按派单 ⛔ 自行取此路。**
2. **（记录，不改码）产品面同形假设。** `internal/config/manager.go:168` 的轮询判据是 `st.ModTime().Equal(m.seenMtime) && st.Size() == m.seenSize` ⇒ "未变"。在 4 ms 量子下，**同尺寸**的外部改动若与上一次观测的 stat 落进同一个 4 ms 桶，会被这枚判据当成"没变"而跳过一次 tick（下一次写会救回来，所以是"延迟一次"，不是"永久丢"）。**这是同一族假设的读数，不是本腿的裁决**；要不要开票由编排者定。
3. **（记录）转盘口径**：本枚在 POSIX 上"三发两红一绿"是**预期**，不是回潮。谁再看到 `:240` 单行红，先按本件 §3 量一次 `mtime` 量子，再决定是不是新病。

## 10. 一句话终态 ＋ 撤销口令

**`TestHandDeletedSchemaVersionOnTheWriteBranchIsForcedBackAndUnclaimed` 在 POSIX 上＝FLAKY（本机容器 100 发 68 FAIL / 32 PASS，失败 100% 同一条 `:240`，病因＝4.00047 ms 的 mtime 量子 vs 2–4 ms 的两次写间隔；Windows 50/50 绿），不是确定性红、不是守卫缺陷、与路径拼写无关；本腿一字节未改、断言零放松，按派单停手交件。**

撤销口令（本腿只有一枚文件，一句话复原）：

```
git rm docs/evidence/s1/config-handed-red-186.md   # 若尚未提交
git revert <本件的提交 sha> --no-edit               # 若已提交（只含本件，零代码字节）
```
