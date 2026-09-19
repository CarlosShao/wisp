# T01 对抗验收报告 — 01-build-chain

- 验收角色：T01-adv（对抗验收，独立上下文，只读仓库 + 系统临时目录实验）
- 验收时间：2026-09-19T06:35Z ~ 06:55Z
- 被验收 HEAD：`29c8e35d3459c5cd8306d8af19085667b4068259`（origin/cnb 的 dev 与本地一致，origin 需重试 1 次，TLS 偶发失败符合已知网络现实）
- 验收环境：Windows 11 x64，Git Bash；Go `go1.27.1`（`D:\work\base\go\bin`）；GCC `(Rev3, Built by MSYS2 project) 16.2.0`（`E:\work\base\msys64\mingw64\bin`）——与 BUILD.md §1 钉死的工具链一致
- 实验副本（全部位于系统临时目录，非仓库）：
  - Clone A：`C:\Users\swq\AppData\Local\Temp\wisp-adv-15393\wisp`（干净构建 + 篡改测试）
  - Clone B：`C:\Users\swq\AppData\Local\Temp\wisp-adv2-15625\wisp`（全冷构建，空 GOMODCACHE/GOCACHE）

## 逐项裁决

### 1. Stub / 假完成扫描 — PASS

- 全仓 `grep -rInE 'todo!|unimplemented|panic\("not implemented"|TODO|FIXME'`（排除 third_party/、build/、.git）：
  代码文件**零命中**。仅命中规则文本自身（`docs/PLAN.md:730,748,949`、`docs/specs/SPEC-10:108`——它们是 D22/验收规则的定义）与 `.scratch/wisp/issues/22-web-tools-d30.md:49`（未来票据的既有备注，非本票产物、非代码）。
- `wisp doctor` 每项 PASS 均为真实检查（读 `cmd/wisp/doctor.go` 逐行核实）：
  - `gcc (build-time)`：真实 `exec gcc --version`；
  - `sherpa-onnx C API`：真实 cgo 调用 `sherpa.GetVersion()` 并与构建时注入 pin 比对；
  - `onnxruntime.dll colocated`：真实读 PE `VS_FIXEDFILEINFO`（version.dll Win32 API），容忍 `.0` 尾差；
  - `DLL colocated`：`os.Stat` exe 同目录；
  - `deps.toml * pin`：解析 deps.toml 与 `debug.ReadBuildInfo()` 链接版本交叉核对；
  - `data dir writable`：真实 mkdir/写探针/删除。
  - **FAIL 分支也被证明是活的**：临时 clone 把 deps.toml `version = "1.13.8"` 改成 `1.13.9` 后，
    `wisp doctor` 输出 `[FAIL] deps.toml sherpa-onnx pin  deps.toml says 1.13.9, binary pinned 1.13.8`，
    末行 `wisp doctor: FAIL`，exit 1（同一次运行里 gcc 不在 PATH 也正确 `[FAIL]`）——不存在常量 PASS。

### 2. 独立复现干净构建 — PASS

Clone A（干净 clone，无 third_party/、无 build/，go module 缓存为既有共享缓存）：

- `powershell -ExecutionPolicy Bypass -File scripts/build.ps1 -Env prod` 一次跑通，**13.9s**。
  直连 GitHub 首次尝试即成功，`archive SHA256 verified`，三个 DLL `sha256 ok` 逐个提取，
  `go build ok (cgo linked against sherpa-onnx C API)`，冒烟 `wisp doctor: PASS`（10 项 PASS，0 FAIL）。
- 产物齐全：`build\wisp.exe`(6,273,998B) + `onnxruntime.dll` + `sherpa-onnx-c-api.dll` + `sherpa-onnx-cxx-api.dll` + `SHA256SUMS`；
  GNU `sha256sum -c SHA256SUMS` → 4× `OK`（LF 行尾写法有效）。
- 二次运行（缓存命中）：**4.4s**，输出 `fetch-deps: cache hit - third_party/sherpa-onnx matches deps.toml (sherpa-onnx 1.13.8)`，约 3.2× 提速。
- **加强项（全冷构建）**：Clone B + 全新 `GOMODCACHE`/`GOCACHE`（空目录）→ 重新从 goproxy 拉取
  `sherpa-onnx-go v1.13.8`/`sherpa-onnx-go-windows v1.13.8`/`x/sys`（sumdb 校验随下载自动进行），
  全量 cgo 编译，**56.5s** 跑通，doctor PASS。这覆盖了「干净机器 = 干净 clone + 空缓存 + 全新依赖下载」的最严格复现
  （真正的裸 OS + CI runner 属票 08 的 CI 范畴，BUILD.md §7 也如此声明）。
