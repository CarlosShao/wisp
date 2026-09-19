# SPEC-11 · 构建、容器化与环境隔离

> 追溯：D17/D24/D41/§14.9（高危缺口）/C29/C30；环境隔离（WISP_ENV）定义见 SPEC-03 §5，
> 本 spec 是其构建/CI/容器侧实现。S0 的 done 判据直接引用本文 §2/§3/§6。

## 1. 容器化边界（先说清楚什么不能容器化）

Wisp 是 Windows 桌面常驻应用：Win32 分层窗口 + WASAPI 音频 + WebView2 + sherpa-onnx cgo。
**没有独立的「后端 API 服务」**——常驻主进程就是后端，LLM 是第三方云 API。因此：

| 对象 | 容器化 | 理由 |
|---|---|---|
| Wisp 主程序运行时 | ❌ | Win32 窗口/音频设备/悬浮球必须真实交互桌面会话；Windows 容器无交互桌面、无音频端点 |
| **Go 构建工具链** | ✅ builder 镜像 | 消灭 §14.9「agent 在构建链上无限试错」——工具链版本钉死，任何机器一键构建 |
| **前端构建** | ✅ Node 镜像 | 产物进 `embed.FS`；宿主机不需要装 Node |
| **测试依赖服务** | ✅ compose | mock LLM / mock 站点 / 模型镜像 / SSRF 靶机，CI 与本地同一套 |
| **CI 流水线** | ✅ 容器化 job | 见 §6 |
| **纯逻辑开发环境** | ✅ devcontainer（可选） | agent 循环/tools/risk/storage/前端可容器内开发+测试；GUI/音频/WebView 必须在 Windows 宿主 |

## 2. 构建工具链（§14.9 定案，S0 冻结进 `docs/BUILD.md`）

### 2.1 工具链选型【SPEC 决策】

- **Go**：pin 于 `go.mod`（`toolchain` 指令），版本写入 BUILD.md。
- **C 工具链（cgo）**：**MSYS2 mingw-w64 GCC（winlibs 发行版，版本钉死）**为唯一支持路径——
  免 Visual Studio 安装（数 GB）、CI 安装快、与 sherpa-onnx 官方 Windows 预编译库链接可行
  （S0 验证；若实测 MSVC 才能链接，回退 VS Build Tools + `CC=cl`，作为 BUILD.md 的备用条目，
  架构不变）。
- **目标架构**：自用期仅 `windows/amd64`；`windows/arm64` 与 `darwin/*` 归 S8（D41d）。
- **原生库来源**：**构建时下载 + SHA256 校验 + 本地缓存**（`scripts/fetch-deps.ps1`，版本与哈希
  钉死在 `deps.toml`），产物落 `third_party/`（git 忽略）。不 vendored 进 git（体积）；不要求
  预装（贡献者门槛）。CI 用 Actions 缓存目录。
  - `sherpa-onnx` C 库 + `onnxruntime` 动态库版本必须与编译时链接一致——版本错配的症状是
    运行时崩溃而非编译错误，故 `deps.toml` 同时是运行时自检的比对源（§7.2）。
  - **D47 追加（2026-09-19）**：`webrtc-audio-processing`（AEC3）为 P15-gated 的第三项原生依赖
    ——P15 通过后纳入 `deps.toml`（版本+SHA256）与 fetch-deps 流水线，DLL 同目录分发规则同上；
    P15 不通过则不引入（Path C barge-in 降级为按键打断，SPEC-01 白名单条目回收）。

### 2.2 构建顺序（一键流程）

```
scripts/build.ps1 [-Env dev|prod] :
  1. fetch-deps.ps1        # 校验/补齐 third_party/（有缓存则秒过）
  2. 前端产物：存在则跳过；--with-frontend 时走 docker/frontend.Dockerfile（§3.2）
  3. go build（cgo: CGO_ENABLED=1, CC=mingw32-gcc; embed assets/web 已就位）
  4. 产物：wisp.exe + onnxruntime.dll + sherpa-onnx c dll 同目录（§7.1）
  5. 输出 SHA256SUMS
```

- CLI 与 GUI 同一二进制：无参 = GUI（`-H=windowsgui`）；`wisp run` 子命令
  `AttachConsole(ATTACH_PARENT_PROCESS)` 输出。【SPEC】
- 构建期注入：版本号、commit、`buildinfo` 里的 C29 minisign 公钥。

## 3. Docker 资产（`docker/`）

