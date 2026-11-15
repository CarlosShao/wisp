# 455 — onnxruntime 绑定复原（recon，只读腿）

**结论一句话：41 枚红不是"摆了 17.8 MB 还被抢"，而是那一发根本没走载具——裸 `go test ./cmd/wisp` 的测试 exe 落在 `%TEMP%\go-buildXXXX\b001\wisp.test.exe`（旁边没有 DLL，也没导出夹具环境）；System32 的 10.5 MB 版按文档顺序先命中。载具本身是对的：同一枚测试经载具 `PASS 2.14s` 并绑到 staged 的 17.8 MB / 1.28.2。那一环＝载具（准确说：载具没有被用上、也不是唯一入口）；不是 sherpa 加载器（它按文档行为工作），也不是 `deps.toml` 与 stage 不一致（sha256 与字节逐位一致）。**

> **对票面一句读数的更正（实测，见 §2.1）**：那 41 红**不是全部**带 `onnxruntime 1.17.1 … shadowing` 句——逐块实测是 **20 枚带该句（绑定层）、20 枚是 `WISP_SPIKE_MODELS is not set`（夹具层）、1 枚是派单已具名的 `TestAC246`**（37.74s）。两层的共同根因都是"这一发没走载具"：载具既 stage DLL，也导出五个夹具变量；裸 `go test` 两样都没有。onnxruntime 那一半（20 枚）才是"绑错版本"的直接受害者。

- 实测绑到：载具形 → `D:\work\workspace\projects plans\Wisp\.scratch\wisp\probes\455\stage\onnxruntime.dll`，**17799168 B**，FileVersion **1.28.2**；裸 `go test` 形 → `C:\Windows\SYSTEM32\onnxruntime.dll`，**10572960 B**，FileVersion **1.17.260311-1434.1**（同一台机、同一条测试、同一个 cwd，只有 exe 所在目录不同）。
- 探针：`.scratch/wisp/probes/455/`（`probe_binding.go`、`loadlib_probe/main.go`、`run-probe.ps1`），一次跑完 A/B/C1/C2/C3 五臂，读数落 `.scratch/wisp/probes/455/run-probe.out.txt`。

---

## 1. 实测绑的是哪一枚（AC#1：要路径与字节，不拿版本号当证据）

探针用仓库已有的同一枚仪器（`EnumProcessModules` + `GetModuleFileNameEx`，见 `internal/speech/cgo_boundary_windows_test.go:34-59`），从**进程内部**打实际加载路径，再 `os.Stat` 取字节，另记 `sherpa.GetOnnxruntimeVersion()`。

```
powershell -NoProfile -ExecutionPolicy Bypass -File .scratch\wisp\probes\455\run-probe.ps1
```

| 臂 | exe 所在目录 | 实际绑定的 `onnxruntime.dll` | 字节 | 版本 |
|---|---|---|---|---|
| **A 载具形** | `probes/455/stage`（三枚 DLL 与 exe 同目录） | `...\probes\455\stage\onnxruntime.dll` | **17799168** | **1.28.2** |
| **B 裸 `go test` 形** | `probes/455/tmp-raw`（无 DLL；sherpa 走 PATH） | `C:\Windows\SYSTEM32\onnxruntime.dll` | **10572960** | **1.17.1** |
| **C1** 显式 `LoadLibraryEx(<stage>\sherpa-onnx-c-api.dll, 0)` | `tmp-raw` | `C:\Windows\SYSTEM32\onnxruntime.dll` | 10572960 | — |
| **C2** 同上，`flags=0x8`（`LOAD_WITH_ALTERED_SEARCH_PATH`） | `tmp-raw` | `...\stage\onnxruntime.dll` | 17799168 | — |
| C3 同上，`flags=0x100`（`LOAD_LIBRARY_SEARCH_DLL_LOAD_DIR`） | `tmp-raw` | `LoadLibraryEx error: The specified module could not be found.`（该臂本机组合未成立，**不用于结论**；C2 已给出同一机制的确定读数） | — | — |

候选文件现量（FileVersion 取自文件版本资源，非报错文本）：

