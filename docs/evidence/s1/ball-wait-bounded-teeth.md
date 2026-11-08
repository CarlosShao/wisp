# 票 408 · `waitStarted` 有界等待的验牙报告（变异三发 + winlive 正控）

- 时刻：2026-11-08 20:2x +08
- 执行腿：Wisp 编队写腿（验牙腿，派单点名列本件）
- 被验对象：产码 `internal/ball/sta_windows.go` 的 `bootWaitBudget`（:227，`30 * time.Second`）与
  `waitStarted`（:235-252，`select{<-s.started / <-timer.C}` 超时返回具名错误，`062e5fe70` 落地）；
  牙件 `internal/ball/sta_wait_started_bounded_windows_test.go`（两测试：楔死形必须超时返回具名错误＋真 boot 错误必须原样穿透）。
- 方法：一律在 `git clone --shared` 导出树（`/d/work/workspace/wisp-mut-waitstarted-408`，HEAD `b7a2b6721`，含 `062e5fe70`）里做变异，主树零改动。

## 0. 起手三核（现跑，不信转述）

1. 牙全文现读：测试一 `TestWaitStartedBoundedReturnsNamedError` 把 `bootWaitBudget` 现场调成 **200ms**（:28-30），用 `newSTAThread(nil)` 造「永不 `close(started)`」的楔死形，goroutine 跑 `waitStarted`，外挂 **5s 守卫** select（:37-49）：必须回非 nil 且含 `ui-sta thread did not report window creation within` 的具名错误；守卫先响 ⇒ 报「unbounded again」。测试二 `TestWaitStartedStillPassesBootErrThrough`（:55-68）钉「真 boot 错误原样穿透，不被超时错误顶掉也不丢成 nil」。
2. 产码现读：`git show 062e5fe70` diff 与盘上现状逐行一致（:227 var 定义、:235-252 select 形）。
3. `grep -rn bootWaitBudget --include=*.go internal/ cmd/`：仅 `sta_windows.go`（定义 :227＋使用 :240/:246）与牙文件（:28-30/:48/:56-58）——无其他耦合点，变异面干净。

## 1. 基线（未变异，导出树）

```
=== RUN   TestWaitStartedBoundedReturnsNamedError
--- PASS: TestWaitStartedBoundedReturnsNamedError (0.20s)
=== RUN   TestWaitStartedStillPassesBootErrThrough
--- PASS: TestWaitStartedStillPassesBootErrThrough (0.00s)
```

## 2. 变异三发（核心读数）

| 发 | 变异 | 预期 | 实测 | 判 |
|---|---|---|---|---|
| A | `waitStarted` 改回裸 `<-s.started`（删 select/timer） | 牙超时红（守卫先响） | **红** 5.00s，红句逐字见 §2.1 | 有牙 |
| B | `bootWaitBudget` 改 `100 * time.Millisecond` | 牙仍 PASS（牙自带楔值，不依赖真 30s） | **PASS** 0.20s | 有牙 |
| C | 超时分支 `return fmt.Errorf(...)` 改 `return nil` | 牙红（证明真断言「具名错误」） | **红** 0.20s，红句逐字见 §2.2 | 有牙 |

### 2.1 变异 A 红句（逐字）

```
    sta_wait_started_bounded_windows_test.go:48: waitStarted did not come back within the 200ms budget + guard: the wait is unbounded again
--- FAIL: TestWaitStartedBoundedReturnsNamedError (5.00s)
```

守卫 5s 先响、按设计响红，不是整跑挂死。测试二在裸接收下仍 PASS（它走的是已信号路径）——正确。

### 2.2 变异 C 红句（逐字）

```
    sta_wait_started_bounded_windows_test.go:42: waitStarted returned nil on a never-signalled boot: the bounded wait is gone
--- FAIL: TestWaitStartedBoundedReturnsNamedError (0.20s)
```

0.20s 即红（budget=200ms 到点走 `return nil` 被断言抓到）——牙断言的是「有界等待返回**具名错误**」，不是「只是不挂」。测试二仍 PASS（穿透路未被动）。

### 2.3 三发还原核验

`git show HEAD:internal/ball/sta_windows.go` 覆盖回去后：导出树 `git diff --stat` 空、`git status --porcelain --untracked-files=no` 空。导出树保留变异脚本 `mutA.sh`/`mutB.sh`/`mutC.sh`/`restore.sh`（临时件只建不删）。