### 3.1 `builder.Dockerfile`（Linux 交叉构建，CI 便利路径）

```dockerfile
FROM golang:1.24-bookworm
RUN apt-get update && apt-get install -y --no-install-recommends \
    gcc-mingw-w64-x86-64 g++-mingw-w64-x86-64 git ca-certificates && rm -rf /var/lib/apt/lists/*
ENV CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64
```

- ⚠ **诚实标注**：mingw 交叉链接 sherpa Windows 预编译库的可行性由 **S0 spike 验证**；
  不可行则交叉路径降级为「仅纯 Go 包的 lint/test」，发布产物一律走 Windows 原生构建（§6 CI 已
  按此设计，不押注交叉路径）。

### 3.2 `frontend.Dockerfile`（前端构建，S5 起需要）

```dockerfile
FROM node:20-alpine AS build
WORKDIR /src
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build          # 产物 /src/dist
FROM scratch AS export
COPY --from=build /src/dist /dist
# 用法：docker buildx build --output type=local,dest=assets/web -f docker/frontend.Dockerfile .
```

- 产物落 `assets/web/` 供 `//go:embed`；宿主机无需 Node；版本随 `package-lock.json` 钉死。
- S5 之前 frontend 目录可空，构建脚本跳过该步。

### 3.3 `compose.dev.yml` / `compose.test.yml`（测试依赖服务）

| 服务 | 镜像/构建 | 端口 | 用途 |
|---|---|---|---|
| `mock-llm` | `tools/mockllm`（本仓 Go 构建） | 18080 | 三协议 mock（§4） |
| `mock-web` | nginx:alpine + 注入载荷静态页 | 18082 | `web.fetch` 正常/注入载荷页（D30/F2 用例） |
| `mock-search` | tools/mockllm 子命令 | 18083 | `web.search` 结果页/API 桩 |
| `model-mirror` | nginx:alpine | 18081 | 正常/损坏/坏清单三档模型文件（C29 用例）；测试微型模型 |
| `ssrf-target` | nginx:alpine | 内网 only（expose 不映射） | 断言 `web.fetch` 私有网段拒绝且**零请求到达** |
| `webhook-receiver` | tools/mockllm 子命令 | 18084 | 外泄断言收据（收到即用例失败） |

- dev 与 test 的差异仅默认端口/数据卷与 `WISP_ENV` 注释示例；同一 compose 文件用 profile 切换。
- `scripts/dev.ps1` = `docker compose -f docker/compose.dev.yml up -d` + `go run ./cmd/wisp`
  （`WISP_ENV=dev`）。

## 4. tools/mockllm（本仓 Go 工具，非运行时依赖）

- 端点：`POST /v1/chat/completions`、`POST /v1/responses`、`POST /v1/messages`——支持
  `stream:true` SSE、tool_calls、usage 字段；**golden 模式**：按 `testdata/golden/<name>.sse`
  脚本回放（与单测回放共用文件格式）。
- **控制端点**（测试编排用）：`POST /__control/fail_next`（次数+error_class）、
  `POST /__control/latency`（首 token/分块延迟）、`POST /__control/truncate`（模拟流断开）、
  `POST /__control/reset`。
- 双形态：`go run ./tools/mockllm`（本地）或容器（compose）；不含任何 Wisp 运行时代码依赖。

## 5. devcontainer（可选，纯逻辑开发）

- `.devcontainer/devcontainer.json`：Go + Node features + docker-in-docker（起 compose mock）。
- 适用：`agent`/`llm`/`tools`/`risk`/`config`/`memory`/`observe` 的单测与集成测试、前端开发。
- **不适用**（明确写进说明）：`ball`/`panel`/`audio`/`speech`（需 Windows 桌面 + 真实设备）；
  SLO 门禁（需 Windows）。

## 6. CI 流水线（GitHub Actions【SPEC】）

