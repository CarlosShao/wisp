# 174-v2 证据：票 174 AC#6 形 乙（`174-r6`）的非实现者判语 —— 成立（强度差如实保留）

- 锚（起手）：HEAD `4c7f0726` / 2026-10-14 14:30 +08；交卷 HEAD `76b9e214`（他腿 253-r6 收件＋编排者 `runtests.sh` 修复；`internal/tools/` 与 `portable-tests.sh` 零重叠）。被测 blob `internal/tools/task_output_pointer_notice_test.go` sha256 `15e0ef42…`（= `a65bb386` 的 blob，工作树与 HEAD 零差）。
- 身份：验收腿（174-v2）。⛔ 未翻框、⛔ 未动 `Status:`/票名、`internal/**`／`scripts/**`／`cmd/**`／`docs/PLAN.md` 零字节；突变一律 `-overlay`（跟踪件零字节，还原 sha256 见 logs/04）。
- 三把尺复核（编排者给的三发）：两笔 commit `a65bb386`/`11cfbc93` 在（13:34:06／13:45:52）｜件齐（`probes/174/r6` 八节＋logs `00–04` 入库＋两枚原始捕获留盘）｜亲跑 `-run 'TestCanonicalizeErrorNoticeMustNotRelayTheFiledPath'` = **ok 1.959s rc=0**（我这一发）。

## 1 结论（先行）

- **形 乙＝成立**（对 `174-r6` 的主张予接受）。同一判据（notice 不许转发 filed path）在 POSIX 上有真可跑形：空白字面量走 `pointerNotice` 的 `case err != nil`（`task.go:839`），同一条 `TrimSpace` 闸（`paths.go:131-134`）平台无关。
- **强度差如实保留**：POSIX 上"路径半边"弱一档（被拒字面量只能是空白串，不是带毒值路径字面量），r6 已用 4 枚 Unicode 空格锐序列＋`%q` 拼法补强；**我另发现新枚还多出一档它自己没声称的强**（M4：`strings.ToLower(raw)` 转发在 r3 臂上漏、新枚抓住——见 §3）。
- **共存成立**：r3 用例／ledger 行（sha `b1ffdb59…`＝r5 落地时逐字节）／`probes/**` 一字节未动；新名不入名册。
- 整票**不可收口**：AC#6 已 `[x]`（13:18 甲形 CI 实证后翻），其余 AC#3／#2b／#2c／#4／#5 仍空 ⇒ 收口建议见 `docs/evidence/s1/ticket174-174-v2-accept-v1.md` §6。

## 2 四证读数（每格都能逐字重跑）

| # | 证 | 命令 | 读数 |
|---|---|---|---|
| 1 | windows 双跑 | `go test -count=1 -v -run 'TestCanonicalizeErrorNoticeMustNotRelayTheFiledPath' ./internal/tools/` | rc=0；r3 真建 junction **PASS**（notice 逐字＝r6 记录）、新枚 **PASS**；`ok 1.959s`。logs/01 |
| 2 | linux 静态 | `GOOS=linux go vet` rc=0；`go test -c` rc=0；二进制含新名 `grep -a -c`=3；`//go:build`=0；新枚体内 `t.Skip`=0 | logs/02 |
| 3 | linux 实跑（Docker） | `golang:1.27` linux/amd64，仓 ro 挂 `/wisp` | r3 **SKIP**（逐字同 r6）、新枚**真跑 PASS**、notice 与 windows 同串；rc=0。logs/03 |
| 4 | 反证四发 | M1/M2（r6 的）＋M3/M4（自造新形） | 见 §3。logs/04 |

## 3 反证：转发的四种拼法，四发的红/绿

