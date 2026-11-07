# 票65 · overlay 弧回归取证与裁定（65-orb-overlay-arc-restore.md）

- 日期：2026-11-07 18:4x–19:1x +08
- 腿：wave68（编排者派单「外圈青色进度弧丢了」，初步定位＝c5cca231a 给 `renderer_windows.go:521`
  `RingColor` 加的 `!r.overlayOnly` 守卫把弧一起关了）
- HEAD 起点：`d902fbc6b`（Jarvis 腿步5），工作区 `git status --porcelain internal/ball` 空
- 本件结论先行：**初步定位不成立——弧从未被关掉，也不需要恢复**。真正落的是「弧存在」钉子测试。

---

## 步1 逐行审计：overlay 分支每个 `!r.overlayOnly` 守卫到底关掉了什么

现跑取证（`git show c5cca231a -- internal/ball/renderer_windows.go` 逐 hunk ＋
`git log -S` 双查），overlay 分支（`renderer_windows.go:473-617`）现存的 6 处守卫：

| 行 | 守卫 | 关掉的是什么 | 类别 | 出生 commit | 裁定 |
|---|---|---|---|---|---|
| `:521` | `!r.overlayOnly && v.RingColor.A > 0` | 步骤6 **状态环**（`RingColor`/`RingPulse` 危险脉冲环，`circle(c, R+5*s, v.RingWidth*s)`） | 图形（非文字，**也非进度弧**） | **`10d7cb3c3`**（与 overlay 架构同生） | **保留关**：源注释原文「Skipped in orb mode: the page renders the state; a D2D ring around the web orb would draw the state twice」——恢复它＝在页面状态面之上**双重画状态环** |
| `:533` | `!r.overlayOnly && v.Icon == IconAudioLines` | 步骤7 听写三波纹 | 图形 | **`10d7cb3c3`** | 保留关：页面音频波就是音频面（同上注释「the page's audio bands are the audio face」） |
| `:546` | `!r.overlayOnly && v.Sweep` | 步骤8 思考光带 | 图形 | **`10d7cb3c3`** | 保留关：页面思考态是面 |
| `:571` | `!r.overlayOnly && !(v.Glass && wensityActive())` | 步骤10 图标线稿 | 图形 | **`10d7cb3c3`** | 保留关：页面承载面 |
| `:588` | `!r.overlayOnly`（badge `drawText`） | 步骤11 badge **白色数字**（Danger 盘本体不关＝无数字通知点） | **文字** | **`c5cca231a`** | 保留关：票65「球面无文字」法（本腿原始目标之一） |
| `:614` | `v.BadgeText != "" && !r.overlayOnly` | 步骤13 **中心文本**（「69%」） | **文字** | **`c5cca231a`** | 保留关：同上 |

**两处关键 git 现跑证据**（逐字命令可复跑）：

1. `git log -S "if !r.overlayOnly && v.RingColor.A > 0" -- internal/ball/renderer_windows.go`
   ⇒ 唯一命中 **`10d7cb3c3`**（票65透明+单球腿），`c5cca231a` **不在结果里**。
2. `git show c5cca231a -- internal/ball/renderer_windows.go` 的全部功能 hunk ＝
   ①`:588` badge 数字包 `if !r.overlayOnly`、②`:614` 中心文本 `&& !r.overlayOnly`、
   ③新增 `overlayGlyphPixels()` 审计函数、④`bits uintptr→unsafe.Pointer`（类型改写，零行为）、
   ⑤注释三段。**`:521`/`:533`/`:546`/`:571` 一个字节都没动**；同笔 `ball_windows.go` 改动为纯注释。

## 弧的真相：步骤9 进度弧**没有守卫**，overlay 模式下一直在画

`renderer_windows.go:556-564`（步骤9，HEAD 现读）：

```go
// 9. Progress ring (Downloading).
if v.Progress >= 0 {
    r.setSolid(mulA(pal.BorderSoft, v.Opacity))
    r.circle(c, R+5*s, 2*s)
    if v.Progress > 0 {
        r.setSolid(mulA(pal.Accent, v.Opacity))
        r.arc(c, R+5*s, -90, -90+360*v.Progress, 2*s)
    }
}
```

- **零 `!r.overlayOnly` 守卫**：`ball_windows.go:436` `b.rend.overlayOnly = b.orb.Load() != nil` 只改剪影分支，
  Visual 完整携带 `v.Progress`（`repaintCurrent:407`）。Downloading 69% 时 overlay 窗照画
  BorderSoft 底环＋**`pal.Accent`＝#86C2B9 青色 69% 弧**。
- 这就是 `orb-FINAL.png` 里那圈青弧：截图描述为「wrapping **partway** around」＝扫 360×0.69°
  的进度弧形状（页面装饰弧是整圈，只有 D2D 进度弧按 Progress 走部分圈）。
- wave65 的证据件（`65-orb-overlay-badge.md`）当时就写明：「进度＝外圈细弧（步骤9 BorderSoft 底环+
  Accent 进度弧，修前就有、截图里的青弧即它，**本腿零新增绘制**）」。
