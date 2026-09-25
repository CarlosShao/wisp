# panel-l2-grant-nail-fix-r5-accept-r1 —— 非实现者对"门钉 r5"这一批的对抗验收

> 被验物：`internal/panel/l2_grant_boundary_test.go` @ 锚点 **`aeba6ff`**（blob `f87301df6be0ba1ceb00cba7f9ddda877df521c6`，2380 行）
> ＋自述件 `docs/evidence/s1/panel-l2-grant-nail-fix-r5.md`（§0–§8）。
> 任务来源：`docs/evidence/s1/panel-l2-grant-nail-fix-r4-accept-r1.md` §12 点名的 **F-R4-1／F-R4-2／F-R4-3**。
> 本程＝**验收程**：被验物一字未改，只 commit 本件；不勾票面／台账；不 push。
> 本程全部变异落在**仓外副本**（§1），原始读数在 `D:\tmp\wisp-r5acc1\out\`（§14 列路径）。

---

## §0 锚点、口径与"读数出在哪台件"——先把这条写死，因为工作树此刻在别人手里

**简报的硬要求**：另一程正在给 `internal/panel/` 加快照泵（会动 `composer.go`／`bridge.go`）⇒ 所有读数取锚点归档副本，
脏工作树的读数单独标注、且**不作任何一格的判据**。本程执行到位，且这条不是形式：我 15:5x 建副本、
16:0x 复量时工作树已经长成这样（逐字 `git status --porcelain`，只列与本包相关的）：

```
建副本时（15:56）  ?? 无；internal/panel/ 干净
复量时（16:04）    M  cmd/wisp/run.go
                   ?? cmd/wisp/panel_pump.go
                   ?? internal/agent/approval/pending_read.go
                   ?? internal/panel/pump.go
                   ?? internal/panel/pump_test.go
⇒ 本程正文里每一枚"四数／红名／asked="都出自 D:\tmp\wisp-r5acc1\{anchor,pre,witnessonly}（§1）；
   工作树口径的读数只有 §8 表里那一行，标着"工作树"，并且它不参与任何一格的裁。
```

```
进场 HEAD = aeba6ff（＝锚点本身）；交件前复量 HEAD = ddf8b6f（共享树在漂，只有 docs 两枚与泵那一程的未提交件）
git rev-parse aeba6ff:internal/panel/l2_grant_boundary_test.go
   -> f87301df6be0ba1ceb00cba7f9ddda877df521c6（2380 行）＝ fcce0d3 那一枚，与实现件 §0 同值   ✅
