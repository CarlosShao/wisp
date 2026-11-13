# 证据件 — 票 266 r2：`slo-nail-census.sh --self-test` 在 CRLF 签出下 14/32 挂（pristine 建源未规范换行；1 行修）

**腿**：266-r2（写腿，第二次派单；第一次返回空，根因由编排者 docker 实验坐实）
**时刻**：2026-11-13 17:0x +08　**基点 HEAD**：`1d689498d`（`dev`）
**写面**：`scripts/slo-nail-census.sh` 一处（原 :438）＋本件＋票 266 腿尾 bullet。
⛔ 阈值 / golden / `allowlist.txt` / `.gitattributes` / 契约：零字节（`git diff` 可复现）。

## 根因（承接派单，本腿逐条复核为真）

1. `.gitattributes` 给 `scripts/slo-check.ps1` 设 `eol=crlf` ⇒ CI（ubuntu）签出时目标是 CRLF。
   本腿现跑 `git check-attr eol text -- scripts/slo-check.ps1 scripts/slo-nail-census.sh`：
   `scripts/slo-check.ps1: eol: crlf`＋`text: set`；`scripts/slo-nail-census.sh: eol: lf`（⇒ control 12 的 bench_census sed 不受影响，与派单预判一致）。
2. 自测全部 sed 变体以 `/^X$/` 锚定行尾删/改行（:458、:472、:541 等）；CRLF 行带 `\r` ⇒ 锚点永不匹配 ⇒ 变体根本没改 ⇒ census 恒绿 ⇒ 所有「必须响」的对照挂。
3. 本机 Windows 工作树恰好 LF（本腿现跑字节统计：`slo-check.ps1` CRLF 行数=0）⇒ 本机恒 32/32，红只出现在 CI。

## 修法（1 行码 + 6 行注释）

- 原 `:438`：`    cp "$target" "$BENCH/pristine.ps1" || cannot_decide "cannot copy $rel into the bench"`
- 现 `:444`：`    tr -d '\015' < "$target" > "$BENCH/pristine.ps1" || cannot_decide "cannot copy $rel into the bench"`
- 机制：`$pristine` 是自测全部 sed 变体的建源；建时就规范成 LF，锚点在 LF/CRLF 两种签出下都能匹配。census 自己读码本就 `tr -d '\015'`（:187）⇒ bench 目标是 LF 毫无影响。
- ★ **CRLF 对照一字未动**（现打印为 control 27，`:523-528`）：它单独用 `awk '{ printf "%s\r\n", $0 }' "$pristine"` 重建 CRLF 副本、期望绿——专门钉「census 本身能吃 CRLF」，修后照常过（见读数 #4 末两行）。

## 验证读数（工具：`D:\work\soft\Git\bin\bash.exe` 5.3.9 cygwin；docker `ubuntu:22.04`；docker 命令＝派单原文那条）

| # | 验证 | 读数 |
|---|---|---|
| 1 | 本地 LF 树自测（**改前基线**） | `slo-nail-census: positive control PASSED - 32 control(s) ran, 0 failed`，rc=0 |
| 2 | **docker CRLF 复现（改前）** | `slo-nail-census: positive control FAILED - 32 control(s) ran, 14 failed`——与 CI 读数一字不差 |
| 3 | 本地 LF 树自测（改后） | `slo-nail-census: positive control PASSED - 32 control(s) ran, 0 failed`，rc=0 |
| 4 | **docker CRLF 复现（改后）** | `CONTROL_FAILED_COUNT=0`；`slo-nail-census: positive control PASSED - 32 control(s) ran, 0 failed`；control 27 `a CRLF copy of the target stays green (rc=0)` 仍过 |
| 5 | 普通 census（无 `--self-test`，对真目标） | `slo-nail-census: PASS - 2/2 nails declared here are defined once, carry their own ok sentence, are called once at column 0, and run before the first sample at scripts/slo-check.ps1:533`，rc=0 |
| 6 | `bash -n scripts/slo-nail-census.sh` | rc=0，无输出 |
| 7 | `sh scripts/d22scan.sh` | `d22scan: clean - no D22 ban violations`，rc=0 |

改前 docker 逐字尾段（`tail -6`）：

```
slo-nail-census: CONTROL FAILED - the red names the ordering token - expected the output to carry [slo-nail-census-ordering]; output at /tmp/tmp.87OJUXbrWw/out-30.txt
slo-nail-census:   (bench kept at /tmp/tmp.87OJUXbrWw)
slo-nail-census: CONTROL FAILED - the red says why the position matters - expected the output to carry [not before the first sampling call]; output at /tmp/tmp.87OJUXbrWw/out-30.txt
slo-nail-census:   (bench kept at /tmp/tmp.87OJUXbrWw)
slo-nail-census: positive control FAILED - 32 control(s) ran, 14 failed
slo-nail-census: positive control RED - the census is not run, because a ruler that cannot go red does not get to print a reading
```

改后 docker 逐字尾段（`tail -6`）：

```
slo-nail-census: control 28 ok - a count that no longer matches its tables turns the census red (rc=1)
slo-nail-census: control 29 ok - the red names the drift token - output carries [slo-nail-census-declaration-drift]
slo-nail-census: control 30 ok - calls moved below the first sample turn the census red (rc=1)
slo-nail-census: control 31 ok - the red names the ordering token - output carries [slo-nail-census-ordering]
slo-nail-census: control 32 ok - the red says why the position matters - output carries [not before the first sampling call]
slo-nail-census: positive control PASSED - 32 control(s) ran, 0 failed
```

（改前挂的 14 个对照与派单引文一致：deleting both nail calls、deleting a nail's definition、the red names the missing definition token、the ordering token、not before the first sampling call 等——全是「删行变体没生效」。）

## 诚实档

- docker 改前读数那行 rc 报 0 是管道尾 `tail` 的返回码；真判定以屏上 `PASSED/FAILED - … failed` 行为准（改前＝FAILED 14，改后＝PASSED 0）。判定行已逐字贴上。
- 中途一次 d22scan 红（`missing go.mod`）是**我自己的验证脚本**把 `MSYS2_ARG_CONV_EXCL="*"` 泄漏给后续进程、弄坏传给 Go 工具的路径所致——环境污染，与 1 行修复无关；清环境重跑即 clean（两次原始读数存 `D:\tmp\wisp\20261113\266-nailcensus-r2\d22scan-postfix.txt`（脏）与 `d22scan-postfix2.txt`（净），只建不删）。
- 正控 bench 目录按脚本自陈 kept on disk（本机 `/tmp/tmp.*`、容器内随 `--rm` 丢弃）；本腿临时件全落 `D:\tmp\wisp\20261113\266-nailcensus-r2\`。
- 边界残留（不归本票、票面结票行已登记）：原地缩声明仍 rc=0。

## 提交

修复与本件同落一枚提交；两枚提交的 hash 记于票 266 腿尾 bullet（复现尺：`git log --oneline -- scripts/slo-nail-census.sh docs/evidence/s1/266-slo-nail-census-crlf.md`）。
