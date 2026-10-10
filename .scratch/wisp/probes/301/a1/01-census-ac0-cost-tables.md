# 301-a1 — `AC#0` 只读代价普查（票 301 · GUARD D 射程 · `internal/audio` windows-tagged 名册）

腿＝`301-a1`；钟＝`2026-10-10 13:2x +08`（见 `00-anchor.md`）；锚＝HEAD `772ff880`，分支 `dev`。
**射程声明（本件每一枚"数"都适用）**：一律量 **HEAD blob**（`git show HEAD:<path>`），⛔ 工作树。
起手 porcelain 尺（`git status --porcelain -- internal cmd scripts .github`）读数 **2 枚**：
` M internal/audio/wasapi_windows.go` ＋ `?? internal/audio/parse_wave_format_300_windows_test.go`
＝票 300 落地腿的在飞面。**那枚未跟踪件不在本件任何分母里**（blob 上没有它）。
本程 ⛔ 任何产码、⛔ 任何 `go build|vet|test|run`；`go` 侧只用过 `git ls-tree`／`git show`／`grep`／`sed`。
★**中途 HEAD 动过一枚**（`772ff880` → `6d81638b`，出自编排者自己那笔立票 commit）⇒ 本件所有 blob **在两枚锚下逐字相同**，
读数在新 HEAD 上重跑过并复现（尺与 diff 回显见 `02-boundary-selfproof.md` §1）。

---

## 表① 会被拉进分母的枚数与名字

**枚数＝9 枚顶层用例，分布在 2 枚文件里。**

尺（⛔ `go list`，那算编译面）：

```
# R1 名册：HEAD blob 下 internal/audio 的全部 _test.go
git ls-tree -r --name-only HEAD -- internal/audio | grep '_test\.go$'
# R2 tag 判据：只看每枚文件前 3 行
while read -r f; do c=$(git show "HEAD:$f" | head -3 | grep -c 'go:build windows'); echo "$c $f"; done < <(R1)
# R3 用例名：对 tagged 文件逐枚
git show "HEAD:$f" | grep -n '^func Test'
```

R1 射程＝`internal/audio/`（blob），**7 枚** `_test.go`：
`capturelevel_windows_test.go` `captureopt_test.go` `gate_test.go` `hotplug_test.go` `level_test.go` `resample_test.go` `wavinjector_test.go`

R2 → **2 枚 tagged**（blob 前 3 行含 `//go:build windows`）：

| 文件 | tag 行 | 按文件名分档的尺会不会漏 |
|---|---|---|
| `internal/audio/capturelevel_windows_test.go` | blob:1 `//go:build windows` | 不漏（名字带 `_windows`） |
| `internal/audio/hotplug_test.go` | blob:1 `//go:build windows` | **会漏**（⛔ `_windows` 后缀，票面 `现量` 第 3 条那枚坑，本程实测复现） |

⇒ **与票面 `现量`（"七枚，其中 2 枚 tagged"）对上了。**

R3 → 9 枚顶层用例（列名带 blob 行号）：

| # | 用例名 | 所在文件 @ blob 行 |
|---|---|---|
| 1 | `TestAC247RealCaptureLoopEmitsLevels` | `capturelevel_windows_test.go:19` |
| 2 | `TestHotplugReenumerateOnce` | `hotplug_test.go:199` |
| 3 | `TestHotplugStaleHandleFailsLoudly` | `hotplug_test.go:260` |
| 4 | `TestHotplugReenumerateFailureNamesDevice` | `hotplug_test.go:334` |
| 5 | `TestOpenOccupiedAndPermissionDenied` | `hotplug_test.go:378` |
| 6 | `TestEndpointsPairQueryable` | `hotplug_test.go:427` |
| 7 | `TestPinnedThreadStable10s` | `hotplug_test.go:443` |
| 8 | `TestLiveWasapiSmoke` | `hotplug_test.go:525` |
| 9 | `TestCaptureLoopWavIntegrity` | `hotplug_test.go:573` |

