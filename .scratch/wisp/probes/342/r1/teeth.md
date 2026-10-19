# 342-r1 去牙反控两发（票 342 AC#3 ⇄ 票 348 AC#3）

★两发都只打**仓外导出树的物理字节**（`treeMut`／`treeMut2`＝`treeLand` 的整拷），⛔ 改跟踪文件、⛔ `-overlay`、★两发**串行**（跨 shell 的 trap/finally 有竞态＝公共块口径）。落点＝**内容锚**（函数体内那枚 token）、⛔ 行号。每发四步齐＝跑前验 pristine 哈希 → 落刀前验旧锚命中＝期望 → 落刀后验旧锚 0 且新锚≥1 → 跑 → 还原 → 报还原哈希并复验 pristine。

## 成对的落点⛔ 相同（这是 AC#3 的硬条件）

| 发 | 文件 | 旧锚（落刀前命中＝期望） | 新锚 | 落刀后 sha256（前 16） |
|---|---|---|---|---|
| 1＝摘掉新加的祖先守卫 | `internal/winsec/pinnedleaf_342_unix.go` | `pinPlacement(dir)` 命中 **1**＝期望 | `os.Open(dir)` 命中 1、旧锚 **0** | `28bbd47e1a9a9ce4` |
| 2＝守卫换成"逢名就拒"反形 | 同文件、**另一枚函数体首行** | `\tpath := filepath.Join(p.name, base)` 命中 **1**＝期望 | `MUTATION: refuse every name` 命中 1 | `427f7cf914b1d507` |

pristine（两发同一枚基线）＝`189e55e8613249db`；★还原后哈希＝发 1 `189e55e8613249db`（pristine_match=YES）、发 2 `189e55e8613249db`（YES）、FINAL `189e55e8613249db`。⇒ 两发落点不同、还原都逐字回到 pristine。

## 发 1 ＝ 摘新加的祖先守卫 ⇒ 常驻钉**必回落到红**（仪器敏感正控）

落刀＝把 `pinnedParent` 里那枚 `pinPlacement(dir)` 换成 `os.Open(dir)`＝父目录改按名字打开 ⇒ 换手的祖先照旧被跟随、复验形同虚设。读数（`-run` 四枚：票 342 钉＋票 348 确定性钉＋票 348 竞态钉＋覆盖写常驻钉；`RUN=10 PASS_TOP=2 FAIL_TOP=2 SKIP=0`）：
- ★`--- FAIL: TestAC1POSIXPrivateFileThroughASwappedAncestorNeverTruncatesAForeignFile342r1`，红句逐字（新句，＝第 16 笔改过文本的那枚）：`AC#1 RED: 0 attempts found the foreign witness's body replaced and 9 found it gone, so this route's act reached THROUGH an ancestor swapped into a link after resolveString had answered`。
- ★`--- FAIL: TestAC2POSIXPrivateFileRaceThroughASwappedAncestorNeverDeletesAForeignFile342r1`，红句逐字：`AC#2 RED: 5 of 1500 attempts found the foreign witness GONE, so the refusal-path unlink in privateFile re-resolved the spelling through an ancestor swapped into a link and removed a name in somebody else's tree`。
- ★liveness 同时在场：`AC#1 RED instrument: the racer proved 159 link-state cycles and the parent pin refused 0 attempts BEFORE any open (walk-after-open answers: 27, …)`＝守卫被摘⇒父钉那支答 0 枚＝仪器⛔ 静默、它照红。
- ★两枚既有牙⛔ 被摘：`TestAC2POSIXRefusalPath…`（确定性那枚）在这发里 **PASS**＝合理（它跑的路径由 seal 侧那支先答、⛔ 依赖父钉），所以本发⛔ "全盘红"、只把新加那两支打回红＝尺指向的就是新加的守卫。⇒ 这一行的意义＝★两枚新钉都⛔ 是"恒绿仪器"。

## 发 2 ＝ 反形"逢名就拒" ⇒ 覆盖写那四枚调用者该红

落刀＝`openPinnedLeaf` 函数体首行之后插一段 `if true { return nil, &os.PathError{Op:"open", Path: path, Err: syscall.EACCES} }`＝每一枚名字都被拒（＝"已存在即拒"那一族的极端形）。
第一版这发⛔ 编译＝我在 `return leafPin{dir: pin, name: dir}, nil` 那行落刀、把 `pin` 变成"声明而未用"⇒ `RUN/PASS/FAIL/SKIP` 全 0＝假读数（★本腿自抓、具名换落点重打、⛔ 计入判据；旧那发的日志留在 `logs/B2-teeth-driver.txt` 里原样、⛔ 删）。重打带 `go vet` 预验 rc=0。读数（`RUN=7 PASS_TOP=0 FAIL_TOP=1 SKIP=0`＝顶层 1 枚红、7 枚 RUN 含子测试）：
- ★`--- FAIL: TestAC2OverwritePromiseHoldsForEveryProductionCaller322r1` 及其四枚调用者子测试**逐枚红**，红句逐字（各一枚，`openleaf_overwrite_322_test.go:70`）：
  - `AC#2 RED: PrivateFile refused an existing name (open …/blob-322r1.bin: permission denied), which breaks storing the same id twice replaces the blob in place`（`internal/secret/store.go:79 Store.Set`）
  - `…which breaks this call is the repair path for a wide backup left by an older build`（`internal/config/migrate.go:93`）
  - `…which breaks a pre-rename temp left by a crash is overwritten, not refused`（`internal/secret/migrate.go:189`）
  - `…which breaks a leftover .retry temp is overwritten (D31 last-writer-wins)`（`internal/agent/spill.go:166`）
  - 另两枚子测试也红（`the entry stays private after the rewrite`＝`:98`／`the exclusive branch…`＝`openleaf_overwrite_322_test.go` 内那支）＝★同一枚反形连"独占仍答 `fs.ErrExist`"那条一起打红⇒尺对两半承诺都敏感。
⇒ ★结论＝**这把尺对"覆盖写被改成拒绝"照红**，而我的落地字节⛔ 让它红（成对四数里它逐名 PASS，见 `readings.md`）＝必答 (iv) 的那句"我这形⛔ 动覆盖写语义"有了仪器证、⛔ 是自我声明。

## 与票 342/348 AC#3 的对照

- ⛔ `-overlay`、★只打仓外树物理字节＝两处都是对 `treeMut`／`treeMut2` 的 `perl -0pi` 落刀；
- ★每发落刀前验旧锚命中＝期望（1）、落刀后旧锚 0 且新锚≥1＝两发都齐（表内）；
- ★跑前跑后各验一次 pristine 哈希、报还原哈希＝`189e55e8613249db` 三次复验一致；
- ★两发串行、⛔ 并发 shell；成对两发落点⛔ 相同（不同函数体、不同 token）；
- ★两发红句逐字进本件；两枚"该红的都红了、该绿的（发 1 里那枚确定性钉）没被误红"＝两把形状都放行＝⛔。