**导出树已知伪红（与本验牙无关，如实记）**：全包 `go test ./internal/ball` 在导出树红 3 枚资产比对测试（`TestOrbAssetShippedUnmodified`/`TestOrbServingConstants`/`TestJarvisAssetShippedAsBuilt`）——原因是 `build/` 与 `third_party/sherpa-onnx` 未被 git 跟踪、导出树里根本没有（`git check-ignore` 实证）；与票 408 Progress log 既有记录（autocrlf 伪红＋DLL 缺失）同族。权威绿测以主树为准（§3）。

## 3. 主树权威读数

- winlive 正控（首跑，`-count=1 -timeout 30m`，339.974s）：**9/10 PASS＋1 红**——红＝`TestBallLiveOrbPageRendersPixels`（95.26s）。该红与 `webview2-nil-env-guard.md` §6:92 记录的「暗页 GPU 抖动族」签名一致（同代码首跑报三试全黑红、重跑绿；同跑 `TestBallLiveOrbCornersReadTheDesktop` 95.17s 同时长却 PASS）。**无任何一枚因 30s 有界等待变红**（预算误伤的形状＝`did not report window creation within 30s` 具名错误，全程未出现）。
- 抖动族取证重跑：首跑曾留下僵尸 `ball.test.exe`（pid 14676，title 带上一测试名，挡住桌面独占检查致单测重跑 0.01s 秒红——测试自带的 foreign-window 守卫，按其处置令杀掉后重跑）→ 清残留后 `TestBallLiveOrbPageRendersPixels` 单测重跑 **PASS 32.58s**。
- winlive 结论：**全家族绿（首跑 9/10＋抖动族重跑绿＝10/10 有效）**；健康启动读数 0.82-7.87s（`TestBallLiveOrbJarvisStyleRendersParticles` 7.87s＝窗口腿最慢读数）远小于 30s 预算，有界等待不误伤正常启动。
- 全包 `go test ./internal/ball -count=1`：主树 **ok 0.520s**（含牙两测试）。

## 4. 30s 预算裁定：**维持 30s，不改**

理由（按用户视角／健康上界／慢机器三候选逐条过）：

1. **健康上界**：本树实测最差整包启动 7.9s（进程内第二枚 orb，= 单跑 6 倍；票 408 r2 读数），且 `waitStarted` 只罩「窗口创建腿」——整包里最快的部分。30s ≈ 最差健康值的 4 倍，余量充足。
2. **慢机器/CI**：窗口创建是本地 Win32 + WebView2 环境建，纯本机操作；就算 CI 慢 4 倍也远够。代价不对称性决定预算**宁长勿短**：预算太短 ⇒ 健康慢启动被误判成楔死、报假错并走降级（坏——把活球判死）；预算太长 ⇒ 真楔死晚 30s 才浮出（可接受——一次性、有界、有具名错误）。偏长的错误代价小。
3. **用户视角**：球 30 秒不出来＝坏了，这个用户判断合理；且到期不是静默——`New` 返回具名错误，常驻层能报死球并降级。
4. **实现细节非契约**：`bootWaitBudget` 是 var（测试可楔短，产线从不写），将来要调是一行改动＋牙全绿即可复核。

## 5. ★如实发现：牙件此前从未入库

三核第 1 步时实证：`git ls-files internal/ball/` 里**没有** `sta_wait_started_bounded_windows_test.go`——它在主工作树是未跟踪状态（`git status --porcelain` `??`），即 `062e5fe70` 只入库了产码、牙漏提交。若不补，任何人 clone 后都没有这枚牙。本腿已把它随证据件一并入库（写面 `internal/ball/**` 本波独占内）。

## 6. 门禁读数（主树，全 LF）

| 门禁 | 读数 |
|---|---|
| gofumpt（`/d/work/base/gopath/bin/gofumpt.exe -l`，按写面 `internal/ball/`） | 空 |
| `go vet ./internal/ball/...` | rc0 |
| `go build ./internal/... ./cmd/...`（按域；`go build ./...` 被存量探针垃圾件 build\\mut_rvw.go 挡红，非本腿面） | rc0 |
| `go test ./internal/ball -count=1` | ok 0.520s |
| d22scan（`scripts/d22scan.sh`） | clean |
| winlive `TestBallLiveOrb` | 首跑 9/10＋1 抖动族红，取证重跑绿 ⇒ 有效 10/10（§3） |

注：全仓 `gofumpt -l .` 列出的百余件全部在 `.scratch/`、`build/`、`forks/` 的存量探针/临时件里（历史各腿产物，非本腿写面，本腿未触碰）；本腿三件（牙/证据/票）gofumpt 零列名。

## 7. 结论

**这枚牙有牙**：三发变异全部按预期落色（A 红/B 绿/C 红），还原干净，winlive 正控不因它变红。`062e5fe70` 的有界等待实现＋牙件成立。