- **cgo 真链接证据**：`objdump -x wisp.exe` 导入表含 `sherpa-onnx-c-api.dll`；
  `sherpa-onnx-c-api.dll` 自身导入表含 `onnxruntime.dll`（wisp.exe → c-api → onnxruntime 的加载链）。
  把三个 DLL 移走后 `./wisp.exe version` 进程根本无法启动（exit 127，Windows 加载器 0xC0000135），
  还原后正常运行——BUILD.md §4「doctor 打印 PASS 即加载器级证明」的论断实测成立。
- 卫生检查：`gofmt -l cmd internal` 无输出；`CGO_ENABLED=1 go vet ./cmd/... ./internal/...` 通过
  （不带 CGO 时 vet 对 sherpa_onnx 包报 build-constraint 错误，属预期环境差异，非代码缺陷）。

### 3. 篡改测试复现 — PASS

在 Clone A 逐一复现（实现报告宣称的两种篡改 + 我追加的第三种）：

1. **改 deps.toml 整包 SHA256 一个字节**（`6dff…`→`7dff…`）：manifest 与 pin 不一致 → 缓存被判无效 → 重新下载 →
   `fetch-deps: FATAL: SHA256 MISMATCH for the downloaded archive. expected: 7dff… actual: 6dff…`，
   附「Do not edit the pin to make this pass」警示，**exit 1**，坏文件已删除。
2. **改 per-DLL pin 一个字节**（c-api `300e…`→`310e…`）：整包校验通过后逐 DLL 比对失败 →
   `fetch-deps: FATAL: SHA256 MISMATCH for lib/sherpa-onnx-c-api.dll`，**exit 1**。
3. **改 third_party 缓存里 c-api.dll 一个字节**（deps.toml 干净）：缓存重哈希发现不符 →
   自动重新下载并逐 DLL 校验后**修复回 pin 哈希**（修复后 `sha256sum` = `300e0c…`，与 pin 一致），exit 0。
   该路径语义是「缓存永不信任，损坏即重取修复」（fetch-deps.ps1 头注与 BUILD.md §3 一致），
   而非静默通过——检测必然发生；若连修复下载也错会落入路径 1 的 FATAL。可接受。
4. 附带收获：路径 3 的重取过程实际复现了 BUILD.md 坑 6——「直连 TLS 失败 → 镜像 TLS 失败 → 直连第 2 轮成功」，
   镜像兜底 + 2 轮重试逻辑真实起效，非过度设计。

### 4. deps.toml 完整性 — PASS

官方 release：`gh api repos/k2-fsa/sherpa-onnx/releases/tags/v1.13.8`（published 2026-09-10）。

- 资产 `sherpa-onnx-v1.13.8-win-x64-shared-MT-Release.tar.bz2` 存在，官方 size=24,805,859B。
- **独立下载复核**（`gh api … -H "Accept: application/octet-stream"`，与 fetch-deps 的 Invoke-WebRequest 完全不同通道）：
  下载体积 24,805,859B 与官方一致，sha256 = `6dffdc715a4465b989446a6105265d2cb345e7101591a17d35534b6758f6e8df`
  **与 deps.toml pin 完全一致**；从该独立归档提取三个 DLL：
  - `bin/onnxruntime.dll` = `7f66f939a881baf4f46a2216496798edf4a1429878b646d12674aa62f27d8a25` ✓ pin
  - `lib/sherpa-onnx-c-api.dll` = `300e0c88400903fc4cfc88be88a8ae587a24cc658653e0e1fd1c2b8dbbc68557` ✓ pin
  - `lib/sherpa-onnx-cxx-api.dll` = `02ca46060f65d7d50f3448e7e67e608b0547c9acb274cba2f5cc6936521352b0` ✓ pin
- 官方 release body **未发布** sha256 清单（上溯无可比对源），因此 pin 的可信度依据 =
  两条独立通道下载（fetch-deps/直连 与 gh api）哈希彼此一致 + 官方资产 size 一致——这是无官方清单时能达到的最强验证。
- onnxruntime 版本：release body 第 20 行 `Update onnxruntime to v1.28.2 by @csukuangfj (PR #3935)` ✓ pin 1.28.2；
  运行时 `sherpa.GetOnnxruntimeVersion()` 亦报 1.28.2。
