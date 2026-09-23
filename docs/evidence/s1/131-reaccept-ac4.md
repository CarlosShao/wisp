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
我把两枚锚点各跑了一遍同一仪器，然后逐枚对名字集合：

```
git archive 56d8026 | tar -x -C /tmp/wisp131-r3-old        (+ 三枚 dll)
go test -count=1 -v ./cmd/wisp/   -> RUN 82 / TOP PASS 46 / FAIL 0 / SKIP 0, rc=0   # 与上一任 §2 的 CI 形状逐字相同
go test -count=1 -v ./cmd/wisp/   -> RUN 100 / TOP PASS 53 / FAIL 0 / SKIP 0, rc=0  # bcb03aa（§1 基线）
comm -23 <旧名字集> <新名字集>    -> 空    # 一枚也没丢
comm -13 <旧名字集> <新名字集>    -> 恰 18 枚
```

那 18 枚逐字是：**票 128** 的 `4e5d240`（`TestAC2EveryLegRefusesTheSameShapeAndWritesNothing128` +5 子、
`TestAC2RealProcessRefusesOnEveryLegWithoutAppData128` +4 子、`TestAC2RefusalMarkersAreNotAShortenableList128`、
`TestAC2ResolveDataDirRefusesInsteadOfFallingBackToCWD128` +2 子、`TestAC2TestDataDirBranchStillResolves128`
⇒ **16 枚 = 5 顶层 + 11 子用例**）与**票 130** 的 `d87905c`
（`TestAC3EarlyRecordLandsOnDiskBeforeTheInstallRecord`、`TestAC3EarlyReplayKeepsTheSinkInsideTheDataRoot` ⇒ **2 枚顶层**）。
逐枚静态反查同数：`func Test` 枚数 `56d8026`=49 → `4e5d240`=54（+5，128）→ `d87905c`=56（+2，130）→ `bcb03aa`=56。
CI 的 `-skip` 名单与这 18 枚零交集（逐名 grep）。

⇒ **数字结论成立、归因写错了名字。** `+18 RUN / +7 顶层 PASS` 这一笔我逐枚对上，**零枚被洗掉、零枚丢失**；
但它是**票 128（16）＋票 130（2）**，不是续单证据与票面写的"全归 129/130"。
`git log 56d8026..b723978 -- cmd/wisp` 里**根本没有 129**（129 的改动面在 `internal/**`，一枚 `cmd/wisp` 用例也没加），
而 128 被漏名。⇒ 登记为 `R-131r3-4`（措辞/归因级，不改判"四数不降"）。
另两枚在该区间动了 `cmd/wisp` 的 commit（`11f3927`＝114 只改 `run.go`、`9b5d64d`＝130 只改既有钉文件、
`717d822`＝本票第 0 件改名）对分母零影响，我也逐枚 `--name-only` 看过。

### 1.2 被驱腿自己的行号也下移了（续单那句话不完整）

续单/票面写"六发红名逐字不变，**只有门与钉自己的 `file:line` 前缀下移**"。前半成立，后半不完整：

```
bcb03aa vs 56d8026 的 install 点：logsink.go:139->144、models.go:276->284、secret.go:219->226、run.go:164->178、resident_windows.go:57 未动
X2 红名里的 emit-site：slog.Info@secret.go:321 -> :328, slog.Info@:461 -> :468, slog.Warn@:461 -> :468
```

⇒ 红名的**文本结构与所点的那条腿一字未变**（我逐字比过，见 §2），但**位移的不止仪器自己的前缀**——
被驱腿文件的行号也被票 128（`secret.go` +11/-4）与票 130（`models.go` +9/-1、`logsink.go` +22/-2）推着走了 7–14 行。
这句话要更正，否则下一个人拿"只有前缀动"去核 X2 会以为读数错了一位。并入 `R-131r3-4`。

---

## 2. 结案判据 1：六发红名逐字不变（X1/X2/X3/X6/X9/X10）—— **PASS**

