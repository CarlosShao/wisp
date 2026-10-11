# 308-a1 · 90-final（票 308 `AC#0` 交件；只读普查腿终局）

锚笔 `02053534`；分支 `dev`（共享工作树，起手 `git status --porcelain -- internal cmd`＝0）。
本腿**只读码＋只跑测试取颜色＋只新建 `probes/308/a1/**`**，⛔ 改任何已跟踪文件。

---

## 1. 四把台面读数（逐发：env verbatim ＋ 在场尺 ＋ 六尺颜色 ＋ 崩位 ＋ panic 逐字）

名册尺（T1–T3）＝`go test ./internal/risk/ -count=1 -timeout 600s -v -run "TestPathResolverUNCAListDenied|TestClassifySpellingInvariance"`；
T4／对照＝同上去掉 `-run`（整包）。台面 env 用 cmd.exe 里 `set` 绕开 MSYS 改写；`GOTMPDIR` 固定给 build，使 test 子进程的 temp 与 build 解耦。

| 发 | 台面 env（verbatim） | 在场尺 `os.MkdirTemp`→`VolumeName` 空？ | rc | 顶层 RUN/PASS/FAIL/SKIP | 子项 PASS/FAIL | 包级终态行（verbatim） | panic | 崩的那枚用例 |
|---|---|---|---|---|---|---|---|---|
| **T1 posix** | `TMP=/tmp TEMP=/tmp`（POSIX 形，票面既有触发形） | `D:\tmp\W308PROBE…` → `"D:"`，**非空** | 0 | 2/2/0/0 | 0/0 | `ok github.com/CarlosShao/wisp/internal/risk 0.138s` | 0 | **⛔ live panic（两枚都跑到且 PASS）** |
| **T1 bare-root**（附加一发，见 `logs/T1-posix-roster.txt` 同形＋`probe-T1-bare-root.txt`） | `TMP=\ TEMP=\` | `D:\W308PROBE…` → `"D:"`，非空 | — | （在场尺发；颜色同 T1 posix） | — | — | — | ⛔ |
| **T1 default**（本工具继承，不设 env） | shell 内 `TMP=/tmp`，MSYS 给子进程改写成 Windows 形 | `C:\Users\swq\AppData\Local\Temp\W308PROBE…` → `"C:"`，非空（`probe-default-D.txt`；直跑两枚名册亦 PASS，见 `logs` 首跑） | 0 | 2/2/0/0 | 0/0 | `ok github.com/CarlosShao/wisp/internal/risk` | 0 | ⛔ live panic |
| **T2 vol**（本机默认） | `TMP=TEMP=C:\Users\swq\AppData\Local\Temp` | 同 default，`"C:"` 非空 | 0 | 2/2/0/0 | 0/0 | `ok github.com/CarlosShao/wisp/internal/risk 0.081s` | 0 | ⛔ live panic |
| **T3 alias**（8.3 短名） | `TMP=TEMP=C:\Users\swq\tmp\RISKAT~1\RUNNER~1`（长形 `…\riskattrib1\RunnerTempAreaForAttribution1`，本腿⛔ 删里面任何东西） | `C:\Users\swq\tmp\RISKAT~1\RUNNER~1\W308PROBE…` → `"C:"`，非空（`T3-alias-probe.txt`） | 0 | 2/2/0/0 | 0/0 | `ok github.com/CarlosShao/wisp/internal/risk 0.060s` | 0 | ⛔ live panic |
| **T4 full**（posix台面·整包名册） | `TMP=/tmp TEMP=/tmp` | 同 T1（非空） | 0 | 132/131/0/1 | 68/0 | `ok github.com/CarlosShao/wisp/internal/risk 2.711s` | 0 | ⛔ live panic |
| **对照 full**（vol台面·整包） | `TMP=TEMP=C:\Users\swq\AppData\Local\Temp` | 非空 | 0 | 132/131/0/1 | 68/0 | `ok github.com/CarlosShao/wisp/internal/risk 4.425s` | 0 | ⛔ live panic |

**混合层级标注**：T4／对照整包"绿"合计数＝`131 顶层 ＋ 68 子项 ＝ 199`（与 `A843` §3 那把"131顶层＋68子项"尺**逐枚同数**）；引它必标"含子项"，⛔ 当纯顶层。
T4 顶层 `RUN=132` ＝ 对照 `132` ⇒ **posix 台面这一发整包其余用例全拿到结果**（⛔ 少跑）。

**本腿没有 live panic 的逐字依据**：每发 `panic|fatal error|test executable failed` 命中＝**0**；包级终态行 `ok …` 在（⛔ 是 `FAIL …` 收尾、⛔ 是名册截断）。

---

## 2. 一句话结论（靠三把尺，⛔ 靠猜）

**在本机（win32＋Go 1.27.1＋这枚 MSYS 工具），票 308 说的两枚崩点用例（`TestClassifySpellingInvariance`、`TestPathResolverUNCAListDenied`）哪一枚都⛔ 真崩——三把台面全部 `PASS`，同包其余用例在整包发里 `RUN=132` 全求值、终态行 `ok` 在、`panic` 命中 0；"崩⇒其余零求值"那一族形状本腿⛔ 能现场触发。**
三把尺一起说＝① 顶层 RUN 枚数⇄正常发**同数**（132⇄132，无早死）② 包级终态行是 `ok github.com/CarlosShao/wisp/internal/risk`（⛔ 缺、⛔ 非 `FAIL`）③ `grep -cE 'panic|fatal error|test executable failed'`＝**0**。
根因＝`os.TempDir()` 在 Windows 走 `GetTempPath`，任何**能建目录成功**的 `TMP/TEMP` 形都被补上当前盘符 ⇒ `filepath.VolumeName(long)` 恒非空 ⇒ `vol[:1]`⛔ 越界；无卷标台面在本机造不出来（相对形／UNC 形虽 drive-less，但 `os.MkdirTemp` 建失败或要写进 `internal/risk`＝本腿写面禁止）。
**历史 live panic 的唯一原始件**＝`probes/72/a2/logs/t03-broken-surface-note.txt`（`rc=1`／名册 2 枚／`PASS=0`／`FAIL=2`），逐字：

    panic: runtime error: slice bounds out of range [:1] with length 0 [recovered, repanicked]

那发崩在 `TestClassifySpellingInvariance`（`anchor_spelling` 体内那处 `vol[:1]`，历史行号 `:78`／HEAD 现量 `:103`），崩时同包其余用例零求值。⚠ 本腿⛔ 复现它。

---

## 3. ⒶⒷⒸ 读数（各带尺名＋射程目录）

### Ⓐ 两枚注释缺陷（复跑 `A843` §2 的尺；详见 `20-comment-defects.md`）
- ① 指针抄错：`anchor_spelling:157`（现量，非 `A843` 的 `:158`）指 `probes/72/r1/20-canonical-form-choice.md`，该件⛔ 在（`ls`＋`find` 双尺坐实；`r1/` 全名册只有 00/10/20-class-b/30 四件 md）；选型真身在 `10-class-a-expected-side.md` 的 `## #2` 一节。
- ② 现时假话：同文件 `:15`（`:12-16` 段）称 `pathresolver_junction_windows_test.go` "byte-for-byte untouched"，而 `git show --numstat --format= f278c57f -- …junction…`＝**`36 7`**（那笔自己就改了它）。
- 归口＝票 308；本腿⛔ 改（裁决者⛔ 动手改实现面）。

