# 33-r13 起手锚（票 33 `:391` 那半句"直到 ⓐ 落地"的措辞归位＝只动一枚 `.md` 的一行）

- 工单＝票 33 `.scratch/wisp/issues/33-panel-host-c27.md`｜腿号＝`33-r13`
- 本腿写面＝`.scratch/wisp/issues/33-panel-host-c27.md`（1 行，行内换字）＋本目录 `.scratch/wisp/probes/33/r13/**`（⛔ 只建 `.md`）
- ⛔ 零 Go 命令（一条都没跑：`go build`/`vet`/`test`/`list`/`env` 全 0 发）｜⛔ 零 AC 框翻动｜⛔ 零 push｜⛔ 不碰 `cmd/wisp/**`、`internal/**`、`docs/**`、`.github/**`、`frontend/**`、`design/**`、`docs/PLAN.md`、`docs/specs/**`、`thresholds.go`、任何 golden
- `probes/33/r13/` 目录起手**不存在**（`ls` 于起手时报"No such file or directory"），由本腿新建＝与"应为空"相符

## 五条起手读数（同发取，逐字原样）

```
$ date '+%Y-%m-%d %H:%M %z'
2026-10-08 13:26 +0800
（同一批发完锚读数后于 13:2x 复跑一次，两读之间无位移）

$ git log --oneline -1
53d73db4 A715 落账：收我自己那枚票 242 射程注（cade79b3，⛔ 零翻框、原判据一字未改，…）

$ git rev-parse --abbrev-ref HEAD
dev

$ git status --porcelain -- .scratch/wisp/issues/33-panel-host-c27.md
（空输出＝该文件工作树干净，无别人的在飞改动 ⇒ 可以动手）

$ wc -l .scratch/wisp/issues/33-panel-host-c27.md
394

$ grep -cE '^[[:space:]]*- \[ \]' .scratch/wisp/issues/33-panel-host-c27.md
13

$ grep -cE '^[[:space:]]*- \[x\]' .scratch/wisp/issues/33-panel-host-c27.md
1

$ grep -n 'ⓒ 一枚会响的钉' .scratch/wisp/issues/33-panel-host-c27.md
391:- **ⓒ 一枚会响的钉＝与 ⓐ 同批**。形状定死：**绿着交**——凡"`Locked` 后缀 ＋ 函数体自取锁"同时出现即红，现存那 2 枚进**具名豁免名册**直到 ⓐ 落地；⛔ 不许留一枚恒红用例让后续程继承（本项目有此定式）。现成尺形照抄 `internal/config/manager_223_test.go:39-51` 那一族，⛔ 不新造机制。
```

## 内容锚（⛔ 不以行号为凭）

要改的那半句逐字＝`现存那 2 枚进**具名豁免名册**直到 ⓐ 落地`，它在票面末节
`## 10-08 11:5x 编排者裁一笔旧账`（`:383` 起）的 **ⓐ／ⓒ／ⓑ 三支柱**里，位置**起手量在 `:391`**；
本腿与后续程一律按上面那句逐字文认它，`grep -n 'ⓒ 一枚会响的钉'` 是可重跑的尺，`:391` 只是起手时刻的读数。

## 改前票面的盘上真状态（本腿自己 `git show`／`grep` 读，⛔ 不信转述；只读不跑 Go）

| 一问 | 我这把尺 | 读数 |
|---|---|---|
| ⓐ 那两枚改名落了没有？ | `grep -n 'firstRoundTrip\|serveNotBuiltNotice' cmd/wisp/panel_host_windows.go` | **已落**＝声明 `:840 func (m *PanelManager) firstRoundTrip(...)`、`:460 func (m *PanelManager) serveNotBuiltNotice()`，调用点 `:448`／`:451` 同形；两枚旧名（`…Locked`）在该文件**零命中** |
| ⓒ 那枚守卫建了没有？ | `ls cmd/wisp/ \| grep -i locked` ＋ `git ls-tree --name-only HEAD cmd/wisp/` | **已建**＝`cmd/wisp/panel_locked_naming_33r11_windows_test.go`（工作树与 HEAD 各一枚同名，**448 行**） |
| 豁免名册今天几枚条目？ | 读该文件 `:66-80` 名册文本（`var lockedNamingRoster33r11 = map[string]string{…}`） | **1 枚**＝`"PanelManager.setPriorFocusLocked": "requires the caller to hold m.mu; body takes no lock"`；件内 `:75-77` 明写本腿那两枚**deliberately ABSENT** |
| "名册条目消失即红"是真牙？ | 读该文件 `:253-257` | 是：`for qualified := range lockedNamingRoster33r11 { … t.Errorf("lockedNamingRoster33r11 names %s, which this package no longer declares - drop the exemption or restore the method") }` |
| `cmd/wisp` 产码里今天还剩几枚 `Locked` 后缀方法？ | `grep -rn --include='*.go' -E '\bfunc \([^)]*\) [A-Za-z0-9_]*Locked\(' cmd/wisp/*.go \| grep -v _test.go` | **1 行**＝`cmd/wisp/panel_host_windows.go:625 setPriorFocusLocked`（＝名册里那一枚）；⚠ 全仓另有 22 枚属票面末段"同形名册"那一账，**不在本枚守卫射程** |

## 凭据（谁判这句该改／谁落的 ⓐ）

- 台账 `docs/reports/pending-and-issues.md` **`A713` §5 ②**（现量 `:14004`）逐字：
  「票面那句"现存那 **2** 枚进名册、直到 ⓐ 落地"与 ⓐ 同批执行＝**没有可落的瞬间**…⇒ 它采了"现存**合法**的进名册"那一读、并补"条目消失即红"＝**这一枚票面措辞要我改**，与①同一批。」
- ⓐ＋ⓒ 的落地腿＝`33-r11`：产码＋守卫一笔 `9995f9b1`（10-08 12:36），证据两笔 `cc5f7627`／`23ade0da`；编排者收件＝`A713`（`bb871191`，10-08 12:59）。
- 同一文件同一族的**文本归位先例**＝`33-r12`：只改票面 `:358` 一枚行、判据＝「只有"用现在时描述代码"的句子才许改，历史叙述一枚不许动」（`probes/33/r12/report.md` §①／§⑤，收件＝`A714`／`040e424d`）。
- ⓑ 那一支今天**一行未做**（票面 `:392` 原样写着"＝否掉"），本腿的措辞⛔ 不把它写成做过。
