# 171-a1 台件：可复制命令序列（2026-09-27，锚 HEAD=d2e84f39，branch=dev）

## step 0
- date; git rev-parse --abbrev-ref HEAD -> dev; git rev-parse HEAD -> d2e84f39ccd1be7a21859d0e408e5febac1ff5c9
- git status --porcelain -- tools/d22scan/ docs/PLAN.md docs/specs/ internal/ cmd/  -> 空

## 名册来源（三档）
1. 语料自登记号（派单 §1 指定的两处）：
   - { cat .scratch/wisp/issues/161-*ever-rung-done.md; cat docs/evidence/s1/161-*.md; } > corpus(2842 行)
   - grep -ohE '\b[0-9a-f]{7,10}\b' corpus | sort -u  -> 58 枚 token
   - 每枚 git cat-file -t：41 行是 commit（去重同 commit 不同长度 7 行 -> **34 枚不同提交**）；17 枚非 commit（123a124 是 diff hunk、20260926 是日期、80471c2/e7c41c2e/5eacc332/39ca85a4/4e595dbd/268f25ab/1322f4a/01461d0/071522cd/1daa1542/1402817a0/1405e7800/1405e9560/1406150c0 是日志/数值噪声）-> 取数前即拒，未进任何名册
   - 34 枚中 6 枚属他票被语料点名：15ff2be(票77) 5afa666(前端log) fbe12c7(r1锚,前端log) 611ae8b(票169) 506cbae(票158) 4267bb3c(票164) -> 161 系自登记 = **28 枚**
2. git log 补登记（写了 161 但没在语料自登记号）：
   - git log --format='%h|%ad|%s' --date=... dev | grep 票161 系 -> 另得 13 枚：0492cd5d cd999794 8995b366 feed0a17 22ecf9db 71cff1cc 0d176473 70a106e5 339b67a1 bc65d975 a7f2d790 1a95eaa6 e0405d57
   - 28+13 = **41 枚**（161 系全集）
3. 本程复算用的最宽口径 = 28+6+13 = **47 枚候选**（含他票，逐枚标归属）
- 区间定义测过：fbe12c7..1a95eaa6=57｜fbe12c7..e0405d57=59｜27f798e..e0405d57=51｜fbe12c7..30121fe=51 -> 无一=37，且区间内混入票 164 族 4 枚 PLAN.md 提交(2871f75 824bb9c ced72f8 4267bb3c) ⇒ 区间口径下"0 命中"本身不成立

## 命中与 numstat
- git log --no-walk --reverse --name-only --format='@C %h %s' <47枚> > roster-files-raw.txt(463 行)
- awk 双判据（冻结 11 项展开 + r1 特别名单 3 项）-> per-commit-hits.txt
- git log --no-walk --reverse --format='@N %h %s' --numstat <47枚> -- tools/d22scan/main.go -> 唯一一行 "12	0	tools/d22scan/main.go"，出处 3e56e680；cat -A 核过删除列是真 0（`12^I0^I...$`，非空行折叠）

## 冻结 11 项在 HEAD 的展开（逐项链数）
docs/PLAN.md=1｜docs/specs=14｜internal/risk=37｜internal/panel=20｜internal/agent/approval=18｜thresholds.go=1｜golden(忽略大小写)=60｜allowlist.txt=1｜scripts/slo-check.ps1=1｜frontend=85｜design=30 -> 合计 268 枚具体文件；11 项指法枚数与票面 :29 逐字一致

## 门禁（只读两把）
- 12:04:23 sh scripts/d22scan.sh -> rc=0（"clean - no D22 ban violations"）
- 12:04:51 bash tools/d22scan/runtests.sh -C tools/d22scan ./... -> rc=0，PASS=34 FAIL=0 SKIP=0，=== RUN=76，[no tests to run]=0
- 名册两向 comm：对 probes/161/r4-pre-runtests.log.roster.txt(34 行) 按测试名（第 3 域、剥耗时括号）比对 -> roster-comm-nameonly.txt = **空**（34 vs 34 同名）
- 第一发整行 comm(roster-comm-diff.txt) 68 行全差 = 本程格式伪影（行尾耗时+空格未剥），留件为证，判定以 name-only 那发为准