### Ⓑ 产码⛔ 在射程内（复跑 `grep -rn '\[:1\]' internal cmd --include=*.go`＝7 枚；详见 `30-codecensus.md`）
- 射程＝**2 枚**：`anchor_spelling:103`、`junction:152`（`vol := filepath.VolumeName(…)` 后⛔ 守卫的 `vol[:1]`）。
- ⛔ 射程＝5 枚：`pathresolver.go:263`（`drive[:1]`，前有 `len(drive)==2 && drive[1]=='$'` 守卫）；`spill_path_invariant_test.go:389/:406`、`mockllm_integ_test.go:190`、`firstrun_198r2_test.go:401`（非 `VolumeName` 崩形／在 `cmd/wisp`＝本腿⛔ 打的面）。
- 产码 `pathresolver_windows.go:60 reparseComponents()`：`:61` 取 `vol` 但⛔ 做 `[:1]`（该文件 `[:1]` 命中 0）。

### Ⓒ CI 档位（只取现量，⛔ 下判语；⛔ 拿"本机跑过"顶替 CI 色）
- 两枚文件 build tag（`head -3`）＝各第 1 行 `//go:build windows` ⇒ 非 windows GOOS 下⛔ 参编译、⛔ 求值。
- `grep -rn 'winlive' internal/risk/*_test.go`＝**0 命中**（`winlive_hits=0`）。
- `.github/workflows/` 里打 `internal/risk` 的步（尺＝`grep -n` 定位后逐行读）＝
  - `ci.yml:797`：`bash scripts/portable-tests.sh --scope=windows`（`internal/risk` 在该 windows scope 内：`portable-tests.sh:186/:202` 列 `github.com/CarlosShao/wisp/internal/risk`）；
  - `ci.yml:813`：`bash tools/d22scan/runtests.sh ./internal/risk/ -run TestPathResolverJunctionWindows`（只打 Case 1，**⛔ 覆盖两枚崩点用例** `TestClassifySpellingInvariance`／`TestPathResolverUNCAListDenied`）。