⚠ 口径补一句：`hotplug_test.go` 里另有 15 枚顶层 `func`（`fakeWatcher`／`fakeStream`／`fakeOpener` 方法＋`makeChunks`／`toneSamples`／`drainFrames` 等 helper），`^func Test` 那把尺不收它们；`t.Run` 子用例**没数进 9**（本票分母口径＝票面点名的"顶层用例名"）。
`AC#1` 之后进入 `test-windows` 分母的**包**枚数＝1（`./internal/audio/` 从 ⛔ 在 windows 档变 ✅ 在），**用例**枚数＝9。
票 300 正要交进来的那一枚（`parse_wave_format_300_windows_test.go`）**不在 9 里**（blob 上没有），但它落进同一枚包 ⇒ 同一笔 `AC#1` 之后它会一并开口——这一格的代价要连票 300 一起算，**归编排者**。

---

## 表② 每一枚用例的前置条件（全部〔读码推的〕，⛔ 跑过任何 go 测试）

尺＝逐枚 `git show HEAD:<path>` 读正文；env 判据尺＝
`for f in $(R1); do git show "HEAD:$f" | grep -n 'WISP_LIVE_MIC' | sed "s|^|$f:|"; done` → **3 行命中，全在 `hotplug_test.go` 523/526/527**；
skip/设备面判据尺＝`git show "HEAD:$f" | grep -nE 't\.Skip|t\.Fatal|Getenv|Live|Enumerat|EnumAudioEndpoints|LockOSThread|syscall\.LoadDLL|NewLazyDLL'`；
盘上侧证（〔盘上现量〕，只读 blob）＝`git show HEAD:.github/workflows/ci.yml | grep -n 'WISP_LIVE_MIC'` → **rc=1，零命中** ⇒ CI 从不设这枚 env（只引变量名，无值）。

| # | 用例 | 读 `WISP_LIVE_MIC`？ | 要真麦克风／活设备？ | 托管 `windows-latest` 上的预判 |
|---|---|---|---|---|
| 1 | `TestAC247RealCaptureLoopEmitsLevels` | ⛔ | ⛔（`fakeWatcher`+`fakeStream`+`fakeOpener` 注入枚举缝；文件自述 `:13-18` "the enumeration seam is injected, the audio-data seam never is"＝注入的是 DC 常值样本） | **求值并 PASS**。唯一 env/文件面＝注入 registry（`observe.NewRegistry()`），⛔ 磁盘 |
| 2 | `TestHotplugReenumerateOnce` | ⛔ | ⛔（全 fake；断言"变更后恰好重开 1 次"） | **求值并 PASS**〔读码推的〕 |
| 3 | `TestHotplugStaleHandleFailsLoudly` | ⛔ | ⛔（`fakeOpener` 返回占用/权限错的假流；断言错误类＝`audio_device`＋点名设备＋带指引） | **求值并 PASS**〔读码推的〕 |
| 4 | `TestHotplugReenumerateFailureNamesDevice` | ⛔ | ⛔（同上，重枚举失败路径用 fake 注入） | **求值并 PASS**〔读码推的〕 |
| 5 | `TestOpenOccupiedAndPermissionDenied` | ⛔ | ⛔（表驱动 `t.Run`，占用/ACCESS_DENIED 由假 opener 造） | **求值并 PASS**〔读码推的〕；⚠ 子用例未计进 9 |
| 6 | `TestEndpointsPairQueryable` | ⛔ | ⛔（`&fakeWatcher{devs:...}` ＋ `&fakeOpener{}`，断言 `Mic X`/`Fake Speakers` 一对端点可查） | **求值并 PASS**〔读码推的〕——名字里的 "Endpoints" ⛔ 指真 COM 枚举 |
| 7 | `TestPinnedThreadStable10s` | ⛔ | ⛔（fake 流；真设备面＝无） | **求值，10 秒实跑**：`testing.Short()` 那一支（blob `:444-446`）在 CI ⛔ 触发——尺＝`git show HEAD:scripts/portable-tests.sh`/`HEAD:.github/workflows/ci.yml`/`HEAD:tools/d22scan/runtests.sh` 三份各 `grep -c -- '-short'` ⇒ **0/0/0**，`portable-tests.sh:694` 那行实调只有 `-count=1 -skip` ⇒ `-short` 从不传 ⇒ 这枚**不会 skip，会占 10s**。风险⛔ 是"缺音频端点"，是**时长**与 `LockOSThread` 独占断言（`:504` 若有别的 goroutine 落到那枚线程即 `t.Fatalf`）〔读码推的〕 |
| 8 | `TestLiveWasapiSmoke` | ✅ **唯一一枚**（blob `:526` `if os.Getenv("WISP_LIVE_MIC") != "1" { t.Skip(...) }`，skip 文案逐字带 `真机冒烟待票 16`） | ✅ **要**（`:529` `NewWASAPIMicrophone()` ＝**真** WASAPI COM 路径，不经任何注入缝） | ⚠ **今天不可能被求值**（见表③）。⚠⛔ 别把它读成"到时会 skip"：本仓口径 **SKIP ⛔ 是绿**——`tools/d22scan/runtests.sh:98-102`（blob）逐字＝skipped≠0 即 `exit 1`（`:24` 注释 "2. any `--- SKIP` → SKIP is NOT a pass"）。⇒ 若 `AC#1` 只动清单而 ⛔ 同时把 `:593` 那行留在 ledger 里，这枚会在 `test-windows` 上**以 SKIP 把该步判红**。它今天之所以还没红，只因它在 ledger 的 `-skip` 名单里、而那枚包⛔ 在 windows `scope` 里（互相抵消）。**这一格是本票最贵的一格，具名归编排者裁（见 §④）** |
| 9 | `TestCaptureLoopWavIntegrity` | ⛔ | ⛔（样本由 `sineI16At` 现场生成、写进 `t.TempDir()+"/fixture.wav"` 再读回来过生产 `parseWav`；仓内 ⛔ 依赖任何 `testdata/*.wav`——尺＝tagged 两份 `grep -nE 'os\.Open\|ReadFile\|testdata\|\.wav\|TempDir'` ⇒ 仅 `:575`/`:577` 两行，均在 TempDir 下） | **求值并 PASS**〔读码推的〕。⚠ 与票 300 在飞那枚**共读同一枚 `parseWav` 面** |