- 旁证（`orb-siri.html` 只读 grep）：页面虽有 `arc/progress/ring` 字样，但页面冻结 API 只有
  `setState/setAudioBands`（orbasset.go 钉死的接缝），**没有进度通道**——页画不出 69% 部分弧。

## `page lit ratio` 21.96%→3.02% 的机理：与弧无关（读数钉死）

- `TestBallLiveOrbCornersReadTheDesktop` 全程 **Sleeping 态**（`:695` `Initial: StateSleeping`，
  测试体从不 `SetProgress`）⇒ 该测试里 **D2D 进度弧从来没有画过**，21.96% 里不可能含弧，
  3.02% 里也不可能缺弧。编排者「19% 亮像素就是外圈弧」的前提在这枚测试的采样对象上不成立。
- 该读数的返回形状（`orb_live_windows_test.go:814`）：`cornersOK && lit > 0.01` **首次过线即返回**——
  页面 WebGPU 淡入期间首次过 1% 线可停在任意低值。同页非黑率的宽区间早有在案读数
  （台账 A1602：全套 count=10 实测曲线 `40.6/…/27.2/13.4/7.7%`，判定为**时序摇摆**非回归）。
- **HEAD（含 c5cca231a）现跑读数**（本腿 2026-11-07 18:5x，逐字 log）：
  - 第 1 发：`page lit ratio 25.91%`，四角 `[[51 51 50]×4]` 精确零差，PASS 27.06s
  - 第 2 发：见下文读数表（补强方差证据）
  - `TestBallLiveOrbOverlayFrameHasNoGlyphPixels` PASS 0.90s（近白像素 0）
  ⇒ **3.02% 是单发早采样 flake，不是回归**；修后基线仍在 ~20% 量级。

## `orb-noglyph2.png` 的「弧没了、球小了」：状态错位解释（真机复核见读数表）

- `orb-FINAL.png`（修前）与 `orb-noglyph2.png`（修后）**不是同一状态的对照**：
  - FINAL＝balldebug 20 态循环的 **Downloading 步**（69% 文本＋D2D 进度弧同帧在）；
  - noglyph2（任务栏时间 18:37）＝**无 Progress 的某个态**（静止带状面、无 badge 盘、无弧）。
    `cmd/balldebug/main.go` 的 `-state` 单态模式只调 `b.SetState(s)` **不调 SetProgress/SetBadgeText**
    （`:317-333`），默认 20 态循环除 Downloading 步外也全无进度 ⇒ 任一非 Downloading 时刻截图都「无弧」。
  - 球 D2D/页面本体在两版本间零改动（页面 sha `d93da649…` 钉测绿、renderer 弧路径零触碰），
    「球变小」＝Downloading 弧＋底环带来的外圈视觉体量差异＋两态页面液面不同。
- 真机复核：本腿用 HEAD 重建 `balldebug-orb.exe -shots` 跑 20 态循环，Downloading 帧
  **弧在、文字无**——逐字读数见下文读数表（这就是编排者要的「弧还在」的屏上证据）。

## 步2 裁定（铁律 ⓪⓪⓪：判据＝最有价值，非动作最小）

**裁定＝不恢复 `:521` 的 `RingColor` 弧（它不是进度弧），产码零改动；把 wave65 漏掉的「弧存在」
断言补成钉子测试。** 理由：

1. **恢复 `:521` 买不到任何价值**：它关的是状态环，不是进度弧；恢复＝页面已渲染的状态环之上
   再叠一圈 D2D 状态环＝`10d7cb3c3` 注释点名的 double-draw 缺陷（画两遍状态）。
2. **弧没有被关**：步骤9 无守卫、overlay 照画，25.91% 的 lit ratio 与真机 Downloading 帧
   双向证明弧在屏上。「过度收窄」的假设与 git 历史、像素读数、屏上截图三方矛盾。
3. **真正的缺口**（编排者 A1606 自己点了）：wave65 只断言「无文字像素」，没断言「进度弧还在」
   ——谁将来在步骤9 上加守卫，现有判据全绿、弧静默消失。补
   `TestBallLiveOrbOverlayFrameKeepsProgressArc`（DIB 审计 Accent 青色像素 ≥50）把这颗牙钉死。
4. 排除项：乙「恢复 :521」＝double-draw 缺陷＋零证据支撑；丙「只补文档不动测试」＝缺口还在，
   下一条腿会再踩同一个坑（本次 wave68 派单本身就是这个坑的第一枚果实）。

## 撤销口令

- 钉子测试＋证据件：`git revert <本腿 commit>`（产码 `renderer_windows.go` 零改动，无产码回退面）。

---

## 读数表（全部本腿现跑，2026-11-07 18:4x–19:1x +08；逐字 log 可复跑）

### 派单三条判据

