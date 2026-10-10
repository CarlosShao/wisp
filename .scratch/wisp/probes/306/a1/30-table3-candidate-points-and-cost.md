# 306-a1 · 30 表③：夹具落点候选＋代价（含硬前提的现量）

世代＝`f994f95`。CI 尺名＝`scripts/portable-tests.sh`（scope 表 `:232-263`、pin 名册 `core_pin` 起 `:164`／`win_pin` 起 `:193`／`cli_pin` 起 `:206`）＋`.github/workflows/ci.yml`（`test-core` `:419-420`＝`runs-on: ubuntu-latest`，调用点 `:487`；`test-windows` `:533-534`＝`runs-on: windows-latest`，调用点 `:761` census／`:797` windows）。

## 硬前提（派单点名"你自己现量这一条再写进表 ③"）＝**已量实，票面那句成立**

| 尺 | 读数 |
|---|---|
| `git show HEAD:internal/audio/wavinjector.go \| grep -n "go:build"` | **零命中**（`rc_tag=1`）⇒ 平台中立件 |
| `sed -n '1,2p' <git show HEAD:internal/audio/wasapi_windows.go>` | 第 1 行逐字 `//go:build windows`；常量在 `:97/:98/:99` |
| `GOOS=linux GOFLAGS= go build ./internal/audio/` | **rc=0**（件＝`logs/50-linux-build.txt`，自落 `rc_linux_build=0`） |
| `GOOS=linux go list -f '{{range .GoFiles}}…' ./internal/audio/` | linux 档 GoFiles 含 **`wavinjector.go`**＋`wasapi_other.go`（⛔ `wasapi_windows.go`） |
| ★决定性一发（**仓外导出树**＝`git archive HEAD \| tar -x -C C:/Users/swq/tmp/306a1-export`，⛔ 碰共享工作树、⛔ 在仓库目录内、件留该树 `ref-const-linux.txt`／`ref-const-windows.txt`） | 把 `:192` 改成 `fmtTag == waveFormatExt`、`:218` 改成 `fmtTag == waveFormatFloat` 后：`GOOS=linux go build ./internal/audio/` ＝ **rc=1**，逐字 `internal\audio\wavinjector.go:192:17: undefined: waveFormatExt` 与 `:218:17: undefined: waveFormatFloat`；同一棵树 `GOOS=windows go build` ＝ **rc=0** |

⇒ **"改成引用常量"这一形今天会让非 windows 档编译不过**＝实测，⛔ 推断。落地腿⛔ 走这一形（票面 §边界第 3 条据此成立）；本票买的是**行内值也有牙**。

## 候选甲＝新增一枚 wav 夹具文件（`internal/audio/testdata/extensible-float32.wav` 之类）

- 动 `third_party`？**⛔ 用不着**（夹具是纯数据，⛔ 碰那三枚冻结 `.dll`；派单/票面边界亦⛔ 动）。
- `//go:build windows`？**⛔ 不需要**（`parseWav` 无标签，读夹具的测试也⛔ 需要标签）⇒ 代价是"⛔ 需要标签"带来的**双档都跑**（见下）。
- CI 哪一档跑到：**两档**。core 档 scope 逐字含 `./internal/audio/...`（`portable-tests.sh:242`）跑在 ubuntu（`ci.yml:420/487`）；windows 档 scope 含 `./internal/audio/`（`:253`）跑在 windows-latest（`ci.yml:534/797`）。夹具⛔ 依赖任何设备/PATH。
- 新增文件要不要动 pin？**⛔ 要不动**：`core_pin`（`:164` 起）与 `win_pin`（`:193` 起）里 `github.com/CarlosShao/wisp/internal/audio` 都是**包名**行（分别 `:168`／`:195`，本腿现量），⛔ 测试名名册。★票面 `:25` 那句"CI 上真会跑（凭据＝`probes/300/orch/r10-ci-third-run-new-test-visibility.txt`，run `38049322056`）"＝**编排者已量、本腿⛔ 复跑**（复跑要 push，越界）。
- 先例成本（这是甲的真代价）：本仓**被跟踪的 `.wav` ＝ 0 枚**（尺＝`git ls-files | grep -icE '\.wav$'`＝0），且 `internal/audio/` **没有 `testdata/` 目录**（尺＝`find internal cmd -type d -name testdata` ⇒ 只有 `internal/agent`、`internal/llm`、`internal/models`、`cmd/wisp` 四枚）。⇒ 甲会同时开两处先例（第一枚二进制夹具资产＋audio 的第一个 testdata），且 `wavinjector_test.go` 家族今天的习惯是 `t.TempDir()`＋in-code 字节。
  ⛔ 阻塞风险：`.gitignore` **⛔ 含 `wav` 规则**（尺＝`grep -nE "wav" .gitignore` 零命中）⇒ 不会被静默跳过。