**小结（代价的形）**：9 枚里 **8 枚设备无关**〔读码推的〕——它们今天不跑，⛔ 因为缺硬件，只因**没有档求值它们**；
唯一带真设备前置的是第 8 枚，而它**已经**在 `:593` 的 ledger 里有名有姓。⇒ `AC#1` 的真实代价 ≈ 8 枚新判进 `test-windows` 分母 ＋ 1 枚靠 ledger 挡着的 opt-in，⛔ "9 枚都会因为托管 runner 没麦克风而红"。

---

## 表③ `scripts/portable-tests.sh:593` 那行夹具 ledger 今天到底 inert 与否

blob `:593` 逐字起头：`"TestLiveWasapiSmoke|./internal/audio/|windows|fixture|live WASAPI capture needs WISP_LIVE_MIC=1 and a physical microphone (hotplug_test.go:527, //go:build windows); a hosted runner has no audio endpoint, so the case has no subject there. Real-hardware smoke is ticket 16's"`

**判定：这枚用例在当前任何一档里都没有被"求值（执行）"的可能；但这行 ledger 本身⛔  inert——它的陈旧性检查今天真的在 windows job 上跑。**
两件事必须分开写：

**甲・用例被求值的可能＝无**〔读码推的〕。尺与链条：
1. `:248-253` windows `scope` 十枚路径里没有 audio（表①/附格已逐枚对拉）；
2. `:241` `core` 的 `scope` 里有 `./internal/audio/...`，而 `core` job＝`test-core:` + `runs-on: ubuntu-latest`（`ci.yml:402-403` blob 现量）；
3. 实调只有 `:694` `sh "$strict" "${scope[@]}" -count=1 -skip "$skip_pattern"`——**跑的只有 `scope`，⛔ `ledger_pkgs`**；
4. `internal/audio` 的 2 枚 tagged 文件在 GOOS=linux 上根本不在编译面里 ⇒ **9 枚用例（含第 8 枚）在 CI 上从未求值一次**，与票面 `现量` 第 1/2 条同结论。

