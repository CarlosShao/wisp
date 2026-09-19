# BUILD.md - 构建流程（S0 冻结）

> 状态：**FROZEN（S0 ticket 01 定稿）**。本文件与 `deps.toml`、`scripts/build.ps1`、
> `scripts/fetch-deps.ps1` 是一体的：改其中任何一个，必须同步改其余并重新走一遍
> 「干净 clone 验证」。工具链版本升级属于决策，不是随手 `go get`。

## 1. 工具链（钉死）

| 工具 | 版本 | 说明 |
|---|---|---|
| Windows | 10/11 x64 | 自用期唯一目标平台（windows/amd64） |
| Go | **1.27.1** | `go.mod` 同时钉了 `go 1.27` + `toolchain go1.27.1`，本地高于此版本也按 toolchain 指令走 |
| C 工具链（cgo） | **MSYS2 mingw-w64 GCC 16.2.0 (Rev3, x86_64-win32-seh)** | SPEC-11 §2.1 定案的唯一支持路径 |
| PowerShell | 5.1+（Windows 自带） | 脚本兼容 5.1，不要求 PS7 |
| MSYS2 tar | `C:\Windows\System32\tar.exe`（libarchive bsdtar） | 解压 .tar.bz2，Windows 10+ 自带 |
| git | 任意近期版本 | build.ps1 用它取 commit hash；缺失时 commit=nogit，不阻塞构建 |

分发版参考：mingw-w64 走 MSYS2 的 mingw64 子系统（`pacman -S mingw-w64-x86_64-gcc`）或
winlibs 等价发行版；**不要**用 MSYS2 的 `usr\bin\gcc`（那是 Cygwin 运行时）。

## 2. deps.toml 钉了什么（版本 + SHA256）

| 依赖 | 版本 | 来源 | SHA256 | 许可 |
|---|---|---|---|---|
| sherpa-onnx（Windows shared release, MT-Release） | v1.13.8 | `https://github.com/k2-fsa/sherpa-onnx/releases/download/v1.13.8/sherpa-onnx-v1.13.8-win-x64-shared-MT-Release.tar.bz2`（镜像前缀 `https://ghfast.top/`） | `6dffdc715a4465b989446a6105265d2cb345e7101591a17d35534b6758f6e8df`（整包）；`deps.toml` 另钉每个 DLL 的独立哈希 | Apache-2.0 |
| onnxruntime.dll（随上包分发） | 1.28.2 | 同上包内 `bin/onnxruntime.dll` | `7f66f939a881baf4f46a2216496798edf4a1429878b646d12674aa62f27d8a25` | MIT |
| sherpa-onnx Go 绑定 | `github.com/k2-fsa/sherpa-onnx-go v1.13.8`（含 `sherpa-onnx-go-windows v1.13.8`） | go.mod + goproxy（go.sum 校验） | go.sum 钉 module 哈希 | Apache-2.0 |

要点：
- **链接时**用的是 Go 绑定模块自带的 mingw 预编译库（`sherpa-onnx-go-windows` 模块内
  `lib/x86_64-pc-windows-gnu/`，cgo LDFLAGS 由模块自己注入，零手工配置）。
- **运行时分发**的是官方 release 的 MSVC（MT=静态 CRT）构建 DLL，由 fetch-deps 下到
  `third_party/sherpa-onnx/`、build.ps1 复制到 exe 同目录。两套同为 sherpa-onnx v1.13.8
  同源码、同 onnxruntime 1.28.2，C ABI 一致；`wisp doctor` 在运行时以真实调用验证这一点。
- 官方 v1.13.8 **没有**发布 mingw 变体的 win-x64-shared 包（只有 MD/MT 的 MSVC 变体），
  所以采用上面的「绑定模块内 mingw 库链接 + 官方 release DLL 分发」组合。
- 网络现实：github.com release 资产从 CN 直连经常 TLS 失败。fetch-deps 直连失败自动换
  `https://ghfast.top/` 镜像前缀重试（直连/镜像各 2 轮）。SHA256 钉死使镜像可信。