| 文件 | 字节 | FileVersion |
|---|---|---|
| `C:\Windows\System32\onnxruntime.dll` | 10572960 | `1.17.260311-1434.1.os-germanium.6bcd20d` |
| `third_party/sherpa-onnx/onnxruntime.dll` | 17799168 | `1.28.2` |
| `.scratch/wisp/stage-cli/onnxruntime.dll` | 17799168 | `1.28.2` |

**端到端对照（同一条测试 `^Test15eSilenceNeverCallsTheHandOff$`，本机现跑 2026-11-15 19:35-19:36）**

| 形 | 命令 | 读数 |
|---|---|---|
| 载具 | `WISP_CLI_STAGE_DIR=<probes/455/stage-cli-run> scripts/wisp-cli-staged-tests.sh '-test.run' '^Test15eSilenceNeverCallsTheHandOff$'` | `--- PASS: … (2.14s)`；`ok github.com/CarlosShao/wisp/cmd/wisp 2.00s`；语音链真起来（silero VAD 装上了、`consumer booked on the D38b roster`）。日志 `probes/455/carrier-15e.log` |
| 裸 `go test` | `$env:PATH="<root>\third_party\sherpa-onnx;"+$env:PATH; go test ./cmd/wisp -run '^Test15eSilenceNeverCallsTheHandOff$' -count=1 -v` | `--- FAIL: … (0.00s)`，逐字 `speech: onnxruntime 1.17.1 is older than the 1.28 … a second onnxruntime.dll is shadowing the one on PATH`。日志 `probes/455/raw-go-test-15e.log` |

⇒ **AC#1 的那一对数字**：staged/钉版是 **17799168 B**；裸 `go test` 形（41 红那一发里的**绑定层 20 枚**）实际绑的是 **10572960 B**（`C:\Windows\SYSTEM32`）。

---

## 2. 那一环在哪（逐条现读，⛔ 猜）

### 2.1 `scripts/stage-sherpa-dlls.sh` 在这一发拒没拒？——**这一发根本没调它**

41 红的那一发是 `.scratch/35g3-wholepkg-skip-death.log`（票 35 G3「发 2」，`PASS=618/FAIL=41/SKIP=10`、`752.948s`）。现读该日志：

| 读什么 | 读数 |
|---|---|
| 载具横幅 `wisp-cli-staged-tests.sh:` | **0 行** |
| `stage-sherpa-dlls: staged …` | **0 行** |
| 夹具未导出的 `WISP_SPIKE_MODELS is not set` | **19 行**（载具会在夹具目录存在时导出；本机 `third_party/spike-models` 在盘上） |
| 末行 | `FAIL	github.com/CarlosShao/wisp/cmd/wisp	752.948s`（`go test` 原生包行；载具会打印自己合成的 `FAIL  github.com/... 752.00s`） |

**41 红分层实测**（按每个 `--- FAIL` 块内是否出现遮蔽句归因）：

| 层 | 枚数 | 逐字签名 | 归因 |
|---|---|---|---|
| 绑定层（onnxruntime） | **20** | `speech: onnxruntime 1.17.1 is older than the 1.28 … shadowing the one on PATH` | 裸跑无 staging → 绑 System32 的 10572960 |
| 夹具层 | **20** | `install the streaming ASR into <tmp>: WISP_SPIKE_MODELS is not set` | 裸跑没导出五个夹具变量（载具会导） |
| 他因已知红 | **1** | `TestAC246DevLegIgnoresTheTestTaskInjection`（37.74s） | 派单已具名，与本票无关 |
| **合计** | **41** | — | — |

⇒ "摆了 17.8 MB 却绑了 10.5 MB"这一对数字，直接落在**绑定层的 20 枚**上；另 20 枚连绑定都还没轮到就先死在夹具层。两者都因为这一发不是载具形。

票 35 G3 证据件自己写明的命令（`docs/evidence/s1/35-g3-backpressure-20261115.md:166-171`）：

```
$env:PATH = "...\third_party\sherpa-onnx;...\build;$env:PATH"
go test ./cmd/wisp -count=1 -timeout 1800s -v
```

`go test -x` 现跑证明裸 `go test` 的测试二进制落在构建缓存临时目录（旁边没有任何 DLL）：

```
WORK=C:\Users\swq\AppData\Local\Temp\go-build2507224625
link.exe -o "$WORK\b001\wisp.test.exe" ...
$WORK\b001\wisp.test.exe -test.run=^$ ...
```

