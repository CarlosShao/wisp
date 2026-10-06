# 票 232 — 写码腿 `232-r2` 取证件

> 判据来源＝`.scratch/wisp/issues/232-the-restart-tier-test-pins-that-the-sentence-exists-not-that-it-carries-the-why.md`
> （本件**不勾它的框、不改它的 Status、不改它的名**）。本件只记读数与判语。
> 落点＝`.scratch/wisp/probes/232/r2/`。骨架先落盘（编排者对 150 轮帽的固定对策），读数逐格追加。

---

## §0 起手锚（三把尺原文）

| 尺 | 命令 | 读数 |
|---|---|---|
| 撤票尺 | `sed -n '3,4p' <票面>` | 第 3 行开头 `- Status: **待派，测试面加固**…`；第 4 行 `- ⚠ **这不是安全漏洞，照三行读**…`。**零 `WITHDRAWN`／零「撤」／零「作废」** ⇒ 未撤票，许开工。 |
| HEAD 锚 | `git log -1 --format='%h %H %ad' --date=iso` | `590fb855 590fb85576d6ef06b79bf445bf26b72653415bb0 2026-10-06 11:34:27 +0800` |
| 写面尺 | `git status --porcelain -- cmd internal` | 空（0 行）⇒ 起手 `cmd`／`internal` 干净。 |
| 框数尺 | `grep -c '^- \[ \]' <票面>` | `6`（六枚全部不归本腿） |
| 产码 md5 | `md5sum cmd/wisp/config_reload.go` vs `git show HEAD:… \| md5sum` | 两者皆 `5ce441ca5e72b64d18a6c26f1c066882` ⇒ 起手一致。 |

### §0.1 那三段文案的 HEAD 原文（`git show HEAD:cmd/wisp/config_reload.go | sed -n '311,330p'` 复量，抄回存档）

```go
func (rt *agentRuntime) reportRestartPending(sections []string) {
	if rt == nil {
		return
	}
	rt.auditf("config: HOT-RELOAD state=restart-pending sections=%v effect=next-process-start tier=restart",
		sections)
	rt.auditf("config: RESTART-PENDING detail=%q keys=%q",
		"这些键只在进程启动时被读一次并交给平台层（开机自启注册、单实例锁、界面语言），本次运行没有重新加载它们的路径；"+
			"内存里继续用旧值，文件里的新值没有被丢掉，重启进程后生效",
		restartTierKeys)
	fmt.Fprintf(rt.stdout,
		"wisp run: 这些段的改动本次运行不会生效（D36 重启档，需要重启进程）：%v。"+
			"原因：这几枚键在进程启动时被读一次就交给平台层（开机自启注册、单实例锁、界面语言），"+
			"本次运行没有重新注册它们的路径。涉及：%s。"+
			"注意两件事都没发生：文件里的新值没有被丢掉，内存里的旧值也没有被偷偷替换。\n",
		sections, strings.Join(restartTierKeys, " / "))
}
```

（落点行号现量：`cmd/wisp/config_reload.go:318-319` 审计 detail 句、`:321-326` 操作员 `Fprintf` 三段。）

---

## §1 现量（票面那张表的五行，本腿自己复跑）

| 票面事实 | 本腿复跑尺 | 本腿读数 |
|---|---|---|
| 内容断言在哪 | `sed -n '482,495p' cmd/wisp/config_reload_223_test.go` | 未修态：`:490 out := r.awaitStdout(t, "本次运行不会生效")`；`:491-495` 三枚 needle 在 **`why+out` 拼接串**上 `strings.Contains`。 |
| ⛔ 瘦句不响 | AC#1 突变 B（见 §2/§3） | 见 §2 AC#1 行 |
| 方向说反**有**牙 | 同上区段 | `:479-481`（等不到那句就 `Fatalf`）＋`:504-506` 的负向断言（重启档不许被说成立即生效）⇒ 该形会红，本腿未动。 |
| ⚠ 前向风险（横幅假绿） | `grep -rn "本次运行不会生效\|重启进程后生效\|开机自启" --include=*.go cmd internal` | 未修态：产码里三枚 needle 的来源是 `config_reload.go:318-319`（审计句）＋`:322-325`（操作员句）＋`:339`（`restartTierKeys` 里 `app.autostart`）⇒ 操作员句被删薄时，`why+out` 仍可能由 stderr 那侧供词。 |
| 助手本身干净 | `sed -n '109,158p' cmd/wisp/config_reload_223_test.go` | `awaitStdout`(`:140-158`) 与 `awaitAudit`(`:109-127`) 同构：同 `reloadCaseBudget`=40s、同 `time.NewTimer`+20ms `time.NewTicker`、同 `Fatalf` 出口、零 `Skip`。 |
| 未修态两枚点名用例 | `PATH=… go test -count=1 -run 'TestTicket223RestartTierSaysItWillNotApply\|TestTicket223R2FailureSentenceRouting' -v ./cmd/wisp` | `--- PASS: TestTicket223RestartTierSaysItWillNotApply (2.48s)` ／ `--- PASS: TestTicket223R2FailureSentenceRouting (0.60s)` ／ `ok github.com/CarlosShao/wisp/cmd/wisp 3.186s`（全文 `.scratch/wisp/probes/232/r2/logs/AC0-baseline-unmutated.txt`） |