## 3. 一键构建（干净 clone 可跑通）

```powershell
# 前置：Go 与 mingw64\bin 都在 PATH（或设 MINGW64_ROOT 指向 mingw64\bin 目录）
git clone <repo> && cd wisp
powershell -ExecutionPolicy Bypass -File scripts\build.ps1 -Env dev     # 或 -Env prod
```

build.ps1 流程（SPEC-11 §2.2）：
1. `fetch-deps.ps1`：校验/补齐 `third_party/sherpa-onnx/`（缓存命中约 1 秒：重哈希缓存
   文件并核对 manifest 与 deps.toml 的一致性）
2. 前端：跳过（S5 之前无 frontend）
3. `go build`：`CGO_ENABLED=1`、`CC=<mingw gcc>`、`-trimpath`，ldflags 注入版本/commit/
   构建时间/`DefaultEnv`/deps.toml 里的两个原生库版本
4. 产物：`build\wisp.exe` + `build\{onnxruntime.dll, sherpa-onnx-c-api.dll,
   sherpa-onnx-cxx-api.dll}`（同目录规则 SPEC-11 §7.1）
5. `build\SHA256SUMS`（LF 行尾，GNU `sha256sum -c` 可直接验证）
6. 冒烟：运行 `build\wisp.exe doctor`，FAIL 则整个脚本 FAIL

工具解析顺序：`go` 必须在 PATH；`CC` 依次找 `$env:CC` → PATH 里的 `gcc.exe` →
`$env:MINGW64_ROOT\gcc.exe` → `C:\msys64\mingw64\bin\gcc.exe`。`GOPROXY` 未设置时脚本
默认 `https://goproxy.cn,direct`（CN 网络；已设置则尊重现值）。

## 4. 验证构建产物

```powershell
cd build
.\wisp.exe                # 无参 = 常驻进程（票 03）：打印版本 → 引导运行时骨架
                          #   （Job Object + 单实例 + goroutine 注册表自检）→ 空事件循环；
                          #   Ctrl+C / 结束信号触发 D38(e) 十步关停后退出。同会话二次启动
                          #   会激活已有实例并退出。悬浮球窗口在票 07
.\wisp.exe run "任务文本"   # CLI 占位：回显任务文本（真实 agent 循环在票 10）
.\wisp.exe doctor         # 自检：工具链/DLL 同目录/版本匹配 deps.toml，输出 PASS/FAIL
sha256sum -c SHA256SUMS   # 或 PowerShell: Get-FileHash 对表
```

`wisp doctor` 各项：
- sherpa-onnx C API 运行时版本 == 构建时从 deps.toml 烘焙进二进制的 pin
- `onnxruntime.dll` FileVersion == pin（读 PE VS_FIXEDFILEINFO，容忍 `.0` 尾差）
- `sherpa-onnx-c-api.dll` / `sherpa-onnx-cxx-api.dll` 存在于 exe 同目录
- deps.toml 在 exe 旁/上级目录/工作目录能找到时，额外核对三个 pin 与二进制一致
- C29 minisign 公钥：占位常量（真实密钥在 C29 落地，密钥材料永不进仓库）
- 数据目录可写（WISP_ENV / portable.txt 规则，SPEC-03 §5.2）

注意：本二进制是**直接链接** sherpa C API 的（导入表里有 `sherpa-onnx-c-api.dll` 与
`onnxruntime.dll`；注意 onnxruntime.dll 是 sherpa-onnx-c-api.dll 的**间接**导入，并非 wisp.exe 的直接导入，但加载器解析规则相同），所以 DLL 不在 exe 旁时进程根本起不来（Windows 报 0xC0000135）——
「doctor 打印 DLL FAIL」这一格在链接语义上由加载器兜底，doctor 打印的 PASS 即加载器
级证明。windowsgui 子系统切换（`-H=windowsgui`）推迟到票 07；届时如需「缺 DLL 仍能
弹出友好错误」，应改为运行时 LoadLibrary 包装，这是那票的设计题。