仪器与上一任同形：`git archive bcb03aa` 纯净快照 + 三枚 dll + `go build ./cmd/wisp/` rc=0 之后才读 `go test -count=1 -v ./cmd/wisp/`。
变异我**按上一任的字节形状重造**（X3/X6/X8/X12 的 `main.go` 与被种的 `fake131.go`/`diag131.go` 直接取上一任快照里那枚文件，
`diff` 证同形；X1/X2/X3 的删行区间逐字同 `284,292d283` / `226,230d225` / `74,81d73`）。

| 发 | 造法（同上一任） | 上一任四数 | **我这次四数**（RUN/顶层PASS/FAIL/SKIP） | 原来那枚红 | 名字换没换 |
| --- | --- | --- | --- | --- | --- |
| **X1** 删 `models.go` 那 9 行 | 删 284-292，`grep -c installLogSink` = 0，334→325 行 | 82/78/**4**/0（3 顶层+1 子） | **100/50/3/0 + 1 子用例红**，rc=1 | `AC#4 RED: nail "TestAC2ModelsLegBooksItsHandOffVerdictOnDisk" claims leg "models", which no longer installs the persistent sink (main.go:74): the nail is now proving something else, or nothing` | **未换**，仅 `214→318`、钉 `307→358`、`455→506` |
| **X2** 删 `secret.go` 那 5 行 | 删 226-230 | 82/78/**4**/0 | **100/48/5/0 + 1 子用例红**，rc=1 | 门给**两条独立读数**（`:298` 要裁决、`:318` 要钉）文本逐字同 | **未换**；多出的 2 枚红点名归票 130 的两枚 early-record 用例（`d87905c`），方向是查得更多 |
| **X3** 删掉 `case "models"` 整支、`models.go` 一字不动 | 删 74-81，生产调用者只剩声明 | 82/81/**1**/0 | **100/52/1/0**，rc=1，唯一一枚红在门上 | `... which main.go does not dispatch to any more` 逐字同 | **未换** |
| **X6** 新腿直呼 install、不给钉 | 取上一任那枚 `fake131.go` + `81a82,83` | 82/81/**1**/0 | **100/52/1/0**，rc=1 | `leg "fake131" (main.go:82) reaches installLogSink on 1 line(s) of this package and no registered nail names it.` 逐字同，行号也是 `main.go:82` | **未换** |
| **X9** 两处 install 挪去 `<根>-elsewhere` | 各改 1 行 | 门 PASS + 两枚钉红 | **100/48/5/0 + 2 子用例红**；门 `--- PASS`、`AC#4 RED` 命中 **0** | 两枚钉各红一次，红名正文带 `dir=...\001-elsewhere\logs` 原文 | **未换**；+2 枚同样是 130 那两枚 |
| **X10** 注册一枚不读盘的用例当钉 | `registerLegNail131("models", TestSecretFlagsAreBoolOnly)` | 82/81/**1**/0 | **100/52/1/0**，rc=1 | `nail "TestSecretFlagsAreBoolOnly" for leg "models" calls none of the shared sink readers (readLegSink131, readResidentSink, readSink) ... (leg site main.go:74)` 逐字同 | **未换**；门内**另多一枚**红＝本轮为 `R-131-2` 加的"没带入口符号"那条（同一枚用例内，顶层枚数不变） |

⇒ **判 PASS。** 六发原来那枚红**逐枚还在、名字一个没换**；新增的红全部点名归到票 130 的用例与本轮 `R-131-2` 的新判据上，
方向是**查得更多**。⚠ 两处措辞要更正（不影响判定）：X2/X9 各多的是**两枚**（`100/48/5` 对上一任 `82/78/4`），
续单证据与 A122① 一处写"两枚"一处写"一枚"，实测为**两枚**；X10 多的那一枚在**同一枚用例内部**，顶层 FAIL 枚数没变。
