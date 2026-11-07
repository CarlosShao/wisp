# 票 65 · 像素测试真 flake 治理（`TestBallLiveOrbPageRendersPixels` 偶发全黑红）

- **日期**：2026-11-07 15:4x–16:3x +08
- **写腿**：Wisp 编队写腿（GLM 5.3 Flash，本件自证）
- **写面**：`internal/ball/orb_live_windows_test.go`（仅像素采样段＋新增测试侧收割 helper）＋本件＋票 65 Progress log 一行
- **提交**：`ea602a61b`（有界重建重试）→ `79fdc5a8b`（泄漏收割，根因对.shape）→ 证据件＋log 追加（只 commit 不 push）
- **禁手自查**：未动 thresholds/golden；未动 frontend/**、design/**；未动 assets/orb-siri.html；
  未动 `internal/ball/orb_windows.go` 产码（工作树里该文件 221 行 diff 是**另一条在飞腿**的活，b576c10a3
  落盘，与本治理无关）；未翻 AC 勾、未改票 65 Status；临时件只建不删。

---

## 1. 复现与读数（全部本腿亲跑，非转述）

### 1.1 编排者既有证据（接票输入，非本腿产出）
- 第 1 发 `go test -tags winlive ./internal/ball/ -run 'TestBallLiveOrb' -count=1` ⇒ **FAIL 32.822s**（0.000% 全黑）
- 第 2 发同命令 ⇒ **PASS 4/4，564/1764 = 31.97%**
- ⇒ 同一代码相邻两发一红一绿 ⇒ 真 flake。

### 1.2 修前基线 A（只跑 PageRendersPixels，count=5）
| 发 | 结果 | 单发时长 | 非黑像素 |
|---|---|---|---|
| 1 | PASS | 33.66s | 31.80% |
| 2 | PASS | 30.29s | 30.22% |
| 3 | PASS | 14.39s | 40.59% |
| 4 | PASS | 28.94s | 29.54% |
| 5 | PASS | 29.85s | 16.44% |

单发首亮延迟实测可到 ≈30s（贴住并偶超当时 25s 预算）——但这只是浅层读数，深层机理在 1.4。

### 1.3 修前基线 B（全套 4 枚 count=10，第一回，管道只留了计数）
**29/40**（应有 40），原始日志未留存（管道设计失误：`grep -c` 吞了明细）。以 1.4 的第二回为准。

### 1.4 ★决定性证据（全套 4 枚 count=10，第二回，全量落盘 `/tmp/orb_suite10.log`）
**39/40 PASS，1 红**；红的就是 PageRendersPixels（103.25s＝3 发重试全黑用尽）。
10 轮像素非黑率曲线：

```
40.59 → 40.31 → 39.29 → 40.65 → 40.25 → 40.65 → 40.59 → 27.15 → 7.71 → 0.000(×3发)
```

**前 7 轮钉在 ~40%，随后坍塌 27 → 7.7 → 黑。** 这不是随机摇摆，是**逐轮累积退化**。

## 2. 根因判读（有决定性读数，非猜测）

1. **等待形状从来不是主因**：测试本来就是「轮询合成屏像素直到非黑或超时」（任务书首选的真信号形状）。
2. **泄漏堆积才是机理**：每轮测试建一个 WebView2 实例，产码 `destroy()`（orb_windows.go:387）把
   `o.browser` 置 nil 后直接 `pDestroyWindow`，**从不关闭 controller**——而 jchv 包装层的
   `ICoreWebView2Controller` **根本没有 Close 方法**（vtable 缺失），所以每个 orb 的
   msedgewebview2.exe 浏览器进程被永久泄漏在本测试进程名下。count=10 连跑时泄漏逐轮堆积 GPU 压力，
   新实例渲染逐轮变暗（1.4 曲线）直到画不出来。这也解释了 standalone 连绿（泄漏少）与
   in-suite 连跑变红的「交互」表象，以及编排者那发 0.000%。
3. **为什么修测试侧而不是产码**：对症的产码动作是 controller.Close，包装层无此方法（需 fork 第三方模块
   ＋replace，超出本缺陷票射程，另立票）；destroy() 的销毁顺序是票 65 实测钉过的形状，动它风险外溢。
4. **为什么不用 NavigationCompleted 事件等待**：包装层回调只能在 `Navigate()`（New 内部）之后安装，
   有错过窗口；像素轮询本身就是断言对象，保留为主信号。
5. **为什么不是裸加长预算**：掩盖非治理；重试＋收割直接拆机理。

## 3. 修法（测试侧，判据零改动）

`orb_live_windows_test.go`：
1. **有界重建重试（3 发）**：`orbPixelAttempt` 单发形状不变（attach 轮询 15s、像素每 300ms 一次、
   单发预算 25s、判据 `ratio > 0.01` 一字未动）；全黑到期返回 false，外环 `b.Close()`（幂等）整套拆掉
   后重建；三发连黑才 Fatal，报三发最好 ratio。attach 超时／窗口 API 失败仍硬 Fatal 不重试。
2. **泄漏收割 `orbTestReapLeakedWebView2`**：Toolhelp32 枚举进程，**父 PID＝本测试进程且映像名＝
   msedgewebview2.exe** 才 Terminate（双条件 ⇒ 其他应用的 WebView2 零命中），测试开始时＋每发重建前
   各清一次，清掉本进程累积的全部泄漏浏览器。
3. backdrop 从 `t.Cleanup` 改为 attempt 内显式销毁：它靠该球 STA 线程泵绘制，必须先于 Close，
   否则 post 到已退出的 STA 会挂死。

## 4. 验证（判定牙）

- **修后全套 `TestBallLiveOrb`（4 枚）`-count=10`：40/40 全绿**（go test rc=0，总 204.4s，
  `/tmp/orb_suite10_v3.log`）。
- 修后像素非黑率曲线：`36.21 13.42 31.01 26.53 13.94 20.29 24.97 29.24 22.89 30.80`
  ——零坍塌、零重试触发（每轮 attempt 1 即亮）。（整体从 ~40% 移到 ~13-36% 区间＝同树另一腿
  b576c10a3 改渲染器所致，与本次治理无关，如实记录。）
- 收割器按预期每轮开工：10 轮各 reap 1 个泄漏进程（`reaped 1 leaked WebView2 browser process(es)`）。
- 修后单枚 PageRendersPixels `-count=10`：10/10 绿（21.92s–35.39s，总 272.2s，reaper 加入前）。
- 门禁：`gofumpt -l internal/ball` 空；`go vet ./internal/ball` rc0；`go vet -tags winlive ./internal/ball`
  rc0（另一腿落盘后）；`go test ./internal/ball -count=1` ok 0.209s；d22scan clean。

## 5. 遗留与如实

- **产码级根治未做**：destroy() 泄漏 controller 是真实缺陷（每建一个 orb 泄一组浏览器进程），
  根治需给 jchv 包装层补 `ICoreWebView2Controller::Close` vtable（fork＋replace，第三方模块）——
  建议另立票；本治理把它的**测试可见后果**（累积黑）拆掉了。
- 基线 A 的「25s 预算 vs 30s 慢 boot」是浅层读数，被 1.4 的累积证据取代；两段都留档不删。
- 基线 B 第一回 29/40 明细未留存（管道失误），判定不依赖它。
- 治理期间工作树有另一条腿在飞（orb_windows.go/renderer_windows.go），已等其落盘（b576c10a3）后
  再跑最终判定；本腿两枚 commit 各自自洽（引用符号在各自时点均存在）。