AttachConsole 通路已按本票代码预先验证：用 `go build -ldflags "-H=windowsgui"` 临时
构建后，`wisp.exe version` 在调用方控制台正常输出（attach+重绑句柄路径生效），且
`wisp.exe version > out.txt` 重定向不受影响（已有合法 stdout 时跳过 attach）。

## 5. 本票真实踩到的坑（排查手册）

1. **mingw gcc「静默失败」**：`gcc.exe` 不在 PATH 时，用全路径调用也会以 exit 1 退出且
   **没有任何 stderr**（其依赖的 mingw64 运行时 DLL/组件按 PATH 解析失败）。cgo 的症状是
   `runtime/cgo: cgo.exe: exit status 2`，同样无输出。解法：把 `mingw64\bin` 加进 PATH
   （或设 `MINGW64_ROOT`）；build.ps1 解析出 CC 后会自动把它的目录前插到 PATH。
2. **PowerShell 函数返回单元素集合被解包**：`Split-TomlKey` 返回 `List[string]`，只有
   一个元素时 PS 把它解包成标量字符串，调用方的 `[0]` 取到的是**第一个字符**
   （`source` → `s`）。解法：`return ,$parts`（前导逗号阻止解包）。
3. **`$PSScriptRoot` 在 param 默认值里可能为空**：`powershell -File script.ps1` 时
   param 默认表达式里 `$PSScriptRoot` 为空串，`Split-Path` 直接抛错。解法：param 只收
   原始参数，脚本体里再解析默认值。
4. **SHA256SUMS 行尾**：`Set-Content` 写出 CRLF，GNU `sha256sum -c` 报
   `'wisp.exe'\r: No such file or directory`。解法：`[IO.File]::WriteAllText` + UTF8
   (no BOM) + LF 连接。
5. **`github.com/k2-fsa/sherpa-onnx-go` 根目录没有 Go 包**：import 必须写子包
   `github.com/k2-fsa/sherpa-onnx-go/sherpa_onnx`，否则 `go mod tidy` 报
   "module ...@latest found, but does not contain package"。
6. **release 资产 TLS 失败是常态**：本票构建中「直连失败 → 镜像失败 → 直连重试成功」
   实际发生过。fetch-deps 的直连+镜像+多轮重试不是过度设计，别删。

## 6. 备用路径（MSVC 回退条件）

**主路径 = mingw-w64**（本文件）。仅当未来 sherpa-onnx 版本的预编译库/Go 绑定实际无法
用 mingw 链接（新版本符号/CRT 不兼容，且官方不再出 `-gnu` 库）时，回退：

1. 安装 VS 2022 Build Tools（勾选 "Desktop development with C++"）
2. 在 "x64 Native Tools Command Prompt"（或先执行 `vcvars64.bat` 的环境）里跑
   `scripts/build.ps1`，并设 `$env:CC = "cl"`（脚本工具解析顺序第一位就是 `$env:CC`）
3. deps.toml 的 DLL 来源不变（MT-Release 本就是 MSVC 构建，回退后链接与分发同源，
   版本错配风险反而更小）

架构不变（SPEC-11 §2.1）：回退只改工具链，不改 deps.toml/脚本结构/doctor 逻辑。
若走回退，把「哪个赢了」更新到本节并解除冻结重新评审。

## 7. CI 备注（票 08 前不建 workflow）

CI runner 上等价于：装 Go 1.27.1 + MSYS2 mingw-w64 gcc（或直接用
`docker/builder.Dockerfile`，交叉链接 sherpa Windows 库的可行性由票 02 spike 验证，
Windows 原生构建始终是主路径）→ 跑同一份 `scripts/build.ps1` → `wisp doctor` 判绿。
缓存目录：Go module/build cache + `third_party/`。
