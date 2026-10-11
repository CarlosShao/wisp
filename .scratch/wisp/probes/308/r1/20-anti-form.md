# 308-r1 · 20-anti-form（`AC#2` 反形敏感，三件一起做，缺一不算交）

尺一律"先落文件再断言、`finally` 还原"：突变⛔ 留在跟踪文件里。⛔ `cmd | head; echo $?`（那取 head 的退码）——每发自己 `> logs/*.txt 2>&1; rc=$?` 直取。⛔ 行号锚开刀（本仓实测：夹具一插就打死锚位）⇒ 全部认**内容锚**。

---

## (a) 拔掉守卫 ⇒ 用例回到 panic（红句逐字）

突变＝删去 `driveShare` 那三行 `if len(vol) == 0 { return "", false }`；改后函数体只剩 `return strings.ToLower(vol[:1]) + "$", true`。尺＝`grep -n 'if len(vol)' internal/risk/pathresolver_volumeless_spelling_308_test.go` ⇒ 命中 `line 24` 是**还原后**的读数（突变那发此尺 0 命中）。

跑（只打新用例，最小面积，⛔ 影响其它件）＝`go test ./internal/risk/ -count=1 -timeout 120s -run TestVolumeInvariantSpellingsVolumeLess308 -v > logs/anti-a-mutation-broken.txt 2>&1; rc=$?`。

读数（逐字，件＝`logs/anti-a-mutation-broken.txt`＋`logs/anti-a-mutation-broken.rc.txt`）：
- `rc=1`
- `--- FAIL: TestVolumeInvariantSpellingsVolumeLess308 (0.00s)`
- **`panic: runtime error: slice bounds out of range [:1] with length 0 [recovered, repanicked]`**（与票 308 §现量 引的历史 `72/a2/logs/t03-broken-surface-note.txt` 那发**逐字同形**）
- 崩位栈＝`internal/risk.driveShare(...)` `pathresolver_volumeless_spelling_308_test.go:26` → `TestVolumeInvariantSpellingsVolumeLess308` `:82`（行号锚＝突变现场的栈行，⛔ 用作后续锚位）

⇒ 牙确实挂在这枚 `len(vol)==0` 守卫上。⛔ 挂在机器台面：本机三把 temp 台面（default/alias/posix-rewritten）⛔ 能现场触发无卷标那形（裁语 ⓑ、`probes/308/a1/90-final.md` §1 表）。

---

## (b) `cmp`／`sha256` 逐字节还原 ⇒ 复绿

还原＝把 §(a) 那三行原样写回。尺（两把并排）：

- `grep -n 'if len(vol)' internal/risk/pathresolver_volumeless_spelling_308_test.go` ⇒ 24 命中（守卫在位）。
- `sha256sum internal/risk/pathresolver_volumeless_spelling_308_test.go` ＝ **`f2c0130702384bfd0c60749cd3163d7492fd5351f64ca32ce4529ca04897070e`**（还原后），与快照件 `logs/anti-b-restored-good.go.txt`（还原后 `cp` 出来的对照）**同 hash**。
- 复绿尺＝`go test ./internal/risk/ -count=1 -run TestVolumeInvariantSpellingsVolumeLess308 -v` → `rc=0`、`--- PASS: TestVolumeInvariantSpellingsVolumeLess308`、`grep -cE 'panic:|slice bounds' logs/anti-b-restored-green.txt`＝0。
- 全⚑ 面复尺＝`go test ./internal/risk/ -count=1 -v` 逐名⇄前次 `logs/post-default-names.txt`（`diff` 空、`IDENTICAL-AFTER-RESTORE`，尺＝`logs/post-restored-default-full.txt`），⛔ 一枚红多、⛔ 一枚红少。
- 跟踪件⛔ 残留：`grep -c 'TEMPORARY MUTATION' internal/risk/pathresolver_volumeless_spelling_308_test.go`＝0（那条注释＝(a) 阶段临时标记，(b) 已剥）。

---

## (c) ★ 别名台面（8.3 那把）的颜色⛔ 因这枚守卫而变

台面（⛔ 我造，⛔ 删里⛔ 任何东西）＝`TEMP/TMP/TMPDIR=C:\Users\swq\tmp\RISKAT~1\RUNNER~1`（长形 `…\riskattrib1\RunnerTempAreaForAttribution1`，本机既有）。绕开 MSYS 路径改写＝用绝对 `cmd.exe` 起子进程（本仓 `cmd` 名被 `%APPDATA%/npm/cmd` 那枚 npm shim 遮蔽，`which cmd` 见 `logs/`；尺＝`/c/windows/system32/cmd.exe //c "set … && go test …"`）。

改前⇄改后逐名作差（两把台面各一枚，⛔ 混层级）：

| 台面 | 前名册（`^--- (PASS|FAIL|SKIP):`） | 后名册（同尺） | 作差 |
|---|---|---|---|
| default | `logs/pre-default-names.txt`＝132 行（131 PASS + 1 既有 SKIP `TestSyncRegistryProbeLive`） | `logs/post-default-names.txt`＝133 行 | **＋1 枚**＝`PASS: TestVolumeInvariantSpellingsVolumeLess308`；⛔ 少⛔ 名、⛔ 变色 |
| **alias**（8.3 那把） | `logs/pre-alias-names.txt`＝132 行（同上） | `logs/post-alias-names.txt`＝133 行 | **＋1 枚**＝同一枚新用例；⛔ 少⛔ 名、⛔ 变色 |

尺＝`diff logs/pre-{default,alias}-names.txt logs/post-{default,alias}-names.txt`（各只输出一行 `126a127 > PASS: TestVolumeInvariantSpellingsVolumeLess308`）。

两把台面**改后⇄改后**⛔ 分⛔ 颜（尺＝`diff logs/post-default-names.txt logs/post-alias-names.txt`＝空、`SURFACES-IDENTICAL`）。

子层级⛔ 变：两枚台面⛔ 前⛔ 后 `grep -cE '^    --- PASS'`＝68；`^    --- FAIL`＝0。⛔ panic／⛔ fatal／⛔ `test executable failed`：⚑ 面 6 份（`pre-{default,alias}`＋`post-{default,alias}`＋`post-restored-default`）每份 `grep -cE 'panic:|fatal error:|test executable failed'` 均＝0。

⇒ **守卫⛔ 关掉⛔ 用例**：三枚崩点用例（`TestClassifySpellingInvariance`／`TestPathResolverUNCAListDenied`／新单元用例）在两把台面都仍 `PASS`；helper 的 ¬ok 分支⛔ 触发（本机制，尺＝`grep -c 'dropped the two volume-derived' logs/post-*-full.txt`＝0，两⚑ 皆 0）。

---

## 三件缺一⛔ 算交 → 三件齐

- (a) ✓：`logs/anti-a-mutation-broken.txt`＋`.rc.txt`（`rc=1`、panic 逐字、栈位具名）。
- (b) ✓：`logs/anti-b-restored-good.go.txt`（sha256 对照快照）＋`logs/anti-b-restored-green.txt`（复绿）。
- (c) ✓：⚑ 面名册改前⇄改后逐名作差，两⚑ 差集同（只 ＋1 枚新用例），⛔ 因守卫失⛔ 名。
