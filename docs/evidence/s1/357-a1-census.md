# 票 357 普查腿 `357-a1` 读数件（只读：`AC#0` ＋ `AC#1`）— part 1/2

起手锚 `c50218bf`，交件时 HEAD `e61861b4`。写面：本件 ＋ `357-a1-census-part2.md` 两枚 ＋ 仓外
`D:\tmp\wisp\20261021\357a1\**`。`.go` **0 枚**（仓内）、`internal/**`/`cmd/**`/`scripts/**`/`.github/**`/
`.scratch/wisp/**`/`docs/reports/**` **0 字节**、未 push、零翻框。
禁码位写作 sec N / "=> " / NOT / (1)(2)(3)。跑门禁前 `unset GOTMPDIR TMPDIR TEMP TMP`、`GOFLAGS=`；未设 `CGO_ENABLED=0`。

## 1 票面前提新鲜度（逐条；尺＝`git cat-file -e HEAD:<path>`，打 HEAD blob，不打工作树、不凭自述）

**1.1 票面/关联件 16 枚逐枚 `git cat-file -e HEAD:<path>` 全 rc=0**（票 109/119/121/342/348/357、
`docs/evidence/s1/109-v2-verdict.md`、`109-adversarial-acceptance.md`、`internal/models/{bridge,downloader}.go`、
`cmd/wisp/models.go`、`internal/winsec`、`docs/reports/pending-and-issues.md`、`models/manifest.json{,.minisig}`、
`scripts/d22scan.sh`）。**正控**：同尺打 `HEAD:docs/evidence/s1/357-a1-census.md`（交件前不存在）与
`HEAD:internal/engines/loader.go` => 各 rc=128，尺是活的。

**1.2 三枚产码 blob == 工作树**：`downloader.go` `b2d7c5df`、`bridge.go` `cfc0ab86`、`models.go` `4259a20f`
（`git rev-parse HEAD:f` == `git hash-object f`）=> 下文所有行号引的是**入库那一版**。

**1.3 `109-v2` 判语锚 `5ae2bae4` 成立**：`5ae2bae4:docs/evidence/s1/109-v2-verdict.md` blob `c24c38f0` == HEAD blob
`c24c38f0`（SAME）。

**1.4 票面三句引语复量（真值方向成立，三处落点要更正）**
- 票面 `:9` 说签名＝"`109-v2` 判语**第 2 条**"。现读：该句在 **sec 5(c)**（`109-v2-verdict.md:109`），sec 2 是"写面与纪律"。
  判语原文逐字："建议：并进「第一个模型装载器」那张票的 AC，NOT 另立一票、也 NOT 并本票。"
  => 归口方向与票面一致，**章节号是票面写错**（不改判据）。
- 票面 `:6` 的 `grep -rln '模型装载器' .scratch/wisp/issues/*.md` "只命中票 109 自己" => **今天不再成立**：票 357 已在盘上，
  现命中 **2 枚**（109 + 357）。"盘上从来没有那张票"这个**立案时刻的历史事实**仍成立（357 是 2026-10-21 才立的），
  但**今天用同一把尺重跑得到 2**，故那句只可当"立案时刻的读数"引。
- 票面 `:11` 说票 342/348 "同族拒绝路径面，`-done`" => **不成立**：两枚文件无 `-done` 后缀，`Status:` 逐字均为
  **"已立，未派"**（342 `:3`、348 `:3`）。票 119/121 的 `-done` 后缀成立（文件名带 `-done`）。

**1.5 行号腐烂（本仓已 ≥5 次，凡引 `file:line` 一律本次现跑复认）**
| 引 | 今天真身 | 判 |
|---|---|---|
| `downloader.go:293 VerifyInstalled` | `:293` 就是 `func (m *Manager) VerifyInstalled` | 未腐 |
| `downloader.go:298-300` 覆盖目录盖掉 store | `:298 dir := filepath.Join(m.opts.DataDir, id)`、`:299 if override, ok`、`:300 dir = override` | 未腐 |
| `downloader.go:628 VerifyDir` / `:663 verifyNothingUnnamed` | `:628` / `:663` 逐字命中 | 未腐 |
| `downloader.go:750` `sha256 mismatch (pins ..., got ...)` | `:750` 逐字命中 | 未腐 |
| `bridge.go:58 VerifyInstalled`（唯一非测试调用） | `:58` 逐字命中 | 未腐 |
| `models.go:249 VerifyInstalled`、`:253` / `:323` 两处 `store=%s` | 三处逐字命中 | 未腐 |
| 台账 `:3534` | 逐字命中（该行仍在；文件 20494 行） | 未腐 |

## 2 `AC#0` 三者逐字对照（两发：无/有 `local_override`）

台子：仓外真 `wisp.exe`（`go build -o D:\tmp\...\bin\wisp.exe ./cmd/wisp`，DLL 与 `models/` 按 `build.ps1` 同款方式并置）
＋一枚**真签名**合成清单（`tools/signmodels sign`，dev 键 `E:\work\base\wisp-minisign\wisp-models.key`，keyid
`a3c8794f3fd94fc5`），id `fixture-a1`，`model.onnx` 12 B，pin `206402cab345…`。配置由产品自己首建
（`wisp run "hello"` => `config.toml`，`[models]` 逐字 `dir = ''` / `mirror = []` / `verify_signature = true`）。
读数件：`logs/ac0-shots.txt`（CLI 台）、`logs/ac0-probe.txt`（真 Go 探针，量 `Ensure` 的**返回值**）。