⇒ 那一发**没有 staging 发生**（不是"拒了"），因为 `go test` 的 exe 不在 `stage-cli` 里；`.scratch/wisp/stage-cli/onnxruntime.dll` 的 17799168 B 是 **18:43:49** 另一发（`wisp.test.exe` 18:43:48、定向 Test46 跑、PASS）留下的，比 19:14:15 起的整包发早半小时。"摆了 17.8 MB"是真的——但**摆在了另一个目录里，不是那一发进程的 exe 目录**。

今天按载具重跑：`stage-sherpa-dlls: staged onnxruntime.dll -> …/stage-cli-run/onnxruntime.dll` 等三枚全部落地（`probes/455/carrier-15e.log` 首段）。

### 2.2 `deps.toml` pin 的那一枚是不是 17799168 B？——**是，且字节/sha256 逐位一致**

| 对象 | 读数 |
|---|---|
| `deps.toml` `[sherpa-onnx.dll."onnxruntime.dll"].sha256` | `7f66f939a881baf4f46a2216496798edf4a1429878b646d12674aa62f27d8a25` |
| `third_party/sherpa-onnx/onnxruntime.dll` 现算 SHA256 | `7f66f939…27d8a25`（**一致**），17799168 B，1.28.2 |
| `.scratch/wisp/stage-cli/onnxruntime.dll` 现算 SHA256 | `7F66F939…27D8A25`（**一致**），17799168 B，1.28.2 |

⇒ **不是 "deps.toml 与 stage 不一致"**：pin、third_party 源件、stage 出来的件是同一枚（同 sha256、同字节）。`deps.toml:onnxruntime.version` 也写着 `1.28.2`。

### 2.3 sherpa 的 `sherpa-onnx-c-api.dll` 用 `LoadLibrary` 按自己的搜索序找 onnxruntime？起点是"加载它的模块所在目录"还是 exe 目录？——**决定一切的那条：起点是 exe 目录，不是加载它的 DLL 所在目录**

Microsoft 官方《Dynamic-link library search order》逐字（2023-02-09 版）：

> "If a DLL has dependencies, then the system searches for the dependent DLLs **as if they were loaded by using only their module names. That's true even if the first DLL was loaded by specifying a full path.**"

> SafeDllSearchMode 开启（默认）时标准序：… 4. loaded-module list → … 7. **the folder from which the application loaded** → 8. `%SystemRoot%\system32` → … 11. current folder → 12. **`PATH`**。

> 只有 `LoadLibraryEx` 带 `LOAD_WITH_ALTERED_SEARCH_PATH`（且 lpFileName 是绝对路径）时，第 7 步才换成"**被加载的那个模块所在的目录**"；`LOAD_LIBRARY_SEARCH_DLL_LOAD_DIR` 同理（这是 `LoadLibraryEx` 的扩展搜索集）。

本机实测（上表）：C1 用绝对路径 `LoadLibraryEx(<stage>\sherpa-onnx-c-api.dll, flags=0)`，onnxruntime **仍落 System32**（10572960）⇒ **加载它的那个 DLL 的目录不在标准序里**；C2 加 `0x8` 后同一调用立刻拿到 `<stage>\onnxruntime.dll`（17799168）⇒ 机制坐实。

⇒ 推论（这就是"载具知识"）：**staging 必须摆在"那个 exe 自己的目录"旁边**。摆在 `.scratch/wisp/stage-cli/` 而 exe 在 `%TEMP%\go-build…`，等于白摆；这也解释了为什么 `PATH` 只救一级导入（sherpa 本体从 PATH 命中），救不了二级依赖（onnxruntime 在 System32 之后才轮到 PATH）。

### 2.4 测试进程里有没有人改 `PATH` / 调 `SetDllDirectory` / `AddDllDirectory`？——**产品与测试 Go 码零调用；只有 shell 侧一处 PATH 前置**

