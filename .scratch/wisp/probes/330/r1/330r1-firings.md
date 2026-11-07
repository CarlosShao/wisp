# 票 330-r1 发记录（G1 对接钉 · cmd/wisp 面）

腿＝330-r1 写腿。日期 2026-11-07 11:02-11:2x +08。锚＝HEAD `766bf71f4`（329 近 15 枚提交零命中、未落地，如实记）。

## 被突变件基线 sha256（突变前现算）

- `internal/agent/testdata/golden/plan-present.sse` = `7d26fa889b59fc3a2b39ff18883ee331f8ff1c06ddb45a5a23bc6820af9079be`（与票面现量 `7d26fa889b59fc3a…` 同头）
- `internal/agent/plan_present_doorbell_165r3_test.go` = `6542773d3d2edcc6d554a47f091ce2dffe35586a92c1f948f79c361454683fca`（与 330-v1 判语 `6542773d…` 同头）

## 正向发（真名与线上名相等 ⇒ 绿）

- 命令：`export PATH="$PWD/build:$PWD/third_party/sherpa-onnx:$PATH" && go test ./cmd/wisp -run 'Test330' -count=1 -timeout 30m -v`
- 读数：`--- PASS: Test330DockPinsTheWireNameToTheAssembledRoster`，rc=0。
- 用例自报读数（`plan_present_roster_dock_330r1_test.go:134`）：
  `对接读数：线上名 "plan.present" ＝ 真名册（22 枚）里的一行；Resident=true`
- 22 枚真名册＝`rt.bridge.Registry().List()`（真装配 `runTextTask → assembleRuntime` 出的那一个），全名单见 G1 发的红句。

## G1 重放发（夹具名＋镜像名一起改 `plan.presentX`、真名册不动 ⇒ 必须红）

- 突变：`sed -i 's/plan\.present/plan.presentX/g'` 打两枚文件——golden 夹具 `:12` 的 `"function":{"name":"plan.presentX"}` 与镜像常量 `plan_present_doorbell_165r3_test.go:38`；真名册（`internal/tools/plan_present.go:72`，经 `run.go:1126` 注册）零字节。
- 同命令重跑：`--- FAIL: Test330DockPinsTheWireNameToTheAssembledRoster`，rc=1。
- 红在 `plan_present_roster_dock_330r1_test.go:128`，红句具名：
  `G1 对接钉：plan-present.sse 线上名 "plan.presentX" 不在装配出来的真名册（22 枚）里：[clipboard.read clipboard.write fs.edit fs.list fs.move fs.read fs.trash fs.write input.type media.control plan.present screen.capture system.get system.power system.set task.cancel task.list task.output task.spawn window.focus window.list window.manage]`
- 165-v2 报的那发绿（`44f51fc8` 导出树、G1）在本钉下变红＝本票唯一判据成立。

## 还原核验

- 还原用 probes 里的基线副本 `cp` 回（⛔ checkout）：上两枚 sha256 逐位回基线（值同上）。
- `git status --porcelain internal/agent` ＝ 空（零字节漂移）。

## 门禁读数表（终态，全在共享工作树）

| 门 | 读数 |
|---|---|
| `gofumpt -l cmd/wisp/`（绝对路径 `/d/work/base/gopath/bin/gofumpt.exe`） | 列空，rc=0 |
| `go vet ./cmd/wisp` | rc=0，0 字节 |
| `go build ./cmd/wisp` | rc=0 |
| `go test ./cmd/wisp -run 'Test330' -count=1 -timeout 30m` | PASS，rc=0（正向读数同上） |
| `sh scripts/d22scan.sh` | clean，rc=0（全文 `d22scan-final.log`） |

## 中途读数（如实）

- 第一版 d22scan rc=1，唯一 finding＝本腿新件 `plan_present_roster_dock_330r1_test.go:123` 字符串里的 `⛔`（U+26D4，ban #8，字符串不豁免）。已把字符串改口为「不要放宽判据」（注释里的 glyph 豁免未动），复跑 rc=0。仪器牙齿是真的：它抓到了本腿自己的字节。

## AC#2 尺（新增名字字面量枚数）

- 尺：`grep -c '"plan\.present' cmd/wisp/plan_present_roster_dock_330r1_test.go` ＝ 0（带引号的工具名字符串字面量零枚；头注里的提及全是注释、注释豁免；夹具文件名 `plan-present.sse` 是连字符路径、不是工具名）。
- 这枚字符串是从哪一行读出来的：
  - 线上名＝`internal/agent/testdata/golden/plan-present.sse:12` 的 `"function":{"name":...}` 字节，由 `wireNamesFromSSE330r1` 结构化抽出；
  - 真名册＝`cmd/wisp/run.go:1126` `tools.BuiltinPlanEntries()` 注册进 `rt.bridge.Registry()` 的那一册，名字最终出自 `internal/tools/plan_present.go:72` 的 `Name()`。

## 基线副本（只建不删）

- `baseline-golden-plan-present.sse`、`baseline-plan_present_doorbell_165r3_test.go`：突变前逐字节副本，还原即 `cp` 回。
