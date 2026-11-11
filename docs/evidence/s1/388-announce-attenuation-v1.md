# 388 — 播报音量衰减 · 非实现者验收腿 v1（只读裁决 · 2026-11-11）

- **裁决者**：非实现者验收腿（388v1 dispatch，model=auto/glm-5.3-flash）
- **证据面**（全部盘上现读，未跑任何 go test/build，未改任何产码/票面）：
  - 票 `.scratch/wisp/issues/388-test-announce-attenuation.md`（Progress log wave121 / 388-closeout-N 两节现读）
  - `cmd/wisp/voice_announce_attenuation_388_windows_test.go`（148 行全文）
  - `cmd/wisp/resident_voice_interrupt_381r2_windows_test.go:453-577`（接线＋逐轮断言现读）
  - `cmd/wisp/resident_voice_interrupt_390r1_windows_test.go:153-189`（接线现读）
  - `internal/audio/aecspike/announcevolume_388_windows.go`（守卫全文）＋ `announcevolume_388_windows_test.go`（三读数＋loopback 对照全文）＋ `dmo_cgo_windows.go:88-108` ＋ `aecsmoke.c:355-405`
  - 只读 git：`git show --stat` 两枚 388 提交（wave121=28763663e、closeout-N=821f4a2d3）＋ `git diff --name-only 28763663e~1..821f4a2d3` ＋ `git status --porcelain` 禁改路径定向检查
- **裁决总表**

| AC | 判定 | 判据（可复跑） |
|---|---|---|
| AC#1 衰减机制 | **PASS** | ①机制在盘：`attenuateAnnounce388(t, attenuated)`（值参数化，调用点传 `aecspike.DefaultAnnounceAttenuatedVolume` 常量 0.15，数字不抄进包内）＝`voice_announce_attenuation_388_windows_test.go:46-61`；守卫 `ApplyAnnounceAttenuation` 先存机主原值再经 cgo 真浮点调 `SetSessionVolume`，回读差 >1e-6 即报错，apply 失败路径先复原再返回（fail-closed）＝`announcevolume_388_windows.go:76-111`。②**整型路零存在**：音量写路只走 `C.aecsmoke_set_client_volume(client, C.float(v))`＝`dmo_cgo_windows.go:91-97` → `aecsmoke.c:365-384` 真 C 调 `ISimpleAudioVolume::SetMasterVolume(sv, v, NULL)`；grep 全包，`syscall.SyscallN` 只在 `rig_windows.go:124` 的 COM vtable 整型参数路（非 float 参数），59-r1 发现②的边界守住。③接线在盘：两枚真设备测试各在**第一发 rv.announce 之前**一行前置——381r2 `:465`（首 announce `:508`）、390r1 `:165`；接线钉 `Test388AnnounceAttenuationWired`（`voice_announce_attenuation_388_windows_test.go:99-148`）从盘上逐枚验「衰减前置在 rv.announce 前」＋Apply/t.Cleanup/Restore 配对，253r1 名册钉同形。④值参数化 CI 钉：`Test388AttenuationValueIsParameterizedAndValidated`（`announcevolume_388_windows_test.go:30-43`）拒 0/负/超 1 值于任何 COM 调用之前。⑤真机三读数在册：closeout 提交信息与票面 wave121/closeout-N 均录 original=1.000000→attenuated=0.150000（回读一致）。 |
| AC#2 出声可闻度对照 | **PARTIAL**（两半拆裁，按 dispatch 拆法） | **①loopback 录回读数那半＝PASS**：wave121 交付的对照读数 750.3→46.5 peakrms（session 0.6→0.15，ratio 0.0619，该机 tap 呈平方律）三处互相印证——票面 Progress log wave121 行、提交 28763663e 信息、盘上测试 `Test388AttenuatedSessionReachesTheEndpointLoopback`（`announcevolume_388_windows_test.go:122-197`，断言 ratio∈(0.01,0.5) 且拒静音拒"没衰减"，原值链式复原）。衰减确实到达端点馈给（＝机主音箱实际吃到的信号）。**②机主亲耳确认那半＝待机主**：听感确认需机主在场，closeout 腿已在票面如实留格（"AC#2 audible-confirmation half is left to the owner/receipt leg honestly"），本腿按 dispatch 如实记**待机主**、不算欠账。 |
| AC#3 反向判据 | **PASS** | ①摘衰减＝Restore 经同一 cgo 浮点路写回原值并回读核（`announcevolume_388_windows.go:116-132`，差 >1e-6 报错）。②三读数在册：closeout 提交信息两测试均录 ATTENUATE original=1.000000→attenuated=0.150000→restored=1.000000；helper `t.Cleanup` 落 `ATTENUATE-388` 一行三读数（`voice_announce_attenuation_388_windows_test.go:52-60`）。③还原核超出票面最低要求：`Test388AnnounceAttenuationTripleReadings`（`announcevolume_388_windows_test.go:50-85`）在 Restore 后经**独立新 client** 复读 session 持久值=fresh re-read=1.000000（不只信守卫自己写的回执）。④守卫链还原核：loopback 测试两级衰减链式解缠到机主真原值（`:158-165` 逐层回读断言）。 |
| AC#4 零回归 | **PASS** | ①被测物语义零触碰（盘上现读）：waveOut 真路径、abort 计时、逐轮 ≤400ms Fatalf 全部原样——381r2 `:550/:555` 两条逐轮断言在原位原语义；衰减只走会话音量路，`announcevolume_388_windows.go` 不触碰 waveOut 流。②分布断言不松界：`assertInterruptDistribution388`（`voice_announce_attenuation_388_windows_test.go:75-90`）用与逐轮断言同一个 `interruptStopBound388=400ms` 常量，每样本同界断言，零样本=Fatalf（防空绿）；两测试各接 unstick＋waveOutReset 两路（381r2 `:576-577`、390r1 `:272-273`）。③重测读数在册（closeout 提交信息＋票面）：381r2 n=5 min=29ms median=39ms max=49ms、390r1 33/39/49ms——衰减后 abort 时延仍在 29-49ms 同量级，界 400ms，"衰减 ⛔ 影响 abort 时延"读数成立。本腿未跑 go test（硬约束），读数以在册提交信息＋票面双源为准；若需第三方复跑可按票面 staged 二进制载体重测（收件腿职责）。 |
| AC#5 纪律格 | **PASS** | 只读 git 现核：①两枚 388 提交的完整文件清单——wave121 `28763663e`（7 件全在 internal/audio/aecspike）、closeout-N `821f4a2d3`（cmd/wisp 4 件＋票 388 进度日志 2 行＋163 家族 3 件〔同提交的另一票，非本票越界〕）——`git diff --name-only 28763663e~1..821f4a2d3` 中 **零条**命中 `internal/risk|internal/session|thresholds.go|golden|PLAN.md|docs/specs`。②工作树定向 `git status --porcelain` 对上述禁改路径＝**空读数**（干净）；`internal/observe/thresholds.go` 最后触碰＝00bbb76fe（票 66 旧波，非本票）；`internal/risk`/`internal/session` 最后触碰＝abf2a35a6（票 123 旧波）。③closeout 腿只 commit 未 push（`origin/dev..HEAD` 核），票面 Status/-done 未动。④本验收腿自身零写入（除本证据件）。 |