| 搜索面 | 读数 |
|---|---|
| `*.go` 里 `SetDllDirectory` / `AddDllDirectory` / `os.Setenv("PATH"` | **0 处**（唯一 `LoadLibraryEx` 命中是 `internal/speech/cgo_boundary_windows_test.go:196`，那是"复制一份同名 DLL 验证计数器能看见两张"的对照腿，不参与加载序） |
| shell 侧 | `scripts/wisp-cli-tests.sh:220` `export PATH="$dll_dir:$PATH"` —— 只把 `third_party/sherpa-onnx` 前置，救一级导入；**载具本身不改 PATH**（它只 stage + 导夹具） |
| 41 红那一发的环境 | 它自己 `$env:PATH = "...\third_party\sherpa-onnx;...\build;$env:PATH"`（票 35 G3 证据件） |

⇒ 没有任何"偷偷改搜索序"的地方；`SetDllDirectory` 空串/null、`AddDllDirectory` 全仓零调用，排除"有人把搜索序改坏"。41 红是**标准序的正常结果**：exe 在无 DLL 的临时目录 → System32 先命中。

### 2.5 旁证：载具形的整包读数是 0 条 onnx 遮蔽红（不是本腿新测，是既有日志复算）

| 日志 | 形 | RUN | PASS | FAIL | SKIP | `shadowing the one on PATH` 命中 | 载具横幅 | staged 行 |
|---|---|---|---|---|---|---|---|---|
| `probes/451/r1/cmdwisp-full.txt`（11-14 整包，载具） | 载具 | 866 | 623 | 3 | 17 | **0** | 4 | 3 |
| `probes/452/cmdwisp-full.log`（11-14 整包，载具） | 载具 | 868 | 626 | 2 | 17 | **0** | 4 | 3 |
| `probes/452/cmdwisp-full-run2.log`（同上第二发） | 载具 | 868 | 626 | 2 | 17 | **0** | 4 | 3 |
| `.scratch/35g3-wholepkg-skip-death.log`（11-15 41 红） | **裸 `go test`** | 894 | 618 | **41** | 10 | **50 行**（按 FAIL 块归因 20 枚，见 §2.1） | **0** | **0** |

⇒ 载具形整包**跨三发、0 条遮蔽红**；41 红那一发是唯一没有载具横幅、没有 staged 行、却出现 19 行"夹具未导出"的整包发。

---

## 3. 最小复现（30 秒内回答"绑的是哪一枚"）

位置（⛔ 不在产码树，全在 `.scratch/wisp/probes/455/`）：

| 件 | 作用 |
|---|---|
| `probe_binding.go` | cgo 版探针：import sherpa 让导入树在进程载入时解析，然后从进程内打 onnxruntime/sherpa 的**实际路径 + 字节 + 版本** |
| `loadlib_probe/main.go` | 无 cgo 版：显式 `LoadLibraryEx(路径, flags)`，回答"加载它的 DLL 的目录算不算搜索起点" |
| `run-probe.ps1` | 一键：建目录 → `go build` 两枚探针 → 把三枚 DLL 拷到 `stage/`（载具形）→ 跑 A/B/C1/C2/C3 五臂 |
| `run-probe.out.txt` | 本机现跑读数（上面两张表的原始输出） |
| `carrier-15e.log` / `raw-go-test-15e.log` | 端到端 PASS vs FAIL 的两发原始日志 |

命令：

```
powershell -NoProfile -ExecutionPolicy Bypass -File .scratch\wisp\probes\455\run-probe.ps1
```

判据（12 行内可读）：
- `ARM A` 行必须出现 `path=...\probes\455\stage\onnxruntime.dll` + `size=17799168 B` + `GetOnnxruntimeVersion(): 1.28.2`；
- `ARM B` 行必须出现 `path=C:\Windows\SYSTEM32\onnxruntime.dll` + `size=10572960 B` + `GetOnnxruntimeVersion(): 1.17.1`；
- 两行只差 `exe dir`，cwd 都是 `cmd\wisp` —— 这就是"绑哪一枚只由 exe 目录决定"的一发复现。

---

## 4. 修法建议（按主流 harness 形状，3 候选各写代价）

> 判据按机主铁律：**最有价值 ＋ 最符合主流 harness**，⛔ 按"风险最小/动作最小"选。签名/契约不用改（本票只动脚本与 test-only 件）。

### 候选① 把载具变成"唯一入口"，并在 stage 后核 sha256（现在只核"存在"）