- **CI 实跑颜色＝本腿拿不到**（⛔ push、⛔ 跑 CI）；上面只交给哪一步能打到它们，判语归编排者／`AC#3`。

---

## 4. ⛔ 做到的事（具名，⛔ 用"应该没问题"含过去时）

- **CI 色本腿拿不到**：本腿⛔ push、⛔ 触发 GitHub Actions；只交了两枚崩点在 CI 哪一步会被／不会被 `-run` 打到（`ci.yml:797` 整包 scope、`ci.yml:813` 只 Case 1），windows 档在 test-windows 里的**实跑 PASS/FAIL 色**＝欠读数。
- **live panic 本腿造不出**：本机 `GetTempPath` 恒给 drive-qualified temp，无卷标台面在场尺（`os.MkdirTemp`→`VolumeName`）在 T1/T2/T3/default/bare-root **全部非空**；§2 的历史 panic 来自 `72-a2` 原始件，⛔ 本腿复现。⇒ "哪一枚真崩"这一问，本腿的回答＝**在本机这两枚都不崩**；真崩形状只能引历史件，⛔ 现量。
- **POSIX 档求不求值本腿⛔ 判**：两枚文件带 `//go:build windows` ⇒ 在 darwin/linux 档⛔ 编译、⛔ 求值；本腿只在 windows 档跑，⛔ 能替"POSIX 档也跑了"背书。
- **⛔ 跑 `go build`／`go vet`／`d22scan`／长度门**：Go 门禁面此刻由另一枚腿 `298-v1` 独占（`A843` §3 在飞那枚）；本腿只调 `go test`，且**⛔ 打 `cmd/wisp`**（那面一次只给一枚腿，`firstrun_198r2_test.go` 只在 Ⓑ 名册里静态列名、⛔ 求值）。
- **共享索引里躺着⛔ 我名的改动**：`git status --porcelain` 整仓 869 行脏（`.gitignore`、`design/**`、别人 `probes/**`、根上 `-` 等），我⛔ 读过其内容、⛔ 动、⛔ 提交；起手 `internal cmd` 干净＝0 行；每笔 commit 带显式 pathspec `.scratch/wisp/probes/308/a1`，越界枚数见 §5（全 0）。
- **⛔ 翻框／⛔ 改产码／⛔ 为变绿改断言**：本腿只交颜色，守卫落点归编排者裁（`AC#1`/`AC#2`/`AC#3` 判语⛔ 由落地腿自勾，本腿亦代⛔ 勾）。

---

## 5. commit 台账（逐笔 numstat ＋ 越界枚数）

| 短号 | 内容 | numstat（加/删） | 越界枚数（`git show --name-only --format=` 中出我写面者） |
|---|---|---|---|
| `c023e8e0` | 00-anchor（commit-first） | `47  0  .scratch/wisp/probes/308/a1/00-anchor.md` | **0** |
| `dd638f98` | 10/20/30 readouts | `59 0`＋`32 0`＋`24 0`（共 3 件） | **0** |
| 本笔（90-final） | 90-final.md | 见 commit | 0（待提交后自量） |

⛔ push；⛔ `git add -A`/`.`/`--amend`/`reset`/`rebase`/`stash`/`checkout .`/`clean`/`--no-verify`。
原始 `go test`／探针输出全在 `logs/`（只建不删、名⛔ 叫 `.out`）：`T1-posix-roster.txt`／`T2-vol-roster.txt`／`T3-alias-roster.txt`／`FULL-t4-posix.txt`／`FULL-vol.txt`／`probe-*.txt`。
