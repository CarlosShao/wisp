# 455 — 写腿收件：让"裸跑 cmd/wisp"不可能（绑定自检＋sha256 门＋注释如实＋一次载具形整包）

**Status:** done（写腿自验收；对抗验收归另一位）
**Leg:** agent-75（2026-11-15 20:0x–2x，单写腿；commit 见本仓 HEAD 记录）
**Spec refs:** SPEC-04（语音链路）、SPEC-10（测试与验收）
**Inputs:** 票面 `.scratch/wisp/issues/455-*.md`（已按侦察改口径 AC#1′–AC#4′）＋侦察件 `455-onnxruntime-binding-recon.md`（已证事实直接用，未重证）

---

## 0. 交付的四件（对票面 AC 编号）

| AC | 交付件 | 落点 |
|---|---|---|
| AC#1′ 绑定自检（test-only fail-fast） | `Test455OnnxruntimeBindingIsThePinnedImage` | `cmd/wisp/a455_onnxruntime_binding_windows_test.go` |
| AC#2′ staged DLL sha256 核，不一致拒跑 | `stage-sherpa-dlls.sh` 每个 cp 后算 sha256 对 deps.toml pin | `scripts/stage-sherpa-dlls.sh:88-141` |
| AC#3′ 载具头注释改如实 | 头注释 :10-56 重写（真正保证＝位置保证，附那对数字） | `scripts/wisp-cli-staged-tests.sh:10-56` |
| AC#4′ 一次载具形整包＋红册分层 | 日志 `.scratch/wisp/probes/455/carrier-fullpkg-455.log` | 本文 §4 |