| 判据 | 读数 | 结论 |
|---|---|---|
| ①球面无文字：`TestBallLiveOrbOverlayFrameHasNoGlyphPixels` | **PASS**（终态复跑 1.47s；全套内 0.90s）——Downloading 69%+"69%"+badge3 缺陷形状帧**不透明近白像素＝0**（近白判据自 wave65 红103→绿0，未动） | ✅ |
| ①b 弧存在（本腿新钉子，同一帧） | **`overlay Downloading frame keeps the progress arc: 283 teal pixels (floor 50)`** | ✅ 弧在 |
| ② `page lit ratio` 回 ~20% 量级 | 修前基线 21.96%（A1604）；本腿 HEAD 现跑 4 发正常启动：**25.91% / 24.77% / 5.20%**（首过线快照，同语义）＋0.000%（黑启动 FAIL，3 发重建重试治愈）；终态全套 7/7 PASS | ✅（页面满渲染值 ~20-26% 与修前同量级；该 log 行本身＝「corners ok ∧ lit>1% 首次过线即返回」的快照，修前 21.96% 亦是同语义快照——21.96 vs 3.02 当稳态对比从未成立；0%/3.02%/5.20%＝同一时序族，与弧无关） |
| ③ 四角零差 | 全部 PASS 发均为 **`corners [[51 51 50] [51 51 50] [51 51 50] [51 51 50]] within +/-12 of (51,51,50)`** 精确零差 | ✅ 透明未破 |

### 3.02% 与 0.000% 的归因（诚实记录，不粉饰）

- 3.02%（编排者引）与本腿第 2 发 0.000%（FAIL）同一族：**WebView2 黑启动/早采样时序**，
  非 renderer 回归——证据链：①Corners 全程 Sleeping、从不设 Progress，D2D 弧在该测试从未画过；
  ②HEAD 含 c5cca231a 的正常启动发读 25.91%/24.77%（高于修前 21.96%）；③同页非黑率宽区间
  在案（A1602 count=10 曲线 13.4–40.6%）。④本腿给 Corners 补了 orbPixelAttempt 同型的
  **3 发有界重建重试（断言零改动）**，全套内黑启动发被重建治愈。
- 「修后判据变红的可能性」（派单如实条款）：不适用——本腿**没有恢复任何被关的绘制**
  （因为没有东西被关），文字判据与弧判据在同一帧上同时成立（近白 0 ∧ teal 283）。

### 红牙证明（判据有牙）

- 突变：临时给步骤9 加 `&& !r.overlayOnly`（模拟「弧被关」）⇒
  `--- FAIL: TestBallLiveOrbOverlayFrameHasNoGlyphPixels`，红字逐字：
  `overlay Downloading frame carries only 0 progress-arc pixels (want >=50): the outer cyan arc the page cannot draw is suppressed (ticket-65 arc regression)`；
  **同帧近白审计仍 0**（弧判据变红不连带文字判据，两判据正交）。
- 还原：`git diff internal/ball/renderer_windows.go` 空（备份比对逐字节还原）。

### 真机证据（屏上，非 DIB）

- HEAD 重建 `build/balldebug-orb.exe -shots` 跑完 20 态循环（handles 全程 445→471 稳定）：
  **`13-Downloading.png`＝球环上有青色进度弧（teal 高亮段）、球心无任何文字数字**——
  与 orb-FINAL.png 的差别只有「69% 文字没了」，弧没有丢。档：`D:\tmp\wisp\wave68\shots\`。

### 门禁（终态）

| 门禁 | 读数 |
|---|---|
| `/d/work/base/gopath/bin/gofumpt.exe -l internal/ball` | 空 |
| `go vet ./internal/ball` | rc0 |
| `go build ./internal/ball ./cmd/wisp ./cmd/balldebug` | rc0 |
| `go test ./internal/ball -count=1` | **ok 0.581s**（终态） |
| `go test -tags winlive -run TestBallLiveOrb -count=1` | **7/7 PASS 69.888s（终态全套）**：Jarvis 11.75s / Attach 1.21s / Placement 0.00s / FollowStays 0.31s / PageRenders 31.22s（6.87% 首过线）/ Corners 24.03s（5.20% 首过线，四角精确零差）/ Glyph+弧 0.84s（近白 0 ∧ teal 283） |
| `go test -tags winlive -run 'TestLiveNeverStealsFocus\|TestLiveTransparentCornerFallsThrough'` | **2/2 PASS**（1.88s / 0.58s） |
| `sh scripts/d22scan.sh` | 真仓扫描 **1 发现＝`internal/panel/config_write_187.go:55` phantom-citation（票 187 腿 untracked 在飞件引了未落盘的证据件）——非本腿写面（⛔ internal/panel），本腿写面 internal/ball 零发现（ban#8 扫 951 Go 文件含新测试件）；按「每包一腿」纪律不碰，归票 187 腿收口时自愈** |

### 落点与写面

- 写面：`internal/ball/orb_live_windows_test.go`（Corners 3 发重建重试＋Glyph 测试弧断言＋
  `overlayArcTealPixels` 助手）＋本证据件＋票 65 Progress log。**产码 `renderer_windows.go` 零改动**。
- commits：`91c0532d2`（取证+裁定）＋本表随测试件与结笔 commit。