- **做什么**：`portable-tests.sh` 的 `cli` 分支已经指向载具（现状成立）；补两件：(a) 门禁/README 明确"裸 `go test ./cmd/wisp` 不是整包读数"（`runtests.sh`/CI step 只认载具）；(b) 载具在 `stage-sherpa-dlls.sh` 之后逐枚比对 `sha256 == deps.toml` pin（现在只 `[ -f ]` 验存在，验不出"拷了错文件"）。
- **代价**：脚本面 ~1 小时，零产品码、零断言改动。**缺点**：拦不住人手敲裸 `go test`（只拦门禁读数），也不改"红句读起来像产品缺陷"这一点——所以单独用不够。
- **撤销口令**：`git revert <本候选 commit>`（脚本面，不动 Go）。

### 候选② 给 `cmd/wisp` 加一枚"绑定自检"（fail-fast 环境前置）——**推荐与①并用，价值最大**

- **做什么**：test-only 一个文件（`cmd/wisp/*_test.go`，⛔ 不改产品行为）：启动时取 `GetModuleFileName(onnxruntime.dll)`＋`os.Stat` 字节，与"exe 同目录那份"比对（期望 17799168）；不符即 `t.Fatalf`，逐字打「bound=%s %d B，pinned=%d B；请走 scripts/wisp-cli-staged-tests.sh」。
- **代价**：~60 行 test-only，零产品码、零断言放宽。**效果**：裸 `go test` 从"41 枚像产品缺陷的红"变成**一条**具名环境红，且自带路径+字节（下次再有人误读整包，红句自己就把答案说了）。这正是主流 harness 的做法：环境前置 fail-fast。
- **代价落点**：门禁口径要一次性调整——今后裸跑整包会稳定出现这 1 条环境红（原本是 41 条），需要把它记为环境红而非产品红。
- **撤销口令**：`git revert <commit>`（删该 test-only 文件，恢复 41 红原状）。

### 候选③ 把产品加载改成显式 `LoadLibraryEx(绝对路径, LOAD_WITH_ALTERED_SEARCH_PATH)`

- **实测可行**：C2 证明 `flags=0x8` 能稳定拿到同目录 1.28.2。
- **代价（高）**：sherpa 现在是 cgo **隐式导入**（导入表在进程载入时解析），要改成运行时 dlopen 的薄 shim / 加一层转发 DLL，动到产品加载路径与 cgo 边界；收益只是"省掉载具"。按"最有价值"判据——**不选**（复杂买不到真价值）。
- **撤销口令**：`git revert <commit>`（但此候选不建议落地，若已落则须连 cgo 边界一起回滚）。

> 追加一句（不单列候选）：`go test -exec <staged-wrapper>` 也能让"裸 go test 自动 stage"，代价中等，但 Windows 引号/空格与 wrapper 维护面比①大；若采纳②，收益已被②的 fail-fast 覆盖大半，故不单列。

---

## 5. 边界与引用

- ⛔ 未改产码／票面／台账；⛔ 未动 `C:\Windows\System32`（只读量它的路径/字节/版本）；⛔ 未 push。
- 写面：本件 `docs/evidence/s1/455-onnxruntime-binding-recon.md` 与 `.scratch/wisp/probes/455/**`。探针 `stage/`、`tmp-raw/`、`stage-cli-run/`（~22 MB DLL 副本 + 63 MB test.exe）是临时构建产物，**不入库**（只 commit `.go`/`.ps1`/日志）。
- 文档引用：
  - Microsoft Learn, *Dynamic-link library search order*（2023-02-09）: https://learn.microsoft.com/en-us/windows/win32/dlls/dynamic-link-library-search-order
  - `scripts/wisp-cli-staged-tests.sh:12-18`（现有头注释断言"exe 目录先于 System32"——对；缺的只是"staging 必须挨着 exe、裸 `go test` 拿不到"这句）。
  - `scripts/stage-sherpa-dlls.sh`（deps.toml 驱动、缺件响亮拒绝）。
  - `docs/evidence/s1/cmd-wisp-dll-staging-carrier.md`（票 409 治法迁到 `cmd/wisp` 的原始证据）。
  - 票 35 G3 读数：`docs/evidence/s1/35-g3-backpressure-20261115.md:164-187`。