---

## §2 六格逐格表（格 ｜ 跑了什么命令 ｜ 原始读数 ｜ 判语）

| 格 | 跑了什么命令 | 原始读数 | 判语 |
|---|---|---|---|
| AC#1 复现瘦句不响 | `-overlay` 突变 B ＋ `-run TestTicket223RestartTierSaysItWillNotApply` | 本行读数随后追加 | 本行判语随后追加 |
| AC#2 断言改到 stdout＋窗口化 | 改 `cmd/wisp/config_reload_223_test.go` | 本行读数随后追加 | 本行判语随后追加 |
| AC#3 反向正控要真响 | 突变 B 复发 ＋ 新突变 C（换掉「原因」段） | 本行读数随后追加 | 本行判语随后追加 |
| AC#4 横幅假绿那一枚 | 突变 D（横幅先含同样字样） | 本行读数随后追加 | 本行判语随后追加 |
| AC#5 零放宽 | 两枚点名用例逐名 | 本行读数随后追加 | 本行判语随后追加 |
| AC#6 整包终态 | `go test ./cmd/wisp ./internal/... -count=1` | 本行读数随后追加 | 本行判语随后追加 |

### §2 末尾安全尺两把

| 尺 | 命令 | 读数 |
|---|---|---|
| 写面终态 | `git status --porcelain -- cmd internal` | 本行读数随后追加 |
| 产码未动 | `md5sum cmd/wisp/config_reload.go` | 本行读数随后追加（须等于 HEAD 值 `5ce441ca5e72b64d18a6c26f1c066882`） |

---

## §3 突变名册（每发：overlay 路径 ＋ 落地证明 ＋ 红/绿原文行 ＋ 还原后 md5）

> 种针一律 `go test -overlay <json> …`，合成副本一律在 `D:/tmp/wisp232/` 下；
> 共享工作树里的 `cmd/wisp/config_reload.go` **一个字节都没被编辑过**（每发跑完对拉 md5）。

| 发 | 类型 | overlay json | 替换 src→dst | 落地证明 | 用例红/绿原文行 | 跑完 md5 |
|---|---|---|---|---|---|---|
| B | AC#1／AC#3 | 本表读数随后追加 | | | | |
| C | AC#3 新突变 | 本表读数随后追加 | | | | |
| D | AC#4 横幅 | 本表读数随后追加 | | | | |

---

## §4 门禁读数

- 本腿写面＝只有 `cmd/wisp/config_reload_223_test.go` ＋ 本件目录 ＋ 票面 Progress log 一行 ⇒ **零脚本面**（没动 `scripts/**`／`.github/**`），
  所以 `bash -n` 与 `sh scripts/d22scan.sh` 两把尺**不适用于脚本改动**；但 d22scan 对**代码形状**的禁令适用，读数见本表下列行。

| 尺 | 命令 | 读数 |
|---|---|---|
| 形状门禁 | `go vet ./cmd/wisp`（`go test` 自带前置） | 本行读数随后追加 |
| d22scan | `sh scripts/d22scan.sh` | 本行读数随后追加 |

---

## §5 判不动的地方（诚实栏）

- 本栏随后追加；若某格两把尺跑第二遍仍拿不到期望读数，这里逐把抄原文，不换尺凑绿。
