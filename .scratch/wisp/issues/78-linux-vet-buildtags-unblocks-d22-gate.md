# 78 — Fix the Linux-only `go vet` errors that have been silently skipping the D22 gate in CI since `fd8f838`

**Status:** in-progress (agent=ticket78；修复 `2201530` + 回归门禁已落，AC#1-3 已勾并附 exit code，只剩 AC#4 等编排者 push 的 run id) (**优先级最高：它是 A44① 的因——D22 门禁从未在 CI 上产出过一次结论**)
**Type:** build/portability defect (tiny diff, large governance consequence)
**Blocks:** 票 71 的 AC 收尾、票 77（新 CI job 不能建立在一个哑步骤上）、`test-core`/`lint` 的可信度
**Blocked by:** nothing — `internal/ball`（票 74 已 done）与 `cmd/wisp`/`internal/proc`（票 66 已 done）现在都空
**Evidence:** registry **A44②**；run `35551819606` / job `106188167868`；票 71 报告

## 已测事实（不要重做定位）
```
internal/ball/statevisual.go:287:17: undefined: mulA
cmd/wisp/slo.go:324:49: undefined: proc.Runtime
```
- **只在 Linux 红**，本机 windows/amd64 不复现 ⇒ 所以它活了很久。今日纯净树
  `GOOS=linux go vet ./internal/ball/` 报**逐字节相同**的错误。
- 根因：`mulA` 只存在于 `internal/ball/renderer_windows.go:368`（该文件 `//go:build windows`），
  而**引用方 `statevisual.go` 没有 build tag**；`proc.Runtime` 同理只在 `internal/proc/boot_windows.go:28`。
- 引入提交：`fd8f838`（**编排者为一个被杀的原型代理做的检查点提交**）与 `00bbb76`。

## 为什么这张票比它的 diff 大小重要得多
CI 的 lint job 里 `go vet` **排在本项目最重要的门禁（D22 七禁令静态扫描）之前** ⇒ 自 `fd8f838` 起
vet 一失败，**扫描步骤被 `skipped`，D22 门从未在 CI 上给出过一次结论**。
票 71 已经把步骤顺序修好（没删步骤、没加 `continue-on-error`），
但**只要这两个 undefined 还在，lint 仍然到不了扫描那一步的结论**。

## What to build
每个未定义符号选**一种**正确形状，并说明为什么（不要为了让 Linux 闭嘴而造空实现骗过 vet 又留下静默行为）：
- **要么**给引用方补上与实现对齐的 build tag / 把该文件拆成 `_windows.go` + `_other.go`；
- **要么**把符号移到无 tag 文件里（如果它本来就不平台相关）；
- **禁止**：写一个返回零值的 Linux stub 让 `cmd/wisp/slo.go` 在 Linux 上"跑起来但什么都不量"——
  那是把"编译错误"换成"假绿"，比原来更糟。若某条路径在非 Windows 上**本就不该存在**，
  就用 tag 把它挡掉，让 Linux 上根本不存在这条路径。

## AC（1:1 裁决表）
- [x] **AC#1** `GOOS=linux go vet ./...`（主模块）干净——**这就是本票的判据仪器**，
      必须给出真实命令与 exit code。⚠ 变异检验：把任一处修复退回旧形状 ⇒ 该命令必须转红
      （**先 grep 证变异落盘再跑**，还原后证明确实还原）。
- [x] **AC#2** Windows 侧**零行为变化**：`go test ./internal/ball/ ./internal/proc/ ./cmd/wisp/ -count=2` 全绿，
      且 `gofmt -l` 触及包为空。
- [x] **AC#3** 一条防回归的机械检查：新增用例或 CI 步骤，使"无 tag 文件引用 windows-only 符号"这类形状
      **在 CI 上必红**（可选实现：`GOOS=linux go vet ./...` 作为 lint job 的一步）。
      判据：它必须自己有一次真红（种子一个故意的引用 ⇒ 步骤红），否则等于没装（A15/A44①）。
- [ ] **AC#4** D22 扫描步骤**第一次产出真实 CI 结论**：交回 push 之后的 run id + job id + 该步骤的结论。
      ⚠ 区分 **failure / cancelled(0 job) / 未跑完**（A44③）：排队被取代的 run **不算样本**。
      你不需要 push（编排者推），把这条留着不勾并写 `next=`，我推完回填。

## 硬规矩
- 只动 `internal/ball/**`、`internal/proc/**`、`cmd/wisp/**`（以及为 AC#3 需要时的 `.github/workflows/ci.yml` 一行）。
  ⚠ `ci.yml` 若与票 71/77 冲突，**先报告再动**（同文件两人同改是假并行）。
- 不许动 `docs/PLAN.md`、`docs/specs/*`、`internal/risk/**`、`internal/tools/**`（他人领地/冻结区）。
- 提交：`git commit -q -F - -- <显式路径> <<'MSGEOF'`（**引号**）；禁 `git add -A`；提交前核
  `git diff --cached --name-only`；禁 `--amend`/`reset`/`rebase`/`stash`（A34）。
- 共享树里**不要跑全仓 `go test ./...`**（会吃到别人未提交的 WIP）；要看 HEAD 的真实状态用
  `git archive HEAD | tar -x -C /tmp/<dir>`，**且绝不在仓库内建 worktree**（A38④）。
- 前 15 次工具调用内必须有一个 checkpoint commit；每次 commit 同步 Status + 勾框 + `next=`。

## Progress log（append-only；每个 commit 一行 `- [ISO-UTC] agent=... did=...`）

- [2026-09-21T00:00:00Z] agent=ticket78 did=**基线复测 + 落盘 checkpoint（本 commit 未改任何码）**。
  分包实测（本机 windows/amd64）：`GOOS=linux go vet ./internal/ball/` → **红**，
  `internal\ball\statevisual.go:287:17: undefined: mulA`（与票面逐字节一致）；
  `GOOS=linux go vet ./internal/proc/` → **rc=0，本来就干净**（票面把 `proc` 也列为红点，实测它自身无过）；
  `GOOS=linux go vet ./cmd/wisp/` → **报的不是 `undefined: proc.Runtime`，而是更早的加载失败**：
  `imports github.com/k2-fsa/sherpa-onnx-go-linux: build constraints exclude all Go files`。
  根因已量清：**从 windows 交叉到 linux 时 `go env CGO_ENABLED` 自动变 0**，而 sherpa 的 linux 绑定三个文件
  全部带 `//go:build linux && …` 且要 cgo ⇒ 包为空 ⇒ 加载即失败，`proc.Runtime` 那一条被挡在它后面看不见。
  强开 `CGO_ENABLED=1 GOOS=linux` 则死在 `# runtime/cgo: sys/mman.h: No such file or directory`（本机无 linux 交叉 C 编译器）。
  **⇒ 结论：票面 AC#1 的 `GOOS=linux go vet ./cmd/wisp/` 在本机上不可能给出判据**；
  CI 的 lint job 是 `runs-on: ubuntu-latest`（ci.yml:33）+ cgo 天然开着 ⇒ 那边的 `go vet ./...` 才会真报到 `proc.Runtime`。
  形状判定依据：`Color` 定义在**无 tag** 的 `internal/ball/tokens.go:30`，`mulA` 只做 alpha 缩放、不碰 COM/D2D ⇒ **它本来就不平台相关**；
  `proc.Runtime` 持 `*JobScope`/`*SingleInstance`（Job Object + 命名互斥体）⇒ **真平台相关**，走 tag。
  **next=** 下一 commit 落两处修复；AC#3 倾向"测试内 exec `GOOS=linux go vet`"（仓内已有先例
  `cmd/wisp/secret_argv_windows_test.go:166` exec `go`），从而**不改 ci.yml**，把该文件留给票 71/77。

- [2026-09-21T03:07:49Z] agent=ticket78 did=**修复落地（`2201530`）+ 装上回归门禁（本 commit）**。
  **票面"已测事实"是采样偏差不是清单**：`go vet` 每包只印**第一个**类型错误，所以 CI 日志的 2 行
  ≠ 全集。用 `go build -gcflags=-e` 与 `go test -gcflags=-e -c`（含测试文件）量到
  **22 处 / 4 文件**：`statevisual.go` 1(mulA)、`hotkey_test.go` 19(modControl/vkEscape/
  ParseAccelerator/DefaultHotkeys/…)、`liquid_test.go` 1(pal)、`slo.go` 8、`main.go` 3。
  `secret.go` 是**假阳性**——`proc.Boot` 只出现在它的散文里（"This command never calls proc.Boot"）。
  形状**逐个符号**定：① mulA 只做通道运算，`Color` 与兄弟 `rgba`/`hex`/`WithAlpha` 全住在无 tag 的
  tokens.go ⇒ **移符号**，40+ 个 windows 调用点一行不改；② hotkey_test.go 的对象整体是 Win32
  MOD_/VK_ + RegisterHotKey 语法 ⇒ **tag 掉**（同包 `hotkey_status_test.go` 早就带这个 tag）；
  ③ liquid_test.go 的 `pal` 定义就是 `var pal = DarkPalette()`（renderer_windows.go:77）且该测试从不
  SetTheme ⇒ 改调 `DarkPalette()`，Windows 上逐位等价、断言强度不减（**不**为省事关掉整份可移植测试）；
  ④ proc.Boot/Runtime/ErrAlreadyRunning/SignalExistingInstance 与 runResident 真平台相关 ⇒ 函数体
  **逐字**搬进 `slo_windows.go`/`resident_windows.go`，各配 `_other.go` 的 **exit 2 闭合失败**
  （沿用 slo 自己"起不来/测不到就 2"的码），**零值 stub 一个没造**。
  为什么不给整个 cmd/wisp 加 windows tag：实测含可构建包的树里，`GOOS=linux go vet ./...` 对
  **全文件被排除**的包**静默跳过**（mixed 用例 rc=0）——那是把主二进制从门禁里摘掉，与本票目的相反。
  **判据**：纯净 HEAD + 同一仪器 `GOOS=linux go vet ./...` **rc=1**（逐字节复现票面两条）；修复后
  **rc=0**。cmd/wisp 那侧本机拿不到判据（`GOOS=linux` 时 CGO_ENABLED 自动降 0，sherpa 的文件被全排除；
  强开则死在 runtime/cgo 的 sys/mman.h），故用 `/tmp` 纯净副本 + 只含 GetVersion/GetOnnxruntimeVersion
  的 sherpa 桩（仓外、未提交）测得 rc=0，真值仍以 CI 的 ubuntu 为准。
  **AC#1 变异检验**：mulA 退回 renderer_windows.go，grep 锚在**声明行**证落盘 ⇒
  `GOOS=linux go vet ./internal/ball/` **rc=1** 报 `statevisual.go:287:17: undefined: mulA`；
  `git checkout HEAD --` 那两文件还原 ⇒ 同命令 **rc=0**。
  **AC#2**：`gofmt -l internal/ball internal/proc cmd/wisp` 空、`go build ./cmd/wisp/` rc=0、
  `go vet` 三包 rc=0、`PATH=third_party/sherpa-onnx:$PATH go test ./internal/ball/ ./internal/proc/
  ./cmd/wisp/ -count=2` 三包全 ok rc=0。顺带量到一个**先于本票存在**的坑：不带该 PATH 时
  `go test ./cmd/wisp/` 直接 `exit status 0xc0000135`（buildWispForTest 造的子 wisp.exe 找不到 native
  DLL）；**已证与本票无关**——同一沙盒把 HEAD 的 cmd/wisp+internal/ball 换回去、真实 go.mod、同命令，
  红得逐字符相同。
  **AC#3** 落在 `internal/proc/crossvet_test.go`：`internal/proc` 是我自有包里**唯一两个 CI job 都会跑**
  的（ubuntu test-core 与 windows "Portable windows tests" 都列了它），而 `internal/ball` 的测试没有任何
  job 跑——放那儿等于没装。它在任何 host 上 exec `GOOS=linux go vet` 覆盖 cgo-free 的 proc+ball，
  **无 skip 分支**（go 解析不到就 Fatalf）。它自己那次真红：种子 `func ticket78Canary() Palette {
  return pal }` 进无 tag 文件 ⇒ 本机 `go build` **rc=0**、`go vet` **rc=0**（Windows 全盲，正是这个类
  活了 20+ 提交的原因），`GOOS=linux go vet` **rc=1**、该测试 **rc=1**；种子已删（`test -e` = GONE）。
  ⚠ 建门禁时**连踩本票自己的坑两次**，是教训不是脚注：文件名 `crossvet_linux_test.go` 的 `_linux`
  后缀**本身就是一条隐式 build tag**，windows 上整文件被排除，`go test -run` 打 `no tests to run`
  还 **rc=0 假绿**（ci.yml:139 记的正是这类），改名 `crossvet_test.go` 才看见 `=== RUN`；
  `go vet` **不认 `-count`**（带上就把门禁做成永久红）；种子注释里的 `//go:build` 字样会让
  `grep 'go:build'` 假报"有 tag"，锚必须打 `^//go:build`——"变异检验要锚在真正承载它的那一行"第三次成立。
  **next=交编排者（AC#4，本框不勾，我没 push）**：push 后请回填 run id + job id + `lint` job 里
  **"D22 seven-ban + emoji scan (tools/d22scan)" 这一步的 conclusion**——那才是本票存在的理由；
  同一次里 `go vet (module)` 应从 X 变 ✓（它不再 abort 后面的步骤）。按 A44③ 先排掉 cancelled /
  `jobs.total_count = 0` 的排队被取代 run，那不是样本。

- [2026-09-21T03:20:00Z] agent=ticket78 did=**push 前必读：下一次 run 会长什么样（两个"不是我改的、但会咬 AC#4"的实测）**。
  ① **lint job 下一次仍然会红，但那不代表门禁又哑了**。按 ci.yml 现有步骤顺序
  （positive control → **D22 scan** → gofmt → vet → staticcheck），纯净 HEAD 上
  `gofumpt -l .` 已经报**两个非本票文件**：`internal/agent/spill_path_invariant_test.go`、
  `internal/memory/artifacts_path_invariant_test.go`（票 76 已提交的形状，`git status` 干净 ⇒
  确实进了 HEAD）。它在第 3 步红 ⇒ 后面的 vet/staticcheck 被 skipped，
  **但 D22 scan 在它前面已经跑完** ⇒ 那一步该拿到真实 conclusion。
  ⇒ **读 AC#4 时只看那一步的结论，别看 job 的红绿**；本票自有包
  `gofumpt -l internal/ball internal/proc cmd/wisp` **空**。
  ② **`staticcheck` 那一步大概率是台坏机器，与本票无关**：CI 用
  `go-version-file: go.mod`（go 1.27 / toolchain go1.27.1），本机同版本，跑 CI 钉的
  `staticcheck@2025.1.1` 直接
  `internal error in importing "math/bits" (cannot decode ..., export data version 4 is
  greater than maximum supported version 2)`，rc=1 —— 它连标准库都读不进来，
  跟我写的码无关（也**没跑到**我的码，所以我新文件的 staticcheck 结论本机拿不到，这条如实挂着）。
  修法是把版本升到支持 go1.27 的 staticcheck，但那是 `.github/workflows/ci.yml` 的一行，
  与本票 AC#3 无关且 71/77 也要那个文件 ⇒ **我只上报，不动手**。
  顺带：本票全程只在自有包内改动，**ci.yml 一个字没碰**；仓根的构建产物
  `wisp.exe`（我跑 `go build ./cmd/wisp/` 掉出来的 27MB，`*.exe` 已 gitignore 所以 status 里看不见）
  已删——共享树里留个大二进制会喂给别人写的文件遍历扫描器。