- 「官方没有 mingw 变体 win-x64 shared 包」：资产全列表核查，win-x64 仅有 MT/MD（MSVC）变体，无 `-gnu`/mingw ✓
  （BUILD.md §2 该论断为真，链接走 Go 绑定模块自带 mingw 库的方案成立）。
- Go 绑定：上游 tag `refs/tags/v1.13.8` 存在（object `4797d195…`）；`go.sum` 钉死
  `sherpa-onnx-go v1.13.8`、`-windows v1.13.8`（另有 linux/macos 间接依赖）；全冷 clone 的模块下载经 sumdb 校验成功。

### 5. BUILD.md 一致性 — PASS（2 条措辞级备注见问题清单）

逐条与代码/实测对表：

| BUILD.md 宣称 | 实测 |
|---|---|
| §1 Go 1.27.1（go.mod `go 1.27`+`toolchain go1.27.1`） | ✓ go.mod；本机即 go1.27.1 |
| §1 mingw-w64 GCC 16.2.0 Rev3 | ✓ 本机 GCC 完全一致 |
| §1 PS 5.1+ / System32 tar.exe | ✓ `#Requires -Version 5.1`；fetch-deps 用 `%SystemRoot%\System32\tar.exe`；实测在 Windows PowerShell 5.1 下跑通 |
| §2 整包+3 DLL SHA256 pin | ✓ 四个哈希全部独立复核一致（见第 4 项） |
| §3 一键命令 `-Env dev\|prod` | ✓ 参数 ValidateSet 一致，prod 实测 |
| §3 CC 解析顺序 `$env:CC`→PATH→`MINGW64_ROOT`→`C:\msys64` | ✓ build.ps1 L48-57 逐行一致 |
| §3 GOPROXY 未设则默认 goproxy.cn | ✓ build.ps1 L119 |
| §3 产物与 SHA256SUMS（LF，GNU sha256sum -c 可验） | ✓ 实测 4× OK |
| §3 缓存命中约 1 秒（重哈希+manifest 核对） | ✓ 二次运行 4.4s（含 go build）；cache hit 行输出 |
| §4 doctor 各项 | ✓ 与 doctor.go 一致，FAIL 分支实测（第 1 项） |
| §5 坑 1 gcc 静默失败 → 自动前插 PATH | ✓ build.ps1 L61-65 存在 |
| §5 坑 2 `return ,$parts` | ✓ fetch-deps L48-50 存在 |
| §5 坑 3 `$PSScriptRoot` 在 param 默认值可能为空 | ✓ fetch-deps 在函数体解析默认值（build.ps1 在体内用 `$PSScriptRoot`，无此问题） |
| §5 坑 4 SHA256SUMS CRLF → WriteAllText+UTF8 no BOM+LF | ✓ build.ps1 L145-147；GNU sha256sum -c 实测通过 |
| §5 坑 5 import 必须子包 `…/sherpa_onnx` | ✓ main.go L16 正是子包路径 |
| §5 坑 6 TLS 失败常态 → 直连+镜像×2 轮 | ✓ 代码一致；验收中实际复现一次（第 3 项） |
| §6 MSVC 回退（条件性，`CC=cl` 为解析顺序第一） | ✓ 文档与代码一致，未越权启用 |
| §7 CI 备注引用 builder.Dockerfile、票 02 spike | ✓ 文件存在且注释一致 |

6 条踩坑全部有真实代码防护，且坑 6 在验收中被网络现场复现。

### 6. D22 禁令 — PASS

- **代码零 emoji**：对全部 `.go/.ps1/.toml/.mod/.sum/.Dockerfile` 扫描 `U+2190–U+2BFF / U+1F300–U+1FAFF / U+FE0F`：**零命中**。
  （`→`/`①`/`≤` 等仅出现在 markdown 散文里——docs 与票据既有行文约定，非代码，非本票引入。）
- **无明文密钥**：唯一的 key 形态常量是 `buildinfo.MinisignPublicKey = "PLACEHOLDER-C29-MINISIGN-PUBLIC-KEY"`，
  注释明确「真实 keypair 离线生成、密钥材料永不进仓库」——占位符合规；正则扫描无 API key/secret/token 明文。
- **goroutine 有名 + recover**：`cmd/`、`internal/` 无任何 `go func`/goroutine——本票无并发面，条款空转满足（后续票引入并发时需重新过闸）。
- **未声明依赖**：go.mod 直接依赖 3 项（sherpa-onnx-go / go-toml v2 / x/sys）全部真实使用；
  脚本仅用 PowerShell 内建 + System32 tar.exe；`ghfast.top` 只是**传输层**兜底前缀，哈希权威始终是仓库内 deps.toml
  （D22「禁从镜像站获取哈希」针对的是 C29 模型分发签名清单，本票不涉及，无违规）。
