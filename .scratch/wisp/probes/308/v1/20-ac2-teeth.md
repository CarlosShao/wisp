# 308-v1 · 20-ac2-teeth（`AC#2` 反形敏感三件＋恒真／空转／缺世界三族判定）

突变纪律＝**先落快照件、`trap` 里还原、当场 `cmp`＋`sha256sum`**；驱动器⛔ 在仓内（写在本腿 `mktemp -d` 的临时目录里，件名见 `logs/teeth*`）。全程 `internal/risk` 工作树在每次脚本结束时 `git status --porcelain` ＝空（已复跑坐实）。

## (a) 拔掉守卫 ⇒ 回到 panic，红句逐字

- 突变＝`sed -i '/^\tif len(vol) == 0 {$/,/^\t}$/d'`（内容锚，⛔ 行号锚）→ 落地尺 `grep -c 'if len(vol) == 0'` ＝ **0**（件 `logs/teeth-mutation-landed.txt`，⛔ 未落地的突变⛔ 算凭据）。
- 只打新用例：`rc=1`，`--- FAIL: TestVolumeInvariantSpellingsVolumeLess308 (0.00s)`，
  **`panic: runtime error: slice bounds out of range [:1] with length 0 [recovered, repanicked]`**（件 `logs/teeth-mut-target.txt`）
  ⇄ 票 §现量 引的历史原始件 `probes/72/a2/logs/t03-broken-surface-note.txt:7` **逐字同形**（本腿把两件的字面串比对过，⛔ 只说"同族"）。
- ★**整包级形状**（这条是票 308 的存在理由，本腿量到的是**这一发**）：同一枚突变下打**整包** `go test ./internal/risk/ -count=1 -v` →
  `rc=1  RUN=35  PASS=27  FAIL=1  SKIP=0  PANIC=1`（件 `logs/teeth-mut-fullpkg.txt`）
  ⇄ 守卫在位的同一包同一尺＝`rc=0 RUN=201 PASS=132 FAIL=0 SKIP=1 PANIC=0`（`logs/v0-full-worktree-guard-intact.txt`）。
  ⇒ 一枚无守卫切片带走 **201−35＝166 枚 RUN / 132−27＝105 枚 PASS 的求值**，且名册⛔ 报错、只看 rc 会读成"红了两枚"。**"看着像红、其实是这一发没跑完"被本腿当场重造出来又当场收回去。**

## (b) 逐字节还原 ⇒ 复绿

- `cmp -s 快照 文件` → `cmp_after_restore_rc=0`；`sha256sum` 两侧同＝**`f2c0130702384bfd0c60749cd3163d7492fd5351f64ca32ce4529ca04897070e`**（件 `logs/teeth-cmp.txt`、`logs/teeth2-cmp.txt`、`logs/teeth3-cmp.txt` 三发同 hash）
  ⇄ `308-r1` 自报的还原 hash（`probes/308/r1/20-anti-form.md:28`）＝同一串 ⇒ 交付物此刻与 r1 落地时**逐字节同一**（单位＝sha256，射程＝工作树文件）。
- 还原后：单用例 `rc=0`（`logs/teeth-restored-target.rc.txt`）、整包 `rc=0 RUN=201 PASS=132 FAIL=0 SKIP=1`（`logs/teeth-restored-fullpkg.*`）；teeth2/teeth3 收尾复绿各 `rc=0`。

## (c) ★别名台面（8.3 那把）的颜色⛔ 因守卫而变

- 台面＝`C:\Users\swq\tmp\RISKAT~1\RUNNER~1`（⇄长形 `…\riskattrib1\RunnerTempAreaForAttribution1`，两者盘上都在）。⛔ 我造⛔ 我删。绕开 MSYS 改写＝`/c/windows/system32/cmd.exe //c "set TMP=…&& go test …"`。
- **先证台面真的用上了**：同一 env 下仓外探针打印 `os.TempDir()="C:\\Users\\swq\\tmp\\RISKAT~1\\RUNNER~1"`（件 `logs/v15-cmd-grid.txt`），且两发名册里各有 **10** 处 `RISKAT~1` 字样（件 `logs/v40`/`v41`）。
- 改前⇄改后逐名作差（尺＝顶层名册排序后 `diff`）：`132 → 133`，差＝**只 `+ PASS TestVolumeInvariantSpellingsVolumeLess308`**（件 `logs/v42-alias-roster-diff.txt`）；子层级 `^    --- PASS` 68⇄68；两发 rc 均＝0、`panic|slice bounds` 命中 0。
- 改后**两把台面互比**＝`diff v41-post-alias v0-worktree` → **空**（`rc=0`）；改前两把台面互比也＝空 ⇒ 别名台面⛔ 因这枚守卫失名、变色、少跑。