- **M1（r6 的：`"+raw+"`）** windows 红＝恰 {r3, 新}（全家谱 16 名其余 14 绿）；linux 红＝新（r3 仍 SKIP）＝r6 读数复现。
- **M2（r6 的：`strings.Join(Roots(), "、")`）** windows 红＝恰新枚一枚：禁词命中的是 **fold 后小写形**（verbatim `…\\001\\rootlist-leak-probe-174r6-a`），r3 臂绿＝"r3 只禁原拼法、fold 漏"复现。
- **M3（我造，新形：`fmt.Sprintf("%q", raw)`）** windows 红＝{r3, 新}（r3 靠引号里原样保留的 token）；linux 红＝新（`must not relay "\"\\u00a0\u2003\u2009\u205f\""`）。
- **M4（我造，新形：`strings.ToLower(raw)`）** windows 红＝**恰新枚一枚**，且 r3 臂绿的同时它的 notice 日志**明示漏了**（`…，c:\\…\\allowlist-leak-probe-174r3）` 全小写路径穿扫而过）；linux 红＝新。⇒ 新枚在"转发行为"上抓的档位 ≥ r3（M2 的根列表 fold、M4 的路径 fold 都只有新枚抓）。
- 第一次 M3 linux 跑成 rc=0 **是我自己的假阴性**：把 host 键 overlay 传进容器（no-op）——已识别、不记作读数、换 `/wisp` 键 overlay 重跑得红（logs/03、logs/04）。⚠ 如实登记这次自摆乌龙。

## 4 门禁（现跑；详情 logs/05）

build `./...` rc=0｜`GOOS=linux go vet` rc=0｜gofumpt v0.12.0 `-l` 空｜d22scan clean rc=0（`ban #8 internal/`=550）｜长度门 VERDICT GREEN rc=0（**60=60**，与派单"现 60=60"一致）。

## 5 风落（非 r6 的账；给编排者）

`bash scripts/d22scan.sh` 在 **GOTMPDIR 未设**时 rc=1：`tools/d22scan/runtests.sh:92: LOCALAPPDATAwisp: unbound variable`——`581dd252` 的"第四处折法"把 `${LOCALAPPDATA//\\//}` 改成了**单反斜杠形**（`"$LOCALAPPDATA\wisp-test-gotmp"`），在 `set -eu` 下反向吞进变量名；三处同类（`portable-tests.sh:155`／`winsec-tests.sh:69`／`wisp-cli-tests.sh:82`）都用**双反斜杠**、我逐字验过都能跑。仅影响**本机 windows**门禁路径（CI ubuntu 不进该分支）。**已闭合**：同窗口 `253-r6` 也撞见同源，编排者 14:49:32 于 `76b9e214` 单→双修复＋登记 `A929`；本腿**修复后复验**（GOTMPDIR 未设）⇒ `d22scan clean rc=0`。六形对照 repro 见 logs/05（仓外脚本 `C:\tmp\wisp174v2gotmp\parse-test.sh`），修复前后逐字对上。

## 6 交件与留盘件

| 文件 | 提交 |
|---|---|
| `docs/evidence/s1/ticket174-174-v2-accept-v1.md`（本腿判语） | 是 |
| 本件、`logs/01`–`05` 五枚 | 是 |
| `raw/`（18 份原始捕获：family 五发、docker 四份、d22scan 两向、gates 两段、sha 证等）、`mut/task-M3/M4.go`、`overlay-m{1,3,4}.json`、`overlay-m{1,3,4}-linux.json` | 留盘未提交（规则 8 只建不删） |
| `C:\tmp\wisp174v2gotmp\parse-test.sh`（风落诊断，仓外） | 留盘 |

复现（任一枚）：`gofumpt -l internal/tools/task_output_pointer_notice_test.go`；`go test -count=1 -run 'TestCanonicalizeErrorNoticeMustNotRelayTheFiledPath' ./internal/tools/`；docker 命令见 logs/03 首行；M4：在其上加 `-overlay=/wisp/.scratch/wisp/probes/174/v2/overlay-m4-linux.json`。

## 7 没做的（具名，不认领）

- **CI `test-core` 的真颜色**：未跑（推不了；docker ≠ ubuntu 宿主，如实区别）。下一发 CI 的判法已写进 logs/05（本对 SKIP=0、新名进 PASS 侧）。
- **整包 `-race`／全量 `portable-tests.sh`**：未跑（他腿同包；脚本无 dry-run）。
- **端到端 `wisp run`**：未跑（本机 DLL 存量事实，174-c1）。
- `internal/risk/**`、`cmd/**`、`docs/PLAN.md`：零字节、未评新内容。
