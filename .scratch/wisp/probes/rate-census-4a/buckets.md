# rate-census-4a — 工单池编号 1–49 开放票分桶（按"owner 会不会感觉到"）

## §0 起手锚

- 分支：`dev`
- HEAD sha（取数时刻）：`6fe300c5035a60586e681f586a4ead5e2b56137b`（短号 `6fe300c5`，2026-10-02 08:41:33 +0800，提交主题＝census(rate-census-3a) 150-199 分桶）
- 取数时刻：`date` = 2026-10-02 08:46 +0800
- 只读声明：本腿**零写入工单池**（票面 AC 勾选框一枚未动、票名一枚未改、`Status:` 行一枚未改），零改动冻结件
  （`docs/PLAN.md` · `docs/specs/**` · `docs/BUILD.md` · `docs/SLO.md` · `internal/observe/thresholds.go` ·
  `tools/d22scan/allowlist.txt` · 三枚冻结测试件），未写 `docs/reports/pending-and-issues.md`。
- 零负载声明：**未跑** `go test` / `go build` / `go vet`，未执行任何 exe，未跑任何计时或压力工具。
  读数全部来自 `ls` / `grep` / `wc` 与读文件本身。
- 零界面声明：`frontend/**` 与 `design/**` **零读取、零引用、零写入**。
  本段里地界是界面的那几枚（34、36、37、38、39、40、77 那类）**只按工单文件名与票面正文分类**，
  没有去界面目录核实"做没做"——这是本腿 §3 若干条判断的共同限制，见 §6。

## §1 分母实测

尺（逐字复跑编排者那一把，未换写法）：

```
cd .scratch/wisp/issues
ls | grep -E "^(0?[1-9]|[1-4][0-9])-" | grep -v -- "-done" | wc -l   -> 33
```

- **开放 33 枚**＝编排者 08:4x 现跑的 33，**一致，无差**。
- 号段 1–49 内共 **49** 枚 `.md`（编号无缺号）；带 `-done` 后缀 **16** 枚；33＋16＝49 自洽。
- 判完成只认文件名 `-done` 后缀。票面 `Status:` 行整批过期，**一枚未采信**。
- 子目录：**不存在** `voided/` 之类目录（`find . -maxdepth 1 -type d` 只返回 `.`）。
  分母里没有任何被排除的候选；`README.md` 不落在号段正则内，天然未计入。
- 工单池总文件枚数 237（含 `README.md`），本腿射程只有 1–49。

名册（33 枚，`nl` 实出，顺序＝`ls` 字典序）：

| # | 文件名 |
|---|---|
| 1 | `12-cli-text-path-s1-gate.md` |
| 2 | `15-speech-engines-cer-harness.md` |
| 3 | `16-s2-acceptance.md` |
| 4 | `20-host-bridge-fs-tools.md` |
| 5 | `21-approval-gates-minimal.md` |
| 6 | `22-web-tools-d30.md` |
| 7 | `23-system-window-input-tools.md` |
| 8 | `24-doc-search-tools.md` |
| 9 | `25-s3-acceptance.md` |
| 10 | `26-tts-output.md` |
| 11 | `27-punctuation.md` |
| 12 | `28-session-scope-warm.md` |
| 13 | `29-memory-l1l2.md` |
| 14 | `30-result-routing-d10.md` |
| 15 | `31-reminders.md` |
| 16 | `32-s4-acceptance.md` |
| 17 | `33-panel-host-c27.md` |
| 18 | `34-frontend-scaffold.md` |
| 19 | `35-panel-bridge-c17.md` |
| 20 | `36-result-history-panel.md` |
| 21 | `37-approval-ui-l2.md` |
| 22 | `38-command-palette-tasks.md` |
| 23 | `39-config-editor-gui.md` |
| 24 | `40-security-privacy-cost-pages.md` |
| 25 | `41-kws-wake-word.md` |
| 26 | `42-watchdog.md` |
| 27 | `43-power-lifecycle-events.md` |
| 28 | `44-costmeter-c23.md` |
| 29 | `45-diagnostics-guards.md` |
| 30 | `46-s6-acceptance.md` |
| 31 | `47-task-scheduler-pathlock.md` |
| 32 | `48-approval-queue-full.md` |
| 33 | `49-session-grants-d45-2.md` |

勾选框口径（先自跑再写）：锚定 `^- [ ]`（未勾）与 `^- \[[xX]\]`（已勾）现量。
本段 33 枚**全部用勾选框格式**（不存在 219 那种全篇零勾选框的枚）：
逐枚 `grep -c '^- \[ \]'` 最小值是 1（票 20），**没有一枚量出 0 格**，所以本段**无"格数量不出来"的条目**。
全段合计：未勾 **185** 格 ／ 已勾 **13** 格（13 格集中在票 12／20／33 三枚，即号段里那三枚大票）。
缩进式勾选框（`  - [ ]`）逐枚 grep 命中 **0** ⇒ 锚定行首没有漏计。

## §2 四桶加总

（本轮尚未分类完，加总在分类跑完后落此节；口径见题面：甲／乙／丙／待定，一枚只进一桶，加总须等于 33。）

## §3 名册

（逐枚一行：`号｜桶｜一句话｜未勾格数｜归界面侧是／否｜疑似被号段外某枚覆盖`。）

## §4 主链表

（本段还挡路的枚数＋链序；跨号段前置只具名点出、本腿未读别人号段的票面。）

## §5 本号段到能用 ≈ N 枚

（结论一句，写明口径与不成立条件。）

## §6 我可能分错的条目

（逐条：如果错了该进哪桶＋错了的后果＝哪个数字会变。）

## §7 占位符自查

读数命令与结果在本节末次写入时现跑，见提交回执。