- **禁修改 D1–D46/C1–C31**：本票仅新增文件与 BUILD.md；PLAN/SPEC 未被改动（git log 核实）。

### 7. 票据对照（.scratch/wisp/issues/01-build-chain.md）— PASS

| 验收标准 | 裁决 | 证据 |
|---|---|---|
| Fresh clone → build.ps1 → 可运行 exe | PASS | 第 2 项（13.9s / 全冷 56.5s，doctor PASS） |
| BUILD.md 定稿冻结（版本、pin、命令、实际踩到的故障排查） | PASS | 头部 FROZEN 标记；6 条真实踩坑；第 5 项对表 |
| deps.toml SHA256 篡改一个字节 → 大声失败 | PASS | 第 3 项（两改两 FATAL） |
| wisp doctor 打印版本 + PASS/FAIL | PASS | 第 1、2 项（含 FAIL 分支实测） |
| BUILD.md 提交并推送双远程 | PASS | origin 与 cnb 的 dev 均 = `29c8e35`（含 ac96b4b BUILD.md 冻结提交）；origin 首查 TLS 失败重试后一致 |

- Progress log：4 条 entry 的内容全部与我的独立复现吻合，无虚报。
- Status 仍为 `in-progress` ✓（符合「验收前不置 done」的流程）。
- 两条格式级 MINOR 见问题清单（日志顺序与 Last update 字段）。

### 8. builder.Dockerfile 裁定 — PASS（调整合理）

- 内容与 SPEC-11 §3.1 逐行同构（apt 包集、ENV 四件套完全一致），唯一偏差：`FROM golang:1.24-bookworm` → `golang:1.27-bookworm`。
- **裁定：必要且正确**。go.mod 钉死 `go 1.27` + `toolchain go1.27.1`，1.24 基镜像要么构建被拒、要么依赖
  `GOTOOLCHAIN=auto` 现场下载工具链（不确定、违背「工具链钉死」的 §14.9 初衷）；基镜像直接对齐 go.mod pin 才是确定性的。
  `golang:1.27-bookworm` 标签已在 Docker Hub 核实存在（HTTP 200）。
- SPEC §3.1 的「诚实标注」要求被继承：Dockerfile 头注写明交叉链接可行性归票 02 spike、原生 Windows 是主路径 ✓；
  BUILD.md §7 同口径 ✓。SPEC 文本本身未回改（D22 禁随手改 SPEC，偏差已在票据 Progress log 记录）。

## 问题清单

无 BLOCKER、无 MAJOR。3 条 MINOR（均不阻塞）：

1. **MINOR（票据卫生）**：`01-build-chain.md` Progress log 头声明 `append-only, newest last`，
   但最新一条（`06:32:38Z` orchestrator 加 Dockerfile）被插在**最前**；且头部 `Last update:` 停在
   `05:33:59Z`（实际最新 entry 为 06:32:38Z）。修复：把最新条目移到末尾、更新 Last update 字段。
2. **MINOR（BUILD.md 措辞）**：§4「导入表里有 sherpa-onnx-c-api.dll 与 onnxruntime.dll」——
   wisp.exe 自身导入表只含 `sherpa-onnx-c-api.dll`，`onnxruntime.dll` 出现在其依赖 c-api DLL 的导入表里。
   加载器级结论不变（缺任一个都起不来，实测 exit 127），但措辞宜改为「导入链覆盖两者」。
3. **MINOR（Dockerfile 标注）**：基镜像 `1.24→1.27` 相对 SPEC-11 §3.1 文本的偏差，未在 Dockerfile/BUILD.md
   显式写一句「SPEC 文本为 1.24，已按 go.mod pin 调整为 1.27」；现在只散落在票据 Progress log 里。
   补一行注释可避免后人「为何与 SPEC 不一致」的困惑。

## 最终裁决

所有 8 项验收全部通过：干净构建可独立复现（13.9s/全冷 56.5s）、doctor 检查真实且 FAIL 分支为活、
三种篡改场景全部被拦截、四个 SHA256 pin 经独立通道复核一致、BUILD.md 与脚本逐条一致且 6 条踩坑均有真实防护、
D22 四类禁令零违规、票据 5 条验收标准全过、Dockerfile 版本调整合理。

**VERDICT: PASS**（3 条 MINOR 为卫生级建议，可随后续票据顺手修复，不构成本票返工项）