**顺手修（非本票 AC，属写面 scripts/**）：path-length 门 RED 复原——`scripts/check-path-length-budget.sh:390` 名册补 `455-*.md`（票面文件是编排者 2026-11-15 挂的、137 字符超帽且不在名册；改名会孤儿化侦察件与派单的引用，按票 262 guard A/B/C/D 记名册条目，69=69 GREEN）。

---

## 1. AC#1′ 绑定自检：它量什么、怎么红

### 1.1 仪器形状

- `GetModuleHandleEx(flags=0, "onnxruntime.dll")`（≡ GetModuleHandle 语义：只认**已映射**模块、绝不自行加载，x/sys v0.49.0 无裸 GetModuleHandle，实测签名）＋`GetModuleFileName` 取实际加载路径，`os.Stat` 取字节。
- 断言两层：①字节＝`deps.toml` pin 的 **17799168**（sha256 那半边由 AC#2′ 的 stage 门逐位保证，测试内不重算 sha256——进程仍 mmap 着该镜像时的重算会引入第二个失败模式）；②路径在仓内（从绑定路径向上找 `deps.toml` 仓根——**不信任 cwd**，go test 与载具两种发射 cwd 都成立；cwd 两级上溯只作 fallback）。
- ⛔ 没用 `os.Executable()` 当断言：exe 目录只是标准序**第一站候选**，不是绑定本身；它只出现在红句的**诊断**里（区分"exe 旁 DLL 被偷换"与"exe 旁根本没有 DLL"两形，让每形的红句都诚实）。
- ⛔ 没写 TestMain：文件名 `a455_` 按字母序排全包第一 ⇒ 该自检在**每次**整包发射的第一发跑（fail-fast 的落点）。
- 拿不到句柄（语音运行时没加载）＝红，票面逐字要求；空路径、stat 失败也＝红，不许糊成 PASS。

### 1.2 三臂实测（本机现跑 2026-11-15 20:30–20:41，同一二进制三形）

绿臂（pinned DLL 摆 exe 旁，`.scratch/wisp/tmp-455-bare/`）：

```
=== RUN   Test455OnnxruntimeBindingIsThePinnedImage
    a455_onnxruntime_binding_windows_test.go:180: 绑定自检通过：本进程绑定的 onnxruntime.dll = D:\...\tmp-455-bare\onnxruntime.dll（17799168 B，＝ deps.toml pin 的那一枚，onnxruntime 1.28.2）
--- PASS: Test455OnnxruntimeBindingIsThePinnedImage (0.00s)
```

红臂①（**A1880 行为牙**：把 System32 那枚 10572960 B 真件偷换到 exe 旁 ⇒ 必须红）：

```
    a455_onnxruntime_binding_windows_test.go:134: 绑定自检（票 455 AC#1′）：本进程实际绑定的是 D:\...\tmp-455-bare\onnxruntime.dll（10572960 B），而不是仓里 pin 的那一枚（17799168 B，onnxruntime 1.28.2）。exe 旁边摆的 DLL 被换成了别的文件——换走 staged DLL ⇒ 必须红（A1880）；stage-sherpa-dlls.sh 的 sha256 门本应先拒跑。请用 scripts/wisp-cli-staged-tests.sh 跑 ./cmd/wisp。
--- FAIL: Test455OnnxruntimeBindingIsThePinnedImage (0.00s)
```

红臂②（**裸 go test 形**：exe 旁无 DLL，第一站空、第二站 System32）：

```
    a455_onnxruntime_binding_windows_test.go:138: 绑定自检（票 455 AC#1′）：本进程实际绑定的是 C:\Windows\SYSTEM32\onnxruntime.dll（10572960 B），而不是仓里 pin 的那一枚（17799168 B，onnxruntime 1.28.2）。这一发是裸 go test 形：测试 exe 在临时构建目录（%TEMP%\go-build…\b001\），旁边没有 DLL，标准序第一站空、第二站 C:\Windows\SYSTEM32\onnxruntime.dll（10572960 B）被先命中。请用 scripts/wisp-cli-staged-tests.sh 跑 ./cmd/wisp（它把 pinned DLL 摆到 exe 旁边，标准序第一站就命中）。
--- FAIL: Test455OnnxruntimeBindingIsThePinnedImage (0.00s)
```

⇒ 绑定真的跟着「exe 旁边是哪枚 DLL」走（同一 exe、只换 DLL，绑定与红绿即时翻转）——这就是 A1880 要的**行为**牙，不是字符串断言。

### 1.3 断言行号（file:line）

- 红臂①（偷换）断言：`cmd/wisp/a455_onnxruntime_binding_windows_test.go:134-136`
- 红臂②（System32）断言：`cmd/wisp/a455_onnxruntime_binding_windows_test.go:138-141`
- 仓内路径断言：`cmd/wisp/a455_onnxruntime_binding_windows_test.go:173-178`；根定位 walk-up：`:150-167`
- 绿通过句：`:180`

---

## 2. AC#2′ sha256 门

- 位置：`scripts/stage-sherpa-dlls.sh:110-141`（cp 落地后逐枚算 sha256，对同一节 deps.toml pin；**MISMATCH ⇒ `exit 1` 拒跑**，不是警告）。
- pin 读取：awk 锁 `[sherpa-onnx.dll."<name>"]` 节内的 `sha256` 行（与 DLL 名同一节，pin 加了自动被核）。
- 两枚测量中抓出的解析坑（都已写进注释）：①GNU sha256sum 对含 `\` 的 Windows 路径整行加 `\` 前缀 ⇒ 传入路径先 `tr '\\' '/'` 归一＋输出剥前缀（不修则绿形都会假红）；②无 hasher 的 shell ⇒ GUARD exit 1，不软跳过。
- 三形实测（2026-11-15 20:38–20:39）：
  - 绿形（真 DLL）`rc=0`，逐字 `stage-sherpa-dlls: staged onnxruntime.dll -> ... (sha256 ok)`
  - 错源件形（沙盒树偷换 third_party 源件 ⇒ rc=1）：`stage-sherpa-dlls: GUARD - ... sha256 MISMATCH: got 98e8288f...ef59, deps.toml pins 7f66f939...27d8a25`（沙盒副本树测的，没碰真 third_party）
  - 反斜杠路径形 rc=0（escape 前缀坑的回归验证）
- ⛔ 未改 deps.toml 任何 pin。

## 3. AC#3′ 载具注释改如实

`scripts/wisp-cli-staged-tests.sh:10-56` 重写：从「摆了 pinned DLL 就不会被抢（假话）」改成**它真正保证的东西**——staging 把 pinned DLL 摆到**测试 exe 旁边**，标准序第一站（exe 自己的目录）命中 pinned 1.28.2，根本走不到 System32；裸 `go test` 的 exe 在 `%TEMP%\go-buildNNN\b001\`，第一站空、第二站 System32 的 1.17（10,572,960 B）应答。两对数字（17,799,168 B/1.28.2 vs 10,572,960 B/1.17.260311）写进注释当证据，并点名两道新牙（自检测试＋sha256 门）。

---

## 4. AC#4′ 载具形整包

**一发，载具形**（`scripts/wisp-cli-staged-tests.sh`，日志 `.scratch/wisp/probes/455/carrier-fullpkg-455.log`，2026-11-15 20:41–21:02，1487s）：

| 读数 | 值 |
|---|---|
| 载具横幅 | 4 行（`wisp-cli-staged-tests.sh:`）＝**这一发真的走了载具** |
| stage | 三枚 DLL 全部 `sha256 ok`（AC#2′ 门在整包里工作，逐字 `staged onnxruntime.dll -> ... (sha256 ok)`） |
| 绑定自检 | **全包第一个 `=== RUN`**＝`Test455OnnxruntimeBindingIsThePinnedImage`，**PASS 0.00s**（fail-fast 落点成立） |
| `=== RUN` 行 | 896（含子测试的 grep 计数） |
| `--- PASS` | 651 |
| `--- FAIL` | **2**（全他因：`TestAC4FocusReturnToPriorWindowGap33r5`＝64259b818 具名的本机红；`TestAC246DevLegIgnoresTheTestTaskInjection`＝侦察件 §2.1 具名他因红） |
| `--- SKIP` | 17 |
| **绑定层红** | **0**（`shadowing the one on PATH` 全日志 **0 行**） |
| **终点锚 `TestZZ315A1*`** | **整包未命中**：包死在 **`Test16r1PremortemModelSha256CorruptRefusesTheHandOff`**（最后一条 `=== RUN`）——模型 hand-off 的网络等待（进程单条 ESTABLISHED 到 `198.18.0.109:443` 代理 fake-ip、CPU 20s Δ0.3s、日志 60s+ 零增长）⇒ **按 A1886 停掉**（kill pid 46964），载具收尾 `FAIL  github.com/CarlosShao/wisp/cmd/wisp	1487.00s` |
| 与票 314 硬形态的区别 | **不是**「整进程消失」：进程一直活着、卡在网络等待，kill 是 A1886 的主动动作；载具合成 FAIL 行存在（314 形是 GUARD B 报"包没有结果行"） |
| 终点锚补证 | 载具形定向 `-test.run '^TestZZ315A1'`：三枚全 PASS（0.33s/0.33s/0.22s）＋ `ok  github.com/CarlosShao/wisp/cmd/wisp	1.00s`，rc=0（日志 `probes/455/carrier-zz315a1-455.log`）⇒ **锚活着，整包未命中是包中段卡死，不是锚坏** |

⇒ 派单问的「可能仍死在票 314 的硬形态」照实答：**这发死的不是 314 的消失形**，是 Test16r1 的模型 hand-off 网络等待（A1886 具名停掉）。4 数格式（RUN/PASS/FAIL/SKIP 的收尾行）因主动停掉而不存在，以上全为 grep 逐行计数、可复跑。


## 5. 门禁（AC#5）

| 门 | 读数 |
|---|---|
| d22scan | `clean - no D22 ban violations`，rc=0（2026-11-15 20:4x，-root 全仓） |
| path-length | 起手 RED（455 票面未名册，pre-existing）→ 名册补条后 `VERDICT GREEN`（roster 69=tree 69），rc=0 |
| go vet ./cmd/wisp | rc=0 |
| gofmt（本腿新件） | clean（仓内另两枚旧产码 `approval_always.go`/`approval_reply.go` gofmt 脏为 pre-existing，未动） |

## 6. 边界

- 写面：`cmd/wisp/a455_onnxruntime_binding_windows_test.go`（新，test-only，零产码）、`scripts/stage-sherpa-dlls.sh`、`scripts/wisp-cli-staged-tests.sh`（仅注释）、`scripts/check-path-length-budget.sh`（名册数据行）、本文。⛔ thresholds/golden/allowlist/PLAN/specs/D43/frontend/design/deps.toml pin/System32 断言全未动；⛔ 未 push；⛔ 未跑裸 `go test ./cmd/wisp`（整包唯一一发经载具，跑前查 `Get-Process wisp.test`）。
- 红臂①的偷换发生在 `.scratch/wisp/tmp-455-bare/`（临时目录，不入库）；sha 错源件形在沙盒副本树测（`.scratch/wisp/tmp-455-sbx/`）。
- 台账（A1880/A1884/A1886 均按令执行）：A1880 行为牙＝红臂①；A1884 两枚 sh 门全程 `D:\work\soft\Git\bin\bash.exe`（`bash -n`、stage 三形、path-length、载具整包）；A1886 整包后台飞、不干等 CPU。