**2.1 先跑正控（票面 `:30` 硬要求）-- 成立**
无 `local_override`、store 里 `model.onnx` 被换成**同长 12 B** 的另一串字节：
```
wisp models verify: model: hand-off re-verification: model fixture-a1 installed files no longer match
the signed manifest in ...\dataA\models\fixture-a1: model.onnx: sha256 mismatch
(pins 206402cab345..., got 950edc0ad300...)      rc=1
```
（先试 14 B 篡改时先报的是 `size 14, manifest pins 12`，故重打成同长，让**哈希**那一跳是开火的那一跳。）
同树不篡改那一发 `rc=0`（A0）=> 台子能测出发育不良。

**2.2 无 `local_override`：三者一致**
| 腿 | 逐字 |
|---|---|
| 守卫验的目录 | `D:\tmp\wisp\20261021\357a1\dataA\models\fixture-a1`（`VerifyDir` 实参，由 `:298` 拼出） |
| CLI 打印 | `wisp models verify: fixture-a1 与已验签清单一致（store=D:\tmp\wisp\20261021\357a1\dataA\models，0.00s）`（`:253`，`store.root`） |
| `Ensure` 真交还 | `D:\tmp\wisp\20261021\357a1\dataA\models\fixture-a1`（探针 P1 现读） |
=> **同一条**。篡改后 CLI `rc=1`（A1），探针 `VerifyInstalled REFUSED`（P2）=> 两腿同向。

**2.3 有 `local_override`：守卫验的目录 与 CLI 打印的目录 不是同一条**（票面味一现场重演）
`[models]` 只加一行：`local_override = { "fixture-a1" = "D:\tmp\wisp\20261021\357a1\override" }`

| 腿 | 逐字 |
|---|---|
| 守卫验的目录 | `D:\tmp\wisp\20261021\357a1\override`（`:299-300` 把 `dir` 换成 override） |
| CLI 打印 | `...（store=D:\tmp\wisp\20261021\357a1\dataB\models，0.00s）`（**仍是 store root**） |
| `Ensure` 真交还 | `D:\tmp\wisp\20261021\357a1\override`（探针 P3/P4：`handed-back == override dir? true`） |

三发读数（每发都是真进程）：

| 发 | store 内容 | override 内容 | CLI 逐字末段 | rc |
|---|---|---|---|---|
| B0 | **目录不存在**（`ABSENT`） | 12 B pin 命中 | `与已验签清单一致（store=...\dataB\models，0.00s）` | **0** |
| B1 | 12 B **篡改**（`950edc0a…`） | 12 B pin 命中 | `与已验签清单一致（store=...\dataB\models，0.00s）` | **0** |
| B2 | 12 B pin 命中 | 12 B **篡改** | `...in ...\override: model.onnx: sha256 mismatch (pins 206402cab345..., got 950edc0ad300...)` | **1** |

**B1 就是票面那味**（"store 里一个字节都没有"的更狠版：字节在、且是篡改的），却报与已验签清单一致，且同一行把
`store=…\dataB\models` 印在"一致"旁边。B2 反证守卫**确实在验 override**（错的一方在 override 时才红）。

**2.4 归因：差的那一格不是"哪次检查"，是"印的是谁"**
`Ensure` 三条返回路径与 `VerifyInstalled` 都调 `handOffCheck(entry, dir)`（`:195/:205/:241/:302`），`dir` 在有 override 时
是**覆盖目录**；而 CLI 两处 `store=%s` 打的都是 `modelStoreDir(...)` 的结果（`models.go:166` => `:191-199`）。
=> "守卫验的目录 == 读取者真正 open 的目录"这一条，今天在 `Ensure`/`VerifyInstalled` **内部是成立的**；
**破裂的是 CLI 的 `store=` 这一行**：它把"验的是谁"说成了"store 里那棵树"。味一的判据应写成这一句。

**2.5 一条今天不成立的前提（必须登记）：交还的目录没有任何人接**
- 唯一生产调用点把返回值**丢掉**：`bridge.go:44` 逐字 `_, err := b.mgr.Ensure(ctx, id, b.onProgress)`。
- `cmd/` 里 `Ensure(` 非测试调用 **0 处**；`cmd/` 与 `internal/` 里对 `<store>/<id>/` 做过 `os.Open|ReadFile|OpenFile`
  的 **0 处**（两把独立尺见 part2 sec 4.2）。
- CLI 打印的是 `ms.root`（`models.go:324/:326`），**不是** `Ensure` 交还的那个值。
=> 味一今天**可观测后果 = 主诉**（CLI 报"与清单一致"而 store 是篡改的）；"读者真的打开错目录"这一后果**可达性 0**
（不是"很小"：根本没有读者）。