## 判据与读数复跑口令（收件腿可逐条重放）

1. 接线核：`grep -n "attenuateAnnounce388" cmd/wisp/resident_voice_interrupt_381r2_windows_test.go cmd/wisp/resident_voice_interrupt_390r1_windows_test.go` → `:465` / `:165`，均在各自首 `rv.announce` 之前。
2. 整型路核：`grep -rn "SyscallN" internal/audio/aecspike` → 仅 `rig_windows.go:124`（vtable 整型路）；音量写路 `grep -n "aecsmoke_set_client_volume" internal/audio/aecspike/dmo_cgo_windows.go` → `:92`，C 侧 `aecsmoke.c:365`。
3. 三读数核：`git log -1 --format=%B 821f4a2d3` → ATTENUATE original=1.000000 attenuated=0.150000 restored=1.000000（两测试）。
4. loopback 对照核：`git log -1 --format=%B 28763663e` → session 0.6 peakrms=750.3 vs 0.15 peakrms=46.5, ratio 0.0619。
5. 禁改路径核：`git diff --name-only 28763663e~1..821f4a2d3 | grep -E "internal/risk|internal/session|thresholds.go|golden|PLAN.md|docs/specs"` → 空输出。

## 遗留（非本腿欠账）

- **AC#2 的机主亲耳确认**：需机主在场听一次衰减后的播报，由收件腿/机主签收——待机主。
- AC#4 读数系 closeout 腿真机重测的在册读数；本腿硬约束禁跑 go test，未做独立复跑（如需，收件腿按票面 staged 载体重放）。
- 本腿只写本证据件，未改产码/票面/Status/-done，未 commit/push。