## 两道必需的正控制（⛔ 一道当两道用）

1. **"这发必须红"确实红**：突变那发 `rc=1`＋panic 逐字（§a），且落地尺＝0 命中。
2. **"pristine"确实绿**：突变之前同一命令序列第一发 `rc_pristine=0`（件 `logs/teeth-pristine.rc.txt`），整包基线 `rc=0`。

## 恒真／空转／缺世界三族（逐族读数）

| 族 | 问的是哪一句 | 读数 | 判 |
|---|---|---|---|
| (a) 永远红 | 断言是否骑在一枚⛔ 能命中的匹配器上 | 全绿基线 `FAIL=0`；`dropped` 那三行判据在突变(见下 M2)下**由绿变红并指名** ⇒ 匹配器可达 | ⛔ 恒真形 |
| (b) 空转绿 | 循环集合是不是空集 | `volumeFree` 那圈＝写死 7 枚（打印＝7）；`for name := range keep` 实测 9／7；**M3b 抽掉 `"lower"` ⇒ 落地尺 1→0、`rc=1`、红句逐字 `volume-free spelling "lower" went missing: only the two volume-derived UNC spellings may be dropped (AC#1)`（`…_test.go:105`），且打印分母自己从 7 掉到 6** ⇒ 集合非空、断言活 | ⛔ 空转 |
| (c) **缺世界** | 分支要不要第二个值才分叉、而⛔ 任何夹具给得出 | ★**有**：`TestClassifySpellingInvariance` 的 `len(dropped) > 0` 一支与 `TestPathResolverUNCAListDenied` 的 `t.Fatalf("…cannot be constructed")` 一支，**在本机能造的每一把台面上到达数＝0**（哨兵 `dropped the two volume-derived`／`cannot be constructed` 在 v0/v13/v41/v26 四发里命中全＝0）。`tmpHome`→`t.TempDir()` 的卷标由 Win32 `GetTempPath` 决定，本腿 env 网格证其⛔ 交⛔ 卷标路径（见 60 件）⇒ 这两支的"红得可解释"目前**只有代码形状＋helper 输入面一枚载具** | **成立，具名 `Q-308-v1-B`** |

- ★另一枚 M2 反证（⛔ 只靠 panic 变红的凭据⛔ 算数）：把 `volumeInvariantSpellings` 的 `return keep, []string{"unc-admin","unc-extended"}` 一行删掉（helper ⛔ 再报告去掉）⇒ `rc=1`，红句＝`dropped=[], want exactly [unc-admin unc-extended]`（`…_test.go:95`），**⛔ panic** ⇒ 判据本身有牙，⛔ 只是"顺带崩了"。
- ⚠ 本腿自己踩到的一枚假凭据具名留下：第一次 M3 用 `sed` 抄了错误空格对齐，**突变没落地**（`grep -c '"lower":'`＝1）却报 `rc=0` ⇒ 那一发**⛔ 是任何证据**；件 `logs/teeth2-M3-landed.txt` 原样保着，随后 `logs/teeth3-M3b-landed.txt`（1→0）才是那发数的正身。

## 本格判语

**成立**（三件齐、两道正控制齐、另加 M2/M3b 两枚反空转反证）。`AC#2(a)` 的红句与历史原始件逐字同形；`AC#2(b)` 的 sha256 与 r1 自报值同一；`AC#2(c)` 别名台面作差＝＋1。
唯一⛔ 被这三件覆盖的东西＝上面 (c) 族的"windows-tagged 两支⛔ 到达"，它⛔ 反证⛔ 是交付缺陷，⛔ 归 `AC#2` 的账，归 `Q-308-v1-B` 的欠读数。
