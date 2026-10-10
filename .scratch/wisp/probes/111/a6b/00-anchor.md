# 111-a6b 起手锚件 — 票 111 `AC#12` 盲区普查（只读腿，重派第 2 枚）

## 0. 起手现量（同一条命令内取，`2026-10-10 17:04 CST`）

| 项 | 读数 | 尺 |
|---|---|---|
| HEAD | `f76091cd 2026-10-10T17:04:15+08:00` | `git log -1 --format='%h %ad' --date=iso-strict` |
| 产码面脏度 | `0` | `git status --porcelain -- cmd internal scripts .github docs \| wc -l` |
| `wisp.exe` 枚数 | `0` | `tasklist //FI "IMAGENAME eq wisp.exe"` 行数 |
| `balldebug.exe` 枚数 | `0` | 同上形 |

⇒ 起手时 `cmd/ internal/ scripts/ .github/ docs/` 工作树＝HEAD（**快照，不是保证**：另两枚腿在飞，
本腿所有产码面读数一律走对象层 `git show HEAD:<path>`／`git grep … HEAD -- <path>`，不拿工作树当 HEAD）。

## 1. 本腿写面 / 读面自报

- 写面：只在 `.scratch/wisp/probes/111/a6b/**` 新建文件。⛔ 改任何产码／票面／台账／HANDOVER。
- **⛔ Go 编译面：本腿至今跑过 `go build`／`go test`／`go vet`／`go run`／`go list` = 0 次**（`303-a1` 独占 `cmd/wisp` 台面）。
  需要 `go list` 类读数的地方一律具名写成"欠读数"交回，⛔ 自造替代尺冒充 `go list` 的口径。
- `go env` 现量：0 次（本腿不需要 GOOS——分档 OS 是从 `ci.yml` 的 `runs-on` 读出来的，不是问编译器）。

## 2. 我打算用的尺（逐枚＋射程）

| 尺名 | 命令形状 | 射程目录 | 抽样还是整族 | 买到什么 |
|---|---|---|---|---|
| R1 文件类别尺 | `git show HEAD:<f>` 取 blob 前 3 行，读 `^//go:build` 行 | 根模块全部 tracked `_test.go`（`cmd/ internal/`，⛔ `.scratch/**`、⛔ 自带 `go.mod` 的子模块） | 整族 | **问文件内容而非文件名**的分类名册（票面点名的陷阱＝`internal/audio/hotplug_test.go` 那形） |
| R1f 文件名尺（对照用） | 路径尾名 `_windows_test.go` / `_posix…` | 同 R1 | 整族 | 只用来量"按文件名分类那把尺会错几枚"，⛔ 用作名册 |
| R2 顶层用例尺 | `git show HEAD:<f> \| grep -c -e '^func Test'` | 同 R1 | 整族 | 用例枚数（**顶层口径**，⛔ 含子测试；与票 111 `PASS=578` 那类"含子测试"尺差一枚数就具名差） |
| R3 档→runner-OS 尺 | `git show HEAD:.github/workflows/ci.yml` 的 `runs-on` × `--scope=` 调用点 | `.github/workflows/ci.yml` | 整族 | 四档各跑在哪台 OS（票面 ⓐ 那枚"散文变数据"的现量） |
| R4 前置条件尺 | 逐枚取函数体，`grep` 门闸字样（`WISP_LIVE_MIC`／`t.Skip`／`winlive`／`exec.Command`／`net.`／`require` 等），**内容锚** | R1 命中的 tagged 文件 | 整族（named 全列） | 连带名册每枚的前置条件，⛔ 按用例名猜（票面 `A807`② 假靶先例） |
| R5 档认领尺 | 把 R1 的包路径与 `core_pin`／`win_pin`／`cli_pin`／`winsec_pin` 逐枚 `grep -qxF` 对拉 | `scripts/portable-tests.sh` blob | 整族 | "这枚包被哪几档认领"＋"那几档跑在哪台 OS"＝盲区判定 |

## 3. 待交三件（票面 `AC#12` 为准，本派单转述次之）

1. 盲区名册（包／文件／用例三级，逐枚标 R1＋R5）。
2. GUARD D 扩到**用例级** vs 扩到**档级 OS 语义**两形代价表（⛔ 推荐、⛔ 裁形，只给代价）。
3. ★"加一枚路径＝连带拉进分母"的名册与前置条件（票 301 那句"⛔ '改一枚清单'是零成本"的现量）。

## 4. 已读来路（原文，非转述）

- `.scratch/wisp/issues/111-ci-tests-20-of-33-packages.md` `AC#12` 整格＋`:429` 起来历节。
- `scripts/portable-tests.sh`（blob）、`.github/workflows/ci.yml`（blob）、`.scratch/wisp/probes/301/a2/`。