**乙・这行 ledger ⛔ 是死行**〔读码推的〕，三条理由，各有尺：
- **平台门**：`:639-642` `case $platform in any|"$goos") ;; *) continue ;;` ⇒ 在 `ubuntu` 的 core/census 上 `platform=windows` 被 `continue` 掉（这行今天**只在 `--scope=windows` 那一 job 上被求值**）；
- **陈旧性检查按包、⛔ 按 scope**：`:651` `listed=$(go test -list "^${name}\$" "$pkg" ...)` 里的 `$pkg` 是从 ledger 行自己取的（`./internal/audio/`），与 `scope` 无关；`:603-605` 那段注释逐字写着 "the scope **PLUS** every ledger package, so an entry whose package is outside a narrowed scope is still checked rather than silently spared"（blob `:616` 的 `-list` 宇宙同理）。⇒ **改名／删用例／再塞一枚 build tag，`test-windows` 那步会红（`:666-673` 的 `stale` 腿 `exit 1`）**。这条与 `internal/risk` 那行（`:596`）的形状不同：risk 行的 ci.yml 说明（`ci.yml:769` 起）明确写了"Widening the scope to all of internal/risk also moved 201 previously windows-untested assertions into the denominator"——**risk 是"行＋档"配成对交的，audio 这枚只交了行、⛔ 档**；
- **`-skip` 腿今天打空**：`active` 收了它 ⇒ `:686` 把它折进 `skip_pattern` ⇒ 但 `:694` 那发 ⛔ 不含 audio 包 ⇒ **这条 `-skip` 分支恒不命中**＝就 `-skip` 这一个用途而言是 no-op（正是 `:669-671` 自己警告的 "a -skip pattern that matches nothing is the green no-op"），只是脚本对这种 no-op ⛔ 判红。

**"登记的 `windows` 档归属与实际清单（`:248-253`）是否矛盾"＝是，矛盾，且矛盾⛔ 致命。**
`|./internal/audio/|windows|` 这三段合起来读出来的事实是"这枚用例属于 windows 档、由 windows 档的 job 负责记账"；盘上事实是 windows 档十枚路径⛔ 收 audio，而**唯一**收 audio 的 `core` 档跑在 ubuntu。⇒ 这行是一枚**跨档悬挂的记账**：账（ledger）在 windows、货（scope）在 core、货在 core 上又⛔ 编得进。它今天**没造成任何假绿**，因为 GUARD B（`:708-747`，逐包要求 `ok\t<import path>\t…s` 那行结果）只对 `scope` 里的包说话，而 audio ⛔ 在 windows `scope` 里 ⇒ audio 的"一次没跑"⛔ 任何门看得见。**这正是票 301 的题面，⛔ 这行 ledger 的错**——它是本票唯一一处"仓里已经按 audio 会在 windows 档那个形状写过字"的地方（票面 `现量` 第 6 条口径）。

**那枚行号锚 `hotplug_test.go:527`（HEAD blob 上现量）**

尺（逐字）：
```
git show HEAD:internal/audio/hotplug_test.go | cat -n | sed -n '523,531p'
git show HEAD:internal/proc/jobscope_windows_test.go | cat -n | sed -n '85,90p'   # 同表 :592 的锚
git show HEAD:internal/memory/concurrent_test.go | cat -n | sed -n '184,189p'     # 同表 :591 的锚
cat -n internal/audio/hotplug_test.go | sed -n '523,531p'                          # blob 与工作树对拉
```
- **行号对上了**：blob `:527` 逐字＝`t.Skip("live WASAPI smoke requires WISP_LIVE_MIC=1 and a real microphone (真机冒烟待票 16)")`。
- **内容＝同一件事**：ledger 那句理由讲的就是"WISP_LIVE_MIC=1 ＋ 一枚真麦克风"，blob `:527` 那枚 `t.Skip` 的文案逐字就是这两件。
- **口径旁证**：同表另两行的锚都指向各自文件里的 **`t.Skip` 那一行**（`jobscope_windows_test.go:87`＝`t.Skip("helper process mode not set")`、`concurrent_test.go:186`＝`t.Skip("crash-writer subprocess; …")`），⛔ 指向 `func` 声明行（`:591` 那枚的 `func` 在 `:184`）。⇒ `:527` 指向 skip 行而 `func TestLiveWasapiSmoke` 在 `:525`，**是这行的既有惯例、⛔ 一枚漂移**。（本格按票面 `与票 255 的关系` 归口票 255 `AC#6` 的既有形，本票⛔ 新造判据。）
- **blob／工作树无漂**：第 4 把尺回显与 blob 逐行一致（porcelain 里 `hotplug_test.go` ⛔ 出现在改动的 2 枚中）⇒ 在飞那枚票 300 的活**没碰**这枚文件。
⇒ **本票 ⛔ 供一枚腐烂锚。票面 `现量` 第 6 条那句"到底 inert 与否、⛔ 现在当成已知"这里交了：锚活，行 inert（指 `-skip` 那条腿），账活（指 `-list` 那条腿）。**

