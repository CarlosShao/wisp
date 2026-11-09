# gate-clauses.sh 全树爆限读数（2026-11-08 23:4x，写腿普查腿留存；零产码、零票改、未立案、只建不删）

发现：共享门禁仪器 `.scratch/wisp/probes/154/gate-clauses.sh` 在当前树上对全编队是坏的。
未立案（票 178 已 -done，其"随单发现①"预言了这次长肥）⇒ 按未定义即停本腿不动仪器，只留读数。

## 读数（本机现跑，全部可重跑）

1. `sh .scratch/wisp/probes/154/gate-clauses.sh` ⇒ 退出码 **10**，`腿数＝14 声明与实测不符＝10`；
   G1/G1b/G2/G5/G5pos/G5neg/G6/G6neg/G7/G7neg 十腿全部 `实测=git-grep-rc-126`（在 10-14 票 178 收口时
   编排者亲跑净树尺为 rc=2；10 月中下旬历次腿读数为 rc=0/2，见各票面）。
2. 旁证 stderr：`line 587: ring: command not found`、`line 590: 改完之后: command not found`
   ＝脚本 587/590 行双引号 echo 内含未转义反引号（`ring 1`／`改完之后 AC#2 那一发仍会红`）被 bash
   当命令替换执行＝纯噪音，与 rc=126 同件不同因。
3. 机制（`bash -x` trace 实锤）：脚本第 67-69 行 `GO='internal/**/*.go'`／`GO2='cmd/**/*.go'`，
   调用点对 `$GO $GO2` **未加引号** ⇒ bash 在工作树上做路径名展开，把 glob 摊成
   **逐枚文件字面 pathspec** 再喂 `git grep <锚> --`。现量展开规模：
   `git ls-files 'internal/**/*.go' 'cmd/**/*.go' | wc -l` ＝ **1209**；
   单行字符数 ＝ **45226** ＞ Windows CreateProcess 32767 命令行上限 ⇒ exec 失败 ⇒ bash rc=126。
   票 178-r1 收单时同展开只有 **496 枚路径**（其随单发现①原文："`-- $@` 未引号⇒射程文件清单来自工作树
   （现量 496 路径，合成树新加文件不可见）"）——10-14 之后树长肥过限，仪器从"能跑"变"全红"。
4. 正控（同命令、glob 形加引号直接喂 git grep，不经展开）：
   `git grep -nEw 'OpenTask|CloseTask' 34cfd4f402cb3d1a44bdaa9c9afc09d89555fdbf -- 'internal/**/*.go' 'cmd/**/*.go' ':!*_test.go' ':!internal/tools/bridge.go'`
   ⇒ rc=0 有命中。⇒ 引号修法可行，但**改仪器须留凭据＋非实现者裁**（178-v1 判语先例），本腿不动。
5. 影响面：凡判据引用这把尺的票——195 AC#5（BAD 腿名册只 G6neg 那句判据照此读数现量不可复现）、
   346 族、154/161/171 家族、以及任何写腿拿它当门禁的那一发，今天都会读到 10 腿 BAD 假红。
6. 锚注：脚本头 `锚点 34cfd4f402cb…` ＝ `git rev-parse HEAD` 在 23:36 现量（即跑尺那一刻的 HEAD），
   不是死锚；与本发现无关。

## 建议归口（编排者裁，本腿不执行）

新立案一枚"仪器程"：修法＝调用点 `"$GO" "$GO2"` 加引号（git grep 原生吃 glob pathspec，正控已证），
顺带把 587/590 两行 echo 的反引号转义或改单引号；凭据＝本件读数＋修后四把（净树 rc 回 2／
AC#2 变异发仍红／G6neg 仍 quiet／flip-declaration.sh 不受影响）；裁决者≠实现者。
