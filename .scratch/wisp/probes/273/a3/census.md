# 273-a3 普查件 — AC#1 未修码读数（带正控）＋ AC#2「乙」那一行的逐名代价（只读腿，零产码改动，⛔ 不选形）

leg id `273-a3`｜票面 `.scratch/wisp/issues/273-shipping-process-builds-the-state-machine-without-a-sink-so-every-d43-side-effect-falls-into-a-no-op.md`（**整份读完 89 行**，含文末「收 `273-a1`」`:43-62`／「收 `273-a2`」`:64-70`／「收 `273-v1`」`:72-88` 三节）｜前序件全部整份读完：`.scratch/wisp/probes/273/a1/census.md`（262 行）／`a2/machines.md`（361 行）／`v1/verdict.md`（187 行）／`correction-a653.md`（23 行，内容与票面 `:43-62` 一节同源）。

> ⛔ 本腿**一枚 AC 框没碰**（票面 `[ ]`/`[x]` 状态与起手一致，见 §4）、⛔ **不选形**（甲/乙/丙的取舍归编排者，本腿只把「乙」那一行的代价量到逐名）、⛔ **零产码/测试文件改动**、⛔ **不 push**。
> ⛔ 本腿**没跑过任何 `go test`**（整包与单枚都没跑；写腿 `274-r1` 在飞构建与门禁，抢整包＝洗它的读数）。用过的命令只有 `grep`／`awk`／`sed`／`wc`／`ls`／`git`／`go list`（票面 §2③ 明令依赖边「用 `go list -deps`／`go list -f '{{.Imports}}'` 现量」）／`sh scripts/d22scan.sh`（§4 授权）。

---

## §0 起手锚（现跑逐字，第 1 笔 commit 之前落的数据）

```
$ date
Wed Oct  7 11:48:59 CST 2026

$ git rev-parse --short HEAD
89114888

$ git status --porcelain -- cmd internal scripts tools .github docs frontend
（零行＝这七棵树工作区全干净，RC=0）
```

### 0.1 两枚非 test 构造点现读逐字（本腿自己在锚 `89114888` 上 `sed` 的原文，⛔ 不是引票面）

```
$ sed -n '298,310p' cmd/wisp/models.go
// This function is the grep target of ticket 121 AC#2: it is the non-test
// caller of models.WireDownloading and therefore of DownloadingBridge.Run,
// whose body holds the only non-test call of Manager.VerifyInstalled
// (internal/models/bridge.go:58). logf receives progress ticks and may be nil.
func (ms *modelStore) handOffModel(io_ modelsIO, id string, logf func(string, ...any)) int {
	machine := statemachine.New(statemachine.Options{Initial: statemachine.StateFirstRun})
	defer machine.Close()
	bridge := models.WireDownloading(ms.mgr, machine)

	ctx, cancel := context.WithTimeout(context.Background(), modelHandoffTimeout)
	defer cancel()

	started := time.Now()
```

⇒ 构造点在 `cmd/wisp/models.go:303`（`sed` 起算 `:298`＋第 6 行），**单行字面量、只有 `Initial`，没有 `Sink`**。

```
$ sed -n '183,193p' cmd/balldebug/main.go
	)

	// Tray Exit and Ctrl+C both land here; the teardown is common.
	exit := make(chan string, 1)

	m = statemachine.New(statemachine.Options{Initial: statemachine.StateSleeping})
	b, err := ball.New(ball.Options{
		SizePx:      *sizePx,
		Initial:     statemachine.StateSleeping,
		WindowTitle: os.Getenv(envBallTitle), // private title when spawned by -diff
		Events: ball.Events{
```

⇒ 构造点在 `cmd/balldebug/main.go:188`，**同样只有 `Initial`**；⚠ 同屏一枚同名陷阱复核为真：`:191` 的 `Initial: statemachine.StateSleeping` 属于 `ball.Options{…}`，**不是**状态机的 `Options`（`balldebug` 里两枚 `Initial` 挨着长在一起）。

### 0.2 唯一真传 Sink 的形状（本腿的正控靶面，现读逐字）

```
$ sed -n '14,20p' internal/models/bridge_test.go
	mgr, _ := newTestManager(t, m, nil, nil)
	var effects []statemachine.Effect
	machine := statemachine.New(statemachine.Options{
		Initial: statemachine.StateFirstRun,
		Sink:    func(e statemachine.Effect) { effects = append(effects, e) },
	})
	return WireDownloading(mgr, machine), &effects, machine
```

⇒ `internal/models/bridge_test.go:16-18` 是**跨 3 行的字面量**（`:16` 起 `New(statemachine.Options{`、`:18` 才是 `Sink:`）。
⚠ 这条形状决定了尺的构造：**只看 `statemachine.New(` 那一行等于零信息**，必须向后读到括号闭合——票面 §1 与 `273-a1` §1.5 要的「正控必须命中」在这一形上才会成立。

### 0.3 `internal/statemachine` 包的文件名册

```
$ git ls-files internal/statemachine | wc -l
8

$ git ls-files internal/statemachine
internal/statemachine/doc.go
internal/statemachine/events.go
internal/statemachine/machine.go
internal/statemachine/machine_test.go
internal/statemachine/states.go
internal/statemachine/table.go
internal/statemachine/table_test.go
internal/statemachine/timeouts.go
```

⇒ 非测试面 **6 枚**（`doc.go`／`events.go`／`machine.go`／`states.go`／`table.go`／`timeouts.go`），测试面 **2 枚**（`machine_test.go`／`table_test.go`）。这 8 枚就是「乙」那一行「动哪几枚文件」的定义域（逐名读数见 §2①）。

### 0.4 工作区里别人的未提交件（⛔ 不是本腿动的，本腿不提交、不还原，只如实登记）

起手时 `cmd internal scripts tools .github docs frontend` 七棵树＝**零行**（见 §0 原文）。本腿不承诺整仓干净：`.scratch/**`（本腿自己那枚目录之外还有别的腿在飞）与 `design/**`、`.gitignore` 属别人，本腿一律不碰、不评（票面 §1.4 Git 纪律与本腿指令 §3 同一条）。

---