---

## 附格 — 改完之后同一笔 commit 必须一起动的清单（`AC#1` 的写面）＋今天各是多少枚

| 面 | 今天枚数 | 尺（射程＝`scripts/portable-tests.sh` 的 **HEAD blob**） | 名字 |
|---|---|---|---|
| `windows)` 档 `scope` 清单 | **10 枚路径** | `git show HEAD:scripts/portable-tests.sh \| sed -n '249,253p' \| tr -s ' ' '\n' \| grep '^\./' \| sed 's\|/\$//' \| sort \| wc -l` | `./cmd/llmrecord ./internal/ball ./internal/config ./internal/perm ./internal/plugin ./internal/proc ./internal/projctx ./internal/risk ./internal/secret ./internal/session` |
| `win_pin` 正文 | **10 枚导入路径** | `git show HEAD:scripts/portable-tests.sh \| sed -n '193,204p' \| grep 'github.com' \| sed 's\|github.com/CarlosShao/wisp/\|./\|' \| sort \| wc -l` | 同上 10 枚（前缀换成导入路径形） |
| **两把尺逐枚对拉** | **同形（差 0 枚）** | `diff /tmp/301a1_tier.md /tmp/301a1_pin.md` → **无输出，rc=0** | — |

`win_pin` 正文在 blob `:193` 起（`:193 win_pin='` … `:204 '`）；`windows)` 档在 `:248-253`，`:254 pinned=$win_pin`。
⇒ `AC#1` 那一笔**必须同时**给这两处各加一枚（`./internal/audio/` 进清单、`github.com/CarlosShao/wisp/internal/audio` 进 pin），因为 `:380` 那条 census 认领判据读的是 **pin**、`:694` 实调读的是 **scope**——只改一处＝一边假认领、一边漏跑。
`internal/audio` 在 blob 里现共 **3 处**引用（尺＝`git show HEAD:scripts/portable-tests.sh | grep -n 'internal/audio'`）：`:168`（`core_pin`）、`:241`（`core` 的 scope）、`:593`（ledger）。**windows 档与 `win_pin` 两处都⛔ 有它**——这就是本票要收的那道缝。

原文那句要求（`AC#1` 引用时⛔ 写成一行"逐字"，**它跨两行**，blob `:417` 尾 ＋ `:418` 头）：

```
:417  echo "portable-tests.sh: zero coverage for four days after their tests landed. Pull the package"
:418  echo "portable-tests.sh: into a named scope and update that tier's pin in the SAME commit, or"
```

尺＝`git show HEAD:scripts/portable-tests.sh | cat -n | sed -n '417,418p'`。两行各自都以 `echo "portable-tests.sh: ` 起头（与票面一致）；拼起来读是 "… Pull the package into a named scope and update that tier's pin in the SAME commit, or …"，**⛔ 是盘上的一行**。

---

## 我跑的每一把尺（逐字，可重跑；射程＝HEAD blob 除注明外）