git rev-parse aeba6ff:internal/panel/bridge.go   -> d2cd6362ecc6941a0cee58073a2bf8ab6669a590
git rev-parse 88eab34:internal/panel/bridge.go   -> d2cd6362ecc6941a0cee58073a2bf8ab6669a590  ⇒ 产品码零字节（跨 r2..r5）
git rev-parse aeba6ff:internal/panel/composer.go = ee2a92d 同值 -> e444d5d117a073c393d9b43f7e3e98e8ede61755
git diff --numstat ee2a92d aeba6ff -- internal/panel/  ->  185   25   internal/panel/l2_grant_boundary_test.go（本包只有这一枚）
git diff --name-only ddea3c3 aeba6ff                  ->  本件之外只有 docs/**（⇒ 锚点相对 ddea3c3 的代码差＝那三枚写死的数）
三枚 commit 的 pathset（逐枚 git show --numstat，不跑区间）：
   ddea3c3  149  24  internal/panel/l2_grant_boundary_test.go
   fcce0d3   36   1  internal/panel/l2_grant_boundary_test.go
   a72c110    2   0  docs/evidence/s1/panel-l2-grant-nail-fix-r4.md
```

**F-R4-3 那一枚文档账我单独核过**：`git diff ee2a92d aeba6ff -- <那枚 r4 件>` 的输出是 `449a450,451`（纯追加），
`## §8` 标题仍在且只有一枚，追加的那一行以 `>` 起头、内容点名"另两枚因子清单自己可被剪窄"、引它自己那两枚发现量
（`asked=462`／`asked=352`，两枚我本程都独立复算到了，见 §3/§4）、出处与去向都对得上 ⇒ **F-R4-3 结清**。
⚠ 一处**过期风险**记在这里（不改判）：那一行写的是"各删任意一枚时全包全绿"，那是 **r4 交付态**的现状；
r5 之后同一发已红（§4 的 K1 十一发）。该句末尾指了 `panel-l2-grant-nail-fix-r5.md §1、§2` 作去向，
所以不是错句，但下一程若只摘那一行会读到过期现状——本程在 §10 把它列成一句口径提醒。

---

## §1 台件自证（ⓐ）：三条仪器坑本程都自己走了一遍，不复用实现者的脚本

```
正路：git -c core.autocrlf=false -c core.eol=lf archive aeba6ff | tar -x -C /d/tmp/wisp-r5acc1/anchor
逐枚比字节（本程自己的尺）：git ls-tree -r aeba6ff（1174 枚）对比 git hash-object --no-filters <副本>
   -> match=1168  mismatch=6
   那 6 枚逐字：scripts/build.ps1 · scripts/dev/ball-cycle.ps1 · scripts/fetch-deps.ps1
               scripts/sign-models.ps1 · scripts/slo-check.ps1 · scripts/spike/run.ps1
   成因＝.gitattributes:2 逐字 `*.ps1 text eol=crlf`（显式属性压过 core.eol=lf）⇒ 仓自己要求的 CRLF，不是污染；
   grep -rn "\.ps1" internal/panel/ -> 0 枚 ⇒ 本包没有一枚测试读 .ps1 ⇒ 这 6 枚不进本程分母。
被验物逐枚（--no-filters）：
   l2_grant_boundary_test.go -> f87301df6be0ba1ceb00cba7f9ddda877df521c6 ＝ 锚点 ✅
   bridge.go                 -> d2cd6362ecc6941a0cee58073a2bf8ab6669a590 ＝ 锚点 ✅
   composer.go               -> e444d5d117a073c393d9b43f7e3e98e8ede61755 ＝ 锚点 ✅
   frontend/fixtures/composer-states.html -> ab39390c…（6179 字节，CR 计数 0）＝ 锚点 ✅

坑①（本程独立复现，第三程）：只给 core.autocrlf=false、不给 core.eol=lf ——
   git -c core.autocrlf=false archive aeba6ff frontend/fixtures/composer-states.html | tar -x
   -> 680b908b…，6187 字节，CR 计数 8  ⇒ 与锚点不同值。本程副本用了正路，未踩那枚假红。
坑②：比字节一律带 --no-filters（上面每一枚都是这么取的；不带它会把 680b908b 读成 ab39390c）。
坑③：临时件只建不删——本程零枚 rm/rmdir/rm -rf（路径清单在 §14，交编排者一次清）。
```

变异台件：`D:\tmp\wisp-r5acc1\scripts\mut.py`（本程自己写的，未复用实现者的 `decoys.py`）。
每发施加前先把两枚文件从 `pristine\<副本>\` 复原；每一刀的 `old` 串在当前文本里**计数不等于 1 就抛 REFUSED、半行不落**
（本程真实拒过 0 次，但同一把尺挡下过我自己两枚坏切法，见 §14 失手登记）；每发结束再复原一次。
末次全树复算（87 发跑完之后）：`match=1168 mismatch=6`，那 6 枚仍是上面点名的 `*.ps1` ⇒ **副本没有被任何一发变异留脏**。
读数总枚数：`out/` 里 **99 枚**（87 发 `go test` ＋ 三门 ＋ 名册/字节比对件）；其中 **4 枚无效**（`H01–H04`，
RUN=0 编译不过，本程第一版"拆正控"的切法漏了变量声明），一律留档不删、正文改引修正后的 `H05–H07`。
判红绿只认 `^--- FAIL` / `^    --- FAIL`（`t.Logf` 也带 `file:line:` 前缀，不作判据）。
