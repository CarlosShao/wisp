# 300-r2 终名册与越界差集（票 300 `AC#6`，落点 C2）

终锚＝`148c0f43`（本节这枚文件是它之后的**第 5 笔**，⚠ 见末行"本件自己那一笔"）。
起手锚＝`05db4bc6`（见 `00-anchor.md`）。

## 门⑦：逐笔名册（尺＝**逐笔** `git show --name-only --format= <sha>`，⛔ 区间尺 `git diff A..B`）

| 笔 | 名册（逐笔全量） |
|---|---|
| `a0339995` | `.scratch/wisp/probes/300/r2/00-anchor.md` |
| `9a442923` | `internal/audio/wave_format_float_300_windows_test.go`、`.scratch/wisp/probes/300/r2/10-shape.md`、`…/logs/gate1-vet.txt`、`…/logs/gate2-targeted-pristine.txt` |
| `cdced117` | `.scratch/wisp/probes/300/r2/20-mutations.md`、`…/30-gates.md`、`…/logs/authority-mmreg.txt`、`…/logs/gate3-full-package-before.txt`、`…/logs/gate3-full-package-after.txt`、`…/logs/gate4-format.txt`、`…/logs/gate5-d22scan.txt`、`…/logs/gate6-census.txt`、`…/logs/gate8-test-counts.txt`、`…/logs/mut-m0-pristine-exporttree.txt`、`…/logs/mut-m1-float-3to1.txt`、`…/logs/mut-m2-float-3to0.txt`、`…/logs/mut-m3-restored-exporttree.txt`、`…/logs/mut-rig-identity.txt`（14 枚） |
| `148c0f43` | `.scratch/wisp/issues/300-parsewaveformat-subformat-offset-reads-two-bytes-past-the-guid.md` |

## 与"允许动的写面"名册作差 ⇒ 越界 **∅**

允许的三枚写面＝① `internal/audio/` 下一枚新增带 `//go:build windows` 的 `_test.go` ② `.scratch/wisp/probes/300/r2/**`（只 `.md`／`.txt`）③ 票面追加一节。
⇒ 上表 19 枚路径⛔ 有一枚落在这三面之外；`internal/audio` 那唯一一枚正是①。

另几把自己加的越界尺（都是只读，逐字可重跑，读数＝本件）：

- 产码⛔ 动过：`git diff --name-only 05db4bc6..HEAD -- internal cmd` 里 `internal/cmd` 只出现 `internal/audio/wave_format_float_300_windows_test.go` 一枚（尺同上表；`wasapi_windows.go`、`wavinjector.go` 均⛔ 命中）。
- 冻结面⛔ 动过：`git diff --name-only 05db4bc6..HEAD -- scripts .github frontend design tools internal/observe/thresholds.go third_party` ⇒ **0 行**。
- 既有那枚凭据本体⛔ 动过：`git diff --name-only 05db4bc6..HEAD -- internal/audio/parse_wave_format_300_windows_test.go internal/audio/wavinjector_test.go` ⇒ **0 行**。
- `.scratch` 里⛔ 落 `.go`、证据件⛔ 叫 `.out`：`find .scratch/wisp/probes/300/r2 -name '*.go' -o -name '*.out' | wc -l` ⇒ **0**。
- 件⛔ 0 字节：`find .scratch/wisp/probes/300/r2 -type f -size 0 | wc -l` ⇒ **0**。
- 件总数（两把各报各的，⛔ 一把合计冒充另一把）：`find .scratch/wisp/probes/300/r2 -type f | wc -l` ⇒ **19**；
  `find .scratch/wisp/probes/300/r2/logs -type f | wc -l` ⇒ **14**（全 `.txt`）；`.md` ⇒ **5** 枚（`00`／`10`／`20`／`30`／本件）。⇒ 14 ＋ 5 ＝ 19，闭合。
  ⚠ 本腿中途在 `30-gates.md` 里写过"件＝`logs/` **14** 枚"，那一枚⛔ 变；两处名册由同一条尺数出来。
- ⛔ push：`git status -sb` 那行 `dev` ⛔ 带 `ahead ... , behind`，且本腿从头到尾⛔ 跑过任何 `git push`（唯一写 ref 的动作是 5 笔 `git commit`）。

## 起手／终了台面差（证明本腿跑期间产码面没漂）

尺＝`git status --porcelain -- internal cmd`（起手 18:38 ⇒ **0 行**；交回 18:5x ⇒ **0 行**）。
⚠ 全仓 `git status --porcelain | wc -l` 起手时＝**834 行**既有 dirt（`design/**` 一批 ` D`、`.gitignore` ` M`、`.scratch/` 一堆未跟踪），⛔ 本腿造的、⛔ 本腿碰的；
本腿所有 commit 都带**显式 pathspec**（⛔ `add -A`／`add .`），所以那 834 行⛔ 进过我的任何一笔名册。

## 本件自己那一笔（具名，⛔ 藏）

门⑦那张表⛔ 含第 5 笔（就是写入本件这一笔）——一枚件⛔ 能列自己的名册。
它的名册只有 `.scratch/wisp/probes/300/r2/40-final.md` 一枚，编排者收件时用**同一条尺**现算即可：
`git show --name-only --format= <第5笔>`。⚠ 这是本腿唯一一处"凭据⛔ 自含"，具名报回。