```
date '+%Y-%m-%d %H:%M:%S %z'
git rev-parse --short HEAD
git status --porcelain -- internal cmd scripts .github
git rev-parse --abbrev-ref HEAD

git ls-tree -r --name-only HEAD -- internal/audio | grep '_test\.go$'
while read -r f; do c=$(git show "HEAD:$f" | head -3 | grep -c 'go:build windows'); echo "$c $f"; done < /tmp/301a1_files.md
git show HEAD:internal/audio/capturelevel_windows_test.go | sed -n '1,6p'   # tag 头
git show HEAD:internal/audio/hotplug_test.go | sed -n '1,6p'                 # tag 头
git show "HEAD:<file>" | grep -n '^func Test'
git show "HEAD:<file>" | grep -n '^func '
for f in $(git ls-tree -r --name-only HEAD -- internal/audio | grep '_test\.go$'); do git show "HEAD:$f" | grep -n 'WISP_LIVE_MIC' | sed "s|^|$f:|"; done
git show HEAD:internal/audio/hotplug_test.go | grep -nE 't\.Skip|t\.Fatal|Getenv|Live|live|Enumerat|EnumAudioEndpoints|eDevice|NewLazyDLL|syscall\.LoadDLL'
git show HEAD:internal/audio/capturelevel_windows_test.go | cat -n
git show HEAD:internal/audio/hotplug_test.go | cat -n | sed -n '420,540p'
git show HEAD:internal/audio/hotplug_test.go | cat -n | sed -n '571,600p'
git show HEAD:scripts/portable-tests.sh | cat -n | sed -n '185,260p'
git show HEAD:scripts/portable-tests.sh | cat -n | sed -n '355,400p'
git show HEAD:scripts/portable-tests.sh | cat -n | sed -n '410,425p'
git show HEAD:scripts/portable-tests.sh | cat -n | sed -n '575,600p'
git show HEAD:scripts/portable-tests.sh | cat -n | sed -n '600,690p'
git show HEAD:scripts/portable-tests.sh | cat -n | sed -n '690,760p'
git show HEAD:scripts/portable-tests.sh | grep -n 'ledger\|LEDGER'
git show HEAD:scripts/portable-tests.sh | grep -n -- '-short\|go test'
git show HEAD:scripts/portable-tests.sh | grep -n 'internal/audio'
git show HEAD:tools/d22scan/runtests.sh | grep -nE 'SKIP|skip|short|exit |PASS'
for f in scripts/portable-tests.sh .github/workflows/ci.yml tools/d22scan/runtests.sh; do git show "HEAD:$f" | grep -cn -- '-short'; done
git show HEAD:.github/workflows/ci.yml | grep -n 'WISP_LIVE_MIC'
git show HEAD:.github/workflows/ci.yml | cat -n | sed -n '395,415p'
git show HEAD:.github/workflows/ci.yml | cat -n | sed -n '760,800p'
git show HEAD:scripts/portable-tests.sh | sed -n '249,253p' | tr -s ' ' '\n' | grep '^\./' | sed 's|/$||' | sort | tee /tmp/301a1_tier.md
git show HEAD:scripts/portable-tests.sh | sed -n '193,204p' | grep 'github.com' | sed 's|github.com/CarlosShao/wisp/|./|' | sort | tee /tmp/301a1_pin.md
diff /tmp/301a1_tier.md /tmp/301a1_pin.md
```

各尺 rc：`date`/`rev-parse`/`porcelain` rc=0；`git ls-tree|grep` rc=0（7 行）；tag 尺 rc=0；
blob `grep -n '^func Test'`：capturelevel rc=0（1 行）、hotplug rc=0（8 行）；
`WISP_LIVE_MIC` 名册尺 rc=0（3 行）；
`ci.yml | grep -n 'WISP_LIVE_MIC'` **rc=1**（零命中，＝判据本身）；
三份文件 `-short` 计数尺：`grep -cn` 逐份回显 0/0/0 ⇒ 末次 rc=1（零命中，＝判据本身）；
tier/pin 计数尺 rc=0（10 与 10）；`diff` **rc=0**（同形）；
`grep -n 'internal/audio'` rc=0（3 行：`:168` `:241` `:593`）。

---

## ⛔ 我能裁的格（具名归编排者）

1. **表② 第 8 枚（`TestLiveWasapiSmoke`）在 `AC#1` 之后的实际颜色**：我只证了"它在 ledger 的 `-skip` 名单里"与"`runtests.sh` 把 SKIP 判红"两条静态事实；
   **"改清单后这一发到底 SKIP=0 还是 SKIP=1"必须真跑一发 `bash scripts/portable-tests.sh --scope=windows` 才裁得掉**——那是 go 编译面（本票 `AC#0` 逐字⛔），且票 300 的落地腿正在同一枚包上有未提交面。缺的东西＝① 一次 `--scope=windows` 的 `-v` 名册前后作差（票面 `AC#2` 那格本来就是干这个的），② 一台能跑 WASAPI 的 host（若判"搬进 ledger 不成立、必须真求值"）。