- 能不能同时买到两枚洞的牙：**能**（前提＝那枚 `.wav` 的头部是 extensible＋float32，即 tag `0xFFFE`／`cbSize 22`／size 40／SubFormat 低字 `3`、bits 32）。⚠ 但字节写死在二进制里＝**改期望侧要重新生成文件**，`AC#1` 凭据形状（"改 `:218` 的 3→2 必须红"）核对时，"夹具是否还合规格"本身要另一把尺（`cmp` 或头部回读测试），甲比乙多这一层。

## 候选乙＝在既有 injector 用例里加形（in-code 40 字节 extensible+float32 面，同包 `_test.go`，⛔ 新标签）

- 动 `third_party`？⛔ 不动。
- `//go:build windows`？⛔ 不需要（`GOOS=linux go vet ./internal/audio/` ＝ **rc=0**，件 `logs/51-linux-vet.txt` ⇒ 无标签测试文件在 linux 档连**测试侧**都能类型检查；今天 audio 包里那两枚 300 钉反而**带**标签：`parse_wave_format_300_windows_test.go:1` 逐字 `//go:build windows`）。
- CI 哪一档跑到：**core(ubuntu)＋windows 两档**（同上 `:242`／`:253`）；⛔ pin 变更（包名名册）。
- 现成形状可搬（⛔ 新算布局）：`parse_wave_format_300_windows_test.go:59-68` 的构造器已把这份 40 字节布局逐字节注释过（`:59` `	le.PutUint16(b[0:], 0xFFFE)`／`:65` `	le.PutUint16(b[16:], 22)`／`:68` `	le.PutUint32(b[f.data1At:], f.data1)`），`:99` 的 `bits: 32, data1: 3, data1At: 24` 正是本票要的形；乙＝把它加上 `RIFF/WAVE/fmt /data` 外壳＋32-bit 样本，喂 `NewWavInjector` 或 `parseWav`（旁路先例＝`hotplug_test.go:581`）。
- 能不能同时买到两枚洞的牙：**能，且一枚夹具一次买两枚（还附赠 `:196`）**。推演依据＝本腿**逐字读到的分支结构**＋**块级 hit 计数**（⛔ 产码，所以"改⇒红"那两发留给 AC#1/AC#2）：
  - extensible+float32 面先经 `:192`（真）→ `:196` 把 `fmtTag` 重写成 SubFormat 低字 `3` → `:218` 命中 → 走 `:224` `FloatToPCM16`。
  - `:192` 的 `0xFFFE`→`0xFFFD` ⇒ 条件假 ⇒ `fmtTag` 留在 `0xFFFE` ⇒ `:211`/`:218` 全不中 ⇒ 落 `:226` 报错 ⇒ **解析值断言红**。
  - `:218` 的 `3`→`2` ⇒ `fmtTag==3` 无 case 命中 ⇒ 同样落 `:226` ⇒ **红**。
  ⇒ **票面"补 float32/extensible 夹具会同时钉住行内 `3`/`0xFFFE`"这句＝本腿验证成立**（并额外把 `:196` 的 SubFormat 偏移 `body+24` 一起钉住——今天它也 hit=0，票面⛔ 数到它）。
  ⚠ 附带读数：断言必须钉**解析出来的值**（票面 `AC#1` 原话），⛔ 断"被调用过"；若只加一枚**非** extensible 的裸 `tag=3/bits=32` 面，则**只**买到 `:218` 一枚、`:192` 仍零执行者。

## 给编排者的取舍（⛔ 本腿裁）

两形都能买到两枚牙。乙的边际成本低（⛔ 二进制资产、⛔ 第一个 audio testdata 先例、⛔ 夹具再生的第二把尺），并且有 `:59-68`／`:99` 的现成形状可搬；甲的边际好处＝夹具字节可被别的包/别的腿复用、以及"真文件"更贴近 `NewWavInjector` 的入口形状。两形⛔ 冲突（可先乙、后按需甲）。
