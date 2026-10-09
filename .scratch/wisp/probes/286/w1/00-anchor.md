# 286-w1 / 00 起手锚（commit-first 闸门：在任何长跑命令之前落第 1 笔）

- 腿代号：`286-w1`（票 286 AC#4 搬期望＋反形自证 ／ AC#5 门禁）
- 落锚时刻：`2026-10-09 11:48 +0800`（`date "+%Y-%m-%d %H:%M %z"` 原文）
- 分支：`dev`
- HEAD：`332e2828`
- HEAD subject 逐字（首 80 字）：`282-r1 读数落件：票 282 AC#1 两发各有指名红（subagent_197.go:262 种回 corr → ticket283:194 红；`

## 起手在飞字节登记（一律不碰）

`git status --porcelain -- internal cmd frontend` 原文：

```
（空输出 ⇒ internal / cmd / frontend 三树起手即净，无别人的字节在其中）
```

`git status --porcelain | wc -l`（全仓在飞条目数）= **759**
`git status --porcelain | head -40` 摘录（具名登记、⛔ 一律不碰）：

```
 M .gitignore
 M .scratch/wisp/issues/286-golden-loop-test-still-nails-corr-equals-taskid-after-the-242-landing.md
 M .scratch/wisp/probes/152/my152.py
 M .scratch/wisp/probes/161/r6/logs/flip-1.txt … flip-6.txt / flip-baseline.txt / flip-restored.txt
 M .scratch/wisp/probes/242/r3/logs/probe-routed.txt
 M .scratch/wisp/probes/268/v1/evidence.md
 D design/assets/base.css / icons.js / theme.js / tokens.css
 M design/doubao/README.md / demo/app.js / demo/index.html / demo/styles.css
 D design/index.html / screens/{app,ball,chat,config,cost,firstrun,palette,privacy,security,states,tasks}.html
?? -
?? .scratch/.scratch/
?? .scratch/ci-logs/*
?? .scratch/commit-msg-167-ac1.txt / commit-msg-167a2-s01.txt
```

★ 具名两笔与本票直接相关，登记以免误判：
1. **票 286 文件本身起手即在飞**（` M`）＝编排者 §AC#1 翻勾与 `:69` 更正的未提交字节。
   ⇒ 本票文末追加段与那些字节同文件、不可分离；我最后一笔 commit 带显式 pathspec 提交该文件时
   会把编排者那些既有字节一并带走（同文件无法局部提交）。已在回报里具名。
   ⛔ 我不改票面任何既有句子，只追加文末。
2. **`.gitignore` 在飞** ＝ 不是我动的；根 `.gitignore:8` 全仓忽略 `*.out`（证据件 ⛔ 不叫 `.out`）。

## 我要动的那一枚（现量，尺＝同一工作树 Read ＋ `git show HEAD:<path>` 同一次取数）

`internal/agent/loop_golden_test.go:69` 逐字：

```
	if call.Req.TaskID != res.TaskID || call.Req.CorrelationID != res.TaskID {
```

`:70` 逐字（红句里报的 `:70` 是这一行的行号，不是 `if` 的）：

```
		t.Errorf("call identity = %+v, want task %s", call.Req, res.TaskID)
```

产码现形（⛔ 一字不动）`internal/agent/loop.go:603` 一带：

```
func callCorr(taskID, callID string, index int) string {
	if callID == "" {
		return fmt.Sprintf("%s#call-%d", taskID, index)
	}
	return taskID + "#" + callID
}
```

唯一调用点 `internal/agent/loop.go:676` 逐字：`			TaskID: taskID, CorrelationID: callCorr(taskID, p.call.ID, i),`

golden 回放里该 call 自己的 id ＝ `call_e1`（同文件 `:66` 已钉：
`	if call.Req.Name != "echo" || call.Req.CallID != "call_e1" {`）。

## ⛔ 禁区复述（本腿自约束）

- 同文件 **`:341`** 逐字 `if row.Tool != "echo" || row.CorrelationID != res.TaskID {` ＝ journal 侧，
  今天靠残余而绿；票 282 `AC#4`(a) 已裁"留"（台账 `A758`）⇒ ⛔ 不碰。
- ⛔ `internal/agent/loop.go:369`、`internal/agent/journal.go:59`／`:81` 一字不动。
- ⛔ `internal/agent/testdata/golden/**` 任何字节；SLO／`internal/observe/thresholds.go`／`docs/SLO.md`／`allowlist.txt`。
- ⛔ 三枚冻结件 `internal/panel/tokens_fourway_test.go`／`internal/panel/l2_grant_boundary_test.go`／`internal/perm/ticket90_persist_test.go`。
- ⛔ `t.Skip`、⛔ 放宽既有断言换绿、⛔ 回退产码、⛔ 改契约文本（D1–D47／C1–C32／R1–R9／D43）。
- ⛔ `frontend/**`／`design/**` 零读零写零转述（上面只登记 status 行，未读内容）。
- 只 commit 不 push；commit 必带显式 pathspec；⛔ `git add -A`／`.`／`--amend`／`reset`／`rebase`／`stash`／`checkout .`／`clean`／merge／worktree。
- 临时件只建不删。

## 读数落盘形制（本腿采用的裁断）

派单正文步骤 2 写了 `logs/before.txt`，同一条派单文末「输出落盘规矩」写了
「只有 `*.md` 会被门禁放过，新建 `.sh`/`.ps1`/`.txt` 会挪 gofumpt／d22scan 的分母，
读数一律写进 `.md` 件里」。两行冲突 ⇒ 按**更严且更靠后的那一条**执行：
基线/改后读数落 `logs/before.md`、`logs/after.md`（命令与 rc 同件并排贴）。⛔ 无 `.out`。

## 本腿承诺的工序

1. 改前基线 `GOFLAGS= go build ./...` ＋ `go test ./internal/agent/ ./internal/tools/ -count=1`（必带 `-count=1`）
2. 搬 `:69` 期望（前缀＋后缀逐字＋不等于 taskID 三件同钉；`TaskID` 那一半等值期望照原样留）⇒ 立刻 commit
3. 改后复跑同一命令（`internal/agent` 应 0 FAIL，PASS 只增不减）
4. 反形自证：盘上种 `callCorr` ＝恒返回 taskID ⇒ `TestGoldenSingleToolCall` 必须红 ⇒ 还原并 hash 对拉＋porcelain 空
5. 门禁：`gofmt -l`／`gofumpt -l` 对动过的件空、`sh scripts/d22scan.sh` rc=0、`git show --stat` 名册、`sh scripts/check-path-length-budget.sh`
6. 票 286 文末追加（只追加；⛔ 不翻任何框）
