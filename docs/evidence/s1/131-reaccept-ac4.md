# 票 131 AC#4 复验（第三任）—— 裁的就是上一任判退回那一格

**验收方**：`acceptor-ticket131-r3`（本文件唯一作者；只读验收 + 自造变异，本文件是唯一写件）
**日期**：2026-09-23（开件 15:5x，逐节落盘，每节一枚 commit）
**被验交件**：票 131 的「续单」三件（续#1/续#2/续#3），实现方 `worker-ticket131-followup` + `worker-ticket131-followup-r2`
**被验 sha**：**`bcb03aa`**（= `717d822` 第 0 件 + `6f702ae` 三形判据 + `c49327c` 追认取证 + `bcb03aa` 票面翻格）
**判据原文**：`docs/evidence/s1/131-adversarial-acceptance.md` §4（X1..X12）、§6（`R-131-1..4`）、§8.3（结案三件事），
以及票面「续单 → 结案判据」那三句。

---

## 0. 锚定与被验面的完整性（先证明我读的是被验版本）

- 开工 `git rev-parse HEAD` = **`ab1462a8380a5a40d78a3d8b2e6c5b6598006a16`**，工作树 `git status --short` **为空** ⇒ 三任代理都入库了。
- HEAD 已前进一格（`bcb03aa..ab1462a` = `ab1462a` 本身）。按简报要求逐枚核被验面：

```
$ git diff --numstat bcb03aa..HEAD -- cmd/wisp
(零行)
$ git diff --numstat bcb03aa..HEAD
33  0   .scratch/wisp/issues/124-posix-harness-hands-the-sealing-floor-an-unresolved-tempdir-79-to-83-reds.md
24  0   docs/reports/pending-and-issues.md
```

⇒ **`cmd/wisp/**` 在 `bcb03aa..HEAD` 之间零字节改动**，唯一增量是两枚文档 ⇒ 被验面未被移动，本文全部读数同时适用于 `ab1462a` 的 `cmd/wisp`。
- **快照**：`git archive bcb03aa | tar -x -C /tmp/wisp131-acc-r3`，随后
  `cp third_party/sherpa-onnx/*.dll /tmp/wisp131-acc-r3/third_party/sherpa-onnx/`（三枚：`onnxruntime.dll`、
  `sherpa-onnx-c-api.dll`、`sherpa-onnx-cxx-api.dll`）。**没拷就是 `0xc0000135`、四数全零**，简报预告过，我按它做了。
- **两记反查**：① `diff -r --exclude=third_party /tmp/wisp131-acc-r3 /tmp/wisp131-verify-acc-r3`（另一枚现取的
  `git archive bcb03aa`）⇒ **rc=0，零差异**；② `diff -r "D:/work/workspace/projects plans/Wisp/cmd/wisp" /tmp/wisp131-acc-r3/cmd/wisp`
  ⇒ **rc=0**，快照的 `cmd/wisp` 与仓库工作树逐字节相同。
- **未在仓内建 worktree / 未 checkout**；全程 `PATH=<快照>/third_party/sherpa-onnx`。
- ⚠ **一轮在飞观察**（不是裁决）：16:0x 我 `git status --short` 第二次量到
  ` M docs/reports/injection-timeline.md`（+21/-0，`git status --short -- cmd/wisp` 仍为空）。
  ⇒ 有兄弟代理正在写那枚文档，**被验面未被碰**，我一格未动它。
- **测前查在飞 run**（本机即 self-hosted runner）：`gh run list --repo CarlosShao/wisp --limit 5` ⇒ 五枚全 `completed`
  （最新 `35832874239` 07:39Z），无 in_progress ⇒ 本轮是安静窗口。整轮读数里 `0xc000013a` 早死 **零命中**（各发日志的
  `exit status` 计数见 §5 附）。墙钟只作旁证。

---

## 1. 基线复算（`bcb03aa` 纯净快照，未变异）—— 新基线是不是真的

简报警告"上一任的 82/82 已被三次移动过，谁拿它当分母谁就错"，也警告"重新基线必须是逐数解释过的"。两边我都量了。

| 仪器 | 续单方自报（`b723978`） | **我复算（`bcb03aa`）** | 判定 |
| --- | --- | --- | --- |
| `go test -count=2 -v ./cmd/wisp/` | `=== RUN` 200 / 顶层 106 / 0 / 0，rc=0 | **RUN 200 / 顶层 PASS 106 / FAIL 0 / SKIP 0**，rc=0，`ok` 112.509s | 同数 |
| `sh scripts/wisp-cli-tests.sh`（CI 形状） | 100 / 53 / 0 / 0，rc=0 | **RUN 100 / 53 / 0 / 0**，rc=0，`ok` 54.953s | 同数 |
| `sh scripts/d22scan.sh` 八 scope | 203/22/40/18/16/40/390/37，rc=0 | **203 / 22 / 40 / 18 / 16 / 40 / 390 / 37**，rc=0 | 同数 |
| `gofmt -l cmd/wisp/` | 空 | **空**，rc=0 | 同 |
| `$(go env GOPATH)/bin/gofumpt.exe -l cmd/wisp/`（v0.7.0，实存） | 空 | **空**，rc=0 | 同 |
| `go vet ./cmd/wisp/`（windows） | rc=0 | **rc=0** | 同 |
| `GOOS=linux go vet ./cmd/wisp/`（交叉） | rc=1，唯一错误行是 sherpa | **rc=1**，逐错误行归因见下 | 同 |
| Linux 原生 `go vet ./cmd/wisp/`（Docker） | rc=0 | **rc=0**，见 §3 | 同 |

基线整份 `-count=2 -v` 里 **`AC#4 RED` 命中数 = 0**（`grep -c 'AC#4 RED' /tmp/wisp131-acc-r3-base-count2.log`）。
⇒ 第二把尺在未变异树上零误伤，这一条我独立量到了（续单方在 `b723978` 报同数）。

### 1.1 "+18 RUN / +7 顶层 PASS 全归 129/130" 这句我逐枚对过

上一任锚 `56d8026` 的 CI 形状是 `82 / 46`，新基线是 `100 / 53` ⇒ 差 **+18 RUN / +7 顶层 PASS**。
逐枚查谁在这段区间往 `cmd/wisp` 加了用例：