2. **表② 里 7 枚"求值并 PASS"全是〔读码推的〕**：⛔ 一枚我都没跑。`AC#2` 的对照读数（改前/改后各一发、逐名作差）是它们唯一的凭据来源。
3. **表① 的分母该⛔ 该扩到 `t.Run` 子用例**：本件按票面口径只数 `^func Test`（9 枚）。若 `AC#1` 的钉要用例级名册，`TestOpenOccupiedAndPermissionDenied`（blob `:378-423`）那枚表驱动会再贡献若干子名——**枚数没算，归编排者裁**。
4. **票 300 那枚在飞件对分母的增量**：blob 上没有它，我⛔ 能把它的用例名算进"会被拉进分母的枚数"。等它交了、或编排者按工作树量一次再说。
5. **GUARD D 用例级盲区要不要收口（票面 `AC#3`）**：本格由编排者裁形，本件⛔ 动门。

## 派单/票面我认为写错或要收紧的地方（必答格）

1. **票面 `:15` 那句"会同时把该包其余 windows-tagged 用例（上面那 2 枚，以及票 300 正要交进来的那一枚）一起拉进分母"——枚数口径错位。** "上面那 2 枚"在 `:11` 指的是**文件**，到 `:15` 被读成**用例**。blob 现量＝**2 枚文件 ＝ 9 枚用例**（不含票 300 那枚）。`AC#0` 的产物就是这一枚数，写歪了正好处在"代价面"那句上。派单 §2 表① 的措辞（"顶层用例名"）⛔ 有此歧义，我以派单为准并在此具名报回。
2. **派单 §2 表② 那句"要不要真麦克风／活设备"对第 6 枚（`TestEndpointsPairQueryable`）是个假靶**：它名字里带 `Endpoints`、正文却全用 `fakeWatcher`/`fakeOpener`（blob `:428-437`），⛔ 碰真枚举。谁按名字把它归入"要设备"就会把代价读高 1 枚。
3. **派单表③ 让我"回答它声称的那枚用例有没有被求值的可能"，预设了这行可能整体 inert——读数是有条件的⛔**：`-skip` 那条腿今天打空（inert），`-list` 陈旧性那条腿今天真的在 windows job 上求值（⛔ inert，`AC#1` 之前改名就红）。把这行整体判成死行会让人在 `AC#1` 里顺手删它，删掉之后第 8 枚会以 SKIP 把 `test-windows` 判红。**这一条是本件最需要被读到的读法。**
4. **`ci.yml:400-401` 那段注释是一枚与本票同根的过期断言**（blob 逐字："Windows-only packages (ball GUI, cgo speech) are out of this job's scope by platform, not skipped: they run in test-windows / slo jobs."）：`internal/audio` 有 2 枚 windows-tagged 文件、被 `core` 认领、跑在 ubuntu ⇒ 它 windows 那一半既⛔ "由平台出档"（它就在 core 的 scope 里）也⛔ "在 test-windows 里跑"（windows 档⛔ 收它）。票面 `现量` 只点了 ledger 那一行（`:593`），没点这行注释 ⇒ **建议 `AC#1` 的写面外加"这行 yml 注释的措辞"**，但那⛔ 碰 `.github/workflows/ci.yml`、票面 `AC#4` 逐字禁 `AC#1` 那笔碰它 ⇒ **⛔ 我自己动，交编排者裁**（拆第二笔票或明确豁免）。
5. **锚点惯例**：派单表③ 给的尺是"`sed -n '520,535p'` 之类"。我把它钉成"`sed -n '523,531p'` ＋ 同表另两行锚的对照"，因为"行号对上了"这件事**只有在读出惯例（指向 `t.Skip` 行、⛔ `func` 行）之后才可裁**；单看 `:525` vs `:527` 会误报一枚漂移。建议这半句进 `AC#0` 的尺库。