| job | runner | 内容 | 门禁级别 |
|---|---|---|---|
| `lint` | ubuntu + builder 镜像 | gofumpt/vet/staticcheck + **D22 七禁令静态扫描**（grep/AST）+ 零 emoji 扫描（design/ 与 frontend/） | 阻塞 |
| `test-core` | ubuntu + compose | 纯 Go 逻辑单测/集成（agent/tools/risk/config/memory）+ 黄金流回放 + 契约测试 | 阻塞 |
| `test-windows` | windows-latest | PathResolver 真实 junction 用例 + cgo 构建冒烟 + 崩溃边界测试 | 阻塞 |
| `frontend` | ubuntu + node 镜像 | tsc/lint/build + 硬编码色值/emoji 扫描 | 阻塞（S5 起） |
| `slo-smoke` | windows-latest | `slo-check.ps1` 内存/句柄子集（无音频场景） | 阻塞 |
| `slo-full` | **self-hosted Windows**（用户机器或专用机） | 六态全量 + 回落 + 延迟 + CER | 合并前必须本地/自托管跑（脚本同份） |
| `asr-cer` | self-hosted + GPU 可选 | 双 wav 集 CER 门禁 | nightly + 模型变更触发 |
| `release` | windows-latest | goreleaser：zip 产物 + SHA256SUMS +（secrets 存在时）SignPath 签名 +（开源后）winget 清单 | tag 触发 |

- 缓存：Go modules/build cache、fetch-deps 的 `third_party` 缓存、node_modules。
- PR 门禁 = lint + test-core + test-windows + slo-smoke；**任何 job 不得配置为可跳过**（D22）。

## 7. 打包、分发与更新（D41）

### 7.1 安装布局

```
%LOCALAPPDATA%\Programs\Wisp\        # per-user 安装，免 UAC
    wisp.exe
    onnxruntime.dll                  # 必须与 exe 同目录：不依赖 PATH（版本错配）+ 防 DLL 劫持
    sherpa-onnx-*.dll                # 同上
    wisp-updater.exe                 # 独立小工具（不 embed 前端资源）
    resources\                       # 前端产物（亦可 embed.FS）
    LICENSE
```

自用期分发 = GitHub Releases zip + SHA256 校验；winget/Scoop 归 S8（D17/D23）。
启动自检（任一失败给明确错误，不得静默降级）：`onnxruntime.dll` 存在且版本匹配 `deps.toml`
→ 自身签名校验（若已签）→ C29 公钥可用 → 数据目录可写。

### 7.2 更新运行中的 exe（D41(b)）

```
下载新版 → staging\ → minisign 验签（C29 同一套密钥）→ 写 update-pending.json
→ 启动 wisp-updater.exe → 主进程按 D38 关停顺序优雅退出
→ updater 等主进程退出（≤5s，超时关 Job Object 强杀）
→ 备份旧版到 backup\<version>\ → 原子替换 → 重启主进程 → updater 退出
```

**回滚（N-1）**：新版启动 60s 内崩溃 ≥2 次 → 自动从 `backup\` 恢复 + 写日志 + 悬浮球提示；
只保留一个备份版本。

### 7.3 两套签名机制，缺一不可（D33/F3 + D41e）

- **CA 代码签名**（SignPath Foundation，免费开源项目通道）：解决 SmartScreen；审批不过 →
  未签名 + 文档化绕过，架构不变（D17）。**secrets-gated**：无 secrets 跳过而不是失败。
- **minisign**：更新包与模型的完整性/来源认证（免费、无 CA）；公钥硬编码 `buildinfo`。
  minisign 不依赖 SignPath，**可先用**（P8）。

### 7.4 卸载（D41c）

删 `%LOCALAPPDATA%\Programs\Wisp\` · 删自启注册项 · Job 随进程消失 · **询问**是否删
`%APPDATA%\wisp\`（默认保留并告知路径；提示中必须列出已下载模型占用大小）。

## 8. 测试决策

- **构建可复现测试**：干净 clone + `build.ps1` 在 CI（windows runner）全绿 = S0 done①；
  BUILD.md 流程与脚本一致性由脚本内 `--check` 模式核对。
- mockllm 契约：golden 回放文件在单测与 compose mock 两条路径下行为一致（同一份测试断言）。
- 更新流程测试：staging→验签→替换→回滚全链路（本地脚本化预演，真实替换文件）。
- 环境分叉在 CI 的可见性：`test-windows`/`test-core` job 显式断言 `WISP_ENV=test` 生效
  （数据目录在临时路径、互斥未注册）。

## 9. 不做什么

- 不容器化 Wisp 运行时（§1）；不做 `wisp run` 的容器化镜像（工具全是 Windows 本机语义，
  D46 CLI 包装在本机执行）。
- 不做 Docker Desktop 硬依赖：无 Docker 时构建退化为「宿主机手工工具链」（BUILD.md 附录记录），
  测试依赖服务退化为 `go run ./tools/mockllm` 直启——容器是默认，不是唯一。
- 不做多版本回滚链（只 N-1）；不做增量更新（全量 zip，体积 ~40–80MB 可接受）。
