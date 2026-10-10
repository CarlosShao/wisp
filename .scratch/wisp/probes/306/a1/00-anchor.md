# 306-a1 · 00 锚（只读普查腿 · 起手台面）

腿名＝`306-a1`。性质＝**只读普查**（票 306 `AC#0`：三张表，⛔ 任何产码）。
票面文件＝`.scratch/wisp/issues/306-inline-authoritative-format-tag-literals-in-wavinjector-go-have-no-ruler-at-all.md`（**已逐字读完**，9837 字节；⛔ 改动，⛔ 翻框）。

## 锚（现量，尺＝`git rev-parse`／`git status --porcelain`）

- `git rev-parse HEAD` ＝ `f994f956647ebf343d861b588f8d379f7f883b6e`
- `git rev-parse --abbrev-ref HEAD` ＝ `dev`
- `git status --porcelain -- internal cmd docs | wc -l` ＝ **0**（这三面工作树＝HEAD，所以本腿「工作树」与「HEAD blob」两个名册在此锚点必然同行号）
- ⚠ 与派单的冲突：**派单里没给锚号**，只给了 `d9864baf`（那是编排者四发攻击所用的**导出树世代**，不是当前 HEAD）。以盘上为准＝`f994f95`。

## 票面 `:192` / `:218` 两枚行内字面量的锚点核对（尺＝`git show HEAD:internal/audio/wavinjector.go | grep -n`）

内容锚（逐字，⛔ 行号当权威）：

- `:192` ＝ `			if fmtTag == 0xFFFE { // WAVE_FORMAT_EXTENSIBLE: real tag is SubFormat[0:2]`
- `:218` ＝ `	case fmtTag == 3 && bits == 32:`

⇒ **与派单/票面写的行号一致**（本锚点 `f994f95`）。同族另两枚也在同文件（派单未列，本腿在表 ① 补上）：

- `:211` ＝ `	case fmtTag == 1 && bits == 16:`（`1`＝WAVE_FORMAT_PCM，`mmreg.h:2418`，同一族）
- `:196` ＝ `				fmtTag = binary.LittleEndian.Uint16(data[body+24 : body+26])`（SubFormat 低字读取，偏移 24 是权威值）

## 进程闸门（现量，尺＝`tasklist | grep -icE "<exe>"`，件＝`logs/00-tasklist.txt`）

- `wisp.exe` ＝ **0**
- `balldebug.exe` ＝ **0**
- `mockllm.exe` ＝ **0**
- `msedgewebview2.exe` ＝ **18**（派单已声明＝机主自己的应用，⛔ 归本腿 ⛔ 归编排者处置）

⇒ 闸门通过；⛔ 开真窗。本腿后续只跑包级 `go test ./internal/audio/`（派单点名允许，且声明本腿是当时唯一用 Go 编译面的腿）。

## 台面与射程（每张名册自带，本仓铁律）

- 表 ① 名册射程＝工作树 `internal/ cmd/ tools/ tests/`（全部 `*.go`），锚点 `f994f95`；因 `git status --porcelain -- internal cmd`＝0，同一读数量在 HEAD blob 上逐字同行（blob 尺也跑了：`git show HEAD:internal/audio/wavinjector.go`）。
- 表 ② 调用链尺＝`grep -rn --include=*.go`（同上射程）＋覆盖率尺**本腿自己复跑**：
  `go test ./internal/audio/ -count=1 -coverprofile=.scratch/wisp/probes/306/a1/logs/40-cover.txt`
  → `logs/40-cover-run.txt` 末行逐字 ＝ `ok  	github.com/CarlosShao/wisp/internal/audio	15.960s	coverage: 59.3% of statements`，该件自落一行 `rc_test=0`。
- 表 ③ 硬前提尺＝`GOOS=linux GOFLAGS= go build ./internal/audio/`（件 `logs/50-linux-build.txt`，自落 `rc_linux_build=0`）＋`GOOS=linux go list -f '{{range .GoFiles}}…'`。
- ⛔ 读 `%APPDATA%\wisp-dev\secrets\`；⛔ 读机主 config 值；凭据⛔ 进件。⛔ 写 `frontend/src/**`、`design/**`、`internal/**`、`cmd/**`、`tools/d22scan/allowlist.txt`、`internal/observe/thresholds.go`、`PLAN.md`、D43 表。
- 允许写面＝仅新建的 `.scratch/wisp/probes/306/a1/**` 下 `.md`／`.txt`（已核对：⛔ 新 `.go`，⛔ 文件名以 `.out` 结尾）。临时件只建不删。

## 我这程复跑的既有读数（对拉结论在 `90-final.md`）

- 编排者四发退码件 `probes/300/orch/logs/r3-orch-v4-attacks-rcfixed-20261010-204046.txt`：**本腿未复跑那四发攻击**（那要动产码，⛔ 本票 AC#0 只读），只做了逐字对拉 ⇒ 见 `90-final.md`。
- 验收腿 `300-v4` 的覆盖尺（`parseWav 70.0%`／`Drain 0.0%`／`convertPacket 0.0%`）：**本腿自己复跑并逐字对上**，凭据＝`logs/41-cover-func.txt`（派单点名 ⛔ 引它的数当凭据，所以复跑了）。
