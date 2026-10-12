# `10-impl.md` — 票 310 增量程 `310-r2` 四枚文件的步子原文与逐枚差异数

锚＝`00-anchor.md`（起手 HEAD `73719706`，`dev`）。代码全在提交 `78782c20`。判语归下一程非实现者腿。

## 1. `go test` 调用点与逐枚差异（本腿自己 `git show` 现跑，⛔ 照 r1 名册抄）

| 枚 | `go test` 真起点 | 本腿加／改的块 | `git show --numstat`（加／删） |
|---|---|---|---|
| `tools/d22scan/runtests.sh` | 逐字 `go test -v -count=1 "$@" >"$out" 2>&1`（`:75`＝r1 点名的第五枚真起点；经 `scripts/d22scan.sh:51` 与 CI "D22 扫描正控"步） | 新增独立块起于 **`:90`**（`export GOTMPDIR` `:94`、成功 echo `:95`、`DEGRADED-cannot-create` `:98`、新 `DEGRADED-LOCALAPPDATA-empty` `:103`），位置＝参数校验之后、`out=$(mktemp)` 之前 ⇒ 在任何 `go test` 之前 | **41 / 0**（`@@ -67,0 +68,41 @@`＝纯增量，一行未删） |
| `scripts/winsec-tests.sh` | 零枚直接 `go test`；唯一出向＝委派 `bash "$portable"`（`:128`） | 把合取 `&& [ -n "${LOCALAPPDATA:-}" ]` 拆成嵌套 `if [ -n LOCALAPPDATA ]…else` 让 windows+空 那一支响；成功 echo `:66`、`DEGRADED-cannot-create` `:69`、新告警 `:74` | **22 / 9**（`@@ -54,6 +54,19 @@`＋`@@ -61,3 +74,3 @@`） |
| `scripts/wisp-cli-tests.sh` | 零枚直接 `go test`；唯一出向＝委派 `bash "$portable" --scope=cli`（`:146`） | 同上拆分（PATH 守卫之前）；`:79`/`:82`/新 `:87` | **22 / 9**（`@@ -67,6 +67,19 @@`＋`@@ -74,3 +87,3 @@`） |
| `scripts/portable-tests.sh` | 真起点三处（r1 已列：`-list`／`-list "^<名>$"`／`sh "$strict"`→runtests:75） | ★**本腿⛔ 动它一字**——它 r1 起就带完整的 windows+空-LOCALAPPDATA `else`（现 `:158`），早已"说实话"，是这一形的参照；blob `7bc4e287`⇄HEAD **SAME**，`git diff --name-only` 不含它 | **0 / 0**（对照件，未改） |

`git diff --name-only HEAD -- '*.go'`＝**0 枚**（本腿只动 shell）。

## 2. runtests 根步子＝`winsec-tests.sh:55-65` 逐字照形
派单写死"现形＝`scripts/winsec-tests.sh:55-65` 那一块，逐字照形，只换 `<entry>` 名与 echo 前缀"。本腿：
- 成功行逐字：`runtests.sh: ticket 310 AC#1: GOTMPDIR=$GOTMPDIR - this run's t.TempDir() family lands here; nothing deletes it.`（＝winsec 逐字同，仅前缀 `runtests.sh:` 与根名 `d22scan-`）；
- `DEGRADED-cannot-create` 行逐字同 winsec；
- 守卫＝`if [ -z "${GOTMPDIR:-}" ] && [ "$(go env GOOS)" = windows ]`（**保留 `-z "${GOTMPDIR:-}"`** ⇒ 委托语义不变：`portable-tests.sh` 已 `export GOTMPDIR` 再调 `sh "$strict"`，runtests 见 GOTMPDIR 非空即**整块跳过、不开第二枚根**；`d22scan.sh` 直调 runtests 时 GOTMPDIR 为空 ⇒ runtests 自开一枚 `d22scan-` 根。一次运行只一枚根，由最外层入口具名）。
- 根命名默认 `${WISP_TEST_GOTMP_ROOT:-"${LOCALAPPDATA//\\//}/wisp-test-gotmp"}/d22scan-<yyyymmdd-HHMMSS>-<pid>`＝落 `%LOCALAPPDATA`（与 `%TEMP%` **同卷**，⛔ 换到仓所在卷 `D:`）。
- ★runtests 是 `#!/bin/sh`（＝本机 `/usr/bin/sh` bash-5.2 posix 模式）；`${LOCALAPPDATA//\\//}` 在 posix 模式下用**真文件字节**实测可展开成 `C:/Users/swq/AppData/Local/...`（件 `logs/`；见派单 §2 台面上"逐字照形"的安全性本腿现验过）。`sh -n tools/d22scan/runtests.sh` rc=0。

## 3. 四枚同一形（empty-LOCALAPPDATA 告警正文逐字节相同）
把四枚的告警 echo 拼接、剥去 `<文件>:` 前缀后归一化比对（件 `logs/15-warn-body-normalized.txt`）：
```
DEGRADED - LOCALAPPDATA is empty, so this step has no measured-safe root to build
(a root on another volume turns the winsec seal cases red, <指针短语>). GOTMPDIR stays unset.
```
- 唯一差异＝**文件特定的指针短语**：portable 用 `see the note above`（长注就在它上方），winsec／cli／runtests 用 `see the note in scripts/portable-tests.sh`（长注在它文件里、不在本文件上方）。除 `PTR` 归一后**四枚正文 SAME**（distinct bodies＝1）。指针短语若强行抹平会在三枚委派件里造成"上方有注"的假陈述，故本腿保留该文件相关差、并在件里具名。

## 4. 语义面（⛔ 动了什么／⛔ fail-closed）
- 只加/改那一段独立块：runtests `+41/-0`、winsec/cli 各 `+22/-9`（删的 9＝原合取 `if` 那一行与其后 8 行被重排为嵌套，逻辑等价、仅把 `&& [ -n LOCALAPPDATA ]` 下移为 `if/else`）。
- ⛔ 任何一处删目录；本腿建的根（`%LOCALAPPDATA%\wisp-test-gotmp\d22scan-*` 等）**全留盘上**（名册见 `30-gates.md`）。
- ⛔ fail-closed：新告警只是 stderr 一行；独立一发 `env -u LOCALAPPDATA sh runtests.sh -C tools/d22scan ./...` 打出告警后 **go test 照跑、`rc=0`**（件 `logs/24-tooth2-after-rt-*`）⇒ "响在嘴上、⛔ 响在退码上"（编排者 `A858 §6` 那句）在本腿入口上二次坐实。
- ⛔ 动非 windows 那两支：四枚的 windows 判定都在块顶（`&& [ "$(go env GOOS)" = windows ]`），非 windows ⇒ 整块跳过、行为与改前逐字节一致（portable 的 ubuntu `left as-is` 行本腿⛔ 碰）；票 `Q-310-v1-2`（ubuntu 那一支的实测）仍⛔ 本腿射程（本机 `go env GOOS=windows`，永远走⛔ 到那一支）。
- 机制依据⛔ 重述（r1 已读 GOROOT：`testing.go:1613` `os.MkdirTemp(os.Getenv("GOTMPDIR"), …)`）。
